package tools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// SafeBuffer is a thread-safe byte buffer with consume semantics.
// Pump goroutines write to it; TakeAll returns and clears contents.
type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

// Write appends data to the buffer. Safe for concurrent use.
func (b *SafeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// TakeAll returns the buffer contents and resets it (consume pattern).
func (b *SafeBuffer) TakeAll() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.buf.String()
	b.buf.Reset()
	return s
}

// Len returns the current byte count in the buffer.
func (b *SafeBuffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Len()
}

// controlFdIndex is the index in ExtraFiles where the control pipe write end
// is placed. ExtraFiles indices map to child fd numbers as: fd = 3 + index.
// We place the control pipe after the deliverable stdout (index 0 = fd 3) and
// material stdin (index 1 = fd 4), so index 2 = fd 5.
//
// This preserves the existing fd 3/fd 4 convention for all shells.
const controlFdIndex = 2 // child fd 5

// controlChildFd is the fd number the control pipe gets in the child process.
const controlChildFd = 3 + controlFdIndex // fd 5

// Session represents a named, persistent shell process with isolated I/O.
//
// Each session has its own /bin/sh process, running in its own process group.
// stdout+stderr are combined and pumped into a SafeBuffer via a goroutine.
// A separate control pipe (fd 5) carries exit code signals, keeping the
// data stream uncontaminated.
type Session struct {
	Name string

	cmd   *exec.Cmd
	stdin io.WriteCloser // write commands here

	// Output capture: pipe from shell stdout+stderr -> pump -> SafeBuffer
	outputPipe io.ReadCloser
	Output     *SafeBuffer

	// Control pipe: shell writes "\x00{exitcode}\x00" to fd 5 after each command.
	// Only the read end is kept by the parent; the write end is passed to the child.
	ctrlR io.ReadCloser

	pgid int  // process group ID for cleanup
	busy bool // true while a command is executing

	// exitCode from the last completed command
	exitCode int

	// done is closed when the control pump detects command completion
	done chan struct{}

	mu sync.Mutex
}

// newSession spawns a new named shell session.
//
// extraFiles should contain:
//   - [0] = deliverable stdout (becomes fd 3 in child)
//   - [1] = material stdin (becomes fd 4 in child)
//
// The control pipe is appended at index 2 (fd 5) automatically.
func newSession(name string, shell string, env []string, extraFiles []*os.File) (*Session, error) {
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if len(env) > 0 {
		cmd.Env = env
	}

	// Create control pipe for exit code signaling.
	// The write end goes to the child as fd 5 (ExtraFiles index 2).
	// The read end stays in the parent for pumpControl().
	ctrlR, ctrlWFile, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("creating control pipe: %w", err)
	}

	// Build ExtraFiles array:
	//   index 0 (fd 3) = deliverable stdout
	//   index 1 (fd 4) = material stdin
	//   index 2 (fd 5) = control pipe write end
	//
	// If extraFiles has fewer than 2 entries, pad with nil so the control
	// pipe lands at the correct index.
	cmdExtraFiles := make([]*os.File, controlFdIndex+1)
	for i := 0; i < len(extraFiles) && i < controlFdIndex; i++ {
		cmdExtraFiles[i] = extraFiles[i]
	}
	cmdExtraFiles[controlFdIndex] = ctrlWFile
	cmd.ExtraFiles = cmdExtraFiles

	// Set up pipes for command I/O
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		ctrlR.Close()
		ctrlWFile.Close()
		return nil, fmt.Errorf("creating stdin pipe: %w", err)
	}

	// Combine stdout and stderr into a single pipe
	outputR, outputW, err := os.Pipe()
	if err != nil {
		ctrlR.Close()
		ctrlWFile.Close()
		stdinPipe.Close()
		return nil, fmt.Errorf("creating output pipe: %w", err)
	}
	cmd.Stdout = outputW
	cmd.Stderr = outputW

	// Start the shell
	if err := cmd.Start(); err != nil {
		ctrlR.Close()
		ctrlWFile.Close()
		stdinPipe.Close()
		outputR.Close()
		outputW.Close()
		return nil, fmt.Errorf("starting shell: %w", err)
	}

	// Close the write ends in the parent (child has them now)
	outputW.Close()
	ctrlWFile.Close()

	s := &Session{
		Name:       name,
		cmd:        cmd,
		stdin:      stdinPipe,
		outputPipe: outputR,
		Output:     &SafeBuffer{},
		ctrlR:      ctrlR,
		pgid:       cmd.Process.Pid,
		done:       make(chan struct{}),
	}

	// Start output pump goroutine
	go s.pumpOutput()

	return s, nil
}

