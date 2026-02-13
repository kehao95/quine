package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseForkArgs_ValidIntent(t *testing.T) {
	args := map[string]any{
		"intent": "Do something useful",
	}
	req, err := ParseForkArgs(args)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if req.Intent != "Do something useful" {
		t.Errorf("Intent = %q, want %q", req.Intent, "Do something useful")
	}
	if req.Wait {
		t.Errorf("Wait = true, want false (default)")
	}
}

func TestParseForkArgs_WithWaitTrue(t *testing.T) {
	args := map[string]any{
		"intent": "Calculate something",
		"wait":   true,
	}
	req, err := ParseForkArgs(args)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if !req.Wait {
		t.Errorf("Wait = false, want true")
	}
}

func TestParseForkArgs_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		errWord string
	}{
		{"MissingIntent", map[string]any{"wait": true}, "intent"},
		{"EmptyIntent", map[string]any{"intent": ""}, "empty"},
		{"WrongIntentType", map[string]any{"intent": 123}, "string"},
		{"WrongWaitType", map[string]any{"intent": "Do something", "wait": "yes"}, "boolean"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseForkArgs(tt.args)
			if err == nil {
				t.Fatalf("expected error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.errWord) {
				t.Errorf("error should mention %q: %v", tt.errWord, err)
			}
		})
	}
}

func TestFilterSessionID(t *testing.T) {
	env := []string{
		"PATH=/usr/bin",
		"QUINE_SESSION_ID=old-session",
		"QUINE_DEPTH=1",
		"HOME=/home/user",
	}
	filtered := filterSessionID(env)

	// Should not contain QUINE_SESSION_ID
	for _, e := range filtered {
		if strings.HasPrefix(e, "QUINE_SESSION_ID=") {
			t.Errorf("filtered env should not contain QUINE_SESSION_ID: %v", filtered)
		}
	}

	// Should contain other entries
	if len(filtered) != 3 {
		t.Errorf("expected 3 entries, got %d: %v", len(filtered), filtered)
	}
}

func TestForkExecutor_CopyTapeForChild(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake tape file
	tapePath := filepath.Join(tmpDir, "test-session.jsonl")
	tapeContent := `{"type":"meta","data":"test"}
{"type":"message","data":"hello"}
`
	if err := os.WriteFile(tapePath, []byte(tapeContent), 0644); err != nil {
		t.Fatalf("failed to write test tape: %v", err)
	}

	f := &ForkExecutor{
		DataDir:   tmpDir,
		SessionID: "test-session",
		TapePath:  tapePath,
	}

	childTapePath, err := f.copyTapeForChild()
	if err != nil {
		t.Fatalf("copyTapeForChild failed: %v", err)
	}
	defer os.Remove(childTapePath)

	// Verify the child tape contains the same content
	childContent, err := os.ReadFile(childTapePath)
	if err != nil {
		t.Fatalf("failed to read child tape: %v", err)
	}
	if string(childContent) != tapeContent {
		t.Errorf("child tape content mismatch.\ngot:\n%s\nwant:\n%s", string(childContent), tapeContent)
	}
}

func TestForkExecutor_CopyTapeForChild_NoTape(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		DataDir:   tmpDir,
		SessionID: "nonexistent-session",
		TapePath:  filepath.Join(tmpDir, "nonexistent.jsonl"),
	}

	childTapePath, err := f.copyTapeForChild()
	if err != nil {
		t.Fatalf("copyTapeForChild should not fail for nonexistent tape: %v", err)
	}
	if childTapePath != "" {
		t.Errorf("expected empty path for nonexistent tape, got %q", childTapePath)
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
		TapePath:       filepath.Join(tmpDir, "test-session.jsonl"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}

	// Create empty tape file
	os.WriteFile(f.TapePath, []byte{}, 0644)

	req := ForkRequest{
		Intent: "test intent",
		Wait:   true,
	}

	result := f.Execute("tool-1", req)
	if !result.IsError {
		t.Errorf("expected error for missing binary")
	}
	if !strings.Contains(result.Content, "FORK ERROR") {
		t.Errorf("expected FORK ERROR in result, got: %s", result.Content)
	}
}

func TestForkExecutor_Execute_AsyncMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Use a real command that will run briefly
	f := &ForkExecutor{
		QuinePath:      "/bin/sleep", // Will fail but that's ok for async
		DataDir:        tmpDir,
		SessionID:      "test-session",
		TapePath:       filepath.Join(tmpDir, "test-session.jsonl"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}

	// Create empty tape file
	os.WriteFile(f.TapePath, []byte{}, 0644)

	req := ForkRequest{
		Intent: "0.1", // sleep argument
		Wait:   false,
	}

	start := time.Now()
	result := f.Execute("tool-1", req)
	elapsed := time.Since(start)

	// Async should return immediately (not wait for child)
	if elapsed > 2*time.Second {
		t.Errorf("async fork took too long: %v", elapsed)
	}

	// Result should indicate child was spawned
	if result.IsError {
		// It's ok if it fails to start, but shouldn't take long
		t.Logf("async fork error (expected for sleep command): %s", result.Content)
	} else {
		if !strings.Contains(result.Content, "[FORK]") {
			t.Errorf("expected [FORK] in result, got: %s", result.Content)
		}
	}
}
