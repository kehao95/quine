package protocol

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

// OpenAIProtocol implements Protocol for OpenAI's Chat Completions API.
// This is also compatible with OpenRouter, Azure OpenAI, and other OpenAI-compatible APIs.
type OpenAIProtocol struct{}

// ---------------------------------------------------------------------------
// API request/response types
// ---------------------------------------------------------------------------

type openaiRequest struct {
	Model    string          `json:"model"`
	Messages []openaiMessage `json:"messages"`
	Tools    []openaiTool    `json:"tools,omitempty"`
}

type openaiMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	ToolCalls  []openaiToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type openaiToolCall struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Function openaiFunctionCall `json:"function"`
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
	Role      string           `json:"role"`
	Content   string           `json:"content"`
	ToolCalls []openaiToolCall `json:"tool_calls,omitempty"`
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

func (p *OpenAIProtocol) EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int) ([]byte, error) {
	apiMsgs := convertOpenAIMessages(messages)
	apiTools := convertOpenAITools(tools)

	req := openaiRequest{
		Model:    model,
		Messages: apiMsgs,
		Tools:    apiTools,
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

func convertOpenAIMessages(msgs []tape.Message) []openaiMessage {
	var out []openaiMessage

	for _, m := range msgs {
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
			for _, tc := range m.ToolCalls {
				argsJSON, _ := json.Marshal(tc.Arguments)
				msg.ToolCalls = append(msg.ToolCalls, openaiToolCall{
					ID:   tc.ID,
					Type: "function",
					Function: openaiFunctionCall{
						Name:      tc.Name,
						Arguments: string(argsJSON),
					},
				})
			}
			out = append(out, msg)

		case tape.RoleToolResult:
			out = append(out, openaiMessage{
				Role:       "tool",
				Content:    m.Content,
				ToolCallID: m.ToolID,
			})
		}
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
	var toolCalls []tape.ToolCall

	for _, tc := range choice.Message.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		toolCalls = append(toolCalls, tape.ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	return tape.Message{
		Role:      tape.RoleAssistant,
		Content:   choice.Message.Content,
		ToolCalls: toolCalls,
		Timestamp: time.Now().UnixMilli(),
	}
}
