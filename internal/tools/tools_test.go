package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
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

func readOverlayWorkspaceFile(t *testing.T, s *subjectiveFS, rel string) string {
	t.Helper()
	exportDir := t.TempDir()
	if err := s.exportCurrentTree(exportDir); err != nil {
		t.Fatalf("exportCurrentTree() error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(exportDir, rel))
	if err != nil {
		t.Fatalf("read exported workspace file %q: %v", rel, err)
	}
	return strings.TrimSpace(string(data))
}

func decodeResultContent(t *testing.T, content any) map[string]any {
	t.Helper()
	var (
		payload map[string]any
		err     error
	)
	switch value := content.(type) {
	case string:
		payload, err = tape.UnmarshalToolResultContent(value)
	case json.RawMessage:
		payload, err = tape.UnmarshalStructuredToolResultContent(value)
	default:
		t.Fatalf("unsupported tool result content type %T", content)
	}
	if err != nil {
		t.Fatalf("decode tool result content: %v\ncontent=%v", err, content)
	}
	return payload
}

func resultString(t *testing.T, payload map[string]any, key string) string {
	t.Helper()
	value, _ := payload[key].(string)
	return value
}

func resultInt(t *testing.T, payload map[string]any, key string) int {
	t.Helper()
	switch value := payload[key].(type) {
	case json.Number:
		n, err := value.Int64()
		if err != nil {
			t.Fatalf("parse %q int: %v", key, err)
		}
		return int(n)
	case float64:
		return int(value)
	case int:
		return value
	case int64:
		return int(value)
	default:
		t.Fatalf("expected numeric %q in payload: %#v", key, payload)
		return 0
	}
}

func resultBool(t *testing.T, payload map[string]any, key string) bool {
	t.Helper()
	value, ok := payload[key].(bool)
	if !ok {
		t.Fatalf("expected boolean %q in payload: %#v", key, payload)
	}
	return value
}

func resultMap(t *testing.T, payload map[string]any, key string) map[string]any {
	t.Helper()
	value, _ := payload[key].(map[string]any)
	if value == nil {
		t.Fatalf("expected %q map in payload: %#v", key, payload)
	}
	return value
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
	payload := decodeResultContent(t, result.Content)
	if resultInt(t, payload, "exit_code") != 0 {
		t.Errorf("expected exit code 0, got:\n%s", result.Content)
	}
	if resultString(t, payload, "stdout") != "hello\n" {
		t.Errorf("expected stdout to contain 'hello', got:\n%s", result.Content)
	}
	if _, ok := payload["fs_mutations"]; ok {
		t.Errorf("did not expect fs_mutations without sandbox, got:\n%s", result.Content)
	}
}

func TestNonZeroExitCode(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-2", "false", 0, 0, false, false, "")

	if !result.IsError {
		t.Errorf("IsError = false, want true for non-zero exit")
	}
	if resultInt(t, decodeResultContent(t, result.Content), "exit_code") != 1 {
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
	if resultInt(t, decodeResultContent(t, result.Content), "exit_code") != 42 {
		t.Errorf("expected exit code 42, got:\n%s", result.Content)
	}
}

func TestStderrCapture(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-3", "echo errormsg >&2", 0, 0, false, false, "")

	payload := decodeResultContent(t, result.Content)
	if !strings.Contains(resultString(t, payload, "stderr"), "errormsg") {
		t.Errorf("expected stderr to contain 'errormsg', got:\n%s", result.Content)
	}
}

func TestOutputTruncation(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	b.MaxOutput = 100 // very small limit for testing

	// Generate output larger than MaxOutput (natural completion, no budget)
	result := b.Execute("tool-6", "python3 -c \"print('A' * 500)\"", 0, 0, false, false, "")
	payload := decodeResultContent(t, result.Content)
	stdout := resultString(t, payload, "stdout")

	if !strings.Contains(stdout, "...[Output Truncated,") {
		t.Errorf("expected truncation notice, got:\n%s", result.Content)
	}
	if !strings.Contains(stdout, "bytes total]") {
		t.Errorf("expected 'bytes total' in truncation notice, got:\n%s", result.Content)
	}
	if !strings.Contains(stdout, "Increase QUINE_OUTPUT_TRUNCATE") {
		t.Errorf("expected QUINE_OUTPUT_TRUNCATE guidance in truncation notice, got:\n%s", result.Content)
	}
	if strings.Contains(stdout, "job directory") {
		t.Errorf("truncation notice must not promise job directory recovery for sync calls, got:\n%s", result.Content)
	}
}

func TestOutputTruncationStderr(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	b.MaxOutput = 100

	result := b.Execute("tool-6b", "python3 -c \"import sys; sys.stderr.write('B' * 500)\"", 0, 0, false, false, "")
	payload := decodeResultContent(t, result.Content)

	if !strings.Contains(resultString(t, payload, "stderr"), "...[Output Truncated,") {
		t.Errorf("expected truncation in stderr, got:\n%s", result.Content)
	}
}

func TestResultFormatExact(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-7", "echo out; echo err >&2", 0, 0, false, false, "")

	payload := decodeResultContent(t, result.Content)
	if _, ok := payload["job"]; ok {
		t.Fatalf("non-detached result should not have job metadata, got:\n%s", result.Content)
	}
	if resultInt(t, payload, "exit_code") != 0 || resultString(t, payload, "stdout") != "out\n" || resultString(t, payload, "stderr") != "err\n" {
		t.Errorf("result format mismatch.\ngot:\n%q", result.Content)
	}
}

func TestResultContentIncludesExitStdoutStderr(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("tool-structured", "echo out; echo err >&2", 0, 0, false, false, "")
	if len(result.StructuredContent) != 0 {
		t.Fatal("expected structured content to be empty on sh result")
	}

	payload := decodeResultContent(t, result.Content)
	if payload["tool"] != "sh" {
		t.Fatalf("tool = %#v, want sh", payload["tool"])
	}
	if payload["mode"] != "sync" {
		t.Fatalf("mode = %#v, want sync", payload["mode"])
	}
	if payload["status"] != "completed" {
		t.Fatalf("status = %#v, want completed", payload["status"])
	}
	if payload["stdout"] != "out\n" {
		t.Fatalf("stdout = %#v, want %q", payload["stdout"], "out\n")
	}
	if payload["stderr"] != "err\n" {
		t.Fatalf("stderr = %#v, want %q", payload["stderr"], "err\n")
	}
}

func TestResultFormatEmptyOutput(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)
	result := b.Execute("tool-8", "true", 0, 0, false, false, "")

	payload := decodeResultContent(t, result.Content)
	if _, ok := payload["job"]; ok {
		t.Fatalf("non-detached result should not have job metadata, got:\n%s", result.Content)
	}
	if resultInt(t, payload, "exit_code") != 0 || resultString(t, payload, "stdout") != "" || resultString(t, payload, "stderr") != "" {
		t.Errorf("result format mismatch for empty output.\ngot:\n%q", result.Content)
	}
}

func TestStartJobPrecreatesMetadataSurface(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	job, err := b.startJob("sleep 10", false, "")
	if err != nil {
		t.Fatalf("startJob() error: %v", err)
	}
	t.Cleanup(func() {
		_ = syscall.Kill(-job.ID, syscall.SIGKILL)
		<-job.doneCh
	})

	cmdData, err := os.ReadFile(filepath.Join(job.canonicalDir, "cmd"))
	if err != nil {
		t.Fatalf("reading cmd file: %v", err)
	}
	if got := string(cmdData); got != "sleep 10" {
		t.Fatalf("cmd file = %q, want %q", got, "sleep 10")
	}
	assertDetachedJobSurface(t, job.canonicalDir, job.ID)
	if _, err := os.Stat(job.exitPath); !os.IsNotExist(err) {
		t.Fatalf("exit should be absent while job is still running; stat err=%v", err)
	}
}

func TestOutputWithoutTrailingNewline(t *testing.T) {
	b := testExecutor()
	defer b.Close(false)

	result := b.Execute("tool-nonl", `printf 'no-newline-here'`, 0, 0, false, false, "")

	if result.IsError {
		t.Fatalf("unexpected error:\n%s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	if resultInt(t, payload, "exit_code") != 0 {
		t.Errorf("expected exit code 0, got:\n%s", result.Content)
	}
	if resultString(t, payload, "stdout") != "no-newline-here" {
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
	if resultInt(t, decodeResultContent(t, result.Content), "exit_code") != 1 {
		t.Errorf("expected exit code 1, got:\n%s", result.Content)
	}

	// Subsequent calls still work (each is ephemeral)
	result2 := b.Execute("tool-exit-2", "echo alive", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("expected success after exit, got error:\n%s", result2.Content)
	}
	if !strings.Contains(string(result2.Content), "alive") {
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
	if !strings.Contains(string(result.Content), "line1") {
		t.Errorf("expected 'line1' in output, got:\n%s", result.Content)
	}
	if !strings.Contains(string(result.Content), "line2") {
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
	if !strings.Contains(string(result.Content), "alpha beta gamma") {
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
	if !strings.Contains(string(result.Content), "3") {
		t.Errorf("expected QUINE_DEPTH=3 in output, got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-env-2", "echo \"SESSION_ID=${QUINE_SESSION_ID:-unset}\"", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(string(result2.Content), "SESSION_ID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset in sh env, got:\n%s", result2.Content)
	}

	result3 := b.Execute("tool-env-3", "echo \"TAPE_ID=${QUINE_TAPE_ID:-unset}\"", 0, 0, false, false, "")
	if result3.IsError {
		t.Fatalf("command failed:\n%s", result3.Content)
	}
	if !strings.Contains(string(result3.Content), "TAPE_ID=unset") {
		t.Errorf("expected QUINE_TAPE_ID to be unset in sh env, got:\n%s", result3.Content)
	}

}

func TestChildEnvDepthIncrement(t *testing.T) {
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", Depth: 2, SessionID: "parent-session-id"}, Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"}, Limits: config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 10, OutputTruncate: 20480}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

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
	if !strings.Contains(string(result.Content), "3") {
		t.Errorf("expected QUINE_DEPTH=3 (parent depth 2 + 1), got:\n%s", result.Content)
	}

	result2 := b.Execute("tool-parent", "echo $QUINE_PARENT_SESSION", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("command failed:\n%s", result2.Content)
	}
	if !strings.Contains(string(result2.Content), "parent-session-id") {
		t.Errorf("expected QUINE_PARENT_SESSION=parent-session-id, got:\n%s", result2.Content)
	}

	result3 := b.Execute("tool-session", "echo \"SID=${QUINE_SESSION_ID:-unset}\"", 0, 0, false, false, "")
	if result3.IsError {
		t.Fatalf("command failed:\n%s", result3.Content)
	}
	if !strings.Contains(string(result3.Content), "SID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be unset, got:\n%s", result3.Content)
	}

	result4 := b.Execute("tool-tape", "echo \"TID=${QUINE_TAPE_ID:-unset}\"", 0, 0, false, false, "")
	if result4.IsError {
		t.Fatalf("command failed:\n%s", result4.Content)
	}
	if !strings.Contains(string(result4.Content), "TID=unset") {
		t.Errorf("expected QUINE_TAPE_ID to be unset, got:\n%s", result4.Content)
	}
}

func TestNewShExecutorWithChildEnv(t *testing.T) {
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", Depth: 1, SessionID: "parent-abc", RunID: "run-parent-abc"}, Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"}, Limits: config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 10, OutputTruncate: 20480}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

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
	if !strings.Contains(string(result.Content), "2") {
		t.Errorf("expected QUINE_DEPTH=2, got:\n%s", result.Content)
	}

	resultRun := b.Execute("tool-run", "echo \"RUN_ID=${QUINE_RUN_ID:-unset} SID=${QUINE_SESSION_ID:-unset}\"", 0, 0, false, false, "")
	if resultRun.IsError {
		t.Fatalf("command failed:\n%s", resultRun.Content)
	}
	if !strings.Contains(string(resultRun.Content), "RUN_ID=run-parent-abc") {
		t.Errorf("expected QUINE_RUN_ID in sh env, got:\n%s", resultRun.Content)
	}
	if !strings.Contains(string(resultRun.Content), "SID=unset") {
		t.Errorf("expected QUINE_SESSION_ID to be absent in sh env, got:\n%s", resultRun.Content)
	}

	result2 := b.Execute("tool-path", "which echo", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("'which echo' failed — PATH not propagated:\n%s", result2.Content)
	}
}

func TestNewShExecutorStripsContextBootstrapEnv(t *testing.T) {
	t.Setenv(ContextBootstrapEnv, "/tmp/bootstrap-context")

	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "parent-abc", RunID: "run-parent-abc"}, Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, []string{ContextBootstrapEnv + "=/tmp/child-bootstrap"})
	defer b.Close(false)

	result := b.Execute("tool-bootstrap", "echo \"CTX_BOOTSTRAP=${QUINE_CONTEXT_BOOTSTRAP:-unset}\"", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("command failed:\n%s", result.Content)
	}
	if !strings.Contains(string(result.Content), "CTX_BOOTSTRAP=unset") {
		t.Errorf("expected QUINE_CONTEXT_BOOTSTRAP to be unset in sh env, got:\n%s", result.Content)
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
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "sandbox-success"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-session-success", WorkspaceOwner: true}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-sandbox-success", "printf 'hello overlay\\n' > result.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected sandbox error:\n%s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	fsMutations := resultString(t, payload, "fs_mutations")
	if fsMutations == "" {
		t.Fatalf("expected FS mutations block, got:\n%s", result.Content)
	}
	if !strings.Contains(fsMutations, "+ result.txt (created)") {
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

func TestWorkspaceOverlayCanCommitWithoutKeepingDetachedJobs(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	dataDir := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "sandbox-signal-commit"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-session-signal-commit", WorkspaceOwner: true}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-signal-commit", "printf 'signal-safe\\n' > result.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected sandbox error:\n%s", result.Content)
	}
	if _, err := os.Stat(filepath.Join(root, "result.txt")); !os.IsNotExist(err) {
		t.Fatalf("host file should not exist before signal-style close, err=%v", err)
	}
	if err := b.CloseWithOptions(false, true); err != nil {
		t.Fatalf("signal-style commit close failed: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "result.txt"))
	if err != nil {
		t.Fatalf("read committed file: %v", err)
	}
	if strings.TrimSpace(string(data)) != "signal-safe" {
		t.Fatalf("committed file = %q, want signal-safe", strings.TrimSpace(string(data)))
	}
}

func TestWorkspaceOverlayCanHideFSMutationTelemetry(t *testing.T) {
	root := t.TempDir()
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		rootReal = root
	}
	dataDir := t.TempDir()
	gates := config.DefaultToolGates()
	gates.FSMutationTelemetry = false
	cfg := &config.Config{

		ToolGates: gates, Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "sandbox-hidden-fs-mutations"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-session-hidden-fs-mutations", WorkspaceOwner: true}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"},
	}

	b := NewShExecutor(cfg, nil)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-sandbox-hidden-fs", "printf 'hello overlay\\n' > result.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected sandbox error:\n%s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	if _, ok := payload["fs_mutations"]; ok {
		t.Fatalf("expected fs_mutations telemetry to be hidden, got:\n%s", result.Content)
	}
	if resultString(t, payload, "world_revision") == "" {
		t.Fatalf("expected world_revision to remain visible, got:\n%s", result.Content)
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

func TestWorkspaceDirectReportsMutationsAndPersistsOnFailure(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-direct-failure"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceBackend: "direct", WorkspaceRevisionMode: config.WorkspaceRevisionNone, WorkspaceSession: "workspace-direct-failure-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-direct-failure", "printf 'shared\\n' > doomed.txt; exit 1", 0, 0, false, false, "")
	if !result.IsError {
		t.Fatal("expected non-zero exit to be an error")
	}
	payload := decodeResultContent(t, result.Content)
	fsMutations := resultString(t, payload, "fs_mutations")
	if !strings.Contains(fsMutations, "+ doomed.txt (created)") {
		t.Fatalf("expected direct workspace mutation for doomed.txt, got:\n%s", result.Content)
	}
	if value := resultString(t, payload, "world_revision"); value != "" {
		t.Fatalf("direct backend should not report world_revision, got %q", value)
	}
	data, err := os.ReadFile(filepath.Join(root, "doomed.txt"))
	if err != nil {
		t.Fatalf("direct workspace file should persist immediately: %v", err)
	}
	if strings.TrimSpace(string(data)) != "shared" {
		t.Fatalf("direct workspace file = %q, want shared", strings.TrimSpace(string(data)))
	}
}

func TestWorkspaceDirectObservesPeerMutationOnNextShellBoundary(t *testing.T) {
	root := t.TempDir()
	dataDir := t.TempDir()

	newDirectCfg := func(sessionID string) *config.Config {
		return &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: sessionID}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceBackend: "direct", WorkspaceRevisionMode: config.WorkspaceRevisionNone, WorkspaceSession: "workspace-direct-shared-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"}}
	}

	observer := NewShExecutor(newDirectCfg("workspace-direct-observer"), nil)
	defer observer.Close(false)
	requireWorkspaceSupport(t, observer)
	peer := NewShExecutor(newDirectCfg("workspace-direct-peer"), nil)
	defer peer.Close(false)
	requireWorkspaceSupport(t, peer)

	peerResult := peer.Execute("tool-direct-peer-write", "printf 'peer\\n' > handoff.txt", 0, 0, false, false, "")
	if peerResult.IsError {
		t.Fatalf("unexpected peer write error:\n%s", peerResult.Content)
	}

	observeResult := observer.Execute("tool-direct-observer-sync", "test -f handoff.txt", 0, 0, false, false, "")
	if observeResult.IsError {
		t.Fatalf("observer sync should succeed:\n%s", observeResult.Content)
	}
	fsMutations := resultString(t, decodeResultContent(t, observeResult.Content), "fs_mutations")
	if !strings.Contains(fsMutations, "+ handoff.txt (created)") {
		t.Fatalf("expected observer to see peer-created file on next boundary, got:\n%s", observeResult.Content)
	}

	noopResult := observer.Execute("tool-direct-observer-noop", "test -f handoff.txt", 0, 0, false, false, "")
	if noopResult.IsError {
		t.Fatalf("observer noop should succeed:\n%s", noopResult.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, noopResult.Content), "fs_mutations"), "[FS MUTATIONS]\n(empty)") {
		t.Fatalf("expected second observer boundary to be empty:\n%s", noopResult.Content)
	}
}

