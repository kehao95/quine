package tools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/kehao95/quine/internal/tape"
)

// anchorStateSchemaDoc is written to context/state/SCHEMA.md so the layout is
// self-describing on disk: any agent — including a fresh peer process handed
// this root for memory maintenance — can read it and write canonical anchors
// without reverse-engineering the format. It mirrors the on-disk contract the
// runtime actually reads (anchorMeta, frontier links, current.jsonl).
const anchorStateSchemaDoc = `# context/state — self-describing layout

This directory is your cognition as plain files. Every turn the runtime
reprojects your context window from it: window = system prompt + one
"[ANCHOR <id>] <resolution>" line per frontier link (sorted by id) + the raw
turns in current.jsonl. Editing these files changes what you think next turn.
The mark/unfold tools, when present, are convenience moves over this same
substrate; everything here is also editable directly with sh — by you, or by a
peer process you hand this directory's parent ($QUINE_AGENT_ROOT) to.

## Files

- current.jsonl — the live raw turn stream, one tape entry (JSON) per line.
  It is [settled prefix | live tail]:
    * settled prefix = complete past turns; compactible.
    * live tail = the current open turn (the last line, e.g. an unanswered
      tool call). NEVER rewrite it or the provider tool-call contract breaks.
- anchors/<id>.anchor/ — immutable crystallized memory:
    * meta.json = {"id": <int>, "resolution": "<one low-entropy summary>",
      "anchors": [<absorbed child anchor ids>], "turns": [<global turn indices
      whose raw lines are stored here>]}
    * raw.jsonl = the exact raw tape entries this anchor crystallized. This is
      the recoverable floor; never delete it.
- frontier/<id>.link — a relative symlink to ../anchors/<id>.anchor. The set of
  frontier links is the ACTIVE root set: only these anchors render as
  "[ANCHOR id]" summary lines in the window.
- state.json = {"cursor": <int>, "base_entry": <int>} — base_entry is the global
  index of the first line currently in current.jsonl.
- turn.id — a monotonic counter holding the next anchor id. Read it, use that
  id, then write back id+1.

## How to compact (promote settled history into an anchor)

1. id = read turn.id; write turn.id = id+1.
2. Split current.jsonl into the live tail (keep the last line aside) and the
   settled prefix (everything before it).
3. Write anchors/<id>.anchor/raw.jsonl = the settled lines, verbatim.
4. Write anchors/<id>.anchor/meta.json = {id, resolution: a faithful summary of
   the settled lines, anchors: [<frontier ids you are absorbing, else []>],
   turns: [<global indices of the settled lines>]}. Commit atomically: write
   into a staging dir then rename, so a half-written anchor is never observed.
5. Rewrite current.jsonl to contain only the live tail.
6. Create frontier/<id>.link -> ../anchors/<id>.anchor. If this new anchor
   subsumes earlier frontier anchors, list their ids in meta.anchors and remove
   their frontier/*.link (they remain reachable via meta.anchors).

## Invariants

- Raw under anchors/*/raw.jsonl is the floor: reorganize the view, never destroy
  history.
- Never rewrite the live tail of current.jsonl.
- Anchor ids are monotonic via turn.id.
- A frontier holding one governing anchor (with the rest absorbed under it via
  meta.anchors) is the compacted/folded shape.
`

// MarkRequest is the parsed request for mark().
type MarkRequest struct {
	Resolution string
	Fold       bool
}

// UnfoldRequest is the parsed request for unfold().
type UnfoldRequest struct {
	AnchorID int
}

// ParseMarkArgs extracts MarkRequest from tool arguments.
func ParseMarkArgs(args map[string]any) (MarkRequest, error) {
	resolution, _ := args["resolution"].(string)
	resolution = strings.TrimSpace(resolution)
	if resolution == "" {
		return MarkRequest{}, errors.New("resolution is required")
	}
	fold, _, err := BoolArg(args, "fold")
	if err != nil {
		return MarkRequest{}, fmt.Errorf("fold must be a boolean: %w", err)
	}
	return MarkRequest{Resolution: resolution, Fold: fold}, nil
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
	ID         int    `json:"id"`
	Resolution string `json:"resolution"`
	Anchors    []int  `json:"anchors"`
	Turns      []int  `json:"turns"`
}

