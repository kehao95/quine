package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

func TestRunResumesCurrentIncarnationInPlace(t *testing.T) {
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
	priorContextRoot := filepath.Join(cfg.SessionIncarnationPath("", 0), "context")
	currentPath := filepath.Join(priorContextRoot, "state", "current.jsonl")
	writeJSONLEntries(t, currentPath, tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "old incarnation"}))
	if err := os.MkdirAll(filepath.Join(priorContextRoot, "prompt"), 0o755); err != nil {
		t.Fatalf("mkdir prior prompt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(priorContextRoot, "prompt", "30-memory.md"), []byte("old incarnation memory\n"), 0o644); err != nil {
		t.Fatalf("write prior memory: %v", err)
	}

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("continued mission", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected provider to be called")
	}
	if got := len(mock.calls[0]); got != 2 {
		t.Fatalf("first provider call len = %d, want 2 (system + preexisting context)", got)
	}
	foundExisting := false
	foundMemory := false
	foundBegin := false
	for _, msg := range mock.calls[0] {
		if msg.Content == "old incarnation" {
			foundExisting = true
		}
		if strings.Contains(msg.Content, "old incarnation memory") {
			foundMemory = true
		}
		if msg.Content == "Begin." {
			foundBegin = true
		}
	}
	if !foundExisting {
		t.Fatal("new incarnation should inherit prior context/state/current.jsonl content")
	}
	if !foundMemory {
		t.Fatal("new incarnation should inherit prior memory.md content")
	}
	if foundBegin {
		t.Fatal("new incarnation should not reseed Begin. when copied context is present")
	}
	if cfg.IncarnationID != 0 {
		t.Fatalf("cfg.IncarnationID = %d, want resumed incarnation 0", cfg.IncarnationID)
	}
	memoryPath := filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "prompt", "30-memory.md")
	memoryData, err := os.ReadFile(memoryPath)
	if err != nil {
		t.Fatalf("read copied memory: %v", err)
	}
	if string(memoryData) != "old incarnation memory\n" {
		t.Fatalf("copied memory = %q, want prior memory", string(memoryData))
	}
}

func TestCurrentSessionIncarnationIDFallbacks(t *testing.T) {
	t.Run("current symlink", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "inc", "2"), 0o755); err != nil {
			t.Fatalf("mkdir inc: %v", err)
		}
		if err := os.Symlink("2", filepath.Join(root, "inc", "current")); err != nil {
			t.Fatalf("symlink current: %v", err)
		}
		id, ok, err := currentSessionIncarnationID(root)
		if err != nil {
			t.Fatalf("currentSessionIncarnationID: %v", err)
		}
		if !ok || id != 2 {
			t.Fatalf("id=%d ok=%v, want 2 true", id, ok)
		}
	})

	t.Run("status fallback", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "status"), 0o755); err != nil {
			t.Fatalf("mkdir status: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "status", "session.json"), []byte(`{"incarnation_id":4}`), 0o644); err != nil {
			t.Fatalf("write status: %v", err)
		}
		id, ok, err := currentSessionIncarnationID(root)
		if err != nil {
			t.Fatalf("currentSessionIncarnationID: %v", err)
		}
		if !ok || id != 4 {
			t.Fatalf("id=%d ok=%v, want 4 true", id, ok)
		}
	})

	t.Run("max incarnation fallback", func(t *testing.T) {
		root := t.TempDir()
		for _, id := range []string{"0", "3", "1"} {
			if err := os.MkdirAll(filepath.Join(root, "inc", id), 0o755); err != nil {
				t.Fatalf("mkdir inc %s: %v", id, err)
			}
		}
		id, ok, err := currentSessionIncarnationID(root)
		if err != nil {
			t.Fatalf("currentSessionIncarnationID: %v", err)
		}
		if !ok || id != 3 {
			t.Fatalf("id=%d ok=%v, want 3 true", id, ok)
		}
	})

	t.Run("empty", func(t *testing.T) {
		id, ok, err := currentSessionIncarnationID(t.TempDir())
		if err != nil {
			t.Fatalf("currentSessionIncarnationID: %v", err)
		}
		if ok || id != 0 {
			t.Fatalf("id=%d ok=%v, want 0 false", id, ok)
		}
	})
}

