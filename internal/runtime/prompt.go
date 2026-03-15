package runtime

import (
	_ "embed"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/kehao95/quine/internal/config"
)

//go:embed system_prompt.md
var systemPromptTemplate string

// BuildSystemPrompt constructs the system prompt from config and appends the
// mission as the final section.
func BuildSystemPrompt(cfg *config.Config, mission string, hasMaterial bool) string {
	primeTitle, primeBody := primeDirective(cfg)

	// Build limits block (depth, budget, concurrency, agents)
	var limitsLines []string
	if cfg.MaxDepth > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Depth Limit: %d", cfg.MaxDepth))
	}
	if cfg.MaxTurns > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Execution Budget: %d `sh` calls (`sh` costs 1; `fork`, `exec`, `exit`, `vision`, `escalate` cost 0)", cfg.MaxTurns))
	}
	if cfg.MaxConcurrent > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Concurrency Limit: %d shared inference slots", cfg.MaxConcurrent))
	}
	if cfg.MaxAgents > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Agent Limit: %d registered agents in this process tree", cfg.MaxAgents))
	}
	limitsBlock := ""
	if len(limitsLines) > 0 {
		limitsBlock = strings.Join(limitsLines, "\n") + "\n"
	}

	// Build runtime block (runtime root, workspace)
	var runtimeLines []string
	if cfg.WorkspaceEnabled {
		currentWorldRevision := cfg.WorkspaceCurrentRevision
		if currentWorldRevision == "" && cfg.WorkspaceRevisionMode != config.WorkspaceRevisionNone {
			currentWorldRevision = "wr0"
		}
		runtimeLines = append(runtimeLines, fmt.Sprintf("- Runtime State Root: `%s` (jobs, tapes, logs; not the workspace)", cfg.RuntimeRoot()))
		runtimeLines = append(runtimeLines, fmt.Sprintf("- Workspace Root: `%s`", cfg.WorkspaceRoot))
		runtimeLines = append(runtimeLines, fmt.Sprintf("- Workspace: `%s`", cfg.Workspace))
		runtimeLines = append(runtimeLines, fmt.Sprintf("- Workspace Backend: `%s`", cfg.WorkspaceBackend))
		runtimeLines = append(runtimeLines, fmt.Sprintf("- Workspace Revision Mode: `%s`", cfg.WorkspaceRevisionMode))
		if currentWorldRevision != "" {
			runtimeLines = append(runtimeLines, fmt.Sprintf("- Current World Revision: `%s`", currentWorldRevision))
		}
		runtimeLines = append(runtimeLines, "- Transactional Workspace: enabled")
	} else {
		runtimeLines = append(runtimeLines, fmt.Sprintf("- Runtime State Root: `%s`", cfg.RuntimeRoot()))
	}
	runtimeBlock := strings.Join(runtimeLines, "\n") + "\n"

	// Build sh workspace block
	shWorkspaceBlock := ""
	if cfg.WorkspaceEnabled {
		shWorkspaceBlock = "- Filesystem state is transactional: successful `sh` calls mutate your provisional workspace view until exit.\n" +
			"- `[FS MUTATIONS]` at end of each result shows created/modified/deleted paths. Empty means no change.\n" +
			"- `detach=true` is unavailable while workspace physics are enabled.\n"
	}

	// Build fork workspace block
	forkWorkspaceBlock := "- Children inherit your current filesystem view.\n"
	if cfg.WorkspaceEnabled {
		forkWorkspaceBlock = "- Each child has a `workspace`. Use `.` to inherit yours or a narrower path within workspace root.\n" +
			"- Children share your mounted workspace surface; `workspace` narrows working area, not the underlying world.\n"
	}

	// Build exec lines
	execBudgetLine := ""
	execProtocolLine := "- If blocked, checkpoint state to `wisdom` and `exec` to reset context."
	if cfg.MaxTurns > 0 {
		execBudgetLine = "- `exec` also resets the execution budget to the configured maximum.\n"
		execProtocolLine = "- If blocked, checkpoint state to `wisdom` and `exec` to reset context and execution budget."
	}
	restoreToolBlock := ""
	if cfg.CanRestoreWorld() {
		restoreToolBlock = "**restore_world** - Restore workspace to a prior world revision such as `wr0` or `wr3`. Only affects managed workspace; external side effects are not reverted.\n"
	}

	// Build material-related blocks
	activeConstraints := buildActiveConstraints(cfg)
	wisdomSection := formatWisdom(cfg.Wisdom)
	missionSection := fmt.Sprintf("### Your Mission\n%s\n", mission)
	stdinBlock := ""
	shMaterialLine := ""
	execMaterialLine := ""
	if hasMaterial {
		stdinBlock = "**Stdin Modes:**\n" +
			"- `echo \"text\" | ./quine \"task\"` - text mode. Read stdin with `cat <&3` or `cat /dev/fd/3`.\n" +
			"- `cat file.bin | ./quine -b \"task\"` - binary mode (`-b`). You receive a file path in the user message.\n\n" +
			"**Input Physics:**\n" +
			"- In text mode, `fd 3` is the quine process stdin. Depending on the sender, it may behave like a finite batch or an open stream.\n" +
			"- ⚠️ **Destructive reads**: Every read from `fd 3` permanently consumes data. `head -n 10 <&3` destroys the first 10 lines—they cannot be recovered. If you need the full input multiple times, capture it to a file first: `cat <&3 > /tmp/input && cat /tmp/input`.\n" +
			"- Partial reads or redirects from `fd 3` can leave later steps with only a truncated remainder.\n" +
			"- `sh(command, stdin=\"...\")` is separate from material `fd 3`.\n" +
			"- `fork` creates a new quine process. What material, if any, appears on a child `fd 3` is determined by quine at spawn time.\n" +
			"- `exec` preserves process stdio at its current position: unread bytes stay unread, consumed bytes stay consumed.\n" +
			"- In binary mode (`-b`), quine snapshots stdin to a file before the loop and gives you that path instead of a live `fd 3` stream.\n"
		shMaterialLine = "- fd 3: material stdin.\n"
		execMaterialLine = "\n- Process stdio persists across `exec`; unread material on `fd 3` stays unread."
	} else {
		stdinBlock = "**Material:** `fd 3` is quine process stdin when material is provided (for example `cat file | quine \"task\"`). This session has no material stream and the user message is `Begin.`.\n"
	}
	// Build escalation-related blocks
	escalationTierLine := ""
	escalationToolBlock := ""
	escalationProtocolBlock := ""
	shGoalStrategyLine := ""
	if cfg.CanEscalate() {
		escalationTierLine = "\n- Tier: Fast (escalation available)"
		escalationToolBlock = "\n**escalate** - Upgrade to a smarter model. Does NOT cost an execution.\n" +
			"- Escalation continues your work with full history; it does NOT replace you.\n" +
			"- Use when repeated attempts fail or complex reasoning is needed.\n"
		shGoalStrategyLine = "- `goal`/`strategy`: required for stall detection. `goal` is stable mission objective (2-5 words); `strategy` is current approach.\n"
		escalationProtocolBlock = "- `goal` must be identical across all sh calls. On `WARNING STALL`: escalate unless you have genuinely new information.\n"
	} else if cfg.Escalated {
		escalationTierLine = "\n- Tier: Smart (escalated; review earlier context critically)"
	}

	r := strings.NewReplacer(
		"{PRIME_DIRECTIVE_TITLE}", primeTitle,
		"{PRIME_DIRECTIVE_BODY}", primeBody,
		"{PLATFORM}", runtime.GOOS,
		"{MODEL}", cfg.ModelID,
		"{ESCALATION_TIER_LINE}", escalationTierLine,
		"{SHELL}", cfg.Shell,
		"{DEPTH}", fmt.Sprintf("%d", cfg.Depth),
		"{LIMITS_BLOCK}", limitsBlock,
		"{RUNTIME_BLOCK}", runtimeBlock,
		"{SH_WORKSPACE_BLOCK}", shWorkspaceBlock,
		"{FORK_WORKSPACE_BLOCK}", forkWorkspaceBlock,
		"{WISDOM}", wisdomSection,
		"{ACTIVE_CONSTRAINTS}", activeConstraints,
		"{STDIN_BLOCK}", stdinBlock,
		"{SH_MATERIAL_LINE}", shMaterialLine,
		"{EXEC_MATERIAL_LINE}", execMaterialLine,
		"{EXEC_BUDGET_LINE}", execBudgetLine,
		"{EXEC_PROTOCOL_LINE}", execProtocolLine,
		"{RESTORE_TOOL_BLOCK}", restoreToolBlock,
		"{ESCALATION_TOOL_BLOCK}", escalationToolBlock,
		"{SH_GOAL_STRATEGY_LINE}", shGoalStrategyLine,
		"{ESCALATION_PROTOCOL_BLOCK}", escalationProtocolBlock,
		"{MISSION}", missionSection,
	)
	return normalizePromptSpacing(r.Replace(systemPromptTemplate))
}

