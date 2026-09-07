---
id: conflict-radius
type: seed
kind: signal
status: SEED — cycle-8 decomposition C (mechanism lens)
origin: cycle-8-recombination collision
scope: theoretical
---
# conflict-radius

**What:** Count the files, validators, and docs that must change for one new instance.

Build this as a brittleness detector: introduce a new legitimate variant and record how far the required repair spreads. A large radius indicates the current structure encodes assumptions in too many places.

It proves useful if high-radius changes predict maintenance friction. It fails if broad changes are rare, intentional, or not costly in practice.

## Connects to
Cycle-8 decomposition C, Signal seed 25: brittleness detector for new instances.
