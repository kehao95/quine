package runtime

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
)

func TestFinalizerCommitsWorkspaceBeforePublicSurfaceCleanup(t *testing.T) {
	root := t.TempDir()
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
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		if !keepDetached || !commitWorkspace {
			t.Fatalf("success finalization close args = keepDetached:%v commitWorkspace:%v, want true/true", keepDetached, commitWorkspace)
		}
		return os.WriteFile(filepath.Join(root, "committed.txt"), []byte("finalized-before-cleanup\n"), 0o644)
	}
	var cleanupErr error
	installFakePublicSurface(rt, &fakePublicSurface{
		cleanupFn: func() error {
			data, err := os.ReadFile(filepath.Join(root, "committed.txt"))
			if err != nil {
				cleanupErr = fmt.Errorf("workspace file missing before public surface cleanup: %w", err)
				return cleanupErr
			}
			if string(data) != "finalized-before-cleanup\n" {
				cleanupErr = fmt.Errorf("workspace file = %q", string(data))
				return cleanupErr
			}
			return nil
		},
	})

	exitCode := rt.Run("commit before public cleanup", "Begin.")
	if cleanupErr != nil {
		t.Fatal(cleanupErr)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
}

func TestFinalizerDoesNotWriteSuccessOutcomeBeforeWorkspaceCommit(t *testing.T) {
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
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		if !keepDetached || !commitWorkspace {
			t.Fatalf("success finalization close args = keepDetached:%v commitWorkspace:%v, want true/true", keepDetached, commitWorkspace)
		}
		if rt.tape != nil && rt.tape.Outcome != nil {
			t.Fatalf("success outcome was written before workspace commit: %#v", rt.tape.Outcome)
		}
		return nil
	}

	exitCode := rt.Run("success outcome ordering", "Begin.")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	commitIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseWorkspaceCommitted)
	outcomeIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseOutcomeWritten)
	if commitIdx < 0 || outcomeIdx < 0 {
		t.Fatalf("missing finalization phases: %v", rt.finalizationPhases)
	}
	if commitIdx > outcomeIdx {
		t.Fatalf("workspace commit phase came after outcome write: %v", rt.finalizationPhases)
	}
}

func TestFinalizerFailsClosedWhenWorkspaceCommitFails(t *testing.T) {
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
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		if !commitWorkspace {
			t.Fatal("success finalization did not request workspace commit")
		}
		return errors.New("commit failed")
	}

	exitCode := rt.Run("fail closed on commit failure", "Begin.")
	if exitCode == 0 {
		t.Fatal("commit failure must not return success")
	}
	if rt.tape.Outcome == nil {
		t.Fatal("missing failure outcome")
	}
	if rt.tape.Outcome.ExitCode == 0 || rt.tape.Outcome.TerminationMode == tape.TermExit {
		t.Fatalf("commit failure wrote success-like outcome: %#v", rt.tape.Outcome)
	}
	if rt.tape.Outcome.TerminationMode != tape.TermFinalizationFailure {
		t.Fatalf("termination mode = %q, want %q", rt.tape.Outcome.TerminationMode, tape.TermFinalizationFailure)
	}
}

func TestRunExitSuccessUsesOwnedFinalizationProtocol(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	root := t.TempDir()
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "printf 'owned-finalizer\\n' > committed.txt",
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
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = root
	cfg.Workspace = root
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceSession = "runtime-owned-finalizer"
	cfg.WorkspaceOwner = true

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("owned finalizer", "Begin.")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if _, err := os.Stat(filepath.Join(root, "committed.txt")); err != nil {
		t.Fatalf("workspace commit missing: %v", err)
	}
	commitIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseWorkspaceCommitted)
	outcomeIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseOutcomeWritten)
	surfaceIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseSurfaceCleanup)
	if commitIdx < 0 || outcomeIdx < 0 || surfaceIdx < 0 {
		t.Fatalf("missing expected finalization phases: %v", rt.finalizationPhases)
	}
	if !(commitIdx < outcomeIdx && outcomeIdx < surfaceIdx) {
		t.Fatalf("unexpected finalization phase order: %v", rt.finalizationPhases)
	}
}

