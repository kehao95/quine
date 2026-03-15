package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// envVars is the full list of environment variables we manage in tests.
var envVars = []string{
	// Core model/provider identity (required runtime wiring).
	"QUINE_MODEL_ID",
	"QUINE_API_TYPE",
	"QUINE_API_BASE",
	"QUINE_API_KEY",

	// Process lineage/session identity.
	"QUINE_MAX_DEPTH",
	"QUINE_DEPTH",
	"QUINE_SESSION_ID",
	"QUINE_TAPE_ID",
	"QUINE_PARENT_SESSION",

	// Runtime capacity limits.
	"QUINE_MAX_CONCURRENT",
	"QUINE_MAX_AGENTS",
	"QUINE_MAX_TURNS",

	// Execution/prompt behavior switches.
	"QUINE_TURN_EXHAUSTION_POLICY",
	"QUINE_PROMPT_METAPHOR",

	// Shell/tool runtime parameters.
	"QUINE_SH_TIMEOUT",
	"QUINE_OUTPUT_TRUNCATE",
	"QUINE_DATA_DIR",
	"QUINE_SHELL",
	"QUINE_WORKSPACE_ROOT",
	"QUINE_WORKSPACE",
	"QUINE_WORKSPACE_BACKEND",
	"QUINE_WORKSPACE_REVISION_MODE",
	"QUINE_WORKSPACE_CURRENT_REVISION",
	"QUINE_WORKSPACE_SOURCE",
	"QUINE_WORKSPACE_SESSION",
	"QUINE_WORKSPACE_OWNER",

	// Context and local config path.
	"QUINE_CONTEXT_WINDOW",
	"QUINE_CONFIG_DIR",
	"QUINE_ANCHOR_MEMORY",
}

// clearEnv unsets all managed env vars and returns a restore function.
func clearEnv(t *testing.T) {
	t.Helper()
	saved := make(map[string]string)
	for _, k := range envVars {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range envVars {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

// setRequired sets the 4 required env vars.
func setRequired(t *testing.T) {
	t.Helper()
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-20250514")
	os.Setenv("QUINE_API_TYPE", "anthropic")
	os.Setenv("QUINE_API_BASE", "https://api.anthropic.com")
	os.Setenv("QUINE_API_KEY", "sk-test-key")
}

func TestHappyPath(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	workspace := root + "/app"
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("realpath root: %v", err)
	}
	workspaceReal, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		t.Fatalf("realpath workspace: %v", err)
	}
	os.Setenv("QUINE_MAX_DEPTH", "10")
	os.Setenv("QUINE_DEPTH", "3")
	os.Setenv("QUINE_SESSION_ID", "my-session")
	os.Setenv("QUINE_TAPE_ID", "my-tape")
	os.Setenv("QUINE_PARENT_SESSION", "parent-session")
	os.Setenv("QUINE_MAX_CONCURRENT", "50")
	os.Setenv("QUINE_MAX_AGENTS", "25")
	os.Setenv("QUINE_SH_TIMEOUT", "60")
	os.Setenv("QUINE_OUTPUT_TRUNCATE", "4096")
	os.Setenv("QUINE_DATA_DIR", "/tmp/data")
	os.Setenv("QUINE_SHELL", "/bin/sh")
	os.Setenv("QUINE_MAX_TURNS", "30")
	os.Setenv("QUINE_TURN_EXHAUSTION_POLICY", "near_death_exec")
	os.Setenv("QUINE_PROMPT_METAPHOR", "thermodynamic")
	os.Setenv("QUINE_ANCHOR_MEMORY", "1")
	if runtime.GOOS == "linux" {
		os.Setenv("QUINE_WORKSPACE_ROOT", root)
		os.Setenv("QUINE_WORKSPACE", workspace)
	}

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"ModelID", c.ModelID, "claude-sonnet-4-20250514"},
		{"APIKey", c.APIKey, "sk-test-key"},
		{"APIBase", c.APIBase, "https://api.anthropic.com"},
		{"Provider", c.Provider, "anthropic"},
		{"MaxDepth", c.MaxDepth, 10},
		{"Depth", c.Depth, 3},
		{"SessionID", c.SessionID, "my-session"},
		{"TapeID", c.TapeID, "my-tape"},
		{"ParentSession", c.ParentSession, "parent-session"},
		{"MaxConcurrent", c.MaxConcurrent, 50},
		{"MaxAgents", c.MaxAgents, 25},
		{"ShTimeout", c.ShTimeout, 60},
		{"OutputTruncate", c.OutputTruncate, 4096},
		{"DataDir", c.DataDir, "/tmp/data"},
		{"Shell", c.Shell, "/bin/sh"},
		{"MaxTurns", c.MaxTurns, 30},
		{"TurnExhaustionPolicy", c.TurnExhaustionPolicy, TurnExhaustionNearDeathExec},
		{"PromptMetaphor", c.PromptMetaphor, PromptMetaphorThermodynamic},
		{"AnchorMemoryEnabled", c.AnchorMemoryEnabled, true},
	}
	if runtime.GOOS == "linux" {
		checks = append(checks,
			struct {
				name string
				got  any
				want any
			}{"WorkspaceEnabled", c.WorkspaceEnabled, true},
			struct {
				name string
				got  any
				want any
			}{"WorkspaceRoot", c.WorkspaceRoot, rootReal},
			struct {
				name string
				got  any
				want any
			}{"Workspace", c.Workspace, workspaceReal},
			struct {
				name string
				got  any
				want any
			}{"WorkspaceRevisionMode", c.WorkspaceRevisionMode, WorkspaceRevisionRestore},
		)
	}
	for _, tc := range checks {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
}

