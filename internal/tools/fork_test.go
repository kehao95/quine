package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

func decodeForkResultContent(t *testing.T, content any) map[string]any {
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
		t.Fatalf("unsupported fork result content type %T", content)
	}
	if err != nil {
		t.Fatalf("decode fork tool result content: %v\ncontent=%v", err, content)
	}
	return payload
}

func requireRelationStatus(t *testing.T, payload map[string]any) map[string]any {
	t.Helper()
	relationRoot, ok := payload["relation_root"].(string)
	if !ok || strings.TrimSpace(relationRoot) == "" {
		t.Fatalf("relation_root missing from payload: %#v", payload)
	}
	data, err := os.ReadFile(filepath.Join(relationRoot, "status.json"))
	if err != nil {
		t.Fatalf("read relation status: %v", err)
	}
	var status map[string]any
	if err := json.Unmarshal(data, &status); err != nil {
		t.Fatalf("decode relation status: %v\ncontent=%s", err, data)
	}
	return status
}

func requireRelationResult(t *testing.T, payload map[string]any) map[string]any {
	t.Helper()
	relationRoot, ok := payload["relation_root"].(string)
	if !ok || strings.TrimSpace(relationRoot) == "" {
		t.Fatalf("relation_root missing from payload: %#v", payload)
	}
	data, err := os.ReadFile(filepath.Join(relationRoot, "result.json"))
	if err != nil {
		t.Fatalf("read relation result: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("decode relation result: %v\ncontent=%s", err, data)
	}
	return result
}

func requireStableRelationCounters(t *testing.T, payload map[string]any, spawned, completed, succeeded, killed int) {
	t.Helper()
	want := map[string]int{
		"spawned":   spawned,
		"completed": completed,
		"succeeded": succeeded,
		"killed":    killed,
	}
	for key, expected := range want {
		if _, ok := payload[key]; !ok {
			t.Fatalf("%s missing from relation payload: %#v", key, payload)
		}
		if got := resultInt(t, payload, key); got != expected {
			t.Fatalf("%s = %d, want %d in payload %#v", key, got, expected, payload)
		}
	}
}

func prepareForkContextRoot(t *testing.T, root string) {
	t.Helper()
	if strings.TrimSpace(root) == "" {
		return
	}
	for _, dir := range []string{root, filepath.Join(root, "prompt"), filepath.Join(root, "state")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir context root: %v", err)
		}
	}
}

func prepareForkPromptRoot(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "prompt"), 0o755); err != nil {
		t.Fatalf("mkdir context root: %v", err)
	}
}

func requireOverlayWorkspaceSupport(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics require Linux")
	}
	if _, err := preflightOverlayMount(t.TempDir()); err != nil {
		t.Skipf("workspace overlay mount unsupported in this Linux environment: %v", err)
	}
}

func decodeContextEntries(t *testing.T, data []byte) []tape.TapeEntry {
	t.Helper()
	lines := bytes.Split(data, []byte("\n"))
	entries := make([]tape.TapeEntry, 0, len(lines))
	for i, raw := range lines {
		line := bytes.TrimSpace(raw)
		if len(line) == 0 {
			continue
		}
		var entry tape.TapeEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Fatalf("decode context entry %d: %v", i, err)
		}
		entries = append(entries, entry)
	}
	return entries
}

func TestNewForkExecutorLeavesDefaultTimeoutDisabled(t *testing.T) {
	cfg := &config.Config{Identity: config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "fork-timeout-default"}, Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"}, Limits: config.Limits{OutputTruncate: 20480, ShTimeout: 123}, Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh", SelfReentryTarget: "/tmp/quine-self-reentry"}}

	fork := NewForkExecutor(cfg)
	if fork.QuinePath != cfg.SelfReentryTarget {
		t.Fatalf("QuinePath = %q, want %q", fork.QuinePath, cfg.SelfReentryTarget)
	}
	if fork.DefaultTimeout != 0 {
		t.Fatalf("DefaultTimeout = %v, want 0", fork.DefaultTimeout)
	}
}

func TestNewForkExecutorUsesConfiguredDefaultTimeout(t *testing.T) {
	cfg := &config.Config{Identity: config.Identity{SessionID: "fork-timeout-configured"}, Limits: config.Limits{ForkDefaultTimeoutSeconds: 7}, Paths: config.Paths{DataDir: t.TempDir(), SelfReentryTarget: "/tmp/quine-self-reentry"}}

	fork := NewForkExecutor(cfg)
	if fork.DefaultTimeout != 7*time.Second {
		t.Fatalf("DefaultTimeout = %v, want 7s", fork.DefaultTimeout)
	}
}

func TestNewForkExecutor_UsesSelfReentryTargetWhenEphemeralBodyEnabled(t *testing.T) {
	cfg := &config.Config{Identity: config.Identity{SessionID: "fork-ephemeral-target"}, ToolGates: config.ToolGates{EphemeralBody: true}, Paths: config.Paths{DataDir: t.TempDir(), ExecutablePath: "/tmp/launch-path-quine", SelfReentryTarget: "/proc/self/exe"}}

	fork := NewForkExecutor(cfg)
	if got, want := fork.QuinePath, cfg.SelfReentryTarget; got != want {
		t.Fatalf("QuinePath = %q, want %q", got, want)
	}
}

func TestParseForkArgs_SingleChild(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "Do something useful", "scope": "."},
		},
	}
	req, err := ParseForkArgs(args, nil)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if len(req.Children) != 1 || req.Children[0].Intent != "Do something useful" {
		t.Errorf("Children = %v", req.Children)
	}
	if req.Mode != ForkModeRace {
		t.Errorf("Mode = %q, want %q (default)", req.Mode, ForkModeRace)
	}
}

func TestParseForkArgs_MultipleChildren(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "Task A", "scope": "."},
			map[string]any{"intent": "Task B", "scope": "sub"},
			map[string]any{"intent": "Task C", "scope": "sub/deeper"},
		},
	}
	req, err := ParseForkArgs(args, nil)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if len(req.Children) != 3 {
		t.Fatalf("got %d children, want 3", len(req.Children))
	}
	if req.Children[0].Intent != "Task A" || req.Children[1].Scope != "sub" || req.Children[2].Scope != "sub/deeper" {
		t.Errorf("Children = %v", req.Children)
	}
}

func TestParseForkArgs_ModeForget(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "Calculate something", "scope": "."},
		},
		"mode": "forget",
	}
	req, err := ParseForkArgs(args, nil)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if req.Mode != ForkModeForget {
		t.Errorf("Mode = %q, want %q", req.Mode, ForkModeForget)
	}
}

func TestParseForkArgs_ModeRace(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "Try approach A", "scope": "."},
			map[string]any{"intent": "Try approach B", "scope": "."},
		},
		"mode": "race",
	}
	req, err := ParseForkArgs(args, nil)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if req.Mode != ForkModeRace {
		t.Errorf("Mode = %q, want %q", req.Mode, ForkModeRace)
	}
}

func TestParseForkArgs_AdoptWinnerRequiresRaceAndAdoptableChildren(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionRestore}}

	_, err := ParseForkArgs(map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "lane", "world": "subjective", "protection": "transactional"},
		},
		"mode":         "wait",
		"adopt_winner": true,
	}, cfg)
	if err == nil || !strings.Contains(err.Error(), "only valid with mode") {
		t.Fatalf("expected adopt_winner mode validation error, got %v", err)
	}

	_, err = ParseForkArgs(map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "lane", "world": "host", "protection": "none", "scope": "."},
		},
		"mode":         "race",
		"adopt_winner": true,
	}, cfg)
	if err == nil || !strings.Contains(err.Error(), "not an adoptable subjective child") {
		t.Fatalf("expected adoptable-child validation error, got %v", err)
	}

	req, err := ParseForkArgs(map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "lane", "world": "subjective", "protection": "transactional"},
		},
		"mode":         "race",
		"adopt_winner": true,
	}, cfg)
	if err != nil {
		t.Fatalf("ParseForkArgs() adopt_winner error = %v", err)
	}
	if !req.AdoptWinner {
		t.Fatal("expected adopt_winner=true")
	}
}

func TestParseForkArgs_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		errWord string
	}{
		{"MissingChildren", map[string]any{"mode": "wait"}, "children"},
		{"EmptyChildren", map[string]any{"children": []interface{}{}}, "at least one"},
		{"WrongChildrenType", map[string]any{"children": "single string"}, "array"},
		{"ChildNotObject", map[string]any{"children": []interface{}{"bad"}}, "object"},
		{"MissingIntent", map[string]any{"children": []interface{}{map[string]any{"scope": "."}}}, "intent"},
		{"MissingScope", map[string]any{"children": []interface{}{map[string]any{"intent": "ok"}}}, "scope"},
		{"EmptyIntent", map[string]any{"children": []interface{}{map[string]any{"intent": "", "scope": "."}}}, "intent"},
		{"EmptyScope", map[string]any{"children": []interface{}{map[string]any{"intent": "ok", "scope": ""}}}, "scope"},
		{"AbsoluteScope", map[string]any{"children": []interface{}{map[string]any{"intent": "ok", "scope": "/app/project"}}}, "relative path"},
		{"WrongModeType", map[string]any{"children": []interface{}{map[string]any{"intent": "task", "scope": "."}}, "mode": 42}, "string"},
		{"InvalidModeValue", map[string]any{"children": []interface{}{map[string]any{"intent": "task", "scope": "."}}, "mode": "yolo"}, "must be one of"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseForkArgs(tt.args, nil)
			if err == nil {
				t.Fatalf("expected error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.errWord) {
				t.Errorf("error should mention %q: %v", tt.errWord, err)
			}
		})
	}
}

func TestParseForkArgs_ForkWorldEnabled_DefaultSubjectiveTransactionalChild(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}}
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "inspect"},
		},
	}

	req, err := ParseForkArgs(args, cfg)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if len(req.Children) != 1 {
		t.Fatalf("got %d children, want 1", len(req.Children))
	}
	if req.Children[0].World != config.WorldSubjective {
		t.Fatalf("child world = %q, want subjective", req.Children[0].World)
	}
	if req.Children[0].Protection != config.ProtectionTransactional {
		t.Fatalf("child protection = %q, want transactional", req.Children[0].Protection)
	}
	if req.Children[0].Scope != "." {
		t.Fatalf("child scope = %q, want .", req.Children[0].Scope)
	}
}

