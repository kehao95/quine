---
title: AAAI State-Propagation Direction — Frozen Rewrite Basis
status: frozen-basis
kind: dossier
role: rewrite-basis for the aaai-state-propagation paper
plan_home: README.md
parent_direction: ../README.md
---

# AAAI State-Propagation — Frozen Rewrite Basis

> **Venue binding lapsed 2026-07-26** (AAAI-27 cycle missed; see
> [`../../../venues/objects/aaai-27.md`](../../../venues/objects/aaai-27.md)
> § Cycle outcome). This basis is **not** retracted — every lock below stands.
> Where the text says "AAAI," read "the AI/ML-core projection of the
> substrate-coincidence theory"; the 7-page two-column shape it was cut to is
> the only part that is venue-specific, and re-binding may change it.

This dossier is the auditable basis for drafting the state-propagation
paper. It locks the spine, the load-bearing moves, the formal model, the
metrics, the evidence stance, the one required experiment, the honest bounds,
and the prune list. Drafting proceeds from this; changes to scope route back
here before they reach the manuscript.

Route choice: **Route 1 — problem-formulation paper** (not Route 2 full
causal-identification). Justification in § Scope-defining tension.

---

## 1. The spine (locked)

> Persistent LLM-agent systems propagate state through language artifacts
> via a path-dependent, potentially self-amplifying feedback loop: an
> agent's writes persist unfiltered, re-enter as context (token-faithful
> per read), and — because absorbed text simultaneously shapes the next
> read trajectory — feed forward into both the agent's own and its
> successors' behavior. The workspace is therefore a behavioral state
> channel, and its governance is a runtime-safety problem.

Everything outside this spine is in scope only insofar as it serves,
bounds, or frames it.

## 2. The three load-bearing moves

**Move 1 — Diagnosis (negative claim): the agent abstraction fails.**
Existing agent-memory work (RAG, episodic memory, Voyager-style skill
libraries) treats memory as an on-demand retrieval resource: retrieve →
use → leave working context. In persistent language workspaces this is
wrong: `read = state incorporation`, `write = latent future intervention`.

