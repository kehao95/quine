package runtime

// config_gate.go is the feedback layer of the child-env override write path:
// the validated write gate behind the public/ctl/env FUSE node. A write is one
// whole config/env/override payload (it IS the file's content); the gate
// validates it against the running binary's compiled capability registry via
// config.ParseEnvOverride — the same validator the envp constructors use, so
// gate and boundary can never diverge — lands the file atomically when legal,
// and rejects the transaction when not. Reads return a computed summary: the
// current policy, its validation state, registry coupling warnings, and the
// violations of the last rejected gate write.
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//
// The gate is convenience, not authority: raw sh writes to config/env/override
// stay legal, and every envp constructor (sh, fork, spawn, exec) revalidates
// the file no matter who wrote it. That is what makes the boundary
// un-fakeable rather than merely well-advertised — enforcement is compiled in,
// and this gate is the polite door beside it.
//
// Validation checks legality, not wisdom (registry brief D4): coupling edges
// are surfaced as warnings on an accepted policy, never as errors.
//
// This file is platform-independent on purpose: only the FUSE node wiring
// (public_surface_fuse_linux.go) is Linux-only, so in the degraded-public/
// regime the gate is simply absent and the raw file remains the universal path.

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

// envGateState is the per-process transaction-feedback memory of the
// public/ctl/env gate: the violations of the most recent rejected write,
// retrievable from the node's read summary until the next accepted write. The
// override file itself is the only durable state; this is ack/reject feedback,
// which dies with the process like the FUSE node that serves it.
type envGateState struct {
	mu            sync.Mutex
	lastRejection []string
}

// applyEnvOverrideWrite is the write half of the public/ctl/env gate: one whole
// child-env policy per call.
//
//   - A non-empty payload is validated with config.ParseEnvOverride. Legal
//     payloads land config/env/override atomically (writeFile temp+rename),
//     REPLACING the file wholesale — matching the whole-file validation
//     semantics; there is no append or merge. Illegal payloads reject the
//     transaction in full: nothing lands, any pre-existing policy survives
//     byte-untouched, and the violations become readable from the gate until
//     the next accepted write.
//   - An empty payload clears the policy: replace-wholesale semantics make an
//     empty file mean "children inherit what I was given", and removing the
//     file is that state's canonical spelling (every constructor treats an
//     absent file as a clean no-op).
func (r *Runtime) applyEnvOverrideWrite(payload string) error {
	path := r.cfg.EnvOverridePath()
	if path == "" {
		return errors.New("env override gate: no agent root configured")
	}
	if len(payload) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("clear env override %s: %w", path, err)
		}
		r.recordEnvGateResult(nil)
		return nil
	}
	if _, err := config.ParseEnvOverride([]byte(payload)); err != nil {
		r.recordEnvGateResult(err)
		return fmt.Errorf("env override rejected (whole payload, nothing landed; violations readable at ctl/env): %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir env override dir: %w", err)
	}
	if err := writeFile(path, []byte(payload)); err != nil {
		return err
	}
	r.recordEnvGateResult(nil)
	return nil
}

// recordEnvGateResult updates the gate's transaction feedback: an accepted
// write (nil) clears the last rejection; a rejected write stores its violations
// for the read summary.
func (r *Runtime) recordEnvGateResult(err error) {
	r.envGate.mu.Lock()
	defer r.envGate.mu.Unlock()
	if err == nil {
		r.envGate.lastRejection = nil
		return
	}
	var rejected *config.EnvOverrideError
	if errors.As(err, &rejected) {
		r.envGate.lastRejection = append([]string(nil), rejected.Violations...)
		return
	}
	r.envGate.lastRejection = []string{err.Error()}
}

