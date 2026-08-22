package protocol

import (
	"crypto/sha256"
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
	System    any                `json:"system,omitempty"`
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

type anthropicSystemBlock struct {
	Type         string            `json:"type"`
	Text         string            `json:"text"`
	CacheControl map[string]string `json:"cache_control,omitempty"`
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

func (b contentBlock) MarshalJSON() ([]byte, error) {
	switch b.Type {
	case "text":
		text := ""
		if b.Text != nil {
			text = *b.Text
		}
		return json.Marshal(struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			Type: b.Type,
			Text: text,
		})
	case "tool_use":
		input := b.Input
		if input == nil {
			input = map[string]any{}
		}
		return json.Marshal(struct {
			Type  string         `json:"type"`
			ID    string         `json:"id"`
			Name  string         `json:"name"`
			Input map[string]any `json:"input"`
		}{
			Type:  b.Type,
			ID:    b.ID,
			Name:  b.Name,
			Input: input,
		})
	case "tool_result":
		return json.Marshal(struct {
			Type      string `json:"type"`
			ToolUseID string `json:"tool_use_id"`
			Content   any    `json:"content"`
		}{
			Type:      b.Type,
			ToolUseID: b.ToolUseID,
			Content:   b.Content,
		})
	case "image":
		return json.Marshal(struct {
			Type   string       `json:"type"`
			Source *imageSource `json:"source"`
		}{
			Type:   b.Type,
			Source: b.Source,
		})
	case "thinking":
		return json.Marshal(struct {
			Type      string  `json:"type"`
			Thinking  *string `json:"thinking"`
			Signature *string `json:"signature,omitempty"`
		}{
			Type:      b.Type,
			Thinking:  b.Thinking,
			Signature: b.Signature,
		})
	case "redacted_thinking":
		return json.Marshal(struct {
			Type string  `json:"type"`
			Data *string `json:"data"`
		}{
			Type: b.Type,
			Data: b.Data,
		})
	default:
		type alias contentBlock
		return json.Marshal(alias(b))
	}
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
	if opts.NoAssistantPrefill {
		apiMsgs = stripTrailingAnthropicAssistantPrefill(apiMsgs)
	}
	apiTools := convertAnthropicTools(tools)
	var systemField any = system
	if opts.ClaudeAgentSDKCompat {
		systemField = claudeAgentSDKSystem(system, apiMsgs)
	}

	req := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    systemField,
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
			// Thinking blocks must lead the assistant turn and carry their
			// original signature so Anthropic accepts the replay before this
			// turn's tool_use (interleaved thinking). Emitted only when the
			// response actually produced them, so non-thinking flows are
			// byte-for-byte unchanged.
			for _, tb := range m.ThinkingBlocks {
				switch tb.Type {
				case "thinking":
					// An unsigned thinking block cannot be validly replayed —
					// Anthropic 400s on a signature-less thinking block before a
					// tool_result, and a missing block is the documented-safe
					// state — so drop it rather than emit it without a signature.
					if tb.Signature == "" {
						continue
					}
					blocks = append(blocks, contentBlock{
						Type:      "thinking",
						Thinking:  strPtr(tb.Thinking),
						Signature: strPtr(tb.Signature),
					})
				case "redacted_thinking":
					blocks = append(blocks, contentBlock{Type: "redacted_thinking", Data: strPtr(tb.Data)})
				}
			}
			trimmed := strings.TrimRight(m.Content, " \t\n\r")
			if trimmed != "" {
				blocks = append(blocks, contentBlock{
					Type: "text",
					Text: strPtr(trimmed),
				})
			}
			for _, tc := range m.ToolCalls {
				blocks = append(blocks, contentBlock{
					Type:  "tool_use",
					ID:    sanitizeToolID(tc.ID),
					Name:  tc.Name,
					Input: sanitizeToolArguments(tc.Arguments),
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
			modelContent := tape.ToolResultModelContent(m.Content, m.StructuredContent)
			// Build tool_result content: always text, optionally followed by image block.
			if m.Image != nil {
				// Multipart: text + image content array
				blocks := []contentBlock{
					{Type: "text", Text: strPtr(modelContent)},
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
							Content:   modelContent,
						},
					},
				})
			}
		}
	}

	return system, out
}

