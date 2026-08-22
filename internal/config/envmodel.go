package config

// envmodel.go is the process-native env boundary model: the compiled tables
// and the single pipeline that build the environment of every process this
// runtime creates.
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//
// There are exactly two levels (brief Layer 1 § 2):
//
//	own env    the birth record. What this process was given at execve.
//	           Immutable in effect, and readable where the OS already
//	           publishes it (/proc/<pid>/environ). Load() reads envp only —
//	           the registry brief's red line 2 is untouched.
//	child env  constructed at each birth this process performs, by one
//	           pipeline shared across every boundary:
//
//	               child = stamps(boundary) ⊕ override ⊕ (environ − mask(boundary))
//
// Precedence runs left to right: stamps always win (they are runtime-owned
// facts about the child, not opinions), an override line beats an inherited
// one, and a masked name never crosses regardless of what the override says.
//
// What this file deliberately does NOT do is synthesize a child env out of
// resolved Config values. A knob the operator never set is ABSENT from a
// child's environ, exactly as it is absent from this process's own. Absence
// means "compiled default", and config/registry.json is the catalog that says
// so. Manufacturing `QUINE_SPAWN_ENABLED=0` lines for defaults nobody chose is
// what taught a founder that constructing a successor was physically
// impossible; that is the failure this model exists to remove, and re-adding a
// synthesizer would re-create it.
//
// Enforcement lives here, compiled into the body — not in the file surfaces
// that project it (brief Layer 1 § 4). Deleting or defacing a projection
// (public/config/registry.json, the ctl/env read-back) changes nothing,
// because nothing in this file reads it. And the pins bind the
// MEDIATED channels only: an agent running `env -i ./quine` through sh is
// generic computation exercising the substrate, not a governance bypass. We do
// not pretend otherwise.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Boundary names a process-construction boundary the runtime owns.
type Boundary string

const (
	// BoundaryShell is an ordinary sh child: generic computation. It carries
	// no lineage marks — a quine launched from a shell is a new founder, not
	// a member of this process tree (brief E4). Depth enforcement never lived
	// in the child's env anyway; it is a parent-side check.
	BoundaryShell Boundary = "sh"
	// BoundaryChild is a managed fork/spawn child: a member of this agent
	// tree, and stamped as one. spawn shares fork's executor and therefore
	// this boundary.
	BoundaryChild Boundary = "fork"
	// BoundarySelf is self-reentry exec: same lineage, new image.
	BoundarySelf Boundary = "exec"
)

// Boundaries lists every boundary, for table-completeness tests.
func Boundaries() []Boundary { return []Boundary{BoundaryShell, BoundaryChild, BoundarySelf} }

// EnvBehavior is what the boundary tables say about one env name.
type EnvBehavior string

const (
	// EnvMasked: never inherited across the boundary, and never settable
	// through the override. The runtime re-stamps the name where it is
	// meaningful, so a child gets a true value instead of a stale one.
	EnvMasked EnvBehavior = "masked"
	// EnvPinned: inherited verbatim; the agent may not change it through any
	// mediated channel. Governance, not physics.
	EnvPinned EnvBehavior = "pinned"
	// EnvFree: inherited; the override may set or unset it.
	EnvFree EnvBehavior = "free"
)

// --- non-registry runtime names ---
//
// The registry models what an operator or agent may AUTHOR. These names are
// emitted by the runtime into processes it builds, so they have no Knob and
// their behavior cannot be derived from Mutability. They are declared here
// instead, and BoundaryBehavior derives their classification from this declaration.
//
// QUINE_CONTEXT_BOOTSTRAP used to live in this bucket; it is a real
// runtime-emitted knob now, which is why the mask below is otherwise pure
// derivation. QUINE_AGENT_ROOT is the one remaining singleton: it is never
// read back by Load() (cfg.AgentRoot() is computed), so promoting it to a knob
// would add a row that describes nothing Load() does.
const (
	// EnvAgentRoot is the live session root the runtime exports to its own
	// children so `$QUINE_AGENT_ROOT/...` works in a shell command. A child
	// must never inherit its parent's: it has a root of its own.
	EnvAgentRoot = "QUINE_AGENT_ROOT"
	// EnvWorkspaceEnabled is emitted by the overlay backend into sh children
	// to describe the workspace physics they are actually running inside.
	EnvWorkspaceEnabled = "QUINE_WORKSPACE_ENABLED"

	// envJobPrefix is sh job-wrapper plumbing (shell, command, session dir,
	// network, interactive, stdin file), stamped per job by the executor.
	envJobPrefix = "QUINE_JOB_"
	// envWorkspaceOverlayPrefix covers the overlay mount coordinates
	// (LOWERDIR, UPPER, WORKDIR, MOUNT_BASE, OVERLAY_EXTRA_OPTS, INIT_ERROR).
	// They describe a specific mount; inheriting a parent's stale coordinates
	// would point a child at the wrong filesystem, so they are re-injected,
	// never inherited.
	envWorkspaceOverlayPrefix = "QUINE_WORKSPACE_OVERLAY_"
	// envTestPrefix is deliberate test-only runtime hooks. Harness-owned:
	// inherited untouched, never governed.
	envTestPrefix = "QUINE_TEST_"
)

