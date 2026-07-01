package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

// OpenAIResponsesProtocol implements Protocol for OpenAI's Responses API.
// This is the newer API shape used by Codex and recommended for agentic workloads.
// Endpoint: POST /v1/responses
type OpenAIResponsesProtocol struct{}

// ---------------------------------------------------------------------------
// API request/response types
// ---------------------------------------------------------------------------

type responsesRequest struct {
	Model        string              `json:"model"`
	Instructions string              `json:"instructions,omitempty"`
	Input        []responsesItem     `json:"input"`
	Tools        []responsesTool     `json:"tools,omitempty"`
	Store        bool                `json:"store"`               // Must be false for Codex models
	Stream       bool                `json:"stream"`              // Must be true for Codex models
	Reasoning    *responsesReasoning `json:"reasoning,omitempty"` // Reasoning configuration
	ServiceTier  string              `json:"service_tier,omitempty"`
}

// responsesReasoning configures reasoning behavior for models that support it
type responsesReasoning struct {
	// Effort controls reasoning depth: none, minimal, low, medium, high, xhigh
	Effort string `json:"effort,omitempty"`
	// Summary controls reasoning summary output: auto, concise, detailed
	Summary string `json:"summary,omitempty"`
}

// responsesTool uses the flat format required by Responses API
// (different from Chat Completions which uses wrapped {type, function: {...}})
type responsesTool struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters"`
	Strict      bool           `json:"strict,omitempty"`
}

type responsesItem struct {
	Type      string             `json:"type,omitempty"`
	Role      string             `json:"role,omitempty"`
	Content   []responsesContent `json:"content,omitempty"`
	CallID    string             `json:"call_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Arguments string             `json:"arguments,omitempty"`
	Output    string             `json:"output,omitempty"`
}

func (i responsesItem) MarshalJSON() ([]byte, error) {
	switch i.Type {
	case "function_call":
		return json.Marshal(struct {
			Type      string `json:"type"`
			CallID    string `json:"call_id"`
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		}{
			Type:      i.Type,
			CallID:    i.CallID,
			Name:      i.Name,
			Arguments: i.Arguments,
		})
	case "function_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			CallID string `json:"call_id"`
			Output string `json:"output"`
		}{
			Type:   i.Type,
			CallID: i.CallID,
			Output: i.Output,
		})
	default:
		type alias responsesItem
		return json.Marshal(alias(i))
	}
}

type responsesContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type responsesResponse struct {
	ID     string            `json:"id"`
	Output []responsesOutput `json:"output"`
	Usage  responsesUsage    `json:"usage"`
}

type responsesOutput struct {
	ID        string             `json:"id,omitempty"`
	Type      string             `json:"type"`
	Role      string             `json:"role,omitempty"`
	Content   []responsesContent `json:"content,omitempty"`
	Summary   []responsesSummary `json:"summary,omitempty"` // For reasoning items
	CallID    string             `json:"call_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Arguments string             `json:"arguments,omitempty"`
}

