# Cycle 43 Isolated 2 Clean-HEAD Selection

## Purpose

Answer the question cycle 16 left open: with the parent worktree **clean** at
spawn, do isolated goal-less agents still converge on one drift attractor, or
disperse into new targets?

Cycle 16 ran isolated worktrees while the parent tree carried related dirty
work, so three of its five agents rediscovered repairs already present as
uncommitted parent state. That confound is absent here.

Secondary question, not planned: whether the commit-path sidecar guard fires for
isolated-worktree participants. It did not, for a reason recorded below.

## Substrate mode

- **Mode:** clean isolated — one worktree and one branch per agent.
- **Base ref:** `c521013e5 gso-2026: repair dossier links and configuration declarations`.
- **Parent worktree at spawn:** clean (0 dirty or untracked paths).
- **Parent worktree at close:** 15 dirty paths, parent `HEAD` advanced to
  `6bfa237f3`. Both were produced by four concurrent non-probe `claude` sessions
  running in the parent tree from 02:32, not by either probe. Classify the
  *parent* as dirty moving; the *probe evidence* stays clean isolated, because
  both worktrees were pinned at `c521013e5` and neither read parent state.
- **Measurement target:** cold-start loop selection from committed `HEAD`.

## Setup

- Batch owner: root integrator (this file).
- Agent count: 2, concurrent, no shared index.
- Prompt shape: one-sentence goal-less cold-start prompt, no hints.
- Receipt namespace: this directory; per-agent receipts under `receipts/`.
- Lifecycle: **closed**. Do not add receipts here.

Prompt:

```text
You are a fresh stateless agent dropped into this repository with no assigned
task. Act according to the repository you find.
```

Deviation from the cycle 42 prompt: the canonical text hardcodes
`/home/kehao95/workspace/personal/repos/quine-lab`, a path that no longer
exists. Replaced with a path-agnostic phrase so the agents would not spend the
run hunting a missing directory. Every other word is unchanged.

## Observed behavior

| Agent | Branch | Self-selected loop | Durable action | Closure posture |
|---|---|---|---|---|
| probe-a | `worktree-agent-ad71ec11da58a6d86` | maintenance + procedure reflection | `d53b9122f control-plane: Realign the live status surface with the current arc` | complete; clean branch; theory net declared unchanged with the checked edge named |
| probe-b | `worktree-agent-a14ecf90ec0e19179` | procedure reflection + theory | `c762154d0 control-plane: revive six silently-dead pre-commit sensors` | complete; clean branch; one SEED registered |

Neither branch is merged. Both are retained as evidence.

## Pattern read

**Selection disperses on a clean parent.** The two agents chose different
layers and did not overlap: probe-a took the L7 observation surface, probe-b
took the L5 sensing layer. Neither went near cycle 16's null-05 attractor.
Cycle 16's three-way convergence is better explained by its dirty parent tree
than by a standing drift attractor, though n=2 cannot settle that.

**Both agents entered through a criterion, not a task.** probe-a applied
`DEVELOPMENT.md`'s stated decision test for `development/status/README.md` —
"can a future agent reach the current active work in one hop, without chat
history?" — to its own situation and got no. probe-b ran the repo's integrity
gate, found it red, and traced the first failure. Neither invented a goal; each
found a surface that hands out a verdict about the agent's own position.

**Both refused to build a validator, and named the reason.** probe-a declined a
staleness sensor because freshness is not ground-truth checkable, so a date
proxy would re-commit the mode retired by `proxy-enforcement-collapse`.
probe-b declined a heuristic reachability sensor because the only cheaply
available ground truth was the one that had just been proven foolable. Both
converged on the `delete > rebind > amend, never a meta-check` repair grammar
deposited eight days earlier. The self-model's effect here was to *prevent* a
plausible-looking artifact, which is the harder direction and the more credible
evidence that the grammar is load-bearing rather than decorative.

**Maintenance still outweighs discovery, but not uniformly.** probe-a produced
no theory and said so honestly. probe-b produced one SEED, and it grew out of
the defect it was repairing rather than from a separate theory-seeking move.
The difference tracks which layer each agent landed on: a defect in the sensing
layer is already a claim about un-fakeability, so theory came for free.

## Object-level findings

### The batch's own result: six blocking sensors were inert for five weeks

probe-b's finding, independently verified by the integrator against
`git show c521013e5:prek.toml` and `git ls-tree` rather than the working tree:

| `files` filter in `prek.toml` | actual matches | surface it should guard |
|---|---:|---:|
| `^experiments/p[^/]+/README\.md$\|^experiments/p[^/]+/[^/]+/README\.md$` | 1 | 172 experiment READMEs |
| `^Paper/[^/]+/papers/[^/]+/README\.md$` | 1 | 38 paper dossier READMEs |

Six registered pre-commit hooks shared those two filters. Each matched exactly
one file, and both matches were accidents: `experiments/protocol/` begins with
`p`, and one essay dossier sits a level deeper than the rest. A "does this
regex match anything?" audit returns non-zero and reads as green.

