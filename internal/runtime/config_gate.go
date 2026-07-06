package runtime

// config_gate.go is the feedback layer of the staged-config write path: the
// validated write gate behind the public/ctl/config FUSE node. A write is one
// whole config/next.env payload (same env grammar — it IS next.env's content);
// the gate validates it against the running binary's compiled capability
// registry via config.ParseStagedEnv (the single source of validation truth,
// shared with the exec boundary), lands the SST atomically when legal, and
// rejects the transaction when not. Reads return a computed summary: current
// staged content, its validation state, registry coupling warnings for staged
// knobs, and the violations of the last rejected gate write.
//
// Design authority:
//   Paper/theory/views/runtime-capability/registry-design-brief.md (§ C
//   feedback layer)
// Work order:
//   Paper/_design/migrations/runtime-capability-registry-execution.md (T3.3)
//
// The gate is convenience, not authority: raw sh writes to config/next.env
// stay legal, and the exec-boundary merge (internal/tools/exec.go via
// config.ReadStagedOverrides) revalidates the file no matter who wrote it.
// Validation checks legality, not wisdom (brief D4): coupling edges are
// surfaced as warnings on accepted stages, never as errors.
//
// This file is platform-independent on purpose: only the FUSE node wiring
// (public_surface_fuse_linux.go) is Linux-only, so in the degraded-public/
// regime the gate is simply absent and next.env remains the universal path.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/kehao95/quine/internal/config"
)

// configGateState is the per-process transaction-feedback memory of the
// public/ctl/config gate: the violations of the most recent rejected write,
// retrievable from the node's read summary until the next accepted write.
// The staged file itself is the only durable state; this is ack/reject
// feedback, which dies with the process like the FUSE node that serves it.
type configGateState struct {
	mu            sync.Mutex
	lastRejection []string
}

// applyConfigStageWrite is the write half of the public/ctl/config gate: one
// whole staged-config transaction per call.
//
//   - A non-empty payload is validated with config.ParseStagedEnv (the exec
//     boundary's own validator — gate and boundary can never diverge). Legal
//     payloads land config/next.env atomically (writeFile temp+rename),
//     REPLACING the staged file wholesale — matching the whole-file
//     validation semantics; there is no append or merge. Illegal payloads
//     reject the transaction in full: nothing lands, any pre-existing staged
//     file survives untouched, and the violations become readable from the
//     gate until the next accepted write.
//   - An empty payload clears the staged file: replace-wholesale semantics
//     make an empty stage mean "no overrides", and removing the file is that
//     state's canonical spelling (the exec boundary treats an absent file as
//     a clean no-op).
func (r *Runtime) applyConfigStageWrite(payload string) error {
	path := r.cfg.StagedNextEnvPath()
	if path == "" {
		return errors.New("staged-config gate: no agent root configured")
	}
	if len(payload) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("clear staged config %s: %w", path, err)
		}
		r.recordConfigGateResult(nil)
		return nil
	}
	if _, err := config.ParseStagedEnv([]byte(payload)); err != nil {
		r.recordConfigGateResult(err)
		return fmt.Errorf("staged config rejected (whole payload, nothing landed; violations readable at ctl/config): %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir staged config dir: %w", err)
	}
	if err := writeFile(path, []byte(payload)); err != nil {
		return err
	}
	r.recordConfigGateResult(nil)
	return nil
}

// recordConfigGateResult updates the gate's transaction feedback: an accepted
// write (nil) clears the last rejection; a rejected write stores its
// violations for the read summary.
func (r *Runtime) recordConfigGateResult(err error) {
	r.configGate.mu.Lock()
	defer r.configGate.mu.Unlock()
	if err == nil {
		r.configGate.lastRejection = nil
		return
	}
	var staged *config.StagedFileError
	if errors.As(err, &staged) {
		r.configGate.lastRejection = append([]string(nil), staged.Violations...)
		return
	}
	r.configGate.lastRejection = []string{err.Error()}
}

