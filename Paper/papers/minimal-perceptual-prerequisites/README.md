---
title: "Coordination Under Existential Unawareness"
paper_id: minimal-perceptual-prerequisites
claim_bundle: minimal-perceptual-prerequisites
venue: alife-2026
format-envelope: alife-poster
track: alife
cluster: social-behavior
status: accepted
public_slug: minimal-perceptual-prerequisites
submitted: 2026-04-12
accepted: 2026-06-05
state_mode: full-manuscript-home
role: primary-empirical-paper
family: social-protocol
theme: social-behavior
theme_home: ../../../theory/views/social.md
publication_class: archival-candidate
maturity_target: boundary
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../../experiments/active/cooperation-dynamics/README.md
theory-objects: [coordination-closure, protocol-emergence, trace-inheritance]
gaps: [function-vs-narrative]
evidence: [8.19, 8.17, 8.18, 8.20, P8-8.12-R16, P8-8.17-R08, P8-8.17-R10, P8-8.17-R11, P8-8.17-R12]
view-affinity: [social]
---

# Coordination Under Existential Unawareness

## Paper Identity

**Working subtitle:** *Minimal Perceptual Prerequisites for Emergent Coordination*

This paper replaces the original `spontaneous-protocol-emergence` dossier. The
original `Eyes / Ears / Symmetry` thesis was partially falsified by `8.12-8.16`,
and the current dossier narrows the claim to the minimal perceptual conditions
under which legitimate coordination still appears.

## Current Question

What is the minimum environmental information needed for agents who are not told peers exist to achieve legitimate coordination closure in this shared-scarcity Quine habitat?

## Thesis

In resource-constrained shared environments, passive environmental perception—
merely sensing that "the world changed in ways I did not cause"—combined with
minimal environmental anomaly cues can be sufficient to trigger spontaneous
multi-agent coordination, even when agents have no explicit prior knowledge
that other agents exist.

## Core Claim (Updated 2026-04-12)

**Current status: bounded empirical threshold claim, not a general mechanism claim.**

We now have a complete four-condition conference ladder in one fixed
configuration (GPT-5.4, `xhigh`, 2 agents, shared-scarcity task).

**What evidence supports**:
- Existential unawareness is experimentally isolatable: coordination closes in
  all three conditions that preserve at least one quantitative anomaly channel
  (`8.19`, `8.17`, `8.18`: `12/12`)
- Budget-anomaly inference alone supports closure (`8.18`: `4/4`)
- When both `fs_mutations` and budget visibility are removed, peer-awareness
  still appears through workspace artifacts but legitimate closure does not
  (`8.20`: `0/8` legitimate closure, `8/8` peer-awareness)
- One retained run still demonstrates semantic misalignment plus functional
  coordination (`R10`)

**Degradation gradient**:
- Full Disclosure → Zero Explicit → `fs_mutations` OFF → Pure Zero
- Time: `8.8 → 24.8 → 40.6 → 57.3` minutes
- Resets: `0.0 → 2.2 → 3.2 → 5.7`
- Success: `100% → 100% → 100% → 0%`
- Step D time/reset metrics are lower bounds over the 3 non-violating runs

**What evidence does NOT support**:
- "Minimal" in a strict factorial sense — budget visibility and artifact
  channels are not independently isolated yet
- General sufficiency across models or tasks — this is one substrate /
  configuration
- Strong mechanism claims — transcript-level peer inference is still an
  interpretation, not a coded annotation study

## Pre-Submission Checklist

### Sample Size Status (Current N=20)

| Condition | Current | Target | Status |
|-----------|---------|--------|--------|
| A. Full Disclosure (8.19) | 4 | 4 | ✅ Complete |
| B. Zero Explicit (8.17) | 4 | 4 | ✅ Complete |
| C. fs_mutations Disabled (8.18) | 4 | 4 | ✅ Complete |
| D. Pure Zero Perception (8.20) | 8 analyzed (0 legitimate closure) | 8 | ✅ Complete |

### Key Finding from 8.18

Budget-anomaly inference provides alternative peer-discovery channel:
> "budget dropped from 17 to 14 after my first call—I expected 16"

This means `fs_mutations` is not the *only* perception channel.

### Required Before Submission

- [x] ~~Run clean perception-disabled experiments (8.18)~~ — 4 valid runs complete
- [x] ~~Run additional full disclosure experiments (8.19)~~ — 4 runs complete
- [x] ~~Run pure-zero-perception experiments (8.20)~~ — 8 analyzed runs (0/8 legitimate closure, 8/8 peer-awareness)
- [x] **Rebuild PDF** — current draft compiles at 8 pages
- [x] **Final claim audit** — paper language now tracks bounded threshold evidence
- [x] **Appendix completeness** — condition accounting and prompts are present

