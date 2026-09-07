# Agent Diff-Review & Theory-Blindspots Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


## Dual Architecture

Agent diff review operates through two strictly separated channels:

```
┌─────────────────────────────────────────────────────────────┐
│                     Agent Diff-Review                       │
├──────────────────────────────┬──────────────────────────────┤
│      GATE (Semantic Guard)   │   MUSE (Advisory Spark)      │
│  - High-severity overclaims  │  - Lateral connections       │
│  - Objective contradictions  │  - Unclaimed evidence leads  │
│  - Fail-open on infra error  │  - Strictly NON-BLOCKING     │
│  - Blocks only on dishonesty │  - Advisory nudges only      │
└──────────────────────────────┴──────────────────────────────┘
```

1. **GATE (Semantic Invariant Guard):**
   - Focuses strictly on issues where regex/AST linters are blind: clear factual contradictions, ungrounded empirical overclaims, or mislabeled artifacts.
   - **Fail-Open:** Network, API, or quota errors must fail open (`exit 0`) so development velocity is never held hostage by external provider availability.
   - High threshold: Only blocks when an objective human reviewer would reject the claim.

2. **MUSE (Advisory & Non-Blocking):**
   - Surfaces subtle blind spots, alternative theoretical interpretations, and unclaimed experimental evidence.
   - **Strictly Non-Blocking:** MUSE commentary is advisory and never halts a commit or automated workflow.

## Valence & Cadence Separation

- **Per-Commit Review (`check-agent-diff-review.sh`):** Focused, fast, low-friction review scoped to the immediate staged diff in conceptual and paper surfaces.
- **Corpus-Level Review (`make theory-blindspots`):** On-demand, coarse-grained analysis across papers, evidence maps, and theory files to identify systemic blind spots and macro directions.

## Review Principles

- **Reviewers as Collaborators:** The review layer exists to sharpen ideas, not multiply bureaucratic gates.
- **Sparse, High-Signal Feedback:** A noisy reviewer gets ignored; keep blocking gates rare and definitive, and keep generative suggestions concise.
