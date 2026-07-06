package runtime

// config_gate_test.go covers the platform-independent half of the
// public/ctl/config validated write gate (work order T3.3): write-transaction
// semantics (land / reject / clear) and the computed read summary (staged
// content, validation state, coupling warnings, rejection read-back). The
// FUSE node wiring on top of this logic is covered by the Linux-gated
// TestBootstrapAgentRootFuseConfigGateTransactions in agent_root_test.go.

import (
	"os"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func newConfigGateRuntime(t *testing.T) *Runtime {
	t.Helper()
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	return rt
}

func TestConfigGateValidWriteLandsNextEnvVerbatim(t *testing.T) {
	rt := newConfigGateRuntime(t)
	payload := "# staged through the gate\nQUINE_MAX_TURNS=64\n"

	if err := rt.applyConfigStageWrite(payload); err != nil {
		t.Fatalf("applyConfigStageWrite(valid) failed: %v", err)
	}
	landed, err := os.ReadFile(rt.cfg.StagedNextEnvPath())
	if err != nil {
		t.Fatalf("read landed next.env: %v", err)
	}
	if string(landed) != payload {
		t.Fatalf("landed next.env = %q, want byte-equal payload %q", landed, payload)
	}

	summary := string(rt.configControlSurfaceSummary())
	for _, want := range []string{
		"control_file: config",
		"staged: 2 line(s), config/next.env verbatim below",
		"QUINE_MAX_TURNS=64",
		"validation: valid against the running capability registry",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("gate summary missing %q:\n%s", want, summary)
		}
	}
	for _, forbidden := range []string{"coupling_warnings", "last_rejected_write"} {
		if strings.Contains(summary, forbidden) {
			t.Errorf("gate summary should not contain %q for a clean coupling-free stage:\n%s", forbidden, summary)
		}
	}
}

func TestConfigGateRejectedWritePreservesStagedFileAndSurfacesViolations(t *testing.T) {
	rt := newConfigGateRuntime(t)
	prior := "QUINE_MAX_TURNS=64\n"
	if err := rt.applyConfigStageWrite(prior); err != nil {
		t.Fatalf("stage prior payload: %v", err)
	}

	// Whole-transaction reject: one unknown knob, one type violation, one
	// forbidden mutability — nothing may land, the prior stage must survive.
	invalid := "QUINE_TOTALLY_UNKNOWN_KNOB=1\nQUINE_MAX_TURNS=abc\nQUINE_API_KEY=stolen\n"
	err := rt.applyConfigStageWrite(invalid)
	if err == nil {
		t.Fatal("applyConfigStageWrite(invalid) should reject")
	}
	surviving, readErr := os.ReadFile(rt.cfg.StagedNextEnvPath())
	if readErr != nil {
		t.Fatalf("read staged file after rejection: %v", readErr)
	}
	if string(surviving) != prior {
		t.Fatalf("rejected replacement clobbered the staged file: got %q, want %q", surviving, prior)
	}

	summary := string(rt.configControlSurfaceSummary())
	for _, want := range []string{
		"last_rejected_write: rejected in full, nothing landed:",
		"QUINE_TOTALLY_UNKNOWN_KNOB: unknown env name",
		"invalid int value \"abc\"",
		"only exec-boundary knobs can be staged",
		// The surviving prior stage is still reported as the staged state.
		"QUINE_MAX_TURNS=64",
		"validation: valid against the running capability registry",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("post-rejection summary missing %q:\n%s", want, summary)
		}
	}

	// The next accepted write clears the rejection read-back.
	if err := rt.applyConfigStageWrite("QUINE_OUTPUT_TRUNCATE=31337\n"); err != nil {
		t.Fatalf("accepted write after rejection failed: %v", err)
	}
	summary = string(rt.configControlSurfaceSummary())
	if strings.Contains(summary, "last_rejected_write") {
		t.Errorf("accepted write should clear the rejection read-back:\n%s", summary)
	}
	if !strings.Contains(summary, "QUINE_OUTPUT_TRUNCATE=31337") {
		t.Errorf("accepted write should replace the stage wholesale:\n%s", summary)
	}
}

func TestConfigGateEmptyWriteClearsStagedFile(t *testing.T) {
	rt := newConfigGateRuntime(t)
	if err := rt.applyConfigStageWrite("QUINE_MAX_TURNS=64\n"); err != nil {
		t.Fatalf("stage payload: %v", err)
	}

	if err := rt.applyConfigStageWrite(""); err != nil {
		t.Fatalf("empty write should clear the stage, got: %v", err)
	}
	if _, err := os.Stat(rt.cfg.StagedNextEnvPath()); !os.IsNotExist(err) {
		t.Fatalf("staged file should be removed by an empty write, got err=%v", err)
	}
	if summary := string(rt.configControlSurfaceSummary()); !strings.Contains(summary, "staged: none") {
		t.Errorf("summary after clear should report staged: none:\n%s", summary)
	}

	// Clearing an already-empty stage is a clean no-op.
	if err := rt.applyConfigStageWrite(""); err != nil {
		t.Fatalf("empty write on empty stage should be a no-op, got: %v", err)
	}
}

func TestConfigGateSummarySurfacesCouplingWarnings(t *testing.T) {
	rt := newConfigGateRuntime(t)
	// EphemeralBodyEnabled carries the registry's first Couples entry
	// (archaeology coupling #17, hazard edge to SelfReentryMode).
	if err := rt.applyConfigStageWrite("QUINE_EPHEMERAL_BODY_ENABLED=1\n"); err != nil {
		t.Fatalf("stage coupled knob (accepted; warnings are not errors): %v", err)
	}

	summary := string(rt.configControlSurfaceSummary())
	for _, want := range []string{
		"validation: valid against the running capability registry",
		"coupling_warnings: staged knobs carry registry coupling edges (warnings, not errors",
		"QUINE_EPHEMERAL_BODY_ENABLED couples with QUINE_SELF_REENTRY_MODE (hazard):",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("coupling summary missing %q:\n%s", want, summary)
		}
	}
}

func TestConfigGateSummaryRevalidatesShWrittenStagedFile(t *testing.T) {
	rt := newConfigGateRuntime(t)
	// The raw sh path stays legal: a directly-written (gate-bypassing)
	// invalid next.env must be reported INVALID by the gate's read summary —
	// validation state is recomputed on every read, never cached from the
	// last gate write.
	path := rt.cfg.StagedNextEnvPath()
	if err := os.MkdirAll(strings.TrimSuffix(path, "/"+config.StagedNextEnvName), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("QUINE_SESSION_ID=forged\n"), 0o644); err != nil {
		t.Fatalf("write staged file directly: %v", err)
	}

	summary := string(rt.configControlSurfaceSummary())
	for _, want := range []string{
		"staged: 1 line(s), config/next.env verbatim below",
		"QUINE_SESSION_ID=forged",
		"validation: INVALID against the running capability registry; an exec now would reject it:",
		"only exec-boundary knobs can be staged",
	} {
		if !strings.Contains(summary, want) {
			t.Errorf("sh-written invalid summary missing %q:\n%s", want, summary)
		}
	}
}
