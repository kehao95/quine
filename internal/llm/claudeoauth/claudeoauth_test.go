package claudeoauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadTokenPrefersEnvToken(t *testing.T) {
	t.Setenv("CLAUDE_CODE_OAUTH_TOKEN", "env-access")
	t.Setenv("CLAUDE_CODE_OAUTH_SCOPES", "user:inference user:profile")
	t.Setenv("CLAUDE_CODE_SUBSCRIPTION_TYPE", "max")
	t.Setenv("CLAUDE_CODE_RATE_LIMIT_TIER", "default_claude_max_20x")

	ts, ok, err := LoadToken(t.TempDir())
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if !ok {
		t.Fatal("LoadToken() ok = false, want true")
	}
	if ts.AccessToken != "env-access" {
		t.Fatalf("AccessToken = %q, want env-access", ts.AccessToken)
	}
	if got, want := ts.Scopes, []string{"user:inference", "user:profile"}; !sameStrings(got, want) {
		t.Fatalf("Scopes = %#v, want %#v", got, want)
	}
	if ts.SubscriptionType != "max" || ts.RateLimitTier != "default_claude_max_20x" {
		t.Fatalf("subscription metadata = %q/%q", ts.SubscriptionType, ts.RateLimitTier)
	}
}

func TestLoadTokenFromClaudeCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)

	writeClaudeCredentials(t, dir, map[string]any{
		"accessToken":      "stored-access",
		"refreshToken":     "stored-refresh",
		"expiresAt":        float64(1893456000000),
		"scopes":           []any{"user:inference", "user:mcp_servers"},
		"subscriptionType": "pro",
		"rateLimitTier":    "default_claude_pro",
	})

	ts, ok, err := LoadToken("")
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if !ok {
		t.Fatal("LoadToken() ok = false, want true")
	}
	if ts.AccessToken != "stored-access" || ts.RefreshToken != "stored-refresh" {
		t.Fatalf("tokens = %q/%q", ts.AccessToken, ts.RefreshToken)
	}
	if ts.ExpiresAt != 1893456000000 {
		t.Fatalf("ExpiresAt = %d, want 1893456000000", ts.ExpiresAt)
	}
	if got, want := ts.Scopes, []string{"user:inference", "user:mcp_servers"}; !sameStrings(got, want) {
		t.Fatalf("Scopes = %#v, want %#v", got, want)
	}
}

func TestSaveTokenUpdatesClaudeCredentialsAndPreservesOtherFields(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	writeClaudeCredentials(t, dir, map[string]any{
		"accessToken":  "old-access",
		"refreshToken": "old-refresh",
		"expiresAt":    float64(1000),
	})

	ts, ok, err := LoadToken("")
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if !ok {
		t.Fatal("LoadToken() ok = false, want true")
	}
	ts.AccessToken = "new-access"
	ts.RefreshToken = "new-refresh"
	ts.ExpiresAt = 2000
	ts.Scopes = []string{"user:inference"}
	ts.SubscriptionType = "max"

	if err := SaveToken("", ts); err != nil {
		t.Fatalf("SaveToken() error: %v", err)
	}

	var raw map[string]any
	readJSON(t, filepath.Join(dir, ".credentials.json"), &raw)
	if raw["other"] != "preserved" {
		t.Fatalf("other field not preserved: %#v", raw)
	}
	oauth := raw["claudeAiOauth"].(map[string]any)
	if oauth["accessToken"] != "new-access" || oauth["refreshToken"] != "new-refresh" {
		t.Fatalf("saved tokens = %#v", oauth)
	}
	if oauth["subscriptionType"] != "max" {
		t.Fatalf("subscriptionType = %#v", oauth["subscriptionType"])
	}
}

