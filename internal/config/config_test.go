package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
)

// envVars is the full list of environment variables we manage in tests.
var envVars = []string{
	"QUINE_MODEL_ID",
	"QUINE_API_TYPE",
	"QUINE_API_BASE",
	"QUINE_API_KEY",
	"QUINE_MAX_DEPTH",
	"QUINE_DEPTH",
	"QUINE_SESSION_ID",
	"QUINE_PARENT_SESSION",
	"QUINE_MAX_CONCURRENT",
	"QUINE_SH_TIMEOUT",
	"QUINE_OUTPUT_TRUNCATE",
	"QUINE_DATA_DIR",
	"QUINE_SHELL",
	"QUINE_MAX_TURNS",
	"QUINE_CONTEXT_WINDOW",
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
	os.Setenv("QUINE_MAX_DEPTH", "10")
	os.Setenv("QUINE_DEPTH", "3")
	os.Setenv("QUINE_SESSION_ID", "my-session")
	os.Setenv("QUINE_PARENT_SESSION", "parent-session")
	os.Setenv("QUINE_MAX_CONCURRENT", "50")
	os.Setenv("QUINE_SH_TIMEOUT", "60")
	os.Setenv("QUINE_OUTPUT_TRUNCATE", "4096")
	os.Setenv("QUINE_DATA_DIR", "/tmp/data")
	os.Setenv("QUINE_SHELL", "/bin/sh")
	os.Setenv("QUINE_MAX_TURNS", "30")

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
		{"ParentSession", c.ParentSession, "parent-session"},
		{"MaxConcurrent", c.MaxConcurrent, 50},
		{"ShTimeout", c.ShTimeout, 60},
		{"OutputTruncate", c.OutputTruncate, 4096},
		{"DataDir", c.DataDir, "/tmp/data"},
		{"Shell", c.Shell, "/bin/sh"},
		{"MaxTurns", c.MaxTurns, 30},
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

	if c.MaxDepth != 5 {
		t.Errorf("MaxDepth = %d, want 5", c.MaxDepth)
	}
	if c.Depth != 0 {
		t.Errorf("Depth = %d, want 0", c.Depth)
	}
	if c.MaxConcurrent != 20 {
		t.Errorf("MaxConcurrent = %d, want 20", c.MaxConcurrent)
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
	if c.MaxTurns != 20 {
		t.Errorf("MaxTurns = %d, want 20", c.MaxTurns)
	}
	if c.ContextWindow != 128_000 {
		t.Errorf("ContextWindow = %d, want 128000", c.ContextWindow)
	}
	if c.SessionID == "" {
		t.Error("SessionID should be auto-generated, got empty")
	}
	if len(c.SessionID) != 36 {
		t.Errorf("SessionID length = %d, want 36", len(c.SessionID))
	}
}

// --- Required field validation tests ---

func TestMissing_ModelID(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_API_TYPE", "openai")
	os.Setenv("QUINE_API_BASE", "https://api.openai.com")
	os.Setenv("QUINE_API_KEY", "sk-test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing QUINE_MODEL_ID")
	}
	if !strings.Contains(err.Error(), "QUINE_MODEL_ID") {
		t.Errorf("error should mention QUINE_MODEL_ID, got: %v", err)
	}
}

func TestMissing_APIType(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "some-model")
	os.Setenv("QUINE_API_BASE", "https://example.com")
	os.Setenv("QUINE_API_KEY", "sk-test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing QUINE_API_TYPE")
	}
	if !strings.Contains(err.Error(), "QUINE_API_TYPE") {
		t.Errorf("error should mention QUINE_API_TYPE, got: %v", err)
	}
}

func TestMissing_APIBase(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "some-model")
	os.Setenv("QUINE_API_TYPE", "openai")
	os.Setenv("QUINE_API_KEY", "sk-test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing QUINE_API_BASE")
	}
	if !strings.Contains(err.Error(), "QUINE_API_BASE") {
		t.Errorf("error should mention QUINE_API_BASE, got: %v", err)
	}
}

