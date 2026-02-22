package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
)

// JobState represents the current state of a job.
type JobState int

const (
	JobRunning JobState = iota
	JobPaused
	JobDone
)

// Job represents a running shell command with a resource budget.
//
// The job is identified by its pgid (process group ID), which is the natural
// OS-assigned unique identifier. The Overseer goroutine monitors two budgets
// concurrently: output_limit (bytes) and timeout (duration). The first to
// trigger sends SIGSTOP to the entire process group.
//
// Output is accumulated across pause/resume cycles in a single SafeBuffer.
// The agent can read partial output at any point via ReadOutput.
//
// In interactive (PTY) mode, ptyMaster is non-nil and is the single I/O stream
// for both stdout and stderr (merged by the PTY). All output flows into stdout;
// stderr remains empty.
type Job struct {
	// ID is the process group ID — the natural job identifier.
	ID int

	cmd        *exec.Cmd
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser

	// ptyMaster is non-nil only in interactive mode. It is the PTY master fd
	// used as the sole I/O channel. Writing to it injects input; reading from
	// it produces merged stdout+stderr output.
	ptyMaster *os.File

	stdout SafeBuffer
	stderr SafeBuffer

	mu       sync.Mutex
	state    JobState
	exitCode int

	// pausedCh is closed when the job transitions to Paused.
	// A new channel is created on each Resume() call.
	pausedCh chan struct{}

	// doneCh is closed when cmd.Wait() returns (job truly finished).
	doneCh chan struct{}

	// limitExceeded is sent on (non-blocking) by pumpOutput when the combined
	// output budget is crossed. The overseer selects on this to freeze immediately
	// without waiting for the next poll tick.
	limitExceeded chan struct{}

	// waitDone carries the exit code from cmd.Wait(). It is written exactly once
	// by a goroutine started in newJob; all overseer generations read from it.
	waitDone chan int
}

// newJob spawns a shell command and returns a Job.
//
// extraFiles should contain:
//   - [0] = deliverable stdout (becomes fd 3 in child)
//   - [1] = material stdin (becomes fd 4 in child)
//
// When interactive is true, a PTY is allocated. The PTY master fd becomes the
// sole I/O channel (stdout+stderr merged); cmd.Stdin is wired to the PTY slave
// so isatty(0)==true inside the child. In non-interactive mode the existing
// pipe-based path is used unchanged.
//
// After newJob returns, the Overseer loop is started with the given budget.
// If outputLimit == 0 and timeout == 0, the job runs until natural completion.
func newJob(shell, command string, env []string, extraFiles []*os.File, timeout time.Duration, outputLimit int, interactive bool) (*Job, error) {
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if len(env) > 0 {
		cmd.Env = env
	}

	// Wire extra fds: fd 3 = deliverable stdout, fd 4 = material stdin.
	if len(extraFiles) > 0 {
		cmd.ExtraFiles = extraFiles
	}

	j := &Job{
		state:         JobRunning,
		pausedCh:      make(chan struct{}),
		doneCh:        make(chan struct{}),
		limitExceeded: make(chan struct{}, 1),
		waitDone:      make(chan int, 1),
	}

	if interactive {
		// PTY mode: allocate a pseudo-terminal and start the command.
		// pty.Start requires Setsid+Setctty; Setpgid is incompatible with Setsid
		// on macOS. We use StartWithAttrs to supply exactly these flags.
		// The process becomes its own session leader (sid==pid==pgid), so
		// Kill(-pid, SIGSTOP) still targets the whole process group correctly.
		ptmx, err := pty.StartWithAttrs(cmd, nil, &syscall.SysProcAttr{
			Setsid:  true,
			Setctty: true,
		})
		if err != nil {
			return nil, fmt.Errorf("starting PTY: %w", err)
		}
		j.ptyMaster = ptmx
		j.ID = cmd.Process.Pid // sid == pid == pgid for session leader

		// Single pump: PTY master merges stdout+stderr; all goes to j.stdout.
		go j.pumpOutput(ptmx, &j.stdout, outputLimit)
	} else {
		// Non-interactive mode: separate pipes for stdout and stderr.
		stdoutR, stdoutW, err := os.Pipe()
		if err != nil {
			return nil, fmt.Errorf("creating stdout pipe: %w", err)
		}
		stderrR, stderrW, err := os.Pipe()
		if err != nil {
			stdoutR.Close()
			stdoutW.Close()
			return nil, fmt.Errorf("creating stderr pipe: %w", err)
		}

		cmd.Stdout = stdoutW
		cmd.Stderr = stderrW

		if err := cmd.Start(); err != nil {
			stdoutR.Close()
			stdoutW.Close()
			stderrR.Close()
			stderrW.Close()
			return nil, fmt.Errorf("starting process: %w", err)
		}

		// Close write ends in the parent (child has them now).
		stdoutW.Close()
		stderrW.Close()

		j.ID = cmd.Process.Pid
		j.stdoutPipe = stdoutR
		j.stderrPipe = stderrR

		go j.pumpOutput(stdoutR, &j.stdout, outputLimit)
		go j.pumpOutput(stderrR, &j.stderr, outputLimit)
	}

	j.cmd = cmd

	// Start the single cmd.Wait() goroutine — called exactly once for the
	// lifetime of the job. All overseer generations read from j.waitDone.
	go func() {
		err := j.cmd.Wait()
		// In PTY mode, close the master fd after Wait so any blocked Read in
		// pumpOutput unblocks. (On Linux/macOS the process exit already causes
		// EIO on the master, but closing is belt-and-suspenders.)
		if j.ptyMaster != nil {
			j.ptyMaster.Close()
		}
		code := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				code = exitErr.ExitCode()
			} else {
				code = 1
			}
		}
		j.waitDone <- code
	}()

	// Start Overseer with given budget.
	go j.overseer(outputLimit, timeout)

	return j, nil
}

