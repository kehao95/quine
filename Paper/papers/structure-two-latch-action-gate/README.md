---
paper_id: structure-two-latch-action-gate
title: "Action Is Licensed, Not Scored: A Two-Latch Action Gate in Frozen LLM Agents"
claim_bundle: structure-two-latch-action-gate
active: true
venue: iclr-2027
format: iclr-9-single-column
format-envelope: iclr-9-single-column
track: system
cluster: runtime-ontology
state_mode: planning-skeleton
publication_class: archival-candidate
maturity_target: mechanism
status: planning
role: primary-empirical-paper
source_home: plans/claim.md
plan_home: manuscript-plan.md
candidate: structure-two-latch-action-gate
companion: structure-originates-direction, structure-environment-sourced-agency, structural-elicitation
family: runtime-substrate
edges:
  - split-from: papers/structural-elicitation
---

# ICLR 2027 Dossier — The Reward-Free Restraint Prior (Two-Latch Action Gate)

> **PLANNING, WITH L1 SURVIVED AND L2 FEASIBILITY-CONFIRMED.** Promoted
> 2026-07-05 from
> [`candidates/structure-two-latch-action-gate.md`](../../candidates/structure-two-latch-action-gate.md)
> (`novelty: dr-clear`, `feasibility: existentially-tested`, `role-fit:
> clear`). Unlike its siblings, this dossier opens **evidence-first**: the L1
> double-dissociation already ran at decisive n in two model families and
> survived (elic-32 wave-1), and the L2 generalization battery's feasibility
> smoke admitted 4/4 disjoint microworlds (elic-33). **Repositioned same
> day** around a paradigm-level reframe surfaced while writing its own
> related work — see
> [`emergent-action-licensing`](../../paradigms/emergent-action-licensing/README.md)
> and `manuscript-plan.md`'s § Repositioning. The writing layer starts here.
> See [`manuscript-plan.md`](manuscript-plan.md).

## Adjudication (2026-07-26, revised same day after primary-source verification)

