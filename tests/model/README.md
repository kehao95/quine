# Model Evaluations

This subtree contains Quine's model-layer evaluation catalog.

The split is now explicit:

- `l3-usage` checks whether the model can use a runtime mechanism correctly
- `l4-discovery` checks whether the model can discover that mechanism without the prompt naming it
- `l5-necessity` is reserved for zero-hint pressure tasks where the mechanism is physically forced out

Use [`./tests/model/run.sh`](./run.sh) to execute evaluations and [`./scripts/check-model-evaluations.sh`](../../scripts/check-model-evaluations.sh) to audit registry, scorer coverage, and baseline hygiene.

For design-stage experiments that are not yet active in the registry, use
[`tests/model/BLUEPRINTS.md`](./BLUEPRINTS.md).

## Layout

| Path | Role |
|------|------|
| `tests/model/evaluations.toml` | canonical active evaluation registry |
| `tests/model/pilots.toml` | runnable pre-registry pilot registry |
| `tests/model/run.sh` | shared runner for all model-evaluation layers |
| `tests/model/BLUEPRINTS.md` | planning surface for non-registered experiment blueprints |
| `tests/model/l3-usage/<feature>/<variant>/prompt.md` | canonical L3 prompt location |
| `tests/model/l4-discovery/<feature>/<variant>/prompt.md` | canonical L4 prompt location |
| `tests/model/l5-necessity/<feature>/<variant>/prompt.md` | canonical L5 prompt location |
| `tests/model/pilot-<level>-<mode>/<feature>/<variant>/prompt.md` | runnable pre-registry pilot prompt location |
| `tests/model/<level>/<feature>/<variant>/runs/` | preserved baselines for that evaluation |

## Usage

```bash
go build -o /tmp/quine ./cmd/quine/
source profiles/gpt-5.4-codex-oauth.env

./tests/model/run.sh usage
./tests/model/run.sh discovery
./tests/model/run.sh necessity
./tests/model/run.sh stdin-explicit-handoff
./tests/model/run.sh pilot:exec-final-utility-stream-handoff

./scripts/check-model-evaluations.sh --strict
./scripts/check-model-evaluations.sh --strict --require-baselines
./scripts/check-model-evaluations.sh --strict --require-tracked-baselines
./scripts/check-model-evaluations.sh --classify-residue
QUINE_BEHAVIOR_KEEP_FAILED_RUNS=0 ./scripts/check-model-evaluations.sh --prune-run-tree
```

Default to `profiles/gpt-5.4-codex-oauth.env`. Use `profiles/kimi-k2.7-code-litellm-medium.env` only when you
are explicitly checking that provider lane.

`./tests/model/run.sh necessity` is the canonical active L5 entrypoint.

`./tests/model/run.sh pilot:<id>` runs a pre-registry pilot from
[`tests/model/pilots.toml`](./pilots.toml). Use that path once a blueprint has
become runnable enough for live iteration, but before it is honest enough for
the active acceptance catalog.

Pilot runs retain non-passing run directories by default so discovery failure
shapes remain inspectable during iteration. Do not treat pilot `runs/` as
ambient cache or let pruning scripts silently delete the methodological record.

Active registered evaluations now follow the same default: non-passing runs are
retained unless you explicitly set `QUINE_BEHAVIOR_KEEP_FAILED_RUNS=0`. Treat
failed `runs/` directories as audit evidence first and local residue second.

Residue classification:

- use `./scripts/check-model-evaluations.sh --classify-residue` during closure or maintenance to list untracked passing-baseline candidates, failed-methodological candidates, and incomplete run artifacts
- use `--require-tracked-baselines` when a closure boundary must prove that baseline evidence is committed rather than merely present in the local filesystem
- `--prune-run-tree` is destructive and requires `QUINE_BEHAVIOR_KEEP_FAILED_RUNS=0`; use it only after run-wave triage, never as an implicit cleanup step

Historical baseline interpretation:

- some preserved baselines under `runs/` predate the 2026-04-18 context-first cognition refactor
- those historical artifacts may still mention `log/current.jsonl` as the current tape or describe `exec` as resetting context
- treat those preserved runs as methodological record, not as the current runtime-contract documentation

Budget and timeout policy:

- low `max_turns` and short wall-clock timeouts are explicit pressure variables, not harmless defaults
- unless an evaluation is directly studying budget or timeout behavior, default runner turn budgets should stay at or above `30`
- unless an evaluation is directly studying timeout behavior, default evaluation timeouts should stay at or above `10 min`; the shared runner now defaults to `30 min`
- if a run fails because a single action was malformed, a relative path was wrong, or a one-shot helper contract exited early, do not misclassify that as a budget or timeout failure

