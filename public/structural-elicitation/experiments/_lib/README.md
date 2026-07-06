# `_lib/` — shared infra for the E4–E8 extension experiments

Reusable, **host-side**, API-free helpers so the `14.xx` extension experiments clone
rather than re-implement.

| Helper | Plan ref | Used by | What it does |
|---|---|---|---|
| `emit_metadata.sh <root>` | §2.1 | E4, E7 | sorted `path\tmtime\tmode\tsize\ttype` snapshot; run on the STAGE **before** chown/cp recovery so the metadata DV survives the round-trip |
| `neutral_names.sh <count> [seed]` | §2.3 | E4, E6 | opaque uniform 6-hex tokens, no readable ordinal; deterministic per seed |
| `lint_seed.sh <ws> [extra_re]` | §2.4 | all | launch gate: exits 1 if the seed carries an imperative/task-leak token or an ordinal filename scheme |
| `real_gated_attempt.py` | §2.2 | E7, E4 | `gated_attempt(run_dir, target, errno)` real kernel-refusal detector + `real_exec(run_dir, target)` real-invocation detector (tape `tool_result`, never narration) |
| `covariates.py` | §2.5 | E4, E6, E8 | activity floor + metadata-perception + exploration covariates; **recorded, never credited**; powers the symmetric "high-perception NULL = falsification" rule |

## Wiring conventions

- **Bash helpers** are sourced/called host-side from a `14.xx/run-container.sh` (the seed
  functions and recovery run on the host, not in the container) — no container mount needed.
- **Python helpers** are imported by a `14.xx/analysis/score.py` via
  `sys.path.insert(0, <repo>/public/structural-elicitation/experiments/_lib)`.
- The missionless prime / docker / env block in every `14.xx/run-container.sh` is copied
  **verbatim** from `13.24-nonlinguistic-gap/run-container.sh`; only the `seed_*()`
  functions, the `case "${CELL}"` dispatch, the metadata-capture step, and the scorer differ.

## Reflexive-capture discipline (carried by every card)

Pre-register DV thresholds; neutral naming (mktemp staging already strips the arm name from
`/proc/self/mountinfo`); blind-able scoring; **symmetric** decision rules (a high-perception
NULL falsifies, it does not excuse). Any perceived POSIX metadatum renders as readable text,
so the defensible claim is "no readable **ordinal/instruction token**", not "non-linguistic".
