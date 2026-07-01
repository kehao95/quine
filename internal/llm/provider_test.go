package llm

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/protocol"
	"github.com/kehao95/quine/internal/tape"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type readCloserFunc struct {
	read  func([]byte) (int, error)
	close func() error
}

func (r readCloserFunc) Read(p []byte) (int, error) {
	return r.read(p)
}

func (r readCloserFunc) Close() error {
	if r.close == nil {
		return nil
	}
	return r.close()
}

func init() {
	// Suppress retry log output during tests.
	stderrOut = io.Discard
}

func providerTestConfig(provider, apiBase, model string, contextWindow int) *config.Config {
	return &config.Config{
		Identity:  config.Identity{ModelID: model},
		Transport: config.Transport{Provider: provider, APIKey: "test-key", APIBase: apiBase},
		Limits:    config.Limits{ContextWindow: contextWindow},
	}
}

// ---------------------------------------------------------------------------
// 1. NewProvider factory tests
// ---------------------------------------------------------------------------

func TestNewProvider_Anthropic(t *testing.T) {
	cfg := providerTestConfig("anthropic", "", "claude-3-5-sonnet-20241022", 200_000)
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
	cfg := providerTestConfig("openai", "", "gpt-4o", 128_000)
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

func TestBuildEndpoint_OpenAICompatibleCustomBases(t *testing.T) {
	tests := []struct {
		name  string
		base  string
		proto protocol.Protocol
		want  string
	}{
		{
			name:  "standard v1 base",
			base:  "https://api.moonshot.ai/v1",
			proto: &protocol.OpenAIProtocol{},
			want:  "https://api.moonshot.ai/v1/chat/completions",
		},
		{
			name:  "google openai compatibility chat completions",
			base:  "https://generativelanguage.googleapis.com/v1beta/openai",
			proto: &protocol.OpenAIProtocol{},
			want:  "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions",
		},
		{
			name:  "google openai compatibility responses",
			base:  "https://generativelanguage.googleapis.com/v1beta/openai",
			proto: &protocol.OpenAIResponsesProtocol{},
			want:  "https://generativelanguage.googleapis.com/v1beta/openai/responses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Transport: config.Transport{APIBase: tt.base}}
			if got := buildEndpoint(cfg, tt.proto); got != tt.want {
				t.Fatalf("buildEndpoint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewHTTPClient_ConservativeTransport(t *testing.T) {
	client := newHTTPClient()
	if client.Timeout != 10*time.Minute {
		t.Fatalf("timeout = %v, want %v", client.Timeout, 10*time.Minute)
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("transport type = %T, want *http.Transport", client.Transport)
	}
	if !transport.DisableCompression {
		t.Fatal("DisableCompression should be true")
	}
	if transport.ForceAttemptHTTP2 {
		t.Fatal("ForceAttemptHTTP2 should be false")
	}
	if !transport.DisableKeepAlives {
		t.Fatal("DisableKeepAlives should be true")
	}
	if transport.MaxIdleConns != 0 {
		t.Fatalf("MaxIdleConns = %d, want 0", transport.MaxIdleConns)
	}
}

func TestNewProvider_Unsupported(t *testing.T) {
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "some-model"},
		Transport: config.Transport{Provider: "fakeprovider", APIKey: "test-key"},
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

	cfg := providerTestConfig("anthropic", srv.URL, "claude-3-5-sonnet-20241022", 200_000)
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

	cfg := providerTestConfig("openai", srv.URL, "gpt-4o", 128_000)
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

	cfg := providerTestConfig("anthropic", srv.URL, "claude-3-5-sonnet-20241022", 200_000)
	cfg.APIKey = "bad-key"
	p, _ := NewProvider(cfg)
	_, _, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "hi"}}, nil)
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

func TestGenerate_RetriesRecoverableBodyReadError(t *testing.T) {
	var calls int32
	cfg := providerTestConfig("openai", "https://example.invalid", "gpt-4o", 128_000)
	raw, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	p := raw.(*provider)
	p.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		call := atomic.AddInt32(&calls, 1)
		if call == 1 {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: readCloserFunc{
					read: func([]byte) (int, error) { return 0, io.ErrUnexpectedEOF },
				},
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"choices": [{"message": {"role": "assistant", "content": "recovered"}}],
				"usage": {"prompt_tokens": 10, "completion_tokens": 2}
			}`)),
		}, nil
	})

	msg, usage, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("calls = %d, want 2", got)
	}
	if msg.Content != "recovered" {
		t.Fatalf("content = %q, want recovered", msg.Content)
	}
	if usage.InputTokens != 10 || usage.OutputTokens != 2 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestGenerate_ReturnsRecoverableAfterBodyReadRetries(t *testing.T) {
	var calls int32
	cfg := providerTestConfig("openai", "https://example.invalid", "gpt-4o", 128_000)
	raw, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	p := raw.(*provider)
	p.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		atomic.AddInt32(&calls, 1)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: readCloserFunc{
				read: func([]byte) (int, error) { return 0, io.ErrUnexpectedEOF },
			},
		}, nil
	})

	_, _, err = p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if !errors.Is(err, ErrRecoverableInference) {
		t.Fatalf("err = %v, want ErrRecoverableInference", err)
	}
	if !strings.Contains(err.Error(), "reading response body: unexpected EOF") {
		t.Fatalf("err = %v, want body read cause", err)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("calls = %d, want 3", got)
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

func TestGenerate_OpenAIResponses_SynthesizesMissingCompletedEventFromDelta(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(200)
		io.WriteString(w, "event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"}\n\n")
	}))
	defer srv.Close()

	cfg := providerTestConfig("openai-responses", srv.URL, "gpt-5.4", 128_000)
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	msg, usage, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("calls = %d, want 1", got)
	}
	if msg.Content != "partial" {
		t.Fatalf("content = %q, want %q", msg.Content, "partial")
	}
	if usage.InputTokens != 0 || usage.OutputTokens != 0 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestGenerate_OpenAIResponses_StopsAfterUnusableTransientSSERetries(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(200)
		io.WriteString(w, "event: response.created\ndata: {\"type\":\"response.created\"}\n\n")
	}))
	defer srv.Close()

	cfg := providerTestConfig("openai-responses", srv.URL, "gpt-5.4", 128_000)
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	_, _, err = p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err == nil {
		t.Fatal("expected SSE parse error")
	}
	if !strings.Contains(err.Error(), "no response.completed event found in SSE stream") {
		t.Fatalf("err = %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("calls = %d, want 3", got)
	}
}

func TestGenerate_OpenAIResponses_FailedEventSurfacesUnderlyingError(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(200)
		io.WriteString(w, strings.Join([]string{
			"event: response.created",
			"data: {\"type\":\"response.created\"}",
			"",
			"event: error",
			"data: {\"type\":\"error\",\"error\":{\"type\":\"insufficient_quota\",\"code\":\"insufficient_quota\",\"message\":\"quota exhausted\"}}",
			"",
			"event: response.failed",
			"data: {\"type\":\"response.failed\",\"response\":{\"id\":\"resp_123\",\"error\":{\"code\":\"insufficient_quota\",\"message\":\"quota exhausted\"}}}",
			"",
		}, "\n"))
	}))
	defer srv.Close()

	cfg := providerTestConfig("openai-responses", srv.URL, "gpt-5.4", 128_000)
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	_, _, err = p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err == nil {
		t.Fatal("expected failed-response error")
	}
	if !strings.Contains(err.Error(), "response.failed") || !strings.Contains(err.Error(), "insufficient_quota") {
		t.Fatalf("err = %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("calls = %d, want 1", got)
	}
}

func TestGenerate_OpenAIResponses_HandlesLeadingSSEKeepaliveComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, strings.Join([]string{
			": keepalive",
			"",
			"event: response.completed",
			"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_123\",\"output\":[{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Hello from comment-prefixed SSE\"}]}],\"usage\":{\"input_tokens\":12,\"output_tokens\":7}}}",
			"",
		}, "\n"))
	}))
	defer srv.Close()

	cfg := providerTestConfig("openai-responses", srv.URL, "gpt-5.4", 128_000)
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	msg, usage, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if msg.Content != "Hello from comment-prefixed SSE" {
		t.Fatalf("content = %q, want %q", msg.Content, "Hello from comment-prefixed SSE")
	}
	if usage.InputTokens != 12 || usage.OutputTokens != 7 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestGenerate_OpenAIResponses_ReconstructsOutputFromOutputItemDone(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, strings.Join([]string{
			"event: response.output_item.done",
			"data: {\"type\":\"response.output_item.done\",\"item\":{\"id\":\"msg_123\",\"type\":\"message\",\"status\":\"completed\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Recovered from output_item.done\"}]},\"output_index\":1,\"sequence_number\":9}",
			"",
			"event: response.completed",
			"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_123\",\"output\":[],\"usage\":{\"input_tokens\":17,\"output_tokens\":45}}}",
			"",
		}, "\n"))
	}))
	defer srv.Close()

	cfg := providerTestConfig("openai-responses", srv.URL, "gpt-5.4", 128_000)
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	msg, usage, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if msg.Content != "Recovered from output_item.done" {
		t.Fatalf("content = %q, want %q", msg.Content, "Recovered from output_item.done")
	}
	if usage.InputTokens != 17 || usage.OutputTokens != 45 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestGenerate_OpenAIResponses_MergesReasoningSummaryFromSSE(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, strings.Join([]string{
			"event: response.reasoning_summary_text.done",
			"data: {\"type\":\"response.reasoning_summary_text.done\",\"item_id\":\"rs_123\",\"summary_index\":0,\"text\":\"First summary.\",\"output_index\":0,\"sequence_number\":1}",
			"",
			"event: response.reasoning_summary_text.done",
			"data: {\"type\":\"response.reasoning_summary_text.done\",\"item_id\":\"rs_123\",\"summary_index\":1,\"text\":\"Second summary.\",\"output_index\":0,\"sequence_number\":2}",
			"",
			"event: response.completed",
			"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_123\",\"output\":[{\"id\":\"rs_123\",\"type\":\"reasoning\"},{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Done\"}]}],\"usage\":{\"input_tokens\":12,\"output_tokens\":7}}}",
			"",
		}, "\n"))
	}))
	defer srv.Close()

	cfg := providerTestConfig("openai-responses", srv.URL, "gpt-5.4", 128_000)
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}

	msg, _, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "Hello"}}, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if msg.Content != "Done" {
		t.Fatalf("content = %q, want %q", msg.Content, "Done")
	}
	if len(msg.ReasoningItems) != 1 {
		t.Fatalf("reasoning item count = %d, want 1", len(msg.ReasoningItems))
	}
	if got := msg.ReasoningItems[0].ID; got != "rs_123" {
		t.Fatalf("reasoning id = %q, want %q", got, "rs_123")
	}
	want := []string{"First summary.", "Second summary."}
	if !reflect.DeepEqual(msg.ReasoningItems[0].Summary, want) {
		t.Fatalf("reasoning summary = %#v, want %#v", msg.ReasoningItems[0].Summary, want)
	}
}

func TestEffectiveThinkingBudget_DefaultsHigh(t *testing.T) {
	p := &provider{}
	if got := p.effectiveThinkingBudget(nil); got != "high" {
		t.Fatalf("effectiveThinkingBudget(nil) = %q, want high", got)
	}
}

func TestEffectiveThinkingBudget_CopilotGPT54ChatWithToolsFallsBackOff(t *testing.T) {
	cfg := providerTestConfig("openai", "https://api.business.githubcopilot.com/v1", "gpt-5.4", 128_000)
	cfg.ThinkingBudget = "high"
	raw, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	p := raw.(*provider)
	if got := p.effectiveThinkingBudget([]ToolSchema{{Name: "sh"}}); got != "off" {
		t.Fatalf("effectiveThinkingBudget(tools) = %q, want off", got)
	}
}

func TestEffectiveThinkingBudget_CopilotGPT54ResponsesKeepsHigh(t *testing.T) {
	cfg := providerTestConfig("openai-responses", "https://api.business.githubcopilot.com/v1", "gpt-5.4", 128_000)
	cfg.ThinkingBudget = "high"
	raw, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	p := raw.(*provider)
	if got := p.effectiveThinkingBudget([]ToolSchema{{Name: "sh"}}); got != "high" {
		t.Fatalf("effectiveThinkingBudget(responses tools) = %q, want high", got)
	}
}

func TestEffectiveThinkingBudget_CopilotOpusChatWithToolsKeepsHigh(t *testing.T) {
	cfg := providerTestConfig("openai", "https://api.business.githubcopilot.com/v1", "claude-opus-4.5", 128_000)
	cfg.ThinkingBudget = "high"
	raw, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	p := raw.(*provider)
	if got := p.effectiveThinkingBudget([]ToolSchema{{Name: "sh"}}); got != "high" {
		t.Fatalf("effectiveThinkingBudget(opus tools) = %q, want high", got)
	}
}
