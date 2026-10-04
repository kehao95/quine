# Track A: System Paper (arXiv Preprint)

**Target**: arXiv cs.OS / cs.AI  
**Length**: ~3-4 pages  
**Status**: PUBLIC PREPRINT (`arXiv:2603.18030v1`, submitted 2026-03-08)

## Navigation

Paper-level control now lives alongside this workspace:

- paper dossier: [`README.md`](README.md)

## Directory Structure

This directory is the **manuscript source** for the active systems preprint.

The build logic is shared with **[`../../templates/arxiv-acmart/`](../../templates/arxiv-acmart)** so future papers can reuse the same scaffold without copying ad hoc commands out of this track.

```
papers/quine-arxiv/
├── Makefile
├── metadata.yaml
├── references.bib
├── data/
│   └── 2603.18030v1.pdf
├── sections/
│   ├── 00-abstract.md
│   ├── 01-introduction.md
│   ├── 02-posix-mapping.md
│   ├── 03-runtime.md
│   ├── 04-properties.md
│   ├── 05-related-work.md
│   ├── 06-conclusion.md
│   └── appendix-a.md
├── figures/
├── support/              # track-local overrides for submission assets
├── output/               # generated artifacts; safe to regenerate
└── templates/
```

Strategic notes, changelogs, and deferred ideas live in
[`archive/track-system-meta/`](../../archive/track-system-meta/README.md).
Claim-boundary decisions should be recorded in the dossier, not only in section files here.

`output/` is treated as disposable build output rather than canonical source material.
`support/` and `templates/` can override the shared scaffold, but the default pattern for new papers is to start from `../../templates/arxiv-acmart/`.
`data/2603.18030v1.pdf` is the local snapshot of the currently released arXiv version and should be preserved as the external publication baseline.

## Build Commands

```bash
# Generate PDF
make pdf

# Generate LaTeX
make tex

# Generate ArXiv submission package
make arxiv

# Clean all generated files
make clean
```

## Manuscript Focus

This version is intentionally narrow:

- strict POSIX mapping
- recursive runtime realization
- structural properties that follow from the mapping
- restrained future-work sketch

## Figures

1. **Architecture.svg** — Host-Guest separation, POSIX mapping, Read-Eval-Print loop
2. **Lifecycle.svg** — `fork()` for parallelism, `exec()` for default self-renewal
