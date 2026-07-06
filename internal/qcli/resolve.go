package qcli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const ContractVersion = "qcli/1"

type ControlAction string

const (
	ControlActionPost      ControlAction = "post"
	ControlActionPoke      ControlAction = "poke"
	ControlActionInject    ControlAction = "inject"
	ControlActionInterrupt ControlAction = "interrupt"
)

var (
	ErrTargetRequired  = errors.New("qcli: exactly one target selector is required")
	ErrTargetNotFound  = errors.New("qcli: target not found")
	ErrNoPeer          = errors.New("qcli: no peer attached")
	ErrEmptyPayload    = errors.New("qcli: empty payload")
	ErrUnreachable     = errors.New("qcli: peer unreachable")
	ErrRegisterTimeout = errors.New("qcli: register timeout")
)

type Agent struct {
	RuntimeRoot    string `json:"runtime_root"`
	AgentRoot      string `json:"agent_root"`
	PublicRoot     string `json:"public_root"`
	Session        string `json:"session"`
	RunID          string `json:"run_id"`
	PID            int    `json:"pid"`
	PPID           int    `json:"ppid,omitempty"`
	ParentSession  string `json:"parent_session,omitempty"`
	Incarnation    int    `json:"incarnation,omitempty"`
	Model          string `json:"model,omitempty"`
	Depth          int    `json:"depth,omitempty"`
	SessionPath    string `json:"session_path"`
	ControlPath    string `json:"control_path"`
	InboxPath      string `json:"inbox_path"`
	ControlLogPath string `json:"control_log_path"`
	ContextPath    string `json:"context_path"`
	LivePath       string `json:"live_path"`
}

type PeerInfo struct {
	Session     string `json:"session"`
	RunID       string `json:"run_id"`
	PID         int    `json:"pid"`
	AgentRoot   string `json:"agent_root"`
	RuntimeRoot string `json:"runtime_root"`
}

func (a Agent) PeerInfo() PeerInfo {
	return PeerInfo{
		Session:     a.Session,
		RunID:       a.RunID,
		PID:         a.PID,
		AgentRoot:   a.AgentRoot,
		RuntimeRoot: a.RuntimeRoot,
	}
}

