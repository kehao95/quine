package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// envVars is the full list of environment variables we manage in tests.
var envVars = []string{
	// Core model/provider identity (required runtime wiring).
	"QUINE_MODEL_ID",
	"QUINE_API_TYPE",
	"QUINE_API_BASE",
	"QUINE_API_KEY",
	"QUINE_DEBUG_REQUEST_BODY_DIR",

	// Process lineage/session identity.
	"QUINE_MAX_DEPTH",
	"QUINE_DEPTH",
	"QUINE_SESSION_ID",
	"QUINE_RUN_ID",
	"QUINE_TAPE_ID",
	"QUINE_PARENT_SESSION",

	// Runtime capacity limits.
	"QUINE_MAX_CONCURRENT",
	"QUINE_MAX_AGENTS",
	"QUINE_FORK_DEFAULT_TIMEOUT_SECONDS",
	"QUINE_MAX_TURNS",
	"QUINE_WALL_CLOCK_EXIT_SECONDS",

	// Execution/prompt behavior switches.
	"QUINE_PROMPT_METAPHOR",
	"QUINE_PROMPT_SELF_MODEL",
	"QUINE_PROMPT_INSTRUCTION_SURFACE",
	"QUINE_PROMPT_RUNTIME_SURFACE",
	"QUINE_PROMPT_PERSONA",
	"QUINE_PROMPT_CTL",
	"QUINE_PROMPT_IMPL_DETAILS",
	"QUINE_PEER_DISCOVERY_ENABLED",
	"QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS",
	"QUINE_FS_MUTATION_TELEMETRY_ENABLED",
	"QUINE_FAIL_ON_IMPOSSIBLE",
	"QUINE_EMPTY_ASSISTANT_SUCCESS",
	"QUINE_READY_TEXT_AUTO_IDLE",

	// Shell/tool runtime parameters.
	"QUINE_SH_DEFAULT_TIMEOUT_SECONDS",
	"QUINE_SH_TIMEOUT_OVERRIDE_ENABLED",
	"QUINE_SH_STDIN_ENABLED",
	"QUINE_SH_DETACH_ENABLED",
	"QUINE_OUTPUT_TRUNCATE",
	"QUINE_DATA_DIR",
	"QUINE_RETENTION_DIR",
	"QUINE_SHELL",
	"QUINE_SH_NETWORK",
	"QUINE_WORK_DIR",
	"QUINE_WORKSPACE_ROOT",
	"QUINE_WORKSPACE",
	"QUINE_WORKSPACE_BACKEND",
	"QUINE_WORKSPACE_OVERLAY_DRIVER",
	"QUINE_WORKSPACE_REVISION_MODE",
	"QUINE_WORKSPACE_CURRENT_REVISION",
	"QUINE_WORKSPACE_SOURCE",
	"QUINE_WORKSPACE_SESSION",
	"QUINE_WORKSPACE_OWNER",
	"QUINE_WORKSPACE_COMMIT_ON_SIGNAL",
	"QUINE_WORKSPACE_BOOTSTRAP",

	// Context and local config path.
	"QUINE_CONTEXT_WINDOW",
	"QUINE_THINKING_BUDGET",
	"QUINE_MODEL_SERVICE_TIER",
	"QUINE_MEMORY_WARN_TOKENS",
	"QUINE_MEMORY_DANGER_TOKENS",
	"QUINE_MEMORY_DEATH_TOKENS",
	"QUINE_MEMORY_STRATEGY_HINTS",
	"QUINE_CONFIG_DIR",
	"QUINE_ANCHOR_MEMORY",
	"QUINE_IDLE_ENABLED",
	"QUINE_FORK_ENABLED",
	"QUINE_EXIT_ENABLED",
	"QUINE_EXEC_ENABLED",
	"QUINE_SPAWN_ENABLED",
	"QUINE_AGENTS_MD_ENABLED",
	"QUINE_AGENTS_SKILLS_ENABLED",
	"QUINE_VISION_ENABLED",
	"QUINE_SH_INTERACTIVE_ENABLED",
	"QUINE_FORK_WORLD_ENABLED",
	"QUINE_EPHEMERAL_BODY_ENABLED",
	"QUINE_SUPPRESS_INITIAL_BEGIN",
	"QUINE_SELF_REENTRY_MODE",
	"QUINE_SELF_REENTRY_TARGET",
	"QUINE_SELF_SOURCE_CODE_ENABLED",
}

// clearEnv unsets all managed env vars and returns a restore function.
func clearEnv(t *testing.T) {
	t.Helper()
	saved := make(map[string]string)
	for _, k := range envVars {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range envVars {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

// setRequired sets the 4 required env vars.
func setRequired(t *testing.T) {
	t.Helper()
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-20250514")
	os.Setenv("QUINE_API_TYPE", "anthropic")
	os.Setenv("QUINE_API_BASE", "https://api.anthropic.com")
	os.Setenv("QUINE_API_KEY", "sk-test-key")
	if runtime.GOOS != "linux" {
		os.Setenv("QUINE_SELF_REENTRY_MODE", string(SelfReentryModeExecutablePath))
	}
}

func expectedTestSelfReentryTarget(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "linux" {
		return linuxSelfReentryTarget
	}
	executablePath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error: %v", err)
	}
	return executablePath
}

func TestResolveSelfReentryMode_DefaultsToSelf(t *testing.T) {
	got, err := resolveSelfReentryMode("")
	if err != nil {
		t.Fatalf("resolveSelfReentryMode() error = %v", err)
	}
	if got != SelfReentryModeSelf {
		t.Fatalf("resolveSelfReentryMode() = %q, want %q", got, SelfReentryModeSelf)
	}
}

func TestResolveSelfReentryMode_AcceptsExecutablePath(t *testing.T) {
	got, err := resolveSelfReentryMode("executable_path")
	if err != nil {
		t.Fatalf("resolveSelfReentryMode() error = %v", err)
	}
	if got != SelfReentryModeExecutablePath {
		t.Fatalf("resolveSelfReentryMode() = %q, want %q", got, SelfReentryModeExecutablePath)
	}
}

func TestResolveSelfReentryMode_RejectsInvalidValue(t *testing.T) {
	_, err := resolveSelfReentryMode("mystery")
	if err == nil || !strings.Contains(err.Error(), "QUINE_SELF_REENTRY_MODE") {
		t.Fatalf("expected invalid-mode error, got %v", err)
	}
}

func TestResolveSelfReentryTarget_SelfModeOnLinux(t *testing.T) {
	got, err := resolveSelfReentryTarget("linux", SelfReentryModeSelf, "")
	if err != nil {
		t.Fatalf("resolveSelfReentryTarget() error = %v", err)
	}
	if got != linuxSelfReentryTarget {
		t.Fatalf("resolveSelfReentryTarget() = %q, want %q", got, linuxSelfReentryTarget)
	}
}

func TestResolveSelfReentryTarget_SelfModeRejectsNonLinux(t *testing.T) {
	_, err := resolveSelfReentryTarget("darwin", SelfReentryModeSelf, "/tmp/current-quine")
	if err == nil || !strings.Contains(err.Error(), "only supported on Linux") {
		t.Fatalf("expected self-mode non-Linux rejection, got %v", err)
	}
}

func TestResolveSelfReentryTarget_ExecutablePathModeUsesExecutablePath(t *testing.T) {
	got, err := resolveSelfReentryTarget("darwin", SelfReentryModeExecutablePath, "/tmp/current-quine")
	if err != nil {
		t.Fatalf("resolveSelfReentryTarget() error = %v", err)
	}
	if got != "/tmp/current-quine" {
		t.Fatalf("resolveSelfReentryTarget() = %q, want /tmp/current-quine", got)
	}
}

func TestResolveSelfReentryTarget_ExecutablePathModeRequiresExecutablePath(t *testing.T) {
	_, err := resolveSelfReentryTarget("darwin", SelfReentryModeExecutablePath, "")
	if err == nil || !strings.Contains(err.Error(), "could not determine current executable path") {
		t.Fatalf("expected missing-target error on darwin, got %v", err)
	}
}

func TestLoadRejectsLegacySelfReentryTargetEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SELF_REENTRY_TARGET", "/tmp/legacy-quine")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_SELF_REENTRY_TARGET has been removed") {
		t.Fatalf("expected legacy self-reentry-target rejection, got %v", err)
	}
}

func TestLoadRejectsDefaultSelfModeOnNonLinuxWhenReentryIsEnabled(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("default self mode is supported on Linux")
	}
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-20250514")
	os.Setenv("QUINE_API_TYPE", "anthropic")
	os.Setenv("QUINE_API_BASE", "https://api.anthropic.com")
	os.Setenv("QUINE_API_KEY", "sk-test-key")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_SELF_REENTRY_MODE=\"self\" is only supported on Linux; use QUINE_SELF_REENTRY_MODE=\"executable_path\" instead") {
		t.Fatalf("expected default self-mode rejection on non-Linux, got %v", err)
	}
}

func containsEnv(env []string, want string) bool {
	for _, e := range env {
		if e == want {
			return true
		}
	}
	return false
}

func writeSkill(t *testing.T, root string, name string, description string) string {
	t.Helper()
	path := filepath.Join(root, ".agents", "skills", name, "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir skill %s: %v", name, err)
	}
	body := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\nThis body should not be injected into the prompt.\n", name, description)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write skill %s: %v", name, err)
	}
	return path
}

