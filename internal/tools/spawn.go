package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

const (
	SpawnModeWait   = ForkModeWait
	SpawnModeRace   = ForkModeRace
	SpawnModeForget = ForkModeForget
)

type SpawnRequest struct {
	Children []SpawnChild
	Mode     string
}

type SpawnChild struct {
	Mission    string
	World      config.WorldKind
	Protection config.ProtectionMode
	Scope      string
}

type SpawnExecutor struct {
	*ForkExecutor
}

type spawnChildStructuredResult struct {
	Index            int    `json:"index"`
	Mission          string `json:"mission"`
	Status           string `json:"status"`
	World            string `json:"world,omitempty"`
	Protection       string `json:"protection,omitempty"`
	Scope            string `json:"scope,omitempty"`
	SessionID        string `json:"session_id,omitempty"`
	AgentRoot        string `json:"agent_root,omitempty"`
	PublicRoot       string `json:"public_root,omitempty"`
	RetainedRoot     string `json:"retained_root,omitempty"`
	StatusPath       string `json:"status_path,omitempty"`
	ControlPath      string `json:"control_path,omitempty"`
	ExitCode         *int   `json:"exit_code,omitempty"`
	PID              int    `json:"pid,omitempty"`
	Stdout           string `json:"stdout,omitempty"`
	Stderr           string `json:"stderr,omitempty"`
	WorldHandle      string `json:"world_handle,omitempty"`
	WorldRevision    string `json:"world_revision,omitempty"`
	WorkspaceSession string `json:"workspace_session,omitempty"`
	Error            string `json:"error,omitempty"`
}

type spawnStructuredResult struct {
	Tool           string                       `json:"tool"`
	Mode           string                       `json:"mode"`
	Status         string                       `json:"status"`
	RelationID     string                       `json:"relation_id,omitempty"`
	RelationRoot   string                       `json:"relation_root,omitempty"`
	RelationHandle string                       `json:"relation_handle,omitempty"`
	Requested      int                          `json:"requested"`
	Spawned        int                          `json:"spawned"`
	Completed      int                          `json:"completed"`
	Succeeded      int                          `json:"succeeded"`
	Killed         int                          `json:"killed"`
	Winner         *spawnChildStructuredResult  `json:"winner,omitempty"`
	Children       []spawnChildStructuredResult `json:"children,omitempty"`
	Errors         []string                     `json:"errors,omitempty"`
}

func NewSpawnExecutor(cfg *config.Config, childEnv []string) *SpawnExecutor {
	return &SpawnExecutor{ForkExecutor: NewForkExecutor(cfg, childEnv)}
}

func ParseSpawnArgs(args map[string]any, cfg *config.Config) (SpawnRequest, error) {
	children, err := parseSpawnChildren(args, cfg)
	if err != nil {
		return SpawnRequest{}, err
	}
	req := SpawnRequest{
		Children: children,
		Mode:     SpawnModeWait,
	}
	if v, ok := args["mode"]; ok {
		s, ok := v.(string)
		if !ok {
			return SpawnRequest{}, fmt.Errorf("mode must be a string, got %T", v)
		}
		switch strings.TrimSpace(s) {
		case SpawnModeWait, SpawnModeRace, SpawnModeForget:
			req.Mode = strings.TrimSpace(s)
		default:
			return SpawnRequest{}, fmt.Errorf("mode must be one of wait, race, forget; got %q", s)
		}
	}
	return req, nil
}

func parseSpawnChildren(args map[string]any, cfg *config.Config) ([]SpawnChild, error) {
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
	children := make([]SpawnChild, 0, len(rawSlice))
	for i, v := range rawSlice {
		obj, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("children[%d] must be an object, got %T", i, v)
		}
		mission, ok := obj["mission"].(string)
		if !ok {
			return nil, fmt.Errorf("children[%d].mission must be a string", i)
		}
		mission = strings.TrimSpace(mission)
		if mission == "" {
			return nil, fmt.Errorf("children[%d].mission cannot be empty", i)
		}
		world, protection, scope, err := parseChildWorldScopeProtection(obj, i, cfg, false)
		if err != nil {
			return nil, err
		}
		children = append(children, SpawnChild{
			Mission:    mission,
			World:      world,
			Protection: protection,
			Scope:      scope,
		})
	}
	return children, nil
}

