package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

// ForkExecutor spawns child quine processes with cloned context.
type ForkExecutor struct {
	// QuinePath is the path to the quine binary. Defaults to "./quine" or
	// the current executable if not set.
	QuinePath string

	// DataDir is the directory where tape files are stored.
	DataDir string

	// SessionID is the current session's ID (becomes PARENT_SESSION for child).
	SessionID string

	// Env contains environment variables for child processes.
	// Should include QUINE_* vars with incremented depth.
	Env []string

	// TapeWriter writes to the current session's tape file (for copying).
	TapePath string

	// DefaultTimeout is the maximum time to wait for a child (when wait=true).
	DefaultTimeout time.Duration

	// MaxOutput limits the captured output size.
	MaxOutput int

	// ProcessStarted is called when a child process starts.
	ProcessStarted func(*os.Process)

	// ProcessEnded is called when a child process ends.
	ProcessEnded func()
}

// NewForkExecutor creates a ForkExecutor from config with the given child
// environment. The childEnv slice should contain QUINE_* overrides.
func NewForkExecutor(cfg *config.Config, childEnv []string) *ForkExecutor {
	// Get the current executable path for spawning children
	quinePath, err := os.Executable()
	if err != nil {
		quinePath = "./quine"
	}

	// Build tape path
	tapePath := filepath.Join(cfg.DataDir, cfg.SessionID+".jsonl")

	return &ForkExecutor{
		QuinePath:      quinePath,
		DataDir:        cfg.DataDir,
		SessionID:      cfg.SessionID,
		Env:            MergeEnv(filterSessionID(os.Environ()), childEnv),
		TapePath:       tapePath,
		DefaultTimeout: time.Duration(cfg.ShTimeout) * time.Second,
		MaxOutput:      cfg.OutputTruncate,
	}
}

// filterSessionID removes QUINE_SESSION_ID from an environment slice.
func filterSessionID(env []string) []string {
	result := make([]string, 0, len(env))
	for _, e := range env {
		if len(e) > 17 && e[:17] == "QUINE_SESSION_ID=" {
			continue
		}
		result = append(result, e)
	}
	return result
}

// ForkRequest represents the parsed arguments from a fork tool call.
type ForkRequest struct {
	Intent string // The task for the child agent (required)
	Wait   bool   // If true, block until child completes (optional, default false)
}

// ParseForkArgs extracts ForkRequest from a ToolCall's Arguments map.
func ParseForkArgs(args map[string]any) (ForkRequest, error) {
	raw, ok := args["intent"]
	if !ok {
		return ForkRequest{}, fmt.Errorf("missing required argument: intent")
	}

	intent, ok := raw.(string)
	if !ok {
		return ForkRequest{}, fmt.Errorf("intent must be a string, got %T", raw)
	}

	if intent == "" {
		return ForkRequest{}, fmt.Errorf("intent cannot be empty")
	}

	req := ForkRequest{Intent: intent}

	if v, ok := args["wait"]; ok {
		b, ok := v.(bool)
		if !ok {
			return ForkRequest{}, fmt.Errorf("wait must be a boolean, got %T", v)
		}
		req.Wait = b
	}

	return req, nil
}

// Execute spawns a child quine process with the given intent.
// Returns a ToolResult with either:
//   - If wait=false: the child's SESSION_ID (fire-and-forget)
//   - If wait=true: the child's stdout/stderr and exit status
func (f *ForkExecutor) Execute(toolID string, req ForkRequest) tape.ToolResult {
	// Step 1: Copy the current tape to a temp file for the child
	childTapePath, err := f.copyTapeForChild()
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[FORK ERROR] Failed to copy tape: %v", err),
			IsError: true,
		}
	}

	// Step 2: Generate a session ID for the child
	// (Child will generate its own via config.Load, but we need to know it for tracking)
	// Actually, we let the child generate its own - we'll parse it from output if needed

	// Step 3: Build the command
	// Intent is passed via argv (becomes the child's mission)
	// stdin is NOT consumed - it remains available to the child via /dev/stdin
	ctx := context.Background()
	var cancel context.CancelFunc
	if req.Wait {
		ctx, cancel = context.WithTimeout(ctx, f.DefaultTimeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, f.QuinePath, req.Intent)
	cmd.Env = append(f.Env, "QUINE_CONTEXT_TAPE="+childTapePath)
	// Do NOT set cmd.Stdin - child inherits parent's stdin (data stream)

	// Set process group for cleanup
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if req.Wait {
		// Synchronous: capture output and wait
		return f.executeSync(toolID, cmd, childTapePath)
	}

	// Asynchronous: fire-and-forget
	return f.executeAsync(toolID, cmd, childTapePath)
}

