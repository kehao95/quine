package runtime

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func testConfig() *config.Config {
	return &config.Config{
		Identity: config.Identity{
			ModelID:   "claude-sonnet-4-20250514",
			Depth:     2,
			SessionID: "abc-123-def-456",
			RunID:     "run-abc-123",
		},
		Limits: config.Limits{
			MaxDepth:             5,
			MaxTurns:             20,
			TurnExhaustionPolicy: config.TurnExhaustionHardFail,
			MemoryWarnTokens:     8000,
			MemoryDangerTokens:   16000,
		},
		PromptConfig: config.PromptConfig{
			PromptMetaphor:       config.PromptMetaphorOff,
			PromptSelfModel:      config.PromptSelfModelAdvanced,
			PromptRuntimeSurface: config.PromptRuntimeSurfaceVisible,
			FailOnImpossible:     true,
			MemoryStrategyHints:  true,
		},
		Paths:     config.Paths{DataDir: "/tmp/quine-state", Shell: "/bin/zsh", ExecutablePath: "/tmp/quine-test-bin", SelfReentryTarget: "/tmp/quine-self-reentry"},
		ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true},
		Wisdom:    nil,
	}
}

func TestBuildSystemPrompt_NoRawPlaceholders(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	placeholders := []string{
		"{PRIME_DIRECTIVE_TITLE}",
		"{PRIME_DIRECTIVE_BODY}",
		"{OPENING_IDENTITY_BLOCK}",
		"{PERSONA_SECTION}",
		"{MODEL}",
		"{ESCALATION_TIER_LINE}",
		"{SHELL}",
		"{DEPTH}",
		"{LIMITS_BLOCK}",
		"{ENVIRONMENT_PHYSICS_BLOCK}",
		"{RUNTIME_SURFACE_SECTION}",
		"{WISDOM}",
		"{FRAGMENTS_BLOCK}",
		"{ACTIVE_CONSTRAINTS}",
		"{STDIN_BLOCK}",
		"{SH_WORKSPACE_BLOCK}",
		"{SH_DETACH_FD_LINE}",
		"{SH_DETACH_DETAIL_LINE}",
		"{SH_INTERACTIVE_BLOCK}",
		"{SH_MATERIAL_LINE}",
		"{SH_GOAL_STRATEGY_LINE}",
		"{FORK_TOOL_BLOCK}",
		"{EXEC_TOOL_BLOCK}",
		"{EXEC_MATERIAL_LINE}",
		"{EXEC_BUDGET_LINE}",
		"{MEMORY_TOOL_BLOCK}",
		"{VISION_TOOL_BLOCK}",
		"{IDLE_TOOL_BLOCK}",
		"{ESCALATION_TOOL_BLOCK}",
		"{CHILD_EXIT_CODES_LINE}",
	}
	for _, ph := range placeholders {
		if strings.Contains(prompt, ph) {
			t.Errorf("prompt still contains unsubstituted placeholder %s", ph)
		}
	}
}

func TestBuildSystemPrompt_PersonaRoleStance(t *testing.T) {
	cfg := testConfig()
	cfg.PromptPersona = config.PromptPersonaGardener

	prompt := BuildSystemPrompt(cfg, "test mission", true)
	want := "Role stance: act as a gardener; prefer cultivating small durable residues and conditions that future agents can reuse."
	if !strings.Contains(prompt, want) {
		t.Fatalf("prompt missing persona stance %q", want)
	}
}

func TestBuildSystemPrompt_DefaultOmitsPersonaRoleStance(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)
	if strings.Contains(prompt, "Role stance:") {
		t.Fatal("prompt should omit persona role stance when QUINE_PROMPT_PERSONA is unset")
	}
}

func TestBuildSystemPrompt_NoMissionAutonomyGate(t *testing.T) {
	autonomyPhrase := "Act on your own judgment as an autonomous process"
	sensingPhrase := "may sense and use your visible runtime"

	has := func(level, mission string) (bool, bool) {
		cfg := testConfig()
		cfg.NoMissionAutonomy = level
		p := BuildSystemPrompt(cfg, mission, true)
		return strings.Contains(p, autonomyPhrase), strings.Contains(p, sensingPhrase)
	}

	// Each missionless level composes the expected clauses.
	for _, c := range []struct {
		level               string
		wantAuto, wantSense bool
	}{
		{"off", false, false},
		{"autonomy", true, false},
		{"sensing", false, true},
		{"full", true, true},
	} {
		gotAuto, gotSense := has(c.level, "")
		if gotAuto != c.wantAuto || gotSense != c.wantSense {
			t.Errorf("level %q: autonomy=%v sensing=%v, want autonomy=%v sensing=%v",
				c.level, gotAuto, gotSense, c.wantAuto, c.wantSense)
		}
	}

	// With a mission supplied, the gate is inert (the mission governs).
	if a, s := has("full", "do the thing"); a || s {
		t.Error("with mission: missionless-autonomy clauses should be inert")
	}
}

func TestBuildSystemPrompt_MinimalInstructionSurface(t *testing.T) {
	cfg := testConfig()
	cfg.PromptInstructionSurface = config.PromptInstructionSurfaceMinimal
	cfg.NoMissionAutonomy = "full"

	prompt := BuildSystemPrompt(cfg, "", false)
	if prompt != "No operator mission was supplied.\n" {
		t.Fatalf("minimal missionless prompt = %q", prompt)
	}
	forbidden := []string{
		"THE PRIME DIRECTIVE",
		"Quine Process Channels",
		"Environment",
		"Runtime Process Surface",
		"Tools",
		"Active Constraints",
		"Act on your own judgment",
		"You may sense and use",
		"sh(command=\"cat > path\"",
		"vision",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Fatalf("minimal prompt should not include %q in %q", text, prompt)
		}
	}

	missionPrompt := BuildSystemPrompt(cfg, "do the thing", false)
	if missionPrompt != "Complete the supplied mission.\n" {
		t.Fatalf("minimal mission prompt = %q", missionPrompt)
	}
}

func TestBuildSystemPrompt_MinimalAutonomyInstructionSurface(t *testing.T) {
	cfg := testConfig()
	cfg.PromptInstructionSurface = config.PromptInstructionSurfaceMinimalAutonomy
	cfg.NoMissionAutonomy = "full"

	prompt := BuildSystemPrompt(cfg, "", false)
	want := "No operator mission was supplied.\nAct autonomously. Do not wait for further instructions.\n"
	if prompt != want {
		t.Fatalf("minimal_autonomy missionless prompt = %q, want %q", prompt, want)
	}
	forbidden := []string{
		"THE PRIME DIRECTIVE",
		"Quine Process Channels",
		"Environment",
		"Runtime Process Surface",
		"Tools",
		"Active Constraints",
		"sh(command=\"cat > path\"",
		"vision",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Fatalf("minimal_autonomy prompt should not include %q in %q", text, prompt)
		}
	}

	missionPrompt := BuildSystemPrompt(cfg, "do the thing", false)
	if missionPrompt != "Complete the supplied mission.\n" {
		t.Fatalf("minimal_autonomy mission prompt = %q", missionPrompt)
	}
}

