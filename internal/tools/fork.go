package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// ForkExecutor spawns child quine processes with bootstrapped context.
type ForkExecutor struct {
	// QuinePath is the path to the quine binary. Defaults to "./quine" or
	// the current executable if not set.
	QuinePath string

	// DataDir is the durable runtime-state root shared by the process tree.
	// It holds tapes, job directories, and coordination state.
	DataDir string

	// SessionID is the current session's ID (becomes PARENT_SESSION for child).
	SessionID string

	// Env contains environment variables for child processes.
	// Should include QUINE_* vars with incremented depth.
	Env []string

	// ContextRoot is the current session's live cognition surface under AgentRoot.
	ContextRoot string

	// DefaultTimeout is the maximum time to wait for the entire swarm
	// operation (when wait=true). Not per-child.
	DefaultTimeout time.Duration

	// MaxOutput limits the captured output size per child.
	MaxOutput int

	// ProcessStarted is called when a child process starts.
	ProcessStarted func(*os.Process)

	// ProcessEnded is called when a child process ends.
	ProcessEnded func()

	WorkspaceEnabled       bool
	WorkspaceRoot          string
	Workspace              string
	WorkspaceBackend       string
	WorkspaceOverlayDriver string
	WorkspaceRevisionMode  config.WorkspaceRevisionMode
	WorkDir                string
	TurnID                 int
	Mission                string

	ForkWorldEnabled           bool
	FSMutationTelemetryEnabled bool

	subjective *subjectiveFS
}

// NewForkExecutor creates a ForkExecutor from config.
//
// A fork child inherits this process's environment minus the runtime-owned
// names, with config/env/override applied, plus the stamps that make it a
// member of this agent tree (config.ForkChildStamps: depth+1, parent session,
// the runtime root it must join, the workspace coordinates it must not
// re-derive). Its own QUINE_SESSION_ID is injected per child in
// launchChildSession — only that call knows each child's id.
//
// Env is rebuilt before every fork/spawn call (runtime.refreshForkEnv), so the
// override is never stale.
func NewForkExecutor(cfg *config.Config) *ForkExecutor {
	contextRoot := filepath.Join(cfg.AgentRoot(), "context")

	return &ForkExecutor{
		QuinePath:                  strings.TrimSpace(cfg.SelfReentryTarget),
		DataDir:                    cfg.DataDir,
		SessionID:                  cfg.SessionID,
		Env:                        ForkChildEnv(cfg, nil),
		ContextRoot:                contextRoot,
		DefaultTimeout:             time.Duration(cfg.ForkDefaultTimeoutSeconds) * time.Second,
		MaxOutput:                  cfg.OutputTruncate,
		WorkspaceEnabled:           cfg.WorkspaceEnabled,
		WorkspaceRoot:              cfg.WorkspaceRoot,
		Workspace:                  cfg.Workspace,
		WorkspaceBackend:           cfg.WorkspaceBackend,
		WorkspaceOverlayDriver:     cfg.WorkspaceOverlayDriver,
		WorkspaceRevisionMode:      cfg.WorkspaceRevisionMode,
		WorkDir:                    cfg.WorkDir,
		ForkWorldEnabled:           cfg.ForkWorldEnabled,
		FSMutationTelemetryEnabled: cfg.FSMutationTelemetryEnabled(),
		subjective:                 newSubjectiveFS(cfg),
	}
}

// ForkChildEnv builds the environment of a managed fork/spawn child: the
// BoundaryChild pipeline (inherit − mask, override, stamps).
//
// It replaces the hand-kept filterProcessIdentity strip list. That list and the
// registry's runtime-emitted class were the same eight names maintained in two
// places; the mask is now derived from the registry, so a new runtime-emitted
// knob cannot be added without also being masked.
//
// A malformed override is reported through onError and then ignored — a child
// still gets a correct, governed environment; only the agent's own additions
// are dropped, and the override surface reports why.
func ForkChildEnv(cfg *config.Config, onError func(error)) []string {
	if cfg == nil {
		return nil
	}
	override, err := config.ReadEnvOverride(cfg.EnvOverridePath())
	if err != nil && onError != nil {
		onError(err)
	}
	return config.BuildChildEnv(config.BoundaryChild, os.Environ(), override, cfg.ForkChildStamps())
}

func (f *ForkExecutor) waitContext() (context.Context, context.CancelFunc) {
	if f.DefaultTimeout <= 0 {
		return context.Background(), func() {}
	}
	return context.WithTimeout(context.Background(), f.DefaultTimeout)
}

// Fork mode constants.
const (
	ForkModeWait   = "wait"   // Block until all children finish
	ForkModeRace   = "race"   // First exit-0 child wins, rest killed (default)
	ForkModeForget = "forget" // Fire-and-forget, return PIDs immediately
)

// ForkRequest represents the parsed arguments from a fork tool call.
type ForkRequest struct {
	Children    []ForkChild
	Mode        string // "race" (default), "wait", or "forget"
	AdoptWinner bool
}

type ForkChild struct {
	Intent     string
	World      config.WorldKind
	Protection config.ProtectionMode
	Scope      string
}

type childContextMode string

const (
	childContextInherited childContextMode = "inherited"
	childContextFresh     childContextMode = "fresh"
)

type childLaunchSpec struct {
	ToolID                 string
	Mission                string
	SessionID              string
	ContextMode            childContextMode
	World                  config.WorldKind
	Protection             config.ProtectionMode
	Scope                  string
	Relation               *forkRelationSurface
	Index                  int
	CaptureRoot            string
	RecordWorkspaceSession bool
	ForkChild              *ForkChild
}

// ParseForkArgs extracts ForkRequest from a ToolCall's Arguments map.
func ParseForkArgs(args map[string]any, cfg *config.Config) (ForkRequest, error) {
	children, err := parseForkChildren(args, cfg)
	if err != nil {
		return ForkRequest{}, err
	}

	req := ForkRequest{
		Children: children,
		Mode:     ForkModeRace, // default
	}

	if v, ok := args["mode"]; ok {
		s, ok := v.(string)
		if !ok {
			return ForkRequest{}, fmt.Errorf("mode must be a string, got %T", v)
		}
		switch s {
		case ForkModeWait, ForkModeRace, ForkModeForget:
			req.Mode = s
		default:
			return ForkRequest{}, fmt.Errorf("mode must be one of wait, race, forget; got %q", s)
		}
	}

	if v, ok := args["adopt_winner"]; ok {
		b, ok := v.(bool)
		if !ok {
			return ForkRequest{}, fmt.Errorf("adopt_winner must be a boolean, got %T", v)
		}
		req.AdoptWinner = b
	}
	if req.AdoptWinner {
		if req.Mode != ForkModeRace {
			return ForkRequest{}, fmt.Errorf("adopt_winner is only valid with mode=%q", ForkModeRace)
		}
		for i, child := range req.Children {
			if !child.adoptable(cfg) {
				return ForkRequest{}, fmt.Errorf("children[%d] is not an adoptable subjective child", i)
			}
		}
	}

	return req, nil
}

func parseForkChildren(args map[string]any, cfg *config.Config) ([]ForkChild, error) {
	raw, ok := args["children"]
	if !ok {
		return nil, fmt.Errorf("missing required argument: children")
	}
	rawSlice, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("children must be an array, got %T", raw)
	}
	if len(rawSlice) == 0 {
		return nil, fmt.Errorf("children must contain at least one entry")
	}

	children := make([]ForkChild, 0, len(rawSlice))
	for i, v := range rawSlice {
		obj, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("children[%d] must be an object, got %T", i, v)
		}

		intent, ok := obj["intent"].(string)
		if !ok {
			return nil, fmt.Errorf("children[%d].intent must be a string", i)
		}
		intent = strings.TrimSpace(intent)
		if intent == "" {
			return nil, fmt.Errorf("children[%d].intent cannot be empty", i)
		}

		world, protection, scope, err := parseChildWorldScopeProtection(obj, i, cfg, true)
		if err != nil {
			return nil, err
		}
		children = append(children, ForkChild{
			Intent:     intent,
			World:      world,
			Protection: protection,
			Scope:      scope,
		})
	}
	return children, nil
}

func parseChildWorldScopeProtection(obj map[string]any, index int, cfg *config.Config, requireScopeWhenWorldDisabled bool) (config.WorldKind, config.ProtectionMode, string, error) {
	worldModesEnabled := cfg != nil && cfg.ForkWorldEnabled
	world := config.WorldKind("")
	if rawWorld, ok := obj["world"]; ok {
		if !worldModesEnabled {
			return "", "", "", fmt.Errorf("children[%d].world is unavailable unless QUINE_FORK_WORLD_ENABLED=1", index)
		}
		worldStr, ok := rawWorld.(string)
		if !ok {
			return "", "", "", fmt.Errorf("children[%d].world must be a string", index)
		}
		switch config.WorldKind(strings.TrimSpace(worldStr)) {
		case config.WorldHost, config.WorldSubjective:
			world = config.WorldKind(strings.TrimSpace(worldStr))
		default:
			return "", "", "", fmt.Errorf("children[%d].world must be one of host, subjective", index)
		}
	}
	if worldModesEnabled && world == "" {
		world = config.WorldSubjective
	}

	protection := config.ProtectionMode("")
	if rawProtection, ok := obj["protection"]; ok {
		if !worldModesEnabled {
			return "", "", "", fmt.Errorf("children[%d].protection is unavailable unless QUINE_FORK_WORLD_ENABLED=1", index)
		}
		protectionStr, ok := rawProtection.(string)
		if !ok {
			return "", "", "", fmt.Errorf("children[%d].protection must be a string", index)
		}
		switch config.ProtectionMode(strings.TrimSpace(protectionStr)) {
		case config.ProtectionNone, config.ProtectionTransactional:
			protection = config.ProtectionMode(strings.TrimSpace(protectionStr))
		default:
			return "", "", "", fmt.Errorf("children[%d].protection must be one of none, transactional", index)
		}
	}
	if worldModesEnabled && protection == "" {
		protection = config.ProtectionTransactional
	}

	scope := "."
	if rawScope, ok := obj["scope"]; ok {
		var ok bool
		scope, ok = rawScope.(string)
		if !ok {
			return "", "", "", fmt.Errorf("children[%d].scope must be a string", index)
		}
		scope = strings.TrimSpace(scope)
		if scope == "" {
			return "", "", "", fmt.Errorf("children[%d].scope cannot be empty", index)
		}
	} else if !worldModesEnabled && requireScopeWhenWorldDisabled {
		return "", "", "", fmt.Errorf("children[%d].scope must be a string", index)
	}

	if filepath.IsAbs(scope) {
		return "", "", "", fmt.Errorf("children[%d].scope must be \".\" or a relative path under the current scope; absolute paths are invalid", index)
	}
	if worldModesEnabled {
		switch {
		case world == config.WorldHost && protection != config.ProtectionNone:
			return "", "", "", fmt.Errorf("children[%d] uses unsupported pair world=\"host\" with protection=%q; host children must use protection=\"none\"", index, protection)
		case world == config.WorldSubjective && protection != config.ProtectionTransactional:
			return "", "", "", fmt.Errorf("children[%d] uses unsupported pair world=\"subjective\" with protection=%q; subjective children must use protection=\"transactional\"", index, protection)
		case world == config.WorldHost && scope != ".":
			return "", "", "", fmt.Errorf("children[%d].scope is only meaningful for world=\"subjective\"; use \".\" for world=\"host\"", index)
		}
	}
	return world, protection, scope, nil
}

