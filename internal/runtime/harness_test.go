package runtime

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

// mockProvider is a test double that returns pre-programmed responses.
type mockProvider struct {
	responses []tape.Message
	callCount int
	calls     [][]tape.Message
}

func (m *mockProvider) Generate(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	call := make([]tape.Message, len(msgs))
	copy(call, msgs)
	m.calls = append(m.calls, call)
	if m.callCount >= len(m.responses) {
		return tape.Message{}, llm.Usage{}, fmt.Errorf("mock: no more responses (call %d)", m.callCount)
	}
	resp := m.responses[m.callCount]
	m.callCount++
	return resp, llm.Usage{InputTokens: 100, OutputTokens: 50}, nil
}

func (m *mockProvider) ContextWindowSize() int { return 200000 }

type fakePublicSurface struct {
	syncFn    func(publicSurfacePaths) error
	cleanupFn func() error
}

func (f *fakePublicSurface) sync(paths publicSurfacePaths) error {
	if f.syncFn != nil {
		return f.syncFn(paths)
	}
	return nil
}

func (f *fakePublicSurface) cleanup() error {
	if f.cleanupFn != nil {
		return f.cleanupFn()
	}
	return nil
}

// installFakePublicSurface wires a fake through the Runtime's public-surface
// test seams so finalization-ordering tests run without a real FUSE mount.
func installFakePublicSurface(rt *Runtime, f *fakePublicSurface) {
	rt.publicSurfaceSyncFn = f.sync
	rt.publicSurfaceCleanupFn = f.cleanup
}

// mockErrorProvider returns errors.
type mockErrorProvider struct {
	err error
}

func (m *mockErrorProvider) Generate(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	return tape.Message{}, llm.Usage{}, m.err
}

func (m *mockErrorProvider) ContextWindowSize() int { return 200000 }

type providerFunc struct {
	generate func([]tape.Message, []llm.ToolSchema) (tape.Message, llm.Usage, error)
}

func (p providerFunc) Generate(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	return p.generate(msgs, tools)
}

func (p providerFunc) ContextWindowSize() int { return 200000 }

func finalizationPhaseIndex(phases []finalizationPhase, phase finalizationPhase) int {
	for i, candidate := range phases {
		if candidate == phase {
			return i
		}
	}
	return -1
}

func readFinalizationPhases(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read finalization state: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	phases := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record finalizationPhaseRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("parse finalization state line %q: %v", line, err)
		}
		phases = append(phases, record.Phase)
	}
	return phases
}

func readWorkspaceCommitIntent(t *testing.T, path string) workspaceCommitIntent {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read workspace commit intent: %v", err)
	}
	var intent workspaceCommitIntent
	if err := json.Unmarshal(data, &intent); err != nil {
		t.Fatalf("parse workspace commit intent: %v", err)
	}
	return intent
}

func testCfg(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Identity: config.Identity{
			ModelID:   "claude-sonnet-4-20250514",
			Depth:     0,
			SessionID: "test-1234-5678",
			RunID:     "run-test-1234",
			TapeID:    "tape-1234-5678",
		},
		Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"},
		Limits: config.Limits{
			MaxDepth:                 5,
			MaxConcurrent:            20,
			MaxAgents:                10,
			ShTimeout:                10,
			OutputTruncate:           20480,
			MaxTurns:                 0, // disabled for existing tests
			MemoryWarnTokens:         8000,
			MemoryDangerTokens:       16000,
			PeerDiscoveryHeartbeatMS: 5000,
		},
		Paths:        config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
		PromptConfig: config.PromptConfig{PromptMetaphor: config.PromptMetaphorOff, FailOnImpossible: true},
		ToolGates:    config.ToolGates{ExecEnabled: true, VisionEnabled: true},
	}
}