func TestMissing_APIKey(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "some-model")
	os.Setenv("QUINE_API_TYPE", "openai")
	os.Setenv("QUINE_API_BASE", "https://example.com")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing QUINE_API_KEY")
	}
	if !strings.Contains(err.Error(), "QUINE_API_KEY") {
		t.Errorf("error should mention QUINE_API_KEY, got: %v", err)
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

func TestDepthExceeded_Greater(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_DEPTH", "3")
	os.Setenv("QUINE_DEPTH", "10")

	_, err := Load()
	if !errors.Is(err, ErrDepthExceeded) {
		t.Errorf("expected ErrDepthExceeded, got: %v", err)
	}
}

// --- Context window tests ---

func TestContextWindow_Default(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ContextWindow != 128_000 {
		t.Errorf("ContextWindow = %d, want 128000", c.ContextWindow)
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
		"QUINE_MAX_CONCURRENT", "QUINE_SH_TIMEOUT", "QUINE_OUTPUT_TRUNCATE",
		"QUINE_DATA_DIR", "QUINE_SHELL", "QUINE_MAX_TURNS",
		"QUINE_CONTEXT_WINDOW",
	}
	for _, k := range expectedKeys {
		if _, ok := m[k]; !ok {
			t.Errorf("missing key %s in ChildEnv", k)
		}
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

	env, err := c.ExecEnv("build the project", 1024)
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

	// Stdin offset should be set
	if m["QUINE_STDIN_OFFSET"] != "1024" {
		t.Errorf("QUINE_STDIN_OFFSET = %q, want 1024", m["QUINE_STDIN_OFFSET"])
	}
}

func TestUUIDFormat(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// UUID v4 format: xxxxxxxx-xxxx-4xxx-[89ab]xxx-xxxxxxxxxxxx
	parts := strings.Split(c.SessionID, "-")
	if len(parts) != 5 {
		t.Fatalf("UUID should have 5 parts, got %d: %q", len(parts), c.SessionID)
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		t.Errorf("UUID part lengths wrong: %q", c.SessionID)
	}
	if parts[2][0] != '4' {
		t.Errorf("UUID version nibble = %c, want '4'", parts[2][0])
	}
	v := parts[3][0]
	if v != '8' && v != '9' && v != 'a' && v != 'b' {
		t.Errorf("UUID variant nibble = %c, want [89ab]", v)
	}
}

func TestWisdomLoading(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	os.Setenv("QUINE_WISDOM_SUMMARY", "User prefers concise answers")
	os.Setenv("QUINE_WISDOM_CONTEXT", "Working on Go project")
	t.Cleanup(func() {
		os.Unsetenv("QUINE_WISDOM_SUMMARY")
		os.Unsetenv("QUINE_WISDOM_CONTEXT")
	})

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(c.Wisdom) != 2 {
		t.Errorf("Wisdom length = %d, want 2", len(c.Wisdom))
	}
	if c.Wisdom["SUMMARY"] != "User prefers concise answers" {
		t.Errorf("Wisdom[SUMMARY] = %q, want %q", c.Wisdom["SUMMARY"], "User prefers concise answers")
	}
	if c.Wisdom["CONTEXT"] != "Working on Go project" {
		t.Errorf("Wisdom[CONTEXT] = %q, want %q", c.Wisdom["CONTEXT"], "Working on Go project")
	}
}

func TestWisdomEmpty(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if c.Wisdom == nil {
		t.Error("Wisdom should not be nil")
	}
	if len(c.Wisdom) != 0 {
		t.Errorf("Wisdom length = %d, want 0", len(c.Wisdom))
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

// --- APIModelID tests ---

func TestAPIModelID(t *testing.T) {
	c := &Config{ModelID: "claude-sonnet-4-20250514"}
	if c.APIModelID() != "claude-sonnet-4-20250514" {
		t.Errorf("APIModelID() = %q, want %q", c.APIModelID(), "claude-sonnet-4-20250514")
	}
}
