//go:build !linux

package runtime

import "errors"

// On non-Linux hosts the public runtime surface — and therefore the entire
// peer-control-write path (ctl/{post,poke,inject,interrupt}) — is unavailable.
// Control-write requires Linux + /dev/fuse + fusermount; there is no fallback.
// See Paper/core/decisions/2026-06-control-surface-fuse-only.md.
const runtimeSurfaceUnsupportedMsg = "runtime public surface requires Linux + /dev/fuse: peer control-write is unavailable on non-Linux hosts"

type fusePublicSurfaceBackend struct{}

func (b *fusePublicSurfaceBackend) sync(publicSurfacePaths) error {
	return errors.New(runtimeSurfaceUnsupportedMsg)
}

func (b *fusePublicSurfaceBackend) cleanup() error { return nil }

func preflightRuntimeSurfaceFUSE() error {
	return errors.New(runtimeSurfaceUnsupportedMsg)
}

// reclaimStaleRuntimeSurfaceMount is a no-op off Linux: without FUSE support
// no previous incarnation can have left a mount behind.
func reclaimStaleRuntimeSurfaceMount(string) (string, error) { return "", nil }

func newFUSEPublicSurfaceBackend(*Runtime) (*fusePublicSurfaceBackend, error) {
	return nil, errors.New(runtimeSurfaceUnsupportedMsg)
}
