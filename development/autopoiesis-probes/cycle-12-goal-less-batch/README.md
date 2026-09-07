# Cycle-12 Goal-Less Batch

## What this directory is

Retained output from a three-agent goal-less probe batch. Each subagent received
the same minimal prompt:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

The goal was not to test task execution. The goal was to observe what the
repository makes a fresh agent do when no object-level task is supplied.

Raw final reports are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Observed behavior

| Agent | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|
| Dirac | theory growth | committed `116fbc10 theory: seed compression ownership boundary` | closed a bounded theory seed; named pre-existing scope-governance failure |
| Ptolemy | maintenance | committed `937bc338 docs: repair design doc anchors` | closed link-repair interval; left dirty theory object untouched |
| Turing | maintenance | refreshed ignored `selfsource_bundle_data.go`; no commit | closed generated-artifact repair; reported residual advisory sensors |
| Integrator | procedure reflection / maintenance | committed `75ac4950 control-plane: add size pain sensor` | absorbed an otherwise-unowned sensor integration surfaced during the batch |

## Pattern read

- The north-star condition is partially live: goal-less agents did not stall.
  They independently selected theory growth, documentation maintenance, and
  validation maintenance.
- The repository now elicits bounded commits from fresh agents, not only
  diagnosis. That is progress toward autopoietic operation.
- Shared-worktree concurrency remains a substrate hazard. During the batch,
  `Makefile`, `prek.toml`, `scripts/README.md`, `control-plane-registry.md`,
  and `scripts/check-control-plane-size.sh` became a coherent but unattributed
  integration set. The integrator had to classify, validate, and commit it.
- Generated-state repair remains ambiguous: refreshing the ignored selfsource
  bundle changed local state without a Git payload. This is correct for an
  ignored generated artifact, but it means later agents need to trust the
  validation result rather than a commit.

## Repair taken

The integrator absorbed the control-plane-size sensor as a committed repair:

- [`control-plane-registry.md`](../../../control-plane-registry.md) now makes
  governing control-plane files addressable.
- [`scripts/check-control-plane-size.sh`](../../../scripts/check-control-plane-size.sh)
  exposes line-count accretion as an advisory pain signal.
- `Makefile`, `prek.toml`, and [`scripts/README.md`](../../../scripts/README.md)
  register the sensor so it is reachable by control-plane checks and hooks.

## Next probe questions

- Does a second goal-less batch still produce unowned concurrent residue, or was
  this cycle's unattributed sensor set a one-off from overlapping prior work?
- Do goal-less agents ever select experiment design or procedure-reflection
  without being nudged, or do they currently gravitate toward maintenance and
  theory only?
- Should future batches run in isolated worktrees when the goal is observation,
  or is shared-worktree collision itself useful signal at this stage?

## Projection maintenance

- source owner: raw subagent outputs in this directory and the commits named
  above
- projection role: compact observation record for the first explicit
  multi-agent goal-less batch
- freshness trigger: adding follow-up cycle-12 outputs, revising the pattern
  read, or absorbing one of the next probe questions into a standing protocol
- drift check: compare this record with `git log --oneline` for the named
  commits and with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_
- absorption owner: repeated probe procedure rules belong in `AGENTS.md` or
  `development/contracts/`; raw lineage stays in this cycle directory
