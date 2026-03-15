package tools

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/kehao95/quine/internal/tape"
)

// MarkRequest is the parsed request for mark().
type MarkRequest struct {
	Summary string
	Fold    bool
}

// UnfoldRequest is the parsed request for unfold().
type UnfoldRequest struct {
	AnchorID int
}

// ParseMarkArgs extracts MarkRequest from tool arguments.
func ParseMarkArgs(args map[string]any) (MarkRequest, error) {
	summary, _ := args["summary"].(string)
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return MarkRequest{}, errors.New("summary is required")
	}
	fold, _ := args["fold"].(bool)
	return MarkRequest{Summary: summary, Fold: fold}, nil
}

// ParseUnfoldArgs extracts UnfoldRequest from tool arguments.
func ParseUnfoldArgs(args map[string]any) (UnfoldRequest, error) {
	raw, ok := args["anchor_id"]
	if !ok {
		return UnfoldRequest{}, errors.New("anchor_id is required")
	}
	id, err := intFromAny(raw)
	if err != nil || id < 0 {
		return UnfoldRequest{}, fmt.Errorf("anchor_id must be a non-negative integer, got %v", raw)
	}
	return UnfoldRequest{AnchorID: id}, nil
}

type memoryState struct {
	Cursor    int `json:"cursor"`
	BaseEntry int `json:"base_entry"`
}

type anchorMeta struct {
	ID      int    `json:"id"`
	Summary string `json:"summary"`
	Anchors []int  `json:"anchors"`
	Turns   []int  `json:"turns"`
}

// AnchorMemoryExecutor manages mark/unfold storage under the runtime context
// surface:
// ${QUINE_AGENT_ROOT}/context/.
type AnchorMemoryExecutor struct {
	AgentRoot string
	TapePath  string
}

// NewAnchorMemoryExecutor creates an anchor-memory executor bound to the
// current session root and current tape file.
func NewAnchorMemoryExecutor(agentRoot, tapePath string) *AnchorMemoryExecutor {
	return &AnchorMemoryExecutor{
		AgentRoot: agentRoot,
		TapePath:  tapePath,
	}
}

func (e *AnchorMemoryExecutor) contextRoot() string { return filepath.Join(e.AgentRoot, "context") }
func (e *AnchorMemoryExecutor) frontierRoot() string {
	return filepath.Join(e.contextRoot(), "frontier")
}
func (e *AnchorMemoryExecutor) anchorsRoot() string { return filepath.Join(e.contextRoot(), "anchors") }
func (e *AnchorMemoryExecutor) turnIDPath() string  { return filepath.Join(e.contextRoot(), "turn.id") }
func (e *AnchorMemoryExecutor) statePath() string {
	return filepath.Join(e.contextRoot(), "state.json")
}
func (e *AnchorMemoryExecutor) currentPath() string {
	return filepath.Join(e.contextRoot(), "current.jsonl")
}

// Mark compresses the current unsummarized tape tail into a new anchor.
func (e *AnchorMemoryExecutor) Mark(toolID string, req MarkRequest) tape.ToolResult {
	state, err := e.syncCurrentFromTape()
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] %v", err),
			IsError: true,
		}
	}

	rawLines, err := readJSONLLines(e.currentPath())
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] read current context: %v", err),
			IsError: true,
		}
	}
	links, err := e.frontierAnchorIDs()
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] read frontier: %v", err),
			IsError: true,
		}
	}

	if !req.Fold && len(rawLines) == 0 {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[MARK ERROR] nothing to mark (no raw context entries)",
			IsError: true,
		}
	}
	if req.Fold && len(rawLines) == 0 && len(links) == 0 {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[MARK ERROR] nothing to mark (no raw context entries and no frontier anchors)",
			IsError: true,
		}
	}

	anchorID, err := e.nextAnchorID()
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] allocate anchor id: %v", err),
			IsError: true,
		}
	}

	stagingDir := filepath.Join(e.anchorsRoot(), fmt.Sprintf(".staging-%d.anchor", anchorID))
	finalDir := filepath.Join(e.anchorsRoot(), fmt.Sprintf("%d.anchor", anchorID))
	if err := os.MkdirAll(stagingDir, 0o755); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] create staging dir: %v", err),
			IsError: true,
		}
	}

	turns := make([]int, len(rawLines))
	for i := range rawLines {
		turns[i] = state.BaseEntry + i
	}
	meta := anchorMeta{
		ID:      anchorID,
		Summary: req.Summary,
		Anchors: []int{},
		Turns:   turns,
	}
	if req.Fold {
		meta.Anchors = append(meta.Anchors, links...)
	}
	if err := writeJSON(stagingDir, "meta.json", meta); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] write meta: %v", err),
			IsError: true,
		}
	}
	if len(rawLines) > 0 {
		if err := writeRawLines(filepath.Join(stagingDir, "raw.jsonl"), rawLines); err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[MARK ERROR] write raw context: %v", err),
				IsError: true,
			}
		}
	}
	if err := os.Rename(stagingDir, finalDir); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] commit anchor: %v", err),
			IsError: true,
		}
	}

	if req.Fold {
		if err := e.removeAllFrontierLinks(); err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[MARK ERROR] clear frontier: %v", err),
				IsError: true,
			}
		}
	}
	linkPath := filepath.Join(e.frontierRoot(), fmt.Sprintf("%d.link", anchorID))
	_ = os.Remove(linkPath)
	if err := os.Symlink(filepath.Join("..", "anchors", fmt.Sprintf("%d.anchor", anchorID)), linkPath); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] create frontier link: %v", err),
			IsError: true,
		}
	}

	// Reset unsummarized context window after successful commit.
	if err := os.WriteFile(e.currentPath(), nil, 0o644); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] reset current context: %v", err),
			IsError: true,
		}
	}
	state.BaseEntry = state.Cursor
	if err := e.saveState(state); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[MARK ERROR] update state: %v", err),
			IsError: true,
		}
	}

	return tape.ToolResult{
		ToolID: toolID,
		Content: fmt.Sprintf("[MARK] anchor=%d fold=%t raw_entries=%d absorbed=%d summary=%q",
			anchorID, req.Fold, len(rawLines), len(meta.Anchors), req.Summary),
	}
}

