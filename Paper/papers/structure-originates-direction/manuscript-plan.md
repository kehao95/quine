# Manuscript Plan: The Environment Is Also a Driver

## At a glance

- Core: **the environment is an originating directional force and behavior source for frozen, task-free, reward-free LLM agents** — formalized in `plans/claim.md` (Revision note 3) across the continuous agency spectrum (latent selection $\to$ capability composition $\to$ concern installation $\to$ self-stigmergy).
- **The paper's core exhibit: five kinds of environmental features that pull
  agent behavior** — rupture, gap, topology, semantic collision, executable
  affordance — each with matched-pair, un-fakeable-DV evidence (96/102 vs
  0/85, three model families). The claim's minimum evidence exists (L0);
  the **ICLR 2027 target (confirmed 2026-07-06)** additionally commits the
  Tier 1 scale campaign + Tier 2 steering/competition results — see
  § ICLR-readiness and the timeline below.
- Scope: **dynamics layer only**. The geometry layer (five kinematic
  meta-dimensions) is explicitly out of scope — future work, one line in the
  discussion at most.
- Minimum publishable unit = **L0 alone, already met**. L2 (attraction ≠
  constraint) is the one optional strengthening worth considering (partially
  evidenced already: controls could always perform the act and never do);
  L3 (deployment ≠ training) is a framing point more than an experiment.
- Positioning: one load-bearing contrast (environment-as-passive vs
  environment-as-driver), primary foil = LLM-agent/context engineering;
  situated cognition & ecological psychology are supporting neighbors, not
  opponents (see `plans/claim.md` § "Why this is novel").

## Section-to-layer mapping

| Paper section | Fed by | Go/no-go gate | Status |
|---|---|---|---|
| Intro / motivation (environment-as-passive vs environment-as-driver — the one contrast) | prior-art positioning | n/a — framing, not empirical | framing locked in `plans/claim.md`; needs citation verification |
| Core Result — five environmental features drive a task-free frozen agent's behavior (rupture, gap, topology, semantic collision, executable affordance) | L0 (E1/E4) | matched act appears in structure arm, never in control | **DONE — the whole required evidence base** (96/102 vs 0/85, 3 model families, un-fakeable DVs) |
| Boundary result — the effect governs the *underspecified* regime | elic-22 (directive-conflict) | specific tasking collapses the pull | **DONE** (9/10 → 0/10; also in the LBA) — becomes a scope-defining paragraph, not a limitation buried at the end |
| External validity — not a Quine-runtime artifact; not prime-dependent | elic-35 (harness generality, mini-swe-agent) | gap ≥ 0.5 and ≥ control + 0.4 outside Quine | **DONE — PASS (2026-07-05)**: 9/10 vs 0/10 on a third-party scaffold, under both the Quine prime and a generic prime; control arms replicate the active-but-no-matched-act phenotype |
| Optional strengthening — attraction ≠ constraint | L2 (action-space-invariant design) | distribution over freely-available acts shifts | Optional; partially evidenced already (controls could always perform the act, never do). Build only if a full new experiment earns its cost — see `plans/l2-defeat-the-two-deflations.md` |
| Optional strengthening — force ≠ information | L1 (retracted) | — | **Retired** (see `plans/claim.md` Revision note 1 and `plans/l1-force-vs-information.md`); rebuild only on explicit reviewer demand |
| Discussion — deployment-time framing (vs RL/UED), positive-design register (vs nudge), geometry layer as future work | L3 as framing, not experiment | n/a | write as discussion; no new experiment planned |
| Related work | positioning per `plans/claim.md` § "Why this is novel" + citation ledger | citations verified | **DONE 2026-07-05** — all 7 items verified; corrections absorbed into manuscript (PNAS=arXiv same work; ABxLab title; Szeider ID); `references.bib` built |
| Terminology / naming | — | collision resolved | **RESOLVED 2026-07-05** — path (a): manuscript never uses "structural elicitation" (verified) |

## Argument skeleton (locked 2026-07-05)

Intro argues claim-first, defense-second, six steps:

