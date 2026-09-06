# Runtime Harness

This subtree holds Quine's live runtime contract surface.

Use it when the claim is about the real binary's observable behavior rather than deterministic internal correctness.

`./tests/runtime/run.sh` is the execution harness.
[`./tests/runtime/COVERAGE_MAP.md`](./COVERAGE_MAP.md) is the planning and backlog surface that maps runtime features onto the 5-layer evaluation ladder.

Drift check: `./scripts/check-tests-entry-docs.sh --strict tests/runtime/README.md` and `./scripts/check-runtime-doc-sync.sh --strict tests/runtime/COVERAGE_MAP.md`, plus comparing `test_*` / `gate_*` functions in `run.sh` against `COVERAGE_MAP.md`.

## Layout

| Path | Role |
|------|------|
| `tests/runtime/run.sh` | API-backed runtime contract harness |
| `tests/runtime/COVERAGE_MAP.md` | runtime feature inventory and coverage backlog |
| `tests/runtime/lib/` | Linux/Lima bridge helpers used by the runtime harness |

## Usage

```bash
go build -o /tmp/quine ./cmd/quine/
source profiles/gpt-5.4-codex-oauth.env

./tests/runtime/run.sh
./tests/runtime/run.sh test_exit_success
./tests/runtime/run.sh test_fd4_delivery test_exec_preserves_mission
./tests/runtime/run.sh gate_overlay_finalization_baseline
```

Default to `profiles/gpt-5.4-codex-oauth.env`. Use `profiles/kimi-k2.7-code-litellm-medium.env` only when you
are explicitly checking that provider lane.

`gate_overlay_finalization_baseline` is the named kernel-overlay baseline for
the current hard requirements: commit, rollback, `switch_world` restore/branch
semantics, child-world adoption, long-lineage materialization, durable shutdown
phase order, and commit-intent recovery. Keep `test_fork_race_adopt_winner` as
a separate subjective-race probe when you specifically want `adopt_winner`
coverage.

Runtime tests are Layer 2 of the validation ladder. They check the real binary's observable contract: exit behavior, channel separation, tape outputs, workspace/world behavior, and runtime-tool continuity.

Retained fork/spawn relation `result.json` and `status.json` files are contract surfaces. Keep `spawned`, `completed`, `succeeded`, and `killed` explicit even when the value is zero.

## Audit Surfaces

Failed runtime runs are retained by default under `/tmp/quine-e2e.*` so harness,
provider, and runtime failures remain inspectable after exit. Set
`QUINE_BEHAVIOR_KEEP_FAILED_RUNS=0` only when you intentionally want pruning.

Inside a preserved runtime run:

- `stdout` and `stderr` are the top-level execution outcome surface
- `tapes/log/<session>/tapes/*.jsonl` is the canonical session tape surface when `QUINE_RETENTION_DIR` is unset
- `tapes/log/sessions/<session>/tapes/*.jsonl` is the canonical session tape surface when `QUINE_RETENTION_DIR=tapes/log`
- `tapes/tapes/<session>/*.jsonl` is only a compatibility alias when `QUINE_RETENTION_DIR` is set
- `tapes/log/<session>/control.jsonl` or `tapes/log/sessions/<session>/control.jsonl` is the retained control-event audit log, depending on whether `QUINE_RETENTION_DIR` is unset or set
- `tapes/log/<session>/runtime.log` or `tapes/log/sessions/<session>/runtime.log` is the retained operational log, depending on whether `QUINE_RETENTION_DIR` is unset or set

Do not glob arbitrary `*.jsonl` under `tapes/` and assume they are all tapes.
`control.jsonl` and tape JSONL have different semantics and must stay separate
in harness helpers and manual audits.

## Coverage Planning

Do not treat the harness file itself as the whole runtime test surface.

Interpretation rules:

- the harness is the executable contract surface
- the coverage map is the planning surface that keeps higher-layer gaps visible
- if a new runtime feature lands, update the coverage map in the same state transition so the feature enters the ladder intentionally
- if runtime behavior changes enough to affect harness semantics, update this README and [`EVALUATION.md`](../../EVALUATION.md) alongside the code