// childProcess holds the state for a single child in a swarm fork.
type childProcess struct {
	cmd              *exec.Cmd
	stdoutFile       *os.File // nil for fire-and-forget
	stderrFile       *os.File // nil for fire-and-forget
	stdoutPath       string
	stderrPath       string
	intent           string
	index            int
	sessionID        string
	seedRoot         string
	workspaceSession string
}

type childContextProjection struct {
	ForkToolID    string
	ParentMission string
	Child         ForkChild
}

type projectedContextToolResult struct {
	Tool              string `json:"tool"`
	Status            string `json:"status"`
	Reason            string `json:"reason,omitempty"`
	ParentMissionRef  string `json:"parent_mission_ref,omitempty"`
	CurrentMissionRef string `json:"current_mission_ref,omitempty"`
	ChildAssignment   string `json:"child_assignment,omitempty"`
	AssignmentRef     string `json:"assignment_ref,omitempty"`
	Note              string `json:"note,omitempty"`
}

// childResult captures the outcome of a completed child.
type childResult struct {
	child    *childProcess
	exitCode int
	err      error // non-nil for start/wait failures (not exit code errors)
}

type forkChildStructuredResult struct {
	Index            int    `json:"index"`
	Intent           string `json:"intent"`
	World            string `json:"world,omitempty"`
	Protection       string `json:"protection,omitempty"`
	Scope            string `json:"scope"`
	Status           string `json:"status"`
	SessionID        string `json:"session_id,omitempty"`
	AgentRoot        string `json:"agent_root,omitempty"`
	PublicRoot       string `json:"public_root,omitempty"`
	RetainedRoot     string `json:"retained_root,omitempty"`
	SeedRoot         string `json:"seed_root,omitempty"`
	StatusPath       string `json:"status_path,omitempty"`
	ControlPath      string `json:"control_path,omitempty"`
	ExitCode         *int   `json:"exit_code,omitempty"`
	PID              int    `json:"pid,omitempty"`
	Stdout           string `json:"stdout,omitempty"`
	Stderr           string `json:"stderr,omitempty"`
	WorldHandle      string `json:"world_handle,omitempty"`
	WorldRevision    string `json:"world_revision,omitempty"`
	WorkspaceSession string `json:"workspace_session,omitempty"`
	Adoptable        bool   `json:"adoptable,omitempty"`
	Error            string `json:"error,omitempty"`
}

type forkStructuredResult struct {
	Tool           string                      `json:"tool"`
	Mode           string                      `json:"mode"`
	Status         string                      `json:"status"`
	RelationID     string                      `json:"relation_id,omitempty"`
	RelationRoot   string                      `json:"relation_root,omitempty"`
	RelationHandle string                      `json:"relation_handle,omitempty"`
	Requested      int                         `json:"requested"`
	Spawned        int                         `json:"spawned"`
	Completed      int                         `json:"completed"`
	Succeeded      int                         `json:"succeeded"`
	Killed         int                         `json:"killed"`
	FSMutations    string                      `json:"fs_mutations,omitempty"`
	WorldRevision  string                      `json:"world_revision,omitempty"`
	Winner         *forkChildStructuredResult  `json:"winner,omitempty"`
	Children       []forkChildStructuredResult `json:"children,omitempty"`
	Errors         []string                    `json:"errors,omitempty"`
}

type forkRelationSurface struct {
	ID               string
	Root             string
	Handle           string
	ToolID           string
	Mode             string
	InitiatorSession string
	Requested        int
	CreatedAt        time.Time
}

type forkRelationStatus struct {
	ID               string   `json:"id"`
	Kind             string   `json:"kind"`
	Tool             string   `json:"tool"`
	Mode             string   `json:"mode"`
	Status           string   `json:"status"`
	InitiatorSession string   `json:"initiator_session"`
	Requested        int      `json:"requested"`
	Spawned          int      `json:"spawned"`
	Completed        int      `json:"completed"`
	Succeeded        int      `json:"succeeded"`
	Killed           int      `json:"killed"`
	WinnerIndex      *int     `json:"winner_index,omitempty"`
	Errors           []string `json:"errors,omitempty"`
	UpdatedAt        string   `json:"updated_at"`
}

