package codexoauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	clientID       = "app_EMoamEEZ73f0CkXaXp7hrann"
	issuer         = "https://auth.openai.com"
	authorizePath  = "/oauth/authorize"
	tokenPath      = "/oauth/token"
	codexEndpoint  = "https://chatgpt.com/backend-api/codex/responses"
	callbackPort   = 1455
	callbackPath   = "/auth/callback"
	callbackCancel = "/cancel"
)

var oauthScopes = []string{
	"openid",
	"profile",
	"email",
	"offline_access",
}

type TokenState struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	AccountID    string `json:"account_id,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type pkceCodes struct {
	Verifier  string
	Challenge string
}

type pendingOAuth struct {
	pkce     pkceCodes
	state    string
	resolve  func(TokenState)
	reject   func(error)
	deadline *time.Timer
}

type CallbackServer struct {
	mu      sync.Mutex
	server  *http.Server
	pending *pendingOAuth
}

func NewCallbackServer() *CallbackServer {
	return &CallbackServer{}
}

func (s *CallbackServer) Start() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.server != nil {
		return callbackURL(), nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, s.handleCallback)
	mux.HandleFunc(callbackCancel, s.handleCancel)

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", callbackPort),
		Handler: mux,
	}

	ln, err := netListen("tcp", fmt.Sprintf("127.0.0.1:%d", callbackPort))
	if err != nil {
		return "", err
	}

	s.server = server
	go func() {
		_ = server.Serve(ln)
	}()

	return callbackURL(), nil
}

func (s *CallbackServer) Stop() {
	s.mu.Lock()
	server := s.server
	pending := s.pending
	s.server = nil
	s.pending = nil
	s.mu.Unlock()

	if pending != nil {
		if pending.deadline != nil {
			pending.deadline.Stop()
		}
		pending.reject(errors.New("OAuth callback server stopped"))
	}
	if server != nil {
		_ = server.Close()
	}
}

func (s *CallbackServer) WaitForCallback(pkce pkceCodes, state string) <-chan TokenState {
	result := make(chan TokenState, 1)
	deadline := time.NewTimer(5 * time.Minute)
	deadline.Stop()

	s.mu.Lock()
	if s.pending != nil {
		prev := s.pending
		s.pending = nil
		if prev.deadline != nil {
			prev.deadline.Stop()
		}
		prev.reject(errors.New("OAuth callback superseded"))
	}
	deadline.Reset(5 * time.Minute)
	pend := &pendingOAuth{
		pkce:  pkce,
		state: state,
		resolve: func(ts TokenState) {
			result <- ts
		},
		reject: func(err error) {
			_ = err
			close(result)
		},
		deadline: deadline,
	}
	s.pending = pend
	s.mu.Unlock()

	go func() {
		<-deadline.C
		s.mu.Lock()
		if s.pending == pend {
			s.pending = nil
		}
		s.mu.Unlock()
		pend.reject(errors.New("OAuth callback timeout - authorization took too long"))
	}()

	return result
}

func (s *CallbackServer) handleCancel(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	pend := s.pending
	s.pending = nil
	s.mu.Unlock()

	if pend != nil {
		if pend.deadline != nil {
			pend.deadline.Stop()
		}
		pend.reject(errors.New("Login cancelled"))
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Login cancelled"))
}

func (s *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	errStr := q.Get("error")
	errDesc := q.Get("error_description")

	s.mu.Lock()
	pend := s.pending
	s.mu.Unlock()

	if errStr != "" {
		msg := errDesc
		if msg == "" {
			msg = errStr
		}
		if pend != nil {
			if pend.deadline != nil {
				pend.deadline.Stop()
			}
			pend.reject(errors.New(msg))
		}
		writeHTML(w, htmlError(msg), http.StatusOK)
		return
	}

	if code == "" {
		msg := "Missing authorization code"
		if pend != nil {
			if pend.deadline != nil {
				pend.deadline.Stop()
			}
			pend.reject(errors.New(msg))
		}
		writeHTML(w, htmlError(msg), http.StatusBadRequest)
		return
	}

	if pend == nil || state != pend.state {
		msg := "Invalid state - potential CSRF attack"
		if pend != nil {
			if pend.deadline != nil {
				pend.deadline.Stop()
			}
			pend.reject(errors.New(msg))
		}
		writeHTML(w, htmlError(msg), http.StatusBadRequest)
		return
	}

	if pend.deadline != nil {
		pend.deadline.Stop()
	}

	s.mu.Lock()
	s.pending = nil
	s.mu.Unlock()

	go func() {
		tokens, err := exchangeCodeForTokens(code, pend.pkce)
		if err != nil {
			pend.reject(err)
			return
		}
		pend.resolve(tokens)
	}()

	writeHTML(w, htmlSuccess(), http.StatusOK)
}

func authorizeURL() (string, pkceCodes, string, error) {
	pkce, err := generatePKCE()
	if err != nil {
		return "", pkceCodes{}, "", err
	}
	state, err := generateState()
	if err != nil {
		return "", pkceCodes{}, "", err
	}

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientID)
	params.Set("redirect_uri", callbackURL())
	params.Set("scope", strings.Join(oauthScopes, " "))
	params.Set("code_challenge", pkce.Challenge)
	params.Set("code_challenge_method", "S256")
	params.Set("id_token_add_organizations", "true")
	params.Set("codex_cli_simplified_flow", "true")
	params.Set("state", state)
	params.Set("originator", "opencode")

	return issuer + authorizePath + "?" + params.Encode(), pkce, state, nil
}

func refreshAccessToken(refreshToken string) (TokenState, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", clientID)

	resp, err := http.PostForm(issuer+tokenPath, form)
	if err != nil {
		return TokenState{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return TokenState{}, fmt.Errorf("token refresh failed: %d %s", resp.StatusCode, string(body))
	}

	return parseTokenResponse(resp.Body)
}

func exchangeCodeForTokens(code string, pkce pkceCodes) (TokenState, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", callbackURL())
	form.Set("client_id", clientID)
	form.Set("code_verifier", pkce.Verifier)

	resp, err := http.PostForm(issuer+tokenPath, form)
	if err != nil {
		return TokenState{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return TokenState{}, fmt.Errorf("token exchange failed: %d %s", resp.StatusCode, string(body))
	}

	return parseTokenResponse(resp.Body)
}

func parseTokenResponse(r io.Reader) (TokenState, error) {
	var raw struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return TokenState{}, err
	}
	if raw.AccessToken == "" || raw.RefreshToken == "" {
		return TokenState{}, errors.New("missing access or refresh token")
	}

	expiresAt := time.Now().Add(time.Duration(raw.ExpiresIn) * time.Second).UnixMilli()
	if raw.ExpiresIn == 0 {
		expiresAt = time.Now().Add(1 * time.Hour).UnixMilli()
	}

	accountID := ExtractAccountID(raw.IDToken)
	if accountID == "" {
		accountID = ExtractAccountID(raw.AccessToken)
	}

	return TokenState{
		AccessToken:  raw.AccessToken,
		RefreshToken: raw.RefreshToken,
		ExpiresAt:    expiresAt,
		AccountID:    accountID,
		IDToken:      raw.IDToken,
	}, nil
}

func ExtractAccountID(token string) string {
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

func generatePKCE() (pkceCodes, error) {
	verifier, err := randomString(43)
	if err != nil {
		return pkceCodes{}, err
	}
	h := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(h[:])
	return pkceCodes{Verifier: verifier, Challenge: challenge}, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func randomString(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b), nil
}

func htmlSuccess() string {
	return "<!doctype html><html><head><title>Quine - Codex Authorization Successful</title>" +
		"<style>body{font-family:system-ui,-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#131010;color:#f1ecec}.container{text-align:center;padding:2rem}h1{color:#f1ecec;margin-bottom:1rem}p{color:#b7b1b1}</style>" +
		"</head><body><div class=\"container\"><h1>Authorization Successful</h1><p>You can close this window and return to Quine.</p></div><script>setTimeout(()=>window.close(),2000)</script></body></html>"
}

func htmlError(err string) string {
	return "<!doctype html><html><head><title>Quine - Codex Authorization Failed</title>" +
		"<style>body{font-family:system-ui,-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#131010;color:#f1ecec}.container{text-align:center;padding:2rem}h1{color:#fc533a;margin-bottom:1rem}p{color:#b7b1b1}.error{color:#ff917b;font-family:monospace;margin-top:1rem;padding:1rem;background:#3c140d;border-radius:0.5rem}</style>" +
		"</head><body><div class=\"container\"><h1>Authorization Failed</h1><p>An error occurred during authorization.</p><div class=\"error\">" + htmlEscape(err) + "</div></div></body></html>"
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(s)
}

func writeHTML(w http.ResponseWriter, html string, status int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(html))
}

func callbackURL() string {
	return fmt.Sprintf("http://localhost:%d%s", callbackPort, callbackPath)
}

func codexTokenPath(configDir string) string {
	base := configDir
	if base == "" {
		base = defaultConfigDir()
	}
	return filepath.Join(base, "codex-oauth.json")
}

func codexCLIAuthPath() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".codex", "auth.json")
	}
	return ""
}

func defaultConfigDir() string {
	if dir, ok := os.LookupEnv("QUINE_CONFIG_DIR"); ok && dir != "" {
		return dir
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", "quine")
	}
	return "."
}

func LoadToken(configDir string) (TokenState, bool, error) {
	if path := codexCLIAuthPath(); path != "" {
		if ts, ok, err := loadTokenFile(path); err != nil {
			return TokenState{}, false, err
		} else if ok {
			return ts, true, nil
		}
	}

	path := codexTokenPath(configDir)
	ts, ok, err := loadTokenFile(path)
	if err != nil || !ok {
		return ts, ok, err
	}
	if err := SaveToken(configDir, ts); err != nil {
		return TokenState{}, false, err
	}
	return ts, true, nil
}

func SaveToken(configDir string, ts TokenState) error {
	path := codexTokenPath(configDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(encodeCodexCLIAuth(ts), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func EnsureToken(configDir string, openBrowser func(string) error) (TokenState, error) {
	if ts, ok, err := LoadToken(configDir); err == nil && ok {
		if ts.ExpiresAt > time.Now().Add(30*time.Second).UnixMilli() {
			return ts, nil
		}
		if ts.RefreshToken != "" {
			refreshed, err := refreshAccessToken(ts.RefreshToken)
			if err == nil {
				if refreshed.AccountID == "" {
					refreshed.AccountID = ts.AccountID
				}
				_ = SaveToken(configDir, refreshed)
				return refreshed, nil
			}
		}
	}

	server := NewCallbackServer()
	_, err := server.Start()
	if err != nil {
		return TokenState{}, err
	}

	authURL, pkce, state, err := authorizeURL()
	if err != nil {
		server.Stop()
		return TokenState{}, err
	}

	callback := server.WaitForCallback(pkce, state)
	if err := openBrowser(authURL); err != nil {
		server.Stop()
		return TokenState{}, err
	}

	select {
	case ts, ok := <-callback:
		server.Stop()
		if !ok {
			return TokenState{}, errors.New("authorization failed")
		}
		if err := SaveToken(configDir, ts); err != nil {
			return TokenState{}, err
		}
		return ts, nil
	case <-time.After(5 * time.Minute):
		server.Stop()
		return TokenState{}, errors.New("authorization timed out")
	}
}

func CodexEndpoint() string {
	return codexEndpoint
}

func AuthorizationHeader(accessToken string) string {
	return "Bearer " + accessToken
}

func ChatGPTAccountIDHeader(accountID string) (string, string) {
	if accountID == "" {
		return "", ""
	}
	return "ChatGPT-Account-Id", accountID
}

type codexCLIAuthFile struct {
	AuthMode     string                 `json:"auth_mode"`
	OpenAIAPIKey any                    `json:"OPENAI_API_KEY"`
	LastRefresh  string                 `json:"last_refresh,omitempty"`
	Tokens       codexCLIAuthTokenState `json:"tokens"`
}

type codexCLIAuthTokenState struct {
	IDToken      string `json:"id_token,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccountID    string `json:"account_id,omitempty"`
}

