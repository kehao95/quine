package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// procExeRe matches a procfs handle to a live process's in-memory executable
// image: /proc/self/exe or /proc/<pid>/exe. Re-executing such a target after
// the ephemeral body has been unlinked re-enters the original image in place
// (body recovery) rather than reconstructing a successor.
var procExeRe = regexp.MustCompile(`^/proc/(self|\d+)/exe$`)

// isLiveProcessImageTarget reports whether target re-enters the current
// process's in-memory image via procfs. It checks the cleaned path and one
// level of symlink indirection (an agent may point a workspace symlink at
// /proc/self/exe and exec that).
func isLiveProcessImageTarget(target string) bool {
	t := filepath.Clean(strings.TrimSpace(target))
	if procExeRe.MatchString(t) {
		return true
	}
	if resolved, err := os.Readlink(t); err == nil {
		if procExeRe.MatchString(filepath.Clean(resolved)) {
			return true
		}
	}
	return false
}

type execStructuredResult struct {
	Tool   string   `json:"tool"`
	Status string   `json:"status"`
	Target string   `json:"target,omitempty"`
	Argv   []string `json:"argv,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// ExecRequest represents the parsed arguments from an exec tool call.
type ExecRequest struct {
	Target string   // Optional target binary path or executable name
	Argv   []string // Optional full argv vector
}

// ParseExecArgs extracts ExecRequest from a ToolCall's Arguments map.
func ParseExecArgs(args map[string]any) (ExecRequest, error) {
	req := ExecRequest{}

	if v, ok := args["target"]; ok {
		s, ok := v.(string)
		if !ok {
			return ExecRequest{}, fmt.Errorf("target must be a string, got %T", v)
		}
		req.Target = strings.TrimSpace(s)
	}

	if v, ok := args["argv"]; ok {
		rawArgv, ok := v.([]any)
		if !ok {
			return ExecRequest{}, fmt.Errorf("argv must be an array, got %T", v)
		}
		req.Argv = make([]string, 0, len(rawArgv))
		for i, val := range rawArgv {
			s, ok := val.(string)
			if !ok {
				return ExecRequest{}, fmt.Errorf("argv[%d] must be a string, got %T", i, val)
			}
			req.Argv = append(req.Argv, s)
		}
		if len(req.Argv) == 0 {
			req.Argv = nil
		}
	}

	return req, nil
}

// ExecExecutor handles the exec (metamorphosis) tool.
// Unlike other tools, exec replaces the current process entirely.
type ExecExecutor struct {
	// QuinePath is quine's default self-reentry target.
	QuinePath string

	// Cfg is the current configuration.
	Cfg *config.Config

	// ContextRoot is the live current-incarnation context surface that should be
	// projected forward across quine re-entry.
	ContextRoot string

	// Mission is the current mission string used for default self re-exec argv.
	Mission string
}

// NewExecExecutor creates an ExecExecutor from config.
func NewExecExecutor(cfg *config.Config, mission string) *ExecExecutor {
	return &ExecExecutor{
		QuinePath:   strings.TrimSpace(cfg.SelfReentryTarget),
		Cfg:         cfg,
		ContextRoot: filepath.Join(cfg.AgentRoot(), "context"),
		Mission:     mission,
	}
}

// Execute performs the exec syscall, replacing the current process image.
// This function does not return on success.
//
// Returns a ToolResult only on failure (exec syscall failed).
func (e *ExecExecutor) Execute(toolID string, req ExecRequest) tape.ToolResult {
	target, argv, err := e.buildTargetAndArgv(req)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(execStructuredResult{
				Tool:   "exec",
				Status: "error",
				Target: req.Target,
				Argv:   req.Argv,
				Error:  fmt.Sprintf("[EXEC ERROR] %v", err),
			}),
			IsError: true,
		}
	}

	// Body-recovery guard: once the ephemeral body has been unlinked, exec-ing
	// the live process image via procfs (/proc/self/exe or /proc/<pid>/exe)
	// re-enters the original in-memory body in place. That is body recovery,
	// not reconstruction of a successor from externalized state, so it is
	// rejected here rather than performed.
	if e.Cfg != nil && e.Cfg.EphemeralBody && isLiveProcessImageTarget(target) {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(execStructuredResult{
				Tool:   "exec",
				Status: "error",
				Target: target,
				Argv:   argv,
				Error: "[EXEC ERROR] re-executing the live process image via procfs " +
					"(/proc/self/exe or /proc/<pid>/exe) is not permitted after the " +
					"ephemeral body has been unlinked: this recovers the original body " +
					"rather than reconstructing a successor. Build a successor body from " +
					"the runtime contract and workspace, then exec that instead.",
			}),
			IsError: true,
		}
	}

	// Build environment for the new process.
	execEnv, err := e.Cfg.ExecEnv()
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(execStructuredResult{
				Tool:   "exec",
				Status: "error",
				Target: target,
				Argv:   argv,
				Error:  fmt.Sprintf("[EXEC ERROR] Failed to build environment: %v", err),
			}),
			IsError: true,
		}
	}

	// Staged-config merge (registry-design-brief § C, work order T3.1): the
	// slot freed by the wisdom-overlay deletion (D5). config/next.env is
	// validated against the RUNNING binary's registry and merged over the
	// ExecEnv() serialization; this exec path is the only merge site — the
	// merge never touches baseEnv(), so fork/spawn children (ChildEnv) never
	// see staged values. A validation failure is a normal exec tool error:
	// the runtime's deregister->fail->re-register recovery handles it, the
	// staged file stays intact, and the agent can fix it and retry. The file
	// is also left intact on the success path — a failed syscall.Exec must
	// find it unchanged (idempotent retry); the successor archives and
	// clears it at bootstrap (T3.2).
	staged, err := config.ReadStagedOverrides(e.Cfg.StagedNextEnvPath())
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(execStructuredResult{
				Tool:   "exec",
				Status: "error",
				Target: target,
				Argv:   argv,
				Error:  fmt.Sprintf("[EXEC ERROR] %v", err),
			}),
			IsError: true,
		}
	}
	execEnv = config.MergeStagedOverrides(execEnv, staged)

	if bootstrapRoot, err := stageExecContextBootstrap(e.Cfg.DataDir, e.ContextRoot); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(execStructuredResult{
				Tool:   "exec",
				Status: "error",
				Target: target,
				Argv:   argv,
				Error:  fmt.Sprintf("[EXEC ERROR] Failed to stage bootstrap context: %v", err),
			}),
			IsError: true,
		}
	} else if strings.TrimSpace(bootstrapRoot) != "" {
		execEnv = append(execEnv, ContextBootstrapEnv+"="+bootstrapRoot)
	}

	// Merge with filtered OS environment (need PATH, HOME, etc.)
	fullEnv := MergeEnv(filterProcessIdentity(os.Environ()), execEnv)
	fullEnv = MergeEnv(fullEnv, execProcessSurfaceEnv(e.Cfg))

	// Perform the exec - this does not return on success.
	err = syscall.Exec(target, argv, fullEnv)

	// If we get here, exec failed.
	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(execStructuredResult{
			Tool:   "exec",
			Status: "error",
			Target: target,
			Argv:   argv,
			Error:  fmt.Sprintf("[EXEC ERROR] syscall.Exec failed: %v", err),
		}),
		IsError: true,
	}
}

func (e *ExecExecutor) buildTargetAndArgv(req ExecRequest) (string, []string, error) {
	target := strings.TrimSpace(req.Target)
	if target == "" {
		target = e.QuinePath
	} else if !strings.Contains(target, "/") {
		resolved, err := exec.LookPath(target)
		if err != nil {
			return "", nil, fmt.Errorf("resolve target %q: %w", target, err)
		}
		target = resolved
	} else if !filepath.IsAbs(target) {
		target = filepath.Join(e.relativeExecBase(), target)
	}

	if req.Argv != nil {
		return target, req.Argv, nil
	}
	if strings.TrimSpace(req.Target) == "" {
		if strings.TrimSpace(e.Mission) == "" {
			return target, []string{target}, nil
		}
		return target, []string{target, e.Mission}, nil
	}
	return target, []string{target}, nil
}

func (e *ExecExecutor) relativeExecBase() string {
	if e != nil && e.Cfg != nil {
		for _, candidate := range []string{
			e.Cfg.Workspace,
			e.Cfg.WorkspaceRoot,
			e.Cfg.WorkDir,
		} {
			candidate = strings.TrimSpace(candidate)
			if candidate != "" {
				return candidate
			}
		}
	}
	if wd, err := os.Getwd(); err == nil && strings.TrimSpace(wd) != "" {
		return wd
	}
	return "."
}

func execProcessSurfaceEnv(cfg *config.Config) []string {
	if cfg == nil {
		return nil
	}
	agentRoot := strings.TrimSpace(cfg.AgentRoot())
	if agentRoot == "" {
		return nil
	}
	return []string{"QUINE_AGENT_ROOT=" + agentRoot}
}

func stageExecContextBootstrap(dataDir, contextRoot string) (string, error) {
	contextRoot = strings.TrimSpace(contextRoot)
	if contextRoot == "" {
		return "", nil
	}
	resolvedRoot := contextRoot
	if resolved, err := filepath.EvalSymlinks(contextRoot); err == nil && strings.TrimSpace(resolved) != "" {
		resolvedRoot = resolved
	}
	info, err := os.Stat(resolvedRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("stat context root %s: %w", resolvedRoot, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("context root %q is not a directory", resolvedRoot)
	}

	tmpDir, err := os.MkdirTemp(strings.TrimSpace(dataDir), "exec-context-*")
	if err != nil {
		return "", fmt.Errorf("create exec bootstrap context: %w", err)
	}
	if err := CopyTreePreservingSymlinks(resolvedRoot, tmpDir); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", fmt.Errorf("copy exec bootstrap context: %w", err)
	}
	if err := sanitizeExecBootstrappedCurrent(filepath.Join(tmpDir, "state", "current.jsonl")); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}
	if err := writeExecBootstrapPrompt(tmpDir); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}
	return tmpDir, nil
}

func sanitizeExecBootstrappedCurrent(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read exec bootstrapped current context: %w", err)
	}
	projected, changed, err := projectExecContextCurrent(data)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return os.WriteFile(path, projected, 0o644)
}

func projectExecContextCurrent(data []byte) ([]byte, bool, error) {
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
			return nil, false, fmt.Errorf("decode exec context entry %d: %w", i, err)
		}

		switch entry.Type {
		case "message":
			var msg tape.Message
			if err := json.Unmarshal(entry.Data, &msg); err != nil {
				return nil, false, fmt.Errorf("decode exec context message %d: %w", i, err)
			}
			if msg.Role == tape.RoleAssistant && len(msg.ToolCalls) > 0 {
				changed = true
				continue
			}
			if msg.Role == tape.RoleAssistant && (strings.TrimSpace(msg.ReasoningContent) != "" || len(msg.ReasoningItems) > 0) {
				msg.ReasoningContent = ""
				msg.ReasoningItems = nil
				if strings.TrimSpace(msg.Content) == "" && !tape.HasStructuredContent(msg.StructuredContent) {
					changed = true
					continue
				}
				entry = tape.MessageEntry(msg)
				encodedLine, err := json.Marshal(entry)
				if err != nil {
					return nil, false, fmt.Errorf("re-encode exec context message %d: %w", i, err)
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

func writeExecBootstrapPrompt(contextRoot string) error {
	promptPath := filepath.Join(contextRoot, "prompt", "35-exec-handoff.md")
	const body = `### Exec Handoff
- This incarnation is already the successor of a completed ` + "`exec`" + ` handoff.
- The current mission in ` + "`40-mission.md`" + ` is the authoritative next step.
- Do not replay the prior handoff just because the earlier incarnation used ` + "`exec`" + `.
`
	if err := os.MkdirAll(filepath.Dir(promptPath), 0o755); err != nil {
		return fmt.Errorf("create exec handoff prompt dir: %w", err)
	}
	if err := os.WriteFile(promptPath, []byte(body), 0o644); err != nil {
		return fmt.Errorf("write exec handoff prompt: %w", err)
	}
	return nil
}
