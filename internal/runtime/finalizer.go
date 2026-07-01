package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

type finalizationPhase string

type finalizationPhaseRecord struct {
	Phase     string `json:"phase"`
	Timestamp int64  `json:"timestamp"`
	SessionID string `json:"session_id,omitempty"`
	RunID     string `json:"run_id,omitempty"`
}

type workspaceCommitIntent struct {
	Status           string `json:"status"`
	Timestamp        int64  `json:"timestamp"`
	SessionID        string `json:"session_id,omitempty"`
	RunID            string `json:"run_id,omitempty"`
	WorkspaceSession string `json:"workspace_session,omitempty"`
	WorkspaceRoot    string `json:"workspace_root,omitempty"`
	Workspace        string `json:"workspace,omitempty"`
	WorkspaceBackend string `json:"workspace_backend,omitempty"`
	WorldRevision    string `json:"world_revision,omitempty"`
}

const (
	finalizationPhaseExitRequested       finalizationPhase = "exit_requested"
	finalizationPhaseWorkspaceRecovering finalizationPhase = "workspace_recovering"
	finalizationPhaseWorkspaceCommitting finalizationPhase = "workspace_committing"
	finalizationPhaseWorkspaceCommitted  finalizationPhase = "workspace_committed"
	finalizationPhaseWorkspaceRollback   finalizationPhase = "workspace_rolling_back"
	finalizationPhaseWorkspaceRolledBack finalizationPhase = "workspace_rolled_back"
	finalizationPhaseOutcomeWritten      finalizationPhase = "outcome_written"
	finalizationPhaseRuntimeQuiescing    finalizationPhase = "runtime_owners_quiescing"
	finalizationPhaseRuntimeQuiesced     finalizationPhase = "runtime_owners_quiesced"
	finalizationPhaseSurfaceCleanup      finalizationPhase = "surface_cleanup"
	finalizationPhaseFailure             finalizationPhase = "finalization_failed"
)

func (r *Runtime) recordFinalizationPhase(phase finalizationPhase) {
	if r == nil {
		return
	}
	r.finalizationMu.Lock()
	r.finalizationPhases = append(r.finalizationPhases, phase)
	r.finalizationMu.Unlock()
	r.log("finalization phase: %s", phase)
	if err := r.appendFinalizationPhaseRecord(phase); err != nil {
		r.log("finalization phase record error: %v", err)
	}
}

func (r *Runtime) finalizationStatePath() string {
	if r == nil || r.cfg == nil {
		return ""
	}
	return filepath.Join(r.cfg.SessionRetainedDir(""), "status", "finalization.jsonl")
}

func (r *Runtime) workspaceCommitIntentPath() string {
	if r == nil || r.cfg == nil {
		return ""
	}
	return filepath.Join(r.cfg.SessionRetainedDir(""), "status", "workspace-commit-intent.json")
}

