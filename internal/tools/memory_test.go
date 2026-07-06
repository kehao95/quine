package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func writeTestContextEntries(t *testing.T, exec *AnchorMemoryExecutor, entries []tape.TapeEntry) {
	t.Helper()
	for _, e := range entries {
		if tape.PartitionOf(e) != tape.PartitionContext {
			continue
		}
		if err := exec.AppendEntry(e); err != nil {
			t.Fatalf("append context entry: %v", err)
		}
	}
}

func minimalContextEntries() []tape.TapeEntry {
	return []tape.TapeEntry{
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "u1"}),
	}
}

func TestAnchorMemory_MarkAndUnfold(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())
	initialStatus, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("initial status block: %v", err)
	}
	initialData := decodeBlockJSON(t, initialStatus)
	if initialData["frontier_estimated_tokens"].(float64) <= 0 {
		t.Fatalf("expected frontier_estimated_tokens > 0 before mark, got %v", initialData["frontier_estimated_tokens"])
	}
	if initialData["frontier_level"].(string) != initialData["level"].(string) {
		t.Fatalf("expected frontier_level to mirror level, got frontier_level=%v level=%v", initialData["frontier_level"], initialData["level"])
	}
	if initialData["anchor_level"].(string) != "ok" {
		t.Fatalf("expected initial anchor_level=ok, got %v", initialData["anchor_level"])
	}
	if initialData["warn_tokens"].(float64) != 32 {
		t.Fatalf("expected warn_tokens=32, got %v", initialData["warn_tokens"])
	}
	if initialData["danger_tokens"].(float64) != 64 {
		t.Fatalf("expected danger_tokens=64, got %v", initialData["danger_tokens"])
	}

	topology, err := exec.TopologyBlock()
	if err != nil {
		t.Fatalf("initial topology block: %v", err)
	}
	topologyData := decodeBlockJSON(t, topology)
	if topologyData["frontier_bytes"].(float64) <= 0 {
		t.Fatalf("expected frontier_bytes > 0 before mark, got %v", topologyData["frontier_bytes"])
	}

	mark := exec.Mark("tool-1", MarkRequest{Resolution: "checkpoint-1"})
	if mark.IsError {
		t.Fatalf("mark failed: %s", mark.Content)
	}
	markPayload := decodeResultContent(t, mark.Content)
	if markPayload["tool"] != "mark" {
		t.Fatalf("mark tool = %#v, want mark", markPayload["tool"])
	}
	if markPayload["status"] != "completed" {
		t.Fatalf("mark status = %#v, want completed", markPayload["status"])
	}
	if resultInt(t, markPayload, "anchor_id") != 0 {
		t.Fatalf("mark anchor_id = %#v, want 0", markPayload["anchor_id"])
	}
	if markPayload["resolution"] != "checkpoint-1" {
		t.Fatalf("mark resolution = %#v, want checkpoint-1", markPayload["resolution"])
	}
	if resultInt(t, markPayload, "absorbed") != 0 {
		t.Fatalf("mark absorbed = %#v, want 0", markPayload["absorbed"])
	}
	if len(mark.StructuredContent) != 0 {
		t.Fatal("mark should leave structured content empty")
	}

	unfold := exec.Unfold("tool-2", UnfoldRequest{AnchorID: 0})
	if unfold.IsError {
		t.Fatalf("unfold failed: %s", unfold.Content)
	}
	unfoldPayload := decodeResultContent(t, unfold.Content)
	if unfoldPayload["resolution"] != "checkpoint-1" {
		t.Fatalf("unfold resolution = %#v, want checkpoint-1", unfoldPayload["resolution"])
	}
	turns, ok := unfoldPayload["turns"].([]any)
	if !ok || len(turns) != 1 {
		t.Fatalf("unfold turns = %#v, want single user entry", unfoldPayload["turns"])
	}
	turn, ok := turns[0].(map[string]any)
	if !ok {
		t.Fatalf("unfold turn item = %#v, want object", turns[0])
	}
	if turn["type"] != "message" {
		t.Fatalf("unfold turn type = %#v, want message", turn["type"])
	}
	data, ok := turn["data"].(map[string]any)
	if !ok {
		t.Fatalf("unfold turn data = %#v, want object", turn["data"])
	}
	if data["role"] != string(tape.RoleUser) {
		t.Fatalf("unfold turn role = %#v, want user", data["role"])
	}
	if data["content"] != "u1" {
		t.Fatalf("unfold turn content = %#v, want u1", data["content"])
	}
	if len(unfold.StructuredContent) != 0 {
		t.Fatal("unfold should leave structured content empty")
	}

	feedback, err := exec.FeedbackBlock()
	if err != nil {
		t.Fatalf("feedback block: %v", err)
	}
	if !strings.Contains(feedback, "\"anchors\":[0]") {
		t.Fatalf("feedback missing frontier anchor: %s", feedback)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["frontier_estimated_tokens"].(float64) != 0 {
		t.Fatalf("expected frontier_estimated_tokens to reset after mark, got %v", statusData["frontier_estimated_tokens"])
	}
	if statusData["level"].(string) != "ok" {
		t.Fatalf("expected level=ok after mark, got %v", statusData["level"])
	}

	topology, err = exec.TopologyBlock()
	if err != nil {
		t.Fatalf("topology block after mark: %v", err)
	}
	topologyData = decodeBlockJSON(t, topology)
	if topologyData["frontier_bytes"].(float64) != 0 {
		t.Fatalf("expected frontier_bytes to reset after mark, got %v", topologyData["frontier_bytes"])
	}
}

