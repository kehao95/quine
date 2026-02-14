package runtime

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

// Runtime orchestrates the agent's execution loop.
type Runtime struct {
	cfg        *config.Config
	provider   llm.Provider
	sh         *tools.ShExecutor
	fork       *tools.ForkExecutor
	exec       *tools.ExecExecutor
	tape       *tape.Tape
	tapeWriter *tape.Writer
	tools      []llm.ToolSchema
	semaphore  *Semaphore
	startTime  time.Time
	log        func(format string, args ...any) // operational log → log file
	logError   func(format string, args ...any) // failure signal → stderr

	// originalInput stores the user's input for this session.
	// Needed for exec to preserve the mission.
	originalInput string

	// stdout/stderr writers (overridable for testing)
	stdout *os.File
	stderr *os.File

	// logFile is the dedicated operational log file (§10.2).
	// Operational (INFO/DEBUG) messages go here, keeping stderr pure
	// for the Agent's semantic gradient (failure signals only).
	logFile *os.File

	// panicMode is set by SIGALRM (§2.2). When set, the next turn injects
	// a "System 1 Override" message forcing the agent to exit immediately.
	// Non-exit tool calls are rejected while in panic mode.
	panicMode atomic.Bool

	// activeProcess tracks the currently running tool subprocess (§2.2).
	// SIGINT is forwarded to this process group when set; otherwise SIGINT
	// triggers graceful shutdown of the agent itself.
	activeProcess atomic.Pointer[os.Process]
}

// SetStdout overrides the Runtime's stdout (fd 3 delivery channel).
// Must be called before Run(). Used by tests to capture deliverables.
func (r *Runtime) SetStdout(f *os.File) {
	r.stdout = f
	r.sh.Stdout = f
}

// SetStderr overrides the Runtime's stderr (failure signal channel).
// Must be called before Run(). Used by tests to capture error output.
func (r *Runtime) SetStderr(f *os.File) {
	r.stderr = f
}

// SetStdin overrides the Runtime's stdin (fd 4 material channel).
// Must be called before Run(). Used by tests to provide piped input.
func (r *Runtime) SetStdin(f *os.File) {
	r.sh.Stdin = f
}

// New creates a Runtime from config. Call Run() to start the loop.
func New(cfg *config.Config) (*Runtime, error) {
	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating provider: %w", err)
	}
	return newRuntime(cfg, provider), nil
}

// NewWithProvider creates a Runtime with a custom provider (for testing).
func NewWithProvider(cfg *config.Config, provider llm.Provider) *Runtime {
	return newRuntime(cfg, provider)
}

