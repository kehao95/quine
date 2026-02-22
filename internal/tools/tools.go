package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// ShExecutor runs shell commands as ephemeral jobs with optional resource budgets.
//
// Each sh(command) spawns a fresh process group. The job is monitored by an
// Overseer goroutine; when a budget (timeout or output_limit) is exhausted,
// the job is SIGSTOP'd and a [PAUSED] result is returned. The agent can
// resume or kill the job via the separate `job` tool.
type ShExecutor struct {
	Shell     string
	MaxOutput int
	Env       []string
	Timeout   time.Duration

	// Stdin is the material stdin file descriptor. Passed as fd 4
	// (ExtraFiles[1]) so the agent can read it via /dev/fd/4 or cat <&4.
	Stdin *os.File

	// Stdout is the deliverable stdout file descriptor. Passed as fd 3
	// (ExtraFiles[0]) so commands can write to >&3.
	Stdout *os.File

	// ProcessStarted is called when a shell starts (for SIGINT forwarding).
	ProcessStarted func(*os.Process)
	// ProcessEnded is called when a shell exits.
	ProcessEnded func()

	// jobs manages all active/paused jobs.
	jobs   *JobManager
	jobsMu sync.Mutex // guards lazy init of jobs
}

// NewShExecutor creates a ShExecutor from config with the given child
// environment. The childEnv slice should contain QUINE_* overrides (from
// Config.ChildEnv). These are merged with os.Environ() so that spawned
// commands inherit a complete environment with QUINE_* vars overlaid.
//
// QUINE_SESSION_ID is stripped from both envs so each child quine process
// generates its own unique session ID.
func NewShExecutor(cfg *config.Config, childEnv []string) *ShExecutor {
	// Filter out QUINE_SESSION_ID from childEnv
	filteredChildEnv := make([]string, 0, len(childEnv))
	for _, entry := range childEnv {
		if !strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			filteredChildEnv = append(filteredChildEnv, entry)
		}
	}

	// Filter out QUINE_SESSION_ID from os.Environ() too
	filteredOsEnv := make([]string, 0, len(os.Environ()))
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			filteredOsEnv = append(filteredOsEnv, entry)
		}
	}

	mergedEnv := MergeEnv(filteredOsEnv, filteredChildEnv)

	timeout := time.Duration(cfg.ShTimeout) * time.Second
	if timeout == 0 {
		timeout = 600 * time.Second // default 10 minutes
	}

	return &ShExecutor{
		Shell:     cfg.Shell,
		MaxOutput: cfg.OutputTruncate,
		Env:       mergedEnv,
		Timeout:   timeout,
	}
}

// MergeEnv takes the OS environment and overlays child overrides.
// Keys from childEnv take precedence over osEnv.
func MergeEnv(osEnv []string, childEnv []string) []string {
	env := make(map[string]string, len(osEnv)+len(childEnv))
	order := make([]string, 0, len(osEnv)+len(childEnv))

	for _, entry := range osEnv {
		key, _, _ := strings.Cut(entry, "=")
		if _, exists := env[key]; !exists {
			order = append(order, key)
		}
		env[key] = entry
	}

	for _, entry := range childEnv {
		key, _, _ := strings.Cut(entry, "=")
		if _, exists := env[key]; !exists {
			order = append(order, key)
		}
		env[key] = entry
	}

	result := make([]string, 0, len(order))
	for _, key := range order {
		result = append(result, env[key])
	}
	return result
}

// initJobs lazily creates the JobManager.
func (b *ShExecutor) initJobs() {
	if b.jobs == nil {
		var extraFiles []*os.File
		if b.Stdout != nil {
			extraFiles = append(extraFiles, b.Stdout)
		}
		if b.Stdin != nil {
			extraFiles = append(extraFiles, b.Stdin)
		}
		b.jobs = NewJobManager(b.Shell, b.Env, extraFiles)
	}
}

// GetJobManager returns the JobManager (creating it if needed).
// Called by runtime.go to wire the job tool dispatcher.
func (b *ShExecutor) GetJobManager() *JobManager {
	b.jobsMu.Lock()
	defer b.jobsMu.Unlock()
	b.initJobs()
	return b.jobs
}

// Execute dispatches a sh(command) call.
func (b *ShExecutor) Execute(toolID string, command string, timeout time.Duration, outputLimit int, interactive bool) tape.ToolResult {
	b.jobsMu.Lock()
	b.initJobs()
	b.jobsMu.Unlock()

	// Resolve effective timeout: explicit > executor default > hard default
	effectiveTimeout := timeout
	if effectiveTimeout == 0 {
		effectiveTimeout = b.Timeout
		if effectiveTimeout == 0 {
			effectiveTimeout = 600 * time.Second
		}
	}

	// Notify caller that a shell is starting (for SIGINT forwarding)
	// We do this synchronously before Launch so the process pointer is
	// tracked from the moment it exists.
	var proc *os.Process
	_ = proc // used via callback below

	// We need the process before the overseer might fire, so we reach
	// into a small shim: Launch returns the job which has the cmd.
	j, err := b.jobs.Launch(command, effectiveTimeout, outputLimit, interactive)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SHELL ERROR] %v", err),
			IsError: true,
		}
	}

	if b.ProcessStarted != nil {
		b.ProcessStarted(j.cmd.Process)
	}

	paused := j.Wait()

	if b.ProcessEnded != nil && !paused {
		b.ProcessEnded()
	}

	if paused {
		// Job is paused — keep in registry.
		stdoutStr := j.ReadStdout()
		stderrStr := j.ReadStderr()
		totalBytes := j.stdout.Len()
		shownBytes := len(stdoutStr)
		content := fmt.Sprintf(
			"[PAUSED] job=%d (process is STOPPED, not exited — no exit code yet)\n[STDOUT] %d bytes shown (%d bytes total in buffer)\n%s\n[STDERR]\n%s\nOptions: job(id=%d, signal=\"cont\", output_limit=N) to resume, job(id=%d, signal=\"kill\") to discard",
			j.ID, shownBytes, totalBytes, stdoutStr, stderrStr, j.ID, j.ID,
		)
		return tape.ToolResult{
			ToolID:  toolID,
			Content: content,
		}
	}

	// Natural completion — remove from registry.
	b.jobs.Remove(j.ID)

	stdoutStr := b.applyOutputLimit(j.ReadStdout())
	stderrStr := b.applyOutputLimit(j.ReadStderr())
	content := fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", j.ExitCode(), stdoutStr, stderrStr)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: j.ExitCode() != 0,
	}
}

