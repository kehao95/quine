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
	"QUINE_API_BASE",
	"QUINE_PROVIDER",
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
	// Provider-specific API keys
	"ANTHROPIC_API_KEY",
	"OPENAI_API_KEY",
	"GOOGLE_GENERATIVE_AI_API_KEY",
	"GEMINI_API_KEY",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_REGION",
	"AWS_DEFAULT_REGION",
	"AWS_SESSION_TOKEN",
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

func setMinimal(t *testing.T) {
	t.Helper()
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-20250514")
	os.Setenv("ANTHROPIC_API_KEY", "sk-test-key")
}

func TestHappyPath(t *testing.T) {
	clearEnv(t)
	setMinimal(t)
	os.Setenv("QUINE_API_BASE", "https://api.example.com")
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
		{"APIBase", c.APIBase, "https://api.example.com"},
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
	setMinimal(t)

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
	if c.SessionID == "" {
		t.Error("SessionID should be auto-generated, got empty")
	}
	// Validate UUID format: 8-4-4-4-12 hex chars
	if len(c.SessionID) != 36 {
		t.Errorf("SessionID length = %d, want 36", len(c.SessionID))
	}
}

func TestDefaultModelID(t *testing.T) {
	clearEnv(t)
	os.Setenv("ANTHROPIC_API_KEY", "sk-test")

	c, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ModelID != "claude-sonnet-4-5-20250929" {
		t.Errorf("ModelID = %q, want %q", c.ModelID, "claude-sonnet-4-5-20250929")
	}
}

func TestMissingAPIKey(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-20250514")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing API_KEY")
	}
	// Should mention provider-specific env var
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Errorf("error should mention ANTHROPIC_API_KEY, got: %v", err)
	}
}

func TestBedrockNoAPIKeyRequired(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "anthropic.claude-3-5-sonnet-20241022-v2:0")
	// Bedrock uses AWS credentials, not API key

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.Provider != "bedrock" {
		t.Errorf("Provider = %q, want bedrock", c.Provider)
	}
	if c.APIKey != "" {
		t.Errorf("APIKey = %q, want empty", c.APIKey)
	}
}

func TestDepthExceeded(t *testing.T) {
	clearEnv(t)
	setMinimal(t)
	os.Setenv("QUINE_MAX_DEPTH", "3")
	os.Setenv("QUINE_DEPTH", "3")

	_, err := Load()
	if !errors.Is(err, ErrDepthExceeded) {
		t.Errorf("expected ErrDepthExceeded, got: %v", err)
	}
}

func TestDepthExceeded_Greater(t *testing.T) {
	clearEnv(t)
	setMinimal(t)
	os.Setenv("QUINE_MAX_DEPTH", "3")
	os.Setenv("QUINE_DEPTH", "10")

	_, err := Load()
	if !errors.Is(err, ErrDepthExceeded) {
		t.Errorf("expected ErrDepthExceeded, got: %v", err)
	}
}

func TestProviderAutoDetect(t *testing.T) {
	cases := []struct {
		model    string
		provider string
		envKey   string
		envVar   string
	}{
		{"claude-sonnet-4-20250514", "anthropic", "ANTHROPIC_API_KEY", "sk-test"},
		{"claude-3-5-haiku-20241022", "anthropic", "ANTHROPIC_API_KEY", "sk-test"},
		{"gpt-4o", "openai", "OPENAI_API_KEY", "sk-test"},
		{"gpt-4-turbo", "openai", "OPENAI_API_KEY", "sk-test"},
		{"gemini-2.0-flash", "google", "GOOGLE_GENERATIVE_AI_API_KEY", "sk-test"},
		{"gemini-1.5-pro", "google", "GOOGLE_GENERATIVE_AI_API_KEY", "sk-test"},
		// Bedrock model IDs (must detect "anthropic." prefix → bedrock, not anthropic)
		{"anthropic.claude-3-5-sonnet-20241022-v2:0", "bedrock", "", ""},
		{"anthropic.claude-sonnet-4-20250514-v1:0", "bedrock", "", ""},
		{"anthropic.claude-3-haiku-20240307-v1:0", "bedrock", "", ""},
		{"us.anthropic.claude-3-5-sonnet-20241022-v2:0", "bedrock", "", ""},
		{"eu.anthropic.claude-sonnet-4-20250514-v1:0", "bedrock", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.model, func(t *testing.T) {
			clearEnv(t)
			os.Setenv("QUINE_MODEL_ID", tc.model)
			if tc.envKey != "" {
				os.Setenv(tc.envKey, tc.envVar)
			}

			c, err := Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}
			if c.Provider != tc.provider {
				t.Errorf("Provider = %q, want %q", c.Provider, tc.provider)
			}
		})
	}
}

