package runtime

// config_surface.go materializes the agent-root config/ surface and owns the
// bootstrap half of the exec-time env-override transaction (archive/clear +
// version-skew drift log).
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//
// Layout under agent/<session>/ :
//
//	config/registry.json — per-build static projection of the compiled
//	  capability registry: the catalog that says what an ABSENT knob means, and
//	  (via each knob's Mutability) which names config/env/override cannot set.
//	  Raw file, 0444, refreshed whenever content drifts from the compiled table.
//	config/env/override — the ONE managed env file: the agent-authored child-env
//	  policy. Raw file, written by the agent with ordinary shell writes (or the
//	  ctl/env peer gate); validated at every use.
//	inc/<n>/override-applied.env — immutable, verbatim archive of the override
//	  that incarnation n staged and its successor applied at exec. Lives under the
//	  retained root (agent-root inc/ is a symlink), so the lineage's capability
//	  trajectory survives agent-root cleanup.
//
// The runtime does NOT render a process's own environment back to it. There is
// no config/resolved.env, no config/env/effective, no inc/<n>/environ snapshot:
// the OS already publishes a process's environment, per-PID and immutably, at
// /proc/<pid>/environ, complete and unredacted. A runtime-owned re-rendering
// could only add content nobody authored (which is precisely how it failed,
// brief Layer 0) or drop content the OS shows in full — either way a worse copy
// of a file that already exists. So the runtime maintains exactly one env file,
// the write surface, and leaves the read surface to the kernel.
//
// Raw files only — the surface works everywhere, including non-Linux hosts and
// degraded-FUSE habitats. The peer-facing FUSE projection is separate and
// read-only.

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kehao95/quine/internal/config"
)

func (r *Runtime) syncConfigSurface(agentRoot string) error {
	configDir := filepath.Join(agentRoot, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("mkdir config surface: %w", err)
	}
	// The agent writes config/env/override with ordinary shell writes, so the
	// directory must exist before it tries — the documented path just works.
	if err := os.MkdirAll(filepath.Join(configDir, config.EnvOverrideDirName), 0o755); err != nil {
		return fmt.Errorf("mkdir env surface: %w", err)
	}
	return r.syncConfigRegistrySurface(configDir)
}

func (r *Runtime) syncConfigRegistrySurface(configDir string) error {
	want, err := config.RegistryJSON()
	if err != nil {
		return fmt.Errorf("render capability registry: %w", err)
	}
	want = append(want, '\n')
	return syncReadOnlyProjection(filepath.Join(configDir, "registry.json"), want, "config registry surface")
}

// syncReadOnlyProjection writes a runtime-owned 0444 projection only when its
// content has drifted from what the compiled body would produce. The skip on
// equal content is what makes a per-safe-point refresh cheap; the rewrite on
// drift is what makes the projection tamper-evident rather than tamper-proof.
func syncReadOnlyProjection(path string, want []byte, label string) error {
	got, readErr := os.ReadFile(path)
	if readErr == nil && bytes.Equal(got, want) {
		return nil
	}
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		return fmt.Errorf("read %s: %w", label, readErr)
	}
	// writeReadOnlyFile stages a 0644 temp file and renames it over any stale
	// (possibly read-only) projection — rename replaces a 0444 target as long
	// as the directory is writable — then the read-only bit is restored.
	return writeReadOnlyFile(path, want, label)
}

// overrideAppliedName is the immutable applied-override archive written into
// the staging incarnation's inc/<n>/ dir when the successor consumes
// config/env/override.
const overrideAppliedName = "override-applied.env"

