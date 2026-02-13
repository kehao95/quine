package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

func init() {
	// Suppress retry log output during tests.
	stderrOut = io.Discard
}

// ---------------------------------------------------------------------------
// 1. Message conversion tests
// ---------------------------------------------------------------------------

func TestConvertMessages_SystemExtraction(t *testing.T) {
	msgs := []tape.Message{
		{Role: tape.RoleSystem, Content: "You are a helpful assistant."},
		{Role: tape.RoleUser, Content: "Hello"},
		{Role: tape.RoleAssistant, Content: "Hi there!"},
	}

	system, apiMsgs := convertMessages(msgs)

	if system != "You are a helpful assistant." {
		t.Errorf("system = %q, want %q", system, "You are a helpful assistant.")
	}
	if len(apiMsgs) != 2 {
		t.Fatalf("got %d messages, want 2", len(apiMsgs))
	}
	if apiMsgs[0].Role != "user" {
		t.Errorf("msg[0].Role = %q, want %q", apiMsgs[0].Role, "user")
	}
	// User messages use a plain string content.
	if s, ok := apiMsgs[0].Content.(string); !ok || s != "Hello" {
		t.Errorf("msg[0].Content = %v, want %q", apiMsgs[0].Content, "Hello")
	}
	if apiMsgs[1].Role != "assistant" {
		t.Errorf("msg[1].Role = %q, want %q", apiMsgs[1].Role, "assistant")
	}
}

func TestConvertMessages_MultipleSystemMessages(t *testing.T) {
	msgs := []tape.Message{
		{Role: tape.RoleSystem, Content: "First system."},
		{Role: tape.RoleSystem, Content: "Second system."},
		{Role: tape.RoleUser, Content: "Hi"},
	}

	system, apiMsgs := convertMessages(msgs)

	if system != "First system.\n\nSecond system." {
		t.Errorf("system = %q", system)
	}
	if len(apiMsgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(apiMsgs))
	}
}

func TestConvertMessages_ToolUseBlocks(t *testing.T) {
	msgs := []tape.Message{
		{
			Role:    tape.RoleAssistant,
			Content: "Let me check that.",
			ToolCalls: []tape.ToolCall{
				{ID: "tc_1", Name: "sh", Arguments: map[string]any{"command": "ls"}},
			},
		},
	}

	_, apiMsgs := convertMessages(msgs)

	if len(apiMsgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(apiMsgs))
	}
	blocks, ok := apiMsgs[0].Content.([]contentBlock)
	if !ok {
		t.Fatalf("expected []contentBlock, got %T", apiMsgs[0].Content)
	}
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	if blocks[0].Type != "text" || blocks[0].Text == nil || *blocks[0].Text != "Let me check that." {
		t.Errorf("block[0] = %+v", blocks[0])
	}
	if blocks[1].Type != "tool_use" || blocks[1].ID != "tc_1" || blocks[1].Name != "sh" {
		t.Errorf("block[1] = %+v", blocks[1])
	}
}

func TestConvertMessages_ToolResult(t *testing.T) {
	msgs := []tape.Message{
		{Role: tape.RoleToolResult, Content: "file.txt", ToolID: "tc_1"},
	}

	_, apiMsgs := convertMessages(msgs)

	if len(apiMsgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(apiMsgs))
	}
	if apiMsgs[0].Role != "user" {
		t.Errorf("role = %q, want %q", apiMsgs[0].Role, "user")
	}
	blocks, ok := apiMsgs[0].Content.([]contentBlock)
	if !ok {
		t.Fatalf("expected []contentBlock, got %T", apiMsgs[0].Content)
	}
	if len(blocks) != 1 {
		t.Fatalf("got %d blocks, want 1", len(blocks))
	}
	if blocks[0].Type != "tool_result" || blocks[0].ToolUseID != "tc_1" || blocks[0].Content != "file.txt" {
		t.Errorf("block = %+v", blocks[0])
	}
}

// ---------------------------------------------------------------------------
// 2. Tool schema conversion
// ---------------------------------------------------------------------------

func TestConvertTools(t *testing.T) {
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

	out := convertTools(tools)

	if len(out) != 1 {
		t.Fatalf("got %d tools, want 1", len(out))
	}
	if out[0].Name != "sh" {
		t.Errorf("name = %q", out[0].Name)
	}
	if out[0].Description != "Run a shell command" {
		t.Errorf("description = %q", out[0].Description)
	}
	if out[0].InputSchema == nil {
		t.Fatal("input_schema is nil")
	}
	if out[0].InputSchema["type"] != "object" {
		t.Errorf("input_schema.type = %v", out[0].InputSchema["type"])
	}
}

func TestConvertTools_NilParameters(t *testing.T) {
	tools := []ToolSchema{
		{Name: "noop", Description: "Does nothing"},
	}

	out := convertTools(tools)
	if out[0].InputSchema == nil {
		t.Fatal("expected default schema for nil parameters")
	}
	if out[0].InputSchema["type"] != "object" {
		t.Errorf("expected type=object, got %v", out[0].InputSchema["type"])
	}
}

