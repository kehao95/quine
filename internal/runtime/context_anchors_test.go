package runtime

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

func TestProviderContextMessagesUseContextCurrent(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "inspect context"
	rt.hasMaterial = true
	rt.tape = tape.NewTape(cfg.SessionID, "", cfg.Depth, cfg.ModelID)
	rt.tape.Append(tape.Message{Role: tape.RoleAssistant, Content: "from tape only"})

	currentPath := filepath.Join(cfg.AgentRoot(), "context", "state", "current.jsonl")
	writeJSONLEntries(t, currentPath,
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "from context user"}),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "from context assistant"}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID:  "tool-1",
			Content: json.RawMessage(`{"tool":"sh","status":"ok"}`),
		}),
	)

	msgs, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages: %v", err)
	}
	if len(msgs) != 3 {
		t.Fatalf("providerContextMessages len = %d, want 3", len(msgs))
	}
	if msgs[0].Role != tape.RoleSystem {
		t.Fatalf("msgs[0].Role = %q, want system", msgs[0].Role)
	}
	if msgs[1].Content != "from context user" {
		t.Fatalf("msgs[1].Content = %q, want context user", msgs[1].Content)
	}
	if msgs[2].Content != "from context assistant" {
		t.Fatalf("msgs[2].Content = %q, want context assistant", msgs[2].Content)
	}
	for _, msg := range msgs {
		if msg.ToolID == "tool-1" {
			t.Fatal("provider context should drop orphan tool results without a preceding tool batch")
		}
	}
	for _, msg := range msgs {
		if msg.Content == "from tape only" {
			t.Fatal("provider context should be assembled from context/state/current.jsonl, not tape messages")
		}
	}
}

func TestProviderContextMessagesRejectIncompleteToolBatch(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "reject dangling tool batch"

	currentPath := filepath.Join(cfg.AgentRoot(), "context", "state", "current.jsonl")
	writeJSONLEntries(t, currentPath,
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "from context user"}),
		tape.MessageEntry(tape.Message{
			Role: tape.RoleAssistant,
			ToolCalls: []tape.ToolCall{
				{ID: "call_fork", Name: "fork", Arguments: map[string]any{"children": []any{}}},
			},
		}),
	)

	if _, err := rt.providerContextMessages(); err == nil || !strings.Contains(err.Error(), "incomplete tool batch") {
		t.Fatalf("providerContextMessages err = %v, want incomplete tool batch error", err)
	}
}

func TestProviderContextMessagesDropOrphanToolResults(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "drop orphan tool results"

	currentPath := filepath.Join(cfg.AgentRoot(), "context", "state", "current.jsonl")
	writeJSONLEntries(t, currentPath,
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "from context user"}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID:  "call_mark",
			Content: json.RawMessage(`{"tool":"mark","status":"completed"}`),
		}),
		tape.MessageEntry(tape.Message{
			Role: tape.RoleAssistant,
			ToolCalls: []tape.ToolCall{
				{ID: "call_sh", Name: "sh", Arguments: map[string]any{"command": "pwd"}},
			},
		}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID:  "call_sh",
			Content: json.RawMessage(`{"tool":"sh","status":"completed"}`),
		}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID:  "call_vision",
			Content: json.RawMessage(`{"tool":"vision","status":"completed"}`),
		}),
	)

	msgs, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages: %v", err)
	}
	if len(msgs) != 4 {
		t.Fatalf("providerContextMessages len = %d, want 4", len(msgs))
	}
	if msgs[1].Role != tape.RoleUser || msgs[1].Content != "from context user" {
		t.Fatalf("msgs[1] = %+v, want retained user message", msgs[1])
	}
	if msgs[2].Role != tape.RoleAssistant || len(msgs[2].ToolCalls) != 1 || msgs[2].ToolCalls[0].ID != "call_sh" {
		t.Fatalf("msgs[2] = %+v, want retained sh batch", msgs[2])
	}
	if msgs[3].Role != tape.RoleToolResult || msgs[3].ToolID != "call_sh" {
		t.Fatalf("msgs[3] = %+v, want retained sh tool result", msgs[3])
	}
	for _, msg := range msgs {
		if msg.ToolID == "call_mark" || msg.ToolID == "call_vision" {
			t.Fatalf("provider context should drop orphan tool result %q", msg.ToolID)
		}
	}
}

