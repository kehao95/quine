package main

import (
	"io"
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
	material, stdinReader, err := handleStdin(cfg, stdinModeText)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Should indicate streaming mode
	if !strings.Contains(material, "Streaming input available") {
		t.Errorf("handleStdin() material = %q, want streaming message", material)
	}

	// Should return a reader
	if stdinReader == nil {
		t.Fatal("handleStdin() stdinReader = nil, want non-nil for text mode")
	}

	// Reading from the reader should yield the original text
	content, err := io.ReadAll(stdinReader)
	if err != nil {
		t.Fatalf("failed to read from stdinReader: %v", err)
	}
	if string(content) != testText {
		t.Errorf("stdinReader content = %q, want %q", string(content), testText)
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
	material, stdinReader, err := handleStdin(cfg, stdinModeBinary)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Should return a reference to the binary file
	expectedPath := filepath.Join(tmpDir, "stdin-test-binary-mode.bin")
	expectedMsg := "User sent a binary file at " + expectedPath
	if material != expectedMsg {
		t.Errorf("handleStdin() material = %q, want %q", material, expectedMsg)
	}

	// Reader should be nil for binary (no streaming)
	if stdinReader != nil {
		t.Errorf("handleStdin() stdinReader = %v, want nil for binary mode", stdinReader)
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
	material, stdinReader, err := handleStdin(cfg, stdinModeBinary)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Empty stdin should return "Begin." even with binary mode
	if material != "Begin." {
		t.Errorf("handleStdin() material = %q, want %q", material, "Begin.")
	}

	if stdinReader != nil {
		t.Errorf("handleStdin() stdinReader = %v, want nil for empty stdin", stdinReader)
	}
}

// TestHandleStdin_TextMode_Unicode tests default text mode with Unicode content
func TestHandleStdin_TextMode_Unicode(t *testing.T) {
	// Create a temp dir for test
	tmpDir := t.TempDir()
	cfg := &config.Config{
		DataDir:   tmpDir,
		SessionID: "test-unicode",
	}

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Write Unicode text
	testText := "你好世界 🌍 مرحبا"
	go func() {
		w.Write([]byte(testText))
		w.Close()
	}()

	// Replace stdin temporarily
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Call handleStdin with TEXT mode (default)
	material, stdinReader, err := handleStdin(cfg, stdinModeText)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Should indicate streaming mode
	if !strings.Contains(material, "Streaming input available") {
		t.Errorf("handleStdin() material = %q, want streaming message", material)
	}

	// Reading from the reader should yield the original text
	if stdinReader == nil {
		t.Fatal("handleStdin() stdinReader = nil, want non-nil")
	}
	content, err := io.ReadAll(stdinReader)
	if err != nil {
		t.Fatalf("failed to read from stdinReader: %v", err)
	}
	if string(content) != testText {
		t.Errorf("stdinReader content = %q, want %q", string(content), testText)
	}
}

// TestHandleStdin_BinaryMode_TextContent tests -b flag treats text as binary
func TestHandleStdin_BinaryMode_TextContent(t *testing.T) {
	// Create a temp dir for test
	tmpDir := t.TempDir()
	cfg := &config.Config{
		DataDir:   tmpDir,
		SessionID: "test-binary-text",
	}

	// Create a pipe to simulate stdin with text content
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Write content that is actually text
	textData := []byte("Hello, World!\nThis is text.")
	go func() {
		w.Write(textData)
		w.Close()
	}()

	// Replace stdin temporarily
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Call handleStdin with BINARY mode - should save as binary anyway
	material, stdinReader, err := handleStdin(cfg, stdinModeBinary)
	if err != nil {
		t.Fatalf("handleStdin() error = %v", err)
	}

	// Should return a reference to the binary file
	expectedPath := filepath.Join(tmpDir, "stdin-test-binary-text.bin")
	expectedMsg := "User sent a binary file at " + expectedPath
	if material != expectedMsg {
		t.Errorf("handleStdin() material = %q, want %q", material, expectedMsg)
	}

	// Reader should be nil for binary (no streaming)
	if stdinReader != nil {
		t.Errorf("handleStdin() stdinReader = %v, want nil for binary mode", stdinReader)
	}

	// Verify binary file was created with correct content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read binary file: %v", err)
	}
	if string(content) != string(textData) {
		t.Errorf("binary file content = %q, want %q", string(content), string(textData))
	}
}
