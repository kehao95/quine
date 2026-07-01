# Facultative Self-Reproduction in Quine

> a Process-Native LLM Medium. Late-Breaking Abstract, accepted at ALIFE 2026.

This entry is the self-contained public record of the late-breaking abstract:
the manuscript source, bibliography, and the built 2-page PDF.

## What the abstract says

Self-reproduction is central to artificial life, but many digital media make it
measurable by fixing the successor's schema in advance — what counts as
offspring, where it may appear, which channel carries heritable state, and under
what norm reproduction is required. This work studies Quine, a process-native LLM
medium, with a **material-by-future-demand assay across real POSIX `exec`
boundaries**. In each condition the post-`exec` demand is designed so that full
self-reproduction *would* satisfy the challenge, but reproduction is neither
required nor guaranteed by the available material. The resulting successors vary
widely — from ordinary POSIX utilities and stream processors to LLM-backed agents
and rebuilt runtimes — and the *same* carried source can become either a rebuilt
runtime or a degenerate stream body depending only on the demand it faces.
High-fidelity self-reproduction thus appears not as a fixed property of the
medium but as a **facultative outcome**: stronger than copying a body (runtime
organization must be reconstructed across mortality), weaker than full autopoiesis
(reconstruction stays scaffolded by LLM inference, external tooling, and
externally posed demand).

## Contents

| Path | What it is |
|---|---|
| [`manuscript.md`](manuscript.md) | The abstract source |
| [`references.bib`](references.bib) | Bibliography |
| [`output/`](output/) | The built 2-page PDF |
| [`metadata.yaml`](metadata.yaml) | Title / authors / abstract metadata |

The runtime this abstract studies is the code at the repository root
([`cmd/quine/`](../../cmd/quine/), [`internal/`](../../internal/)).

## Experiments

The material-by-future-demand succession assays ship under
[`experiments/`](experiments/) as reproduction skeletons (README, prompts,
`run.sh`, analysis, and run tapes as [DVC](https://dvc.org) pointers):

| Directory | Role |
|---|---|
| [`succ-30-fixed-substrate-carrier-ladder`](experiments/succ-30-fixed-substrate-carrier-ladder) | the fixed-substrate carrier ladder — the same carried source becoming a rebuilt runtime or a degenerate stream body |
| [`succ-35-aggregate-morphology-metrics`](experiments/succ-35-aggregate-morphology-metrics) | aggregate successor morphospace readout |

Fetch the run tapes with `./scripts/pull-public-dvc-manifest.sh` from the
repository root.

## Cite

```bibtex
@inproceedings{ke2026facultative,
  title     = {Facultative Self-Reproduction in Quine, a Process-Native LLM Medium},
  author    = {Hao Ke and Jingyun Wu},
  booktitle = {Late-Breaking Abstracts, Artificial Life Conference (ALIFE)},
  year      = {2026}
}
```
