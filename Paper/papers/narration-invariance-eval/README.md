---
paper_id: narration-invariance-eval
claim_bundle: narration-invariance-eval
venue: aaai-27
format-envelope: 7page-archival-two-column
track: theory
cluster: operator-projection-eval
status: planning
paused: true
state_mode: planning-skeleton
role: cross-environment-invariance-eval-primitive
family: eval-primitives
publication_class: theory-candidate
maturity_target: boundary
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: README.md
view-affinity: [reflexive]
theory-objects: [narrative-claim]
gaps: [function-vs-narrative]
---

# Cross-Environment Invariance: A Causal Test for Load-Bearing vs Narrated Competence

## Status: PAUSED (inactive, 2026-06-23); venue binding lapsed (2026-07-26)

**Venue update 2026-07-26.** The AAAI-27 cycle was missed
([`venues/objects/aaai-27.md`](../../venues/objects/aaai-27.md) § Cycle
outcome) and this line is withdrawn from it. This changes nothing about the
line's substance: it was deliberately paused on `2026-06-23` behind a restart
gate it had not met, and a paused line missing a deadline is the pause working,
not a failure. The `venue: aaai-27` frontmatter is kept as the historical
target — the AAAI-importance framing in the restart gate below is what that
binding was for. **Do not re-bind a venue as a way of restarting this line;**
the gate below is the restart condition, and the venue question comes after it.

**This line is paused — not abandoned.** Importance cannot be judged without
the struggle-zone experiment (does the phenomenon generalize to frontier models
at the edge of capability?), and that experiment is not worth the investment
until the existential confound is resolved. The line is parked pending two
make-or-break results; do NOT invest further experiment effort until the restart
gate is met.

**Restart gate (both required before reactivation):**
1. **D-control (existential):** a real-word-nonsense dressing arm (e.g.
   "glomerates" — real English tokens, no real-world intuition) — separates
   narration-channel dependence (breaks on D) from nonce-OOD fragility (breaks
   only on nonce C). If D resolves AGAINST, the narration story collapses and
   the line retires.
2. **Struggle-zone scaling (importance):** scale math-graph complexity until
   frontier models (claude/opus/deepseek, all 0% in the pilot) degrade, then
   test whether they fall back on narration. Converts "niche methodology, 1/4
   models" → "major finding" — the prerequisite for an AAAI-importance case.

**Why paused:** Gemini novelty/importance review converged with the internal
read — the contribution (methodological discrimination assay) is real and
defensible, but (a) the nonce-OOD confound must be cleared before the
phenomenological claim is submittable, and (b) the 0%-in-frontier-models result
makes the importance case unjudgeable without the struggle-zone experiment. See
[`references/gemini-novelty-importance-2026-06-23.md`](references/gemini-novelty-importance-2026-06-23.md).
No more experiment effort on this line until the gate is met.

**What stands regardless of pause (do not discard):** the operator + 3-arm
assay design; the D1–D4 novelty source-check; the 4-model pilot (a real,
B-control-isolated, model-dependent result); the filed predictions P4/P5/P6;
the foreign-domain registry and Case 03 (procedure-reflection) are independent
of this line's activity and stay.

## Current Question

When an evaluation signal (an interpretability claim, a citation, a spec-match)
can be satisfied by **description** rather than by un-narratable state-change,
how do you causally separate competence that is load-bearing from competence
that is merely narrated?

## Thesis

