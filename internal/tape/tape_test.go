package tape

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewTape(t *testing.T) {
	before := time.Now().UnixMilli()
	tp := NewTape("sess-1", "parent-1", 2, "claude-sonnet-4-20250514")
	after := time.Now().UnixMilli()

	if tp.SessionID != "sess-1" {
		t.Errorf("SessionID = %q, want %q", tp.SessionID, "sess-1")
	}
	if tp.ParentSessionID != "parent-1" {
		t.Errorf("ParentSessionID = %q, want %q", tp.ParentSessionID, "parent-1")
	}
	if tp.Depth != 2 {
		t.Errorf("Depth = %d, want 2", tp.Depth)
	}
	if tp.ModelID != "claude-sonnet-4-20250514" {
		t.Errorf("ModelID = %q, want %q", tp.ModelID, "claude-sonnet-4-20250514")
	}
	if tp.CreatedAt < before || tp.CreatedAt > after {
		t.Errorf("CreatedAt = %d, want between %d and %d", tp.CreatedAt, before, after)
	}
	if tp.Len() != 0 {
		t.Errorf("Len() = %d, want 0", tp.Len())
	}
	if tp.Outcome != nil {
		t.Error("Outcome should be nil for a new tape")
	}
}

func TestAppend(t *testing.T) {
	tp := NewTape("s", "", 0, "m")

	// Message with explicit timestamp should keep it.
	m1 := Message{Role: RoleSystem, Content: "hello", Timestamp: 42}
	tp.Append(m1)
	if tp.Len() != 1 {
		t.Fatalf("Len() = %d after first Append, want 1", tp.Len())
	}
	msgs := tp.Messages()
	if msgs[0].Timestamp != 42 {
		t.Errorf("explicit timestamp = %d, want 42", msgs[0].Timestamp)
	}

	// Message with zero timestamp should be set automatically.
	before := time.Now().UnixMilli()
	m2 := Message{Role: RoleUser, Content: "world"}
	tp.Append(m2)
	after := time.Now().UnixMilli()

	if tp.Len() != 2 {
		t.Fatalf("Len() = %d after second Append, want 2", tp.Len())
	}
	msgs = tp.Messages()
	if msgs[1].Timestamp < before || msgs[1].Timestamp > after {
		t.Errorf("auto timestamp = %d, want between %d and %d", msgs[1].Timestamp, before, after)
	}
}

func TestMessagesCopy(t *testing.T) {
	tp := NewTape("s", "", 0, "m")
	tp.Append(Message{Role: RoleUser, Content: "original"})

	msgs := tp.Messages()
	// Mutate the returned slice.
	msgs[0].Content = "mutated"

	// Internal state must be unchanged.
	internal := tp.Messages()
	if internal[0].Content != "original" {
		t.Errorf("internal Content = %q, want %q — Messages() did not return a copy",
			internal[0].Content, "original")
	}

	// Appending to the returned slice must not affect the tape.
	msgs = append(msgs, Message{Role: RoleAssistant, Content: "extra"})
	_ = msgs
	if tp.Len() != 1 {
		t.Errorf("Len() = %d after mutating returned slice, want 1", tp.Len())
	}
}

func TestAddUsage(t *testing.T) {
	tp := NewTape("s", "", 0, "m")

	tp.AddUsage(100, 50)
	// AddUsage only tracks tokens, NOT turns (turns are tracked separately via IncrementTurn)
	if tp.TokensIn != 100 || tp.TokensOut != 50 {
		t.Errorf("after first AddUsage: in=%d out=%d, want 100/50",
			tp.TokensIn, tp.TokensOut)
	}

	tp.AddUsage(200, 75)
	if tp.TokensIn != 300 || tp.TokensOut != 125 {
		t.Errorf("after second AddUsage: in=%d out=%d, want 300/125",
			tp.TokensIn, tp.TokensOut)
	}
}

func TestIncrementTurn(t *testing.T) {
	tp := NewTape("s", "", 0, "m")

	if tp.TurnCount != 0 {
		t.Errorf("initial TurnCount = %d, want 0", tp.TurnCount)
	}

	tp.IncrementTurn()
	if tp.TurnCount != 1 {
		t.Errorf("after first IncrementTurn: TurnCount = %d, want 1", tp.TurnCount)
	}

	tp.IncrementTurn()
	if tp.TurnCount != 2 {
		t.Errorf("after second IncrementTurn: TurnCount = %d, want 2", tp.TurnCount)
	}
}

