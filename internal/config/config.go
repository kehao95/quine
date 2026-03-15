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

// processStartTime captures the host process birth time for deriving a stable
// process-shaped session identity when QUINE_SESSION_ID is not supplied.
var processStartTime = time.Now()

// TurnExhaustionPolicy controls runtime behavior when execution budget is exhausted.
type TurnExhaustionPolicy string

const (
	// TurnExhaustionHardFail terminates immediately when executions reach zero.
	TurnExhaustionHardFail TurnExhaustionPolicy = "hard_fail"
	// TurnExhaustionNearDeathExec allows one final inference where only exec is accepted.
	TurnExhaustionNearDeathExec TurnExhaustionPolicy = "near_death_exec"
)

// PromptMetaphor controls optional metaphor framing in the system prompt.
type PromptMetaphor string

const (
	PromptMetaphorOff           PromptMetaphor = "off"
	PromptMetaphorThermodynamic PromptMetaphor = "thermodynamic"
)

// WorkspaceRevisionMode controls which explicit revision primitives are
// available once workspace physics are enabled.
type WorkspaceRevisionMode string

const (
	WorkspaceRevisionNone    WorkspaceRevisionMode = "none"
	WorkspaceRevisionRestore WorkspaceRevisionMode = "restore"
)

// Config holds all runtime configuration for Quine.
// Every field is populated from environment variables by Load().
type Config struct {
	ModelID        string // QUINE_MODEL_ID (required)
	APIKey         string // QUINE_API_KEY (required)
	APIBase        string // QUINE_API_BASE (required)
	Provider       string // QUINE_API_TYPE (required): "openai", "anthropic", or "openai-responses"
	MaxDepth       int    // QUINE_MAX_DEPTH (default 0 = disabled)
	Depth          int    // QUINE_DEPTH (default 0)
	SessionID      string // QUINE_SESSION_ID (default auto <YYYYMMDD-HHMMSS>_<ppid>_<pid>)
	TapeID         string // QUINE_TAPE_ID (default auto 4-digit increment per session; changes across exec incarnations)
	ParentSession  string // QUINE_PARENT_SESSION
	MaxConcurrent  int    // QUINE_MAX_CONCURRENT (default 0 = disabled)
	MaxAgents      int    // QUINE_MAX_AGENTS (default 0 = disabled)
	ShTimeout      int    // QUINE_SH_TIMEOUT in seconds (default 600)
	OutputTruncate int    // QUINE_OUTPUT_TRUNCATE in bytes (default 20480)
	DataDir        string // QUINE_DATA_DIR: durable runtime-state root (default ".quine/"); must stay outside QUINE_WORKSPACE_ROOT
	WorkDir        string // QUINE_WORK_DIR (default pwd at startup)
	Shell          string // QUINE_SHELL (default "/bin/sh")
	MaxTurns       int    // QUINE_MAX_TURNS (default 0 = disabled)
	// QUINE_TURN_EXHAUSTION_POLICY (default "hard_fail")
	TurnExhaustionPolicy TurnExhaustionPolicy
	// QUINE_PROMPT_METAPHOR (default "off")
	PromptMetaphor PromptMetaphor
	ContextWindow  int               // QUINE_CONTEXT_WINDOW (default 128000)
	Wisdom         map[string]string // QUINE_WISDOM_* env vars (key without prefix -> value)
	OriginalIntent string            // QUINE_ORIGINAL_INTENT (preserved across exec for mission continuity)

	// Escalation tier (optional — all empty = single-model mode)
	SmartModelID   string // QUINE_SMART_MODEL_ID
	SmartAPIKey    string // QUINE_SMART_API_KEY (falls back to APIKey if empty)
	SmartAPIBase   string // QUINE_SMART_API_BASE (falls back to APIBase if empty)
	SmartProvider  string // QUINE_SMART_API_TYPE (falls back to Provider if empty)
	StallThreshold int    // QUINE_STALL_THRESHOLD (default 5): triggers STALL warning after N turns on same goal
	Escalated      bool   // Runtime flag, not from env. Tracks whether escalation has occurred.

	// Optional HTTP headers
	UserAgent string // QUINE_USER_AGENT (optional): custom User-Agent header for API requests

	// Thinking budget for reasoning models (Kimi, o1, etc.)
	ThinkingBudget string // QUINE_THINKING_BUDGET: "off", "low", "medium", "high" (default: "high")

	// Anchor memory tools (mark/unfold). Controls tool exposure only; paths stay stable.
	AnchorMemoryEnabled bool // QUINE_ANCHOR_MEMORY (default false)

	// Linux-only overlay workspace physics (Evolution 11.1)
	WorkspaceEnabled      bool                  // implied by QUINE_WORKSPACE* configuration
	WorkspaceRoot         string                // Stable host world boundary for task-visible file operations (canonical absolute path)
	Workspace             string                // Current writable workspace scope within WorkspaceRoot (canonical absolute path)
	WorkspaceBackend      string                // QUINE_WORKSPACE_BACKEND: "overlay" (default) or "direct"
	WorkspaceRevisionMode WorkspaceRevisionMode // QUINE_WORKSPACE_REVISION_MODE: "none" or "restore"
	WorkspaceCurrentRevision string             // QUINE_WORKSPACE_CURRENT_REVISION: runtime-propagated current world revision handle
	WorkspaceSession      string                // Stable overlay-state namespace under QUINE_DATA_DIR shared across exec/fork
	WorkspaceOwner        bool                  // True only for the process responsible for committing/rolling back workspace state
}

