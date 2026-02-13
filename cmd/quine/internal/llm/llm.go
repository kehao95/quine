package llm

import (
	"errors"
	"fmt"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
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

// Provider is the interface that all LLM backends must implement.
type Provider interface {
	// Generate sends a conversation and available tools to the model,
	// returning the assistant's response message and token usage.
	Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error)

	// ContextWindowSize returns the model's context window size in tokens.
	// If configured from the registry, this returns the registry value;
	// otherwise falls back to provider-specific defaults.
	ContextWindowSize() int
}

// Sentinel errors for callers to match with errors.Is.
var (
	ErrAuth            = errors.New("authentication failed")
	ErrContextOverflow = errors.New("context window exceeded")
)

// NewProvider constructs the appropriate Provider for the given config.
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "anthropic":
		return newAnthropicProvider(cfg), nil
	case "openai":
		return newOpenAIProvider(cfg), nil
	case "google":
		return newGoogleProvider(cfg), nil
	case "bedrock":
		return newBedrockProvider(cfg)
	default:
		return nil, fmt.Errorf("provider %s not yet implemented", cfg.Provider)
	}
}
