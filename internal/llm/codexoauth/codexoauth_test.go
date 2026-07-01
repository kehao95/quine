package codexoauth

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadTokenPrefersCodexCLIAuth(t *testing.T) {
	home := t.TempDir()
	configDir := t.TempDir()
	t.Setenv("HOME", home)

	cliPath := filepath.Join(home, ".codex", "auth.json")
	if err := os.MkdirAll(filepath.Dir(cliPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(cliPath, []byte(`{
  "auth_mode": "chatgpt",
  "OPENAI_API_KEY": null,
  "tokens": {
    "access_token": "`+testJWT(`{"exp":2000000000,"https://api.openai.com/auth":{"chatgpt_account_id":"cli-account"}}`)+`",
    "refresh_token": "cli-refresh",
    "id_token": "`+testJWT(`{"https://api.openai.com/auth":{"chatgpt_account_id":"cli-account"}}`)+`"
  }
}`), 0o600); err != nil {
		t.Fatalf("WriteFile(cli auth) error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(configDir, "codex-oauth.json"), []byte(`{
  "access_token": "legacy-access",
  "refresh_token": "legacy-refresh",
  "expires_at": 123
}`), 0o600); err != nil {
		t.Fatalf("WriteFile(legacy token) error = %v", err)
	}

	got, ok, err := LoadToken(configDir)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}
	if !ok {
		t.Fatal("LoadToken() ok = false, want true")
	}
	if got.RefreshToken != "cli-refresh" {
		t.Fatalf("RefreshToken = %q, want %q", got.RefreshToken, "cli-refresh")
	}
	if got.AccountID != "cli-account" {
		t.Fatalf("AccountID = %q, want %q", got.AccountID, "cli-account")
	}
}

func TestLoadTokenMigratesLegacyConfigToCodexCLIShape(t *testing.T) {
	home := t.TempDir()
	configDir := t.TempDir()
	t.Setenv("HOME", home)

	legacyAccess := testJWT(`{"exp":2000000000,"https://api.openai.com/auth":{"chatgpt_account_id":"legacy-account"}}`)
	legacyID := testJWT(`{"https://api.openai.com/auth":{"chatgpt_account_id":"legacy-account"}}`)
	path := filepath.Join(configDir, "codex-oauth.json")
	if err := os.WriteFile(path, []byte(`{
  "access_token": "`+legacyAccess+`",
  "refresh_token": "legacy-refresh",
  "expires_at": 123,
  "account_id": "legacy-account",
  "id_token": "`+legacyID+`"
}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, ok, err := LoadToken(configDir)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}
	if !ok {
		t.Fatal("LoadToken() ok = false, want true")
	}
	if got.RefreshToken != "legacy-refresh" {
		t.Fatalf("RefreshToken = %q, want %q", got.RefreshToken, "legacy-refresh")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if strings.Contains(string(data), `"expires_at"`) {
		t.Fatalf("migrated config still contains legacy flat shape: %s", string(data))
	}

	var raw codexCLIAuthFile
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if raw.AuthMode != "chatgpt" {
		t.Fatalf("auth_mode = %q, want %q", raw.AuthMode, "chatgpt")
	}
	if raw.Tokens.RefreshToken != "legacy-refresh" {
		t.Fatalf("tokens.refresh_token = %q, want %q", raw.Tokens.RefreshToken, "legacy-refresh")
	}
	if raw.Tokens.AccountID != "legacy-account" {
		t.Fatalf("tokens.account_id = %q, want %q", raw.Tokens.AccountID, "legacy-account")
	}
}

func testJWT(payload string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	body := base64.RawURLEncoding.EncodeToString([]byte(payload))
	return header + "." + body + ".sig"
}
