package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	stdio "io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kehao95/quine/internal/config"
)

const overlayStateVersion = "overlay-lineage-v1"
const directStateVersion = "workspace-direct-v1"

type subjectiveFS struct {
	enabled bool

	sessionID        string
	workspaceRoot    string
	workspace        string
	workspaceBackend string
	overlayDriver    string
	revisionMode     config.WorkspaceRevisionMode
	dataDir          string
	workspaceSession string
	workspaceOwner   bool
	bootstrapSource  string
	currentRevision  string
	overlayMountOpts string

	initialized bool
}

type worldRevision struct {
	ID             string `json:"id"`
	Parent         string `json:"parent,omitempty"`
	Kind           string `json:"kind"`
	Turn           int    `json:"turn,omitempty"`
	SourceSession  string `json:"source_session,omitempty"`
	SourceRevision string `json:"source_revision,omitempty"`
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

type turnFinalizeResult struct {
	Changed   bool
	Mutations string
	Revision  worldRevision
}

func newSubjectiveFS(cfg *config.Config) *subjectiveFS {
	if cfg == nil || !cfg.WorkspaceEnabled {
		return &subjectiveFS{}
	}
	return &subjectiveFS{
		enabled:          true,
		sessionID:        cfg.SessionID,
		workspaceRoot:    cfg.WorkspaceRoot,
		workspace:        cfg.Workspace,
		workspaceBackend: cfg.WorkspaceBackend,
		overlayDriver:    cfg.WorkspaceOverlayDriver,
		revisionMode:     cfg.WorkspaceRevisionMode,
		dataDir:          cfg.DataDir,
		workspaceSession: cfg.WorkspaceSession,
		workspaceOwner:   cfg.WorkspaceOwner,
		bootstrapSource:  cfg.WorkspaceBootstrap,
	}
}

func (s *subjectiveFS) effectiveWorkspaceBackend() string {
	if strings.TrimSpace(s.workspaceBackend) == "" {
		return "overlay"
	}
	return s.workspaceBackend
}

func (s *subjectiveFS) usesOverlayBackend() bool {
	return s.effectiveWorkspaceBackend() == "overlay"
}

func (s *subjectiveFS) usesDirectBackend() bool {
	return s.effectiveWorkspaceBackend() == "direct"
}

func (s *subjectiveFS) canTrackWorldRevisions() bool {
	return s.enabled && s.usesOverlayBackend() && s.revisionMode != config.WorkspaceRevisionNone
}

func (s *subjectiveFS) canRestoreWorld() bool {
	return s.enabled && s.usesOverlayBackend() && s.revisionMode == config.WorkspaceRevisionRestore
}

func (s *subjectiveFS) currentWorldRevision() string {
	if s.canTrackWorldRevisions() && s.initialized {
		if revision, err := s.loadCurrentWorldRevision(); err == nil && revision.ID != "" {
			s.currentRevision = revision.ID
			return revision.ID
		}
	}
	if s.currentRevision != "" {
		return s.currentRevision
	}
	if s.canTrackWorldRevisions() {
		return "wr0"
	}
	return ""
}

func (s *subjectiveFS) readWorkspaceFile(rawPath string) ([]byte, error) {
	if s == nil || !s.enabled {
		return nil, fs.ErrNotExist
	}
	if !s.initialized {
		if err := s.init(s.dataDir, s.sessionID); err != nil {
			return nil, err
		}
	}

	rel, ok, err := s.workspaceRelativePath(rawPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fs.ErrNotExist
	}

	if s.usesDirectBackend() {
		return os.ReadFile(filepath.Join(s.workspaceRoot, rel))
	}
	if s.usesOverlayBackend() {
		return s.readOverlayWorkspaceFile(rel)
	}
	return nil, fs.ErrNotExist
}

func (s *subjectiveFS) workspaceRelativePath(rawPath string) (string, bool, error) {
	rawPath = strings.TrimSpace(rawPath)
	if rawPath == "" {
		return "", false, fs.ErrInvalid
	}

	root := strings.TrimSpace(s.workspaceRoot)
	if root == "" {
		return "", false, nil
	}
	root = filepath.Clean(root)
	workspace := strings.TrimSpace(s.workspace)
	if workspace == "" {
		workspace = root
	}
	workspace = filepath.Clean(workspace)

	relativeFromRoot := func(absPath string) (string, bool, error) {
		absPath = filepath.Clean(absPath)
		if !isPathWithin(root, absPath) {
			return "", false, nil
		}
		rel, err := filepath.Rel(root, absPath)
		if err != nil {
			return "", false, err
		}
		rel = filepath.Clean(rel)
		if rel == "." {
			return "", false, fs.ErrInvalid
		}
		return rel, true, nil
	}

	if filepath.IsAbs(rawPath) {
		if rel, ok, err := relativeFromRoot(rawPath); ok || err != nil {
			return rel, ok, err
		}
		if cwd, err := os.Getwd(); err == nil && isPathWithin(filepath.Clean(cwd), filepath.Clean(rawPath)) {
			cwdRel, relErr := filepath.Rel(filepath.Clean(cwd), filepath.Clean(rawPath))
			if relErr != nil {
				return "", false, relErr
			}
			return relativeFromRoot(filepath.Join(workspace, cwdRel))
		}
		return "", false, nil
	}

	return relativeFromRoot(filepath.Join(workspace, rawPath))
}

func (s *subjectiveFS) observedDir() string {
	return filepath.Join(s.directWorkspaceStateDir(), "observed")
}

func (s *subjectiveFS) observedTreePath() string {
	observer := strings.TrimSpace(s.sessionID)
	if observer == "" {
		observer = strings.TrimSpace(s.workspaceSession)
	}
	if observer == "" {
		observer = "observer"
	}
	return filepath.Join(s.observedDir(), observer+".json")
}

func (s *subjectiveFS) directWorkspaceStateDir() string {
	return filepath.Join(s.dataDir, "workspaces", s.workspaceSession)
}

func (s *subjectiveFS) directStateVersionPath() string {
	return filepath.Join(s.directWorkspaceStateDir(), "STATE_VERSION")
}

func (s *subjectiveFS) initDirectState() error {
	if err := os.MkdirAll(s.directWorkspaceStateDir(), 0o755); err != nil {
		return fmt.Errorf("create direct workspace state dir: %w", err)
	}
	if data, err := os.ReadFile(s.directStateVersionPath()); err == nil {
		if strings.TrimSpace(string(data)) != directStateVersion {
			return fmt.Errorf("workspace session %q uses unsupported direct state version %q", s.workspaceSession, strings.TrimSpace(string(data)))
		}
	} else if errorsIsNotExist(err) {
		if err := os.WriteFile(s.directStateVersionPath(), []byte(directStateVersion+"\n"), 0o644); err != nil {
			return fmt.Errorf("write direct workspace state version: %w", err)
		}
	} else {
		return fmt.Errorf("stat direct workspace state version: %w", err)
	}
	if err := os.MkdirAll(s.observedDir(), 0o755); err != nil {
		return fmt.Errorf("create direct observed dir: %w", err)
	}
	if _, err := os.Lstat(s.observedTreePath()); err == nil {
		return nil
	} else if !errorsIsNotExist(err) {
		return fmt.Errorf("stat direct observed tree: %w", err)
	}
	tree, err := snapshotTree(s.workspace)
	if err != nil {
		return fmt.Errorf("snapshot direct workspace: %w", err)
	}
	if err := s.saveObservedTree(tree); err != nil {
		return err
	}
	return nil
}

func (s *subjectiveFS) loadObservedTree() (map[string]treeEntry, error) {
	data, err := os.ReadFile(s.observedTreePath())
	if err != nil {
		if errorsIsNotExist(err) {
			return map[string]treeEntry{}, nil
		}
		return nil, fmt.Errorf("read observed tree: %w", err)
	}
	var tree map[string]treeEntry
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, fmt.Errorf("parse observed tree: %w", err)
	}
	if tree == nil {
		tree = map[string]treeEntry{}
	}
	return tree, nil
}

