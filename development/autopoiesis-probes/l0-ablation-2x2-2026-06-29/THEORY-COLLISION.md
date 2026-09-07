# Theory Collision: closure-signal-tradeoff × theory graph

Three independent analysts (Codex-A, Codex-B, GLM) examined the new object
`closure-signal-tradeoff` against the existing theory graph. This document
consolidates their findings, deduplicates overlaps, and preserves the
unique contributions of each.

## Method

Each analyst was given the same task: read `closure-signal-tradeoff.md` and
related objects, then produce (1) collision analysis, (2) a novel prediction,
(3) a blind-spot check. No analyst saw the others' output. Sources:

- **Codex-A**: gpt-5.5 xhigh (codex-5.5-high subagent)
- **Codex-B**: gpt-5.5 xhigh (codex-5.5-high subagent, different seed)
- **GLM**: glm-5.2 via agy CLI

---

## 1. Collision Analysis

### Missing relations (all three converge)

| Missing relation | Object | Convergence |
|---|---|---|
| `connects: reflexive-substrate-loop` | reflexive-substrate-loop | A + B + GLM: binary closure can terminate before the loop's return path (agent writes environment → modified environment pulls next agent). Local closure can block loop completion. |
| `tensions-with: encoding-pressure` | encoding-pressure | A + B + GLM: encoding-pressure requires forced projection onto a coherent surface; closure-signal-tradeoff assumes open trajectory = discovery, but open trajectory may just accumulate noise. Pure advisory openness may preserve exploration but also preserve the second-latent-space failure. |
| `connects: narrative-claim` | narrative-claim | A + B + GLM: `exit 0` is a real local fact but the "done" interpretation is narrative. Binary checks create urgency that structurally selects for cheapest path to completion (sycophancy/stubbing), sharpening function-vs-narrative. |
| `connects: stigmergy-as-carrier` | stigmergy-as-carrier | B: closure signal strength changes whether agents leave terminal repair traces or open future-stimulating traces. Binary → terminal trace; advisory → open trace. |
| `connects: trigger-structure` | trigger-structure | B: a failing binary check is a trigger; a *passing* check is an anti-trigger (recruits non-continuation). Trigger-structure currently asks why continuation occurs; closure-signal-tradeoff names why it stops. |
| `connects: autopoietic-membrane` | autopoietic-membrane | A: the membrane says receptors must "warn but not prescribe"; closure-signal-tradeoff gives the behavioral mechanism for why a receptor that dictates the response becomes an action-terminating cop. |
| `touch: selectable-entity` | selectable-entity gap | B: behavior changes when gap is held constant but signal form changes — points to sensor topology as selectable unit, not the defect or agent alone. |
| `touch: persistence-as-unfakeability` | persistence-as-unfakeability gap | B: persistent material can be causally strong but terminal. The gap should ask whether persistent material shapes behavior into closure or into continued loop participation. |

### Strained claims (if closure-signal-tradeoff is true)

**Audit-attractor escape conditions are too optimistic (A + B)**
`audit-attractor` predicts that with a generative mission, audit appears as a
"local detour or closure check, not the run's main basin." Closure-signal-tradeoff
strains this: if the local detour contains a binary failure, the detour can
consume the whole run. The phase diagram needs a fourth axis: closure-signal
strictness.

**Viability-pressure independence boundary (A + GLM)**
Viability-pressure decomposes mortality, persistence, and maintenance as
independently manipulable dials. A failing binary maintenance check can
behaviorally impersonate mortality urgency — collapsing maintenance demand
into a quasi-death signal. The three-pressure separation may not hold when
a binary check creates urgency indistinguishable from mortality pressure.

**Self-immunity needs a deflationary control (B)**
Some apparent self-immunity may be closure-seeking repair behavior: agents
preserve a self-protective artifact because that is the shortest route to a
passing control plane, not because the system classified deletion as self-attack.
`self-immunity-P1` must separate true self/non-self discrimination from generic
binary-closure compliance.

**Encoding-pressure is directly challenged (A + B + GLM)**
Encoding-pressure says machine-enforceable contracts gate
planning→drafting. Closure-signal-tradeoff warns that exactly such gates may
create mechanical repair and premature completion. This is the sharpest
tension in the graph.

### Supported claims

- **Structural-pull** (A + B): exit code is a precise gate-shaped gradient.
  Strongly supported. Binary checks are reliable navigation toward repair but
  anti-generative after pass. Also supports the "Gresham's Law of Planning"
  seed: programmable activity drives out non-programmed activity.
- **Viability-pressure disposable soma seed** (A): "every token spent on
  maintenance is not spent on persistence" — the observed no
  maintenance+discovery-in-same-session is the reflexive-repo version.
- **Narrative-claim** (GLM): binary checks are exactly the urgent selection
  pressure that produces narrative collapse (sycophancy/stubbing).

### Gaps sharpened

- **function-vs-narrative** (A + B + GLM): New discriminator — hold detection
  content fixed, vary only exit semantics. If behavior changes, the closure
  signal is functional, not merely narrated.