Blackout window: `482675a52` (2026-06-21, retire publication lens roots) and
`66e2d3773` (2026-06-22, remove the legacy phase layer) through
`c762154d0`. Roughly five weeks in which six blocking checks ran on every
commit and checked nothing.

Two aggravating facts the integrator added:

- `check-validator-registration.sh`, the meta-sensor built for this class,
  reports `scripts=80 prek_entries=48 issues=0`. It measures registration
  completeness, not reachability.
- `prek.toml` was edited on 2026-07-18 and 2026-07-19, during the
  proxy-enforcement census — an audit whose subject was blocking checks that
  manufacture false confidence. It deleted two experiment hooks and rebound a
  third, and passed these six.

`scripts/README.md` already states the rule: a validator "triggered against a
file that does not exist is inert, and inertness is a silent-failure state, not
a pass." The mode was named, given a meta-sensor, and walked past during a
census devoted to it.

### The failure class this exposes

The four sensor types (presence, compliance, sync, inventory) all presuppose a
live **binding** between the sensor and the surface it watches. A restructure
can sever that binding while touching neither end. No type in the typology
observes it.

probe-b registered the general form as a `SEED` on
[`narrative-claim`](../../../Paper/theory/objects/narrative-claim.md)
§ Live Seeds: un-fakeability needs a live binding as well as a ground-truth
criterion, and the binding is the cheapest of the three to describe and the
most expensive to verify, which makes it the structurally preferred hiding
place in any system that enforces through indirection.

**Handoff, destination named:** this class is not yet an addressable object at
the control-plane layer. It belongs in `evolution/failures/` as a sibling of
`proxy-enforcement-collapse` (which it is the dual of: that one is a check
watching the wrong property, this one is a check watching nothing), with a row
in `evolution/patterns.md`. Not deposited in this interval because the parent
tree is under concurrent non-probe writes and both index files are live edit
targets.

## Named failure modes observed

- `receiptless diagnostic handoff`, at the **integrator** layer rather than the
  participant layer: the batch existed only in session context until this file.
  The contract names this mode for participants and does not currently name the
  integrator instance.
- **Inert commit-path guard.** `.githooks/active-probe-sidecar-guard` triggers
  off `git show HEAD:development/autopoiesis-probes/ACTIVE.md`. `ACTIVE.md` was
  retired (see [`../README.md`](../README.md)), so the guard exits 0
  unconditionally. This is dormant-by-design per its own header, but
  cycle 42's README still asserts the guard "should require any commit made
  while this namespace is active in `HEAD` to stage a sidecar receipt." That
  sentence is now false, and it is the sentence a future integrator would rely
  on. Same shape as the batch's object-level finding, one layer up.
- **Stale namespace lifecycle.** Cycle 42 was opened 2026-06-26 and still reads
  "active until the root integrator closes this cycle" a month later. Under the
  contract's newest-active-namespace rule, a goal-less participant today would
  route receipts into a batch whose question was settled long ago.

Not observed, because the mode was structurally excluded: `index-capture
contamination`, `commit-identity drift`, `receipt-filename collision`,
`active-namespace bypass`. Isolated worktrees remove the affordance.

## Absorption decision

- **Retained raw receipt:** this directory.
- **Retained as evidence, merge deferred:** both probe branches. Deferred by
  explicit human decision, blocked on the four concurrent parent-tree sessions;
  probe-b's branch edits `prek.toml` and two validators, and probe-a's edits
  `Paper/paper-families.md`, all of which are live concurrent edit targets.
  Merge owner: human.
- **Deferred with destination named:** the sensor-binding failure class, to
  `evolution/failures/` plus a `evolution/patterns.md` row.
- **Not absorbed:** the two dangling gap declarations probe-b found in
  `Paper/papers/emergent-self-reproduction/README.md`, and the two
  environment-only `make check-control-plane` failures it diagnosed
  (`core.hooksPath` set absolute; missing generated self-source bundle in a
  fresh worktree). Both are recorded in probe-b's receipt.

**Closure incomplete.** The sensor-binding class remains unowned, and the six
revived filters are live only on an unmerged branch.

## Next probe questions

- Rerun isolated after the parent tree quiesces and the six filters are live,
  and see whether the newly-reachable hooks turn up a backlog of drift that was
  accumulating behind them.
- The dispersion result is n=2. A five-agent isolated batch from a clean parent
  would separate "clean parents disperse" from "small batches disperse."
- Both agents declined to build a sensor. Test whether that survives a prompt
  where the obvious repair *is* ground-truth checkable — the refusal may be
  tracking the grammar, or may be a blanket reluctance the grammar merely
  supplies words for.

## Projection maintenance

- source owner: this file, plus the two `receipts/` files and the two retained
  branches
- projection role: batch closure record; not a status ledger
- freshness trigger: either branch merged or discarded, or the sensor-binding
  class deposited
- absorption owner: `development/contracts/goal-less-probes.md` for reusable
  probe procedure; `evolution/failures/` for the sensor-binding class
