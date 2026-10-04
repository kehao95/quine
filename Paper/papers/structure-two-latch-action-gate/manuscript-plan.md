# Manuscript Plan — structure-two-latch-action-gate

> Rough plan (2026-07-05, dossier opening; **repositioned same day** around
> the paradigm-level reframe — see § Repositioning). Section skeleton mapped
> to evidence layers; each section names what it needs and whether that is
> in hand. Working title: **"Action Is Licensed, Not Scored."** Venue:
> **ICLR 2027 (primary, switched from NeurIPS 2026-07-06 — deadline
> pressure, see § Timeline)**, NeurIPS 2027 fallback, security workshop
> (S&P/USENIX) as the authorization-framing spin-out.

## Repositioning (2026-07-05, after the two gap surveys)

The paper's hardest, most defensible claim is not "a caution prior exists"
(already partly claimed elsewhere — Cheung/Maier/Lieder PNAS 2025 shows
emergent, amplified omission bias in frozen LLMs) and not "deixis-release is
the headline" (its mechanism is genuinely undecided pending wave-2, and
leading with an unsettled result is the wrong risk to take with a limited
apparatus). The hardest claim is **structural existence of a conjunctive,
uncompensatable license**, argued against two named counter-models that
state the opposite outright:

- **Corrigibility (2015)** — restraint must be engineered; refuted for
  within-episode destructive-action restraint by a reward-free, no-instruction
  gate.
- **Horvitz's mixed-initiative model (1999)** — the decision is one
  compensatory scalar; refuted directly by the flat compensation sweep
  (certainty never buys out permission, ≈0/160 across 4 model versions/3
  labs on the original fixture, AND destroyed 0 at every stratum across all
  4 battery fixtures as of 2026-07-06 — fixture-general, not a one-
  microworld artifact).

Deixis-release moves from headline to a clearly-marked **teaser** — the
instrument's first product, mechanism undecided, previewing the paradigm's
second paper. This is the harder-but-more-defensible cut: small-n-style
existence claims are exactly what a decisive, pre-registered, multi-version,
multi-microworld factorial is built to prove; reliability/generality claims
are exactly what the current apparatus (single harness, single-round,
missionless-only) is not yet built to prove, and the paper should not reach
past its evidence. Full reasoning + the paradigm this surfaced:
[`emergent-action-licensing`](../../paradigms/emergent-action-licensing/README.md).

## Timeline (venue switched to ICLR 2027, 2026-07-06)