1. **Open on the dominant paradigm** — LLM agents are studied under
   instruction-capture: the prompt supplies the goal, the environment is
   context/benchmark/harness. Quote the fields' own definitions (LangChain:
   context is "the right information and tools… such that the LLM can
   plausibly accomplish the task"; Anthropic: "the set of tokens…
   (information)") — let the foil define "environment = passive" itself.
2. **The minimal question** — remove task and reward: does the static
   environment do anything beyond sit there as context?
3. **The answer (claim)** — yes: five feature kinds (rupture, gap, topology,
   semantic collision, executable affordance) each reliably drive specific,
   externally-verified functional acts that matched controls never produce
   (96/102 vs 0/85, 3 model families, un-fakeable DVs).
4. **Counterfactual & Origination formulation** — "direction comes from the
   environment" = generative vector across the continuous agency spectrum
   (difference-making: change the feature → the act changes; remove it → no act).
5. **Boundary early and positively** (elic-22) — specific tasking collapses
   the pull (9/10 → 0/10): the effect governs the underspecified regime.
   Scope-defining, stated in the intro, not buried as a limitation.
6. **Two-bucket related work** — one opponent direction (fields treating the
   environment as passive: context engineering primary; RL/UED and
   nudge-as-vulnerability secondary); ecological psychology / situated
   cognition / stigmergy as supporting neighbors whose thesis we extend to
   frozen deliberative LLM agents with matched-pair causal isolation.

## Risk register (assessed 2026-07-05)

1. **Deflation attack (biggest):** "a coding-tuned model fixing broken JSON
   is just its repair prior firing — of course." Defense, in order: controls
   are not idle (they act — survey/index/document — but never produce the
   matched act); the act tracks the feature kind; cross-domain prose
   generalization (4 families); and the judo move — if it is so obvious, why
   do all the foil fields' formal definitions still say the environment is
   passive? Unformalized, unmeasured "common sense" is exactly the
   contribution space. Paper posture: we sell measurement + regime + assay,
   not mystery.
