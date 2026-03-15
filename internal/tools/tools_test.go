package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
)

// testExecutor returns a ShExecutor with test-friendly defaults.
func testExecutor() *ShExecutor {
	dataDir, err := os.MkdirTemp("", "quine-tools-test-*")
	if err != nil {
		panic(err)
	}
	return &ShExecutor{
		Shell:     "/bin/sh",
		MaxOutput: 20480,
		Timeout:   30 * time.Second,
		DataDir:   dataDir,
		SessionID: fmt.Sprintf("test-%d", time.Now().UnixNano()),
	}
}

func requireWorkspaceSupport(t *testing.T, b *ShExecutor) {
	t.Helper()
	if err := b.Prepare(); err != nil {
		t.Skipf("workspace physics unsupported in this environment: %v", err)
	}
}

// --- Basic execution tests ---

func TestSimpleCommandExecution(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-1", "echo hello", 0, 0, false, false, "")

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
	if strings.Contains(result.Content, "[FS MUTATIONS]") {
		t.Errorf("did not expect FS mutations block without sandbox, got:\n%s", result.Content)
	}
}

func TestNonZeroExitCode(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-2", "false", 0, 0, false, false, "")

	if !result.IsError {
		t.Errorf("IsError = false, want true for non-zero exit")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 1") {
		t.Errorf("expected exit code 1, got:\n%s", result.Content)
	}
}

func TestNonZeroExitCodeFromCommand(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-2b", "exit 42", 0, 0, false, false, "")

	if !result.IsError {
		t.Errorf("IsError = false, want true for non-zero exit")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 42") {
		t.Errorf("expected exit code 42, got:\n%s", result.Content)
	}
}

