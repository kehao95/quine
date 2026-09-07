# Active Doc Addressability Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


Active control-plane docs should not promise more addressability than the repository actually preserves.

Interpretation rules:

- active markdown docs should use repo-relative links for in-repo references; workstation-specific absolute paths are drift, not convenience
- a clickable markdown link in an active doc implies that the target is intentionally kept and currently materialized
- repo-root entry docs and `tests/` entry docs should keep a stable H1/H2 role spine even when their prose evolves, so navigation and audits remain machine-checkable
- if a root or tests entry doc changes role enough that the stable spine no longer fits, update the owning validator and hook docs in the same state transition rather than letting the doc drift silently
- if a cited artifact is no longer kept as a live file, do not leave a broken link behind; either archive and relink it, or downgrade the citation to a non-clickable runid/code path while naming its retention status honestly
- if an active doc links directly into `runs/` or similar generated evidence
  trees, that target should be repo-retained either directly in git or through
  a git-tracked adjacent `.dvc` pointer for the run/batch/campaign bundle;
  local-only material is not stable enough to advertise as clickable canonical
  state
- active clickable evidence should prefer stable top-level surfaces such as the run directory itself or top-level result/report/summary files; deep paths into `agent/`, `tapes/`, `context/`, `workspace/`, or per-command traces should usually be plain code paths instead
- historical captured outputs inside archived/generated evidence surfaces may preserve old absolute paths verbatim, but those snapshots should not be copied back into active control-plane prose as if they were portable links
- when a doc-heavy change touches active control-plane markdown, run `./scripts/check-active-doc-links.sh` and treat failures as control-plane drift rather than editorial polish
- run `./scripts/check-stable-evidence-surfaces.sh` when active docs cite experiment evidence; treat failures as a sign that the doc is leaning on a brittle artifact surface instead of a stable retained one
- when control-plane work spans experiments, papers, scenario docs, or active owner docs together, prefer `make check-control-plane` over ad hoc validator selection so the default audit surface stays centralized

Default maintenance question:

- if a future agent clicks every markdown link in the active doc surface around this change, do they reach live canonical state, or a path that only worked on one machine at one moment in time?

If the answer is the latter, the maintenance task is still open.