// APIModelID returns the model ID to use in API calls.
func (c *Config) APIModelID() string {
	return c.ModelID
}

// CanEscalate returns true if escalation is configured and hasn't occurred yet.
func (c *Config) CanEscalate() bool {
	return c.SmartModelID != "" && !c.Escalated
}

// UsesNearDeathContinuation returns true when execution exhaustion should enter
// a final exec-only continuation window.
func (c *Config) UsesNearDeathContinuation() bool {
	return c.MaxTurns > 0 && c.TurnExhaustionPolicy == TurnExhaustionNearDeathExec
}

// ThermodynamicMetaphorEnabled reports whether thermodynamic framing should be
// added to the system prompt.
func (c *Config) ThermodynamicMetaphorEnabled() bool {
	return c.PromptMetaphor == PromptMetaphorThermodynamic
}

// RuntimeRoot returns the durable runtime-state root shared by the process tree.
// It holds tapes, session logs, job directories, coordination locks, and
// workspace overlay state. It is distinct from the task workspace surface.
func (c *Config) RuntimeRoot() string {
	return c.DataDir
}

// LockDir returns the coordination directory used for shared semaphore and
// agent-registry files under the runtime root.
func (c *Config) LockDir() string {
	return filepath.Join(c.RuntimeRoot(), "locks")
}

// SessionLogPath returns the flat operational log path for a session.
func (c *Config) SessionLogPath(sessionID string) string {
	if sessionID == "" {
		sessionID = c.SessionID
	}
	return filepath.Join(c.RuntimeRoot(), sessionID+".log")
}

// TapeDir returns the append-only tape directory for a session.
func (c *Config) TapeDir(sessionID string) string {
	if sessionID == "" {
		sessionID = c.SessionID
	}
	return filepath.Join(c.RuntimeRoot(), "tapes", sessionID)
}

// TapePath returns the append-only JSONL tape path for a tape incarnation
// under the current stable session.
func (c *Config) TapePath(tapeID string) string {
	if tapeID == "" {
		tapeID = c.TapeID
	}
	return filepath.Join(c.TapeDir(""), tapeID+".jsonl")
}

// TapeYAMLPath returns the YAML mirror path for a tape incarnation under the
// current stable session.
func (c *Config) TapeYAMLPath(tapeID string) string {
	if tapeID == "" {
		tapeID = c.TapeID
	}
	return filepath.Join(c.TapeDir(""), tapeID+".log.yaml")
}

