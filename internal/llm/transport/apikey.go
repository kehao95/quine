package transport

import "net/http"

// APIKeyHeader implements Transport using a custom header for the API key.
// Used by: Anthropic (x-api-key header).
type APIKeyHeader struct {
	HeaderName   string
	APIKey       string
	ExtraHeaders map[string]string
}

func (t *APIKeyHeader) Sign(req *http.Request, body []byte) error {
	req.Header.Set(t.HeaderName, t.APIKey)
	for k, v := range t.ExtraHeaders {
		req.Header.Set(k, v)
	}
	return nil
}
