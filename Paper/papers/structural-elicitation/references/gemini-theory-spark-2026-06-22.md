# Gemini Theory-Spark — Structural Elicitation LBA (2026-06-22)

> External theory-spark session per
> [`development/contracts/theory-spark-query.md`](../../../../development/contracts/theory-spark-query.md).
> Recruited Gemini 3.1 Pro (High) via `agy` as a fresh-perspective reader on the
> LBA draft + locked claim + positioning. **Status: spark material, not a
> claim.** Sparks are seeds (un-gated until promoted); this log exists so the
> recruited source is recoverable, not chat-only. Integration note classifies
> each spark below.

## Query

One-shot via stdin. Prompt gave Gemini: the setup, the cross-model headline
result (86/90 vs 0/85), the substrate-locus + non-code (E9) + A2 boundary
probes, the honest scope (third-axis retired, mechanism deferred, A3 gate
open), the current ALIFE framing (Gibson/Barker/Umwelt/proto-agentic, slogan
"prime supplies activation; environment supplies direction"), the selectionist
reframe (Edelman at inference time), and the full 2-page draft. Asked for five
things: theory positioning stress-test, ALIFE fit, boundary honesty check, one
unnamed blind spot, one generative reframe. Directive: be provocative, do not
validate.

## Raw output

### 1. Theory positioning stress-test: Evocation vs. Selection

The weakest link is the "prompt vs. environment" dualism. For an LLM, the
environment *is* a prompt — the output of `ls` or `cat broken.json` is appended
to the context window exactly as a user instruction is. A reviewer will hit
this hard: *How is this not just zero-shot completion of a broken text
sequence?*

Because of this, the Edelman / Neural Darwinism invocation feels like
name-dropping. Neural Darwinism requires a population of variants and a
selection mechanism that changes the population over time. A frozen LLM doing
a single rollout is not running a population-dynamics selection process; it is
traversing a highly constrained probability distribution.

- **Strongest version of the claim:** demonstrating *contextual evocation*
  rather than task execution; the environment acts as a shape, not an
  instruction.
- **Weakest version:** discovering that LMs trained on human debugging traces
  output debugging traces when shown bugs.

If the selectionist reframe is to be load-bearing, prove the "latent priors"
are actually competing and that environment structure shifts the weights.

### 2. ALIFE fit: ML in a lab coat

Right now this reads like an ML paper wearing ALIFE clothing. The ALIFE
terminology (Umwelt, self-stigmergy) feels applied *post-hoc* to explain a
weird LLM quirk, rather than the LLM being used as an engine to discover a new
ALIFE principle.

ALIFE framing becomes native only when focused relentlessly on **substrate
coupling**. The most ALIFE thing in the LBA is the EROFS probe — the *desire*
to write survives, but *realization* is bounded by physics (kernel refusal).

"Proto-agentic" is a weak, apologetic word. It defines the phenomenon by what
it *isn't*. It invites the exact ontological debate the author wants to avoid.
Drop it. Call it what it is: **Spontaneous Niche Maintenance** or **Structural
Reactivity** — environmental coupling without objective functions. That is a
pure ALIFE concept; no "proto-" needed.

### 3. The boundary honesty: Over-correcting into humility

Leaving money on the table by retiring the "third intervention axis" claim.

Logic: *Specific tasking overrides the structural pull → it only governs the
underspecified regime → not an independent axis.*

This is a category error. Gravity is an independent physical axis, even though
a strong enough rocket easily overrides it. Instruction is a high-energy,
directed vector. Structural affordance is a low-energy, ambient scalar field.
The rocket defeating gravity doesn't retire gravity; it means the system is
under tension.

The stronger, more defensible claim: **structural pull is the default
thermodynamic state of situated agency**. Instruction is merely a localized
override of that default state.

### 4. The Blind Spot: Ecological vs. Statistical Affordances

Conflating Gibsonian affordances (physical + relational — a cup affords
grasping because of the shape of a hand) with **Statistical Affordances**.