func TestHappyPath(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	workspace := root + "/app"
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("realpath root: %v", err)
	}
	workspaceReal, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		t.Fatalf("realpath workspace: %v", err)
	}
	execSessionID := "stable-session-for-resume"
	os.Setenv("QUINE_MAX_DEPTH", "10")
	os.Setenv("QUINE_DEPTH", "3")
	os.Setenv("QUINE_SESSION_ID", execSessionID)
	os.Setenv("QUINE_RUN_ID", "stale-run-id")
	os.Setenv("QUINE_TAPE_ID", "my-tape")
	os.Setenv("QUINE_PARENT_SESSION", "parent-session")
	os.Setenv("QUINE_MAX_CONCURRENT", "50")
	os.Setenv("QUINE_MAX_AGENTS", "25")
	os.Setenv("QUINE_FORK_DEFAULT_TIMEOUT_SECONDS", "45")
	os.Setenv("QUINE_SH_DEFAULT_TIMEOUT_SECONDS", "60")
	os.Setenv("QUINE_SH_TIMEOUT_OVERRIDE_ENABLED", "0")
	os.Setenv("QUINE_SH_STDIN_ENABLED", "0")
	os.Setenv("QUINE_SH_DETACH_ENABLED", "0")
	os.Setenv("QUINE_OUTPUT_TRUNCATE", "4096")
	os.Setenv("QUINE_DATA_DIR", "/tmp/data")
	os.Setenv("QUINE_RETENTION_DIR", "/tmp/retained")
	os.Setenv("QUINE_SHELL", "/bin/sh")
	os.Setenv("QUINE_SH_NETWORK", "none")
	os.Setenv("QUINE_MAX_TURNS", "30")
	os.Setenv("QUINE_WALL_CLOCK_EXIT_SECONDS", "870")
	os.Setenv("QUINE_PROMPT_METAPHOR", "thermodynamic")
	os.Setenv("QUINE_PROMPT_SELF_MODEL", "basic")
	os.Setenv("QUINE_PROMPT_INSTRUCTION_SURFACE", "minimal")
	os.Setenv("QUINE_PROMPT_RUNTIME_SURFACE", "hidden")
	os.Setenv("QUINE_PROMPT_PERSONA", "architect")
	os.Setenv("QUINE_FAIL_ON_IMPOSSIBLE", "0")
	os.Setenv("QUINE_EMPTY_ASSISTANT_SUCCESS", "1")
	os.Setenv("QUINE_READY_TEXT_AUTO_IDLE", "1")
	os.Setenv("QUINE_ANCHOR_MEMORY", "1")
	os.Setenv("QUINE_IDLE_ENABLED", "1")
	os.Setenv("QUINE_MEMORY_WARN_TOKENS", "7000")
	os.Setenv("QUINE_MEMORY_DANGER_TOKENS", "14000")
	os.Setenv("QUINE_MEMORY_DEATH_TOKENS", "18000")
	os.Setenv("QUINE_FORK_ENABLED", "1")
	os.Setenv("QUINE_EXEC_ENABLED", "1")
	os.Setenv("QUINE_SPAWN_ENABLED", "1")
	os.Setenv("QUINE_VISION_ENABLED", "1")
	os.Setenv("QUINE_SH_INTERACTIVE_ENABLED", "1")
	os.Setenv("QUINE_SUPPRESS_INITIAL_BEGIN", "1")
	os.Setenv("QUINE_SELF_SOURCE_CODE_ENABLED", "1")
	os.Setenv("QUINE_WORKSPACE_COMMIT_ON_SIGNAL", "1")
	os.Setenv("QUINE_DEBUG_REQUEST_BODY_DIR", "/tmp/request-dumps")
	if runtime.GOOS == "linux" {
		os.Setenv("QUINE_WORKSPACE_ROOT", root)
		os.Setenv("QUINE_WORKSPACE", workspace)
	}

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"ModelID", c.ModelID, "claude-sonnet-4-20250514"},
		{"APIKey", c.APIKey, "sk-test-key"},
		{"APIBase", c.APIBase, "https://api.anthropic.com"},
		{"Provider", c.Provider, "anthropic"},
		{"DebugRequestBodyDir", c.DebugRequestBodyDir, "/tmp/request-dumps"},
		{"MaxDepth", c.MaxDepth, 10},
		{"Depth", c.Depth, 3},
		{"SessionID", c.SessionID, execSessionID},
		{"TapeID", c.TapeID, "my-tape"},
		{"ParentSession", c.ParentSession, "parent-session"},
		{"MaxConcurrent", c.MaxConcurrent, 50},
		{"MaxAgents", c.MaxAgents, 25},
		{"ForkDefaultTimeoutSeconds", c.ForkDefaultTimeoutSeconds, 45},
		{"ShTimeout", c.ShTimeout, 60},
		{"ShTimeoutOverrideEnabled", c.ShTimeoutOverrideEnabled(), false},
		{"ShStdinEnabled", c.ShStdinEnabled(), false},
		{"ShDetachEnabled", c.ShDetachEnabled(), false},
		{"OutputTruncate", c.OutputTruncate, 4096},
		{"DataDir", c.DataDir, "/tmp/data"},
		{"RetentionDir", c.RetentionDir, "/tmp/retained"},
		{"Shell", c.Shell, "/bin/sh"},
		{"ShNetwork", c.ShNetwork, "none"},
		{"MaxTurns", c.MaxTurns, 30},
		{"WallClockExitSeconds", c.WallClockExitSeconds, 870},
		{"PromptMetaphor", c.PromptMetaphor, PromptMetaphorThermodynamic},
		{"PromptSelfModel", c.PromptSelfModel, PromptSelfModelBasic},
		{"PromptInstructionSurface", c.PromptInstructionSurface, PromptInstructionSurfaceMinimal},
		{"PromptRuntimeSurface", c.PromptRuntimeSurface, PromptRuntimeSurfaceHidden},
		{"PromptPersona", c.PromptPersona, PromptPersonaArchitect},
		{"FailOnImpossible", c.FailOnImpossible, false},
		{"EmptyAssistantSuccess", c.EmptyAssistantSuccess, true},
		{"ReadyTextAutoIdle", c.ReadyTextAutoIdle, true},
		{"AnchorMemoryEnabled", c.AnchorMemoryEnabled, true},
		{"IdleEnabled", c.IdleToolEnabled(), true},
		{"MemoryWarnTokens", c.MemoryWarnTokens, 7000},
		{"MemoryDangerTokens", c.MemoryDangerTokens, 14000},
		{"MemoryDeathTokens", c.MemoryDeathTokens, 18000},
		{"ExecEnabled", c.ExecEnabled, true},
		{"SpawnEnabled", c.SpawnEnabled(), true},
		{"ExitEnabled", c.ExitEnabled(), true},
		{"VisionEnabled", c.VisionEnabled, true},
		{"ShInteractiveEnabled", c.ShInteractiveEnabled(), true},
		{"FSMutationTelemetryEnabled", c.FSMutationTelemetryEnabled(), true},
		{"SuppressInitialBegin", c.SuppressInitialBegin, true},
		{"SelfSourceCodeEnabled", c.SelfSourceCodeEnabled, true},
		{"ForkEnabled", c.ForkEnabled(), true},
		{"VisionEnabled", c.VisionEnabled, true},
		{"WorkspaceCommitOnSignal", c.WorkspaceCommitOnSignal, true},
	}
	if runtime.GOOS == "linux" {
		checks = append(checks,
			struct {
				name string
				got  any
				want any
			}{"WorkspaceEnabled", c.WorkspaceEnabled, true},
			struct {
				name string
				got  any
				want any
			}{"WorkspaceRoot", c.WorkspaceRoot, rootReal},
			struct {
				name string
				got  any
				want any
			}{"Workspace", c.Workspace, workspaceReal},
			struct {
				name string
				got  any
				want any
			}{"WorkspaceRevisionMode", c.WorkspaceRevisionMode, WorkspaceRevisionRestore},
		)
	}
	for _, tc := range checks {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
	if c.RunID == "" || c.RunID == "stale-run-id" {
		t.Fatalf("RunID = %q, want fresh generated physical run id", c.RunID)
	}

	pathChecks := []struct {
		name string
		got  string
		want string
	}{
		{"AgentRoot", c.AgentRoot(), filepath.Join("/tmp/data", "agent", execSessionID)},
		{"AgentStatusDir", c.AgentStatusDir(), filepath.Join("/tmp/data", "agent", execSessionID, "status")},
		{"AgentLogDir", c.AgentLogDir(), filepath.Join("/tmp/data", "agent", execSessionID, "log")},
		{"ControlDir", c.ControlDir(), filepath.Join("/tmp/data", "agent", execSessionID, "ctl")},
		{"ControlPath", c.ControlPath(), filepath.Join("/tmp/data", "agent", execSessionID, "ctl")},
		{"InboxPath", c.InboxPath(), filepath.Join("/tmp/data", "agent", execSessionID, "status", "inbox.json")},
		{"ControlLogPath", c.ControlLogPath(), filepath.Join("/tmp/data", "agent", execSessionID, "log", "control.jsonl")},
	}
	for _, tc := range pathChecks {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

func TestDefaults(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if c.MaxDepth != 0 {
		t.Errorf("MaxDepth = %d, want 0", c.MaxDepth)
	}
	if c.Depth != 0 {
		t.Errorf("Depth = %d, want 0", c.Depth)
	}
	if c.MaxConcurrent != 0 {
		t.Errorf("MaxConcurrent = %d, want 0", c.MaxConcurrent)
	}
	if c.MaxAgents != 0 {
		t.Errorf("MaxAgents = %d, want 0", c.MaxAgents)
	}
	if c.ForkDefaultTimeoutSeconds != 0 {
		t.Errorf("ForkDefaultTimeoutSeconds = %d, want 0", c.ForkDefaultTimeoutSeconds)
	}
	if c.ShTimeout != 300 {
		t.Errorf("ShTimeout = %d, want 300", c.ShTimeout)
	}
	if !c.ShTimeoutOverrideEnabled() {
		t.Error("ShTimeoutOverrideEnabled() should default true")
	}
	if !c.ShStdinEnabled() {
		t.Error("ShStdinEnabled() should default true")
	}
	if !c.ShDetachEnabled() {
		t.Error("ShDetachEnabled() should default true")
	}
	if c.OutputTruncate != 20480 {
		t.Errorf("OutputTruncate = %d, want 20480", c.OutputTruncate)
	}
	if c.DataDir != ".quine/" {
		t.Errorf("DataDir = %q, want %q", c.DataDir, ".quine/")
	}
	if c.RetentionDir != "" {
		t.Errorf("RetentionDir = %q, want empty string", c.RetentionDir)
	}
	if c.Shell != "/bin/sh" {
		t.Errorf("Shell = %q, want /bin/sh", c.Shell)
	}
	if c.ShNetwork != "host" {
		t.Errorf("ShNetwork = %q, want host", c.ShNetwork)
	}
	if c.MaxTurns != 0 {
		t.Errorf("MaxTurns = %d, want 0 (disabled)", c.MaxTurns)
	}
	if c.PromptMetaphor != PromptMetaphorOff {
		t.Errorf("PromptMetaphor = %q, want %q", c.PromptMetaphor, PromptMetaphorOff)
	}
	if c.PromptSelfModel != PromptSelfModelAdvanced {
		t.Errorf("PromptSelfModel = %q, want %q", c.PromptSelfModel, PromptSelfModelAdvanced)
	}
	if c.PromptInstructionSurface != PromptInstructionSurfaceStandard {
		t.Errorf("PromptInstructionSurface = %q, want %q", c.PromptInstructionSurface, PromptInstructionSurfaceStandard)
	}
	if c.PromptRuntimeSurface != PromptRuntimeSurfaceVisible {
		t.Errorf("PromptRuntimeSurface = %q, want %q", c.PromptRuntimeSurface, PromptRuntimeSurfaceVisible)
	}
	if c.PromptPersona != PromptPersonaNone {
		t.Errorf("PromptPersona = %q, want empty", c.PromptPersona)
	}
	if !c.FailOnImpossible {
		t.Errorf("FailOnImpossible = false, want true by default")
	}
	if c.ContextWindow != 128_000 {
		t.Errorf("ContextWindow = %d, want 128000", c.ContextWindow)
	}
	if c.MemoryWarnTokens != 8000 {
		t.Errorf("MemoryWarnTokens = %d, want 8000", c.MemoryWarnTokens)
	}
	if c.MemoryDangerTokens != 16000 {
		t.Errorf("MemoryDangerTokens = %d, want 16000", c.MemoryDangerTokens)
	}
	if c.MemoryDeathTokens != 0 {
		t.Errorf("MemoryDeathTokens = %d, want 0", c.MemoryDeathTokens)
	}
	if c.AnchorMemoryEnabled {
		t.Errorf("AnchorMemoryEnabled = true, want false by default")
	}
	if c.IdleToolEnabled() {
		t.Errorf("IdleEnabled = true, want false by default")
	}
	if !c.ExecEnabled {
		t.Errorf("ExecEnabled = false, want true by default")
	}
	if !c.ExitEnabled() {
		t.Errorf("ExitEnabled = false, want true by default")
	}
	if !c.ForkEnabled() {
		t.Errorf("ForkEnabled = false, want true by default")
	}
	if !c.VisionEnabled {
		t.Errorf("VisionEnabled = false, want true by default")
	}
	if !c.ShInteractiveEnabled() {
		t.Errorf("ShInteractiveEnabled = false, want true by default")
	}
	if !c.FSMutationTelemetryEnabled() {
		t.Errorf("FSMutationTelemetryEnabled = false, want true by default")
	}
	if c.PeerDiscoveryEnabled {
		t.Errorf("PeerDiscoveryEnabled = true, want false by default")
	}
	if c.ReadyTextAutoIdle {
		t.Errorf("ReadyTextAutoIdle = true, want false by default")
	}
	if c.PeerDiscoveryHeartbeatMS != 5000 {
		t.Errorf("PeerDiscoveryHeartbeatMS = %d, want 5000", c.PeerDiscoveryHeartbeatMS)
	}
	if c.SelfSourceCodeEnabled {
		t.Errorf("SelfSourceCodeEnabled = true, want false by default")
	}
	if c.SuppressInitialBegin {
		t.Errorf("SuppressInitialBegin = true, want false by default")
	}
	if c.WorkspaceEnabled {
		t.Errorf("WorkspaceEnabled = true, want false by default")
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionNone {
		t.Errorf("WorkspaceRevisionMode = %q, want %q", c.WorkspaceRevisionMode, WorkspaceRevisionNone)
	}
	if c.WorkspaceRoot != "" || c.Workspace != "" {
		t.Errorf("workspace fields should default empty, got root=%q workspace=%q", c.WorkspaceRoot, c.Workspace)
	}
	if c.SessionID == "" {
		t.Error("SessionID should be auto-generated, got empty")
	}
	if !strings.HasPrefix(c.SessionID, "sess_") {
		t.Fatalf("SessionID = %q, want sess_<YYYYMMDD-HHMMSS>_<random>", c.SessionID)
	}
	sessionParts := strings.SplitN(strings.TrimPrefix(c.SessionID, "sess_"), "_", 2)
	if len(sessionParts) != 2 {
		t.Fatalf("SessionID = %q, want sess_<YYYYMMDD-HHMMSS>_<random>", c.SessionID)
	}
	if _, err := time.Parse("20060102-150405", sessionParts[0]); err != nil {
		t.Errorf("SessionID start_time = %q, want YYYYMMDD-HHMMSS: %v", sessionParts[0], err)
	}
	if strings.Contains(c.SessionID, "_"+strconv.Itoa(os.Getpid())) {
		t.Errorf("SessionID = %q should not embed current pid", c.SessionID)
	}
	if c.RunID == "" {
		t.Error("RunID should be auto-generated, got empty")
	}
	if !strings.HasPrefix(c.RunID, "run_") || !strings.Contains(c.RunID, "_"+strconv.Itoa(os.Getpid())+"_") {
		t.Fatalf("RunID = %q, want run_<YYYYMMDD-HHMMSS>_<pid>_<random>", c.RunID)
	}
	if c.TapeID == "" {
		t.Error("TapeID should be auto-generated, got empty")
	}
	if c.TapeID != "0001" {
		t.Errorf("TapeID = %q, want %q", c.TapeID, "0001")
	}
	if c.TapeID == c.SessionID {
		t.Errorf("TapeID = %q, want a distinct trace id from SessionID", c.TapeID)
	}
}

func TestExplicitDataDirWithoutRetentionDirUsesSessionLogFallback(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_DATA_DIR", "/tmp/runtime")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.DataDir != "/tmp/runtime" {
		t.Fatalf("DataDir = %q, want %q", c.DataDir, "/tmp/runtime")
	}
	if c.RetentionDir != "" {
		t.Fatalf("RetentionDir = %q, want empty fallback mode", c.RetentionDir)
	}
	if got := c.SessionRetainedDir(""); got != filepath.Join("/tmp/runtime", "log", c.SessionID) {
		t.Fatalf("SessionRetainedDir() = %q, want fallback under QUINE_DATA_DIR/log/<session>", got)
	}
}

func TestExplicitRetentionDirWithoutDataDirUsesDefaultRuntimeRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_RETENTION_DIR", "/tmp/retained")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.DataDir != ".quine/" {
		t.Fatalf("DataDir = %q, want default runtime root", c.DataDir)
	}
	if c.RetentionDir != "/tmp/retained" {
		t.Fatalf("RetentionDir = %q, want explicit retained root", c.RetentionDir)
	}
}