// AgentRoot returns the synthetic process-filesystem root for this session.
// The root remains stable across exec because SessionID remains stable.
func (c *Config) AgentRoot() string {
	return filepath.Join(c.RuntimeRoot(), "agent", c.SessionID)
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
//   - QUINE_API_TYPE:   Wire protocol: "openai", "anthropic", or "codex-oauth"
//   - QUINE_API_BASE:   API base URL (e.g. "https://api.anthropic.com", "https://api.openai.com")
//   - QUINE_API_KEY:    API key
func Load() (*Config, error) {
	c := &Config{}

	// --- 4 required fields ---
	c.ModelID = os.Getenv("QUINE_MODEL_ID")
	if c.ModelID == "" {
		return nil, fmt.Errorf("QUINE_MODEL_ID is required")
	}

	c.Provider = os.Getenv("QUINE_API_TYPE")
	if c.Provider == "" {
		return nil, fmt.Errorf("QUINE_API_TYPE is required (\"openai\", \"anthropic\", or \"openai-responses\")")
	}
	if c.Provider != "openai" && c.Provider != "anthropic" && c.Provider != "openai-responses" {
		return nil, fmt.Errorf("unsupported QUINE_API_TYPE=%q: must be \"openai\", \"anthropic\", or \"openai-responses\"", c.Provider)
	}

	c.APIBase = os.Getenv("QUINE_API_BASE")
	if c.APIBase == "" {
		return nil, fmt.Errorf("QUINE_API_BASE is required")
	}

	c.APIKey = os.Getenv("QUINE_API_KEY")
	if c.APIKey == "" {
		return nil, fmt.Errorf("QUINE_API_KEY is required")
	}

	// --- Optional string fields ---
	c.ParentSession = os.Getenv("QUINE_PARENT_SESSION")

	// --- Linux-only overlay workspace physics ---
	c.WorkspaceRoot = strings.TrimSpace(os.Getenv("QUINE_WORKSPACE_ROOT"))
	c.Workspace = strings.TrimSpace(os.Getenv("QUINE_WORKSPACE"))
	c.WorkspaceBackend = strings.TrimSpace(os.Getenv("QUINE_WORKSPACE_BACKEND"))
	revisionMode := strings.TrimSpace(os.Getenv("QUINE_WORKSPACE_REVISION_MODE"))
	c.WorkspaceCurrentRevision = strings.TrimSpace(os.Getenv("QUINE_WORKSPACE_CURRENT_REVISION"))
	c.WorkspaceSession = strings.TrimSpace(os.Getenv("QUINE_WORKSPACE_SESSION"))
	c.WorkspaceOwner = envBoolDefault("QUINE_WORKSPACE_OWNER", true)
	if strings.TrimSpace(os.Getenv("QUINE_WORKSPACE_SOURCE")) != "" {
		return nil, fmt.Errorf("QUINE_WORKSPACE_SOURCE has been removed; use QUINE_WORKSPACE_ROOT and QUINE_WORKSPACE")
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
		if runtime.GOOS != "linux" {
			return nil, fmt.Errorf("workspace physics are only supported on Linux")
		}
		if c.WorkspaceBackend == "" {
			c.WorkspaceBackend = "overlay"
		}
		switch c.WorkspaceBackend {
		case "overlay", "direct":
		default:
			return nil, fmt.Errorf("QUINE_WORKSPACE_BACKEND=%q: must be \"overlay\" or \"direct\"", c.WorkspaceBackend)
		}
		if revisionMode == "" {
			c.WorkspaceRevisionMode = WorkspaceRevisionRestore
		} else {
			switch WorkspaceRevisionMode(revisionMode) {
			case WorkspaceRevisionNone, WorkspaceRevisionRestore:
				c.WorkspaceRevisionMode = WorkspaceRevisionMode(revisionMode)
			default:
				return nil, fmt.Errorf("QUINE_WORKSPACE_REVISION_MODE=%q: must be %q or %q",
					revisionMode, WorkspaceRevisionNone, WorkspaceRevisionRestore)
			}
		}
		if c.WorkspaceRoot == "" {
			return nil, fmt.Errorf("QUINE_WORKSPACE_ROOT or QUINE_WORKSPACE is required when workspace physics are enabled")
		}
		if c.Workspace == "" {
			c.Workspace = c.WorkspaceRoot
		}
		root, err := canonicalPath(c.WorkspaceRoot)
		if err != nil {
			return nil, fmt.Errorf("canonicalize workspace root: %w", err)
		}
		workspace, err := canonicalPath(c.Workspace)
		if err != nil {
			return nil, fmt.Errorf("canonicalize workspace: %w", err)
		}
		if !isPathWithin(root, workspace) {
			return nil, fmt.Errorf("workspace %q must be within workspace root %q", workspace, root)
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
				return nil, fmt.Errorf("QUINE_WORKSPACE_REVISION_MODE=%q requires workspace physics to be enabled", revisionMode)
			}
		}
	}

	// --- Integer fields with defaults ---
	var err error

	c.ContextWindow, err = envInt("QUINE_CONTEXT_WINDOW", 128_000)
	if err != nil {
		return nil, err
	}

	c.MaxDepth, err = envInt("QUINE_MAX_DEPTH", 0)
	if err != nil {
		return nil, err
	}

	c.Depth, err = envInt("QUINE_DEPTH", 0)
	if err != nil {
		return nil, err
	}

	c.MaxConcurrent, err = envInt("QUINE_MAX_CONCURRENT", 0)
	if err != nil {
		return nil, err
	}

	c.MaxAgents, err = envInt("QUINE_MAX_AGENTS", 0)
	if err != nil {
		return nil, err
	}

	c.ShTimeout, err = envInt("QUINE_SH_TIMEOUT", 600)
	if err != nil {
		return nil, err
	}

	c.OutputTruncate, err = envInt("QUINE_OUTPUT_TRUNCATE", 20480)
	if err != nil {
		return nil, err
	}

	c.MaxTurns, err = envInt("QUINE_MAX_TURNS", 0)
	if err != nil {
		return nil, err
	}

	// --- Depth check (disabled when MaxDepth <= 0) ---
	if c.MaxDepth > 0 && c.Depth >= c.MaxDepth {
		return nil, ErrDepthExceeded
	}

	// --- Session ID ---
	c.SessionID = os.Getenv("QUINE_SESSION_ID")
	if c.SessionID == "" {
		c.SessionID = processSessionID()
	}

	// --- Data dir ---
	c.DataDir = os.Getenv("QUINE_DATA_DIR")
	if c.DataDir == "" {
		c.DataDir = ".quine/"
	}
	c.TapeID = strings.TrimSpace(os.Getenv("QUINE_TAPE_ID"))
	if c.TapeID == "" {
		c.TapeID, err = nextTapeID(c.TapeDir(""))
		if err != nil {
			return nil, fmt.Errorf("generating tape ID: %w", err)
		}
	}
	if c.WorkspaceEnabled {
		if c.WorkspaceSession == "" {
			c.WorkspaceSession = c.SessionID
		}
		dataDirReal, err := canonicalPath(c.DataDir)
		if err != nil {
			return nil, fmt.Errorf("canonicalize data dir: %w", err)
		}
		if isPathWithin(c.WorkspaceRoot, dataDirReal) {
			return nil, fmt.Errorf("QUINE_DATA_DIR %q must be outside workspace root %q", dataDirReal, c.WorkspaceRoot)
		}
		if _, err := os.Stat(c.WorkspaceRoot); err != nil {
			return nil, fmt.Errorf("workspace root %q must exist: %w", c.WorkspaceRoot, err)
		}
	}

	// --- Shell ---
	c.Shell = os.Getenv("QUINE_SHELL")
	if c.Shell == "" {
		c.Shell = "/bin/sh"
	}

	// --- Work dir ---
	c.WorkDir = os.Getenv("QUINE_WORK_DIR")
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

	// --- Wisdom (QUINE_WISDOM_* env vars) ---
	c.Wisdom = loadWisdom()

	// --- Original Intent (preserved across exec for mission continuity) ---
	c.OriginalIntent = os.Getenv("QUINE_ORIGINAL_INTENT")

	// --- Smart model config for escalation (optional) ---
	c.SmartModelID = os.Getenv("QUINE_SMART_MODEL_ID")
	c.SmartAPIKey = os.Getenv("QUINE_SMART_API_KEY")
	if c.SmartAPIKey == "" {
		c.SmartAPIKey = c.APIKey
	}
	c.SmartAPIBase = os.Getenv("QUINE_SMART_API_BASE")
	if c.SmartAPIBase == "" {
		c.SmartAPIBase = c.APIBase
	}
	c.SmartProvider = os.Getenv("QUINE_SMART_API_TYPE")
	if c.SmartProvider == "" {
		c.SmartProvider = c.Provider
	}

	c.StallThreshold, err = envInt("QUINE_STALL_THRESHOLD", 5)
	if err != nil {
		return nil, err
	}

	// --- Execution exhaustion policy ---
	// Policy is only meaningful when execution budget is enabled.
	// QUINE_MAX_TURNS=0 means disabled/unlimited, so ignore policy value.
	c.TurnExhaustionPolicy = TurnExhaustionHardFail
	if c.MaxTurns > 0 {
		policy := os.Getenv("QUINE_TURN_EXHAUSTION_POLICY")
		if policy == "" {
			c.TurnExhaustionPolicy = TurnExhaustionHardFail
		} else {
			switch TurnExhaustionPolicy(policy) {
			case TurnExhaustionHardFail, TurnExhaustionNearDeathExec:
				c.TurnExhaustionPolicy = TurnExhaustionPolicy(policy)
			default:
				return nil, fmt.Errorf("QUINE_TURN_EXHAUSTION_POLICY=%q: must be %q or %q",
					policy, TurnExhaustionHardFail, TurnExhaustionNearDeathExec)
			}
		}
	}

	// --- Prompt metaphor mode ---
	metaphor := os.Getenv("QUINE_PROMPT_METAPHOR")
	if metaphor == "" {
		c.PromptMetaphor = PromptMetaphorOff
	} else {
		switch PromptMetaphor(metaphor) {
		case PromptMetaphorOff, PromptMetaphorThermodynamic:
			c.PromptMetaphor = PromptMetaphor(metaphor)
		default:
			return nil, fmt.Errorf("QUINE_PROMPT_METAPHOR=%q: must be %q or %q",
				metaphor, PromptMetaphorOff, PromptMetaphorThermodynamic)
		}
	}

	// --- Optional User-Agent ---
	c.UserAgent = os.Getenv("QUINE_USER_AGENT")

	// --- Optional Thinking Budget ---
	c.ThinkingBudget = os.Getenv("QUINE_THINKING_BUDGET")
	if c.ThinkingBudget == "" {
		c.ThinkingBudget = "high" // Default to high reasoning effort
	} else {
		switch c.ThinkingBudget {
		case "off", "low", "medium", "high":
			// valid values
		default:
			return nil, fmt.Errorf("QUINE_THINKING_BUDGET=%q: must be \"off\", \"low\", \"medium\", or \"high\"", c.ThinkingBudget)
		}
	}

	// --- Optional anchor-memory tool gate ---
	c.AnchorMemoryEnabled = envBoolDefault("QUINE_ANCHOR_MEMORY", false)

	return c, nil
}