func TestAnchorMemory_MarkFoldAbsorbsFrontier(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())

	first := exec.Mark("tool-1", MarkRequest{Resolution: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	// Append one more context entry so there is new raw context before fold.
	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a1"})); err != nil {
		t.Fatalf("append assistant entry: %v", err)
	}

	second := exec.Mark("tool-2", MarkRequest{Resolution: "folded", Fold: true})
	if second.IsError {
		t.Fatalf("second mark failed: %s", second.Content)
	}
	secondPayload := decodeResultContent(t, second.Content)
	if secondPayload["fold"] != true {
		t.Fatalf("folded mark should report fold=true, got %#v", secondPayload["fold"])
	}
	if resultInt(t, secondPayload, "absorbed") != 1 {
		t.Fatalf("folded mark absorbed = %#v, want 1", secondPayload["absorbed"])
	}
	absorbedIDs, ok := secondPayload["absorbed_anchor_ids"].([]any)
	if !ok || len(absorbedIDs) != 1 {
		t.Fatalf("folded mark absorbed_anchor_ids = %#v, want [0]", secondPayload["absorbed_anchor_ids"])
	}
	switch value := absorbedIDs[0].(type) {
	case json.Number:
		if value.String() != "0" {
			t.Fatalf("folded mark absorbed_anchor_ids[0] = %#v, want 0", absorbedIDs[0])
		}
	case float64:
		if value != 0 {
			t.Fatalf("folded mark absorbed_anchor_ids[0] = %#v, want 0", absorbedIDs[0])
		}
	default:
		t.Fatalf("folded mark absorbed_anchor_ids[0] has unexpected type %#v", absorbedIDs[0])
	}
	if len(second.StructuredContent) != 0 {
		t.Fatal("folded mark should leave structured content empty")
	}

	ids, err := exec.frontierAnchorIDs()
	if err != nil {
		t.Fatalf("frontier ids: %v", err)
	}
	if len(ids) != 1 || ids[0] != 1 {
		t.Fatalf("frontier should contain only folded anchor 1, got %+v", ids)
	}
}

func TestAnchorMemory_FoldWithoutPriorAnchorsReportsNoConsolidation(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())

	result := exec.Mark("tool-1", MarkRequest{Resolution: "first", Fold: true})
	if result.IsError {
		t.Fatalf("fold-first mark failed: %s", result.Content)
	}
	payload := decodeResultContent(t, result.Content)
	if resultInt(t, payload, "absorbed") != 0 {
		t.Fatalf("expected absorbed=0 for fold without prior anchors, got %s", result.Content)
	}
	if payload["note"] == "" {
		t.Fatalf("expected structured note for ineffective fold, got %#v", payload["note"])
	}
	if !strings.Contains(payload["note"].(string), "no consolidation occurred") {
		t.Fatalf("expected no-consolidation note, got %#v", payload["note"])
	}
	if len(result.StructuredContent) != 0 {
		t.Fatal("fold-without-prior-anchors result should leave structured content empty")
	}
}