func TestRunImportsCopiedAgentSurfaceForNewSessionID(t *testing.T) {
	sourceCfg := testCfg(t)
	sourceCfg.SessionID = "source-session"
	sourceCfg.RunID = "run-source"
	sourceCfg.TapeID = "tape-source"
	sourceRT := NewWithProvider(sourceCfg, &mockProvider{})
	silenceRuntime(sourceRT)
	sourceRT.originalInput = "source mission"
	if err := sourceRT.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrap source agent root: %v", err)
	}

	sourceContextPath := filepath.Join(sourceCfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl")
	writeJSONLEntries(t, sourceContextPath,
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "copied agent context"}),
	)
	sourceMemoryPath := filepath.Join(sourceCfg.SessionIncarnationPath("", 0), "context", "prompt", "30-memory.md")
	if err := os.WriteFile(sourceMemoryPath, []byte("copied agent memory\n"), 0o644); err != nil {
		t.Fatalf("write source memory: %v", err)
	}

	targetCfg := *sourceCfg
	targetCfg.SessionID = "copied-session"
	targetCfg.RunID = "run-copied"
	targetCfg.TapeID = "tape-copied"
	targetAgentRoot := targetCfg.AgentRoot()
	copyAgentRootForResumeTest(t, sourceCfg.AgentRoot(), targetAgentRoot)

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:        "call_exit",
						Name:      "exit",
						Arguments: map[string]any{"status": "success"},
					},
				},
			},
		},
	}
	rt := NewWithProvider(&targetCfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("copied mission", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected provider to be called")
	}
	var foundContext, foundMemory, foundBegin bool
	for _, msg := range mock.calls[0] {
		if msg.Content == "copied agent context" {
			foundContext = true
		}
		if strings.Contains(msg.Content, "copied agent memory") {
			foundMemory = true
		}
		if msg.Content == "Begin." {
			foundBegin = true
		}
	}
	if !foundContext {
		t.Fatalf("provider input missing copied context: %#v", mock.calls[0])
	}
	if !foundMemory {
		t.Fatalf("provider input missing copied memory: %#v", mock.calls[0])
	}
	if foundBegin {
		t.Fatalf("copied agent context should suppress Begin. reseed: %#v", mock.calls[0])
	}
	if targetCfg.IncarnationID != 0 {
		t.Fatalf("target incarnation = %d, want 0", targetCfg.IncarnationID)
	}
	importedContextPath := filepath.Join(targetCfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl")
	importedContext, err := os.ReadFile(importedContextPath)
	if err != nil {
		t.Fatalf("read imported context: %v", err)
	}
	if !strings.Contains(string(importedContext), "copied agent context") {
		t.Fatalf("imported context = %q, want copied context", string(importedContext))
	}
	statusData, err := os.ReadFile(filepath.Join(targetCfg.SessionRetainedDir(""), "status", "session.json"))
	if err != nil {
		t.Fatalf("read target status: %v", err)
	}
	var status map[string]any
	if err := json.Unmarshal(statusData, &status); err != nil {
		t.Fatalf("unmarshal target status: %v", err)
	}
	if status["session_id"] != targetCfg.SessionID {
		t.Fatalf("status session_id = %v, want %q", status["session_id"], targetCfg.SessionID)
	}
	if status["run_id"] != targetCfg.RunID {
		t.Fatalf("status run_id = %v, want %q", status["run_id"], targetCfg.RunID)
	}
	if _, err := os.Stat(filepath.Join(targetCfg.SessionRetainedDir(""), ".resume")); !os.IsNotExist(err) {
		t.Fatalf("resume must not create .resume checkpoint, stat err=%v", err)
	}
}

func TestRunImportsCopiedAgentSurfaceWithPendingToolBatch(t *testing.T) {
	sourceCfg := testCfg(t)
	sourceCfg.SessionID = "source-pending-session"
	sourceCfg.RunID = "run-source-pending"
	sourceCfg.TapeID = "tape-source-pending"
	sourceRT := NewWithProvider(sourceCfg, &mockProvider{})
	silenceRuntime(sourceRT)
	sourceRT.originalInput = "source pending mission"
	if err := sourceRT.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrap source agent root: %v", err)
	}
	seedPendingToolBatch(t, sourceCfg, tape.ToolCall{
		ID:        "call_exec",
		Name:      "exec",
		Arguments: map[string]any{},
	})

	targetCfg := *sourceCfg
	targetCfg.SessionID = "copied-pending-session"
	targetCfg.RunID = "run-copied-pending"
	targetCfg.TapeID = "tape-copied-pending"
	copyAgentRootForResumeTest(t, sourceCfg.AgentRoot(), targetCfg.AgentRoot())

	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(&targetCfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("copied pending mission", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected provider to be called")
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_exec")
	if got := toolString(t, payload, "status"); got != "unknown" {
		t.Fatalf("copied pending exec status = %q, want unknown; payload=%#v", got, payload)
	}
	if _, err := os.Stat(filepath.Join(targetCfg.SessionRetainedDir(""), ".resume")); !os.IsNotExist(err) {
		t.Fatalf("resume must not create .resume checkpoint, stat err=%v", err)
	}
}

func TestRunDoesNotImportCopiedAgentSurfaceWhenRetainedStateExists(t *testing.T) {
	sourceCfg := testCfg(t)
	sourceCfg.SessionID = "source-session"
	sourceCfg.RunID = "run-source"
	sourceCfg.TapeID = "tape-source"
	sourceRT := NewWithProvider(sourceCfg, &mockProvider{})
	silenceRuntime(sourceRT)
	sourceRT.originalInput = "source mission"
	if err := sourceRT.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrap source agent root: %v", err)
	}
	writeJSONLEntries(t, filepath.Join(sourceCfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "copied agent context"}),
	)

	targetCfg := *sourceCfg
	targetCfg.SessionID = "retained-target"
	targetCfg.RunID = "run-retained"
	targetCfg.TapeID = "tape-retained"
	copyAgentRootForResumeTest(t, sourceCfg.AgentRoot(), targetCfg.AgentRoot())
	writeJSONLEntries(t, filepath.Join(targetCfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "retained target context"}),
	)

	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(&targetCfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("retained target mission", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected provider to be called")
	}
	var foundRetained, foundCopied bool
	for _, msg := range mock.calls[0] {
		if msg.Content == "retained target context" {
			foundRetained = true
		}
		if msg.Content == "copied agent context" {
			foundCopied = true
		}
	}
	if !foundRetained {
		t.Fatalf("provider input missing retained target context: %#v", mock.calls[0])
	}
	if foundCopied {
		t.Fatalf("provider input should not import copied context over retained state: %#v", mock.calls[0])
	}
}

