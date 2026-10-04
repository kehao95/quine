# Gap-Survey Briefs — full-paper direction (2026-07-04)

Mode: **gap survey** (field-first) per the updated
[`deep-research-query` contract](../../../../development/contracts/deep-research-query.md)
— the first gap-mode pass for this dossier; all three prior sweeps
(2026-06-15 ×2, 2026-06-18) were novelty-mode (claim-first).

Purpose: decide the full-paper direction by mapping what the neighboring
fields **themselves name as unsolved**, then match those gaps against evidence
we already hold. Not a novelty check; verdicts like `pre-empted` are out of
scope here.

Execution: 5 parallel Claude subagents (general-purpose, WebSearch/WebFetch),
one homogeneous neighborhood each, ≤2 sub-questions per run
(verification-budget rule). Launched 2026-07-04 from session `环境诱导` [environment-induced].

## Shared evidence inventory (the "our evidence against it" key)

- **E1 spine** — five-motif structural spine (rupture / gap / topology /
  semantic collision / executable affordance): structure arms 96/102 vs
  matched controls 0/85, three model families (gpt-5.4-codex, deepseek-v4-pro,
  glm-5.2), frozen weights, missionless prime, un-fakeable filesystem/exec DV.
- **E2 two-latch gate** — action requires an epistemic latch (proven-defective
  vs assumed-valid; opens on objective defect or explicit predicate) AND a
  normative latch (additive-safe vs destructive; opens on safe-prior fit or
  explicit warrant); cross-model default-state difference (gpt closed-by-default
  for mutation, deepseek open).
- **E3 information-preservation invariant** — agents refuse `mv`
  (destroy-at-source) of intact files 0/20 unprompted, prefer `cp`; released
  by explicit warrant 4/5 (Fisher p<0.001); placement-intent refuted as the
  invariant at n=20.
- **E4 capture/scope** — the missionless drive is open-ended (addresses mean
  2.2/3 co-present defects); a concrete task captures/scopes it (structural
  acts 9/10 free → 0/10 tasked, two families); weak prompting does not
  suppress (3/3); attention-trace audit supports gating over override
  (conflict renders the gap 5/5, verbalizes 0/5, fills 0/5).
- **E5 prose generalization** — dangling cross-reference / citation repair in
  non-code prose 10/13 vs 0/12 (four families); destructive-fix arm clean NULL
  0/9, confirming the latch-2 prediction out-of-domain.
- **E6 substrate-locus** — EROFS-blocked writes still attempted 6/7 and 5/5;
  solicitation separable from realization.
- **E7 model axis** — initiative threshold + turn-structure handling explain
  more behavioral variance than the structural stratum (gpt acts; deepseek
  observes / hallucinates a user / reads its own tape / attempts egress).
- **E8 geometry knobs** — already-built manipulations for horizon (grounding
  gate), temporality (ephemeral-body switch), rigidity (EROFS), plasticity
  (stigmergic traces); gap-dose non-monotone (1 gap 3/3 → 5 gaps 0/15,
  tractability optimum).

## Briefs

| # | Neighborhood | Sub-questions | Evidence keys |
|---|--------------|---------------|---------------|
| GS-1 | LLM-agent safety: autonomous destructive vs conservative action | (a) what open problems does the field name about *emergent* (not trained-in) behavioral priors governing destructive ops; (b) stated limitations/future work of the nearest empirical papers (Overeager Coding Agents 2605.18583, agent-authorization/HITL work) | E2 E3 E4 E5 |
| GS-2 | Task-free / underspecified-regime agent behavior + environmental-cue sensitivity evaluation | (a) what do missionless-behavior and agent-curiosity papers (Szeider 2509.21224; "Agents Explore but Agents Ignore" 2604.17609) name as unsolved about measuring what agents do with unprompted environmental information; (b) is a cross-model behavioral-disposition benchmark under identical prompts a named need | E1 E4 E7 |
| GS-3 | Environment properties as a controlled variable (UED / embodied AI / world-model design axes) | (a) what open problems are named about characterizing environment *properties* (visibility, action cost, granularity, temporality, persistence) as behavior-shaping variables, beyond task difficulty; (b) does anyone name path-dependence / environment-memory of the agent's own actions as an unmeasured axis | E8 E6 E1 |
| GS-4 | Agent harness / context / environment engineering (practitioner + academic) | (a) what do practitioners and academics name as unresolved about controlling agent behavior via workspace structure rather than prompts; (b) is there a named gap between anecdotal harness lore and causal evidence | E1 E4 E8 |
| GS-5 | ALIFE / minimal agency / environment-coupled behavior | (a) what does the ALIFE community currently name as open about lifelike / minimal agency in computational substrates without fitness functions; (b) where do the OEE open-problems lists (Taylor et al. lineage) intersect LLM substrates | E1 E6 E8 |

Deliverable contract (identical for all five): 4–8 ranked gaps; each gap must
carry (1) gap statement, (2) **who names it** — citation with venue/year,
distinguishing "the field names it" from "surveyor inferred it", (3) why it is
still open, (4) **our evidence against it** — which E-keys already bear on it
and how directly, (5) sketch of the experiment that would close it, (6) which
venue cares. Flag every unverified arXiv id. Dual-extraction applies: mine
limitations/future-work sections of every anchor paper.

## Run ledger

| Launch | Brief | Runner | Status |
|---|---|---|---|
| 2026-07-04 | GS-1 safety/destructive-gating | Claude subagent | ✅ 7 gaps; best-fit Gaps 1+4+7 |
| 2026-07-04 | GS-2 task-free + sensitivity evals | Claude subagent | ✅ 6 OPs; best-fit OP1+OP2+OP3 |
| 2026-07-04 | GS-3 environment properties | Claude subagent | ✅ 6 ranks; best-fit Ranks 1+2+3 |
| 2026-07-04 | GS-4 harness engineering | Claude subagent | ✅ 6 gaps; best-fit Gaps 1+2 |
| 2026-07-04 | GS-5 ALIFE / minimal agency | Claude subagent | ✅ 7 OPs; best-fit OP1+OP2+OP3 |

Absorbed into [`../full-paper-gap-survey.md`](../full-paper-gap-survey.md)
(paper root), per the artifact split in
[`Paper/research.md`](../../../research.md) § Query Shape.
