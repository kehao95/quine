package config

// envmodel_test.go is the L1 substrate property suite for the env boundary
// model (envmodel.go).
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//   (Layer 6 "L1 — substrate" is the checklist this file implements)
//
// Deliberately build-tag-free (brief E7/OQ-9): these raw-table and raw-parse
// tests are the degraded-path coverage for non-Linux / no-FUSE hosts.

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// --- helpers ---

// childEnvMap turns a BuildChildEnv result into name->value for assertions.
func childEnvMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, kv := range env {
		if name, value, ok := strings.Cut(kv, "="); ok {
			m[name] = value
		}
	}
	return m
}

func mustParseOverride(t *testing.T, content string) *EnvOverride {
	t.Helper()
	o, err := ParseEnvOverride([]byte(content))
	if err != nil {
		t.Fatalf("ParseEnvOverride(%q) rejected a policy this test needs accepted: %v", content, err)
	}
	return o
}

// overrideViolations asserts rejection and returns the enumerated violations.
func overrideViolations(t *testing.T, content string) []string {
	t.Helper()
	o, err := ParseEnvOverride([]byte(content))
	if err == nil {
		t.Fatalf("ParseEnvOverride(%q) accepted; want whole-file rejection", content)
	}
	if o != nil {
		t.Fatalf("ParseEnvOverride(%q): rejection must not also return a partial override (nothing partial ever lands), got %+v", content, o)
	}
	var oe *EnvOverrideError
	if !errors.As(err, &oe) {
		t.Fatalf("ParseEnvOverride(%q) error is %T, want *EnvOverrideError (callers need the enumerated violations)", content, err)
	}
	if len(oe.Violations) == 0 {
		t.Fatalf("ParseEnvOverride(%q): EnvOverrideError with zero violations is unreportable", content)
	}
	return oe.Violations
}

// --- 1. table completeness / agreement ---

// TestBoundaryBehaviorTableComplete: every registry knob classifies to exactly
// one EnvBehavior at every boundary, and the classification is a pure function
// of Mutability per the governance table (brief Layer 3):
//
//	runtime-emitted               -> masked
//	operator-only/substrate-pinned -> pinned
//	exec-boundary                 -> free
//
// with exactly one boundary-local exception: the E4 shell lineage mask.
// A knob added with a NEW Mutability value has no row in wantByMutability and
// fails loudly here instead of silently falling through to EnvFree.
func TestBoundaryBehaviorTableComplete(t *testing.T) {
	wantByMutability := map[Mutability]EnvBehavior{
		MutRuntimeEmitted:  EnvMasked,
		MutOperatorOnly:    EnvPinned,
		MutSubstratePinned: EnvPinned,
		MutExecBoundary:    EnvFree,
	}
	for _, k := range Registry {
		want, known := wantByMutability[k.Mutability]
		if !known {
			t.Errorf("%s: Mutability %q has no boundary-behavior mapping — a new mutability class must be added to the governance table (envmodel.go BoundaryBehavior) and to this test, or knobs will silently classify as free", k.Env, k.Mutability)
			continue
		}
		for _, b := range Boundaries() {
			expect := want
			if b == BoundaryShell && (k.Env == EnvDepth || k.Env == EnvParentSession) {
				expect = EnvMasked // brief E4: sh children are unmarked new roots
			}
			if got := BoundaryBehavior(k.Env, b); got != expect {
				t.Errorf("BoundaryBehavior(%s, %s) = %q, want %q: a %q knob must not drift out of the governance table (brief Layer 3)", k.Env, b, got, expect, k.Mutability)
			}
		}
	}
}

