package tools

import (
	"fmt"
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
	payload := decodeResultContent(t, result.Content)
	if resultString(t, payload, "stdout") != "block_out\n" {
		t.Fatalf("expected stdout in result, got: %s", result.Content)
	}
	if resultString(t, payload, "stderr") != "block_err\n" {
		t.Fatalf("expected stderr in result, got: %s", result.Content)
	}
	if resultInt(t, payload, "exit_code") != 0 {
		t.Fatalf("expected exit code 0 in result, got: %s", result.Content)
	}
	if len(result.StructuredContent) != 0 {
		t.Fatalf("expected empty structured content, got %q", result.StructuredContent)
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

func TestBlockingJobTimeoutInterruptsAndCanBeResumed(t *testing.T) {
	b := testExecutor()
	// Give the shell enough runway to emit the initial partial output even when
	// this test runs under nested go test load, while still forcing interruption
	// well before the command can finish naturally.
	b.Timeout = time.Second
	defer b.Close(false)

	start := time.Now()
	result := b.Execute("jobfs-timeout", "printf start; sleep 2; printf end", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("expected interrupted note, got error: %s", result.Content)
	}
	if elapsed := time.Since(start); elapsed >= 3*time.Second {
		t.Fatalf("timeout path took too long: %v", elapsed)
	}

	payload := decodeResultContent(t, result.Content)
	if got := resultString(t, payload, "status"); got != "interrupted" {
		t.Fatalf("status = %q, want interrupted", got)
	}
	if got := resultString(t, payload, "cause"); got != "timeout" {
		t.Fatalf("cause = %q, want timeout", got)
	}
	if timeoutSeconds := resultInt(t, payload, "timeout_seconds"); timeoutSeconds != 1 {
		t.Fatalf("timeout_seconds = %d, want 1", timeoutSeconds)
	}
	if stdout := resultString(t, payload, "stdout_so_far"); stdout != "start" {
		t.Fatalf("stdout_so_far = %q, want %q", stdout, "start")
	}
	if stderr := resultString(t, payload, "stderr_so_far"); stderr != "" {
		t.Fatalf("stderr_so_far = %q, want empty", stderr)
	}
	if _, ok := payload["exit_code"]; ok {
		t.Fatalf("interrupted result should omit exit_code: %#v", payload)
	}

	job := resultMap(t, payload, "job")
	pid := resultInt(t, job, "pid")
	jobDir := strings.TrimSuffix(resultString(t, job, "path"), "/")
	assertDetachedJobSurface(t, jobDir, pid)
	waitForProcessStateContains(t, pid, "T", 2*time.Second)

	exitPath := filepath.Join(jobDir, "exit")
	if _, err := os.Stat(exitPath); !os.IsNotExist(err) {
		t.Fatalf("exit should be absent while interrupted job is stopped; stat err=%v", err)
	}

	sessionDir := filepath.Join(b.DataDir, "jobs", b.SessionID)
	entries, err := os.ReadDir(sessionDir)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("unexpected error reading session dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected interrupted job directory to be retained, found %d entries", len(entries))
	}

	resume := b.Execute("jobfs-resume", fmt.Sprintf("python3 -c 'import os, signal; os.killpg(%d, signal.SIGCONT)'", pid), 5*time.Second, 0, false, false, "")
	if resume.IsError {
		t.Fatalf("resume failed: %s", resume.Content)
	}
	resumePayload := decodeResultContent(t, resume.Content)
	if got := resultString(t, resumePayload, "status"); got != "completed" {
		t.Fatalf("resume status = %q, want completed; content=%v", got, resume.Content)
	}
	if exitCode := resultInt(t, resumePayload, "exit_code"); exitCode != 0 {
		t.Fatalf("resume exit_code = %d, want 0; content=%v", exitCode, resume.Content)
	}

	waitForFileExists(t, exitPath, 3*time.Second)
	outData, err := os.ReadFile(filepath.Join(jobDir, "out.log"))
	if err != nil {
		t.Fatalf("reading out.log after resume: %v", err)
	}
	if string(outData) != "startend" {
		t.Fatalf("out.log = %q, want %q", string(outData), "startend")
	}
	data, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit after resume: %v", err)
	}
	if strings.TrimSpace(string(data)) != "0" {
		t.Fatalf("exit after resume = %q, want 0", string(data))
	}
}