func TestStderrCapture(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-3", "echo errormsg >&2", 0, 0, false, false, "")

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
	defer b.Close(false)
	b.MaxOutput = 100 // very small limit for testing

	// Generate output larger than MaxOutput (natural completion, no budget)
	result := b.Execute("tool-6", "python3 -c \"print('A' * 500)\"", 0, 0, false, false, "")

	if !strings.Contains(result.Content, "...[Output Truncated,") {
		t.Errorf("expected truncation notice, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "bytes total]") {
		t.Errorf("expected 'bytes total' in truncation notice, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "Increase QUINE_OUTPUT_TRUNCATE") {
		t.Errorf("expected QUINE_OUTPUT_TRUNCATE guidance in truncation notice, got:\n%s", result.Content)
	}
	if strings.Contains(result.Content, "job directory") {
		t.Errorf("truncation notice must not promise job directory recovery for sync calls, got:\n%s", result.Content)
	}
}

func TestOutputTruncationStderr(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	b.MaxOutput = 100

	result := b.Execute("tool-6b", "python3 -c \"import sys; sys.stderr.write('B' * 500)\"", 0, 0, false, false, "")

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
	defer b.Close(false)
	result := b.Execute("tool-7", "echo out; echo err >&2", 0, 0, false, false, "")

	// Non-detached mode should NOT include [JOB] header
	if strings.Contains(result.Content, "[JOB] pid=") {
		t.Fatalf("non-detached result should not have [JOB] header, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0") || !strings.Contains(result.Content, "[STDOUT]\nout\n") || !strings.Contains(result.Content, "[STDERR]\nerr\n") {
		t.Errorf("result format mismatch.\ngot:\n%q", result.Content)
	}
}

func TestResultFormatEmptyOutput(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-8", "true", 0, 0, false, false, "")

	// Non-detached mode should NOT include [JOB] header
	if strings.Contains(result.Content, "[JOB] pid=") {
		t.Fatalf("non-detached result should not have [JOB] header, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 0\n[STDOUT]\n[STDERR]\n") {
		t.Errorf("result format mismatch for empty output.\ngot:\n%q", result.Content)
	}
}

func TestOutputWithoutTrailingNewline(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("tool-nonl", `printf 'no-newline-here'`, 0, 0, false, false, "")

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
	defer b.Close(false)

	// Each sh call is ephemeral — exit 1 just gives exit code 1
	result := b.Execute("tool-exit-1", "exit 1", 0, 0, false, false, "")
	if !result.IsError {
		t.Errorf("expected error from exit 1")
	}
	if !strings.Contains(result.Content, "[EXIT CODE] 1") {
		t.Errorf("expected exit code 1, got:\n%s", result.Content)
	}

	// Subsequent calls still work (each is ephemeral)
	result2 := b.Execute("tool-exit-2", "echo alive", 0, 0, false, false, "")
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
	defer b.Close(false)
	cmd := fmt.Sprintf(`mkdir -p "$(dirname %q)" && printf '%%s\n' "hello world" > %q`, testFile, testFile)
	result := b.Execute("tool-9", cmd, 0, 0, false, false, "")

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
	defer b.Close(false)
	cmd := fmt.Sprintf(`cat -n %q`, testFile)
	result := b.Execute("tool-10", cmd, 0, 0, false, false, "")

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
	defer b.Close(false)
	cmd := fmt.Sprintf(`printf '%%s\n' "alpha beta gamma" > %q && cat -n %q`, testFile, testFile)
	result := b.Execute("tool-11", cmd, 0, 0, false, false, "")

	if result.IsError {
		t.Fatalf("roundtrip failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "alpha beta gamma") {
		t.Errorf("expected roundtrip content, got:\n%s", result.Content)
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
	defer b.Close(false)
	b.Env = MergeEnv(os.Environ(), []string{
		"QUINE_DEPTH=3",
	})

	result := b.Execute("tool-env-1", "echo $QUINE_DEPTH", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "3") {
		t.Errorf("expected QUINE_DEPTH=3 in output, got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-env-2", "echo \"SESSION_ID=${QUINE_SESSION_ID:-unset}\"", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "SESSION_ID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset in sh env, got:\n%s", result2.Content)
	}

	result3 := b.Execute("tool-env-3", "echo \"TAPE_ID=${QUINE_TAPE_ID:-unset}\"", 0, 0, false, false, "")
	if result3.IsError {
		t.Fatalf("command failed:\n%s", result3.Content)
	}
	if !strings.Contains(result3.Content, "TAPE_ID=unset") {
		t.Errorf("expected QUINE_TAPE_ID to be unset in sh env, got:\n%s", result3.Content)
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
	defer b.Close(false)

	result := b.Execute("tool-depth", "echo $QUINE_DEPTH", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "3") {
		t.Errorf("expected QUINE_DEPTH=3 (parent depth 2 + 1), got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-parent", "echo $QUINE_PARENT_SESSION", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "parent-session-id") {
		t.Errorf("expected QUINE_PARENT_SESSION=parent-session-id, got:\n%s", result2.Content)
	}

	result3 := b.Execute("tool-session", "echo \"SID=${QUINE_SESSION_ID:-unset}\"", 0, 0, false, false, "")
	if result3.IsError {
		t.Fatalf("command failed:\n%s", result3.Content)
	}
	if !strings.Contains(result3.Content, "SID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset, got:\n%s", result3.Content)
	}

	result4 := b.Execute("tool-tape", "echo \"TID=${QUINE_TAPE_ID:-unset}\"", 0, 0, false, false, "")
	if result4.IsError {
		t.Fatalf("command failed:\n%s", result4.Content)
	}
	if !strings.Contains(result4.Content, "TID=unset") {
		t.Errorf("expected QUINE_TAPE_ID to be unset, got:\n%s", result4.Content)
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
	defer b.Close(false)

	result := b.Execute("tool-ctor", "echo $QUINE_DEPTH", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "2") {
		t.Errorf("expected QUINE_DEPTH=2, got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-path", "which echo", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("'which echo' failed — PATH not propagated:\n%s", result2.Content)
	}
}

func TestWorkspaceOverlayReportsMutationsAndCommitsOnSuccess(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("realpath root: %v", err)
	}
	dataDir := t.TempDir()
	cfg := &config.Config{
		ModelID:          "claude-sonnet-4-20250514",
		APIKey:           "test-key",
		APIBase:          "https://api.example.com",
		Provider:         "anthropic",
		SessionID:        "sandbox-success",
		OutputTruncate:   20480,
		ShTimeout:        10,
		DataDir:          dataDir,
		Shell:            "/bin/sh",
		WorkspaceEnabled: true,
		WorkspaceRoot:    root,
		Workspace:        root,
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession: "workspace-session-success",
		WorkspaceOwner:   true,
	}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-sandbox-success", "printf 'hello overlay\\n' > result.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected sandbox error:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "[FS MUTATIONS]") {
		t.Fatalf("expected FS mutations block, got:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "+ result.txt (created)") {
		t.Fatalf("expected created workspace mutation for result.txt, got:\n%s", result.Content)
	}
	if _, err := os.Stat(filepath.Join(root, "result.txt")); !os.IsNotExist(err) {
		t.Fatalf("host file should not exist before commit, err=%v", err)
	}
	if err := b.Close(true); err != nil {
		t.Fatalf("commit close failed: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(rootReal, "result.txt"))
	if err != nil {
		t.Fatalf("read committed file: %v", err)
	}
	if strings.TrimSpace(string(data)) != "hello overlay" {
		t.Fatalf("committed file = %q, want hello overlay", strings.TrimSpace(string(data)))
	}
}

func TestWorkspaceOverlayRollsBackOnFailure(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{
		ModelID:          "claude-sonnet-4-20250514",
		APIKey:           "test-key",
		APIBase:          "https://api.example.com",
		Provider:         "anthropic",
		SessionID:        "sandbox-failure",
		OutputTruncate:   20480,
		ShTimeout:        10,
		DataDir:          t.TempDir(),
		Shell:            "/bin/sh",
		WorkspaceEnabled: true,
		WorkspaceRoot:    root,
		Workspace:        root,
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession: "workspace-session-failure",
		WorkspaceOwner:   true,
	}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-sandbox-failure", "printf 'temporary\\n' > doomed.txt; exit 1", 0, 0, false, false, "")
	if !result.IsError {
		t.Fatalf("expected non-zero exit to be an error")
	}
	if err := b.Close(false); err != nil {
		t.Fatalf("rollback close failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "doomed.txt")); !os.IsNotExist(err) {
		t.Fatalf("host file should have been rolled back, err=%v", err)
	}
}

func TestWorkspaceOverlayTracksAbsolutePaths(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{
		ModelID:          "claude-sonnet-4-20250514",
		APIKey:           "test-key",
		APIBase:          "https://api.example.com",
		Provider:         "anthropic",
		SessionID:        "sandbox-absolute",
		OutputTruncate:   20480,
		ShTimeout:        10,
		DataDir:          t.TempDir(),
		Shell:            "/bin/sh",
		WorkspaceEnabled: true,
		WorkspaceRoot:    root,
		Workspace:        root,
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession: "workspace-session-absolute",
		WorkspaceOwner:   true,
	}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	absFile := filepath.Join(root, "absolute.txt")
	result := b.Execute("tool-absolute", fmt.Sprintf("cd %q && printf 'abs\\n' > absolute.txt", root), 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected absolute-path error:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "~ absolute.txt (modified)") && !strings.Contains(result.Content, "+ absolute.txt (created)") {
		t.Fatalf("expected absolute-path mutation for absolute.txt, got:\n%s", result.Content)
	}
	if _, err := os.Stat(absFile); !os.IsNotExist(err) {
		t.Fatalf("host file should not exist before commit, err=%v", err)
	}
	if err := b.Close(true); err != nil {
		t.Fatalf("commit close failed: %v", err)
	}
	data, err := os.ReadFile(absFile)
	if err != nil {
		t.Fatalf("absolute-path file should be committed: %v", err)
	}
	if strings.TrimSpace(string(data)) != "abs" {
		t.Fatalf("absolute-path file = %q, want abs", strings.TrimSpace(string(data)))
	}
}

func TestWorkspaceDirectReportsMutationsAndCommitsOnSuccess(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{
		ModelID:          "claude-sonnet-4-20250514",
		APIKey:           "test-key",
		APIBase:          "https://api.example.com",
		Provider:         "anthropic",
		SessionID:        "workspace-direct-success",
		OutputTruncate:   20480,
		ShTimeout:        10,
		DataDir:          t.TempDir(),
		Shell:            "/bin/sh",
		WorkspaceEnabled: true,
		WorkspaceRoot:    root,
		Workspace:        root,
		WorkspaceBackend: "direct",
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession: "workspace-direct-session-success",
		WorkspaceOwner:   true,
	}

	b := NewShExecutor(cfg, nil)
	if err := b.Prepare(); err != nil {
		t.Fatalf("Prepare() error: %v", err)
	}
	result := b.Execute("tool-workspace-direct-success", "printf 'hello direct\\n' > result.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected direct workspace error:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "+ result.txt (created)") {
		t.Fatalf("expected created workspace mutation for result.txt, got:\n%s", result.Content)
	}
	data, err := os.ReadFile(filepath.Join(root, "result.txt"))
	if err != nil {
		t.Fatalf("read direct workspace file: %v", err)
	}
	if strings.TrimSpace(string(data)) != "hello direct" {
		t.Fatalf("direct workspace file = %q, want hello direct", strings.TrimSpace(string(data)))
	}
	if err := b.Close(true); err != nil {
		t.Fatalf("commit close failed: %v", err)
	}
}

func TestWorkspaceDirectRollsBackOnFailure(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{
		ModelID:          "claude-sonnet-4-20250514",
		APIKey:           "test-key",
		APIBase:          "https://api.example.com",
		Provider:         "anthropic",
		SessionID:        "workspace-direct-failure",
		OutputTruncate:   20480,
		ShTimeout:        10,
		DataDir:          t.TempDir(),
		Shell:            "/bin/sh",
		WorkspaceEnabled: true,
		WorkspaceRoot:    root,
		Workspace:        root,
		WorkspaceBackend: "direct",
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession: "workspace-direct-session-failure",
		WorkspaceOwner:   true,
	}

	b := NewShExecutor(cfg, nil)
	if err := b.Prepare(); err != nil {
		t.Fatalf("Prepare() error: %v", err)
	}
	result := b.Execute("tool-workspace-direct-failure", "printf 'temporary\\n' > doomed.txt; exit 1", 0, 0, false, false, "")
	if !result.IsError {
		t.Fatalf("expected non-zero exit to be an error")
	}
	if _, err := os.Stat(filepath.Join(root, "doomed.txt")); err != nil {
		t.Fatalf("expected provisional file before rollback, err=%v", err)
	}
	if err := b.Close(false); err != nil {
		t.Fatalf("rollback close failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "doomed.txt")); !os.IsNotExist(err) {
		t.Fatalf("direct workspace file should have been rolled back, err=%v", err)
	}
}

func TestWorkspaceDirectWorldRevisionsAndRestoreWorld(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{
		ModelID:          "claude-sonnet-4-20250514",
		APIKey:           "test-key",
		APIBase:          "https://api.example.com",
		Provider:         "anthropic",
		SessionID:        "workspace-direct-checkpoints",
		OutputTruncate:   20480,
		ShTimeout:        10,
		DataDir:          t.TempDir(),
		Shell:            "/bin/sh",
		WorkspaceEnabled: true,
		WorkspaceRoot:    root,
		Workspace:        root,
		WorkspaceBackend: "direct",
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession: "workspace-direct-session-checkpoints",
		WorkspaceOwner:   true,
	}

	b := NewShExecutor(cfg, nil)
	if err := b.Prepare(); err != nil {
		t.Fatalf("Prepare() error: %v", err)
	}

	b.TurnID = 1
	result1 := b.Execute("tool-checkpoint-1", "printf 'v1\\n' > state.txt", 0, 0, false, false, "")
	if result1.IsError {
		t.Fatalf("turn 1 failed:\n%s", result1.Content)
	}
	if !strings.Contains(result1.Content, "[WORLD REVISION] created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("turn 1 missing world revision block:\n%s", result1.Content)
	}

	b.TurnID = 2
	result2 := b.Execute("tool-checkpoint-2", "printf 'v2\\n' > state.txt", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("turn 2 failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "[WORLD REVISION] created=wr2 parent=wr1 current=wr2") {
		t.Fatalf("turn 2 missing world revision block:\n%s", result2.Content)
	}

	data, err := os.ReadFile(filepath.Join(root, "state.txt"))
	if err != nil {
		t.Fatalf("read current state: %v", err)
	}
	if strings.TrimSpace(string(data)) != "v2" {
		t.Fatalf("current state = %q, want v2", strings.TrimSpace(string(data)))
	}

	restore := b.RestoreWorld("tool-restore-world", "wr1")
	if restore.IsError {
		t.Fatalf("restore failed:\n%s", restore.Content)
	}
	if !strings.Contains(restore.Content, "[WORLD REVISION] wr2 -> wr1") {
		t.Fatalf("restore missing world revision transition:\n%s", restore.Content)
	}

	data, err = os.ReadFile(filepath.Join(root, "state.txt"))
	if err != nil {
		t.Fatalf("read restored state: %v", err)
	}
	if strings.TrimSpace(string(data)) != "v1" {
		t.Fatalf("restored state = %q, want v1", strings.TrimSpace(string(data)))
	}
}

func TestWorkspaceNoopShellDoesNotAdvanceWorldRevision(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{
		ModelID:               "claude-sonnet-4-20250514",
		APIKey:                "test-key",
		APIBase:               "https://api.example.com",
		Provider:              "anthropic",
		SessionID:             "workspace-noop-revision",
		OutputTruncate:        20480,
		ShTimeout:             10,
		DataDir:               t.TempDir(),
		Shell:                 "/bin/sh",
		WorkspaceEnabled:      true,
		WorkspaceRoot:         root,
		Workspace:             root,
		WorkspaceBackend:      "direct",
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkspaceSession:      "workspace-noop-revision-session",
		WorkspaceOwner:        true,
	}

	b := NewShExecutor(cfg, nil)
	if err := b.Prepare(); err != nil {
		t.Fatalf("Prepare() error: %v", err)
	}

	b.TurnID = 1
	result1 := b.Execute("tool-noop-revision-1", "printf 'v1\\n' > state.txt", 0, 0, false, false, "")
	if result1.IsError {
		t.Fatalf("turn 1 failed:\n%s", result1.Content)
	}
	if !strings.Contains(result1.Content, "[WORLD REVISION] created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("turn 1 missing revision creation:\n%s", result1.Content)
	}

	b.TurnID = 2
	result2 := b.Execute("tool-noop-revision-2", "test -f state.txt", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("turn 2 failed:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "[FS MUTATIONS]\n(empty)") {
		t.Fatalf("turn 2 should report empty mutations:\n%s", result2.Content)
	}
	if !strings.Contains(result2.Content, "[WORLD REVISION] current=wr1 (unchanged)") {
		t.Fatalf("turn 2 should keep the same revision:\n%s", result2.Content)
	}
	if strings.Contains(result2.Content, "created=wr2") {
		t.Fatalf("turn 2 should not create wr2 on a no-op:\n%s", result2.Content)
	}
	if got := b.CurrentWorldRevision(); got != "wr1" {
		t.Fatalf("CurrentWorldRevision() = %q, want wr1", got)
	}
}

func TestParseRestoreWorldArgs(t *testing.T) {
	req, err := ParseRestoreWorldArgs(map[string]any{"revision": "wr3"})
	if err != nil {
		t.Fatalf("ParseRestoreWorldArgs() error = %v", err)
	}
	if req.Revision != "wr3" {
		t.Fatalf("Revision = %q, want wr3", req.Revision)
	}

	if _, err := ParseRestoreWorldArgs(map[string]any{"revision": 3}); err == nil {
		t.Fatal("expected type error for numeric revision")
	}
	if _, err := ParseRestoreWorldArgs(map[string]any{"revision": "3"}); err == nil {
		t.Fatal("expected prefix validation error for revision without wr prefix")
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
		ParentSession:  "root-session",
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
	if envMap["QUINE_PARENT_SESSION"] != "root-session" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want root-session", envMap["QUINE_PARENT_SESSION"])
	}
	if envMap["QUINE_ORIGINAL_INTENT"] != originalIntent {
		t.Errorf("QUINE_ORIGINAL_INTENT = %q, want %q", envMap["QUINE_ORIGINAL_INTENT"], originalIntent)
	}
	if envMap["QUINE_SESSION_ID"] != "pre-exec-session" {
		t.Errorf("QUINE_SESSION_ID = %q, want pre-exec-session", envMap["QUINE_SESSION_ID"])
	}
	if envMap["QUINE_PARENT_SESSION"] != "root-session" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want root-session", envMap["QUINE_PARENT_SESSION"])
	}
	if envMap["QUINE_WISDOM_SUMMARY"] != "Found 3 bugs" {
		t.Errorf("QUINE_WISDOM_SUMMARY = %q, want 'Found 3 bugs'", envMap["QUINE_WISDOM_SUMMARY"])
	}
	if envMap["QUINE_WISDOM_PROGRESS"] != "50%" {
		t.Errorf("QUINE_WISDOM_PROGRESS = %q, want '50%%'", envMap["QUINE_WISDOM_PROGRESS"])
	}
	if _, exists := envMap["QUINE_TAPE_ID"]; exists {
		t.Errorf("QUINE_TAPE_ID should not be set in exec env")
	}
	if envMap["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want claude-sonnet-4-20250514", envMap["QUINE_MODEL_ID"])
	}
	if envMap["QUINE_MAX_DEPTH"] != "5" {
		t.Errorf("QUINE_MAX_DEPTH = %q, want 5", envMap["QUINE_MAX_DEPTH"])
	}
}

// --- helpers ---

func extractPID(content string) int {
	var id int
	fmt.Sscanf(content, "[JOB] pid=%d", &id)
	return id
}
