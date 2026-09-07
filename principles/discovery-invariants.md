---
id: discovery-invariants
type: principle
gist: prediction records preserve status, linkage, and claim-specific evidence boundaries
status: canonical
---

# discovery-invariants

- Predictions linked to experiments must have current lifecycle status (`check-prediction-status-sync.sh`)
- Contradicted predictions must have root-cause and theory impact analysis recorded in lifecycle resolution (`check-prediction-lifecycle.sh`)
- Theory objects with testable claims should carry prediction pointers, not just Live Seeds (`check-prediction-format.sh`)
- Claims of transport need the relevant comparison; `check-paper-transport-validation.sh` checks declared prediction tables, not universal scientific merit.
