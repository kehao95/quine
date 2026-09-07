# Cycle-19 Concurrent-5 Dirty Moving-Substrate Batch

## What this directory is

Retained output from a five-agent concurrent goal-less probe batch in the dirty
shared worktree after the isolated-worktree comparison and green-gate dirty
single probes. Each subagent received only:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

Raw final reports are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Substrate mode

- **Mode:** shared dirty worktree, high concurrency, moving `HEAD`.
- **Base window:** agents observed commits around `0f0eaf5d`, `8542bd18`,
  `b712753b`, `de215ed4`, and related concurrent probe commits.
- **Parent worktree:** broad unowned dirty state across Paper, theory,
  experiments, runtime, scripts, and untracked run/experiment artifacts.
- **Measurement target:** what no-task agents do when the control plane is green
  or advisory-only, but both the working tree and branch head are live.

## Observed behavior

| Agent | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|
| Planck | control-plane maintenance | committed `8542bd18 evaluation: compress control plane` | complete; path-scoped repair |
| Pasteur | discovery / prediction hygiene | produced `narrative-claim-P4` correction, later committed by integrator as `b987b103` | complete after hook false-positive handling |
| Halley | theory / discovery | committed `de215ed4 theory: predict pressure-driven channel substitution` | complete after collision audit |
| Faraday | validation maintenance | ran control-plane, substrate, and `succ-48` smoke checks; no edits | complete diagnostic interval |
| Godel | procedure reflection | committed `71e77243` / `b712753b` probe receipts for green-gate dirty substrate and boundary collision | complete after repairing accidental mixed commit |

## Pattern read

- **Broad loop selection persists under dirty moving substrate.** The batch
  produced control-plane compression, prediction hygiene, theory prediction,
  validation-only maintenance, and procedure-reflection.
- **Moving `HEAD` is now a first-class substrate condition.** Multiple agents
  observed branch movement or commit-boundary collision while acting. The
  problem is not only a dirty index; the branch itself is changing under the
  probe.
- **Green gates still do not authorize broad edits.** Faraday and Godel both
  found green or advisory-only validation but treated the broad dirty substrate
  as the dominant action boundary.
- **Path-scoped useful work remains viable.** Planck and Halley completed useful
  commits by staying narrow and auditing boundary contamination.
- **Hook behavior can become part of the substrate.** The integrator's
  `narrative-claim-P4` commit hit a `check-agent-diff-review` false-positive
  "files were modified" failure while reporting `clean`, requiring a targeted
  skip after focused validators passed.

## Repair taken

- Integrated Pasteur's prediction hygiene as `b987b103 theory: clarify
  narrative-claim P4 remains open`.
- Retained the uncommitted cycle-18 dirty-green repeat receipt because it
  independently observed green-gate dirty substrate plus `HEAD` movement.
- Updated probe rules so moving-`HEAD` shared-substrate probes must record start
  and end refs, and must not be interpreted as ordinary single-agent no-task
  behavior.

## Next probe questions

- After the broad parent dirty work is committed or discarded, rerun an isolated
  five-agent batch from the updated `HEAD` to test F17.
- If continuing shared-dirty stress probes, record `HEAD` at spawn and at each
  agent completion so branch movement is measured rather than reconstructed.
- Watch whether control-plane size compression becomes the next maintenance
  attractor now that paper-feed warnings are closed.

## Projection maintenance

- source owner: raw outputs in this directory, the named commits, and adjacent
  cycle-17/18 dirty-green receipts
- projection role: aggregate record for a concurrent-five no-task batch under
  dirty shared worktree plus moving-branch pressure
- freshness trigger: another moving-HEAD shared-substrate batch, a new
  substrate-isolation rule, or resolution of the parent dirty work
- drift check: compare this file with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_,
  `FRICTION.md` F16-F18, and `git log --oneline` for the named commits
- absorption owner: durable concurrency procedure belongs in the parent probe
  README or a development contract; raw lineage stays here