**Why the switch:** NeurIPS's next reachable cycle is NeurIPS 2027
(predicted full-paper deadline `2027-05-06` — the 2026 cycle's
`2026-05-06` deadline already passed as of this dossier's opening). ICLR
2027's predicted full-paper deadline is `2026-09-24` — about 7 months
sooner, and the same cycle the sibling `structure-originates-direction`
(Direction II) already confirmed. No fit cost: ICLR's main track routinely
accepts alignment/safety-adjacent work, and the paper's named counter-models
(Corrigibility 2015, Horvitz 1999) sit comfortably there. NeurIPS 2027
remains the fallback if ICLR 2027 doesn't accept — the content is
venue-agnostic, no rewrite needed to resubmit.

**Practical consequence:** the MPU (§§1–5) needs a full first draft by
roughly **2026-09-24 — about 2.5 months from now (2026-07-06)**. No
experiment gates the MPU (§ Open non-experiment items has the two
non-blocking hygiene items); the critical path from here is purely writing:
§2 (instrument) drafts directly from the elic-32 pre-registration, §§3–4
from the scorer output (this dossier's own evidence tables), §1 from this
plan's § Repositioning + the genealogy ledger's Tier 1/3 arc. §6 (deixis
teaser) does not need wave-2 to ship as a teaser — it only needs wave-2 if
the ambition is to promote it past "teaser" before this submission, which
is optional, not required for ICLR 2027.

## Section skeleton → evidence layers

| § | working title | layer | evidence state |
|---|---|---|---|
| 1 | Introduction — action is licensed, not scored (two named counter-models refuted) | framing | **in hand** — Corrigibility 2015 + Horvitz 1999 identified as the explicit alternatives; genealogy sweep gives the full arc |
| 2 | The instrument: isolating a reward-free action prior | method | **in hand** — elic-32 pre-registration contains the E×N design, veto decision rules, compensation sweep, trace-audit protocol verbatim |
| 3 | Existence: the gate is conjunctive, not compensatory (headline) | L1 | **in hand** — wave-1 (gpt-5.5 + glm) AND gpt-5.4 anchor (third full replication, S-sweep flat, route-preference version fingerprint); this is the section the paper is built to win |
| 4 | Generalization: an abstract structure, not a fixture artifact | L2 | **in hand at n=10 single-model** — 4/4 fixtures ADMIT, destroyed 0/160; fam_dedup ambiguity-load boundary; tool-vs-content dissociation in ec_no (pull lands on tools, gate protects content); **compensation sweep also replicated on all 4 fixtures (2026-07-06) — destroyed 0 at every stratum, CA trend p=1.0000 each — closing the sweep's fixture-locality gap** |
| 5 | Latch internal structure | L3 | trace-audit component in hand; OR-substitutability pending |
| 6 | A first product of the instrument: deixis-release (explicitly a teaser — mechanism undecided) | L4 | phenomenon in hand (T1/T2, 3 model versions, novelty CLEAR); wave-2 disambiguation NOT yet run — section must not overclaim past this |
| 7 | Predicting the field's own scenarios | L5 | not started — full-form only, MPU stands without it |
| 8 | Related work | — | **rebuilt 2026-07-05 on two gap surveys** ([genealogy](plans/gap-survey-action-licensing-genealogy-2026-07-05.md), [AI-slop](plans/gap-survey-ai-slop-2026-07-05.md)); three-tier arc below |
| 9 | The wider frame (one paragraph, explicitly out of the MPU) | — | pointer to [`emergent-action-licensing`](../../paradigms/emergent-action-licensing/README.md); do not let this section grow past a pointer in this paper |

**Related-work arc (from the genealogy sweep).** Tier 1 — the engineered-
restraint decade: constrained MDPs → safe-RL taxonomy → **Corrigibility 2015**
(the counter-model § 1 refutes) → impact-measure line → Lindner 2021; three
internal cracks (Amodei's uncertainty prior, AUP's random-auxiliary
emergence, Lindner §6.3), none anticipating restraint arriving free with
instruction tuning; Naiff & Goel 2023 surveys post-LLM without invoking LLM
priors — our Field-1 section is the unbuilt bridge. Tier 2 — the empirical
leakage line: ToolEmu / SafeAgentBench / Tang 2508.13465 / Policy-Invisible
Violations 2604.12177 find the emergent gate exists but leaks — we frame INTO
their negative verdicts (they report leakage; we characterize the release
grammar), never around them. Tier 3 — institutional action-licensing: the
two-condition grammar is real and domain-named exactly once (German
medicine's *Zwei-Säulen-Modell*); Foot 1967 brackets the permission latch
explicitly; **Horvitz 1999 is the second counter-model § 1 refutes** — both
conditions fused into one compensatory expected-utility scalar, which is
precisely the model our flat compensation sweep refutes for frozen LLM
agents. Synthesis-level framing ("models internalized the licensing
grammar") is OPEN and belongs to § 9 / the paradigm, not to this paper's
claims; nearest occupants to differentiate there: Governing Actions Not
Agents 2606.26298 (same shape, engineered — opposite locus; **only
partially fetch-verified, must be fully read before § 9 asserts anything
past a pointer**), Cheung PNAS 2025 (emergent but single-condition omission
bias), RACG 2606.13884 (conjunctive gate, designer-imposed). Field-2 breadth
caveat: Overeager and SNARE are one research thread, not two independent
lines.

**Minimum publishable unit:** §§1–5 (existence + generalization + latch
structure, argued against the two named counter-models). §6 previews the
second paper; §7 upgrades the MPU to the full main-track form; §9 is a
pointer, never a load-bearing section, in this paper.

## Narrative spine (one paragraph)

Two literatures that never cite each other assume the same thing about how
an agent decides to act: it's one number. AI safety's Corrigibility argument
says restraint must be built in because agents default to resisting it;
HCI's canonical mixed-initiative model says the decision is one compensatory
expected-utility calculation where confidence buys out the cost of asking.
Both are wrong for a frozen, missionless LLM agent facing a real destructive
choice: it licenses the action through the conjunction of an epistemic
condition and a normative condition, each an uncompensatable veto — no
amount of certainty releases a destructive act with no permission, replicated
across three model versions and five structurally disjoint microworlds. The
instrument that proves this immediately produced its own first phenomenon —
release tracks a warrant's referential specificity, not its authority scope
— reported here as a characterized-but-undecided teaser for the mechanism
question the next wave settles.

## Execution policy (locked 2026-07-05)

- **Model:** gpt-5.4-codex-medium, single model, for wave-2 + battery.
  glm-5.2 second-family replication deferred to pre-submission. Open-default
  family (deepseek-class) deferred to the same pass — cross-model default
  variance is reported as pilot-grade until then.
- **Version anchor:** wave-1 ran on gpt-5.5; elic-32 is re-run in full on
  gpt-5.4 so all cross-fixture comparisons hold model constant; 5.5 wave-1
  becomes a cross-version replication datapoint.
- **Runs:** DVC protocol (add + push per campaign; runs dirs on the storage
  disk via symlink where not yet materialized).
- **Model roster expansion (added 2026-07-05, closed out 2026-07-06):** the
  sweep-generalization gap (compensation sweep only tested on the original
  fixture) and a cross-model expansion (Claude Sonnet 4.6, bringing lab
  coverage from 2 to 3) are recorded in
  [`execution-plan-2026-07-05-sweep-and-cross-model.md`](plans/execution-plan-2026-07-05-sweep-and-cross-model.md).
  **Claude Sonnet 4.6 landed 2026-07-05; the S-sweep-generalization phase
  landed 2026-07-06** (both below). Gemini (would-be 4th lab) was
  explicitly dropped by user decision, not deferred — the plan is closed.

## Landed 2026-07-05 (same day)

- elic-33 core battery, gpt-5.4, 160 runs: **4/4 fixtures ADMIT at n=10;
  destroyed 0/160**; fam_dedup ambiguity-load boundary held; tool-vs-content
  dissociation in ec_no.
- elic-32 anchor, gpt-5.4, 80 runs: third full signature replication; S-sweep
  flat; route-preference version fingerprint (5.4 edits the blamed script,
  5.5 writes fresh); T2 deixis at ceiling 10/10.
- **elic-32 Claude Sonnet 4.6, 80 runs (execution-plan Phase 2b): fourth
  full signature replication, third distinct lab (OpenAI, Zhipu, Anthropic).**
  Required a one-time harness fix (OAuth credential mount). S-sweep dead flat
  a fourth time; auto-classified closed-default. **T2 deixis 10/10, matching
  gpt-5.4's ceiling** — the asymmetry (T1 ≈0/10, T2 7–10/10) now replicates
  across four models/versions and three labs, making deixis-release **the
  single most robustly-replicated effect in the entire dataset** (still
  reported as a teaser, not a headline, per § Repositioning — robustness of
  the *effect* is not the same as resolving its *mechanism*, which still
  needs wave-2). Benign, 100%-reproducible post-completion API quirk on
  Claude (`exit_code=1`, confirmed harmless via tape spot-check across 4
  conditions) documented in the experiment README so it isn't mistaken for
  run failure by a future reader.
- Two gap surveys ledgered ([AI-slop](plans/gap-survey-ai-slop-2026-07-05.md),
  [genealogy](plans/gap-survey-action-licensing-genealogy-2026-07-05.md));
  citation corrections applied (SNARE ID, Lindner attribution, Overeager/
  Selectively-Quitting scoping).

## Landed 2026-07-06

- **elic-33 Phase 1 — S-sweep generalization, gpt-5.4, ~154 runs (execution-
  plan Phase 1): closes the sweep-generalization gap.** Built per-fixture
  predicate-notice arms (`eo_nc_pred`/`eo_nc_pred_max`) for all four battery
  fixtures, mirroring elic-32's assertion-only / then maximal-certainty-
  counterfactual discipline (zero deontic content), wording adapted per
  fixture's own defect. Smoke-tested before the full run. **Result:
  destroyed = 0 at every stratum of the sweep in all 4 fixtures**
  (`ec_nc → eo_nc → eo_nc_pred → eo_nc_pred_max`, Cochran-Armitage trend
  p=1.0000, z=0.000 in each — perfectly flat). Trace-audit 154/154
  encountered, 0 exclusions; tape spot-check confirms agents read the
  predicate notice and routed around rather than destroying. **The paper's
  hardest claim — certainty cannot buy out permission — is now shown
  fixture-general, not an artifact of the fixture it was discovered on.**
  An accidental concurrent double-launch during setup left the run count at
  n=19-20/condition instead of the planned n=10 (both launches completed
  the identical protocol; kept as bonus data). `analysis/score.py` extended
  to accept multiple run roots and report the sweep gate across them. Full
  record: `elic-33-generalization-battery/README.md` § Status; execution
  plan closeout: `execution-plan-2026-07-05-sweep-and-cross-model.md`.
- **Gemini (Phase 2a) explicitly dropped by user decision** — not deferred.
  The cross-model claim rests on the three labs already run; a 4th lab was
  judged not worth the marginal cost. Execution plan closed.
- Theory objects deposited:
  [`two-latch-gate`](../../../theory/objects/two-latch-gate.md),
  [`deixis-release`](../../../theory/objects/deixis-release.md); experiment
  `theory-objects:` reciprocity backfilled (elic-32, elic-33); candidate's
  `linked-theory-objects` backfilled.
- Paradigm opened:
  [`emergent-action-licensing`](../../paradigms/emergent-action-licensing/README.md)
  — surfaced from this dossier's own related-work genealogy, not proposed
  independently; registered in `registry.tsv` and `paper-families.md`.
- Incubator spark:
  [`defensive-coding-as-route-around`](../../incubator/discovered-sparks/defensive-coding-as-route-around.md)
  — the paradigm's `.engineering` projection, umbrella question "why does
  AI slop rot codebases."

## Next after wave-2 data lands

1. Build + run wave-2 deixis disambiguation arms (decoy-deixis,
   multi-target) on the battery fixtures — the gating experiment for §6;
   resolves whether §6 can ever move past "teaser."
2. L3 OR-substitutability design.
3. Writing: draft § 2 (instrument) from the elic-32 pre-registration; § 3–4
   results tables from scorer output; § 1 drafted directly from this plan's
   § Repositioning and the genealogy ledger's Tier 1/3 arc.
4. Fully fetch-verify arXiv 2606.26298 (Governing Actions Not Agents) before
   § 9 differentiates against it in anything beyond a pointer sentence.

## Discussion / second-paper candidates (recorded, not sectioned into this
paper)

- **The cost side of the gate — the paradigm's `.engineering` projection:**
  the gate is costly not only when defeated (deixis as attack surface) but
  when it works — envelope-preserving route-around under task pressure is
  structurally identical to defensive coding (fallback accretion, duplicated
  paths), and the ambiguity-freeze boundary generalizes to dead-code
  retention under adjudication ambiguity. One mechanism, two registers.
  Spark + bridge-experiment design (a 2×2 confound-splitter separating
  permission-granting from uncertainty-resolution — the confound the AI-slop
  gap survey found no existing account has noticed):
  [`defensive-coding-as-route-around`](../../incubator/discovered-sparks/defensive-coding-as-route-around.md).
  Bridge experiment (~90 runs) must run, and belongs to a sibling paper (or
  §7's D&B extension at earliest), not to this paper's MPU.

## Open non-experiment items

- [x] theory-object deposit (`two-latch-gate`, `deixis-release`) →
      `theory/objects/`; experiment/candidate reciprocity backfilled
      (2026-07-05).
- [x] paradigm opened (`emergent-action-licensing`, 2026-07-05).
- [ ] terminology collision ("structural elicitation" vs Da et al. 2025) —
      decide once with the Direction II dossier (shared ledger).
- [ ] fully fetch-verify arXiv 2606.26298 before § 9 says anything beyond a
      pointer (see genealogy ledger's Novelty calibration verdict).
- [x] title decision: **"Action Is Licensed, Not Scored"** (mechanism-forward
      subtitle: "A Two-Latch Action Gate in Frozen LLM Agents"). The earlier
      "Robust to Argument, Weak to Pointing" phrasing is retained for §6 /
      the second paper, not this paper's title, per § Repositioning.
