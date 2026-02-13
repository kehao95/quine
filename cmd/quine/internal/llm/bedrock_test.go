package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

// ---------------------------------------------------------------------------
// Helper: set AWS env vars for bedrock tests
// ---------------------------------------------------------------------------

func setBedrockEnv(t *testing.T) {
	t.Helper()
	awsVars := []string{
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN", "AWS_DEFAULT_REGION", "AWS_REGION",
	}
	saved := make(map[string]string)
	for _, k := range awsVars {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
	}
	t.Cleanup(func() {
		for _, k := range awsVars {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})

	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	os.Setenv("AWS_SESSION_TOKEN", "FwoGZXIvYXdzEBYaDHqa0AP")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
}

// ---------------------------------------------------------------------------
// TestNewBedrockProvider – env var reading
// ---------------------------------------------------------------------------

func TestNewBedrockProvider(t *testing.T) {
	setBedrockEnv(t)

	cfg := &config.Config{
		Provider: "bedrock",
		ModelID:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
	}

	p, err := newBedrockProvider(cfg)
	if err != nil {
		t.Fatalf("newBedrockProvider: %v", err)
	}
	if p.region != "us-east-1" {
		t.Errorf("region = %q, want us-east-1", p.region)
	}
	if p.accessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("accessKey = %q", p.accessKey)
	}
	if p.sessionToken != "FwoGZXIvYXdzEBYaDHqa0AP" {
		t.Errorf("sessionToken = %q", p.sessionToken)
	}
	if !strings.Contains(p.endpoint, "bedrock-runtime.us-east-1.amazonaws.com") {
		t.Errorf("endpoint = %q", p.endpoint)
	}
}

func TestNewBedrockProvider_MissingRegion(t *testing.T) {
	awsVars := []string{
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN", "AWS_DEFAULT_REGION", "AWS_REGION",
	}
	saved := make(map[string]string)
	for _, k := range awsVars {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range awsVars {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})

	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")

	cfg := &config.Config{Provider: "bedrock", ModelID: "anthropic.claude-3-5-sonnet-20241022-v2:0"}
	_, err := newBedrockProvider(cfg)
	if err == nil || !strings.Contains(err.Error(), "AWS_DEFAULT_REGION") {
		t.Errorf("expected region error, got: %v", err)
	}
}

func TestNewBedrockProvider_MissingAccessKey(t *testing.T) {
	awsVars := []string{
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN", "AWS_DEFAULT_REGION", "AWS_REGION",
	}
	saved := make(map[string]string)
	for _, k := range awsVars {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range awsVars {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})

	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	// AWS_ACCESS_KEY_ID deliberately not set

	cfg := &config.Config{Provider: "bedrock", ModelID: "anthropic.claude-3-5-sonnet-20241022-v2:0"}
	_, err := newBedrockProvider(cfg)
	if err == nil || !strings.Contains(err.Error(), "AWS_ACCESS_KEY_ID") {
		t.Errorf("expected access key error, got: %v", err)
	}
}

func TestNewBedrockProvider_CustomEndpoint(t *testing.T) {
	setBedrockEnv(t)

	cfg := &config.Config{
		Provider: "bedrock",
		ModelID:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
		APIBase:  "http://localhost:9999",
	}
	p, err := newBedrockProvider(cfg)
	if err != nil {
		t.Fatalf("newBedrockProvider: %v", err)
	}
	if p.endpoint != "http://localhost:9999" {
		t.Errorf("endpoint = %q, want http://localhost:9999", p.endpoint)
	}
}

// ---------------------------------------------------------------------------
// TestBedrockContextWindowSize
// ---------------------------------------------------------------------------

