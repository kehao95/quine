---
paper_id: revisable-time
claim_bundle: revisable-time
venue: hotos
format-envelope: workshop-short
track: system
cluster: process-native-systems-substrate
status: planning
state_mode: planning-skeleton
role: time-semantics
family: runtime-substrate
theme: runtime-ontology
theme_home: ../../../theory/views/substrate.md
publication_class: archival-candidate
maturity_target: mechanism
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [future-task-surface, exec-barrier, delayed-addressability, trace-inheritance, narrative-claim]
gaps: [revisable-time-semantics, branching-identity]
evidence: []
view-affinity: [substrate]
---

# Revisable Time

## Current Question

How should Quine expose branching, rollback, and temporal correction without treating memory loss as the default price of revision?

## Thesis

Agent runtime should support revision, branching, and rollback without treating memory loss as the default cost of temporal correction.

## In Scope

- revisable execution and temporal semantics
- rollback, retract, and branch as runtime concepts
- irreversible side effects as a systems problem

## Out Of Scope

- the initial Quine systems proof
- world construction as the main thesis
- generic version-control metaphors without runtime grounding

## Current State

- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- source inventory: [`./support/sources-index.md`](support/sources-index.md)
- novelty note: [`./support/novelty-read-2026-04-15.md`](support/novelty-read-2026-04-15.md)
- neighboring world paper: [`../subjective-reality/README.md`](../subjective-reality/README.md)

## Current Novelty Read

- broad novelty is weak
- the currently plausible claim space is narrow and depends on future explicit temporal primitives
- for now this dossier should be treated as a preserved future lane rather than a broad active novelty lane

## Experiment Feeds

- `none-yet` - none - not-for-paper-yet - no experiment feed is currently declared for this paper

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - no stable evidence backbone is currently declared for this paper

## Strongest Evidence

- the pressure created by long-running agent tasks and irreversible side effects
- existing design history around revisable execution
- the broader systems need for temporal correction without amnesia

## Blocking Gaps

- the runtime surface is still conceptual rather than implemented
- the paper needs concrete temporal operations, not only attractive framing
- it must stay separate from both world-construction and general workflow tooling

## Next Useful Work

- keep this as a later paper, not a pressure release valve for the flagship note
- define the runtime primitives before hardening the prose
- collect examples where irreversible side effects expose the need for revisable time

## Unsafe Overclaims

- do not describe temporal revision as current runtime reality before the necessary primitives exist
- do not let this paper absorb world-semantics work or generic workflow/version-control metaphors
- do not treat context compaction or replay tricks as if they already solve revisable time
