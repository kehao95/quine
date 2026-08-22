package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

// handleExit processes an exit tool call. Returns (exitCode, true) if the
// process should exit, or (0, false) if the exit was rejected (e.g. failure
// without a reason) and a rejection tool result was sent back to the agent.
func (r *Runtime) handleExit(tc tape.ToolCall) (int, bool) {
	exitReq, err := tools.ParseExitArgs(tc.Arguments)
	if err != nil {
		if !r.cfg.FailOnImpossible {
			rejectMsg := runtimeToolResultMessage(tc.ID, "exit", "rejected", map[string]any{
				"error": fmt.Sprintf("Exit rejected: invalid exit args: %v", err),
			})
			r.appendRuntimeToolMessage(rejectMsg, true)
			r.log("exit rejected: invalid args: %v", err)
			return 0, false
		}
		r.log("failed to parse exit args: %v", err)
		exitReq = tools.ExitRequest{Status: tools.StatusFailure, Stderr: fmt.Sprintf("invalid exit args: %v", err)}
	}

	if !r.cfg.FailOnImpossible && exitReq.Status == tools.StatusFailure {
		rejectMsg := runtimeToolResultMessage(tc.ID, "exit", "rejected", map[string]any{
			"error": "Exit rejected: status=\"failure\" is unavailable in this runtime; only status=\"success\" is allowed.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("exit rejected: failure status unavailable when fail-on-impossible is disabled")
		return 0, false
	}

	// Validate semantic constraints - bounce back to agent on violation.
	if err := exitReq.Validate(); err != nil {
		rejectMsg := runtimeToolResultMessage(tc.ID, "exit", "rejected", map[string]any{
			"error": fmt.Sprintf("Exit rejected: %s", err),
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("exit rejected: %v", err)
		return 0, false
	}

	// Log exit.
	var argParts []string
	argParts = append(argParts, fmt.Sprintf("status=%s", exitReq.Status))
	if exitReq.Stderr != "" {
		argParts = append(argParts, fmt.Sprintf("stderr=%s", truncateOneLine(exitReq.Stderr, 80)))
	}
	r.log("exit(%s)", joinArgs(argParts))

	code := r.finalizeExitRequest(exitReq)
	totalTokens := r.tape.TokensIn + r.tape.TokensOut
	r.log("session ended (exit=%d, %d turns, %.1fs, %d tokens)",
		code, r.tape.TurnCount, time.Since(r.startTime).Seconds(), totalTokens)

	return code, true
}

// handleSh processes a sh tool call and appends the result to the tape.
// This is the ONLY tool that consumes executions.
// Returns true if the execution budget is exhausted after this call.
func (r *Runtime) handleSh(tc tape.ToolCall) bool {
	// Increment execution counter BEFORE execution (sh is the only execution-consuming tool).
	r.tape.IncrementTurn()
	// Extract command and budget parameters.
	command, _ := tc.Arguments["command"].(string)
	detach, _ := tc.Arguments["detach"].(bool)
	interactive, _ := tc.Arguments["interactive"].(bool)
	stdin, _ := tc.Arguments["stdin"].(string)
	// Reject a missing/blank/non-string command instead of silently coercing it
	// to "" and executing an exit-0 no-op the agent would read as success.
	if reason, bad := invalidShCommandReason(tc.Arguments); bad {
		rejectMsg := runtimeToolResultMessage(tc.ID, "sh", "rejected", map[string]any{
			"error": "Rejected: " + reason,
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("sh rejected: %s", reason)
		return r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns
	}
	timeout := time.Duration(0)
	if _, hasTimeout := tc.Arguments["timeout"]; hasTimeout && !r.cfg.ShTimeoutOverrideEnabled() {
		rejectMsg := runtimeToolResultMessage(tc.ID, "sh", "rejected", map[string]any{
			"error": "Rejected: sh timeout override is disabled in this runtime.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("sh rejected: timeout override disabled by config")
		return r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns
	}
	if _, hasStdin := tc.Arguments["stdin"]; hasStdin && !r.cfg.ShStdinEnabled() {
		rejectMsg := runtimeToolResultMessage(tc.ID, "sh", "rejected", map[string]any{
			"error": "Rejected: sh stdin is disabled in this runtime.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("sh rejected: stdin disabled by config")
		return r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns
	}
	if detach && !r.cfg.ShDetachEnabled() {
		rejectMsg := runtimeToolResultMessage(tc.ID, "sh", "rejected", map[string]any{
			"error": "Rejected: sh detach is disabled in this runtime.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("sh rejected: detach disabled by config")
		return r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns
	}
	if interactive && !r.cfg.ShInteractiveEnabled() {
		rejectMsg := runtimeToolResultMessage(tc.ID, "sh", "rejected", map[string]any{
			"error": "Rejected: sh interactive mode is disabled in this runtime.",
		})
		r.appendRuntimeToolMessage(rejectMsg, true)
		r.log("sh rejected: interactive mode disabled by config")
		return r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns
	}
	if r.cfg.ShTimeoutOverrideEnabled() {
		secs, present, err := tools.IntArg(tc.Arguments, "timeout")
		if err != nil {
			rejectMsg := runtimeToolResultMessage(tc.ID, "sh", "rejected", map[string]any{
				"error": "Rejected: sh timeout must be an integer number of seconds.",
			})
			r.appendRuntimeToolMessage(rejectMsg, true)
			r.log("sh rejected: invalid timeout argument: %v", err)
			return r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns
		}
		if present {
			timeout = time.Duration(secs) * time.Second
		}
	}
	// Execute.
	r.sh.TurnID = r.tape.TurnCount
	result := r.sh.Execute(tc.ID, command, timeout, 0, interactive, detach, stdin)

	if mutationBlock := extractFSMutationsBlock(result.Content); mutationBlock != "" {
		r.lastFSMutations = mutationBlock
	}
	r.syncWorldRevisionSurface()
	if r.cfg.MaxTurns > 0 {
		remaining := r.cfg.MaxTurns - r.tape.TurnCount
		updated, err := setRuntimeField(result.Content, "executions_left", remaining)
		if err != nil {
			r.log("tool result executions-left update error: %v", err)
		} else {
			result.Content = updated
		}
	}

	// Append tool result to tape.
	r.appendToolResult(result)

	// Check if execution budget is now exhausted.
	if r.cfg.MaxTurns > 0 && r.tape.TurnCount >= r.cfg.MaxTurns {
		return true
	}
	return false
}

// invalidShCommandReason reports why a sh `command` argument is unusable, if it
// is. A missing key, a non-string value, or a blank string would otherwise be
// coerced to "" and executed as an exit-0 no-op that the agent reads as success.
func invalidShCommandReason(args map[string]any) (string, bool) {
	raw, present := args["command"]
	if !present {
		return "sh requires a command", true
	}
	s, ok := raw.(string)
	if !ok {
		return "sh command must be a string", true
	}
	if strings.TrimSpace(s) == "" {
		return "sh command must not be empty", true
	}
	return "", false
}

// handleIdle processes an idle tool call and appends the resumed result to the tape.
// Idle does not consume a turn; it blocks until an external control event arrives.
func (r *Runtime) handleIdle(tc tape.ToolCall) {
	if r.cfg == nil || !r.cfg.IdleToolEnabled() {
		if spec, ok := r.toolSpec("idle"); ok {
			r.rejectDisabledTool(tc, spec)
		}
		return
	}
	r.handleIdleEnabled(tc)
}

func (r *Runtime) handleIdleEnabled(tc tape.ToolCall) {
	delivery, pendingCount := r.waitForIdleResume()
	fields := map[string]any{
		"delivery": string(delivery),
	}
	if pendingCount > 0 {
		fields["pending_count"] = pendingCount
	}
	if delivery == controlDeliveryInterrupt && pendingCount == 0 {
		fields["interrupt_notice"] = "Current operation was interrupted by peer control input."
	}

	r.log("idle resumed (delivery=%s, pending=%d)", delivery, pendingCount)
	r.appendRuntimeToolMessage(runtimeToolResultMessage(tc.ID, "idle", "completed", fields), false)
}

// handleFork processes a fork tool call and appends the result to the tape.
func (r *Runtime) handleFork(tc tape.ToolCall) {
	if r.cfg == nil || !r.cfg.ForkEnabled() {
		if spec, ok := r.toolSpec("fork"); ok {
			r.rejectDisabledTool(tc, spec)
		}
		return
	}
	r.handleForkEnabled(tc)
}

func (r *Runtime) handleForkEnabled(tc tape.ToolCall) {
	forkReq, err := tools.ParseForkArgs(tc.Arguments, r.cfg)
	if err != nil {
		r.log("fork parse error: %v", err)
		errMsg := runtimeToolResultMessage(tc.ID, "fork", "error", map[string]any{
			"error": fmt.Sprintf("[FORK ERROR] %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}

	if errMsg := r.precheckProcessCreation(tc.ID, "fork", len(forkReq.Children)); errMsg != nil {
		r.appendRuntimeToolMessage(*errMsg, true)
		return
	}

	r.refreshForkEnv()
	r.fork.Mission = r.originalInput
	r.fork.TurnID = r.tape.TurnCount
	result := r.fork.Execute(tc.ID, forkReq)

	if result.IsError {
		r.log("fork done [error]: %s", truncateStr(string(result.Content), 100))
	} else {
		r.log("fork done (n=%d, mode=%s)", len(forkReq.Children), forkReq.Mode)
	}

	if mutationBlock := extractFSMutationsBlock(result.Content); mutationBlock != "" {
		r.lastFSMutations = mutationBlock
	}
	r.syncWorldRevisionSurface()

	r.appendToolResult(result)
}

// handleSpawn processes a spawn tool call and appends the result to the tape.
func (r *Runtime) handleSpawn(tc tape.ToolCall) {
	if r.cfg == nil || !r.cfg.SpawnEnabled() {
		if spec, ok := r.toolSpec("spawn"); ok {
			r.rejectDisabledTool(tc, spec)
		}
		return
	}
	r.handleSpawnEnabled(tc)
}

func (r *Runtime) handleSpawnEnabled(tc tape.ToolCall) {
	spawnReq, err := tools.ParseSpawnArgs(tc.Arguments, r.cfg)
	if err != nil {
		r.log("spawn parse error: %v", err)
		errMsg := runtimeToolResultMessage(tc.ID, "spawn", "error", map[string]any{
			"error": fmt.Sprintf("[SPAWN ERROR] %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}

	if errMsg := r.precheckProcessCreation(tc.ID, "spawn", len(spawnReq.Children)); errMsg != nil {
		r.appendRuntimeToolMessage(*errMsg, true)
		return
	}

	if r.spawn == nil {
		r.spawn = tools.NewSpawnExecutor(r.cfg)
		r.spawn.ProcessStarted = func(proc *os.Process) {
			r.activeProcess.Store(proc)
		}
		r.spawn.ProcessEnded = func() {
			r.activeProcess.Store(nil)
		}
	}
	r.refreshSpawnEnv()
	result := r.spawn.Execute(tc.ID, spawnReq)
	if result.IsError {
		r.log("spawn done [error]: %s", truncateStr(string(result.Content), 100))
	} else {
		r.log("spawn done (n=%d, mode=%s)", len(spawnReq.Children), spawnReq.Mode)
	}
	r.appendToolResult(result)
}

// forkChildEnv rebuilds the fork/spawn boundary environment from the CURRENT
// process environ, the CURRENT config, and the override as it stands on disk
// right now. Called immediately before every fork and spawn: the agent may have
// rewritten config/env/override in the shell command before this one, and a
// child born from a cached policy would be born from a policy the agent has
// already replaced.
func (r *Runtime) forkChildEnv() []string {
	return tools.ForkChildEnv(r.cfg, func(err error) {
		r.log("child-env override ignored for this child: %v", err)
	})
}

func (r *Runtime) refreshForkEnv() {
	if r.fork == nil {
		return
	}
	r.fork.Env = r.forkChildEnv()
}

func (r *Runtime) refreshSpawnEnv() {
	if r.spawn == nil {
		return
	}
	r.spawn.Env = r.forkChildEnv()
}

// handleExec processes an exec tool call.
// Note: On success, this function does NOT return - the process is replaced.
// On failure, it appends an error result to the tape.
func (r *Runtime) handleExec(tc tape.ToolCall) {
	execReq, err := tools.ParseExecArgs(tc.Arguments)
	if err != nil {
		r.log("exec parse error: %v", err)
		errMsg := runtimeToolResultMessage(tc.ID, "exec", "error", map[string]any{
			"error": fmt.Sprintf("[EXEC ERROR] %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}

	// Log the call: show target/argv.
	{
		var parts []string
		if execReq.Target != "" {
			parts = append(parts, fmt.Sprintf("target=%s", execReq.Target))
		}
		if execReq.Argv != nil {
			parts = append(parts, fmt.Sprintf("argv=%q", execReq.Argv))
		}
		r.log("exec(%s)", joinArgs(parts))
	}

	registryReleased := false
	if r.agentRegistry != nil {
		r.stopPeerDiscoveryHeartbeat()
		if err := r.agentRegistry.UnpublishSelfPID(); err != nil {
			r.log("exec pre-unpublish error: %v", err)
		}
		if err := r.agentRegistry.Deregister(); err != nil {
			r.log("exec deregistration error: %v", err)
			errMsg := runtimeToolResultMessage(tc.ID, "exec", "error", map[string]any{
				"error": fmt.Sprintf("[EXEC ERROR] failed to release active agent registration before exec: %v", err),
			})
			r.appendRuntimeToolMessage(errMsg, true)
			return
		}
		registryReleased = true
	}

	// Write outcome before exec (we're about to be replaced).
	r.flushPendingToolResult()
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        0,
		Stderr:          "exec: process image replaced",
		DurationMs:      duration.Milliseconds(),
		TerminationMode: tape.TermExec,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())

	if r.tapeWriter != nil {
		r.tapeWriter.Close()
	}
	if r.logFile != nil {
		r.logFile.Close()
	}

	// Execute the exec; this does not return on success.
	result := r.exec.Execute(tc.ID, execReq)

	// If we get here, exec failed. Reopen the log file to record the error.
	logPath := r.cfg.SessionLogPath("")
	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)
	r.logFile, _ = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if registryReleased && r.agentRegistry != nil {
		if err := r.agentRegistry.Register(); err != nil {
			r.log("agent re-registration after exec failure failed: %v", err)
			r.appendRuntimeToolMessage(runtimeToolResultMessage(tc.ID, "exec", "error", map[string]any{
				"error": fmt.Sprintf("[EXEC ERROR] failed to re-register after exec: %v", err),
			}), true)
			return
		}
		if err := r.agentRegistry.PublishSelfPID(); err != nil {
			r.log("agent pid republish after exec failure failed: %v", err)
			r.appendRuntimeToolMessage(runtimeToolResultMessage(tc.ID, "exec", "error", map[string]any{
				"error": fmt.Sprintf("[EXEC ERROR] failed to republish pid route after exec: %v", err),
			}), true)
			return
		}
		if r.cfg != nil && r.cfg.PeerDiscoveryEnabled {
			r.startPeerDiscoveryHeartbeat()
		}
	}

	r.log("exec failed: %s", truncateStr(string(result.Content), 100))

	tw, err := tape.NewWriter(r.cfg.TapeDir(""), r.cfg.TapeID)
	if err == nil {
		r.tapeWriter = tw
	}

	r.appendToolResult(result)
}

func (r *Runtime) handleSwitchWorld(tc tape.ToolCall) {
	req, err := tools.ParseSwitchWorldArgs(tc.Arguments)
	if err != nil {
		r.log("switch_world parse error: %v", err)
		errMsg := runtimeToolResultMessage(tc.ID, "switch_world", "error", map[string]any{
			"error": fmt.Sprintf("[SWITCH WORLD ERROR] %v", err),
		})
		r.appendRuntimeToolMessage(errMsg, true)
		return
	}

	result := r.sh.SwitchWorld(tc.ID, req.Target)
	if mutationBlock := extractFSMutationsBlock(result.Content); mutationBlock != "" {
		r.lastFSMutations = mutationBlock
	}
	r.syncWorldRevisionSurface()
	r.appendToolResult(result)
}

func (r *Runtime) syncWorldRevisionSurface() {
	if r.cfg == nil || !r.cfg.WorkspaceEnabled || r.sh == nil {
		return
	}
	revision := r.sh.CurrentWorldRevision()
	if strings.TrimSpace(revision) == "" {
		return
	}
	revisionChanged := revision != r.cfg.WorkspaceCurrentRevision
	r.cfg.WorkspaceCurrentRevision = revision
	if revisionChanged {
		// The fork executor used to be rebuilt here, solely to refresh a cached
		// child env that carried QUINE_WORKSPACE_CURRENT_REVISION. Fork children
		// no longer receive that name — a revision handle names a place in THIS
		// process's workspace state, and a child mounts its own view — and the
		// fork env is rebuilt before every call regardless. Nothing to refresh.
		if r.tape != nil {
			systemPrompt, err := r.currentSystemPrompt()
			if err == nil {
				r.tape.SetSystemPrompt(systemPrompt)
			}
		}
	}
	if !r.agentRootBootstrapped {
		return
	}
	if err := r.syncSessionStatusAndWorldSurfaces(); err != nil {
		r.log("agent status surface sync error: %v", err)
	}
}
