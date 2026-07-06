package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// --- grammar ---

func TestParseStagedEnvAcceptsEnvFileGrammar(t *testing.T) {
	content := strings.Join([]string{
		"# staged for the next incarnation",
		"",
		"QUINE_MAX_TURNS=12",
		"   # indented comment",
		"QUINE_PROMPT_METAPHOR=thermodynamic\r", // CRLF tolerated
		"QUINE_INITIAL_USER_MESSAGE=note with spaces and = signs",
		"QUINE_FORK_ENABLED=0",
	}, "\n")

	staged, err := ParseStagedEnv([]byte(content))
	if err != nil {
		t.Fatalf("ParseStagedEnv() error: %v", err)
	}
	want := map[string]string{
		"QUINE_MAX_TURNS":            "12",
		"QUINE_PROMPT_METAPHOR":      "thermodynamic",
		"QUINE_INITIAL_USER_MESSAGE": "note with spaces and = signs",
		"QUINE_FORK_ENABLED":         "0",
	}
	if !reflect.DeepEqual(staged, want) {
		t.Fatalf("ParseStagedEnv() = %#v, want %#v", staged, want)
	}
}

func TestParseStagedEnvEmptyAndCommentOnlyContentIsNoOp(t *testing.T) {
	for _, content := range []string{"", "\n\n", "# only a comment\n\n  \n"} {
		staged, err := ParseStagedEnv([]byte(content))
		if err != nil {
			t.Fatalf("ParseStagedEnv(%q) error: %v", content, err)
		}
		if len(staged) != 0 {
			t.Fatalf("ParseStagedEnv(%q) = %#v, want empty", content, staged)
		}
	}
}

func TestParseStagedEnvRejectsSyntaxViolations(t *testing.T) {
	cases := []struct {
		name    string
		line    string
		errWant string
	}{
		{"no assignment", "just some words", "not a KEY=VALUE assignment"},
		{"leading whitespace", "  QUINE_MAX_TURNS=5", "not a legal env name"},
		{"whitespace before equals", "QUINE_MAX_TURNS =5", "not a legal env name"},
		{"empty key", "=5", "not a legal env name"},
		{"illegal key char", "QUINE-MAX-TURNS=5", "not a legal env name"},
		{"export prefix", "export QUINE_MAX_TURNS=5", "not a legal env name"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			staged, err := ParseStagedEnv([]byte(tc.line + "\n"))
			if err == nil {
				t.Fatalf("ParseStagedEnv(%q) accepted, want syntax rejection (got %#v)", tc.line, staged)
			}
			if staged != nil {
				t.Fatalf("rejected file must return a nil map, got %#v", staged)
			}
			if !strings.Contains(err.Error(), tc.errWant) {
				t.Fatalf("error %q does not name the violation %q", err, tc.errWant)
			}
			if !strings.Contains(err.Error(), "line 1") {
				t.Fatalf("error %q does not name the offending line", err)
			}
		})
	}
}

func TestParseStagedEnvRejectsDuplicateAssignment(t *testing.T) {
	content := "QUINE_MAX_TURNS=5\nQUINE_MAX_TURNS=9\n"
	_, err := ParseStagedEnv([]byte(content))
	if err == nil {
		t.Fatal("duplicate assignment accepted, want rejection")
	}
	if !strings.Contains(err.Error(), "duplicate assignment for QUINE_MAX_TURNS") ||
		!strings.Contains(err.Error(), "line 1") {
		t.Fatalf("duplicate error should name the knob and the first line, got %q", err)
	}
}

// --- registry validation ---

func TestParseStagedEnvRejectsUnknownName(t *testing.T) {
	_, err := ParseStagedEnv([]byte("QUINE_NOT_A_KNOB=1\n"))
	if err == nil {
		t.Fatal("unknown env name accepted, want rejection")
	}
	if !strings.Contains(err.Error(), "QUINE_NOT_A_KNOB") ||
		!strings.Contains(err.Error(), "unknown env name") {
		t.Fatalf("unknown-name error should name the knob and the reason, got %q", err)
	}
}

func TestParseStagedEnvRejectsRetiredNameWithReplacementPointer(t *testing.T) {
	_, err := ParseStagedEnv([]byte("QUINE_SH_TIMEOUT=600\n"))
	if err == nil {
		t.Fatal("retired env name accepted, want rejection")
	}
	if !strings.Contains(err.Error(), "retired knob") ||
		!strings.Contains(err.Error(), EnvShDefaultTimeout) {
		t.Fatalf("retired-name error should point at the replacement %s, got %q", EnvShDefaultTimeout, err)
	}
}