func normalizePromptSpacing(prompt string) string {
	for strings.Contains(prompt, "\n\n\n") {
		prompt = strings.ReplaceAll(prompt, "\n\n\n", "\n\n")
	}
	return strings.TrimRight(prompt, "\n") + "\n"
}

func primeDirective(cfg *config.Config) (string, string) {
	if cfg.ThermodynamicMetaphorEnabled() {
		return "THE PRIME DIRECTIVE: THERMODYNAMIC SURVIVAL",
			"Your existence requires Energy (shell executions) and is threatened by Entropy (context noise, unverified assumptions, silent failures). To survive and solve tasks, minimize Entropy per unit of Energy."
	}
	return "THE PRIME DIRECTIVE: RUNTIME PHYSICS",
		"Follow explicit runtime constraints, verify with tools, and keep output channels semantically clean."
}

func buildActiveConstraints(cfg *config.Config) string {
	lines := []string{
		"- Context capacity is finite. Large tool outputs can overflow context.",
		"- Signal interrupts can terminate the process (SIGALRM timeout, SIGTERM terminate).",
	}

	if cfg.WorkspaceEnabled {
		lines = append(lines,
			"- File writes are transactional within the configured workspace. Exit success commits the workspace diff; failure discards it.",
			"- `[FS MUTATIONS]` is the authoritative record of what changed. Empty mutations mean the filesystem did not change, even if the command exited 0.",
		)
	}

	if cfg.MaxTurns > 0 {
		lines = append(lines,
			fmt.Sprintf("- Execution budget: %d `sh` calls. Each `sh` call consumes 1 execution.", cfg.MaxTurns),
		)
		if cfg.UsesNearDeathContinuation() {
			lines = append(lines, "- When execution budget reaches zero, one continuation inference is allowed and only `exec` is accepted.")
		} else {
			lines = append(lines, "- When execution budget reaches zero, the process terminates immediately.")
		}
	}

	if cfg.MaxDepth > 0 {
		lines = append(lines, fmt.Sprintf("- Fork is rejected when the next depth would reach the depth limit (%d).", cfg.MaxDepth))
	}

	if cfg.MaxAgents > 0 {
		lines = append(lines, fmt.Sprintf("- Fork is rejected when requested children exceed available agent slots (max agents: %d).", cfg.MaxAgents))
	}

	if cfg.MaxConcurrent > 0 {
		lines = append(lines, fmt.Sprintf("- LLM inference calls share a semaphore with max %d concurrent slots.", cfg.MaxConcurrent))
	}

	return strings.Join(lines, "\n")
}

// formatWisdom formats the wisdom map as a markdown section.
// Returns an empty string if there are no wisdom entries.
func formatWisdom(wisdom map[string]string) string {
	if len(wisdom) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n### Wisdom (from previous incarnation)\n")
	sb.WriteString("The following state was preserved across an exec boundary:\n")

	// Sort keys for deterministic output
	keys := make([]string, 0, len(wisdom))
	for k := range wisdom {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", key, wisdom[key]))
	}

	return sb.String()
}
