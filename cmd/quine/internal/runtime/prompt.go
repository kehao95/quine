package runtime

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/kehao95/quine/cmd/quine/internal/config"
)

//go:embed system_prompt.md
var systemPromptTemplate string

// BuildSystemPrompt constructs the system prompt from config and the template in §9.
// The mission parameter is appended as a "### Your Mission" section.
func BuildSystemPrompt(cfg *config.Config, mission string) string {
	maxTurns := "unlimited"
	if cfg.MaxTurns > 0 {
		maxTurns = fmt.Sprintf("%d", cfg.MaxTurns)
	}

	// Build wisdom section if there are any wisdom entries
	wisdomSection := formatWisdom(cfg.Wisdom)

	// Build mission section (Harvard Architecture: mission is code, not data)
	missionSection := fmt.Sprintf("\n### Your Mission\n%s\n", mission)

	r := strings.NewReplacer(
		"{DEPTH}", fmt.Sprintf("%d", cfg.Depth),
		"{MAX_DEPTH}", fmt.Sprintf("%d", cfg.MaxDepth),
		"{MAX_TURNS}", maxTurns,
		"{MODEL_ID}", cfg.ModelID,
		"{SESSION_ID}", cfg.SessionID,
		"{SHELL}", cfg.Shell,
		"{WISDOM}", wisdomSection,
		"{MISSION}", missionSection,
	)
	return r.Replace(systemPromptTemplate)
}

// formatWisdom formats the wisdom map as a markdown section.
// Returns an empty string if there are no wisdom entries.
func formatWisdom(wisdom map[string]string) string {
	if len(wisdom) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n### Wisdom (from previous incarnation)\n")
	sb.WriteString("The following state was preserved across an exec boundary:\n")

	// Sort keys for deterministic output
	keys := make([]string, 0, len(wisdom))
	for k := range wisdom {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", key, wisdom[key]))
	}

	return sb.String()
}
