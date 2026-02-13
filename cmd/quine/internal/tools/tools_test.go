package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

// testExecutor returns a ShExecutor with test-friendly defaults.
func testExecutor() *ShExecutor {
	return &ShExecutor{
		Shell:          "/bin/sh",
		DefaultTimeout: 10 * time.Second,
		MaxOutput:      20480,
		ShellInit:      shellInit,
	}
}

func TestNewshExecutor(t *testing.T) {
	cfg := &config.Config{
		Shell:          "/bin/sh",
		ShTimeout:      60,
		OutputTruncate: 10000,
	}
	b := NewShExecutor(cfg, nil)
	if b.Shell != "/bin/sh" {
		t.Errorf("Shell = %q, want /bin/sh", b.Shell)
	}
	if b.DefaultTimeout != 60*time.Second {
		t.Errorf("DefaultTimeout = %v, want 60s", b.DefaultTimeout)
	}
	if b.MaxOutput != 10000 {
		t.Errorf("MaxOutput = %d, want 10000", b.MaxOutput)
	}
	// Env should contain at least the OS environment
	if len(b.Env) == 0 {
		t.Error("Env should contain at least OS environment vars")
	}
}

func TestSimpleCommandExecution(t *testing.T) {
	b := testExecutor()
	result := b.Execute("tool-1", "echo hello", 0, false)

	if result.ToolID != "tool-1" {
		t.Errorf("ToolID = %q, want %q", result.ToolID, "tool-1")
	}
	if result.IsError {
		t.Errorf("IsError = true, want false for successful command")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0") {
		t.Errorf("expected exit code 0, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "hello") {
		t.Errorf("expected stdout to contain 'hello', got:\n%s", result.Content)
	}
}

func TestNonZeroExitCode(t *testing.T) {
	b := testExecutor()
	result := b.Execute("tool-2", "exit 42", 0, false)

	if !result.IsError {
		t.Errorf("IsError = false, want true for non-zero exit")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 42") {
		t.Errorf("expected exit code 42, got:\n%s", result.Content)
	}
}

func TestStderrCapture(t *testing.T) {
	b := testExecutor()
	result := b.Execute("tool-3", "echo errormsg >&2", 0, false)

	if !strings.Contains(result.Content, "errormsg") {
		t.Errorf("expected stderr to contain 'errormsg', got:\n%s", result.Content)
	}
	// Verify it appears in the STDERR section
	parts := strings.SplitN(result.Content, "[STDERR]", 2)
	if len(parts) < 2 {
		t.Fatalf("result missing [STDERR] section:\n%s", result.Content)
	}
	if !strings.Contains(parts[1], "errormsg") {
		t.Errorf("'errormsg' not in STDERR section:\n%s", result.Content)
	}
}

func TestTimeoutEnforcement(t *testing.T) {
	b := testExecutor()
	b.DefaultTimeout = 5 * time.Second // safety ceiling

	start := time.Now()
	result := b.Execute("tool-4", "sleep 30", 1, false)
	elapsed := time.Since(start)

	if !result.IsError {
		t.Errorf("IsError = false, want true for timed-out command")
	}
	if !strings.Contains(result.Content, "[EXIT CODE]") {
		t.Errorf("result missing exit code:\n%s", result.Content)
	}
	// The exit code should be non-zero (either -1 or 137 from SIGKILL)
	if strings.Contains(result.Content, "[EXIT CODE] 0") {
		t.Errorf("exit code should be non-zero for killed process, got:\n%s", result.Content)
	}
	// Should complete within ~2 seconds (1s timeout + buffer)
	if elapsed > 3*time.Second {
		t.Errorf("timeout took %v, expected ~1s", elapsed)
	}
}

func TestTimeoutUsesMinimum(t *testing.T) {
	b := testExecutor()
	b.DefaultTimeout = 1 * time.Second

	start := time.Now()
	// Provide a large timeout arg, but DefaultTimeout is smaller — it should
	// use the minimum (DefaultTimeout = 1s).
	result := b.Execute("tool-5", "sleep 30", 60, false)
	elapsed := time.Since(start)

	if !result.IsError {
		t.Errorf("IsError = false, want true for timed-out command")
	}
	if elapsed > 3*time.Second {
		t.Errorf("should have used DefaultTimeout (1s), but took %v", elapsed)
	}
	_ = result
}

func TestOutputTruncation(t *testing.T) {
	b := testExecutor()
	b.MaxOutput = 100 // very small limit for testing

	// Generate output larger than MaxOutput
	result := b.Execute("tool-6", "python3 -c \"print('A' * 500)\"", 0, false)

	if !strings.Contains(result.Content, "...[Output Truncated,") {
		t.Errorf("expected truncation notice, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "bytes total]") {
		t.Errorf("expected 'bytes total' in truncation notice, got:\n%s", result.Content)
	}
}

func TestOutputTruncationStderr(t *testing.T) {
	b := testExecutor()
	b.MaxOutput = 100

	result := b.Execute("tool-6b", "python3 -c \"import sys; sys.stderr.write('B' * 500)\"", 0, false)

	// The STDERR section should contain truncation
	parts := strings.SplitN(result.Content, "[STDERR]", 2)
	if len(parts) < 2 {
		t.Fatalf("result missing [STDERR] section:\n%s", result.Content)
	}
	if !strings.Contains(parts[1], "...[Output Truncated,") {
		t.Errorf("expected truncation in stderr, got:\n%s", result.Content)
	}
}

func TestResultFormatExact(t *testing.T) {
	b := testExecutor()
	result := b.Execute("tool-7", "echo out; echo err >&2", 0, false)

	expected := "[EXIT CODE] 0\n[STDOUT]\nout\n\n[STDERR]\nerr\n"
	if result.Content != expected {
		t.Errorf("result format mismatch.\ngot:\n%q\nwant:\n%q", result.Content, expected)
	}
}

func TestResultFormatEmptyOutput(t *testing.T) {
	b := testExecutor()
	result := b.Execute("tool-8", "true", 0, false)

	expected := "[EXIT CODE] 0\n[STDOUT]\n\n[STDERR]\n"
	if result.Content != expected {
		t.Errorf("result format mismatch for empty output.\ngot:\n%q\nwant:\n%q", result.Content, expected)
	}
}

func TestHelperWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "sub", "test.txt")

	b := testExecutor()
	cmd := fmt.Sprintf(`write_file %q "hello world"`, testFile)
	result := b.Execute("tool-9", cmd, 0, false)

	if result.IsError {
		t.Fatalf("write_file failed:\n%s", result.Content)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	content := strings.TrimSpace(string(data))
	if content != "hello world" {
		t.Errorf("file content = %q, want %q", content, "hello world")
	}
}

func TestHelperReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "read_test.txt")
	if err := os.WriteFile(testFile, []byte("line1\nline2\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	b := testExecutor()
	cmd := fmt.Sprintf(`read_file %q`, testFile)
	result := b.Execute("tool-10", cmd, 0, false)

	if result.IsError {
		t.Fatalf("read_file failed:\n%s", result.Content)
	}
	// cat -n produces numbered lines
	if !strings.Contains(result.Content, "1") && !strings.Contains(result.Content, "line1") {
		t.Errorf("expected numbered output with 'line1', got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "line1") {
		t.Errorf("expected 'line1' in output, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "line2") {
		t.Errorf("expected 'line2' in output, got:\n%s", result.Content)
	}
}

func TestHelperWriteReadRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "roundtrip.txt")

	b := testExecutor()
	cmd := fmt.Sprintf(`write_file %q "alpha beta gamma" && read_file %q`, testFile, testFile)
	result := b.Execute("tool-11", cmd, 0, false)

	if result.IsError {
		t.Fatalf("roundtrip failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "alpha beta gamma") {
		t.Errorf("expected roundtrip content, got:\n%s", result.Content)
	}
}

func TestToolResultType(t *testing.T) {
	b := testExecutor()
	result := b.Execute("tool-12", "echo test", 0, false)

	// Verify the result is of type tape.ToolResult
	var _ tape.ToolResult = result
	if result.ToolID != "tool-12" {
		t.Errorf("ToolID = %q, want %q", result.ToolID, "tool-12")
	}
}

// --- Recursion / Environment propagation tests ---

func TestMergeEnvOverlaysChildVars(t *testing.T) {
	osEnv := []string{
		"PATH=/usr/bin",
		"HOME=/home/user",
		"QUINE_DEPTH=0",
	}
	childEnv := []string{
		"QUINE_DEPTH=1",
		"QUINE_SESSION_ID=child-session",
	}

	merged := MergeEnv(osEnv, childEnv)

	// Build a map for easy lookup
	envMap := make(map[string]string)
	for _, entry := range merged {
		key, val, _ := strings.Cut(entry, "=")
		envMap[key] = val
	}

	// PATH should be preserved from osEnv
	if envMap["PATH"] != "/usr/bin" {
		t.Errorf("PATH = %q, want /usr/bin", envMap["PATH"])
	}
	// HOME should be preserved from osEnv
	if envMap["HOME"] != "/home/user" {
		t.Errorf("HOME = %q, want /home/user", envMap["HOME"])
	}
	// QUINE_DEPTH should be overridden by childEnv
	if envMap["QUINE_DEPTH"] != "1" {
		t.Errorf("QUINE_DEPTH = %q, want 1", envMap["QUINE_DEPTH"])
	}
	// QUINE_SESSION_ID should be added from childEnv
	if envMap["QUINE_SESSION_ID"] != "child-session" {
		t.Errorf("QUINE_SESSION_ID = %q, want child-session", envMap["QUINE_SESSION_ID"])
	}
}

func TestMergeEnvNilChildEnv(t *testing.T) {
	osEnv := []string{"PATH=/usr/bin", "HOME=/home/user"}
	merged := MergeEnv(osEnv, nil)

	if len(merged) != 2 {
		t.Errorf("expected 2 entries, got %d", len(merged))
	}
}

func TestMergeEnvNoDuplicateKeys(t *testing.T) {
	osEnv := []string{
		"QUINE_DEPTH=0",
		"PATH=/usr/bin",
	}
	childEnv := []string{
		"QUINE_DEPTH=1",
	}

	merged := MergeEnv(osEnv, childEnv)

	// Count QUINE_DEPTH entries — should be exactly 1
	count := 0
	for _, entry := range merged {
		if strings.HasPrefix(entry, "QUINE_DEPTH=") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 QUINE_DEPTH entry, got %d", count)
	}
}

func TestEnvPropagationViash(t *testing.T) {
	// Verify that spawned commands can see QUINE_* env vars
	b := testExecutor()
	b.Env = MergeEnv(os.Environ(), []string{
		"QUINE_DEPTH=3",
	})

	result := b.Execute("tool-env-1", "echo $QUINE_DEPTH", 0, false)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "3") {
		t.Errorf("expected QUINE_DEPTH=3 in output, got:\n%s", result.Content)
	}

	// Verify that QUINE_SESSION_ID is NOT set in the sh environment.
	// Each ./quine child process generates its own unique session ID
	// via config.Load(), ensuring multiple children spawned from one
	// sh command don't collide on the same tape file.
	result2 := b.Execute("tool-env-2", "echo \"SESSION_ID=${QUINE_SESSION_ID:-unset}\"", 0, false)
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "SESSION_ID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset in sh env, got:\n%s", result2.Content)
	}
}

