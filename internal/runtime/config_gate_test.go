package runtime

// config_gate_test.go covers the platform-independent half of the public/ctl/env
// validated write gate: write-transaction semantics (land / reject / clear) and
// the computed read summary (current policy, validation state, coupling
// warnings, rejection read-back). The FUSE node wiring on top of this logic is
// covered by the Linux-gated TestBootstrapAgentRootFuseEnvGateTransactions in
// agent_root_test.go.
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//
// The gate is the polite door, not the wall. Raw sh writes to
// config/env/override stay legal, and every envp constructor revalidates the
// file no matter who wrote it — which is why the read summary recomputes
// validation on every read instead of caching what the last gate write decided.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newEnvGateRuntime(t *testing.T) *Runtime {
	t.Helper()
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	return rt
}

func TestEnvGateValidWriteLandsOverrideVerbatim(t *testing.T) {
	rt := newEnvGateRuntime(t)
	// The staged names and values are deliberately NOT the ones in the gate's
	// static usage/example help text: this test forbids value echo by substring,
	// so a payload that reused the example (QUINE_MAX_TURNS=64, FOO=bar) would
	// trip on the help line, not on a real leak. QUINE_OUTPUT_TRUNCATE carries no
	// coupling edges, so a clean valid write emits no coupling_warnings.
	payload := "# authored through the gate\nQUINE_OUTPUT_TRUNCATE=51234\nGATE_PROBE_XYZ=zzz9\n"

	if err := rt.applyEnvOverrideWrite(payload); err != nil {
		t.Fatalf("applyEnvOverrideWrite(valid) failed: %v", err)
	}
	landed, err := os.ReadFile(rt.cfg.EnvOverridePath())
	if err != nil {
		t.Fatalf("read landed config/env/override: %v", err)
	}
	if string(landed) != payload {
		t.Fatalf("landed override = %q, want byte-equal payload %q", landed, payload)
	}

	summary := string(rt.envControlSurfaceSummary())
	for _, want := range []string{
		"control_file: env",
		"mode: validated-child-env-policy",
		// Names and operations, not values: this node is READ BY PEERS, and the
		// override is the agent's own file that the FUSE config projection
		// deliberately keeps off the peer surface.
		"policy_names: QUINE_OUTPUT_TRUNCATE (set), GATE_PROBE_XYZ (set)",
		"validation: valid against the running capability registry",
		// The non-claim travels with the gate: it shapes children, not the peer.
		"it does not change the peer's own environment, which is fixed at its birth",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("gate summary missing %q:\n%s", want, summary)
		}
	}
	// The staged VALUES must NOT be echoed to a peer-readable node (the KEY=VALUE
	// forms and the raw value strings), and a clean policy carries no
	// coupling/rejection block.
	for _, forbidden := range []string{"QUINE_OUTPUT_TRUNCATE=51234", "GATE_PROBE_XYZ=zzz9", "zzz9", "51234", "coupling_warnings", "last_rejected_write"} {
		if strings.Contains(summary, forbidden) {
			t.Errorf("gate summary should not contain %q — values are the agent's, not the peer's; and no coupling/rejection for a clean policy:\n%s", forbidden, summary)
		}
	}
}

// TestEnvGateSummaryNeverEchoesForeignValuesToPeers is the regression test for
// the peer disclosure the gate read-back created. The override's advertised
// purpose is handing env to children, and foreign names are free-form (E9), so
// an agent legitimately writes GITHUB_TOKEN=… there. The read side of this gate
// is a peer-readable node under public/; echoing the file verbatim handed every
// peer on the host the token. The FUSE config projection goes to deliberate
// lengths to keep config/env/override off the peer surface — this node must not
// undo that.
func TestEnvGateSummaryNeverEchoesForeignValuesToPeers(t *testing.T) {
	rt := newEnvGateRuntime(t)
	if err := rt.applyEnvOverrideWrite("GITHUB_TOKEN=ghp_PEER_MUST_NOT_SEE\nQUINE_MAX_TURNS=64\nLC_ALL\n"); err != nil {
		t.Fatalf("a legal policy with a foreign credential-shaped value: %v", err)
	}
	summary := string(rt.envControlSurfaceSummary())
	if strings.Contains(summary, "ghp_PEER_MUST_NOT_SEE") {
		t.Fatalf("ctl/env echoed a foreign value to a peer-readable node:\n%s", summary)
	}
	// A peer still learns which names the policy speaks to, and that it is valid.
	for _, want := range []string{
		"GITHUB_TOKEN (set)", "QUINE_MAX_TURNS (set)", "LC_ALL (unset)",
		"validation: valid against the running capability registry",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("gate summary missing %q — a peer needs the names and the verdict, just not the values:\n%s", want, summary)
		}
	}
}

