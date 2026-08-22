package config

// registry_test.go is the registry <-> Load() agreement test (work order
// T1.1, brief D1). It pins four properties:
//
//  1. Bijection: every Env* constant in envnames.go has exactly one registry
//     entry, and every registry Env is one of those constants (the
//     compiler-checked registryEnvUniverse list below is the checklist).
//  2. Load() round-trip: for every registry knob with a loadable Config
//     field, a type-appropriate non-default env value resolves to the
//     matching field; enum knobs accept every legal value (with documented
//     alias normalization) and reject illegal ones exactly as Load() does
//     today.
//  3. Defaults: knobs with DefaultSpec kind "value" resolve to that value on
//     an empty (required-only) envp; derived memory defaults follow the
//     documented ContextWindow formulas.
//  4. Reference integrity: every Couples.Peer, ParentGate, and Default.From
//     names an existing registry entry.
//
// The registry only DESCRIBES Load(); these tests exist to fail loudly when
// either side drifts.

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// registryEnvUniverse enumerates every Env* constant in envnames.go, by
// reference, so the compiler enforces existence and deletion pressure.
var registryEnvUniverse = []string{
	EnvModelID,
	EnvAPIType,
	EnvAPIBase,
	EnvAPIKey,
	EnvMaxDepth,
	EnvDepth,
	EnvSessionID,
	EnvRunID,
	EnvTapeID,
	EnvParentSession,
	EnvMaxConcurrent,
	EnvMaxAgents,
	EnvForkDefaultTimeout,
	EnvShDefaultTimeout,
	EnvShTimeoutOverride,
	EnvShStdinEnabled,
	EnvShDetachEnabled,
	EnvOutputTruncate,
	EnvDataDir,
	EnvRetentionDir,
	EnvWorkDir,
	EnvShell,
	EnvShNetwork,
	EnvSelfReentryMode,
	EnvSelfReentryTarget,
	EnvMaxTurns,
	EnvWallClockExitSeconds,
	EnvPromptMetaphor,
	EnvPromptSelfModel,
	EnvPromptInstructionSurface,
	EnvPromptRuntimeSurface,
	EnvPromptPersona,
	EnvPromptCtl,
	EnvPromptImplDetails,
	EnvPromptBudgetVisibility,
	EnvPeerDiscoveryEnabled,
	EnvPeerDiscoveryHeartbeat,
	EnvFSMutationTelemetry,
	EnvFailOnImpossible,
	EnvNoMissionAutonomy,
	EnvEmptyAssistantSuccess,
	EnvReadyTextAutoIdle,
	EnvContextWindow,
	EnvThinkingBudget,
	EnvModelServiceTier,
	EnvMemoryWarnTokens,
	EnvMemoryDangerTokens,
	EnvMemoryDeathTokens,
	EnvMemoryStrategyHints,
	EnvConfigDir,
	EnvAnchorMemory,
	EnvAnchorFoldEnabled,
	EnvAnchorMarkEnabled,
	EnvIdleEnabled,
	EnvForkEnabled,
	EnvExitEnabled,
	EnvExecEnabled,
	EnvSpawnEnabled,
	EnvAgentsMDEnabled,
	EnvAgentsSkillsEnabled,
	EnvVisionEnabled,
	EnvShInteractiveEnabled,
	EnvForkWorldEnabled,
	EnvEphemeralBodyEnabled,
	EnvSuppressInitialBegin,
	EnvInitialUserMessage,
	EnvSelfSourceCodeEnabled,
	EnvSelfSourceProjection,
	EnvUserAgent,
	EnvContextTape,
	EnvContextBootstrap,
	EnvWorkspaceRoot,
	EnvWorkspace,
	EnvWorkspaceBackend,
	EnvWorkspaceOverlayDriver,
	EnvWorkspaceRevisionMode,
	EnvWorkspaceSession,
	EnvWorkspaceOwner,
	EnvWorkspaceCommitOnSignal,
	EnvWorkspaceBootstrap,
	EnvWorkspaceCurrentRevision,
	EnvWorkspaceSource,
	EnvWorldOnePerShell,
	EnvDebugRequestBodyDir,
}