type responsesSummary struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type responsesUsage struct {
	InputTokens      int `json:"input_tokens"`
	OutputTokens     int `json:"output_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type responsesError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// ---------------------------------------------------------------------------
// Protocol implementation
// ---------------------------------------------------------------------------

func (p *OpenAIResponsesProtocol) ContentType() string {
	return "application/json"
}

func (p *OpenAIResponsesProtocol) EndpointPath() string {
	return "/v1/responses"
}

func (p *OpenAIResponsesProtocol) EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int, opts RequestOptions) ([]byte, error) {
	input, instructions := convertResponsesInput(messages)
	if len(input) == 0 {
		input = []responsesItem{{
			Role: "user",
			Content: []responsesContent{
				{Type: "input_text", Text: " "},
			},
		}}
	}

	req := responsesRequest{
		Model:        model,
		Instructions: instructions,
		Input:        input,
		Tools:        convertResponsesTools(tools),
		Store:        false, // Required for Codex models
		Stream:       true,  // Required for Codex models
		Reasoning: &responsesReasoning{
			Summary: "detailed", // Copilot Responses has proven unreliable about emitting readable summaries under "concise"
		},
		ServiceTier: opts.ServiceTier,
	}

	// Apply thinking budget if specified
	if opts.ThinkingBudget != "" && req.Reasoning != nil {
		// Map our budget levels to Responses API effort levels
		switch opts.ThinkingBudget {
		case "off":
			req.Reasoning.Effort = "none"
		case "low":
			req.Reasoning.Effort = "low"
		case "medium":
			req.Reasoning.Effort = "medium"
		case "high":
			req.Reasoning.Effort = "high"
		case "xhigh":
			req.Reasoning.Effort = "xhigh"
		}
	}

	return json.Marshal(req)
}

func (p *OpenAIResponsesProtocol) DecodeResponse(body []byte) (tape.Message, Usage, error) {
	if looksLikeResponsesSSE(body) {
		var err error
		body, err = extractSSECompletedResponse(bytes.NewReader(body))
		if err != nil {
			return tape.Message{}, Usage{}, fmt.Errorf("parsing SSE stream: %w", err)
		}
	}
	return decodeCompletedResponsesBody(body)
}

// DecodeStream parses a Responses API SSE stream incrementally, invoking
// onDelta as each text/reasoning/tool-call delta arrives, and returns the
// authoritative final (message, usage). The final cell is derived from the same
// completed-event assembly the buffered DecodeResponse path uses, so the two
// paths crystallize an equivalent cell (reproducibility boundary).
func (p *OpenAIResponsesProtocol) DecodeStream(r io.Reader, onDelta func(StreamDelta)) (tape.Message, Usage, error) {
	body, err := extractSSECompletedResponseStream(r, onDelta)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("parsing SSE stream: %w", err)
	}
	return decodeCompletedResponsesBody(body)
}

// decodeCompletedResponsesBody turns a (possibly SSE-wrapped) completed
// Responses payload into a tape message + usage. It is the shared final-assembly
// tail of both DecodeResponse (buffered) and DecodeStream (streaming).
func decodeCompletedResponsesBody(body []byte) (tape.Message, Usage, error) {
	// Handle SSE event wrapper: {"type":"response.completed","response":{...}}
	var wrapper struct {
		Type     string            `json:"type"`
		Response responsesResponse `json:"response"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	// Use the nested response if this is an SSE event wrapper
	var resp responsesResponse
	if wrapper.Type == "response.completed" {
		resp = wrapper.Response
	} else {
		// Fall back to direct response format (non-streaming)
		if err := json.Unmarshal(body, &resp); err != nil {
			return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
		}
	}

	msg := parseResponsesOutput(resp)
	usage := Usage{
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
	}
	// Fallback to legacy field names if new ones are zero
	if usage.InputTokens == 0 && resp.Usage.PromptTokens != 0 {
		usage.InputTokens = resp.Usage.PromptTokens
	}
	if usage.OutputTokens == 0 && resp.Usage.CompletionTokens != 0 {
		usage.OutputTokens = resp.Usage.CompletionTokens
	}

	return msg, usage, nil
}

func looksLikeResponsesSSE(body []byte) bool {
	body = bytes.TrimPrefix(body, []byte{0xEF, 0xBB, 0xBF})
	for len(body) > 0 {
		lineEnd := bytes.IndexByte(body, '\n')
		var line []byte
		if lineEnd == -1 {
			line = body
			body = nil
		} else {
			line = body[:lineEnd]
			body = body[lineEnd+1:]
		}
		line = bytes.TrimRight(line, "\r")
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		return bytes.HasPrefix(line, []byte("event:")) ||
			bytes.HasPrefix(line, []byte("data:")) ||
			line[0] == ':'
	}
	return false
}

// extractSSECompletedResponse parses a Responses API SSE stream and extracts
// the final response.completed payload, merging incremental events that some
// vendors omit from the terminal wrapper. It is the buffered entrypoint and
// emits no deltas.
func extractSSECompletedResponse(r io.Reader) ([]byte, error) {
	return extractSSECompletedResponseStream(r, nil)
}