// baseEnv returns the common environment variable slice shared by
// ChildEnv and ExecEnv. depth and parentSession are parameterized
// since they differ between the two callers.
func (c *Config) baseEnv(depth int, parentSession string) []string {
	env := []string{
		"QUINE_MODEL_ID=" + c.ModelID,
		"QUINE_API_TYPE=" + c.Provider,
		"QUINE_API_BASE=" + c.APIBase,
		"QUINE_API_KEY=" + c.APIKey,
		"QUINE_MAX_DEPTH=" + strconv.Itoa(c.MaxDepth),
		"QUINE_DEPTH=" + strconv.Itoa(depth),
		"QUINE_PARENT_SESSION=" + parentSession,
		"QUINE_MAX_CONCURRENT=" + strconv.Itoa(c.MaxConcurrent),
		"QUINE_MAX_AGENTS=" + strconv.Itoa(c.MaxAgents),
		"QUINE_SH_TIMEOUT=" + strconv.Itoa(c.ShTimeout),
		"QUINE_OUTPUT_TRUNCATE=" + strconv.Itoa(c.OutputTruncate),
		"QUINE_DATA_DIR=" + c.DataDir,
		"QUINE_SHELL=" + c.Shell,
		"QUINE_MAX_TURNS=" + strconv.Itoa(c.MaxTurns),
		"QUINE_TURN_EXHAUSTION_POLICY=" + string(c.TurnExhaustionPolicy),
		"QUINE_PROMPT_METAPHOR=" + string(c.PromptMetaphor),
		"QUINE_CONTEXT_WINDOW=" + strconv.Itoa(c.ContextWindow),
	}

	if c.WorkspaceEnabled {
		env = append(env,
			"QUINE_WORKSPACE_ROOT="+c.WorkspaceRoot,
			"QUINE_WORKSPACE="+c.Workspace,
			"QUINE_WORKSPACE_BACKEND="+c.WorkspaceBackend,
			"QUINE_WORKSPACE_REVISION_MODE="+string(c.WorkspaceRevisionMode),
			"QUINE_WORKSPACE_CURRENT_REVISION="+c.WorkspaceCurrentRevision,
			"QUINE_WORKSPACE_SESSION="+c.WorkspaceSession,
			"QUINE_WORKSPACE_OWNER="+strconv.FormatBool(c.WorkspaceOwner),
		)
	}

	// Pass through QUINE_WISDOM_* env vars for state transfer across exec boundaries
	for key, value := range c.Wisdom {
		env = append(env, wisdomPrefix+key+"="+value)
	}

	if configDir := os.Getenv("QUINE_CONFIG_DIR"); configDir != "" {
		env = append(env, "QUINE_CONFIG_DIR="+configDir)
	}

	// Propagate smart config to children for escalation
	if c.SmartModelID != "" {
		env = append(env,
			"QUINE_SMART_MODEL_ID="+c.SmartModelID,
			"QUINE_SMART_API_KEY="+c.SmartAPIKey,
			"QUINE_SMART_API_BASE="+c.SmartAPIBase,
			"QUINE_SMART_API_TYPE="+c.SmartProvider,
		)
	}

	// Propagate custom User-Agent if set
	if c.UserAgent != "" {
		env = append(env, "QUINE_USER_AGENT="+c.UserAgent)
	}

	// Propagate thinking budget if set
	if c.ThinkingBudget != "" {
		env = append(env, "QUINE_THINKING_BUDGET="+c.ThinkingBudget)
	}
	if c.AnchorMemoryEnabled {
		env = append(env, "QUINE_ANCHOR_MEMORY=1")
	}

	return env
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
	return c.baseEnv(c.Depth+1, c.SessionID), nil
}

