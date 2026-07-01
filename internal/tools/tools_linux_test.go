//go:build linux

package tools

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func TestWorkspaceOverlayStartJobRequestsMountNamespace(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-overlay-newns"},
		Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"},
		Limits:    config.Limits{OutputTruncate: 20480, ShTimeout: 10},
		WorkspaceConfig: config.WorkspaceConfig{
			WorkspaceEnabled:      true,
			WorkspaceRoot:         root,
			Workspace:             root,
			WorkspaceBackend:      "overlay",
			WorkspaceRevisionMode: config.WorkspaceRevisionRestore,
			WorkspaceSession:      "workspace-overlay-newns-session",
			WorkspaceOwner:        true,
		},
		Paths: config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
	}

	b := NewShExecutor(cfg, nil)
	if err := b.Prepare(); err != nil {
		t.Skipf("workspace physics unsupported in this Linux environment: %v", err)
	}
	job, err := b.startJob("true", false, "")
	if err != nil {
		t.Fatalf("startJob() error: %v", err)
	}
	defer b.Close(false)

	if job.cmd.SysProcAttr == nil {
		t.Fatal("job SysProcAttr should not be nil")
	}
	if job.cmd.SysProcAttr.Cloneflags&syscall.CLONE_NEWNS == 0 {
		t.Fatalf("overlay workspace startJob should request CLONE_NEWNS: %#x", job.cmd.SysProcAttr.Cloneflags)
	}
	<-job.doneCh
	_ = os.RemoveAll(job.canonicalDir)
}

