---
paper_id: material-channels
claim_bundle: material-channels
venue: aamas
format-envelope: 7page-archival-two-column
track: system
cluster: social-behavior
status: planning
state_mode: planning-skeleton
role: communication-surface-follow-on
family: social-protocol
theme: social-behavior
theme_home: ../../../theory/views/social.md
publication_class: archival-candidate
maturity_target: mechanism
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [coordination-closure, protocol-emergence, carrier]
gaps: [collective-continuity]
evidence: [2.1b-shell-mapreduce]
view-affinity: [social]
---

# Material Channels

## Current Question

Which communication and material surfaces should become named runtime objects instead of staying hidden behind privileged transport like `stdin` or broker-managed message buses?

## Thesis

Agent material and inter-agent communication should become named, inspectable filesystem channels rather than remain hidden behind privileged `stdin` or broker-managed buses.

## In Scope

- channel objects as public runtime resources
- the continuity between single-agent input and multi-agent communication
- inspectable, composable, permissioned channel semantics

## Out Of Scope

- world semantics as the primary frame
- time/rollback semantics as the primary frame
- generic distributed-systems messaging abstractions
- empirical claims about whether agents already invent callback / routing /
  delegation conventions over the current minimal runtime; that belongs in
  [`../minimal-runtime-protocols/README.md`](../minimal-runtime-protocols/README.md)

## Current State

- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- upstream dependency: [`../process-as-fs/README.md`](../process-as-fs/README.md)
- empirical neighbor: [`../minimal-runtime-protocols/README.md`](../minimal-runtime-protocols/README.md)
- related neighbor: [`../subjective-reality/README.md`](../subjective-reality/README.md)

## Experiment Feeds

- `experiments/active/early-existence/exst-04-shell-mapreduce/README.md` - future-home - exploratory-only - points toward later work on explicit communication and coordination surfaces beyond privileged process calls

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - no stable evidence backbone is currently declared for this paper

## Strongest Evidence

- Quine's existing process contract around `stdin` / `stdout`
- design pressure from material streams and mesh-collaboration evolution notes
- the need to expose communication as runtime topology instead of framework policy

## Blocking Gaps

- the runtime does not yet expose these channels as stable public objects
- the paper needs implementation slices, not only architecture sketches
- communication semantics must stay distinct from broader world-construction claims

## Next Useful Work

- treat this as a downstream paper after `process-as-fs`
- collect concrete channel-surface experiments before tightening prose
- keep mesh topology and channel semantics in one bounded argument
- let current evidence about callback / routing / protocol emergence live in
  [`../minimal-runtime-protocols/README.md`](../minimal-runtime-protocols/README.md)
  so this paper can stay architectural

## Unsafe Overclaims

- do not treat architectural pressure for material channels as proof that the runtime already exposes them as stable first-class objects
- do not let this paper absorb empirical protocol-emergence claims that belong in `minimal-runtime-protocols`
- do not broaden filesystem-native channel semantics into generic distributed-systems messaging rhetoric
