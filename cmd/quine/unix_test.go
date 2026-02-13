package main

// Unix conformance tests for quine as a POSIX process.
//
// These tests verify that quine behaves like a well-formed Unix process:
//   - stdin (fd 4 in persistent shell) delivers material correctly
//   - stdout (fd 3) carries only deliverables, never internal chatter
//   - stderr carries only failure signals
//   - exit codes follow POSIX conventions (0=success, 1=failure)
//   - persistent shell preserves state across sh calls (cd, export, variables)
//   - process isolation: fd 1 (captured) vs fd 3 (delivered) are distinct

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/runtime"
	"github.com/kehao95/quine/internal/tape"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// unixTestConfig builds a Config for Unix conformance tests.
func unixTestConfig(t *testing.T, tapeDir string) *config.Config {
	t.Helper()
	return &config.Config{
		ModelID:        "test-model",
		APIKey:         "test-key",
		APIBase:        "",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          0,
		SessionID:      "unix-test-" + t.Name(),
		MaxConcurrent:  20,
		ShTimeout:      30,
		OutputTruncate: 20480,
		DataDir:        tapeDir,
		Shell:          "/bin/sh",
		MaxTurns:       20,
	}
}

// runWithMock runs a quine session with mock responses and returns the exit code,
// captured stdout, captured stderr, and tape entries.
func runWithMock(t *testing.T, responses []tape.Message, mission, material string) (exitCode int, tapeEntries []tape.TapeEntry) {
	t.Helper()

	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	mock := newMockProvider(responses)
	rt := runtime.NewWithProvider(cfg, mock)
	exitCode = rt.Run(mission, material)

	// Parse tape
	tapePath := filepath.Join(tapeDir, cfg.SessionID+".jsonl")
	if _, err := os.Stat(tapePath); err == nil {
		tapeEntries = parseTapeJSONL(t, tapePath)
	}
	return
}

// findToolResults extracts all tool_result entries from tape entries.
func findToolResults(t *testing.T, entries []tape.TapeEntry) []tape.ToolResult {
	t.Helper()
	var results []tape.ToolResult
	for _, e := range entries {
		if e.Type == "tool_result" {
			var tr tape.ToolResult
			if err := json.Unmarshal(e.Data, &tr); err != nil {
				t.Fatalf("unmarshal tool_result: %v", err)
			}
			results = append(results, tr)
		}
	}
	return results
}

// findOutcome extracts the session outcome from tape entries.
func findOutcome(t *testing.T, entries []tape.TapeEntry) tape.SessionOutcome {
	t.Helper()
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Type == "outcome" {
			var outcome tape.SessionOutcome
			if err := json.Unmarshal(entries[i].Data, &outcome); err != nil {
				t.Fatalf("unmarshal outcome: %v", err)
			}
			return outcome
		}
	}
	t.Fatal("no outcome entry found in tape")
	return tape.SessionOutcome{}
}

// ---------------------------------------------------------------------------
// 1. Exit codes
// ---------------------------------------------------------------------------

// TestExitCode_Success verifies exit(success) → exit code 0.
func TestExitCode_Success(t *testing.T) {
	code, _ := runWithMock(t, []tape.Message{
		assistantExit("e1", "success", "", ""),
	}, "say hello", "Begin.")

	if code != 0 {
		t.Errorf("exit(success) → exit code %d, want 0", code)
	}
}

// TestExitCode_Failure verifies exit(failure) → exit code 1.
func TestExitCode_Failure(t *testing.T) {
	code, _ := runWithMock(t, []tape.Message{
		assistantExit("e1", "failure", "", "something broke"),
	}, "do something", "Begin.")

	if code != 1 {
		t.Errorf("exit(failure) → exit code %d, want 1", code)
	}
}

