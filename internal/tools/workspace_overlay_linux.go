//go:build linux

package tools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
)

func (s *subjectiveFS) init(dataDir, sessionID string) error {
	if !s.enabled || s.initialized {
		return nil
	}
	s.sessionID = sessionID

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

	if s.usesDirectBackend() {
		if err := s.initDirectState(); err != nil {
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
	if s.overlayDriver == "" {
		s.overlayDriver = "kernel"
	}

	if err := s.bootstrapWorkspaceState(); err != nil {
		return err
	}
	if err := s.ensureOverlayLayout(); err != nil {
		return err
	}
	if err := s.ensureWorldRevisionLedger(); err != nil {
		return err
	}
	if err := s.preflight(); err != nil {
		return err
	}

	s.initialized = true
	return nil
}

func (s *subjectiveFS) commandEnv() []string {
	if !s.enabled || s.usesDirectBackend() {
		return nil
	}
	lowerdir, err := s.currentLowerDirString()
	if err != nil {
		return []string{
			"QUINE_WORKSPACE_ENABLED=1",
			"QUINE_WORKSPACE_INIT_ERROR=" + err.Error(),
		}
	}
	stateDir := s.workspaceStateDir()
	mountStateDir, err := s.overlayMountStateDir()
	if err != nil {
		return []string{
			"QUINE_WORKSPACE_ENABLED=1",
			"QUINE_WORKSPACE_INIT_ERROR=" + err.Error(),
		}
	}
	lowerdir = rewriteOverlayLowerDirString(lowerdir, stateDir, mountStateDir)
	return []string{
		"QUINE_WORKSPACE_ENABLED=1",
		config.EnvWorkspaceRoot + "=" + s.workspaceRoot,
		config.EnvWorkspace + "=" + s.workspace,
		config.EnvWorkspaceBackend + "=overlay",
		config.EnvWorkspaceOverlayDriver + "=" + s.overlayDriver,
		config.EnvWorkspaceSession + "=" + s.workspaceSession,
		"QUINE_WORKSPACE_LOWERDIR=" + lowerdir,
		"QUINE_WORKSPACE_UPPER=" + rewriteOverlayStatePath(s.liveUpperDir(), stateDir, mountStateDir),
		"QUINE_WORKSPACE_WORKDIR=" + rewriteOverlayStatePath(s.liveWorkDir(), stateDir, mountStateDir),
		"QUINE_WORKSPACE_MOUNT_BASE=" + rewriteOverlayStatePath(s.mountBase(), stateDir, mountStateDir),
		"QUINE_WORKSPACE_OVERLAY_EXTRA_OPTS=" + s.overlayMountOpts,
	}
}

func (s *subjectiveFS) childEnvOverrides() []string {
	if !s.enabled || s.usesDirectBackend() {
		return nil
	}
	return []string{
		config.EnvWorkspaceRoot + "=" + s.workspaceRoot,
		config.EnvWorkspace + "=" + s.workspace,
		config.EnvWorkspaceBackend + "=overlay",
		config.EnvWorkspaceOverlayDriver + "=" + s.overlayDriver,
		config.EnvWorkspaceSession + "=" + s.workspaceSession,
		config.EnvWorkspaceCurrentRevision + "=" + s.currentWorldRevision(),
		config.EnvWorkspaceOwner + "=" + fmt.Sprintf("%t", s.workspaceOwner),
	}
}

func (s *subjectiveFS) commit() error {
	if !s.enabled || !s.workspaceOwner || s.usesDirectBackend() {
		return nil
	}
	viewRoot, err := os.MkdirTemp(s.workspaceStateDir(), "commit-")
	if err != nil {
		return fmt.Errorf("create overlay commit view: %w", err)
	}
	defer os.RemoveAll(viewRoot)

	if err := s.exportCurrentTree(viewRoot); err != nil {
		return err
	}
	if err := syncTree(viewRoot, s.workspaceRoot); err != nil {
		return err
	}
	return forceRemoveTree(s.workspaceStateDir())
}

func (s *subjectiveFS) rollback() error {
	if !s.enabled || !s.workspaceOwner || s.usesDirectBackend() {
		return nil
	}
	return forceRemoveTree(s.workspaceStateDir())
}

func (s *subjectiveFS) preflight() error {
	tempRoot, err := os.MkdirTemp(s.workspaceStateDir(), "preflight-")
	if err != nil {
		return fmt.Errorf("create overlay preflight dir: %w", err)
	}
	defer os.RemoveAll(tempRoot)

	view := filepath.Join(tempRoot, "view")
	if err := s.exportCurrentTree(view); err != nil {
		return fmt.Errorf("workspace physics unsupported in this Linux environment: %w", err)
	}
	switch s.overlayDriver {
	case "", "kernel":
		overlayMountOpts, err := preflightOverlayMount(tempRoot)
		if err != nil {
			return fmt.Errorf("workspace physics unsupported in this Linux environment: %w", err)
		}
		s.overlayMountOpts = overlayMountOpts
	case "fuse":
		if _, err := probeFuseOverlayFS(tempRoot); err != nil {
			return fmt.Errorf("workspace physics unsupported in this Linux environment: %w", err)
		}
		s.overlayMountOpts = ""
	default:
		return fmt.Errorf("unsupported workspace overlay driver %q", s.overlayDriver)
	}
	return nil
}

func preflightOverlayMount(root string) (string, error) {
	userxattrRoot := filepath.Join(root, "userxattr")
	defaultRoot := filepath.Join(root, "default")
	if err := preflightOverlayMountWithOptions(userxattrRoot, "userxattr"); err == nil {
		return "userxattr", nil
	} else if fallbackErr := preflightOverlayMountWithOptions(defaultRoot, ""); fallbackErr == nil {
		return "", nil
	} else {
		return "", fmt.Errorf("userxattr: %w; default: %w", err, fallbackErr)
	}
}

func preflightOverlayMountWithOptions(root, extraOptions string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", overlayMountPreflightScript, "sh", root, extraOptions)
	cmd.SysProcAttr = jobSysProcAttr(false, true, false)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("overlay mount preflight timed out")
	}
	if err != nil {
		return fmt.Errorf("overlay mount preflight failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

const overlayMountPreflightScript = `
set -eu
root="$1"
extra_options="${2-userxattr}"
lower="$root/lower"
upper="$root/upper"
work="$root/work"
merged="$root/merged"

mkdir -p "$lower/opaque" "$lower/swap" "$lower/dir-to-file" "$upper" "$work" "$merged"
printf 'old\n' > "$lower/opaque/old.txt"
printf 'child\n' > "$lower/swap/child.txt"
printf 'child\n' > "$lower/dir-to-file/child.txt"
printf 'leaf\n' > "$lower/file-to-dir"

cleanup() {
  cd / 2>/dev/null || true
  umount "$merged" 2>/dev/null || umount -l "$merged" 2>/dev/null || true
  chmod -R u+rwx "$upper" "$work" "$merged" 2>/dev/null || true
  rm -rf "$upper" "$work" "$merged" 2>/dev/null || true
}
trap cleanup EXIT

overlay_options="lowerdir=$lower,upperdir=$upper,workdir=$work"
if [ -n "$extra_options" ]; then
  overlay_options="$overlay_options,$extra_options"
fi
mount -t overlay overlay -o "$overlay_options" "$merged"
cd "$merged"

rm -rf opaque
mkdir opaque
printf 'new\n' > opaque/new.txt

rm -rf swap
printf 'file\n' > swap

rm -rf dir-to-file
printf 'leaf\n' > dir-to-file

rm file-to-dir
mkdir file-to-dir
printf 'child\n' > file-to-dir/child.txt

test "$(cat opaque/new.txt)" = "new"
test "$(cat swap)" = "file"
test "$(cat dir-to-file)" = "leaf"
test "$(cat file-to-dir/child.txt)" = "child"
`

func (s *subjectiveFS) exportCurrentTree(targetDir string) error {
	if err := forceRemoveTree(targetDir); err != nil && !errorsIsNotExist(err) {
		return fmt.Errorf("reset overlay export dir: %w", err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create overlay export dir: %w", err)
	}
	if _, err := os.Lstat(s.baseDir()); err == nil {
		if err := syncTree(s.baseDir(), targetDir); err != nil {
			return fmt.Errorf("seed overlay export from base snapshot: %w", err)
		}
	} else if !errorsIsNotExist(err) {
		return fmt.Errorf("stat overlay base snapshot: %w", err)
	}

	if s.canTrackWorldRevisions() {
		ledger, err := s.loadWorldRevisionLedger()
		if err != nil {
			return err
		}
		chain, err := s.revisionChain(ledger, ledger.Current)
		if err != nil {
			return err
		}
		for _, revision := range chain {
			if revision.ID == "wr0" {
				continue
			}
			if err := applyLayerDirToDisk(targetDir, s.layerDir(revision.ID)); err != nil {
				return err
			}
		}
	}
	if err := applyLayerDirToDisk(targetDir, s.liveUpperDir()); err != nil {
		return err
	}
	return nil
}

func (s *subjectiveFS) workspaceStateDir() string {
	return filepath.Join(s.dataDir, "workspaces", s.workspaceSession)
}

func (s *subjectiveFS) overlayMountStateDir() (string, error) {
	stateDir := s.workspaceStateDir()
	if len(stateDir) <= 100 {
		return stateDir, nil
	}
	// Kernel overlay mount options include every lower/upper/work path.
	// A short symlink keeps public workspace_session names stable while
	// avoiding environment-specific mount failures on long test/runtime paths.
	sum := sha256.Sum256([]byte(stateDir))
	aliasRoot := filepath.Join(os.TempDir(), "quine-overlay")
	alias := filepath.Join(aliasRoot, hex.EncodeToString(sum[:])[:16])
	if err := os.MkdirAll(aliasRoot, 0o755); err != nil {
		return "", fmt.Errorf("create overlay mount alias root: %w", err)
	}
	if target, err := os.Readlink(alias); err == nil {
		if target == stateDir {
			return alias, nil
		}
		if removeErr := os.Remove(alias); removeErr != nil {
			return "", fmt.Errorf("replace stale overlay mount alias: %w", removeErr)
		}
	} else if !errorsIsNotExist(err) {
		return "", fmt.Errorf("stat overlay mount alias: %w", err)
	}
	if err := os.Symlink(stateDir, alias); err != nil {
		if target, readErr := os.Readlink(alias); readErr == nil && target == stateDir {
			return alias, nil
		}
		return "", fmt.Errorf("create overlay mount alias: %w", err)
	}
	return alias, nil
}

func rewriteOverlayLowerDirString(lowerdir, stateDir, mountStateDir string) string {
	if stateDir == mountStateDir || strings.TrimSpace(lowerdir) == "" {
		return lowerdir
	}
	parts := strings.Split(lowerdir, ":")
	for i, part := range parts {
		parts[i] = rewriteOverlayStatePath(part, stateDir, mountStateDir)
	}
	return strings.Join(parts, ":")
}

func rewriteOverlayStatePath(path, stateDir, mountStateDir string) string {
	if stateDir == mountStateDir {
		return path
	}
	rel, err := filepath.Rel(stateDir, path)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || rel == ".." {
		return path
	}
	return filepath.Join(mountStateDir, rel)
}

func (s *subjectiveFS) stateVersionPath() string {
	return filepath.Join(s.workspaceStateDir(), "STATE_VERSION")
}

func (s *subjectiveFS) baseDir() string {
	return filepath.Join(s.workspaceStateDir(), "base")
}

func (s *subjectiveFS) layersDir() string {
	return filepath.Join(s.workspaceStateDir(), "layers")
}

func (s *subjectiveFS) layerDir(revision string) string {
	return filepath.Join(s.layersDir(), revision)
}

func (s *subjectiveFS) liveDir() string {
	return filepath.Join(s.workspaceStateDir(), "live")
}

func (s *subjectiveFS) liveUpperDir() string {
	return filepath.Join(s.liveDir(), "upper")
}

func (s *subjectiveFS) liveWorkDir() string {
	return filepath.Join(s.liveDir(), "work")
}

func (s *subjectiveFS) mountBase() string {
	return filepath.Join(s.workspaceStateDir(), "mounts")
}

func (s *subjectiveFS) ensureOverlayLayout() error {
	if err := os.MkdirAll(s.workspaceStateDir(), 0o755); err != nil {
		return fmt.Errorf("create overlay workspace state dir: %w", err)
	}

	if _, err := os.Lstat(s.stateVersionPath()); err == nil {
		data, readErr := os.ReadFile(s.stateVersionPath())
		if readErr != nil {
			return fmt.Errorf("read overlay state version: %w", readErr)
		}
		if strings.TrimSpace(string(data)) != overlayStateVersion {
			return fmt.Errorf("workspace session %q uses unsupported overlay state version %q", s.workspaceSession, strings.TrimSpace(string(data)))
		}
	} else if errorsIsNotExist(err) {
		if entries, readErr := os.ReadDir(s.workspaceStateDir()); readErr == nil && len(entries) > 0 {
			return fmt.Errorf("workspace session %q uses unsupported pre-lineage overlay state; remove %q to recreate it", s.workspaceSession, s.workspaceStateDir())
		}
		if writeErr := os.WriteFile(s.stateVersionPath(), []byte(overlayStateVersion+"\n"), 0o644); writeErr != nil {
			return fmt.Errorf("write overlay state version: %w", writeErr)
		}
	} else {
		return fmt.Errorf("stat overlay state version: %w", err)
	}

	if err := os.MkdirAll(s.layersDir(), 0o755); err != nil {
		return fmt.Errorf("create overlay layers dir: %w", err)
	}
	if err := os.MkdirAll(s.mountBase(), 0o755); err != nil {
		return fmt.Errorf("create overlay mount base: %w", err)
	}
	if err := s.ensureBaseSnapshot(); err != nil {
		return err
	}
	if err := s.resetLiveState(); err != nil {
		return err
	}
	return nil
}

func (s *subjectiveFS) ensureBaseSnapshot() error {
	if _, err := os.Lstat(s.baseDir()); err == nil {
		return nil
	} else if !errorsIsNotExist(err) {
		return fmt.Errorf("stat overlay base dir: %w", err)
	}
	if err := syncTree(s.workspaceRoot, s.baseDir()); err != nil {
		return fmt.Errorf("seed overlay base snapshot: %w", err)
	}
	return nil
}

func (s *subjectiveFS) resetLiveState() error {
	if err := forceRemoveTree(s.liveDir()); err != nil {
		return fmt.Errorf("reset overlay live dir: %w", err)
	}
	if err := os.MkdirAll(s.liveUpperDir(), 0o755); err != nil {
		return fmt.Errorf("create overlay live upperdir: %w", err)
	}
	if err := os.MkdirAll(s.liveWorkDir(), 0o755); err != nil {
		return fmt.Errorf("create overlay live workdir: %w", err)
	}
	return nil
}

func (s *subjectiveFS) bootstrapWorkspaceState() error {
	if s.usesDirectBackend() {
		return nil
	}
	if s.bootstrapSource == "" || s.bootstrapSource == s.workspaceSession {
		return nil
	}

	target := s.workspaceStateDir()
	if _, err := os.Lstat(target); err == nil {
		return nil
	} else if !errorsIsNotExist(err) {
		return fmt.Errorf("stat target workspace state dir: %w", err)
	}

	source := filepath.Join(s.dataDir, "workspaces", s.bootstrapSource)
	if _, err := os.Lstat(source); err != nil {
		if errorsIsNotExist(err) {
			return fmt.Errorf("bootstrap workspace %q: source state %q does not exist", s.workspaceSession, s.bootstrapSource)
		}
		return fmt.Errorf("stat source workspace state dir: %w", err)
	}
	versionData, err := os.ReadFile(filepath.Join(source, "STATE_VERSION"))
	if err != nil {
		if errorsIsNotExist(err) {
			return fmt.Errorf("bootstrap workspace %q: source state %q is pre-lineage overlay state", s.workspaceSession, s.bootstrapSource)
		}
		return fmt.Errorf("read source overlay state version: %w", err)
	}
	if strings.TrimSpace(string(versionData)) != overlayStateVersion {
		return fmt.Errorf("bootstrap workspace %q: source state %q uses unsupported overlay state version %q", s.workspaceSession, s.bootstrapSource, strings.TrimSpace(string(versionData)))
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return fmt.Errorf("create target workspace state dir: %w", err)
	}
	if err := os.WriteFile(filepath.Join(target, "STATE_VERSION"), versionData, 0o644); err != nil {
		return fmt.Errorf("write target overlay state version: %w", err)
	}
	if err := copyWorkspaceStateDir(filepath.Join(source, "base"), filepath.Join(target, "base")); err != nil {
		return err
	}
	if err := copyWorkspaceStateDir(filepath.Join(source, "layers"), filepath.Join(target, "layers")); err != nil {
		return err
	}
	if err := copyWorkspaceStateFile(filepath.Join(source, "world-revisions.json"), filepath.Join(target, "world-revisions.json")); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(target, "live", "upper"), 0o755); err != nil {
		return fmt.Errorf("create bootstrapped overlay live upperdir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(target, "live", "work"), 0o755); err != nil {
		return fmt.Errorf("create bootstrapped overlay live workdir: %w", err)
	}
	return nil
}

func copyWorkspaceStateDir(sourceDir, targetDir string) error {
	if _, err := os.Lstat(sourceDir); err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat workspace state dir %q: %w", sourceDir, err)
	}
	if err := syncTree(sourceDir, targetDir); err != nil {
		return fmt.Errorf("copy workspace state dir %q: %w", sourceDir, err)
	}
	return nil
}

func copyWorkspaceStateFile(sourcePath, targetPath string) error {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read workspace state file %q: %w", sourcePath, err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("create workspace state file parent: %w", err)
	}
	if err := os.WriteFile(targetPath, data, 0o644); err != nil {
		return fmt.Errorf("write workspace state file %q: %w", targetPath, err)
	}
	return nil
}

func (s *subjectiveFS) worldRevisionLedgerPath() string {
	return filepath.Join(s.workspaceStateDir(), "world-revisions.json")
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
	if err := os.RemoveAll(s.layerDir(revision.ID)); err != nil {
		return worldRevision{}, fmt.Errorf("reset overlay layer dir %s: %w", revision.ID, err)
	}
	if err := os.Rename(s.liveUpperDir(), s.layerDir(revision.ID)); err != nil {
		return worldRevision{}, fmt.Errorf("freeze overlay layer %s: %w", revision.ID, err)
	}
	if err := os.MkdirAll(s.liveUpperDir(), 0o755); err != nil {
		return worldRevision{}, fmt.Errorf("recreate overlay live upperdir: %w", err)
	}
	if err := forceRemoveTree(s.liveWorkDir()); err != nil {
		return worldRevision{}, fmt.Errorf("reset overlay live workdir: %w", err)
	}
	if err := os.MkdirAll(s.liveWorkDir(), 0o755); err != nil {
		return worldRevision{}, fmt.Errorf("recreate overlay live workdir: %w", err)
	}
	ledger.Next++
	ledger.Current = revision.ID
	ledger.Revisions[revision.ID] = revision
	if err := s.saveWorldRevisionLedger(ledger); err != nil {
		return worldRevision{}, err
	}
	return revision, nil
}

func forceRemoveTree(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		if info.Mode()&os.ModeSymlink == 0 {
			_ = os.Chmod(path, 0o600)
		}
		if err := os.Remove(path); err != nil && !errorsIsNotExist(err) {
			return err
		}
		return nil
	}
	if err := os.Chmod(path, 0o700); err != nil && !errorsIsNotExist(err) {
		return err
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if err := forceRemoveTree(filepath.Join(path, entry.Name())); err != nil {
			return err
		}
	}
	if err := os.Remove(path); err != nil && !errorsIsNotExist(err) {
		return err
	}
	return nil
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

func (s *subjectiveFS) switchWorld(target string) (string, string, error) {
	if !s.canRestoreWorld() {
		return "", "", nil
	}
	if foreignSession, foreignRevision, ok := parseWorldHandle(target); ok {
		return s.adoptForeignWorld(foreignSession, foreignRevision)
	}
	return s.switchLocalRevision(target)
}

func (s *subjectiveFS) switchLocalRevision(revision string) (string, string, error) {
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return "", "", err
	}
	if _, ok := ledger.Revisions[revision]; !ok {
		return "", "", fmt.Errorf("world revision %s does not exist", revision)
	}
	previous := ledger.Current
	if err := s.resetLiveState(); err != nil {
		return "", "", err
	}
	ledger.Current = revision
	if err := s.saveWorldRevisionLedger(ledger); err != nil {
		return "", "", err
	}
	return previous, revision, nil
}

func (s *subjectiveFS) adoptForeignWorld(foreignSession, foreignRevision string) (string, string, error) {
	if strings.TrimSpace(foreignSession) == "" || strings.TrimSpace(foreignRevision) == "" {
		return "", "", fmt.Errorf("foreign world handle is incomplete")
	}
	if foreignSession == s.workspaceSession {
		return s.switchLocalRevision(foreignRevision)
	}

	source := *s
	source.workspaceSession = foreignSession
	source.initialized = true

	sourceLedger, err := source.loadWorldRevisionLedger()
	if err != nil {
		return "", "", fmt.Errorf("load foreign world ledger: %w", err)
	}
	chain, err := source.revisionChain(sourceLedger, foreignRevision)
	if err != nil {
		return "", "", err
	}

	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return "", "", err
	}
	previous := ledger.Current
	parent := ""
	importedCurrent := ""
	for _, foreign := range chain {
		if foreign.ID == "wr0" {
			continue
		}
		local := worldRevision{
			ID:             "wr" + strconv.Itoa(ledger.Next),
			Parent:         parent,
			Kind:           "import",
			SourceSession:  foreignSession,
			SourceRevision: foreign.ID,
		}
		if err := os.RemoveAll(s.layerDir(local.ID)); err != nil {
			return "", "", fmt.Errorf("reset imported layer %s: %w", local.ID, err)
		}
		if err := copyWorkspaceStateDir(source.layerDir(foreign.ID), s.layerDir(local.ID)); err != nil {
			return "", "", fmt.Errorf("import layer %s from %s: %w", foreign.ID, foreignSession, err)
		}
		ledger.Revisions[local.ID] = local
		ledger.Next++
		parent = local.ID
		importedCurrent = local.ID
	}
	if importedCurrent == "" {
		importedCurrent = "wr0"
	}
	if err := s.resetLiveState(); err != nil {
		return "", "", err
	}
	ledger.Current = importedCurrent
	if err := s.saveWorldRevisionLedger(ledger); err != nil {
		return "", "", err
	}
	return previous, importedCurrent, nil
}

func (s *subjectiveFS) finalizeTurn(kind string, turnID int) (turnFinalizeResult, error) {
	if s.usesDirectBackend() {
		return s.finalizeDirectTurn()
	}
	mutations, changed, err := s.currentTurnMutationBlock()
	if err != nil {
		return turnFinalizeResult{}, err
	}
	if !changed {
		if err := s.resetLiveState(); err != nil {
			return turnFinalizeResult{}, err
		}
		revision, err := s.loadCurrentWorldRevision()
		if err != nil {
			return turnFinalizeResult{}, err
		}
		return turnFinalizeResult{
			Mutations: mutations,
			Revision:  revision,
		}, nil
	}

	revision, err := s.captureWorldRevision(kind, turnID)
	if err != nil {
		return turnFinalizeResult{}, err
	}
	return turnFinalizeResult{
		Changed:   true,
		Mutations: mutations,
		Revision:  revision,
	}, nil
}

func (s *subjectiveFS) importHostWorkspaceChanges(kind string, turnID int) (turnFinalizeResult, error) {
	if s.usesDirectBackend() {
		return s.finalizeDirectTurn()
	}
	if !s.canTrackWorldRevisions() {
		return turnFinalizeResult{}, nil
	}

	baseTree, err := snapshotTree(s.baseDir())
	if err != nil {
		return turnFinalizeResult{}, fmt.Errorf("snapshot overlay base tree: %w", err)
	}
	hostTree, err := snapshotTree(s.workspaceRoot)
	if err != nil {
		return turnFinalizeResult{}, fmt.Errorf("snapshot host workspace tree: %w", err)
	}
	observedHostTree, err := s.loadObservedHostTree(baseTree)
	if err != nil {
		return turnFinalizeResult{}, err
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return turnFinalizeResult{}, err
	}
	visibleTree, err := s.visibleTreeForLedgerRevision(ledger, ledger.Current)
	if err != nil {
		return turnFinalizeResult{}, err
	}

	candidates := diffTree(observedHostTree, hostTree)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Path < candidates[j].Path
	})
	skipDescendants := make([]string, 0)
	for _, mutation := range candidates {
		if pathHasSkippedAncestor(mutation.Path, skipDescendants) {
			continue
		}

		hostEntry, hostExists := hostTree[mutation.Path]
		visibleEntry, visibleExists := visibleTree[mutation.Path]
		if hostExists && visibleExists && hostEntry == visibleEntry {
			continue
		}
		if !hostExists && !visibleExists {
			continue
		}

		if hostExists {
			if err := s.copyHostEntryToLiveUpper(mutation.Path, hostEntry); err != nil {
				return turnFinalizeResult{}, fmt.Errorf("import host workspace path %q: %w", mutation.Path, err)
			}
			if hostEntry.Kind != "dir" {
				skipDescendants = append(skipDescendants, mutation.Path)
			}
			continue
		}

		if err := s.createLiveUpperWhiteout(mutation.Path); err != nil {
			return turnFinalizeResult{}, fmt.Errorf("import host workspace deletion %q: %w", mutation.Path, err)
		}
		if observedEntry, ok := observedHostTree[mutation.Path]; ok && observedEntry.Kind == "dir" {
			skipDescendants = append(skipDescendants, mutation.Path)
		}
	}

	mutations, changed, err := s.currentTurnMutationBlock()
	if err != nil {
		return turnFinalizeResult{}, err
	}
	if !changed {
		if err := s.resetLiveState(); err != nil {
			return turnFinalizeResult{}, err
		}
		revision, err := s.loadCurrentWorldRevision()
		if err != nil {
			return turnFinalizeResult{}, err
		}
		if err := s.saveObservedHostTree(hostTree); err != nil {
			return turnFinalizeResult{}, err
		}
		return turnFinalizeResult{
			Mutations: mutations,
			Revision:  revision,
		}, nil
	}

	revision, err := s.captureWorldRevision(kind, turnID)
	if err != nil {
		return turnFinalizeResult{}, err
	}
	if err := s.saveObservedHostTree(hostTree); err != nil {
		return turnFinalizeResult{}, err
	}
	return turnFinalizeResult{
		Changed:   true,
		Mutations: mutations,
		Revision:  revision,
	}, nil
}