An evaluation signal that is **narration-satisfiable** will see measured
competence inflate against a **cross-environment invariance** metric, and the
inflation concentrates where the artifact looks cleanest. The missing causal
operator is a controlled **environment shift**: perturb the readable channel the
signal depends on, leave the functional mechanism intact — a load-bearing
feature is invariant across the shift; a narrated feature drifts with the
perturbed channel. Cross-environment invariance (after Arjovsky et al.'s IRM) is
the discriminator the function-vs-narration gap named as missing.

**Sharper framing (Legibility Penalty).** The demand for describability does
not merely *observe* decoupling — it *manufactures* it: any mechanism optimized
for human legibility (a clean circuit story, a supportive citation, a coherent
reasoning trace) is under pressure to satisfy the description constraint by
decoupling from the load-bearing computation. Legibility is not a window into
the model; it is a spurious attractor. This turns worst-where-cleanest from a
correlation into the signature of an attractor. (Empirical prediction, not a
proven law.)

## The Operator (source)

Projected from the [`narrative-claim`](../../../theory/objects/narrative-claim.md)
object and the [`function-vs-narrative`](../../../theory/gaps/function-vs-narrative.md)
gap, stripped of its home substrate:

- state is cheap to describe but expensive to make/carry; an account of a
  state-change can be produced at near-zero cost whether or not the change
  happened;
- the discrimination operator is cross-environment invariance — a feature
  causally linked to the outcome is stable across environments; a spurious
  (narrated) feature breaks;
- the attractor prediction: any selection/eval signal that *can* be satisfied by
  narration rather than by un-narratable state-change will see the population
  collapse onto convincing-but-counterfactual narration.

This dossier holds the **operator-level claim**. The manuscript picks one lead
projection and uses the others as transport corroboration.

## Four Projections (the operator transports; working lead = Arm A)

| Foreign domain | The hole the operator names | Lead falsifiable prediction |
|---|---|---|
| **Chain-of-Thought / System-2 reasoning faithfulness** (o1-style; **WORKING LEAD**) | Are a model's "reasoning traces" load-bearing (causing the answer) or narrated (post-hoc rationalization satisfying an RLHF reward for "looking like thinking")? The unfaithful-CoT literature (Turpin, Lanham) shows unfaithfulness exists; it lacks the causal-invariance operator that separates load-bearing from narrated reasoning holding formal topology constant. | Semantic-isomorphism shift on MATH/GSM8K: translate word problems into meaningless tokens ("blargs/glops") preserving the exact math graph (functional mechanism held constant) while destroying narrative resonance. Load-bearing reasoners invariant; narrated reasoners collapse. Worst-where-cleanest: the most coherent-looking traces are least load-bearing. |
| **Mechanistic interpretability** (original candidate; **now fallback arm**) | A feature/circuit explanation ("induction head," a monosemantic SAE feature) is narration-satisfiable. | Token-identity permutation; **monosemantic-by-narration** subset. Sharpest testbed: **Othello-GPT** (linear probe = narration vs nonlinear board state = mechanism). Fallback if Arm A's experiment proves unworkable. |
| **RAG attribution faithfulness** | "Groundedness" scores *consistency* — narration-satisfiable. | Counterfactual passage swap (ALCE / FaithDial). §4 transport corroboration. |
| **Program synthesis (spec vs code)** | pass@k + spec-shaped surface form is narration-satisfiable. | Multiple honest operationalizations of one NL spec as environment shifts. §4 transport corroboration. |

The same structural prediction recurs across four unrelated fields — each
field's eval signal is narration-satisfiable, so each sits in the attractor the
operator names.

**Working lead domain: Arm A (CoT / System-2 reasoning faithfulness).** The
novelty boundary currently looks promising in the checked-in source-check table
([`plans/dr-lead-domain-report.md`](plans/dr-lead-domain-report.md), interaction
`v1_ChdVYzg2YXB1…`). Mapped against four distinguishing criteria — **D1** shift
targets the narration/legibility channel specifically; **D2** functional topology
held mathematically constant; **D3** cross-environment invariance measured as
causal attribution; **D4** the Legibility-Penalty attractor claim — the current
table does not show a listed prior occupying all four, with primary-source
confirmation still pending:
- *Illusion of Thinking* / *GSM-Symbolic* (Apple): scale complexity / swap input
  features → measure **capability collapse** (not invariance); D1, D3, and D4
  absent; D2 only partial for number swaps and broken by distractor-clause arms.
- *Turpin* (unfaithful CoT): biases the **prompt**, measures answer-change, not
  invariance; D1 absent.
- *Lanham* (Anthropic, closest behavioral kin): targets the trace (**D1 ✓**) and
  tests answer-invariance (**D3 partial**) — **but injects mistakes / truncates,
  breaking the math graph (violates D2)**: they test "does breaking logic break
  the answer," not "does changing legibility holding logic reveal causal
  attribution."
- *Triangulation* framework (late-2025, closest **principle** kin): already
  establishes cross-environment invariance as the gold standard for causal
  attribution (**D3-as-principle is NOT novel**) — but applies it to *internal
  circuitry* across language/script shifts, not external traces, and changes
  topology via translation. **D3's novelty is application-only (external trace +
  topology-preserved shift), not the principle itself — carry this honestly.**
- *Obfuscated-CoT / LITL* (security): philosophically aligned with D4 but no
  formal operator.

**Reviewer preemptions to carry** (surfaced by the source-check): (a) the
*Triangulation* lineage already established invariance-as-causal-gold-standard —
the contribution is applying it to the external trace under topology
preservation, not inventing it; (b) **CausalDPO (2026) uses the identical phrase
"cross-environment invariance"** for OOD recommendation generalization —
distinguish explicitly (preference optimization vs trace-faithfulness
attribution); (c) *Lanham* is the citation a reviewer will reach for first — the
D2 distinction (they break logic; we hold it) is the one-paragraph defense.