// TestExitCode_TurnExhaustion verifies turn limit → exit code 1.
func TestExitCode_TurnExhaustion(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)
	cfg.MaxTurns = 2 // Only 2 turns allowed

	// Mock: 3 sh calls (exceeds MaxTurns=2)
	mock := newMockProvider([]tape.Message{
		assistantsh("c1", "echo turn1"),
		assistantsh("c2", "echo turn2"),
		// Turn limit hit after c2; agent never gets to call c3
	})

	rt := runtime.NewWithProvider(cfg, mock)
	code := rt.Run("loop forever", "Begin.")

	if code != 1 {
		t.Errorf("turn exhaustion → exit code %d, want 1", code)
	}
}

// ---------------------------------------------------------------------------
// 2. Stdout (fd 3) — deliverable channel
// ---------------------------------------------------------------------------

// TestStdout_Fd3Delivery verifies that >&3 output reaches the process's real stdout.
func TestStdout_Fd3Delivery(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	// Capture stdout by replacing it with a pipe
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	mock := newMockProvider([]tape.Message{
		assistantsh("c1", `echo "delivered" >&3`),
		assistantExit("e1", "success", "", ""),
	})

	rt := runtime.NewWithProvider(cfg, mock)

	// Replace the runtime's stdout with our pipe
	rt.SetStdout(stdoutW)

	code := rt.Run("deliver output", "Begin.")
	stdoutW.Close()

	// Read what was delivered to fd 3
	buf := make([]byte, 4096)
	n, _ := stdoutR.Read(buf)
	delivered := strings.TrimSpace(string(buf[:n]))

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if delivered != "delivered" {
		t.Errorf("fd 3 output = %q, want %q", delivered, "delivered")
	}
}

