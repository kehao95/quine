package runtime

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
	"github.com/kehao95/quine/internal/tape"
)

// mockProvider is a test double that returns pre-programmed responses.
type mockProvider struct {
	responses []tape.Message
	callCount int
}

func (m *mockProvider) Generate(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	if m.callCount >= len(m.responses) {
		return tape.Message{}, llm.Usage{}, fmt.Errorf("mock: no more responses (call %d)", m.callCount)
	}
	resp := m.responses[m.callCount]
	m.callCount++
	return resp, llm.Usage{InputTokens: 100, OutputTokens: 50}, nil
}

func (m *mockProvider) ContextWindowSize() int { return 200000 }

// mockErrorProvider returns errors.
type mockErrorProvider struct {
	err error
}

func (m *mockErrorProvider) Generate(msgs []tape.Message, tools []llm.ToolSchema) (tape.Message, llm.Usage, error) {
	return tape.Message{}, llm.Usage{}, m.err
}

func (m *mockErrorProvider) ContextWindowSize() int { return 200000 }

func testCfg(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		ModelID:        "claude-sonnet-4-20250514",
		APIKey:         "test-key",
		Provider:       "anthropic",
		MaxDepth:       5,
		Depth:          0,
		SessionID:      "test-1234-5678",
		MaxConcurrent:  20,
		ShTimeout:      10,
		OutputTruncate: 20480,
		DataDir:        t.TempDir(),
		Shell:          "/bin/sh",
		MaxTurns:       0, // unlimited for existing tests
	}
}

// silenceRuntime suppresses all runtime output for clean test output.
func silenceRuntime(rt *Runtime) {
	devnull, _ := os.Open(os.DevNull)
	rt.stderr = devnull
	rt.log = func(format string, args ...any) {}
	rt.logError = func(format string, args ...any) {}
	if rt.logFile != nil {
		rt.logFile.Close()
		rt.logFile = nil
	}
}

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

func TestTurnLimitKillsProcess(t *testing.T) {
	// Set MaxTurns=2. Agent does sh (turn 1), then sh (turn 2).
	// After turn 2, the near-death warning fires and the agent gets ONE
	// final inference. If it doesn't call exec, it dies.
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
			// Near-death inference: agent calls exit (not exec), so it dies.
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

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for turn exhaustion, got %d", exitCode)
	}

	// 3 LLM calls: turn 1 sh, turn 2 sh, near-death inference
	if mock.callCount != 3 {
		t.Errorf("expected exactly 3 LLM calls (2 turns + near-death), got %d", mock.callCount)
	}

	if rt.tape.Outcome == nil {
		t.Fatal("expected outcome to be set")
	}
	if rt.tape.Outcome.TerminationMode != tape.TermTurnExhaustion {
		t.Errorf("expected termination mode %q, got %q",
			tape.TermTurnExhaustion, rt.tape.Outcome.TerminationMode)
	}
	if !strings.Contains(rt.tape.Outcome.Stderr, "turn limit exhausted") {
		t.Errorf("expected stderr to contain 'turn limit exhausted', got %q", rt.tape.Outcome.Stderr)
	}
}

func TestTurnLimitFeedbackMessages(t *testing.T) {
	// With MaxTurns=3, after sh (turn 1) the tool result should contain
	// "[TURNS LEFT] 2" appended to it.
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

	// Verify "[TURNS LEFT] 2" and "[CONTEXT USED]" were appended to the sh tool result
	msgs := rt.tape.Messages()
	foundBudget := false
	foundContext := false
	for _, m := range msgs {
		if m.Role == tape.RoleToolResult {
			if strings.Contains(m.Content, "[TURNS LEFT]") {
				foundBudget = true
				if !strings.Contains(m.Content, "[TURNS LEFT] 2") {
					t.Errorf("expected '[TURNS LEFT] 2' in tool result, got %q", m.Content)
				}
			}
			if strings.Contains(m.Content, "[CONTEXT USED]") {
				foundContext = true
			}
		}
	}
	if !foundBudget {
		t.Error("expected [TURNS LEFT] in a tool result message, found none")
	}
	if !foundContext {
		t.Error("expected [CONTEXT USED] in a tool result message, found none")
	}
}

func TestTurnLimitZeroMeansUnlimited(t *testing.T) {
	// MaxTurns=0 should not inject any turn budget into tool results
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

	// Verify NO [TURNS LEFT] in any tool result (but [CONTEXT USED] should still be present)
	msgs := rt.tape.Messages()
	foundContext := false
	for _, m := range msgs {
		if m.Role == tape.RoleToolResult {
			if strings.Contains(m.Content, "[TURNS LEFT]") {
				t.Error("did not expect [TURNS LEFT] in tool results when MaxTurns=0")
			}
			if strings.Contains(m.Content, "[CONTEXT USED]") {
				foundContext = true
			}
		}
	}
	if !foundContext {
		t.Error("expected [CONTEXT USED] in tool result even when MaxTurns=0")
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
		if m.Role == tape.RoleToolResult && strings.Contains(m.Content, "Exit rejected") {
			foundRejection = true
			break
		}
	}
	if !foundRejection {
		t.Error("expected a rejection tool result in the tape, found none")
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
		if m.Role == tape.RoleToolResult && strings.Contains(m.Content, "Exit rejected") {
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
		if m.Role == tape.RoleToolResult && strings.Contains(m.Content, "SIGALRM") {
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

func TestProcessTrackingCallbacks(t *testing.T) {
	// Verify that shExecutor's ProcessStarted/ProcessEnded callbacks
	// are wired correctly to Runtime's activeProcess tracking.
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

	exitCode := rt.Run("test process tracking", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// After Run completes, activeProcess should be nil
	if proc := rt.activeProcess.Load(); proc != nil {
		t.Errorf("expected activeProcess to be nil after Run, got pid=%d", proc.Pid)
	}
}
