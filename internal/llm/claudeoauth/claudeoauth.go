package claudeoauth

import (
	"bytes"
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
	"time"
)

const (
	AnthropicBeta = "oauth-2025-04-20,interleaved-thinking-2025-05-14"
	UserAgent     = "claude-cli/2.1.178 (external, sdk-cli)"

	defaultAuthorizeURL = "https://claude.ai/oauth/authorize"
	defaultClientID     = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	defaultTokenURL     = "https://platform.claude.com/v1/oauth/token"
	callbackPath        = "/oauth/callback"

	sourceEnv               = "env"
	sourceClaudeCredentials = "claude_credentials"
	sourceQuineToken        = "quine_token"
)

var authorizeURL = defaultAuthorizeURL
var tokenURL = defaultTokenURL

type TokenState struct {
	AccessToken      string   `json:"access_token"`
	RefreshToken     string   `json:"refresh_token,omitempty"`
	ExpiresAt        int64    `json:"expires_at,omitempty"`
	Scopes           []string `json:"scopes,omitempty"`
	SubscriptionType string   `json:"subscription_type,omitempty"`
	RateLimitTier    string   `json:"rate_limit_tier,omitempty"`
	ClientID         string   `json:"client_id,omitempty"`

	sourcePath string
	sourceKind string
}

type pkceCodes struct {
	Verifier  string
	Challenge string
}

func LoadToken(configDir string) (TokenState, bool, error) {
	if ts, ok := tokenFromEnv(); ok {
		return ts, true, nil
	}

	for _, path := range credentialCandidatePaths(configDir) {
		ts, ok, err := loadTokenFile(path)
		if err != nil {
			return TokenState{}, false, err
		}
		if ok {
			return ts, true, nil
		}
	}

	return TokenState{}, false, nil
}

func SaveToken(configDir string, ts TokenState) error {
	if ts.sourceKind == sourceEnv {
		return nil
	}

	path := ts.sourcePath
	kind := ts.sourceKind
	if path == "" {
		path = quineTokenPath(configDir)
		kind = sourceQuineToken
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	var data []byte
	var err error
	switch kind {
	case sourceClaudeCredentials:
		data, err = encodeClaudeCredentials(path, ts)
	default:
		data, err = json.MarshalIndent(ts, "", "  ")
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func EnsureToken(configDir string, openBrowser func(string) error) (TokenState, error) {
	var prior TokenState
	if ts, ok, err := LoadToken(configDir); err != nil {
		return TokenState{}, err
	} else if ok {
		prior = ts
		if tokenUsable(ts, time.Now()) {
			return ts, nil
		}
		if ts.RefreshToken != "" {
			refreshed, err := RefreshToken(ts.RefreshToken, ts)
			if err == nil {
				if err := SaveToken(configDir, refreshed); err != nil {
					return TokenState{}, err
				}
				return refreshed, nil
			}
		}
	}

	if openBrowser == nil {
		return TokenState{}, errors.New("Claude OAuth login requires a browser opener")
	}
	ts, err := loginWithBrowser(openBrowser)
	if err != nil {
		return TokenState{}, err
	}
	ts.sourcePath = prior.sourcePath
	ts.sourceKind = prior.sourceKind
	if err := SaveToken(configDir, ts); err != nil {
		return TokenState{}, err
	}
	return ts, nil
}

func RefreshToken(refreshToken string, prior TokenState) (TokenState, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return TokenState{}, errors.New("missing Claude OAuth refresh token")
	}

	respBody, err := postTokenRequest(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     resolveClientID(prior.ClientID),
		"scope":         strings.Join(resolveScopes(prior.Scopes), " "),
	})
	if err != nil {
		return TokenState{}, err
	}

	refreshed, err := parseRefreshResponse(respBody)
	if err != nil {
		return TokenState{}, err
	}
	if len(refreshed.Scopes) == 0 {
		refreshed.Scopes = prior.Scopes
	}
	if refreshed.RefreshToken == "" {
		refreshed.RefreshToken = prior.RefreshToken
	}
	if refreshed.SubscriptionType == "" {
		refreshed.SubscriptionType = prior.SubscriptionType
	}
	if refreshed.RateLimitTier == "" {
		refreshed.RateLimitTier = prior.RateLimitTier
	}
	if refreshed.ClientID == "" {
		refreshed.ClientID = resolveClientID(prior.ClientID)
	}
	refreshed.sourcePath = prior.sourcePath
	refreshed.sourceKind = prior.sourceKind
	return refreshed, nil
}

func ExchangeCodeForTokens(code, codeVerifier, redirectURI string) (TokenState, error) {
	code = parseAuthorizationCode(code)
	if code == "" {
		return TokenState{}, errors.New("missing Claude OAuth authorization code")
	}
	if strings.TrimSpace(codeVerifier) == "" {
		return TokenState{}, errors.New("missing Claude OAuth PKCE verifier")
	}
	if strings.TrimSpace(redirectURI) == "" {
		return TokenState{}, errors.New("missing Claude OAuth redirect URI")
	}

	respBody, err := postTokenRequest(map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"code_verifier": codeVerifier,
		"client_id":     resolveClientID(""),
		"redirect_uri":  redirectURI,
	})
	if err != nil {
		return TokenState{}, err
	}
	ts, err := parseRefreshResponse(respBody)
	if err != nil {
		return TokenState{}, err
	}
	if ts.ClientID == "" {
		ts.ClientID = resolveClientID("")
	}
	return ts, nil
}