func TestParseForkArgs_ForkWorldEnabled_HostChildRejectsScopedPath(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}}
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "inspect", "world": "host", "protection": "none", "scope": "subdir"},
		},
	}

	_, err := ParseForkArgs(args, cfg)
	if err == nil {
		t.Fatal("expected host child scope validation error")
	}
	if !strings.Contains(err.Error(), "only meaningful for world=\"subjective\"") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseForkArgs_ForkWorldEnabled_RejectsUnsupportedPairs(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}}
	cases := []map[string]any{
		{"intent": "inspect", "world": "host", "protection": "transactional", "scope": "."},
		{"intent": "inspect", "world": "subjective", "protection": "none", "scope": "."},
	}
	for _, child := range cases {
		_, err := ParseForkArgs(map[string]any{"children": []interface{}{child}}, cfg)
		if err == nil {
			t.Fatalf("expected unsupported pair error for %#v", child)
		}
		if !strings.Contains(err.Error(), "unsupported pair") {
			t.Fatalf("unexpected error for %#v: %v", child, err)
		}
	}
}

func TestParseForkArgs_ForkWorldEnabled_WorldRejectedWhenDisabled(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "inspect", "world": "subjective", "scope": "."},
		},
	}

	_, err := ParseForkArgs(args, &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: false}})
	if err == nil {
		t.Fatal("expected world gating error")
	}
	if !strings.Contains(err.Error(), "QUINE_FORK_WORLD_ENABLED=1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestForkChildEnvMaskSubsumesProcessIdentityFilter is the successor to
// TestFilterProcessIdentity. That function is gone: its hand-kept strip list and
// the registry's runtime-emitted class were the same eight names maintained in
// two places, and the mask in config.BoundaryBehavior is now derived from the
// registry instead. This test is the proof of subsumption — it plants the exact
// eight names the old filter listed, in the real process environ, and drives the
// real fork-boundary constructor.
//
// It also asserts what the old filter could NOT: QUINE_DEPTH is not merely
// stripped-and-reintroduced-by-synthesis, it is stamped from the parent's
// in-memory depth and wins over whatever the environ said.
func TestForkChildEnvMaskSubsumesProcessIdentityFilter(t *testing.T) {
	// The eight names the deleted filterProcessIdentity stripped by hand.
	stale := map[string]string{
		"QUINE_SESSION_ID":                 "old-session",
		"QUINE_RUN_ID":                     "old-run",
		"QUINE_TAPE_ID":                    "tape-session",
		ContextBootstrapEnv:                "bootstrap-dir",
		"QUINE_WORKSPACE_SESSION":          "workspace-session",
		"QUINE_WORKSPACE_OWNER":            "1",
		"QUINE_WORKSPACE_BOOTSTRAP":        "parent-workspace",
		"QUINE_WORKSPACE_CURRENT_REVISION": "wr4",
	}
	for name, value := range stale {
		t.Setenv(name, value)
	}
	t.Setenv("QUINE_DEPTH", "1")
	t.Setenv("FORK_MASK_FOREIGN", "kept")

	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "test-model", SessionID: "mask-parent", Depth: 2},
		Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"},
		Limits:    config.Limits{MaxDepth: 5, OutputTruncate: 20480},
		Paths:     config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}
	env := envSliceToMap(ForkChildEnv(cfg, nil))

	for name := range stale {
		if got, present := env[name]; present {
			t.Errorf("%s must not be inherited by a fork child (mask), got %q", name, got)
		}
	}

	// The parent's own environ still crosses the boundary untouched.
	if env["FORK_MASK_FOREIGN"] != "kept" {
		t.Errorf("foreign var = %q, want kept — the mask must only remove runtime-owned names", env["FORK_MASK_FOREIGN"])
	}
	if env["PATH"] == "" {
		t.Error("PATH missing from fork child env")
	}

	// Lineage: stamped from the parent's in-memory depth, not inherited.
	if env["QUINE_DEPTH"] != "3" {
		t.Errorf("QUINE_DEPTH = %q, want 3 (parent depth 2 + 1, stamped over the inherited 1)", env["QUINE_DEPTH"])
	}
	if env["QUINE_PARENT_SESSION"] != "mask-parent" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want mask-parent", env["QUINE_PARENT_SESSION"])
	}
}

// TestForkChildEnvOmitsUnsetKnobs is the fork-boundary half of the
// manufactured-evidence inversion (the sh half is TestShChildOmitsUnsetKnobs).
// The deleted ChildEnv serialized every registry knob into every child, defaults
// included. A knob nobody set is now ABSENT, and absence is what "compiled
// default" is spelled as.
func TestForkChildEnvOmitsUnsetKnobs(t *testing.T) {
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "test-model", SessionID: "unset-knob-parent"},
		Transport: config.Transport{APIKey: "test-key", Provider: "anthropic"},
		// Resolved values an operator never authored. None of them may appear.
		Limits:    config.Limits{MaxDepth: 5, MaxConcurrent: 20, ShTimeout: 10, OutputTruncate: 20480, MaxTurns: 20},
		ToolGates: config.ToolGates{},
		Paths:     config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}
	assertNoUnauthoredKnobs(t, "fork", envSliceToMap(ForkChildEnv(cfg, nil)), cfg.ForkChildStamps())
}

func TestCurrentRevisionFromWorldRevisionBlock(t *testing.T) {
	cases := map[string]string{
		"[WORLD REVISION] created=wr3 parent=wr1 current=wr3": "wr3",
		"[WORLD REVISION] current=wr2 (unchanged)":            "wr2",
		"[WORLD REVISION] wr4 -> wr1":                         "wr1",
		"":                                                    "",
	}
	for input, want := range cases {
		if got := currentRevisionFromWorldRevisionBlock(input); got != want {
			t.Fatalf("currentRevisionFromWorldRevisionBlock(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestForkExecutorChildCurrentWorldRevisionFallsBackToTape(t *testing.T) {
	tmpDir := t.TempDir()
	sessionID := "child-session"
	tapeDir := childTapeDir(tmpDir, nil, sessionID)
	if err := os.MkdirAll(tapeDir, 0o755); err != nil {
		t.Fatalf("mkdir tape dir: %v", err)
	}
	tapePath := filepath.Join(tapeDir, "0001.jsonl")
	tape := strings.Join([]string{
		`{"type":"message","data":{"role":"assistant","content":"...","tool_calls":[{"name":"sh"}]}}`,
		`{"type":"tool_result","data":{"content":"{\"tool\":\"sh\",\"mode\":\"sync\",\"status\":\"completed\",\"exit_code\":0,\"world_revision\":\"[WORLD REVISION] created=wr2 parent=wr1 current=wr2\"}"}}`,
		"",
	}, "\n")
	if err := os.WriteFile(tapePath, []byte(tape), 0o644); err != nil {
		t.Fatalf("write tape: %v", err)
	}

	f := &ForkExecutor{DataDir: tmpDir}
	got, err := f.childCurrentWorldRevision(sessionID)
	if err != nil {
		t.Fatalf("childCurrentWorldRevision() error = %v", err)
	}
	if got != "wr2" {
		t.Fatalf("childCurrentWorldRevision() = %q, want wr2", got)
	}
}

func TestForkExecutorChildCurrentWorldRevisionPrefersTapeOverSeedLedger(t *testing.T) {
	tmpDir := t.TempDir()
	sessionID := "child-session"

	stateDir := filepath.Join(tmpDir, "workspaces", sessionID)
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatalf("mkdir state dir: %v", err)
	}
	seedLedger := `{
  "current": "wr0",
  "next": 1,
  "revisions": {
    "wr0": {"id":"wr0","kind":"baseline"}
  }
}`
	if err := os.WriteFile(filepath.Join(stateDir, "world-revisions.json"), []byte(seedLedger), 0o644); err != nil {
		t.Fatalf("write seed ledger: %v", err)
	}

	tapeDir := childTapeDir(tmpDir, nil, sessionID)
	if err := os.MkdirAll(tapeDir, 0o755); err != nil {
		t.Fatalf("mkdir tape dir: %v", err)
	}
	tapePath := filepath.Join(tapeDir, "0001.jsonl")
	tape := strings.Join([]string{
		`{"type":"tool_result","data":{"content":"{\"tool\":\"sh\",\"mode\":\"sync\",\"status\":\"completed\",\"exit_code\":0,\"world_revision\":\"[WORLD REVISION] created=wr1 parent=wr0 current=wr1\"}"}}`,
		"",
	}, "\n")
	if err := os.WriteFile(tapePath, []byte(tape), 0o644); err != nil {
		t.Fatalf("write tape: %v", err)
	}

	f := &ForkExecutor{DataDir: tmpDir}
	got, err := f.childCurrentWorldRevision(sessionID)
	if err != nil {
		t.Fatalf("childCurrentWorldRevision() error = %v", err)
	}
	if got != "wr1" {
		t.Fatalf("childCurrentWorldRevision() = %q, want wr1", got)
	}
}

func TestForkExecutorResolveChildSessionIDFromPIDIndex(t *testing.T) {
	tmpDir := t.TempDir()
	sessionID := "20260406-120000_111_222"
	agentRoot := filepath.Join(tmpDir, "agent", sessionID)
	if err := os.MkdirAll(filepath.Join(tmpDir, "pid"), 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}
	if err := os.MkdirAll(agentRoot, 0o755); err != nil {
		t.Fatalf("mkdir agent root: %v", err)
	}
	if err := os.Symlink(agentRoot, filepath.Join(tmpDir, "pid", "222")); err != nil {
		t.Fatalf("symlink pid entry: %v", err)
	}

	f := &ForkExecutor{DataDir: tmpDir, SessionID: "parent-session"}
	if got := f.resolveChildSessionID(222); got != sessionID {
		t.Fatalf("resolveChildSessionID() = %q, want %q", got, sessionID)
	}
}

func TestForkExecutorResolveChildSessionIDFromPIDIndexPublicSurface(t *testing.T) {
	tmpDir := t.TempDir()
	sessionID := "20260406-120000_111_222"
	agentRoot := filepath.Join(tmpDir, "agent", sessionID)
	if err := os.MkdirAll(filepath.Join(tmpDir, "pid"), 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(agentRoot, "public"), 0o755); err != nil {
		t.Fatalf("mkdir public agent root: %v", err)
	}
	if err := os.Symlink(filepath.Join(agentRoot, "public"), filepath.Join(tmpDir, "pid", "222")); err != nil {
		t.Fatalf("symlink pid entry: %v", err)
	}

	f := &ForkExecutor{DataDir: tmpDir, SessionID: "parent-session"}
	if got := f.resolveChildSessionID(222); got != sessionID {
		t.Fatalf("resolveChildSessionID() = %q, want %q", got, sessionID)
	}
}

func TestForkExecutorResolveChildSessionIDFallsBackToAgentStatus(t *testing.T) {
	tmpDir := t.TempDir()
	sessionID := "20260406-120001_333_444"
	statusDir := filepath.Join(tmpDir, "agent", sessionID, "status")
	if err := os.MkdirAll(statusDir, 0o755); err != nil {
		t.Fatalf("mkdir status dir: %v", err)
	}
	status := map[string]any{
		"session_id":     sessionID,
		"parent_session": "parent-session",
		"pid":            444,
	}
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("marshal status: %v", err)
	}
	if err := os.WriteFile(filepath.Join(statusDir, "session.json"), data, 0o644); err != nil {
		t.Fatalf("write session.json: %v", err)
	}

	f := &ForkExecutor{DataDir: tmpDir, SessionID: "parent-session"}
	if got := f.resolveChildSessionID(444); got != sessionID {
		t.Fatalf("resolveChildSessionID() = %q, want %q", got, sessionID)
	}
}

func TestForkExecutor_CopyContextBootstrap(t *testing.T) {
	tmpDir := t.TempDir()
	contextRoot := filepath.Join(tmpDir, "context")
	frontierDir := filepath.Join(contextRoot, "state", "frontier")
	anchorDir := filepath.Join(contextRoot, "state", "anchors", "7.anchor")
	if err := os.MkdirAll(frontierDir, 0o755); err != nil {
		t.Fatalf("mkdir frontier: %v", err)
	}
	if err := os.MkdirAll(anchorDir, 0o755); err != nil {
		t.Fatalf("mkdir anchor dir: %v", err)
	}
	currentPath := filepath.Join(contextRoot, "state", "current.jsonl")
	currentContent := `{"type":"message","data":{"role":"assistant","content":"hello"}}
`
	if err := os.WriteFile(currentPath, []byte(currentContent), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}
	prepareForkPromptRoot(t, contextRoot)
	if err := os.WriteFile(filepath.Join(contextRoot, "prompt", "30-memory.md"), []byte("FORK_MEMORY_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write memory context: %v", err)
	}
	if err := os.WriteFile(filepath.Join(anchorDir, "meta.json"), []byte(`{"id":7,"resolution":"anchor-seven"}`), 0o644); err != nil {
		t.Fatalf("write anchor meta: %v", err)
	}
	if err := os.Symlink(filepath.Join("..", "anchors", "7.anchor"), filepath.Join(frontierDir, "7.link")); err != nil {
		t.Fatalf("symlink frontier link: %v", err)
	}

	f := &ForkExecutor{
		DataDir:     tmpDir,
		SessionID:   "test-session",
		ContextRoot: contextRoot,
	}

	bootstrapRoot, err := f.copyContextBootstrap(nil)
	if err != nil {
		t.Fatalf("copyContextBootstrap failed: %v", err)
	}
	defer os.RemoveAll(bootstrapRoot)

	childContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped current context: %v", err)
	}
	if string(childContent) != currentContent {
		t.Errorf("bootstrapped current content mismatch.\ngot:\n%s\nwant:\n%s", string(childContent), currentContent)
	}
	rawContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "bootstrap", "current.parent.raw.jsonl"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped raw current snapshot: %v", err)
	}
	if string(rawContent) != currentContent {
		t.Fatalf("bootstrapped raw current snapshot mismatch.\ngot:\n%s\nwant:\n%s", string(rawContent), currentContent)
	}
	memoryContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "prompt", "30-memory.md"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped memory context: %v", err)
	}
	if string(memoryContent) != "FORK_MEMORY_MARKER\n" {
		t.Fatalf("bootstrapped memory content mismatch: %q", string(memoryContent))
	}
	linkTarget, err := os.Readlink(filepath.Join(bootstrapRoot, "state", "frontier", "7.link"))
	if err != nil {
		t.Fatalf("read bootstrapped frontier link: %v", err)
	}
	if linkTarget != filepath.Join("..", "anchors", "7.anchor") {
		t.Fatalf("frontier link target = %q, want ../anchors/7.anchor", linkTarget)
	}
}