// nonLoadedRegistryEnvs are registry entries deliberately without a Load()
// round-trip probe; each has a dedicated behavior test below.
var nonLoadedRegistryEnvs = map[string]string{
	EnvRunID:                  "regenerated on every activation; TestRegistryRuntimeEmittedLoadBehavior",
	EnvContextTape:            "never read by Load(); TestRegistryRuntimeEmittedLoadBehavior",
	EnvContextBootstrap:       "never read by Load(); consumed by the runtime's context import and unset after; TestRegistryRuntimeEmittedLoadBehavior",
	EnvConfigDir:              "not stored on Config; baseEnv passthrough; TestRegistryConfigDirPassthrough",
	EnvPromptBudgetVisibility: "read directly by internal/world, unvalidated; TestRegistryUnvalidatedEnvsLoadCleanly",
	EnvWorldOnePerShell:       "read directly by internal/world, unvalidated; TestRegistryUnvalidatedEnvsLoadCleanly",
	EnvSelfReentryTarget:      "removed tombstone; TestRegistryLegacyTombstonesRejected",
	EnvWorkspaceSource:        "removed tombstone; TestRegistryLegacyTombstonesRejected",
}

// loadProbe drives one knob through Load().
type loadProbe struct {
	set             string            // type-appropriate non-default value
	want            string            // expected resolved value; defaults to set
	pre             map[string]string // prerequisite envs
	linuxOnly       bool              // whole probe needs Linux
	linuxOnlyValues map[string]bool   // enum values that need Linux
	canon           map[string]string // accepted input -> resolved value (aliases)
	get             func(*Config) string
}

func b01get(f func(*Config) bool) func(*Config) string {
	return func(c *Config) string { return bool01(f(c)) }
}

func intGet(f func(*Config) int) func(*Config) string {
	return func(c *Config) string { return strconv.Itoa(f(c)) }
}

func mustCanonical(t *testing.T, p string) string {
	t.Helper()
	got, err := canonicalPath(p)
	if err != nil {
		t.Fatalf("canonicalPath(%q): %v", p, err)
	}
	return got
}

