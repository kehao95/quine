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
// Special case: if apiKey is "codex-oauth", OAuth transport is used
// regardless of apiType (for Codex CLI integration).
func For(apiType, apiKey string, cfg *config.Config) (Transport, error) {
	// Special sentinel: "codex-oauth" API key triggers OAuth flow
	if apiKey == "codex-oauth" {
		return NewCodexOAuthTransport(cfg)
	}

	switch apiType {
	case "anthropic":
		return &APIKeyHeader{
			HeaderName: "x-api-key",
			APIKey:     apiKey,
			ExtraHeaders: map[string]string{
				"anthropic-version": "2023-06-01",
			},
		}, nil
	case "openai", "openai-responses":
		return &BearerToken{APIKey: apiKey}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", apiType)
	}
}
