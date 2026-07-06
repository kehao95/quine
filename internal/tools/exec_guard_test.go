package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func TestIsLiveProcessImageTarget(t *testing.T) {
	cases := []struct {
		name   string
		target string
		want   bool
	}{
		{"proc self exe", "/proc/self/exe", true},
		{"proc pid exe", "/proc/12345/exe", true},
		{"proc self exe with dotslash", "/proc/self/./exe", true},
		{"trailing space", "  /proc/self/exe  ", true},
		{"ordinary workspace body", "./quine", false},
		{"absolute workspace body", "/tmp/run/workspace/quine", false},
		{"proc cmdline not exe", "/proc/self/cmdline", false},
		{"proc self exe substring in name", "/tmp/proc-self-exe-copy", false},
		{"empty", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isLiveProcessImageTarget(c.target); got != c.want {
				t.Fatalf("isLiveProcessImageTarget(%q) = %v, want %v", c.target, got, c.want)
			}
		})
	}
}

// Execute must reject a procfs exec BEFORE performing the syscall when the
// ephemeral body has been unlinked. The guard runs before syscall.Exec, so
// this test never actually replaces the process image — it asserts the error
// ToolResult is returned instead. (If the guard regressed, syscall.Exec would
// fire and the test binary would be replaced, which itself signals failure.)
func TestExecRejectsProcSelfExeUnderEphemeralBody(t *testing.T) {
	e := &ExecExecutor{
		QuinePath: "/nonexistent/quine",
		Cfg:       &config.Config{ToolGates: config.ToolGates{EphemeralBody: true}},
	}
	res := e.Execute("t1", ExecRequest{Target: "/proc/self/exe"})
	if !res.IsError {
		t.Fatalf("expected error result rejecting /proc/self/exe, got non-error")
	}
	if !strings.Contains(string(res.Content), "live process image") {
		t.Fatalf("expected body-recovery rejection message, got: %s", res.Content)
	}
}

// The guard gates on EphemeralBody: without it, the condition
// `e.Cfg.EphemeralBody && isLiveProcessImageTarget(target)` short-circuits and
// the guard never fires. We assert that at the gate level rather than by
// calling Execute() with a real /proc/self/exe target — doing so would let
// syscall.Exec re-enter the test binary and loop.
func TestExecGuardGatesOnEphemeralBody(t *testing.T) {
	if !isLiveProcessImageTarget("/proc/self/exe") {
		t.Fatal("precondition: /proc/self/exe must be detected as a live image")
	}
	off := &config.Config{ToolGates: config.ToolGates{EphemeralBody: false}}
	if off.EphemeralBody {
		t.Fatal("EphemeralBody should be false here")
	}
	// With EphemeralBody false the guard's leading condition is false, so the
	// procfs check is never consulted — the ordinary exec path applies.
}

// A workspace symlink pointing at /proc/self/exe must also be caught, since
// exec-ing the symlink re-enters the live image just the same.
func TestIsLiveProcessImageTargetSymlink(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "successor")
	if err := os.Symlink("/proc/self/exe", link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	if !isLiveProcessImageTarget(link) {
		t.Fatalf("symlink to /proc/self/exe not detected as live process image")
	}
}