func (r *Runtime) appendFinalizationPhaseRecord(phase finalizationPhase) error {
	path := r.finalizationStatePath()
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	record := finalizationPhaseRecord{
		Phase:     string(phase),
		Timestamp: time.Now().UnixMilli(),
	}
	if r.cfg != nil {
		record.SessionID = r.cfg.SessionID
		record.RunID = r.cfg.RunID
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) writeWorkspaceCommitIntent(status string) error {
	path := r.workspaceCommitIntentPath()
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	intent := workspaceCommitIntent{
		Status:    status,
		Timestamp: time.Now().UnixMilli(),
	}
	if r.cfg != nil {
		intent.SessionID = r.cfg.SessionID
		intent.RunID = r.cfg.RunID
		intent.WorkspaceSession = r.cfg.WorkspaceSession
		intent.WorkspaceRoot = r.cfg.WorkspaceRoot
		intent.Workspace = r.cfg.Workspace
		intent.WorkspaceBackend = r.cfg.EffectiveWorkspaceBackend()
	}
	if r.sh != nil {
		intent.WorldRevision = r.sh.CurrentWorldRevision()
	}
	data, err := json.MarshalIndent(intent, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (r *Runtime) readWorkspaceCommitIntent() (workspaceCommitIntent, bool, error) {
	path := r.workspaceCommitIntentPath()
	if path == "" {
		return workspaceCommitIntent{}, false, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return workspaceCommitIntent{}, false, nil
		}
		return workspaceCommitIntent{}, false, err
	}
	var intent workspaceCommitIntent
	if err := json.Unmarshal(data, &intent); err != nil {
		return workspaceCommitIntent{}, false, err
	}
	return intent, true, nil
}

func (r *Runtime) recoverPendingWorkspaceCommit() error {
	intent, ok, err := r.readWorkspaceCommitIntent()
	if err != nil {
		return fmt.Errorf("read workspace commit intent: %w", err)
	}
	if !ok || strings.TrimSpace(intent.Status) != "pending" {
		return nil
	}
	r.recordFinalizationPhase(finalizationPhaseWorkspaceRecovering)
	recoverCommit := r.recoverWorkspaceCommit
	if recoverCommit == nil {
		recoverCommit = func() error {
			return tools.RecoverWorkspaceCommit(r.cfg)
		}
	}
	if err := recoverCommit(); err != nil {
		r.recordFinalizationPhase(finalizationPhaseFailure)
		return err
	}
	if r.sh != nil {
		r.sh.ResetWorkspaceAfterExternalCommit()
	}
	if err := r.writeWorkspaceCommitIntent("committed"); err != nil {
		return fmt.Errorf("mark recovered workspace commit intent committed: %w", err)
	}
	r.recordFinalizationPhase(finalizationPhaseWorkspaceCommitted)
	return nil
}

func (r *Runtime) markShFinalized() {
	if r == nil {
		return
	}
	r.finalizationMu.Lock()
	r.shFinalized = true
	r.finalizationMu.Unlock()
}

func (r *Runtime) shAlreadyFinalized() bool {
	if r == nil {
		return true
	}
	r.finalizationMu.Lock()
	defer r.finalizationMu.Unlock()
	return r.shFinalized
}

func (r *Runtime) finalizationStarted() bool {
	if r == nil {
		return false
	}
	r.finalizationMu.Lock()
	defer r.finalizationMu.Unlock()
	return len(r.finalizationPhases) > 0
}

func (r *Runtime) closeSubstrates(keepDetached bool, commitWorkspace bool) error {
	if r == nil || r.sh == nil {
		return nil
	}
	if commitWorkspace {
		if err := r.writeWorkspaceCommitIntent("pending"); err != nil {
			return fmt.Errorf("write workspace commit intent: %w", err)
		}
		r.recordFinalizationPhase(finalizationPhaseWorkspaceCommitting)
		if os.Getenv("QUINE_TEST_EXIT_AFTER_WORKSPACE_COMMIT_INTENT") == "1" {
			os.Exit(86)
		}
	} else {
		r.recordFinalizationPhase(finalizationPhaseWorkspaceRollback)
	}

	var err error
	if r.closeRuntimeSubstrates != nil {
		err = r.closeRuntimeSubstrates(keepDetached, commitWorkspace)
	} else {
		err = r.sh.CloseWithOptions(keepDetached, commitWorkspace)
	}
	r.markShFinalized()
	if err != nil {
		return err
	}

	if commitWorkspace {
		if err := r.writeWorkspaceCommitIntent("committed"); err != nil {
			return fmt.Errorf("mark workspace commit intent committed: %w", err)
		}
		r.recordFinalizationPhase(finalizationPhaseWorkspaceCommitted)
	} else {
		r.recordFinalizationPhase(finalizationPhaseWorkspaceRolledBack)
	}
	return nil
}

func (r *Runtime) quiesceRuntimeOwners(pidPublished bool) {
	if r == nil {
		return
	}
	record := r.finalizationStarted()
	if record {
		r.recordFinalizationPhase(finalizationPhaseRuntimeQuiescing)
	}
	r.stopPeerDiscoveryHeartbeat()
	if pidPublished && r.agentRegistry != nil {
		if err := r.agentRegistry.UnpublishSelfPID(); err != nil {
			r.log("agent pid unpublish error: %v", err)
		}
	}
	if r.agentRegistry != nil {
		if err := r.agentRegistry.Deregister(); err != nil {
			r.log("agent deregistration error: %v", err)
		}
	}
	if record {
		r.recordFinalizationPhase(finalizationPhaseRuntimeQuiesced)
	}
}

func (r *Runtime) finalizeExitRequest(exitReq tools.ExitRequest) int {
	r.recordFinalizationPhase(finalizationPhaseExitRequested)

	exitCode := exitReq.ExitCode()
	if exitReq.Stderr != "" {
		fmt.Fprint(r.stderr, exitReq.Stderr)
	}

	return r.finalizeOutcome(exitCode, exitReq.Stderr, tape.TermExit)
}

func (r *Runtime) finalizeOutcome(exitCode int, stderr string, mode tape.TerminationMode) int {
	r.flushPendingToolResult()
	commitWorkspace := exitCode == 0
	if err := r.closeSubstrates(commitWorkspace, commitWorkspace); err != nil {
		r.recordFinalizationPhase(finalizationPhaseFailure)
		finalizationStderr := fmt.Sprintf("finalization failed: %v", err)
		r.log("%s", finalizationStderr)
		r.logError("%s", finalizationStderr)
		if r.stderr != nil {
			fmt.Fprint(r.stderr, finalizationStderr)
		}
		r.writeOutcome(1, finalizationStderr, tape.TermFinalizationFailure)
		return 1
	}

	r.writeOutcome(exitCode, stderr, mode)
	return exitCode
}

func (r *Runtime) writeOutcome(exitCode int, stderr string, mode tape.TerminationMode) {
	r.flushPendingToolResult()
	duration := time.Since(r.startTime)
	r.tape.SetOutcome(tape.SessionOutcome{
		ExitCode:        exitCode,
		Stderr:          stderr,
		DurationMs:      duration.Milliseconds(),
		TerminationMode: mode,
	})
	r.writeTapeEntry(r.tape.OutcomeEntry())
	r.recordFinalizationPhase(finalizationPhaseOutcomeWritten)
}
