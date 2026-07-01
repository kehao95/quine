package transport

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/kimioauth"
)

// KimiOAuthTransport implements Transport using Kimi OAuth Device Authorization.
// It impersonates Kimi CLI by sending the same headers and using OAuth tokens.
type KimiOAuthTransport struct {
	configDir string
	client    *http.Client
	mu        sync.Mutex
	state     kimiOAuthState
	device    kimioauth.DeviceInfo
}

type kimiOAuthState struct {
	accessToken  string
	refreshToken string
	expiresAt    int64
}

// NewKimiOAuthTransport creates a new KimiOAuthTransport.
func NewKimiOAuthTransport(_ *config.Config) (Transport, error) {
	configDir := os.Getenv("QUINE_CONFIG_DIR")
	return &KimiOAuthTransport{
		configDir: configDir,
		client:    &http.Client{Timeout: 2 * time.Minute},
	}, nil
}

// Sign adds Kimi OAuth authentication and impersonation headers to the request.
func (t *KimiOAuthTransport) Sign(req *http.Request, body []byte) error {
	ts, device, err := t.ensureToken()
	if err != nil {
		return err
	}

	// Set Authorization header
	req.Header.Set("Authorization", kimioauth.AuthorizationHeader(ts.accessToken))

	// Set Kimi CLI impersonation headers
	for k, v := range kimioauth.CommonHeaders(device) {
		req.Header.Set(k, v)
	}

	return nil
}

// persistRefreshedToken writes a freshly refreshed token to disk, surfacing a
// save failure to stderr instead of discarding it. A swallowed save error lets
// durable token state silently diverge from the in-memory state, so a later
// session fails to refresh with no breadcrumb explaining why.
func (t *KimiOAuthTransport) persistRefreshedToken(refreshed kimioauth.TokenState) {
	if err := kimioauth.SaveToken(t.configDir, refreshed); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to persist refreshed Kimi token: %v\n", err)
	}
}

func (t *KimiOAuthTransport) ensureToken() (kimiOAuthState, kimioauth.DeviceInfo, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UnixMilli()

	// Return cached token if still valid
	if t.state.accessToken != "" && t.state.expiresAt > now+30*1000 {
		return t.state, t.device, nil
	}

	// Try to refresh if we have a refresh token
	if t.state.refreshToken != "" {
		device := t.device
		if device.DeviceID == "" {
			device, _ = kimioauth.LoadDevice(t.configDir)
		}
		refreshed, err := kimioauth.RefreshToken(device, t.state.refreshToken)
		if err == nil {
			t.state = kimiOAuthStateFromToken(refreshed)
			t.device = device
			t.persistRefreshedToken(refreshed)
			return t.state, t.device, nil
		}
	}

	// Try to load from disk
	loaded, ok, err := kimioauth.LoadToken(t.configDir)
	if err == nil && ok {
		device, _ := kimioauth.LoadDevice(t.configDir)
		if loaded.ExpiresAt > now+30*1000 {
			t.state = kimiOAuthStateFromToken(loaded)
			t.device = device
			return t.state, t.device, nil
		}
		// Try to refresh the loaded token
		if loaded.RefreshToken != "" {
			refreshed, err := kimioauth.RefreshToken(device, loaded.RefreshToken)
			if err == nil {
				t.state = kimiOAuthStateFromToken(refreshed)
				t.device = device
				t.persistRefreshedToken(refreshed)
				return t.state, t.device, nil
			}
		}
	}

	// Need to do full OAuth flow
	newToken, device, err := kimioauth.EnsureToken(t.configDir, kimioauth.OpenBrowser)
	if err != nil {
		return kimiOAuthState{}, kimioauth.DeviceInfo{}, err
	}

	t.state = kimiOAuthStateFromToken(newToken)
	t.device = device
	return t.state, t.device, nil
}

func kimiOAuthStateFromToken(ts kimioauth.TokenState) kimiOAuthState {
	return kimiOAuthState{
		accessToken:  ts.AccessToken,
		refreshToken: ts.RefreshToken,
		expiresAt:    ts.ExpiresAt,
	}
}
