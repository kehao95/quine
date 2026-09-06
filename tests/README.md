# Evaluation Harnesses

This tree is part of the public research surface. Git carries the harnesses,
prompts, expected results, DVC pointers, and documentation; retained run data is
restored only through the credential-scanned public DVC manifest. Checked-out
run directories, local links, and ignored residue never enter the projection.

The deterministic substrate and public Go packages run in a fresh public clone.
Live runtime and model evaluations may additionally require Linux or Lima, a
provider credential, private profiles, or research-control files that are not
in the current snapshot.

Quine now uses a five-layer evaluation ladder:

1. `substrate` - deterministic Go tests
2. `runtime` - live binary contract tests
3. `usage` - model-layer mechanism use
4. `discovery` - unnamed mechanism discovery
5. `necessity` - zero-hint pressure selection

Use [`./tests/validate.sh`](./validate.sh) as the runtime evaluation entrypoint. Control-plane changes use
`make check-control-plane`; `make test-control-plane-routing` exercises actual
hook domains and active-experiment frontmatter acceptance/rejection.

The ladder is the acceptance-side companion to Quine's
[operating philosophy](../Paper/philosophy/operating-philosophy.md): it checks
whether implementation-backed runtime physics are usable or discoverable by a
model, not just present in code.

Drift check: `./scripts/check-tests-entry-docs.sh --strict tests/README.md` and `./scripts/check-validation-surface.sh --strict`.

## Layout

| Path | Role |
|------|------|
| `tests/validate.sh` | unified evaluation dispatcher across all 5 layers |
| [`tests/runtime/README.md`](./runtime/README.md) | runtime-harness index and local control plane |
| `tests/runtime/` | live runtime contract harness and Linux bridge helpers |
| `tests/runtime/COVERAGE_MAP.md` | runtime feature inventory mapped onto the evaluation ladder |
| `tests/model/` | model-facing evaluation registry, runner, docs, and baselines |
| [`tests/test_control_plane_routing.py`](./test_control_plane_routing.py) | local routing regression; no model calls |
| `tests/fixtures/` | dispatcher smoke for `scripts/test-prediction.sh` |

## Canonical commands

```bash
./tests/validate.sh --change substrate
./tests/validate.sh --change runtime --runtime test_fd4_delivery
./tests/validate.sh --change usage --runtime test_fd4_delivery --usage stdin-explicit-handoff
./tests/validate.sh --change discovery --runtime test_workspace_overlay_commit --usage workspace-overlay-relative-path-explicit --discovery sandbox-unknown-format-boldness
```

`--change` selects the highest required layer; the script always prints the full ladder so higher layers remain visible even when they are not required for the current change.

## Audit And Retention

- Failed `runtime` and active `model` runs are retained by default because non-passing traces are audit evidence, not disposable cache.
- Set `QUINE_BEHAVIOR_KEEP_FAILED_RUNS=0` only when you deliberately want local pruning; set `QUINE_BEHAVIOR_KEEP_ALL_RUNS=1` to retain every run.
- For runtime audit, treat `QUINE_RETENTION_DIR/sessions/<session>/tapes/*.jsonl` as the canonical tape surface when `QUINE_RETENTION_DIR` is set; otherwise use `QUINE_DATA_DIR/log/<session>/tapes/*.jsonl`. `QUINE_DATA_DIR/tapes/<session>/*.jsonl` is only a compatibility alias when `QUINE_RETENTION_DIR` is set.
- Treat `control.jsonl` as a separate control-event audit stream, not as a tape; helpers must not glob arbitrary `*.jsonl` under the runtime root and assume they are all equivalent.