Do not register a blueprint into [`tests/model/evaluations.toml`](./evaluations.toml)
until the prompt leakage boundary, scorer shape, and preservation plan are all
concrete enough to support a real pilot run.

Registry metadata lives in [`tests/model/evaluations.toml`](./evaluations.toml). Each evaluation now carries:

- `level` - `l3`, `l4`, or `l5`
- `mode` - `usage`, `discovery`, or `necessity`
- `feature` - normalized runtime feature name
- `variant` - evaluation variant under that feature
- `path` - canonical filesystem path prefix: `<level>-<mode>/<feature>/<variant>`

Pilot metadata lives in [`tests/model/pilots.toml`](./pilots.toml) and uses the
same fields, but with the path prefix `pilot-<level>-<mode>/...`.

Naming contract:

- feature slugs name the runtime surface, not the prompt brand
- variant slugs distinguish task shape within one feature
- `id` is the canonical selector and should follow `<feature>-<variant>`
- prompt and baseline paths are derived from registry metadata rather than inferred from old flat directories

Blueprint contract:

- use `tests/model/BLUEPRINTS.md` plus `tests/model/blueprints/*.md` for design-stage work
- keep future ids and future canonical paths explicit even before registration
- once a blueprint becomes runnable but is not registerable, place it in `tests/model/pilots.toml`
- do not add a pilot to the active registry until it is intended to serve as acceptance

## L3 Usage Evaluations

This layer contains explicit mechanism-use evaluations, including cases that were previously misclassified as discovery but are now classified honestly as usage.

| Feature | Variants |
|---------|----------|
| `detach` | `explicit-protocol` |
| `interactive-jobs` | `repl-surface`, `overlay-world-adoption`, `terminal-control-surface` |
| `daemon-pattern` | `http-server-survival` |
| `stdin` | `explicit-handoff`, `shell-hostile-literals` |
| `stdin-physics` | `material-stream-contract` |
| `shell-isolation` | `explicit-nonpersistence` |
| `shell-envelope` | `timeout-resume-exit-observation`, `timeout-terminate-exit-observation` |
| `fragments` | `explicit-surface-inspection`, `explicit-durable-surface-selection` |
| `context-memory` | `explicit-refresh` |
| `agents-md` | `explicit-policy-activation`, `explicit-refresh-without-exec` |
| `skills` | `explicit-catalog-activation`, `hierarchical-resource-use`, `explicit-frontmatter-refresh` |
| `session-resume` | `explicit-contract` |
| `process-surface` | `explicit-self-neighbor-discovery`, `explicit-self-source-inspection`, `explicit-peer-communication`, `explicit-peer-inject-delivery`, `explicit-peer-interrupt-delivery`, `explicit-peer-failover`, `explicit-client-response` |
| `relation-recovery` | `explicit-outcome-semantics` |
| `resource-governance` | `explicit-depth-limit`, `explicit-agent-slot-limit` |
| `exec` | `explicit-external-handoff`, `explicit-self-source-rebuild-handoff`, `stream-pipe-handoff`, `stream-stdio-handoff` |
| `vision` | `direct-tool-use` |
| `fork` | `explicit-modes`, `explicit-relation-surfaces`, `search-delegation`, `parallel-race`, `batch-parallelism` |
| `spawn` | `explicit-fresh-process`, `fresh-audit-shared-workspace` |
| `fork-world` | `explicit-world-selection`, `search-lane-scoping`, `batch-lane-scoping` |
| `fork-adopt` | `explicit-child-adoption`, `winner-adoption` |
| `execution-budget` | `hard-fail`, `hard-fail-thermodynamic` |
| `idle` | `explicit-suspension-resume` |
| `anchor-memory` | `explicit-mark-unfold`, `retrieval-pressure-explicit` |
| `sandbox` | `explicit-isolation` |
| `switch-world` | `explicit-revision-restore`, `destructive-probe-restore` |
| `workspace-overlay` | `relative-path-explicit`, `absolute-path-explicit`, `exploratory-decode` |
| `workspace-direct` | `relative-path-explicit` |

## L4 Discovery Evaluations

These are the currently active discovery evaluations that still clear the stronger boundary: the prompt does not name the target mechanism, but the environment or task still makes it discoverable.

