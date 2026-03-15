package runtime

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func testConfig() *config.Config {
	return &config.Config{
		ModelID:              "claude-sonnet-4-20250514",
		Depth:                2,
		MaxDepth:             5,
		MaxTurns:             20,
		TurnExhaustionPolicy: config.TurnExhaustionHardFail,
		PromptMetaphor:       config.PromptMetaphorOff,
		SessionID:            "abc-123-def-456",
		DataDir:              "/tmp/quine-state",
		Shell:                "/bin/zsh",
		Wisdom:               nil,
	}
}

func TestBuildSystemPrompt_NoRawPlaceholders(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	placeholders := []string{
		"{PRIME_DIRECTIVE_TITLE}",
		"{PRIME_DIRECTIVE_BODY}",
		"{MODEL}",
		"{ESCALATION_TIER_LINE}",
		"{SHELL}",
		"{DEPTH}",
		"{LIMITS_BLOCK}",
		"{RUNTIME_BLOCK}",
		"{WISDOM}",
		"{ACTIVE_CONSTRAINTS}",
		"{STDIN_BLOCK}",
		"{SH_WORKSPACE_BLOCK}",
		"{SH_MATERIAL_LINE}",
		"{SH_GOAL_STRATEGY_LINE}",
		"{FORK_WORKSPACE_BLOCK}",
		"{EXEC_MATERIAL_LINE}",
		"{EXEC_BUDGET_LINE}",
		"{EXEC_PROTOCOL_LINE}",
		"{ESCALATION_TOOL_BLOCK}",
		"{ESCALATION_PROTOCOL_BLOCK}",
		"{MISSION}",
	}
	for _, ph := range placeholders {
		if strings.Contains(prompt, ph) {
			t.Errorf("prompt still contains unsubstituted placeholder %s", ph)
		}
	}
}

func TestBuildSystemPrompt_CorrectValues(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := map[string]string{
		"Depth":   fmt.Sprintf("- Depth: %d", cfg.Depth),
		"ModelID": cfg.ModelID,
		"Shell":   cfg.Shell,
	}

	for name, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %s value %q", name, want)
		}
	}
}

func TestBuildSystemPrompt_WorkspacePhysics(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.DataDir = "/tmp/quine-state"
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project/app"
	cfg.WorkspaceBackend = "direct"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Runtime State Root: `/tmp/quine-state` (jobs, tapes, logs; not the workspace)",
		"Workspace Root: `/tmp/project`",
		"Workspace: `/tmp/project/app`",
		"Workspace Backend: `direct`",
		"Workspace Revision Mode: `restore`",
		"Current World Revision: `wr0`",
		"Transactional Workspace: enabled",
		"`[FS MUTATIONS]` is the authoritative record",
		"`[FS MUTATIONS]` at end of each result",
		"Filesystem state is transactional",
		"share your mounted workspace surface",
		"**restore_world** - Restore workspace to a prior world revision such as `wr0` or `wr3`.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing workspace physics text %q", want)
		}
	}
}

func TestBuildSystemPrompt_DefaultPhysicsMode(t *testing.T) {
	cfg := testConfig()
	cfg.PromptMetaphor = config.PromptMetaphorOff
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "### THE PRIME DIRECTIVE: RUNTIME PHYSICS") {
		t.Error("prompt should contain runtime physics title when metaphor is off")
	}
	if strings.Contains(prompt, "THERMODYNAMIC SURVIVAL") {
		t.Error("prompt should not contain thermodynamic title when metaphor is off")
	}
	if strings.Contains(prompt, "Energy") || strings.Contains(prompt, "Entropy") {
		t.Error("prompt should not contain thermodynamic wording when metaphor is off")
	}
	if strings.Contains(prompt, "Filesystem state is transactional") {
		t.Error("prompt should not describe transactional workspace physics when workspace physics are disabled")
	}
	if strings.Contains(prompt, "Every completed `sh` result ends with `[FS MUTATIONS]`") {
		t.Error("prompt should not mention FS mutations when workspace physics are disabled")
	}
	if strings.Contains(prompt, "**restore_world**") {
		t.Error("prompt should not mention restore_world when workspace revision is unavailable")
	}
}

