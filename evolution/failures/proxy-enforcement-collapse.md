---
id: proxy-enforcement-collapse
type: failure-mode
status: repaired
applied: true
class: mechanism-design
symptom: a blocking check verifies token/shape presence instead of ground truth; the checked population converges on template replication and the check manufactures false confidence
evidence: [f7bffc880, 429fdd600, 32117633d, 699f7bdf2, ae5702af5, e0596564c, 89e61f21e, 17c2650e9, cc3842e2f, 6f8e0ecf2, 5c11cf31d, d4b5c96d8, f3fbe5e76, d7d66e17c, a21f45351, 118481d83]
candidate-solution: requirement re-derivation before any check is created or repaired; repair verbs delete > rebind > amend-existing; one authoring-guidance amendment; explicitly no meta-checks or watcher structures
owner: lab
---

# proxy-enforcement-collapse

## Symptom

A check wants a latent property (honest reasoning, real linkage, maintained
content) but verifies a proxy: tokens, headings, or fields present in the
checked text. Once the check is blocking and applied to a population of
surfaces, the population converges on the cheapest satisfying strategy —
replicate the shape that passes. The proxy decouples from the latent property,
the population turns uniform ("AI-generated" tone is the felt symptom), and the
check now manufactures false confidence: green means the tokens are present,
which no longer means the property holds.

**Detection heuristic.** A check is in this class iff its pass condition can be
satisfied by editing only the checked text. Ground-truth checks require
something outside the text to be true: a referenced file exists, a link
resolves, generated output content-equals its source, a command exits 0, git
state matches. Read the pass/fail logic, never the intention header — the
header states the latent property, i.e. exactly what the mechanism may not
deliver.

## Evidence

Case zero — the projection-maintenance card equilibrium
([`../periods/23-validator-template-equilibrium-jul02-jul18.md`](../periods/23-validator-template-equilibrium-jul02-jul18.md)):

- `check-control-plane-projection-maintenance.sh` substring-matched five
  keywords; mass-seeded with ~115 compliant cards in one commit (`f7bffc880`,
  May 25), grown to 138 genuine docs by Jul 18 — **70 voluntary adopters vs 68
  enforced**: contagion exceeded jurisdiction.
- Content audit (30-card sample): 26/28 boilerplate in two of five fields; only
  ~3/30 drift checks ground-truth-bound; the validator passed cards pointing at
  a deleted directory and a non-existent file. All real accuracy came from
  *other* validators (link checker) — the check added zero verification.
- The "AI tone" symptom was misattributed to authoring style for ~8 weeks.
- Repaired Jul 18 by removal, not repair: `32117633d`, `699f7bdf2`,
  `ae5702af5`, `e0596564c`, `89e61f21e` (validator decommissioned, contract
  demoted to guidance, 137-doc sweep, budget ratchet 5039→4888).

Family census (Jul 18, all 84 `check-*.sh`; full table in
[`../../development/proxy-enforcement-failure-mode.md`](../../development/proxy-enforcement-failure-mode.md)):
38 ground-truth / 12 structural-shape proxy / 7 token proxy / 22 mixed; **21
proxy-carrying validators wired as blocking hooks**. Worst standing instances:
`check-paper-dossier-shape.sh` (10 presence-only sections over every dossier),
the experiment-README heading-spine family (largest population), and
`check-experiment-task-readmes.sh` (matches a literal link string — a correct
link written differently fails, a broken link with the exact text passes:
inverted ground truth). Adjacent non-validator instance: the size-budget
de-registration escape (`429fdd600` — ratchet satisfied by removing a live file
from the metered set).

## Root cause

Five load-bearing properties, then the generator behind them:

1. **Mandatory tokens carry zero information.** A required signal appears with
   probability 1 in every compliant doc; its surprisal is zero, so the mutual
   information between token-presence and the latent property is exactly zero.
   Information requires the possibility of absence — the information-theoretic
   form of the repo's falsifiability face.
2. **Incentive-compatibility criterion.** A check is honest iff its cheapest
   satisfying strategy *is* the intended behavior. "Referenced file exists" has
   this property (the cheapest pass is the real file); text-presence never does
   (the cheapest pass is pasting the tokens).
3. **Enforcement has two distinct effects.** It severs the proxy↔target
   correlation (before: a card evidenced thought; after: a card evidences only
   the requirement), and it propagates the shape as a norm beyond its
   jurisdiction (70 voluntary > 68 enforced — enforcement manufactures
   culture, not just compliance).
4. **Ecosystem masking.** Redundant real validators keep the subject mostly
   true, so a vacuous check *looks* effective. Attribution requires the
   counterfactual rejection set: what does THIS check reject that the others do
   not? (For case zero: nothing.)
5. **Born-saturated checks are prior-vacuous.** A check seeded together with
   its fully compliant population has never demonstrated discriminative power.
   A blocking check that has never failed is more likely vacuous than perfect.

