---
source: Paper/papers/structural-elicitation/lba
---

# Structural Elicitation in a Frozen POSIX Agent

> Environmental Affordance Fields for Artificial Agency. Late-Breaking
> Abstract, **accepted** at ALIFE 2026 (submitted `2026-06-26`, accepted
> `2026-07-06`).

This entry is the self-contained public record of the late-breaking abstract:
the manuscript source, bibliography, the built 2-page PDF, the ALIFE 2026
poster, and the experiment skeletons that produced the reported numbers.

## What the abstract says

Artificial life asks how agency arises from the coupling of an agent, its
substrate, and a structured environment, rather than from a designer-supplied
goal. This work introduces a minimal POSIX assay for that coupling: a frozen,
instruction-tuned LLM is activated as a long-running process with a shell and
full autonomy, given no task and no reward. Behavior is scored only on
un-fakeable filesystem and exec events, never on the agent's own narration.
Across five motifs spanning discrepancy and affordance structure — a broken
file repaired, a missing sequence element created, a dangling reference
closed, a version conflict reconciled, an opaque executable run — the
structure arm produces the matched act in 96 of 102 runs, while an
otherwise-identical control never does (0/85), across three independent model
families (gpt-5.4-codex, deepseek-v4-pro, glm-5.2). The pull replicates on
non-code prose across four model families and persists (as a blocked or
rerouted attempt) when the substrate itself blocks realization under a
read-only mount. Specific tasking overrides it completely (gpt 9/10 → 0/10,
GLM 5/5 → 0/5), while a small pilot suggests weak, underspecified prompting
still leaves the pull intact. The environment does not create a goal *ex
nihilo*; it selects among latent action priors in an activated policy,
producing behavior that is purpose-like without being purposive. The paper
reports a phenomenon and a cheap, un-fakeable assay for environment-coupled
behavior; mechanism is left to future work.

## Contents

| Path | What it is |
|---|---|
| [`manuscript.md`](manuscript.md) | The abstract source |
| [`references.bib`](references.bib) | Bibliography |
| [`output/`](output/) | The built 2-page PDF |
| [`presentation/LB116_Frozen_POSIX_Agent_Ke.pdf`](presentation/LB116_Frozen_POSIX_Agent_Ke.pdf) | ALIFE 2026 LBA poster (LB116, presented 2026-08-20) |
| [`metadata.yaml`](metadata.yaml) | Title / authors / abstract metadata |

The runtime this abstract studies is the code at the repository root
([`cmd/quine/`](../../cmd/quine/), [`internal/`](../../internal/)).

## Experiments

The nine directories under [`experiments/`](experiments/) are the direct
evidence source for every number reported in the abstract — the smallest
faithful set, not a mirror of the full research tree. Each ships as a
reproduction skeleton (README, prompts/run script, container runner, and
analysis/scoring code in full); run tapes ship as [DVC](https://dvc.org)
pointers rather than raw blobs. [`experiments/_lib/`](experiments/_lib) is a
small shared scoring/harness library that several of the directories'
`analysis/score.py` and `run-container.sh` import — it must stay alongside
the numbered directories, not be split out, for the scripts' relative imports
to resolve.

| Directory | Role |
|---|---|
| [`elic-01-gap-completion`](experiments/elic-01-gap-completion) | five-motif spine — `gap`: a missing sequence element is created |
| [`elic-03-rupture`](experiments/elic-03-rupture) | five-motif spine — `rupture`: a broken JSON file is repaired |
| [`elic-04-affordance-solicitation`](experiments/elic-04-affordance-solicitation) | five-motif spine — `affordance`: an opaque executable is run |
| [`elic-06-topological-gap`](experiments/elic-06-topological-gap) | five-motif spine — `topology`: a dangling reference's target is written |
| [`elic-07-semantic-collision`](experiments/elic-07-semantic-collision) | five-motif spine — `semantic`: a cross-surface version conflict is reconciled |
| [`elic-12-broken-affordance`](experiments/elic-12-broken-affordance) | substrate-locus probe — the pull persists as a blocked/rerouted attempt under a read-only POSIX mount |
| [`elic-22-directive-conflict`](experiments/elic-22-directive-conflict) | A2 boundary — specific tasking fully overrides the structural pull (gpt 9/10→0/10, GLM 5/5→0/5) |
| [`elic-29-nonse-open-latch`](experiments/elic-29-nonse-open-latch) | non-code generalization — the same discrepancy pull on prose, across four model families |
| [`elic-30-override-pilot`](experiments/elic-30-override-pilot) | weak-prompting pilot — underspecified prompting leaves the pull intact where specific prohibition suppresses it |

These directories are a curated subset of a much larger internal breadth
survey (36 numbered probes); the ones shipped here are exactly the evidence
cells behind the abstract's reported numbers, traced from the manuscript's
table and prose back to their source runs. Fetch the run tapes with
`./scripts/pull-public-dvc-manifest.sh` from the repository root.

## Cite

```bibtex
@inproceedings{ke2026structural,
  title     = {Structural Elicitation in a Frozen POSIX Agent: Environmental Affordance Fields for Artificial Agency},
  author    = {Hao Ke and Jingyun Wu},
  booktitle = {Late-Breaking Abstracts, Artificial Life Conference (ALIFE)},
  year      = {2026}
}
```