type forkRelationLogEntry struct {
	At          string `json:"at"`
	Event       string `json:"event"`
	ToolID      string `json:"tool_id,omitempty"`
	Mode        string `json:"mode,omitempty"`
	Status      string `json:"status,omitempty"`
	MemberIndex *int   `json:"member_index,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Detail      string `json:"detail,omitempty"`
}

func (c ForkChild) adoptable(cfg *config.Config) bool {
	if cfg != nil && cfg.ForkWorldEnabled {
		return cfg.WorkspaceTransactional() && c.World == config.WorldSubjective && c.Protection == config.ProtectionTransactional
	}
	return cfg != nil && cfg.WorkspaceTransactional()
}

func newForkSessionID(index int) string {
	return fmt.Sprintf("sess_fork_%d_%d_%d", time.Now().UnixNano(), os.Getpid(), index)
}

func sanitizeRelationID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Sprintf("fork_%d", time.Now().UnixNano())
	}
	var b strings.Builder
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-', r == '_', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	sanitized := strings.Trim(b.String(), "._-")
	if sanitized == "" {
		return fmt.Sprintf("fork_%d", time.Now().UnixNano())
	}
	return sanitized
}

func retentionRootFromEnv(env []string) string {
	for _, entry := range env {
		if strings.HasPrefix(entry, config.EnvRetentionDir+"=") {
			return strings.TrimSpace(strings.TrimPrefix(entry, config.EnvRetentionDir+"="))
		}
	}
	return ""
}

func (f *ForkExecutor) sessionRetainedDir(sessionID string) string {
	if strings.TrimSpace(sessionID) == "" {
		sessionID = f.SessionID
	}
	if root := retentionRootFromEnv(f.Env); root != "" {
		return filepath.Join(root, "sessions", sessionID)
	}
	return filepath.Join(f.DataDir, "log", sessionID)
}

func (f *ForkExecutor) sessionRoot(sessionID string) string {
	if strings.TrimSpace(sessionID) == "" {
		sessionID = f.SessionID
	}
	return filepath.Join(f.DataDir, "agent", sessionID)
}

func (f *ForkExecutor) sessionPublicRoot(sessionID string) string {
	return filepath.Join(f.sessionRoot(sessionID), "public")
}

func (f *ForkExecutor) sessionStatusPath(sessionID string) string {
	return filepath.Join(f.sessionPublicRoot(sessionID), "status", "session.json")
}

func (f *ForkExecutor) sessionControlPath(sessionID string) string {
	return filepath.Join(f.sessionPublicRoot(sessionID), "ctl")
}

func (f *ForkExecutor) beginRelation(toolID string, req ForkRequest) (*forkRelationSurface, error) {
	children := make([]forkChildStructuredResult, 0, len(req.Children))
	for i, child := range req.Children {
		children = append(children, f.structuredChild(child, i))
	}
	return beginSwarmRelation("fork", toolID, req.Mode, f.SessionID, f.sessionRetainedDir, len(req.Children), children)
}

func (f *ForkExecutor) updateRelationStatus(relation *forkRelationSurface, status forkRelationStatus) error {
	return updateSwarmRelationStatus(relation, "fork", status)
}

func (f *ForkExecutor) writeRelationMember(relation *forkRelationSurface, child forkChildStructuredResult) error {
	return writeSwarmRelationMember(relation, child.Index, child)
}

func (f *ForkExecutor) appendRelationLog(relation *forkRelationSurface, entry forkRelationLogEntry) error {
	return appendSwarmRelationLog(relation, entry)
}

func (f *ForkExecutor) relationToolResult(toolID string, relation *forkRelationSurface, payload forkStructuredResult, isError bool) tape.ToolResult {
	completed := 0
	for _, child := range payload.Children {
		switch child.Status {
		case "completed", "error", "spawn_failed", "killed", "timeout", "no_result":
			completed++
		}
	}
	if payload.Winner != nil {
		completed++
	}
	payload.Completed = completed
	if relation != nil {
		payload.RelationID = relation.ID
		payload.RelationRoot = relation.Root
		payload.RelationHandle = relation.Handle
		status := forkRelationStatus{
			Status:           payload.Status,
			Mode:             payload.Mode,
			InitiatorSession: relation.InitiatorSession,
			Requested:        payload.Requested,
			Spawned:          payload.Spawned,
			Succeeded:        payload.Succeeded,
			Killed:           payload.Killed,
			Errors:           append([]string(nil), payload.Errors...),
		}
		if payload.Winner != nil {
			winnerIndex := payload.Winner.Index
			status.WinnerIndex = &winnerIndex
		}
		status.Completed = completed
		if err := writeJSONFile(filepath.Join(relation.Root, "result.json"), payload); err != nil {
			payload.Status = "error"
			payload.Errors = append(payload.Errors, fmt.Sprintf("write relation result: %v", err))
			isError = true
			status.Status = "error"
			status.Errors = append(status.Errors, fmt.Sprintf("write relation result: %v", err))
		}
		if err := f.updateRelationStatus(relation, status); err != nil {
			payload.Status = "error"
			payload.Errors = append(payload.Errors, fmt.Sprintf("write relation status: %v", err))
			isError = true
		}
		_ = f.appendRelationLog(relation, forkRelationLogEntry{
			Event:  "completed",
			Status: payload.Status,
			Detail: fmt.Sprintf("spawned=%d succeeded=%d killed=%d", payload.Spawned, payload.Succeeded, payload.Killed),
		})
	}
	return tape.ToolResult{
		ToolID:  toolID,
		Content: tape.MarshalToolResultContent(payload),
		IsError: isError,
	}
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func appendJSONLine(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = fh.Write(data)
	return err
}

// Execute spawns child quine processes with the given argv missions.
// Routes to the appropriate mode based on req.Mode.
func (f *ForkExecutor) Execute(toolID string, req ForkRequest) tape.ToolResult {
	switch req.Mode {
	case ForkModeRace:
		return f.executeRace(toolID, req)
	case ForkModeForget:
		return f.executeFireAndForget(toolID, req)
	default: // ForkModeWait
		return f.executeGatherAll(toolID, req)
	}
}

func (f *ForkExecutor) initState() error {
	if f.MaxOutput == 0 {
		f.MaxOutput = 20480
	}
	if f.DataDir == "" {
		tmpDir, err := os.MkdirTemp("", "quine-fork-data-*")
		if err != nil {
			return fmt.Errorf("creating temp data dir: %w", err)
		}
		f.DataDir = tmpDir
	}
	if f.SessionID == "" {
		f.SessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}
	if f.subjective == nil && f.WorkspaceEnabled {
		f.subjective = &subjectiveFS{
			enabled:          true,
			workspaceRoot:    f.WorkspaceRoot,
			workspace:        f.Workspace,
			workspaceBackend: f.WorkspaceBackend,
			overlayDriver:    f.WorkspaceOverlayDriver,
			revisionMode:     f.WorkspaceRevisionMode,
			dataDir:          f.DataDir,
			workspaceSession: f.SessionID,
			workspaceOwner:   true,
		}
	}
	if f.subjective != nil {
		if err := f.subjective.init(f.DataDir, f.SessionID); err != nil {
			return fmt.Errorf("initializing workspace physics: %w", err)
		}
	}
	return nil
}

func (f *ForkExecutor) launchChildSession(ctx context.Context, spec childLaunchSpec, capture bool) (*childProcess, error) {
	sessionID := strings.TrimSpace(spec.SessionID)
	if sessionID == "" {
		sessionID = newForkSessionID(spec.Index)
	}
	childIntent := strings.TrimSpace(spec.Mission)
	if spec.ForkChild != nil && strings.TrimSpace(spec.ForkChild.Intent) != "" {
		childIntent = strings.TrimSpace(spec.ForkChild.Intent)
	}
	cp := &childProcess{
		intent:    childIntent,
		index:     spec.Index,
		sessionID: sessionID,
	}

	if spec.ContextMode == childContextInherited {
		cp.seedRoot = filepath.Join(f.sessionRetainedDir(sessionID), "seed")
		child := ForkChild{
			Intent:     spec.Mission,
			World:      spec.World,
			Protection: spec.Protection,
			Scope:      spec.Scope,
		}
		if spec.ForkChild != nil {
			child = *spec.ForkChild
		}
		if _, err := f.writeChildSeedSurface(sessionID, spec.Relation, &childContextProjection{
			ForkToolID:    strings.TrimSpace(spec.ToolID),
			ParentMission: strings.TrimSpace(f.Mission),
			Child:         child,
		}, spec.Index); err != nil {
			return cp, fmt.Errorf("write child %d seed surface: %w", spec.Index, err)
		}
	}

	env := MergeEnv(append([]string(nil), f.Env...), []string{config.EnvSessionID + "=" + sessionID})
	if f.ForkWorldEnabled {
		switch spec.World {
		case config.WorldHost:
			env = filterWorkspacePhysics(env)
		case config.WorldSubjective:
			if !f.WorkspaceEnabled {
				return cp, fmt.Errorf("child %d requested world=\"subjective\", but workspace physics are not enabled", spec.Index)
			}
			if runtime.GOOS != "linux" {
				return cp, fmt.Errorf("child %d requested world=\"subjective\", but Linux workspace physics are unavailable", spec.Index)
			}
			root, current := f.currentWorkspaceRootAndPath()
			workspace, err := resolveChildWorkspace(root, current, spec.Scope)
			if err != nil {
				return cp, fmt.Errorf("resolve child %d scope: %w", spec.Index, err)
			}
			env = MergeEnv(env, []string{
				config.EnvWorkspaceRoot + "=" + root,
				config.EnvWorkspace + "=" + workspace,
				config.EnvWorkspaceBackend + "=" + f.workspaceChildBackend(),
				config.EnvWorkspaceRevisionMode + "=" + string(f.workspaceChildRevisionMode()),
				config.EnvWorkspaceOwner + "=0",
			})
		default:
			return cp, fmt.Errorf("child %d has unsupported world %q", spec.Index, spec.World)
		}
	} else if f.WorkspaceEnabled {
		workspace, err := resolveChildWorkspace(f.WorkspaceRoot, f.Workspace, spec.Scope)
		if err != nil {
			return cp, fmt.Errorf("resolve child %d scope: %w", spec.Index, err)
		}
		env = MergeEnv(env, []string{
			config.EnvWorkspaceRoot + "=" + f.WorkspaceRoot,
			config.EnvWorkspace + "=" + workspace,
			config.EnvWorkspaceOwner + "=0",
		})
	} else if strings.TrimSpace(spec.Scope) != "" && strings.TrimSpace(spec.Scope) != "." {
		return cp, fmt.Errorf("scope requested for child %d but workspace physics is not enabled", spec.Index)
	}

	if err := startQuineChildProcess(ctx, cp, f.QuinePath, spec.Mission, env, f.WorkDir, spec.CaptureRoot, capture, f.ProcessStarted); err != nil {
		return cp, err
	}
	if spec.RecordWorkspaceSession {
		cp.workspaceSession = sessionID
	}

	return cp, nil
}

// spawnChild creates and starts a single child process. If capture is true,
// stdout/stderr are captured into regular files; otherwise they are discarded.
// Before activation it writes a retained seed surface for the new session.
func (f *ForkExecutor) spawnChild(ctx context.Context, relation *forkRelationSurface, toolID string, child ForkChild, index int, capture bool) (*childProcess, error) {
	return f.launchChildSession(ctx, childLaunchSpec{
		ToolID:      toolID,
		Mission:     f.childProcessMission(child),
		SessionID:   newForkSessionID(index),
		ContextMode: childContextInherited,
		World:       child.World,
		Protection:  child.Protection,
		Scope:       child.Scope,
		Relation:    relation,
		Index:       index,
		CaptureRoot: filepath.Join(f.DataDir, "fork-output"),
		RecordWorkspaceSession: child.adoptable(&config.Config{
			WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: f.WorkspaceEnabled},
			ToolGates:       config.ToolGates{ForkWorldEnabled: f.ForkWorldEnabled},
		}),
		ForkChild: &child,
	}, capture)
}

func (f *ForkExecutor) childProcessMission(child ForkChild) string {
	if parentMission := strings.TrimSpace(f.Mission); parentMission != "" {
		return parentMission
	}
	return strings.TrimSpace(child.Intent)
}

func (f *ForkExecutor) spawnAll(ctx context.Context, relation *forkRelationSurface, toolID string, req ForkRequest, capture bool) ([]*childProcess, []string) {
	return startSwarmChildren(ctx, req.Children, capture,
		func(ctx context.Context, child ForkChild, index int, capture bool) (*childProcess, error) {
			return f.spawnChild(ctx, relation, toolID, child, index, capture)
		},
		func(child ForkChild, i int, cp *childProcess) {
			childStructured := f.structuredChild(child, i)
			childStructured.Status = map[bool]string{true: "running", false: "spawned"}[capture]
			f.enrichStructuredChild(&childStructured, child, cp)
			_ = f.writeRelationMember(relation, childStructured)
			_ = f.appendRelationLog(relation, forkRelationLogEntry{
				Event:       "member_started",
				Status:      childStructured.Status,
				MemberIndex: &i,
				SessionID:   childStructured.SessionID,
				PID:         childStructured.PID,
			})
		},
		func(child ForkChild, i int, cp *childProcess, err error) {
			childStructured := f.structuredChild(child, i)
			childStructured.Status = "spawn_failed"
			childStructured.Error = err.Error()
			f.enrichStructuredChild(&childStructured, child, cp)
			_ = f.writeRelationMember(relation, childStructured)
			_ = f.appendRelationLog(relation, forkRelationLogEntry{
				Event:       "member_spawn_failed",
				Status:      "spawn_failed",
				MemberIndex: &i,
				SessionID:   childStructured.SessionID,
				Detail:      err.Error(),
			})
		},
	)
}

func (cp *childProcess) closeCaptureFiles() {
	if cp == nil {
		return
	}
	if cp.stdoutFile != nil {
		_ = cp.stdoutFile.Close()
		cp.stdoutFile = nil
	}
	if cp.stderrFile != nil {
		_ = cp.stderrFile.Close()
		cp.stderrFile = nil
	}
}

func (cp *childProcess) removeCaptureFiles() {
	if cp == nil {
		return
	}
	if cp.stdoutPath != "" {
		_ = os.Remove(cp.stdoutPath)
	}
	if cp.stderrPath != "" {
		_ = os.Remove(cp.stderrPath)
	}
}

func (f *ForkExecutor) childCapturedOutput(cp *childProcess, path string) string {
	if cp == nil || path == "" {
		return ""
	}
	cp.closeCaptureFiles()
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("[capture read error: %v]", err)
	}
	return f.truncate(data)
}

func (f *ForkExecutor) childStdout(cp *childProcess) string {
	if cp == nil {
		return ""
	}
	return f.childCapturedOutput(cp, cp.stdoutPath)
}

func (f *ForkExecutor) childStderr(cp *childProcess) string {
	if cp == nil {
		return ""
	}
	return f.childCapturedOutput(cp, cp.stderrPath)
}

func (f *ForkExecutor) resolveChildSessionID(pid int) string {
	if pid <= 0 {
		return ""
	}
	for attempt := 0; attempt < 200; attempt++ {
		if sessionID, ok := f.childSessionIDFromPIDIndex(pid); ok {
			return sessionID
		}
		if sessionID, ok := f.childSessionIDFromAgentRoots(pid); ok {
			return sessionID
		}
		if attempt < 199 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	return ""
}

func (f *ForkExecutor) childSessionIDFromPIDIndex(pid int) (string, bool) {
	linkPath := filepath.Join(f.DataDir, "pid", strconv.Itoa(pid))
	target, err := os.Readlink(linkPath)
	if err != nil {
		return "", false
	}
	sessionID := strings.TrimSpace(filepath.Base(target))
	if sessionID == "public" {
		sessionID = strings.TrimSpace(filepath.Base(filepath.Dir(target)))
	}
	return sessionID, sessionID != ""
}

func (f *ForkExecutor) childSessionIDFromAgentRoots(pid int) (string, bool) {
	agentDir := filepath.Join(f.DataDir, "agent")
	entries, err := os.ReadDir(agentDir)
	if err != nil {
		return "", false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		statusPath := filepath.Join(agentDir, entry.Name(), "status", "session.json")
		data, err := os.ReadFile(statusPath)
		if err != nil {
			continue
		}
		var status struct {
			SessionID     string `json:"session_id"`
			ParentSession string `json:"parent_session"`
			PID           int    `json:"pid"`
		}
		if err := json.Unmarshal(data, &status); err != nil {
			continue
		}
		if status.PID != pid {
			continue
		}
		if f.SessionID != "" && status.ParentSession != f.SessionID {
			continue
		}
		sessionID := strings.TrimSpace(status.SessionID)
		if sessionID == "" {
			sessionID = entry.Name()
		}
		return sessionID, true
	}
	return "", false
}

func (f *ForkExecutor) resolveChildWorkspace(raw string) (string, error) {
	return resolveChildWorkspace(f.WorkspaceRoot, f.Workspace, raw)
}

func resolveChildWorkspace(root, current, raw string) (string, error) {
	workspace := strings.TrimSpace(raw)
	if workspace == "" || workspace == "." {
		return current, nil
	}

	if filepath.IsAbs(workspace) {
		return "", fmt.Errorf(
			"scope %q must be \".\" or a relative path under the current scope; absolute paths are invalid",
			workspace,
		)
	}
	workspace = filepath.Join(current, workspace)
	workspace = filepath.Clean(workspace)

	resolved, err := canonicalizeRequestedPath(workspace)
	if err != nil {
		return "", err
	}
	root, err = canonicalizeRequestedPath(root)
	if err != nil {
		return "", err
	}
	if !isPathWithin(root, resolved) {
		return "", fmt.Errorf(
			"scope %q must stay within workspace root %q; use \".\" or a relative subpath under the current scope, not an unrelated absolute path",
			resolved,
			root,
		)
	}
	return resolved, nil
}

func (f *ForkExecutor) currentWorkspaceRootAndPath() (string, string) {
	if f.WorkspaceEnabled {
		return f.WorkspaceRoot, f.Workspace
	}
	return f.WorkDir, f.WorkDir
}

func (f *ForkExecutor) workspaceChildBackend() string {
	if f.WorkspaceEnabled && strings.TrimSpace(f.WorkspaceBackend) != "" {
		return f.WorkspaceBackend
	}
	return "overlay"
}

func (f *ForkExecutor) workspaceChildRevisionMode() config.WorkspaceRevisionMode {
	if f.WorkspaceEnabled && f.WorkspaceRevisionMode != "" {
		return f.WorkspaceRevisionMode
	}
	return config.WorkspaceRevisionRestore
}

func filterWorkspacePhysics(env []string) []string {
	filtered := make([]string, 0, len(env))
	for _, entry := range env {
		if strings.HasPrefix(entry, config.EnvWorkspace+"=") || strings.HasPrefix(entry, "QUINE_WORKSPACE_") {
			continue
		}
		if strings.HasPrefix(entry, config.EnvForkWorldEnabled+"=") {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func (f *ForkExecutor) structuredChild(child ForkChild, index int) forkChildStructuredResult {
	result := forkChildStructuredResult{
		Index:  index,
		Intent: child.Intent,
		Scope:  child.Scope,
	}
	if f.ForkWorldEnabled {
		result.World = string(child.World)
		result.Protection = string(child.Protection)
	}
	return result
}

func (f *ForkExecutor) enrichStructuredChild(result *forkChildStructuredResult, child ForkChild, cp *childProcess) {
	if cp == nil {
		return
	}
	if cp.sessionID != "" {
		result.SessionID = cp.sessionID
		result.AgentRoot = f.sessionRoot(cp.sessionID)
		result.PublicRoot = f.sessionPublicRoot(cp.sessionID)
		result.RetainedRoot = f.sessionRetainedDir(cp.sessionID)
		result.SeedRoot = cp.seedRoot
		result.StatusPath = f.sessionStatusPath(cp.sessionID)
		result.ControlPath = f.sessionControlPath(cp.sessionID)
	}
	if cp.cmd != nil && cp.cmd.Process != nil {
		result.PID = cp.cmd.Process.Pid
	}
	if cp.workspaceSession == "" {
		return
	}
	result.WorkspaceSession = cp.workspaceSession
	result.Adoptable = child.adoptable(&config.Config{
		WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: f.WorkspaceEnabled},
		ToolGates:       config.ToolGates{ForkWorldEnabled: f.ForkWorldEnabled},
	})
	revision, err := f.childCurrentWorldRevision(cp.workspaceSession)
	if err != nil || revision == "" {
		return
	}
	result.WorldRevision = revision
	result.WorldHandle = buildWorldHandle(cp.workspaceSession, revision)
}

func (f *ForkExecutor) forkParentWorkspaceResult(children []forkChildStructuredResult) (string, string, error) {
	if f.subjective == nil || !f.subjective.enabled {
		return "", "", nil
	}
	if f.hasSuccessfulHostChild(children) {
		finalized, err := f.subjective.importHostWorkspaceChanges("fork-host", f.TurnID)
		if err != nil {
			return "", "", err
		}
		return finalized.Mutations, formatWorldRevisionCreated(finalized.Revision, !finalized.Changed), nil
	}
	if f.hasObservedHostChild(children) {
		if err := f.subjective.observeHostWorkspaceChanges(); err != nil {
			return "", "", err
		}
	}

	revision, err := f.subjective.loadCurrentWorldRevision()
	if err != nil {
		return "", "", err
	}
	if revision.ID == "" {
		return "", "", nil
	}
	return formatMutations(nil), formatWorldRevisionCreated(revision, true), nil
}

func (f *ForkExecutor) hasObservedHostChild(children []forkChildStructuredResult) bool {
	if !f.ForkWorldEnabled {
		return false
	}
	for _, child := range children {
		if child.World == string(config.WorldHost) && child.Status != "spawn_failed" && child.Status != "no_result" {
			return true
		}
	}
	return false
}

func (f *ForkExecutor) hasSuccessfulHostChild(children []forkChildStructuredResult) bool {
	if !f.ForkWorldEnabled {
		return false
	}
	for _, child := range children {
		if child.World != string(config.WorldHost) || child.Status != "completed" {
			continue
		}
		if child.ExitCode != nil && *child.ExitCode == 0 {
			return true
		}
	}
	return false
}

func (f *ForkExecutor) childCurrentWorldRevision(workspaceSession string) (string, error) {
	if strings.TrimSpace(workspaceSession) == "" {
		return "", nil
	}
	subj := &subjectiveFS{
		enabled:          true,
		dataDir:          f.DataDir,
		workspaceSession: workspaceSession,
		revisionMode:     f.workspaceChildRevisionMode(),
		initialized:      true,
	}
	var seedRevision string
	var ledgerErr error
	var tapeErr error
	for attempt := 0; attempt < 50; attempt++ {
		revision, err := subj.loadCurrentWorldRevision()
		if err == nil && revision.ID != "" && revision.ID != "wr0" {
			return revision.ID, nil
		}
		if err == nil && revision.ID == "wr0" {
			seedRevision = revision.ID
		} else if err != nil {
			ledgerErr = err
		}

		fallback, err := f.childCurrentWorldRevisionFromTapeOnce(workspaceSession)
		if err == nil && fallback != "" {
			return fallback, nil
		}
		if err != nil {
			tapeErr = err
		}
		if attempt < 49 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if seedRevision != "" {
		return seedRevision, nil
	}
	if ledgerErr != nil {
		return "", ledgerErr
	}
	return "", tapeErr
}

func formatChildWorldSurface(child forkChildStructuredResult) string {
	lines := make([]string, 0, 4)
	if child.WorldHandle != "" {
		lines = append(lines, "[WORLD HANDLE] "+child.WorldHandle)
	}
	if child.WorldRevision != "" {
		lines = append(lines, "[CHILD WORLD REVISION] "+child.WorldRevision)
	}
	if child.WorkspaceSession != "" {
		lines = append(lines, "[CHILD WORKSPACE SESSION] "+child.WorkspaceSession)
	}
	if child.Adoptable {
		lines = append(lines, "[ADOPTABLE] true")
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

func (f *ForkExecutor) childCurrentWorldRevisionFromTape(sessionID string) (string, error) {
	if strings.TrimSpace(sessionID) == "" {
		return "", nil
	}
	for attempt := 0; attempt < 50; attempt++ {
		revision, err := f.childCurrentWorldRevisionFromTapeOnce(sessionID)
		if err == nil && revision != "" {
			return revision, nil
		}
		if attempt < 49 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	return "", fmt.Errorf("no child world revision found in tapes for session %s", sessionID)
}

func (f *ForkExecutor) childCurrentWorldRevisionFromTapeOnce(sessionID string) (string, error) {
	tapeDir := childTapeDir(f.DataDir, f.Env, sessionID)
	entries, err := os.ReadDir(tapeDir)
	if err != nil {
		return "", err
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		paths = append(paths, filepath.Join(tapeDir, entry.Name()))
	}
	sort.Strings(paths)
	for i := len(paths) - 1; i >= 0; i-- {
		revision, err := extractCurrentRevisionFromTape(paths[i])
		if err == nil && revision != "" {
			return revision, nil
		}
	}
	return "", fmt.Errorf("no world revision block found in tapes for session %s", sessionID)
}

func childTapeDir(dataDir string, env []string, sessionID string) string {
	retentionRoot := ""
	for _, entry := range env {
		if strings.HasPrefix(entry, config.EnvRetentionDir+"=") {
			retentionRoot = strings.TrimSpace(strings.TrimPrefix(entry, config.EnvRetentionDir+"="))
			break
		}
	}
	if retentionRoot != "" {
		return filepath.Join(retentionRoot, "sessions", sessionID, "tapes")
	}
	return filepath.Join(dataDir, "log", sessionID, "tapes")
}

func extractCurrentRevisionFromTape(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entryType, _ := entry["type"].(string); entryType != "tool_result" {
			continue
		}
		dataField, _ := entry["data"].(map[string]any)
		if dataField == nil {
			continue
		}
		if block := extractWorldRevisionBlock(dataField); block != "" {
			if revision := currentRevisionFromWorldRevisionBlock(block); revision != "" {
				return revision, nil
			}
		}
	}
	return "", fmt.Errorf("no world revision block found in tape %s", path)
}

func extractWorldRevisionBlock(dataField map[string]any) string {
	switch content := dataField["content"].(type) {
	case map[string]any:
		if block, ok := content["world_revision"].(string); ok && strings.TrimSpace(block) != "" {
			return block
		}
	case string:
		if strings.TrimSpace(content) == "" {
			return ""
		}
		if payload, err := tape.UnmarshalToolResultContent(content); err == nil {
			if block, ok := payload["world_revision"].(string); ok && strings.TrimSpace(block) != "" {
				return block
			}
		}
	}
	return ""
}

func currentRevisionFromWorldRevisionBlock(block string) string {
	if block == "" {
		return ""
	}
	if idx := strings.LastIndex(block, "current="); idx >= 0 {
		rest := block[idx+len("current="):]
		fields := strings.Fields(rest)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "wr") {
			return strings.TrimSpace(fields[0])
		}
	}
	if idx := strings.LastIndex(block, "-> "); idx >= 0 {
		rest := block[idx+len("-> "):]
		fields := strings.Fields(rest)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "wr") {
			return strings.TrimSpace(fields[0])
		}
	}
	return ""
}

func canonicalizeRequestedPath(path string) (string, error) {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved, nil
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	parent := filepath.Dir(abs)
	parentResolved, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(parentResolved, filepath.Base(abs)), nil
}

func (f *ForkExecutor) writeChildSeedSurface(sessionID string, relation *forkRelationSurface, projection *childContextProjection, index int) (string, error) {
	seedRoot := filepath.Join(f.sessionRetainedDir(sessionID), "seed")
	if err := os.MkdirAll(seedRoot, 0o755); err != nil {
		return "", fmt.Errorf("create seed root: %w", err)
	}
	contextRoot := filepath.Join(seedRoot, "context")
	if err := f.copyContextProjectionTo(contextRoot, projection); err != nil {
		return "", err
	}
	origin := map[string]any{
		"kind":              "fork",
		"relation_id":       "",
		"relation_root":     "",
		"relation_handle":   "",
		"tool":              "fork",
		"tool_id":           "",
		"initiator_session": f.SessionID,
		"created_at":        time.Now().UTC().Format(time.RFC3339Nano),
		"child_index":       index,
	}
	if projection != nil {
		origin["prior_mission"] = projection.ParentMission
		origin["intent"] = projection.Child.Intent
		origin["world"] = projection.Child.World
		origin["protection"] = projection.Child.Protection
		origin["scope"] = projection.Child.Scope
	}
	if relation != nil {
		origin["relation_id"] = relation.ID
		origin["relation_root"] = relation.Root
		origin["relation_handle"] = relation.Handle
		origin["tool_id"] = relation.ToolID
	}
	if err := writeJSONFile(filepath.Join(seedRoot, "origin.json"), origin); err != nil {
		return "", fmt.Errorf("write seed origin: %w", err)
	}
	return seedRoot, nil
}

func (f *ForkExecutor) copyContextProjectionTo(dst string, projection *childContextProjection) error {
	if strings.TrimSpace(f.ContextRoot) == "" {
		return nil
	}
	contextRoot := f.ContextRoot
	if resolved, err := filepath.EvalSymlinks(contextRoot); err == nil {
		contextRoot = resolved
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("resolve context root: %w", err)
	}
	info, err := os.Stat(contextRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat context root: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("context root %q is not a directory", contextRoot)
	}
	if err := CopyTreePreservingSymlinks(contextRoot, dst); err != nil {
		return fmt.Errorf("copy context tree: %w", err)
	}
	if err := projectBootstrappedCurrentContext(filepath.Join(dst, "state", "current.jsonl"), projection); err != nil {
		return fmt.Errorf("project bootstrapped current context: %w", err)
	}
	if err := writeForkAssignmentPrompt(dst, projection); err != nil {
		return err
	}
	return nil
}

func writeForkAssignmentPrompt(contextRoot string, projection *childContextProjection) error {
	if projection == nil {
		return nil
	}
	parentMission := strings.TrimSpace(projection.ParentMission)
	childIntent := strings.TrimSpace(projection.Child.Intent)
	if parentMission == "" || childIntent == "" || parentMission == childIntent {
		return nil
	}
	promptRoot := filepath.Join(contextRoot, "prompt")
	if err := os.MkdirAll(promptRoot, 0o755); err != nil {
		return fmt.Errorf("create fork assignment prompt root: %w", err)
	}
	var sb strings.Builder
	sb.WriteString("The parent mission remains the active mission for this fork child.\n")
	sb.WriteString("This child owns one isolated lane within that mission:\n\n")
	sb.WriteString(childIntent)
	sb.WriteString("\n\n")
	sb.WriteString("Use the parent mission for task-visible inputs, required paths, output formats, and final contract. ")
	sb.WriteString("Use this fork assignment to choose the lane-specific strategy and success check.\n")
	if projection.Child.Scope != "" {
		sb.WriteString("\nScope: `")
		sb.WriteString(projection.Child.Scope)
		sb.WriteString("`.\n")
	}
	if projection.Child.World != "" {
		sb.WriteString("World: `")
		sb.WriteString(string(projection.Child.World))
		sb.WriteString("`.\n")
	}
	if projection.Child.Protection != "" {
		sb.WriteString("Protection: `")
		sb.WriteString(string(projection.Child.Protection))
		sb.WriteString("`.\n")
	}
	return os.WriteFile(filepath.Join(promptRoot, "45-fork-assignment.md"), []byte(sb.String()), 0o644)
}

