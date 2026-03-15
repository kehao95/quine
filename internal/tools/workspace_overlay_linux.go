//go:build linux

package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func (s *subjectiveFS) init(dataDir, sessionID string) error {
	if !s.enabled || s.initialized {
		return nil
	}

	if s.workspaceRoot == "" || s.workspace == "" {
		return fmt.Errorf("workspace physics require workspace root and workspace")
	}

	rootReal, err := realishPath(s.workspaceRoot)
	if err != nil {
		return fmt.Errorf("canonicalize workspace root: %w", err)
	}
	workspaceReal, err := realishPath(s.workspace)
	if err != nil {
		return fmt.Errorf("canonicalize workspace: %w", err)
	}
	dataDirReal, err := realishPath(dataDir)
	if err != nil {
		return fmt.Errorf("canonicalize data dir: %w", err)
	}

	s.workspaceRoot = rootReal
	s.workspace = workspaceReal
	s.dataDir = dataDirReal
	if s.workspaceSession == "" {
		s.workspaceSession = sessionID
	}

	if s.directBackend() {
		if err := os.MkdirAll(s.workspaceStateDir(), 0o755); err != nil {
			return fmt.Errorf("create direct workspace state dir: %w", err)
		}
		if err := s.ensureDirectBaseline(); err != nil {
			return err
		}
		if err := s.ensureWorldRevisionLedger(); err != nil {
			return err
		}
		s.initialized = true
		return nil
	}

	if _, err := exec.LookPath("mount"); err != nil {
		return fmt.Errorf("workspace physics require mount(8): %w", err)
	}
	if _, err := exec.LookPath("umount"); err != nil {
		return fmt.Errorf("workspace physics require umount(8): %w", err)
	}
	if _, err := exec.LookPath("tar"); err != nil {
		return fmt.Errorf("workspace physics require tar(1): %w", err)
	}
	if _, err := exec.LookPath("python3"); err != nil {
		return fmt.Errorf("workspace checkpoint restore requires python3: %w", err)
	}

	if err := os.MkdirAll(s.upperDir(), 0o755); err != nil {
		return fmt.Errorf("create overlay upperdir: %w", err)
	}
	if err := os.MkdirAll(s.mountBase(), 0o755); err != nil {
		return fmt.Errorf("create overlay mount base: %w", err)
	}
	if err := os.MkdirAll(s.viewBase(), 0o755); err != nil {
		return fmt.Errorf("create overlay view base: %w", err)
	}

	if err := s.preflight(); err != nil {
		return err
	}
	if err := s.ensureWorldRevisionLedger(); err != nil {
		return err
	}

	s.initialized = true
	return nil
}

func (s *subjectiveFS) commandEnv() []string {
	if !s.enabled {
		return nil
	}
	env := []string{
		"QUINE_WORKSPACE_ENABLED=1",
		"QUINE_WORKSPACE_ROOT=" + s.workspaceRoot,
		"QUINE_WORKSPACE=" + s.workspace,
		"QUINE_WORKSPACE_BACKEND=" + s.workspaceBackend,
		"QUINE_WORKSPACE_SESSION=" + s.workspaceSession,
	}
	if !s.directBackend() {
		env = append(env,
			"QUINE_WORKSPACE_UPPER="+s.upperDir(),
			"QUINE_WORKSPACE_MOUNT_BASE="+s.mountBase(),
		)
	}
	return env
}

func (s *subjectiveFS) childEnvOverrides() []string {
	if !s.enabled {
		return nil
	}
	return []string{
		"QUINE_WORKSPACE_ROOT=" + s.workspaceRoot,
		"QUINE_WORKSPACE=" + s.workspace,
		"QUINE_WORKSPACE_BACKEND=" + s.workspaceBackend,
		"QUINE_WORKSPACE_SESSION=" + s.workspaceSession,
		"QUINE_WORKSPACE_CURRENT_REVISION=" + s.currentWorldRevision(),
		"QUINE_WORKSPACE_OWNER=" + fmt.Sprintf("%t", s.workspaceOwner),
	}
}

