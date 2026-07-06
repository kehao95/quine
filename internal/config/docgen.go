package config

// docgen.go renders the two generated projections of the capability registry:
//
//   (a) Paper/core/registries/env-controls.md — the canonical env matrix
//   (b) .env.example                          — the operator-facing preset
//
// Work order:
//   Paper/_design/migrations/runtime-capability-registry-execution.md (T1.2)
// Design authority:
//   Paper/theory/views/runtime-capability/registry-design-brief.md
//   (§ A "Consumers/projections" item 2 — the hand-maintained doc becomes a
//   projection; the doc<->code drift window closes)
//
// The renderers are pure functions over Registry (registry.go), the retired
// table (registry_retired.go), and the presentation data below. Row CONTENT
// comes from the registry; row ORDER, section membership, the per-row
// validation-owner column, and the surrounding prose are presentation facts
// owned here. docgen_test.go asserts the committed files byte-match the
// rendered output (content-equality freshness); scripts/gen-env-docs is the
// thin filesystem shell that writes/checks them.
//
// The validation-owner column is doc-layer data because the Knob schema has
// no owner field (T1.1); if that column starts drifting, promote it into the
// registry instead of growing the map here.

import (
	"fmt"
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------
// Presentation data: section layout
// ---------------------------------------------------------------------------

type envDocRowKind int

const (
	rowLive envDocRowKind = iota
	rowRetired
	rowExternal
)

type envDocRow struct {
	kind envDocRowKind
	env  string
}

func live(env string) envDocRow     { return envDocRow{kind: rowLive, env: env} }
func retired(env string) envDocRow  { return envDocRow{kind: rowRetired, env: env} }
func external(env string) envDocRow { return envDocRow{kind: rowExternal, env: env} }

type envDocSection struct {
	title string
	intro string // optional prose paragraph between heading and table
	rows  []envDocRow
}

// envControlsLayout mirrors the section structure and row order of
// env-controls.md. Every live registry knob must appear exactly once, every
// RetiredRegistry / ExternalLabels entry exactly once (renderer-enforced).
var envControlsLayout = []envDocSection{
	{
		title: "Runtime Identity And Protocol",
		rows: []envDocRow{
			live(EnvModelID),
			live(EnvAPIType),
			live(EnvAPIBase),
			live(EnvAPIKey),
			external("QUINE_PROVIDER"),
			live(EnvConfigDir),
			live(EnvUserAgent),
			live(EnvThinkingBudget),
			live(EnvModelServiceTier),
			retired("QUINE_SMART_MODEL_ID"),
			retired("QUINE_SMART_API_TYPE"),
			retired("QUINE_SMART_API_BASE"),
			retired("QUINE_SMART_API_KEY"),
			live(EnvDebugRequestBodyDir),
		},
	},
	{
		title: "Workspace And Runtime-Surface Backends",
		rows: []envDocRow{
			live(EnvWorkspaceRoot),
			live(EnvWorkspace),
			live(EnvWorkspaceBackend),
			live(EnvWorkspaceOverlayDriver),
			live(EnvWorkspaceRevisionMode),
			live(EnvWorkspaceCurrentRevision),
			live(EnvWorkspaceSession),
			live(EnvWorkspaceOwner),
			live(EnvWorkspaceBootstrap),
			live(EnvWorkspaceCommitOnSignal),
			live(EnvWorkspaceSource),
			retired("QUINE_RUNTIME_SURFACE_BACKEND"),
			live(EnvWorldOnePerShell),
		},
	},
	{
		title: "Capability Exposure",
		rows: []envDocRow{
			live(EnvForkEnabled),
			live(EnvForkWorldEnabled),
			live(EnvSpawnEnabled),
			live(EnvExecEnabled),
			live(EnvExitEnabled),
			live(EnvVisionEnabled),
			live(EnvShInteractiveEnabled),
			live(EnvShStdinEnabled),
			live(EnvShDetachEnabled),
			live(EnvIdleEnabled),
			live(EnvAnchorMemory),
			live(EnvAnchorFoldEnabled),
			live(EnvAnchorMarkEnabled),
			live(EnvAgentsMDEnabled),
			live(EnvAgentsSkillsEnabled),
			live(EnvSelfSourceCodeEnabled),
			live(EnvPeerDiscoveryEnabled),
			live(EnvPeerDiscoveryHeartbeat),
		},
	},
	{
		title: "Prompt Disclosure",
		intro: "These are disclosure gates. They change model-facing explanation without changing runtime acceptance, emitted telemetry, or substrate physics.",
		rows: []envDocRow{
			live(EnvPromptMetaphor),
			live(EnvPromptSelfModel),
			live(EnvPromptRuntimeSurface),
			live(EnvPromptInstructionSurface),
			live(EnvInitialUserMessage),
			live(EnvPromptPersona),
			live(EnvPromptCtl),
			live(EnvPromptImplDetails),
			live(EnvNoMissionAutonomy),
			live(EnvPromptBudgetVisibility),
		},
	},
	{
		title: "Runtime Evidence And Telemetry",
		rows: []envDocRow{
			live(EnvFSMutationTelemetry),
		},
	},
	{
		title: "Response Governance And Execution Envelope",
		rows: []envDocRow{
			live(EnvFailOnImpossible),
			live(EnvMaxTurns),
			retired("QUINE_TURN_EXHAUSTION_POLICY"),
			live(EnvMaxDepth),
			live(EnvDepth),
			live(EnvMaxAgents),
			live(EnvMaxConcurrent),
			live(EnvForkDefaultTimeout),
			live(EnvShDefaultTimeout),
			retired("QUINE_SH_TIMEOUT"),
			live(EnvShTimeoutOverride),
			live(EnvWallClockExitSeconds),
			live(EnvOutputTruncate),
			retired("QUINE_STALL_THRESHOLD"),
			live(EnvEmptyAssistantSuccess),
			live(EnvEphemeralBodyEnabled),
			live(EnvReadyTextAutoIdle),
			live(EnvSuppressInitialBegin),
		},
	},
	{
		title: "Runtime Paths, Context Pressure, And Lineage",
		rows: []envDocRow{
			live(EnvContextWindow),
			live(EnvMemoryWarnTokens),
			live(EnvMemoryDangerTokens),
			live(EnvMemoryDeathTokens),
			live(EnvMemoryStrategyHints),
			retired("QUINE_WISDOM_*"),
			live(EnvDataDir),
			live(EnvRetentionDir),
			live(EnvWorkDir),
			live(EnvShell),
			live(EnvShNetwork),
			live(EnvSelfReentryMode),
			live(EnvSelfReentryTarget),
			live(EnvSessionID),
			live(EnvRunID),
			live(EnvTapeID),
			live(EnvContextTape),
			live(EnvParentSession),
		},
	},
}

// envDocValidationOwners is the "Primary validation owner" column, keyed by
// env name. Doc-layer data (see the file header); retired rows default to
// "—" when absent.
var envDocValidationOwners = map[string]string{
	EnvModelID:          "`internal/config/config_test.go`, `internal/llm/*_test.go`",
	EnvAPIType:          "`internal/config/config_test.go`, `internal/llm/*_test.go`, `internal/runtime/prompt_test.go`",
	EnvAPIBase:          "`internal/config/config_test.go`, `internal/llm/*_test.go`, `internal/runtime/prompt_test.go`",
	EnvAPIKey:           "`internal/config/config_test.go`, `internal/llm/transport/*_test.go`, `internal/runtime/prompt_test.go`",
	"QUINE_PROVIDER":    "`.env.example`, `tests/runtime/run.sh`, `tests/model/run.sh`",
	EnvConfigDir:        "`internal/config/config_test.go`, OAuth transport tests",
	EnvUserAgent:        "`internal/config/config_test.go`, `internal/llm/transport/*_test.go`",
	EnvThinkingBudget:   "`internal/config/config_test.go`, `internal/llm/*_test.go`, `EVALUATION.md`",
	EnvModelServiceTier: "`internal/config/config_test.go`, `internal/llm/protocol/openai_responses_test.go`",

	EnvDebugRequestBodyDir: "`internal/llm/provider.go`",

	EnvWorkspaceRoot:                "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `tests/runtime/COVERAGE_MAP.md`",
	EnvWorkspace:                    "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, runtime workspace tests",
	EnvWorkspaceBackend:             "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `internal/tools/fork_test.go`, runtime workspace tests",
	EnvWorkspaceOverlayDriver:       "`internal/config/config_test.go`, `internal/tools/workspace_fuse_overlayfs_linux_test.go`, runtime workspace tests",
	EnvWorkspaceRevisionMode:        "`internal/config/config_test.go`, `internal/tools/tools_test.go`, runtime workspace tests",
	EnvWorkspaceCurrentRevision:     "`internal/config/config_test.go`, `internal/tools/fork_test.go`, runtime workspace tests",
	EnvWorkspaceSession:             "`internal/config/config_test.go`, runtime workspace tests",
	EnvWorkspaceOwner:               "`internal/config/config_test.go`, workspace tests",
	EnvWorkspaceBootstrap:           "`internal/config/config_test.go`, workspace tests",
	EnvWorkspaceCommitOnSignal:      "`internal/config/config_test.go`, runtime finalization tests",
	EnvWorkspaceSource:              "`internal/config/config_test.go`",
	"QUINE_RUNTIME_SURFACE_BACKEND": "`internal/config/config_test.go`",
	EnvWorldOnePerShell:             "`internal/config/config_test.go`, `internal/world/state_test.go`, `habitat/world/app/main.go`",

	EnvForkEnabled:            "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvForkWorldEnabled:       "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/tools/fork_test.go`, `internal/runtime/prompt_test.go`",
	EnvSpawnEnabled:           "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/tools/spawn_test.go`, `internal/runtime/prompt_test.go`, `tests/runtime/COVERAGE_MAP.md`",
	EnvExecEnabled:            "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvExitEnabled:            "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvVisionEnabled:          "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`",
	EnvShInteractiveEnabled:   "`internal/config/config_test.go`, `internal/tools/interactive_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/prompt_test.go`, `internal/runtime/shell_test.go`",
	EnvShStdinEnabled:         "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/prompt_test.go`, `internal/runtime/shell_test.go`",
	EnvShDetachEnabled:        "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/prompt_test.go`, `internal/runtime/shell_test.go`",
	EnvIdleEnabled:            "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvAnchorMemory:           "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvAnchorFoldEnabled:      "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvAnchorMarkEnabled:      "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvAgentsMDEnabled:        "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`, `tests/runtime/COVERAGE_MAP.md`, `tests/runtime/run.sh`, `tests/model/run.sh`",
	EnvAgentsSkillsEnabled:    "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`, `tests/runtime/COVERAGE_MAP.md`, `tests/runtime/run.sh`, `tests/model/run.sh`",
	EnvSelfSourceCodeEnabled:  "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`",
	EnvPeerDiscoveryEnabled:   "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `tests/runtime/COVERAGE_MAP.md`, `tests/runtime/run.sh`",
	EnvPeerDiscoveryHeartbeat: "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `tests/runtime/COVERAGE_MAP.md`, `tests/runtime/run.sh`",

	EnvPromptMetaphor:           "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvPromptSelfModel:          "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvPromptRuntimeSurface:     "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvPromptInstructionSurface: "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvInitialUserMessage:       "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `cmd/quine/main.go`",
	EnvPromptPersona:            "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvPromptCtl:                "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvPromptImplDetails:        "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvNoMissionAutonomy:        "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `tests/model/run.sh`",
	EnvPromptBudgetVisibility:   "`internal/world/state_test.go`, `internal/config/config_test.go`",

	EnvFSMutationTelemetry: "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `internal/tools/tools_test.go`, `internal/tools/schemas_test.go`, `internal/tools/fork_test.go`",

	EnvFailOnImpossible:      "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/runtime_test.go`",
	EnvMaxTurns:              "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, runtime budget tests",
	EnvMaxDepth:              "`internal/config/config_test.go`, `cmd/quine/integration_test.go`",
	EnvDepth:                 "`internal/config/config_test.go`, `cmd/quine/integration_test.go`, runtime fork/spawn tests",
	EnvMaxAgents:             "`internal/config/config_test.go`, `internal/runtime/semaphore_test.go`, runtime fork tests",
	EnvMaxConcurrent:         "`internal/config/config_test.go`, `internal/runtime/semaphore_test.go`",
	EnvForkDefaultTimeout:    "`internal/config/config_test.go`, `internal/tools/fork_test.go`",
	EnvShDefaultTimeout:      "`internal/config/config_test.go`, `internal/tools/tools_test.go`, `EVALUATION.md`",
	EnvShTimeoutOverride:     "`internal/config/config_test.go`, `internal/tools/schemas_test.go`, `internal/runtime/prompt_test.go`, `internal/runtime/shell_test.go`",
	EnvWallClockExitSeconds:  "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`",
	EnvOutputTruncate:        "`internal/config/config_test.go`, `internal/tools/tools_test.go`, `EVALUATION.md`",
	EnvEmptyAssistantSuccess: "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`",
	EnvEphemeralBodyEnabled:  "`cmd/quine/body_test.go`, `internal/config/config_test.go`, `internal/runtime/prompt_test.go`",
	EnvReadyTextAutoIdle:     "`internal/config/config_test.go`, `internal/runtime/run_loop_test.go`",
	EnvSuppressInitialBegin:  "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `internal/runtime/run_loop_test.go`",

	EnvContextWindow:       "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `EVALUATION.md`",
	EnvMemoryWarnTokens:    "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `EVALUATION.md`",
	EnvMemoryDangerTokens:  "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `EVALUATION.md`",
	EnvMemoryDeathTokens:   "`internal/config/config_test.go`, `internal/runtime/run_loop_test.go`",
	EnvMemoryStrategyHints: "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `internal/tools/memory_test.go`",
	EnvDataDir:             "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, runtime path tests",
	EnvRetentionDir:        "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `tests/runtime/COVERAGE_MAP.md`",
	EnvWorkDir:             "`internal/config/config_test.go`, `internal/tools/tools_test.go`, `EVALUATION.md`",
	EnvShell:               "`internal/config/config_test.go`, `internal/tools/tools_test.go`",
	EnvShNetwork:           "`internal/config/config_test.go`, `internal/tools/job_attrs_linux_test.go`, `internal/runtime/prompt_test.go`",
	EnvSelfReentryMode:     "`internal/config/config_test.go`",
	EnvSelfReentryTarget:   "`internal/config/config_test.go`",
	EnvSessionID:           "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `tests/runtime/COVERAGE_MAP.md`",
	EnvRunID:               "`internal/config/config_test.go`, `internal/runtime/prompt_test.go`, `internal/runtime/runtime_test.go`",
	EnvTapeID:              "`internal/config/config_test.go`, tape tests, `tests/runtime/COVERAGE_MAP.md`",
	EnvContextTape:         "`internal/config/config_test.go`, `internal/tools/fork_test.go`, `internal/runtime/tool_handlers.go`",
	EnvParentSession:       "`internal/config/config_test.go`, `internal/runtime/runtime_test.go`, `tests/runtime/COVERAGE_MAP.md`",
}

// ---------------------------------------------------------------------------
// Markdown cell helpers
// ---------------------------------------------------------------------------

// envTokenRe matches standalone QUINE_* env tokens (optionally with an
// =value or /subpath tail) so registry prose renders them as code spans the
// way the hand-written doc did.
var envTokenRe = regexp.MustCompile(`\bQUINE_[A-Z0-9_]+\*?(?:[=/][A-Za-z0-9_./<>-]+)?`)

// mdCell renders plain registry prose (Scope/Notes strings) as a markdown
// table cell: QUINE_* tokens become code spans, pipes are escaped.
func mdCell(s string) string {
	s = envTokenRe.ReplaceAllString(s, "`$0`")
	return strings.ReplaceAll(s, "|", "\\|")
}

// mdCellPre renders pre-formatted markdown (retired notes, owner cells) as a
// table cell: only pipe escaping, no token rewriting.
func mdCellPre(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// defaultCell renders the Default column for a live knob.
func defaultCell(k Knob) (string, error) {
	switch k.Default.Kind {
	case DefaultRequired:
		return "required", nil
	case DefaultValue:
		if k.Default.Value == "" {
			return "unset", nil
		}
		return "`" + k.Default.Value + "`", nil
	case DefaultDerived:
		from, ok := KnobByName(k.Default.From)
		if !ok {
			return "", fmt.Errorf("knob %s: derived default references unknown knob %q", k.Name, k.Default.From)
		}
		return "derived from `" + from.Env + "`", nil
	case DefaultRuntimeEmitted:
		return "runtime-emitted", nil
	case DefaultLegacy:
		return "removed", nil
	case DefaultExternalLabel:
		return "unset", nil
	}
	return "", fmt.Errorf("knob %s: unknown default kind %q", k.Name, k.Default.Kind)
}

func ownerCell(env, fallback string) string {
	if owner, ok := envDocValidationOwners[env]; ok {
		return owner
	}
	return fallback
}

// ---------------------------------------------------------------------------
// env-controls.md renderer
// ---------------------------------------------------------------------------

const envControlsTableHeader = "| Env | Default | Scope | Prompt | Schema | Runtime | Transport | Primary validation owner | Notes |\n" +
	"|-----|---------|-------|--------|--------|---------|-----------|--------------------------|-------|\n"

// envControlsPreamble is the doc prose before the first knob section. It is
// template text (presentation, not knob content); keep edits here honest with
// the registry semantics they describe.
var envControlsPreamble = strings.Join([]string{
	"# Runtime Env Controls",
	"",
	"<!-- GENERATED FILE — rendered from internal/config/registry.go (+ registry_retired.go)",
	"     by `go run ./scripts/gen-env-docs` (renderer: internal/config/docgen.go).",
	"     Do not edit by hand: edit the registry and regenerate. -->",
	"",
	"> Canonical registry for the active `QUINE_*` control surface that shapes runtime behavior, prompt disclosure, transport wiring, and retained lineage.",
	"",
	"Use this registry when you need one answer to \"what does this env actually change?\" without reconstructing the answer from `config.Load()`, prompt code, tool schemas, and transport helpers.",
	"",
	"## Projection Maintenance",
	"",
	"| Field | Value |",
	"|-------|-------|",
	"| Source owner | `internal/config/registry.go` (the compiled capability registry) and `internal/config/registry_retired.go` own row content; `internal/config/docgen.go` owns layout/prose; code, prompt construction, tool schemas, transport helpers, and runtime tests own actual env behavior |",
	"| Projection role | Generated canonical matrix for active runtime/profile `QUINE_*` controls, selectors, propagated lineage knobs, and removed legacy names — a pure projection of the compiled registry (registry-design-brief § A, consumers item 2) |",
	"| Freshness trigger | Any change to `internal/config/registry.go`, `registry_retired.go`, `docgen.go`, or `envnames.go`; regenerate with `go run ./scripts/gen-env-docs` |",
	"| Drift check | `go run ./scripts/gen-env-docs -check` (content equality; enforced by the `check-env-docs-freshness` hook and `internal/config/docgen_test.go`); stale-phrase greps remain in `./scripts/check-runtime-doc-sync.sh --strict` |",
	"| Absorption owner | Code/tests own behavior; the registry owns the matrix content; this file is a rendered projection — edit the registry, never this file |",
	"",
	"## Terms",
	"",
	"- `profile` means a repo `.env.*` preset such as `.env.kimi`; profiles are convenience bundles, not runtime ontology",
	"- `protocol` means `QUINE_API_TYPE`; it selects the wire format and endpoint path family",
	"- `transport` means authentication/signing behavior such as raw API key headers or OAuth sentinel flows",
	"- `backend` means `QUINE_WORKSPACE_BACKEND`; it selects workspace physics, not API routing",
	"- `provider label` means `QUINE_PROVIDER`; it is a harness/profile label and is not read by the `quine` binary",
	"",
	"## Gate Classes",
	"",
	"Gate class names the semantic role of a control. Impact columns below only say which surfaces change.",
	"",
	"| Class | Meaning | Runtime rejection implied? | Naming rule |",
	"|-------|---------|----------------------------|-------------|",
	"| `disclosure` | Changes what the model-facing prompt explains or frames | no | use `QUINE_PROMPT_*` |",
	"| `model-facing exposure` | Changes prompt/schema exposure so the LLM does not generate a tool or argument | no, unless separately marked | use `*_ENABLED` when it toggles a feature surface |",
	"| `runtime acceptance` | Changes whether the runtime accepts or rejects a call, exit, or operation | yes | use `*_ENABLED` or a behavior-specific name |",
	"| `telemetry feature` | Changes whether runtime evidence appears in prompt/schema and tool results | no call rejection; result shape changes | use a feature name, not `QUINE_PROMPT_*` |",
	"| `physics/backend` | Changes substrate, namespace, workspace, process, or filesystem behavior | yes, by changing reality rather than call admission | use backend/physics vocabulary |",
	"| `transport/profile` | Changes provider transport or harness/profile routing | no tool-call semantics | use protocol/transport/profile vocabulary |",
	"",
	"Rename rule: if a control's name and gate class diverge, prefer a complete cutover over compatibility aliases. The active registry must keep describing current code until implementation, tests, presets, and docs move together. Compatibility aliases are exceptional and require explicit Human approval.",
	"",
	"## Name And Valve Style",
	"",
	"Use these styles for new controls and for future renames:",
	"",
	"- positive boolean gates use `*_ENABLED=1|0`; `1` means the capability,",
	"  exposure surface, or telemetry is available, and `0` means it is removed,",
	"  rejected, or suppressed according to that row's gate class",
	"- prompt disclosure families use `QUINE_PROMPT_<SURFACE>=visible|hidden`;",
	"  prompt styling or multi-state self-model selectors use semantic enum values",
	"  such as `off|thermodynamic` or `basic|advanced`",
	"- physics selectors use noun-style names with semantic valve values:",
	"  `*_BACKEND`, `*_MODE`, or `*_POLICY`, not boolean phrasing",
	"- runtime facts and propagated state use noun-style names such as `*_ID`,",
	"  `*_SESSION`, `*_REVISION`, `*_DIR`, or `*_ROOT`; they are not gates",
	"- negative-polarity booleans should not be introduced for new controls; use a",
	"  positive `*_ENABLED=1|0` gate or a policy enum instead",
	"",
	"Current style exceptions are registry facts, not templates for new envs:",
	"",
	"- `QUINE_ANCHOR_MEMORY` is a boolean feature or exposure gate without `_ENABLED`",
	"- `QUINE_PROMPT_CTL` and `QUINE_PROMPT_IMPL_DETAILS` are prompt disclosure",
	"  valves expressed as `1|0` rather than `visible|hidden`",
	"- `QUINE_NO_MISSION_AUTONOMY` is a missionless-opening posture enum retained",
	"  without the `QUINE_PROMPT_*` prefix because it is already an experiment-facing",
	"  ablation key",
	"- `QUINE_FAIL_ON_IMPOSSIBLE` is a mixed behavior/runtime-acceptance policy",
	"  expressed as a boolean",
	"- (No compatibility aliases for these envs)",
	"",
	"Do not add compatibility aliases for these names casually. If one is renamed,",
	"cut over code, tests, presets, and docs in one coherent state transition.",
	"",
	"## Impact Columns",
	"",
	"| Column | Meaning |",
	"|--------|---------|",
	"| `Prompt` | changes what the system prompt teaches or discloses |",
	"| `Schema` | changes advertised tool inventory or argument schema |",
	"| `Runtime` | changes live acceptance, enforcement, filesystem/process physics, or retained state |",
	"| `Transport` | changes protocol selection, auth/signing, endpoint construction, or request shaping |",
	"",
	"## Interpretation Rules",
	"",
	"- Gate class is the semantic contract; the `Prompt`, `Schema`, `Runtime`, and `Transport` columns only report affected surfaces.",
	"- `*_ENABLED` envs may be model-facing exposure gates, runtime acceptance gates, or telemetry feature gates; the row notes must name the distinction when it is not obvious.",
	"- Unless a row is an enum/string mode, boolean runtime envs are strict `1|0` controls; invalid non-empty values are configuration errors rather than silent falsey fallbacks.",
	"- For LLM-originated tool calls, prompt and schema omission are sufficient enforcement for harmless, internal, or legacy-compatible surfaces; runtime rejection is required only when the hidden capability is unsafe, externally reachable, or registered as a runtime acceptance gate.",
	"- `QUINE_PROMPT_*` envs are disclosure controls. Do not use a prompt-prefixed name for controls that change emitted telemetry, retained state, or substrate physics.",
	"- Disclosure gates that shape the running process self-model should propagate across `fork` / `exec` transitions unless the row explicitly says they are launch-local only.",
	"- `QUINE_FAIL_ON_IMPOSSIBLE` is a behavior/runtime acceptance control: it changes prompt posture, `exit` schema, and runtime validation of failure exits.",
	"- `QUINE_FS_MUTATION_TELEMETRY_ENABLED=1|0` is the filesystem-mutation telemetry feature gate. It changes prompt/schema disclosure and emitted `fs_mutations` result telemetry without changing workspace physics or runtime call acceptance.",
	"- `world=\"subjective\"` only guarantees private lineage under the `overlay` backend. Under `direct`, the child request still narrows scope but writes remain host-visible and non-adoptable.",
	"- Helper or test-only envs that belong to a specific harness stay documented with that harness. This file covers the active runtime/profile env surface that shapes the `quine` binary and its canonical presets.",
	"- Runtime-emitted entrypoints such as `QUINE_AGENT_ROOT` are owned by",
	"  [`runtime-surface.md`](./runtime-surface.md) unless they also become",
	"  user-authored controls or selectors.",
	"",
}, "\n")

var envControlsRelated = strings.Join([]string{
	"## Related",
	"",
	"- Runtime inventory: [`runtime-surface.md`](./runtime-surface.md)",
	"- World semantics: [`../architecture/worlds.md`](../architecture/worlds.md)",
	"- Public fork contract: [`../primitives/fork.md`](../primitives/fork.md)",
	"- Developer entrypoint: [`../../../DEVELOPMENT.md`](../../../DEVELOPMENT.md)",
	"",
}, "\n")

func liveRowCells(k Knob) ([]string, error) {
	def, err := defaultCell(k)
	if err != nil {
		return nil, err
	}
	return []string{
		"`" + k.Env + "`",
		def,
		mdCell(k.Scope),
		yesNo(k.Surfaces.Prompt),
		yesNo(k.Surfaces.Schema),
		yesNo(k.Surfaces.Runtime),
		yesNo(k.Surfaces.Transport),
		mdCellPre(ownerCell(k.Env, "—")),
		mdCell(k.Notes),
	}, nil
}

func retiredRowCells(k RetiredKnob) []string {
	return []string{
		"`" + k.Env + "`",
		"retired",
		"legacy / no-op",
		"no", "no", "no", "no",
		mdCellPre(ownerCell(k.Env, "—")),
		mdCellPre(k.Note),
	}
}

func externalRowCells(k ExternalLabelKnob) []string {
	return []string{
		"`" + k.Env + "`",
		"unset",
		mdCellPre(k.Scope),
		"no", "no", "no", "no",
		mdCellPre(ownerCell(k.Env, "—")),
		mdCellPre(k.Note),
	}
}

// RenderEnvControlsDoc renders the full generated content of
// Paper/core/registries/env-controls.md. It errors when the layout and the
// registry disagree (a knob without a row, a row without a knob) so a
// registry change cannot silently fall out of the doc.
func RenderEnvControlsDoc() (string, error) {
	var b strings.Builder
	b.WriteString(envControlsPreamble)
	b.WriteString("\n")

	placedLive := make(map[string]bool, len(Registry))
	placedRetired := make(map[string]bool, len(RetiredRegistry))
	placedExternal := make(map[string]bool, len(ExternalLabels))

	for _, sec := range envControlsLayout {
		b.WriteString("## " + sec.title + "\n\n")
		if sec.intro != "" {
			b.WriteString(sec.intro + "\n\n")
		}
		b.WriteString(envControlsTableHeader)
		for _, row := range sec.rows {
			var cells []string
			switch row.kind {
			case rowLive:
				k, ok := KnobByEnv(row.env)
				if !ok {
					return "", fmt.Errorf("env-controls layout: live row %q is not in the registry", row.env)
				}
				if placedLive[row.env] {
					return "", fmt.Errorf("env-controls layout: live row %q placed twice", row.env)
				}
				placedLive[row.env] = true
				var err error
				cells, err = liveRowCells(k)
				if err != nil {
					return "", err
				}
			case rowRetired:
				k, ok := RetiredKnobByEnv(row.env)
				if !ok {
					return "", fmt.Errorf("env-controls layout: retired row %q is not in RetiredRegistry", row.env)
				}
				if placedRetired[row.env] {
					return "", fmt.Errorf("env-controls layout: retired row %q placed twice", row.env)
				}
				placedRetired[row.env] = true
				cells = retiredRowCells(k)
			case rowExternal:
				var found *ExternalLabelKnob
				for i := range ExternalLabels {
					if ExternalLabels[i].Env == row.env {
						found = &ExternalLabels[i]
						break
					}
				}
				if found == nil {
					return "", fmt.Errorf("env-controls layout: external row %q is not in ExternalLabels", row.env)
				}
				if placedExternal[row.env] {
					return "", fmt.Errorf("env-controls layout: external row %q placed twice", row.env)
				}
				placedExternal[row.env] = true
				cells = externalRowCells(*found)
			}
			b.WriteString("| " + strings.Join(cells, " | ") + " |\n")
		}
		b.WriteString("\n")
	}

	for _, k := range Registry {
		if !placedLive[k.Env] {
			return "", fmt.Errorf("env-controls layout: registry knob %s (%s) has no doc row", k.Name, k.Env)
		}
	}
	for _, k := range RetiredRegistry {
		if !placedRetired[k.Env] {
			return "", fmt.Errorf("env-controls layout: retired knob %s has no doc row", k.Env)
		}
	}
	for _, k := range ExternalLabels {
		if !placedExternal[k.Env] {
			return "", fmt.Errorf("env-controls layout: external label %s has no doc row", k.Env)
		}
	}

	b.WriteString(envControlsRelated)
	return b.String(), nil
}

// ---------------------------------------------------------------------------
// .env.example renderer
// ---------------------------------------------------------------------------

// envExampleRequired pins the presentation of the four required transport
// knobs (the current operator-facing example values). The renderer errors if
// the registry's required set diverges from this list.
var envExampleRequired = []struct{ env, example string }{
	{EnvModelID, "claude-sonnet-4-6"},
	{EnvAPIType, "anthropic"},
	{EnvAPIBase, "https://api.anthropic.com"},
	{EnvAPIKey, "sk-your-key-here"},
}

var envExampleHeader = strings.Join([]string{
	"# Quine Configuration",
	"#",
	"# GENERATED from internal/config/registry.go by `go run ./scripts/gen-env-docs`",
	"# (renderer: internal/config/docgen.go). Do not edit by hand.",
	"#",
	"# Copy this file to .env and fill in your values:",
	"#   cp .env.example .env",
	"#",
	"# Then source it before running quine:",
	"#   source .env",
	"#   quine \"Hello\"",
	"#",
	"# IMPORTANT: Every line must start with 'export' so that",
	"# 'source .env' exports variables to child processes.",
	"",
}, "\n")

var envExampleFooter = strings.Join([]string{
	"# ── Example: Anthropic via Claude Code subscription OAuth ─",
	"# Reuses Claude Code credentials when available. If no fresh token can be",
	"# refreshed, Quine starts a browser OAuth flow. You can also export a",
	"# `CLAUDE_CODE_OAUTH_TOKEN` generated by `claude setup-token`.",
	"# export QUINE_PROVIDER=anthropic     # Harness label only",
	"# export QUINE_MODEL_ID=claude-sonnet-4-6",
	"# export QUINE_API_TYPE=anthropic",
	"# export QUINE_API_BASE=https://api.anthropic.com",
	"# export QUINE_API_KEY=claude-oauth",
	"",
	"# ── Example: OpenAI Responses API with Codex OAuth ───────",
	"# export QUINE_PROVIDER=openai        # Harness label only",
	"# export QUINE_MODEL_ID=gpt-5.1-codex-mini",
	"# export QUINE_API_TYPE=openai-responses",
	"# export QUINE_API_BASE=https://chatgpt.com/backend-api/codex",
	"# export QUINE_API_KEY=codex-oauth",
	"",
	"# ── Example: Google Gemini 3.5 Flash via OpenAI compatibility ─",
	"# export QUINE_PROVIDER=gemini",
	"# export QUINE_MODEL_ID=gemini-3.5-flash",
	"# export QUINE_API_TYPE=openai",
	"# export QUINE_API_BASE=https://generativelanguage.googleapis.com/v1beta/openai",
	"# export QUINE_API_KEY=\"${GEMINI_API_KEY:-${GOOGLE_GENERATIVE_AI_API_KEY:-}}\"",
	"# export QUINE_THINKING_BUDGET=low",
	"",
	"# ── Example: GitHub Copilot OAuth with Claude Sonnet 4.6 ─",
	"# export COPILOT_OAUTH_CLIENT_ID=your_github_oauth_app_client_id",
	"# export QUINE_PROVIDER=copilot       # Harness label only",
	"# export QUINE_MODEL_ID=claude-sonnet-4.6",
	"# export QUINE_API_TYPE=openai",
	"# export QUINE_API_BASE=https://api.business.githubcopilot.com",
	"# export QUINE_API_KEY=copilot-oauth",
	"# export QUINE_THINKING_BUDGET=off",
	"",
}, "\n")

// envExampleSectionHeading renders a section banner in the current file's
// "# ── Title ─────" style, padded to a fixed width.
func envExampleSectionHeading(title string) string {
	const width = 60
	head := "# ── " + title + " "
	pad := width - len([]rune(head))
	if pad < 2 {
		pad = 2
	}
	return head + strings.Repeat("─", pad)
}

// envExampleOperatorSettable reports whether a live knob belongs in
// .env.example: only knobs an operator legitimately authors at launch.
// Runtime-emitted identity, substrate-pinned lineage counters, and legacy
// load-error tombstones stay out.
func envExampleOperatorSettable(k Knob) bool {
	if k.Default.Kind == DefaultLegacy || k.Default.Kind == DefaultRequired {
		return false // required knobs render in the dedicated required block
	}
	return k.Mutability == MutOperatorOnly || k.Mutability == MutExecBoundary
}

func envExampleKnobLine(k Knob) (string, error) {
	value := ""
	derivedFrom := ""
	switch k.Default.Kind {
	case DefaultValue:
		value = k.Default.Value
	case DefaultDerived:
		from, ok := KnobByName(k.Default.From)
		if !ok {
			return "", fmt.Errorf("knob %s: derived default references unknown knob %q", k.Name, k.Default.From)
		}
		derivedFrom = from.Env
	default:
		return "", fmt.Errorf("knob %s: default kind %q not renderable in .env.example", k.Name, k.Default.Kind)
	}

	comment := k.Scope
	switch k.Type.Kind {
	case TypeBool:
		comment += " (1 | 0)"
	case TypeEnum:
		comment += " (" + strings.Join(k.Type.Enum, " | ") + ")"
	}
	if derivedFrom != "" {
		comment += "; default derived from " + derivedFrom
	}

	assignment := "# export " + k.Env + "=" + value
	const commentCol = 50
	if pad := commentCol - len(assignment); pad > 0 {
		assignment += strings.Repeat(" ", pad)
	} else {
		assignment += " "
	}
	return assignment + " # " + comment, nil
}

// RenderEnvExample renders the full generated content of .env.example.
// Only live, operator-settable registry knobs appear; retired names never do
// (check-authored-env-consistency scans .env* for unknown names).
func RenderEnvExample() (string, error) {
	var b strings.Builder
	b.WriteString(envExampleHeader)
	b.WriteString("\n")

	// Required transport knobs, presented as live example values.
	requiredInRegistry := map[string]bool{}
	for _, k := range Registry {
		if k.Default.Kind == DefaultRequired {
			requiredInRegistry[k.Env] = true
		}
	}
	if len(requiredInRegistry) != len(envExampleRequired) {
		return "", fmt.Errorf(".env.example: registry has %d required knobs, template presents %d — update envExampleRequired", len(requiredInRegistry), len(envExampleRequired))
	}
	b.WriteString("# ── Required (all four must be set) ──────────────────────\n")
	for _, req := range envExampleRequired {
		if !requiredInRegistry[req.env] {
			return "", fmt.Errorf(".env.example: template required knob %s is not required in the registry", req.env)
		}
		b.WriteString("export " + req.env + "=" + req.example + "\n")
	}
	b.WriteString("\n")

	// External harness labels (deliberately not runtime knobs).
	b.WriteString("# ── Optional: Harness Label ──────────────────────────────\n")
	for _, l := range ExternalLabels {
		b.WriteString("# export " + l.Env + "=anthropic" + strings.Repeat(" ", 12) + " # " + l.Scope + "; `quine` ignores this env\n")
	}
	b.WriteString("\n")

	// One block per env-controls.md section, operator-settable knobs only.
	for _, sec := range envControlsLayout {
		var lines []string
		for _, row := range sec.rows {
			if row.kind != rowLive {
				continue
			}
			k, ok := KnobByEnv(row.env)
			if !ok {
				return "", fmt.Errorf(".env.example: layout row %q is not in the registry", row.env)
			}
			if !envExampleOperatorSettable(k) {
				continue
			}
			line, err := envExampleKnobLine(k)
			if err != nil {
				return "", err
			}
			lines = append(lines, line)
		}
		if len(lines) == 0 {
			continue
		}
		b.WriteString(envExampleSectionHeading("Optional: "+sec.title) + "\n")
		for _, line := range lines {
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(envExampleFooter)
	return b.String(), nil
}