func TestWorkspaceJobWrapperLocalizesMountIsolation(t *testing.T) {
	if strings.Contains(jobWrapperScript, "mount --make-rprivate /") || strings.Contains(jobWrapperScript, "mount --make-private /") {
		t.Fatalf("job wrapper should not mutate global root mount propagation:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `mount --make-rslave /`) {
		t.Fatalf("job wrapper should make the child mount namespace propagation one-way before workspace mounts:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `mount --bind "$QUINE_WORKSPACE_ROOT" "$QUINE_WORKSPACE_ROOT"`) {
		t.Fatalf("job wrapper should self-bind the workspace root before mounting overlay:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `mount --make-private "$QUINE_WORKSPACE_ROOT"`) {
		t.Fatalf("job wrapper should privatize the workspace mountpoint:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `QUINE_WORKSPACE_OVERLAY_EXTRA_OPTS`) {
		t.Fatalf("job wrapper should honor preflight-selected overlay metadata options:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `QUINE_WORKSPACE_OVERLAY_DRIVER:-kernel`) {
		t.Fatalf("job wrapper should default overlay mount driver to kernel:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `fuse-overlayfs -f -o "$overlay_options" "$merged_dir"`) {
		t.Fatalf("job wrapper should support forced fuse-overlayfs mounts:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `fusermount3 -u "$merged_dir"`) {
		t.Fatalf("job wrapper should unmount fuse-overlayfs mounts through fusermount when available:\n%s", jobWrapperScript)
	}
}

func TestWorkspaceJobWrapperStopsWhenWorkspaceEntryFails(t *testing.T) {
	if got := strings.Count(jobWrapperScript, `enter_workspace || return "$?"`); got != 2 {
		t.Fatalf("job wrapper should stop both sync and interactive jobs when workspace entry fails, got %d:\n%s", got, jobWrapperScript)
	}
}

func TestWorkspaceOverlayPreflightProbesUserNamespaceWhiteouts(t *testing.T) {
	if _, err := preflightOverlayMount(t.TempDir()); err != nil {
		t.Skipf("workspace overlay mount unsupported in this Linux environment: %v", err)
	}
}

func TestWorkspaceOverlayPreflightScriptCoversDeletionAndReplacement(t *testing.T) {
	checks := []string{
		`userxattr`,
		`rm -rf opaque`,
		`rm -rf swap`,
		`rm -rf dir-to-file`,
		`rm file-to-dir`,
	}
	for _, want := range checks {
		if !strings.Contains(overlayMountPreflightScript, want) {
			t.Fatalf("overlay preflight script should contain %q:\n%s", want, overlayMountPreflightScript)
		}
	}
}

func TestWorkspaceJobWrapperCleanupAvoidsBlockingUnmounts(t *testing.T) {
	if !strings.Contains(jobWrapperScript, `cd / 2>/dev/null || true`) {
		t.Fatalf("job wrapper cleanup should leave the workspace tree before unmounting:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `if command -v python3 >/dev/null 2>&1; then`) {
		t.Fatalf("job wrapper cleanup should prefer the python daemon helper when available:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `pid = os.fork()`) || !strings.Contains(jobWrapperScript, `os.setsid()`) {
		t.Fatalf("job wrapper cleanup should daemonize the unmount helper to avoid shell job waits:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `safe_run(["umount", "-l", target])`) {
		t.Fatalf("job wrapper cleanup should lazily unmount both workspace targets in the daemon helper:\n%s", jobWrapperScript)
	}
	if !strings.Contains(jobWrapperScript, `shutil.rmtree(mount_root, ignore_errors=True)`) {
		t.Fatalf("job wrapper cleanup should still best-effort remove the per-job mount root:\n%s", jobWrapperScript)
	}
}

func TestWorkspaceOverlayExportCurrentTreeMaterializesLineage(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "base.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "dir"), 0o755); err != nil {
		t.Fatalf("mkdir base dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "dir", "old.txt"), []byte("old\n"), 0o644); err != nil {
		t.Fatalf("write old file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "swap"), 0o755); err != nil {
		t.Fatalf("mkdir swap dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "swap", "nested.txt"), []byte("nested\n"), 0o644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	s := &subjectiveFS{
		enabled:          true,
		dataDir:          t.TempDir(),
		workspaceRoot:    root,
		workspace:        root,
		workspaceSession: "export-lineage-session",
	}
	if err := s.ensureOverlayLayout(); err != nil {
		t.Fatalf("ensureOverlayLayout() error: %v", err)
	}
	if err := s.ensureWorldRevisionLedger(); err != nil {
		t.Fatalf("ensureWorldRevisionLedger() error: %v", err)
	}

	if err := os.MkdirAll(s.layerDir("wr1"), 0o755); err != nil {
		t.Fatalf("mkdir wr1 layer: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.layerDir("wr1"), ".wh.base.txt"), nil, 0o644); err != nil {
		t.Fatalf("write whiteout: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(s.layerDir("wr1"), "dir"), 0o755); err != nil {
		t.Fatalf("mkdir wr1 dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.layerDir("wr1"), "dir", ".wh..wh..opq"), nil, 0o644); err != nil {
		t.Fatalf("write opaque marker: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.layerDir("wr1"), "dir", "new.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatalf("write new file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.layerDir("wr1"), "swap"), []byte("now-file\n"), 0o644); err != nil {
		t.Fatalf("write dir-to-file replacement: %v", err)
	}

	ledger := worldRevisionLedger{
		Current: "wr1",
		Next:    2,
		Revisions: map[string]worldRevision{
			"wr0": {ID: "wr0", Kind: "baseline"},
			"wr1": {ID: "wr1", Parent: "wr0", Kind: "sh", Turn: 1},
		},
	}
	if err := s.saveWorldRevisionLedger(ledger); err != nil {
		t.Fatalf("saveWorldRevisionLedger() error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.liveUpperDir(), "live.txt"), []byte("live\n"), 0o644); err != nil {
		t.Fatalf("write live upper file: %v", err)
	}

	exportDir := t.TempDir()
	if err := s.exportCurrentTree(exportDir); err != nil {
		t.Fatalf("exportCurrentTree() error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(exportDir, "base.txt")); !os.IsNotExist(err) {
		t.Fatalf("base.txt should be deleted by whiteout, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(exportDir, "dir", "old.txt")); !os.IsNotExist(err) {
		t.Fatalf("dir/old.txt should be removed by opaque dir, stat err = %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(exportDir, "dir", "new.txt")); err != nil {
		t.Fatalf("read dir/new.txt: %v", err)
	} else if strings.TrimSpace(string(got)) != "new" {
		t.Fatalf("dir/new.txt = %q, want new", strings.TrimSpace(string(got)))
	}
	if info, err := os.Lstat(filepath.Join(exportDir, "swap")); err != nil {
		t.Fatalf("stat swap: %v", err)
	} else if info.IsDir() {
		t.Fatal("swap should be materialized as a file")
	}
	if got, err := os.ReadFile(filepath.Join(exportDir, "swap")); err != nil {
		t.Fatalf("read swap: %v", err)
	} else if strings.TrimSpace(string(got)) != "now-file" {
		t.Fatalf("swap = %q, want now-file", strings.TrimSpace(string(got)))
	}
	if got, err := os.ReadFile(filepath.Join(exportDir, "live.txt")); err != nil {
		t.Fatalf("read live.txt: %v", err)
	} else if strings.TrimSpace(string(got)) != "live" {
		t.Fatalf("live.txt = %q, want live", strings.TrimSpace(string(got)))
	}
}

func TestForceRemoveTreeHandlesInaccessibleOverlayWorkdir(t *testing.T) {
	root := t.TempDir()
	liveDir := filepath.Join(root, "live")
	workDir := filepath.Join(liveDir, "work")
	internalDir := filepath.Join(workDir, "work")
	if err := os.MkdirAll(internalDir, 0o755); err != nil {
		t.Fatalf("mkdir overlay workdir tree: %v", err)
	}
	if err := os.WriteFile(filepath.Join(internalDir, "placeholder"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write placeholder: %v", err)
	}
	if err := os.Chmod(internalDir, 0); err != nil {
		t.Fatalf("chmod internal workdir: %v", err)
	}

	if err := forceRemoveTree(liveDir); err != nil {
		t.Fatalf("forceRemoveTree(liveDir) error: %v", err)
	}
	if _, err := os.Lstat(liveDir); !os.IsNotExist(err) {
		t.Fatalf("liveDir should be removed, got err=%v", err)
	}
}