// pumpOutput continuously reads from the shell's stdout+stderr pipe
// and writes to the SafeBuffer. Runs until the pipe is closed (shell exit).
func (s *Session) pumpOutput() {
	buf := make([]byte, 4096)
	for {
		n, err := s.outputPipe.Read(buf)
		if n > 0 {
			s.Output.Write(buf[:n])
		}
		if err != nil {
			return // pipe closed (shell exited or killed)
		}
	}
}

// pumpControl reads the control pipe for exit code signals.
// Format: \x00{exitcode}\x00
// When detected, marks the session as not-busy and records the exit code.
func (s *Session) pumpControl() {
	buf := make([]byte, 64)
	var acc []byte
	for {
		n, err := s.ctrlR.Read(buf)
		if n > 0 {
			acc = append(acc, buf[:n]...)
			// Look for \x00{digits}\x00 pattern
			for {
				start := bytes.IndexByte(acc, 0x00)
				if start < 0 {
					break
				}
				end := bytes.IndexByte(acc[start+1:], 0x00)
				if end < 0 {
					break
				}
				// Extract exit code
				codeStr := string(acc[start+1 : start+1+end])
				code := 0
				fmt.Sscanf(codeStr, "%d", &code)
				acc = acc[start+1+end+1:]

				s.mu.Lock()
				s.exitCode = code
				s.busy = false
				done := s.done
				s.mu.Unlock()

				// Signal completion
				select {
				case <-done:
				default:
					close(done)
				}
			}
		}
		if err != nil {
			return
		}
	}
}

// Execute runs a command in this session. Returns the exit code and output.
// Blocks until the command completes (detected via control pipe).
func (s *Session) Execute(command string, timeout time.Duration) (exitCode int, output string, err error) {
	s.mu.Lock()
	if s.busy {
		s.mu.Unlock()
		return -1, "", fmt.Errorf("session %q is busy", s.Name)
	}
	s.busy = true
	s.done = make(chan struct{})
	s.mu.Unlock()

	// Drain any accumulated output before this command
	// (so output only contains this command's results)
	s.Output.TakeAll()

	// Wrap command: run in brace group, then signal exit code on the control fd.
	// The control pipe is fd 5 in named sessions (preserving fd 3/fd 4 for
	// deliverable stdout and material stdin).
	wrappedCmd := fmt.Sprintf("{ %s\n}; printf '\\x00%%d\\x00' $? >&%d\n", command, controlChildFd)

	if _, err := io.WriteString(s.stdin, wrappedCmd); err != nil {
		s.mu.Lock()
		s.busy = false
		s.mu.Unlock()
		return -1, "", fmt.Errorf("writing command to session %q: %w", s.Name, err)
	}

	// Wait for completion signal from control pipe
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	s.mu.Lock()
	done := s.done
	s.mu.Unlock()

	select {
	case <-done:
		// Command completed
	case <-timer.C:
		s.mu.Lock()
		s.busy = false
		s.mu.Unlock()
		return -1, "", fmt.Errorf("session %q: command timed out after %v", s.Name, timeout)
	}

	// Small delay to let output pump flush
	time.Sleep(10 * time.Millisecond)

	s.mu.Lock()
	code := s.exitCode
	s.mu.Unlock()

	return code, s.Output.TakeAll(), nil
}

// ReadOutput returns accumulated output without executing a command (consume pattern).
func (s *Session) ReadOutput() string {
	return s.Output.TakeAll()
}

// IsBusy returns whether the session has a command in progress.
func (s *Session) IsBusy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.busy
}

// Close kills the session's shell process group and cleans up.
func (s *Session) Close() error {
	// Kill the entire process group
	_ = syscall.Kill(-s.pgid, syscall.SIGTERM)

	// Wait briefly for graceful exit
	done := make(chan error, 1)
	go func() {
		done <- s.cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		_ = syscall.Kill(-s.pgid, syscall.SIGKILL)
		<-done
	}

	// Close pipes
	if s.stdin != nil {
		s.stdin.Close()
	}
	if s.outputPipe != nil {
		s.outputPipe.Close()
	}
	if s.ctrlR != nil {
		s.ctrlR.Close()
	}

	return nil
}