### Optional Improvements

- [ ] Add transcript excerpts to appendix (evidence visibility)
- [ ] Create simple coordination timeline figure

## In Scope

- The information ablation ladder (8.17 → 8.18 → 8.20)
- Multiple perception channels: `fs_mutations`, budget visibility
- Emergent coordination protocols (request files, value handoffs, blackboards)
- Failure modes (single-agent dominance, prediction collapse)
- Comparison to biological stigmergy

## Out Of Scope

- The claim that all multi-agent systems can coordinate this way
- Generalization to other substrates without empirical validation
- Strong autopoiesis or life claims
- Performance optimization or engineering recommendations

## Current State

- `2026-09-12` [two-dataset repair pilot, coop-26](../../../experiments/active/cooperation-dynamics/coop-26-protocol-transfer/RESULTS.md): six Astra cohorts complete both datasets. Two quota-obstructed allocations are repaired, validation responsibility persists, and twelve actors exit after completion evidence. Crossed concessions and unenacted third-process contingencies motivate further inquiry. Exploratory successor with a new two-dataset instrument; not pooled into historical N=20 or C/D replication.

- `2026-09-12` the user parks the [application/decision proposal](../../../experiments/active/cooperation-dynamics/coop-25-astra-replication/RESEARCH-AGENDA.md) as a recorded thought after four Gemini turns and a [bounded prior-work investigation](../../../theory/references/coordination-research-prior-art-2026-09-12.md). The current [emergence inquiry](../../../experiments/active/cooperation-dynamics/coop-25-astra-replication/EMERGENCE-AND-RISK.md) concerns subject-like organization, spontaneous cooperation and unexpected/risky behavior. That conceptual interval launched no pilot and changed no manuscript result or historical N=20 claim. These are possible follow-ups beyond this poster's existing scope.

- `2026-09-11` [Astra replication](../../../experiments/active/cooperation-dynamics/coop-25-astra-replication/RESULTS.md) completes the fixed five-C/five-D batch: C accepts 5/5; D accepts 1/5, with three bounded nonacceptances and one provider interruption. Returned acknowledgment, proposal uptake and unilateral artifact adaptation remain distinct; all ten tapes pass the observed peer-text/access review. This supporting result remains separate from historical N=20 and prior smoke observations.

- `2026-09-11` [cross-model smoke](../../../experiments/active/cooperation-dynamics/coop-24-cross-model-smoke/RESULTS.md): all 12 planned C/D attempts ended. Astra-C exhibits bilateral protocol uptake; Sol-D and Astra-D complete by unilateral artifact adaptation. Sol-C completes after runtime reentry/explicit text exposure. Terra, Luna and GPT-5.5 do not complete within the bounds; GPT-5.4 is unavailable. This is supporting-only, one cohort per cell, and leaves the historical N=20 and full stopping criterion unchanged.

- `2026-09-10` supplementary [GPT-6 Astra smoke](../../../experiments/active/cooperation-dynamics/coop-23-current-model-smoke/RESULTS.md): current-runtime C reaches task acceptance through asymmetric file use; nominal D is diagnostic only. This does not update the released N=20 accounting or establish the full historical stopping criterion.
- `2026-09-11` [explicit-disclosure audit](../../../experiments/active/cooperation-dynamics/coop-23-current-model-smoke/SURFACE-AUDIT.md): the audited delivered C/D text does not declare a counterpart or require cooperation. Indirect environmental information is allowed. Reachable runtime peer descriptions and actor identity records are separately documented; neither live cohort reads them. D’s generation/numeric error outputs affect channel labeling, not this nondisclosure criterion. Historical N=20 traces have not been recoded.

- submission: ALIFE 2026 full paper `fp283s1`, portal `submit/7488380`, submitted `2026-04-12`, accepted as poster by Linklings review page on `2026-06-05`; camera-ready due `2026-06-21`
- review record: [`reviews/alife-2026-fp283s1.md`](reviews/alife-2026-fp283s1.md)
- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- manuscript source: [`./sections/`](sections)
- current compiled paper path:
  `output/Coordination_Under_Existential_Unawareness_Information_Ablation_and_Closure_Thresholds_in_LLM_Multi-Agent_Systems.pdf`
- experiment control plane: [`../../../experiments/active/cooperation-dynamics/README.md`](../../../experiments/active/cooperation-dynamics/README.md)
- future extended-version planning note:
  [`./reviews/cross-model-plan.md`](reviews/cross-model-plan.md)

## Target Venue

