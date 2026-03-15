package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseForkArgs_SingleChild(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "Do something useful", "workspace": "."},
		},
	}
	req, err := ParseForkArgs(args)
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
			map[string]any{"intent": "Task A", "workspace": "."},
			map[string]any{"intent": "Task B", "workspace": "sub"},
			map[string]any{"intent": "Task C", "workspace": "sub/deeper"},
		},
	}
	req, err := ParseForkArgs(args)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if len(req.Children) != 3 {
		t.Fatalf("got %d children, want 3", len(req.Children))
	}
	if req.Children[0].Intent != "Task A" || req.Children[1].Workspace != "sub" || req.Children[2].Workspace != "sub/deeper" {
		t.Errorf("Children = %v", req.Children)
	}
}

func TestParseForkArgs_ModeForget(t *testing.T) {
	args := map[string]any{
		"children": []interface{}{
			map[string]any{"intent": "Calculate something", "workspace": "."},
		},
		"mode": "forget",
	}
	req, err := ParseForkArgs(args)
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
			map[string]any{"intent": "Try approach A", "workspace": "."},
			map[string]any{"intent": "Try approach B", "workspace": "."},
		},
		"mode": "race",
	}
	req, err := ParseForkArgs(args)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if req.Mode != ForkModeRace {
		t.Errorf("Mode = %q, want %q", req.Mode, ForkModeRace)
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
		{"MissingIntent", map[string]any{"children": []interface{}{map[string]any{"workspace": "."}}}, "intent"},
		{"MissingWorkspace", map[string]any{"children": []interface{}{map[string]any{"intent": "ok"}}}, "workspace"},
		{"EmptyIntent", map[string]any{"children": []interface{}{map[string]any{"intent": "", "workspace": "."}}}, "intent"},
		{"EmptyWorkspace", map[string]any{"children": []interface{}{map[string]any{"intent": "ok", "workspace": ""}}}, "workspace"},
		{"WrongModeType", map[string]any{"children": []interface{}{map[string]any{"intent": "task", "workspace": "."}}, "mode": 42}, "string"},
		{"InvalidModeValue", map[string]any{"children": []interface{}{map[string]any{"intent": "task", "workspace": "."}}, "mode": "yolo"}, "must be one of"},
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

// Backward compatibility: old wait/race booleans should map to mode.
func TestParseForkArgs_BackwardCompat_WaitFalse(t *testing.T) {
	args := map[string]any{
		"argv": []interface{}{"task"},
		"wait": false,
	}
	req, err := ParseForkArgs(args)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if req.Mode != ForkModeForget {
		t.Errorf("Mode = %q, want %q (wait=false → forget)", req.Mode, ForkModeForget)
	}
}

func TestParseForkArgs_BackwardCompat_RaceTrue(t *testing.T) {
	args := map[string]any{
		"argv": []interface{}{"A", "B"},
		"race": true,
	}
	req, err := ParseForkArgs(args)
	if err != nil {
		t.Fatalf("ParseForkArgs failed: %v", err)
	}
	if req.Mode != ForkModeRace {
		t.Errorf("Mode = %q, want %q (race=true → race)", req.Mode, ForkModeRace)
	}
}

func TestFilterProcessIdentity(t *testing.T) {
	env := []string{
		"PATH=/usr/bin",
		"QUINE_SESSION_ID=old-session",
		"QUINE_TAPE_ID=tape-session",
		"QUINE_DEPTH=1",
		"HOME=/home/user",
	}
	filtered := filterProcessIdentity(env)

	// Should not contain QUINE_SESSION_ID or QUINE_TAPE_ID
	for _, e := range filtered {
		if strings.HasPrefix(e, "QUINE_SESSION_ID=") {
			t.Errorf("filtered env should not contain QUINE_SESSION_ID: %v", filtered)
		}
		if strings.HasPrefix(e, "QUINE_TAPE_ID=") {
			t.Errorf("filtered env should not contain QUINE_TAPE_ID: %v", filtered)
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
		Children: []ForkChild{{Intent: "test intent", Workspace: "."}},
		Mode:     ForkModeWait,
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
		Children: []ForkChild{{Intent: "0.1", Workspace: "."}}, // sleep argument
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
		if !strings.Contains(result.Content, "[FORK OK]") {
			t.Errorf("expected [FORK OK] in result, got: %s", result.Content)
		}
	}
}

func TestForkExecutor_Execute_GatherAll_MissingBinary(t *testing.T) {
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

	os.WriteFile(f.TapePath, []byte{}, 0644)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "task A", Workspace: "."}, {Intent: "task B", Workspace: "."}},
		Mode:     ForkModeWait,
	}

	result := f.Execute("tool-1", req)
	if !result.IsError {
		t.Errorf("expected error for missing binary")
	}
	if !strings.Contains(result.Content, "FORK ERROR") {
		t.Errorf("expected FORK ERROR in result, got: %s", result.Content)
	}
}

func TestForkExecutor_Execute_Race_MissingBinary(t *testing.T) {
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

	os.WriteFile(f.TapePath, []byte{}, 0644)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "approach A", Workspace: "."}, {Intent: "approach B", Workspace: "."}},
		Mode:     ForkModeRace,
	}

	result := f.Execute("tool-1", req)
	if !result.IsError {
		t.Errorf("expected error for missing binary")
	}
	if !strings.Contains(result.Content, "FORK ERROR") {
		t.Errorf("expected FORK ERROR in result, got: %s", result.Content)
	}
}

func TestForkExecutor_Execute_MultipleAsync(t *testing.T) {
	tmpDir := t.TempDir()

	f := &ForkExecutor{
		QuinePath:      "/bin/sleep",
		DataDir:        tmpDir,
		SessionID:      "test-session",
		TapePath:       filepath.Join(tmpDir, "test-session.jsonl"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}

	os.WriteFile(f.TapePath, []byte{}, 0644)

	req := ForkRequest{
		Children: []ForkChild{{Intent: "0.1", Workspace: "."}, {Intent: "0.1", Workspace: "."}, {Intent: "0.1", Workspace: "."}},
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
		if !strings.Contains(result.Content, "3 children spawned") {
			t.Errorf("expected '3 children spawned' in result, got: %s", result.Content)
		}
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
		TapePath:       filepath.Join(tmpDir, "test-session.jsonl"),
		DefaultTimeout: 5 * time.Second,
		MaxOutput:      10000,
		Env:            []string{},
	}

	if err := os.WriteFile(f.TapePath, []byte{}, 0o644); err != nil {
		t.Fatalf("write tape: %v", err)
	}

	result := f.Execute("tool-1", ForkRequest{
		Children: []ForkChild{{Intent: "ignored", Workspace: "."}},
		Mode:     ForkModeWait,
	})

	if result.IsError {
		t.Fatalf("expected child helper to succeed, got error:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "NO_FD3") {
		t.Fatalf("expected child output to report missing fd 3, got:\n%s", result.Content)
	}
	if strings.Contains(result.Content, "HAS_FD3") {
		t.Fatalf("child unexpectedly inherited live fd 3:\n%s", result.Content)
	}
}