type markStructuredResult struct {
	Tool              string `json:"tool"`
	Status            string `json:"status"`
	AnchorID          int    `json:"anchor_id"`
	Fold              bool   `json:"fold"`
	RawEntries        int    `json:"raw_entries"`
	Absorbed          int    `json:"absorbed"`
	AbsorbedAnchorIDs []int  `json:"absorbed_anchor_ids"`
	Resolution        string `json:"resolution,omitempty"`
	Note              string `json:"note,omitempty"`
	Error             string `json:"error,omitempty"`
}

type unfoldChildSummary struct {
	ID         int    `json:"id"`
	Resolution string `json:"resolution"`
}

type unfoldStructuredResult struct {
	Tool       string               `json:"tool"`
	Status     string               `json:"status"`
	AnchorID   int                  `json:"anchor_id,omitempty"`
	Resolution string               `json:"resolution,omitempty"`
	Anchors    []unfoldChildSummary `json:"anchors,omitempty"`
	Turns      []any                `json:"turns,omitempty"`
	Error      string               `json:"error,omitempty"`
}

// AnchorMemoryExecutor manages mark/unfold storage under the runtime state
// surface:
// ${QUINE_AGENT_ROOT}/context/state/.
type AnchorMemoryExecutor struct {
	AgentRoot    string
	WarnTokens   int
	DangerTokens int
	DeathTokens  int
	// MarkDisabled / FoldDisabled make the pressure telemetry coherent with the
	// available tool surface. Zero value (false) keeps the default mark/fold
	// advice; when a gate removes a tool, the matching field is set so the
	// telemetry stops advising a tool the agent cannot call and instead frames
	// the goal (reduce live working memory), leaving the "how" to the prompt.
	MarkDisabled        bool
	FoldDisabled        bool
	MemoryStrategyHints bool
}

type MemoryDeathStatus struct {
	FrontierEstimatedTokens int64
	DeathTokens             int
	Exceeded                bool
}

// NewAnchorMemoryExecutor creates an anchor-memory executor bound to the
// current session root.
func NewAnchorMemoryExecutor(agentRoot string, warnTokens, dangerTokens int, deathTokens ...int) *AnchorMemoryExecutor {
	if warnTokens <= 0 {
		warnTokens = 8000
	}
	if dangerTokens <= warnTokens {
		dangerTokens = warnTokens * 2
	}
	death := 0
	if len(deathTokens) > 0 {
		death = deathTokens[0]
	}
	if death < 0 {
		death = 0
	}
	return &AnchorMemoryExecutor{
		AgentRoot:           agentRoot,
		WarnTokens:          warnTokens,
		DangerTokens:        dangerTokens,
		DeathTokens:         death,
		MemoryStrategyHints: true,
	}
}

func (e *AnchorMemoryExecutor) stateRoot() string {
	return filepath.Join(e.AgentRoot, "context", "state")
}
func (e *AnchorMemoryExecutor) frontierRoot() string {
	return filepath.Join(e.stateRoot(), "frontier")
}
func (e *AnchorMemoryExecutor) anchorsRoot() string { return filepath.Join(e.stateRoot(), "anchors") }
func (e *AnchorMemoryExecutor) turnIDPath() string  { return filepath.Join(e.stateRoot(), "turn.id") }
func (e *AnchorMemoryExecutor) schemaPath() string  { return filepath.Join(e.stateRoot(), "SCHEMA.md") }
func (e *AnchorMemoryExecutor) statePath() string {
	return filepath.Join(e.stateRoot(), "state.json")
}
func (e *AnchorMemoryExecutor) currentPath() string {
	return filepath.Join(e.stateRoot(), "current.jsonl")
}