func TestAnchorMemory_StatusAndTopologyExposeGoverningAnchor(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())

	first := exec.Mark("tool-1", MarkRequest{Resolution: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a1"})); err != nil {
		t.Fatalf("append assistant entry: %v", err)
	}

	second := exec.Mark("tool-2", MarkRequest{Resolution: "folded", Fold: true})
	if second.IsError {
		t.Fatalf("second mark failed: %s", second.Content)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["governing_anchor"].(float64) != 1 {
		t.Fatalf("status governing_anchor = %#v, want 1", statusData["governing_anchor"])
	}

	topology, err := exec.TopologyBlock()
	if err != nil {
		t.Fatalf("topology block: %v", err)
	}
	topologyData := decodeBlockJSON(t, topology)
	if topologyData["governing_anchor"].(float64) != 1 {
		t.Fatalf("topology governing_anchor = %#v, want 1", topologyData["governing_anchor"])
	}
}

func TestAnchorMemory_StatusWarnsOnLargeFrontier(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, append(minimalContextEntries(), tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: strings.Repeat("x", 400)})))
	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["level"].(string) != "danger" {
		t.Fatalf("expected level=danger, got %v", statusData["level"])
	}
	if _, ok := statusData["warning"]; !ok {
		t.Fatalf("expected warning text in danger status, got %v", statusData)
	}
	if _, ok := statusData["frontier_warning"]; !ok {
		t.Fatalf("expected frontier_warning text in danger status, got %v", statusData)
	}
	if statusData["anchor_level"].(string) != "ok" {
		t.Fatalf("expected anchor_level=ok with no anchors, got %v", statusData["anchor_level"])
	}
}

func TestAnchorMemory_StatusRedactsImagePayloadsForTelemetry(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 8000, 16000)
	imageData := strings.Repeat("a", 40000)
	if err := exec.AppendEntry(tape.ToolResultEntry(tape.ToolResult{
		ToolID:  "vision-1",
		Content: tape.MarshalToolResultContent(map[string]any{"tool": "vision", "status": "completed"}),
		Image: &tape.ImagePart{
			MIMEType: "image/png",
			Data:     imageData,
		},
	})); err != nil {
		t.Fatalf("append image tool result: %v", err)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["level"].(string) != "ok" {
		t.Fatalf("expected level=ok after image telemetry redaction, got %v", statusData["level"])
	}
	if tokens := statusData["frontier_estimated_tokens"].(float64); tokens >= 8000 {
		t.Fatalf("expected redacted image telemetry below warn threshold, got %v", tokens)
	}
	if _, ok := statusData["next_action"]; ok {
		t.Fatalf("expected image payload alone not to trigger next_action, got %v", statusData["next_action"])
	}

	stored, err := os.ReadFile(exec.currentPath())
	if err != nil {
		t.Fatalf("read stored current context: %v", err)
	}
	if !strings.Contains(string(stored), imageData) {
		t.Fatalf("expected stored context to preserve image data")
	}
}