func (s *subjectiveFS) snapshot() (fsSnapshot, error) {
	if !s.enabled {
		return fsSnapshot{}, nil
	}
	if s.directBackend() {
		entries, err := snapshotTree(s.workspaceRoot)
		if err != nil {
			return fsSnapshot{}, err
		}
		return fsSnapshot{workspace: entries}, nil
	}
	viewRoot, err := os.MkdirTemp(s.viewBase(), "snapshot-")
	if err != nil {
		return fsSnapshot{}, fmt.Errorf("create overlay snapshot view: %w", err)
	}
	defer os.RemoveAll(viewRoot)

	if err := s.exportMergedTree(s.upperDir(), viewRoot); err != nil {
		return fsSnapshot{}, err
	}

	entries, err := snapshotTree(viewRoot)
	if err != nil {
		return fsSnapshot{}, err
	}
	return fsSnapshot{workspace: entries}, nil
}

func (s *subjectiveFS) formatMutations(before, after fsSnapshot) string {
	if !s.enabled {
		return ""
	}
	return formatMutations(diffTree(before.workspace, after.workspace))
}

func (s *subjectiveFS) commit() error {
	if !s.enabled || !s.workspaceOwner {
		return nil
	}
	if s.directBackend() {
		return os.RemoveAll(s.workspaceStateDir())
	}
	viewRoot, err := os.MkdirTemp(s.viewBase(), "commit-")
	if err != nil {
		return fmt.Errorf("create overlay commit view: %w", err)
	}
	defer os.RemoveAll(viewRoot)

	if err := s.exportMergedTree(s.upperDir(), viewRoot); err != nil {
		return err
	}
	if err := syncTree(viewRoot, s.workspaceRoot); err != nil {
		return err
	}
	return os.RemoveAll(s.workspaceStateDir())
}

func (s *subjectiveFS) rollback() error {
	if !s.enabled || !s.workspaceOwner {
		return nil
	}
	if s.directBackend() {
		if err := s.restoreDirectBaseline(); err != nil {
			return err
		}
		return os.RemoveAll(s.workspaceStateDir())
	}
	return os.RemoveAll(s.workspaceStateDir())
}

func (s *subjectiveFS) preflight() error {
	if s.directBackend() {
		return nil
	}
	tempRoot, err := os.MkdirTemp(s.workspaceStateDir(), "preflight-")
	if err != nil {
		return fmt.Errorf("create overlay preflight dir: %w", err)
	}
	defer os.RemoveAll(tempRoot)

	upper := filepath.Join(tempRoot, "upper")
	view := filepath.Join(tempRoot, "view")
	if err := os.MkdirAll(upper, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(view, 0o755); err != nil {
		return err
	}
	if err := s.exportMergedTree(upper, view); err != nil {
		return fmt.Errorf("workspace physics unsupported in this Linux environment: %w", err)
	}
	return nil
}

func (s *subjectiveFS) exportMergedTree(upperDir, targetDir string) error {
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create overlay export dir: %w", err)
	}

	const script = `
set -eu
_tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/quine-ov.XXXXXX")"
_upper_parent="$(dirname "${QUINE_WORKSPACE_UPPER}")"
_work_dir="$(mktemp -d "${_upper_parent}/.quine-ov-work.XXXXXX")"
_merged_dir="${_tmp_dir}/merged"
cleanup() {
  umount "${_merged_dir}" 2>/dev/null || true
  rm -rf "${_work_dir}" 2>/dev/null || true
  rm -rf "${_tmp_dir}" 2>/dev/null || true
}
trap cleanup EXIT INT TERM
mkdir -p "${_work_dir}" "${_merged_dir}"
mount -t overlay overlay -o "lowerdir=${QUINE_WORKSPACE_ROOT},upperdir=${QUINE_WORKSPACE_UPPER},workdir=${_work_dir}" "${_merged_dir}"
tar -C "${_merged_dir}" -cf - . | tar -C "${QUINE_EXPORT_TARGET}" -xf -
`

	cmd := exec.Command("/bin/sh", "-ceu", script)
	cmd.SysProcAttr = jobSysProcAttr(false, true)
	cmd.Env = MergeEnv(os.Environ(), []string{
		"QUINE_WORKSPACE_ROOT=" + s.workspaceRoot,
		"QUINE_WORKSPACE_UPPER=" + upperDir,
		"QUINE_EXPORT_TARGET=" + targetDir,
	})

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errText := stderr.String()
		if errText == "" {
			errText = err.Error()
		}
		return fmt.Errorf("mount/export merged workspace: %s", errText)
	}
	return nil
}

