# Cycle-8 Seed Curation Index

## Role

This is the compression layer over the 136 raw cycle-8 seed files. It orders the
material without flattening the projections that generated it.

The raw directories are not duplicates to clean up. They are independent
projection files:

- `A-concepts/`: 38 concept seeds
- `B-tensions/`: 40 tension/paradox seeds
- `C-mechanisms/`: 58 mechanism, signal, prediction, and experiment seeds

Total retained seeds: 136. Unique seed IDs by filename: 127. Exact repeated IDs:
6.

## Compression rule

Do not merge same-named raw seed files unless a later promoted object explicitly
absorbs their obligations. A same-name recurrence across projections is a
common-refinement candidate: it means different projections independently found
the same axis.

Promotion order:

1. Exact repeats across projections: highest confidence that the axis is not a
   lens artifact.
2. Near-repeat clusters that connect concept + tension + mechanism: best next
   source for theory objects, gaps, or predictions.
3. Single-projection seeds: retain as searchable material until reuse, friction,
   or transport pressure makes them promotion candidates.

## Exact repeats

| Seed ID | Projections | Compression read |
|---|---|---|
| `commutator-friction` | A, B, C | Strongest common-refinement signal: order sensitivity as evidence that the current representation is missing a dimension. Promote before subtler clusters. |
| `conflict-radius` | A, B, C | Measures how much structure must move to admit one honest new instance; useful as a repair-cost / topology signal. |
| `selection-shadow` | A, B, C | Early evaluation suppresses novelty before recombination; direct method-risk for future divergent-collision runs. |
| `canonical-target-branched-path` | A, B | Reconciles single canonical destination with divergent paths; compresses the main governance tension. |
| `lineage-retention` | A, C | Names the retention condition for future recombination; likely belongs with addressability/deposition metrics. |
| `narrative-hybrid` | A, C | Failure mode: child sounds like parents but preserves no obligations; should pair with parent-recoverability tests. |

## Common-refinement clusters

These clusters are ordered by next useful promotion, not by theoretical grandeur.

| Promotion target | Raw seeds | Why first / next action |
|---|---|---|
| Missing-dimension signal | `commutator-friction`, `order-sensitivity-index`, `common-refinement-child-test` | Gives the method a detector: if valid operations do not commute, the repo may lack a dimension. Next action: derive a prediction or assay before hardening it into a control-plane rule. |
| Functional child validity | `parent-recoverability`, `fusion-addressability`, `formal-novelty-functional-continuity`, `narrative-hybrid`, `parent-recoverability-contract`, `fusion-addressability-assay`, `parent-recoverability-utility-prediction` | Prevents recombination from becoming prose synthesis. Next action: define a parent-obligation checklist for one real promoted child object. |
| Reproductive surface / membrane | `addressable-genotype`, `recombinability-not-openness`, `openness-recombinability`, `deposited-lineage`, `semi-permeable-seed-surface`, `permeability-noise`, `porosity-overload`, `porosity-ab-assay` | Names the substrate condition: open notes are not enough; material must be bounded, labeled, heritable, and not over-porous. Next action: use this index as the first semi-permeable seed surface and measure fresh-agent usefulness later. |
| Divergence substrate | `projection-object-type`, `intra-snapshot-branching`, `divergence-preserving-deposition`, `addressability-divergence-flattening`, `git-whole-tree-node-divergence`, `projection-object`, `versioned-concept-fragments`, `collision-record-carrier`, `collision-record` | Directly supports `divergence-substrate`: one file per concept flattens projection diversity. Next action: compare `projection-object` vs `versioned-concept-fragments` on one future concept. |
| Timing / selection | `branching-potential`, `compression-branching-potential`, `short-window-selection-blindness`, `anti-compression-window`, `anti-compression-window-experiment`, `selection-shadow` | Keeps compression from firing before recombination yield is visible. Next action: keep this curation index instead of deleting raw seeds, then test reuse. |
| Graph primitive | `recombinational-neighborhood`, `recombination-hyperedge`, `parentage-signature`, `child-structural-signature`, `functional-dependency-signal`, `functional-dependency-removal-test` | Points to a many-to-one relation type, but needs one real child object first. Next action: trial parentage/obligation fields manually before adding schema or validators. |
| Procedure as material | `collision-not-consensus`, `consensus-recombination`, `tension-preserving-integration`, `procedure-reification`, `procedure-manifest`, `procedure-reification-experiment`, `synthesis-before-refutation-assay` | Extends the discovery loop from a live habit into recombinable procedure objects. Next action: keep procedure changes in `development/contracts/divergent-collision.md`, not only in narrative synthesis. |
| Hygiene / contamination | `gamete-contamination`, `gamete-contamination-clean-collision`, `read-only-probe-sandbox`, `light-prompt-divergence`, `light-heavy-prompt`, `light-vs-heavy-prompt-assay` | Protects attribution. Next action: future probes should use read-only or disposable worktrees before claiming lineage independence. |
| Trigger vs validation | `trigger-validation-split`, `internal-collision-external-validation`, `mechanism-science-validation`, `multi-lens-yield-assay`, `deposition-vs-chat-assay`, `prediction-density-prediction` | Holds the boundary between a generative collision and a validated scientific claim. Next action: promote one C prediction seed through the discovery loop. |

## Immediate absorption

The synthesis-level child dimension is now addressable as
[`reproductive-code-topology`](../../../../Paper/theory/objects/reproductive-code-topology.md).
That object is deliberately a SEED and not an umbrella for the whole seed set.
Future promotions should point back to the raw seeds they absorb rather than
deleting those seed files.

## Projection maintenance

- source owner: files under `A-concepts/`, `B-tensions/`, and `C-mechanisms/`
- projection role: ordered curation over retained seed projections after positive overload
- freshness trigger: seed file add/remove/rename, exact-repeat count change, or promotion of a seed into a theory object/gap/prediction
- drift check: recompute exact repeats from seed filenames and compare promoted objects against the immediate absorption section
- absorption owner: raw projection files remain lineage evidence; promoted claims move to `Paper/theory/objects/`, `Paper/theory/gaps/`, or `Paper/predictions/`
