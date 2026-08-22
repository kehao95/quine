package config

// registry.go is the declarative capability registry: the canonical model of
// every QUINE_* knob the runtime understands. One entry per surviving env
// constant in envnames.go, including runtime-emitted process-identity envs.
//
// Design authority:
//   Paper/theory/views/runtime-capability/registry-design-brief.md (§ A, D1)
// Work order:
//   Paper/_design/migrations/runtime-capability-registry-execution.md (T1.1)
//
// The registry DESCRIBES Load(); it does not drive it. Load() keeps its
// hand-written parsing; internal/config/registry_test.go is the agreement
// test that keeps the two in sync. Projections (generated env-controls.md,
// .env.example, config/registry.json) are later phases.
//
// JSON encoding: camelCase field names via struct tags, chosen for the
// Phase-2 config/registry.json materialization. RegistryJSON() emits the
// whole table.
//
// Locator convention for ImplSites: repo-relative path, optionally suffixed
// with ":<function>" for function-scoped behavior — the same convention as
// Paper/core/mechanisms/implementation-mechanism-map.md, preferring function
// names over line numbers so locators survive drift.

import "encoding/json"

// TypeKind enumerates the value shapes a knob's env encoding can take.
type TypeKind string

const (
	TypeBool   TypeKind = "bool"   // strict "0"|"1" unless Notes say otherwise
	TypeInt    TypeKind = "int"    // strconv.Atoi
	TypeString TypeKind = "string" // free-form (validation, if any, in Notes)
	TypeEnum   TypeKind = "enum"   // closed value set in KnobType.Enum
)

// KnobType is a knob's value type; Enum carries the legal values for
// TypeEnum knobs (canonical spellings; alias normalization lives in Notes).
type KnobType struct {
	Kind TypeKind `json:"kind"`
	Enum []string `json:"enum,omitempty"`
}

// DefaultKind enumerates how a knob's resolved value defaults when unset.
type DefaultKind string

const (
	// DefaultRequired: Load() fails when the env is unset/empty.
	DefaultRequired DefaultKind = "required"
	// DefaultValue: unset resolves to DefaultSpec.Value (canonical string form).
	DefaultValue DefaultKind = "value"
	// DefaultDerived: unset resolves from another knob (DefaultSpec.From);
	// the derivation formula is described in Notes.
	DefaultDerived DefaultKind = "derived"
	// DefaultRuntimeEmitted: the runtime generates/propagates the value;
	// user-authored values are adopted, ignored, or rewritten per Notes.
	DefaultRuntimeEmitted DefaultKind = "runtime-emitted"
	// DefaultLegacy: a retired name kept only as a tombstone;
	// DefaultSpec.Value carries the legacy sub-kind ("removed-error").
	DefaultLegacy DefaultKind = "legacy"
	// DefaultExternalLabel: declared for external harness/profile labeling,
	// never read by the binary. (Anticipated by the brief; currently unused —
	// every envnames.go constant is consumed somewhere in the binary.)
	DefaultExternalLabel DefaultKind = "external-label"
)

// DefaultSpec describes a knob's default resolution.
type DefaultSpec struct {
	Kind  DefaultKind `json:"kind"`
	Value string      `json:"value,omitempty"` // canonical default for "value"; sub-kind for "legacy"
	From  string      `json:"from,omitempty"`  // source knob Name for "derived"
}

// GateClass names the semantic role of a knob (orthogonal to Surfaces).
type GateClass string

const (
	ClassCapability GateClass = "capability" // gates what the runtime/tools can do (incl. physics selectors)
	ClassDisclosure GateClass = "disclosure" // shapes model-facing narration only
	ClassBudget     GateClass = "budget"     // numeric envelopes, limits, thresholds, cadences
	ClassIdentity   GateClass = "identity"   // process/lineage identity and propagated state
	ClassTransport  GateClass = "transport"  // provider wiring, auth, request shaping
	ClassPath       GateClass = "path"       // filesystem locations
)

// Surfaces is the impact matrix: which runtime surfaces a knob changes.
// Mirrors the Prompt/Schema/Runtime/Transport columns of
// Paper/core/registries/env-controls.md.
type Surfaces struct {
	Prompt    bool `json:"prompt"`
	Schema    bool `json:"schema"`
	Runtime   bool `json:"runtime"`
	Transport bool `json:"transport"`
}

// Mutability is the authority boundary for changing a knob's value.
type Mutability string

const (
	// MutSubstratePinned: never agent-stageable; the body enforces or
	// rewrites it (lineage counters, removed-name tombstones).
	MutSubstratePinned Mutability = "substrate-pinned"
	// MutExecBoundary: agent-changeable, but only across an exec boundary
	// (the staged-override tier of the brief's two-tier model).
	MutExecBoundary Mutability = "exec-boundary"
	// MutOperatorOnly: operator infrastructure (API auth, transport wiring,
	// durable state roots); not an agent capability position.
	MutOperatorOnly Mutability = "operator-only"
	// MutRuntimeEmitted: the runtime writes it; agents never author it
	// (the ProcessIdentityEnvNames family).
	MutRuntimeEmitted Mutability = "runtime-emitted"
)

// CouplingKind types the edge between two knobs.
type CouplingKind string

const (
	// CoupleRequires: violating the constraint is a config-load error.
	CoupleRequires CouplingKind = "requires"
	// CoupleSilentWithout: this knob silently no-ops without the peer.
	CoupleSilentWithout CouplingKind = "silent-without"
	// CoupleOverriddenBy: the peer's value can negate/dominate this knob.
	CoupleOverriddenBy CouplingKind = "overridden-by"
	// CoupleOverrides: this knob's value can negate/dominate the peer.
	CoupleOverrides CouplingKind = "overrides"
	// CoupleHazard: the combination loads fine but fails later at enactment.
	CoupleHazard CouplingKind = "hazard"
	// CoupleShares: the knobs share a consumer/envelope; changing one
	// affects the mechanism the other names.
	CoupleShares CouplingKind = "shares"
)

// Coupling is a typed edge to a peer knob, visible at the edit site.
type Coupling struct {
	Peer string       `json:"peer"` // registry Name of the peer knob
	Kind CouplingKind `json:"kind"`
	Note string       `json:"note"`
}

// Knob is one entry in the capability registry (brief § A, post-F-4.7
// amendments: DefaultSpec, Scope, Surfaces, Notes, runtime-emitted).
type Knob struct {
	Name       string      `json:"name"`                 // stable ID; theory docs reference this
	Env        string      `json:"env"`                  // env encoding
	Type       KnobType    `json:"type"`                 // bool | int | string | enum(values)
	Default    DefaultSpec `json:"default"`              // required | value | derived | runtime-emitted | legacy | external-label
	Scope      string      `json:"scope"`                // fine-grained one-line role (env-controls "Scope" column)
	Axes       []string    `json:"axes,omitempty"`       // tags into projection-axes; cheap to re-cut
	Class      GateClass   `json:"class"`                // semantic role
	Surfaces   Surfaces    `json:"surfaces"`             // impact flags, orthogonal to Class
	Mutability Mutability  `json:"mutability"`           // authority boundary
	ParentGate string      `json:"parentGate,omitempty"` // gate trees are first-class
	Couples    []Coupling  `json:"couples,omitempty"`    // typed edges
	Notes      string      `json:"notes,omitempty"`      // knob-intrinsic prose: special values, platform gating
	ImplSites  []string    `json:"implSites,omitempty"`  // locators, pair with implementation-mechanism-map
}

