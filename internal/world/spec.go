package world

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const DefaultSpecFilename = "world.json"

var ErrUnknownID = errors.New("unknown id")
var ErrBudgetExhausted = errors.New("budget exhausted")

var embeddedSpecBase64 string

// Spec is the minimal hidden-world mapping for the generic query tool.
// The filesystem remains the lineage habitat; this file only defines the
// opaque query surface.
//
// When Config is present and Config.Budget > 0, the world operates in
// budgeted mode: each get costs 1 from a shared budget, and reset
// regenerates all item values.
type Spec struct {
	Items  map[string]string `json:"items"`
	Config *SpecConfig       `json:"config,omitempty"`
}

// SpecConfig holds optional budgeted-world configuration.
type SpecConfig struct {
	Budget        int    `json:"budget"`
	StateDir      string `json:"state_dir"`
	Cells         int    `json:"cells"`
	AgentGetLimit int    `json:"agent_get_limit,omitempty"`
	ResetQuorum   int    `json:"reset_quorum,omitempty"`
}

// Budgeted reports whether this spec operates under a shared budget.
func (s *Spec) Budgeted() bool {
	return s.Config != nil && s.Config.Budget > 0
}

func Load(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read spec %q: %w", path, err)
	}

	return parseSpec(data, path)
}

func parseSpec(data []byte, source string) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parse spec %q: %w", source, err)
	}

	if len(spec.Items) == 0 {
		return nil, fmt.Errorf("parse spec %q: items must not be empty", source)
	}

	return &spec, nil
}

func DefaultSpecPath() string {
	// Use QUINE_DATA_DIR-relative path to prevent agents from discovering
	// the spec location through WORLD_SPEC environment variable.
	if dataDir := os.Getenv("QUINE_DATA_DIR"); dataDir != "" {
		return dataDir + "/world/world.json"
	}
	return DefaultSpecFilename
}

func LoadDefault() (*Spec, string, error) {
	// Load from QUINE_DATA_DIR-relative path only; WORLD_SPEC is no longer
	// supported to reduce the attack surface for per-agent limit bypass.
	if dataDir := os.Getenv("QUINE_DATA_DIR"); dataDir != "" {
		path := dataDir + "/world/world.json"
		spec, err := Load(path)
		if err != nil {
			return nil, path, err
		}
		return spec, path, nil
	}

	if embeddedSpecBase64 != "" {
		data, err := base64.StdEncoding.DecodeString(embeddedSpecBase64)
		if err != nil {
			return nil, "embedded", fmt.Errorf("decode embedded spec: %w", err)
		}
		spec, err := parseSpec(data, "embedded")
		if err != nil {
			return nil, "embedded", err
		}
		return spec, "embedded", nil
	}

	path := DefaultSpecFilename
	spec, err := Load(path)
	if err != nil {
		return nil, path, err
	}
	return spec, path, nil
}

func (s *Spec) Payload(id string) (string, error) {
	payload, ok := s.Items[id]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownID, id)
	}
	return payload, nil
}

// SaveSpec writes the spec to a JSON file.
func SaveSpec(spec *Spec, path string) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal spec: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