// silenceRuntime suppresses all runtime output for clean test output.
func silenceRuntime(rt *Runtime) {
	devnull, _ := os.Open(os.DevNull)
	rt.stderr = devnull
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}
	if rt.logFile != nil {
		rt.logFile.Close()
		rt.logFile = nil
	}
	// Stub the public surface by default so generic tests stay hermetic: the
	// control surface is FUSE-only, and mounting real FUSE in every bootstrap
	// would leave lingering mounts and require /dev/fuse everywhere. Tests that
	// exercise the real projection call useRealPublicSurface to opt back in.
	installFakePublicSurface(rt, &fakePublicSurface{})
}

// useRealPublicSurface re-enables the real FUSE public-surface backend for a
// runtime that silenceRuntime stubbed out. Pair with requireRuntimeSurfaceFUSESupport.
func useRealPublicSurface(rt *Runtime) {
	rt.publicSurfaceSyncFn = nil
	rt.publicSurfaceCleanupFn = nil
}

func decodeToolContent(t *testing.T, content any) map[string]any {
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

func requireRuntimeSurfaceFUSESupport(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("fuse runtime surface is Linux-only")
	}
	if err := preflightRuntimeSurfaceFUSE(); err != nil {
		t.Skipf("fuse runtime surface unsupported in this Linux environment: %v", err)
	}
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	useRealPublicSurface(rt)
	rt.originalInput = "probe fuse runtime surface support"
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Skipf("fuse runtime surface unsupported in this Linux environment: %v", err)
	}
	// A denied mount no longer fails bootstrapAgentRoot — it degrades. Probe
	// the recorded state so FUSE tests still skip where the mount cannot land.
	degradedReason := rt.publicSurfaceUnavailableReason()
	if err := rt.cleanupAgentRoot(); err != nil {
		t.Fatalf("cleanup fuse runtime surface probe: %v", err)
	}
	if degradedReason != "" {
		t.Skipf("fuse runtime surface unsupported in this Linux environment: %s", degradedReason)
	}
}

func requireOverlayWorkspaceSupport(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}
	root := t.TempDir()
	cfg := testCfg(t)
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = root
	cfg.Workspace = root
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceSession = "overlay-workspace-probe"
	cfg.WorkspaceOwner = true
	sh := tools.NewShExecutor(cfg)
	if err := sh.Prepare(); err != nil {
		t.Skipf("overlay workspace unsupported in this Linux environment: %v", err)
	}
	if err := sh.Close(false); err != nil {
		t.Fatalf("cleanup overlay workspace probe: %v", err)
	}
}

func writeControlActionFile(t *testing.T, path string, payload string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(payload), 0o666); err != nil {
		t.Fatalf("write control action file: %v", err)
	}
}

func toolString(t *testing.T, payload map[string]any, key string) string {
	t.Helper()
	value, _ := payload[key].(string)
	return value
}

func toolMap(t *testing.T, payload map[string]any, key string) map[string]any {
	t.Helper()
	value, _ := payload[key].(map[string]any)
	if value == nil {
		t.Fatalf("expected %q map in payload: %#v", key, payload)
	}
	return value
}