// consumeAppliedEnvOverride is phase 2 of the exec-time override transaction.
// When this bootstrap is a staged exec reentry (execReentryStaged: the
// QUINE_CONTEXT_BOOTSTRAP handover signal, still set here —
// importBootstrappedContext unsets it only after bootstrapAgentRoot) and
// config/env/override exists, its content is already baked into THIS process's
// envp: the predecessor applied it pre-syscall.Exec. All that remains is
// bookkeeping — archive the file into the PREDECESSOR's incarnation dir
// (inc/<n-1>, the incarnation that wrote it) as the immutable
// override-applied.env, then remove config/env/override so the policy surface
// is empty for this incarnation's own future.
//
// If the bootstrap is NOT a staged reentry (fresh start, session resume,
// fork/spawn child), any override is untouched: it belongs to this
// incarnation's own future children and its own future exec.
//
// Ordering: called from bootstrapAgentRoot BEFORE syncConfigSurface, for two
// reasons. (1) Transaction-before-declaration: an archive failure aborts
// bootstrap before the config/ surface claims a clean position with a
// consumed-but-present override sitting next to it. (2) It must run after
// ensureIncarnation (it needs the incarnation id to address inc/<n-1>), which
// bootstrapAgentRoot guarantees.
func (r *Runtime) consumeAppliedEnvOverride() error {
	if !execReentryStaged() {
		return nil
	}
	overridePath := r.cfg.EnvOverridePath()
	content, err := os.ReadFile(overridePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read applied env override %s: %w", overridePath, err)
	}

	archiveID := r.cfg.IncarnationID - 1
	if archiveID < 0 {
		// A staged exec reentry into a lineage with retained state always
		// advances past the predecessor's id, so id 0 here means the lineage
		// state is off (e.g. the retained root vanished between exec and
		// bootstrap). Keep the evidence anyway, under the current dir.
		r.log("env override: exec reentry at incarnation id %d has no predecessor dir; archiving under the current incarnation", r.cfg.IncarnationID)
		archiveID = r.cfg.IncarnationID
	}
	archiveDir := r.cfg.SessionIncarnationPath("", archiveID)
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return fmt.Errorf("mkdir env override archive dir: %w", err)
	}
	archivePath := filepath.Join(archiveDir, overrideAppliedName)
	// Archive the override VERBATIM. It is the agent's own authored policy file,
	// not the operator's inherited environment: the parser already rejects every
	// pinned/runtime-emitted knob (QUINE_API_KEY among them, operator-only), so
	// the file cannot carry a registry credential, and a foreign line the agent
	// chose to set (brief E9) is the agent's own authorship — withholding it here
	// would make the lineage record LIE about what policy was applied. The
	// operator-secret concern that would justify redaction does not arise: the
	// runtime never renders the inherited environment onto any surface (there is
	// no effective view and no birth snapshot); it only copies a file the agent
	// wrote.
	if err := writeReadOnlyFile(archivePath, content, "env override archive"); err != nil {
		return err
	}
	if err := os.Remove(overridePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("clear consumed env override %s: %w", overridePath, err)
	}
	return nil
}

// nonRegistryRuntimeEnv reports whether an inherited QUINE_* env name is
// legitimately outside the compiled capability registry: runtime-emitted
// process-surface signals and runtime-internal child plumbing that were never
// registry knobs. Authoritative sets: internal/config/registry.go (the registry
// proper, via config.KnobByEnv) plus the runtime-emitted family recorded in
// scripts/check-authored-env-consistency.sh's allowlist.
//
// The same family is classified for boundary purposes in
// config.BoundaryBehavior (runtimeOwnedNonRegistryEnv); this copy exists only
// so the drift log below can stay silent about them.
func nonRegistryRuntimeEnv(name string) bool {
	switch name {
	case "QUINE_AGENT_ROOT", // exec.go:execProcessSurfaceEnv and config.ShellStamps emit it
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
// paragraph). This is the observability half of the version-skew stance: a
// self-upgraded successor whose registry renamed or removed a knob silently
// ignores the inherited name in Load() by design — self-upgrade must not
// deadlock — and this bootstrap log line is what keeps that silence
// diagnosable in the lineage record (D4).
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