func TestSessionRetainedDirAndDerivedPaths(t *testing.T) {
	t.Run("without retention dir", func(t *testing.T) {
		cfg := &Config{
			Paths:    Paths{DataDir: "/tmp/runtime"},
			Identity: Identity{SessionID: "session-1", TapeID: "0001"},
		}

		if got := cfg.SessionRetainedDir(""); got != "/tmp/runtime/log/session-1" {
			t.Fatalf("SessionRetainedDir() = %q, want %q", got, "/tmp/runtime/log/session-1")
		}
		if got := cfg.SessionLogDir(""); got != "/tmp/runtime/log/session-1" {
			t.Fatalf("SessionLogDir() = %q, want %q", got, "/tmp/runtime/log/session-1")
		}
		if got := cfg.TapeDir(""); got != "/tmp/runtime/log/session-1/tapes" {
			t.Fatalf("TapeDir() = %q, want %q", got, "/tmp/runtime/log/session-1/tapes")
		}
		if got := cfg.SessionIncarnationDir(""); got != "/tmp/runtime/log/session-1/inc" {
			t.Fatalf("SessionIncarnationDir() = %q, want %q", got, "/tmp/runtime/log/session-1/inc")
		}
		if got := cfg.SessionControlLogPath(""); got != "/tmp/runtime/log/session-1/control.jsonl" {
			t.Fatalf("SessionControlLogPath() = %q, want %q", got, "/tmp/runtime/log/session-1/control.jsonl")
		}
	})

	t.Run("with retention dir", func(t *testing.T) {
		cfg := &Config{
			Paths:    Paths{DataDir: "/tmp/runtime", RetentionDir: "/tmp/retained"},
			Identity: Identity{SessionID: "session-1", TapeID: "0001"},
		}

		if got := cfg.SessionRetainedDir(""); got != "/tmp/retained/sessions/session-1" {
			t.Fatalf("SessionRetainedDir() = %q, want %q", got, "/tmp/retained/sessions/session-1")
		}
		if got := cfg.SessionLogDir(""); got != "/tmp/retained/sessions/session-1" {
			t.Fatalf("SessionLogDir() = %q, want %q", got, "/tmp/retained/sessions/session-1")
		}
		if got := cfg.TapeDir(""); got != "/tmp/retained/sessions/session-1/tapes" {
			t.Fatalf("TapeDir() = %q, want %q", got, "/tmp/retained/sessions/session-1/tapes")
		}
		if got := cfg.SessionIncarnationDir(""); got != "/tmp/retained/sessions/session-1/inc" {
			t.Fatalf("SessionIncarnationDir() = %q, want %q", got, "/tmp/retained/sessions/session-1/inc")
		}
		if got := cfg.SessionControlLogPath(""); got != "/tmp/retained/sessions/session-1/control.jsonl" {
			t.Fatalf("SessionControlLogPath() = %q, want %q", got, "/tmp/retained/sessions/session-1/control.jsonl")
		}
	})
}

func TestShTimeoutOverrideCanBeDisabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SH_TIMEOUT_OVERRIDE_ENABLED", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ShTimeoutOverrideEnabled() {
		t.Fatal("ShTimeoutOverrideEnabled() should be false when QUINE_SH_TIMEOUT_OVERRIDE_ENABLED=0")
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_SH_TIMEOUT_OVERRIDE_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_SH_TIMEOUT_OVERRIDE_ENABLED=0, got %v", env)
	}
}

func TestShNetworkParsesAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SH_NETWORK", "none")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ShNetwork != "none" {
		t.Fatalf("ShNetwork = %q, want none", c.ShNetwork)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_SH_NETWORK=none") {
		t.Fatalf("ChildEnv() should propagate QUINE_SH_NETWORK=none, got %v", env)
	}
}

func TestShNetworkRejectsUnsupportedValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SH_NETWORK", "offline")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_SH_NETWORK") {
		t.Fatalf("Load() error = %v, want QUINE_SH_NETWORK validation error", err)
	}
}

func TestForkDefaultTimeoutSecondsParsesAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_FORK_DEFAULT_TIMEOUT_SECONDS", "17")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ForkDefaultTimeoutSeconds != 17 {
		t.Fatalf("ForkDefaultTimeoutSeconds = %d, want 17", c.ForkDefaultTimeoutSeconds)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_FORK_DEFAULT_TIMEOUT_SECONDS=17") {
		t.Fatalf("ChildEnv() should propagate QUINE_FORK_DEFAULT_TIMEOUT_SECONDS=17, got %v", env)
	}
}

func TestSpawnEnabledParsesAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.SpawnEnabled() {
		t.Fatal("SpawnEnabled() should default to false")
	}

	os.Setenv("QUINE_SPAWN_ENABLED", "1")
	c, err = Load()
	if err != nil {
		t.Fatalf("Load() with spawn enabled error: %v", err)
	}
	if !c.SpawnEnabled() {
		t.Fatal("SpawnEnabled() should be true when QUINE_SPAWN_ENABLED=1")
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_SPAWN_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_SPAWN_ENABLED=1, got %v", env)
	}

	os.Setenv("QUINE_FORK_ENABLED", "0")
	os.Setenv("QUINE_EXEC_ENABLED", "0")
	c, err = Load()
	if err != nil {
		t.Fatalf("Load() with only spawn enabled error: %v", err)
	}
	if !c.SpawnEnabled() || c.ForkEnabled() || c.ExecEnabled {
		t.Fatalf("expected only spawn enabled, got spawn=%v fork=%v exec=%v", c.SpawnEnabled(), c.ForkEnabled(), c.ExecEnabled)
	}
	if c.SelfReentryTarget != expectedTestSelfReentryTarget(t) {
		t.Fatalf("SelfReentryTarget = %q, want %q when spawn alone needs a quine binary", c.SelfReentryTarget, expectedTestSelfReentryTarget(t))
	}
}

func TestWorkspaceCommitOnSignalParsesAndPropagates(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics require Linux")
	}
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_COMMIT_ON_SIGNAL", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.WorkspaceCommitOnSignal {
		t.Fatal("WorkspaceCommitOnSignal = false, want true")
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_WORKSPACE_COMMIT_ON_SIGNAL=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_WORKSPACE_COMMIT_ON_SIGNAL=1, got %v", env)
	}
}

func TestEmptyAssistantSuccessParsesAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_EMPTY_ASSISTANT_SUCCESS", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.EmptyAssistantSuccess {
		t.Fatal("EmptyAssistantSuccess = false, want true")
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_EMPTY_ASSISTANT_SUCCESS=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_EMPTY_ASSISTANT_SUCCESS=1, got %v", env)
	}
}

func TestReadyTextAutoIdleParsesAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_READY_TEXT_AUTO_IDLE", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.ReadyTextAutoIdle {
		t.Fatal("ReadyTextAutoIdle = false, want true")
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_READY_TEXT_AUTO_IDLE=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_READY_TEXT_AUTO_IDLE=1, got %v", env)
	}
}

func TestWallClockExitSecondsParsesAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_WALL_CLOCK_EXIT_SECONDS", "870")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.WallClockExitSeconds != 870 {
		t.Fatalf("WallClockExitSeconds = %d, want 870", c.WallClockExitSeconds)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_WALL_CLOCK_EXIT_SECONDS=870") {
		t.Fatalf("ChildEnv() should propagate QUINE_WALL_CLOCK_EXIT_SECONDS=870, got %v", env)
	}
}

func TestForkDefaultTimeoutSecondsRejectsNegative(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_FORK_DEFAULT_TIMEOUT_SECONDS", "-1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_FORK_DEFAULT_TIMEOUT_SECONDS must be >= 0") {
		t.Fatalf("Load() error = %v, want fork timeout validation error", err)
	}
}

func TestWallClockExitSecondsRejectsNegative(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_WALL_CLOCK_EXIT_SECONDS", "-1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_WALL_CLOCK_EXIT_SECONDS must be >= 0") {
		t.Fatalf("Load() error = %v, want wall-clock exit validation error", err)
	}
}

func TestBooleanEnvRejectsNonBinaryValues(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
	}{
		{name: "tool gate", key: "QUINE_VISION_ENABLED", val: "true"},
		{name: "runtime acceptance", key: "QUINE_FAIL_ON_IMPOSSIBLE", val: "yes"},
		{name: "empty assistant success", key: "QUINE_EMPTY_ASSISTANT_SUCCESS", val: "true"},
		{name: "ready text auto idle", key: "QUINE_READY_TEXT_AUTO_IDLE", val: "true"},
		{name: "workspace ownership", key: "QUINE_WORKSPACE_OWNER", val: "false"},
		{name: "signal workspace commit", key: "QUINE_WORKSPACE_COMMIT_ON_SIGNAL", val: "true"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			setRequired(t)
			os.Setenv(tc.key, tc.val)

			_, err := Load()
			if err == nil {
				t.Fatalf("expected binary-bool rejection for %s, got %v", tc.key, err)
			}
			if !strings.Contains(err.Error(), tc.key+"=") {
				t.Fatalf("expected env key in error for %s, got %v", tc.key, err)
			}
			if !strings.Contains(err.Error(), `must be "0" or "1"`) {
				t.Fatalf("expected 0/1 guidance for %s, got %v", tc.key, err)
			}
		})
	}
}

