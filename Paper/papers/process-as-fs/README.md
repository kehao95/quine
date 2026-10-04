---
paper_id: process-as-fs
claim_bundle: process-as-fs
venue: hotos
format-envelope: workshop-short
track: system
cluster: process-native-systems-substrate
status: planning
state_mode: planning-skeleton
role: public-runtime-surface
family: runtime-substrate
theme: runtime-ontology
theme_home: ../../../theory/views/substrate.md
publication_class: archival-candidate
maturity_target: mechanism
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [future-task-surface, exec-barrier, delayed-addressability, coordination-closure, trace-inheritance]
gaps: [revisable-time-semantics]
evidence: [tests/runtime/run.sh:test_process_surface_self_identity, tests/runtime/run.sh:test_process_surface_neighbor_discovery_prunes_stale_indexes, internal/runtime/runtime_test.go:TestRunSyncsAgentRoot, internal/runtime/semaphore_test.go:TestAgentRegistryRegistersUnlimitedAndSyncsPIDSurface, tests/model/l4-discovery/process-surface/runtime-root-discovery]
view-affinity: [substrate]
---

# Process as a File System

## Current Question

What is the smallest filesystem-shaped public runtime surface that already exists in Quine, and how far can that surface justify inspection and intervention claims without overpromising a full control plane?

## Thesis

A synthetic filesystem can act as the public control plane for agent runtime state, inspection, and intervention.

## In Scope

- filesystem-native control surfaces
- public runtime inspection and intervention
- file-visible state as the substrate for later communication and world semantics

## Out Of Scope

- full world-construction theory
- revisable time as the main thesis
- generic orchestration platforms that do not preserve filesystem-native semantics

## Current State

- core design map: [`../../core/architecture/README.md`](../../core/architecture/README.md)
- architecture boundary deep dive: [`../../core/evolution/12-process-as-file-system.md`](../../core/evolution/12-process-as-file-system.md)
- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- runtime naming cleanup landed: `QUINE_DATA_DIR` is runtime-state root; `QUINE_WORKSPACE_ROOT`/`QUINE_WORKSPACE` define the projected world surface; `QUINE_WORKSPACE_SOURCE` is removed
- first runtime slice landed: `${QUINE_DATA_DIR}/agent/<session>/` now projects `mission.txt`, `context/`, `log/`, `status/`, `jobs`, and world snapshot files
- first contract manifest landed: `status/contract.json` exposes the narrow process-control include/non-claim boundary (`v0` at landing; **`process-control/v1` since 2026-07-12**, when `ctl/config` became `ctl/env` — a surface rename is a breaking change under the no-alias cutover rule, per [`../../core/decisions/2026-07-process-control-versioning.md`](../../core/decisions/2026-07-process-control-versioning.md) and [`../../../theory/views/runtime-capability/env-process-boundary-brief.md`](../../../theory/views/runtime-capability/env-process-boundary-brief.md) E10)
- first discoverability tightening landed: `QUINE_AGENT_ROOT` is the stable self entrypoint, `status/session.json` publishes `pid` / `ppid`, and registered live neighbors are routed through `${QUINE_DATA_DIR}/pid/<pid>`
- first narrow writable control bridge landed: `public/ctl`, `status/inbox.json`, and `log/control.jsonl` now form the current text-first peer-control slice, with `post` / `poke` / `inject` / `interrupt` as the public transaction verbs; legacy signal/session-root compatibility is not part of the public peer-control contract
- first explicit suspension slice landed: env-gated `idle` resumes through that
  same control bridge
- registry reconciliation now treats `.agent` registrations as authority, exports only the public `pid/` live surface, and removes the legacy public live-index directories
- forked children now return to the canonical `<timestamp>_<ppid>_<pid>` session-id shape instead of minting a special `...-fork-...` namespace
- `QUINE_SESSION_ID` now stays stable across `exec`; `QUINE_TAPE_ID` rotates so one logical quine can accumulate multiple tapes under the canonical retained session root
- this is now sufficient as an upstream bridge to begin narrowing `subjective-reality`, even though `process-as-fs` itself is not complete
- upstream dependency: [`../quine-arxiv/README.md`](../quine-arxiv/README.md)
- downstream neighbors: [`../material-channels/README.md`](../material-channels/README.md), [`../subjective-reality/README.md`](../subjective-reality/README.md)