func TestWorkspaceOverlayPreservesFailureAsRestorableRevision(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "sandbox-failure"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-session-failure", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-sandbox-failure", "printf 'temporary\\n' > doomed.txt; exit 1", 0, 0, false, false, "")
	if !result.IsError {
		t.Fatalf("expected non-zero exit to be an error")
	}
	payload := decodeResultContent(t, result.Content)
	if !strings.Contains(resultString(t, payload, "fs_mutations"), "+ doomed.txt (created)") {
		t.Fatalf("expected failed shell side effect to be reported:\n%s", result.Content)
	}
	if !strings.Contains(resultString(t, payload, "world_revision"), "[WORLD REVISION] created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("expected failed shell side effect to create a revision:\n%s", result.Content)
	}
	if got := readOverlayWorkspaceFile(t, b.subjective, "doomed.txt"); got != "temporary" {
		t.Fatalf("current failed-world file = %q, want temporary", got)
	}

	restore := b.SwitchWorld("tool-sandbox-failure-restore", "wr0")
	if restore.IsError {
		t.Fatalf("switch back to baseline failed:\n%s", restore.Content)
	}
	exportDir := t.TempDir()
	if err := b.subjective.exportCurrentTree(exportDir); err != nil {
		t.Fatalf("export restored workspace: %v", err)
	}
	if _, err := os.Stat(filepath.Join(exportDir, "doomed.txt")); !os.IsNotExist(err) {
		t.Fatalf("restored baseline should not contain doomed.txt, err=%v", err)
	}

	if err := b.Close(false); err != nil {
		t.Fatalf("rollback close failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "doomed.txt")); !os.IsNotExist(err) {
		t.Fatalf("host file should have been rolled back, err=%v", err)
	}
}

