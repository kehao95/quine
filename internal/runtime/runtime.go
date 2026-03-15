package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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
	cfg           *config.Config
	provider      llm.Provider
	sh            *tools.ShExecutor
	fork          *tools.ForkExecutor
	exec          *tools.ExecExecutor
	memory        *tools.AnchorMemoryExecutor
	tape          *tape.Tape
	tapeWriter    *tape.Writer
	tools         []llm.ToolSchema
	semaphore     *Semaphore
	agentRegistry *AgentRegistry
	startTime     time.Time
	log           func(format string, args ...any) // operational log → log file
	logError      func(format string, args ...any) // failure signal → stderr

	// originalInput stores the user's input for this session.
	// Needed for exec to preserve the mission.
	originalInput string
	hasMaterial   bool

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

	// Goal stall tracking for escalation hints
	lastGoal   string
	stallCount int

	childEnvBase    []string
	lastFSMutations string
}

// SetStdout overrides the Runtime's stdout (fd 4 delivery channel).
// Must be called before Run(). Used by tests to capture deliverables.
func (r *Runtime) SetStdout(f *os.File) {
	r.stdout = f
	r.sh.Stdout = f
}

// SetStderr overrides the Runtime's stderr (failure signal channel).
// Must be called before Run(). Used by tests to capture error output.
func (r *Runtime) SetStderr(f *os.File) {
	r.stderr = f
	r.sh.Stderr = f
}

// SetStdin overrides the Runtime's stdin (fd 3 material channel).
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
	rt := newRuntime(cfg, provider)
	if err := rt.sh.Prepare(); err != nil {
		return nil, err
	}
	return rt, nil
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

	lockDir := cfg.LockDir()

	// Create dedicated log file for operational messages (§10.2).
	// Location: ${QUINE_DATA_DIR}/${SESSION_ID}.log (flat structure under the runtime root)
	os.MkdirAll(cfg.RuntimeRoot(), 0o755)
	logPath := cfg.SessionLogPath("")
	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)

	r := &Runtime{
		cfg:           cfg,
		provider:      provider,
		sh:            tools.NewShExecutor(cfg, childEnv),
		fork:          tools.NewForkExecutor(cfg, childEnv),
		tools:         tools.AllToolSchemas(cfg),
		semaphore:     NewSemaphore(lockDir, cfg.MaxConcurrent, cfg.SessionID),
		agentRegistry: NewAgentRegistry(lockDir, cfg.MaxAgents, cfg.SessionID),
		stdout:        os.Stdout,
		stderr:        os.Stderr,
		logFile:       logFile,
		childEnvBase:  childEnv,
	}
	if cfg.AnchorMemoryEnabled {
		r.memory = tools.NewAnchorMemoryExecutor(cfg.AgentRoot(), cfg.TapePath(""))
	}
	r.sh.Env = tools.MergeEnv(r.sh.Env, []string{
		"QUINE_AGENT_ROOT=" + cfg.AgentRoot(),
	})

	// Wire the process's real stdin/stdout to the sh executor so that
	// commands can use runtime side channels:
	// fd 3 = material stdin, fd 4 = deliverable stdout, fd 5 = failure stderr.
	r.sh.Stdin = os.Stdin
	r.sh.Stdout = r.stdout
	r.sh.Stderr = r.stderr

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

	// Deregister from agent registry
	if r.agentRegistry != nil {
		r.agentRegistry.Deregister()
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
		// If exiting successfully (0), spare any detached jobs.
		// If failing (!0), kill everything.
		keepDetached := (exitCode == 0)
		r.sh.Close(keepDetached)
	}

	// Close log file before exit (deferred close won't run after os.Exit).
	if r.logFile != nil {
		r.logFile.Close()
	}

	os.Exit(exitCode)
}

