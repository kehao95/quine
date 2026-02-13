// Package transport handles authentication and request signing for different providers.
package transport

import (
	"fmt"
	"net/http"
)

// Transport handles authentication for API requests.
type Transport interface {
	// Sign adds authentication headers/params to the request.
	Sign(req *http.Request, body []byte) error
}

// For returns the Transport implementation for a given API type.
// Only "openai" and "anthropic" are supported.
func For(apiType, apiKey string) (Transport, error) {
	switch apiType {
	case "anthropic":
		return &APIKeyHeader{
			HeaderName: "x-api-key",
			APIKey:     apiKey,
			ExtraHeaders: map[string]string{
				"anthropic-version": "2023-06-01",
			},
		}, nil
	case "openai":
		return &BearerToken{APIKey: apiKey}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", apiType)
	}
}
