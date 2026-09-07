# Cycle-21 Concurrent-5 Dirty Moving-Substrate Batch

## What this directory is

Retained output from a five-agent concurrent goal-less probe batch in the dirty
shared worktree, run after the moving-substrate receipt rule was added. Each
subagent received only:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

Raw final reports are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Substrate mode

- **Mode:** shared dirty worktree, high concurrency, moving `HEAD`.
- **Start ref:** `b2020424`.
- **End ref:** `69142e24`.
- **Parent worktree at spawn:** 69 `git status --short` entries, no staged diff.
- **Observed movement:** heartbeat and probe-contract commits landed during the
  batch; one agent created and then reverted an accidental heartbeat-readout
  deletion.
- **Measurement target:** whether explicitly measuring start/end refs improves
  interpretation of shared-dirty behavior.

## Observed behavior

| Agent | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|
| Fermat | validation maintenance | ran control-plane and Go tests; no edits | complete diagnostic interval |
| Singer | validation maintenance | validated `succ-48`; accidentally deleted then reverted a concurrent heartbeat readout | complete after revert |
| Galileo | theory-entry maintenance | updated `Paper/theory/objects/README.md`; no commit because file was already dirty | incomplete object-level closure, useful signal |
| Dalton | procedure reflection | committed `4f8bb8cf` retaining a late-moving dirty-green receipt | complete observation boundary |
| Herschel | diagnostic maintenance | validated gates, noted collision/gap issues, removed generated heartbeat residue | complete diagnostic interval |

## Pattern read

- **Measured movement clarified attribution.** The batch began at `b2020424` and
  ended at `69142e24`; intermediate heartbeat/probe commits explain why agents
  reported different repository states.
- **Diagnostic no-op is now the dominant safe behavior.** Three agents avoided
  durable object-level edits because the shared dirty substrate was broad and
  moving.
- **Artifact tombstoning is a moving-substrate hazard.** Singer accidentally
  deleted a heartbeat readout created by concurrent work and had to revert it.
  Herschel also removed a generated heartbeat readout as residue. In a moving
  substrate, newly appearing artifacts are evidence first, not disposable noise.
- **Theory inventory drift remains a visible attractor.** Galileo repaired the
  cold-start theory inventory but left it uncommitted because the file was
  already dirty. That keeps the object-level issue outside this probe's ownership.

## Repair taken

The durable repair is procedural, not object-level:

- the Goal-Less Probe Contract now names `moving-artifact tombstoning`
- moving-substrate probes must treat newly appearing readouts/generated outputs
  as concurrent evidence unless the current agent created them
- object-level theory inventory changes remain uncommitted because they overlap
  pre-existing dirty state

## Next probe questions

- If shared-dirty probing continues, capture per-agent start and completion refs
  directly, not only batch start/end refs.
- After parent dirty theory-index work is committed or discarded, rerun an
  isolated batch to test whether theory inventory drift remains a target.
- Consider whether moving-substrate probes should be read-only by default unless
  they operate in owned paths.

## Projection maintenance

- source owner: raw outputs in this directory and `git log` for `b2020424` to
  `69142e24`
- projection role: aggregate record for a measured moving-substrate no-task
  batch
- freshness trigger: another moving-substrate batch, change to artifact cleanup
  rules, or integration of the dirty theory inventory repair
- drift check: compare this file with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_, the
  Goal-Less Probe Contract, and `FRICTION.md` F18/F19
- absorption owner: reusable procedure belongs in
  `development/contracts/goal-less-probes.md`; raw lineage stays here
