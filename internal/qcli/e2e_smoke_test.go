package qcli

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/llm/codexoauth"
)

func TestE2ELiveSmoke(t *testing.T) {
	if os.Getenv("QCLI_E2E") != "1" {
		t.Skip("set QCLI_E2E=1 to run the live qcli smoke against profiles/gpt-5.5-codex-oauth.env")
	}
	if err := requireFreshCodexOAuthToken(); err != nil {
		t.Skipf("codex OAuth credentials are not non-interactively fresh: %v", err)
	}

	repoRoot := findRepoRoot(t)
	quineBin := filepath.Join(t.TempDir(), "quine")
	build := exec.Command("go", "build", "-o", quineBin, "./cmd/quine")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build quine: %v\n%s", err, out)
	}

	profileEnv := loadProfileEnv(t, repoRoot)
	runtimeRoot := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	kernel, err := NewKernel(ctx, KernelOptions{
		RuntimeRoot:  runtimeRoot,
		QuineBinary:  quineBin,
		PollInterval: 50 * time.Millisecond,
		CtlTimeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer kernel.Close()

	mission := "When you receive a message containing a reply_ctl path, immediately write the word pong to that path, then exit."
	agent, err := kernel.Spawn(SpawnOptions{
		Mission:     mission,
		RuntimeRoot: runtimeRoot,
		QuineBinary: quineBin,
		WorkDir:     repoRoot,
		Env:         profileEnv,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Stdin:       os.Stdin,
	})
	if err != nil {
		t.Fatalf("spawn: %v", err)
	}
	defer terminatePID(agent.PID)

	ch, unsub := kernel.session.Subscribe()
	defer unsub()
	drainInitial(t, ch)
	queued, err := kernel.session.Send(ControlActionInject, "qcli e2e: reply pong to the reply_ctl path, then exit")
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if queued.Type != "queued" || queued.ClientRef == "" {
		t.Fatalf("queued = %#v", queued)
	}

	waitE2EEvent[ReceiptEvent](t, ctx, ch, func(ev ReceiptEvent) bool {
		return ev.Stage == "received" && ev.ClientRef != nil && *ev.ClientRef == queued.ClientRef
	})
	waitE2EEvent[StreamDeltaEvent](t, ctx, ch, func(ev StreamDeltaEvent) bool {
		return ev.Kind == "text_delta" || ev.Kind == "reasoning_delta"
	})
	waitE2EEvent[StreamDeltaEvent](t, ctx, ch, func(ev StreamDeltaEvent) bool {
		return ev.Kind == "completed"
	})
	waitE2EEvent[Cell](t, ctx, ch, func(ev Cell) bool {
		return ev.Kind == CellMessage || ev.Kind == CellReasoning || ev.Kind == CellToolCall || ev.Kind == CellToolResult
	})
	waitE2EEvent[ReceiptEvent](t, ctx, ch, func(ev ReceiptEvent) bool {
		return ev.Stage == "delivered" && ev.ClientRef != nil && *ev.ClientRef == queued.ClientRef
	})
	waitE2EEvent[Cell](t, ctx, ch, func(ev Cell) bool {
		return ev.Kind == CellReply && ev.Text != nil && strings.TrimSpace(*ev.Text) == "pong"
	})
}

func requireFreshCodexOAuthToken() error {
	ts, ok, err := codexoauth.LoadToken("")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("no codex OAuth token found")
	}
	if ts.AccessToken == "" {
		return errors.New("codex OAuth token has no access token")
	}
	if ts.ExpiresAt <= time.Now().Add(2*time.Minute).UnixMilli() {
		return errors.New("codex OAuth access token is expired or too close to expiry")
	}
	return nil
}

func loadProfileEnv(t *testing.T, repoRoot string) []string {
	t.Helper()
	profile := filepath.Join(repoRoot, "profiles", "gpt-5.5-codex-oauth.env")
	cmd := exec.Command("bash", "-c", "set -a; source \"$1\"; set +a; env", "bash", profile)
	cmd.Dir = repoRoot
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("load profile env: %v", err)
	}
	env := envLines(out)
	env = setEnvKey(env, "QUINE_IDLE_ENABLED", "1")
	env = setEnvKey(env, "QUINE_MAX_TURNS", "6")
	env = setEnvKey(env, "QUINE_WALL_CLOCK_EXIT_SECONDS", "150")
	return env
}

func envLines(data []byte) []string {
	var env []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			env = append(env, line)
		}
	}
	return env
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func terminatePID(pid int) {
	if pid <= 0 {
		return
	}
	_ = syscall.Kill(pid, syscall.SIGTERM)
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if !pidLive(pid) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	_ = syscall.Kill(pid, syscall.SIGKILL)
}

func waitE2EEvent[T any](t *testing.T, ctx context.Context, ch <-chan any, match func(T) bool) T {
	t.Helper()
	for {
		select {
		case ev := <-ch:
			if typed, ok := ev.(T); ok && match(typed) {
				return typed
			}
		case <-ctx.Done():
			var zero T
			t.Fatalf("timed out waiting for E2E event %T: %v", zero, ctx.Err())
		}
	}
}