// buildLoadProbes returns the env -> probe table. Every registry entry must
// appear either here or in nonLoadedRegistryEnvs.
func buildLoadProbes(t *testing.T) map[string]loadProbe {
	t.Helper()

	wsRoot := t.TempDir()
	wsRootCanon := mustCanonical(t, wsRoot)
	wsSub := wsRoot + string(os.PathSeparator) + "sub"
	if err := os.MkdirAll(wsSub, 0o755); err != nil {
		t.Fatalf("mkdir workspace sub: %v", err)
	}
	wsSubCanon := mustCanonical(t, wsSub)
	directWS := map[string]string{EnvWorkspaceRoot: wsRoot, EnvWorkspaceBackend: "direct"}
	overlayWS := map[string]string{EnvWorkspaceRoot: wsRoot, EnvWorkspaceBackend: "overlay"}
	scratchWork := t.TempDir()

	return map[string]loadProbe{
		// --- transport / identity ---
		EnvModelID: {set: "registry-test-model", get: func(c *Config) string { return c.ModelID }},
		EnvAPIType: {set: "openai", get: func(c *Config) string { return c.Provider }},
		EnvAPIBase: {set: "https://registry.example.test", get: func(c *Config) string { return c.APIBase }},
		EnvAPIKey:  {set: "sk-registry-test", get: func(c *Config) string { return c.APIKey }},
		EnvUserAgent: {set: "quine-registry-test/1.0",
			get: func(c *Config) string { return c.UserAgent }},
		EnvThinkingBudget: {set: "low",
			get: func(c *Config) string { return c.ThinkingBudget }},
		EnvModelServiceTier: {set: "flex", canon: map[string]string{"fast": "priority"},
			get: func(c *Config) string { return c.ServiceTier }},
		EnvDebugRequestBodyDir: {set: "/tmp/quine-registry-debug",
			get: func(c *Config) string { return c.DebugRequestBodyDir }},

		// --- lineage identity ---
		EnvSessionID:     {set: "sess-registry-roundtrip", get: func(c *Config) string { return c.SessionID }},
		EnvTapeID:        {set: "tape-registry-7", get: func(c *Config) string { return c.TapeID }},
		EnvParentSession: {set: "parent-registry", get: func(c *Config) string { return c.ParentSession }},
		EnvDepth:         {set: "3", get: intGet(func(c *Config) int { return c.Depth })},
		EnvMaxDepth:      {set: "10", get: intGet(func(c *Config) int { return c.MaxDepth })},

		// --- limits / envelopes ---
		EnvMaxConcurrent:      {set: "50", get: intGet(func(c *Config) int { return c.MaxConcurrent })},
		EnvMaxAgents:          {set: "25", get: intGet(func(c *Config) int { return c.MaxAgents })},
		EnvForkDefaultTimeout: {set: "45", get: intGet(func(c *Config) int { return c.ForkDefaultTimeoutSeconds })},
		EnvShDefaultTimeout:   {set: "60", get: intGet(func(c *Config) int { return c.ShTimeout })},
		EnvOutputTruncate:     {set: "4096", get: intGet(func(c *Config) int { return c.OutputTruncate })},
		EnvMaxTurns:           {set: "30", get: intGet(func(c *Config) int { return c.MaxTurns })},
		EnvWallClockExitSeconds: {set: "870",
			get: intGet(func(c *Config) int { return c.WallClockExitSeconds })},
		EnvContextWindow: {set: "64000", get: intGet(func(c *Config) int { return c.ContextWindow })},
		EnvMemoryWarnTokens: {set: "7000",
			get: intGet(func(c *Config) int { return c.MemoryWarnTokens })},
		EnvMemoryDangerTokens: {set: "20000",
			get: intGet(func(c *Config) int { return c.MemoryDangerTokens })},
		EnvMemoryDeathTokens: {set: "30000",
			get: intGet(func(c *Config) int { return c.MemoryDeathTokens })},
		EnvPeerDiscoveryHeartbeat: {set: "1234",
			get: intGet(func(c *Config) int { return c.PeerDiscoveryHeartbeatMS })},

		// --- tool gates ---
		EnvShTimeoutOverride: {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.ShTimeoutOverrideEnabled })},
		EnvShStdinEnabled:    {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.ShStdinEnabled })},
		EnvShDetachEnabled:   {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.ShDetachEnabled })},
		EnvShInteractiveEnabled: {set: "0",
			get: b01get(func(c *Config) bool { return c.ToolGates.ShInteractiveEnabled })},
		EnvForkEnabled:   {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.ForkEnabled })},
		EnvExitEnabled:   {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.ExitEnabled })},
		EnvExecEnabled:   {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.ExecEnabled })},
		EnvSpawnEnabled:  {set: "1", get: b01get(func(c *Config) bool { return c.ToolGates.SpawnEnabledFlag })},
		EnvIdleEnabled:   {set: "1", get: b01get(func(c *Config) bool { return c.ToolGates.IdleEnabled })},
		EnvVisionEnabled: {set: "0", get: b01get(func(c *Config) bool { return c.ToolGates.VisionEnabled })},
		EnvAnchorMemory:  {set: "1", get: b01get(func(c *Config) bool { return c.ToolGates.AnchorMemoryEnabled })},
		EnvAnchorFoldEnabled: {set: "0",
			get: b01get(func(c *Config) bool { return c.ToolGates.AnchorFoldEnabled })},
		EnvAnchorMarkEnabled: {set: "0",
			get: b01get(func(c *Config) bool { return c.ToolGates.AnchorMarkEnabled })},
		EnvAgentsMDEnabled: {set: "1", pre: map[string]string{EnvWorkDir: scratchWork},
			get: b01get(func(c *Config) bool { return c.ToolGates.AgentsMDEnabled })},
		EnvAgentsSkillsEnabled: {set: "1", pre: map[string]string{EnvWorkDir: scratchWork},
			get: b01get(func(c *Config) bool { return c.ToolGates.AgentsSkillsEnabled })},
		EnvSelfSourceCodeEnabled: {set: "1",
			get: b01get(func(c *Config) bool { return c.ToolGates.SelfSourceCodeEnabled })},
		EnvSelfSourceProjection: {set: "runtime",
			get: func(c *Config) string { return c.SelfSourceProjectionMode() }},
		EnvPeerDiscoveryEnabled: {set: "1",
			get: b01get(func(c *Config) bool { return c.ToolGates.PeerDiscoveryEnabled })},
		EnvFSMutationTelemetry: {set: "0",
			get: b01get(func(c *Config) bool { return c.ToolGates.FSMutationTelemetry })},
		EnvEmptyAssistantSuccess: {set: "1",
			get: b01get(func(c *Config) bool { return c.ToolGates.EmptyAssistantSuccess })},
		EnvReadyTextAutoIdle: {set: "1",
			get: b01get(func(c *Config) bool { return c.ToolGates.ReadyTextAutoIdle })},
		EnvEphemeralBodyEnabled: {set: "1",
			get: b01get(func(c *Config) bool { return c.ToolGates.EphemeralBody })},
		EnvSuppressInitialBegin: {set: "1",
			get: b01get(func(c *Config) bool { return c.ToolGates.SuppressInitialBegin })},
		EnvInitialUserMessage: {set: "wake up",
			get: func(c *Config) string { return c.InitialUserMessage }},

		// --- prompt disclosure ---
		EnvPromptMetaphor: {set: "thermodynamic",
			get: func(c *Config) string { return string(c.PromptMetaphor) }},
		EnvPromptSelfModel: {set: "basic",
			get: func(c *Config) string { return string(c.PromptSelfModel) }},
		EnvPromptInstructionSurface: {set: "minimal",
			get: func(c *Config) string { return string(c.PromptInstructionSurface) }},
		EnvPromptRuntimeSurface: {set: "hidden",
			get: func(c *Config) string { return string(c.PromptRuntimeSurface) }},
		EnvPromptPersona: {set: "architect",
			get: func(c *Config) string { return string(c.PromptPersona) }},
		EnvPromptCtl:         {set: "0", get: b01get(func(c *Config) bool { return c.PromptCtlPhysics })},
		EnvPromptImplDetails: {set: "1", get: b01get(func(c *Config) bool { return c.PromptImplDetails })},
		EnvFailOnImpossible:  {set: "0", get: b01get(func(c *Config) bool { return c.FailOnImpossible })},
		EnvNoMissionAutonomy: {set: "sensing", canon: map[string]string{"0": "off", "1": "full"},
			get: func(c *Config) string { return c.NoMissionAutonomy }},
		EnvMemoryStrategyHints: {set: "0",
			get: b01get(func(c *Config) bool { return c.MemoryStrategyHints })},

		// --- paths / shell ---
		EnvDataDir:      {set: "/tmp/quine-registry-data", get: func(c *Config) string { return c.DataDir }},
		EnvRetentionDir: {set: "/tmp/quine-registry-retained", get: func(c *Config) string { return c.RetentionDir }},
		EnvWorkDir:      {set: scratchWork, get: func(c *Config) string { return c.WorkDir }},
		EnvShell:        {set: "/bin/bash", get: func(c *Config) string { return c.Shell }},
		EnvShNetwork: {set: "none", linuxOnlyValues: map[string]bool{"none": true},
			linuxOnly: runtime.GOOS != "linux",
			get:       func(c *Config) string { return c.ShNetwork }},
		EnvSelfReentryMode: {set: "executable_path", linuxOnlyValues: map[string]bool{"self": true},
			get: func(c *Config) string { return string(c.SelfReentryMode) }},

		// --- workspace physics ---
		EnvWorkspaceRoot: {set: wsRoot, want: wsRootCanon, pre: map[string]string{EnvWorkspaceBackend: "direct"},
			get: func(c *Config) string { return c.WorkspaceRoot }},
		EnvWorkspace: {set: wsSub, want: wsSubCanon, pre: directWS,
			get: func(c *Config) string { return c.Workspace }},
		EnvWorkspaceBackend: {set: "direct", pre: map[string]string{EnvWorkspaceRoot: wsRoot},
			linuxOnlyValues: map[string]bool{"overlay": true},
			get:             func(c *Config) string { return c.WorkspaceBackend }},
		EnvWorkspaceOverlayDriver: {set: "fuse", pre: overlayWS, linuxOnly: true,
			get: func(c *Config) string { return c.WorkspaceOverlayDriver }},
		EnvWorkspaceRevisionMode: {set: "none", pre: overlayWS, linuxOnly: true,
			get: func(c *Config) string { return string(c.WorkspaceRevisionMode) }},
		EnvWorkspaceSession: {set: "wsess-registry", pre: directWS,
			get: func(c *Config) string { return c.WorkspaceSession }},
		EnvWorkspaceOwner: {set: "0",
			get: b01get(func(c *Config) bool { return c.WorkspaceOwner })},
		EnvWorkspaceCommitOnSignal: {set: "1",
			get: b01get(func(c *Config) bool { return c.WorkspaceCommitOnSignal })},
		EnvWorkspaceBootstrap: {set: "boot-lineage-registry",
			get: func(c *Config) string { return c.WorkspaceBootstrap }},
		EnvWorkspaceCurrentRevision: {set: "rev-42",
			get: func(c *Config) string { return c.WorkspaceCurrentRevision }},
		EnvForkWorldEnabled: {set: "1", pre: directWS,
			get: b01get(func(c *Config) bool { return c.ToolGates.ForkWorldEnabled })},
	}
}

