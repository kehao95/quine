package world

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndPayload(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "world.json")
	if err := os.WriteFile(path, []byte(`{"items":{"1":"alpha","2":"beta"}}`), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	spec, err := Load(path)
	if err != nil {
		t.Fatalf("Load(): %v", err)
	}

	got, err := spec.Payload("2")
	if err != nil {
		t.Fatalf("Payload(): %v", err)
	}
	if got != "beta" {
		t.Fatalf("Payload() = %q, want %q", got, "beta")
	}
}

func TestPayloadUnknownID(t *testing.T) {
	t.Parallel()

	spec := &Spec{Items: map[string]string{"1": "alpha"}}
	_, err := spec.Payload("9")
	if !errors.Is(err, ErrUnknownID) {
		t.Fatalf("Payload() error = %v, want ErrUnknownID", err)
	}
}

func TestDefaultSpecPathUsesQuineDataDir(t *testing.T) {
	t.Setenv("QUINE_DATA_DIR", "/tmp/runtime")
	if got := DefaultSpecPath(); got != "/tmp/runtime/world/world.json" {
		t.Fatalf("DefaultSpecPath() = %q, want %q", got, "/tmp/runtime/world/world.json")
	}
}

func TestLoadDefaultUsesEmbeddedSpec(t *testing.T) {
	original := embeddedSpecBase64
	embeddedSpecBase64 = "eyJpdGVtcyI6eyIxIjoiYWxwaGEifX0="
	t.Cleanup(func() {
		embeddedSpecBase64 = original
	})

	spec, source, err := LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault(): %v", err)
	}
	if source != "embedded" {
		t.Fatalf("LoadDefault() source = %q, want %q", source, "embedded")
	}
	got, err := spec.Payload("1")
	if err != nil {
		t.Fatalf("Payload(): %v", err)
	}
	if got != "alpha" {
		t.Fatalf("Payload() = %q, want %q", got, "alpha")
	}
}