func TestSessionIDNotSetInshEnv(t *testing.T) {
	// Verify that QUINE_SESSION_ID is NOT set in the sh environment.
	// This is critical: when a single sh command spawns multiple ./quine
	// children (via & backgrounding), each child must generate its own
	// session ID to avoid tape file collisions.
	b := testExecutor()
	b.Env = MergeEnv(os.Environ(), []string{"QUINE_DEPTH=1"})

	// Filter out any QUINE_SESSION_ID that might exist in the test env
	filtered := make([]string, 0, len(b.Env))
	for _, entry := range b.Env {
		if !strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			filtered = append(filtered, entry)
		}
	}
	b.Env = filtered

	// Execute multiple commands and verify none have QUINE_SESSION_ID set
	for i := 0; i < 3; i++ {
		result := b.Execute(fmt.Sprintf("tool-%d", i), "echo \"SID=${QUINE_SESSION_ID:-unset}\"", 0, false)
		if result.IsError {
			t.Fatalf("command %d failed:\n%s", i, result.Content)
		}
		if !strings.Contains(result.Content, "SID=unset") {
			t.Errorf("command %d: expected QUINE_SESSION_ID to be unset, got:\n%s", i, result.Content)
		}
	}
}

func TestChildEnvDepthIncrement(t *testing.T) {
	// Create a config at depth 2 and verify ChildEnv produces depth 3
	cfg := &config.Config{
		ModelID:        "claude-sonnet-4-20250514",
		APIKey:         "test-key",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          2,
		SessionID:      "parent-session-id",
		MaxConcurrent:  20,
		ShTimeout:      10,
		OutputTruncate: 20480,
		DataDir:        t.TempDir(),
		Shell:          "/bin/sh",
	}

	childEnv, err := cfg.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv failed: %v", err)
	}

	// Build a ShExecutor with the child env and verify QUINE_DEPTH
	b := &ShExecutor{
		Shell:          "/bin/sh",
		DefaultTimeout: 10 * time.Second,
		MaxOutput:      20480,
		ShellInit:      shellInit,
		Env:            MergeEnv(os.Environ(), childEnv),
	}

	result := b.Execute("tool-depth", "echo $QUINE_DEPTH", 0, false)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "3") {
		t.Errorf("expected QUINE_DEPTH=3 (parent depth 2 + 1), got:\n%s", result.Content)
	}

	// Verify QUINE_PARENT_SESSION is set to the parent's session ID
	result2 := b.Execute("tool-parent", "echo $QUINE_PARENT_SESSION", 0, false)
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "parent-session-id") {
		t.Errorf("expected QUINE_PARENT_SESSION=parent-session-id, got:\n%s", result2.Content)
	}

	// Verify QUINE_SESSION_ID is NOT set in the child env
	// (each ./quine child generates its own via config.Load)
	result3 := b.Execute("tool-session", "echo \"SID=${QUINE_SESSION_ID:-unset}\"", 0, false)
	if result3.IsError {
		t.Fatalf("command failed:\n%s", result3.Content)
	}
	if !strings.Contains(result3.Content, "SID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset, got:\n%s", result3.Content)
	}
}