// Axis tags (projection-axes.md vocabulary; the seven discovered axes).
const (
	axMortality    = "mortality"     // Axis 1
	axContinuity   = "continuity"    // Axis 2 (generation continuity)
	axDisclosure   = "disclosure"    // Axis 3 (perceptual disclosure)
	axSelfRelation = "self-relation" // Axis 4
	axDirective    = "directive"     // Axis 5 (ecological directive)
	axSocial       = "social"        // Axis 6 (social density)
	axEnaction     = "enaction"      // Axis 7 (enaction fidelity)
)

// --- small constructors to keep the table readable ---

func boolT() KnobType               { return KnobType{Kind: TypeBool} }
func intT() KnobType                { return KnobType{Kind: TypeInt} }
func strT() KnobType                { return KnobType{Kind: TypeString} }
func enumT(vals ...string) KnobType { return KnobType{Kind: TypeEnum, Enum: vals} }

func defRequired() DefaultSpec      { return DefaultSpec{Kind: DefaultRequired} }
func defValue(v string) DefaultSpec { return DefaultSpec{Kind: DefaultValue, Value: v} }
func defDerived(from string) DefaultSpec {
	return DefaultSpec{Kind: DefaultDerived, From: from}
}
func defRuntime() DefaultSpec { return DefaultSpec{Kind: DefaultRuntimeEmitted} }
func defLegacy(subKind string) DefaultSpec {
	return DefaultSpec{Kind: DefaultLegacy, Value: subKind}
}

func surf(prompt, schema, rt, transport bool) Surfaces {
	return Surfaces{Prompt: prompt, Schema: schema, Runtime: rt, Transport: transport}
}