func TestWorkspaceOverlayTimeoutKillsShellAndPreservesBoundaryEffects(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "sandbox-timeout"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-session-timeout", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)

	start := time.Now()
	result := b.Execute("tool-overlay-timeout", "printf start; printf temporary > timed.txt; sleep 30; printf end", 150*time.Millisecond, 0, false, false, "")
	if result.IsError {
		t.Fatalf("expected interrupted timeout result, got error:\n%s", result.Content)
	}
	if elapsed := time.Since(start); elapsed >= 5*time.Second {
		t.Fatalf("overlay timeout path took too long: %v", elapsed)
	}

	payload := decodeResultContent(t, result.Content)
	if got := resultString(t, payload, "status"); got != "interrupted" {
		t.Fatalf("status = %q, want interrupted; content=%s", got, result.Content)
	}
	if got := resultString(t, payload, "cause"); got != "timeout" {
		t.Fatalf("cause = %q, want timeout", got)
	}
	if stdout := resultString(t, payload, "stdout_so_far"); stdout != "start" {
		t.Fatalf("stdout_so_far = %q, want start", stdout)
	}
	if !strings.Contains(resultString(t, payload, "fs_mutations_so_far"), "+ timed.txt (created)") {
		t.Fatalf("timeout should report workspace side effects that reached the boundary:\n%s", result.Content)
	}
	if !strings.Contains(resultString(t, payload, "world_revision"), "[WORLD REVISION] created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("timeout side effect should create a world revision:\n%s", result.Content)
	}

	observe := b.Execute("tool-overlay-timeout-observe", "test -e timed.txt && printf preserved", 5*time.Second, 0, false, false, "")
	if observe.IsError {
		t.Fatalf("expected next overlay turn to observe preserved timeout writes:\n%s", observe.Content)
	}
	observePayload := decodeResultContent(t, observe.Content)
	if stdout := resultString(t, observePayload, "stdout"); stdout != "preserved" {
		t.Fatalf("observe stdout = %q, want preserved", stdout)
	}

	restore := b.SwitchWorld("tool-overlay-timeout-restore", "wr0")
	if restore.IsError {
		t.Fatalf("switch back to baseline failed:\n%s", restore.Content)
	}
	clean := b.Execute("tool-overlay-timeout-clean", "test ! -e timed.txt && printf clean", 5*time.Second, 0, false, false, "")
	if clean.IsError {
		t.Fatalf("expected restored baseline to remove timeout writes:\n%s", clean.Content)
	}
	cleanPayload := decodeResultContent(t, clean.Content)
	if stdout := resultString(t, cleanPayload, "stdout"); stdout != "clean" {
		t.Fatalf("restored clean stdout = %q, want clean", stdout)
	}
}

