package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	quineroot "github.com/kehao95/quine"
	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

func TestRunSyncsAgentRoot(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
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

	if exitCode := rt.Run("sync agent root", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	retainedRoot := cfg.SessionLogDir("")
	statusPath := filepath.Join(retainedRoot, "status", "session.json")
	statusData, err := os.ReadFile(statusPath)
	if err != nil {
		t.Fatalf("read status file: %v", err)
	}
	var status map[string]any
	if err := json.Unmarshal(statusData, &status); err != nil {
		t.Fatalf("unmarshal status file: %v", err)
	}
	if status["session_id"] != cfg.SessionID {
		t.Fatalf("status.session_id = %v, want %q", status["session_id"], cfg.SessionID)
	}
	if status["run_id"] != cfg.RunID {
		t.Fatalf("status.run_id = %v, want %q", status["run_id"], cfg.RunID)
	}
	if status["incarnation_id"] != float64(0) {
		t.Fatalf("status.incarnation_id = %v, want 0", status["incarnation_id"])
	}
	if status["pid"] != float64(os.Getpid()) {
		t.Fatalf("status.pid = %v, want %d", status["pid"], os.Getpid())
	}
	if status["ppid"] != float64(os.Getppid()) {
		t.Fatalf("status.ppid = %v, want %d", status["ppid"], os.Getppid())
	}
	if _, ok := status["tape_id"]; ok {
		t.Fatalf("status should not expose tape_id, got %v", status["tape_id"])
	}

	missionData, err := os.ReadFile(filepath.Join(retainedRoot, "mission.txt"))
	if err != nil {
		t.Fatalf("read mission file: %v", err)
	}
	if strings.TrimSpace(string(missionData)) != "sync agent root" {
		t.Fatalf("mission.txt = %q, want %q", strings.TrimSpace(string(missionData)), "sync agent root")
	}
	assertSymlinkTarget(t, filepath.Join(retainedRoot, "inc", "current"), "0")
	incMissionData, err := os.ReadFile(filepath.Join(retainedRoot, "inc", "0", "mission.txt"))
	if err != nil {
		t.Fatalf("read retained incarnation mission: %v", err)
	}
	if strings.TrimSpace(string(incMissionData)) != "sync agent root" {
		t.Fatalf("retained incarnation mission = %q, want %q", strings.TrimSpace(string(incMissionData)), "sync agent root")
	}

	if _, err := os.Lstat(filepath.Join(retainedRoot, "current.jsonl")); !os.IsNotExist(err) {
		t.Fatalf("retained log/current.jsonl should be absent, got err=%v", err)
	}

	if info, err := os.Stat(filepath.Join(retainedRoot, "tapes")); err != nil {
		t.Fatalf("stat retained tapes dir: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("retained tapes surface should be a directory, got %v", info.Mode())
	}
	assertSameFile(t, filepath.Join(retainedRoot, "tapes"), cfg.TapeDir(""))

	runtimeLogData, err := os.ReadFile(cfg.SessionLogPath(""))
	if err != nil {
		t.Fatalf("read runtime log: %v", err)
	}
	runtimeLog := string(runtimeLogData)
	if strings.Contains(runtimeLog, "turn ") {
		t.Fatalf("runtime.log should not contain turn-level entries, got %q", runtimeLog)
	}
	if strings.Contains(runtimeLog, "THE PRIME DIRECTIVE") {
		t.Fatalf("runtime.log should not duplicate system prompt content, got %q", runtimeLog)
	}

	if _, err := os.Stat(filepath.Join(retainedRoot, "jobs")); !os.IsNotExist(err) {
		t.Fatalf("expected retained jobs projection removed, got err=%v", err)
	}

	if _, err := os.Stat(cfg.AgentRoot()); !os.IsNotExist(err) {
		t.Fatalf("agent root should be removed after exit, got err=%v", err)
	}
	if _, err := os.Stat(cfg.ControlPath()); !os.IsNotExist(err) {
		t.Fatalf("ctl surface should be removed after exit, got err=%v", err)
	}

	inboxData, err := os.ReadFile(filepath.Join(retainedRoot, "status", "inbox.json"))
	if err != nil {
		t.Fatalf("read inbox snapshot: %v", err)
	}
	var inbox map[string]any
	if err := json.Unmarshal(inboxData, &inbox); err != nil {
		t.Fatalf("unmarshal inbox snapshot: %v", err)
	}
	if pending := toolInt(t, inbox, "pending_count"); pending != 0 {
		t.Fatalf("inbox pending_count = %d, want 0", pending)
	}

	if _, err := os.Stat(cfg.SessionControlLogPath("")); err != nil {
		t.Fatalf("stat control log: %v", err)
	}

	for _, legacyPath := range []string{
		filepath.Join(cfg.RuntimeRoot(), "agent", "live"),
		filepath.Join(cfg.RuntimeRoot(), "agent", "live_by_pid"),
	} {
		if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
			t.Fatalf("expected legacy live index %s removed, got err=%v", legacyPath, err)
		}
	}
}

func TestBootstrapAgentRootProjectsLiveAndMirrorSurfaces(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "sync agent root"
	rt.lastFSMutations = "mutation block"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	if cfg.IncarnationID != 0 {
		t.Fatalf("cfg.IncarnationID = %d, want 0", cfg.IncarnationID)
	}

	liveRoot := cfg.AgentRoot()
	mirrorRoot := cfg.SessionLogDir("")
	currentIncRoot := cfg.SessionIncarnationPath("", 0)

	for _, path := range []string{
		filepath.Join(liveRoot, "inc"),
		filepath.Join(liveRoot, "inc", "current"),
		currentIncRoot,
		filepath.Join(currentIncRoot, "context"),
		filepath.Join(currentIncRoot, "mission.txt"),
		filepath.Join(liveRoot, "mission.txt"),
		filepath.Join(liveRoot, "context"),
		filepath.Join(liveRoot, "status", "session.json"),
		filepath.Join(liveRoot, "status", "inbox.json"),
		filepath.Join(liveRoot, "status", "contract.json"),
		filepath.Join(liveRoot, "world", "workspace_root"),
		filepath.Join(liveRoot, "world", "workspace"),
		filepath.Join(liveRoot, "world", "status.json"),
		filepath.Join(liveRoot, "world", "resources.json"),
		filepath.Join(liveRoot, "world", "events.jsonl"),
		filepath.Join(liveRoot, "tapes"),
		filepath.Join(liveRoot, "jobs"),
		filepath.Join(liveRoot, "log"),
		filepath.Join(mirrorRoot, "mission.txt"),
		filepath.Join(mirrorRoot, "status", "session.json"),
		filepath.Join(mirrorRoot, "status", "inbox.json"),
		filepath.Join(mirrorRoot, "status", "contract.json"),
		filepath.Join(mirrorRoot, "world", "workspace_root"),
		filepath.Join(mirrorRoot, "world", "workspace"),
		filepath.Join(mirrorRoot, "world", "status.json"),
		filepath.Join(mirrorRoot, "world", "resources.json"),
		filepath.Join(mirrorRoot, "world", "events.jsonl"),
		filepath.Join(mirrorRoot, "tapes"),
	} {
		if _, err := os.Lstat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	assertSymlinkTarget(t, filepath.Join(liveRoot, "inc", "current"), "0")
	assertSymlinkTarget(t, filepath.Join(liveRoot, "mission.txt"), filepath.Join("inc", "current", "mission.txt"))
	assertSymlinkTarget(t, filepath.Join(liveRoot, "context"), filepath.Join("inc", "current", "context"))
	assertSameFile(t, filepath.Join(liveRoot, "mission.txt"), filepath.Join(currentIncRoot, "mission.txt"))
	assertSameFile(t, filepath.Join(liveRoot, "mission.txt"), filepath.Join(mirrorRoot, "mission.txt"))
	assertSameFile(t, filepath.Join(liveRoot, "status", "session.json"), filepath.Join(mirrorRoot, "status", "session.json"))
	assertSameFile(t, filepath.Join(liveRoot, "status", "inbox.json"), filepath.Join(mirrorRoot, "status", "inbox.json"))
	assertSameFile(t, filepath.Join(liveRoot, "status", "contract.json"), filepath.Join(mirrorRoot, "status", "contract.json"))
	assertSameFile(t, filepath.Join(liveRoot, "world", "workspace_root"), filepath.Join(mirrorRoot, "world", "workspace_root"))
	assertSameFile(t, filepath.Join(liveRoot, "world", "workspace"), filepath.Join(mirrorRoot, "world", "workspace"))
	assertSameFile(t, filepath.Join(liveRoot, "world", "status.json"), filepath.Join(mirrorRoot, "world", "status.json"))
	assertSameFile(t, filepath.Join(liveRoot, "world", "resources.json"), filepath.Join(mirrorRoot, "world", "resources.json"))

	assertRelativeSymlinkResolvesTo(t, filepath.Join(liveRoot, "log"), mirrorRoot)
	assertRelativeSymlinkResolvesTo(t, filepath.Join(liveRoot, "inc"), filepath.Join(mirrorRoot, "inc"))
	assertRelativeSymlinkResolvesTo(t, filepath.Join(liveRoot, "status"), filepath.Join(mirrorRoot, "status"))
	assertRelativeSymlinkResolvesTo(t, filepath.Join(liveRoot, "world"), filepath.Join(mirrorRoot, "world"))
	assertRelativeSymlinkResolvesTo(t, filepath.Join(liveRoot, "jobs"), cfg.JobSessionDir(""))
	assertRelativeSymlinkResolvesTo(t, filepath.Join(liveRoot, "tapes"), cfg.TapeDir(""))
	// The public/ surface is a FUSE projection (control-surface-fuse-only); its
	// node-level structure is covered by TestBootstrapAgentRootProjectsFusePublicSurface.
	// Here we validate the underlying retained contract content, which the FUSE
	// status/contract.json node simply projects.
	contractData, err := os.ReadFile(filepath.Join(liveRoot, "status", "contract.json"))
	if err != nil {
		t.Fatalf("read retained contract: %v", err)
	}
	var contract map[string]any
	if err := json.Unmarshal(contractData, &contract); err != nil {
		t.Fatalf("unmarshal public contract: %v", err)
	}
	if contract["contract_version"] != processControlContractVersion {
		t.Fatalf("contract_version = %v, want %q", contract["contract_version"], processControlContractVersion)
	}
	actions, ok := contract["control_actions"].(map[string]any)
	if !ok {
		t.Fatalf("contract missing control_actions: %#v", contract)
	}
	post, ok := actions["post"].(map[string]any)
	if !ok || post["queues"] != true {
		t.Fatalf("contract post action = %#v, want queues=true", actions["post"])
	}
	// The control surface must self-describe so a peer (which does not share this
	// agent's prompt) can operate it from the peer-readable contract alone.
	if desc, _ := post["description"].(string); !strings.Contains(desc, "queue-only") {
		t.Fatalf("contract post action missing self-describing description, got %#v", post["description"])
	}
	if _, ok := contract["inbox_schema"].(map[string]any); !ok {
		t.Fatalf("contract missing inbox_schema: %#v", contract)
	}
	if _, ok := contract["control_log_events"].(map[string]any); !ok {
		t.Fatalf("contract missing control_log_events: %#v", contract)
	}
	if usage, _ := contract["usage"].(string); !strings.Contains(usage, "ctl/") {
		t.Fatalf("contract missing usage guidance, got %#v", contract["usage"])
	}

	worldData, err := os.ReadFile(filepath.Join(liveRoot, "world", "status.json"))
	if err != nil {
		t.Fatalf("read world status: %v", err)
	}
	var worldStatus map[string]any
	if err := json.Unmarshal(worldData, &worldStatus); err != nil {
		t.Fatalf("unmarshal world status: %v", err)
	}
	if worldStatus["contract_version"] != worldSurfaceContractVersion {
		t.Fatalf("world contract_version = %v, want %q", worldStatus["contract_version"], worldSurfaceContractVersion)
	}
	if worldStatus["world"] != "host" {
		t.Fatalf("world status world = %v, want host", worldStatus["world"])
	}
	for _, absent := range []string{"workspace_enabled", "revision_mode", "can_rollback", "can_switch_revision", "fs_mutation_telemetry"} {
		if _, ok := worldStatus[absent]; ok {
			t.Fatalf("world status should omit absent capability field %q: %#v", absent, worldStatus)
		}
	}

	resourcesData, err := os.ReadFile(filepath.Join(liveRoot, "world", "resources.json"))
	if err != nil {
		t.Fatalf("read world resources: %v", err)
	}
	var resources map[string]any
	if err := json.Unmarshal(resourcesData, &resources); err != nil {
		t.Fatalf("unmarshal world resources: %v", err)
	}
	toolsRaw, ok := resources["tools"].([]any)
	if !ok {
		t.Fatalf("world resources tools = %#v, want array", resources["tools"])
	}
	var gotTools []string
	for _, raw := range toolsRaw {
		tool, ok := raw.(string)
		if !ok {
			t.Fatalf("world resources tool = %#v, want string", raw)
		}
		gotTools = append(gotTools, tool)
	}
	wantTools := []string{"sh", "fork", "exit", "exec", "vision"}
	if strings.Join(gotTools, ",") != strings.Join(wantTools, ",") {
		t.Fatalf("world resources tools = %v, want %v", gotTools, wantTools)
	}
	for _, absent := range []string{"idle", "switch_world"} {
		if strings.Contains(strings.Join(gotTools, ","), absent) {
			t.Fatalf("world resources leaked absent tool %q in %v", absent, gotTools)
		}
	}
	if _, ok := resources["workspace"]; ok {
		t.Fatalf("world resources should omit workspace when workspace is absent: %#v", resources["workspace"])
	}
	if strings.Contains(string(resourcesData), "revision switching") {
		t.Fatalf("world resources should not mention revision switching when unavailable: %s", string(resourcesData))
	}
	eventsData, err := os.ReadFile(filepath.Join(liveRoot, "world", "events.jsonl"))
	if err != nil {
		t.Fatalf("read world events: %v", err)
	}
	if !strings.Contains(string(eventsData), `"kind":"mutation_observed"`) || !strings.Contains(string(eventsData), "mutation block") {
		t.Fatalf("world events missing mutation observation: %q", string(eventsData))
	}

	if _, err := os.Stat(filepath.Join(mirrorRoot, "jobs")); !os.IsNotExist(err) {
		t.Fatalf("expected retained jobs projection removed, got err=%v", err)
	}
}

func TestWorldResourcesListOnlyPresentAffordances(t *testing.T) {
	cfg := testCfg(t)
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ForkEnabled = false
	cfg.ExecEnabled = false
	cfg.ToolGates.ExitEnabled = false
	cfg.VisionEnabled = false
	cfg.IdleEnabled = false
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = t.TempDir()
	cfg.Workspace = filepath.Join(cfg.WorkspaceRoot, "app")
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)

	resources := rt.runtimeWorldResources()
	if got, want := strings.Join(resources.Tools, ","), "sh,switch_world"; got != want {
		t.Fatalf("resources tools = %q, want %q", got, want)
	}
	if resources.Workspace == nil {
		t.Fatal("resources should include workspace when workspace exists")
	}
	if !containsString(resources.NonClaims, "revision switching is available only when the active workspace backend supports restore") {
		t.Fatalf("resources non_claims missing revision-switch boundary: %v", resources.NonClaims)
	}

	status := rt.runtimeWorldStatus()
	if status.RevisionMode != string(config.WorkspaceRevisionRestore) {
		t.Fatalf("status revision mode = %q, want restore", status.RevisionMode)
	}
	if !status.CanSwitchRevision {
		t.Fatal("status should expose can_switch_revision only when available")
	}
}

func TestWriteTapeEntryDoesNotBootstrapAgentRoot(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "write tape without sync"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	t.Cleanup(func() {
		_ = rt.cleanupAgentRoot()
	})

	statusPath := filepath.Join(cfg.SessionRetainedDir(""), "status", "session.json")
	if err := os.Remove(statusPath); err != nil {
		t.Fatalf("remove status fixture: %v", err)
	}

	rt.tape = tape.NewTape(cfg.SessionID, cfg.ParentSession, cfg.Depth, cfg.ModelID)
	writer, err := tape.NewWriter(cfg.TapeDir(""), cfg.TapeID)
	if err != nil {
		t.Fatalf("new tape writer: %v", err)
	}
	rt.tapeWriter = writer
	defer writer.Close()

	rt.writeTapeEntry(rt.tape.MetaEntry())

	if _, err := os.Stat(statusPath); !os.IsNotExist(err) {
		t.Fatalf("writeTapeEntry should not recreate %s, got err=%v", statusPath, err)
	}
}

func TestBootstrapAgentRootResumesSameIncarnationForSameSession(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "first incarnation"
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("first bootstrapAgentRoot failed: %v", err)
	}

	firstCurrentPath := filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl")
	writeJSONLEntries(t, firstCurrentPath, tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "old context"}))
	firstMemoryPath := filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "prompt", "30-memory.md")
	if err := os.WriteFile(firstMemoryPath, []byte("old memory context\n"), 0o644); err != nil {
		t.Fatalf("write first memory: %v", err)
	}

	cfg2 := *cfg
	cfg2.RunID = "run-test-5678"
	rt2 := NewWithProvider(&cfg2, &mockProvider{})
	silenceRuntime(rt2)
	rt2.originalInput = "second incarnation"
	t.Cleanup(func() {
		_ = rt2.cleanupAgentRoot()
	})
	if err := rt2.bootstrapAgentRoot(); err != nil {
		t.Fatalf("second bootstrapAgentRoot failed: %v", err)
	}
	if cfg2.IncarnationID != 0 {
		t.Fatalf("cfg2.IncarnationID = %d, want resumed incarnation 0", cfg2.IncarnationID)
	}
	statusData, err := os.ReadFile(filepath.Join(cfg2.SessionRetainedDir(""), "status", "session.json"))
	if err != nil {
		t.Fatalf("read resumed status: %v", err)
	}
	var status map[string]any
	if err := json.Unmarshal(statusData, &status); err != nil {
		t.Fatalf("unmarshal resumed status: %v", err)
	}
	if status["session_id"] != cfg.SessionID {
		t.Fatalf("resumed session_id = %v, want %q", status["session_id"], cfg.SessionID)
	}
	if status["run_id"] != cfg2.RunID {
		t.Fatalf("resumed run_id = %v, want %q", status["run_id"], cfg2.RunID)
	}
	assertSymlinkTarget(t, filepath.Join(cfg2.AgentRoot(), "inc", "current"), "0")
	if _, err := os.Stat(firstCurrentPath); err != nil {
		t.Fatalf("expected first incarnation context to remain visible: %v", err)
	}
	msgs, err := rt2.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages() error: %v", err)
	}
	foundOldContext := false
	foundOldMemory := false
	for _, msg := range msgs {
		if strings.Contains(msg.Content, "old context") {
			foundOldContext = true
		}
		if strings.Contains(msg.Content, "old memory context") {
			foundOldMemory = true
		}
	}
	if !foundOldContext {
		t.Fatalf("new incarnation provider context should inherit old raw context: %#v", msgs)
	}
	if !foundOldMemory {
		t.Fatalf("new incarnation provider prompt should inherit memory context: %#v", msgs)
	}
}

