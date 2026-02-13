package config

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kehao95/quine/cmd/quine/internal/models"
)

// ErrDepthExceeded is returned when QUINE_DEPTH >= QUINE_MAX_DEPTH.
var ErrDepthExceeded = errors.New("max recursion depth exceeded")

// Config holds all runtime configuration for Quine.
// Every field is populated from environment variables by Load().
type Config struct {
	ModelID        string            // QUINE_MODEL_ID (required) - original model ID for child processes
	ActualModelID  string            // Model ID to use in API calls (may differ when provider prefix is used)
	APIKey         string            // Resolved from provider-specific env var (e.g. ANTHROPIC_API_KEY)
	APIBase        string            // QUINE_API_BASE (optional override)
	Provider       string            // Auto-detected from models.dev registry
	MaxDepth       int               // QUINE_MAX_DEPTH (default 5)
	Depth          int               // QUINE_DEPTH (default 0)
	SessionID      string            // QUINE_SESSION_ID (default auto UUID v4)
	ParentSession  string            // QUINE_PARENT_SESSION
	MaxConcurrent  int               // QUINE_MAX_CONCURRENT (default 20)
	ShTimeout      int               // QUINE_SH_TIMEOUT in seconds (default 600)
	OutputTruncate int               // QUINE_OUTPUT_TRUNCATE in bytes (default 20480)
	DataDir        string            // QUINE_DATA_DIR (default ".quine/")
	Shell          string            // QUINE_SHELL (default "/bin/sh")
	MaxTurns       int               // QUINE_MAX_TURNS (default 20, 0 = unlimited)
	MaxReadLines   int               // QUINE_MAX_READ_LINES (default 500)
	ContextWindow  int               // From models.dev registry
	Wisdom         map[string]string // QUINE_WISDOM_* env vars (key without prefix -> value)
	OriginalIntent string            // QUINE_ORIGINAL_INTENT (preserved across exec for mission continuity)
	StdinOffset    int64             // QUINE_STDIN_OFFSET (preserved across exec for stdin position continuity)
}

// APIModelID returns the model ID to use in API calls.
// This may differ from ModelID when a provider prefix is used (e.g., "openrouter/model-name").
func (c *Config) APIModelID() string {
	if c.ActualModelID != "" {
		return c.ActualModelID
	}
	return c.ModelID
}

// Load reads all configuration from environment variables and returns
// a validated Config. It returns an error if required variables are
// missing, if provider cannot be detected, or if depth is exceeded.
func Load() (*Config, error) {
	c := &Config{}

	// --- Model ID (default: claude-sonnet-4-5-20250514) ---
	c.ModelID = os.Getenv("QUINE_MODEL_ID")
	if c.ModelID == "" {
		c.ModelID = "claude-sonnet-4-5-20250929"
	}

	// --- Optional string fields ---
	c.APIBase = os.Getenv("QUINE_API_BASE")
	c.ParentSession = os.Getenv("QUINE_PARENT_SESSION")

	// --- Load models registry and lookup model ---
	registry, err := models.Load()
	if err != nil {
		return nil, fmt.Errorf("loading models registry: %w", err)
	}

	result, err := registry.Lookup(c.ModelID)
	if err != nil {
		// Fall back to legacy detection for unregistered models
		return loadLegacy(c)
	}

	// Set provider from registry
	c.Provider = result.InternalProvider()
	c.ContextWindow = result.ContextWindow

	// Keep original ModelID for child processes (preserves provider prefix like "openrouter/")
	// Use ActualModelID for API calls
	c.ActualModelID = result.ActualModelID

	// Use API base from registry if not overridden by env var
	if c.APIBase == "" && result.APIBase != "" {
		c.APIBase = result.APIBase
	}

	// API key is not required for bedrock (uses AWS credentials)
	if result.APIKey == "" && c.Provider != "bedrock" {
		envVarNames := strings.Join(result.EnvVars, " or ")
		return nil, fmt.Errorf("API key required: set %s", envVarNames)
	}
	c.APIKey = result.APIKey

	// --- Integer fields with defaults ---
	c.MaxDepth, err = envInt("QUINE_MAX_DEPTH", 5)
	if err != nil {
		return nil, err
	}

	c.Depth, err = envInt("QUINE_DEPTH", 0)
	if err != nil {
		return nil, err
	}

	c.MaxConcurrent, err = envInt("QUINE_MAX_CONCURRENT", 20)
	if err != nil {
		return nil, err
	}

	c.ShTimeout, err = envInt("QUINE_SH_TIMEOUT", 600)
	if err != nil {
		return nil, err
	}

	c.OutputTruncate, err = envInt("QUINE_OUTPUT_TRUNCATE", 20480)
	if err != nil {
		return nil, err
	}

	c.MaxTurns, err = envInt("QUINE_MAX_TURNS", 20)
	if err != nil {
		return nil, err
	}

	c.MaxReadLines, err = envInt("QUINE_MAX_READ_LINES", 500)
	if err != nil {
		return nil, err
	}

	// --- Depth check ---
	if c.Depth >= c.MaxDepth {
		return nil, ErrDepthExceeded
	}

	// --- Session ID ---
	c.SessionID = os.Getenv("QUINE_SESSION_ID")
	if c.SessionID == "" {
		c.SessionID, err = uuidV4()
		if err != nil {
			return nil, fmt.Errorf("generating session ID: %w", err)
		}
	}

	// --- Data dir ---
	c.DataDir = os.Getenv("QUINE_DATA_DIR")
	if c.DataDir == "" {
		c.DataDir = ".quine/"
	}

	// --- Shell ---
	c.Shell = os.Getenv("QUINE_SHELL")
	if c.Shell == "" {
		c.Shell = "/bin/sh"
	}

	// --- Wisdom (QUINE_WISDOM_* env vars) ---
	c.Wisdom = loadWisdom()

	// --- Original Intent (preserved across exec for mission continuity) ---
	c.OriginalIntent = os.Getenv("QUINE_ORIGINAL_INTENT")

	// --- Stdin Offset (preserved across exec for stdin position continuity) ---
	c.StdinOffset, _ = envInt64("QUINE_STDIN_OFFSET", 0)

	return c, nil
}

