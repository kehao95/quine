package runtime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

const RecoverableInferenceExitCode = 75

// Runtime orchestrates the agent's execution loop.
type Runtime struct {
	cfg           *config.Config
	provider      llm.Provider
	sh            *tools.ShExecutor
	fork          *tools.ForkExecutor
	spawn         *tools.SpawnExecutor
	exec          *tools.ExecExecutor
	memory        *tools.AnchorMemoryExecutor
	tape          *tape.Tape
	tapeWriter    *tape.Writer
	tools         []llm.ToolSchema
	toolRegistry  map[string]toolSpec
	pendingTool   *pendingToolResult
	semaphore     *Semaphore
	agentRegistry *AgentRegistry
	startTime     time.Time
	log           func(format string, args ...any) // operational log → log file
	logError      func(format string, args ...any) // failure signal → stderr

	// originalInput stores the user's input for this session.
	// Used as default mission argv for self re-exec.
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

	// contextIntegrityErr is set when an integrity-critical entry (an assistant
	// tool_use turn or a tool_result) fails to persist to the provider-input
	// context surface (current.jsonl). A lost entry breaks tool_use/tool_result
	// pairing on the next request, so the turn loop checks this before the next
	// provider call and fails loudly instead of sending a corrupt context.
	contextIntegrityErr error

	// inferenceInFlight is true while a provider request is waiting for an LLM
	// response. Transport SIGPIPE in this window is recoverable from the same
	// durable context boundary.
	inferenceInFlight atomic.Bool

	// activeProcess tracks the currently running tool subprocess (§2.2).
	// SIGINT is forwarded to this process group when set; otherwise SIGINT
	// triggers graceful shutdown of the agent itself.
	activeProcess atomic.Pointer[os.Process]

	peerDiscoveryHeartbeatMu   sync.Mutex
	peerDiscoveryHeartbeatStop chan struct{}
	peerDiscoveryHeartbeatDone chan struct{}
	peerDiscoveryMu            sync.Mutex
	peerDiscoveryKnown         map[int]peerDiscoveryPeer
	peerDiscoveryInitialized   bool

	// Goal stall tracking for escalation hints
	lastGoal   string
	stallCount int

	incarnationID         int
	childEnvBase          []string
	lastFSMutations       string
	agentRootBootstrapped bool
	control               *controlState
	publicSurface         *fusePublicSurfaceBackend
	// Test seams for the public surface. When non-nil they replace the real
	// FUSE mount/unmount, letting finalization-ordering tests run without
	// /dev/fuse. Nil in production.
	publicSurfaceSyncFn    func(publicSurfacePaths) error
	publicSurfaceCleanupFn func() error

	finalizationMu         sync.Mutex
	finalizationPhases     []finalizationPhase
	shFinalized            bool
	closeRuntimeSubstrates func(keepDetached bool, commitWorkspace bool) error
	recoverWorkspaceCommit func() error
}

type pendingToolResult struct {
	result tape.ToolResult
}

func isAssistantNoOp(msg tape.Message) bool {
	return len(msg.ToolCalls) == 0 &&
		strings.TrimSpace(msg.Content) == "" &&
		len(bytes.TrimSpace(msg.StructuredContent)) == 0 &&
		strings.TrimSpace(msg.ReasoningContent) == "" &&
		len(msg.ReasoningItems) == 0 &&
		msg.Image == nil
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
	if strings.TrimSpace(cfg.RunID) == "" {
		cfg.RunID = fmt.Sprintf("run_%s_%d", cfg.SessionID, os.Getpid())
	}
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
	// Location: SessionLogPath(), owned by QUINE_RETENTION_DIR/sessions/<session>
	// when QUINE_RETENTION_DIR is set, otherwise by QUINE_DATA_DIR/log/<session>.
	logPath := cfg.SessionLogPath("")
	os.MkdirAll(filepath.Dir(logPath), 0o755)
	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)

	r := &Runtime{
		cfg:           cfg,
		provider:      provider,
		sh:            tools.NewShExecutor(cfg, childEnv),
		tools:         tools.AllToolSchemas(cfg),
		semaphore:     NewSemaphore(lockDir, cfg.MaxConcurrent, cfg.RunID),
		agentRegistry: NewAgentRegistry(lockDir, cfg.RuntimeRoot(), cfg.MaxAgents, cfg.SessionID, cfg.RunID),
		stdout:        os.Stdout,
		stderr:        os.Stderr,
		logFile:       logFile,
		incarnationID: -1,
		childEnvBase:  childEnv,
		control:       newControlState(),
	}
	r.toolRegistry = newToolRegistry()
	if cfg.ForkEnabled() {
		r.fork = tools.NewForkExecutor(cfg, childEnv)
	}
	if cfg.SpawnEnabled() {
		r.spawn = tools.NewSpawnExecutor(cfg, childEnv)
	}
	if cfg.AnchorMemoryEnabled {
		r.memory = tools.NewAnchorMemoryExecutor(
			cfg.AgentRoot(),
			cfg.MemoryWarnTokens,
			cfg.MemoryDangerTokens,
			cfg.MemoryDeathTokens,
		)
		r.memory.MarkDisabled = !cfg.AnchorMarkEnabled()
		r.memory.FoldDisabled = !cfg.AnchorFoldEnabled()
		r.memory.MemoryStrategyHints = cfg.MemoryStrategyHints
	}
	r.sh.Env = tools.MergeEnv(r.sh.Env, []string{
		"QUINE_AGENT_ROOT=" + cfg.AgentRoot(),
		config.EnvRunID + "=" + cfg.RunID,
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

	// Wire fork executor process tracking for SIGINT forwarding when fork is enabled.
	if r.fork != nil {
		r.fork.ProcessStarted = func(proc *os.Process) {
			r.activeProcess.Store(proc)
		}
		r.fork.ProcessEnded = func() {
			r.activeProcess.Store(nil)
		}
	}
	if r.spawn != nil {
		r.spawn.ProcessStarted = func(proc *os.Process) {
			r.activeProcess.Store(proc)
		}
		r.spawn.ProcessEnded = func() {
			r.activeProcess.Store(nil)
		}
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
		r.agentRegistry.logWriter = logFile
	}

	// Redirect LLM retry logs to the log file.
	llm.SetLogOutput(logFile)

	return r
}

// setupSignalHandler installs handlers for lifecycle and peer-control signals (§2.2).
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
				if r.inferenceInFlight.Load() {
					r.log("SIGPIPE received during inference, marking recoverable")
					r.gracefulShutdownWithOutcome(
						RecoverableInferenceExitCode,
						"recoverable inference signal: SIGPIPE",
						tape.TermRecoverableInference,
					)
				}
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
	r.gracefulShutdownWithOutcome(exitCode, fmt.Sprintf("terminated by signal (exit %d)", exitCode), tape.TermSignal)
}

