# Gemini Spark — narration-invariance-eval novelty + importance pressure-test

> Recruited via `agy` (Gemini 3.1 Pro High), 2026-06-23, per the
> [Theory-Spark Query Contract](../../../../development/contracts/theory-spark-query.md).
> Spark source, not executor; this is the raw spark + owner integration. Do not
> cite Gemini directly — cite the integrated deltas (abstract v2, roadmap).

## Prompt (compressed)

Fed the core hypothesis, abstract v1, the honest 4-model pilot (17.4% in
gpt-5.4-mini, 0% in deepseek/sonnet/opus; B-control 0%; narration-vs-nonce-OOD
confound open), and the source-checked prior-art landscape. Asked for the
sharpest read on (1) novelty and (2) importance; biggest threat to each.

## Raw sparks (Gemini, unshaped)

### Novelty
- novelty is **methodological, not conceptual**: shortcut learning + unfaithful
  CoT are known; the contribution is the topology-preserving invariance
  *assay* on external traces. Defensible if you sell the scalpel, not the
  pathology.
- **biggest novelty threat: GSM-Symbolic** — reviewers reflexively "this is just
  GSM-Symbolic." Defense = weaponize the B-vs-C delta (B 0%, C 17.4%) in
  abstract sentence 2; GSM-Symbolic measures generic variance, this isolates
  load-bearing semantic reliance. Lanham defense (topology preservation) is
  airtight, dispatch in related work.
- "narration-satisfiability" is a **useful high-resolution rebrand** of shortcut
  learning, not a discovery — name + diagnose a suspected-but-unisolated
  mechanism; don't claim you discovered gravity.
- confound undercuts **phenomenological** novelty ("we found narrated traces")
  but NOT **methodological** novelty (the 3-arm design surfaced the confound — a
  feature).

### Importance
- **not AAAI-worthy as framed with D unresolved** — reviewers spot the
  nonce-OOD confound and reject "come back when finished." A methodological
  paper must isolate the variable it claims. Clear D first.
- AAAI does love rigorous causal eval primitives; "honest boundary" framing is a
  strength (maturity), not a weakness.
- **biggest importance threat: 0% in 3 frontier models** → "you built an assay
  to prove a pathology that only exists in a cheap mini model. Who cares?"
  Fix: **find the struggle zone** — scale math-graph complexity until frontier
  models degrade; show they fall back on narration at the edge of capability.
  Niche-methodology → major finding.
- if D resolves AGAINST: paper collapses to "CoT brittle to OOD lexical
  substitution" — weaker, less novel, probably sub-AAAI.

### Strategy (Gemini)
1. run Arm D (existential). 2. find the frontier struggle zone. 3. rewrite
abstract to distance from GSM-Symbolic via B-vs-C delta in sentence 2.

## Integration (owner)

**Adopted:**
- methodological framing (sell the scalpel); "narration-satisfiability" =
  rebrand not discovery → abstract v2 rewritten accordingly.
- B-vs-C vs GSM-Symbolic in abstract sentence 2 → abstract v2.
- D-control existential (matches my own read).
- struggle-zone scaling as the importance strategy — operationalizes my
  "regime-gated" reading into a concrete plan (scale graph complexity until
  frontier models enter the struggle zone, then test narration fallback).
- honest-boundary-as-strength kept.

**Corrected:**
- "do not submit with D unresolved" — correct for AAAI main + phenomenological
  claim; the assay-with-open-confound is a workshop/short-paper fallback, not a
  total dead end (B-vs-C delta is a residual diagnostic even if OOD).
- if D against: narrative collapses (Gemini right) but B-vs-C delta still a
  weaker OOD diagnostic — honest carry, not full collapse.

**Net:** Gemini converges with and sharpens my own read. Two next experiments
are make-or-break, not optional: **D-control** (existential — does the phenomenon
survive as narration-channel?) and **struggle-zone scaling** (importance — does
it generalize to frontier models at the edge of capability?). Abstract updated
to v2 (methodological framing, B-vs-C vs GSM-Symbolic foregrounded).
