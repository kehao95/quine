package runtime

// config_surface_test.go covers the agent-root config/ surface under the
// process-native env boundary model:
//
//	config/registry.json          the knob catalog — what an ABSENT name means
//	config/env/override            the ONE managed env file: the agent-authored
//	                               child-env policy (agent-writable, never
//	                               created by the runtime)
//	inc/<n>/override-applied.env  the immutable, VERBATIM exec-time archive of
//	                               an applied override
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//
// These tests are deliberately build-tag-free: they are the ONLY coverage of
// this surface on non-Linux / degraded-FUSE hosts, where the raw files are the
// whole mechanism and the peer projection does not exist (brief E7).
//
// There is no config/env/pinned, no config/env/effective, and no inc/<n>/environ
// birth snapshot, and the absence is the point. The runtime used to render a
// process's own environment back to it out of resolved Config values — which is
// how `QUINE_SPAWN_ENABLED=0`, a line no operator ever authored, taught a
// founder that constructing a successor was impossible. The OS already
// publishes a process's own environment, per-PID and immutably, at
// /proc/<pid>/environ; a runtime-owned re-rendering could only add content
// nobody wrote or drop content the OS already shows in full.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tools"
)

func bootstrapConfigSurfaceRuntime(t *testing.T, cfg *config.Config, mission string) *Runtime {
	t.Helper()
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = mission
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	return rt
}

// clearEnvForTest removes name from the process env for the duration of the
// test, restoring whatever was there afterwards. Absence is a load-bearing
// assertion in this file, so a test that claims "this name is not in the child
// env" must first guarantee it is not in the parent's either.
func clearEnvForTest(t *testing.T, name string) {
	t.Helper()
	t.Setenv(name, "") // registers the restore-to-original cleanup
	if err := os.Unsetenv(name); err != nil {
		t.Fatalf("unset %s: %v", name, err)
	}
}

// writeEnvOverrideFile writes config/env/override the way the agent does: a raw
// shell write, straight to the file, with no gate in the path.
func writeEnvOverrideFile(t *testing.T, cfg *config.Config, content string) string {
	t.Helper()
	path := cfg.EnvOverridePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir config/env: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config/env/override: %v", err)
	}
	return path
}

// TestBootstrapAgentRootMaterializesConfigSurface asserts the whole surface
// bootstrap materializes: config/registry.json (the catalog, byte-exact and
// read-only) and an env/ directory the agent can write its override into —
// and, just as load-bearing, that the deleted surfaces (config/env/pinned,
// config/env/effective, inc/0/environ) do NOT reappear.
func TestBootstrapAgentRootMaterializesConfigSurface(t *testing.T) {
	t.Setenv(config.EnvAPIKey, "sk-live-do-not-leak")
	t.Setenv("QUINE_MAX_TURNS", "7") // operator-authored: irrelevant now, but harmless
	cfg := testCfg(t)
	bootstrapConfigSurfaceRuntime(t, cfg, "materialize config surface")

	configDir := filepath.Join(cfg.AgentRoot(), "config")

	// registry.json: the catalog that gives absence its meaning.
	regPath := filepath.Join(configDir, "registry.json")
	regData, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read config/registry.json: %v", err)
	}
	wantReg, err := config.RegistryJSON()
	if err != nil {
		t.Fatalf("render compiled registry: %v", err)
	}
	wantReg = append(wantReg, '\n')
	if !bytes.Equal(regData, wantReg) {
		t.Fatalf("config/registry.json drifts from the compiled RegistryJSON payload")
	}
	assertPerm(t, regPath, 0o444)

	// config/env/ must exist so the agent's raw shell write to
	// config/env/override just works, with no directory to create first.
	envDir := filepath.Join(configDir, config.EnvOverrideDirName)
	info, err := os.Stat(envDir)
	if err != nil {
		t.Fatalf("stat config/env: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("config/env is not a directory")
	}

	// The deleted surfaces must not reappear.
	for _, path := range []string{
		filepath.Join(envDir, "pinned"),
		filepath.Join(envDir, "effective"),
		filepath.Join(cfg.AgentRoot(), "inc", "0", "environ"),
	} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("%s must not exist: the runtime no longer renders a process's own environment back to it (stat err = %v)", path, err)
		}
	}

	// The runtime never authors the override: no policy is the normal state, and
	// an empty file the agent did not write would be a line nobody chose.
	if _, err := os.Lstat(cfg.EnvOverridePath()); !os.IsNotExist(err) {
		t.Fatalf("config/env/override must not be created by the runtime (stat err = %v)", err)
	}
}

