package tape

import (
	"bytes"
	"encoding/json"
	"strings"
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
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Arguments    map[string]any  `json:"arguments"`
	ExtraContent json.RawMessage `json:"extra_content,omitempty"` // Provider-specific passthrough metadata
}

// MalformedArgumentsKey is the sentinel key a protocol decoder stores in a
// ToolCall's Arguments when the provider sent a non-empty arguments payload that
// could not be parsed as JSON (typically a mid-JSON truncation at an output
// limit). Discarding the decode error would leave Arguments empty, which the
// runtime would misread as "the model intentionally passed no arguments" and run
// a success-shaped no-op. The sentinel lets the dispatch layer reject the call
// loudly instead. The value is the raw, undecodable payload.
const MalformedArgumentsKey = "__quine_malformed_arguments__"

// MalformedArguments reports whether this tool call carries undecodable
// arguments, returning the raw payload that failed to parse.
func (tc ToolCall) MalformedArguments() (raw string, ok bool) {
	if tc.Arguments == nil {
		return "", false
	}
	v, present := tc.Arguments[MalformedArgumentsKey]
	if !present {
		return "", false
	}
	s, _ := v.(string)
	return s, true
}

// ReasoningItem captures reasoning output from models that expose it
// (e.g., OpenAI Responses API with reasoning models).
type ReasoningItem struct {
	ID      string   `json:"id,omitempty"`
	Summary []string `json:"summary,omitempty"` // High-level reasoning summary
}

// ThinkingBlock captures an Anthropic extended-thinking block exactly as the
// model emitted it. Anthropic requires these blocks (with their signature) to be
// replayed before the tool_use of the same assistant turn; dropping them breaks
// multi-turn tool loops under interleaved thinking (HTTP 400) or silently loses
// reasoning. The block is opaque and provider-specific — preserve it verbatim.
type ThinkingBlock struct {
	Type      string `json:"type"`                // "thinking" or "redacted_thinking"
	Thinking  string `json:"thinking,omitempty"`  // for type="thinking"
	Signature string `json:"signature,omitempty"` // for type="thinking"
	Data      string `json:"data,omitempty"`      // for type="redacted_thinking"
}

// Message is a single turn in the conversation tape.
type Message struct {
	Role              Role            `json:"role"`
	Content           string          `json:"content"`
	StructuredContent json.RawMessage `json:"structured_content,omitempty"`
	ReasoningContent  string          `json:"reasoning_content,omitempty"` // For Chat Completions (e.g., o3)
	ReasoningItems    []ReasoningItem `json:"reasoning_items,omitempty"`   // For Responses API reasoning output
	ThinkingBlocks    []ThinkingBlock `json:"thinking_blocks,omitempty"`   // Anthropic extended-thinking blocks (replayed verbatim)
	ToolCalls         []ToolCall      `json:"tool_calls,omitempty"`
	ToolID            string          `json:"tool_id,omitempty"`
	Image             *ImagePart      `json:"image,omitempty"` // Optional image payload (vision tool results)
	Timestamp         int64           `json:"timestamp"`
}

// TerminationMode describes how a session ended.
type TerminationMode string

const (
	TermExit                 TerminationMode = "exit"
	TermContextDeath         TerminationMode = "context_death"
	TermContextExhaustion    TerminationMode = "context_exhaustion"
	TermTurnExhaustion       TerminationMode = "turn_exhaustion"
	TermTimeout              TerminationMode = "timeout"
	TermSignal               TerminationMode = "signal"
	TermRecoverableInference TerminationMode = "recoverable_inference"
	TermFinalizationFailure  TerminationMode = "finalization_failure"
	TermExec                 TerminationMode = "exec" // Process replaced via exec syscall
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

// ImagePart carries a base64-encoded image for delivery to vision-capable models.
// It is attached to a ToolResult (or Message) via the Image field and routed
// through the protocol layer as a native image content block.
type ImagePart struct {
	MIMEType string `json:"mime_type"` // e.g. "image/png"
	Data     string `json:"data"`      // base64-encoded bytes (no data: URI prefix)
}

// ToolResult holds the output of a tool execution.
type ToolResult struct {
	ToolID            string          `json:"tool_id"`
	Content           json.RawMessage `json:"content"`
	StructuredContent json.RawMessage `json:"structured_content,omitempty"`
	IsError           bool            `json:"is_error"`
	Image             *ImagePart      `json:"image,omitempty"` // Optional image payload (vision tool)
}

// SyntheticAssistantToolBatch rebuilds an assistant tool-call message after the
// original context line has been compacted or split away. Preserve reasoning
// metadata so providers that require assistant/tool-call reasoning continuity
// still accept the resumed batch.
func SyntheticAssistantToolBatch(source Message, calls []ToolCall) Message {
	msg := Message{
		Role:             RoleAssistant,
		Content:          "",
		ReasoningContent: source.ReasoningContent,
		ToolCalls:        append([]ToolCall(nil), calls...),
	}
	if len(source.ReasoningItems) > 0 {
		msg.ReasoningItems = append([]ReasoningItem(nil), source.ReasoningItems...)
	}
	return msg
}

func HasStructuredContent(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && !bytes.Equal(trimmed, []byte("null"))
}

func CompactStructuredJSON(raw json.RawMessage) string {
	if !HasStructuredContent(raw) {
		return ""
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, raw); err != nil {
		return ""
	}
	return buf.String()
}

// ToolResultModelContent returns the provider-facing tool result payload.
// When structure is available, send a compact JSON envelope that preserves the
// human-readable text while exposing the machine-readable fields.
func ToolResultModelContent(content string, structured json.RawMessage) string {
	compact := CompactStructuredJSON(structured)
	if compact == "" {
		return content
	}
	if strings.TrimSpace(content) == "" {
		return compact
	}

	payload := struct {
		Text       string          `json:"text,omitempty"`
		Structured json.RawMessage `json:"structured"`
	}{
		Structured: json.RawMessage(compact),
	}
	if strings.TrimSpace(content) != "" {
		payload.Text = content
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return content
	}
	return string(data)
}

// Tape is the in-memory structured trace for one runtime session.
// It backs retained trace artifacts and any later exports derived from them,
// not the canonical live provider-input contract, which is assembled from
// context/.
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

// SetSystemPrompt replaces the content of the first message (which must be RoleSystem).
// Used by escalation to update the system prompt after a model hot-swap.
func (t *Tape) SetSystemPrompt(content string) {
	if len(t.messages) > 0 && t.messages[0].Role == RoleSystem {
		t.messages[0].Content = content
	}
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