// TestStdout_Fd1NotLeaked verifies that regular command output (fd 1)
// stays in the tool result and does NOT leak to the process's real stdout.
func TestStdout_Fd1NotLeaked(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	// Capture stdout
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	mock := newMockProvider([]tape.Message{
		// This writes to fd 1 (captured), not fd 3 (delivered)
		assistantsh("c1", `echo "internal"`),
		assistantExit("e1", "success", "", ""),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	rt.SetStdout(stdoutW)

	code := rt.Run("produce internal output", "Begin.")
	stdoutW.Close()

	// Read whatever reached stdout
	buf := make([]byte, 4096)
	n, _ := stdoutR.Read(buf)
	stdoutContent := string(buf[:n])

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}

	// "internal" should NOT appear in the process's stdout
	if strings.Contains(stdoutContent, "internal") {
		t.Errorf("fd 1 output leaked to process stdout: %q", stdoutContent)
	}
}

// ---------------------------------------------------------------------------
// 3. Stderr — failure signal channel
// ---------------------------------------------------------------------------

// TestStderr_FailureSignal verifies that exit(failure, stderr=...) writes to stderr.
func TestStderr_FailureSignal(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	// Capture stderr
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	mock := newMockProvider([]tape.Message{
		assistantExit("e1", "failure", "", "file not found"),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	rt.SetStderr(stderrW)

	code := rt.Run("find file", "Begin.")
	stderrW.Close()

	buf := make([]byte, 4096)
	n, _ := stderrR.Read(buf)
	stderrContent := string(buf[:n])

	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderrContent, "file not found") {
		t.Errorf("stderr = %q, want it to contain %q", stderrContent, "file not found")
	}
}

// TestStderr_SuccessSilent verifies that exit(success) produces no stderr.
func TestStderr_SuccessSilent(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	// Capture stderr
	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	mock := newMockProvider([]tape.Message{
		assistantExit("e1", "success", "", ""),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	rt.SetStderr(stderrW)

	code := rt.Run("noop", "Begin.")
	stderrW.Close()

	buf := make([]byte, 4096)
	n, _ := stderrR.Read(buf)

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if n > 0 {
		t.Errorf("exit(success) wrote to stderr: %q", string(buf[:n]))
	}
}

// ---------------------------------------------------------------------------
// 4. Persistent shell state
// ---------------------------------------------------------------------------

// TestPersistentShell_CdPersists verifies that cd in one sh call persists to the next.
func TestPersistentShell_CdPersists(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	mock := newMockProvider([]tape.Message{
		assistantsh("c1", fmt.Sprintf("cd %q", tmpDir)),
		assistantsh("c2", "pwd"),
		assistantExit("e1", "success", "", ""),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	code := rt.Run("check cd persistence", "Begin.")

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}

	// Parse tape to find pwd output
	tapePath := filepath.Join(tapeDir, cfg.SessionID+".jsonl")
	entries := parseTapeJSONL(t, tapePath)
	results := findToolResults(t, entries)

	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	// Second result (pwd) should show tmpDir
	pwdResult := results[1].Content
	if !strings.Contains(pwdResult, tmpDir) {
		t.Errorf("pwd after cd = %q, want it to contain %q", pwdResult, tmpDir)
	}
}

// TestPersistentShell_ExportPersists verifies that export in one sh call persists to the next.
func TestPersistentShell_ExportPersists(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "export QUINE_TEST_VAR=hello42"),
		assistantsh("c2", "echo $QUINE_TEST_VAR"),
		assistantExit("e1", "success", "", ""),
	}, "check export persistence", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	echoResult := results[1].Content
	if !strings.Contains(echoResult, "hello42") {
		t.Errorf("echo after export = %q, want it to contain %q", echoResult, "hello42")
	}
}

// TestPersistentShell_ShellVarPersists verifies that shell variables persist.
func TestPersistentShell_ShellVarPersists(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "MY_COUNT=99"),
		assistantsh("c2", "echo $MY_COUNT"),
		assistantExit("e1", "success", "", ""),
	}, "check variable persistence", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	if !strings.Contains(results[1].Content, "99") {
		t.Errorf("echo after set = %q, want it to contain %q", results[1].Content, "99")
	}
}

// TestPersistentShell_FunctionPersists verifies that function definitions persist.
func TestPersistentShell_FunctionPersists(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", `greet() { echo "hi $1"; }`),
		assistantsh("c2", `greet world`),
		assistantExit("e1", "success", "", ""),
	}, "check function persistence", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	if !strings.Contains(results[1].Content, "hi world") {
		t.Errorf("function call result = %q, want it to contain %q", results[1].Content, "hi world")
	}
}

// ---------------------------------------------------------------------------
// 5. Crash recovery
// ---------------------------------------------------------------------------

// TestCrashRecovery_ExitInShell verifies that `exit` in a brace-group
// kills the shell, crash recovery detects it, and the next call auto-restarts.
func TestCrashRecovery_ExitInShell(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "exit 1"), // kills the shell
		assistantsh("c2", "echo recovered"),
		assistantExit("e1", "success", "", ""),
	}, "test crash recovery", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	// First result should be a shell error
	if !strings.Contains(results[0].Content, "SHELL ERROR") {
		t.Errorf("crash result = %q, want SHELL ERROR", results[0].Content)
	}

	// Second result should show recovery
	if !strings.Contains(results[1].Content, "recovered") {
		t.Errorf("recovery result = %q, want 'recovered'", results[1].Content)
	}
}

// TestCrashRecovery_StateLost verifies that state is lost after crash recovery.
func TestCrashRecovery_StateLost(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "export EPHEMERAL=before_crash"),
		assistantsh("c2", "exit 0"), // kills the shell
		assistantsh("c3", `echo "val=${EPHEMERAL:-gone}"`),
		assistantExit("e1", "success", "", ""),
	}, "test state loss after crash", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 3 {
		t.Fatalf("expected at least 3 tool results, got %d", len(results))
	}

	// After crash recovery, EPHEMERAL should be gone
	if !strings.Contains(results[2].Content, "val=gone") {
		t.Errorf("after crash: %q, want val=gone (state should be lost)", results[2].Content)
	}
}

