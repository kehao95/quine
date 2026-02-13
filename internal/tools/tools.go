package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// shellInit defines helper shell functions that are prepended to every command.
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

// ShExecutor runs shell commands and returns structured results.
type ShExecutor struct {
	Shell          string
	DefaultTimeout time.Duration
	MaxOutput      int
	ShellInit      string   // Shell initialization script (helper functions)
	Env            []string // Base environment variables (without QUINE_SESSION_ID)

	// Stdin is the process's real stdin file descriptor. Commands executed
	// via sh can read from this to access the data stream (Material channel).
	Stdin *os.File

	// Stdout is the process's real stdout file descriptor. When a command
	// is executed with passthrough=true, the child's stdout is wired directly
	// to this file instead of being captured — enabling binary output.
	Stdout *os.File

	// ProcessStarted is called when a child process starts. The caller can
	// use this to track the active process (e.g., for SIGINT forwarding).
	// ProcessEnded is called when the process exits (or fails to start).
	ProcessStarted func(*os.Process)
	ProcessEnded   func()
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
		Shell:          cfg.Shell,
		DefaultTimeout: time.Duration(cfg.ShTimeout) * time.Second,
		MaxOutput:      cfg.OutputTruncate,
		ShellInit:      shellInit,
		Env:            MergeEnv(filteredOsEnv, filteredChildEnv),
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

// Execute runs a shell command and returns a ToolResult.
//
// When passthrough is true and the command succeeds (exit code 0), stdout is
// flushed to the process's real stdout (b.Stdout). This enables binary output
// that flows from the child command to the parent's fd 1 without string
// conversion, truncation, or context pollution.
//
// IMPORTANT: stdout is always captured first into a buffer. Only on success
// is it flushed to real stdout. This prevents partial output on failure —
// e.g., if a shell command outputs some content before encountering a syntax
// error, that partial output won't pollute the real stdout.
//
// QUINE_SESSION_ID is intentionally NOT set in the child environment.
// Each child ./quine process generates its own unique session ID via
// config.Load(). This is necessary because a single sh command can
// spawn multiple ./quine children (e.g. via & backgrounding), and each
// must write to its own tape file.
func (b *ShExecutor) Execute(toolID string, command string, timeout int, passthrough bool) tape.ToolResult {
	// Determine effective timeout: use the smaller of the provided timeout
	// and DefaultTimeout. If timeout arg is 0, use DefaultTimeout.
	effectiveTimeout := b.DefaultTimeout
	if timeout > 0 {
		provided := time.Duration(timeout) * time.Second
		if provided < effectiveTimeout {
			effectiveTimeout = provided
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), effectiveTimeout)
	defer cancel()

	// Prepend shell init functions to the command.
	fullCommand := b.ShellInit + command

	cmd := exec.CommandContext(ctx, b.Shell, "-c", fullCommand)

	// Set environment for child processes. This includes the full OS
	// environment merged with QUINE_* overrides so that:
	// - Regular commands (ls, cat, etc.) work because PATH etc. are present
	// - Recursive ./quine invocations get incremented DEPTH, PARENT_SESSION, etc.
	// - Non-quine commands simply ignore the QUINE_* vars
	// Note: QUINE_SESSION_ID is absent — each ./quine generates its own.
	if len(b.Env) > 0 {
		cmd.Env = b.Env
	}

	// Ensure child processes are killed on timeout by using a process group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Connect the child's stdin to the process's real stdin (Material channel).
	// This allows commands like `cat` or `./quine` to read the data stream.
	if b.Stdin != nil {
		cmd.Stdin = b.Stdin
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	// Always capture stdout first, even in passthrough mode.
	// We only flush to real stdout on success to prevent partial output on failure.
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Start()
	if err == nil {
		// Notify caller that a child process is running.
		if b.ProcessStarted != nil {
			b.ProcessStarted(cmd.Process)
		}
		err = cmd.Wait()
	}

	// Notify caller that the child process has exited.
	if b.ProcessEnded != nil {
		b.ProcessEnded()
	}

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// Context deadline exceeded or other error — send SIGKILL to
			// the entire process group to clean up orphans.
			if cmd.Process != nil {
				_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
			exitCode = -1
		}
	}

	stderr := b.truncate(stderrBuf.Bytes())

	var content string
	if passthrough && b.Stdout != nil {
		// Passthrough mode: only flush to real stdout if command succeeded.
		// This prevents partial output on failure (e.g., shell quoting errors
		// that output some content before failing).
		if exitCode == 0 {
			b.Stdout.Write(stdoutBuf.Bytes())
			content = fmt.Sprintf("[EXIT CODE] %d\n[STDOUT] (passthrough)\n[STDERR]\n%s", exitCode, stderr)
		} else {
			// Command failed — report the captured output in the tool result
			// instead of sending to real stdout. This allows the agent to see
			// what went wrong and retry with a different approach.
			stdout := b.truncate(stdoutBuf.Bytes())
			content = fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", exitCode, stdout, stderr)
		}
	} else {
		stdout := b.truncate(stdoutBuf.Bytes())
		content = fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", exitCode, stdout, stderr)
	}

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
