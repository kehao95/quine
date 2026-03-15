# Testing Quine

This is the canonical testing control plane for Quine.

The testing stack is part of Quine's public engineering surface. `main` should
carry these commands, harnesses, and scenario docs by default unless a specific
component is genuinely private or unstable.

Quine now has four validation layers:

1. `substrate` - deterministic Go correctness
2. `runtime` - live binary contract against a real model
3. `instructional` - model behavior when the protocol is explicit
4. `emergent` - highest-level acceptance for discoverability and physical understanding

Use [`./tests/validate.sh`](./tests/validate.sh) as the default entrypoint after any change.

## Core Rule

Do not hand off a change after running only `go test ./...` unless the change is honestly `substrate`-only.

The point of the four-layer ladder is that higher layers stay visible even when they are not required for the current change. `./tests/validate.sh` always prints all four layers before it runs anything.

## Canonical Commands

```bash
./tests/validate.sh --change substrate
./tests/validate.sh --change runtime --runtime test_exit_success
./tests/validate.sh --change instructional --runtime test_fd4_delivery --instructional stdin
./tests/validate.sh --change emergent --runtime test_workspace_overlay_commit --instructional workspace-shadow --emergent workspace-shadow-emergent
```

If you want the script to infer the minimum layer from the current git worktree, use:

```bash
./tests/validate.sh --from-git --runtime test_fd4_delivery --instructional stdin
```

`--from-git` only infers the minimum required layer. It does not guess the exact runtime tests or scenarios for you.

## Layer Intent

### Layer 1: `substrate`

Command:

```bash
go test ./...
```

Questions answered:

- Did deterministic internal logic break?
- Did parsing, tape persistence, config loading, or tool argument handling regress?

### Layer 2: `runtime`

Command:

```bash
./tests/runtime/run.sh
./tests/runtime/run.sh test_exit_success
```

Questions answered:

- Does the real binary still behave correctly under live execution?
- Are exit codes, stdout/stderr separation, fd 3/fd 4 semantics, tape creation, and workspace contracts still intact?

### Layer 3: `instructional`

Command:

```bash
./tests/model/run.sh instructional
./tests/model/run.sh stdin
```

Questions answered:

- Can the model use the protocol correctly when the prompt names the intended primitive?
- Did we break tool legibility, explicit protocol compliance, or documented workflow wiring?

### Layer 4: `emergent`

Command:

```bash
./tests/model/run.sh emergent
./tests/model/run.sh stdin-emergent
```

Questions answered:

- Does the environment itself nudge the model toward the intended primitive?
- Can the model discover and use the physics without direct coaching?

This is the highest acceptance layer for model-facing physics.

## Validation Contract

| Change layer | Minimum validation |
|--------------|--------------------|
| `substrate` | `./tests/validate.sh --change substrate` |
| `runtime` | `./tests/validate.sh --change runtime --runtime <test>` |
| `instructional` | `./tests/validate.sh --change instructional --runtime <test> --instructional <scenario>` |
| `emergent` | `./tests/validate.sh --change emergent --runtime <test> --instructional <scenario> --emergent <scenario>` |

Interpretation rules:

- `instructional` is where you prove explicit protocol correctness.
- `emergent` is where you prove discoverability and physical understanding.
- A feature is not accepted as runtime physics on model behavior evidence alone if only the instructional layer passes.

## Environment Profiles

| Profile | Use for | Notes |
|---------|---------|-------|
| `mac-kimi` | most substrate, runtime, and non-Linux model checks | default development lane; canonical env is `.env.kimi-oauth` |
| `linux-kimi` | sandbox, workspace, and other Linux-only physics | on macOS, runtime workspace checks bridge through Lima `colima` by default; canonical env is `.env.kimi-oauth`, with OAuth config forwarded into guest/root runs via `QUINE_CONFIG_DIR` when needed |

## Model Scenario Control Plane

Model scenarios now live under [`tests/model/`](./tests/model/README.md).

If you are unsure what runtime surface a change touched, start with
[`tests/runtime/COVERAGE_MAP.md`](./tests/runtime/COVERAGE_MAP.md).
That inventory is the canonical map of runtime features vs layer coverage.

Canonical inventory:

- scenario registry: [`tests/model/scenarios.toml`](./tests/model/scenarios.toml)
- runner: [`tests/model/run.sh`](./tests/model/run.sh)
- audit: `./scripts/check-model-scenarios.sh`

Audit commands:

```bash
./scripts/check-model-scenarios.sh --strict
./scripts/check-model-scenarios.sh --strict --require-baselines
./scripts/check-model-scenarios.sh --prune-run-tree
```

The registry is now the source of truth. Prompts, scorer coverage, and baseline run trees must agree with it.

## Runtime Notes

For Linux-only workspace validation:

- `./tests/runtime/run.sh` is the runtime contract harness
- on macOS, workspace tests bridge into Lima via the `colima` guest
- workspace overlay validation still prefers rootless user-namespace mounts
- if the Linux run happens inside a guest, the binary used inside that guest must also be Linux-compatible

For OAuth-backed runs, ensure token cache is available through `QUINE_CONFIG_DIR`.
If a Linux guest or `sudo` path cannot yet carry that config through, treat that
as a harness defect to repair rather than switching the shared default away from
OAuth.

## Runtime Coverage Inventory

[`tests/runtime/COVERAGE_MAP.md`](./tests/runtime/COVERAGE_MAP.md) is the planning surface for runtime coverage.

When runtime implementation grows, update that map alongside code and tests so the new feature enters the four-layer ladder explicitly.

## Design Rules

Add or update a Go test when:

- the behavior is deterministic
- no real API call is required
- the claim is about substrate correctness

Add or update a runtime test when:

- the binary boundary matters
- the observable contract matters more than internals
- the claim is about live runtime behavior

Add or update an instructional scenario when:

- the feature needs an explicit protocol regression check
- you need a narrow diagnostic scenario that says "use primitive X"

Add or update an emergent scenario when:

- the real claim is discoverability
- the feature is being accepted as runtime physics
- you want to know whether the environment rewards the intended primitive without direct coaching

## Operational References

- [`tests/README.md`](./tests/README.md) - 4-layer test map
- [`tests/runtime/run.sh`](./tests/runtime/run.sh) - API-backed runtime harness
- [`tests/runtime/COVERAGE_MAP.md`](./tests/runtime/COVERAGE_MAP.md) - runtime feature inventory and coverage backlog
- [`tests/model/README.md`](./tests/model/README.md) - instructional and emergent scenario catalog
- [`tests/model/scenarios.toml`](./tests/model/scenarios.toml) - scenario registry
