package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kehao95/quine/internal/tools"
)

func (r *Runtime) ensureIncarnation() error {
	if r.incarnationID >= 0 {
		r.cfg.IncarnationID = r.incarnationID
		return nil
	}
	retainedRoot := r.cfg.SessionRetainedDir("")
	// Exec re-entry: a successor started by exec carries a staged bootstrap
	// context (QUINE_CONTEXT_BOOTSTRAP). When it re-enters a session that already
	// has retained incarnation state, it is the next incarnation of that lineage,
	// not a re-run of the current body. Advance to the next id so the successor is
	// projected as inc/N+1 and inherits the prior incarnation's context, instead
	// of overwriting inc/N in place. Session-resume and fork children do not hit
	// this: resume carries no bootstrap context, and a fork child opens a fresh
	// session with no retained incarnation state yet.
	if execReentryStaged() && hasRetainedSessionState(retainedRoot) {
		id, err := nextSessionIncarnationID(retainedRoot)
		if err != nil {
			return fmt.Errorf("compute next incarnation id: %w", err)
		}
		r.incarnationID = id
		r.cfg.IncarnationID = id
		return nil
	}
	id, ok, err := currentSessionIncarnationID(retainedRoot)
	if err != nil {
		return fmt.Errorf("read current incarnation id: %w", err)
	}
	if !ok {
		id = 0
	}
	r.incarnationID = id
	r.cfg.IncarnationID = id
	return nil
}

// execReentryStaged reports whether this process was started by an exec handoff
// that staged a bootstrap context for adoption. It is the signal that the
// current body is a successor incarnation rather than a fresh or resumed body.
func execReentryStaged() bool {
	return strings.TrimSpace(os.Getenv(tools.ContextBootstrapEnv)) != ""
}

func (r *Runtime) currentIncarnationRoot() string {
	if r.incarnationID < 0 {
		// Path projections must stay pinned to the same incarnation even when a
		// caller touches the surface before the first explicit sync. Best-effort
		// lazy allocation keeps those pre-sync writes from accidentally creating a
		// phantom prior incarnation and forcing the real sync onto the next id.
		_ = r.ensureIncarnation()
	}
	return r.cfg.SessionIncarnationPath("", r.cfg.IncarnationID)
}

func (r *Runtime) currentIncarnationContextRoot() string {
	return filepath.Join(r.currentIncarnationRoot(), "context")
}

func (r *Runtime) currentIncarnationPromptRoot() string {
	return filepath.Join(r.currentIncarnationContextRoot(), "prompt")
}

func (r *Runtime) currentIncarnationPromptFile(name string) string {
	return filepath.Join(r.currentIncarnationPromptRoot(), name)
}

func (r *Runtime) currentIncarnationStateRoot() string {
	return filepath.Join(r.currentIncarnationContextRoot(), "state")
}

func (r *Runtime) currentIncarnationMissionPath() string {
	return filepath.Join(r.currentIncarnationRoot(), "mission.txt")
}

func nextSessionIncarnationID(retainedRoot string) (int, error) {
	maxID := -1

	statusPath := filepath.Join(retainedRoot, "status", "session.json")
	statusData, err := os.ReadFile(statusPath)
	if err == nil {
		var status struct {
			IncarnationID int `json:"incarnation_id"`
		}
		if err := json.Unmarshal(statusData, &status); err == nil && status.IncarnationID > maxID {
			maxID = status.IncarnationID
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("read prior session status %s: %w", statusPath, err)
	}

	entries, err := os.ReadDir(filepath.Join(retainedRoot, "inc"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return maxID + 1, nil
		}
		return 0, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1, nil
}

func currentSessionIncarnationID(retainedRoot string) (int, bool, error) {
	currentPath := filepath.Join(retainedRoot, "inc", "current")
	if target, err := os.Readlink(currentPath); err == nil {
		id, parseErr := strconv.Atoi(filepath.Base(target))
		if parseErr == nil {
			return id, true, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, false, fmt.Errorf("read current incarnation link %s: %w", currentPath, err)
	}

	statusPath := filepath.Join(retainedRoot, "status", "session.json")
	statusData, err := os.ReadFile(statusPath)
	if err == nil {
		var status struct {
			IncarnationID int `json:"incarnation_id"`
		}
		if err := json.Unmarshal(statusData, &status); err == nil {
			return status.IncarnationID, true, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, false, fmt.Errorf("read prior session status %s: %w", statusPath, err)
	}

	entries, err := os.ReadDir(filepath.Join(retainedRoot, "inc"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, false, nil
		}
		return 0, false, err
	}
	maxID := -1
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}
	if maxID < 0 {
		return 0, false, nil
	}
	return maxID, true, nil
}