// pumpOutput reads from r and appends to buf until EOF or error.
// If outputLimit > 0, it signals j.limitExceeded (non-blocking) as soon as the
// combined stdout+stderr bytes cross the limit — allowing the overseer to freeze
// the job immediately without waiting for the next poll tick.
func (j *Job) pumpOutput(r io.Reader, buf *SafeBuffer, outputLimit int) {
	data := make([]byte, 4096)
	for {
		n, err := r.Read(data)
		if n > 0 {
			buf.Write(data[:n])
			if outputLimit > 0 && j.stdout.Len()+j.stderr.Len() >= outputLimit {
				// Read the current limitExceeded channel under the lock and
				// signal it non-blocking. We re-read each time so Resume()
				// can safely replace the channel between pump iterations.
				j.mu.Lock()
				ch := j.limitExceeded
				j.mu.Unlock()
				select {
				case ch <- struct{}{}:
				default:
				}
			}
		}
		if err != nil {
			return
		}
	}
}

// overseer monitors the job's output_limit and timeout budgets concurrently.
// The first budget to fire sends SIGSTOP and marks the job Paused.
// When both budgets are 0, overseer only watches for natural completion.
//
// cmd.Wait() is called exactly once (in a goroutine started by newJob); the
// result is read from j.waitDone so multiple Resume() generations are safe.
func (j *Job) overseer(outputLimit int, timeout time.Duration) {
	// Set up budget channels.
	var timeoutCh <-chan time.Time
	if timeout > 0 {
		t := time.NewTimer(timeout)
		defer t.Stop()
		timeoutCh = t.C
	}

	// output_limit is signalled by pumpOutput via limitExceeded channel.
	var limitCh <-chan struct{}
	if outputLimit > 0 {
		limitCh = j.limitExceeded
	}

	for {
		select {
		case code := <-j.waitDone:
			// Natural completion. Only one overseer generation will receive
			// this (channel has capacity 1, sent exactly once). Mark done and
			// close doneCh so any other waiting goroutines also unblock.
			j.mu.Lock()
			j.state = JobDone
			j.exitCode = code
			j.mu.Unlock()
			close(j.doneCh)
			return

		case <-j.doneCh:
			// Another overseer generation already handled completion; exit.
			return

		case <-timeoutCh:
			// Timeout budget exhausted.
			j.freeze()
			return

		case <-limitCh:
			// Output limit crossed (signalled by pumpOutput).
			j.freeze()
			return
		}
	}
}

// freeze sends SIGSTOP to the job's process group and marks it Paused.
func (j *Job) freeze() {
	_ = syscall.Kill(-j.ID, syscall.SIGSTOP)

	j.mu.Lock()
	if j.state == JobRunning {
		j.state = JobPaused
		close(j.pausedCh)
	}
	j.mu.Unlock()
}

// Wait blocks until the job transitions out of Running state
// (either Done or Paused). Returns true if paused, false if done.
func (j *Job) Wait() (paused bool) {
	j.mu.Lock()
	pausedCh := j.pausedCh
	doneCh := j.doneCh
	j.mu.Unlock()

	select {
	case <-pausedCh:
		return true
	case <-doneCh:
		return false
	}
}

// Resume sends SIGCONT and starts a new Overseer with the given budget.
// If input is non-empty and the job is interactive, input is written to the
// PTY master BEFORE SIGCONT so the data is in the kernel TTY buffer when the
// process wakes up and calls read(). Returns an error if the job is not Paused.
func (j *Job) Resume(timeout time.Duration, outputLimit int, input string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.state != JobPaused {
		return fmt.Errorf("job %d is not paused (state=%d)", j.ID, j.state)
	}

	// Reset to Running with a fresh pausedCh and limitExceeded channel.
	j.state = JobRunning
	j.pausedCh = make(chan struct{})
	j.limitExceeded = make(chan struct{}, 1)

	// Write input BEFORE SIGCONT so data is in the TTY buffer when the
	// process resumes. Only valid in interactive (PTY) mode.
	if input != "" && j.ptyMaster != nil {
		if _, err := j.ptyMaster.Write([]byte(input)); err != nil {
			// Non-fatal: log and continue — process may have exited already.
			_ = err
		}
	}

	_ = syscall.Kill(-j.ID, syscall.SIGCONT)

	// Start new Overseer.
	go j.overseer(outputLimit, timeout)

	return nil
}

// Kill sends SIGKILL to the job's process group.
func (j *Job) Kill() {
	_ = syscall.Kill(-j.ID, syscall.SIGKILL)
}

// State returns the current job state (thread-safe).
func (j *Job) State() JobState {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.state
}

// ExitCode returns the exit code. Only valid when State() == JobDone.
func (j *Job) ExitCode() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.exitCode
}

// ReadStdout returns accumulated stdout without consuming it.
func (j *Job) ReadStdout() string {
	return j.stdout.Peek()
}

// ReadStderr returns accumulated stderr without consuming it.
func (j *Job) ReadStderr() string {
	return j.stderr.Peek()
}
