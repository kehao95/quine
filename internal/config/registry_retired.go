package config

// registry_retired.go models the RETIRED half of the env namespace: names the
// runtime once read (or that authored surfaces once set) that were deliberately
// removed and now resolve to nothing. They are NOT live registry knobs — the
// registry <-> envnames.go bijection (registry_test.go) must stay clean — but
// the generated env-controls.md keeps their rows so the deletion decisions
// survive regeneration.
//
// Work order:
//   Paper/_design/migrations/runtime-capability-registry-execution.md (T1.2)
// Design authority:
//   Paper/theory/views/runtime-capability/registry-design-brief.md (D5/D8/D9)
//
// Scope rule: only silently-ignored retired names belong here. The two
// load-error tombstones (QUINE_SELF_REENTRY_TARGET, QUINE_WORKSPACE_SOURCE)
// are LIVE registry entries with Default kind "legacy" — the runtime still
// reads them to reject them — so they are deliberately absent from this table.

// RetiredKnob is one retired env name kept as documentation data.
//
// Note carries the env-controls.md Notes cell verbatim (markdown; relative
// links resolve from Paper/core/registries/). Decision/When/Replacement are
// the structured form for future consumers (e.g. suggesting the successor
// name when an authored surface still sets a retired env).
type RetiredKnob struct {
	Env         string `json:"env"`                   // retired env name; may be a family pattern such as QUINE_WISDOM_*
	Decision    string `json:"decision"`              // owning decision record (registry-design-brief D-id or decision doc slug)
	When        string `json:"when"`                  // retirement date (commit date of the deleting change)
	Replacement string `json:"replacement,omitempty"` // successor env, when one exists
	Note        string `json:"note"`                  // env-controls.md Notes cell, verbatim markdown
}

// RetiredRegistry lists every retired env name the generated env-controls.md
// documents. Setting any of these is ignored by the runtime (unknown env
// names never error); check-authored-env-consistency keeps authored surfaces
// from still using them.
var RetiredRegistry = []RetiredKnob{
	{
		Env: "QUINE_SMART_MODEL_ID", Decision: "brief D9", When: "2026-07-03",
		Note: "retired with the `escalate` feature deletion ([`registry-design-brief.md` D9](../../theory/views/runtime-capability/registry-design-brief.md)); setting it is ignored (unknown env names never error), not an error",
	},
	{
		Env: "QUINE_SMART_API_TYPE", Decision: "brief D9", When: "2026-07-03",
		Note: "retired with the `escalate` feature deletion (brief D9); setting it is ignored",
	},
	{
		Env: "QUINE_SMART_API_BASE", Decision: "brief D9", When: "2026-07-03",
		Note: "retired with the `escalate` feature deletion (brief D9); setting it is ignored",
	},
	{
		Env: "QUINE_SMART_API_KEY", Decision: "brief D9", When: "2026-07-03",
		Note: "retired with the `escalate` feature deletion (brief D9); setting it is ignored",
	},
	{
		Env: "QUINE_RUNTIME_SURFACE_BACKEND", Decision: "2026-06-control-surface-fuse-only", When: "2026-06-27",
		Note: "retired by the FUSE-only control-surface collapse ([`runtime-surface.md`](./runtime-surface.md), [`../decisions/2026-06-control-surface-fuse-only.md`](../decisions/2026-06-control-surface-fuse-only.md)). The public surface is now unconditionally FUSE; the `legacy` symlink backend and this selector are gone. Setting it is ignored (not propagated to child env), not an error",
	},
	{
		Env: "QUINE_TURN_EXHAUSTION_POLICY", Decision: "brief D8", When: "2026-07-03",
		Note: "retired by the registry cleanup ([`registry-design-brief.md` D8](../../theory/views/runtime-capability/registry-design-brief.md)): `hard_fail` was the only legal value, so the enum carried no information; hard-fail at budget zero is now the implicit behavior; setting it is ignored, not an error",
	},
	{
		Env: "QUINE_SH_TIMEOUT", Decision: "brief D8", When: "2026-07-03", Replacement: EnvShDefaultTimeout,
		Note: "retired by the registry cleanup ([`registry-design-brief.md` D8](../../theory/views/runtime-capability/registry-design-brief.md)): pre-2026-04-19 name of the sh timeout, since then consumed by nothing (declared only for defensive child-env stripping, itself now removed); setting it is ignored, not an error; active env is `QUINE_SH_DEFAULT_TIMEOUT_SECONDS`",
	},
	{
		Env: "QUINE_STALL_THRESHOLD", Decision: "brief D9", When: "2026-07-03",
		Note: "retired with the `escalate` feature deletion ([`registry-design-brief.md` D9](../../theory/views/runtime-capability/registry-design-brief.md)): stall detection's only consumer was the escalate hint; setting it is ignored",
	},
	{
		Env: "QUINE_WISDOM_*", Decision: "brief D5", When: "2026-07-03", Replacement: EnvInitialUserMessage,
		Note: "retired with the wisdom mechanism deletion ([`registry-design-brief.md` D5](../../theory/views/runtime-capability/registry-design-brief.md)): channel tri-partition (argv = mission, env = capability position, `context/` = cognition carry-forward); operator note injections migrated to `QUINE_INITIAL_USER_MESSAGE`; setting it is ignored, not an error",
	},
}

// ExternalLabelKnob describes an env name that authored surfaces use as a
// harness/profile label but the quine binary deliberately never reads (the
// brief's anticipated external-label default kind). It cannot enter the live
// registry today without an envnames.go constant (the bijection test), so it
// lives here as documentation data; flagged in the T1.2 report as a T1.1
// follow-up if the owner wants it registry-resident.
type ExternalLabelKnob struct {
	Env   string `json:"env"`
	Scope string `json:"scope"`
	Note  string `json:"note"` // env-controls.md Notes cell, verbatim markdown
}

// ExternalLabels lists the deliberate non-knob label envs env-controls.md
// documents alongside the live matrix.
var ExternalLabels = []ExternalLabelKnob{
	{
		Env:   "QUINE_PROVIDER",
		Scope: "profile / harness label only",
		Note:  "not read by `quine`; use only for bench/harness/profile routing",
	},
}

// RetiredKnobByEnv returns the retired-knob entry for an env name.
func RetiredKnobByEnv(env string) (RetiredKnob, bool) {
	for _, k := range RetiredRegistry {
		if k.Env == env {
			return k, true
		}
	}
	return RetiredKnob{}, false
}
