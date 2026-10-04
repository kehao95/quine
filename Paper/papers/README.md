# Papers — Family-Grouped Dossier Index

Public snapshot: see [scope, path normalization, and dependencies](PUBLIC-RELEASE.md).

Paper dossiers live flat at `papers/<id>/`. The 2026-06-29 three-layer-flatten
migration (see
[`../_design/migrations/three-layer-flatten/closeout-2026-06-29.md`](../_design/migrations/three-layer-flatten/closeout-2026-06-29.md))
retired the earlier `papers/<family>/<id>/` nesting; this index **groups them by
publication family** for navigation instead. A family is a
**publication-management unit** — sibling role split, anti-duplication boundary,
sequencing (see
[`../paper-families.md`](../paper-families.md)). It holds **no theory content**: a
dossier's scientific affiliation is carried by **frontmatter** (`theory-objects`,
`gaps`, `view-affinity`), **not by its family grouping**. `family ≠ view` — the
family a dossier is grouped under here is the publication sibling unit; the
`view-affinity` frontmatter is the theory lens, and they need not align. To see the
science a dossier projects, read its frontmatter and follow the object / gap / view
IDs into [`../../theory/README.md`](../../theory/README.md).

This index is **index-like**: it routes to dossiers and records affiliation,
theory-objects, and status at a glance. It holds **no mutable per-paper prose** —
thesis, scope, evidence backbone, plan, and draft state live in each dossier's own
`README.md`.

## Strict `papers/` rule

A directory earns a home in `papers/` **only if** it has a `manuscript-plan.md`,
OR a `manuscript.md`, OR a `plans/` subdirectory with a non-trivial plan file, OR
is a multi-configuration dossier whose subdirectories each meet that bar (see
[`../contracts.md`](../contracts.md) § Strict `papers/` rule). A README-only
skeleton (a `WORKSPACE.md` alone does not count) **retreats to
[`../candidates/`](../candidates/README.md)** as `status: held`; its original
path is preserved as a redirect pointer under
[`_retreated/`](_retreated/) (`kind: retreated-pointer`) that links to the
candidate. Re-promotion restores the dossier from the candidate's
`dossier-provenance.frontmatter-snapshot`.

### Retreated dossiers (now in `../candidates/`)

These skeleton-only homes retreated to [`../candidates/`](../candidates/README.md)
on 2026-06-29 (pointers in [`_retreated/`](_retreated/)):
<!-- BEGIN GENERATED: papers-index:retreated (generate-paper-projections.py — edit frontmatter, not this table) -->
`audit-attractor`, `communication-emergence`, `interregnum-corrigibility`, `reproductive-specialization`, `self-immunity`, `stigmergic-trace-causality`, `symmetry-trap-phase`.
<!-- END GENERATED: papers-index:retreated -->


## `runtime-substrate` (12)

The process / world / time runtime substrate chain (runtime-ontology lineage).

<!-- BEGIN GENERATED: papers-index:runtime-substrate (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`agent-mechanisms-as-os-primitives`](agent-mechanisms-as-os-primitives/README.md) | `substrate` | `future-task-surface`, `exec-barrier`, `coordination-closure`, `carrier` | `planning` |
| [`agentic-alife-workshop`](agentic-alife-workshop/README.md) | `substrate` | `exec-barrier`, `viability-pressure`, `trace-inheritance`, `coordination-closure`, `protocol-emergence`, `carrier`, `delayed-addressability`, `future-task-surface`, `narrative-claim` | `accepted` |
| [`institutional-action-licensing`](institutional-action-licensing/README.md) | `substrate` | `two-latch-gate`, `deixis-release` | `planning` |
| [`process-as-fs`](process-as-fs/README.md) | `substrate` | `future-task-surface`, `exec-barrier`, `delayed-addressability`, `coordination-closure`, `trace-inheritance` | `planning` |
| [`quine-arxiv`](quine-arxiv/README.md) | `substrate` | `future-task-surface`, `exec-barrier`, `carrier`, `trace-inheritance` | `published` |
| [`quine-hotos`](quine-hotos/README.md) | `substrate` | `carrier`, `exec-barrier`, `trace-inheritance`, `delayed-addressability`, `coordination-closure`, `future-task-surface`, `structural-pull` | `planning` |
| [`quine-notes-series`](quine-notes-series/README.md) | `substrate` | `future-task-surface`, `narrative-claim` | `planning` |
| [`revisable-time`](revisable-time/README.md) | `substrate` | `future-task-surface`, `exec-barrier`, `delayed-addressability`, `trace-inheritance`, `narrative-claim` | `planning` |
| [`structural-elicitation`](structural-elicitation/README.md) | `substrate` | `structural-pull`, `narrative-claim` | `accepted` |
| [`structure-originates-direction`](structure-originates-direction/README.md) | `—` | `—` | `planning` |
| [`structure-two-latch-action-gate`](structure-two-latch-action-gate/README.md) | `—` | `—` | `planning` |
| [`subjective-reality`](subjective-reality/README.md) | `substrate` | `future-task-surface`, `exec-barrier` | `planning` |
| [`systems-track-b`](systems-track-b/README.md) | `substrate` | `future-task-surface`, `exec-barrier`, `coordination-closure` | `planning` |
<!-- END GENERATED: papers-index:runtime-substrate -->