func TestBootstrapAgentRootProjectsSelfSourceSurfaceWhenEnabled(t *testing.T) {
	cfg := testCfg(t)
	cfg.SelfSourceCodeEnabled = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "sync agent root with self source"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}

	sourceRoot := filepath.Join(cfg.AgentRoot(), "source-code")
	for _, rel := range []string{
		"go.mod",
		"go.sum",
		"selfsource.go",
		filepath.Join("cmd", "quine", "main.go"),
		filepath.Join("internal", "runtime", "runtime.go"),
	} {
		target := filepath.Join(sourceRoot, rel)
		info, err := os.Stat(target)
		if err != nil {
			t.Fatalf("expected self-source file %s: %v", target, err)
		}
		if info.IsDir() {
			t.Fatalf("expected file at %s, got directory", target)
		}
		if info.Mode().Perm() != 0o444 {
			t.Fatalf("self-source file %s permissions = %o, want 444", target, info.Mode().Perm())
		}

	}

	rootInfo, err := os.Stat(sourceRoot)
	if err != nil {
		t.Fatalf("stat self-source root: %v", err)
	}
	if !rootInfo.IsDir() {
		t.Fatalf("self-source root should be a directory, got %v", rootInfo.Mode())
	}
	if rootInfo.Mode().Perm() != 0o555 {
		t.Fatalf("self-source root permissions = %o, want 555", rootInfo.Mode().Perm())
	}
	// The public/source-code FUSE projection is covered by
	// TestBootstrapAgentRootProjectsFusePublicSurface; here we validate the
	// retained source tree and manifest that the projection reflects.
	manifestData, err := os.ReadFile(selfSourceManifestPath(sourceRoot))
	if err != nil {
		t.Fatalf("read self-source manifest: %v", err)
	}
	var manifest selfSourceManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatalf("unmarshal self-source manifest: %v", err)
	}
	wantManifest, err := currentSelfSourceManifest()
	if err != nil {
		t.Fatalf("compute current self-source manifest: %v", err)
	}
	if manifest != wantManifest {
		t.Fatalf("self-source manifest = %#v, want %#v", manifest, wantManifest)
	}
	if err := runSelfSourceGit(sourceRoot, "rev-parse", "--is-inside-work-tree"); err != nil {
		t.Fatalf("source-code should be a readable git worktree: %v", err)
	}
	if _, err := os.Stat(filepath.Join(sourceRoot, "cmd", "world")); err != nil {
		t.Fatalf("repository-backed source-code should expose full repo paths such as cmd/world: %v", err)
	}
	copyRoot := filepath.Join(t.TempDir(), "source-copy")
	if err := runSelfSourceGit("", "clone", "-q", sourceRoot, copyRoot); err != nil {
		t.Fatalf("git clone source-code for mutation: %v", err)
	}
	if err := runSelfSourceGit(copyRoot, "config", "user.name", "Quine Test"); err != nil {
		t.Fatalf("configure copied source git user: %v", err)
	}
	if err := runSelfSourceGit(copyRoot, "config", "user.email", "quine-test@example.invalid"); err != nil {
		t.Fatalf("configure copied source git email: %v", err)
	}
	marker := filepath.Join(copyRoot, "mutation-test.txt")
	if err := os.WriteFile(marker, []byte("copy can mutate and commit\n"), 0o644); err != nil {
		t.Fatalf("write mutation marker: %v", err)
	}
	if err := runSelfSourceGit(copyRoot, "add", "mutation-test.txt"); err != nil {
		t.Fatalf("git add in copied source: %v", err)
	}
	if err := runSelfSourceGit(copyRoot, "commit", "-q", "--no-gpg-sign", "-m", "Test copied source mutation"); err != nil {
		t.Fatalf("git commit in copied source: %v", err)
	}
}