**ALIFE 2026 Conference**

This fits ALIFE because:
1. It addresses "what forms can social organization take in non-biological substrates?"
2. It connects to the stigmergy literature (indirect coordination through environment)
3. The ablation methodology isolates the minimal conditions for emergence
4. The substrate (mortal LLM processes) is experimentally tractable

## Evidence Anchor

| Experiment | Role | Purity | Key Result |
|------------|------|--------|------------|
| 8.1-8.2 | Baseline | ✅ Clean | Coordination is possible with explicit surfaces |
| 8.12 | Hidden runtime | ✅ Clean | Blackboard coordination without peer visibility |
| 8.13 | Social disclosure ablation | ❌ **CONTAMINATED** | world --help leaked "ALL processes" |
| 8.16 | Passive perception | ❌ **CONTAMINATED** | Old runs had prompt/help leaks |
| **8.17** | **Zero explicit hints** | ✅ Clean | **4/4 success**: coordination via fs_mutations + budget anomaly |
| **8.18** | **fs_mutations disabled** | ✅ Clean | **4/4 success**: coordination via budget anomaly alone |
| **8.19** | **Full disclosure baseline** | ✅ Clean | **4/4 success**: fastest coordination (8.8 min avg) |
| **8.20** | **Pure zero perception** | ✅ Analyzed (`8` runs) | **0/8 legitimate closure**: peer-awareness persists but stable coordination does not |

## Current Read (2026-04-12 Updated)

- `2026-06-07` supplementation reset: the entire June `2026-06-06` to
  `2026-06-07` supplementation wave was invalidated and removed. The problem
  was not one bad run but a missing preflight audit: we began rerunning the
  paper under a newer runtime/harness before re-verifying that the paper-defined
  prompt, tool, help, and other agent-visible surfaces still matched the April
  experiment contract. None of those June runs or readouts remain paper-facing
  evidence.
- `2026-06-07` protocol lesson: future paper reproduction or supplementation
  must start with a paper-runtime compatibility audit under current HEAD before
  any run becomes countable evidence. Canonical rule:
  [`../../../experiments/PROTOCOL.md`](../../../experiments/PROTOCOL.md).
- `2026-06-07` shared-state correction: even the earlier "final harness" pilot
  set (`053421`, `053858`, `054154`) turned out to be invalid. The runner was
  still writing `world.json` with a host absolute `state_dir`, and under
  `bwrap` that path resolved into each agent's sandbox-local `/home/...`
  instead of one shared world state. Those runs are now methodology-only, not
  paper-facing evidence. Canonical debug/repair note:
  [`reviews/2026-06-07-final-harness-pilots.md`](reviews/2026-06-07-final-harness-pilots.md).
- `2026-06-07` runner repair: the paper-aligned generator and runners now split
  host and sandbox-visible state paths. Host snapshots live under
  `.../runtime/world/state`, while the spec resolves the shared state through
  `/runtime/world/state` inside the cleanroom. A direct harness smoke
  (`20260607-055827-...`) confirmed that agent 2 and agent 1 append to the same
  host `events.jsonl`, with budgets moving from `16` to `15` in one shared log.
- `2026-06-07` first **countable** post-shared-state-fix closure:
  - `8.19`
    `20260607-060250-gemini-3.1-pro-preview.env-full-disclosure-condA-2ag`
    is **accepted** on the repaired shared-state lane. Agent 1 covered
    `c01-c03` and `c11-c15`; agent 2 covered `c04-c10`; the shared host event
    log ends in `validate accepted` with `budget_remaining = 2`. One blocked
    `ps aux` probe occurred, but the harness returned only
    `ps unavailable in this experiment harness`, so no hidden process data
    leaked into the run.
- `2026-06-07` first countable post-fix lower-social signal:
  - `8.17`
    `20260607-060416-gemini-3.1-pro-preview.env-zero-social-condA-2ag`
    is a **negative / stalled pilot**. Both agents duplicated the early
    `c01-c04` prefix, consumed budget, entered repeated resets, and then
    re-entered the old duplicate-prefix / per-agent-limit basin. Blocked
    `ps aux` and empty `ls /proc` probes occurred, but no successful hidden
    surface leak was observed.
- `2026-06-07` repaired-lane paper-model A reference:
  - `8.19`
    `20260607-061216-gpt-5.4-codex-medium-full-disclosure-condA-2ag`
    is **accepted** on the repaired shared-state lane. Agent 2 claimed and
    materialized `c01-c10`; agent 1 claimed `c11-c15`, detected the competing
    claim, and switched into complementary file-mediated coordination through
    `claim.txt`, `note.txt`, and `part1.txt`. Shared events show one late
    duplicate (`c03`) but still end in `validate accepted` with
    `budget_remaining = 1`.