// configControlSurfaceSummary renders the read side of public/ctl/config: a
// computed summary of the gate's semantics, the current staged content (or
// "none"), that content's validation state against the RUNNING registry
// (re-checked on every read — sh may have written next.env directly), registry
// coupling warnings for staged knobs, and the violations of the last rejected
// gate write.
func (r *Runtime) configControlSurfaceSummary() []byte {
	lines := []string{
		fmt.Sprintf("backend: %s", runtimeSurfaceBackendName),
		"control_file: config",
		"mode: validated-config-stage",
		"usage: write ONE whole config/next.env payload (env syntax: KEY=VALUE lines, # comments, blank lines; values verbatim, no quoting/expansion); each write REPLACES the staged file wholesale",
		`example: printf 'QUINE_MAX_TURNS=64\n' > ctl/config`,
		"empty_write: : > ctl/config clears the staged file (an empty stage means no overrides)",
		"rejection: an illegal payload rejects the whole transaction at close (EINVAL) and lands nothing; its violations stay readable below until the next accepted write",
		"enforcement: this gate is write-time feedback; the exec boundary revalidates config/next.env against the running registry no matter who wrote it (raw sh writes to config/next.env stay legal)",
	}

	path := r.cfg.StagedNextEnvPath()
	var staged []byte
	readErr := os.ErrNotExist
	if path != "" {
		staged, readErr = os.ReadFile(path)
	}
	switch {
	case errors.Is(readErr, os.ErrNotExist):
		lines = append(lines, "staged: none")
	case readErr != nil:
		lines = append(lines, fmt.Sprintf("staged: unreadable (%v)", readErr))
	default:
		content := strings.TrimSuffix(string(staged), "\n")
		contentLines := []string{}
		if content != "" {
			contentLines = strings.Split(content, "\n")
		}
		lines = append(lines, fmt.Sprintf("staged: %d line(s), config/next.env verbatim below", len(contentLines)))
		lines = append(lines, contentLines...)
		if overrides, err := config.ParseStagedEnv(staged); err != nil {
			lines = append(lines, "validation: INVALID against the running capability registry; an exec now would reject it:")
			lines = append(lines, stagedViolationLines(err)...)
		} else {
			lines = append(lines, "validation: valid against the running capability registry")
			if warnings := stagedCouplingWarnings(overrides); len(warnings) > 0 {
				lines = append(lines, "coupling_warnings: staged knobs carry registry coupling edges (warnings, not errors — legality passed, wisdom is yours):")
				for _, warning := range warnings {
					lines = append(lines, "  - "+warning)
				}
			}
		}
	}

	r.configGate.mu.Lock()
	rejection := append([]string(nil), r.configGate.lastRejection...)
	r.configGate.mu.Unlock()
	if len(rejection) > 0 {
		lines = append(lines, "last_rejected_write: rejected in full, nothing landed:")
		for _, violation := range rejection {
			lines = append(lines, "  - "+violation)
		}
	}

	return []byte(strings.Join(lines, "\n") + "\n")
}

// stagedViolationLines flattens a staged-config validation error into
// summary bullet lines, preserving the per-violation enumeration of
// *config.StagedFileError.
func stagedViolationLines(err error) []string {
	var staged *config.StagedFileError
	if errors.As(err, &staged) {
		out := make([]string, 0, len(staged.Violations))
		for _, violation := range staged.Violations {
			out = append(out, "  - "+violation)
		}
		return out
	}
	return []string{"  - " + err.Error()}
}

// stagedCouplingWarnings renders one warning per registry Couples edge of
// each staged knob — constraint 2's edit-site visibility: the combination
// that fails at the next exec (or silently no-ops) is named where the value
// is staged. Deterministic order: staged env names sorted, edges in registry
// declaration order. Peer knobs are named by their env encoding (the spelling
// a staged file uses), falling back to the registry Name if the peer entry is
// ever missing.
func stagedCouplingWarnings(staged map[string]string) []string {
	envs := make([]string, 0, len(staged))
	for env := range staged {
		envs = append(envs, env)
	}
	sort.Strings(envs)
	var warnings []string
	for _, env := range envs {
		knob, ok := config.KnobByEnv(env)
		if !ok {
			continue
		}
		for _, edge := range knob.Couples {
			peerEnv := edge.Peer
			if peer, ok := config.KnobByName(edge.Peer); ok {
				peerEnv = peer.Env
			}
			warnings = append(warnings, fmt.Sprintf("%s couples with %s (%s): %s", env, peerEnv, edge.Kind, edge.Note))
		}
	}
	return warnings
}
