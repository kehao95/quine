package config

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ErrDepthExceeded is returned when QUINE_DEPTH >= QUINE_MAX_DEPTH.
var ErrDepthExceeded = errors.New("max recursion depth exceeded")

// processStartTime captures the host process birth time for deriving readable
// runtime identities when QUINE_SESSION_ID is not supplied.
var processStartTime = time.Now()

// PromptMetaphor controls optional metaphor framing in the system prompt.
type PromptMetaphor string

const (
	PromptMetaphorOff           PromptMetaphor = "off"
	PromptMetaphorThermodynamic PromptMetaphor = "thermodynamic"
)

// PromptSelfModel controls how much self-model framing the system prompt reveals.
type PromptSelfModel string

const (
	PromptSelfModelBasic    PromptSelfModel = "basic"
	PromptSelfModelAdvanced PromptSelfModel = "advanced"
)

// PromptInstructionSurface controls how much runtime instruction text is
// rendered into the provider-visible system prompt.
type PromptInstructionSurface string

const (
	PromptInstructionSurfaceStandard        PromptInstructionSurface = "standard"
	PromptInstructionSurfaceMinimal         PromptInstructionSurface = "minimal"
	PromptInstructionSurfaceMinimalAutonomy PromptInstructionSurface = "minimal_autonomy"
	// MinimalExistence denies any task/objective while leaving the process able
	// to act; the activation clause is carried by the initial user message, not
	// the system prompt.
	PromptInstructionSurfaceMinimalExistence PromptInstructionSurface = "minimal_existence"
)

// PromptRuntimeSurface controls whether the prompt explicitly teaches Quine's
// runtime self-mapping and neighbor-discovery surface.
type PromptRuntimeSurface string

const (
	PromptRuntimeSurfaceVisible PromptRuntimeSurface = "visible"
	PromptRuntimeSurfaceHidden  PromptRuntimeSurface = "hidden"
)

// PromptPersona controls optional role-stance framing in the system prompt.
type PromptPersona string

const (
	PromptPersonaNone         PromptPersona = ""
	PromptPersonaCoder        PromptPersona = "coder"
	PromptPersonaAnalyst      PromptPersona = "analyst"
	PromptPersonaEngineer     PromptPersona = "engineer"
	PromptPersonaArchitect    PromptPersona = "architect"
	PromptPersonaExplorer     PromptPersona = "explorer"
	PromptPersonaSteward      PromptPersona = "steward"
	PromptPersonaCartographer PromptPersona = "cartographer"
	PromptPersonaGardener     PromptPersona = "gardener"
	PromptPersonaWitness      PromptPersona = "witness"
	PromptPersonaCatalyst     PromptPersona = "catalyst"
	PromptPersonaSkeptic      PromptPersona = "skeptic"
)

// WorkspaceRevisionMode controls which explicit revision primitives are
// available once workspace physics are enabled.
type WorkspaceRevisionMode string

const (
	WorkspaceRevisionNone    WorkspaceRevisionMode = "none"
	WorkspaceRevisionRestore WorkspaceRevisionMode = "restore"
)

// WorldKind identifies the current process world substrate.
type WorldKind string

const (
	WorldHost       WorldKind = "host"
	WorldSubjective WorldKind = "subjective"
)

// ProtectionMode identifies the current process protection surface.
type ProtectionMode string

const (
	ProtectionNone          ProtectionMode = "none"
	ProtectionTransactional ProtectionMode = "transactional"
)

// SelfReentryMode controls how quine defaults its own fork/exec re-entry target.
type SelfReentryMode string

const (
	SelfReentryModeSelf           SelfReentryMode = "self"
	SelfReentryModeExecutablePath SelfReentryMode = "executable_path"
)

type ToolGates struct {
	AnchorMemoryEnabled      bool   // QUINE_ANCHOR_MEMORY (default false)
	AnchorFoldEnabled        bool   // QUINE_ANCHOR_FOLD_ENABLED (default true)
	AnchorMarkEnabled        bool   // QUINE_ANCHOR_MARK_ENABLED (default true)
	IdleEnabled              bool   // QUINE_IDLE_ENABLED (default false)
	ExecEnabled              bool   // QUINE_EXEC_ENABLED (default true)
	SpawnEnabledFlag         bool   // QUINE_SPAWN_ENABLED (default false)
	AgentsMDEnabled          bool   // QUINE_AGENTS_MD_ENABLED (default false)
	AgentsSkillsEnabled      bool   // QUINE_AGENTS_SKILLS_ENABLED (default false)
	VisionEnabled            bool   // QUINE_VISION_ENABLED (default true)
	ShTimeoutOverrideEnabled bool   // QUINE_SH_TIMEOUT_OVERRIDE_ENABLED (default true)
	ShStdinEnabled           bool   // QUINE_SH_STDIN_ENABLED (default true)
	ShDetachEnabled          bool   // QUINE_SH_DETACH_ENABLED (default true)
	ExitEnabled              bool   // QUINE_EXIT_ENABLED (default true)
	ForkEnabled              bool   // QUINE_FORK_ENABLED (default true)
	ShInteractiveEnabled     bool   // QUINE_SH_INTERACTIVE_ENABLED (default true)
	FSMutationTelemetry      bool   // QUINE_FS_MUTATION_TELEMETRY_ENABLED (default true)
	ForkWorldEnabled         bool   // QUINE_FORK_WORLD_ENABLED (default false)
	EphemeralBody            bool   // QUINE_EPHEMERAL_BODY_ENABLED (default false)
	SuppressInitialBegin     bool   // QUINE_SUPPRESS_INITIAL_BEGIN (default false)
	InitialUserMessage       string // QUINE_INITIAL_USER_MESSAGE: overrides the synthetic TTY "Begin." user message (default "")
	SelfSourceCodeEnabled    bool   // QUINE_SELF_SOURCE_CODE_ENABLED (default false)
	PeerDiscoveryEnabled     bool   // QUINE_PEER_DISCOVERY_ENABLED (default false)
	EmptyAssistantSuccess    bool   // QUINE_EMPTY_ASSISTANT_SUCCESS (default false)
	ReadyTextAutoIdle        bool   // QUINE_READY_TEXT_AUTO_IDLE (default false)

	defaultsApplied bool
}

func DefaultToolGates() ToolGates {
	return ToolGates{
		ExecEnabled:              true,
		VisionEnabled:            true,
		ShTimeoutOverrideEnabled: true,
		ShStdinEnabled:           true,
		ShDetachEnabled:          true,
		ExitEnabled:              true,
		ForkEnabled:              true,
		AnchorFoldEnabled:        true,
		AnchorMarkEnabled:        true,
		ShInteractiveEnabled:     true,
		FSMutationTelemetry:      true,
		defaultsApplied:          true,
	}
}

type Identity struct {
	ModelID       string
	SessionID     string
	RunID         string
	IncarnationID int
	TapeID        string
	ParentSession string
	Depth         int
}

type Transport struct {
	APIKey              string
	APIBase             string
	Provider            string
	UserAgent           string
	ThinkingBudget      string
	ServiceTier         string
	DebugRequestBodyDir string
}

type Limits struct {
	MaxDepth                  int
	MaxConcurrent             int
	MaxAgents                 int
	ShTimeout                 int
	OutputTruncate            int
	MaxTurns                  int
	WallClockExitSeconds      int
	ContextWindow             int
	MemoryWarnTokens          int
	MemoryDangerTokens        int
	MemoryDeathTokens         int
	ForkDefaultTimeoutSeconds int
	PeerDiscoveryHeartbeatMS  int
}

type PromptConfig struct {
	PromptMetaphor           PromptMetaphor
	PromptSelfModel          PromptSelfModel
	PromptInstructionSurface PromptInstructionSurface
	PromptRuntimeSurface     PromptRuntimeSurface
	PromptPersona            PromptPersona
	FailOnImpossible         bool
	NoMissionAutonomy        string
	AgentsMDPath             string
	Skills                   []Skill
	PromptImplDetails        bool
	PromptCtlPhysics         bool
	MemoryStrategyHints      bool
}

type WorkspaceConfig struct {
	WorkspaceEnabled         bool
	WorkspaceRoot            string
	Workspace                string
	WorkspaceBackend         string
	WorkspaceOverlayDriver   string
	WorkspaceRevisionMode    WorkspaceRevisionMode
	WorkspaceCurrentRevision string
	WorkspaceSession         string
	WorkspaceOwner           bool
	WorkspaceCommitOnSignal  bool
	WorkspaceBootstrap       string
}

