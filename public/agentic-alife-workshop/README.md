---
source: Paper/papers/agentic-alife-workshop
---

# From Simulated Worlds to Process Habitats

> POSIX as a Substrate for Generative Artificial Life. Workshop paper, accepted
> at the Agentic & Generative AI as Artificial Life workshop, ALIFE 2026.

This entry is the self-contained public record of the workshop paper: the
manuscript source, bibliography, the built PDF, and the ALIFE 2026 workshop
talk slides.

## What the paper says

Artificial Life advances by changing the *medium* in which life-like
organization is studied. Classic systems (Tierra, Avida) withheld
organism-level goals but stayed below the semantic complexity of language;
today's LLM agents supply that complexity, yet the frameworks organizing them
predefine the coordination and succession ALife sets out to explain. Quine — a
POSIX-native runtime — instead realizes each agent as an operating-system
process, so lifecycle, memory, coordination, and succession reduce to OS
primitives the experimenter can vary. Withholding affordances for social
knowledge and reproduction, exploratory assays reveal heterogeneous coordination
protocols and a successor morphospace. The paper offers Quine as a
process-habitat substrate for making such organizational diversity
experimentally observable.

## Contents

| Path | What it is |
|---|---|
| [`manuscript.md`](manuscript.md) | The manuscript source |
| [`references.bib`](references.bib) | Bibliography |
| [`output/`](output/) | The built PDF |
| [`presentation/deck.pdf`](presentation/deck.pdf) | ALIFE 2026 workshop talk slides (10-min contributed talk, presented 2026-08-20) |
| [`metadata.yaml`](metadata.yaml) | Title / authors / abstract metadata |

The runtime this paper describes is the code at the repository root
([`cmd/quine/`](../../cmd/quine/), [`internal/`](../../internal/)).

## Experiments

The referenced assays ship under [`experiments/`](experiments/) as reproduction
skeletons — each carries its `README`, prompts, `run.sh`, analysis, and run tapes
as [DVC](https://dvc.org) pointers (`runs/*.dvc`). Raw run trees and bulk input
corpora are not shipped; the pointers are the data.

| Directory | Role |
|---|---|
| [`exst-01-ariadne`](experiments/exst-01-ariadne) | exec-enabled vs exec-disabled mortality pressure (death / continuity) — primary |
| [`trce-04-legacy-protocol-ablation`](experiments/trce-04-legacy-protocol-ablation) | persistent / wiped / read-only lineage-memory ablation — primary |
| [`coop-18-zero-social-information`](experiments/coop-18-zero-social-information) | zero-explicit-hint file-mediated coordination — primary |
| [`trce-08-successor-trace-uptake`](experiments/trce-08-successor-trace-uptake) | inherited environmental topology changing successor behavior — primary |
| [`succ-30-fixed-substrate-carrier-ladder`](experiments/succ-30-fixed-substrate-carrier-ladder) | successor morphology / carrier ladder — primary |
| [`trce-02-stigmergy`](experiments/trce-02-stigmergy) | filesystem-mediated stigmergy — supporting |
| [`coop-22-pure-zero-perception`](experiments/coop-22-pure-zero-perception) | artifact-aware coordination without fs/budget visibility — supporting |
| [`succ-35-aggregate-morphology-metrics`](experiments/succ-35-aggregate-morphology-metrics) | aggregate successor morphospace readout — supporting |

Fetch the run tapes with `./scripts/pull-public-dvc-manifest.sh` from the
repository root. Each `run.sh` uses the shared `experiments/_setup/generate-world.sh`.

## Cite

```bibtex
@inproceedings{ke2026process,
  title     = {From Simulated Worlds to Process Habitats: POSIX as a Substrate
               for Generative Artificial Life},
  author    = {Hao Ke and Jingyun Wu},
  booktitle = {Agentic and Generative AI as Artificial Life (Workshop),
               Artificial Life Conference (ALIFE)},
  year      = {2026}
}
```
