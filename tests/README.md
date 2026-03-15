# Tests

Quine now has a four-layer validation ladder:

1. `substrate` — deterministic Go tests
2. `runtime` — live binary contract tests
3. `instructional` — model behavior with explicit protocol guidance
4. `emergent` — highest-level acceptance for physical understanding and discoverability

Use [`./tests/validate.sh`](./validate.sh) as the canonical entrypoint after any change.

## Layout

| Path | Role |
|------|------|
| `tests/validate.sh` | unified validation dispatcher across all 4 layers |
| `tests/runtime/` | live runtime contract harness and Linux bridge helpers |
| `tests/runtime/COVERAGE_MAP.md` | runtime feature inventory mapped onto the 4 validation layers |
| `tests/model/` | model-facing scenario registry, runner, docs, and baselines |

## Canonical commands

```bash
./tests/validate.sh --change substrate
./tests/validate.sh --change runtime --runtime test_fd4_delivery
./tests/validate.sh --change instructional --runtime test_fd4_delivery --instructional stdin
./tests/validate.sh --change emergent --runtime test_workspace_overlay_commit --instructional workspace-shadow --emergent workspace-shadow-emergent
```

`--change` selects the highest required layer; the script always prints all four layers so they cannot disappear from view.