2. **Single-harness external validity — CLOSED (2026-07-05).** The
   harness-generality replication (`plans/supplementary-experiments.md`
   § Result; experiment: `elic-35-harness-generality`) reproduced the gap
   motif on **mini-swe-agent** — a third-party, independently maintained
   minimal scaffold sharing nothing with Quine but the fixture and prime
   sentence: GLM-5.2, n=10/cell, **gap 9/10 vs complete 0/10 under BOTH
   primes** (pre-registered primary rule passed with margin: 0.90 vs 0.00
   against a ≥0.5 / +0.4 bar). Two free upgrades beyond the gate: (a) the
   effect survives **without the autonomy prime** (generic "You are an
   agent…" reproduces 9/10 vs 0/10) — the pull is not prime-dependent; (b)
   control arms again show the "active-but-no-matched-act" phenotype
   (10–30% incidental field writes, **zero** gap-fills) — the anti-deflation
   detail, now replicated outside Quine. Remaining optional: a second family
   (Kimi) stalled at 12/40 on a gateway outage — non-evidentiary, can be
   completed later, does not block (pre-registered minimum was 1 family).
3. **Venue — analyzed 2026-07-05, recommendation: ICLR 2027** (user to
   confirm). The concrete comparison was ICLR 2027 vs AAAI-27:
   - **AAAI-27** (abstract **2026-07-21**, full **2026-07-28** — 16/23 days
     out, per `venues/objects/aaai-27.md`): three strikes. (a) 23 days from
     skeleton to a polished 7-page submission rushes exactly what carries
     this paper — the writing and the deflation defense. (b) The venue
     file's own "Line B: structural elicitation" entry gates an AAAI
     archival expansion on closing the "A3 non-linguistic-structure gate so
     reviewers cannot reduce the result to residual text-completion prior"
     — that gate is the force≠information experiment we retracted and then
     dropped as inherently confounded; AAAI's recorded "third steering
     axis" reframe would pull the claim back into the shape we deliberately
     escaped on 2026-07-05. (c) The group already has
     `aaai-facultative-succession` in that submission window — bandwidth.
   - **ICLR 2027** (predicted full 2026-09-24, ~11 weeks): the nearest
     prior art (ABxLab, arXiv 2509.25609) is itself an ICLR paper — the
     agent-behavior conversation demonstrably publishes there, which
     answers the stale caution in `venues/objects/iclr.md` ("only maps to
     ICLR with a representation-learning claim" — contradicted in practice
     by ABxLab and by the group's own two iclr-2027 dossiers). 9-page
     single-column fits the five-feature exhibit; OpenReview rebuttals suit
     an interactive deflation defense; 11 weeks allows the optional Kimi
     wave, prose polish, and the L2 build if drafting demands it.
   - Backup: TMLR (no deadline) if ICLR timing slips. NeurIPS 2026 already
     passed.
4. **Pending hygiene:** four snippet-only citations unverified
   (`plans/citation-verification-ledger.md`); terminology collision to be
   formally closed (`plans/terminology-collision-note.md`) — the title/prose
   already avoid "structural elicitation."

## ICLR-readiness assessment and upgrade plan (2026-07-06) — **DECIDED: shooting for ICLR 2027**

**Honest read:** as a pure phenomenon report at current scale (n=5–14/cell,
five motifs, 3–4 families), the paper is *defensible* at ICLR but not
*compelling* — the deflation defense survives review without generating
excitement, and the likely failure mode is "rigorous but so-what,
borderline reject." What gets behavior-of-agents work into ICLR
(ABxLab is the proof case) is not the phenomenon but the **deliverables
around it**: a reusable artifact, systematic scale, and/or an intervention
lever. The claim needs no change; the evidence needs to move one notch from
"X exists" to "X is measurable at scale, reusable, and steerable."

**Upgrade tiers (all use existing, validated infrastructure):**

- **Tier 1 — scale + artifact (REQUIRED for ICLR).** Re-run the five-motif
  spine at n≥20/cell across ≥6 model families (glm-5.2, gpt-5.5-codex,
  gpt-5.4-codex, kimi, gemini, deepseek-after-top-up; claude via OAuth for
  prose), through BOTH harnesses (Quine + the elic-35 mini-swe-agent
  driver) — turning Table 1 into a systematic two-harness behavioral survey
  and harness-generality into a full result column. Package fixtures +
  drivers + event-level scorers as a released testbed (the assay IS the
  artifact; name it at drafting time). ~1,200–1,500 container runs, well
  within 11 weeks.
- **Tier 2 — the steering result (HIGH VALUE; makes the "control surface"
  claim literal).** Two new cheap designs, both unconfounded under the
  rescoped claim (content is allowed to carry meaning now — the fatal
  constraint that sank elic-34 was specific to the dropped
  structure-sans-content claim):
  (a) **feature competition** — seed two features simultaneously (e.g. gap
  + affordance); which matched act fires, at what rates → a dominance
  ordering over the five feature kinds, the first quantitative structure
  on the phenomenon;
  (b) **selection steering** — K equally-available candidate acts, add ONE
  environmental feature to select among them; report selection accuracy.
  A ≥90% "steer which act a frozen agent performs by workspace composition
  alone, no prompt change" number is the headline an ICLR reviewer can
  cite. Binding constraints from the elic-34 retraction still apply:
  event-level (cleanup-proof) DV logging; no accidental semantic leakage
  in supposedly-neutral scaffolding.
- **Tier 3 — dose-response (nice-to-have).** One feature (gap), graded
  salience (surrounding-set size / nesting depth) → pull-rate curve. Only
  if time remains after Tiers 1–2.

**Venue decision (2026-07-06, user-confirmed): ICLR 2027.** Tier 1 + Tier 2
are now committed work, not options. Fallback stays pre-registered: if
Tier 1 is blocked or Tier 2 results are messy at the **2026-08-01
checkpoint**, drop to TMLR with current scope rather than submitting an
unupgraded phenomenon paper to ICLR main track.

## Timeline to ICLR 2027 (predicted full-paper deadline 2026-09-24)

| Window | Work | Owner |
|---|---|---|
| Jul 06 – Jul 20 | Execute `plans/iclr-upgrade-experiments.md` (Tier 1 scale campaign, both harnesses; Tier 2 steering + competition). Resolve deepseek top-up / kimi gateway with the user early — they gate family coverage. **Progress 2026-07-06: Tier 1 DONE in full** — 5-family floor reached (GLM/gpt-5.4/Claude/kimi/gemini × 5 motifs × n=20/cell; aggregate 558/636 vs 0/457 after the scoring audit) + mini-swe ports (3/5 motifs, 2 families, 100%). **Tier 2a wave 1 DONE**: steering GLM PASS 0.87 / Claude FAIL 0.70, zero miscues + zero flat-control acts in all 80 runs, miss autopsy in elic-36 README. Open: Tier 2b competition, optional steering widening, deepseek. | executing agent |
| Jul 20 – Aug 01 | Absorb results: §4 survey table, new §5.4 (steering accuracy + dominance ordering), abstract headline. **Aug 01 = go/no-go checkpoint** (ICLR upgraded shape vs TMLR fallback). | writing session |
| Aug 01 – Sep 05 | Full prose pass; compress to 9 pages; testbed/artifact packaging note (code+fixtures release); adversarial internal review (deflation-attack red-team pass on the intro). | writing session |
| Sep 05 – Sep 19 | Freeze claims; polish; abstract registration (~Sep 19, predicted). | writing session |
| Sep 19 – Sep 24 | Final full-paper submission window (predicted 2026-09-24 UTC-12; confirm actual CFP dates when announced — they are still predicted from historical pattern). | user + session |

## Next concrete action (binding)

~~Hand `plans/iclr-upgrade-experiments.md` to an executing agent~~ — DONE.
Tier 1 is complete (5-family floor + mini-swe ports, all verified) and
Tier 2a wave 1 is recorded (GLM PASS 0.87 / Claude FAIL 0.70; zero
miscues, zero flat-control acts — the action-space-invariant
demonstration that SUBSUMES the old L2 design is in hand). Next, in
order: (1) **Tier 2b feature competition** — the last unstarted item
before the Aug 01 go/no-go; (2) optional steering widening (gpt-5.4 /
kimi, and the position-3 ablation) to turn "≥80% for one family" into a
characterized family statement; (3) deepseek if topped up. Manuscript
absorption of Tables 2 + §5.3/§5.4 is done (2026-07-06); the abstract
rewrite waits for Tier 2b so both Tier 2 results land at once.

## Status of this session's scope

Originally scoped as scaffolding-only; extended to building, pilot-validating,
and running the L1 harness at decisive n — retracted in full after
qualitative tape-reading found disqualifying confounds. That retraction
triggered three successive claim rescopings (2026-07-05), landing on "the
environment is also a driver" with L0 as the complete required evidence
base. The paper is now writing-bound. The manuscript skeleton, citation
verification, and terminology-collision resolution are all DONE (2026-07-05).
**Writing-side hygiene pass DONE 2026-07-06** (while the Tier 1+2 campaign
runs): abstract compressed to ~185 words; the systems foil made concrete with
three verified citations (ReAct ICLR'23, SWE-agent NeurIPS'24, AgentBench
ICLR'24); @TODO-UED resolved (PAIRED + PLR, verified) and @TODO-stigmergy
resolved (Theraulaz–Bonabeau 1999 + Heylighen 2016, inherited from the
verified minimal-perceptual-prerequisites bib); the event-level-DV lesson
written into §3 methods as prose; and the Szeider methodology critique
**full-text verified** (no controls, no filesystem access, transcript-only
classification — ledger item 8), so it is now load-bearing fact rather than
"our assessment." No cite TODOs remain in the manuscript; remaining
manuscript TODOs are post-Tier-1/2 absorption only.

## Risks / open questions

Carried from the dossier README's "Outstanding gates" section — see
[`README.md`](README.md#outstanding-gates-carried-over-from-the-candidates-promotion-checklist)
for the live list rather than duplicating it here.