func TestBuildSystemPrompt_WorkspaceRestoreDoesNotMentionResetWorld(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project"
	cfg.WorkspaceBackend = "direct"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if !strings.Contains(prompt, "**restore_world** - Restore workspace to a prior world revision such as `wr0` or `wr3`.") {
		t.Error("prompt should expose restore_world in restore mode")
	}
	if strings.Contains(prompt, "reset_world=true") {
		t.Error("prompt should not mention reset_world")
	}
}

func TestBuildSystemPrompt_ThermodynamicOverlay(t *testing.T) {
	cfg := testConfig()
	cfg.PromptMetaphor = config.PromptMetaphorThermodynamic
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "### THE PRIME DIRECTIVE: THERMODYNAMIC SURVIVAL") {
		t.Error("prompt should contain thermodynamic title when metaphor is enabled")
	}
	if !strings.Contains(prompt, "Energy") || !strings.Contains(prompt, "Entropy") {
		t.Error("prompt should contain thermodynamic wording when metaphor is enabled")
	}
}

func TestBuildSystemPrompt_ExecutionBudgetVisibleWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.MaxTurns = 15
	cfg.TurnExhaustionPolicy = config.TurnExhaustionHardFail
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "Execution Budget: 15") {
		t.Error("prompt should show execution budget when MaxTurns > 0")
	}
	if !strings.Contains(prompt, "Each `sh` call consumes 1 execution") {
		t.Error("prompt should describe sh-only execution counting")
	}
	if !strings.Contains(prompt, "process terminates immediately") {
		t.Error("hard_fail policy should be disclosed")
	}
	if !strings.Contains(prompt, "resets the execution budget") {
		t.Error("prompt should describe exec budget reset when budget is enabled")
	}
}

func TestBuildSystemPrompt_ExecutionBudgetHiddenWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.MaxTurns = 0
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if strings.Contains(prompt, "Execution Budget:") {
		t.Error("prompt should omit execution budget when MaxTurns is 0")
	}
	if strings.Contains(prompt, "Execution budget:") {
		t.Error("prompt should omit execution constraint bullets when MaxTurns is 0")
	}
	if strings.Contains(prompt, "resets the execution budget") {
		t.Error("prompt should not claim exec resets execution budget when budget is disabled")
	}
}

func TestBuildSystemPrompt_NearDeathPolicyVisibility(t *testing.T) {
	cfg := testConfig()
	cfg.MaxTurns = 8
	cfg.TurnExhaustionPolicy = config.TurnExhaustionNearDeathExec
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "only `exec` is accepted") {
		t.Error("near_death_exec policy should be disclosed in active constraints")
	}
}

func TestBuildSystemPrompt_ConditionalLimitsVisibility(t *testing.T) {
	cfg := testConfig()
	cfg.MaxDepth = 0
	cfg.MaxConcurrent = 0
	cfg.MaxAgents = 0
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if strings.Contains(prompt, "Depth Limit:") {
		t.Error("prompt should omit depth limit when MaxDepth is 0")
	}
	if strings.Contains(prompt, "Concurrency Limit:") {
		t.Error("prompt should omit concurrency limit when MaxConcurrent is 0")
	}
	if strings.Contains(prompt, "Agent Limit:") {
		t.Error("prompt should omit agent limit when MaxAgents is 0")
	}
}

func TestBuildSystemPrompt_KeySections(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	sections := []string{
		"### Quine Process Channels",
		"### Environment",
		"### Active Constraints",
		"### Tools",
		"### Protocols",
	}

	for _, sec := range sections {
		if !strings.Contains(prompt, sec) {
			t.Errorf("prompt missing section %q", sec)
		}
	}
}

func TestBuildSystemPrompt_OpeningElevatesMissionAndFailure(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	checks := []string{
		"Complete your mission (argv)",
	}

	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing mission/failure text %q", want)
		}
	}
}

func TestBuildSystemPrompt_ExitToolDefinesTermination(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	checks := []string{
		"**exit** - Terminate explicitly",
		"Does NOT write to stdout",
		"emit deliverables via `sh` with `>&4`",
	}

	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing exit-tool text %q", want)
		}
	}
}

func TestBuildSystemPrompt_OpeningMentionsQuineRecursion(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	if !strings.Contains(prompt, "quine") {
		t.Error("prompt opening should identify the agent as quine")
	}
	if !strings.Contains(prompt, "POSIX") {
		t.Error("prompt should mention POSIX operating system")
	}
}

