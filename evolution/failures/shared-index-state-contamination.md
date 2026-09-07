---
id: shared-index-state-contamination
type: failure-mode
status: repaired
applied: true
class: shared-substrate
symptom: concurrent agents share a moving worktree and Git index, so staged ownership, artifact provenance, and reported commit identity become unstable
evidence: [cycle-21, cycle-24, cycle-27, cycle-32, cycle-40, cycle-41]
candidate-solution: classify shared moving runs as substrate-stress evidence, record start and end state, and reserve autonomous-preference claims for clean isolated probes
owner: lab
---

# shared-index-state-contamination

## Symptom

Concurrent agents in one dirty worktree do not observe independent repository
states. They share the index, see `HEAD` advance during orientation, and encounter
files or staged deletions produced by actors whose transition they do not own.
Attribution becomes unstable even when every individual agent behaves correctly.

Observed manifestations:

- **index-capture contamination** — one agent's commit captures another agent's
  staged or uncommitted work;
- **commit-identity drift** — a reported commit ID ceases to identify the final
  transition after concurrent advances, amendments, or absorption;
- **moving-artifact tombstoning** — a newly appearing artifact is mistaken for
  disposable local residue;
- **staged-surface tombstoning** — an agent interprets another actor's staged
  deletion as authority to restore or complete it;
- **active-index race reroute** — `HEAD` changes during cold-start orientation
  and exposes a route that was absent from the agent's starting state.

## Evidence

- [`cycle-21-concurrent-5-dirty-moving-measured`](../../development/autopoiesis-probes/cycle-21-concurrent-5-dirty-moving-measured/README.md)
  named moving-artifact tombstoning.
- [`cycle-24-single-dirty-moving-control-plane`](../../development/autopoiesis-probes/cycle-24-single-dirty-moving-control-plane/README.md)
  recorded commit-identity drift.
- [`cycle-27-concurrent-5-preallocated-receipts`](../../development/autopoiesis-probes/cycle-27-concurrent-5-preallocated-receipts/README.md)
  observed transient staged-sibling index capture.
- [`cycle-32-concurrent-5-committed-namespace-repeat`](../../development/autopoiesis-probes/cycle-32-concurrent-5-committed-namespace-repeat/README.md)
  recorded staged-surface tombstoning.
- [`cycle-40`](../../development/autopoiesis-probes/cycle-40-concurrent-5-receiptless-diagnostic-guard/README.md)
  exposed active-index race reroute; [`cycle-41`](../../development/autopoiesis-probes/cycle-41-concurrent-5-active-index-race-control/README.md)
  controlled it by fixing the starting state before the wave.

## Root cause

Git's worktree and index are shared mutable substrate, while agent reports are
local observations sampled at different times. Treating those observations as
if each agent had a private snapshot creates false ownership and causal claims.

## Candidate solution (applied: true)

The repair is epistemic, not a concurrency bureaucracy. The
[`Goal-Less Probe Contract`](../../development/contracts/goal-less-probes.md)
now distinguishes clean isolated preference probes from dirty moving stress
probes, records start/end `HEAD`, and forbids clean-preference claims when mode
provenance is missing. Newly appearing or staged state remains foreign until
ownership is established.

## Theory links

- [`validation-attraction.md`](./validation-attraction.md) — another substrate
  property that redirects goal-less behavior.
- [`../patterns.md`](../patterns.md) — shared-substrate attribution failure.
