package transport

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/claudeoauth"
)

type ClaudeOAuthTransport struct {
	configDir string
	userAgent string
	mu        sync.Mutex
	state     claudeoauth.TokenState
}

func NewClaudeOAuthTransport(cfg *config.Config) (Transport, error) {
	var userAgent string
	if cfg != nil {
		userAgent = strings.TrimSpace(cfg.UserAgent)
	}
	return &ClaudeOAuthTransport{
		configDir: os.Getenv("QUINE_CONFIG_DIR"),
		userAgent: userAgent,
	}, nil
}

func (t *ClaudeOAuthTransport) Sign(req *http.Request, body []byte) error {
	ts, err := t.ensureToken()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", claudeoauth.AuthorizationHeader(ts.AccessToken))
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", claudeoauth.AnthropicBeta)
	req.Header.Set("x-app", "cli")
	req.Header.Set("anthropic-dangerous-direct-browser-access", "true")
	if t.userAgent != "" {
		req.Header.Set("User-Agent", t.userAgent)
	} else {
		req.Header.Set("User-Agent", claudeoauth.UserAgent)
	}
	return nil
}

func (t *ClaudeOAuthTransport) ensureToken() (claudeoauth.TokenState, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UnixMilli()
	if t.state.AccessToken != "" && tokenUsable(t.state, now) {
		return t.state, nil
	}

	ts, err := claudeoauth.EnsureToken(t.configDir, claudeoauth.OpenBrowser)
	if err != nil {
		return claudeoauth.TokenState{}, err
	}
	t.state = ts
	return t.state, nil
}

func tokenUsable(ts claudeoauth.TokenState, now int64) bool {
	if ts.AccessToken == "" {
		return false
	}
	return ts.ExpiresAt == 0 || ts.ExpiresAt > now+30*1000
}