func TestForkExecutor_CopyContextBootstrap_ProjectsPendingToolBatchForChild(t *testing.T) {
	tmpDir := t.TempDir()
	contextRoot := filepath.Join(tmpDir, "context")
	prepareForkContextRoot(t, contextRoot)

	currentPath := filepath.Join(contextRoot, "state", "current.jsonl")
	currentContent := strings.Join([]string{
		`{"type":"message","data":{"role":"user","content":"Begin."}}`,
		`{"type":"message","data":{"role":"assistant","content":"","reasoning_content":"keep the remaining tool batch resumable","tool_calls":[{"id":"call_prev","name":"sh","arguments":{"command":"echo one"}},{"id":"call_fork","name":"fork","arguments":{"children":[{"intent":"child","scope":"lane"}]}},{"id":"call_after","name":"sh","arguments":{"command":"echo later"}}]}}`,
		`{"type":"tool_result","data":{"tool_id":"call_prev","content":{"status":"completed"},"is_error":false}}`,
		"",
	}, "\n")
	if err := os.WriteFile(currentPath, []byte(currentContent), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}

	f := &ForkExecutor{
		DataDir:     tmpDir,
		SessionID:   "test-session",
		ContextRoot: contextRoot,
	}

	bootstrapRoot, err := f.copyContextBootstrap(&childContextProjection{
		ForkToolID:    "call_fork",
		ParentMission: "parent mission",
		Child: ForkChild{
			Intent:     "child",
			World:      config.WorldSubjective,
			Protection: config.ProtectionTransactional,
			Scope:      "lane",
		},
	})
	if err != nil {
		t.Fatalf("copyContextBootstrap failed: %v", err)
	}
	defer os.RemoveAll(bootstrapRoot)

	childContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped current context: %v", err)
	}
	rawContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "bootstrap", "current.parent.raw.jsonl"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped raw current snapshot: %v", err)
	}
	if string(rawContent) != currentContent {
		t.Fatalf("bootstrapped raw current snapshot mismatch.\ngot:\n%s\nwant:\n%s", string(rawContent), currentContent)
	}
	assignmentContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "prompt", "45-fork-assignment.md"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped fork assignment prompt: %v", err)
	}
	assignment := string(assignmentContent)
	if !strings.Contains(assignment, "parent mission remains the active mission") {
		t.Fatalf("fork assignment should preserve parent mission semantics:\n%s", assignment)
	}
	if !strings.Contains(assignment, "child") {
		t.Fatalf("fork assignment should include child intent:\n%s", assignment)
	}

	entries := decodeContextEntries(t, childContent)
	if len(entries) != 4 {
		t.Fatalf("bootstrapped context entry count = %d, want 4\n%s", len(entries), string(childContent))
	}

	var projectedBatch tape.Message
	if err := json.Unmarshal(entries[1].Data, &projectedBatch); err != nil {
		t.Fatalf("decode projected tool batch: %v", err)
	}
	if projectedBatch.Role != tape.RoleAssistant {
		t.Fatalf("projected batch role = %q, want assistant", projectedBatch.Role)
	}
	if len(projectedBatch.ToolCalls) != 2 {
		t.Fatalf("projected batch tool call count = %d, want 2", len(projectedBatch.ToolCalls))
	}
	if projectedBatch.ReasoningContent != "keep the remaining tool batch resumable" {
		t.Fatalf("projected batch reasoning = %q, want preserved reasoning", projectedBatch.ReasoningContent)
	}
	if projectedBatch.ToolCalls[0].ID != "call_fork" || projectedBatch.ToolCalls[1].ID != "call_after" {
		t.Fatalf("projected batch tool ids = %#v, want [call_fork call_after]", []string{
			projectedBatch.ToolCalls[0].ID,
			projectedBatch.ToolCalls[1].ID,
		})
	}

	var forkBootstrap tape.ToolResult
	if err := json.Unmarshal(entries[2].Data, &forkBootstrap); err != nil {
		t.Fatalf("decode fork bootstrap result: %v", err)
	}
	if forkBootstrap.ToolID != "call_fork" {
		t.Fatalf("fork bootstrap ToolID = %q, want call_fork", forkBootstrap.ToolID)
	}
	forkPayload := decodeForkResultContent(t, forkBootstrap.Content)
	if got := forkPayload["status"]; got != "child_bootstrap" {
		t.Fatalf("fork bootstrap status = %v, want child_bootstrap", got)
	}
	if _, ok := forkPayload["prior_mission"]; ok {
		t.Fatalf("fork bootstrap should not expose prior_mission content, got %v", forkPayload["prior_mission"])
	}
	if got := forkPayload["parent_mission_ref"]; got != "context/prompt/40-mission.md" {
		t.Fatalf("fork bootstrap parent_mission_ref = %v, want context/prompt/40-mission.md", got)
	}
	if got := forkPayload["current_mission_ref"]; got != "context/prompt/40-mission.md" {
		t.Fatalf("fork bootstrap current_mission_ref = %v, want context/prompt/40-mission.md", got)
	}
	if got := forkPayload["child_assignment"]; got != "child" {
		t.Fatalf("fork bootstrap child_assignment = %v, want child", got)
	}
	if got := forkPayload["assignment_ref"]; got != "context/prompt/45-fork-assignment.md" {
		t.Fatalf("fork bootstrap assignment_ref = %v, want context/prompt/45-fork-assignment.md", got)
	}
	if _, ok := forkPayload["current_mission"]; ok {
		t.Fatalf("fork bootstrap should not replace current_mission with child intent, got %v", forkPayload["current_mission"])
	}
	if _, ok := forkPayload["projection_scope"]; ok {
		t.Fatalf("fork bootstrap projection_scope should be absent, got %v", forkPayload["projection_scope"])
	}
	if _, ok := forkPayload["assigned_child"]; ok {
		t.Fatalf("fork bootstrap assigned_child should be absent, got %v", forkPayload["assigned_child"])
	}

	var projectedAway tape.ToolResult
	if err := json.Unmarshal(entries[3].Data, &projectedAway); err != nil {
		t.Fatalf("decode projected-away result: %v", err)
	}
	if projectedAway.ToolID != "call_after" {
		t.Fatalf("projected-away ToolID = %q, want call_after", projectedAway.ToolID)
	}
	projectedPayload := decodeForkResultContent(t, projectedAway.Content)
	if got := projectedPayload["status"]; got != "projected_away" {
		t.Fatalf("projected-away status = %v, want projected_away", got)
	}
	if got := projectedPayload["reason"]; got != "outside_child_thread" {
		t.Fatalf("projected-away reason = %v, want outside_child_thread", got)
	}
}

