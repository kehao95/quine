# Quine Note 01 — Cognitive Harvard Architecture

**Series**: Quine Notes  
**Status**: first-wave priority  
**Target length**: 3–5 pages

---

## Working Title Options

### Preferred

**"Cognitive Harvard Architecture: Instruction–Material Separation in Process-Native Agents"**

### Alternatives

1. **"Quine Note: A Cognitive Harvard Architecture for Agent Runtime"**
2. **"Instruction and Material as Separate Channels in Process-Native Agents"**
3. **"argv and stdin as a Cognitive Harvard Architecture"**

---

## Core Claim

Quine's `argv` / `stdin` split realizes a **cognitive Harvard architecture**: instructions and material enter through structurally distinct channels, so agent control and agent data are separated at the runtime boundary rather than only in prompt formatting.

## One-Sentence Thesis

The important move is to treat instruction–data separation as a **runtime property** of the agent boundary, not merely as a prompt-writing convention.

---

## Draft Abstract (~140–180 words)

> Contemporary agent systems often flatten mission, material, retrieved context, and tool output into a single prompt surface. This makes instruction priority largely conventional: the boundary between what the agent is told to do and what the agent is asked to inspect exists mainly in formatting. Quine exposes a sharper design point by realizing agents as native POSIX processes. In this setting, mission enters through `argv`, while material enters through `stdin`. This note argues that the split is not cosmetic but architectural. It realizes a cognitive analogue of Harvard architecture: control and material cross the runtime boundary through distinct channels with different intended semantics. The consequence is a cleaner vocabulary for delegation, substitution, replay, and supervision. The note does not claim complete security against prompt manipulation, nor that OS-level channel separation survives unchanged inside model context assembly. The narrower claim is that process-native agents can inherit a principled instruction–material boundary from the operating system itself, giving agent runtime a more stable input semantics than flattened application-layer prompts.

---

## Why This Note Exists

Most agent systems flatten instructions, context, and tool output into one prompt surface. Quine instead inherits an existing systems distinction:

- `argv` = mission / instruction
- `stdin` = material / data

This note argues that the distinction is architectural, not cosmetic.

More specifically, it gives you a way to say:

- what the agent is here to do
- what the agent is here to inspect
- why those are different categories
- why that difference should exist before prompt assembly, not only inside it

---

## Intended Contribution Shape

This is a **sharp interface-semantics note**.

It should be:

- more precise than a metaphor paper
- lighter than a full security paper
- framed as runtime semantics, not prompt-engineering advice

---

## Argument Spine

### A. Existing agent runtimes usually flatten instruction and material

The dominant prompt surface mixes:

- mission
- input data
- retrieved context
- prior tool outputs
- formatting instructions

This makes the distinction between control and material weakly represented.

### B. Quine inherits a stronger boundary from POSIX

At the process interface:

- `argv` is the mission channel
- `stdin` is the material channel

The split precedes prompt serialization.

### C. This makes instruction–material separation structural

The distinction is not “please put the task above the context.” It is “these inputs arrive through different runtime channels with different roles.”

### D. Structural separation has downstream consequences

- delegation can vary mission without rewriting material
- replay can hold material constant while swapping mission
- composition can preserve pipeline data while changing task intent
- safety reasoning gains a cleaner boundary condition

### E. But the claim remains narrow

This does **not** prove:

- robust prompt-injection immunity
- total preservation of separation after model-side serialization
- equivalence to hardware Harvard architectures beyond analogy

---

## Expanded Outline

### 1. Framework Paragraph

Restate the Quine process mapping briefly, then specify that this note isolates the input-side semantics of `argv` and `stdin`.

**Target content**:

- Quine models agents as native POSIX processes
- mission arrives via `argv`
- material arrives via `stdin`
- output, diagnostics, and verdict travel separately
- this note focuses only on the semantics of the two input channels

**Possible opening sentence**:

> Quine models an agent as a native POSIX process, and that choice makes the input boundary itself semantically structured.

### 2. The Flattening Problem

- application-layer agents often collapse instruction and material into one buffer
- when everything is one blob, priority is conventional rather than structural
- the runtime loses a clean boundary between mission and evidence

**Subsections**:

#### 2.1 Prompt Monism

- many systems have one main prompt surface
- distinctions are represented by templates, delimiters, or role labels
- useful, but still application-defined

#### 2.2 What Gets Lost

- ambiguous status of retrieved text
- weaker distinction between mission and evidence
- harder to reason about substitution and replay

#### 2.3 Why This Matters

- the system cannot clearly express “same task, new material” vs “same material, new task”
- delegation contracts stay muddy

### 3. The POSIX Separation

- `argv` as immutable mission descriptor
- `stdin` as material stream
- host maps them to privileged instruction vs task material
- relation to instruction hierarchy without overclaiming security guarantees

