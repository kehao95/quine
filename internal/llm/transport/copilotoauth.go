package transport

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/copilotoauth"
)

type CopilotOAuthTransport struct {
	configDir     string
	clientID      string
	userAgent     string
	integrationID string
	interactionID string
	mu            sync.Mutex
	token         copilotoauth.TokenState
}

func NewCopilotOAuthTransport(cfg *config.Config) (Transport, error) {
	return &CopilotOAuthTransport{
		configDir:     os.Getenv("QUINE_CONFIG_DIR"),
		clientID:      copilotoauth.ResolveClientID(""),
		userAgent:     copilotUserAgent(cfg),
		integrationID: copilotIntegrationID(),
		interactionID: copilotInteractionID(),
	}, nil
}

func (t *CopilotOAuthTransport) Sign(req *http.Request, body []byte) error {
	ts, err := t.ensureToken()
	if err != nil {
		return err
	}

	// GitHub Copilot exposes Chat Completions at /chat/completions rather than
	// the usual /v1/chat/completions used by generic OpenAI-compatible servers.
	if req.URL != nil && req.URL.Path == "/v1/chat/completions" {
		req.URL.Path = "/chat/completions"
	}

	req.Header.Set("Authorization", copilotoauth.AuthorizationHeader(ts.AccessToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Openai-Intent", "conversation-agent")
	req.Header.Set("X-Initiator", "user")
	req.Header.Set("X-GitHub-Api-Version", "2026-01-09")
	req.Header.Set("Copilot-Integration-Id", t.integrationID)
	req.Header.Set("X-Interaction-Id", t.interactionID)
	if t.userAgent != "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	return nil
}

func (t *CopilotOAuthTransport) ensureToken() (copilotoauth.TokenState, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if strings.TrimSpace(t.token.AccessToken) != "" {
		return t.token, nil
	}
	if loaded, ok, err := copilotoauth.LoadToken(t.configDir); err == nil && ok {
		t.token = loaded
		return t.token, nil
	}
	ts, err := copilotoauth.EnsureToken(t.configDir, t.clientID, copilotoauth.OpenBrowser)
	if err != nil {
		return copilotoauth.TokenState{}, err
	}
	t.token = ts
	return t.token, nil
}

func copilotIntegrationID() string {
	if v := strings.TrimSpace(os.Getenv("GITHUB_COPILOT_INTEGRATION_ID")); v != "" {
		return v
	}
	return "copilot-developer-cli"
}

func copilotUserAgent(cfg *config.Config) string {
	if cfg != nil && strings.TrimSpace(cfg.UserAgent) != "" {
		return cfg.UserAgent
	}
	return fmt.Sprintf("quine (client/github/cli %s %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func copilotInteractionID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("quine-%d", time.Now().UnixNano())
}
