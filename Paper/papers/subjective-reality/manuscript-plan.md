# Manuscript Plan: Subjective Reality

**Working title:** "Subjective Reality: Constructed Worlds for Agent Runtime"  
**Role:** world-semantics paper downstream of Process as FS  
**Status:** planning stage

**Dependency note:**

- the minimal workspace-centered implementation can start directly from Quine
- the broader paper-level version likely sits **after Process as a File System**, because richer constructed worlds become clearer once projections, diffs, and retained structures are themselves file-visible
- this paper should stay separate from **Revisable Time**

---

## Core Thesis

> **Agent reality should be constructed by the runtime, not merely inherited from the host.**

More concretely:

> **What is present in the runtime's projection is what is present in the agent's world.**

This paper is about **world semantics**, not rollback semantics.

---

## What This Paper Is About

- constructed worlds for agents
- filesystem projection as the primary operational medium
- curated tool surfaces and safe bins
- network / capability scoping as part of world construction
- mutation visibility as part of world definition
- child-specific world pruning

## What This Paper Is Not About

- full rollback theory
- branch/retract semantics
- undo without amnesia
- irreversible side-effect recovery after acceptance

Those belong to **Revisable Time**.

---

## Problem Statement

Quine establishes that an agent can be realized as a native POSIX process.
But a process alone still inherits an ill-defined world:

- too many files
- too many tools
- too much ambient authority
- poor observability of physical change

This yields:

- pollution
- blindness
- misdirected competence

The problem is not only sandboxing.
It is the absence of task-relative world construction.

---

## Main Claim Structure

### Claim 1

Inherited host reality is the wrong default abstraction for agents.

### Claim 2

Filesystem projection is the cleanest first mechanism for constructing agent worlds.

### Claim 3

Workspace projection is only the first instantiation of a broader constructed-world runtime.

### Claim 4

Tool subsets, safe bins, and network scoping are also forms of world construction, not just policy hardening.

### Claim 5

Observable mutation channels such as `[FS MUTATIONS]` are part of the world's semantics because they define what the agent can know about action.

---

## Implementation Frame

### First Instantiation: Workspace-Centered World

- `QUINE_WORKSPACE`
- `workspace root`
- `workspace session`
- `workspace diff`

This provides:

- a stable writable world
- explicit commit/discard semantics
- a concrete place to attach mutation visibility

### Stronger Future Instantiations

- synthetic projected roots
- child-specific pruned worlds
- curated executable surfaces
- network-scoped worlds
- capability-specific evaluation sandboxes

---

## Paper Structure

## 1. Introduction

- Quine proved `agent = process`
- process is necessary but not sufficient
- the runtime must also define the world
- contribution: constructed worlds for agent runtime

## 2. Background: From Process to World

- recap Quine mapping
- why process identity does not solve world ambiguity
- sandbox vs world-construction distinction

## 3. Thesis: Reality Should Be Constructed

- formal statement of constructed-world claim
- task-relative projection
- semantic honesty over fake parity

## 4. Filesystem Projection as World Construction

- projected filesystem
- stable world root
- child-specific narrowing
- mutation visibility

## 5. Workspace as First Instantiation

- why workspace first
- recursive workspace inheritance
- commit surface vs execution surface
- macOS and Linux as different substrate depths under one contract

## 6. Beyond Workspace

- safe bins
- subset of tools
- capability scoping
- network scoping
- agent-constructed child worlds

## 7. Evaluation Plan

- wrong-path edit reduction
- tool-selection accuracy
- navigation reduction
- mutation-awareness gains from `[FS MUTATIONS]`
- world-pruning effects on task success

## 8. Boundaries

- not a time paper
- not a rollback paper
- not a claim that all platforms share identical substrate

## 9. Conclusion

- process gives identity
- subjective reality gives world

---

## Evaluation Questions

1. Does constructed world definition reduce blind edits?
2. Does curated tool visibility improve tool choice?
3. Does mutation visibility reduce false beliefs after "successful" commands?
4. Do child-specific pruned worlds improve delegated task efficiency?

---

## Publication Angle

This should read as a **systems/runtime paper**, not a workspace feature note.

The framing should be:

- theoretical layer: constructed worlds
- implementation layer: workspace projection as the first serious slice

That keeps the current implementation meaningful without collapsing the paper into a product feature.
