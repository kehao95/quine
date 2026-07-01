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

func writeSpawnStubQuine(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "stub-quine.sh")
	body := `#!/bin/sh
printf 'MISSION=%s\n' "$1"
printf 'SESSION=%s\n' "${QUINE_SESSION_ID:-}"
printf 'PARENT=%s\n' "${QUINE_PARENT_SESSION:-unset}"
printf 'CTX_BOOTSTRAP=%s\n' "${QUINE_CONTEXT_BOOTSTRAP:-unset}"
printf 'CTX_TAPE=%s\n' "${QUINE_CONTEXT_TAPE:-unset}"
printf 'WORKSPACE_ROOT=%s\n' "${QUINE_WORKSPACE_ROOT:-unset}"
printf 'WORKSPACE=%s\n' "${QUINE_WORKSPACE:-unset}"
printf 'WORKSPACE_BACKEND=%s\n' "${QUINE_WORKSPACE_BACKEND:-unset}"
printf 'WORKSPACE_REVISION_MODE=%s\n' "${QUINE_WORKSPACE_REVISION_MODE:-unset}"
printf 'WORKSPACE_OWNER=%s\n' "${QUINE_WORKSPACE_OWNER:-unset}"
printf 'WORKSPACE_BOOTSTRAP=%s\n' "${QUINE_WORKSPACE_BOOTSTRAP:-unset}"
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write stub quine: %v", err)
	}
	return path
}

func writeSpawnRaceStubQuine(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "stub-quine-race.sh")
	body := `#!/bin/sh
case "$1" in
  slow*) sleep 5; printf 'SLOW_DONE\n' ;;
  fast*) printf 'FAST_DONE\n' ;;
  *) printf 'MISSION=%s\n' "$1" ;;
esac
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write race stub quine: %v", err)
	}
	return path
}

func TestParseSpawnArgs(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}}
	req, err := ParseSpawnArgs(map[string]any{
		"children": []interface{}{
			map[string]any{"mission": "fresh review", "world": "subjective", "protection": "transactional", "scope": "."},
			map[string]any{"mission": "fresh implementation", "world": "host", "protection": "none"},
		},
		"mode": "race",
	}, cfg)
	if err != nil {
		t.Fatalf("ParseSpawnArgs() error = %v", err)
	}
	if req.Mode != SpawnModeRace {
		t.Fatalf("Mode = %q, want race", req.Mode)
	}
	if len(req.Children) != 2 || req.Children[0].Mission != "fresh review" {
		t.Fatalf("Children = %#v", req.Children)
	}
	if req.Children[0].World != config.WorldSubjective || req.Children[0].Protection != config.ProtectionTransactional || req.Children[0].Scope != "." {
		t.Fatalf("child[0] workspace fields = %#v", req.Children[0])
	}
	if req.Children[1].World != config.WorldHost || req.Children[1].Protection != config.ProtectionNone || req.Children[1].Scope != "." {
		t.Fatalf("child[1] workspace fields = %#v", req.Children[1])
	}
}

func TestParseSpawnArgsRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		errWord string
	}{
		{"MissingChildren", map[string]any{"mode": "wait"}, "children"},
		{"EmptyChildren", map[string]any{"children": []interface{}{}}, "at least one"},
		{"ChildNotObject", map[string]any{"children": []interface{}{"bad"}}, "object"},
		{"MissingMission", map[string]any{"children": []interface{}{map[string]any{"intent": "old fork shape"}}}, "mission"},
		{"EmptyMission", map[string]any{"children": []interface{}{map[string]any{"mission": " "}}}, "mission"},
		{"InvalidMode", map[string]any{"children": []interface{}{map[string]any{"mission": "ok"}}, "mode": "later"}, "must be one of"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSpawnArgs(tt.args, nil)
			if err == nil || !strings.Contains(err.Error(), tt.errWord) {
				t.Fatalf("ParseSpawnArgs() error = %v, want containing %q", err, tt.errWord)
			}
		})
	}
}