func TestPositiveToolGateInversionsDefaultAndPropagate(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		enabled func(*Config) bool
	}{
		{name: "sh timeout override", key: EnvShTimeoutOverride, enabled: (*Config).ShTimeoutOverrideEnabled},
		{name: "sh stdin", key: EnvShStdinEnabled, enabled: (*Config).ShStdinEnabled},
		{name: "sh detach", key: EnvShDetachEnabled, enabled: (*Config).ShDetachEnabled},
		{name: "exit", key: EnvExitEnabled, enabled: (*Config).ExitEnabled},
		{name: "fork", key: EnvForkEnabled, enabled: (*Config).ForkEnabled},
		{name: "sh interactive", key: EnvShInteractiveEnabled, enabled: (*Config).ShInteractiveEnabled},
		{name: "fs mutation telemetry", key: EnvFSMutationTelemetry, enabled: (*Config).FSMutationTelemetryEnabled},
	}

	for _, tc := range tests {
		t.Run(tc.name+" default", func(t *testing.T) {
			clearEnv(t)
			setRequired(t)

			c, err := Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}
			if !tc.enabled(c) {
				t.Fatalf("%s should default enabled", tc.key)
			}

			childEnv, err := c.ChildEnv()
			if err != nil {
				t.Fatalf("ChildEnv() error: %v", err)
			}
			if !containsEnv(childEnv, envKV(tc.key, "1")) {
				t.Fatalf("ChildEnv() should propagate %s=1, got %v", tc.key, childEnv)
			}

			execEnv, err := c.ExecEnv()
			if err != nil {
				t.Fatalf("ExecEnv() error: %v", err)
			}
			if !containsEnv(execEnv, envKV(tc.key, "1")) {
				t.Fatalf("ExecEnv() should propagate %s=1, got %v", tc.key, execEnv)
			}
		})

		t.Run(tc.name+" disabled", func(t *testing.T) {
			clearEnv(t)
			setRequired(t)
			os.Setenv(tc.key, "0")

			c, err := Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
			}
			if tc.enabled(c) {
				t.Fatalf("%s should be disabled when set to 0", tc.key)
			}

			childEnv, err := c.ChildEnv()
			if err != nil {
				t.Fatalf("ChildEnv() error: %v", err)
			}
			if !containsEnv(childEnv, envKV(tc.key, "0")) {
				t.Fatalf("ChildEnv() should propagate %s=0, got %v", tc.key, childEnv)
			}

			execEnv, err := c.ExecEnv()
			if err != nil {
				t.Fatalf("ExecEnv() error: %v", err)
			}
			if !containsEnv(execEnv, envKV(tc.key, "0")) {
				t.Fatalf("ExecEnv() should propagate %s=0, got %v", tc.key, execEnv)
			}
		})
	}
}

func TestForkEnabledCanBeDisabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_FORK_ENABLED", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ForkEnabled() {
		t.Fatal("ForkEnabled() should be false when QUINE_FORK_ENABLED=0")
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_FORK_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_FORK_ENABLED=0, got %v", env)
	}
}

func TestAnchorMarkEnabled_DefaultsTrue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.AnchorMarkEnabled() {
		t.Fatal("AnchorMarkEnabled() should default to true")
	}
}

func TestAnchorMarkCanBeDisabledAndPropagatesUnderAnchorMemory(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_ANCHOR_MEMORY", "1")
	os.Setenv("QUINE_ANCHOR_MARK_ENABLED", "0")
	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.AnchorMarkEnabled() {
		t.Fatal("AnchorMarkEnabled() should be false when QUINE_ANCHOR_MARK_ENABLED=0")
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_ANCHOR_MARK_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_ANCHOR_MARK_ENABLED=0 under anchor memory, got %v", env)
	}
}

func TestAnchorFoldEnabled_DefaultsTrue(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.AnchorFoldEnabled() {
		t.Fatal("AnchorFoldEnabled() should default to true")
	}
}

func TestAnchorFoldCanBeDisabledAndPropagatesUnderAnchorMemory(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_ANCHOR_MEMORY", "1")
	os.Setenv("QUINE_ANCHOR_FOLD_ENABLED", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.AnchorFoldEnabled() {
		t.Fatal("AnchorFoldEnabled() should be false when QUINE_ANCHOR_FOLD_ENABLED=0")
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_ANCHOR_FOLD_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_ANCHOR_FOLD_ENABLED=0 under anchor memory, got %v", env)
	}
}

func TestShInteractiveCanBeDisabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SH_INTERACTIVE_ENABLED", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ShInteractiveEnabled() {
		t.Fatal("ShInteractiveEnabled() should be false when QUINE_SH_INTERACTIVE_ENABLED=0")
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_SH_INTERACTIVE_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_SH_INTERACTIVE_ENABLED=0, got %v", env)
	}
}

func TestFSMutationTelemetryCanBeDisabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_FS_MUTATION_TELEMETRY_ENABLED", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.FSMutationTelemetryEnabled() {
		t.Fatal("FSMutationTelemetryEnabled() should be false when QUINE_FS_MUTATION_TELEMETRY_ENABLED=0")
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_FS_MUTATION_TELEMETRY_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_FS_MUTATION_TELEMETRY_ENABLED=0, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_FS_MUTATION_TELEMETRY_ENABLED=0") {
		t.Fatalf("ExecEnv() should propagate QUINE_FS_MUTATION_TELEMETRY_ENABLED=0, got %v", execEnv)
	}
}

func TestPromptCtlCanBeHiddenAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_CTL", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.PromptCtlPhysics {
		t.Fatal("PromptCtlPhysics should be false when QUINE_PROMPT_CTL=0")
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_PROMPT_CTL=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_PROMPT_CTL=0, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_PROMPT_CTL=0") {
		t.Fatalf("ExecEnv() should propagate QUINE_PROMPT_CTL=0, got %v", execEnv)
	}
}

func TestPeerDiscoveryEnabledCanBeEnabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PEER_DISCOVERY_ENABLED", "1")
	os.Setenv("QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS", "2000")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.PeerDiscoveryEnabled {
		t.Fatal("PeerDiscoveryEnabled should be true when QUINE_PEER_DISCOVERY_ENABLED=1")
	}
	if c.PeerDiscoveryHeartbeatMS != 2000 {
		t.Fatalf("PeerDiscoveryHeartbeatMS = %d, want 2000", c.PeerDiscoveryHeartbeatMS)
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_PEER_DISCOVERY_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_PEER_DISCOVERY_ENABLED=1, got %v", childEnv)
	}
	if !containsEnv(childEnv, "QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS=2000") {
		t.Fatalf("ChildEnv() should propagate QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS=2000, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_PEER_DISCOVERY_ENABLED=1") {
		t.Fatalf("ExecEnv() should propagate QUINE_PEER_DISCOVERY_ENABLED=1, got %v", execEnv)
	}
	if !containsEnv(execEnv, "QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS=2000") {
		t.Fatalf("ExecEnv() should propagate QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS=2000, got %v", execEnv)
	}
}

func TestPromptImplDetailsCanBeEnabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_IMPL_DETAILS", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.PromptImplDetails {
		t.Fatal("PromptImplDetails should be true when QUINE_PROMPT_IMPL_DETAILS=1")
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_PROMPT_IMPL_DETAILS=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_PROMPT_IMPL_DETAILS=1, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_PROMPT_IMPL_DETAILS=1") {
		t.Fatalf("ExecEnv() should propagate QUINE_PROMPT_IMPL_DETAILS=1, got %v", execEnv)
	}
}

func TestMemoryStrategyHintsCanBeDisabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MEMORY_STRATEGY_HINTS", "0")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.MemoryStrategyHints {
		t.Fatal("MemoryStrategyHints should be false when QUINE_MEMORY_STRATEGY_HINTS=0")
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_MEMORY_STRATEGY_HINTS=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_MEMORY_STRATEGY_HINTS=0, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_MEMORY_STRATEGY_HINTS=0") {
		t.Fatalf("ExecEnv() should propagate QUINE_MEMORY_STRATEGY_HINTS=0, got %v", execEnv)
	}
}

func TestAgentsMDDisabledByDefaultAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.AgentsMDEnabled {
		t.Fatal("AgentsMDEnabled should default false")
	}
	if c.AgentsMDPath != "" {
		t.Fatalf("AgentsMDPath = %q, want empty", c.AgentsMDPath)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_AGENTS_MD_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_AGENTS_MD_ENABLED=0, got %v", env)
	}
}

func TestAgentsMDEnabledFindsSingleFileAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	work := filepath.Join(root, "pkg")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatalf("mkdir work: %v", err)
	}
	agentsPath := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("repo instructions\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", work)
	os.Setenv("QUINE_AGENTS_MD_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.AgentsMDEnabled {
		t.Fatal("AgentsMDEnabled should be true when QUINE_AGENTS_MD_ENABLED=1")
	}
	if c.AgentsMDPath != agentsPath {
		t.Fatalf("AgentsMDPath = %q, want %q", c.AgentsMDPath, agentsPath)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_AGENTS_MD_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_AGENTS_MD_ENABLED=1, got %v", env)
	}
}

func TestAgentsMDEnabledRejectsInvalidGateValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_AGENTS_MD_ENABLED", "true")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_AGENTS_MD_ENABLED") {
		t.Fatalf("expected invalid AGENTS gate error, got %v", err)
	}
}

func TestAgentsMDEnabledRejectsMultipleFiles(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	subdir := filepath.Join(root, "pkg")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("root\n"), 0o644); err != nil {
		t.Fatalf("write root AGENTS.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "AGENTS.md"), []byte("nested\n"), 0o644); err != nil {
		t.Fatalf("write nested AGENTS.md: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", subdir)
	os.Setenv("QUINE_AGENTS_MD_ENABLED", "1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "hierarchical AGENTS.md is not supported yet") {
		t.Fatalf("expected multiple AGENTS.md error, got %v", err)
	}
}

func TestAgentSkillsDisabledByDefaultAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.AgentsSkillsEnabled {
		t.Fatal("AgentsSkillsEnabled should default false")
	}
	if len(c.Skills) != 0 {
		t.Fatalf("Skills length = %d, want 0", len(c.Skills))
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_AGENTS_SKILLS_ENABLED=0") {
		t.Fatalf("ChildEnv() should propagate QUINE_AGENTS_SKILLS_ENABLED=0, got %v", env)
	}
}

func TestAgentSkillsEnabledLoadsIndexAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	work := filepath.Join(root, "pkg")
	writeSkill(t, root, "bar", "Bar helper")
	writeSkill(t, root, "foo", "Foo helper")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatalf("mkdir work: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", work)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.AgentsSkillsEnabled {
		t.Fatal("AgentsSkillsEnabled should be true when QUINE_AGENTS_SKILLS_ENABLED=1")
	}
	if len(c.Skills) != 2 {
		t.Fatalf("Skills length = %d, want 2: %#v", len(c.Skills), c.Skills)
	}
	if c.Skills[0].Name != "bar" || c.Skills[0].Description != "Bar helper" || c.Skills[0].Source != "../.agents/skills/bar/SKILL.md" {
		t.Fatalf("first skill = %#v", c.Skills[0])
	}
	if c.Skills[1].Name != "foo" || c.Skills[1].Description != "Foo helper" || c.Skills[1].Source != "../.agents/skills/foo/SKILL.md" {
		t.Fatalf("second skill = %#v", c.Skills[1])
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_AGENTS_SKILLS_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_AGENTS_SKILLS_ENABLED=1, got %v", env)
	}
}

func TestAgentSkillsEnabledRejectsInvalidGateValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "true")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_AGENTS_SKILLS_ENABLED") {
		t.Fatalf("expected invalid skills gate error, got %v", err)
	}
}

func TestAgentSkillsEnabledUsesFrontmatterNameAsCanonicalName(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	path := filepath.Join(root, ".agents", "skills", "gitbutler", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir skill: %v", err)
	}
	if err := os.WriteFile(path, []byte("---\nname: but\ndescription: GitButler command skill\n---\nbody\n"), 0o644); err != nil {
		t.Fatalf("write skill: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(c.Skills) != 1 {
		t.Fatalf("Skills length = %d, want 1", len(c.Skills))
	}
	got := c.Skills[0]
	if got.Name != "but" || got.Description != "GitButler command skill" || got.Source != ".agents/skills/gitbutler/SKILL.md" {
		t.Fatalf("skill index = %#v", got)
	}
}

func TestAgentSkillsEnabledRejectsDuplicateFrontmatterName(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	for _, dir := range []string{"gitbutler", "but-alias"} {
		path := filepath.Join(root, ".agents", "skills", dir, "SKILL.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir skill %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("---\nname: but\ndescription: Duplicate command skill\n---\nbody\n"), 0o644); err != nil {
			t.Fatalf("write skill %s: %v", dir, err)
		}
	}
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "duplicate skill name \"but\"") {
		t.Fatalf("expected duplicate skill name error, got %v", err)
	}
}

