---
paper_id: unfakeable-selection
claim_bundle: unfakeable-selection
venue: arxiv
format-envelope: arxiv-preprint
track: theory
cluster: ecological-selection
status: planning
state_mode: planning-skeleton
role: unfakeable-selection-theory-backbone
family: open-ended-ecology
publication_class: theory-candidate
maturity_target: theory-anchor
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: WORKSPACE.md
theory-objects: [narrative-claim, carrier, exec-barrier, delayed-addressability]
gaps: [function-vs-narrative, persistence-as-unfakeability]
provenance: ../../../theory/references/multi-lens-collision-2026-06-29.md
view-affinity: [reflexive, ecology]
---

# A Separation Theory of Un-Fakeable Selection

> **Theoretical Foundation & Archival Anchor.** This dossier establishes the
> unified **separation theory of un-fakeable selection**: demonstrating how five
> foundational lenses (Rice's Theorem in TCS, Bennett's Logical Depth, Mechanism
> Design Single-Crossing, Non-Equilibrium Thermodynamic Work $\Delta W$, and
> Fisher Information) compose into a single rigorous mathematical operator for
> un-fakeable evolutionary and coordination selection in computational agent ecologies.
>
> As a pure theoretical anchor, this framework is venue-independent: it formalizes the
> structural necessity of physical execution barriers (`exec` as $\Delta \to \infty$)
> and provides the mathematical foundation for un-fakeable coordination across all
> language-mediated agent ecosystems.

**Provenance:** 5-lens convergence (Finding 1) in
[`multi-lens-collision-2026-06-29`](../../../theory/references/multi-lens-collision-2026-06-29.md).
**Feeds:** the **exec-unfakeable** direction
([`substrate-coincidence/DIRECTIONS.md`](../substrate-coincidence/DIRECTIONS.md)
#3; incubator writeup
[`aaai-c1-exec-unfakeable-writeup`](../../incubator/aaai-c1-exec-unfakeable-writeup.md)).
This dossier owns the **theory backbone** of un-fakeable selection.


## Current Question

When can a selection or coordination signal in a language-mediated agent
ecology **not** be satisfied by narration — and is there a single graded
quantity, measurable on the runtime, that decides it?

## Thesis

Un-fakeability is not a binary property of an action but a **separation
condition on a cost/depth gradient**, and the same condition is approached from
five independent directions that compose into one operator. A signal resists
narrative collapse iff the cost gap between faking and doing satisfies
single-crossing (Δ = cost(fake) − cost(real) large); the un-fakeable residue is
the **logically deep**, non-shortcuttable, behaviorally-load-bearing component
of an artifact; and `exec` is simply the **Δ→∞ corner** of that continuum — the
unique incentive-compatible mechanism when a substrate has no transfers and no
future. Specification gaming is then "narrative decoupling on the low-Δ region
of the gradient," and the cure is structural (push the load-bearing signal to
high Δ), not a better judge.

## The five composing results

From [Finding 1](../../../theory/references/multi-lens-collision-2026-06-29.md):