func TestParseSpawnArgsRejectsInvalidWorldPairs(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}}
	tests := []struct {
		name    string
		child   map[string]any
		errWord string
	}{
		{"HostTransactional", map[string]any{"mission": "bad", "world": "host", "protection": "transactional"}, "host"},
		{"SubjectiveNone", map[string]any{"mission": "bad", "world": "subjective", "protection": "none"}, "subjective"},
		{"HostScoped", map[string]any{"mission": "bad", "world": "host", "protection": "none", "scope": "subdir"}, "scope"},
		{"AbsoluteScope", map[string]any{"mission": "bad", "world": "subjective", "protection": "transactional", "scope": "/tmp"}, "absolute"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSpawnArgs(map[string]any{"children": []interface{}{tt.child}}, cfg)
			if err == nil || !strings.Contains(err.Error(), tt.errWord) {
				t.Fatalf("ParseSpawnArgs() error = %v, want containing %q", err, tt.errWord)
			}
		})
	}
}

func TestSpawnExecutor_WaitStartsFreshProcessWithoutContextBootstrap(t *testing.T) {
	tmpDir := t.TempDir()
	stub := writeSpawnStubQuine(t, tmpDir)
	t.Setenv(ContextBootstrapEnv, "/parent/context/bootstrap")
	t.Setenv("QUINE_CONTEXT_TAPE", "/parent/context/tape.jsonl")

	cfg := &config.Config{Identity: config.Identity{SessionID: "parent-session"}, Limits: config.Limits{OutputTruncate: 20480}, Paths: config.Paths{DataDir: filepath.Join(tmpDir, "data"), RetentionDir: filepath.Join(tmpDir, "retained"), SelfReentryTarget: stub}}
	spawn := NewSpawnExecutor(cfg, []string{"QUINE_PARENT_SESSION=parent-session"})

	result := spawn.Execute("call_spawn", SpawnRequest{
		Mode: SpawnModeWait,
		Children: []SpawnChild{
			{Mission: "fresh child mission"},
		},
	})
	if result.IsError {
		t.Fatalf("spawn result IsError=true: %s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	if got := payload["tool"]; got != "spawn" {
		t.Fatalf("tool = %v, want spawn", got)
	}
	if got := payload["status"]; got != "completed" {
		t.Fatalf("status = %v, want completed", got)
	}
	children, ok := payload["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("children = %#v, want one child", payload["children"])
	}
	child, ok := children[0].(map[string]any)
	if !ok {
		t.Fatalf("child payload = %#v", children[0])
	}
	stdout, _ := child["stdout"].(string)
	for _, want := range []string{
		"MISSION=fresh child mission",
		"PARENT=parent-session",
		"CTX_BOOTSTRAP=unset",
		"CTX_TAPE=unset",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("child stdout missing %q:\n%s", want, stdout)
		}
	}
	if _, ok := child["seed_root"]; ok {
		t.Fatalf("spawn child result must not expose fork seed_root: %#v", child)
	}
	relationRoot, _ := payload["relation_root"].(string)
	if relationRoot == "" {
		t.Fatalf("relation_root missing in payload: %#v", payload)
	}
	relationDoc, err := os.ReadFile(filepath.Join(relationRoot, "relation.json"))
	if err != nil {
		t.Fatalf("read relation.json: %v", err)
	}
	if !strings.Contains(string(relationDoc), `"tool": "spawn"`) {
		t.Fatalf("relation.json should identify spawn tool:\n%s", relationDoc)
	}
}

