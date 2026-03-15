package tools

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
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
		filepath.Join(jobDir, "in"),
		filepath.Join(jobDir, "screen.txt"),
		filepath.Join(jobDir, "screen.png"),
		filepath.Join(jobDir, "screen.meta"),
		filepath.Join(jobDir, "winsize"),
		filepath.Join(jobDir, "events.log"),
		filepath.Join(jobDir, "exit"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	waitForFileContains(t, filepath.Join(jobDir, "screen.txt"), "HELLO_INTERACTIVE", 2*time.Second)
	waitForFileContains(t, filepath.Join(jobDir, "events.log"), "HELLO_INTERACTIVE", 2*time.Second)
	meta := waitForMeta(t, filepath.Join(jobDir, "screen.meta"), 2*time.Second, func(m interactiveMeta) bool {
		return m.Generation > 0
	})
	if meta.Rows != defaultInteractiveRows || meta.Cols != defaultInteractiveCols {
		t.Fatalf("meta rows/cols = %dx%d, want %dx%d", meta.Cols, meta.Rows, defaultInteractiveCols, defaultInteractiveRows)
	}

	exitData, err := exec.Command("/bin/sh", "-c", "cat "+shellQuote(filepath.Join(jobDir, "exit"))).Output()
	if err != nil {
		t.Fatalf("cat exit: %v", err)
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

	waitForFileContains(t, screenPath, ">>>", 3*time.Second)

	if err := os.WriteFile(inputPath, []byte("print(6*7)<enter>"), 0o644); err != nil {
		t.Fatalf("write interactive input: %v", err)
	}
	waitForFileContains(t, screenPath, "42", 3*time.Second)

	if err := os.WriteFile(inputPath, []byte("exit()<enter>"), 0o644); err != nil {
		t.Fatalf("write interactive exit input: %v", err)
	}
	exitData, err := exec.Command("/bin/sh", "-c", "cat "+shellQuote(exitPath)).Output()
	if err != nil {
		t.Fatalf("cat exit: %v", err)
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