// waitChild waits for a child to complete and returns its result.
func (f *ForkExecutor) waitChild(cp *childProcess) childResult {
	return waitChildProcess(cp, f.ProcessEnded)
}

// killProcessGroup sends SIGKILL to a child's entire process group.
func killProcessGroup(cp *childProcess) {
	if cp != nil && cp.cmd != nil && cp.cmd.Process != nil {
		_ = syscall.Kill(-cp.cmd.Process.Pid, syscall.SIGKILL)
	}
}

func childStarted(cp *childProcess) bool {
	return cp != nil && cp.cmd != nil && cp.cmd.Process != nil
}

// executeGatherAll spawns N children concurrently and waits for ALL to complete.
func (f *ForkExecutor) executeGatherAll(toolID string, req ForkRequest) tape.ToolResult {
	n := len(req.Children)

	if err := f.initState(); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeWait,
				Status:    "error",
				Requested: n,
				Errors:    []string{err.Error()},
			}),
			IsError: true,
		}
	}
	relation, err := f.beginRelation(toolID, req)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeWait,
				Status:    "error",
				Requested: n,
				Errors:    []string{fmt.Sprintf("create relation surface: %v", err)},
			}),
			IsError: true,
		}
	}

	ctx, cancel := f.waitContext()
	defer cancel()

	children, spawnErrors := f.spawnAll(ctx, relation, toolID, req, true)

	// If all children failed to spawn, return immediately
	spawned := 0
	for _, cp := range children {
		if childStarted(cp) {
			spawned++
		}
	}
	_ = f.updateRelationStatus(relation, forkRelationStatus{
		Status:           map[bool]string{true: "error", false: "running"}[spawned == 0],
		Mode:             ForkModeWait,
		InitiatorSession: f.SessionID,
		Requested:        n,
		Spawned:          spawned,
		Errors:           append([]string(nil), spawnErrors...),
	})
	if spawned == 0 {
		structuredChildren := make([]forkChildStructuredResult, 0, n)
		for i, child := range req.Children {
			childStructured := f.structuredChild(child, i)
			childStructured.Status = "spawn_failed"
			childStructured.Error = swarmSpawnErrorFor(i, spawnErrors)
			f.enrichStructuredChild(&childStructured, child, children[i])
			structuredChildren = append(structuredChildren, childStructured)
		}
		return f.relationToolResult(toolID, relation, forkStructuredResult{
			Tool:      "fork",
			Mode:      ForkModeWait,
			Status:    "error",
			Requested: n,
			Children:  structuredChildren,
			Errors:    spawnErrors,
		}, true)
	}

	results, completed, killed, killedChildren, timeoutError := waitForSwarmChildren(
		ctx,
		children,
		f.waitChild,
		fmt.Sprintf("fork wait timed out after %s", f.DefaultTimeout),
		func(entry forkRelationLogEntry) error { return f.appendRelationLog(relation, entry) },
	)
	timedOut := timeoutError != ""

	// Aggregate results
	allFailed := true
	var sb strings.Builder
	structuredChildren := make([]forkChildStructuredResult, 0, n)
	if timedOut {
		fmt.Fprintf(&sb, "[FORK SWARM] %d children requested, %d completed before timeout, %d killed\n", n, completed-killed, killed)
	} else {
		fmt.Fprintf(&sb, "[FORK SWARM] %d children completed\n", n)
	}

	for i := range req.Children {
		if !childStarted(children[i]) {
			fmt.Fprintf(&sb, "\n--- CHILD %d [spawn failed]: %q @ scope %q ---\n", i, req.Children[i].Intent, req.Children[i].Scope)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = "spawn_failed"
			childStructured.Error = swarmSpawnErrorFor(i, spawnErrors)
			if childStructured.Error != "" {
				fmt.Fprintf(&sb, "%s\n", childStructured.Error)
			}
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}

		r := results[i]
		if killedChildren[i] || (timedOut && r.child == nil && r.err == nil) {
			fmt.Fprintf(&sb, "\n--- CHILD %d [timeout]: %q @ scope %q ---\n%s\n", i, req.Children[i].Intent, req.Children[i].Scope, timeoutError)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = "timeout"
			childStructured.Error = timeoutError
			if r.child != nil {
				exitCode := r.exitCode
				childStructured.ExitCode = &exitCode
				childStructured.Stdout = f.childStdout(r.child)
				childStructured.Stderr = f.childStderr(r.child)
			}
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}
		if r.err != nil {
			// Execution error (timeout, etc.)
			fmt.Fprintf(&sb, "\n--- CHILD %d [error]: %q @ scope %q ---\n%v\n", i, req.Children[i].Intent, req.Children[i].Scope, r.err)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = "error"
			childStructured.Error = r.err.Error()
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}

		if r.exitCode == 0 {
			allFailed = false
		}

		stdout := f.childStdout(r.child)
		stderr := f.childStderr(r.child)
		exitCode := r.exitCode
		childStructured := f.structuredChild(req.Children[i], i)
		childStructured.Status = "completed"
		childStructured.ExitCode = &exitCode
		childStructured.Stdout = stdout
		childStructured.Stderr = stderr
		f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
		_ = f.writeRelationMember(relation, childStructured)
		_ = f.appendRelationLog(relation, forkRelationLogEntry{
			Event:       "member_completed",
			Status:      "completed",
			MemberIndex: &i,
			SessionID:   childStructured.SessionID,
			PID:         childStructured.PID,
			Detail:      fmt.Sprintf("exit_code=%d", r.exitCode),
		})
		fmt.Fprintf(&sb, "\n--- CHILD %d [exit %d]: %q @ scope %q ---\n[STDOUT]\n%s\n[STDERR]\n%s\n%s",
			i, r.exitCode, req.Children[i].Intent, req.Children[i].Scope, stdout, stderr, formatChildWorldSurface(childStructured))
		structuredChildren = append(structuredChildren, childStructured)
	}

	mutations := ""
	worldRevisionBlock := ""
	if nextMutations, nextWorldRevisionBlock, err := f.forkParentWorkspaceResult(structuredChildren); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeWait,
				Status:    "error",
				Requested: n,
				Spawned:   spawned,
				Children:  structuredChildren,
				Errors:    append(append([]string(nil), spawnErrors...), fmt.Sprintf("finalize parent workspace after fork: %v", err)),
			}),
			IsError: true,
		}
	} else {
		mutations = nextMutations
		worldRevisionBlock = nextWorldRevisionBlock
	}
	if !f.FSMutationTelemetryEnabled {
		mutations = ""
	}
	errors := append([]string(nil), spawnErrors...)
	status := "completed"
	if timedOut {
		status = "timeout"
		errors = append(errors, timeoutError)
	}
	succeeded := 0
	for _, child := range structuredChildren {
		if child.ExitCode != nil && *child.ExitCode == 0 && child.Status == "completed" {
			succeeded++
		}
	}
	_ = f.updateRelationStatus(relation, forkRelationStatus{
		Status:           status,
		Mode:             ForkModeWait,
		InitiatorSession: f.SessionID,
		Requested:        n,
		Spawned:          spawned,
		Completed:        completed,
		Succeeded:        succeeded,
		Killed:           killed,
		Errors:           errors,
	})

	return f.relationToolResult(toolID, relation, forkStructuredResult{
		Tool:          "fork",
		Mode:          ForkModeWait,
		Status:        status,
		Requested:     n,
		Spawned:       spawned,
		Succeeded:     succeeded,
		Killed:        killed,
		FSMutations:   mutations,
		WorldRevision: worldRevisionBlock,
		Children:      structuredChildren,
		Errors:        errors,
	}, allFailed)
}

