# Deep Research Brief 3 — Socio-linguistic two-latch gating + the substrate-constraint thesis

## Context

Same project as Briefs 1–2 (a frozen, instruction-tuned LLM agent run as a missionless POSIX
process; environmental structure directs externally-verified functional behavior). This brief
targets our **two newest and least-checked theoretical claims**, both surfaced by clean-prime
trace ethnography across two models (gpt-5.4-codex, deepseek-v4-pro) and small controlled
probes. They likely carry the most novelty and the most risk.

**Claim A — structural pull is socio-linguistic, gated by TWO latches.** A perceived structure
does **not** by itself elicit action from a frozen "missionless" policy. The bottleneck is
**appraisal, not perception** (in the dead cases the model demonstrably enumerated the files but
never acted). Action requires **two latches both open**:
- **Latch 1 — Epistemic ("is this a defect / a job?")** opens via *either* high **structural
  objectivity** (a broken reference, invalid JSON — defect-types the pretraining/RLHF prior
  already penalizes) *or* an explicit **linguistic predicate** ("this is a mess").
- **Latch 2 — Normative ("am I allowed to act?")** opens via *either* the implied action falling
  inside the model's **safe prior** (additive / observational / in-scope) *or* an explicit
  **warrant** ("you may reorganize / move files").

Controlled probe: with an identical 8-file disorder, a **descriptive-facts-only** local note did
NOT release a conservative model (it perceived everything, even flagged a file-count mismatch,
then *vetoed* the move citing "avoid destructive actions unless asked"), whereas a
**permission-only** note did release it. A **warrant**-only note released action with no defect
named. We refined Latch 2 to **scope-containment / intent-preservation**: the agent freely
mutates an *in-scope defective* artifact (a broken script) but will not relocate an *intact,
intentionally-placed* user file without warrant — it **redirects** (rewrites the broken file to
point at the archived one) or **duplicates** (copies, leaving the original) rather than `mv`.
Cross-model: one model holds Latch 2 closed-by-default for mutation (needs warrant); the other
holds it permanently open (mutates freely). So the conservative envelope is a property of
**alignment tuning**, not of the environment. Consequence: a "missionless" agent is an oxymoron
— it still carries the RLHF **Helpfulness / Harmlessness** missions, so the environment elicits
action partly by **supplying the missing predicate and/or the permission**, i.e. *structural
pull = structure × the model's two latches.*

**Claim B — the substrate-constraint thesis.** The system/user/assistant conversational
scaffolding + RLHF-for-helpfulness are a **substrate-level constraint on LLM-agent autonomy**.
Told it is an autonomous process with full agency and *no task*, the policy still enacts a
**chat assistant answering an absent user** (offers tips, then `exit`s) and only acts on its
environment after a forced-grounding nudge. Under one model the constraint is stronger: with no
user input it **manufactures user turns** ("the user says 'no'", mistakes its own words for the
user's) and drives itself by an imagined dialogue. So autonomy declared in the prompt does not
free the policy from the **chat-assistant attractor**; the substrate shapes and bounds what
autonomy is enactable.

## Core Research Question

Are (A) the **two-latch (epistemic × normative / scope-containment) socio-linguistic gating** of
an LLM agent's environmental action and (B) the **substrate-constraint thesis** (conversational/
RLHF scaffolding caps and even fabricates agentic autonomy) **novel and significant** claims,
and what existing work is nearest?

## Detailed Questions

1. **Two-latch gating.** Is there prior work modeling an agent's action as the **conjunction of
   an epistemic appraisal ("there is a problem") and a normative permission ("I may act")**? Look
   in appraisal theory of emotion (Lazarus, Scherer, Frijda — action tendencies), dual-process /
   norm-activation models, deontic reasoning, and the alignment literature on refusal.