// extractSSECompletedResponseStream is the shared SSE event loop. When onDelta
// is non-nil it emits an incremental StreamDelta per text/reasoning/tool-call
// event as it arrives; in all cases it returns the same merged completed
// payload the buffered path produces, so the final cell is identical whether or
// not deltas were observed.
func extractSSECompletedResponseStream(r io.Reader, onDelta func(StreamDelta)) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var currentEvent string
	var currentData []string
	var lastCompletedData []byte
	var terminalFailure error
	reasoningSummaries := map[string]map[int]string{}
	reasoningSummaryDeltas := map[string]map[int]string{}
	outputItems := map[int]map[string]any{}
	outputTextDeltas := map[int]map[int]string{}

	flushEvent := func() error {
		if len(currentData) == 0 {
			currentEvent = ""
			currentData = nil
			return nil
		}
		data := strings.Join(currentData, "\n")
		if strings.TrimSpace(data) == "[DONE]" {
			currentEvent = ""
			currentData = nil
			return nil
		}
		eventName := currentEvent
		if eventName == "" {
			var evt struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil {
				eventName = evt.Type
			}
		}
		switch eventName {
		case "response.completed":
			lastCompletedData = []byte(data)
		case "response.output_item.done":
			var evt struct {
				OutputIndex int            `json:"output_index"`
				Item        map[string]any `json:"item"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.Item != nil {
				outputItems[evt.OutputIndex] = evt.Item
			}
		case "response.output_text.delta":
			var evt struct {
				OutputIndex  int    `json:"output_index"`
				ContentIndex int    `json:"content_index"`
				Delta        string `json:"delta"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.Delta != "" {
				putOutputTextDelta(outputTextDeltas, evt.OutputIndex, evt.ContentIndex, evt.Delta)
				if onDelta != nil {
					onDelta(StreamDelta{Kind: StreamDeltaText, Text: evt.Delta})
				}
			}
		case "response.function_call_arguments.delta":
			// Display-only: the authoritative tool call is reassembled from
			// response.output_item.done, so this case never touches the merge
			// maps. It exists solely to surface tool-call generation as a live
			// delta.
			if onDelta != nil {
				var evt struct {
					ItemID string `json:"item_id"`
					CallID string `json:"call_id"`
					Delta  string `json:"delta"`
				}
				if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.Delta != "" {
					id := evt.CallID
					if id == "" {
						id = evt.ItemID
					}
					onDelta(StreamDelta{Kind: StreamDeltaToolCall, Text: evt.Delta, ToolID: id})
				}
			}
		case "response.output_text.done":
			var evt struct {
				OutputIndex  int    `json:"output_index"`
				ContentIndex int    `json:"content_index"`
				Text         string `json:"text"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.Text != "" {
				setOutputText(outputTextDeltas, evt.OutputIndex, evt.ContentIndex, evt.Text)
			}
		case "":
			appendChatCompletionChunk(outputTextDeltas, data, onDelta)
		case "error":
			var evt struct {
				Error struct {
					Type    string `json:"type"`
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil {
				terminalFailure = newSSETerminalError("error", evt.Error.Code, evt.Error.Message, data)
			}
		case "response.failed":
			var evt struct {
				Response struct {
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				} `json:"response"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil {
				terminalFailure = newSSETerminalError("response.failed", evt.Response.Error.Code, evt.Response.Error.Message, data)
			}
		case "response.reasoning_summary_text.done":
			var evt struct {
				ItemID       string `json:"item_id"`
				SummaryIndex int    `json:"summary_index"`
				Text         string `json:"text"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.ItemID != "" && evt.Text != "" {
				putReasoningSummaryText(reasoningSummaries, evt.ItemID, evt.SummaryIndex, evt.Text)
			}
		case "response.reasoning_summary_part.done":
			var evt struct {
				ItemID       string `json:"item_id"`
				SummaryIndex int    `json:"summary_index"`
				Part         struct {
					Text string `json:"text"`
				} `json:"part"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.ItemID != "" && evt.Part.Text != "" {
				putReasoningSummaryText(reasoningSummaries, evt.ItemID, evt.SummaryIndex, evt.Part.Text)
			}
		case "response.reasoning_summary_text.delta":
			var evt struct {
				ItemID       string `json:"item_id"`
				SummaryIndex int    `json:"summary_index"`
				Delta        string `json:"delta"`
			}
			if err := json.Unmarshal([]byte(data), &evt); err == nil && evt.ItemID != "" && evt.Delta != "" {
				putReasoningSummaryText(reasoningSummaryDeltas, evt.ItemID, evt.SummaryIndex, evt.Delta)
				if onDelta != nil {
					onDelta(StreamDelta{Kind: StreamDeltaReasoning, Text: evt.Delta})
				}
			}
		}
		currentEvent = ""
		currentData = nil
		return nil
	}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			if err := flushEvent(); err != nil {
				return nil, err
			}
			currentEvent = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			currentData = append(currentData, strings.TrimPrefix(line, "data: "))
		}
		if line == "" {
			if err := flushEvent(); err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning SSE stream: %w", err)
	}
	if err := flushEvent(); err != nil {
		return nil, err
	}

	if lastCompletedData == nil {
		if terminalFailure != nil {
			return nil, terminalFailure
		}
		var err error
		lastCompletedData, err = synthesizeCompletedResponse(outputItems, outputTextDeltas)
		if err != nil {
			return nil, err
		}
	}

	completed, err := mergeCompletedOutputItems(lastCompletedData, outputItems)
	if err != nil {
		return nil, fmt.Errorf("merging output items: %w", err)
	}

	merged, err := mergeReasoningSummaries(completed, reasoningSummaries, reasoningSummaryDeltas)
	if err != nil {
		return nil, fmt.Errorf("merging reasoning summaries: %w", err)
	}

	return merged, nil
}

func putOutputTextDelta(dst map[int]map[int]string, outputIndex, contentIndex int, delta string) {
	if dst[outputIndex] == nil {
		dst[outputIndex] = map[int]string{}
	}
	dst[outputIndex][contentIndex] += delta
}

func setOutputText(dst map[int]map[int]string, outputIndex, contentIndex int, text string) {
	if dst[outputIndex] == nil {
		dst[outputIndex] = map[int]string{}
	}
	dst[outputIndex][contentIndex] = text
}

func appendChatCompletionChunk(dst map[int]map[int]string, data string, onDelta func(StreamDelta)) {
	var chunk struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return
	}
	for _, choice := range chunk.Choices {
		if choice.Delta.Content != "" {
			putOutputTextDelta(dst, 0, 0, choice.Delta.Content)
			if onDelta != nil {
				onDelta(StreamDelta{Kind: StreamDeltaText, Text: choice.Delta.Content})
			}
		}
	}
}