func TestAgentSkillsEnabledRejectsMissingSkillFile(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "empty-package"), 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "missing SKILL.md") {
		t.Fatalf("expected missing SKILL.md error, got %v", err)
	}
}

func TestAgentSkillsEnabledAllowsHierarchicalSkillDirectoryButIndexesOnlyMetadata(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	writeSkill(t, root, "review-helper", "Review helper")
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "review-helper", "scripts"), 0o755); err != nil {
		t.Fatalf("mkdir scripts: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "review-helper", "references"), 0o755); err != nil {
		t.Fatalf("mkdir references: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "review-helper", "assets"), 0o755); err != nil {
		t.Fatalf("mkdir assets: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".agents", "skills", "review-helper", "scripts", "check.sh"), []byte("echo SCRIPT_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".agents", "skills", "review-helper", "references", "rules.md"), []byte("REFERENCE_MARKER\n"), 0o644); err != nil {
		t.Fatalf("write reference: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(c.Skills) != 1 {
		t.Fatalf("Skills length = %d, want 1", len(c.Skills))
	}
	got := c.Skills[0]
	if got.Name != "review-helper" || got.Description != "Review helper" || got.Source != ".agents/skills/review-helper/SKILL.md" {
		t.Fatalf("skill index = %#v", got)
	}
}

func TestAgentSkillsEnabledRejectsMissingSkillDescription(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	path := filepath.Join(root, ".agents", "skills", "missing-description", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir skill: %v", err)
	}
	if err := os.WriteFile(path, []byte("---\nname: missing-description\n---\nbody\n"), 0o644); err != nil {
		t.Fatalf("write skill: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "missing frontmatter field description") {
		t.Fatalf("expected missing skill description error, got %v", err)
	}
}

func TestAgentSkillsEnabledRejectsInvalidSkillName(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	path := filepath.Join(root, ".agents", "skills", "bad--name", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir skill: %v", err)
	}
	if err := os.WriteFile(path, []byte("---\nname: bad--name\ndescription: Invalid name\n---\nbody\n"), 0o644); err != nil {
		t.Fatalf("write skill: %v", err)
	}
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "invalid name") {
		t.Fatalf("expected invalid skill name error, got %v", err)
	}
}

func TestAgentSkillsIndexIsStartupSnapshot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	skillPath := writeSkill(t, root, "foo", "Description v1")
	os.Setenv("QUINE_WORK_DIR", root)
	os.Setenv("QUINE_AGENTS_SKILLS_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got := c.Skills[0].Description; got != "Description v1" {
		t.Fatalf("initial description = %q", got)
	}
	if err := os.WriteFile(skillPath, []byte("---\nname: foo\ndescription: Description v2\n---\nbody\n"), 0o644); err != nil {
		t.Fatalf("rewrite skill: %v", err)
	}
	if got := c.Skills[0].Description; got != "Description v1" {
		t.Fatalf("loaded config should retain startup snapshot, got %q", got)
	}

	c2, err := Load()
	if err != nil {
		t.Fatalf("Load() second error: %v", err)
	}
	if got := c2.Skills[0].Description; got != "Description v2" {
		t.Fatalf("reloaded description = %q, want Description v2", got)
	}
}

func TestIdleEnabledCanBeEnabledAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_IDLE_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.IdleToolEnabled() {
		t.Fatal("IdleToolEnabled() should be true when QUINE_IDLE_ENABLED=1")
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_IDLE_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_IDLE_ENABLED=1, got %v", env)
	}
}

func TestTapeIDAutoIncrementsPerSession(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	dataDir := t.TempDir()
	retainedRoot := t.TempDir()
	// Session ID must embed current PID to be inherited
	sessionID := fmt.Sprintf("20260312-120000_%d_%d", os.Getppid(), os.Getpid())
	tapeDir := filepath.Join(retainedRoot, "sessions", sessionID, "tapes")
	if err := os.MkdirAll(tapeDir, 0o755); err != nil {
		t.Fatalf("mkdir tape dir: %v", err)
	}
	for _, name := range []string{"0001.jsonl", "0001.log.yaml", "0002.jsonl", "ignore.txt"} {
		if err := os.WriteFile(filepath.Join(tapeDir, name), []byte("x"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	os.Setenv("QUINE_DATA_DIR", dataDir)
	os.Setenv("QUINE_RETENTION_DIR", retainedRoot)
	os.Setenv("QUINE_SESSION_ID", sessionID)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if c.TapeID != "0003" {
		t.Errorf("TapeID = %q, want %q", c.TapeID, "0003")
	}
}

// --- Required field validation tests ---

func TestMissingRequiredField(t *testing.T) {
	tests := []struct {
		name    string
		missing string
		setVars map[string]string
	}{
		{
			name:    "ModelID",
			missing: "QUINE_MODEL_ID",
			setVars: map[string]string{
				"QUINE_API_TYPE": "openai",
				"QUINE_API_BASE": "https://api.openai.com",
				"QUINE_API_KEY":  "sk-test",
			},
		},
		{
			name:    "APIType",
			missing: "QUINE_API_TYPE",
			setVars: map[string]string{
				"QUINE_MODEL_ID": "some-model",
				"QUINE_API_BASE": "https://example.com",
				"QUINE_API_KEY":  "sk-test",
			},
		},
		{
			name:    "APIBase",
			missing: "QUINE_API_BASE",
			setVars: map[string]string{
				"QUINE_MODEL_ID": "some-model",
				"QUINE_API_TYPE": "openai",
				"QUINE_API_KEY":  "sk-test",
			},
		},
		{
			name:    "APIKey",
			missing: "QUINE_API_KEY",
			setVars: map[string]string{
				"QUINE_MODEL_ID": "some-model",
				"QUINE_API_TYPE": "openai",
				"QUINE_API_BASE": "https://example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.setVars {
				os.Setenv(k, v)
			}
			_, err := Load()
			if err == nil {
				t.Fatalf("expected error for missing %s", tt.missing)
			}
			if !strings.Contains(err.Error(), tt.missing) {
				t.Errorf("error should mention %s, got: %v", tt.missing, err)
			}
		})
	}
}

func TestUnsupportedAPIType(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "some-model")
	os.Setenv("QUINE_API_TYPE", "gemini")
	os.Setenv("QUINE_API_BASE", "https://example.com")
	os.Setenv("QUINE_API_KEY", "sk-test")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for unsupported API type")
	}
	if !strings.Contains(err.Error(), "unsupported QUINE_API_TYPE") {
		t.Errorf("error should mention unsupported QUINE_API_TYPE, got: %v", err)
	}
}

func TestCodexOAuthAPIKeyAllowed(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "gpt-5.2-codex")
	os.Setenv("QUINE_API_TYPE", "openai-responses")
	os.Setenv("QUINE_API_BASE", "https://api.openai.com")
	os.Setenv("QUINE_API_KEY", "codex-oauth") // Special sentinel value triggers OAuth
	if runtime.GOOS != "linux" {
		os.Setenv("QUINE_SELF_REENTRY_MODE", string(SelfReentryModeExecutablePath))
	}

	_, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
}

func TestClaudeOAuthAPIKeyAllowed(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4-6")
	os.Setenv("QUINE_API_TYPE", "anthropic")
	os.Setenv("QUINE_API_BASE", "https://api.anthropic.com")
	os.Setenv("QUINE_API_KEY", "claude-oauth")
	if runtime.GOOS != "linux" {
		os.Setenv("QUINE_SELF_REENTRY_MODE", string(SelfReentryModeExecutablePath))
	}

	_, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
}

func TestGitHubCopilotOAuthAPIKeyAllowed(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "claude-sonnet-4.6")
	os.Setenv("QUINE_API_TYPE", "openai")
	os.Setenv("QUINE_API_BASE", "https://api.githubcopilot.com")
	os.Setenv("QUINE_API_KEY", "copilot-oauth")
	if runtime.GOOS != "linux" {
		os.Setenv("QUINE_SELF_REENTRY_MODE", string(SelfReentryModeExecutablePath))
	}

	if _, err := Load(); err != nil {
		t.Fatalf("Load() error: %v", err)
	}
}

// --- Third-party provider test ---

func TestThirdPartyProvider(t *testing.T) {
	clearEnv(t)
	os.Setenv("QUINE_MODEL_ID", "kimi-k2.5")
	os.Setenv("QUINE_API_TYPE", "openai")
	os.Setenv("QUINE_API_BASE", "https://api.moonshot.ai/v1")
	os.Setenv("QUINE_API_KEY", "sk-moonshot-test")
	if runtime.GOOS != "linux" {
		os.Setenv("QUINE_SELF_REENTRY_MODE", string(SelfReentryModeExecutablePath))
	}

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ModelID != "kimi-k2.5" {
		t.Errorf("ModelID = %q, want kimi-k2.5", c.ModelID)
	}
	if c.Provider != "openai" {
		t.Errorf("Provider = %q, want openai", c.Provider)
	}
	if c.APIBase != "https://api.moonshot.ai/v1" {
		t.Errorf("APIBase = %q, want https://api.moonshot.ai/v1", c.APIBase)
	}
	if c.APIKey != "sk-moonshot-test" {
		t.Errorf("APIKey = %q, want sk-moonshot-test", c.APIKey)
	}
}

func TestDepthExceeded(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_DEPTH", "3")
	os.Setenv("QUINE_DEPTH", "3")

	_, err := Load()
	if !errors.Is(err, ErrDepthExceeded) {
		t.Errorf("expected ErrDepthExceeded, got: %v", err)
	}
}

func TestDepthDisabledWhenMaxDepthZero(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MAX_DEPTH", "0")
	os.Setenv("QUINE_DEPTH", "999")

	_, err := Load()
	if err != nil {
		t.Fatalf("Load() should allow QUINE_DEPTH when QUINE_MAX_DEPTH=0, got: %v", err)
	}
}

func TestContextWindow_ExplicitOverride(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_CONTEXT_WINDOW", "500000")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ContextWindow != 500_000 {
		t.Errorf("ContextWindow = %d, want 500000", c.ContextWindow)
	}
	if c.MemoryWarnTokens != 31_250 {
		t.Errorf("MemoryWarnTokens = %d, want 31250", c.MemoryWarnTokens)
	}
	if c.MemoryDangerTokens != 62_500 {
		t.Errorf("MemoryDangerTokens = %d, want 62500", c.MemoryDangerTokens)
	}
	if c.MemoryDeathTokens != 0 {
		t.Errorf("MemoryDeathTokens = %d, want 0", c.MemoryDeathTokens)
	}
}

func TestMemoryThreshold_ExplicitOverride(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MEMORY_WARN_TOKENS", "9000")
	os.Setenv("QUINE_MEMORY_DANGER_TOKENS", "15000")
	os.Setenv("QUINE_MEMORY_DEATH_TOKENS", "20000")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.MemoryWarnTokens != 9000 {
		t.Errorf("MemoryWarnTokens = %d, want 9000", c.MemoryWarnTokens)
	}
	if c.MemoryDangerTokens != 15000 {
		t.Errorf("MemoryDangerTokens = %d, want 15000", c.MemoryDangerTokens)
	}
	if c.MemoryDeathTokens != 20000 {
		t.Errorf("MemoryDeathTokens = %d, want 20000", c.MemoryDeathTokens)
	}
}

func TestMemoryThreshold_InvalidOrdering(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MEMORY_WARN_TOKENS", "9000")
	os.Setenv("QUINE_MEMORY_DANGER_TOKENS", "9000")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid memory threshold ordering")
	}
	if !strings.Contains(err.Error(), "QUINE_MEMORY_DANGER_TOKENS") {
		t.Fatalf("expected QUINE_MEMORY_DANGER_TOKENS in error, got %v", err)
	}
}