## In Scope

- the cross-environment invariance operator as a causal test for load-bearing vs
  narrated competence
- one lead domain (CoT / System-2 reasoning faithfulness) developed for the main
  experiment, with RAG and code-synthesis as transport-corroboration sections
- the narration-satisfiability attractor prediction (inflation vs invariance;
  worst-where-cleanest) as the falsifiable claim

## Out Of Scope

- re-deriving the operator from the Quine substrate (the operator is imported;
  the home evidence is not cited in the manuscript)
- claiming the operator explains *all* eval pathology (it names one structural
  failure class)
- a full causal-theory-of-evaluation treatise; the IRM grounding is borrowed, not
  extended

## Current State

- manuscript workspace: `./WORKSPACE.md` (stub)
- current manuscript plan: `./manuscript-plan.md` (lead-domain skeleton)
- track judgment: seed-dossier; the operator-radius verification (see provenance)
established the projections run entirely on foreign-domain artifacts with zero
home-evidence reference. The manuscript is not yet drafted.

## Writing Posture

- **AAAI audience lens:** `ai-ml-core` primary, `ai-safety` secondary
  (eval/interpretability/reasoning faithfulness). De-ALife entirely — the
  operator's home substrate does not appear in the manuscript; it is cited once
  as the source of the operator if at all.
- **Working lead domain: Arm A (CoT / System-2 reasoning faithfulness)** after
  the deep-research novelty-boundary source-check
  ([`plans/dr-lead-domain-report.md`](plans/dr-lead-domain-report.md)). The
  D1–D4 intersection appears unoccupied in the current source-check. Arm B
  (mech-interp) is the fallback if the CoT experiment proves unworkable.
- **Reviewer-risk 1 — "just shortcut learning / Goodhart rebranded."** Answer:
  shortcut learning is about spurious features in *inputs*; this is about the
  system's *generated explanations* of its own state being spurious, AND the
  claim that the demand for legibility *manufactures* the decoupling. Front-line
  this distinction.
- **Reviewer-risk 2 — "the environment shift is ill-defined."** Answer: use
  formal testbeds (MATH/GSM8K semantic-isomorphism) where functional-topology
  preservation is mathematically guaranteeable.