func TestForkExecutor_CopyContextBootstrap_NoContext(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		DataDir:     tmpDir,
		SessionID:   "nonexistent-session",
		ContextRoot: filepath.Join(tmpDir, "nonexistent-context"),
	}

	bootstrapRoot, err := f.copyContextBootstrap(nil)
	if err != nil {
		t.Fatalf("copyContextBootstrap should not fail for nonexistent context: %v", err)
	}
	if bootstrapRoot != "" {
		t.Errorf("expected empty path for nonexistent context, got %q", bootstrapRoot)
	}
}

func TestForkExecutor_CopyContextBootstrap_ContextProjectionSymlink(t *testing.T) {
	tmpDir := t.TempDir()
	sessionRoot := filepath.Join(tmpDir, "agent", "session")
	incRoot := filepath.Join(sessionRoot, "inc")
	contextRoot := filepath.Join(incRoot, "0", "context")
	frontierDir := filepath.Join(contextRoot, "state", "frontier")
	anchorDir := filepath.Join(contextRoot, "state", "anchors", "7.anchor")
	if err := os.MkdirAll(frontierDir, 0o755); err != nil {
		t.Fatalf("mkdir frontier: %v", err)
	}
	if err := os.MkdirAll(anchorDir, 0o755); err != nil {
		t.Fatalf("mkdir anchor dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(contextRoot, "state", "current.jsonl"), []byte("{\"type\":\"message\",\"data\":{\"role\":\"assistant\",\"content\":\"hello\"}}\n"), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}
	if err := os.Symlink(filepath.Join("..", "anchors", "7.anchor"), filepath.Join(frontierDir, "7.link")); err != nil {
		t.Fatalf("symlink frontier link: %v", err)
	}
	if err := os.MkdirAll(incRoot, 0o755); err != nil {
		t.Fatalf("mkdir incarnation root: %v", err)
	}
	if err := os.Symlink("0", filepath.Join(incRoot, "current")); err != nil {
		t.Fatalf("symlink inc/current: %v", err)
	}
	if err := os.Symlink(filepath.Join("inc", "current", "context"), filepath.Join(sessionRoot, "context")); err != nil {
		t.Fatalf("symlink session context projection: %v", err)
	}

	f := &ForkExecutor{
		DataDir:     tmpDir,
		SessionID:   "test-session",
		ContextRoot: filepath.Join(sessionRoot, "context"),
	}

	bootstrapRoot, err := f.copyContextBootstrap(nil)
	if err != nil {
		t.Fatalf("copyContextBootstrap failed for projection symlink: %v", err)
	}
	defer os.RemoveAll(bootstrapRoot)

	childContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("failed to read bootstrapped current context: %v", err)
	}
	if !strings.Contains(string(childContent), "\"content\":\"hello\"") {
		t.Fatalf("bootstrapped current content mismatch:\n%s", string(childContent))
	}
	linkTarget, err := os.Readlink(filepath.Join(bootstrapRoot, "state", "frontier", "7.link"))
	if err != nil {
		t.Fatalf("read bootstrapped frontier link: %v", err)
	}
	if linkTarget != filepath.Join("..", "anchors", "7.anchor") {
		t.Fatalf("frontier link target = %q, want ../anchors/7.anchor", linkTarget)
	}
}

func TestForkExecutor_CopyContextBootstrap_StripsDanglingToolBatch(t *testing.T) {
	tmpDir := t.TempDir()
	contextRoot := filepath.Join(tmpDir, "context")
	prepareForkContextRoot(t, contextRoot)

	currentContent := strings.Join([]string{
		`{"type":"message","data":{"role":"system","content":"sys"}}`,
		`{"type":"message","data":{"role":"user","content":"user"}}`,
		`{"type":"message","data":{"role":"assistant","content":"","tool_calls":[{"id":"call_sh","name":"sh","arguments":{"command":"pwd"}},{"id":"call_fork","name":"fork","arguments":{"children":[{"intent":"child"}]}}]}}`,
		`{"type":"tool_result","data":{"tool_id":"call_sh","content":"{\"status\":\"completed\"}","is_error":false}}`,
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(contextRoot, "state", "current.jsonl"), []byte(currentContent), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}

	f := &ForkExecutor{
		DataDir:     tmpDir,
		SessionID:   "test-session",
		ContextRoot: contextRoot,
	}

	bootstrapRoot, err := f.copyContextBootstrap(nil)
	if err != nil {
		t.Fatalf("copyContextBootstrap failed: %v", err)
	}
	defer os.RemoveAll(bootstrapRoot)

	childContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("read bootstrapped current context: %v", err)
	}
	got := string(childContent)
	if strings.Contains(got, `"tool_calls"`) || strings.Contains(got, `call_sh`) || strings.Contains(got, `call_fork`) {
		t.Fatalf("bootstrapped current retained dangling tool batch:\n%s", got)
	}
	if !strings.Contains(got, `"role":"system"`) || !strings.Contains(got, `"role":"user"`) {
		t.Fatalf("bootstrapped current lost stable prefix:\n%s", got)
	}
}

func TestForkExecutor_CopyContextBootstrap_StripsToolMechanicsButKeepsAssistantText(t *testing.T) {
	tmpDir := t.TempDir()
	contextRoot := filepath.Join(tmpDir, "context")
	prepareForkContextRoot(t, contextRoot)

	currentContent := strings.Join([]string{
		`{"type":"message","data":{"role":"system","content":"sys"}}`,
		`{"type":"message","data":{"role":"assistant","content":"checked pwd","tool_calls":[{"id":"call_sh","name":"sh","arguments":{"command":"pwd"}}]}}`,
		`{"type":"tool_result","data":{"tool_id":"call_sh","content":"{\"status\":\"completed\"}","is_error":false}}`,
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(contextRoot, "state", "current.jsonl"), []byte(currentContent), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}

	f := &ForkExecutor{
		DataDir:     tmpDir,
		SessionID:   "test-session",
		ContextRoot: contextRoot,
	}

	bootstrapRoot, err := f.copyContextBootstrap(nil)
	if err != nil {
		t.Fatalf("copyContextBootstrap failed: %v", err)
	}
	defer os.RemoveAll(bootstrapRoot)

	childContent, err := os.ReadFile(filepath.Join(bootstrapRoot, "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("read bootstrapped current context: %v", err)
	}
	got := string(childContent)
	if strings.Contains(got, `"tool_calls"`) || strings.Contains(got, `"tool_id":"call_sh"`) {
		t.Fatalf("bootstrapped current should strip tool mechanics:\n%s", got)
	}
	if !strings.Contains(got, `"content":"checked pwd"`) {
		t.Fatalf("bootstrapped current should preserve assistant text:\n%s", got)
	}
}

func TestForkExecutor_ResolveChildWorkspace_InheritsCurrentWorkspace(t *testing.T) {
	root := t.TempDir()
	current := filepath.Join(root, "subdir")
	if err := os.MkdirAll(current, 0o755); err != nil {
		t.Fatalf("mkdir current workspace: %v", err)
	}

	rootCanon, err := canonicalizeRequestedPath(root)
	if err != nil {
		t.Fatalf("canonicalize root: %v", err)
	}
	currentCanon, err := canonicalizeRequestedPath(current)
	if err != nil {
		t.Fatalf("canonicalize current workspace: %v", err)
	}
	f := &ForkExecutor{
		WorkspaceEnabled: true,
		WorkspaceRoot:    rootCanon,
		Workspace:        currentCanon,
	}

	got, err := f.resolveChildWorkspace(".")
	if err != nil {
		t.Fatalf("resolveChildWorkspace(.) failed: %v", err)
	}
	if got != currentCanon {
		t.Fatalf("resolveChildWorkspace(.) = %q, want %q", got, currentCanon)
	}
}

func TestForkExecutor_ResolveChildWorkspace_RelativeToCurrentWorkspace(t *testing.T) {
	root := t.TempDir()
	current := filepath.Join(root, "subdir")
	child := filepath.Join(current, "nested")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir nested workspace: %v", err)
	}

	rootCanon, err := canonicalizeRequestedPath(root)
	if err != nil {
		t.Fatalf("canonicalize root: %v", err)
	}
	currentCanon, err := canonicalizeRequestedPath(current)
	if err != nil {
		t.Fatalf("canonicalize current workspace: %v", err)
	}
	f := &ForkExecutor{
		WorkspaceEnabled: true,
		WorkspaceRoot:    rootCanon,
		Workspace:        currentCanon,
	}

	want, err := canonicalizeRequestedPath(child)
	if err != nil {
		t.Fatalf("canonicalize nested workspace: %v", err)
	}
	got, err := f.resolveChildWorkspace("nested")
	if err != nil {
		t.Fatalf("resolveChildWorkspace(nested) failed: %v", err)
	}
	if got != want {
		t.Fatalf("resolveChildWorkspace(nested) = %q, want %q", got, want)
	}
}

func TestForkExecutor_ResolveChildWorkspace_RejectsAbsolutePath(t *testing.T) {
	root := t.TempDir()
	current := filepath.Join(root, "subdir")
	child := filepath.Join(root, "other")
	if err := os.MkdirAll(current, 0o755); err != nil {
		t.Fatalf("mkdir current workspace: %v", err)
	}
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir child workspace: %v", err)
	}

	rootCanon, err := canonicalizeRequestedPath(root)
	if err != nil {
		t.Fatalf("canonicalize root: %v", err)
	}
	currentCanon, err := canonicalizeRequestedPath(current)
	if err != nil {
		t.Fatalf("canonicalize current workspace: %v", err)
	}
	f := &ForkExecutor{
		WorkspaceEnabled: true,
		WorkspaceRoot:    rootCanon,
		Workspace:        currentCanon,
	}

	_, err = f.resolveChildWorkspace(child)
	if err == nil {
		t.Fatalf("resolveChildWorkspace should reject absolute paths")
	}
	if !strings.Contains(err.Error(), "absolute paths are invalid") {
		t.Fatalf("resolveChildWorkspace(abs) error = %q, want absolute-path remediation", err)
	}
}

