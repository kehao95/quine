package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

func init() {
	// Suppress retry log output during tests.
	stderrOut = io.Discard
}

// ---------------------------------------------------------------------------
// 1. NewProvider factory tests
// ---------------------------------------------------------------------------

func TestNewProvider_Anthropic(t *testing.T) {
	cfg := &config.Config{
		Provider:      "anthropic",
		APIKey:        "test-key",
		ModelID:       "claude-3-5-sonnet-20241022",
		ContextWindow: 200_000,
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

func TestNewProvider_OpenAI(t *testing.T) {
	cfg := &config.Config{
		Provider:      "openai",
		APIKey:        "test-key",
		ModelID:       "gpt-4o",
		ContextWindow: 128_000,
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("provider is nil")
	}
	if p.ContextWindowSize() != 128_000 {
		t.Errorf("context window = %d, want 128000", p.ContextWindowSize())
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
	if !strings.Contains(err.Error(), "unknown provider") {
		t.Errorf("error = %q, want 'unknown provider'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// 2. Generate integration tests with mock servers
// ---------------------------------------------------------------------------

func TestGenerate_Anthropic_FullRoundTrip(t *testing.T) {
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
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad request body: %v", err)
		}
		if req["system"] != "You are helpful." {
			t.Errorf("system = %q", req["system"])
		}
		if req["model"] != "claude-3-5-sonnet-20241022" {
			t.Errorf("model = %q", req["model"])
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
		Provider:      "anthropic",
		APIKey:        "test-key",
		APIBase:       srv.URL,
		ModelID:       "claude-3-5-sonnet-20241022",
		ContextWindow: 200_000,
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

func TestGenerate_OpenAI_FullRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers.
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Authorization = %q, want 'Bearer test-key'", auth)
		}

		// Return response.
		resp := `{
			"choices": [{
				"message": {
					"role": "assistant",
					"content": "Hello!",
					"tool_calls": [{
						"id": "call_abc",
						"type": "function",
						"function": {
							"name": "sh",
							"arguments": "{\"command\":\"ls\"}"
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
		Provider:      "openai",
		APIKey:        "test-key",
		APIBase:       srv.URL,
		ModelID:       "gpt-4o",
		ContextWindow: 128_000,
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	msgs := []tape.Message{
		{Role: tape.RoleUser, Content: "Hello"},
	}

	msg, usage, err := p.Generate(msgs, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if msg.Role != tape.RoleAssistant {
		t.Errorf("role = %q", msg.Role)
	}
	if msg.Content != "Hello!" {
		t.Errorf("content = %q", msg.Content)
	}
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("tool_calls len = %d", len(msg.ToolCalls))
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
		Provider:      "anthropic",
		APIKey:        "bad-key",
		APIBase:       srv.URL,
		ModelID:       "claude-3-5-sonnet-20241022",
		ContextWindow: 200_000,
	}
	p, _ := NewProvider(cfg)
	_, _, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "hi"}}, nil)
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

// ---------------------------------------------------------------------------
// 3. Retry logic tests
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
