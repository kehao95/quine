package kimioauth

import (
	"testing"
)

func TestCommonHeaders(t *testing.T) {
	device := DeviceInfo{
		DeviceID:    "abc123",
		DeviceName:  "test-host",
		DeviceModel: "macOS arm64",
		OSVersion:   "darwin/arm64",
	}

	headers := CommonHeaders(device)

	// Check required headers
	if headers["User-Agent"] != "KimiCLI/"+cliVersion {
		t.Errorf("User-Agent = %q, want %q", headers["User-Agent"], "KimiCLI/"+cliVersion)
	}
	if headers["X-Msh-Platform"] != "kimi_cli" {
		t.Errorf("X-Msh-Platform = %q, want %q", headers["X-Msh-Platform"], "kimi_cli")
	}
	if headers["X-Msh-Version"] != cliVersion {
		t.Errorf("X-Msh-Version = %q, want %q", headers["X-Msh-Version"], cliVersion)
	}
	if headers["X-Msh-Device-Id"] != "abc123" {
		t.Errorf("X-Msh-Device-Id = %q, want %q", headers["X-Msh-Device-Id"], "abc123")
	}
	if headers["X-Msh-Device-Name"] != "test-host" {
		t.Errorf("X-Msh-Device-Name = %q, want %q", headers["X-Msh-Device-Name"], "test-host")
	}
}

func TestAsciiHeaderValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with spaces", "with spaces"},
		{"émoji 🎉", "moji"},
		{"中文", "unknown"},
		{"", "unknown"},
		{"  trimmed  ", "trimmed"},
	}

	for _, tt := range tests {
		got := asciiHeaderValue(tt.input)
		if got != tt.want {
			t.Errorf("asciiHeaderValue(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGenerateDeviceID(t *testing.T) {
	id1 := generateDeviceID()
	id2 := generateDeviceID()

	if id1 == "" {
		t.Error("generateDeviceID() returned empty string")
	}
	if len(id1) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("generateDeviceID() returned %d chars, want 32", len(id1))
	}
	if id1 == id2 {
		t.Error("generateDeviceID() returned same ID twice")
	}
}

func TestDeviceModel(t *testing.T) {
	model := deviceModel()
	if model == "" {
		t.Error("deviceModel() returned empty string")
	}
	// Should contain OS and architecture
	t.Logf("deviceModel() = %q", model)
}

func TestAuthorizationHeader(t *testing.T) {
	got := AuthorizationHeader("test-token")
	want := "Bearer test-token"
	if got != want {
		t.Errorf("AuthorizationHeader() = %q, want %q", got, want)
	}
}