func (s *subjectiveFS) observeHostWorkspaceChanges() error {
	if s.usesDirectBackend() {
		current, err := snapshotTree(s.workspace)
		if err != nil {
			return fmt.Errorf("snapshot direct host workspace: %w", err)
		}
		return s.saveObservedTree(current)
	}
	if !s.canTrackWorldRevisions() {
		return nil
	}
	hostTree, err := snapshotTree(s.workspaceRoot)
	if err != nil {
		return fmt.Errorf("snapshot host workspace tree: %w", err)
	}
	return s.saveObservedHostTree(hostTree)
}

func (s *subjectiveFS) observedHostTreePath() string {
	return filepath.Join(s.workspaceStateDir(), "host-observed.json")
}

func (s *subjectiveFS) loadObservedHostTree(baseTree map[string]treeEntry) (map[string]treeEntry, error) {
	data, err := os.ReadFile(s.observedHostTreePath())
	if err != nil {
		if errorsIsNotExist(err) {
			return cloneTree(baseTree), nil
		}
		return nil, fmt.Errorf("read observed host tree: %w", err)
	}
	var tree map[string]treeEntry
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, fmt.Errorf("parse observed host tree: %w", err)
	}
	if tree == nil {
		tree = map[string]treeEntry{}
	}
	return tree, nil
}

