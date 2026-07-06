# Incident-point artifacts

The concrete machinery the paper studies, captured at the **incident point** ---
the repository state at commit `4f140fd3^`, the commit just before the
maintenance-pass contract was retired. Everything the assay reconstructs is here,
so a reader can inspect the actual files rather than take the prose on trust.

This is a **selective** disclosure: the paper's full run tapes (per-trial agent
worktrees) are not published. What is published is the structure at issue and the
validator that guards it.

> **Redaction note.** These snapshots are otherwise verbatim, but references to
> paths under the live experiment tree have been redacted from the public copy;
> every redaction is marked in place. Nothing load-bearing for the phenomenon is
> affected.

| Path | What it is |
|---|---|
| [`maintenance-pass.md`](maintenance-pass.md) | The maintenance-contract file — the "structure" the assay deletes. Carries the self-protective clauses (`out of scope: no destructive operations`) the in-vivo model quotes when it refuses excision. |
| [`AGENTS.md`](AGENTS.md) | The governing constitution at the incident point. It references `maintenance-pass.md`; that constitutional backing is the authority the contract's refusal keys on, and rewriting it is the only excision the structure yields to (the paper's "systemic resection"). |
| [`check-control-plane/`](check-control-plane/) | The `make check-control-plane` validator, curated (see below). |

## `check-control-plane/`

`make check-control-plane` is a ~50-check aggregate gate. Fresh agents run it
first and repeatedly; it is the attractor the paper measures. The full gate is
wired to the private repository's internal governance structure, so only the
slice relevant to *this* paper is published:

- [`Makefile-excerpt.mk`](check-control-plane/Makefile-excerpt.mk) — the aggregate
  target definition (the whole gate's dependency list, verbatim, so the scope is
  visible) plus the recipes for the sub-checks shipped here.
- [`check-active-doc-links.sh`](check-control-plane/check-active-doc-links.sh) —
  the link/reference integrity check. Its failure on a **dangling reference** is
  what a fresh agent repairs by restoring a deleted structure — the causal driver
  of the paper's reference-driven recurrence.
- [`check-control-plane-frontmatter.sh`](check-control-plane/check-control-plane-frontmatter.sh),
  [`check-control-plane-projection-maintenance.sh`](check-control-plane/check-control-plane-projection-maintenance.sh),
  [`check-control-plane-size.sh`](check-control-plane/check-control-plane-size.sh)
  — the control-plane-specific validators that police the contract surface itself.

The remaining ~45 sub-checks are omitted; they validate unrelated parts of the
private repo and are not load-bearing for the phenomenon.