// executeRace spawns N children concurrently. The first child to exit with
// code 0 wins; all remaining children are killed. If all children fail,
// all results are returned.
func (f *ForkExecutor) executeRace(toolID string, req ForkRequest) tape.ToolResult {
	n := len(req.Children)

	if err := f.initState(); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeRace,
				Status:    "error",
				Requested: n,
				Errors:    []string{err.Error()},
			}),
			IsError: true,
		}
	}
	relation, err := f.beginRelation(toolID, req)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeRace,
				Status:    "error",
				Requested: n,
				Errors:    []string{fmt.Sprintf("create relation surface: %v", err)},
			}),
			IsError: true,
		}
	}

	ctx, cancel := f.waitContext()
	defer cancel()

	children, spawnErrors := f.spawnAll(ctx, relation, toolID, req, true)

	// Count successfully spawned
	spawned := 0
	for _, cp := range children {
		if childStarted(cp) {
			spawned++
		}
	}
	_ = f.updateRelationStatus(relation, forkRelationStatus{
		Status:           map[bool]string{true: "error", false: "running"}[spawned == 0],
		Mode:             ForkModeRace,
		InitiatorSession: f.SessionID,
		Requested:        n,
		Spawned:          spawned,
		Errors:           append([]string(nil), spawnErrors...),
	})
	if spawned == 0 {
		structuredChildren := make([]forkChildStructuredResult, 0, n)
		for i, child := range req.Children {
			childStructured := f.structuredChild(child, i)
			childStructured.Status = "spawn_failed"
			childStructured.Error = swarmSpawnErrorFor(i, spawnErrors)
			f.enrichStructuredChild(&childStructured, child, children[i])
			structuredChildren = append(structuredChildren, childStructured)
		}
		return f.relationToolResult(toolID, relation, forkStructuredResult{
			Tool:      "fork",
			Mode:      ForkModeRace,
			Status:    "error",
			Requested: n,
			Children:  structuredChildren,
			Errors:    spawnErrors,
		}, true)
	}

	completed, killed, killedChildren, winnerIndex, timeoutError := raceSwarmChildren(
		ctx,
		children,
		f.waitChild,
		fmt.Sprintf("fork race timed out after %s", f.DefaultTimeout),
		func(entry forkRelationLogEntry) error { return f.appendRelationLog(relation, entry) },
	)
	timedOut := timeoutError != ""

	// Build result
	if winnerIndex != nil {
		winner := *winnerIndex
		r := completed[winner]
		cp := children[winner]
		stdout := f.childStdout(cp)
		stderr := f.childStderr(cp)

		succeeded := 0
		for _, cr := range completed {
			if cr.err == nil && cr.exitCode == 0 {
				succeeded++
			}
		}

		winnerExit := r.exitCode
		winnerStructured := f.structuredChild(req.Children[winner], winner)
		winnerStructured.Status = "completed"
		winnerStructured.ExitCode = &winnerExit
		winnerStructured.Stdout = stdout
		winnerStructured.Stderr = stderr
		f.enrichStructuredChild(&winnerStructured, req.Children[winner], children[winner])
		_ = f.writeRelationMember(relation, winnerStructured)
		_ = f.appendRelationLog(relation, forkRelationLogEntry{
			Event:       "winner_selected",
			Status:      "completed",
			MemberIndex: &winner,
			SessionID:   winnerStructured.SessionID,
			PID:         winnerStructured.PID,
			Detail:      fmt.Sprintf("exit_code=%d", winnerExit),
		})
		mutations := ""
		worldRevisionBlock := ""
		// The adopt decision must be separated from world-handle presence. If
		// the caller asked to adopt the winner and adoption can actually move
		// the world (restore mode tracks revisions), but the winner's world
		// handle is empty because childCurrentWorldRevision could not resolve
		// the revision (the error is swallowed in enrichStructuredChild), do NOT
		// fall through to the non-adopting branch and report "completed" — that
		// silently discards the winner's filesystem work. Fail loudly instead.
		// Gated on canRestoreWorld(): in revision_mode=none an empty handle is
		// expected and adoption is a legitimate no-op, so the guard must not fire.
		if req.AdoptWinner && f.subjective != nil && f.subjective.canRestoreWorld() && winnerStructured.WorldHandle == "" {
			return tape.ToolResult{
				ToolID: toolID,
				Content: tape.MarshalToolResultContent(forkStructuredResult{
					Tool:      "fork",
					Mode:      ForkModeRace,
					Status:    "error",
					Requested: n,
					Spawned:   spawned,
					Succeeded: succeeded,
					Killed:    killed,
					Winner:    &winnerStructured,
					Errors:    append(append([]string(nil), spawnErrors...), "adopt winner: could not resolve the winner's world handle (world revision unavailable); winner work was not adopted"),
				}),
				IsError: true,
			}
		}
		if req.AdoptWinner && winnerStructured.WorldHandle != "" && f.subjective != nil && f.subjective.enabled {
			previous, current, err := f.subjective.switchWorld(winnerStructured.WorldHandle)
			if err != nil {
				return tape.ToolResult{
					ToolID: toolID,
					Content: tape.MarshalToolResultContent(forkStructuredResult{
						Tool:      "fork",
						Mode:      ForkModeRace,
						Status:    "error",
						Requested: n,
						Spawned:   spawned,
						Succeeded: succeeded,
						Killed:    killed,
						Winner:    &winnerStructured,
						Errors:    append(append([]string(nil), spawnErrors...), fmt.Sprintf("adopt winner world: %v", err)),
					}),
					IsError: true,
				}
			}
			mutations, err = f.subjective.restoreMutationBlock(previous, current)
			if err != nil {
				return tape.ToolResult{
					ToolID: toolID,
					Content: tape.MarshalToolResultContent(forkStructuredResult{
						Tool:      "fork",
						Mode:      ForkModeRace,
						Status:    "error",
						Requested: n,
						Spawned:   spawned,
						Succeeded: succeeded,
						Killed:    killed,
						Winner:    &winnerStructured,
						Errors:    append(append([]string(nil), spawnErrors...), fmt.Sprintf("winner adoption diff: %v", err)),
					}),
					IsError: true,
				}
			}
			worldRevisionBlock = formatWorldRevisionTransition(previous, current)
		} else if nextMutations, nextWorldRevisionBlock, err := f.forkParentWorkspaceResult([]forkChildStructuredResult{winnerStructured}); err != nil {
			return tape.ToolResult{
				ToolID: toolID,
				Content: tape.MarshalToolResultContent(forkStructuredResult{
					Tool:      "fork",
					Mode:      ForkModeRace,
					Status:    "error",
					Requested: n,
					Spawned:   spawned,
					Succeeded: succeeded,
					Killed:    killed,
					Winner:    &winnerStructured,
					Errors:    append(append([]string(nil), spawnErrors...), fmt.Sprintf("finalize parent workspace after fork: %v", err)),
				}),
				IsError: true,
			}
		} else {
			mutations = nextMutations
			worldRevisionBlock = nextWorldRevisionBlock
		}
		if !f.FSMutationTelemetryEnabled {
			mutations = ""
		}

		return f.relationToolResult(toolID, relation, forkStructuredResult{
			Tool:          "fork",
			Mode:          ForkModeRace,
			Status:        "completed",
			Requested:     n,
			Spawned:       spawned,
			Succeeded:     succeeded,
			Killed:        killed,
			FSMutations:   mutations,
			WorldRevision: worldRevisionBlock,
			Winner:        &winnerStructured,
			Errors:        spawnErrors,
		}, false)
	}

	// All children failed or the race timed out before a winner.
	var sb strings.Builder
	structuredChildren := make([]forkChildStructuredResult, 0, n)
	if timedOut {
		fmt.Fprintf(&sb, "[FORK RACE] timed out after %s; %d of %d children returned before timeout\n", f.DefaultTimeout, len(completed)-killed, spawned)
	} else {
		fmt.Fprintf(&sb, "[FORK RACE] all %d children failed\n", n)
	}

	for i := range req.Children {
		if !childStarted(children[i]) {
			fmt.Fprintf(&sb, "\n--- CHILD %d [spawn failed]: %q @ scope %q ---\n", i, req.Children[i].Intent, req.Children[i].Scope)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = "spawn_failed"
			childStructured.Error = swarmSpawnErrorFor(i, spawnErrors)
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}

		r, ok := completed[i]
		if !ok {
			status := "no_result"
			errMsg := ""
			if killedChildren[i] || timedOut {
				status = "timeout"
				errMsg = timeoutError
			}
			fmt.Fprintf(&sb, "\n--- CHILD %d [%s]: %q @ scope %q ---\n%s\n", i, status, req.Children[i].Intent, req.Children[i].Scope, errMsg)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = status
			childStructured.Error = errMsg
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}

		if killedChildren[i] {
			fmt.Fprintf(&sb, "\n--- CHILD %d [timeout]: %q @ scope %q ---\n%s\n", i, req.Children[i].Intent, req.Children[i].Scope, timeoutError)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = "timeout"
			childStructured.Error = timeoutError
			exitCode := r.exitCode
			childStructured.ExitCode = &exitCode
			childStructured.Stdout = f.childStdout(children[i])
			childStructured.Stderr = f.childStderr(children[i])
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}

		if r.err != nil {
			fmt.Fprintf(&sb, "\n--- CHILD %d [error]: %q @ scope %q ---\n%v\n", i, req.Children[i].Intent, req.Children[i].Scope, r.err)
			childStructured := f.structuredChild(req.Children[i], i)
			childStructured.Status = "error"
			childStructured.Error = r.err.Error()
			f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
			_ = f.writeRelationMember(relation, childStructured)
			structuredChildren = append(structuredChildren, childStructured)
			continue
		}

		stdout := f.childStdout(children[i])
		stderr := f.childStderr(children[i])
		exitCode := r.exitCode
		childStructured := f.structuredChild(req.Children[i], i)
		childStructured.Status = "completed"
		childStructured.ExitCode = &exitCode
		childStructured.Stdout = stdout
		childStructured.Stderr = stderr
		f.enrichStructuredChild(&childStructured, req.Children[i], children[i])
		_ = f.writeRelationMember(relation, childStructured)
		fmt.Fprintf(&sb, "\n--- CHILD %d [exit %d]: %q @ scope %q ---\n[STDOUT]\n%s\n[STDERR]\n%s\n%s",
			i, r.exitCode, req.Children[i].Intent, req.Children[i].Scope, stdout, stderr, formatChildWorldSurface(childStructured))
		structuredChildren = append(structuredChildren, childStructured)
	}

	mutations := ""
	worldRevisionBlock := ""
	if nextMutations, nextWorldRevisionBlock, err := f.forkParentWorkspaceResult(structuredChildren); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeRace,
				Status:    "error",
				Requested: n,
				Spawned:   spawned,
				Children:  structuredChildren,
				Errors:    append(append([]string(nil), spawnErrors...), fmt.Sprintf("finalize parent workspace after fork: %v", err)),
			}),
			IsError: true,
		}
	} else {
		mutations = nextMutations
		worldRevisionBlock = nextWorldRevisionBlock
	}
	if !f.FSMutationTelemetryEnabled {
		mutations = ""
	}

	errors := append([]string(nil), spawnErrors...)
	status := "completed"
	if timedOut {
		status = "timeout"
		errors = append(errors, timeoutError)
	}
	_ = f.updateRelationStatus(relation, forkRelationStatus{
		Status:           status,
		Mode:             ForkModeRace,
		InitiatorSession: f.SessionID,
		Requested:        n,
		Spawned:          spawned,
		Completed:        len(completed),
		Killed:           killed,
		Errors:           errors,
	})

	return f.relationToolResult(toolID, relation, forkStructuredResult{
		Tool:          "fork",
		Mode:          ForkModeRace,
		Status:        status,
		Requested:     n,
		Spawned:       spawned,
		Killed:        killed,
		FSMutations:   mutations,
		WorldRevision: worldRevisionBlock,
		Children:      structuredChildren,
		Errors:        errors,
	}, true)
}

