package llm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm/protocol"
	"github.com/kehao95/quine/internal/llm/transport"
	"github.com/kehao95/quine/internal/tape"
)

// Re-export types from protocol package for API compatibility
type (
	Usage      = protocol.Usage
	ToolSchema = protocol.ToolSchema
)

// Re-export errors from protocol package
var (
	ErrAuth            = protocol.ErrAuth
	ErrContextOverflow = protocol.ErrContextOverflow
)

// Provider is the interface that all LLM backends must implement.
type Provider interface {
	Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error)
	ContextWindowSize() int
}

// provider implements Provider using composable protocol and transport.
type provider struct {
	proto         protocol.Protocol
	trans         transport.Transport
	endpoint      string
	model         string
	maxTokens     int
	contextWindow int
	client        *http.Client
}

// NewProvider constructs a Provider for the given config.
func NewProvider(cfg *config.Config) (Provider, error) {
	// Get protocol for this API type
	proto, err := protocol.For(cfg.Provider, cfg.APIModelID())
	if err != nil {
		return nil, err
	}

	// Get transport for this API type
	trans, err := transport.For(cfg.Provider, cfg.APIKey, cfg)
	if err != nil {
		return nil, err
	}

	// Build endpoint URL
	endpoint := buildEndpoint(cfg, proto)

	return &provider{
		proto:         proto,
		trans:         trans,
		endpoint:      endpoint,
		model:         cfg.APIModelID(),
		maxTokens:     defaultMaxTokens(cfg.Provider),
		contextWindow: cfg.ContextWindow,
		client:        &http.Client{Timeout: 10 * time.Minute},
	}, nil
}

// buildEndpoint constructs the full API endpoint URL from base + protocol path.
func buildEndpoint(cfg *config.Config, proto protocol.Protocol) string {
	base := cfg.APIBase
	if base == "" {
		base = defaultAPIBase(cfg.Provider)
	}
	base = strings.TrimRight(base, "/")

	path := proto.EndpointPath()

	// For OpenAI-compatible APIs with custom base URLs,
	// the base may already include /v1 — avoid doubling it.
	if strings.HasPrefix(path, "/v1/") && strings.HasSuffix(base, "/v1") {
		path = path[len("/v1"):]
	}

	return base + path
}

func defaultAPIBase(apiType string) string {
	switch apiType {
	case "anthropic":
		return "https://api.anthropic.com"
	case "openai", "openai-responses":
		return "https://api.openai.com"
	default:
		return ""
	}
}

func defaultMaxTokens(apiType string) int {
	switch apiType {
	case "anthropic":
		return 16384
	default:
		return 4096
	}
}

// Generate sends a conversation and available tools to the model.
func (p *provider) Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error) {
	// Encode request using protocol
	body, err := p.proto.EncodeRequest(messages, tools, p.model, p.maxTokens)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("encoding request: %w", err)
	}

	// Execute with retry
	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		req, err := http.NewRequest("POST", p.endpoint, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", p.proto.ContentType())

		if err := p.trans.Sign(req, body); err != nil {
			return nil, fmt.Errorf("signing request: %w", err)
		}

		return p.client.Do(req)
	})
	if err != nil {
		return tape.Message{}, Usage{}, err
	}
	defer resp.Body.Close()

	// Read the full response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("reading response body: %w", err)
	}

	// Handle SSE streaming responses (used by OpenAI Responses API)
	// Some servers (like ChatGPT) return SSE but with Content-Type: application/json
	// so we detect SSE by looking for the "event: " prefix in the body
	if bytes.HasPrefix(respBody, []byte("event: ")) {
		respBody, err = extractSSECompletedResponse(bytes.NewReader(respBody))
		if err != nil {
			return tape.Message{}, Usage{}, fmt.Errorf("parsing SSE stream: %w", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return tape.Message{}, Usage{}, p.proto.ClassifyError(resp.StatusCode, respBody)
	}

	return p.proto.DecodeResponse(respBody)
}

// ContextWindowSize returns the model's context window size in tokens.
func (p *provider) ContextWindowSize() int {
	if p.contextWindow > 0 {
		return p.contextWindow
	}
	// Fallback defaults
	return 128_000
}

// extractSSECompletedResponse parses an SSE stream and extracts the final
// response.completed event's data payload. This is used for the OpenAI
// Responses API which requires streaming.
func extractSSECompletedResponse(r io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	// Increase buffer size for large SSE events
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var currentEvent string
	var lastCompletedData []byte

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			// We want the response.completed event which contains the full response
			if currentEvent == "response.completed" {
				lastCompletedData = []byte(data)
			}
		}
		// Empty line marks end of event - reset for next event
		if line == "" {
			currentEvent = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning SSE stream: %w", err)
	}

	if lastCompletedData == nil {
		return nil, fmt.Errorf("no response.completed event found in SSE stream")
	}

	return lastCompletedData, nil
}
