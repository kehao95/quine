# DR Query Set: Future Addressability And Computational Individuality

This query set starts the next `dr` research round for the theory paper:

> Computational individuality is indexed by future addressability.

The set is intentionally merged into four broad queries rather than many narrow
ones.

## Run IDs

| Query | Status | `dr` interaction id |
|-------|--------|---------------------|
| Query 1: Theory Of Computational Individuality | completed and recovered | `v1_ChdGWl92YWZURk5mblRqTWNQekxHWHVBWRIXRlpfdmFmVEZOZm5Uak1jUHpMR1h1QVk` |
| Query 2: Systems Carriers And Addressability | completed and recovered | `v1_ChdGWl92YWYtZ080MmwxTWtQNVlfZm1BcxIXRlpfdmFmLWdPNDJsMU1rUDVZX2ZtQXM` |
| Query 3: Formal Vocabulary And Models | completed and recovered | `v1_ChdGcF92YWI2YUdNaVk5TW9QbU5pdGdBaxIXRnBfdmFiNmFHTWlZOU1vUG1OaXRnQWs` |
| Query 4: Experimental And Evaluation Program | completed and recovered | `v1_ChdGWl92YWJDX090U24xTWtQamJlTHFRdxIXRlpfdmFiQ19PdFNuMU1rUGpiZUxxUXc` |

Recovered reports live in [`../references/README.md`](../references/README.md).

## Query 1: Theory Of Computational Individuality

Research computational individuality as a continuity criterion. We are
developing the thesis that computational individuality is indexed by future
addressability: a successor is continuous only relative to a future ecology in
which future obligations can still reach it through the relevant carrier and be
acted on.

Please find theories and canonical sources that support or challenge this
claim across philosophy of computation, artificial life, philosophy of
identity, organizational continuity, genidentity, functional continuity,
process ontology, and distributed computational systems.

Deliver:

- a conceptual map of identity criteria: artifact, source, memory, process,
  endpoint, behavior, obligation, carrier, organization
- strongest arguments for future addressability as an operational continuity
  criterion
- strongest objections and limits
- manuscript-safe language for presenting the thesis as criterion/theory rather
  than universal law
- 10-20 primary or canonical sources with stable bibliographic details

## Query 2: Systems Carriers And Addressability

Research OS and distributed-systems theories of continuity, identity, and
addressability. Focus on POSIX `exec`, inherited file descriptors, capability
systems, actor model/mailboxes, process migration, checkpoint/restart, service
failover, endpoint identity, distributed object references, and naming systems.

Our target concept is carrier topology: source carrier, execution carrier,
future-task carrier, state carrier, tool carrier, endpoint/mailbox carrier, and
authority/capability carrier.

Deliver:

- comparison table: system, boundary event, future carrier, continuity object,
  false-positive identity handle
- which systems preserve state versus reconstruct organization
- how addressability, authority, and mailbox/endpoint continuity are treated
  in canonical sources
- how to use these literatures without overclaiming life or agency
- 10-20 high-value sources with stable URLs/DOIs where possible

## Query 3: Formal Vocabulary And Models

Research formal models that can make post-handoff addressability precise.
Focus on process calculi, labeled transition systems, trace equivalence,
bisimulation, coalgebraic behavioral equivalence, interface automata,
session types, temporal logic/liveness, actor semantics, and capability
semantics.

We want to formalize:

```text
successor S is continuous under ecology E iff,
after boundary B,
S remains addressable through carrier C_E
and can perform E's required future-obligation trace.
```

For Quine this trace is:

```text
idle -> wake(task) -> tool/action -> report -> idle
```

Deliver:

- which formal traditions best fit carrier-indexed continuity
- minimal notation suitable for a theory paper
- definitions of future-task ecology, carrier, post-handoff addressability,
  future-obligation trace, and carrier-indexed self-object
- warnings about where the formalization would be misleading
- primary/canonical sources

## Query 4: Experimental And Evaluation Program

Research how to turn future addressability into an empirical program and
evaluation criterion. Use Quine as motivating case, but generalize to other
computational media.

Current candidate experiments:

- successor-mission ablation
- multi-wake closure
- carrier perturbation
- minimal non-Quine `ctl/wake` adapter
- source degradation / repair
- multi-generation recursion
- actor-mailbox and service-endpoint analogs
- checkpoint/restart contrast

Deliver:

- ranked experiment roadmap by theoretical payoff and feasibility
- pass/fail criteria for future addressability
- failure-mode taxonomy
- what evidence would support carrier-indexed individuality beyond Quine
- what evidence would falsify or weaken the thesis
- safety and containment boundaries