**Subsections**:

#### 3.1 `argv` as Mission Channel

- compact, intentional, task-defining input
- stable across one process invocation
- naturally suited to “what must be done”

#### 3.2 `stdin` as Material Channel

- potentially large, streamed, task-relative input
- naturally suited to “what must be examined”
- pipeline-compatible by construction

#### 3.3 Host Mapping

- runtime maps mission into privileged instruction layer
- runtime maps material into task content layer
- resulting distinction is inherited from interface structure, not merely formatting style

#### 3.4 Relation to Instruction Hierarchy

- complementary to model-level instruction hierarchy work
- provides a structural boundary condition rather than a full defense guarantee

### 4. Structural Consequences

- cleaner mission/material vocabulary
- better delegation contracts
- easier substitution and replay
- stronger basis for reasoning about prompt manipulation boundaries

**Subsections**:

#### 4.1 Delegation Clarity

- parent can keep material fixed while changing mission
- or preserve mission while streaming different material batches

#### 4.2 Reproducibility and Replay

- experiments can vary one channel at a time
- note-level examples become clearer and more falsifiable

#### 4.3 Compositional Stability

- stdin remains pipeline material even as mission changes
- supports cognitive filters with stable external interfaces

#### 4.4 Safety Framing

- cleaner reasoning about where manipulation pressure enters
- not a full security result, but a better runtime vocabulary for discussing one

### 5. Illustrative Patterns

- fixed mission, changing material
- same material, different mission
- parent delegates by rewriting `argv` while preserving pipeline input

**Candidate examples**:

#### Example A: Fixed mission, changing material

```bash
cat log-a.txt | ./quine "Extract auth failures"
cat log-b.txt | ./quine "Extract auth failures"
```

#### Example B: Same material, different mission

```bash
cat incident.md | ./quine "Summarize the failure mechanism"
cat incident.md | ./quine "List operator action items"
```

#### Example C: Delegation with mission rewrite only

Parent consumes one stream, then spawns children with different `argv` missions over the same staged material.

#### Example D: Evaluation logic

Hold `stdin` constant while testing whether alternate missions produce consistent behavior.

---

## Proposed Section-Level Skeleton

### 1. Introduction

- one paragraph: field tendency to flatten inputs
- one paragraph: Quine's sharper design point
- one paragraph: contribution and scope

### 2. From Flattened Prompts to Structured Process Input

- define the problem crisply
- explain why formatting alone is weaker than interface separation

### 3. A Cognitive Harvard Architecture

- define the analogy
- map `argv` ↔ control, `stdin` ↔ material
- specify what the analogy does and does not mean

### 4. Structural Consequences

- delegation
- replay
- composition
- safety framing

### 5. Boundaries

- no full security claim
- no claim that serialization preserves purity absolutely
- analogy is architectural, not literal hardware equivalence

### 6. Conclusion

- restate the narrow contribution
- connect back to Quine Notes series

## Key Boundary Reminder

- not a full security proof
- not equivalent to CPU Harvard architecture in all respects
- does not eliminate all channel contamination once material is reserialized into model context

**Boundary sentence to preserve**:

> The claim is not that Quine makes prompt injection disappear, but that it gives the agent runtime a cleaner and more principled instruction–material boundary than flattened prompt architectures.

## Closing Emphasis

The important contribution is not the metaphor itself, but the claim that input stratification should be treated as runtime structure.

---

## Related Work / Citation Pointers

- Quine main paper for the five-channel mapping
- instruction hierarchy work as supporting context for privileged instruction layers
- optionally one brief nod to Harvard architecture only if it sharpens the analogy rather than overburdening the note

---

## Writing Risks

1. **Too metaphor-heavy**  
   Fix: keep returning to the runtime boundary and concrete channels.

2. **Drifting into security-paper claims**  
   Fix: say “structural basis” or “cleaner boundary,” not “solves prompt injection.”

3. **Re-explaining all of Quine**  
   Fix: use one tight framework paragraph, then move immediately to the specific primitive.

4. **Overstating the Harvard analogy**  
   Fix: present it as a cognitive-runtime analogy, not hardware identity.

---

## Likely Figures / Tables

1. `argv` / `stdin` split mapped onto instruction / material
2. comparison table: flattened prompt vs structured process input

### Optional third figure

3. “same material / different mission” and “same mission / different material” matrix

---

## Dependencies

- cites the main Quine paper for the five-channel mapping
- optionally cites instruction hierarchy work for supporting context

---

## Suggested Next Step

Write this note in the following order:

1. framework paragraph
2. Section 2 flattening problem
3. Section 3 core mapping
4. 150-word abstract
5. examples and boundaries last