// TestEnvGateAcceptsForeignNamesAndUnsets: the override is a general child-env
// manager, not a QUINE_* knob poker (brief E9). A foreign name is free-form, and
// a bare KEY line unsets an inherited var — the persistent-`export` capability
// that did not exist under the old staged-config channel.
func TestEnvGateAcceptsForeignNamesAndUnsets(t *testing.T) {
	rt := newEnvGateRuntime(t)
	payload := "FOO=bar\nEMPTY=\nLANG\nQUINE_MAX_TURNS\n"

	if err := rt.applyEnvOverrideWrite(payload); err != nil {
		t.Fatalf("foreign names, empty sets, and bare-KEY unsets are all legal: %v", err)
	}
	landed, err := os.ReadFile(rt.cfg.EnvOverridePath())
	if err != nil {
		t.Fatalf("read landed override: %v", err)
	}
	if string(landed) != payload {
		t.Fatalf("landed override = %q, want byte-equal %q", landed, payload)
	}
	summary := string(rt.envControlSurfaceSummary())
	if !strings.Contains(summary, "validation: valid against the running capability registry") {
		t.Fatalf("gate should accept the policy:\n%s", summary)
	}
}

func TestEnvGateRejectedWritePreservesOverrideAndSurfacesViolations(t *testing.T) {
	rt := newEnvGateRuntime(t)
	prior := "QUINE_MAX_TURNS=64\n"
	if err := rt.applyEnvOverrideWrite(prior); err != nil {
		t.Fatalf("land the prior policy: %v", err)
	}

	// Whole-transaction reject, enumerating every violation: an unknown knob, a
	// type violation, a runtime-emitted (masked) name, a pinned operator-only
	// knob, and a runtime-owned non-registry name. Nothing may land; the prior
	// policy must survive byte-untouched.
	invalid := strings.Join([]string{
		"QUINE_TOTALLY_UNKNOWN_KNOB=1",
		"QUINE_MAX_TURNS=abc",
		"QUINE_SESSION_ID=forged",
		"QUINE_MAX_DEPTH=99",
		"QUINE_AGENT_ROOT=/elsewhere",
		"",
	}, "\n")
	err := rt.applyEnvOverrideWrite(invalid)
	if err == nil {
		t.Fatal("applyEnvOverrideWrite(invalid) should reject")
	}
	surviving, readErr := os.ReadFile(rt.cfg.EnvOverridePath())
	if readErr != nil {
		t.Fatalf("read override after rejection: %v", readErr)
	}
	if string(surviving) != prior {
		t.Fatalf("a rejected replacement clobbered the override: got %q, want %q", surviving, prior)
	}

	summary := string(rt.envControlSurfaceSummary())
	for _, want := range []string{
		"last_rejected_write: rejected in full, nothing landed:",
		"QUINE_TOTALLY_UNKNOWN_KNOB: unknown env name",
		`invalid int value "abc"`,
		"QUINE_SESSION_ID: runtime-emitted knob",
		// E5: tree budgets are the operator's, not the lineage's — a lineage
		// must not relax its own depth budget through a mediated channel.
		`QUINE_MAX_DEPTH: mutability "operator-only" — pinned`,
		"QUINE_AGENT_ROOT: runtime-owned name",
		// The surviving prior policy is still reported as the current state.
		"QUINE_MAX_TURNS=64",
		"validation: valid against the running capability registry",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("post-rejection summary missing %q:\n%s", want, summary)
		}
	}

	// The next accepted write clears the rejection read-back and replaces the
	// policy wholesale — there is no append, and no merge with what was there.
	replacement := "QUINE_OUTPUT_TRUNCATE=31337\n"
	if err := rt.applyEnvOverrideWrite(replacement); err != nil {
		t.Fatalf("accepted write after rejection failed: %v", err)
	}
	landed, err := os.ReadFile(rt.cfg.EnvOverridePath())
	if err != nil {
		t.Fatalf("read override after the replacement: %v", err)
	}
	if string(landed) != replacement {
		t.Errorf("wholesale replacement: override = %q, want exactly %q (the prior policy is gone, not merged)", landed, replacement)
	}
	summary = string(rt.envControlSurfaceSummary())
	if strings.Contains(summary, "last_rejected_write") {
		t.Errorf("an accepted write should clear the rejection read-back:\n%s", summary)
	}
	if !strings.Contains(summary, "policy_names: QUINE_OUTPUT_TRUNCATE (set)") {
		t.Errorf("summary should report the replacement policy's names as the only one in force:\n%s", summary)
	}
	if strings.Contains(summary, "QUINE_OUTPUT_TRUNCATE=31337") {
		t.Errorf("summary must not echo the policy's values to a peer-readable node:\n%s", summary)
	}
}

