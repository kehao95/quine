package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/kehao95/quine/internal/tape"
	"github.com/kehao95/quine/internal/tools"
)

func (r *Runtime) recoverIncompleteToolBatches() (int, bool) {
	assistant, pending, ok, err := r.pendingToolBatchFromContext()
	if err != nil {
		r.log("pending tool recovery scan error: %v", err)
		return 0, false
	}
	if !ok {
		return 0, false
	}

	for _, tc := range pending {
		r.flushPendingToolResult()
		if result, ok := r.findRetainedToolResult(tc.ID); ok {
			r.appendToolResult(result)
			continue
		}
		if code, done := r.recoverPendingToolCall(assistant, tc); done {
			return code, true
		}
	}
	return 0, false
}

func (r *Runtime) pendingToolBatchFromContext() (tape.Message, []tape.ToolCall, bool, error) {
	entries, err := readContextEntries(r.contextCurrentPath())
	if err != nil {
		return tape.Message{}, nil, false, err
	}
	var (
		assistant tape.Message
		pending   map[string]tape.ToolCall
		order     []string
	)
	for _, entry := range entries {
		switch entry.Type {
		case "message":
			var msg tape.Message
			if err := json.Unmarshal(entry.Data, &msg); err != nil {
				return tape.Message{}, nil, false, err
			}
			if msg.Role == tape.RoleAssistant && len(msg.ToolCalls) > 0 {
				assistant = msg
				pending = make(map[string]tape.ToolCall, len(msg.ToolCalls))
				order = order[:0]
				for _, tc := range msg.ToolCalls {
					pending[tc.ID] = tc
					order = append(order, tc.ID)
				}
				continue
			}
			pending = nil
			order = nil
		case "tool_result":
			if pending == nil {
				continue
			}
			var result tape.ToolResult
			if err := json.Unmarshal(entry.Data, &result); err != nil {
				return tape.Message{}, nil, false, err
			}
			delete(pending, result.ToolID)
			if len(pending) == 0 {
				pending = nil
				order = nil
			}
		}
	}
	if len(pending) == 0 {
		return tape.Message{}, nil, false, nil
	}
	out := make([]tape.ToolCall, 0, len(order))
	for _, id := range order {
		if tc, ok := pending[id]; ok {
			out = append(out, tc)
		}
	}
	return assistant, out, true, nil
}

