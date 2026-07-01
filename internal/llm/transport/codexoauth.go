package transport

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/codexoauth"
)

type CodexOAuthTransport struct {
	configDir string
	client    *http.Client
	mu        sync.Mutex
	state     codexoauthState
}

type codexoauthState struct {
	accessToken  string
	refreshToken string
	expiresAt    int64
	accountID    string
	idToken      string
}

func NewCodexOAuthTransport(_ *config.Config) (Transport, error) {
	configDir := os.Getenv("QUINE_CONFIG_DIR")
	return &CodexOAuthTransport{
		configDir: configDir,
		client:    &http.Client{Timeout: 2 * time.Minute},
	}, nil
}

func (t *CodexOAuthTransport) Sign(req *http.Request, body []byte) error {
	ts, err := t.ensureToken()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", codexoauth.AuthorizationHeader(ts.accessToken))
	if header, value := codexoauth.ChatGPTAccountIDHeader(ts.accountID); header != "" {
		req.Header.Set(header, value)
	}

	if shouldRewriteCodexEndpoint(req.URL.Path) {
		target, err := url.Parse(codexoauth.CodexEndpoint())
		if err != nil {
			return err
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.Host = target.Host
		req.Header.Set("Originator", "opencode")
	}

	return nil
}

func (t *CodexOAuthTransport) ensureToken() (codexoauthState, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UnixMilli()
	if t.state.accessToken != "" && t.state.expiresAt > now+30*1000 {
		return t.state, nil
	}

	if t.state.refreshToken != "" {
		req, err := t.refreshToken(t.state.refreshToken)
		if err == nil {
			// A refresh response may omit id_token, leaving the new state's
			// accountID empty. Carry the prior account ID forward so the
			// ChatGPT-Account-Id routing header is not silently dropped on all
			// subsequent requests (mirrors codexoauth.EnsureToken).
			if req.accountID == "" {
				req.accountID = t.state.accountID
			}
			t.state = req
			return t.state, nil
		}
	}

	loaded, ok, err := codexoauth.LoadToken(t.configDir)
	if err == nil && ok {
		if loaded.ExpiresAt > now+30*1000 {
			t.state = codexoauthStateFromToken(loaded)
			return t.state, nil
		}
		if loaded.RefreshToken != "" {
			refreshed, err := t.refreshToken(loaded.RefreshToken)
			if err == nil {
				if refreshed.accountID == "" {
					refreshed.accountID = codexoauthStateFromToken(loaded).accountID
				}
				t.state = refreshed
				return t.state, nil
			}
		}
	}

	newToken, err := codexoauth.EnsureToken(t.configDir, codexoauth.OpenBrowser)
	if err != nil {
		return codexoauthState{}, err
	}
	state := codexoauthStateFromToken(newToken)
	t.state = state
	return state, nil
}

func (t *CodexOAuthTransport) refreshToken(refreshToken string) (codexoauthState, error) {
	form := makeRefreshRequest(refreshToken)
	resp, err := t.client.Post("https://auth.openai.com/oauth/token", "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		return codexoauthState{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return codexoauthState{}, fmt.Errorf("token refresh failed: %d %s", resp.StatusCode, string(body))
	}

	state, err := parseTokenResponse(resp.Body)
	if err != nil {
		return codexoauthState{}, err
	}
	if err := codexoauth.SaveToken(t.configDir, codexoauthTokenFromState(state)); err != nil {
		return state, err
	}
	return state, nil
}

func makeRefreshRequest(refreshToken string) string {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", "app_EMoamEEZ73f0CkXaXp7hrann")
	return form.Encode()
}

func parseTokenResponse(r io.Reader) (codexoauthState, error) {
	var raw struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return codexoauthState{}, err
	}
	if raw.AccessToken == "" || raw.RefreshToken == "" {
		return codexoauthState{}, errors.New("missing access or refresh token")
	}
	expiresAt := time.Now().Add(time.Duration(raw.ExpiresIn) * time.Second).UnixMilli()
	if raw.ExpiresIn == 0 {
		expiresAt = time.Now().Add(1 * time.Hour).UnixMilli()
	}

	accountID := extractAccountID(raw.IDToken)
	if accountID == "" {
		accountID = extractAccountID(raw.AccessToken)
	}

	return codexoauthState{
		accessToken:  raw.AccessToken,
		refreshToken: raw.RefreshToken,
		expiresAt:    expiresAt,
		accountID:    accountID,
		idToken:      raw.IDToken,
	}, nil
}

func codexoauthStateFromToken(ts codexoauth.TokenState) codexoauthState {
	return codexoauthState{
		accessToken:  ts.AccessToken,
		refreshToken: ts.RefreshToken,
		expiresAt:    ts.ExpiresAt,
		accountID:    ts.AccountID,
		idToken:      ts.IDToken,
	}
}

func codexoauthTokenFromState(state codexoauthState) codexoauth.TokenState {
	return codexoauth.TokenState{
		AccessToken:  state.accessToken,
		RefreshToken: state.refreshToken,
		ExpiresAt:    state.expiresAt,
		AccountID:    state.accountID,
		IDToken:      state.idToken,
	}
}

func shouldRewriteCodexEndpoint(path string) bool {
	return strings.Contains(path, "/v1/responses") || strings.Contains(path, "/chat/completions")
}

func extractAccountID(token string) string {
	if token == "" {
		return ""
	}
	claims, err := parseJWTClaims(token)
	if err != nil {
		return ""
	}
	if v, ok := claims["chatgpt_account_id"].(string); ok && v != "" {
		return v
	}
	if v, ok := claims["https://api.openai.com/auth"].(map[string]any); ok {
		if vv, ok := v["chatgpt_account_id"].(string); ok && vv != "" {
			return vv
		}
	}
	if v, ok := claims["organizations"].([]any); ok && len(v) > 0 {
		if org, ok := v[0].(map[string]any); ok {
			if id, ok := org["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

func parseJWTClaims(token string) (map[string]any, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	return claims, nil
}
