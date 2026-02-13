package tools

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
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

// shellInit defines helper shell functions that are run once when the
// persistent shell starts.
const shellInit = `
write_file() {
    local path="$1"
    shift
    mkdir -p "$(dirname "$path")"
    printf '%s\n' "$*" > "$path"
}

read_file() {
    cat -n "$1"
}

`

// ShExecutor runs shell commands via a persistent /bin/sh process.
// Commands are written to the shell's stdin pipe, and output is read
// until a sentinel marker is detected.
type ShExecutor struct {
	Shell     string
	MaxOutput int
	ShellInit string   // Shell initialization script (helper functions)
	Env       []string // Base environment variables (without QUINE_SESSION_ID)

	// Stdin is the material stdin file descriptor. With the persistent shell,
	// this is passed as fd 4 (ExtraFiles[1]) so the agent can read it via
	// /dev/fd/4 or cat <&4.
	Stdin *os.File

	// Stdout is the deliverable stdout file descriptor. Passed as fd 3
	// (ExtraFiles[0]) so commands can write to >&3.
	Stdout *os.File

	// ProcessStarted is called when the persistent shell starts.
	ProcessStarted func(*os.Process)
	// ProcessEnded is called when the persistent shell exits.
	ProcessEnded func()

	// Persistent shell process
	cmd        *exec.Cmd
	stdinPipe  io.WriteCloser // Go writes commands here
	stdoutPipe io.ReadCloser  // Go reads output+sentinel here
	mu         sync.Mutex     // Serializes Execute() calls
	started    bool
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

	return &ShExecutor{
		Shell:     cfg.Shell,
		MaxOutput: cfg.OutputTruncate,
		ShellInit: shellInit,
		Env:       MergeEnv(filteredOsEnv, filteredChildEnv),
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

// generateNonce returns a random 16-character hex string for sentinel markers.
func generateNonce() string {
	var b [8]byte
	rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// Start spawns the persistent /bin/sh process. It is safe to call Start
// multiple times; subsequent calls are no-ops if the shell is already running.
func (b *ShExecutor) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.startLocked()
}

// startLocked spawns the persistent shell. Caller must hold b.mu.
func (b *ShExecutor) startLocked() error {
	if b.started {
		return nil
	}

	shell := b.Shell
	if shell == "" {
		shell = "/bin/sh"
	}

	b.cmd = exec.Command(shell)

	// Set environment
	if len(b.Env) > 0 {
		b.cmd.Env = b.Env
	}

	// Process group for signal forwarding
	b.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Set up extra file descriptors:
	// fd 3 = b.Stdout (deliverable stdout)
	// fd 4 = b.Stdin  (material stdin)
	var extraFiles []*os.File
	if b.Stdout != nil {
		extraFiles = append(extraFiles, b.Stdout)
	}
	if b.Stdin != nil {
		extraFiles = append(extraFiles, b.Stdin)
	}
	if len(extraFiles) > 0 {
		b.cmd.ExtraFiles = extraFiles
	}

	// Shell stderr → discard (per-command stderr goes to temp files)
	b.cmd.Stderr = io.Discard

	// Set up pipes for command I/O
	var err error
	b.stdinPipe, err = b.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("creating stdin pipe: %w", err)
	}

	b.stdoutPipe, err = b.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	// Start the shell process
	if err := b.cmd.Start(); err != nil {
		return fmt.Errorf("starting shell: %w", err)
	}

	b.started = true

	// Notify caller that the persistent shell is running
	if b.ProcessStarted != nil {
		b.ProcessStarted(b.cmd.Process)
	}

	// Run shellInit to define helper functions, consuming its output.
	// We write the init script followed by a sentinel echo. The init
	// script defines functions (which produce no output), so we just
	// need to consume until the sentinel appears.
	if b.ShellInit != "" {
		initNonce := generateNonce()
		sentinel := fmt.Sprintf("___QUINE_DONE_%s", initNonce)
		// Write init commands directly (not wrapped in { }) since they
		// contain function definitions with their own braces.
		initCmd := b.ShellInit + fmt.Sprintf("\necho \"%s_0___\"\n", sentinel)
		if _, err := io.WriteString(b.stdinPipe, initCmd); err != nil {
			b.closeLocked()
			return fmt.Errorf("writing shell init: %w", err)
		}

		// Consume init output until sentinel
		scanner := bufio.NewScanner(b.stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, sentinel+"_") && strings.HasSuffix(line, "___") {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			b.closeLocked()
			return fmt.Errorf("reading shell init output: %w", err)
		}
	}

	return nil
}

// Close shuts down the persistent shell process gracefully.
func (b *ShExecutor) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closeLocked()
}