func TestWorkspaceOverlayTracksAbsolutePaths(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "sandbox-absolute"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-session-absolute", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	absFile := filepath.Join(root, "absolute.txt")
	result := b.Execute("tool-absolute", fmt.Sprintf("cd %q && printf 'abs\\n' > absolute.txt", root), 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected absolute-path error:\n%s", result.Content)
	}
	fsMutations := resultString(t, decodeResultContent(t, result.Content), "fs_mutations")
	if !strings.Contains(fsMutations, "~ absolute.txt (modified)") && !strings.Contains(fsMutations, "+ absolute.txt (created)") {
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

func TestWorkspaceOverlayAllowsNarrowerWorkspaceAbsentFromHostRoot(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	subdir := filepath.Join(root, "subdir")
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-overlay-narrow-child"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: subdir, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-overlay-narrow-child-session", WorkspaceOwner: false}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-overlay-narrow-child", "printf 'child\\n' > child.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected overlay narrow child error:\n%s", result.Content)
	}
	fsMutations := resultString(t, decodeResultContent(t, result.Content), "fs_mutations")
	if !strings.Contains(fsMutations, "+ child.txt (created)") && !strings.Contains(fsMutations, "+ subdir/child.txt (created)") {
		t.Fatalf("expected created workspace mutation for child.txt, got:\n%s", result.Content)
	}
	if _, err := os.Stat(subdir); !os.IsNotExist(err) {
		t.Fatalf("host workspace subdir should not exist before owner commit, err=%v", err)
	}
}