func TestDefaults(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if c.MaxDepth != 0 {
		t.Errorf("MaxDepth = %d, want 0", c.MaxDepth)
	}
	if c.Depth != 0 {
		t.Errorf("Depth = %d, want 0", c.Depth)
	}
	if c.MaxConcurrent != 0 {
		t.Errorf("MaxConcurrent = %d, want 0", c.MaxConcurrent)
	}
	if c.MaxAgents != 0 {
		t.Errorf("MaxAgents = %d, want 0", c.MaxAgents)
	}
	if c.ShTimeout != 600 {
		t.Errorf("ShTimeout = %d, want 600", c.ShTimeout)
	}
	if c.OutputTruncate != 20480 {
		t.Errorf("OutputTruncate = %d, want 20480", c.OutputTruncate)
	}
	if c.DataDir != ".quine/" {
		t.Errorf("DataDir = %q, want %q", c.DataDir, ".quine/")
	}
	if c.Shell != "/bin/sh" {
		t.Errorf("Shell = %q, want /bin/sh", c.Shell)
	}
	if c.MaxTurns != 0 {
		t.Errorf("MaxTurns = %d, want 0 (disabled)", c.MaxTurns)
	}
	if c.TurnExhaustionPolicy != TurnExhaustionHardFail {
		t.Errorf("TurnExhaustionPolicy = %q, want %q", c.TurnExhaustionPolicy, TurnExhaustionHardFail)
	}
	if c.PromptMetaphor != PromptMetaphorOff {
		t.Errorf("PromptMetaphor = %q, want %q", c.PromptMetaphor, PromptMetaphorOff)
	}
	if c.ContextWindow != 128_000 {
		t.Errorf("ContextWindow = %d, want 128000", c.ContextWindow)
	}
	if c.AnchorMemoryEnabled {
		t.Errorf("AnchorMemoryEnabled = true, want false by default")
	}
	if c.WorkspaceEnabled {
		t.Errorf("WorkspaceEnabled = true, want false by default")
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionNone {
		t.Errorf("WorkspaceRevisionMode = %q, want %q", c.WorkspaceRevisionMode, WorkspaceRevisionNone)
	}
	if c.WorkspaceRoot != "" || c.Workspace != "" {
		t.Errorf("workspace fields should default empty, got root=%q workspace=%q", c.WorkspaceRoot, c.Workspace)
	}
	if c.SessionID == "" {
		t.Error("SessionID should be auto-generated, got empty")
	}
	parts := strings.Split(c.SessionID, "_")
	if len(parts) != 3 {
		t.Fatalf("SessionID = %q, want <YYYYMMDD-HHMMSS>_<ppid>_<pid>", c.SessionID)
	}
	if _, err := time.Parse("20060102-150405", parts[0]); err != nil {
		t.Errorf("SessionID start_time = %q, want YYYYMMDD-HHMMSS: %v", parts[0], err)
	}
	if parts[1] != strconv.Itoa(os.Getppid()) {
		t.Errorf("SessionID ppid = %q, want %d", parts[1], os.Getppid())
	}
	if parts[2] != strconv.Itoa(os.Getpid()) {
		t.Errorf("SessionID pid = %q, want %d", parts[2], os.Getpid())
	}
	if c.TapeID == "" {
		t.Error("TapeID should be auto-generated, got empty")
	}
	if c.TapeID != "0001" {
		t.Errorf("TapeID = %q, want %q", c.TapeID, "0001")
	}
	if c.TapeID == c.SessionID {
		t.Errorf("TapeID = %q, want a distinct incarnation ID from SessionID", c.TapeID)
	}
}

func TestTapeIDAutoIncrementsPerSession(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	dataDir := t.TempDir()
	sessionID := "20260312-120000_123_456"
	tapeDir := filepath.Join(dataDir, "tapes", sessionID)
	if err := os.MkdirAll(tapeDir, 0o755); err != nil {
		t.Fatalf("mkdir tape dir: %v", err)
	}
	for _, name := range []string{"0001.jsonl", "0001.log.yaml", "0002.jsonl", "ignore.txt"} {
		if err := os.WriteFile(filepath.Join(tapeDir, name), []byte("x"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	os.Setenv("QUINE_DATA_DIR", dataDir)
	os.Setenv("QUINE_SESSION_ID", sessionID)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if c.TapeID != "0003" {
		t.Errorf("TapeID = %q, want %q", c.TapeID, "0003")
	}
}

// --- Required field validation tests ---

func TestMissingRequiredField(t *testing.T) {
	tests := []struct {
		name    string
		missing string
		setVars map[string]string
	}{
		{
			name:    "ModelID",
			missing: "QUINE_MODEL_ID",
			setVars: map[string]string{
				"QUINE_API_TYPE": "openai",
				"QUINE_API_BASE": "https://api.openai.com",
				"QUINE_API_KEY":  "sk-test",
			},
		},
		{
			name:    "APIType",
			missing: "QUINE_API_TYPE",
			setVars: map[string]string{
				"QUINE_MODEL_ID": "some-model",
				"QUINE_API_BASE": "https://example.com",
				"QUINE_API_KEY":  "sk-test",
			},
		},
		{
			name:    "APIBase",
			missing: "QUINE_API_BASE",
			setVars: map[string]string{
				"QUINE_MODEL_ID": "some-model",
				"QUINE_API_TYPE": "openai",
				"QUINE_API_KEY":  "sk-test",
			},
		},
		{
			name:    "APIKey",
			missing: "QUINE_API_KEY",
			setVars: map[string]string{
				"QUINE_MODEL_ID": "some-model",
				"QUINE_API_TYPE": "openai",
				"QUINE_API_BASE": "https://example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.setVars {
				os.Setenv(k, v)
			}
			_, err := Load()
			if err == nil {
				t.Fatalf("expected error for missing %s", tt.missing)
			}
			if !strings.Contains(err.Error(), tt.missing) {
				t.Errorf("error should mention %s, got: %v", tt.missing, err)
			}
		})
	}
}

func TestUnsupportedAPIType(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "some-model")
	os.Setenv("QUINE_API_TYPE", "gemini")
	os.Setenv("QUINE_API_BASE", "https://example.com")
	os.Setenv("QUINE_API_KEY", "sk-test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for unsupported API type")
	}
	if !strings.Contains(err.Error(), "unsupported QUINE_API_TYPE") {
		t.Errorf("error should mention unsupported QUINE_API_TYPE, got: %v", err)
	}
}

func TestCodexOAuthAPIKeyAllowed(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "gpt-5.2-codex")
	os.Setenv("QUINE_API_TYPE", "openai-responses")
	os.Setenv("QUINE_API_BASE", "https://api.openai.com")
	os.Setenv("QUINE_API_KEY", "codex-oauth") // Special sentinel value triggers OAuth

	_, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
}

