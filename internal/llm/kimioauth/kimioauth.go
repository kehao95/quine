// Package kimioauth implements OAuth Device Authorization Grant for Kimi API.
// This allows Quine to authenticate as Kimi CLI for API access.
package kimioauth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// Kimi OAuth configuration (matching Kimi CLI)
	clientID       = "17e5f671-d194-4dfb-9706-5516cb48c098"
	defaultAuthURL = "https://auth.kimi.com"
	cliVersion     = "1.13.0"
)

// TokenState holds the OAuth tokens and metadata.
type TokenState struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Scope        string `json:"scope,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

// DeviceInfo holds device identification for Kimi CLI impersonation.
type DeviceInfo struct {
	DeviceID    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	DeviceModel string `json:"device_model"`
	OSVersion   string `json:"os_version"`
}

// DeviceAuthorization holds the response from device authorization request.
type DeviceAuthorization struct {
	UserCode                string `json:"user_code"`
	DeviceCode              string `json:"device_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

func authHost() string {
	if host := os.Getenv("KIMI_OAUTH_HOST"); host != "" {
		return strings.TrimRight(host, "/")
	}
	return defaultAuthURL
}

// CommonHeaders returns the headers that Kimi CLI sends with every request.
func CommonHeaders(device DeviceInfo) map[string]string {
	return map[string]string{
		"User-Agent":         fmt.Sprintf("KimiCLI/%s", cliVersion),
		"X-Msh-Platform":     "kimi_cli",
		"X-Msh-Version":      cliVersion,
		"X-Msh-Device-Name":  asciiHeaderValue(device.DeviceName),
		"X-Msh-Device-Model": asciiHeaderValue(device.DeviceModel),
		"X-Msh-Os-Version":   asciiHeaderValue(device.OSVersion),
		"X-Msh-Device-Id":    device.DeviceID,
	}
}

// asciiHeaderValue sanitizes a string for use in HTTP headers.
func asciiHeaderValue(value string) string {
	var result strings.Builder
	for _, r := range value {
		if r >= 32 && r < 127 {
			result.WriteRune(r)
		}
	}
	s := strings.TrimSpace(result.String())
	if s == "" {
		return "unknown"
	}
	return s
}

