# Manuscript Plan: Process as a File System

**Working title:** "Process as a File System: Public Runtime Surfaces for Agents"  
**Role:** middle paper between Quine and the world/time papers  
**Status:** planning stage

**Dependency note:**

- depends on Quine's process ontology
- can stand as an independent paper before Subjective Reality and Revisable Time
- likely becomes the shared operational substrate for both later papers

---

## Core Thesis

> **Agent runtime state should be reified as a synthetic filesystem, so that observation, intervention, composition, and temporal control become ordinary file operations.**

Or more compactly:

> **The agent should not only be a process you run; it should also be a directory you inspect and manipulate.**

This paper is about the **public runtime surface** of agents.

---

## What This Paper Is About

- synthetic filesystem as agent control plane
- public runtime surfaces for state, status, and control
- observation and intervention through files
- composition through filesystem structure
- making agent state Unix-native and inspectable

## What This Paper Is Not About

- world construction as the main thesis
- rollback/time semantics as the main thesis
- a generic storage format for traces

Those are later papers.

---

## Problem Statement

Quine establishes that an agent can be a native POSIX process.

But process identity alone still leaves a major systems problem:

- the agent's internal runtime state is opaque
- supervision requires bespoke APIs or logs
- intervention is awkward or destructive
- composition across agents lacks a clean public object model

So the problem becomes:

> **what should the public surface of an agent be?**

The proposed answer is:

> **a synthetic filesystem.**

---

## Main Claim Structure

### Claim 1

Agent runtime state should be public, structured, and inspectable rather than hidden in framework memory.

### Claim 2

The filesystem is the right public object model because it already supports inspection, composition, access control, and tooling.

### Claim 3

Files such as `ctl`, `status`, `memory/`, `fs/`, `diff/`, and `children/` can form a unified control plane for agents.

### Claim 4

Once agent state becomes a filesystem surface, existing Unix tooling becomes agent tooling.

### Claim 5

This public filesystem surface becomes the common substrate for later world semantics and time semantics.

---

## Paper Structure

## 1. Introduction

- Quine proved `agent = process`
- process identity is necessary but not enough
- runtime still lacks a public surface
- contribution: synthetic filesystem as public runtime surface

## 2. Background: Why Opaque Agent State Is a Systems Problem

- logs are not enough
- RPC-heavy control planes are bespoke and brittle
- hidden memory blocks inspection and intervention

## 3. Thesis: The Agent as a Synthetic Filesystem

- the agent as directory, not just PID
- public runtime surface
- filesystem as control plane

## 4. Filesystem Schema

- `ctl`
- `status`
- `mission`
- `memory/`
- `fs/`
- `diff/`
- `children/`
- possible `history/`

## 5. Operational Consequences

- inspectability
- intervention
- debugging
- composition
- operator tooling reuse

## 6. Examples

- `cat status`
- editing a memory file
- watching `diff/`
- writing commands to `ctl`
- inspecting child directories

## 7. Relation to Later Papers

- **Subjective Reality**: the world's projection becomes file-visible
- **Revisable Time**: checkpoint and restore become file-visible

This paper is the common operational substrate, not the whole story.

## 8. Evaluation Plan

- debugging efficiency vs log-only baselines
- intervention latency / simplicity
- operator task complexity
- composition and observability wins from filesystem exposure

## 9. Boundaries

- not yet a full world-construction paper
- not yet a rollback/time paper
- not a claim that every internal model state is recoverable or editable

## 10. Conclusion

- process gives identity
- process-as-fs gives publicity
- later papers add world and time semantics

---

## Publication Angle

This should be framed as a systems paper about:

- **public runtime surface**
- **publicity of mind**
- **filesystem as control plane**

not merely "saving agent state in files."

The strongest sequence is:

```text
Quine
-> Process as a File System
-> Subjective Reality
-> Revisable Time
```