// TestBoundaryMaskAgreement: three-way agreement at the managed boundaries —
// masked set at fork/exec == registry MutRuntimeEmitted class ==
// ProcessIdentityEnvNames. A knob added later with runtime-emitted mutability
// but missing from the identity slice (or vice versa) fails here, so the
// shared mask can never silently diverge from its derivation source.
func TestBoundaryMaskAgreement(t *testing.T) {
	identity := make(map[string]bool, len(ProcessIdentityEnvNames))
	for _, env := range ProcessIdentityEnvNames {
		identity[env] = true
		if _, ok := KnobByEnv(env); !ok {
			t.Errorf("ProcessIdentityEnvNames contains %q which is not a registry knob — the identity slice must stay a pure projection of the registry", env)
		}
	}
	for _, b := range []Boundary{BoundaryChild, BoundarySelf} {
		for _, k := range Registry {
			masked := BoundaryBehavior(k.Env, b) == EnvMasked
			if masked != (k.Mutability == MutRuntimeEmitted) {
				t.Errorf("%s at %s: masked=%v but Mutability=%q — the mask must derive from MutRuntimeEmitted exactly", k.Env, b, masked, k.Mutability)
			}
			if masked != identity[k.Env] {
				t.Errorf("%s at %s: masked=%v but in ProcessIdentityEnvNames=%v — the mask and the identity slice must agree", k.Env, b, masked, identity[k.Env])
			}
		}
	}
	// QUINE_CONTEXT_BOOTSTRAP was a hardcoded mask exception before it was
	// promoted to a runtime-emitted knob; the whole point of the promotion is
	// that its mask now DERIVES. Assert it explicitly at every boundary.
	for _, b := range Boundaries() {
		if got := BoundaryBehavior(EnvContextBootstrap, b); got != EnvMasked {
			t.Errorf("BoundaryBehavior(%s, %s) = %q, want %q: the context-bootstrap mask must derive from its registry mutability, not from a hardcoded exception", EnvContextBootstrap, b, got, EnvMasked)
		}
	}
}

// --- 2. pipeline properties, per boundary ---

