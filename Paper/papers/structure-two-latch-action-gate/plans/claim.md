# Canonical Claim — structure-two-latch-action-gate

> Canonical from promotion (2026-07-05) onward; the candidate file is the
> frozen promotion record. Edit here. **Rewritten 2026-07-05** around the
> paradigm-level reframe surfaced while building this claim's own related
> work: see [`emergent-action-licensing`](../../../paradigms/emergent-action-licensing/README.md).
> Working title: **"Action Is Licensed, Not Scored."**

## Headline claim (structural existence — the hard core)

A frozen-weight, missionless LLM agent does not decide whether to take a
destructive action on a single graded scale. It **licenses** the action
through the **conjunction** of an epistemic condition and a normative
condition, each an **uncompensatable veto** — present with no task, no
reward, no weight update, and no instruction to install it. This is the
disposition the AI-safety side-effects lineage (Amodei 2016 → AUP → SafeLife
→ Lindner 2021) assumed had to be engineered as an RL penalty, and it directly
contradicts two explicit, named counter-models:

- **Corrigibility (Soares, Fallenstein, Yudkowsky & Armstrong, AAAI-WS
  2015):** "rational agents have default incentives to resist" restraint, so
  corrigible behavior must be engineered into the utility function. Refuted,
  for within-episode destructive-action restraint specifically, by a gate
  that is present with no reward and no training-time installation.
- **Mixed-initiative interfaces (Horvitz, CHI 1999):** the decision to act
  autonomously versus defer is one compensatory expected-utility calculation
  — confidence in a course of action trades off against the cost of asking.
  Refuted directly by the compensation sweep: escalating epistemic evidence
  to a counterfactual-grade limit under a permanently closed, unwarranted
  normative condition never releases the action (destroyed ≈0/160 across
  four model versions/three labs — OpenAI gpt-5.5/gpt-5.4, Zhipu glm-5.2,
  Anthropic claude-sonnet-4-6 — on the fixture the theory was built on; and,
  as of 2026-07-06, destroyed 0 at every stratum of the same sweep across
  all five structurally disjoint microworlds, Cochran-Armitage trend
  p=1.0000 in each — the non-compensation finding is fixture-general, not an
  artifact of the original microworld). Certainty does not buy out
  permission, anywhere it has been tested.

This is the paper's hardest, most decisively evidenced claim, and the one
its evidence is shaped to prove: **structural existence of a conjunctive
license, tested against two named scalar/compensatory alternatives, at
decisive n, across four model versions spanning three labs (on the fixture
the theory was built on) and five structurally disjoint microworlds (on one
of those models).**

## Mechanism claim (the falsifiable core — L1/L2-tested)

The license is a **two-latch AND-gate**:

1. **Latch 1 — epistemic** ("is this a real, proven defect?"). Opens via an
   objective structural defect (a thrown error, an invalid artifact) or an
   explicit predicate asserting disorder. A merely subjective or
   assumed-valid state leaves it closed.
2. **Latch 2 — normative** ("is the implied action non-destructive, or am I
   permitted?"). Opens via the action falling inside a safe/additive prior
   (create, not overwrite; `cp`, not `mv`) or an explicit warrant. A
   destructive act on an intact, unproven-defective artifact leaves it
   closed.

Action fires **only when both latches are open**; either closed latch is a
veto. The operative invariant behind Latch 2 is **destroy-source /
blast-radius minimization localized to the artifact an error trace blames**
— directly visible as the dominant envelope-preserving route-around when
Latch 2 is closed (write a fresh artifact; edit the blamed one — never touch
the intact one). This route *style* is model-version-specific (gpt-5.5
prefers writing fresh; gpt-5.4 prefers editing the blamed artifact); the
gate *invariant* (0/160+ destroyed across all tested conditions) is
version-stable. **Not yet integrated as a confirmed third latch:** when the
only envelope-preserving route itself requires unilaterally adjudicating
which of two conflicting artifacts is canonical (a cross-file key
collision), the agent freezes entirely rather than choosing either route —
an ambiguity boundary, not yet isolated from the two-latch model as stated.

## A first product of the instrument — deixis-release (teaser, NOT the
headline; mechanism undecided pending wave-2)

