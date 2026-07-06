package config

// staged.go is the validation layer of the staged-config write path: the
// agent stages capability overrides for its NEXT incarnation in the raw SST
// <agentRoot>/config/next.env; the exec boundary validates them against the
// RUNNING binary's compiled registry and merges them over the ExecEnv()
// serialization immediately before syscall.Exec.
//
// Design authority:
//   Paper/theory/views/runtime-capability/registry-design-brief.md (§ C)
// Work order:
//   Paper/_design/migrations/runtime-capability-registry-execution.md (T3.1)
//
// This file is the single source of validation truth: the exec boundary
// (internal/tools/exec.go) is the non-bypassable enforcement point, and the
// Phase-3 FUSE write gate (public/ctl/config, T3.3) reuses ParseStagedEnv as
// its write-time feedback validator. Parse/validate (ParseStagedEnv) is
// deliberately separable from merge (MergeStagedOverrides) so the gate can
// validate without merging.
//
// Red lines (brief "Preserved invariants"):
//   - Load() NEVER reads next.env: the file authors the next process's envp
//     at the exec boundary only; it is not a config source (invariant 2).
//   - The merge never touches baseEnv(): ChildEnv (fork/spawn) shares
//     baseEnv with ExecEnv, and children must never see staged values —
//     children are new sessions whose agent roots start empty (F-3).
//   - Enforcement is compiled into the body (this file), not prompt
//     convention (invariant 4).
//
// Validation checks legality, not wisdom (brief D4): a well-typed suicidal
// value passes. Range/coupling constraints living in Knob.Notes prose are not
// machine-checked here; a legal-typed value that Load() later rejects kills
// the successor at startup, which is the selection-honesty stance D4 records.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// StagedNextEnvName is the staged-override file name under the agent-root
// config/ directory.
const StagedNextEnvName = "next.env"

// StagedNextEnvPath returns the agent-root staged-override SST path
// (<agentRoot>/config/next.env), or "" when no agent root is configured.
// The path definition is single-sourced here for the exec boundary
// (internal/tools/exec.go) and the successor's consume/clear bootstrap step
// (internal/runtime, T3.2).
func (c *Config) StagedNextEnvPath() string {
	root := strings.TrimSpace(c.AgentRoot())
	if root == "" {
		return ""
	}
	return filepath.Join(root, "config", StagedNextEnvName)
}

// StagedFileError is the whole-file rejection: every violation found in a
// staged-override file, agent-legible (knob name, reason, legal form). A
// staged file with any violation is rejected in full — no partial strip —
// so a retry after fixing the named lines applies exactly what was written.
type StagedFileError struct {
	Violations []string
}