func TestWorkspaceBootstrapClonesSourceLineageIntoFreshSession(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	dataDir := t.TempDir()
	sourceState := filepath.Join(dataDir, "workspaces", "parent-session")
	if err := os.MkdirAll(filepath.Join(sourceState, "base"), 0o755); err != nil {
		t.Fatalf("mkdir source base: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(sourceState, "layers", "wr4"), 0o755); err != nil {
		t.Fatalf("mkdir source layer: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceState, "STATE_VERSION"), []byte(overlayStateVersion+"\n"), 0o644); err != nil {
		t.Fatalf("write source state version: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceState, "base", "base.txt"), []byte("baseline\n"), 0o644); err != nil {
		t.Fatalf("write source base file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceState, "layers", "wr4", "child.txt"), []byte("snapshot wr4\n"), 0o644); err != nil {
		t.Fatalf("write source layer file: %v", err)
	}
	ledger := worldRevisionLedger{
		Current: "wr4",
		Next:    5,
		Revisions: map[string]worldRevision{
			"wr0": {ID: "wr0", Kind: "baseline"},
			"wr4": {ID: "wr4", Parent: "wr3", Kind: "sh", Turn: 7},
		},
	}
	ledgerData, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		t.Fatalf("marshal ledger: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceState, "world-revisions.json"), ledgerData, 0o644); err != nil {
		t.Fatalf("write source ledger: %v", err)
	}

	s := &subjectiveFS{
		enabled:          true,
		dataDir:          dataDir,
		workspaceSession: "child-session",
		bootstrapSource:  "parent-session",
	}
	if err := s.bootstrapWorkspaceState(); err != nil {
		t.Fatalf("bootstrapWorkspaceState() error: %v", err)
	}

	targetState := filepath.Join(dataDir, "workspaces", "child-session")
	if data, err := os.ReadFile(filepath.Join(targetState, "STATE_VERSION")); err != nil {
		t.Fatalf("read target state version: %v", err)
	} else if strings.TrimSpace(string(data)) != overlayStateVersion {
		t.Fatalf("target state version = %q, want %q", strings.TrimSpace(string(data)), overlayStateVersion)
	}
	if data, err := os.ReadFile(filepath.Join(targetState, "base", "base.txt")); err != nil {
		t.Fatalf("read target base file: %v", err)
	} else if strings.TrimSpace(string(data)) != "baseline" {
		t.Fatalf("target base file = %q, want baseline", strings.TrimSpace(string(data)))
	}
	if data, err := os.ReadFile(filepath.Join(targetState, "layers", "wr4", "child.txt")); err != nil {
		t.Fatalf("read target layer file: %v", err)
	} else if strings.TrimSpace(string(data)) != "snapshot wr4" {
		t.Fatalf("target layer file = %q, want wr4 snapshot copy", strings.TrimSpace(string(data)))
	}
	if _, err := os.Stat(filepath.Join(targetState, "live", "upper")); err != nil {
		t.Fatalf("expected bootstrapped live upperdir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(targetState, "live", "work")); err != nil {
		t.Fatalf("expected bootstrapped live workdir: %v", err)
	}
	clonedLedger, err := os.ReadFile(filepath.Join(targetState, "world-revisions.json"))
	if err != nil {
		t.Fatalf("read target ledger: %v", err)
	}
	if string(clonedLedger) != string(ledgerData) {
		t.Fatalf("target ledger mismatch\ngot:\n%s\nwant:\n%s", string(clonedLedger), string(ledgerData))
	}
}

func TestWorkspaceOverlaySwitchWorldReusesTargetHandle(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-overlay-restore"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-overlay-restore-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)

	b.TurnID = 1
	result1 := b.Execute("tool-checkpoint-1", "printf 'v1\\n' > state.txt", 0, 0, false, false, "")
	if result1.IsError {
		t.Fatalf("turn 1 failed:\n%s", result1.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, result1.Content), "world_revision"), "[WORLD REVISION] created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("turn 1 missing world revision block:\n%s", result1.Content)
	}

	b.TurnID = 2
	result2 := b.Execute("tool-checkpoint-2", "printf 'v2\\n' > state.txt", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("turn 2 failed:\n%s", result2.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, result2.Content), "world_revision"), "[WORLD REVISION] created=wr2 parent=wr1 current=wr2") {
		t.Fatalf("turn 2 missing world revision block:\n%s", result2.Content)
	}

	if got := readOverlayWorkspaceFile(t, b.subjective, "state.txt"); got != "v2" {
		t.Fatalf("current state = %q, want v2", got)
	}

	restore := b.SwitchWorld("tool-switch-world", "wr1")
	if restore.IsError {
		t.Fatalf("switch failed:\n%s", restore.Content)
	}
	restorePayload := decodeResultContent(t, restore.Content)
	if resultString(t, restorePayload, "fs_mutations") == "" {
		t.Fatalf("switch missing fs mutations block:\n%s", restore.Content)
	}
	if !strings.Contains(resultString(t, restorePayload, "fs_mutations"), "state.txt (modified)") {
		t.Fatalf("switch should report restored file mutation:\n%s", restore.Content)
	}
	if !strings.Contains(resultString(t, restorePayload, "world_revision"), "[WORLD REVISION] wr2 -> wr1") {
		t.Fatalf("switch missing world revision transition:\n%s", restore.Content)
	}
	if len(restore.StructuredContent) != 0 {
		t.Fatal("switch structured content should be empty")
	}
	if restorePayload["tool"] != "switch_world" {
		t.Fatalf("switch tool = %#v, want switch_world", restorePayload["tool"])
	}
	if !strings.Contains(restorePayload["fs_mutations"].(string), "state.txt (modified)") {
		t.Fatalf("switch structured fs_mutations missing state.txt modification: %#v", restorePayload["fs_mutations"])
	}
	if !strings.Contains(restorePayload["world_revision"].(string), "wr2 -> wr1") {
		t.Fatalf("switch structured world_revision missing transition: %#v", restorePayload["world_revision"])
	}

	if got := b.CurrentWorldRevision(); got != "wr1" {
		t.Fatalf("CurrentWorldRevision() after switch = %q, want wr1", got)
	}
	if got := readOverlayWorkspaceFile(t, b.subjective, "state.txt"); got != "v1" {
		t.Fatalf("switched state = %q, want v1", got)
	}

	b.TurnID = 3
	result3 := b.Execute("tool-checkpoint-3", "printf 'v3\\n' > state.txt", 0, 0, false, false, "")
	if result3.IsError {
		t.Fatalf("turn 3 failed:\n%s", result3.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, result3.Content), "world_revision"), "[WORLD REVISION] created=wr3 parent=wr1 current=wr3") {
		t.Fatalf("turn 3 should create monotonic revision off switched head:\n%s", result3.Content)
	}
	if got := readOverlayWorkspaceFile(t, b.subjective, "state.txt"); got != "v3" {
		t.Fatalf("post-restore state = %q, want v3", got)
	}
}

func TestWorkspaceOverlayVisionReaderSeesCommittedWorkspaceFiles(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-overlay-vision"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-overlay-vision-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)

	b.TurnID = 1
	result := b.Execute("tool-write-image", "printf 'P6\\n1 1\\n255\\n\\377\\000\\000' > image.ppm", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("write image failed:\n%s", result.Content)
	}
	if _, err := os.Stat(filepath.Join(root, "image.ppm")); !os.IsNotExist(err) {
		t.Fatalf("committed overlay file should not be host-visible at workspace root, err=%v", err)
	}

	vision := HandleVisionWithReader("vision-1", map[string]any{"path": "image.ppm"}, b.ReadWorkspaceFile)
	if vision.IsError {
		t.Fatalf("vision could not read committed overlay image:\n%s", vision.Content)
	}
	if vision.Image == nil {
		t.Fatal("vision returned nil image payload")
	}
	if vision.Image.MIMEType != "image/png" {
		t.Fatalf("vision image MIME = %q, want image/png", vision.Image.MIMEType)
	}

	absoluteVision := HandleVisionWithReader("vision-2", map[string]any{"path": filepath.Join(root, "image.ppm")}, b.ReadWorkspaceFile)
	if absoluteVision.IsError {
		t.Fatalf("vision could not read committed overlay image by absolute workspace path:\n%s", absoluteVision.Content)
	}
}

func TestWorkspaceOverlaySwitchWorldCanHideFSMutationTelemetry(t *testing.T) {
	root := t.TempDir()
	gates := config.DefaultToolGates()
	gates.FSMutationTelemetry = false
	cfg := &config.Config{

		ToolGates: gates, Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-switch-hidden-fs-mutations"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-switch-hidden-fs-mutations-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}

	b := NewShExecutor(cfg, nil)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)

	b.TurnID = 1
	result1 := b.Execute("tool-checkpoint-1", "printf 'v1\\n' > state.txt", 0, 0, false, false, "")
	if result1.IsError {
		t.Fatalf("turn 1 failed:\n%s", result1.Content)
	}
	b.TurnID = 2
	result2 := b.Execute("tool-checkpoint-2", "printf 'v2\\n' > state.txt", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("turn 2 failed:\n%s", result2.Content)
	}

	restore := b.SwitchWorld("tool-switch-world", "wr1")
	if restore.IsError {
		t.Fatalf("switch failed:\n%s", restore.Content)
	}
	restorePayload := decodeResultContent(t, restore.Content)
	if _, ok := restorePayload["fs_mutations"]; ok {
		t.Fatalf("switch_world should hide fs_mutations telemetry when disabled, got:\n%s", restore.Content)
	}
	if !strings.Contains(resultString(t, restorePayload, "world_revision"), "[WORLD REVISION] wr2 -> wr1") {
		t.Fatalf("switch_world should keep world_revision transition when telemetry is disabled:\n%s", restore.Content)
	}
}

func TestWorkspaceOverlaySwitchWorldAdoptsChildHandle(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	dataDir := t.TempDir()
	parentCfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-parent"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-parent-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"}}
	parent := NewShExecutor(parentCfg, nil)
	requireWorkspaceSupport(t, parent)

	childCfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-child"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-child-session", WorkspaceOwner: true, WorkspaceBootstrap: parentCfg.WorkspaceSession}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"}}
	child := NewShExecutor(childCfg, nil)
	requireWorkspaceSupport(t, child)

	child.TurnID = 1
	result := child.Execute("tool-child-write", "printf 'from child\\n' > adopted.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("child write failed:\n%s", result.Content)
	}
	if got := parent.CurrentWorldRevision(); got != "wr0" {
		t.Fatalf("parent current revision before switch = %q, want wr0", got)
	}

	exportDir := t.TempDir()
	if err := parent.subjective.exportCurrentTree(exportDir); err != nil {
		t.Fatalf("export parent tree before switch: %v", err)
	}
	if _, err := os.Stat(filepath.Join(exportDir, "adopted.txt")); !os.IsNotExist(err) {
		t.Fatalf("adopted.txt should be absent before switch, stat err = %v", err)
	}

	target := buildWorldHandle(childCfg.WorkspaceSession, child.CurrentWorldRevision())
	switchResult := parent.SwitchWorld("tool-switch-child", target)
	if switchResult.IsError {
		t.Fatalf("switch child world failed:\n%s", switchResult.Content)
	}
	if !strings.Contains(string(switchResult.Content), "adopted.txt (created)") {
		t.Fatalf("switch result should report adopted file creation:\n%s", switchResult.Content)
	}
	if got := readOverlayWorkspaceFile(t, parent.subjective, "adopted.txt"); got != "from child" {
		t.Fatalf("adopted file contents = %q, want from child", got)
	}
	if got := parent.CurrentWorldRevision(); got == "wr0" {
		t.Fatalf("parent current revision after adopt = %q, want imported revision", got)
	}
}

