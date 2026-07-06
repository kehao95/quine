package runtime

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

// quineEnvTokenPattern extracts QUINE_*-shaped tokens from rendered prompt
// text. It requires at least one [A-Z0-9_] character after the QUINE_
// prefix so the bare wildcard reference "every `QUINE_*` knob" (used once in
// the runtime-surface section to describe config/registry.json itself,
// rather than to name one specific env var) is not picked up as a token.
var quineEnvTokenPattern = regexp.MustCompile(`QUINE_[A-Z0-9_]+`)

// knownNonRegistryRuntimeEntrypoints are QUINE_* names the prompt discloses
// that are deliberately NOT config/registry.go knobs: they are runtime-emitted
// identity facts about the current process (e.g. its own live session root),
// not agent-settable or exec-boundary-stageable capability controls.
// internal/config/docgen.go documents this distinction explicitly:
// "Runtime-emitted entrypoints such as `QUINE_AGENT_ROOT` are owned by
// runtime-surface.md unless they also become user-authored controls or
// selectors." Confirmed by reading registry.go: QUINE_AGENT_ROOT has no
// envnames.go constant and no Registry entry at all — it is never read via
// os.Getenv; cfg.AgentRoot() is a computed path. This is the one documented
// case as of T4.3.
var knownNonRegistryRuntimeEntrypoints = map[string]bool{
	"QUINE_AGENT_ROOT": true,
}

// TestBuildSystemPrompt_QuineEnvNamesMatchRegistry cross-checks every
// QUINE_*-shaped token the assembled system prompt names literally against
// internal/config's compiled Registry (config.KnobByEnv). It exists to catch
// a prompt.go/fragments.go change that hardcodes a made-up, stale, or
// renamed QUINE_* name that is not (or is no longer) a real registry entry —
// a drift that would otherwise only be caught by a human proofreading the
// prompt against the registry by eye.
//
// The cfg below turns on every branch (identified by reading prompt.go in
// full) that names an individual QUINE_* env var literally in rendered text
// but is not already active on testConfig()'s zero-valued fields: provider
// transport disclosure (Provider/APIBase/APIKey/UserAgent), ephemeral-body
// physics, and the QUINE_SH_NETWORK=none shell-network constraint line.
// Everything else that names a QUINE_* var (agent-root/run-id/session-id/
// data-dir identity lines, the AGENTS_MD/AGENTS_SKILLS_ENABLED context-files
// lines, and the FORK_ENABLED=0/SPAWN_ENABLED=0 disabled-tool lines) is
// either unconditional given a visible runtime surface or already reachable
// from testConfig()'s defaults (fork/spawn are disabled there).
func TestBuildSystemPrompt_QuineEnvNamesMatchRegistry(t *testing.T) {
	cfg := testConfig()
	cfg.Provider = "openai-responses"
	cfg.APIBase = "http://127.0.0.1:18080"
	cfg.APIKey = "secret-test-key"
	cfg.UserAgent = "test-agent"
	cfg.EphemeralBody = true
	cfg.ShNetwork = "none"

	prompt := BuildSystemPrompt(cfg, "test mission", true)

	matches := quineEnvTokenPattern.FindAllString(prompt, -1)
	if len(matches) == 0 {
		t.Fatal("expected the rendered prompt to name at least one QUINE_* env var literally; extraction regex or this test's cfg has drifted from prompt.go's actual disclosure branches")
	}

	seen := map[string]bool{}
	for _, m := range matches {
		seen[m] = true
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)

	var unknown []string
	for _, name := range names {
		if knownNonRegistryRuntimeEntrypoints[name] {
			continue
		}
		if _, ok := config.KnobByEnv(name); !ok {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		t.Errorf("prompt names QUINE_* env var(s) absent from internal/config's compiled Registry (and not in knownNonRegistryRuntimeEntrypoints): %s\nAll QUINE_* names found in rendered prompt: %s",
			strings.Join(unknown, ", "), strings.Join(names, ", "))
	}

	// Guard the allowlist against going stale: if QUINE_AGENT_ROOT ever
	// becomes a real registry knob, its allowlist entry should be removed so
	// this test still exercises the "is it actually registered" check for
	// that name instead of silently short-circuiting on it forever.
	for name := range knownNonRegistryRuntimeEntrypoints {
		if !seen[name] {
			continue
		}
		if _, ok := config.KnobByEnv(name); ok {
			t.Errorf("knownNonRegistryRuntimeEntrypoints[%q] is stale: it is now a real registry knob; remove it from the allowlist", name)
		}
	}
}
