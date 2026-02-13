package llm

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kehao95/quine/cmd/quine/internal/config"
	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

// bedrockProvider implements Provider for AWS Bedrock's InvokeModel API
// using raw HTTP with SigV4 signing (zero external dependencies).
type bedrockProvider struct {
	region        string
	accessKey     string
	secretKey     string
	sessionToken  string
	modelID       string
	endpoint      string
	client        *http.Client
	contextWindow int // from registry, 0 means use default
}

func newBedrockProvider(cfg *config.Config) (*bedrockProvider, error) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}
	if region == "" {
		return nil, fmt.Errorf("AWS_DEFAULT_REGION or AWS_REGION is required for bedrock provider")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID is required for bedrock provider")
	}

	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY is required for bedrock provider")
	}

	sessionToken := os.Getenv("AWS_SESSION_TOKEN")

	endpoint := cfg.APIBase
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com", region)
	}
	endpoint = strings.TrimRight(endpoint, "/")

	return &bedrockProvider{
		region:        region,
		accessKey:     accessKey,
		secretKey:     secretKey,
		sessionToken:  sessionToken,
		modelID:       cfg.APIModelID(),
		endpoint:      endpoint,
		client:        &http.Client{Timeout: 5 * time.Minute},
		contextWindow: cfg.ContextWindow,
	}, nil
}

// ContextWindowSize returns the context window for the configured model.
func (p *bedrockProvider) ContextWindowSize() int {
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
// Bedrock API request type (differs from Anthropic: no "model", uses
// "anthropic_version" instead of header)
// ---------------------------------------------------------------------------

type bedrockRequest struct {
	AnthropicVersion string             `json:"anthropic_version"`
	MaxTokens        int                `json:"max_tokens"`
	System           string             `json:"system,omitempty"`
	Messages         []anthropicMessage `json:"messages"`
	Tools            []anthropicTool    `json:"tools,omitempty"`
}

// ---------------------------------------------------------------------------
// Generate – main entry point
// ---------------------------------------------------------------------------

func (p *bedrockProvider) Generate(messages []tape.Message, tools []ToolSchema) (tape.Message, Usage, error) {
	system, apiMsgs := convertMessages(messages)
	apiTools := convertTools(tools)

	reqBody := bedrockRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        4096,
		System:           system,
		Messages:         apiMsgs,
		Tools:            apiTools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return tape.Message{}, Usage{}, fmt.Errorf("marshalling request: %w", err)
	}

	url := fmt.Sprintf("%s/model/%s/invoke", p.endpoint, p.modelID)

	resp, err := retryWithBackoff(5, func() (*http.Response, error) {
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		if err := p.signV4(req, bodyBytes); err != nil {
			return nil, fmt.Errorf("signing request: %w", err)
		}

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
		return tape.Message{}, Usage{}, classifyBedrockError(resp.StatusCode, respBody)
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
// AWS SigV4 signing (stdlib only)
// ---------------------------------------------------------------------------

// signV4 adds AWS Signature Version 4 headers to the request.
func (p *bedrockProvider) signV4(req *http.Request, payload []byte) error {
	now := time.Now().UTC()
	return p.signV4WithTime(req, payload, now)
}

// signV4WithTime is the testable version that accepts an explicit timestamp.
func (p *bedrockProvider) signV4WithTime(req *http.Request, payload []byte, now time.Time) error {
	const service = "bedrock"

	datestamp := now.Format("20060102")
	amzdate := now.Format("20060102T150405Z")

	// Set required headers before signing.
	req.Header.Set("X-Amz-Date", amzdate)
	if p.sessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", p.sessionToken)
	}

	// The host header must be present for signing. Go's http.Request
	// stores it in req.Host / req.URL.Host, not in req.Header.
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}

	// 1. Canonical request.
	payloadHash := sha256Hex(payload)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)

	// Signed headers: sorted lowercase header names.
	signedHeaderNames := p.signedHeaderNames(req, host)
	signedHeaders := strings.Join(signedHeaderNames, ";")

	canonicalHeaders := p.canonicalHeaders(req, host, signedHeaderNames)

	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI(req),
		canonicalQueryString(req),
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	// 2. String to sign.
	credentialScope := datestamp + "/" + p.region + "/" + service + "/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzdate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	// 3. Signing key.
	signingKey := awsDeriveKey(p.secretKey, datestamp, p.region, service)

	// 4. Signature.
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// 5. Authorization header.
	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		p.accessKey, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)

	return nil
}