func TestRunImportsBootstrappedContextBeforeSeeding(t *testing.T) {
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
	cfg.AnchorMemoryEnabled = true
	bootstrapRoot := filepath.Join(t.TempDir(), "bootstrap-context")
	stateRoot := filepath.Join(bootstrapRoot, "state")
	frontierDir := filepath.Join(stateRoot, "frontier")
	anchorDir := filepath.Join(stateRoot, "anchors", "2.anchor")
	if err := os.MkdirAll(frontierDir, 0o755); err != nil {
		t.Fatalf("mkdir frontier: %v", err)
	}
	if err := os.MkdirAll(anchorDir, 0o755); err != nil {
		t.Fatalf("mkdir anchor dir: %v", err)
	}
	writeJSONLEntries(t, filepath.Join(stateRoot, "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "bootstrapped current"}),
	)
	if err := os.WriteFile(filepath.Join(frontierDir, "2.link"), []byte(""), 0o644); err != nil {
		t.Fatalf("write frontier link: %v", err)
	}
	if err := os.WriteFile(filepath.Join(anchorDir, "meta.json"), []byte(`{"id":2,"resolution":"bootstrapped anchor"}`), 0o644); err != nil {
		t.Fatalf("write anchor meta: %v", err)
	}
	t.Setenv(tools.ContextBootstrapEnv, bootstrapRoot)

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("child mission", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected provider to be called")
	}
	var foundCurrent, foundAnchor, foundBegin bool
	for _, msg := range mock.calls[0] {
		switch msg.Content {
		case "bootstrapped current":
			foundCurrent = true
		case "[ANCHOR 2] bootstrapped anchor":
			foundAnchor = true
		case "Begin.":
			foundBegin = true
		}
	}
	if !foundCurrent {
		t.Fatal("provider input should include bootstrapped context/state/current.jsonl content")
	}
	if !foundAnchor {
		t.Fatal("provider input should include imported frontier anchor summary")
	}
	if foundBegin {
		t.Fatal("runtime should not reseed Begin. when bootstrapped context is present")
	}
}

func TestRunImportsRetainedSeedContextForFirstIncarnation(t *testing.T) {
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
	cfg.AnchorMemoryEnabled = true
	seedRoot := filepath.Join(cfg.SessionRetainedDir(""), "seed", "context")
	stateRoot := filepath.Join(seedRoot, "state")
	frontierDir := filepath.Join(stateRoot, "frontier")
	anchorDir := filepath.Join(stateRoot, "anchors", "3.anchor")
	if err := os.MkdirAll(frontierDir, 0o755); err != nil {
		t.Fatalf("mkdir seed frontier: %v", err)
	}
	if err := os.MkdirAll(anchorDir, 0o755); err != nil {
		t.Fatalf("mkdir seed anchor dir: %v", err)
	}
	writeJSONLEntries(t, filepath.Join(stateRoot, "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "seeded current"}),
	)
	if err := os.WriteFile(filepath.Join(frontierDir, "3.link"), []byte(""), 0o644); err != nil {
		t.Fatalf("write seed frontier link: %v", err)
	}
	if err := os.WriteFile(filepath.Join(anchorDir, "meta.json"), []byte(`{"id":3,"resolution":"seeded anchor"}`), 0o644); err != nil {
		t.Fatalf("write seed anchor meta: %v", err)
	}

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("seeded mission", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected provider to be called")
	}
	var foundCurrent, foundAnchor, foundBegin bool
	for _, msg := range mock.calls[0] {
		switch msg.Content {
		case "seeded current":
			foundCurrent = true
		case "[ANCHOR 3] seeded anchor":
			foundAnchor = true
		case "Begin.":
			foundBegin = true
		}
	}
	if !foundCurrent {
		t.Fatal("provider input should include retained seed context/state/current.jsonl content")
	}
	if !foundAnchor {
		t.Fatal("provider input should include retained seed frontier anchor summary")
	}
	if foundBegin {
		t.Fatal("runtime should not reseed Begin. when retained seed context is present")
	}
}

func TestResumeSplicesRetainedToolResultFromTape(t *testing.T) {
	cfg := testCfg(t)
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:        "call_sh",
		Name:      "sh",
		Arguments: map[string]any{"command": "printf retained"},
	})
	writeJSONLEntries(t, filepath.Join(cfg.TapeDir(""), "prior.jsonl"),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID: "call_sh",
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":      "sh",
				"mode":      "sync",
				"status":    "completed",
				"exit_code": 0,
				"stdout":    "retained stdout",
				"stderr":    "",
			}),
		}),
	)
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume retained tape result", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("provider was not called")
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
	if got := toolString(t, payload, "status"); got != "completed" {
		t.Fatalf("status = %q, want completed; payload=%#v", got, payload)
	}
	if got := toolString(t, payload, "stdout"); got != "retained stdout" {
		t.Fatalf("stdout = %q, want retained stdout", got)
	}
}

func TestResumeSplicesNewestRetainedToolResultFromTape(t *testing.T) {
	cfg := testCfg(t)
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:        "call_sh",
		Name:      "sh",
		Arguments: map[string]any{"command": "printf retained"},
	})
	writeJSONLEntries(t, filepath.Join(cfg.TapeDir(""), "001-old.jsonl"),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID: "call_sh",
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":   "sh",
				"status": "completed",
				"stdout": "old retained stdout",
			}),
		}),
	)
	writeJSONLEntries(t, filepath.Join(cfg.TapeDir(""), "999-new.jsonl"),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID: "call_sh",
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":   "sh",
				"status": "completed",
				"stdout": "new retained stdout",
			}),
		}),
	)
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume newest retained tape result", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
	if got := toolString(t, payload, "stdout"); got != "new retained stdout" {
		t.Fatalf("stdout = %q, want newest retained stdout", got)
	}
}