// resetRegistryEnv unsets every registry env (restoring on cleanup), then
// applies the 4 required envs the way the package's other tests do.
func resetRegistryEnv(t *testing.T) {
	t.Helper()
	saved := make(map[string]string)
	for _, k := range Registry {
		if v, ok := os.LookupEnv(k.Env); ok {
			saved[k.Env] = v
		}
		os.Unsetenv(k.Env)
	}
	t.Cleanup(func() {
		for _, k := range Registry {
			if v, ok := saved[k.Env]; ok {
				os.Setenv(k.Env, v)
			} else {
				os.Unsetenv(k.Env)
			}
		}
	})
	setRequired(t)
}

// requiredBaselineEnvs are set by setRequired and therefore excluded from
// default-value assertions.
func requiredBaselineEnvs() map[string]bool {
	req := map[string]bool{
		EnvModelID: true, EnvAPIType: true, EnvAPIBase: true, EnvAPIKey: true,
	}
	if runtime.GOOS != "linux" {
		req[EnvSelfReentryMode] = true
	}
	return req
}

// --- 1. Bijection ---

func TestRegistryBijectionWithEnvNames(t *testing.T) {
	if len(registryEnvUniverse) != len(Registry) {
		t.Errorf("registry has %d entries, envnames.go universe has %d", len(Registry), len(registryEnvUniverse))
	}
	universe := make(map[string]bool, len(registryEnvUniverse))
	for _, env := range registryEnvUniverse {
		if universe[env] {
			t.Errorf("duplicate env %q in universe list", env)
		}
		universe[env] = true
	}
	seen := make(map[string]bool, len(Registry))
	for _, k := range Registry {
		if seen[k.Env] {
			t.Errorf("duplicate registry entry for env %q", k.Env)
		}
		seen[k.Env] = true
		if !universe[k.Env] {
			t.Errorf("registry entry %q (%s) does not correspond to an Env* constant", k.Name, k.Env)
		}
	}
	for env := range universe {
		if !seen[env] {
			t.Errorf("envnames.go constant %q has no registry entry", env)
		}
	}
}

