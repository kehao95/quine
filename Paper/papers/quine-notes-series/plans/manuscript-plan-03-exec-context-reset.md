# Quine Note 03 — exec as Context Reset

**Series**: Quine Notes  
**Status**: first-wave priority  
**Target length**: 3–5 pages

---

## Core Claim

`exec` should be understood as a **context reset primitive** for agent runtime: it replaces the running image while preserving process identity and selected scoped state, enabling hard re-instantiation rather than incremental prompt cleanup.

---

## Why This Note Exists

Long-running agents accumulate context, drift, and residue. Many systems answer this with summarization, truncation, or ad hoc “clear context” routines. Quine has a sharper primitive already available:

- checkpoint what must survive
- `exec` into a fresh image
- resume under the same process identity

---

## Minimal Outline

### 1. Framework Paragraph

Restate Quine's lifecycle mapping and isolate the `exec` operation.

### 2. The Context Accumulation Problem

- finite context windows
- instruction drift and residue
- difference between fresh state and merely shorter prompts

### 3. What `exec` Actually Does

- image replacement, not in-place cleanup
- preserved PID / environment / inherited descriptors
- discarded address-space state

### 4. Reset as Re-instantiation

- checkpoint / exec / resume pattern
- deterministic reconstruction vs historical carryover
- why this is stronger than summarization alone

### 5. Structural Consequences

- long tasks decomposed into generations
- identity continuity without memory continuity
- clearer model of cognitive mortality and rebirth

### 6. Boundaries

- reset is only as good as the checkpoint contract
- not all useful state fits naturally into environment or filesystem
- `exec` is not a substitute for deliberate decomposition

### 7. Conclusion

The note should argue that agent continuation deserves a syscall-grade primitive, not a framework convenience method.

---

## Likely Figures / Tables

1. checkpoint → exec → resume lifecycle diagram
2. table of preserved vs discarded state across `exec`