func TestResumeCompletesPartialMultiToolBatch(t *testing.T) {
	cfg := testCfg(t)
	writeJSONLEntries(t, filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "before partial batch"}),
		tape.MessageEntry(tape.Message{
			Role: tape.RoleAssistant,
			ToolCalls: []tape.ToolCall{
				{ID: "call_done", Name: "sh", Arguments: map[string]any{"command": "printf done"}},
				{ID: "call_missing", Name: "exec", Arguments: map[string]any{}},
			},
		}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID: "call_done",
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":   "sh",
				"status": "completed",
				"stdout": "already done",
			}),
		}),
	)
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume partial batch", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("provider was not called")
	}
	if count := countToolResultsInMessages(mock.calls[0], "call_done"); count != 1 {
		t.Fatalf("call_done result count = %d, want 1", count)
	}
	donePayload := requireToolResultInMessages(t, mock.calls[0], "call_done")
	if got := toolString(t, donePayload, "stdout"); got != "already done" {
		t.Fatalf("existing result stdout = %q, want already done", got)
	}
	missingPayload := requireToolResultInMessages(t, mock.calls[0], "call_missing")
	if got := toolString(t, missingPayload, "status"); got != "unknown" {
		t.Fatalf("missing result status = %q, want unknown; payload=%#v", got, missingPayload)
	}
}

func TestResumeLeavesCompleteToolBatchUnchanged(t *testing.T) {
	cfg := testCfg(t)
	writeJSONLEntries(t, filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "state", "current.jsonl"),
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "before complete batch"}),
		tape.MessageEntry(tape.Message{
			Role:      tape.RoleAssistant,
			ToolCalls: []tape.ToolCall{{ID: "call_done", Name: "sh", Arguments: map[string]any{"command": "printf done"}}},
		}),
		tape.ToolResultEntry(tape.ToolResult{
			ToolID: "call_done",
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":   "sh",
				"status": "completed",
				"stdout": "done once",
			}),
		}),
	)
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume complete batch", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if len(mock.calls) == 0 {
		t.Fatal("provider was not called")
	}
	if count := countToolResultsInMessages(mock.calls[0], "call_done"); count != 1 {
		t.Fatalf("call_done result count = %d, want 1", count)
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_done")
	if got := toolString(t, payload, "stdout"); got != "done once" {
		t.Fatalf("stdout = %q, want done once", got)
	}
}

func TestResumeDisabledPendingToolsReturnUnknown(t *testing.T) {
	tests := []struct {
		name   string
		cfg    func(*config.Config)
		tc     tape.ToolCall
		reason string
	}{
		{
			name: "exit disabled",
			cfg: func(cfg *config.Config) {
				cfg.ToolGates = config.DefaultToolGates()
				cfg.ToolGates.ExitEnabled = false
			},
			tc:     tape.ToolCall{ID: "call_exit", Name: "exit", Arguments: map[string]any{"status": "success"}},
			reason: "exit is disabled",
		},
		{
			name:   "idle disabled",
			cfg:    func(cfg *config.Config) { cfg.IdleEnabled = false },
			tc:     tape.ToolCall{ID: "call_idle", Name: "idle", Arguments: map[string]any{}},
			reason: "idle is disabled",
		},
		{
			name:   "vision disabled",
			cfg:    func(cfg *config.Config) { cfg.VisionEnabled = false },
			tc:     tape.ToolCall{ID: "call_vision", Name: "vision", Arguments: map[string]any{"path": "missing.png"}},
			reason: "vision is disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			tt.cfg(cfg)
			seedPendingToolBatch(t, cfg, tt.tc)
			rt := NewWithProvider(cfg, &mockProvider{})
			silenceRuntime(rt)
			rt.originalInput = "recover disabled pending tool"
			rt.tape = tape.NewTape(cfg.SessionID, cfg.ParentSession, cfg.Depth, cfg.ModelID)
			if err := rt.bootstrapAgentRoot(); err != nil {
				t.Fatalf("bootstrapAgentRoot: %v", err)
			}

			if code, done := rt.recoverIncompleteToolBatches(); done {
				t.Fatalf("recoverIncompleteToolBatches done=%v code=%d, want continue", done, code)
			}
			msgs, err := rt.providerContextMessages()
			if err != nil {
				t.Fatalf("providerContextMessages: %v", err)
			}
			payload := requireToolResultInMessages(t, msgs, tt.tc.ID)
			if got := toolString(t, payload, "status"); got != "unknown" {
				t.Fatalf("status = %q, want unknown; payload=%#v", got, payload)
			}
			if got := toolString(t, payload, "error"); !strings.Contains(got, tt.reason) {
				t.Fatalf("error = %q, want contains %q", got, tt.reason)
			}
		})
	}
}

func TestResumeCompletesOpaquePendingToolsAsUnknown(t *testing.T) {
	tests := []struct {
		name string
		tc   tape.ToolCall
	}{
		{
			name: "sh",
			tc: tape.ToolCall{ID: "call_sh", Name: "sh", Arguments: map[string]any{
				"command": "printf side-effect",
			}},
		},
		{
			name: "fork",
			tc: tape.ToolCall{ID: "call_fork", Name: "fork", Arguments: map[string]any{
				"children": []any{map[string]any{"intent": "child", "scope": "."}},
			}},
		},
		{
			name: "spawn",
			tc: tape.ToolCall{ID: "call_spawn", Name: "spawn", Arguments: map[string]any{
				"children": []any{map[string]any{"mission": "child"}},
			}},
		},
		{
			name: "exec",
			tc:   tape.ToolCall{ID: "call_exec", Name: "exec", Arguments: map[string]any{}},
		},
		{
			name: "mark",
			tc: tape.ToolCall{ID: "call_mark", Name: "mark", Arguments: map[string]any{
				"resolution": "opaque mark",
			}},
		},
		{
			name: "unknown",
			tc:   tape.ToolCall{ID: "call_mystery", Name: "mystery", Arguments: map[string]any{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			seedPendingToolBatch(t, cfg, tt.tc)
			mock := &mockProvider{responses: []tape.Message{{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{{
					ID:        "call_exit",
					Name:      "exit",
					Arguments: map[string]any{"status": "success"},
				}},
			}}}
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume opaque pending tool", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			if len(mock.calls) == 0 {
				t.Fatal("provider was not called")
			}
			payload := requireToolResultInMessages(t, mock.calls[0], tt.tc.ID)
			if got := toolString(t, payload, "status"); got != "unknown" {
				t.Fatalf("status = %q, want unknown; payload=%#v", got, payload)
			}
			if got := toolString(t, payload, "side_effects"); got != "unknown" {
				t.Fatalf("side_effects = %q, want unknown", got)
			}
		})
	}
}