func (s *subjectiveFS) saveObservedTree(tree map[string]treeEntry) error {
	if tree == nil {
		tree = map[string]treeEntry{}
	}
	data, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("marshal observed tree: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.observedTreePath()), 0o755); err != nil {
		return fmt.Errorf("create observed tree parent: %w", err)
	}
	if err := os.WriteFile(s.observedTreePath(), data, 0o644); err != nil {
		return fmt.Errorf("write observed tree: %w", err)
	}
	return nil
}

func (s *subjectiveFS) finalizeDirectTurn() (turnFinalizeResult, error) {
	previous, err := s.loadObservedTree()
	if err != nil {
		return turnFinalizeResult{}, err
	}
	current, err := snapshotTree(s.workspace)
	if err != nil {
		return turnFinalizeResult{}, fmt.Errorf("snapshot direct workspace: %w", err)
	}
	if err := s.saveObservedTree(current); err != nil {
		return turnFinalizeResult{}, err
	}
	mutations := diffTree(previous, current)
	return turnFinalizeResult{
		Changed:   len(mutations) > 0,
		Mutations: formatMutations(mutations),
	}, nil
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
	buf := make([]byte, 32*1024)
	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			if _, err := h.Write(buf[:n]); err != nil {
				return "", err
			}
		}
		if readErr == nil {
			continue
		}
		if readErr == stdio.EOF {
			break
		}
		return "", readErr
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
	if _, err := stdio.Copy(out, in); err != nil {
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