// executeSync runs the child and waits for completion.
func (f *ForkExecutor) executeSync(toolID string, cmd *exec.Cmd, childTapePath string) tape.ToolResult {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Start()
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[FORK ERROR] Failed to start child: %v", err),
			IsError: true,
		}
	}

	// Notify caller that child is running
	if f.ProcessStarted != nil {
		f.ProcessStarted(cmd.Process)
	}

	// Wait for completion
	err = cmd.Wait()

	if f.ProcessEnded != nil {
		f.ProcessEnded()
	}

	// Clean up temp tape file
	os.Remove(childTapePath)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// Timeout or other error - kill process group
			if cmd.Process != nil {
				_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[FORK ERROR] Child execution failed: %v", err),
				IsError: true,
			}
		}
	}

	stdout := f.truncate(stdoutBuf.Bytes())
	stderr := f.truncate(stderrBuf.Bytes())

	content := fmt.Sprintf("[EXIT CODE] %d\n[STDOUT]\n%s\n[STDERR]\n%s", exitCode, stdout, stderr)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: exitCode != 0,
	}
}

// executeAsync starts the child and returns immediately.
func (f *ForkExecutor) executeAsync(toolID string, cmd *exec.Cmd, childTapePath string) tape.ToolResult {
	// For async, we need to capture the child's session ID.
	// The child writes its session ID to a temp file we can read.
	// Or we can parse it from the tape file name.

	// Actually, since the child generates its own session ID, we can't know it ahead of time.
	// We'll start the process and let it run independently.
	// The child's tape will be written to DataDir with its own session ID.

	// Discard stdout/stderr for async - child writes to its own tape
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	err := cmd.Start()
	if err != nil {
		os.Remove(childTapePath)
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[FORK ERROR] Failed to start child: %v", err),
			IsError: true,
		}
	}

	// Fire-and-forget: don't wait, don't track
	// The child will clean up its own resources
	go func() {
		cmd.Wait()
		os.Remove(childTapePath)
	}()

	// Return the child's PID for reference (session ID is unknown until child starts)
	content := fmt.Sprintf("[FORK] Child spawned (pid=%d, wait=false)\n"+
		"Child is running independently. Check %s for child tapes.",
		cmd.Process.Pid, f.DataDir)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: false,
	}
}

// copyTapeForChild copies the current tape file to a temp file that the child
// can read for context. Returns the path to the temp file.
func (f *ForkExecutor) copyTapeForChild() (string, error) {
	// Read current tape
	data, err := os.ReadFile(f.TapePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No tape yet - child starts fresh
			return "", nil
		}
		return "", fmt.Errorf("reading tape: %w", err)
	}

	// Create temp file in the same directory
	tmpFile, err := os.CreateTemp(f.DataDir, "fork-tape-*.jsonl")
	if err != nil {
		return "", fmt.Errorf("creating temp tape: %w", err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("writing temp tape: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("closing temp tape: %w", err)
	}

	return tmpFile.Name(), nil
}

// truncate returns the string representation of data, truncating if needed.
func (f *ForkExecutor) truncate(data []byte) string {
	if len(data) <= f.MaxOutput {
		return string(data)
	}
	total := len(data)
	truncated := string(data[:f.MaxOutput])
	return truncated + fmt.Sprintf("\n...[Output Truncated, %d bytes total]", total)
}

// ForkResultEntry returns a TapeEntry for a fork tool result.
func ForkResultEntry(sessionID string, childPID int, waited bool) tape.TapeEntry {
	data, _ := json.Marshal(map[string]any{
		"child_pid": childPID,
		"waited":    waited,
	})
	return tape.TapeEntry{Type: "fork", Data: data}
}
