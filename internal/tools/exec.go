package tools

import (
	"fmt"
	"os"
	"syscall"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// ExecRequest represents the parsed arguments from an exec tool call.
type ExecRequest struct {
	Persona string            // Optional persona name
	Wisdom  map[string]string // Key-value pairs to pass to the new instance
}

// ParseExecArgs extracts ExecRequest from a ToolCall's Arguments map.
func ParseExecArgs(args map[string]any) (ExecRequest, error) {
	req := ExecRequest{}

	if v, ok := args["persona"]; ok {
		s, ok := v.(string)
		if !ok {
			return ExecRequest{}, fmt.Errorf("persona must be a string, got %T", v)
		}
		req.Persona = s
	}

	// Parse wisdom map
	if v, ok := args["wisdom"]; ok {
		wisdomMap, ok := v.(map[string]any)
		if !ok {
			return ExecRequest{}, fmt.Errorf("wisdom must be an object, got %T", v)
		}
		req.Wisdom = make(map[string]string)
		for k, val := range wisdomMap {
			strVal, ok := val.(string)
			if !ok {
				return ExecRequest{}, fmt.Errorf("wisdom values must be strings, key %q has %T", k, val)
			}
			req.Wisdom[k] = strVal
		}
	}

	return req, nil
}

// ExecExecutor handles the exec (metamorphosis) tool.
// Unlike other tools, exec replaces the current process entirely.
type ExecExecutor struct {
	// QuinePath is the path to the quine binary.
	QuinePath string

	// Cfg is the current configuration.
	Cfg *config.Config

	// OriginalIntent is the original input from stdin that must be preserved.
	OriginalIntent string
}

// NewExecExecutor creates an ExecExecutor from config.
func NewExecExecutor(cfg *config.Config, originalIntent string) *ExecExecutor {
	quinePath, err := os.Executable()
	if err != nil {
		quinePath = "./quine"
	}

	return &ExecExecutor{
		QuinePath:      quinePath,
		Cfg:            cfg,
		OriginalIntent: originalIntent,
	}
}

// Execute performs the exec syscall, replacing the current process with a
// fresh quine instance. This function does not return on success.
//
// The new process gets:
//   - Fresh tape (new SESSION_ID)
//   - Same mission (passed via argv, preserved from original startup)
//   - All QUINE_WISDOM_* vars preserved (learned insights survive)
//   - New wisdom from the exec call merged in (overwrites existing keys)
//   - QUINE_PARENT_SESSION set for lineage tracking
//   - QUINE_DEPTH reset to 0 (fresh brain, not deeper recursion)
//
// Returns a ToolResult only on failure (exec syscall failed).
func (e *ExecExecutor) Execute(toolID string, req ExecRequest) tape.ToolResult {
	// Build environment for the new process
	execEnv, err := e.Cfg.ExecEnv(e.OriginalIntent)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[EXEC ERROR] Failed to build environment: %v", err),
			IsError: true,
		}
	}

	// Add new wisdom from the exec request (these override existing keys)
	for key, value := range req.Wisdom {
		execEnv = append(execEnv, "QUINE_WISDOM_"+key+"="+value)
	}

	// Merge with filtered OS environment (need PATH, HOME, etc.)
	fullEnv := MergeEnv(filterSessionID(os.Environ()), execEnv)

	// The exec syscall replaces the current process image.
	// Mission is passed via argv (argv[0] = binary, argv[1] = mission)
	//
	// The new process will:
	// 1. Read mission from argv[1] (or QUINE_ORIGINAL_INTENT if set)
	// 2. Generate a fresh SESSION_ID
	// 3. Start with an empty tape
	// 4. stdin remains available (data stream)

	// Perform the exec - this does not return on success
	err = syscall.Exec(e.QuinePath, []string{e.QuinePath, e.OriginalIntent}, fullEnv)

	// If we get here, exec failed
	return tape.ToolResult{
		ToolID:  toolID,
		Content: fmt.Sprintf("[EXEC ERROR] syscall.Exec failed: %v", err),
		IsError: true,
	}
}
