package runtime

import (
	"strings"
	"testing"
)

// TestAssembleSystemPromptEmptyLeadingFragment pins E1: an earlier fragment that
// renders empty (the default empty memory fragment is routine) must not leave a
// stray double blank line before the first real block.
func TestAssembleSystemPromptEmptyLeadingFragment(t *testing.T) {
	base := "BASE"
	fragments := []promptFragment{
		{Name: "30-memory.md", Content: ""}, // renders empty → skipped
		{Name: "10-agents.md", Content: "Be good."},
	}
	got := assembleSystemPrompt(base, fragments)

	if strings.Contains(got, "\n\n\n") {
		t.Fatalf("stray blank line in assembled prompt:\n%q", got)
	}
	want := "BASE\n\n### AGENTS.md\nBe good.\n"
	if got != want {
		t.Fatalf("assembled prompt mismatch:\n got %q\nwant %q", got, want)
	}
}

func TestAssembleSystemPromptSeparatesEmittedBlocks(t *testing.T) {
	got := assembleSystemPrompt("BASE", []promptFragment{
		{Name: "10-agents.md", Content: "A"},
		{Name: "30-memory.md", Content: ""}, // skipped, no spurious separator
		{Name: "40-mission.md", Content: "M"},
	})
	want := "BASE\n\n### AGENTS.md\nA\n\n### Your Mission\nM\n"
	if got != want {
		t.Fatalf("assembled prompt mismatch:\n got %q\nwant %q", got, want)
	}
}

// TestInvalidShCommandReason pins A2: a missing, non-string, or blank command is
// rejected rather than coerced to "" and run as an exit-0 no-op.
func TestInvalidShCommandReason(t *testing.T) {
	bad := []map[string]any{
		{},                 // absent
		{"command": 42},    // non-string
		{"command": ""},    // empty
		{"command": "   "}, // blank
	}
	for _, args := range bad {
		if _, isBad := invalidShCommandReason(args); !isBad {
			t.Fatalf("expected %v to be rejected", args)
		}
	}
	if _, isBad := invalidShCommandReason(map[string]any{"command": "ls"}); isBad {
		t.Fatal("a real command must not be rejected")
	}
}
