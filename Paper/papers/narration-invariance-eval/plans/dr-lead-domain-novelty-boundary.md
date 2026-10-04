# Deep Research Brief — narration-invariance-eval lead-domain novelty boundary

## Context

We are deciding the **lead domain** for an AAAI-27 paper (`narration-invariance-eval`).
The paper projects an operator — **cross-environment invariance as a causal test
for load-bearing vs narrated competence** — into a mainstream AI/ML evaluation
conversation. Two candidate lead domains remain:

- **Arm A — Chain-of-Thought / System-2 reasoning faithfulness** (o1-style
  reasoning traces). Hotter, and its experiment is cleaner (functional mechanism
  formally verifiable in math/logic). BUT it is adjacent to a body of prior art
  whose novelty boundary must be fought explicitly.
- **Arm B — Mechanistic interpretability** (circuit/feature explanations).

The decision turns on ONE question: **can our operator be distinguished from the
CoT-perturbation prior art in a single clean paragraph a reviewer will accept?**
If yes → Arm A wins (hotter, cleaner). If the boundary is too thin or already
occupied → Arm B wins (novelty-denser). This research resolves that.

## Core Research Question

For the adjacent prior art on perturbation-based tests of LLM reasoning /
chain-of-thought, determine exactly **what each perturbs, what each holds
constant, what each measures, and what each claims** — so we can locate whether
our specific operator is already occupied or genuinely novel.

## Our Operator (the thing being compared against prior art)

A controlled **environment shift** that perturbs the *narration / legibility
channel* an evaluation signal depends on, while holding the *functional
mechanism / formal topology* intact. A **load-bearing** feature/trace is
**invariant** across the shift; a **narrated** one **drifts**. Two accompanying
claims: (1) **worst-where-cleanest** — the cleanest-looking explanations/traces
are the least load-bearing (negative correlation between legibility and causal
robustness); (2) **the Legibility Penalty** — the demand for describability
*manufactures* the decoupling (legibility is a spurious attractor, not a window).

The crucial distinguishing features to check against prior art:
- (D1) the shift targets the **narration/legibility channel specifically**, not
  the input distribution;
- (D2) **functional topology is held mathematically constant** (not just
  "similar");
- (D3) the measured signal is **cross-environment invariance as causal
  attribution** (not accuracy/capability under perturbation);
- (D4) the **worst-where-cleanest** correlation and the **legibility-manufactures-decoupling**
  attractor claim.

## Detailed Questions (per prior-art item)

For EACH item below, extract: (a) exact perturbation performed; (b) what is held
constant; (c) what is measured (metric); (d) the central claim; (e) does it
occupy any of D1–D4?

1. **Apple — "The Illusion of Thinking: ... in Reasoning Models" (2025).** What
   do they perturb (puzzle variants? formatting? tokens?); is functional
   topology preserved or changed; do they measure capability-collapse or
   causal-attribution; do they make any worst-where-cleanest / legibility claim?
2. **Mirzadeh et al. (Apple) — GSM-Symbolic (2024).** Same questions. What
   exactly is the symbolic substitution; is the math graph preserved; what does
   the performance drop claim about reasoning vs pattern-matching?
3. **Turpin et al. — "Language Models Are Active Explainers" / the unfaithful-CoT
   line (2023+).** What intervention tests CoT faithfulness (biasing prompts,
   paraphrasing, truncation); do they ever hold the final-answer-computation
   constant while perturbing only the trace; do they measure invariance?
4. **Lanham et al. (Anthropic) — CoT faithfulness / early-exit /
   thought-injection tests (2023).** Same questions; does perturbing/removing
   CoT tokens and measuring answer-change constitute a causal-attribution test?
   How does it differ from D3?
5. **Any prior work doing semantic-isomorphism / formal-equivalence
   perturbations on math benchmarks** (preserving the solution graph while
   changing surface dressing) — e.g. counterfactual task variants, template
   perturbations, "variance under reformulation" studies.
6. **Any causal-attribution / mediation work on CoT tokens specifically** — path
   patching, activation patching, ablation that asks "is the CoT doing the work"
   (rather than "is the model right").
7. **Any work explicitly claiming a negative correlation between
   explanation-quality / legibility and faithfulness / load-bearingness** (i.e.,
   someone who already occupies D4).

## Likely anchors / search directions

- Apple ML research page + arXiv for Illusion of Thinking / GSM-Symbolic /
  GSM8K-platinum / Math-Lite
- Anthropic alignment forum + arXiv for Lanham CoT faithfulness, Turpin
  unfaithful explanations, Anthropic 2024-2025 CoT-verbosity / CoT-need work
- "chain of thought faithfulness", "unfaithful chain of thought", "does CoT
  cause the answer", "CoT ablation", "thought editing", "reasoning trace
  causal", "semantic perturbation math benchmark", "counterfactual invariance
  reasoning"
- IRM / invariant-predictor lineage applied to LLM reasoning (if anyone has)

## Relevance filters

- **In scope:** perturbation / intervention studies on LLM reasoning traces or
  math benchmarks that test faithfulness, causal role, or robustness of the
  *reasoning* (not just final-answer accuracy).
- **Primary:** peer-reviewed or strong-lab arXiv (Apple, Anthropic, DeepMind,
  top venues) 2023–2026.
- **Out of scope:** general benchmark-construction without a
  faithfulness/causal claim; pure capability scaling papers; mech-interp
  circuit work (that is Arm B's separate source-check, not this one).

## Deliverables Requested

1. **A comparison table**: rows = the prior-art items (1–7), columns =
   {perturbation, held-constant, metric, central claim, occupies D1?, D2?, D3?,
   D4?}.
2. **A one-paragraph novelty-boundary statement** that EITHER cleanly
   distinguishes our operator from all located prior art (naming the precise
   axis of difference) OR identifies the prior work that already occupies the
   operator and names the residual novelty (if any).
3. **A lead-domain recommendation** grounded in the table: Arm A (CoT) viable iff
   the novelty-boundary paragraph is clean; else Arm B (mech-interp).
4. **Any additional adjacent work** we missed that a reviewer would cite against
   Arm A.

## What to extract

- the exact operational definition of each prior perturbation (so we can state
  "they change X; we hold X constant and change Y instead")
- whether anyone measures *invariance* (not *change*) as the diagnostic
- whether anyone makes the legibility-vs-load-bearing negative-correlation claim

## Manuscript payoff

This resolves the dossier's first blocking decision (lead domain) and feeds the
§6 related-work paragraph + the §1 novelty statement of the manuscript plan. The
output table becomes the related-work comparison artifact.

## Red line / non-claim

Do NOT report "this is novel" without the table — the whole point is to falsify
our novelty assumption. If prior art occupies D1–D4, say so plainly; the paper
is better served by Arm B than by a fightable-but-losing novelty boundary at
Arm A.

## Stop conditions

- the comparison table + novelty-boundary paragraph stop changing the
  lead-domain recommendation, AND
- at least the four named anchors (Illusion of Thinking, GSM-Symbolic, Turpin,
  Lanham) are located with their exact perturbation defined, AND
- the D1–D4 occupancy is determined for each.
