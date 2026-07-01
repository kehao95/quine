package copilotoauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultClientID       = "Iv1.b507a08c87ecfe98"
	defaultDeviceCodeURL  = "https://github.com/login/device/code"
	defaultAccessTokenURL = "https://github.com/login/oauth/access_token"
	defaultGitHubAPIURL   = "https://api.github.com"
	tokenFilename         = "copilot-oauth.json"
)

type TokenState struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

type DeviceAuthorization struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type Profile struct {
	Login string `json:"login"`
	Name  string `json:"name,omitempty"`
	ID    int64  `json:"id,omitempty"`
}

func DeviceCodeURL() string {
	if v := strings.TrimSpace(os.Getenv("GITHUB_DEVICE_CODE_URL")); v != "" {
		return v
	}
	return defaultDeviceCodeURL
}

func AccessTokenURL() string {
	if v := strings.TrimSpace(os.Getenv("GITHUB_ACCESS_TOKEN_URL")); v != "" {
		return v
	}
	return defaultAccessTokenURL
}

func GitHubAPIURL() string {
	if v := strings.TrimSpace(os.Getenv("GITHUB_API_URL")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return defaultGitHubAPIURL
}

func ResolveClientID(explicit string) string {
	for _, candidate := range []string{
		strings.TrimSpace(explicit),
		strings.TrimSpace(os.Getenv("COPILOT_OAUTH_CLIENT_ID")),
		strings.TrimSpace(os.Getenv("GITHUB_COPILOT_OAUTH_CLIENT_ID")),
		strings.TrimSpace(os.Getenv("GITHUB_OAUTH_CLIENT_ID")),
		defaultClientID,
	} {
		if candidate != "" {
			return candidate
		}
	}
	return ""
}

func LoadToken(configDir string) (TokenState, bool, error) {
	path := tokenPath(configDir)
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
	if strings.TrimSpace(ts.AccessToken) == "" {
		return TokenState{}, false, errors.New("invalid GitHub Copilot token file")
	}
	return ts, true, nil
}

func SaveToken(configDir string, ts TokenState) error {
	path := tokenPath(configDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func EnsureToken(configDir, explicitClientID string, openBrowser func(string) error) (TokenState, error) {
	if ts, ok, err := LoadToken(configDir); err == nil && ok {
		return ts, nil
	}

	clientID := ResolveClientID(explicitClientID)
	auth, err := RequestDeviceAuthorization(clientID)
	if err != nil {
		return TokenState{}, err
	}

	fmt.Fprintf(os.Stderr, "\nTo authorize Quine with GitHub Copilot, visit:\n  %s\n\n", auth.VerificationURI)
	fmt.Fprintf(os.Stderr, "Then enter code: %s\n\n", auth.UserCode)
	fmt.Fprintf(os.Stderr, "Waiting for authorization...\n")

	if openBrowser != nil {
		if err := openBrowser(auth.VerificationURI); err != nil {
			fmt.Fprintf(os.Stderr, "Could not open browser automatically: %v\n", err)
		}
	}

	ts, err := PollForToken(clientID, auth)
	if err != nil {
		return TokenState{}, err
	}
	if err := SaveToken(configDir, ts); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save GitHub Copilot token: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "Authorization successful!\n\n")
	return ts, nil
}

func RequestDeviceAuthorization(clientID string) (*DeviceAuthorization, error) {
	form := url.Values{}
	form.Set("client_id", clientID)

	req, err := http.NewRequest(http.MethodPost, DeviceCodeURL(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "quine-github-copilot-auth")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device authorization request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device authorization failed: %d %s", resp.StatusCode, string(body))
	}

	var auth DeviceAuthorization
	if err := json.Unmarshal(body, &auth); err != nil {
		return nil, fmt.Errorf("parsing device authorization response: %w", err)
	}
	if auth.DeviceCode == "" || auth.UserCode == "" || auth.VerificationURI == "" {
		return nil, errors.New("device authorization response missing required fields")
	}
	return &auth, nil
}

func PollForToken(clientID string, auth *DeviceAuthorization) (TokenState, error) {
	interval := auth.Interval
	if interval < 1 {
		interval = 5
	}

	deadline := time.Now().Add(time.Duration(auth.ExpiresIn) * time.Second)
	if auth.ExpiresIn <= 0 {
		deadline = time.Now().Add(15 * time.Minute)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	for time.Now().Before(deadline) {
		form := url.Values{}
		form.Set("client_id", clientID)
		form.Set("device_code", auth.DeviceCode)
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequest(http.MethodPost, AccessTokenURL(), strings.NewReader(form.Encode()))
		if err != nil {
			return TokenState{}, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", "quine-github-copilot-auth")

		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var tokenResp struct {
			AccessToken      string `json:"access_token"`
			TokenType        string `json:"token_type,omitempty"`
			Scope            string `json:"scope,omitempty"`
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			if resp.StatusCode != http.StatusOK {
				return TokenState{}, fmt.Errorf("token request failed: %d %s", resp.StatusCode, string(body))
			}
			return TokenState{}, fmt.Errorf("parsing token response: %w", err)
		}

		if strings.TrimSpace(tokenResp.AccessToken) != "" {
			return TokenState{
				AccessToken: tokenResp.AccessToken,
				TokenType:   tokenResp.TokenType,
				Scope:       tokenResp.Scope,
			}, nil
		}

		if tokenResp.Error != "" {
			switch tokenResp.Error {
			case "authorization_pending":
			case "slow_down":
				interval += 5
			case "expired_token":
				return TokenState{}, errors.New("device code expired; restart login")
			case "access_denied":
				return TokenState{}, errors.New("authorization denied by user")
			default:
				return TokenState{}, fmt.Errorf("token request failed: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
			}
		} else if resp.StatusCode != http.StatusOK {
			return TokenState{}, fmt.Errorf("token request failed: %d %s", resp.StatusCode, string(body))
		} else {
			return TokenState{}, errors.New("missing access token in response")
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return TokenState{}, errors.New("authorization timed out")
}

func FetchProfile(accessToken string) (Profile, error) {
	req, err := http.NewRequest(http.MethodGet, GitHubAPIURL()+"/user", nil)
	if err != nil {
		return Profile{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", AuthorizationHeader(accessToken))
	req.Header.Set("User-Agent", "quine-github-copilot-auth")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("fetching GitHub profile: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf("fetching GitHub profile failed: %d %s", resp.StatusCode, string(body))
	}

	var profile Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		return Profile{}, fmt.Errorf("parsing GitHub profile: %w", err)
	}
	if strings.TrimSpace(profile.Login) == "" {
		return Profile{}, errors.New("GitHub profile response missing login")
	}
	return profile, nil
}

func AuthorizationHeader(accessToken string) string {
	return "Bearer " + accessToken
}

func tokenPath(configDir string) string {
	base := strings.TrimSpace(configDir)
	if base == "" {
		base = defaultConfigDir()
	}
	return filepath.Join(base, tokenFilename)
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