// Registry is the declarative table: one entry per envnames.go constant.
// Grouping mirrors Paper/core/registries/env-controls.md sections.
var Registry = []Knob{
	// ------------------------------------------------------------------
	// Runtime identity and protocol
	// ------------------------------------------------------------------
	{
		Name: "ModelID", Env: EnvModelID,
		Type: strT(), Default: defRequired(),
		Scope: "model identity", Class: ClassTransport,
		Surfaces: surf(false, false, true, true), Mutability: MutOperatorOnly,
		Notes: "model sent to the selected API",
		ImplSites: []string{
			"internal/config/config.go:loadRequiredIdentityAndTransport",
			"internal/llm/provider.go:NewProvider",
		},
	},
	{
		Name: "APIType", Env: EnvAPIType,
		Type: enumT("openai", "anthropic", "openai-responses"), Default: defRequired(),
		Scope: "protocol selector", Class: ClassTransport,
		Surfaces: surf(true, false, true, true), Mutability: MutOperatorOnly,
		Notes: "canonical runtime selector for the wire format and endpoint path family; prompt may disclose the selected provider transport",
		ImplSites: []string{
			"internal/config/config.go:loadRequiredIdentityAndTransport",
			"internal/llm/provider.go:NewProvider",
		},
	},
	{
		Name: "APIBase", Env: EnvAPIBase,
		Type: strT(), Default: defRequired(),
		Scope: "API base URL", Class: ClassTransport,
		Surfaces: surf(true, false, true, true), Mutability: MutOperatorOnly,
		Notes: "combined with protocol endpoint path; prompt may disclose the configured provider base",
		ImplSites: []string{
			"internal/config/config.go:loadRequiredIdentityAndTransport",
			"internal/llm/provider.go:NewProvider",
		},
	},
	{
		Name: "APIKey", Env: EnvAPIKey,
		Type: strT(), Default: defRequired(),
		Scope: "auth / sentinel", Class: ClassTransport,
		Surfaces: surf(true, false, true, true), Mutability: MutOperatorOnly,
		Notes: "raw key or OAuth sentinel such as codex-oauth, kimi-oauth, copilot-oauth, claude-oauth; prompt discloses only presence/absence, never the value. The runtime never renders a process's own environment back to it (self-readout is /proc/<pid>/environ), so the key value only ever exists in envp, never on a config surface",
		ImplSites: []string{
			"internal/config/config.go:loadRequiredIdentityAndTransport",
			"internal/llm/transport",
		},
	},
	{
		Name: "ConfigDir", Env: EnvConfigDir,
		Type: strT(), Default: defValue(""),
		Scope: "OAuth token/config storage", Class: ClassTransport,
		Surfaces: surf(false, false, true, true), Mutability: MutOperatorOnly,
		Notes: "not loaded into Config: OAuth transports read the process env directly. Free at every boundary, so a child inherits it when the operator set it and does not see it when they did not",
		ImplSites: []string{
			"internal/config/envmodel.go:BuildChildEnv",
			"internal/llm/transport/codexoauth.go",
			"internal/llm/claudeoauth/claudeoauth.go",
		},
	},
	{
		Name: "UserAgent", Env: EnvUserAgent,
		Type: strT(), Default: defValue(""),
		Scope: "HTTP header shaping", Class: ClassTransport,
		Surfaces: surf(false, false, true, true), Mutability: MutOperatorOnly,
		Notes: "custom User-Agent header; propagated to children only when non-empty",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/llm/provider.go:NewProvider",
		},
	},
	{
		Name: "ThinkingBudget", Env: EnvThinkingBudget,
		Type: enumT("off", "low", "medium", "high", "xhigh"), Default: defValue("high"),
		Scope: "request shaping", Class: ClassTransport,
		Surfaces: surf(false, false, true, true), Mutability: MutOperatorOnly,
		Notes: "request option, not a prompt-only knob; some transports fall back to off; value is matched verbatim (no whitespace trimming)",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/llm/provider.go:NewProvider",
		},
	},
	{
		Name: "ModelServiceTier", Env: EnvModelServiceTier,
		Type: enumT("priority", "flex", "fast"), Default: defValue(""),
		Scope: "request shaping", Class: ClassTransport,
		Surfaces: surf(false, false, true, true), Mutability: MutOperatorOnly,
		Notes: "OpenAI Responses service tier; empty means provider default; legacy input fast normalizes to priority at load",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/llm/protocol",
		},
	},
	{
		Name: "DebugRequestBodyDir", Env: EnvDebugRequestBodyDir,
		Type: strT(), Default: defValue(""),
		Scope: "debug dump", Class: ClassPath,
		Surfaces: surf(false, false, true, true), Mutability: MutOperatorOnly,
		Notes: "debug-only request/response dump directory for failed calls; free at every boundary, so children inherit it exactly when the operator set it",
		ImplSites: []string{
			"internal/config/config.go:loadRequiredIdentityAndTransport",
			"internal/llm/provider.go",
		},
	},

	// ------------------------------------------------------------------
	// Workspace and runtime-surface backends
	// ------------------------------------------------------------------
	{
		Name: "WorkspaceRoot", Env: EnvWorkspaceRoot,
		Type: strT(), Default: defValue(""),
		Scope: "managed workspace root", Axes: []string{axEnaction}, Class: ClassPath,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "enables workspace physics when set (also enabled via WORKSPACE or WORKSPACE_SESSION); canonicalized at load; the directory must exist; when only WORKSPACE is set, root adopts it",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/tools/workspace_overlay_linux.go",
		},
	},
	{
		Name: "Workspace", Env: EnvWorkspace,
		Type: strT(), Default: defDerived("WorkspaceRoot"),
		Scope: "managed workspace scope", Axes: []string{axEnaction}, Class: ClassPath,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "WorkspaceRoot", Kind: CoupleRequires, Note: "must canonicalize to a path within the workspace root"},
		},
		Notes: "current writable scope within the workspace root; unset defaults to the startup cwd when the process starts inside the root, else to the root itself",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
		},
	},
	{
		Name: "WorkspaceBackend", Env: EnvWorkspaceBackend,
		Type: enumT("overlay", "direct"), Default: defDerived("WorkspaceRoot"),
		Scope: "workspace backend", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "defaults to overlay when workspace physics are enabled, empty otherwise; overlay is the Linux-only transactional lineage/rollback backend (load error on non-Linux); direct is the shared host-visible backend",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/config/config.go:WorkspaceTransactional",
		},
	},
	{
		Name: "WorkspaceOverlayDriver", Env: EnvWorkspaceOverlayDriver,
		Type: enumT("kernel", "fuse"), Default: defDerived("WorkspaceBackend"),
		Scope: "overlay mount driver", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "WorkspaceBackend", Kind: CoupleRequires, Note: "setting a driver with BACKEND=direct is a config-load error (archaeology coupling #2)"},
		},
		Notes: "defaults to kernel under the overlay backend; kernel uses kernel overlayfs, fuse uses fuse-overlayfs; Linux-only by inheritance from the overlay backend",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/tools/workspace_overlay_linux.go",
		},
	},
	{
		Name: "WorkspaceRevisionMode", Env: EnvWorkspaceRevisionMode,
		Type: enumT("none", "restore"), Default: defDerived("WorkspaceBackend"),
		Scope: "revision surface", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "WorkspaceBackend", Kind: CoupleRequires, Note: "restore requires BACKEND=overlay; direct only supports none; any explicit value requires workspace physics (archaeology coupling #3)"},
		},
		Notes: "defaults to restore under overlay, none under direct (backend-conditional default) and none when workspace physics are disabled; restore exposes switch_world",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/tools/tools.go:SwitchWorld",
		},
	},
	{
		Name: "WorkspaceCurrentRevision", Env: EnvWorkspaceCurrentRevision,
		Type: strT(), Default: defRuntime(),
		Scope: "world lineage state", Axes: []string{axEnaction}, Class: ClassIdentity,
		Surfaces: surf(true, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "runtime-emitted current revision handle, updated after sh/switch_world (the only in-process cfg mutation). Masked at every boundary: it names a revision in THIS process's workspace state, so an sh or fork child inheriting it would hold a handle into a workspace it is not in. Re-stamped only at exec, where the same process continues in the same workspace. Not a user-authored knob",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/config/envstamps.go:SelfReentryStamps",
			"internal/runtime/tool_handlers.go",
		},
	},
	{
		Name: "WorkspaceSession", Env: EnvWorkspaceSession,
		Type: strT(), Default: defDerived("SessionID"),
		Scope: "lineage namespace", Axes: []string{axEnaction}, Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Couples: []Coupling{
			{Peer: "WorkspaceRoot", Kind: CoupleRequires, Note: "setting WORKSPACE_SESSION alone enables workspace physics, which then requires a root"},
		},
		Notes: "stable overlay-state namespace under QUINE_DATA_DIR; defaults to SessionID when workspace physics are enabled. Masked at every boundary and re-stamped only at exec: a fork child gets a fresh session, or the parent's as a WORKSPACE_BOOTSTRAP lineage handle, never as its own",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/config/envstamps.go:SelfReentryStamps",
		},
	},
	{
		Name: "WorkspaceOwner", Env: EnvWorkspaceOwner,
		Type: boolT(), Default: defValue("1"),
		Scope: "ownership bit", Axes: []string{axEnaction}, Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "ownership bit for the process that commits/rolls back workspace state. Masked at every boundary and re-stamped where it means something: fork/spawn children are always stamped 0 (a child borrows a view of its parent's workspace, it does not own it), and exec carries the current value forward",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/config/envstamps.go:ForkChildStamps",
			"internal/config/envstamps.go:SelfReentryStamps",
			"internal/tools/workspace_overlay_linux.go",
		},
	},
	{
		Name: "WorkspaceCommitOnSignal", Env: EnvWorkspaceCommitOnSignal,
		Type: boolT(), Default: defValue("0"),
		Scope: "signal shutdown workspace finalization", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Notes: "when enabled, external signal shutdown commits the latest completed workspace revision instead of rolling it back with unfinished work",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/runtime/runtime.go:setupSignalHandler",
		},
	},
	{
		Name: "WorkspaceBootstrap", Env: EnvWorkspaceBootstrap,
		Type: strT(), Default: defRuntime(),
		Scope: "bootstrap lineage source", Axes: []string{axEnaction}, Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "runtime-emitted lineage source identity: the fork/spawn boundary stamps it with the parent's WorkspaceSession so a child adopts the lineage instead of starting a fresh one. Masked everywhere else — an inherited bootstrap handle would adopt a lineage nobody chose",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/config/envstamps.go:ForkChildStamps",
		},
	},
	{
		Name: "WorkspaceSource", Env: EnvWorkspaceSource,
		Type: strT(), Default: defLegacy("removed-error"),
		Scope: "legacy", Class: ClassPath,
		Surfaces: surf(false, false, true, false), Mutability: MutSubstratePinned,
		Notes: "removed: setting any non-empty value is a config-load error; use WORKSPACE_ROOT and WORKSPACE instead",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
		},
	},
	{
		Name: "WorldOnePerShell", Env: EnvWorldOnePerShell,
		Type: boolT(), Default: defValue("0"),
		Scope: "one-world-per-shell guard", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutOperatorOnly,
		Notes: "Linux-only; gates one-world-per-shell workspace physics via a per-parent-pid stamp file; not loaded into Config — internal/world reads the process env directly, and only the literal value 1 enables it (any other value silently disables; laxer than the strict 0|1 loader bools)",
		ImplSites: []string{
			"internal/world/state.go:EnforceSingleWorldInvocationPerShell",
			"habitat/world/app/main.go",
		},
	},

	// ------------------------------------------------------------------
	// Capability exposure (tool gates)
	// ------------------------------------------------------------------
	{
		Name: "ForkEnabled", Env: EnvForkEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "fork capability", Axes: []string{axContinuity}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "removing fork from prompt/schema does not weaken runtime rejection (dispatch and handler both re-check)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/fork.go",
			"internal/tools/schemas.go",
		},
	},
	{
		Name: "ForkWorldEnabled", Env: EnvForkWorldEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "child world / protection request surface", Axes: []string{axContinuity, axEnaction, axSocial}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "WorkspaceRoot", Kind: CoupleRequires, Note: "=1 without explicit workspace physics is a config-load error (archaeology coupling #1)"},
		},
		Notes: "only exposes the child world/protection request surface; the workspace backend still determines actual lineage guarantees",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/tools/fork.go",
		},
	},
	{
		Name: "SpawnEnabled", Env: EnvSpawnEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "fresh spawn capability", Axes: []string{axContinuity, axSocial}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "exposes fresh Quine process creation without importing parent active context, seed, or anchor-memory surface; children still belong to the same depth/agent-slot governed process tree",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/spawn.go",
		},
	},
	{
		Name: "ExecEnabled", Env: EnvExecEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "exec capability", Axes: []string{axContinuity}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "schema + runtime gate for process-image replacement (self-reentry)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/exec.go",
		},
	},
	{
		Name: "ExitEnabled", Env: EnvExitEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "exit capability", Axes: []string{axMortality}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "disabling exit removes voluntary death: termination is left to execution budget exhaustion or signal (inverse loading on the mortality axis)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/exit.go",
		},
	},
	{
		Name: "VisionEnabled", Env: EnvVisionEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "vision capability", Axes: []string{axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "schema + runtime gate for the image-read tool",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/vision.go:HandleVision",
		},
	},
	{
		Name: "ShInteractiveEnabled", Env: EnvShInteractiveEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "sh(interactive=true) request surface", Axes: []string{axDisclosure}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "removes and rejects PTY interactive mode when disabled; argument-level gate on the sh tool (rejection still costs an execution turn)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/interactive.go",
		},
	},
	{
		Name: "ShStdinEnabled", Env: EnvShStdinEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "sh(stdin=...) request surface", Axes: []string{axDisclosure}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "removes and rejects structured per-call sh stdin when disabled; process material stdin (fd 3) is separate and remains governed by material input posture",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "ShDetachEnabled", Env: EnvShDetachEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "sh(detach=true) request surface", Axes: []string{axDisclosure}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "WorkspaceBackend", Kind: CoupleOverriddenBy, Note: "sh(detach=true) is rejected under the overlay backend even when this gate is on (archaeology coupling #12)"},
		},
		Notes: "removes and rejects detached shell jobs when disabled",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "IdleEnabled", Env: EnvIdleEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "explicit suspension capability", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "exposes idle plus the post/poke/inject/interrupt control slice; the system prompt owns the full peer-control file map, while the schema only names the idle-local resume contract",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/tool_handlers.go:handleIdleEnabled",
			"internal/runtime/control.go:waitForIdleResume",
		},
	},
	{
		Name: "AnchorMemory", Env: EnvAnchorMemory,
		Type: boolT(), Default: defValue("0"),
		Scope: "memory tools", Axes: []string{axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "parent gate of the anchor-memory tree: exposes mark/unfold and constructs the memory executor that all memory-pressure knobs depend on; runtime state surface remains present for internal compatibility; conditionally propagated (children see the sub-gates only when this is 1)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/memory.go",
		},
	},
	{
		Name: "AnchorMarkEnabled", Env: EnvAnchorMarkEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "mark crystallization tool", Axes: []string{axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		ParentGate: "AnchorMemory",
		Notes:      "only meaningful with QUINE_ANCHOR_MEMORY=1; =0 removes the mark tool from schema/prompt and rejects mark calls, leaving unfold + the raw context/state substrate so the only context-load reduction path is the filesystem itself (directly or via a recruited peer)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/memory.go",
		},
	},
	{
		Name: "AnchorFoldEnabled", Env: EnvAnchorFoldEnabled,
		Type: boolT(), Default: defValue("1"),
		Scope: "fold move on mark", Axes: []string{axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		ParentGate: "AnchorMarkEnabled",
		Notes:      "third level of the anchor-memory gate tree (AnchorMemory -> AnchorMark -> AnchorFold): the fold move rides the mark tool, so it needs both ancestors; =0 removes the fold param from the mark schema/prompt and rejects fold calls, leaving plain mark + unfold over the raw context/state file substrate",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/memory.go",
		},
	},
	{
		Name: "AgentsMDEnabled", Env: EnvAgentsMDEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "project AGENTS.md context projection", Axes: []string{axDisclosure, axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Notes: "startup/refresh gate for projecting one discoverable project AGENTS.md into context/prompt/10-agents.md; if more than one is discoverable, startup fails until hierarchical support exists; discovery is bounded by the workspace root when workspace physics are on, else walks up from WORK_DIR",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/config/agents_md.go:DiscoverSingleAgentsMD",
		},
	},
	{
		Name: "AgentsSkillsEnabled", Env: EnvAgentsSkillsEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "project SKILLS context projection", Axes: []string{axDisclosure, axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Notes: "startup/refresh gate for generating context/prompt/20-skills.md from visible .agents/skills/<name>/SKILL.md frontmatter only; hierarchical skill bodies/resources are read later by the agent when relevant",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/config/skills.go:LoadSkills",
		},
	},
	{
		Name: "SelfSourceCodeEnabled", Env: EnvSelfSourceCodeEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "source-code self surface", Axes: []string{axSelfRelation}, Class: ClassCapability,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Notes: "legacy binary gate for the read-only source-code/ self-description surface; when QUINE_SELF_SOURCE_PROJECTION is set, the three-valued projection selector takes precedence",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/self_source.go:syncSelfSourceSurface",
		},
	},
	{
		Name: "SelfSourceProjection", Env: EnvSelfSourceProjection,
		Type: enumT("none", "runtime", "repo"), Default: defValue("none"),
		Scope: "source-code self-surface projection shape", Axes: []string{axSelfRelation, axDisclosure}, Class: ClassCapability,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "SelfSourceCodeEnabled", Kind: CoupleOverrides, Note: "when explicitly set, none disables projection and runtime/repo enable it regardless of the legacy boolean gate"},
		},
		Notes: "none omits source-code; runtime projects only the embedded buildable Quine runtime tree; repo projects the complete embedded repository bundle",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/self_source.go:syncSelfSourceSurface",
		},
	},
	{
		Name: "PeerDiscoveryEnabled", Env: EnvPeerDiscoveryEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "peer topology observation", Axes: []string{axDisclosure, axSocial}, Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Notes: "exposes local runtime.peer_discovery snapshots from existing pid/<pid> discovery at this process's own safe points; also enables stale pid-lock heartbeat pruning; does not broadcast through ctl or create a new public surface",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/peer_discovery.go:startPeerDiscoveryHeartbeat",
		},
	},
	{
		Name: "PeerDiscoveryHeartbeat", Env: EnvPeerDiscoveryHeartbeat,
		Type: intT(), Default: defValue("5000"),
		Scope: "peer stale-lock heartbeat cadence", Axes: []string{axDisclosure, axSocial}, Class: ClassBudget,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		ParentGate: "PeerDiscoveryEnabled",
		Notes:      "scan interval (ms) for heartbeat pruning when QUINE_PEER_DISCOVERY_ENABLED=1; must be positive (load error otherwise); not an independent gate",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/peer_discovery.go:startPeerDiscoveryHeartbeat",
		},
	},

	// ------------------------------------------------------------------
	// Prompt disclosure
	// ------------------------------------------------------------------
	{
		Name: "PromptMetaphor", Env: EnvPromptMetaphor,
		Type: enumT("off", "thermodynamic"), Default: defValue("off"),
		Scope: "framing overlay", Axes: []string{axSelfRelation, axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptInstructionSurface", Kind: CoupleOverriddenBy, Note: "no-op under any minimal* instruction surface (archaeology coupling #9)"},
		},
		Notes: "prompt-only styling",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go",
		},
	},
	{
		Name: "PromptSelfModel", Env: EnvPromptSelfModel,
		Type: enumT("basic", "advanced"), Default: defValue("advanced"),
		Scope: "continuity / cognition disclosure", Axes: []string{axSelfRelation, axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptInstructionSurface", Kind: CoupleOverriddenBy, Note: "no-op under any minimal* instruction surface (archaeology coupling #9)"},
		},
		Notes: "controls richer self-model narration only",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go",
		},
	},
	{
		Name: "PromptInstructionSurface", Env: EnvPromptInstructionSurface,
		Type: enumT("standard", "minimal", "minimal_autonomy", "minimal_existence"), Default: defValue("standard"),
		Scope: "instruction-surface disclosure", Axes: []string{axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptMetaphor", Kind: CoupleOverrides, Note: "minimal* short-circuits the prompt builder, no-opping this peer"},
			{Peer: "PromptSelfModel", Kind: CoupleOverrides, Note: "minimal* short-circuits the prompt builder, no-opping this peer"},
			{Peer: "PromptRuntimeSurface", Kind: CoupleOverrides, Note: "minimal* short-circuits the prompt builder, no-opping this peer"},
			{Peer: "PromptPersona", Kind: CoupleOverrides, Note: "minimal* short-circuits the prompt builder, no-opping this peer"},
			{Peer: "PromptCtl", Kind: CoupleOverrides, Note: "minimal* short-circuits the prompt builder, no-opping this peer"},
			{Peer: "PromptImplDetails", Kind: CoupleOverrides, Note: "minimal* short-circuits the prompt builder, no-opping this peer"},
		},
		Notes: "the disclosure kill switch: minimal emits only the mission/no-mission opening; minimal_autonomy adds only generic missionless self-activation and do-not-wait text; minimal_existence denies any task/objective while declaring a running process with sh and full agency (activation clause carried by the initial user message); all three suppress runtime-teaching instruction blocks for prompt-contamination ablations",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go",
		},
	},
	{
		Name: "PromptRuntimeSurface", Env: EnvPromptRuntimeSurface,
		Type: enumT("visible", "hidden"), Default: defValue("visible"),
		Scope: "runtime-surface disclosure", Axes: []string{axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptInstructionSurface", Kind: CoupleOverriddenBy, Note: "no-op under any minimal* instruction surface (archaeology coupling #9)"},
		},
		Notes: "hiding does not remove the underlying files/env vars",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go:buildRuntimeSurfaceSection",
		},
	},
	{
		Name: "PromptPersona", Env: EnvPromptPersona,
		Type: enumT("coder", "analyst", "engineer", "architect", "explorer", "steward", "cartographer", "gardener", "witness", "catalyst", "skeptic"), Default: defValue(""),
		Scope: "role-stance overlay", Axes: []string{axDisclosure, axDirective}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptInstructionSurface", Kind: CoupleOverriddenBy, Note: "no-op under any minimal* instruction surface (archaeology coupling #9)"},
		},
		Notes: "adds one role-stance sentence; empty means no persona; RETAINED per brief D6-SUSPENDED (dorm-04/05/07/08/09 use it as an experimental variable); value is whitespace-trimmed before enum validation",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go",
		},
	},
	{
		Name: "PromptCtl", Env: EnvPromptCtl,
		Type: boolT(), Default: defValue("1"),
		Scope: "peer-control disclosure", Axes: []string{axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptInstructionSurface", Kind: CoupleOverriddenBy, Note: "no-op under any minimal* instruction surface (archaeology coupling #9)"},
		},
		Notes: "hides the control-surface explanation without disabling the surface; style exception: a 1|0 disclosure valve rather than visible|hidden",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go",
		},
	},
	{
		Name: "PromptImplDetails", Env: EnvPromptImplDetails,
		Type: boolT(), Default: defValue("0"),
		Scope: "implementation-physics disclosure", Axes: []string{axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "PromptInstructionSurface", Kind: CoupleOverriddenBy, Note: "no-op under any minimal* instruction surface (archaeology coupling #9)"},
		},
		Notes: "reveals fd/exec/detach implementation physics; style exception: a 1|0 disclosure valve rather than visible|hidden",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/runtime/prompt.go",
		},
	},
	{
		Name: "PromptBudgetVisibility", Env: EnvPromptBudgetVisibility,
		Type: enumT("visible", "hidden"), Default: defValue("visible"),
		Scope: "budget disclosure", Axes: []string{axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Notes: "despite the QUINE_PROMPT_* prefix, implementation lives in internal/world/state.go and governs tool-result / world --help budget framing rather than the system prompt; hidden is the strictest-floor toggle for zero-perception experiments; not loaded into Config and never validated — only the literal value hidden is significant, anything else behaves as visible",
		ImplSites: []string{
			"internal/world/state.go:LoadState",
			"internal/world/state.go:FormatBudgetedHelp",
		},
	},
	{
		Name: "NoMissionAutonomy", Env: EnvNoMissionAutonomy,
		Type: enumT("off", "autonomy", "sensing", "full"), Default: defValue("off"),
		Scope: "missionless opening posture", Axes: []string{axDirective}, Class: ClassDisclosure,
		Surfaces: surf(true, false, false, false), Mutability: MutExecBoundary,
		Notes: "legacy 0|1 accepted as off|full for backward compatibility; only affects no-mission prompts (silently irrelevant when a mission argv is present — archaeology coupling #10): autonomy adds an autonomous-process/do-not-wait clause, sensing adds visible-runtime/cwd attention, full composes both; supplies no task goal; propagates across child/re-entry envs for ablation consistency",
		ImplSites: []string{
			"internal/config/config.go:parseNoMissionAutonomy",
			"internal/runtime/prompt.go:buildOpeningIdentityBlock",
		},
	},
	{
		Name: "MemoryStrategyHints", Env: EnvMemoryStrategyHints,
		Type: boolT(), Default: defValue("1"),
		Scope: "memory pressure prompt", Axes: []string{axSelfRelation, axDisclosure}, Class: ClassDisclosure,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "AnchorMemory", Kind: CoupleSilentWithout, Note: "the pressure telemetry it reshapes only exists when the memory executor is constructed"},
		},
		Notes: "when off, gated pressure telemetry frames the goal neutrally instead of advising peer-delegation moves",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/tools/memory.go",
		},
	},

	// ------------------------------------------------------------------
	// Runtime evidence and telemetry
	// ------------------------------------------------------------------
	{
		Name: "FSMutationTelemetry", Env: EnvFSMutationTelemetry,
		Type: boolT(), Default: defValue("1"),
		Scope: "filesystem mutation telemetry", Axes: []string{axDisclosure}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "telemetry feature gate: hides fs_mutations prompt/schema language and omits fs_mutations / fs_mutations_so_far fields from tool results when 0; does not disable filesystem writes, revision tracking, or workspace physics",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/tools/tools.go",
		},
	},

	// ------------------------------------------------------------------
	// Response governance and execution envelope
	// ------------------------------------------------------------------
	{
		Name: "FailOnImpossible", Env: EnvFailOnImpossible,
		Type: boolT(), Default: defValue("1"),
		Scope: "impossibility posture", Axes: []string{axDirective}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "mixed triple-surface control (archaeology coupling #16): prompt posture, exit schema enum, and runtime validation of failure exits change together",
		ImplSites: []string{
			"internal/config/config.go:loadPromptAndTransportOptions",
			"internal/tools/schemas.go",
			"internal/runtime/tool_handlers.go",
		},
	},
	{
		Name: "MaxTurns", Env: EnvMaxTurns,
		Type: intT(), Default: defValue("0"),
		Scope: "execution budget", Axes: []string{axMortality}, Class: ClassBudget,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "0 disables the budget; sh is the only budget-consuming tool; exhaustion hard-fails the incarnation",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/runtime.go:handleExecutionBudgetExhaustion",
		},
	},
	{
		Name: "MaxDepth", Env: EnvMaxDepth,
		Type: intT(), Default: defValue("0"),
		Scope: "recursion limit", Axes: []string{axContinuity}, Class: ClassBudget,
		Surfaces: surf(true, false, true, false), Mutability: MutOperatorOnly,
		Couples: []Coupling{
			{Peer: "Depth", Kind: CoupleRequires, Note: "Load fails with ErrDepthExceeded when MaxDepth > 0 and Depth >= MaxDepth"},
		},
		Notes: "absent or 0 means no depth enforcement — it is not a budget of zero. Operator-only: a process-tree resource bound belongs to whoever launched the tree, so a lineage cannot relax its own limit through a mediated channel",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/tool_dispatch.go:precheckProcessCreation",
		},
	},
	{
		Name: "Depth", Env: EnvDepth,
		Type: intT(), Default: defValue("0"),
		Scope: "current process-tree depth", Axes: []string{axContinuity}, Class: ClassIdentity,
		Surfaces: surf(true, false, true, false), Mutability: MutSubstratePinned,
		Notes: "lineage counter: fork/spawn children are stamped Depth+1, and exec PRESERVES the current depth — exec is a new image of the same process, not a birth, and resetting it there would refill the MaxDepth budget enforcement reads from memory. Masked at the sh boundary: a program started from a shell is not a member of this agent tree. Pinned against agent staging even though Load() adopts inbound values, because faking it would fake lineage",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/config/envstamps.go:ForkChildStamps",
			"internal/config/envstamps.go:SelfReentryStamps",
		},
	},
	{
		Name: "MaxAgents", Env: EnvMaxAgents,
		Type: intT(), Default: defValue("0"),
		Scope: "process-tree slot limit", Axes: []string{axSocial}, Class: ClassBudget,
		Surfaces: surf(true, false, true, false), Mutability: MutOperatorOnly,
		Notes: "absent or 0 means no slot enforcement — it is not a budget of zero. Operator-only for the same reason as MaxDepth: the bound is on the tree, not on the individual",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/semaphore.go:NewAgentRegistry",
		},
	},
	{
		Name: "MaxConcurrent", Env: EnvMaxConcurrent,
		Type: intT(), Default: defValue("0"),
		Scope: "shared inference concurrency", Axes: []string{axSocial}, Class: ClassBudget,
		Surfaces: surf(false, false, true, false), Mutability: MutOperatorOnly,
		Notes: "absent or 0 means no semaphore enforcement — it is not a budget of zero. Operator-only: the inference slots are shared across the tree, so one member must not widen them for itself",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/semaphore.go:NewSemaphore",
		},
	},
	{
		Name: "ForkDefaultTimeout", Env: EnvForkDefaultTimeout,
		Type: intT(), Default: defValue("0"),
		Scope: "fork/spawn join envelope", Axes: []string{axContinuity}, Class: ClassBudget,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "ForkEnabled", Kind: CoupleShares, Note: "join deadline for the fork executor's wait/race joins"},
			{Peer: "SpawnEnabled", Kind: CoupleShares, Note: "spawn embeds the fork executor and inherits the same join deadline (archaeology coupling #18)"},
		},
		Notes: "max wall time (seconds) for fork/spawn wait and race joins; 0 disables the join deadline; must be >= 0 (load error otherwise)",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/tools/fork.go",
		},
	},
	{
		Name: "ShDefaultTimeout", Env: EnvShDefaultTimeout,
		Type: intT(), Default: defValue("300"),
		Scope: "shell envelope", Class: ClassBudget,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "per-call shell timeout (seconds); timeout is runtime physics, not prompt-only advice",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "ShTimeoutOverride", Env: EnvShTimeoutOverride,
		Type: boolT(), Default: defValue("1"),
		Scope: "shell envelope capability", Axes: []string{axDisclosure}, Class: ClassCapability,
		Surfaces: surf(true, true, true, false), Mutability: MutExecBoundary,
		Notes: "when disabled, removes and rejects per-call sh(timeout=...); the supported default timeout behavior still applies",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "WallClockExitSeconds", Env: EnvWallClockExitSeconds,
		Type: intT(), Default: defValue("0"),
		Scope: "process self-exit deadline", Axes: []string{axMortality}, Class: ClassBudget,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Notes: "opt-in wall-clock deadline that triggers graceful runtime self-exit (exit 0); 0 disables; must be >= 0 (load error otherwise); RETAINED per brief D7 (dorm family + bench-05 usage)",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/runtime.go",
		},
	},
	{
		Name: "OutputTruncate", Env: EnvOutputTruncate,
		Type: intT(), Default: defValue("20480"),
		Scope: "tool result truncation", Class: ClassBudget,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Notes: "caps retained tool output size (bytes)",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "EmptyAssistantSuccess", Env: EnvEmptyAssistantSuccess,
		Type: boolT(), Default: defValue("0"),
		Scope: "empty assistant success policy", Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Notes: "when enabled, an empty assistant response after prior tool/workspace progress is finalized as success rather than provider failure",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/ready_text.go:canTreatEmptyAssistantAsSuccess",
		},
	},
	{
		Name: "EphemeralBodyEnabled", Env: EnvEphemeralBodyEnabled,
		Type: boolT(), Default: defValue("0"),
		Scope: "launch-path body consumption", Axes: []string{axContinuity, axEnaction}, Class: ClassCapability,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "SelfReentryMode", Kind: CoupleHazard, Note: "under executable_path mode (or when the launch path is the reentry target) exec self-reentry fails after the body is unlinked until a replacement body exists; /proc/self/exe reentry survives an unlinked body (archaeology coupling #17)"},
		},
		Notes: "unlinks the launch path during startup; behavior/governance control, not prompt-only disclosure; the prompt warns about the reentry interaction",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"cmd/quine/main.go",
		},
	},
	{
		Name: "ReadyTextAutoIdle", Env: EnvReadyTextAutoIdle,
		Type: boolT(), Default: defValue("0"),
		Scope: "ready-text auto-suspend policy", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		ParentGate: "IdleEnabled",
		Notes:      "auto-triggers idle after a ready-like text-only response; silently no-ops unless QUINE_IDLE_ENABLED=1 (archaeology coupling #11); launch-local",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/ready_text.go:shouldReadyTextAutoIdle",
		},
	},
	{
		Name: "SuppressInitialBegin", Env: EnvSuppressInitialBegin,
		Type: boolT(), Default: defValue("0"),
		Scope: "synthetic initial-turn policy", Axes: []string{axDirective}, Class: ClassDisclosure,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "InitialUserMessage", Kind: CoupleShares, Note: "suppression compares the incoming material against the configured synthetic initial message to decide whether to omit the turn"},
		},
		Notes: "omits the synthetic Begin. user turn injected when stdin is a TTY (or empty pipe); launch-local",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/runtime.go",
		},
	},
	{
		Name: "InitialUserMessage", Env: EnvInitialUserMessage,
		Type: strT(), Default: defValue(""),
		Scope: "synthetic TTY initial user message", Axes: []string{axDirective}, Class: ClassDisclosure,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "SuppressInitialBegin", Kind: CoupleShares, Note: "the runtime compares against this value when deciding whether to flip hasMaterial and whether suppression should omit the turn"},
		},
		Notes: "overrides the legacy Begin. user turn injected when stdin is a TTY (or empty pipe); treated as synthetic, not material; propagated to children only when non-empty; the operator-injection successor of the deleted wisdom channel (brief D5)",
		ImplSites: []string{
			"internal/config/config.go:loadToolGates",
			"internal/runtime/runtime.go",
		},
	},

	// ------------------------------------------------------------------
	// Runtime paths, context pressure, and lineage
	// ------------------------------------------------------------------
	{
		Name: "ContextWindow", Env: EnvContextWindow,
		Type: intT(), Default: defValue("128000"),
		Scope: "model context budget", Axes: []string{axSelfRelation}, Class: ClassBudget,
		Surfaces: surf(true, false, true, true), Mutability: MutExecBoundary,
		Notes: "runtime pressure signal and provider request bound; also the derivation source for the memory warn/danger defaults",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/runtime.go",
		},
	},
	{
		Name: "MemoryWarnTokens", Env: EnvMemoryWarnTokens,
		Type: intT(), Default: defDerived("ContextWindow"),
		Scope: "memory pressure", Axes: []string{axSelfRelation}, Class: ClassBudget,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "AnchorMemory", Kind: CoupleSilentWithout, Note: "the memory executor that enforces pressure thresholds is only constructed when anchor memory is on"},
		},
		Notes: "must be > 0 (load error otherwise); default derives as ContextWindow/16 with a 2048 floor",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/config/config.go:defaultMemoryWarnTokens",
			"internal/tools/memory.go",
		},
	},
	{
		Name: "MemoryDangerTokens", Env: EnvMemoryDangerTokens,
		Type: intT(), Default: defDerived("ContextWindow"),
		Scope: "memory pressure", Axes: []string{axSelfRelation}, Class: ClassBudget,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "MemoryWarnTokens", Kind: CoupleRequires, Note: "must be greater than the warn threshold (load error otherwise)"},
			{Peer: "AnchorMemory", Kind: CoupleSilentWithout, Note: "the memory executor that enforces pressure thresholds is only constructed when anchor memory is on"},
		},
		Notes: "default derives as ContextWindow/8 with a 4096 floor",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/config/config.go:defaultMemoryDangerTokens",
			"internal/tools/memory.go",
		},
	},
	{
		Name: "MemoryDeathTokens", Env: EnvMemoryDeathTokens,
		Type: intT(), Default: defValue("0"),
		Scope: "memory pressure", Axes: []string{axMortality, axSelfRelation}, Class: ClassBudget,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "MemoryDangerTokens", Kind: CoupleRequires, Note: "when > 0 must be greater than the danger threshold (load error otherwise)"},
			{Peer: "AnchorMemory", Kind: CoupleSilentWithout, Note: "silently ignored when the memory executor is nil — no config-load error (archaeology coupling #7)"},
		},
		Notes: "0 disables; must be >= 0; crossing it terminates the incarnation (context_death)",
		ImplSites: []string{
			"internal/config/config.go:loadLimitConfig",
			"internal/runtime/runtime.go:memoryDeathStatus",
		},
	},
	{
		Name: "DataDir", Env: EnvDataDir,
		Type: strT(), Default: defValue(".quine/"),
		Scope: "runtime root", Class: ClassPath,
		Surfaces: surf(true, true, true, false), Mutability: MutOperatorOnly,
		Couples: []Coupling{
			{Peer: "WorkspaceRoot", Kind: CoupleRequires, Note: "must resolve outside the workspace root (load error otherwise)"},
		},
		Notes: "the durable runtime-state root (live session surfaces, pid routing, jobs, locks, overlay state); set-but-blank also resolves to .quine/; retained session state falls back to QUINE_DATA_DIR/log/<session>/ when QUINE_RETENTION_DIR is unset",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/config/config.go:RuntimeRoot",
		},
	},
	{
		Name: "RetentionDir", Env: EnvRetentionDir,
		Type: strT(), Default: defValue(""),
		Scope: "retained owner root", Class: ClassPath,
		Surfaces: surf(false, false, true, false), Mutability: MutOperatorOnly,
		Couples: []Coupling{
			{Peer: "WorkspaceRoot", Kind: CoupleRequires, Note: "must resolve outside the workspace root (load error otherwise)"},
		},
		Notes: "when set, canonical retained session state lives under QUINE_RETENTION_DIR/sessions/<session>; propagated to children only when non-empty",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/config/config.go:SessionRetainedDir",
		},
	},
	{
		Name: "WorkDir", Env: EnvWorkDir,
		Type: strT(), Default: defDerived("Workspace"),
		Scope: "shell cwd default", Class: ClassPath,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Notes: "startup shell cwd; distinct from workspace root; unset defaults to the current workspace when workspace physics are enabled, else to the startup cwd",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "Shell", Env: EnvShell,
		Type: strT(), Default: defValue("/bin/sh"),
		Scope: "shell binary", Class: ClassPath,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Notes: "binary used for sh tool execution; not validated at load (existence is checked at execution time)",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/tools/tools.go:NewShExecutor",
		},
	},
	{
		Name: "ShNetwork", Env: EnvShNetwork,
		Type: enumT("host", "none"), Default: defValue("host"),
		Scope: "shell network namespace", Axes: []string{axEnaction}, Class: ClassCapability,
		Surfaces: surf(true, false, true, false), Mutability: MutExecBoundary,
		Notes: "host leaves shell networking unchanged; none requests an isolated network namespace (CLONE_NEWNET) while preserving local loopback when supported; none is Linux-only and config validation fails closed on non-Linux rather than silently granting host network (archaeology coupling #4); RETAINED per brief D7 (bench-04 default)",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/tools/job_attrs_linux.go",
		},
	},
	{
		Name: "SelfReentryMode", Env: EnvSelfReentryMode,
		Type: enumT("self", "executable_path"), Default: defValue("self"),
		Scope: "default fork / exec re-entry mode", Axes: []string{axContinuity, axEnaction}, Class: ClassCapability,
		Surfaces: surf(false, false, true, false), Mutability: MutExecBoundary,
		Couples: []Coupling{
			{Peer: "ExecEnabled", Kind: CoupleSilentWithout, Note: "the reentry target is only resolved when exec, fork, or spawn is enabled (any one suffices; archaeology coupling #13)"},
			{Peer: "ForkEnabled", Kind: CoupleSilentWithout, Note: "the reentry target is only resolved when exec, fork, or spawn is enabled (any one suffices; archaeology coupling #13)"},
			{Peer: "SpawnEnabled", Kind: CoupleSilentWithout, Note: "the reentry target is only resolved when exec, fork, or spawn is enabled (any one suffices; archaeology coupling #13)"},
			{Peer: "EphemeralBodyEnabled", Kind: CoupleHazard, Note: "executable_path reentry fails after the ephemeral body is unlinked until a replacement exists; self (/proc/self/exe) survives an unlinked body (archaeology coupling #17)"},
		},
		Notes: "self resolves to /proc/self/exe and is Linux-only (load error elsewhere when a lifecycle gate is on); executable_path resolves to the current executable path and is the non-Linux opt-in mode",
		ImplSites: []string{
			"internal/config/config.go:loadSelfReentryConfig",
			"internal/config/config.go:resolveSelfReentryTarget",
		},
	},
	{
		Name: "SelfReentryTarget", Env: EnvSelfReentryTarget,
		Type: strT(), Default: defLegacy("removed-error"),
		Scope: "legacy", Class: ClassPath,
		Surfaces: surf(false, false, true, false), Mutability: MutSubstratePinned,
		Notes: "removed: setting any non-empty value is a config-load error; use QUINE_SELF_REENTRY_MODE instead",
		ImplSites: []string{
			"internal/config/config.go:loadSelfReentryConfig",
		},
	},
	{
		Name: "SessionID", Env: EnvSessionID,
		Type: strT(), Default: defRuntime(),
		Scope: "stable session identity", Axes: []string{axSocial}, Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "auto-generated when unset; passing it at bootstrap is the explicit resume/adoption mechanism. Masked at every boundary and re-stamped where it means something: exec carries it forward (same quine, new image), while each fork/spawn child is injected with its own by the fork executor, and an sh child gets none at all",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/config/envstamps.go:SelfReentryStamps",
			"internal/tools/fork.go:launchChildSession",
		},
	},
	{
		Name: "RunID", Env: EnvRunID,
		Type: strT(), Default: defRuntime(),
		Scope: "physical run identity", Class: ClassIdentity,
		Surfaces: surf(true, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "always regenerated on every process activation — stale inbound env values are ignored; exported to tools for the current process activation; children/re-entry generate new run ids",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/config/config.go:processRunID",
		},
	},
	{
		Name: "TapeID", Env: EnvTapeID,
		Type: strT(), Default: defRuntime(),
		Scope: "tape/log identity", Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "auto-incremented when unset (next numeric id in the session tape dir); re-stamped at exec for tape continuity; internal trace identity, not a live context contract. Masked at the sh and fork boundaries as process-private",
		ImplSites: []string{
			"internal/config/config.go:loadIdentityAndPathConfig",
			"internal/config/config.go:nextTapeID",
			"internal/config/envstamps.go:SelfReentryStamps",
		},
	},
	{
		Name: "ContextTape", Env: EnvContextTape,
		Type: strT(), Default: defRuntime(),
		Scope: "exec-only process identity", Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "runtime-emitted, process-private; not read by Load() at all — it exists to stage context across the exec barrier. Masked at every boundary by its runtime-emitted mutability, so no process this runtime builds can inherit another's; runtime-emitted entrypoints are owned by runtime-surface.md",
		ImplSites: []string{
			"internal/config/envmodel.go:BoundaryBehavior",
			"internal/runtime/tool_handlers.go",
		},
	},
	{
		Name: "ContextBootstrap", Env: EnvContextBootstrap,
		Type: strT(), Default: defRuntime(),
		Scope: "context handover signal", Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutRuntimeEmitted,
		Notes: "runtime-emitted, process-private; not read by Load() at all — the runtime points a successor or child at a staged context tree with it and the receiver unsets it after import, so it must never be inherited by an unrelated process. Registered as a knob (rather than a hardcoded exception) so the env mask derives from Mutability alone",
		ImplSites: []string{
			"internal/tools/fs_copy.go",
			"internal/tools/exec.go:Execute",
			"internal/runtime/runtime.go:importBootstrappedContext",
		},
	},
	{
		Name: "ParentSession", Env: EnvParentSession,
		Type: strT(), Default: defValue(""),
		Scope: "parent/child lineage link", Axes: []string{axSocial}, Class: ClassIdentity,
		Surfaces: surf(false, false, true, false), Mutability: MutSubstratePinned,
		Notes: "process-tree linkage surface: fork/spawn children are stamped with the current SessionID as parent, and exec preserves the existing one (a new image does not acquire a new parent). Masked at the sh boundary — a program started from a shell is not a member of this tree. Pinned against agent staging: faking it would fake lineage, even though Load() adopts inbound values",
		ImplSites: []string{
			"internal/config/config.go:loadWorkspaceConfig",
			"internal/config/envstamps.go:ForkChildStamps",
			"internal/config/envstamps.go:SelfReentryStamps",
		},
	},
}

// KnobByEnv returns the registry entry for an env name.
func KnobByEnv(env string) (Knob, bool) {
	for _, k := range Registry {
		if k.Env == env {
			return k, true
		}
	}
	return Knob{}, false
}

// KnobByName returns the registry entry for a stable knob Name.
func KnobByName(name string) (Knob, bool) {
	for _, k := range Registry {
		if k.Name == name {
			return k, true
		}
	}
	return Knob{}, false
}

// RegistryJSON renders the full registry table as indented JSON (camelCase
// keys), the exact payload the Phase-2 config/registry.json materialization
// will write.
func RegistryJSON() ([]byte, error) {
	return json.MarshalIndent(Registry, "", "  ")
}
