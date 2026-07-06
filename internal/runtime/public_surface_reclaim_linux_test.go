//go:build linux

package runtime

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	fusefs "github.com/hanwen/go-fuse/v2/fs"
	fusepkg "github.com/hanwen/go-fuse/v2/fuse"
)

// Captured from a live reproduction of the post-exec stale mount (probe on
// 2026-07-03; same shape as the l0-ablation-2x2-2026-06-29 diagnosis).
const staleMountinfoFixture = `21 26 0:20 / /sys rw,nosuid,nodev,noexec,relatime shared:2 - sysfs sysfs rw
26 1 8:1 / / rw,relatime shared:1 - ext4 /dev/sda1 rw,errors=remount-ro
808 39 0:77 / /tmp/agent/sess_1/public rw,nosuid,nodev,relatime shared:578 - fuse.quine-public quine-public rw,user_id=1000,group_id=1000,max_read=131072
812 39 0:78 / /tmp/with\040space/public rw,relatime shared:580 - fuse.quine-public quine-public rw
`

func TestMountinfoEntryForMountpoint(t *testing.T) {
	fstype, ok := mountinfoEntryForMountpoint(staleMountinfoFixture, "/tmp/agent/sess_1/public")
	if !ok {
		t.Fatal("expected mountinfo entry for /tmp/agent/sess_1/public")
	}
	if fstype != "fuse.quine-public" {
		t.Fatalf("fstype = %q, want fuse.quine-public", fstype)
	}

	if _, ok := mountinfoEntryForMountpoint(staleMountinfoFixture, "/tmp/agent/sess_1"); ok {
		t.Fatal("parent of a mountpoint must not match")
	}
	if _, ok := mountinfoEntryForMountpoint(staleMountinfoFixture, "/tmp/agent/sess_1/pub"); ok {
		t.Fatal("path prefix of a mountpoint must not match")
	}

	// proc(5) octal-escapes spaces in mountinfo paths.
	fstype, ok = mountinfoEntryForMountpoint(staleMountinfoFixture, "/tmp/with space/public")
	if !ok || fstype != "fuse.quine-public" {
		t.Fatalf("escaped mountpoint: fstype=%q ok=%v, want fuse.quine-public true", fstype, ok)
	}
}