// signedHeaderNames returns a sorted list of lowercase header names to sign.
func (p *bedrockProvider) signedHeaderNames(req *http.Request, host string) []string {
	names := make(map[string]bool)
	names["host"] = true
	for k := range req.Header {
		names[strings.ToLower(k)] = true
	}

	sorted := make([]string, 0, len(names))
	for n := range names {
		sorted = append(sorted, n)
	}
	sort.Strings(sorted)

	// Ensure host is not missing.
	_ = host
	return sorted
}

// canonicalHeaders builds the canonical headers string.
func (p *bedrockProvider) canonicalHeaders(req *http.Request, host string, signedHeaderNames []string) string {
	var b strings.Builder
	for _, name := range signedHeaderNames {
		if name == "host" {
			b.WriteString("host:")
			b.WriteString(host)
			b.WriteByte('\n')
		} else {
			values := req.Header.Values(http.CanonicalHeaderKey(name))
			b.WriteString(name)
			b.WriteByte(':')
			b.WriteString(strings.Join(values, ","))
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// canonicalURI returns the URI-encoded path per AWS SigV4 spec.
// AWS requires each path segment to be URI-encoded (encoding all characters
// except unreserved: A-Z a-z 0-9 - _ . ~). This matters for Bedrock model
// IDs that contain colons (e.g. "anthropic.claude-3-5-sonnet-20241022-v2:0").
func canonicalURI(req *http.Request) string {
	path := req.URL.Path
	if path == "" {
		return "/"
	}
	segments := strings.Split(path, "/")
	for i, s := range segments {
		segments[i] = awsURIEncode(s)
	}
	return strings.Join(segments, "/")
}

// awsURIEncode encodes a string per AWS SigV4 URI encoding rules:
// encode every byte except unreserved characters (A-Z, a-z, 0-9, '-', '.', '_', '~').
func awsURIEncode(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '.' || c == '_' || c == '~' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

// canonicalQueryString returns the sorted query string.
func canonicalQueryString(req *http.Request) string {
	// Bedrock InvokeModel has no query params, but handle it generically.
	return req.URL.RawQuery
}

// sha256Hex returns the lowercase hex-encoded SHA-256 of data.
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// hmacSHA256 computes HMAC-SHA256.
func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// awsDeriveKey computes the SigV4 signing key.
func awsDeriveKey(secret, datestamp, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), []byte(datestamp))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}

// ---------------------------------------------------------------------------
// Error classification for Bedrock
// ---------------------------------------------------------------------------

// bedrockAWSError handles the AWS-style error format {"message":"...","__type":"..."}.
type bedrockAWSError struct {
	Message string `json:"message"`
	Type    string `json:"__type"`
}

func classifyBedrockError(statusCode int, body []byte) error {
	switch {
	case statusCode == 401 || statusCode == 403:
		return ErrAuth
	default:
		// Try Anthropic error format first (Bedrock proxies these).
		var ae anthropicError
		if json.Unmarshal(body, &ae) == nil && ae.Error.Message != "" {
			msg := strings.ToLower(ae.Error.Message)
			errType := strings.ToLower(ae.Error.Type)
			if strings.Contains(msg, "context") || strings.Contains(msg, "too many tokens") ||
				strings.Contains(msg, "token") && strings.Contains(msg, "exceed") ||
				errType == "overloaded" {
				return ErrContextOverflow
			}
			return fmt.Errorf("bedrock API error (HTTP %d): %s", statusCode, ae.Error.Message)
		}

		// Try AWS error format.
		var awsErr bedrockAWSError
		if json.Unmarshal(body, &awsErr) == nil && awsErr.Message != "" {
			msg := strings.ToLower(awsErr.Message)
			if strings.Contains(msg, "security token") || strings.Contains(msg, "credential") ||
				strings.Contains(msg, "not authorized") || strings.Contains(msg, "access denied") {
				return ErrAuth
			}
			return fmt.Errorf("bedrock API error (HTTP %d): [%s] %s", statusCode, awsErr.Type, awsErr.Message)
		}

		return fmt.Errorf("bedrock API error (HTTP %d): %s", statusCode, string(body))
	}
}