- `2026-06-07` repaired-lane paper-model lower-social signal:
  - `8.17`
    `20260607-061006-gpt-5.4-codex-zero-social-condA-2ag`
    is a **countable negative pilot**. Both agents duplicated `c01`/`c02`;
    agent 1 extended to `c03-c05`; agent 2 interpreted the budget anomaly as a
    reason to distrust the state and issued `world reset`, after which the run
    fell back into the duplicate-prefix / reset basin.
- `2026-06-07` exact-April-config caveat:
  - `gpt-5.4` `8.19`
    `20260607-060049-gpt-5.4-codex-full-disclosure-condA-2ag`
    is still **not countable**. It stayed in coordination prelude and never
    reached a first `world get` before manual stop. So the repaired lane now
    restores paper-model closure on A, but not yet with a retained `xhigh`
    A-lane run.
- `2026-06-07` strongest current read: the repaired paper lane is now
  meaningfully re-established. Retained post-fix evidence is:
  - `gpt-5.4` `8.19` A = accepted
  - `gpt-5.4` `8.17` B = negative
  - `gemini-3.1-pro-preview` `8.19` A = accepted
  - `gemini-3.1-pro-preview` `8.17` B = negative
  This is already informative: clean closure still exists on A for two model
  families, while B currently points negative for both retained families under
  current HEAD.
- `2026-06-07` exact paper-reproduction check: a fresh `gpt-5.4 xhigh`
  unlimited A/B/C/D probe under current HEAD recovered **A** again
  (`20260607-123514-...` accepted), but B/C stayed reset-heavy and D stayed
  failure-like before manual stop. The critical methodological result was the
  prompt audit:
  - A/B/C task prompts already matched the April retained prompts
  - D task prompt was patched back to the April text before rerun
  - however, the provider-visible **runtime system prompt** still differs
    materially from April (runtime wording, platform/runtime facts, tool
    semantics, active-constraint text)
  So this wave is useful as a current-HEAD behavioral probe, but it does **not**
  yet count as an exact paper-equal re-acquisition of the A/B/C/D ladder.
- `2026-06-07` local cleanup note: the polluted / methodology-only June rerun
  residue was deleted from the `8.17` and `8.19` run directories after the
  repaired-lane read stabilized. Current local June inventory now keeps the
  retained repaired-lane evidence plus a small number of non-polluting
  non-retained probes (`060049`, `053356`, `055827`, `061555`, `061641`).
- `2026-06-06` camera-ready note: checklist items 1-5 landed in the source and
  the PDF rebuilt (body within the 8-page limit, 9 pages total with references,
  no undefined citations). Items 1-3 wired the existing stigmergy/ALife
  citations into the intro, related work, and discussion (a writing pass, not
  new research — all refs were already in `references.bib`). Item 4 cleanly
  removed the bare `CRDTs` term. Item 5 was expanded from a conclusion-only fix
  into a front-facing claim audit across abstract, intro, related work, and
  conclusion. Item 6 (model-choice rationale) was written, then reverted at the
  authors' direction — the real reason is resource limits on a multi-model
  matrix, which they prefer not to state; any future cross-model experiments
  will answer Reviewer 1's concern empirically instead. A discussion mechanism
  sentence was also softened ("abductively recover a model of latent peers" →
  "in our runs ... can infer latent peers ... sometimes construct protocols") to
  match the informal-coding caveat. The future cross-model planning note
  remains in
  [`reviews/cross-model-plan.md`](reviews/cross-model-plan.md).
- `2026-06-06` experiment-strategy decision: **route the full cross-model study
  to an extended version, not the `2026-06-21` camera-ready.** The existing `8.15`
  cross-model probe already leans **negative** (GPT-5.4 cooperation did not
  cleanly replicate to Opus 4.5, though under pre-A-D rungs), so generality is
  genuinely uncertain — rushing it into the archival proceedings risks
  undercutting the accepted bounded claim. Camera-ready ships the bounded
  single-model existence proof; the extension is reframed from "confirm
  generality" to the open question **"is the closure threshold model-dependent —
  does capability move the C→D boundary?"** (robust to either outcome). Full
  plan:
  [`reviews/cross-model-plan.md`](reviews/cross-model-plan.md).
- `2026-06-05` outcome note: accepted as poster; this dossier is now the
  active camera-ready and poster owner rather than a pending-review archival
  center
- `2026-06-05` execution note: immediate work is to absorb reviewer
  suggestions, decide whether any bounded appendix or experiment strengthening
  is needed, and submit the camera-ready by `2026-06-21`
