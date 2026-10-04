# Sections Index: Quine arXiv Paper

## Purpose

This directory holds the sectionized manuscript source for the systems preprint.
Each file is one ordered manuscript section, split out so the paper can be edited
without flattening the whole draft into a single prose file.

## Ordering

- `00-abstract.md` is the entry abstract and should stay first
- `01-*.md` through `06-*.md` are the main narrative sections in paper order
- `appendix-*.md` holds supplemental material that supports the main argument

The numeric prefix is the ordering contract. Rename or insert sections by
preserving that prefix so the manuscript stays navigable.

## Primary Entry Points

- `00-abstract.md` for the paper-level abstract
- `01-introduction.md` for the opening claim and setup
- `02-posix-mapping.md` for the core POSIX correspondence
- `03-runtime.md` for the runtime implementation surface
- `04-properties.md` for emergent structural properties
- `05-related-work.md` for comparison and context
- `06-conclusion.md` for the closing claim boundary
- `appendix-a.md` for supplemental material

## Editing Rule

Use this surface for section-level prose only. Keep paper boundary decisions in
[`../README.md`](../README.md), workspace guidance in [`../WORKSPACE.md`](../WORKSPACE.md),
and argument ordering in [`../manuscript-plan.md`](../manuscript-plan.md).
