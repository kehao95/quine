# XX: Runtime Capability Topology (Jun 22)

## Arc

This period began from a diagnosis of the gap between the theory generation
loop (seed → hypothesis → theory object growth, already working) and the full
research loop (experiment design → execution → evidence mining → theory
writeback → manuscript prose, never validated end-to-end).

### Loop test (diagnosis)

Two agents were dispatched to walk the research loop on real targets:
- Agent A: stigmergy-as-carrier (evidence → theory → prose)
- Agent B: trigger-structure seed (seed → experiment design → theory)

Result: agents did not hard-block but produced uncontrollable-quality output.
Both independently reported the structure pulling toward audit/deepening
rather than execution/generation — the audit-attractor pattern documented in
`theory/views/reflexive/deep/cases/01`.

### Root cause

The runtime's capability space is a flat bag of ~90 env vars with no named
compositional structure. "Mortality" as a concept does not exist in code —
it is the emergent property of 4+ env vars. Agents cannot determine whether
a proposed ecological condition is expressible because the compositional
relationships are implicit in experiment runners, not in any named surface.

### Round 1: latent space discovery

Four parallel agents produced `theory/views/runtime-capability/`:
- **env-archaeology** (716 lines): code-level trace of all 88 env vars
- **regime-atlas** (365 lines): 9 named regimes from 89 experiment-regime associations
- **philosophy-coherence** (519 lines): 5 tensions between env-config and everything-is-file
- **projection-axes** (921 lines): 7 named axes (Mortality, Generation Continuity, Perceptual Disclosure, Self-Relation, Ecological Directive, Social Density, Enaction Fidelity)

### Round 2: refactoring preparation

Three targeted agents deepened the investigation:
- **axis-code-mapping** (804 lines): per-axis Go code ownership, refactoring grain, readiness ranking
- **philosophy-of-capability** (1018 lines): per-axis philosophical commitments, 5 cross-axis conflicts, 7 blind spots
- **axis-asymmetry** (792 lines): dimensional depth profile, dark matter, phantom demand, priority ranking

Key finding: Axis 4 (Self-Relation) is the largest dark matter — fully implemented (`anchor_memory` + `mark` + `fold` + `self_source_code`), universally disabled in experiments, zero theory demand.

### Coherence entropy reduction

A focused audit identified 16 coherence issues (broken pointers, stale citations, doc-code mismatches, reciprocity gaps). 13 were fixed in `7121cd58` (COH-01 to COH-14 bounded issues; COH-16 was a negative finding); the remaining 3 structural issues were addressed in `61a717e9` (COH-07 legacy ID resolution across 9 theory objects, COH-14 residual null-04 aspirational trim, COH-15 cross-family reciprocity judgment).

**Note:** `61a717e9` was subsequently reverted (`21d9af09`) by a concurrent work stream. The COH-07 legacy ID debt (9 theory objects with unresolved `6.NN`/`6M.NN`/`6R.NN` citations) is therefore **still open**. The COH-15 reciprocity fixes may also have been undone. This revert should be investigated and the fixes re-applied if the revert was unintentional.

### Loop re-test (validation)

Two agents re-walked the same loop with structural support (access to the
runtime-capability view). Result: **structural naming measurably improved
loop quality.**

| Metric | Without structure | With structure |
|--------|-------------------|----------------|
| Agent A Stage 3 (mine) | blocked (DVC wall) | partial — P7G run actually mined from materialized artifacts |
| Agent A Stage 5 (prose) | blocked (no home) | partial — full grounded paragraph written |
| Agent B Stage 3 (design) | blocked (capability floor) | partial — full README skeleton + runnable two-run approximation |
| Blocker naming | ad hoc | axis vocabulary |
| Workaround discovery | invented from scratch | read from axis-asymmetry.md |

### Four blocker fixes

Four physical-layer blockers identified by the loop test were fixed (`5591ee71`):
1. DVC evidence pull affordance documented in PROTOCOL.md §2b
2. Stigmergic-trace-causality seed dossier created as prose home
3. Trajectory experiment contract added to PROTOCOL.md §9
4. Regime atlas precision fixed (trce-08/09 added)

## Commits

- `0264e151` — Paper: map runtime capability topology (Round 1 + Round 2 research notes)
- `7121cd58` — Close runtime-capability coherence debt (COH-01 to COH-14)
- `61a717e9` — Close structural coherence debt (COH-07/14/15) [REVERTED by 21d9af09]
- `5591ee71` — Close four loop-test blockers

## Open debt carried forward

Canonical detailed queue: `Paper/theory/views/runtime-capability/axis-asymmetry.md` § Carry-Forward Pain Points.

- **COH-07**: legacy ID resolution across 9 theory objects was reverted; the debt is back
- **COH-15**: cross-family reciprocity for affordance-null-ecology may have been undone by the same revert
- **Axis/code gap**: axes are named but code structure (`ToolGates`, `Config`) does not yet reflect them
- **File-addressability gap**: config layer is still env-only, not file-operable by the agent
- **Theory/runtime decision queue**: Axis 4 dark matter, Axis 6 opacity, static axis-space trajectory, and capability scope-condition absorption now live in the canonical runtime-capability queue