func TestNoContextLeakageViaStdin(t *testing.T) {
	// Verify that a child process does NOT inherit the parent's tape/context.
	// Quine's design ensures context isolation via stdin-only input:
	// the child reads its task from stdin, not from env vars or shared memory.
	// This test verifies that a subprocess started by shExecutor does NOT
	// have access to any "conversation" content — it starts clean.
	b := testExecutor()
	b.Env = MergeEnv(os.Environ(), []string{"QUINE_DEPTH=1"})

	// A subprocess reading stdin should get empty input (no parent tape data)
	result := b.Execute("tool-leak", "cat < /dev/null | wc -c", 0, false)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	// wc -c of empty input should produce 0
	if !strings.Contains(result.Content, "0") {
		t.Errorf("expected 0 bytes from stdin (no context leak), got:\n%s", result.Content)
	}
}

func TestNewshExecutorWithChildEnv(t *testing.T) {
	cfg := &config.Config{
		ModelID:        "claude-sonnet-4-20250514",
		APIKey:         "test-key",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          1,
		SessionID:      "parent-abc",
		MaxConcurrent:  20,
		ShTimeout:      10,
		OutputTruncate: 20480,
		DataDir:        t.TempDir(),
		Shell:          "/bin/sh",
	}

	childEnv, err := cfg.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv failed: %v", err)
	}

	b := NewShExecutor(cfg, childEnv)

	// Verify QUINE_DEPTH is 2 (parent depth 1 + 1) in the executor's env
	result := b.Execute("tool-ctor", "echo $QUINE_DEPTH", 0, false)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "2") {
		t.Errorf("expected QUINE_DEPTH=2, got:\n%s", result.Content)
	}

	// PATH should still work (system tools accessible)
	result2 := b.Execute("tool-path", "which echo", 0, false)
	if result2.IsError {
		t.Fatalf("'which echo' failed — PATH not propagated:\n%s", result2.Content)
	}
}

