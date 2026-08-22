package runtime

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	quineroot "github.com/kehao95/quine"
)

var selfSourceSentinelFiles = []string{
	"go.mod",
	"go.sum",
	"selfsource.go",
	filepath.Join("cmd", "quine", "main.go"),
	filepath.Join("internal", "runtime", "runtime.go"),
}

var selfSourceGitSentinelFiles = []string{
	"selfsource_bundle_data.go",
	filepath.Join(".git", "HEAD"),
	filepath.Join(".git", "index"),
}

type selfSourceManifest struct {
	Format     string `json:"format"`
	Projection string `json:"projection"`
	Files      int    `json:"files"`
	Size       int    `json:"size"`
	SHA256     string `json:"sha256"`
}

func (r *Runtime) syncSelfSourceSurface(agentRoot string) error {
	sourceRoot := filepath.Join(agentRoot, "source-code")
	projection := r.cfg.SelfSourceProjectionMode()
	if projection == "none" {
		if err := removeSelfSourceSurface(sourceRoot); err != nil {
			return fmt.Errorf("remove disabled self-source surface: %w", err)
		}
		return nil
	}
	if selfSourceSurfaceReady(sourceRoot, projection) {
		return nil
	}
	if err := removeSelfSourceSurface(sourceRoot); err != nil {
		return fmt.Errorf("clear incomplete self-source surface: %w", err)
	}
	if err := os.MkdirAll(sourceRoot, 0o755); err != nil {
		return fmt.Errorf("mkdir self-source root: %w", err)
	}
	if err := materializeSelfSourceSurface(sourceRoot, projection); err != nil {
		_ = removeSelfSourceSurface(sourceRoot)
		return err
	}
	if err := setReadOnlyTree(sourceRoot); err != nil {
		_ = removeSelfSourceSurface(sourceRoot)
		return err
	}
	return nil
}

func selfSourceSurfaceReady(root, projection string) bool {
	want, err := currentSelfSourceManifestForProjection(projection)
	if err != nil {
		return false
	}
	data, err := os.ReadFile(selfSourceManifestPath(root))
	if err != nil {
		return false
	}
	var manifest selfSourceManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return false
	}
	if manifest != want {
		return false
	}
	sentinelFiles := append(append([]string{}, selfSourceSentinelFiles...), selfSourceGitSentinelFiles...)
	for _, rel := range sentinelFiles {
		info, err := os.Stat(filepath.Join(root, rel))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

func materializeSelfSourceSurface(root, projection string) error {
	switch projection {
	case "runtime":
		return materializeRuntimeSelfSourceSurface(root)
	case "repo":
		return materializeSelfSourceRepoSurface(root)
	default:
		return fmt.Errorf("unsupported self-source projection %q", projection)
	}
}

func materializeSelfSourceRepoSurface(root string) error {
	bundle, err := quineroot.ReadSelfSourceBundle()
	if err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp(filepath.Dir(root), ".quine-source-bundle-*")
	if err != nil {
		return fmt.Errorf("create temporary self-source bundle dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	bundlePath := filepath.Join(tmpDir, "source.bundle")
	if err := os.WriteFile(bundlePath, bundle, 0o600); err != nil {
		return fmt.Errorf("write temporary self-source bundle: %w", err)
	}
	if err := runSelfSourceGit("", "clone", "-q", bundlePath, root); err != nil {
		return err
	}
	if err := overlayEmbeddedSelfSource(root); err != nil {
		return err
	}
	if err := writeSelfSourceBundleBuildData(root, bundle); err != nil {
		return err
	}
	if err := commitSelfSourceOverlay(root); err != nil {
		return err
	}
	return writeSelfSourceManifest(root, "repo")
}

func materializeRuntimeSelfSourceSurface(root string) error {
	if err := runSelfSourceGit("", "init", "-q", root); err != nil {
		return err
	}
	if err := overlayEmbeddedSelfSource(root); err != nil {
		return err
	}
	if err := commitSelfSourceOverlay(root); err != nil {
		return err
	}
	return writeSelfSourceManifest(root, "runtime")
}

func overlayEmbeddedSelfSource(root string) error {
	return quineroot.WalkSelfSource(func(rel string, d fs.DirEntry, _ error) error {
		target := filepath.Join(root, filepath.FromSlash(rel))
		if d.IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("mkdir embedded source overlay dir %s: %w", target, err)
			}
			return nil
		}
		data, err := quineroot.ReadSelfSource(rel)
		if err != nil {
			return fmt.Errorf("read embedded self-source %s: %w", rel, err)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("mkdir parent for embedded source overlay %s: %w", target, err)
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			return fmt.Errorf("write embedded source overlay %s: %w", target, err)
		}
		return nil
	})
}

func writeSelfSourceBundleBuildData(root string, bundle []byte) error {
	data := fmt.Sprintf(`package quine

import "encoding/base64"

const embeddedSelfSourceBundleBase64 = %q

var SelfSourceBundle = mustDecodeSelfSourceBundle(embeddedSelfSourceBundleBase64)

func mustDecodeSelfSourceBundle(encoded string) []byte {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err)
	}
	return data
}
`, base64.StdEncoding.EncodeToString(bundle))
	if err := os.WriteFile(filepath.Join(root, "selfsource_bundle_data.go"), []byte(data), 0o644); err != nil {
		return fmt.Errorf("write self-source bundle build data: %w", err)
	}
	return nil
}

func commitSelfSourceOverlay(root string) error {
	if err := runSelfSourceGit(root, "config", "user.name", "Quine Runtime"); err != nil {
		return err
	}
	if err := runSelfSourceGit(root, "config", "user.email", "quine-runtime@example.invalid"); err != nil {
		return err
	}
	if err := runSelfSourceGit(root, "add", "-A"); err != nil {
		return err
	}
	changed, err := selfSourceGitHasStagedChanges(root)
	if err != nil {
		return err
	}
	if changed {
		return runSelfSourceGit(root, "commit", "-q", "--no-gpg-sign", "-m", "Overlay live embedded source")
	}
	return nil
}

func selfSourceGitHasStagedChanges(root string) (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return false, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return true, nil
	}
	return false, fmt.Errorf("git diff --cached --quiet: %w: %s", err, strings.TrimSpace(stderr.String()))
}

