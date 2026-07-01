package copilotoauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadSaveToken(t *testing.T) {
	dir := t.TempDir()
	want := TokenState{
		AccessToken: "gho_test",
		TokenType:   "bearer",
		Scope:       "read:user",
	}
	if err := SaveToken(dir, want); err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	got, ok, err := LoadToken(dir)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}
	if !ok {
		t.Fatal("LoadToken() ok = false, want true")
	}
	if got != want {
		t.Fatalf("LoadToken() = %+v, want %+v", got, want)
	}
}

func TestResolveClientID(t *testing.T) {
	t.Setenv("COPILOT_OAUTH_CLIENT_ID", "")
	t.Setenv("GITHUB_COPILOT_OAUTH_CLIENT_ID", "")
	t.Setenv("GITHUB_OAUTH_CLIENT_ID", "")
	if got := ResolveClientID("explicit"); got != "explicit" {
		t.Fatalf("ResolveClientID(explicit) = %q, want %q", got, "explicit")
	}

	t.Setenv("COPILOT_OAUTH_CLIENT_ID", "copilot")
	if got := ResolveClientID(""); got != "copilot" {
		t.Fatalf("ResolveClientID() = %q, want %q", got, "copilot")
	}

	t.Setenv("COPILOT_OAUTH_CLIENT_ID", "")
	t.Setenv("GITHUB_OAUTH_CLIENT_ID", "generic")
	if got := ResolveClientID(""); got != "generic" {
		t.Fatalf("ResolveClientID() = %q, want %q", got, "generic")
	}

	t.Setenv("GITHUB_OAUTH_CLIENT_ID", "")
	if got := ResolveClientID(""); got != defaultClientID {
		t.Fatalf("ResolveClientID() = %q, want %q", got, defaultClientID)
	}
}

func TestPollForTokenHandlesPendingThenSuccess(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.Method; got != http.MethodPost {
			t.Fatalf("method = %q, want %q", got, http.MethodPost)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("Accept = %q, want %q", got, "application/json")
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("Content-Type = %q, want %q", got, "application/x-www-form-urlencoded")
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if got := r.PostForm.Get("client_id"); got != defaultClientID {
			t.Fatalf("client_id = %q, want %q", got, defaultClientID)
		}
		if got := r.PostForm.Get("device_code"); got != "device-code" {
			t.Fatalf("device_code = %q, want %q", got, "device-code")
		}
		if got := r.PostForm.Get("grant_type"); got != "urn:ietf:params:oauth:grant-type:device_code" {
			t.Fatalf("grant_type = %q, want device code grant", got)
		}

		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			_, _ = w.Write([]byte(`{"error":"authorization_pending","error_description":"waiting for user"}`))
			return
		}
		_, _ = w.Write([]byte(`{"access_token":"gho_test","token_type":"bearer","scope":"read:user"}`))
	}))
	defer server.Close()

	t.Setenv("GITHUB_ACCESS_TOKEN_URL", server.URL)

	got, err := PollForToken(defaultClientID, &DeviceAuthorization{
		DeviceCode: "device-code",
		ExpiresIn:  10,
		Interval:   1,
	})
	if err != nil {
		t.Fatalf("PollForToken() error = %v", err)
	}
	want := TokenState{
		AccessToken: "gho_test",
		TokenType:   "bearer",
		Scope:       "read:user",
	}
	if got != want {
		t.Fatalf("PollForToken() = %+v, want %+v", got, want)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
}