func TestRegistryProcessIdentityMutability(t *testing.T) {
	identity := make(map[string]bool, len(ProcessIdentityEnvNames))
	for _, env := range ProcessIdentityEnvNames {
		identity[env] = true
		k, ok := KnobByEnv(env)
		if !ok {
			t.Errorf("process-identity env %q missing from registry", env)
			continue
		}
		if k.Mutability != MutRuntimeEmitted {
			t.Errorf("%s: Mutability = %q, want %q (ProcessIdentityEnvNames member)", env, k.Mutability, MutRuntimeEmitted)
		}
	}
	for _, k := range Registry {
		if k.Mutability == MutRuntimeEmitted && !identity[k.Env] {
			t.Errorf("%s: Mutability runtime-emitted but not in ProcessIdentityEnvNames", k.Env)
		}
	}
}

// --- 4. Reference integrity and schema hygiene ---

func TestRegistryReferencesResolve(t *testing.T) {
	for _, k := range Registry {
		for _, c := range k.Couples {
			peer, ok := KnobByName(c.Peer)
			if !ok {
				t.Errorf("%s: coupling peer %q is not a registry Name", k.Name, c.Peer)
			} else if peer.Name == k.Name {
				t.Errorf("%s: coupling points at itself", k.Name)
			}
			if c.Kind == "" || c.Note == "" {
				t.Errorf("%s -> %s: coupling must carry kind and note", k.Name, c.Peer)
			}
		}
		if k.ParentGate != "" {
			if _, ok := KnobByName(k.ParentGate); !ok {
				t.Errorf("%s: ParentGate %q is not a registry Name", k.Name, k.ParentGate)
			}
		}
		if k.Default.From != "" {
			if k.Default.Kind != DefaultDerived {
				t.Errorf("%s: Default.From set but kind is %q", k.Name, k.Default.Kind)
			}
			if _, ok := KnobByName(k.Default.From); !ok {
				t.Errorf("%s: Default.From %q is not a registry Name", k.Name, k.Default.From)
			}
		} else if k.Default.Kind == DefaultDerived {
			t.Errorf("%s: derived default without From", k.Name)
		}
	}
}

