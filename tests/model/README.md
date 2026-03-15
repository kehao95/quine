# Model Scenarios

This subtree contains Quine's model-layer tests.

The split is intentional:

- `instructional` checks whether the model can follow the explicit protocol correctly
- `emergent` checks whether the environment itself causes the model to discover the intended primitive

Use [`./tests/model/run.sh`](./run.sh) to execute scenarios and [`./scripts/check-model-scenarios.sh`](../../scripts/check-model-scenarios.sh) to audit registry, scorer coverage, and baseline hygiene.

## Layout

| Path | Role |
|------|------|
| `tests/model/scenarios.toml` | canonical scenario registry |
| `tests/model/run.sh` | shared runner for both layers |
| `tests/model/instructional/prompts/` | prompts for Layer 3 scenarios |
| `tests/model/instructional/runs/` | preserved passing baselines for Layer 3 |
| `tests/model/emergent/prompts/` | prompts for Layer 4 scenarios |
| `tests/model/emergent/runs/` | preserved passing baselines for Layer 4 |

## Usage

```bash
go build -o /tmp/quine ./cmd/quine/
source .env.kimi-oauth

./tests/model/run.sh instructional
./tests/model/run.sh emergent
./tests/model/run.sh stdin
./tests/model/run.sh stdin-emergent

./scripts/check-model-scenarios.sh --strict
./scripts/check-model-scenarios.sh --strict --require-baselines
./scripts/check-model-scenarios.sh --prune-run-tree
```

Default model-scenario env is `.env.kimi-oauth`, including `linux-kimi`
scenarios. If a scenario crosses a guest or `sudo` boundary, stage the OAuth
config through `QUINE_CONFIG_DIR` or the harness-local equivalent. Treat any
scenario that still requires a plain API-key env as a harness defect to repair
rather than a model-scenario exception.

By default the runner only keeps complete passing runs. Failed or incomplete runs are pruned unless `QUINE_BEHAVIOR_KEEP_FAILED_RUNS=1` or `QUINE_BEHAVIOR_KEEP_ALL_RUNS=1` is set.

Scenario control-plane metadata lives in [`tests/model/scenarios.toml`](./scenarios.toml). When a scenario needs a non-default execution budget or exhaustion policy, encode it there rather than hiding it in runner-local shell logic. Only budget-focused scenarios should use aggressively tight turn limits; other behavior scenarios should inherit the default generous budget unless there is a clear reason not to.

## Instructional Scenarios

| ID | Feature | Notes |
|----|---------|-------|
| `detach` | detached jobs | explicit detach protocol |
| `interactive` | interactive jobs | PTY-backed REPL surface |
| `daemon` | daemon pattern | detach plus server survival |
| `stdin` | stdin parameter | explicit stdin hint |
| `stdin-physics` | stdin physics | explicit fd 3 / stdin protocol |
| `vision` | vision | direct vision tool usage |
| `escalate` | escalate | explicit handoff instruction |
| `swarm-fork` | fork | explicit fork modes |
| `budget-hard-fail` | execution budget | hard-fail planning |
| `budget-near-death` | execution budget | exec-only continuation |
| `budget-hard-fail-thermo` | execution budget | metaphor overlay contract |
| `anchor-memory` | anchor memory | explicit checkpoint protocol |
| `sandbox` | sandbox | explicit isolation physics |
| `restore-world` | restore world revision | explicit intra-run rewind protocol |
| `workspace-shadow` | workspace overlay | relative-path transactional writes |
| `workspace-absolute` | workspace overlay | absolute-path transactional writes |

Instructional scenarios are Layer 3. They are the right place for wiring checks and explicit protocol regressions.

## Emergent Scenarios

| ID | Feature | Notes |
|----|---------|-------|
| `detach-emergent` | detached jobs | overlap slow and fast lanes without coaching |
| `detach-overlap-emergent` | detached jobs | forced overlap acceptance |
| `stdin-emergent` | stdin parameter | discover stdin instead of quoting tricks |
| `escalate-emergent` | escalate | discover escalation or impossibility |
| `fork-search` | fork | discover delegation for search |
| `fork-race` | fork | discover parallel strategy race |
| `fork-batch` | fork | discover batch parallelism |
| `sandbox-emergent` | sandbox | exploit isolation boldly under vague task |
| `workspace-shadow-emergent` | workspace overlay | exploratory decoding in transactional workspace |
| `restore-world-emergent` | restore world revision | destructive first probe, restore `wr0`, helper-driven recovery |
| `logic-bomb` | logic-bomb containment | identify danger, contain execution, survive |

Emergent scenarios are Layer 4. They are the highest acceptance layer for model-facing physics.

## Design Rules

- Register every scenario in [`tests/model/scenarios.toml`](./scenarios.toml).
- Put prompts under the correct layer directory; do not encode layer only through filename suffix.
- Add scorer logic in [`tests/model/run.sh`](./run.sh) before calling a scenario done.
- Preserve at least one complete passing baseline run for every current scenario.
- Prefer an instructional scenario for protocol regressions.
- Prefer an emergent scenario for discoverability and physics claims.