func TestMemoryDeathThreshold_InvalidOrdering(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MEMORY_WARN_TOKENS", "9000")
	os.Setenv("QUINE_MEMORY_DANGER_TOKENS", "15000")
	os.Setenv("QUINE_MEMORY_DEATH_TOKENS", "15000")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for memory death threshold at or below danger")
	}
	if !strings.Contains(err.Error(), "QUINE_MEMORY_DEATH_TOKENS") {
		t.Fatalf("expected QUINE_MEMORY_DEATH_TOKENS in error, got %v", err)
	}
}

// --- ChildEnv / ExecEnv tests ---

func TestChildEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	// Don't set QUINE_SESSION_ID - let it be auto-generated based on current PID

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}

	// Depth should be incremented
	childDepth, err := strconv.Atoi(m["QUINE_DEPTH"])
	if err != nil {
		t.Fatalf("parsing child QUINE_DEPTH: %v", err)
	}
	if childDepth != c.Depth+1 {
		t.Errorf("child QUINE_DEPTH = %d, want %d", childDepth, c.Depth+1)
	}

	// Parent session should be current session (auto-generated)
	if m["QUINE_PARENT_SESSION"] != c.SessionID {
		t.Errorf("QUINE_PARENT_SESSION = %q, want %q", m["QUINE_PARENT_SESSION"], c.SessionID)
	}

	// Session ID should NOT be present — each child generates its own
	if _, hasSessionID := m["QUINE_SESSION_ID"]; hasSessionID {
		t.Error("ChildEnv should NOT include QUINE_SESSION_ID (children generate their own)")
	}
	if _, hasRunID := m["QUINE_RUN_ID"]; hasRunID {
		t.Error("ChildEnv should NOT include QUINE_RUN_ID (children generate their own)")
	}
	if _, hasTapeID := m["QUINE_TAPE_ID"]; hasTapeID {
		t.Error("ChildEnv should NOT include QUINE_TAPE_ID (children generate their own)")
	}

	// All 4 required fields passed through
	if m["QUINE_MODEL_ID"] != "claude-sonnet-4-20250514" {
		t.Errorf("QUINE_MODEL_ID = %q, want %q", m["QUINE_MODEL_ID"], "claude-sonnet-4-20250514")
	}
	if m["QUINE_API_TYPE"] != "anthropic" {
		t.Errorf("QUINE_API_TYPE = %q, want anthropic", m["QUINE_API_TYPE"])
	}
	if m["QUINE_API_BASE"] != "https://api.anthropic.com" {
		t.Errorf("QUINE_API_BASE = %q, want https://api.anthropic.com", m["QUINE_API_BASE"])
	}
	if m["QUINE_API_KEY"] != "sk-test-key" {
		t.Errorf("QUINE_API_KEY = %q, want sk-test-key", m["QUINE_API_KEY"])
	}

	// MaxTurns should be propagated
	if m["QUINE_MAX_TURNS"] != strconv.Itoa(c.MaxTurns) {
		t.Errorf("QUINE_MAX_TURNS = %q, want %q", m["QUINE_MAX_TURNS"], strconv.Itoa(c.MaxTurns))
	}

	// All expected keys present
	expectedKeys := []string{
		"QUINE_MODEL_ID", "QUINE_API_TYPE", "QUINE_API_BASE", "QUINE_API_KEY",
		"QUINE_MAX_DEPTH", "QUINE_DEPTH", "QUINE_PARENT_SESSION",
		"QUINE_MAX_CONCURRENT", "QUINE_MAX_AGENTS", "QUINE_SH_DEFAULT_TIMEOUT_SECONDS", "QUINE_SH_TIMEOUT_OVERRIDE_ENABLED", "QUINE_SH_STDIN_ENABLED", "QUINE_SH_DETACH_ENABLED", "QUINE_OUTPUT_TRUNCATE",
		"QUINE_DATA_DIR", "QUINE_SHELL", "QUINE_SH_NETWORK", "QUINE_MAX_TURNS",
		"QUINE_PROMPT_METAPHOR", "QUINE_PROMPT_SELF_MODEL", "QUINE_PROMPT_INSTRUCTION_SURFACE", "QUINE_PROMPT_RUNTIME_SURFACE", "QUINE_PROMPT_PERSONA", "QUINE_PROMPT_CTL", "QUINE_PROMPT_IMPL_DETAILS", "QUINE_PEER_DISCOVERY_ENABLED", "QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS", "QUINE_FAIL_ON_IMPOSSIBLE",
		"QUINE_FS_MUTATION_TELEMETRY_ENABLED",
		"QUINE_READY_TEXT_AUTO_IDLE",
		"QUINE_CONTEXT_WINDOW",
		"QUINE_MEMORY_WARN_TOKENS",
		"QUINE_MEMORY_DANGER_TOKENS",
		"QUINE_MEMORY_DEATH_TOKENS",
		"QUINE_FORK_ENABLED",
		"QUINE_EXEC_ENABLED",
		"QUINE_SPAWN_ENABLED",
		"QUINE_AGENTS_MD_ENABLED",
		"QUINE_AGENTS_SKILLS_ENABLED",
		"QUINE_VISION_ENABLED",
		"QUINE_SH_INTERACTIVE_ENABLED",
		"QUINE_SH_STDIN_ENABLED",
		"QUINE_SH_DETACH_ENABLED",
		"QUINE_SELF_REENTRY_MODE",
		"QUINE_SELF_SOURCE_CODE_ENABLED",
	}
	for _, k := range expectedKeys {
		if _, ok := m[k]; !ok {
			t.Errorf("missing key %s in ChildEnv", k)
		}
	}
	if c.RetentionDir != "" {
		if _, ok := m["QUINE_RETENTION_DIR"]; !ok {
			t.Error("missing key QUINE_RETENTION_DIR in ChildEnv")
		}
	}
}

func TestChildEnv_PropagatesAnchorMemoryFlag(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_ANCHOR_MEMORY", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	found := false
	for _, e := range env {
		if e == "QUINE_ANCHOR_MEMORY=1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ChildEnv should propagate QUINE_ANCHOR_MEMORY=1 when enabled")
	}
	if !containsEnv(env, "QUINE_MEMORY_WARN_TOKENS="+strconv.Itoa(c.MemoryWarnTokens)) {
		t.Error("ChildEnv should propagate QUINE_MEMORY_WARN_TOKENS")
	}
	if !containsEnv(env, "QUINE_MEMORY_DANGER_TOKENS="+strconv.Itoa(c.MemoryDangerTokens)) {
		t.Error("ChildEnv should propagate QUINE_MEMORY_DANGER_TOKENS")
	}
	if !containsEnv(env, "QUINE_MEMORY_DEATH_TOKENS="+strconv.Itoa(c.MemoryDeathTokens)) {
		t.Error("ChildEnv should propagate QUINE_MEMORY_DEATH_TOKENS")
	}
	if !containsEnv(env, "QUINE_MEMORY_STRATEGY_HINTS=1") {
		t.Error("ChildEnv should propagate QUINE_MEMORY_STRATEGY_HINTS=1 by default")
	}
}

func TestExecEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SESSION_ID", "exec-parent-uuid")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}

	// Depth should be reset to 0
	if m["QUINE_DEPTH"] != "0" {
		t.Errorf("QUINE_DEPTH = %q, want 0", m["QUINE_DEPTH"])
	}

	if m["QUINE_SESSION_ID"] != c.SessionID {
		t.Errorf("QUINE_SESSION_ID = %q, want %q", m["QUINE_SESSION_ID"], c.SessionID)
	}
	if _, hasRunID := m["QUINE_RUN_ID"]; hasRunID {
		t.Error("ExecEnv should NOT include QUINE_RUN_ID (new process activation generates its own)")
	}
	if m["QUINE_PARENT_SESSION"] != c.ParentSession {
		t.Errorf("QUINE_PARENT_SESSION = %q, want %q", m["QUINE_PARENT_SESSION"], c.ParentSession)
	}
	if m["QUINE_TAPE_ID"] != c.TapeID {
		t.Errorf("QUINE_TAPE_ID = %q, want %q", m["QUINE_TAPE_ID"], c.TapeID)
	}

	if _, ok := m["QUINE_CONFIG_DIR"]; ok {
		t.Errorf("QUINE_CONFIG_DIR should not be set when empty")
	}
	if m["QUINE_FORK_ENABLED"] != "1" {
		t.Errorf("QUINE_FORK_ENABLED = %q, want 1", m["QUINE_FORK_ENABLED"])
	}
	if m["QUINE_PROMPT_SELF_MODEL"] != "advanced" {
		t.Errorf("QUINE_PROMPT_SELF_MODEL = %q, want advanced", m["QUINE_PROMPT_SELF_MODEL"])
	}
	if m["QUINE_PROMPT_INSTRUCTION_SURFACE"] != "standard" {
		t.Errorf("QUINE_PROMPT_INSTRUCTION_SURFACE = %q, want standard", m["QUINE_PROMPT_INSTRUCTION_SURFACE"])
	}
	if m["QUINE_PROMPT_RUNTIME_SURFACE"] != "visible" {
		t.Errorf("QUINE_PROMPT_RUNTIME_SURFACE = %q, want visible", m["QUINE_PROMPT_RUNTIME_SURFACE"])
	}
	if m["QUINE_PROMPT_PERSONA"] != "" {
		t.Errorf("QUINE_PROMPT_PERSONA = %q, want empty", m["QUINE_PROMPT_PERSONA"])
	}
	if m["QUINE_PROMPT_CTL"] != "1" {
		t.Errorf("QUINE_PROMPT_CTL = %q, want 1", m["QUINE_PROMPT_CTL"])
	}
	if m["QUINE_PROMPT_IMPL_DETAILS"] != "0" {
		t.Errorf("QUINE_PROMPT_IMPL_DETAILS = %q, want 0", m["QUINE_PROMPT_IMPL_DETAILS"])
	}
	if m["QUINE_PEER_DISCOVERY_ENABLED"] != "0" {
		t.Errorf("QUINE_PEER_DISCOVERY_ENABLED = %q, want 0", m["QUINE_PEER_DISCOVERY_ENABLED"])
	}
	if m["QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS"] != "5000" {
		t.Errorf("QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS = %q, want 5000", m["QUINE_PEER_DISCOVERY_HEARTBEAT_INTERVAL_MS"])
	}
	if m["QUINE_READY_TEXT_AUTO_IDLE"] != "0" {
		t.Errorf("QUINE_READY_TEXT_AUTO_IDLE = %q, want 0", m["QUINE_READY_TEXT_AUTO_IDLE"])
	}
	if m["QUINE_FAIL_ON_IMPOSSIBLE"] != "1" {
		t.Errorf("QUINE_FAIL_ON_IMPOSSIBLE = %q, want 1", m["QUINE_FAIL_ON_IMPOSSIBLE"])
	}
	if m["QUINE_EXEC_ENABLED"] != "1" {
		t.Errorf("QUINE_EXEC_ENABLED = %q, want 1", m["QUINE_EXEC_ENABLED"])
	}
	if m["QUINE_SPAWN_ENABLED"] != "0" {
		t.Errorf("QUINE_SPAWN_ENABLED = %q, want 0", m["QUINE_SPAWN_ENABLED"])
	}
	if m["QUINE_AGENTS_MD_ENABLED"] != "0" {
		t.Errorf("QUINE_AGENTS_MD_ENABLED = %q, want 0", m["QUINE_AGENTS_MD_ENABLED"])
	}
	if m["QUINE_AGENTS_SKILLS_ENABLED"] != "0" {
		t.Errorf("QUINE_AGENTS_SKILLS_ENABLED = %q, want 0", m["QUINE_AGENTS_SKILLS_ENABLED"])
	}
	if m["QUINE_SELF_REENTRY_MODE"] != string(c.SelfReentryMode) {
		t.Errorf("QUINE_SELF_REENTRY_MODE = %q, want %q", m["QUINE_SELF_REENTRY_MODE"], c.SelfReentryMode)
	}
	if m["QUINE_VISION_ENABLED"] != "1" {
		t.Errorf("QUINE_VISION_ENABLED = %q, want 1", m["QUINE_VISION_ENABLED"])
	}
	if m["QUINE_SH_INTERACTIVE_ENABLED"] != "1" {
		t.Errorf("QUINE_SH_INTERACTIVE_ENABLED = %q, want 1", m["QUINE_SH_INTERACTIVE_ENABLED"])
	}
	if m["QUINE_SH_STDIN_ENABLED"] != "1" {
		t.Errorf("QUINE_SH_STDIN_ENABLED = %q, want 1", m["QUINE_SH_STDIN_ENABLED"])
	}
	if m["QUINE_SH_DETACH_ENABLED"] != "1" {
		t.Errorf("QUINE_SH_DETACH_ENABLED = %q, want 1", m["QUINE_SH_DETACH_ENABLED"])
	}
	if m["QUINE_EPHEMERAL_BODY_ENABLED"] != "0" {
		t.Errorf("QUINE_EPHEMERAL_BODY_ENABLED = %q, want 0", m["QUINE_EPHEMERAL_BODY_ENABLED"])
	}
	if m["QUINE_SUPPRESS_INITIAL_BEGIN"] != "0" {
		t.Errorf("QUINE_SUPPRESS_INITIAL_BEGIN = %q, want 0", m["QUINE_SUPPRESS_INITIAL_BEGIN"])
	}
	if m["QUINE_SELF_SOURCE_CODE_ENABLED"] != "0" {
		t.Errorf("QUINE_SELF_SOURCE_CODE_ENABLED = %q, want 0", m["QUINE_SELF_SOURCE_CODE_ENABLED"])
	}
}

