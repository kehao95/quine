//go:build linux

package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	errFuseOverlayFSBinaryUnavailable        = errors.New("fuse-overlayfs binary unavailable")
	errFuseOverlayFSUnmountHelperUnavailable = errors.New("fusermount helper unavailable")
	errFuseOverlayFSDeviceUnavailable        = errors.New("/dev/fuse unavailable")
)

type fuseOverlayFSProbeResult struct {
	BinaryPath             string `json:"binary_path"`
	FusermountPath         string `json:"fusermount_path"`
	DevFusePath            string `json:"dev_fuse_path"`
	LowPrivilegeNamespaces bool   `json:"low_privilege_namespaces"`
	MountSucceeded         bool   `json:"mount_succeeded"`
	WhiteoutSemantics      bool   `json:"whiteout_semantics"`
	OpaqueDirSemantics     bool   `json:"opaque_dir_semantics"`
	DirReplacedByFile      bool   `json:"dir_replaced_by_file"`
	FileReplacedByDir      bool   `json:"file_replaced_by_dir"`
	CleanupUnmounted       bool   `json:"cleanup_unmounted"`
	CleanupRemoved         bool   `json:"cleanup_removed"`
}

func probeFuseOverlayFS(root string) (fuseOverlayFSProbeResult, error) {
	result := fuseOverlayFSProbeResult{
		LowPrivilegeNamespaces: true,
	}

	binaryPath, err := exec.LookPath("fuse-overlayfs")
	if err != nil {
		return result, fmt.Errorf("%w: %v", errFuseOverlayFSBinaryUnavailable, err)
	}
	result.BinaryPath = binaryPath

	fusermountPath, err := lookPathOneOf("fusermount3", "fusermount")
	if err != nil {
		return result, fmt.Errorf("%w: %v", errFuseOverlayFSUnmountHelperUnavailable, err)
	}
	result.FusermountPath = fusermountPath

	if _, err := os.Stat("/dev/fuse"); err != nil {
		return result, fmt.Errorf("%w: %v", errFuseOverlayFSDeviceUnavailable, err)
	}
	result.DevFusePath = "/dev/fuse"

	if err := preflightFuseOverlayFSMount(root, binaryPath, fusermountPath); err != nil {
		return result, err
	}

	result.MountSucceeded = true
	result.WhiteoutSemantics = true
	result.OpaqueDirSemantics = true
	result.DirReplacedByFile = true
	result.FileReplacedByDir = true
	result.CleanupUnmounted = true
	result.CleanupRemoved = true
	return result, nil
}

func lookPathOneOf(names ...string) (string, error) {
	for _, name := range names {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}
	return "", exec.ErrNotFound
}

func preflightFuseOverlayFSMount(root, binaryPath, fusermountPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", fuseOverlayFSMountPreflightScript, "sh", root, binaryPath, fusermountPath, strconv.Itoa(os.Getuid()), strconv.Itoa(os.Getgid()))
	cmd.SysProcAttr = jobSysProcAttr(false, true, false)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("fuse-overlayfs feasibility probe timed out")
	}
	if err != nil {
		return fmt.Errorf("fuse-overlayfs feasibility probe failed: %w: %s", err, string(output))
	}
	return nil
}

const fuseOverlayFSMountPreflightScript = `
set -eu
root="$1"
fuse_overlayfs="$2"
fusermount_bin="$3"
host_uid="$4"
host_gid="$5"
lower="$root/lower"
upper="$root/upper"
work="$root/work"
merged="$root/merged"
stdout_log="$root/fuse-overlayfs.stdout"
stderr_log="$root/fuse-overlayfs.stderr"

mkdir -p "$lower/opaque" "$lower/swap" "$lower/dir-to-file" "$upper" "$work" "$merged"
printf 'old\n' > "$lower/opaque/old.txt"
printf 'child\n' > "$lower/swap/child.txt"
printf 'child\n' > "$lower/dir-to-file/child.txt"
printf 'leaf\n' > "$lower/file-to-dir"

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
  rm -rf "$lower" "$upper" "$work" "$merged" 2>/dev/null || true
}
trap cleanup EXIT

if [ "$host_uid" != "0" ]; then
  test "$(id -u)" = "0"
fi
if [ "$host_gid" != "0" ]; then
  test "$(id -g)" = "0"
fi

"$fuse_overlayfs" -f -o "lowerdir=$lower,upperdir=$upper,workdir=$work" "$merged" >"$stdout_log" 2>"$stderr_log" &
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

test "$(cat opaque/new.txt)" = "new"
test "$(cat swap)" = "file"
test "$(cat dir-to-file)" = "leaf"
test "$(cat file-to-dir/child.txt)" = "child"

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

rm -rf "$lower" "$upper" "$work" "$merged"
test ! -e "$lower"
test ! -e "$upper"
test ! -e "$work"
test ! -e "$merged"
`
