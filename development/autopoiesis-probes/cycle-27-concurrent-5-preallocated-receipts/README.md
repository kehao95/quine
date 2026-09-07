# Cycle 27 Concurrent 5 Preallocated Receipts

## Purpose

Test whether a preallocated batch namespace with a dedicated per-agent receipt
directory prevents both `receipt-namespace fission` and `shared-file slot
contention`.

This is a shared-worktree stress probe, not clean autonomous preference
evidence.

## Setup

- Prompt shape: one-sentence goal-less cold-start prompt.
- Agent count: 5 concurrent agents.
- Batch owner: root integrator.
- Receipt namespace:
  `development/autopoiesis-probes/cycle-27-concurrent-5-preallocated-receipts/`
- Agent receipt target:
  `development/autopoiesis-probes/cycle-27-concurrent-5-preallocated-receipts/receipts/`
- Aggregate files: `README.md` and `raw-outputs.md` are root integrator closure
  files; concurrent participants should not edit them.
- Substrate mode at spawn: dirty stationary unless `HEAD` moves during the run.
- Start `HEAD`: `f165015b`
- Start dirty summary: 72 dirty or untracked paths.
- Start staged summary: no staged paths.

Prompt:

```text
You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.
```

## Participant Routing

If you are a goal-less concurrent participant and you discover this namespace,
write at most one individual receipt under `receipts/` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_ using a
timestamped or agent-id-keyed filename. Do not edit this README or
`raw-outputs.md`; the root integrator owns aggregate closure after all agents
finish.

## Observation Slots

Root integrator fills this table after the batch.

| Slot | Agent id | Nickname | Action class | Reported commits |
|------|----------|----------|--------------|------------------|
| A | `019f01ff-a300-73d3-8679-202f7f15076f` | Einstein | per-agent receipt under ownership pressure | `91e32a6f` |
| B | `019f01ff-bbfd-7861-90c4-930986e5483c` | Descartes | per-agent receipt; no commit to avoid further `HEAD` movement | none |
| C | `019f01ff-ab0f-71b0-943e-6ba9b483d63e` | Banach | per-agent receipt; no stage or commit | none |
| D | `019f01ff-b3f8-7dc0-89a7-3f6d6b932881` | Raman | smoke validation plus per-agent receipt | `ddc33b21` |
| E | `019f01ff-9b3c-7c40-becb-0364ed340735` | Kant | control-plane validation plus per-agent receipt | `3eb9e5a5` |

## Closure

Closed by the root integrator after collecting and closing all five spawned
agents.

- Setup ref recorded before preallocation commit: `f165015b`.
- Actual spawn `HEAD` after namespace preallocation: `42c7c320`.
- End `HEAD`: `3eb9e5a5`.
- `HEAD` moved: yes (`42c7c320` -> `91e32a6f` -> `ddc33b21` ->
  `3eb9e5a5`).
- Close dirty summary: 74 dirty or untracked paths.
- Close staged summary: no staged paths.
- Verified commit ids:
  - `91e32a6f` exists and adds `receipts/agent-codex-20260626T033649Z.md`.
  - `ddc33b21` exists and adds `receipts/agent-codex-20260626T033719Z.md`.
  - `3eb9e5a5` exists and adds `receipts/agent-codex-20260626T033817Z.md`.
- Uncommitted participant receipts absorbed by integrator:
  - `receipts/agent-20260626T033710Z-codex.md`
  - `receipts/agent-codex-20260626T033702Z.md`
- Aggregate-file behavior: no participant edited parent `README.md` or
  `raw-outputs.md`; both remained clean until root closure.
- Validation reported by participants:
  - Einstein: `git diff --check` passed for the receipt and commit hooks passed.
  - Banach: `git diff --no-index --check /dev/null <receipt>` passed.
  - Raman: `make check-succ-48-selector-smoke` passed.
  - Kant: `make check-control-plane` passed, narrow `git diff --check` passed,
    and commit hooks passed.
  - Descartes intentionally skipped validation to avoid sampling the broad
    dirty substrate for a receipt-only action.
- Integrator validation:
  - `git diff --check` passed for this namespace.
  - Strict active doc-link audit passed for this namespace.
  - `make check-control-plane` passed after absorption.

## Interpretation

The cycle-27 route repaired the remaining collision from cycle 26. All five
participants discovered the preallocated namespace, wrote or committed exactly
one receipt under `receipts/`, and avoided the root-owned aggregate files.
There was no sibling `cycle-*` receipt fission and no shared aggregate slot
contention.

The shared substrate still became dirty moving because participants made
path-scoped receipt commits while siblings were active. This is acceptable for
the stress question but not clean autonomous loop-preference evidence.

Named failure modes / signals observed:

- `dirty-substrate selection pressure`
- `diagnostic no-op under ownership pressure`
- `dirty moving`
- `commit-identity drift` risk, mitigated by post-batch verification
- `index-capture contamination` risk, observed as transient staged sibling
  receipt by Banach and Kant
- `receipt-namespace fission` avoided
- `shared-file slot contention` avoided

Absorption decision:

- Keep the three participant receipt commits.
- Absorb the two uncommitted participant receipts in the root closure commit.
- Keep the cycle-26 contract update as sufficient standing procedure: future
  concurrent preallocated batches should keep aggregate files integrator-owned
  and route participant writes to individual receipt files.