func TestAnchorMemory_StatusWarnsOnParallelAnchors(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())
	first := exec.Mark("tool-1", MarkRequest{Resolution: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a1"})); err != nil {
		t.Fatalf("append assistant entry: %v", err)
	}
	second := exec.Mark("tool-2", MarkRequest{Resolution: "second"})
	if second.IsError {
		t.Fatalf("second mark failed: %s", second.Content)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["anchor_level"].(string) != "warn" {
		t.Fatalf("expected anchor_level=warn, got %v", statusData["anchor_level"])
	}
	if statusData["topology"].(string) != "flat" {
		t.Fatalf("expected topology=flat, got %v", statusData["topology"])
	}
	if _, ok := statusData["anchor_warning"]; !ok {
		t.Fatalf("expected anchor_warning in status, got %v", statusData)
	}
	if !strings.Contains(statusData["anchor_warning"].(string), "do not force a fold") {
		t.Fatalf("expected anchor_warning to clarify that two anchors are not enough, got %v", statusData["anchor_warning"])
	}
	if !strings.Contains(statusData["anchor_warning"].(string), "newer higher-order resolution") {
		t.Fatalf("expected anchor_warning to require a newer higher-order resolution, got %v", statusData["anchor_warning"])
	}
	if _, ok := statusData["next_action"]; ok {
		t.Fatalf("expected no next_action at warn-level anchor pressure, got %v", statusData["next_action"])
	}
	if _, ok := statusData["next_action_timing"]; ok {
		t.Fatalf("expected no next_action_timing at warn-level anchor pressure, got %v", statusData["next_action_timing"])
	}
	if !strings.Contains(statusData["anchor_warning"].(string), "fold=true") {
		t.Fatalf("expected anchor_warning to mention fold=true, got %v", statusData["anchor_warning"])
	}
	if !strings.Contains(statusData["anchor_warning"].(string), "child findings") {
		t.Fatalf("expected anchor_warning to mention child findings as a valid fold case, got %v", statusData["anchor_warning"])
	}
}

func TestAnchorMemory_StatusRaisesFoldActionOnlyAtDangerParallelAnchors(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())
	if result := exec.Mark("tool-1", MarkRequest{Resolution: "first"}); result.IsError {
		t.Fatalf("first mark failed: %s", result.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a1"})); err != nil {
		t.Fatalf("append first assistant entry: %v", err)
	}
	if result := exec.Mark("tool-2", MarkRequest{Resolution: "second"}); result.IsError {
		t.Fatalf("second mark failed: %s", result.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a2"})); err != nil {
		t.Fatalf("append second assistant entry: %v", err)
	}
	if result := exec.Mark("tool-3", MarkRequest{Resolution: "third"}); result.IsError {
		t.Fatalf("third mark failed: %s", result.Content)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["anchor_level"].(string) != "danger" {
		t.Fatalf("expected anchor_level=danger, got %v", statusData["anchor_level"])
	}
	if statusData["topology"].(string) != "flat" {
		t.Fatalf("expected topology=flat, got %v", statusData["topology"])
	}
	if statusData["next_action"].(string) != "fold" {
		t.Fatalf("expected next_action=fold at danger-level anchor pressure, got %v", statusData["next_action"])
	}
	if statusData["next_action_timing"].(string) != "before_next_sh" {
		t.Fatalf("expected next_action_timing=before_next_sh, got %v", statusData["next_action_timing"])
	}
	if !strings.Contains(statusData["anchor_warning"].(string), "parent-level conclusion") {
		t.Fatalf("expected danger anchor_warning to mention parent-level conclusion folding, got %v", statusData["anchor_warning"])
	}
}

func TestAnchorMemory_FrontierWarningStillPointsToPlainMarkWithExistingAnchor(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())
	first := exec.Mark("tool-1", MarkRequest{Resolution: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: strings.Repeat("x", 400)})); err != nil {
		t.Fatalf("append assistant entry: %v", err)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["topology"].(string) != "forming" {
		t.Fatalf("expected topology=forming, got %v", statusData["topology"])
	}
	if statusData["frontier_level"].(string) == "ok" {
		t.Fatalf("expected non-ok frontier level, got %v", statusData["frontier_level"])
	}
	frontierWarning, ok := statusData["frontier_warning"].(string)
	if !ok {
		t.Fatalf("expected frontier_warning in status, got %v", statusData)
	}
	if !strings.Contains(frontierWarning, "another plain mark") {
		t.Fatalf("expected frontier_warning to keep pointing to plain mark, got %q", frontierWarning)
	}
	if statusData["next_action"].(string) != "mark" {
		t.Fatalf("expected next_action=mark, got %v", statusData["next_action"])
	}
	if statusData["next_action_timing"].(string) != "before_next_sh" {
		t.Fatalf("expected next_action_timing=before_next_sh, got %v", statusData["next_action_timing"])
	}
	if !strings.Contains(frontierWarning, "before the next `sh`") {
		t.Fatalf("expected frontier_warning to specify timing, got %q", frontierWarning)
	}
}

func TestAnchorMemory_FrontierPressureOverridesAnchorFoldHint(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())
	first := exec.Mark("tool-1", MarkRequest{Resolution: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a1"})); err != nil {
		t.Fatalf("append first assistant entry: %v", err)
	}
	second := exec.Mark("tool-2", MarkRequest{Resolution: "second"})
	if second.IsError {
		t.Fatalf("second mark failed: %s", second.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: strings.Repeat("x", 400)})); err != nil {
		t.Fatalf("append large assistant entry: %v", err)
	}

	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if statusData["anchor_level"].(string) != "warn" {
		t.Fatalf("expected anchor_level=warn, got %v", statusData["anchor_level"])
	}
	if statusData["frontier_level"].(string) == "ok" {
		t.Fatalf("expected non-ok frontier_level, got %v", statusData["frontier_level"])
	}
	if statusData["next_action"].(string) != "mark" {
		t.Fatalf("expected frontier pressure to override anchor fold hint with next_action=mark, got %v", statusData["next_action"])
	}
	if statusData["next_action_timing"].(string) != "before_next_sh" {
		t.Fatalf("expected next_action_timing=before_next_sh, got %v", statusData["next_action_timing"])
	}
}

func TestAnchorMemory_ActionAndWarningBlocksFollowFrontierPressure(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	writeTestContextEntries(t, exec, minimalContextEntries())
	first := exec.Mark("tool-1", MarkRequest{Resolution: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: strings.Repeat("x", 400)})); err != nil {
		t.Fatalf("append large assistant entry: %v", err)
	}

	action, err := exec.ActionBlock()
	if err != nil {
		t.Fatalf("action block: %v", err)
	}
	if !strings.Contains(action, "next_action=mark") {
		t.Fatalf("expected action block to surface mark interrupt, got %q", action)
	}
	if !strings.Contains(action, "next_action_timing=before_next_sh") {
		t.Fatalf("expected action timing in action block, got %q", action)
	}

	warning, err := exec.WarningBlock()
	if err != nil {
		t.Fatalf("warning block: %v", err)
	}
	if !strings.Contains(warning, "[MEMORY WARNING]") {
		t.Fatalf("expected warning block prefix, got %q", warning)
	}
	if !strings.Contains(warning, "call `mark` before the next `sh`") {
		t.Fatalf("expected warning block to preserve imperative text, got %q", warning)
	}
}

