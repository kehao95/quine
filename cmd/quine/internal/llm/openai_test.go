package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

// ---------------------------------------------------------------------------
// 1. TestOpenAIConvertMessages
// ---------------------------------------------------------------------------

func TestOpenAIConvertMessages(t *testing.T) {
	msgs := []tape.Message{
		{Role: tape.RoleSystem, Content: "You are a helpful assistant."},
		{Role: tape.RoleUser, Content: "Hello"},
		{
			Role:    tape.RoleAssistant,
			Content: "Let me check.",
			ToolCalls: []tape.ToolCall{
				{ID: "call_1", Name: "sh", Arguments: map[string]any{"command": "ls"}},
			},
		},
		{Role: tape.RoleToolResult, Content: "file.txt", ToolID: "call_1"},
	}

	apiMsgs := openaiConvertMessages(msgs)

	if len(apiMsgs) != 4 {
		t.Fatalf("got %d messages, want 4", len(apiMsgs))
	}

	// System stays in messages array (not extracted).
	if apiMsgs[0].Role != "system" || apiMsgs[0].Content != "You are a helpful assistant." {
		t.Errorf("msg[0] = %+v, want system message", apiMsgs[0])
	}

	// User message.
	if apiMsgs[1].Role != "user" || apiMsgs[1].Content != "Hello" {
		t.Errorf("msg[1] = %+v, want user message", apiMsgs[1])
	}

	// Assistant with tool_calls using function wrapper.
	if apiMsgs[2].Role != "assistant" {
		t.Errorf("msg[2].Role = %q, want assistant", apiMsgs[2].Role)
	}
	if apiMsgs[2].Content != "Let me check." {
		t.Errorf("msg[2].Content = %q, want 'Let me check.'", apiMsgs[2].Content)
	}
	if len(apiMsgs[2].ToolCalls) != 1 {
		t.Fatalf("msg[2].ToolCalls len = %d, want 1", len(apiMsgs[2].ToolCalls))
	}
	tc := apiMsgs[2].ToolCalls[0]
	if tc.ID != "call_1" {
		t.Errorf("tool_call.ID = %q, want call_1", tc.ID)
	}
	if tc.Type != "function" {
		t.Errorf("tool_call.Type = %q, want function", tc.Type)
	}
	if tc.Function.Name != "sh" {
		t.Errorf("tool_call.Function.Name = %q, want sh", tc.Function.Name)
	}

	// Tool result uses role:"tool" with tool_call_id.
	if apiMsgs[3].Role != "tool" {
		t.Errorf("msg[3].Role = %q, want tool", apiMsgs[3].Role)
	}
	if apiMsgs[3].ToolCallID != "call_1" {
		t.Errorf("msg[3].ToolCallID = %q, want call_1", apiMsgs[3].ToolCallID)
	}
	if apiMsgs[3].Content != "file.txt" {
		t.Errorf("msg[3].Content = %q, want file.txt", apiMsgs[3].Content)
	}
}

// ---------------------------------------------------------------------------
// 2. TestOpenAIConvertTools
// ---------------------------------------------------------------------------

func TestOpenAIConvertTools(t *testing.T) {
	tools := []ToolSchema{
		{
			Name:        "sh",
			Description: "Run a shell command",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]any{"type": "string"},
				},
				"required": []any{"command"},
			},
		},
	}

	out := openaiConvertTools(tools)

	if len(out) != 1 {
		t.Fatalf("got %d tools, want 1", len(out))
	}
	if out[0].Type != "function" {
		t.Errorf("tool.Type = %q, want function", out[0].Type)
	}
	if out[0].Function.Name != "sh" {
		t.Errorf("tool.Function.Name = %q, want sh", out[0].Function.Name)
	}
	if out[0].Function.Description != "Run a shell command" {
		t.Errorf("tool.Function.Description = %q", out[0].Function.Description)
	}
	if out[0].Function.Parameters == nil {
		t.Fatal("tool.Function.Parameters is nil")
	}
	if out[0].Function.Parameters["type"] != "object" {
		t.Errorf("parameters.type = %v", out[0].Function.Parameters["type"])
	}
}