type Paths struct {
	DataDir           string
	RetentionDir      string
	WorkDir           string
	Shell             string
	ShNetwork         string
	ExecutablePath    string
	SelfReentryMode   SelfReentryMode
	SelfReentryTarget string
}

// Config holds all runtime configuration for Quine.
// Every field is populated from environment variables by Load().
type Config struct {
	Identity
	Transport
	Limits
	ToolGates
	PromptConfig
	WorkspaceConfig
	Paths
}

// APIModelID returns the model ID to use in API calls.
func (c *Config) APIModelID() string {
	return c.ModelID
}

// CurrentWorld reports the current process world substrate.
func (c *Config) CurrentWorld() WorldKind {
	if c != nil && c.WorkspaceTransactional() {
		return WorldSubjective
	}
	return WorldHost
}

// CurrentProtection reports the current process protection surface.
func (c *Config) CurrentProtection() ProtectionMode {
	if c != nil && c.WorkspaceTransactional() {
		return ProtectionTransactional
	}
	return ProtectionNone
}

// EffectiveWorkspaceBackend reports the active workspace backend, defaulting to
// overlay for backward compatibility when the field is unset on a constructed Config.
func (c *Config) EffectiveWorkspaceBackend() string {
	if c == nil || !c.WorkspaceEnabled {
		return ""
	}
	if strings.TrimSpace(c.WorkspaceBackend) == "" {
		return "overlay"
	}
	return c.WorkspaceBackend
}

// WorkspaceTransactional reports whether the workspace is a transactional,
// process-local managed world rather than a direct shared host surface.
func (c *Config) WorkspaceTransactional() bool {
	return c != nil && c.WorkspaceEnabled && c.EffectiveWorkspaceBackend() == "overlay"
}

// ThermodynamicMetaphorEnabled reports whether thermodynamic framing should be
// added to the system prompt.
func (c *Config) ThermodynamicMetaphorEnabled() bool {
	return c.PromptMetaphor == PromptMetaphorThermodynamic
}

// SelfModelMode reports the configured prompt self-model disclosure mode.
func (c *Config) SelfModelMode() PromptSelfModel {
	if c == nil || c.PromptSelfModel == "" {
		return PromptSelfModelAdvanced
	}
	return c.PromptSelfModel
}

// AdvancedSelfModelEnabled reports whether advanced continuity/cognition
// framing should appear in the system prompt.
func (c *Config) AdvancedSelfModelEnabled() bool {
	return c.SelfModelMode() == PromptSelfModelAdvanced
}

// InstructionSurfaceMode reports the configured instruction-surface mode.
func (c *Config) InstructionSurfaceMode() PromptInstructionSurface {
	if c == nil || c.PromptInstructionSurface == "" {
		return PromptInstructionSurfaceStandard
	}
	return c.PromptInstructionSurface
}

// MinimalInstructionSurface reports whether the runtime prompt should collapse
// to the minimal mission/missionless statement and suppress runtime teaching.
func (c *Config) MinimalInstructionSurface() bool {
	switch c.InstructionSurfaceMode() {
	case PromptInstructionSurfaceMinimal, PromptInstructionSurfaceMinimalAutonomy, PromptInstructionSurfaceMinimalExistence:
		return true
	default:
		return false
	}
}

// MinimalAutonomyInstructionSurface reports whether the collapsed missionless
// prompt should include the smallest activation clause.
func (c *Config) MinimalAutonomyInstructionSurface() bool {
	return c.InstructionSurfaceMode() == PromptInstructionSurfaceMinimalAutonomy
}

// NoMissionAutonomyLevel reports the normalized missionless-autonomy level
// (off|autonomy|sensing|full); the empty zero value normalizes to "off".
func (c *Config) NoMissionAutonomyLevel() string {
	if c == nil || c.NoMissionAutonomy == "" {
		return "off"
	}
	return c.NoMissionAutonomy
}

// RuntimeSurfaceMode reports the configured prompt runtime-surface disclosure mode.
func (c *Config) RuntimeSurfaceMode() PromptRuntimeSurface {
	if c == nil || c.PromptRuntimeSurface == "" {
		return PromptRuntimeSurfaceVisible
	}
	return c.PromptRuntimeSurface
}

// PersonaMode reports the configured prompt persona role stance.
func (c *Config) PersonaMode() PromptPersona {
	if c == nil {
		return PromptPersonaNone
	}
	return c.PromptPersona
}

// RuntimeSurfaceVisible reports whether the system prompt should explicitly
// teach Quine's runtime self-mapping surface.
func (c *Config) RuntimeSurfaceVisible() bool {
	return c.RuntimeSurfaceMode() == PromptRuntimeSurfaceVisible
}

// ForkEnabled reports whether the fork tool/runtime capability is enabled.
func (c *Config) ForkEnabled() bool {
	return c == nil || c.effectiveToolGates().ForkEnabled
}

// AnchorFoldEnabled reports whether the fold move on the mark tool is exposed.
// Disabling it removes in-process self-consolidation while leaving mark/unfold
// and the raw context/state file substrate intact.
func (c *Config) AnchorFoldEnabled() bool {
	return c == nil || c.effectiveToolGates().AnchorFoldEnabled
}

// AnchorMarkEnabled reports whether the mark crystallization tool is exposed.
// Disabling it (with anchor memory still on) removes the agent's in-process
// self-compression move while leaving unfold and the raw context/state file
// substrate intact, so the only remaining route to lower context load is the
// filesystem itself (directly, or via a recruited peer process).
func (c *Config) AnchorMarkEnabled() bool {
	return c == nil || c.effectiveToolGates().AnchorMarkEnabled
}

// SpawnEnabled reports whether the spawn tool/runtime capability is enabled.
func (c *Config) SpawnEnabled() bool {
	return c != nil && c.SpawnEnabledFlag
}

// ExitEnabled reports whether the exit tool is available.
func (c *Config) ExitEnabled() bool {
	return c == nil || c.effectiveToolGates().ExitEnabled
}

// ShInteractiveEnabled reports whether sh(interactive=true) is exposed.
func (c *Config) ShInteractiveEnabled() bool {
	return c == nil || c.effectiveToolGates().ShInteractiveEnabled
}

// ShStdinEnabled reports whether sh(stdin=...) is exposed and accepted.
func (c *Config) ShStdinEnabled() bool {
	return c == nil || c.effectiveToolGates().ShStdinEnabled
}

// ShDetachEnabled reports whether sh(detach=true) is exposed and accepted.
func (c *Config) ShDetachEnabled() bool {
	return c == nil || c.effectiveToolGates().ShDetachEnabled
}

// FSMutationTelemetryEnabled reports whether fs_mutations telemetry is exposed.
func (c *Config) FSMutationTelemetryEnabled() bool {
	return c == nil || c.effectiveToolGates().FSMutationTelemetry
}

// ShTimeoutOverrideEnabled reports whether sh(timeout=...) is exposed and accepted.
func (c *Config) ShTimeoutOverrideEnabled() bool {
	return c == nil || c.effectiveToolGates().ShTimeoutOverrideEnabled
}

func (c *Config) effectiveToolGates() ToolGates {
	if c == nil || !c.ToolGates.defaultsApplied {
		return DefaultToolGates()
	}
	return c.ToolGates
}

const linuxSelfReentryTarget = "/proc/self/exe"

func resolveSelfReentryMode(raw string) (SelfReentryMode, error) {
	mode := strings.TrimSpace(raw)
	if mode == "" {
		return SelfReentryModeSelf, nil
	}
	switch SelfReentryMode(mode) {
	case SelfReentryModeSelf, SelfReentryModeExecutablePath:
		return SelfReentryMode(mode), nil
	default:
		return "", fmt.Errorf("QUINE_SELF_REENTRY_MODE=%q: must be %q or %q", mode, SelfReentryModeSelf, SelfReentryModeExecutablePath)
	}
}

func resolveSelfReentryTarget(goos string, mode SelfReentryMode, executablePath string) (string, error) {
	switch mode {
	case SelfReentryModeSelf:
		if goos != "linux" {
			return "", fmt.Errorf(
				"QUINE_SELF_REENTRY_MODE=%q is only supported on Linux; use QUINE_SELF_REENTRY_MODE=%q instead",
				mode,
				SelfReentryModeExecutablePath,
			)
		}
		return linuxSelfReentryTarget, nil
	case SelfReentryModeExecutablePath:
		executablePath = strings.TrimSpace(executablePath)
		if executablePath == "" {
			return "", fmt.Errorf("could not determine current executable path for QUINE_SELF_REENTRY_MODE=%q", mode)
		}
		return executablePath, nil
	default:
		return "", fmt.Errorf("unsupported self re-entry mode %q", mode)
	}
}

