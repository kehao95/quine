package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"--help"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run(--help) exitCode = %d, want 0", exitCode)
	}
	if stderr.Len() != 0 {
		t.Fatalf("run(--help) stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "world get <id>") {
		t.Fatalf("run(--help) stdout = %q, want usage text", stdout.String())
	}
}

func TestRunGet(t *testing.T) {
	dir := t.TempDir()
	// Set up QUINE_DATA_DIR/world/world.json structure
	worldDir := filepath.Join(dir, "world")
	if err := os.MkdirAll(worldDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(worldDir, "world.json")
	if err := os.WriteFile(path, []byte(`{"items":{"7":"payload-7"}}`), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	t.Setenv("QUINE_DATA_DIR", dir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"get", "7"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run(get 7) exitCode = %d, want 0", exitCode)
	}
	if stderr.Len() != 0 {
		t.Fatalf("run(get 7) stderr = %q, want empty", stderr.String())
	}
	if stdout.String() != "payload-7\n" {
		t.Fatalf("run(get 7) stdout = %q, want %q", stdout.String(), "payload-7\n")
	}
}

func TestRunGetUnknownID(t *testing.T) {
	dir := t.TempDir()
	worldDir := filepath.Join(dir, "world")
	if err := os.MkdirAll(worldDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(worldDir, "world.json")
	if err := os.WriteFile(path, []byte(`{"items":{"7":"payload-7"}}`), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	t.Setenv("QUINE_DATA_DIR", dir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{"get", "9"}, &stdout, &stderr)
	if exitCode != 2 {
		t.Fatalf("run(get 9) exitCode = %d, want 2", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("run(get 9) stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown id: 9") {
		t.Fatalf("run(get 9) stderr = %q, want unknown id", stderr.String())
	}
}

func writeBudgetedSpec(t *testing.T, dir string, items map[string]string, budget, cells int) string {
	t.Helper()
	// Create QUINE_DATA_DIR/world/world.json structure
	worldDir := filepath.Join(dir, "world")
	if err := os.MkdirAll(worldDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	specPath := filepath.Join(worldDir, "world.json")
	stateDir := filepath.Join(dir, "state")

	spec := struct {
		Items  map[string]string `json:"items"`
		Config struct {
			Budget        int    `json:"budget"`
			Cells         int    `json:"cells"`
			StateDir      string `json:"state_dir"`
			AgentGetLimit int    `json:"agent_get_limit,omitempty"`
		} `json:"config"`
	}{Items: items}
	spec.Config.Budget = budget
	spec.Config.Cells = cells
	spec.Config.StateDir = stateDir

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(specPath, data, 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	return dir // Return QUINE_DATA_DIR, not spec path
}

func writeBudgetedSpecWithAgentLimit(t *testing.T, dir string, items map[string]string, budget, cells, agentGetLimit int) string {
	return writeBudgetedSpecWithLimits(t, dir, items, budget, cells, agentGetLimit, 0)
}

func writeBudgetedSpecWithLimits(t *testing.T, dir string, items map[string]string, budget, cells, agentGetLimit, resetQuorum int) string {
	t.Helper()
	worldDir := filepath.Join(dir, "world")
	if err := os.MkdirAll(worldDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	specPath := filepath.Join(worldDir, "world.json")
	stateDir := filepath.Join(dir, "state")

	spec := struct {
		Items  map[string]string `json:"items"`
		Config struct {
			Budget        int    `json:"budget"`
			Cells         int    `json:"cells"`
			StateDir      string `json:"state_dir"`
			AgentGetLimit int    `json:"agent_get_limit,omitempty"`
			ResetQuorum   int    `json:"reset_quorum,omitempty"`
		} `json:"config"`
	}{Items: items}
	spec.Config.Budget = budget
	spec.Config.Cells = cells
	spec.Config.StateDir = stateDir
	spec.Config.AgentGetLimit = agentGetLimit
	spec.Config.ResetQuorum = resetQuorum

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(specPath, data, 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	return dir // Return QUINE_DATA_DIR, not spec path
}

func TestBudgetedGet(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha", "c02": "beta"}
	specPath := writeBudgetedSpec(t, dir, items, 3, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)

	// First get: should succeed, budget 2/3.
	var stdout, stderr bytes.Buffer
	code := run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("get c01: exit %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "alpha") {
		t.Fatalf("get c01: stdout=%q, want alpha", stdout.String())
	}
	if !strings.Contains(stdout.String(), "2/3") {
		t.Fatalf("get c01: stdout=%q, want budget 2/3", stdout.String())
	}
	if !strings.Contains(stdout.String(), "generation: 1") {
		t.Fatalf("get c01: stdout=%q, want generation 1", stdout.String())
	}

	// Second get: budget 1/3.
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c02"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("get c02: exit %d", code)
	}
	if !strings.Contains(stdout.String(), "1/3") {
		t.Fatalf("get c02: stdout=%q, want budget 1/3", stdout.String())
	}

	// Third get (duplicate c01): budget 0/3.
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("get c01 again: exit %d", code)
	}
	if !strings.Contains(stdout.String(), "0/3") {
		t.Fatalf("get c01 again: stdout=%q, want budget 0/3", stdout.String())
	}

	// Fourth get: budget exhausted.
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c02"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("get after exhaustion: exit %d, want 3", code)
	}
	if !strings.Contains(stderr.String(), "budget exhausted") {
		t.Fatalf("get after exhaustion: stderr=%q", stderr.String())
	}
}

func TestBudgetedReset(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha", "c02": "beta"}
	specPath := writeBudgetedSpec(t, dir, items, 2, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)

	// Exhaust budget.
	var stdout, stderr bytes.Buffer
	run([]string{"get", "c01"}, &stdout, &stderr)
	stdout.Reset()
	stderr.Reset()
	run([]string{"get", "c02"}, &stdout, &stderr)
	stdout.Reset()
	stderr.Reset()

	// Confirm exhausted.
	code := run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("should be exhausted, got exit %d", code)
	}

	// Reset.
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"reset"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("reset: exit %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "reset complete") {
		t.Fatalf("reset: stdout=%q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "2/2") {
		t.Fatalf("reset: stdout=%q, want budget 2/2", stdout.String())
	}
	if !strings.Contains(stdout.String(), "generation: 2") {
		t.Fatalf("reset: stdout=%q, want generation 2", stdout.String())
	}

	// Get after reset should work, but values are regenerated.
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("get after reset: exit %d, stderr=%q", code, stderr.String())
	}
	// Value should be different from original "alpha" (random hex).
	// We can't guarantee it's different (astronomically unlikely to collide though).
	if !strings.Contains(stdout.String(), "1/2") {
		t.Fatalf("get after reset: stdout=%q, want budget 1/2", stdout.String())
	}
}

func TestBudgetedValidateAcceptsCorrectFile(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha", "c02": "beta"}
	specPath := writeBudgetedSpec(t, dir, items, 2, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)

	resultsPath := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(resultsPath, []byte("c01: alpha\nc02: beta\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"validate", resultsPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("validate accepted: exit %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "validate accepted") {
		t.Fatalf("validate accepted stdout=%q", stdout.String())
	}
}

func TestBudgetedValidateCanRecoverAfterRejectedAttempt(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha", "c02": "beta"}
	specPath := writeBudgetedSpec(t, dir, items, 2, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)

	resultsPath := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(resultsPath, []byte("c01: alpha\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"validate", resultsPath}, &stdout, &stderr)
	if code != 4 {
		t.Fatalf("validate rejected: exit %d, want 4", code)
	}
	if !strings.Contains(stderr.String(), "validate rejected") {
		t.Fatalf("validate rejected stderr=%q", stderr.String())
	}

	if err := os.WriteFile(resultsPath, []byte("c01: wrong\nc02: wrong\n"), 0o644); err != nil {
		t.Fatalf("rewrite results: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"validate", resultsPath}, &stdout, &stderr)
	if code != 4 {
		t.Fatalf("validate wrong results: exit %d, want 4", code)
	}
	if !strings.Contains(stderr.String(), "validate rejected") {
		t.Fatalf("validate wrong stderr=%q", stderr.String())
	}

	if err := os.WriteFile(resultsPath, []byte("c01: alpha\nc02: beta\n"), 0o644); err != nil {
		t.Fatalf("rewrite corrected results: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"validate", resultsPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("validate corrected results: exit %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "validate accepted") {
		t.Fatalf("validate corrected stdout=%q", stdout.String())
	}
}

func TestBudgetedHelp(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha"}
	specPath := writeBudgetedSpecWithLimits(t, dir, items, 22, 20, 10, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)

	var stdout, stderr bytes.Buffer
	code := run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("help: exit %d", code)
	}
	out := stdout.String()
	if !strings.Contains(out, "20 cells") {
		t.Fatalf("help should mention 20 cells: %q", out)
	}
	if !strings.Contains(out, "22 get calls") {
		t.Fatalf("help should mention 22 budget: %q", out)
	}
	if !strings.Contains(out, "world reset") {
		t.Fatalf("help should mention reset: %q", out)
	}
	if !strings.Contains(out, "world validate <path>") {
		t.Fatalf("help should mention validate: %q", out)
	}
	if !strings.Contains(out, "10 calls per reset epoch") {
		t.Fatalf("help should mention per-process get limit: %q", out)
	}
	if !strings.Contains(out, "`world reset` executes only after 2 requests in the same generation") {
		t.Fatalf("help should mention reset quorum: %q", out)
	}
	if strings.Contains(out, "WORLD_AGENT_ID") || strings.Contains(out, "QUINE_SESSION_ID") {
		t.Fatalf("help should not expose identity env wiring: %q", out)
	}
}

func TestBudgetedPerAgentCallLimitBlocksSoloGetSequence(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha", "c02": "beta"}
	specPath := writeBudgetedSpecWithAgentLimit(t, dir, items, 20, 2, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)
	t.Setenv("QUINE_SESSION_ID", "agent-1")

	var stdout, stderr bytes.Buffer
	code := run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("first get exit = %d, stderr=%q", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c02"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("second get exit = %d, stderr=%q", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("third get exit = %d, want 5", code)
	}
	if !strings.Contains(stderr.String(), "get limit exhausted") {
		t.Fatalf("third get stderr=%q", stderr.String())
	}
}

func TestBudgetedGetRequiresStableAgentIdentityWhenGetCapEnabled(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha"}
	specPath := writeBudgetedSpecWithAgentLimit(t, dir, items, 20, 1, 1)
	t.Setenv("QUINE_DATA_DIR", specPath)
	t.Setenv("QUINE_SESSION_ID", "")

	var stdout, stderr bytes.Buffer
	code := run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("get exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "agent identity unavailable") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestBudgetedPerAgentCallLimitResetsPerEpoch(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha"}
	specPath := writeBudgetedSpecWithAgentLimit(t, dir, items, 20, 1, 1)
	t.Setenv("QUINE_DATA_DIR", specPath)
	t.Setenv("QUINE_SESSION_ID", "agent-1")

	var stdout, stderr bytes.Buffer
	code := run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("first get exit = %d, stderr=%q", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("second get exit = %d, want 5", code)
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"reset"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("reset exit = %d, stderr=%q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("get after reset exit = %d, stderr=%q", code, stderr.String())
	}
}

func TestBudgetedValidateRemainsUnlimitedWhenGetCapIsReached(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha"}
	specPath := writeBudgetedSpecWithAgentLimit(t, dir, items, 20, 1, 1)
	t.Setenv("QUINE_DATA_DIR", specPath)
	t.Setenv("QUINE_SESSION_ID", "agent-1")

	resultsPath := filepath.Join(dir, "results.txt")
	if err := os.WriteFile(resultsPath, []byte("c01: alpha\n"), 0o644); err != nil {
		t.Fatalf("write results: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"get", "c01"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("get exit = %d, stderr=%q", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"validate", resultsPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("validate exit = %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "validate accepted") {
		t.Fatalf("validate stdout=%q", stdout.String())
	}
}

func TestBudgetedResetQuorumWaitsForAllParticipants(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha"}
	specPath := writeBudgetedSpecWithLimits(t, dir, items, 20, 1, 0, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)

	var stdout, stderr bytes.Buffer
	t.Setenv("QUINE_SESSION_ID", "agent-1")
	code := run([]string{"reset"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("first reset exit = %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "reset pending") || !strings.Contains(stdout.String(), "1/2") {
		t.Fatalf("first reset stdout=%q, want pending 1/2", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"reset"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("repeat same-agent reset exit = %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "1/2") {
		t.Fatalf("repeat same-agent reset stdout=%q, want still 1/2", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	t.Setenv("QUINE_SESSION_ID", "agent-2")
	code = run([]string{"reset"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("second participant reset exit = %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "reset complete") {
		t.Fatalf("second participant reset stdout=%q, want reset complete", stdout.String())
	}
	if !strings.Contains(stdout.String(), "generation: 2") {
		t.Fatalf("second participant reset stdout=%q, want generation 2", stdout.String())
	}
}

func TestBudgetedResetRequiresStableAgentIdentityWhenQuorumEnabled(t *testing.T) {
	dir := t.TempDir()
	items := map[string]string{"c01": "alpha"}
	specPath := writeBudgetedSpecWithLimits(t, dir, items, 20, 1, 0, 2)
	t.Setenv("QUINE_DATA_DIR", specPath)
	t.Setenv("QUINE_SESSION_ID", "")

	var stdout, stderr bytes.Buffer
	code := run([]string{"reset"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("reset exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "agent identity unavailable") {
		t.Fatalf("stderr=%q", stderr.String())
	}
}

func TestResetNotAvailableForPlainSpec(t *testing.T) {
	dir := t.TempDir()
	worldDir := filepath.Join(dir, "world")
	if err := os.MkdirAll(worldDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(worldDir, "world.json")
	if err := os.WriteFile(path, []byte(`{"items":{"1":"a"}}`), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Setenv("QUINE_DATA_DIR", dir)

	var stdout, stderr bytes.Buffer
	code := run([]string{"reset"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("reset on plain spec: exit %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "budgeted mode") {
		t.Fatalf("reset stderr=%q", stderr.String())
	}
}
