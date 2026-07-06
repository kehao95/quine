package runtime

import (
	"fmt"
	"os"
	"path/filepath"
)

// runtimeSurfaceBackendName is the single, unconditional public-surface backend.
// The control-write substrate was collapsed to FUSE-only; see
// Paper/core/decisions/2026-06-control-surface-fuse-only.md.
const runtimeSurfaceBackendName = "fuse"

// fuseDevicePath is the FUSE control device the Linux preflight probes.
// Package variable so tests can point it at a nonexistent path to exercise
// degraded mode on hosts that do have FUSE. Unused on non-Linux builds,
// where the preflight is unconditionally unsupported.
var fuseDevicePath = "/dev/fuse"

// publicSurfaceUnavailableMarkerName is the single file materialized inside a
// degraded public/ directory so peers and agents see an honest explanation
// instead of an empty or missing surface.
const publicSurfaceUnavailableMarkerName = "UNAVAILABLE"

// publicSurfacePaths describes the retained targets the public FUSE surface
// projects: the status SST directory, the control log, the config/ read
// surface, and (optionally) the self-source tree. PublicDir is the mountpoint.
type publicSurfacePaths struct {
	PublicDir        string
	StatusTarget     string
	ControlLogTarget string
	ConfigTarget     string
	SourceRoot       string
}

// publicSurfaceUnavailableReason reports why the public FUSE surface cannot be
// served in this environment ("" = available). The preflight probe runs once
// per process and the result is memoized; a later mount failure also lands its
// reason here (see syncPublicSurface), so prompt disclosure and the surface
// sync consume one consistent state instead of re-probing.
func (r *Runtime) publicSurfaceUnavailableReason() string {
	if r.publicSurfaceSyncFn != nil {
		// Test seam installed: the fake backend serves the surface.
		return ""
	}
	if !r.publicSurfaceProbed {
		r.publicSurfaceProbed = true
		if err := preflightRuntimeSurfaceFUSE(); err != nil {
			r.markPublicSurfaceDegraded(err)
		}
	}
	return r.publicSurfaceDegraded
}

// markPublicSurfaceDegraded records the degradation reason and logs the
// warning once. The run continues without public/; only peer-facing
// projection and control-write are lost.
func (r *Runtime) markPublicSurfaceDegraded(cause error) {
	r.publicSurfaceProbed = true
	if r.publicSurfaceDegraded != "" {
		return
	}
	r.publicSurfaceDegraded = cause.Error()
	r.log("public surface degraded: %v; continuing without public/ (peer reads and control-write unavailable)", cause)
	r.logError("public surface degraded: %v; continuing without public/ (peer reads and control-write unavailable)", cause)
}

// reclaimStalePublicSurfaceMount frees publicDir from a mount abandoned by a
// previous process image before bootstrapAgentRoot recreates the agent-root
// tree: after the exec handover the predecessor's in-process FUSE server is
// gone and its mount survives disconnected, so MkdirAll on publicDir — and any
// fresh mount over it — fails with ENOTCONN/EEXIST (diagnosed 2026-06-29 in
// development/autopoiesis-probes/l0-ablation-2x2-2026-06-29). A surface this
// process is actively serving is left alone. Failures are logged, not fatal:
// the mount attempt and its degradation handling decide downstream.
func (r *Runtime) reclaimStalePublicSurfaceMount(publicDir string) {
	if r.publicSurfaceSyncFn != nil || r.publicSurface != nil {
		return
	}
	signal, err := reclaimStaleRuntimeSurfaceMount(publicDir)
	if err != nil {
		r.log("public surface: %v", err)
		return
	}
	if signal != "" {
		r.log("public surface: reclaimed stale mount left by a previous incarnation at %s (%s)", publicDir, signal)
	}
}

func (r *Runtime) syncPublicSurface(paths publicSurfacePaths) error {
	if r.publicSurfaceSyncFn != nil {
		return r.publicSurfaceSyncFn(paths)
	}
	if reason := r.publicSurfaceUnavailableReason(); reason != "" {
		return writePublicSurfaceUnavailableMarker(paths.PublicDir, reason)
	}
	if r.publicSurface == nil {
		surface, err := newFUSEPublicSurfaceBackend(r)
		if err != nil {
			r.markPublicSurfaceDegraded(err)
			return writePublicSurfaceUnavailableMarker(paths.PublicDir, r.publicSurfaceDegraded)
		}
		r.publicSurface = surface
	}
	if err := r.publicSurface.sync(paths); err != nil {
		// The mount attempt itself was denied (container without FUSE
		// privileges, seccomp, missing fusermount fallback, ...). Degrade
		// instead of killing the process: startup must survive hosts where
		// the public surface cannot exist.
		r.publicSurface = nil
		r.markPublicSurfaceDegraded(err)
		return writePublicSurfaceUnavailableMarker(paths.PublicDir, r.publicSurfaceDegraded)
	}
	return nil
}

// writePublicSurfaceUnavailableMarker materializes public/ as a plain
// directory holding a single UNAVAILABLE marker stating why the FUSE
// projection is absent, so peers scanning the surface see an honest
// explanation rather than nothing.
func writePublicSurfaceUnavailableMarker(publicDir string, reason string) error {
	if err := os.MkdirAll(publicDir, 0o755); err != nil {
		return fmt.Errorf("mkdir degraded public surface %s: %w", publicDir, err)
	}
	body := "public runtime surface unavailable in this environment.\n" +
		"reason: " + reason + "\n" +
		"This process runs normally, but peers cannot read public/{status,log,source-code} or write public/ctl here.\n" +
		"The public projection requires Linux with FUSE (/dev/fuse + fusermount); see Paper/core/decisions/2026-06-control-surface-fuse-only.md.\n"
	return writeTextFile(filepath.Join(publicDir, publicSurfaceUnavailableMarkerName), body)
}

func (r *Runtime) cleanupPublicSurface() error {
	if r.publicSurfaceCleanupFn != nil {
		return r.publicSurfaceCleanupFn()
	}
	if r.publicSurface == nil {
		return nil
	}
	err := r.publicSurface.cleanup()
	r.publicSurface = nil
	return err
}
