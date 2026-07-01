package protocol

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

// OpenAIProtocol implements Protocol for OpenAI's Chat Completions API.
// This is also compatible with OpenRouter, Azure OpenAI, and other OpenAI-compatible APIs.
type OpenAIProtocol struct{}

func isKimiModel(model string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), "kimi")
}

// ---------------------------------------------------------------------------
// API request/response types
// ---------------------------------------------------------------------------

type openaiRequest struct {
	Model           string          `json:"model"`
	Messages        []openaiMessage `json:"messages"`
	Tools           []openaiTool    `json:"tools,omitempty"`
	ReasoningEffort string          `json:"reasoning_effort,omitempty"` // For o1/o3 models
	// Kimi-specific thinking config (embedded in extra_body by JSON marshaling)
	Thinking *openaiThinkingConfig `json:"thinking,omitempty"`
}

// openaiThinkingConfig is Kimi's thinking configuration.
type openaiThinkingConfig struct {
	Type string `json:"type"` // "enabled" or "disabled"
}

type openaiMessage struct {
	Role             string           `json:"role"`
	Content          any              `json:"content,omitempty"` // string or []openaiContentPart
	ReasoningContent string           `json:"reasoning_content,omitempty"`
	ToolCalls        []openaiToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string           `json:"tool_call_id,omitempty"`
}

// openaiContentPart is a single element in a multipart content array.
type openaiContentPart struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *openaiImageURL `json:"image_url,omitempty"`
}

// openaiImageURL carries a data URI for image_url content parts.
type openaiImageURL struct {
	URL string `json:"url"` // "data:<mime>;base64,<data>"
}

type openaiToolCall struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"`
	Function     openaiFunctionCall `json:"function"`
	ExtraContent json.RawMessage    `json:"extra_content,omitempty"`
}

type openaiFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openaiTool struct {
	Type     string         `json:"type"`
	Function openaiFunction `json:"function"`
}

type openaiFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type openaiResponse struct {
	Choices []openaiChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

type openaiChoice struct {
	Message      openaiResponseMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

type openaiResponseMessage struct {
	Role             string           `json:"role"`
	Content          string           `json:"content"`
	ExtraContent     json.RawMessage  `json:"extra_content,omitempty"`
	ReasoningContent string           `json:"reasoning_content,omitempty"`
	ToolCalls        []openaiToolCall `json:"tool_calls,omitempty"`
}

type openaiError struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// ---------------------------------------------------------------------------
// Protocol implementation
// ---------------------------------------------------------------------------

func (p *OpenAIProtocol) ContentType() string {
	return "application/json"
}

func (p *OpenAIProtocol) EndpointPath() string {
	return "/v1/chat/completions"
}

func (p *OpenAIProtocol) EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int, opts RequestOptions) ([]byte, error) {
	apiMsgs := convertOpenAIMessages(messages, opts.NoAssistantPrefill, opts.StripAssistantReasoning)
	apiTools := convertOpenAITools(tools)

	req := openaiRequest{
		Model:    model,
		Messages: apiMsgs,
		Tools:    apiTools,
	}

	// Apply thinking budget if specified.
	//
	// `thinking` is a Kimi-specific extension on OpenAI-style chat completions.
	// Generic OpenAI chat/completions accepts `reasoning_effort` for supported
	// models but rejects the extra `thinking` object.
	if opts.ThinkingBudget != "" {
		kimiModel := isKimiModel(model)
		switch opts.ThinkingBudget {
		case "off":
			if kimiModel {
				req.Thinking = &openaiThinkingConfig{Type: "disabled"}
			}
		case "low", "medium", "high", "xhigh":
			req.ReasoningEffort = opts.ThinkingBudget
			if kimiModel {
				req.Thinking = &openaiThinkingConfig{Type: "enabled"}
			}
		}
	}

	return json.Marshal(req)
}

func (p *OpenAIProtocol) DecodeResponse(body []byte) (tape.Message, Usage, error) {
	var resp openaiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	msg := parseOpenAIResponse(resp)
	usage := Usage{
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
	}

	return msg, usage, nil
}

func (p *OpenAIProtocol) ClassifyError(statusCode int, body []byte) error {
	switch {
	case statusCode == 401 || statusCode == 403:
		return ErrAuth
	default:
		var oe openaiError
		if json.Unmarshal(body, &oe) == nil {
			msg := strings.ToLower(oe.Error.Message)
			code := strings.ToLower(oe.Error.Code)
			if strings.Contains(msg, "context") || strings.Contains(msg, "too many tokens") ||
				strings.Contains(msg, "maximum context length") ||
				(strings.Contains(msg, "token") && strings.Contains(msg, "exceed")) ||
				code == "context_length_exceeded" {
				return ErrContextOverflow
			}
		}
		return fmt.Errorf("openai API error (HTTP %d): %s", statusCode, string(body))
	}
}

// ---------------------------------------------------------------------------
// Message conversion: tape → OpenAI
// ---------------------------------------------------------------------------

func convertOpenAIMessages(msgs []tape.Message, noAssistantPrefill, stripAssistantReasoning bool) []openaiMessage {
	msgs = trimIncompleteOpenAIToolBatches(msgs)
	if noAssistantPrefill {
		for len(msgs) > 0 {
			last := msgs[len(msgs)-1]
			if last.Role != tape.RoleAssistant || len(last.ToolCalls) != 0 {
				break
			}
			msgs = msgs[:len(msgs)-1]
		}
	}

	var out []openaiMessage
	// pendingImages holds invisible user messages with image data that must be
	// appended AFTER the full contiguous batch of tool messages ends.
	// OpenAI requires all tool responses for a tool_calls batch to be
	// contiguous before any user message — inserting a user message mid-batch
	// causes a 400 "tool_call_id without response" error.
	var pendingImages []openaiMessage

	flushImages := func() {
		out = append(out, pendingImages...)
		pendingImages = nil
	}

	for _, m := range msgs {
		// When we leave a tool_result run, flush any deferred image messages.
		if m.Role != tape.RoleToolResult {
			flushImages()
		}

		switch m.Role {
		case tape.RoleSystem:
			out = append(out, openaiMessage{
				Role:    "system",
				Content: m.Content,
			})

		case tape.RoleUser:
			out = append(out, openaiMessage{
				Role:    "user",
				Content: m.Content,
			})

		case tape.RoleAssistant:
			msg := openaiMessage{
				Role:    "assistant",
				Content: m.Content,
			}
			if !stripAssistantReasoning {
				msg.ReasoningContent = m.ReasoningContent
			}
			for _, tc := range m.ToolCalls {
				msg.ToolCalls = append(msg.ToolCalls, openaiToolCall{
					ID:           tc.ID,
					Type:         "function",
					ExtraContent: tc.ExtraContent,
					Function: openaiFunctionCall{
						Name:      tc.Name,
						Arguments: marshalToolArguments(tc.Arguments),
					},
				})
			}
			out = append(out, msg)

		case tape.RoleToolResult:
			modelContent := tape.ToolResultModelContent(m.Content, m.StructuredContent)
			// OpenAI Chat Completions does not support images in role=tool messages.
			// Workaround: send text as a normal tool message (preserving tool_call_id
			// pairing), then defer an invisible role=user message carrying the image
			// as a data: URI. The user message is flushed after the full tool batch.
			out = append(out, openaiMessage{
				Role:       "tool",
				Content:    modelContent,
				ToolCallID: m.ToolID,
			})
			if m.Image != nil {
				dataURI := "data:" + m.Image.MIMEType + ";base64," + m.Image.Data
				carrierText := modelContent
				if strings.TrimSpace(carrierText) == "" {
					carrierText = " "
				}
				parts := []openaiContentPart{
					{
						Type: "text",
						Text: carrierText,
					},
					{
						Type: "image_url",
						ImageURL: &openaiImageURL{
							URL: dataURI,
						},
					},
				}
				pendingImages = append(pendingImages, openaiMessage{
					Role:    "user",
					Content: parts,
				})
			}
		}
	}

	// Flush any remaining deferred images at end of message list.
	flushImages()

	return out
}

func trimIncompleteOpenAIToolBatches(msgs []tape.Message) []tape.Message {
	if len(msgs) == 0 {
		return nil
	}

	out := make([]tape.Message, 0, len(msgs))
	var (
		pending   map[string]struct{}
		batchFrom = -1
	)

	for _, msg := range msgs {
		if pending != nil && msg.Role != tape.RoleToolResult {
			out = out[:batchFrom]
			pending = nil
			batchFrom = -1
		}

		switch {
		case msg.Role == tape.RoleAssistant && len(msg.ToolCalls) > 0:
			pending = make(map[string]struct{}, len(msg.ToolCalls))
			for _, tc := range msg.ToolCalls {
				pending[tc.ID] = struct{}{}
			}
			batchFrom = len(out)
			out = append(out, msg)
		case msg.Role == tape.RoleToolResult:
			out = append(out, msg)
			if pending != nil {
				delete(pending, msg.ToolID)
				if len(pending) == 0 {
					pending = nil
					batchFrom = -1
				}
			}
		default:
			out = append(out, msg)
		}
	}

	if pending != nil && batchFrom >= 0 {
		out = out[:batchFrom]
	}
	return out
}

func convertOpenAITools(tools []ToolSchema) []openaiTool {
	if len(tools) == 0 {
		return nil
	}
	out := make([]openaiTool, len(tools))
	for i, t := range tools {
		params := t.Parameters
		if params == nil {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out[i] = openaiTool{
			Type: "function",
			Function: openaiFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Response parsing: OpenAI → tape
// ---------------------------------------------------------------------------

func parseOpenAIResponse(resp openaiResponse) tape.Message {
	if len(resp.Choices) == 0 {
		return tape.Message{
			Role:      tape.RoleAssistant,
			Timestamp: time.Now().UnixMilli(),
		}
	}

	choice := resp.Choices[0]
	for _, candidate := range resp.Choices {
		if len(candidate.Message.ToolCalls) > 0 {
			choice = candidate
			break
		}
	}
	var toolCalls []tape.ToolCall

	for _, tc := range choice.Message.ToolCalls {
		toolCalls = append(toolCalls, tape.ToolCall{
			ID:           tc.ID,
			Name:         tc.Function.Name,
			Arguments:    decodeToolArguments(tc.Function.Arguments),
			ExtraContent: append(json.RawMessage(nil), tc.ExtraContent...),
		})
	}

	content, extractedReasoning := normalizeOpenAIThoughtContent(choice.Message.Content, choice.Message.ExtraContent)
	reasoning := strings.TrimSpace(choice.Message.ReasoningContent)
	switch {
	case reasoning == "":
		reasoning = extractedReasoning
	case extractedReasoning != "" && reasoning != extractedReasoning:
		reasoning = reasoning + "\n\n" + extractedReasoning
	}

	return tape.Message{
		Role:             tape.RoleAssistant,
		Content:          content,
		ReasoningContent: reasoning,
		ToolCalls:        toolCalls,
		Timestamp:        time.Now().UnixMilli(),
	}
}

var openAIThoughtBlockPattern = regexp.MustCompile(`(?s)<thought>(.*?)</thought>`)

func normalizeOpenAIThoughtContent(content string, extra json.RawMessage) (string, string) {
	if !shouldExtractOpenAIThoughtContent(content, extra) {
		return content, ""
	}

	matches := openAIThoughtBlockPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return strings.TrimSpace(content), ""
	}

	reasoningParts := make([]string, 0, len(matches))
	for _, match := range matches {
		part := strings.TrimSpace(match[1])
		if part == "" {
			continue
		}
		reasoningParts = append(reasoningParts, part)
	}

	visible := strings.TrimSpace(openAIThoughtBlockPattern.ReplaceAllString(content, ""))
	reasoning := strings.Join(reasoningParts, "\n\n")
	return visible, reasoning
}

func shouldExtractOpenAIThoughtContent(content string, extra json.RawMessage) bool {
	if hasGoogleThoughtMarker(extra) {
		return true
	}
	trimmed := strings.TrimSpace(content)
	return strings.HasPrefix(trimmed, "<thought>") && strings.Contains(trimmed, "</thought>")
}

func hasGoogleThoughtMarker(extra json.RawMessage) bool {
	if len(extra) == 0 {
		return false
	}
	var decoded struct {
		Google struct {
			Thought bool `json:"thought"`
		} `json:"google"`
	}
	if err := json.Unmarshal(extra, &decoded); err != nil {
		return false
	}
	return decoded.Google.Thought
}