func (r *Runtime) gracefulShutdownWithOutcome(exitCode int, stderr string, mode tape.TerminationMode) {
	// Kill active child process to prevent orphans.
	// Use negative PID to kill the entire process group (children use Setpgid: true).
	if proc := r.activeProcess.Load(); proc != nil {
		_ = syscall.Kill(-proc.Pid, syscall.SIGKILL)
	}

	// Set session outcome if tape is initialized
	if r.tape != nil {
		r.flushPendingToolResult()
		duration := time.Since(r.startTime)
		r.tape.SetOutcome(tape.SessionOutcome{
			ExitCode:        exitCode,
			Stderr:          stderr,
			DurationMs:      duration.Milliseconds(),
			TerminationMode: mode,
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
		commitWorkspace := keepDetached || (mode == tape.TermSignal && r.cfg.WorkspaceCommitOnSignal)
		if err := r.sh.CloseWithOptions(keepDetached, commitWorkspace); err != nil {
			r.log("workspace close failed during signal shutdown: %v", err)
		}
	}

	r.quiesceRuntimeOwners(true)
	_ = r.cleanupAgentRoot()

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
	if err := r.recoverPendingWorkspaceCommit(); err != nil {
		r.logError("workspace commit recovery failed: %v", err)
		return 1
	}
	setupFailed := true
	registered := false
	agentRootSynced := false
	pidPublished := false
	defer func() {
		if !setupFailed {
			return
		}
		r.stopPeerDiscoveryHeartbeat()
		if pidPublished && r.agentRegistry != nil {
			if err := r.agentRegistry.UnpublishSelfPID(); err != nil {
				r.log("agent pid unpublish error after setup failure: %v", err)
			}
		}
		if registered && r.agentRegistry != nil {
			if err := r.agentRegistry.Deregister(); err != nil {
				r.log("agent deregistration error after setup failure: %v", err)
			}
		}
		if agentRootSynced {
			if err := r.cleanupAgentRoot(); err != nil {
				r.log("agent-root cleanup error: %v", err)
			}
		}
	}()

	// Register before syncing the shared session root so concurrent resume of a
	// live session fails before touching that session's durable surface.
	if err := r.agentRegistry.Register(); err != nil {
		r.logError("agent registration failed: %v", err)
		return 1
	}
	registered = true
	if err := r.importRetainedStateFromAgentRootIfNeeded(); err != nil {
		r.logError("agent-root retained-state import failed: %v", err)
		return 1
	}
	if err := r.loadRetainedControlState(); err != nil {
		r.logError("control state recovery failed: %v", err)
		return 1
	}
	if err := r.bootstrapAgentRoot(); err != nil {
		r.logError("agent-root init failed: %v", err)
		return 1
	}
	agentRootSynced = true
	if err := r.agentRegistry.PublishSelfPID(); err != nil {
		r.logError("agent pid publish failed: %v", err)
		return 1
	}
	pidPublished = true
	if err := r.importBootstrappedContext(); err != nil {
		r.logError("context bootstrap import failed: %v", err)
		return 1
	}
	setupFailed = false

	// Initialize exec executor now that we have the original input
	r.exec = tools.NewExecExecutor(r.cfg, mission)

	// Close the persistent shell when Run exits.
	// On success (exit code 0), spare detached jobs so daemons survive.
	// On failure (non-zero), kill everything to be safe.
	// gracefulShutdown handles the signal-based exit path separately.
	runExitCode := 1 // default: kill everything (all error paths return 1)
	defer func() {
		if r.shAlreadyFinalized() {
			return
		}
		if err := r.sh.Close(runExitCode == 0); err != nil {
			r.log("sh close error: %v", err)
			r.logError("sh close error: %v", err)
		}
	}()

	// Close the operational log file when Run exits.
	if r.logFile != nil {
		defer r.logFile.Close()
	}
	defer func() {
		r.quiesceRuntimeOwners(pidPublished)
		if err := r.cleanupAgentRoot(); err != nil {
			r.log("agent-root cleanup error: %v", err)
		}
	}()

	// Initialize tape
	r.tape = tape.NewTape(r.cfg.SessionID, r.cfg.ParentSession, r.cfg.Depth, r.cfg.ModelID)

	// Initialize legacy trace writer for JSONL persistence (§10).
	// Live provider input is assembled from context/prompt and context/state, not from this sidecar.
	tw, err := tape.NewWriter(r.cfg.TapeDir(""), r.cfg.TapeID)
	if err != nil {
		r.log("failed to create tape writer: %v", err)
	} else {
		r.tapeWriter = tw
		defer r.tapeWriter.Close()
	}

	// Write meta entry
	r.writeTapeEntry(r.tape.MetaEntry())

	// Build and append system prompt from the unified context surface.
	// The synthetic initial user message ("Begin." or QUINE_INITIAL_USER_MESSAGE)
	// is not real piped material, so it must not flip hasMaterial (which would
	// otherwise emit the fd3 material block on the standard surface).
	syntheticInitial := "Begin."
	if r.cfg != nil && r.cfg.InitialUserMessage != "" {
		syntheticInitial = r.cfg.InitialUserMessage
	}
	r.hasMaterial = material != syntheticInitial
	systemPrompt, err := r.currentSystemPrompt()
	if err != nil {
		r.logError("system prompt assembly failed: %v", err)
		return 1
	}
	systemMsg := tape.Message{
		Role:    tape.RoleSystem,
		Content: systemPrompt,
	}
	r.tape.Append(systemMsg)
	r.writeTapeEntry(tape.MessageEntry(systemMsg))

	// Append user message (material from stdin). In experimental no-Begin
	// sessions, omit the synthetic no-stdin initial turn entirely.
	if !(r.cfg.SuppressInitialBegin && material == syntheticInitial) {
		userMsg := tape.Message{
			Role:    tape.RoleUser,
			Content: material,
		}
		r.tape.Append(userMsg)
		r.writeTapeEntry(tape.MessageEntry(userMsg))
		if !r.contextHasVisibleState() {
			r.appendContextMessage(userMsg)
		}
	}

	if err := r.setupControlSurface(); err != nil {
		r.logError("control surface setup failed: %v", err)
		return 1
	}
	if code, done := r.recoverIncompleteToolBatches(); done {
		runExitCode = code
		return code
	}
	if r.agentRegistry != nil && r.cfg != nil && r.cfg.PeerDiscoveryEnabled {
		r.startPeerDiscoveryHeartbeat()
	}

	r.log("session started (depth=%d, model=%s)", r.cfg.Depth, r.cfg.ModelID)
	if len(r.tools) > 0 {
		names := make([]string, 0, len(r.tools))
		for _, tool := range r.tools {
			names = append(names, tool.Name)
		}
		r.log("tools available: %s", strings.Join(names, ","))
	}

	// Install signal handler for graceful shutdown (§7.3)
	r.setupSignalHandler()
	var wallClockExitTimer *time.Timer
	if r.cfg.WallClockExitSeconds > 0 {
		deadline := time.Duration(r.cfg.WallClockExitSeconds) * time.Second
		wallClockExitTimer = time.AfterFunc(deadline, func() {
			r.log("wall-clock self-exit deadline reached after %s", deadline)
			r.gracefulShutdownWithOutcome(
				0,
				fmt.Sprintf("wall-clock self-exit after %s", deadline),
				tape.TermSignal,
			)
		})
		defer wallClockExitTimer.Stop()
	}

	// Turn loop
	for {
		r.flushPendingToolResult()

		// SIGALRM panic mode (§2.2): inject a system override message
		// forcing the agent to exit with its best current answer.
		if r.panicMode.Load() {
			r.log("panic mode active, injecting system override")
			panicContent := "System interrupt: Time limit reached. Stop reasoning. Output your best current answer immediately using the exit tool. You MUST call exit now."
			if !r.cfg.ExitEnabled() {
				panicContent = "System interrupt: Time limit reached. Write any final output to >&4 now. The process will terminate."
			}
			panicMsg := tape.Message{
				Role:    tape.RoleUser,
				Content: panicContent,
			}
			r.tape.Append(panicMsg)
			r.writeTapeEntry(tape.MessageEntry(panicMsg))
			r.appendContextMessage(panicMsg)
		}

		// Abort before the next provider call if an integrity-critical entry
		// (tool_use turn / tool_result) failed to persist to the context
		// surface; continuing would send a context with broken tool pairing.
		if r.contextIntegrityErr != nil {
			code := r.handleError(fmt.Errorf("context surface integrity write failed: %w", r.contextIntegrityErr))
			runExitCode = code
			return code
		}

		if status, ok := r.memoryDeathStatus(); ok && status.Exceeded {
			code := r.terminateContextDeath(status)
			runExitCode = code
			return code
		}

		// Acquire concurrency slot before calling the LLM (§8.2)
		if err := r.semaphore.Acquire(); err != nil {
			r.log("semaphore acquire failed: %v", err)
		}

		// 1. Call provider.Generate. Streaming-capable providers route through
		// generateStreaming, which emits in-flight deltas to the transient
		// live.jsonl surface while crystallizing the same final cell. The live
		// writes happen inside the held semaphore slot (acquired above), so they
		// stay within the generate critical section.
		var (
			assistantMsg tape.Message
			usage        llm.Usage
			err          error
		)
		if sp, ok := r.provider.(llm.StreamingProvider); ok {
			assistantMsg, usage, err = r.generateStreaming(sp)
		} else {
			assistantMsg, usage, err = r.generateProviderResponse()
		}

		// Release concurrency slot immediately after LLM returns
		if releaseErr := r.semaphore.Release(); releaseErr != nil {
			r.log("semaphore release failed: %v", releaseErr)
		}

		if err != nil {
			code := r.handleError(err)
			runExitCode = code
			return code
		}

		readyTextAutoIdle := r.shouldReadyTextAutoIdle(assistantMsg)

		// 2. Append assistant message to Tape
		r.tape.Append(assistantMsg)
		r.writeTapeEntry(tape.MessageEntry(assistantMsg))
		if !readyTextAutoIdle {
			r.appendContextMessage(assistantMsg)
		}

		// 3. Accumulate usage
		r.tape.AddUsage(usage.InputTokens, usage.OutputTokens)

		if isAssistantNoOp(assistantMsg) {
			if r.canTreatEmptyAssistantAsSuccess() {
				code := r.finishEmptyAssistantSuccess()
				runExitCode = code
				return code
			}
			code := r.handleError(fmt.Errorf("provider returned an empty assistant response with no tool calls"))
			runExitCode = code
			return code
		}

		// 4. Inspect assistant message
		if len(assistantMsg.ToolCalls) == 0 {
			if r.shouldTerminateChildTextOnly(assistantMsg) {
				return r.terminateChildTextOnly(assistantMsg)
			}
			if readyTextAutoIdle {
				r.handleReadyTextAutoIdle()
				continue
			}
			// Text-only response: continue loop for next inference.
			continue
		}

		// Process tool calls sequentially
		for tcIdx, tc := range assistantMsg.ToolCalls {
			r.flushPendingToolResult()

			// In panic mode, reject any tool call that isn't exit (§2.2).
			// When exit is disabled, panic mode terminates immediately.
			if r.panicMode.Load() {
				if !r.cfg.ExitEnabled() {
					r.log("panic mode with exit disabled: terminating immediately")
					return r.terminateExecutionBudgetExhausted()
				}
				if tc.Name != "exit" {
					rejectMsg := runtimeToolResultMessage(tc.ID, tc.Name, "rejected", map[string]any{
						"error": "Rejected: time limit reached (SIGALRM). You MUST call exit immediately with your best current answer.",
					})
					r.appendRuntimeToolMessage(rejectMsg, true)
					r.log("panic mode: rejected non-exit tool call %q", tc.Name)
					continue
				}
			}

			// Reject a tool call whose arguments the provider sent but could
			// not be decoded as JSON (e.g. truncated at an output limit). The
			// protocol decoder flags this instead of silently emptying the
			// arguments, which would otherwise run as a success-shaped no-op.
			if _, malformed := tc.MalformedArguments(); malformed {
				rejectMsg := runtimeToolResultMessage(tc.ID, tc.Name, "rejected", map[string]any{
					"error": "Rejected: tool arguments could not be parsed as JSON (they may have been truncated by an output limit). Resend the tool call with complete, valid JSON arguments.",
				})
				r.appendRuntimeToolMessage(rejectMsg, true)
				r.log("%s rejected: malformed/undecodable arguments", tc.Name)
				continue
			}

			spec, ok := r.toolRegistry[tc.Name]
			if !ok {
				// Unknown tool — return error result
				unknownMsg := runtimeToolResultMessage(tc.ID, tc.Name, "error", map[string]any{
					"error": fmt.Sprintf("unknown tool: %s", tc.Name),
				})
				r.appendRuntimeToolMessage(unknownMsg, true)
				continue
			}
			if spec.enabled != nil && !spec.enabled(r.cfg) {
				r.rejectDisabledTool(tc, spec)
				continue
			}

			outcome := spec.handle(r, assistantMsg, tc)
			if outcome.terminate {
				runExitCode = outcome.exitCode
				return outcome.exitCode
			}
			if outcome.budgetExhausted {
				code := r.handleExecutionBudgetExhaustion(assistantMsg.ToolCalls[tcIdx+1:])
				runExitCode = code
				return code
			}
		}

		// Inject resource usage into the last tool result so the agent
		// sees its situation without breaking tool_use/tool_result pairing.
		if last := r.pendingToolResultMessage(); last != nil {
			var deathStatus tools.MemoryDeathStatus
			deathExceeded := false
			extraEntry, hasExtraEntry := r.pendingToolResultEntry(last)
			if hasExtraEntry {
				r.appendMemoryStatus(last, extraEntry)
				if status, ok := r.memoryDeathStatus(extraEntry); ok && status.Exceeded {
					deathStatus = status
					deathExceeded = true
				}
			} else {
				r.appendMemoryStatus(last)
				if status, ok := r.memoryDeathStatus(); ok && status.Exceeded {
					deathStatus = status
					deathExceeded = true
				}
			}
			if updated, err := setRuntimeField(last.StructuredContent, "context_used_k", usage.InputTokens/1000); err == nil {
				last.StructuredContent = updated
				syncToolResultMessageContent(last)
			} else {
				r.log("tool result context update error: %v", err)
			}
			if deathExceeded {
				code := r.terminateContextDeath(deathStatus)
				runExitCode = code
				return code
			}
		}
	}
}

func (r *Runtime) shouldTerminateChildTextOnly(msg tape.Message) bool {
	return r.cfg != nil &&
		r.cfg.Depth > 0 &&
		strings.TrimSpace(r.cfg.ParentSession) != "" &&
		strings.TrimSpace(msg.Content) != ""
}

func (r *Runtime) terminateChildTextOnly(msg tape.Message) int {
	if r.stdout != nil && msg.Content != "" {
		_, _ = fmt.Fprint(r.stdout, msg.Content)
	}

	duration := time.Since(r.startTime)
	totalTokens := r.tape.TokensIn + r.tape.TokensOut
	r.log("child text-only response terminated as success (%d turns, %.1fs, %d tokens)",
		r.tape.TurnCount, duration.Seconds(), totalTokens)

	return r.finalizeOutcome(0, "", tape.TermExit)
}

// handleExecutionBudgetExhaustion rejects any remaining same-batch tool calls,
// exposes the zero-budget edge in the last sh result, then grants one final
// inference in which only exit is accepted (when exit is enabled).
// When exit is disabled, terminates directly after rejecting remaining calls.
func (r *Runtime) handleExecutionBudgetExhaustion(remaining []tape.ToolCall) int {
	// Reject remaining tool calls in this batch to maintain protocol pairing.
	for _, remainingTC := range remaining {
		rejectMsg := runtimeToolResultMessage(remainingTC.ID, remainingTC.Name, "rejected", map[string]any{
			"error": "Rejected: execution budget exhausted before this tool call could be processed.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
	}

	// When exit is disabled, terminate directly — no final inference needed.
	if !r.cfg.ExitEnabled() {
		return r.terminateExecutionBudgetExhausted()
	}

	// Expose the zero-budget edge explicitly in the last tool result, then
	// grant one final inference in which only exit is accepted.
	if last := r.pendingToolResultMessage(); last != nil {
		if updated, err := setRuntimeField(last.StructuredContent, "no_turns_left", "0 turns remain. Do not call `sh`. Call `exit` now."); err == nil {
			last.StructuredContent = updated
			syncToolResultMessageContent(last)
		} else {
			r.log("tool result no-turns-left update error: %v", err)
		}
	}

	if err := r.semaphore.Acquire(); err != nil {
		r.log("semaphore acquire failed (exit-only continuation): %v", err)
	}
	finalMsg, finalUsage, err := r.generateProviderResponse()
	if releaseErr := r.semaphore.Release(); releaseErr != nil {
		r.log("semaphore release failed (exit-only continuation): %v", releaseErr)
	}
	if err != nil {
		code := r.handleError(err)
		return code
	}
	r.tape.Append(finalMsg)
	r.writeTapeEntry(tape.MessageEntry(finalMsg))
	r.appendContextMessage(finalMsg)
	r.tape.AddUsage(finalUsage.InputTokens, finalUsage.OutputTokens)

	for _, tc := range finalMsg.ToolCalls {
		if tc.Name == "exit" {
			r.log("execution-exhausted continuation: exit selected")
			code, ok := r.handleExit(tc)
			if ok {
				return code
			}
			continue
		}

		rejectMsg := runtimeToolResultMessage(tc.ID, tc.Name, "rejected", map[string]any{
			"error": "Rejected: no turns left. Only exit is accepted after execution budget is exhausted.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
	}
	return r.terminateExecutionBudgetExhausted()
}

func (r *Runtime) terminateExecutionBudgetExhausted() int {
	r.flushPendingToolResult()
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

func (r *Runtime) memoryDeathStatus(extraEntries ...tape.TapeEntry) (tools.MemoryDeathStatus, bool) {
	if r.memory == nil || r.cfg == nil || r.cfg.MemoryDeathTokens <= 0 {
		return tools.MemoryDeathStatus{}, false
	}
	status, err := r.memory.DeathStatusWithEntries(extraEntries...)
	if err != nil {
		r.log("memory death check error: %v", err)
		return tools.MemoryDeathStatus{}, false
	}
	return status, status.DeathTokens > 0
}

func (r *Runtime) terminateContextDeath(status tools.MemoryDeathStatus) int {
	r.flushPendingToolResult()
	stderr := fmt.Sprintf("context death: frontier_estimated_tokens=%d >= death_tokens=%d",
		status.FrontierEstimatedTokens, status.DeathTokens)
	r.log("%s", stderr)
	r.logError("%s", stderr)
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        1,
		Stderr:          stderr,
		DurationMs:      duration.Milliseconds(),
		TerminationMode: tape.TermContextDeath,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())
	return 1
}

func (r *Runtime) generateProviderResponse() (tape.Message, llm.Usage, error) {
	r.flushPendingToolResult()
	messages, err := r.providerContextMessages()
	if err != nil {
		return tape.Message{}, llm.Usage{}, fmt.Errorf("assemble provider context: %w", err)
	}
	r.inferenceInFlight.Store(true)
	defer r.inferenceInFlight.Store(false)
	return r.provider.Generate(messages, r.tools)
}

// liveEntry is one line of the transient live-generation stream (live.jsonl).
// It is display-only: never provider input, never recovered.
type liveEntry struct {
	Seq      int    `json:"seq"`
	Kind     string `json:"kind"`
	Text     string `json:"text,omitempty"`
	ToolID   string `json:"tool_id,omitempty"`
	ToolName string `json:"tool_name,omitempty"`
	TS       int64  `json:"ts"`
}

// generateStreaming runs a streaming generation, mirroring generateProviderResponse
// for the crystallized cell while emitting in-flight deltas to live.jsonl. The
// returned (message, usage) is the authoritative final cell — byte-equivalent to
// the buffered path — and the live stream is pure transient display signal that
// the downstream tape/current.jsonl write path never reads.
func (r *Runtime) generateStreaming(sp llm.StreamingProvider) (tape.Message, llm.Usage, error) {
	r.flushPendingToolResult()
	messages, err := r.providerContextMessages()
	if err != nil {
		return tape.Message{}, llm.Usage{}, fmt.Errorf("assemble provider context: %w", err)
	}
	r.inferenceInFlight.Store(true)
	defer r.inferenceInFlight.Store(false)

	// Truncate the live surface at generation start so seq resets and no stale
	// deltas from a prior generation remain. Best-effort: live.jsonl is never a
	// correctness surface, so a truncate failure only degrades the preview.
	livePath := r.contextLivePath()
	if err := os.Remove(livePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		r.log("live surface truncate error: %v", err)
	}

	ch, err := sp.GenerateStream(messages, r.tools)
	if err != nil {
		return tape.Message{}, llm.Usage{}, err
	}

	var (
		finalMsg tape.Message
		usage    llm.Usage
		genErr   error
		seq      int
	)
	for ev := range ch {
		switch ev.Kind {
		case llm.StreamText:
			seq++
			r.appendLiveEntry(livePath, liveEntry{Seq: seq, Kind: "text_delta", Text: ev.Text})
		case llm.StreamReasoning:
			seq++
			r.appendLiveEntry(livePath, liveEntry{Seq: seq, Kind: "reasoning_delta", Text: ev.Text})
		case llm.StreamToolCall:
			seq++
			r.appendLiveEntry(livePath, liveEntry{Seq: seq, Kind: "tool_call_delta", Text: ev.Text, ToolID: ev.ToolID, ToolName: ev.ToolName})
		case llm.StreamDone:
			seq++
			r.appendLiveEntry(livePath, liveEntry{Seq: seq, Kind: "completed"})
			finalMsg = ev.Message
			usage = ev.Usage
		case llm.StreamError:
			genErr = ev.Err
		}
	}
	if genErr != nil {
		return tape.Message{}, llm.Usage{}, genErr
	}
	return finalMsg, usage, nil
}

// appendLiveEntry writes one delta line to the transient live surface. Failures
// are non-fatal and only logged: live.jsonl is display-only and a lost line can
// never corrupt provider input or recovery (unlike appendContextEntry).
func (r *Runtime) appendLiveEntry(path string, entry liveEntry) {
	entry.TS = time.Now().UnixMilli()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		r.log("live surface mkdir error: %v", err)
		return
	}
	line, err := json.Marshal(entry)
	if err != nil {
		r.log("live surface marshal error: %v", err)
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		r.log("live surface open error: %v", err)
		return
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		r.log("live surface write error: %v", err)
	}
}

func (r *Runtime) providerContextMessages() ([]tape.Message, error) {
	r.flushPendingToolResult()
	systemPrompt, err := r.currentSystemPrompt()
	if err != nil {
		return nil, err
	}
	msgs := []tape.Message{{
		Role:    tape.RoleSystem,
		Content: systemPrompt,
	}}

	if r.cfg != nil && r.cfg.AnchorMemoryEnabled {
		summaries, err := r.anchorSummaryMessages()
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, summaries...)
	}

	current, err := readContextMessages(r.contextCurrentPath())
	if err != nil {
		return nil, err
	}
	msgs = append(msgs, current...)
	return normalizeProviderToolMessages(msgs)
}

func (r *Runtime) contextCurrentPath() string {
	return filepath.Join(r.cfg.AgentRoot(), "context", "state", "current.jsonl")
}

// contextLivePath is the transient live-generation stream. It carries in-flight
// generation deltas for display only and is truncated at the start of every
// generation. It is NEVER provider input and NEVER recovered — see
// generateStreaming and the live-surface isolation validator.
func (r *Runtime) contextLivePath() string {
	return filepath.Join(r.cfg.AgentRoot(), "context", "state", "live.jsonl")
}

func (r *Runtime) contextFrontierRoot() string {
	return filepath.Join(r.cfg.AgentRoot(), "context", "state", "frontier")
}

func (r *Runtime) contextAnchorsRoot() string {
	return filepath.Join(r.cfg.AgentRoot(), "context", "state", "anchors")
}

func (r *Runtime) contextHasVisibleState() bool {
	if info, err := os.Stat(r.contextCurrentPath()); err == nil && info.Size() > 0 {
		return true
	}
	entries, err := os.ReadDir(r.contextFrontierRoot())
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".link") {
			return true
		}
	}
	return false
}

func (r *Runtime) importBootstrappedContext() error {
	bootstrapRoot := strings.TrimSpace(os.Getenv(tools.ContextBootstrapEnv))
	if bootstrapRoot == "" {
		return nil
	}
	_ = os.Unsetenv(tools.ContextBootstrapEnv)

	contextRoot := r.currentIncarnationContextRoot()
	if err := os.RemoveAll(contextRoot); err != nil {
		return fmt.Errorf("reset context root %s: %w", contextRoot, err)
	}
	if err := tools.CopyTreePreservingSymlinks(bootstrapRoot, contextRoot); err != nil {
		return fmt.Errorf("import %s -> %s: %w", bootstrapRoot, contextRoot, err)
	}
	r.log("imported bootstrapped context from %s", bootstrapRoot)
	return nil
}

func (r *Runtime) anchorSummaryMessages() ([]tape.Message, error) {
	entries, err := os.ReadDir(r.contextFrontierRoot())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	type anchorMeta struct {
		ID         int    `json:"id"`
		Resolution string `json:"resolution"`
	}

	ids := make([]int, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".link") {
			continue
		}
		id, err := strconv.Atoi(strings.TrimSuffix(entry.Name(), ".link"))
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	sort.Ints(ids)

	out := make([]tape.Message, 0, len(ids))
	for _, id := range ids {
		path := filepath.Join(r.contextAnchorsRoot(), fmt.Sprintf("%d.anchor", id), "meta.json")
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var meta anchorMeta
		if err := json.Unmarshal(data, &meta); err != nil {
			return nil, err
		}
		resolution := strings.TrimSpace(meta.Resolution)
		if resolution == "" {
			continue
		}
		out = append(out, tape.Message{
			Role:    tape.RoleAssistant,
			Content: fmt.Sprintf("[ANCHOR %d] %s", id, resolution),
		})
	}
	return out, nil
}

// readContextMessages assembles provider-input messages from the complete-cell
// context surface. Invariant (C2): it reads ONLY current.jsonl (contextCurrentPath);
// it must never read the transient live.jsonl (contextLivePath) — live deltas are
// partial cells that would corrupt tool_use/tool_result pairing.
func readContextMessages(path string) ([]tape.Message, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	lines := bytes.Split(data, []byte("\n"))
	out := make([]tape.Message, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var entry tape.TapeEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, err
		}
		switch entry.Type {
		case "message":
			var msg tape.Message
			if err := json.Unmarshal(entry.Data, &msg); err != nil {
				return nil, err
			}
			out = append(out, msg)
		case "tool_result":
			var result tape.ToolResult
			if err := json.Unmarshal(entry.Data, &result); err != nil {
				return nil, err
			}
			out = append(out, tape.Message{
				Role:              tape.RoleToolResult,
				Content:           "",
				StructuredContent: result.Content,
				ToolID:            result.ToolID,
				Image:             result.Image,
			})
		}
	}
	return out, nil
}

func normalizeProviderToolMessages(msgs []tape.Message) ([]tape.Message, error) {
	if len(msgs) == 0 {
		return nil, nil
	}

	out := make([]tape.Message, 0, len(msgs))
	var pending map[string]struct{}

	for _, msg := range msgs {
		switch {
		case msg.Role == tape.RoleAssistant && len(msg.ToolCalls) > 0:
			if len(pending) > 0 {
				return nil, fmt.Errorf("incomplete tool batch before assistant message: missing results for %s", strings.Join(sortedPendingToolIDs(pending), ", "))
			}
			pending = make(map[string]struct{}, len(msg.ToolCalls))
			for _, tc := range msg.ToolCalls {
				pending[tc.ID] = struct{}{}
			}
			out = append(out, msg)
		case msg.Role == tape.RoleToolResult:
			if pending == nil {
				continue
			}
			if _, ok := pending[msg.ToolID]; !ok {
				continue
			}
			out = append(out, msg)
			delete(pending, msg.ToolID)
			if len(pending) == 0 {
				pending = nil
			}
		default:
			if len(pending) > 0 {
				return nil, fmt.Errorf("incomplete tool batch before %s message: missing results for %s", msg.Role, strings.Join(sortedPendingToolIDs(pending), ", "))
			}
			out = append(out, msg)
		}
	}

	if len(pending) > 0 {
		return nil, fmt.Errorf("incomplete tool batch at context tail: missing results for %s", strings.Join(sortedPendingToolIDs(pending), ", "))
	}
	return out, nil
}

func sortedPendingToolIDs(pending map[string]struct{}) []string {
	ids := make([]string, 0, len(pending))
	for id := range pending {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
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

func (r *Runtime) appendToolResult(result tape.ToolResult) {
	r.flushPendingToolResult()
	if !tape.HasStructuredContent(result.Content) {
		result.Content = result.StructuredContent
	}
	result.StructuredContent = nil
	msg := tape.Message{
		Role:              tape.RoleToolResult,
		Content:           "",
		StructuredContent: result.Content,
		ToolID:            result.ToolID,
		Image:             result.Image,
	}
	r.tape.Append(msg)
	r.pendingTool = &pendingToolResult{result: result}
}

func (r *Runtime) appendRuntimeToolMessage(msg tape.Message, isError bool) {
	r.flushPendingToolResult()
	r.tape.Append(msg)
	r.pendingTool = &pendingToolResult{result: tape.ToolResult{
		ToolID:  msg.ToolID,
		Content: msg.StructuredContent,
		IsError: isError,
		Image:   msg.Image,
	}}
}

func (r *Runtime) pendingToolResultMessage() *tape.Message {
	if r.pendingTool == nil || r.tape == nil {
		return nil
	}
	last := r.tape.LastMessage()
	if last == nil || last.Role != tape.RoleToolResult {
		r.log("pending tool result missing last tool message")
		return nil
	}
	if last.ToolID != r.pendingTool.result.ToolID {
		r.log("pending tool result mismatch: pending=%s last=%s", r.pendingTool.result.ToolID, last.ToolID)
		return nil
	}
	return last
}

func (r *Runtime) flushPendingToolResult() {
	if r.pendingTool == nil {
		return
	}
	pending := r.pendingTool
	r.pendingTool = nil

	msg := r.pendingToolResultMessageFor(pending)
	if msg == nil {
		return
	}
	r.appendControlStatus(msg)
	r.appendPeerDiscoveryStatus(msg)

	result := pending.result
	result.ToolID = msg.ToolID
	result.Content = msg.StructuredContent
	result.StructuredContent = nil
	result.Image = msg.Image
	r.writeTapeEntry(tape.ToolResultEntry(result))
	r.appendContextToolResult(result)
}

func (r *Runtime) pendingToolResultMessageFor(pending *pendingToolResult) *tape.Message {
	if pending == nil || r.tape == nil {
		return nil
	}
	last := r.tape.LastMessage()
	if last == nil || last.Role != tape.RoleToolResult {
		r.log("pending tool result missing last tool message")
		return nil
	}
	if last.ToolID != pending.result.ToolID {
		r.log("pending tool result mismatch: pending=%s last=%s", pending.result.ToolID, last.ToolID)
		return nil
	}
	return last
}

func (r *Runtime) pendingToolResultEntry(msg *tape.Message) (tape.TapeEntry, bool) {
	if r.pendingTool == nil || msg == nil || msg.Role != tape.RoleToolResult || msg.ToolID != r.pendingTool.result.ToolID {
		return tape.TapeEntry{}, false
	}
	result := r.pendingTool.result
	result.ToolID = msg.ToolID
	result.Content = msg.StructuredContent
	result.StructuredContent = nil
	result.Image = msg.Image
	return tape.ToolResultEntry(result), true
}

func (r *Runtime) appendContextMessage(msg tape.Message) {
	if msg.Role == tape.RoleSystem {
		return
	}
	if err := r.appendContextEntry(tape.MessageEntry(msg)); err != nil && len(msg.ToolCalls) > 0 {
		// A lost assistant tool_use entry orphans the tool calls on the
		// provider-input surface; escalate so the loop fails rather than
		// silently dropping the turn.
		r.recordContextIntegrityFailure(err)
	}
}

func (r *Runtime) appendContextToolResult(result tape.ToolResult) {
	if err := r.appendContextEntry(tape.ToolResultEntry(result)); err != nil {
		// A lost tool_result breaks tool_use/tool_result pairing on the next
		// request; escalate rather than send a context the provider rejects.
		r.recordContextIntegrityFailure(err)
	}
}

// recordContextIntegrityFailure latches the first integrity-critical context
// write failure so the turn loop can abort before the next provider call.
func (r *Runtime) recordContextIntegrityFailure(err error) {
	if r.contextIntegrityErr == nil {
		r.contextIntegrityErr = err
	}
}

func (r *Runtime) appendContextEntry(entry tape.TapeEntry) error {
	var err error
	if r.memory != nil {
		err = r.memory.AppendEntry(entry)
	} else {
		err = appendJSONL(r.contextCurrentPath(), entry)
	}
	if err != nil {
		r.log("context append error: %v", err)
	}
	return err
}

func appendJSONL(path string, entry tape.TapeEntry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(line); err != nil {
		return err
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) bootstrapAgentRoot() error {
	if err := r.ensureIncarnation(); err != nil {
		return err
	}

	agentRoot := r.cfg.AgentRoot()
	retainedRoot := r.cfg.SessionRetainedDir("")
	retainedIncarnationDir := r.cfg.SessionIncarnationDir("")
	currentIncarnationLink := r.cfg.SessionCurrentIncarnationPath("")
	currentIncarnationRoot := r.currentIncarnationRoot()
	currentContextDir := r.currentIncarnationContextRoot()
	currentPromptDir := r.currentIncarnationPromptRoot()
	currentStateDir := r.currentIncarnationStateRoot()
	currentMissionPath := r.currentIncarnationMissionPath()
	retainedStatusDir := filepath.Join(retainedRoot, "status")
	retainedWorldDir := filepath.Join(retainedRoot, "world")
	retainedTapesDir := r.cfg.TapeDir("")
	jobSessionDir := r.cfg.JobSessionDir("")
	retainedMissionPath := filepath.Join(retainedRoot, "mission.txt")
	retainedContextPath := filepath.Join(retainedRoot, "context")
	sessionStatusDir := r.cfg.SessionStatusDir("")
	sessionWorldDir := r.cfg.SessionWorldDir("")
	contextDir := filepath.Join(agentRoot, "context")
	contextFrontierDir := filepath.Join(currentStateDir, "frontier")
	contextAnchorsDir := filepath.Join(currentStateDir, "anchors")
	publicDir := filepath.Join(agentRoot, "public")
	missionPath := r.cfg.SessionMissionPath("")
	runtimeLogCompatRoot := filepath.Join(r.cfg.RuntimeRoot(), "log")
	tapesCompatRoot := filepath.Join(r.cfg.RuntimeRoot(), "tapes")
	compatLogPath := filepath.Join(runtimeLogCompatRoot, r.cfg.SessionID)
	compatTapesPath := filepath.Join(tapesCompatRoot, r.cfg.SessionID)
	contextFresh := false
	if _, err := os.Stat(currentContextDir); errors.Is(err, os.ErrNotExist) {
		contextFresh = true
	} else if err != nil {
		return fmt.Errorf("stat context root %s: %w", currentContextDir, err)
	}

	for _, dir := range []string{
		agentRoot,
		retainedRoot,
		retainedIncarnationDir,
		currentIncarnationRoot,
		currentContextDir,
		currentPromptDir,
		currentStateDir,
		contextFrontierDir,
		contextAnchorsDir,
		publicDir,
		retainedStatusDir,
		retainedWorldDir,
		retainedTapesDir,
		jobSessionDir,
		runtimeLogCompatRoot,
		tapesCompatRoot,
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	if err := replaceRelativeSymlink(filepath.Join(agentRoot, "log"), retainedRoot); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(filepath.Join(agentRoot, "inc"), retainedIncarnationDir); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(sessionStatusDir, retainedStatusDir); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(sessionWorldDir, retainedWorldDir); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(filepath.Join(agentRoot, "jobs"), jobSessionDir); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(filepath.Join(agentRoot, "tapes"), retainedTapesDir); err != nil {
		return err
	}
	if err := r.ensureControlSurface(); err != nil {
		return fmt.Errorf("sync control surface: %w", err)
	}
	if err := replaceRelativeSymlinkAtomic(currentIncarnationLink, strconv.Itoa(r.cfg.IncarnationID)); err != nil {
		return fmt.Errorf("sync current incarnation link: %w", err)
	}
	if contextFresh {
		if err := r.copyInitialIncarnationContext(currentContextDir); err != nil {
			return fmt.Errorf("copy initial incarnation context: %w", err)
		}
	}
	if err := replaceRelativeSymlink(contextDir, filepath.Join(agentRoot, "inc", "current", "context")); err != nil {
		return err
	}
	if err := writeTextFile(currentMissionPath, r.originalInput); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(missionPath, filepath.Join(agentRoot, "inc", "current", "mission.txt")); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(retainedMissionPath, filepath.Join(retainedRoot, "inc", "current", "mission.txt")); err != nil {
		return err
	}
	if err := replaceRelativeSymlink(retainedContextPath, filepath.Join(retainedRoot, "inc", "current", "context")); err != nil {
		return err
	}
	if err := r.syncPromptFragments(); err != nil {
		return fmt.Errorf("sync context files: %w", err)
	}
	if err := r.syncSelfSourceSurface(agentRoot); err != nil {
		return fmt.Errorf("sync self-source surface: %w", err)
	}
	if err := removeSelfSourceSurface(filepath.Join(agentRoot, "genome")); err != nil {
		return fmt.Errorf("remove stale genome surface: %w", err)
	}
	sourceRoot := ""
	if r.cfg.SelfSourceCodeEnabled {
		sourceRoot = filepath.Join(agentRoot, "source-code")
	}
	if err := r.syncPublicSurface(publicSurfacePaths{
		PublicDir:        publicDir,
		StatusTarget:     sessionStatusDir,
		ControlLogTarget: r.cfg.SessionControlLogPath(""),
		SourceRoot:       sourceRoot,
	}); err != nil {
		return fmt.Errorf("sync public surface: %w", err)
	}
	if strings.TrimSpace(r.cfg.RetentionDir) != "" {
		if err := replaceRelativeSymlink(compatLogPath, retainedRoot); err != nil {
			return err
		}
		if err := replaceRelativeSymlink(compatTapesPath, retainedTapesDir); err != nil {
			return err
		}
	}
	_ = os.Remove(filepath.Join(retainedRoot, "current.jsonl"))
	_ = os.Remove(filepath.Join(retainedRoot, "current.log.yaml"))
	if err := os.RemoveAll(filepath.Join(retainedRoot, "jobs")); err != nil {
		return fmt.Errorf("remove retained jobs projection: %w", err)
	}

	if err := r.syncAgentStatusSurfaces(); err != nil {
		return err
	}
	r.agentRootBootstrapped = true
	return nil
}

func (r *Runtime) syncAgentStatusSurfaces() error {
	if err := r.syncSessionStatusAndWorldSurfaces(); err != nil {
		return err
	}

	return r.writeControlInboxSnapshot(r.controlSnapshot())
}

func (r *Runtime) syncSessionStatusAndWorldSurfaces() error {
	if err := r.writeSessionStatus(); err != nil {
		return err
	}

	retainedStatusDir := filepath.Join(r.cfg.SessionRetainedDir(""), "status")
	if err := r.writeRuntimeContract(retainedStatusDir); err != nil {
		return err
	}

	retainedWorldDir := filepath.Join(r.cfg.SessionRetainedDir(""), "world")
	if err := r.syncInspectableWorldSurface(retainedWorldDir); err != nil {
		return err
	}

	return nil
}

func (r *Runtime) writeSessionStatus() error {
	retainedStatusDir := filepath.Join(r.cfg.SessionRetainedDir(""), "status")
	if err := ensureDirPath(retainedStatusDir, "session status dir"); err != nil {
		return err
	}

	// session.json: core identity fields only. Derived paths (ctl, inbox,
	// control_log) are omitted; peers can infer them from agent_root using the
	// standard public layout.
	status := map[string]any{
		"session_id":            r.cfg.SessionID,
		"run_id":                r.cfg.RunID,
		"incarnation_id":        r.cfg.IncarnationID,
		"pid":                   os.Getpid(),
		"ppid":                  os.Getppid(),
		"parent_session":        r.cfg.ParentSession,
		"depth":                 r.cfg.Depth,
		"model_id":              r.cfg.ModelID,
		"runtime_root":          r.cfg.RuntimeRoot(),
		"agent_root":            r.cfg.AgentRoot(),
		"workspace_enabled":     r.cfg.WorkspaceEnabled,
		"workspace_root":        r.cfg.WorkspaceRoot,
		"workspace":             r.cfg.Workspace,
		"workspace_session":     r.cfg.WorkspaceSession,
		"anchor_memory_enabled": r.cfg.AnchorMemoryEnabled,
		"retention_dir":         r.cfg.RetentionDir,
	}
	return writeJSONFile(filepath.Join(retainedStatusDir, "session.json"), status)
}

func (r *Runtime) copyPriorIncarnationContext(dst string) error {
	if r.cfg.IncarnationID <= 0 {
		return nil
	}
	src := filepath.Join(r.cfg.SessionIncarnationPath("", r.cfg.IncarnationID-1), "context")
	info, err := os.Stat(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("prior context root %q is not a directory", src)
	}
	return tools.CopyTreePreservingSymlinks(src, dst)
}

func (r *Runtime) copyInitialIncarnationContext(dst string) error {
	if r.cfg.IncarnationID > 0 {
		return r.copyPriorIncarnationContext(dst)
	}
	return r.copyRetainedSeedContext(dst)
}

func (r *Runtime) copyRetainedSeedContext(dst string) error {
	src := filepath.Join(r.cfg.SessionRetainedDir(""), "seed", "context")
	info, err := os.Stat(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("retained seed context root %q is not a directory", src)
	}
	return tools.CopyTreePreservingSymlinks(src, dst)
}

func (r *Runtime) cleanupAgentRoot() error {
	agentRoot := strings.TrimSpace(r.cfg.AgentRoot())
	if agentRoot == "" {
		return nil
	}
	if r.finalizationStarted() {
		r.recordFinalizationPhase(finalizationPhaseSurfaceCleanup)
	}
	if err := removeSelfSourceSurface(filepath.Join(agentRoot, "source-code")); err != nil {
		return fmt.Errorf("remove self-source surface: %w", err)
	}
	if err := removeSelfSourceSurface(filepath.Join(agentRoot, "genome")); err != nil {
		return fmt.Errorf("remove stale genome surface: %w", err)
	}
	if err := r.cleanupPublicSurface(); err != nil {
		return err
	}
	if err := makeTreeRemovable(agentRoot); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("prepare agent root %s for removal: %w", agentRoot, err)
	}
	if err := os.RemoveAll(agentRoot); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove agent root %s: %w", agentRoot, err)
	}
	r.agentRootBootstrapped = false
	return nil
}

func makeTreeRemovable(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		mode := os.FileMode(0o644)
		if d.IsDir() {
			mode = 0o755
		}
		if err := os.Chmod(path, mode); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("chmod %s to %o: %w", path, mode, err)
		}
		return nil
	})
}

