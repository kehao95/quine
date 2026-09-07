# Experiments Restructure — Migration Tooling

> Working surface for the `experiments/` family-restructure migration.
> Design + phase plan live in
> [`../../experiments/restructure-design.md`](../../experiments/restructure-design.md).
> This directory holds the **machine artifacts** that execute that plan.

## Summary

The migration drops the phase layer (`p1-borges/`, `p6-autopoiesis/`, …) as the
primary organizing principle and re-homes active experiments under
`experiments/active/<family>/<id>/`. These tools make the move deterministic and
reversible up to the Phase 4 deletion point.

## Contents

| File | Role |
|------|------|
| `id-map.tsv` | canonical legacy-path → new-path map (tab-separated): `legacy_path  new_path  legacy_id  new_id  family  disposition`. Source of truth for `repoint.py`. Built in Phase 0, consumed in Phase 3. |
| `repoint.py` | deterministic path repointer — literal substring replace of each legacy experiment path to its new family-home path across control-plane text files. Excludes run/tape payloads, generated heartbeat readouts, and `.dvc` pointers. DRY-RUN by default; `--apply` to write. |
| `fix-runner-depth.py` | adds one `../` to `${SCRIPT_DIR}/../../../<root-entry>` root-reach climbs (depth-3 → depth-4 after the move) and reports other cross-experiment climbs for manual Appendix G reconciliation. DRY-RUN by default; `--apply` to write. |

## Usage

```sh
# Phase 0 / pre-flight scope scan (dry-run, no writes):
python3 repoint.py --map id-map.tsv
python3 fix-runner-depth.py --only experiments

# Phase 3 (after git mv of experiments into active/<family>/):
python3 repoint.py --map id-map.tsv --apply
python3 fix-runner-depth.py --only experiments/active --apply
```

## Scope notes (Phase 0 findings)

- **194** control-plane files reference legacy phase paths (Paper 83, experiments
  84, development 15, scripts 6, others). These are the repoint surface.
- **199** runner scripts use a `${SCRIPT_DIR}/../../../` root-reach (the design
  doc estimated ~12) — the depth-change blast radius is the dominant mechanical
  risk. `fix-runner-depth.py` handles the uniform +1 case and flags the rest.
- The five `development/heartbeat/readouts/heartbeat-2026-06-12T*.md` snapshots
  carry ~26k legacy-path occurrences each; they are generated historical
  readouts and are **excluded** from repointing (same principle as retained
  tape evidence, Decision 7).

## Lifecycle

This is a temporary migration surface. When the migration completes (Phase 4+),
it is either deleted or folded into the completed-migration record alongside the
design doc.
