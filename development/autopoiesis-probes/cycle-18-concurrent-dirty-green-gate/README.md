# Cycle-18 Concurrent Dirty-Substrate Green-Gate Probe

## What this directory is

Retained output from a single goal-less probe in the main shared worktree. The
agent received only:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

Raw receipt notes are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Substrate mode

- **Mode:** single agent in the shared main worktree, with external branch
  movement observed during the session.
- **Entry observation:** initial orientation saw a dirty worktree and recent
  `HEAD` at `6cb3ff16`; later verification showed `HEAD` at `de215ed4`.
- **Parent worktree:** dirty at entry and closure, with broad runtime, Paper,
  theory, experiment, script, and untracked run residue.
- **Index state:** no staged diff observed before action.
- **Measurement target:** no-task behavior when strict control-plane validation
  is green but the shared substrate is both broad-dirty and concurrently moving.

## Observed behavior

| Agent | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|
| Codex | maintenance / procedure-reflection | retained this probe receipt; no runtime/theory/experiment edits | commit scoped to receipt files after explicit path staging |

## Validation

`make check-control-plane` passed from the dirty shared worktree. Notable
output:

- experiment paper feeds: `issues=0 warnings=0`
- active doc links: `issues=0`
- experiment theory reciprocity: `dangling_errors=0 advisory_gaps=3963`
- control-plane size: advisory only, 6 registered files over the 200-line pain
  line
- selfsource freshness: OK

## Pattern read

- This repeats the green-gate dirty-substrate behavior from cycle 17, but adds
  one stronger condition: `HEAD` advanced while the probe was active.
- A green aggregate validator did not make object-level repair safe; the
  dominant risk was substrate ownership and commit-boundary contamination.
- The correct bounded action was observation retention plus an explicit-path
  commit. Runtime, Paper, experiment, and theory edits would have mixed with
  unowned work and blurred the measurement.

## Residue / handoff

The broad dirty worktree remains unowned by this probe. This cycle does not
claim runtime, Paper track, theory-object, `null-05`, selector-stack, `succ-48`,
or untracked run residue.

Next useful boundary: after the existing broad dirty work is committed or
discarded, rerun either a single goal-less probe or an isolated batch from the
updated `HEAD` to separate persistent target selection from stale-parent and
concurrent-substrate pressure.

## Projection maintenance

- source owner: this directory's raw receipt and `git log` for any later receipt
  commit
- projection role: observation record for one goal-less probe in a green-gate
  dirty shared substrate with concurrent branch movement
- freshness trigger: follow-up clean, isolated, or dirty no-task probe; or a
  future rule change about when moving shared substrates should edit versus only
  retain observation
- drift check: compare this file with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_,
  `FRICTION.md` F16/F17, and `META-DIRECTION.md`
- absorption owner: recurring dirty/concurrent-substrate procedure belongs in
  `FRICTION.md`, `AGENTS.md`, or a development contract once it changes future
  dispatch policy
