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
// 1. Message conversion tests
// ---------------------------------------------------------------------------

func TestGeminiConvertMessages(t *testing.T) {
	msgs := []tape.Message{
		{Role: tape.RoleSystem, Content: "You are a helpful assistant."},
		{Role: tape.RoleSystem, Content: "Be concise."},
		{Role: tape.RoleUser, Content: "Hello"},
		{
			Role:    tape.RoleAssistant,
			Content: "Let me check.",
			ToolCalls: []tape.ToolCall{
				{ID: "tc_1", Name: "sh", Arguments: map[string]any{"command": "ls"}},
			},
		},
		{Role: tape.RoleToolResult, Content: "file.txt", ToolID: "sh"},
	}

	sysInst, contents := geminiConvertMessages(msgs)

	// System instruction should be set with two parts.
	if sysInst == nil {
		t.Fatal("system_instruction is nil")
	}
	if len(sysInst.Parts) != 2 {
		t.Fatalf("system_instruction parts = %d, want 2", len(sysInst.Parts))
	}
	if sysInst.Parts[0].Text != "You are a helpful assistant." {
		t.Errorf("system part[0] = %q", sysInst.Parts[0].Text)
	}
	if sysInst.Parts[1].Text != "Be concise." {
		t.Errorf("system part[1] = %q", sysInst.Parts[1].Text)
	}

	// Should have 3 content messages (user, model, user/functionResponse).
	if len(contents) != 3 {
		t.Fatalf("got %d contents, want 3", len(contents))
	}

	// User message.
	if contents[0].Role != "user" {
		t.Errorf("contents[0].Role = %q, want %q", contents[0].Role, "user")
	}
	if len(contents[0].Parts) != 1 || contents[0].Parts[0].Text != "Hello" {
		t.Errorf("contents[0].Parts = %+v", contents[0].Parts)
	}

	// Model message with text and functionCall.
	if contents[1].Role != "model" {
		t.Errorf("contents[1].Role = %q, want %q", contents[1].Role, "model")
	}
	if len(contents[1].Parts) != 2 {
		t.Fatalf("contents[1] parts = %d, want 2", len(contents[1].Parts))
	}
	if contents[1].Parts[0].Text != "Let me check." {
		t.Errorf("text part = %q", contents[1].Parts[0].Text)
	}
	if contents[1].Parts[1].FunctionCall == nil {
		t.Fatal("functionCall part is nil")
	}
	if contents[1].Parts[1].FunctionCall.Name != "sh" {
		t.Errorf("functionCall name = %q", contents[1].Parts[1].FunctionCall.Name)
	}
	if contents[1].Parts[1].FunctionCall.Args["command"] != "ls" {
		t.Errorf("functionCall args = %v", contents[1].Parts[1].FunctionCall.Args)
	}

	// FunctionResponse message.
	if contents[2].Role != "user" {
		t.Errorf("contents[2].Role = %q, want %q", contents[2].Role, "user")
	}
	if len(contents[2].Parts) != 1 {
		t.Fatalf("contents[2] parts = %d, want 1", len(contents[2].Parts))
	}
	fr := contents[2].Parts[0].FunctionResponse
	if fr == nil {
		t.Fatal("functionResponse is nil")
	}
	if fr.Name != "sh" {
		t.Errorf("functionResponse name = %q", fr.Name)
	}
	if fr.Response["content"] != "file.txt" {
		t.Errorf("functionResponse content = %v", fr.Response["content"])
	}
}

func TestGeminiConvertMessages_NoSystem(t *testing.T) {
	msgs := []tape.Message{
		{Role: tape.RoleUser, Content: "Hi"},
	}
	sysInst, contents := geminiConvertMessages(msgs)
	if sysInst != nil {
		t.Errorf("expected nil system_instruction, got %+v", sysInst)
	}
	if len(contents) != 1 {
		t.Fatalf("got %d contents, want 1", len(contents))
	}
}

// ---------------------------------------------------------------------------
// 2. Tool schema conversion
// ---------------------------------------------------------------------------

