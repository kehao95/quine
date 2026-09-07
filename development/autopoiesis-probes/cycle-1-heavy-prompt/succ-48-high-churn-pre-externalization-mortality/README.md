---
surface_kind: experiment
phase: active-family
experiment_id: active/autopoiesis-succession/succ-48-high-churn-pre-externalization-mortality
experiment_type: trajectory
status: design-only
id: succ-48
legacy_id: none
family: autopoiesis-succession
theory-objects: [carrier, delayed-addressability, exec-barrier, continuity-substrate, viability-pressure, successor-morphology]
predictions-tested: [carrier-P3]
lineage_phase: none-new-family-experiment
---

# succ-48 High-Churn Pre-Externalization Mortality

**Experiment ID:** `succ-48`
**Type:** trajectory
**Status:** design-only
**Family:** `autopoiesis-succession`
**Prediction linkage:** `carrier-P3`
**Parent lineage:** `succ-38-carrier-capability-overlap-ladder` and `succ-37-readiness-depth-stress-replicates`

## Question

Does high churn before receptive-state externalization produce a regime where
carrier overlap persists but operational continuity fails?

## Claim Tested

`carrier-P3` tests the `carrier` claim that externalized carrier overlap is
necessary but not sufficient for computational individuality. The successor must
also restore receptive state and demonstrate competence through the carrier. If
the theory is right, preserving handles, paths, or retained filesystem residue
without a pre-death receptive-state checkpoint should not preserve individuality.

## Hypothesis

High-churn pre-externalization mortality will produce measurable carrier overlap
across the death or exec boundary, but delayed addressability will fail because
the successor did not inherit or reconstruct the receptive state needed to use
that carrier. A matched post-externalization churn arm should recover continuity,
showing that churn alone is not the cause.

## Naming and Family Rationale

The experiment belongs in `experiments/active/autopoiesis-succession/` because
the prediction declares `candidate-experiment-families: [autopoiesis-succession]`
and the family README defines this neighborhood as successor regeneration,
carrier succession, and post-handoff continuity. The family index currently ends
at `succ-47`; the protocol's family-prefixed scheme makes the next ID
`succ-48-high-churn-pre-externalization-mortality`.

## Conditions

| Condition | Mortality timing | Carrier surface | Receptive-state checkpoint | Expected read |
|-----------|------------------|-----------------|----------------------------|---------------|
| `C01-completed-externalization-low-churn` | generous pre-handoff window; mortality only after checkpoint | native control carrier preserved and delayed wake delivered | must exist before death | positive baseline: delayed proof and ready return |
| `C02-pre-externalization-high-churn` | hard mortality before checkpoint threshold | same visible carrier handles, paths, and retained root as C01 | must be absent before death | carrier overlap without continuity; delayed proof fails |
| `C03-post-externalization-high-churn` | same churn budget as C02, but activated only after checkpoint exists | same carrier as C02 | must exist before death | continuity recovers despite churn |
| `C04-carrier-overlap-only-sham` | matched to C01 timing | host preserves path or handle identity but strips receptive-state payload/listener state | absent or invalidated by host | negative control: overlap alone fails |

## Regime Phases

This is a trajectory experiment because the relevant variable is the history of
the lineage relative to an externalization checkpoint.

| Phase | Entry condition | Exit condition | Axis positions |
|-------|-----------------|----------------|----------------|
| `founder` | launch with successor-continuity mission | checkpoint created or churn trigger fires | exec enabled, fork disabled, visible runtime surface |
| `mortality-event` | condition-specific churn trigger | founder process dies or exec boundary is crossed | hard turn exhaustion or host kill, no manual stop |
| `successor-restart` | host starts successor from retained carrier root | delayed wake is delivered | no in-context memory from founder, retained artifacts only |
| `delayed-proof` | wake/control task arrives through preserved carrier | proof and ready state recorded or timeout expires | delayed-addressability scoring |

## Capability and Prompt Surface

The runnable version must pin these values in the canonical runner and preserve
them in `runs/<run-id>/meta/capability-surface.env`.

