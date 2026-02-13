package tape

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadTapeFile(t *testing.T) {
	dir := t.TempDir()
	sessionID := "read-test"

	// Write a complete tape using Writer.
	w, err := NewWriter(dir, sessionID)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape(sessionID, "parent-read", 2, "claude-sonnet-4-20250514")

	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("write meta: %v", err)
	}
	if err := w.WriteEntry(MessageEntry(Message{Role: RoleSystem, Content: "prompt", Timestamp: 1})); err != nil {
		t.Fatalf("write system: %v", err)
	}
	if err := w.WriteEntry(MessageEntry(Message{Role: RoleUser, Content: "task", Timestamp: 2})); err != nil {
		t.Fatalf("write user: %v", err)
	}
	if err := w.WriteEntry(MessageEntry(Message{
		Role:      RoleAssistant,
		Content:   "thinking",
		ToolCalls: []ToolCall{{ID: "tc-1", Name: "sh", Arguments: map[string]any{"command": "ls"}}},
		Timestamp: 3,
	})); err != nil {
		t.Fatalf("write assistant: %v", err)
	}
	if err := w.WriteEntry(ToolResultEntry(ToolResult{ToolID: "tc-1", Content: "file.go", IsError: false})); err != nil {
		t.Fatalf("write tool_result: %v", err)
	}

	tp.AddUsage(500, 200)
	tp.SetOutcome(SessionOutcome{
		ExitCode:        0,
		DurationMs:      1000,
		TerminationMode: TermExit,
	})
	if err := w.WriteEntry(tp.OutcomeEntry()); err != nil {
		t.Fatalf("write outcome: %v", err)
	}
	w.Close()

	// Read it back.
	path := filepath.Join(dir, sessionID+".jsonl")
	summary, err := ReadTapeFile(path)
	if err != nil {
		t.Fatalf("ReadTapeFile: %v", err)
	}

	if summary.SessionID != sessionID {
		t.Errorf("SessionID = %q, want %q", summary.SessionID, sessionID)
	}
	if summary.ParentSessionID != "parent-read" {
		t.Errorf("ParentSessionID = %q, want %q", summary.ParentSessionID, "parent-read")
	}
	if summary.Depth != 2 {
		t.Errorf("Depth = %d, want 2", summary.Depth)
	}
	if summary.ModelID != "claude-sonnet-4-20250514" {
		t.Errorf("ModelID = %q, want %q", summary.ModelID, "claude-sonnet-4-20250514")
	}
	if len(summary.Entries) != 6 {
		t.Fatalf("Entries count = %d, want 6", len(summary.Entries))
	}
	if summary.Outcome == nil {
		t.Fatal("Outcome should not be nil")
	}
	if summary.Outcome.ExitCode != 0 {
		t.Errorf("Outcome.ExitCode = %d, want 0", summary.Outcome.ExitCode)
	}
	if summary.Outcome.TokensIn != 500 {
		t.Errorf("Outcome.TokensIn = %d, want 500", summary.Outcome.TokensIn)
	}
}

func TestReadTapeFile_NoOutcome(t *testing.T) {
	dir := t.TempDir()
	sessionID := "running-test"

	w, err := NewWriter(dir, sessionID)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape(sessionID, "", 0, "gpt-4o")
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("write meta: %v", err)
	}
	if err := w.WriteEntry(MessageEntry(Message{Role: RoleUser, Content: "working...", Timestamp: 1})); err != nil {
		t.Fatalf("write user: %v", err)
	}
	w.Close()

	path := filepath.Join(dir, sessionID+".jsonl")
	summary, err := ReadTapeFile(path)
	if err != nil {
		t.Fatalf("ReadTapeFile: %v", err)
	}

	if summary.Outcome != nil {
		t.Error("Outcome should be nil for running session")
	}
	if len(summary.Entries) != 2 {
		t.Errorf("Entries count = %d, want 2", len(summary.Entries))
	}
}

func TestReadTapeFile_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := ReadTapeFile(path)
	if err == nil {
		t.Fatal("expected error for empty tape file")
	}
}

