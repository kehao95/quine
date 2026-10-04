# Subjective Reality Sources Index

This file records the current control-plane, implementation, coverage, and experiment surfaces relevant to the `subjective-reality` dossier.

It is an orientation index, not a claim of paper-readiness.

## Dossier And Program Control Plane

- dossier: [`../README.md`](../README.md)
- workspace: [`../WORKSPACE.md`](../WORKSPACE.md)
- manuscript plan: [`../manuscript-plan.md`](../manuscript-plan.md)
- paper index: [`../../../papers/README.md`](../../../papers/README.md)
- program map: [`../../../program-map.md`](../../../program-map.md)
- family registry: [`../../../paper-families.md`](../../../paper-families.md)
- evidence map: [`../../../evidence-map.md`](../../../evidence-map.md)

## Core Definitions And Boundary Docs

- implementation contract: [`../../../core/implementation.md`](../../../core/implementation.md)
- evolution line: [`../../../core/evolution/11-subjective-reality.md`](../../../core/evolution/11-subjective-reality.md)
- process-as-fs bridge: [`../../process-as-fs/README.md`](../../process-as-fs/README.md)
- fork/world decision: [`../../../core/decisions/2026-03-fork-worlds.md`](../../../core/decisions/2026-03-fork-worlds.md)
- runtime surface registry: [`../../../core/registries/runtime-surface.md`](../../../core/registries/runtime-surface.md)

## Runtime And Tooling Surfaces

- prompt disclosure of workspace/world semantics: [`../../../../internal/runtime/prompt.go`](../../../../internal/runtime/prompt.go)
- runtime dispatch and world switching: [`../../../../internal/runtime/runtime.go`](../../../../internal/runtime/runtime.go)
- workspace physics integration: [`../../../../internal/tools/workspace_physics.go`](../../../../internal/tools/workspace_physics.go)
- overlay subjective workspace backend: [`../../../../internal/tools/workspace_overlay_linux.go`](../../../../internal/tools/workspace_overlay_linux.go)
- switch-world tool surface: [`../../../../internal/tools/switch_world.go`](../../../../internal/tools/switch_world.go)
- tool schema surface: [`../../../../internal/tools/schemas.go`](../../../../internal/tools/schemas.go)

## Deterministic Tests

- prompt/runtime contract tests: [`../../../../internal/runtime/prompt_test.go`](../../../../internal/runtime/prompt_test.go)
- runtime identity/root projection tests:
  [`../../../../internal/runtime/agent_root_test.go`](../../../../internal/runtime/agent_root_test.go)
- workspace and switch-world tests: [`../../../../internal/tools/tools_test.go`](../../../../internal/tools/tools_test.go)
- fork world / parent-child lineage tests: [`../../../../internal/tools/fork_test.go`](../../../../internal/tools/fork_test.go)
- Linux workspace job attributes tests: [`../../../../internal/tools/job_attrs_linux_test.go`](../../../../internal/tools/job_attrs_linux_test.go)

## Runtime Coverage Ladder

- coverage inventory: [`../../../../tests/runtime/COVERAGE_MAP.md`](../../../../tests/runtime/COVERAGE_MAP.md)
- live runtime harness: [`../../../../tests/runtime/run.sh`](../../../../tests/runtime/run.sh)
- model evaluation registry: [`../../../../tests/model/evaluations.toml`](../../../../tests/model/evaluations.toml)
- model blueprint status: [`../../../../tests/model/BLUEPRINTS.md`](../../../../tests/model/BLUEPRINTS.md)

## Current Model-Facing Evaluation Families

- switch-world usage tasks: [`../../../../tests/model/README.md`](../../../../tests/model/README.md)
- switch-world explicit restore family: [`../../../../tests/model/l3-usage/switch-world`](../../../../tests/model/l3-usage/switch-world)
- workspace usage families: [`../../../../tests/model/l3-usage/workspace-overlay/relative-path-explicit`](../../../../tests/model/l3-usage/workspace-overlay/relative-path-explicit), [`../../../../tests/model/l3-usage/workspace-overlay/absolute-path-explicit`](../../../../tests/model/l3-usage/workspace-overlay/absolute-path-explicit)
- sandbox / emergent world-facing discovery families: [`../../../../tests/model/l4-discovery/sandbox`](../../../../tests/model/l4-discovery/sandbox)
- fork adopt / winning world promotion pilot: [`../../../../tests/model/pilot-l4-discovery/fork-adopt/winning-world-promotion`](../../../../tests/model/pilot-l4-discovery/fork-adopt/winning-world-promotion)

## Experiment Feeds Worth Revisiting

- local WAL recovery and reset-world precursor evidence: [`../../../../experiments/active/benchmarks/bench-03-local-wal-recovery/README.md`](../../../../experiments/active/benchmarks/bench-03-local-wal-recovery/README.md)
- world actuation pressure: [`../../../../experiments/archive/legacy-phase-6/6.8-world-actuation-pressure/README.md`](../../../../experiments/archive/legacy-phase-6/6.8-world-actuation-pressure/README.md)
- runtime-contract success path: [`../../../../experiments/archive/legacy-phase-6/6.11-runtime-contract-success-path/README.md`](../../../../experiments/archive/legacy-phase-6/6.11-runtime-contract-success-path/README.md)
- runtime-contract cue ablation: [`../../../../experiments/archive/legacy-phase-6/6.12-runtime-contract-cue-ablation/README.md`](../../../../experiments/archive/legacy-phase-6/6.12-runtime-contract-cue-ablation/README.md)
- runtime-physics absorption note: [`../../../../experiments/archive/legacy-phase-6/6.35-runtime-physics-absorption/README.md`](../../../../experiments/archive/legacy-phase-6/6.35-runtime-physics-absorption/README.md)

## Current Read

- This dossier already has real runtime substrate and real test coverage.
- The main missing piece is not raw material; it is a stable paper-facing evidence backbone and a narrower claim boundary.
- The strongest current implemented slice is `workspace as first subjective world slice`, not a general theory of capability/resource/world projection.