// runtimeOwnedOverlayNames are the overlay-physics names that do not share the
// QUINE_WORKSPACE_OVERLAY_ prefix but are just as mount-specific.
var runtimeOwnedOverlayNames = []string{
	EnvWorkspaceEnabled,
	"QUINE_WORKSPACE_LOWERDIR",
	"QUINE_WORKSPACE_UPPER",
	"QUINE_WORKSPACE_WORKDIR",
	"QUINE_WORKSPACE_MOUNT_BASE",
	"QUINE_WORKSPACE_INIT_ERROR",
}

// runtimeOwnedNonRegistryEnv reports whether a name is runtime-emitted
// plumbing outside the registry: masked from inheritance everywhere, and
// re-stamped by whichever executor owns it.
func runtimeOwnedNonRegistryEnv(name string) bool {
	if name == EnvAgentRoot {
		return true
	}
	for _, n := range runtimeOwnedOverlayNames {
		if name == n {
			return true
		}
	}
	if isRegistryKnob(name) {
		// A registry knob is governed by its Mutability, never by the prefix
		// rules below: QUINE_WORKSPACE_OVERLAY_DRIVER shares the overlay
		// mechanics prefix but is an exec-boundary knob an operator or agent
		// authors, not mount plumbing the runtime stamps.
		return false
	}
	return strings.HasPrefix(name, envJobPrefix) || strings.HasPrefix(name, envWorkspaceOverlayPrefix)
}

// shellLineageMasked reports whether a name is a lineage mark that must not
// cross into an ordinary shell child (brief E4). Depth and ParentSession are
// facts about membership in THIS agent tree; a program started from a shell is
// not a member, and letting it inherit a depth would make a founder look like
// a descendant.
func shellLineageMasked(name string) bool {
	return name == EnvDepth || name == EnvParentSession
}

// BoundaryBehavior classifies one env name at one boundary.
//
// It is derived from the registry rather than from a hand-kept list:
// Mutability already IS the authority taxonomy (registry brief red line 4), so
// a knob cannot drift out of governance without failing the agreement test in
// envmodel_test.go.
func BoundaryBehavior(name string, b Boundary) EnvBehavior {
	if runtimeOwnedNonRegistryEnv(name) {
		return EnvMasked
	}
	if b == BoundaryShell && shellLineageMasked(name) {
		return EnvMasked
	}
	knob, ok := KnobByEnv(name)
	if !ok {
		// Foreign vars (PATH, HOME, an operator's own) and QUINE_* names this
		// binary does not know are inherited untouched. The unknown-QUINE_*
		// case is the version-skew stance: a successor whose registry dropped
		// a knob ignores the inbound name rather than deadlocking on it.
		return EnvFree
	}
	switch knob.Mutability {
	case MutRuntimeEmitted:
		return EnvMasked
	case MutOperatorOnly, MutSubstratePinned:
		return EnvPinned
	case MutExecBoundary:
		return EnvFree
	}
	return EnvFree
}

// --- the override ---

// EnvOverrideDirName / EnvOverrideName locate the agent-authored child-env
// policy under the agent root.
const (
	EnvOverrideDirName = "env"
	EnvOverrideName    = "override"
)

// EnvOverrideDir returns <agentRoot>/config/env, or "" without an agent root.
func (c *Config) EnvOverrideDir() string {
	root := strings.TrimSpace(c.AgentRoot())
	if root == "" {
		return ""
	}
	return filepath.Join(root, "config", EnvOverrideDirName)
}

// EnvOverridePath returns <agentRoot>/config/env/override, or "".
func (c *Config) EnvOverridePath() string {
	dir := c.EnvOverrideDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, EnvOverrideName)
}

// EnvOverride is a validated child-env policy: what this process passes on.
//
// Sets and Unsets are disjoint (a name appears in at most one), and Names
// preserves file order so every rendering of the policy is stable.
type EnvOverride struct {
	Sets   map[string]string
	Unsets map[string]bool
	Names  []string
}

// IsEmpty reports whether the override changes nothing.
func (o *EnvOverride) IsEmpty() bool {
	return o == nil || (len(o.Sets) == 0 && len(o.Unsets) == 0)
}