func TestShutdownStateRecordsDurablePhaseTransitions(t *testing.T) {
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
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		if !keepDetached || !commitWorkspace {
			t.Fatalf("success finalization close args = keepDetached:%v commitWorkspace:%v, want true/true", keepDetached, commitWorkspace)
		}
		return nil
	}

	exitCode := rt.Run("durable finalization state", "Begin.")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}

	phases := readFinalizationPhases(t, rt.finalizationStatePath())
	wantOrder := []string{
		string(finalizationPhaseExitRequested),
		string(finalizationPhaseWorkspaceCommitting),
		string(finalizationPhaseWorkspaceCommitted),
		string(finalizationPhaseOutcomeWritten),
		string(finalizationPhaseSurfaceCleanup),
	}
	next := 0
	for _, phase := range phases {
		if next < len(wantOrder) && phase == wantOrder[next] {
			next++
		}
	}
	if next != len(wantOrder) {
		t.Fatalf("finalization phases %v do not contain required ordered subsequence %v", phases, wantOrder)
	}
}

func TestWorkspaceCommitIntentIsDurableAndIdempotent(t *testing.T) {
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
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = t.TempDir()
	cfg.Workspace = cfg.WorkspaceRoot
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceSession = "runtime-commit-intent"
	cfg.WorkspaceOwner = true

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)
	var sawPending bool
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		if !commitWorkspace {
			t.Fatal("success finalization did not request workspace commit")
		}
		intent := readWorkspaceCommitIntent(t, rt.workspaceCommitIntentPath())
		if intent.Status != "pending" {
			t.Fatalf("intent status before commit = %q, want pending", intent.Status)
		}
		if intent.SessionID != cfg.SessionID || intent.RunID != cfg.RunID || intent.WorkspaceSession != cfg.WorkspaceSession {
			t.Fatalf("intent identity = %#v, want session/run/workspace metadata", intent)
		}
		sawPending = true
		return nil
	}

	exitCode := rt.Run("workspace commit intent", "Begin.")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !sawPending {
		t.Fatal("commit did not observe pending intent")
	}
	intent := readWorkspaceCommitIntent(t, rt.workspaceCommitIntentPath())
	if intent.Status != "committed" {
		t.Fatalf("intent status after commit = %q, want committed", intent.Status)
	}
	if err := rt.writeWorkspaceCommitIntent("committed"); err != nil {
		t.Fatalf("repeat committed intent write failed: %v", err)
	}
	intent = readWorkspaceCommitIntent(t, rt.workspaceCommitIntentPath())
	if intent.Status != "committed" {
		t.Fatalf("intent status after repeated committed write = %q, want committed", intent.Status)
	}
}