func TestSetOutcome(t *testing.T) {
	tp := NewTape("s", "", 0, "m")
	tp.AddUsage(500, 200)
	tp.AddUsage(300, 100)
	tp.IncrementTurn() // Turns are now tracked separately
	tp.IncrementTurn()

	outcome := SessionOutcome{
		ExitCode:        0,
		Stderr:          "",
		DurationMs:      1234,
		TerminationMode: TermExit,
	}
	tp.SetOutcome(outcome)

	if tp.Outcome == nil {
		t.Fatal("Outcome should not be nil after SetOutcome")
	}
	if tp.Outcome.TokensIn != 800 {
		t.Errorf("Outcome.TokensIn = %d, want 800", tp.Outcome.TokensIn)
	}
	if tp.Outcome.TokensOut != 300 {
		t.Errorf("Outcome.TokensOut = %d, want 300", tp.Outcome.TokensOut)
	}
	if tp.Outcome.TurnCount != 2 {
		t.Errorf("Outcome.TurnCount = %d, want 2", tp.Outcome.TurnCount)
	}
	if tp.Outcome.ExitCode != 0 {
		t.Errorf("Outcome.ExitCode = %d, want 0", tp.Outcome.ExitCode)
	}
	if tp.Outcome.TerminationMode != TermExit {
		t.Errorf("Outcome.TerminationMode = %q, want %q", tp.Outcome.TerminationMode, TermExit)
	}
}

// ---------------------------------------------------------------------------
// JSON round-trip tests for all entry types
// ---------------------------------------------------------------------------

func TestMetaEntryJSON(t *testing.T) {
	tp := NewTape("sess-rt", "parent-rt", 3, "gpt-4o")

	entry := tp.MetaEntry()
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Marshal MetaEntry: %v", err)
	}

	var decoded TapeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal MetaEntry: %v", err)
	}
	if decoded.Type != "meta" {
		t.Errorf("Type = %q, want %q", decoded.Type, "meta")
	}

	// Decode the inner data.
	var meta struct {
		SessionID       string `json:"session_id"`
		ParentSessionID string `json:"parent_session_id"`
		Depth           int    `json:"depth"`
		ModelID         string `json:"model_id"`
		CreatedAt       int64  `json:"created_at"`
	}
	if err := json.Unmarshal(decoded.Data, &meta); err != nil {
		t.Fatalf("Unmarshal meta data: %v", err)
	}
	if meta.SessionID != "sess-rt" {
		t.Errorf("meta.SessionID = %q, want %q", meta.SessionID, "sess-rt")
	}
	if meta.Depth != 3 {
		t.Errorf("meta.Depth = %d, want 3", meta.Depth)
	}
}

func TestMessageEntryJSON(t *testing.T) {
	msg := Message{
		Role:    RoleAssistant,
		Content: "thinking...",
		ToolCalls: []ToolCall{
			{ID: "tc-1", Name: "sh", Arguments: map[string]any{"command": "ls"}},
		},
		Timestamp: 1700000000000,
	}

	entry := MessageEntry(msg)
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Marshal MessageEntry: %v", err)
	}

	var decoded TapeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal MessageEntry: %v", err)
	}
	if decoded.Type != "message" {
		t.Errorf("Type = %q, want %q", decoded.Type, "message")
	}

	var m Message
	if err := json.Unmarshal(decoded.Data, &m); err != nil {
		t.Fatalf("Unmarshal message data: %v", err)
	}
	if m.Role != RoleAssistant {
		t.Errorf("Role = %q, want %q", m.Role, RoleAssistant)
	}
	if m.Content != "thinking..." {
		t.Errorf("Content = %q, want %q", m.Content, "thinking...")
	}
	if len(m.ToolCalls) != 1 || m.ToolCalls[0].Name != "sh" {
		t.Errorf("ToolCalls mismatch: %+v", m.ToolCalls)
	}
}