// --- Third-party provider test ---

func TestThirdPartyProvider(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "kimi-k2.5")
	os.Setenv("QUINE_API_TYPE", "openai")
	os.Setenv("QUINE_API_BASE", "https://api.moonshot.ai/v1")
	os.Setenv("QUINE_API_KEY", "sk-moonshot-test")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ModelID != "kimi-k2.5" {
		t.Errorf("ModelID = %q, want kimi-k2.5", c.ModelID)
	}
	if c.Provider != "openai" {
		t.Errorf("Provider = %q, want openai", c.Provider)
	}
	if c.APIBase != "https://api.moonshot.ai/v1" {
		t.Errorf("APIBase = %q, want https://api.moonshot.ai/v1", c.APIBase)
	}
	if c.APIKey != "sk-moonshot-test" {
		t.Errorf("APIKey = %q, want sk-moonshot-test", c.APIKey)
	}
}

func TestDepthExceeded(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_DEPTH", "3")
	os.Setenv("QUINE_DEPTH", "3")

	_, err := Load()
	if !errors.Is(err, ErrDepthExceeded) {
		t.Errorf("expected ErrDepthExceeded, got: %v", err)
	}
}

func TestDepthDisabledWhenMaxDepthZero(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_DEPTH", "0")
	os.Setenv("QUINE_DEPTH", "999")

	_, err := Load()
	if err != nil {
		t.Fatalf("Load() should allow QUINE_DEPTH when QUINE_MAX_DEPTH=0, got: %v", err)
	}
}