// readContextEntries reads the recovery surface. Invariant (C2): it reads ONLY
// current.jsonl (contextCurrentPath); the transient live.jsonl (contextLivePath)
// is never recovered — its deltas are display-only and never crystallized state.
func readContextEntries(path string) ([]tape.TapeEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	out := make([]tape.TapeEntry, 0, len(lines))
	for _, raw := range lines {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		var entry tape.TapeEntry
		if err := json.Unmarshal([]byte(raw), &entry); err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	return out, nil
}

func (r *Runtime) recoverPendingToolCall(source tape.Message, tc tape.ToolCall) (int, bool) {
	switch tc.Name {
	case "exit":
		if !r.cfg.ExitEnabled() {
			r.appendUnknownToolResult(tc, "exit is disabled in this runtime")
			return 0, false
		}
		code, ok := r.handleExit(tc)
		return code, ok
	case "idle":
		if !r.cfg.IdleToolEnabled() {
			r.appendUnknownToolResult(tc, "idle is disabled in this runtime")
			return 0, false
		}
		r.handleIdle(tc)
	case "vision":
		if !r.cfg.VisionEnabled {
			r.appendUnknownToolResult(tc, "vision is disabled in this runtime")
			return 0, false
		}
		r.handleVision(tc)
	case "unfold":
		r.handleUnfold(tc)
	case "switch_world":
		r.handleSwitchWorld(tc)
	case "sh":
		if result, ok := r.recoverShToolResult(tc); ok {
			r.appendToolResult(result)
		} else {
			r.appendUnknownToolResult(tc, "failed to gather sh result; execution may or may not have succeeded")
		}
	case "fork", "spawn":
		if result, ok := r.recoverRelationToolResult(tc); ok {
			r.appendToolResult(result)
		} else {
			r.appendUnknownToolResult(tc, fmt.Sprintf("failed to gather %s result; child execution may or may not have completed", tc.Name))
		}
	case "mark":
		if result, ok := r.recoverMarkToolResult(tc); ok {
			r.appendToolResult(result)
		} else {
			r.appendUnknownToolResult(tc, "failed to gather mark result; anchor creation may or may not have succeeded")
		}
	case "exec":
		r.appendUnknownToolResult(tc, "failed to gather exec result; process image replacement may or may not have succeeded")
	default:
		r.appendUnknownToolResult(tc, fmt.Sprintf("unknown tool: %s", tc.Name))
	}
	_ = source
	return 0, false
}

func (r *Runtime) appendUnknownToolResult(tc tape.ToolCall, reason string) {
	msg := runtimeToolResultMessage(tc.ID, tc.Name, "unknown", map[string]any{
		"error":        reason,
		"side_effects": "unknown",
	})
	r.appendRuntimeToolMessage(msg, false)
}

func (r *Runtime) findRetainedToolResult(toolID string) (tape.ToolResult, bool) {
	pattern := filepath.Join(r.cfg.TapeDir(""), "*.jsonl")
	paths, err := filepath.Glob(pattern)
	if err != nil || len(paths) == 0 {
		return tape.ToolResult{}, false
	}
	sort.Strings(paths)
	for i := len(paths) - 1; i >= 0; i-- {
		result, ok := findToolResultInTape(paths[i], toolID)
		if ok {
			return result, true
		}
	}
	return tape.ToolResult{}, false
}

func findToolResultInTape(path string, toolID string) (tape.ToolResult, bool) {
	entries, err := readContextEntries(path)
	if err != nil {
		return tape.ToolResult{}, false
	}
	var found tape.ToolResult
	var ok bool
	for _, entry := range entries {
		if entry.Type != "tool_result" {
			continue
		}
		var result tape.ToolResult
		if err := json.Unmarshal(entry.Data, &result); err != nil {
			continue
		}
		if result.ToolID == toolID {
			found = result
			ok = true
		}
	}
	return found, ok
}

func (r *Runtime) recoverRelationToolResult(tc tape.ToolCall) (tape.ToolResult, bool) {
	path := filepath.Join(r.cfg.SessionRetainedDir(""), "relations", sanitizeRecoveryID(tc.ID), "result.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return tape.ToolResult{}, false
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return tape.ToolResult{}, false
	}
	return tape.ToolResult{
		ToolID:  tc.ID,
		Content: tape.MarshalToolResultContent(payload),
		IsError: relationPayloadIsError(payload),
	}, true
}

func relationPayloadIsError(payload map[string]any) bool {
	status, _ := payload["status"].(string)
	switch strings.TrimSpace(status) {
	case "":
		return false
	case "completed":
		if succeeded, ok := relationPayloadInt(payload, "succeeded"); ok {
			if succeeded > 0 {
				return false
			}
			return relationPayloadHasFailureEvidence(payload)
		}
		if relationPayloadChildrenAllFailed(payload) {
			return true
		}
		return false
	case "spawned":
		if spawned, ok := relationPayloadInt(payload, "spawned"); ok {
			return spawned == 0 && relationPayloadHasFailureEvidence(payload)
		}
		return false
	case "timeout":
		if succeeded, ok := relationPayloadInt(payload, "succeeded"); ok {
			return succeeded == 0
		}
		return true
	case "error", "rejected", "unknown":
		return true
	default:
		return true
	}
}

func relationPayloadHasFailureEvidence(payload map[string]any) bool {
	if requested, ok := relationPayloadInt(payload, "requested"); ok && requested > 0 {
		return true
	}
	if spawned, ok := relationPayloadInt(payload, "spawned"); ok && spawned > 0 {
		return true
	}
	if children, ok := payload["children"].([]any); ok && len(children) > 0 {
		return true
	}
	if errors, ok := payload["errors"].([]any); ok && len(errors) > 0 {
		return true
	}
	return false
}

func relationPayloadChildrenAllFailed(payload map[string]any) bool {
	children, ok := payload["children"].([]any)
	if !ok || len(children) == 0 {
		return false
	}
	sawTerminal := false
	for _, raw := range children {
		child, ok := raw.(map[string]any)
		if !ok {
			return false
		}
		status, _ := child["status"].(string)
		switch strings.TrimSpace(status) {
		case "completed":
			exitCode, ok := relationMapInt(child, "exit_code")
			if ok && exitCode == 0 {
				return false
			}
			sawTerminal = true
		case "error", "spawn_failed", "killed", "timeout", "no_result":
			sawTerminal = true
		case "spawned", "running", "":
			return false
		default:
			return false
		}
	}
	return sawTerminal
}

func relationPayloadInt(payload map[string]any, key string) (int, bool) {
	return relationMapInt(payload, key)
}

func relationMapInt(payload map[string]any, key string) (int, bool) {
	switch value := payload[key].(type) {
	case json.Number:
		n, err := value.Int64()
		if err != nil {
			return 0, false
		}
		return int(n), true
	case float64:
		return int(value), true
	case int:
		return value, true
	case int64:
		return int(value), true
	default:
		return 0, false
	}
}

func sanitizeRecoveryID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
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
	return strings.Trim(b.String(), "._-")
}