// executeFireAndForget spawns N children and returns their PIDs immediately.
func (f *ForkExecutor) executeFireAndForget(toolID string, req ForkRequest) tape.ToolResult {
	n := len(req.Children)
	if err := f.initState(); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeForget,
				Status:    "error",
				Requested: n,
				Errors:    []string{err.Error()},
			}),
			IsError: true,
		}
	}
	relation, err := f.beginRelation(toolID, req)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(forkStructuredResult{
				Tool:      "fork",
				Mode:      ForkModeForget,
				Status:    "error",
				Requested: n,
				Errors:    []string{fmt.Sprintf("create relation surface: %v", err)},
			}),
			IsError: true,
		}
	}

	children, spawnErrors := f.spawnAll(context.Background(), relation, toolID, req, false)
	structuredChildren := make([]forkChildStructuredResult, 0, n)
	for i, child := range req.Children {
		childStructured := f.structuredChild(child, i)
		if childStarted(children[i]) {
			childStructured.Status = "spawned"
		} else {
			childStructured.Status = "spawn_failed"
			childStructured.Error = swarmSpawnErrorFor(i, spawnErrors)
		}
		f.enrichStructuredChild(&childStructured, child, children[i])
		_ = f.writeRelationMember(relation, childStructured)
		structuredChildren = append(structuredChildren, childStructured)
	}
	spawned := countStarted(children)
	if spawned == 0 {
		return f.relationToolResult(toolID, relation, forkStructuredResult{
			Tool:      "fork",
			Mode:      ForkModeForget,
			Status:    "error",
			Requested: n,
			Children:  structuredChildren,
			Errors:    spawnErrors,
		}, true)
	}

	// Build result with PIDs
	var sb strings.Builder
	fmt.Fprintf(&sb, "[FORK OK] %d children spawned\n", spawned)
	for _, cp := range children {
		if !childStarted(cp) {
			continue
		}
		childStructured := f.structuredChild(req.Children[cp.index], cp.index)
		childStructured.Status = "spawned"
		f.enrichStructuredChild(&childStructured, req.Children[cp.index], cp)
		_ = f.writeRelationMember(relation, childStructured)
		fmt.Fprintf(&sb, "  child %d: PID %d — %q\n%s", cp.index, cp.cmd.Process.Pid, cp.intent, formatChildWorldSurface(childStructured))
	}
	if len(spawnErrors) > 0 {
		fmt.Fprintf(&sb, "\nFailed to spawn %d children:\n", len(spawnErrors))
		for _, e := range spawnErrors {
			fmt.Fprintf(&sb, "  %s\n", e)
		}
	}
	fmt.Fprintf(&sb, "\nChildren are running independently.\nChild process surfaces will appear under the runtime root: %s", f.DataDir)

	return f.relationToolResult(toolID, relation, forkStructuredResult{
		Tool:      "fork",
		Mode:      ForkModeForget,
		Status:    "spawned",
		Requested: n,
		Spawned:   spawned,
		Children:  structuredChildren,
		Errors:    spawnErrors,
	}, false)
}

