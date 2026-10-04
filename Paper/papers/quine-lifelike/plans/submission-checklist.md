# ALIFE Conference Submission Checklist

This file is build-facing and submission-facing. It should use the same public language standard as the manuscript itself.

- do not mirror internal experiment numbering, run IDs, condition codes, or redesign chronology here
- refer to manuscript-visible contrasts by their public labels (`continuity contrast`, `persistence ablation`, `maintenance-horizon comparison`) rather than by dossier shorthand
- if a checklist item needs internal provenance, link back to the owning dossier instead of importing the internal wording into this file

## Current Build State

- the manuscript currently compiles to
  `output/Separable_Viability_Pressures_in_a_Process-Native_ALIFE_Habitat.pdf`
- the latest local `make pdf` run produced a 9-page draft including bibliography
- the build is usable but not warning-free: Tectonic reports one `lineno.sty` UTF-8 warning plus non-blocking font and underfull-box warnings

## Pre-Build Checks That Still Hold

- title and abstract in [`metadata.yaml`](../metadata.yaml) match the current pressure-centered claim
- the section contract now uses the current five-section main-text structure plus appendix from [`manuscript-plan.md`](../manuscript-plan.md)
- bibliography maintenance continues in [`../references.bib`](../references.bib)
- the build should target the current files:
  - `sections/01-introduction.md`
  - `sections/02-substrate-pressures-evidence-rule.md`
  - `sections/03-results.md`
  - `sections/04-discussion.md`
  - `sections/05-conclusion.md`

## Current Build Audit

- current output PDF path:
  `output/Separable_Viability_Pressures_in_a_Process-Native_ALIFE_Habitat.pdf`
- compiled page count: `9`
- citation-clean build status: bibliography is present; no undefined-citation warning was surfaced in the current log
- figure placement and caption audit
- non-blocking TeX warnings worth ignoring: `lineno.sty` UTF-8 warning, font-shape fallbacks, underfull box warnings, and current `xdvipdfmx` duplicate-object warnings

## Current Metadata Defaults

- [`metadata.yaml`](../metadata.yaml) uses a shared affiliation line for both listed authors; no second affiliation line is currently needed
- the submission draft intentionally omits an acknowledgements section
- the submission draft currently includes no funding statement
- `code_url` remains the public repository URL in [`metadata.yaml`](../metadata.yaml)

## Remaining Human Confirmation

- confirm that the current omission of acknowledgements is intentional for submission
- confirm that no funding disclosure is required
- confirm that the shared-affiliation author block matches the intended submission metadata

## Last Useful Checks Before Submission

- do one full progressive-disclosure pass:
  - each section should present the main effect before auxiliary controls, boundary conditions, or philosophical caveats
  - each subsection should move from strongest causal contrast -> conditions of uptake or failure -> boundary of generalization
  - discussion should focus on `so what` rather than re-running the Results ladder in the same wording
- manually skim the first page of the regenerated PDF for author block spacing and abstract wrap
- manually skim the references page for capitalization, ordering, and broken URLs
- verify figure captions use the same terms as the body: `environmental organ`, `post-exit maintenance`, `process-native ALIFE habitat`
- manually inspect every figure for text overlap, clipped labels, legend collisions, and caption/figure terminology mismatch
- do one last read for claim discipline: no drift toward `life`, `autopoiesis`, or model-universal claims
