# Manuscript Plan: Revisable Time

**Working title:** "Revisable Time: Rollback Without Amnesia in Agent Runtime"  
**Role:** time-semantics paper downstream of Subjective Reality and Process as FS  
**Status:** planning stage

**Dependency note:**

- depends on Quine's process ontology
- depends on **Subjective Reality** for provisional world semantics
- likely sits after **Process as a File System** for the strongest operational form, because checkpoint, restore, and branch become much cleaner once mind and world are projected as files
- the likely paper-level chain is `Quine -> Process as a File System -> Subjective Reality -> Revisable Time`, even if some workspace-centered engineering slices land earlier

---

## Core Thesis

> **Agent runtime should give actions revisable temporal semantics rather than treating side-effects as immediately final reality.**

Or more sharply:

> **Rollback should revise world-state without necessarily erasing cognition.**

This paper is about **time semantics**, not world construction.

---

## What This Paper Is About

- rollback without amnesia
- branch / revise / retract
- accepted state vs provisional state
- temporal boundaries for operations
- irreversible side-effects as a runtime problem
- mind/world synchronization under revision

## What This Paper Is Not About

- world construction as such
- curated tools as such
- safe-bin or network scoping as such
- workspace projection as the central object

Those belong to **Subjective Reality**.

---

## Problem Statement

Even in a well-constructed world, time is still broken by default.

Why?

- commands execute irreversibly
- failed trajectories leave residue
- context can be reset, but side-effects often cannot
- the runtime lacks a clean boundary between provisional and accepted action

This yields the central failure:

> the agent's cognition and the world's state diverge across failed time.

Examples:

- destructive probes of fragile data
- partial refactors that leave a broken tree
- useful reasoning embedded inside failed action sequences

---

## Main Claim Structure

### Claim 1

Rollback is not only a storage feature; it is a temporal semantic of runtime.

### Claim 2

The correct target is not "undo everything blindly" but "revise action history while preserving selected cognition."

### Claim 3

Operation boundaries should become explicit temporal boundaries.

### Claim 4

Pending diff state, accepted base state, and retained reasoning state must be modeled separately.

### Claim 5

Irreversible side-effects are not edge cases; they define the need for revisable time.

---

## Relation to Subjective Reality

The clean split is:

- **Subjective Reality** asks: what world is present?
- **Revisable Time** asks: how can actions on that world be revised?

Interfaces between the two:

- projected diff state
- operation checkpoints
- accepted base vs provisional layer
- retained runtime history

But the papers should stay separate because the evaluation targets differ:

- world paper: blindness reduction
- time paper: irreversibility reduction

---

## Paper Structure

## 1. Introduction

- even a constructed world still suffers from irreversible time
- agent runtime needs revisable execution semantics
- contributions: rollback without amnesia, explicit temporal boundaries, runtime separation of cognition vs side-effects

## 2. Background: Why Time Is Broken

- irreversible side-effects in ordinary process execution
- context reset vs world reset mismatch
- failure residue and branch collapse

## 3. Thesis: Time Should Be Revisable

- provisional vs accepted action
- retained cognition vs revised world state
- temporal semantics rather than ad hoc cleanup

## 4. Temporal Primitives

- operation boundary
- checkpoint
- pending diff
- accepted base
- retained history

## 5. Rollback Without Amnesia

- why naive rollback is insufficient
- selective cognition retention
- world rollback with explicit surviving state

## 6. Branch / Retract / Revise

- branching failed attempts
- retracting a provisional line
- revising one branch without contaminating accepted state

## 7. Irreversible Side Effects

- dangerous probes
- external effects
- when rollback cannot restore the external world
- what the runtime can still do semantically

## 8. Evaluation Plan

- recovery from failed action sequences
- rollback latency and overhead
- branch exploration cost
- success after destructive probe scenarios
- cognition preservation vs full reset baselines

## 9. Boundaries

- not a paper about world construction
- not a claim that every side-effect is physically reversible
- not a claim that rollback eliminates the need for careful task design

## 10. Conclusion

- process gives identity
- subjective reality gives world
- revisable time gives recoverable execution

---

## Evaluation Questions

1. Can rollback preserve useful reasoning while revising failed world state?
2. Does explicit temporal structure reduce failure residue?
3. What is the cost of checkpoints and branch maintenance?
4. How often do destructive task probes benefit from revisable time?

---

## Publication Angle

This should read as a **runtime semantics paper about time**, not as a storage-engine appendix.

The strong framing is:

- agent systems need revisable execution, not just safer file writes
- rollback is part of temporal ontology
- cognition and side-effects must be modeled separately