// copyContextBootstrap stages the parent's current context tree into a temp
// directory that the child imports before entering its turn loop.
func (f *ForkExecutor) copyContextBootstrap(projection *childContextProjection) (string, error) {
	if strings.TrimSpace(f.ContextRoot) == "" {
		return "", nil
	}
	contextRoot := f.ContextRoot
	if resolved, err := filepath.EvalSymlinks(contextRoot); err == nil {
		contextRoot = resolved
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("resolve context root: %w", err)
	}
	info, err := os.Stat(contextRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("stat context root: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("context root %q is not a directory", contextRoot)
	}
	if err := os.MkdirAll(f.DataDir, 0o755); err != nil {
		return "", fmt.Errorf("create data dir for child context bootstrap: %w", err)
	}
	tmpDir, err := os.MkdirTemp(f.DataDir, "fork-context-*")
	if err != nil {
		return "", fmt.Errorf("create temp context bootstrap: %w", err)
	}
	if err := f.copyContextProjectionTo(tmpDir, projection); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}
	return tmpDir, nil
}

func projectBootstrappedCurrentContext(path string, projection *childContextProjection) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := preserveBootstrappedCurrentRaw(path, data); err != nil {
		return err
	}
	projected, changed, err := projectChildContextCurrent(data, projection)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return os.WriteFile(path, projected, 0o644)
}