// closeLocked shuts down the persistent shell. Caller must hold b.mu.
func (b *ShExecutor) closeLocked() error {
	if !b.started {
		return nil
	}

	b.started = false

	// Close stdin pipe — shell will exit when it reads EOF
	if b.stdinPipe != nil {
		b.stdinPipe.Close()
	}

	// Wait for process to finish with a short deadline
	done := make(chan error, 1)
	go func() {
		done <- b.cmd.Wait()
	}()

	select {
	case <-done:
		// Process exited cleanly
	case <-time.After(2 * time.Second):
		// Force kill the process group
		if b.cmd.Process != nil {
			_ = syscall.Kill(-b.cmd.Process.Pid, syscall.SIGKILL)
		}
		<-done // Wait for Wait() to return after kill
	}

	// Close stdout pipe
	if b.stdoutPipe != nil {
		b.stdoutPipe.Close()
	}

	// Notify caller that the persistent shell has exited
	if b.ProcessEnded != nil {
		b.ProcessEnded()
	}

	return nil
}

// handleCrash cleans up after the persistent shell has crashed.
// Caller must hold b.mu.
func (b *ShExecutor) handleCrash() {
	b.started = false

	if b.stdinPipe != nil {
		b.stdinPipe.Close()
	}
	if b.stdoutPipe != nil {
		b.stdoutPipe.Close()
	}
	if b.cmd != nil && b.cmd.Process != nil {
		_ = syscall.Kill(-b.cmd.Process.Pid, syscall.SIGKILL)
		b.cmd.Wait()
	}

	// Notify caller that the shell died
	if b.ProcessEnded != nil {
		b.ProcessEnded()
	}
}

// Execute runs a shell command in the persistent shell and returns a ToolResult.
//
// The persistent shell has these file descriptors:
//   - fd 0 (stdin): pipe from Go (for receiving commands)
//   - fd 1 (stdout): pipe to Go (for sending output + sentinel)
//   - fd 3: deliverable stdout (b.Stdout, via ExtraFiles[0])
//   - fd 4: material stdin (b.Stdin, via ExtraFiles[1])
func (b *ShExecutor) Execute(toolID string, command string) tape.ToolResult {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Auto-start if not started (or if previously crashed)
	if !b.started {
		if err := b.startLocked(); err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[SHELL ERROR] %v", err),
				IsError: true,
			}
		}
	}

	// Generate unique nonce for this command
	nonce := generateNonce()

	// Build wrapped command:
	// - Run the user's command in a brace group { ...; } so it executes in
	//   the current shell context (cd, export, variables all persist)
	// - Redirect stderr to a nonce-named temp file for clean capture
	// - Echo a sentinel line with the exit code
	//
	// Risk: if the user command calls `exit N`, it kills the persistent shell.
	// This is handled by crash recovery (handleCrash → auto-restart on next call).
	// The system prompt instructs the agent not to use bare `exit` in sh commands.
	stderrFile := fmt.Sprintf("/tmp/__quine_stderr_%s", nonce)
	sentinel := fmt.Sprintf("___QUINE_DONE_%s", nonce)
	wrappedCmd := fmt.Sprintf(
		"{ %s\n} 2>\"%s\"; echo \"%s_${?}___\"\n",
		command, stderrFile, sentinel,
	)

	// Write command to shell stdin
	if _, err := io.WriteString(b.stdinPipe, wrappedCmd); err != nil {
		// Shell probably died
		b.handleCrash()
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[SHELL ERROR] Shell process died. State lost. Retrying on next call.",
			IsError: true,
		}
	}

	// Read stdout until sentinel
	var stdout strings.Builder
	scanner := bufio.NewScanner(b.stdoutPipe)
	// Increase scanner buffer for large outputs (1MB max per line)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	exitCode := 0
	foundSentinel := false

	for scanner.Scan() {
		line := scanner.Text()
		// Check if this line IS the sentinel: ___QUINE_DONE_{nonce}_{exitcode}___
		if strings.HasPrefix(line, sentinel+"_") && strings.HasSuffix(line, "___") {
			// Parse exit code
			codeStr := line[len(sentinel)+1 : len(line)-3]
			fmt.Sscanf(codeStr, "%d", &exitCode)
			foundSentinel = true
			break
		}
		stdout.WriteString(line)
		stdout.WriteString("\n")
	}

	if !foundSentinel {
		// EOF without sentinel — shell crashed
		b.handleCrash()
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[SHELL ERROR] Shell process terminated unexpectedly. State lost.",
			IsError: true,
		}
	}

	// Read stderr from temp file
	stderrBytes, _ := os.ReadFile(stderrFile)
	os.Remove(stderrFile) // Clean up

	// Truncate and format output
	stdoutStr := b.truncate([]byte(stdout.String()))
	stderrStr := b.truncate(stderrBytes)
	content := fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", exitCode, stdoutStr, stderrStr)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: exitCode != 0,
	}
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