| Source | Role in the operator |
|--------|----------------------|
| **Rice's theorem** (TCS) | No *syntactic* function/narrative decider can exist; the operator is necessarily behavioral and **one-sided** (r.e. not co-r.e.): a carrier is confirmable, pure narrative only suspectable. |
| **Logical depth** (Bennett) | The graded magnitude — narration shallow, carrier deep; NOT Kolmogorov (random = incompressible = fakeable). Unifies persistence + function into one axis. |
| **Δ single-crossing** (mechanism design) | The equilibrium condition for separation; reframes [`unfakeability-gradient`](../../core/mechanisms/unfakeability-gradient.md) as measuring level when it should measure slope. |
| **No-IC mechanism** (mechanism design) | Why the carrier is forced: no transfers + no future ⇒ the revelation principle fails ⇒ a carrier-restricted message space is the unique incentive-compatible mechanism. |
| **Signed ΔW** (non-equilibrium physics) / **Fisher information** (ML) | Two measurable estimators: ΔW = free-energy transfer to the consumer (adversarial falsehood ΔW>0, worse than neutral); Fisher = ∂(next agent's action)/∂(content) with files as weights. |

## Core Claim Boundary

This dossier establishes the **5-lens separation theory of un-fakeable
selection** and its **multidisciplinary measurement framework** (Δ / depth /
Fisher / ΔW as convergent estimators of one underlying quantity), grounded in the
`narrative-claim` × `carrier` spine. It proves the Rice impossibility of syntactic
function deciders and the No-IC uniqueness of carrier-restricted message spaces.

## In Scope

- the separation condition (single-crossing in Δ) and its comparative statics
- logical depth as the graded un-fakeability measure, contra Kolmogorov/MDL
- the No-IC uniqueness argument (no transfers, no future → carrier message-space)
- convergent estimators (Fisher info, signed ΔW) and their agreement/divergence
- the Goodhart-ladder connection (causal mode cured by exec → regressional
  residue = the symmetry trap; see sibling `symmetry-trap-phase`)
- the hollow-payload resolution as a *prediction* of the theory (covert channel
  = lower-Δ carrier, still on the continuum)

## Out Of Scope

- Claiming model-class-independent universality before cross-architecture empirical verification
- Conflating theoretical separation conditions with specific single-venue formatting requirements

## Relation To Siblings

- [`substrate-coincidence/DIRECTIONS.md`](../substrate-coincidence/DIRECTIONS.md) #3 (exec-unfakeable) provides the execution-grounded physical layer that implements this theory.
- [`symmetry-trap-phase`](../../candidates/symmetry-trap-phase.md) (candidate) formalizes the *regressional residue* predicted after the Δ→∞ barrier closes the causal mode.
- [`narration-invariance-eval`](../narration-invariance-eval/README.md) provides invariant risk minimization as one observable estimator.

## Current State

- provenance / full derivation:
  [`multi-lens-collision-2026-06-29`](../../../theory/references/multi-lens-collision-2026-06-29.md)
  Finding 1 + the per-lens candidate entries
- working notes: [`WORKSPACE.md`](WORKSPACE.md)
- seeds planted on `function-vs-narrative`, `persistence-as-unfakeability`,
  `narrative-claim`, `carrier-selection`

## Writing Posture

- Treat the five lenses as convergent mathematical derivations of a single un-fakeable selection operator.
- Ground the theoretical architecture in measurable computational observables (Δ, depth, Fisher info, signed ΔW).
- Lead with the constructive operator, using Rice's theorem to establish the necessity of behavioral rather than syntactic evaluation.

## Experiment Feeds

- `none-yet` - none - not-for-paper-yet - the dorm-04/dorm-11 Δ two-arm and the cfi-01 IRM/capability-dependence check are the intended feeds but are unrun for this dossier and not yet reciprocally declared; bind once the experiments add `## Paper Feeds` rows (see Next Useful Work)

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - the two-arm Δ experiment is designed but unrun; all falsifiers are proposed, none executed

## Strongest Planned Contribution

- a single separation operator with four agreeing estimators (Δ, depth, Fisher,
  ΔW) and a stated impossibility (Rice) that explains why the operator must be
  behavioral
- the No-IC uniqueness result: on a transfer-free, future-free substrate the
  carrier-restricted message space is the *only* incentive-compatible mechanism
- a comparative static that no single-ecology run can currently make: collapse
  probability is monotone in Δ and **invariant to selection intensity**
- a principled dissolution of the hollow-payload objection as a theory
  prediction rather than a patch

## Blocking Gaps

- external novelty boundary (Rice/Bennett/von Neumann applied to LLM ecologies;
  Δ-single-crossing vs the spec-gaming canon) — UNADJUDICATED, route to deep
  research before any claim
- the Δ-vs-selection-intensity two-arm must be run (no executed falsifier yet)
- whether Fisher and signed-ΔW are the *same* estimator up to sign, or genuinely
  separable measurements, needs a worked case
- the paradigm proposals (`unfakeability-gradient` level→slope; `carrier` =
  logically-deep object) need a human-directed decision

## Next Useful Work

- source-checked novelty review (the gating step)
- design + run the dorm-04/dorm-11 Δ two-arm (Δ≈∞ carrier vs Δ≈0 narratable
  field, selection intensity fixed)
- a single worked example computing Δ, a depth proxy, and Fisher on the *same*
  artifact to show the estimators agree
- draft the operator definition + the Rice/No-IC results from the provenance file