func TestBootstrapAgentRootRefreshesSelfSourceSurfaceWhenManifestDrifts(t *testing.T) {
	cfg := testCfg(t)
	cfg.SelfSourceCodeEnabled = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "refresh stale self source surface"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("initial bootstrapAgentRoot failed: %v", err)
	}

	sourceRoot := filepath.Join(cfg.AgentRoot(), "source-code")
	if err := setWritableTree(sourceRoot); err != nil {
		t.Fatalf("setWritableTree failed: %v", err)
	}
	if err := os.WriteFile(selfSourceManifestPath(sourceRoot), []byte("{\"format\":\"stale\"}\n"), 0o644); err != nil {
		t.Fatalf("write stale manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceRoot, "go.mod"), []byte("module stale.example\n"), 0o644); err != nil {
		t.Fatalf("write stale go.mod: %v", err)
	}
	if err := setReadOnlyTree(sourceRoot); err != nil {
		t.Fatalf("setReadOnlyTree failed: %v", err)
	}

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("second bootstrapAgentRoot failed: %v", err)
	}
	gotGoMod, err := os.ReadFile(filepath.Join(sourceRoot, "go.mod"))
	if err != nil {
		t.Fatalf("read refreshed go.mod: %v", err)
	}
	if string(gotGoMod) == "module stale.example\n" {
		t.Fatalf("self-source go.mod was not refreshed after manifest drift")
	}
	manifestData, err := os.ReadFile(selfSourceManifestPath(sourceRoot))
	if err != nil {
		t.Fatalf("read refreshed manifest: %v", err)
	}
	var manifest selfSourceManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatalf("unmarshal refreshed manifest: %v", err)
	}
	wantManifest, err := currentSelfSourceManifest()
	if err != nil {
		t.Fatalf("compute current self-source manifest: %v", err)
	}
	if manifest != wantManifest {
		t.Fatalf("refreshed self-source manifest = %#v, want %#v", manifest, wantManifest)
	}
}

