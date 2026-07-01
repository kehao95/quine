package world

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateItems(t *testing.T) {
	t.Parallel()
	items, err := GenerateItems(20)
	if err != nil {
		t.Fatalf("GenerateItems: %v", err)
	}
	if len(items) != 20 {
		t.Fatalf("got %d items, want 20", len(items))
	}
	for i := 1; i <= 20; i++ {
		id := "c" + padInt(i)
		if _, ok := items[id]; !ok {
			t.Fatalf("missing cell %s", id)
		}
	}
}

func padInt(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func TestStateLoadSaveRoundtrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	// Fresh load.
	st, err := LoadState(dir, 22)
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}
	if st.BudgetRemaining != 22 {
		t.Fatalf("budget = %d, want 22", st.BudgetRemaining)
	}
	if len(st.Collected) != 0 {
		t.Fatalf("collected = %d, want 0", len(st.Collected))
	}
	if len(st.AgentGets) != 0 {
		t.Fatalf("agent gets = %d, want 0", len(st.AgentGets))
	}
	if len(st.ResetVotes) != 0 {
		t.Fatalf("reset votes = %d, want 0", len(st.ResetVotes))
	}

	// Modify and save.
	st.BudgetRemaining = 18
	st.Collected["c01"] = CollectedCell{Value: "alpha", PID: 12345}
	st.AgentGets["agent-1"] = 2
	st.ResetVotes["agent-2"] = true
	if err := st.Save(dir); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Reload.
	st2, err := LoadState(dir, 22)
	if err != nil {
		t.Fatalf("LoadState after save: %v", err)
	}
	if st2.BudgetRemaining != 18 {
		t.Fatalf("budget = %d, want 18", st2.BudgetRemaining)
	}
	if st2.Collected["c01"].Value != "alpha" {
		t.Fatalf("collected c01 = %q, want %q", st2.Collected["c01"].Value, "alpha")
	}
	if st2.AgentGets["agent-1"] != 2 {
		t.Fatalf("agent-1 gets = %d, want 2", st2.AgentGets["agent-1"])
	}
	if !st2.ResetVotes["agent-2"] {
		t.Fatal("expected persisted reset vote for agent-2")
	}
}

func TestLockUnlock(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	f, err := LockState(dir)
	if err != nil {
		t.Fatalf("LockState: %v", err)
	}
	// Lock file should exist.
	if _, err := os.Stat(filepath.Join(dir, "lock")); err != nil {
		t.Fatalf("lock file missing: %v", err)
	}
	UnlockState(f)
}

func TestResolveAgentIDUsesRunID(t *testing.T) {
	t.Parallel()

	getenv := func(key string) string {
		if key == "QUINE_RUN_ID" {
			return "run-env"
		}
		return ""
	}

	if got := ResolveAgentID("", getenv); got != "run-env" {
		t.Fatalf("ResolveAgentID() = %q, want %q", got, "run-env")
	}
}

func TestResolveAgentIDPrefersExecutableAgentID(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	exePath := filepath.Join(dir, "world")
	if err := os.WriteFile(filepath.Join(dir, "agent-id.txt"), []byte("sealed-agent\n"), 0o644); err != nil {
		t.Fatalf("write agent-id.txt: %v", err)
	}
	getenv := func(key string) string {
		if key == "QUINE_RUN_ID" {
			return "run-env"
		}
		return ""
	}
	if got := ResolveAgentID(exePath, getenv); got != "sealed-agent" {
		t.Fatalf("ResolveAgentID() = %q, want sealed-agent", got)
	}
}

func TestResolveExecutablePathFallsBackToArgv0(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	worldPath := filepath.Join(dir, "world")
	if err := os.WriteFile(worldPath, []byte{}, 0o755); err != nil {
		t.Fatalf("write world: %v", err)
	}
	got := resolveExecutablePath(
		func() (string, error) { return "", os.ErrNotExist },
		func(name string) (string, error) {
			if name != "world" {
				t.Fatalf("lookPath name = %q, want world", name)
			}
			return worldPath, nil
		},
		[]string{"world"},
	)
	if got != worldPath {
		t.Fatalf("resolveExecutablePath() = %q, want %q", got, worldPath)
	}
}