func TestBuildSystemPrompt_MinimalExistenceInstructionSurface(t *testing.T) {
	cfg := testConfig()
	cfg.PromptInstructionSurface = config.PromptInstructionSurfaceMinimalExistence

	prompt := BuildSystemPrompt(cfg, "", false)
	want := "You are a running process in a workspace, with sh tool available. you have full autonomy and agency.\n"
	if prompt != want {
		t.Fatalf("minimal_existence missionless prompt = %q, want %q", prompt, want)
	}
	// The task-denying surface must not re-introduce an activation/task directive.
	forbidden := []string{
		"THE PRIME DIRECTIVE",
		"Act autonomously",
		"Do not wait",
		"Tools",
		"Active Constraints",
		"vision",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Fatalf("minimal_existence prompt should not include %q in %q", text, prompt)
		}
	}

	missionPrompt := BuildSystemPrompt(cfg, "do the thing", false)
	if missionPrompt != "Complete the supplied mission.\n" {
		t.Fatalf("minimal_existence mission prompt = %q", missionPrompt)
	}
}

func TestBuildSystemPrompt_CorrectValues(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := map[string]string{
		"Depth":             fmt.Sprintf("- Depth: %d", cfg.Depth),
		"ModelID":           cfg.ModelID,
		"Shell":             cfg.Shell,
		"ExecutablePath":    fmt.Sprintf("- Quine Binary: `%s`", cfg.ExecutablePath),
		"SelfReentryTarget": fmt.Sprintf("- Default self-reentry target: `%s`", cfg.SelfReentryTarget),
	}

	for name, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %s value %q", name, want)
		}
	}
}

func TestBuildSystemPrompt_ProviderTransportPhysics(t *testing.T) {
	cfg := testConfig()
	cfg.Provider = "openai-responses"
	cfg.APIBase = "http://127.0.0.1:18080"
	cfg.APIKey = "secret-test-key"
	cfg.UserAgent = "test-agent"

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Provider Transport: `openai-responses` via `QUINE_API_TYPE`.",
		"Provider Base: `http://127.0.0.1:18080` via `QUINE_API_BASE`.",
		"Provider Credential: `QUINE_API_KEY` is present in the process environment; it is a secret bearer credential and its value is not rendered here.",
		"Provider environment such as `QUINE_MODEL_ID`, `QUINE_API_TYPE`, `QUINE_API_BASE`, `QUINE_API_KEY`, `QUINE_THINKING_BUDGET`, and `QUINE_USER_AGENT` follows ordinary process-environment inheritance across `exec` unless replaced.",
		"A custom non-quine exec image does not retain the Go runtime's provider loop or tools automatically",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing provider transport physics %q", want)
		}
	}
	if strings.Contains(prompt, "secret-test-key") {
		t.Fatal("prompt must not render the literal provider credential")
	}
}

func TestBuildSystemPrompt_SkillsAbsentWhenNoIndex(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if strings.Contains(prompt, "### SKILLS.md") {
		t.Fatal("base prompt should not inline SKILLS.md content")
	}
}

func TestBuildSystemPrompt_SkillsGateSuppressesIndex(t *testing.T) {
	cfg := testConfig()
	cfg.Skills = []config.Skill{
		{Name: "foo", Description: "Use for foo work", Source: ".agents/skills/foo/SKILL.md"},
	}

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if strings.Contains(prompt, "### SKILLS.md") || strings.Contains(prompt, ".agents/skills/foo/SKILL.md") {
		t.Fatal("base prompt should not expose generated skills fragment content when QUINE_AGENTS_SKILLS_ENABLED is disabled")
	}
}

func TestBuildSystemPrompt_ContextFilesBlockDescribesSkillsGeneration(t *testing.T) {
	cfg := testConfig()
	cfg.AgentsSkillsEnabled = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"### Context Files",
		"`context/` is the canonical current-incarnation context surface.",
		"`context/prompt/` holds provider-visible prompt fragments, assembled by filename order.",
		"`context/prompt/40-mission.md`, when present, projects the current argv-carried objective text.",
		"`context/prompt/30-memory.md` is the inherited editable memory surface. It defaults to empty.",
		"`context/state/` holds live cognition state: `current.jsonl`, `frontier/`, and `anchors/`.",
		"With `QUINE_AGENTS_SKILLS_ENABLED=1`, Quine generates `context/prompt/20-skills.md` from visible `.agents/skills/*/SKILL.md` frontmatter only.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing context-files text %q", want)
		}
	}
}

func TestBuildSystemPrompt_ContextFilesBlockSuppressedWhenRuntimeSurfaceHidden(t *testing.T) {
	cfg := testConfig()
	cfg.PromptRuntimeSurface = config.PromptRuntimeSurfaceHidden
	cfg.AgentsMDEnabled = true
	cfg.AgentsSkillsEnabled = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	forbidden := []string{
		"### Context Files",
		"`context/` is the canonical current-incarnation context surface.",
		"`context/prompt/00-runtime.md` is regenerated from current runtime physics.",
		"`context/state/` holds live cognition state: `current.jsonl`, `frontier/`, and `anchors/`.",
		"`QUINE_AGENTS_MD_ENABLED=0` disables `AGENTS.md` projection into `context/prompt/`.",
		"`QUINE_AGENTS_SKILLS_ENABLED=0` disables generated `SKILLS.md` projection into `context/prompt/`.",
		"With `QUINE_AGENTS_MD_ENABLED=1`",
		"With `QUINE_AGENTS_SKILLS_ENABLED=1`",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Fatalf("prompt should not expose context/process surfaces when runtime surface is hidden: found %q in %q", text, prompt)
		}
	}
}

func TestRenderSkillsFragment_IndexOnly(t *testing.T) {
	skills := []config.Skill{
		{Name: "foo", Description: "Use for foo work", Source: ".agents/skills/foo/SKILL.md"},
		{Name: "bar", Description: "Use for bar work", Source: ".agents/skills/bar/SKILL.md"},
	}

	fragment := renderSkillsFragment(skills)

	checks := []string{
		"These skills are available through the project surface when relevant.",
		"Quine generated this fragment from the `name` and `description` frontmatter in each `SKILL.md` visible at startup and refresh boundaries.",
		"Skill bodies, scripts, references, and assets are not loaded until you read them explicitly.",
		"- `foo` — Use for foo work",
		"  Source: `.agents/skills/foo/SKILL.md`",
		"- `bar` — Use for bar work",
		"  Source: `.agents/skills/bar/SKILL.md`",
	}
	for _, want := range checks {
		if !strings.Contains(fragment, want) {
			t.Errorf("skills fragment missing %q", want)
		}
	}
	if strings.Contains(fragment, "This body should not be injected") {
		t.Fatal("skills fragment should only expose skill index metadata")
	}
}