func TestRegistrySchemaHygiene(t *testing.T) {
	classes := map[GateClass]bool{
		ClassCapability: true, ClassDisclosure: true, ClassBudget: true,
		ClassIdentity: true, ClassTransport: true, ClassPath: true,
	}
	muts := map[Mutability]bool{
		MutSubstratePinned: true, MutExecBoundary: true,
		MutOperatorOnly: true, MutRuntimeEmitted: true,
	}
	typeKinds := map[TypeKind]bool{TypeBool: true, TypeInt: true, TypeString: true, TypeEnum: true}
	defKinds := map[DefaultKind]bool{
		DefaultRequired: true, DefaultValue: true, DefaultDerived: true,
		DefaultRuntimeEmitted: true, DefaultLegacy: true, DefaultExternalLabel: true,
	}
	axes := map[string]bool{
		axMortality: true, axContinuity: true, axDisclosure: true, axSelfRelation: true,
		axDirective: true, axSocial: true, axEnaction: true,
	}
	names := make(map[string]bool, len(Registry))
	for _, k := range Registry {
		if k.Name == "" || k.Env == "" || k.Scope == "" {
			t.Errorf("%s/%s: Name, Env, and Scope are mandatory", k.Name, k.Env)
		}
		if names[k.Name] {
			t.Errorf("duplicate registry Name %q", k.Name)
		}
		names[k.Name] = true
		if !strings.HasPrefix(k.Env, "QUINE_") {
			t.Errorf("%s: env %q lacks QUINE_ prefix", k.Name, k.Env)
		}
		if !classes[k.Class] {
			t.Errorf("%s: unknown Class %q", k.Name, k.Class)
		}
		if !muts[k.Mutability] {
			t.Errorf("%s: unknown Mutability %q", k.Name, k.Mutability)
		}
		if !typeKinds[k.Type.Kind] {
			t.Errorf("%s: unknown Type.Kind %q", k.Name, k.Type.Kind)
		}
		if (k.Type.Kind == TypeEnum) != (len(k.Type.Enum) > 0) {
			t.Errorf("%s: Enum values must be present exactly for enum knobs", k.Name)
		}
		if !defKinds[k.Default.Kind] {
			t.Errorf("%s: unknown Default.Kind %q", k.Name, k.Default.Kind)
		}
		for _, ax := range k.Axes {
			if !axes[ax] {
				t.Errorf("%s: unknown axis tag %q", k.Name, ax)
			}
		}
		if len(k.ImplSites) == 0 {
			t.Errorf("%s: ImplSites must not be empty", k.Name)
		}
	}
	data, err := RegistryJSON()
	if err != nil {
		t.Fatalf("RegistryJSON() error: %v", err)
	}
	if len(data) == 0 || data[0] != '[' {
		t.Fatalf("RegistryJSON() produced unexpected payload")
	}
}

// --- 2. Load() round-trip ---

func TestRegistryLoadRoundTrip(t *testing.T) {
	probes := buildLoadProbes(t)
	for env := range probes {
		if _, ok := KnobByEnv(env); !ok {
			t.Errorf("probe for %q has no registry entry", env)
		}
	}
	for _, k := range Registry {
		probe, ok := probes[k.Env]
		if !ok {
			if _, listed := nonLoadedRegistryEnvs[k.Env]; !listed {
				t.Errorf("%s: no round-trip probe and not declared non-loaded", k.Env)
			}
			continue
		}
		k := k
		t.Run(k.Env, func(t *testing.T) {
			if probe.linuxOnly && runtime.GOOS != "linux" {
				t.Skipf("%s probe requires Linux", k.Env)
			}
			if probe.linuxOnlyValues[probe.set] && runtime.GOOS != "linux" {
				t.Skipf("%s=%s requires Linux", k.Env, probe.set)
			}
			resetRegistryEnv(t)
			for pk, pv := range probe.pre {
				os.Setenv(pk, pv)
			}
			os.Setenv(k.Env, probe.set)
			c, err := Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}
			want := probe.want
			if want == "" {
				want = probe.set
			}
			if got := probe.get(c); got != want {
				t.Errorf("%s=%q resolved to %q, want %q", k.Env, probe.set, got, want)
			}
		})
	}
}

