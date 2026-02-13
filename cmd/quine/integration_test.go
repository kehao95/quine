package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/runtime"
	"github.com/kehao95/quine/internal/tape"
)

// ---------------------------------------------------------------------------
// Mock provider
// ---------------------------------------------------------------------------

// mockProvider implements llm.Provider and returns pre-programmed responses.
type mockProvider struct {
	mu        sync.Mutex
	responses []tape.Message
	callIndex int
}

func newMockProvider(responses []tape.Message) *mockProvider {
	return &mockProvider{responses: responses}
}

func (m *mockProvider) Generate(messages []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.callIndex >= len(m.responses) {
		return tape.Message{}, llm.Usage{}, errors.New("mock provider: no more responses")
	}
	msg := m.responses[m.callIndex]
	m.callIndex++
	return msg, llm.Usage{InputTokens: 100, OutputTokens: 50}, nil
}

func (m *mockProvider) ContextWindowSize() int {
	return 200_000
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// testConfig builds a Config for testing without touching real env vars.
func testConfig(t *testing.T, tapeDir string) *config.Config {
	t.Helper()
	return &config.Config{
		ModelID:        "claude-test",
		APIKey:         "test-key",
		APIBase:        "",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          0,
		SessionID:      "test-session-0001",
		ParentSession:  "",
		MaxConcurrent:  20,
		ShTimeout:      30,
		OutputTruncate: 20480,
		DataDir:        tapeDir,
		Shell:          "/bin/sh",
	}
}

// assistantsh returns a tape.Message that calls the sh tool.
func assistantsh(id, command string) tape.Message {
	return tape.Message{
		Role: tape.RoleAssistant,
		ToolCalls: []tape.ToolCall{
			{
				ID:   id,
				Name: "sh",
				Arguments: map[string]any{
					"command": command,
				},
			},
		},
	}
}

// assistantExit returns a tape.Message that calls the exit tool.
func assistantExit(id string, status string, output, errMsg string) tape.Message {
	args := map[string]any{
		"status": status,
	}
	if output != "" {
		args["output"] = output
	}
	if errMsg != "" {
		args["stderr"] = errMsg
	}
	return tape.Message{
		Role: tape.RoleAssistant,
		ToolCalls: []tape.ToolCall{
			{
				ID:        id,
				Name:      "exit",
				Arguments: args,
			},
		},
	}
}

// parseTapeJSONL reads a JSONL tape file and returns parsed entries.
func parseTapeJSONL(t *testing.T, path string) []tape.TapeEntry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open tape JSONL: %v", err)
	}
	defer f.Close()

	var entries []tape.TapeEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry tape.TapeEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("unmarshal tape entry: %v (line: %s)", err, line)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}
	return entries
}

// ---------------------------------------------------------------------------
// Test 1: The Survival Test (§11.1)
// ---------------------------------------------------------------------------