// ExecEnv returns a slice of "KEY=VALUE" environment variable strings
// suitable for exec'ing a fresh process (metamorphosis). Unlike ChildEnv:
//   - DEPTH is NOT incremented (fresh context = restart)
//   - SESSION_ID is preserved (same logical quine across incarnations)
//   - PARENT_SESSION is preserved unchanged
//   - ORIGINAL_INTENT is set to preserve the mission
//   - All QUINE_WISDOM_* vars are preserved (learned insights survive)
//
// Note: QUINE_TAPE_ID is intentionally NOT included. The new incarnation
// generates the next per-session tape ID via config.Load().
func (c *Config) ExecEnv(originalIntent string) ([]string, error) {
	env := c.baseEnv(0, c.ParentSession)
	env = append(env,
		"QUINE_SESSION_ID="+c.SessionID,
		"QUINE_ORIGINAL_INTENT="+originalIntent,
	)
	return env, nil
}

// CanRestoreWorld reports whether the current workspace configuration
// exposes restore_world semantics.
func (c *Config) CanRestoreWorld() bool {
	if c == nil || !c.WorkspaceEnabled {
		return false
	}
	return c.WorkspaceRevisionMode == WorkspaceRevisionRestore
}

func envBool(key string) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func envBoolDefault(key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes" || v == "on"
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

// processSessionID derives a readable identity from the current process start
// time and Unix lineage. Exec preserves this via QUINE_SESSION_ID, while forked
// children generate their own value when they call Load().
func processSessionID() string {
	return fmt.Sprintf("%s_%d_%d", processStartTime.Format("20060102-150405"), os.Getppid(), os.Getpid())
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

// wisdomPrefix is the environment variable prefix for wisdom transfer.
const wisdomPrefix = "QUINE_WISDOM_"

// loadWisdom scans all environment variables and collects those starting
// with QUINE_WISDOM_. It returns a map with keys stripped of the prefix.
func loadWisdom() map[string]string {
	wisdom := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, wisdomPrefix) {
			// Split on first "=" to get key=value
			key, value, found := strings.Cut(env, "=")
			if !found {
				continue
			}
			// Strip the prefix from the key
			shortKey := strings.TrimPrefix(key, wisdomPrefix)
			if shortKey != "" && value != "" {
				wisdom[shortKey] = value
			}
		}
	}
	return wisdom
}