func TestForkExecutor_ResolveChildWorkspace_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	current := filepath.Join(root, "subdir", "deeper")
	if err := os.MkdirAll(current, 0o755); err != nil {
		t.Fatalf("mkdir current workspace: %v", err)
	}

	rootCanon, err := canonicalizeRequestedPath(root)
	if err != nil {
		t.Fatalf("canonicalize root: %v", err)
	}
	currentCanon, err := canonicalizeRequestedPath(current)
	if err != nil {
		t.Fatalf("canonicalize current workspace: %v", err)
	}
	f := &ForkExecutor{
		WorkspaceEnabled: true,
		WorkspaceRoot:    rootCanon,
		Workspace:        currentCanon,
	}

	_, err = f.resolveChildWorkspace("../../../../escape")
	if err == nil {
		t.Fatalf("resolveChildWorkspace should reject escapes")
	}
	if !strings.Contains(err.Error(), "use \".\" or a relative subpath") {
		t.Fatalf("resolveChildWorkspace escape error = %q, want remediation hint", err)
	}
}

func TestForkExecutor_Truncate(t *testing.T) {
	f := &ForkExecutor{MaxOutput: 100}

	// Short content - no truncation
	short := "hello world"
	if result := f.truncate([]byte(short)); result != short {
		t.Errorf("truncate(%q) = %q, want %q", short, result, short)
	}

	// Long content - should truncate
	long := strings.Repeat("A", 200)
	result := f.truncate([]byte(long))
	if !strings.Contains(result, "...[Output Truncated,") {
		t.Errorf("truncate should add truncation notice, got: %s", result)
	}
	if !strings.Contains(result, "200 bytes total]") {
		t.Errorf("truncate should show total bytes, got: %s", result)
	}
}

// Integration test - requires actual quine binary
func TestForkExecutor_Execute_MissingBinary(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		QuinePath:      "/nonexistent/quine",
		DataDir:        tmpDir,
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}
	prepareForkContextRoot(t, f.ContextRoot)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "test intent", Scope: "."}},
		Mode:     ForkModeWait,
	}

	result := f.Execute("tool-1", req)
	if !result.IsError {
		t.Errorf("expected error for missing binary")
	}
	payload := decodeForkResultContent(t, result.Content)
	if payload["tool"] != "fork" {
		t.Fatalf("tool = %#v, want fork", payload["tool"])
	}
	if payload["status"] != "error" {
		t.Fatalf("status = %#v, want error", payload["status"])
	}
	if resultInt(t, payload, "requested") != 1 {
		t.Fatalf("requested = %#v, want 1", payload["requested"])
	}
	errors, ok := payload["errors"].([]any)
	if !ok || len(errors) != 1 {
		t.Fatalf("errors = %#v, want single entry", payload["errors"])
	}
	if !strings.Contains(errors[0].(string), "no such file or directory") {
		t.Fatalf("unexpected fork error payload: %#v", errors[0])
	}
	if len(result.StructuredContent) != 0 {
		t.Fatal("expected structured content to be empty for missing-binary fork result")
	}
}

func TestForkExecutor_Execute_AsyncMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Use a real command that will run briefly
	f := &ForkExecutor{
		QuinePath:      "/bin/sleep", // Will fail but that's ok for async
		DataDir:        tmpDir,
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}
	prepareForkContextRoot(t, f.ContextRoot)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "0.1", Scope: "."}}, // sleep argument
		Mode:     ForkModeForget,
	}

	start := time.Now()
	result := f.Execute("tool-1", req)
	elapsed := time.Since(start)

	// Async should return immediately (not wait for child)
	if elapsed > 2*time.Second {
		t.Errorf("async fork took too long: %v", elapsed)
	}

	// Result should indicate children were spawned
	if result.IsError {
		// It's ok if it fails to start, but shouldn't take long
		t.Logf("async fork error (expected for sleep command): %s", result.Content)
	} else {
		payload := decodeForkResultContent(t, result.Content)
		if payload["status"] != "spawned" {
			t.Errorf("expected spawned status in result, got: %s", result.Content)
		}
		children, ok := payload["children"].([]any)
		if !ok || len(children) != 1 {
			t.Fatalf("children = %#v, want single child", payload["children"])
		}
		child, ok := children[0].(map[string]any)
		if !ok {
			t.Fatalf("child payload = %#v, want object", children[0])
		}
		if child["scope"] != "." {
			t.Fatalf("child scope = %#v, want \".\"", child["scope"])
		}
	}
}

func TestForkExecutor_Execute_GatherAll_MissingBinary(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		QuinePath:      "/nonexistent/quine",
		DataDir:        tmpDir,
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}
	prepareForkContextRoot(t, f.ContextRoot)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "task A", Scope: "."}, {Intent: "task B", Scope: "."}},
		Mode:     ForkModeWait,
	}

	result := f.Execute("tool-1", req)
	if !result.IsError {
		t.Errorf("expected error for missing binary")
	}
	payload := decodeForkResultContent(t, result.Content)
	if payload["status"] != "error" {
		t.Fatalf("status = %#v, want error", payload["status"])
	}
	if resultInt(t, payload, "requested") != 2 {
		t.Fatalf("requested = %#v, want 2", payload["requested"])
	}
	errors, ok := payload["errors"].([]any)
	if !ok || len(errors) != 2 {
		t.Fatalf("errors = %#v, want two entries", payload["errors"])
	}
	for _, raw := range errors {
		if !strings.Contains(raw.(string), "no such file or directory") {
			t.Fatalf("unexpected fork error payload: %#v", raw)
		}
	}
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 2 {
		t.Fatalf("children = %#v, want two entries", payload["children"])
	}
	for i, raw := range children {
		child, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("child %d payload = %#v, want object", i, raw)
		}
		if got := resultString(t, child, "error"); !strings.Contains(got, "no such file or directory") {
			t.Fatalf("child %d error = %q, want concrete spawn failure", i, got)
		}
	}
	if len(result.StructuredContent) != 0 {
		t.Fatal("expected structured content to be empty for missing-binary fork result")
	}
}

func TestForkExecutor_Execute_GatherAll_TimeoutReturnsPartialResults(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "fork-child.sh")
	script := `#!/bin/sh
case "$1" in
  fast)
    printf 'fast child done\n'
    exit 0
    ;;
  slow)
    sleep 10
    ;;
esac
`
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:      helper,
		DataDir:        filepath.Join(tmpDir, "runtime"),
		SessionID:      "wait-timeout-parent",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 100 * time.Millisecond,
		MaxOutput:      10000,
		Env:            []string{},
		WorkDir:        tmpDir,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	start := time.Now()
	result := f.Execute("tool-wait-timeout", ForkRequest{
		Children: []ForkChild{{Intent: "fast", Scope: "."}, {Intent: "slow", Scope: "."}},
		Mode:     ForkModeWait,
	})
	elapsed := time.Since(start)
	if elapsed > 3*time.Second {
		t.Fatalf("fork wait timeout took %s, want bounded return", elapsed)
	}
	if result.IsError {
		t.Fatalf("expected partial success not to be an error, got:\n%s", result.Content)
	}

	payload := decodeForkResultContent(t, result.Content)
	if payload["status"] != "timeout" {
		t.Fatalf("status = %#v, want timeout", payload["status"])
	}
	if resultInt(t, payload, "succeeded") != 1 {
		t.Fatalf("succeeded = %#v, want 1", payload["succeeded"])
	}
	if resultInt(t, payload, "killed") != 1 {
		t.Fatalf("killed = %#v, want 1", payload["killed"])
	}
	requireStableRelationCounters(t, payload, 2, 2, 1, 1)
	relationResult := requireRelationResult(t, payload)
	requireStableRelationCounters(t, relationResult, 2, 2, 1, 1)
	status := requireRelationStatus(t, payload)
	requireStableRelationCounters(t, status, 2, 2, 1, 1)
	children := payload["children"].([]any)
	if children[0].(map[string]any)["status"] != "completed" {
		t.Fatalf("fast child status = %#v", children[0])
	}
	if children[1].(map[string]any)["status"] != "timeout" {
		t.Fatalf("slow child status = %#v", children[1])
	}
}

func TestForkExecutor_Execute_Race_MissingBinary(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		QuinePath:      "/nonexistent/quine",
		DataDir:        tmpDir,
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}
	prepareForkContextRoot(t, f.ContextRoot)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "approach A", Scope: "."}, {Intent: "approach B", Scope: "."}},
		Mode:     ForkModeRace,
	}

	result := f.Execute("tool-1", req)
	if !result.IsError {
		t.Errorf("expected error for missing binary")
	}
	payload := decodeForkResultContent(t, result.Content)
	if payload["status"] != "error" {
		t.Fatalf("status = %#v, want error", payload["status"])
	}
	if payload["mode"] != ForkModeRace {
		t.Fatalf("mode = %#v, want %q", payload["mode"], ForkModeRace)
	}
	if resultInt(t, payload, "requested") != 2 {
		t.Fatalf("requested = %#v, want 2", payload["requested"])
	}
	errors, ok := payload["errors"].([]any)
	if !ok || len(errors) != 2 {
		t.Fatalf("errors = %#v, want two entries", payload["errors"])
	}
	for _, raw := range errors {
		if !strings.Contains(raw.(string), "no such file or directory") {
			t.Fatalf("unexpected fork error payload: %#v", raw)
		}
	}
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 2 {
		t.Fatalf("children = %#v, want two entries", payload["children"])
	}
	for i, raw := range children {
		child, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("child %d payload = %#v, want object", i, raw)
		}
		if got := resultString(t, child, "error"); !strings.Contains(got, "no such file or directory") {
			t.Fatalf("child %d error = %q, want concrete spawn failure", i, got)
		}
	}
	if len(result.StructuredContent) != 0 {
		t.Fatal("expected structured content to be empty for missing-binary fork result")
	}
}