func TestContextWindow_ExplicitOverride(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_CONTEXT_WINDOW", "500000")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ContextWindow != 500_000 {
		t.Errorf("ContextWindow = %d, want 500000", c.ContextWindow)
	}
}

// --- ChildEnv / ExecEnv tests ---

func TestChildEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SESSION_ID", "parent-uuid")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}

	// Depth should be incremented
	childDepth, err := strconv.Atoi(m["QUINE_DEPTH"])
	if err != nil {
		t.Fatalf("parsing child QUINE_DEPTH: %v", err)
	}
	if childDepth != c.Depth+1 {
		t.Errorf("child QUINE_DEPTH = %d, want %d", childDepth, c.Depth+1)
	}

	// Parent session should be current session
	if m["QUINE_PARENT_SESSION"] != "parent-uuid" {
		t.Errorf("QUINE_PARENT_SESSION = %q, want %q", m["QUINE_PARENT_SESSION"], "parent-uuid")
	}

	// Session ID should NOT be present — each child generates its own
	if _, hasSessionID := m["QUINE_SESSION_ID"]; hasSessionID {
		t.Error("ChildEnv should NOT include QUINE_SESSION_ID (children generate their own)")
	}
	if _, hasTapeID := m["QUINE_TAPE_ID"]; hasTapeID {
		t.Error("ChildEnv should NOT include QUINE_TAPE_ID (children generate their own)")
	}

	// All 4 required fields passed through
	if m["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want %q", m["QUINE_MODEL_ID"], "claude-sonnet-4-20250514")
	}
	if m["QUINE_API_TYPE"] != "anthropic" {
		t.Errorf("QUINE_API_TYPE = %q, want anthropic", m["QUINE_API_TYPE"])
	}
	if m["QUINE_API_BASE"] != "https://api.anthropic.com" {
		t.Errorf("QUINE_API_BASE = %q, want https://api.anthropic.com", m["QUINE_API_BASE"])
	}
	if m["QUINE_API_KEY"] != "sk-test-key" {
		t.Errorf("QUINE_API_KEY = %q, want sk-test-key", m["QUINE_API_KEY"])
	}

	// MaxTurns should be propagated
	if m["QUINE_MAX_TURNS"] != strconv.Itoa(c.MaxTurns) {
		t.Errorf("QUINE_MAX_TURNS = %q, want %q", m["QUINE_MAX_TURNS"], strconv.Itoa(c.MaxTurns))
	}

	// All expected keys present
	expectedKeys := []string{
		"QUINE_MODEL_ID", "QUINE_API_TYPE", "QUINE_API_BASE", "QUINE_API_KEY",
		"QUINE_MAX_DEPTH", "QUINE_DEPTH", "QUINE_PARENT_SESSION",
		"QUINE_MAX_CONCURRENT", "QUINE_MAX_AGENTS", "QUINE_SH_TIMEOUT", "QUINE_OUTPUT_TRUNCATE",
		"QUINE_DATA_DIR", "QUINE_SHELL", "QUINE_MAX_TURNS", "QUINE_TURN_EXHAUSTION_POLICY",
		"QUINE_PROMPT_METAPHOR",
		"QUINE_CONTEXT_WINDOW",
	}
	for _, k := range expectedKeys {
		if _, ok := m[k]; !ok {
			t.Errorf("missing key %s in ChildEnv", k)
		}
	}
}