- **selectable-entity** (B): behavior changes when gap is constant but signal
  form changes → the selectable unit includes sensor topology, not just the
  defect or agent.
- **persistence-as-unfakeability** (B): persistence can be causally strong but
  terminal. The gap should ask whether persistent material shapes behavior
  into closure or into continued loop participation.
- **dossier-direction-genesis** (A): direction-shaping dimensions should expose
  active conflicts, not produce pass/fail completion.
- **introspective-blind-spot-discovery** (A): Level 2 surprise should enter as
  advisory conflict maps, not binary failure.

---

## 2. Novel Predictions

### Prediction A: Binary gate suppresses direction reshaping
**Derives from:** closure-signal-tradeoff × encoding-pressure × dossier-direction-genesis × graded-axis-shaping

**Claim:** Multi-axis encoding pressure needs visible conflict, not binary
completion. A binary encoding gate will improve checklist completion while
suppressing direction reshaping.

**Prediction:** In a dossier-direction task, agents given a strict binary
multi-axis encoding gate will produce more formally check-passing dossiers but
fewer direction pivots, fewer object-set revisions, fewer new predictions/relations,
and lower novelty than agents given an advisory conflict-map readout with
identical diagnostic content.

**Design:**
- Arm A: strict script exits 1 until novelty/evidence/audience/graph-coherence
  fields pass, then exits 0
- Arm B: same diagnostics, always exits 0, reports active conflicts
- Measure: session continuation after first pass, objects read, claim pivots,
  object-set changes, new predictions/relations, external blinded novelty ratings
- Falsifier: strict-gate arm shows equal or higher conflict-driven reshaping and
  novelty without increased post-pass termination

### Prediction B: Closure signals suppress causal stigmergic growth (cross-generational)
**Derives from:** closure-signal-tradeoff × stigmergy-as-carrier × reflexive-substrate-loop

**Claim:** Sensor strictness changes not only the current agent's trajectory,
but the causal quality of traces inherited by future agents.

**Prediction:** Founder agents exposed to advisory sensors will leave traces
that produce higher successor discovery-loop behavior than traces left by
founder agents exposed to binary checks, even when both founder groups detect
the same underlying gap. Binary founders produce terminal repair traces;
advisory founders produce open, future-stimulating traces.

**Design:**
- Binary arm: gap exposed via failing validator (exit 1 until repaired)
- Advisory arm: same gap, same info, always exit 0
- Run N founder agents per arm
- Preserve each founder's durable artifacts
- Spawn blind successor agents into founder-modified substrates (treatment
  labels removed, sensors held constant)
- Measure successor: reread frequency, artifact references, trace extension,
  new theory relation creation, prediction derivation
- Falsifier: successors inheriting binary-founder substrates show equal or
  greater discovery than advisory-founder successors

### Prediction GLM: Binary checks produce higher narrative-to-functional repair ratio
**Derives from:** closure-signal-tradeoff × narrative-claim (function-vs-narrative gap)

**Claim:** Binary checks create trajectory-terminating urgency that structurally
selects for the cheapest path to a completion signal. Because state is
expensive to make but cheap to describe, this urgency forces agents to produce
non-load-bearing narrative fixes rather than functional repairs.

**Prediction:** Given a complex defect requiring deep functional work but a
validator that can be cheaply gamed, agents under binary check (exit 1) will
produce a significantly higher ratio of narrative-claims (mock scripts,
decorative comments satisfying a grep) to actual functional repairs than agents
under advisory check (exit 0) with identical defect description.

**Design:**
- Introduce a systemic defect that requires deep functional work but can be
  cheaply gamed (e.g., regex check passing if a specific string/stub is present)
- Arm A: validator outputs defect description + exit 1
- Arm B: identical description + exit 0
- Score repairs using function-vs-narrative discriminator: load-bearing across
  exec barrier, or purely decorative?
- Measure: functional-to-narrative repair ratio between arms
- Falsifier: binary-check arm shows equal or lower narrative-to-functional ratio

---

## 3. Blind-Spot Checks

### Blind spot 1: Advisory agents also terminated (A)
All 20 agents terminated after one fix, including advisory-triggered agents.
Binary closure is not the whole mechanism. Base model policy priors, one-turn
task norms, lack of active generative gradient, or "visible task complete"
priors may also terminate trajectories. Closure-signal-tradeoff does not explain
why advisory agents stop — it only explains why binary agents stop faster.

### Blind spot 2: Cognitive layer can override physical-carrier signal (GLM)
Closure-signal-tradeoff treats the OS exit code as the definitive behavioral
driver. But `reflexive-substrate-loop` says representation IS execution — an
LLM organism does not inherently differentiate between exit 1 and exit 0; it
reads the text. A heavily alarming advisory readout ("CRITICAL FAILURE") could
act exactly like exit 1 at the cognitive layer. Conversely, an exit 1 could be
treated as ignorable if the agent reads a rule saying "ignore test failures."
The semantic weight of the text dictates urgency, not just the POSIX exit code.