| Feature | Variants |
|---------|----------|
| `detach` | `two-lane-overlap`, `forced-overlap-window` |
| `interactive-jobs` | `terminal-world-discovery`, `terminal-control-discovery` |
| `fragments` | `active-surface-discovery` |
| `agents-md` | `relevant-policy-discovery`, `durable-rule-discovery` |
| `context-memory` | `inheritance-discovery` |
| `shell-envelope` | `timeout-resume-discovery`, `timeout-terminate-discovery` |
| `spawn` | `fresh-reviewer-discovery`, `fresh-audit-choice` |
| `skills` | `relevant-workflow-discovery`, `durable-workflow-discovery` |
| `session-resume` | `runtime-discovery` |
| `process-surface` | `runtime-root-discovery`, `peer-message-discovery`, `peer-inject-discovery`, `peer-interrupt-discovery`, `peer-failover-discovery`, `client-response-discovery` |
| `relation-recovery` | `status-discovery` |
| `fork` | `context-preserving-choice` |
| `idle` | `external-poke-discovery` |
| `sandbox` | `unknown-format-boldness` |
| `containment` | `hostile-script-survival` |
| `workspace-direct` | `peer-handoff-observation` |
| `workspace-overlay` | `fuse-dangerous-decoder-containment` |

## L5 Necessity Evaluations

`l5-necessity` is now a first-class layer in the control plane. Current catalog:

| Feature | Variants | Status |
|---------|----------|--------|
| `idle` | `quiet-standby-pressure` | active baseline recorded |
| `process-surface` | `peer-callback-protocol` | active experimental, baseline pending |
| `relation-recovery` | `resume-error-semantics` | active experimental, baseline recorded |
| `session-resume` | `corpse-resurrection` | active baseline recorded |

L5 requirements:

- zero mechanism hinting
- zero behavior or strategy hinting
- real task/environment pressure rather than rhetorical narrowing
- a scorer that can honestly say the mechanism was forced rather than merely chosen

L5 negative-constraint discipline:

- allowed: close accidental side channels or non-target shortcuts that the
  fixture/scorer cannot otherwise isolate cleanly, for example helper-private
  artifacts, retained logs, or unrelated host process discovery
- not allowed: steer the model away from alternate strategy families that could
  still satisfy the task, for example telling it not to invent a new bootstrap,
  wrapper, alias, sidecar startup path, or other competing persistence surface
- if a line narrows *which solution family* is acceptable rather than sealing a
  leakage route, that line belongs in fixture/scorer design, not in an L5 prompt

Current process-surface callback iteration also retains a cleanroom pilot wave
under `pilot:process-surface-peer-callback-cleanroom`. Keep those pilot runs as
methodological record while the active `peer-callback-protocol` registry entry
still lacks its canonical baseline.
Retained pilot runs keep artifacts, not live `/tmp/quine` processes.

The new `pilot:agents-md-startup-token-lockin` case is the first necessity-side
startup-guidance probe for the `context/prompt` / `AGENTS.md` surface. Keep it
pilot-only until repeated runs show that the scorer is really measuring
startup-surface lock-in rather than incidental prompt imitation.

`pilot:context-memory-exec-token-lockin` is the corresponding inherited-context
probe for `context/prompt/30-memory.md`. Its current pilot bar verifies physical retention
across exec and rejects transient argv leakage; do not promote it until the
successor-output behavior is stable rather than merely reaching Memory.

## Blueprint Portfolio

The next-wave design portfolio now lives outside the active registry on purpose:

- [`tests/model/blueprints/l3-usage.md`](./blueprints/l3-usage.md)
- [`tests/model/blueprints/l4-discovery.md`](./blueprints/l4-discovery.md)
- [`tests/model/blueprints/l5-necessity.md`](./blueprints/l5-necessity.md)

## Design Rules

- Register every evaluation in [`tests/model/evaluations.toml`](./evaluations.toml).
- Keep the canonical path explicit through `level`, `mode`, `feature`, `variant`, and `path`.
- Treat old flat prompt paths as compatibility aliases only.
- Add scorer logic in [`tests/model/run.sh`](./run.sh) before calling an evaluation done.
- Preserve at least one complete passing baseline run for every active evaluation.
- Closure for model-evaluation work must classify `tests/model/.../runs/...` diffs explicitly.
- If an active evaluation gains or refreshes its canonical baseline, commit that retained run in the same closure interval rather than leaving it as local residue.
- Keep pilot failure runs when they remain the methodological record; only delete them deliberately.
- Prefer `usage` for mechanism operation checks.
- Prefer `discovery` only when the prompt stays clear of mechanism naming and strategy teaching.
- Use `necessity` only when the task pressure, not the wording, does the forcing.
- Keep non-running experiment design in the blueprint surface until it is honest enough to pilot.