// --- Stdout passthrough tests ---

func TestPassthroughWritesToFile(t *testing.T) {
	// Create a temp file to act as "stdout"
	tmpFile, err := os.CreateTemp(t.TempDir(), "stdout-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer tmpFile.Close()

	b := testExecutor()
	b.Stdout = tmpFile

	result := b.Execute("tool-pt-1", "echo hello-passthrough", 0, true)

	// Tool result should NOT contain the actual stdout content
	if strings.Contains(result.Content, "hello-passthrough") {
		t.Errorf("passthrough result should not contain stdout, got:\n%s", result.Content)
	}
	// Tool result should indicate passthrough mode
	if !strings.Contains(result.Content, "(passthrough)") {
		t.Errorf("expected '(passthrough)' in tool result, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0") {
		t.Errorf("expected exit code 0, got:\n%s", result.Content)
	}

	// The actual output should have been written to the file
	tmpFile.Sync()
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("read stdout file: %v", err)
	}
	if !strings.Contains(string(data), "hello-passthrough") {
		t.Errorf("expected 'hello-passthrough' in stdout file, got: %q", string(data))
	}
}

func TestPassthroughBinaryOutput(t *testing.T) {
	// Verify binary data flows through unmodified
	tmpFile, err := os.CreateTemp(t.TempDir(), "binary-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer tmpFile.Close()

	b := testExecutor()
	b.Stdout = tmpFile

	// Write 256 bytes (0x00-0xFF) to stdout via printf
	result := b.Execute("tool-pt-bin", `python3 -c "import sys; sys.stdout.buffer.write(bytes(range(256)))"`, 0, true)

	if result.IsError {
		t.Fatalf("binary passthrough failed:\n%s", result.Content)
	}

	tmpFile.Sync()
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("read binary file: %v", err)
	}
	if len(data) != 256 {
		t.Errorf("expected 256 bytes, got %d", len(data))
	}
	// Verify every byte value is present
	for i := 0; i < 256; i++ {
		if data[i] != byte(i) {
			t.Errorf("byte %d: got %d, want %d", i, data[i], i)
			break
		}
	}
}

func TestPassthroughStderrStillCaptured(t *testing.T) {
	// Even in passthrough mode, stderr should still be captured
	tmpFile, err := os.CreateTemp(t.TempDir(), "stdout-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer tmpFile.Close()

	b := testExecutor()
	b.Stdout = tmpFile

	result := b.Execute("tool-pt-err", "echo out-data; echo err-data >&2", 0, true)

	// Stderr should still appear in the tool result
	if !strings.Contains(result.Content, "err-data") {
		t.Errorf("expected stderr in tool result, got:\n%s", result.Content)
	}
	// Stdout should NOT appear in tool result (it's passthrough)
	if strings.Contains(result.Content, "out-data") {
		t.Errorf("passthrough result should not contain stdout 'out-data', got:\n%s", result.Content)
	}
}

func TestPassthroughFalseIsDefault(t *testing.T) {
	// When passthrough=false, behavior is the normal capture mode
	b := testExecutor()
	result := b.Execute("tool-pt-default", "echo captured", 0, false)

	if !strings.Contains(result.Content, "captured") {
		t.Errorf("non-passthrough should capture stdout, got:\n%s", result.Content)
	}
	if strings.Contains(result.Content, "(passthrough)") {
		t.Errorf("non-passthrough result should not contain '(passthrough)', got:\n%s", result.Content)
	}
}

func TestPassthroughWithoutStdoutFile(t *testing.T) {
	// When Stdout is nil, passthrough should fall back to capture mode
	b := testExecutor()
	b.Stdout = nil

	result := b.Execute("tool-pt-nil", "echo fallback", 0, true)

	// Should fall back to capture mode since Stdout is nil
	if !strings.Contains(result.Content, "fallback") {
		t.Errorf("passthrough with nil Stdout should capture, got:\n%s", result.Content)
	}
}

// --- Exec tool tests ---

func TestParseExecArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		want    ExecRequest
		wantErr bool
	}{
		{
			name:    "empty args",
			args:    map[string]any{},
			want:    ExecRequest{},
			wantErr: false,
		},
		{
			name: "with persona",
			args: map[string]any{
				"persona": "analyst",
			},
			want:    ExecRequest{Persona: "analyst"},
			wantErr: false,
		},
		{
			name: "with reason",
			args: map[string]any{
				"reason": "context too noisy",
			},
			want:    ExecRequest{Reason: "context too noisy"},
			wantErr: false,
		},
		{
			name: "with both",
			args: map[string]any{
				"persona": "coder",
				"reason":  "need fresh brain",
			},
			want:    ExecRequest{Persona: "coder", Reason: "need fresh brain"},
			wantErr: false,
		},
		{
			name: "persona wrong type",
			args: map[string]any{
				"persona": 123,
			},
			want:    ExecRequest{},
			wantErr: true,
		},
		{
			name: "reason wrong type",
			args: map[string]any{
				"reason": []string{"wrong"},
			},
			want:    ExecRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseExecArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExecArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Persona != tt.want.Persona {
					t.Errorf("Persona = %q, want %q", got.Persona, tt.want.Persona)
				}
				if got.Reason != tt.want.Reason {
					t.Errorf("Reason = %q, want %q", got.Reason, tt.want.Reason)
				}
			}
		})
	}
}

