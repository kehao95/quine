package runtime

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kehao95/quine/internal/tools"
)

func (r *Runtime) importRetainedStateFromAgentRootIfNeeded() error {
	if r == nil || r.cfg == nil {
		return nil
	}
	if hasRetainedSessionState(r.cfg.SessionRetainedDir("")) {
		return nil
	}
	agentContext := filepath.Join(r.cfg.AgentRoot(), "context")
	sourceContext := agentContext
	if resolved, err := filepath.EvalSymlinks(agentContext); err == nil && resolved != "" {
		sourceContext = resolved
	}
	info, err := os.Stat(sourceContext)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat agent context surface: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("agent context surface %s is not a directory", sourceContext)
	}
	contextRoot := filepath.Join(r.cfg.SessionIncarnationPath("", 0), "context")
	if err := os.RemoveAll(contextRoot); err != nil {
		return fmt.Errorf("reset imported context root: %w", err)
	}
	if err := tools.CopyTreePreservingSymlinks(sourceContext, contextRoot); err != nil {
		return fmt.Errorf("import agent context surface: %w", err)
	}
	if err := replaceRelativeSymlinkAtomic(r.cfg.SessionCurrentIncarnationPath(""), "0"); err != nil {
		return fmt.Errorf("sync imported current incarnation link: %w", err)
	}
	return nil
}

func hasRetainedSessionState(retainedRoot string) bool {
	for _, path := range []string{
		filepath.Join(retainedRoot, "inc"),
		filepath.Join(retainedRoot, "status", "session.json"),
	} {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}
