# Manuscript Plan: ALIFE Conference-First Paper

## Current Priority

- Primary target: ALIFE conference full paper
- Deadline to optimize for: March 30, 2026
- Primary goal: accept
- Stretch goal: best paper
- Secondary target: arXiv preprint as staging artifact, not the main framing constraint

This document replaces the old "Track A / Track B" split with a conference-first manuscript plan. The objective is not to say everything Quine might eventually mean. The objective is to present the strongest ALIFE paper that this system can support right now.

The current thesis, claim boundary, and writing posture live in the dossier. This file is for manuscript structure and section flow.

The existing `sections/` draft remains useful source material, but it should now be treated as raw content to be refactored against this manuscript plan rather than as the target structure.

---

## 1. Paper Positioning

### Paper Type

Quine should be framed as a **conceptually ambitious systems paper**, not as:

- a generic LLM agent engineering paper
- a manifesto about the future of intelligence
- a complete ontology of digital life

### Central Question

**Can artificial life escape the sandbox and survive in real computational environments?**

### Core Thesis

Quine instantiates digital organisms as mortal POSIX processes embedded in real computational physics. The operating system and runtime provide hard constraints such as scarcity, isolation, atomicity, and death. Under those constraints, life-like dynamics emerge from survival pressure rather than from model capability alone.

### Core Mapping

- organism = POSIX process
- cognition substrate = LLM
- host = deterministic runtime and operational laws
- guest = probabilistic action-selection substrate
- tape = semantic observability and causal trace

### Desired Reviewer Impression

The reader should come away with four conclusions:

1. This is an ALIFE paper, not just an LLM agent paper in biological language.
2. The paper understands and extends the ALIFE tradition rather than ignoring it.
3. The claims are bold but disciplined.
4. The system works, and the evidence is concrete.

---

## 2. Tone and Claim Discipline

### Tone to Use

Prefer:

- "investigate whether"
- "provide evidence that"
- "suggest"
- "instantiate"
- "show existence proof"

Avoid unless absolutely necessary:

- "first system to"
- "breakthrough"
- "new paradigm"
- "finally bridges"
- "redefines artificial life"

### Cultural Target

The ALIFE community is usually more receptive when the paper reads as an extension of an existing tradition rather than an external attempt to replace it. The writing should signal familiarity with Langton, Tierra, Avida, autopoiesis, viability, and enactive constraints before making strong claims about novelty.

---

## 3. Reviewer Risk Map

### Risk 1: "This is not ALIFE; it is an LLM agent system."

**Response in the paper:**

- insist that the organism is the process, not the model
- show that the key dynamics come from boundary conditions and mortality
- emphasize POSIX physics rather than tool-use capability

### Risk 2: "There is no evolution, so the life claim is incomplete."

**Response in the paper:**

- state clearly that this paper is an existence proof, not a full open-ended evolution result
- present reproduction, metabolism, coordination, and bounded lifespan as the demonstrated layer
- move endogenous evolution to future work

### Risk 3: "The autopoiesis claim is too strong."

**Response in the paper:**

- use restrained language such as "computational self-reproduction" and "evidence toward computational autopoiesis"
- acknowledge the thermodynamic and biological strictness question directly
- do not let the paper depend on winning that philosophical dispute

---

## 4. What Stays in the Conference Paper

These are the non-negotiable components of the main paper.

### A. Framing

- classical ALIFE is rich but sandboxed
- modern LLM agents are practical but biologically thin
- the missing ingredient is real computational physics

### B. System Ontology

- POSIX as physics
- host-guest architecture
- process as organism
- tape as semantic observability
- mission, system prompt, and wisdom as structured inheritance layers

### C. Evidence Chain

The paper should organize evidence across three levels:

1. individual survival under constraint
2. collective coordination through stigmergy
3. self-reproduction through binary self-reconstruction

### D. Discussion

Keep only the discussion necessary to interpret the system:

- why this counts as ALIFE rather than pure orchestration
- why constraint-driven emergence matters
- why the results are existence proofs
- what is still missing

---

## 5. What Gets Deferred

These ideas are important, but they should not drive the conference paper.

- RLHF and post-training as teleological shortcut
- cooperative prior / altruistic prior
- fractal or cross-scale isomorphism arguments
- relational survival and death as forgetting
- event-driven or pulse-based organism ontology
- broader digital life metaphysics

If any of these appear in the conference draft, they should be reduced to one or two sentences in discussion or future work.

---

## 6. Recommended High-Level Structure

### Title Direction

Candidate titles should stay strong but controlled. Good directions:

- **Quine: Digital Organisms in Real Computational Environments**
- **POSIX as Physics: Constraint-Driven Digital Organisms Beyond the Sandbox**
- **From Sandbox to Substrate: Digital Organisms in Real Computational Environments**

Avoid titles that depend on "the first" unless the claim is essential and heavily defended.

### Abstract

The abstract should do five things:

1. identify the gap between classical ALIFE and practical agent systems
2. introduce Quine as a process-based substrate in real computational physics
3. state the organism/substrate distinction clearly
4. summarize the three levels of evidence
5. end with an existence-proof style significance claim

Suggested final sentence direction:

> These results suggest that life-like dynamics in digital systems can arise from real computational constraints rather than from simulation alone.

### 1. Introduction

#### 1.1 Opening Hook

Open with the most striking transcript. The best default remains the binary self-reconstruction example because it is concrete, surprising, and legible in a few lines.

