# Quine Note 06 — Filesystem as Information Medium

**Series**: Quine Notes  
**Status**: second-wave priority  
**Target length**: 4–5 pages

---

## Core Claim

In process-native agent systems, the filesystem is not only persistence. It is an **information medium** for durable coordination, artifact exchange, and shared situational state across processes and time.

---

## Why This Note Exists

Streams are excellent for immediate flow, but many agent interactions need:

- durable artifacts
- asynchronous handoff
- inspectable intermediate state
- shared but loosely coupled coordination

The filesystem already provides these semantics.

---

## Minimal Outline

### 1. Framework Paragraph

Restate the Quine state model: memory, environment, filesystem. Isolate the persistent shared layer.

### 2. Beyond Streams

- why some coordination cannot live purely in pipes
- persistence, naming, revisitation, and asynchronous access

### 3. Filesystem Semantics as Coordination Primitives

- paths as references
- files as artifacts
- directories as shared task spaces
- append, overwrite, atomic rename, and lockfile-style patterns

### 4. Structural Consequences

- parent-child handoff through artifacts
- auditability and replay
- loose coupling between agents that need not share a direct pipe

### 5. Examples

- work queue via directory convention
- shared scratchpad or patch bundle
- postmortem through saved artifacts and diagnostics

### 6. Boundaries

- consistency and cleanup become real concerns
- shared filesystem can reintroduce coupling and contention
- this note should not overclaim distributed consensus properties

### 7. Conclusion

The note should position filesystem semantics as the durable counterpart to stream composition.

---

## Likely Figures / Tables

1. streams vs filesystem comparison table
2. directory-based coordination sketch