func TestResumeReplaysSafePendingTools(t *testing.T) {
	tests := []struct {
		name string
		cfg  func(*config.Config)
		tc   tape.ToolCall
	}{
		{
			name: "vision",
			tc: tape.ToolCall{ID: "call_vision", Name: "vision", Arguments: map[string]any{
				"path": filepath.Join(t.TempDir(), "missing.png"),
			}},
		},
		{
			name: "unfold",
			cfg:  func(cfg *config.Config) { cfg.AnchorMemoryEnabled = true },
			tc: tape.ToolCall{ID: "call_unfold", Name: "unfold", Arguments: map[string]any{
				"anchor_id": float64(99),
			}},
		},
		{
			name: "switch_world",
			tc: tape.ToolCall{ID: "call_switch", Name: "switch_world", Arguments: map[string]any{
				"target": "wr0",
			}},
		},
		{
			name: "escalate",
			tc: tape.ToolCall{ID: "call_escalate", Name: "escalate", Arguments: map[string]any{
				"reason": "recover pending escalation",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			if tt.cfg != nil {
				tt.cfg(cfg)
			}
			seedPendingToolBatch(t, cfg, tt.tc)
			mock := &mockProvider{responses: []tape.Message{{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{{
					ID:        "call_exit",
					Name:      "exit",
					Arguments: map[string]any{"status": "success"},
				}},
			}}}
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume safe pending tool", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			payload := requireToolResultInMessages(t, mock.calls[0], tt.tc.ID)
			if got := toolString(t, payload, "tool"); got != tt.tc.Name {
				t.Fatalf("tool = %q, want %q", got, tt.tc.Name)
			}
			if got := toolString(t, payload, "status"); got == "unknown" {
				t.Fatalf("safe replay returned unknown: %#v", payload)
			}
		})
	}
}

func TestResumeReexecutesPendingExitWithoutProvider(t *testing.T) {
	cfg := testCfg(t)
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:        "call_exit",
		Name:      "exit",
		Arguments: map[string]any{"status": "success"},
	})
	mock := &mockProvider{}
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume pending exit", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if mock.callCount != 0 {
		t.Fatalf("provider calls = %d, want 0", mock.callCount)
	}
}

