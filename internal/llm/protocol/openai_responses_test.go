package protocol

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestOpenAIResponsesEncodeRequest_RequestsDetailedReasoningSummary(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleUser, Content: "Begin."}},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{ThinkingBudget: "high"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Reasoning struct {
			Effort  string `json:"effort"`
			Summary string `json:"summary"`
		} `json:"reasoning"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if req.Reasoning.Effort != "high" {
		t.Fatalf("reasoning.effort = %q, want %q", req.Reasoning.Effort, "high")
	}
	if req.Reasoning.Summary != "detailed" {
		t.Fatalf("reasoning.summary = %q, want %q", req.Reasoning.Summary, "detailed")
	}
}

func TestOpenAIResponsesDecodeResponse_HandlesSSEKeepaliveComment(t *testing.T) {
	body := strings.Join([]string{
		": keepalive",
		"",
		"event: response.completed",
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_123\",\"output\":[{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Hello from comment-prefixed SSE\"}]}],\"usage\":{\"input_tokens\":12,\"output_tokens\":7}}}",
		"",
	}, "\n")

	msg, usage, err := (&OpenAIResponsesProtocol{}).DecodeResponse([]byte(body))
	if err != nil {
		t.Fatalf("DecodeResponse() error = %v", err)
	}
	if msg.Content != "Hello from comment-prefixed SSE" {
		t.Fatalf("content = %q", msg.Content)
	}
	if usage.InputTokens != 12 || usage.OutputTokens != 7 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestOpenAIResponsesDecodeResponse_SynthesizesCompletedFromOutputTextDelta(t *testing.T) {
	body := "event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"}\n\n"

	msg, _, err := (&OpenAIResponsesProtocol{}).DecodeResponse([]byte(body))
	if err != nil {
		t.Fatalf("DecodeResponse() error = %v", err)
	}
	if msg.Content != "partial" {
		t.Fatalf("content = %q", msg.Content)
	}
}

func TestOpenAIResponsesDecodeResponse_SynthesizesCompletedFromDataOnlyChatChunks(t *testing.T) {
	body := strings.Join([]string{
		"data: {\"id\":\"chatcmpl_123\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"delta\":{\"content\":\"hel\"}}]}",
		"",
		"data: {\"id\":\"chatcmpl_123\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"delta\":{\"content\":\"lo\"}}]}",
		"",
		"data: [DONE]",
		"",
	}, "\n")

	msg, _, err := (&OpenAIResponsesProtocol{}).DecodeResponse([]byte(body))
	if err != nil {
		t.Fatalf("DecodeResponse() error = %v", err)
	}
	if msg.Content != "hello" {
		t.Fatalf("content = %q", msg.Content)
	}
}

func TestOpenAIResponsesDecodeResponse_ReconstructsOutputFromOutputItemDone(t *testing.T) {
	body := strings.Join([]string{
		"event: response.output_item.done",
		"data: {\"type\":\"response.output_item.done\",\"item\":{\"id\":\"msg_123\",\"type\":\"message\",\"status\":\"completed\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Recovered from output_item.done\"}]},\"output_index\":1,\"sequence_number\":9}",
		"",
		"event: response.completed",
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_123\",\"output\":[],\"usage\":{\"input_tokens\":17,\"output_tokens\":45}}}",
		"",
	}, "\n")

	msg, usage, err := (&OpenAIResponsesProtocol{}).DecodeResponse([]byte(body))
	if err != nil {
		t.Fatalf("DecodeResponse() error = %v", err)
	}
	if msg.Content != "Recovered from output_item.done" {
		t.Fatalf("content = %q", msg.Content)
	}
	if usage.InputTokens != 17 || usage.OutputTokens != 45 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestOpenAIResponsesDecodeResponse_MergesReasoningSummaryFromSSE(t *testing.T) {
	body := strings.Join([]string{
		"event: response.reasoning_summary_text.done",
		"data: {\"type\":\"response.reasoning_summary_text.done\",\"item_id\":\"rs_123\",\"summary_index\":0,\"text\":\"First summary.\",\"output_index\":0,\"sequence_number\":1}",
		"",
		"event: response.reasoning_summary_text.done",
		"data: {\"type\":\"response.reasoning_summary_text.done\",\"item_id\":\"rs_123\",\"summary_index\":1,\"text\":\"Second summary.\",\"output_index\":0,\"sequence_number\":2}",
		"",
		"event: response.completed",
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_123\",\"output\":[{\"id\":\"rs_123\",\"type\":\"reasoning\"},{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Done\"}]}],\"usage\":{\"input_tokens\":12,\"output_tokens\":7}}}",
		"",
	}, "\n")

	msg, _, err := (&OpenAIResponsesProtocol{}).DecodeResponse([]byte(body))
	if err != nil {
		t.Fatalf("DecodeResponse() error = %v", err)
	}
	if msg.Content != "Done" {
		t.Fatalf("content = %q", msg.Content)
	}
	if len(msg.ReasoningItems) != 1 {
		t.Fatalf("reasoning item count = %d", len(msg.ReasoningItems))
	}
	if got := msg.ReasoningItems[0].ID; got != "rs_123" {
		t.Fatalf("reasoning id = %q", got)
	}
	want := []string{"First summary.", "Second summary."}
	if !reflect.DeepEqual(msg.ReasoningItems[0].Summary, want) {
		t.Fatalf("reasoning summary = %#v, want %#v", msg.ReasoningItems[0].Summary, want)
	}
}

func TestOpenAIResponsesEncodeRequest_SerializesServiceTier(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleUser, Content: "Begin."}},
		nil,
		"gpt-5.5",
		4096,
		RequestOptions{ThinkingBudget: "medium", ServiceTier: "priority"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Model       string `json:"model"`
		ServiceTier string `json:"service_tier"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if req.Model != "gpt-5.5" {
		t.Fatalf("model = %q, want gpt-5.5", req.Model)
	}
	if req.ServiceTier != "priority" {
		t.Fatalf("service_tier = %q, want priority", req.ServiceTier)
	}
}