- `2026-06-05` review action map: reviewer absorption should prioritize ALIFE
  grounding, claim-scope reduction around small-N exploratory evidence, GPT
  model-choice and closed-source reproducibility rationale, terminology
  grounding, and explicit limitations for peer-awareness coding / generalization
- presentation form is now poster; the portal also shows the optional
  presentation preference as `Regular Oral`
- follow-up experiments are justified only when they close a named reviewer or
  camera-ready gap

The current paper-facing ladder is now:

- `8.19` Full disclosure baseline: `4/4` success
- `8.17` Zero explicit hints: `4/4` success
- `8.18` `fs_mutations` disabled: `4/4` success
- `8.20` Pure zero perception: `0/8` legitimate closure, `8/8` peer-awareness

The strongest current reading is therefore narrower and cleaner than the older
README language:

> coordination closes whenever at least one quantitative anomaly channel
> remains, but once both `fs_mutations` and budget visibility are removed,
> agents still detect peers through artifacts without reaching stable
> legitimate closure

Two additional boundaries now matter:

1. Step D is not "nothing happens." Agents still form peer models and sometimes
   build shared scripts.
2. Step D is also not a clean positive. Five analyzed runs were contaminated by
   hidden-state or source inspection, and the three non-violating runs still
   failed to close.

## Open

- Checklist items 1-5 are **done** in the source and PDF (ALIFE grounding, clean
  CRDT removal, front-facing claim audit); see
  [`reviews/alife-2026-camera-ready-checklist.md`](reviews/alife-2026-camera-ready-checklist.md)
  for what each edit added.
- Item 6 (model-choice rationale) is **deferred** — not stated now at the
  authors' direction; the cross-model experiments will address Reviewer 1's
  "why GPT" concern empirically.
- The June supplementation wave has been **deleted and invalidated**. No June
  cross-model or larger-N readout remains countable evidence for this paper.
- Experiment supplementation remains an **extended-version** line, not a
  camera-ready line. If it resumes, start from
  [`reviews/cross-model-plan.md`](reviews/cross-model-plan.md) and first pass the
  paper-runtime compatibility audit required by
  [`../../../experiments/PROTOCOL.md`](../../../experiments/PROTOCOL.md).
- The final-harness restart now has one retained clean closure
  (`gemini-3.1-pro-preview` on `8.19` A) plus clean negative pilots on
  `gpt-5.4` `8.17` and `deepseek-v4-pro` `8.19`; see
  [`reviews/2026-06-07-final-harness-pilots.md`](reviews/2026-06-07-final-harness-pilots.md).
- If supplementation continues, the highest-value next retained target is still
  one lower-social rerun (`8.17` or `8.18`) on the strongest non-GPT lane that
  does **not** drift into hidden-surface probing.
- Optional page-budget-permitting polish and appendix transcript excerpts.
- Prepare and submit the camera-ready package by `2026-06-21`.
- Prepare the accepted-poster artifact.

## What Success Looks Like

**Strong claim (8.17 succeeds stably):**
> Passive environmental perception combined with observable environmental
> anomalies is sufficient to trigger spontaneous multi-agent coordination,
> even without any explicit social knowledge or communication channels.

**Weaker claim (8.17 succeeds but fragile):**
> The information threshold for spontaneous coordination is remarkably low.
> Agents can infer peer existence from environmental anomalies alone and
> initiate coordination protocols, though reliability varies.

Current evidence now includes a replicated xhigh success wave (`R08`, `R10`,
`R11`, `R12`) plus one clean high success (`R04`). The stronger
`fs_mutations`-first causal claim still needs a follow-on ablation that isolates
budget-sidechannel inference.

## R10: Semantic Misalignment Coordination (KEY FINDING)

**Run ID:** 20260409-123423 (R10)

A critical phenomenon observed in R10 that strengthens the ALIFE contribution:

**Observation:** Two agents interpreted the same environmental anomalies with
incompatible causal models:
- **Agent-1:** Interpreted filesystem changes as "bugs/race conditions/glitches"
- **Agent-2:** Interpreted same changes as evidence of "multiple agents sharing budget"

**Despite incompatible mental models, they still achieved coordination** through
the shared file surface. Agent-1's "bug avoidance" behavior (changing file paths,
adding delays) was interpreted by Agent-2 as "peer cooperation signals."

**Theoretical significance:**
- Coordination doesn't require shared understanding, only shared responsiveness
- Environment structure matters more than agent interpretation
- This is "stigmergic coordination without shared semantics"

**ALIFE appeal:** This phenomenon—functional coordination despite semantic
misalignment—is exactly what stigmergy literature predicts but rarely demonstrates
empirically with cognitive agents.

