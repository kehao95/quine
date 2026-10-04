# Figures

This directory holds generated figure assets for the ALIFE conference paper.

Current generated figure set:

1. `fig-01-host-guest.{pdf,png}` - substrate schematic for the host/process boundary and survival channels
2. `fig-02-ariadne-contrast.{pdf,png}` - `1.1` continuity contrast (`exec` enabled vs blocked)
3. `fig-03-pressure-gradient.{pdf,png}` - optional pressure-gradient view for the continuity subsection
4. `fig-04-guardian-tiers.{pdf,png}` - cross-model same-prompt horizon summary across `H0 / H1 / H2`, including clean baseline, live restoration, post-exit restoration, and invalid horizon misses
5. `fig-05-guardian-phenotypes.{pdf,png}` - maintenance-layer gradient from one-shot repair to host-level persistence
6. `fig-06-ablation-summary.{pdf,png}` - persistence-ablation summary for the environmental-organ section

Current paper status:

- section drafts were regenerated on `2026-04-12`, and the current draft compiles against them
- figure placement remains provisional until the manuscript prose stabilizes, but numbering should follow the regenerated `sections/` files and current five-section manuscript plan rather than deleted pre-reframe section labels

Workflow:

- regenerate or update figures through [`../scripts/generate_figures.py`](../scripts/generate_figures.py) when the current manuscript needs them
- if a figure caption or purpose changes materially, update this file so figure intent stays aligned with the current paper framing
