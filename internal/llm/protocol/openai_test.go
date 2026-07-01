package protocol

import (
	"encoding/json"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestOpenAIEncodeRequest_StripsTrailingAssistantPrefillWhenRequested(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{Role: tape.RoleAssistant, Content: "2"},
		},
		nil,
		"claude-sonnet-4.6",
		4096,
		RequestOptions{NoAssistantPrefill: true},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Messages) != 1 {
		t.Fatalf("message count = %d, want 1", len(req.Messages))
	}
	if req.Messages[0].Role != "user" || req.Messages[0].Content != "Begin." {
		t.Fatalf("messages[0] = %+v, want trailing assistant removed", req.Messages[0])
	}
}

func TestOpenAIEncodeRequest_StripsIncompleteToolBatch(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{ID: "call_fork", Name: "fork", Arguments: map[string]any{"children": []any{}}},
				},
			},
		},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Messages []struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []any  `json:"tool_calls"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Messages) != 1 {
		t.Fatalf("message count = %d, want 1", len(req.Messages))
	}
	if req.Messages[0].Role != "user" || req.Messages[0].Content != "Begin." {
		t.Fatalf("messages[0] = %+v, want dangling tool batch removed", req.Messages[0])
	}
}

func TestOpenAIEncodeRequest_PreservesCompletedToolBatch(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:           "call_sh",
						Name:         "sh",
						Arguments:    map[string]any{"command": "pwd"},
						ExtraContent: json.RawMessage(`{"google":{"thought_signature":"sig-a"}}`),
					},
				},
			},
			{Role: tape.RoleToolResult, ToolID: "call_sh", Content: `{"status":"completed"}`},
		},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Messages []struct {
			Role       string `json:"role"`
			ToolCallID string `json:"tool_call_id"`
			ToolCalls  []struct {
				ID           string          `json:"id"`
				ExtraContent json.RawMessage `json:"extra_content"`
			} `json:"tool_calls"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Messages) != 3 {
		t.Fatalf("message count = %d, want 3", len(req.Messages))
	}
	if len(req.Messages[1].ToolCalls) != 1 || req.Messages[1].ToolCalls[0].ID != "call_sh" {
		t.Fatalf("assistant tool batch missing from request: %+v", req.Messages[1])
	}
	if got := string(req.Messages[1].ToolCalls[0].ExtraContent); got != `{"google":{"thought_signature":"sig-a"}}` {
		t.Fatalf("assistant tool extra_content = %s, want thought signature preserved", got)
	}
	if req.Messages[2].Role != "tool" || req.Messages[2].ToolCallID != "call_sh" {
		t.Fatalf("tool result missing from request: %+v", req.Messages[2])
	}
}

func TestOpenAIEncodeRequest_NormalizesNilToolArguments(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{{
					ID:   "call_empty",
					Name: "empty",
				}},
			},
			{Role: tape.RoleToolResult, ToolID: "call_empty"},
		},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Messages []struct {
			ToolCalls []struct {
				Function struct {
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Messages) != 3 || len(req.Messages[1].ToolCalls) != 1 {
		t.Fatalf("completed tool batch missing from request: %s", string(body))
	}
	if got := req.Messages[1].ToolCalls[0].Function.Arguments; got != "{}" {
		t.Fatalf("tool arguments = %q, want {}", got)
	}
}

func TestOpenAIEncodeRequest_OmitsKimiThinkingForGenericOpenAI(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleUser, Content: "Begin."}},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{ThinkingBudget: "high"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if _, ok := req["thinking"]; ok {
		t.Fatalf("generic OpenAI request should omit kimi-only thinking field: %s", string(body))
	}
	if got := req["reasoning_effort"]; got != "high" {
		t.Fatalf("reasoning_effort = %v, want high", got)
	}
}