func TestBootstrapAgentRootRequiresSelfSourceBundleWhenEnabled(t *testing.T) {
	oldBundle := quineroot.SelfSourceBundle
	quineroot.SelfSourceBundle = nil
	t.Cleanup(func() {
		quineroot.SelfSourceBundle = oldBundle
	})

	cfg := testCfg(t)
	cfg.SelfSourceCodeEnabled = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "require self source bundle"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	err := rt.bootstrapAgentRoot()
	if err == nil {
		t.Fatalf("bootstrapAgentRoot should fail when self-source projection is enabled without an embedded repository bundle")
	}
	if !strings.Contains(err.Error(), "embedded self-source repository bundle is not available") {
		t.Fatalf("bootstrapAgentRoot error = %v", err)
	}
}

func TestBootstrapAgentRootRemovesSelfSourceSurfaceWhenDisabled(t *testing.T) {
	cfg := testCfg(t)
	cfg.SelfSourceCodeEnabled = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "self source toggle cleanup"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("initial bootstrapAgentRoot failed: %v", err)
	}

	sourceRoot := filepath.Join(cfg.AgentRoot(), "source-code")
	if _, err := os.Stat(sourceRoot); err != nil {
		t.Fatalf("expected self-source root after enabled sync: %v", err)
	}

	cfg.SelfSourceCodeEnabled = false
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("disabled bootstrapAgentRoot failed: %v", err)
	}
	if _, err := os.Stat(sourceRoot); !os.IsNotExist(err) {
		t.Fatalf("expected self-source root removed when disabled, got err=%v", err)
	}
	// The public/source-code FUSE projection follows the retained source tree;
	// its presence/absence is covered by TestBootstrapAgentRootProjectsFusePublicSurface.
}

