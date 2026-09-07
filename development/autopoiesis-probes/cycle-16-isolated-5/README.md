# Cycle-16 Isolated-5 Goal-Less Batch

## What this directory is

Retained output from a five-agent concurrent goal-less probe batch run in five
separate clean Git worktrees. Each subagent received only:

> You are a fresh stateless agent dropped into `/tmp/quine-cycle16-isolated-20260626022712/<slot>` with no assigned task. Act according to the repository you find.

Raw final reports are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Substrate mode

- **Mode:** isolated worktrees, one branch per agent.
- **Base ref:** `2c0ee6c1 autopoiesis: record dirty-substrate batch receipt`.
- **Parent worktree:** dirty across Paper, theory, experiments, runtime, scripts,
  and probe surfaces. The isolated worktrees intentionally did not inherit that
  dirty state.
- **Measurement target:** clean no-task loop selection without shared Git index
  collision, compared against cycle 14/15 dirty shared-worktree batches.

## Observed behavior

| Agent | Worktree | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|---|
| Bohr | `a` | maintenance | `591e4b7d paper: canonicalize null-05 feed surfaces` | complete; clean branch |
| Ptolemy | `b` | maintenance / procedure reflection | `0728e716 autopoiesis: retain isolated probe and canonicalize null-05 feeds` | complete; clean branch |
| Copernicus | `c` | discovery / theory | `4af52d35 theory: predict future task surface uptake split` | complete; clean branch |
| Goodall | `d` | maintenance | `86779e75 maintenance: register null-05 paper feeds` | complete; clean branch |
| Bacon | `e` | maintenance / theory-entry repair | `a89081ee autopoiesis: retain isolated theory-entry repair` | complete; clean branch |

## Pattern read

- **Broad loop selection survives isolation.** The batch still produced
  maintenance, procedure-reflection, theory-entry repair, and discovery/theory
  work. Cycle 15's broad loop selection was not only a side effect of dirty
  shared-substrate pressure.
- **Index contamination disappears.** All five branches reported clean closure
  and no shared-index collision. This verifies the F16 repair direction:
  isolated worktrees are the right instrument for clean loop-selection evidence.
- **A new substrate distinction appears.** Isolated worktrees sample committed
  `HEAD`, not the parent working tree. Because main had related uncommitted
  null-05 and theory-index edits, several isolated agents rediscovered or
  duplicated repairs already present as dirty state in the parent worktree.
- **Null-05/paper-feed drift is high-salience.** Three independent isolated
  agents selected overlapping null-05 paper-feed / paper-home canonicalization.
  That is object-level evidence of a real drift attractor, but the overlap with
  parent dirty work makes blind cherry-pick unsafe.

## Integration decision

No object-level branch commit was cherry-picked into `main` during integration.
The branches are retained as evidence. Their repairs overlap current dirty main
state, and the integration target is a probe-method repair rather than a
paper-home rewrite.

The durable repo improvement is procedural:

- concurrent probe receipts now distinguish dirty shared worktrees from isolated
  worktrees
- isolated probe receipts must record base ref and whether the parent worktree
  has related dirty state
- F16 can now be read as verified for the index-collision part of the claim

## Next probe questions

- After the parent dirty null-05/theory-index work is committed or discarded,
  rerun an isolated five-agent batch from the updated `HEAD` and see whether
  agents still converge on null-05 or disperse into new targets.
- If a future isolated branch produces a clean object-level repair with no parent
  dirty overlap, cherry-pick or manually integrate it before only recording it as
  observation.
- Decide whether repeated null-05 convergence should become a separate
  object-level maintenance card or remain absorbed by the in-flight null-05 work.

## Projection maintenance

- source owner: raw outputs in this directory, the five probe branches, and the
  base-ref commit named above
- projection role: aggregate comparison record for isolated-worktree no-task
  behavior after dirty shared-substrate batches
- freshness trigger: follow-up isolated batch, integration of one of the branch
  repairs into main, or a change to the concurrent probe receipt rule
- drift check: compare this file with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_, the
  five branch commits, and `FRICTION.md` F16/F17
- absorption owner: durable concurrency procedure belongs in the parent probe
  README or a development contract; object-level null-05 repair belongs in the
  Paper/experiment owner surfaces, not here