// envControlSurfaceSummary renders the read side of public/ctl/env: a computed
// summary of the gate's semantics, the current policy (or "none"), that
// policy's validation state against the RUNNING registry (re-checked on every
// read — sh may have written the file directly), registry coupling warnings,
// and the violations of the last rejected gate write.
func (r *Runtime) envControlSurfaceSummary() []byte {
	lines := []string{
		fmt.Sprintf("backend: %s", runtimeSurfaceBackendName),
		"control_file: env",
		"mode: validated-child-env-policy",
		"usage: write ONE whole config/env/override payload (KEY=VALUE sets, a bare KEY unsets, # comments, blank lines; values verbatim, no quoting/expansion); each write REPLACES the file wholesale",
		`example: printf 'QUINE_MAX_TURNS=64\nFOO=bar\nLANG\n' > ctl/env`,
		"empty_write: : > ctl/env clears the policy (children inherit unchanged)",
		"rejection: an illegal payload rejects the whole transaction at close (EINVAL) and lands nothing; its violations stay readable below until the next accepted write",
		"enforcement: this gate is write-time feedback; every process this peer builds (sh, fork, spawn, exec) revalidates config/env/override against its running registry no matter who wrote it (raw sh writes stay legal)",
		"scope: the policy shapes processes this peer CONSTRUCTS; it does not change the peer's own environment, which is fixed at its birth",
	}

	path := r.cfg.EnvOverridePath()
	var content []byte
	readErr := os.ErrNotExist
	if path != "" {
		content, readErr = os.ReadFile(path)
	}
	switch {
	case errors.Is(readErr, os.ErrNotExist):
		lines = append(lines, "policy: none")
	case readErr != nil:
		lines = append(lines, fmt.Sprintf("policy: unreadable (%v)", readErr))
	default:
		body := strings.TrimSuffix(string(content), "\n")
		contentLines := []string{}
		if body != "" {
			contentLines = strings.Split(body, "\n")
		}
		lines = append(lines, fmt.Sprintf("policy: %d line(s); names only below — this node is READ BY PEERS, and config/env/override is the agent's own file, deliberately not projected to them", len(contentLines)))
		if override, err := config.ParseEnvOverride(content); err != nil {
			lines = append(lines, "validation: INVALID against the running capability registry; it applies to no child and an exec would reject it:")
			lines = append(lines, envViolationLines(err)...)
		} else {
			lines = append(lines, "policy_names: "+envPolicyNames(override))
			lines = append(lines, "validation: valid against the running capability registry")
			if warnings := envCouplingWarnings(override); len(warnings) > 0 {
				lines = append(lines, "coupling_warnings: named knobs carry registry coupling edges (warnings, not errors — legality passed, wisdom is yours):")
				for _, warning := range warnings {
					lines = append(lines, "  - "+warning)
				}
			}
		}
	}

	r.envGate.mu.Lock()
	rejection := append([]string(nil), r.envGate.lastRejection...)
	r.envGate.mu.Unlock()
	if len(rejection) > 0 {
		lines = append(lines, "last_rejected_write: rejected in full, nothing landed:")
		for _, violation := range rejection {
			lines = append(lines, "  - "+violation)
		}
	}

	return []byte(strings.Join(lines, "\n") + "\n")
}

// envPolicyNames renders the policy as NAMES and operations, never values.
//
// The read side of this gate is a peer-facing node under public/: any quine on
// the host can cat it. The FUSE config projection goes to deliberate lengths to
// keep config/env/override off that surface — fuseConfigEnvDirNode is a fixed
// 2-name enumeration rather than a directory scan precisely so the agent's own
// policy file never becomes peer-readable — and the prompt tells the agent as
// much ("your config/env/override is not projected: it is yours"). Echoing the
// file's content verbatim here handed peers the very bytes that node refuses to
// show them, and the override is where a legal foreign line like
// GITHUB_TOKEN=ghp_… lives, since foreign names are free-form by design (E9).
//
// Names are what a peer needs to reason about the gate — which names this policy
// speaks to, and whether the runtime accepted it. Values are the agent's.
func envPolicyNames(override *config.EnvOverride) string {
	if override == nil || len(override.Names) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(override.Names))
	for _, name := range override.Names {
		op := "set"
		if override.Unsets[name] {
			op = "unset"
		}
		parts = append(parts, fmt.Sprintf("%s (%s)", name, op))
	}
	return strings.Join(parts, ", ")
}

// envViolationLines flattens an override validation error into summary bullet
// lines, preserving the per-violation enumeration of *config.EnvOverrideError.
func envViolationLines(err error) []string {
	var rejected *config.EnvOverrideError
	if errors.As(err, &rejected) {
		out := make([]string, 0, len(rejected.Violations))
		for _, violation := range rejected.Violations {
			out = append(out, "  - "+violation)
		}
		return out
	}
	return []string{"  - " + err.Error()}
}

// envCouplingWarnings renders one warning per registry Couples edge of each
// knob the policy names — edit-site visibility: the combination that fails at
// the next boundary (or silently no-ops) is named where the value is written.
// Unset lines are included: reverting a knob to its compiled default moves the
// same coupled edge as setting it. Deterministic order: named knobs sorted,
// edges in registry declaration order. Peer knobs are named by their env
// encoding (the spelling the policy file uses), falling back to the registry
// Name if the peer entry is ever missing.
func envCouplingWarnings(override *config.EnvOverride) []string {
	if override == nil {
		return nil
	}
	names := append([]string(nil), override.Names...)
	sort.Strings(names)
	var warnings []string
	for _, env := range names {
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