func (s *subjectiveFS) saveObservedHostTree(tree map[string]treeEntry) error {
	if tree == nil {
		tree = map[string]treeEntry{}
	}
	data, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("marshal observed host tree: %w", err)
	}
	if err := os.WriteFile(s.observedHostTreePath(), data, 0o644); err != nil {
		return fmt.Errorf("write observed host tree: %w", err)
	}
	return nil
}

func cloneTree(tree map[string]treeEntry) map[string]treeEntry {
	cloned := make(map[string]treeEntry, len(tree))
	for path, entry := range tree {
		cloned[path] = entry
	}
	return cloned
}

func pathHasSkippedAncestor(path string, ancestors []string) bool {
	slashed := filepath.ToSlash(path)
	for _, ancestor := range ancestors {
		ancestor = filepath.ToSlash(ancestor)
		if slashed == ancestor || strings.HasPrefix(slashed, ancestor+"/") {
			return true
		}
	}
	return false
}

func (s *subjectiveFS) copyHostEntryToLiveUpper(rel string, entry treeEntry) error {
	sourcePath := filepath.Join(s.workspaceRoot, rel)
	targetPath := filepath.Join(s.liveUpperDir(), rel)
	if err := forceRemoveTree(targetPath); err != nil {
		return err
	}

	switch entry.Kind {
	case "dir":
		if err := os.MkdirAll(targetPath, entry.Mode.Perm()); err != nil {
			return err
		}
		return os.Chmod(targetPath, entry.Mode.Perm())
	case "file":
		return copyFile(sourcePath, targetPath, entry.Mode.Perm())
	case "symlink":
		link, err := os.Readlink(sourcePath)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}
		return os.Symlink(link, targetPath)
	default:
		return nil
	}
}