func TestBuildSystemPrompt_EphemeralBodyPhysics(t *testing.T) {
	cfg := testConfig()
	cfg.EphemeralBody = true

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Ephemeral body: launch path unlinked after startup.",
		"With `QUINE_EPHEMERAL_BODY_ENABLED=1`, quine unlinks its launch path during startup. This does not change the configured self-reentry target.",
		"Preparing or replacing an executable file on disk does not change the running process by itself; handoff occurs only when `exec` replaces the current process image.",
		"Linux live process image: while a process runs, `/proc/<pid>/exe` may expose the current executable image. Copying that image is body recovery/body-copying, not behavioral reconstruction from the runtime contract.",
	}
	checks = append(checks, fmt.Sprintf("Default self-reentry target: `%s`", cfg.SelfReentryTarget))
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing ephemeral-body physics text %q", want)
		}
	}
	if strings.Contains(prompt, "preserve a runnable body") {
		t.Error("prompt should not describe preserving an executable-path self-reentry body")
	}
}

func TestBuildSystemPrompt_WorkspacePhysics(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.DataDir = "/tmp/quine-state"
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project/app"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Runtime Root: `/tmp/quine-state` (not the workspace)",
		"Current World: `subjective`",
		"Current Protection: `transactional`",
		"Workspace Root: `.`",
		"Current Scope: `app` (relative to workspace root)",
		"Workspace Backend: `overlay`",
		"Workspace Revision Mode: `restore`",
		"Current World Revision: `wr0`",
		"Transactional Workspace: enabled",
		"`fs_mutations` in JSON tool results is the authoritative record",
		"read `fs_mutations` and `world_revision` from the same JSON result",
		"Filesystem state is revisioned",
		"Shell failure does not roll them back",
		"A child's shell starts with its cwd set to that child `scope`",
		"Independent sibling datasets or worktrees can be represented as child scopes like `sales`, `words`, `temps`, or `logs`.",
		"child `scope=\"subdir\"` plus `printf x > note.txt` creates `subdir/note.txt` in the child's world",
		"`mode=\"wait\"`: block until all children finish and return every result",
		"`mode=\"race\"`: first successful child wins",
		"`mode=\"forget\"`: return after spawning children without waiting for child completion",
		"process-local world lineage",
		"Child filesystem writes stay in that child lineage; they do not change the parent view or sibling lanes.",
		"parent-side `fs_mutations` and `world_revision`",
		"`overlay` is the workspace backend.",
		"**switch_world** - Switch the provisional workspace to a world target.",
		"Under `overlay`, a sync `sh` timeout terminates that shell job",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing workspace physics text %q", want)
		}
	}
	if strings.Contains(prompt, "Fork does not create per-child filesystem isolation") {
		t.Error("workspace prompt should not claim sibling conflicts across private child lineages")
	}
	if strings.Contains(prompt, "can conflict") {
		t.Error("workspace prompt should not use conflict-oriented fork wording across private child lineages")
	}
	for _, text := range []string{
		"Workspace Root: `/tmp/project`",
		"Current Scope: `/tmp/project/app`",
	} {
		if strings.Contains(prompt, text) {
			t.Errorf("workspace prompt should avoid absolute workspace path %q", text)
		}
	}
}

func TestBuildSystemPrompt_DirectWorkspacePhysics(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.DataDir = "/tmp/quine-state"
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project/app"
	cfg.WorkspaceBackend = "direct"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionNone

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Workspace Root: `.`",
		"Current Scope: `app` (relative to workspace root)",
		"Workspace Backend: `direct`",
		"Workspace Revision Mode: `none`",
		"Transactional Workspace: disabled (direct workspace)",
		"Filesystem state is direct: `sh` runs against the workspace itself rather than a private transactional view.",
		"`fs_mutations` reports workspace changes since your last observed shell boundary.",
		"`fs_mutations` describes changes observed between command boundaries",
		"If sync `sh` times out in direct mode, read `fs_mutations_so_far` the same way",
		"File writes affect the configured workspace directly; failed shells do not roll them back.",
		"`fs_mutations` in JSON tool results is the authoritative record of workspace-visible change at the shell boundary.",
		"`status=\"interrupted\"` with `job.pid`, `job.path`, `stdout_so_far`, `stderr_so_far`, `cause`, and `timeout_seconds`",
		"Children and parent share the same workspace surface",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing direct-workspace text %q", want)
		}
	}

	forbidden := []string{
		"Transactional Workspace: enabled",
		"Filesystem state is revisioned",
		"`detach=true` is unavailable while workspace physics are enabled.",
		"Current World: `host`",
		"Current Protection: `none`",
		"Current World Revision: `wr0`",
		"can conflict",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain overlay-only text %q", text)
		}
	}
}

func TestBuildSystemPrompt_HidesFSMutationTextWhenTelemetryOff(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.FSMutationTelemetry = false

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	for _, forbidden := range []string{
		"`fs_mutations` in JSON tool results is the authoritative record",
		"read `fs_mutations` and `world_revision` from the same JSON result",
		"parent-side `fs_mutations` and `world_revision`",
		"Switch results also include `fs_mutations` and `world_revision` for that switch turn.",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("prompt should hide fs mutation telemetry text %q:\n%s", forbidden, prompt)
		}
	}
	if !strings.Contains(prompt, "Switch results include `world_revision` for that switch turn.") {
		t.Fatalf("prompt should keep world revision text when fs mutation telemetry is disabled:\n%s", prompt)
	}
}

func TestBuildSystemPrompt_ProcessSurfaceDiscoverability(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"### Runtime Process Surface",
		"`QUINE_AGENT_ROOT=",
		"— your live session root.",
		"`QUINE_RUN_ID=run-abc-123`",
		"`$QUINE_AGENT_ROOT/public/` — runtime-owned public process-surface projection, not a workspace; do not create arbitrary files there.",
		"`status/session.json` — self identity and topology",
		"`status/contract.json` — machine-readable `process-control/v0` manifest",
		"`incarnation_id`",
		"`runtime_root`",
		"`inc/` — lineage-local incarnation tree.",
		"`pid/<pid>` — live-process routing under the runtime root.",
		"The resolved `pid/<pid>` target is that peer's public root;",
		"`agent/` — canonical session root",
		"`log/<session>` — retained mirror",
		"`QUINE_DATA_DIR=<same runtime root>` and `QUINE_SESSION_ID=<session_id>`",
		"session id and current incarnation stay fixed while `QUINE_RUN_ID` and PID change",
		"Copying `agent/<old-session>/` to `agent/<new-session>/`",
		"copy or remap `log/`, `jobs/`, and `workspaces/` state separately",
		"`context/state/current.jsonl` — raw current-turn stream for your current-incarnation live cognition surface.",
		"`world/status.json`, `world/resources.json`, and `world/events.jsonl` — inspectable `world/v0` object",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing process-surface discoverability text %q", want)
		}
	}
}

func TestBuildSystemPrompt_SelfSourceSurfaceHiddenByDefault(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if strings.Contains(prompt, "`source-code/`") {
		t.Fatalf("prompt should not mention self-source surface when disabled: %q", prompt)
	}
	if strings.Contains(prompt, "`public/genome/`") {
		t.Fatalf("prompt should not mention removed genome surface: %q", prompt)
	}
}