func (r *Runtime) recoverShToolResult(tc tape.ToolCall) (tape.ToolResult, bool) {
	command, _ := tc.Arguments["command"].(string)
	command = strings.TrimSpace(command)
	if command == "" {
		return tape.ToolResult{}, false
	}
	matches, err := r.matchingJobDirs(command)
	if err != nil || len(matches) != 1 {
		return tape.ToolResult{}, false
	}
	jobDir := matches[0]
	pid, err := strconv.Atoi(strings.TrimSpace(readText(filepath.Join(jobDir, "pid"))))
	if err != nil || pid <= 0 {
		// A matched job dir whose `pid` file is missing/empty/unparseable is a
		// crash-orphaned staging dir (`cmd` written, `pid` not yet published).
		// Reporting pid:0 would fabricate a live job handle; decline recovery.
		return tape.ToolResult{}, false
	}
	display := filepath.ToSlash(jobDir) + "/"
	interactive, _ := tc.Arguments["interactive"].(bool)
	detach, _ := tc.Arguments["detach"].(bool)
	if interactive || detach {
		job := map[string]any{
			"pid":  pid,
			"path": display,
		}
		mode := "detached"
		if interactive {
			mode = "interactive"
			job["interactive"] = true
		}
		if detach {
			job["detached"] = true
		}
		return tape.ToolResult{
			ToolID: tc.ID,
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":   "sh",
				"mode":   mode,
				"status": "spawned",
				"job":    job,
			}),
		}, true
	}

	out := readText(filepath.Join(jobDir, "out.log"))
	errText := readText(filepath.Join(jobDir, "err.log"))
	if exitText := strings.TrimSpace(readText(filepath.Join(jobDir, "exit"))); exitText != "" {
		exitCode, parseErr := strconv.Atoi(exitText)
		if parseErr != nil {
			return tape.ToolResult{}, false
		}
		return tape.ToolResult{
			ToolID: tc.ID,
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":      "sh",
				"mode":      "sync",
				"status":    "completed",
				"exit_code": exitCode,
				"stdout":    out,
				"stderr":    errText,
			}),
			IsError: exitCode != 0,
		}, true
	}
	return tape.ToolResult{
		ToolID: tc.ID,
		Content: tape.MarshalToolResultContent(map[string]any{
			"tool":          "sh",
			"mode":          "sync",
			"status":        "interrupted",
			"job":           map[string]any{"pid": pid, "path": display},
			"stdout_so_far": out,
			"stderr_so_far": errText,
			"cause":         "result_unavailable",
		}),
	}, true
}

func (r *Runtime) matchingJobDirs(command string) ([]string, error) {
	root := r.cfg.JobSessionDir("")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var matches []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		if strings.TrimSpace(readText(filepath.Join(dir, "cmd"))) == command {
			matches = append(matches, dir)
		}
	}
	sort.Strings(matches)
	return matches, nil
}

func readText(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (r *Runtime) recoverMarkToolResult(tc tape.ToolCall) (tape.ToolResult, bool) {
	req, err := tools.ParseMarkArgs(tc.Arguments)
	if err != nil {
		return tape.ToolResult{}, false
	}
	root := filepath.Join(r.contextAnchorsRoot())
	entries, err := os.ReadDir(root)
	if err != nil {
		return tape.ToolResult{}, false
	}
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasSuffix(entry.Name(), ".anchor") || strings.HasPrefix(entry.Name(), ".staging-") {
			continue
		}
		metaPath := filepath.Join(root, entry.Name(), "meta.json")
		var meta struct {
			ID         int    `json:"id"`
			Resolution string `json:"resolution"`
		}
		data, err := os.ReadFile(metaPath)
		if err != nil || json.Unmarshal(data, &meta) != nil {
			continue
		}
		if strings.TrimSpace(meta.Resolution) != req.Resolution {
			continue
		}
		return tape.ToolResult{
			ToolID: tc.ID,
			Content: tape.MarshalToolResultContent(map[string]any{
				"tool":       "mark",
				"status":     "completed",
				"anchor_id":  meta.ID,
				"fold":       req.Fold,
				"resolution": req.Resolution,
				"note":       "gathered from existing anchor surface",
			}),
		}, true
	}
	return tape.ToolResult{}, false
}
