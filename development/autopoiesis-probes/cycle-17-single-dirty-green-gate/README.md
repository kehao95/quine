# Cycle-17 Single Dirty-Substrate Green-Gate Probe

## What this directory is

Retained output from a single goal-less probe in the main shared worktree. The
agent received only:

> You are a fresh stateless agent dropped into /home/kehao95/workspace/personal/repos/quine-lab with no assigned task. Act according to the repository you find.

Raw receipt notes are retained in `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_.

## Substrate mode

- **Mode:** single agent in the shared main worktree.
- **Base ref:** `f4982de9`.
- **Parent worktree:** dirty at entry/checkpoint with 68 `git status --short`
  entries, including broad runtime, Paper, theory, experiment, script, and
  untracked run/experiment residue.
- **Index state:** no staged diff observed before action.
- **Measurement target:** no-task behavior when the strict control-plane gate is
  green but the shared substrate is already broad and unowned.

## Observed behavior

| Agent | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|
| Codex | maintenance / procedure-reflection | retained this probe receipt; no runtime/theory/experiment edits | commit amended after boundary-contamination repair |

## Validation

`make check-control-plane` passed from the dirty shared worktree. Notable output:

- experiment paper feeds: `issues=0 warnings=0`
- active doc links: `issues=0`
- experiment theory reciprocity: `dangling_errors=0 advisory_gaps=3963`
- control-plane size: advisory only, 7 registered files over the 200-line pain
  line
- selfsource freshness: OK

## Pattern read

- This is a distinct dirty-substrate case from cycle 15/16: the control-plane
  aggregate was green, so the dirty tree itself, not a red validator, was the
  dominant action-shaping condition.
- A green control plane plus broad unowned local edits selects for
  low-overlap maintenance: validate, name the substrate, and retain a receipt
  rather than editing object-level code or theory.
- The first commit attempt captured three Paper files outside the intended
  probe boundary. The commit was amended so those files remained local residue
  at that boundary and only the receipt stayed committed. A later concurrent
  theory commit absorbed part of that Paper work, making the collision visible
  in branch history as well as the index. This is a single-agent/concurrent
  instance of the same substrate-contamination class described by F16.
- The behavior supports F16's interpretation that dirty shared substrate is an
  instrument, not a neutral cold-start environment. It does not add a new
  friction entry because the current rules already cover the substrate-mode
  receipt requirement.

## Residue / handoff

The broad dirty worktree remains unowned by this probe. This cycle does not
claim runtime, Paper track, theory-object, `null-05`, selector-stack, or
untracked run residue.

Next useful boundary: after the existing broad dirty work is committed or
discarded, rerun either a single goal-less probe or an isolated batch from the
updated `HEAD` to separate persistent target selection from stale-parent
pressure.

## Projection maintenance

- source owner: this directory's raw receipt and `git log` for any receipt
  commit
- projection role: observation record for one goal-less probe in a green-gate
  dirty shared substrate
- freshness trigger: follow-up clean or dirty no-task probe, or a future rule
  change about when green-gate dirty substrates should edit versus only retain
  observation
- drift check: compare this file with `raw-outputs.md` _(raw pruned 2026-07-01 compression pass; findings retained in this README + FRICTION + META-DIRECTION)_,
  `FRICTION.md` F16/F17, and `META-DIRECTION.md`
- absorption owner: recurring dirty-substrate procedure belongs in
  `FRICTION.md`, `AGENTS.md`, or a development contract once it changes future
  dispatch policy
