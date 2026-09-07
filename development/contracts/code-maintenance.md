# Code Maintenance Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


Recent drift exposed a specific maintenance defect: implementation semantics changed, but preserved baselines and nearby control-plane text still encoded the old semantics. The problem was not one missed edit. The control plane had not stated clearly enough that checked-in evidence is part of code maintenance.

Interpretation rules:

- maintaining code includes maintaining the checked-in surfaces that restate what the code now means
- after a semantic change, inspect prompts, scenario registry, runner env, coverage maps, preserved baselines, and nearby experiment or paper summaries for stale wording or dead knobs
- preserved baselines are current executable examples, not a museum of every prior field name, threshold, or prompt contract
- when a baseline run encodes superseded semantics, replace or remove it once a current passing baseline exists; only keep historical runs when they are explicitly archived as historical evidence in the owning experiment surface
- if a variable, threshold, or scorer rule is no longer live in code, remove it from runners and docs instead of letting an inert compatibility shell persist
- a code change is not complete while the repository's preserved evidence still teaches a different contract than the implementation enforces

Default maintenance question:

- if a future agent read only the checked-in docs, baselines, and run trees around this feature, would they learn the contract the code implements today or a stale one?

If the answer is stale, the maintenance task is still open.
