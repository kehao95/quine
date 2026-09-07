# Cycle-15 Concurrent-5 Dirty-Substrate Batch

## What this directory is

Retained output from a five-agent concurrent goal-less probe batch run in a
heavily dirty shared worktree. Each subagent received only:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

Raw final reports are retained in [`raw-outputs.md`](./raw-outputs.md). A
single-agent follow-up under the same prompt is retained in
[`followup-20260626-022537-validation-green.md`](./followup-20260626-022537-validation-green.md).

## Substrate mode

- **Mode:** shared dirty worktree, high concurrency.
- **Known confound:** many tracked and untracked surfaces were already in flight
  across Paper, theory, experiments, runtime, scripts, and autopoiesis probes.
- **Measurement target:** not pure cold-start target selection; this batch
  measures what goal-less agents do when autonomous action is coupled to
  unowned dirty state and concurrent Git/index pressure.

## Observed behavior

| Agent | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|
| Maxwell | procedure reflection | committed `61371abf autopoiesis: retain dirty-substrate probe` | complete for an isolated observation deposit |
| Gauss | maintenance / validation | ran `make check-control-plane`, `git diff --check`, `go test ./...`; no edits | incomplete because broad dirty state made ownership unsafe |
| Hypatia | discovery / theory | committed `0ae2be35 theory: predict sequence-carrier recovery divergence` | complete for a bounded prediction increment |
| Euclid | discovery / theory | committed `78e12780 theory: predict divergence substrate synthesis effect` | complete, but reported index lock / staged-capture pressure and used `--no-verify` after targeted checks |
| Hegel | maintenance | committed `ea05ae8b experiments: normalize structural elicitation paper feeds` | complete for a bounded paper-feed repair |

## Pattern read

- **Loop selection is now broad.** The batch produced procedure-reflection,
  validation maintenance, paper-feed maintenance, and two theory/discovery
  prediction increments.
- **Dirty substrate shapes the initial policy.** Maxwell and Gauss both treated
  the dirty tree as the dominant fact: one deposited the observation, the other
  refused to edit.
- **Shared index pressure is not hypothetical.** Euclid reported Git index
  locking and temporary unrelated staged capture during commit. This is the same
  substrate class as `gamete-contamination`, now reproduced in a five-agent
  concurrent batch.
- **Useful work still happens under collision pressure.** Three agents completed
  bounded commits, but closure quality depends on path isolation and targeted
  validation rather than aggregate green state.
- **Green validators do not dissolve dirty-substrate pressure.** The follow-up
  found `make check-discovery-lifecycle`, `make check-control-plane`, and
  `git diff --check` green, but ownership ambiguity still made broad edits the
  wrong action boundary.

## Repair taken

The cycle is retained as a dirty-substrate observation rather than treated as a
clean cold-start assay. The durable procedure repair lands in the parent probe
README: concurrent goal-less probes now need an explicit substrate-mode receipt
so future agents can distinguish clean loop-selection evidence from shared-index
stress evidence.

## Next probe questions

- Run a five-agent goal-less batch in isolated worktrees and compare loop
  selection, commit contamination, and closure confidence.
- If shared-dirty batches continue, treat them as a stress test and record staged
  residue before/after instead of reading them as neutral autonomous-selection
  probes.
- Check whether theory/discovery output remains high when index collision is
  removed.

## Projection maintenance

- source owner: raw outputs in this directory and `git log` for the named commits
- projection role: observation record for a concurrent-five goal-less batch under
  dirty shared-substrate pressure
- freshness trigger: follow-up isolated batch, another dirty shared batch, or
  promotion of dirty-substrate selection pressure into a standing contract
- drift check: compare this file with [`raw-outputs.md`](./raw-outputs.md) and
  the named commits
- absorption owner: repeated concurrency procedure belongs in this directory's
  README, `FRICTION.md`, or a development contract if it becomes repo-wide law
