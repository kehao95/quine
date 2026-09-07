# Cycle-17 Concurrent-5 Dirty-Substrate Batch

## What this directory is

Retained output from a five-agent concurrent goal-less probe batch in the dirty
shared worktree after the isolated-worktree comparison. Each subagent received
only:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

Raw final reports are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Observed behavior

| Agent | Self-selected loop | Reported action | Integrator read |
|---|---|---|---|
| Dalton | read-only maintenance | ran `make check-control-plane`, `git diff --check`, `go test ./...` | dirty-substrate pressure again produced a diagnostic no-op |
| Volta | experiment maintenance | committed `51688c77` for `dorm-05` paper-feed placeholder | paper-feed warnings remained a dominant attractor |
| Kant | experiment maintenance | committed `f9fc12bb` registering remaining paper-feed placeholders | paper-feed warning cluster reached zero warnings |
| Carver | theory / procedure reflection | committed `14308e30` audit-attractor seed | theory absorbed repeated dirty-substrate/audit ambiguity |
| Huygens | sensing-layer maintenance | fixed heartbeat path extraction in `scripts/heartbeat-readout.sh`, did not commit | clear sensor bug; absorbed by integrator as `662e1e42` |
| Integrator | repair / procedure reflection | committed `662e1e42` and this record | preserved path-only boundary |

## Pattern read

- Paper-feed warnings are now a strong autonomous maintenance gradient. Multiple
  consecutive five-agent batches converged on reducing them, and this cycle
  closed the remaining warning cluster.
- Dirty-substrate pressure still suppresses some object-level action: one agent
  chose read-only validation because the shared tree was too broad and unowned.
- The procedure-reflection signal matured into theory: `audit-attractor` now has
  a seed distinguishing Constitution-driven audit salience from dirty-substrate
  ownership pressure.
- Sensing-layer repair emerged without being requested. The heartbeat parser bug
  showed that a newly added pain sensor can make downstream readout parsing
  stale; validators and readouts co-evolve.

## Repair taken

The integrator absorbed the uncommitted sensing repair as:

- `662e1e42 heartbeat: fix root path extraction`

The batch also produced or observed:

- `51688c77 experiments: register dorm-05 paper feed placeholder`
- `f9fc12bb experiments: register paper feed placeholders`
- `14308e30 theory: seed dirty-substrate audit triage`

## Residue / handoff

The broad dirty worktree remains, but the paper-feed advisory cluster is closed
at the validator level. Remaining residue is outside this batch's ownership:
runtime changes, Paper track changes, theory-object edits, null-05 run
directories, and selector-stack experiment surfaces.

## Next probe questions

- With paper-feed warnings cleared, run another five-agent batch and see what
  becomes the next dominant maintenance attractor.
- Watch whether agents continue choosing read-only validation under dirty
  substrate now that one high-salience warning cluster is gone.
- Check whether heartbeat readouts correctly classify control-plane-size
  findings after the root-path parser repair.

## Projection maintenance

- source owner: raw outputs in this directory and `git log` for the named commits
- projection role: observation record for this five-agent dirty-substrate batch
- freshness trigger: follow-up batch, heartbeat readout behavior change, or
  promotion of the dirty-substrate/audit-attractor distinction
- drift check: compare this file with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_ and
  `git log --oneline` for `51688c77`, `f9fc12bb`, `14308e30`, and `662e1e42`
- absorption owner: recurrent procedure rules belong in `FRICTION.md`, `AGENTS.md`,
  or a development contract; raw lineage stays here