// loadLegacy handles models not in the registry using the old detection method
func loadLegacy(c *Config) (*Config, error) {
	var err error

	// --- Provider (auto-detect if omitted) ---
	c.Provider = os.Getenv("QUINE_PROVIDER")
	if c.Provider == "" {
		c.Provider, err = detectProvider(c.ModelID)
		if err != nil {
			return nil, err
		}
	}

	// --- API key from provider-specific env vars ---
	c.APIKey = resolveAPIKey(c.Provider)
	if c.APIKey == "" && c.Provider != "bedrock" {
		return nil, fmt.Errorf("API key required for provider %s", c.Provider)
	}

	// --- Integer fields with defaults ---
	c.MaxDepth, err = envInt("QUINE_MAX_DEPTH", 5)
	if err != nil {
		return nil, err
	}

	c.Depth, err = envInt("QUINE_DEPTH", 0)
	if err != nil {
		return nil, err
	}

	c.MaxConcurrent, err = envInt("QUINE_MAX_CONCURRENT", 20)
	if err != nil {
		return nil, err
	}

	c.ShTimeout, err = envInt("QUINE_SH_TIMEOUT", 600)
	if err != nil {
		return nil, err
	}

	c.OutputTruncate, err = envInt("QUINE_OUTPUT_TRUNCATE", 20480)
	if err != nil {
		return nil, err
	}

	c.MaxTurns, err = envInt("QUINE_MAX_TURNS", 20)
	if err != nil {
		return nil, err
	}

	c.MaxReadLines, err = envInt("QUINE_MAX_READ_LINES", 500)
	if err != nil {
		return nil, err
	}

	// --- Depth check ---
	if c.Depth >= c.MaxDepth {
		return nil, ErrDepthExceeded
	}

	// --- Session ID ---
	c.SessionID = os.Getenv("QUINE_SESSION_ID")
	if c.SessionID == "" {
		c.SessionID, err = uuidV4()
		if err != nil {
			return nil, fmt.Errorf("generating session ID: %w", err)
		}
	}

	// --- Data dir ---
	c.DataDir = os.Getenv("QUINE_DATA_DIR")
	if c.DataDir == "" {
		c.DataDir = ".quine/"
	}

	// --- Shell ---
	c.Shell = os.Getenv("QUINE_SHELL")
	if c.Shell == "" {
		c.Shell = "/bin/sh"
	}

	// --- Wisdom (QUINE_WISDOM_* env vars) ---
	c.Wisdom = loadWisdom()

	// --- Original Intent (preserved across exec for mission continuity) ---
	c.OriginalIntent = os.Getenv("QUINE_ORIGINAL_INTENT")

	// --- Stdin Offset (preserved across exec for stdin position continuity) ---
	c.StdinOffset, _ = envInt64("QUINE_STDIN_OFFSET", 0)

	return c, nil
}