func TestOpenAIConvertTools_NilParameters(t *testing.T) {
	tools := []ToolSchema{
		{Name: "noop", Description: "Does nothing"},
	}

	out := openaiConvertTools(tools)
	if out[0].Function.Parameters == nil {
		t.Fatal("expected default schema for nil parameters")
	}
	if out[0].Function.Parameters["type"] != "object" {
		t.Errorf("expected type=object, got %v", out[0].Function.Parameters["type"])
	}
}

func TestOpenAIConvertTools_Empty(t *testing.T) {
	out := openaiConvertTools(nil)
	if out != nil {
		t.Errorf("expected nil for empty tools, got %v", out)
	}
}

// ---------------------------------------------------------------------------
// 3. TestOpenAIGenerate – Full round-trip with httptest mock server
// ---------------------------------------------------------------------------

func TestOpenAIGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers.
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Authorization = %q, want 'Bearer test-key'", auth)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}

		// Verify request body.
		body, _ := io.ReadAll(r.Body)
		var req openaiRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad request body: %v", err)
		}
		if req.Model != "gpt-4o" {
			t.Errorf("model = %q, want gpt-4o", req.Model)
		}
		// System message should be in the messages array.
		if len(req.Messages) < 2 {
			t.Fatalf("expected at least 2 messages, got %d", len(req.Messages))
		}
		if req.Messages[0].Role != "system" || req.Messages[0].Content != "You are helpful." {
			t.Errorf("messages[0] = %+v, want system", req.Messages[0])
		}

		// Return a response with text and tool_calls.
		resp := `{
			"choices": [{
				"message": {
					"role": "assistant",
					"content": "I'll run that command.",
					"tool_calls": [{
						"id": "call_abc",
						"type": "function",
						"function": {
							"name": "sh",
							"arguments": "{\"command\":\"ls -la\"}"
						}
					}]
				},
				"finish_reason": "tool_calls"
			}],
			"usage": {"prompt_tokens": 100, "completion_tokens": 50}
		}`
		w.WriteHeader(200)
		w.Write([]byte(resp))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Provider: "openai",
		APIKey:   "test-key",
		APIBase:  srv.URL,
		ModelID:  "gpt-4o",
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	msgs := []tape.Message{
		{Role: tape.RoleSystem, Content: "You are helpful."},
		{Role: tape.RoleUser, Content: "Run ls"},
	}
	tools := []ToolSchema{
		{
			Name:        "sh",
			Description: "Run a command",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]any{"type": "string"},
				},
			},
		},
	}

	msg, usage, err := p.Generate(msgs, tools)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if msg.Role != tape.RoleAssistant {
		t.Errorf("role = %q", msg.Role)
	}
	if msg.Content != "I'll run that command." {
		t.Errorf("content = %q", msg.Content)
	}
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("tool_calls len = %d", len(msg.ToolCalls))
	}
	if msg.ToolCalls[0].ID != "call_abc" || msg.ToolCalls[0].Name != "sh" {
		t.Errorf("tool_call = %+v", msg.ToolCalls[0])
	}
	if msg.ToolCalls[0].Arguments["command"] != "ls -la" {
		t.Errorf("arguments = %v", msg.ToolCalls[0].Arguments)
	}
	if usage.InputTokens != 100 || usage.OutputTokens != 50 {
		t.Errorf("usage = %+v", usage)
	}
}

// ---------------------------------------------------------------------------
// 4. TestOpenAIToolArgsAreJSONString
// ---------------------------------------------------------------------------