func AuthorizationHeader(accessToken string) string {
	return "Bearer " + accessToken
}

func tokenFromEnv() (TokenState, bool) {
	accessToken := strings.TrimSpace(os.Getenv("CLAUDE_CODE_OAUTH_TOKEN"))
	if accessToken == "" {
		return TokenState{}, false
	}
	return TokenState{
		AccessToken:      accessToken,
		Scopes:           scopesFromEnv(),
		SubscriptionType: strings.TrimSpace(os.Getenv("CLAUDE_CODE_SUBSCRIPTION_TYPE")),
		RateLimitTier:    strings.TrimSpace(os.Getenv("CLAUDE_CODE_RATE_LIMIT_TIER")),
		ClientID:         resolveClientID(""),
		sourceKind:       sourceEnv,
	}, true
}

func scopesFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("CLAUDE_CODE_OAUTH_SCOPES"))
	if raw == "" {
		return defaultScopes()
	}
	return strings.Fields(raw)
}

func resolveScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return defaultScopes()
	}
	return scopes
}

func defaultScopes() []string {
	return []string{
		"user:profile",
		"user:inference",
		"user:sessions:claude_code",
		"user:mcp_servers",
		"user:file_upload",
	}
}

func resolveClientID(explicit string) string {
	for _, candidate := range []string{
		strings.TrimSpace(explicit),
		strings.TrimSpace(os.Getenv("CLAUDE_CODE_OAUTH_CLIENT_ID")),
		defaultClientID,
	} {
		if candidate != "" {
			return candidate
		}
	}
	return defaultClientID
}

func loginWithBrowser(openBrowser func(string) error) (TokenState, error) {
	pkce, err := generatePKCE()
	if err != nil {
		return TokenState{}, err
	}
	state, err := randomString(32)
	if err != nil {
		return TokenState{}, err
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return TokenState{}, err
	}
	defer ln.Close()

	redirectURI := fmt.Sprintf("http://127.0.0.1:%d%s", ln.Addr().(*net.TCPAddr).Port, callbackPath)
	result := make(chan string, 1)
	errs := make(chan error, 1)

	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}
	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("state"); got != state {
			writeHTML(w, htmlError("invalid OAuth state"), http.StatusBadRequest)
			errs <- errors.New("Claude OAuth callback state mismatch")
			return
		}
		if oauthErr := query.Get("error"); oauthErr != "" {
			description := query.Get("error_description")
			if description == "" {
				description = oauthErr
			}
			writeHTML(w, htmlError(description), http.StatusOK)
			errs <- fmt.Errorf("Claude OAuth callback error: %s", description)
			return
		}
		code := query.Get("code")
		if code == "" {
			writeHTML(w, htmlError("missing authorization code"), http.StatusBadRequest)
			errs <- errors.New("Claude OAuth callback missing authorization code")
			return
		}
		writeHTML(w, htmlSuccess(), http.StatusOK)
		result <- code
	})
	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()
	defer server.Close()

	authURL := buildAuthorizeURL(redirectURI, pkce, state)
	fmt.Fprintf(os.Stderr, "\nTo authorize Quine with Claude OAuth, open:\n  %s\n\nWaiting for authorization callback...\n", authURL)
	if err := openBrowser(authURL); err != nil {
		fmt.Fprintf(os.Stderr, "Could not open browser automatically: %v\nOpen the URL above manually.\n", err)
	}

	select {
	case code := <-result:
		return ExchangeCodeForTokens(code, pkce.Verifier, redirectURI)
	case err := <-errs:
		return TokenState{}, err
	case <-time.After(5 * time.Minute):
		return TokenState{}, errors.New("Claude OAuth authorization timed out")
	}
}