func newRuntime(cfg *config.Config, provider llm.Provider) *Runtime {
	shortID := cfg.SessionID
	if len(shortID) > 4 {
		shortID = shortID[:4]
	}

	// Compute child environment for recursive invocations.
	// ChildEnv() returns QUINE_* vars with DEPTH+1, fresh SESSION_ID, etc.
	// If it fails (e.g., crypto/rand error), fall back to no child overrides.
	childEnv, err := cfg.ChildEnv()
	if err != nil {
		childEnv = nil
	}

	// Derive lock directory from data dir: {dataDir}/locks/
	lockDir := filepath.Join(cfg.DataDir, "locks")

	// Create dedicated log file for operational messages (§10.2).
	// Location: ${QUINE_DATA_DIR}/${SESSION_ID}.log (flat structure)
	os.MkdirAll(cfg.DataDir, 0o755)
	logPath := filepath.Join(cfg.DataDir, cfg.SessionID+".log")
	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)

	r := &Runtime{
		cfg:       cfg,
		provider:  provider,
		sh:        tools.NewShExecutor(cfg, childEnv),
		fork:      tools.NewForkExecutor(cfg, childEnv),
		tools:     tools.AllToolSchemas(),
		semaphore: NewSemaphore(lockDir, cfg.MaxConcurrent, cfg.SessionID),
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		logFile:   logFile,
	}

	// Wire the process's real stdin/stdout to the sh executor so that
	// commands can read from /dev/stdin and write to /dev/stdout.
	r.sh.Stdin = os.Stdin
	r.sh.Stdout = r.stdout

	// Wire process tracking callbacks so SIGINT can be forwarded to
	// the active tool subprocess (§2.2).
	r.sh.ProcessStarted = func(proc *os.Process) {
		r.activeProcess.Store(proc)
	}
	r.sh.ProcessEnded = func() {
		r.activeProcess.Store(nil)
	}

	// Wire fork executor process tracking for SIGINT forwarding.
	r.fork.ProcessStarted = func(proc *os.Process) {
		r.activeProcess.Store(proc)
	}
	r.fork.ProcessEnded = func() {
		r.activeProcess.Store(nil)
	}

	// Operational logs → log file (silent if file creation failed).
	r.log = func(format string, args ...any) {
		if r.logFile != nil {
			msg := fmt.Sprintf(format, args...)
			fmt.Fprintf(r.logFile, "quine[%s]: %s\n", shortID, msg)
		}
	}

	// Failure signals → OS stderr (parent's gradient).
	r.logError = func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(r.stderr, "quine[%s]: %s\n", shortID, msg)
	}

	// Route semaphore operational logs to the log file.
	if logFile != nil {
		r.semaphore.logWriter = logFile
	}

	// Redirect LLM retry logs to the log file.
	llm.SetLogOutput(logFile)

	return r
}

// setupSignalHandler installs handlers for SIGINT, SIGTERM, and SIGALRM (§2.2).
//
// Signal behavior:
//   - SIGALRM: Sets panicMode flag. The turn loop will inject a "System 1
//     Override" message forcing the agent to exit with its best current answer.
//   - SIGINT: If a tool subprocess is running, forwards SIGINT to its process
//     group (letting e.g. python handle Ctrl+C). If no tool is running, triggers
//     graceful shutdown (same as SIGTERM).
//   - SIGTERM: Flushes the Tape to disk and exits with code 143.
//   - SIGPIPE: Downstream pipe closed. Flushes the Tape and exits with code 141.
//   - SIGHUP: Terminal hangup. Flushes the Tape and exits with code 129.
func (r *Runtime) setupSignalHandler() {
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGALRM, syscall.SIGPIPE, syscall.SIGHUP)

	go func() {
		for sig := range sigCh {
			switch sig {
			case syscall.SIGALRM:
				// Time pressure: set panic mode flag.
				// The turn loop checks this and injects the override message.
				r.panicMode.Store(true)
				r.log("SIGALRM received, entering panic mode")

			case os.Interrupt: // SIGINT
				// If a tool subprocess is running, forward SIGINT to it.
				if proc := r.activeProcess.Load(); proc != nil {
					r.log("SIGINT received, forwarding to active process (pid=%d)", proc.Pid)
					// Send to the process group so child trees also get it.
					_ = syscall.Kill(-proc.Pid, syscall.SIGINT)
					continue
				}
				// No tool running — treat as graceful shutdown.
				r.log("SIGINT received, no active tool, shutting down")
				r.gracefulShutdown(130) // 128 + 2

			case syscall.SIGHUP:
				r.log("SIGHUP received, terminal hangup")
				r.gracefulShutdown(129) // 128 + 1

			case syscall.SIGPIPE:
				r.log("SIGPIPE received, downstream pipe closed")
				r.gracefulShutdown(141) // 128 + 13

			case syscall.SIGTERM:
				r.log("SIGTERM received, shutting down")
				r.gracefulShutdown(143) // 128 + 15
			}
		}
	}()
}