func stripTrailingAnthropicAssistantPrefill(messages []anthropicMessage) []anthropicMessage {
	for len(messages) > 0 && isAnthropicAssistantPrefill(messages[len(messages)-1]) {
		messages = messages[:len(messages)-1]
	}
	return messages
}

func isAnthropicAssistantPrefill(msg anthropicMessage) bool {
	if msg.Role != "assistant" {
		return false
	}
	switch content := msg.Content.(type) {
	case string:
		return true
	case []contentBlock:
		for _, block := range content {
			switch block.Type {
			case "thinking", "redacted_thinking", "text":
				// These blocks are all generation content. Without a tool_use
				// block, the turn is an unsupported assistant prefill and may
				// be removed so the conversation ends on a user message.
			default:
				return false
			}
		}
		return true
	default:
		return false
	}
}

const (
	claudeAgentSDKIdentity = "You are a Claude agent, built on Anthropic's Claude Agent SDK."
	claudeCodeVersion      = "2.1.178"
	claudeCodeEntrypoint   = "sdk-cli"
	claudeCCHSalt          = "59cf53e54c78"
)

var claudeCCHPositions = []int{4, 7, 20}

func claudeAgentSDKSystem(system string, messages []anthropicMessage) []anthropicSystemBlock {
	firstUserText := firstAnthropicUserText(messages)
	blocks := []anthropicSystemBlock{{
		Type: "text",
		Text: buildClaudeBillingHeader(firstUserText),
	}, {
		Type: "text",
		Text: claudeAgentSDKIdentity,
		CacheControl: map[string]string{
			"type": "ephemeral",
			"ttl":  "1h",
		},
	}}
	if strings.TrimSpace(system) != "" {
		blocks = append(blocks, anthropicSystemBlock{
			Type: "text",
			Text: system,
			CacheControl: map[string]string{
				"type": "ephemeral",
				"ttl":  "1h",
			},
		})
	}
	return blocks
}

func buildClaudeBillingHeader(firstUserText string) string {
	versionSuffix := claudeVersionSuffix(firstUserText, claudeCodeVersion)
	return fmt.Sprintf(
		"x-anthropic-billing-header: cc_version=%s.%s; cc_entrypoint=%s; cch=%s;",
		claudeCodeVersion,
		versionSuffix,
		claudeCodeEntrypoint,
		claudeCCH(firstUserText),
	)
}

func claudeCCH(messageText string) string {
	sum := sha256.Sum256([]byte(messageText))
	return fmt.Sprintf("%x", sum)[:5]
}

func claudeVersionSuffix(messageText, version string) string {
	var sampled strings.Builder
	for _, index := range claudeCCHPositions {
		if index >= 0 && index < len(messageText) {
			sampled.WriteByte(messageText[index])
		} else {
			sampled.WriteByte('0')
		}
	}
	sum := sha256.Sum256([]byte(claudeCCHSalt + sampled.String() + version))
	return fmt.Sprintf("%x", sum)[:3]
}

func firstAnthropicUserText(messages []anthropicMessage) string {
	for _, msg := range messages {
		if msg.Role != "user" {
			continue
		}
		switch content := msg.Content.(type) {
		case string:
			return content
		case []contentBlock:
			for _, block := range content {
				if block.Type == "text" && block.Text != nil {
					return *block.Text
				}
			}
		}
		return ""
	}
	return ""
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
	var thinkingBlocks []tape.ThinkingBlock

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
		case "thinking":
			// Capture the signed thinking block. Anthropic requires it to be
			// replayed (signature included) before this turn's tool_use on the
			// next request; dropping it breaks interleaved-thinking tool loops.
			thinkingBlocks = append(thinkingBlocks, tape.ThinkingBlock{
				Type:      "thinking",
				Thinking:  block.Thinking,
				Signature: block.Signature,
			})
		case "redacted_thinking":
			thinkingBlocks = append(thinkingBlocks, tape.ThinkingBlock{
				Type: "redacted_thinking",
				Data: block.Data,
			})
		}
	}

	return tape.Message{
		Role:           tape.RoleAssistant,
		Content:        strings.Join(textParts, ""),
		ToolCalls:      toolCalls,
		ThinkingBlocks: thinkingBlocks,
		Timestamp:      time.Now().UnixMilli(),
	}
}