func TestProviderContextMessagesIncludeFrontierAnchors(t *testing.T) {
	cfg := testCfg(t)
	cfg.AnchorMemoryEnabled = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "inspect anchors"

	currentPath := filepath.Join(cfg.AgentRoot(), "context", "state", "current.jsonl")
	writeJSONLEntries(t, currentPath, tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "live current"}))

	frontierDir := filepath.Join(cfg.AgentRoot(), "context", "state", "frontier")
	anchorDir := filepath.Join(cfg.AgentRoot(), "context", "state", "anchors", "2.anchor")
	if err := os.MkdirAll(frontierDir, 0o755); err != nil {
		t.Fatalf("mkdir frontier: %v", err)
	}
	if err := os.MkdirAll(anchorDir, 0o755); err != nil {
		t.Fatalf("mkdir anchor: %v", err)
	}
	if err := os.WriteFile(filepath.Join(frontierDir, "2.link"), []byte(""), 0o644); err != nil {
		t.Fatalf("write frontier link: %v", err)
	}
	if err := os.WriteFile(filepath.Join(anchorDir, "meta.json"), []byte(`{"id":2,"resolution":"checkpoint-2"}`), 0o644); err != nil {
		t.Fatalf("write anchor meta: %v", err)
	}

	msgs, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages: %v", err)
	}
	if len(msgs) != 3 {
		t.Fatalf("providerContextMessages len = %d, want 3", len(msgs))
	}
	if msgs[1].Role != tape.RoleAssistant || msgs[1].Content != "[ANCHOR 2] checkpoint-2" {
		t.Fatalf("msgs[1] = %#v, want anchor summary", msgs[1])
	}
	if msgs[2].Role != tape.RoleUser || msgs[2].Content != "live current" {
		t.Fatalf("msgs[2] = %#v, want live current user message", msgs[2])
	}
}