func toolInt(t *testing.T, payload map[string]any, key string) int {
	t.Helper()
	switch value := payload[key].(type) {
	case json.Number:
		n, err := value.Int64()
		if err != nil {
			t.Fatalf("parse %q int from json.Number: %v", key, err)
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

func waitForCondition(t *testing.T, timeout time.Duration, condition func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func assertSameFile(t *testing.T, a, b string) {
	t.Helper()
	infoA, err := os.Stat(a)
	if err != nil {
		t.Fatalf("stat %s: %v", a, err)
	}
	infoB, err := os.Stat(b)
	if err != nil {
		t.Fatalf("stat %s: %v", b, err)
	}
	if !os.SameFile(infoA, infoB) {
		t.Fatalf("%s and %s are not the same file", a, b)
	}
}

func assertSymlinkTarget(t *testing.T, path, want string) {
	t.Helper()
	target, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("read symlink %s: %v", path, err)
	}
	if target != want {
		t.Fatalf("symlink %s = %q, want %q", path, target, want)
	}
}

func assertRelativeSymlinkResolvesTo(t *testing.T, path, want string) {
	t.Helper()
	target, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("read symlink %s: %v", path, err)
	}
	if filepath.IsAbs(target) {
		t.Fatalf("symlink %s target is absolute: %q", path, target)
	}
	assertSameFile(t, path, want)
}

func writeJSONLEntries(t *testing.T, path string, entries ...tape.TapeEntry) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	data := make([]byte, 0, 256)
	for _, entry := range entries {
		line, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("marshal entry for %s: %v", path, err)
		}
		data = append(data, line...)
		data = append(data, '\n')
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func copyAgentRootForResumeTest(t *testing.T, src, dst string) {
	t.Helper()
	info, err := os.Lstat(src)
	if err != nil {
		t.Fatalf("stat source agent root: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("source agent root %s is not a directory", src)
	}
	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		t.Fatalf("mkdir target agent root: %v", err)
	}
	if err := filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == src {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		switch mode := info.Mode(); {
		case mode&os.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(target, targetPath)
		case d.IsDir():
			return os.MkdirAll(targetPath, info.Mode().Perm())
		case mode&os.ModeNamedPipe != 0:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			return syscall.Mkfifo(targetPath, uint32(info.Mode().Perm()))
		case mode.IsRegular():
			return copyFile(path, targetPath)
		default:
			return fmt.Errorf("unsupported test agent-root entry %s with mode %v", path, mode)
		}
	}); err != nil {
		t.Fatalf("copy agent root: %v", err)
	}
}

func seedPendingToolBatch(t *testing.T, cfg *config.Config, calls ...tape.ToolCall) {
	t.Helper()
	writeJSONLEntries(t, filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "before pending tool"}),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, ToolCalls: calls}),
	)
}

func requireToolResultInMessages(t *testing.T, msgs []tape.Message, toolID string) map[string]any {
	t.Helper()
	for _, msg := range msgs {
		if msg.Role != tape.RoleToolResult || msg.ToolID != toolID {
			continue
		}
		return decodeToolContent(t, msg.StructuredContent)
	}
	t.Fatalf("missing tool result %q in %#v", toolID, msgs)
	return nil
}

func countToolResultsInMessages(msgs []tape.Message, toolID string) int {
	count := 0
	for _, msg := range msgs {
		if msg.Role == tape.RoleToolResult && msg.ToolID == toolID {
			count++
		}
	}
	return count
}

func writeJobSurface(t *testing.T, cfg *config.Config, pid string, fields map[string]string) string {
	t.Helper()
	jobDir := filepath.Join(cfg.JobSessionDir(""), pid)
	if err := os.MkdirAll(jobDir, 0o755); err != nil {
		t.Fatalf("mkdir job dir: %v", err)
	}
	if _, ok := fields["pid"]; !ok {
		fields["pid"] = pid + "\n"
	}
	if _, ok := fields["started_at"]; !ok {
		fields["started_at"] = time.Now().UTC().Format(time.RFC3339Nano) + "\n"
	}
	for name, data := range fields {
		if err := os.WriteFile(filepath.Join(jobDir, name), []byte(data), 0o644); err != nil {
			t.Fatalf("write job %s: %v", name, err)
		}
	}
	return jobDir
}

func exitAfterRecoveryProvider() *mockProvider {
	return &mockProvider{responses: []tape.Message{{
		Role: tape.RoleAssistant,
		ToolCalls: []tape.ToolCall{{
			ID:        "call_exit",
			Name:      "exit",
			Arguments: map[string]any{"status": "success"},
		}},
	}}}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	return data
}

func restoreEnv(key, value string) {
	if value == "" {
		_ = os.Unsetenv(key)
		return
	}
	_ = os.Setenv(key, value)
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