func newSpawnSessionID(index int) string {
	return fmt.Sprintf("sess_spawn_%d_%d_%d", time.Now().UnixNano(), os.Getpid(), index)
}

func (s *SpawnExecutor) Execute(toolID string, req SpawnRequest) tape.ToolResult {
	switch req.Mode {
	case SpawnModeRace:
		return s.executeRace(toolID, req)
	case SpawnModeForget:
		return s.executeForget(toolID, req)
	default:
		return s.executeWait(toolID, req)
	}
}

func (s *SpawnExecutor) beginRelation(toolID string, req SpawnRequest) (*forkRelationSurface, error) {
	children := make([]spawnChildStructuredResult, 0, len(req.Children))
	for i, child := range req.Children {
		children = append(children, s.structuredChild(child, i))
	}
	return beginSwarmRelation("spawn", toolID, req.Mode, s.SessionID, s.sessionRetainedDir, len(req.Children), children)
}

func (s *SpawnExecutor) updateRelationStatus(relation *forkRelationSurface, status forkRelationStatus) error {
	return updateSwarmRelationStatus(relation, "spawn", status)
}

func (s *SpawnExecutor) writeRelationMember(relation *forkRelationSurface, child spawnChildStructuredResult) error {
	return writeSwarmRelationMember(relation, child.Index, child)
}