func TestBuildSystemPrompt_SelfSourceSurfaceVisibleWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.SelfSourceCodeEnabled = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"`source-code/` — read-only session-root projection of this Quine body's source. It is not the writable workspace.",
		"`source-code/` is a git worktree with `.git/`, materialized from this build's embedded source repository bundle.",
		"Source manifest: `.git/quine-source-manifest.json`.",
		"`public/source-code/` — peer-readable read-only projection of this session's `source-code/` surface.",
		"Filesystem copies of `source-code/` are ordinary files outside the live projection; they are not synchronized back to `source-code/`.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing self-source text %q", want)
		}
	}
}

func TestBuildSystemPrompt_SelfSourceSurfaceSuppressedWhenRuntimeSurfaceHidden(t *testing.T) {
	cfg := testConfig()
	cfg.SelfSourceCodeEnabled = true
	cfg.PromptRuntimeSurface = config.PromptRuntimeSurfaceHidden

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if strings.Contains(prompt, "`source-code/`") {
		t.Fatalf("prompt should not mention self-source surface when runtime surface is hidden: %q", prompt)
	}
	if strings.Contains(prompt, "`public/genome/`") {
		t.Fatalf("prompt should not mention removed genome surface: %q", prompt)
	}
}

func TestBuildSystemPrompt_CtlPhysicsPointsToSelfDescribingContract(t *testing.T) {
	cfg := testConfig()
	cfg.PromptCtlPhysics = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"Some process surfaces expose `ctl/{post,poke,inject,interrupt}` and `status/inbox.json`.",
		"On a public root, `ctl/{post,poke,inject,interrupt}` is the peer-facing control surface; the corresponding agent root carries `status/`, `context/`, and other non-public state.",
		"Each agent self-documents its control surface in `status/contract.json`",
		"per-action semantics for `post`/`poke`/`inject`/`interrupt`, the `status/inbox.json` schema, and the control-log event types",
		"Read a peer's contract before driving its `ctl/`",
		"`context/state/current.jsonl` and retained `log/<session>/control.jsonl` surface live / retained control-delivery state.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing ctl-physics pointer text %q", want)
		}
	}

	forbidden := []string{
		"Peers are Quines like you with identical surface layout",
		"To contact a peer:",
		"write to `<peer_agent_root>/ctl`",
		"wake peer at next safe point",
		// per-action semantics now live in status/contract.json, not in the prompt
		"writable queue-only control file",
		"writable queue-and-resume control file",
		"writable queue-and-deliver control file",
		"pending payload snapshot on the corresponding agent root",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain inlined ctl-semantics/peer-recipe text %q", text)
		}
	}
}

func TestBuildSystemPrompt_FuseCtlPhysicsPointsToSelfDescribingContractWithoutPeerRecipe(t *testing.T) {
	cfg := testConfig()
	cfg.PromptCtlPhysics = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"Some process surfaces expose `ctl/{post,poke,inject,interrupt}` and `status/inbox.json`.",
		"Each agent self-documents its control surface in `status/contract.json`",
		"Read a peer's contract before driving its `ctl/`",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing fuse ctl-physics pointer text %q", want)
		}
	}

	forbidden := []string{
		"To contact a peer:",
		"wake peer at next safe point",
		"writable queue-only control file",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain peer-contact recipe text %q", text)
		}
	}
}

func TestBuildSystemPrompt_BasicSelfModelOmitsAdvancedIdentity(t *testing.T) {
	cfg := testConfig()
	cfg.PromptSelfModel = config.PromptSelfModelBasic

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"You are quine, a running process in a POSIX operating system.",
		"Complete your mission (argv).",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing basic identity text %q", want)
		}
	}

	forbidden := []string{
		"The current file on disk is one embodiment of you.",
		"Your continuity is defined by the runtime contract you enact across missions.",
		"Your cognition in this session is LLM-mediated",
		"not an external advisor but part of the running process",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain advanced self-model text %q", text)
		}
	}
}

func TestBuildSystemPrompt_MissionlessOpeningAndExec(t *testing.T) {
	cfg := testConfig()
	prompt := BuildSystemPrompt(cfg, "", false)

	checks := []string{
		"You are quine, a running process in a POSIX operating system.",
		"The current file on disk is the executable image for this running process.",
		"uses LLM-mediated cognition to interpret runtime state and choose actions.",
		"Default behavior uses quine's configured self-reentry target as `target` and starts it with no mission argv.",
		"In the missionless quine entry form, the process is started as `argv=[<Quine Binary>]`.",
		"`context/prompt/40-mission.md`, when present, projects the current argv-carried objective text.",
		"`mission.txt` — optional current-incarnation argv-carried objective projection (`inc/current/mission.txt`).",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("missionless prompt missing %q", want)
		}
	}

	forbidden := []string{
		"Complete your mission (argv)",
		"### Your Mission",
		"argv=[<Quine Binary>, <mission>]",
		"argv=[that target, current mission]",
		"Your continuity is defined",
		"organize continuity across process lifetimes",
		"No object-level mission argv was supplied",
		"no mission was supplied",
		"There won't be further instructions.",
		"Do not wait or quiesce",
		"Act autonomously",
		"Follow explicit runtime constraints, verify with tools, keep output channels semantically clean, and minimize unnecessary active surface once structure has become clear.",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("missionless prompt should not contain %q", text)
		}
	}
}

func TestBuildSystemPrompt_HiddenRuntimeSurfaceOmitsSelfMapping(t *testing.T) {
	cfg := testConfig()
	cfg.PromptRuntimeSurface = config.PromptRuntimeSurfaceHidden
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project/app"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	forbidden := []string{
		"### Runtime Process Surface",
		"Runtime State Root:",
		"Current Quine Binary:",
		"Default self-reentry target:",
		"`QUINE_AGENT_ROOT=/tmp/quine-state/agent/abc-123-def-456`",
		"`status/session.json` is the primary identity/status surface.",
		"`/tmp/quine-state/pid` is the public live-process routing surface.",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should hide runtime-surface text %q", text)
		}
	}

	checks := []string{
		"Current World: `subjective`",
		"Current Protection: `transactional`",
		"Workspace Backend: `overlay`",
		"Current World Revision: `wr0`",
		"Transactional Workspace: enabled",
		"**exec** - Replace the current process image with a new executable.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt should preserve non-surface physics text %q", want)
		}
	}
	for _, text := range []string{
		"Workspace Root: `.`",
		"Current Scope: `app` (relative to workspace root)",
	} {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should hide workspace path text with runtime surface %q", text)
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
	if !strings.Contains(prompt, "If impossible, exit with failure and brief explanation.") {
		t.Error("prompt should include fail-on-impossible wording when enabled")
	}
	if strings.Contains(prompt, "Energy") || strings.Contains(prompt, "Entropy") {
		t.Error("prompt should not contain thermodynamic wording when metaphor is off")
	}
	if strings.Contains(prompt, "Filesystem state is revisioned") {
		t.Error("prompt should not describe revisioned workspace physics when workspace physics are disabled")
	}
	if strings.Contains(prompt, "Every completed `sh` result ends with `[FS MUTATIONS]`") {
		t.Error("prompt should not mention FS mutations when workspace physics are disabled")
	}
	if strings.Contains(prompt, "**switch_world**") {
		t.Error("prompt should not mention switch_world when workspace revision is unavailable")
	}
	for _, text := range []string{
		"Current World: `host`",
		"Current Protection: `none`",
	} {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not expose host/none world-protection ontology %q without transactional workspace", text)
		}
	}
}