// Unfold returns one-level structured view for an anchor.
func (e *AnchorMemoryExecutor) Unfold(toolID string, req UnfoldRequest) tape.ToolResult {
	if _, err := e.syncCurrentFromTape(); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[UNFOLD ERROR] %v", err),
			IsError: true,
		}
	}
	meta, err := e.readAnchorMeta(req.AnchorID)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[UNFOLD ERROR] %v", err),
			IsError: true,
		}
	}

	children := make([]map[string]any, 0, len(meta.Anchors))
	for _, childID := range meta.Anchors {
		childMeta, err := e.readAnchorMeta(childID)
		if err != nil {
			continue
		}
		children = append(children, map[string]any{
			"id":      childMeta.ID,
			"summary": childMeta.Summary,
		})
	}

	turnItems, err := e.readAnchorRaw(req.AnchorID)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[UNFOLD ERROR] %v", err),
			IsError: true,
		}
	}

	view := map[string]any{
		"id":      meta.ID,
		"summary": meta.Summary,
		"anchors": children,
		"turns":   turnItems,
	}
	data, _ := json.MarshalIndent(view, "", "  ")
	return tape.ToolResult{
		ToolID:  toolID,
		Content: "[UNFOLD]\n" + string(data),
	}
}

// FeedbackBlock returns compact memory feedback suitable for tool-result suffixes.
func (e *AnchorMemoryExecutor) FeedbackBlock() (string, error) {
	state, err := e.syncCurrentFromTape()
	if err != nil {
		return "", err
	}
	rawCount, err := countJSONLLines(e.currentPath())
	if err != nil {
		return "", err
	}
	links, err := e.frontierAnchorIDs()
	if err != nil {
		return "", err
	}

	pending := map[string]any{
		"anchors": links,
	}
	if rawCount > 0 {
		pending["mark"] = map[string]any{
			"first_entry": state.BaseEntry,
			"last_entry":  state.BaseEntry + rawCount - 1,
		}
	}
	out := map[string]any{
		"pending": pending,
	}
	data, _ := json.Marshal(out)
	return "[MEMORY META] " + string(data), nil
}

func (e *AnchorMemoryExecutor) syncCurrentFromTape() (memoryState, error) {
	if err := e.ensureLayout(); err != nil {
		return memoryState{}, err
	}
	state, err := e.loadState()
	if err != nil {
		return memoryState{}, err
	}

	summary, err := tape.ReadTapeFile(e.TapePath)
	if err != nil {
		// Tape may not exist during early bootstrap.
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		return memoryState{}, err
	}
	entries := tape.ContextEntries(summary.Entries)
	if state.Cursor > len(entries) {
		state.Cursor = len(entries)
		state.BaseEntry = state.Cursor
		if err := os.WriteFile(e.currentPath(), nil, 0o644); err != nil {
			return memoryState{}, err
		}
		return state, e.saveState(state)
	}
	if state.Cursor == len(entries) {
		return state, nil
	}

	existingCount, err := countJSONLLines(e.currentPath())
	if err != nil {
		return memoryState{}, err
	}
	if existingCount == 0 {
		state.BaseEntry = state.Cursor
	}

	f, err := os.OpenFile(e.currentPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return memoryState{}, err
	}
	defer f.Close()
	for _, entry := range entries[state.Cursor:] {
		line, _ := json.Marshal(entry)
		if _, err := f.Write(line); err != nil {
			return memoryState{}, err
		}
		if _, err := f.Write([]byte("\n")); err != nil {
			return memoryState{}, err
		}
	}
	state.Cursor = len(entries)
	if err := e.saveState(state); err != nil {
		return memoryState{}, err
	}
	return state, nil
}