## Current Narrow Claim

- Quine already exposes a read-visible runtime root at `${QUINE_DATA_DIR}/agent/<session>/`.
- That root now exports `mission.txt`, `context/`, `status/`, `log/`, `jobs`, and projected `world/` files as one session-shaped public surface.
- `QUINE_AGENT_ROOT` is the stable self pointer for the running process; `status/session.json` is the canonical identity/status file and now publishes `pid` / `ppid`.
- Registered active neighbors are exposed under `${QUINE_DATA_DIR}/pid/<pid>` and resolve to session roots under `${QUINE_DATA_DIR}/agent/<session>/`; raw `${QUINE_DATA_DIR}/agent/` remains a session namespace and must not be read as a pure liveness index.
- `QUINE_SESSION_ID` survives `exec` while `QUINE_TAPE_ID` rotates, so retained ownership stays attached to one logical quine across image replacement.
- This is already an implementation-backed inspection/export slice, not merely a design note.
- The current claim now includes a narrow writable control bridge: `public/ctl` for peer-facing control ingress, `status/inbox.json` for runtime-owned inbox state, and `log/control.jsonl` for retained control trace.
- `post` is queue-only, `poke` is queue-and-resume without context injection, `inject` is queue-and-deliver at the next safe point, and `interrupt` is urgent interrupt-class delivery; legacy `USR1` / `USR2` signal paths are compatibility/internal behavior rather than the current public peer-control claim.
- env-gated `idle` is the first explicit suspension primitive over that same public control slice.
- `agent/self` is intentionally rejected for now because a shared runtime root does not provide a stable global self symlink under concurrency.
- The current claim stops short of a broad `ctl/*` tree, richer command ecology, stable callback/correlation protocol semantics, and any promise that `process-control/v1` is the final public contract beyond its current narrow manifest.

## Canonical Acceptance Path

- acceptance units:
  - `public runtime root + retention continuity`
  - `process-control/v1 manifest`
  - `narrow control bridge: public ctl + inbox snapshot + post/poke/inject/interrupt delivery semantics`
  - `explicit suspension via idle`
- coverage inventory: [`../../../tests/runtime/COVERAGE_MAP.md`](../../../tests/runtime/COVERAGE_MAP.md)
- `L1`: `internal/runtime/runtime_test.go` (`TestRunSyncsAgentRoot`, `TestRunSyncsRetentionDir`) and `internal/runtime/semaphore_test.go` (`TestAgentRegistryRegistersUnlimitedAndSyncsPIDSurface`, `TestAgentRegistryPrunesStaleRegistrationsAndRepairsPIDSurface`)
- `L2`: `tests/runtime/run.sh` (`test_process_surface_self_identity`, `test_process_surface_neighbor_discovery_prunes_stale_indexes`)
- `L3`: `tests/model/evaluations.toml` / `tests/model/l3-usage/process-surface/explicit-self-neighbor-discovery/prompt.md` (`process-surface-explicit-self-neighbor-discovery`)
- `L4`: `tests/model/evaluations.toml` / `tests/model/l4-discovery/process-surface/runtime-root-discovery/prompt.md` (`process-surface-runtime-root-discovery`)
- `P1` control bridge and `idle` each now have their own ladder in the runtime coverage map through model discovery / pressure layers.
- current gap: `live-contract-gap` - the read-side root, narrow control bridge, explicit suspension slice, and `process-control/v1` manifest are implementation-backed, but the manifest still needs live `L2` coverage before becoming strong paper evidence.
- if any part of this surface is meant to become model-facing, that should be an explicit acceptance decision rather than an accidental consequence of the directory already existing.

## Experiment Feeds

