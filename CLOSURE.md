# Closure Protocol

This document defines the protocol for closing development intervals and absorbing work into durable repository state.

## Purpose

Closure makes work available to the next participant: implemented changes,
evidence, decisions, discussions, and unresolved questions. The retained form
should support continuation at the level the work needs.

Active development naturally produces exploratory scratchpads, temporary scripts, and messy intermediate residue. Closure does not demand the paralyzing triage of every ephemeral scratch file; rather, it ensures that the *durable signal* is captured and that the repository's ground-truth state remains coherent.

## When To Run

- When explicitly requested (e.g. `Closure` command).
- At the culmination of a feature, bug fix, or refactor.
- When an experiment campaign, model evaluation, or theoretical milestone is reached.
- Before handing off work across agent boundaries or ending a major line of development.

## Closure Criteria

A work interval is ready for closure when the following conditions are met:

1. **Substantive Work Absorbed**:
   - Implementation changes are tested, functional, and checked into their canonical code locations.
   - Any modifications to public runtime physics, tool schemas, or environment controls are reflected in their owning control-plane documentation (e.g. [`DEVELOPMENT.md`](./DEVELOPMENT.md), [`Paper/core/registries/env-controls.md`](./Paper/core/registries/env-controls.md)).

2. **Durable Evidence & Baselines Preserved**:
   - Curated proof runs, benchmark baselines, or key diagnostic evidence are recorded.
   - For bulky experiment payloads, use DVC tracking (`.dvc` pointers in git) per [`experiments/PROTOCOL.md`](./experiments/PROTOCOL.md); discard redundant or non-claim raw runs.
   - Active model evaluations retain their passing canonical baselines under `tests/model/`.

3. **Theory & Research Alignment**:
   - When scientific or architectural findings emerge, update the relevant theory objects, gaps, or live seeds in [`Paper/theory/`](./Paper/theory/README.md).
   - If active experiments were completed or updated, confirm reciprocity with their associated theory claims.
   - For conceptual work, retain the developed account and consequential alternatives; record decisions where made and update affected routes.
   - Run checks appropriate to the change per [EVALUATION.md](./EVALUATION.md); do not turn a documentation interval into an unrelated runtime evaluation campaign.

4. **Observation & Status Updated**:
   - If the interval alters the active roadmap, current blockers, or next transitions, update [`development/status/README.md`](./development/status/README.md).
   - Keep status entries concise and resume-ready for future agents.

## Residue & Scratchpad Policy

- **Useful work stays findable**: Retain implementation, evidence, discussions,
  and open alternatives with enough context to continue them. A research record
  need not become a single settled account or a standing instruction.
- **Ephemeral Scratch Does Not Block**: Disposable debug logs, failed local
  attempts, and scratch files with no retained information value need no
  item-by-item classification and may be discarded at will.
- **No Unintended Commit Bloat**: Ensure untracked temporary scratch files or massive uncurated raw directories are not accidentally committed to git.
- **Explicit Handoff**: If meaningful work or unmerged prototypes are intentionally left in-progress, state the open context clearly in the closure report or status log.

## Output Contract

State the audit outcome clearly:

- `Closure complete` — when core changes, evidence, and docs are absorbed into durable state.
- `Closure incomplete` — when critical validation is pending, blocking errors remain, or state transitions are unresolved.

Summarize what changed, the supporting validation/evidence, and any next steps.