// TestEnvPipelineForeignPassthrough: operator-authored foreign vars cross
// every boundary byte-exact — including values containing '=' and spaces, and
// unknown QUINE_* names (version-skew stance: tolerated in environ even though
// the override rejects them).
func TestEnvPipelineForeignPassthrough(t *testing.T) {
	environ := []string{
		"PATH=/usr/bin:/bin",
		"FOO=bar",
		"WEIRD=a=b c ",
		"QUINE_FROM_THE_FUTURE=1", // unknown QUINE_*: inherited untouched
	}
	for _, b := range Boundaries() {
		got := BuildChildEnv(b, environ, nil, nil)
		for _, kv := range environ {
			found := false
			for _, out := range got {
				if out == kv {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("boundary %s: foreign var %q did not pass through byte-exact; got %v", b, kv, got)
			}
		}
	}
}

// TestEnvPipelineMasksRuntimeEmitted: masked names never appear in child env
// even when present in environ. QUINE_WORKSPACE_CURRENT_REVISION leaks into sh
// children under the old synthesis (brief Layer 2 #8); closing that leak is an
// intended fix, so sh is asserted alongside fork/exec.
func TestEnvPipelineMasksRuntimeEmitted(t *testing.T) {
	environ := []string{
		EnvSessionID + "=sess-forged",
		EnvContextBootstrap + "=/tmp/ctx",
		EnvWorkspaceCurrentRevision + "=abc123",
		"PATH=/usr/bin",
	}
	maskedNames := []string{EnvSessionID, EnvContextBootstrap, EnvWorkspaceCurrentRevision}
	for _, b := range Boundaries() {
		out := childEnvMap(BuildChildEnv(b, environ, nil, nil))
		for _, name := range maskedNames {
			if v, present := out[name]; present {
				t.Errorf("boundary %s: masked name %s crossed the boundary with value %q — masked names must never be inherited (for %s at sh this is the Layer 2 #8 leak, reopened)", b, name, v, EnvWorkspaceCurrentRevision)
			}
		}
		if out["PATH"] != "/usr/bin" {
			t.Errorf("boundary %s: PATH should survive the mask untouched, got %q", b, out["PATH"])
		}
	}
}

// TestEnvPipelineUnsetKnobAbsent: AN UNSET KNOB IS ABSENT from child env.
//
// This is the inversion of the old TestChildEnv, which asserted every registry
// knob is ALWAYS present, defaults included. Manufactured `KNOB=0` lines are
// what taught a succ-52 founder that process creation was physically
// impossible (brief Layer 0: QUINE_SPAWN_ENABLED=0 / QUINE_MAX_DEPTH=0, two
// lines no operator authored). Absence means "compiled default", and the
// pipeline must preserve absence.
func TestEnvPipelineUnsetKnobAbsent(t *testing.T) {
	environ := []string{
		EnvModelID + "=claude-x",
		"PATH=/usr/bin",
	}
	neverSet := []string{EnvSpawnEnabled, EnvMaxDepth, EnvForkEnabled}
	for _, b := range Boundaries() {
		out := childEnvMap(BuildChildEnv(b, environ, nil, nil))
		for _, name := range neverSet {
			if v, present := out[name]; present {
				t.Errorf("boundary %s: %s=%q appeared in child env but was never set anywhere — the pipeline manufactured evidence (brief Layer 0: this is the synthesis defect the model exists to remove)", b, name, v)
			}
		}
		if out[EnvModelID] != "claude-x" {
			t.Errorf("boundary %s: the one knob that WAS set (%s) must pass through, got %q", b, EnvModelID, out[EnvModelID])
		}
	}
}

// TestEnvPipelinePrecedence: stamps beat override beats environ.
func TestEnvPipelinePrecedence(t *testing.T) {
	environ := []string{
		"FOO=from-environ",
		"BAR=from-environ",
		"BAZ=from-environ",
	}
	override := mustParseOverride(t, "FOO=from-override\nBAR=from-override\n")
	stamps := []string{"FOO=from-stamp"}
	for _, b := range Boundaries() {
		out := childEnvMap(BuildChildEnv(b, environ, override, stamps))
		if out["FOO"] != "from-stamp" {
			t.Errorf("boundary %s: FOO=%q, want from-stamp — stamps are runtime-owned facts and must win everything", b, out["FOO"])
		}
		if out["BAR"] != "from-override" {
			t.Errorf("boundary %s: BAR=%q, want from-override — an override line must beat an inherited value", b, out["BAR"])
		}
		if out["BAZ"] != "from-environ" {
			t.Errorf("boundary %s: BAZ=%q, want from-environ — an untouched inherited value must survive verbatim", b, out["BAZ"])
		}
	}
}

// TestEnvPipelineShellLineageMask (brief E4, both directions): QUINE_DEPTH and
// QUINE_PARENT_SESSION are masked at the sh boundary — a program started from
// a shell is an unmarked new root, not a tree member — but NOT masked at
// fork/exec, where they are pinned inherits the executors stamp over. The
// fork/exec direction also exercises E11's raw material: a genuinely nonzero
// QUINE_DEPTH=2 survives the exec pipeline instead of being reset.
func TestEnvPipelineShellLineageMask(t *testing.T) {
	environ := []string{
		EnvDepth + "=2",
		EnvParentSession + "=parent-abc",
		"PATH=/bin",
	}
	out := childEnvMap(BuildChildEnv(BoundaryShell, environ, nil, nil))
	for _, name := range []string{EnvDepth, EnvParentSession} {
		if v, present := out[name]; present {
			t.Errorf("sh child inherited lineage mark %s=%q — an sh child is an unmarked new root (brief E4); letting it inherit depth makes a founder look like a descendant", name, v)
		}
	}
	for _, b := range []Boundary{BoundaryChild, BoundarySelf} {
		out := childEnvMap(BuildChildEnv(b, environ, nil, nil))
		if out[EnvDepth] != "2" {
			t.Errorf("boundary %s: %s=%q, want \"2\" — lineage names must NOT be masked at managed boundaries (the executor stamps the true value on top; at exec, depth is preserved, not reset — brief E11)", b, EnvDepth, out[EnvDepth])
		}
		if out[EnvParentSession] != "parent-abc" {
			t.Errorf("boundary %s: %s=%q, want \"parent-abc\" — lineage names must NOT be masked at managed boundaries", b, EnvParentSession, out[EnvParentSession])
		}
	}
}

// TestEnvPipelineBareKeyUnset: a bare-KEY override line removes an inherited
// var — foreign or free knob alike — at every boundary.
func TestEnvPipelineBareKeyUnset(t *testing.T) {
	environ := []string{
		"FOO=bar",
		EnvMaxTurns + "=9",
		"PATH=/bin",
	}
	override := mustParseOverride(t, "FOO\n"+EnvMaxTurns+"\n")
	for _, b := range Boundaries() {
		out := childEnvMap(BuildChildEnv(b, environ, override, nil))
		if v, present := out["FOO"]; present {
			t.Errorf("boundary %s: FOO=%q survived a bare-KEY unset — the override's unset op must remove inherited vars", b, v)
		}
		if v, present := out[EnvMaxTurns]; present {
			t.Errorf("boundary %s: %s=%q survived a bare-KEY unset — unsetting a free knob means the child takes the compiled default", b, EnvMaxTurns, v)
		}
		if out["PATH"] != "/bin" {
			t.Errorf("boundary %s: PATH should be untouched by unrelated unsets, got %q", b, out["PATH"])
		}
	}
}

// TestEnvPipelineGovernedOverrideLinesInert: an override line naming a pinned
// or masked name is INERT at use time. The parser rejects such a file at write
// time, but the file is raw-writable through sh, so use-time enforcement is
// what makes the boundary un-fakeable rather than merely well-advertised. The
// override here is built by hand to simulate exactly that raw write.
func TestEnvPipelineGovernedOverrideLinesInert(t *testing.T) {
	override := &EnvOverride{
		Sets: map[string]string{
			EnvSessionID: "forged-session", // masked everywhere
			EnvMaxDepth:  "99",             // pinned (E5)
			EnvDepth:     "7",              // masked at sh, pinned at fork/exec
		},
		Unsets: map[string]bool{
			EnvAPIKey: true, // pinned: an unset must be just as inert as a set
		},
		Names: []string{EnvSessionID, EnvMaxDepth, EnvDepth, EnvAPIKey},
	}
	environ := []string{
		EnvAPIKey + "=real-key",
		EnvMaxDepth + "=3",
		EnvDepth + "=1",
		"PATH=/bin",
	}
	for _, b := range Boundaries() {
		out := childEnvMap(BuildChildEnv(b, environ, override, nil))
		if v, present := out[EnvSessionID]; present {
			t.Errorf("boundary %s: raw-written override forged %s=%q into a child — masked names must be inert in the override at use time", b, EnvSessionID, v)
		}
		if out[EnvMaxDepth] != "3" {
			t.Errorf("boundary %s: %s=%q, want the inherited \"3\" — a pinned name's override set must be inert at use time (raw write bypassed the parser)", b, EnvMaxDepth, out[EnvMaxDepth])
		}
		if out[EnvAPIKey] != "real-key" {
			t.Errorf("boundary %s: %s=%q, want the inherited \"real-key\" — a pinned name's override UNSET must be inert at use time (dropping the credential would orphan the child)", b, EnvAPIKey, out[EnvAPIKey])
		}
		switch b {
		case BoundaryShell:
			if v, present := out[EnvDepth]; present {
				t.Errorf("sh: %s=%q crossed via override — masked-at-sh lineage names must be inert in the override too", EnvDepth, v)
			}
		default:
			if out[EnvDepth] != "1" {
				t.Errorf("boundary %s: %s=%q, want the inherited \"1\" — pinned at managed boundaries, so the override's forged 7 must be inert", b, EnvDepth, out[EnvDepth])
			}
		}
	}
}

// --- 3. override parser ---

func TestOverrideParserGrammar(t *testing.T) {
	t.Run("accepted forms", func(t *testing.T) {
		cases := []struct {
			name    string
			content string
			sets    map[string]string
			unsets  []string
			order   []string
		}{
			{"simple set", "FOO=bar\n", map[string]string{"FOO": "bar"}, nil, []string{"FOO"}},
			{"bare KEY unsets", "FOO\n", nil, []string{"FOO"}, []string{"FOO"}},
			{"comments and blanks skipped", "# a comment\n\n   \n  # indented comment\nFOO=bar\n", map[string]string{"FOO": "bar"}, nil, []string{"FOO"}},
			{"crlf tolerated", "FOO=bar\r\nBAR=baz\r\n", map[string]string{"FOO": "bar", "BAR": "baz"}, nil, []string{"FOO", "BAR"}},
			{"no quote stripping", "FOO=\"quoted\"\n", map[string]string{"FOO": "\"quoted\""}, nil, []string{"FOO"}},
			{"no expansion", "FOO=$HOME/x\n", map[string]string{"FOO": "$HOME/x"}, nil, []string{"FOO"}},
			{"empty value is a legal SET", "FOO=\n", map[string]string{"FOO": ""}, nil, []string{"FOO"}},
			{"value verbatim after first =", "FOO=a=b=c\n", map[string]string{"FOO": "a=b=c"}, nil, []string{"FOO"}},
			{"value keeps interior and trailing spaces", "FOO=hello world \n", map[string]string{"FOO": "hello world "}, nil, []string{"FOO"}},
			{"file order preserved in Names", "B=1\nA=2\nC\n", map[string]string{"B": "1", "A": "2"}, []string{"C"}, []string{"B", "A", "C"}},
			{"no trailing newline", "FOO=bar", map[string]string{"FOO": "bar"}, nil, []string{"FOO"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				o := mustParseOverride(t, tc.content)
				if len(o.Sets) != len(tc.sets) {
					t.Errorf("Sets = %v, want %v", o.Sets, tc.sets)
				}
				for k, want := range tc.sets {
					got, present := o.Sets[k]
					if !present {
						t.Errorf("set %s missing from parse result", k)
					} else if got != want {
						t.Errorf("Sets[%s] = %q, want %q — VALUE must be verbatim after the first '='", k, got, want)
					}
					if o.Unsets[k] {
						t.Errorf("%s is in both Sets and Unsets — they must be disjoint", k)
					}
				}
				if len(o.Unsets) != len(tc.unsets) {
					t.Errorf("Unsets = %v, want %v", o.Unsets, tc.unsets)
				}
				for _, k := range tc.unsets {
					if !o.Unsets[k] {
						t.Errorf("bare line %q did not parse as an unset", k)
					}
				}
				if fmt.Sprint(o.Names) != fmt.Sprint(tc.order) {
					t.Errorf("Names = %v, want %v — file order must be preserved so renderings are stable", o.Names, tc.order)
				}
			})
		}
	})

	t.Run("illegal key shapes rejected loudly", func(t *testing.T) {
		for _, content := range []string{
			"FOO =bar\n", // whitespace inside key
			"1FOO=bar\n", // leading digit
			"=bar\n",     // empty key
			"FO O\n",     // whitespace inside bare key
		} {
			violations := overrideViolations(t, content)
			if !strings.Contains(violations[0], "not a legal env name") {
				t.Errorf("ParseEnvOverride(%q): violation %q should name the legal key shape, never silently skip the line", content, violations[0])
			}
		}
	})
}

func TestOverrideParserDuplicateNameRejected(t *testing.T) {
	for _, content := range []string{
		"FOO=a\nFOO=b\n", // set twice
		"FOO=a\nFOO\n",   // set then unset
		"FOO\nFOO\n",     // unset twice
	} {
		violations := overrideViolations(t, content)
		if !strings.Contains(violations[0], "named twice") {
			t.Errorf("ParseEnvOverride(%q): violation %q should reject the duplicate by name — a policy has no last-wins semantics", content, violations[0])
		}
	}
}

// TestOverrideParserWholeFileRejection: one bad line rejects the entire file,
// and the error enumerates EVERY violation, not just the first.
func TestOverrideParserWholeFileRejection(t *testing.T) {
	content := "FOO=fine\n" + EnvSessionID + "=x\nQUINE_NO_SUCH_KNOB=1\n"
	violations := overrideViolations(t, content)
	if len(violations) != 2 {
		t.Fatalf("got %d violations %v, want 2 — the rejection must enumerate every bad line so one retry fixes the whole file", len(violations), violations)
	}
	joined := strings.Join(violations, "\n")
	if !strings.Contains(joined, EnvSessionID) {
		t.Errorf("violations %v do not name %s", violations, EnvSessionID)
	}
	if !strings.Contains(joined, "QUINE_NO_SUCH_KNOB") {
		t.Errorf("violations %v do not name QUINE_NO_SUCH_KNOB", violations)
	}
	// The rendered error must also carry both, since gates surface Error().
	errText := (&EnvOverrideError{Violations: violations}).Error()
	if !strings.Contains(errText, EnvSessionID) || !strings.Contains(errText, "QUINE_NO_SUCH_KNOB") {
		t.Errorf("EnvOverrideError.Error() = %q must enumerate both violations", errText)
	}
}

func TestOverrideParserQuineNameValidation(t *testing.T) {
	t.Run("unknown QUINE_* rejected", func(t *testing.T) {
		violations := overrideViolations(t, "QUINE_NO_SUCH_KNOB=1\n")
		if !strings.Contains(violations[0], "unknown") {
			t.Errorf("violation %q should say the name is unknown and point at the catalog", violations[0])
		}
	})
	t.Run("retired QUINE_* rejected with replacement pointer", func(t *testing.T) {
		violations := overrideViolations(t, "QUINE_SH_TIMEOUT=30\n")
		if !strings.Contains(violations[0], "retired") || !strings.Contains(violations[0], EnvShDefaultTimeout) {
			t.Errorf("violation %q should say retired and point at the successor %s", violations[0], EnvShDefaultTimeout)
		}
	})
	t.Run("retired QUINE_* without replacement rejected", func(t *testing.T) {
		violations := overrideViolations(t, "QUINE_SMART_MODEL_ID=x\n")
		if !strings.Contains(violations[0], "retired") {
			t.Errorf("violation %q should say the knob is retired", violations[0])
		}
	})
	t.Run("runtime-emitted rejected", func(t *testing.T) {
		overrideViolations(t, EnvSessionID+"=forged\n")
	})
	t.Run("pinned rejected, set and unset alike", func(t *testing.T) {
		for _, content := range []string{
			EnvAPIKey + "=sk-mine\n",
			EnvAPIKey + "\n", // bare unset of a pinned name is just as illegal
			EnvMaxDepth + "=5\n",
			EnvDepth + "=3\n",
		} {
			overrideViolations(t, content)
		}
	})
}

func TestOverrideParserTypeValidation(t *testing.T) {
	accepted := []string{
		EnvForkEnabled + "=1\n",
		EnvForkEnabled + "=0\n",
		EnvForkEnabled + "=\n", // empty value is a legal SET for every type
		EnvForkEnabled + "\n",  // bare unset of a free knob = take compiled default
		EnvMaxTurns + "=25\n",
		EnvShNetwork + "=none\n",
		EnvInitialUserMessage + "=hello world, free-form string\n",
	}
	for _, content := range accepted {
		if _, err := ParseEnvOverride([]byte(content)); err != nil {
			t.Errorf("ParseEnvOverride(%q) rejected a well-typed free-knob line: %v", content, err)
		}
	}
	rejected := []struct{ content, wantSubstr string }{
		{EnvForkEnabled + "=true\n", `must be "0" or "1"`},
		{EnvForkEnabled + "=2\n", `must be "0" or "1"`},
		{EnvMaxTurns + "=abc\n", "invalid int"},
		{EnvShNetwork + "=offline\n", "invalid enum"},
	}
	for _, tc := range rejected {
		violations := overrideViolations(t, tc.content)
		if !strings.Contains(violations[0], tc.wantSubstr) {
			t.Errorf("ParseEnvOverride(%q): violation %q should contain %q so the writer learns the legal form", tc.content, violations[0], tc.wantSubstr)
		}
	}
	// An enum violation must teach the legal values.
	violations := overrideViolations(t, EnvShNetwork+"=offline\n")
	if !strings.Contains(violations[0], "host") || !strings.Contains(violations[0], "none") {
		t.Errorf("enum violation %q should enumerate the legal values host|none", violations[0])
	}
}

// TestOverrideParserForeignNamesAccepted (brief E9): the override is a general
// child-env manager, so non-QUINE_ names are free-form — set and unset alike.
func TestOverrideParserForeignNamesAccepted(t *testing.T) {
	o := mustParseOverride(t, "FOO=bar\nPATH=/opt/bin\nMY_TOOL_CONFIG=/etc/x\nLC_ALL\n")
	if o.Sets["FOO"] != "bar" || o.Sets["PATH"] != "/opt/bin" || o.Sets["MY_TOOL_CONFIG"] != "/etc/x" {
		t.Errorf("foreign sets did not parse: %v — the override is a general child-env manager (brief E9)", o.Sets)
	}
	if !o.Unsets["LC_ALL"] {
		t.Errorf("bare foreign name LC_ALL did not parse as an unset: %v", o.Unsets)
	}
}

// TestOverrideParserBudgetPinsE5: the E5 regression guard. Tree budgets
// (MaxDepth/MaxAgents/MaxConcurrent) reclassified operator-only — a lineage
// must not relax its own tree budget through the mediated channel — while
// MaxTurns stays agent-overridable (mortality self-modification remains an
// agent power).
func TestOverrideParserBudgetPinsE5(t *testing.T) {
	for _, env := range []string{EnvMaxDepth, EnvMaxAgents, EnvMaxConcurrent} {
		violations := overrideViolations(t, env+"=5\n")
		if !strings.Contains(violations[0], "pinned") {
			t.Errorf("%s: violation %q should say the knob is pinned (E5 reclassified tree budgets to operator-only)", env, violations[0])
		}
	}
	if _, err := ParseEnvOverride([]byte(EnvMaxTurns + "=5\n")); err != nil {
		t.Errorf("%s=5 rejected: %v — MaxTurns must stay agent-overridable (E5 keeps mortality self-modification an agent power)", EnvMaxTurns, err)
	}
}

// --- 6. stamps ---

// TestStampedEnvNamesMatchBuilders pins the declared stamp vocabulary
// (StampedEnvNames) to what the builders actually emit for a fully-populated
// Config. A stamp is the only place a resolved Config value may enter a child
// env, so StampedEnvNames is the drift guard / review checkpoint on that set:
// every addition to it deserves the "is this a runtime-owned fact, or a policy
// default nobody chose" question before it is allowed to exist.
//
// This is the guard the four re-synthesized workspace stamps slipped past: they
// were emitted into every workspace-enabled fork child while appearing on no
// governance surface at all, so nothing — no test, no table, no log — could tell
// the agent that its override of QUINE_WORKSPACE_BACKEND would be silently
// overwritten. A stamp that is not declared here cannot be described to the
// reader, so it is not allowed to exist.
func TestStampedEnvNamesMatchBuilders(t *testing.T) {
	cfg := &Config{
		Identity: Identity{SessionID: "sess-1", RunID: "run-1", TapeID: "0007", ParentSession: "parent-0", Depth: 2},
		Paths:    Paths{DataDir: "/tmp/data", RetentionDir: "/tmp/retained", WorkDir: "/tmp/work"},
		WorkspaceConfig: WorkspaceConfig{
			WorkspaceEnabled:         true,
			WorkspaceRoot:            "/tmp/ws",
			Workspace:                "/tmp/ws/scope",
			WorkspaceBackend:         "overlay",
			WorkspaceOverlayDriver:   "kernel",
			WorkspaceRevisionMode:    WorkspaceRevisionRestore,
			WorkspaceCommitOnSignal:  true,
			WorkspaceSession:         "ws-sess",
			WorkspaceOwner:           true,
			WorkspaceCurrentRevision: "rev-3",
		},
	}
	builders := map[Boundary][]string{
		BoundaryShell: cfg.ShellStamps(),
		BoundaryChild: cfg.ForkChildStamps(),
		BoundarySelf:  cfg.SelfReentryStamps(),
	}
	for b, stamps := range builders {
		emitted := map[string]bool{}
		for _, kv := range stamps {
			name, _, _ := strings.Cut(kv, "=")
			emitted[name] = true
		}
		declared := map[string]bool{}
		for _, name := range StampedEnvNames(b) {
			declared[name] = true
			if !emitted[name] {
				t.Errorf("boundary %s: StampedEnvNames declares %s but the builder does not emit it for a fully-populated Config — a declared stamp must correspond to something the builder actually stamps", b, name)
			}
		}
		for name := range emitted {
			if !declared[name] {
				t.Errorf("boundary %s: the builder stamps %s but StampedEnvNames does not declare it — it would be written into every child while appearing on no governance surface, which is exactly how four workspace knobs came to silently defeat the override", b, name)
			}
		}
	}
}

// TestForkStampsDoNotSynthesizeFreeKnobs is the regression test for the
// synthesis that survived the cutover.
//
// ForkChildStamps stamped QUINE_WORKSPACE_{BACKEND,OVERLAY_DRIVER,REVISION_MODE,
// COMMIT_ON_SIGNAL} at their RESOLVED values whenever workspace physics were on.
// All four are exec-boundary (free) knobs whose resolved values are pure
// backend-conditional Load() defaults, so the child derives them identically from
// the workspace root it is handed — the stamps carried no information at all.
// What they did carry was QUINE_WORKSPACE_COMMIT_ON_SIGNAL=0: a default-valued
// negative no operator authored, in every workspace-enabled fork child's environ
// and its permanent birth record. That is the morphology of QUINE_SPAWN_ENABLED=0
// (brief Layer 0), reintroduced under the name "stamp".
//
// The pre-existing exactness fixture could not catch it: it built its Config
// without WorkspaceEnabled, so the entire stamp block was unreachable.
func TestForkStampsDoNotSynthesizeFreeKnobs(t *testing.T) {
	cfg := &Config{
		Identity: Identity{SessionID: "sess-1", Depth: 0},
		Paths:    Paths{DataDir: "/tmp/data"},
		WorkspaceConfig: WorkspaceConfig{
			WorkspaceEnabled:        true,
			WorkspaceRoot:           "/tmp/ws",
			Workspace:               "/tmp/ws",
			WorkspaceBackend:        "overlay", // resolved default, not authored
			WorkspaceOverlayDriver:  "kernel",  // resolved default, not authored
			WorkspaceRevisionMode:   WorkspaceRevisionRestore,
			WorkspaceCommitOnSignal: false, // resolved default => the "=0" line
			WorkspaceSession:        "ws-sess",
		},
	}
	// The operator authored the root and nothing else. This is the whole environ.
	environ := []string{envKV(EnvWorkspaceRoot, "/tmp/ws"), envKV(EnvAPIKey, "sk")}
	child := childEnvMap(BuildChildEnv(BoundaryChild, environ, nil, cfg.ForkChildStamps()))

	for _, name := range []string{
		EnvWorkspaceBackend, EnvWorkspaceOverlayDriver, EnvWorkspaceRevisionMode, EnvWorkspaceCommitOnSignal,
	} {
		if v, present := child[name]; present {
			t.Errorf("fork child carries %s=%q, which no operator authored: it is a compiled default the child would derive for itself. A stamp must be a fact the child CANNOT derive (envstamps.go's own test); anything else is the synthesizer coming back under a new name", name, v)
		}
	}
	// Specifically the manufactured negative — the brief's Layer 4 "never again".
	if v, present := child[EnvWorkspaceCommitOnSignal]; present && v == "0" {
		t.Errorf("%s=0 reached a child env: a QUINE_*=0 line nobody authored is precisely what taught a founder its own reproduction was impossible", EnvWorkspaceCommitOnSignal)
	}

	// The stamps that ARE facts survive: scope and lineage the child cannot derive.
	for name, want := range map[string]string{
		EnvWorkspaceRoot:      "/tmp/ws",
		EnvWorkspace:          "/tmp/ws",
		EnvWorkspaceOwner:     "0",
		EnvWorkspaceBootstrap: "ws-sess",
		EnvDepth:              "1",
	} {
		if child[name] != want {
			t.Errorf("fork child %s=%q, want %q — this one IS a runtime-owned fact and must still be stamped", name, child[name], want)
		}
	}
}

// TestForkOverrideReachesFreeWorkspaceKnobs: the other half of the same defect.
// Because the stamps used to beat the override, an agent that wrote
// QUINE_WORKSPACE_BACKEND=direct into config/env/override had it accepted by the
// gate — QUINE_WORKSPACE_BACKEND is a FREE knob — but forked and got a child
// running overlay anyway, silently. Post-fix, the runtime does not stamp that
// knob at all, so what the agent authored in config/env/override is exactly
// what the fork child receives. No rejection, no log line, no silent override:
// the surface said one thing and the boundary did another, which is the
// failure class this whole model exists to remove.
func TestForkOverrideReachesFreeWorkspaceKnobs(t *testing.T) {
	cfg := &Config{
		Identity: Identity{SessionID: "sess-1"},
		Paths:    Paths{DataDir: "/tmp/data"},
		WorkspaceConfig: WorkspaceConfig{
			WorkspaceEnabled: true, WorkspaceRoot: "/tmp/ws", Workspace: "/tmp/ws",
			WorkspaceBackend: "overlay", WorkspaceRevisionMode: WorkspaceRevisionRestore,
		},
	}
	override := mustParseOverride(t, EnvWorkspaceBackend+"=direct\n"+EnvWorkspaceCommitOnSignal+"=1\n")
	child := childEnvMap(BuildChildEnv(BoundaryChild, []string{envKV(EnvWorkspaceRoot, "/tmp/ws")}, override, cfg.ForkChildStamps()))

	if child[EnvWorkspaceBackend] != "direct" {
		t.Errorf("%s=%q, want \"direct\": the override names a FREE knob and the gate accepts it — since the runtime does not stamp this knob, a fork child that silently gets the opposite value means the agent was lied to by the mechanism built to stop lying to it", EnvWorkspaceBackend, child[EnvWorkspaceBackend])
	}
	if child[EnvWorkspaceCommitOnSignal] != "1" {
		t.Errorf("%s=%q, want \"1\": same defect", EnvWorkspaceCommitOnSignal, child[EnvWorkspaceCommitOnSignal])
	}
}

// TestRootStampsAreAbsolute: QUINE_DATA_DIR names the runtime root this process
// JOINED, and the stamp only means that if it means the same thing from any cwd.
// The compiled default is the relative ".quine/", and children do not share our
// cwd — an sh child runs inside the workspace mount, a fork child in the
// executor's WorkDir. Stamping ".quine/" therefore reproduced the divergence the
// stamp exists to prevent (a different root, joined silently, no error anywhere)
// while looking like it had prevented it.
func TestRootStampsAreAbsolute(t *testing.T) {
	cfg := &Config{
		Identity: Identity{SessionID: "sess-1"},
		Paths:    Paths{DataDir: ".quine/", RetentionDir: "retained/"},
	}
	// An sh child needs the runtime root (the `world` binary resolves its spec from
	// it); only the managed boundaries carry the retention root.
	builders := map[Boundary][]string{
		BoundaryShell: cfg.ShellStamps(),
		BoundaryChild: cfg.ForkChildStamps(),
		BoundarySelf:  cfg.SelfReentryStamps(),
	}
	for b, stamps := range builders {
		got := childEnvMap(stamps)
		for _, name := range StampedEnvNames(b) {
			if name != EnvDataDir && name != EnvRetentionDir {
				continue
			}
			v, present := got[name]
			if !present {
				t.Fatalf("boundary %s: %s is declared stamped but absent", b, name)
			}
			if !filepath.IsAbs(v) {
				t.Errorf("boundary %s: %s=%q is relative — the child resolves it against ITS cwd, not ours, so it joins a different root and the process tree splits with no error. The stamp must name a path that survives the change of cwd", b, name, v)
			}
		}
	}
}