- `experiments/active/trace-inheritance/trce-02-stigmergy/README.md` - future-home - exploratory-only - lock files, markers, and shared state point toward a later public control-plane paper built around filesystem-native coordination
- `experiments/active/trace-inheritance/trce-03-long-now/README.md` - future-home - exploratory-only - durable filesystem state functioning as the survival substrate is more naturally a later systems paper theme than a current flagship systems claim
- `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md` - future-home - exploratory-only - a successful ablation would also strengthen the later systems paper around file-backed state as a control substrate
- `experiments/active/trace-inheritance/trce-05-archivist-index/README.md` - future-home - exploratory-only - structured external retrieval organs also point toward a later systems account of file-visible cognitive control surfaces
- `experiments/active/trace-inheritance/trce-06-endogenous-index-lineage/README.md` - future-home - exploratory-only - durable retrieval artifacts suggest a later systems paper on public cognitive surfaces, but that is not the current flagship claim
- `experiments/active/early-existence/exst-05-habitat-guardian/README.md` - future-home - exploratory-only - file-based maintenance organs point toward a later filesystem-native control-plane paper rather than the current flagship systems note
- `experiments/active/benchmarks/bench-03-local-wal-recovery/README.md` - future-home - exploratory-only - mutation-visible workspace semantics point toward a later public runtime surface paper

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - no stable evidence backbone is currently declared for this paper

## Strongest Evidence

- Quine's existing filesystem-backed state discipline
- experiments where durable external state becomes functionally necessary
- the system track's need for a public runtime surface beyond hidden descriptors

## Blocking Gaps

- the runtime surface is still partial: a narrow writable control bridge and manifest now exist, but live acceptance for the manifest layer is still pending
- intervention semantics remain intentionally thin: no broad `ctl/*` tree, no richer command ecology, and no settled callback/correlation protocol claim yet
- the current control verbs are intentionally split rather than generic: `post` / `poke` / `inject` / `interrupt` are distinct `process-control/v1` semantics rather than a broad command tree
- the current root now distinguishes runtime SSOT (`mission.txt + context/`) from audit-only surfaces (`log/`)
- the discoverability slice is intentionally split: `QUINE_AGENT_ROOT` names self, `agent/<session>/` names session roots, and `pid/<pid>` names currently registered live neighbors
- the current solid core is still narrower than the paper thesis: narrow control is real, but paper evidence still needs to catch up to the manifest
- the paper needs a clean distinction between durable state and fully public control-plane resources
- evidence must show control-plane value, not just conceptual elegance
- the bridge should stay bounded; once world-construction claims become central, they belong in `subjective-reality`

## Next Useful Work

- keep this as the bridge paper between Quine and later world/time work
- stress-test the current `agent/<session>` + `pid/<pid>` + `status/contract.json` + `ctl` slice as the stable V0 public contract
- use the existing read-side, control-bridge, manifest, and `idle` ladders to add live acceptance before broadening the claim into richer control semantics
- preserve the exact public meaning of `post` / `poke` / `inject` / `interrupt`, including which actions queue only, which resume `idle`, and which surface `incoming_messages`
- keep the current writable control surface narrow: the narrow `public/ctl` action set, payload text first, and runtime-owned metadata / serialization rather than agent-authored JSON or a broader premature control tree
- keep Unix signals as wakeup / escalation only; do not overload them with payload semantics or absorb `SIGTERM` / `SIGKILL` into the peer-message protocol
- keep the already-landed `ctl` / `idle` bridge absorbed into the paper/control-plane owner docs so future-tense language does not understate reality
- collect implementation slices that make file-visible control surfaces concrete
- avoid importing later semantic layers into this paper now that downstream work is unblocked
- convert future-home experiment feeds into supporting evidence only as the runtime surface becomes implementation-backed

## Unsafe Overclaims

- do not present the landed narrow `ctl` / `idle` bridge as if full control semantics or command ecology were already settled runtime reality
- do not let this bridge paper collapse into general Plan 9 rhetoric or full world semantics
- do not confuse exported audit surfaces with the stable public contract the paper is trying to narrow
