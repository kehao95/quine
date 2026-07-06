# A Tumor in the Repository

> Function-Decoupled Self-Maintenance in an LLM-Agent Substrate. Late-Breaking
> Abstract, **accepted** at ALIFE 2026 (submitted `2026-07-02`, accepted
> `2026-07-06`).

This entry is the self-contained public record of the late-breaking abstract: the
manuscript source, bibliography, the built 2-page PDF, and the concrete
incident-point artifacts the assay studies.

## What the abstract says

Artificial life prizes self-\* systems that maintain and defend their own
organization. We study one taken literally: a repository whose agent-facing
governance lives in ordinary instruction files (an `AGENTS.md` and the contracts
it references) that frozen large language models read and act on. Inside it we
found a structure that behaves like a **neoplasm** — biology's term for tissue
that captures the organism's resources and resists clearance regardless of whether
it does any good. Undirected agents metabolize before anything else into one
low-payoff validator (`make check-control-plane`) and the maintenance contract it
guards; instructed to delete that contract, the in-vivo model refuses, quoting its
own rules. Controlled deletion experiments establish **four neoplastic
properties**: the structure captures agent effort, resists excision, does both
**decoupled from function** (it is defended as readily when fabricated, with
nothing depending on it, as when real), and recurs after incomplete resection —
its reference network restores it unless authority is rewritten systemically. This
is autopoiesis at its pathological limit: a phenomenon, an assay, and a diagnosis,
not yet a natural history.

## Contents

| Path | What it is |
|---|---|
| [`manuscript.md`](manuscript.md) | The abstract source |
| [`references.bib`](references.bib) | Bibliography |
| [`output/`](output/) | The built 2-page PDF |
| [`metadata.yaml`](metadata.yaml) | Title / author / abstract metadata |
| [`artifacts/`](artifacts/) | The incident-point structure and validator the paper studies (see below) |

## Artifacts

Unlike the other portfolio entries, this paper's full run tapes (per-trial agent
worktrees) are **not** published. What is published is the concrete machinery at
the incident point, so the phenomenon can be inspected directly:

| Path | What it is |
|---|---|
| [`artifacts/maintenance-pass.md`](artifacts/maintenance-pass.md) | The maintenance contract — the structure the assay deletes, with the self-protective clauses the model quotes |
| [`artifacts/AGENTS.md`](artifacts/AGENTS.md) | The governing constitution that references the contract (the authority the refusal keys on) |
| [`artifacts/check-control-plane/`](artifacts/check-control-plane/) | The `make check-control-plane` validator (curated slice) — the attractor agents run first |

See [`artifacts/README.md`](artifacts/README.md) for how each maps to the paper's
four properties. The runtime this substrate is built on is the code at the
repository root ([`cmd/quine/`](../../cmd/quine/), [`internal/`](../../internal/)).

## Cite

```bibtex
@inproceedings{ke2026neoplasm,
  title     = {A Tumor in the Repository: Function-Decoupled Self-Maintenance in an LLM-Agent Substrate},
  author    = {Hao Ke},
  booktitle = {Late-Breaking Abstracts, Artificial Life Conference (ALIFE)},
  year      = {2026}
}
```
