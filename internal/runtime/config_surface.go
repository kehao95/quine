package runtime

// config_surface.go materializes the agent-root config/ read surface (the
// agent's inspectable capability position) and owns the bootstrap half of
// the staged-config transaction (consume/clear + version-skew drift log).
//
// Design authority:
//   Paper/theory/views/runtime-capability/registry-design-brief.md (§ B, § C)
// Work orders:
//   Paper/_design/migrations/runtime-capability-registry-execution.md
//   (T2.1 read surface, T3.2 consume/clear)
//
// Layout under agent/<session>/ :
//
//	config/registry.json — per-build static projection of the compiled
//	  capability registry; raw file, 0444 (the source-code/ read-only
//	  materialization vehicle — FUSE would be static-content misuse).
//	  Refreshed whenever content drifts from the compiled table.
//	config/resolved.env — runtime-owned cache-file of the current resolved
//	  capability position in env syntax (atomic-rename writes). NOT
//	  read-only: the runtime rewrites it on in-process config mutation
//	  (post-D9 the only such mutation is the workspace revision switch).
//	config/next.env — agent-writable staged capability overrides for the
//	  NEXT incarnation; validated+merged at the exec boundary (T3.1,
//	  internal/config/staged.go + internal/tools/exec.go), then archived
//	  and cleared here at the successor's bootstrap (T3.2).
//	inc/<n>/config.env — immutable per-incarnation birth snapshot of the
//	  resolved position; written once at bootstrap, 0444, never rewritten.
//	  Lives under the retained root (agent-root inc/ is a symlink), so the
//	  lineage's capability trajectory survives agent-root cleanup.
//	inc/<n>/config-applied.env — immutable archive of the staged overrides
//	  incarnation n wrote and its successor applied; written by the
//	  successor when it consumes next.env.
//
// Raw files only — the surface works everywhere, including non-Linux hosts
// and degraded-FUSE habitats. The peer-facing FUSE projection is Phase 2.2.

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tools"
)

func (r *Runtime) syncConfigSurface(agentRoot string) error {
	configDir := filepath.Join(agentRoot, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("mkdir config surface: %w", err)
	}
	if err := r.syncConfigRegistrySurface(configDir); err != nil {
		return err
	}
	if err := r.syncResolvedConfigSurface(configDir); err != nil {
		return err
	}
	return r.writeIncarnationConfigSnapshot()
}

func (r *Runtime) syncConfigRegistrySurface(configDir string) error {
	want, err := config.RegistryJSON()
	if err != nil {
		return fmt.Errorf("render capability registry: %w", err)
	}
	want = append(want, '\n')
	path := filepath.Join(configDir, "registry.json")
	got, readErr := os.ReadFile(path)
	if readErr == nil && bytes.Equal(got, want) {
		return nil
	}
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		return fmt.Errorf("read config registry surface: %w", readErr)
	}
	// writeReadOnlyFile stages a 0644 temp file and renames it over any stale
	// (possibly read-only) projection — rename replaces a 0444 target as long
	// as the directory is writable — then the read-only bit is restored.
	return writeReadOnlyFile(path, want, "config registry surface")
}

func (r *Runtime) syncResolvedConfigSurface(configDir string) error {
	return writeFile(filepath.Join(configDir, "resolved.env"), r.cfg.ResolvedEnv())
}

func (r *Runtime) writeIncarnationConfigSnapshot() error {
	path := filepath.Join(r.currentIncarnationRoot(), "config.env")
	if _, err := os.Lstat(path); err == nil {
		// Immutable birth snapshot: write once, never rewrite.
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat incarnation config snapshot: %w", err)
	}
	return writeReadOnlyFile(path, r.cfg.ResolvedEnv(), "incarnation config snapshot")
}

// onConfigMutated re-renders config/resolved.env after an in-process config
// mutation. Post-D9 the only in-process mutation is the workspace revision
// switch (syncWorldRevisionSurface). registry.json is per-build static and
// the inc/<n>/config.env birth snapshot is immutable history, so neither is
// touched here.
func (r *Runtime) onConfigMutated() {
	if !r.agentRootBootstrapped {
		return
	}
	configDir := filepath.Join(r.cfg.AgentRoot(), "config")
	if err := r.syncResolvedConfigSurface(configDir); err != nil {
		r.log("config surface resolved.env sync error: %v", err)
	}
}

// stagedConfigAppliedName is the immutable applied-overrides archive written
// into the staging incarnation's inc/<n>/ dir when the successor consumes
// config/next.env.
const stagedConfigAppliedName = "config-applied.env"

