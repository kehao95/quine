package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type childLauncher[T any] func(context.Context, T, int, bool) (*childProcess, error)
type childSpawnHook[T any] func(T, int, *childProcess)
type childSpawnErrorHook[T any] func(T, int, *childProcess, error)

type indexedChildResult struct {
	index  int
	result childResult
}

func beginSwarmRelation(toolName, toolID, mode, sessionID string, retainedDir func(string) string, requested int, children any) (*forkRelationSurface, error) {
	relationID := sanitizeRelationID(toolID)
	root := filepath.Join(retainedDir(sessionID), "relations", relationID)
	if err := os.MkdirAll(filepath.Join(root, "members"), 0o755); err != nil {
		return nil, fmt.Errorf("create relation root: %w", err)
	}
	relation := &forkRelationSurface{
		ID:               relationID,
		Root:             root,
		Handle:           fmt.Sprintf("relation://%s/%s", sessionID, relationID),
		ToolID:           toolID,
		Mode:             mode,
		InitiatorSession: sessionID,
		Requested:        requested,
		CreatedAt:        time.Now().UTC(),
	}
	relationDoc := map[string]any{
		"id":                relation.ID,
		"kind":              toolName,
		"tool":              toolName,
		"tool_id":           toolID,
		"mode":              mode,
		"initiator_session": sessionID,
		"created_at":        relation.CreatedAt.Format(time.RFC3339Nano),
		"requested":         requested,
		"children":          children,
	}
	if err := writeJSONFile(filepath.Join(root, "relation.json"), relationDoc); err != nil {
		return nil, fmt.Errorf("write relation metadata: %w", err)
	}
	if err := updateSwarmRelationStatus(relation, toolName, forkRelationStatus{
		Status:           "created",
		Mode:             relation.Mode,
		InitiatorSession: relation.InitiatorSession,
		Requested:        relation.Requested,
		UpdatedAt:        relation.CreatedAt.Format(time.RFC3339Nano),
	}); err != nil {
		return nil, fmt.Errorf("write relation status: %w", err)
	}
	if err := appendSwarmRelationLog(relation, forkRelationLogEntry{
		At:     relation.CreatedAt.Format(time.RFC3339Nano),
		Event:  "created",
		ToolID: toolID,
		Mode:   mode,
		Status: "created",
	}); err != nil {
		return nil, fmt.Errorf("write relation log: %w", err)
	}
	return relation, nil
}

func updateSwarmRelationStatus(relation *forkRelationSurface, toolName string, status forkRelationStatus) error {
	if relation == nil {
		return nil
	}
	status.ID = relation.ID
	status.Kind = toolName
	status.Tool = toolName
	if status.Mode == "" {
		status.Mode = relation.Mode
	}
	if status.InitiatorSession == "" {
		status.InitiatorSession = relation.InitiatorSession
	}
	if status.Requested == 0 {
		status.Requested = relation.Requested
	}
	if status.UpdatedAt == "" {
		status.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	return writeJSONFile(filepath.Join(relation.Root, "status.json"), status)
}

func writeSwarmRelationMember(relation *forkRelationSurface, index int, child any) error {
	if relation == nil {
		return nil
	}
	return writeJSONFile(filepath.Join(relation.Root, "members", fmt.Sprintf("%03d.json", index)), child)
}

func appendSwarmRelationLog(relation *forkRelationSurface, entry forkRelationLogEntry) error {
	if relation == nil {
		return nil
	}
	if entry.At == "" {
		entry.At = time.Now().UTC().Format(time.RFC3339Nano)
	}
	if entry.ToolID == "" {
		entry.ToolID = relation.ToolID
	}
	if entry.Mode == "" {
		entry.Mode = relation.Mode
	}
	return appendJSONLine(filepath.Join(relation.Root, "log.jsonl"), entry)
}

func startQuineChildProcess(ctx context.Context, cp *childProcess, quinePath, mission string, env []string, workDir, captureRoot string, capture bool, processStarted func(*os.Process)) error {
	cmd := exec.CommandContext(ctx, quinePath, mission)
	cmd.Env = append([]string(nil), env...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cp.cmd = cmd

	if capture {
		stdoutFile, stdoutPath, err := createChildCaptureFile(captureRoot, cp.sessionID, "stdout")
		if err != nil {
			return fmt.Errorf("create child %d stdout capture: %w", cp.index, err)
		}
		stderrFile, stderrPath, err := createChildCaptureFile(captureRoot, cp.sessionID, "stderr")
		if err != nil {
			_ = stdoutFile.Close()
			_ = os.Remove(stdoutPath)
			return fmt.Errorf("create child %d stderr capture: %w", cp.index, err)
		}
		cp.stdoutFile = stdoutFile
		cp.stderrFile = stderrFile
		cp.stdoutPath = stdoutPath
		cp.stderrPath = stderrPath
		cmd.Stdout = stdoutFile
		cmd.Stderr = stderrFile
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if err := cmd.Start(); err != nil {
		cp.closeCaptureFiles()
		cp.removeCaptureFiles()
		return fmt.Errorf("start child %d: %w", cp.index, err)
	}
	if processStarted != nil {
		processStarted(cmd.Process)
	}
	return nil
}

func createChildCaptureFile(root, sessionID, stream string) (*os.File, string, error) {
	dir := filepath.Join(root, sessionID)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, "", err
	}
	file, err := os.CreateTemp(dir, stream+"-*.log")
	if err != nil {
		return nil, "", err
	}
	return file, file.Name(), nil
}

func waitChildProcess(cp *childProcess, processEnded func()) childResult {
	err := cp.cmd.Wait()
	cp.closeCaptureFiles()
	if processEnded != nil {
		processEnded()
	}
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			err = nil // exit code errors are child outcomes, not wait failures.
		}
	}
	return childResult{
		child:    cp,
		exitCode: exitCode,
		err:      err,
	}
}

