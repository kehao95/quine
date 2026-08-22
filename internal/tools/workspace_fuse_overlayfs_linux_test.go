//go:build linux

package tools

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/config"
)

func TestFuseOverlayFSProbeScriptCoversDeletionReplacementAndCleanup(t *testing.T) {
	checks := []string{
		`"$fuse_overlayfs" -f -o "lowerdir=$lower,upperdir=$upper,workdir=$work" "$merged"`,
		`"$fusermount_bin" -u "$merged"`,
		`rm -rf opaque`,
		`rm -rf swap`,
		`rm -rf dir-to-file`,
		`rm file-to-dir`,
		`wait_for_mount`,
		`wait_for_unmount`,
		`test "$(id -u)" = "0"`,
	}
	for _, want := range checks {
		if !strings.Contains(fuseOverlayFSMountPreflightScript, want) {
			t.Fatalf("fuse-overlayfs feasibility probe should contain %q:\n%s", want, fuseOverlayFSMountPreflightScript)
		}
	}
}

func TestProbeFuseOverlayFSCurrentHost(t *testing.T) {
	result, err := probeFuseOverlayFS(t.TempDir())
	if err != nil {
		switch {
		case errors.Is(err, errFuseOverlayFSBinaryUnavailable):
			t.Skipf("fuse-overlayfs binary unavailable in this environment: %v", err)
		case errors.Is(err, errFuseOverlayFSUnmountHelperUnavailable):
			t.Skipf("fusermount helper unavailable in this environment: %v", err)
		case errors.Is(err, errFuseOverlayFSDeviceUnavailable):
			t.Skipf("/dev/fuse unavailable in this environment: %v", err)
		default:
			t.Fatalf("probeFuseOverlayFS() error: %v", err)
		}
	}
	if result.BinaryPath == "" {
		t.Fatal("probe result should record the fuse-overlayfs binary path")
	}
	if result.FusermountPath == "" {
		t.Fatal("probe result should record the fusermount helper path")
	}
	if result.DevFusePath != "/dev/fuse" {
		t.Fatalf("probe result dev_fuse_path = %q, want /dev/fuse", result.DevFusePath)
	}
	if !result.LowPrivilegeNamespaces || !result.MountSucceeded {
		t.Fatalf("probe result = %+v, want low-privilege namespace mount success", result)
	}
	if !result.WhiteoutSemantics || !result.OpaqueDirSemantics || !result.DirReplacedByFile || !result.FileReplacedByDir {
		t.Fatalf("probe result = %+v, want semantic checks recorded as passing", result)
	}
	if !result.CleanupUnmounted || !result.CleanupRemoved {
		t.Fatalf("probe result = %+v, want cleanup checks recorded as passing", result)
	}
}

func TestApplyFuseOverlayFSLayerMaterializesReplacementDialect(t *testing.T) {
	upper := createFuseOverlayFSReplacementUpper(t)
	targetRoot := t.TempDir()
	mustWriteFile(t, filepath.Join(targetRoot, "opaque", "old.txt"), "old\n")
	mustWriteFile(t, filepath.Join(targetRoot, "swap", "child.txt"), "child\n")
	mustWriteFile(t, filepath.Join(targetRoot, "dir-to-file", "child.txt"), "child\n")
	mustWriteFile(t, filepath.Join(targetRoot, "file-to-dir"), "leaf\n")

	if err := applyLayerDirToDisk(targetRoot, upper); err != nil {
		t.Fatalf("applyLayerDirToDisk() error: %v", err)
	}

	assertFileContent(t, filepath.Join(targetRoot, "opaque", "new.txt"), "new\n")
	assertNotExist(t, filepath.Join(targetRoot, "opaque", "old.txt"))
	assertFileContent(t, filepath.Join(targetRoot, "swap"), "file\n")
	assertFileContent(t, filepath.Join(targetRoot, "dir-to-file"), "leaf\n")
	assertFileContent(t, filepath.Join(targetRoot, "file-to-dir", "child.txt"), "child\n")
}

func TestWorkspaceOverlayFuseDriverRunsAndCommits(t *testing.T) {
	root := t.TempDir()
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("realpath root: %v", err)
	}
	cfg := &config.Config{
		Identity:  config.Identity{ModelID: "claude-sonnet-4-20250514", SessionID: "workspace-fuse-driver"},
		Transport: config.Transport{APIKey: "test-key", APIBase: "https://api.example.com", Provider: "anthropic"},
		Limits:    config.Limits{OutputTruncate: 20480, ShTimeout: 10},
		Paths:     config.Paths{DataDir: t.TempDir(), Shell: "/bin/sh"},
		WorkspaceConfig: config.WorkspaceConfig{
			WorkspaceEnabled:       true,
			WorkspaceRoot:          root,
			Workspace:              root,
			WorkspaceBackend:       "overlay",
			WorkspaceOverlayDriver: "fuse",
			WorkspaceRevisionMode:  config.WorkspaceRevisionRestore,
			WorkspaceSession:       "workspace-fuse-driver-session",
			WorkspaceOwner:         true,
		},
	}

	b := NewShExecutor(cfg)
	defer b.Close(false)
	requireWorkspaceSupport(t, b)
	result := b.Execute("tool-fuse-driver", "test \"${QUINE_WORKSPACE_OVERLAY_DRIVER:-}\" = fuse && printf 'fuse ok\\n' > result.txt", 0, 0, false, false, "")
	if result.IsError {
		t.Fatalf("unexpected fuse workspace error:\n%s", result.Content)
	}
	if _, err := os.Stat(filepath.Join(rootReal, "result.txt")); !os.IsNotExist(err) {
		t.Fatalf("host file should not exist before commit, err=%v", err)
	}
	if err := b.Close(true); err != nil {
		t.Fatalf("commit close failed: %v", err)
	}
	assertFileContent(t, filepath.Join(rootReal, "result.txt"), "fuse ok\n")
}

