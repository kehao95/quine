---
id: joint-mdl-recombination-gate
type: seed
kind: mechanism
status: SEED — cycle-8 decomposition C (mechanism lens)
origin: cycle-8-recombination collision
scope: theoretical
---
# joint-mdl-recombination-gate

**What:** A child is useful if `MDL(child) + integration_cost < MDL(A) + MDL(B)`.

Build this as a formal scoring gate for promoted recombinations, not for first deposition. Estimate description length of the child, the cost of integrating it into the graph, and the description length of leaving parents separate.

It proves useful if it distinguishes real compression from rhetorical blending. It is disproved if high-value recombinations consistently fail the inequality or if low-value hybrids pass it because costs are poorly measured.

## Connects to
Cycle-8 decomposition C, Mechanism seed 11: formal scoring gate for recombination.
