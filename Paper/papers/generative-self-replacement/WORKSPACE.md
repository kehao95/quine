# Generative Self-Replacement Workspace

## Navigation

- paper dossier: [`README.md`](README.md)
- HotOS draft: [`manuscript.md`](manuscript.md)
- manuscript plan: [`manuscript-plan.md`](manuscript-plan.md)
- HotOS one-page pitch: [`plans/hotos-one-page-pitch.md`](plans/hotos-one-page-pitch.md)

This file is the manuscript working-surface guide for the paper.

## Directory Role

Use this directory for manuscript source, figures, build assets, and venue-local working files.

Keep thesis shifts, evidence-promotion decisions, and paper-boundary decisions in the dossier rather than only in manuscript prose.

## Voice Boundary

This is a hardcore systems paper. The voice should be:

- precise technical language
- POSIX semantics as the primary framing
- empirical claims backed by verifiable artifacts
- no A-Life or philosophical framing in the main text
- if `6.64` / `6.65` are promoted, describe them as control-surface
  runtime-rebuild mechanism evidence, not as autopoiesis evidence

## Contents

| Path | Purpose |
|------|---------|
| `manuscript-plan.md` | current argument and section plan |
| `manuscript.md` | canonical prose source (to be created) |
| `figures/` | figure sources and paper-local visual assets |

## Evidence Sources

The original custom-body experiments live in:

- `experiments/active/autopoiesis-succession/6.4-self-replacement-handoff/`

Canonical runs:

- `20260321-112549-5.4` - first strict verification baseline
- `20260321-112852-5.4` - repeatability confirmation

Candidate stronger mechanism evidence now lives in:

- `experiments/active/autopoiesis-succession/succ-16-control-surface-viability-pressure/`
- `experiments/active/autopoiesis-succession/succ-17-minimal-control-surface-pressure/`

Current read:

- `6.4` shows non-isomorphic custom-body synthesis plus `exec`
- `6.64` / `6.65` show source rebuild of the fuller quine runtime body plus
  post-handoff control-surface viability
- the dossier should decide whether the paper is about the broader mechanism
  family or only the older non-isomorphic custom-body slice before any
  manuscript prose is drafted

## Build Notes

- target: arXiv or systems conference
- format: TBD based on venue selection
