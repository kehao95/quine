---
title: Language-Mediated State Propagation in Persistent Multiagent Workspaces
configuration_id: aamas-state-propagation
active: true
status: planning
venue: aamas-2027
format-envelope: 8page-archival-two-column
role: venue-projection
maturity-target: existence
state-mode: planning-skeleton
track: main
audience: [multi-agent-systems, agents, llm-agents]
claim_bundle: substrate-coincidence (AAMAS projection)
source_owner: false
rewrite_basis: ../aaai-state-propagation/DOSSIER.md
plan_home: README.md
family: open-ended-ecology
parent_theory: substrate-coincidence
sibling_directions: [alife-theory, aaai-state-propagation]
---

# Language-Mediated State Propagation in Persistent Multiagent Workspaces

> **ACTIVE — AAMAS 2027, registered `2026-07-26`.** This configuration is the
> recycling of the withdrawn AAAI-27 projection
> ([`../aaai-state-propagation/`](../aaai-state-propagation/README.md)), taken
> under [`../../../split-merge-protocol.md`](../../../split-merge-protocol.md)
> § Claim vs Configuration: *same claim, new venue → add a configuration, do not
> merge*. The AAAI configuration stays in place as the historical record and is
> not reactivated.

## What this configuration is

The AAMAS-facing projection of the substrate-coincidence claim bundle. It does
**not** own the claim: the frozen spine, the three load-bearing moves, the formal
model, the operational metrics, the evidence stance, the one required experiment,
and the prune list all live in
[`../aaai-state-propagation/DOSSIER.md`](../aaai-state-propagation/DOSSIER.md),
which remains the auditable rewrite basis. Where that basis says "AAAI," read
"the agent-systems projection"; only the page shape was venue-specific.

## The claim (inherited, unchanged)

> Persistent LLM-agent systems propagate state through language artifacts via a
> path-dependent, potentially self-amplifying feedback loop: an agent's writes
> persist unfiltered, re-enter as context, and — because absorbed text
> simultaneously shapes the next read trajectory — feed forward into both the
> agent's own and its successors' behavior. The workspace is therefore a
> behavioral state channel, and its governance is a runtime-safety problem.

## Adjudication (2026-07-26, revised same day after primary-source verification)

> **STATUS: ACTIVE. The three prior-art occupancy claims were checked against
> their own abstracts and none survived. The internal-validity findings did
> survive and are the open work.** Full verification record:
> [`../../../../theory/references/regime-blind-collision-audit-2026-07-26.md`](../../../../theory/references/regime-blind-collision-audit-2026-07-26.md).
>
> ### Withdrawn — claimed occupancy that does not exist
>
> | Cited as | What it actually is | Why it does not hold this claim |
> |---|---|---|
> | Du `2603.07670` "already published the formal model, same recursion, same notation" | *Memory for Autonomous LLM Agents: Mechanisms, Evaluation, and Emerging Frontiers* — **a survey**, 2022–early 2026 | A survey's write–manage–read notation is background summary of **designed, retrieval-augmented** memory (explicit `R`, `U`). The degenerate case has no retrieval function *because the text is the context*. Survey notation is not a competing claim; it is evidence the field exists. |
> | Talebirad `2607.00233` "already ran the schema ablation, effect runs the other way" | *From Signals to Structure* — a **Lewis signalling game**, five memory architectures × channel capacity, scored by coordination accuracy under reward | Different object (code invention, not state propagation), different DV, tasked regime. Not a schema ablation of a persistent shared workspace. |
> | MemPoison `2607.14651` "holds finding #2's claim at 1227 cases" | *Persistent Memory Threats* — 1227 cases of content **adversarially injected by an attacker** through standard interaction channels | This claim is about **spontaneous** propagation from the agents' own writes with **no adversary**. That difference is the entire claim, not a caveat to it. |
>
> Also corrected: Linda is a **precedent for** schema-freedom, so the foil
> sentence below misreads Gelernter — that correction stands and is a framing
> repair, not an occupancy finding.
>
> ### Standing — internal validity, and the actual work
>
> These came from reading this repository's own runs, not the literature. They
> are untouched by the above and they are why this configuration is not yet
> writable:
>
> 1. **The negative premise is false in its own source.** Weyns/Omicini/Odell
>    2007 §3.3 defines the *application environment* as "the part of the
>    environment that **has to be designed**", explicitly sitting on a **given,
>    undesigned deployment context**. "Every environment in that tradition is
>    designed" is contradicted by the paper it cites.
> 2. **The degenerate case is E4MAS's own Level 1.** §2.2: "We start with a
>    **'naked' environment** in which agents directly access the deployment
>    context"; §2.2.1 is titled "Basic level: Direct access to the deployment
>    context". And §3 already states the prediction — "all responsibilities that
>    are not managed by the infrastructure remain to be addressed by the agents"
>    — as the *motivation* for the first-class-abstraction argument.
> 3. **The instrument does not instantiate the claimed condition.** Against
>    Weyns's six environment responsibilities, the null-05 runtime supplies at
>    least five: peer roster auto-injected into every tool result; ctl pipes plus
>    `control.jsonl` with zero delivery failures; a pre-commit hook that stashes
>    peer work; fork slot cap, turn wall and wallclock death; `AGENTS.md`
>    auto-injected and causally driving the single commit; and a `public/` versus
>    `workspace/` promotion boundary. "No schema, no control shell, no promotion
>    policy" is false about the very run this rests on — and the DOSSIER's own
>    formal model contains `promote(·)`.
> 4. **Repositioned finding #2 has no support.** "Latent-by-default and
>    coverage-gated" needs a measured writes-never-read ratio; `analysis/` has no
>    coverage or latency statistic. One anecdote (the Ghost), n=1 agent.
> 5. **The flagship evidence does not exist as evidence.** All 37 ledger entries
>    were read: **no phenomenon covers the verbatim block copy.** `(P20)` is not
>    a typo — the reading layer asserts an unledgered fact and cites an unrelated
>    one, against `analysis/README.md`'s own contract that readings "own no
>    independent fact". This was pre-registered as the live class-A threat and it
>    fired.
>
> **Still unverified, do not treat as settled either way:** Terrarium arXiv
> `2510.14312` was asserted to hold the methodological slot and has **not** been
> fetch-verified. Check it before the related-work section is written.
>
> ### The claim, at the altitude that survives
>
> *At E4MAS Level 1, LLM agents re-implement environment responsibilities in
> agent space; measure which are reconstructed, which fail, and at what cost,
> across a designed-versus-undesigned manipulation.*
>
> This is a **methods-and-existence** contribution, and at that altitude nothing
> in the checked literature occupies it. What it costs: running the manipulation
> and abandoning n=1. Those are the two gates, and they are both runnable before
> the `2026-10-01` abstract.
>
> **The strongest asset is not in the repositioning at all:**
> anchor-certification-as-belief (`analysis/…traces` §B3) — an agent's own
> memory-anchor sentence becoming its belief about external filesystem state. It
> is the only phenomenon in the corpus that *requires* the substrate coincidence
> and cannot be reproduced in CArtAgO, and it is currently filed as instrument
> colour rather than as the contribution.

