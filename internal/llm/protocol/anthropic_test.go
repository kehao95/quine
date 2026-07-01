package protocol

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestConvertAnthropicMessages_TrimsAssistantTrailingWhitespace(t *testing.T) {
	_, out := convertAnthropicMessages([]tape.Message{
		{Role: tape.RoleAssistant, Content: "hello world  \n\t"},
	})

	if len(out) != 1 {
		t.Fatalf("got %d output messages, want 1", len(out))
	}
	blocks, ok := out[0].Content.([]contentBlock)
	if !ok {
		t.Fatalf("assistant content type = %T, want []contentBlock", out[0].Content)
	}
	if len(blocks) != 1 {
		t.Fatalf("got %d content blocks, want 1", len(blocks))
	}
	if blocks[0].Text == nil || *blocks[0].Text != "hello world" {
		t.Fatalf("assistant text = %#v, want trimmed content", blocks[0].Text)
	}
}

func TestAnthropicEncodeRequest_ClaudeAgentSDKCompatAddsAttributionBlocks(t *testing.T) {
	body, err := (&AnthropicProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleSystem, Content: "You are quine."},
			{Role: tape.RoleUser, Content: "Begin."},
		},
		nil,
		"claude-sonnet-4-6",
		1024,
		RequestOptions{ClaudeAgentSDKCompat: true},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error: %v", err)
	}

	var req struct {
		System []anthropicSystemBlock `json:"system"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if len(req.System) != 3 {
		t.Fatalf("system block count = %d, want 3: %s", len(req.System), string(body))
	}
	if !strings.HasPrefix(req.System[0].Text, "x-anthropic-billing-header: ") ||
		!strings.Contains(req.System[0].Text, "cc_entrypoint=sdk-cli") ||
		!strings.Contains(req.System[0].Text, "cch=") {
		t.Fatalf("billing block = %q", req.System[0].Text)
	}
	if req.System[1].Text != claudeAgentSDKIdentity {
		t.Fatalf("identity block = %q", req.System[1].Text)
	}
	if req.System[1].CacheControl["type"] != "ephemeral" || req.System[1].CacheControl["ttl"] != "1h" {
		t.Fatalf("identity cache_control = %#v", req.System[1].CacheControl)
	}
	if req.System[2].Text != "You are quine." {
		t.Fatalf("original system block = %q", req.System[2].Text)
	}
}

func TestAnthropicEncodeRequest_ClaudeAgentSDKCompatPreservesTrailingAssistantText(t *testing.T) {
	body, err := (&AnthropicProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleSystem, Content: "You are quine."},
			{Role: tape.RoleUser, Content: "Begin."},
			{Role: tape.RoleAssistant, Content: "hello world"},
		},
		nil,
		"claude-sonnet-4-6",
		1024,
		RequestOptions{ClaudeAgentSDKCompat: true},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error: %v", err)
	}

	var req struct {
		Messages []struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if len(req.Messages) != 2 {
		t.Fatalf("message count = %d, want trailing assistant preserved: %s", len(req.Messages), string(body))
	}
	if req.Messages[1].Role != "assistant" {
		t.Fatalf("last retained role = %q, want assistant", req.Messages[1].Role)
	}
	var blocks []contentBlock
	if err := json.Unmarshal(req.Messages[1].Content, &blocks); err != nil {
		t.Fatalf("unmarshal assistant content: %v", err)
	}
	if len(blocks) != 1 || blocks[0].Text == nil || *blocks[0].Text != "hello world" {
		t.Fatalf("assistant blocks = %+v, want preserved text prefill", blocks)
	}
}

func TestStripTrailingAnthropicAssistantPrefillPreservesToolUse(t *testing.T) {
	messages := []anthropicMessage{
		{Role: "user", Content: "Begin."},
		{
			Role: "assistant",
			Content: []contentBlock{
				{Type: "text", Text: strPtr("I'll inspect it.")},
				{Type: "tool_use", ID: "call_sh", Name: "sh", Input: map[string]any{"command": "pwd"}},
			},
		},
	}

	got := stripTrailingAnthropicAssistantPrefill(messages)
	if len(got) != 2 {
		t.Fatalf("message count = %d, want trailing tool_use assistant preserved", len(got))
	}
	if got[1].Role != "assistant" {
		t.Fatalf("last role = %q, want assistant", got[1].Role)
	}
	blocks, ok := got[1].Content.([]contentBlock)
	if !ok {
		t.Fatalf("assistant content type = %T, want []contentBlock", got[1].Content)
	}
	if len(blocks) != 2 || blocks[1].Type != "tool_use" || blocks[1].ID != "call_sh" {
		t.Fatalf("assistant blocks = %+v, want preserved tool_use", blocks)
	}
}

func TestAnthropicEncodeRequest_PreservesEmptyToolUseInput(t *testing.T) {
	body, err := (&AnthropicProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{ID: "call_idle", Name: "idle", Arguments: map[string]any{}},
				},
			},
		},
		nil,
		"claude-sonnet-4-6",
		1024,
		RequestOptions{},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error: %v", err)
	}

	var req struct {
		Messages []struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if len(req.Messages) != 2 {
		t.Fatalf("message count = %d, want 2: %s", len(req.Messages), string(body))
	}
	var content []map[string]json.RawMessage
	if err := json.Unmarshal(req.Messages[1].Content, &content); err != nil {
		t.Fatalf("unmarshal assistant content: %v", err)
	}
	block := content[0]
	if string(block["type"]) != `"tool_use"` {
		t.Fatalf("block type = %s, want tool_use", block["type"])
	}
	input, ok := block["input"]
	if !ok {
		t.Fatalf("tool_use block missing input: %s", string(body))
	}
	if string(input) != `{}` {
		t.Fatalf("tool_use input = %s, want {}", input)
	}
}

func TestAnthropicEncodeRequest_TextBlockOmitsToolUseFields(t *testing.T) {
	body, err := (&AnthropicProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{Role: tape.RoleAssistant, Content: "hello"},
		},
		nil,
		"claude-sonnet-4-6",
		1024,
		RequestOptions{},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error: %v", err)
	}

	var req struct {
		Messages []struct {
			Content json.RawMessage `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	var content []map[string]json.RawMessage
	if err := json.Unmarshal(req.Messages[1].Content, &content); err != nil {
		t.Fatalf("unmarshal assistant content: %v", err)
	}
	block := content[0]
	for _, key := range []string{"id", "name", "input", "tool_use_id", "source"} {
		if _, ok := block[key]; ok {
			t.Fatalf("text block unexpectedly has %q: %s", key, string(body))
		}
	}
}

func TestConvertAnthropicMessages_OmitsWhitespaceOnlyAssistantTextBeforeToolUse(t *testing.T) {
	_, out := convertAnthropicMessages([]tape.Message{
		{
			Role:    tape.RoleAssistant,
			Content: " \n\t",
			ToolCalls: []tape.ToolCall{
				{ID: "sh:0", Name: "sh", Arguments: map[string]any{"command": "pwd"}},
			},
		},
	})

	if len(out) != 1 {
		t.Fatalf("got %d output messages, want 1", len(out))
	}
	blocks, ok := out[0].Content.([]contentBlock)
	if !ok {
		t.Fatalf("assistant content type = %T, want []contentBlock", out[0].Content)
	}
	if len(blocks) != 1 {
		t.Fatalf("got %d content blocks, want only tool_use block", len(blocks))
	}
	if blocks[0].Type != "tool_use" {
		t.Fatalf("first block type = %q, want tool_use", blocks[0].Type)
	}
	if blocks[0].ID != "sh_0" {
		t.Fatalf("sanitized tool id = %q, want sh_0", blocks[0].ID)
	}
}
