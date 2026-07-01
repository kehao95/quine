// Package transport handles authentication and request signing for different providers.
package transport

import (
	"fmt"
	"net/http"

	"github.com/kehao95/quine/internal/config"
)

// Transport handles authentication for API requests.
type Transport interface {
	// Sign adds authentication headers/params to the request.
	Sign(req *http.Request, body []byte) error
}

// For returns the Transport implementation for a given API type.
// Supported: "openai", "anthropic", "openai-responses".
//
// Special cases for OAuth:
//   - apiKey "codex-oauth" triggers Codex CLI OAuth flow
//   - apiKey "claude-oauth" uses Claude Code subscription OAuth credentials
//   - apiKey "kimi-oauth" triggers Kimi CLI OAuth flow (impersonates Kimi CLI)
//   - apiKey "copilot-oauth" triggers GitHub Copilot OAuth device flow
func For(apiType, apiKey string, cfg *config.Config) (Transport, error) {
	// Special sentinel: "codex-oauth" API key triggers Codex OAuth flow
	if apiKey == "codex-oauth" {
		return NewCodexOAuthTransport(cfg)
	}

	if apiKey == "claude-oauth" {
		return NewClaudeOAuthTransport(cfg)
	}

	// Special sentinel: "kimi-oauth" API key triggers Kimi OAuth flow
	if apiKey == "kimi-oauth" {
		return NewKimiOAuthTransport(cfg)
	}

	if apiKey == "copilot-oauth" {
		return NewCopilotOAuthTransport(cfg)
	}

	// Extract optional User-Agent from config
	var userAgent string
	if cfg != nil {
		userAgent = cfg.UserAgent
	}

	switch apiType {
	case "anthropic":
		extraHeaders := map[string]string{
			"anthropic-version": "2023-06-01",
		}
		if userAgent != "" {
			extraHeaders["User-Agent"] = userAgent
		}
		return &APIKeyHeader{
			HeaderName:   "x-api-key",
			APIKey:       apiKey,
			ExtraHeaders: extraHeaders,
		}, nil
	case "openai", "openai-responses":
		return &BearerToken{APIKey: apiKey, UserAgent: userAgent}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", apiType)
	}
}