func TestForkExecutor_Execute_Race_TimeoutReturnsWithoutWinner(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "slow-child.sh")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\nsleep 10\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:      helper,
		DataDir:        filepath.Join(tmpDir, "runtime"),
		SessionID:      "race-timeout-parent",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 100 * time.Millisecond,
		MaxOutput:      10000,
		Env:            []string{},
		WorkDir:        tmpDir,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	start := time.Now()
	result := f.Execute("tool-race-timeout", ForkRequest{
		Children: []ForkChild{{Intent: "a", Scope: "."}, {Intent: "b", Scope: "."}},
		Mode:     ForkModeRace,
	})
	elapsed := time.Since(start)
	if elapsed > 3*time.Second {
		t.Fatalf("fork race timeout took %s, want bounded return", elapsed)
	}
	if !result.IsError {
		t.Fatalf("expected timed out race with no winner to be an error")
	}

	payload := decodeForkResultContent(t, result.Content)
	if payload["status"] != "timeout" {
		t.Fatalf("status = %#v, want timeout", payload["status"])
	}
	if resultInt(t, payload, "killed") != 2 {
		t.Fatalf("killed = %#v, want 2", payload["killed"])
	}
	requireStableRelationCounters(t, payload, 2, 2, 0, 2)
	relationResult := requireRelationResult(t, payload)
	requireStableRelationCounters(t, relationResult, 2, 2, 0, 2)
	status := requireRelationStatus(t, payload)
	requireStableRelationCounters(t, status, 2, 2, 0, 2)
	if payload["winner"] != nil {
		t.Fatalf("winner = %#v, want nil", payload["winner"])
	}
	children := payload["children"].([]any)
	for i, raw := range children {
		if raw.(map[string]any)["status"] != "timeout" {
			t.Fatalf("child %d status = %#v, want timeout", i, raw)
		}
	}
}

func TestForkExecutor_Execute_RaceWinnerPreservesSpawnErrors(t *testing.T) {
	tmpDir := t.TempDir()
	workspace := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
	stub := writeSpawnRaceStubQuine(t, tmpDir)
	contextRoot := filepath.Join(tmpDir, "context")
	prepareForkContextRoot(t, contextRoot)

	f := &ForkExecutor{
		QuinePath:        stub,
		DataDir:          filepath.Join(tmpDir, "runtime"),
		SessionID:        "race-spawn-errors-parent",
		ContextRoot:      contextRoot,
		DefaultTimeout:   time.Second,
		MaxOutput:        10000,
		Env:              []string{},
		WorkspaceEnabled: true,
		WorkspaceRoot:    workspace,
		Workspace:        workspace,
		WorkspaceBackend: "direct",
		WorkDir:          workspace,
	}

	result := f.Execute("tool-race-spawn-errors", ForkRequest{
		Children: []ForkChild{
			{Intent: "bad scope", Scope: "../outside"},
			{Intent: "fast child", Scope: "."},
		},
		Mode: ForkModeRace,
	})
	if result.IsError {
		t.Fatalf("fork race should still return the successful winner: %s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	winner, ok := payload["winner"].(map[string]any)
	if !ok {
		t.Fatalf("winner missing from payload: %#v", payload)
	}
	if got := resultInt(t, winner, "index"); got != 1 {
		t.Fatalf("winner index = %d, want 1", got)
	}
	errors, ok := payload["errors"].([]any)
	if !ok || len(errors) != 1 {
		t.Fatalf("errors = %#v, want one retained spawn error", payload["errors"])
	}
	if got := fmt.Sprint(errors[0]); !strings.Contains(got, "resolve child 0 scope") {
		t.Fatalf("spawn error = %q, want child scope failure", got)
	}
	requireStableRelationCounters(t, payload, 1, 1, 1, 0)
	relationResult := requireRelationResult(t, payload)
	requireStableRelationCounters(t, relationResult, 1, 1, 1, 0)
	status := requireRelationStatus(t, payload)
	requireStableRelationCounters(t, status, 1, 1, 1, 0)
	statusErrors, ok := status["errors"].([]any)
	if !ok || len(statusErrors) != 1 {
		t.Fatalf("relation errors = %#v, want one retained spawn error", status["errors"])
	}
}

func TestForkExecutor_Execute_MultipleAsync(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		QuinePath:      "/bin/sleep",
		DataDir:        tmpDir,
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}
	prepareForkContextRoot(t, f.ContextRoot)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "0.1", Scope: "."}, {Intent: "0.1", Scope: "."}, {Intent: "0.1", Scope: "."}},
		Mode:     ForkModeForget,
	}

	start := time.Now()
	result := f.Execute("tool-1", req)
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("async fork with 3 children took too long: %v", elapsed)
	}

	if result.IsError {
		t.Logf("async fork error: %s", result.Content)
	} else {
		payload := decodeForkResultContent(t, result.Content)
		if got := payload["spawned"]; got != json.Number("3") && got != float64(3) {
			t.Errorf("expected spawned=3 in result, got: %s", result.Content)
		}
	}
}

func TestForkExecutor_Execute_Forget_WritesRelationAndSeedSurface(t *testing.T) {
	tmpDir := t.TempDir()
	retainedRoot := filepath.Join(tmpDir, "retained")
	contextRoot := filepath.Join(tmpDir, "context")
	prepareForkContextRoot(t, contextRoot)
	currentContent := `{"type":"message","data":{"role":"assistant","content":"hello seed"}}` + "\n"
	if err := os.WriteFile(filepath.Join(contextRoot, "state", "current.jsonl"), []byte(currentContent), 0o644); err != nil {
		t.Fatalf("write current context: %v", err)
	}

	helper := filepath.Join(tmpDir, "sleep.sh")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\nsleep 0.2\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:   helper,
		DataDir:     filepath.Join(tmpDir, "runtime"),
		SessionID:   "parent-session",
		ContextRoot: contextRoot,
		Env:         []string{"QUINE_RETENTION_DIR=" + retainedRoot},
		WorkDir:     tmpDir,
	}

	result := f.Execute("tool:1/relation", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Scope: "."}},
		Mode:     ForkModeForget,
	})
	if result.IsError {
		t.Fatalf("expected forget fork to succeed, got:\n%s", result.Content)
	}

	payload := decodeForkResultContent(t, result.Content)
	if got := payload["relation_id"]; got != "tool_1_relation" {
		t.Fatalf("relation_id = %#v, want tool_1_relation", got)
	}
	relationRoot, ok := payload["relation_root"].(string)
	if !ok || relationRoot == "" {
		t.Fatalf("relation_root missing from payload: %#v", payload["relation_root"])
	}
	wantRelationRoot := filepath.Join(retainedRoot, "sessions", "parent-session", "relations", "tool_1_relation")
	if relationRoot != wantRelationRoot {
		t.Fatalf("relation_root = %q, want %q", relationRoot, wantRelationRoot)
	}
	if _, err := os.Stat(filepath.Join(relationRoot, "relation.json")); err != nil {
		t.Fatalf("stat relation.json: %v", err)
	}
	if _, err := os.Stat(filepath.Join(relationRoot, "status.json")); err != nil {
		t.Fatalf("stat status.json: %v", err)
	}
	if _, err := os.Stat(filepath.Join(relationRoot, "result.json")); err != nil {
		t.Fatalf("stat result.json: %v", err)
	}
	requireStableRelationCounters(t, payload, 1, 0, 0, 0)
	relationResult := requireRelationResult(t, payload)
	requireStableRelationCounters(t, relationResult, 1, 0, 0, 0)
	status := requireRelationStatus(t, payload)
	requireStableRelationCounters(t, status, 1, 0, 0, 0)
	if _, err := os.Stat(filepath.Join(relationRoot, "members", "000.json")); err != nil {
		t.Fatalf("stat member surface: %v", err)
	}

	children, ok := payload["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("children = %#v, want single child", payload["children"])
	}
	child, ok := children[0].(map[string]any)
	if !ok {
		t.Fatalf("child payload = %#v, want object", children[0])
	}
	childSession, ok := child["session_id"].(string)
	if !ok || childSession == "" {
		t.Fatalf("child session_id missing: %#v", child["session_id"])
	}
	wantRetainedRoot := filepath.Join(retainedRoot, "sessions", childSession)
	if got := child["retained_root"]; got != wantRetainedRoot {
		t.Fatalf("retained_root = %#v, want %q", got, wantRetainedRoot)
	}
	wantSeedRoot := filepath.Join(wantRetainedRoot, "seed")
	if got := child["seed_root"]; got != wantSeedRoot {
		t.Fatalf("seed_root = %#v, want %q", got, wantSeedRoot)
	}

	seedCurrent, err := os.ReadFile(filepath.Join(wantSeedRoot, "context", "state", "current.jsonl"))
	if err != nil {
		t.Fatalf("read seed current context: %v", err)
	}
	if string(seedCurrent) != currentContent {
		t.Fatalf("seed current context mismatch.\ngot:\n%s\nwant:\n%s", string(seedCurrent), currentContent)
	}
	originData, err := os.ReadFile(filepath.Join(wantSeedRoot, "origin.json"))
	if err != nil {
		t.Fatalf("read seed origin: %v", err)
	}
	var origin map[string]any
	if err := json.Unmarshal(originData, &origin); err != nil {
		t.Fatalf("decode seed origin: %v", err)
	}
	if got := origin["relation_id"]; got != "tool_1_relation" {
		t.Fatalf("origin relation_id = %#v, want tool_1_relation", got)
	}
	if got := origin["initiator_session"]; got != "parent-session" {
		t.Fatalf("origin initiator_session = %#v, want parent-session", got)
	}
	if got := origin["intent"]; got != "ignored" {
		t.Fatalf("origin intent = %#v, want ignored", got)
	}
}

func TestForkExecutor_Execute_WaitPreservesParentMissionAsChildArgv(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "print-argv.sh")
	script := "#!/bin/sh\nprintf 'ARG1=%s\\n' \"$1\"\n"
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:      helper,
		DataDir:        filepath.Join(tmpDir, "runtime"),
		SessionID:      "parent-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		Env:            []string{},
		WorkDir:        tmpDir,
		MaxOutput:      10000,
		DefaultTimeout: 2 * time.Second,
		Mission:        "parent mission with task URL and /app/solution.txt",
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-parent-mission", ForkRequest{
		Children: []ForkChild{{Intent: "child lane assignment only", Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected wait fork to succeed, got:\n%s", result.Content)
	}

	payload := decodeForkResultContent(t, result.Content)
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("children = %#v, want single child", payload["children"])
	}
	child, ok := children[0].(map[string]any)
	if !ok {
		t.Fatalf("child payload = %#v, want object", children[0])
	}
	stdout, _ := child["stdout"].(string)
	if !strings.Contains(stdout, "ARG1=parent mission with task URL and /app/solution.txt") {
		t.Fatalf("child argv should preserve parent mission, got:\n%s", stdout)
	}
	if strings.Contains(stdout, "child lane assignment only") {
		t.Fatalf("child argv should not replace mission with lane assignment, got:\n%s", stdout)
	}
}

