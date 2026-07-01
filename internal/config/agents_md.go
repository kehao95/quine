package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const agentsMDFileName = "AGENTS.md"

// DiscoverSingleAgentsMD returns the one project AGENTS.md visible from the
// configured working surface. Multiple discoverable files are rejected until
// hierarchical AGENTS support lands explicitly.
func DiscoverSingleAgentsMD(cfg *Config) (string, error) {
	start, ceiling, err := discoverAgentsMDSearchBounds(cfg)
	if err != nil {
		return "", err
	}

	matches := make([]string, 0, 1)
	for dir := start; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, agentsMDFileName)
		info, err := os.Stat(candidate)
		if err == nil {
			if info.IsDir() {
				return "", fmt.Errorf("%s exists but is a directory", candidate)
			}
			matches = append(matches, candidate)
			if len(matches) > 1 {
				return "", fmt.Errorf("multiple AGENTS.md files discovered (%s, %s); hierarchical AGENTS.md is not supported yet", matches[0], matches[1])
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("stat AGENTS.md %q: %w", candidate, err)
		}

		if dir == ceiling {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	if len(matches) == 0 {
		return "", nil
	}
	return matches[0], nil
}

func discoverAgentsMDSearchBounds(cfg *Config) (string, string, error) {
	start := ""
	ceiling := ""
	if cfg != nil && cfg.WorkspaceEnabled {
		root := strings.TrimSpace(cfg.WorkspaceRoot)
		if root == "" {
			return "", "", fmt.Errorf("QUINE_WORKSPACE_ROOT is required when QUINE_AGENTS_MD_ENABLED=1 and workspace physics are enabled")
		}
		rootAbs, err := filepath.Abs(root)
		if err != nil {
			return "", "", fmt.Errorf("absolute workspace root %q: %w", root, err)
		}
		ceiling = rootAbs
		start = rootAbs
		if work := strings.TrimSpace(cfg.WorkDir); work != "" {
			workAbs, err := filepath.Abs(work)
			if err != nil {
				return "", "", fmt.Errorf("absolute AGENTS.md discovery start %q: %w", work, err)
			}
			rel, err := filepath.Rel(rootAbs, workAbs)
			if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
				start = workAbs
			}
		}
		startDir, _, err := ensureSearchDir(start)
		if err != nil {
			return "", "", err
		}
		return startDir, ceiling, nil
	}

	if cfg != nil {
		start = strings.TrimSpace(cfg.WorkDir)
	}
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("get working directory for AGENTS.md discovery: %w", err)
		}
	}

	startDir, _, err := ensureSearchDir(start)
	if err != nil {
		return "", "", err
	}
	return startDir, volumeRoot(startDir), nil
}

func ensureSearchDir(path string) (string, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", fmt.Errorf("stat AGENTS.md discovery start %q: %w", path, err)
	}
	if !info.IsDir() {
		path = filepath.Dir(path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", "", fmt.Errorf("absolute AGENTS.md discovery start %q: %w", path, err)
	}
	return abs, volumeRoot(abs), nil
}

func volumeRoot(path string) string {
	vol := filepath.VolumeName(path)
	if vol == "" {
		return string(os.PathSeparator)
	}
	return vol + string(os.PathSeparator)
}