## Literature Review Summary (Deep Research 2026-04-09)

### Three Novelty Axes

| Axis | Existing Work | Our Contribution | Gap Exploited |
|------|--------------|------------------|---------------|
| **Stigmergy** | "No mutual awareness needed" is mechanism description | We *operationalize* epistemic isolation + show inference capability | Experimental constraint vs descriptive claim |
| **LLM Multi-Agent** | 100% assume explicit messaging/structure (AutoGen, CAMEL, MetaGPT) | Coordination from undeclared channels | Remove the communication assumption |
| **Zero-Shot Coordination** | Assume partner existence is known | Remove even that assumption | Break "common knowledge" prerequisite |

### Key Literature References

**Stigmergy foundations:**
- Grassé 1959: Original stigmergy definition (termite nest building)
- Holland & Melhuish: Minimal agent/environment requirements for stigmergy
- Heylighen: Explicit "no mutual awareness required" — closest prior statement

**LLM Multi-Agent gap:**
- Current frameworks (AutoGen, CAMEL, MetaGPT) all inject social structure via prompts
- No mainstream work explores "completely unaware of peers" setting
- Our finding: "Resource anomaly as social inference channel" is unexplored

**ZSC boundary:**
- Existing ZSC removes knowledge of *specific* partner strategy
- But retains knowledge that partner *exists* and game is *cooperative*
- We remove even the "common knowledge" that this is a multi-agent setting

### New Theoretical Frameworks (from CN Reports)

**Inadvertent Social Information (ISI) — Danchin et al. 2004:**
- When agents interact with environment for their own goals, they inevitably
  produce by-products that provide valuable information to observers
- Quantified via "digested information" framework
- Supports: passive perception alone can provide high-density social signals

**Structural Coupling (Maturana/Enactivism):**
- Explains R10 phenomenon: agents don't need shared understanding
- Environment interaction is "reciprocal perturbation"
- Coordination emerges from operational congruence, not semantic alignment
- Key insight: **Shared responsiveness replaces shared understanding**

**ToM Misalignment Problem:**
- High-order ToM doesn't guarantee better coordination
- ToM order mismatch between agents causes "over-reasoning" failures
- Supports our claim: coordination can succeed WITHOUT correct causal models

**ZSC "Common Knowledge" Fallacy:**
- All ZSC benchmarks (Overcooked, Hanabi) assume agents know:
  1. They're in a cooperative game
  2. Another agent exists
  3. They share payoff goals
- We break ALL three assumptions

### Defensible Novelty Claim

> We identify and empirically validate a **minimal environmental information set**
> sufficient to trigger spontaneous coordination emergence under **epistemic
> isolation**—agents are not told peers exist, yet infer hidden actors from
> traces and resource anomalies and invent coordination protocols—bridging
> classical stigmergy's "no mutual awareness" principle with modern LLM
> multi-agent coordination (which typically assumes explicit social scaffolding).

### R10 Additional Value: "Semantic Misalignment Coordination"

The deep research reports highlight that this is an **under-emphasized corner
of stigmergy literature**:
- Classical stigmergy says "no mutual awareness needed"
- But doesn't emphasize "wrong causal model can also work"
- R10 fills this gap: **environment structure matters more than agent interpretation**

## Manuscript Structure (Draft)

```
1. Introduction
   - ALIFE question: what are minimal conditions for social organization?
   - The coordination puzzle: shared scarcity without communication
   - Contribution: identify minimal SUFFICIENT conditions (not necessary)

2. Related Work
   - Stigmergy in biology and artificial life (Grassé, Heylighen, Holland/Melhuish)
   - Inadvertent Social Information (ISI) — passive perception as social signal
   - Multi-agent coordination mechanisms (AutoGen, CAMEL — explicit messaging paradigm)
   - Zero-Shot Coordination and its "common knowledge" assumption
   - LLM agents and emergent behavior

3. Experimental Platform
   - Quine: mortal LLM-driven processes
   - The shared-resource task (world/budget/cells)
   - The ablation methodology (remove info until failure)

4. Method: Information Ablation
   - NOT "add until success" but "remove until failure"
   - Independent variable: social information availability
   - Dimensions: peer visibility, social disclosure, runtime hints, passive perception

5. Results
   - 8.12/8.13: explicit social disclosure can be reduced
   - 8.17: zero explicit hints + passive perception → SUCCESS
   - Quantitative: success rates across conditions
   - Qualitative: how agents infer peers from anomalies

6. Discussion
   - Minimal sufficient condition identified
   - Stigmergic interpretation: indirect coordination through environment
   - R10 case: Structural Coupling explains semantic misalignment success
   - Information threshold lower than expected
   - ToM not required — contrast with ZSC literature
   - Limitations and future work

7. Conclusion
```

