//go:build linux

package runtime

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	fusefs "github.com/hanwen/go-fuse/v2/fs"
	fusepkg "github.com/hanwen/go-fuse/v2/fuse"
)

type fusePublicSurfaceBackend struct {
	rt         *Runtime
	mu         sync.Mutex
	server     *fusepkg.Server
	mountpoint string
}

func newFUSEPublicSurfaceBackend(r *Runtime) (*fusePublicSurfaceBackend, error) {
	if err := preflightRuntimeSurfaceFUSE(); err != nil {
		return nil, err
	}
	return &fusePublicSurfaceBackend{rt: r}, nil
}

func preflightRuntimeSurfaceFUSE() error {
	f, err := os.OpenFile(fuseDevicePath, os.O_RDWR, 0)
	if err != nil {
		// err is a *PathError that already names the device path.
		return fmt.Errorf("runtime surface FUSE unsupported in this Linux environment: %w", err)
	}
	_ = f.Close()
	// Root can mount(2) directly (DirectMount); everyone else needs the
	// fusermount/fusermount3 setuid helper as the go-fuse fallback path.
	if os.Geteuid() != 0 {
		if _, err3 := exec.LookPath("fusermount3"); err3 != nil {
			if _, err1 := exec.LookPath("fusermount"); err1 != nil {
				return fmt.Errorf("runtime surface FUSE unsupported in this Linux environment: neither fusermount3 nor fusermount found in PATH")
			}
		}
	}
	return nil
}

// detectStaleRuntimeSurfaceMount reports the signal ("" = clean) that the
// mountpoint is occupied by a mount this process is not serving. After the
// exec handover the predecessor image's in-process FUSE server is gone, so its
// mount survives disconnected: stat returns ENOTCONN and /proc/self/mountinfo
// still lists the path (diagnosed 2026-06-29 in development/autopoiesis-probes/
// l0-ablation-2x2-2026-06-29/agent-d2-stdout.log). Callers reclaim before
// mounting so a stale mount is never misread as "mount denied → degrade".
func detectStaleRuntimeSurfaceMount(mountpoint string) string {
	if _, err := os.Stat(mountpoint); isDisconnectedMountStatError(err) {
		return err.Error()
	}
	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	if fstype, ok := mountinfoEntryForMountpoint(string(data), mountpoint); ok {
		return fmt.Sprintf("%s mount listed in /proc/self/mountinfo", fstype)
	}
	return ""
}

// isDisconnectedMountStatError classifies a stat error as the
// disconnected-FUSE signal: ENOTCONN means the mountpoint's filesystem lost
// its userspace server.
func isDisconnectedMountStatError(err error) bool {
	return err != nil && errors.Is(err, syscall.ENOTCONN)
}

// mountinfoEntryForMountpoint scans /proc/self/mountinfo content for an entry
// whose mount point (field 5, octal-escaped per proc(5)) equals mountpoint and
// returns its filesystem type (the field after the "-" separator).
func mountinfoEntryForMountpoint(mountinfo, mountpoint string) (string, bool) {
	target := filepath.Clean(mountpoint)
	for _, line := range strings.Split(mountinfo, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 || unescapeMountinfoField(fields[4]) != target {
			continue
		}
		fstype := "unknown"
		for i := 5; i < len(fields)-1; i++ {
			if fields[i] == "-" {
				fstype = fields[i+1]
				break
			}
		}
		return fstype, true
	}
	return "", false
}

