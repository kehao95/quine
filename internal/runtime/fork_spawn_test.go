package runtime

import (
	"os"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

func TestHandleForkRejectedWhenDisabled(t *testing.T) {
	cfg := testCfg(t)
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ForkEnabled = false
	rt := &Runtime{
		cfg:  cfg,
		tape: tape.NewTape(cfg.SessionID, "", cfg.Depth, cfg.ModelID),
		log:  func(string, ...any) {},
	}

	rt.handleFork(tape.ToolCall{
		ID:   "fork-1",
		Name: "fork",
		Arguments: map[string]any{
			"children": []any{
				map[string]any{"intent": "noop", "scope": "."},
			},
		},
	})

	messages := rt.tape.Messages()
	if len(messages) != 1 {
		t.Fatalf("expected 1 tape message, got %d", len(messages))
	}
	got := messages[0]
	if got.Role != tape.RoleToolResult {
		t.Fatalf("fork rejection role = %q, want %q", got.Role, tape.RoleToolResult)
	}
	if !strings.Contains(string(got.StructuredContent), "fork is disabled in this runtime") {
		t.Fatalf("fork rejection content = %s, want disabled message", string(got.StructuredContent))
	}
}

func TestHandleSpawnRejectedWhenDisabled(t *testing.T) {
	cfg := testCfg(t)
	cfg.SpawnEnabledFlag = false
	rt := &Runtime{
		cfg:  cfg,
		tape: tape.NewTape(cfg.SessionID, "", cfg.Depth, cfg.ModelID),
		log:  func(string, ...any) {},
	}

	rt.handleSpawn(tape.ToolCall{
		ID:   "spawn-1",
		Name: "spawn",
		Arguments: map[string]any{
			"children": []any{
				map[string]any{"mission": "noop"},
			},
		},
	})

	messages := rt.tape.Messages()
	if len(messages) != 1 {
		t.Fatalf("expected 1 tape message, got %d", len(messages))
	}
	got := messages[0]
	if got.Role != tape.RoleToolResult {
		t.Fatalf("spawn rejection role = %q, want %q", got.Role, tape.RoleToolResult)
	}
	if !strings.Contains(string(got.StructuredContent), "spawn is disabled in this runtime") {
		t.Fatalf("spawn rejection content = %s, want disabled message", string(got.StructuredContent))
	}
}

func TestRefreshForkEnvDoesNotLeakParentWorkspaceIdentity(t *testing.T) {
	root := t.TempDir()

	cfg := testCfg(t)
	cfg.SpawnEnabledFlag = true
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = root
	cfg.Workspace = root
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceCurrentRevision = "wr3"
	cfg.WorkspaceSession = "parent-workspace-session"
	cfg.WorkspaceOwner = true

	prevSession := os.Getenv("QUINE_WORKSPACE_SESSION")
	prevOwner := os.Getenv("QUINE_WORKSPACE_OWNER")
	prevBootstrap := os.Getenv("QUINE_WORKSPACE_BOOTSTRAP")
	prevRevision := os.Getenv("QUINE_WORKSPACE_CURRENT_REVISION")
	prevContextBootstrap := os.Getenv(tools.ContextBootstrapEnv)
	prevContextTape := os.Getenv("QUINE_CONTEXT_TAPE")
	t.Cleanup(func() {
		restoreEnv("QUINE_WORKSPACE_SESSION", prevSession)
		restoreEnv("QUINE_WORKSPACE_OWNER", prevOwner)
		restoreEnv("QUINE_WORKSPACE_BOOTSTRAP", prevBootstrap)
		restoreEnv("QUINE_WORKSPACE_CURRENT_REVISION", prevRevision)
		restoreEnv(tools.ContextBootstrapEnv, prevContextBootstrap)
		restoreEnv("QUINE_CONTEXT_TAPE", prevContextTape)
	})
	os.Setenv("QUINE_WORKSPACE_SESSION", "leaked-parent-session")
	os.Setenv("QUINE_WORKSPACE_OWNER", "1")
	os.Setenv("QUINE_WORKSPACE_BOOTSTRAP", "leaked-bootstrap")
	os.Setenv("QUINE_WORKSPACE_CURRENT_REVISION", "wr999")
	os.Setenv(tools.ContextBootstrapEnv, "leaked-context-bootstrap")
	os.Setenv("QUINE_CONTEXT_TAPE", "leaked-context-tape")

	rt := NewWithProvider(cfg, &mockProvider{})
	rt.refreshForkEnv()
	rt.refreshSpawnEnv()

	envMap := make(map[string]string)
	for _, entry := range rt.fork.Env {
		k, v, _ := strings.Cut(entry, "=")
		envMap[k] = v
	}

	if got := envMap["QUINE_WORKSPACE_BOOTSTRAP"]; got != cfg.WorkspaceSession {
		t.Fatalf("QUINE_WORKSPACE_BOOTSTRAP = %q, want %q", got, cfg.WorkspaceSession)
	}
	if got := envMap["QUINE_WORKSPACE_CURRENT_REVISION"]; got != cfg.WorkspaceCurrentRevision {
		t.Fatalf("QUINE_WORKSPACE_CURRENT_REVISION = %q, want %q", got, cfg.WorkspaceCurrentRevision)
	}
	if _, ok := envMap["QUINE_WORKSPACE_SESSION"]; ok {
		t.Fatalf("fork env should not include QUINE_WORKSPACE_SESSION, got %q", envMap["QUINE_WORKSPACE_SESSION"])
	}
	if _, ok := envMap["QUINE_WORKSPACE_OWNER"]; ok {
		t.Fatalf("fork env should not include QUINE_WORKSPACE_OWNER, got %q", envMap["QUINE_WORKSPACE_OWNER"])
	}
	if _, ok := envMap[tools.ContextBootstrapEnv]; ok {
		t.Fatalf("fork env should not include %s, got %q", tools.ContextBootstrapEnv, envMap[tools.ContextBootstrapEnv])
	}
	if _, ok := envMap["QUINE_CONTEXT_TAPE"]; ok {
		t.Fatalf("fork env should not include QUINE_CONTEXT_TAPE, got %q", envMap["QUINE_CONTEXT_TAPE"])
	}
	if rt.spawn == nil {
		t.Fatal("spawn executor should be initialized when spawn is enabled")
	}
	spawnEnvMap := make(map[string]string)
	for _, entry := range rt.spawn.Env {
		k, v, _ := strings.Cut(entry, "=")
		spawnEnvMap[k] = v
	}
	if _, ok := spawnEnvMap[tools.ContextBootstrapEnv]; ok {
		t.Fatalf("spawn env should not include %s, got %q", tools.ContextBootstrapEnv, spawnEnvMap[tools.ContextBootstrapEnv])
	}
	if _, ok := spawnEnvMap["QUINE_CONTEXT_TAPE"]; ok {
		t.Fatalf("spawn env should not include QUINE_CONTEXT_TAPE, got %q", spawnEnvMap["QUINE_CONTEXT_TAPE"])
	}
}

func TestWorldResourcesIncludesSpawnWhenEnabled(t *testing.T) {
	cfg := testCfg(t)
	cfg.SpawnEnabledFlag = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)

	resources := rt.runtimeWorldResources()
	if !containsString(resources.Tools, "spawn") {
		t.Fatalf("resources tools = %v, want spawn when enabled", resources.Tools)
	}
}

func TestProcessTrackingCallbacks(t *testing.T) {
	// Verify that shExecutor's ProcessStarted/ProcessEnded callbacks
	// are wired correctly to Runtime's activeProcess tracking.
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "echo hello",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
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
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("test process tracking", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// After Run completes, activeProcess should be nil
	if proc := rt.activeProcess.Load(); proc != nil {
		t.Errorf("expected activeProcess to be nil after Run, got pid=%d", proc.Pid)
	}
}

// ---------------------------------------------------------------------------
// stdin parameter tests (LLM-behaviour: does the runtime wire stdin through?)
// ---------------------------------------------------------------------------

// TestShStdinParameter verifies that when the LLM sends a sh tool call with
// a "stdin" argument, the runtime parses it and passes it to Execute so the
// data is piped into the command.
