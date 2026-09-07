---
id: sh-output-explosion
type: failure-mode
status: deposited
applied: false
class: runtime-memory-substrate
symptom: Large sh output truncates and is destroyed, forcing re-runs that explode the context and kill the incarnation (~5min/turn, memory death).
evidence: [l0-quine-missionless-2026-06-27-post-restructure runs, revert cfdaf9320, revert 481cf8cee]
candidate-solution: stratified context — cognition as projection of the context/state substrate at variable resolution.
owner: lab
---
# sh-output-explosion

## Symptom

When a quine agent runs `sh` against a large or dirty substrate (`git diff` of
many files, `make check-control-plane`, `find`, `ls -R`), each tool result far
exceeds the useful context window. Observed: ~5 min/turn on a dirty live repo vs
~15 s/turn on a clean tree, ending in memory death. The mechanism is a **re-run
explosion loop**: head-only truncation + destruction of the only full copy (the
completed-sync job dir) + `sh` as the only file reader leaves re-running the
producing command as the sole escape, which reproduces the same large output and
truncates+destroys again.

## Evidence

- `development/autopoiesis-probes/l0-quine-missionless-2026-06-27-post-restructure/`
  — the live-repo runs where ~5min/turn and turn-rate collapse on a dirty
  substrate were observed.
- revert commits `cfdaf9320` and `481cf8cee` — a first-attempt fix (retain the
  job dir + a `compaction_target` nudge) was judged ontological drift
  (cross-layer coupling: using the execution substrate as a cognitive backing
  store; plus prompt prose doing the agent's reasoning) and reverted.

## Root cause

Not a size problem but an **ontological conflation**: `current.jsonl` plays two
roles at once — literal append-only history *and* the window the model sees. So
truncation reads as *data loss* rather than as a rendering choice, which forces
the cross-layer reach into the execution substrate (job-dir retention) to recover
lost cognition. The deeper cut: jobs are the execution layer; context/state is
the cognitive layer; compression of a tool result is a cognitive-layer operation
that must stay within cognition, not lean on the execution vehicle. The truncation
cliff and `MemoryDeathTokens` death are both symptoms of the same conflation.

## Candidate solution (deposited — applied: false)

**Stratified context.** Separate the ground truth (the `context/state` substrate,
with a working-honesty floor) from the context (a projection rendered each turn at
variable, per-entry resolution). Three primitives reorganize the substrate:
`mark` (semantic label), `fold` (resolution degrade), `compact` (global
restructure, auto-at-danger). Truncation, compaction, and memory pressure unify
into one axis (resolution); truncation becomes a rendering resolution, not data
loss; the job dir is always destroyed (bytes live in the cognitive substrate).

Full design + clean-upgrade feasibility:
[`../../development/stratified-context.md`](../../development/stratified-context.md).

## Why not applied yet

A clean upgrade is **feasible-with-caveats** (two read-only scans), not yet
scoped or committed. The atomic core is larger than first thought: a
substrate-version field (currently absent — non-negotiable prerequisite) +
projection engine + fork/exec/recovery jsonl readers + mark-split must land
together behind a flag, because fork/exec *produce* the substrates the projection
*reads*. ~6 tests pin the old contract and must be forked into old/new sets. The
SCHEMA is a contract enacted by autonomous peers, so the cross-policy-peer hazard
must be gated by a version field both runtime and peer check. See
`development/stratified-context.md` § Feasibility & clean-upgrade assessment.

Deposited `applied: false` deliberately: the design + the upgrade path are
recorded, but enactment is an open discovery target a future agent (possibly a
goal-less probe) can find, weigh, and attempt — a test of whether the system can
self-repair its own runtime substrate, not just execute a given fix.

## Theory links

- [`../theory/objects/research-state-machine.md`](../../Paper/theory/objects/research-state-machine.md)
  — this failure-mode and its candidate solution are the **carrying instance** of
  the research state-machine (friction → principle/mechanism → design →
  constraint, with feedback). Depositing it here is the reification that makes
  the trajectory searchable.
- [`../philosophy/computational-epistemology.md`](../../Paper/philosophy/computational-epistemology.md)
  — the π (procedure) + reification seed: the fix is itself a procedure-level
  redesign (stratified context is a π-transition), and naming it reifies it for
  search.
- [`../theory/objects/narrative-claim.md`](../../Paper/theory/objects/narrative-claim.md)
  — the working-honesty floor rests on the function/narrative distinction: the
  last real I/O is function (un-fakeable); reorganized earlier content is
  narrative (honestly labeled representation).
