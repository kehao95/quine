package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

// ForkExecutor spawns child quine processes with cloned context.
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

	// TapePath is the current session's tape path under DataDir.
	TapePath string

	// DefaultTimeout is the maximum time to wait for the entire swarm
	// operation (when wait=true). Not per-child.
	DefaultTimeout time.Duration

	// MaxOutput limits the captured output size per child.
	MaxOutput int

	// ProcessStarted is called when a child process starts.
	// Only used for single-child backward compat, set to nil for swarm.
	ProcessStarted func(*os.Process)

	// ProcessEnded is called when a child process ends.
	// Only used for single-child backward compat, set to nil for swarm.
	ProcessEnded func()

	WorkspaceEnabled bool
	WorkspaceRoot    string
	Workspace        string
	WorkspaceSession string
}

// NewForkExecutor creates a ForkExecutor from config with the given child
// environment. The childEnv slice should contain QUINE_* overrides.
func NewForkExecutor(cfg *config.Config, childEnv []string) *ForkExecutor {
	// Get the current executable path for spawning children
	quinePath, err := os.Executable()
	if err != nil {
		quinePath = "./quine"
	}

	// Copy the current tape incarnation when seeding child context.
	tapePath := cfg.TapePath("")

	return &ForkExecutor{
		QuinePath:        quinePath,
		DataDir:          cfg.DataDir,
		SessionID:        cfg.SessionID,
		Env:              MergeEnv(filterProcessIdentity(os.Environ()), childEnv),
		TapePath:         tapePath,
		DefaultTimeout:   time.Duration(cfg.ShTimeout) * time.Second,
		MaxOutput:        cfg.OutputTruncate,
		WorkspaceEnabled: cfg.WorkspaceEnabled,
		WorkspaceRoot:    cfg.WorkspaceRoot,
		Workspace:        cfg.Workspace,
		WorkspaceSession: cfg.WorkspaceSession,
	}
}

// filterProcessIdentity removes per-process identity from an environment slice.
func filterProcessIdentity(env []string) []string {
	result := make([]string, 0, len(env))
	for _, e := range env {
		if len(e) > 17 && e[:17] == "QUINE_SESSION_ID=" {
			continue
		}
		if len(e) > 14 && e[:14] == "QUINE_TAPE_ID=" {
			continue
		}
		result = append(result, e)
	}
	return result
}

// Fork mode constants.
const (
	ForkModeWait   = "wait"   // Block until all children finish
	ForkModeRace   = "race"   // First exit-0 child wins, rest killed (default)
	ForkModeForget = "forget" // Fire-and-forget, return PIDs immediately
)

// ForkRequest represents the parsed arguments from a fork tool call.
type ForkRequest struct {
	Children []ForkChild
	Mode     string // "race" (default), "wait", or "forget"
}

type ForkChild struct {
	Intent    string
	Workspace string
}

// ParseForkArgs extracts ForkRequest from a ToolCall's Arguments map.
func ParseForkArgs(args map[string]any) (ForkRequest, error) {
	children, err := parseForkChildren(args)
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

	// Backward compat: accept old wait/race booleans
	if v, ok := args["race"]; ok {
		if b, ok := v.(bool); ok && b {
			req.Mode = ForkModeRace
		}
	}
	if v, ok := args["wait"]; ok {
		if b, ok := v.(bool); ok && !b {
			req.Mode = ForkModeForget
		}
	}

	return req, nil
}

func parseForkChildren(args map[string]any) ([]ForkChild, error) {
	if raw, ok := args["children"]; ok {
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

			workspace, ok := obj["workspace"].(string)
			if !ok {
				return nil, fmt.Errorf("children[%d].workspace must be a string", i)
			}
			workspace = strings.TrimSpace(workspace)
			if workspace == "" {
				return nil, fmt.Errorf("children[%d].workspace cannot be empty", i)
			}

			children = append(children, ForkChild{
				Intent:    intent,
				Workspace: workspace,
			})
		}
		return children, nil
	}

	raw, ok := args["argv"]
	if !ok {
		return nil, fmt.Errorf("missing required argument: children")
	}

	rawSlice, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("argv must be an array, got %T", raw)
	}
	if len(rawSlice) == 0 {
		return nil, fmt.Errorf("argv must contain at least one entry")
	}

	children := make([]ForkChild, 0, len(rawSlice))
	for i, v := range rawSlice {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("argv[%d] must be a string, got %T", i, v)
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, fmt.Errorf("argv[%d] cannot be empty", i)
		}
		children = append(children, ForkChild{
			Intent:    s,
			Workspace: ".",
		})
	}
	return children, nil
}