func TestWorkspaceNoopShellDoesNotAdvanceWorldRevision(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-noop-revision"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-noop-revision-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)

	b.TurnID = 1
	result1 := b.Execute("tool-noop-revision-1", "printf 'v1\\n' > state.txt", 0, 0, false, false, "")
	if result1.IsError {
		t.Fatalf("turn 1 failed:\n%s", result1.Content)
	}
	if !strings.Contains(resultString(t, decodeResultContent(t, result1.Content), "world_revision"), "[WORLD REVISION] created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("turn 1 missing revision creation:\n%s", result1.Content)
	}

	b.TurnID = 2
	result2 := b.Execute("tool-noop-revision-2", "test -f state.txt", 0, 0, false, false, "")
	if result2.IsError {
		t.Fatalf("turn 2 failed:\n%s", result2.Content)
	}
	payload := decodeResultContent(t, result2.Content)
	if !strings.Contains(resultString(t, payload, "fs_mutations"), "[FS MUTATIONS]\n(empty)") {
		t.Fatalf("turn 2 should report empty mutations:\n%s", result2.Content)
	}
	if !strings.Contains(resultString(t, payload, "world_revision"), "[WORLD REVISION] current=wr1 (unchanged)") {
		t.Fatalf("turn 2 should keep the same revision:\n%s", result2.Content)
	}
	if strings.Contains(resultString(t, payload, "world_revision"), "created=wr2") {
		t.Fatalf("turn 2 should not create wr2 on a no-op:\n%s", result2.Content)
	}
	if got := b.CurrentWorldRevision(); got != "wr1" {
		t.Fatalf("CurrentWorldRevision() = %q, want wr1", got)
	}
}

func TestWorkspaceOverlayStateVersionBreakRejected(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	dataDir := t.TempDir()
	stateDir := filepath.Join(dataDir, "workspaces", "legacy-session")
	if err := os.MkdirAll(filepath.Join(stateDir, "upper"), 0o755); err != nil {
		t.Fatalf("mkdir legacy state: %v", err)
	}

	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-legacy-state"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "legacy-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: dataDir, Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	if err := b.Prepare(); err == nil {
		t.Fatal("Prepare() unexpectedly succeeded for pre-lineage overlay state")
	} else if !strings.Contains(err.Error(), "pre-lineage overlay state") {
		t.Fatalf("Prepare() error = %v, want pre-lineage overlay state rejection", err)
	}
}

func TestWorkspaceOverlayReportsDeletionAndReplacementMutations(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	mustWrite := func(path string, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	mustWrite(filepath.Join(root, "opaque", "old.txt"), "old\n")
	mustWrite(filepath.Join(root, "swap", "child.txt"), "child\n")
	mustWrite(filepath.Join(root, "dir-to-file", "child.txt"), "child\n")
	mustWrite(filepath.Join(root, "file-to-dir"), "leaf\n")

	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-overlay-mutations"}, Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 10}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: root, Workspace: root, WorkspaceRevisionMode: config.WorkspaceRevisionRestore, WorkspaceSession: "workspace-overlay-mutations-session", WorkspaceOwner: true}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}}

	b := NewShExecutor(cfg, nil)
	requireWorkspaceSupport(t, b)

	result := b.Execute("tool-overlay-mutations", `
rm -rf opaque
mkdir opaque
printf 'new\n' > opaque/new.txt
rm -rf swap
printf 'file\n' > swap
rm -rf dir-to-file
printf 'leaf\n' > dir-to-file
rm file-to-dir
mkdir file-to-dir
printf 'child\n' > file-to-dir/child.txt
`, 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected overlay mutation error:\n%s", result.Content)
	}
	checks := []string{
		"- opaque/old.txt (deleted)",
		"+ opaque/new.txt (created)",
		"- swap/child.txt (deleted)",
		"~ swap (modified)",
		"- dir-to-file/child.txt (deleted)",
		"~ dir-to-file (modified)",
		"~ file-to-dir (modified)",
		"+ file-to-dir/child.txt (created)",
	}
	for _, want := range checks {
		if !strings.Contains(string(result.Content), want) {
			t.Fatalf("result missing mutation %q:\n%s", want, result.Content)
		}
	}
}

func TestParseSwitchWorldArgs(t *testing.T) {
	req, err := ParseSwitchWorldArgs(map[string]any{"target": "wr3"})
	if err != nil {
		t.Fatalf("ParseSwitchWorldArgs() error = %v", err)
	}
	if req.Target != "wr3" {
		t.Fatalf("Target = %q, want wr3", req.Target)
	}

	req, err = ParseSwitchWorldArgs(map[string]any{"target": "world://child-session/wr4"})
	if err != nil {
		t.Fatalf("ParseSwitchWorldArgs() world handle error = %v", err)
	}
	if req.Target != "world://child-session/wr4" {
		t.Fatalf("Target = %q, want world handle", req.Target)
	}

	if _, err := ParseSwitchWorldArgs(map[string]any{"target": 3}); err == nil {
		t.Fatal("expected type error for numeric revision")
	}
	if _, err := ParseSwitchWorldArgs(map[string]any{"target": "3"}); err == nil {
		t.Fatal("expected validation error for invalid target")
	}
	if _, err := ParseSwitchWorldArgs(map[string]any{"target": "world://child-session/not-a-rev"}); err == nil {
		t.Fatal("expected validation error for malformed world handle")
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
			name: "with target argv",
			args: map[string]any{
				"target": "/bin/sh",
				"argv":   []any{"/bin/sh", "-c", "echo hi"},
			},
			want: ExecRequest{
				Target: "/bin/sh",
				Argv:   []string{"/bin/sh", "-c", "echo hi"},
			},
			wantErr: false,
		},
		{
			name: "target wrong type",
			args: map[string]any{
				"target": 123,
			},
			want:    ExecRequest{},
			wantErr: true,
		},
		{
			name: "argv wrong type",
			args: map[string]any{
				"argv": "not-array",
			},
			want:    ExecRequest{},
			wantErr: true,
		},
		{
			name: "argv element wrong type",
			args: map[string]any{
				"argv": []any{"/bin/sh", 42},
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
				if got.Target != tt.want.Target {
					t.Errorf("Target = %q, want %q", got.Target, tt.want.Target)
				}
				if !reflect.DeepEqual(got.Argv, tt.want.Argv) {
					t.Errorf("Argv = %v, want %v", got.Argv, tt.want.Argv)
				}
			}
		})
	}
}

func TestExecEnv(t *testing.T) {
	cfg := &config.Config{
		Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", Depth: 3, SessionID: "pre-exec-session", RunID: "pre-exec-run", TapeID: "legacy-tape", ParentSession: "root-session"}, Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"}, Limits: config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 600, OutputTruncate: 20480, MaxTurns: 20}, Paths: config.Paths{DataDir: "/tmp/quine-test", Shell: "/bin/sh"},
	}

	execEnv, err := cfg.ExecEnv()
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
	if envMap["QUINE_SESSION_ID"] != "pre-exec-session" {
		t.Errorf("QUINE_SESSION_ID = %q, want pre-exec-session", envMap["QUINE_SESSION_ID"])
	}
	if _, ok := envMap["QUINE_RUN_ID"]; ok {
		t.Errorf("QUINE_RUN_ID should not be propagated across exec")
	}
	if envMap["QUINE_PARENT_SESSION"] != "root-session" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want root-session", envMap["QUINE_PARENT_SESSION"])
	}
	if envMap["QUINE_TAPE_ID"] != "legacy-tape" {
		t.Errorf("QUINE_TAPE_ID = %q, want legacy-tape", envMap["QUINE_TAPE_ID"])
	}
	if envMap["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want claude-sonnet-4-20250514", envMap["QUINE_MODEL_ID"])
	}
	if envMap["QUINE_MAX_DEPTH"] != "5" {
		t.Errorf("QUINE_MAX_DEPTH = %q, want 5", envMap["QUINE_MAX_DEPTH"])
	}
}

