# Quine Note 05 — stdin/stdout as the Compositional Interface

**Series**: Quine Notes  
**Status**: second-wave priority  
**Target length**: 4–5 pages

---

## Core Claim

`stdin` and `stdout` form the minimal **compositional membrane** of process-native agents: if that contract is preserved, agents become substitutable filters in pipelines, chains, and larger process graphs.

---

## Why This Note Exists

Quine's process model is not only about isolation. It is also about composition. The central question is not symmetry with other channels, but why pipelines remain possible when cognition enters the loop.

---

## Minimal Outline

### 1. Framework Paragraph

Restate the five-channel interface, then focus specifically on material-in / deliverable-out as the composition surface visible to peers.

### 2. Composition as Interface Preservation

- why composition depends on stable external contract rather than shared internal representation
- POSIX filters as the historical model
- Quine as cognitive filter

### 3. Contract Dimensions

- representation: text / binary
- temporal mode: stream / batch
- framing and granularity
- human-legible vs machine-legible output

### 4. Structural Consequences

- drop-in substitution in shell pipelines
- implicit DAGs without orchestrator state
- delegation and chaining through shared membrane

### 5. Examples

- semantic grep replacement
- multi-stage reasoning chain
- parent-child substitution without API redesign

### 6. Boundaries

- not every cognitive task fits stream composition cleanly
- framing mismatches can break interoperability
- stdout discipline is a real requirement, not a free consequence

### 7. Conclusion

This note should argue that composability is inherited from interface discipline, not added by orchestration software.

---

## Likely Figures / Tables

1. contract-dimensions table
2. pipeline diagram with Quine instances and classic Unix tools
