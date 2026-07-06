package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
)

func TestSimpleExit(t *testing.T) {
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
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("say hello", "Begin.")

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	if mock.callCount != 1 {
		t.Errorf("expected 1 LLM call, got %d", mock.callCount)
	}
}

func TestNonZeroExit(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "exit",
						Arguments: map[string]any{
							"status": "failure",
							"stderr": "something went wrong",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	// Capture stderr output from exit
	oldStderr := rt.stderr
	r, w, _ := os.Pipe()
	rt.stderr = w
	// Re-silence loggers but keep stderr pipe for exit tool output
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}

	exitCode := rt.Run("fail please", "Begin.")

	w.Close()
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()
	rt.stderr = oldStderr

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	stderr := string(buf[:n])
	if stderr != "something went wrong" {
		t.Errorf("expected stderr %q, got %q", "something went wrong", stderr)
	}
}

func TestAuthError(t *testing.T) {
	provider := &mockErrorProvider{err: llm.ErrAuth}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, provider)
	silenceRuntime(rt)

	exitCode := rt.Run("hello", "Begin.")

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for auth error, got %d", exitCode)
	}
}

func TestContextOverflowError(t *testing.T) {
	provider := &mockErrorProvider{err: llm.ErrContextOverflow}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, provider)
	silenceRuntime(rt)

	exitCode := rt.Run("hello", "Begin.")

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for context overflow, got %d", exitCode)
	}

	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be set")
	}
	if rt.tape.Outcome.TerminationMode != tape.TermContextExhaustion {
		t.Errorf("expected termination mode %q, got %q",
			tape.TermContextExhaustion, rt.tape.Outcome.TerminationMode)
	}
}

func TestRecoverableInferenceError(t *testing.T) {
	provider := &mockErrorProvider{
		err: fmt.Errorf("%w: reading response body: unexpected EOF", llm.ErrRecoverableInference),
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, provider)
	silenceRuntime(rt)

	exitCode := rt.Run("hello", "Begin.")

	if exitCode != RecoverableInferenceExitCode {
		t.Errorf("expected exit code %d for recoverable inference error, got %d", RecoverableInferenceExitCode, exitCode)
	}

	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be set")
	}
	if rt.tape.Outcome.TerminationMode != tape.TermRecoverableInference {
		t.Errorf("expected termination mode %q, got %q",
			tape.TermRecoverableInference, rt.tape.Outcome.TerminationMode)
	}
	if rt.tape.Outcome.ExitCode != RecoverableInferenceExitCode {
		t.Errorf("outcome exit code = %d, want %d", rt.tape.Outcome.ExitCode, RecoverableInferenceExitCode)
	}
	if len(rt.tape.Messages()) != 2 {
		t.Fatalf("recoverable inference failure should not append an assistant/error message, got %d messages", len(rt.tape.Messages()))
	}
}

func TestExecutionBudgetHardFailAllowsFinalExitResponse(t *testing.T) {
	// With hard_fail policy, reaching MaxTurns after a sh call should allow one
	// final exit-only response.
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "echo hello",
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
							"command": "echo world",
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
	cfg.MaxTurns = 2
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for final exit response, got %d", exitCode)
	}

	if mock.callCount != 3 {
		t.Errorf("expected exactly 3 LLM calls with hard_fail policy, got %d", mock.callCount)
	}

	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be set")
	}
	if rt.tape.Outcome.TerminationMode != tape.TermExit {
		t.Errorf("expected termination mode %q, got %q",
			tape.TermExit, rt.tape.Outcome.TerminationMode)
	}

	foundNoTurnsLeft := false
	for _, m := range rt.tape.Messages() {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		runtimePayload, _ := payload["runtime"].(map[string]any)
		if runtimePayload == nil {
			continue
		}
		if noTurnsLeft, _ := runtimePayload["no_turns_left"].(string); strings.Contains(noTurnsLeft, "Do not call `sh`. Call `exit` now.") {
			foundNoTurnsLeft = true
			break
		}
	}
	if !foundNoTurnsLeft {
		t.Error("expected final sh tool result to include [NO TURNS LEFT] guidance")
	}
}