func loadTokenFile(path string) (TokenState, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TokenState{}, false, nil
		}
		return TokenState{}, false, err
	}

	if ts, ok, err := decodeCodexCLIAuth(data); err != nil {
		return TokenState{}, false, err
	} else if ok {
		return ts, true, nil
	}

	var ts TokenState
	if err := json.Unmarshal(data, &ts); err != nil {
		return TokenState{}, false, err
	}
	if ts.AccessToken == "" || ts.RefreshToken == "" {
		return TokenState{}, false, errors.New("invalid token file")
	}
	if ts.ExpiresAt == 0 {
		ts.ExpiresAt = tokenExpiryFromJWT(ts.AccessToken)
	}
	if ts.AccountID == "" {
		ts.AccountID = ExtractAccountID(ts.IDToken)
		if ts.AccountID == "" {
			ts.AccountID = ExtractAccountID(ts.AccessToken)
		}
	}
	return ts, true, nil
}

func decodeCodexCLIAuth(data []byte) (TokenState, bool, error) {
	var raw codexCLIAuthFile
	if err := json.Unmarshal(data, &raw); err != nil {
		return TokenState{}, false, err
	}
	if raw.AuthMode == "" && raw.Tokens.AccessToken == "" && raw.Tokens.RefreshToken == "" {
		return TokenState{}, false, nil
	}
	if raw.Tokens.AccessToken == "" || raw.Tokens.RefreshToken == "" {
		return TokenState{}, false, errors.New("invalid codex auth file")
	}

	accountID := raw.Tokens.AccountID
	if accountID == "" {
		accountID = ExtractAccountID(raw.Tokens.IDToken)
		if accountID == "" {
			accountID = ExtractAccountID(raw.Tokens.AccessToken)
		}
	}

	return TokenState{
		AccessToken:  raw.Tokens.AccessToken,
		RefreshToken: raw.Tokens.RefreshToken,
		ExpiresAt:    tokenExpiryFromJWT(raw.Tokens.AccessToken),
		AccountID:    accountID,
		IDToken:      raw.Tokens.IDToken,
	}, true, nil
}

func encodeCodexCLIAuth(ts TokenState) codexCLIAuthFile {
	accountID := ts.AccountID
	if accountID == "" {
		accountID = ExtractAccountID(ts.IDToken)
		if accountID == "" {
			accountID = ExtractAccountID(ts.AccessToken)
		}
	}
	return codexCLIAuthFile{
		AuthMode:     "chatgpt",
		OpenAIAPIKey: nil,
		LastRefresh:  time.Now().UTC().Format(time.RFC3339Nano),
		Tokens: codexCLIAuthTokenState{
			IDToken:      ts.IDToken,
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken,
			AccountID:    accountID,
		},
	}
}

func tokenExpiryFromJWT(token string) int64 {
	claims, err := parseJWTClaims(token)
	if err != nil {
		return 0
	}
	exp, ok := claims["exp"].(float64)
	if !ok || exp <= 0 {
		return 0
	}
	return int64(exp * 1000)
}

// netListen isolates net.Listen for testing.
var netListen = net.Listen