func TestProviderAutoDetectUnknown(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "llama-3-70b")
	os.Setenv("ANTHROPIC_API_KEY", "sk-test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
	// For models not in registry, falls back to legacy detection
	if !strings.Contains(err.Error(), "cannot auto-detect provider") && !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention auto-detect or not found, got: %v", err)
	}
}

func TestProviderExplicit(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "llama-3-70b")
	os.Setenv("OPENAI_API_KEY", "sk-test")
	os.Setenv("QUINE_PROVIDER", "openai")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.Provider != "openai" {
		t.Errorf("Provider = %q, want %q", c.Provider, "openai")
	}
}

func TestChildEnv(t *testing.T) {
	clearEnv(t)
	setMinimal(t)
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
	// via config.Load() to avoid tape file collisions when multiple
	// children are spawned from a single sh command.
	if _, hasSessionID := m["QUINE_SESSION_ID"]; hasSessionID {
		t.Error("ChildEnv should NOT include QUINE_SESSION_ID (children generate their own)")
	}

	// Model inherited
	if m["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want %q", m["QUINE_MODEL_ID"], "claude-sonnet-4-20250514")
	}

	// API key passed through provider-specific env var
	if m["ANTHROPIC_API_KEY"] != "sk-test-key" {
		t.Errorf("ANTHROPIC_API_KEY = %q, want %q", m["ANTHROPIC_API_KEY"], "sk-test-key")
	}

	// MaxTurns should be propagated
	if m["QUINE_MAX_TURNS"] != strconv.Itoa(c.MaxTurns) {
		t.Errorf("QUINE_MAX_TURNS = %q, want %q", m["QUINE_MAX_TURNS"], strconv.Itoa(c.MaxTurns))
	}

	// All expected keys present
	expectedKeys := []string{
		"QUINE_MODEL_ID", "QUINE_API_BASE", "QUINE_PROVIDER",
		"QUINE_MAX_DEPTH", "QUINE_DEPTH", "QUINE_PARENT_SESSION",
		"QUINE_MAX_CONCURRENT", "QUINE_SH_TIMEOUT", "QUINE_OUTPUT_TRUNCATE",
		"QUINE_DATA_DIR", "QUINE_SHELL", "QUINE_MAX_TURNS",
	}
	for _, k := range expectedKeys {
		if _, ok := m[k]; !ok {
			t.Errorf("missing key %s in ChildEnv", k)
		}
	}
}

func TestUUIDFormat(t *testing.T) {
	clearEnv(t)
	setMinimal(t)

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
	// Version nibble should be '4'
	if parts[2][0] != '4' {
		t.Errorf("UUID version nibble = %c, want '4'", parts[2][0])
	}
	// Variant nibble should be 8, 9, a, or b
	v := parts[3][0]
	if v != '8' && v != '9' && v != 'a' && v != 'b' {
		t.Errorf("UUID variant nibble = %c, want [89ab]", v)
	}
}

// TestContextWindowFromRegistry tests that context window is set from registry
func TestContextWindowFromRegistry(t *testing.T) {
	clearEnv(t)
	setMinimal(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// claude-sonnet-4-20250514 should have 200000 context window
	if c.ContextWindow != 200000 {
		t.Errorf("ContextWindow = %d, want 200000", c.ContextWindow)
	}
}

func TestWisdomLoading(t *testing.T) {
	clearEnv(t)
	setMinimal(t)

	// Set some wisdom env vars
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

	// Check wisdom is loaded
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
	setMinimal(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Wisdom should be empty but not nil
	if c.Wisdom == nil {
		t.Error("Wisdom should not be nil")
	}
	if len(c.Wisdom) != 0 {
		t.Errorf("Wisdom length = %d, want 0", len(c.Wisdom))
	}
}

func TestWisdomChildEnv(t *testing.T) {
	clearEnv(t)
	setMinimal(t)

	// Set wisdom env vars
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

	// Check wisdom is passed through
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
	setMinimal(t)

	// Set wisdom with empty value
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

	// Empty values should be ignored
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
