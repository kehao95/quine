// Package protocol defines wire format conversion between tape messages
// and provider-specific API formats.
package protocol

import (
	"fmt"

	"github.com/kehao95/quine/internal/tape"
)

// Usage reports token consumption for a single LLM call.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// ToolSchema describes a tool that can be offered to the model.
type ToolSchema struct {
	Name        string
	Description string
	Parameters  map[string]any // JSON Schema object
}

// RequestOptions contains optional parameters for API requests.
type RequestOptions struct {
	// ThinkingBudget controls the reasoning effort for thinking models.
	// Values: "off", "low", "medium", "high", "xhigh", or "" (model default).
	ThinkingBudget string

	// ServiceTier controls provider-specific latency/cost tiers. For Codex
	// GPT-5.5 fast mode this is "priority".
	ServiceTier string

	// NoAssistantPrefill strips unsupported trailing assistant messages from
	// the request body for providers that reject assistant-prefill history.
	NoAssistantPrefill bool

	// StripAssistantReasoning omits prior assistant turns' reasoning_content
	// from outgoing requests. Some providers (z.ai/GLM) reject conversations
	// that echo accumulated reasoning back; the reasoning is still parsed from
	// responses and stored in the tape. DeepSeek requires the echo, so this is
	// provider-scoped rather than a global default.
	StripAssistantReasoning bool

	// ClaudeAgentSDKCompat adds Claude Agent SDK attribution blocks for
	// Claude Code OAuth traffic. It is intentionally opt-in so normal
	// Anthropic API-key traffic keeps the plain Messages API shape.
	ClaudeAgentSDKCompat bool
}

// Protocol defines how to encode/decode messages for a specific API format.
type Protocol interface {
	// EncodeRequest converts tape messages to provider-specific request body.
	EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int, opts RequestOptions) ([]byte, error)

	// DecodeResponse parses provider response into tape message + usage.
	DecodeResponse(body []byte) (tape.Message, Usage, error)

	// ClassifyError interprets an error response body.
	ClassifyError(statusCode int, body []byte) error

	// ContentType returns the request Content-Type header value.
	ContentType() string

	// EndpointPath returns the API endpoint path (e.g., "/v1/messages").
	EndpointPath() string
}

// For returns the Protocol implementation for a given API type.
// Supported: "openai", "anthropic", "openai-responses".
func For(apiType, model string) (Protocol, error) {
	switch apiType {
	case "anthropic":
		return &AnthropicProtocol{}, nil
	case "openai":
		return &OpenAIProtocol{}, nil
	case "openai-responses":
		return &OpenAIResponsesProtocol{}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", apiType)
	}
}