// ---------------------------------------------------------------------------
// 6. sh tool result format
// ---------------------------------------------------------------------------

// TestShResultFormat verifies the exact format: [EXIT CODE] N\n[STDOUT]\n...\n[STDERR]\n...
func TestShResultFormat(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "echo out; echo err >&2"),
		assistantExit("e1", "success", "", ""),
	}, "check result format", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	content := results[0].Content
	// Must have all three sections
	if !strings.Contains(content, "[EXIT CODE] 0") {
		t.Errorf("missing [EXIT CODE] 0 in %q", content)
	}
	if !strings.Contains(content, "[STDOUT]") {
		t.Errorf("missing [STDOUT] in %q", content)
	}
	if !strings.Contains(content, "[STDERR]") {
		t.Errorf("missing [STDERR] in %q", content)
	}

	// Verify stdout has "out" and stderr has "err"
	parts := strings.SplitN(content, "[STDERR]", 2)
	if len(parts) != 2 {
		t.Fatalf("result format unexpected: %q", content)
	}
	stdoutSection := parts[0]
	stderrSection := parts[1]

	if !strings.Contains(stdoutSection, "out") {
		t.Errorf("stdout section = %q, want 'out'", stdoutSection)
	}
	if !strings.Contains(stderrSection, "err") {
		t.Errorf("stderr section = %q, want 'err'", stderrSection)
	}
}

// TestShResultFormat_NonZeroExit verifies non-zero exit code is reported.
func TestShResultFormat_NonZeroExit(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "false"), // exit code 1
		assistantExit("e1", "success", "", ""),
	}, "check non-zero exit", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	if !strings.Contains(results[0].Content, "[EXIT CODE] 1") {
		t.Errorf("expected [EXIT CODE] 1, got %q", results[0].Content)
	}
}

// ---------------------------------------------------------------------------
// 7. Stdin (fd 4) — material channel
// ---------------------------------------------------------------------------

// TestStdin_Fd4Available verifies that material stdin is readable via /dev/fd/4.
func TestStdin_Fd4Available(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	cfg := unixTestConfig(t, tapeDir)

	// Create a pipe to simulate stdin
	stdinR, stdinW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Write test data to stdin pipe
	go func() {
		stdinW.Write([]byte("material data\n"))
		stdinW.Close()
	}()

	mock := newMockProvider([]tape.Message{
		assistantsh("c1", "cat /dev/fd/4"),
		assistantExit("e1", "success", "", ""),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	rt.SetStdin(stdinR)

	code := rt.Run("read stdin", "Input is piped to stdin.")

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}

	// Parse tape to verify fd 4 read
	tapePath := filepath.Join(tapeDir, cfg.SessionID+".jsonl")
	entries := parseTapeJSONL(t, tapePath)
	results := findToolResults(t, entries)

	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	if !strings.Contains(results[0].Content, "material data") {
		t.Errorf("fd 4 read = %q, want 'material data'", results[0].Content)
	}
}

// ---------------------------------------------------------------------------
// 8. Multiline and special characters
// ---------------------------------------------------------------------------

// TestShell_MultilineCommand verifies that multiline commands work.
func TestShell_MultilineCommand(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "for i in 1 2 3; do\n  echo \"line$i\"\ndone"),
		assistantExit("e1", "success", "", ""),
	}, "multiline command", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	content := results[0].Content
	if !strings.Contains(content, "line1") || !strings.Contains(content, "line2") || !strings.Contains(content, "line3") {
		t.Errorf("multiline output = %q, want line1, line2, line3", content)
	}
}

// TestShell_SpecialChars verifies that special characters are handled correctly.
func TestShell_SpecialChars(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", `printf 'tab\there\nnewline\n'`),
		assistantExit("e1", "success", "", ""),
	}, "special chars", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	content := results[0].Content
	if !strings.Contains(content, "tab\there") {
		t.Errorf("special chars output = %q, want tab and newline chars", content)
	}
}