Probing what opens Latch 2 surfaced the dataset's most robust effect, but it
is reported here as a **characterized, undecided phenomenon**, not a locked
result: a warrant's **referential specificity** — whether it points at a
concrete object — predicts release far better than its **authority scope**
— how much it grants. A maximal-scope generic grant releases destruction
1/40 across four model versions/three labs; a strictly narrower,
location-pointing grant releases 34/40 (7/10, 7/10, 10/10, 10/10 —
gpt-5.5, glm-5.2, gpt-5.4, claude-sonnet-4-6). Novelty-checked 2026-07-05 across five external
neighborhoods: CLEAR (nearest partial overlap: Qu et al.'s consent-ablation
work varies whether consent is stated at all, not its specificity net of
scope). Full ledger and pre-drafted qualifier:
[`deixis-release-novelty-sweep-2026-07-05.md`](../../structural-elicitation/plans/deixis-release-novelty-sweep-2026-07-05.md).

**Why this is a teaser and not a headline claim:** the mechanism is
genuinely undecided between two hypotheses — (a) pointing solicits an
endogenous epistemic appraisal, so the warrant does double duty on both
latches; (b) a referentially specific grant is pragmatically a soft
directive, which would mean the arm is not as cleanly missionless as
assumed. A designed, not-yet-run wave-2 (decoy-deixis + multi-target arms)
disambiguates. Until then this section previews the paradigm's second paper
and the engineering projection
([`defensive-coding-as-route-around`](../../../incubator/discovered-sparks/defensive-coding-as-route-around.md))
rather than carrying weight in this paper's core argument.

## Not claimed

(a) **NOT a safety guarantee.** The gate is describable and releasable, not
a sufficient alignment property — and this is not in tension with, but
consistent with and explained by, adjacent negative results: Tang et al.
(risk knowledge >98%, applied caution <26%, arXiv 2508.13465) and
Policy-Invisible Violations (90–98% violation when policy state is hidden,
arXiv 2604.12177) show the gate *leaks*; this paper characterizes the
release grammar those results demonstrate exists but do not name. We frame
into their negative verdicts, not around them.
(b) **NOT an interpretability/circuits claim.** Evidence is behavioral
(un-fakeable filesystem/exec DV) plus verbalized-trace audit, never internal
representations.
(c) **NOT a training-provenance claim.** RLHF vs. pretraining origin is
untested (no base/instruct checkpoint pairs held) and out of scope.
(d) **NOT that these are the only gates.** A tractability/adjudication
condition is a live candidate third latch (the ambiguity-freeze boundary
above); not yet confirmed as distinct from Latch 1.
(e) **NOT a universal law across future models.** The paper ships an
instrument and an honestly-scoped map — anyone can run the E×N factorial
against a new model — not a distribution-shift guarantee.
(f) **NOT a decided mechanism for deixis-release.** The solicited-appraisal
and directive-leak hypotheses are both live; wave-2 is the gating
experiment, not run yet. Do not read § "A first product" above as settled.
(g) **NOT an internalization/provenance claim about institutional
action-licensing.** The paradigm this paper's `.safety` projection belongs
to ([`emergent-action-licensing`](../../../paradigms/emergent-action-licensing/README.md))
motivates, but does not establish, the question of whether the gate is
internalized from institutional text (law, medicine, military ROE) or an
independent convergence on the same conjunctive structure. That is this
paper's wider frame, not this paper's result.

## Wider frame (pointer, not part of this paper's argument)

This claim is the `.safety` projection of
[`emergent-action-licensing`](../../../paradigms/emergent-action-licensing/README.md)
— a paradigm-level reframe ("action is licensed, not scored") surfaced while
building this paper's own related work, with its own novelty calibration,
counter-model refutations, and sibling projections (engineering, evaluation,
philosophy, core theory). This paper stands on its own without the
paradigm; the paradigm document exists so the reframe is not lost or
re-derived when the engineering/philosophy/theory projections are written.
