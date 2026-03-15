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
	Thinking  *anthropicThinking `json:"thinking,omitempty"`
}

// anthropicThinking controls extended thinking for Anthropic models.
// See https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking
type anthropicThinking struct {
	Type         string `json:"type"`                    // "enabled" or "disabled"
	BudgetTokens int    `json:"budget_tokens,omitempty"` // required when type="enabled"
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string or []contentBlock
}

type contentBlock struct {
	Type      string         `json:"type"`
	Text      *string        `json:"text,omitempty"`
	Thinking  *string        `json:"thinking,omitempty"`  // for type="thinking"
	Signature *string        `json:"signature,omitempty"` // for type="thinking"
	Data      *string        `json:"data,omitempty"`      // for type="redacted_thinking"
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
	Content   any            `json:"content,omitempty"` // string or []contentBlock for tool_result
	Source    *imageSource   `json:"source,omitempty"`  // for type="image"
}

// imageSource is the Anthropic base64 image source format.
type imageSource struct {
	Type      string `json:"type"`       // always "base64"
	MediaType string `json:"media_type"` // e.g. "image/png"
	Data      string `json:"data"`       // base64-encoded bytes
}

func strPtr(s string) *string { return &s }

// sanitizeToolID ensures tool ID matches Anthropic's required pattern ^[a-zA-Z0-9_-]+$
// Some providers (e.g., Kimi) generate IDs like "sh:0" which contain invalid characters.
func sanitizeToolID(id string) string {
	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		} else {
			result.WriteRune('_') // Replace invalid chars with underscore
		}
	}
	return result.String()
}

type anthropicTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type      string         `json:"type"`
		Text      string         `json:"text,omitempty"`
		Thinking  string         `json:"thinking,omitempty"`
		Signature string         `json:"signature,omitempty"`
		Data      string         `json:"data,omitempty"` // redacted_thinking
		ID        string         `json:"id,omitempty"`
		Name      string         `json:"name,omitempty"`
		Input     map[string]any `json:"input,omitempty"`
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

func (p *AnthropicProtocol) EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int, opts RequestOptions) ([]byte, error) {
	system, apiMsgs := convertAnthropicMessages(messages)
	apiTools := convertAnthropicTools(tools)

	req := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    system,
		Messages:  apiMsgs,
		Tools:     apiTools,
	}

	// Map ThinkingBudget to Anthropic extended thinking parameter.
	// budget_tokens values: low=5000, medium=15000, high=50000.
	// budget_tokens must be < max_tokens, so we cap it.
	switch opts.ThinkingBudget {
	case "low":
		budget := min(5000, maxTokens-1)
		req.Thinking = &anthropicThinking{Type: "enabled", BudgetTokens: budget}
	case "medium":
		budget := min(15000, maxTokens-1)
		req.Thinking = &anthropicThinking{Type: "enabled", BudgetTokens: budget}
	case "high":
		budget := min(50000, maxTokens-1)
		req.Thinking = &anthropicThinking{Type: "enabled", BudgetTokens: budget}
		// "off" or anything else: omit the thinking field entirely
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
					ID:    sanitizeToolID(tc.ID),
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
			// Build tool_result content: always text, optionally followed by image block.
			if m.Image != nil {
				// Multipart: text + image content array
				blocks := []contentBlock{
					{Type: "text", Text: strPtr(m.Content)},
					{
						Type: "image",
						Source: &imageSource{
							Type:      "base64",
							MediaType: m.Image.MIMEType,
							Data:      m.Image.Data,
						},
					},
				}
				out = append(out, anthropicMessage{
					Role: "user",
					Content: []contentBlock{
						{
							Type:      "tool_result",
							ToolUseID: sanitizeToolID(m.ToolID),
							Content:   blocks,
						},
					},
				})
			} else {
				// Plain text tool result
				out = append(out, anthropicMessage{
					Role: "user",
					Content: []contentBlock{
						{
							Type:      "tool_result",
							ToolUseID: sanitizeToolID(m.ToolID),
							Content:   m.Content,
						},
					},
				})
			}
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
			// "thinking" and "redacted_thinking" blocks are internal reasoning;
			// we skip them here as they are not part of the visible response content.
		}
	}

	return tape.Message{
		Role:      tape.RoleAssistant,
		Content:   strings.Join(textParts, ""),
		ToolCalls: toolCalls,
		Timestamp: time.Now().UnixMilli(),
	}
}