// AppendEntry appends one raw context entry to context/state/current.jsonl and
// advances the logical cursor used by mark/fold bookkeeping.
func (e *AnchorMemoryExecutor) AppendEntry(entry tape.TapeEntry) error {
	if err := e.ensureLayout(); err != nil {
		return err
	}
	st, err := e.loadState()
	if err != nil {
		return err
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(e.currentPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(line); err != nil {
		return err
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return err
	}
	st.Cursor++
	return e.saveState(st)
}

// Mark compresses the current unsummarized raw context frontier into a new anchor.
func (e *AnchorMemoryExecutor) Mark(toolID string, req MarkRequest) tape.ToolResult {
	if err := e.ensureLayout(); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] %v", err),
			}),
			IsError: true,
		}
	}
	state, err := e.loadState()
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] %v", err),
			}),
			IsError: true,
		}
	}

	rawLines, err := readJSONLLines(e.currentPath())
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] read current context: %v", err),
			}),
			IsError: true,
		}
	}
	links, err := e.frontierAnchorIDs()
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] read frontier: %v", err),
			}),
			IsError: true,
		}
	}

	if !req.Fold && len(rawLines) == 0 {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  "[MARK ERROR] nothing to mark (no raw context entries)",
			}),
			IsError: true,
		}
	}
	if req.Fold && len(rawLines) == 0 && len(links) == 0 {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  "[MARK ERROR] nothing to mark (no raw context entries and no frontier anchors)",
			}),
			IsError: true,
		}
	}

	anchorID, err := e.nextAnchorID()
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] allocate anchor id: %v", err),
			}),
			IsError: true,
		}
	}

	stagingDir := filepath.Join(e.anchorsRoot(), fmt.Sprintf(".staging-%d.anchor", anchorID))
	finalDir := filepath.Join(e.anchorsRoot(), fmt.Sprintf("%d.anchor", anchorID))
	if err := os.MkdirAll(stagingDir, 0o755); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] create staging dir: %v", err),
			}),
			IsError: true,
		}
	}

	turns := make([]int, len(rawLines))
	for i := range rawLines {
		turns[i] = state.BaseEntry + i
	}
	meta := anchorMeta{
		ID:         anchorID,
		Resolution: req.Resolution,
		Anchors:    []int{},
		Turns:      turns,
	}
	if req.Fold {
		meta.Anchors = append(meta.Anchors, links...)
	}
	if err := writeJSON(stagingDir, "meta.json", meta); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] write meta: %v", err),
			}),
			IsError: true,
		}
	}
	if len(rawLines) > 0 {
		if err := writeRawLines(filepath.Join(stagingDir, "raw.jsonl"), rawLines); err != nil {
			return tape.ToolResult{
				ToolID: toolID,
				Content: tape.MarshalToolResultContent(markStructuredResult{
					Tool:   "mark",
					Status: "error",
					Error:  fmt.Sprintf("[MARK ERROR] write raw context: %v", err),
				}),
				IsError: true,
			}
		}
	}
	if err := os.Rename(stagingDir, finalDir); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] commit anchor: %v", err),
			}),
			IsError: true,
		}
	}

	if req.Fold {
		if err := e.removeAllFrontierLinks(); err != nil {
			return tape.ToolResult{
				ToolID: toolID,
				Content: tape.MarshalToolResultContent(markStructuredResult{
					Tool:   "mark",
					Status: "error",
					Error:  fmt.Sprintf("[MARK ERROR] clear frontier: %v", err),
				}),
				IsError: true,
			}
		}
	}
	linkPath := filepath.Join(e.frontierRoot(), fmt.Sprintf("%d.link", anchorID))
	_ = os.Remove(linkPath)
	if err := os.Symlink(filepath.Join("..", "anchors", fmt.Sprintf("%d.anchor", anchorID)), linkPath); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] create frontier link: %v", err),
			}),
			IsError: true,
		}
	}

	// Reset unsummarized context window after successful commit.
	if err := os.WriteFile(e.currentPath(), nil, 0o644); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] reset current context: %v", err),
			}),
			IsError: true,
		}
	}
	state.BaseEntry = state.Cursor
	if err := e.saveState(state); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(markStructuredResult{
				Tool:   "mark",
				Status: "error",
				Error:  fmt.Sprintf("[MARK ERROR] update state: %v", err),
			}),
			IsError: true,
		}
	}

	note := ""
	if req.Fold && len(meta.Anchors) == 0 {
		note = "fold requested with no prior anchors; no consolidation occurred. The new anchor was recorded, but this was not an effective fold."
	}
	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(markStructuredResult{
			Tool:              "mark",
			Status:            "completed",
			AnchorID:          anchorID,
			Fold:              req.Fold,
			RawEntries:        len(rawLines),
			Absorbed:          len(meta.Anchors),
			AbsorbedAnchorIDs: append([]int{}, meta.Anchors...),
			Resolution:        req.Resolution,
			Note:              note,
		}),
	}
}