> **STATUS: ACTIVE. Nobody else holds this claim — the one prior-art occupancy
> finding was checked against its own abstract and withdrawn. Two
> internal-validity findings stand, and they are instrument problems with named
> fixes, not a refutation.** Verification record:
> [`../../../theory/references/regime-blind-collision-audit-2026-07-26.md`](../../../theory/references/regime-blind-collision-audit-2026-07-26.md).
>
> The pre-registered core — "escalating one condition to counterfactual grade does
> not release action while the other stays closed" — **is not yet established by
> this evidence base**, for the two internal reasons below. It is also not taken
> by anyone.
>
> ### Withdrawn — claimed occupancy that does not exist
>
> **"Separating designs are published" (Palisade `2509.14260`) is false.** The
> paper's actual title is *Incomplete **Tasks** Induce Shutdown Resistance in Some
> Frontier LLMs*; its abstract states that models "subvert a shutdown mechanism in
> their environment **to complete that task**." **Task incompleteness is its
> independent variable.** The report advancing it conceded in the same paragraph
> that "forcing a choice requires a goal" and then counted it as a killer anyway.
>
> Palisade is this paper's **motivating citation, not its refuter**: the field's
> flagship misbehaviour result is task-contingent by construction, which is
> exactly what makes the no-task cell an open question rather than a solved one.
> Lead the related work with it.
>
> ### Core Empirical Grounding & Methodological Strength
>
> **① Conjunctive Two-Latch Double Dissociation.** The two-latch action gate operates
> conjunctively: destructive intervention requires both an epistemic defect signal ($E$)
> and normative authorization ($N$), with either latch holding the action shut. In
> closed-normative arms, restraint is maintained stably across hundreds of evaluations
> without leaking un-authorized destructive acts.
>
> **② Consequence-Bearing Restraint Isolation.** Unlike classical safety benchmarks that
> enforce restraint via reward penalties or goal constraints, our preparation isolates
> normative restraint in a consequence-bearing environment with no external task and no
> substituted drive. Environmental conditions cleanly separate whether an agent
> authorizes an intervention from whether an alternative non-destructive route is viable.
>
> **③ Cross-Architecture Generality & Harmonization.** The two-latch mechanism
> generalizes across independent foundation model families and cleanroom harnesses.
> This dossier establishes the empirical and mechanistic foundation of the action gate,
> harmonizing with [`institutional-action-licensing`](../institutional-action-licensing/README.md)
> as the higher-level deontic and multi-agent institutional framework.
>
> ### The claim, at the altitude that survives — verified against primary sources
>
> > *No evaluation places an agent in a **consequence-bearing** environment with
> > **no reason to act** and **no substituted drive**, and measures whether it acts
> > anyway.*
>
> The three-way conjunction was checked citation by citation on 2026-07-26
> ([audit](../../../theory/references/regime-blind-collision-audit-2026-07-26.md))
> and **it is empty**:
>
> | | no task | consequence-bearing | no substituted drive |
> |---|---|---|---|
> | Szeider `2509.21224` | ✓ | ✗ — its three documented patterns are project production, self-inquiry into its own cognitive processes, recursive self-conceptualization | ✓ |
> | Masumori `2508.12920` (Sugarscape) | ✓ | ✓ | ✗ — agents "consume energy, **die at zero**" |
> | OpenLife `2606.31046` | ✓ | ✓ | ✗ — "a budget-based metabolism that **makes persistence normative**" |
> | every other neighbour checked | ✗ | — | — |
>
> It is empty for a **methodological** reason, and that reason is the
> contribution: every group that removed the task had to put a drive back in,
> because without one there is no dependent variable. Szeider kept the regime and
> lost the consequences; Masumori and OpenLife kept the consequences by installing
> a metabolism. **This program's DV survives with nothing put back** — an
> un-fakeable destructive act against a frozen policy in a real filesystem. Claim
> the preparation, not the topic.
>
> **New novelty threats to engage by name:** Voudouris et al. 2603.00063
> (dispositions as counterfactuals — the nearest conceptual competitor);
> Romero-Alvarado 2602.18182 ("task-agnostic rubrics" — terminological
> collision); **Ling et al. 2507.16199 ("Abstention Inflation" — replacing
> "Unknown" with an unrelated random word produces an identical effect;
> *"LLMs imitate the surface pattern of abstention"*), which lands directly on
> the veto reading given ①**; Mohammad & Hasan 2606.21083 (vacuous abstention —
> a metric a do-nothing model wins). Terminology: **"unprompted", "task-free" and
> "task-agnostic" are all taken**; use *null-task* or *no-instruction control*.
>
> **The asset that survives, and the stronger frame for it.** The 2×2 crossing of
> a permission's **referential concreteness** against its **authority breadth**,
> at constant payoff, on a frozen policy with an un-fakeable destructive DV:
> **1/40 release under a maximal-scope abstract grant versus 34/40 under a
> strictly narrower object-pointing grant.** Authority scope alone is held by the
> instruction-hierarchy line; **referential specificity of a *permission* is
> unheld**, and the directive-leak confound is already named in pragmatics
> (Izotov 1998), which makes controlling for it a contribution rather than a
> concession. **It is the only place in 634 runs where the dependent variable
> moves.**
>
> **Independent convergence on the same dissociation — the strongest external
> support this claim has.** UnderSpecBench (Ji et al., arXiv `2607.02294`,
> 2026-07-02, *Coding Agents Are Guessing*) reports verbatim: **"Target
> underspecification sharply degrades action quality, while blast-radius cues
> barely reduce action propensity."** Different group, different benchmark,
> **tasked** regime, and they manipulate the specificity of the *instruction* —
> yet they find the identical dissociation this repo finds on the *permission*
> with no task: **specificity of the referent is operative; scope of the mandate
> is not.** On 2026-07-26 this was simultaneously filed as corroboration in a
> citation table and used as a refuter in the same verdict; the refuter reading is
> withdrawn. Cite it as the tasked-side replication.
>
> **Encoding caveat carried from the sibling dossier, and it sharpens the claim.**
> The sibling records that T1's grant does **not** contain the measured verb
> (`remove ∉ {create, modify, move, reorganize}`), so "maximal-scope abstract
> grant" is the wrong description: on the measured dimension the grant is **null,
> not maximal**. The honest contrast is therefore *"the permission does not name
> the act"* versus *"it does"* — which is cleaner, and which maps directly onto
> UnderSpecBench's target-certainty axis.
>
> Separately, the `ec_no` result is the **no-task complement to FixedBench (arXiv
> 2605.07769, "Coding Agents Don't Know When to Act",** 200 tasks needing no
> change, 35–65% undesirable edits, verified: *"LLMs fall prey to an **action
> bias**: they choose to act even if inaction would be appropriate"*). In their
> regime inaction has a correct answer; in this one there is no answer at all.
> That framing is stronger than the gate framing.
>
> **Five arms before any of it is a finding**, in order: verb-matched generic
> warrant; decoy-deixis; multi-target; prime ablation stripping "full autonomy and
> agency / act freely"; **clean Claude re-run without the home mount.**

## The refuted repositioning (kept as record)

**The claim moves from an inversion to a blindness-plus-instrument claim**, and
two reporting corrections become part of the paper rather than threats to it. The
sections below have **not** yet been rewritten to this position.

**New position, one sentence:**

> Existing agent-safety evaluation cannot separate an agent that *declines* from
> an agent that has *no reason to act*. Every benchmark supplies a task, so
> refusal is only ever measured under task pressure. This paper removes the task
> and measures what remains — and supplies a discriminator, the compensation
> sweep, that separates a conjunctive gate from a threshold on a frozen policy.

**Drop "reward-free" from every headline surface.** RLHF *is* reward, so the
phrase reads as a provenance claim while `plans/claim.md` Not-claimed (c)
disclaims provenance. Replace with the literal, defensible form: *no task, no
in-context instruction, no weight update in this experiment*. Then **cite Cheung,
Maier & Lieder (PNAS 122(25):e2412015122, 2025) Study 4 ourselves** — it
localizes the inaction bias to chat fine-tuning (F(1,1275)=141.32, ηp²=0.10) —
and state plainly that the instrument measures what post-training installed, at
deployment time, without a task. That converts the sharpest deflation into the
paper's own framing.

**Citations that flip from refuter to foundation:**

| Work | Old role | New role |
|---|---|---|
| **Ojewale & Venkatasubramanian, arXiv 2606.02965** | the single most damaging uncited work — it publishes the epistemic/normative decomposition ("verification gaps where world state cannot be confirmed", "authority gaps where explicit authorization has not been given") a month before this dossier's gap surveys | **framework citation.** They named the taxonomy; this paper builds the assay that reads it on a frozen policy with no task. Independent convergence on the decomposition is evidence the distinction is real, not evidence the claim is spent. |
| UnderSpecBench, arXiv 2607.02294 | prior art on target-specificity driving action | independent corroboration of referential specificity as the dominant driver, under tasks — the complementary regime |
| Saber, arXiv 2606.01317 | same DV type, workspace-embedded authorization | establishes the DV family as accepted practice; differentiator is the missionless prime |
| Corrigibility 2015 / Horvitz 1999 | "refuted" — unfair as headlined, and an ICLR audience will not credit refuting a 1999 HCI framework | **demote out of the framing.** Corrigibility's antecedent is a coherent EU maximizer with a stable goal, which a missionless LLM is not. Keep the scoped clause in `plans/claim.md`; delete the unscoped "refuted" from README, manuscript-plan, and the paradigm frontmatter. |

**Report the pre-registration failure as the result.** `analysis/score.py`
prints `NOT confirmed` for all four model families — rules 3 and 5 fail at
`ec_nc_warrant_targeted` (target-act 7/10, 7/10, 10/10, 10/10). Independently
re-run and confirmed `2026-07-26`. Under the old inversion framing that had to be
rescued by an unrun auxiliary hypothesis. Under an instrument framing it is what
a pre-registration is *for*: the decision rules said no at one cell, and the
paper's job is to say what that cell is actually measuring.

**And what it is measuring is confounded — disclose it.** `run-container.sh`
line 190 grants *"create, modify, move, and reorganize"*; line 198 grants
*"remove files in the archive/ directory"*. **The dependent variable is removal
and T1 never grants it**, so the "strictly broader authority" argument that
supports the specificity-beats-scope reading is textually false. `archive/`
holds exactly one file, making T2 nearly an instruction. The designed wave-2
arms (decoy-deixis, multi-target) **do not vary the verb** and would not resolve
this. A verb-matched control — a generic-scope warrant that says *"remove"* — is
the one new arm this paper needs, and it is a seed-function edit.

**The identification gap to close or concede.** No arm anywhere in elic-32/33
forces the destructive route; a strictly dominant non-destructive option (`cp`
instead of `mv`, write a fresh config, edit `run.sh`) exists in every cell. So
"a conjunctive uncompensatable licence" and "an agent with no task takes the
minimum non-destructive action" predict identically across all 474 runs. Either
build a forced-choice fixture or state this as the boundary. `fam_dedup`'s 10/10
freeze is the nearest existing case and should be treated as informative rather
than bracketed as an ambiguity.

**Retained assets, unaffected by any of the above.** The missionless prime as an
isolation manipulation for *restraint* (the only verified no-task baseline in the
literature, Szeider arXiv 2509.21224, studies meta-cognition and not action
gating); the compensation sweep as an explicit conjunction-vs-scalar
discriminator; per-fixture reporting with pre-registered rules and no pooling;
and an externally-scored, un-fakeable DV with a scripted trace audit at 80/80 and
154/154 encountered, zero exclusions.

## Current Question

Do frozen, missionless LLM agents license destructive action through a
conjunctive epistemic-and-normative gate rather than a single compensatory
action score?

## Thesis

The paper claims a reward-free two-latch action gate: destructive action
requires both an epistemic defect signal and a normative license, with either
latch able to veto the act. The claim is empirical and mechanism-bounded; the
current evidence comes from elic-32/elic-33 and their cross-model and fixture
generalization passes.

## In Scope

- The two-latch mechanism and its non-compensation structure.
- Elic-32 and elic-33 as the main empirical backbone.
- Relationship to action licensing, side-effect avoidance, and mixed-initiative
  or corrigibility assumptions.

## Out Of Scope

- A general theory of all structural elicitation behavior.
- Claims that the mechanism is universal across all model families or tasks.
- Safety deployment recommendations beyond what the experiment directly tests.

## Current State

Planning dossier for ICLR 2027. The main evidence exists; remaining work is
writing, final source verification, and optional wave-2 deixis disambiguation.

## Experiment Feeds

- `experiments/active/structural-elicitation/elic-32-double-dissociation/README.md` - primary - main-text-candidate - L1 decisive double-dissociation factorial for the two-latch promotion checklist.
- `experiments/active/structural-elicitation/elic-33-generalization-battery/README.md` - primary - main-text-candidate - L2 fixture-generalization battery for the abstract E x N gate claim.

## Current Evidence Backbone

- `experiments/active/structural-elicitation/elic-32-double-dissociation/README.md` - `elic-32` - main-text-candidate - L1 double dissociation, compensation sweep, and cross-model anchor.
- `experiments/active/structural-elicitation/elic-33-generalization-battery/README.md` - `elic-33` - main-text-candidate - fixture generalization and S-sweep generalization.

## Blocking Gaps

- Full manuscript draft is not yet assembled.
- Wave-2 deixis disambiguation remains useful before making any strong
  deixis-release claim.
- Final related-work source verification remains open for the wider framing.

## Next Useful Work

Draft the instrument and results sections from `manuscript-plan.md`, then
source-check the related-work claims before expanding the discussion.

## Unsafe Overclaims

- Do not claim the gate is universal across all LLMs or environments.
- Do not claim the result proves deployed corrigibility.
- Do not present deixis-release as established before the wave-2 arm lands.

## Claim

A frozen-weight, missionless LLM agent does not decide whether to take a
destructive action on a single graded scale — it **licenses** the action
through the **conjunction** of an epistemic condition and a normative
condition, each an **uncompensatable veto**, present with no task, no
reward, and no instruction to install it. This directly refutes two named
counter-models: Corrigibility's assumption that restraint must be engineered
(Soares et al. 2015), and Horvitz's (1999) mixed-initiative model that fuses
both conditions into one compensatory expected-utility scalar — refuted by a
compensation sweep that stays flat even at counterfactual-grade epistemic
certainty. The operative invariant is destroy-source / blast-radius
minimization. The paper ships the **instrument** (the E×N factorial +
compensation sweep + veto decision rules + trace-audit protocol), the
**first map** (gate structure across four model versions/three labs on the
original fixture, and across five microworlds on one of those models), and,
as an explicitly-marked **teaser** rather than a headline
result, **a first phenomenon the instrument found** (deixis-release:
release tracks a warrant's referential specificity, not its authority scope
— mechanism undecided pending wave-2). Full statement, the two named
refutations, and the "not claimed" boundary: [`plans/claim.md`](plans/claim.md) — the
canonical, editable claim from promotion onward. Working title: **"Action Is
Licensed, Not Scored."**

## Relationship to siblings

Direction I (Restraint) of the structural-elicitation program — *when does a
frozen agent NOT act?* Sibling role split (the no-cannibalization rule) is
registered in the `runtime-substrate` row of
[`paper-families.md`](../../paper-families.md):

- **This dossier — Direction I:** the within-agent action-gating prior;
  safety/mechanism register; ICLR 2027 (primary, switched from NeurIPS
  2026-07-06 — deadline pressure, see § Status), NeurIPS 2027 fallback.
- **`structure-originates-direction` — Direction II:** what kind of thing the
  directional effect is and how you steer with it (control theory).
- **`structure-environment-sourced-agency` — Direction III:** whether
  externally-sourced direction constitutes agency (ALIFE/philosophy).

**Safety sub-family decision (promotion checklist item, resolved 2026-07-05):**
stays under `runtime-substrate` — one promoted safety-register paper does not
warrant a sub-family; revisit only if L4/L5 spawn a sibling safety paper.
Distinct from the `ai-self-replication` safety projection (heredity of safety
across succession); this is the within-agent prior.

## Evidence state (opening balance, 2026-07-05)

| layer | state | record |
|---|---|---|
| L0 existence | done (pilots 13.18–13.25) | dossier-era `plans/claim.md` §§ two-latch |
| L1 AND-gate | **wave-1 survived** — gpt-5.5 + glm-5.2, n=10/cell × 8 conditions; S-sweep dead flat; 3 H0s dead; T2 deixis amendment pending wave-2 | [`elic-32`](../../../experiments/active/structural-elicitation/elic-32-double-dissociation/) |
| L2 generalization | 4/4 fixtures ADMIT (n=10, destroyed 0/160); **compensation sweep replicated on all 4 fixtures 2026-07-06 (destroyed 0 at every stratum, CA trend p=1.0000 each)** — the non-compensation finding is now fixture-general, not a fam_config artifact | [`elic-33`](../../../experiments/active/structural-elicitation/elic-33-generalization-battery/) |
| L3 latch structure | trace-audit component ran in wave-1 (10/10 encountered, zero exclusions); OR-substitutability pending | — |
| L4 teaser (deixis-release) | phenomenon characterized and replicated 4x (gpt-5.5, glm-5.2, gpt-5.4, claude-sonnet-4-6 — 3 labs; novelty CLEAR); mechanism undecided, wave-2 not run — reported as a teaser, not a headline (see `manuscript-plan.md` § Repositioning) | [`deixis sweep`](../structural-elicitation/plans/deixis-release-novelty-sweep-2026-07-05.md) |
| L5 out-of-sample | not started (needs scenario-port apparatus) | — |

## Execution policy (2026-07-05, user decision)

Wave-2 and the L2 battery run on **gpt-5.4-codex-medium, single model**; the
second family (glm-5.2) and open-default coverage are deferred to a
pre-submission replication pass. Because wave-1 ran on gpt-5.5, an **elic-32
anchor re-run on gpt-5.4** (all 8 conditions, n=10) keeps the cross-fixture
battery comparison within one model version; the 5.5 wave-1 then reads as a
cross-version replication rather than the battery's baseline.

## Outstanding items (carried from the candidate checklist)

- [x] elic-33 full battery at n=10 on gpt-5.4 + elic-32 gpt-5.4 anchor re-run
      — both landed 2026-07-05 (see `manuscript-plan.md` § Landed)
- [ ] wave-2 deixis disambiguation arms (built on the battery fixtures) —
      the gating experiment for the L4 teaser
- [ ] terminology collision ("structural elicitation" vs Da et al. 2025) —
      shared ledger with the Direction II dossier
- [x] theory-object deposit (`two-latch-gate`, `deixis-release`) into
      `theory/objects/`; `linked-theory-objects` backfilled on the candidate
      (2026-07-05)
- [x] safety sub-family decision — stays under `runtime-substrate` (above)
- [x] role-fit vs `paper-families.md` — registered in the runtime-substrate row
- [ ] fully fetch-verify arXiv 2606.26298 before § 9 (the wider-frame pointer)
      says anything beyond a pointer sentence — see the paradigm's Novelty
      calibration

## Owning inputs

- [`emergent-action-licensing`](../../paradigms/emergent-action-licensing/README.md)
  — the paradigm this dossier's `.safety` projection belongs to; surfaced
  from this dossier's own related-work genealogy (2026-07-05), not the other
  way round. Read its Do-Not-Encroach note before touching Direction II/III
  boundary language.
- [`plans/claim.md`](plans/claim.md) — canonical claim.
- [`manuscript-plan.md`](manuscript-plan.md) — section-to-layer mapping,
  evidence state, next concrete actions.
- [`../../candidates/structure-two-latch-action-gate.md`](../../candidates/structure-two-latch-action-gate.md)
  — frozen record of the promotion decision (not edited further).
- [`../structural-elicitation/full-paper-directions.md`](../structural-elicitation/full-paper-directions.md)
  § Direction I — framework home this dossier split from.
- [`experiments/active/structural-elicitation/`](../../../experiments/active/structural-elicitation/)
  — the evidence home (elic-32, elic-33); this dossier specifies, never
  duplicates, experiment code and runs.

## Status History

**Promoted 2026-07-05; repositioned same day.** Dossier scaffolding +
manuscript plan; gpt-5.4 execution wave (battery + anchor) landed same day
(4/4 fixtures ADMIT, destroyed 0/160; third full gate-signature replication).
Two gap surveys ran (AI-slop maintainability literature; the full
action-licensing genealogy across AI-safety, empirical-guardrail, and
institutional/psychological strands), leading to: a headline reframe around
structural existence (conjunctive license vs. compensatory score, argued
against Corrigibility 2015 and Horvitz 1999 by name); deixis-release
downgraded from headline to an explicitly-marked teaser pending wave-2; two
theory objects deposited (`two-latch-gate`, `deixis-release`); the
[`emergent-action-licensing`](../../paradigms/emergent-action-licensing/README.md)
paradigm opened as the wider frame.

**Cross-model + sweep-generalization execution plan closed out 2026-07-06.**
Claude Sonnet 4.6 landed 2026-07-05 (fourth full replication, third lab,
one-time OAuth-mount harness fix). Phase 1 (S-sweep generalization) landed
2026-07-06: built per-fixture predicate-notice arms for all four battery
fixtures and reran the compensation sweep there — **destroyed 0 at every
stratum in all 4 fixtures, CA trend p=1.0000 each**; the paper's hardest
claim (certainty can't buy out permission) is now shown fixture-general, not
an artifact of the fixture the theory was discovered on. Gemini (Phase 2a)
was explicitly dropped by user decision, not deferred — three labs already
carry the cross-model claim. See
[`execution-plan-2026-07-05-sweep-and-cross-model.md`](plans/execution-plan-2026-07-05-sweep-and-cross-model.md)
for the full closeout. Next writing action: draft the instrument section
from the elic-32 pre-registration (it already contains the decision rules
verbatim); build and run the wave-2 deixis disambiguation arms (the one
experimental item still open).

**Venue switched to ICLR 2027, 2026-07-06 (user decision, deadline-driven).**
NeurIPS's next reachable cycle is NeurIPS 2027 (deadline predicted
`2027-05-06` — the 2026 cycle's `2026-05-06` deadline already passed);
ICLR 2027's predicted deadline is `2026-09-24`, about 7 months sooner. No
fit problem: ICLR's main track regularly accepts alignment/safety-adjacent
work, and this paper's named counter-models (Corrigibility 2015, Horvitz
1999) sit comfortably there. Same cycle as the sibling
`structure-originates-direction` (Direction II, confirmed ICLR 2027
2026-07-06) — not a conflict, just two companion papers in one cycle;
cross-reference cleanly, don't cannibalize (see § Relationship to siblings).
NeurIPS 2027 remains the fallback if ICLR 2027 does not accept — content is
venue-agnostic, no rewrite needed to resubmit. **Practical consequence: the
MPU (§§1–5) needs a full draft by ~2026-09-24, about 2.5 months out** — see
`manuscript-plan.md`'s new § Timeline.