func TestConvertTools_Empty(t *testing.T) {
	out := convertTools(nil)
	if out != nil {
		t.Errorf("expected nil for empty tools, got %v", out)
	}
}

// ---------------------------------------------------------------------------
// 3. NewProvider factory
// ---------------------------------------------------------------------------

func TestNewProvider_Anthropic(t *testing.T) {
	cfg := &config.Config{
		Provider: "anthropic",
		APIKey:   "test-key",
		ModelID:  "claude-3-5-sonnet-20241022",
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("provider is nil")
	}
	if p.ContextWindowSize() != 200_000 {
		t.Errorf("context window = %d, want 200000", p.ContextWindowSize())
	}
}

func TestNewProvider_Unsupported(t *testing.T) {
	cfg := &config.Config{
		Provider: "fakeprovider",
		APIKey:   "test-key",
		ModelID:  "some-model",
	}
	_, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("error = %q, want 'not yet implemented'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// 4. Retry logic with mock HTTP server
// ---------------------------------------------------------------------------

func TestRetryWithBackoff_SuccessOnFirst(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	resp, err := retryWithBackoff(3, func() (*http.Response, error) {
		return http.Get(srv.URL)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestRetryWithBackoff_429Retries(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(429)
			w.Write([]byte(`{"error":{"type":"rate_limit","message":"rate limited"}}`))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		return http.Get(srv.URL)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if c := atomic.LoadInt32(&calls); c < 3 {
		t.Errorf("expected at least 3 calls, got %d", c)
	}
}

func TestRetryWithBackoff_5xxRetries(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n <= 2 {
			w.WriteHeader(500)
			w.Write([]byte(`internal server error`))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		return http.Get(srv.URL)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// 5. Error classification
// ---------------------------------------------------------------------------

func TestClassifyError_Auth401(t *testing.T) {
	err := classifyError(401, []byte(`{"error":{"type":"authentication_error","message":"invalid api key"}}`))
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

func TestClassifyError_Auth403(t *testing.T) {
	err := classifyError(403, []byte(`{"error":{"type":"permission_error","message":"not allowed"}}`))
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

func TestClassifyError_ContextOverflow(t *testing.T) {
	body := `{"error":{"type":"invalid_request_error","message":"prompt is too long: your prompt has too many tokens"}}`
	err := classifyError(400, []byte(body))
	if err != ErrContextOverflow {
		t.Errorf("err = %v, want ErrContextOverflow", err)
	}
}

func TestClassifyError_Overloaded(t *testing.T) {
	body := `{"error":{"type":"overloaded","message":"Overloaded"}}`
	err := classifyError(529, []byte(body))
	if err != ErrContextOverflow {
		t.Errorf("err = %v, want ErrContextOverflow", err)
	}
}

func TestClassifyError_GenericServerError(t *testing.T) {
	err := classifyError(500, []byte(`{"error":{"type":"server_error","message":"internal error"}}`))
	if err == ErrAuth || err == ErrContextOverflow {
		t.Errorf("unexpected sentinel error: %v", err)
	}
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("error should mention HTTP 500: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Integration: Generate with mock server
// ---------------------------------------------------------------------------

func TestGenerate_FullRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers.
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("wrong anthropic-version: %s", r.Header.Get("anthropic-version"))
		}
		if r.Header.Get("content-type") != "application/json" {
			t.Errorf("wrong content-type: %s", r.Header.Get("content-type"))
		}

		// Verify request body.
		body, _ := io.ReadAll(r.Body)
		var req anthropicRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad request body: %v", err)
		}
		if req.System != "You are helpful." {
			t.Errorf("system = %q", req.System)
		}
		if req.Model != "claude-3-5-sonnet-20241022" {
			t.Errorf("model = %q", req.Model)
		}

		// Return a response with text and tool_use.
		resp := `{
			"content": [
				{"type": "text", "text": "I'll run that command."},
				{"type": "tool_use", "id": "tu_123", "name": "sh", "input": {"command": "ls -la"}}
			],
			"usage": {"input_tokens": 100, "output_tokens": 50},
			"stop_reason": "tool_use"
		}`
		w.WriteHeader(200)
		w.Write([]byte(resp))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Provider: "anthropic",
		APIKey:   "test-key",
		APIBase:  srv.URL,
		ModelID:  "claude-3-5-sonnet-20241022",
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
	if msg.ToolCalls[0].ID != "tu_123" || msg.ToolCalls[0].Name != "sh" {
		t.Errorf("tool_call = %+v", msg.ToolCalls[0])
	}
	if msg.ToolCalls[0].Arguments["command"] != "ls -la" {
		t.Errorf("arguments = %v", msg.ToolCalls[0].Arguments)
	}
	if usage.InputTokens != 100 || usage.OutputTokens != 50 {
		t.Errorf("usage = %+v", usage)
	}
}

func TestGenerate_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":{"type":"authentication_error","message":"invalid api key"}}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Provider: "anthropic",
		APIKey:   "bad-key",
		APIBase:  srv.URL,
		ModelID:  "claude-3-5-sonnet-20241022",
	}
	p, _ := NewProvider(cfg)
	_, _, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "hi"}}, nil)
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}