func TestEnvGateEmptyWriteClearsOverride(t *testing.T) {
	rt := newEnvGateRuntime(t)
	if err := rt.applyEnvOverrideWrite("QUINE_MAX_TURNS=64\n"); err != nil {
		t.Fatalf("land a policy: %v", err)
	}

	if err := rt.applyEnvOverrideWrite(""); err != nil {
		t.Fatalf("an empty write should clear the policy, got: %v", err)
	}
	if _, err := os.Stat(rt.cfg.EnvOverridePath()); !os.IsNotExist(err) {
		t.Fatalf("config/env/override should be removed by an empty write, got err=%v", err)
	}
	if summary := string(rt.envControlSurfaceSummary()); !strings.Contains(summary, "policy: none") {
		t.Errorf("summary after clear should report policy: none:\n%s", summary)
	}

	// Clearing an already-empty policy is a clean no-op.
	if err := rt.applyEnvOverrideWrite(""); err != nil {
		t.Fatalf("an empty write on an empty policy should be a no-op, got: %v", err)
	}
}

func TestEnvGateSummarySurfacesCouplingWarnings(t *testing.T) {
	rt := newEnvGateRuntime(t)
	// EphemeralBodyEnabled carries a registry hazard edge to SelfReentryMode
	// (archaeology coupling #17). Legality passed; wisdom is the agent's.
	if err := rt.applyEnvOverrideWrite("QUINE_EPHEMERAL_BODY_ENABLED=1\n"); err != nil {
		t.Fatalf("a coupled knob is accepted; warnings are not errors: %v", err)
	}

	summary := string(rt.envControlSurfaceSummary())
	for _, want := range []string{
		"validation: valid against the running capability registry",
		"coupling_warnings: named knobs carry registry coupling edges (warnings, not errors",
		"QUINE_EPHEMERAL_BODY_ENABLED couples with QUINE_SELF_REENTRY_MODE (hazard):",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("coupling summary missing %q:\n%s", want, summary)
		}
	}
}

func TestEnvGateSummaryRevalidatesShWrittenOverride(t *testing.T) {
	rt := newEnvGateRuntime(t)
	// The raw sh path stays legal, and the gate has no hook on it. A
	// directly-written (gate-bypassing) illegal policy must therefore be
	// reported INVALID by the read summary: validation state is recomputed on
	// every read, never cached from the last gate write.
	path := rt.cfg.EnvOverridePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir config/env: %v", err)
	}
	if err := os.WriteFile(path, []byte("QUINE_SESSION_ID=forged\n"), 0o644); err != nil {
		t.Fatalf("write override directly: %v", err)
	}

	summary := string(rt.envControlSurfaceSummary())
	// The forged VALUE is not echoed to peers even on the invalid path; the
	// violation names the offending knob without reproducing what was written.
	if strings.Contains(summary, "QUINE_SESSION_ID=forged") {
		t.Errorf("an invalid sh-written policy must not have its raw line echoed to a peer-readable node:\n%s", summary)
	}
	for _, want := range []string{
		"validation: INVALID against the running capability registry; it applies to no child and an exec would reject it:",
		"QUINE_SESSION_ID: runtime-emitted knob",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("sh-written invalid summary missing %q:\n%s", want, summary)
		}
	}
}
