# Workspace: Decomposing at the Wrong Level

The working-surface guide for the ICLR evaluation paper (six-dimension regime
decomposition of AI self-replication evaluation).

## Navigation

- paper dossier: [`README.md`](README.md)
- manuscript plan: [`manuscript-plan.md`](manuscript-plan.md)
- manuscript metadata: [`metadata.yaml`](metadata.yaml)
- manuscript references (DRAFT): [`references.bib`](references.bib)
- section source: [`sections/README.md`](sections/README.md)
- planning notes: [`plans/README.md`](plans/README.md)
  — incl. [`reconciliation-and-decisions-2026-06-30.md`](plans/reconciliation-and-decisions-2026-06-30.md)
- survey: [`survey/synthesis.md`](survey/synthesis.md) + [`survey/reports/`](survey/reports)
- citation provenance / source-check pool:
  [`references/citation-pool-2026-06-30.md`](references/citation-pool-2026-06-30.md)
- companion paper (phenomenon): [`../../archive/frozen-2026-07/aaai-facultative-succession/README.md`](../../archive/frozen-2026-07/aaai-facultative-succession/README.md) (frozen; AAAI-27 venue lost 2026-07-26)
- companion paper (safety): [`../../archive/frozen-2026-07/iclr-safeguard-composition/README.md`](../../archive/frozen-2026-07/iclr-safeguard-composition/README.md) (frozen)
- paradigm (authoritative): [`../../paradigms/ai-self-replication/README.md`](../../paradigms/ai-self-replication/README.md)
- theory objects: [`regime-decomposition`](../../../theory/objects/regime-decomposition.md),
  [`selection-signal-collapse`](../../../theory/objects/selection-signal-collapse.md),
  [`emergent-succession`](../../../theory/objects/emergent-succession.md)
- predictions: [`../../predictions/open/`](../../predictions/open/) (`regime-decomposition-P1..P7`)
- paper registry: [`../README.md`](../README.md)

## Directory role

The home for the **framework** paper of the three-paper self-replication arc.
It owns the six-dimension reference frame, the two formal results (Compression
Blindness Principle, Selection-Signal Collapse Theorem), the wrong-level
compression survey, the per-dimension demonstration, the evaluation-redesign
protocols, and predictions P1/P3/P4/P5/P6/P7.

## Current work mode (2026-06-30)

- status: **drafting** — scaffolding + section drafts under way (was: planning)
- framing: **A (reference frame)**, not a Kuhnian paradigm-shift claim (see
  `plans/reconciliation-and-decisions-2026-06-30.md`, D-2)
- six dimensions D1–D6 (the candidate's "five" is stale — D-1)
- evidence stance: **hybrid** — demonstrate via existing waves; propose the
  6-condition full-inversion as the discriminating future test (D-5)

## Paper boundary (do not cross — see reconciliation note D-3)

This paper OWNS the framework + CBP + Selection-Signal Collapse + P1,P3,P4,P5,P6,P7.

It must NOT develop (reserved for companions):

- the boundary-relative safeguard-composition result, guard-placement assay,
  inheritance/composition distinction, and continuity-only successor test → all
  reserved for [`iclr-safeguard-composition`](../../archive/frozen-2026-07/iclr-safeguard-composition/README.md)
- the facultative-morphology assay *details* (material × demand design) →
  reserved for [`aaai-facultative-succession`](../../archive/frozen-2026-07/aaai-facultative-succession/README.md);
  this paper cites the 6M-wave *result* as the D2 demonstration only.

This paper uses dorm-06 ONLY as D4 (selection-signal) / D6 (emergent) evidence,
never as a safety claim.

## Build

`make pdf` concatenates `sections/01..08` + `metadata.yaml` + `references.bib`
into `output/` and renders via the shared `../../templates/arxiv-acmart`. The
assembly step (`output/full_paper.md`) is verified working. The PDF step needs
the pandoc data templates the shared `acmart.latex` includes (e.g.
`document-metadata.latex`) present on pandoc's data path / `--data-dir`; this
build env lacks them, so the final `tectonic` render is currently blocked at the
template stage — a shared-template/environment dependency that affects every
arxiv-acmart paper here, not this paper's sources. Section headings are bare
(no manual `3.1` numbering); pandoc/LaTeX auto-numbers them.

## Control-plane rules

- paper-shaped theory → `README.md` + `manuscript-plan.md`
- manuscript source → `sections/` built via `Makefile` (`make pdf`); generated
  files under ignored `output/`
- do NOT auto-edit `scope: paradigm` theory objects / gaps / predictions; the
  paradigm README and `regime-decomposition` / `emergent-succession` objects and
  the `regime-decomposition-P*` predictions are human-directed — flag deltas as
  proposals (see reconciliation note S-1, S-2)
- `references.bib` is DRAFT (survey-extracted, many 2026 ids); source-check
  every entry against primary sources before submission