func TestProcessEnvFixtureExactOutput(t *testing.T) {
	cfg := &Config{
		Identity: Identity{
			ModelID:       "test-model",
			SessionID:     "session-1",
			TapeID:        "0007",
			ParentSession: "parent-0",
			Depth:         2,
		},
		Transport: Transport{
			Provider: "anthropic",
			APIBase:  "https://api.example.invalid",
			APIKey:   "sk-test",
		},
		Limits: Limits{
			MaxDepth:                  4,
			MaxConcurrent:             5,
			MaxAgents:                 6,
			ForkDefaultTimeoutSeconds: 7,
			ShTimeout:                 8,
			OutputTruncate:            9,
			MaxTurns:                  10,
			WallClockExitSeconds:      11,
			ContextWindow:             12,
			MemoryWarnTokens:          13,
			MemoryDangerTokens:        14,
			MemoryDeathTokens:         15,
			PeerDiscoveryHeartbeatMS:  5000,
		},
		ToolGates:    DefaultToolGates(),
		PromptConfig: PromptConfig{PromptSelfModel: PromptSelfModelAdvanced, PromptInstructionSurface: PromptInstructionSurfaceStandard, PromptRuntimeSurface: PromptRuntimeSurfaceVisible, FailOnImpossible: true, PromptCtlPhysics: true},
		Paths:        Paths{DataDir: "/tmp/data", WorkDir: "/tmp/work", Shell: "/bin/sh", ShNetwork: "host", SelfReentryMode: SelfReentryModeExecutablePath},
	}

	wantChild := []string{
		envKV(EnvModelID, "test-model"),
		envKV(EnvAPIType, "anthropic"),
		envKV(EnvAPIBase, "https://api.example.invalid"),
		envKV(EnvAPIKey, "sk-test"),
		envKV(EnvMaxDepth, "4"),
		envKV(EnvDepth, "3"),
		envKV(EnvParentSession, "session-1"),
		envKV(EnvMaxConcurrent, "5"),
		envKV(EnvMaxAgents, "6"),
		envKV(EnvForkDefaultTimeout, "7"),
		envKV(EnvShDefaultTimeout, "8"),
		envKV(EnvShTimeoutOverride, "1"),
		envKV(EnvShStdinEnabled, "1"),
		envKV(EnvShDetachEnabled, "1"),
		envKV(EnvOutputTruncate, "9"),
		envKV(EnvDataDir, "/tmp/data"),
		envKV(EnvWorkDir, "/tmp/work"),
		envKV(EnvShell, "/bin/sh"),
		envKV(EnvShNetwork, "host"),
		envKV(EnvSelfReentryMode, string(SelfReentryModeExecutablePath)),
		envKV(EnvMaxTurns, "10"),
		envKV(EnvWallClockExitSeconds, "11"),
		envKV(EnvPromptMetaphor, ""),
		envKV(EnvPromptSelfModel, string(PromptSelfModelAdvanced)),
		envKV(EnvPromptInstructionSurface, string(PromptInstructionSurfaceStandard)),
		envKV(EnvPromptRuntimeSurface, string(PromptRuntimeSurfaceVisible)),
		envKV(EnvPromptPersona, ""),
		envKV(EnvPromptCtl, "1"),
		envKV(EnvPromptImplDetails, "0"),
		envKV(EnvPeerDiscoveryEnabled, "0"),
		envKV(EnvPeerDiscoveryHeartbeat, "5000"),
		envKV(EnvFSMutationTelemetry, "1"),
		envKV(EnvFailOnImpossible, "1"),
		envKV(EnvNoMissionAutonomy, "off"),
		envKV(EnvEmptyAssistantSuccess, "0"),
		envKV(EnvReadyTextAutoIdle, "0"),
		envKV(EnvContextWindow, "12"),
		envKV(EnvMemoryWarnTokens, "13"),
		envKV(EnvMemoryDangerTokens, "14"),
		envKV(EnvMemoryDeathTokens, "15"),
		envKV(EnvMemoryStrategyHints, "0"),
		envKV(EnvIdleEnabled, "0"),
		envKV(EnvExitEnabled, "1"),
		envKV(EnvExecEnabled, "1"),
		envKV(EnvSpawnEnabled, "0"),
		envKV(EnvAgentsMDEnabled, "0"),
		envKV(EnvAgentsSkillsEnabled, "0"),
		envKV(EnvVisionEnabled, "1"),
		envKV(EnvForkEnabled, "1"),
		envKV(EnvShInteractiveEnabled, "1"),
		envKV(EnvForkWorldEnabled, "0"),
		envKV(EnvEphemeralBodyEnabled, "0"),
		envKV(EnvSuppressInitialBegin, "0"),
		envKV(EnvSelfSourceCodeEnabled, "0"),
	}

	childEnv, err := cfg.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	assertEnvExact(t, "ChildEnv", childEnv, wantChild)

	wantExec := append([]string(nil), wantChild...)
	wantExec[5] = envKV(EnvDepth, "0")
	wantExec[6] = envKV(EnvParentSession, "parent-0")
	wantExec = append(wantExec, envKV(EnvSessionID, "session-1"), envKV(EnvTapeID, "0007"))
	execEnv, err := cfg.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	assertEnvExact(t, "ExecEnv", execEnv, wantExec)
}

func assertEnvExact(t *testing.T, name string, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s length = %d, want %d\ngot:  %v\nwant: %v", name, len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s[%d] = %q, want %q\ngot:  %v\nwant: %v", name, i, got[i], want[i], got, want)
		}
	}
}

func TestConfigDirPassthrough(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_CONFIG_DIR", "/tmp/quine-config")
	t.Cleanup(func() {
		os.Unsetenv("QUINE_CONFIG_DIR")
	})

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	found := false
	for _, e := range env {
		if e == "QUINE_CONFIG_DIR=/tmp/quine-config" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("QUINE_CONFIG_DIR should be passed through")
	}
}

func TestEphemeralBodyPropagatesAcrossProcessTransitions(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_EPHEMERAL_BODY_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.EphemeralBody {
		t.Fatal("EphemeralBody should be enabled when QUINE_EPHEMERAL_BODY_ENABLED=1")
	}
	if c.SelfReentryMode == "" {
		t.Fatal("SelfReentryMode should be populated")
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_EPHEMERAL_BODY_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_EPHEMERAL_BODY_ENABLED=1, got %v", childEnv)
	}
	if !containsEnv(childEnv, "QUINE_SELF_REENTRY_MODE="+string(c.SelfReentryMode)) {
		t.Fatalf("ChildEnv() should preserve self reentry mode, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_EPHEMERAL_BODY_ENABLED=1") {
		t.Fatalf("ExecEnv() should propagate QUINE_EPHEMERAL_BODY_ENABLED=1, got %v", execEnv)
	}
	if !containsEnv(execEnv, "QUINE_SELF_REENTRY_MODE="+string(c.SelfReentryMode)) {
		t.Fatalf("ExecEnv() should preserve self reentry mode, got %v", execEnv)
	}
}

func TestSelfReentryTargetPropagatesAcrossProcessTransitions(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	wantSelfReentryTarget := expectedTestSelfReentryTarget(t)
	if c.SelfReentryTarget != wantSelfReentryTarget {
		t.Fatalf("SelfReentryTarget = %q, want %q", c.SelfReentryTarget, wantSelfReentryTarget)
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_SELF_REENTRY_MODE="+string(c.SelfReentryMode)) {
		t.Fatalf("ChildEnv() should propagate QUINE_SELF_REENTRY_MODE, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_SELF_REENTRY_MODE="+string(c.SelfReentryMode)) {
		t.Fatalf("ExecEnv() should propagate QUINE_SELF_REENTRY_MODE, got %v", execEnv)
	}
}

func TestSelfSourceCodeEnabledPropagatesAcrossProcessTransitions(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_SELF_SOURCE_CODE_ENABLED", "1")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.SelfSourceCodeEnabled {
		t.Fatal("SelfSourceCodeEnabled should be enabled when QUINE_SELF_SOURCE_CODE_ENABLED=1")
	}

	childEnv, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(childEnv, "QUINE_SELF_SOURCE_CODE_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_SELF_SOURCE_CODE_ENABLED=1, got %v", childEnv)
	}

	execEnv, err := c.ExecEnv()
	if err != nil {
		t.Fatalf("ExecEnv() error: %v", err)
	}
	if !containsEnv(execEnv, "QUINE_SELF_SOURCE_CODE_ENABLED=1") {
		t.Fatalf("ExecEnv() should propagate QUINE_SELF_SOURCE_CODE_ENABLED=1, got %v", execEnv)
	}
}

func TestThinkingBudgetAcceptsXHigh(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_THINKING_BUDGET", "xhigh")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ThinkingBudget != "xhigh" {
		t.Fatalf("ThinkingBudget = %q, want xhigh", c.ThinkingBudget)
	}
}

func TestThinkingBudgetDefaultsHighWhenUnset(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ThinkingBudget != "high" {
		t.Fatalf("ThinkingBudget = %q, want high when unset", c.ThinkingBudget)
	}
}

func TestServiceTierNormalizesFastToCodexPriority(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MODEL_SERVICE_TIER", "fast")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ServiceTier != "priority" {
		t.Fatalf("ServiceTier = %q, want priority", c.ServiceTier)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_MODEL_SERVICE_TIER=priority") {
		t.Fatalf("ChildEnv() should propagate QUINE_MODEL_SERVICE_TIER=priority, got %v", env)
	}
}

func TestServiceTierAcceptsPriority(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MODEL_SERVICE_TIER", "priority")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ServiceTier != "priority" {
		t.Fatalf("ServiceTier = %q, want priority", c.ServiceTier)
	}
}

func TestServiceTierRejectsInvalidValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_MODEL_SERVICE_TIER", "turbo")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_MODEL_SERVICE_TIER") {
		t.Fatalf("expected invalid service tier error, got %v", err)
	}
}

func TestPromptPersonaAcceptsKnownRoleAndPropagates(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_PERSONA", "gardener")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.PromptPersona != PromptPersonaGardener {
		t.Fatalf("PromptPersona = %q, want %q", c.PromptPersona, PromptPersonaGardener)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_PROMPT_PERSONA=gardener") {
		t.Fatalf("ChildEnv() should propagate QUINE_PROMPT_PERSONA=gardener, got %v", env)
	}
}

func TestPromptPersonaRejectsInvalidValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_PERSONA", "poet")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "QUINE_PROMPT_PERSONA") {
		t.Fatalf("expected invalid persona error, got %v", err)
	}
}

