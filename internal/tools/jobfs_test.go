package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestBlockingJobCleansUpAfterReturn(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("jobfs-blocking", "echo block_out; echo block_err >&2", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}

	// Verify output is captured in result
	if !strings.Contains(result.Content, "block_out") {
		t.Fatalf("expected stdout in result, got: %s", result.Content)
	}
	if !strings.Contains(result.Content, "block_err") {
		t.Fatalf("expected stderr in result, got: %s", result.Content)
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0") {
		t.Fatalf("expected exit code 0 in result, got: %s", result.Content)
	}

	// Non-detached jobs should be cleaned up after returning results.
	// The session directory may exist but should be empty (no job subdirs).
	sessionDir := filepath.Join(b.DataDir, "jobs", b.SessionID)
	entries, err := os.ReadDir(sessionDir)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("unexpected error reading session dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected job directory to be cleaned up, found %d entries", len(entries))
	}
}

func TestDetachWaitViaExitFile(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("jobfs-detach", "sleep 1 && echo DETACH_DONE", 0, 0, false, true, "")
	if result.IsError {
		t.Fatalf("detach failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	if pid <= 0 {
		t.Fatalf("could not extract pid from: %s", result.Content)
	}

	exitPath := filepath.Join(jobDirFromExecutor(b, pid), "exit")
	waitCmd := exec.Command("/bin/sh", "-c", "cat "+shellQuote(exitPath))
	waitOut, err := waitCmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	if err := waitCmd.Start(); err != nil {
		t.Fatalf("start cat exit: %v", err)
	}

	select {
	case <-time.After(200 * time.Millisecond):
	case <-waitProcessDone(waitCmd.Process.Pid):
		t.Fatal("cat exit returned too early; expected blocking wait")
	}

	data, err := io.ReadAll(waitOut)
	if err != nil {
		t.Fatalf("read wait output: %v", err)
	}
	if err := waitCmd.Wait(); err != nil {
		t.Fatalf("cat exit wait failed: %v", err)
	}
	if strings.TrimSpace(string(data)) != "0" {
		t.Fatalf("cat exit = %q, want 0", string(data))
	}

	outData, err := os.ReadFile(filepath.Join(jobDirFromExecutor(b, pid), "out.log"))
	if err != nil {
		t.Fatalf("reading out.log: %v", err)
	}
	if !strings.Contains(string(outData), "DETACH_DONE") {
		t.Fatalf("out.log missing DETACH_DONE: %q", string(outData))
	}
}

func TestExitFilePersistsAfterCompletion(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("jobfs-persist", "echo persist", 0, 0, false, true, "")
	if result.IsError {
		t.Fatalf("detach failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	exitPath := filepath.Join(jobDirFromExecutor(b, pid), "exit")

	first, err := exec.Command("/bin/sh", "-c", "cat "+shellQuote(exitPath)).Output()
	if err != nil {
		t.Fatalf("first cat exit: %v", err)
	}
	second, err := exec.Command("/bin/sh", "-c", "cat "+shellQuote(exitPath)).Output()
	if err != nil {
		t.Fatalf("second cat exit: %v", err)
	}
	if strings.TrimSpace(string(first)) != "0" || strings.TrimSpace(string(second)) != "0" {
		t.Fatalf("persisted exit mismatch: first=%q second=%q", first, second)
	}
}

func TestDetachKillProducesExitFile(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("jobfs-kill", "sleep 60", 0, 0, false, true, "")
	if result.IsError {
		t.Fatalf("detach failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
		t.Fatalf("kill detach group: %v", err)
	}

	exitPath := filepath.Join(jobDirFromExecutor(b, pid), "exit")
	data, err := exec.Command("/bin/sh", "-c", "cat "+shellQuote(exitPath)).Output()
	if err != nil {
		t.Fatalf("cat exit after kill: %v", err)
	}
	if strings.TrimSpace(string(data)) == "" {
		t.Fatalf("empty exit after kill")
	}
}

func TestCloseKeepDetachedPreservesProcess(t *testing.T) {
	b := testExecutor()

	result := b.Execute("jobfs-keep", "sleep 60", 0, 0, false, true, "")
	if result.IsError {
		t.Fatalf("detach failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	if err := b.Close(true); err != nil {
		t.Fatalf("close(true): %v", err)
	}
	if err := syscall.Kill(pid, 0); err != nil {
		t.Fatalf("expected detached pid %d to survive close(true): %v", pid, err)
	}
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}

func TestCloseFalseKillsDetached(t *testing.T) {
	b := testExecutor()

	result := b.Execute("jobfs-close-kill", "sleep 60", 0, 0, false, true, "")
	if result.IsError {
		t.Fatalf("detach failed: %s", result.Content)
	}

	pid := extractPID(result.Content)
	if err := b.Close(false); err != nil {
		t.Fatalf("close(false): %v", err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(pid, 0); err != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("expected detached pid %d to die after close(false)", pid)
}

func pidString(pid int) string {
	return fmt.Sprintf("%d", pid)
}

func jobDirFromExecutor(b *ShExecutor, pid int) string {
	return filepath.Join(b.DataDir, "jobs", b.SessionID, pidString(pid))
}

func shellQuote(path string) string {
	return "'" + strings.ReplaceAll(path, "'", "'\"'\"'") + "'"
}

func waitProcessDone(pid int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if err := syscall.Kill(pid, 0); err != nil {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return done
}