func TestResumeRestoresPendingIdleBeforeProvider(t *testing.T) {
	cfg := testCfg(t)
	cfg.IdleEnabled = true
	seedPendingToolBatch(t, cfg, tape.ToolCall{ID: "call_idle", Name: "idle", Arguments: map[string]any{}})
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("resume pending idle", "Begin.")
	}()
	waitForCondition(t, time.Second, func() bool {
		if rt.control == nil {
			return false
		}
		rt.control.mu.Lock()
		defer rt.control.mu.Unlock()
		return rt.control.started
	}, "control surface ready")
	rt.requestControlPoke()

	select {
	case exitCode := <-done:
		if exitCode != 0 {
			t.Fatalf("exit code = %d, want 0", exitCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runtime did not resume from recovered idle")
	}
	if len(mock.calls) == 0 {
		t.Fatal("provider was not called after idle resumed")
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_idle")
	if got := toolString(t, payload, "status"); got != "completed" {
		t.Fatalf("idle status = %q, want completed; payload=%#v", got, payload)
	}
	if got := toolString(t, payload, "delivery"); got != string(controlDeliveryPoke) {
		t.Fatalf("idle delivery = %q, want poke", got)
	}
}

func TestResumeIdlePreservesRetainedInboxState(t *testing.T) {
	cfg := testCfg(t)
	cfg.IdleEnabled = true
	seedPendingToolBatch(t, cfg, tape.ToolCall{ID: "call_idle", Name: "idle", Arguments: map[string]any{}})
	snapshot := controlInboxSnapshot{
		PendingCount: 1,
		Messages: []controlInboxMessage{{
			ID:         "message-7",
			Payload:    "queued before crash",
			ReceivedAt: 123,
		}},
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(cfg.SessionRetainedDir(""), "status", "inbox.json")), 0o755); err != nil {
		t.Fatalf("mkdir retained status: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfg.SessionRetainedDir(""), "status", "inbox.json"), mustJSON(t, snapshot), 0o644); err != nil {
		t.Fatalf("write retained inbox: %v", err)
	}
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("resume pending idle with retained inbox", "Begin.")
	}()
	waitForCondition(t, time.Second, func() bool {
		if rt.control == nil {
			return false
		}
		rt.control.mu.Lock()
		defer rt.control.mu.Unlock()
		return rt.control.started
	}, "control surface ready")
	rt.requestControlPoke()

	select {
	case exitCode := <-done:
		if exitCode != 0 {
			t.Fatalf("exit code = %d, want 0", exitCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runtime did not resume from recovered idle")
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_idle")
	if got := toolString(t, payload, "delivery"); got != string(controlDeliveryPoke) {
		t.Fatalf("idle delivery = %q, want poke", got)
	}
	if got := toolInt(t, payload, "pending_count"); got != 1 {
		t.Fatalf("idle pending_count = %d, want retained inbox count 1", got)
	}
}

func TestLoadRetainedControlStateContinuesMessageIDs(t *testing.T) {
	cfg := testCfg(t)
	snapshot := controlInboxSnapshot{
		PendingCount: 1,
		Messages: []controlInboxMessage{{
			ID:         "message-7",
			Payload:    "queued before crash",
			ReceivedAt: 123,
		}},
	}
	if err := os.MkdirAll(filepath.Join(cfg.SessionRetainedDir(""), "status"), 0o755); err != nil {
		t.Fatalf("mkdir retained status: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfg.SessionRetainedDir(""), "status", "inbox.json"), mustJSON(t, snapshot), 0o644); err != nil {
		t.Fatalf("write retained inbox: %v", err)
	}
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	if err := rt.loadRetainedControlState(); err != nil {
		t.Fatalf("loadRetainedControlState: %v", err)
	}
	queued, err := rt.enqueueControlPayload(controlActionPost, "after resume")
	if err != nil {
		t.Fatalf("enqueueControlPayload: %v", err)
	}
	if queued.ID != "message-8" {
		t.Fatalf("queued id = %q, want message-8", queued.ID)
	}
}

func TestResumeGathersPendingShResultFromJobSurface(t *testing.T) {
	cfg := testCfg(t)
	const command = "printf gathered"
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:        "call_sh",
		Name:      "sh",
		Arguments: map[string]any{"command": command},
	})
	writeJobSurface(t, cfg, "12345", map[string]string{
		"cmd":        command,
		"pid":        "12345\n",
		"out.log":    "gathered stdout",
		"err.log":    "",
		"exit":       "0\n",
		"started_at": time.Now().UTC().Format(time.RFC3339Nano) + "\n",
	})
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume sh job result", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
	if got := toolString(t, payload, "status"); got != "completed" {
		t.Fatalf("sh status = %q, want completed; payload=%#v", got, payload)
	}
	if got := toolString(t, payload, "stdout"); got != "gathered stdout" {
		t.Fatalf("stdout = %q, want gathered stdout", got)
	}
}

func TestResumeGathersPendingShNonZeroExitFromJobSurface(t *testing.T) {
	cfg := testCfg(t)
	const command = "exit 7"
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:        "call_sh",
		Name:      "sh",
		Arguments: map[string]any{"command": command},
	})
	writeJobSurface(t, cfg, "777", map[string]string{
		"cmd":     command,
		"pid":     "777\n",
		"out.log": "",
		"err.log": "failed",
		"exit":    "7\n",
	})
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume sh nonzero job result", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
	if got := toolString(t, payload, "status"); got != "completed" {
		t.Fatalf("status = %q, want completed; payload=%#v", got, payload)
	}
	if got := toolInt(t, payload, "exit_code"); got != 7 {
		t.Fatalf("exit_code = %d, want 7", got)
	}
	if got := toolString(t, payload, "stderr"); got != "failed" {
		t.Fatalf("stderr = %q, want failed", got)
	}
}

func TestResumeUnknownForAmbiguousOrInvalidShJobSurface(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*testing.T, *config.Config)
	}{
		{
			name: "ambiguous matching jobs",
			setup: func(t *testing.T, cfg *config.Config) {
				writeJobSurface(t, cfg, "101", map[string]string{"cmd": "printf duplicate", "out.log": "first"})
				writeJobSurface(t, cfg, "102", map[string]string{"cmd": "printf duplicate", "out.log": "second"})
			},
		},
		{
			name: "invalid exit code",
			setup: func(t *testing.T, cfg *config.Config) {
				writeJobSurface(t, cfg, "201", map[string]string{
					"cmd":     "printf duplicate",
					"out.log": "stdout",
					"err.log": "",
					"exit":    "not-an-int\n",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			seedPendingToolBatch(t, cfg, tape.ToolCall{
				ID:        "call_sh",
				Name:      "sh",
				Arguments: map[string]any{"command": "printf duplicate"},
			})
			tt.setup(t, cfg)
			mock := exitAfterRecoveryProvider()
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume bad sh surface", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
			if got := toolString(t, payload, "status"); got != "unknown" {
				t.Fatalf("status = %q, want unknown; payload=%#v", got, payload)
			}
			if got := toolString(t, payload, "side_effects"); got != "unknown" {
				t.Fatalf("side_effects = %q, want unknown", got)
			}
		})
	}
}

func TestResumeUnknownForEmptyShCommand(t *testing.T) {
	cfg := testCfg(t)
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:        "call_sh",
		Name:      "sh",
		Arguments: map[string]any{"command": "   "},
	})
	mock := exitAfterRecoveryProvider()
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume empty sh command", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
	if got := toolString(t, payload, "status"); got != "unknown" {
		t.Fatalf("status = %q, want unknown; payload=%#v", got, payload)
	}
}

