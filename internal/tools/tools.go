package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// ShExecutor runs shell commands via either ephemeral (anonymous) or persistent
// (named session) shells.
//
// Anonymous mode: spawns a fresh /bin/sh, runs the command, collects output via
// pipe EOF (no sentinels), kills the process group, returns.
//
// Named session mode: delegates to a SessionManager which maintains persistent
// shells with control-pipe-based completion detection.
type ShExecutor struct {
	Shell     string
	MaxOutput int
	Env       []string // Base environment variables (without QUINE_SESSION_ID)
	Timeout   time.Duration

	// Stdin is the material stdin file descriptor. Passed as fd 4
	// (ExtraFiles[1]) so the agent can read it via /dev/fd/4 or cat <&4.
	Stdin *os.File

	// Stdout is the deliverable stdout file descriptor. Passed as fd 3
	// (ExtraFiles[0]) so commands can write to >&3.
	Stdout *os.File

	// ProcessStarted is called when an anonymous shell starts.
	ProcessStarted func(*os.Process)
	// ProcessEnded is called when an anonymous shell exits.
	ProcessEnded func()

	// sessions manages named persistent shells.
	sessions *SessionManager

	mu sync.Mutex // Serializes anonymous Execute() calls
}

// NewShExecutor creates a ShExecutor from config with the given child
// environment. The childEnv slice should contain QUINE_* overrides (from
// Config.ChildEnv). These are merged with os.Environ() so that spawned
// commands inherit a complete environment with QUINE_* vars overlaid.
//
// Note: QUINE_SESSION_ID is stripped from BOTH childEnv and os.Environ() so
// that each child ./quine process generates its own unique session ID via
// config.Load(). This is critical because a single sh command can spawn
// multiple ./quine children (e.g. via backgrounding with &), and they must
// each have distinct session IDs to write to separate tape files.
func NewShExecutor(cfg *config.Config, childEnv []string) *ShExecutor {
	// Filter out QUINE_SESSION_ID from childEnv
	filteredChildEnv := make([]string, 0, len(childEnv))
	for _, entry := range childEnv {
		if !strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			filteredChildEnv = append(filteredChildEnv, entry)
		}
	}

	// Filter out QUINE_SESSION_ID from os.Environ() too — the parent's
	// session ID must not leak into children.
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
// Keys from childEnv take precedence over osEnv. This ensures spawned
// processes have a full environment (PATH, HOME, etc.) with QUINE_*
// variables set for recursive invocations.
func MergeEnv(osEnv []string, childEnv []string) []string {
	env := make(map[string]string, len(osEnv)+len(childEnv))
	order := make([]string, 0, len(osEnv)+len(childEnv))

	// Load OS environment first
	for _, entry := range osEnv {
		key, _, _ := strings.Cut(entry, "=")
		if _, exists := env[key]; !exists {
			order = append(order, key)
		}
		env[key] = entry
	}

	// Overlay child env vars (QUINE_* take precedence)
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

// initSessions lazily creates the SessionManager on first named-session use.
func (b *ShExecutor) initSessions() {
	if b.sessions == nil {
		var extraFiles []*os.File
		if b.Stdout != nil {
			extraFiles = append(extraFiles, b.Stdout)
		}
		if b.Stdin != nil {
			extraFiles = append(extraFiles, b.Stdin)
		}
		b.sessions = NewSessionManager(b.Shell, b.Env, extraFiles)
	}
}

// Execute dispatches a shell command based on the presence of a session name.
//
// Three modes:
//   - command only → anonymous ephemeral execution
//   - command + session → execute in named session (error if busy)
//   - session only → read accumulated output from named session
func (b *ShExecutor) Execute(toolID string, command string, session string) tape.ToolResult {
	if session == "" {
		// Anonymous ephemeral mode
		return b.executeAnonymous(toolID, command)
	}

	if command == "" {
		// Read-only: return accumulated output from named session
		return b.readSession(toolID, session)
	}

	// Named session: execute command in persistent shell
	return b.executeSession(toolID, command, session)
}

// executeAnonymous spawns a fresh shell, runs the command, collects output
// via pipe EOF (stdout+stderr combined), and kills the process group.
// No sentinels, no temp files, no persistent state.
func (b *ShExecutor) executeAnonymous(toolID string, command string) tape.ToolResult {
	b.mu.Lock()
	defer b.mu.Unlock()

	shell := b.Shell
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if len(b.Env) > 0 {
		cmd.Env = b.Env
	}

	// Set up extra file descriptors (same as old persistent shell):
	// fd 3 = deliverable stdout, fd 4 = material stdin
	var extraFiles []*os.File
	if b.Stdout != nil {
		extraFiles = append(extraFiles, b.Stdout)
	}
	if b.Stdin != nil {
		extraFiles = append(extraFiles, b.Stdin)
	}
	if len(extraFiles) > 0 {
		cmd.ExtraFiles = extraFiles
	}

	// Capture stdout and stderr separately via pipes
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SHELL ERROR] creating stdout pipe: %v", err),
			IsError: true,
		}
	}
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		stdoutR.Close()
		stdoutW.Close()
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SHELL ERROR] creating stderr pipe: %v", err),
			IsError: true,
		}
	}

	cmd.Stdout = stdoutW
	cmd.Stderr = stderrW

	// Start the shell
	if err := cmd.Start(); err != nil {
		stdoutR.Close()
		stdoutW.Close()
		stderrR.Close()
		stderrW.Close()
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SHELL ERROR] %v", err),
			IsError: true,
		}
	}

	// Close write ends in parent (child has them)
	stdoutW.Close()
	stderrW.Close()

	// Notify caller that an anonymous shell is running
	if b.ProcessStarted != nil {
		b.ProcessStarted(cmd.Process)
	}

	// Read stdout and stderr concurrently
	var stdoutBytes, stderrBytes []byte
	var readWg sync.WaitGroup
	readWg.Add(2)
	go func() {
		defer readWg.Done()
		stdoutBytes, _ = io.ReadAll(stdoutR)
		stdoutR.Close()
	}()
	go func() {
		defer readWg.Done()
		stderrBytes, _ = io.ReadAll(stderrR)
		stderrR.Close()
	}()

	// Wait for command with timeout
	doneCh := make(chan error, 1)
	go func() {
		doneCh <- cmd.Wait()
	}()

	timeout := b.Timeout
	if timeout == 0 {
		timeout = 600 * time.Second
	}

	var exitCode int
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case err := <-doneCh:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 1
			}
		}
	case <-timer.C:
		// Timeout: kill process group
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-doneCh // wait for Wait() to return
		readWg.Wait()

		if b.ProcessEnded != nil {
			b.ProcessEnded()
		}

		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SHELL ERROR] command timed out after %v", timeout),
			IsError: true,
		}
	}

	// Wait for read goroutines to finish
	readWg.Wait()

	// Notify caller that the anonymous shell has exited
	if b.ProcessEnded != nil {
		b.ProcessEnded()
	}

	// Format output
	stdoutStr := b.truncate(stdoutBytes)
	stderrStr := b.truncate(stderrBytes)
	content := fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", exitCode, stdoutStr, stderrStr)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: exitCode != 0,
	}
}