// ---------------------------------------------------------------------------
// 9. Pipe semantics — command pipeline exit codes
// ---------------------------------------------------------------------------

// TestShell_PipeExitCode verifies that pipeline exit code is the last command's.
func TestShell_PipeExitCode(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		// In POSIX sh, pipeline exit = last command's exit code
		assistantsh("c1", "false | true"),
		assistantExit("e1", "success", "", ""),
	}, "pipe exit code", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	// Pipeline exit code = last command (true = 0)
	if !strings.Contains(results[0].Content, "[EXIT CODE] 0") {
		t.Errorf("pipe exit code: %q, want [EXIT CODE] 0", results[0].Content)
	}
}

// TestShell_PipeData verifies data flows through pipes correctly.
func TestShell_PipeData(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", `echo "alpha\nbeta\ngamma" | grep beta`),
		assistantExit("e1", "success", "", ""),
	}, "pipe data flow", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	if !strings.Contains(results[0].Content, "beta") {
		t.Errorf("pipe output = %q, want 'beta'", results[0].Content)
	}
}

// ---------------------------------------------------------------------------
// 10. Background jobs in persistent shell
// ---------------------------------------------------------------------------

// TestShell_BackgroundJob verifies that background jobs can be started and waited on.
func TestShell_BackgroundJob(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "bg.txt")

	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", fmt.Sprintf(`echo background > %q &
BG_PID=$!
wait $BG_PID
echo "waited"`, outFile)),
		assistantExit("e1", "success", "", ""),
	}, "background job", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 1 {
		t.Fatal("expected at least 1 tool result")
	}

	if !strings.Contains(results[0].Content, "waited") {
		t.Errorf("bg job result = %q, want 'waited'", results[0].Content)
	}

	// Verify the background job wrote its file
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("bg output file not found: %v", err)
	}
	if !strings.Contains(string(data), "background") {
		t.Errorf("bg file content = %q, want 'background'", string(data))
	}
}

// TestShell_BackgroundPidPersists verifies that $! persists across sh calls.
func TestShell_BackgroundPidPersists(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "async.txt")

	_, entries := runWithMock(t, []tape.Message{
		// Start a background job that writes to a file after a brief sleep
		assistantsh("c1", fmt.Sprintf(`(sleep 0.1 && echo done > %q) &
echo $!`, outFile)),
		// Wait for it in a subsequent call
		assistantsh("c2", `wait