func TestSurvivalTest(t *testing.T) {
	// Setup temp directories
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")
	workDir := filepath.Join(tmpDir, "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatalf("create work dir: %v", err)
	}

	cfg := testConfig(t, tapeDir)

	// Mock LLM: sh(echo 'world' > hello.txt), sh(cat hello.txt), exit(0)
	mock := newMockProvider([]tape.Message{
		assistantsh("call-1", "cd "+workDir+" && echo 'world' > hello.txt"),
		assistantsh("call-2", "cat "+filepath.Join(workDir, "hello.txt")),
		assistantExit("call-3", "success", "world", ""),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	exitCode := rt.Run("Create hello.txt with 'world' and read it back", "Begin.", nil)

	// Verify exit code
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Verify hello.txt was created
	helloPath := filepath.Join(workDir, "hello.txt")
	data, err := os.ReadFile(helloPath)
	if err != nil {
		t.Fatalf("hello.txt not created: %v", err)
	}
	content := strings.TrimSpace(string(data))
	if content != "world" {
		t.Errorf("hello.txt content = %q, want %q", content, "world")
	}

	// Verify tape JSONL file exists
	tapePath := filepath.Join(tapeDir, cfg.SessionID+".jsonl")
	if _, err := os.Stat(tapePath); os.IsNotExist(err) {
		t.Fatalf("tape JSONL file does not exist: %s", tapePath)
	}

	// Parse tape entries and verify structure
	entries := parseTapeJSONL(t, tapePath)
	if len(entries) == 0 {
		t.Fatal("tape JSONL has no entries")
	}

	// Verify entry types sequence:
	// meta, message(system), message(user), message(assistant/sh1),
	// tool_result(sh1), message(assistant/sh2), tool_result(sh2),
	// message(assistant/exit), outcome
	expectedTypes := []string{
		"meta",
		"message",     // system
		"message",     // user
		"message",     // assistant (sh call-1)
		"tool_result", // sh result 1
		"message",     // assistant (sh call-2)
		"tool_result", // sh result 2
		"message",     // assistant (exit)
		"outcome",     // session outcome
	}

	if len(entries) != len(expectedTypes) {
		t.Errorf("tape entry count = %d, want %d", len(entries), len(expectedTypes))
		for i, e := range entries {
			t.Logf("  entry[%d]: type=%s", i, e.Type)
		}
	} else {
		for i, want := range expectedTypes {
			if entries[i].Type != want {
				t.Errorf("entry[%d].Type = %q, want %q", i, entries[i].Type, want)
			}
		}
	}

	// Verify the last entry is an outcome with exit_code 0
	lastEntry := entries[len(entries)-1]
	if lastEntry.Type != "outcome" {
		t.Fatalf("last tape entry type = %q, want %q", lastEntry.Type, "outcome")
	}
	var outcome tape.SessionOutcome
	if err := json.Unmarshal(lastEntry.Data, &outcome); err != nil {
		t.Fatalf("unmarshal outcome: %v", err)
	}
	if outcome.ExitCode != 0 {
		t.Errorf("outcome exit_code = %d, want 0", outcome.ExitCode)
	}
	if outcome.TerminationMode != tape.TermExit {
		t.Errorf("outcome termination_mode = %q, want %q", outcome.TerminationMode, tape.TermExit)
	}
}

// ---------------------------------------------------------------------------
// Test 2: The Resilience Test (§11.3)
// ---------------------------------------------------------------------------

func TestResilienceTest(t *testing.T) {
	tmpDir := t.TempDir()
	tapeDir := filepath.Join(tmpDir, "tapes")

	cfg := testConfig(t, tapeDir)

	// Mock LLM: sh(cat nonexistent), then exit(1, stderr)
	mock := newMockProvider([]tape.Message{
		assistantsh("call-1", "cat /tmp/nonexistent_ghost.txt"),
		assistantExit("call-2", "failure", "", "File /tmp/nonexistent_ghost.txt does not exist"),
	})

	rt := runtime.NewWithProvider(cfg, mock)
	exitCode := rt.Run("Read /tmp/nonexistent_ghost.txt", "Begin.", nil)

	// Verify exit code is 1
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	// Verify tape JSONL exists and has an outcome
	tapePath := filepath.Join(tapeDir, cfg.SessionID+".jsonl")
	entries := parseTapeJSONL(t, tapePath)

	// Find the tool_result entry — it should contain an error indication
	var foundErrorResult bool
	for _, e := range entries {
		if e.Type == "tool_result" {
			var tr tape.ToolResult
			if err := json.Unmarshal(e.Data, &tr); err != nil {
				t.Fatalf("unmarshal tool_result: %v", err)
			}
			// The sh command should have failed (cat nonexistent file)
			if strings.Contains(tr.Content, "[EXIT CODE]") && !strings.Contains(tr.Content, "[EXIT CODE] 0") {
				foundErrorResult = true
			}
		}
	}
	if !foundErrorResult {
		t.Error("expected a tool_result with non-zero exit code from failed cat command")
	}

	// Verify outcome
	lastEntry := entries[len(entries)-1]
	if lastEntry.Type != "outcome" {
		t.Fatalf("last tape entry type = %q, want %q", lastEntry.Type, "outcome")
	}
	var outcome tape.SessionOutcome
	if err := json.Unmarshal(lastEntry.Data, &outcome); err != nil {
		t.Fatalf("unmarshal outcome: %v", err)
	}
	if outcome.ExitCode != 1 {
		t.Errorf("outcome exit_code = %d, want 1", outcome.ExitCode)
	}
	if outcome.Stderr != "File /tmp/nonexistent_ghost.txt does not exist" {
		t.Errorf("outcome stderr = %q, want error message about nonexistent file", outcome.Stderr)
	}
}

// ---------------------------------------------------------------------------
// Test 3: The Bomb Test (§11.4) — Depth limit enforcement
// ---------------------------------------------------------------------------

func TestBombTestDepthLimit(t *testing.T) {
	// Save and restore environment variables.
	envVars := []string{
		"QUINE_DEPTH", "QUINE_MAX_DEPTH", "QUINE_MODEL_ID",
		"QUINE_API_TYPE", "QUINE_API_BASE", "QUINE_API_KEY",
	}
	saved := make(map[string]string)
	for _, key := range envVars {
		saved[key] = os.Getenv(key)
	}
	t.Cleanup(func() {
		for _, key := range envVars {
			if saved[key] == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, saved[key])
			}
		}
	})

	// Set depth == max_depth to trigger ErrDepthExceeded
	os.Setenv("QUINE_DEPTH", "5")
	os.Setenv("QUINE_MAX_DEPTH", "5")
	// Required fields to avoid earlier validation errors
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-20250514")
	os.Setenv("QUINE_API_TYPE", "anthropic")
	os.Setenv("QUINE_API_BASE", "https://api.anthropic.com")
	os.Setenv("QUINE_API_KEY", "test-key")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected ErrDepthExceeded, got nil")
	}
	if !errors.Is(err, config.ErrDepthExceeded) {
		t.Errorf("expected ErrDepthExceeded, got: %v", err)
	}

	// Verify that main.go would produce exit code 126 for this error.
	// (The main function checks errors.Is(err, config.ErrDepthExceeded) → os.Exit(126))
	// We can't call os.Exit in a test, but we verify the error matches what
	// triggers the 126 path.
	expectedExitCode := 126
	if !errors.Is(err, config.ErrDepthExceeded) {
		t.Errorf("error should trigger exit code %d path", expectedExitCode)
	}
}

