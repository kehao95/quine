package transport

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/claudeoauth"
)

func TestClaudeOAuthTransportSignAddsAnthropicOAuthHeaders(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	writeTransportClaudeCredentials(t, dir)

	tr, err := NewClaudeOAuthTransport(&config.Config{
		Transport: config.Transport{UserAgent: "quine-test"},
	})
	if err != nil {
		t.Fatalf("NewClaudeOAuthTransport() error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := tr.Sign(req, nil); err != nil {
		t.Fatalf("Sign() error: %v", err)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer stored-access" {
		t.Fatalf("Authorization = %q", got)
	}
	if got := req.Header.Get("anthropic-version"); got != "2023-06-01" {
		t.Fatalf("anthropic-version = %q", got)
	}
	if got := req.Header.Get("anthropic-beta"); got != claudeoauth.AnthropicBeta {
		t.Fatalf("anthropic-beta = %q", got)
	}
	if got := req.Header.Get("x-app"); got != "cli" {
		t.Fatalf("x-app = %q", got)
	}
	if got := req.Header.Get("anthropic-dangerous-direct-browser-access"); got != "true" {
		t.Fatalf("anthropic-dangerous-direct-browser-access = %q", got)
	}
	if got := req.Header.Get("User-Agent"); got != "quine-test" {
		t.Fatalf("User-Agent = %q", got)
	}
}

func TestClaudeOAuthTransportSignUsesDefaultAgentSDKUserAgent(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	writeTransportClaudeCredentials(t, dir)

	tr, err := NewClaudeOAuthTransport(&config.Config{})
	if err != nil {
		t.Fatalf("NewClaudeOAuthTransport() error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := tr.Sign(req, nil); err != nil {
		t.Fatalf("Sign() error: %v", err)
	}
	if got := req.Header.Get("User-Agent"); got != claudeoauth.UserAgent {
		t.Fatalf("User-Agent = %q, want %q", got, claudeoauth.UserAgent)
	}
}

func TestTransportForClaudeOAuthSentinel(t *testing.T) {
	tr, err := For("anthropic", "claude-oauth", &config.Config{})
	if err != nil {
		t.Fatalf("For() error: %v", err)
	}
	if _, ok := tr.(*ClaudeOAuthTransport); !ok {
		t.Fatalf("transport type = %T, want *ClaudeOAuthTransport", tr)
	}
}

func writeTransportClaudeCredentials(t *testing.T, dir string) {
	t.Helper()
	data, err := json.Marshal(map[string]any{
		"claudeAiOauth": map[string]any{
			"accessToken":  "stored-access",
			"refreshToken": "stored-refresh",
			"expiresAt":    float64(1893456000000),
			"scopes":       []string{"user:inference"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".credentials.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}