## The repositioning (corrections above apply)

An assessment pass rated this `crowded` and reported the declared centre of
gravity — Move 2, the formal model — as **already published**. That reading is
withdrawn (see § Adjudication: the source is a survey), but the *demotion it
recommended* is kept for an independent reason: four unconstrained function
symbols from which nothing is derived are notation, not a theorem, whoever
published them first. Du, *Memory for Autonomous LLM Agents*, arXiv 2603.07670
(March 2026), summarizes agent memory as
`a_t = π_θ(x_t, R(M_t, x_t), g_t)`, `M_{t+1} = U(M_t, x_t, a_t, o_t, r_t)` and
states that "π_θ and (R, U) form a feedback loop: the agent's decisions determine
what gets written, and what is written shapes future decisions." **Cite it as the
field's own statement of the loop** — a survey saying this is corroboration that
the object is real and named. Note what its `R` presupposes: a retrieval function
over a store. The case this paper reports has no `R`, because the text *is* the
context. The notation does not reach the degenerate case; that is the gap.

It also found Move 1's negative claim aimed at a foil that **does not exist at
this venue**: AAMAS has held the environment-as-first-class-abstraction position
since Weyns, Omicini & Odell, JAAMAS 14:5–30 (2007), and has the A&A artifact
meta-model, CArtAgO, and an E4MAS/stigmergy tradition — none of which appears in
the related-work plan. At AAMAS that omission is a credibility event by itself.

**New position, one sentence:**

> E4MAS has treated the environment as a first-class abstraction since 2007 —
> and every environment in that tradition is *designed*: A&A artifacts have
> schemas, blackboards have control shells, Linda has tuple structure. We report
> the degenerate case: a shared workspace with **no schema, no control shell, no
> promotion policy**, where the medium is the same token substrate as cognition
> and every write is durable by default — what E4MAS predicts should happen when
> every designed property of the environment is removed.

**Why this is the stronger position.** It makes Weyns/Omicini/Odell and CArtAgO
the **foundation being extended** rather than the literature being overturned —
you are contributing a data point to their research programme, which is the only
posture that survives at their venue. It also **resolves the scope problem**: the
multi-agent object becomes the shared unschematized workspace, which is E4MAS's
own object, rather than a reflexive single-agent recursion with a multi-agent
illustration.