func (s *subjectiveFS) createLiveUpperWhiteout(rel string) error {
	targetPath := filepath.Join(s.liveUpperDir(), rel)
	if err := forceRemoveTree(targetPath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	if strings.Contains(s.overlayMountOpts, "userxattr") {
		if err := createXattrOverlayWhiteout(targetPath, "user.overlay.whiteout"); err == nil {
			return nil
		}
		_ = os.Remove(targetPath)
	}
	if err := syscall.Mknod(targetPath, syscall.S_IFCHR|0o600, 0); err == nil {
		return nil
	}
	if err := createXattrOverlayWhiteout(targetPath, "trusted.overlay.whiteout"); err == nil {
		return nil
	}
	_ = os.Remove(targetPath)
	if err := createXattrOverlayWhiteout(targetPath, "user.overlay.whiteout"); err == nil {
		return nil
	}
	_ = os.Remove(targetPath)

	whiteoutPath := filepath.Join(filepath.Dir(targetPath), ".wh."+filepath.Base(targetPath))
	if err := forceRemoveTree(whiteoutPath); err != nil {
		return err
	}
	return os.WriteFile(whiteoutPath, nil, 0o600)
}

func createXattrOverlayWhiteout(path string, attr string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if closeErr := file.Close(); closeErr != nil {
		_ = os.Remove(path)
		return closeErr
	}
	if err := syscall.Setxattr(path, attr, []byte("y"), 0); err != nil {
		_ = os.Remove(path)
		return err
	}
	return nil
}

func (s *subjectiveFS) currentTurnMutationBlock() (string, bool, error) {
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return "", false, err
	}
	visible, err := s.visibleTreeForLedgerRevision(ledger, ledger.Current)
	if err != nil {
		return "", false, err
	}
	mutations, err := s.liveUpperMutations(visible)
	if err != nil {
		return "", false, err
	}
	return formatMutations(mutations), len(mutations) > 0, nil
}