// gracefulShutdown flushes the tape, closes the log file, and exits.
func (r *Runtime) gracefulShutdown(exitCode int) {
	// Kill active child process to prevent orphans.
	// Use negative PID to kill the entire process group (children use Setpgid: true).
	if proc := r.activeProcess.Load(); proc != nil {
		_ = syscall.Kill(-proc.Pid, syscall.SIGKILL)
	}

	// Set session outcome if tape is initialized
	if r.tape != nil {
		duration := time.Since(r.startTime)
		r.tape.SetOutcome(tape.SessionOutcome{
			ExitCode:        exitCode,
			Stderr:          fmt.Sprintf("terminated by signal (exit %d)", exitCode),
			DurationMs:      duration.Milliseconds(),
			TerminationMode: tape.TermSignal,
		})

		// Flush outcome to disk
		if r.tapeWriter != nil {
			r.writeTapeEntry(r.tape.OutcomeEntry())
			r.tapeWriter.Close()
		}
	}

	// Close persistent shell before exit.
	if r.sh != nil {
		r.sh.Close()
	}

	// Close log file before exit (deferred close won't run after os.Exit).
	if r.logFile != nil {
		r.logFile.Close()
	}

	os.Exit(exitCode)
}

// Run executes the full agent lifecycle:
//  1. Initialize Tape with system + user messages
//  2. Enter Turn Loop
//  3. Return the exit code
//
// Parameters:
//   - mission: The task/goal from argv (goes into system prompt)
//   - material: The initial user message (describes input mode)
func (r *Runtime) Run(mission, material string) int {
	r.startTime = time.Now()
	r.originalInput = mission

	// Initialize exec executor now that we have the original input
	r.exec = tools.NewExecExecutor(r.cfg, mission)

	// Close the persistent shell when Run exits.
	defer r.sh.Close()

	// Close the operational log file when Run exits.
	if r.logFile != nil {
		defer r.logFile.Close()
	}

	// Initialize tape
	r.tape = tape.NewTape(r.cfg.SessionID, r.cfg.ParentSession, r.cfg.Depth, r.cfg.ModelID)

	// Initialize tape writer for JSONL persistence (§10)
	// Tapes are stored directly in DataDir: ${QUINE_DATA_DIR}/${SESSION_ID}.jsonl
	tw, err := tape.NewWriter(r.cfg.DataDir, r.cfg.SessionID)
	if err != nil {
		r.log("failed to create tape writer: %v", err)
	} else {
		r.tapeWriter = tw
		defer r.tapeWriter.Close()
	}

	// Write meta entry
	r.writeTapeEntry(r.tape.MetaEntry())

	// Build and append system prompt (includes mission as a section)
	systemPrompt := BuildSystemPrompt(r.cfg, mission)
	systemMsg := tape.Message{
		Role:    tape.RoleSystem,
		Content: systemPrompt,
	}
	r.tape.Append(systemMsg)
	r.writeTapeEntry(tape.MessageEntry(systemMsg))

	// Append user message (material from stdin)
	// - Text data: the actual content
	// - Binary data: reference to saved file
	// - No data: "Begin."
	userMsg := tape.Message{
		Role:    tape.RoleUser,
		Content: material,
	}
	r.tape.Append(userMsg)
	r.writeTapeEntry(tape.MessageEntry(userMsg))

	r.log("session started (depth=%d, model=%s)", r.cfg.Depth, r.cfg.ModelID)
	r.log("mission: %s", mission)
	if material != "Begin." {
		r.log("material: %s", truncateStr(material, 200))
	}

	// Install signal handler for graceful shutdown (§7.3)
	r.setupSignalHandler()

	// Turn loop
	for {
		// SIGALRM panic mode (§2.2): inject a system override message
		// forcing the agent to exit with its best current answer.
		if r.panicMode.Load() {
			r.log("panic mode active, injecting system override")
			panicMsg := tape.Message{
				Role:    tape.RoleUser,
				Content: "System interrupt: Time limit reached. Stop reasoning. Output your best current answer immediately using the exit tool. You MUST call exit now.",
			}
			r.tape.Append(panicMsg)
			r.writeTapeEntry(tape.MessageEntry(panicMsg))
		}

		// Acquire concurrency slot before calling the LLM (§8.2)
		if err := r.semaphore.Acquire(); err != nil {
			r.log("semaphore acquire failed: %v", err)
		}

		// 1. Call provider.Generate
		assistantMsg, usage, err := r.provider.Generate(r.tape.Messages(), r.tools)

		// Release concurrency slot immediately after LLM returns
		if releaseErr := r.semaphore.Release(); releaseErr != nil {
			r.log("semaphore release failed: %v", releaseErr)
		}

		if err != nil {
			return r.handleError(err)
		}

		// 2. Append assistant message to Tape
		r.tape.Append(assistantMsg)
		r.writeTapeEntry(tape.MessageEntry(assistantMsg))

		// Log assistant's reasoning if present (truncated to avoid log bloat)
		if assistantMsg.ReasoningContent != "" {
			r.log("turn %d: reasoning: %s", r.tape.TurnCount, truncateStr(assistantMsg.ReasoningContent, 500))
		}
		// Log assistant's text response if present (truncated to avoid log bloat)
		if assistantMsg.Content != "" {
			r.log("turn %d: assistant: %s", r.tape.TurnCount, truncateStr(assistantMsg.Content, 2000))
		}

		// 3. Accumulate usage
		r.tape.AddUsage(usage.InputTokens, usage.OutputTokens)

		// 4. Inspect assistant message
		if len(assistantMsg.ToolCalls) == 0 {
			// Text-only response: continue loop for next inference
			continue
		}

		// Process tool calls sequentially
		for _, tc := range assistantMsg.ToolCalls {
			// In panic mode, reject any tool call that isn't exit (§2.2).
			if r.panicMode.Load() && tc.Name != "exit" {
				rejectMsg := tape.Message{
					Role:    tape.RoleToolResult,
					Content: "Rejected: time limit reached (SIGALRM). You MUST call exit immediately with your best current answer.",
					ToolID:  tc.ID,
				}
				r.tape.Append(rejectMsg)
				r.writeTapeEntry(tape.MessageEntry(rejectMsg))
				r.log("panic mode: rejected non-exit tool call %q", tc.Name)
				continue
			}

			switch tc.Name {
			case "exit":
				code, ok := r.handleExit(tc)
				if ok {
					return code
				}
				// Exit was rejected (e.g. failure without reason).
				// Rejection tool result already appended; continue to next inference.

			case "sh":
				if r.handleSh(tc) {
					// Turn limit exhausted — give the agent a near-death experience.
					// Inject a warning into the last tool result, then allow ONE final
					// inference. If the agent calls exec, it survives. Otherwise it dies.
					r.log("turn limit reached (%d/%d) — near-death warning issued", r.tape.TurnCount, r.cfg.MaxTurns)
					if last := r.tape.LastMessage(); last != nil && last.Role == tape.RoleToolResult {
						last.Content += "\n[RESOURCE EXHAUSTION IMMINENT] Process will be terminated after this response. To survive, call exec now with wisdom to preserve your state. This is your last chance."
					}

					// One final inference
					if err := r.semaphore.Acquire(); err != nil {
						r.log("semaphore acquire failed (near-death): %v", err)
					}
					finalMsg, finalUsage, err := r.provider.Generate(r.tape.Messages(), r.tools)
					if releaseErr := r.semaphore.Release(); releaseErr != nil {
						r.log("semaphore release failed (near-death): %v", releaseErr)
					}
					if err != nil {
						return r.handleError(err)
					}
					r.tape.Append(finalMsg)
					r.writeTapeEntry(tape.MessageEntry(finalMsg))
					r.tape.AddUsage(finalUsage.InputTokens, finalUsage.OutputTokens)
					if finalMsg.Content != "" {
						r.log("near-death response: %s", truncateStr(finalMsg.Content, 2000))
					}

					// Check if the agent called exec in its final breath
					execCalled := false
					for _, lastTC := range finalMsg.ToolCalls {
						if lastTC.Name == "exec" {
							r.log("near-death exec — agent chose survival")
							r.handleExec(lastTC)
							execCalled = true
							break // handleExec does not return (it calls syscall.Exec)
						}
						// Reject any non-exec tool call
						rejectMsg := tape.Message{
							Role:    tape.RoleToolResult,
							Content: "Rejected: resource exhaustion. Only exec is accepted at this point.",
							ToolID:  lastTC.ID,
						}
						r.tape.Append(rejectMsg)
						r.writeTapeEntry(tape.MessageEntry(rejectMsg))
						r.log("near-death: rejected tool call %q (only exec accepted)", lastTC.Name)
					}

					if !execCalled {
						r.log("turn limit exhausted (%d/%d)", r.tape.TurnCount, r.cfg.MaxTurns)
						r.logError("turn limit exhausted (%d/%d)", r.tape.TurnCount, r.cfg.MaxTurns)
						duration := time.Since(r.startTime)
						r.tape.SetOutcome(tape.SessionOutcome{
							ExitCode:        1,
							Stderr:          fmt.Sprintf("turn limit exhausted (%d/%d)", r.tape.TurnCount, r.cfg.MaxTurns),
							DurationMs:      duration.Milliseconds(),
							TerminationMode: tape.TermTurnExhaustion,
						})
						r.writeTapeEntry(r.tape.OutcomeEntry())
						return 1
					}
				}

			case "fork":
				r.handleFork(tc)

			case "exec":
				r.handleExec(tc)

			default:
				// Unknown tool — return error result
				unknownMsg := tape.Message{
					Role:    tape.RoleToolResult,
					Content: fmt.Sprintf("unknown tool: %s", tc.Name),
					ToolID:  tc.ID,
				}
				r.tape.Append(unknownMsg)
				r.writeTapeEntry(tape.MessageEntry(unknownMsg))
			}
		}

		// Inject resource usage into the last tool result so the agent
		// sees its situation without breaking tool_use/tool_result pairing.
		if last := r.tape.LastMessage(); last != nil && last.Role == tape.RoleToolResult {
			if r.cfg.MaxTurns > 0 {
				remaining := r.cfg.MaxTurns - r.tape.TurnCount
				last.Content += fmt.Sprintf("\n[TURNS LEFT] %d", remaining)
			}
			last.Content += fmt.Sprintf("\n[CONTEXT USED] %dK", usage.InputTokens/1000)
		}
	}
}

