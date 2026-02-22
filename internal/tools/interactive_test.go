package tools

import (
	"strings"
	"testing"
	"time"
)

// --- Interactive (PTY) mode tests ---

// TestInteractivePTYBasicOutput verifies that interactive=true allocates a PTY,
// captures output, and returns exit code 0 on natural completion.
func TestInteractivePTYBasicOutput(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	result := b.Execute("pty-basic", "echo hello_pty", 5*time.Second, 0, true)

	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0") {
		t.Errorf("expected exit code 0, got:\n%s", result.Content)
	}
	// PTY echoes input and output; "hello_pty" must appear somewhere in stdout.
	if !strings.Contains(result.Content, "hello_pty") {
		t.Errorf("expected 'hello_pty' in output, got:\n%s", result.Content)
	}
	// In PTY mode stderr is always empty.
	parts := strings.SplitN(result.Content, "[STDERR]", 2)
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		t.Errorf("expected empty stderr in PTY mode, got: %q", parts[1])
	}
}

// TestInteractivePTYExitCode verifies that non-zero exit codes propagate correctly
// through the PTY path (EIO handling on process exit).
func TestInteractivePTYExitCode(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	result := b.Execute("pty-exit42", "exit 42", 5*time.Second, 0, true)

	if !result.IsError {
		t.Errorf("expected IsError=true for exit 42, got false")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 42") {
		t.Errorf("expected exit code 42, got:\n%s", result.Content)
	}
}

// TestInteractivePTYTimeoutPauses verifies that a PTY job can be SIGSTOP'd by
// the timeout budget, returning [PAUSED] with a job ID.
func TestInteractivePTYTimeoutPauses(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	start := time.Now()
	result := b.Execute("pty-timeout", "sleep 30", 1*time.Second, 0, true)
	elapsed := time.Since(start)

	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Errorf("expected [PAUSED] after timeout, got:\n%s", result.Content)
	}
	if elapsed > 3*time.Second {
		t.Errorf("took %v, expected ~1s for timeout", elapsed)
	}

	// Clean up
	pgid := extractJobID(result.Content)
	if pgid > 0 {
		b.HandleJob("cleanup", map[string]any{"id": float64(pgid), "signal": "kill"})
	}
}

// TestInteractivePTYInputInjection verifies the core interactive workflow:
//  1. Start a shell in PTY mode with a short timeout so it pauses.
//  2. Resume with input="echo injected\n" — the shell should execute the command.
//  3. The final output should contain "injected".
func TestInteractivePTYInputInjection(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// Start an interactive shell session that will pause after 1 second.
	// We use a shell that just waits for input (read loop).
	result := b.Execute("pty-input", "cat", 1*time.Second, 0, true)

	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Skipf("job completed before pause (unexpected): %s", result.Content)
	}

	pgid := extractJobID(result.Content)
	if pgid <= 0 {
		t.Fatalf("could not extract job id from: %s", result.Content)
	}

	// Resume with input — cat should echo it back.
	resumeResult := b.HandleJob("pty-cont", map[string]any{
		"id":      float64(pgid),
		"signal":  "cont",
		"input":   "injected_text\n",
		"timeout": float64(2),
	})

	if resumeResult.IsError {
		t.Errorf("resume failed: %s", resumeResult.Content)
	}

	// Clean up if still paused
	if strings.Contains(resumeResult.Content, "[PAUSED]") {
		pgid2 := extractJobID(resumeResult.Content)
		if pgid2 > 0 {
			b.HandleJob("cleanup", map[string]any{"id": float64(pgid2), "signal": "kill"})
		}
		// Check accumulated output contains the injected text
		if !strings.Contains(resumeResult.Content, "injected_text") {
			t.Errorf("expected 'injected_text' echoed back in output, got:\n%s", resumeResult.Content)
		}
	} else {
		// Completed (cat exited after EOF) — check output
		if !strings.Contains(resumeResult.Content, "injected_text") {
			t.Errorf("expected 'injected_text' echoed back in output, got:\n%s", resumeResult.Content)
		}
	}
}

// TestInteractivePTYIsatty verifies that the child process sees isatty(0)==true
// in PTY mode, which is the primary motivation for the feature.
func TestInteractivePTYIsatty(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// python3 -c "import sys; print('TTY' if sys.stdin.isatty() else 'NOTTY')"
	result := b.Execute("pty-isatty", `python3 -c "import sys; print('TTY' if sys.stdin.isatty() else 'NOTTY')"`, 5*time.Second, 0, true)

	if result.IsError {
		t.Fatalf("command failed: %s", result.Content)
	}
	if !strings.Contains(result.Content, "TTY") {
		t.Errorf("expected isatty(0)==true (TTY) in PTY mode, got:\n%s", result.Content)
	}
	if strings.Contains(result.Content, "NOTTY") {
		t.Errorf("stdin is not a TTY in PTY mode, got:\n%s", result.Content)
	}
}

// TestNonInteractivePTYIsatty confirms the baseline: non-interactive mode has
// isatty(0)==false (stdin from /dev/null).
func TestNonInteractivePTYIsatty(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	result := b.Execute("pipe-isatty", `python3 -c "import sys; print('TTY' if sys.stdin.isatty() else 'NOTTY')"`, 5*time.Second, 0, false)

	if result.IsError {
		t.Fatalf("command failed: %s", result.Content)
	}
	if !strings.Contains(result.Content, "NOTTY") {
		t.Errorf("expected isatty(0)==false (NOTTY) in non-interactive mode, got:\n%s", result.Content)
	}
}
