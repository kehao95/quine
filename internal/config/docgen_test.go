package config

// docgen_test.go pins the generated-projection contract (work order T1.2):
//
//  1. Content equality: the committed env-controls.md and .env.example
//     byte-match the rendered output. This is the doc-freshness gate for
//     generated files (replaces presence-based gating; brief Phase 1).
//  2. Structural sanity: the layout covers the registry exactly once, the
//     retired table never shadows a live knob, and no retired name leaks
//     into .env.example (the check-authored-env-consistency contract).
//
// Embedded-tree caveat: TestEmbeddedSelfSourceBuildsAndTests extracts
// internal/+cmd/+go.mod into a temp dir and runs `go test ./...` there.
// Paper/ and .env.example do not exist in that extraction, so the equality
// tests locate the repo root by walking up for the env-controls.md marker
// and t.Skip when it is absent. The extraction lands under the system temp
// dir, never inside the repo, so the walk cannot false-positive onto the
// real repo root.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// docRepoRoot walks up from the test working directory looking for the
// generated-doc marker. Returns ok=false outside the full repo checkout
// (e.g. inside the embedded self-source extraction).
func docRepoRoot(t *testing.T) (string, bool) {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		marker := filepath.Join(dir, "Paper", "core", "registries", "env-controls.md")
		if _, err := os.Stat(marker); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func firstDiffLine(got, want string) (int, string, string) {
	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")
	n := len(gotLines)
	if len(wantLines) < n {
		n = len(wantLines)
	}
	for i := 0; i < n; i++ {
		if gotLines[i] != wantLines[i] {
			return i + 1, gotLines[i], wantLines[i]
		}
	}
	return n + 1, "", ""
}

func assertRenderedMatchesFile(t *testing.T, rendered, rel string) {
	t.Helper()
	root, ok := docRepoRoot(t)
	if !ok {
		t.Skipf("%s not reachable above cwd (embedded self-source extraction); content-equality runs only in the full repo", rel)
	}
	path := filepath.Join(root, filepath.FromSlash(rel))
	onDisk, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if string(onDisk) == rendered {
		return
	}
	line, got, want := firstDiffLine(string(onDisk), rendered)
	t.Fatalf("%s is stale relative to the registry renderer (first difference at line %d)\n  on disk:  %q\n  rendered: %q\nregenerate with: go run ./scripts/gen-env-docs", rel, line, got, want)
}

func TestGeneratedEnvControlsDocIsFresh(t *testing.T) {
	rendered, err := RenderEnvControlsDoc()
	if err != nil {
		t.Fatalf("RenderEnvControlsDoc: %v", err)
	}
	assertRenderedMatchesFile(t, rendered, "Paper/core/registries/env-controls.md")
}

func TestGeneratedEnvExampleIsFresh(t *testing.T) {
	rendered, err := RenderEnvExample()
	if err != nil {
		t.Fatalf("RenderEnvExample: %v", err)
	}
	assertRenderedMatchesFile(t, rendered, ".env.example")
}

// TestEnvControlsLayoutCoversRegistry runs everywhere (including the embedded
// tree): the renderer's own completeness checks must hold — every live knob
// has exactly one row, every retired/external entry is placed, and rendering
// is deterministic.
func TestEnvControlsLayoutCoversRegistry(t *testing.T) {
	first, err := RenderEnvControlsDoc()
	if err != nil {
		t.Fatalf("RenderEnvControlsDoc: %v", err)
	}
	second, err := RenderEnvControlsDoc()
	if err != nil {
		t.Fatalf("RenderEnvControlsDoc (second render): %v", err)
	}
	if first != second {
		t.Fatal("RenderEnvControlsDoc is not deterministic")
	}
	for _, k := range Registry {
		if !strings.Contains(first, "| `"+k.Env+"` |") {
			t.Errorf("registry knob %s (%s) missing from rendered env-controls.md", k.Name, k.Env)
		}
	}
	for _, k := range RetiredRegistry {
		if !strings.Contains(first, "| `"+k.Env+"` | retired |") {
			t.Errorf("retired knob %s missing from rendered env-controls.md", k.Env)
		}
	}
}

// TestRetiredRegistryDisjointFromLiveRegistry guards the tombstone trap: the
// load-error legacy knobs (SELF_REENTRY_TARGET, WORKSPACE_SOURCE) live in the
// registry proper and must never be duplicated as retired rows.
func TestRetiredRegistryDisjointFromLiveRegistry(t *testing.T) {
	seen := map[string]bool{}
	for _, k := range RetiredRegistry {
		if seen[k.Env] {
			t.Errorf("retired knob %s listed twice", k.Env)
		}
		seen[k.Env] = true
		if _, ok := KnobByEnv(k.Env); ok {
			t.Errorf("retired knob %s is also a live registry knob", k.Env)
		}
		if k.Decision == "" || k.When == "" || k.Note == "" {
			t.Errorf("retired knob %s missing Decision/When/Note", k.Env)
		}
		if k.Replacement != "" {
			if _, ok := KnobByEnv(k.Replacement); !ok {
				t.Errorf("retired knob %s names replacement %s which is not a live registry knob", k.Env, k.Replacement)
			}
		}
	}
	for _, l := range ExternalLabels {
		if _, ok := KnobByEnv(l.Env); ok {
			t.Errorf("external label %s is also a live registry knob", l.Env)
		}
	}
}

// TestEnvExampleContainsOnlyLiveNames: retired names must never reappear in
// .env.example (check-authored-env-consistency scans .env* for unknown
// names), and non-operator knobs (runtime-emitted identity, substrate-pinned
// lineage, legacy tombstones) stay out.
func TestEnvExampleContainsOnlyLiveNames(t *testing.T) {
	rendered, err := RenderEnvExample()
	if err != nil {
		t.Fatalf("RenderEnvExample: %v", err)
	}
	for _, k := range RetiredRegistry {
		// Token-boundary match (the check-authored-env-consistency convention):
		// QUINE_SH_TIMEOUT must not flag the live QUINE_SH_TIMEOUT_OVERRIDE_ENABLED.
		name := strings.TrimSuffix(k.Env, "*")
		re := regexp.MustCompile(`(^|[^A-Za-z0-9_])` + regexp.QuoteMeta(name) + `([^A-Z0-9_]|$)`)
		if strings.HasSuffix(k.Env, "*") {
			re = regexp.MustCompile(`(^|[^A-Za-z0-9_])` + regexp.QuoteMeta(name))
		}
		if re.MatchString(rendered) {
			t.Errorf("retired env %s appears in rendered .env.example", k.Env)
		}
	}
	for _, k := range Registry {
		listed := strings.Contains(rendered, "export "+k.Env+"=")
		operator := envExampleOperatorSettable(k) || k.Default.Kind == DefaultRequired
		if operator && !listed {
			t.Errorf("operator-settable knob %s missing from rendered .env.example", k.Env)
		}
		if !operator && listed {
			t.Errorf("non-operator knob %s (%s/%s) must not appear in rendered .env.example", k.Env, k.Mutability, k.Default.Kind)
		}
	}
}
