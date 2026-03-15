package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func writeTestTape(t *testing.T, tapePath string, entries []tape.TapeEntry) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(tapePath), 0o755); err != nil {
		t.Fatalf("mkdir tape dir: %v", err)
	}
	f, err := os.OpenFile(tapePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		t.Fatalf("open tape file: %v", err)
	}
	defer f.Close()
	for _, e := range entries {
		line, _ := json.Marshal(e)
		if _, err := f.Write(line); err != nil {
			t.Fatalf("write tape entry: %v", err)
		}
		if _, err := f.Write([]byte("\n")); err != nil {
			t.Fatalf("write newline: %v", err)
		}
	}
}

func minimalTapeEntries() []tape.TapeEntry {
	tp := tape.NewTape("sess-1", "", 0, "test-model")
	return []tape.TapeEntry{
		tp.MetaEntry(),
		tape.MessageEntry(tape.Message{Role: tape.RoleSystem, Content: "sys"}),
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "u1"}),
	}
}

func TestAnchorMemory_MarkAndUnfold(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	tapePath := filepath.Join(tmp, "tapes", "sess-1", "0000.jsonl")
	writeTestTape(t, tapePath, minimalTapeEntries())

	exec := NewAnchorMemoryExecutor(agentRoot, tapePath)
	mark := exec.Mark("tool-1", MarkRequest{Summary: "checkpoint-1"})
	if mark.IsError {
		t.Fatalf("mark failed: %s", mark.Content)
	}
	if !strings.Contains(mark.Content, "anchor=0") {
		t.Fatalf("unexpected mark content: %s", mark.Content)
	}

	unfold := exec.Unfold("tool-2", UnfoldRequest{AnchorID: 0})
	if unfold.IsError {
		t.Fatalf("unfold failed: %s", unfold.Content)
	}
	if !strings.Contains(unfold.Content, "\"id\": 0") {
		t.Fatalf("unfold content missing anchor id: %s", unfold.Content)
	}
	if !strings.Contains(unfold.Content, "checkpoint-1") {
		t.Fatalf("unfold content missing summary: %s", unfold.Content)
	}
	if strings.Contains(unfold.Content, "\"session_id\"") {
		t.Fatalf("unfold content should not include tape meta entries: %s", unfold.Content)
	}
	if strings.Contains(unfold.Content, "\"sys\"") {
		t.Fatalf("unfold content should not include system prompt entries: %s", unfold.Content)
	}
	if !strings.Contains(unfold.Content, "\"u1\"") {
		t.Fatalf("unfold content missing user context entry: %s", unfold.Content)
	}

	feedback, err := exec.FeedbackBlock()
	if err != nil {
		t.Fatalf("feedback block: %v", err)
	}
	if !strings.Contains(feedback, "\"anchors\":[0]") {
		t.Fatalf("feedback missing frontier anchor: %s", feedback)
	}
}

func TestAnchorMemory_MarkFoldAbsorbsFrontier(t *testing.T) {
	tmp := t.TempDir()
	agentRoot := filepath.Join(tmp, "agent", "sess-1")
	tapePath := filepath.Join(tmp, "tapes", "sess-1", "0000.jsonl")
	writeTestTape(t, tapePath, minimalTapeEntries())
	exec := NewAnchorMemoryExecutor(agentRoot, tapePath)

	first := exec.Mark("tool-1", MarkRequest{Summary: "first"})
	if first.IsError {
		t.Fatalf("first mark failed: %s", first.Content)
	}

	// Append one more tape entry so there is new raw context before fold.
	entries := minimalTapeEntries()
	entries = append(entries, tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "a1"}))
	writeTestTape(t, tapePath, entries)

	second := exec.Mark("tool-2", MarkRequest{Summary: "folded", Fold: true})
	if second.IsError {
		t.Fatalf("second mark failed: %s", second.Content)
	}
	if !strings.Contains(second.Content, "absorbed=1") {
		t.Fatalf("second mark should absorb prior frontier: %s", second.Content)
	}

	ids, err := exec.frontierAnchorIDs()
	if err != nil {
		t.Fatalf("frontier ids: %v", err)
	}
	if len(ids) != 1 || ids[0] != 1 {
		t.Fatalf("frontier should contain only folded anchor 1, got %+v", ids)
	}
}

func TestParseMarkArgs_Validation(t *testing.T) {
	if _, err := ParseMarkArgs(map[string]any{}); err == nil {
		t.Fatal("expected missing summary error")
	}
	req, err := ParseMarkArgs(map[string]any{"summary": "  hello  ", "fold": true})
	if err != nil {
		t.Fatalf("parse mark args: %v", err)
	}
	if req.Summary != "hello" || !req.Fold {
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