func copyFile(src string, dst string) error {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("abs source path %q: %w", src, err)
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("abs destination path %q: %w", dst, err)
	}
	if srcAbs == dstAbs {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return fmt.Errorf("stat source file %q: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("open destination file %q: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %q -> %q: %w", src, dst, err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("sync destination file %q: %w", dst, err)
	}
	return nil
}

func hardlinkOrCopyFile(src string, dst string) error {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("abs source path %q: %w", src, err)
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("abs destination path %q: %w", dst, err)
	}
	if srcAbs == dstAbs {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
	}
	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("remove destination %q: %w", dst, err)
	}
	if err := os.Link(src, dst); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrExist) {
		if err := copyFile(src, dst); err != nil {
			return err
		}
		return nil
	}
	if err := copyFile(src, dst); err != nil {
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
	var lastErr error
	for attempt := 0; attempt < 64; attempt++ {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
		if err := os.Symlink(target, path); err == nil {
			return nil
		} else {
			lastErr = err
			if existing, readErr := os.Readlink(path); readErr == nil && existing == target {
				return nil
			}
			if !errors.Is(err, os.ErrExist) {
				return fmt.Errorf("symlink %s -> %s: %w", path, target, err)
			}
		}
	}
	if existing, readErr := os.Readlink(path); readErr == nil && existing == target {
		return nil
	}
	return fmt.Errorf("symlink %s -> %s: %w", path, target, lastErr)
}