func TestWorkspaceCommitIntentRecoveryMaterializesBeforeProviderAndOutcome(t *testing.T) {
	root := t.TempDir()
	cfg := testCfg(t)

	var rt *Runtime
	rt = NewWithProvider(cfg, providerFunc{
		generate: func(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
			if rt.tape != nil && rt.tape.Outcome != nil {
				t.Fatalf("outcome was written before provider observed recovered workspace: %#v", rt.tape.Outcome)
			}
			data, err := os.ReadFile(filepath.Join(root, "recovered.txt"))
			if err != nil {
				t.Fatalf("provider ran before pending workspace commit recovery materialized file: %v", err)
			}
			if string(data) != "recovered-before-provider\n" {
				t.Fatalf("recovered file = %q", string(data))
			}
			return tape.Message{
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
			}, llm.Usage{}, nil
		},
	})
	silenceRuntime(rt)
	if err := rt.writeWorkspaceCommitIntent("pending"); err != nil {
		t.Fatalf("seed pending commit intent: %v", err)
	}
	rt.recoverWorkspaceCommit = func() error {
		if rt.tape != nil && rt.tape.Outcome != nil {
			t.Fatalf("outcome was written before recovery: %#v", rt.tape.Outcome)
		}
		return os.WriteFile(filepath.Join(root, "recovered.txt"), []byte("recovered-before-provider\n"), 0o644)
	}
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		if !keepDetached || !commitWorkspace {
			t.Fatalf("success finalization close args = keepDetached:%v commitWorkspace:%v, want true/true", keepDetached, commitWorkspace)
		}
		return nil
	}

	exitCode := rt.Run("recover pending workspace commit", "Begin.")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	intent := readWorkspaceCommitIntent(t, rt.workspaceCommitIntentPath())
	if intent.Status != "committed" {
		t.Fatalf("intent status after recovery = %q, want committed", intent.Status)
	}
	recoverIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseWorkspaceRecovering)
	commitIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseWorkspaceCommitted)
	outcomeIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseOutcomeWritten)
	if recoverIdx < 0 || commitIdx < 0 || outcomeIdx < 0 {
		t.Fatalf("missing recovery finalization phases: %v", rt.finalizationPhases)
	}
	if !(recoverIdx < commitIdx && commitIdx < outcomeIdx) {
		t.Fatalf("unexpected recovery phase order: %v", rt.finalizationPhases)
	}
}

func TestWorkspaceCommitIntentRecoveryFailureDoesNotWriteSuccessOutcome(t *testing.T) {
	cfg := testCfg(t)
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = t.TempDir()
	cfg.Workspace = cfg.WorkspaceRoot
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceSession = "runtime-commit-intent-recovery-failure"
	cfg.WorkspaceOwner = true

	providerCalled := false
	rt := NewWithProvider(cfg, providerFunc{
		generate: func(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
			providerCalled = true
			return tape.Message{}, llm.Usage{}, errors.New("provider must not be called after recovery failure")
		},
	})
	silenceRuntime(rt)
	if err := rt.writeWorkspaceCommitIntent("pending"); err != nil {
		t.Fatalf("seed pending commit intent: %v", err)
	}
	rt.recoverWorkspaceCommit = func() error {
		return errors.New("synthetic recovery failure")
	}

	exitCode := rt.Run("recover pending workspace commit failure", "Begin.")
	if exitCode == 0 {
		t.Fatal("recovery failure returned success")
	}
	if providerCalled {
		t.Fatal("provider was called after recovery failure")
	}
	if rt.tape != nil && rt.tape.Outcome != nil && rt.tape.Outcome.ExitCode == 0 {
		t.Fatalf("recovery failure wrote success outcome: %#v", rt.tape.Outcome)
	}
	if finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseFailure) < 0 {
		t.Fatalf("recovery failure did not record finalization failure: %v", rt.finalizationPhases)
	}
	intent := readWorkspaceCommitIntent(t, rt.workspaceCommitIntentPath())
	if intent.Status != "pending" {
		t.Fatalf("failed recovery changed intent status = %q, want pending", intent.Status)
	}
}

func TestFusePublicSurfaceUnmountRequiresQuiescedOwnerSet(t *testing.T) {
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
	var cleanupErr error
	installFakePublicSurface(rt, &fakePublicSurface{
		cleanupFn: func() error {
			quiescedIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseRuntimeQuiesced)
			cleanupIdx := finalizationPhaseIndex(rt.finalizationPhases, finalizationPhaseSurfaceCleanup)
			if quiescedIdx < 0 {
				cleanupErr = fmt.Errorf("public surface cleanup ran before runtime owner quiescence was recorded: %v", rt.finalizationPhases)
				return cleanupErr
			}
			if cleanupIdx < 0 || quiescedIdx > cleanupIdx {
				cleanupErr = fmt.Errorf("runtime owner quiescence did not precede surface cleanup: %v", rt.finalizationPhases)
				return cleanupErr
			}
			return nil
		},
	})

	exitCode := rt.Run("quiesce before surface cleanup", "Begin.")
	if cleanupErr != nil {
		t.Fatal(cleanupErr)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
}