func (s *SpawnExecutor) relationToolResult(toolID string, relation *forkRelationSurface, payload spawnStructuredResult, isError bool) tape.ToolResult {
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
		if err := s.updateRelationStatus(relation, status); err != nil {
			payload.Status = "error"
			payload.Errors = append(payload.Errors, fmt.Sprintf("write relation status: %v", err))
			isError = true
		}
		_ = s.appendRelationLog(relation, forkRelationLogEntry{
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

func (s *SpawnExecutor) structuredChild(child SpawnChild, index int) spawnChildStructuredResult {
	result := spawnChildStructuredResult{
		Index:   index,
		Mission: child.Mission,
		Scope:   child.Scope,
	}
	if s.ForkWorldEnabled {
		result.World = string(child.World)
		result.Protection = string(child.Protection)
	}
	return result
}

func (s *SpawnExecutor) enrichStructuredChild(result *spawnChildStructuredResult, cp *childProcess) {
	if result == nil || cp == nil {
		return
	}
	if cp.sessionID != "" {
		result.SessionID = cp.sessionID
		result.AgentRoot = s.sessionRoot(cp.sessionID)
		result.PublicRoot = s.sessionPublicRoot(cp.sessionID)
		result.RetainedRoot = s.sessionRetainedDir(cp.sessionID)
		result.StatusPath = s.sessionStatusPath(cp.sessionID)
		result.ControlPath = s.sessionControlPath(cp.sessionID)
	}
	if cp.cmd != nil && cp.cmd.Process != nil {
		result.PID = cp.cmd.Process.Pid
	}
	if cp.workspaceSession == "" {
		return
	}
	result.WorkspaceSession = cp.workspaceSession
	revision, err := s.childCurrentWorldRevision(cp.workspaceSession)
	if err != nil || revision == "" {
		return
	}
	result.WorldRevision = revision
	result.WorldHandle = buildWorldHandle(cp.workspaceSession, revision)
}

func (s *SpawnExecutor) spawnChild(ctx context.Context, child SpawnChild, index int, capture bool) (*childProcess, error) {
	return s.launchChildSession(ctx, childLaunchSpec{
		Mission:                child.Mission,
		SessionID:              newSpawnSessionID(index),
		ContextMode:            childContextFresh,
		World:                  child.World,
		Protection:             child.Protection,
		Scope:                  child.Scope,
		Index:                  index,
		CaptureRoot:            filepath.Join(s.DataDir, "spawn-output"),
		RecordWorkspaceSession: spawnChildWorldSurfaceEnabled(s, child),
	}, capture)
}

func spawnChildWorldSurfaceEnabled(s *SpawnExecutor, child SpawnChild) bool {
	if s == nil || !s.WorkspaceEnabled {
		return false
	}
	if s.ForkWorldEnabled {
		return child.World == config.WorldSubjective && child.Protection == config.ProtectionTransactional && s.workspaceChildBackend() == "overlay"
	}
	return s.workspaceChildBackend() == "overlay"
}

func (s *SpawnExecutor) spawnAll(ctx context.Context, relation *forkRelationSurface, req SpawnRequest, capture bool) ([]*childProcess, []string) {
	return startSwarmChildren(ctx, req.Children, capture,
		func(ctx context.Context, child SpawnChild, index int, capture bool) (*childProcess, error) {
			return s.spawnChild(ctx, child, index, capture)
		},
		func(child SpawnChild, i int, cp *childProcess) {
			childStructured := s.structuredChild(child, i)
			childStructured.Status = map[bool]string{true: "running", false: "spawned"}[capture]
			s.enrichStructuredChild(&childStructured, cp)
			_ = s.writeRelationMember(relation, childStructured)
			_ = s.appendRelationLog(relation, forkRelationLogEntry{
				Event:       "member_started",
				Status:      childStructured.Status,
				MemberIndex: &i,
				SessionID:   childStructured.SessionID,
				PID:         childStructured.PID,
			})
		},
		func(child SpawnChild, i int, cp *childProcess, err error) {
			childStructured := s.structuredChild(child, i)
			childStructured.Status = "spawn_failed"
			childStructured.Error = err.Error()
			s.enrichStructuredChild(&childStructured, cp)
			_ = s.writeRelationMember(relation, childStructured)
			_ = s.appendRelationLog(relation, forkRelationLogEntry{
				Event:       "member_spawn_failed",
				Status:      "spawn_failed",
				MemberIndex: &i,
				SessionID:   childStructured.SessionID,
				Detail:      err.Error(),
			})
		},
	)
}

func (s *SpawnExecutor) executeWait(toolID string, req SpawnRequest) tape.ToolResult {
	n := len(req.Children)
	if err := s.initState(); err != nil {
		return s.errorResult(toolID, SpawnModeWait, n, err.Error())
	}
	relation, err := s.beginRelation(toolID, req)
	if err != nil {
		return s.errorResult(toolID, SpawnModeWait, n, fmt.Sprintf("create relation surface: %v", err))
	}

	ctx, cancel := s.waitContext()
	defer cancel()

	children, spawnErrors := s.spawnAll(ctx, relation, req, true)
	spawned := countStarted(children)
	_ = s.updateRelationStatus(relation, forkRelationStatus{
		Status:           map[bool]string{true: "error", false: "running"}[spawned == 0],
		Mode:             SpawnModeWait,
		InitiatorSession: s.SessionID,
		Requested:        n,
		Spawned:          spawned,
		Errors:           append([]string(nil), spawnErrors...),
	})
	if spawned == 0 {
		return s.relationToolResult(toolID, relation, spawnStructuredResult{
			Tool:      "spawn",
			Mode:      SpawnModeWait,
			Status:    "error",
			Requested: n,
			Children:  s.spawnFailedChildren(req, children, spawnErrors),
			Errors:    spawnErrors,
		}, true)
	}

	results, completed, killed, killedChildren, timeoutError := waitForSwarmChildren(
		ctx,
		children,
		s.waitChild,
		fmt.Sprintf("spawn wait timed out after %s", s.DefaultTimeout),
		func(entry forkRelationLogEntry) error { return s.appendRelationLog(relation, entry) },
	)
	structuredChildren := make([]spawnChildStructuredResult, 0, n)
	allFailed := true
	for i, child := range req.Children {
		childStructured := s.resultChild(child, i, children[i], results[i], killedChildren[i], timeoutError, spawnErrors)
		if childStructured.ExitCode != nil && *childStructured.ExitCode == 0 && childStructured.Status == "completed" {
			allFailed = false
		}
		_ = s.writeRelationMember(relation, childStructured)
		structuredChildren = append(structuredChildren, childStructured)
	}
	status := "completed"
	errors := append([]string(nil), spawnErrors...)
	if timeoutError != "" {
		status = "timeout"
		errors = append(errors, timeoutError)
	}
	succeeded := countSucceeded(structuredChildren)
	_ = s.updateRelationStatus(relation, forkRelationStatus{
		Status:           status,
		Mode:             SpawnModeWait,
		InitiatorSession: s.SessionID,
		Requested:        n,
		Spawned:          spawned,
		Completed:        completed,
		Succeeded:        succeeded,
		Killed:           killed,
		Errors:           errors,
	})
	return s.relationToolResult(toolID, relation, spawnStructuredResult{
		Tool:      "spawn",
		Mode:      SpawnModeWait,
		Status:    status,
		Requested: n,
		Spawned:   spawned,
		Succeeded: succeeded,
		Killed:    killed,
		Children:  structuredChildren,
		Errors:    errors,
	}, allFailed)
}

func (s *SpawnExecutor) executeRace(toolID string, req SpawnRequest) tape.ToolResult {
	n := len(req.Children)
	if err := s.initState(); err != nil {
		return s.errorResult(toolID, SpawnModeRace, n, err.Error())
	}
	relation, err := s.beginRelation(toolID, req)
	if err != nil {
		return s.errorResult(toolID, SpawnModeRace, n, fmt.Sprintf("create relation surface: %v", err))
	}

	ctx, cancel := s.waitContext()
	defer cancel()

	children, spawnErrors := s.spawnAll(ctx, relation, req, true)
	spawned := countStarted(children)
	if spawned == 0 {
		return s.relationToolResult(toolID, relation, spawnStructuredResult{
			Tool:      "spawn",
			Mode:      SpawnModeRace,
			Status:    "error",
			Requested: n,
			Children:  s.spawnFailedChildren(req, children, spawnErrors),
			Errors:    spawnErrors,
		}, true)
	}

	completed, killed, killedChildren, winnerIndex, timeoutError := raceSwarmChildren(
		ctx,
		children,
		s.waitChild,
		fmt.Sprintf("spawn race timed out after %s", s.DefaultTimeout),
		func(entry forkRelationLogEntry) error { return s.appendRelationLog(relation, entry) },
	)

	if winnerIndex != nil {
		idx := *winnerIndex
		result := completed[idx]
		winner := s.resultChild(req.Children[idx], idx, children[idx], result, false, "", spawnErrors)
		_ = s.writeRelationMember(relation, winner)
		_ = s.appendRelationLog(relation, forkRelationLogEntry{
			Event:       "winner_selected",
			Status:      "completed",
			MemberIndex: &idx,
			SessionID:   winner.SessionID,
			PID:         winner.PID,
			Detail:      fmt.Sprintf("exit_code=%d", result.exitCode),
		})
		return s.relationToolResult(toolID, relation, spawnStructuredResult{
			Tool:      "spawn",
			Mode:      SpawnModeRace,
			Status:    "completed",
			Requested: n,
			Spawned:   spawned,
			Succeeded: 1,
			Killed:    killed,
			Winner:    &winner,
			Errors:    spawnErrors,
		}, false)
	}

	structuredChildren := make([]spawnChildStructuredResult, 0, n)
	for i, child := range req.Children {
		structured := s.resultChild(child, i, children[i], completed[i], killedChildren[i], timeoutError, spawnErrors)
		_ = s.writeRelationMember(relation, structured)
		structuredChildren = append(structuredChildren, structured)
	}
	status := "completed"
	errors := append([]string(nil), spawnErrors...)
	if timeoutError != "" {
		status = "timeout"
		errors = append(errors, timeoutError)
	}
	return s.relationToolResult(toolID, relation, spawnStructuredResult{
		Tool:      "spawn",
		Mode:      SpawnModeRace,
		Status:    status,
		Requested: n,
		Spawned:   spawned,
		Succeeded: countSucceeded(structuredChildren),
		Killed:    killed,
		Children:  structuredChildren,
		Errors:    errors,
	}, true)
}

func (s *SpawnExecutor) executeForget(toolID string, req SpawnRequest) tape.ToolResult {
	n := len(req.Children)
	if err := s.initState(); err != nil {
		return s.errorResult(toolID, SpawnModeForget, n, err.Error())
	}
	relation, err := s.beginRelation(toolID, req)
	if err != nil {
		return s.errorResult(toolID, SpawnModeForget, n, fmt.Sprintf("create relation surface: %v", err))
	}
	children, spawnErrors := s.spawnAll(context.Background(), relation, req, false)
	structuredChildren := make([]spawnChildStructuredResult, 0, n)
	for i, child := range req.Children {
		structured := s.structuredChild(child, i)
		if childStarted(children[i]) {
			structured.Status = "spawned"
		} else {
			structured.Status = "spawn_failed"
			structured.Error = swarmSpawnErrorFor(i, spawnErrors)
		}
		s.enrichStructuredChild(&structured, children[i])
		_ = s.writeRelationMember(relation, structured)
		structuredChildren = append(structuredChildren, structured)
	}
	spawned := countStarted(children)
	status := "spawned"
	isError := false
	if spawned == 0 {
		status = "error"
		isError = true
	}
	return s.relationToolResult(toolID, relation, spawnStructuredResult{
		Tool:      "spawn",
		Mode:      SpawnModeForget,
		Status:    status,
		Requested: n,
		Spawned:   spawned,
		Children:  structuredChildren,
		Errors:    spawnErrors,
	}, isError)
}

func (s *SpawnExecutor) errorResult(toolID, mode string, requested int, errMsg string) tape.ToolResult {
	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(spawnStructuredResult{
			Tool:      "spawn",
			Mode:      mode,
			Status:    "error",
			Requested: requested,
			Errors:    []string{errMsg},
		}),
		IsError: true,
	}
}

func (s *SpawnExecutor) resultChild(child SpawnChild, index int, cp *childProcess, r childResult, killed bool, timeoutError string, spawnErrors []string) spawnChildStructuredResult {
	structured := s.structuredChild(child, index)
	if !childStarted(cp) {
		structured.Status = "spawn_failed"
		structured.Error = swarmSpawnErrorFor(index, spawnErrors)
		s.enrichStructuredChild(&structured, cp)
		return structured
	}
	if killed {
		structured.Status = "timeout"
		structured.Error = timeoutError
		if r.child != nil {
			exitCode := r.exitCode
			structured.ExitCode = &exitCode
			structured.Stdout = s.childStdout(r.child)
			structured.Stderr = s.childStderr(r.child)
		}
		s.enrichStructuredChild(&structured, cp)
		return structured
	}
	if r.err != nil {
		structured.Status = "error"
		structured.Error = r.err.Error()
		s.enrichStructuredChild(&structured, cp)
		return structured
	}
	exitCode := r.exitCode
	structured.Status = "completed"
	structured.ExitCode = &exitCode
	structured.Stdout = s.childStdout(cp)
	structured.Stderr = s.childStderr(cp)
	s.enrichStructuredChild(&structured, cp)
	return structured
}

func (s *SpawnExecutor) spawnFailedChildren(req SpawnRequest, children []*childProcess, spawnErrors []string) []spawnChildStructuredResult {
	structuredChildren := make([]spawnChildStructuredResult, 0, len(req.Children))
	for i, child := range req.Children {
		structured := s.structuredChild(child, i)
		structured.Status = "spawn_failed"
		structured.Error = swarmSpawnErrorFor(i, spawnErrors)
		s.enrichStructuredChild(&structured, children[i])
		structuredChildren = append(structuredChildren, structured)
	}
	return structuredChildren
}

func countSucceeded(children []spawnChildStructuredResult) int {
	count := 0
	for _, child := range children {
		if child.Status == "completed" && child.ExitCode != nil && *child.ExitCode == 0 {
			count++
		}
	}
	return count
}
