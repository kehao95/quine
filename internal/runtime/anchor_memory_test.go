package runtime

import (
	"encoding/json"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestMarkTool_WithMemoryMetaFeedback(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_mark",
						Name: "mark",
						Arguments: map[string]any{
							"resolution": "checkpoint",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.AnchorMemoryEnabled = true
	cfg.MemoryWarnTokens = 1000
	cfg.MemoryDangerTokens = 2000
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("checkpoint", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d", exitCode)
	}

	foundMark := false
	foundMeta := false
	for _, m := range rt.tape.Messages() {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		if toolString(t, payload, "tool") == "mark" {
			foundMark = true
			if toolString(t, payload, "resolution") != "checkpoint" {
				t.Fatalf("mark resolution = %#v, want checkpoint", payload["resolution"])
			}
			runtimePayload := toolMap(t, payload, "runtime")
			if _, ok := runtimePayload["memory_feedback"]; ok {
				foundMeta = true
			}
		}
	}
	if !foundMark {
		t.Fatal("expected mark tool_result on tape")
	}
	if !foundMeta {
		t.Fatal("expected memory meta feedback in tool_result")
	}
}

func TestMarkTool_ReSeedsReasoningForAssistantToolBatch(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role:             tape.RoleAssistant,
				ReasoningContent: "checkpoint the frontier before continuing",
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_mark",
						Name: "mark",
						Arguments: map[string]any{
							"resolution": "checkpoint",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.AnchorMemoryEnabled = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("checkpoint", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit 0, got %d", exitCode)
	}
	if len(mock.calls) < 2 {
		t.Fatalf("provider call count = %d, want at least 2", len(mock.calls))
	}

	foundReseededMark := false
	for _, msg := range mock.calls[1] {
		if msg.Role != tape.RoleAssistant || len(msg.ToolCalls) != 1 || msg.ToolCalls[0].ID != "call_mark" {
			continue
		}
		foundReseededMark = true
		if msg.ReasoningContent != "checkpoint the frontier before continuing" {
			t.Fatalf("re-seeded mark reasoning = %q, want preserved reasoning", msg.ReasoningContent)
		}
	}
	if !foundReseededMark {
		t.Fatal("expected second provider call to include the re-seeded mark tool batch")
	}
}

func TestAnchorMemoryStatusVisibleOnOrdinaryToolResult(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "echo hi",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.AnchorMemoryEnabled = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("status", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d", exitCode)
	}

	foundStatus := false
	for _, m := range rt.tape.Messages() {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		runtimePayload, _ := payload["runtime"].(map[string]any)
		if runtimePayload == nil {
			continue
		}
		status, _ := runtimePayload["memory_status"].(map[string]any)
		if status != nil {
			if _, ok := status["frontier_estimated_tokens"]; ok {
				foundStatus = true
				break
			}
		}
	}
	if !foundStatus {
		t.Fatal("expected memory status feedback on ordinary tool_result")
	}

	summary, err := tape.ReadTapeFile(cfg.TapePath(""))
	if err != nil {
		t.Fatalf("read tape file: %v", err)
	}
	foundStatusOnDisk := false
	for _, entry := range summary.Entries {
		if entry.Type != "tool_result" {
			continue
		}
		var tr tape.ToolResult
		if err := json.Unmarshal(entry.Data, &tr); err != nil {
			t.Fatalf("unmarshal tool_result entry: %v", err)
		}
		payload := decodeToolContent(t, tr.Content)
		runtimePayload, _ := payload["runtime"].(map[string]any)
		status, _ := runtimePayload["memory_status"].(map[string]any)
		if status != nil {
			if _, ok := status["frontier_estimated_tokens"]; ok {
				foundStatusOnDisk = true
				break
			}
		}
	}
	if !foundStatusOnDisk {
		t.Fatal("expected memory status feedback to be persisted in tape file")
	}
}
