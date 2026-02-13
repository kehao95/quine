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

const defaultGoogleBase = "https://generativelanguage.googleapis.com"

// googleProvider implements Provider for Google's Gemini generateContent API.
type googleProvider struct {
	apiKey        string
	apiBase       string
	modelID       string
	client        *http.Client
	contextWindow int // from registry, 0 means use default
}

func newGoogleProvider(cfg *config.Config) *googleProvider {
	base := cfg.APIBase
	if base == "" {
		base = defaultGoogleBase
	}
	base = strings.TrimRight(base, "/")
	return &googleProvider{
		apiKey:        cfg.APIKey,
		apiBase:       base,
		modelID:       cfg.APIModelID(),
		client:        &http.Client{Timeout: 5 * time.Minute},
		contextWindow: cfg.ContextWindow,
	}
}

// ContextWindowSize returns the context window for the configured Gemini model.
func (p *googleProvider) ContextWindowSize() int {
	// Use registry value if available
	if p.contextWindow > 0 {
		return p.contextWindow
	}
	// Fallback to defaults
	m := strings.ToLower(p.modelID)
	switch {
	case strings.Contains(m, "gemini-1.5-pro"):
		return 2_097_152
	case strings.Contains(m, "gemini-2.0-flash"):
		return 1_048_576
	case strings.Contains(m, "gemini-1.5-flash"):
		return 1_048_576
	default:
		return 1_048_576
	}
}

// ---------------------------------------------------------------------------
// Google Gemini API request/response types
// ---------------------------------------------------------------------------

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	Tools             []geminiToolSet `json:"tools,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text             string              `json:"text,omitempty"`
	FunctionCall     *geminiFunctionCall `json:"functionCall,omitempty"`
	FunctionResponse *geminiFuncResponse `json:"functionResponse,omitempty"`
}

type geminiFunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

type geminiFuncResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type geminiToolSet struct {
	FunctionDeclarations []geminiFunctionDecl `json:"functionDeclarations"`
}

type geminiFunctionDecl struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
	} `json:"usageMetadata"`
}

type geminiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// ---------------------------------------------------------------------------
// Generate – main entry point
// ---------------------------------------------------------------------------

func (p *googleProvider) Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error) {
	systemInst, contents := geminiConvertMessages(messages)
	apiTools := geminiConvertTools(tools)

	reqBody := geminiRequest{
		SystemInstruction: systemInst,
		Contents:          contents,
		Tools:             apiTools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("marshalling request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", p.apiBase, p.modelID, p.apiKey)

	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		req, err := http.NewRequest("POST", endpoint, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
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
		return tape.Message{}, Usage{}, geminiClassifyError(resp.StatusCode, respBody)
	}

	var apiResp geminiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	msg := geminiParseResponse(apiResp)
	usage := Usage{
		InputTokens:  apiResp.UsageMetadata.PromptTokenCount,
		OutputTokens: apiResp.UsageMetadata.CandidatesTokenCount,
	}

	return msg, usage, nil
}

// ---------------------------------------------------------------------------
// Message conversion: tape → Gemini
// ---------------------------------------------------------------------------

// geminiConvertMessages extracts the system instruction and converts the
// remaining tape messages into Gemini API content format.
func geminiConvertMessages(msgs []tape.Message) (*geminiContent, []geminiContent) {
	var systemParts []geminiPart
	var out []geminiContent

	for _, m := range msgs {
		switch m.Role {
		case tape.RoleSystem:
			systemParts = append(systemParts, geminiPart{Text: m.Content})

		case tape.RoleUser:
			out = append(out, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: m.Content}},
			})

		case tape.RoleAssistant:
			var parts []geminiPart
			if m.Content != "" {
				parts = append(parts, geminiPart{Text: m.Content})
			}
			for _, tc := range m.ToolCalls {
				parts = append(parts, geminiPart{
					FunctionCall: &geminiFunctionCall{
						Name: tc.Name,
						Args: tc.Arguments,
					},
				})
			}
			if len(parts) == 0 {
				parts = append(parts, geminiPart{Text: ""})
			}
			out = append(out, geminiContent{
				Role:  "model",
				Parts: parts,
			})

		case tape.RoleToolResult:
			out = append(out, geminiContent{
				Role: "user",
				Parts: []geminiPart{
					{
						FunctionResponse: &geminiFuncResponse{
							Name:     m.ToolID,
							Response: map[string]any{"content": m.Content},
						},
					},
				},
			})
		}
	}

	var systemInst *geminiContent
	if len(systemParts) > 0 {
		systemInst = &geminiContent{Parts: systemParts}
	}

	return systemInst, out
}

// geminiConvertTools maps ToolSchema values to Gemini functionDeclarations format.
func geminiConvertTools(tools []ToolSchema) []geminiToolSet {
	if len(tools) == 0 {
		return nil
	}
	decls := make([]geminiFunctionDecl, len(tools))
	for i, t := range tools {
		decls[i] = geminiFunctionDecl{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		}
	}
	return []geminiToolSet{{FunctionDeclarations: decls}}
}

// ---------------------------------------------------------------------------
// Response parsing: Gemini → tape
// ---------------------------------------------------------------------------

func geminiParseResponse(resp geminiResponse) tape.Message {
	var textParts []string
	var toolCalls []tape.ToolCall

	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.FunctionCall != nil {
				toolCalls = append(toolCalls, tape.ToolCall{
					ID:        fmt.Sprintf("gc_%s_%d", part.FunctionCall.Name, time.Now().UnixNano()),
					Name:      part.FunctionCall.Name,
					Arguments: part.FunctionCall.Args,
				})
			} else if part.Text != "" {
				textParts = append(textParts, part.Text)
			}
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

func geminiClassifyError(statusCode int, body []byte) error {
	switch {
	case statusCode == 401 || statusCode == 403:
		return ErrAuth
	default:
		var ge geminiErrorResponse
		if json.Unmarshal(body, &ge) == nil {
			msg := strings.ToLower(ge.Error.Message)
			if strings.Contains(msg, "context") || strings.Contains(msg, "too many tokens") ||
				strings.Contains(msg, "token") && strings.Contains(msg, "exceed") {
				return ErrContextOverflow
			}
		}
		return fmt.Errorf("google API error (HTTP %d): %s", statusCode, string(body))
	}
}