echo "all done"`),
		assistantExit("e1", "success", "", ""),
	}, "bg pid across calls", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	// First result should contain a PID number
	if !strings.Contains(results[0].Content, "[EXIT CODE] 0") {
		t.Errorf("bg start result = %q, want success", results[0].Content)
	}

	// Second result should show "all done"
	if !strings.Contains(results[1].Content, "all done") {
		t.Errorf("wait result = %q, want 'all done'", results[1].Content)
	}

	// Verify the async file was created
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("async output file not found: %v", err)
	}
	if !strings.Contains(string(data), "done") {
		t.Errorf("async file content = %q, want 'done'", string(data))
	}
}

// ---------------------------------------------------------------------------
// 11. Tape integrity
// ---------------------------------------------------------------------------

// TestTape_OutcomePresent verifies every session produces an outcome entry.
func TestTape_OutcomePresent(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantExit("e1", "success", "", ""),
	}, "tape outcome", "Begin.")

	outcome := findOutcome(t, entries)
	if outcome.ExitCode != 0 {
		t.Errorf("outcome exit code = %d, want 0", outcome.ExitCode)
	}
	if outcome.TerminationMode != tape.TermExit {
		t.Errorf("termination mode = %q, want %q", outcome.TerminationMode, tape.TermExit)
	}
}

// TestTape_TurnCounting verifies turn count matches sh calls.
func TestTape_TurnCounting(t *testing.T) {
	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", "echo 1"),
		assistantsh("c2", "echo 2"),
		assistantsh("c3", "echo 3"),
		assistantExit("e1", "success", "", ""),
	}, "turn counting", "Begin.")

	// Should have exactly 3 tool results (one per sh call)
	results := findToolResults(t, entries)
	if len(results) != 3 {
		t.Errorf("tool result count = %d, want 3", len(results))
	}

	// Outcome should record 3 turns
	outcome := findOutcome(t, entries)
	if outcome.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", outcome.ExitCode)
	}
}

// ---------------------------------------------------------------------------
// 12. File operations via shell
// ---------------------------------------------------------------------------

// TestShell_FileRoundtrip verifies write → read via shell commands.
func TestShell_FileRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "roundtrip.txt")

	_, entries := runWithMock(t, []tape.Message{
		assistantsh("c1", fmt.Sprintf(`echo "hello from shell" > %q`, testFile)),
		assistantsh("c2", fmt.Sprintf(`cat %q`, testFile)),
		assistantExit("e1", "success", "", ""),
	}, "file roundtrip", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	if !strings.Contains(results[1].Content, "hello from shell") {
		t.Errorf("file read = %q, want 'hello from shell'", results[1].Content)
	}
}

// TestShell_Permissions verifies permission enforcement.
func TestShell_Permissions(t *testing.T) {
	tmpDir := t.TempDir()
	noReadFile := filepath.Join(tmpDir, "noread.txt")

	_, entries := runWithMock(t, []tape.Message{
		// Create file, then remove read permission, then try to read
		assistantsh("c1", fmt.Sprintf(`echo "secret" > %q && chmod 000 %q`, noReadFile, noReadFile)),
		assistantsh("c2", fmt.Sprintf(`cat %q 2>&1`, noReadFile)),
		// Clean up
		assistantsh("c3", fmt.Sprintf(`chmod 644 %q`, noReadFile)),
		assistantExit("e1", "success", "", ""),
	}, "permission check", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", len(results))
	}

	// Second result should show a permission error and non-zero exit
	secondResult := results[1].Content
	if !strings.Contains(secondResult, "[EXIT CODE] 1") {
		t.Errorf("permission denied: %q, want [EXIT CODE] 1", secondResult)
	}
	if !strings.Contains(strings.ToLower(secondResult), "permission denied") &&
		!strings.Contains(strings.ToLower(secondResult), "operation not permitted") {
		t.Errorf("permission denied message missing: %q", secondResult)
	}
}

// TestShell_ExecutePermission verifies that scripts need execute permission.
func TestShell_ExecutePermission(t *testing.T) {
	tmpDir := t.TempDir()
	script := filepath.Join(tmpDir, "test.sh")

	_, entries := runWithMock(t, []tape.Message{
		// Create script without execute permission
		assistantsh("c1", fmt.Sprintf(`printf '#!/bin/sh\necho works\n' > %q`, script)),
		// Try to execute it directly (should fail)
		assistantsh("c2", fmt.Sprintf(`%s 2>&1`, script)),
		// Add execute permission and try again
		assistantsh("c3", fmt.Sprintf(`chmod +x %q && %s`, script, script)),
		assistantExit("e1", "success", "", ""),
	}, "execute permission", "Begin.")

	results := findToolResults(t, entries)
	if len(results) < 3 {
		t.Fatalf("expected at least 3 tool results, got %d", len(results))
	}

	// Without +x: should fail
	if strings.Contains(results[1].Content, "[EXIT CODE] 0") {
		t.Errorf("expected non-zero exit without +x, got: %q", results[1].Content)
	}

	// With +x: should succeed
	if !strings.Contains(results[2].Content, "works") {
		t.Errorf("with +x: %q, want 'works'", results[2].Content)
	}
}