func replaceRelativeSymlink(path string, target string) error {
	linkDir := filepath.Dir(path)
	if !filepath.IsAbs(linkDir) {
		absLinkDir, err := filepath.Abs(linkDir)
		if err != nil {
			return fmt.Errorf("resolve symlink parent %s: %w", linkDir, err)
		}
		linkDir = absLinkDir
	}
	if !filepath.IsAbs(target) {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			return fmt.Errorf("resolve symlink target %s: %w", target, err)
		}
		target = absTarget
	}
	relTarget, err := filepath.Rel(linkDir, target)
	if err != nil {
		return fmt.Errorf("relativize symlink target %s from %s: %w", target, linkDir, err)
	}
	return replaceRelativeSymlinkAtomic(path, relTarget)
}

func replaceRelativeSymlinkAtomic(path string, target string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir parent for %s: %w", path, err)
	}
	tmpPath := filepath.Join(filepath.Dir(path), "."+filepath.Base(path)+".tmp")
	if err := os.RemoveAll(tmpPath); err != nil {
		return fmt.Errorf("remove temp symlink %s: %w", tmpPath, err)
	}
	if err := os.Symlink(target, tmpPath); err != nil {
		return fmt.Errorf("symlink %s -> %s: %w", tmpPath, target, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename %s -> %s: %w", tmpPath, path, err)
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

// writeFile writes the full file contents atomically: it stages the data in a
// temp file in the same directory and renames it into place. Cross-process
// surfaces written here (e.g. status/session.json) are read concurrently by peer
// processes; a plain os.WriteFile truncates in place, so a reader can observe an
// empty or partial file mid-write and decode a half-populated record. The
// rename makes every read see either the complete old or complete new content.
func writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	tmpName := tmp.Name()
	cleanup := func() {
		tmp.Close()
		os.Remove(tmpName)
	}
	if _, err := tmp.Write(data); err != nil {
		cleanup()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := tmp.Chmod(0o644); err != nil {
		cleanup()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
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
	updated, err := setRuntimeField(result.Content, "memory_feedback", unwrapLabeledRuntimeValue(block))
	if err != nil {
		r.log("tool result memory feedback update error: %v", err)
		return
	}
	result.Content = updated
}

func (r *Runtime) appendMemoryStatus(msg *tape.Message, extraEntries ...tape.TapeEntry) {
	if r.memory == nil || msg == nil || msg.Role != tape.RoleToolResult {
		return
	}
	block, err := r.memory.StatusBlockWithEntries(extraEntries...)
	if err != nil {
		return
	}
	if strings.TrimSpace(block) == "" {
		return
	}
	updated, err := setRuntimeField(msg.StructuredContent, "memory_status", unwrapLabeledRuntimeValue(block))
	if err != nil {
		r.log("tool result memory status update error: %v", err)
		return
	}
	msg.StructuredContent = updated
	syncToolResultMessageContent(msg)
	r.log("%s", block)
	action, err := r.memory.ActionBlockWithEntries(extraEntries...)
	if err == nil && strings.TrimSpace(action) != "" {
		updated, updateErr := setRuntimeField(msg.StructuredContent, "memory_action", unwrapLabeledRuntimeValue(action))
		if updateErr == nil {
			msg.StructuredContent = updated
			syncToolResultMessageContent(msg)
		} else {
			r.log("tool result memory action update error: %v", updateErr)
		}
		r.log("%s", action)
	}
	warning, err := r.memory.WarningBlockWithEntries(extraEntries...)
	if err == nil && strings.TrimSpace(warning) != "" {
		updated, updateErr := setRuntimeField(msg.StructuredContent, "memory_warning", unwrapLabeledRuntimeValue(warning))
		if updateErr == nil {
			msg.StructuredContent = updated
			syncToolResultMessageContent(msg)
		} else {
			r.log("tool result memory warning update error: %v", updateErr)
		}
		r.log("%s", warning)
	}
	exposeTopology, err := r.memory.ShouldExposeTopologyWithEntries(extraEntries...)
	if err != nil || !exposeTopology {
		return
	}
	topology, err := r.memory.TopologyBlockWithEntries(extraEntries...)
	if err != nil || strings.TrimSpace(topology) == "" {
		return
	}
	updated, updateErr := setRuntimeField(msg.StructuredContent, "memory_topology", unwrapLabeledRuntimeValue(topology))
	if updateErr == nil {
		msg.StructuredContent = updated
		syncToolResultMessageContent(msg)
	} else {
		r.log("tool result memory topology update error: %v", updateErr)
	}
	r.log("%s", topology)
}

// handleMark processes a mark tool call and appends the result to the tape.
func (r *Runtime) handleMark(source tape.Message, tc tape.ToolCall) {
	if r.memory == nil {
		errMsg := runtimeToolResultMessage(tc.ID, "mark", "error", map[string]any{
			"error": "mark is not available in this runtime.",
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}
	if r.cfg == nil || !r.cfg.AnchorMarkEnabled() {
		errMsg := runtimeToolResultMessage(tc.ID, "mark", "error", map[string]any{
			"error": "[MARK ERROR] mark is disabled (QUINE_ANCHOR_MARK_ENABLED=0).",
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}
	req, err := tools.ParseMarkArgs(tc.Arguments)
	if err != nil {
		errMsg := runtimeToolResultMessage(tc.ID, "mark", "error", map[string]any{
			"error": fmt.Sprintf("[MARK ERROR] %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}
	if req.Fold && (r.cfg == nil || !r.cfg.AnchorFoldEnabled()) {
		errMsg := runtimeToolResultMessage(tc.ID, "mark", "error", map[string]any{
			"error": "[MARK ERROR] fold is disabled (QUINE_ANCHOR_FOLD_ENABLED=0); record a plain mark instead.",
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}
	result := r.memory.Mark(tc.ID, req)
	r.appendMemoryFeedback(&result)
	if !result.IsError {
		// mark() compacts live current state, so re-seed the completed tool batch
		// into current.jsonl before appending its tool result.
		r.appendContextMessage(tape.SyntheticAssistantToolBatch(source, []tape.ToolCall{tc}))
	}
	r.appendToolResult(result)
}

// handleUnfold processes an unfold tool call and appends the result to the tape.
func (r *Runtime) handleUnfold(tc tape.ToolCall) {
	if r.memory == nil {
		errMsg := runtimeToolResultMessage(tc.ID, "unfold", "error", map[string]any{
			"error": "unfold is not available in this runtime.",
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}
	req, err := tools.ParseUnfoldArgs(tc.Arguments)
	if err != nil {
		errMsg := runtimeToolResultMessage(tc.ID, "unfold", "error", map[string]any{
			"error": fmt.Sprintf("[UNFOLD ERROR] %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}
	result := r.memory.Unfold(tc.ID, req)
	r.appendMemoryFeedback(&result)
	r.appendToolResult(result)
}

// handleVision processes a vision tool call and appends the result to the tape.
// Vision does NOT consume a turn — it is a pure observation tool.
func (r *Runtime) handleVision(tc tape.ToolCall) {
	readFile := os.ReadFile
	if r.sh != nil {
		readFile = r.sh.ReadWorkspaceFile
	}
	result := tools.HandleVisionWithReader(tc.ID, tc.Arguments, readFile)

	r.appendToolResult(result)
}

// handleEscalate processes an escalate tool call, hot-swapping to a smarter model.
// Escalate does NOT consume a turn — it is a model upgrade operation.
func (r *Runtime) handleEscalate(tc tape.ToolCall) {
	reason, _ := tc.Arguments["reason"].(string)

	// Guard: already escalated or not configured
	if r.cfg.Escalated || r.cfg.SmartModelID == "" {
		errMsg := runtimeToolResultMessage(tc.ID, "escalate", "error", map[string]any{
			"error": "Escalation not available.",
		})
		r.appendRuntimeToolMessage(errMsg, true)
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
		errMsg := runtimeToolResultMessage(tc.ID, "escalate", "error", map[string]any{
			"error": fmt.Sprintf("Escalation failed: %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
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
	systemPrompt, err := r.currentSystemPrompt()
	if err == nil {
		r.tape.SetSystemPrompt(systemPrompt)
	}

	// Tool result — reason is already in tc.Arguments (the handoff note)
	resultMsg := runtimeToolResultMessage(tc.ID, "escalate", "completed", map[string]any{
		"model":  r.cfg.SmartModelID,
		"reason": reason,
	})
	r.appendRuntimeToolMessage(resultMsg, false)
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

func runtimeToolResultContent(toolName string, status string, fields map[string]any) json.RawMessage {
	payload := map[string]any{
		"tool":   toolName,
		"status": status,
	}
	for k, v := range fields {
		if v == nil {
			continue
		}
		payload[k] = v
	}
	return tape.MarshalToolResultContent(payload)
}

func runtimeToolResultMessage(toolID string, toolName string, status string, fields map[string]any) tape.Message {
	structured := runtimeToolResultContent(toolName, status, fields)
	return tape.Message{
		Role:              tape.RoleToolResult,
		Content:           "",
		StructuredContent: structured,
		ToolID:            toolID,
	}
}

func rewriteToolResultContent(content json.RawMessage, mutate func(map[string]any) error) (json.RawMessage, error) {
	return tape.RewriteStructuredToolResultContent(content, mutate)
}

func ensureRuntimeFields(payload map[string]any) map[string]any {
	runtimeFields, _ := payload["runtime"].(map[string]any)
	if runtimeFields == nil {
		runtimeFields = map[string]any{}
		payload["runtime"] = runtimeFields
	}
	return runtimeFields
}

func setRuntimeField(content json.RawMessage, key string, value any) (json.RawMessage, error) {
	return rewriteToolResultContent(content, func(payload map[string]any) error {
		ensureRuntimeFields(payload)[key] = value
		return nil
	})
}

func syncToolResultMessageContent(msg *tape.Message) {
	if msg == nil || msg.Role != tape.RoleToolResult {
		return
	}
	if tape.HasStructuredContent(msg.StructuredContent) {
		msg.Content = ""
	}
}

func unwrapLabeledRuntimeValue(block string) any {
	trimmed := strings.TrimSpace(block)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "[") {
		if idx := strings.Index(trimmed, "]"); idx >= 0 {
			trimmed = strings.TrimSpace(trimmed[idx+1:])
		}
	}

	var payload any
	if err := json.Unmarshal([]byte(trimmed), &payload); err == nil {
		return payload
	}
	return trimmed
}

func extractFSMutationsBlock(content json.RawMessage) string {
	payload, err := tape.UnmarshalStructuredToolResultContent(content)
	if err != nil {
		return ""
	}
	block, _ := payload["fs_mutations"].(string)
	if strings.TrimSpace(block) == "" {
		return ""
	}
	if strings.HasSuffix(block, "\n") {
		return block
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

// exitCodeFromResult extracts the exit code from a structured tool result.
func exitCodeFromResult(r tape.ToolResult) int {
	payload, err := tape.UnmarshalStructuredToolResultContent(r.Content)
	if err != nil {
		return 0
	}
	switch code := payload["exit_code"].(type) {
	case json.Number:
		n, err := code.Int64()
		if err == nil {
			return int(n)
		}
	case float64:
		return int(code)
	case int:
		return code
	case int64:
		return int(code)
	}
	return 0
}