func TestParseStagedEnvRejectsTypeViolations(t *testing.T) {
	cases := []struct {
		name    string
		line    string
		errWant string
	}{
		{"bool", "QUINE_FORK_ENABLED=yes", `invalid bool value "yes"`},
		{"int", "QUINE_MAX_TURNS=soon", `invalid int value "soon"`},
		{"enum", "QUINE_SH_NETWORK=wifi", `invalid enum value "wifi" (legal: host|none)`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseStagedEnv([]byte(tc.line + "\n"))
			if err == nil {
				t.Fatalf("ParseStagedEnv(%q) accepted, want type rejection", tc.line)
			}
			if !strings.Contains(err.Error(), tc.errWant) {
				t.Fatalf("error %q does not carry the legal form %q", err, tc.errWant)
			}
		})
	}
}

func TestParseStagedEnvMutabilityWhitelistOnlyExecBoundary(t *testing.T) {
	cases := []struct {
		name string
		line string
		mut  Mutability
	}{
		{"substrate-pinned", "QUINE_DEPTH=5", MutSubstratePinned},
		{"runtime-emitted", "QUINE_SESSION_ID=faked-session", MutRuntimeEmitted},
		{"operator-only auth", "QUINE_API_KEY=stolen", MutOperatorOnly},
		{"operator-only path", "QUINE_DATA_DIR=/tmp/elsewhere", MutOperatorOnly},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseStagedEnv([]byte(tc.line + "\n"))
			if err == nil {
				t.Fatalf("ParseStagedEnv(%q) accepted, want mutability rejection", tc.line)
			}
			if !strings.Contains(err.Error(), string(tc.mut)) ||
				!strings.Contains(err.Error(), "only exec-boundary knobs") {
				t.Fatalf("mutability error should name the class %q and the whitelist, got %q", tc.mut, err)
			}
		})
	}

	// The whitelist itself: an exec-boundary knob passes.
	staged, err := ParseStagedEnv([]byte("QUINE_OUTPUT_TRUNCATE=31337\n"))
	if err != nil {
		t.Fatalf("exec-boundary knob rejected: %v", err)
	}
	if staged["QUINE_OUTPUT_TRUNCATE"] != "31337" {
		t.Fatalf("staged = %#v, want QUINE_OUTPUT_TRUNCATE=31337", staged)
	}
}

func TestParseStagedEnvRegistryAgreesOnMutabilityClasses(t *testing.T) {
	// Every knob used by the whitelist tests above must keep the mutability
	// class the test assumes; if the registry re-classifies one, fail here
	// with a clear pointer instead of silently weakening the tests.
	for env, want := range map[string]Mutability{
		"QUINE_DEPTH":           MutSubstratePinned,
		"QUINE_SESSION_ID":      MutRuntimeEmitted,
		"QUINE_API_KEY":         MutOperatorOnly,
		"QUINE_DATA_DIR":        MutOperatorOnly,
		"QUINE_OUTPUT_TRUNCATE": MutExecBoundary,
	} {
		knob, ok := KnobByEnv(env)
		if !ok {
			t.Fatalf("registry lost %s; update staged_test.go fixtures", env)
		}
		if knob.Mutability != want {
			t.Fatalf("%s mutability = %q, test fixtures assume %q; update staged_test.go", env, knob.Mutability, want)
		}
	}
}

func TestParseStagedEnvWholeFileRejectEnumeratesAllViolations(t *testing.T) {
	content := strings.Join([]string{
		"QUINE_MAX_TURNS=12",     // valid — must still not be applied
		"QUINE_NOT_A_KNOB=1",     // unknown
		"QUINE_FORK_ENABLED=yes", // type
		"QUINE_API_KEY=stolen",   // mutability
		"no assignment here",     // syntax
	}, "\n")

	staged, err := ParseStagedEnv([]byte(content))
	if err == nil {
		t.Fatal("invalid file accepted, want whole-file rejection")
	}
	if staged != nil {
		t.Fatalf("whole-file reject must strip nothing and apply nothing; got %#v", staged)
	}
	var sfe *StagedFileError
	ok := false
	if sfe, ok = err.(*StagedFileError); !ok {
		t.Fatalf("error type = %T, want *StagedFileError", err)
	}
	if len(sfe.Violations) != 4 {
		t.Fatalf("violations = %d, want 4 (every violation enumerated):\n%v", len(sfe.Violations), err)
	}
	for _, want := range []string{
		"QUINE_NOT_A_KNOB",
		"QUINE_FORK_ENABLED",
		"QUINE_API_KEY",
		"not a KEY=VALUE assignment",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("enumerated error missing %q:\n%v", want, err)
		}
	}
}