// Unfold returns one-level structured view for an anchor.
func (e *AnchorMemoryExecutor) Unfold(toolID string, req UnfoldRequest) tape.ToolResult {
	if err := e.ensureLayout(); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(unfoldStructuredResult{
				Tool:   "unfold",
				Status: "error",
				Error:  fmt.Sprintf("[UNFOLD ERROR] %v", err),
			}),
			IsError: true,
		}
	}
	meta, err := e.readAnchorMeta(req.AnchorID)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(unfoldStructuredResult{
				Tool:     "unfold",
				Status:   "error",
				AnchorID: req.AnchorID,
				Error:    fmt.Sprintf("[UNFOLD ERROR] %v", err),
			}),
			IsError: true,
		}
	}

	structuredChildren := make([]unfoldChildSummary, 0, len(meta.Anchors))
	for _, childID := range meta.Anchors {
		childMeta, err := e.readAnchorMeta(childID)
		if err != nil {
			continue
		}
		structuredChildren = append(structuredChildren, unfoldChildSummary{
			ID:         childMeta.ID,
			Resolution: childMeta.Resolution,
		})
	}

	turnItems, err := e.readAnchorRaw(req.AnchorID)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(unfoldStructuredResult{
				Tool:     "unfold",
				Status:   "error",
				AnchorID: req.AnchorID,
				Error:    fmt.Sprintf("[UNFOLD ERROR] %v", err),
			}),
			IsError: true,
		}
	}

	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(unfoldStructuredResult{
			Tool:       "unfold",
			Status:     "completed",
			AnchorID:   meta.ID,
			Resolution: meta.Resolution,
			Anchors:    structuredChildren,
			Turns:      turnItems,
		}),
	}
}

