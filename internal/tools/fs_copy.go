package tools

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kehao95/quine/internal/config"
)

// ContextBootstrapEnv points a child process at a staged context tree that it
// should import into its live agent root before entering the main turn loop.
//
// The name is owned by the capability registry (a runtime-emitted knob) so the
// env mask derives from Mutability alone; this is the in-package spelling.
const ContextBootstrapEnv = config.EnvContextBootstrap

// CopyTreePreservingSymlinks recursively copies a tree, preserving directory
// structure, file modes, and symlink targets.
func CopyTreePreservingSymlinks(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}
	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == src {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		switch mode := info.Mode(); {
		case mode&os.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(target, targetPath)
		case d.IsDir():
			return os.MkdirAll(targetPath, info.Mode().Perm())
		default:
			return copyFileWithMode(path, targetPath, info.Mode().Perm())
		}
	})
}

func copyFileWithMode(src, dst string, mode fs.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = out.Close()
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	closed = true
	return nil
}
