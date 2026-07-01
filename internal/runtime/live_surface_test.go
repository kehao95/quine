package runtime

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
)

// mockStreamingProvider is a test double implementing llm.StreamingProvider. Each
// GenerateStream call pops the next scripted event sequence and replays it on a
// buffered channel that is closed immediately, so consumers drain a complete,
// deterministic stream.
type mockStreamingProvider struct {
	scripts   [][]llm.StreamEvent
	callCount int
	setupErr  error
}

func (m *mockStreamingProvider) Generate(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	return tape.Message{}, llm.Usage{}, nil
}

func (m *mockStreamingProvider) ContextWindowSize() int { return 200000 }

func (m *mockStreamingProvider) GenerateStream(msgs []tape.Message, tools []llm.ToolSchema) (<-chan llm.StreamEvent, error) {
	if m.setupErr != nil {
		return nil, m.setupErr
	}
	var evs []llm.StreamEvent
	if m.callCount < len(m.scripts) {
		evs = m.scripts[m.callCount]
	}
	m.callCount++
	ch := make(chan llm.StreamEvent, len(evs)+1)
	for _, e := range evs {
		ch <- e
	}
	close(ch)
	return ch, nil
}

func readLiveEntries(t *testing.T, path string) []liveEntry {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read live.jsonl %s: %v", path, err)
	}
	var out []liveEntry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e liveEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("unmarshal live entry %q: %v", line, err)
		}
		out = append(out, e)
	}
	return out
}

func newStreamingRuntime(t *testing.T, sp llm.StreamingProvider) *Runtime {
	t.Helper()
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, sp)
	silenceRuntime(rt)
	rt.originalInput = "live surface test"
	return rt
}

// T-LIVE-WRITE: generateStreaming records the full delta sequence plus a
// terminal completed entry, and returns the authoritative final cell carried by
// the StreamDone event.
func TestGenerateStreaming_WritesDeltaSequence(t *testing.T) {
	final := tape.Message{
		Role:    tape.RoleAssistant,
		Content: "Hello world",
		ToolCalls: []tape.ToolCall{
			{ID: "call_1", Name: "bash", Arguments: map[string]any{"command": "echo hi"}},
		},
	}
	wantUsage := llm.Usage{InputTokens: 120, OutputTokens: 42}
	sp := &mockStreamingProvider{scripts: [][]llm.StreamEvent{{
		{Kind: llm.StreamText, Text: "Hello "},
		{Kind: llm.StreamText, Text: "world"},
		{Kind: llm.StreamReasoning, Text: "deciding to run echo"},
		{Kind: llm.StreamToolCall, Text: `{"command":"echo hi"}`, ToolID: "call_1", ToolName: "bash"},
		{Kind: llm.StreamDone, Message: final, Usage: wantUsage},
	}}}

	rt := newStreamingRuntime(t, sp)

	msg, usage, err := rt.generateStreaming(sp)
	if err != nil {
		t.Fatalf("generateStreaming: %v", err)
	}
	if !reflect.DeepEqual(msg, final) {
		t.Fatalf("final message = %#v, want %#v", msg, final)
	}
	if usage != wantUsage {
		t.Fatalf("usage = %#v, want %#v", usage, wantUsage)
	}

	entries := readLiveEntries(t, rt.contextLivePath())
	wantKinds := []string{"text_delta", "text_delta", "reasoning_delta", "tool_call_delta", "completed"}
	if len(entries) != len(wantKinds) {
		t.Fatalf("live entries = %d, want %d (%+v)", len(entries), len(wantKinds), entries)
	}
	for i, e := range entries {
		if e.Seq != i+1 {
			t.Fatalf("entry %d seq = %d, want %d", i, e.Seq, i+1)
		}
		if e.Kind != wantKinds[i] {
			t.Fatalf("entry %d kind = %q, want %q", i, e.Kind, wantKinds[i])
		}
		if e.TS == 0 {
			t.Fatalf("entry %d has zero timestamp", i)
		}
	}
	tc := entries[3]
	if tc.ToolID != "call_1" || tc.ToolName != "bash" || tc.Text != `{"command":"echo hi"}` {
		t.Fatalf("tool_call_delta entry = %#v, want call_1/bash with args", tc)
	}
	if entries[4].Text != "" {
		t.Fatalf("completed entry should carry no text, got %q", entries[4].Text)
	}
}

