package runtime

// runtimeSurfaceBackendName is the single, unconditional public-surface backend.
// The control-write substrate was collapsed to FUSE-only; see
// Paper/core/decisions/2026-06-control-surface-fuse-only.md.
const runtimeSurfaceBackendName = "fuse"

// publicSurfacePaths describes the retained targets the public FUSE surface
// projects: the status SST directory, the control log, and (optionally) the
// self-source tree. PublicDir is the mountpoint.
type publicSurfacePaths struct {
	PublicDir        string
	StatusTarget     string
	ControlLogTarget string
	SourceRoot       string
}

func (r *Runtime) syncPublicSurface(paths publicSurfacePaths) error {
	if r.publicSurfaceSyncFn != nil {
		return r.publicSurfaceSyncFn(paths)
	}
	if r.publicSurface == nil {
		surface, err := newFUSEPublicSurfaceBackend(r)
		if err != nil {
			return err
		}
		r.publicSurface = surface
	}
	return r.publicSurface.sync(paths)
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