// EnvOverrideError is a whole-file rejection carrying every violation found.
// Nothing partial ever lands: a file with one bad line is rejected entire, so
// a retry after fixing that line applies exactly what was written.
type EnvOverrideError struct {
	Violations []string
}

func (e *EnvOverrideError) Error() string {
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

// envKeyRE is the accepted key shape. Anything else is a loud syntax
// violation, never a silently skipped line.
var envKeyRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ParseEnvOverride parses and validates a child-env policy against the running
// binary's compiled registry.
//
// Grammar (POSIX-env style, because agents write this file with sh):
//
//	KEY=VALUE   set. VALUE is everything after the first '=' verbatim to end of
//	            line: no quoting, no quote stripping, no ${...} expansion, no
//	            escapes, no `export` prefix. An empty value is a legal SET (the
//	            child receives the name, set and empty).
//	KEY         unset. A line with no '=' is not a legal assignment in any
//	            env-file dialect, so the meaning is unambiguous: remove this
//	            name from what the child inherits.
//	# comment   comments; blank and whitespace-only lines are skipped.
//
// A trailing "\r" is tolerated per line. Naming the same key twice is rejected:
// this is a policy, not a shell script, so last-wins ambiguity bounces instead
// of being guessed at.
//
// Validation checks legality, not wisdom (registry brief D4). A well-typed
// suicidal policy passes; real death keeps selection honest.
func ParseEnvOverride(content []byte) (*EnvOverride, error) {
	override := &EnvOverride{Sets: map[string]string{}, Unsets: map[string]bool{}}
	firstLine := map[string]int{}
	var violations []string

	for i, raw := range strings.Split(string(content), "\n") {
		lineNo := i + 1
		line := strings.TrimSuffix(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, isSet := strings.Cut(line, "=")
		if !isSet {
			// Unset form: the whole (trimmed) line is the name.
			key = trimmed
		}
		if !envKeyRE.MatchString(key) {
			violations = append(violations, fmt.Sprintf(
				"line %d: key %q is not a legal env name (want [A-Za-z_][A-Za-z0-9_]* with no surrounding whitespace; KEY=VALUE sets, a bare KEY unsets)", lineNo, key))
			continue
		}
		if prev, dup := firstLine[key]; dup {
			violations = append(violations, fmt.Sprintf(
				"line %d: %s named twice (first at line %d); state each name exactly once", lineNo, key, prev))
			continue
		}
		firstLine[key] = lineNo

		if v := validateOverrideName(key, value, isSet); v != "" {
			violations = append(violations, fmt.Sprintf("line %d: %s", lineNo, v))
			continue
		}
		override.Names = append(override.Names, key)
		if isSet {
			override.Sets[key] = value
		} else {
			override.Unsets[key] = true
		}
	}

	if len(violations) > 0 {
		return nil, &EnvOverrideError{Violations: violations}
	}
	return override, nil
}

// validateOverrideName checks one policy line against the compiled registry.
// Returns "" when legal, otherwise an agent-legible violation naming what is
// wrong and what the legal form is.
func validateOverrideName(name, value string, isSet bool) string {
	// A name governed at ANY boundary is rejected here: the override is one
	// policy applied to every process this one builds, so it may only carry
	// names that are free everywhere.
	for _, b := range Boundaries() {
		switch BoundaryBehavior(name, b) {
		case EnvMasked:
			if runtimeOwnedNonRegistryEnv(name) || !isRegistryKnob(name) {
				return fmt.Sprintf("%s: runtime-owned name — the runtime writes it into each process it builds; it is not yours to pass on", name)
			}
			return fmt.Sprintf("%s: runtime-emitted knob — the runtime writes it at each boundary where it means something (see config/registry.json)", name)
		case EnvPinned:
			knob, _ := KnobByEnv(name)
			return fmt.Sprintf("%s: mutability %q — pinned, not settable through this file (see config/registry.json)", name, knob.Mutability)
		}
	}

	if !strings.HasPrefix(name, "QUINE_") {
		// Foreign names are free-form: this file is a general child-env
		// manager, and `FOO=bar` for a tool you built is a legitimate use.
		return ""
	}

	knob, ok := KnobByEnv(name)
	if !ok {
		if retired, wasKnob := RetiredKnobByEnv(name); wasKnob {
			if retired.Replacement != "" {
				return fmt.Sprintf("%s: retired knob (%s); use %s instead", name, retired.Decision, retired.Replacement)
			}
			return fmt.Sprintf("%s: retired knob (%s); it no longer exists", name, retired.Decision)
		}
		return fmt.Sprintf("%s: unknown env name — not a capability-registry knob (see config/registry.json)", name)
	}
	if !isSet {
		// Unsetting a knob is legal and means "let the child take the
		// compiled default", which is exactly what absence means.
		return ""
	}
	if value == "" {
		// Set-but-empty is a real state for some knobs (QUINE_DATA_DIR
		// resolves it to .quine/), so it is legal for every type and left to
		// the child's own Load() to resolve.
		return ""
	}
	switch knob.Type.Kind {
	case TypeBool:
		if value != "0" && value != "1" {
			return fmt.Sprintf("%s: invalid bool value %q (must be \"0\" or \"1\")", name, value)
		}
	case TypeInt:
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Sprintf("%s: invalid int value %q", name, value)
		}
	case TypeEnum:
		for _, legal := range knob.Type.Enum {
			if value == legal {
				return ""
			}
		}
		return fmt.Sprintf("%s: invalid enum value %q (legal: %s)", name, value, strings.Join(knob.Type.Enum, "|"))
	case TypeString:
		// free-form
	}
	return ""
}

func isRegistryKnob(name string) bool {
	_, ok := KnobByEnv(name)
	return ok
}

// ReadEnvOverride reads and validates the child-env policy at path.
//
// An absent file — or no agent root at all — is a clean no-op, not an error:
// no policy is the normal state, and it means "children inherit what I was
// given". The returned override is always non-nil so callers never have to
// distinguish "no file" from "empty file".
//
// The file is never rewritten here. A rejected file is left exactly as the
// agent wrote it, so fixing the named line and retrying applies precisely what
// was written — no partial application, no silent repair, nothing lost.
//
// Every envp constructor calls this at construction time, not once at startup:
// the agent edits config/env/override with ordinary shell writes, and a cached
// copy would mean a policy that has been written but does not yet apply. That
// is the class of lie this model exists to remove.
func ReadEnvOverride(path string) (*EnvOverride, error) {
	empty := &EnvOverride{Sets: map[string]string{}, Unsets: map[string]bool{}}
	path = strings.TrimSpace(path)
	if path == "" {
		return empty, nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return empty, nil
		}
		return nil, fmt.Errorf("read child-env override %s: %w", path, err)
	}
	override, err := ParseEnvOverride(content)
	if err != nil {
		return nil, fmt.Errorf("child-env override %s rejected (whole file; nothing applied; the file is left intact — fix or remove the named lines, then retry): %w", path, err)
	}
	return override, nil
}