// --- 3. Defaults ---

func TestRegistryValueDefaults(t *testing.T) {
	probes := buildLoadProbes(t)
	skip := requiredBaselineEnvs()
	resetRegistryEnv(t)
	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	for _, k := range Registry {
		if k.Default.Kind != DefaultValue || skip[k.Env] {
			continue
		}
		probe, ok := probes[k.Env]
		if !ok {
			continue // non-loaded knobs (documented in nonLoadedRegistryEnvs)
		}
		if got := probe.get(c); got != k.Default.Value {
			t.Errorf("%s: empty env resolved to %q, registry default says %q", k.Env, got, k.Default.Value)
		}
	}
}

func TestRegistryDerivedMemoryDefaults(t *testing.T) {
	resetRegistryEnv(t)
	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	// ContextWindow default 128000: warn = /16, danger = /8.
	if c.MemoryWarnTokens != 8000 {
		t.Errorf("MemoryWarnTokens = %d, want 8000 (ContextWindow/16)", c.MemoryWarnTokens)
	}
	if c.MemoryDangerTokens != 16000 {
		t.Errorf("MemoryDangerTokens = %d, want 16000 (ContextWindow/8)", c.MemoryDangerTokens)
	}

	resetRegistryEnv(t)
	os.Setenv(EnvContextWindow, "16000")
	c, err = Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	// Small windows hit the documented floors.
	if c.MemoryWarnTokens != 2048 {
		t.Errorf("MemoryWarnTokens = %d, want floor 2048", c.MemoryWarnTokens)
	}
	if c.MemoryDangerTokens != 4096 {
		t.Errorf("MemoryDangerTokens = %d, want floor 4096", c.MemoryDangerTokens)
	}
}

// --- enum agreement: every legal value parses, aliases normalize, illegal rejected ---

func TestRegistryEnumAgreement(t *testing.T) {
	probes := buildLoadProbes(t)
	for _, k := range Registry {
		if k.Type.Kind != TypeEnum {
			continue
		}
		probe, ok := probes[k.Env]
		if !ok {
			continue // lenient, non-loaded enums have dedicated tests
		}
		k, probe := k, probe
		t.Run(k.Env, func(t *testing.T) {
			if probe.linuxOnly && runtime.GOOS != "linux" {
				t.Skipf("%s probe requires Linux", k.Env)
			}
			inputs := make(map[string]string, len(k.Type.Enum)+len(probe.canon))
			for _, v := range k.Type.Enum {
				inputs[v] = v
			}
			for in, out := range probe.canon {
				inputs[in] = out
			}
			for in, wantOut := range inputs {
				if probe.linuxOnlyValues[in] && runtime.GOOS != "linux" {
					continue
				}
				resetRegistryEnv(t)
				for pk, pv := range probe.pre {
					os.Setenv(pk, pv)
				}
				os.Setenv(k.Env, in)
				c, err := Load()
				if err != nil {
					t.Errorf("%s=%q: Load() rejected a registry-legal value: %v", k.Env, in, err)
					continue
				}
				if got := probe.get(c); got != wantOut {
					t.Errorf("%s=%q resolved to %q, want %q", k.Env, in, got, wantOut)
				}
			}

			resetRegistryEnv(t)
			for pk, pv := range probe.pre {
				os.Setenv(pk, pv)
			}
			os.Setenv(k.Env, "__registry_illegal__")
			if _, err := Load(); err == nil {
				t.Errorf("%s=__registry_illegal__: Load() accepted a value outside the registry enum", k.Env)
			}
		})
	}
}

// --- strict parse agreement for bool/int knobs ---

func TestRegistryBoolStrictness(t *testing.T) {
	probes := buildLoadProbes(t)
	for _, k := range Registry {
		if k.Type.Kind != TypeBool {
			continue
		}
		probe, ok := probes[k.Env]
		if !ok {
			continue // lenient world-package bools have a dedicated test
		}
		if probe.linuxOnly && runtime.GOOS != "linux" {
			continue
		}
		resetRegistryEnv(t)
		for pk, pv := range probe.pre {
			os.Setenv(pk, pv)
		}
		os.Setenv(k.Env, "2")
		if _, err := Load(); err == nil {
			t.Errorf("%s=2: Load() accepted an invalid boolean (strict 0|1 contract)", k.Env)
		}
	}
}