func TestResumeGathersPendingShInterruptedAndSpawnedJobSurfaces(t *testing.T) {
	tests := []struct {
		name       string
		args       map[string]any
		job        map[string]string
		wantStatus string
		wantMode   string
		wantField  string
	}{
		{
			name:       "sync interrupted",
			args:       map[string]any{"command": "printf partial"},
			job:        map[string]string{"cmd": "printf partial", "pid": "111\n", "out.log": "partial out", "err.log": "partial err"},
			wantStatus: "interrupted",
			wantMode:   "sync",
			wantField:  "stdout_so_far",
		},
		{
			name:       "interactive spawned",
			args:       map[string]any{"command": "bash", "interactive": true},
			job:        map[string]string{"cmd": "bash", "pid": "222\n", "out.log": "", "err.log": ""},
			wantStatus: "spawned",
			wantMode:   "interactive",
			wantField:  "interactive",
		},
		{
			name:       "detached spawned",
			args:       map[string]any{"command": "sleep 10", "detach": true},
			job:        map[string]string{"cmd": "sleep 10", "pid": "333\n", "out.log": "", "err.log": ""},
			wantStatus: "spawned",
			wantMode:   "detached",
			wantField:  "detached",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			seedPendingToolBatch(t, cfg, tape.ToolCall{ID: "call_sh", Name: "sh", Arguments: tt.args})
			writeJobSurface(t, cfg, strings.TrimSpace(tt.job["pid"]), tt.job)
			mock := exitAfterRecoveryProvider()
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume sh job variant", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			payload := requireToolResultInMessages(t, mock.calls[0], "call_sh")
			if got := toolString(t, payload, "status"); got != tt.wantStatus {
				t.Fatalf("status = %q, want %q; payload=%#v", got, tt.wantStatus, payload)
			}
			if got := toolString(t, payload, "mode"); got != tt.wantMode {
				t.Fatalf("mode = %q, want %q; payload=%#v", got, tt.wantMode, payload)
			}
			switch tt.wantField {
			case "stdout_so_far":
				if got := toolString(t, payload, "stdout_so_far"); got != "partial out" {
					t.Fatalf("stdout_so_far = %q, want partial out", got)
				}
			default:
				job := toolMap(t, payload, "job")
				if got, _ := job[tt.wantField].(bool); !got {
					t.Fatalf("job.%s = %#v, want true; payload=%#v", tt.wantField, job[tt.wantField], payload)
				}
			}
		})
	}
}

func TestResumeGathersPendingRelationResults(t *testing.T) {
	for _, toolName := range []string{"fork", "spawn"} {
		t.Run(toolName, func(t *testing.T) {
			cfg := testCfg(t)
			toolID := "call_" + toolName
			args := map[string]any{}
			if toolName == "fork" {
				args["children"] = []any{map[string]any{"intent": "child", "scope": "."}}
			} else {
				args["children"] = []any{map[string]any{"mission": "child"}}
			}
			seedPendingToolBatch(t, cfg, tape.ToolCall{ID: toolID, Name: toolName, Arguments: args})
			relationRoot := filepath.Join(cfg.SessionRetainedDir(""), "relations", toolID)
			if err := os.MkdirAll(relationRoot, 0o755); err != nil {
				t.Fatalf("mkdir relation root: %v", err)
			}
			result := map[string]any{
				"tool":      toolName,
				"mode":      "wait",
				"status":    "completed",
				"requested": 1,
				"spawned":   1,
				"succeeded": 1,
			}
			if err := os.WriteFile(filepath.Join(relationRoot, "result.json"), mustJSON(t, result), 0o644); err != nil {
				t.Fatalf("write relation result: %v", err)
			}
			mock := &mockProvider{responses: []tape.Message{{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{{
					ID:        "call_exit",
					Name:      "exit",
					Arguments: map[string]any{"status": "success"},
				}},
			}}}
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume relation result", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			payload := requireToolResultInMessages(t, mock.calls[0], toolID)
			if got := toolString(t, payload, "status"); got != "completed" {
				t.Fatalf("relation status = %q, want completed; payload=%#v", got, payload)
			}
		})
	}
}

func TestRecoverRelationToolResultPreservesIsErrorSemantics(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		payload     map[string]any
		wantIsError bool
	}{
		{
			name:     "fork forget spawned",
			toolName: "fork",
			payload: map[string]any{
				"tool":      "fork",
				"mode":      "forget",
				"status":    "spawned",
				"requested": 1,
				"spawned":   1,
			},
		},
		{
			name:     "spawn forget spawned",
			toolName: "spawn",
			payload: map[string]any{
				"tool":      "spawn",
				"mode":      "forget",
				"status":    "spawned",
				"requested": 1,
				"spawned":   1,
			},
		},
		{
			name:     "fork wait completed all failed",
			toolName: "fork",
			payload: map[string]any{
				"tool":      "fork",
				"mode":      "wait",
				"status":    "completed",
				"requested": 1,
				"spawned":   1,
				"succeeded": 0,
				"children": []any{map[string]any{
					"index":     0,
					"status":    "completed",
					"exit_code": 2,
				}},
			},
			wantIsError: true,
		},
		{
			name:     "fork wait completed all failed omits succeeded zero",
			toolName: "fork",
			payload: map[string]any{
				"tool":      "fork",
				"mode":      "wait",
				"status":    "completed",
				"requested": 1,
				"spawned":   1,
				"children": []any{map[string]any{
					"index":     0,
					"status":    "completed",
					"exit_code": 2,
				}},
			},
			wantIsError: true,
		},
		{
			name:     "spawn wait completed all failed",
			toolName: "spawn",
			payload: map[string]any{
				"tool":      "spawn",
				"mode":      "wait",
				"status":    "completed",
				"requested": 1,
				"spawned":   1,
				"succeeded": 0,
				"children": []any{map[string]any{
					"index":     0,
					"status":    "completed",
					"exit_code": 2,
				}},
			},
			wantIsError: true,
		},
		{
			name:     "fork wait timeout partial success",
			toolName: "fork",
			payload: map[string]any{
				"tool":      "fork",
				"mode":      "wait",
				"status":    "timeout",
				"requested": 2,
				"spawned":   2,
				"succeeded": 1,
				"killed":    1,
			},
		},
		{
			name:     "spawn wait timeout all failed",
			toolName: "spawn",
			payload: map[string]any{
				"tool":      "spawn",
				"mode":      "wait",
				"status":    "timeout",
				"requested": 2,
				"spawned":   2,
				"succeeded": 0,
				"killed":    2,
			},
			wantIsError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			toolID := "call_" + strings.ReplaceAll(tt.name, " ", "_")
			relationRoot := filepath.Join(cfg.SessionRetainedDir(""), "relations", toolID)
			if err := os.MkdirAll(relationRoot, 0o755); err != nil {
				t.Fatalf("mkdir relation root: %v", err)
			}
			if err := os.WriteFile(filepath.Join(relationRoot, "result.json"), mustJSON(t, tt.payload), 0o644); err != nil {
				t.Fatalf("write relation result: %v", err)
			}
			rt := NewWithProvider(cfg, &mockProvider{})
			silenceRuntime(rt)

			result, ok := rt.recoverRelationToolResult(tape.ToolCall{ID: toolID, Name: tt.toolName})
			if !ok {
				t.Fatalf("recoverRelationToolResult returned ok=false")
			}
			if result.IsError != tt.wantIsError {
				t.Fatalf("IsError = %v, want %v; payload=%s", result.IsError, tt.wantIsError, result.Content)
			}
		})
	}
}

