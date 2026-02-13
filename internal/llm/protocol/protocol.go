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

// Protocol defines how to encode/decode messages for a specific API format.
type Protocol interface {
	// EncodeRequest converts tape messages to provider-specific request body.
	EncodeRequest(messages []tape.Message, tools []ToolSchema, model string, maxTokens int) ([]byte, error)

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
// Only "openai" and "anthropic" are supported.
func For(apiType, model string) (Protocol, error) {
	switch apiType {
	case "anthropic":
		return &AnthropicProtocol{}, nil
	case "openai":
		return &OpenAIProtocol{}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", apiType)
	}
}
