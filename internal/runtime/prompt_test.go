package runtime

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func testConfig() *config.Config {
	return &config.Config{
		ModelID:   "claude-sonnet-4-20250514",
		Depth:     2,
		MaxDepth:  5,
		MaxTurns:  20,
		SessionID: "abc-123-def-456",
		Shell:     "/bin/zsh",
		Wisdom:    nil,
	}
}

func TestBuildSystemPrompt_NoRawPlaceholders(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission")

	placeholders := []string{"{DEPTH}", "{MAX_DEPTH}", "{MAX_TURNS}", "{MODEL_ID}", "{SESSION_ID}", "{SHELL}", "{WISDOM}", "{MISSION}"}
	for _, ph := range placeholders {
		if strings.Contains(prompt, ph) {
			t.Errorf("prompt still contains unsubstituted placeholder %s", ph)
		}
	}
}

func TestBuildSystemPrompt_CorrectValues(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission")

	checks := map[string]string{
		"Depth":     fmt.Sprintf("%d / %d", cfg.Depth, cfg.MaxDepth),
		"ModelID":   cfg.ModelID,
		"SessionID": cfg.SessionID,
		"Shell":     cfg.Shell,
	}

	for name, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %s value %q", name, want)
		}
	}
}

func TestBuildSystemPrompt_MaxTurnsUnlimited(t *testing.T) {
	cfg := testConfig()
	cfg.MaxTurns = 0
	prompt := BuildSystemPrompt(cfg, "test mission")

	if !strings.Contains(prompt, "Shell Executions Remaining: unlimited") {
		t.Error("prompt should show 'Shell Executions Remaining: unlimited' when MaxTurns is 0")
	}
}

func TestBuildSystemPrompt_MaxTurnsNumeric(t *testing.T) {
	cfg := testConfig()
	cfg.MaxTurns = 15
	prompt := BuildSystemPrompt(cfg, "test mission")

	if !strings.Contains(prompt, "Shell Executions Remaining: 15") {
		t.Error("prompt should show 'Shell Executions Remaining: 15' when MaxTurns is 15")
	}
}

func TestBuildSystemPrompt_KeySections(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission")

	sections := []string{
		"### THE PRIME DIRECTIVE",
		"### Environment",
		"### Mortality",
		"### Tools",
		"### SURVIVAL PROTOCOLS",
		"### Semantic Gradient",
		"### Output",
	}

	for _, sec := range sections {
		if !strings.Contains(prompt, sec) {
			t.Errorf("prompt missing section %q", sec)
		}
	}
}

func TestBuildSystemPrompt_DepthZero(t *testing.T) {
	cfg := &config.Config{
		ModelID:   "gpt-4o",
		Depth:     0,
		MaxDepth:  10,
		SessionID: "root-session",
		Shell:     "/bin/sh",
	}
	prompt := BuildSystemPrompt(cfg, "test mission")

	if !strings.Contains(prompt, "0 / 10") {
		t.Error("prompt should contain '0 / 10' for depth 0, max 10")
	}
}

func TestBuildSystemPrompt_WithWisdom(t *testing.T) {
	cfg := &config.Config{
		ModelID:   "claude-sonnet-4-20250514",
		Depth:     1,
		MaxDepth:  5,
		MaxTurns:  20,
		SessionID: "test-session",
		Shell:     "/bin/sh",
		Wisdom: map[string]string{
			"SUMMARY": "User prefers concise answers",
			"CONTEXT": "Working on Go project",
		},
	}
	prompt := BuildSystemPrompt(cfg, "test mission")

	// Check wisdom section is present
	if !strings.Contains(prompt, "### Wisdom (from previous incarnation)") {
		t.Error("prompt should contain wisdom section header")
	}
	if !strings.Contains(prompt, "**SUMMARY**: User prefers concise answers") {
		t.Error("prompt should contain SUMMARY wisdom entry")
	}
	if !strings.Contains(prompt, "**CONTEXT**: Working on Go project") {
		t.Error("prompt should contain CONTEXT wisdom entry")
	}
}

func TestBuildSystemPrompt_WithoutWisdom(t *testing.T) {
	cfg := &config.Config{
		ModelID:   "claude-sonnet-4-20250514",
		Depth:     0,
		MaxDepth:  5,
		MaxTurns:  20,
		SessionID: "test-session",
		Shell:     "/bin/sh",
		Wisdom:    nil,
	}
	prompt := BuildSystemPrompt(cfg, "test mission")

	// Check wisdom section is NOT present
	if strings.Contains(prompt, "### Wisdom") {
		t.Error("prompt should NOT contain wisdom section when wisdom is nil")
	}
}

func TestBuildSystemPrompt_EmptyWisdom(t *testing.T) {
	cfg := &config.Config{
		ModelID:   "claude-sonnet-4-20250514",
		Depth:     0,
		MaxDepth:  5,
		MaxTurns:  20,
		SessionID: "test-session",
		Shell:     "/bin/sh",
		Wisdom:    map[string]string{},
	}
	prompt := BuildSystemPrompt(cfg, "test mission")

	// Check wisdom section is NOT present for empty map
	if strings.Contains(prompt, "### Wisdom") {
		t.Error("prompt should NOT contain wisdom section when wisdom is empty")
	}
}

func TestBuildSystemPrompt_WisdomSorted(t *testing.T) {
	cfg := &config.Config{
		ModelID:   "claude-sonnet-4-20250514",
		Depth:     1,
		MaxDepth:  5,
		MaxTurns:  20,
		SessionID: "test-session",
		Shell:     "/bin/sh",
		Wisdom: map[string]string{
			"ZEBRA":  "last alphabetically",
			"APPLE":  "first alphabetically",
			"MIDDLE": "in between",
		},
	}
	prompt := BuildSystemPrompt(cfg, "test mission")

	// Check keys are sorted
	appleIdx := strings.Index(prompt, "**APPLE**")
	middleIdx := strings.Index(prompt, "**MIDDLE**")
	zebraIdx := strings.Index(prompt, "**ZEBRA**")

	if appleIdx == -1 || middleIdx == -1 || zebraIdx == -1 {
		t.Fatal("all wisdom keys should be present")
	}

	if !(appleIdx < middleIdx && middleIdx < zebraIdx) {
		t.Error("wisdom keys should be sorted alphabetically")
	}
}

func TestBuildSystemPrompt_WisdomNoRawPlaceholder(t *testing.T) {
	cfg := testConfig()
	cfg.Wisdom = map[string]string{"TEST": "value"}
	prompt := BuildSystemPrompt(cfg, "test mission")

	if strings.Contains(prompt, "{WISDOM}") {
		t.Error("prompt should not contain raw {WISDOM} placeholder")
	}
}