func (e *StagedFileError) Error() string {
	if len(e.Violations) == 1 {
		return e.Violations[0]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d violations:", len(e.Violations))
	for _, v := range e.Violations {
		b.WriteString("\n  - ")
		b.WriteString(v)
	}
	return b.String()
}

// stagedKeyRE is the accepted assignment-key shape. Anything else on an
// assignment line is a loud syntax violation, never silently skipped.
var stagedKeyRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ParseStagedEnv parses AND validates staged-override content against the
// running binary's compiled capability registry. On success it returns the
// override map (env name -> value). On any violation it returns a nil map
// and a *StagedFileError enumerating every violation found (syntax and
// semantic) — the whole file is rejected, no partial strip.
//
// Accepted grammar (unsurprising POSIX-env style; the file is written by
// agents via sh):
//
//	KEY=VALUE     one assignment per line; KEY matches [A-Za-z_][A-Za-z0-9_]*
//	              with no surrounding whitespace; VALUE is everything after
//	              the first '=' verbatim to end of line — no quoting, no
//	              quote stripping, no ${...} expansion, no escapes, no
//	              multi-line values, no `export` prefix
//	# comment     lines whose first non-whitespace character is '#'
//	              (blank lines and whitespace-only lines are also skipped)
//
// A trailing "\r" is tolerated per line (CRLF files). Assigning the same KEY
// twice is rejected: this is a capability-staging transaction, not a shell
// script, so last-wins ambiguity bounces instead of guessing.
//
// Validation per assignment, against the compiled registry:
//   - the env name must be a registry knob (unknown names rejected; retired
//     names get a pointer at their replacement when one exists);
//   - mutability whitelist: ONLY `exec-boundary` knobs are stageable —
//     `substrate-pinned`, `runtime-emitted`, and `operator-only` are all
//     rejected (brief § C as amended 2026-07-03: operator-only covers
//     transport/credentials/paths, and agent-staging those would resurrect
//     the deleted escalate transport-mutation channel, D9);
//   - the value must fit Knob.Type: bool is strictly "0"|"1", int is
//     strconv.Atoi, enum must be one of the registry's legal values, string
//     is free-form. An EMPTY value (`KEY=`) is legal for any type: it
//     replaces the serialized entry with an empty value, which Load()
//     uniformly resolves to the knob's default — the explicit
//     reset-to-default idiom.
func ParseStagedEnv(content []byte) (map[string]string, error) {
	staged := make(map[string]string)
	firstLine := make(map[string]int)
	var violations []string

	for i, raw := range strings.Split(string(content), "\n") {
		lineNo := i + 1
		line := strings.TrimSuffix(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			violations = append(violations,
				fmt.Sprintf("line %d: %q is not a KEY=VALUE assignment (accepted: KEY=VALUE lines, # comments, blank lines)", lineNo, line))
			continue
		}
		if !stagedKeyRE.MatchString(key) {
			violations = append(violations,
				fmt.Sprintf("line %d: key %q is not a legal env name (want [A-Za-z_][A-Za-z0-9_]* with no surrounding whitespace)", lineNo, key))
			continue
		}
		if prev, dup := firstLine[key]; dup {
			violations = append(violations,
				fmt.Sprintf("line %d: duplicate assignment for %s (first assigned at line %d); stage each knob exactly once", lineNo, key, prev))
			continue
		}
		firstLine[key] = lineNo

		if v := validateStagedKnob(key, value); v != "" {
			violations = append(violations, fmt.Sprintf("line %d: %s", lineNo, v))
			continue
		}
		staged[key] = value
	}

	if len(violations) > 0 {
		return nil, &StagedFileError{Violations: violations}
	}
	return staged, nil
}

// validateStagedKnob checks one assignment against the compiled registry.
// Returns "" when legal, otherwise an agent-legible violation.
func validateStagedKnob(env, value string) string {
	knob, ok := KnobByEnv(env)
	if !ok {
		if retired, wasKnob := RetiredKnobByEnv(env); wasKnob {
			if retired.Replacement != "" {
				return fmt.Sprintf("%s: retired knob (%s); use %s instead", env, retired.Decision, retired.Replacement)
			}
			return fmt.Sprintf("%s: retired knob (%s); it no longer exists", env, retired.Decision)
		}
		return fmt.Sprintf("%s: unknown env name — not a capability-registry knob (see config/registry.json)", env)
	}
	if knob.Mutability != MutExecBoundary {
		return fmt.Sprintf("%s: mutability %q — only exec-boundary knobs can be staged in next.env (see config/registry.json)", env, knob.Mutability)
	}
	if value == "" {
		// Explicit reset-to-default: Load() resolves an empty value to the
		// knob's default for every type.
		return ""
	}
	switch knob.Type.Kind {
	case TypeBool:
		if value != "0" && value != "1" {
			return fmt.Sprintf("%s: invalid bool value %q (must be \"0\" or \"1\")", env, value)
		}
	case TypeInt:
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Sprintf("%s: invalid int value %q", env, value)
		}
	case TypeEnum:
		for _, legal := range knob.Type.Enum {
			if value == legal {
				return ""
			}
		}
		return fmt.Sprintf("%s: invalid enum value %q (legal: %s)", env, value, strings.Join(knob.Type.Enum, "|"))
	case TypeString:
		// free-form
	}
	return ""
}

// MergeStagedOverrides overlays validated staged overrides onto an env pair
// list ("KEY=VALUE" strings, normally the ExecEnv() output): an existing
// KEY= entry is replaced in place (staged values win — they are the agent's
// explicit statement about the next incarnation), and staged keys with no
// serialized entry are appended in sorted order (a knob at its compiled
// default has no ExecEnv line). Merge priority thereby lands as: registry
// defaults -> launch envp -> in-process cfg -> staged next.env.
//
// The input MUST be exec-path env only. ChildEnv() output must never pass
// through here: staged overrides apply to self-reentry exec, never to
// fork/spawn children (brief § C, F-3).
func MergeStagedOverrides(env []string, staged map[string]string) []string {
	if len(staged) == 0 {
		return env
	}
	out := make([]string, 0, len(env)+len(staged))
	replaced := make(map[string]bool, len(staged))
	for _, kv := range env {
		key, _, _ := strings.Cut(kv, "=")
		if value, ok := staged[key]; ok {
			out = append(out, envKV(key, value))
			replaced[key] = true
			continue
		}
		out = append(out, kv)
	}
	appended := make([]string, 0, len(staged))
	for key := range staged {
		if !replaced[key] {
			appended = append(appended, key)
		}
	}
	sort.Strings(appended)
	for _, key := range appended {
		out = append(out, envKV(key, staged[key]))
	}
	return out
}

// ReadStagedOverrides reads and validates the staged-override file at path.
// An absent file (or empty path) is a clean no-op: (nil, nil). The file is
// never deleted or modified here — a failed syscall.Exec must find it intact
// so retry is idempotent; the successor archives and clears it at bootstrap
// (T3.2).
func ReadStagedOverrides(path string) (map[string]string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read staged config %s: %w", path, err)
	}
	staged, err := ParseStagedEnv(content)
	if err != nil {
		return nil, fmt.Errorf("staged config %s rejected (whole file; nothing applied; the file is left intact — fix or remove it, then retry): %w", path, err)
	}
	return staged, nil
}