## `successor-continuity` (9)

Bounded self-production, renewal, and lineage continuity (autopoiesis-renewal lineage).

<!-- BEGIN GENERATED: papers-index:successor-continuity (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`autoconstructive-successor-morphology`](autoconstructive-successor-morphology/README.md) | `continuity` | `carrier`, `exec-barrier`, `delayed-addressability`, `viability-pressure`, `future-task-surface` | `planning` |
| [`boundary-paper`](boundary-paper/README.md) | `continuity` | `carrier`, `exec-barrier`, `delayed-addressability`, `trace-inheritance`, `viability-pressure`, `future-task-surface` | `planning` |
| [`computational-morphogenesis`](computational-morphogenesis/README.md) | `continuity` | `carrier`, `exec-barrier`, `delayed-addressability`, `future-task-surface`, `viability-pressure` | `planning` |
| [`emergent-self-reproduction`](emergent-self-reproduction/README.md) | `continuity` | `emergent-succession`, `self-description-resource`, `identity-criterion` | `drafting` |
| [`future-addressability`](future-addressability/README.md) | `continuity` | `carrier`, `exec-barrier`, `delayed-addressability`, `future-task-surface`, `trace-inheritance` | `planning` |
| [`generative-self-replacement`](generative-self-replacement/README.md) | `continuity` | `carrier`, `exec-barrier`, `delayed-addressability`, `future-task-surface`, `viability-pressure` | `drafting` |
| [`iclr-regime-decomposition`](iclr-regime-decomposition/README.md) | `continuity` | `exec-barrier`, `carrier`, `delayed-addressability`, `unsampled-function`, `viability-as-distortion-measure`, `carrier-membership-heritability-predictor`, `narrative-claim`, `viability-pressure`, `future-task-surface`, `reproduced-object` | `drafting` |
| [`lifecycle-bottleneck-reproduction`](lifecycle-bottleneck-reproduction/README.md) | `continuity` | `carrier`, `exec-barrier`, `delayed-addressability`, `trace-inheritance`, `viability-pressure` | `planning` |
<!-- END GENERATED: papers-index:successor-continuity -->


## `eval-primitives` (1)

Operator-projection evaluation primitives for mainstream AI/ML venues.

<!-- BEGIN GENERATED: papers-index:eval-primitives (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`narration-invariance-eval`](narration-invariance-eval/README.md) | `reflexive` | `narrative-claim` | `planning` |
<!-- END GENERATED: papers-index:eval-primitives -->


## `viability` (5)

Bounded life-like organization under separable viability pressures (viability-habitat lineage).

<!-- BEGIN GENERATED: papers-index:viability (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`alife-priority-arxiv`](alife-priority-arxiv/README.md) | `viability` | `viability-pressure`, `carrier`, `exec-barrier`, `coordination-closure`, `protocol-emergence`, `trace-inheritance` | `maintained` |
| [`quine-lifelike`](quine-lifelike/README.md) | `viability` | `viability-pressure`, `carrier`, `exec-barrier`, `trace-inheritance` | `rejected` |
| [`quine-process-habitat`](quine-process-habitat/README.md) | `viability` | `viability-pressure`, `carrier`, `trace-inheritance`, `future-task-surface`, `exec-barrier`, `coordination-closure`, `protocol-emergence` | `planning` |
| [`routing-surface-coevolution`](routing-surface-coevolution/README.md) | `viability` | `viability-pressure`, `trace-inheritance`, `stigmergy-as-carrier` | `drafting` |
<!-- END GENERATED: papers-index:viability -->


