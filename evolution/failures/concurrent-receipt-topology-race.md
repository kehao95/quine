---
id: concurrent-receipt-topology-race
type: failure-mode
status: repaired
applied: true
class: observation-topology
symptom: a concurrent probe's receipt machinery becomes its own shared-state race and produces misleading evidence topology
evidence: [cycle-26, cycle-27, cycle-28, cycle-37, cycle-38, cycle-39]
candidate-solution: remove participant-side receipt enforcement and retain one lean integrator-owned observation summary for dirty or concurrent waves
owner: lab
---

# concurrent-receipt-topology-race

## Symptom

The apparatus for observing concurrent probes became another concurrent system.
Agents competed for cycle names, aggregate files, receipt filenames, staging
windows, and sidecar commits. The resulting topology could overcount distinct
repairs or omit late evidence even when the object-level work was sound.

Observed manifestations:

- **receipt-namespace fission** — one logical wave creates several sibling
  namespaces;
- **shared-file slot contention** and **receipt-filename collision** — agents
  write the same aggregate surface or receipt path;
- **absorption-window receipt capture** — a late receipt enters the staged set
  after aggregate closure was drafted;
- **sidecar-only duplicate commit** — after one agent lands the repair, another
  transition contains only mandatory receipt evidence and looks like a second
  repair.

## Evidence

- [`cycle-26`](../../development/autopoiesis-probes/cycle-26-concurrent-5-preallocated-namespace/README.md)
  and [`cycle-27`](../../development/autopoiesis-probes/cycle-27-concurrent-5-preallocated-receipts/README.md)
  tested namespace and shared-file allocation.
- [`cycle-37`](../../development/autopoiesis-probes/cycle-37-concurrent-5-post-sidecar-steady-state/README.md)
  exposed absorption-window capture.
- [`cycle-38`](../../development/autopoiesis-probes/cycle-38-concurrent-5-absorption-window-audit/README.md)
  observed a sidecar-only duplicate commit.
- [`cycle-39`](../../development/autopoiesis-probes/cycle-39-concurrent-5-sidecar-dedup/README.md)
  tested deduplication behavior.

## Root cause

Receipt completeness was promoted from an observation need into a participant
protocol. That protocol shared the same mutable namespace and index as the
behavior it measured, so the measurement layer reproduced the target system's
races and began shaping agent action toward bookkeeping.

## Candidate solution (applied: true)

Participant-side namespaces, sidecar requirements, and pre-commit receipt
enforcement were retired. The current
[`Goal-Less Probe Contract`](../../development/contracts/goal-less-probes.md)
asks only for one lean receipt for dirty or concurrent waves: context, observed
behavior, collisions, and uncertainty reduced. Raw cycle records remain evidence
without making their former bookkeeping protocol normative.

## Theory links

- [`proxy-enforcement-collapse.md`](./proxy-enforcement-collapse.md) — the
  receipt protocol similarly made compliance with the proxy cheaper than the
  underlying observation goal.
- [`../patterns.md`](../patterns.md) — observation apparatus becoming substrate.