// T-LIVE-NOT-INPUT: provider-input assembly is identical whether or not
// live.jsonl exists — the transient surface never feeds the model.
func TestLiveSurfaceNeverProviderInput(t *testing.T) {
	sp := &mockStreamingProvider{}
	rt := newStreamingRuntime(t, sp)

	currentPath := rt.contextCurrentPath()
	writeJSONLEntries(t, currentPath,
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "from context user"}),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "from context assistant"}),
	)

	before, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages (no live): %v", err)
	}

	// Seed live.jsonl with deltas that would corrupt pairing if ever read.
	livePath := rt.contextLivePath()
	rt.appendLiveEntry(livePath, liveEntry{Seq: 1, Kind: "text_delta", Text: "partial cell that must be ignored"})
	rt.appendLiveEntry(livePath, liveEntry{Seq: 2, Kind: "tool_call_delta", Text: `{"command":`, ToolID: "ghost", ToolName: "bash"})
	if _, err := os.Stat(livePath); err != nil {
		t.Fatalf("live.jsonl should exist after seeding: %v", err)
	}

	after, err := rt.providerContextMessages()
	if err != nil {
		t.Fatalf("providerContextMessages (with live): %v", err)
	}

	if !reflect.DeepEqual(before, after) {
		t.Fatalf("provider input changed when live.jsonl present:\n before=%#v\n after =%#v", before, after)
	}
	for _, m := range after {
		if strings.Contains(m.Content, "must be ignored") || m.ToolID == "ghost" {
			t.Fatalf("provider input leaked a live delta: %#v", m)
		}
	}
}

// T-TRUNCATE: each generation truncates live.jsonl, so seq resets and no stale
// entries from the prior generation survive.
func TestGenerateStreaming_TruncatesPerGeneration(t *testing.T) {
	gen1 := tape.Message{Role: tape.RoleAssistant, Content: "first"}
	gen2 := tape.Message{Role: tape.RoleAssistant, Content: "second"}
	sp := &mockStreamingProvider{scripts: [][]llm.StreamEvent{
		{
			{Kind: llm.StreamText, Text: "AAAA"},
			{Kind: llm.StreamText, Text: "BBBB"},
			{Kind: llm.StreamDone, Message: gen1},
		},
		{
			{Kind: llm.StreamText, Text: "CCCC"},
			{Kind: llm.StreamDone, Message: gen2},
		},
	}}
	rt := newStreamingRuntime(t, sp)

	if _, _, err := rt.generateStreaming(sp); err != nil {
		t.Fatalf("generateStreaming gen1: %v", err)
	}
	if _, _, err := rt.generateStreaming(sp); err != nil {
		t.Fatalf("generateStreaming gen2: %v", err)
	}

	entries := readLiveEntries(t, rt.contextLivePath())
	// gen2 = one text_delta + one completed.
	if len(entries) != 2 {
		t.Fatalf("post-gen2 live entries = %d, want 2 (stale gen1 entries not truncated?): %+v", len(entries), entries)
	}
	if entries[0].Seq != 1 {
		t.Fatalf("first entry seq = %d, want 1 (seq did not reset)", entries[0].Seq)
	}
	for _, e := range entries {
		if e.Text == "AAAA" || e.Text == "BBBB" {
			t.Fatalf("stale gen1 delta survived truncation: %#v", e)
		}
	}
	if entries[0].Text != "CCCC" {
		t.Fatalf("entry[0].Text = %q, want CCCC", entries[0].Text)
	}
}

// T-RECOVERY: the recovery reader consumes only current.jsonl; live deltas
// sitting beside it are never recovered.
func TestLiveSurfaceNotRecovered(t *testing.T) {
	sp := &mockStreamingProvider{}
	rt := newStreamingRuntime(t, sp)

	currentPath := rt.contextCurrentPath()
	writeJSONLEntries(t, currentPath,
		tape.MessageEntry(tape.Message{Role: tape.RoleUser, Content: "recover me"}),
		tape.MessageEntry(tape.Message{Role: tape.RoleAssistant, Content: "and me"}),
	)

	livePath := rt.contextLivePath()
	rt.appendLiveEntry(livePath, liveEntry{Seq: 1, Kind: "text_delta", Text: "do not recover"})
	rt.appendLiveEntry(livePath, liveEntry{Seq: 2, Kind: "reasoning_delta", Text: "ephemeral"})
	rt.appendLiveEntry(livePath, liveEntry{Seq: 3, Kind: "completed"})

	recovered, err := readContextEntries(rt.contextCurrentPath())
	if err != nil {
		t.Fatalf("readContextEntries: %v", err)
	}
	if len(recovered) != 2 {
		t.Fatalf("recovered entries = %d, want 2 (live deltas must not be recovered): %+v", len(recovered), recovered)
	}
	for _, e := range recovered {
		if e.Type != "message" {
			t.Fatalf("recovered a non-message entry (live leak?): %#v", e)
		}
	}
}