// TestBootstrapAgentRootRefreshesConfigSurfaceOnDrift: config/registry.json is
// tamper-EVIDENT, not tamper-proof — the envp constructors never read it
// (brief Layer 1 § 4), so defacing it changes nothing about what a child
// receives — but the next bootstrap repairs it back to the compiled content.
func TestBootstrapAgentRootRefreshesConfigSurfaceOnDrift(t *testing.T) {
	cfg := testCfg(t)
	rt := bootstrapConfigSurfaceRuntime(t, cfg, "refresh drifted config surface")

	configDir := filepath.Join(cfg.AgentRoot(), "config")
	regPath := filepath.Join(configDir, "registry.json")

	if err := os.Chmod(regPath, 0o644); err != nil {
		t.Fatalf("chmod %s writable: %v", regPath, err)
	}
	if err := os.WriteFile(regPath, []byte("QUINE_STALE=1\n"), 0o644); err != nil {
		t.Fatalf("deface %s: %v", regPath, err)
	}

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("second bootstrapAgentRoot failed: %v", err)
	}

	wantReg, err := config.RegistryJSON()
	if err != nil {
		t.Fatalf("render compiled registry: %v", err)
	}
	wantReg = append(wantReg, '\n')

	got, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read repaired %s: %v", regPath, err)
	}
	if strings.Contains(string(got), "QUINE_STALE") {
		t.Fatalf("%s was not repaired after defacement", regPath)
	}
	if !bytes.Equal(got, wantReg) {
		t.Fatalf("%s was not restored to the compiled content", regPath)
	}
	assertPerm(t, regPath, 0o444)
}

// TestBootstrapArchivesAndClearsEnvOverrideOnExecReentry: the exec-time
// archive of an applied override is now written VERBATIM — the raw bytes of
// config/env/override, unrendered and unredacted. The override is the agent's
// OWN authored file (the parser already rejects pinned/runtime-emitted
// registry knobs, so it can never carry a registry credential), and a foreign
// line the agent chose to set (E9) is the agent's own authorship: redacting it
// here would make the lineage record lie about what was actually applied.
func TestBootstrapArchivesAndClearsEnvOverrideOnExecReentry(t *testing.T) {
	t.Setenv(tools.ContextBootstrapEnv, "")
	cfg := testCfg(t)

	// Predecessor incarnation (inc/0) bootstraps and authors a policy for the
	// successor it is about to exec into. The policy carries a legal foreign
	// line with a credential-shaped value: foreign names are free-form by
	// design (E9), so nothing rejects GITHUB_TOKEN=…, and it must survive the
	// archive byte-for-byte — verbatim means verbatim.
	bootstrapConfigSurfaceRuntime(t, cfg, "predecessor stages an env policy")
	policy := "# staged for my successor\nQUINE_OUTPUT_TRUNCATE=31337\nFOO=bar\nGITHUB_TOKEN=ghp_SECRET_NEVER_RETAINED\n"
	overridePath := writeEnvOverrideFile(t, cfg, policy)

	// Successor: exec handover signaled by the staged bootstrap context. In the
	// real flow the policy is already baked into this process's envp (the
	// predecessor applied it pre-syscall.Exec); the bootstrap step is pure
	// bookkeeping — archive it, then clear it.
	t.Setenv(tools.ContextBootstrapEnv, t.TempDir())
	successor := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(successor)
	successor.originalInput = "successor consumes the policy"
	if err := successor.bootstrapAgentRoot(); err != nil {
		t.Fatalf("successor bootstrapAgentRoot failed: %v", err)
	}

	if cfg.IncarnationID != 1 {
		t.Fatalf("successor incarnation id = %d, want 1", cfg.IncarnationID)
	}
	if _, err := os.Lstat(overridePath); !os.IsNotExist(err) {
		t.Fatalf("config/env/override should be cleared after a staged reentry (stat err = %v)", err)
	}

	// The archive lands in the PREDECESSOR's incarnation dir: inc/0 is the
	// incarnation that authored the policy, and the record belongs to it.
	archivePath := filepath.Join(cfg.SessionIncarnationPath("", 0), overrideAppliedName)
	archived, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read inc/0/%s archive: %v", overrideAppliedName, err)
	}
	// Verbatim: byte-exact against the authored policy, comment and credential
	// line included — no rendering, no redaction.
	if string(archived) != policy {
		t.Fatalf("archive should be the VERBATIM authored policy:\n got:  %q\n want: %q", archived, policy)
	}
	assertPerm(t, archivePath, 0o444)
	// Reachable through the agent-root inc/ symlink too.
	if _, err := os.ReadFile(filepath.Join(cfg.AgentRoot(), "inc", "0", overrideAppliedName)); err != nil {
		t.Fatalf("archive not visible through the agent-root inc/ symlink: %v", err)
	}

	// Idempotence: a re-run of the bootstrap (env still set, file gone) is a
	// clean no-op.
	if err := successor.bootstrapAgentRoot(); err != nil {
		t.Fatalf("re-bootstrap after consumption failed: %v", err)
	}
}