// --- the pipeline ---

// BuildChildEnv constructs the environment of a process this runtime creates:
//
//	child = stamps ⊕ override ⊕ (environ − mask(boundary))
//
// This is REAL envp. It carries secret values a child legitimately needs — a
// child quine cannot reach the provider without QUINE_API_KEY. It is not a
// display surface, and it must never be used as one: the rendering path is
// RenderEffectiveEnv, a separate function returning a different type, so that
// redaction cannot be lost to a refactor that "reuses" this one (brief E13).
//
// stamps are runtime-owned facts about the child ("KEY=VALUE" form), supplied
// by the caller because only the caller knows them (the child's depth, its
// agent root, the mount it will run inside). They are applied last and win
// everything.
func BuildChildEnv(b Boundary, environ []string, override *EnvOverride, stamps []string) []string {
	values := make(map[string]string, len(environ)+len(stamps))
	order := make([]string, 0, len(environ)+len(stamps))

	put := func(name, value string) {
		if _, seen := values[name]; !seen {
			order = append(order, name)
		}
		values[name] = value
	}

	// 1. inherit, minus the mask.
	for _, kv := range environ {
		name, value, ok := strings.Cut(kv, "=")
		if !ok || name == "" {
			continue
		}
		if BoundaryBehavior(name, b) == EnvMasked {
			continue
		}
		put(name, value)
	}

	// 2. the override: unset removes, set replaces or introduces. A line the
	//    validator would have rejected is inert here — the file may have been
	//    written raw through sh, and use-time enforcement is what makes the
	//    boundary un-fakeable rather than merely well-advertised.
	if override != nil {
		for _, name := range override.Names {
			if BoundaryBehavior(name, b) != EnvFree {
				continue
			}
			if override.Unsets[name] {
				if _, present := values[name]; present {
					delete(values, name)
					order = removeName(order, name)
				}
				continue
			}
			put(name, override.Sets[name])
		}
	}

	// 3. stamps win.
	for _, kv := range stamps {
		name, value, ok := strings.Cut(kv, "=")
		if !ok || name == "" {
			continue
		}
		put(name, value)
	}

	out := make([]string, 0, len(order))
	for _, name := range order {
		out = append(out, name+"="+values[name])
	}
	return out
}

func removeName(order []string, name string) []string {
	for i, n := range order {
		if n == name {
			return append(order[:i:i], order[i+1:]...)
		}
	}
	return order
}