func startSwarmChildren[T any](ctx context.Context, children []T, capture bool, launch childLauncher[T], onStarted childSpawnHook[T], onFailed childSpawnErrorHook[T]) ([]*childProcess, []string) {
	processes := make([]*childProcess, len(children))
	var spawnErrors []string
	for i, child := range children {
		cp, err := launch(ctx, child, i, capture)
		processes[i] = cp
		if err != nil {
			spawnErrors = append(spawnErrors, fmt.Sprintf("child %d: %v", i, err))
			if onFailed != nil {
				onFailed(child, i, cp, err)
			}
			continue
		}
		if onStarted != nil {
			onStarted(child, i, cp)
		}
	}
	return processes, spawnErrors
}

func waitForSwarmChildren(ctx context.Context, children []*childProcess, waitChild func(*childProcess) childResult, timeoutLabel string, appendLog func(forkRelationLogEntry) error) ([]childResult, int, int, map[int]bool, string) {
	spawned := countStarted(children)
	results := make([]childResult, len(children))
	resultCh := make(chan indexedChildResult, spawned)
	for i, cp := range children {
		if !childStarted(cp) {
			results[i] = childResult{err: fmt.Errorf("failed to spawn")}
			continue
		}
		go func(idx int, child *childProcess) {
			resultCh <- indexedChildResult{index: idx, result: waitChild(child)}
		}(i, cp)
	}

	completed := 0
	killed := 0
	killedChildren := make(map[int]bool)
	timeoutError := ""
waitLoop:
	for completed < spawned {
		select {
		case ir := <-resultCh:
			results[ir.index] = ir.result
			completed++
		case <-ctx.Done():
			timeoutError = timeoutLabel
			if timeoutError == "" {
				timeoutError = "swarm wait timed out"
			}
			for i, cp := range children {
				if !childStarted(cp) {
					continue
				}
				if results[i].child != nil || results[i].err != nil {
					continue
				}
				killProcessGroup(cp)
				killedChildren[i] = true
				killed++
				if appendLog != nil {
					_ = appendLog(forkRelationLogEntry{
						Event:       "member_killed",
						Status:      "timeout",
						MemberIndex: &i,
						SessionID:   cp.sessionID,
						PID:         cp.cmd.Process.Pid,
						Detail:      timeoutError,
					})
				}
			}
			deadline := time.After(2 * time.Second)
			for completed < spawned {
				select {
				case ir := <-resultCh:
					results[ir.index] = ir.result
					completed++
				case <-deadline:
					break waitLoop
				}
			}
			break waitLoop
		}
	}
	return results, completed, killed, killedChildren, timeoutError
}

func raceSwarmChildren(ctx context.Context, children []*childProcess, waitChild func(*childProcess) childResult, timeoutLabel string, appendLog func(forkRelationLogEntry) error) (map[int]childResult, int, map[int]bool, *int, string) {
	spawned := countStarted(children)
	resultCh := make(chan indexedChildResult, spawned)
	for i, cp := range children {
		if !childStarted(cp) {
			continue
		}
		go func(idx int, child *childProcess) {
			resultCh <- indexedChildResult{index: idx, result: waitChild(child)}
		}(i, cp)
	}

	completed := make(map[int]childResult)
	killedChildren := make(map[int]bool)
	killed := 0
	var winnerIndex *int
	timeoutError := ""
	for len(completed) < spawned {
		select {
		case ir := <-resultCh:
			completed[ir.index] = ir.result
			if ir.result.err == nil && ir.result.exitCode == 0 && winnerIndex == nil {
				idx := ir.index
				winnerIndex = &idx
				for i, cp := range children {
					if i == idx || !childStarted(cp) {
						continue
					}
					if _, done := completed[i]; done {
						continue
					}
					killProcessGroup(cp)
					killedChildren[i] = true
					killed++
				}
			}
		case <-ctx.Done():
			timeoutError = timeoutLabel
			if timeoutError == "" {
				timeoutError = "swarm race timed out"
			}
			for i, cp := range children {
				if !childStarted(cp) {
					continue
				}
				if _, done := completed[i]; done {
					continue
				}
				killProcessGroup(cp)
				killedChildren[i] = true
				killed++
				if appendLog != nil {
					_ = appendLog(forkRelationLogEntry{
						Event:       "member_killed",
						Status:      "timeout",
						MemberIndex: &i,
						SessionID:   cp.sessionID,
						PID:         cp.cmd.Process.Pid,
						Detail:      timeoutError,
					})
				}
			}
		}
		if winnerIndex != nil || timeoutError != "" {
			break
		}
	}

	drainDeadline := time.After(2 * time.Second)
	for len(completed) < spawned {
		select {
		case ir := <-resultCh:
			completed[ir.index] = ir.result
		case <-drainDeadline:
			return completed, killed, killedChildren, winnerIndex, timeoutError
		}
	}
	return completed, killed, killedChildren, winnerIndex, timeoutError
}

func countStarted(children []*childProcess) int {
	count := 0
	for _, cp := range children {
		if childStarted(cp) {
			count++
		}
	}
	return count
}

func swarmSpawnErrorFor(index int, spawnErrors []string) string {
	prefix := fmt.Sprintf("child %d:", index)
	for _, errMsg := range spawnErrors {
		if strings.HasPrefix(errMsg, prefix) {
			return errMsg
		}
	}
	return "failed to spawn"
}
