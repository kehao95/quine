# Cycle-11 Recombination Probe

## What this directory is

Retained output from a divergent-collision theory-growth pass on
`carrier-selection`, the open question of which carrier class wins under which
bottleneck and what makes a trace become load-bearing.

Status: generative `SEED` work. The outputs below are theoretical growth and
compression artifacts, not transport-validated scientific claims.

## Navigation order

1. [`target-selection.md`](./target-selection.md) - target pick and reasoning
2. `growth/` - three divergent lens growths
3. [`synthesis.md`](./synthesis.md) - tension-preserving collision synthesis
4. [`seeds/README.md`](./seeds/README.md) - curation index over the 31 retained
   seed files
5. [`compression.md`](./compression.md) - ordering, promotion decision, and next moves

Raw seed directories preserve the three decomposition projections:

- [`seeds/A-concepts/`](./seeds/A-concepts/) - concept seeds
- [`seeds/B-tensions/`](./seeds/B-tensions/) - tension/gap seeds
- [`seeds/C-mechanisms/`](./seeds/C-mechanisms/) - mechanism, assay, and
  prediction seeds

## Compression decision

Do not merge repeated raw seed files during the first compression pass. Repeated
concepts across projections are common-refinement evidence. The compression
layer is the curation index, plus one promoted SEED object:
[`selector-stack-topology`](../../../Paper/theory/objects/selector-stack-topology.md).

## Projection maintenance

- source owner: raw outputs in this directory and promoted theory objects under
  `Paper/theory/objects/`
- projection role: navigation and compression index for the cycle-11
  divergent-collision pass
- freshness trigger: adding/removing growth files, decomposed seeds, or promoted
  theory objects from this cycle
- drift check: compare this navigation order with files under this directory and
  any cycle-11 references in `Paper/theory/objects/`
- absorption owner: stable theory lands in `Paper/theory/objects/`; testable
  claims land in `Paper/predictions/`; raw probe artifacts remain here as lineage
  evidence
