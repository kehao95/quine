package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func TestConsumeExecutableBodyIfConfigured_RemovesLaunchPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "quine")
	if err := os.WriteFile(path, []byte("body"), 0o755); err != nil {
		t.Fatalf("write launch path: %v", err)
	}

	cfg := &config.Config{
		ToolGates: config.ToolGates{EphemeralBody: true},
		Paths:     config.Paths{ExecutablePath: path},
	}

	if err := consumeExecutableBodyIfConfigured(cfg); err != nil {
		t.Fatalf("consumeExecutableBodyIfConfigured() error: %v", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("launch path should be removed, stat error = %v", err)
	}
}

func TestConsumeExecutableBodyIfConfigured_ExecutablePathModeDoesNotPreserveBody(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "quine")
	if err := os.WriteFile(path, []byte("body"), 0o755); err != nil {
		t.Fatalf("write launch path: %v", err)
	}

	cfg := &config.Config{
		Identity:  config.Identity{SessionID: "sess-body", RunID: "run-body"},
		ToolGates: config.ToolGates{EphemeralBody: true},
		Paths: config.Paths{
			DataDir:           filepath.Join(root, "runtime"),
			ExecutablePath:    path,
			SelfReentryMode:   config.SelfReentryModeExecutablePath,
			SelfReentryTarget: path,
		},
	}

	if err := consumeExecutableBodyIfConfigured(cfg); err != nil {
		t.Fatalf("consumeExecutableBodyIfConfigured() error: %v", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("launch path should be removed, stat error = %v", err)
	}
	if cfg.SelfReentryTarget != path {
		t.Fatalf("SelfReentryTarget = %q, want unchanged launch path %q", cfg.SelfReentryTarget, path)
	}
	if _, err := os.Stat(filepath.Join(cfg.DataDir, "body")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("runtime body directory should not be created, stat error = %v", err)
	}
}

func TestConsumeExecutableBodyIfConfigured_LeavesLaunchPathWhenDisabled(t *testing.T) {
	path := filepath.Join(t.TempDir(), "quine")
	if err := os.WriteFile(path, []byte("body"), 0o755); err != nil {
		t.Fatalf("write launch path: %v", err)
	}

	cfg := &config.Config{
		ToolGates: config.ToolGates{EphemeralBody: false},
		Paths:     config.Paths{ExecutablePath: path},
	}

	if err := consumeExecutableBodyIfConfigured(cfg); err != nil {
		t.Fatalf("consumeExecutableBodyIfConfigured() error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("launch path should remain when disabled, stat error = %v", err)
	}
}