// handleExit processes an exit tool call. Returns (exitCode, true) if the
// process should exit, or (0, false) if the exit was rejected (e.g. failure
// without a reason) and a rejection tool result was sent back to the agent.
func (r *Runtime) handleExit(tc tape.ToolCall) (int, bool) {
	exitReq, err := tools.ParseExitArgs(tc.Arguments)
	if err != nil {
		r.log("failed to parse exit args: %v", err)
		exitReq = tools.ExitRequest{Status: tools.StatusFailure, Stderr: fmt.Sprintf("invalid exit args: %v", err)}
	}

	// Validate semantic constraints — bounce back to agent on violation
	if err := exitReq.Validate(); err != nil {
		rejectMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("Exit rejected: %s", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(rejectMsg)
		r.writeTapeEntry(tape.MessageEntry(rejectMsg))
		r.log("exit rejected: %v", err)
		return 0, false
	}

	exitCode := exitReq.ExitCode()
	turnNum := r.tape.TurnCount

	// Log exit
	var argParts []string
	argParts = append(argParts, fmt.Sprintf("status=%s", exitReq.Status))
	if exitReq.Stderr != "" {
		argParts = append(argParts, fmt.Sprintf("stderr=%q", exitReq.Stderr))
	}
	r.log("turn %d: assistant called exit(%s)", turnNum, joinArgs(argParts))

	// Write stderr (stdout is only via sh passthrough)
	if exitReq.Stderr != "" {
		fmt.Fprint(r.stderr, exitReq.Stderr)
	}

	// Set outcome
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        exitCode,
		Stderr:          exitReq.Stderr,
		DurationMs:      duration.Milliseconds(),
		TerminationMode: tape.TermExit,
	})

	totalTokens := r.tape.TokensIn + r.tape.TokensOut
	r.log("session ended (exit=%d, %d turns, %.1fs, %d tokens)",
		exitCode, r.tape.TurnCount, duration.Seconds(), totalTokens)

	// Write outcome to tape file
	r.writeTapeEntry(r.tape.OutcomeEntry())

	return exitCode, true
}