func TestSpawnExecutor_AppliesWorkspaceProjectionLikeFork(t *testing.T) {
	tmpDir := t.TempDir()
	stub := writeSpawnStubQuine(t, tmpDir)
	cfg := &config.Config{Identity: config.Identity{SessionID: "parent-session"}, Limits: config.Limits{OutputTruncate: 20480}, ToolGates: config.ToolGates{ForkWorldEnabled: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: filepath.Join(tmpDir, "workspace"), Workspace: filepath.Join(tmpDir, "workspace"), WorkspaceBackend: "direct", WorkspaceRevisionMode: config.WorkspaceRevisionRestore}, Paths: config.Paths{DataDir: filepath.Join(tmpDir, "data"), RetentionDir: filepath.Join(tmpDir, "retained"), SelfReentryTarget: stub}}
	if err := os.MkdirAll(cfg.Workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	spawn := NewSpawnExecutor(cfg, []string{
		"QUINE_PARENT_SESSION=parent-session",
		"QUINE_WORKSPACE_BOOTSTRAP=parent-session",
		"QUINE_WORKSPACE_CURRENT_REVISION=wr7",
	})

	result := spawn.Execute("call_spawn_workspace", SpawnRequest{
		Mode: SpawnModeWait,
		Children: []SpawnChild{{
			Mission:    "fresh workspace child",
			World:      config.WorldSubjective,
			Protection: config.ProtectionTransactional,
			Scope:      ".",
		}},
	})
	if result.IsError {
		t.Fatalf("spawn result IsError=true: %s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	children := payload["children"].([]any)
	child := children[0].(map[string]any)
	stdout, _ := child["stdout"].(string)
	for _, want := range []string{
		"WORKSPACE_ROOT=" + cfg.WorkspaceRoot,
		"WORKSPACE=" + cfg.Workspace,
		"WORKSPACE_BACKEND=direct",
		"WORKSPACE_REVISION_MODE=restore",
		"WORKSPACE_OWNER=0",
		"WORKSPACE_BOOTSTRAP=parent-session",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("child stdout missing %q:\n%s", want, stdout)
		}
	}
	if got := child["world"]; got != "subjective" {
		t.Fatalf("child world = %#v, want subjective", got)
	}
	if got := child["protection"]; got != "transactional" {
		t.Fatalf("child protection = %#v, want transactional", got)
	}
	if _, ok := child["seed_root"]; ok {
		t.Fatalf("spawn child result must not expose fork seed_root: %#v", child)
	}
}

func TestSpawnExecutor_ForgetReturnsProcessHandles(t *testing.T) {
	tmpDir := t.TempDir()
	stub := writeSpawnStubQuine(t, tmpDir)
	cfg := &config.Config{Identity: config.Identity{SessionID: "parent-session"}, Limits: config.Limits{OutputTruncate: 20480}, Paths: config.Paths{DataDir: filepath.Join(tmpDir, "data"), RetentionDir: filepath.Join(tmpDir, "retained"), SelfReentryTarget: stub}}
	spawn := NewSpawnExecutor(cfg, []string{"QUINE_PARENT_SESSION=parent-session"})

	result := spawn.Execute("call_spawn_forget", SpawnRequest{
		Mode:     SpawnModeForget,
		Children: []SpawnChild{{Mission: "background fresh child"}},
	})
	if result.IsError {
		t.Fatalf("spawn forget IsError=true: %s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	if got := payload["status"]; got != "spawned" {
		t.Fatalf("status = %v, want spawned", got)
	}
	if got := payload["spawned"]; fmt.Sprint(got) != "1" {
		t.Fatalf("spawned = %v, want 1", got)
	}
	requireStableRelationCounters(t, payload, 1, 0, 0, 0)
	relationResult := requireRelationResult(t, payload)
	requireStableRelationCounters(t, relationResult, 1, 0, 0, 0)
	status := requireRelationStatus(t, payload)
	requireStableRelationCounters(t, status, 1, 0, 0, 0)
}

func TestSpawnExecutor_RaceReturnsFirstSuccessfulFreshProcess(t *testing.T) {
	tmpDir := t.TempDir()
	stub := writeSpawnRaceStubQuine(t, tmpDir)
	cfg := &config.Config{Identity: config.Identity{SessionID: "parent-session"}, Limits: config.Limits{OutputTruncate: 20480}, Paths: config.Paths{DataDir: filepath.Join(tmpDir, "data"), RetentionDir: filepath.Join(tmpDir, "retained"), SelfReentryTarget: stub}}
	spawn := NewSpawnExecutor(cfg, []string{"QUINE_PARENT_SESSION=parent-session"})

	result := spawn.Execute("call_spawn_race", SpawnRequest{
		Mode: SpawnModeRace,
		Children: []SpawnChild{
			{Mission: "slow child"},
			{Mission: "fast child"},
		},
	})
	if result.IsError {
		t.Fatalf("spawn race IsError=true: %s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	if got := payload["status"]; got != "completed" {
		t.Fatalf("status = %v, want completed", got)
	}
	winner, ok := payload["winner"].(map[string]any)
	if !ok {
		t.Fatalf("winner missing from payload: %#v", payload)
	}
	if got := fmt.Sprint(winner["index"]); got != "1" {
		t.Fatalf("winner index = %s, want 1", got)
	}
	if stdout, _ := winner["stdout"].(string); !strings.Contains(stdout, "FAST_DONE") {
		t.Fatalf("winner stdout missing FAST_DONE:\n%s", stdout)
	}
	if got := fmt.Sprint(payload["killed"]); got != "1" {
		t.Fatalf("killed = %v, want 1", payload["killed"])
	}
}

func TestSpawnExecutor_WaitTimeoutRelationStatusCountsTimeoutChildren(t *testing.T) {
	tmpDir := t.TempDir()
	helper := filepath.Join(tmpDir, "spawn-child.sh")
	script := `#!/bin/sh
case "$1" in
  fast)
    printf 'fast spawn done\n'
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

	spawn := &SpawnExecutor{ForkExecutor: &ForkExecutor{
		QuinePath:      helper,
		DataDir:        filepath.Join(tmpDir, "runtime"),
		SessionID:      "spawn-wait-timeout-parent",
		DefaultTimeout: 100 * time.Millisecond,
		MaxOutput:      10000,
		Env:            []string{},
		WorkDir:        tmpDir,
	}}

	result := spawn.Execute("tool-spawn-wait-timeout", SpawnRequest{
		Children: []SpawnChild{{Mission: "fast"}, {Mission: "slow"}},
		Mode:     SpawnModeWait,
	})
	if result.IsError {
		t.Fatalf("expected partial spawn success not to be an error: %s", result.Content)
	}
	payload := decodeForkResultContent(t, result.Content)
	if got := payload["status"]; got != "timeout" {
		t.Fatalf("status = %#v, want timeout", got)
	}
	requireStableRelationCounters(t, payload, 2, 2, 1, 1)
	relationResult := requireRelationResult(t, payload)
	requireStableRelationCounters(t, relationResult, 2, 2, 1, 1)
	status := requireRelationStatus(t, payload)
	requireStableRelationCounters(t, status, 2, 2, 1, 1)
}

func TestSpawnExecutor_WaitMissingBinaryPreservesPerChildErrors(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{Identity: config.Identity{SessionID: "spawn-missing-parent"}, Limits: config.Limits{OutputTruncate: 20480}, Paths: config.Paths{DataDir: filepath.Join(tmpDir, "data"), RetentionDir: filepath.Join(tmpDir, "retained"), SelfReentryTarget: "/nonexistent/quine"}}
	spawn := NewSpawnExecutor(cfg, nil)

	result := spawn.Execute("tool-spawn-missing", SpawnRequest{
		Mode:     SpawnModeWait,
		Children: []SpawnChild{{Mission: "task A"}, {Mission: "task B"}},
	})
	if !result.IsError {
		t.Fatalf("expected missing spawn binary to be an error")
	}
	payload := decodeForkResultContent(t, result.Content)
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
}

func TestSpawnExecutor_RaceWinnerPreservesSpawnErrors(t *testing.T) {
	tmpDir := t.TempDir()
	workspace := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
	stub := writeSpawnRaceStubQuine(t, tmpDir)
	cfg := &config.Config{Identity: config.Identity{SessionID: "spawn-race-spawn-errors-parent"}, Limits: config.Limits{OutputTruncate: 20480}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRoot: workspace, Workspace: workspace, WorkspaceBackend: "direct"}, Paths: config.Paths{DataDir: filepath.Join(tmpDir, "data"), RetentionDir: filepath.Join(tmpDir, "retained"), SelfReentryTarget: stub}}
	spawn := NewSpawnExecutor(cfg, nil)

	result := spawn.Execute("tool-spawn-race-errors", SpawnRequest{
		Mode: SpawnModeRace,
		Children: []SpawnChild{
			{Mission: "bad scope", Scope: "../outside"},
			{Mission: "fast child", Scope: "."},
		},
	})
	if result.IsError {
		t.Fatalf("spawn race should still return the successful winner: %s", result.Content)
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