func TestUnescapeMountinfoField(t *testing.T) {
	cases := map[string]string{
		"/plain/path":       "/plain/path",
		`/with\040space`:    "/with space",
		`/tab\011and\134bs`: "/tab\tand\\bs",
		`/trailing\04`:      `/trailing\04`, // truncated escape stays literal
	}
	for in, want := range cases {
		if got := unescapeMountinfoField(in); got != want {
			t.Errorf("unescapeMountinfoField(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsDisconnectedMountStatError(t *testing.T) {
	if !isDisconnectedMountStatError(&os.PathError{Op: "stat", Path: "/x", Err: syscall.ENOTCONN}) {
		t.Fatal("wrapped ENOTCONN must classify as disconnected mount")
	}
	if !isDisconnectedMountStatError(syscall.ENOTCONN) {
		t.Fatal("bare ENOTCONN must classify as disconnected mount")
	}
	if isDisconnectedMountStatError(nil) {
		t.Fatal("nil error must not classify as disconnected mount")
	}
	if isDisconnectedMountStatError(&os.PathError{Op: "stat", Path: "/x", Err: syscall.ENOENT}) {
		t.Fatal("ENOENT must not classify as disconnected mount")
	}
}

func TestDetectStaleRuntimeSurfaceMountCleanPaths(t *testing.T) {
	dir := t.TempDir()
	if signal := detectStaleRuntimeSurfaceMount(dir); signal != "" {
		t.Fatalf("plain directory reported stale: %q", signal)
	}
	if signal := detectStaleRuntimeSurfaceMount(filepath.Join(dir, "missing")); signal != "" {
		t.Fatalf("missing path reported stale: %q", signal)
	}
}

// TestAbandonedFuseMountHelperProcess is not a test: it is the child body for
// TestBootstrapAgentRootReclaimsAbandonedPublicSurfaceMount. It mounts a FUSE
// filesystem at the directed mountpoint, then blocks until SIGKILL without
// ever unmounting — the moral equivalent of the exec handover killing the
// predecessor's in-process FUSE server.
func TestAbandonedFuseMountHelperProcess(t *testing.T) {
	mountpoint := os.Getenv("QUINE_TEST_ABANDONED_FUSE_MOUNTPOINT")
	if mountpoint == "" {
		t.Skip("helper process for TestBootstrapAgentRootReclaimsAbandonedPublicSurfaceMount")
	}
	if err := os.MkdirAll(mountpoint, 0o755); err != nil {
		fmt.Printf("HELPER_MOUNT_FAILED: %v\n", err)
		os.Exit(3)
	}
	timeout := time.Duration(0)
	if _, err := fusefs.Mount(mountpoint, &fusefs.Inode{}, &fusefs.Options{
		EntryTimeout: &timeout,
		AttrTimeout:  &timeout,
		MountOptions: fusepkg.MountOptions{
			FsName:        "quine-public",
			Name:          "quine-public",
			DisableXAttrs: true,
			DirectMount:   true,
		},
	}); err != nil {
		fmt.Printf("HELPER_MOUNT_FAILED: %v\n", err)
		os.Exit(3)
	}
	fmt.Println("HELPER_MOUNTED")
	select {}
}

// abandonFuseMountAt leaves a disconnected FUSE mount at mountpoint by
// mounting from a child process and killing it, then waits for the kernel to
// reflect the dead connection (stat → ENOTCONN).
func abandonFuseMountAt(t *testing.T, mountpoint string) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestAbandonedFuseMountHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(), "QUINE_TEST_ABANDONED_FUSE_MOUNTPOINT="+mountpoint)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("helper stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start abandoned-mount helper: %v", err)
	}

	mounted := false
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "HELPER_MOUNT_FAILED") {
			break
		}
		if strings.Contains(line, "HELPER_MOUNTED") {
			mounted = true
			break
		}
	}
	if !mounted {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatalf("abandoned-mount helper never mounted at %s", mountpoint)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill abandoned-mount helper: %v", err)
	}
	_ = cmd.Wait()

	deadline := time.Now().Add(10 * time.Second)
	for {
		_, err := os.Stat(mountpoint)
		if isDisconnectedMountStatError(err) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("mount at %s never became disconnected after helper death: stat err = %v", mountpoint, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// TestBootstrapAgentRootReclaimsAbandonedPublicSurfaceMount reproduces the
// post-exec stale-mount defect (l0-ablation-2x2-2026-06-29: the successor saw
// "Transport endpoint is not connected" on public/ and could not bootstrap)
// and asserts the fix: a disconnected predecessor mount at public/ is
// reclaimed and the fresh mount lands, with no degradation concluded. Three
// abandon→bootstrap cycles stand in for a repeated exec loop.
func TestBootstrapAgentRootReclaimsAbandonedPublicSurfaceMount(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	cfg := testCfg(t)
	publicDir := filepath.Join(cfg.AgentRoot(), "public")
	t.Cleanup(func() {
		// Never leave a stale mount on the host, even on failure.
		if helper, err := exec.LookPath("fusermount3"); err == nil {
			_ = exec.Command(helper, "-uz", publicDir).Run()
		} else if helper, err := exec.LookPath("fusermount"); err == nil {
			_ = exec.Command(helper, "-uz", publicDir).Run()
		}
	})

	for generation := 1; generation <= 3; generation++ {
		abandonFuseMountAt(t, publicDir)
		if _, err := os.Stat(publicDir); !isDisconnectedMountStatError(err) {
			t.Fatalf("generation %d: expected disconnected mount before bootstrap, stat err = %v", generation, err)
		}

		// A fresh Runtime over the same session stands in for the exec
		// successor bootstrapping over its predecessor's dead mount.
		rt := NewWithProvider(cfg, &mockProvider{})
		silenceRuntime(rt)
		useRealPublicSurface(rt)
		rt.originalInput = "reclaim stale public surface mount"

		if err := rt.bootstrapAgentRoot(); err != nil {
			t.Fatalf("generation %d: bootstrapAgentRoot over abandoned mount failed: %v", generation, err)
		}
		if reason := rt.publicSurfaceUnavailableReason(); reason != "" {
			t.Fatalf("generation %d: stale mount was misread as degradation: %s", generation, reason)
		}
		sessionJSON := filepath.Join(publicDir, "status", "session.json")
		if data, err := os.ReadFile(sessionJSON); err != nil {
			t.Fatalf("generation %d: successor public status unreadable: %v", generation, err)
		} else if !strings.Contains(string(data), "session") {
			t.Fatalf("generation %d: unexpected session.json content: %q", generation, string(data))
		}
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("generation %d: cleanupAgentRoot failed: %v", generation, err)
		}
	}
}