## `social-protocol` (3)

Coordination closure, protocol emergence, and channel semantics (social-behavior lineage).

<!-- BEGIN GENERATED: papers-index:social-protocol (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`material-channels`](material-channels/README.md) | `social` | `coordination-closure`, `protocol-emergence`, `carrier` | `planning` |
| [`minimal-perceptual-prerequisites`](minimal-perceptual-prerequisites/README.md) | `social` | `coordination-closure`, `protocol-emergence`, `trace-inheritance` | `accepted` |
| [`minimal-runtime-protocols`](minimal-runtime-protocols/README.md) | `social` | `protocol-emergence`, `coordination-closure`, `carrier` | `planning` |
<!-- END GENERATED: papers-index:social-protocol -->


## `open-ended-ecology` (4)

Missionless, affordance-rich ecologies and stigmergic trace causality (open-ended-ecology lineage).

<!-- BEGIN GENERATED: papers-index:open-ended-ecology (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`substrate-coincidence`](substrate-coincidence/README.md) | `runtime-capability` | `substrate-coincidence`, `reflexive-substrate-loop`, `narrative-claim` | `drafting` |
| [`unfakeable-selection`](unfakeable-selection/README.md) | `reflexive`, `ecology` | `narrative-claim`, `carrier`, `exec-barrier`, `delayed-addressability` | `planning` |
<!-- END GENERATED: papers-index:open-ended-ecology -->


## `reflexive-repository` (7)

Repository self-maintenance, structural pull, and public explanation (repository-self-model + public-explanation lineage).

<!-- BEGIN GENERATED: papers-index:reflexive-repository (generate-paper-projections.py — edit frontmatter, not this table) -->
| dossier | view-affinity | theory-objects | status |
|---------|---------------|----------------|--------|
| [`computational-neoplasm`](computational-neoplasm/README.md) | `reflexive` | `computational-neoplasm`, `self-immunity`, `audit-attractor`, `authority-override`, `narrative-claim` | `accepted` |
| [`contentsquare-cli-one-day`](contentsquare-cli-one-day/README.md) | `reflexive` | `structural-pull`, `narrative-claim` | `published` |
| [`harness-engineering-is-architectural-amnesia`](harness-engineering-is-architectural-amnesia/README.md) | `reflexive` | `carrier`, `exec-barrier`, `protocol-emergence`, `structural-pull`, `narrative-claim` | `published` |
| [`open-sourcing-a-research-workbench`](open-sourcing-a-research-workbench/README.md) | `reflexive` | `carrier`, `trace-inheritance`, `delayed-addressability`, `protocol-emergence`, `structural-pull`, `narrative-claim` | `maintained` |
| [`structure-grows`](structure-grows/README.md) | `reflexive` | `structural-pull`, `narrative-claim`, `viability-pressure` | `published` |
| [`the-autopoietic-repository`](the-autopoietic-repository/README.md) | `reflexive` | `structural-pull`, `narrative-claim` | `published` |
| [`tumor-in-the-repository-essay`](tumor-in-the-repository-essay/README.md) | `reflexive` | `computational-neoplasm`, `self-immunity`, `audit-attractor`, `authority-override`, `narrative-claim` | `published` |
<!-- END GENERATED: papers-index:reflexive-repository -->

## Routing Rule

If you are looking for the **scientific object** a dossier projects, do not read it
from this index — read the dossier's frontmatter and follow the `theory-objects` /
`gaps` / `view-affinity` IDs into [`../../theory/README.md`](../../theory/README.md). This index
only routes to dossiers under their publication family and shows declared
affiliation and status. A dossier's directory name may differ from its `paper_id`
(e.g. `boundary-paper/` → `quine-cognitive-autopoiesis`); the link resolves to the
actual directory.

<!-- BEGIN GENERATED: papers-index:count (generate-paper-projections.py — edit frontmatter, not this table) -->
*(38 dossiers across 7 families. `terminal-bench-2-love-letter` stays in `essays/papers/`.)*
<!-- END GENERATED: papers-index:count -->