// unescapeMountinfoField decodes the \ooo octal escapes proc(5) applies to
// space, tab, newline, and backslash in mountinfo paths.
func unescapeMountinfoField(field string) string {
	if !strings.Contains(field, `\`) {
		return field
	}
	var out strings.Builder
	for i := 0; i < len(field); i++ {
		if field[i] == '\\' && i+3 < len(field) &&
			field[i+1] >= '0' && field[i+1] <= '7' &&
			field[i+2] >= '0' && field[i+2] <= '7' &&
			field[i+3] >= '0' && field[i+3] <= '7' {
			out.WriteByte((field[i+1]-'0')<<6 | (field[i+2]-'0')<<3 | (field[i+3] - '0'))
			i += 3
			continue
		}
		out.WriteByte(field[i])
	}
	return out.String()
}

// reclaimStaleRuntimeSurfaceMount unmounts whatever occupies mountpoint so a
// fresh mount can land. Unmount chain: fusermount3/fusermount -u, then lazy
// -uz (the setuid helper is the only road for non-root mounts), then direct
// umount2 — plain, then MNT_DETACH — for the root/direct-mount case the
// preflight exempts from the fusermount requirement. Returns the detection
// signal ("" when nothing was stale) and an error when the mountpoint could
// not be freed.
func reclaimStaleRuntimeSurfaceMount(mountpoint string) (string, error) {
	signal := detectStaleRuntimeSurfaceMount(mountpoint)
	if signal == "" {
		return "", nil
	}
	var attempts []string
	unmounted := func(err error, attempt string) bool {
		if err == nil {
			if detectStaleRuntimeSurfaceMount(mountpoint) == "" {
				return true
			}
			err = errors.New("mountpoint still occupied")
		}
		attempts = append(attempts, fmt.Sprintf("%s: %v", attempt, err))
		return false
	}
	helper := ""
	for _, name := range []string{"fusermount3", "fusermount"} {
		if path, err := exec.LookPath(name); err == nil {
			helper = path
			break
		}
	}
	if helper == "" {
		attempts = append(attempts, "fusermount3/fusermount: not found in PATH")
	} else {
		for _, flag := range []string{"-u", "-uz"} {
			out, err := exec.Command(helper, flag, mountpoint).CombinedOutput()
			if err != nil && len(out) > 0 {
				err = fmt.Errorf("%w (%s)", err, strings.TrimSpace(string(out)))
			}
			if unmounted(err, filepath.Base(helper)+" "+flag) {
				return signal, nil
			}
		}
	}
	for _, flags := range []int{0, syscall.MNT_DETACH} {
		if unmounted(syscall.Unmount(mountpoint, flags), fmt.Sprintf("umount2 flags=%#x", flags)) {
			return signal, nil
		}
	}
	return signal, fmt.Errorf("stale mount at %s not reclaimed (%s): %s", mountpoint, signal, strings.Join(attempts, "; "))
}

func (b *fusePublicSurfaceBackend) sync(paths publicSurfacePaths) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.server != nil {
		if b.mountpoint != paths.PublicDir {
			return fmt.Errorf("runtime surface FUSE already mounted at %s", b.mountpoint)
		}
		return nil
	}

	// A predecessor image's mount survives the exec handover as a disconnected
	// mountpoint that refuses a fresh mount over it. Reclaim it first; only if
	// the mount below still fails does the caller conclude degradation.
	if signal, err := reclaimStaleRuntimeSurfaceMount(paths.PublicDir); err != nil {
		b.rt.log("public surface: %v", err)
	} else if signal != "" {
		b.rt.log("public surface: reclaimed stale mount at %s (%s)", paths.PublicDir, signal)
	}

	timeout := time.Duration(0)
	server, err := fusefs.Mount(paths.PublicDir, &fusePublicSurfaceRoot{
		rt:    b.rt,
		paths: paths,
	}, &fusefs.Options{
		EntryTimeout: &timeout,
		AttrTimeout:  &timeout,
		MountOptions: fusepkg.MountOptions{
			FsName:        "quine-public",
			Name:          "quine-public",
			DisableXAttrs: true,
			DirectMount:   true,
		},
	})
	if err != nil {
		return fmt.Errorf("mount runtime public surface: %w", err)
	}
	b.server = server
	b.mountpoint = paths.PublicDir
	return nil
}

func (b *fusePublicSurfaceBackend) cleanup() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.server == nil {
		return nil
	}
	err := b.server.Unmount()
	b.server = nil
	b.mountpoint = ""
	if err != nil {
		return fmt.Errorf("unmount runtime public surface: %w", err)
	}
	return nil
}

type fusePublicSurfaceRoot struct {
	fusefs.Inode
	rt    *Runtime
	paths publicSurfacePaths
}

var _ = (fusefs.NodeOnAdder)((*fusePublicSurfaceRoot)(nil))

func (r *fusePublicSurfaceRoot) OnAdd(ctx context.Context) {
	r.AddChild("ctl", r.NewPersistentInode(ctx, &fuseControlDirNode{rt: r.rt}, fusefs.StableAttr{
		Mode: syscall.S_IFDIR,
	}), false)
	r.AddChild("status", r.NewPersistentInode(ctx, &fuseStatusDirNode{
		statusDirTarget: r.paths.StatusTarget,
	}, fusefs.StableAttr{
		Mode: syscall.S_IFDIR,
	}), false)
	r.AddChild("log", r.NewPersistentInode(ctx, &fuseLogDirNode{
		controlLogTarget: r.paths.ControlLogTarget,
	}, fusefs.StableAttr{
		Mode: syscall.S_IFDIR,
	}), false)
	if r.paths.ConfigTarget != "" {
		r.AddChild("config", r.NewPersistentInode(ctx, &fuseConfigDirNode{
			configDirTarget: r.paths.ConfigTarget,
		}, fusefs.StableAttr{
			Mode: syscall.S_IFDIR,
		}), false)
	}
	if r.paths.SourceRoot != "" {
		r.AddChild("source-code", r.NewPersistentInode(ctx, &fuseProjectedDirNode{
			targetDir: r.paths.SourceRoot,
		}, fusefs.StableAttr{
			Mode: syscall.S_IFDIR,
		}), false)
	}
}

type fuseLogDirNode struct {
	fusefs.Inode
	controlLogTarget string
}

var _ = (fusefs.NodeOnAdder)((*fuseLogDirNode)(nil))

func (n *fuseLogDirNode) OnAdd(ctx context.Context) {
	n.AddChild("control.jsonl", n.NewPersistentInode(ctx, &fuseProjectedFileNode{
		targetPath: n.controlLogTarget,
	}, fusefs.StableAttr{
		Mode: syscall.S_IFREG,
	}), false)
}

type fuseStatusDirNode struct {
	fusefs.Inode
	statusDirTarget string
}

var _ = (fusefs.NodeOnAdder)((*fuseStatusDirNode)(nil))

func (n *fuseStatusDirNode) OnAdd(ctx context.Context) {
	for _, name := range []string{"session.json", "inbox.json", "contract.json"} {
		n.AddChild(name, n.NewPersistentInode(ctx, &fuseProjectedFileNode{
			targetPath: filepath.Join(n.statusDirTarget, name),
		}, fusefs.StableAttr{
			Mode: syscall.S_IFREG,
		}), false)
	}
}

// fuseConfigDirNode projects the peer-readable slice of the agent-root
// config/ read surface (registry-design-brief § B: the peer-readable
// capability position) as live views over the backing SSTs — the status/*
// computed-on-read vehicle: every read fetches the current backing file, so a
// peer reading resolved.env after onConfigMutated sees the mutated position.
// The child set is a fixed enumeration: the Phase-3 agent-writable
// config/next.env staging file must never become peer-visible here.
type fuseConfigDirNode struct {
	fusefs.Inode
	configDirTarget string
}

var _ = (fusefs.NodeOnAdder)((*fuseConfigDirNode)(nil))

func (n *fuseConfigDirNode) OnAdd(ctx context.Context) {
	for _, name := range []string{"registry.json", "resolved.env"} {
		n.AddChild(name, n.NewPersistentInode(ctx, &fuseProjectedFileNode{
			targetPath: filepath.Join(n.configDirTarget, name),
		}, fusefs.StableAttr{
			Mode: syscall.S_IFREG,
		}), false)
	}
}

type fuseProjectedFileNode struct {
	fusefs.Inode
	targetPath string
}

type fuseProjectedDirNode struct {
	fusefs.Inode
	targetDir string
}

var _ = (fusefs.NodeOnAdder)((*fuseProjectedDirNode)(nil))
var _ = (fusefs.NodeGetattrer)((*fuseProjectedDirNode)(nil))

func (n *fuseProjectedDirNode) OnAdd(ctx context.Context) {
	entries, err := os.ReadDir(n.targetDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		name := entry.Name()
		target := filepath.Join(n.targetDir, name)
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() {
			n.AddChild(name, n.NewPersistentInode(ctx, &fuseProjectedDirNode{
				targetDir: target,
			}, fusefs.StableAttr{
				Mode: syscall.S_IFDIR,
			}), false)
			continue
		}
		if info.Mode().IsRegular() {
			n.AddChild(name, n.NewPersistentInode(ctx, &fuseProjectedFileNode{
				targetPath: target,
			}, fusefs.StableAttr{
				Mode: syscall.S_IFREG,
			}), false)
		}
	}
}

func (n *fuseProjectedDirNode) Getattr(ctx context.Context, fh fusefs.FileHandle, out *fusepkg.AttrOut) syscall.Errno {
	info, err := os.Stat(n.targetDir)
	if err != nil {
		return errnoFromError(err)
	}
	if !info.IsDir() {
		return syscall.ENOTDIR
	}
	out.Mode = syscall.S_IFDIR | 0o555
	return 0
}

var _ = (fusefs.NodeGetattrer)((*fuseProjectedFileNode)(nil))
var _ = (fusefs.NodeOpener)((*fuseProjectedFileNode)(nil))

func (n *fuseProjectedFileNode) Open(ctx context.Context, flags uint32) (fusefs.FileHandle, uint32, syscall.Errno) {
	if int(flags)&syscall.O_ACCMODE != syscall.O_RDONLY {
		return nil, 0, syscall.EACCES
	}
	return &fuseProjectedFileHandle{targetPath: n.targetPath}, fusepkg.FOPEN_DIRECT_IO, 0
}

func (n *fuseProjectedFileNode) Getattr(ctx context.Context, fh fusefs.FileHandle, out *fusepkg.AttrOut) syscall.Errno {
	info, err := os.Stat(n.targetPath)
	if err != nil {
		return errnoFromError(err)
	}
	out.Mode = syscall.S_IFREG | 0o444
	out.Size = uint64(info.Size())
	return 0
}

type fuseProjectedFileHandle struct {
	targetPath string
}

var _ = (fusefs.FileReader)((*fuseProjectedFileHandle)(nil))

func (h *fuseProjectedFileHandle) Read(ctx context.Context, dest []byte, off int64) (fusepkg.ReadResult, syscall.Errno) {
	data, err := os.ReadFile(h.targetPath)
	if err != nil {
		return nil, errnoFromError(err)
	}
	if off >= int64(len(data)) {
		return fusepkg.ReadResultData(nil), 0
	}
	end := off + int64(len(dest))
	if end > int64(len(data)) {
		end = int64(len(data))
	}
	return fusepkg.ReadResultData(data[off:end]), 0
}

func errnoFromError(err error) syscall.Errno {
	if err == nil {
		return 0
	}
	if errors.Is(err, fs.ErrNotExist) {
		return syscall.ENOENT
	}
	if errors.Is(err, fs.ErrPermission) {
		return syscall.EACCES
	}
	return syscall.EIO
}

type fuseControlDirNode struct {
	fusefs.Inode
	rt *Runtime
}

var _ = (fusefs.NodeOnAdder)((*fuseControlDirNode)(nil))

func (n *fuseControlDirNode) OnAdd(ctx context.Context) {
	for _, action := range controlSurfaceActions {
		n.AddChild(string(action), n.NewPersistentInode(ctx, &fuseControlActionNode{
			rt:     n.rt,
			action: action,
		}, fusefs.StableAttr{
			Mode: syscall.S_IFREG,
		}), false)
	}
}

type fuseControlActionNode struct {
	fusefs.Inode
	rt                *Runtime
	action            controlSurfaceAction
	mu                sync.Mutex
	activeWriteHandle *fuseControlActionHandle
}

var _ = (fusefs.NodeOpener)((*fuseControlActionNode)(nil))
var _ = (fusefs.NodeGetattrer)((*fuseControlActionNode)(nil))
var _ = (fusefs.NodeSetattrer)((*fuseControlActionNode)(nil))

func (n *fuseControlActionNode) Open(ctx context.Context, flags uint32) (fusefs.FileHandle, uint32, syscall.Errno) {
	writeAccess := int(flags)&syscall.O_ACCMODE != syscall.O_RDONLY
	handle := &fuseControlActionHandle{node: n, writeAccess: writeAccess}
	if writeAccess {
		n.mu.Lock()
		n.activeWriteHandle = handle
		n.mu.Unlock()
	}
	if int(flags)&syscall.O_TRUNC != 0 {
		handle.reset(true)
	}
	return handle, fusepkg.FOPEN_DIRECT_IO, 0
}

func (n *fuseControlActionNode) Getattr(ctx context.Context, fh fusefs.FileHandle, out *fusepkg.AttrOut) syscall.Errno {
	data := n.rt.publicControlSurfaceSummary(n.action)
	out.Mode = syscall.S_IFREG | 0o666
	out.Size = uint64(len(data))
	return 0
}

func (n *fuseControlActionNode) Setattr(ctx context.Context, fh fusefs.FileHandle, in *fusepkg.SetAttrIn, out *fusepkg.AttrOut) syscall.Errno {
	if sz, ok := in.GetSize(); ok && sz == 0 {
		if handle, ok := fh.(*fuseControlActionHandle); ok {
			handle.reset(true)
		} else {
			n.mu.Lock()
			handle := n.activeWriteHandle
			n.mu.Unlock()
			if handle != nil {
				handle.reset(true)
			}
		}
	}
	return n.Getattr(ctx, fh, out)
}

func (n *fuseControlActionNode) clearActiveWriteHandle(handle *fuseControlActionHandle) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.activeWriteHandle == handle {
		n.activeWriteHandle = nil
	}
}

type fuseControlActionHandle struct {
	node        *fuseControlActionNode
	writeAccess bool
	mu          sync.Mutex
	buf         []byte
	dirty       bool
	committed   bool
}

var _ = (fusefs.FileReader)((*fuseControlActionHandle)(nil))
var _ = (fusefs.FileWriter)((*fuseControlActionHandle)(nil))
var _ = (fusefs.FileFlusher)((*fuseControlActionHandle)(nil))
var _ = (fusefs.FileReleaser)((*fuseControlActionHandle)(nil))

func (h *fuseControlActionHandle) Read(ctx context.Context, dest []byte, off int64) (fusepkg.ReadResult, syscall.Errno) {
	data := h.node.rt.publicControlSurfaceSummary(h.node.action)
	if off >= int64(len(data)) {
		return fusepkg.ReadResultData(nil), 0
	}
	end := off + int64(len(dest))
	if end > int64(len(data)) {
		end = int64(len(data))
	}
	return fusepkg.ReadResultData(data[off:end]), 0
}

func (h *fuseControlActionHandle) Write(ctx context.Context, data []byte, off int64) (uint32, syscall.Errno) {
	if off < 0 {
		return 0, syscall.EINVAL
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	end := int(off) + len(data)
	if end > len(h.buf) {
		next := make([]byte, end)
		copy(next, h.buf)
		h.buf = next
	}
	copy(h.buf[int(off):end], data)
	h.dirty = true
	return uint32(len(data)), 0
}

func (h *fuseControlActionHandle) Flush(ctx context.Context) syscall.Errno {
	// The config gate is a validated transaction: commit at flush so the
	// writer's close(2) observes acceptance or rejection synchronously
	// (RELEASE carries no reply back to any syscall, so a release-time
	// rejection would be invisible to the writer). The message actions keep
	// their commit at release: their payloads cannot be rejected, so nothing
	// is lost by the asynchronous ack.
	if h.node.action == controlActionConfig {
		return h.commit()
	}
	return 0
}

func (h *fuseControlActionHandle) Release(ctx context.Context) syscall.Errno {
	if h.writeAccess {
		h.node.clearActiveWriteHandle(h)
	}
	return h.commit()
}

func (h *fuseControlActionHandle) commit() syscall.Errno {
	h.mu.Lock()
	payload := string(h.buf)
	dirty := h.dirty
	committed := h.committed
	if dirty && !committed {
		h.committed = true
	}
	h.mu.Unlock()
	if !dirty || committed {
		return 0
	}
	if err := h.node.rt.applyControlSurfaceAction(h.node.action, payload); err != nil {
		h.node.rt.log("public ctl %s error: %v", h.node.action, err)
		return syscall.EINVAL
	}
	return 0
}

func (h *fuseControlActionHandle) reset(dirty bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.buf = nil
	h.dirty = dirty
}