func synthesizeCompletedResponse(outputItems map[int]map[string]any, outputTextDeltas map[int]map[int]string) ([]byte, error) {
	output := make([]any, 0)

	if len(outputItems) > 0 {
		indexes := make([]int, 0, len(outputItems))
		for idx := range outputItems {
			indexes = append(indexes, idx)
		}
		sort.Ints(indexes)
		for _, idx := range indexes {
			output = append(output, outputItems[idx])
		}
	}

	if len(outputTextDeltas) > 0 {
		indexes := make([]int, 0, len(outputTextDeltas))
		for idx := range outputTextDeltas {
			if _, alreadyHaveItem := outputItems[idx]; alreadyHaveItem {
				continue
			}
			indexes = append(indexes, idx)
		}
		sort.Ints(indexes)

		for _, idx := range indexes {
			partsByIndex := outputTextDeltas[idx]
			partIndexes := make([]int, 0, len(partsByIndex))
			for partIdx, text := range partsByIndex {
				if text == "" {
					continue
				}
				partIndexes = append(partIndexes, partIdx)
			}
			sort.Ints(partIndexes)
			if len(partIndexes) == 0 {
				continue
			}

			content := make([]any, 0, len(partIndexes))
			for _, partIdx := range partIndexes {
				content = append(content, map[string]any{
					"type": "output_text",
					"text": partsByIndex[partIdx],
				})
			}
			output = append(output, map[string]any{
				"id":      fmt.Sprintf("msg_synth_%d", idx),
				"type":    "message",
				"status":  "completed",
				"role":    "assistant",
				"content": content,
			})
		}
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("no response.completed event found in SSE stream")
	}

	return json.Marshal(map[string]any{
		"type": "response.completed",
		"response": map[string]any{
			"id":     "resp_synthesized_from_sse",
			"output": output,
			"usage":  map[string]any{},
		},
	})
}

