package protocol

import (
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestConvertAnthropicMessages_UsesStructuredToolEnvelope(t *testing.T) {
	_, out := convertAnthropicMessages([]tape.Message{
		{
			Role:              tape.RoleToolResult,
			Content:           "human tool text",
			StructuredContent: []byte(`{"exit_code":0,"stdout":"ok"}`),
			ToolID:            "sh:0",
		},
	})

	if len(out) != 1 {
		t.Fatalf("got %d output messages, want 1", len(out))
	}
	blocks, ok := out[0].Content.([]contentBlock)
	if !ok || len(blocks) != 1 {
		t.Fatalf("content type = %T, want single tool_result block", out[0].Content)
	}
	got, ok := blocks[0].Content.(string)
	if !ok {
		t.Fatalf("tool_result content type = %T, want string", blocks[0].Content)
	}
	want := `{"text":"human tool text","structured":{"exit_code":0,"stdout":"ok"}}`
	if got != want {
		t.Fatalf("tool_result content = %q, want %q", got, want)
	}
}

func TestConvertOpenAIMessages_UsesStructuredToolEnvelope(t *testing.T) {
	out := convertOpenAIMessages([]tape.Message{
		{
			Role:              tape.RoleToolResult,
			Content:           "human tool text",
			StructuredContent: []byte(`{"exit_code":0,"stdout":"ok"}`),
			ToolID:            "sh:0",
		},
	}, false, false)

	if len(out) != 1 {
		t.Fatalf("got %d output messages, want 1", len(out))
	}
	want := `{"text":"human tool text","structured":{"exit_code":0,"stdout":"ok"}}`
	if out[0].Content != want {
		t.Fatalf("tool content = %#v, want %q", out[0].Content, want)
	}
}

func TestConvertResponsesInput_UsesStructuredToolEnvelope(t *testing.T) {
	input, _ := convertResponsesInput([]tape.Message{
		{
			Role:              tape.RoleToolResult,
			Content:           "human tool text",
			StructuredContent: []byte(`{"exit_code":0,"stdout":"ok"}`),
			ToolID:            "sh:0",
		},
	})

	if len(input) != 1 {
		t.Fatalf("got %d input items, want 1", len(input))
	}
	want := `{"text":"human tool text","structured":{"exit_code":0,"stdout":"ok"}}`
	if input[0].Output != want {
		t.Fatalf("tool output = %q, want %q", input[0].Output, want)
	}
}