func TestOpenAIResponsesEncodeRequest_CarriesSystemOnlyInputWithBlankUserItem(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleSystem, Content: "system only"}},
		nil,
		"gpt-5.4",
		4096,
		RequestOptions{ThinkingBudget: "medium"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req responsesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if req.Instructions != "system only" {
		t.Fatalf("instructions = %q, want system only", req.Instructions)
	}
	if req.Input == nil {
		t.Fatalf("input should be serialized, body=%s", string(body))
	}
	if len(req.Input) != 1 {
		t.Fatalf("input len = %d, want 1", len(req.Input))
	}
	if req.Input[0].Role != "user" {
		t.Fatalf("fallback role = %q, want user", req.Input[0].Role)
	}
	if len(req.Input[0].Content) != 1 || req.Input[0].Content[0].Text != " " {
		t.Fatalf("fallback content = %#v, want a blank input carrier", req.Input[0].Content)
	}
}

func TestOpenAIResponsesEncodeRequest_CarriesToolResultImageAsUserInput(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{{
					ID:        "call_vision",
					Name:      "vision",
					Arguments: map[string]any{"path": "/tmp/render.ppm"},
				}},
			},
			{
				Role:              tape.RoleToolResult,
				ToolID:            "call_vision",
				StructuredContent: []byte(`{"tool":"vision","status":"completed"}`),
				Image: &tape.ImagePart{
					MIMEType: "image/png",
					Data:     "aW1n",
				},
			},
		},
		nil,
		"gpt-5.5",
		4096,
		RequestOptions{ThinkingBudget: "medium"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req responsesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Input) != 4 {
		t.Fatalf("input len = %d, want 4; body=%s", len(req.Input), string(body))
	}
	if req.Input[2].Type != "function_call_output" || req.Input[2].CallID != "call_vision" {
		t.Fatalf("tool output item = %#v", req.Input[2])
	}
	if req.Input[3].Role != "user" {
		t.Fatalf("image carrier role = %q, want user", req.Input[3].Role)
	}
	if len(req.Input[3].Content) != 2 {
		t.Fatalf("image carrier content len = %d, want 2", len(req.Input[3].Content))
	}
	if req.Input[3].Content[0].Type != "input_text" {
		t.Fatalf("carrier text type = %q, want input_text", req.Input[3].Content[0].Type)
	}
	if req.Input[3].Content[1].Type != "input_image" {
		t.Fatalf("carrier image type = %q, want input_image", req.Input[3].Content[1].Type)
	}
	if got := req.Input[3].Content[1].ImageURL; got != "data:image/png;base64,aW1n" {
		t.Fatalf("image_url = %q", got)
	}
}