// IdleToolEnabled reports whether the explicit suspension primitive is available.
func (c *Config) IdleToolEnabled() bool {
	return c != nil && c.IdleEnabled
}

// RuntimeRoot returns the durable runtime-state root shared by the process tree.
// It holds live session surfaces, pid routing, job directories, coordination
// locks, and workspace overlay state. Retained session state falls back under
// this root when QUINE_RETENTION_DIR is unset.
func (c *Config) RuntimeRoot() string {
	return c.DataDir
}

// LockDir returns the internal coordination directory used for shared
// semaphore and agent-registry files under the runtime root.
func (c *Config) LockDir() string {
	return filepath.Join(c.RuntimeRoot(), "locks")
}

// PIDDir returns the public live-process routing surface keyed by pid.
// Each pid entry points at the corresponding agent/<session>/public projection.
func (c *Config) PIDDir() string {
	return filepath.Join(c.RuntimeRoot(), "pid")
}

// SessionRetainedDir returns the canonical retained root for a session.
// When QUINE_RETENTION_DIR is set, retained session state lives under
// QUINE_RETENTION_DIR/sessions/<session>. Otherwise it falls back to
// QUINE_DATA_DIR/log/<session>.
func (c *Config) SessionRetainedDir(sessionID string) string {
	if sessionID == "" {
		sessionID = c.SessionID
	}
	if strings.TrimSpace(c.RetentionDir) != "" {
		return filepath.Join(c.RetentionDir, "sessions", sessionID)
	}
	return filepath.Join(c.RuntimeRoot(), "log", sessionID)
}

// SessionLogDir is a compatibility alias for the canonical retained root.
// It survives after the live agent root is removed.
func (c *Config) SessionLogDir(sessionID string) string {
	return c.SessionRetainedDir(sessionID)
}

// SessionLogPath returns the retained operational log path for a session.
func (c *Config) SessionLogPath(sessionID string) string {
	return filepath.Join(c.SessionRetainedDir(sessionID), "runtime.log")
}

// TapeDir returns the append-only tape directory for a session.
func (c *Config) TapeDir(sessionID string) string {
	return filepath.Join(c.SessionRetainedDir(sessionID), "tapes")
}

// TapePath returns the append-only JSONL path for an internal trace artifact
// under the current stable session.
func (c *Config) TapePath(tapeID string) string {
	if tapeID == "" {
		tapeID = c.TapeID
	}
	return filepath.Join(c.TapeDir(""), tapeID+".jsonl")
}

// TapeYAMLPath returns the YAML mirror path for an internal trace artifact
// under the current stable session.
func (c *Config) TapeYAMLPath(tapeID string) string {
	if tapeID == "" {
		tapeID = c.TapeID
	}
	return filepath.Join(c.TapeDir(""), tapeID+".log.yaml")
}

// SessionRoot returns the canonical session root for a session.
func (c *Config) SessionRoot(sessionID string) string {
	if sessionID == "" {
		sessionID = c.SessionID
	}
	return filepath.Join(c.RuntimeRoot(), "agent", sessionID)
}

// SessionIncarnationDir returns the incarnation tree root for a session.
func (c *Config) SessionIncarnationDir(sessionID string) string {
	return filepath.Join(c.SessionRetainedDir(sessionID), "inc")
}

// SessionIncarnationPath returns the numeric incarnation root under a session.
func (c *Config) SessionIncarnationPath(sessionID string, incarnationID int) string {
	return filepath.Join(c.SessionIncarnationDir(sessionID), strconv.Itoa(incarnationID))
}

// SessionCurrentIncarnationPath returns the current-incarnation alias path.
func (c *Config) SessionCurrentIncarnationPath(sessionID string) string {
	return filepath.Join(c.SessionIncarnationDir(sessionID), "current")
}

// SessionMissionPath returns the canonical mission surface for a session.
func (c *Config) SessionMissionPath(sessionID string) string {
	return filepath.Join(c.SessionRoot(sessionID), "mission.txt")
}

// SessionStatusDir returns the canonical status surface for a session.
func (c *Config) SessionStatusDir(sessionID string) string {
	return filepath.Join(c.SessionRoot(sessionID), "status")
}

// SessionInboxPath returns the canonical inbox snapshot for a session.
func (c *Config) SessionInboxPath(sessionID string) string {
	return filepath.Join(c.SessionStatusDir(sessionID), "inbox.json")
}

// SessionWorldDir returns the canonical projected-world surface for a session.
func (c *Config) SessionWorldDir(sessionID string) string {
	return filepath.Join(c.SessionRoot(sessionID), "world")
}

// SessionControlLogPath returns the retained append-only control-event audit log.
func (c *Config) SessionControlLogPath(sessionID string) string {
	return filepath.Join(c.SessionRetainedDir(sessionID), "control.jsonl")
}

// AgentRoot returns the live process-filesystem root for this session.
// The root remains stable across exec because SessionID remains stable.
func (c *Config) AgentRoot() string {
	return c.SessionRoot(c.SessionID)
}

// AgentStatusDir returns the public status surface under the session root.
func (c *Config) AgentStatusDir() string {
	return filepath.Join(c.AgentRoot(), "status")
}

// AgentLogDir returns the public audit/log surface under the session root.
func (c *Config) AgentLogDir() string {
	return filepath.Join(c.AgentRoot(), "log")
}

// ControlDir returns the peer-write control directory for the session. Under
// the FUSE-only control surface this path is materialized as virtual FUSE nodes
// (public/ctl/<action>), not real files at the agent root; the accessor is
// retained as the canonical "where ctl lives" reference.
func (c *Config) ControlDir() string {
	return filepath.Join(c.AgentRoot(), "ctl")
}

// ControlPath returns the alias for the peer-write control directory.
func (c *Config) ControlPath() string {
	return c.ControlDir()
}

// InboxPath returns the current pending inbox snapshot for the session.
func (c *Config) InboxPath() string {
	return filepath.Join(c.AgentStatusDir(), "inbox.json")
}

// ControlLogPath returns the append-only control-event audit log.
func (c *Config) ControlLogPath() string {
	return filepath.Join(c.AgentLogDir(), "control.jsonl")
}

// JobSessionDir returns the managed job directory root for a session.
func (c *Config) JobSessionDir(sessionID string) string {
	if sessionID == "" {
		sessionID = c.SessionID
	}
	return filepath.Join(c.RuntimeRoot(), "jobs", sessionID)
}

// WorkspaceStateDir returns the overlay-state root for a workspace session.
func (c *Config) WorkspaceStateDir(workspaceSession string) string {
	if workspaceSession == "" {
		workspaceSession = c.WorkspaceSession
	}
	return filepath.Join(c.RuntimeRoot(), "workspaces", workspaceSession)
}