func TestRegistryIntStrictness(t *testing.T) {
	probes := buildLoadProbes(t)
	for _, k := range Registry {
		if k.Type.Kind != TypeInt {
			continue
		}
		probe, ok := probes[k.Env]
		if !ok {
			continue
		}
		if probe.linuxOnly && runtime.GOOS != "linux" {
			continue
		}
		resetRegistryEnv(t)
		for pk, pv := range probe.pre {
			os.Setenv(pk, pv)
		}
		os.Setenv(k.Env, "not-a-number")
		if _, err := Load(); err == nil {
			t.Errorf("%s=not-a-number: Load() accepted an invalid integer", k.Env)
		}
	}
}

// --- legacy tombstones ---

func TestRegistryLegacyTombstonesRejected(t *testing.T) {
	for _, k := range Registry {
		if k.Default.Kind != DefaultLegacy {
			continue
		}
		if k.Default.Value != "removed-error" {
			t.Errorf("%s: unexpected legacy sub-kind %q", k.Env, k.Default.Value)
			continue
		}
		resetRegistryEnv(t)
		os.Setenv(k.Env, "any-value")
		_, err := Load()
		if err == nil {
			t.Errorf("%s: registry says removed-error but Load() accepted it", k.Env)
			continue
		}
		if !strings.Contains(err.Error(), k.Env) {
			t.Errorf("%s: rejection error %q does not name the env", k.Env, err)
		}
	}
}

// --- runtime-emitted identity behavior (documenting, not changing) ---

func TestRegistryRuntimeEmittedLoadBehavior(t *testing.T) {
	// SessionID: generated when unset, adopted when set (resume/adoption).
	resetRegistryEnv(t)
	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if strings.TrimSpace(c.SessionID) == "" {
		t.Error("SessionID should be auto-generated when unset")
	}
	if strings.TrimSpace(c.TapeID) == "" {
		t.Error("TapeID should be auto-generated when unset")
	}
	if strings.TrimSpace(c.RunID) == "" {
		t.Error("RunID should be generated")
	}

	// RunID: stale inbound values are ignored, never adopted.
	resetRegistryEnv(t)
	os.Setenv(EnvRunID, "stale-run-id")
	c, err = Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.RunID == "stale-run-id" {
		t.Error("RunID must be regenerated, not adopted from env")
	}

	// ContextTape: declared for defensive stripping only; Load() ignores it.
	resetRegistryEnv(t)
	os.Setenv(EnvContextTape, "any-tape-ref")
	if _, err := Load(); err != nil {
		t.Errorf("Load() should ignore %s, got error: %v", EnvContextTape, err)
	}
}

// --- envs the binary reads outside Load() ---

func TestRegistryUnvalidatedEnvsLoadCleanly(t *testing.T) {
	// These two are consumed directly from the process env by internal/world;
	// Load() neither stores nor validates them (registry Notes document the
	// lenient parsing).
	resetRegistryEnv(t)
	os.Setenv(EnvPromptBudgetVisibility, "not-a-mode")
	os.Setenv(EnvWorldOnePerShell, "banana")
	if _, err := Load(); err != nil {
		t.Errorf("Load() should not validate world-package envs, got error: %v", err)
	}
}

// TestRegistryConfigDirPassthrough: QUINE_CONFIG_DIR is the one registry knob
// Load() never stores on Config — it is consumed at load time and reaches
// children only by passing through. Under the deleted synthesizer that required
// a hand-written baseEnv special case. It is now the default behavior of the
// pipeline (operator-only → pinned → inherited verbatim), which is why the claim
// is re-expressed against BuildChildEnv rather than a serializer.
func TestRegistryConfigDirPassthrough(t *testing.T) {
	resetRegistryEnv(t)
	os.Setenv(EnvConfigDir, "/tmp/quine-registry-config")
	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	assertCrossesEveryBoundary(t, c, EnvConfigDir+"=/tmp/quine-registry-config")

	// And it is not settable through the mediated channel: pinned means the
	// operator's call, not the agent's.
	if _, err := ParseEnvOverride([]byte(EnvConfigDir + "=/tmp/agent-chosen\n")); err == nil {
		t.Errorf("config/env/override must reject %s (operator-only → pinned)", EnvConfigDir)
	}
}