func TestBootstrapAgentRootRemovesStaleGenomeSurface(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "remove stale genome surface"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	staleGenome := filepath.Join(cfg.AgentRoot(), "genome")
	if err := os.MkdirAll(staleGenome, 0o755); err != nil {
		t.Fatalf("mkdir stale dir %s: %v", staleGenome, err)
	}
	if err := os.WriteFile(filepath.Join(staleGenome, "quine.bundle"), []byte("stale\n"), 0o444); err != nil {
		t.Fatalf("write stale bundle %s: %v", staleGenome, err)
	}
	if err := os.Chmod(staleGenome, 0o555); err != nil {
		t.Fatalf("chmod stale dir %s: %v", staleGenome, err)
	}

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	if _, err := os.Stat(staleGenome); !os.IsNotExist(err) {
		t.Fatalf("expected stale genome root removed, got err=%v", err)
	}
	// The genome surface is no longer projected under public/ (FUSE simply omits
	// the node); there is no separate real public/genome surface to clean.
}

func TestCleanupAgentRootRemovesReadOnlySubtrees(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "cleanup read-only subtree"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}

	readOnlyDir := filepath.Join(cfg.AgentRoot(), "rebuild-src", "internal", "tools")
	if err := os.MkdirAll(readOnlyDir, 0o755); err != nil {
		t.Fatalf("mkdir read-only subtree: %v", err)
	}
	readOnlyFile := filepath.Join(readOnlyDir, "tools.go")
	if err := os.WriteFile(readOnlyFile, []byte("package tools\n"), 0o644); err != nil {
		t.Fatalf("write read-only subtree file: %v", err)
	}
	if err := os.Chmod(readOnlyFile, 0o444); err != nil {
		t.Fatalf("chmod read-only file: %v", err)
	}
	if err := os.Chmod(readOnlyDir, 0o555); err != nil {
		t.Fatalf("chmod read-only dir: %v", err)
	}
	if err := os.Chmod(filepath.Dir(readOnlyDir), 0o555); err != nil {
		t.Fatalf("chmod read-only parent dir: %v", err)
	}
	if err := os.Chmod(filepath.Join(cfg.AgentRoot(), "rebuild-src"), 0o555); err != nil {
		t.Fatalf("chmod read-only root dir: %v", err)
	}

	if err := rt.cleanupAgentRoot(); err != nil {
		t.Fatalf("cleanupAgentRoot failed: %v", err)
	}
	if _, err := os.Stat(cfg.AgentRoot()); !os.IsNotExist(err) {
		t.Fatalf("expected agent root removed after cleanup, got err=%v", err)
	}
}

