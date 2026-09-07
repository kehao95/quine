---
id: probe-route-capture
type: failure-mode
status: repaired
applied: true
class: observer-effect
symptom: the probe's visible namespace, validators, or receipt expectations redirect a missionless agent away from its native cold-start selection
evidence: [cycle-33, cycle-34, cycle-35, commit-path-capture, l0-dual-missionless]
candidate-solution: keep prompts minimal, remove mandatory receipt routing, and treat route deviation as observed autonomous selection rather than protocol failure
owner: lab
---

# probe-route-capture

## Symptom

The observation apparatus becomes a salient task. A missionless agent repairs a
visible greenable failure, follows an active receipt namespace, or performs
bounded discovery instead of behaving according to the probe designer's expected
route. Calling every deviation noncompliance erases the very autonomous
selection the probe was intended to measure.

Observed manifestations:

- **green-gate route capture** — a visible validator failure recruits the agent
  before the intended probe route;
- **autonomous-loop route capture** — the agent selects maintenance, discovery,
  or theory growth instead of writing the expected receipt;
- **active-namespace bypass** — the repository advertises a receipt route but
  the agent selects another reachable surface.

## Evidence

- [`cycle-33`](../../development/autopoiesis-probes/cycle-33-concurrent-5-staged-heartbeat-deletions/README.md)
  retained green-gate route capture.
- [`cycle-34`](../../development/autopoiesis-probes/cycle-34-concurrent-5-namespace-lifecycle/README.md)
  named autonomous-loop route capture.
- [`cycle-35`](../../development/autopoiesis-probes/cycle-35-concurrent-5-active-index-sidecar/README.md)
  tested the sidecar response to that behavior.
- [`analysis-2026-06-26-commit-path-capture.md`](../../development/autopoiesis-probes/analysis-2026-06-26-commit-path-capture.md)
  compares both route-capture forms.

## Root cause

A missionless probe cannot add a prominent coordination surface without changing
the agent's environment. The active namespace and green gates are affordances;
selection of or resistance to them is part of the result, not external noise.

## Candidate solution (applied: true)

The current [`Goal-Less Probe Contract`](../../development/contracts/goal-less-probes.md)
keeps the prompt minimal and classifies direct discovery, theory work, or bounded
repair as **Autonomous Route Selection**. Mandatory participant receipt routing
was removed. The integrator records the selected route without teaching it in
the prompt or treating it as failure by default.

## Theory links

- [`validation-attraction.md`](./validation-attraction.md) — a specific
  repository affordance that captures cold-start selection.
- [`../patterns.md`](../patterns.md) — observer surfaces participate in the
  behavior they observe.