func (s *subjectiveFS) syncMergedTreeIntoUpper(sourceDir, upperDir string) error {
	const script = `
set -eu
_tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/quine-ov.XXXXXX")"
_work_dir="$(mktemp -d "${_tmp_dir}/work.XXXXXX")"
_merged_dir="${_tmp_dir}/merged"
cleanup() {
  umount "${_merged_dir}" 2>/dev/null || true
  rm -rf "${_work_dir}" 2>/dev/null || true
  rm -rf "${_tmp_dir}" 2>/dev/null || true
}
trap cleanup EXIT INT TERM
mkdir -p "${_work_dir}" "${_merged_dir}" "${QUINE_WORKSPACE_UPPER}"
mount -t overlay overlay -o "lowerdir=${QUINE_WORKSPACE_ROOT},upperdir=${QUINE_WORKSPACE_UPPER},workdir=${_work_dir}" "${_merged_dir}"
python3 - "${QUINE_SYNC_SOURCE}" "${_merged_dir}" <<'PY'
import os
import shutil
import sys

src = sys.argv[1]
dst = sys.argv[2]

for root, dirs, files in os.walk(dst, topdown=False):
    rel = os.path.relpath(root, dst)
    src_root = src if rel == "." else os.path.join(src, rel)
    for name in files:
        dst_path = os.path.join(root, name)
        src_path = os.path.join(src_root, name)
        if not os.path.lexists(src_path):
            os.remove(dst_path)
    for name in dirs:
        dst_path = os.path.join(root, name)
        src_path = os.path.join(src_root, name)
        if not os.path.lexists(src_path):
            shutil.rmtree(dst_path)

for root, dirs, files in os.walk(src):
    rel = os.path.relpath(root, src)
    dst_root = dst if rel == "." else os.path.join(dst, rel)
    os.makedirs(dst_root, exist_ok=True)
    for name in dirs:
        os.makedirs(os.path.join(dst_root, name), exist_ok=True)
    for name in files:
        src_path = os.path.join(root, name)
        dst_path = os.path.join(dst_root, name)
        if os.path.islink(src_path):
            if os.path.lexists(dst_path):
                if os.path.isdir(dst_path) and not os.path.islink(dst_path):
                    shutil.rmtree(dst_path)
                else:
                    os.remove(dst_path)
            os.symlink(os.readlink(src_path), dst_path)
            continue
        if os.path.lexists(dst_path):
            if os.path.isdir(dst_path) and not os.path.islink(dst_path):
                shutil.rmtree(dst_path)
            else:
                os.remove(dst_path)
        shutil.copy2(src_path, dst_path)
PY
`

	cmd := exec.Command("/bin/sh", "-ceu", script)
	cmd.SysProcAttr = jobSysProcAttr(false, true)
	cmd.Env = MergeEnv(os.Environ(), []string{
		"QUINE_WORKSPACE_ROOT=" + s.workspaceRoot,
		"QUINE_WORKSPACE_UPPER=" + upperDir,
		"QUINE_SYNC_SOURCE=" + sourceDir,
	})

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errText := stderr.String()
		if errText == "" {
			errText = err.Error()
		}
		return fmt.Errorf("sync merged checkpoint into overlay upperdir: %s", errText)
	}
	return nil
}