// executeSession runs a command in a named persistent session.
func (b *ShExecutor) executeSession(toolID string, command string, sessionName string) tape.ToolResult {
	b.mu.Lock()
	b.initSessions()
	b.mu.Unlock()

	sess, err := b.sessions.GetOrCreate(sessionName)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SESSION ERROR] creating session %q: %v", sessionName, err),
			IsError: true,
		}
	}

	timeout := b.Timeout
	if timeout == 0 {
		timeout = 600 * time.Second
	}

	exitCode, output, err := sess.Execute(command, timeout)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SESSION ERROR] %v", err),
			IsError: true,
		}
	}

	// Format output (named sessions combine stdout+stderr)
	outputStr := b.truncate([]byte(output))
	content := fmt.Sprintf("[EXIT CODE] %d\n[OUTPUT]\n%s", exitCode, outputStr)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: exitCode != 0,
	}
}

// readSession returns accumulated output from a named session without
// executing a command. Always succeeds (returns empty if no output).
func (b *ShExecutor) readSession(toolID string, sessionName string) tape.ToolResult {
	b.mu.Lock()
	b.initSessions()
	b.mu.Unlock()

	sess := b.sessions.Get(sessionName)
	if sess == nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SESSION ERROR] session %q not found", sessionName),
			IsError: true,
		}
	}

	output := sess.ReadOutput()
	return tape.ToolResult{
		ToolID:  toolID,
		Content: fmt.Sprintf("[OUTPUT]\n%s", output),
	}
}

// Close shuts down all sessions.
func (b *ShExecutor) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.sessions != nil {
		return b.sessions.CloseAll()
	}
	return nil
}

// truncate returns the string representation of data, truncating it if it
// exceeds MaxOutput bytes with a trailing notice.
func (b *ShExecutor) truncate(data []byte) string {
	if len(data) <= b.MaxOutput {
		return string(data)
	}
	total := len(data)
	truncated := string(data[:b.MaxOutput])
	return truncated + fmt.Sprintf("\n...[Output Truncated, %d bytes total]", total)
}