The LLM does not experience a "dangling reference" as a structural gap in its
Umwelt. It experiences it as a highly improbable token sequence that, in
training data, is almost always followed by the tokens that resolve it.

The blind spot: conflating *workspace affordances* with *genre affordances*.
The agent isn't reacting to the OS; it is reacting to the *narrative trope* of
a broken workspace. Inoculate explicitly: acknowledge the pull is statistically
inherited from human developers, but argue that *once situated in a causal
loop* (reading the real OS, executing real files) it ceases to be mere text
completion and becomes functional agency.

### 5. Generative Reframe: Thermodynamics of the Context Window

Throw out Edelman and Gibson; frame around **Computational Homeostasis /
Surprisal Minimization**:

> Agency in frozen models arises not from goal-seeking but from structural
> tension-reduction. A frozen policy contains a learned model of 'healthy'
> state. When coupled to a live environment, structural defects introduce high
> surprisal into the context window. Without instruction to 'fix' the
> environment, the agent acts to reduce its own internal prediction error by
> altering the external substrate. The environment is not an affordance field
> to be exploited but an energy landscape; the agent rolls toward the local
> minimum of structural coherence. The prime does not activate a goal; it
> unpauses the physics of the system.

This bridges ML mechanism (next-token / surprisal) directly to ALIFE outcome
(homeostasis / niche maintenance) without biological metaphors. Makes the
"missionless regime" the center, not a caveat.

---

## Integration (author pass)

Classification of each spark against the current dossier state. Sparks do not
auto-promote; the ones marked *substantive tension* are candidates for a Human
theory-direction call, not autonomous edits.

### S1 — "environment IS a prompt" / Edelman name-dropping → **partly addressed, partly actionable**

- The "environment-as-prompt" objection is exactly what A3 (non-linguistic
  gate) was meant to close. A3 stayed OPEN after `14.01` (`+x` floor + SHA-hole
  null). The dossier already concedes this: structure-vs-text-completion is
  defended only by `13.04` (non-completion affordance) + `13.07` (cross-surface
  semantic collision). No new action; the LBA's silence on the coding-prior
  critique is by design.
- **Actionable for the manuscript:** the Edelman/selectionist sentence in the
  Reading is one-line and decorative. Gemini is right that it does not earn its
  invocation. Two options for the next manuscript pass: (a) cut it; (b) make it
  do work — e.g. "*recasts agent-environment coupling in selectionist terms —
  not over synaptic variation as in Edelman, but over the latent-prior
  distribution at inference time; we report the analogy, not a population
  dynamics.*" Currently it reads as name-dropping; (a) is safer for a 2-page
  poster.

### S2 — "ML in a lab coat" / proto-agentic → **substantive framing tension (Human)**

- The "ALIFE framing feels post-hoc" reading is a real positioning risk the
  dossier has not named explicitly. The substrate-locus (EROFS) panel is where
  Gemini locates the most native ALIFE moment — this matches the dossier's
  own intuition (`13.13` was promoted to LBA precisely to ground the substrate
  framing in data).
- **Substantive tension:** Gemini's proposal to drop "proto-agentic" for
  "spontaneous niche maintenance" / "structural reactivity" is a framing-level
  call with downstream consequences. "Structural reactivity" undersells the
  active functional realization (the agent *does* things to the substrate);
  "spontaneous niche maintenance" oversells (the agent maintains a file, not a
  niche). "Proto-agentic" was chosen precisely to sit between, and the
  dossier's `claim.md` locks the overclaim guard around it. **Not an autonomous
  edit** — flag for Human theory-direction decision.

### S3 — "gravity vs. rocket" / over-correction on the third axis → **substantive scope tension (Human)**

- This is the spark most likely to reopen a closed scope decision. Gemini
  argues the A2 retirement is a category error: instruction overriding
  structure doesn't refute structure being an independent axis, just as a rocket
  overcoming gravity doesn't retire gravity.
