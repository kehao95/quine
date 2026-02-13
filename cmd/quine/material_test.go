package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

// TestHandleStdin_TextMode tests default text/streaming mode
func TestHandleStdin_TextMode(t *testing.T) {
	// Create a temp dir for test
	tmpDir := t.TempDir()
	cfg := &config.Config{
		DataDir:   tmpDir,
		SessionID: "test-text-mode",
	}

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Write text data
	testText := "Hello, World!\nThis is a test."
	go func() {
		w.Write([]byte(testText))
		w.Close()
	}()

	// Replace stdin temporarily
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Call handleStdin with TEXT mode (default)
	material, err := handleStdin(cfg, stdinModeText)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Should indicate stdin is piped
	if !strings.Contains(material, "Input is piped to stdin") {
		t.Errorf("handleStdin() material = %q, want stdin piped message", material)
	}

	// Verify no binary file was created
	files, _ := filepath.Glob(filepath.Join(tmpDir, "stdin-*.bin"))
	if len(files) > 0 {
		t.Errorf("unexpected binary file created: %v", files)
	}
}

// TestHandleStdin_BinaryMode tests -b flag forcing binary mode
func TestHandleStdin_BinaryMode(t *testing.T) {
	// Create a temp dir for test
	tmpDir := t.TempDir()
	cfg := &config.Config{
		DataDir:   tmpDir,
		SessionID: "test-binary-mode",
	}

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Write binary data
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	go func() {
		w.Write(binaryData)
		w.Close()
	}()

	// Replace stdin temporarily
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Call handleStdin with BINARY mode (-b flag)
	material, err := handleStdin(cfg, stdinModeBinary)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Should return a reference to the binary file
	expectedPath := filepath.Join(tmpDir, "stdin-test-binary-mode.bin")
	expectedMsg := "User sent a binary file at " + expectedPath
	if material != expectedMsg {
		t.Errorf("handleStdin() material = %q, want %q", material, expectedMsg)
	}

	// Verify binary file was created with correct content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read binary file: %v", err)
	}
	if string(content) != string(binaryData) {
		t.Errorf("binary file content = %v, want %v", content, binaryData)
	}
}

// TestHandleStdin_BinaryMode_Empty tests -b flag with empty stdin
func TestHandleStdin_BinaryMode_Empty(t *testing.T) {
	// Create a temp dir for test
	tmpDir := t.TempDir()
	cfg := &config.Config{
		DataDir:   tmpDir,
		SessionID: "test-binary-empty",
	}

	// Create a pipe to simulate empty stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Close immediately (empty stdin)
	w.Close()

	// Replace stdin temporarily
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Call handleStdin with BINARY mode
	material, err := handleStdin(cfg, stdinModeBinary)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Empty stdin should return "Begin." even with binary mode
	if material != "Begin." {
		t.Errorf("handleStdin() material = %q, want %q", material, "Begin.")
	}
}
