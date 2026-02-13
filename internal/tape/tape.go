package tape

import (
	"encoding/json"
	"time"
)

// Role identifies the sender of a message in the conversation tape.
type Role string

const (
	RoleSystem     Role = "system"
	RoleUser       Role = "user"
	RoleAssistant  Role = "assistant"
	RoleToolResult Role = "tool_result"
)

// ToolCall represents a tool invocation requested by the assistant.
type ToolCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// Message is a single turn in the conversation tape.
type Message struct {
	Role             Role       `json:"role"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	ToolID           string     `json:"tool_id,omitempty"`
	Timestamp        int64      `json:"timestamp"`
}

// TerminationMode describes how a session ended.
type TerminationMode string

const (
	TermExit              TerminationMode = "exit"
	TermContextExhaustion TerminationMode = "context_exhaustion"
	TermTurnExhaustion    TerminationMode = "turn_exhaustion"
	TermTimeout           TerminationMode = "timeout"
	TermSignal            TerminationMode = "signal"
	TermExec              TerminationMode = "exec" // Process replaced via exec syscall
)

// SessionOutcome captures the final result of a session.
type SessionOutcome struct {
	ExitCode        int             `json:"exit_code"`
	Stderr          string          `json:"stderr"`
	DurationMs      int64           `json:"duration_ms"`
	TokensIn        int             `json:"tokens_in"`
	TokensOut       int             `json:"tokens_out"`
	TurnCount       int             `json:"turn_count"`
	TerminationMode TerminationMode `json:"termination_mode"`
}

// ToolResult holds the output of a tool execution.
type ToolResult struct {
	ToolID  string `json:"tool_id"`
	Content string `json:"content"`
	IsError bool   `json:"is_error"`
}

// Tape is the append-only conversation log for a single session.
type Tape struct {
	SessionID       string          `json:"session_id"`
	ParentSessionID string          `json:"parent_session_id"`
	Depth           int             `json:"depth"`
	ModelID         string          `json:"model_id"`
	CreatedAt       int64           `json:"created_at"`
	messages        []Message       // unexported, append-only
	Outcome         *SessionOutcome `json:"outcome,omitempty"`

	// Token tracking (excluded from JSON — persisted via SessionOutcome)
	TokensIn  int `json:"-"`
	TokensOut int `json:"-"`
	TurnCount int `json:"-"`
}

// NewTape creates a fresh Tape with CreatedAt set to the current time.
func NewTape(sessionID, parentSessionID string, depth int, modelID string) *Tape {
	return &Tape{
		SessionID:       sessionID,
		ParentSessionID: parentSessionID,
		Depth:           depth,
		ModelID:         modelID,
		CreatedAt:       time.Now().UnixMilli(),
	}
}

// Append adds a message to the tape. If msg.Timestamp is zero it is set
// to the current time in milliseconds since epoch.
func (t *Tape) Append(msg Message) {
	if msg.Timestamp == 0 {
		msg.Timestamp = time.Now().UnixMilli()
	}
	t.messages = append(t.messages, msg)
}

// Messages returns a shallow copy of the internal messages slice so that
// callers cannot mutate the tape's state.
func (t *Tape) Messages() []Message {
	out := make([]Message, len(t.messages))
	copy(out, t.messages)
	return out
}

// Len returns the number of messages on the tape.
func (t *Tape) Len() int {
	return len(t.messages)
}

// LastMessage returns a pointer to the last message on the tape,
// allowing the caller to mutate it in place. Returns nil if the tape is empty.
func (t *Tape) LastMessage() *Message {
	if len(t.messages) == 0 {
		return nil
	}
	return &t.messages[len(t.messages)-1]
}

// SetOutcome records the final session outcome. Running token and turn
// totals are copied into the outcome struct.
func (t *Tape) SetOutcome(outcome SessionOutcome) {
	outcome.TokensIn = t.TokensIn
	outcome.TokensOut = t.TokensOut
	outcome.TurnCount = t.TurnCount
	t.Outcome = &outcome
}

// AddUsage accumulates token counts. Does NOT increment turn counter
// (turns are only consumed by sh tool calls, tracked separately).
func (t *Tape) AddUsage(tokensIn, tokensOut int) {
	t.TokensIn += tokensIn
	t.TokensOut += tokensOut
}

// IncrementTurn increments the turn counter. Called only when sh tool is used.
func (t *Tape) IncrementTurn() {
	t.TurnCount++
}

// ---------------------------------------------------------------------------
// JSONL entry types (§10.1)
// ---------------------------------------------------------------------------

// TapeEntry is a single line in the JSONL tape file. The Type field
// discriminates the payload stored in Data.
type TapeEntry struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// MetaEntry returns a TapeEntry of type "meta" containing the tape header.
func (t *Tape) MetaEntry() TapeEntry {
	type meta struct {
		SessionID       string `json:"session_id"`
		ParentSessionID string `json:"parent_session_id"`
		Depth           int    `json:"depth"`
		ModelID         string `json:"model_id"`
		CreatedAt       int64  `json:"created_at"`
	}
	data, _ := json.Marshal(meta{
		SessionID:       t.SessionID,
		ParentSessionID: t.ParentSessionID,
		Depth:           t.Depth,
		ModelID:         t.ModelID,
		CreatedAt:       t.CreatedAt,
	})
	return TapeEntry{Type: "meta", Data: data}
}

// MessageEntry returns a TapeEntry of type "message" wrapping msg.
func MessageEntry(msg Message) TapeEntry {
	data, _ := json.Marshal(msg)
	return TapeEntry{Type: "message", Data: data}
}

// ToolResultEntry returns a TapeEntry of type "tool_result" wrapping tr.
func ToolResultEntry(tr ToolResult) TapeEntry {
	data, _ := json.Marshal(tr)
	return TapeEntry{Type: "tool_result", Data: data}
}

// OutcomeEntry returns a TapeEntry of type "outcome" wrapping the session outcome.
// It returns a zero-value TapeEntry if no outcome has been set.
func (t *Tape) OutcomeEntry() TapeEntry {
	if t.Outcome == nil {
		return TapeEntry{Type: "outcome"}
	}
	data, _ := json.Marshal(t.Outcome)
	return TapeEntry{Type: "outcome", Data: data}
}