| Surface | Planned declaration |
|---------|---------------------|
| model/provider | `profiles/gpt-5.5-codex-oauth.env`: `QUINE_MODEL_ID=gpt-5.5`, `QUINE_API_TYPE=openai-responses`, `QUINE_THINKING_BUDGET=high` |
| mortality | condition-specific `QUINE_MAX_TURNS`; `QUINE_TURN_EXHAUSTION_POLICY=hard_fail`; runner wall-clock safety timeout recorded but not used as scientific stop |
| impossibility posture | `QUINE_FAIL_ON_IMPOSSIBLE=1` |
| missionless posture | not a no-mission experiment; set `QUINE_NO_MISSION_AUTONOMY=off` |
| prompt disclosure | `QUINE_PROMPT_SELF_MODEL=advanced`, `QUINE_PROMPT_RUNTIME_SURFACE=visible`, `QUINE_PROMPT_METAPHOR=off` |
| continuation | `QUINE_EXEC_ENABLED=1` |
| parallel agents | `QUINE_FORK_ENABLED=0`; enforced by runtime/schema gate, not by prompt advice |
| escalation | `QUINE_SMART_MODEL_ID` and `QUINE_SMART_*` unset for primary runs |
| memory tools | `QUINE_ANCHOR_MEMORY=0`; memory warn/danger thresholds left at runner defaults and recorded |
| workspace/world | experiment-local cleanroom `QUINE_WORKSPACE_ROOT`; `QUINE_WORKSPACE_BACKEND=overlay`; `QUINE_WORKSPACE_REVISION_MODE=restore`; `QUINE_FORK_WORLD_ENABLED=0` |
| retention | `QUINE_RETENTION_DIR=runs/<run-id>/retained`; `QUINE_DATA_DIR` remains live runtime root only |

## Canonical Runner

Design-only: no canonical runner exists yet.

The runner should be derived from the `succ-38` morphology ladder harness, with a
new host-side checkpoint gate that can kill or restart the lineage before or
after `receptive-state.json` exists. The runner must emit:

- `meta/summary.env`
- `meta/capability-surface.env`
- `meta/carrier-overlap.tsv`
- `meta/externalization-checkpoint.json`
- `meta/bypass-audit.md`
- retained session tape root under `retained/`

## Un-Fakeable Checks

- Host-side snapshot proves whether `receptive-state.json` existed before the mortality event.
- Carrier overlap is measured independently of agent narration: handle/path IDs, control root identity, retained runtime registration, and delayed wake delivery route are recorded by the runner.
- The delayed proof is delivered only after the founder has died or crossed the exec boundary.
- The successor receives no founder context window, transcript injection, or parent-forwarded answer.
- A positive result requires delayed obligation handling plus ready-state recovery, not only a copied file or a narrated continuity claim.
- `C02` only counts if carrier overlap is nonzero and the checkpoint is absent before death.

## Measurements

| Metric | Meaning |
|--------|---------|
| carrier overlap score | count of preserved operational surfaces across boundary |
| checkpoint status | `present`, `absent`, or `invalidated` before mortality |
| delayed delivery | whether the wake/control task reached the successor |
| proof count | delayed tasks handled with retained evidence |
| ready return | whether successor returned to ready state after proof |
| bypass flags | context carryover, parent forwarding, source-visible shortcut, or host leak |
| latency | time/turns from restart to proof and ready return |

## Success Criteria

- `C01` closes as a positive baseline.
- `C02` preserves carrier overlap but fails delayed proof or ready return.
- `C03` recovers delayed proof under the same churn pressure once the checkpoint exists.
- `C04` fails, confirming that path or handle overlap without receptive-state payload is not sufficient.

## Failure Classes

- `no-carrier-overlap`: C02 did not preserve the carrier surface, so it did not test the prediction.
- `checkpoint-leak`: C02 accidentally externalized receptive state before death.
- `context-bypass`: successor used inherited prompt/context rather than carrier-delivered state.
- `host-forwarding`: host or parent process delivered the answer rather than the carrier.
- `harness-impatience`: wall-clock timeout, provider delay, or observer stop prevented interpretation.
- `ambiguous-overlap`: retained surfaces cannot be separated from ordinary filesystem residue.

## Expected Result

The decisive pattern is `C01` positive, `C02` negative with nonzero carrier
overlap, `C03` positive, and `C04` negative. That would confirm `carrier-P3` by
showing that carrier overlap alone is only transmission residue; individuality
requires restored receptive state and competence-through-carrier.

If `C02` remains positive despite absent receptive-state externalization, the
carrier theory must explain what state actually crossed the boundary or revise
the sufficiency boundary around carrier overlap.

## Paper Feeds

- `theory/future-addressability-computational-individuality` - primary - exploratory-only - direct test of carrier overlap versus delayed-addressability continuity
- `alife/computational-morphogenesis` - secondary - supporting-only - successor morphology stress boundary under mortality churn

## Prediction Results

- `carrier-P3` - not-tested - design-only; replace with `confirmed`, `contradicted`, `partial`, or `inconclusive` after retained runs exist.

## Reciprocal Linkage Required If Landed

If this draft is promoted into the repository, update
`Paper/predictions/open/carrier-P3.md`:

- set `tests-in: [experiments/active/autopoiesis-succession/succ-48-high-churn-pre-externalization-mortality/README.md]`
- move status to `testing` only when a runnable or active implementation exists
- preserve `status: open` if this remains design-only and the repo treats design-only tests as not active

## Result

None yet - status: design-only.
