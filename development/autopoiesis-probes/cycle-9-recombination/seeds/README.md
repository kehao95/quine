# Cycle-9 Seed Curation Index

## Role

This is the compression layer over the 36 raw cycle-9 seed files. It orders the
material without flattening the three decomposition projections that generated
it.

The raw directories are not duplicates to clean up. They are independent
projection files:

- `A-concepts/`: 12 concept seeds
- `B-tensions/`: 12 tension/paradox seeds
- `C-mechanisms/`: 12 mechanism, signal, prediction, and experiment seeds

Total retained seeds: 36. Unique seed IDs by filename: 36. Exact repeated IDs: 0.

## Compression Rule

Do not merge near-repeated raw seed files unless a later promoted object
explicitly absorbs their obligations. In cycle 9, robustness appears mostly as
near-repeat clusters rather than exact filename repeats: the same dimension is
visible as a concept, a preserved tension, and an assay path.

Promotion order:

1. Near-repeat clusters that connect concept + tension + mechanism.
2. Mechanism seeds that could become discovery-loop predictions or assays.
3. Single-projection seeds retained as searchable material until reuse,
   friction, or transport pressure makes them promotion candidates.

## Common-Refinement Clusters

These clusters are ordered by next useful promotion, not by theoretical grandeur.

| Promotion target | Raw seeds | Why first / next action |
|---|---|---|
| Sequence as carrier | `sequence-carrier-topology`, `sequence-carrier-assay`, `repair-vs-reproduction`, `ownership-location-vs-sequence`, `ownership-handoff-order-signal` | This is the primary child dimension from the synthesis. Immediate action: promoted as `Paper/theory/objects/sequence-carrier-topology.md`; next action is a repair-order assay. |
| Path viability equivalence | `path-viability-equivalence`, `endpoint-vs-path-viability`, `path-viability-equivalence-vector`, `intermediate-viability`, `intermediate-viability-signal` | Separates endpoint equality from viability through the path. Next action: define an equivalence vector before any scalar order-sensitivity score. |
| Reclosure and readiness | `reclosure-holonomy`, `reclosure-readiness`, `closure-vs-readiness`, `readiness-window-prediction`, `commutation-vs-confluence-vs-holonomy` | Names path-dependent closure residue. Next action: test whether a green endpoint is ready for the next perturbation. |
| Sensor / validator subjecthood | `sensor-subjecthood`, `sensor-subject-critical-pair`, `validator-judge-vs-organism-part`, `validator-hessian-coupling`, `repair-apparatus-mutation`, `repair-apparatus-mutation-classifier` | Keeps validators inside the organism rather than outside it. Next action: use a false-positive or stale-trigger case to decide whether content or sensor is the repair subject. |
| Scar and memory | `scar-memory`, `scar-debt-vs-memory`, `scar-memory-classifier` | Preserves the unresolved debt-vs-memory distinction. Next action: classify one real exception rule or routing detour as debt, immune memory, or unresolved wound. |
| Absorption / digestion lag | `addressability-vs-digestion`, `absorption-lag-assay`, `carrier-succession` | Distinguishes deposited text from metabolized behavior. Next action: test whether a freshly deposited seed becomes available to a later agent without direct prompt context. |
| Ownership sequence | `ownership-handoff-topology`, `ownership-location-vs-sequence`, `ownership-handoff-order-signal`, `local-fix-vs-control-plane-law` | Shows ownership as an operation-order relation, not only a final location. Next action: record a maintenance case where local fix first and control-plane fix first would preserve different evidence. |
| Coordinate / metric split | `equivalence-metric-family`, `coordinate-refinement-vs-metric-relaxation`, `quotient-vs-metric-test`, `commutator-basis-prediction` | Keeps the geometry tension alive: add a coordinate or relax the metric. Next action: compare several repair-pair residues before adding a new rule. |

## Immediate Absorption

The synthesis-level child dimension is now addressable as
[`sequence-carrier-topology`](../../../../Paper/theory/objects/sequence-carrier-topology.md).
That object is deliberately a SEED and not an umbrella for the whole seed set.
Future promotions should point back to the raw seeds they absorb rather than
deleting those seed files.