// checkGoalStall returns true if the model has been working on the same goal for StallThreshold+ turns
func (r *Runtime) checkGoalStall(goal string) bool {
	if goal == "" {
		return false // Missing goal, skip tracking
	}

	if goal == r.lastGoal {
		r.stallCount++
	} else {
		r.lastGoal = goal
		r.stallCount = 1
	}

	return r.stallCount >= r.cfg.StallThreshold
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
	r.syncWorldRevisionSurface()

	// Register this agent in the global registry
	if err := r.agentRegistry.Register(); err != nil {
		r.logError("agent registration failed: %v", err)
		return 1
	}
	defer r.agentRegistry.Deregister()

	// Initialize exec executor now that we have the original input
	r.exec = tools.NewExecExecutor(r.cfg, mission)

	// Close the persistent shell when Run exits.
	// On success (exit code 0), spare detached jobs so daemons survive.
	// On failure (non-zero), kill everything to be safe.
	// gracefulShutdown handles the signal-based exit path separately.
	runExitCode := 1 // default: kill everything (all error paths return 1)
	defer func() {
		r.sh.Close(runExitCode == 0)
	}()

	// Close the operational log file when Run exits.
	if r.logFile != nil {
		defer r.logFile.Close()
	}

	// Initialize tape
	r.tape = tape.NewTape(r.cfg.SessionID, r.cfg.ParentSession, r.cfg.Depth, r.cfg.ModelID)

	// Initialize tape writer for JSONL persistence (§10).
	// A stable session can accumulate multiple tapes across exec incarnations.
	tw, err := tape.NewWriter(r.cfg.TapeDir(""), r.cfg.TapeID)
	if err != nil {
		r.log("failed to create tape writer: %v", err)
	} else {
		r.tapeWriter = tw
		defer r.tapeWriter.Close()
	}

	// Write meta entry
	r.writeTapeEntry(r.tape.MetaEntry())

	// Build and append system prompt (includes mission as a section)
	r.hasMaterial = material != "Begin."
	systemPrompt := BuildSystemPrompt(r.cfg, mission, r.hasMaterial)
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

	// Install signal handler for graceful shutdown (§7.3)
	r.setupSignalHandler()

	// Turn loop
	consecutiveTextOnly := 0
	const maxConsecutiveTextOnly = 3
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

		// 3. Accumulate usage
		r.tape.AddUsage(usage.InputTokens, usage.OutputTokens)

		// 4. Inspect assistant message
		if len(assistantMsg.ToolCalls) == 0 {
			// Text-only response: continue loop for next inference.
			// Guard against infinite loops where the model never calls a tool.
			consecutiveTextOnly++
			if consecutiveTextOnly >= maxConsecutiveTextOnly {
				r.log("aborting: %d consecutive text-only responses without tool calls", consecutiveTextOnly)
				r.logError("agent stuck: %d consecutive text-only responses, forcing exit", consecutiveTextOnly)
				return 1
			}
			continue
		}
		consecutiveTextOnly = 0 // reset on any tool call

		// Process tool calls sequentially
		for tcIdx, tc := range assistantMsg.ToolCalls {
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
					runExitCode = code
					return code
				}
				// Exit was rejected (e.g. failure without reason).
				// Rejection tool result already appended; continue to next inference.

			case "sh":
				if r.handleSh(tc) {
					return r.handleExecutionBudgetExhaustion(assistantMsg.ToolCalls[tcIdx+1:])
				}

			case "fork":
				r.handleFork(tc)

			case "exec":
				r.handleExec(tc)

			case "restore_world":
				r.handleRestoreWorld(tc)

			case "vision":
				r.handleVision(tc)

			case "mark":
				r.handleMark(tc)

			case "unfold":
				r.handleUnfold(tc)

			case "escalate":
				r.handleEscalate(tc)

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
				last.Content += fmt.Sprintf("\n[EXECUTIONS LEFT] %d", remaining)
			}
			last.Content += fmt.Sprintf("\n[CONTEXT USED] %dK", usage.InputTokens/1000)
		}
	}
}