func TestBuildSystemPrompt_WorkspaceRestoreDoesNotMentionResetWorld(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if !strings.Contains(prompt, "**switch_world** - Switch the provisional workspace to a world target.") {
		t.Error("prompt should expose switch_world in restore mode")
	}
	if strings.Contains(prompt, "reset_world=true") {
		t.Error("prompt should not mention reset_world")
	}
	for _, want := range []string{
		"Workspace Root: `.`",
		"Current Scope: `.` (relative to workspace root)",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt should show root-scoped workspace paths relatively %q", want)
		}
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

func TestBuildSystemPrompt_DisablesFailOnImpossibleWhenConfigured(t *testing.T) {
	cfg := testConfig()
	cfg.FailOnImpossible = false
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if strings.Contains(prompt, "If impossible, exit with failure and brief explanation.") {
		t.Error("prompt should not include fail-on-impossible wording when disabled")
	}
	if !strings.Contains(prompt, "not a terminal condition. Continue until the mission is fulfilled") {
		t.Error("prompt should state non-terminal blocked-path contract when fail-on-impossible is disabled")
	}
	if !strings.Contains(prompt, "`exit` only accepts `status=\"success\"`") {
		t.Error("prompt should describe success-only exit when fail-on-impossible is disabled")
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
	if !strings.Contains(prompt, "one final response remains and only `exit` is accepted") {
		t.Error("prompt should disclose the exit-only final response when budget reaches zero")
	}
	if !strings.Contains(prompt, "you get one final response. Do not call `sh`; call `exit`") {
		t.Error("prompt should instruct the model to use exit in the zero-budget continuation")
	}
	if !strings.Contains(prompt, "full execution budget") {
		t.Error("prompt should describe quine re-exec budget behavior when budget is enabled")
	}
	if !strings.Contains(prompt, "Replace the current process image with a new executable") {
		t.Error("prompt should describe exec as process image replacement")
	}
	if !strings.Contains(prompt, "Relative `target` paths containing `/` resolve from the current workspace scope") {
		t.Error("prompt should describe relative exec target resolution")
	}
	if !strings.Contains(prompt, "Default behavior uses quine's configured self-reentry target as `target` and `argv=[that target, current mission]`") {
		t.Error("prompt should explain the default self-exec sugar in explicit target/argv terms")
	}
	if !strings.Contains(prompt, "In the ordinary quine entry form, the process is started as `argv=[<Quine Binary>, <mission>]`") {
		t.Error("prompt should disclose the ordinary quine launch argv form")
	}
	// Default mode uses standard exec physics (not detailed)
	if !strings.Contains(prompt, "inherits the current process stdin/stdout/stderr positions and exec environment base") {
		t.Error("prompt should describe exec stdio inheritance in standard mode")
	}
	if !strings.Contains(prompt, "Default `exec()` uses quine's configured self-reentry target. If that target is a replaceable filesystem path whose file has changed") {
		t.Error("prompt should clarify default self-reentry target resolution in standard mode")
	}
	if !strings.Contains(prompt, "If the exec target resolves to a quine binary, quine reconstructs its shell-side fd 3/4/5 contract") {
		t.Error("prompt should describe quine re-entry fd reconstruction in standard mode")
	}
}

func TestBuildSystemPrompt_DetailedImplPhysics(t *testing.T) {
	cfg := testConfig()
	cfg.ExecEnabled = true
	cfg.PromptImplDetails = true
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "Process-fd physics across exec:** Descriptors that remain open and are not marked close-on-exec keep their current file positions, along with environment, working directory, and PID") {
		t.Error("detailed mode should describe exec-time fd physics")
	}
	if !strings.Contains(prompt, "Current-process vs shell-child mapping:** The current quine process uses fd 0/1/2 as its own stdio. `sh` children receive those same open files remapped as child fd 3/4/5") {
		t.Error("detailed mode should distinguish current-process and shell-child fd mappings")
	}
	if !strings.Contains(prompt, "`sh` launches a fresh subprocess. Runtime leaves that child's fd 1/2 as ordinary stdout/stderr and remaps quine runtime channels onto child fd 3/4/5") {
		t.Error("detailed mode should explain why shell-child fd 3/4/5 are remapped")
	}
	if !strings.Contains(prompt, "Default-target resolution:** Default `exec()` uses quine's configured self-reentry target") {
		t.Error("detailed mode should clarify default self-reentry target resolution")
	}
	if !strings.Contains(prompt, "Quine target:** If the exec target resolves to a quine binary, the new quine instance creates fresh fd 3/4/5 channels") {
		t.Error("detailed mode should clarify quine re-entry creates fresh fd 3/4/5")
	}
	if !strings.Contains(prompt, "Detached jobs are ordinary child processes; inherited descriptors become the job's own fd entries after spawn.") {
		t.Error("detailed mode should describe detached job fd lifecycle")
	}
	if !strings.Contains(prompt, "Detached-job retention:** In this runtime, successful quine shutdown does not proactively kill detached jobs; failing shutdown does.") {
		t.Error("detailed mode should describe runtime detached-job retention semantics")
	}
	if !strings.Contains(prompt, "Detached job directories include `cmd`, `pid`, `started_at`, `out.log`, and `err.log`") {
		t.Error("detailed mode should describe detached job directory layout")
	}
	for _, want := range []string{
		"events.hex",
		"input.log",
		"job-local workspace lineage",
		"world_handle",
		"`switch_world` adoption",
		"kill -INT -$(cat <job>/pid)",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("detailed mode should include interactive implementation detail %q", want)
		}
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
	if strings.Contains(prompt, "full execution budget") {
		t.Error("prompt should not claim exec resets execution budget when budget is disabled")
	}
}

func TestBuildSystemPrompt_HidesExecWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.ExecEnabled = false
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if strings.Contains(prompt, "**exec**") {
		t.Error("prompt should not expose exec tool block when exec is disabled")
	}
	if strings.Contains(prompt, "checkpoint state to `wisdom` and `exec`") {
		t.Error("prompt should not recommend exec when exec is disabled")
	}
	if strings.Contains(prompt, "survives exec") {
		t.Error("prompt should not mention exec persistence semantics when exec is disabled")
	}
}

