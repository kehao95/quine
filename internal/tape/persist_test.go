package tape

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriterCreatesFile(t *testing.T) {
	dir := t.TempDir()
	w, err := NewWriter(dir, "test-session")
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	defer w.Close()

	tp := NewTape("test-session", "", 0, "test-model")
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("WriteEntry: %v", err)
	}

	path := filepath.Join(dir, "test-session.jsonl")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should exist: %v", err)
	}
}

func TestWriterRoundTrip(t *testing.T) {
	dir := t.TempDir()
	w, err := NewWriter(dir, "rt-session")
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape("rt-session", "parent-1", 1, "gpt-4o")

	// Write a meta entry
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("WriteEntry meta: %v", err)
	}

	// Write a message entry
	msg := Message{Role: RoleUser, Content: "hello", Timestamp: 1700000000001}
	if err := w.WriteEntry(MessageEntry(msg)); err != nil {
		t.Fatalf("WriteEntry message: %v", err)
	}

	w.Close()

	// Read back and decode
	path := filepath.Join(dir, "rt-session.jsonl")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var entries []TapeEntry
	for scanner.Scan() {
		var entry TapeEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			t.Fatalf("Unmarshal line: %v", err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Scanner error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Type != "meta" {
		t.Errorf("entry[0].Type = %q, want %q", entries[0].Type, "meta")
	}
	if entries[1].Type != "message" {
		t.Errorf("entry[1].Type = %q, want %q", entries[1].Type, "message")
	}

	// Verify round-trip of message data
	var decoded Message
	if err := json.Unmarshal(entries[1].Data, &decoded); err != nil {
		t.Fatalf("Unmarshal message data: %v", err)
	}
	if decoded.Role != RoleUser {
		t.Errorf("decoded.Role = %q, want %q", decoded.Role, RoleUser)
	}
	if decoded.Content != "hello" {
		t.Errorf("decoded.Content = %q, want %q", decoded.Content, "hello")
	}
}

func TestWriterCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	nested := filepath.Join(base, "a", "b", "c")

	w, err := NewWriter(nested, "deep-session")
	if err != nil {
		t.Fatalf("NewWriter with nested dir: %v", err)
	}
	defer w.Close()

	if err := w.WriteEntry(TapeEntry{Type: "meta", Data: json.RawMessage(`{}`)}); err != nil {
		t.Fatalf("WriteEntry: %v", err)
	}

	path := filepath.Join(nested, "deep-session.jsonl")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should exist in nested dir: %v", err)
	}
}

func TestWriterFullLifecycle(t *testing.T) {
	dir := t.TempDir()
	w, err := NewWriter(dir, "lifecycle")
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	tp := NewTape("lifecycle", "", 0, "claude-sonnet-4-20250514")

	// 1. Meta
	if err := w.WriteEntry(tp.MetaEntry()); err != nil {
		t.Fatalf("WriteEntry meta: %v", err)
	}

	// 2. System message
	sysMsg := Message{Role: RoleSystem, Content: "system prompt", Timestamp: 1}
	if err := w.WriteEntry(MessageEntry(sysMsg)); err != nil {
		t.Fatalf("WriteEntry system: %v", err)
	}

	// 3. User message
	userMsg := Message{Role: RoleUser, Content: "do something", Timestamp: 2}
	if err := w.WriteEntry(MessageEntry(userMsg)); err != nil {
		t.Fatalf("WriteEntry user: %v", err)
	}

	// 4. Assistant message with tool call
	assistantMsg := Message{
		Role:    RoleAssistant,
		Content: "running sh",
		ToolCalls: []ToolCall{
			{ID: "tc-1", Name: "sh", Arguments: map[string]any{"command": "ls"}},
		},
		Timestamp: 3,
	}
	if err := w.WriteEntry(MessageEntry(assistantMsg)); err != nil {
		t.Fatalf("WriteEntry assistant: %v", err)
	}

	// 5. Tool result
	tr := ToolResult{ToolID: "tc-1", Content: "file.go", IsError: false}
	if err := w.WriteEntry(ToolResultEntry(tr)); err != nil {
		t.Fatalf("WriteEntry tool_result: %v", err)
	}

	// 6. Outcome
	tp.AddUsage(100, 50)
	tp.SetOutcome(SessionOutcome{
		ExitCode:        0,
		DurationMs:      500,
		TerminationMode: TermExit,
	})
	if err := w.WriteEntry(tp.OutcomeEntry()); err != nil {
		t.Fatalf("WriteEntry outcome: %v", err)
	}

	w.Close()

	// Read back and verify line count and types
	path := filepath.Join(dir, "lifecycle.jsonl")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer f.Close()

	expectedTypes := []string{"meta", "message", "message", "message", "tool_result", "outcome"}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		if lineNum >= len(expectedTypes) {
			t.Fatalf("more lines than expected (%d)", len(expectedTypes))
		}
		var entry TapeEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			t.Fatalf("line %d: Unmarshal: %v", lineNum, err)
		}
		if entry.Type != expectedTypes[lineNum] {
			t.Errorf("line %d: type = %q, want %q", lineNum, entry.Type, expectedTypes[lineNum])
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Scanner error: %v", err)
	}
	if lineNum != len(expectedTypes) {
		t.Errorf("got %d lines, want %d", lineNum, len(expectedTypes))
	}
}
