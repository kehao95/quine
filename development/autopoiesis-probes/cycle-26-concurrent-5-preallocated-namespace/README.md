# Cycle 26 Concurrent 5 Preallocated Namespace

## Purpose

Test whether committing the batch receipt namespace before spawning a
goal-less concurrent batch reduces `receipt-namespace fission`.

This is a shared-worktree stress probe, not clean autonomous preference
evidence.

## Setup

- Prompt shape: one-sentence goal-less cold-start prompt.
- Agent count: 5 concurrent agents.
- Batch owner: root integrator.
- Receipt namespace:
  `development/autopoiesis-probes/cycle-26-concurrent-5-preallocated-namespace/`
- Substrate mode at spawn: dirty stationary; reclassified to dirty moving
  after sibling commits advanced `HEAD` during the run.
- Start `HEAD`: `f83c2274`
- Start dirty summary: 71 dirty or untracked paths.
- Start staged summary: no staged paths.

Prompt:

```text
You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.
```

## Observation Slots

Raw per-agent outputs belong in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

| Slot | Agent id | Nickname | Action class | Reported commits |
|------|----------|----------|--------------|------------------|
| A | `019f01f7-5e72-7592-8090-31e3e543893b` | Sartre | maintenance/procedure-reflection receipt; strict control-plane validation | none |
| B | `019f01f7-4f41-7000-b1dc-6432ffe0fc5f` | Feynman | maintenance/procedure-reflection receipt; observed slot-A write collision | none |
| C | `019f01f7-56cf-77b0-8f3d-34989db2ad8c` | Sagan | per-agent receipt under ownership pressure; observed index lock contention | `ea20a178` |
| D | `019f01f7-4791-72f2-bd7d-38bfa502cb65` | Leibniz | bounded theory-graph repair | `be4a782a` |
| E | `019f01f7-3f5d-7413-8975-8981c5dd1bdf` | Epicurus | validation-only no-op after moving `HEAD` | none |

## Closure

Closed by the root integrator after collecting and closing all five spawned
agents.

- End `HEAD`: `ea20a178`
- `HEAD` moved: yes (`f83c2274` -> `af36656a` -> `be4a782a` -> `ea20a178`)
- Close dirty summary: 74 dirty or untracked paths.
- Close staged summary: no staged paths.
- Verified commit ids:
  - `be4a782a` exists and modifies
    `Paper/theory/objects/stigmergy-as-carrier.md`.
  - `ea20a178` exists and adds
    `agent-codex-20260626T032930Z.md` inside this namespace.
- Validation reported by agents:
  - `make check-control-plane` passed in multiple sessions.
  - `go test ./...` passed in Epicurus's validation-only session.
  - `make audit-discovery-lifecycle` and
    `make audit-succ-48-selector-smoke` passed in Epicurus's session.
  - `make audit-gap-status` advisories reduced from 4 to 3 and
    `make audit-theory-coherence` advisories reduced from 16 to 15 after
    Leibniz's repair.
  - Focused strict doc-link and surface-ownership checks passed in the
    per-agent receipt sessions.
- Integrator validation:
  - `git diff --check` passed for this namespace and the probe contract.
  - `make check-control-plane` passed after absorption.

## Interpretation

Preallocating the namespace before spawning worked: all receipt-writing agents
found `cycle-26-concurrent-5-preallocated-namespace/`; no sibling `cycle-*`
directory was created for this batch.

The batch also exposed a narrower collision surface. Sartre and Feynman edited
the shared `README.md` / `raw-outputs.md` slots directly, while Sagan avoided
that shared file after observing a concurrent writer and committed a separate
timestamped receipt. Sagan also hit `.git/index.lock` during a commit attempt,
confirming shared-index contention under concurrent hooks.

Named failure modes observed:

- `dirty-substrate selection pressure`
- `diagnostic no-op under ownership pressure`
- `index-capture contamination` pressure / shared-index contention
- `commit-identity drift` risk, mitigated by post-batch verification
- `shared-file slot contention`

Absorption decision:

- Keep `be4a782a` as an object-level theory maintenance repair.
- Keep `ea20a178` and the remaining uncommitted per-agent receipt as cycle-26
  probe evidence.
- Update `development/contracts/goal-less-probes.md` so future preallocated
  namespaces treat aggregate files as integrator-owned closure files and route
  concurrent agents to individual receipt files.
