// Package quine is the module root. In the public (trimmed) build the
// self-source / self-reproduction feature is compiled out: the runtime embeds
// no copy of its own source and no git bundle, so a public binary cannot leak
// repository history. These stubs exist only so internal/runtime compiles; the
// feature stays inert (QUINE_SELF_SOURCE_CODE_ENABLED defaults to false).
package quine

import (
	"fmt"
	"io/fs"
)

// SelfSourceBundle is empty in the trimmed public build.
var SelfSourceBundle []byte

// WalkSelfSource is a no-op in the trimmed public build.
func WalkSelfSource(fn fs.WalkDirFunc) error { return nil }

// ReadSelfSource is unavailable in the trimmed public build.
func ReadSelfSource(path string) ([]byte, error) {
	return nil, fmt.Errorf("self-source projection is not available in this build")
}

// SelfSourceBundleAvailable reports false in the trimmed public build.
func SelfSourceBundleAvailable() bool { return false }

// ReadSelfSourceBundle is unavailable in the trimmed public build.
func ReadSelfSourceBundle() ([]byte, error) {
	return nil, fmt.Errorf("embedded self-source repository bundle is not available in this build")
}