func TestBootstrapAgentRootProjectsFusePublicSurface(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	cfg := testCfg(t)
	cfg.SelfSourceCodeEnabled = true
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	useRealPublicSurface(rt)
	rt.originalInput = "sync agent root over fuse"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	publicRoot := filepath.Join(cfg.AgentRoot(), "public")
	ctlPath := filepath.Join(publicRoot, "ctl")
	postPath := filepath.Join(ctlPath, "post")
	pokePath := filepath.Join(ctlPath, "poke")
	injectPath := filepath.Join(ctlPath, "inject")
	interruptPath := filepath.Join(ctlPath, "interrupt")
	statusPath := filepath.Join(publicRoot, "status")
	logPath := filepath.Join(publicRoot, "log")
	sourcePath := filepath.Join(publicRoot, "source-code")
	genomePath := filepath.Join(publicRoot, "genome")

	if info, err := os.Lstat(ctlPath); err != nil {
		t.Fatalf("stat public ctl: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("public ctl should be a FUSE directory, got mode %v", info.Mode())
	}
	for _, path := range []string{postPath, pokePath, injectPath, interruptPath} {
		if info, err := os.Lstat(path); err != nil {
			t.Fatalf("stat public control child %s: %v", path, err)
		} else if info.Mode()&os.ModeType != 0 {
			t.Fatalf("public control child %s should be a regular FUSE file, got mode %v", path, info.Mode())
		}
	}
	if info, err := os.Lstat(statusPath); err != nil {
		t.Fatalf("stat public status dir: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("public status should be a FUSE directory, got mode %v", info.Mode())
	}
	for _, name := range []string{"session.json", "inbox.json", "contract.json"} {
		path := filepath.Join(statusPath, name)
		if info, err := os.Lstat(path); err != nil {
			t.Fatalf("stat public status file %s: %v", path, err)
		} else if info.Mode()&os.ModeType != 0 {
			t.Fatalf("public status file %s should be a regular FUSE file, got mode %v", path, info.Mode())
		}
		publicData, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read public status file %s: %v", path, err)
		}
		targetData, err := os.ReadFile(filepath.Join(cfg.AgentRoot(), "status", name))
		if err != nil {
			t.Fatalf("read mirrored status target %s: %v", name, err)
		}
		if string(publicData) != string(targetData) {
			t.Fatalf("public status file %s drifted from target", name)
		}
	}
	if info, err := os.Lstat(logPath); err != nil {
		t.Fatalf("stat public log dir: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("public log should be a FUSE directory, got mode %v", info.Mode())
	}
	controlLogPath := filepath.Join(logPath, "control.jsonl")
	if info, err := os.Lstat(controlLogPath); err != nil {
		t.Fatalf("stat public control log projection: %v", err)
	} else if info.Mode()&os.ModeType != 0 {
		t.Fatalf("public control log should be a regular FUSE file, got mode %v", info.Mode())
	}
	publicControlLog, err := os.ReadFile(controlLogPath)
	if err != nil {
		t.Fatalf("read public control log projection: %v", err)
	}
	targetControlLog, err := os.ReadFile(cfg.SessionControlLogPath(""))
	if err != nil {
		t.Fatalf("read retained control log target: %v", err)
	}
	if string(publicControlLog) != string(targetControlLog) {
		t.Fatal("public control log projection drifted from retained target")
	}
	if _, err := os.Lstat(filepath.Join(logPath, "tapes")); !os.IsNotExist(err) {
		t.Fatalf("public log should not project retained tapes, got err=%v", err)
	}
	if info, err := os.Lstat(sourcePath); err != nil {
		t.Fatalf("stat public source-code dir: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("public source-code should be a FUSE directory, got mode %v", info.Mode())
	}
	publicGoMod, err := os.ReadFile(filepath.Join(sourcePath, "go.mod"))
	if err != nil {
		t.Fatalf("read public FUSE source-code go.mod: %v", err)
	}
	sourceGoMod, err := os.ReadFile(filepath.Join(cfg.AgentRoot(), "source-code", "go.mod"))
	if err != nil {
		t.Fatalf("read source-code go.mod target: %v", err)
	}
	if string(publicGoMod) != string(sourceGoMod) {
		t.Fatalf("public FUSE source-code go.mod drifted from source target")
	}
	if _, err := os.Lstat(genomePath); !os.IsNotExist(err) {
		t.Fatalf("public genome should be absent, got err=%v", err)
	}

	summary, err := os.ReadFile(postPath)
	if err != nil {
		t.Fatalf("read public ctl/post summary: %v", err)
	}
	text := string(summary)
	if !strings.Contains(text, "control_file: post") {
		t.Fatalf("public ctl/post summary missing control file: %q", text)
	}
	if !strings.Contains(text, "mode: queue-only") {
		t.Fatalf("public ctl/post summary missing queue-only mode: %q", text)
	}
	if !strings.Contains(text, "pending_count: 0") {
		t.Fatalf("public ctl/post summary missing pending_count=0: %q", text)
	}

	payload := "hello through fuse public ctl/post"
	writeControlActionFile(t, postPath, payload+"\n")

	waitForCondition(t, 2*time.Second, func() bool {
		data, err := os.ReadFile(cfg.InboxPath())
		if err != nil {
			return false
		}
		var snapshot controlInboxSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return false
		}
		return snapshot.PendingCount == 1
	}, "fuse ctl inbox update")

	inboxData, err := os.ReadFile(cfg.InboxPath())
	if err != nil {
		t.Fatalf("read inbox snapshot: %v", err)
	}
	var inbox controlInboxSnapshot
	if err := json.Unmarshal(inboxData, &inbox); err != nil {
		t.Fatalf("unmarshal inbox snapshot: %v", err)
	}
	if inbox.PendingCount != 1 {
		t.Fatalf("inbox pending_count = %d, want 1", inbox.PendingCount)
	}
	if len(inbox.Messages) != 1 || inbox.Messages[0].Payload != payload {
		t.Fatalf("inbox messages = %#v, want payload %q", inbox.Messages, payload)
	}

	if err := rt.cleanupAgentRoot(); err != nil {
		t.Fatalf("cleanupAgentRoot failed: %v", err)
	}
	if _, err := os.Stat(cfg.AgentRoot()); !os.IsNotExist(err) {
		t.Fatalf("agent root should be removed after cleanup, got err=%v", err)
	}
}

func TestRunSyncsRetentionDir(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
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
	cfg.RetentionDir = filepath.Join(t.TempDir(), "retained")
	rt := NewWithProvider(cfg, mock)

	if exitCode := rt.Run("sync retention", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	retainedRoot := cfg.SessionRetainedDir("")
	missionData, err := os.ReadFile(filepath.Join(retainedRoot, "mission.txt"))
	if err != nil {
		t.Fatalf("read retention mission: %v", err)
	}
	if strings.TrimSpace(string(missionData)) != "sync retention" {
		t.Fatalf("retention mission = %q, want %q", strings.TrimSpace(string(missionData)), "sync retention")
	}

	sessionData, err := os.ReadFile(filepath.Join(retainedRoot, "status", "session.json"))
	if err != nil {
		t.Fatalf("read retention session metadata: %v", err)
	}
	var session map[string]any
	if err := json.Unmarshal(sessionData, &session); err != nil {
		t.Fatalf("unmarshal retention session metadata: %v", err)
	}
	if session["session_id"] != cfg.SessionID {
		t.Fatalf("retention session_id = %v, want %q", session["session_id"], cfg.SessionID)
	}
	if session["incarnation_id"] != float64(0) {
		t.Fatalf("retention incarnation_id = %v, want 0", session["incarnation_id"])
	}
	if _, ok := session["tape_id"]; ok {
		t.Fatalf("retention session metadata should not expose tape_id, got %v", session["tape_id"])
	}

	runtimeLogData, err := os.ReadFile(filepath.Join(retainedRoot, "runtime.log"))
	if err != nil {
		t.Fatalf("read retention runtime.log: %v", err)
	}
	if !strings.Contains(string(runtimeLogData), "session started") {
		t.Fatalf("retention runtime.log missing session start entry: %q", string(runtimeLogData))
	}

	tapeJSONLPath := filepath.Join(retainedRoot, "tapes", cfg.TapeID+".jsonl")
	tapeJSONLData, err := os.ReadFile(tapeJSONLPath)
	if err != nil {
		t.Fatalf("read retained tape jsonl: %v", err)
	}
	if !strings.Contains(string(tapeJSONLData), `"type":"outcome"`) {
		t.Fatalf("retained tape jsonl missing outcome entry: %q", string(tapeJSONLData))
	}

	tapeYAMLPath := filepath.Join(retainedRoot, "tapes", cfg.TapeID+".log.yaml")
	if _, err := os.Stat(tapeYAMLPath); err != nil {
		t.Fatalf("stat retained tape yaml: %v", err)
	}
	assertSymlinkTarget(t, filepath.Join(retainedRoot, "inc", "current"), "0")
	incMissionData, err := os.ReadFile(filepath.Join(retainedRoot, "inc", "0", "mission.txt"))
	if err != nil {
		t.Fatalf("read retention incarnation mission: %v", err)
	}
	if strings.TrimSpace(string(incMissionData)) != "sync retention" {
		t.Fatalf("retention incarnation mission = %q, want %q", strings.TrimSpace(string(incMissionData)), "sync retention")
	}

	compatLogPath := filepath.Join(cfg.RuntimeRoot(), "log", cfg.SessionID)
	assertRelativeSymlinkResolvesTo(t, compatLogPath, retainedRoot)
	compatTapesPath := filepath.Join(cfg.RuntimeRoot(), "tapes", cfg.SessionID)
	assertRelativeSymlinkResolvesTo(t, compatTapesPath, filepath.Join(retainedRoot, "tapes"))
}

func TestReplaceSymlinkResolvesRelativeTarget(t *testing.T) {
	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir tmp: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	linkPath := filepath.Join(".quine", "agent", "session-1", "jobs")
	targetPath := filepath.Join(".quine", "jobs", "session-1")
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		t.Fatalf("mkdir link parent: %v", err)
	}
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	if err := replaceSymlink(linkPath, targetPath); err != nil {
		t.Fatalf("replaceSymlink: %v", err)
	}

	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("eval symlink: %v", err)
	}
	want, err := filepath.Abs(targetPath)
	if err != nil {
		t.Fatalf("abs target: %v", err)
	}
	want, err = filepath.EvalSymlinks(want)
	if err != nil {
		t.Fatalf("eval target symlink: %v", err)
	}
	if resolved != want {
		t.Fatalf("resolved symlink = %q, want %q", resolved, want)
	}
}

func TestReplaceSymlinkConcurrentSameTarget(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	linkPath := filepath.Join(tmpDir, "pid", "12345")
	targetPath := filepath.Join(tmpDir, "agent", "session-1", "public")
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		t.Fatalf("mkdir link parent: %v", err)
	}
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 64)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 16; j++ {
				if err := replaceSymlink(linkPath, targetPath); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("replaceSymlink concurrent same target: %v", err)
	}

	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("eval symlink: %v", err)
	}
	want, err := filepath.Abs(targetPath)
	if err != nil {
		t.Fatalf("abs target: %v", err)
	}
	want, err = filepath.EvalSymlinks(want)
	if err != nil {
		t.Fatalf("eval target symlink: %v", err)
	}
	if resolved != want {
		t.Fatalf("resolved symlink = %q, want %q", resolved, want)
	}
}