func createFuseOverlayFSReplacementUpper(t *testing.T) string {
	t.Helper()

	fuseOverlayFSPath, err := exec.LookPath("fuse-overlayfs")
	if err != nil {
		t.Skipf("fuse-overlayfs binary unavailable in this environment: %v", err)
	}
	fusermountPath, err := lookPathOneOf("fusermount3", "fusermount")
	if err != nil {
		t.Skipf("fusermount helper unavailable in this environment: %v", err)
	}
	if _, err := os.Stat("/dev/fuse"); err != nil {
		t.Skipf("/dev/fuse unavailable in this environment: %v", err)
	}

	root := t.TempDir()
	lower := filepath.Join(root, "lower")
	upper := filepath.Join(root, "upper")
	work := filepath.Join(root, "work")
	merged := filepath.Join(root, "merged")
	mustWriteFile(t, filepath.Join(lower, "opaque", "old.txt"), "old\n")
	mustWriteFile(t, filepath.Join(lower, "swap", "child.txt"), "child\n")
	mustWriteFile(t, filepath.Join(lower, "dir-to-file", "child.txt"), "child\n")
	mustWriteFile(t, filepath.Join(lower, "file-to-dir"), "leaf\n")
	for _, dir := range []string{upper, work, merged} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", fuseOverlayFSTestMutationScript, "sh", fuseOverlayFSPath, fusermountPath, lower, upper, work, merged)
	cmd.SysProcAttr = jobSysProcAttr(false, true, false)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("fuse-overlayfs test mutation timed out")
	}
	if err != nil {
		t.Skipf("fuse-overlayfs mount unavailable in this environment: %v: %s", err, string(output))
	}
	return upper
}

func mustWriteFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir parent for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertFileContent(t *testing.T, path string, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", path, string(got), want)
	}
}

func assertNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); err == nil {
		t.Fatalf("%s exists, want absent", path)
	} else if !errorsIsNotExist(err) {
		t.Fatalf("stat %s: %v", path, err)
	}
}

const fuseOverlayFSTestMutationScript = `
set -eu
fuse_overlayfs="$1"
fusermount_bin="$2"
lower="$3"
upper="$4"
work="$5"
merged="$6"
stderr_log="$upper/fuse-overlayfs.stderr"

is_mounted() {
  grep -qs " $merged " /proc/self/mounts
}

wait_for_mount() {
  i=0
  while [ "$i" -lt 100 ]; do
    if is_mounted; then
      return 0
    fi
    if ! kill -0 "$fuse_pid" 2>/dev/null; then
      break
    fi
    sleep 0.05
    i=$((i + 1))
  done
  return 1
}

wait_for_unmount() {
  i=0
  while [ "$i" -lt 100 ]; do
    if ! is_mounted; then
      return 0
    fi
    sleep 0.05
    i=$((i + 1))
  done
  return 1
}

cleanup() {
  cd / 2>/dev/null || true
  if is_mounted; then
    "$fusermount_bin" -u "$merged" 2>/dev/null || \
      "$fusermount_bin" -uz "$merged" 2>/dev/null || \
      umount "$merged" 2>/dev/null || \
      umount -l "$merged" 2>/dev/null || true
    wait_for_unmount || true
  fi
  if [ -n "${fuse_pid:-}" ]; then
    kill "$fuse_pid" 2>/dev/null || true
    wait "$fuse_pid" 2>/dev/null || true
  fi
}
trap cleanup EXIT

"$fuse_overlayfs" -f -o "lowerdir=$lower,upperdir=$upper,workdir=$work" "$merged" >/dev/null 2>"$stderr_log" &
fuse_pid=$!
wait_for_mount || {
  cat "$stderr_log" >&2 || true
  echo "fuse-overlayfs mount did not appear" >&2
  exit 1
}

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
cd /

"$fusermount_bin" -u "$merged"
wait_for_unmount || {
  cat "$stderr_log" >&2 || true
  echo "fuse-overlayfs mount cleanup failed" >&2
  exit 1
}
kill "$fuse_pid" 2>/dev/null || true
wait "$fuse_pid" 2>/dev/null || true
fuse_pid=""
`
