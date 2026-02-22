package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
)

// testExecutor returns a ShExecutor with test-friendly defaults.
func testExecutor() *ShExecutor {
	return &ShExecutor{
		Shell:     "/bin/sh",
		MaxOutput: 20480,
		Timeout:   30 * time.Second,
	}
}

// --- Basic execution tests ---

func TestSimpleCommandExecution(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	result := b.Execute("tool-1", "echo hello", 0, 0)

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
	defer b.Close()
	result := b.Execute("tool-2", "false", 0, 0)

	if !result.IsError {
		t.Errorf("IsError = false, want true for non-zero exit")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 1") {
		t.Errorf("expected exit code 1, got:\n%s", result.Content)
	}
}

func TestNonZeroExitCodeFromCommand(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	result := b.Execute("tool-2b", "exit 42", 0, 0)

	if !result.IsError {
		t.Errorf("IsError = false, want true for non-zero exit")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 42") {
		t.Errorf("expected exit code 42, got:\n%s", result.Content)
	}
}

func TestStderrCapture(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	result := b.Execute("tool-3", "echo errormsg >&2", 0, 0)

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

func TestOutputTruncation(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	b.MaxOutput = 100 // very small limit for testing

	// Generate output larger than MaxOutput (natural completion, no budget)
	result := b.Execute("tool-6", "python3 -c \"print('A' * 500)\"", 0, 0)

	if !strings.Contains(result.Content, "...[Output Truncated,") {
		t.Errorf("expected truncation notice, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "bytes total]") {
		t.Errorf("expected 'bytes total' in truncation notice, got:\n%s", result.Content)
	}
}

func TestOutputTruncationStderr(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	b.MaxOutput = 100

	result := b.Execute("tool-6b", "python3 -c \"import sys; sys.stderr.write('B' * 500)\"", 0, 0)

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
	defer b.Close()
	result := b.Execute("tool-7", "echo out; echo err >&2", 0, 0)

	expected := "[EXIT CODE] 0\n[STDOUT]\nout\n\n[STDERR]\nerr\n"
	if result.Content != expected {
		t.Errorf("result format mismatch.\ngot:\n%q\nwant:\n%q", result.Content, expected)
	}
}

func TestResultFormatEmptyOutput(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	result := b.Execute("tool-8", "true", 0, 0)

	expected := "[EXIT CODE] 0\n[STDOUT]\n\n[STDERR]\n"
	if result.Content != expected {
		t.Errorf("result format mismatch for empty output.\ngot:\n%q\nwant:\n%q", result.Content, expected)
	}
}

func TestOutputWithoutTrailingNewline(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	result := b.Execute("tool-nonl", `printf 'no-newline-here'`, 0, 0)

	if result.IsError {
		t.Fatalf("unexpected error:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0") {
		t.Errorf("expected exit code 0, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "no-newline-here") {
		t.Errorf("expected stdout to contain 'no-newline-here', got:\n%s", result.Content)
	}
}

func TestExitDoesNotCrash(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// Each sh call is ephemeral — exit 1 just gives exit code 1
	result := b.Execute("tool-exit-1", "exit 1", 0, 0)
	if !result.IsError {
		t.Errorf("expected error from exit 1")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 1") {
		t.Errorf("expected exit code 1, got:\n%s", result.Content)
	}

	// Subsequent calls still work (each is ephemeral)
	result2 := b.Execute("tool-exit-2", "echo alive", 0, 0)
	if result2.IsError {
		t.Fatalf("expected success after exit, got error:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "alive") {
		t.Errorf("expected 'alive' in output, got:\n%s", result2.Content)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "sub", "test.txt")

	b := testExecutor()
	defer b.Close()
	cmd := fmt.Sprintf(`mkdir -p "$(dirname %q)" && printf '%%s\n' "hello world" > %q`, testFile, testFile)
	result := b.Execute("tool-9", cmd, 0, 0)

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

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "read_test.txt")
	if err := os.WriteFile(testFile, []byte("line1\nline2\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	b := testExecutor()
	defer b.Close()
	cmd := fmt.Sprintf(`cat -n %q`, testFile)
	result := b.Execute("tool-10", cmd, 0, 0)

	if result.IsError {
		t.Fatalf("read_file failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "line1") {
		t.Errorf("expected 'line1' in output, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "line2") {
		t.Errorf("expected 'line2' in output, got:\n%s", result.Content)
	}
}

func TestWriteReadRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "roundtrip.txt")

	b := testExecutor()
	defer b.Close()
	cmd := fmt.Sprintf(`printf '%%s\n' "alpha beta gamma" > %q && cat -n %q`, testFile, testFile)
	result := b.Execute("tool-11", cmd, 0, 0)

	if result.IsError {
		t.Fatalf("roundtrip failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "alpha beta gamma") {
		t.Errorf("expected roundtrip content, got:\n%s", result.Content)
	}
}

// --- output_limit budget tests (Pain Principle) ---

func TestOutputLimitPausesJob(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// Generate lots of output with a small output_limit
	// seq 1 1000 produces ~4000 bytes, limit to 100
	result := b.Execute("tool-limit", "seq 1 1000", 0, 100)

	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Errorf("expected [PAUSED] with output_limit=100, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "job=") {
		t.Errorf("expected job= in [PAUSED] result, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "signal=\"cont\"") {
		t.Errorf("expected resume hint, got:\n%s", result.Content)
	}

	// Clean up: kill the paused job
	pgid := extractJobID(result.Content)
	if pgid > 0 {
		b.HandleJob("cleanup", map[string]any{"id": float64(pgid), "signal": "kill"})
	}
}

func TestJobKill(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// Start a job with a small limit so it pauses
	result := b.Execute("tool-killtest", "seq 1 10000", 0, 50)
	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Skipf("job completed before pause (output too small?): %s", result.Content)
	}

	pgid := extractJobID(result.Content)
	if pgid <= 0 {
		t.Fatalf("could not extract job id from: %s", result.Content)
	}

	// Kill the job
	killResult := b.HandleJob("tool-kill", map[string]any{"id": float64(pgid), "signal": "kill"})
	if killResult.IsError {
		t.Errorf("kill failed: %s", killResult.Content)
	}
	if !strings.Contains(killResult.Content, "killed") {
		t.Errorf("expected 'killed' in result, got: %s", killResult.Content)
	}

	// Job should no longer be in the registry
	j := b.jobs.Get(pgid)
	if j != nil {
		t.Errorf("job %d still in registry after kill", pgid)
	}
}

func TestJobReadOutput(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	result := b.Execute("tool-readtest", "seq 1 5000", 0, 50)
	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Skipf("job completed before pause")
	}

	pgid := extractJobID(result.Content)
	if pgid <= 0 {
		t.Fatalf("could not extract job id from: %s", result.Content)
	}

	// Read without resuming
	readResult := b.HandleJob("tool-read", map[string]any{"id": float64(pgid)})
	if readResult.IsError {
		t.Errorf("read failed: %s", readResult.Content)
	}
	if !strings.Contains(readResult.Content, "[JOB") {
		t.Errorf("expected [JOB ...] header, got: %s", readResult.Content)
	}

	// Clean up
	b.HandleJob("cleanup", map[string]any{"id": float64(pgid), "signal": "kill"})
}

func TestJobResume(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// Pause at 100 bytes
	result := b.Execute("tool-resume", "seq 1 10000", 0, 100)
	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Skipf("job completed before pause")
	}

	pgid := extractJobID(result.Content)
	if pgid <= 0 {
		t.Fatalf("could not extract job id")
	}

	// Resume with a larger limit
	resumeResult := b.HandleJob("tool-cont", map[string]any{
		"id":           float64(pgid),
		"signal":       "cont",
		"output_limit": float64(10000),
	})

	// Should either pause again or complete — not error
	if resumeResult.IsError {
		t.Errorf("resume failed: %s", resumeResult.Content)
	}

	// If still paused, kill it
	if strings.Contains(resumeResult.Content, "[PAUSED]") {
		pgid2 := extractJobID(resumeResult.Content)
		if pgid2 > 0 {
			b.HandleJob("cleanup", map[string]any{"id": float64(pgid2), "signal": "kill"})
		}
	}
}

func TestTimeoutPausesJob(t *testing.T) {
	b := testExecutor()
	defer b.Close()

	// Run a sleep command with 1 second timeout
	start := time.Now()
	result := b.Execute("tool-timeout", "sleep 10", 1*time.Second, 0)
	elapsed := time.Since(start)

	if !strings.Contains(result.Content, "[PAUSED]") {
		t.Errorf("expected [PAUSED] after timeout, got:\n%s", result.Content)
	}
	if elapsed > 3*time.Second {
		t.Errorf("took %v, expected ~1s for timeout", elapsed)
	}

	pgid := extractJobID(result.Content)
	if pgid > 0 {
		b.HandleJob("cleanup", map[string]any{"id": float64(pgid), "signal": "kill"})
	}
}

// --- Environment propagation tests ---

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

	envMap := make(map[string]string)
	for _, entry := range merged {
		key, val, _ := strings.Cut(entry, "=")
		envMap[key] = val
	}

	if envMap["PATH"] != "/usr/bin" {
		t.Errorf("PATH = %q, want /usr/bin", envMap["PATH"])
	}
	if envMap["HOME"] != "/home/user" {
		t.Errorf("HOME = %q, want /home/user", envMap["HOME"])
	}
	if envMap["QUINE_DEPTH"] != "1" {
		t.Errorf("QUINE_DEPTH = %q, want 1", envMap["QUINE_DEPTH"])
	}
	if envMap["QUINE_SESSION_ID"] != "child-session" {
		t.Errorf("QUINE_SESSION_ID = %q, want child-session", envMap["QUINE_SESSION_ID"])
	}
}

func TestEnvPropagationViaSh(t *testing.T) {
	b := testExecutor()
	defer b.Close()
	b.Env = MergeEnv(os.Environ(), []string{
		"QUINE_DEPTH=3",
	})

	result := b.Execute("tool-env-1", "echo $QUINE_DEPTH", 0, 0)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "3") {
		t.Errorf("expected QUINE_DEPTH=3 in output, got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-env-2", "echo \"SESSION_ID=${QUINE_SESSION_ID:-unset}\"", 0, 0)
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "SESSION_ID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset in sh env, got:\n%s", result2.Content)
	}
}

func TestChildEnvDepthIncrement(t *testing.T) {
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

	b := &ShExecutor{
		Shell:     "/bin/sh",
		MaxOutput: 20480,
		Timeout:   30 * time.Second,
		Env:       MergeEnv(os.Environ(), childEnv),
	}
	defer b.Close()

	result := b.Execute("tool-depth", "echo $QUINE_DEPTH", 0, 0)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "3") {
		t.Errorf("expected QUINE_DEPTH=3 (parent depth 2 + 1), got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-parent", "echo $QUINE_PARENT_SESSION", 0, 0)
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "parent-session-id") {
		t.Errorf("expected QUINE_PARENT_SESSION=parent-session-id, got:\n%s", result2.Content)
	}

	result3 := b.Execute("tool-session", "echo \"SID=${QUINE_SESSION_ID:-unset}\"", 0, 0)
	if result3.IsError {
		t.Fatalf("command failed:\n%s", result3.Content)
	}
	if !strings.Contains(result3.Content, "SID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset, got:\n%s", result3.Content)
	}
}