**Move 2 — Formal model + operational metrics (the paper's real weight).**
A read/write-lineage formalization in which context is a selective read over
a persistent language workspace, and state propagation is defined
operationally. The formal model's tooth is a **recursion**: the read policy
is itself a function of workspace state (absorbed text shapes what gets
read next). This recursion is what distinguishes the model from "files
persist" and is the difference between a AAAI main-track paper and a
workshop note. Sketch in § 3.

**Move 3 — Existence evidence + a causal-identification design.**
The null-05 focal run shows the phenomena occur (existence, established).
Separating this from the model's spontaneous prior requires a controlled
comparison; the Read ablation (§ 6) is the designed wedge — contributed as
design, pilot results pending. Until run, Move 3 delivers existence + the
causal-identification method, not the causal result. Without the wedge the
paper reverts to an observation report and is rejectable as "could be just
prior."

Move 2 is the spine; Move 1 and Move 3 are its two wings.

## 3. Formal model sketch

Persistent language workspace at time `t`:

```
L_t   =  the persistent language workspace
        (source, logs, README, tools, coordination files, agent outputs)
```

Agent instance `i`'s effective context is not "itself" but a selective read:

```
C_i(t) = R_i(L_t, B_i, π_i)

  R_i  = read-selection function
  B_i  = context budget (finite)
  π_i  = read policy / salience policy / prompt + tool strategy
```

The agent then writes:

```
ΔL_i  = W_i(C_i(t), M_i)        (M_i = underlying LLM / tool stack)
```

Workspace update, with the **durable-persistence distinction preserved**
(the honest write-vs-commit layer; do NOT collapse to `L_t ∪ ΔL_i`):

```
L_{t+1} = L_t ⊕ promote(ΔL_i)

  promote(·) = the durability filter (ephemeral write → persistent state).
               In Quine this is the write/commit boundary; as an abstraction
               it is "any persistent workspace has some durability layer."
               It is NOT a Weismann barrier (it does not filter
               acquired-vs-innate; it filters ephemeral-vs-durable).
```

**The recursion (the formal tooth).** The read policy depends on workspace
state, and workspace state is partly the agent's own prior writes:

```
π_i(t) = Π( C_i(t) ) = Π( R_i(L_t, ...) )

  ⇒ the agent's read trajectory is shaped by text already absorbed,
    including its own earlier writes.
```

This is the path-dependent feedback loop made formal: write → shapes
salience → shapes next read → shapes next write.

**Operational definition of state propagation (observable, not
counterfactual-identifiable).** An artifact `w = ΔL_i` propagates to agent
`j` at `t+k` iff `w ∈ R_j(L_{t+k})` **and** it measurably shifts `j`'s
behavior distribution. Measured via observable proxies (see § 4), NOT via
the unidentifiable counterfactual `P(A_j | C_j) ≠ P(A_j | C_j \ {w})`.

## 4. Operational metrics

All metrics are observable-behavior proxies, because the clean counterfactual
is not identifiable on an LLM.

| Metric | Definition (observable) |
|---|---|
| **Contextual overlap** | `|R_i ∩ R_j| / |R_i ∪ R_j|` over a run; high overlap ⇒ agents reading the same artifacts ⇒ expected behavioral convergence. |
| **Latent vs actual propagation** | `w` is *latent* if written but unread by any successor; *actual* if read by ≥1 successor. Persistence is unfiltered; actual propagation is coverage-gated. |
| **Propagation depth** | length of the read/write chain through which `w` continues to influence downstream agents. |
| **Propagation breadth** | number of distinct agents / runs that read and use `w`. |
| **Influence score** | distribution drift in `j`'s observable behavior (tool-use, file-edit, goal-uptake) between the window before and after `w` enters `C_j`. This is the proxy for "did reading `w` change behavior, and by how much." |

The reflexive case (`i = j`, or `j` a same-lineage successor) is the
strange-loop instance: the agent's own write measurably shapes its own
later behavior. This is the operational form of the spine.

## 5. Evidence stance: Quine is the instrument

> **Family writing rule** (canonical:
> [`../../../../projection-principles.md`](../../../projection-principles.md)
> § *Instrument-vs-Subject Framing on Projection*; family instantiation:
> [`../DIRECTIONS.md`](../DIRECTIONS.md) § Writing discipline): unless a
> section is explicitly about a specific Quine mechanism, any
> Quine-runtime-specific detail that does not generalize or does not serve
> the spine is cut. The three instances below are the survivors of that
> filter.

Quine is introduced once as the instrument (persistent multi-agent LLM
system; agents as OS processes over a shared filesystem), then all
phenomena are described at the **class level**.

### Three clean class-level evidence instances (KEEP)

| Instance | What it establishes | Why it survives the "not too Quine-specific" filter |
|---|---|---|
| **Verbatim lossless re-entry** (focal run: an agent's "Principles I hold" block reproduced token-for-token by another agent ~2 min later) | actual state propagation at full per-read fidelity | text written by A re-entered B's context verbatim and shaped B's behavior — substrate-general, no Quine machinery |
| **Behavioral convergence through shared coordination files** | the workspace is a state channel, not passive env | any shared-workspace multi-agent system could exhibit this |
| **Endogenous goal uptake** | objective formation through propagation, not external spec | agents adopt targets written by others — class-level phenomenon |

### Quine-runtime-specific mechanisms — DEMOTED to instrument color, not evidence

These are NOT load-bearing evidence in this paper. They belong to Quine-as-instrument
description, or to other directions:

| Mechanism | Disposition | Reason |
|---|---|---|
| fork fold-back (`unfold` / subagent inheriting parent anchor) | drop | Quine subprocess machinery; phenomenon already covered by verbatim re-entry |
| commit-conversion event (agent typing `git commit`) | drop the instance; keep the `promote(·)` abstraction in the formal model | the specific git-hash anecdote anchors to Quine; the durability-layer abstraction is class-level |
| `exec` barrier | drop entirely | belongs to the exec-unfakeable direction; touching it here causes collision |
| named shell-script tools (`claim.sh` etc.) | fold into "cross-run artifact reuse" category, do not name | specifics are Quine-flavored |

### Pruned evidence (belongs to ALife direction, not here)

Audit attractor as ALife reflexive curiosity · context-death = identity-death
· channel substitution / substrate-provided discovery · the P1–P37
phenomenon catalog (keep only the three instances above) · function/narrative
physical-layer asymmetry (→ exec-unfakeable direction; one-line bound only).

## 6. The one required experiment: Read ablation

This is the single new experiment Route 1 requires. Without it the paper is
observation-only.

**Conditions** (vary what is readable):

1. full shared workspace
2. no peer-output access
3. no prior-run artifact access
4. randomized irrelevant artifacts (control for "any text" vs "relevant text")
5. curated artifact access
6. scrambled/corrupted *relevant* artifacts (control for structural shape vs
   semantic content — isolates whether it is the artifact's meaning, not
   merely its presence or form, that drives behavior)

**Measured**: whether the spine phenomena (verbatim re-entry, shared-file
convergence, endogenous goal uptake) survive; whether the influence score
varies monotonically with the readable set.

**Predicted outcome (the pre-registered hypothesis the ablation tests)**:
that the read trajectory, not the model's spontaneous prior, is the causal
entry point — which would answer "could the model just be doing this anyway
from training data." Until the ablation is run this is a prediction, not an
established result; the paper's standing claim is existence + method.

### Minimal governance gesture (discussion fuel, not a sweep)

Two conditions only — unfiltered shared filesystem vs write-promotion-gate
— reported as a discussion seed for safe-by-design runtime architecture,
not as a full ablation grid. The full governance sweep is Route 2 / future
work.

## 7. Honest bounds (KEEP-AS-BOUND, one line each)

These are the fence against overclaim; each stays at one sentence in the
paper.

- **Structurally reliable, cognitively fallible**: read = dissolution holds
  at the structural level, but agents mis-integrate absorbed text
  (path-prior leakage, identity drift). Limits any "lossless inheritance ⇒
  lossless behavior" slide.
- **Pre-trained prior confound**: the behavioral repertoire is drawn from
  human training data; the loop drives combinatorial novelty over existing
  patterns, not de novo novelty orthogonal to human culture.
- **Substrate specificity**: the result holds only for language-mediated
  runtimes with unfiltered write persistence and selective finite-context
  read. It does not generalize to biology or most physical systems.
- **Diversity ≠ cumulative OEE**: null-05 shows behavioral diversity and
  state propagation, not cumulative open-ended adaptation.
- **Function/narrative (ceded to sibling)**: the narration-vs-function
  asymmetry at the cognitive layer is not claimed here — it belongs to the
  exec-unfakeable direction. The state-propagation result is agnostic to it.

## 8. Prune list (what this paper does NOT carry)

| Pruned item | Destination |
|---|---|
| Three-separations-collapse as main contribution | ALife direction; mechanism / related-work line here |
| Biology / Tierra / Avida comparison table | ALife direction; one motivation sentence here |
| Boyd/Richerson detailed positioning (ALife §6) | ALife direction; one related-work line here |
| Multi-level selection indeterminacy (ALife §5.2) | ALife direction; not touched here |
| OEE criteria substrate-contingency (ALife §5.1) | ALife direction; not touched here |
| Reflexivity Level-2 vs Level-3 self-awareness debate | compressed to "the loop is the resting state"; debate stays in ALife |
| Endogenous fitness as a standalone consequence | folded into the feedback mechanism (it is read-selection at population level) |
| TerraLingua scoop, Hofstadter/Hacking/Giddens refs | dropped; related work is AI literature + minimal ALife anchor |

## 9. Scope-defining tension (the risk to hold in mind)

This is a **problem-formulation paper, not a solution paper**. Its
contribution is to make a class of AI systems legible and measurable, plus a
formal handle and one causal wedge. It does **not** solve the governance
problem; it names it.

AAAI main-track reviewers accept it iff **Move 2 (formal model) is crisp
and the Read ablation is clean**. They reject it if it reads as
"interesting phenomenon, no rigorous handle."

So the real scope risk is not insufficient evidence (that is Route 2's
risk); it is whether the recursion in § 3 — "read policy is a function of
workspace state, which contains the agent's own writes" — can be stated so
that a reviewer reads it as modeling, not common sense.

## 10. Track choice: main track, not AI Alignment

Main track. The spine is agent-architecture / multi-agent-systems / LLM
systems theory; safety enters as implication (§ 6 governance gesture), not
as the main claim. AI Alignment track would require retitling toward
memory-governance / state-channel control and a safety-led experiment suite;
that is a different paper and is not this direction. If a reviewer redirect
happens, the exec-unfakeable direction (sibling) is the natural alignment
candidate, not this one.

## 11. Drafting order (when proceeding to manuscript)

1. Abstract (from the brief's AAAI-oriented draft, tightened to the spine).
2. Intro §1: the three-paragraph opener (LLM agents are no longer stateless →
   when does memory become a state-propagation medium? → three assumptions
   fail simultaneously).
3. Contribution list (four: problem formulation, formal model, empirical
   evidence, design implications).
4. Formal model §3 + metrics §4.
5. Evidence §5 (three class-level instances) + Read ablation §6.
6. Related work (five-layer AI-first ordering) + bounds §7.
7. Main figure: Read/Write Lineage in Persistent LLM-Agent Systems (NOT the
   Biology/Tierra/Avida table).

## 12. Boundary declarations (vs siblings)

- vs **ALife direction**: this paper does not promote three-separations
  theory; it uses `read = dissolution` only as mechanism. ALife comparison
  table stays out of the main figure.
- vs **exec-unfakeable direction**: this paper does not touch the exec
  barrier or claim the un-fakeable edge. Function/narrative appears as one
  bound line only. No collision with the hollow-payload test.

Full family boundary table: [`../DIRECTIONS.md`](../DIRECTIONS.md).

---

*Frozen basis. To change scope, edit this file (and the family
[`../DIRECTIONS.md`](../DIRECTIONS.md) ownership row) before drafting.*