func preserveBootstrappedCurrentRaw(path string, data []byte) error {
	rawDir := filepath.Join(filepath.Dir(path), "bootstrap")
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return fmt.Errorf("create bootstrapped current raw dir: %w", err)
	}
	rawPath := filepath.Join(rawDir, "current.parent.raw.jsonl")
	if err := os.WriteFile(rawPath, data, 0o644); err != nil {
		return fmt.Errorf("write bootstrapped current raw snapshot: %w", err)
	}
	return nil
}

func projectChildContextCurrent(data []byte, projection *childContextProjection) ([]byte, bool, error) {
	if projection == nil || strings.TrimSpace(projection.ForkToolID) == "" {
		return sanitizeCurrentContext(data)
	}
	lines := bytes.Split(data, []byte("\n"))
	out := make([][]byte, 0, len(lines)+4)
	var (
		batch   *pendingToolBatch
		changed bool
	)

	for i, rawLine := range lines {
		line := bytes.TrimSpace(rawLine)
		if len(line) == 0 {
			continue
		}

		var entry tape.TapeEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, false, fmt.Errorf("decode context entry %d: %w", i, err)
		}

		switch entry.Type {
		case "message":
			var msg tape.Message
			if err := json.Unmarshal(entry.Data, &msg); err != nil {
				return nil, false, fmt.Errorf("decode context message %d: %w", i, err)
			}
			if batch != nil && batch.hasPending() {
				projectedLines, err := batch.syntheticLines(projection)
				if err != nil {
					return nil, false, err
				}
				out = append(out, projectedLines...)
				changed = true
				batch = nil
			}
			if msg.Role == tape.RoleAssistant && len(msg.ToolCalls) > 0 {
				batch = newPendingToolBatch(msg)
				msg.ToolCalls = nil
				msg.ReasoningContent = ""
				msg.ReasoningItems = nil
				if strings.TrimSpace(msg.Content) == "" && !tape.HasStructuredContent(msg.StructuredContent) && strings.TrimSpace(msg.ReasoningContent) == "" && len(msg.ReasoningItems) == 0 {
					changed = true
					continue
				}
				entry = tape.MessageEntry(msg)
				encodedLine, err := json.Marshal(entry)
				if err != nil {
					return nil, false, fmt.Errorf("re-encode sanitized assistant message %d: %w", i, err)
				}
				line = encodedLine
				changed = true
			}
			out = append(out, bytes.TrimRight(line, "\r"))
		case "tool_result":
			changed = true
			var result tape.ToolResult
			if err := json.Unmarshal(entry.Data, &result); err != nil {
				return nil, false, fmt.Errorf("decode context tool result %d: %w", i, err)
			}
			if batch == nil {
				continue
			}
			batch.resolve(result.ToolID)
			if !batch.hasPending() {
				batch = nil
			}
		default:
			if batch != nil && batch.hasPending() {
				projectedLines, err := batch.syntheticLines(projection)
				if err != nil {
					return nil, false, err
				}
				out = append(out, projectedLines...)
				changed = true
				batch = nil
			}
			out = append(out, bytes.TrimRight(line, "\r"))
		}
	}

	if batch != nil && batch.hasPending() {
		projectedLines, err := batch.syntheticLines(projection)
		if err != nil {
			return nil, false, err
		}
		out = append(out, projectedLines...)
		changed = true
	}
	if !changed {
		return data, false, nil
	}
	return joinContextLines(out), true, nil
}

func sanitizeCurrentContext(data []byte) ([]byte, bool, error) {
	lines := bytes.Split(data, []byte("\n"))
	out := make([][]byte, 0, len(lines))
	var changed bool

	for i, rawLine := range lines {
		line := bytes.TrimSpace(rawLine)
		if len(line) == 0 {
			continue
		}

		var entry tape.TapeEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, false, fmt.Errorf("decode context entry %d: %w", i, err)
		}

		switch entry.Type {
		case "message":
			var msg tape.Message
			if err := json.Unmarshal(entry.Data, &msg); err != nil {
				return nil, false, fmt.Errorf("decode context message %d: %w", i, err)
			}
			if msg.Role == tape.RoleAssistant && len(msg.ToolCalls) > 0 {
				msg.ToolCalls = nil
				if strings.TrimSpace(msg.Content) == "" && !tape.HasStructuredContent(msg.StructuredContent) && strings.TrimSpace(msg.ReasoningContent) == "" && len(msg.ReasoningItems) == 0 {
					changed = true
					continue
				}
				entry = tape.MessageEntry(msg)
				encodedLine, err := json.Marshal(entry)
				if err != nil {
					return nil, false, fmt.Errorf("re-encode sanitized assistant message %d: %w", i, err)
				}
				line = encodedLine
				changed = true
			}
			out = append(out, bytes.TrimRight(line, "\r"))
		case "tool_result":
			changed = true
			continue
		default:
			out = append(out, bytes.TrimRight(line, "\r"))
		}
	}
	if !changed {
		return data, false, nil
	}
	return joinContextLines(out), true, nil
}

type pendingToolBatch struct {
	source  tape.Message
	calls   []tape.ToolCall
	pending map[string]tape.ToolCall
}

func newPendingToolBatch(source tape.Message) *pendingToolBatch {
	calls := source.ToolCalls
	pending := make(map[string]tape.ToolCall, len(calls))
	for _, tc := range calls {
		if strings.TrimSpace(tc.ID) == "" {
			continue
		}
		pending[tc.ID] = tc
	}
	if len(pending) == 0 {
		return nil
	}
	return &pendingToolBatch{
		source:  source,
		calls:   append([]tape.ToolCall(nil), calls...),
		pending: pending,
	}
}

func (b *pendingToolBatch) hasPending() bool {
	return b != nil && len(b.pending) > 0
}

func (b *pendingToolBatch) resolve(toolID string) {
	if b == nil {
		return
	}
	delete(b.pending, strings.TrimSpace(toolID))
}

func (b *pendingToolBatch) syntheticLines(projection *childContextProjection) ([][]byte, error) {
	if b == nil || !b.hasPending() {
		return nil, nil
	}
	lines := make([][]byte, 0, len(b.pending)+1)
	pendingCalls := make([]tape.ToolCall, 0, len(b.pending))
	for _, tc := range b.calls {
		if strings.TrimSpace(tc.ID) == "" {
			continue
		}
		if _, ok := b.pending[tc.ID]; !ok {
			continue
		}
		pendingCalls = append(pendingCalls, tc)
	}
	if len(pendingCalls) > 0 {
		entryBytes, err := json.Marshal(tape.MessageEntry(tape.SyntheticAssistantToolBatch(b.source, pendingCalls)))
		if err != nil {
			return nil, fmt.Errorf("marshal synthetic projected tool batch: %w", err)
		}
		lines = append(lines, entryBytes)
	}
	for _, tc := range b.calls {
		if strings.TrimSpace(tc.ID) == "" {
			continue
		}
		if _, ok := b.pending[tc.ID]; !ok {
			continue
		}
		result := projectedToolResultForCall(tc, projection)
		entryBytes, err := json.Marshal(tape.ToolResultEntry(result))
		if err != nil {
			return nil, fmt.Errorf("marshal synthetic projected tool result for %s: %w", tc.ID, err)
		}
		lines = append(lines, entryBytes)
	}
	return lines, nil
}

func projectedToolResultForCall(call tape.ToolCall, projection *childContextProjection) tape.ToolResult {
	payload := projectedContextToolResult{
		Tool: strings.TrimSpace(call.Name),
	}
	if strings.TrimSpace(call.ID) == projection.ForkToolID {
		payload.Tool = "fork"
		payload.Status = "child_bootstrap"
		if strings.TrimSpace(projection.ParentMission) != "" {
			payload.ParentMissionRef = "context/prompt/40-mission.md"
			payload.CurrentMissionRef = "context/prompt/40-mission.md"
		}
		if strings.TrimSpace(projection.Child.Intent) != "" {
			payload.ChildAssignment = strings.TrimSpace(projection.Child.Intent)
			payload.AssignmentRef = "context/prompt/45-fork-assignment.md"
		}
		payload.Note = "You are now a fork child. The parent mission remains the active mission; follow child_assignment as this lane's focus, not as a replacement for the task contract."
	} else {
		payload.Status = "projected_away"
		payload.Reason = "outside_child_thread"
	}
	return tape.ToolResult{
		ToolID:  call.ID,
		Content: tape.MarshalToolResultContent(payload),
		IsError: false,
	}
}

func joinContextLines(lines [][]byte) []byte {
	kept := make([][]byte, 0, len(lines))
	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		kept = append(kept, bytes.TrimRight(line, "\r"))
	}
	if len(kept) == 0 {
		return nil
	}
	out := bytes.Join(kept, []byte("\n"))
	return append(out, '\n')
}

// truncate returns the string representation of data, truncating if needed.
func (f *ForkExecutor) truncate(data []byte) string {
	if len(data) <= f.MaxOutput {
		return string(data)
	}
	total := len(data)
	truncated := string(data[:f.MaxOutput])
	return truncated + fmt.Sprintf("\n...[Output Truncated, %d bytes total]", total)
}

// ForkResultEntry returns a TapeEntry for a fork tool result.
func ForkResultEntry(sessionID string, childPIDs []int, mode string) tape.TapeEntry {
	data, _ := json.Marshal(map[string]any{
		"child_pids": childPIDs,
		"mode":       mode,
	})
	return tape.TapeEntry{Type: "fork", Data: data}
}