// handleExecutionBudgetExhaustion enforces the configured policy when the
// execution budget reaches zero after a sh call. It always terminates unless
// an exec tool call replaces the process image.
func (r *Runtime) handleExecutionBudgetExhaustion(remaining []tape.ToolCall) int {
	// Reject remaining tool calls in this batch to maintain protocol pairing.
	for _, remainingTC := range remaining {
		rejectMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: "Rejected: execution budget exhausted before this tool call could be processed.",
			ToolID:  remainingTC.ID,
		}
		r.tape.Append(rejectMsg)
		r.writeTapeEntry(tape.MessageEntry(rejectMsg))
	}

	if !r.cfg.UsesNearDeathContinuation() {
		return r.terminateExecutionBudgetExhausted()
	}

	r.log("execution budget reached (%d/%d) — continuation window opened", r.tape.TurnCount, r.cfg.MaxTurns)
	if last := r.tape.LastMessage(); last != nil && last.Role == tape.RoleToolResult {
		last.Content += "\n[EXECUTION BUDGET EXHAUSTED] One continuation action is available: call exec with wisdom. Only exec is accepted in the next response."
	}

	// One final inference for exec-only continuation.
	if err := r.semaphore.Acquire(); err != nil {
		r.log("semaphore acquire failed (exec-only continuation): %v", err)
	}
	finalMsg, finalUsage, err := r.provider.Generate(r.tape.Messages(), r.tools)
	if releaseErr := r.semaphore.Release(); releaseErr != nil {
		r.log("semaphore release failed (exec-only continuation): %v", releaseErr)
	}
	if err != nil {
		return r.handleError(err)
	}
	r.tape.Append(finalMsg)
	r.writeTapeEntry(tape.MessageEntry(finalMsg))
	r.tape.AddUsage(finalUsage.InputTokens, finalUsage.OutputTokens)

	for _, tc := range finalMsg.ToolCalls {
		if tc.Name == "exec" {
			r.log("execution-exhausted continuation: exec selected")
			r.handleExec(tc)
			return 1 // unreachable on successful exec
		}

		rejectMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: "Rejected: execution budget exhausted. Only exec continuation is accepted.",
			ToolID:  tc.ID,
		}
		r.tape.Append(rejectMsg)
		r.writeTapeEntry(tape.MessageEntry(rejectMsg))
	}

	return r.terminateExecutionBudgetExhausted()
}

func (r *Runtime) terminateExecutionBudgetExhausted() int {
	r.log("execution budget exhausted (%d/%d)", r.tape.TurnCount, r.cfg.MaxTurns)
	r.logError("execution budget exhausted (%d/%d)", r.tape.TurnCount, r.cfg.MaxTurns)
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        1,
		Stderr:          fmt.Sprintf("execution budget exhausted (%d/%d)", r.tape.TurnCount, r.cfg.MaxTurns),
		DurationMs:      duration.Milliseconds(),
		TerminationMode: tape.TermTurnExhaustion,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())
	return 1
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

	// Log exit
	var argParts []string
	argParts = append(argParts, fmt.Sprintf("status=%s", exitReq.Status))
	if exitReq.Stderr != "" {
		argParts = append(argParts, fmt.Sprintf("stderr=%s", truncateOneLine(exitReq.Stderr, 80)))
	}
	r.log("exit(%s)", joinArgs(argParts))

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
// This is the ONLY tool that consumes executions.
// Returns true if the execution budget is exhausted after this call.
func (r *Runtime) handleSh(tc tape.ToolCall) bool {
	// Increment execution counter BEFORE execution (sh is the only execution-consuming tool).
	r.tape.IncrementTurn()
	// Extract command and budget parameters
	command, _ := tc.Arguments["command"].(string)
	detach, _ := tc.Arguments["detach"].(bool)
	interactive, _ := tc.Arguments["interactive"].(bool)
	stdin, _ := tc.Arguments["stdin"].(string)
	goal, _ := tc.Arguments["goal"].(string)
	// strategy is required in schema but not used for stall detection (only goal matters)
	_, _ = tc.Arguments["strategy"].(string)

	// Execute
	r.sh.TurnID = r.tape.TurnCount
	result := r.sh.Execute(tc.ID, command, 0, 0, interactive, detach, stdin)

	// Check for goal stall and inject escalation hint if applicable
	if goal != "" && r.cfg.CanEscalate() && r.checkGoalStall(goal) {
		result.Content += fmt.Sprintf("\n\n⚠️ STALL: %d turns on same goal without progress. STOP. Your variations are not working. Call `escalate` now — it costs nothing and a smarter model may see what you're missing.", r.stallCount)
	}
	if mutationBlock := extractFSMutationsBlock(result.Content); mutationBlock != "" {
		r.lastFSMutations = mutationBlock
	}
	r.syncWorldRevisionSurface()

	// Append tool result to tape
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))

	// Check if execution budget is now exhausted.
	if r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns {
		return true
	}
	return false
}