func TestPublicSurfaceCleanupCannotBlockWorkspaceMaterialization(t *testing.T) {
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
	committed := make(chan struct{})
	cleanupEntered := make(chan struct{})
	releaseCleanup := make(chan struct{})
	rt.closeRuntimeSubstrates = func(keepDetached bool, commitWorkspace bool) error {
		close(committed)
		return nil
	}
	installFakePublicSurface(rt, &fakePublicSurface{
		cleanupFn: func() error {
			select {
			case <-committed:
			default:
				return errors.New("public surface cleanup started before workspace commit")
			}
			close(cleanupEntered)
			<-releaseCleanup
			return nil
		},
	})

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("blocking public cleanup", "Begin.")
	}()
	select {
	case <-committed:
	case <-time.After(2 * time.Second):
		t.Fatal("workspace commit did not complete before cleanup")
	}
	select {
	case <-cleanupEntered:
	case <-time.After(2 * time.Second):
		t.Fatal("public surface cleanup did not start")
	}
	if rt.tape == nil || rt.tape.Outcome == nil || rt.tape.Outcome.ExitCode != 0 {
		t.Fatalf("success outcome not written before blocked cleanup: %#v", rt.tape)
	}
	close(releaseCleanup)
	select {
	case exitCode := <-done:
		if exitCode != 0 {
			t.Fatalf("exit code = %d, want 0", exitCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runtime did not return after public cleanup released")
	}
}

func TestProgramBenchLikeExitCannotLeaveSuccessWithUncommittedWorkspace(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	root := t.TempDir()
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": strings.Join([]string{
								"mkdir -p source",
								"printf '#!/bin/sh\\necho built\\n' > compile.sh",
								"chmod +x compile.sh",
								"printf 'package main\\nfunc main(){}\\n' > source/main.go",
							}, " && "),
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
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = root
	cfg.Workspace = root
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceSession = "runtime-programbench-finalizer"
	cfg.WorkspaceOwner = true

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)
	var cleanupErr error
	installFakePublicSurface(rt, &fakePublicSurface{
		cleanupFn: func() error {
			if rt.tape == nil || rt.tape.Outcome == nil || rt.tape.Outcome.ExitCode != 0 {
				return nil
			}
			for _, rel := range []string{"compile.sh", filepath.Join("source", "main.go")} {
				if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
					cleanupErr = fmt.Errorf("success outcome visible before %s was materialized: %w", rel, err)
					return cleanupErr
				}
			}
			return nil
		},
	})

	exitCode := rt.Run("programbench-like finalizer", "Begin.")
	if cleanupErr != nil {
		t.Fatal(cleanupErr)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	for _, rel := range []string{"compile.sh", filepath.Join("source", "main.go")} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("workspace file %s missing after success: %v", rel, err)
		}
	}
	if rt.tape.Outcome == nil || rt.tape.Outcome.ExitCode != 0 {
		t.Fatalf("missing success outcome: %#v", rt.tape.Outcome)
	}
}

func TestExecutionBudgetFinalExitCommitsWorkspaceOnSuccess(t *testing.T) {
	requireOverlayWorkspaceSupport(t)

	root := t.TempDir()
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "printf 'budget-final-exit\\n' > committed.txt",
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
	cfg.MaxTurns = 1
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = root
	cfg.Workspace = root
	cfg.WorkspaceBackend = "overlay"
	cfg.WorkspaceRevisionMode = config.WorkspaceRevisionRestore
	cfg.WorkspaceSession = "runtime-budget-final-exit-commit"
	cfg.WorkspaceOwner = true

	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("commit after final exit", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(filepath.Join(root, "committed.txt"))
	if err != nil {
		t.Fatalf("read committed workspace file: %v", err)
	}
	if string(data) != "budget-final-exit\n" {
		t.Fatalf("committed workspace file = %q, want %q", string(data), "budget-final-exit\n")
	}
}