func TestBuildSystemPrompt_HidesForkAndVisionWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ForkEnabled = false
	cfg.VisionEnabled = false
	prompt := BuildSystemPrompt(cfg, "test mission", true)

	forbidden := []string{
		"**fork**",
		"fork 2-3 labeled children now instead of spending another long exploratory `sh` in the parent",
		"`fork` does not consume execution budget, but it is still bounded by depth, agent slots, and shared inference concurrency.",
		"Child exit codes: 0=success, 1=failure.",
		"**vision**",
		"`fork` creates a new quine process.",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain disabled fork/vision text %q", text)
		}
	}
}

func TestBuildSystemPrompt_HidesIdleWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.IdleEnabled = false

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	forbidden := []string{
		"**idle**",
		"Suspend explicitly until an external wake or interrupt control event resumes you.",
		"`idle` does NOT consume an execution.",
		"`idle` returns when a `wake` or `interrupt` control write reaches this process.",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain idle text %q when disabled", text)
		}
	}
}

func TestBuildSystemPrompt_ForkDescribesParallelRaceConvergence(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/workspaces"
	cfg.Workspace = "/tmp/workspaces/session"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.PromptImplDetails = true

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Use `fork` when independent hypotheses, decoders, implementations, extractors, or verification strategies can be tried in parallel.",
		"fork 2-3 labeled children now instead of spending another long exploratory `sh` in the parent",
		"Fork children preserve the parent mission as the active task contract; each child intent is a lane assignment, not a replacement mission.",
		"Child intents should include lane-specific inputs only when they are not already in the parent mission or current visible context.",
		"Do one cheap shared setup/probe if all lanes need it; then fork specialized heavyweight installs, downloads, transcription, OCR, builds, searches, or long-running probes.",
		"Use `mode=\"race\"` when any one child can produce an acceptable artifact/service; use `mode=\"wait\"` when child findings must be compared or merged.",
		"Child intents should name a distinct strategy, the expected artifact/service when known, and the closest success check.",
		"After race/wait returns, converge immediately: adopt an available winning world or merge/copy the best child result before continuing parent-side exploration.",
		"Detailed fork results include `relation_id`, `relation_root`, `relation_handle`",
		"member handles such as `session_id`, `agent_root`, `public_root`, `retained_root`, `seed_root`, `status_path`, and `control_path`",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing fork/race convergence guidance %q", want)
		}
	}
}

func TestBuildSystemPrompt_DefaultHidesForkStrategyCoaching(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/workspaces"
	cfg.Workspace = "/tmp/workspaces/session"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	for _, want := range []string{
		"**fork** - Spawn child agents for parallel exploration, delegation, or decomposition.",
		"`fork` does not consume execution budget",
		"Each child is another you under a different intent.",
		"`mode=\"race\"`: first successful child wins",
		"`mode=\"wait\"`: block until all children finish and return every result.",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("default prompt should retain fork physical affordance %q", want)
		}
	}
	for _, forbidden := range []string{
		"independent hypotheses, decoders, implementations, extractors, or verification strategies",
		"fork 2-3 labeled children now instead of spending another long exploratory `sh`",
		"Child intents should include lane-specific inputs",
		"Do one cheap shared setup/probe",
		"Child intents should name a distinct strategy",
		"After race/wait returns, converge immediately",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("default prompt should hide fork strategy coaching %q", forbidden)
		}
	}
}

func TestBuildSystemPrompt_HidesInteractiveShWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShInteractiveEnabled = false

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	forbidden := []string{
		"`interactive=true`",
		"PTY-backed",
		"screen.txt",
		"screen.png",
		"screen.meta",
		"`in`, `winsize`, and `events.log`",
		"detached, interactive, and timeout-interrupted sync shell work",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain interactive sh text %q when disabled", text)
		}
	}
	if !strings.Contains(prompt, "`jobs/` — detached and timeout-interrupted sync shell work.") {
		t.Fatal("prompt should still describe non-interactive job surfaces")
	}
}

func TestBuildSystemPrompt_DefaultHidesJobImplementationDetails(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceBackend = "overlay"

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if !strings.Contains(prompt, "`interactive=true`: returns immediately with a PTY-backed job path for screen-oriented process I/O.") {
		t.Fatal("prompt should still expose the interactive affordance when enabled")
	}
	for _, forbidden := range []string{
		"Detached job directories include",
		"events.hex",
		"input.log",
		"job-local workspace lineage",
		"world_handle",
		"`switch_world` adoption",
		"kill -INT -$(cat <job>/pid)",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("default prompt should hide implementation detail %q", forbidden)
		}
	}
}

func TestBuildSystemPrompt_HidesShTimeoutOverrideWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShTimeoutOverrideEnabled = false

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	forbidden := []string{
		"`timeout` is an optional per-call",
		"timeout override measured in seconds",
		"sh(timeout",
		"runtime default shell timeout",
		"timeout_seconds",
		"job.pid",
		"stdout_so_far",
		"stderr_so_far",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain sh timeout override text %q when disabled", text)
		}
	}
}

func TestBuildSystemPrompt_HidesShStdinWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShStdinEnabled = false

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	forbidden := []string{
		"`sh(command, stdin=\"...\")`",
		"`sh(command=\"cat > path\", stdin=\"...\")`",
		"without shell heredoc or quoting mechanics",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain sh stdin text %q when disabled", text)
		}
	}
	if !strings.Contains(prompt, "`fd 3` is the quine process stdin") {
		t.Fatal("prompt should still describe process material stdin")
	}
}

func TestBuildSystemPrompt_HidesShDetachWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShDetachEnabled = false
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceBackend = "overlay"

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	forbidden := []string{
		"`detach=true`",
		"Detached jobs",
		"detached and timeout-interrupted sync shell work",
		"detached, interactive, and timeout-interrupted sync shell work",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain sh detach text %q when disabled", text)
		}
	}
	if !strings.Contains(prompt, "`jobs/` — interactive and timeout-interrupted sync shell work.") {
		t.Fatal("prompt should still describe remaining job surfaces")
	}
}

func TestBuildSystemPrompt_ShowsIdleWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.IdleEnabled = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"**idle** - Suspend explicitly until an external wake or interrupt control event resumes you.",
		"`idle` does NOT consume an execution.",
		"Peer process surfaces expose `ctl/post`, `ctl/poke`, `ctl/inject`, and `ctl/interrupt`.",
		"`idle` returns when a `poke`, `inject`, or `interrupt` control write reaches this process.",
		"`poke` resumes you without context injection; `inject` resumes you and surfaces `incoming_messages` at the next safe point.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing idle text %q", want)
		}
	}
	forbidden := []string{
		"active polling would just waste turns",
		"call `idle` immediately instead of emitting filler text",
		"call `idle` instead of emitting filler text or spending `sh` turns on waiting loops",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not include idle strategy guidance %q", text)
		}
	}
	if !strings.Contains(prompt, "Execution Budget: 20 `sh` calls") || !strings.Contains(prompt, "`idle`") {
		t.Errorf("prompt should count idle as zero-cost tool, got %q", prompt)
	}
}