// handleFork processes a fork tool call and appends the result to the tape.
func (r *Runtime) handleFork(tc tape.ToolCall) {
	// Parse fork arguments
	forkReq, err := tools.ParseForkArgs(tc.Arguments)
	if err != nil {
		r.log("fork parse error: %v", err)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[FORK ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// Check depth limit before forking (disabled when MaxDepth <= 0).
	if r.cfg.MaxDepth > 0 && r.cfg.Depth+1 >= r.cfg.MaxDepth {
		r.log("fork rejected - depth limit exceeded (%d/%d)", r.cfg.Depth+1, r.cfg.MaxDepth)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[FORK ERROR] Max recursion depth exceeded (%d/%d). Cannot spawn child.", r.cfg.Depth+1, r.cfg.MaxDepth),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// All-or-Nothing slot allocation: if the agent requests N children but
	// fewer than N slots are available, the entire fork operation fails.
	// No partial swarms — the agent must make a conscious decision.
	requested := len(forkReq.Children)
	if r.cfg.MaxAgents > 0 {
		current := r.agentRegistry.Count()
		available := r.cfg.MaxAgents - current
		if available < requested {
			r.log("fork rejected - insufficient slots (requested %d, available %d, current %d/%d)",
				requested, available, current, r.cfg.MaxAgents)
			errMsg := tape.Message{
				Role:    tape.RoleToolResult,
				Content: fmt.Sprintf("[FORK ERROR] Insufficient slots. Requested %d, available %d. Release agents or request fewer.", requested, available),
				ToolID:  tc.ID,
			}
			r.tape.Append(errMsg)
			r.writeTapeEntry(tape.MessageEntry(errMsg))
			return
		}
	}

	// Flush the tape before forking so child gets complete context
	if r.tapeWriter != nil {
		r.tapeWriter.Close()
		// Reopen for continued writing
		tw, err := tape.NewWriter(r.cfg.TapeDir(""), r.cfg.TapeID)
		if err != nil {
			r.log("failed to reopen tape writer after fork: %v", err)
		} else {
			r.tapeWriter = tw
		}
	}

	// Execute fork
	r.refreshForkEnv()
	result := r.fork.Execute(tc.ID, forkReq)

	// Log completion
	if result.IsError {
		r.log("fork done [error]: %s", truncateStr(result.Content, 100))
	} else {
		r.log("fork done (n=%d, mode=%s)", len(forkReq.Children), forkReq.Mode)
	}

	// Append tool result to tape
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

func (r *Runtime) refreshForkEnv() {
	if r.fork == nil {
		return
	}
	env := append([]string(nil), r.childEnvBase...)
	if r.sh != nil {
		env = append(env, r.sh.ChildEnvOverrides()...)
	}
	r.fork.Env = tools.MergeEnv(withoutProcessIdentity(os.Environ()), env)
}

func withoutProcessIdentity(env []string) []string {
	out := make([]string, 0, len(env))
	for _, entry := range env {
		if strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			continue
		}
		if strings.HasPrefix(entry, "QUINE_TAPE_ID=") {
			continue
		}
		out = append(out, entry)
	}
	return out
}

// handleExec processes an exec tool call.
// Note: On success, this function does NOT return — the process is replaced.
// On failure, it appends an error result to the tape.
func (r *Runtime) handleExec(tc tape.ToolCall) {
	// Parse exec arguments
	execReq, err := tools.ParseExecArgs(tc.Arguments)
	if err != nil {
		r.log("exec parse error: %v", err)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[EXEC ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// Log the call — show persona and wisdom keys (truncate values)
	{
		var parts []string
		if execReq.Persona != "" {
			parts = append(parts, fmt.Sprintf("persona=%s", execReq.Persona))
		}
		for k, v := range execReq.Wisdom {
			parts = append(parts, fmt.Sprintf("wisdom.%s=%s", k, truncateOneLine(v, 40)))
		}
		r.log("exec(%s)", joinArgs(parts))
	}

	// Write outcome before exec (we're about to be replaced)
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        0,
		Stderr:          "exec: metamorphosis to fresh context",
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
	logPath := r.cfg.SessionLogPath("")
	r.logFile, _ = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)

	r.log("exec failed: %s", truncateStr(result.Content, 100))

	// Append error result to tape (need to reopen writer)
	tw, err := tape.NewWriter(r.cfg.TapeDir(""), r.cfg.TapeID)
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

func (r *Runtime) handleRestoreWorld(tc tape.ToolCall) {
	req, err := tools.ParseRestoreWorldArgs(tc.Arguments)
	if err != nil {
		r.log("restore_world parse error: %v", err)
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[RESTORE WORLD ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	result := r.sh.RestoreWorld(tc.ID, req.Revision)
	if !result.IsError {
		r.lastFSMutations = "[FS MUTATIONS]\n(restored from world revision)"
	}
	r.syncWorldRevisionSurface()
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

func (r *Runtime) syncWorldRevisionSurface() {
	if r.cfg == nil || !r.cfg.WorkspaceEnabled || r.sh == nil {
		return
	}
	revision := r.sh.CurrentWorldRevision()
	if strings.TrimSpace(revision) == "" {
		return
	}
	r.cfg.WorkspaceCurrentRevision = revision
	if childEnv, err := r.cfg.ChildEnv(); err == nil {
		r.childEnvBase = childEnv
		if r.fork != nil {
			started := r.fork.ProcessStarted
			ended := r.fork.ProcessEnded
			r.fork = tools.NewForkExecutor(r.cfg, childEnv)
			r.fork.ProcessStarted = started
			r.fork.ProcessEnded = ended
		}
	}
	if r.tape != nil {
		r.tape.SetSystemPrompt(BuildSystemPrompt(r.cfg, r.originalInput, r.hasMaterial))
	}
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
	if err := r.syncAgentRoot(); err != nil {
		r.log("agent-root sync error: %v", err)
	}
}

func (r *Runtime) syncAgentRoot() error {
	agentRoot := r.cfg.AgentRoot()
	contextDir := filepath.Join(agentRoot, "context")
	contextFrontierDir := filepath.Join(contextDir, "frontier")
	contextAnchorsDir := filepath.Join(contextDir, "anchors")
	logDir := filepath.Join(agentRoot, "log")
	statusDir := filepath.Join(agentRoot, "status")
	worldDir := filepath.Join(agentRoot, "world")
	missionPath := filepath.Join(agentRoot, "mission.txt")

	for _, dir := range []string{contextDir, contextFrontierDir, contextAnchorsDir, logDir, statusDir, worldDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}

	if err := writeTextFile(missionPath, r.originalInput); err != nil {
		return err
	}

	status := map[string]any{
		"session_id":            r.cfg.SessionID,
		"tape_id":               r.cfg.TapeID,
		"parent_session":        r.cfg.ParentSession,
		"depth":                 r.cfg.Depth,
		"model_id":              r.cfg.ModelID,
		"runtime_root":          r.cfg.RuntimeRoot(),
		"agent_root":            agentRoot,
		"workspace_enabled":     r.cfg.WorkspaceEnabled,
		"workspace_root":        r.cfg.WorkspaceRoot,
		"workspace":             r.cfg.Workspace,
		"workspace_session":     r.cfg.WorkspaceSession,
		"anchor_memory_enabled": r.cfg.AnchorMemoryEnabled,
	}
	if err := writeJSONFile(filepath.Join(statusDir, "session.json"), status); err != nil {
		return err
	}

	if err := replaceSymlink(filepath.Join(logDir, "current.jsonl"), r.cfg.TapePath("")); err != nil {
		return err
	}
	if err := replaceSymlink(filepath.Join(logDir, "current.log.yaml"), r.cfg.TapeYAMLPath("")); err != nil {
		return err
	}
	if err := replaceSymlink(filepath.Join(logDir, "runtime.log"), r.cfg.SessionLogPath("")); err != nil {
		return err
	}
	if err := replaceSymlink(filepath.Join(logDir, "tapes"), r.cfg.TapeDir("")); err != nil {
		return err
	}
	if err := replaceSymlink(filepath.Join(agentRoot, "jobs"), r.cfg.JobSessionDir("")); err != nil {
		return err
	}

	if err := writeTextFile(filepath.Join(worldDir, "workspace_root"), r.cfg.WorkspaceRoot); err != nil {
		return err
	}
	if err := writeTextFile(filepath.Join(worldDir, "workspace"), r.cfg.Workspace); err != nil {
		return err
	}
	fsMutationsPath := filepath.Join(worldDir, "fs-mutations.latest")
	if strings.TrimSpace(r.lastFSMutations) == "" {
		if err := os.Remove(fsMutationsPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove fs mutations snapshot: %w", err)
		}
	} else if err := writeTextFile(fsMutationsPath, r.lastFSMutations); err != nil {
		return err
	}

	return nil
}

func replaceSymlink(path string, target string) error {
	if !filepath.IsAbs(target) {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			return fmt.Errorf("resolve symlink target %s: %w", target, err)
		}
		target = absTarget
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	if err := os.Symlink(target, path); err != nil {
		return fmt.Errorf("symlink %s -> %s: %w", path, target, err)
	}
	return nil
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	return writeFile(path, data)
}

func writeTextFile(path string, value string) error {
	if value == "" {
		value = "\n"
	} else if !strings.HasSuffix(value, "\n") {
		value += "\n"
	}
	return writeFile(path, []byte(value))
}

func writeFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func (r *Runtime) appendMemoryFeedback(result *tape.ToolResult) {
	if r.memory == nil || result == nil {
		return
	}
	block, err := r.memory.FeedbackBlock()
	if err != nil {
		return
	}
	if strings.TrimSpace(block) == "" {
		return
	}
	result.Content += "\n" + block
}

// handleMark processes a mark tool call and appends the result to the tape.
func (r *Runtime) handleMark(tc tape.ToolCall) {
	if r.memory == nil {
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: "mark is not available in this runtime.",
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}
	req, err := tools.ParseMarkArgs(tc.Arguments)
	if err != nil {
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[MARK ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}
	result := r.memory.Mark(tc.ID, req)
	r.appendMemoryFeedback(&result)
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

// handleUnfold processes an unfold tool call and appends the result to the tape.
func (r *Runtime) handleUnfold(tc tape.ToolCall) {
	if r.memory == nil {
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: "unfold is not available in this runtime.",
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}
	req, err := tools.ParseUnfoldArgs(tc.Arguments)
	if err != nil {
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("[UNFOLD ERROR] %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}
	result := r.memory.Unfold(tc.ID, req)
	r.appendMemoryFeedback(&result)
	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

// handleVision processes a vision tool call and appends the result to the tape.
// Vision does NOT consume a turn — it is a pure observation tool.
func (r *Runtime) handleVision(tc tape.ToolCall) {
	result := tools.HandleVision(tc.ID, tc.Arguments)

	r.tape.Append(tape.Message{
		Role:    tape.RoleToolResult,
		Content: result.Content,
		ToolID:  result.ToolID,
		Image:   result.Image,
	})
	r.writeTapeEntry(tape.ToolResultEntry(result))
}

// handleEscalate processes an escalate tool call, hot-swapping to a smarter model.
// Escalate does NOT consume a turn — it is a model upgrade operation.
func (r *Runtime) handleEscalate(tc tape.ToolCall) {
	reason, _ := tc.Arguments["reason"].(string)

	// Guard: already escalated or not configured
	if r.cfg.Escalated || r.cfg.SmartModelID == "" {
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: "Escalation not available.",
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	oldModel := r.cfg.ModelID
	r.log("escalate: %s -> %s (reason=%s)", oldModel, r.cfg.SmartModelID, truncateOneLine(reason, 80))

	// Build new provider from smart config
	smartCfg := *r.cfg
	smartCfg.ModelID = r.cfg.SmartModelID
	smartCfg.APIKey = r.cfg.SmartAPIKey
	smartCfg.APIBase = r.cfg.SmartAPIBase
	smartCfg.Provider = r.cfg.SmartProvider

	newProvider, err := llm.NewProvider(&smartCfg)
	if err != nil {
		errMsg := tape.Message{
			Role:    tape.RoleToolResult,
			Content: fmt.Sprintf("Escalation failed: %v", err),
			ToolID:  tc.ID,
		}
		r.tape.Append(errMsg)
		r.writeTapeEntry(tape.MessageEntry(errMsg))
		return
	}

	// Hot-swap
	r.provider = newProvider
	r.cfg.ModelID = r.cfg.SmartModelID
	r.cfg.Escalated = true

	// Update tape metadata so JSONL records attribute post-escalation turns correctly
	r.tape.ModelID = r.cfg.SmartModelID

	// Remove escalate tool (CanEscalate() now returns false)
	r.tools = tools.AllToolSchemas(r.cfg)

	// Update system prompt in tape via SetSystemPrompt (not Messages() which returns a copy)
	r.tape.SetSystemPrompt(BuildSystemPrompt(r.cfg, r.originalInput, r.hasMaterial))

	// Tool result — reason is already in tc.Arguments (the handoff note)
	resultMsg := tape.Message{
		Role:    tape.RoleToolResult,
		Content: fmt.Sprintf("Escalated to %s.", r.cfg.SmartModelID),
		ToolID:  tc.ID,
	}
	r.tape.Append(resultMsg)
	r.writeTapeEntry(tape.ToolResultEntry(tape.ToolResult{
		ToolID:  tc.ID,
		Content: resultMsg.Content,
	}))
}

// truncateStr truncates s to maxLen characters, appending "..." if truncated.
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// truncateOneLine replaces newlines with spaces and truncates to maxLen.
// Used for log fields that could be multi-line or very long (commands, stdin, etc.).
func truncateOneLine(s string, maxLen int) string {
	// Replace newlines/tabs with spaces to keep log lines single-line
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' {
			b.WriteByte(' ')
		} else {
			b.WriteRune(r)
		}
	}
	return truncateStr(b.String(), maxLen)
}

func extractFSMutationsBlock(content string) string {
	idx := strings.LastIndex(content, "[FS MUTATIONS]")
	if idx < 0 {
		return ""
	}
	block := strings.TrimSpace(content[idx:])
	if block == "" {
		return ""
	}
	return block + "\n"
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