func TestParseStagedEnvEmptyValueIsResetToDefault(t *testing.T) {
	// KEY= is the explicit reset-to-default idiom: it must pass every type
	// check and merge as an empty value, which Load() resolves to the default.
	staged, err := ParseStagedEnv([]byte("QUINE_MAX_TURNS=\nQUINE_FORK_ENABLED=\nQUINE_SH_NETWORK=\n"))
	if err != nil {
		t.Fatalf("empty values rejected: %v", err)
	}
	for _, env := range []string{"QUINE_MAX_TURNS", "QUINE_FORK_ENABLED", "QUINE_SH_NETWORK"} {
		if v, ok := staged[env]; !ok || v != "" {
			t.Fatalf("staged[%s] = %q,%v, want empty value present", env, v, ok)
		}
	}
}

// --- merge ---

func TestMergeStagedOverridesReplacesAndAppends(t *testing.T) {
	env := []string{
		"QUINE_MODEL_ID=model-a",
		"QUINE_MAX_TURNS=40",
		"QUINE_OUTPUT_TRUNCATE=20480",
	}
	staged := map[string]string{
		"QUINE_OUTPUT_TRUNCATE":         "31337", // replace in place
		"QUINE_WALL_CLOCK_EXIT_SECONDS": "90",    // append (no serialized line)
		"QUINE_ANCHOR_MEMORY":           "1",     // append, sorted before WALL_CLOCK
	}

	got := MergeStagedOverrides(env, staged)
	want := []string{
		"QUINE_MODEL_ID=model-a",
		"QUINE_MAX_TURNS=40",
		"QUINE_OUTPUT_TRUNCATE=31337",
		"QUINE_ANCHOR_MEMORY=1",
		"QUINE_WALL_CLOCK_EXIT_SECONDS=90",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeStagedOverrides() = %#v, want %#v", got, want)
	}
}

func TestMergeStagedOverridesEmptyStagedIsIdentity(t *testing.T) {
	env := []string{"QUINE_MODEL_ID=model-a"}
	if got := MergeStagedOverrides(env, nil); !reflect.DeepEqual(got, env) {
		t.Fatalf("MergeStagedOverrides(env, nil) = %#v, want unchanged input", got)
	}
	if got := MergeStagedOverrides(env, map[string]string{}); !reflect.DeepEqual(got, env) {
		t.Fatalf("MergeStagedOverrides(env, {}) = %#v, want unchanged input", got)
	}
}

func TestMergeStagedOverridesEmptyValueReplacesEntry(t *testing.T) {
	env := []string{"QUINE_MAX_TURNS=40"}
	got := MergeStagedOverrides(env, map[string]string{"QUINE_MAX_TURNS": ""})
	want := []string{"QUINE_MAX_TURNS="}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MergeStagedOverrides() = %#v, want %#v (reset-to-default)", got, want)
	}
}

// --- file-level reader ---

func TestReadStagedOverridesAbsentFileIsNoOp(t *testing.T) {
	staged, err := ReadStagedOverrides(filepath.Join(t.TempDir(), "config", "next.env"))
	if err != nil {
		t.Fatalf("absent staged file should be a clean no-op, got %v", err)
	}
	if staged != nil {
		t.Fatalf("absent staged file should yield nil overrides, got %#v", staged)
	}
	if staged, err = ReadStagedOverrides(""); err != nil || staged != nil {
		t.Fatalf("empty path should be a clean no-op, got %#v, %v", staged, err)
	}
}

func TestReadStagedOverridesKeepsInvalidFileIntact(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "next.env")
	content := []byte("QUINE_NOT_A_KNOB=1\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write staged file: %v", err)
	}

	staged, err := ReadStagedOverrides(path)
	if err == nil {
		t.Fatalf("invalid staged file accepted: %#v", staged)
	}
	if !strings.Contains(err.Error(), "left intact") {
		t.Fatalf("rejection should tell the agent the file is intact for fix-and-retry, got %q", err)
	}
	after, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("re-read staged file: %v", readErr)
	}
	if string(after) != string(content) {
		t.Fatalf("rejection must not modify the staged file: %q -> %q", content, after)
	}
}

func TestStagedNextEnvPath(t *testing.T) {
	cfg := &Config{
		Identity: Identity{SessionID: "staged-path-session"},
		Paths:    Paths{DataDir: t.TempDir()},
	}
	want := filepath.Join(cfg.AgentRoot(), "config", "next.env")
	if got := cfg.StagedNextEnvPath(); got != want {
		t.Fatalf("StagedNextEnvPath() = %q, want %q", got, want)
	}
}
