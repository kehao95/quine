package transport

import "net/http"

// BearerToken implements Transport using Bearer token authentication.
// Used by: OpenAI, OpenRouter, Azure OpenAI.
type BearerToken struct {
	APIKey string
}

func (t *BearerToken) Sign(req *http.Request, body []byte) error {
	req.Header.Set("Authorization", "Bearer "+t.APIKey)
	return nil
}