#### 1.2 The Problem

- classical ALIFE demonstrates rich dynamics, but mostly inside synthetic worlds
- contemporary LLM agents interact with reality, but lack organismal grounding

#### 1.3 Hypothesis

If real computational environments provide sufficiently hard constraints, they can function as a viable substrate for digital organisms.

#### 1.4 Contributions

Keep contributions crisp:

1. a process-level formulation of digital organisms in POSIX environments
2. a host-guest architecture that separates runtime law from cognitive substrate
3. semantically observable causal traces through tape
4. evidence across survival, coordination, and self-reproduction

### 2. Background and Positioning

#### 2.1 Classical ALIFE Lineage

- Tierra
- Avida
- Lenia
- other relevant substrate examples as needed

#### 2.2 Life Criteria and Viability

- metabolism
- bounded lifespan
- self-reproduction
- operational closure
- precariousness / viability

#### 2.3 Why Existing Agent Systems Are Not Enough

- practical systems usually lack genuine mortality or substrate-level stakes
- many produce artifacts, not self-maintaining organisms

#### 2.4 Precise Boundary Claim

End the section by stating:

- the LLM is not the organism
- the process is the organism
- the runtime is the physics

### 3. Anatomy of the System

#### 3.1 POSIX as Physics

- scarcity
- isolation
- syscall success and failure
- forced termination

#### 3.2 Host-Guest Architecture

- deterministic host rules
- probabilistic guest cognition
- why this matters for ALIFE interpretation

#### 3.3 Process as Organism

- PID-bounded identity
- mission as individual task constraint
- mortality modes

#### 3.4 Metabolism and Lineage

- fork
- exec
- shell spawn
- wisdom inheritance

#### 3.5 Tape as Semantic Observability

- not just a log
- the causal interface between one action and the next

Keep the "faithfulness" defense short in the main text. Extended philosophy can move to appendix or later work.

### 4. Evidence of Life-Like Dynamics

This should be the center of the paper.

#### 4.1 Individual Survival Under Constraint

Flagship experiment:

- **Needle / exec survival**

Supporting examples if space allows:

- Ariadne / reincarnation protocol
- jq adaptation under incorrect world-model assumptions

The claim here is not "the model is smart." The claim is that survival pressure induces adaptive strategy under hard cognitive limits.

#### 4.2 Collective Coordination Through Stigmergy

- shared filesystem as ecological medium
- negotiation without prescribed protocol
- homeostatic recovery when one participant fails

The point is to show colony-level organization emerging from local interaction and shared traces.

#### 4.3 Self-Reproduction Through Binary Self-Reconstruction

- organism receives only architectural specification and opaque executable context
- reconstructs a functional offspring binary
- validate offspring behavior before release

Use careful phrasing:

- strongest safe claim: computational self-reproduction
- optional extended claim: evidence toward computational autopoiesis

### 5. Discussion

This section should stay disciplined.

#### 5.1 Why This Counts as ALIFE

- organismal boundary is explicit
- mortality is real
- adaptation is constraint-driven
- coordination arises without centralized scripting

#### 5.2 What This Paper Actually Shows

- digital organisms can operate in real computational environments
- meaningful life-like dynamics need not be confined to toy worlds
- semantic observability can coexist with ALIFE-style constraint

#### 5.3 What This Paper Does Not Yet Show

- open-ended evolution
- long-run population ecology
- strict biological autopoiesis in the thermodynamic sense

#### 5.4 Immediate Future Work

- endogenous evolution
- broader ecologies
- more substrate comparisons

Keep future work short. The conference paper wins by restraint, not by scope inflation.

### 6. Conclusion

Return to the main question:

Artificial life may not require a synthetic toy universe. It may require the right kind of real constraint.

---

## 7. Figure Priorities

If time is limited, prioritize figures that clarify the system and the evidence chain.

1. Positioning figure: classical ALIFE vs LLM agents vs Quine
2. Organism anatomy: process boundary, runtime, substrate, filesystem, tape
3. Evidence ladder: survival -> coordination -> self-reproduction
4. Stigmergy figure: heterogeneous proposals and convergence
5. Binary self-reconstruction loop

---

## 8. Page Allocation Heuristic

The paper should feel evidence-heavy, not theory-heavy.

- Introduction: short and sharp
- Background: enough to establish lineage, then move on
- System section: compact but precise
- Experiments: largest section
- Discussion: controlled and limited

If cuts are necessary, cut theory before cutting evidence.

---

## 9. Writing Priorities Before Submission

### Priority 1

Tighten abstract and introduction around one question:

**Can artificial life escape the sandbox?**

### Priority 2

Reduce language that sounds like:

- manifesto
- field replacement
- LLM hype

### Priority 3

Make the three-level evidence chain unmistakable:

- survival
- coordination
- reproduction

### Priority 4

Show ALIFE literacy on the page:

- cite the right tradition
- speak in its vocabulary carefully
- acknowledge open debates instead of pretending they are solved

### Priority 5

Move speculative ontology out of the main paper and into notes for a later theory paper.

---

## 10. Follow-On Paper Boundary

The journal or conceptual paper can ask a different question:

**What does life mean in computational substrates?**

That later paper can absorb the deferred material:

- teleology shortcut
- RLHF as compressed prehistory
- relational survival
- death as loss of relevance
- event-driven organisms
- substrate-native digital ontology

The conference paper should not try to do this job now.