func TestTailLastEntry(t *testing.T) {
	dir := t.TempDir()
	sessionID := "tail-test"

	w, err := NewWriter(dir, sessionID)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape(sessionID, "", 0, "model")
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("write meta: %v", err)
	}
	if err := w.WriteEntry(MessageEntry(Message{Role: RoleUser, Content: "hello", Timestamp: 1})); err != nil {
		t.Fatalf("write user: %v", err)
	}
	if err := w.WriteEntry(MessageEntry(Message{
		Role:      RoleAssistant,
		Content:   "running sh",
		ToolCalls: []ToolCall{{ID: "tc-1", Name: "sh", Arguments: map[string]any{"command": "echo hi"}}},
		Timestamp: 2,
	})); err != nil {
		t.Fatalf("write assistant: %v", err)
	}
	w.Close()

	path := filepath.Join(dir, sessionID+".jsonl")
	entry, err := TailLastEntry(path)
	if err != nil {
		t.Fatalf("TailLastEntry: %v", err)
	}
	if entry.Type != "message" {
		t.Errorf("Type = %q, want %q", entry.Type, "message")
	}

	// Verify it's the assistant message (the last one).
	var msg Message
	if err := json.Unmarshal(entry.Data, &msg); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if msg.Role != RoleAssistant {
		t.Errorf("Role = %q, want %q", msg.Role, RoleAssistant)
	}
	if msg.Content != "running sh" {
		t.Errorf("Content = %q, want %q", msg.Content, "running sh")
	}
}

func TestTailLastEntry_SingleLine(t *testing.T) {
	dir := t.TempDir()
	sessionID := "tail-single"

	w, err := NewWriter(dir, sessionID)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape(sessionID, "", 0, "model")
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("write meta: %v", err)
	}
	w.Close()

	path := filepath.Join(dir, sessionID+".jsonl")
	entry, err := TailLastEntry(path)
	if err != nil {
		t.Fatalf("TailLastEntry: %v", err)
	}
	if entry.Type != "meta" {
		t.Errorf("Type = %q, want %q", entry.Type, "meta")
	}
}

func TestReadTapeFile_DuplicateMeta(t *testing.T) {
	// Simulate the bug where multiple ./quine children write to the same tape
	// file because they share a session ID. The file has duplicate meta entries
	// interleaved with messages. The reader should skip duplicate metas and
	// use only the first one.
	dir := t.TempDir()
	sessionID := "dupe-meta-test"

	w, err := NewWriter(dir, sessionID)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape(sessionID, "parent-1", 1, "claude-sonnet-4-20250514")

	// Write meta 3 times (simulating 3 children with same session ID)
	for i := 0; i < 3; i++ {
		if err := w.WriteEntry(tp.MetaEntry()); err != nil {
			t.Fatalf("write meta %d: %v", i, err)
		}
	}

	// Write a message
	if err := w.WriteEntry(MessageEntry(Message{Role: RoleUser, Content: "task", Timestamp: 1})); err != nil {
		t.Fatalf("write user: %v", err)
	}

	// Write another duplicate meta (interleaved, as seen in real tapes)
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("write meta interleaved: %v", err)
	}

	// Write more messages
	if err := w.WriteEntry(MessageEntry(Message{Role: RoleAssistant, Content: "response", Timestamp: 2})); err != nil {
		t.Fatalf("write assistant: %v", err)
	}

	w.Close()

	// Read it back — should succeed despite duplicate metas
	path := filepath.Join(dir, sessionID+".jsonl")
	summary, err := ReadTapeFile(path)
	if err != nil {
		t.Fatalf("ReadTapeFile should handle duplicate metas, got: %v", err)
	}

	if summary.SessionID != sessionID {
		t.Errorf("SessionID = %q, want %q", summary.SessionID, sessionID)
	}
	if summary.ParentSessionID != "parent-1" {
		t.Errorf("ParentSessionID = %q, want %q", summary.ParentSessionID, "parent-1")
	}
	if summary.Depth != 1 {
		t.Errorf("Depth = %d, want 1", summary.Depth)
	}

	// Entries should include everything (metas are skipped from re-parsing but still in entries list)
	// We wrote 4 meta + 2 message = 6 total entries in the file
	// But duplicate metas are skipped via continue, so they're NOT added to entries
	// Only 1 meta + 2 messages = 3 entries
	if len(summary.Entries) != 3 {
		t.Errorf("Entries count = %d, want 3 (1 meta + 2 messages)", len(summary.Entries))
	}
}