// handleSh processes a sh tool call and appends the result to the tape.
// This is the ONLY tool that consumes turns.
// Returns true if the process should terminate (turn limit reached after this call).
func (r *Runtime) handleSh(tc tape.ToolCall) bool {
	// Increment turn counter BEFORE execution (sh is the only turn-consuming tool)
	r.tape.IncrementTurn()
	turnNum := r.tape.TurnCount

	// Extract command from arguments
	command, _ := tc.Arguments["command"].(string)

	// Log the call
	argSummary := truncateStr(command, 60)
	r.log("turn %d: assistant called %s(\"%s\")", turnNum, "sh", argSummary)

	// Execute
	result := r.sh.Execute(tc.ID, command)

	// Log completion
	r.log("turn %d: sh completed (exit=%d, %d bytes)", turnNum, exitCodeFromResult(result), len(result.Content))

	// Append tool result to tape
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))

	// Check if turn limit is now exhausted
	if r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns {
		return true // Signal to terminate
	}
	return false
}

// handleFork processes a fork tool call and appends the result to the tape.
func (r *Runtime) handleFork(tc tape.ToolCall) {
	turnNum := r.tape.TurnCount

	// Parse fork arguments
	forkReq, err := tools.ParseForkArgs(tc.Arguments)
	if err != nil {
		r.log("turn %d: fork parse error: %v", turnNum, err)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[FORK ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// Check depth limit before forking
	if r.cfg.Depth+1 >= r.cfg.MaxDepth {
		r.log("turn %d: fork rejected - depth limit exceeded (%d/%d)", turnNum, r.cfg.Depth+1, r.cfg.MaxDepth)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[FORK ERROR] Max recursion depth exceeded (%d/%d). Cannot spawn child.", r.cfg.Depth+1, r.cfg.MaxDepth),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// Log the call
	waitStr := "false"
	if forkReq.Wait {
		waitStr = "true"
	}
	intentSummary := truncateStr(forkReq.Intent, 60)
	r.log("turn %d: assistant called fork(intent=%q, wait=%s)", turnNum, intentSummary, waitStr)

	// Flush the tape before forking so child gets complete context
	if r.tapeWriter != nil {
		r.tapeWriter.Close()
		// Reopen for continued writing
		tw, err := tape.NewWriter(r.cfg.DataDir, r.cfg.SessionID)
		if err != nil {
			r.log("failed to reopen tape writer after fork: %v", err)
		} else {
			r.tapeWriter = tw
		}
	}

	// Execute fork
	result := r.fork.Execute(tc.ID, forkReq)

	// Log completion
	if result.IsError {
		r.log("turn %d: fork failed: %s", turnNum, truncateStr(result.Content, 100))
	} else {
		r.log("turn %d: fork completed (wait=%s)", turnNum, waitStr)
	}

	// Append tool result to tape
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

// handleExec processes an exec tool call.
// Note: On success, this function does NOT return — the process is replaced.
// On failure, it appends an error result to the tape.
func (r *Runtime) handleExec(tc tape.ToolCall) {
	turnNum := r.tape.TurnCount

	// Parse exec arguments
	execReq, err := tools.ParseExecArgs(tc.Arguments)
	if err != nil {
		r.log("turn %d: exec parse error: %v", turnNum, err)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[EXEC ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// Log the call
	reasonStr := ""
	if execReq.Reason != "" {
		reasonStr = fmt.Sprintf(", reason=%q", truncateStr(execReq.Reason, 40))
	}
	personaStr := ""
	if execReq.Persona != "" {
		personaStr = fmt.Sprintf(", persona=%q", execReq.Persona)
	}
	r.log("turn %d: assistant called exec(%s%s)", turnNum, reasonStr, personaStr)

	// Write outcome before exec (we're about to be replaced)
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        0,
		Stderr:          fmt.Sprintf("exec: metamorphosis to fresh context (reason: %s)", execReq.Reason),
		DurationMs:      duration.Milliseconds(),
		TerminationMode: tape.TermExec,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())

	// Close tape writer before exec
	if r.tapeWriter != nil {
		r.tapeWriter.Close()
	}

	// Close log file before exec
	if r.logFile != nil {
		r.logFile.Close()
	}

	// Execute the exec — this does not return on success
	result := r.exec.Execute(tc.ID, execReq)

	// If we get here, exec failed — reopen log file to record the error
	logPath := filepath.Join(r.cfg.DataDir, r.cfg.SessionID+".log")
	r.logFile, _ = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)

	r.log("turn %d: exec failed: %s", turnNum, truncateStr(result.Content, 100))

	// Append error result to tape (need to reopen writer)
	tw, err := tape.NewWriter(r.cfg.DataDir, r.cfg.SessionID)
	if err == nil {
		r.tapeWriter = tw
	}

	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

// handleError handles LLM errors and returns the appropriate exit code.
// Failure signals are written to stderr (not the log file) so parent
// processes can see why the child died (§10.2).
func (r *Runtime) handleError(err error) int {
	duration := time.Since(r.startTime)

	if errors.Is(err, llm.ErrAuth) {
		r.logError("authentication failed: %v", err)
		r.tape.SetOutcome(tape.SessionOutcome{
			ExitCode:        1,
			Stderr:          err.Error(),
			DurationMs:      duration.Milliseconds(),
			TerminationMode: tape.TermExit,
		})
		r.writeTapeEntry(r.tape.OutcomeEntry())
		return 1
	}

	if errors.Is(err, llm.ErrContextOverflow) {
		r.logError("context exhausted: %v", err)
		r.tape.SetOutcome(tape.SessionOutcome{
			ExitCode:        1,
			Stderr:          fmt.Sprintf("context exhausted: %v", err),
			DurationMs:      duration.Milliseconds(),
			TerminationMode: tape.TermContextExhaustion,
		})
		r.writeTapeEntry(r.tape.OutcomeEntry())
		return 1
	}

	r.logError("LLM error: %v", err)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        1,
		Stderr:          err.Error(),
		DurationMs:      duration.Milliseconds(),
		TerminationMode: tape.TermExit,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())
	return 1
}

// writeTapeEntry writes an entry to the tape writer if available.
// Errors are logged but do not halt execution.
func (r *Runtime) writeTapeEntry(entry tape.TapeEntry) {
	if r.tapeWriter == nil {
		return
	}
	if err := r.tapeWriter.WriteEntry(entry); err != nil {
		r.log("tape write error: %v", err)
	}
}

// truncateStr truncates s to maxLen characters, appending "..." if truncated.
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// joinArgs joins argument parts with ", ".
func joinArgs(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}

// exitCodeFromResult extracts the exit code from a tool result content string.
// The content format is "[EXIT CODE] %d\n...".
func exitCodeFromResult(r tape.ToolResult) int {
	var code int
	fmt.Sscanf(r.Content, "[EXIT CODE] %d", &code)
	return code
}