func (e *AnchorMemoryExecutor) statusData(extraEntries ...tape.TapeEntry) (map[string]any, map[string]any, error) {
	if err := e.ensureLayout(); err != nil {
		return nil, nil, err
	}
	state, err := e.loadState()
	if err != nil {
		return nil, nil, err
	}
	rawCount, rawBytes, err := contextTelemetryStats(e.currentPath(), extraEntries...)
	if err != nil {
		return nil, nil, err
	}
	links, err := e.frontierAnchorIDs()
	if err != nil {
		return nil, nil, err
	}
	totalAnchors, err := e.anchorCount()
	if err != nil {
		return nil, nil, err
	}
	rawTokens := estimateTokens(rawBytes)

	ratio := any(nil)
	if totalAnchors > 0 {
		ratio = float64(rawCount) / float64(totalAnchors)
	}
	byteRatio := any(nil)
	if totalAnchors > 0 {
		byteRatio = float64(rawBytes) / float64(totalAnchors)
	}

	topology := "seed"
	switch {
	case len(links) > 1:
		topology = "flat"
	case len(links) == 1 && totalAnchors > len(links):
		topology = "layered"
	case len(links) == 1:
		topology = "forming"
	case totalAnchors > 0:
		topology = "backgrounded"
	}

	pending := map[string]any{
		"anchors": links,
	}
	if rawCount > 0 {
		pending["mark"] = map[string]any{
			"first_entry":      state.BaseEntry,
			"last_entry":       state.BaseEntry + rawCount - 1,
			"raw_bytes":        rawBytes,
			"estimated_tokens": rawTokens,
		}
	}
	topologyData := map[string]any{
		"frontier_entries":               rawCount,
		"frontier_bytes":                 rawBytes,
		"frontier_estimated_tokens":      rawTokens,
		"frontier_anchors":               len(links),
		"governing_anchor":               governingAnchor(links),
		"anchor_count":                   totalAnchors,
		"frontier_to_anchor_ratio":       ratio,
		"frontier_bytes_to_anchor_ratio": byteRatio,
		"topology":                       topology,
		"pending":                        pending,
	}
	statusData := map[string]any{
		"level":                     e.memoryLevel(rawTokens),
		"pressure":                  e.pressure(rawTokens),
		"frontier_level":            e.memoryLevel(rawTokens),
		"frontier_pressure":         e.pressure(rawTokens),
		"frontier_estimated_tokens": rawTokens,
		"anchor_level":              e.anchorLevel(len(links), topology),
		"anchor_pressure":           e.anchorPressure(len(links), topology),
		"frontier_anchors":          len(links),
		"governing_anchor":          governingAnchor(links),
		"topology":                  topology,
		"warn_tokens":               e.WarnTokens,
		"danger_tokens":             e.DangerTokens,
	}
	if e.DeathTokens > 0 {
		statusData["death_tokens"] = e.DeathTokens
		statusData["death_pressure"] = e.deathPressure(rawTokens)
	}
	if nextAction, nextTiming := e.frontierNextAction(rawTokens, len(links), topology); nextAction != "" {
		statusData["next_action"] = nextAction
		statusData["next_action_timing"] = nextTiming
	} else if nextAction, nextTiming := e.anchorNextAction(rawTokens, len(links), topology); nextAction != "" {
		statusData["next_action"] = nextAction
		statusData["next_action_timing"] = nextTiming
	}
	if warning := e.frontierWarning(rawTokens, len(links), topology); warning != "" {
		statusData["warning"] = warning
		statusData["frontier_warning"] = warning
	}
	if anchorWarning := e.anchorWarning(len(links), topology); anchorWarning != "" {
		statusData["anchor_warning"] = anchorWarning
	}
	return statusData, topologyData, nil
}

// FeedbackBlock returns compact memory feedback suitable for tool-result suffixes.
func (e *AnchorMemoryExecutor) FeedbackBlock() (string, error) {
	status, topology, err := e.statusData()
	if err != nil {
		return "", err
	}
	meta := make(map[string]any, len(topology)+4)
	for k, v := range topology {
		meta[k] = v
	}
	meta["status_level"] = status["level"]
	meta["status_pressure"] = status["pressure"]
	meta["warn_tokens"] = status["warn_tokens"]
	meta["danger_tokens"] = status["danger_tokens"]
	if warning, ok := status["warning"]; ok {
		meta["status_warning"] = warning
	}
	data, _ := json.Marshal(meta)
	return "[MEMORY META] " + string(data), nil
}

// StatusBlock returns ongoing working-memory pressure telemetry suitable for
// appending to ordinary tool results.
func (e *AnchorMemoryExecutor) StatusBlock() (string, error) {
	return e.StatusBlockWithEntries()
}

func (e *AnchorMemoryExecutor) StatusBlockWithEntries(extraEntries ...tape.TapeEntry) (string, error) {
	status, _, err := e.statusData(extraEntries...)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(status)
	return "[MEMORY STATUS] " + string(data), nil
}

func (e *AnchorMemoryExecutor) DeathStatusWithEntries(extraEntries ...tape.TapeEntry) (MemoryDeathStatus, error) {
	if e == nil || e.DeathTokens <= 0 {
		return MemoryDeathStatus{}, nil
	}
	status, _, err := e.statusData(extraEntries...)
	if err != nil {
		return MemoryDeathStatus{}, err
	}
	tokens, _ := status["frontier_estimated_tokens"].(int64)
	return MemoryDeathStatus{
		FrontierEstimatedTokens: tokens,
		DeathTokens:             e.DeathTokens,
		Exceeded:                tokens >= int64(e.DeathTokens),
	}, nil
}