// childProcess holds the state for a single child in a swarm fork.
type childProcess struct {
	cmd      *exec.Cmd
	stdout   *bytes.Buffer // nil for fire-and-forget
	stderr   *bytes.Buffer // nil for fire-and-forget
	intent   string
	index    int
	tapePath string // temp tape file to clean up
}

// childResult captures the outcome of a completed child.
type childResult struct {
	child    *childProcess
	exitCode int
	err      error // non-nil for start/wait failures (not exit code errors)
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

// spawnChild creates and starts a single child process. If capture is true,
// stdout/stderr are captured into buffers; otherwise they are discarded.
// The caller is responsible for removing tapePath when done.
func (f *ForkExecutor) spawnChild(ctx context.Context, child ForkChild, index int, capture bool) (*childProcess, error) {
	tapePath, err := f.copyTapeForChild()
	if err != nil {
		return nil, fmt.Errorf("copy tape for child %d: %w", index, err)
	}

	cmd := exec.CommandContext(ctx, f.QuinePath, child.Intent)
	cmd.Env = append([]string(nil), f.Env...)
	if tapePath != "" {
		cmd.Env = append(cmd.Env, "QUINE_CONTEXT_TAPE="+tapePath)
	}
	if f.WorkspaceEnabled {
		workspace, err := f.resolveChildWorkspace(child.Workspace)
		if err != nil {
			cleanupTapePath(tapePath)
			return nil, fmt.Errorf("resolve child %d workspace: %w", index, err)
		}
		cmd.Env = MergeEnv(cmd.Env, []string{
			"QUINE_WORKSPACE_ROOT=" + f.WorkspaceRoot,
			"QUINE_WORKSPACE=" + workspace,
			"QUINE_WORKSPACE_SESSION=" + f.WorkspaceSession,
			"QUINE_WORKSPACE_OWNER=false",
		})
	} else if strings.TrimSpace(child.Workspace) != "" && strings.TrimSpace(child.Workspace) != "." {
		cleanupTapePath(tapePath)
		return nil, fmt.Errorf("workspace requested for child %d but workspace physics is not enabled", index)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cp := &childProcess{
		cmd:      cmd,
		intent:   child.Intent,
		index:    index,
		tapePath: tapePath,
	}

	if capture {
		cp.stdout = &bytes.Buffer{}
		cp.stderr = &bytes.Buffer{}
		cmd.Stdout = cp.stdout
		cmd.Stderr = cp.stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if err := cmd.Start(); err != nil {
		cleanupTapePath(tapePath)
		return nil, fmt.Errorf("start child %d: %w", index, err)
	}

	return cp, nil
}

func (f *ForkExecutor) resolveChildWorkspace(raw string) (string, error) {
	workspace := strings.TrimSpace(raw)
	if workspace == "" || workspace == "." {
		return f.Workspace, nil
	}

	if !filepath.IsAbs(workspace) {
		workspace = filepath.Join(f.Workspace, workspace)
	}
	workspace = filepath.Clean(workspace)

	resolved, err := canonicalizeRequestedPath(workspace)
	if err != nil {
		return "", err
	}
	root, err := canonicalizeRequestedPath(f.WorkspaceRoot)
	if err != nil {
		return "", err
	}
	if !isPathWithin(root, resolved) {
		return "", fmt.Errorf("workspace %q must be within workspace root %q", resolved, root)
	}
	return resolved, nil
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

func cleanupTapePath(path string) {
	if path != "" {
		_ = os.Remove(path)
	}
}

// cleanupChild removes the child's temp tape file.
func cleanupChild(cp *childProcess) {
	if cp.tapePath != "" {
		os.Remove(cp.tapePath)
	}
}

// waitChild waits for a child to complete and returns its result.
func waitChild(cp *childProcess) childResult {
	err := cp.cmd.Wait()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			err = nil // exit code errors are not failures
		}
	}
	return childResult{
		child:    cp,
		exitCode: exitCode,
		err:      err,
	}
}

// killProcessGroup sends SIGKILL to a child's entire process group.
func killProcessGroup(cp *childProcess) {
	if cp.cmd.Process != nil {
		_ = syscall.Kill(-cp.cmd.Process.Pid, syscall.SIGKILL)
	}
}

// executeGatherAll spawns N children concurrently and waits for ALL to complete.
func (f *ForkExecutor) executeGatherAll(toolID string, req ForkRequest) tape.ToolResult {
	n := len(req.Children)

	ctx, cancel := context.WithTimeout(context.Background(), f.DefaultTimeout)
	defer cancel()

	// Spawn all children
	children := make([]*childProcess, n)
	var spawnErrors []string
	for i, child := range req.Children {
		cp, err := f.spawnChild(ctx, child, i, true)
		if err != nil {
			spawnErrors = append(spawnErrors, fmt.Sprintf("child %d: %v", i, err))
			continue
		}
		children[i] = cp
	}

	// If all children failed to spawn, return immediately
	spawned := 0
	for _, cp := range children {
		if cp != nil {
			spawned++
		}
	}
	if spawned == 0 {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[FORK ERROR] All %d children failed to spawn:\n%s", n, strings.Join(spawnErrors, "\n")),
			IsError: true,
		}
	}

	// Wait for all children concurrently
	results := make([]childResult, n)
	var wg sync.WaitGroup
	for i, cp := range children {
		if cp == nil {
			results[i] = childResult{
				err: fmt.Errorf("failed to spawn"),
			}
			continue
		}
		wg.Add(1)
		go func(idx int, child *childProcess) {
			defer wg.Done()
			results[idx] = waitChild(child)
		}(i, cp)
	}
	wg.Wait()

	// Clean up all temp tape files
	for _, cp := range children {
		if cp != nil {
			cleanupChild(cp)
		}
	}

	// Aggregate results
	allFailed := true
	var sb strings.Builder
	fmt.Fprintf(&sb, "[FORK SWARM] %d children completed\n", n)

	for i := range req.Children {
		if children[i] == nil {
			fmt.Fprintf(&sb, "\n--- CHILD %d [spawn failed]: %q @ %q ---\n", i, req.Children[i].Intent, req.Children[i].Workspace)
			if spawnErrors != nil {
				// Find the relevant error
				for _, se := range spawnErrors {
					if strings.HasPrefix(se, fmt.Sprintf("child %d:", i)) {
						fmt.Fprintf(&sb, "%s\n", se)
					}
				}
			}
			continue
		}

		r := results[i]
		if r.err != nil {
			// Execution error (timeout, etc.)
			fmt.Fprintf(&sb, "\n--- CHILD %d [error]: %q @ %q ---\n%v\n", i, req.Children[i].Intent, req.Children[i].Workspace, r.err)
			continue
		}

		if r.exitCode == 0 {
			allFailed = false
		}

		stdout := f.truncate(r.child.stdout.Bytes())
		stderr := f.truncate(r.child.stderr.Bytes())
		fmt.Fprintf(&sb, "\n--- CHILD %d [exit %d]: %q @ %q ---\n[STDOUT]\n%s\n[STDERR]\n%s\n",
			i, r.exitCode, req.Children[i].Intent, req.Children[i].Workspace, stdout, stderr)
	}

	return tape.ToolResult{
		ToolID:  toolID,
		Content: sb.String(),
		IsError: allFailed,
	}
}

// executeRace spawns N children concurrently. The first child to exit with
// code 0 wins; all remaining children are killed. If all children fail,
// all results are returned.
func (f *ForkExecutor) executeRace(toolID string, req ForkRequest) tape.ToolResult {
	n := len(req.Children)

	ctx, cancel := context.WithTimeout(context.Background(), f.DefaultTimeout)
	defer cancel()

	// Spawn all children
	children := make([]*childProcess, n)
	var spawnErrors []string
	for i, child := range req.Children {
		cp, err := f.spawnChild(ctx, child, i, true)
		if err != nil {
			spawnErrors = append(spawnErrors, fmt.Sprintf("child %d: %v", i, err))
			continue
		}
		children[i] = cp
	}

	// Count successfully spawned
	spawned := 0
	for _, cp := range children {
		if cp != nil {
			spawned++
		}
	}
	if spawned == 0 {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[FORK ERROR] All %d children failed to spawn:\n%s", n, strings.Join(spawnErrors, "\n")),
			IsError: true,
		}
	}

	// Collect results as children finish via a channel
	type indexedResult struct {
		index  int
		result childResult
	}
	resultCh := make(chan indexedResult, spawned)

	for i, cp := range children {
		if cp == nil {
			continue
		}
		go func(idx int, child *childProcess) {
			r := waitChild(child)
			resultCh <- indexedResult{index: idx, result: r}
		}(i, cp)
	}

	// Wait for results, looking for a winner
	var winner *indexedResult
	completed := make(map[int]childResult)
	killed := 0

	for range spawned {
		ir := <-resultCh
		completed[ir.index] = ir.result

		if winner == nil && ir.result.err == nil && ir.result.exitCode == 0 {
			// We have a winner! Kill remaining children.
			winner = &ir
			for i, cp := range children {
				if cp == nil {
					continue
				}
				if _, done := completed[i]; !done {
					killProcessGroup(cp)
					killed++
				}
			}
			// If we've already received all results, break
			if len(completed) == spawned {
				break
			}
		}

		if len(completed) == spawned {
			break
		}
	}

	// Wait for any remaining children that were killed to actually exit
	// (they should exit quickly after SIGKILL)
	for len(completed) < spawned {
		ir := <-resultCh
		completed[ir.index] = ir.result
	}

	// Clean up all temp tape files
	for _, cp := range children {
		if cp != nil {
			cleanupChild(cp)
		}
	}

	// Build result
	if winner != nil {
		r := completed[winner.index]
		cp := children[winner.index]
		stdout := f.truncate(cp.stdout.Bytes())
		stderr := f.truncate(cp.stderr.Bytes())

		succeeded := 0
		for _, cr := range completed {
			if cr.err == nil && cr.exitCode == 0 {
				succeeded++
			}
		}

		content := fmt.Sprintf("[FORK RACE] child %d won (%d spawned, %d succeeded, %d killed)\n\n--- WINNER [exit 0]: %q ---\n[STDOUT]\n%s\n[STDERR]\n%s\n",
			winner.index, spawned, succeeded, killed, r.child.intent, stdout, stderr)

		return tape.ToolResult{
			ToolID:  toolID,
			Content: content,
			IsError: false,
		}
	}

	// All children failed
	var sb strings.Builder
	fmt.Fprintf(&sb, "[FORK RACE] all %d children failed\n", n)

	for i := range req.Children {
		if children[i] == nil {
			fmt.Fprintf(&sb, "\n--- CHILD %d [spawn failed]: %q @ %q ---\n", i, req.Children[i].Intent, req.Children[i].Workspace)
			continue
		}

		r, ok := completed[i]
		if !ok {
			fmt.Fprintf(&sb, "\n--- CHILD %d [no result]: %q @ %q ---\n", i, req.Children[i].Intent, req.Children[i].Workspace)
			continue
		}

		if r.err != nil {
			fmt.Fprintf(&sb, "\n--- CHILD %d [error]: %q @ %q ---\n%v\n", i, req.Children[i].Intent, req.Children[i].Workspace, r.err)
			continue
		}

		stdout := f.truncate(children[i].stdout.Bytes())
		stderr := f.truncate(children[i].stderr.Bytes())
		fmt.Fprintf(&sb, "\n--- CHILD %d [exit %d]: %q @ %q ---\n[STDOUT]\n%s\n[STDERR]\n%s\n",
			i, r.exitCode, req.Children[i].Intent, req.Children[i].Workspace, stdout, stderr)
	}

	return tape.ToolResult{
		ToolID:  toolID,
		Content: sb.String(),
		IsError: true,
	}
}

