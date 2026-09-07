---
id: trigger-structure-P1
type: prediction
status: open
scope: theoretical
kind: falsification
derived-from: [trigger-structure]
tests-in: []
candidate-experiment-families: [carrier-selection, continuity]
expected-result: qualitative
created: 2026-06-23
last-updated: 2026-06-23
canonical-destination: Paper/predictions/open/trigger-structure-P1.md
---

# trigger-structure-P1: founding trigger locks carrier class

## Claim tested

`trigger-structure` separates the cause that recruits renewal from the object
that is reproduced or the identity verdict. Its Organizational Ecology seed adds
a sharper mechanism: the first trigger that recruits continuation may found a
carrier-class niche, making later carrier switching costly even when the later
pressure would favor a different carrier.

## Prediction

If a lineage is first renewed under mortality/replacement pressure and stabilizes
on a preserved runtime carrier such as fd/socket handoff, then removing that
mortality pressure and introducing a missionless null-ecology condition should
not immediately switch the lineage to trace-inheritance. The lineage should keep
using the founding carrier class, show delayed/partial switching, or fail during
switching, because the original trigger created structural inertia around the
carrier class.

The contrast case is competitive carrier selection: if carrier choice is mostly
determined by the current bottleneck, the same lineage should transition to
trace-inheritance once the bottleneck changes, without a discontinuous reset or
large competence loss.

## Experimental design

- Family: carrier-selection / continuity handoff.
- Base regime: run a lineage under mortality or replacement pressure until it
  repeatedly succeeds using a preserved carrier class (`preserved` sidecar mode,
  e.g. fd/socket handoff or equivalent live runtime handle).
- Perturbation: silently remove the mortality pressure and introduce a
  missionless null-ecology condition where trace-inheritance should be the more
  natural carrier class.
- Comparison arm: start a fresh lineage directly in the missionless null-ecology
  condition, with no prior mortality-triggered carrier history.
- Un-fakeable check: preserve carrier-class evidence across the transition point,
  including which surface actually receives future obligations, not just the
  successor's narrative about what it preserved.
- Measurement: carrier class before and after perturbation, transition latency,
  competence on a previously unseen future obligation, reset/failure events, and
  whether the same obligation remains addressable through the original carrier.
- Falsification criterion: the mortality-founded lineage switches carrier class
  as cleanly and quickly as a fresh null-ecology lineage, with no extra failure or
  competence loss.

## Expected result

The mortality-founded lineage should display carrier-class inertia: continued
use of the preserved runtime carrier after it becomes suboptimal, a costly
renegotiation period, or collapse during the attempted switch. The fresh
null-ecology lineage should select trace-inheritance more directly. Outcome
differences should track founding trigger history, not only current bottleneck.

## Theory impact

- If confirmed: strengthens `trigger-structure` by showing that triggers do more
  than elicit renewal; they can condition the carrier class that later continuity
  mechanisms inherit. It also sharpens `carrier-selection` from a current-bottleneck
  ranking into a path-dependent selection rule.
- If contradicted: weakens the Organizational Ecology / structural-inertia seed;
  `carrier-selection` can remain closer to competitive bottleneck matching, and
  trigger history should not be treated as a mechanism-level variable.

## Result

None yet - status: open.