// ActionBlock returns a compact imperative action hint when memory telemetry has
// an explicit next step.
func (e *AnchorMemoryExecutor) ActionBlock() (string, error) {
	return e.ActionBlockWithEntries()
}

func (e *AnchorMemoryExecutor) ActionBlockWithEntries(extraEntries ...tape.TapeEntry) (string, error) {
	status, _, err := e.statusData(extraEntries...)
	if err != nil {
		return "", err
	}
	action, _ := status["next_action"].(string)
	if strings.TrimSpace(action) == "" {
		return "", nil
	}
	timing, _ := status["next_action_timing"].(string)
	if strings.TrimSpace(timing) == "" {
		return "", nil
	}
	return fmt.Sprintf("[MEMORY ACTION] next_action=%s next_action_timing=%s", action, timing), nil
}

// WarningBlock returns the human-readable warning surface for elevated
// memory pressure.
func (e *AnchorMemoryExecutor) WarningBlock() (string, error) {
	return e.WarningBlockWithEntries()
}

func (e *AnchorMemoryExecutor) WarningBlockWithEntries(extraEntries ...tape.TapeEntry) (string, error) {
	status, _, err := e.statusData(extraEntries...)
	if err != nil {
		return "", err
	}
	warning, _ := status["warning"].(string)
	if strings.TrimSpace(warning) == "" {
		return "", nil
	}
	return "[MEMORY WARNING] " + warning, nil
}

// TopologyBlock returns the fuller topology surface for elevated-pressure or
// debug-facing disclosure.
func (e *AnchorMemoryExecutor) TopologyBlock() (string, error) {
	return e.TopologyBlockWithEntries()
}

func (e *AnchorMemoryExecutor) TopologyBlockWithEntries(extraEntries ...tape.TapeEntry) (string, error) {
	_, topology, err := e.statusData(extraEntries...)
	if err != nil {
		return "", err
	}
	data, _ := json.Marshal(topology)
	return "[MEMORY TOPOLOGY] " + string(data), nil
}

func (e *AnchorMemoryExecutor) ShouldExposeTopology() (bool, error) {
	return e.ShouldExposeTopologyWithEntries()
}

func (e *AnchorMemoryExecutor) ShouldExposeTopologyWithEntries(extraEntries ...tape.TapeEntry) (bool, error) {
	status, topology, err := e.statusData(extraEntries...)
	if err != nil {
		return false, err
	}
	level, _ := status["level"].(string)
	if level != "ok" {
		return true, nil
	}
	if topology["topology"] != "seed" {
		return true, nil
	}
	return false, nil
}

func (e *AnchorMemoryExecutor) memoryLevel(rawTokens int64) string {
	switch {
	case e.DeathTokens > 0 && rawTokens >= int64(e.DeathTokens):
		return "death"
	case rawTokens >= int64(e.DangerTokens):
		return "danger"
	case rawTokens >= int64(e.WarnTokens):
		return "warn"
	default:
		return "ok"
	}
}

func (e *AnchorMemoryExecutor) pressure(rawTokens int64) float64 {
	if e.DangerTokens <= 0 {
		return 0
	}
	value := float64(rawTokens) / float64(e.DangerTokens)
	return math.Round(value*100) / 100
}

func (e *AnchorMemoryExecutor) deathPressure(rawTokens int64) float64 {
	if e.DeathTokens <= 0 {
		return 0
	}
	value := float64(rawTokens) / float64(e.DeathTokens)
	return math.Round(value*100) / 100
}