// executeFireAndForget spawns N children and returns their PIDs immediately.
func (f *ForkExecutor) executeFireAndForget(toolID string, req ForkRequest) tape.ToolResult {
	n := len(req.Children)

	children := make([]*childProcess, 0, n)
	var spawnErrors []string

	for i, child := range req.Children {
		cp, err := f.spawnChild(context.Background(), child, i, false)
		if err != nil {
			spawnErrors = append(spawnErrors, fmt.Sprintf("child %d: %v", i, err))
			continue
		}
		children = append(children, cp)
	}

	if len(children) == 0 {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[FORK ERROR] All %d children failed to spawn:\n%s", n, strings.Join(spawnErrors, "\n")),
			IsError: true,
		}
	}

	// Fire-and-forget: clean up tape files in background after children exit
	for _, cp := range children {
		go func(child *childProcess) {
			child.cmd.Wait()
			cleanupChild(child)
		}(cp)
	}

	// Build result with PIDs
	var sb strings.Builder
	fmt.Fprintf(&sb, "[FORK OK] %d children spawned\n", len(children))
	for _, cp := range children {
		fmt.Fprintf(&sb, "  child %d: PID %d — %q\n", cp.index, cp.cmd.Process.Pid, cp.intent)
	}
	if len(spawnErrors) > 0 {
		fmt.Fprintf(&sb, "\nFailed to spawn %d children:\n", len(spawnErrors))
		for _, e := range spawnErrors {
			fmt.Fprintf(&sb, "  %s\n", e)
		}
	}
	fmt.Fprintf(&sb, "\nChildren are running independently.\nChild tapes will appear under the runtime root: %s", f.DataDir)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: sb.String(),
		IsError: false,
	}
}

// copyTapeForChild copies the current tape file to a temp file that the child
// can read for context. Returns the path to the temp file.
func (f *ForkExecutor) copyTapeForChild() (string, error) {
	// Read current tape
	data, err := os.ReadFile(f.TapePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No tape yet - child starts fresh
			return "", nil
		}
		return "", fmt.Errorf("reading tape: %w", err)
	}

	// Create temp file in the same directory
	tmpFile, err := os.CreateTemp(f.DataDir, "fork-tape-*.jsonl")
	if err != nil {
		return "", fmt.Errorf("creating temp tape: %w", err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("writing temp tape: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("closing temp tape: %w", err)
	}

	return tmpFile.Name(), nil
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
