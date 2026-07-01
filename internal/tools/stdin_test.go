package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// --- stdin parameter tests ---

// TestStdinBasic verifies that stdin data is piped into the command.
func TestStdinBasic(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("stdin-basic", "cat", 0, 0, false, false, "hello from stdin\n")

	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	if !strings.Contains(resultString(t, payload, "stdout"), "hello from stdin") {
		t.Errorf("expected stdin content in output, got:\n%s", result.Content)
	}
	if len(result.StructuredContent) != 0 {
		t.Fatalf("expected empty structured content, got %q", result.StructuredContent)
	}
}

// TestStdinNoEscaping verifies that special shell characters in stdin are
// delivered verbatim — no escaping needed.
func TestStdinNoEscaping(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	// Content with quotes, backslashes, dollar signs — all would break heredoc
	tricky := "key = \"value\"\npath = C:\\data\nprice = $100\n"
	result := b.Execute("stdin-special", "cat", 0, 0, false, false, tricky)

	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	stdout := resultString(t, payload, "stdout")
	if !strings.Contains(stdout, `key = "value"`) {
		t.Errorf("expected double-quoted value, got:\n%s", result.Content)
	}
	if !strings.Contains(stdout, `path = C:\data`) {
		t.Errorf("expected backslash path, got:\n%s", result.Content)
	}
	if !strings.Contains(stdout, "price = $100") {
		t.Errorf("expected dollar sign (no expansion), got:\n%s", result.Content)
	}
}

// TestStdinWriteFile verifies the canonical use case: creating a file with
// content that would be hard to escape in a shell command.
func TestStdinWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "script.py")

	content := "#!/usr/bin/env python3\nprint('hello, world!')\nprint(\"it's alive\")\n"

	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("stdin-write", "cat > "+outFile, 0, 0, false, false, content)
	if result.IsError {
		t.Fatalf("write failed: %s", result.Content)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading written file: %v", err)
	}
	if string(data) != content {
		t.Errorf("file content mismatch.\ngot:  %q\nwant: %q", string(data), content)
	}
}

// TestStdinPipedToProgram verifies that stdin is consumed by a program
// that reads from stdin (not just cat).
func TestStdinPipedToProgram(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	// wc -l counts lines — should report 3
	result := b.Execute("stdin-wc", "wc -l", 0, 0, false, false, "line1\nline2\nline3\n")
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, result.Content), "stdout"), "3") {
		t.Errorf("expected wc -l to count 3 lines, got:\n%s", result.Content)
	}
}

// TestStdinEmptyString verifies that empty stdin ("") behaves normally
// — no stdin pipe is wired, command runs as usual.
func TestStdinEmptyString(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("stdin-empty", "echo hello", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, result.Content), "stdout"), "hello") {
		t.Errorf("expected hello in output, got:\n%s", result.Content)
	}
}

// TestStdinWithMultilineScript verifies piping a full script to an interpreter.
func TestStdinWithMultilineScript(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	script := "import sys\nprint('version:', sys.version_info.major)\nprint('done')\n"
	result := b.Execute("stdin-python", "python3 -", 0, 0, false, false, script)
	if result.IsError {
		t.Fatalf("python3 failed: %s", result.Content)
	}
	stdout := resultString(t, decodeResultContent(t, result.Content), "stdout")
	if !strings.Contains(stdout, "version:") {
		t.Errorf("expected 'version:' in output, got:\n%s", result.Content)
	}
	if !strings.Contains(stdout, "done") {
		t.Errorf("expected 'done' in output, got:\n%s", result.Content)
	}
}

// TestStdinCompatibleWithDetach verifies stdin works with detach=true.
func TestStdinCompatibleWithDetach(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "detached.txt")

	b := testExecutor()
	defer b.Close(false)

	// Write stdin content to a file via a detached job
	result := b.Execute("stdin-detach", "cat > "+outFile, 0, 0, false, true, "detached content\n")
	if result.IsError {
		t.Fatalf("detach launch failed: %s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	if payload["mode"] != "detached" {
		t.Fatalf("mode = %#v, want detached", payload["mode"])
	}
	job := resultMap(t, payload, "job")
	jobID := resultInt(t, job, "pid")
	if jobID <= 0 {
		t.Fatalf("could not extract job ID from: %s", result.Content)
	}

	exitPath := filepath.Join(jobDirFromExecutor(b, jobID), "exit")
	waitForFileExists(t, exitPath, 3*time.Second)
	if out, err := os.ReadFile(exitPath); err != nil {
		t.Fatalf("reading exit: %v", err)
	} else if strings.TrimSpace(string(out)) != "0" {
		t.Fatalf("exit code = %q, want 0", out)
	}

	// Check the file was written
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "detached content") {
		t.Errorf("file content = %q, want 'detached content'", string(data))
	}
}