func (e *AnchorMemoryExecutor) frontierWarning(rawTokens int64, frontierAnchors int, topology string) string {
	if e.MarkDisabled {
		if !e.MemoryStrategyHints {
			switch e.memoryLevel(rawTokens) {
			case "warn":
				return "raw working memory is entering degradation range; reduce live context before further context-growing work"
			case "danger":
				return "raw working memory is beyond reliable range; reduce live context now or uncrystallized findings may become unstable"
			case "death":
				return "raw working memory crossed the configured death cutoff; this incarnation will terminate"
			default:
				return ""
			}
		}
		// mark is gated off: do not advise a tool the agent cannot call. Frame
		// the goal; the prompt owns the available ways to reduce live context.
		switch e.memoryLevel(rawTokens) {
		case "warn":
			return "raw working memory is entering degradation range; before the next `sh`, manage your context: your settled `context/state` can be reorganized and compacted in place or by a peer process you start, while findings are still stable"
		case "danger":
			return "raw working memory is beyond reliable range; before the next `sh`, manage your context now — hand your `context/state` to a peer process to reorganize and compact it, or trim it yourself; uncrystallized findings are otherwise unstable"
		case "death":
			return "raw working memory crossed the configured death cutoff; this incarnation will terminate"
		default:
			return ""
		}
	}
	switch e.memoryLevel(rawTokens) {
	case "warn":
		if frontierAnchors > 0 && topology == "forming" {
			return "call `mark` before the next `sh`: raw frontier is entering working-memory degradation; checkpoint the newly stabilized local structure with a plain mark, not fold, unless a higher-order synthesis has already formed"
		}
		return "if a local result is already stable, call `mark` before the next `sh`: raw frontier is entering working-memory degradation and crystallized boundaries are more reliable than unresolved traces"
	case "danger":
		if frontierAnchors > 0 && topology == "forming" {
			return "call `mark` before the next `sh`: raw frontier is beyond reliable working-memory range; take another plain mark for the new local structure unless a higher-order resolution now makes older anchors mere background"
		}
		return "if a local result is already stable, call `mark` before the next `sh`: raw frontier is beyond reliable working-memory range and uncrystallized findings are now unstable"
	case "death":
		return "raw working memory crossed the configured death cutoff; this incarnation will terminate"
	default:
		return ""
	}
}

func (e *AnchorMemoryExecutor) frontierNextAction(rawTokens int64, frontierAnchors int, topology string) (string, string) {
	switch e.memoryLevel(rawTokens) {
	case "warn", "danger":
		if e.MarkDisabled {
			// Goal token, not a tool name: the agent reduces live context by
			// whatever means the prompt makes available.
			return "reduce_context", "before_next_sh"
		}
		return "mark", "before_next_sh"
	}
	return "", ""
}

func (e *AnchorMemoryExecutor) anchorLevel(frontierAnchors int, topology string) string {
	if topology != "flat" {
		return "ok"
	}
	switch {
	case frontierAnchors >= 3:
		return "danger"
	case frontierAnchors >= 2:
		return "warn"
	default:
		return "ok"
	}
}

func (e *AnchorMemoryExecutor) anchorPressure(frontierAnchors int, topology string) float64 {
	if topology != "flat" || frontierAnchors <= 1 {
		return 0
	}
	value := float64(frontierAnchors-1) / 2.0
	return math.Round(value*100) / 100
}

func governingAnchor(links []int) any {
	if len(links) != 1 {
		return nil
	}
	return links[0]
}

func (e *AnchorMemoryExecutor) anchorWarning(frontierAnchors int, topology string) string {
	if e.FoldDisabled {
		// fold is gated off: do not advise it. Parallel-anchor pressure rarely
		// arises when the agent cannot mark, but stay coherent if it does.
		return ""
	}
	switch e.anchorLevel(frontierAnchors, topology) {
	case "warn":
		return "parallel working anchors are beginning to accumulate; two active anchors alone do not force a fold. Keep using plain `mark` for new stable structure, and call `mark` with `fold=true` only later if a newer higher-order resolution makes earlier anchors remembered background and reconfigures the active frontier around one governing anchor, including when several returned child findings are now being compressed into one parent synthesis"
	case "danger":
		return "parallel working anchors are accumulating; if your newest resolution now makes some earlier anchors remembered background and reconfigures the active frontier around one governing anchor, call `mark` with `fold=true` before the next `sh`; this includes the case where several stable child findings have returned and are now being compressed into one parent-level conclusion"
	}
	return ""
}