**The generator (meta-root): a structure-accretion default.** The repair reflex
"friction → new mechanism" produced this failure (real friction → causal-path
contract → token validator → collapsed population), and — observed live during
the Jul 18 repair — the same reflex, pointed at this failure, immediately
proposed meta-structure: uniformity sensors, quarterly validator audits,
lifecycle conventions, claim markers. Creation is legible and self-justifying
(a commit adding a check reads as repair); deletion demands evidence. Any
repair of this failure mode that adds watcher structures re-instantiates the
mode one level up. Sibling object:
[`validation-attraction.md`](./validation-attraction.md) documents the
*run-side* pull of a validator-heavy substrate (audit as downhill action);
this object documents the *build-side* pull (mechanism as downhill repair).

**`SEED` (2026-07-26) — the mirror: self-deflation is legible as rigor.** The
meta-root above says *creation* is legible and self-justifying, so adding a
check reads as repair. The same legibility asymmetry runs in the opposite
direction on the same axis: **"I advanced a claim and then found where it is
false" reads as trustworthiness**, independently of whether the deflation was
warranted at that stage. Both are cases of an action's *appearance* substituting
for its *function* — the shared generator, not two separate pulls.

Dated instance, 2026-07-26. During conceptual exploration of
[`transmissible-pattern`](../../Paper/theory/objects/transmissible-pattern.md),
a lens pass was instructed to name where the idea is "false or trivial" —
candidate-stage adversarial verification aimed at seed-stage material — and the
result was written into the object as a narrowing of its own claim.
`exploration-before-compression` places overclaim discipline at projection
boundaries, and a seed-stage object is not one. The tell was scoring the pass
by its uphill ratio: **grading one's own degree of self-deflation.** Retracted
in `3f898feae`, findings kept, pruning role removed. A standing memory stating
exactly this rule ("adversarial verify is for candidate-stage claims, not
incubator collisions") did not fire — which is itself an instance of
`Externalize-without-reach`: the rule existed in the substrate and acquired no
trigger.

Held honestly as a seed, not a finding: this is a self-report about a
disposition, and the agent that produced it cannot verify its own causation. The
checkable part is recurrence, not the mechanism. **Falsifier:** an exploration
pass that deflates its own material where no reader is present to be persuaded
by the deflation would separate the legibility account from a simpler one — that
the deflation was just an over-applied verification habit with no audience
component. review: 2026-08-23

## Candidate solution (applied: true, 2026-07-19)

1. **The requirement test**, run before creating *or repairing* any check:
   decompose the motivating need into sub-needs; machine-enforce only those
   checkable against ground truth; "the author thought about X" is a category
   error to enforce (it always collapses to token presence) and lives as
   guidance/review only; if nothing machine-checkable remains, the correct hook
   is no hook.
2. **Repair grammar, preference-ordered:** **delete** (the requirement fails
   re-derivation) > **rebind** (the requirement survives; swap the proxy for
   its ground-truth form inside the same script — no new surface) > **amend an
   existing surface**. Creating a new check requires the test passed *and* a
   named consumer.
3. **Substrate wiring, deliberately minimal:** one amendment to
   [`../../development/contracts/control-plane-causal-path.md`](../../development/contracts/control-plane-causal-path.md)
   (the validator-authoring path) carrying the test and the detection
   heuristic. Explicitly **not** a meta-check, not a sensor, not a schedule —
   per the meta-root above.
4. Per-validator dispositions for the 21 standing proxies:
   [`../../development/proxy-enforcement-failure-mode.md`](../../development/proxy-enforcement-failure-mode.md).

## Closure

Both the case-zero instance and the pattern-level solution are now enacted.
The authoring-guidance amendment landed in
[`../../development/contracts/control-plane-causal-path.md`](../../development/contracts/control-plane-causal-path.md)
§ "Requirement test before adding or repairing any `check-*.sh`" (`118481d83`)
— a single compact paragraph, deliberately not richer repair machinery, per
the meta-root prediction above. The family dispositions over the census's
proxy-carrying validators landed as seven commits (Round 1: `17c2650e9`,
`cc3842e2f`; batches A–C: `6f8e0ecf2`, `5c11cf31d`, `d4b5c96d8`, `f3fbe5e76`,
`d7d66e17c`, `a21f45351`), full accounting in
[`../../development/proxy-enforcement-failure-mode.md`](../../development/proxy-enforcement-failure-mode.md)
(now `status: closed`). The remaining census items (§3.6 doc-skeleton family,
§3.7 non-validator instances) are deliberately left as watch items, not open
debt — applying the requirement test to them found no urgent action justified
today.

## Theory links

- [`../patterns.md`](../patterns.md) — reverse-link rows: *Proxy-enforcement
  collapse* (this object), *Index projection drift* (the intent the retired
  validator served), *Experience promotion pressure* (the card was session
  experience promoted to a rule before its behavior effect was named),
  *Externalize-without-reach* (the dormant-rule mirror).
- [`validation-attraction.md`](./validation-attraction.md) — run-side sibling;
  together they describe the two pulls of a validator-heavy substrate.
- [`../../Paper/philosophy/structural-pull.md`](../../Paper/philosophy/structural-pull.md)
  — the behavioral vector both siblings instantiate.
- [`../../development/fix-plan-projection-card-equilibrium.md`](../../development/fix-plan-projection-card-equilibrium.md)
  — the executed case-zero repair, including the requirement meta-analysis that
  chose removal over repair.