func TestOpenAIToolArgsAreJSONString(t *testing.T) {
	msgs := []tape.Message{
		{
			Role: tape.RoleAssistant,
			ToolCalls: []tape.ToolCall{
				{
					ID:   "call_1",
					Name: "sh",
					Arguments: map[string]any{
						"command": "echo hello",
						"timeout": float64(30),
					},
				},
			},
		},
	}

	apiMsgs := openaiConvertMessages(msgs)

	if len(apiMsgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(apiMsgs))
	}
	if len(apiMsgs[0].ToolCalls) != 1 {
		t.Fatalf("got %d tool_calls, want 1", len(apiMsgs[0].ToolCalls))
	}

	argsStr := apiMsgs[0].ToolCalls[0].Function.Arguments

	// Arguments must be a JSON string, not a raw object.
	var parsed map[string]any
	if err := json.Unmarshal([]byte(argsStr), &parsed); err != nil {
		t.Fatalf("arguments is not valid JSON string: %v (got %q)", err, argsStr)
	}
	if parsed["command"] != "echo hello" {
		t.Errorf("command = %v", parsed["command"])
	}
	if parsed["timeout"] != float64(30) {
		t.Errorf("timeout = %v", parsed["timeout"])
	}

	// Verify the wire format: marshal the message and check arguments is a string.
	wireBytes, err := json.Marshal(apiMsgs[0])
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var wireMsg map[string]any
	json.Unmarshal(wireBytes, &wireMsg)
	toolCalls := wireMsg["tool_calls"].([]any)
	tc := toolCalls[0].(map[string]any)
	fn := tc["function"].(map[string]any)
	if _, ok := fn["arguments"].(string); !ok {
		t.Errorf("wire arguments should be a string, got %T", fn["arguments"])
	}
}

// ---------------------------------------------------------------------------
// 5. TestNewProviderOpenAI – Factory returns openai provider
// ---------------------------------------------------------------------------

func TestNewProviderOpenAI(t *testing.T) {
	cfg := &config.Config{
		Provider: "openai",
		APIKey:   "test-key",
		ModelID:  "gpt-4o",
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("provider is nil")
	}
	if _, ok := p.(*openaiProvider); !ok {
		t.Errorf("expected *openaiProvider, got %T", p)
	}
	if p.ContextWindowSize() != 128_000 {
		t.Errorf("context window = %d, want 128000", p.ContextWindowSize())
	}
}

// ---------------------------------------------------------------------------
// OpenAI context window sizes
// ---------------------------------------------------------------------------

func TestOpenAIContextWindowSizes(t *testing.T) {
	tests := []struct {
		model string
		want  int
	}{
		{"gpt-4o", 128_000},
		{"gpt-4o-mini", 128_000},
		{"gpt-4-turbo", 128_000},
		{"gpt-4-turbo-preview", 128_000},
		{"gpt-4", 8_192},
		{"gpt-4-0613", 8_192},
		{"gpt-3.5-turbo", 16_385},
		{"gpt-3.5-turbo-16k", 16_385},
		{"o1", 200_000},
		{"o1-preview", 200_000},
		{"o3", 200_000},
		{"o3-mini", 200_000},
		{"o4-mini", 200_000},
		{"some-unknown-model", 128_000},
	}

	for _, tt := range tests {
		p := &openaiProvider{modelID: tt.model}
		got := p.ContextWindowSize()
		if got != tt.want {
			t.Errorf("ContextWindowSize(%q) = %d, want %d", tt.model, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// OpenAI error classification
// ---------------------------------------------------------------------------

func TestOpenAIClassifyError_Auth(t *testing.T) {
	err := openaiClassifyError(401, []byte(`{"error":{"message":"invalid api key","type":"invalid_request_error"}}`))
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

func TestOpenAIClassifyError_ContextOverflow(t *testing.T) {
	body := `{"error":{"message":"This model's maximum context length is 128000 tokens","type":"invalid_request_error","code":"context_length_exceeded"}}`
	err := openaiClassifyError(400, []byte(body))
	if err != ErrContextOverflow {
		t.Errorf("err = %v, want ErrContextOverflow", err)
	}
}

func TestOpenAIGenerate_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":{"message":"Incorrect API key provided","type":"invalid_request_error"}}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Provider: "openai",
		APIKey:   "bad-key",
		APIBase:  srv.URL,
		ModelID:  "gpt-4o",
	}
	p, _ := NewProvider(cfg)
	_, _, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "hi"}}, nil)
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}