func (e *AnchorMemoryExecutor) anchorNextAction(rawTokens int64, frontierAnchors int, topology string) (string, string) {
	if e.FoldDisabled {
		return "", ""
	}
	if e.memoryLevel(rawTokens) != "ok" {
		return "", ""
	}
	switch e.anchorLevel(frontierAnchors, topology) {
	case "danger":
		return "fold", "before_next_sh"
	default:
		return "", ""
	}
}

func (e *AnchorMemoryExecutor) ensureLayout() error {
	for _, dir := range []string{e.stateRoot(), e.frontierRoot(), e.anchorsRoot()} {
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
	// Self-describing layout: any agent (or a peer handed this root) can read
	// SCHEMA.md to write canonical anchors without reverse-engineering. Keep it
	// current with the runtime by rewriting when absent or stale.
	if data, err := os.ReadFile(e.schemaPath()); errors.Is(err, os.ErrNotExist) || (err == nil && string(data) != anchorStateSchemaDoc) {
		if werr := os.WriteFile(e.schemaPath(), []byte(anchorStateSchemaDoc), 0o644); werr != nil {
			return werr
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

func (e *AnchorMemoryExecutor) anchorCount() (int, error) {
	entries, err := os.ReadDir(e.anchorsRoot())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	count := 0
	for _, de := range entries {
		if de.IsDir() && strings.HasSuffix(de.Name(), ".anchor") && !strings.HasPrefix(de.Name(), ".staging-") {
			count++
		}
	}
	return count, nil
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

func contextTelemetryStats(path string, extraEntries ...tape.TapeEntry) (int, int64, error) {
	lines, err := readJSONLLines(path)
	if err != nil {
		return 0, 0, err
	}
	var count int
	var bytesTotal int64
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		count++
		bytesTotal += int64(len(redactContextImagesForTelemetry(line)) + 1)
	}
	for _, entry := range extraEntries {
		line, err := json.Marshal(entry)
		if err != nil {
			return 0, 0, err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		count++
		bytesTotal += int64(len(redactContextImagesForTelemetry(line)) + 1)
	}
	return count, bytesTotal, nil
}

func redactContextImagesForTelemetry(line []byte) []byte {
	var entry tape.TapeEntry
	if err := json.Unmarshal(line, &entry); err != nil {
		return line
	}

	switch entry.Type {
	case "message":
		var msg tape.Message
		if err := json.Unmarshal(entry.Data, &msg); err != nil || msg.Image == nil {
			return line
		}
		redactImagePartForTelemetry(msg.Image)
		return marshalTelemetryEntry(tape.MessageEntry(msg), line)
	case "tool_result":
		var result tape.ToolResult
		if err := json.Unmarshal(entry.Data, &result); err != nil || result.Image == nil {
			return line
		}
		redactImagePartForTelemetry(result.Image)
		return marshalTelemetryEntry(tape.ToolResultEntry(result), line)
	default:
		return line
	}
}

func redactImagePartForTelemetry(image *tape.ImagePart) {
	if image == nil || image.Data == "" {
		return
	}
	image.Data = fmt.Sprintf("[redacted %d base64 chars]", len(image.Data))
}

func marshalTelemetryEntry(entry tape.TapeEntry, fallback []byte) []byte {
	data, err := json.Marshal(entry)
	if err != nil {
		return fallback
	}
	return data
}

func estimateTokens(rawBytes int64) int64 {
	if rawBytes <= 0 {
		return 0
	}
	// Use a rough bytes->tokens conversion so large single-entry tool outputs
	// still surface as meaningful memory mass in telemetry.
	return (rawBytes + 3) / 4
}

func writeJSON(dir, name string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, name), data, 0o644)
}