func (e *AnchorMemoryExecutor) ensureLayout() error {
	for _, dir := range []string{e.contextRoot(), e.frontierRoot(), e.anchorsRoot()} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(e.turnIDPath()); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(e.turnIDPath(), []byte("0\n"), 0o644); err != nil {
			return err
		}
	}
	if _, err := os.Stat(e.statePath()); errors.Is(err, os.ErrNotExist) {
		if err := writeJSON(filepath.Dir(e.statePath()), filepath.Base(e.statePath()), memoryState{}); err != nil {
			return err
		}
	}
	if _, err := os.Stat(e.currentPath()); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(e.currentPath(), nil, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func (e *AnchorMemoryExecutor) loadState() (memoryState, error) {
	var st memoryState
	data, err := os.ReadFile(e.statePath())
	if err != nil {
		return st, err
	}
	if len(data) == 0 {
		return st, nil
	}
	if err := json.Unmarshal(data, &st); err != nil {
		return st, err
	}
	return st, nil
}

func (e *AnchorMemoryExecutor) saveState(st memoryState) error {
	return writeJSON(filepath.Dir(e.statePath()), filepath.Base(e.statePath()), st)
}

func (e *AnchorMemoryExecutor) nextAnchorID() (int, error) {
	data, err := os.ReadFile(e.turnIDPath())
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(data))
	if s == "" {
		s = "0"
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	next := strconv.Itoa(id + 1)
	if err := os.WriteFile(e.turnIDPath(), []byte(next+"\n"), 0o644); err != nil {
		return 0, err
	}
	return id, nil
}

func anchorIDsFromDirEntries(entries []os.DirEntry) []int {
	var ids []int
	for _, de := range entries {
		name := de.Name()
		if !strings.HasSuffix(name, ".link") {
			continue
		}
		idStr := strings.TrimSuffix(name, ".link")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}

func (e *AnchorMemoryExecutor) frontierEntries() ([]os.DirEntry, error) {
	entries, err := os.ReadDir(e.frontierRoot())
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (e *AnchorMemoryExecutor) frontierAnchorIDs() ([]int, error) {
	entries, err := e.frontierEntries()
	if err != nil {
		return nil, err
	}
	return anchorIDsFromDirEntries(entries), nil
}

func (e *AnchorMemoryExecutor) removeAllFrontierLinks() error {
	entries, err := e.frontierEntries()
	if err != nil {
		return err
	}
	for _, de := range entries {
		if strings.HasSuffix(de.Name(), ".link") {
			if err := os.Remove(filepath.Join(e.frontierRoot(), de.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *AnchorMemoryExecutor) readAnchorMeta(anchorID int) (anchorMeta, error) {
	var m anchorMeta
	path := filepath.Join(e.anchorsRoot(), fmt.Sprintf("%d.anchor", anchorID), "meta.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return m, fmt.Errorf("anchor %d not found", anchorID)
		}
		return m, err
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, err
	}
	return m, nil
}

func (e *AnchorMemoryExecutor) readAnchorRaw(anchorID int) ([]any, error) {
	path := filepath.Join(e.anchorsRoot(), fmt.Sprintf("%d.anchor", anchorID), "raw.jsonl")
	lines, err := readJSONLLines(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]any, 0, len(lines))
	for _, line := range lines {
		var v any
		if err := json.Unmarshal(line, &v); err != nil {
			out = append(out, map[string]any{
				"raw":          string(line),
				"decode_error": err.Error(),
			})
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

func intFromAny(v any) (int, error) {
	switch x := v.(type) {
	case int:
		return x, nil
	case int64:
		return int(x), nil
	case float64:
		if x != float64(int(x)) {
			return 0, errors.New("not an integer")
		}
		return int(x), nil
	case json.Number:
		n, err := x.Int64()
		return int(n), err
	case string:
		return strconv.Atoi(strings.TrimSpace(x))
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

func readJSONLLines(path string) ([][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	var lines [][]byte
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func countJSONLLines(path string) (int, error) {
	lines, err := readJSONLLines(path)
	if err != nil {
		return 0, err
	}
	return len(lines), nil
}

func writeRawLines(path string, lines [][]byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range lines {
		if _, err := f.Write(line); err != nil {
			return err
		}
		if _, err := f.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(dir, name string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, name), data, 0o644)
}
