package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
		ModelID:              "claude-sonnet-4-20250514",
		APIKey:               "test-key",
		Provider:             "anthropic",
		MaxDepth:             5,
		Depth:                0,
		SessionID:            "test-1234-5678",
		TapeID:               "tape-1234-5678",
		MaxConcurrent:        20,
		MaxAgents:            10,
		ShTimeout:            10,
		OutputTruncate:       20480,
		DataDir:              t.TempDir(),
		Shell:                "/bin/sh",
		MaxTurns:             0, // disabled for existing tests
		TurnExhaustionPolicy: config.TurnExhaustionHardFail,
		PromptMetaphor:       config.PromptMetaphorOff,
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
							"summary": "checkpoint",
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
		if strings.Contains(m.Content, "[MARK]") {
			foundMark = true
		}
		if strings.Contains(m.Content, "[MEMORY META]") {
			foundMeta = true
		}
	}
	if !foundMark {
		t.Fatal("expected mark tool_result on tape")
	}
	if !foundMeta {
		t.Fatal("expected memory meta feedback in tool_result")
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

func TestRunSyncsAgentRoot(t *testing.T) {
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

	if exitCode := rt.Run("sync agent root", "Begin."); exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	statusPath := cfg.AgentRoot() + "/status/session.json"
	statusData, err := os.ReadFile(statusPath)
	if err != nil {
		t.Fatalf("read status file: %v", err)
	}
	var status map[string]any
	if err := json.Unmarshal(statusData, &status); err != nil {
		t.Fatalf("unmarshal status file: %v", err)
	}
	if status["session_id"] != cfg.SessionID {
		t.Fatalf("status.session_id = %v, want %q", status["session_id"], cfg.SessionID)
	}
	if status["tape_id"] != cfg.TapeID {
		t.Fatalf("status.tape_id = %v, want %q", status["tape_id"], cfg.TapeID)
	}

	missionData, err := os.ReadFile(cfg.AgentRoot() + "/mission.txt")
	if err != nil {
		t.Fatalf("read mission file: %v", err)
	}
	if strings.TrimSpace(string(missionData)) != "sync agent root" {
		t.Fatalf("mission.txt = %q, want %q", strings.TrimSpace(string(missionData)), "sync agent root")
	}

	currentTape, err := os.Readlink(cfg.AgentRoot() + "/log/current.jsonl")
	if err != nil {
		t.Fatalf("read current tape symlink: %v", err)
	}
	if currentTape != cfg.TapePath("") {
		t.Fatalf("current tape symlink = %q, want %q", currentTape, cfg.TapePath(""))
	}

	tapeDir, err := os.Readlink(cfg.AgentRoot() + "/log/tapes")
	if err != nil {
		t.Fatalf("read tapes symlink: %v", err)
	}
	if tapeDir != cfg.TapeDir("") {
		t.Fatalf("log/tapes symlink = %q, want %q", tapeDir, cfg.TapeDir(""))
	}

	runtimeLogTarget, err := os.Readlink(cfg.AgentRoot() + "/log/runtime.log")
	if err != nil {
		t.Fatalf("read runtime log symlink: %v", err)
	}
	runtimeLogData, err := os.ReadFile(runtimeLogTarget)
	if err != nil {
		t.Fatalf("read runtime log: %v", err)
	}
	runtimeLog := string(runtimeLogData)
	if strings.Contains(runtimeLog, "turn ") {
		t.Fatalf("runtime.log should not contain turn-level entries, got %q", runtimeLog)
	}
	if strings.Contains(runtimeLog, "THE PRIME DIRECTIVE") {
		t.Fatalf("runtime.log should not duplicate system prompt content, got %q", runtimeLog)
	}

	jobsDir, err := os.Readlink(cfg.AgentRoot() + "/jobs")
	if err != nil {
		t.Fatalf("read jobs symlink: %v", err)
	}
	if jobsDir != cfg.JobSessionDir("") {
		t.Fatalf("jobs symlink = %q, want %q", jobsDir, cfg.JobSessionDir(""))
	}
}

func TestReplaceSymlinkResolvesRelativeTarget(t *testing.T) {
	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir tmp: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	linkPath := filepath.Join(".quine", "agent", "session-1", "jobs")
	targetPath := filepath.Join(".quine", "jobs", "session-1")
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		t.Fatalf("mkdir link parent: %v", err)
	}
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	if err := replaceSymlink(linkPath, targetPath); err != nil {
		t.Fatalf("replaceSymlink: %v", err)
	}

	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("eval symlink: %v", err)
	}
	want, err := filepath.Abs(targetPath)
	if err != nil {
		t.Fatalf("abs target: %v", err)
	}
	if resolved != want {
		t.Fatalf("resolved symlink = %q, want %q", resolved, want)
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

func TestExecutionBudgetHardFailTerminatesImmediately(t *testing.T) {
	// With hard_fail policy, reaching MaxTurns should terminate immediately
	// without an extra continuation inference.
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
	cfg.TurnExhaustionPolicy = config.TurnExhaustionHardFail
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")

	if exitCode != 1 {
		t.Errorf("expected exit code 1 for execution exhaustion, got %d", exitCode)
	}

	// Only first two responses should be used (no continuation inference).
	if mock.callCount != 2 {
		t.Errorf("expected exactly 2 LLM calls with hard_fail policy, got %d", mock.callCount)
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
}

func TestExecutionBudgetNearDeathContinuation(t *testing.T) {
	// With near_death_exec policy, reaching MaxTurns opens one final inference
	// window where only exec is accepted.
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
	cfg.TurnExhaustionPolicy = config.TurnExhaustionNearDeathExec
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("do something", "Begin.")
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if mock.callCount != 3 {
		t.Errorf("expected exactly 3 LLM calls (2 sh + continuation), got %d", mock.callCount)
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
		if m.Role == tape.RoleToolResult {
			if strings.Contains(m.Content, "[EXECUTIONS LEFT]") {
				foundBudget = true
				if !strings.Contains(m.Content, "[EXECUTIONS LEFT] 2") {
					t.Errorf("expected '[EXECUTIONS LEFT] 2' in tool result, got %q", m.Content)
				}
			}
			if strings.Contains(m.Content, "[CONTEXT USED]") {
				foundContext = true
			}
		}
	}
	if !foundBudget {
		t.Error("expected [EXECUTIONS LEFT] in a tool result message, found none")
	}
	if !foundContext {
		t.Error("expected [CONTEXT USED] in a tool result message, found none")
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
		if m.Role == tape.RoleToolResult {
			if strings.Contains(m.Content, "[EXECUTIONS LEFT]") {
				t.Error("did not expect [EXECUTIONS LEFT] in tool results when MaxTurns=0")
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

// ---------------------------------------------------------------------------
// stdin parameter tests (LLM-behaviour: does the runtime wire stdin through?)
// ---------------------------------------------------------------------------

// TestShStdinParameter verifies that when the LLM sends a sh tool call with
// a "stdin" argument, the runtime parses it and passes it to Execute so the
// data is piped into the command.
func TestShStdinParameter(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := tmpDir + "/result.txt"

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "cat > " + outFile,
							"stdin":   "hello from llm stdin\n",
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

	exitCode := rt.Run("write a file via stdin", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "hello from llm stdin") {
		t.Errorf("file content = %q, expected 'hello from llm stdin'", string(data))
	}
}

// TestShStdinSpecialChars verifies that characters that would break shell
// quoting (quotes, backslashes, dollar signs) pass through unchanged.
func TestShStdinSpecialChars(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := tmpDir + "/special.txt"

	tricky := "key = \"value\"\npath = C:\\data\nprice = $100\n"

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "cat > " + outFile,
							"stdin":   tricky,
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

	exitCode := rt.Run("write a file with special chars", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if string(data) != tricky {
		t.Errorf("content mismatch.\ngot:  %q\nwant: %q", string(data), tricky)
	}
}