// writeStagedNextEnv plants a staged-override file under the agent-root
// config/ dir the way an agent's sh call would.
func writeStagedNextEnv(t *testing.T, cfg *config.Config, content string) string {
	t.Helper()
	path := cfg.StagedNextEnvPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir staged config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write staged config: %v", err)
	}
	return path
}

// TestStagedOverridesExecPathOnlyChildIsolation is the child-isolation
// regression test required by registry-design-brief § C (F-3): the staged
// next.env merge lives strictly in the exec path — fork/spawn children built
// from ChildEnv() (which shares baseEnv with ExecEnv) must NEVER see staged
// values.
func TestStagedOverridesExecPathOnlyChildIsolation(t *testing.T) {
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "test-model", SessionID: "staged-isolation-session", TapeID: "tape-staged"},
		Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"},
		Limits:    config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 10, OutputTruncate: 20480},
		Paths:     config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}
	writeStagedNextEnv(t, cfg, "QUINE_OUTPUT_TRUNCATE=31337\n")

	// Exec path: the staged value must be merged over the ExecEnv() output.
	execEnv, err := cfg.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	staged, err := config.ReadStagedOverrides(cfg.StagedNextEnvPath())
	if err != nil {
		t.Fatalf("ReadStagedOverrides() error: %v", err)
	}
	merged := config.MergeStagedOverrides(execEnv, staged)
	mergedMap := envSliceToMap(merged)
	if mergedMap["QUINE_OUTPUT_TRUNCATE"] != "31337" {
		t.Fatalf("exec path QUINE_OUTPUT_TRUNCATE = %q, want staged 31337", mergedMap["QUINE_OUTPUT_TRUNCATE"])
	}

	// Child path: fork/spawn children never see the staged value.
	childEnv, err := cfg.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	childMap := envSliceToMap(childEnv)
	if childMap["QUINE_OUTPUT_TRUNCATE"] != "20480" {
		t.Fatalf("child QUINE_OUTPUT_TRUNCATE = %q, want the in-process 20480 — staged next.env leaked into ChildEnv", childMap["QUINE_OUTPUT_TRUNCATE"])
	}
}

func TestExecRejectsInvalidStagedConfigAndKeepsFileIntact(t *testing.T) {
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "test-model", SessionID: "staged-reject-session", TapeID: "tape-staged-reject"},
		Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"},
		Limits:    config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 10, OutputTruncate: 20480},
		Paths:     config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}
	content := "QUINE_API_KEY=stolen\nQUINE_MAX_TURNS=soon\n"
	path := writeStagedNextEnv(t, cfg, content)

	exec := &ExecExecutor{QuinePath: "/nonexistent/quine", Cfg: cfg, Mission: "staged reject"}
	result := exec.Execute("tool-exec-staged", ExecRequest{})
	if !result.IsError {
		t.Fatal("exec with an invalid staged config should fail as a normal tool error")
	}
	text := string(result.Content)
	for _, want := range []string{
		"[EXEC ERROR]",
		"QUINE_API_KEY",            // mutability violation named
		"QUINE_MAX_TURNS",          // type violation named
		"only exec-boundary knobs", // the whitelist, agent-legible
		"left intact",              // fix-and-retry guidance
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("exec rejection missing %q:\n%s", want, text)
		}
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("re-read staged file: %v", err)
	}
	if string(after) != content {
		t.Fatalf("rejected exec must leave next.env intact: %q -> %q", content, after)
	}
}

func TestExecAbsentStagedConfigIsNoOp(t *testing.T) {
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "test-model", SessionID: "staged-absent-session", TapeID: "tape-staged-absent"},
		Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"},
		Limits:    config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 10, OutputTruncate: 20480},
		Paths:     config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}

	// No next.env exists: exec proceeds to the syscall and fails only on the
	// bogus target, not on staged validation.
	exec := &ExecExecutor{QuinePath: "/nonexistent/quine", Cfg: cfg, Mission: "staged absent"}
	result := exec.Execute("tool-exec-absent", ExecRequest{})
	if !result.IsError {
		t.Fatal("expected exec failure on the nonexistent target")
	}
	text := string(result.Content)
	if !strings.Contains(text, "syscall.Exec failed") {
		t.Fatalf("failure should come from the exec syscall, not staged validation:\n%s", text)
	}
	if strings.Contains(text, "staged config") {
		t.Fatalf("absent staged file must be a silent no-op:\n%s", text)
	}
}

func envSliceToMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, entry := range env {
		key, val, _ := strings.Cut(entry, "=")
		m[key] = val
	}
	return m
}

func TestExecProcessSurfaceEnvIncludesAgentRoot(t *testing.T) {
	cfg := &config.Config{Identity: config.Identity{SessionID: "exec-surface-session", RunID: "exec-surface-run"}, Paths: config.Paths{DataDir: t.TempDir()}}

	env := execProcessSurfaceEnv(cfg)
	envMap := make(map[string]string)
	for _, entry := range env {
		key, val, _ := strings.Cut(entry, "=")
		envMap[key] = val
	}

	if envMap["QUINE_AGENT_ROOT"] != cfg.AgentRoot() {
		t.Fatalf("QUINE_AGENT_ROOT = %q, want %q", envMap["QUINE_AGENT_ROOT"], cfg.AgentRoot())
	}
	if _, ok := envMap["QUINE_RUN_ID"]; ok {
		t.Fatal("exec process surface env should not propagate QUINE_RUN_ID")
	}
}