func (s *subjectiveFS) restoreMutationBlock(previous, current string) (string, error) {
	if strings.TrimSpace(previous) == "" || previous == current {
		return formatMutations(nil), nil
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return "", err
	}
	beforeTree, err := s.visibleTreeForLedgerRevision(ledger, previous)
	if err != nil {
		return "", err
	}
	afterTree, err := s.visibleTreeForLedgerRevision(ledger, current)
	if err != nil {
		return "", err
	}
	return formatMutations(diffTree(beforeTree, afterTree)), nil
}

func (s *subjectiveFS) currentLowerDirString() (string, error) {
	if !s.canTrackWorldRevisions() {
		return s.baseDir(), nil
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return "", err
	}
	return s.lowerDirStringForLedgerRevision(ledger, ledger.Current)
}

func (s *subjectiveFS) readOverlayWorkspaceFile(rel string) ([]byte, error) {
	rel = filepath.Clean(rel)
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return nil, os.ErrNotExist
	}
	ledger, err := s.loadWorldRevisionLedger()
	if err != nil {
		return nil, err
	}
	visible, err := s.visibleTreeForLedgerRevision(ledger, ledger.Current)
	if err != nil {
		return nil, err
	}
	entry, ok := visible[rel]
	if !ok {
		return nil, os.ErrNotExist
	}
	if entry.Kind != "file" {
		return nil, fmt.Errorf("workspace path %q is %s, not a file", rel, entry.Kind)
	}

	chain, err := s.revisionChain(ledger, ledger.Current)
	if err != nil {
		return nil, err
	}
	for i := len(chain) - 1; i >= 0; i-- {
		rev := chain[i]
		if rev.ID == "wr0" {
			continue
		}
		candidate := filepath.Join(s.layerDir(rev.ID), rel)
		if data, err := os.ReadFile(candidate); err == nil {
			return data, nil
		} else if !errorsIsNotExist(err) {
			return nil, err
		}
	}

	return os.ReadFile(filepath.Join(s.baseDir(), rel))
}