func TestExecutionBudgetHardFailRejectsNonExitFinalResponse(t *testing.T) {
	// After MaxTurns is reached, one final response is allowed but only exit is
	// accepted.
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "echo hello",
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
							"command": "echo world",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_3",
						Name: "sh",
						Arguments: map[string]any{
							"command": "echo should_not_run",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.MaxTurns = 2
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for execution exhaustion, got %d", exitCode)
	}

	if mock.callCount != 3 {
		t.Errorf("expected exactly 3 LLM calls with hard_fail policy, got %d", mock.callCount)
	}

	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be set")
	}
	if rt.tape.Outcome.TerminationMode != tape.TermTurnExhaustion {
		t.Errorf("expected termination mode %q, got %q",
			tape.TermTurnExhaustion, rt.tape.Outcome.TerminationMode)
	}
	if !strings.Contains(rt.tape.Outcome.Stderr, "execution budget exhausted") {
		t.Errorf("expected stderr to contain 'execution budget exhausted', got %q", rt.tape.Outcome.Stderr)
	}

	foundRejection := false
	for _, m := range rt.tape.Messages() {
		if m.Role == tape.RoleToolResult && strings.Contains(string(m.StructuredContent), "Only exit is accepted after execution budget is exhausted.") {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Error("expected non-exit final response to be rejected")
	}
}

func TestExecutionBudgetFeedbackMessages(t *testing.T) {
	// With MaxTurns=3, after first sh call the tool result should contain
	// \"[EXECUTIONS LEFT] 2\".
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
	cfg.MaxTurns = 3
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Verify execution budget and context feedback were appended.
	msgs := rt.tape.Messages()
	foundBudget := false
	foundContext := false
	for _, m := range msgs {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		runtimePayload, _ := payload["runtime"].(map[string]any)
		if runtimePayload == nil {
			continue
		}
		if value, ok := runtimePayload["executions_left"]; ok {
			foundBudget = true
			if got := toolInt(t, map[string]any{"executions_left": value}, "executions_left"); got != 2 {
				t.Errorf("expected executions_left=2 in tool result, got %s", string(m.StructuredContent))
			}
		}
		if _, ok := runtimePayload["context_used_k"]; ok {
			foundContext = true
		}
	}
	if !foundBudget {
		t.Error("expected [EXECUTIONS LEFT] in a tool result message, found none")
	}
	if !foundContext {
		t.Error("expected [CONTEXT USED] in a tool result message, found none")
	}

	summary, err := tape.ReadTapeFile(cfg.TapePath(""))
	if err != nil {
		t.Fatalf("read tape file: %v", err)
	}
	foundPersistedBudget := false
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
		if runtimePayload == nil {
			continue
		}
		if value, ok := runtimePayload["executions_left"]; ok {
			if toolInt(t, map[string]any{"executions_left": value}, "executions_left") == 2 {
				foundPersistedBudget = true
				break
			}
		}
	}
	if !foundPersistedBudget {
		t.Error("expected persisted tape to contain runtime.executions_left=2, found none")
	}
}