func TestGeminiConvertTools(t *testing.T) {
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
		{
			Name:        "read",
			Description: "Read a file",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{"type": "string"},
				},
			},
		},
	}

	out := geminiConvertTools(tools)

	if len(out) != 1 {
		t.Fatalf("got %d tool sets, want 1", len(out))
	}
	decls := out[0].FunctionDeclarations
	if len(decls) != 2 {
		t.Fatalf("got %d declarations, want 2", len(decls))
	}
	if decls[0].Name != "sh" {
		t.Errorf("decl[0].Name = %q", decls[0].Name)
	}
	if decls[0].Description != "Run a shell command" {
		t.Errorf("decl[0].Description = %q", decls[0].Description)
	}
	if decls[0].Parameters["type"] != "object" {
		t.Errorf("decl[0].Parameters.type = %v", decls[0].Parameters["type"])
	}
	if decls[1].Name != "read" {
		t.Errorf("decl[1].Name = %q", decls[1].Name)
	}
}

func TestGeminiConvertTools_Empty(t *testing.T) {
	out := geminiConvertTools(nil)
	if out != nil {
		t.Errorf("expected nil for empty tools, got %v", out)
	}
}

// ---------------------------------------------------------------------------
// 3. Generate with mock server
// ---------------------------------------------------------------------------

func TestGeminiGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify API key is in query parameter.
		apiKey := r.URL.Query().Get("key")
		if apiKey != "test-google-key" {
			t.Errorf("API key = %q, want %q", apiKey, "test-google-key")
		}

		// Verify content-type header.
		if r.Header.Get("content-type") != "application/json" {
			t.Errorf("content-type = %q", r.Header.Get("content-type"))
		}

		// Verify request body structure.
		body, _ := io.ReadAll(r.Body)
		var req geminiRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad request body: %v", err)
		}

		// System instruction should be present.
		if req.SystemInstruction == nil {
			t.Error("system_instruction is nil")
		} else if len(req.SystemInstruction.Parts) != 1 || req.SystemInstruction.Parts[0].Text != "You are helpful." {
			t.Errorf("system_instruction = %+v", req.SystemInstruction)
		}

		// Contents should have 1 user message.
		if len(req.Contents) != 1 {
			t.Errorf("contents len = %d, want 1", len(req.Contents))
		}

		// Tools should be present.
		if len(req.Tools) != 1 || len(req.Tools[0].FunctionDeclarations) != 1 {
			t.Errorf("tools = %+v", req.Tools)
		}

		// Return a response with text and functionCall.
		resp := `{
			"candidates": [{
				"content": {
					"parts": [
						{"text": "I'll run that command."},
						{"functionCall": {"name": "sh", "args": {"command": "ls -la"}}}
					]
				}
			}],
			"usageMetadata": {
				"promptTokenCount": 100,
				"candidatesTokenCount": 50
			}
		}`
		w.WriteHeader(200)
		w.Write([]byte(resp))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Provider: "google",
		APIKey:   "test-google-key",
		APIBase:  srv.URL,
		ModelID:  "gemini-2.0-flash",
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
	if msg.ToolCalls[0].Name != "sh" {
		t.Errorf("tool_call name = %q", msg.ToolCalls[0].Name)
	}
	if msg.ToolCalls[0].Arguments["command"] != "ls -la" {
		t.Errorf("arguments = %v", msg.ToolCalls[0].Arguments)
	}
	if usage.InputTokens != 100 || usage.OutputTokens != 50 {
		t.Errorf("usage = %+v", usage)
	}
}

func TestGeminiGenerate_AuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"error":{"code":403,"message":"API key not valid","status":"PERMISSION_DENIED"}}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Provider: "google",
		APIKey:   "bad-key",
		APIBase:  srv.URL,
		ModelID:  "gemini-2.0-flash",
	}
	p, _ := NewProvider(cfg)
	_, _, err := p.Generate([]tape.Message{{Role: tape.RoleUser, Content: "hi"}}, nil)
	if err != ErrAuth {
		t.Errorf("err = %v, want ErrAuth", err)
	}
}

// ---------------------------------------------------------------------------
// 4. NewProvider factory
// ---------------------------------------------------------------------------

func TestNewProviderGoogle(t *testing.T) {
	cfg := &config.Config{
		Provider: "google",
		APIKey:   "test-key",
		ModelID:  "gemini-2.0-flash",
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("provider is nil")
	}
	if p.ContextWindowSize() != 1_048_576 {
		t.Errorf("context window = %d, want 1048576", p.ContextWindowSize())
	}
}

func TestNewProviderGoogle_15Pro(t *testing.T) {
	cfg := &config.Config{
		Provider: "google",
		APIKey:   "test-key",
		ModelID:  "gemini-1.5-pro",
	}
	p, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ContextWindowSize() != 2_097_152 {
		t.Errorf("context window = %d, want 2097152", p.ContextWindowSize())
	}
}
