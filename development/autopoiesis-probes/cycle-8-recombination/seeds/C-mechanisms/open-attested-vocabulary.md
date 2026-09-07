---
id: open-attested-vocabulary
type: seed
kind: mechanism
status: SEED — cycle-8 decomposition C (mechanism lens)
origin: cycle-8-recombination collision
scope: theoretical
---
# open-attested-vocabulary

**What:** Validators use a preferred core plus attested extensions instead of closed enums.

Build this as a vocabulary pattern where known values are recommended and machine-checked for spelling, while novel values can be admitted if they are locally attested and indexed. The validator should flag unattested novelty, not forbid all novelty.

It proves useful when new recombination variants appear without breaking the control plane. It fails if openness erodes consistency or if closed enums produce less friction in practice.

## Connects to
Cycle-8 decomposition C, Mechanism seed 13: validator design with preferred core and attested extensions.