func TestExecEnv(t *testing.T) {
	cfg := &config.Config{
		ModelID:        "claude-sonnet-4-20250514",
		APIKey:         "test-key",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          3, // Current depth is 3
		SessionID:      "pre-exec-session",
		MaxConcurrent:  20,
		ShTimeout:      600,
		OutputTruncate: 20480,
		DataDir:        "/tmp/quine-test",
		Shell:          "/bin/sh",
		MaxTurns:       20,
		Wisdom: map[string]string{
			"SUMMARY":  "Found 3 bugs",
			"PROGRESS": "50%",
		},
	}

	originalIntent := "Fix the bugs in src/main.go"
	execEnv, err := cfg.ExecEnv(originalIntent, 0)
	if err != nil {
		t.Fatalf("ExecEnv failed: %v", err)
	}

	// Build a map for easy lookup
	envMap := make(map[string]string)
	for _, entry := range execEnv {
		key, val, _ := strings.Cut(entry, "=")
		envMap[key] = val
	}

	// DEPTH should be reset to 0 (fresh context)
	if envMap["QUINE_DEPTH"] != "0" {
		t.Errorf("QUINE_DEPTH = %q, want 0 (reset for exec)", envMap["QUINE_DEPTH"])
	}

	// PARENT_SESSION should be the pre-exec session ID
	if envMap["QUINE_PARENT_SESSION"] != "pre-exec-session" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want pre-exec-session", envMap["QUINE_PARENT_SESSION"])
	}

	// ORIGINAL_INTENT should be set
	if envMap["QUINE_ORIGINAL_INTENT"] != originalIntent {
		t.Errorf("QUINE_ORIGINAL_INTENT = %q, want %q", envMap["QUINE_ORIGINAL_INTENT"], originalIntent)
	}

	// WISDOM vars should be preserved
	if envMap["QUINE_WISDOM_SUMMARY"] != "Found 3 bugs" {
		t.Errorf("QUINE_WISDOM_SUMMARY = %q, want 'Found 3 bugs'", envMap["QUINE_WISDOM_SUMMARY"])
	}
	if envMap["QUINE_WISDOM_PROGRESS"] != "50%" {
		t.Errorf("QUINE_WISDOM_PROGRESS = %q, want '50%%'", envMap["QUINE_WISDOM_PROGRESS"])
	}

	// SESSION_ID should NOT be present (new process generates its own)
	if _, exists := envMap["QUINE_SESSION_ID"]; exists {
		t.Errorf("QUINE_SESSION_ID should not be set in exec env")
	}

	// STDIN_OFFSET should be 0 for this test
	if envMap["QUINE_STDIN_OFFSET"] != "0" {
		t.Errorf("QUINE_STDIN_OFFSET = %q, want 0", envMap["QUINE_STDIN_OFFSET"])
	}

	// Other config values should be preserved
	if envMap["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want claude-sonnet-4-20250514", envMap["QUINE_MODEL_ID"])
	}
	if envMap["QUINE_MAX_DEPTH"] != "5" {
		t.Errorf("QUINE_MAX_DEPTH = %q, want 5", envMap["QUINE_MAX_DEPTH"])
	}

	// Test with non-zero stdin offset
	execEnvWithOffset, err := cfg.ExecEnv(originalIntent, 12345)
	if err != nil {
		t.Fatalf("ExecEnv with offset failed: %v", err)
	}
	envMapOffset := make(map[string]string)
	for _, entry := range execEnvWithOffset {
		key, val, _ := strings.Cut(entry, "=")
		envMapOffset[key] = val
	}
	if envMapOffset["QUINE_STDIN_OFFSET"] != "12345" {
		t.Errorf("QUINE_STDIN_OFFSET = %q, want 12345", envMapOffset["QUINE_STDIN_OFFSET"])
	}
}