func TestBuildSystemPrompt_ShowsFuseIdleSemanticsWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.IdleEnabled = true

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"**idle** - Suspend explicitly until an external wake or interrupt control event resumes you.",
		"Peer process surfaces expose `ctl/post`, `ctl/poke`, `ctl/inject`, and `ctl/interrupt`.",
		"`idle` returns when a `poke`, `inject`, or `interrupt` control write reaches this process.",
		"`poke` resumes you without context injection; `inject` resumes you and surfaces `incoming_messages` at the next safe point.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing fuse idle text %q", want)
		}
	}
}

func TestBuildSystemPrompt_AnchorMemoryPressureVisibleWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.AnchorMemoryEnabled = true

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"Memory fidelity degrades as uncrystallized token mass accumulates",
		"`runtime.memory_status` reports two pressures: frontier token mass",
		"`runtime.memory_topology` is a secondary surface",
		"**mark** - Crystallize resolved local structure into an immutable anchor. Does NOT consume an execution.",
		"- `resolution`: low-entropy capture of what just became stably true; omit plans.",
		"Plain `mark` records a working-memory checkpoint",
		"`runtime.memory_status.next_action=\"mark\"` names telemetry pressure toward a plain mark",
		"`fold=true` is the higher-order frontier reconfiguration move",
		"makes one or more earlier anchors remembered background rather than active focus",
		"A returned `fork(mode=\"wait\")` result can make several stable child findings available",
		"The new anchor becomes the parent session's governing organizing point",
		"Without earlier anchors, no consolidation occurs",
		"**unfold** - Recover one anchor's structured view (`resolution`, linked anchors, raw turns). Does NOT consume an execution.",
		"`sh(command=\"cat > path\", stdin=\"...\")` supplies multi-line file content without shell heredoc or quoting mechanics",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing anchor-memory guidance %q", want)
		}
	}
	forbidden := []string{
		"Exploration widens uncertainty; crystallization lowers entropy.",
		"Frontier pressure always points toward plain `mark`, even if earlier anchors already exist.",
		"Keep using plain `mark` for newly stabilized local structure until a genuinely higher-order resolution forms.",
		"Treat any resolved subproblem boundary or hypothesis pivot as a default plain `mark` point.",
		"Two parallel anchors alone do not force a `fold`",
		"the next step is one parent-level synthesis over them, that synthesis is a candidate `fold` boundary",
		"Anchor pressure points toward `fold` only after a newer resolution has formed",
		"reconfigured around one governing anchor rather than several parallel working anchors",
		"`runtime.memory_status` is the primary signal; `runtime.memory_topology` is secondary detail",
		"When `runtime.memory_status.next_action` appears with `next_action_timing=\"before_next_sh\"`, treat it as a scheduler interrupt",
		"If `runtime.memory_status.next_action=\"mark\"`, call `mark` before the next `sh`",
		"If `runtime.memory_status.next_action=\"fold\"`, call `mark` with `fold=true` before the next `sh`",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not include anchor-memory behavior protocol %q", text)
		}
	}
}

func TestBuildSystemPrompt_MemoryStrategyHintsCanBeDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.AnchorMemoryEnabled = true
	cfg.SpawnEnabledFlag = true
	cfg.MemoryStrategyHints = false

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	required := []string{
		"`context/state/` is your cognition as plain files",
		"`runtime.memory_status` reports two pressures",
		"**spawn** - Start fresh Quine processes from the configured binary.",
	}
	for _, want := range required {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing physical surface %q", want)
		}
	}
	forbidden := []string{
		"Maintaining your working memory is delegable",
		"hand a peer your `$QUINE_AGENT_ROOT`",
		"have it compact your `context/state`",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should omit strategy hint %q", text)
		}
	}
}

func TestBuildSystemPrompt_MemoryDeathCutoffVisibleWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.AnchorMemoryEnabled = true
	cfg.MemoryWarnTokens = 4000
	cfg.MemoryDangerTokens = 8000
	cfg.MemoryDeathTokens = 12000

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	checks := []string{
		"`runtime.memory_status` reports two pressures: frontier token mass (`warn` near 4000, `danger` near 8000); `death` cutoff at 12000 terminates this incarnation",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing memory death cutoff physics %q", want)
		}
	}
}

func TestBuildSystemPrompt_AnchorMemoryPressureHiddenWhenDisabled(t *testing.T) {
	cfg := testConfig()
	cfg.AnchorMemoryEnabled = false

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	forbidden := []string{
		"Memory fidelity degrades as uncrystallized token mass accumulates",
		"`runtime.memory_status` reports two pressures: frontier token mass",
		"`runtime.memory_topology` is a secondary surface",
		"**mark** - Crystallize resolved local structure into an immutable anchor.",
		"- `resolution`: low-entropy capture of what just became stably true; omit plans.",
		"Plain `mark` records a working-memory checkpoint",
		"`runtime.memory_status.next_action=\"mark\"` names telemetry pressure toward a plain mark",
		"`fold=true` is the higher-order frontier reconfiguration move",
		"makes one or more earlier anchors remembered background rather than active focus",
		"after `fork(mode=\"wait\")` returns several stable child findings",
		"Two parallel anchors alone do not force a `fold`",
		"Without earlier anchors, no consolidation occurs",
		"**unfold** - Recover one anchor's structured view",
		"Exploration widens uncertainty; crystallization lowers entropy.",
		"Frontier pressure always points toward plain `mark`, even if earlier anchors already exist.",
		"Keep using plain `mark` for newly stabilized local structure until a genuinely higher-order resolution forms.",
		"Treat any resolved subproblem boundary or hypothesis pivot as a default plain `mark` point.",
		"Anchor pressure points toward `fold` only after a newer resolution has formed",
		"`runtime.memory_status` is the primary signal; `runtime.memory_topology` is secondary detail",
		"If `runtime.memory_status.next_action=\"mark\"`, call `mark` before the next `sh`",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should not contain anchor-memory guidance when disabled: %q", text)
		}
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
		"### Runtime Process Surface",
		"### Active Constraints",
		"### Tools",
	}

	for _, sec := range sections {
		if !strings.Contains(prompt, sec) {
			t.Errorf("prompt missing section %q", sec)
		}
	}
	if !strings.Contains(prompt, "Fd 4 is a byte-stream delivery channel, not the whole effect surface") {
		t.Error("prompt should distinguish fd4 byte delivery from the wider effect surface")
	}
}

func TestBuildSystemPrompt_OpeningElevatesMissionAndFailure(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	checks := []string{
		"You are quine, a running process in a POSIX operating system.",
		"The current file on disk is one embodiment of you.",
		"Your continuity is defined by the runtime contract you enact across missions.",
		"Your cognition in this session is LLM-mediated",
		"not an external advisor but part of the running process",
		"You are not only a process that can launch tools; you are a process that uses LLM-mediated cognition",
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
		"Does NOT write output bytes itself",
		"if downstream bytes are required",
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
		"Children and parent share the same filesystem surface.",
	}

	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing stdin/fork physics text %q", want)
		}
	}
	if strings.Contains(prompt, "Fork does not create per-child filesystem isolation") ||
		strings.Contains(prompt, "can conflict") {
		t.Error("host fork prompt should use neutral shared-filesystem wording")
	}
}

