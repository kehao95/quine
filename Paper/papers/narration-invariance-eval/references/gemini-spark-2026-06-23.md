# Gemini Spark — narration-invariance-eval lead-domain pressure-test

> Recruited external perspective via `agy` (Gemini 3.1 Pro High), 2026-06-23, per
> the [Theory-Spark Query Contract](../../../../development/contracts/theory-spark-query.md).
> Gemini is a spark source, not an executor; this log is the raw spark, retained
> durably so the recruited source is recoverable. Integration (what is valuable
> vs overclaim) lives in the dossier README and at the bottom of this file — do
> not cite Gemini's claims directly; cite the integrated deltas.

## The prompt (compressed)

Fed the operator essence (function-vs-narrative → cross-environment invariance;
IRM-inspired), the three projections (mech-interp lead, RAG, code-synthesis),
the worst-where-cleanest prediction. Asked for: sharpest framing, foreign
collisions missed, lead-domain argument, strongest reviewer kill-shot, named
testbeds. Generative genre; no audit.

## Raw sparks (Gemini's output, unshaped)

### Sharper framing — "The Legibility Penalty"
- legibility is not a window into the model; it is a **spurious attractor**;
  demanding a state-change be cheap to describe makes the system satisfy the
  description constraint by decoupling from the load-bearing computation;
- provocative form: "semantic coherence is inversely proportional to causal
  robustustness" — the features/citations/steps human judges rate most coherent
  are the first to break under an environment shift;
- surgical test: a transformation T that scrambles the semantic/narrative
  surface while holding the formal/topological structure constant; load-bearing
  → invariant under T; narrated → collapses.

### Foreign collisions missed (the big one)
- **Agentic reasoning traces / System-2 CoT (o1-style)** — the hottest 2025-2026
  zone. Are hidden reasoning tokens load-bearing (causing the answer) or
  post-hoc rationalization satisfying an RLHF reward for "looking like thinking"?
  Shift: swap a logical operator to its inverse (break load-bearing logic) but
  keep "Therefore"/"Thus"/confident tone (narration); or scramble narrative glue
  to gibberish while keeping the math operations intact. Prediction: model
  follows confident narration even when logic is broken.
- secondary: **Reward-Model collapse in RLHF** (RM rewards polite authoritative
  text not truth); **VLM "blindness"** (narrating what it "sees" from language
  priors).

### Lead-domain argument — pivot to CoT
- mech-interp is a niche subculture with infighting over what a "true" circuit /
  monosemantic feature is; 4 of 7 pages spent defining terms and defending
  baselines;
- CoT is hotter (every lab building o1-clones; AAAI desperate for rigorous
  faithfulness eval); experiment cleaner (functional mechanism formally
  verifiable in math/logic, so the environment shift is better-defined).

### Kill-shots
- **"Just shortcut learning / Goodhart rebranded."** Answer: shortcut learning
  is about spurious features in INPUTS; this is about the system's GENERATED
  EXPLANATIONS of its own state being spurious, and the demand for legibility
  MANUFACTURING the decoupling — not finding shortcuts in data but proving
  evaluations force the model to decouple explanation from mechanism.
- **"The shift is ill-defined."** Answer: use formal testbeds (code/math) where
  functional topology preservation is mathematically guaranteeable.

### Named testbeds
- CoT: MATH, GSM8K under programmatic semantic-isomorphism shift (translate word
  problems into meaningless tokens preserving the math graph); Bamboogle.
- RAG: ALCE, FaithDial, HaluEval.
- mech-interp: IOI (Indirect Object Identification); Othello-GPT (linear probe =
  narration vs nonlinear board state = mechanism; perturb the read-out without
  changing next-move prediction).

---

## Integration (owner's shaping — what is valuable vs overclaim)

**Adopted:**
1. **Lead domain is now a genuine open decision, not settled.** CoT/reasoning
   faithfulness is a real contender — hotter, cleaner experiment, shift
   better-defined, most-provocative worst-where-cleanest. This is the single
   most valuable spark. The dossier's "mech-interp recommended" is downgraded to
   "mech-interp vs CoT, real tradeoff."
2. **"Legibility Penalty" as the sharper framing** — legibility *manufactures*
   the decoupling (not just observes it). Adopted as thesis sharpening. Dropped:
   "thermodynamic law / mathematically guaranteed" (we have no proof; it is an
   empirical prediction).
3. **Kill-shot answers** (shortcut-learning distinction; formal-testbed answer)
   adopted into Writing Posture / Unsafe Overclaims.
4. **Testbeds** adopted: Othello-GPT (the sharpest mech-interp testbed — linear
   probe vs nonlinear board state is exactly narration-vs-mechanism); IOI; ALCE;
   FaithDial; MATH/GSM8K/Bamboogle for the CoT arm.
5. **New domains** for the foreign-domains registry: CoT/reasoning-faithfulness
   (primary), RLHF-RM-collapse and VLM-blindness (secondary).

**Rejected / corrected (Gemini overclaim):**
1. **CoT novelty overstated.** Gemini framed CoT as if uncontested ("no one has
   the test"). It is NOT — there is a substantial **unfaithful-CoT literature**
   (Turpin et al. 2023; Lanham et al. 2023) AND adjacent prior art that perturbs
   math problems and shows collapse: **Apple's "The Illusion of Thinking" (2025)
   and GSM-Symbolic (2024)**. This is the serious obstacle Gemini missed. The
   operator IS distinguishable (Illusion-of-Thinking tests capability under
   input distribution shift; this tests causal attribution of the *reasoning
   trace* holding functional topology constant), but a reviewer WILL conflate
   them, so the novelty boundary must be fought explicitly. Net: CoT is hotter
   but more crowded; mech-interp is colder but novelty-denser for this exact
   operator. Genuine tradeoff, not the cakewalk Gemini implied.
2. **"Mathematically guaranteed" inverse proportionality** — dropped (no proof).

**Net effect on the paper:** the lead-domain decision is now the key blocking
choice, and the manuscript-plan must carry BOTH arms with the tradeoff stated.
The sharper framing (Legibility Penalty) and the testbed inventory are
immediate upgrades. The CoT adjacent-prior-art (Illusion of Thinking,
unfaithful-CoT) is a mandatory source-check before any CoT lead is chosen.