// ---------------------------------------------------------------------------
// Test 4: Context Isolation — ChildEnv verification
// ---------------------------------------------------------------------------

func TestContextIsolation(t *testing.T) {
	parentCfg := &config.Config{
		ModelID:        "claude-test",
		APIKey:         "test-key",
		APIBase:        "https://api.example.com",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          2,
		SessionID:      "parent-session-abcd",
		ParentSession:  "grandparent-session-1234",
		MaxConcurrent:  20,
		ShTimeout:      120,
		OutputTruncate: 20480,
		DataDir:        ".quine/",
		Shell:          "/bin/sh",
	}

	childEnv, err := parentCfg.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	// Parse child env into a map for easy lookup.
	envMap := make(map[string]string)
	for _, entry := range childEnv {
		key, val, _ := strings.Cut(entry, "=")
		envMap[key] = val
	}

	// QUINE_DEPTH should be parent depth + 1
	if got := envMap["QUINE_DEPTH"]; got != "3" {
		t.Errorf("QUINE_DEPTH = %q, want %q", got, "3")
	}

	// QUINE_SESSION_ID should NOT be in childEnv — each child generates its own
	// via config.Load() to prevent tape file collisions when multiple children
	// are spawned from a single sh command.
	if _, hasSessionID := envMap["QUINE_SESSION_ID"]; hasSessionID {
		t.Error("ChildEnv should NOT include QUINE_SESSION_ID (children generate their own)")
	}

	// QUINE_PARENT_SESSION should be the parent's session ID
	if got := envMap["QUINE_PARENT_SESSION"]; got != parentCfg.SessionID {
		t.Errorf("QUINE_PARENT_SESSION = %q, want %q", got, parentCfg.SessionID)
	}

	// Inherited fields should be preserved
	if got := envMap["QUINE_MODEL_ID"]; got != parentCfg.ModelID {
		t.Errorf("QUINE_MODEL_ID = %q, want %q", got, parentCfg.ModelID)
	}
	// API key is passed through QUINE_API_KEY
	if got := envMap["QUINE_API_KEY"]; got != parentCfg.APIKey {
		t.Errorf("QUINE_API_KEY = %q, want %q", got, parentCfg.APIKey)
	}
	if got := envMap["QUINE_MAX_DEPTH"]; got != "5" {
		t.Errorf("QUINE_MAX_DEPTH = %q, want %q", got, "5")
	}
	if got := envMap["QUINE_API_TYPE"]; got != parentCfg.Provider {
		t.Errorf("QUINE_API_TYPE = %q, want %q", got, parentCfg.Provider)
	}
}
