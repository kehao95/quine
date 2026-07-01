package kimioauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRefreshTokenCarriesForwardRefreshToken pins H1: a non-rotating server that
// returns a fresh access token but omits refresh_token must not wipe the
// still-valid refresh token. Carrying it forward keeps future refresh working
// instead of forcing a full interactive re-login.
func TestRefreshTokenCarriesForwardRefreshToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"new-access","expires_in":3600}`))
	}))
	defer srv.Close()
	t.Setenv("KIMI_OAUTH_HOST", srv.URL)

	ts, err := RefreshToken(DeviceInfo{}, "existing-refresh")
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if ts.AccessToken != "new-access" {
		t.Fatalf("access token = %q, want new-access", ts.AccessToken)
	}
	if ts.RefreshToken != "existing-refresh" {
		t.Fatalf("refresh token must be carried forward, got %q", ts.RefreshToken)
	}
}

// TestRefreshTokenKeepsRotatedToken ensures a rotating server still wins: a
// returned refresh_token replaces the old one.
func TestRefreshTokenKeepsRotatedToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"rotated","expires_in":3600}`))
	}))
	defer srv.Close()
	t.Setenv("KIMI_OAUTH_HOST", srv.URL)

	ts, err := RefreshToken(DeviceInfo{}, "existing-refresh")
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if ts.RefreshToken != "rotated" {
		t.Fatalf("rotated refresh token must win, got %q", ts.RefreshToken)
	}
}
