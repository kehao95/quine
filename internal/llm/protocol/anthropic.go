package protocol

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

const anthropicVersion = "2023-06-01"

// AnthropicProtocol implements Protocol for Anthropic's Messages API.
type AnthropicProtocol struct{}

// ---------------------------------------------------------------------------
// API request/response types
// ---------------------------------------------------------------------------

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string or []contentBlock
}

type contentBlock struct {
	Type      string         `json:"type"`
	Text      *string        `json:"text,omitempty"`
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
	Content   string         `json:"content,omitempty"`
}

func strPtr(s string) *string { return &s }

type anthropicTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type  string         `json:"type"`
		Text  string         `json:"text,omitempty"`
		ID    string         `json:"id,omitempty"`
		Name  string         `json:"name,omitempty"`
		Input map[string]any `json:"input,omitempty"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	StopReason string `json:"stop_reason"`
}

type anthropicError struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// ---------------------------------------------------------------------------
// Protocol implementation
// ---------------------------------------------------------------------------

func (p *AnthropicProtocol) ContentType() string {
	return "application/json"
}

func (p *AnthropicProtocol) EndpointPath() string {
	return "/v1/messages"
}

func (p *AnthropicProtocol) EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int) ([]byte, error) {
	system, apiMsgs := convertAnthropicMessages(messages)
	apiTools := convertAnthropicTools(tools)

	req := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    system,
		Messages:  apiMsgs,
		Tools:     apiTools,
	}

	return json.Marshal(req)
}

func (p *AnthropicProtocol) DecodeResponse(body []byte) (tape.Message, Usage, error) {
	var resp anthropicResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	msg := parseAnthropicResponse(resp)
	usage := Usage{
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
	}

	return msg, usage, nil
}

func (p *AnthropicProtocol) ClassifyError(statusCode int, body []byte) error {
	switch {
	case statusCode == 401 || statusCode == 403:
		return ErrAuth
	default:
		var ae anthropicError
		if json.Unmarshal(body, &ae) == nil {
			msg := strings.ToLower(ae.Error.Message)
			errType := strings.ToLower(ae.Error.Type)
			if strings.Contains(msg, "context") || strings.Contains(msg, "too many tokens") ||
				strings.Contains(msg, "token") && strings.Contains(msg, "exceed") ||
				errType == "overloaded" {
				return ErrContextOverflow
			}
		}
		return fmt.Errorf("anthropic API error (HTTP %d): %s", statusCode, string(body))
	}
}

// ---------------------------------------------------------------------------
// Message conversion: tape → Anthropic
// ---------------------------------------------------------------------------

func convertAnthropicMessages(msgs []tape.Message) (string, []anthropicMessage) {
	var system string
	var out []anthropicMessage

	for _, m := range msgs {
		switch m.Role {
		case tape.RoleSystem:
			if system != "" {
				system += "\n\n"
			}
			system += m.Content

		case tape.RoleUser:
			out = append(out, anthropicMessage{
				Role:    "user",
				Content: m.Content,
			})

		case tape.RoleAssistant:
			var blocks []contentBlock
			if m.Content != "" {
				blocks = append(blocks, contentBlock{
					Type: "text",
					Text: strPtr(strings.TrimRight(m.Content, " \t\n\r")),
				})
			}
			for _, tc := range m.ToolCalls {
				blocks = append(blocks, contentBlock{
					Type:  "tool_use",
					ID:    tc.ID,
					Name:  tc.Name,
					Input: tc.Arguments,
				})
			}
			if len(blocks) == 0 {
				blocks = append(blocks, contentBlock{
					Type: "text",
					Text: strPtr(""),
				})
			}
			out = append(out, anthropicMessage{
				Role:    "assistant",
				Content: blocks,
			})

		case tape.RoleToolResult:
			out = append(out, anthropicMessage{
				Role: "user",
				Content: []contentBlock{
					{
						Type:      "tool_result",
						ToolUseID: m.ToolID,
						Content:   m.Content,
					},
				},
			})
		}
	}

	return system, out
}

func convertAnthropicTools(tools []ToolSchema) []anthropicTool {
	if len(tools) == 0 {
		return nil
	}
	out := make([]anthropicTool, len(tools))
	for i, t := range tools {
		schema := t.Parameters
		if schema == nil {
			schema = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out[i] = anthropicTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: schema,
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Response parsing: Anthropic → tape
// ---------------------------------------------------------------------------

func parseAnthropicResponse(resp anthropicResponse) tape.Message {
	var textParts []string
	var toolCalls []tape.ToolCall

	for _, block := range resp.Content {
		switch block.Type {
		case "text":
			if block.Text != "" {
				textParts = append(textParts, block.Text)
			}
		case "tool_use":
			toolCalls = append(toolCalls, tape.ToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: block.Input,
			})
		}
	}

	return tape.Message{
		Role:      tape.RoleAssistant,
		Content:   strings.Join(textParts, ""),
		ToolCalls: toolCalls,
		Timestamp: time.Now().UnixMilli(),
	}
}
