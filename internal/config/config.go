package config

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ErrDepthExceeded is returned when QUINE_DEPTH >= QUINE_MAX_DEPTH.
var ErrDepthExceeded = errors.New("max recursion depth exceeded")

// Config holds all runtime configuration for Quine.
// Every field is populated from environment variables by Load().
type Config struct {
	ModelID        string            // QUINE_MODEL_ID (required)
	APIKey         string            // QUINE_API_KEY (required)
	APIBase        string            // QUINE_API_BASE (required)
	Provider       string            // QUINE_API_TYPE (required): "openai" or "anthropic"
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
	ContextWindow  int               // QUINE_CONTEXT_WINDOW (default 128000)
	Wisdom         map[string]string // QUINE_WISDOM_* env vars (key without prefix -> value)
	OriginalIntent string            // QUINE_ORIGINAL_INTENT (preserved across exec for mission continuity)
	StdinOffset    int64             // QUINE_STDIN_OFFSET (preserved across exec for stdin position continuity)
}

// APIModelID returns the model ID to use in API calls.
func (c *Config) APIModelID() string {
	return c.ModelID
}

// Load reads all configuration from environment variables and returns
// a validated Config. It returns an error if required variables are
// missing or if depth is exceeded.
//
// Four variables are required:
//   - QUINE_MODEL_ID:   Model name (e.g. "claude-sonnet-4-5-20250929", "gpt-4o", "kimi-k2.5")
//   - QUINE_API_TYPE:   Wire protocol: "openai" or "anthropic"
//   - QUINE_API_BASE:   API base URL (e.g. "https://api.anthropic.com", "https://api.openai.com")
//   - QUINE_API_KEY:    API key
func Load() (*Config, error) {
	c := &Config{}

	// --- 4 required fields ---
	c.ModelID = os.Getenv("QUINE_MODEL_ID")
	if c.ModelID == "" {
		return nil, fmt.Errorf("QUINE_MODEL_ID is required")
	}

	c.Provider = os.Getenv("QUINE_API_TYPE")
	if c.Provider == "" {
		return nil, fmt.Errorf("QUINE_API_TYPE is required (\"openai\" or \"anthropic\")")
	}
	if c.Provider != "openai" && c.Provider != "anthropic" {
		return nil, fmt.Errorf("unsupported QUINE_API_TYPE=%q: must be \"openai\" or \"anthropic\"", c.Provider)
	}

	c.APIBase = os.Getenv("QUINE_API_BASE")
	if c.APIBase == "" {
		return nil, fmt.Errorf("QUINE_API_BASE is required")
	}

	c.APIKey = os.Getenv("QUINE_API_KEY")
	if c.APIKey == "" {
		return nil, fmt.Errorf("QUINE_API_KEY is required")
	}

	// --- Optional string fields ---
	c.ParentSession = os.Getenv("QUINE_PARENT_SESSION")

	// --- Integer fields with defaults ---
	var err error

	c.ContextWindow, err = envInt("QUINE_CONTEXT_WINDOW", 128_000)
	if err != nil {
		return nil, err
	}

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

// baseEnv returns the common environment variable slice shared by
// ChildEnv and ExecEnv. depth and parentSession are parameterized
// since they differ between the two callers.
func (c *Config) baseEnv(depth int, parentSession string) []string {
	env := []string{
		"QUINE_MODEL_ID=" + c.ModelID,
		"QUINE_API_TYPE=" + c.Provider,
		"QUINE_API_BASE=" + c.APIBase,
		"QUINE_API_KEY=" + c.APIKey,
		"QUINE_MAX_DEPTH=" + strconv.Itoa(c.MaxDepth),
		"QUINE_DEPTH=" + strconv.Itoa(depth),
		"QUINE_PARENT_SESSION=" + parentSession,
		"QUINE_MAX_CONCURRENT=" + strconv.Itoa(c.MaxConcurrent),
		"QUINE_SH_TIMEOUT=" + strconv.Itoa(c.ShTimeout),
		"QUINE_OUTPUT_TRUNCATE=" + strconv.Itoa(c.OutputTruncate),
		"QUINE_DATA_DIR=" + c.DataDir,
		"QUINE_SHELL=" + c.Shell,
		"QUINE_MAX_TURNS=" + strconv.Itoa(c.MaxTurns),
		"QUINE_MAX_READ_LINES=" + strconv.Itoa(c.MaxReadLines),
		"QUINE_CONTEXT_WINDOW=" + strconv.Itoa(c.ContextWindow),
	}

	// Pass through QUINE_WISDOM_* env vars for state transfer across exec boundaries
	for key, value := range c.Wisdom {
		env = append(env, wisdomPrefix+key+"="+value)
	}

	return env
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
	return c.baseEnv(c.Depth+1, c.SessionID), nil
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
	env := c.baseEnv(0, c.SessionID)
	env = append(env,
		"QUINE_ORIGINAL_INTENT="+originalIntent,
		"QUINE_STDIN_OFFSET="+strconv.FormatInt(stdinOffset, 10),
	)
	return env, nil
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