func (s *subjectiveFS) workspaceStateDir() string {
	return filepath.Join(s.dataDir, "workspaces", s.workspaceSession)
}

func (s *subjectiveFS) worldRevisionsDir() string {
	return filepath.Join(s.workspaceStateDir(), "revisions")
}

func (s *subjectiveFS) worldRevisionDir(revision string) string {
	return filepath.Join(s.worldRevisionsDir(), revision)
}

func (s *subjectiveFS) worldRevisionLedgerPath() string {
	return filepath.Join(s.workspaceStateDir(), "world-revisions.json")
}

func (s *subjectiveFS) captureCurrentTree(targetDir string) error {
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("reset revision dir: %w", err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create revision dir: %w", err)
	}
	if s.directBackend() {
		if err := syncTree(s.workspaceRoot, targetDir); err != nil {
			return fmt.Errorf("capture direct workspace tree: %w", err)
		}
		return nil
	}
	if err := s.exportMergedTree(s.upperDir(), targetDir); err != nil {
		return fmt.Errorf("capture overlay workspace tree: %w", err)
	}
	return nil
}

func (s *subjectiveFS) loadWorldRevisionLedger() (worldRevisionLedger, error) {
	var ledger worldRevisionLedger
	data, err := os.ReadFile(s.worldRevisionLedgerPath())
	if err != nil {
		return ledger, err
	}
	if err := json.Unmarshal(data, &ledger); err != nil {
		return ledger, fmt.Errorf("parse world revision ledger: %w", err)
	}
	if ledger.Revisions == nil {
		ledger.Revisions = make(map[string]worldRevision)
	}
	return ledger, nil
}

func (s *subjectiveFS) saveWorldRevisionLedger(ledger worldRevisionLedger) error {
	if ledger.Revisions == nil {
		ledger.Revisions = make(map[string]worldRevision)
	}
	data, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal world revision ledger: %w", err)
	}
	if err := os.WriteFile(s.worldRevisionLedgerPath(), data, 0o644); err != nil {
		return fmt.Errorf("write world revision ledger: %w", err)
	}
	s.currentRevision = ledger.Current
	return nil
}

func (s *subjectiveFS) ensureWorldRevisionLedger() error {
	if !s.canTrackWorldRevisions() {
		return nil
	}
	if err := os.MkdirAll(s.worldRevisionsDir(), 0o755); err != nil {
		return fmt.Errorf("create world revisions dir: %w", err)
	}
	if _, err := os.Lstat(s.worldRevisionLedgerPath()); err == nil {
		ledger, loadErr := s.loadWorldRevisionLedger()
		if loadErr != nil {
			return loadErr
		}
		s.currentRevision = ledger.Current
		return nil
	} else if !errorsIsNotExist(err) {
		return fmt.Errorf("stat world revision ledger: %w", err)
	}

	base := worldRevision{ID: "wr0", Kind: "baseline"}
	if err := s.captureCurrentTree(s.worldRevisionDir(base.ID)); err != nil {
		return fmt.Errorf("capture baseline world revision: %w", err)
	}
	ledger := worldRevisionLedger{
		Current: "wr0",
		Next:    1,
		Revisions: map[string]worldRevision{
			base.ID: base,
		},
	}
	return s.saveWorldRevisionLedger(ledger)
}