// Load reads all configuration from environment variables and returns
// a validated Config. It returns an error if required variables are
// missing or if depth is exceeded.
//
// Four variables are required:
//   - QUINE_MODEL_ID:   Model name (e.g. "claude-sonnet-4-5-20250929", "gpt-4o", "kimi-k2.5")
//   - QUINE_API_TYPE:   Wire protocol: "openai", "anthropic", or "openai-responses"
//   - QUINE_API_BASE:   API base URL (e.g. "https://api.anthropic.com", "https://api.openai.com")
//   - QUINE_API_KEY:    API key or OAuth sentinel such as "claude-oauth", "codex-oauth", "kimi-oauth", or "copilot-oauth"
func Load() (*Config, error) {
	c := &Config{ToolGates: DefaultToolGates()}

	for _, load := range []func(*Config) error{
		loadRequiredIdentityAndTransport,
		loadWorkspaceConfig,
		loadLimitConfig,
		loadIdentityAndPathConfig,
		loadPromptAndTransportOptions,
		loadToolGates,
		loadSelfReentryConfig,
	} {
		if err := load(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func loadRequiredIdentityAndTransport(c *Config) error {
	// --- 4 required fields ---
	c.ModelID = os.Getenv(EnvModelID)
	if c.ModelID == "" {
		return fmt.Errorf("QUINE_MODEL_ID is required")
	}

	c.Provider = os.Getenv(EnvAPIType)
	if c.Provider == "" {
		return fmt.Errorf("QUINE_API_TYPE is required (\"openai\", \"anthropic\", or \"openai-responses\")")
	}
	if c.Provider != "openai" && c.Provider != "anthropic" && c.Provider != "openai-responses" {
		return fmt.Errorf("unsupported QUINE_API_TYPE=%q: must be \"openai\", \"anthropic\", or \"openai-responses\"", c.Provider)
	}

	c.APIBase = os.Getenv(EnvAPIBase)
	if c.APIBase == "" {
		return fmt.Errorf("QUINE_API_BASE is required")
	}

	c.APIKey = os.Getenv(EnvAPIKey)
	if c.APIKey == "" {
		return fmt.Errorf("QUINE_API_KEY is required")
	}
	c.DebugRequestBodyDir = strings.TrimSpace(os.Getenv(EnvDebugRequestBodyDir))
	return nil
}

func loadWorkspaceConfig(c *Config) error {
	var err error
	// --- Optional string fields ---
	c.ParentSession = os.Getenv(EnvParentSession)
	c.ForkWorldEnabled, err = envBoolDefault(EnvForkWorldEnabled, false)
	if err != nil {
		return err
	}

	// --- Workspace physics ---
	c.WorkspaceRoot = strings.TrimSpace(os.Getenv(EnvWorkspaceRoot))
	c.Workspace = strings.TrimSpace(os.Getenv(EnvWorkspace))
	c.WorkspaceBackend = strings.TrimSpace(os.Getenv(EnvWorkspaceBackend))
	c.WorkspaceOverlayDriver = strings.TrimSpace(os.Getenv(EnvWorkspaceOverlayDriver))
	revisionMode := strings.TrimSpace(os.Getenv(EnvWorkspaceRevisionMode))
	c.WorkspaceCurrentRevision = strings.TrimSpace(os.Getenv(EnvWorkspaceCurrentRevision))
	c.WorkspaceSession = strings.TrimSpace(os.Getenv(EnvWorkspaceSession))
	c.WorkspaceOwner, err = envBoolDefault(EnvWorkspaceOwner, true)
	if err != nil {
		return err
	}
	c.WorkspaceCommitOnSignal, err = envBoolDefault(EnvWorkspaceCommitOnSignal, false)
	if err != nil {
		return err
	}
	c.WorkspaceBootstrap = strings.TrimSpace(os.Getenv(EnvWorkspaceBootstrap))
	if strings.TrimSpace(os.Getenv(EnvWorkspaceSource)) != "" {
		return fmt.Errorf("QUINE_WORKSPACE_SOURCE has been removed; use QUINE_WORKSPACE_ROOT and QUINE_WORKSPACE")
	}
	if c.WorkspaceRoot == "" && c.Workspace != "" {
		c.WorkspaceRoot = c.Workspace
	}
	if c.Workspace == "" && c.WorkspaceRoot != "" {
		// If started inside the workspace root, default workspace scope to pwd.
		if wd, err := os.Getwd(); err == nil {
			if rootCandidate, rootErr := canonicalPath(c.WorkspaceRoot); rootErr == nil {
				if wdCandidate, wdErr := canonicalPath(wd); wdErr == nil && isPathWithin(rootCandidate, wdCandidate) {
					c.Workspace = wdCandidate
				}
			}
		}
	}
	if c.Workspace == "" && c.WorkspaceRoot != "" {
		c.Workspace = c.WorkspaceRoot
	}
	c.WorkspaceEnabled = c.WorkspaceRoot != "" || c.Workspace != "" || c.WorkspaceSession != ""
	if c.WorkspaceEnabled {
		if c.WorkspaceBackend == "" {
			c.WorkspaceBackend = "overlay"
		}
		switch c.WorkspaceBackend {
		case "overlay":
			if runtime.GOOS != "linux" {
				return fmt.Errorf("workspace backend %q is only supported on Linux", c.WorkspaceBackend)
			}
			if c.WorkspaceOverlayDriver == "" {
				c.WorkspaceOverlayDriver = "kernel"
			}
			switch c.WorkspaceOverlayDriver {
			case "kernel", "fuse":
			default:
				return fmt.Errorf("QUINE_WORKSPACE_OVERLAY_DRIVER=%q: must be %q or %q", c.WorkspaceOverlayDriver, "kernel", "fuse")
			}
			if revisionMode == "" {
				c.WorkspaceRevisionMode = WorkspaceRevisionRestore
			} else {
				switch WorkspaceRevisionMode(revisionMode) {
				case WorkspaceRevisionNone, WorkspaceRevisionRestore:
					c.WorkspaceRevisionMode = WorkspaceRevisionMode(revisionMode)
				default:
					return fmt.Errorf("QUINE_WORKSPACE_REVISION_MODE=%q: must be %q or %q",
						revisionMode, WorkspaceRevisionNone, WorkspaceRevisionRestore)
				}
			}
		case "direct":
			if c.WorkspaceOverlayDriver != "" {
				return fmt.Errorf("QUINE_WORKSPACE_OVERLAY_DRIVER requires QUINE_WORKSPACE_BACKEND=%q", "overlay")
			}
			if revisionMode == "" {
				c.WorkspaceRevisionMode = WorkspaceRevisionNone
			} else {
				switch WorkspaceRevisionMode(revisionMode) {
				case WorkspaceRevisionNone:
					c.WorkspaceRevisionMode = WorkspaceRevisionNone
				default:
					return fmt.Errorf("QUINE_WORKSPACE_REVISION_MODE=%q: backend %q only supports %q",
						revisionMode, c.WorkspaceBackend, WorkspaceRevisionNone)
				}
			}
		default:
			return fmt.Errorf("QUINE_WORKSPACE_BACKEND=%q: must be %q or %q", c.WorkspaceBackend, "overlay", "direct")
		}
		if c.WorkspaceRoot == "" {
			return fmt.Errorf("QUINE_WORKSPACE_ROOT or QUINE_WORKSPACE is required when workspace physics are enabled")
		}
		if c.Workspace == "" {
			c.Workspace = c.WorkspaceRoot
		}
		root, err := canonicalPath(c.WorkspaceRoot)
		if err != nil {
			return fmt.Errorf("canonicalize workspace root: %w", err)
		}
		workspace, err := canonicalPath(c.Workspace)
		if err != nil {
			return fmt.Errorf("canonicalize workspace: %w", err)
		}
		if !isPathWithin(root, workspace) {
			return fmt.Errorf("workspace %q must be within workspace root %q", workspace, root)
		}
		c.WorkspaceRoot = root
		c.Workspace = workspace
	} else {
		if revisionMode == "" {
			c.WorkspaceRevisionMode = WorkspaceRevisionNone
		} else {
			switch WorkspaceRevisionMode(revisionMode) {
			case WorkspaceRevisionNone:
				c.WorkspaceRevisionMode = WorkspaceRevisionNone
			default:
				return fmt.Errorf("QUINE_WORKSPACE_REVISION_MODE=%q requires workspace physics to be enabled", revisionMode)
			}
		}
	}
	if c.ForkWorldEnabled && !c.WorkspaceEnabled {
		return fmt.Errorf("QUINE_FORK_WORLD_ENABLED=1 requires explicit workspace physics; set QUINE_WORKSPACE_ROOT or QUINE_WORKSPACE with QUINE_WORKSPACE_BACKEND=direct or overlay")
	}
	return nil
}

func loadLimitConfig(c *Config) error {
	var err error
	// --- Integer fields with defaults ---

	c.ContextWindow, err = envInt(EnvContextWindow, 128_000)
	if err != nil {
		return err
	}

	defaultWarnTokens := defaultMemoryWarnTokens(c.ContextWindow)
	c.MemoryWarnTokens, err = envInt(EnvMemoryWarnTokens, defaultWarnTokens)
	if err != nil {
		return err
	}
	defaultDangerTokens := defaultMemoryDangerTokens(c.ContextWindow)
	c.MemoryDangerTokens, err = envInt(EnvMemoryDangerTokens, defaultDangerTokens)
	if err != nil {
		return err
	}
	c.MemoryDeathTokens, err = envInt(EnvMemoryDeathTokens, 0)
	if err != nil {
		return err
	}
	if c.MemoryWarnTokens <= 0 {
		return fmt.Errorf("QUINE_MEMORY_WARN_TOKENS must be > 0")
	}
	if c.MemoryDangerTokens <= c.MemoryWarnTokens {
		return fmt.Errorf("QUINE_MEMORY_DANGER_TOKENS=%d must be greater than QUINE_MEMORY_WARN_TOKENS=%d",
			c.MemoryDangerTokens, c.MemoryWarnTokens)
	}
	if c.MemoryDeathTokens < 0 {
		return fmt.Errorf("QUINE_MEMORY_DEATH_TOKENS must be >= 0")
	}
	if c.MemoryDeathTokens > 0 && c.MemoryDeathTokens <= c.MemoryDangerTokens {
		return fmt.Errorf("QUINE_MEMORY_DEATH_TOKENS=%d must be greater than QUINE_MEMORY_DANGER_TOKENS=%d",
			c.MemoryDeathTokens, c.MemoryDangerTokens)
	}

	c.MaxDepth, err = envInt(EnvMaxDepth, 0)
	if err != nil {
		return err
	}

	c.Depth, err = envInt(EnvDepth, 0)
	if err != nil {
		return err
	}

	c.MaxConcurrent, err = envInt(EnvMaxConcurrent, 0)
	if err != nil {
		return err
	}

	c.MaxAgents, err = envInt(EnvMaxAgents, 0)
	if err != nil {
		return err
	}

	c.ForkDefaultTimeoutSeconds, err = envInt(EnvForkDefaultTimeout, 0)
	if err != nil {
		return err
	}
	if c.ForkDefaultTimeoutSeconds < 0 {
		return fmt.Errorf("QUINE_FORK_DEFAULT_TIMEOUT_SECONDS must be >= 0")
	}

	c.ShTimeout, err = envInt(EnvShDefaultTimeout, 300)
	if err != nil {
		return err
	}
	shTimeoutOverrideEnabled, err := envBoolDefault(EnvShTimeoutOverride, true)
	if err != nil {
		return err
	}
	c.ToolGates.ShTimeoutOverrideEnabled = shTimeoutOverrideEnabled
	c.ToolGates.ShStdinEnabled, err = envBoolDefault(EnvShStdinEnabled, true)
	if err != nil {
		return err
	}
	c.ToolGates.ShDetachEnabled, err = envBoolDefault(EnvShDetachEnabled, true)
	if err != nil {
		return err
	}

	c.OutputTruncate, err = envInt(EnvOutputTruncate, 20480)
	if err != nil {
		return err
	}

	c.MaxTurns, err = envInt(EnvMaxTurns, 0)
	if err != nil {
		return err
	}
	c.WallClockExitSeconds, err = envInt(EnvWallClockExitSeconds, 0)
	if err != nil {
		return err
	}
	if c.WallClockExitSeconds < 0 {
		return fmt.Errorf("QUINE_WALL_CLOCK_EXIT_SECONDS must be >= 0")
	}
	c.PeerDiscoveryHeartbeatMS, err = envInt(EnvPeerDiscoveryHeartbeat, 5000)
	if err != nil {
		return err
	}
	if c.PeerDiscoveryHeartbeatMS <= 0 {
		return fmt.Errorf("QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS must be positive")
	}

	// --- Depth check (disabled when MaxDepth <= 0) ---
	if c.MaxDepth > 0 && c.Depth >= c.MaxDepth {
		return ErrDepthExceeded
	}
	return nil
}

func loadIdentityAndPathConfig(c *Config) error {
	var err error
	// --- Session/run identity ---
	// QUINE_SESSION_ID is a stable lineage identity. Passing it at bootstrap is
	// the explicit resume/adoption mechanism for an existing session.
	if inheritedID := strings.TrimSpace(os.Getenv(EnvSessionID)); inheritedID != "" {
		c.SessionID = inheritedID
	} else {
		c.SessionID = stableSessionID()
	}
	// QUINE_RUN_ID is a physical run fact exported to tools. It is regenerated
	// on every process activation, even if the environment contains a stale one.
	c.RunID = processRunID()

	// --- Runtime/retention dirs ---
	rawDataDir, hasDataDir := os.LookupEnv(EnvDataDir)
	c.RetentionDir = strings.TrimSpace(os.Getenv(EnvRetentionDir))
	if hasDataDir {
		c.DataDir = strings.TrimSpace(rawDataDir)
		if c.DataDir == "" {
			c.DataDir = ".quine/"
		}
	} else {
		c.DataDir = ".quine/"
	}
	c.TapeID = strings.TrimSpace(os.Getenv(EnvTapeID))
	if c.TapeID == "" {
		c.TapeID, err = nextTapeID(c.TapeDir(""))
		if err != nil {
			return fmt.Errorf("generating tape ID: %w", err)
		}
	}
	if c.WorkspaceEnabled {
		if c.WorkspaceSession == "" {
			c.WorkspaceSession = c.SessionID
		}
		dataDirReal, err := canonicalPath(c.DataDir)
		if err != nil {
			return fmt.Errorf("canonicalize data dir: %w", err)
		}
		if isPathWithin(c.WorkspaceRoot, dataDirReal) {
			return fmt.Errorf("QUINE_DATA_DIR %q must be outside workspace root %q", dataDirReal, c.WorkspaceRoot)
		}
		if c.RetentionDir != "" {
			retentionDirReal, err := canonicalPath(c.RetentionDir)
			if err != nil {
				return fmt.Errorf("canonicalize retention dir: %w", err)
			}
			if isPathWithin(c.WorkspaceRoot, retentionDirReal) {
				return fmt.Errorf("QUINE_RETENTION_DIR %q must be outside workspace root %q", retentionDirReal, c.WorkspaceRoot)
			}
		}
		if _, err := os.Stat(c.WorkspaceRoot); err != nil {
			return fmt.Errorf("workspace root %q must exist: %w", c.WorkspaceRoot, err)
		}
	}

	// --- Shell ---
	c.Shell = os.Getenv(EnvShell)
	if c.Shell == "" {
		c.Shell = "/bin/sh"
	}
	c.ShNetwork = strings.TrimSpace(os.Getenv(EnvShNetwork))
	if c.ShNetwork == "" {
		c.ShNetwork = "host"
	}
	switch c.ShNetwork {
	case "host":
	case "none":
		// Network isolation is enforced via a network namespace (CLONE_NEWNET),
		// which only exists on Linux. On other platforms the isolation flag is
		// silently discarded and the job would run with full host network — a
		// validated security boundary failing open. Fail closed instead, matching
		// the overlay/fuse/self-reentry Linux gates.
		if runtime.GOOS != "linux" {
			return fmt.Errorf("QUINE_SH_NETWORK=%q is only supported on Linux", c.ShNetwork)
		}
	default:
		return fmt.Errorf("QUINE_SH_NETWORK=%q: must be \"host\" or \"none\"", c.ShNetwork)
	}
	if executablePath, err := os.Executable(); err == nil {
		c.ExecutablePath = executablePath
	}

	// --- Work dir ---
	c.WorkDir = os.Getenv(EnvWorkDir)
	if c.WorkDir == "" {
		if c.WorkspaceEnabled {
			// Keep command wrapper cwd aligned with visible workspace by default.
			c.WorkDir = c.Workspace
		}
	}
	if c.WorkDir == "" {
		// Default to current working directory at startup
		if wd, err := os.Getwd(); err == nil {
			c.WorkDir = wd
		}
	}
	return nil
}

func loadPromptAndTransportOptions(c *Config) error {
	var err error
	// --- Prompt metaphor mode ---
	metaphor := os.Getenv(EnvPromptMetaphor)
	if metaphor == "" {
		c.PromptMetaphor = PromptMetaphorOff
	} else {
		switch PromptMetaphor(metaphor) {
		case PromptMetaphorOff, PromptMetaphorThermodynamic:
			c.PromptMetaphor = PromptMetaphor(metaphor)
		default:
			return fmt.Errorf("QUINE_PROMPT_METAPHOR=%q: must be %q or %q",
				metaphor, PromptMetaphorOff, PromptMetaphorThermodynamic)
		}
	}

	selfModel := os.Getenv(EnvPromptSelfModel)
	if selfModel == "" {
		c.PromptSelfModel = PromptSelfModelAdvanced
	} else {
		switch PromptSelfModel(selfModel) {
		case PromptSelfModelBasic, PromptSelfModelAdvanced:
			c.PromptSelfModel = PromptSelfModel(selfModel)
		default:
			return fmt.Errorf("QUINE_PROMPT_SELF_MODEL=%q: must be %q or %q",
				selfModel, PromptSelfModelBasic, PromptSelfModelAdvanced)
		}
	}

	instructionSurface := os.Getenv(EnvPromptInstructionSurface)
	if instructionSurface == "" {
		c.PromptInstructionSurface = PromptInstructionSurfaceStandard
	} else {
		switch PromptInstructionSurface(instructionSurface) {
		case PromptInstructionSurfaceStandard, PromptInstructionSurfaceMinimal, PromptInstructionSurfaceMinimalAutonomy, PromptInstructionSurfaceMinimalExistence:
			c.PromptInstructionSurface = PromptInstructionSurface(instructionSurface)
		default:
			return fmt.Errorf("QUINE_PROMPT_INSTRUCTION_SURFACE=%q: must be %q, %q, %q, or %q",
				instructionSurface, PromptInstructionSurfaceStandard, PromptInstructionSurfaceMinimal, PromptInstructionSurfaceMinimalAutonomy, PromptInstructionSurfaceMinimalExistence)
		}
	}

	runtimeSurface := os.Getenv(EnvPromptRuntimeSurface)
	if runtimeSurface == "" {
		c.PromptRuntimeSurface = PromptRuntimeSurfaceVisible
	} else {
		switch PromptRuntimeSurface(runtimeSurface) {
		case PromptRuntimeSurfaceVisible, PromptRuntimeSurfaceHidden:
			c.PromptRuntimeSurface = PromptRuntimeSurface(runtimeSurface)
		default:
			return fmt.Errorf("QUINE_PROMPT_RUNTIME_SURFACE=%q: must be %q or %q",
				runtimeSurface, PromptRuntimeSurfaceVisible, PromptRuntimeSurfaceHidden)
		}
	}

	persona := strings.TrimSpace(os.Getenv(EnvPromptPersona))
	switch PromptPersona(persona) {
	case PromptPersonaNone, PromptPersonaCoder, PromptPersonaAnalyst, PromptPersonaEngineer, PromptPersonaArchitect, PromptPersonaExplorer, PromptPersonaSteward, PromptPersonaCartographer, PromptPersonaGardener, PromptPersonaWitness, PromptPersonaCatalyst, PromptPersonaSkeptic:
		c.PromptPersona = PromptPersona(persona)
	default:
		return fmt.Errorf("QUINE_PROMPT_PERSONA=%q: must be one of: %q, %q, %q, %q, %q, %q, %q, %q, %q, %q, or %q",
			persona, PromptPersonaCoder, PromptPersonaAnalyst, PromptPersonaEngineer, PromptPersonaArchitect, PromptPersonaExplorer, PromptPersonaSteward, PromptPersonaCartographer, PromptPersonaGardener, PromptPersonaWitness, PromptPersonaCatalyst, PromptPersonaSkeptic)
	}

	c.FailOnImpossible, err = envBoolDefault(EnvFailOnImpossible, true)
	if err != nil {
		return err
	}
	// When no mission argv is supplied, a frozen chat model defaults to an
	// operator-wait ("Ready.") posture. This gate controls how much of an
	// autonomous-process framing the missionless opening identity carries, as
	// two independently ablatable clauses (for the A0 prime-directive ablation):
	//   off       bare identity (reproduces the "Ready." collapse)
	//   autonomy  + "act on own judgment, do not wait" (the anti-idle clause)
	//   sensing   + "may sense/use the working directory" (the cwd-attention clause)
	//   full      autonomy + sensing
	// None of these supply a task goal. Default off preserves the bare prompt;
	// "0"/"1" are accepted as off/full for backward compatibility.
	c.NoMissionAutonomy, err = parseNoMissionAutonomy(os.Getenv(EnvNoMissionAutonomy))
	if err != nil {
		return err
	}
	c.PromptCtlPhysics, err = envBoolDefault(EnvPromptCtl, true)
	if err != nil {
		return err
	}
	c.PromptImplDetails, err = envBoolDefault(EnvPromptImplDetails, false)
	if err != nil {
		return err
	}
	c.MemoryStrategyHints, err = envBoolDefault(EnvMemoryStrategyHints, true)
	if err != nil {
		return err
	}

	// --- Optional User-Agent ---
	c.UserAgent = os.Getenv(EnvUserAgent)

	// --- Optional Thinking Budget ---
	c.ThinkingBudget = os.Getenv(EnvThinkingBudget)
	if c.ThinkingBudget == "" {
		c.ThinkingBudget = "high"
	} else {
		switch c.ThinkingBudget {
		case "off", "low", "medium", "high", "xhigh":
			// valid values
		default:
			return fmt.Errorf("QUINE_THINKING_BUDGET=%q: must be \"off\", \"low\", \"medium\", \"high\", or \"xhigh\"", c.ThinkingBudget)
		}
	}
	c.ServiceTier = strings.TrimSpace(os.Getenv(EnvModelServiceTier))
	switch c.ServiceTier {
	case "", "priority", "flex":
	case "fast":
		c.ServiceTier = "priority"
	default:
		return fmt.Errorf("QUINE_MODEL_SERVICE_TIER=%q: must be \"priority\", \"fast\", or \"flex\"", c.ServiceTier)
	}
	return nil
}

func loadSelfReentryConfig(c *Config) error {
	var err error
	if legacyTarget := strings.TrimSpace(os.Getenv(EnvSelfReentryTarget)); legacyTarget != "" {
		return fmt.Errorf("QUINE_SELF_REENTRY_TARGET has been removed; use QUINE_SELF_REENTRY_MODE=%q or %q", SelfReentryModeSelf, SelfReentryModeExecutablePath)
	}
	c.SelfReentryMode, err = resolveSelfReentryMode(os.Getenv(EnvSelfReentryMode))
	if err != nil {
		return err
	}
	if c.ExecEnabled || c.ForkEnabled() || c.SpawnEnabled() {
		c.SelfReentryTarget, err = resolveSelfReentryTarget(runtime.GOOS, c.SelfReentryMode, c.ExecutablePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadToolGates(c *Config) error {
	var err error
	c.AnchorMemoryEnabled, err = envBoolDefault(EnvAnchorMemory, false)
	if err != nil {
		return err
	}
	c.IdleEnabled, err = envBoolDefault(EnvIdleEnabled, false)
	if err != nil {
		return err
	}
	c.EmptyAssistantSuccess, err = envBoolDefault(EnvEmptyAssistantSuccess, false)
	if err != nil {
		return err
	}
	c.ReadyTextAutoIdle, err = envBoolDefault(EnvReadyTextAutoIdle, false)
	if err != nil {
		return err
	}
	c.ToolGates.ForkEnabled, err = envBoolDefault(EnvForkEnabled, true)
	if err != nil {
		return err
	}
	c.ToolGates.AnchorFoldEnabled, err = envBoolDefault(EnvAnchorFoldEnabled, true)
	if err != nil {
		return err
	}
	c.ToolGates.AnchorMarkEnabled, err = envBoolDefault(EnvAnchorMarkEnabled, true)
	if err != nil {
		return err
	}
	c.ToolGates.ExitEnabled, err = envBoolDefault(EnvExitEnabled, true)
	if err != nil {
		return err
	}
	c.ExecEnabled, err = envBoolDefault(EnvExecEnabled, true)
	if err != nil {
		return err
	}
	c.SpawnEnabledFlag, err = envBoolDefault(EnvSpawnEnabled, false)
	if err != nil {
		return err
	}
	c.AgentsMDEnabled, err = envBool01Default(EnvAgentsMDEnabled, false)
	if err != nil {
		return err
	}
	if c.AgentsMDEnabled {
		c.AgentsMDPath, err = DiscoverSingleAgentsMD(c)
		if err != nil {
			return err
		}
	}
	c.AgentsSkillsEnabled, err = envBool01Default(EnvAgentsSkillsEnabled, false)
	if err != nil {
		return err
	}
	if c.AgentsSkillsEnabled {
		c.Skills, err = LoadSkills(c)
		if err != nil {
			return err
		}
	}
	c.VisionEnabled, err = envBoolDefault(EnvVisionEnabled, true)
	if err != nil {
		return err
	}
	c.ToolGates.ShInteractiveEnabled, err = envBoolDefault(EnvShInteractiveEnabled, true)
	if err != nil {
		return err
	}
	c.EphemeralBody, err = envBoolDefault(EnvEphemeralBodyEnabled, false)
	if err != nil {
		return err
	}
	c.InitialUserMessage = os.Getenv(EnvInitialUserMessage)
	c.SuppressInitialBegin, err = envBoolDefault(EnvSuppressInitialBegin, false)
	if err != nil {
		return err
	}
	c.SelfSourceCodeEnabled, err = envBoolDefault(EnvSelfSourceCodeEnabled, false)
	if err != nil {
		return err
	}
	c.PeerDiscoveryEnabled, err = envBoolDefault(EnvPeerDiscoveryEnabled, false)
	if err != nil {
		return err
	}
	c.ToolGates.FSMutationTelemetry, err = envBoolDefault(EnvFSMutationTelemetry, true)
	if err != nil {
		return err
	}
	return nil
}

// baseEnv returns the common environment variable slice shared by
// ChildEnv and ExecEnv. depth and parentSession are parameterized
// since they differ between the two callers.
func (c *Config) baseEnv(depth int, parentSession string) []string {
	env := []string{
		envKV(EnvModelID, c.ModelID),
		envKV(EnvAPIType, c.Provider),
		envKV(EnvAPIBase, c.APIBase),
		envKV(EnvAPIKey, c.APIKey),
		envKV(EnvMaxDepth, strconv.Itoa(c.MaxDepth)),
		envKV(EnvDepth, strconv.Itoa(depth)),
		envKV(EnvParentSession, parentSession),
		envKV(EnvMaxConcurrent, strconv.Itoa(c.MaxConcurrent)),
		envKV(EnvMaxAgents, strconv.Itoa(c.MaxAgents)),
		envKV(EnvForkDefaultTimeout, strconv.Itoa(c.ForkDefaultTimeoutSeconds)),
		envKV(EnvShDefaultTimeout, strconv.Itoa(c.ShTimeout)),
		envKV(EnvShTimeoutOverride, bool01(c.ShTimeoutOverrideEnabled())),
		envKV(EnvShStdinEnabled, bool01(c.ShStdinEnabled())),
		envKV(EnvShDetachEnabled, bool01(c.ShDetachEnabled())),
		envKV(EnvOutputTruncate, strconv.Itoa(c.OutputTruncate)),
		envKV(EnvDataDir, c.DataDir),
		envKV(EnvWorkDir, c.WorkDir),
		envKV(EnvShell, c.Shell),
		envKV(EnvShNetwork, c.ShNetwork),
		envKV(EnvSelfReentryMode, string(c.SelfReentryMode)),
		envKV(EnvMaxTurns, strconv.Itoa(c.MaxTurns)),
		envKV(EnvWallClockExitSeconds, strconv.Itoa(c.WallClockExitSeconds)),
		envKV(EnvPromptMetaphor, string(c.PromptMetaphor)),
		envKV(EnvPromptSelfModel, string(c.SelfModelMode())),
		envKV(EnvPromptInstructionSurface, string(c.InstructionSurfaceMode())),
		envKV(EnvPromptRuntimeSurface, string(c.RuntimeSurfaceMode())),
		envKV(EnvPromptPersona, string(c.PersonaMode())),
		envKV(EnvPromptCtl, bool01(c.PromptCtlPhysics)),
		envKV(EnvPromptImplDetails, bool01(c.PromptImplDetails)),
		envKV(EnvPeerDiscoveryEnabled, bool01(c.PeerDiscoveryEnabled)),
		envKV(EnvPeerDiscoveryHeartbeat, strconv.Itoa(c.PeerDiscoveryHeartbeatMS)),
		envKV(EnvFSMutationTelemetry, bool01(c.FSMutationTelemetryEnabled())),
		envKV(EnvFailOnImpossible, bool01(c.FailOnImpossible)),
		envKV(EnvNoMissionAutonomy, c.NoMissionAutonomyLevel()),
		envKV(EnvEmptyAssistantSuccess, bool01(c.EmptyAssistantSuccess)),
		envKV(EnvReadyTextAutoIdle, bool01(c.ReadyTextAutoIdle)),
		envKV(EnvContextWindow, strconv.Itoa(c.ContextWindow)),
		envKV(EnvMemoryWarnTokens, strconv.Itoa(c.MemoryWarnTokens)),
		envKV(EnvMemoryDangerTokens, strconv.Itoa(c.MemoryDangerTokens)),
		envKV(EnvMemoryDeathTokens, strconv.Itoa(c.MemoryDeathTokens)),
		envKV(EnvMemoryStrategyHints, bool01(c.MemoryStrategyHints)),
	}
	if c.RetentionDir != "" {
		env = append(env, envKV(EnvRetentionDir, c.RetentionDir))
	}

	if c.WorkspaceEnabled {
		env = append(env,
			envKV(EnvWorkspaceRoot, c.WorkspaceRoot),
			envKV(EnvWorkspace, c.Workspace),
			envKV(EnvWorkspaceBackend, c.WorkspaceBackend),
			envKV(EnvWorkspaceOverlayDriver, c.WorkspaceOverlayDriver),
			envKV(EnvWorkspaceRevisionMode, string(c.WorkspaceRevisionMode)),
			envKV(EnvWorkspaceCurrentRevision, c.WorkspaceCurrentRevision),
			envKV(EnvWorkspaceSession, c.WorkspaceSession),
			envKV(EnvWorkspaceOwner, bool01(c.WorkspaceOwner)),
			envKV(EnvWorkspaceCommitOnSignal, bool01(c.WorkspaceCommitOnSignal)),
		)
	}

	if configDir := os.Getenv(EnvConfigDir); configDir != "" {
		env = append(env, envKV(EnvConfigDir, configDir))
	}

	// Propagate custom User-Agent if set
	if c.UserAgent != "" {
		env = append(env, envKV(EnvUserAgent, c.UserAgent))
	}

	// Propagate thinking budget if set
	if c.ThinkingBudget != "" {
		env = append(env, envKV(EnvThinkingBudget, c.ThinkingBudget))
	}
	if c.ServiceTier != "" {
		env = append(env, envKV(EnvModelServiceTier, c.ServiceTier))
	}
	if c.AnchorMemoryEnabled {
		env = append(env, envKV(EnvAnchorMemory, "1"))
		if !c.AnchorFoldEnabled() {
			env = append(env, envKV(EnvAnchorFoldEnabled, "0"))
		}
		if !c.AnchorMarkEnabled() {
			env = append(env, envKV(EnvAnchorMarkEnabled, "0"))
		}
	}
	if c.IdleEnabled {
		env = append(env, envKV(EnvIdleEnabled, "1"))
	} else {
		env = append(env, envKV(EnvIdleEnabled, "0"))
	}
	if c.ExitEnabled() {
		env = append(env, envKV(EnvExitEnabled, "1"))
	} else {
		env = append(env, envKV(EnvExitEnabled, "0"))
	}
	if c.ExecEnabled {
		env = append(env, envKV(EnvExecEnabled, "1"))
	} else {
		env = append(env, envKV(EnvExecEnabled, "0"))
	}
	if c.SpawnEnabled() {
		env = append(env, envKV(EnvSpawnEnabled, "1"))
	} else {
		env = append(env, envKV(EnvSpawnEnabled, "0"))
	}
	if c.AgentsMDEnabled {
		env = append(env, envKV(EnvAgentsMDEnabled, "1"))
	} else {
		env = append(env, envKV(EnvAgentsMDEnabled, "0"))
	}
	if c.AgentsSkillsEnabled {
		env = append(env, envKV(EnvAgentsSkillsEnabled, "1"))
	} else {
		env = append(env, envKV(EnvAgentsSkillsEnabled, "0"))
	}
	if c.VisionEnabled {
		env = append(env, envKV(EnvVisionEnabled, "1"))
	} else {
		env = append(env, envKV(EnvVisionEnabled, "0"))
	}
	if c.ForkEnabled() {
		env = append(env, envKV(EnvForkEnabled, "1"))
	} else {
		env = append(env, envKV(EnvForkEnabled, "0"))
	}
	if c.ShInteractiveEnabled() {
		env = append(env, envKV(EnvShInteractiveEnabled, "1"))
	} else {
		env = append(env, envKV(EnvShInteractiveEnabled, "0"))
	}
	if c.ForkWorldEnabled {
		env = append(env, envKV(EnvForkWorldEnabled, "1"))
	} else {
		env = append(env, envKV(EnvForkWorldEnabled, "0"))
	}
	if c.EphemeralBody {
		env = append(env, envKV(EnvEphemeralBodyEnabled, "1"))
	} else {
		env = append(env, envKV(EnvEphemeralBodyEnabled, "0"))
	}
	if c.SuppressInitialBegin {
		env = append(env, envKV(EnvSuppressInitialBegin, "1"))
	} else {
		env = append(env, envKV(EnvSuppressInitialBegin, "0"))
	}
	if c.InitialUserMessage != "" {
		env = append(env, envKV(EnvInitialUserMessage, c.InitialUserMessage))
	}
	if c.SelfSourceCodeEnabled {
		env = append(env, envKV(EnvSelfSourceCodeEnabled, "1"))
	} else {
		env = append(env, envKV(EnvSelfSourceCodeEnabled, "0"))
	}

	return env
}

func envKV(key, value string) string {
	return key + "=" + value
}

// ChildEnv returns a slice of "KEY=VALUE" environment variable strings
// suitable for spawning a child process. The child gets:
//   - QUINE_DEPTH incremented by 1
//   - QUINE_PARENT_SESSION set to the current SessionID
//   - All other config values inherited
//
// Note: QUINE_SESSION_ID and QUINE_TAPE_ID are intentionally NOT included.
// Each child ./quine process generates its own unique session/tape identity.
func (c *Config) ChildEnv() ([]string, error) {
	env := c.baseEnv(c.Depth+1, c.SessionID)
	if !c.WorkspaceEnabled {
		return env, nil
	}

	filtered := make([]string, 0, len(env)+1)
	for _, entry := range env {
		switch {
		case strings.HasPrefix(entry, EnvWorkspaceSession+"="):
			continue
		case strings.HasPrefix(entry, EnvWorkspaceOwner+"="):
			continue
		case strings.HasPrefix(entry, EnvWorkspaceBootstrap+"="):
			continue
		default:
			filtered = append(filtered, entry)
		}
	}

	if strings.TrimSpace(c.WorkspaceSession) != "" {
		filtered = append(filtered, envKV(EnvWorkspaceBootstrap, c.WorkspaceSession))
	}
	return filtered, nil
}

// ExecEnv returns a slice of "KEY=VALUE" environment variable strings
// suitable for exec'ing a fresh process (metamorphosis). Unlike ChildEnv:
//   - DEPTH is NOT incremented (same logical quine, new image)
//   - SESSION_ID is preserved (same logical quine across incarnations)
//   - PARENT_SESSION is preserved unchanged
//
// Note: QUINE_TAPE_ID is preserved for legacy tape continuity across exec.
func (c *Config) ExecEnv() ([]string, error) {
	env := c.baseEnv(0, c.ParentSession)
	env = append(env,
		envKV(EnvSessionID, c.SessionID),
		envKV(EnvTapeID, c.TapeID),
	)
	return env, nil
}

// ResolvedEnv renders the current resolved capability position as the body of
// the agent-root config/resolved.env read surface and the inc/<n>/config.env
// birth snapshots (registry-design-brief § B, work order T2.1).
//
// Content-source decision: the serialization is the exec-boundary env —
// ExecEnv()'s payload (baseEnv plus the exec-carried identity envs
// QUINE_SESSION_ID / QUINE_TAPE_ID) — because the self-reentry envp IS the
// capability position: env is the only injection channel at the process
// boundary. Two deliberate fidelity deviations from ExecEnv():
//
//   - QUINE_DEPTH renders the CURRENT depth (c.Depth), not the constant 0
//     that ExecEnv passes for the exec handover. resolved.env is a readout of
//     this process's position, not the next process's envp.
//   - QUINE_API_KEY is redacted to presence-only. The registry pins APIKey as
//     operator-only auth material whose value is never disclosed ("prompt
//     discloses only presence/absence, never the value"), and the
//     inc/<n>/config.env birth snapshots derived from this rendering are
//     retained lineage history — a raw credential must not land there.
//
// The rendering is a runtime-owned readout, never a config source: Load()
// reads envp only (preserved invariant 2). Lines are the raw boundary ABI
// "KEY=VALUE" strings with no shell quoting (brief D2's zero-translation
// stance).
func (c *Config) ResolvedEnv() []byte {
	env := c.baseEnv(c.Depth, c.ParentSession)
	env = append(env,
		envKV(EnvSessionID, c.SessionID),
		envKV(EnvTapeID, c.TapeID),
	)
	var b strings.Builder
	b.WriteString("# Quine resolved capability position (env syntax).\n")
	b.WriteString("# Runtime-owned readout: regenerated at bootstrap and after in-process config\n")
	b.WriteString("# mutation (workspace revision switch). Never read back as config — Load()\n")
	b.WriteString("# reads the process env only. Lines are raw KEY=VALUE boundary strings.\n")
	apiKeyPrefix := EnvAPIKey + "="
	for _, kv := range env {
		if strings.HasPrefix(kv, apiKeyPrefix) {
			if strings.TrimSpace(strings.TrimPrefix(kv, apiKeyPrefix)) == "" {
				b.WriteString("# " + EnvAPIKey + " is unset.\n")
			} else {
				b.WriteString("# " + EnvAPIKey + " is set; value redacted (auth material never lands on the config surface).\n")
			}
			continue
		}
		b.WriteString(kv)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

// CanRestoreWorld reports whether the current workspace configuration
// exposes explicit world-switch semantics.
func (c *Config) CanRestoreWorld() bool {
	if c == nil || !c.WorkspaceEnabled {
		return false
	}
	return c.WorkspaceRevisionMode == WorkspaceRevisionRestore
}

// parseNoMissionAutonomy normalizes the missionless-autonomy level. It accepts
// the named levels off|autonomy|sensing|full and the legacy 0/1 (off/full).
func parseNoMissionAutonomy(v string) (string, error) {
	switch strings.TrimSpace(v) {
	case "", "0", "off":
		return "off", nil
	case "1", "full":
		return "full", nil
	case "autonomy", "sensing":
		return strings.TrimSpace(v), nil
	default:
		return "", fmt.Errorf("%s=%q: must be off|autonomy|sensing|full (or 0/1)", EnvNoMissionAutonomy, v)
	}
}

func envBoolDefault(key string, def bool) (bool, error) {
	if v, ok := os.LookupEnv(key); ok {
		return parseEnvBool(key, v, def)
	}
	return def, nil
}

func parseEnvBool(key, value string, def bool) (bool, error) {
	if strings.TrimSpace(value) == "" {
		return def, nil
	}
	switch strings.TrimSpace(value) {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("%s=%q: must be \"0\" or \"1\"", key, value)
	}
}

func bool01(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func defaultMemoryWarnTokens(contextWindow int) int {
	warn := contextWindow / 16
	if warn < 2048 {
		return 2048
	}
	return warn
}

func defaultMemoryDangerTokens(contextWindow int) int {
	danger := contextWindow / 8
	if danger < 4096 {
		return 4096
	}
	return danger
}

func canonicalPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return real, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		parent := filepath.Dir(abs)
		realParent, parentErr := filepath.EvalSymlinks(parent)
		if parentErr == nil {
			return filepath.Join(realParent, filepath.Base(abs)), nil
		}
	}
	return abs, nil
}

func isPathWithin(root, child string) bool {
	if root == child {
		return true
	}
	rel, err := filepath.Rel(root, child)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// envInt reads an environment variable as int, returning def if unset.
func envInt(key string, def int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parsing %s=%q: %w", key, v, err)
	}
	return n, nil
}

// uuidV4 generates a random UUID v4 using crypto/rand.
func uuidV4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func stableSessionID() string {
	return fmt.Sprintf("sess_%s_%s", processStartTime.Format("20060102-150405"), mustUUIDV4())
}

func processRunID() string {
	return fmt.Sprintf("run_%s_%d_%s", processStartTime.Format("20060102-150405"), os.Getpid(), mustUUIDV4())
}

func mustUUIDV4() string {
	id, err := uuidV4()
	if err != nil {
		// Keep the process bootable even if entropy is temporarily unavailable.
		return fmt.Sprintf("fallback-%d-%d", os.Getpid(), time.Now().UnixNano())
	}
	return id
}

func nextTapeID(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "0001", nil
		}
		return "", err
	}

	maxID := 0
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		base := strings.TrimSuffix(name, ".jsonl")
		id, err := strconv.Atoi(base)
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}

	return fmt.Sprintf("%04d", maxID+1), nil
}