func mergeCompletedOutputItems(completed []byte, outputItems map[int]map[string]any) ([]byte, error) {
	if len(outputItems) == 0 {
		return completed, nil
	}

	var wrapper map[string]any
	if err := json.Unmarshal(completed, &wrapper); err != nil {
		return nil, err
	}

	response, ok := wrapper["response"].(map[string]any)
	if !ok {
		return completed, nil
	}
	output, ok := response["output"].([]any)
	if ok && len(output) > 0 {
		return completed, nil
	}

	indexes := make([]int, 0, len(outputItems))
	for idx := range outputItems {
		indexes = append(indexes, idx)
	}
	sort.Ints(indexes)

	merged := make([]any, 0, len(indexes))
	for _, idx := range indexes {
		merged = append(merged, outputItems[idx])
	}
	response["output"] = merged

	return json.Marshal(wrapper)
}

func newSSETerminalError(event, code, message, fallback string) error {
	parts := []string{event}
	if strings.TrimSpace(code) != "" {
		parts = append(parts, strings.TrimSpace(code))
	}
	if strings.TrimSpace(message) != "" {
		parts = append(parts, strings.TrimSpace(message))
	}
	if len(parts) > 1 {
		return fmt.Errorf("%s", strings.Join(parts, ": "))
	}
	if strings.TrimSpace(fallback) != "" {
		return fmt.Errorf("%s: %s", event, strings.TrimSpace(fallback))
	}
	return fmt.Errorf("%s", event)
}

func putReasoningSummaryText(dst map[string]map[int]string, itemID string, summaryIndex int, text string) {
	if dst[itemID] == nil {
		dst[itemID] = map[int]string{}
	}
	dst[itemID][summaryIndex] += text
}

func mergeReasoningSummaries(completed []byte, summaries, deltas map[string]map[int]string) ([]byte, error) {
	if len(summaries) == 0 && len(deltas) == 0 {
		return completed, nil
	}

	var wrapper map[string]any
	if err := json.Unmarshal(completed, &wrapper); err != nil {
		return nil, err
	}

	response, ok := wrapper["response"].(map[string]any)
	if !ok {
		return completed, nil
	}
	output, ok := response["output"].([]any)
	if !ok {
		return completed, nil
	}

	changed := false
	for _, item := range output {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if msg["type"] != "reasoning" {
			continue
		}
		if existing, ok := msg["summary"].([]any); ok && len(existing) > 0 {
			continue
		}

		itemID, _ := msg["id"].(string)
		textByIndex := summaries[itemID]
		if len(textByIndex) == 0 {
			textByIndex = deltas[itemID]
		}
		if len(textByIndex) == 0 {
			continue
		}

		indices := make([]int, 0, len(textByIndex))
		for idx, text := range textByIndex {
			if strings.TrimSpace(text) == "" {
				continue
			}
			indices = append(indices, idx)
		}
		if len(indices) == 0 {
			continue
		}
		sort.Ints(indices)

		parts := make([]any, 0, len(indices))
		for _, idx := range indices {
			parts = append(parts, map[string]any{
				"type": "summary_text",
				"text": textByIndex[idx],
			})
		}
		msg["summary"] = parts
		changed = true
	}

	if !changed {
		return completed, nil
	}
	return json.Marshal(wrapper)
}

func (p *OpenAIResponsesProtocol) ClassifyError(statusCode int, body []byte) error {
	switch {
	case statusCode == 401 || statusCode == 403:
		return ErrAuth
	default:
		var re responsesError
		if json.Unmarshal(body, &re) == nil {
			code := strings.ToLower(re.Error.Code)
			msg := strings.ToLower(re.Error.Message)

			// Context overflow detection
			if code == "context_length_exceeded" ||
				strings.Contains(msg, "maximum context length") ||
				strings.Contains(msg, "too many tokens") ||
				(strings.Contains(msg, "token") && strings.Contains(msg, "exceed")) {
				return ErrContextOverflow
			}
		}
		return fmt.Errorf("openai responses API error (HTTP %d): %s", statusCode, string(body))
	}
}

// ---------------------------------------------------------------------------
// Tool conversion: ToolSchema → Responses API format (flat)
// ---------------------------------------------------------------------------

