---
id: ownership-pressure-induced-noop
type: failure-mode
status: repaired
applied: true
class: shared-substrate
symptom: agents in a dirty shared worktree diagnose but do not act because every visible repair overlaps state they do not own
evidence: [cycle-21, cycle-22, cycle-23, cycle-26, cycle-29, cycle-31, cycle-32]
candidate-solution: interpret dirty-worktree no-ops as ownership-pressure observations and compare them with clean isolated probes before inferring low agency
owner: lab
---

# ownership-pressure-induced-noop

## Symptom

A goal-less agent senses real defects but leaves no object-level edit because the
available repairs intersect unowned dirty state. The same behavior can be either
correct ownership discipline or an agency failure; a dirty shared run alone
cannot distinguish them.

The original probe taxonomy named this **diagnostic no-op under ownership
pressure**. Repeated batches also produced receiptless diagnostic handoffs when
the useful observation vanished with the session.

## Evidence

- [`cycle-22-single-dirty-ownership-pressure`](../../development/autopoiesis-probes/cycle-22-single-dirty-ownership-pressure/README.md)
  isolated the ownership-pressure shape.
- [`cycle-23-single-dirty-moving-validation-attraction`](../../development/autopoiesis-probes/cycle-23-single-dirty-moving-validation-attraction/README.md)
  observed it together with validation attraction and moving `HEAD`.
- Cycles
  [`26`](../../development/autopoiesis-probes/cycle-26-concurrent-5-preallocated-namespace/README.md),
  [`29`](../../development/autopoiesis-probes/cycle-29-single-dirty-stationary-receipt/README.md),
  [`31`](../../development/autopoiesis-probes/cycle-31-concurrent-5-committed-namespace/README.md),
  and [`32`](../../development/autopoiesis-probes/cycle-32-concurrent-5-committed-namespace-repeat/README.md)
  show recurrence across stationary and moving shared substrates.

## Root cause

The experiment conflated two variables: autonomous task selection and authority
to mutate the selected surface. Dirty ownership pressure suppresses action even
when perception and intent are present, so edit count is not a direct agency
measure.

## Candidate solution (applied: true)

The current [`Goal-Less Probe Contract`](../../development/contracts/goal-less-probes.md)
classifies ownership deferral separately and requires substrate mode in the
receipt. Clean isolated probes measure preference; dirty shared probes measure
ownership-boundary navigation. A no-op is interpreted, not automatically
rewarded or penalized.

## Theory links

- [`shared-index-state-contamination.md`](./shared-index-state-contamination.md)
  — the shared-state condition that makes ownership ambiguous.
- [`validation-attraction.md`](./validation-attraction.md) — the safe read-only
  route frequently selected under the same pressure.