// resolveAPIKey tries provider-specific environment variables
func resolveAPIKey(provider string) string {
	switch provider {
	case "anthropic":
		return os.Getenv("ANTHROPIC_API_KEY")
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "google":
		if key := os.Getenv("GOOGLE_GENERATIVE_AI_API_KEY"); key != "" {
			return key
		}
		return os.Getenv("GEMINI_API_KEY")
	case "bedrock":
		return "" // Uses AWS credentials
	default:
		return ""
	}
}

// ChildEnv returns a slice of "KEY=VALUE" environment variable strings
// suitable for spawning a child process. The child gets:
//   - QUINE_DEPTH incremented by 1
//   - QUINE_PARENT_SESSION set to the current SessionID
//   - All other config values inherited
//
// Note: QUINE_SESSION_ID is intentionally NOT included. Each child ./quine
// process generates its own unique session ID via config.Load(). This ensures
// that multiple children spawned from a single sh command (e.g. via &
// backgrounding) each get distinct session IDs and write to separate tape files.
func (c *Config) ChildEnv() ([]string, error) {
	env := []string{
		"QUINE_MODEL_ID=" + c.ModelID,
		"QUINE_API_BASE=" + c.APIBase,
		"QUINE_PROVIDER=" + c.Provider,
		"QUINE_MAX_DEPTH=" + strconv.Itoa(c.MaxDepth),
		"QUINE_DEPTH=" + strconv.Itoa(c.Depth+1),
		"QUINE_PARENT_SESSION=" + c.SessionID,
		"QUINE_MAX_CONCURRENT=" + strconv.Itoa(c.MaxConcurrent),
		"QUINE_SH_TIMEOUT=" + strconv.Itoa(c.ShTimeout),
		"QUINE_OUTPUT_TRUNCATE=" + strconv.Itoa(c.OutputTruncate),
		"QUINE_DATA_DIR=" + c.DataDir,
		"QUINE_SHELL=" + c.Shell,
		"QUINE_MAX_TURNS=" + strconv.Itoa(c.MaxTurns),
		"QUINE_MAX_READ_LINES=" + strconv.Itoa(c.MaxReadLines),
	}

	// Pass through provider-specific API key env vars
	switch c.Provider {
	case "anthropic":
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			env = append(env, "ANTHROPIC_API_KEY="+key)
		}
	case "openai":
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			env = append(env, "OPENAI_API_KEY="+key)
		}
	case "google":
		if key := os.Getenv("GOOGLE_GENERATIVE_AI_API_KEY"); key != "" {
			env = append(env, "GOOGLE_GENERATIVE_AI_API_KEY="+key)
		}
		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			env = append(env, "GEMINI_API_KEY="+key)
		}
	case "bedrock":
		// Pass through AWS credentials
		for _, awsVar := range []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION", "AWS_SESSION_TOKEN"} {
			if val := os.Getenv(awsVar); val != "" {
				env = append(env, awsVar+"="+val)
			}
		}
	}

	// Pass through QUINE_WISDOM_* env vars for state transfer across exec boundaries
	for key, value := range c.Wisdom {
		env = append(env, wisdomPrefix+key+"="+value)
	}

	return env, nil
}