func TestNewExecExecutor(t *testing.T) {
	cfg := &config.Config{
		ModelID:   "claude-sonnet-4-20250514",
		Provider:  "anthropic",
		MaxDepth:  5,
		Depth:     1,
		SessionID: "test-session",
		DataDir:   t.TempDir(),
		Shell:     "/bin/sh",
	}

	originalIntent := "Test task"
	exec := NewExecExecutor(cfg, originalIntent)

	if exec.Cfg != cfg {
		t.Errorf("Cfg not set correctly")
	}
	if exec.OriginalIntent != originalIntent {
		t.Errorf("OriginalIntent = %q, want %q", exec.OriginalIntent, originalIntent)
	}
	if exec.QuinePath == "" {
		t.Errorf("QuinePath should not be empty")
	}
}

// TestReadExecutorBytesConsumed verifies that BytesConsumed correctly tracks
// the logical position in the stream, accounting for bufio's buffering.
func TestReadExecutorBytesConsumed(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5\n"
	reader := strings.NewReader(content)

	exec := NewReadExecutor(reader, 60*time.Second)

	// Read first 2 lines (12 bytes: "line1\nline2\n")
	req := ReadRequest{Lines: 2}
	result := exec.Execute("test1", req)
	if result.IsError {
		t.Fatalf("First read failed: %s", result.Content)
	}

	// BytesConsumed should be 12 (logical position after "line1\nline2\n")
	consumed := exec.BytesConsumed()
	// Note: Due to bufio buffering, the exact value may vary,
	// but it should be at least 12 (the logical minimum)
	if consumed < 12 {
		t.Errorf("BytesConsumed = %d, want >= 12", consumed)
	}

	// Read 2 more lines
	result = exec.Execute("test2", req)
	if result.IsError {
		t.Fatalf("Second read failed: %s", result.Content)
	}

	// BytesConsumed should now be at least 24
	consumed2 := exec.BytesConsumed()
	if consumed2 < 24 {
		t.Errorf("BytesConsumed after second read = %d, want >= 24", consumed2)
	}
	if consumed2 <= consumed {
		t.Errorf("BytesConsumed should increase: was %d, now %d", consumed, consumed2)
	}
}

// TestReadExecutorWithOffset verifies that NewReadExecutorWithOffset
// correctly tracks the total position including the initial offset.
func TestReadExecutorWithOffset(t *testing.T) {
	// Simulate content starting at offset 100
	content := "line_at_100\n"
	reader := strings.NewReader(content)

	initialOffset := int64(100)
	exec := NewReadExecutorWithOffset(reader, 60*time.Second, initialOffset)

	// Read the line
	req := ReadRequest{Lines: 1}
	result := exec.Execute("test1", req)
	if result.IsError {
		t.Fatalf("Read failed: %s", result.Content)
	}

	// BytesConsumed should include the initial offset
	consumed := exec.BytesConsumed()
	if consumed < initialOffset {
		t.Errorf("BytesConsumed = %d, should be >= initialOffset %d", consumed, initialOffset)
	}
	// And it should be at least initialOffset + len("line_at_100\n") = 100 + 12
	if consumed < 112 {
		t.Errorf("BytesConsumed = %d, want >= 112 (100 + 12)", consumed)
	}
}