func TestAnchorMemory_MarkDisabledWarningCanOmitStrategyHints(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 32, 64)
	exec.MarkDisabled = true
	exec.MemoryStrategyHints = false
	writeTestContextEntries(t, exec, minimalContextEntries())
	if err := exec.AppendEntry(tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: strings.Repeat("x", 400)})); err != nil {
		t.Fatalf("append large assistant entry: %v", err)
	}

	warning, err := exec.WarningBlock()
	if err != nil {
		t.Fatalf("warning block: %v", err)
	}
	if !strings.Contains(warning, "[MEMORY WARNING]") {
		t.Fatalf("expected warning block prefix, got %q", warning)
	}
	if !strings.Contains(warning, "reduce live context") {
		t.Fatalf("expected neutral reduce-context goal, got %q", warning)
	}
	forbidden := []string{
		"peer process",
		"hand your `context/state`",
		"compact",
		"mark",
	}
	for _, text := range forbidden {
		if strings.Contains(warning, text) {
			t.Fatalf("neutral warning should omit %q, got %q", text, warning)
		}
	}
}

func TestAnchorMemory_StatusOmitsNextActionWhenFrontierOK(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	exec := NewAnchorMemoryExecutor(agentRoot, 8000, 16000)
	writeTestContextEntries(t, exec, minimalContextEntries())
	status, err := exec.StatusBlock()
	if err != nil {
		t.Fatalf("status block: %v", err)
	}
	statusData := decodeBlockJSON(t, status)
	if _, ok := statusData["next_action"]; ok {
		t.Fatalf("expected no next_action when frontier is ok, got %v", statusData["next_action"])
	}
	if _, ok := statusData["next_action_timing"]; ok {
		t.Fatalf("expected no next_action_timing when frontier is ok, got %v", statusData["next_action_timing"])
	}
}

func TestParseMarkArgs_Validation(t *testing.T) {
	if _, err := ParseMarkArgs(map[string]any{}); err == nil {
		t.Fatal("expected missing resolution error")
	}
	req, err := ParseMarkArgs(map[string]any{"resolution": "  hello  ", "fold": true})
	if err != nil {
		t.Fatalf("parse mark args: %v", err)
	}
	if req.Resolution != "hello" || !req.Fold {
		t.Fatalf("unexpected mark request: %+v", req)
	}
}

func TestParseUnfoldArgs_Validation(t *testing.T) {
	if _, err := ParseUnfoldArgs(map[string]any{}); err == nil {
		t.Fatal("expected missing anchor_id error")
	}
	req, err := ParseUnfoldArgs(map[string]any{"anchor_id": float64(3)})
	if err != nil {
		t.Fatalf("parse unfold args: %v", err)
	}
	if req.AnchorID != 3 {
		t.Fatalf("unexpected anchor id: %d", req.AnchorID)
	}
}

func decodeBlockJSON(t *testing.T, block string) map[string]any {
	t.Helper()
	idx := strings.Index(block, "{")
	if idx < 0 {
		t.Fatalf("block missing JSON payload: %q", block)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(block[idx:]), &out); err != nil {
		t.Fatalf("decode block json: %v", err)
	}
	return out
}
