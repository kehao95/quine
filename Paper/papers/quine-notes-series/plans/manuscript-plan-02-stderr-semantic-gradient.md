# Quine Note 02 — stderr as Semantic Gradient

**Series**: Quine Notes  
**Status**: first-wave priority  
**Target length**: 3–5 pages

---

## Working Title Options

### Preferred

**"stderr as Semantic Gradient: Diagnostic Feedback in Process-Native Agents"**

### Alternatives

1. **"Quine Note: stderr as Semantic Gradient"**
2. **"Diagnostic Streams and Supervisory Control in Process-Native Agents"**
3. **"From Logs to Semantic Gradient: Reinterpreting stderr for Agent Runtime"**

---

## Core Claim

In a process-native agent runtime, `stderr` is not merely a log stream. It is a **semantic gradient**: a continuous diagnostics channel that lets supervisors observe partial failure, uncertainty, and recovery state before the final `exit` verdict arrives.

## One-Sentence Thesis

The key move is to treat diagnostics as an explicit **supervisory interface**, not as incidental textual exhaust.

---

## Draft Abstract (~140–180 words)

> In classical Unix process semantics, `stdout`, `stderr`, and `exit status` already form a three-part communication pattern: result, diagnostics, and verdict. Contemporary agent frameworks often reconstruct similar distinctions through custom event schemas, traces, or callback objects. Quine inherits a simpler design point by realizing agents as native POSIX processes. This note argues that, in such a runtime, `stderr` should be understood as a semantic gradient rather than a mere log sink. Unlike `exit status`, which provides a discrete terminal judgment, `stderr` can expose partial failure, uncertainty, progress, local diagnosis, and recovery intent while execution is still underway. That distinction matters because supervision is rarely binary: parents must decide whether to retry, repair, escalate, ignore, or continue waiting. The note does not claim that all failure states are narratable, nor that free-form diagnostics are sufficient for robust automation by themselves. The narrower claim is that process-native agent runtimes inherit a standard and underappreciated channel for graded supervisory feedback, reducing the need to invent parallel observability protocols above the operating system.

---

## Why This Note Exists

If an agent process exposes only `stdout` and `exit`, supervision becomes binary and coarse. Quine inherits a richer pattern:

- `stdout` = deliverable
- `stderr` = diagnostics in progress
- `exit` = terminal verdict

This creates a graded control loop rather than a single success/failure bit.

It also lets Quine say something stronger than “the agent logs what happened.”

The sharper claim is:

- `stdout` answers: what was produced?
- `stderr` answers: what is going wrong, being checked, or being repaired?
- `exit` answers: what verdict should the parent use for control flow?

---

## Intended Contribution Shape

This is a **supervision-semantics note**.

It should read as:

- stronger than a logging convention
- narrower than a full observability paper
- rooted in process interface semantics rather than orchestration middleware design

---

## Argument Spine

### A. Agent supervision needs more than terminal success/failure

Parents often need to know:

- whether a child is making progress
- what kind of failure is happening
- whether repair is plausible
- whether escalation is warranted

### B. POSIX already provides a graded interface

- `stdout` carries deliverable content
- `stderr` carries diagnostics
- `exit` closes with a verdict

### C. In agent runtime, `stderr` becomes semantically important

Because the child is cognitive, diagnostics are not only stack traces or warnings. They can encode search state, uncertainty, failure mode, or recovery attempt.

### D. This changes supervision design

Parents can route behavior based on the gradient before or alongside the verdict.

### E. But the claim remains disciplined

The note should not claim that free-form diagnostics automatically yield reliable machine supervision.

---

## Expanded Outline

### 1. Framework Paragraph

Reintroduce the five-channel contract and isolate the difference between deliverable, diagnostics, and outcome.

**Target content**:

- Quine exposes deliverable, diagnostics, and outcome as separate channels
- this note focuses on the supervisory meaning of `stderr`
- the central contrast is continuous diagnostic flow vs discrete terminal verdict

**Possible opening sentence**:

> Once an agent is realized as a process, its feedback to supervisors is already stratified by the operating system.

### 2. Why Binary Outcome Is Not Enough

- delegation often needs progress-signals before completion
- retry, escalation, and repair depend on seeing failure modes, not only end states
- framework runtimes often reconstruct this channel through custom event objects

**Subsections**:

#### 2.1 Delegation Is Rarely Purely Binary

- some failures are temporary
- some are terminal
- some indicate wrong decomposition rather than wrong execution

#### 2.2 Verdict Without Context Is Weak

- non-zero alone does not explain whether to retry, rewrite, or escalate
- binary termination is necessary but insufficient for supervision

#### 2.3 Framework Reconstruction

- many systems invent trace events, callbacks, or monitoring objects
- Quine can inherit an OS-native distinction instead

### 3. stderr as Continuous Diagnostic Surface

