package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kehao95/quine/internal/config"
)

type subjectiveFS struct {
	enabled bool

	workspaceRoot    string
	workspace        string
	workspaceBackend string
	revisionMode     config.WorkspaceRevisionMode
	dataDir          string
	workspaceSession string
	workspaceOwner   bool
	currentRevision  string

	initialized bool
}

type fsSnapshot struct {
	workspace map[string]treeEntry
}

type worldRevision struct {
	ID     string `json:"id"`
	Parent string `json:"parent,omitempty"`
	Kind   string `json:"kind"`
	Turn   int    `json:"turn,omitempty"`
}

type worldRevisionLedger struct {
	Current   string                   `json:"current"`
	Next      int                      `json:"next"`
	Revisions map[string]worldRevision `json:"revisions"`
}

type treeEntry struct {
	Kind string
	Mode fs.FileMode
	Size int64
	Hash string
	Link string
}

type fsMutation struct {
	Path string
	Kind string
}

func newSubjectiveFS(cfg *config.Config) *subjectiveFS {
	if cfg == nil || !cfg.WorkspaceEnabled {
		return &subjectiveFS{}
	}
	return &subjectiveFS{
		enabled:          true,
		workspaceRoot:    cfg.WorkspaceRoot,
		workspace:        cfg.Workspace,
		workspaceBackend: cfg.WorkspaceBackend,
		revisionMode:     cfg.WorkspaceRevisionMode,
		dataDir:          cfg.DataDir,
		workspaceSession: cfg.WorkspaceSession,
		workspaceOwner:   cfg.WorkspaceOwner,
	}
}

func (s *subjectiveFS) directBackend() bool {
	return s.workspaceBackend == "direct"
}

func (s *subjectiveFS) canTrackWorldRevisions() bool {
	return s.enabled && s.revisionMode != config.WorkspaceRevisionNone
}

func (s *subjectiveFS) canRestoreWorld() bool {
	return s.enabled && s.revisionMode == config.WorkspaceRevisionRestore
}

func (s *subjectiveFS) currentWorldRevision() string {
	if s.currentRevision != "" {
		return s.currentRevision
	}
	if s.canTrackWorldRevisions() {
		return "wr0"
	}
	return ""
}

func diffTree(before, after map[string]treeEntry) []fsMutation {
	if len(before) == 0 && len(after) == 0 {
		return nil
	}
	keys := make(map[string]struct{}, len(before)+len(after))
	for path := range before {
		keys[path] = struct{}{}
	}
	for path := range after {
		keys[path] = struct{}{}
	}
	out := make([]fsMutation, 0, len(keys))
	for path := range keys {
		oldEntry, hadOld := before[path]
		newEntry, hadNew := after[path]
		switch {
		case !hadOld && hadNew:
			out = append(out, fsMutation{Path: path, Kind: "created"})
		case hadOld && !hadNew:
			out = append(out, fsMutation{Path: path, Kind: "deleted"})
		case hadOld && hadNew && oldEntry != newEntry:
			out = append(out, fsMutation{Path: path, Kind: "modified"})
		}
	}
	return out
}

func formatMutations(mutations []fsMutation) string {
	if len(mutations) == 0 {
		return "[FS MUTATIONS]\n(empty)"
	}
	sort.Slice(mutations, func(i, j int) bool {
		return mutations[i].Path < mutations[j].Path
	})
	lines := make([]string, 0, len(mutations)+1)
	lines = append(lines, "[FS MUTATIONS]")
	for _, mutation := range mutations {
		lines = append(lines, mutationPrefix(mutation.Kind)+" "+filepath.ToSlash(mutation.Path)+" ("+mutation.Kind+")")
	}
	return strings.Join(lines, "\n")
}

func mutationPrefix(kind string) string {
	switch kind {
	case "created":
		return "+"
	case "deleted":
		return "-"
	default:
		return "~"
	}
}

func snapshotTree(root string) (map[string]treeEntry, error) {
	entries := make(map[string]treeEntry)
	if root == "" {
		return entries, nil
	}
	if _, err := os.Lstat(root); err != nil {
		if errorsIsNotExist(err) {
			return entries, nil
		}
		return nil, err
	}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		entry, entryErr := buildTreeEntry(path, d)
		if entryErr != nil {
			return entryErr
		}
		entries[rel] = entry
		return nil
	})
	return entries, err
}

func buildTreeEntry(path string, d fs.DirEntry) (treeEntry, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return treeEntry{}, err
	}
	entry := treeEntry{
		Mode: info.Mode(),
		Size: info.Size(),
	}
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		entry.Kind = "symlink"
		link, err := os.Readlink(path)
		if err != nil {
			return treeEntry{}, err
		}
		entry.Link = link
	case info.IsDir():
		entry.Kind = "dir"
	case info.Mode().IsRegular():
		entry.Kind = "file"
		hash, err := hashFile(path)
		if err != nil {
			return treeEntry{}, err
		}
		entry.Hash = hash
	default:
		entry.Kind = "other"
	}
	return entry, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func syncTree(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	if err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}

		target := filepath.Join(dst, rel)
		info, infoErr := os.Lstat(path)
		if infoErr != nil {
			return infoErr
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			if err := os.RemoveAll(target); err != nil && !errorsIsNotExist(err) {
				return err
			}
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			return os.Symlink(link, target)
		case info.IsDir():
			return os.MkdirAll(target, info.Mode().Perm())
		case info.Mode().IsRegular():
			return copyFile(path, target, info.Mode().Perm())
		default:
			return nil
		}
	}); err != nil {
		return err
	}

	return filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(dst, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		sourcePath := filepath.Join(src, rel)
		if _, statErr := os.Lstat(sourcePath); statErr == nil {
			return nil
		}
		if err := os.RemoveAll(path); err != nil && !errorsIsNotExist(err) {
			return err
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}

func realishPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return real, nil
	}
	if errorsIsNotExist(err) {
		parent := filepath.Dir(abs)
		parentReal, parentErr := filepath.EvalSymlinks(parent)
		if parentErr == nil {
			return filepath.Join(parentReal, filepath.Base(abs)), nil
		}
	}
	return abs, nil
}

func isPathWithin(root, child string) bool {
	if root == "" || child == "" {
		return false
	}
	if root == child {
		return true
	}
	rel, err := filepath.Rel(root, child)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func errorsIsNotExist(err error) bool {
	return err != nil && (os.IsNotExist(err) || errors.Is(err, fs.ErrNotExist))
}