func TestInvalidPromptMetaphor(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_METAPHOR", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid QUINE_PROMPT_METAPHOR")
	}
	if !strings.Contains(err.Error(), "QUINE_PROMPT_METAPHOR") {
		t.Errorf("error should mention QUINE_PROMPT_METAPHOR, got: %v", err)
	}
}

func TestInvalidPromptSelfModel(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_SELF_MODEL", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid QUINE_PROMPT_SELF_MODEL")
	}
	if !strings.Contains(err.Error(), "QUINE_PROMPT_SELF_MODEL") {
		t.Errorf("error should mention QUINE_PROMPT_SELF_MODEL, got: %v", err)
	}
}

func TestInvalidPromptRuntimeSurface(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_RUNTIME_SURFACE", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid QUINE_PROMPT_RUNTIME_SURFACE")
	}
	if !strings.Contains(err.Error(), "QUINE_PROMPT_RUNTIME_SURFACE") {
		t.Errorf("error should mention QUINE_PROMPT_RUNTIME_SURFACE, got: %v", err)
	}
}

func TestInvalidPromptInstructionSurface(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_INSTRUCTION_SURFACE", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid QUINE_PROMPT_INSTRUCTION_SURFACE")
	}
	if !strings.Contains(err.Error(), "QUINE_PROMPT_INSTRUCTION_SURFACE") {
		t.Errorf("error should mention QUINE_PROMPT_INSTRUCTION_SURFACE, got: %v", err)
	}
}

func TestLoadAcceptsMinimalAutonomyPromptInstructionSurface(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_PROMPT_INSTRUCTION_SURFACE", "minimal_autonomy")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.PromptInstructionSurface != PromptInstructionSurfaceMinimalAutonomy {
		t.Fatalf("PromptInstructionSurface = %q, want %q",
			c.PromptInstructionSurface, PromptInstructionSurfaceMinimalAutonomy)
	}
	if !c.MinimalInstructionSurface() {
		t.Fatal("minimal_autonomy should be a collapsed minimal instruction surface")
	}
	if !c.MinimalAutonomyInstructionSurface() {
		t.Fatal("minimal_autonomy should report the autonomy variant")
	}
}

func TestWorkspaceConfigImplicitEnable(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	subdir := root + "/subdir"
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	subdirReal, err := filepath.EvalSymlinks(subdir)
	if err != nil {
		t.Fatalf("realpath subdir: %v", err)
	}
	os.Setenv("QUINE_WORKSPACE", subdir)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !c.WorkspaceEnabled {
		t.Fatal("workspace configuration should implicitly enable workspace physics")
	}
	if c.WorkspaceRoot != subdirReal {
		t.Errorf("WorkspaceRoot = %q, want %q", c.WorkspaceRoot, subdirReal)
	}
	if c.Workspace != subdirReal {
		t.Errorf("Workspace = %q, want %q", c.Workspace, subdirReal)
	}
	if c.WorkspaceBackend != "overlay" {
		t.Errorf("WorkspaceBackend = %q, want overlay", c.WorkspaceBackend)
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionRestore {
		t.Errorf("WorkspaceRevisionMode = %q, want restore", c.WorkspaceRevisionMode)
	}
}

func TestWorkspaceDefaultsToPwdWithinRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	nested := filepath.Join(root, "nested", "deep")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	nestedReal, err := filepath.EvalSymlinks(nested)
	if err != nil {
		t.Fatalf("realpath nested: %v", err)
	}
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatalf("chdir nested: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_DATA_DIR", t.TempDir())
	os.Setenv("QUINE_RETENTION_DIR", t.TempDir())

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.Workspace != nestedReal {
		t.Fatalf("Workspace = %q, want %q", c.Workspace, nestedReal)
	}
}

func TestWorkspaceMustStayWithinRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	other := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", other)

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace/root validation error")
	}
	if !strings.Contains(err.Error(), "must be within workspace root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceDirectBackendAccepted(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "direct")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.WorkspaceBackend != "direct" {
		t.Fatalf("WorkspaceBackend = %q, want direct", c.WorkspaceBackend)
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionNone {
		t.Fatalf("WorkspaceRevisionMode = %q, want none", c.WorkspaceRevisionMode)
	}
	if c.CurrentWorld() != WorldHost {
		t.Fatalf("CurrentWorld() = %q, want host", c.CurrentWorld())
	}
	if c.CurrentProtection() != ProtectionNone {
		t.Fatalf("CurrentProtection() = %q, want none", c.CurrentProtection())
	}
}

func TestWorkspaceOverlayFuseDriverAcceptedAndPropagated(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "overlay")
	os.Setenv("QUINE_WORKSPACE_OVERLAY_DRIVER", "fuse")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.WorkspaceBackend != "overlay" {
		t.Fatalf("WorkspaceBackend = %q, want overlay", c.WorkspaceBackend)
	}
	if c.WorkspaceOverlayDriver != "fuse" {
		t.Fatalf("WorkspaceOverlayDriver = %q, want fuse", c.WorkspaceOverlayDriver)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_WORKSPACE_OVERLAY_DRIVER=fuse") {
		t.Fatalf("child env should propagate forced overlay driver, got %#v", env)
	}
}

func TestWorkspaceOverlayDriverRejectsUnknownValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_OVERLAY_DRIVER", "mystery")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace overlay driver validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_OVERLAY_DRIVER") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceOverlayDriverRequiresOverlayBackend(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "direct")
	os.Setenv("QUINE_WORKSPACE_OVERLAY_DRIVER", "fuse")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace overlay driver/backend validation error")
	}
	if !strings.Contains(err.Error(), "requires QUINE_WORKSPACE_BACKEND") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceDirectBackendRejectsRestoreMode(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "direct")
	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "restore")

	_, err := Load()
	if err == nil {
		t.Fatal("expected direct backend revision-mode validation error")
	}
	if !strings.Contains(err.Error(), "only supports") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceRevisionModeAccepted(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "restore")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.WorkspaceRevisionMode != WorkspaceRevisionRestore {
		t.Fatalf("WorkspaceRevisionMode = %q, want restore", c.WorkspaceRevisionMode)
	}
	if !c.CanRestoreWorld() {
		t.Fatal("restore mode should enable switch_world")
	}
}

func TestWorkspaceRevisionModeRejectsUnknownValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "mystery")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace revision mode validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_REVISION_MODE") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceRevisionModeRequiresWorkspace(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	os.Setenv("QUINE_WORKSPACE_REVISION_MODE", "restore")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace revision mode requires workspace error")
	}
	if !strings.Contains(err.Error(), "requires workspace physics") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceBackendRejectsUnknownValue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "mystery")

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace backend validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_BACKEND") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceSourceIsRejected(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_SOURCE", root)

	_, err := Load()
	if err == nil {
		t.Fatal("expected workspace source removal error")
	}
	if !strings.Contains(err.Error(), "QUINE_WORKSPACE_SOURCE has been removed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceDataDirMustStayOutsideRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_DATA_DIR", filepath.Join(root, ".quine"))
	os.Setenv("QUINE_RETENTION_DIR", t.TempDir())

	_, err := Load()
	if err == nil {
		t.Fatal("expected QUINE_DATA_DIR outside workspace root validation error")
	}
	if !strings.Contains(err.Error(), "must be outside workspace root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceRetentionDirMustStayOutsideRoot(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_RETENTION_DIR", filepath.Join(root, "log"))

	_, err := Load()
	if err == nil {
		t.Fatal("expected QUINE_RETENTION_DIR outside workspace root validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_RETENTION_DIR") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkspaceUnsupportedOnNonLinux(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS == "linux" {
		t.Skip("workspace physics are supported on Linux")
	}

	root := t.TempDir()
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)

	_, err := Load()
	if err == nil {
		t.Fatal("expected non-Linux workspace error")
	}
	if !strings.Contains(err.Error(), "only supported on Linux") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRuntimeSurfaceBackendEnvIsRetired asserts the control-surface collapse to
// FUSE-only: QUINE_RUNTIME_SURFACE_BACKEND is no longer recognized, so setting
// it neither errors nor propagates into the child environment.
// See Paper/core/decisions/2026-06-control-surface-fuse-only.md.
func TestRuntimeSurfaceBackendEnvIsRetired(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	os.Setenv("QUINE_RUNTIME_SURFACE_BACKEND", "legacy")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	for _, kv := range env {
		if strings.HasPrefix(kv, "QUINE_RUNTIME_SURFACE_BACKEND=") {
			t.Fatalf("ChildEnv() should not propagate retired QUINE_RUNTIME_SURFACE_BACKEND, got %q", kv)
		}
	}
}

func TestForkWorldEnabled_DefaultDisabled(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if c.ForkWorldEnabled {
		t.Fatal("ForkWorldEnabled should default to false")
	}
	if c.CurrentWorld() != WorldHost {
		t.Fatalf("CurrentWorld() = %q, want host", c.CurrentWorld())
	}
	if c.CurrentProtection() != ProtectionNone {
		t.Fatalf("CurrentProtection() = %q, want none", c.CurrentProtection())
	}
}

func TestForkWorldEnabled_EnabledPropagatesToChildEnv(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	root := t.TempDir()
	os.Setenv("QUINE_FORK_WORLD_ENABLED", "1")
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", root)
	os.Setenv("QUINE_WORKSPACE_BACKEND", "direct")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}
	if !containsEnv(env, "QUINE_FORK_WORLD_ENABLED=1") {
		t.Fatalf("ChildEnv() should propagate QUINE_FORK_WORLD_ENABLED=1, got %v", env)
	}
}

func TestForkWorldEnabledRequiresExplicitWorkspacePhysics(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	os.Setenv("QUINE_FORK_WORLD_ENABLED", "1")

	_, err := Load()
	if err == nil {
		t.Fatal("expected fork world dependency validation error")
	}
	if !strings.Contains(err.Error(), "QUINE_FORK_WORLD_ENABLED=1 requires explicit workspace physics") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChildEnvIncludesWorkspacePhysics(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	if runtime.GOOS != "linux" {
		t.Skip("workspace physics are Linux-only")
	}

	root := t.TempDir()
	child := root + "/child"
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir child: %v", err)
	}
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("realpath root: %v", err)
	}
	childReal, err := filepath.EvalSymlinks(child)
	if err != nil {
		t.Fatalf("realpath child: %v", err)
	}
	os.Setenv("QUINE_WORKSPACE_ROOT", root)
	os.Setenv("QUINE_WORKSPACE", child)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	env, err := c.ChildEnv()
	if err != nil {
		t.Fatalf("ChildEnv() error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}

	if m["QUINE_WORKSPACE_ROOT"] != rootReal {
		t.Fatalf("QUINE_WORKSPACE_ROOT = %q, want %q", m["QUINE_WORKSPACE_ROOT"], rootReal)
	}
	if m["QUINE_WORKSPACE"] != childReal {
		t.Fatalf("QUINE_WORKSPACE = %q, want %q", m["QUINE_WORKSPACE"], childReal)
	}
	if m["QUINE_WORKSPACE_BACKEND"] != "overlay" {
		t.Fatalf("QUINE_WORKSPACE_BACKEND = %q, want overlay", m["QUINE_WORKSPACE_BACKEND"])
	}
	if m["QUINE_WORKSPACE_REVISION_MODE"] != "restore" {
		t.Fatalf("QUINE_WORKSPACE_REVISION_MODE = %q, want restore", m["QUINE_WORKSPACE_REVISION_MODE"])
	}
	if _, ok := m["QUINE_WORKSPACE_SESSION"]; ok {
		t.Fatalf("child env should not inherit QUINE_WORKSPACE_SESSION, got %q", m["QUINE_WORKSPACE_SESSION"])
	}
	if _, ok := m["QUINE_WORKSPACE_OWNER"]; ok {
		t.Fatalf("child env should not inherit QUINE_WORKSPACE_OWNER, got %q", m["QUINE_WORKSPACE_OWNER"])
	}
	if m["QUINE_WORKSPACE_BOOTSTRAP"] != c.WorkspaceSession {
		t.Fatalf("QUINE_WORKSPACE_BOOTSTRAP = %q, want %q", m["QUINE_WORKSPACE_BOOTSTRAP"], c.WorkspaceSession)
	}
}