- stream semantics rather than terminal semantics
- human-readable and machine-usable diagnostics
- distinction from ordinary verbose logging

**Subsections**:

#### 3.1 Why Stream Semantics Matter

- diagnostics can be emitted before completion
- parent can observe evolving state, not just postmortem result

#### 3.2 What Makes Diagnostics Semantic

- not just noise or debug text
- statements about uncertainty, blockage, search progress, or repair strategy
- can be shaped by runtime policy and prompting

#### 3.3 Diagnostic Channel vs Logging Channel

- logs are often archival and developer-facing
- this note treats `stderr` as part of the agent's external contract
- the distinction is supervisory, not merely operational

### 4. exit as Discrete Verdict

- `exit` closes the task
- `stderr` shapes interpretation before closure
- together they form gradient + verdict

**Subsections**:

#### 4.1 Why Verdict Still Matters

- final machine-checkable branch condition
- shell-native compatibility
- minimal delegation contract remains intact

#### 4.2 Gradient + Verdict

- `stderr` informs interpretation
- `exit` determines closure
- one is graded, one is discrete

#### 4.3 Separation of Roles

- `stdout` should not be polluted with supervision chatter
- `stderr` carries diagnostics without breaking pipeline deliverables

### 5. Supervisory Patterns

- retry with altered mission
- fallback to another child
- escalate to human or stronger model
- aggregate child diagnostics in fan-out runs

**Candidate patterns**:

#### Pattern A: Retry with repair hint

Child reports on `stderr` that a file path is ambiguous or an assumption failed; parent reissues a narrower mission.

#### Pattern B: Fallback routing

Child exits non-zero and `stderr` indicates capability mismatch; parent dispatches to a stronger or more specialized child.

#### Pattern C: Fan-out and synthesis

Multiple children emit diagnostic streams; parent aggregates them to identify shared blockage or disagreement.

#### Pattern D: Human escalation

Parent forwards deliverable candidate plus diagnostic context for review when verdict is uncertain or repeated attempts fail.

---

## Proposed Section-Level Skeleton

### 1. Introduction

- supervision problem
- why binary verdicts are insufficient
- note's central claim

### 2. The Three-Part Output Contract

- `stdout` as deliverable
- `stderr` as diagnostics
- `exit` as outcome

### 3. stderr as Semantic Gradient

- define the term carefully
- explain why continuous diagnostics matter
- distinguish from logs

### 4. Supervisory Consequences

- retry
- repair
- escalation
- aggregation across children

### 5. Boundaries

- diagnostics can be low quality
- not every failure is narratable
- structured supervision may still want more formal schemas later

### 6. Conclusion

- reuse OS-native diagnostics channel
- retain minimalism of process interface

## Key Boundary Reminder

- `stderr` quality depends on agent discipline and runtime shaping
- not every failure is narratable
- continuous diagnostics do not replace objective metrics

**Boundary sentence to preserve**:

> The claim is not that stderr alone yields reliable automated supervision, but that it provides a standard graded feedback channel that sits naturally between deliverable and verdict.

## Closing Emphasis

The primitive matters because it lets supervision reuse a standard OS channel instead of inventing a parallel observability protocol.

---

## Related Work / Citation Pointers

- Quine main paper for the five-channel contract
- optionally cite one or two agent-framework examples only to illustrate that many systems reconstruct richer event layers above the OS
- avoid turning this note into a full observability literature review

---

## Writing Risks

1. **Collapsing into “stderr is just logs”**  
   Fix: emphasize supervisory semantics and relation to `exit`.

2. **Overclaiming automation reliability**  
   Fix: explicitly separate availability of signal from quality of interpretation.

3. **Letting the note become about generic observability**  
   Fix: stay focused on process interface roles.

4. **Ignoring `stdout` / `exit` contrast**  
   Fix: keep the tripartite distinction central throughout.

---

## Likely Figures / Tables

1. `stdout` / `stderr` / `exit` separation diagram
2. supervisor decision table: retry / skip / escalate / abort

### Optional third figure

3. timeline view: child emits diagnostics over time, then terminates with final verdict

---

## Suggested Example Snippets

### Example 1: clean deliverable vs diagnostics

```bash
./quine "Summarize the incident and exit 0 only if root cause is identified" < incident.md \
  > summary.txt \
  2> diagnosis.log
```

### Example 2: verdict-driven branching with diagnostic review

```bash
if ./quine "Apply patch and exit 0 only if tests pass" < task.md 2>repair.log; then
  ./quine "Prepare commit message" < task.md
else
  ./quine "Review failure diagnostics and propose next step" < repair.log
fi
```

---

## Suggested Next Step

Write this note in the following order:

1. Section 2 binary-outcome problem
2. Section 3 define “semantic gradient” carefully
3. Section 4 gradient + verdict distinction
4. abstract
5. example snippets and boundaries last