func TestResolveAgentIDLegacyFallbacks(t *testing.T) {
	t.Parallel()

	getenv := func(key string) string {
		switch key {
		case "WORLD_AGENT_ID":
			return "world-agent"
		case "QUINE_SESSION_ID":
			return "legacy-session"
		default:
			return ""
		}
	}
	if got := ResolveAgentID("", getenv); got != "world-agent" {
		t.Fatalf("ResolveAgentID() = %q, want world-agent", got)
	}
	getenv = func(key string) string {
		if key == "QUINE_SESSION_ID" {
			return "legacy-session"
		}
		return ""
	}
	if got := ResolveAgentID("", getenv); got != "legacy-session" {
		t.Fatalf("ResolveAgentID() = %q, want legacy-session", got)
	}
}

func TestResolveAgentIDReturnsEmptyWithoutIdentity(t *testing.T) {
	t.Parallel()

	getenv := func(key string) string {
		return ""
	}
	if got := ResolveAgentID("", getenv); got != "" {
		t.Fatalf("ResolveAgentID() = %q, want empty", got)
	}
}

func TestEnforceSingleWorldInvocationPerShell(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	t.Setenv("QUINE_WORLD_ONE_PER_SHELL", "1")
	t.Setenv("QUINE_SESSION_ID", "sess-test")

	if err := EnforceSingleWorldInvocationPerShell(os.Getenv); err != nil {
		t.Fatalf("first guard = %v, want nil", err)
	}
	if err := EnforceSingleWorldInvocationPerShell(os.Getenv); err == nil {
		t.Fatal("second guard = nil, want error")
	}
}

func TestAppendEvent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	ev := Event{Time: "t1", Action: "get", Cell: "c01", PID: 1, BudgetAfter: 21, Result: "ok"}
	if err := AppendEvent(dir, ev); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	ev2 := Event{Time: "t2", Action: "get", Cell: "c02", PID: 2, BudgetAfter: 20, Result: "ok"}
	if err := AppendEvent(dir, ev2); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "events.jsonl"))
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}
	if lines != 2 {
		t.Fatalf("got %d lines, want 2", lines)
	}
}

func TestBudgetedSpec(t *testing.T) {
	t.Parallel()
	plain := &Spec{Items: map[string]string{"1": "a"}}
	if plain.Budgeted() {
		t.Fatal("plain spec should not be budgeted")
	}

	budgeted := &Spec{
		Items:  map[string]string{"c01": "x"},
		Config: &SpecConfig{Budget: 22, Cells: 20, StateDir: "/tmp/test", AgentGetLimit: 10},
	}
	if !budgeted.Budgeted() {
		t.Fatal("budgeted spec should be budgeted")
	}
}

func TestSaveSpec(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "world.json")

	spec := &Spec{
		Items:  map[string]string{"c01": "alpha"},
		Config: &SpecConfig{Budget: 22, Cells: 20, StateDir: ".state", AgentGetLimit: 10},
	}
	if err := SaveSpec(spec, path); err != nil {
		t.Fatalf("SaveSpec: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Items["c01"] != "alpha" {
		t.Fatalf("item = %q, want %q", loaded.Items["c01"], "alpha")
	}
	if !loaded.Budgeted() {
		t.Fatal("loaded spec should be budgeted")
	}
	if loaded.Config.Budget != 22 {
		t.Fatalf("budget = %d, want 22", loaded.Config.Budget)
	}
	if loaded.Config.AgentGetLimit != 10 {
		t.Fatalf("agent get limit = %d, want 10", loaded.Config.AgentGetLimit)
	}
}