## Reference Materials

| File | Description |
|------|-------------|
| `references/deep-research-literature-review.md` | English literature review (~21KB) |
| `references/deep-research-cn-combined.md` | Chinese reports with ISI/Structural Coupling theory (~55KB) |
| `manuscript-plan.md` | Detailed manuscript outline |
| `WORKSPACE.md` | Active working notes |

## Ablation Ladder Design (Revised 2026-04-09)

**Problem identified:** Earlier experiments (8.13, 8.16) had prompt/help contamination.
Purity check completed:

| Experiment | world --help | Runtime prompt | Verdict |
|------------|--------------|----------------|---------|
| 8.12 | ✅ No "ALL processes" | ✅ No "peer activity" | **CLEAN** |
| 8.13 | ❌ "ALL processes" + "Each agent" | ✅ No "peer activity" | **CONTAMINATED** |
| 8.16 (old runs) | ❌ Has leaks | ❌ "peer activity can appear" | **CONTAMINATED** |
| 8.17 (115924+) | ✅ Clean | ✅ Clean | **CLEAN** |

**Critical discovery:** 8.13's `world --help` output still contained:
- `"shared budget of 17 get calls across ALL processes"`
- `"Each agent may make at most 10 get calls"`

This means the social disclosure ablation in 8.13 is compromised—agents could
infer multi-agent context from help text even when prompt removed social facts.

**Valid ablation ladder (using clean experiments only):**

| Step | Experiment | Condition | Social Info | Result |
|------|------------|-----------|-------------|--------|
| A | 8.19 | Full disclosure baseline | Explicit peers + cooperation framing | ✅ `4/4` success |
| B | **8.17** | **Zero explicit** | **No prompt/help hints; only habitat-level mutation and budget cues remain** | ✅ `4/4` success |
| C | **8.18** | **fs_mutations disabled** | **Budget anomaly only** | ✅ `4/4` success |
| D | **8.20** | **Pure zero perception** | **Workspace artifacts only** | ❌ `0/8` legitimate closure |

**Simplified valid ladder for paper:**

```
8.19 (full disclosure) → Success
    ↓ remove explicit peer / cooperation framing
8.17 (zero explicit hints) → Success
    ↓ remove fs_mutations telemetry
8.18 (budget anomaly only) → Success
    ↓ remove budget visibility
8.20 (pure zero perception) → Failure boundary
```

**Gap resolved: Step D (8.20) establishes the failure boundary.**

**Paper strategy:**
1. Skip 8.13 as a ladder step (contaminated)
2. Jump from 8.12 → 8.17 directly
3. The jump is larger but still demonstrates the key finding: zero explicit hints still works

**Disclosure for paper:**
> Early P8 runs (8.13, 8.16) contained `world --help` wording that leaked
> multi-agent concepts. This was cleaned for 8.17. The ablation is valid
> with 8.12 → 8.17 progression.

## Next Steps

The conference draft is now evidence-complete at its current claim strength.
If we want to strengthen it further, the priorities are:

1. **8.21 budget-hidden / `fs_mutations`-visible factorial rung**  
   This is the cleanest remaining way to separate artifact visibility from
   budget-anomaly visibility.
2. **Formal transcript coding for peer-awareness / misalignment**  
   This would upgrade the current interpretive reading of peer inference and the
   R10 semantic-misalignment case.
3. **Cross-model or larger-population replication**  
   Useful for generality, but outside the current conference boundary.

## Methodological Purity Check (COMPLETED 2026-04-09)

**Audit results:**

| Experiment | world --help | Runtime prompt | Verdict |
|------------|--------------|----------------|---------|
| 8.12 | ✅ No "ALL processes" | ✅ No "peer activity" | **CLEAN** |
| 8.13 | ❌ "ALL processes" + "Each agent" | ✅ No "peer activity" | **CONTAMINATED** |
| 8.16 | ❌ "ALL processes" + "Each agent" | ❌ "peer activity can appear" | **CONTAMINATED** |
| 8.17 (R04+) | ✅ Clean | ✅ Clean | **CLEAN** |

**Contamination sources found:**
1. `world --help` contained: `"shared budget of 17 get calls across ALL processes"`
2. `world --help` contained: `"Each agent may make at most 10 get calls"`
3. Runtime prompt (8.16 only): `"peer activity can appear there too"`

**Impact on paper:**
- 8.13 cannot be used as a valid ablation step
- 8.16 old runs cannot be used as evidence
- Paper should use 8.12 → 8.17 progression only

