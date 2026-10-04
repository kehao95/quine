# Revisable Time Sources Index

This file records the current control-plane, implementation, coverage, and experiment surfaces relevant to the `revisable-time` dossier.

It is an orientation index, not a claim that the paper already has an implementation-backed spine.

## Dossier And Program Control Plane

- dossier: [`../README.md`](../README.md)
- workspace: [`../WORKSPACE.md`](../WORKSPACE.md)
- manuscript plan: [`../manuscript-plan.md`](../manuscript-plan.md)
- paper index: [`../../../papers/README.md`](../../../papers/README.md)
- program map: [`../../../program-map.md`](../../../program-map.md)
- family registry: [`../../../paper-families.md`](../../../paper-families.md)
- evidence map: [`../../../evidence-map.md`](../../../evidence-map.md)

## Core Definitions And Boundary Docs

- dossier thesis boundary: [`../README.md`](../README.md)
- manuscript thesis structure: [`../manuscript-plan.md`](../manuscript-plan.md)
- evolution line: [`../../../core/evolution/13-revisable-time.md`](../../../core/evolution/13-revisable-time.md)
- implementation contract it currently depends on: [`../../../core/implementation.md`](../../../core/implementation.md)
- neighboring world-semantics paper: [`../../subjective-reality/README.md`](../../subjective-reality/README.md)
- process-as-fs bridge: [`../../process-as-fs/README.md`](../../process-as-fs/README.md)

## Current Runtime Precursors

- overlay workspace rollback substrate: [`../../../../internal/tools/workspace_overlay_linux.go`](../../../../internal/tools/workspace_overlay_linux.go)
- workspace/world-state integration: [`../../../../internal/tools/workspace_physics.go`](../../../../internal/tools/workspace_physics.go)
- switch-world restore surface: [`../../../../internal/tools/switch_world.go`](../../../../internal/tools/switch_world.go)
- runtime dispatch for world switching: [`../../../../internal/runtime/runtime.go`](../../../../internal/runtime/runtime.go)
- prompt disclosure of revision handles and restore semantics: [`../../../../internal/runtime/prompt.go`](../../../../internal/runtime/prompt.go)
- tool schema wording for restore mode: [`../../../../internal/tools/schemas.go`](../../../../internal/tools/schemas.go)

## Deterministic Tests For Current Precursors

- workspace and switch-world unit tests: [`../../../../internal/tools/tools_test.go`](../../../../internal/tools/tools_test.go)
- fork parent/child world revision behavior: [`../../../../internal/tools/fork_test.go`](../../../../internal/tools/fork_test.go)
- prompt/runtime contract tests: [`../../../../internal/runtime/prompt_test.go`](../../../../internal/runtime/prompt_test.go)
- live runtime rollback/restore harness: [`../../../../tests/runtime/run.sh`](../../../../tests/runtime/run.sh)
- runtime coverage inventory: [`../../../../tests/runtime/COVERAGE_MAP.md`](../../../../tests/runtime/COVERAGE_MAP.md)

## Current Model-Facing Evaluation Families

- switch-world usage family: [`../../../../tests/model/l3-usage/switch-world`](../../../../tests/model/l3-usage/switch-world)
- switch-world pilot backlog and variants: [`../../../../tests/model/BLUEPRINTS.md`](../../../../tests/model/BLUEPRINTS.md)
- model evaluation registry entries: [`../../../../tests/model/evaluations.toml`](../../../../tests/model/evaluations.toml), [`../../../../tests/model/pilots.toml`](../../../../tests/model/pilots.toml)

## Experiments With Relevant Time / Rollback Pressure

- local WAL recovery and reset-world ad hoc replay: [`../../../../experiments/active/benchmarks/bench-03-local-wal-recovery/README.md`](../../../../experiments/active/benchmarks/bench-03-local-wal-recovery/README.md)
- historical track notes mentioning rollback direction: [`../../../archive/track-system-meta/CHANGES-2026-03-06.md`](../../../archive/track-system-meta/CHANGES-2026-03-06.md), [`../../../archive/track-system-meta/CHANGELOG-outline-track-a.md`](../../../archive/track-system-meta/CHANGELOG-outline-track-a.md)

## Current Read

- This dossier is real at the control-plane level, but most of its runtime substance is still borrowed from workspace/world revision precursors.
- The currently implemented substrate supports `minimal provisional time`: provisional view, accepted base, and revision restore.
- The repo does not yet expose a strong independent time-semantics surface such as checkpoints, retract, retained alternate histories, or cognition/world separation as a stable primitive set.