func TestNewShExecutorWithChildEnv(t *testing.T) {
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
	defer b.Close()

	result := b.Execute("tool-ctor", "echo $QUINE_DEPTH", 0, 0)
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "2") {
		t.Errorf("expected QUINE_DEPTH=2, got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-path", "which echo", 0, 0)
	if result2.IsError {
		t.Fatalf("'which echo' failed — PATH not propagated:\n%s", result2.Content)
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
			name: "persona wrong type",
			args: map[string]any{
				"persona": 123,
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
		Depth:          3,
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
	execEnv, err := cfg.ExecEnv(originalIntent)
	if err != nil {
		t.Fatalf("ExecEnv failed: %v", err)
	}

	envMap := make(map[string]string)
	for _, entry := range execEnv {
		key, val, _ := strings.Cut(entry, "=")
		envMap[key] = val
	}

	if envMap["QUINE_DEPTH"] != "0" {
		t.Errorf("QUINE_DEPTH = %q, want 0 (reset for exec)", envMap["QUINE_DEPTH"])
	}
	if envMap["QUINE_PARENT_SESSION"] != "pre-exec-session" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want pre-exec-session", envMap["QUINE_PARENT_SESSION"])
	}
	if envMap["QUINE_ORIGINAL_INTENT"] != originalIntent {
		t.Errorf("QUINE_ORIGINAL_INTENT = %q, want %q", envMap["QUINE_ORIGINAL_INTENT"], originalIntent)
	}
	if envMap["QUINE_WISDOM_SUMMARY"] != "Found 3 bugs" {
		t.Errorf("QUINE_WISDOM_SUMMARY = %q, want 'Found 3 bugs'", envMap["QUINE_WISDOM_SUMMARY"])
	}
	if envMap["QUINE_WISDOM_PROGRESS"] != "50%" {
		t.Errorf("QUINE_WISDOM_PROGRESS = %q, want '50%%'", envMap["QUINE_WISDOM_PROGRESS"])
	}
	if _, exists := envMap["QUINE_SESSION_ID"]; exists {
		t.Errorf("QUINE_SESSION_ID should not be set in exec env")
	}
	if envMap["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want claude-sonnet-4-20250514", envMap["QUINE_MODEL_ID"])
	}
	if envMap["QUINE_MAX_DEPTH"] != "5" {
		t.Errorf("QUINE_MAX_DEPTH = %q, want 5", envMap["QUINE_MAX_DEPTH"])
	}
}

// --- helpers ---

// extractJobID parses the job ID from a [PAUSED] result content string.
func extractJobID(content string) int {
	var id int
	fmt.Sscanf(content, "[PAUSED] job=%d", &id)
	return id
}