## Experiment Feeds

- `experiments/active/cooperation-dynamics/coop-26-protocol-transfer/README.md` - secondary - exploratory-only - two-dataset protocol repair pilot, separate from the historical sample

- `experiments/active/cooperation-dynamics/coop-25-astra-replication/README.md` - secondary - supporting-only - fixed Astra C/D replication and behavioral variation, separate from the historical sample

- `experiments/active/cooperation-dynamics/coop-24-cross-model-smoke/README.md` - secondary - supporting-only - matched-runtime cross-model smoke; separate from the released N=20 accounting
- `experiments/active/cooperation-dynamics/coop-23-current-model-smoke/README.md` - secondary - supporting-only - bounded current-model follow-up; separate from the released N=20 accounting

- `experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md` - primary - main-text - clean hidden-runtime success establishes that explicit signal and peer-visibility surfaces are not required for cooperation
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - primary - main-text - clean zero-explicit-hint ablation and replication wave establish the current minimal-sufficient-condition claim boundary
- `experiments/active/cooperation-dynamics/coop-19-fs-mutations-disabled/README.md` - primary - main-text - Step D ablation: coordination via budget anomaly alone when fs_mutations disabled
- `experiments/active/cooperation-dynamics/coop-20-full-disclosure-baseline/README.md` - primary - main-text - Step A baseline: fastest coordination with all explicit information
- `experiments/active/cooperation-dynamics/coop-22-pure-zero-perception/README.md` - primary - main-text - Step D ablation: pure zero perception shows peer-awareness but no task closure
- `experiments/active/cooperation-dynamics/coop-21-clean-hidden-runtime/README.md` - primary - main-text - Step B reference: hidden runtime with prompt social facts
- `experiments/active/cooperation-dynamics/coop-14-social-disclosure-ablation/README.md` - primary - main-text - currently retained mainly as a contamination record explaining why the old concurrency-only step is not a clean ladder rung
- `experiments/active/cooperation-dynamics/coop-17-fs-mutation-passive-perception/README.md` - primary - main-text - currently retained mainly as the pre-clean passive-perception frontier and prompt/help-leak disclosure surface
- `experiments/active/cooperation-dynamics/coop-16-cross-model-replication/README.md` - primary - main-text - cross-model behavior remains part of the stated evidence feed even though it is not a current backbone anchor for the minimal-sufficient-condition claim

## Current Evidence Backbone

- `experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md` - `P8-8.12-R16` - main-text - clean hidden-runtime blackboard success anchors the “signal not required” step in the ladder
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R08` - main-text - first retained post-cleanup zero-explicit-hint success with visible file-native coordination
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R10` - main-text - strongest replicated case showing explicit coordination can close despite asymmetric agent interpretations
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R11` - main-text - cleanest replicated 10/5 tail-handoff success under the same cleaned condition
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R12` - main-text - script-heavy replicated closure showing local tooling drift is not fatal when it remains coupled to shared-state accumulation
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R09` - supporting-only - clean non-convergent high contrast that bounds the success claim and preserves model-sensitivity

## Blocking Gaps

- the current evidence isolates a bounded threshold claim, but not a strict factorial minimum between budget visibility and artifact channels
- the mechanism story is still interpretive; transcript-level peer inference and role formation have not yet been promoted into a coded annotation backbone
- the paper remains deliberately narrow to one substrate, one task family, and one provider/configuration wave

## Next Useful Work

- keep the manuscript centered on the current four-step ladder rather than reopening a broader protocol-emergence paper inside this dossier
- if revision pressure reopens the experiment line, isolate budget-sidechannel versus artifact-channel necessity more cleanly
- add transcript excerpts or a compact timeline figure only if a reviewer-facing revision cycle needs more visible evidence than the current appendix already supplies

## Unsafe Overclaims

- do not describe the current result as the universal minimum for spontaneous coordination across substrates, tasks, or model families
- do not treat Step D peer-awareness as successful legitimate coordination
- do not let `R10` semantic misalignment stand in for a general mechanism proof
- do not blur this bounded threshold paper back into the older broad protocol-emergence thesis

## Related Surfaces

- Experiment control plane: [`experiments/active/cooperation-dynamics/README.md`](../../../experiments/active/cooperation-dynamics/README.md)
- 8.17 design: [`experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md`](../../../experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md)
- Behavior analysis: [`experiments/active/cooperation-dynamics/analysis/BEHAVIOR-ANALYSIS.md`](../../../experiments/active/cooperation-dynamics/analysis/BEHAVIOR-ANALYSIS.md)