func TestSyncPromptFragmentsWritesManagedSurface(t *testing.T) {
	cfg := testCfg(t)
	cfg.WorkDir = filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(cfg.WorkDir, 0o755); err != nil {
		t.Fatalf("mkdir workdir: %v", err)
	}

	agentsSource := filepath.Join(cfg.WorkDir, "AGENTS.md")
	if err := os.WriteFile(agentsSource, []byte("ROOT_AGENTS_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	cfg.AgentsMDEnabled = true
	cfg.AgentsMDPath = agentsSource
	cfg.AgentsSkillsEnabled = true
	cfg.Skills = []config.Skill{
		{
			Name:        "foo",
			Description: "Skill description marker",
			Source:      ".agents/skills/foo/SKILL.md",
		},
	}

	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "raw mission text"

	if err := rt.syncPromptFragments(); err != nil {
		t.Fatalf("syncPromptFragments: %v", err)
	}

	fragmentsRoot := rt.currentIncarnationPromptRoot()
	missionPath := filepath.Join(fragmentsRoot, "40-mission.md")
	agentsPath := filepath.Join(fragmentsRoot, "10-agents.md")
	skillsPath := filepath.Join(fragmentsRoot, "20-skills.md")
	memoryPath := filepath.Join(fragmentsRoot, "30-memory.md")
	runtimePath := filepath.Join(fragmentsRoot, "00-runtime.md")

	missionData, err := os.ReadFile(missionPath)
	if err != nil {
		t.Fatalf("read mission fragment: %v", err)
	}
	if string(missionData) != "raw mission text\n" {
		t.Fatalf("mission fragment = %q, want raw mission text with trailing newline", string(missionData))
	}

	assertSymlinkTarget(t, agentsPath, agentsSource)
	assertSameFile(t, agentsPath, agentsSource)

	skillsData, err := os.ReadFile(skillsPath)
	if err != nil {
		t.Fatalf("read skills fragment: %v", err)
	}
	skillsText := string(skillsData)
	if !strings.Contains(skillsText, "Skill description marker") {
		t.Fatalf("skills fragment missing description: %q", skillsText)
	}
	if !strings.Contains(skillsText, "Source: `.agents/skills/foo/SKILL.md`") {
		t.Fatalf("skills fragment missing source path: %q", skillsText)
	}
	if memoryData, err := os.ReadFile(memoryPath); err != nil {
		t.Fatalf("read memory surface: %v", err)
	} else if string(memoryData) != "\n" {
		t.Fatalf("memory surface = %q, want single trailing newline by default", string(memoryData))
	}
	if runtimeData, err := os.ReadFile(runtimePath); err != nil {
		t.Fatalf("read runtime surface: %v", err)
	} else if !strings.Contains(string(runtimeData), "### Context Files") {
		t.Fatalf("runtime surface missing context-files block: %q", string(runtimeData))
	}
}

func TestSyncPromptFragmentsRemovesDisabledOptionalFragments(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "mission only"

	fragmentsRoot := rt.currentIncarnationPromptRoot()
	if err := os.MkdirAll(fragmentsRoot, 0o755); err != nil {
		t.Fatalf("mkdir fragments root: %v", err)
	}
	if err := os.WriteFile(filepath.Join(fragmentsRoot, "10-agents.md"), []byte("stale agents"), 0o644); err != nil {
		t.Fatalf("write stale AGENTS.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(fragmentsRoot, "20-skills.md"), []byte("stale skills"), 0o644); err != nil {
		t.Fatalf("write stale SKILLS.md: %v", err)
	}

	if err := rt.syncPromptFragments(); err != nil {
		t.Fatalf("syncPromptFragments: %v", err)
	}

	fragmentsRoot = rt.currentIncarnationPromptRoot()
	if _, err := os.Stat(filepath.Join(fragmentsRoot, "40-mission.md")); err != nil {
		t.Fatalf("40-mission.md should exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(fragmentsRoot, "30-memory.md")); err != nil {
		t.Fatalf("30-memory.md should exist: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(fragmentsRoot, "10-agents.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("10-agents.md should be removed when gate is disabled, got %v", err)
	}
	if _, err := os.Lstat(filepath.Join(fragmentsRoot, "20-skills.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("20-skills.md should be removed when gate is disabled, got %v", err)
	}
}

func TestProviderContextMessagesAssemblePromptFragments(t *testing.T) {
	cfg := testCfg(t)
	cfg.WorkDir = filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(cfg.WorkDir, 0o755); err != nil {
		t.Fatalf("mkdir workdir: %v", err)
	}

	agentsSource := filepath.Join(cfg.WorkDir, "AGENTS.md")
	if err := os.WriteFile(agentsSource, []byte("ASSEMBLED_AGENTS_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	cfg.AgentsMDEnabled = true
	cfg.AgentsMDPath = agentsSource
	cfg.AgentsSkillsEnabled = true
	cfg.Skills = []config.Skill{
		{
			Name:        "foo",
			Description: "Assembled skills marker",
			Source:      ".agents/skills/foo/SKILL.md",
		},
	}

	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "assembled mission marker"
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot: %v", err)
	}
	t.Cleanup(func() {
		_ = rt.cleanupAgentRoot()
	})
	if err := os.MkdirAll(rt.currentIncarnationPromptRoot(), 0o755); err != nil {
		t.Fatalf("mkdir context root: %v", err)
	}
	if err := os.WriteFile(rt.currentIncarnationPromptFile("30-memory.md"), []byte("ASSEMBLED_MEMORY_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write memory: %v", err)
	}

	msgs, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("providerContextMessages len = %d, want 1", len(msgs))
	}
	if msgs[0].Role != tape.RoleSystem {
		t.Fatalf("msgs[0].Role = %q, want system", msgs[0].Role)
	}

	systemPrompt := msgs[0].Content
	for _, want := range []string{
		"### Context Files",
		"### Your Mission\nassembled mission marker",
		"### AGENTS.md\nASSEMBLED_AGENTS_MARKER",
		"### SKILLS.md\nThese skills are available through the project surface when relevant.",
		"Assembled skills marker",
		"### Memory\nASSEMBLED_MEMORY_MARKER",
	} {
		if !strings.Contains(systemPrompt, want) {
			t.Fatalf("system prompt missing %q\nprompt=%s", want, systemPrompt)
		}
	}
}

func TestProviderContextMessagesRefreshMemoryBetweenTurns(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_write_memory",
						Name: "sh",
						Arguments: map[string]any{
							"command": "printf '%s\\n' MEMORY_REFRESH_MARKER_UNIT > \"$QUINE_AGENT_ROOT/context/prompt/30-memory.md\"",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.MaxTurns = 3
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("refresh memory between turns", "Begin."); exitCode != 0 {
		t.Fatalf("Run exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) < 2 {
		t.Fatalf("provider calls = %d, want at least 2", len(mock.calls))
	}

	secondSystemPrompt := mock.calls[1][0].Content
	for _, want := range []string{
		"### Memory\nMEMORY_REFRESH_MARKER_UNIT",
		"### Your Mission\nrefresh memory between turns",
	} {
		if !strings.Contains(secondSystemPrompt, want) {
			t.Fatalf("second provider system prompt missing %q\nprompt=%s", want, secondSystemPrompt)
		}
	}
}

func TestSyncPromptFragmentsPreservesMemoryFile(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "memory preservation mission"

	if err := os.MkdirAll(rt.currentIncarnationPromptRoot(), 0o755); err != nil {
		t.Fatalf("mkdir context root: %v", err)
	}
	memoryPath := rt.currentIncarnationPromptFile("30-memory.md")
	if err := os.WriteFile(memoryPath, []byte("PRESERVED_MEMORY_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write memory: %v", err)
	}

	if err := rt.syncPromptFragments(); err != nil {
		t.Fatalf("syncPromptFragments: %v", err)
	}

	memoryData, err := os.ReadFile(memoryPath)
	if err != nil {
		t.Fatalf("read memory: %v", err)
	}
	if string(memoryData) != "PRESERVED_MEMORY_MARKER\n" {
		t.Fatalf("memory.md = %q, want preserved content", string(memoryData))
	}
}