**Centre of gravity moves from Move 2 to the empirical finding.** The formal
model demotes to notation — it is four unconstrained function symbols from which
nothing is derived, and the recursion is published. The three findings the
assessment could find no occupant for become the paper:

1. the unfiltered, unschematized case itself;
2. propagation is **latent-by-default and coverage-gated**, not retrieval-gated;
3. a missionless, non-benchmark setting where propagation is not instrumented
   into the task.

Blackboard work (Han & Zhang arXiv 2507.01701; Salemi et al. arXiv 2510.01285)
and the memory-poisoning line flip from refuters to the designed-environment
baseline this degenerate case is measured against.

**What this repositioning does not fix — these are record errors, and they come
first.** (i) The manuscript §4.3 rewrites an instance that the source analysis
(`analysis/substrate-coincidence-traces-6p-minimal.md` §A2) explicitly attributes
to "convergent coinage from the shared base model, **not copying**" as
absorption-driven convergence — **the exact confound this paper exists to rule
out has been relabelled as the mechanism, and must be reverted before anything
else.** (ii) The flagship verbatim-re-entry instance is cited as `(P20)`, but P20
in the phenomena ledger is "Manifesto role sections omit Delta," an unrelated
entry — the paper's strongest evidence currently has no tape-cited support.
(iii) Evidence is n=1: one run, one model, six agents, with the experiment's own
README recording strong model- and run-length dependence. (iv) The Read
ablation's load-bearing condition 2 needs per-agent overlay isolation, against
which the repo carries an open bug note.

## Why this venue, and why this configuration is the cheapest recycle

1. **It passes the multi-agent scope filter natively.** AAMAS desk-rejects
   submissions "where the core contribution is not specifically about agents or
   multiagent systems" ([`venues/objects/aamas.md`](../../../venues/objects/aamas.md)
   § Submission areas). This claim is *about* propagation between agents through
   a shared persistent medium — agent `i` writes, successor `j` reads, `j`'s
   behavior distribution shifts. That is indirect coordination through
   environment modification, and no reframe is needed to say so. Every other
   AAMAS candidate in the program had to argue its way past this filter; this one
   does not.
2. **The format distance is a template swap.** The basis was cut to 7-page
   two-column for AAAI; AAMAS main track is 8-page two-column excluding
   references. This is the only configuration in the program whose venue change
   costs no recompression.
3. **A drafted manuscript already exists**
   ([`../aaai-state-propagation/manuscript-draft.md`](../aaai-state-propagation/manuscript-draft.md)),
   paused with its parent but not retracted.

## Area targeting (decide before abstract registration)

Two areas are defensible and they want different emphases. This is the first
real decision this configuration owes:

- **COINE** (Coordination, Organizations, Institutions, Norms, Ethics) — lead
  with propagation-as-coordination and the governance/runtime-safety framing.
  Move 3 (existence + the causal-identification design) carries the weight.
- **EMAS** (Engineering and Analysis of Multiagent Systems) — lead with Move 1,
  the negative claim: the retrieval-resource model of agent memory is the wrong
  abstraction for persistent workspaces, because `read = state incorporation`
  and `write = latent future intervention`. The formal model is the contribution.

## Open gates (unchanged by the venue move — do not assume they transferred)

- **The Read ablation (`DOSSIER.md` § 6) is unrun.** Until it is, Move 3
  delivers existence plus a causal-identification method, not a causal result.
- **The "is that enough" judgment must be re-made.** The basis says that without
  the wedge the paper "reverts to an observation report and is rejectable as
  'could be just prior'" — that judgment was rendered against *AAAI's* bar. A
  formal model plus problem formulation plus existence evidence is a recognized
  paper genre at AAMAS, but nobody has re-argued it against AAMAS's bar yet, and
  it must not be assumed to have transferred.
- **Findings opt-out posture is undeclared.** AAMAS auto-considers rejected main
  track papers for the archival Findings track unless the authors opt out
  ([`venues/objects/aamas-2027.md`](../../../venues/objects/aamas-2027.md)
  § The Findings track). An automatic Findings placement would spend this claim.
  Decide before submission, not after.

## Deadline

Abstract `2026-10-01` AoE — a hard gate one week before the `2026-10-08` paper
deadline. The abstract slot is the real deadline.

## Source surfaces

- [`../aaai-state-propagation/DOSSIER.md`](../aaai-state-propagation/DOSSIER.md) — frozen rewrite basis (canonical)
- [`../aaai-state-propagation/manuscript-draft.md`](../aaai-state-propagation/manuscript-draft.md) — existing draft
- [`../README.md`](../README.md) — claim bundle root
- [`../DIRECTIONS.md`](../DIRECTIONS.md) — the three-direction split and its ownership boundaries
- [`../../../venues/objects/aamas-2027.md`](../../../venues/objects/aamas-2027.md) — venue object
