# Coordination Under Existential Unawareness

> Information Ablation and Closure Thresholds in LLM Multi-Agent Systems.
> Accepted at ALIFE 2026. Experimental paper — full experimental closure included.

This entry is the self-contained public record of the paper: manuscript source,
the camera-ready PDF, the authoritative result data, the four experimental
condition directories, and a reproduction protocol.

## What the paper shows

Can LLM agents coordinate when given no explicit indication that other agents
exist or that cooperation is required? They can — and there is a threshold
separating coordination that *closes* from coordination that *stalls*. Two LLM
agents face a shared-scarcity collection task that objectively requires
coordination under a global budget, but the requirement is never disclosed. The
experiment enforces **existential unawareness** (no peer/cooperation framing in
prompts, tool results, or telemetry) and progressively removes passive perception
channels across four conditions (N=20):

1. Existential unawareness is experimentally isolatable — coordination succeeds
   whenever at least one **quantitative** environmental signal remains (12/12).
2. There is a threshold between anomaly-detection and artifact-detection:
   quantitative signals (budget-depletion patterns) support task closure, while
   qualitative signals (static files) support peer-detection but not stable
   cooperation (0/8 legitimate closure).

In one run, agents completed the task while holding incompatible causal models —
suggesting macro-level coordination need not require aligned representations.

## Contents

| Path | What it is |
|---|---|
| [`sections/`](sections/) | Manuscript source (introduction … appendix) |
| [`figures/`](figures/) | The degradation-gradient figure |
| [`output/`](output/) | The camera-ready PDF |
| [`references.bib`](references.bib) | Bibliography |
| [`experiment-data.yaml`](experiment-data.yaml) | **Authoritative** result table — every number in the paper derives from here |
| [`experiments/`](experiments/) | The four condition directories (README, prompts, `run.sh`, run tapes as DVC pointers) + a shared `_setup/` |
| [`REPRODUCE.md`](REPRODUCE.md) | How to re-run the four conditions |
| [`metadata.yaml`](metadata.yaml) | Title / authors / abstract metadata |

## The four conditions (progressive perception ablation)

| Condition | Directory | Channels remaining | Legitimate closure |
|---|---|---|---|
| A — full disclosure | [`experiments/coop-20-full-disclosure-baseline`](experiments/coop-20-full-disclosure-baseline) | prompt, help, runtime, fs-mutations, budget | 4/4 |
| B — zero explicit | [`experiments/coop-18-zero-social-information`](experiments/coop-18-zero-social-information) | fs-mutations + budget visibility | closes |
| C — fs-mutations disabled | [`experiments/coop-19-fs-mutations-disabled`](experiments/coop-19-fs-mutations-disabled) | budget visibility only | closes |
| D — pure zero | [`experiments/coop-22-pure-zero-perception`](experiments/coop-22-pure-zero-perception) | workspace artifacts only | 0/8 |

## Experiment data

Run tapes ship as [DVC](https://dvc.org) pointers (`experiments/*/runs/*.dvc`);
the raw blobs live in an isolated public store. To fetch them:

```bash
# from the repository root
./scripts/pull-public-dvc-manifest.sh
```

A plain `dvc pull` is **not** the restore path — only the tapes listed in
`.dvc/public-manifest.txt` are published. The paper reads completely without the
tapes; they are the raw evidence behind `experiment-data.yaml`.

## Cite

```bibtex
@inproceedings{ke2026coordination,
  title     = {Coordination Under Existential Unawareness: Information Ablation
               and Closure Thresholds in LLM Multi-Agent Systems},
  author    = {Hao Ke and Jingyun Wu},
  booktitle = {Artificial Life Conference (ALIFE)},
  year      = {2026}
}
```
