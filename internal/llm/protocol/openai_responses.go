package protocol

import (
	"encoding/json"
	"fmt"
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
	Input        []responsesItem     `json:"input,omitempty"`
	Tools        []responsesTool     `json:"tools,omitempty"`
	Store        bool                `json:"store"`               // Must be false for Codex models
	Stream       bool                `json:"stream"`              // Must be true for Codex models
	Reasoning    *responsesReasoning `json:"reasoning,omitempty"` // Reasoning configuration
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
	Parameters  map[string]any `json:"parameters,omitempty"`
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

type responsesContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
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

	req := responsesRequest{
		Model:        model,
		Instructions: instructions,
		Input:        input,
		Tools:        convertResponsesTools(tools),
		Store:        false, // Required for Codex models
		Stream:       true,  // Required for Codex models
		Reasoning: &responsesReasoning{
			Summary: "detailed", // Request detailed reasoning summaries
		},
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
		}
	}

	return json.Marshal(req)
}

func (p *OpenAIResponsesProtocol) DecodeResponse(body []byte) (tape.Message, Usage, error) {
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

	for _, m := range msgs {
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
				argsJSON, _ := json.Marshal(tc.Arguments)
				input = append(input, responsesItem{
					Type:      "function_call",
					CallID:    tc.ID,
					Name:      tc.Name,
					Arguments: string(argsJSON),
				})
			}

		case tape.RoleToolResult:
			// TODO(vision): The Responses API function_call_output type does not
			// currently support image content. Image data is silently dropped here.
			// When OpenAI adds image support to function_call_output, attach the
			// image as an additional content part.
			input = append(input, responsesItem{
				Type:   "function_call_output",
				CallID: m.ToolID,
				Output: m.Content,
			})
		}
	}

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
			var args map[string]any
			_ = json.Unmarshal([]byte(item.Arguments), &args)
			toolCalls = append(toolCalls, tape.ToolCall{
				ID:        item.CallID,
				Name:      item.Name,
				Arguments: args,
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