func TestExecExecutor_FailureUsesJSONContentOnly(t *testing.T) {
	exec := &ExecExecutor{
		QuinePath: "/nonexistent/quine",
		Cfg:       &config.Config{Identity: config.Identity{SessionID: "exec-failure-session", ModelID: "test-model"}, Transport: config.Transport{Provider: "anthropic", APIKey: "test-key"}, Limits: config.Limits{OutputTruncate: 1024, ShTimeout: 10}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"}},
		Mission:   "test exec failure",
	}

	result := exec.Execute("tool-exec", ExecRequest{
		Target: "definitely-missing-executable-for-structured-content",
	})
	if !result.IsError {
		t.Fatal("expected exec failure result")
	}
	if len(result.StructuredContent) != 0 {
		t.Fatal("expected exec failure structured content to be empty")
	}
	payload := decodeResultContent(t, result.Content)
	if payload["tool"] != "exec" {
		t.Fatalf("exec tool = %#v, want exec", payload["tool"])
	}
	if payload["status"] != "error" {
		t.Fatalf("exec status = %#v, want error", payload["status"])
	}
	if payload["target"] != "definitely-missing-executable-for-structured-content" {
		t.Fatalf("exec target = %#v, want definitely-missing-executable-for-structured-content", payload["target"])
	}
}

func TestExecExecutor_DefaultSelfReentryArgvOmitsAbsentMission(t *testing.T) {
	for _, tt := range []struct {
		name    string
		mission string
		want    []string
	}{
		{name: "with mission", mission: "continue current work", want: []string{"/tmp/quine", "continue current work"}},
		{name: "without mission", mission: "", want: []string{"/tmp/quine"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			exec := &ExecExecutor{
				QuinePath: "/tmp/quine",
				Mission:   tt.mission,
			}
			target, argv, err := exec.buildTargetAndArgv(ExecRequest{})
			if err != nil {
				t.Fatalf("buildTargetAndArgv error: %v", err)
			}
			if target != "/tmp/quine" {
				t.Fatalf("target = %q, want /tmp/quine", target)
			}
			if !reflect.DeepEqual(argv, tt.want) {
				t.Fatalf("argv = %#v, want %#v", argv, tt.want)
			}
		})
	}
}

func TestExecExecutor_RelativeTargetResolvesFromWorkspace(t *testing.T) {
	workspace := t.TempDir()
	exec := &ExecExecutor{
		QuinePath: "/tmp/quine",
		Cfg:       &config.Config{WorkspaceConfig: config.WorkspaceConfig{Workspace: workspace}},
		Mission:   "test mission",
	}

	target, argv, err := exec.buildTargetAndArgv(ExecRequest{Target: "outbox/quine-next"})
	if err != nil {
		t.Fatalf("buildTargetAndArgv error: %v", err)
	}
	want := filepath.Join(workspace, "outbox", "quine-next")
	if target != want {
		t.Fatalf("target = %q, want %q", target, want)
	}
	if !reflect.DeepEqual(argv, []string{want}) {
		t.Fatalf("argv = %#v, want %#v", argv, []string{want})
	}
}

func TestExecExecutor_PathLookupStillHandlesBareExecutableNames(t *testing.T) {
	exec := &ExecExecutor{
		QuinePath: "/tmp/quine",
		Cfg:       &config.Config{WorkspaceConfig: config.WorkspaceConfig{Workspace: t.TempDir()}},
		Mission:   "test mission",
	}

	target, argv, err := exec.buildTargetAndArgv(ExecRequest{Target: "sh"})
	if err != nil {
		t.Fatalf("buildTargetAndArgv error: %v", err)
	}
	if filepath.Base(target) != "sh" || !filepath.IsAbs(target) {
		t.Fatalf("target = %q, want absolute sh path", target)
	}
	if !reflect.DeepEqual(argv, []string{target}) {
		t.Fatalf("argv = %#v, want %#v", argv, []string{target})
	}
}

func TestNewExecExecutor_UsesSelfReentryTargetWhenEphemeralBodyEnabled(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{EphemeralBody: true}, Paths: config.Paths{ExecutablePath: "/tmp/launch-path-quine", SelfReentryTarget: "/proc/self/exe"}}

	exec := NewExecExecutor(cfg, "test mission")
	if got, want := exec.QuinePath, cfg.SelfReentryTarget; got != want {
		t.Fatalf("QuinePath = %q, want %q", got, want)
	}
}

func TestStageExecContextBootstrapStripsToolMechanics(t *testing.T) {
	contextRoot := filepath.Join(t.TempDir(), "context")
	if err := os.MkdirAll(filepath.Join(contextRoot, "prompt"), 0o755); err != nil {
		t.Fatalf("mkdir prompt: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(contextRoot, "state"), 0o755); err != nil {
		t.Fatalf("mkdir state: %v", err)
	}
	if err := os.WriteFile(filepath.Join(contextRoot, "prompt", "30-memory.md"), []byte("carry this forward\n"), 0o644); err != nil {
		t.Fatalf("write memory: %v", err)
	}

	entries := []tape.TapeEntry{
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "Begin."}),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "plain checkpoint"}),
		tape.MessageEntry(tape.Message{
			Role:             tape.RoleAssistant,
			ReasoningContent: "internal chain of thought that should not survive exec",
		}),
		tape.MessageEntry(tape.Message{
			Role:      tape.RoleAssistant,
			Content:   "Now executing the rebuilt binary.",
			ToolCalls: []tape.ToolCall{{ID: "exec:1", Name: "exec", Arguments: map[string]any{"target": "/tmp/rebuilt-quine"}}},
		}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID:  "exec:1",
			Content: tape.MarshalToolResultContent(map[string]any{"tool": "exec", "status": "error"}),
			IsError: true,
		}),
		tape.MessageEntry(tape.Message{
			Role:      tape.RoleAssistant,
			ToolCalls: []tape.ToolCall{{ID: "exec:2", Name: "exec", Arguments: map[string]any{"target": "/tmp/rebuilt-quine"}}},
		}),
	}
	currentPath := filepath.Join(contextRoot, "state", "current.jsonl")
	var lines [][]byte
	for _, entry := range entries {
		line, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("marshal entry: %v", err)
		}
		lines = append(lines, line)
	}
	if err := os.WriteFile(currentPath, append(bytes.Join(lines, []byte("\n")), '\n'), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}
	contextAlias := filepath.Join(t.TempDir(), "context-link")
	if err := os.Symlink(contextRoot, contextAlias); err != nil {
		t.Fatalf("symlink context root: %v", err)
	}

	bootstrapRoot, err := stageExecContextBootstrap(t.TempDir(), contextAlias)
	if err != nil {
		t.Fatalf("stageExecContextBootstrap error: %v", err)
	}

	memoryData, err := os.ReadFile(filepath.Join(bootstrapRoot, "prompt", "30-memory.md"))
	if err != nil {
		t.Fatalf("read copied memory: %v", err)
	}
	if string(memoryData) != "carry this forward\n" {
		t.Fatalf("copied memory = %q, want preserved memory", string(memoryData))
	}
	handoffPrompt, err := os.ReadFile(filepath.Join(bootstrapRoot, "prompt", "35-exec-handoff.md"))
	if err != nil {
		t.Fatalf("read exec handoff prompt: %v", err)
	}
	if !strings.Contains(string(handoffPrompt), "already the successor of a completed") {
		t.Fatalf("exec handoff prompt = %q, want successor guidance", string(handoffPrompt))
	}

	currentData, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("read projected current context: %v", err)
	}
	got := string(currentData)
	if !strings.Contains(got, "Begin.") {
		t.Fatalf("projected current context should keep user material, got %q", got)
	}
	if !strings.Contains(got, "plain checkpoint") {
		t.Fatalf("projected current context should keep plain assistant text, got %q", got)
	}
	if strings.Contains(got, "Now executing the rebuilt binary.") {
		t.Fatalf("projected current context should drop assistant tool-call narration, got %q", got)
	}
	if strings.Contains(got, "internal chain of thought that should not survive exec") {
		t.Fatalf("projected current context should strip assistant reasoning, got %q", got)
	}
	if strings.Contains(got, "\"tool_calls\"") {
		t.Fatalf("projected current context should drop tool calls, got %q", got)
	}
	if strings.Contains(got, "\"type\":\"tool_result\"") {
		t.Fatalf("projected current context should drop tool results, got %q", got)
	}
}

// --- helpers ---

func extractPID(content any) int {
	var (
		payload map[string]any
		err     error
	)
	switch value := content.(type) {
	case string:
		payload, err = tape.UnmarshalToolResultContent(value)
	case json.RawMessage:
		payload, err = tape.UnmarshalStructuredToolResultContent(value)
	default:
		return 0
	}
	if err != nil {
		return 0
	}
	job, _ := payload["job"].(map[string]any)
	if job == nil {
		return 0
	}
	switch value := job["pid"].(type) {
	case json.Number:
		n, _ := value.Int64()
		return int(n)
	case float64:
		return int(value)
	}
	return 0
}
