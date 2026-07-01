package runtime

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

func TestRunSuppressesInitialBegin(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
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
	cfg.SuppressInitialBegin = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	if exitCode := rt.Run("", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if len(mock.calls) != 1 {
		t.Fatalf("expected 1 LLM call, got %d", len(mock.calls))
	}
	for _, msg := range mock.calls[0] {
		if msg.Role == tape.RoleUser || msg.Content == "Begin." {
			t.Fatalf("provider input should not include synthetic Begin user message: %#v", mock.calls[0])
		}
	}
}

func TestShThenExit(t *testing.T) {
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
						ID:   "call_2",
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
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("run echo hi then exit", "Begin.")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls, got %d", mock.callCount)
	}

	// Verify sh tool result was appended to the tape
	msgs := rt.tape.Messages()
	foundToolResult := false
	for _, m := range msgs {
		if m.Role == tape.RoleToolResult {
			foundToolResult = true
			break
		}
	}
	if !foundToolResult {
		t.Error("expected tool_result message in tape after sh execution")
	}
}

func TestStructuredToolResultPersistsAfterRewrite(t *testing.T) {
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
						ID:   "call_2",
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

	exitCode := rt.Run("run echo hi then exit", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var inMemory *tape.Message
	for _, m := range rt.tape.Messages() {
		if m.Role == tape.RoleToolResult && m.ToolID == "call_1" {
			msg := m
			inMemory = &msg
			break
		}
	}
	if inMemory == nil {
		t.Fatal("expected sh tool_result on tape")
	}
	if len(inMemory.StructuredContent) == 0 {
		t.Fatal("expected in-memory tool_result structured content to be populated")
	}
	payload := decodeToolContent(t, inMemory.StructuredContent)
	runtimePayload := toolMap(t, payload, "runtime")
	if _, ok := runtimePayload["memory_status"]; !ok {
		t.Fatalf("expected rewritten tool_result content to include runtime.memory_status, got %s", string(inMemory.StructuredContent))
	}

	summary, err := tape.ReadTapeFile(cfg.TapePath(""))
	if err != nil {
		t.Fatalf("read tape file: %v", err)
	}
	foundOnDisk := false
	for _, entry := range summary.Entries {
		if entry.Type != "tool_result" {
			continue
		}
		var tr tape.ToolResult
		if err := json.Unmarshal(entry.Data, &tr); err != nil {
			t.Fatalf("unmarshal tool_result entry: %v", err)
		}
		if tr.ToolID != "call_1" {
			continue
		}
		foundOnDisk = true
		if len(tr.StructuredContent) != 0 {
			t.Fatal("expected on-disk tool_result structured content to remain empty")
		}
		payload := decodeToolContent(t, tr.Content)
		runtimePayload := toolMap(t, payload, "runtime")
		if _, ok := runtimePayload["memory_status"]; !ok {
			t.Fatalf("expected on-disk tool_result content to include runtime.memory_status, got %q", tr.Content)
		}
		break
	}
	if !foundOnDisk {
		t.Fatal("expected sh tool_result entry on disk")
	}
}

func TestToolResultIncludesMemoryActionAndWarningWhenFrontierWarns(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": `python3 - <<'PY'
print("x" * 40000)
PY`,
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
						Name: "sh",
						Arguments: map[string]any{
							"command": `python3 - <<'PY'
print("y" * 40000)
PY`,
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_3",
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

	exitCode := rt.Run("emit large output", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit 0, got %d", exitCode)
	}

	found := false
	var seen []string
	for _, m := range rt.tape.Messages() {
		if m.Role != tape.RoleToolResult {
			continue
		}
		seen = append(seen, m.Content)
		payload := decodeToolContent(t, m.StructuredContent)
		runtimePayload, _ := payload["runtime"].(map[string]any)
		if runtimePayload == nil {
			continue
		}
		_, hasStatus := runtimePayload["memory_status"].(map[string]any)
		action, _ := runtimePayload["memory_action"].(string)
		warning, _ := runtimePayload["memory_warning"].(string)
		if hasStatus &&
			strings.Contains(action, "next_action=mark") &&
			strings.Contains(warning, "before the next `sh`") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected tool_result to surface memory status, action, and warning blocks together; saw %q", seen)
	}
}

func TestMemoryDeathCutoffTerminatesAfterToolResult(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": `i=0; while [ "$i" -lt 2000 ]; do printf x; i=$((i + 1)); done; printf '\n'`,
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
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
	cfg.MemoryWarnTokens = 100
	cfg.MemoryDangerTokens = 200
	cfg.MemoryDeathTokens = 300
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("cross memory death cutoff", "Begin.")
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if mock.callCount != 1 {
		t.Fatalf("provider calls = %d, want 1", mock.callCount)
	}
	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome")
	}
	if rt.tape.Outcome.TerminationMode != tape.TermContextDeath {
		t.Fatalf("termination mode = %q, want %q", rt.tape.Outcome.TerminationMode, tape.TermContextDeath)
	}
	if !strings.Contains(rt.tape.Outcome.Stderr, "context death: frontier_estimated_tokens=") ||
		!strings.Contains(rt.tape.Outcome.Stderr, "death_tokens=300") {
		t.Fatalf("outcome stderr did not record death cause: %q", rt.tape.Outcome.Stderr)
	}

	var toolMsg *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult && msg.ToolID == "call_1" {
			copy := msg
			toolMsg = &copy
			break
		}
	}
	if toolMsg == nil {
		t.Fatal("expected sh tool result")
	}
	payload := decodeToolContent(t, toolMsg.StructuredContent)
	runtimePayload := toolMap(t, payload, "runtime")
	status := toolMap(t, runtimePayload, "memory_status")
	if got := toolString(t, status, "level"); got != "death" {
		t.Fatalf("memory level = %q, want death", got)
	}
	if got := toolInt(t, status, "death_tokens"); got != 300 {
		t.Fatalf("death_tokens = %d, want 300", got)
	}
}

func TestTextOnlyResponseContinuesLoop(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role:    tape.RoleAssistant,
				Content: "Let me think about this...",
				// No tool calls — text only
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
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
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls (text-only + exit), got %d", mock.callCount)
	}
}

func TestChildTextOnlyResponseTerminatesWithStdout(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role:    tape.RoleAssistant,
				Content: "child discriminator report",
			},
		},
	}

	cfg := testCfg(t)
	cfg.Depth = 1
	cfg.ParentSession = "parent-session"
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	stdout, err := os.CreateTemp(t.TempDir(), "stdout-*")
	if err != nil {
		t.Fatalf("create stdout capture: %v", err)
	}
	defer stdout.Close()
	rt.SetStdout(stdout)

	exitCode := rt.Run("child task", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if mock.callCount != 1 {
		t.Fatalf("expected one LLM call, got %d", mock.callCount)
	}
	if _, err := stdout.Seek(0, 0); err != nil {
		t.Fatalf("seek stdout capture: %v", err)
	}
	data, err := os.ReadFile(stdout.Name())
	if err != nil {
		t.Fatalf("read stdout capture: %v", err)
	}
	if string(data) != "child discriminator report" {
		t.Fatalf("stdout = %q", string(data))
	}
	if rt.tape.Outcome == nil || rt.tape.Outcome.ExitCode != 0 || rt.tape.Outcome.TerminationMode != tape.TermExit {
		t.Fatalf("outcome = %#v, want successful exit", rt.tape.Outcome)
	}
}

func TestEmptyAssistantResponseWithoutToolCallsFailsFast(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if mock.callCount != 1 {
		t.Fatalf("expected 1 LLM call, got %d", mock.callCount)
	}
	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be recorded")
	}
	if !strings.Contains(rt.tape.Outcome.Stderr, "empty assistant response") {
		t.Fatalf("expected stderr to mention empty assistant response, got %q", rt.tape.Outcome.Stderr)
	}
}

func TestReadyTextAutoIdleResumesWithInjectedPayload(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role:    tape.RoleAssistant,
				Content: "Ready. Awaiting your objective.",
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
	cfg.IdleEnabled = true
	cfg.ReadyTextAutoIdle = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("auto idle ready text", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool {
		return rt.control != nil && rt.control.started && mock.callCount == 1
	}, "ready-like response to enter auto-idle")
	select {
	case exitCode := <-done:
		t.Fatalf("runtime exited before inject, exit=%d", exitCode)
	case <-time.After(100 * time.Millisecond):
	}
	if mock.callCount != 1 {
		t.Fatalf("provider call count before inject = %d, want 1", mock.callCount)
	}

	if _, err := rt.enqueueControlPayload(controlActionInject, "wake payload"); err != nil {
		t.Fatalf("enqueue control payload: %v", err)
	}
	rt.requestControlInject()

	select {
	case exitCode := <-done:
		if exitCode != 0 {
			t.Fatalf("exit code = %d, want 0", exitCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runtime did not resume after inject")
	}
	if mock.callCount != 2 {
		t.Fatalf("provider call count after inject = %d, want 2", mock.callCount)
	}

	var sawPayload bool
	for _, msg := range mock.calls[1] {
		if msg.Role == tape.RoleAssistant && strings.Contains(msg.Content, "Ready. Awaiting") {
			t.Fatalf("ready-like assistant text leaked into provider context: %#v", msg)
		}
		if msg.Role == tape.RoleUser && strings.Contains(msg.Content, "wake payload") {
			sawPayload = true
		}
	}
	if !sawPayload {
		t.Fatalf("second provider call did not include injected payload: %#v", mock.calls[1])
	}
}

func TestEmptyAssistantResponseCanSucceedAfterToolProgress(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_sh",
						Name: "sh",
						Arguments: map[string]any{
							"command": "true",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
			},
		},
	}

	cfg := testCfg(t)
	cfg.EmptyAssistantSuccess = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if mock.callCount != 2 {
		t.Fatalf("expected 2 LLM calls, got %d", mock.callCount)
	}
	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be recorded")
	}
	if rt.tape.Outcome.ExitCode != 0 || rt.tape.Outcome.TerminationMode != tape.TermExit {
		t.Fatalf("outcome = %#v, want successful exit", rt.tape.Outcome)
	}
}