func TestBuildSystemPrompt_NoTripleBlankLines(t *testing.T) {
	withMaterial := BuildSystemPrompt(testConfig(), "test mission", true)
	withoutMaterial := BuildSystemPrompt(testConfig(), "test mission", false)

	if strings.Contains(withMaterial, "\n\n\n") {
		t.Error("prompt with material should not contain triple blank lines")
	}
	if strings.Contains(withoutMaterial, "\n\n\n") {
		t.Error("prompt without material should not contain triple blank lines")
	}
}

func TestBuildSystemPrompt_StdinAndForkPhysics(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	checks := []string{
		"`fd 3` is the quine process stdin",
		"finite batch or an open stream",
		"permanently consumes data",
		"truncated remainder",
		"`sh(command, stdin=\"...\")` is separate from material `fd 3`",
		"determined by quine at spawn time",
		"preserves process stdio at its current position",
		"binary mode (`-b`), quine snapshots stdin to a file before the loop",
	}

	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing stdin/fork physics text %q", want)
		}
	}
}

func TestBuildSystemPrompt_NoMaterialOmitsStdinSections(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", false)

	omitted := []string{
		"**Stdin Modes:**",
		"**Input Physics:**",
		"`fd 3` is the quine process stdin",
		"- fd 3: material stdin.",
		"unread material on `fd 3` stays unread",
	}

	for _, text := range omitted {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should omit no-material stdin text %q", text)
		}
	}

	if !strings.Contains(prompt, "`fd 3` is quine process stdin when material is provided") {
		t.Error("prompt should retain a minimal material principle in no-material sessions")
	}
	if !strings.Contains(prompt, "This session has no material stream and the user message is `Begin.`.") {
		t.Error("prompt should explicitly describe the no-material case")
	}
}

func TestBuildSystemPrompt_WithWisdom(t *testing.T) {
	cfg := testConfig()
	cfg.Wisdom = map[string]string{
		"SUMMARY": "User prefers concise answers",
		"CONTEXT": "Working on Go project",
	}
	prompt := BuildSystemPrompt(cfg, "test mission", true)

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
	cfg := testConfig()
	cfg.Wisdom = nil
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if strings.Contains(prompt, "### Wisdom") {
		t.Error("prompt should not contain wisdom section when wisdom is nil")
	}
}

func TestBuildSystemPrompt_WisdomSorted(t *testing.T) {
	cfg := testConfig()
	cfg.Wisdom = map[string]string{
		"ZEBRA":  "last alphabetically",
		"APPLE":  "first alphabetically",
		"MIDDLE": "in between",
	}
	prompt := BuildSystemPrompt(cfg, "test mission", true)

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

func TestBuildSystemPrompt_EscalationFastMode(t *testing.T) {
	cfg := testConfig()
	cfg.SmartModelID = "claude-opus"
	cfg.Escalated = false

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "Tier: Fast") {
		t.Error("prompt should contain 'Tier: Fast' when escalation is available")
	}
	if !strings.Contains(prompt, "**escalate**") {
		t.Error("prompt should contain escalate tool documentation")
	}
	if !strings.Contains(prompt, "`goal`/`strategy`") {
		t.Error("prompt should contain goal/strategy documentation in escalation mode")
	}
	if !strings.Contains(prompt, "WARNING STALL") {
		t.Error("prompt should contain stall warning protocol in escalation mode")
	}
}

func TestBuildSystemPrompt_EscalationSmartMode(t *testing.T) {
	cfg := testConfig()
	cfg.SmartModelID = "claude-opus"
	cfg.Escalated = true

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "Tier: Smart") {
		t.Error("prompt should contain 'Tier: Smart' when escalated")
	}
	if strings.Contains(prompt, "**escalate**") {
		t.Error("prompt should not contain escalate tool docs after escalation")
	}
	if strings.Contains(prompt, "goal=\"...\"") {
		t.Error("prompt should not contain goal/strategy docs after escalation")
	}
}

func TestBuildSystemPrompt_SingleModelMode(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if strings.Contains(prompt, "Tier:") {
		t.Error("prompt should not contain tier hints in single-model mode")
	}
	if strings.Contains(prompt, "**escalate**") {
		t.Error("prompt should not mention escalate tool docs in single-model mode")
	}
}
