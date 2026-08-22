package tools

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
)

func waitForFileContains(t *testing.T, path string, needle string, timeout time.Duration) string {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(path)
		if err == nil && strings.Contains(string(data), needle) {
			return string(data)
		}
		time.Sleep(25 * time.Millisecond)
	}
	data, _ := os.ReadFile(path)
	t.Fatalf("timed out waiting for %q in %s; last content=%q", needle, path, string(data))
	return ""
}

func waitForMeta(t *testing.T, path string, timeout time.Duration, pred func(interactiveMeta) bool) interactiveMeta {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(path)
		if err == nil {
			var meta interactiveMeta
			if json.Unmarshal(data, &meta) == nil && pred(meta) {
				return meta
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	data, _ := os.ReadFile(path)
	t.Fatalf("timed out waiting for meta predicate on %s; last content=%q", path, string(data))
	return interactiveMeta{}
}

func TestInteractiveJobMaterializesScreenFiles(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("interactive-materialize", "printf 'HELLO_INTERACTIVE\\n'; sleep 1", 0, 0, true, false, "")
	if result.IsError {
		t.Fatalf("interactive start failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	jobDir := jobDirFromExecutor(b, pid)

	for _, path := range []string{
		filepath.Join(jobDir, "pid"),
		filepath.Join(jobDir, "started_at"),
		filepath.Join(jobDir, "in"),
		filepath.Join(jobDir, "screen.txt"),
		filepath.Join(jobDir, "screen.png"),
		filepath.Join(jobDir, "screen.meta"),
		filepath.Join(jobDir, "winsize"),
		filepath.Join(jobDir, "events.log"),
		filepath.Join(jobDir, "events.hex"),
		filepath.Join(jobDir, "input.log"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
	assertJobIdentityFiles(t, jobDir, pid)
	if _, err := os.Stat(filepath.Join(jobDir, "exit")); !os.IsNotExist(err) {
		t.Fatalf("exit should be absent while interactive job is running; stat err=%v", err)
	}

	waitForFileContains(t, filepath.Join(jobDir, "screen.txt"), "HELLO_INTERACTIVE", 2*time.Second)
	waitForFileContains(t, filepath.Join(jobDir, "events.log"), "HELLO_INTERACTIVE", 2*time.Second)
	meta := waitForMeta(t, filepath.Join(jobDir, "screen.meta"), 2*time.Second, func(m interactiveMeta) bool {
		return m.Generation > 0
	})
	if meta.Rows != defaultInteractiveRows || meta.Cols != defaultInteractiveCols {
		t.Fatalf("meta rows/cols = %dx%d, want %dx%d", meta.Cols, meta.Rows, defaultInteractiveCols, defaultInteractiveRows)
	}

	exitPath := filepath.Join(jobDir, "exit")
	waitForFileExists(t, exitPath, 3*time.Second)
	exitData, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit: %v", err)
	}
	if strings.TrimSpace(string(exitData)) != "0" {
		t.Fatalf("exit = %q, want 0", string(exitData))
	}
}

func TestInteractiveInputDrivesREPL(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("interactive-repl", "python3 -q", 0, 0, true, false, "")
	if result.IsError {
		t.Fatalf("interactive start failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	jobDir := jobDirFromExecutor(b, pid)
	screenPath := filepath.Join(jobDir, "screen.txt")
	inputPath := filepath.Join(jobDir, "in")
	exitPath := filepath.Join(jobDir, "exit")

	waitForFileContains(t, screenPath, ">>>", 10*time.Second)
	assertJobIdentityFiles(t, jobDir, pid)
	if _, err := os.Stat(exitPath); !os.IsNotExist(err) {
		t.Fatalf("exit should be absent before the REPL exits; stat err=%v", err)
	}

	if err := os.WriteFile(inputPath, []byte("print(6*7)<enter>"), 0o644); err != nil {
		t.Fatalf("write interactive input: %v", err)
	}
	waitForFileContains(t, screenPath, "42", 10*time.Second)
	waitForFileContains(t, filepath.Join(jobDir, "input.log"), "print(6*7)\r", 3*time.Second)

	if err := os.WriteFile(inputPath, []byte("exit()<enter>"), 0o644); err != nil {
		t.Fatalf("write interactive exit input: %v", err)
	}
	waitForFileExists(t, exitPath, 3*time.Second)
	waitForFileContains(t, filepath.Join(jobDir, "events.hex"), "3432", 3*time.Second)
	exitData, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit: %v", err)
	}
	if strings.TrimSpace(string(exitData)) != "0" {
		t.Fatalf("exit = %q, want 0", string(exitData))
	}
}

func TestInteractiveWinsizeUpdatesMeta(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("interactive-winsize", "sleep 60", 0, 0, true, false, "")
	if result.IsError {
		t.Fatalf("interactive start failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	jobDir := jobDirFromExecutor(b, pid)
	winsizePath := filepath.Join(jobDir, "winsize")
	metaPath := filepath.Join(jobDir, "screen.meta")
	assertJobIdentityFiles(t, jobDir, pid)

	if err := os.WriteFile(winsizePath, []byte("120x40\n"), 0o644); err != nil {
		t.Fatalf("write winsize: %v", err)
	}
	meta := waitForMeta(t, metaPath, 3*time.Second, func(m interactiveMeta) bool {
		return m.Rows == 40 && m.Cols == 120
	})
	if meta.Rows != 40 || meta.Cols != 120 {
		t.Fatalf("meta rows/cols = %dx%d, want 120x40", meta.Cols, meta.Rows)
	}

	if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
		t.Fatalf("kill interactive process group: %v", err)
	}
}

func TestInteractiveSIGINTProcessGroupTriggersTrap(t *testing.T) {
	b := testExecutor()
	resultPath := filepath.Join(b.DataDir, "interactive-sigint-result.txt")
	b.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"SIGNAL_RESULT=" + resultPath,
	}
	defer b.Close(false)

	command := `trap 'printf "INT_TRAPPED\n" > "$SIGNAL_RESULT"; exit 0' INT; printf 'READY_INT\n'; while :; do sleep 1; done`
	result := b.Execute("interactive-sigint", command, 0, 0, true, false, "")
	if result.IsError {
		t.Fatalf("interactive start failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	jobDir := jobDirFromExecutor(b, pid)
	waitForFileContains(t, filepath.Join(jobDir, "screen.txt"), "READY_INT", 3*time.Second)
	assertJobIdentityFiles(t, jobDir, pid)

	if err := syscall.Kill(-pid, syscall.SIGINT); err != nil {
		t.Fatalf("send SIGINT to interactive process group: %v", err)
	}

	waitForFileContains(t, resultPath, "INT_TRAPPED", 3*time.Second)
	exitPath := filepath.Join(jobDir, "exit")
	waitForFileExists(t, exitPath, 3*time.Second)
	exitData, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit: %v", err)
	}
	if strings.TrimSpace(string(exitData)) != "0" {
		t.Fatalf("exit = %q, want 0", string(exitData))
	}
}

func TestInteractiveJobStartsWithNetworkNoneWhenNamespaceAvailable(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("network namespace isolation is Linux-only")
	}
	probe := exec.Command("true")
	probe.SysProcAttr = jobSysProcAttr(false, false, true)
	if err := probe.Run(); err != nil {
		t.Skipf("network namespace unsupported in this environment: %v", err)
	}

	b := testExecutor()
	b.Network = "none"
	defer b.Close(false)

	result := b.Execute("interactive-network-none", "printf 'NETWORK_NONE_INTERACTIVE\\n'", 0, 0, true, false, "")
	if result.IsError {
		t.Fatalf("interactive start with network=none failed: %s", result.Content)
	}
	pid := extractPID(result.Content)
	jobDir := jobDirFromExecutor(b, pid)
	waitForFileContains(t, filepath.Join(jobDir, "screen.txt"), "NETWORK_NONE_INTERACTIVE", 3*time.Second)
	waitForFileExists(t, filepath.Join(jobDir, "exit"), 3*time.Second)
}

func TestInteractiveOverlayJobCreatesAdoptableWorld(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	dataDir := t.TempDir()
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "interactive-overlay-parent"},
		Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"},
		Limits:    config.Limits{OutputTruncate: 20480, ShTimeout: 10},
		WorkspaceConfig: config.WorkspaceConfig{
			WorkspaceEnabled:      true,
			WorkspaceRoot:         root,
			Workspace:             root,
			WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
			WorkspaceSession:      "interactive-overlay-parent-session",
			WorkspaceOwner:        true,
		},
		Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"},
	}
	b := NewShExecutor(cfg)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)

	b.TurnID = 1
	seed := b.Execute("interactive-overlay-seed", "printf 'parent\\n' > parent.txt", 0, 0, false, false, "")
	if seed.IsError {
		t.Fatalf("seed write failed:\n%s", seed.Content)
	}

	b.TurnID = 2
	result := b.Execute("interactive-overlay", "cat parent.txt > seen.txt; printf 'from interactive\\n' > interactive.txt; printf 'DONE\\n'", 0, 0, true, false, "")
	if result.IsError {
		t.Fatalf("interactive overlay start failed:\n%s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	job, _ := payload["job"].(map[string]any)
	if job == nil {
		t.Fatalf("interactive result missing job payload:\n%s", result.Content)
	}
	if job["workspace_session"] == "" || job["adoptable"] != true {
		t.Fatalf("interactive overlay job should expose adoptable workspace session, got %#v", job)
	}

	pid := extractPID(result.Content)
	jobDir := jobDirFromExecutor(b, pid)
	waitForFileContains(t, filepath.Join(jobDir, "screen.txt"), "DONE", 5*time.Second)
	waitForFileExists(t, filepath.Join(jobDir, "exit"), 5*time.Second)

	handleData, err := os.ReadFile(filepath.Join(jobDir, "world_handle"))
	if err != nil {
		t.Fatalf("read interactive world_handle: %v", err)
	}
	handle := strings.TrimSpace(string(handleData))
	if !strings.HasPrefix(handle, "world://") {
		t.Fatalf("world_handle = %q, want world:// handle", handle)
	}
	mutations, err := os.ReadFile(filepath.Join(jobDir, "fs_mutations"))
	if err != nil {
		t.Fatalf("read interactive fs_mutations: %v", err)
	}
	if !strings.Contains(string(mutations), "interactive.txt (created)") || !strings.Contains(string(mutations), "seen.txt (created)") {
		t.Fatalf("interactive fs_mutations missing created files:\n%s", string(mutations))
	}
	if _, err := os.Stat(filepath.Join(root, "interactive.txt")); !os.IsNotExist(err) {
		t.Fatalf("interactive file should not be host-visible before adoption, stat err=%v", err)
	}
	before := t.TempDir()
	if err := b.subjective.exportCurrentTree(before); err != nil {
		t.Fatalf("export parent world before adoption: %v", err)
	}
	if _, err := os.Stat(filepath.Join(before, "interactive.txt")); !os.IsNotExist(err) {
		t.Fatalf("interactive file should not be in parent world before adoption, stat err=%v", err)
	}

	switched := b.SwitchWorld("interactive-overlay-switch", handle)
	if switched.IsError {
		t.Fatalf("switch to interactive world failed:\n%s", switched.Content)
	}
	if got := readOverlayWorkspaceFile(t, b.subjective, "interactive.txt"); got != "from interactive" {
		t.Fatalf("adopted interactive.txt = %q, want from interactive", got)
	}
	if got := readOverlayWorkspaceFile(t, b.subjective, "seen.txt"); got != "parent" {
		t.Fatalf("adopted seen.txt = %q, want parent", got)
	}
}