func TestOpenAIEncodeRequest_KeepsKimiThinkingForKimiModels(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleUser, Content: "Begin."}},
		nil,
		"kimi-k2.5",
		4096,
		RequestOptions{ThinkingBudget: "high"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		ReasoningEffort string `json:"reasoning_effort"`
		Thinking        *struct {
			Type string `json:"type"`
		} `json:"thinking"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if req.ReasoningEffort != "high" {
		t.Fatalf("reasoning_effort = %q, want high", req.ReasoningEffort)
	}
	if req.Thinking == nil || req.Thinking.Type != "enabled" {
		t.Fatalf("thinking = %+v, want enabled", req.Thinking)
	}
}

func TestOpenAIEncodeRequest_AllowsXHighReasoningEffort(t *testing.T) {
	body, err := (&OpenAIProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleUser, Content: "Begin."}},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{ThinkingBudget: "xhigh"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got := req["reasoning_effort"]; got != "xhigh" {
		t.Fatalf("reasoning_effort = %v, want xhigh", got)
	}
}

func TestOpenAIEncodeRequest_AssistantReasoningEcho(t *testing.T) {
	convo := []tape.Message{
		{Role: tape.RoleUser, Content: "Begin."},
		{Role: tape.RoleAssistant, Content: "Done.", ReasoningContent: "deliberating at length"},
	}

	encodedReasoning := func(t *testing.T, strip bool) (string, bool) {
		t.Helper()
		body, err := (&OpenAIProtocol{}).EncodeRequest(
			convo, nil, "glm-5.2", 4096,
			RequestOptions{StripAssistantReasoning: strip},
		)
		if err != nil {
			t.Fatalf("EncodeRequest() error = %v", err)
		}
		var req struct {
			Messages []struct {
				Role             string `json:"role"`
				ReasoningContent string `json:"reasoning_content"`
			} `json:"messages"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}
		for _, m := range req.Messages {
			if m.Role == "assistant" {
				return m.ReasoningContent, true
			}
		}
		t.Fatalf("no assistant message in request: %s", string(body))
		return "", false
	}

	t.Run("echoed when not stripped", func(t *testing.T) {
		got, _ := encodedReasoning(t, false)
		if got != "deliberating at length" {
			t.Fatalf("reasoning_content = %q, want it echoed for default providers", got)
		}
	})

	t.Run("stripped when flagged", func(t *testing.T) {
		got, _ := encodedReasoning(t, true)
		if got != "" {
			t.Fatalf("reasoning_content = %q, want it omitted when StripAssistantReasoning is set", got)
		}
	})
}

func TestParseOpenAIResponse_PrefersChoiceWithToolCalls(t *testing.T) {
	resp := openaiResponse{
		Choices: []openaiChoice{
			{
				FinishReason: "tool_calls",
				Message: openaiResponseMessage{
					Role:    "assistant",
					Content: "Let me examine these more carefully:",
				},
			},
			{
				FinishReason: "tool_calls",
				Message: openaiResponseMessage{
					Role: "assistant",
					ToolCalls: []openaiToolCall{
						{
							ID:           "toolu_123",
							Type:         "function",
							ExtraContent: json.RawMessage(`{"google":{"thought_signature":"sig-b"}}`),
							Function: openaiFunctionCall{
								Name:      "sh",
								Arguments: `{"command":"pwd"}`,
							},
						},
					},
				},
			},
		},
	}

	msg := parseOpenAIResponse(resp)
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("tool call count = %d, want 1", len(msg.ToolCalls))
	}
	if msg.ToolCalls[0].Name != "sh" {
		t.Fatalf("tool name = %q, want sh", msg.ToolCalls[0].Name)
	}
	if got := msg.ToolCalls[0].Arguments["command"]; got != "pwd" {
		t.Fatalf("tool argument command = %v, want pwd", got)
	}
	if got := string(msg.ToolCalls[0].ExtraContent); got != `{"google":{"thought_signature":"sig-b"}}` {
		t.Fatalf("tool extra_content = %s, want thought signature preserved", got)
	}
}

func TestParseOpenAIResponse_ExtractsGeminiThoughtContent(t *testing.T) {
	resp := openaiResponse{
		Choices: []openaiChoice{
			{
				FinishReason: "stop",
				Message: openaiResponseMessage{
					Role:         "assistant",
					Content:      "<thought>Inspect the constraints first.</thought>READY",
					ExtraContent: json.RawMessage(`{"google":{"thought":true}}`),
				},
			},
		},
	}

	msg := parseOpenAIResponse(resp)
	if msg.Content != "READY" {
		t.Fatalf("content = %q, want READY", msg.Content)
	}
	if msg.ReasoningContent != "Inspect the constraints first." {
		t.Fatalf("reasoning = %q, want extracted Gemini thought", msg.ReasoningContent)
	}
}

func TestParseOpenAIResponse_ExtractsGeminiThoughtOnlyToolCall(t *testing.T) {
	resp := openaiResponse{
		Choices: []openaiChoice{
			{
				FinishReason: "tool_calls",
				Message: openaiResponseMessage{
					Role:         "assistant",
					Content:      "<thought>Use the shell before answering.</thought>",
					ExtraContent: json.RawMessage(`{"google":{"thought":true}}`),
					ToolCalls: []openaiToolCall{
						{
							ID:           "toolu_456",
							Type:         "function",
							ExtraContent: json.RawMessage(`{"google":{"thought_signature":"sig-c"}}`),
							Function: openaiFunctionCall{
								Name:      "sh",
								Arguments: `{"command":"echo hi"}`,
							},
						},
					},
				},
			},
		},
	}

	msg := parseOpenAIResponse(resp)
	if msg.Content != "" {
		t.Fatalf("content = %q, want empty visible content after stripping thought block", msg.Content)
	}
	if msg.ReasoningContent != "Use the shell before answering." {
		t.Fatalf("reasoning = %q, want extracted Gemini thought", msg.ReasoningContent)
	}
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("tool call count = %d, want 1", len(msg.ToolCalls))
	}
	if got := string(msg.ToolCalls[0].ExtraContent); got != `{"google":{"thought_signature":"sig-c"}}` {
		t.Fatalf("tool extra_content = %s, want thought signature preserved", got)
	}
}
