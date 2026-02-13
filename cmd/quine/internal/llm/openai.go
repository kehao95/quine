package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

const defaultOpenAIBase = "https://api.openai.com"

// openaiProvider implements Provider for OpenAI's Chat Completions API.
type openaiProvider struct {
	apiKey        string
	apiBase       string
	modelID       string
	client        *http.Client
	contextWindow int // from registry, 0 means use default
}

func newOpenAIProvider(cfg *config.Config) *openaiProvider {
	base := cfg.APIBase
	if base == "" {
		base = defaultOpenAIBase
	}
	base = strings.TrimRight(base, "/")
	return &openaiProvider{
		apiKey:        cfg.APIKey,
		apiBase:       base,
		modelID:       cfg.APIModelID(),
		client:        &http.Client{Timeout: 5 * time.Minute},
		contextWindow: cfg.ContextWindow,
	}
}

// ContextWindowSize returns the context window for the configured model.
func (p *openaiProvider) ContextWindowSize() int {
	// Use registry value if available
	if p.contextWindow > 0 {
		return p.contextWindow
	}
	// Fallback to defaults
	m := strings.ToLower(p.modelID)
	switch {
	case strings.HasPrefix(m, "o1"),
		strings.HasPrefix(m, "o3"),
		strings.HasPrefix(m, "o4"):
		return 200_000
	case strings.HasPrefix(m, "gpt-4o"):
		return 128_000
	case strings.HasPrefix(m, "gpt-4-turbo"):
		return 128_000
	case strings.HasPrefix(m, "gpt-4"):
		return 8_192
	case strings.HasPrefix(m, "gpt-3.5-turbo"):
		return 16_385
	default:
		return 128_000
	}
}

// ---------------------------------------------------------------------------
// OpenAI API request/response types
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
// Generate – main entry point
// ---------------------------------------------------------------------------

func (p *openaiProvider) Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error) {
	apiMsgs := openaiConvertMessages(messages)
	apiTools := openaiConvertTools(tools)

	reqBody := openaiRequest{
		Model:    p.modelID,
		Messages: apiMsgs,
		Tools:    apiTools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("marshalling request: %w", err)
	}

	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		// Build endpoint URL. If apiBase already ends with /v1, just append /chat/completions
		// Otherwise append /v1/chat/completions (for OpenAI-style base URLs like https://api.openai.com)
		var endpoint string
		if strings.HasSuffix(p.apiBase, "/v1") {
			endpoint = p.apiBase + "/chat/completions"
		} else {
			endpoint = p.apiBase + "/v1/chat/completions"
		}
		req, err := http.NewRequest("POST", endpoint, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
		req.Header.Set("Content-Type", "application/json")
		return p.client.Do(req)
	})
	if err != nil {
		return tape.Message{}, Usage{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return tape.Message{}, Usage{}, openaiClassifyError(resp.StatusCode, respBody)
	}

	var apiResp openaiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	msg := openaiParseResponse(apiResp)
	usage := Usage{
		InputTokens:  apiResp.Usage.PromptTokens,
		OutputTokens: apiResp.Usage.CompletionTokens,
	}

	return msg, usage, nil
}

// ---------------------------------------------------------------------------
// Message conversion: tape → OpenAI
// ---------------------------------------------------------------------------

func openaiConvertMessages(msgs []tape.Message) []openaiMessage {
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

// openaiConvertTools maps ToolSchema values to OpenAI function-calling format.
func openaiConvertTools(tools []ToolSchema) []openaiTool {
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

func openaiParseResponse(resp openaiResponse) tape.Message {
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

// ---------------------------------------------------------------------------
// Error classification
// ---------------------------------------------------------------------------

func openaiClassifyError(statusCode int, body []byte) error {
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