func TestDisabledForkAndVisionAreRejected(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "fork",
						Arguments: map[string]any{
							"children": []map[string]any{{"intent": "probe"}},
						},
					},
					{
						ID:   "call_2",
						Name: "vision",
						Arguments: map[string]any{
							"image_path": "/tmp/demo.png",
							"question":   "what is this?",
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
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ForkEnabled = false
	cfg.VisionEnabled = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("stay minimal", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	msgs := rt.tape.Messages()
	var sawForkReject, sawVisionReject bool
	for _, m := range msgs {
		if m.Role != tape.RoleToolResult {
			continue
		}
		if strings.Contains(string(m.StructuredContent), "Rejected: fork is disabled in this runtime.") {
			sawForkReject = true
		}
		if strings.Contains(string(m.StructuredContent), "Rejected: vision is disabled in this runtime.") {
			sawVisionReject = true
		}
	}
	if !sawForkReject {
		t.Fatal("expected disabled fork rejection in tape")
	}
	if !sawVisionReject {
		t.Fatal("expected disabled vision rejection in tape")
	}
}

func TestExecutionBudgetZeroMeansDisabled(t *testing.T) {
	// MaxTurns=0 should not inject execution budget feedback into tool results.
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
	cfg.MaxTurns = 0
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Verify NO execution budget marker in tool results (but context should remain).
	msgs := rt.tape.Messages()
	foundContext := false
	for _, m := range msgs {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		runtimePayload, _ := payload["runtime"].(map[string]any)
		if runtimePayload == nil {
			continue
		}
		if _, ok := runtimePayload["executions_left"]; ok {
			t.Error("did not expect runtime.executions_left in tool results when MaxTurns=0")
		}
		if _, ok := runtimePayload["context_used_k"]; ok {
			foundContext = true
		}
	}
	if !foundContext {
		t.Error("expected runtime.context_used_k in tool result even when MaxTurns=0")
	}
}

func TestFailureWithoutReasonIsRejected(t *testing.T) {
	// Agent first tries to exit with status="failure" but no stderr.
	// Runtime rejects the exit and sends a tool result.
	// Agent retries with a reason in stderr.
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "exit",
						Arguments: map[string]any{
							"status": "failure",
							// no stderr — should be rejected
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
							"status": "failure",
							"stderr": "file not found",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	// Capture stderr
	oldStderr := rt.stderr
	r, w, _ := os.Pipe()
	rt.stderr = w
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}

	exitCode := rt.Run("do something that fails", "Begin.")

	w.Close()
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()
	rt.stderr = oldStderr

	// Should have made 2 LLM calls (rejected first, accepted second)
	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls (rejection + retry), got %d", mock.callCount)
	}

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	stderr := string(buf[:n])
	if stderr != "file not found" {
		t.Errorf("expected stderr %q, got %q", "file not found", stderr)
	}

	// Verify rejection tool result was added to the tape
	msgs := rt.tape.Messages()
	foundRejection := false
	for _, m := range msgs {
		if m.Role == tape.RoleToolResult && strings.Contains(string(m.StructuredContent), "Exit rejected") {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Error("expected a rejection tool result in the tape, found none")
	}
}

func TestFailureExitIsUnavailableWhenFailOnImpossibleDisabled(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "exit",
						Arguments: map[string]any{
							"status": "failure",
							"stderr": "blocked",
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
	cfg.FailOnImpossible = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("stay alive until success", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0 after retrying with success, got %d", exitCode)
	}
	if mock.callCount != 2 {
		t.Fatalf("expected 2 LLM calls, got %d", mock.callCount)
	}

	msgs := rt.tape.Messages()
	foundRejection := false
	for _, m := range msgs {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		if toolString(t, payload, "error") == "Exit rejected: status=\"failure\" is unavailable in this runtime; only status=\"success\" is allowed." {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Fatal("expected rejection tool result for failure exit in success-only mode")
	}
}

func TestInvalidExitArgsAreRejectedWhenFailOnImpossibleDisabled(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "exit",
						Arguments: map[string]any{
							"status": "progress",
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
	cfg.FailOnImpossible = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("stay alive until success", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0 after invalid exit rejection, got %d", exitCode)
	}
	if mock.callCount != 2 {
		t.Fatalf("expected 2 LLM calls, got %d", mock.callCount)
	}

	msgs := rt.tape.Messages()
	foundRejection := false
	for _, m := range msgs {
		if m.Role != tape.RoleToolResult {
			continue
		}
		payload := decodeToolContent(t, m.StructuredContent)
		if strings.Contains(toolString(t, payload, "error"), "invalid exit args") {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Fatal("expected rejection tool result for invalid exit args in success-only mode")
	}
}

func TestSuccessWithStderrIsRejected(t *testing.T) {
	// Agent tries to exit with status="success" but includes stderr.
	// Runtime rejects. Agent retries correctly with status="failure".
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
							"output": "",
							"stderr": "context window exceeded",
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
							"status": "failure",
							"stderr": "context window exceeded",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	oldStderr := rt.stderr
	r, w, _ := os.Pipe()
	rt.stderr = w
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}

	exitCode := rt.Run("task that incorrectly claims success", "Begin.")

	w.Close()
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	r.Close()
	rt.stderr = oldStderr

	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls (rejection + retry), got %d", mock.callCount)
	}

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}

	stderr := string(buf[:n])
	if stderr != "context window exceeded" {
		t.Errorf("expected stderr %q, got %q", "context window exceeded", stderr)
	}

	// Verify rejection tool result in tape
	msgs := rt.tape.Messages()
	foundRejection := false
	for _, m := range msgs {
		if m.Role == tape.RoleToolResult && strings.Contains(string(m.StructuredContent), "Exit rejected") {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Error("expected a rejection tool result in the tape, found none")
	}
}

func TestProgressStatusIsRejected(t *testing.T) {
	// "progress" status was removed. Agent tries it, gets parse error, exit code 1.
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "exit",
						Arguments: map[string]any{
							"status": "progress",
							"stderr": "context window at 90%, delegating remaining work",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	// Capture stderr
	oldStderr := rt.stderr
	rErr, wErr, _ := os.Pipe()
	rt.stderr = wErr
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}

	exitCode := rt.Run("find all needles", "Begin.")

	wErr.Close()
	errBuf := make([]byte, 4096)
	nErr, _ := rErr.Read(errBuf)
	rErr.Close()
	rt.stderr = oldStderr

	// "progress" is now an invalid status — parsed as failure (exit code 1)
	if exitCode != 1 {
		t.Errorf("expected exit code 1 (invalid status treated as failure), got %d", exitCode)
	}

	stderr := string(errBuf[:nErr])
	if !strings.Contains(stderr, "invalid exit args") {
		t.Errorf("expected stderr to contain 'invalid exit args', got %q", stderr)
	}
}

// TestProgressWithoutReasonIsRejected is removed — "progress" status no longer exists.
// The equivalent behavior is tested by TestProgressStatusIsRejected above.

// ---------------------------------------------------------------------------
// SIGALRM / Panic Mode tests (§2.2)
// ---------------------------------------------------------------------------

func TestPanicModeInjectsOverrideMessage(t *testing.T) {
	// When panicMode is set before the turn loop runs, the runtime should
	// inject a "System 1 Override" message into the tape.
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
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	// Simulate SIGALRM: set panic mode before Run starts the loop
	rt.panicMode.Store(true)

	exitCode := rt.Run("some task", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Verify the override message was injected into the tape
	msgs := rt.tape.Messages()
	foundOverride := false
	for _, m := range msgs {
		if m.Role == tape.RoleUser && strings.Contains(m.Content, "System interrupt") && strings.Contains(m.Content, "Time limit reached") {
			foundOverride = true
			break
		}
	}
	if !foundOverride {
		t.Error("expected System 1 Override message in tape, found none")
	}
}

func TestPanicModeRejectsNonExitToolCalls(t *testing.T) {
	// In panic mode, sh tool calls should be rejected with a message
	// telling the agent to call exit immediately.
	mock := &mockProvider{
		responses: []tape.Message{
			// First response: agent tries sh in panic mode → rejected
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "echo still working",
						},
					},
				},
			},
			// Second response: agent complies and exits
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

	// Set panic mode
	rt.panicMode.Store(true)

	exitCode := rt.Run("some task", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls (rejected sh + exit), got %d", mock.callCount)
	}

	// Verify sh was rejected
	msgs := rt.tape.Messages()
	foundRejection := false
	for _, m := range msgs {
		if m.Role == tape.RoleToolResult && strings.Contains(string(m.StructuredContent), "SIGALRM") {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Error("expected SIGALRM rejection tool result in tape, found none")
	}
}

func TestPanicModeAllowsExitToolCall(t *testing.T) {
	// In panic mode, exit tool calls should still be processed normally.
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "exit",
						Arguments: map[string]any{
							"status": "failure",
							"stderr": "interrupted by SIGALRM",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	// Capture stderr
	oldStderr := rt.stderr
	rErr, wErr, _ := os.Pipe()
	rt.stderr = wErr
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}

	rt.panicMode.Store(true)
	exitCode := rt.Run("task under time pressure", "Begin.")

	wErr.Close()
	errBuf := make([]byte, 4096)
	nErr, _ := rErr.Read(errBuf)
	rErr.Close()
	rt.stderr = oldStderr

	if exitCode != 1 {
		t.Errorf("expected exit code 1 (failure), got %d", exitCode)
	}

	stderr := string(errBuf[:nErr])
	if stderr != "interrupted by SIGALRM" {
		t.Errorf("expected stderr %q, got %q", "interrupted by SIGALRM", stderr)
	}
}

// ---------------------------------------------------------------------------
// SIGINT forwarding / process tracking tests (§2.2)
// ---------------------------------------------------------------------------