func TestChildEnv_PropagatesAnchorMemoryFlag(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_ANCHOR_MEMORY", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	found := false
	for _, e := range env {
		if e == "QUINE_ANCHOR_MEMORY=1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ChildEnv should propagate QUINE_ANCHOR_MEMORY=1 when enabled")
	}
}

func TestExecEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SESSION_ID", "exec-parent-uuid")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ExecEnv("build the project")
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}

	// Depth should be reset to 0
	if m["QUINE_DEPTH"] != "0" {
		t.Errorf("QUINE_DEPTH = %q, want 0", m["QUINE_DEPTH"])
	}

	// Original intent should be set
	if m["QUINE_ORIGINAL_INTENT"] != "build the project" {
		t.Errorf("QUINE_ORIGINAL_INTENT = %q, want %q", m["QUINE_ORIGINAL_INTENT"], "build the project")
	}
	if m["QUINE_SESSION_ID"] != c.SessionID {
		t.Errorf("QUINE_SESSION_ID = %q, want %q", m["QUINE_SESSION_ID"], c.SessionID)
	}
	if m["QUINE_PARENT_SESSION"] != c.ParentSession {
		t.Errorf("QUINE_PARENT_SESSION = %q, want %q", m["QUINE_PARENT_SESSION"], c.ParentSession)
	}
	if _, hasTapeID := m["QUINE_TAPE_ID"]; hasTapeID {
		t.Error("ExecEnv should NOT include QUINE_TAPE_ID (exec creates a fresh tape)")
	}

	if _, ok := m["QUINE_CONFIG_DIR"]; ok {
		t.Errorf("QUINE_CONFIG_DIR should not be set when empty")
	}
}

func TestConfigDirPassthrough(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_CONFIG_DIR", "/tmp/quine-config")
	t.Cleanup(func() {
		os.Unsetenv("QUINE_CONFIG_DIR")
	})

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	found := false
	for _, e := range env {
		if e == "QUINE_CONFIG_DIR=/tmp/quine-config" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("QUINE_CONFIG_DIR should be passed through")
	}
}

func TestWisdomChildEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	os.Setenv("QUINE_WISDOM_STATE", "processing chunk 5")
	t.Cleanup(func() {
		os.Unsetenv("QUINE_WISDOM_STATE")
	})

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	found := false
	for _, e := range env {
		if e == "QUINE_WISDOM_STATE=processing chunk 5" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ChildEnv should include QUINE_WISDOM_STATE")
	}
}

func TestWisdomIgnoresEmptyValues(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	os.Setenv("QUINE_WISDOM_EMPTY", "")
	os.Setenv("QUINE_WISDOM_VALID", "has value")
	t.Cleanup(func() {
		os.Unsetenv("QUINE_WISDOM_EMPTY")
		os.Unsetenv("QUINE_WISDOM_VALID")
	})

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(c.Wisdom) != 1 {
		t.Errorf("Wisdom length = %d, want 1 (empty values ignored)", len(c.Wisdom))
	}
	if _, exists := c.Wisdom["EMPTY"]; exists {
		t.Error("Wisdom should not contain EMPTY key")
	}
	if c.Wisdom["VALID"] != "has value" {
		t.Errorf("Wisdom[VALID] = %q, want %q", c.Wisdom["VALID"], "has value")
	}
}

func TestInvalidTurnExhaustionPolicy(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_TURNS", "1")
	os.Setenv("QUINE_TURN_EXHAUSTION_POLICY", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid QUINE_TURN_EXHAUSTION_POLICY")
	}
	if !strings.Contains(err.Error(), "QUINE_TURN_EXHAUSTION_POLICY") {
		t.Errorf("error should mention QUINE_TURN_EXHAUSTION_POLICY, got: %v", err)
	}
}

func TestTurnExhaustionPolicyIgnoredWhenBudgetDisabled(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_TURNS", "0")
	os.Setenv("QUINE_TURN_EXHAUSTION_POLICY", "invalid")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() should ignore QUINE_TURN_EXHAUSTION_POLICY when QUINE_MAX_TURNS=0, got: %v", err)
	}
	if c.TurnExhaustionPolicy != TurnExhaustionHardFail {
		t.Errorf("TurnExhaustionPolicy = %q, want %q when budget disabled", c.TurnExhaustionPolicy, TurnExhaustionHardFail)
	}
}

func TestInvalidPromptMetaphor(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_METAPHOR", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid QUINE_PROMPT_METAPHOR")
	}
	if !strings.Contains(err.Error(), "QUINE_PROMPT_METAPHOR") {
		t.Errorf("error should mention QUINE_PROMPT_METAPHOR, got: %v", err)
	}
}

func TestWorkspaceConfigImplicitEnable(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	subdir := root + "/subdir"
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	subdirReal, err := filepath.EvalSymlinks(subdir)
	if err != nil {
		t.Fatalf("realpath subdir: %v", err)
	}
	os.Setenv("QUINE_WORKSPACE", subdir)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.WorkspaceEnabled {
		t.Fatal("workspace configuration should implicitly enable workspace physics")
	}
	if c.WorkspaceRoot != subdirReal {
		t.Errorf("WorkspaceRoot = %q, want %q", c.WorkspaceRoot, subdirReal)
	}
	if c.Workspace != subdirReal {
		t.Errorf("Workspace = %q, want %q", c.Workspace, subdirReal)
	}
	if c.WorkspaceBackend != "overlay" {
		t.Errorf("WorkspaceBackend = %q, want overlay", c.WorkspaceBackend)
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionRestore {
		t.Errorf("WorkspaceRevisionMode = %q, want restore", c.WorkspaceRevisionMode)
	}
}

func TestWorkspaceDefaultsToPwdWithinRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	nested := filepath.Join(root, "nested", "deep")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	nestedReal, err := filepath.EvalSymlinks(nested)
	if err != nil {
		t.Fatalf("realpath nested: %v", err)
	}
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatalf("chdir nested: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	os.Setenv("QUINE_WORKSPACE_ROOT", root)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.Workspace != nestedReal {
		t.Fatalf("Workspace = %q, want %q", c.Workspace, nestedReal)
	}
}

func TestWorkspaceMustStayWithinRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	other := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", other)

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace/root validation error")
	}
	if !strings.Contains(err.Error(), "must be within workspace root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceDirectBackendAccepted(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "direct")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.WorkspaceBackend != "direct" {
		t.Fatalf("WorkspaceBackend = %q, want direct", c.WorkspaceBackend)
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionRestore {
		t.Fatalf("WorkspaceRevisionMode = %q, want restore", c.WorkspaceRevisionMode)
	}
}

func TestWorkspaceRevisionModeAccepted(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "restore")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionRestore {
		t.Fatalf("WorkspaceRevisionMode = %q, want restore", c.WorkspaceRevisionMode)
	}
	if !c.CanRestoreWorld() {
		t.Fatal("restore mode should enable restore_world")
	}
}

func TestWorkspaceRevisionModeRejectsUnknownValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "mystery")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace revision mode validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_REVISION_MODE") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceRevisionModeRequiresWorkspace(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "restore")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace revision mode requires workspace error")
	}
	if !strings.Contains(err.Error(), "requires workspace physics") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceBackendRejectsUnknownValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "mystery")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace backend validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_BACKEND") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceSourceIsRejected(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_SOURCE", root)

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace source removal error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_SOURCE has been removed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceDataDirMustStayOutsideRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_DATA_DIR", filepath.Join(root, ".quine"))

	_, err := Load()
	if err == nil {
		t.Fatal("expected QUINE_DATA_DIR outside workspace root validation error")
	}
	if !strings.Contains(err.Error(), "must be outside workspace root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceUnsupportedOnNonLinux(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS == "linux" {
		t.Skip("workspace physics are supported on Linux")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)

	_, err := Load()
	if err == nil {
		t.Fatal("expected non-Linux workspace error")
	}
	if !strings.Contains(err.Error(), "only supported on Linux") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChildEnvIncludesWorkspacePhysics(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	child := root + "/child"
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir child: %v", err)
	}
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("realpath root: %v", err)
	}
	childReal, err := filepath.EvalSymlinks(child)
	if err != nil {
		t.Fatalf("realpath child: %v", err)
	}
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", child)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}

	if m["QUINE_WORKSPACE_ROOT"] != rootReal {
		t.Fatalf("QUINE_WORKSPACE_ROOT = %q, want %q", m["QUINE_WORKSPACE_ROOT"], rootReal)
	}
	if m["QUINE_WORKSPACE"] != childReal {
		t.Fatalf("QUINE_WORKSPACE = %q, want %q", m["QUINE_WORKSPACE"], childReal)
	}
	if m["QUINE_WORKSPACE_BACKEND"] != "overlay" {
		t.Fatalf("QUINE_WORKSPACE_BACKEND = %q, want overlay", m["QUINE_WORKSPACE_BACKEND"])
	}
	if m["QUINE_WORKSPACE_REVISION_MODE"] != "full" {
		t.Fatalf("QUINE_WORKSPACE_REVISION_MODE = %q, want full", m["QUINE_WORKSPACE_REVISION_MODE"])
	}
	if m["QUINE_WORKSPACE_SESSION"] != c.WorkspaceSession {
		t.Fatalf("QUINE_WORKSPACE_SESSION = %q, want %q", m["QUINE_WORKSPACE_SESSION"], c.WorkspaceSession)
	}
	if m["QUINE_WORKSPACE_OWNER"] != "true" {
		t.Fatalf("QUINE_WORKSPACE_OWNER = %q, want true", m["QUINE_WORKSPACE_OWNER"])
	}
}