func (s *subjectiveFS) lowerDirStringForLedgerRevision(ledger worldRevisionLedger, revision string) (string, error) {
	chain, err := s.revisionChain(ledger, revision)
	if err != nil {
		return "", err
	}
	dirs := make([]string, 0, len(chain))
	for i := len(chain) - 1; i >= 0; i-- {
		if chain[i].ID == "wr0" {
			continue
		}
		dirs = append(dirs, s.layerDir(chain[i].ID))
	}
	dirs = append(dirs, s.baseDir())
	return strings.Join(dirs, ":"), nil
}

func (s *subjectiveFS) visibleTreeForLedgerRevision(ledger worldRevisionLedger, revision string) (map[string]treeEntry, error) {
	visible, err := snapshotTree(s.baseDir())
	if err != nil {
		return nil, fmt.Errorf("capture overlay base tree: %w", err)
	}
	chain, err := s.revisionChain(ledger, revision)
	if err != nil {
		return nil, err
	}
	for _, rev := range chain {
		if rev.ID == "wr0" {
			continue
		}
		if err := s.applyLayerDir(visible, s.layerDir(rev.ID)); err != nil {
			return nil, err
		}
	}
	return visible, nil
}

func (s *subjectiveFS) revisionChain(ledger worldRevisionLedger, revision string) ([]worldRevision, error) {
	current, ok := ledger.Revisions[revision]
	if !ok {
		return nil, fmt.Errorf("world revision %s does not exist", revision)
	}
	chain := []worldRevision{current}
	for current.Parent != "" {
		parent, ok := ledger.Revisions[current.Parent]
		if !ok {
			return nil, fmt.Errorf("world revision %s has missing parent %s", current.ID, current.Parent)
		}
		chain = append(chain, parent)
		current = parent
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain, nil
}

func (s *subjectiveFS) liveUpperMutations(visible map[string]treeEntry) ([]fsMutation, error) {
	mutationKinds := make(map[string]string)
	err := filepath.WalkDir(s.liveUpperDir(), func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(s.liveUpperDir(), path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if target, ok, err := overlayWhiteoutTarget(path, rel, d); err != nil {
			return err
		} else if ok {
			if target != "" {
				recordDeletion(mutationKinds, visible, target)
			}
			return nil
		}

		if d.IsDir() {
			opaque, err := overlayDirIsOpaque(path)
			if err != nil {
				return err
			}
			if opaque {
				recordDescendantDeletions(mutationKinds, visible, rel)
			}
		}

		entry, err := buildTreeEntry(path, d)
		if err != nil {
			return err
		}
		previous, hadPrevious := visible[rel]
		if entry.Kind != "dir" {
			recordDescendantDeletions(mutationKinds, visible, rel)
		}
		switch {
		case !hadPrevious:
			mutationKinds[rel] = "created"
		case previous != entry:
			mutationKinds[rel] = "modified"
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(mutationKinds) == 0 {
		return nil, nil
	}
	paths := make([]string, 0, len(mutationKinds))
	for path := range mutationKinds {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	mutations := make([]fsMutation, 0, len(paths))
	for _, path := range paths {
		mutations = append(mutations, fsMutation{
			Path: path,
			Kind: mutationKinds[path],
		})
	}
	return mutations, nil
}

func (s *subjectiveFS) applyLayerDir(visible map[string]treeEntry, layerDir string) error {
	return filepath.WalkDir(layerDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(layerDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if target, ok, err := overlayWhiteoutTarget(path, rel, d); err != nil {
			return err
		} else if ok {
			if target != "" {
				deleteSubtree(visible, target)
			}
			return nil
		}

		if d.IsDir() {
			opaque, err := overlayDirIsOpaque(path)
			if err != nil {
				return err
			}
			if opaque {
				deleteDescendants(visible, rel)
			}
		}

		entry, err := buildTreeEntry(path, d)
		if err != nil {
			return err
		}
		if entry.Kind != "dir" {
			deleteSubtree(visible, rel)
		}
		visible[rel] = entry
		return nil
	})
}

func recordDeletion(mutationKinds map[string]string, visible map[string]treeEntry, target string) {
	recordDescendantDeletions(mutationKinds, visible, target)
	if _, ok := visible[target]; ok {
		mutationKinds[target] = "deleted"
	}
}

func recordDescendantDeletions(mutationKinds map[string]string, visible map[string]treeEntry, root string) {
	prefix := filepath.ToSlash(root) + "/"
	for path := range visible {
		slashed := filepath.ToSlash(path)
		if slashed == filepath.ToSlash(root) || strings.HasPrefix(slashed, prefix) {
			if slashed != filepath.ToSlash(root) {
				mutationKinds[path] = "deleted"
			}
		}
	}
}

func deleteSubtree(visible map[string]treeEntry, root string) {
	delete(visible, root)
	deleteDescendants(visible, root)
}

func deleteDescendants(visible map[string]treeEntry, root string) {
	prefix := filepath.ToSlash(root) + "/"
	for path := range visible {
		if strings.HasPrefix(filepath.ToSlash(path), prefix) {
			delete(visible, path)
		}
	}
}

func overlayWhiteoutTarget(path string, rel string, d os.DirEntry) (string, bool, error) {
	base := filepath.Base(rel)
	switch {
	case base == ".wh..wh..opq":
		return "", true, nil
	case strings.HasPrefix(base, ".wh."):
		target := filepath.Join(filepath.Dir(rel), strings.TrimPrefix(base, ".wh."))
		if target == "." {
			target = strings.TrimPrefix(base, ".wh.")
		}
		return target, true, nil
	}

	info, err := d.Info()
	if err != nil {
		return "", false, err
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Rdev == 0 {
			return rel, true, nil
		}
	}
	whiteout, err := overlayHasBoolXattr(path, "trusted.overlay.whiteout", "user.overlay.whiteout")
	if err != nil {
		return "", false, err
	}
	if whiteout {
		return rel, true, nil
	}
	return "", false, nil
}

func overlayDirIsOpaque(path string) (bool, error) {
	if opaque, err := overlayHasBoolXattr(path, "trusted.overlay.opaque", "user.overlay.opaque"); err != nil {
		return false, err
	} else if opaque {
		return true, nil
	}
	if _, err := os.Lstat(filepath.Join(path, ".wh..wh..opq")); err == nil {
		return true, nil
	} else if !errorsIsNotExist(err) {
		return false, err
	}
	return false, nil
}

func applyLayerDirToDisk(targetRoot string, layerDir string) error {
	if _, err := os.Lstat(layerDir); err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat overlay layer %q: %w", layerDir, err)
	}
	return filepath.WalkDir(layerDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(layerDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if target, ok, err := overlayWhiteoutTarget(path, rel, d); err != nil {
			return err
		} else if ok {
			if target != "" {
				if err := os.RemoveAll(filepath.Join(targetRoot, target)); err != nil && !errorsIsNotExist(err) {
					return fmt.Errorf("apply overlay deletion %q: %w", target, err)
				}
			}
			return nil
		}

		targetPath := filepath.Join(targetRoot, rel)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			if err := os.RemoveAll(targetPath); err != nil && !errorsIsNotExist(err) {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, targetPath)
		case info.IsDir():
			if existing, err := os.Lstat(targetPath); err == nil && !existing.IsDir() {
				if err := os.RemoveAll(targetPath); err != nil {
					return err
				}
			} else if err != nil && !errorsIsNotExist(err) {
				return err
			}
			if err := os.MkdirAll(targetPath, info.Mode().Perm()); err != nil {
				return err
			}
			opaque, err := overlayDirIsOpaque(path)
			if err != nil {
				return err
			}
			if opaque {
				if err := deleteDescendantsOnDisk(targetPath); err != nil {
					return err
				}
			}
			return os.Chmod(targetPath, info.Mode().Perm())
		case info.Mode().IsRegular():
			if err := os.RemoveAll(targetPath); err != nil && !errorsIsNotExist(err) {
				return err
			}
			return copyFile(path, targetPath, info.Mode().Perm())
		default:
			return nil
		}
	})
}

func deleteDescendantsOnDisk(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		if errorsIsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(root, entry.Name())); err != nil && !errorsIsNotExist(err) {
			return err
		}
	}
	return nil
}

func overlayHasBoolXattr(path string, attrs ...string) (bool, error) {
	for _, attr := range attrs {
		size, err := syscall.Getxattr(path, attr, nil)
		switch {
		case err == nil:
			if size == 0 {
				return false, nil
			}
		case err == syscall.ENODATA || err == syscall.EPERM || err == syscall.EOPNOTSUPP || err == syscall.ENOTSUP:
			continue
		default:
			return false, err
		}

		buf := make([]byte, size)
		n, err := syscall.Getxattr(path, attr, buf)
		switch {
		case err == nil:
			value := strings.ToLower(strings.TrimSpace(string(buf[:n])))
			return value == "y" || value == "1" || value == "true", nil
		case err == syscall.ENODATA || err == syscall.EPERM || err == syscall.EOPNOTSUPP || err == syscall.ENOTSUP:
			continue
		default:
			return false, err
		}
	}
	return false, nil
}
