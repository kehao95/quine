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

	// The absence assertions below are about what the runtime does NOT
	// manufacture, so the ambient environment must not be able to supply these
	// names behind the test's back: a developer shell exporting QUINE_MAX_DEPTH
	// would otherwise make the assertion pass or fail for the wrong reason.
	for _, name := range []string{"QUINE_MAX_DEPTH", "QUINE_SPAWN_ENABLED", "QUINE_FORK_ENABLED", "QUINE_OUTPUT_TRUNCATE"} {
		clearEnvForTest(t, name)
	}

	rt := NewWithProvider(cfg, &mockProvider{})
	rt.refreshForkEnv()
	rt.refreshSpawnEnv()

	envMap := make(map[string]string)
	for _, entry := range rt.fork.Env {
		k, v, _ := strings.Cut(entry, "=")
		envMap[k] = v
	}

	// The parent's own workspace identity is runtime-emitted, therefore masked:
	// none of the four leaked values above may cross the fork boundary. What the
	// child gets instead is what the runtime stamps, and only that.

	// Adoption, not inheritance: the child takes the parent's workspace SESSION
	// as its BOOTSTRAP so it joins the lineage — and it is the parent's real
	// session, not the "leaked-bootstrap" that was sitting in the environ.
	if got := envMap["QUINE_WORKSPACE_BOOTSTRAP"]; got != cfg.WorkspaceSession {
		t.Fatalf("QUINE_WORKSPACE_BOOTSTRAP = %q, want the parent's workspace session %q", got, cfg.WorkspaceSession)
	}
	// A child never owns its parent's workspace; it borrows a view of one. The
	// stamp says so explicitly, and it beats the leaked owner=1 in the environ.
	if got, ok := envMap["QUINE_WORKSPACE_OWNER"]; !ok || got != "0" {
		t.Fatalf("QUINE_WORKSPACE_OWNER = %q (present=%v), want the stamped \"0\" — a child borrows a workspace, it does not own one", got, ok)
	}
	// INVERSION (brief E-matrix, fork row): this used to assert that the child
	// inherits the PARENT's QUINE_WORKSPACE_CURRENT_REVISION ("wr3"). That is
	// exactly the leak this test is named after. A revision handle names a place
	// in the parent's own workspace state; the child mounts its own view and
	// computes its own revision. Neither the parent's real revision nor the
	// leaked "wr999" may appear.
	if got, ok := envMap["QUINE_WORKSPACE_CURRENT_REVISION"]; ok {
		t.Fatalf("QUINE_WORKSPACE_CURRENT_REVISION = %q, want ABSENT: a revision handle is a place in the PARENT's workspace state, not a fact about the child", got)
	}
	// The child gets its own session id from the fork executor, never the
	// parent's workspace session.
	if _, ok := envMap["QUINE_WORKSPACE_SESSION"]; ok {
		t.Fatalf("fork env should not include QUINE_WORKSPACE_SESSION, got %q", envMap["QUINE_WORKSPACE_SESSION"])
	}
	if _, ok := envMap[tools.ContextBootstrapEnv]; ok {
		t.Fatalf("fork env should not include %s, got %q", tools.ContextBootstrapEnv, envMap[tools.ContextBootstrapEnv])
	}
	if _, ok := envMap["QUINE_CONTEXT_TAPE"]; ok {
		t.Fatalf("fork env should not include QUINE_CONTEXT_TAPE, got %q", envMap["QUINE_CONTEXT_TAPE"])
	}
	// Tree membership: the fork boundary is where lineage marks are stamped, and
	// this is the boundary where they mean something (contrast BoundaryShell,
	// where they are masked — a program started from a shell is not a member of
	// this agent tree).
	if got := envMap["QUINE_DEPTH"]; got != "1" {
		t.Fatalf("QUINE_DEPTH = %q, want \"1\" (parent depth 0 + 1)", got)
	}
	if got := envMap["QUINE_PARENT_SESSION"]; got != cfg.SessionID {
		t.Fatalf("QUINE_PARENT_SESSION = %q, want %q", got, cfg.SessionID)
	}
	// The manufactured-evidence inversion, at the fork boundary: a knob nobody
	// set is ABSENT from the child's env. cfg.MaxDepth is 5 and the tool gates
	// are off, but no operator authored those, so no child may read them as if
	// someone had.
	for _, unset := range []string{"QUINE_MAX_DEPTH", "QUINE_SPAWN_ENABLED", "QUINE_FORK_ENABLED", "QUINE_OUTPUT_TRUNCATE"} {
		if got, ok := envMap[unset]; ok {
			t.Fatalf("fork child env manufactures %s=%q for a knob nobody set; absence is how an unset knob is spelled", unset, got)
		}
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