func TestInterruptedJobCanBeTerminatedViaOrdinarySh(t *testing.T) {
	b := testExecutor()
	b.Timeout = 150 * time.Millisecond
	defer b.Close(false)

	result := b.Execute("jobfs-interrupt-term", "sleep 30", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("expected interrupted note, got error: %s", result.Content)
	}

	payload := decodeResultContent(t, result.Content)
	job := resultMap(t, payload, "job")
	pid := resultInt(t, job, "pid")
	jobDir := strings.TrimSuffix(resultString(t, job, "path"), "/")
	waitForProcessStateContains(t, pid, "T", 2*time.Second)

	terminate := b.Execute("jobfs-term", fmt.Sprintf("python3 -c 'import os, signal; os.killpg(%d, signal.SIGKILL)'", pid), 5*time.Second, 0, false, false, "")
	if terminate.IsError {
		t.Fatalf("terminate failed: %s", terminate.Content)
	}
	terminatePayload := decodeResultContent(t, terminate.Content)
	if got := resultString(t, terminatePayload, "status"); got != "completed" {
		t.Fatalf("terminate status = %q, want completed; content=%v", got, terminate.Content)
	}
	if exitCode := resultInt(t, terminatePayload, "exit_code"); exitCode != 0 {
		t.Fatalf("terminate exit_code = %d, want 0; content=%v", exitCode, terminate.Content)
	}

	exitPath := filepath.Join(jobDir, "exit")
	waitForFileExists(t, exitPath, 3*time.Second)
	data, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit after terminate: %v", err)
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "0" {
		t.Fatalf("exit after terminate = %q, want non-zero", trimmed)
	}
}

func TestCloseTrueKillsInterruptedPendingJobs(t *testing.T) {
	b := testExecutor()
	b.Timeout = 150 * time.Millisecond

	result := b.Execute("jobfs-close-pending", "sleep 30", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("expected interrupted note, got error: %s", result.Content)
	}

	payload := decodeResultContent(t, result.Content)
	job := resultMap(t, payload, "job")
	pid := resultInt(t, job, "pid")
	waitForProcessStateContains(t, pid, "T", 2*time.Second)

	if err := b.Close(true); err != nil {
		t.Fatalf("close(true): %v", err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(pid, 0); err != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("expected interrupted pid %d to die after close(true)", pid)
}

func TestDetachPublishesFactSurfaceImmediately(t *testing.T) {
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

	jobDir := jobDirFromExecutor(b, pid)
	assertDetachedJobSurface(t, jobDir, pid)

	exitPath := filepath.Join(jobDir, "exit")
	if _, err := os.Stat(exitPath); !os.IsNotExist(err) {
		t.Fatalf("exit should be absent while job is still running; stat err=%v", err)
	}

	waitForFileExists(t, exitPath, 3*time.Second)
	data, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit file: %v", err)
	}
	if strings.TrimSpace(string(data)) != "0" {
		t.Fatalf("exit = %q, want 0", string(data))
	}

	outData, err := os.ReadFile(filepath.Join(jobDir, "out.log"))
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
	jobDir := jobDirFromExecutor(b, pid)
	assertDetachedJobSurface(t, jobDir, pid)
	exitPath := filepath.Join(jobDir, "exit")
	waitForFileExists(t, exitPath, 3*time.Second)

	first, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("first read exit: %v", err)
	}
	second, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("second read exit: %v", err)
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
	jobDir := jobDirFromExecutor(b, pid)
	assertDetachedJobSurface(t, jobDir, pid)
	if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
		t.Fatalf("kill detach group: %v", err)
	}

	exitPath := filepath.Join(jobDir, "exit")
	waitForFileExists(t, exitPath, 15*time.Second)
	data, err := os.ReadFile(exitPath)
	if err != nil {
		t.Fatalf("read exit after kill: %v", err)
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		t.Fatalf("empty exit after kill")
	}
	if trimmed == "0" {
		t.Fatalf("exit after kill = %q, want non-zero", trimmed)
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

func assertJobIdentityFiles(t *testing.T, jobDir string, pid int) {
	t.Helper()

	cmdData, err := os.ReadFile(filepath.Join(jobDir, "cmd"))
	if err != nil {
		t.Fatalf("read cmd file: %v", err)
	}
	if strings.TrimSpace(string(cmdData)) == "" {
		t.Fatalf("cmd file should not be empty")
	}

	pidData, err := os.ReadFile(filepath.Join(jobDir, "pid"))
	if err != nil {
		t.Fatalf("read pid file: %v", err)
	}
	if got := strings.TrimSpace(string(pidData)); got != pidString(pid) {
		t.Fatalf("pid file = %q, want %q", got, pidString(pid))
	}

	startedAtData, err := os.ReadFile(filepath.Join(jobDir, "started_at"))
	if err != nil {
		t.Fatalf("read started_at file: %v", err)
	}
	if _, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(string(startedAtData))); err != nil {
		t.Fatalf("parse started_at: %v", err)
	}

}

func assertDetachedJobSurface(t *testing.T, jobDir string, pid int) {
	t.Helper()
	assertJobIdentityFiles(t, jobDir, pid)
	for _, path := range []string{
		filepath.Join(jobDir, "out.log"),
		filepath.Join(jobDir, "err.log"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func waitForFileExists(t *testing.T, path string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s after %v", path, timeout)
}

func waitForProcessStateContains(t *testing.T, pid int, want string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		out, err := exec.Command("ps", "-o", "stat=", "-p", fmt.Sprintf("%d", pid)).Output()
		if err == nil && strings.Contains(string(out), want) {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for pid %d state to contain %q", pid, want)
}
