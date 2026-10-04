# Quine Notes — Series Manuscript Plan

**Series role**: extend the Quine systems paper into a family of compact notes on interface semantics, composition, and coordination.  
**Default length**: 3–5 pages each.  
**Status**: planned.

---

## Series Thesis

Once an agent is realized as a native POSIX process, familiar process primitives stop being mere implementation details and become **cognitive runtime semantics**.

The notes isolate those semantics one by one.

---

## Organizing Logic

### A. Interface Semantics

These notes explain what the core process interface already means:

1. **Cognitive Harvard Architecture** — instruction and material separation
2. **stderr as Semantic Gradient** — graded supervision feedback
3. **exec as Context Reset** — hard re-instantiation without identity loss
4. **exit as Delegation Verdict** — minimal outcome contract for delegated work

### B. Composition and Coordination

These notes explain how agent processes combine:

5. **stdin/stdout as the Compositional Interface** — substitutable pipeline membrane
6. **Filesystem as Information Medium** — durable shared coordination surface
7. **Fork, Wait, and Race in Agent Runtime** — native parallel exploration and supervision

### C. Data / Learning

These notes explain what the runtime records make possible:

8. **Tape and Exit for Synthetic Trajectory Generation** — execution traces as filterable learning data

---

## Recommended Publication Sequence

### Phase 1 — sharp primitive notes

1. Cognitive Harvard Architecture
2. stderr as Semantic Gradient
3. exec as Context Reset
4. exit as Delegation Verdict

### Phase 2 — composition notes

5. stdin/stdout as the Compositional Interface
6. Filesystem as Information Medium
7. Fork, Wait, and Race in Agent Runtime

### Phase 3 — experiment-sensitive note

8. Tape and Exit for Synthetic Trajectory Generation

---

## Shared Note Skeleton

Each note should, in minimal form, contain:

1. **Framework paragraph**
   - restate the Quine mapping in 4–6 sentences
   - specify which primitive this note isolates

2. **Problem statement**
   - what disappears when the primitive is treated as mere plumbing

3. **Core claim**
   - one sentence, technically sharp, venue-safe

4. **Mechanism**
   - how the POSIX primitive behaves at runtime

5. **Structural consequence**
   - what new property follows for agent systems

6. **Illustrative examples**
   - shell compositions, supervision loops, or runtime sketches

7. **Boundaries**
   - what the note does not claim

8. **Conclusion**
   - fold the primitive back into the broader Quine vocabulary

---

## Cross-Note Rules

- Avoid re-proving the whole Quine paper.
- Reuse one stable vocabulary across notes.
- Keep each note independently readable.
- Prefer mechanism → consequence over manifesto language.
- Defer benchmark-heavy claims unless a note genuinely requires them.

---

## Common Framework Paragraph Template

> Quine models an agent as a native POSIX process. The mission is passed through `argv`, material through `stdin`, deliverables through `stdout`, diagnostics through `stderr`, and final outcome through `exit status`; lifecycle is expressed through `fork`, `exec`, and `exit`. The host runtime maps this contract onto an LLM supervision loop while preserving the external process interface. This note isolates one consequence of that mapping: <primitive-specific claim>.

---

## Near-Term Deliverables

- complete Wave 1 manuscript plans
- draft abstracts for notes 1 and 2 first
- keep notes 6–8 in queued state until the mother ship and first wave are stable