func buildAuthorizeURL(redirectURI string, pkce pkceCodes, state string) string {
	params := url.Values{}
	params.Set("code", "true")
	params.Set("client_id", resolveClientID(""))
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", strings.Join(defaultScopes(), " "))
	params.Set("code_challenge", pkce.Challenge)
	params.Set("code_challenge_method", "S256")
	params.Set("state", state)
	return authorizeURL + "?" + params.Encode()
}

func postTokenRequest(payload map[string]string) ([]byte, error) {
	body, status, err := postTokenJSON(payload)
	if err == nil && status == http.StatusOK {
		return body, nil
	}
	firstErr := tokenEndpointError("Claude OAuth token request", status, body, err)

	body, status, err = postTokenForm(payload)
	if err == nil && status == http.StatusOK {
		return body, nil
	}
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		return nil, tokenEndpointError("Claude OAuth token request", status, body, nil)
	}
	return nil, firstErr
}

func postTokenJSON(payload map[string]string) ([]byte, int, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest(http.MethodPost, tokenURL, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return doTokenRequest(req)
}

func postTokenForm(payload map[string]string) ([]byte, int, error) {
	form := url.Values{}
	for key, value := range payload {
		form.Set(key, value)
	}
	req, err := http.NewRequest(http.MethodPost, tokenURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	return doTokenRequest(req)
}

func doTokenRequest(req *http.Request) ([]byte, int, error) {
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("Claude OAuth token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func tokenEndpointError(prefix string, status int, body []byte, err error) error {
	if err != nil {
		return err
	}
	if status == 0 {
		return errors.New(prefix + " failed")
	}
	return fmt.Errorf("%s failed: %d %s", prefix, status, string(body))
}

func tokenUsable(ts TokenState, now time.Time) bool {
	if strings.TrimSpace(ts.AccessToken) == "" {
		return false
	}
	return ts.ExpiresAt == 0 || ts.ExpiresAt > now.Add(30*time.Second).UnixMilli()
}

func generatePKCE() (pkceCodes, error) {
	verifier, err := randomString(64)
	if err != nil {
		return pkceCodes{}, err
	}
	sum := sha256.Sum256([]byte(verifier))
	return pkceCodes{
		Verifier:  verifier,
		Challenge: base64.RawURLEncoding.EncodeToString(sum[:]),
	}, nil
}

func randomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

func parseAuthorizationCode(raw string) string {
	code := strings.TrimSpace(raw)
	if code == "" {
		return ""
	}
	if strings.Contains(code, "?") {
		if parsed, err := url.Parse(code); err == nil {
			if value := parsed.Query().Get("code"); value != "" {
				code = value
			}
		}
	}
	if index := strings.IndexAny(code, "#&"); index >= 0 {
		code = code[:index]
	}
	return strings.TrimSpace(code)
}

func credentialCandidatePaths(configDir string) []string {
	var paths []string
	if dir := strings.TrimSpace(os.Getenv("CLAUDE_CONFIG_DIR")); dir != "" {
		paths = append(paths, filepath.Join(dir, ".credentials.json"))
	}
	if configDir = strings.TrimSpace(configDir); configDir != "" {
		paths = append(paths, quineTokenPath(configDir))
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		paths = append(paths, filepath.Join(home, ".claude", ".credentials.json"))
	}

	seen := make(map[string]bool, len(paths))
	var out []string
	for _, path := range paths {
		clean := filepath.Clean(path)
		if clean == "." || seen[clean] {
			continue
		}
		seen[clean] = true
		out = append(out, clean)
	}
	return out
}

func quineTokenPath(configDir string) string {
	base := strings.TrimSpace(configDir)
	if base == "" {
		base = defaultConfigDir()
	}
	return filepath.Join(base, "claude-oauth.json")
}

func defaultConfigDir() string {
	if dir := strings.TrimSpace(os.Getenv("QUINE_CONFIG_DIR")); dir != "" {
		return dir
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", "quine")
	}
	return "."
}

func loadTokenFile(path string) (TokenState, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TokenState{}, false, nil
		}
		return TokenState{}, false, err
	}

	if ts, ok, err := decodeClaudeCredentials(data); err != nil {
		return TokenState{}, false, err
	} else if ok {
		ts.sourcePath = path
		ts.sourceKind = sourceClaudeCredentials
		return ts, true, nil
	}

	var ts TokenState
	if err := json.Unmarshal(data, &ts); err != nil {
		return TokenState{}, false, err
	}
	if strings.TrimSpace(ts.AccessToken) == "" {
		return TokenState{}, false, errors.New("invalid Claude OAuth token file")
	}
	ts.ClientID = resolveClientID(ts.ClientID)
	ts.sourcePath = path
	ts.sourceKind = sourceQuineToken
	return ts, true, nil
}

func decodeClaudeCredentials(data []byte) (TokenState, bool, error) {
	var raw struct {
		ClaudeAIOAuth *struct {
			AccessToken      string   `json:"accessToken"`
			RefreshToken     string   `json:"refreshToken,omitempty"`
			ExpiresAt        int64    `json:"expiresAt,omitempty"`
			Scopes           []string `json:"scopes,omitempty"`
			SubscriptionType string   `json:"subscriptionType,omitempty"`
			RateLimitTier    string   `json:"rateLimitTier,omitempty"`
			ClientID         string   `json:"clientId,omitempty"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return TokenState{}, false, err
	}
	if raw.ClaudeAIOAuth == nil {
		return TokenState{}, false, nil
	}
	if strings.TrimSpace(raw.ClaudeAIOAuth.AccessToken) == "" {
		return TokenState{}, false, errors.New("invalid Claude credentials file")
	}
	return TokenState{
		AccessToken:      raw.ClaudeAIOAuth.AccessToken,
		RefreshToken:     raw.ClaudeAIOAuth.RefreshToken,
		ExpiresAt:        raw.ClaudeAIOAuth.ExpiresAt,
		Scopes:           raw.ClaudeAIOAuth.Scopes,
		SubscriptionType: raw.ClaudeAIOAuth.SubscriptionType,
		RateLimitTier:    raw.ClaudeAIOAuth.RateLimitTier,
		ClientID:         resolveClientID(raw.ClaudeAIOAuth.ClientID),
	}, true, nil
}

func encodeClaudeCredentials(path string, ts TokenState) ([]byte, error) {
	raw := map[string]any{}
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
	}

	oauth, _ := raw["claudeAiOauth"].(map[string]any)
	if oauth == nil {
		oauth = map[string]any{}
	}
	oauth["accessToken"] = ts.AccessToken
	if ts.RefreshToken != "" {
		oauth["refreshToken"] = ts.RefreshToken
	}
	if ts.ExpiresAt != 0 {
		oauth["expiresAt"] = ts.ExpiresAt
	}
	if len(ts.Scopes) > 0 {
		oauth["scopes"] = ts.Scopes
	}
	if ts.SubscriptionType != "" {
		oauth["subscriptionType"] = ts.SubscriptionType
	}
	if ts.RateLimitTier != "" {
		oauth["rateLimitTier"] = ts.RateLimitTier
	}
	if ts.ClientID != "" && ts.ClientID != defaultClientID {
		oauth["clientId"] = ts.ClientID
	}
	raw["claudeAiOauth"] = oauth

	return json.MarshalIndent(raw, "", "  ")
}

func parseRefreshResponse(data []byte) (TokenState, error) {
	var raw struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int64  `json:"expires_in"`
		Scope            string `json:"scope"`
		SubscriptionType string `json:"subscription_type"`
		RateLimitTier    string `json:"rate_limit_tier"`
		ClientID         string `json:"client_id"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return TokenState{}, err
	}
	if raw.AccessToken == "" {
		return TokenState{}, errors.New("refresh response missing access token")
	}
	if raw.ExpiresIn <= 0 {
		return TokenState{}, errors.New("refresh response missing positive expires_in")
	}
	return TokenState{
		AccessToken:      raw.AccessToken,
		RefreshToken:     raw.RefreshToken,
		ExpiresAt:        time.Now().Add(time.Duration(raw.ExpiresIn) * time.Second).UnixMilli(),
		Scopes:           strings.Fields(raw.Scope),
		SubscriptionType: raw.SubscriptionType,
		RateLimitTier:    raw.RateLimitTier,
		ClientID:         raw.ClientID,
	}, nil
}

func htmlSuccess() string {
	return "<!doctype html><html><head><title>Quine - Claude Authorization Successful</title>" +
		"<style>body{font-family:system-ui,-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#131010;color:#f1ecec}.container{text-align:center;padding:2rem}h1{color:#f1ecec;margin-bottom:1rem}p{color:#b7b1b1}</style>" +
		"</head><body><div class=\"container\"><h1>Authorization Successful</h1><p>You can close this window and return to Quine.</p></div><script>setTimeout(()=>window.close(),2000)</script></body></html>"
}

func htmlError(err string) string {
	return "<!doctype html><html><head><title>Quine - Claude Authorization Failed</title>" +
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