func convertResponsesTools(tools []ToolSchema) []responsesTool {
	if len(tools) == 0 {
		return nil
	}
	out := make([]responsesTool, len(tools))
	for i, t := range tools {
		params := t.Parameters
		if params == nil {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out[i] = responsesTool{
			Type:        "function",
			Name:        t.Name,
			Description: t.Description,
			Parameters:  params,
			// Responses API defaults to strict=true, we follow that convention
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Message conversion: tape → Responses API
// ---------------------------------------------------------------------------

func convertResponsesInput(msgs []tape.Message) ([]responsesItem, string) {
	var input []responsesItem
	var instructionsParts []string
	var pendingImages []responsesItem

	flushImages := func() {
		input = append(input, pendingImages...)
		pendingImages = nil
	}

	for _, m := range msgs {
		if m.Role != tape.RoleToolResult {
			flushImages()
		}

		switch m.Role {
		case tape.RoleSystem:
			// System messages go to top-level instructions field
			if strings.TrimSpace(m.Content) != "" {
				instructionsParts = append(instructionsParts, m.Content)
			}

		case tape.RoleUser:
			content := m.Content
			if strings.TrimSpace(content) == "" {
				content = "(no user input)"
			}
			input = append(input, responsesItem{
				Role: "user",
				Content: []responsesContent{
					{Type: "input_text", Text: content},
				},
			})

		case tape.RoleAssistant:
			// Add text content as a message item
			if strings.TrimSpace(m.Content) != "" {
				input = append(input, responsesItem{
					Role: "assistant",
					Content: []responsesContent{
						{Type: "output_text", Text: m.Content},
					},
				})
			}
			// Add each tool call as a separate function_call item
			for _, tc := range m.ToolCalls {
				input = append(input, responsesItem{
					Type:      "function_call",
					CallID:    tc.ID,
					Name:      tc.Name,
					Arguments: marshalToolArguments(tc.Arguments),
				})
			}

		case tape.RoleToolResult:
			modelContent := tape.ToolResultModelContent(m.Content, m.StructuredContent)
			input = append(input, responsesItem{
				Type:   "function_call_output",
				CallID: m.ToolID,
				Output: modelContent,
			})
			if m.Image != nil {
				dataURI := "data:" + m.Image.MIMEType + ";base64," + m.Image.Data
				carrierText := modelContent
				if strings.TrimSpace(carrierText) == "" {
					carrierText = " "
				}
				pendingImages = append(pendingImages, responsesItem{
					Role: "user",
					Content: []responsesContent{
						{Type: "input_text", Text: carrierText},
						{Type: "input_image", ImageURL: dataURI},
					},
				})
			}
		}
	}

	flushImages()

	return input, strings.Join(instructionsParts, "\n\n")
}

// ---------------------------------------------------------------------------
// Response parsing: Responses API → tape
// ---------------------------------------------------------------------------

func parseResponsesOutput(resp responsesResponse) tape.Message {
	var contentParts []string
	var toolCalls []tape.ToolCall
	var reasoningItems []tape.ReasoningItem

	for _, item := range resp.Output {
		switch item.Type {
		case "reasoning":
			// Capture reasoning items for the tape
			ri := tape.ReasoningItem{ID: item.ID}
			for _, s := range item.Summary {
				if s.Type == "summary_text" && s.Text != "" {
					ri.Summary = append(ri.Summary, s.Text)
				}
			}
			if len(ri.Summary) > 0 || ri.ID != "" {
				reasoningItems = append(reasoningItems, ri)
			}

		case "message":
			if item.Role == "assistant" {
				for _, c := range item.Content {
					if c.Type == "output_text" && c.Text != "" {
						contentParts = append(contentParts, c.Text)
					}
				}
			}

		case "function_call":
			toolCalls = append(toolCalls, tape.ToolCall{
				ID:        item.CallID,
				Name:      item.Name,
				Arguments: decodeToolArguments(item.Arguments),
			})
		}
	}

	return tape.Message{
		Role:           tape.RoleAssistant,
		Content:        strings.Join(contentParts, ""),
		ReasoningItems: reasoningItems,
		ToolCalls:      toolCalls,
		Timestamp:      time.Now().UnixMilli(),
	}
}