func TestForkExecutor_Execute_Wait_RecordsChildProcessHandles(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "print-session.sh")
	script := "#!/bin/sh\nprintf 'SESSION=%s\\n' \"$QUINE_SESSION_ID\"\n"
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:   helper,
		DataDir:     filepath.Join(tmpDir, "runtime"),
		SessionID:   "parent-session",
		ContextRoot: filepath.Join(tmpDir, "context"),
		Env:         []string{},
		WorkDir:     tmpDir,
		MaxOutput:   10000,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-handles", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected wait fork to succeed, got:\n%s", result.Content)
	}

	payload := decodeForkResultContent(t, result.Content)
	if got := payload["relation_id"]; got != "tool-handles" {
		t.Fatalf("relation_id = %#v, want tool-handles", got)
	}
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("children = %#v, want single child", payload["children"])
	}
	child, ok := children[0].(map[string]any)
	if !ok {
		t.Fatalf("child payload = %#v, want object", children[0])
	}
	sessionID, ok := child["session_id"].(string)
	if !ok || sessionID == "" {
		t.Fatalf("child session_id missing: %#v", child["session_id"])
	}
	stdout, _ := child["stdout"].(string)
	if !strings.Contains(stdout, "SESSION="+sessionID) {
		t.Fatalf("stdout should include injected QUINE_SESSION_ID %q, got:\n%s", sessionID, stdout)
	}
	wantAgentRoot := filepath.Join(f.DataDir, "agent", sessionID)
	if got := child["agent_root"]; got != wantAgentRoot {
		t.Fatalf("agent_root = %#v, want %q", got, wantAgentRoot)
	}
	wantPublicRoot := filepath.Join(wantAgentRoot, "public")
	if got := child["public_root"]; got != wantPublicRoot {
		t.Fatalf("public_root = %#v, want %q", got, wantPublicRoot)
	}
	if got := child["status_path"]; got != filepath.Join(wantPublicRoot, "status", "session.json") {
		t.Fatalf("status_path = %#v", got)
	}
	if got := child["control_path"]; got != filepath.Join(wantPublicRoot, "ctl") {
		t.Fatalf("control_path = %#v", got)
	}

	relationRoot, _ := payload["relation_root"].(string)
	resultData, err := os.ReadFile(filepath.Join(relationRoot, "result.json"))
	if err != nil {
		t.Fatalf("read relation result: %v", err)
	}
	if !strings.Contains(string(resultData), sessionID) {
		t.Fatalf("relation result should retain child session handle, got:\n%s", string(resultData))
	}
}

func TestForkExecutor_ChildrenDoNotInheritLiveSideChannelFDs(t *testing.T) {
	tmpDir := t.TempDir()

	helper := filepath.Join(tmpDir, "fork-helper.sh")
	script := `#!/bin/sh
if cat /dev/fd/3 >/dev/null 2>/tmp/fork-fd3.err; then
  echo HAS_FD3
else
  echo NO_FD3
fi
`
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:      helper,
		DataDir:        tmpDir,
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Scope: "."}},
		Mode:     ForkModeWait,
	})

	if result.IsError {
		t.Fatalf("expected child helper to succeed, got error:\n%s", result.Content)
	}
	if !strings.Contains(string(result.Content), "NO_FD3") {
		t.Fatalf("expected child output to report missing fd 3, got:\n%s", result.Content)
	}
	if strings.Contains(string(result.Content), "HAS_FD3") {
		t.Fatalf("child unexpectedly inherited live fd 3:\n%s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	if payload["tool"] != "fork" {
		t.Fatalf("tool = %#v, want fork", payload["tool"])
	}
	if payload["mode"] != "wait" {
		t.Fatalf("mode = %#v, want wait", payload["mode"])
	}
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("children = %#v, want single child", payload["children"])
	}
}

func TestForkExecutor_WaitDoesNotDependOnGrandchildStdoutEOF(t *testing.T) {
	tmpDir := t.TempDir()

	helper := filepath.Join(tmpDir, "grandchild-stdout-holder.sh")
	script := `#!/bin/sh
(
  sleep 2
  printf 'late grandchild output\n'
) &
echo "$!" > bg.pid
printf 'direct child output\n'
exit 0
`
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:      helper,
		DataDir:        filepath.Join(tmpDir, "runtime"),
		SessionID:      "test-session",
		ContextRoot:    filepath.Join(tmpDir, "context"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
		WorkDir:        tmpDir,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	start := time.Now()
	result := f.Execute("tool-stdout-grandchild", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Scope: "."}},
		Mode:     ForkModeWait,
	})
	elapsed := time.Since(start)

	if data, err := os.ReadFile(filepath.Join(tmpDir, "bg.pid")); err == nil {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			if proc, err := os.FindProcess(pid); err == nil {
				_ = proc.Kill()
			}
		}
	}

	if result.IsError {
		t.Fatalf("expected wait fork to succeed, got:\n%s", result.Content)
	}
	if elapsed >= time.Second {
		t.Fatalf("fork wait took %s; it likely waited for inherited stdout EOF from a grandchild", elapsed)
	}
	payload := decodeForkResultContent(t, result.Content)
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("children = %#v, want single child", payload["children"])
	}
	child := children[0].(map[string]any)
	stdout, _ := child["stdout"].(string)
	if !strings.Contains(stdout, "direct child output") {
		t.Fatalf("stdout missing direct child output:\n%s", stdout)
	}
}