// consumeAppliedStagedConfig is phase 2 of the staged-config transaction
// (registry-design-brief § C; work order T3.2). When this bootstrap is a
// staged exec reentry (execReentryStaged: the QUINE_CONTEXT_BOOTSTRAP
// handover signal, still set here — importBootstrappedContext unsets it only
// after bootstrapAgentRoot) and config/next.env exists, its content is
// already baked into THIS process's envp: the predecessor validated and
// merged it pre-syscall.Exec (T3.1). All that remains is bookkeeping:
// archive the file into the PREDECESSOR's incarnation dir — inc/<n-1>, the
// incarnation that staged it — as the immutable config-applied.env, then
// remove config/next.env so the staging surface is empty for this
// incarnation's own future.
//
// If the bootstrap is NOT a staged reentry (fresh start, session resume,
// fork/spawn child), any next.env is untouched: it belongs to this
// incarnation's own future exec.
//
// Ordering: called from bootstrapAgentRoot BEFORE syncConfigSurface, for two
// reasons. (1) Transaction-before-declaration: the config/ surface and the
// inc/<n>/config.env birth snapshot are born coherent — the snapshot records
// the applied position (from this process's envp) with no consumed-but-
// present next.env sitting next to it, and an archive failure aborts
// bootstrap before any surface claims a clean position. (2) It must run
// after ensureIncarnation (needs the incarnation id to address inc/<n-1>),
// which bootstrapAgentRoot guarantees.
func (r *Runtime) consumeAppliedStagedConfig() error {
	if !execReentryStaged() {
		return nil
	}
	stagedPath := r.cfg.StagedNextEnvPath()
	content, err := os.ReadFile(stagedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read applied staged config %s: %w", stagedPath, err)
	}

	archiveID := r.cfg.IncarnationID - 1
	if archiveID < 0 {
		// A staged exec reentry into a lineage with retained state always
		// advances past the predecessor's id, so id 0 here means the lineage
		// state is off (e.g. the retained root vanished between exec and
		// bootstrap). Keep the evidence anyway, under the current dir.
		r.log("staged config: exec reentry at incarnation id %d has no predecessor dir; archiving under the current incarnation", r.cfg.IncarnationID)
		archiveID = r.cfg.IncarnationID
	}
	archiveDir := r.cfg.SessionIncarnationPath("", archiveID)
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return fmt.Errorf("mkdir staged-config archive dir: %w", err)
	}
	archivePath := filepath.Join(archiveDir, stagedConfigAppliedName)
	if err := writeReadOnlyFile(archivePath, content, "staged-config archive"); err != nil {
		return err
	}
	if err := os.Remove(stagedPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("clear consumed staged config %s: %w", stagedPath, err)
	}
	return nil
}

// nonRegistryRuntimeEnv reports whether an inherited QUINE_* env name is
// legitimately outside the compiled capability registry: runtime-emitted
// process-surface signals and runtime-internal child plumbing that were
// never registry knobs. Authoritative sets: internal/config/registry.go (the
// registry proper, via config.KnobByEnv) plus the runtime-emitted family
// recorded in scripts/check-authored-env-consistency.sh's allowlist (T1.3).
//
// This exclusion list is a stopgap for the T1.3 "registry completeness gap"
// follow-up (runtime-emitted envs without envnames.go constants:
// QUINE_AGENT_ROOT, QUINE_CONTEXT_BOOTSTRAP, QUINE_JOB_*,
// QUINE_WORKSPACE_ENABLED, QUINE_TEST_*): when those gain registry entries,
// shrink this list in the same change.
func nonRegistryRuntimeEnv(name string) bool {
	switch name {
	case "QUINE_AGENT_ROOT", // exec.go:execProcessSurfaceEnv emits it into the agent env
		tools.ContextBootstrapEnv, // exec/fork handover signal (QUINE_CONTEXT_BOOTSTRAP)
		"QUINE_WORKSPACE_ENABLED": // overlay backend emits it into sh-child envs
		return true
	}
	for _, prefix := range []string{
		"QUINE_JOB_",  // sh job-wrapper plumbing (interactive.go/tools.go)
		"QUINE_TEST_", // deliberate test-only runtime hooks
	} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// logCapabilityEnvDrift records inherited QUINE_* env names the CURRENT
// binary's registry does not know (registry-design-brief § C version-skew
// paragraph; work order T3.2). This is the observability half of the
// version-skew stance: a self-upgraded successor whose registry renamed or
// removed a knob silently ignores the inherited name in Load() by design —
// self-upgrade must not deadlock — and this bootstrap log line is what keeps
// that silence diagnosable in the lineage record (D4).
func (r *Runtime) logCapabilityEnvDrift() {
	var drifted []string
	for _, kv := range os.Environ() {
		name, _, ok := strings.Cut(kv, "=")
		if !ok || !strings.HasPrefix(name, "QUINE_") {
			continue
		}
		if _, known := config.KnobByEnv(name); known {
			continue
		}
		if nonRegistryRuntimeEnv(name) {
			continue
		}
		drifted = append(drifted, name)
	}
	if len(drifted) == 0 {
		return
	}
	sort.Strings(drifted)
	r.log("config drift: inherited QUINE_* env names unknown to this binary's capability registry (knob renamed/removed by a self-upgrade? Load() ignores them silently): %s", strings.Join(drifted, ", "))
}
