---
title: "Quine: Realizing LLM Agents as Native POSIX Processes"
venue_label: "arXiv:2603.18030"
paper_id: quine-arxiv
claim_bundle: quine-arxiv
venue: arxiv
format-envelope: arxiv-preprint
track: system
cluster: process-native-systems-substrate
status: published
public_slug: quine-arxiv
published: 2026-03-08
preprint: true
state_mode: full-manuscript-home
role: flagship-minimal-proof
family: runtime-substrate
theme: runtime-ontology
theme_home: ../../../theory/views/substrate.md
publication_class: released-anchor
maturity_target: anchor
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [future-task-surface, exec-barrier, carrier, trace-inheritance]
gaps: [revisable-time-semantics]
evidence: [exp:p3-mrcr/3.1-needle-retrieval, exp:p3-mrcr/3.2-mixed-needle, exp:p4-harbor/4.3-local-wal-recovery, exp:p1-borges/3.2-fork-born-lineage]
view-affinity: [substrate]
---

# Quine (arXiv Systems Note)

## Current Question

What is the narrowest implementation-backed systems claim Quine can defend publicly today without borrowing authority from later runtime layers?

## Thesis

POSIX primitives are sufficient substrate for realizing LLM agents as systems software.

## In Scope

- agent = process as the core systems claim
- mission / material / deliverable / signal separation
- lifecycle semantics via `fork`, `exec`, and `exit`
- implementation-backed evidence and bounded behavior traces

## Out Of Scope

- full world semantics
- revisable time as a first-class runtime surface
- benchmark-heavy evaluation programs
- paradigm claims that depend on future runtime layers

## Current State

- public preprint: [`arXiv:2603.18030`](https://arxiv.org/abs/2603.18030) / [`10.48550/arXiv.2603.18030`](https://doi.org/10.48550/arXiv.2603.18030)
- released version: `v1`, submitted on `2026-03-08`
- local PDF snapshot: `data/2603.18030v1.pdf`
- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)

## Experiment Feeds

- `experiments/active/early-existence/exst-07-needle-retrieval/README.md` - primary - supporting-only - mission-first streaming and early exit are direct implementation-backed evidence for the current systems paper
- `experiments/active/early-existence/exst-09-mixed-needle/README.md` - secondary - supporting-only - explicit state tracking under multi-target retrieval supports the current systems claim about substrate-level agent sufficiency
- `experiments/active/benchmarks/bench-03-local-wal-recovery/README.md` - secondary - supporting-only - honest failure under visible runtime physics supports the current systems paper's implementation-backed substrate claim
- `experiments/active/early-existence/exst-06-fork-born-lineage/README.md` - secondary - supporting-only - repaired fork continuation sharpens the systems claim that native lifecycle operators can support bounded descendant continuation after parent death

## Current Evidence Backbone

- `experiments/active/early-existence/exst-07-needle-retrieval/README.md` - `evidence-set:mrcr-eight-sample-2026-02-12` - supporting-only - the current flagship systems note relies on the eight-sample Quine versus Baseline comparison as its clearest retrieval-sufficiency evidence
- `experiments/active/early-existence/exst-09-mixed-needle/README.md` - `evidence-set:mixed-needle-two-sample-2026-02-12` - supporting-only - multi-target retrieval provides a compact explicit-state-tracking support case
- `experiments/active/benchmarks/bench-03-local-wal-recovery/README.md` - `condition-set:wal-forced-claude-2026-03-09` - supporting-only - local WAL recovery gives the cleanest current honest-failure contrast under visible workspace physics

## Strongest Evidence

- current runtime implementation and tool surface
- behavior evidence that the physics can be discovered from prompt and schema
- compact manuscript examples that map one POSIX primitive to one bounded claim

## Blocking Gaps

- behavior baseline coverage is still uneven across newer physics surfaces
- the paper boundary must stay clean against later world/time papers
- the manuscript should keep implementation-backed sufficiency separate from broader paradigm framing
- reverse-link note: broader benchmark evidence should stay in later systems papers instead of inflating the flagship note

## Next Useful Work

- treat the public arXiv v1 as the released baseline for future diffs
- keep the manuscript narrow and reproducible
- promote only evidence that already has a visible artifact path
- route later semantics into their own dossiers instead of inflating this paper

## Unsafe Overclaims

- do not import `process-as-fs`, `subjective-reality`, or `revisable-time` as if they were already part of the released flagship claim
- do not treat implementation-backed sufficiency for this runtime as proof that full agent architecture is solved
- do not let a small number of strong runs stand in for universal model-independent substrate conclusions