type ResolveOptions struct {
	RuntimeRoot string `json:"runtime_root,omitempty"`
	AgentRoot   string `json:"agent_root,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Session     string `json:"session,omitempty"`
	SearchFrom  string `json:"-"`
}

type sessionStatus struct {
	SessionID      string `json:"session_id"`
	RunID          string `json:"run_id"`
	IncarnationID  int    `json:"incarnation_id"`
	PID            int    `json:"pid"`
	PPID           int    `json:"ppid"`
	ParentSession  string `json:"parent_session"`
	RuntimeRoot    string `json:"runtime_root"`
	AgentRoot      string `json:"agent_root"`
	ModelID        string `json:"model_id"`
	Depth          int    `json:"depth"`
	RetentionDir   string `json:"retention_dir"`
	WorkspaceRoot  string `json:"workspace_root"`
	Workspace      string `json:"workspace"`
	WorkspaceOwner bool   `json:"workspace_owner"`
}

type inboxStatus struct {
	PendingCount int `json:"pending_count"`
}

func Resolve(opts ResolveOptions) (Agent, error) {
	selectors := 0
	if strings.TrimSpace(opts.AgentRoot) != "" {
		selectors++
	}
	if opts.PID > 0 {
		selectors++
	}
	if strings.TrimSpace(opts.Session) != "" {
		selectors++
	}
	if selectors != 1 {
		return Agent{}, ErrTargetRequired
	}

	if strings.TrimSpace(opts.AgentRoot) != "" {
		return loadAgentFromRoot(opts.AgentRoot)
	}

	roots, err := candidateRuntimeRoots(strings.TrimSpace(opts.RuntimeRoot), opts.SearchFrom)
	if err != nil {
		return Agent{}, err
	}
	for _, root := range roots {
		var agentRoot string
		switch {
		case opts.PID > 0:
			link := filepath.Join(root, "pid", strconv.Itoa(opts.PID))
			resolved, err := filepath.EvalSymlinks(link)
			if err != nil {
				continue
			}
			agentRoot = resolved
		case strings.TrimSpace(opts.Session) != "":
			agentRoot = filepath.Join(root, "agent", strings.TrimSpace(opts.Session))
			if _, err := os.Stat(filepath.Join(agentRoot, "status", "session.json")); err != nil {
				if _, publicErr := os.Stat(filepath.Join(agentRoot, "public", "status", "session.json")); publicErr != nil {
					continue
				}
			}
		}
		agent, err := loadAgentFromRoot(agentRoot)
		if err == nil {
			return agent, nil
		}
	}
	return Agent{}, fmt.Errorf("%w: pid=%d session=%q runtime_root=%q", ErrTargetNotFound, opts.PID, opts.Session, opts.RuntimeRoot)
}

func DiscoverRuntimeRoot(explicit string) (string, error) {
	return defaultRuntimeRoot(explicit, "")
}

func defaultRuntimeRoot(explicitRoot, workDir string) (string, error) {
	if strings.TrimSpace(explicitRoot) != "" {
		return filepath.Abs(strings.TrimSpace(explicitRoot))
	}
	if envRoot := strings.TrimSpace(os.Getenv("QCLI_RUNTIME_ROOT")); envRoot != "" {
		return filepath.Abs(envRoot)
	}
	if envRoot := strings.TrimSpace(os.Getenv("QUINE_DATA_DIR")); envRoot != "" {
		return filepath.Abs(envRoot)
	}
	base := workDir
	if strings.TrimSpace(base) == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("qcli: getwd: %w", err)
		}
		base = wd
	}
	return filepath.Abs(filepath.Join(base, ".quine"))
}

func candidateRuntimeRoots(explicitRoot, searchFrom string) ([]string, error) {
	if strings.TrimSpace(explicitRoot) != "" {
		root, err := filepath.Abs(strings.TrimSpace(explicitRoot))
		if err != nil {
			return nil, fmt.Errorf("qcli: resolve runtime root: %w", err)
		}
		return []string{root}, nil
	}

	var roots []string
	seen := map[string]struct{}{}
	add := func(path string) {
		if path == "" {
			return
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return
		}
		if _, ok := seen[abs]; ok {
			return
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			seen[abs] = struct{}{}
			roots = append(roots, abs)
		}
	}
	add(os.Getenv("QCLI_RUNTIME_ROOT"))
	add(os.Getenv("QUINE_DATA_DIR"))

	start := strings.TrimSpace(searchFrom)
	if start == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("qcli: getwd: %w", err)
		}
		start = wd
	}
	start, err := filepath.Abs(start)
	if err != nil {
		return nil, fmt.Errorf("qcli: resolve search root: %w", err)
	}
	dir := start
	for {
		add(filepath.Join(dir, ".quine"))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("%w: no runtime root candidates discovered", ErrTargetNotFound)
	}
	return roots, nil
}

func loadAgentFromRoot(root string) (Agent, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return Agent{}, fmt.Errorf("qcli: empty agent root")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Agent{}, fmt.Errorf("qcli: resolve agent root: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(absRoot); err == nil {
		absRoot = resolved
	}

	statusPath := filepath.Join(absRoot, "status", "session.json")
	data, err := os.ReadFile(statusPath)
	if err != nil && filepath.Base(absRoot) != "public" {
		statusPath = filepath.Join(absRoot, "public", "status", "session.json")
		data, err = os.ReadFile(statusPath)
	}
	if err != nil {
		return Agent{}, fmt.Errorf("qcli: read session status: %w", err)
	}

	var status sessionStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return Agent{}, fmt.Errorf("qcli: parse session status: %w", err)
	}
	if status.PID <= 0 {
		return Agent{}, fmt.Errorf("qcli: invalid pid in %s", statusPath)
	}

	loadedRoot := absRoot
	canonicalRoot := strings.TrimSpace(status.AgentRoot)
	if canonicalRoot != "" {
		canonicalRoot, err = filepath.Abs(canonicalRoot)
		if err != nil {
			return Agent{}, fmt.Errorf("qcli: resolve canonical agent root: %w", err)
		}
		if resolvedRoot, resolveErr := filepath.EvalSymlinks(canonicalRoot); resolveErr == nil {
			canonicalRoot = resolvedRoot
		}
	} else if filepath.Base(loadedRoot) == "public" {
		canonicalRoot = filepath.Dir(loadedRoot)
	} else if filepath.Base(filepath.Dir(statusPath)) == "status" && filepath.Base(filepath.Dir(filepath.Dir(statusPath))) == "public" {
		canonicalRoot = filepath.Dir(filepath.Dir(filepath.Dir(statusPath)))
	} else {
		canonicalRoot = loadedRoot
	}

	runtimeRoot := strings.TrimSpace(status.RuntimeRoot)
	if runtimeRoot == "" {
		runtimeRoot = filepath.Dir(filepath.Dir(canonicalRoot))
	}
	publicRoot := filepath.Join(canonicalRoot, "public")
	if filepath.Base(loadedRoot) == "public" {
		publicRoot = loadedRoot
	}
	if resolvedPublic, err := filepath.EvalSymlinks(publicRoot); err == nil && filepath.Base(resolvedPublic) == "public" {
		publicRoot = resolvedPublic
	}

	return Agent{
		RuntimeRoot:    runtimeRoot,
		AgentRoot:      canonicalRoot,
		PublicRoot:     publicRoot,
		Session:        status.SessionID,
		RunID:          status.RunID,
		PID:            status.PID,
		PPID:           status.PPID,
		ParentSession:  status.ParentSession,
		Incarnation:    status.IncarnationID,
		Model:          status.ModelID,
		Depth:          status.Depth,
		SessionPath:    filepath.Join(publicRoot, "status", "session.json"),
		ControlPath:    filepath.Join(publicRoot, "ctl"),
		InboxPath:      filepath.Join(publicRoot, "status", "inbox.json"),
		ControlLogPath: filepath.Join(publicRoot, "log", "control.jsonl"),
		ContextPath:    filepath.Join(canonicalRoot, "context", "state", "current.jsonl"),
		LivePath:       filepath.Join(canonicalRoot, "context", "state", "live.jsonl"),
	}, nil
}

func pidLive(pid int) bool {
	if pid <= 0 {
		return false
	}
	if syscall.Kill(pid, 0) != nil {
		return false
	}
	return !pidZombie(pid)
}

// pidZombie reports whether pid is a defunct (zombie) process. Signal-0
// succeeds on zombies — an unreaped exited agent would otherwise read as live
// forever — so liveness must also consult the /proc state field. On platforms
// without /proc the signal-0 verdict stands.
func pidZombie(pid int) bool {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return false
	}
	// Layout: pid (comm) state ... — comm may contain spaces and parens, so
	// the state byte sits two past the last ')'.
	i := bytes.LastIndexByte(data, ')')
	if i < 0 || i+2 >= len(data) {
		return false
	}
	return data[i+2] == 'Z'
}

func verifyLivePID(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("qcli: invalid pid %d", pid)
	}
	if err := syscall.Kill(pid, 0); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return fmt.Errorf("%w: pid %d is not live", ErrUnreachable, pid)
		}
		return fmt.Errorf("qcli: pid %d liveness check failed: %w", pid, err)
	}
	if pidZombie(pid) {
		return fmt.Errorf("%w: pid %d is defunct", ErrUnreachable, pid)
	}
	return nil
}

func validControlAction(action ControlAction) bool {
	switch action {
	case ControlActionPost, ControlActionPoke, ControlActionInject, ControlActionInterrupt:
		return true
	default:
		return false
	}
}

func controlPathForAction(agent Agent, action ControlAction) (string, bool, error) {
	controlDir := strings.TrimSpace(agent.ControlPath)
	if controlDir == "" {
		return "", false, fmt.Errorf("qcli: missing control path")
	}
	switch action {
	case ControlActionPost:
		return filepath.Join(controlDir, "post"), false, nil
	case ControlActionPoke:
		return filepath.Join(controlDir, "poke"), true, nil
	case ControlActionInject:
		return filepath.Join(controlDir, "inject"), false, nil
	case ControlActionInterrupt:
		return filepath.Join(controlDir, "interrupt"), true, nil
	default:
		return "", false, fmt.Errorf("qcli: unknown control action %q", action)
	}
}
