package llm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
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

// StreamKind classifies a StreamEvent surfaced during a streaming generation.
type StreamKind int

const (
	// StreamText is an incremental piece of assistant-visible output text.
	StreamText StreamKind = iota
	// StreamReasoning is an incremental piece of reasoning-summary text.
	StreamReasoning
	// StreamToolCall is an incremental piece of a tool call's arguments.
	StreamToolCall
	// StreamDone is the terminal event carrying the authoritative final cell.
	StreamDone
	// StreamError is the terminal event carrying a generation error.
	StreamError
)

// StreamEvent is one event on a streaming generation channel. Delta events
// (text/reasoning/tool-call) are transient display signal; the terminal event
// is either StreamDone (carrying the authoritative Message + Usage) or
// StreamError (carrying Err). Exactly one terminal event is sent before close.
type StreamEvent struct {
	Kind     StreamKind
	Text     string
	ToolID   string
	ToolName string
	Message  tape.Message // populated on StreamDone
	Usage    Usage        // populated on StreamDone
	Err      error        // populated on StreamError
}

// StreamingProvider is the opt-in streaming superset of Provider. The runtime
// feature-detects it; providers that do not implement it keep the buffered
// path and simply never emit a live preview. The final authoritative cell is
// delivered via a StreamDone event on the returned channel and is derived from
// the same terminal payload the buffered Generate path uses (shared
// final-assembly — the crystallized cell stays reproducible).
type StreamingProvider interface {
	Provider
	// GenerateStream returns a channel of StreamEvents. Delta events arrive as
	// generation proceeds; a single terminal StreamDone/StreamError closes the
	// channel. A non-nil returned error is a synchronous setup/transport failure
	// (no channel); once a channel is returned, all outcomes flow through it.
	GenerateStream(messages []tape.Message, tools []ToolSchema) (<-chan StreamEvent, error)
}

// provider implements Provider using composable protocol and transport.
type provider struct {
	proto                protocol.Protocol
	trans                transport.Transport
	endpoint             string
	model                string
	maxTokens            int
	contextWindow        int
	thinkingBudget       string
	serviceTier          string
	claudeAgentSDKCompat bool
	debugRequestBodyDir  string
	quirks               endpointQuirks
	client               *http.Client
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

	apiBase := resolvedAPIBase(cfg)
	quirks := quirksFor(apiBase, cfg.APIKey)
	if proto.EndpointPath() != "/v1/chat/completions" {
		quirks.thinkingBudgetFallback = nil
	}
	endpoint := buildEndpointWithQuirks(apiBase, proto, quirks)

	return &provider{
		proto:                proto,
		trans:                trans,
		endpoint:             endpoint,
		model:                cfg.APIModelID(),
		maxTokens:            defaultMaxTokens(cfg.Provider),
		contextWindow:        cfg.ContextWindow,
		thinkingBudget:       cfg.ThinkingBudget,
		serviceTier:          cfg.ServiceTier,
		claudeAgentSDKCompat: cfg.Provider == "anthropic" && cfg.APIKey == "claude-oauth",
		debugRequestBodyDir:  cfg.DebugRequestBodyDir,
		quirks:               quirks,
		client:               newHTTPClient(),
	}, nil
}

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          0,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
		DisableKeepAlives:     true,
	}
	return &http.Client{
		Timeout:   10 * time.Minute,
		Transport: transport,
	}
}

// buildEndpoint constructs the full API endpoint URL from base + protocol path.
func buildEndpoint(cfg *config.Config, proto protocol.Protocol) string {
	base := resolvedAPIBase(cfg)
	return buildEndpointWithQuirks(base, proto, quirksFor(base, cfg.APIKey))
}

func resolvedAPIBase(cfg *config.Config) string {
	if cfg.APIBase != "" {
		return cfg.APIBase
	}
	return defaultAPIBase(cfg.Provider)
}

func buildEndpointWithQuirks(base string, proto protocol.Protocol, quirks endpointQuirks) string {
	base = strings.TrimRight(base, "/")

	path := proto.EndpointPath()

	// For OpenAI-compatible APIs with custom base URLs, the base may already
	// include a versioned OpenAI surface. Google's Gemini compatibility base,
	// for example, ends at /v1beta/openai and expects /chat/completions rather
	// than /v1/chat/completions.
	if quirks.trimV1Prefix && strings.HasPrefix(path, "/v1/") {
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
	if transport, ok := p.client.Transport.(*http.Transport); ok {
		defer transport.CloseIdleConnections()
	}

	// Build request options
	opts := protocol.RequestOptions{
		ThinkingBudget:          p.effectiveThinkingBudget(tools),
		ServiceTier:             p.serviceTier,
		NoAssistantPrefill:      p.quirks.noAssistantPrefill,
		StripAssistantReasoning: p.quirks.stripAssistantReasoning,
		ClaudeAgentSDKCompat:    p.claudeAgentSDKCompat,
	}

	// Encode request using protocol
	body, err := p.proto.EncodeRequest(messages, tools, p.model, p.maxTokens, opts)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("encoding request: %w", err)
	}

	const recoverableInferenceRetryLimit = 2
	for attempt := 0; ; attempt++ {
		msg, usage, err := p.generateOnce(body)
		if err == nil {
			return msg, usage, nil
		}
		if !errors.Is(err, ErrRecoverableInference) || attempt >= recoverableInferenceRetryLimit {
			return tape.Message{}, Usage{}, err
		}
		logRetry(attempt+1, recoverableInferenceRetryLimit, err.Error())
		backoffSleep(attempt)
	}
}