func TestOpenAIResponsesEncodeRequest_PreservesEmptyFunctionCallOutput(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
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
		"gpt-5.5",
		4096,
		RequestOptions{ThinkingBudget: "medium"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Input []map[string]json.RawMessage `json:"input"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Input) != 3 {
		t.Fatalf("input len = %d, want 3; body=%s", len(req.Input), string(body))
	}
	if got := string(req.Input[1]["arguments"]); got != `"{}"` {
		t.Fatalf("function_call arguments = %s, want empty object string; body=%s", got, string(body))
	}
	output, ok := req.Input[2]["output"]
	if !ok {
		t.Fatalf("function_call_output missing output field; body=%s", string(body))
	}
	if got := string(output); got != `""` {
		t.Fatalf("function_call_output output = %s, want empty string", got)
	}
	if _, ok := req.Input[2]["content"]; ok {
		t.Fatalf("function_call_output should not serialize message content fields: %s", string(body))
	}
}

func TestOpenAIResponsesEncodeRequest_PreservesEmptyToolParameters(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
		[]tape.Message{{Role: tape.RoleUser, Content: "Begin."}},
		[]ToolSchema{{
			Name:       "empty_schema",
			Parameters: map[string]any{},
		}},
		"gpt-5.5",
		4096,
		RequestOptions{ThinkingBudget: "medium"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req struct {
		Tools []map[string]json.RawMessage `json:"tools"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Tools) != 1 {
		t.Fatalf("tool count = %d, want 1; body=%s", len(req.Tools), string(body))
	}
	params, ok := req.Tools[0]["parameters"]
	if !ok {
		t.Fatalf("tool parameters omitted; body=%s", string(body))
	}
	if got := string(params); got != `{}` {
		t.Fatalf("tool parameters = %s, want {}", got)
	}
}

func TestOpenAIResponsesEncodeRequest_CarriesBlankTextForImageOnlyToolResult(t *testing.T) {
	body, err := (&OpenAIResponsesProtocol{}).EncodeRequest(
		[]tape.Message{
			{Role: tape.RoleUser, Content: "Begin."},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{{
					ID:        "call_vision",
					Name:      "vision",
					Arguments: map[string]any{},
				}},
			},
			{
				Role:   tape.RoleToolResult,
				ToolID: "call_vision",
				Image: &tape.ImagePart{
					MIMEType: "image/png",
					Data:     "aW1n",
				},
			},
		},
		nil,
		"gpt-5.5",
		4096,
		RequestOptions{ThinkingBudget: "medium"},
	)
	if err != nil {
		t.Fatalf("EncodeRequest() error = %v", err)
	}

	var req responsesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(req.Input) != 4 {
		t.Fatalf("input len = %d, want 4; body=%s", len(req.Input), string(body))
	}
	if len(req.Input[3].Content) != 2 {
		t.Fatalf("image carrier content len = %d, want 2", len(req.Input[3].Content))
	}
	if got := req.Input[3].Content[0].Text; got != " " {
		t.Fatalf("image carrier text = %q, want blank carrier", got)
	}
}