func TestForkExecutor_ForkWorldEnabled_HostChildFiltersWorkspacePhysics(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "print-env.sh")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\nenv | sort\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:        helper,
		DataDir:          tmpDir,
		SessionID:        "test-session",
		ContextRoot:      filepath.Join(tmpDir, "context"),
		DefaultTimeout:   5 * time.Second,
		MaxOutput:        10000,
		Env:              []string{"QUINE_WORKSPACE_ROOT=/tmp/root", "QUINE_WORKSPACE=/tmp/root/app", "QUINE_WORKSPACE_BACKEND=overlay", "QUINE_WORKSPACE_REVISION_MODE=restore"},
		WorkDir:          tmpDir,
		ForkWorldEnabled: true,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", World: config.WorldHost, Protection: config.ProtectionNone, Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected host child to run, got error:\n%s", result.Content)
	}
	if strings.Contains(string(result.Content), "QUINE_WORKSPACE_ROOT=") {
		t.Fatalf("host child should not inherit workspace env:\n%s", result.Content)
	}
	if strings.Contains(string(result.Content), "QUINE_WORKSPACE=") {
		t.Fatalf("host child should not inherit bare workspace env:\n%s", result.Content)
	}
	if strings.Contains(string(result.Content), "QUINE_FORK_WORLD_ENABLED=") {
		t.Fatalf("host child should not inherit fork-world env without workspace physics:\n%s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	children := payload["children"].([]any)
	child := children[0].(map[string]any)
	if child["world"] != "host" {
		t.Fatalf("child world = %#v, want host", child["world"])
	}
}

func TestForkExecutor_ForkWorldEnabled_HostChildImportsHostHandoffIntoSubjectiveParent(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	tmpDir := t.TempDir()
	workspaceRoot := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(workspaceRoot, 0o755); err != nil {
		t.Fatalf("mkdir workspace root: %v", err)
	}
	helper := filepath.Join(tmpDir, "host-handoff.sh")
	script := `#!/bin/sh
set -eu
mkdir -p handoff
printf 'from-host\n' > handoff/child.txt
printf 'host handoff complete\n'
`
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:                  helper,
		DataDir:                    filepath.Join(tmpDir, "data"),
		SessionID:                  "test-session",
		ContextRoot:                filepath.Join(tmpDir, "context"),
		DefaultTimeout:             5 * time.Second,
		MaxOutput:                  10000,
		WorkspaceEnabled:           true,
		WorkspaceRoot:              workspaceRoot,
		Workspace:                  workspaceRoot,
		WorkspaceBackend:           "overlay",
		WorkspaceRevisionMode:      config.WorkspaceRevisionRestore,
		WorkDir:                    workspaceRoot,
		ForkWorldEnabled:           true,
		FSMutationTelemetryEnabled: true,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "write a host handoff file", World: config.WorldHost, Protection: config.ProtectionNone, Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected host handoff fork to succeed, got:\n%s", result.Content)
	}

	payload := decodeForkResultContent(t, result.Content)
	fsMutations, ok := payload["fs_mutations"].(string)
	if !ok || !strings.Contains(fsMutations, "+ handoff/child.txt (created)") {
		t.Fatalf("expected imported host handoff mutation, got:\n%s", result.Content)
	}
	worldRevision, ok := payload["world_revision"].(string)
	if !ok || !strings.Contains(worldRevision, "created=wr1 parent=wr0 current=wr1") {
		t.Fatalf("expected host handoff import to create a parent revision, got:\n%s", result.Content)
	}
	if got := readOverlayWorkspaceFile(t, f.subjective, "handoff/child.txt"); got != "from-host" {
		t.Fatalf("parent subjective workspace file = %q, want from-host", got)
	}
	hostData, err := os.ReadFile(filepath.Join(workspaceRoot, "handoff", "child.txt"))
	if err != nil {
		t.Fatalf("host handoff file should also exist on host surface: %v", err)
	}
	if strings.TrimSpace(string(hostData)) != "from-host" {
		t.Fatalf("host handoff file = %q, want from-host", strings.TrimSpace(string(hostData)))
	}

	if err := os.MkdirAll(filepath.Join(f.subjective.liveUpperDir(), "handoff"), 0o755); err != nil {
		t.Fatalf("mkdir parent live upper handoff: %v", err)
	}
	if err := os.WriteFile(filepath.Join(f.subjective.liveUpperDir(), "handoff", "child.txt"), []byte("from-parent\n"), 0o644); err != nil {
		t.Fatalf("write parent live upper handoff: %v", err)
	}
	if _, err := f.subjective.finalizeTurn("test-parent", 2); err != nil {
		t.Fatalf("finalize parent overlay change: %v", err)
	}

	noopHelper := filepath.Join(tmpDir, "host-noop.sh")
	if err := os.WriteFile(noopHelper, []byte("#!/bin/sh\nprintf 'noop host child\\n'\n"), 0o755); err != nil {
		t.Fatalf("write noop helper: %v", err)
	}
	f.QuinePath = noopHelper
	result = f.Execute("tool-2", ForkRequest{
		Children: []ForkChild{{Intent: "no-op host child", World: config.WorldHost, Protection: config.ProtectionNone, Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected no-op host fork to succeed, got:\n%s", result.Content)
	}
	if got := readOverlayWorkspaceFile(t, f.subjective, "handoff/child.txt"); got != "from-parent" {
		t.Fatalf("no-op host import should not replay stale host state, got %q", got)
	}
}

func TestForkExecutor_ForkWorldEnabled_FailedHostChildDoesNotImportOnLaterSuccess(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	tmpDir := t.TempDir()
	workspaceRoot := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(workspaceRoot, 0o755); err != nil {
		t.Fatalf("mkdir workspace root: %v", err)
	}
	failHelper := filepath.Join(tmpDir, "host-fails-after-write.sh")
	failScript := `#!/bin/sh
set -eu
printf 'failed host residue\n' > failed.txt
exit 2
`
	if err := os.WriteFile(failHelper, []byte(failScript), 0o755); err != nil {
		t.Fatalf("write failing helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:             failHelper,
		DataDir:               filepath.Join(tmpDir, "data"),
		SessionID:             "test-session",
		ContextRoot:           filepath.Join(tmpDir, "context"),
		DefaultTimeout:        5 * time.Second,
		MaxOutput:             10000,
		WorkspaceEnabled:      true,
		WorkspaceRoot:         workspaceRoot,
		Workspace:             workspaceRoot,
		WorkspaceBackend:      "overlay",
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
		WorkDir:               workspaceRoot,
		ForkWorldEnabled:      true,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "fail after host write", World: config.WorldHost, Protection: config.ProtectionNone, Scope: "."}},
		Mode:     ForkModeWait,
	})
	if !result.IsError {
		t.Fatalf("expected failing host fork to be an error, got:\n%s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	fsMutations, _ := payload["fs_mutations"].(string)
	if strings.Contains(fsMutations, "failed.txt") {
		t.Fatalf("failed host child residue should not be imported, got:\n%s", result.Content)
	}

	noopHelper := filepath.Join(tmpDir, "host-noop.sh")
	if err := os.WriteFile(noopHelper, []byte("#!/bin/sh\nprintf 'noop host child\\n'\n"), 0o755); err != nil {
		t.Fatalf("write noop helper: %v", err)
	}
	f.QuinePath = noopHelper
	result = f.Execute("tool-2", ForkRequest{
		Children: []ForkChild{{Intent: "no-op host child", World: config.WorldHost, Protection: config.ProtectionNone, Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected later no-op host fork to succeed, got:\n%s", result.Content)
	}
	payload = decodeForkResultContent(t, result.Content)
	fsMutations, _ = payload["fs_mutations"].(string)
	if strings.Contains(fsMutations, "failed.txt") {
		t.Fatalf("later successful host child should not import older failed residue, got:\n%s", result.Content)
	}
	if _, err := f.subjective.readOverlayWorkspaceFile("failed.txt"); err == nil {
		t.Fatalf("failed host residue should remain outside parent subjective workspace")
	}
}

func TestForkExecutor_ForkWorldEnabled_SubjectiveChildInjectsWorkspacePhysicsOrErrors(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "print-env.sh")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\nenv | sort\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:             helper,
		DataDir:               tmpDir,
		SessionID:             "test-session",
		ContextRoot:           filepath.Join(tmpDir, "context"),
		DefaultTimeout:        5 * time.Second,
		MaxOutput:             10000,
		Env:                   []string{},
		WorkDir:               tmpDir,
		WorkspaceEnabled:      true,
		WorkspaceRoot:         filepath.Join(tmpDir, "workspace"),
		Workspace:             filepath.Join(tmpDir, "workspace"),
		WorkspaceBackend:      "direct",
		WorkspaceRevisionMode: config.WorkspaceRevisionNone,
		ForkWorldEnabled:      true,
	}
	if err := os.MkdirAll(f.Workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", World: config.WorldSubjective, Protection: config.ProtectionTransactional, Scope: "."}},
		Mode:     ForkModeWait,
	})

	if runtime.GOOS != "linux" {
		if !result.IsError {
			t.Fatalf("expected non-Linux subjective child request to fail, got:\n%s", result.Content)
		}
		payload := decodeForkResultContent(t, result.Content)
		if payload["status"] != "error" {
			t.Fatalf("status = %#v, want error", payload["status"])
		}
		errors, ok := payload["errors"].([]any)
		if !ok || len(errors) != 1 {
			t.Fatalf("errors = %#v, want single entry", payload["errors"])
		}
		if !strings.Contains(errors[0].(string), `world="subjective"`) {
			t.Fatalf("unexpected non-Linux error:\n%s", result.Content)
		}
		return
	}

	if result.IsError {
		t.Fatalf("expected subjective child to run on Linux, got error:\n%s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	children := payload["children"].([]any)
	child := children[0].(map[string]any)
	stdout, _ := child["stdout"].(string)
	checks := []string{
		"QUINE_WORKSPACE_ROOT=" + f.WorkspaceRoot,
		"QUINE_WORKSPACE=" + f.Workspace,
		"QUINE_WORKSPACE_BACKEND=direct",
		"QUINE_WORKSPACE_REVISION_MODE=none",
	}
	for _, want := range checks {
		if !strings.Contains(stdout, want) {
			t.Fatalf("subjective child output missing %q:\n%s", want, result.Content)
		}
	}
	if child["world"] != "subjective" {
		t.Fatalf("child world = %#v, want subjective", child["world"])
	}
	if child["protection"] != "transactional" {
		t.Fatalf("child protection = %#v, want transactional", child["protection"])
	}
}

func TestForkExecutor_Execute_GatherAll_ReportsFSMutations(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	tmpDir := t.TempDir()
	workspaceRoot := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(workspaceRoot, 0o755); err != nil {
		t.Fatalf("mkdir workspace root: %v", err)
	}
	sessionID := "test-session"
	helper := filepath.Join(tmpDir, "mutate-view.sh")
	script := `#!/bin/sh
set -eu
printf 'fork wrote this\n' > child.txt
printf 'mutated\n'
`
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:                  helper,
		DataDir:                    filepath.Join(tmpDir, "data"),
		SessionID:                  sessionID,
		ContextRoot:                filepath.Join(tmpDir, "context"),
		DefaultTimeout:             5 * time.Second,
		MaxOutput:                  10000,
		WorkspaceEnabled:           true,
		WorkspaceRoot:              workspaceRoot,
		Workspace:                  workspaceRoot,
		WorkspaceBackend:           "overlay",
		WorkspaceRevisionMode:      config.WorkspaceRevisionRestore,
		WorkDir:                    tmpDir,
		FSMutationTelemetryEnabled: true,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected mutation fork to succeed, got:\n%s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	fsMutations, ok := payload["fs_mutations"].(string)
	if !ok || fsMutations == "" {
		t.Fatalf("expected fork result to include FS mutations, got:\n%s", result.Content)
	}
	if !strings.Contains(fsMutations, "[FS MUTATIONS]\n(empty)") {
		t.Fatalf("expected parent fork result to report no parent-side mutations, got:\n%s", result.Content)
	}
	worldRevision, ok := payload["world_revision"].(string)
	if !ok || !strings.Contains(worldRevision, "[WORLD REVISION] current=wr0 (unchanged)") {
		t.Fatalf("expected fork result to leave parent world revision unchanged, got:\n%s", result.Content)
	}
	if !strings.Contains(fsMutations, "[FS MUTATIONS]\n(empty)") {
		t.Fatalf("structured fs_mutations should report no parent-side mutations: %q", fsMutations)
	}
	if !strings.Contains(worldRevision, "current=wr0 (unchanged)") {
		t.Fatalf("structured world_revision should keep wr0 unchanged: %q", worldRevision)
	}
}

func TestForkExecutor_Execute_GatherAll_CanHideFSMutationTelemetry(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	tmpDir := t.TempDir()
	workspaceRoot := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(workspaceRoot, 0o755); err != nil {
		t.Fatalf("mkdir workspace root: %v", err)
	}
	sessionID := "test-session-hidden-fs"
	helper := filepath.Join(tmpDir, "mutate-view.sh")
	script := `#!/bin/sh
set -eu
printf 'fork wrote this\n' > child.txt
printf 'mutated\n'
`
	if err := os.WriteFile(helper, []byte(script), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	f := &ForkExecutor{
		QuinePath:                  helper,
		DataDir:                    filepath.Join(tmpDir, "data"),
		SessionID:                  sessionID,
		ContextRoot:                filepath.Join(tmpDir, "context"),
		DefaultTimeout:             5 * time.Second,
		MaxOutput:                  10000,
		WorkspaceEnabled:           true,
		WorkspaceRoot:              workspaceRoot,
		Workspace:                  workspaceRoot,
		WorkspaceBackend:           "overlay",
		WorkspaceRevisionMode:      config.WorkspaceRevisionRestore,
		WorkDir:                    tmpDir,
		FSMutationTelemetryEnabled: false,
	}
	prepareForkContextRoot(t, f.ContextRoot)

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Scope: "."}},
		Mode:     ForkModeWait,
	})
	if result.IsError {
		t.Fatalf("expected mutation fork to succeed, got:\n%s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	if _, ok := payload["fs_mutations"]; ok {
		t.Fatalf("expected fork result to hide fs_mutations telemetry, got:\n%s", result.Content)
	}
	worldRevision, ok := payload["world_revision"].(string)
	if !ok || !strings.Contains(worldRevision, "[WORLD REVISION] current=wr0 (unchanged)") {
		t.Fatalf("expected fork result to keep parent world revision visible, got:\n%s", result.Content)
	}
}

func TestChildTapeDirPrefersRetentionSessionsTree(t *testing.T) {
	sessionID := "child-session"

	if got := childTapeDir("/tmp/runtime", nil, sessionID); got != "/tmp/runtime/log/child-session/tapes" {
		t.Fatalf("childTapeDir() without retention dir = %q, want %q", got, "/tmp/runtime/log/child-session/tapes")
	}

	env := []string{
		"QUINE_MODEL_ID=test-model",
		"QUINE_RETENTION_DIR=/tmp/retained",
	}
	if got := childTapeDir("/tmp/runtime", env, sessionID); got != "/tmp/retained/sessions/child-session/tapes" {
		t.Fatalf("childTapeDir() with retention dir = %q, want %q", got, "/tmp/retained/sessions/child-session/tapes")
	}
}