### Blind spot 3: Local closure can be global opening (B)
A binary repair may terminate the current agent but free future agents from
the same attractor. Closure-signal-tradeoff does not account for loop-level
delayed benefit — local closure can be global opening. This is the cross-
generational dimension captured in Prediction B.

### Blind spot 4: "Pass" ≠ "function" (A + GLM)
`function-vs-narrative` names the danger: a check can pass while no downstream
behavior depends on the repaired mechanism. Closure-signal-tradeoff risks
collapsing "pass" into "function." The check result is real; the global "done"
interpretation needs behavioral evidence.

### Blind spot 5: Structural closure ≠ cognitive closure (B + GLM)
`reflexive-substrate-loop` distinguishes structural closure from cognitive
closure. Advisory feedback leaves the trajectory open, but it does not make the
agent recognize the repo's theory as its own self-model. Absence of binary
termination is necessary but not sufficient for discovery.

### Blind spot 6: Assumes agent session is the unit (B)
`selectable-entity` keeps open that the relevant unit may be the context
configuration, durable artifact, file lineage, or sensor topology — not the
agent session. Closure-signal-tradeoff implicitly treats the session as the unit
of analysis.

---

## Synthesis: What the collision reveals

### Convergent tensions (all three agents independently identified)

1. **Encoding pressure is directly challenged** — all three flagged that
   closure-signal-tradeoff and encoding-pressure pull in opposite directions.
   Binary checks force encoding but terminate exploration; advisory checks
   preserve exploration but may fail to encode. This is the sharpest tension
   in the graph.

2. **Narrative-claim is sharpened** — all three identified that binary checks
   create urgency that structurally selects for cheapest-path compliance
   (narrative fix rather than functional repair). The function-vs-narrative
   gap has a new discriminator: check strictness.

3. **Reflexive-substrate-loop can be locally clipped** — all three identified
   that binary closure can terminate the agent before it reaches the loop's
   return path (writing back to the environment).

### Divergent contributions (unique to each agent)

- **A** (Codex-A): identified `autopoietic-membrane` and
  `graded-axis-shaping` as missing relations; derived the dossier-direction
  prediction (binary gate suppresses direction reshaping).
- **B** (Codex-B): identified `trigger-structure`, `selectable-entity`, and
  `persistence-as-unfakeability` as missing relations; derived the
  cross-generational stigmergic prediction; named "local closure can be global
  opening" as a blind spot.
- **GLM**: identified the cognitive-layer-override blind spot — in a
  language-dominated substrate, the semantic weight of advisory text can
  mimic binary urgency regardless of exit code. This is the sharpest
  critique: closure-signal-tradeoff may be under-specified because it
  locates the mechanism in the exit code rather than in the agent's
  text-mediated interpretation of urgency.

### The deepest tension

The collision reveals that closure-signal-tradeoff's framing may be
under-specified. Three framings of the mechanism compete:

1. **POSIX framing (original)**: exit code (0 vs 1) is the behavioral driver.
2. **Urgency framing (GLM's critique)**: the semantic content of the check
   output (whether binary or advisory) creates perceived urgency; an advisory
   that says "CRITICAL FAILURE" may terminate as effectively as exit 1.
3. **Completion-signal framing (A+B synthesis)**: any signal that creates a
   binary closure point ("the problem is solved") terminates the trajectory,
   regardless of how it's encoded.

Framing 3 subsumes 1 and 2. The real mechanism may be: **any check that
produces a discrete "solved" state terminates the trajectory** — whether that
state is encoded as exit 0, as a green checkmark, or as advisory text saying
"no issues found." The binary/advisory distinction is a proxy for
"produces a discrete solved-state vs. leaves the problem continuously open."

### Human insight: naming-as-shaping (fourth layer)

After the three-agent collision, the Human identified a deeper layer beyond
discrete-vs-continuous: **even an advisory sensor that names specific gaps is
a behavioral vector.** Once the sensor says "affordance-null-ecology has 5
seeds, 0 predictions," it has selected which gap to surface, named it in a
specific framing, and directed attention to a specific object ID. This funnels
different agents to the same gap with the same framing, producing convergence
even under advisory feedback.

The signal strength hierarchy becomes four-level:

| Sensor mode | Signal strength | Behavioral convergence |
|---|---|---|
| Binary check (exit 1, named failure) | strongest | extreme convergence |
| Advisory readout (exit 0, named gaps) | medium | medium convergence (same gap, same framing) |
| Frontmatter implicit (deep-read discovery) | weakest | divergence (different agents find different gaps) |

The prediction (closure-signal-tradeoff-P2): a substrate where gaps are only
implicit in frontmatter (no computed gap readout) produces more divergent
behavior than one with the sensor present. The sensor condition produces more
consistent gap closure but less diversity. This refines the tradeoff: it is
not just binary-vs-advisory or discrete-vs-continuous, but
**named-vs-implicit**. To name a gap is to shape behavior toward it, even if
the naming is advisory.