func TestBedrockContextWindowSize(t *testing.T) {
	tests := []struct {
		model string
		want  int
	}{
		{"anthropic.claude-3-5-sonnet-20241022-v2:0", 200_000},
		{"anthropic.claude-sonnet-4-20250514-v1:0", 200_000},
		{"anthropic.claude-3-haiku-20240307-v1:0", 200_000},
		{"us.anthropic.claude-3-5-sonnet-20241022-v2:0", 200_000},
		{"anthropic.some-future-model", 200_000},
	}
	for _, tt := range tests {
		p := &bedrockProvider{modelID: tt.model}
		got := p.ContextWindowSize()
		if got != tt.want {
			t.Errorf("ContextWindowSize(%q) = %d, want %d", tt.model, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// TestBedrockGenerate – Full round-trip with httptest mock server
// ---------------------------------------------------------------------------

func TestBedrockGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL path contains model ID.
		if !strings.Contains(r.URL.Path, "/model/anthropic.claude-3-5-sonnet-20241022-v2:0/invoke") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify SigV4 headers are present.
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 ") {
			t.Errorf("Authorization header missing SigV4: %q", auth)
		}
		if !strings.Contains(auth, "Credential=AKIAIOSFODNN7EXAMPLE") {
			t.Errorf("Authorization missing access key: %q", auth)
		}
		if !strings.Contains(auth, "bedrock/aws4_request") {
			t.Errorf("Authorization missing bedrock service scope: %q", auth)
		}
		if r.Header.Get("X-Amz-Date") == "" {
			t.Error("missing X-Amz-Date header")
		}
		if r.Header.Get("X-Amz-Security-Token") != "FwoGZXIvYXdzEBYaDHqa0AP" {
			t.Errorf("X-Amz-Security-Token = %q", r.Header.Get("X-Amz-Security-Token"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}

		// Verify request body format (Bedrock-specific).
		body, _ := io.ReadAll(r.Body)
		var req bedrockRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad request body: %v", err)
		}
		if req.AnthropicVersion != "bedrock-2023-05-31" {
			t.Errorf("anthropic_version = %q", req.AnthropicVersion)
		}
		if req.System != "You are helpful." {
			t.Errorf("system = %q", req.System)
		}
		if req.MaxTokens != 4096 {
			t.Errorf("max_tokens = %d, want 4096", req.MaxTokens)
		}

		// Verify there is no "model" field in request body.
		var raw map[string]any
		json.Unmarshal(body, &raw)
		if _, hasModel := raw["model"]; hasModel {
			t.Error("bedrock request body should NOT contain 'model' field")
		}

		// Return Anthropic-format response.
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

	setBedrockEnv(t)

	cfg := &config.Config{
		Provider: "bedrock",
		ModelID:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
		APIBase:  srv.URL,
	}
	p, err := newBedrockProvider(cfg)
	if err != nil {
		t.Fatalf("newBedrockProvider: %v", err)
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

// ---------------------------------------------------------------------------
// TestBedrockGenerate_AuthError
// ---------------------------------------------------------------------------

func TestBedrockGenerate_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"message":"The security token included in the request is invalid.","__type":"UnrecognizedClientException"}`))
	}))
	defer srv.Close()

	setBedrockEnv(t)

	cfg := &config.Config{
		Provider: "bedrock",
		ModelID:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
		APIBase:  srv.URL,
	}
	p, err := newBedrockProvider(cfg)
	if err != nil {
		t.Fatalf("newBedrockProvider: %v", err)
	}

	_, _, err = p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "hi"}}, nil)
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

// ---------------------------------------------------------------------------
// TestBedrockClassifyError
// ---------------------------------------------------------------------------

func TestBedrockClassifyError_Auth(t *testing.T) {
	err := classifyBedrockError(403, []byte(`{"message":"not authorized","__type":"AccessDeniedException"}`))
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

func TestBedrockClassifyError_SecurityToken(t *testing.T) {
	err := classifyBedrockError(400, []byte(`{"message":"The security token included in the request is invalid.","__type":"UnrecognizedClientException"}`))
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

func TestBedrockClassifyError_AnthropicContextOverflow(t *testing.T) {
	body := `{"error":{"type":"invalid_request_error","message":"prompt is too long: your prompt has too many tokens"}}`
	err := classifyBedrockError(400, []byte(body))
	if err != ErrContextOverflow {
		t.Errorf("err = %v, want ErrContextOverflow", err)
	}
}

func TestBedrockClassifyError_AWSValidation(t *testing.T) {
	body := `{"message":"1 validation error detected","__type":"ValidationException"}`
	err := classifyBedrockError(400, []byte(body))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "ValidationException") {
		t.Errorf("error should contain type: %v", err)
	}
}

func TestBedrockClassifyError_GenericBody(t *testing.T) {
	err := classifyBedrockError(500, []byte(`something went wrong`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("error should mention HTTP 500: %v", err)
	}
}

// ---------------------------------------------------------------------------
// TestNewProviderBedrock – Factory returns bedrock provider
// ---------------------------------------------------------------------------

func TestNewProviderBedrock(t *testing.T) {
	setBedrockEnv(t)

	cfg := &config.Config{
		Provider: "bedrock",
		ModelID:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("provider is nil")
	}
	if _, ok := p.(*bedrockProvider); !ok {
		t.Errorf("expected *bedrockProvider, got %T", p)
	}
	if p.ContextWindowSize() != 200_000 {
		t.Errorf("context window = %d, want 200000", p.ContextWindowSize())
	}
}

// ---------------------------------------------------------------------------
// TestBedrockSigV4Headers – Verify signing produces correct header structure
// ---------------------------------------------------------------------------

func TestBedrockSigV4Headers(t *testing.T) {
	setBedrockEnv(t)

	cfg := &config.Config{
		Provider: "bedrock",
		ModelID:  "anthropic.claude-3-5-sonnet-20241022-v2:0",
	}
	p, err := newBedrockProvider(cfg)
	if err != nil {
		t.Fatalf("newBedrockProvider: %v", err)
	}

	payload := []byte(`{"test":"data"}`)
	req, _ := http.NewRequest("POST", "https://bedrock-runtime.us-east-1.amazonaws.com/model/test/invoke", nil)
	req.Header.Set("Content-Type", "application/json")

	if err := p.signV4(req, payload); err != nil {
		t.Fatalf("signV4: %v", err)
	}

	// Check Authorization header format.
	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 ") {
		t.Errorf("bad auth prefix: %q", auth)
	}
	if !strings.Contains(auth, "Credential=AKIAIOSFODNN7EXAMPLE/") {
		t.Errorf("missing credential: %q", auth)
	}
	if !strings.Contains(auth, "/us-east-1/bedrock/aws4_request") {
		t.Errorf("missing scope: %q", auth)
	}
	if !strings.Contains(auth, "SignedHeaders=") {
		t.Errorf("missing SignedHeaders: %q", auth)
	}
	if !strings.Contains(auth, "Signature=") {
		t.Errorf("missing Signature: %q", auth)
	}

	// Check required SigV4 headers.
	if req.Header.Get("X-Amz-Date") == "" {
		t.Error("missing X-Amz-Date")
	}
	if req.Header.Get("X-Amz-Security-Token") == "" {
		t.Error("missing X-Amz-Security-Token")
	}
	if req.Header.Get("X-Amz-Content-Sha256") == "" {
		t.Error("missing X-Amz-Content-Sha256")
	}
}

// ---------------------------------------------------------------------------
// TestBedrockSigV4_NoSessionToken – Verify no token header when empty
// ---------------------------------------------------------------------------

func TestBedrockSigV4_NoSessionToken(t *testing.T) {
	p := &bedrockProvider{
		region:       "us-east-1",
		accessKey:    "AKIA",
		secretKey:    "secret",
		sessionToken: "",
		modelID:      "anthropic.claude-3-5-sonnet-20241022-v2:0",
		endpoint:     "https://bedrock-runtime.us-east-1.amazonaws.com",
	}

	payload := []byte(`{}`)
	req, _ := http.NewRequest("POST", "https://bedrock-runtime.us-east-1.amazonaws.com/model/test/invoke", nil)
	req.Header.Set("Content-Type", "application/json")

	if err := p.signV4(req, payload); err != nil {
		t.Fatalf("signV4: %v", err)
	}

	if req.Header.Get("X-Amz-Security-Token") != "" {
		t.Error("X-Amz-Security-Token should be empty when sessionToken is empty")
	}
	// Authorization should still be present.
	if req.Header.Get("Authorization") == "" {
		t.Error("missing Authorization header")
	}
}