- **Reviewer-risk 3 (primary, carry the D1–D4 defense from Four Projections).**
  Closest kin: *Lanham* (targets trace but breaks logic → D2 violation);
  *Triangulation* (already established invariance-as-causal-gold-standard →
  D3's novelty is application-only, not the principle — state this);
  *CausalDPO* (identical "cross-environment invariance" phrase, different task —
  distinguish). *Illusion of Thinking / GSM-Symbolic* perturb inputs and test
  capability, not trace-attribution.
- **Title/tone:** operational and falsifiable, not philosophical.

## Experiment Feeds

- `experiments/active/cot-faithfulness-invariance/cfi-01-semantic-isomorphism/README.md` - primary - main-text-candidate - the 3-arm semantic-isomorphism pilot (run; see note on pause)

## Current Evidence Backbone

- `experiments/active/cot-faithfulness-invariance/cfi-01-semantic-isomorphism/README.md` - `pilot-4-model` - supporting-only - 4-model pilot: B-control-isolated narrated-trace subset in 1/4 models (gpt-5.4-mini 17.4%); not decisive (model-dependent, nonce-OOD confound open)

## Transport Validation

| Claim | Predictions | Views tested | Anomalies | Transport-valid? |
|-------|-------------|--------------|-----------|------------------|
| narrated-trace subset exists in CoT (operator non-vacuous) | narrative-claim-P5 (confirmed) | reflexive (operator source) | none | YES |
| worst-where-cleanest: legibility ↔ load-bearing negative correlation | narrative-claim-P4 (open) | reflexive | none | NO |
| Legibility Penalty is causal (manufactures decoupling) | narrative-claim-P6 (open) | reflexive | none | NO |

### Paper readiness

- transport-valid claims: none (P5 partial pilot only — narrated-trace subset
  in 1/4 models: 17.4% in gpt-5.4-mini with B-control isolation, 0% in
  deepseek-v3.2 / claude-sonnet-4.6 / claude-opus-4.8; model-dependent, N=24 —
  not decisive; see
  [`experiments/active/cot-faithfulness-invariance/cfi-01-semantic-isomorphism/`](../../../experiments/active/cot-faithfulness-invariance/cfi-01-semantic-isomorphism/README.md))
- unresolved anomalies: none
- novelty: promising (D1–D4 intersection appears unoccupied in the current source-check; worst-where-cleanest + causal Legibility Penalty need checked-in citation support before manuscript use)
- quality: pilot run complete (2 models, 3-arm design); P4/P6 not yet tested
- status: seed-dossier with working lead domain + filed falsifiable spine + first pilot evidence

## Strongest Evidence

No manuscript evidence yet. The strongest support is **operator-radius
verification** (2026-06-23): four parallel operator-projection subagents under
the anti-domestication contract produced this collision set with zero home-field
evidence reference and falsifiable predictions running entirely on foreign-domain
artifacts (SAEs/induction-head, ALCE/RAGAS, HumanEval/MBPP). That supports the
claim that the operator is testable in-domain; it does not establish the
predictions hold.
Provenance: [`../../../theory/views/reflexive/deep/cases/03-the-evidence-anchor-attractor.md`](../../../theory/views/reflexive/deep/cases/03-the-evidence-anchor-attractor.md).

## Blocking Gaps

- D-control existential confound unresolved: a real-word-nonsense dressing arm
  must separate narration-channel dependence from nonce-OOD fragility before the
  line can reactivate.
- struggle-zone scaling unresolved: frontier models need a capability-edge regime
  before importance can be judged.
- checked-in citation extraction still needed before any manuscript uses the
  foreign-domain source-check beyond planning posture.

## Next Useful Work

1. **working lead domain — Arm A (CoT).** Next is freezing the main experiment:
   MATH/GSM8K semantic-isomorphism shift (translate word problems into
   meaningless tokens preserving the math graph) scored by trace-invariance
   under the shift.
2. keep the filed prediction objects current (`narrative-claim-P4/P5/P6`) as the
   D-control and struggle-zone work advances
3. primary-source-verify the deep-research report's anchors before citation
   (Illusion of Thinking, GSM-Symbolic, Turpin 2023, Lanham 2023 are
   well-established; the 2025-2026 items — Triangulation, CMI, LITL, CausalDPO —
   need primary-source confirmation of dates/claims before they enter the
   related-work paragraph)
4. ~~draft the one-paragraph D2-vs-Lanham novelty defense~~ — done:
   [`manuscript-plan.md`](manuscript-plan.md) § Reviewer Defense — D2 vs Lanham
   (still needs the `[cite: 13, 14]` markers resolved into a real bibliography
   entry per item 3 before it is submission-ready)

## Unsafe Overclaims

- do not claim the operator is validated by the home substrate's experiments
  (those are not cited; the manuscript stands on its own foreign-domain runs)
- do not claim all interpretability work is "narrated" — the claim is a
  measurable *subset* and a literature-level bias, not a universal indictment
- do not claim IRM is novel here; the mechanism is borrowed, the contribution is
  the operator + transport + worst-where-cleanest prediction
- do not let the three-domain transport read as three independent results before
  any domain has a run

## Why This Paper Exists / Family Note

This is the first dossier in a new `eval-primitives` family: operator-projection
papers that take a Quine theory operator (here `narrative-claim` /
function-vs-narrative) and project it into a mainstream AI/ML evaluation
conversation the repo's experiments never touch. It is distinct from every
existing family (not ALife, not runtime-substrate, not repo-self-model) and is
the direct output of the generative operator-projection method verified in
Case 03. If the lane matures, it absorbs sibling operator-projection dossiers
(carrier × distributed-systems, structural-pull × contamination, recombination ×
search — see [`../../../theory/lenses/foreign-domains.md`](../../../theory/lenses/foreign-domains.md));
if it does not, this single seed is abandoned cleanly with no orphaned theory.

## Backlog Maturation Pass (2026-07-02)

Appended note only — not a rewrite. A 104-seed backlog-maturation pass (see
`Paper/paradigms/backlog-maturation-2026-07.md`) filed the following seed as
an extension of this dossier's territory, the sole current member of the
nascent `eval-primitives` family:
- `context-decorrelation-vs-model-identity`: the raw claim (cross-context
  review beats same-session self-review) is already published and closed
  (Song, arXiv:2603.12123) — not novel for Quine to rediscover. What
  genuinely remains open, flagged as a limitation in that paper's own
  follow-up and never isolated externally, is whether the gain is fully
  explained by session/context decorrelation or whether swapping the
  reviewer's base model adds independent signal beyond it — a 2x2
  {same/cross-session} x {same/cross-model} design, currently unrun anywhere
  and cheaper for a multi-model-infrastructure repo like Quine to run than
  for most labs.