func TestBootstrapLeavesEnvOverrideWithoutExecReentry(t *testing.T) {
	t.Setenv(tools.ContextBootstrapEnv, "")
	cfg := testCfg(t)
	rt := bootstrapConfigSurfaceRuntime(t, cfg, "non-staged bootstrap")

	policy := "QUINE_OUTPUT_TRUNCATE=31337\n"
	overridePath := writeEnvOverrideFile(t, cfg, policy)

	// A non-staged re-bootstrap (fresh start / resume shape) must not touch the
	// file: the policy belongs to THIS incarnation's own future children and its
	// own future exec, not to a handover that never happened.
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("re-bootstrap failed: %v", err)
	}

	after, err := os.ReadFile(overridePath)
	if err != nil {
		t.Fatalf("config/env/override should survive a non-staged bootstrap: %v", err)
	}
	if string(after) != policy {
		t.Fatalf("config/env/override content changed: %q -> %q", policy, after)
	}
	if _, err := os.Lstat(filepath.Join(cfg.SessionIncarnationPath("", 0), overrideAppliedName)); !os.IsNotExist(err) {
		t.Fatalf("no %s archive may appear without a staged reentry (stat err = %v)", overrideAppliedName, err)
	}
}

// TestConsumeAppliedEnvOverrideWithoutPredecessorArchivesUnderCurrent: a staged
// reentry at incarnation 0 (no predecessor dir — off-nominal lineage state)
// keeps the evidence under the current incarnation rather than dropping it, and
// the archived content is the VERBATIM authored policy — no header, no
// rendering, no redaction.
func TestConsumeAppliedEnvOverrideWithoutPredecessorArchivesUnderCurrent(t *testing.T) {
	t.Setenv(tools.ContextBootstrapEnv, t.TempDir())
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "orphan staged reentry"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	policy := "QUINE_OUTPUT_TRUNCATE=31337\n"
	overridePath := writeEnvOverrideFile(t, cfg, policy)

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	if cfg.IncarnationID != 0 {
		t.Fatalf("incarnation id = %d, want 0 (no retained state)", cfg.IncarnationID)
	}
	if _, err := os.Lstat(overridePath); !os.IsNotExist(err) {
		t.Fatalf("config/env/override should still be consumed (stat err = %v)", err)
	}
	archived, err := os.ReadFile(filepath.Join(cfg.SessionIncarnationPath("", 0), overrideAppliedName))
	if err != nil {
		t.Fatalf("fallback archive under the current incarnation missing: %v", err)
	}
	if string(archived) != policy {
		t.Fatalf("fallback archive content = %q, want %q", archived, policy)
	}
}

func TestLogCapabilityEnvDrift(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	var logged []string
	rt.log = func(format string, args ...any) {
		logged = append(logged, fmt.Sprintf(format, args...))
	}

	// One genuine drift name (registry-unknown), plus every excluded
	// non-registry runtime namespace — the latter must stay silent.
	t.Setenv("QUINE_BOGUS_RENAMED_KNOB", "1")
	t.Setenv("QUINE_JOB_SHELL", "/bin/sh")
	t.Setenv("QUINE_TEST_HOOK", "1")
	t.Setenv("QUINE_AGENT_ROOT", "/somewhere")
	t.Setenv("QUINE_WORKSPACE_ENABLED", "1")
	t.Setenv(tools.ContextBootstrapEnv, "/staged")

	rt.logCapabilityEnvDrift()

	all := strings.Join(logged, "\n")
	if !strings.Contains(all, "QUINE_BOGUS_RENAMED_KNOB") {
		t.Fatalf("drift log should name the registry-unknown env, got %q", all)
	}
	for _, excluded := range []string{
		"QUINE_JOB_SHELL",
		"QUINE_TEST_HOOK",
		"QUINE_AGENT_ROOT",
		"QUINE_WORKSPACE_ENABLED",
		tools.ContextBootstrapEnv,
	} {
		if strings.Contains(all, excluded) {
			t.Fatalf("drift log must stay silent for excluded namespace %s, got %q", excluded, all)
		}
	}

	// After the unknown name is gone it must not be reported again. (Total
	// silence is not asserted: the ambient test environment may carry
	// legitimate drift such as a profile's QUINE_PROVIDER label.)
	t.Setenv("QUINE_BOGUS_RENAMED_KNOB", "")
	if err := os.Unsetenv("QUINE_BOGUS_RENAMED_KNOB"); err != nil {
		t.Fatalf("unset bogus env: %v", err)
	}
	logged = nil
	rt.logCapabilityEnvDrift()
	if strings.Contains(strings.Join(logged, "\n"), "QUINE_BOGUS_RENAMED_KNOB") {
		t.Fatalf("cleared env still reported as drift: %q", logged)
	}
}

func assertPerm(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if info.Mode().Perm() != want {
		t.Fatalf("%s permissions = %o, want %o", path, info.Mode().Perm(), want)
	}
}
