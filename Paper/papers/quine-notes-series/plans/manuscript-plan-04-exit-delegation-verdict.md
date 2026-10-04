# Quine Note 04 — exit as Delegation Verdict

**Series**: Quine Notes  
**Status**: first-wave priority  
**Target length**: 3–4 pages

---

## Core Claim

For delegated agent work, `exit status` provides the minimal **verdict channel**: a task-level outcome interface orthogonal to both deliverable (`stdout`) and explanation (`stderr`).

---

## Why This Note Exists

Delegation requires more than output. A parent must know whether the child claims completion, failure, refusal, or uncertainty. Shell composition already carries a minimal answer to this need.

This note reframes `exit` as a cognitive delegation contract rather than an implementation detail.

---

## Minimal Outline

### 1. Framework Paragraph

Reintroduce the five-channel contract and isolate outcome as separate from content.

### 2. Why Output Is Not Enough

- a child may emit text and still fail the task
- supervisors need a machine-checkable outcome channel
- application-layer agent frameworks often invent explicit status envelopes for this

### 3. exit as Verdict Interface

- `0` = success path
- non-zero = failure path
- optional conventions for richer internal taxonomy without weakening the minimal contract

### 4. Shell-Native Delegation

- `&&`, `||`, `if`, `wait` as cognitive control flow
- how a parent can route next steps from verdict alone
- interplay with `stderr` for interpretation

### 5. Structural Consequences

- simpler parent-child contracts
- easier substitution into existing automation
- clearer separation of result, rationale, and decision

### 6. Boundaries

- binary verdicts are intentionally narrow
- rich task semantics may still need typed output on `stdout`
- note should avoid overclaiming that all delegation reduces to `0/1`

### 7. Conclusion

The important point is not that exit codes are expressive enough for everything, but that they are sufficient for a robust minimal delegation layer.

---

## Likely Figures / Tables

1. deliverable / diagnostics / verdict tripartite table
2. shell control-flow snippets driven by exit status
