package runtime

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncPromptFragmentsOmitsMissionFragmentWhenMissionAbsent(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = ""

	fragmentsRoot := rt.currentIncarnationPromptRoot()
	if err := os.MkdirAll(fragmentsRoot, 0o755); err != nil {
		t.Fatalf("mkdir fragments root: %v", err)
	}
	missionPath := filepath.Join(fragmentsRoot, "40-mission.md")
	if err := os.WriteFile(missionPath, []byte("stale mission\n"), 0o644); err != nil {
		t.Fatalf("write stale mission: %v", err)
	}

	if err := rt.syncPromptFragments(); err != nil {
		t.Fatalf("syncPromptFragments: %v", err)
	}

	if _, err := os.Stat(missionPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("40-mission.md should be absent for missionless runs, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(fragmentsRoot, "30-memory.md")); err != nil {
		t.Fatalf("30-memory.md should still exist: %v", err)
	}
}

func TestProviderContextMessagesOmitMissionSectionWhenMissionAbsent(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = ""

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot: %v", err)
	}
	t.Cleanup(func() {
		_ = rt.cleanupAgentRoot()
	})

	msgs, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("providerContextMessages len = %d, want 1", len(msgs))
	}
	systemPrompt := msgs[0].Content
	if strings.Contains(systemPrompt, "### Your Mission") {
		t.Fatalf("missionless system prompt should omit Your Mission section:\n%s", systemPrompt)
	}
	if !strings.Contains(systemPrompt, "You are quine, a running process in a POSIX operating system.") {
		t.Fatalf("missionless system prompt missing opening:\n%s", systemPrompt)
	}
	if strings.Contains(systemPrompt, "There won't be further instructions.") ||
		strings.Contains(systemPrompt, "Act autonomously") ||
		strings.Contains(systemPrompt, "Do not wait or quiesce") {
		t.Fatalf("missionless system prompt should not force autonomous action:\n%s", systemPrompt)
	}
}

func TestRenderPromptFragmentBlockForkAssignment(t *testing.T) {
	got := renderPromptFragmentBlock(promptFragment{
		Name:    "45-fork-assignment.md",
		Content: "lane focus",
	})
	if got != "### Fork Assignment\nlane focus" {
		t.Fatalf("fork assignment fragment render = %q", got)
	}
}