// RequestDeviceAuthorization initiates the device authorization flow.
func RequestDeviceAuthorization(device DeviceInfo) (*DeviceAuthorization, error) {
	form := url.Values{}
	form.Set("client_id", clientID)

	req, err := http.NewRequest("POST", authHost()+"/api/oauth/device_authorization", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range CommonHeaders(device) {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device authorization request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device authorization failed: %d %s", resp.StatusCode, string(body))
	}

	var auth DeviceAuthorization
	if err := json.Unmarshal(body, &auth); err != nil {
		return nil, fmt.Errorf("parsing device authorization response: %w", err)
	}

	return &auth, nil
}

// PollForToken polls the token endpoint until the user completes authorization.
func PollForToken(device DeviceInfo, auth *DeviceAuthorization) (TokenState, error) {
	interval := auth.Interval
	if interval < 1 {
		interval = 5
	}

	deadline := time.Now().Add(time.Duration(auth.ExpiresIn) * time.Second)
	if auth.ExpiresIn == 0 {
		deadline = time.Now().Add(5 * time.Minute)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	for time.Now().Before(deadline) {
		form := url.Values{}
		form.Set("client_id", clientID)
		form.Set("device_code", auth.DeviceCode)
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequest("POST", authHost()+"/api/oauth/token", strings.NewReader(form.Encode()))
		if err != nil {
			return TokenState{}, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for k, v := range CommonHeaders(device) {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return parseTokenResponse(body)
		}

		// Check for pending/slow_down errors
		var errResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(body, &errResp); err == nil {
			switch errResp.Error {
			case "authorization_pending", "slow_down":
				// Keep polling
				if errResp.Error == "slow_down" {
					interval += 5
				}
			case "expired_token":
				return TokenState{}, errors.New("device code expired - please restart authorization")
			case "access_denied":
				return TokenState{}, errors.New("authorization denied by user")
			default:
				return TokenState{}, fmt.Errorf("token request failed: %s - %s", errResp.Error, errResp.ErrorDescription)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return TokenState{}, errors.New("authorization timed out")
}

// RefreshToken exchanges a refresh token for a new access token.
func RefreshToken(device DeviceInfo, refreshToken string) (TokenState, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", authHost()+"/api/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return TokenState{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range CommonHeaders(device) {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return TokenState{}, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenState{}, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return TokenState{}, errors.New("refresh token rejected - please re-authenticate")
	}

	if resp.StatusCode != http.StatusOK {
		return TokenState{}, fmt.Errorf("token refresh failed: %d %s", resp.StatusCode, string(body))
	}

	ts, err := parseTokenResponse(body)
	if err != nil {
		return TokenState{}, err
	}
	// A refresh response may legitimately omit refresh_token (RFC 6749 §6: the
	// refresh token is optional and non-rotating servers reuse it). parseTokenResponse
	// would then return an empty RefreshToken, which callers persist over the still-valid
	// one — disabling all future refresh and forcing a full interactive re-login. Carry
	// the existing token forward when the server didn't rotate it.
	if ts.RefreshToken == "" {
		ts.RefreshToken = refreshToken
	}
	return ts, nil
}

func parseTokenResponse(body []byte) (TokenState, error) {
	var raw struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return TokenState{}, fmt.Errorf("parsing token response: %w", err)
	}

	if raw.AccessToken == "" {
		return TokenState{}, errors.New("missing access token in response")
	}

	expiresAt := time.Now().Add(time.Duration(raw.ExpiresIn) * time.Second).UnixMilli()
	if raw.ExpiresIn == 0 {
		expiresAt = time.Now().Add(1 * time.Hour).UnixMilli()
	}

	return TokenState{
		AccessToken:  raw.AccessToken,
		RefreshToken: raw.RefreshToken,
		ExpiresAt:    expiresAt,
		Scope:        raw.Scope,
		TokenType:    raw.TokenType,
	}, nil
}

// --- Token and Device persistence ---

func kimiTokenPath(configDir string) string {
	base := configDir
	if base == "" {
		base = defaultConfigDir()
	}
	return filepath.Join(base, "kimi-oauth.json")
}

func kimiDevicePath(configDir string) string {
	base := configDir
	if base == "" {
		base = defaultConfigDir()
	}
	return filepath.Join(base, "kimi-device.json")
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

// LoadToken loads saved token state from disk.
func LoadToken(configDir string) (TokenState, bool, error) {
	path := kimiTokenPath(configDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TokenState{}, false, nil
		}
		return TokenState{}, false, err
	}
	var ts TokenState
	if err := json.Unmarshal(data, &ts); err != nil {
		return TokenState{}, false, err
	}
	if ts.AccessToken == "" {
		return TokenState{}, false, errors.New("invalid token file")
	}
	return ts, true, nil
}

// SaveToken persists token state to disk.
func SaveToken(configDir string, ts TokenState) error {
	path := kimiTokenPath(configDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadDevice loads saved device info from disk, or generates new if not found.
func LoadDevice(configDir string) (DeviceInfo, error) {
	path := kimiDevicePath(configDir)
	data, err := os.ReadFile(path)
	if err == nil {
		var device DeviceInfo
		if err := json.Unmarshal(data, &device); err == nil && device.DeviceID != "" {
			return device, nil
		}
	}

	// Generate new device info
	device := generateDeviceInfo()
	if err := SaveDevice(configDir, device); err != nil {
		// Non-fatal, we can continue with the generated device
		_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to save device info: %v\n", err)
	}
	return device, nil
}

// SaveDevice persists device info to disk.
func SaveDevice(configDir string, device DeviceInfo) error {
	path := kimiDevicePath(configDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func generateDeviceInfo() DeviceInfo {
	deviceID := generateDeviceID()

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	return DeviceInfo{
		DeviceID:    deviceID,
		DeviceName:  hostname,
		DeviceModel: deviceModel(),
		OSVersion:   osVersion(),
	}
}

func generateDeviceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func deviceModel() string {
	system := runtime.GOOS
	arch := runtime.GOARCH

	switch system {
	case "darwin":
		return fmt.Sprintf("macOS %s", arch)
	case "windows":
		return fmt.Sprintf("Windows %s", arch)
	case "linux":
		return fmt.Sprintf("Linux %s", arch)
	default:
		return fmt.Sprintf("%s %s", system, arch)
	}
}

func osVersion() string {
	// This is a simplified version - Kimi CLI uses platform-specific APIs
	// for more detailed version info, but this is sufficient for impersonation
	return runtime.GOOS + "/" + runtime.GOARCH
}

// EnsureToken ensures a valid access token is available, refreshing or
// re-authenticating as needed.
func EnsureToken(configDir string, openBrowser func(string) error) (TokenState, DeviceInfo, error) {
	device, err := LoadDevice(configDir)
	if err != nil {
		return TokenState{}, DeviceInfo{}, fmt.Errorf("loading device info: %w", err)
	}

	// Try to load existing token
	if ts, ok, err := LoadToken(configDir); err == nil && ok {
		// Token valid for at least 30 more seconds?
		if ts.ExpiresAt > time.Now().Add(30*time.Second).UnixMilli() {
			return ts, device, nil
		}
		// Try to refresh
		if ts.RefreshToken != "" {
			refreshed, err := RefreshToken(device, ts.RefreshToken)
			if err == nil {
				_ = SaveToken(configDir, refreshed)
				return refreshed, device, nil
			}
			// Refresh failed, need to re-auth
			_, _ = fmt.Fprintf(os.Stderr, "Token refresh failed: %v\n", err)
		}
	}

	// Need to do full device authorization flow
	auth, err := RequestDeviceAuthorization(device)
	if err != nil {
		return TokenState{}, device, fmt.Errorf("device authorization: %w", err)
	}

	// Display authorization instructions
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "To authorize Quine with Kimi, please visit:\n")
	fmt.Fprintf(os.Stderr, "  %s\n", auth.VerificationURIComplete)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Or enter code: %s\n", auth.UserCode)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Waiting for authorization...\n")

	// Try to open browser
	if openBrowser != nil {
		if err := openBrowser(auth.VerificationURIComplete); err != nil {
			fmt.Fprintf(os.Stderr, "Could not open browser automatically: %v\n", err)
		}
	}

	// Poll for token
	ts, err := PollForToken(device, auth)
	if err != nil {
		return TokenState{}, device, fmt.Errorf("polling for token: %w", err)
	}

	// Save token
	if err := SaveToken(configDir, ts); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save token: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "Authorization successful!\n\n")
	return ts, device, nil
}

// AuthorizationHeader returns the Authorization header value for API requests.
func AuthorizationHeader(accessToken string) string {
	return "Bearer " + accessToken
}
