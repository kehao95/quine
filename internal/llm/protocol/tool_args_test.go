package protocol

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestDecodeToolArguments(t *testing.T) {
	t.Run("empty string yields empty map, not nil", func(t *testing.T) {
		got := decodeToolArguments("")
		if got == nil {
			t.Fatal("expected non-nil empty map")
		}
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
		if _, malformed := (tape.ToolCall{Arguments: got}).MalformedArguments(); malformed {
			t.Fatal("empty args must not be flagged malformed")
		}
	})

	t.Run("valid JSON decodes", func(t *testing.T) {
		got := decodeToolArguments(`{"command":"ls -la"}`)
		if got["command"] != "ls -la" {
			t.Fatalf("expected command, got %v", got)
		}
	})

	t.Run("truncated JSON is flagged, not silently emptied", func(t *testing.T) {
		raw := `{"command": "very long comm`
		got := decodeToolArguments(raw)
		stored, malformed := (tape.ToolCall{Arguments: got}).MalformedArguments()
		if !malformed {
			t.Fatal("truncated arguments must be flagged malformed")
		}
		if stored != raw {
			t.Fatalf("malformed marker must preserve raw payload, got %q", stored)
		}
	})
}

// TestAnthropicThinkingRoundTrip locks in D1: a signed thinking block must be
// captured on decode and replayed (signature included) before the turn's
// tool_use on re-encode, or Anthropic rejects interleaved-thinking tool loops.
func TestAnthropicThinkingRoundTrip(t *testing.T) {
	body := `{"content":[
		{"type":"thinking","thinking":"weighing options","signature":"sig-xyz"},
		{"type":"text","text":"Listing the directory."},
		{"type":"tool_use","id":"tu_1","name":"sh","input":{"command":"ls"}}
	]}`
	var resp anthropicResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	msg := parseAnthropicResponse(resp)
	if len(msg.ThinkingBlocks) != 1 {
		t.Fatalf("expected 1 thinking block captured, got %d", len(msg.ThinkingBlocks))
	}
	if msg.ThinkingBlocks[0].Signature != "sig-xyz" || msg.ThinkingBlocks[0].Thinking != "weighing options" {
		t.Fatalf("thinking block not captured verbatim: %+v", msg.ThinkingBlocks[0])
	}

	_, out := convertAnthropicMessages([]tape.Message{msg})
	if len(out) != 1 {
		t.Fatalf("expected 1 assistant message, got %d", len(out))
	}
	encoded, err := json.Marshal(out[0].Content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	s := string(encoded)
	thinkingIdx := strings.Index(s, `"type":"thinking"`)
	toolUseIdx := strings.Index(s, `"type":"tool_use"`)
	if thinkingIdx < 0 {
		t.Fatalf("re-encoded content missing thinking block: %s", s)
	}
	if !strings.Contains(s, "sig-xyz") {
		t.Fatalf("re-encoded thinking block lost its signature: %s", s)
	}
	if toolUseIdx < 0 || thinkingIdx > toolUseIdx {
		t.Fatalf("thinking block must precede tool_use: %s", s)
	}
}

// TestSanitizeToolArgumentsStripsSentinel pins the review fix: the internal
// malformed-arguments sentinel must never reach any provider wire.
func TestSanitizeToolArgumentsStripsSentinel(t *testing.T) {
	if got := marshalToolArguments(map[string]any{tape.MalformedArgumentsKey: "{trunc"}); got != "{}" {
		t.Fatalf("sentinel-only args must marshal to {}, got %q", got)
	}
	got := marshalToolArguments(map[string]any{"command": "ls", tape.MalformedArgumentsKey: "x"})
	if strings.Contains(got, "__quine_malformed") {
		t.Fatalf("sentinel leaked into marshaled args: %s", got)
	}
	if !strings.Contains(got, "command") {
		t.Fatalf("real arg dropped while stripping sentinel: %s", got)
	}
}

func TestAnthropicEncodeStripsMalformedSentinel(t *testing.T) {
	msg := tape.Message{
		Role:      tape.RoleAssistant,
		ToolCalls: []tape.ToolCall{{ID: "tu_1", Name: "sh", Arguments: map[string]any{tape.MalformedArgumentsKey: "{trunc"}}},
	}
	_, out := convertAnthropicMessages([]tape.Message{msg})
	encoded, _ := json.Marshal(out[0].Content)
	if strings.Contains(string(encoded), "__quine_malformed") {
		t.Fatalf("anthropic encode leaked the malformed sentinel: %s", encoded)
	}
}

// TestAnthropicDropsUnsignedThinkingBlock pins the review fix: a thinking block
// without a signature cannot be validly replayed and must be dropped, not emitted.
func TestAnthropicDropsUnsignedThinkingBlock(t *testing.T) {
	msg := tape.Message{
		Role:           tape.RoleAssistant,
		Content:        "hi",
		ThinkingBlocks: []tape.ThinkingBlock{{Type: "thinking", Thinking: "x", Signature: ""}},
		ToolCalls:      []tape.ToolCall{{ID: "tu_1", Name: "sh", Arguments: map[string]any{"command": "ls"}}},
	}
	_, out := convertAnthropicMessages([]tape.Message{msg})
	encoded, _ := json.Marshal(out[0].Content)
	if strings.Contains(string(encoded), "thinking") {
		t.Fatalf("unsigned thinking block must be dropped from encode: %s", encoded)
	}
}

// TestAnthropicNoThinkingUnchanged guarantees the D1 change is inert when no
// thinking blocks are present (thinking disabled): assistant content is exactly
// text + tool_use with no spurious thinking block.
func TestAnthropicNoThinkingUnchanged(t *testing.T) {
	msg := tape.Message{
		Role:      tape.RoleAssistant,
		Content:   "Hi",
		ToolCalls: []tape.ToolCall{{ID: "tu_1", Name: "sh", Arguments: map[string]any{"command": "ls"}}},
	}
	_, out := convertAnthropicMessages([]tape.Message{msg})
	encoded, err := json.Marshal(out[0].Content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	if strings.Contains(string(encoded), "thinking") {
		t.Fatalf("non-thinking flow must not emit a thinking block: %s", encoded)
	}
}