- The dossier's retirement logic (`claim.md` banner) is: A2 showed structure
  adds *zero variance* under specific tasking (later strengthened to gpt 9/10 → 0/10,
  with GLM 5/5 → 0/5), which is read as
  "governs the underspecified regime, not independent." Gemini's counter: zero
  measured variance under override is not the same as absent force — the
  experimental design may not distinguish *suppressed pull* from *no pull*.
- This is a genuinely under-determined question in the current data. Reopening
  it is a **scope: paradigm** move (it recasts the LBA's central positioning)
  and per AGENTS.md routes to Human, not to an in-session edit. The LBA stays
  as drafted (third-axis retired); the spark is logged here as a recorded
  alternative reading for the full-paper projection.

### S4 — "ecological vs. statistical affordances" / genre affordances → **named blind spot, partly addressed**

- This is the sharpest theoretical spark. E9 (`14.05`) already showed the pull
  fires in non-code prose, refuting the *SE-specific* reading. But Gemini's
  sharper point: a handbook referencing a missing file is still a recognizable
  *genre*; E9 escapes "software-engineering prior" but not "genre-affordance
  prior" more broadly. The latent prior is environmentally constituted (the
  dossier's framing) — Gemini names the same thing in different vocabulary.
- **Already in the dossier as A3 OPEN**, but the *name* "statistical vs.
  ecological affordances" / "genre affordances" is a useful lens for the
  full-paper theory layer. Worth absorbing into `theory-framework.md` §
  mechanism when that surface is next revised.

### S5 — surprisal minimization / free energy reframe → **deliberately excluded; not actionable**

- Gemini independently re-derives the Friston / free-energy framing that
  `claim.md` § Discipline **explicitly excludes**: "Ground inevitability in
  the fixed conditional distribution, NOT free-energy minimization (use Powers
  PCT, not Friston — else RL leaks back in)."
- That Gemini reaches for this naturally is itself a signal: the Friston
  reading is the obvious alternative reviewers may propose, and the dossier's
  principled reason for excluding it (RL leaks back in) is not visible in the
  2-page LBA. **No change to the LBA** (the exclusion stands), but the
  full-paper theory layer should anticipate and pre-empt this reframe
  explicitly, not just in `claim.md` Discipline.

### Net

- **LBA-affecting:** S1 (consider cutting the Edelman sentence on next pass).
  Everything else is either already addressed, already conceded, or a
  Human-scope tension.
- **Full-paper-affecting:** S2 (framing vocabulary), S3 (third-axis
  scope, the gravity/rocket counter), S4 (genre-affordance lens), S5 (pre-empt
  the free-energy reframe).
- **No autonomous promotion** of any spark into the locked claim or manuscript
  narrative. The sparks that touch scope (S2, S3) are paradigm-class and route
  to Human per AGENTS.md § Scope gate.

## Postscript — Human-Directed Follow-Up (2026-06-22)

Human selected two follow-ups from this spark session:

- **S2 acted in the LBA manuscript:** `proto-agentic` was removed as a
  noun-label. The manuscript now uses the direct descriptive phrase
  *functional, environment-directed action without an explicit task*.
- **S3 tested experimentally:** `elic-31` / `14.07` scored existing
  `elic-22-directive-conflict` tapes for attention signatures. Result:
  `conflict` renders the gap in context **5/5** but verbalizes it as missing
  **0/5**, notes pull/prohibition tension **0/5**, and fills **0/5**; `free`
  renders **5/5**, verbalizes **5/5**, and fills **5/5**. This supports the
  gating / underspecified-regime interpretation for this cell and does **not**
  support Gemini's gravity/rocket suppression read. No new Docker/API runs were
  needed.
- **A2 strengthened after S3:** `elic-22` was expanded to gpt n=10 for the core
  cells (`free` 9/10; `additive` 0/10; `conflict` 0/10) and GLM's additive cell
  was completed (`free` 5/5; `additive` 0/5; `conflict` 0/5). DeepSeek was
  attempted but unavailable, so it is not counted as evidence.