func TestToolResultEntryJSON(t *testing.T) {
	tr := ToolResult{
		ToolID:  "tc-1",
		Content: "file1.go\nfile2.go\n",
		IsError: false,
	}

	entry := ToolResultEntry(tr)
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Marshal ToolResultEntry: %v", err)
	}

	var decoded TapeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal ToolResultEntry: %v", err)
	}
	if decoded.Type != "tool_result" {
		t.Errorf("Type = %q, want %q", decoded.Type, "tool_result")
	}

	var result ToolResult
	if err := json.Unmarshal(decoded.Data, &result); err != nil {
		t.Fatalf("Unmarshal tool_result data: %v", err)
	}
	if result.ToolID != "tc-1" {
		t.Errorf("ToolID = %q, want %q", result.ToolID, "tc-1")
	}
	if result.Content != "file1.go\nfile2.go\n" {
		t.Errorf("Content = %q, want %q", result.Content, "file1.go\nfile2.go\n")
	}
	if result.IsError {
		t.Error("IsError should be false")
	}
}

func TestOutcomeEntryJSON(t *testing.T) {
	tp := NewTape("s", "", 0, "m")
	tp.AddUsage(1000, 500)
	tp.IncrementTurn() // Turns tracked separately
	tp.SetOutcome(SessionOutcome{
		ExitCode:        1,
		Stderr:          "err",
		DurationMs:      5000,
		TerminationMode: TermTimeout,
	})

	entry := tp.OutcomeEntry()
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Marshal OutcomeEntry: %v", err)
	}

	var decoded TapeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal OutcomeEntry: %v", err)
	}
	if decoded.Type != "outcome" {
		t.Errorf("Type = %q, want %q", decoded.Type, "outcome")
	}

	var outcome SessionOutcome
	if err := json.Unmarshal(decoded.Data, &outcome); err != nil {
		t.Fatalf("Unmarshal outcome data: %v", err)
	}
	if outcome.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", outcome.ExitCode)
	}
	if outcome.TokensIn != 1000 {
		t.Errorf("TokensIn = %d, want 1000", outcome.TokensIn)
	}
	if outcome.TokensOut != 500 {
		t.Errorf("TokensOut = %d, want 500", outcome.TokensOut)
	}
	if outcome.TurnCount != 1 {
		t.Errorf("TurnCount = %d, want 1", outcome.TurnCount)
	}
	if outcome.TerminationMode != TermTimeout {
		t.Errorf("TerminationMode = %q, want %q", outcome.TerminationMode, TermTimeout)
	}
}

func TestToolResultEntryJSON_Error(t *testing.T) {
	tr := ToolResult{
		ToolID:  "tc-err",
		Content: "command not found",
		IsError: true,
	}

	entry := ToolResultEntry(tr)
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded TapeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	var result ToolResult
	if err := json.Unmarshal(decoded.Data, &result); err != nil {
		t.Fatalf("Unmarshal data: %v", err)
	}
	if !result.IsError {
		t.Error("IsError should be true")
	}
}

func TestMessageEntryJSON_ToolID(t *testing.T) {
	msg := Message{
		Role:      RoleToolResult,
		Content:   "output here",
		ToolID:    "tc-42",
		Timestamp: 1700000000000,
	}

	entry := MessageEntry(msg)
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded TapeEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	var m Message
	if err := json.Unmarshal(decoded.Data, &m); err != nil {
		t.Fatalf("Unmarshal data: %v", err)
	}
	if m.ToolID != "tc-42" {
		t.Errorf("ToolID = %q, want %q", m.ToolID, "tc-42")
	}
	if m.ToolCalls != nil {
		t.Errorf("ToolCalls should be nil (omitempty), got %+v", m.ToolCalls)
	}
}

func TestMessageOmitemptyJSON(t *testing.T) {
	// A plain user message should not include tool_calls or tool_id in JSON.
	msg := Message{
		Role:      RoleUser,
		Content:   "hi",
		Timestamp: 1,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map: %v", err)
	}
	if _, ok := raw["tool_calls"]; ok {
		t.Error("tool_calls should be omitted for empty slice")
	}
	if _, ok := raw["tool_id"]; ok {
		t.Error("tool_id should be omitted for empty string")
	}
}
