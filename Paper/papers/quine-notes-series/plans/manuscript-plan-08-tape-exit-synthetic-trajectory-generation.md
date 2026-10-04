# Quine Note 08 — Tape and Exit for Synthetic Trajectory Generation

**Series**: Quine Notes  
**Status**: deferred / experiment-sensitive  
**Target length**: 4–6 pages

---

## Core Claim

Execution tape paired with `exit status` yields a natural unit of **synthetic trajectory data**: process traces can be filtered by outcome, inspected by diagnostics, and reused as candidate learning material.

---

## Why This Note Is Deferred

This idea is strong but not note-complete on theory alone. It quickly raises empirical questions about:

- label quality
- filtering validity
- learning usefulness
- data contamination and bias

So this note should stay later in the sequence unless paired with real experiments.

---

## Minimal Outline

### 1. Framework Paragraph

Restate Quine's tape, diagnostics, and outcome channels. Isolate the runtime record produced by one process execution.

### 2. The Data Opportunity

- every run yields a trajectory
- `exit` supplies a weak but useful task-level label
- `stderr` helps interpret failure modes

### 3. Candidate Data Unit

- mission
- material
- tool trace / tape
- deliverable
- diagnostics
- outcome

### 4. Filtering and Use Cases

- keep only successful or partially recoverable runs
- mine common failure trajectories
- bootstrap repair or critique datasets

### 5. What Must Be Proven Empirically

- whether verdict correlates with true task quality
- whether traces improve downstream agents
- how much manual filtering is still required

### 6. Boundaries

- synthetic does not mean valid
- success labels can be gamed or shallow
- data-generation claims should remain narrow until measured

### 7. Conclusion

Frame this note as a bridge from runtime semantics to learning systems, but keep it explicitly conditional on experimental validation.

---

## Likely Figures / Tables

1. trajectory record schema
2. filtering pipeline from raw runs to candidate training data
