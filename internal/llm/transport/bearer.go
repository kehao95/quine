package transport

import "net/http"

// BearerToken implements Transport using Bearer token authentication.
// Used by: OpenAI, OpenRouter, Azure OpenAI.
type BearerToken struct {
	APIKey    string
	UserAgent string // Optional custom User-Agent header
}

func (t *BearerToken) Sign(req *http.Request, body []byte) error {
	req.Header.Set("Authorization", "Bearer "+t.APIKey)
	if t.UserAgent != "" {
		req.Header.Set("User-Agent", t.UserAgent)
	}
	return nil
}