func writeSelfSourceManifest(root, projection string) error {
	currentManifest, err := currentSelfSourceManifestForProjection(projection)
	if err != nil {
		return err
	}
	manifest, err := json.MarshalIndent(currentManifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal self-source manifest: %w", err)
	}
	manifest = append(manifest, '\n')
	if err := os.WriteFile(selfSourceManifestPath(root), manifest, 0o644); err != nil {
		return fmt.Errorf("write self-source manifest: %w", err)
	}
	return nil
}

func currentSelfSourceManifest() (selfSourceManifest, error) {
	return currentSelfSourceManifestForProjection("repo")
}

func currentSelfSourceManifestForProjection(projection string) (selfSourceManifest, error) {
	bundle, err := quineroot.ReadSelfSourceBundle()
	if err != nil {
		return selfSourceManifest{}, err
	}
	hash := sha256.New()
	hash.Write([]byte("projection\x00" + projection + "\x00"))
	size := 0
	if projection == "repo" {
		hash.Write([]byte("bundle\x00"))
		hash.Write(bundle)
		size += len(bundle)
	} else if projection != "runtime" {
		return selfSourceManifest{}, fmt.Errorf("unsupported self-source projection %q", projection)
	}
	hash.Write([]byte("\x00embedded\x00"))
	embeddedSize, err := hashEmbeddedSelfSource(hash)
	if err != nil {
		return selfSourceManifest{}, err
	}
	return selfSourceManifest{
		Format:     "quine-source-repo/v1",
		Projection: projection,
		Size:       size + embeddedSize,
		SHA256:     hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

func hashEmbeddedSelfSource(hash hashWriter) (int, error) {
	size := 0
	if err := quineroot.WalkSelfSource(func(rel string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}
		data, err := quineroot.ReadSelfSource(rel)
		if err != nil {
			return fmt.Errorf("read embedded self-source %s: %w", rel, err)
		}
		hash.Write([]byte(rel))
		hash.Write([]byte{0})
		hash.Write(data)
		hash.Write([]byte{0})
		size += len(data)
		return nil
	}); err != nil {
		return 0, err
	}
	return size, nil
}

type hashWriter interface {
	Write([]byte) (int, error)
}

func selfSourceManifestPath(root string) string {
	return filepath.Join(root, ".git", "quine-source-manifest.json")
}

func runSelfSourceGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if strings.TrimSpace(dir) != "" {
		cmd.Dir = dir
	}
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func setReadOnlyTree(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		mode := os.FileMode(0o444)
		if info.IsDir() {
			mode = 0o555
		}
		if err := os.Chmod(path, mode); err != nil {
			return fmt.Errorf("chmod %s to %o: %w", path, mode, err)
		}
		return nil
	})
}

func removeSelfSourceSurface(root string) error {
	if _, err := os.Stat(root); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err := setWritableTree(root); err != nil {
		return err
	}
	if err := os.RemoveAll(root); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func setWritableTree(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		mode := os.FileMode(0o644)
		if info.IsDir() {
			mode = 0o755
		}
		if err := os.Chmod(path, mode); err != nil {
			return fmt.Errorf("chmod %s to %o: %w", path, mode, err)
		}
		return nil
	})
}