func TestRefreshTokenPostsFormEncodedOAuthRequestAndInheritsMetadata(t *testing.T) {
	oldTokenURL := tokenURL
	defer func() { tokenURL = oldTokenURL }()

	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("Accept = %q, want application/json", got)
		}

		var payload url.Values
		switch requests {
		case 1:
			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("first Content-Type = %q, want application/json", got)
			}
			var raw map[string]string
			if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
				t.Fatalf("decode JSON body: %v", err)
			}
			payload = url.Values{}
			for key, value := range raw {
				payload.Set(key, value)
			}
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		case 2:
			if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
				t.Fatalf("second Content-Type = %q, want application/x-www-form-urlencoded", got)
			}
			if err := r.ParseForm(); err != nil {
				t.Fatalf("ParseForm(): %v", err)
			}
			payload = url.Values(r.PostForm)
		default:
			t.Fatalf("unexpected request count %d", requests)
		}

		if got := payload.Get("grant_type"); got != "refresh_token" {
			t.Fatalf("grant_type = %q", got)
		}
		if got := payload.Get("refresh_token"); got != "old-refresh" {
			t.Fatalf("refresh_token = %q", got)
		}
		if got := payload.Get("client_id"); got != defaultClientID {
			t.Fatalf("client_id = %q", got)
		}
		if got := payload.Get("scope"); got != "user:inference" {
			t.Fatalf("scope = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "new-access",
			"expires_in":   3600,
		})
	}))
	defer srv.Close()
	tokenURL = srv.URL

	prior := TokenState{
		RefreshToken:     "old-refresh",
		Scopes:           []string{"user:inference"},
		SubscriptionType: "max",
		RateLimitTier:    "default_claude_max_20x",
	}
	got, err := RefreshToken("old-refresh", prior)
	if err != nil {
		t.Fatalf("RefreshToken() error: %v", err)
	}
	if got.AccessToken != "new-access" || got.RefreshToken != "old-refresh" {
		t.Fatalf("tokens = %q/%q", got.AccessToken, got.RefreshToken)
	}
	if got.ExpiresAt <= time.Now().UnixMilli() {
		t.Fatalf("ExpiresAt = %d, want future", got.ExpiresAt)
	}
	if got.SubscriptionType != "max" || got.RateLimitTier != "default_claude_max_20x" {
		t.Fatalf("metadata = %q/%q", got.SubscriptionType, got.RateLimitTier)
	}
	if got.ClientID != defaultClientID {
		t.Fatalf("ClientID = %q, want %q", got.ClientID, defaultClientID)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want JSON attempt plus form fallback", requests)
	}
}

func TestEnsureTokenRunsBrowserLoginWhenClaudeCredentialsExpiredWithoutRefresh(t *testing.T) {
	oldTokenURL := tokenURL
	defer func() { tokenURL = oldTokenURL }()

	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	writeClaudeCredentials(t, dir, map[string]any{
		"accessToken":  "old-access",
		"refreshToken": "",
		"expiresAt":    float64(1000),
		"scopes":       []any{"user:inference", "user:profile"},
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode token exchange body: %v", err)
		}
		if got := payload["grant_type"]; got != "authorization_code" {
			t.Fatalf("grant_type = %q", got)
		}
		if got := payload["code"]; got != "test-code" {
			t.Fatalf("code = %q", got)
		}
		if payload["code_verifier"] == "" {
			t.Fatal("code_verifier is empty")
		}
		if payload["redirect_uri"] == "" {
			t.Fatal("redirect_uri is empty")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "login-access",
			"refresh_token": "login-refresh",
			"expires_in":    3600,
			"scope":         "user:inference user:profile",
		})
	}))
	defer srv.Close()
	tokenURL = srv.URL

	openBrowser := func(authURL string) error {
		parsed, err := url.Parse(authURL)
		if err != nil {
			t.Fatalf("parse auth URL: %v", err)
		}
		query := parsed.Query()
		if query.Get("code") != "true" {
			t.Fatalf("code param = %q", query.Get("code"))
		}
		if query.Get("client_id") != defaultClientID {
			t.Fatalf("client_id = %q", query.Get("client_id"))
		}
		if query.Get("code_challenge") == "" {
			t.Fatal("code_challenge is empty")
		}
		redirectURI := query.Get("redirect_uri")
		state := query.Get("state")
		if redirectURI == "" || state == "" {
			t.Fatalf("redirect_uri/state missing in %s", authURL)
		}
		go func() {
			callbackURL := redirectURI + "?code=test-code&state=" + url.QueryEscape(state)
			resp, err := http.Get(callbackURL)
			if err == nil {
				_ = resp.Body.Close()
			}
		}()
		return nil
	}

	got, err := EnsureToken("", openBrowser)
	if err != nil {
		t.Fatalf("EnsureToken() error: %v", err)
	}
	if got.AccessToken != "login-access" || got.RefreshToken != "login-refresh" {
		t.Fatalf("tokens = %q/%q", got.AccessToken, got.RefreshToken)
	}

	var raw map[string]any
	readJSON(t, filepath.Join(dir, ".credentials.json"), &raw)
	oauth := raw["claudeAiOauth"].(map[string]any)
	if oauth["accessToken"] != "login-access" || oauth["refreshToken"] != "login-refresh" {
		t.Fatalf("stored oauth = %#v", oauth)
	}
}

func writeClaudeCredentials(t *testing.T, dir string, oauth map[string]any) {
	t.Helper()
	data, err := json.MarshalIndent(map[string]any{
		"claudeAiOauth": oauth,
		"other":         "preserved",
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".credentials.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func readJSON(t *testing.T, path string, dst any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		t.Fatal(err)
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
