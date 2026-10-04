---
title: Language-Mediated State Propagation in Persistent LLM Agent Systems
configuration_id: aaai-state-propagation
status: planning
paused: true
kind: AI/ML-core problem-formulation paper (venue binding lapsed)
venue: aaai-27
format-envelope: 7page-archival-two-column
role: venue-projection
maturity-target: existence
state-mode: planning-skeleton
track: main
audience: [ai-ml-core, multi-agent-systems, llm-agents]
claim_bundle: substrate-coincidence (AAAI projection)
source_owner: true
plan_home: README.md
family: open-ended-ecology
parent_theory: substrate-coincidence
sibling_directions: [alife-theory, exec-unfakeable]
---

# Language-Mediated State Propagation in Persistent LLM Agent Systems

> **PAUSED — venue binding lapsed 2026-07-26.** The AAAI-27 cycle was missed
> (abstract deadline `2026-07-21` passed with no registration; see
> [`venues/objects/aaai-27.md`](../../../venues/objects/aaai-27.md) § Cycle
> outcome). This configuration is **withdrawn from AAAI-27 and not re-bound to
> another venue yet**. Nothing below is retracted: the spine, the three
> load-bearing moves, the formal model, the metrics, and the evidence stance in
> [`DOSSIER.md`](DOSSIER.md) all stand as written and stay the rewrite basis.
> What lapsed is a venue and a page budget, not a claim.
>
> The `venue:` and `format-envelope:` fields below are kept as the historical
> binding so the AAAI reframe stays legible; read them as "what this was shaped
> for," not as a live target. Re-binding is a decision, not a default — take it
> only when someone is actually going to draft, and record the new venue here
> and in [`../DIRECTIONS.md`](../DIRECTIONS.md) at the same time.

## What this paper is

An AI-systems paper arguing that **persistent LLM-agent systems violate a
common abstraction in agent design**: the assumption that memory, logs,
files, and peer outputs are passive environment or auxiliary retrieval
resources. In language-mediated runtimes over a shared persistent workspace,
reading a workspace artifact incorporates external language state into the
agent's effective context, and writing an artifact creates a latent
intervention on future agents. This makes the workspace a **behavioral
state-propagation channel**, not passive memory.

This is **not** an artificial-life paper. ALife theory
([substrate-coincidence](../alife-theory/manuscript-draft.md)) is retained as a
theoretical source; it does not appear in the title, abstract, or first
layer of contributions.

## Central claim (spine)

> Persistent LLM-agent systems propagate state through language artifacts
> via a path-dependent, potentially self-amplifying feedback loop: an
> agent's writes persist unfiltered, re-enter as context (token-faithful
> per read), and — because absorbed text simultaneously shapes the next
> read trajectory — feed forward into both the agent's own and its
> successors' behavior. The workspace is therefore a behavioral state
> channel, and its governance is a runtime-safety problem.

## Audience and venue

- **Venue**: AAAI-27 main technical track
  ([card](../../../venues/objects/aaai-27.md)). Deadline: abstract
  `2026-07-21`, full `2026-07-28` UTC-12.
- **Format**: 7-page two-column + references
  ([envelope](../../../venues/envelopes/7page-archival-two-column.md)).
- **Audience lens**: primary `ai-ml-core`, secondary multi-agent systems /
  LLM agents. Safety enters as implication, not as the main claim.
- **Track choice**: main track (not AI Alignment). Rationale in
  [DOSSIER.md § Track](DOSSIER.md).

## Quine is the instrument, not the subject

The paper demonstrates the regime via the Quine runtime — persistent
multi-agent LLM agents running as OS processes over a shared filesystem —
introduced once as the instrument. All phenomena are described at the
**class level** (persistent language-mediated agent systems), not as
Quine-runtime machinery. See [DOSSIER.md § Evidence stance](DOSSIER.md).

## Status

**Frozen rewrite basis, paused, unbound.** See [DOSSIER.md](DOSSIER.md). The
basis is complete and does not need revisiting; the missing step was always the
same one and it never happened — draft abstract + intro + contribution list +
formal-model skeleton from that basis. That remains the next step whenever this
line is picked back up, under whatever venue it is re-bound to.

Honest read on why it stalled: the dossier locked scope thoroughly enough to
feel like progress, and no one converted the lock into prose. A frozen basis is
a precondition for drafting, not a substitute for it.

## Boundary against sibling directions

- Owns: read/write-lineage formal model, operational metrics,
  abstraction-failure diagnosis, class-level propagation phenomena.
- Does **not** own: three-separations-collapse theory (→ ALife direction),
  exec barrier / spec-gaming cure (→ exec-unfakeable direction).
- Full boundary table: [`../DIRECTIONS.md`](../DIRECTIONS.md).

## Source surfaces

- Frozen basis: [./DOSSIER.md](DOSSIER.md)
- Family index: [../DIRECTIONS.md](../DIRECTIONS.md)
- Theory object: [`../../../../theory/objects/substrate-coincidence.md`](../../../../theory/objects/substrate-coincidence.md)
- Experimental basis: [`../../../../../experiments/active/null-ecology/null-05-reflexive-self-embedding/`](../../../../experiments/active/null-ecology/null-05-reflexive-self-embedding)
