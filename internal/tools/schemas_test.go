package tools

import (
	"testing"
)

func TestAllToolSchemas_Count(t *testing.T) {
	schemas := AllToolSchemas()
	if len(schemas) != 5 {
		t.Fatalf("AllToolSchemas() returned %d schemas, want 5", len(schemas))
	}
}

func TestShSchema_RequiredCommand(t *testing.T) {
	s := ShToolSchema()
	if s.Name != "sh" {
		t.Fatalf("Name = %q, want %q", s.Name, "sh")
	}
	required, ok := s.Parameters["required"].([]string)
	if !ok {
		t.Fatal("Parameters[\"required\"] is not []string")
	}
	found := false
	for _, r := range required {
		if r == "command" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("sh schema required fields %v do not include \"command\"", required)
	}
}

func TestExitSchema_RequiredStatus(t *testing.T) {
	s := ExitToolSchema()
	if s.Name != "exit" {
		t.Fatalf("Name = %q, want %q", s.Name, "exit")
	}
	required, ok := s.Parameters["required"].([]string)
	if !ok {
		t.Fatal("Parameters[\"required\"] is not []string")
	}
	found := false
	for _, r := range required {
		if r == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("exit schema required fields %v do not include \"status\"", required)
	}
}

func TestForkSchema_RequiredIntent(t *testing.T) {
	s := ForkToolSchema()
	if s.Name != "fork" {
		t.Fatalf("Name = %q, want %q", s.Name, "fork")
	}
	required, ok := s.Parameters["required"].([]string)
	if !ok {
		t.Fatal("Parameters[\"required\"] is not []string")
	}
	found := false
	for _, r := range required {
		if r == "intent" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("fork schema required fields %v do not include \"intent\"", required)
	}
}

func TestForkSchema_HasWaitProperty(t *testing.T) {
	s := ForkToolSchema()
	props, ok := s.Parameters["properties"].(map[string]any)
	if !ok {
		t.Fatal("Parameters[\"properties\"] is not map[string]any")
	}
	wait, ok := props["wait"].(map[string]any)
	if !ok {
		t.Fatal("wait property is not map[string]any")
	}
	if wait["type"] != "boolean" {
		t.Errorf("wait type = %q, want \"boolean\"", wait["type"])
	}
}