func (s *subjectiveFS) captureWorldRevision(kind string, turnID int) (worldRevision, error) {
	if !s.canTrackWorldRevisions() {
		return worldRevision{}, nil
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return worldRevision{}, err
	}
	revision := worldRevision{
		ID:     "wr" + strconv.Itoa(ledger.Next),
		Parent: ledger.Current,
		Kind:   kind,
		Turn:   turnID,
	}
	if err := s.captureCurrentTree(s.worldRevisionDir(revision.ID)); err != nil {
		return worldRevision{}, fmt.Errorf("capture world revision %s: %w", revision.ID, err)
	}
	ledger.Next++
	ledger.Current = revision.ID
	ledger.Revisions[revision.ID] = revision
	if err := s.saveWorldRevisionLedger(ledger); err != nil {
		return worldRevision{}, err
	}
	return revision, nil
}

func (s *subjectiveFS) loadCurrentWorldRevision() (worldRevision, error) {
	if !s.canTrackWorldRevisions() {
		return worldRevision{}, nil
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return worldRevision{}, err
	}
	revision, ok := ledger.Revisions[ledger.Current]
	if !ok {
		return worldRevision{}, fmt.Errorf("current world revision %s does not exist", ledger.Current)
	}
	return revision, nil
}

func (s *subjectiveFS) restoreWorld(revision string) (string, string, error) {
	if !s.canRestoreWorld() {
		return "", "", nil
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return "", "", err
	}
	target, ok := ledger.Revisions[revision]
	if !ok {
		return "", "", fmt.Errorf("world revision %s does not exist", revision)
	}
	sourceDir := s.worldRevisionDir(target.ID)
	if _, err := os.Lstat(sourceDir); err != nil {
		if errorsIsNotExist(err) {
			return "", "", fmt.Errorf("world revision %s is missing its snapshot", revision)
		}
		return "", "", fmt.Errorf("stat world revision %s: %w", revision, err)
	}
	previous := ledger.Current
	if s.directBackend() {
		if err := syncTree(sourceDir, s.workspaceRoot); err != nil {
			return "", "", fmt.Errorf("restore direct world revision %s: %w", revision, err)
		}
	} else {
		if err := os.RemoveAll(s.upperDir()); err != nil {
			return "", "", fmt.Errorf("reset overlay upperdir: %w", err)
		}
		if err := os.MkdirAll(s.upperDir(), 0o755); err != nil {
			return "", "", fmt.Errorf("recreate overlay upperdir: %w", err)
		}
		if err := s.syncMergedTreeIntoUpper(sourceDir, s.upperDir()); err != nil {
			return "", "", fmt.Errorf("restore overlay world revision %s: %w", revision, err)
		}
	}
	ledger.Current = revision
	if err := s.saveWorldRevisionLedger(ledger); err != nil {
		return "", "", err
	}
	return previous, revision, nil
}

func (s *subjectiveFS) baselineDir() string {
	return filepath.Join(s.workspaceStateDir(), "baseline")
}

func (s *subjectiveFS) upperDir() string {
	return filepath.Join(s.workspaceStateDir(), "upper")
}

func (s *subjectiveFS) mountBase() string {
	return filepath.Join(s.workspaceStateDir(), "mounts")
}

func (s *subjectiveFS) viewBase() string {
	return filepath.Join(s.workspaceStateDir(), "views")
}

func (s *subjectiveFS) ensureDirectBaseline() error {
	base := s.baselineDir()
	if _, err := os.Lstat(base); err == nil {
		return nil
	} else if !errorsIsNotExist(err) {
		return fmt.Errorf("stat direct workspace baseline: %w", err)
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return fmt.Errorf("create direct workspace baseline dir: %w", err)
	}
	if err := syncTree(s.workspaceRoot, base); err != nil {
		return fmt.Errorf("seed direct workspace baseline: %w", err)
	}
	return nil
}

func (s *subjectiveFS) restoreDirectBaseline() error {
	base := s.baselineDir()
	if _, err := os.Lstat(base); err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat direct workspace baseline: %w", err)
	}
	if err := syncTree(base, s.workspaceRoot); err != nil {
		return fmt.Errorf("restore direct workspace baseline: %w", err)
	}
	return nil
}
