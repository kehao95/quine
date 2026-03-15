package tools

import (
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func TestAllToolSchemas_Count(t *testing.T) {
	// Without config (nil), workspace revision tools should not be included.
	schemas := AllToolSchemas(nil)
	if len(schemas) != 5 {
		t.Fatalf("AllToolSchemas(nil) returned %d schemas, want 5", len(schemas))
	}
}

func TestAllToolSchemas_WithEscalation(t *testing.T) {
	cfg := &config.Config{
		SmartModelID: "claude-opus", // This enables escalation
		Escalated:    false,
	}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 6 {
		t.Fatalf("AllToolSchemas() with escalation should return 6 schemas, got %d", len(schemas))
	}

	// Verify escalate is in the list
	found := false
	for _, s := range schemas {
		if s.Name == "escalate" {
			found = true
			break
		}
	}
	if !found {
		t.Error("AllToolSchemas() should include escalate tool when escalation is configured")
	}
}

func TestAllToolSchemas_PostEscalation(t *testing.T) {
	cfg := &config.Config{
		SmartModelID: "claude-opus",
		Escalated:    true, // Already escalated
	}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 5 {
		t.Fatalf("AllToolSchemas() post-escalation should return 5 schemas, got %d", len(schemas))
	}

	// Verify escalate is NOT in the list
	for _, s := range schemas {
		if s.Name == "escalate" {
			t.Error("AllToolSchemas() should NOT include escalate tool after escalation")
		}
	}
}

func TestAllToolSchemas_WithAnchorMemory(t *testing.T) {
	cfg := &config.Config{
		AnchorMemoryEnabled: true,
	}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 7 {
		t.Fatalf("AllToolSchemas() with anchor memory should return 7 schemas, got %d", len(schemas))
	}
	var foundMark, foundUnfold bool
	for _, s := range schemas {
		if s.Name == "mark" {
			foundMark = true
		}
		if s.Name == "unfold" {
			foundUnfold = true
		}
	}
	if !foundMark || !foundUnfold {
		t.Fatalf("expected mark/unfold schemas, got mark=%v unfold=%v", foundMark, foundUnfold)
	}
}

func TestAllToolSchemas_WithAnchorMemoryAndEscalation(t *testing.T) {
	cfg := &config.Config{
		AnchorMemoryEnabled: true,
		SmartModelID:        "claude-opus",
	}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 8 {
		t.Fatalf("AllToolSchemas() with anchor memory + escalation should return 8 schemas, got %d", len(schemas))
	}
}

func TestRestoreWorldToolSchema(t *testing.T) {
	schema := RestoreWorldToolSchema(&config.Config{
		WorkspaceEnabled:      true,
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
	})
	if schema.Name != "restore_world" {
		t.Fatalf("schema name = %q, want restore_world", schema.Name)
	}
	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["revision"]; !ok {
		t.Fatal("restore_world schema should expose revision")
	}
	if !strings.Contains(schema.Description, "provisional workspace") {
		t.Fatalf("unexpected restore_world description: %q", schema.Description)
	}
}

func TestExecToolSchema_DoesNotExposeResetWorld(t *testing.T) {
	schema := ExecToolSchema(&config.Config{
		WorkspaceEnabled:      true,
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
	})
	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["reset_world"]; ok {
		t.Fatal("exec schema should not expose reset_world")
	}
	if strings.Contains(schema.Description, "reset_world") {
		t.Fatalf("exec schema should not mention reset_world, got: %q", schema.Description)
	}
}

func TestAllToolSchemas_WithWorkspaceRestoreMode(t *testing.T) {
	cfg := &config.Config{
		WorkspaceEnabled:      true,
		WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
	}
	schemas := AllToolSchemas(cfg)
	var foundRestore bool
	for _, s := range schemas {
		if s.Name == "restore_world" {
			foundRestore = true
		}
	}
	if !foundRestore {
		t.Fatal("restore mode should expose restore_world")
	}
}

func TestAllToolSchemas_WithWorkspaceNoneModeHidesRestore(t *testing.T) {
	cfg := &config.Config{
		WorkspaceEnabled:      true,
		WorkspaceRevisionMode: config.WorkspaceRevisionNone,
	}
	schemas := AllToolSchemas(cfg)
	for _, s := range schemas {
		if s.Name == "restore_world" {
			t.Fatal("none mode should not expose restore_world")
		}
	}
}

func TestShToolSchema_SingleModelMode(t *testing.T) {
	// Single-model mode: no SmartModelID configured
	schema := ShToolSchema(nil)

	props := schema.Parameters["properties"].(map[string]any)
	required := schema.Parameters["required"].([]string)

	// Should have command but NOT goal/strategy
	if _, ok := props["command"]; !ok {
		t.Error("sh schema should have 'command' property")
	}
	if _, ok := props["interactive"]; !ok {
		t.Error("sh schema should have 'interactive' property")
	}
	if _, ok := props["goal"]; ok {
		t.Error("sh schema should NOT have 'goal' in single-model mode")
	}
	if _, ok := props["strategy"]; ok {
		t.Error("sh schema should NOT have 'strategy' in single-model mode")
	}

	// Only command should be required
	if len(required) != 1 || required[0] != "command" {
		t.Errorf("sh schema required should be ['command'] in single-model mode, got %v", required)
	}
}

func TestShToolSchema_EscalationMode(t *testing.T) {
	// Escalation mode: SmartModelID configured and not yet escalated
	cfg := &config.Config{
		SmartModelID: "claude-opus",
		Escalated:    false,
	}
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	required := schema.Parameters["required"].([]string)

	// Should have command, goal, and strategy
	if _, ok := props["command"]; !ok {
		t.Error("sh schema should have 'command' property")
	}
	if _, ok := props["interactive"]; !ok {
		t.Error("sh schema should have 'interactive' in escalation mode")
	}
	if _, ok := props["goal"]; !ok {
		t.Error("sh schema should have 'goal' in escalation mode")
	}
	if _, ok := props["strategy"]; !ok {
		t.Error("sh schema should have 'strategy' in escalation mode")
	}

	// All three should be required
	if len(required) != 3 {
		t.Errorf("sh schema required should have 3 items in escalation mode, got %v", required)
	}
}

func TestShToolSchema_PostEscalation(t *testing.T) {
	// Post-escalation: SmartModelID configured but already escalated
	cfg := &config.Config{
		SmartModelID: "claude-opus",
		Escalated:    true,
	}
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	required := schema.Parameters["required"].([]string)

	// Post-escalation behaves like single-model mode (no STALL detection needed)
	if _, ok := props["goal"]; ok {
		t.Error("sh schema should NOT have 'goal' post-escalation")
	}
	if _, ok := props["strategy"]; ok {
		t.Error("sh schema should NOT have 'strategy' post-escalation")
	}

	// Only command should be required
	if len(required) != 1 || required[0] != "command" {
		t.Errorf("sh schema required should be ['command'] post-escalation, got %v", required)
	}
}