func (p *provider) effectiveThinkingBudget(tools []ToolSchema) string {
	budget := p.thinkingBudget
	if budget == "" {
		budget = "high"
	}
	if budget == "off" {
		return budget
	}
	if p.quirks.requiresThinkingBudgetFallback(p.model, len(tools) > 0) {
		return "off"
	}
	return budget
}

func (p *provider) generateOnce(body []byte) (tape.Message, Usage, error) {
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
		readErr := fmt.Errorf("reading response body: %w", err)
		if isTransientInferenceTransportError(readErr) {
			return tape.Message{}, Usage{}, newRecoverableInferenceError(readErr)
		}
		return tape.Message{}, Usage{}, readErr
	}

	if resp.StatusCode != http.StatusOK {
		maybeDumpFailedRequestBody(body, respBody, p.debugRequestBodyDir)
		return tape.Message{}, Usage{}, p.proto.ClassifyError(resp.StatusCode, respBody)
	}

	msg, usage, err := p.proto.DecodeResponse(respBody)
	if err != nil {
		if isTransientInferenceDecodeError(err) {
			return tape.Message{}, Usage{}, newRecoverableInferenceError(err)
		}
		return tape.Message{}, Usage{}, err
	}
	return msg, usage, nil
}

// GenerateStream implements StreamingProvider. When the underlying protocol can
// decode incrementally it streams deltas off the live HTTP body; otherwise it
// falls back to the buffered Generate path, surfacing the result as a single
// terminal event. Either way the final cell flows through a StreamDone event so
// the caller has one consumption shape.
func (p *provider) GenerateStream(messages []tape.Message, tools []ToolSchema) (<-chan StreamEvent, error) {
	streamer, ok := p.proto.(protocol.StreamingProtocol)
	if !ok {
		// No incremental decode for this protocol — fall back to buffered and
		// surface the result as a one-shot stream (no deltas).
		ch := make(chan StreamEvent, 1)
		msg, usage, err := p.Generate(messages, tools)
		if err != nil {
			ch <- StreamEvent{Kind: StreamError, Err: err}
		} else {
			ch <- StreamEvent{Kind: StreamDone, Message: msg, Usage: usage}
		}
		close(ch)
		return ch, nil
	}

	opts := protocol.RequestOptions{
		ThinkingBudget:          p.effectiveThinkingBudget(tools),
		ServiceTier:             p.serviceTier,
		NoAssistantPrefill:      p.quirks.noAssistantPrefill,
		StripAssistantReasoning: p.quirks.stripAssistantReasoning,
		ClaudeAgentSDKCompat:    p.claudeAgentSDKCompat,
	}
	body, err := p.proto.EncodeRequest(messages, tools, p.model, p.maxTokens, opts)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	resp, err := p.openStreamResponse(body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		maybeDumpFailedRequestBody(body, respBody, p.debugRequestBodyDir)
		return nil, p.proto.ClassifyError(resp.StatusCode, respBody)
	}

	ch := make(chan StreamEvent, 32)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		if t, ok := p.client.Transport.(*http.Transport); ok {
			defer t.CloseIdleConnections()
		}
		msg, usage, err := streamer.DecodeStream(resp.Body, func(d protocol.StreamDelta) {
			ch <- streamEventFromDelta(d)
		})
		if err != nil {
			if isTransientInferenceDecodeError(err) {
				err = newRecoverableInferenceError(err)
			}
			ch <- StreamEvent{Kind: StreamError, Err: err}
			return
		}
		ch <- StreamEvent{Kind: StreamDone, Message: msg, Usage: usage}
	}()
	return ch, nil
}

// openStreamResponse performs the HTTP request and returns the live response
// with its body unread, so the protocol can decode it incrementally. Setup and
// connect retries match the buffered path; mid-stream retry is intentionally
// not attempted (the stream has already begun emitting deltas).
func (p *provider) openStreamResponse(body []byte) (*http.Response, error) {
	return retryWithBackoff(5, func() (*http.Response, error) {
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
}

func streamEventFromDelta(d protocol.StreamDelta) StreamEvent {
	switch d.Kind {
	case protocol.StreamDeltaReasoning:
		return StreamEvent{Kind: StreamReasoning, Text: d.Text}
	case protocol.StreamDeltaToolCall:
		return StreamEvent{Kind: StreamToolCall, Text: d.Text, ToolID: d.ToolID, ToolName: d.ToolName}
	default:
		return StreamEvent{Kind: StreamText, Text: d.Text}
	}
}

func maybeDumpFailedRequestBody(body []byte, respBody []byte, dir string) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	stamp := fmt.Sprintf("%d", time.Now().UnixNano())
	_ = os.WriteFile(filepath.Join(dir, stamp+".request.json"), body, 0o644)
	_ = os.WriteFile(filepath.Join(dir, stamp+".response.txt"), respBody, 0o644)
}

// ContextWindowSize returns the model's context window size in tokens.
func (p *provider) ContextWindowSize() int {
	if p.contextWindow > 0 {
		return p.contextWindow
	}
	// Fallback defaults
	return 128_000
}