// applyOutputLimit truncates with a visible notice (only for naturally-completed
// output that somehow exceeds MaxOutput — this should be rare with output_limit).
func (b *ShExecutor) applyOutputLimit(s string) string {
	if b.MaxOutput <= 0 || len(s) <= b.MaxOutput {
		return s
	}
	total := len(s)
	return s[:b.MaxOutput] + fmt.Sprintf("\n...[Output Truncated, %d bytes total]", total)
}

// HandleJob processes a job(id=N, ...) tool call and returns the result.
func (b *ShExecutor) HandleJob(toolID string, args map[string]any) tape.ToolResult {
	b.jobsMu.Lock()
	b.initJobs()
	b.jobsMu.Unlock()

	// Parse id
	idRaw, ok := args["id"]
	if !ok {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[JOB ERROR] missing required parameter: id",
			IsError: true,
		}
	}
	pgid := ToInt(idRaw)
	if pgid <= 0 {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[JOB ERROR] invalid id: %v", idRaw),
			IsError: true,
		}
	}

	j := b.jobs.Get(pgid)
	if j == nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[JOB ERROR] job %d not found (already completed or killed)", pgid),
			IsError: true,
		}
	}

	signal, _ := args["signal"].(string)

	switch signal {
	case "kill":
		j.Kill()
		b.jobs.Remove(pgid)
		if b.ProcessEnded != nil {
			b.ProcessEnded()
		}
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[JOB] %d killed", pgid),
		}

	case "cont":
		timeout := time.Duration(ToInt(args["timeout"])) * time.Second
		outputLimit := ToInt(args["output_limit"])
		input, _ := args["input"].(string)
		// If no budget given, use executor defaults
		if timeout == 0 {
			timeout = b.Timeout
		}
		if err := j.Resume(timeout, outputLimit, input); err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[JOB ERROR] %v", err),
				IsError: true,
			}
		}

		paused := j.Wait()
		if b.ProcessEnded != nil && !paused {
			b.ProcessEnded()
		}

		if paused {
			stdoutStr := j.ReadStdout()
			stderrStr := j.ReadStderr()
			totalBytes := j.stdout.Len()
			shownBytes := len(stdoutStr)
			content := fmt.Sprintf(
				"[PAUSED] job=%d (process is STOPPED, not exited — no exit code yet)\n[STDOUT] %d bytes shown (%d bytes total in buffer)\n%s\n[STDERR]\n%s\nOptions: job(id=%d, signal=\"cont\", output_limit=N) to resume, job(id=%d, signal=\"kill\") to discard",
				j.ID, shownBytes, totalBytes, stdoutStr, stderrStr, j.ID, j.ID,
			)
			return tape.ToolResult{ToolID: toolID, Content: content}
		}

		// Completed
		b.jobs.Remove(pgid)
		stdoutStr := b.applyOutputLimit(j.ReadStdout())
		stderrStr := b.applyOutputLimit(j.ReadStderr())
		content := fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", j.ExitCode(), stdoutStr, stderrStr)
		return tape.ToolResult{
			ToolID:  toolID,
			Content: content,
			IsError: j.ExitCode() != 0,
		}

	default:
		// Read-only: return current output without resuming
		stdoutStr := j.ReadStdout()
		stderrStr := j.ReadStderr()
		stateStr := "paused"
		if j.State() == JobRunning {
			stateStr = "running"
		}
		content := fmt.Sprintf("[JOB %d - %s]\n[STDOUT]\n%s\n[STDERR]\n%s", pgid, stateStr, stdoutStr, stderrStr)
		return tape.ToolResult{ToolID: toolID, Content: content}
	}
}

// Close kills all active jobs.
func (b *ShExecutor) Close() error {
	b.jobsMu.Lock()
	defer b.jobsMu.Unlock()
	if b.jobs != nil {
		return b.jobs.KillAll()
	}
	return nil
}

// ToInt converts various numeric types from JSON unmarshalling to int.
func ToInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	}
	return 0
}

// toDuration converts an int (seconds) to time.Duration.
func toDuration(v any) time.Duration {
	secs := ToInt(v)
	if secs <= 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}

// Keep exec.Cmd accessible for ProcessStarted callback.
var _ = (*exec.Cmd)(nil)
var _ = syscall.SIGSTOP // ensure syscall is used
