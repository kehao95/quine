package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

const (
	defaultAnthropicBase = "https://api.anthropic.com"
	anthropicVersion     = "2023-06-01"
)

// anthropicProvider implements Provider for Anthropic's Messages API.
type anthropicProvider struct {
	apiKey        string
	apiBase       string
	modelID       string
	client        *http.Client
	contextWindow int // from registry, 0 means use default
}

func newAnthropicProvider(cfg *config.Config) *anthropicProvider {
	base := cfg.APIBase
	if base == "" {
		base = defaultAnthropicBase
	}
	// Strip trailing slash for consistent URL construction.
	base = strings.TrimRight(base, "/")
	return &anthropicProvider{
		apiKey:        cfg.APIKey,
		apiBase:       base,
		modelID:       cfg.APIModelID(),
		client:        &http.Client{Timeout: 10 * time.Minute},
		contextWindow: cfg.ContextWindow,
	}
}

// ContextWindowSize returns the context window for the configured model.
func (p *anthropicProvider) ContextWindowSize() int {
	// Use registry value if available
	if p.contextWindow > 0 {
		return p.contextWindow
	}
	// Fallback to defaults
	m := strings.ToLower(p.modelID)
	switch {
	case strings.Contains(m, "claude-3-5-sonnet"):
		return 200_000
	case strings.Contains(m, "claude-sonnet-4"):
		return 200_000
	case strings.Contains(m, "claude-3-haiku"):
		return 200_000
	default:
		return 200_000
	}
}

// ---------------------------------------------------------------------------
// Anthropic API request/response types
// ---------------------------------------------------------------------------

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string or []contentBlock
}

type contentBlock struct {
	Type      string         `json:"type"`
	Text      *string        `json:"text,omitempty"`
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
	Content   string         `json:"content,omitempty"`
}

func strPtr(s string) *string { return &s }

type anthropicTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type  string         `json:"type"`
		Text  string         `json:"text,omitempty"`
		ID    string         `json:"id,omitempty"`
		Name  string         `json:"name,omitempty"`
		Input map[string]any `json:"input,omitempty"`
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
// Generate – main entry point
// ---------------------------------------------------------------------------

func (p *anthropicProvider) Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error) {
	system, apiMsgs := convertMessages(messages)
	apiTools := convertTools(tools)

	reqBody := anthropicRequest{
		Model:     p.modelID,
		MaxTokens: 16384,
		System:    system,
		Messages:  apiMsgs,
		Tools:     apiTools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("marshalling request: %w", err)
	}

	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		req, err := http.NewRequest("POST", p.apiBase+"/v1/messages", bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("x-api-key", p.apiKey)
		req.Header.Set("anthropic-version", anthropicVersion)
		req.Header.Set("content-type", "application/json")
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
		return tape.Message{}, Usage{}, classifyError(resp.StatusCode, respBody)
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	msg := parseResponse(apiResp)
	usage := Usage{
		InputTokens:  apiResp.Usage.InputTokens,
		OutputTokens: apiResp.Usage.OutputTokens,
	}

	return msg, usage, nil
}

// ---------------------------------------------------------------------------
// Message conversion: tape → Anthropic
// ---------------------------------------------------------------------------

// convertMessages extracts the system prompt and converts the remaining
// tape messages into Anthropic API message format.
func convertMessages(msgs []tape.Message) (string, []anthropicMessage) {
	var system string
	var out []anthropicMessage

	for _, m := range msgs {
		switch m.Role {
		case tape.RoleSystem:
			// Collect all system messages into one string.
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
					ID:    tc.ID,
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
			out = append(out, anthropicMessage{
				Role: "user",
				Content: []contentBlock{
					{
						Type:      "tool_result",
						ToolUseID: m.ToolID,
						Content:   m.Content,
					},
				},
			})
		}
	}

	return system, out
}

// convertTools maps ToolSchema values to Anthropic tool format.
func convertTools(tools []ToolSchema) []anthropicTool {
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

func parseResponse(resp anthropicResponse) tape.Message {
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
		}
	}

	return tape.Message{
		Role:      tape.RoleAssistant,
		Content:   strings.Join(textParts, ""),
		ToolCalls: toolCalls,
		Timestamp: time.Now().UnixMilli(),
	}
}