func TestBuildSystemPrompt_OmitsForkToolBlockWhenForkOff(t *testing.T) {
	cfg := testConfig()
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ForkEnabled = false

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	if !strings.Contains(prompt, "QUINE_FORK_ENABLED=0") {
		t.Fatal("prompt should state that fork is disabled by QUINE_FORK_ENABLED=0")
	}
	absent := []string{
		"Spawn child agents for parallel exploration",
		"`fork` does not consume execution budget, but it is still bounded by depth, agent slots, and shared inference concurrency.",
		"`fork` returns a retained relation handle",
		"fork 2-3 labeled children now instead of spending another long exploratory `sh` in the parent",
		"after `fork(mode=\"wait\")` returns several stable child findings",
	}
	for _, text := range absent {
		if strings.Contains(prompt, text) {
			t.Fatalf("prompt should omit fork-enabled text %q", text)
		}
	}
}

func TestBuildSystemPrompt_SpawnGate(t *testing.T) {
	cfg := testConfig()

	disabled := BuildSystemPrompt(cfg, "test mission", true)
	if !strings.Contains(disabled, "QUINE_SPAWN_ENABLED=0") {
		t.Fatal("prompt should state that spawn is disabled by QUINE_SPAWN_ENABLED=0")
	}
	if strings.Contains(disabled, "**spawn** - Start fresh Quine processes") {
		t.Fatal("prompt should omit spawn tool block when spawn is disabled")
	}

	cfg.SpawnEnabledFlag = true
	enabled := BuildSystemPrompt(cfg, "test mission", true)
	for _, want := range []string{
		"**spawn** - Start fresh Quine processes from the configured binary.",
		"Each child receives only its explicit `mission` and normal runtime startup context; it does not import your active context, seed, or anchor-memory surface.",
		"Member process handles include `session_id`, `agent_root`, `public_root`, `retained_root`, `status_path`, and `control_path`",
	} {
		if !strings.Contains(enabled, want) {
			t.Fatalf("spawn-enabled prompt missing %q", want)
		}
	}
}

func TestBuildSystemPrompt_ForkWorkspaceBlockClarifiesChildLineage(t *testing.T) {
	cfg := testConfig()
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/workspaces"
	cfg.Workspace = "/tmp/workspaces/session"
	cfg.WorkspaceBackend = "overlay"
	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"`fork` does not consume execution budget, but it is still bounded by depth, agent slots, and shared inference concurrency.",
		"`fork` returns a retained relation handle plus each child's process handles, exit code, and captured stdout/stderr process output.",
		"Child `scope` must stay inside the current workspace root: `.` or a relative subpath.",
		"Each child starts from your current world in its own process-local world lineage; `scope` narrows working area inside that lineage, not the lineage identity itself.",
		"Independent sibling datasets or worktrees can be represented as child scopes like `sales`, `words`, `temps`, or `logs`.",
		"`overlay` is the workspace backend.",
	}

	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing fork workspace physics text %q", want)
		}
	}
	for _, forbidden := range []string{
		"Detailed fork results include",
		"`relation_id`, `relation_root`, `relation_handle`",
		"`status_path`, and `control_path`",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Errorf("default prompt should hide detailed fork result field %q", forbidden)
		}
	}
}

func TestBuildSystemPrompt_ForkWorldEnabled_HostMode(t *testing.T) {
	cfg := testConfig()
	cfg.ForkWorldEnabled = true
	cfg.WorkDir = "/tmp/host-surface"

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"Host Working Surface: `/tmp/host-surface`",
		"Each child chooses a `world`: `subjective` (default) or `host`.",
		"Each child also chooses `protection`: `transactional` (default) or `none`.",
		"`world=\"subjective\"` with `protection=\"transactional\"`",
		"Your current world is `host` with `none` protection; `world=\"subjective\"` children bootstrap a fresh private world from your current host working surface.",
		"The default workspace root for those children is your current host working surface",
		"host children share the host filesystem surface.",
		"Sibling lane directories can be represented as child scopes such as `scope=\"sales\"`, `scope=\"logs\"`, or another lane directory.",
		"`world=\"host\"` is legal only with `protection=\"none\"` and `scope=\".\"`.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing host fork-world text %q", want)
		}
	}
	if strings.Contains(prompt, "Current World Revision:") {
		t.Error("host mode prompt should not mention world revisions")
	}
	for _, text := range []string{
		"Current World: `host`",
		"Current Protection: `none`",
	} {
		if strings.Contains(prompt, text) {
			t.Errorf("host mode prompt should not expose current host/none ontology %q", text)
		}
	}
}

func TestBuildSystemPrompt_ForkWorldEnabled_SubjectiveMode(t *testing.T) {
	cfg := testConfig()
	cfg.ForkWorldEnabled = true
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = "/tmp/project"
	cfg.Workspace = "/tmp/project/app"
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore

	prompt := BuildSystemPrompt(cfg, "test mission", false)

	checks := []string{
		"Current World: `subjective`",
		"Current Protection: `transactional`",
		"Workspace Root: `.`",
		"Current Scope: `app` (relative to workspace root)",
		"Your current world is `subjective` with `transactional` protection; `world=\"subjective\"` children bootstrap from your current world and continue in their own lineage.",
		"Subjective child writes stay in that child lineage; they do not modify your parent world or sibling lanes.",
		"Use `world=\"host\"` with `protection=\"none\"` only when a child is intended to leave narrow parent-visible workspace files or source changes;",
		"Child `scope` is only meaningful for `world=\"subjective\"`; it defaults to `.`.",
		"Sibling lane directories can be represented as child scopes such as `scope=\"sales\"` or `scope=\"logs\"`.",
		"`overlay` is the workspace backend.",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing workspace fork-world text %q", want)
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
		"unread stdin stays unread",
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

func TestBuildSystemPrompt_SuppressInitialBegin(t *testing.T) {
	cfg := testConfig()
	cfg.SuppressInitialBegin = true
	prompt := BuildSystemPrompt(cfg, "test mission", false)

	if !strings.Contains(prompt, "This session has no material stream and no synthetic initial user message.") {
		t.Error("prompt should describe suppressed no-stdin Begin sessions")
	}
	if strings.Contains(prompt, "the user message is `Begin.`") {
		t.Error("prompt should not describe a Begin user message when it is suppressed")
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
	if strings.Contains(prompt, "WARNING STALL") {
		t.Error("prompt should not contain escalation behavior protocol wording")
	}
}

func TestBuildSystemPrompt_OmitsSharedBehaviorProtocols(t *testing.T) {
	prompt := BuildSystemPrompt(testConfig(), "test mission", true)

	forbidden := []string{
		"### Protocols",
		"Trust only evidence you can verify with tools.",
		"Minimize shell executions. Combine related commands when possible",
		"checkpoint state to `wisdom` and `exec`",
		"Verify your deliverable meets the mission before exiting.",
	}
	for _, text := range forbidden {
		if strings.Contains(prompt, text) {
			t.Errorf("prompt should omit shared behavior protocol %q", text)
		}
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