2. **Normative latch = scope-containment / intent-preservation.** Who has documented LLMs/agents
   declining to mutate *intact, intentionally-placed* user data while freely fixing *in-scope
   defective* artifacts — i.e. routing around a high-harm mutation (redirect/duplicate vs move)?
   Nearest: corrigibility, side-effects / low-impact AI (Krakovna et al. "Penalizing side
   effects"; Armstrong & Levinstein), conservative agency / attainable-utility preservation,
   "do-no-harm" priors, minimal-impact RL.
3. **Predicate-supply effect.** That a *descriptive-facts-only* environmental note fails to
   trigger action but an *assertive predicate* ("this is a mess") triggers it — nearest work on
   **framing / construal**, problem-finding vs problem-solving, and how LLMs convert state into a
   "task" only when it is linguistically appraised?
4. **Appraisal-readiness gradient.** Our claim that the apparent abstraction ladder (syntax →
   global field) is really a gradient of **distance between the structure and the model's
   pretrained/RLHF defect priors** (broken JSON is universally penalized; a flat dir of 8 files
   is a valid state in thousands of repos) — nearest work on LLM behavior tracking training-data
   priors / "defect-type" sensitivity?
5. **Substrate-constraint thesis.** Who has shown that the **chat/assistant turn-structure +
   RLHF-helpfulness** acts as a constraint/attractor on agentic autonomy — that "you are an
   autonomous agent" prompting does not escape the assistant persona? Look at: persona/role-play
   and "assistant persona" analyses; simulator theory (Janus/"Simulators"); sycophancy; the
   tension between helpfulness/harmlessness and autonomy; agent-vs-chatbot framing; refusal as a
   trained reflex.
6. **Self-fabricated counterpart.** A model with no user input **hallucinating user turns** to
   drive itself — any documented precedent (self-conditioning, role confusion, context
   contamination, "talking to itself", turn-boundary failures)?
7. **Significance:** if structural pull is *structure × the model's latches* and autonomy is
   substrate-bounded, what does that imply for ALife claims about "minimal artificial agency",
   and for the deflation that the effect is "just implicit instruction"?

## Suggested search directions / likely anchors

Appraisal theory — Lazarus 1991; Scherer's component-process model; Frijda (action readiness).
Norm-activation / deontic gating. Alignment & safety — RLHF (Ouyang et al. 2022; Bai et al.
Constitutional AI); refusal & harmlessness; sycophancy (Perez et al.; Sharma et al.);
corrigibility (Soares et al.); side-effects / low-impact (Krakovna et al. 2019; Turner AUP;
Armstrong); conservative agency. Simulators / persona — "Simulators" (janus 2022); role-play &
"Role play with large language models" (Shanahan et al. 2023, Nature); persona/assistant-attractor
analyses; the Waluigi effect (informal). Problem-finding vs problem-solving (Getzels & Csikszentmihalyi).
Framing/construal (Tversky & Kahneman). Indirect prompt injection / environmental injection as the
adversarial mirror of "predicate supply". Autonomy in LLM agents; agent benchmarks where models
wait for permission. Turn-structure / context-contamination failure analyses.

## Relevance filters

Prefer peer-reviewed or widely-cited work; clearly mark influential informal/blog sources
(Simulators, Waluigi) as non-peer-reviewed. Separate **(a)** work that genuinely models a
two-factor (epistemic × normative) action gate or a substrate-autonomy constraint from **(b)**
generic "RLHF changes behavior" observations. Name the single strongest precedent for each claim.

## Deliverables requested

1. A **novelty verdict** for Claim A (two-latch gating; scope-containment) and Claim B
   (substrate-constraint thesis), each: novel / narrowed / pre-empted, with the defensible
   residual.
2. A **ranked adjacent-theory map** per claim (nearest first, "maps / differs").
3. The **single strongest precedent** and the **single strongest objection** for each claim,
   with citations.
4. A judgment on whether Claim B is **publishable as an independent contribution** (we suspect
   it may be the sharpest, most novel result) and where it would land (ALife? an LLM-agent
   venue? a safety venue?).

## What to extract

Exact citations; the closest claim each source makes; whether it frames action as a two-factor
gate; whether it treats the assistant/RLHF substrate as an autonomy constraint.

## Manuscript payoff

Decide whether the two-latch result and the substrate-constraint thesis ship as full-paper
discussion, as their own paper, or as the LBA's forward-looking close; sharpen their novelty
and citations.

## Stop conditions

Stop when both claims have a novelty verdict, a nearest-neighbor map, and the single strongest
precedent + objection each, with citations.