// ---------------------------------------------------------------------------
// Error classification
// ---------------------------------------------------------------------------

func classifyError(statusCode int, body []byte) error {
	switch {
	case statusCode == 401 || statusCode == 403:
		return ErrAuth
	default:
		// Try to parse Anthropic error body for context overflow hints.
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
// Retry with exponential backoff + jitter
// ---------------------------------------------------------------------------

// retryWithBackoff executes fn with retry logic. The retry behaviour depends
// on the HTTP status code returned:
//
//   - 429  → up to maxRetries retries with exponential backoff + jitter
//   - 5xx  → up to 3 retries
//   - 401/403 → no retry, return ErrAuth
//   - Network error → up to 3 retries
//   - Malformed / unexpected → retry once
//
// On a retryable status the response body is drained and closed before
// sleeping so the connection can be reused.
func retryWithBackoff(maxRetries int, fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := fn()

		if err != nil {
			// Network-level error → retry up to 3 times.
			lastErr = err
			retryLimit := 3
			if attempt >= retryLimit {
				return nil, lastErr
			}
			fmt.Fprintf(io.Discard, "")
			logRetry(attempt+1, retryLimit, err.Error())
			backoffSleep(attempt)
			continue
		}

		// Successful HTTP response.
		switch {
		case resp.StatusCode == http.StatusOK:
			return resp, nil

		case resp.StatusCode == 429:
			// Rate limited → retry up to maxRetries times.
			drainAndClose(resp)
			lastErr = fmt.Errorf("rate limited (429)")
			if attempt >= maxRetries {
				// Return a final response so caller can classify.
				// Re-do the request one last time to get a response body.
				return fn()
			}
			logRetry(attempt+1, maxRetries, "rate limited (429)")
			backoffSleep(attempt)
			continue

		case resp.StatusCode == 401 || resp.StatusCode == 403:
			// Auth error → no retry.
			return resp, nil

		case resp.StatusCode >= 500:
			// Server error → retry up to 3 times.
			drainAndClose(resp)
			lastErr = fmt.Errorf("server error (%d)", resp.StatusCode)
			retryLimit := 3
			if attempt >= retryLimit {
				// Return a fresh response for error classification.
				return fn()
			}
			logRetry(attempt+1, retryLimit, fmt.Sprintf("server error (%d)", resp.StatusCode))
			backoffSleep(attempt)
			continue

		default:
			// Unexpected status → retry once (malformed response).
			if attempt >= 1 {
				return resp, nil
			}
			drainAndClose(resp)
			lastErr = fmt.Errorf("unexpected status %d", resp.StatusCode)
			logRetry(attempt+1, 1, fmt.Sprintf("unexpected status (%d)", resp.StatusCode))
			backoffSleep(attempt)
			continue
		}
	}

	if lastResp != nil {
		return lastResp, nil
	}
	return nil, lastErr
}

func logRetry(attempt, max int, reason string) {
	fmt.Fprintf(stderrWriter(), "quine: LLM retry %d/%d (%s)\n", attempt, max, reason)
}

// stderrWriter returns os.Stderr. Extracted so tests can optionally
// override, but by default logs go to stderr.
func stderrWriter() io.Writer {
	return stderrOut
}

// stderrOut can be overridden in tests or via SetLogOutput.
var stderrOut io.Writer = os.Stderr

// SetLogOutput redirects operational log messages (e.g. retry warnings)
// to w instead of os.Stderr. Pass nil to discard.
func SetLogOutput(w io.Writer) {
	if w == nil {
		stderrOut = io.Discard
		return
	}
	stderrOut = w
}

func backoffSleep(attempt int) {
	base := time.Duration(1<<uint(attempt)) * 500 * time.Millisecond
	jitter := time.Duration(rand.Int63n(int64(base / 2)))
	time.Sleep(base + jitter)
}

func drainAndClose(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
