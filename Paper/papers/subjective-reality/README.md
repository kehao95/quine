---
paper_id: subjective-reality
claim_bundle: subjective-reality
venue: hotos
format-envelope: workshop-short
track: system
cluster: process-native-systems-substrate
status: planning
state_mode: planning-skeleton
role: world-semantics
family: runtime-substrate
theme: runtime-ontology
theme_home: ../../../theory/views/substrate.md
publication_class: archival-candidate
maturity_target: mechanism
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [future-task-surface, exec-barrier]
gaps: [revisable-time-semantics]
evidence: [none-yet]
view-affinity: [substrate]
---

# Subjective Reality

## Current Question

What makes an agent's operative world a projected local reality rather than a single objective environment, and which current Quine surfaces already support that narrower claim?

## Thesis

Agent runtimes construct operative worlds by projecting tools, resources, and interfaces into a bounded local reality rather than by acting inside a single objective environment.

## In Scope

- world construction at the runtime surface
- the workspace as one implementation slice, not the whole theory
- capability and resource projection as world semantics

## Out Of Scope

- the initial Quine systems proof
- revisable time as the main argument
- generic philosophical realism debates without runtime consequences

## Current State

- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- source inventory: [`./support/sources-index.md`](support/sources-index.md)
- novelty note: [`./support/novelty-read-2026-04-15.md`](support/novelty-read-2026-04-15.md)
- Linux workspace physics already provide a narrow implemented slice: `QUINE_WORKSPACE_ROOT`, `QUINE_WORKSPACE`, transactional mutation visibility, and projected `world/v0` entries under `agent/<session>/world/`
- upstream `process-as-fs` now has a coherent public process/control bridge, so this paper is no longer blocked on a totally missing runtime surface
- upstream bridge: [`../process-as-fs/README.md`](../process-as-fs/README.md)
- neighboring time paper: [`../revisable-time/README.md`](../revisable-time/README.md)

## Current Novelty Read

- broad novelty is constrained
- the main surviving angle is not sandboxing or branching by itself, but a first-class runtime world object that unifies scope, capability, observation, mutation visibility, revision addressing, and child-world adoption
- `Plan 9` / Unix namespace lineage is part of the deliberate design ancestry here, so the paper should treat projected-world mechanisms as ancestry and focus its claim at the agent-runtime semantics layer

## Current Narrow Claim

- On Linux, Quine already runs inside a projected workspace reality rooted at `QUINE_WORKSPACE_ROOT` and `QUINE_WORKSPACE`.
- Workspace writes are transactional and revision-addressable rather than naive host mutation.
- `switch_world` already gives that workspace slice an explicit rewind/return mechanic.
- The runtime also projects `agent/<session>/world/{status.json,resources.json,events.jsonl}`, so the workspace is no longer just hidden plumbing; it is the first inspectable `world/v0` object.
- This is enough to claim `workspace as first subjective world slice`, but not yet enough to claim a complete theory of capability/resource/world projection.

## Canonical Acceptance Path

- acceptance unit: `workspace-projected subjective world slice`
- coverage inventory: [`../../../tests/runtime/COVERAGE_MAP.md`](../../../tests/runtime/COVERAGE_MAP.md)
- `L1`: `internal/tools/tools_test.go`, `internal/tools/job_attrs_linux_test.go`
- `L2`: `test_workspace_overlay_commit`, `test_workspace_overlay_rollback`, `test_workspace_overlay_absolute_path`, `test_workspace_overlay_lima_probe`, `test_switch_world_restores_prior_revision`, `test_workspace_unsupported_on_non_linux`
- `L3`: `sandbox`, `switch-world`, `workspace-shadow`, `workspace-absolute`
- `L4`: `sandbox-emergent`, `workspace-shadow-emergent`, `switch-world-emergent`, `logic-bomb`
- current gap: `closure-gap` - current evidence still bundles workspace overlay physics, sandbox isolation, revision rewind, and projected world export into one broad claim.
- current gap: `harness-gap` - guest-scoped sandbox fixtures and dedicated live `L2` checks for `world/v0` are still weak.
- current gap: `emergent-gap` - passing baselines remain uneven for revision-handle uptake and explicit `[FS MUTATIONS]`-driven workflows.

## Experiment Feeds

- `none-yet` - none - not-for-paper-yet - no experiment feed is currently declared for this paper

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - no stable evidence backbone is currently declared for this paper

## Strongest Evidence

- existing Linux workspace physics and `[FS MUTATIONS]` behavior
- projected world surface under `agent/<session>/world/`
- the emerging need to model capability surfaces as locally constructed world facts
- filesystem/public-surface pressure from the system track

## Blocking Gaps

- runtime support exists, but the paper-ready evidence spine is still missing
- the paper risks drifting into philosophy unless it stays anchored to concrete runtime mechanisms
- workspace must be framed as first instance, not total theory
- the current claim is stronger for workspace projection and revision mechanics than for general capability or resource projection
- capability and resource projection are still mostly design pressure rather than explicit runtime objects

## Next Useful Work

- treat the existing workspace/world slice as the first concrete instance, not just prehistory
- harden the new inspectable `world/v0` object with live acceptance before broadening rhetoric
- decide which `world/v0` dimension should be strengthened first: capabilities, tools, resources, or observation surfaces
- separate world semantics cleanly from time semantics

## Unsafe Overclaims

- do not present workspace projection as if Quine already provides a total private world or complete capability ontology
- do not let this paper drift into abstract metaphysics without explicit runtime consequences
- do not borrow revisable-time claims before the world slice itself is more narrowly stabilized