// ExecEnv returns a slice of "KEY=VALUE" environment variable strings
// suitable for exec'ing a fresh process (metamorphosis). Unlike ChildEnv:
//   - DEPTH is NOT incremented (fresh context = restart)
//   - PARENT_SESSION tracks lineage to the pre-exec session
//   - ORIGINAL_INTENT is set to preserve the mission
//   - STDIN_OFFSET preserves the position in the input stream
//   - All QUINE_WISDOM_* vars are preserved (learned insights survive)
//
// Note: QUINE_SESSION_ID is intentionally NOT included. The new process
// generates its own unique session ID via config.Load().
func (c *Config) ExecEnv(originalIntent string, stdinOffset int64) ([]string, error) {
	env := []string{
		"QUINE_MODEL_ID=" + c.ModelID,
		"QUINE_API_BASE=" + c.APIBase,
		"QUINE_PROVIDER=" + c.Provider,
		"QUINE_MAX_DEPTH=" + strconv.Itoa(c.MaxDepth),
		"QUINE_DEPTH=0", // Reset depth — exec is a "fresh brain", not a deeper recursion
		"QUINE_PARENT_SESSION=" + c.SessionID,
		"QUINE_MAX_CONCURRENT=" + strconv.Itoa(c.MaxConcurrent),
		"QUINE_SH_TIMEOUT=" + strconv.Itoa(c.ShTimeout),
		"QUINE_OUTPUT_TRUNCATE=" + strconv.Itoa(c.OutputTruncate),
		"QUINE_DATA_DIR=" + c.DataDir,
		"QUINE_SHELL=" + c.Shell,
		"QUINE_MAX_TURNS=" + strconv.Itoa(c.MaxTurns),
		"QUINE_MAX_READ_LINES=" + strconv.Itoa(c.MaxReadLines),
		"QUINE_ORIGINAL_INTENT=" + originalIntent,
		"QUINE_STDIN_OFFSET=" + strconv.FormatInt(stdinOffset, 10),
	}

	// Pass through provider-specific API key env vars
	switch c.Provider {
	case "anthropic":
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			env = append(env, "ANTHROPIC_API_KEY="+key)
		}
	case "openai":
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			env = append(env, "OPENAI_API_KEY="+key)
		}
	case "google":
		if key := os.Getenv("GOOGLE_GENERATIVE_AI_API_KEY"); key != "" {
			env = append(env, "GOOGLE_GENERATIVE_AI_API_KEY="+key)
		}
		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			env = append(env, "GEMINI_API_KEY="+key)
		}
	case "bedrock":
		// Pass through AWS credentials
		for _, awsVar := range []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION", "AWS_SESSION_TOKEN"} {
			if val := os.Getenv(awsVar); val != "" {
				env = append(env, awsVar+"="+val)
			}
		}
	}

	// Pass through QUINE_WISDOM_* env vars — learned insights survive exec
	for key, value := range c.Wisdom {
		env = append(env, wisdomPrefix+key+"="+value)
	}

	return env, nil
}

// detectProvider infers the LLM provider from the model ID string.
func detectProvider(model string) (string, error) {
	m := strings.ToLower(model)
	switch {
	// Bedrock model IDs have "anthropic." prefix (e.g. "anthropic.claude-3-5-sonnet-20241022-v2:0")
	// or regional prefix like "us.anthropic.claude-...". Must come BEFORE "claude" check.
	case strings.Contains(m, "anthropic."):
		return "bedrock", nil
	case strings.Contains(m, "claude"):
		return "anthropic", nil
	case strings.Contains(m, "gpt"),
		strings.Contains(m, "o1"),
		strings.Contains(m, "o3"),
		strings.Contains(m, "o4"):
		return "openai", nil
	case strings.Contains(m, "gemini"):
		return "google", nil
	default:
		return "", fmt.Errorf("cannot auto-detect provider from model ID %q", model)
	}
}

// envInt reads an environment variable as int, returning def if unset.
func envInt(key string, def int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parsing %s=%q: %w", key, v, err)
	}
	return n, nil
}

// envInt64 reads an environment variable as int64, returning def if unset.
func envInt64(key string, def int64) (int64, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s=%q: %w", key, v, err)
	}
	return n, nil
}

// uuidV4 generates a random UUID v4 using crypto/rand.
func uuidV4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// wisdomPrefix is the environment variable prefix for wisdom transfer.
const wisdomPrefix = "QUINE_WISDOM_"

// loadWisdom scans all environment variables and collects those starting
// with QUINE_WISDOM_. It returns a map with keys stripped of the prefix.
func loadWisdom() map[string]string {
	wisdom := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, wisdomPrefix) {
			// Split on first "=" to get key=value
			key, value, found := strings.Cut(env, "=")
			if !found {
				continue
			}
			// Strip the prefix from the key
			shortKey := strings.TrimPrefix(key, wisdomPrefix)
			if shortKey != "" && value != "" {
				wisdom[shortKey] = value
			}
		}
	}
	return wisdom
}
