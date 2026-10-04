# Revisable Time Novelty Read (2026-04-15)

This note records the current literature-based novelty judgment for the `revisable-time` dossier.

It is a working read, not a final paper-position commitment.

## Current Read

Broad novelty for `revisable-time` looks weak.

The general research space is already crowded by:

- transactional / saga-style agent recovery
- OS-level branch / explore / commit / abort primitives
- checkpoint / restore semantics under irreversible tool effects
- web-agent rollback and reset+replay methods
- reasoning rollback / backtracking methods that retain some useful failure knowledge

This means the following broad claims are not safely defensible as headline novelty:

- rollback for agents is new
- branch / restore / retract for agents is new
- rollback should be treated as runtime semantics rather than only storage
- irreversible side effects are central
- correction should not imply total amnesia

## Strongest Prior-Art Threats

- `SagaLLM` - transaction-style checkpointing, compensation, persistent state, and explicit anti-amnesia framing
- `Fork, Explore, Commit` - explicit OS primitives for branch contexts, commit/abort, and nested exploration
- `ACRFence` - checkpoint/restore under irreversible tool effects with replay-or-fork semantics
- `Toward Systems Foundations for Agentic Exploration` - explicit committed vs tentative semantics and external-side-effect pressure
- web-agent rollback/search work such as tree-search agents, WILBUR, and WebRollback

## Claim-Level Read

### Weak / already-covered claim space

- rollback as runtime semantics in the broad sense
- branch / retract / restore as explicit runtime concepts
- irreversible side effects as a first-class problem

### Narrower claim space that may still survive

- explicit separation of:
  - provisional world state
  - accepted world baseline
  - retained cognition state
- dependency-aware selective cognition retention after world revision
- typed tool-effect boundaries:
  - reversible
  - compensable
  - replayable
  - irreversible
- evaluation of revisability quality itself rather than only task success

## Repo-Level Implication

The current repo has meaningful precursor substrate:

- transactional workspace behavior
- revision-addressable workspace state
- `switch_world`
- rollback-shaped local evidence

But it does not yet have a clearly independent, paper-ready `revisable-time` primitive set.

So the current best judgment is:

- keep the dossier
- do not treat it as a broad active novelty lane yet
- only revisit promotion if a sharper temporal primitive set and a sharper evaluation story appear

## Scope Posture For Now

Current dossier posture should remain:

- real as a control-plane lane
- later than `subjective-reality`
- explicitly narrower than its manuscript rhetoric currently suggests

Until stronger primitives land, this dossier should be treated mainly as:

- a preserved future paper lane
- a place to collect temporal-semantics pressure
- not yet a broad active archival candidate

## Most Defensible Residual Gap

If this lane becomes publishable later, the strongest surviving thesis is likely to be a narrow conjunction:

> revisable execution semantics over projected world state with dependency-aware cognition retention under typed tool-effect boundaries.

That is much narrower than “rollback for agents.”