func TestResumeUnknownForMalformedRelationResult(t *testing.T) {
	for _, toolName := range []string{"fork", "spawn"} {
		t.Run(toolName, func(t *testing.T) {
			cfg := testCfg(t)
			toolID := "call_" + toolName
			args := map[string]any{}
			if toolName == "fork" {
				args["children"] = []any{map[string]any{"intent": "child", "scope": "."}}
			} else {
				args["children"] = []any{map[string]any{"mission": "child"}}
			}
			seedPendingToolBatch(t, cfg, tape.ToolCall{ID: toolID, Name: toolName, Arguments: args})
			relationRoot := filepath.Join(cfg.SessionRetainedDir(""), "relations", toolID)
			if err := os.MkdirAll(relationRoot, 0o755); err != nil {
				t.Fatalf("mkdir relation root: %v", err)
			}
			if err := os.WriteFile(filepath.Join(relationRoot, "result.json"), []byte(`{bad json`), 0o644); err != nil {
				t.Fatalf("write malformed relation result: %v", err)
			}
			mock := exitAfterRecoveryProvider()
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume malformed relation", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			payload := requireToolResultInMessages(t, mock.calls[0], toolID)
			if got := toolString(t, payload, "status"); got != "unknown" {
				t.Fatalf("status = %q, want unknown; payload=%#v", got, payload)
			}
		})
	}
}

func TestResumeGathersPendingMarkFromAnchorSurface(t *testing.T) {
	cfg := testCfg(t)
	cfg.AnchorMemoryEnabled = true
	seedPendingToolBatch(t, cfg, tape.ToolCall{
		ID:   "call_mark",
		Name: "mark",
		Arguments: map[string]any{
			"resolution": "already anchored",
		},
	})
	anchorDir := filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "state", "anchors", "5.anchor")
	if err := os.MkdirAll(anchorDir, 0o755); err != nil {
		t.Fatalf("mkdir anchor dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(anchorDir, "meta.json"), []byte(`{"id":5,"resolution":"already anchored"}`), 0o644); err != nil {
		t.Fatalf("write anchor meta: %v", err)
	}
	mock := &mockProvider{responses: []tape.Message{{
		Role: tape.RoleAssistant,
		ToolCalls: []tape.ToolCall{{
			ID:        "call_exit",
			Name:      "exit",
			Arguments: map[string]any{"status": "success"},
		}},
	}}}
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("resume mark result", "Begin."); exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	payload := requireToolResultInMessages(t, mock.calls[0], "call_mark")
	if got := toolString(t, payload, "status"); got != "completed" {
		t.Fatalf("mark status = %q, want completed; payload=%#v", got, payload)
	}
	if got := toolInt(t, payload, "anchor_id"); got != 5 {
		t.Fatalf("anchor_id = %d, want 5", got)
	}
}

func TestResumeUnknownForMalformedOrMismatchedMarkSurface(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
		meta string
	}{
		{
			name: "invalid args",
			args: map[string]any{},
			meta: `{"id":5,"resolution":"already anchored"}`,
		},
		{
			name: "mismatched anchor",
			args: map[string]any{"resolution": "wanted anchor"},
			meta: `{"id":5,"resolution":"other anchor"}`,
		},
		{
			name: "malformed anchor meta",
			args: map[string]any{"resolution": "wanted anchor"},
			meta: `{bad json`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testCfg(t)
			cfg.AnchorMemoryEnabled = true
			seedPendingToolBatch(t, cfg, tape.ToolCall{ID: "call_mark", Name: "mark", Arguments: tt.args})
			anchorDir := filepath.Join(cfg.SessionIncarnationPath("", 0), "context", "state", "anchors", "5.anchor")
			if err := os.MkdirAll(anchorDir, 0o755); err != nil {
				t.Fatalf("mkdir anchor dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(anchorDir, "meta.json"), []byte(tt.meta), 0o644); err != nil {
				t.Fatalf("write anchor meta: %v", err)
			}
			mock := exitAfterRecoveryProvider()
			rt := NewWithProvider(cfg, mock)
			silenceRuntime(rt)

			if exitCode := rt.Run("resume malformed mark surface", "Begin."); exitCode != 0 {
				t.Fatalf("exit code = %d, want 0", exitCode)
			}
			payload := requireToolResultInMessages(t, mock.calls[0], "call_mark")
			if got := toolString(t, payload, "status"); got != "unknown" {
				t.Fatalf("status = %q, want unknown; payload=%#v", got, payload)
			}
		})
	}
}
