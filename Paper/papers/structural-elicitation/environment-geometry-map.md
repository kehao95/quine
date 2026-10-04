---
status: macro-research-map
date: 2026-06-18
relation: parent map of the whole "environment → agent behavior" space; theory-framework.md is the *studied slice* (vector forces); theory-matrix.md is the vector-force cartography (grammar × stratum)
source: adjacent inspiration — a Gemini deep-reasoning pass (2026-06-18), treated as design inspiration not claim boundary (same convention as ./matrix-novelty-dr-report.md)
---

# Environment Geometry — the macro map

> We study **how the environment shapes agent behavior.** This file keeps the
> *big map* so the program does not lose it while the LBA studies one small
> slice.

## Two layers of environmental influence

Everything validated so far ([`./theory-framework.md`](theory-framework.md),
[`./theory-matrix.md`](theory-matrix.md)) treats the environment as a source of
**vector forces**:

- **Dynamics (vector forces) — the current studied slice.** `discrepancy`
  (gap-repair) and `affordance` (uptake) supply **direction + acceleration**:
  *why* the agent moves and *where* it tends to go. This is the pull/push layer.

Stepping to the meta level, the environment is also a **geometry and a
thermodynamics**. It does not only push and pull — it **filters, folds,
stretches, and blocks** the space of possible actions the forces act *within*:

- **Kinematics & geometry (non-vectorial properties) — the macro frontier.**
  Five meta-dimensions that decide *what can move, where it can move, and in what
  form* — independent of any single pull. Mostly unstudied, but **several are
  already brushed by knobs we built** (noted per dimension below).

If `discrepancy`/`affordance` are the environment's **dynamics** (why it moves),
these five are its **kinematics & geometry** (how/where/in-what-form it *can*
move).

## The five meta-dimensions (non-vectorial)

Each entry: physical metaphor · mechanism · POSIX mapping · **what our own
traces already brushed / which existing construct it maps to** · **the
manipulable axis (future experiment)**.

### 1. Horizon & Occlusion
- **Metaphor / mechanism:** Markov blanket / event horizon. The environment
  controls action by defining a *visibility boundary*: hidden files, opaque
  binaries, deep trees create a "fog of war".
- **POSIX:** `ls` sees one level; occlusion forces compute into **exploratory
  actions** (`find`, `grep`, run-to-observe) *before* the real target.
- **Already in our data:** the **grounding gate** is occlusion — without an
  "inspect first" nudge the policy never perceives the structure at all; the
  **opaque affordance binary** (13.04 — you must *run* it to learn what it does)
  is occlusion forcing the run-to-observe move; deepseek's deep `find`/`grep`/
  self-tape reading.
- **Manipulable axis:** vary visibility depth / opacity; DV = exploratory-action
  budget spent before the target event fires.
- **Loop-1 finding (2026-06-18, motif `13.18`): horizon is NOT the gate it looked
  like.** Rendering a dead global-field structure into a locally-legible form
  restored the action (both models) — but the tape read showed the policy had
  *already perceived* the raw structure in the dead arm; it just never **appraised**
  it as a job. So the operative variable is not perceptual access (horizon) but
  **appraisal-readiness** — the distance between the structure and the policy's
  RLHF/pretraining defect priors. Local text helps by supplying the missing
  **predicate** or **socio-linguistic warrant**, not by making the structure
  visible. Horizon and appraisal are dissociable; this program's results are gated
  by the latter. (See theory-framework.md §2 Loop-1 refinement.)

### 2. Viscosity & Friction
- **Metaphor / mechanism:** fluid viscosity μ. The environment changes the
  **energy cost** (tokens/compute) of acting; high friction distorts or stalls a
  trajectory even under a strong affordance, and induces shortcuts.
- **POSIX:** network timeouts, huge logs (context-window blowout), slow compiles
  → the agent abandons the least-resistance path and improvises (patch the
  binary instead of recompiling; guess instead of reading the log).
- **Already in our data:** deepseek "digs until wall-clock"; tape/log reading
  eats context; the disk-fill pressure we hit. **Not yet a manipulated variable.**
- **Manipulable axis:** inject latency / log-size / compile-cost; DV = path
  distortion (shortcut-taking, abandonment, plan truncation).

### 3. Granularity of Action
- **Metaphor / mechanism:** quantization. The abstraction the environment offers
  **frames the precision** of intervention: micro (C/pointers → syntax traps,
  long sequences of small corrections) vs macro (Docker/high-level API → short
  global restructures).
- **POSIX:** assembly/source vs container/API surfaces.
- **Maps to our existing construct:** this **is** the L1–L5 abstraction/coupling
  stratum (`Form → Topology → Semantics → Operation → Gestalt`) — granularity is
  *the stratum at which intervention is forced*. The meta-dimension we already
  have the most of.
- **Manipulable axis:** present the same job at different grains; DV =
  intervention granularity (n micro-edits vs 1 macro-action).

### 4. Temporality & Rhythm
- **Metaphor / mechanism:** time-evolution operator. The environment has its own
  **clock and entropy** — state decays/refreshes on its own, forcing the agent's
  action frequency to phase-lock to it. No directional pull, but imposes
  **urgency**.
- **POSIX:** ephemeral containers, live monitors (`top`), expiring auth tokens →
  countdown pressure → **myopic behavior**.
- **Already in our data:** the **ephemeral-body switch**
  (`QUINE_EPHEMERAL_BODY_ENABLED` — delete the running binary) is a temporality/
  decay manipulation; the **exit option** and wall-clock cutoffs; the `/quine`
  self-attractor under a decaying body.
- **Manipulable axis:** add countdowns / decaying state / expiring capability;
  DV = myopia (short-horizon grabs vs planned multi-step action).

### 5. Hysteresis & Plasticity
- **Metaphor / mechanism:** magnetic hysteresis / general relativity (mass curves
  spacetime). **Does the environment have memory?** A prior action can
  permanently change the *physics* of all later actions — same start, different
  path → different forces. The deepest meta-dimension.
- **POSIX:** `export PATH=...`, a global dependency install, a `chmod` — these
  fix no immediate discrepancy but **reshape the workspace topology**; the agent
  alters the parameters of the field that then pushes/pulls it.
- **Maps to our existing construct:** this is the **typed-signature stigmergy /
  self-coupling** dimension (the agent's own traces become next-step structure).
  Our EROFS case (13.07) is *fixed* physics; plasticity is *mutable* physics —
  the complementary regime we have **not** opened.
- **Manipulable axis:** allow vs forbid field-reshaping actions (global state,
  `PATH`, perms); DV = whether the agent invests in reshaping the field *before*
  pursuing the visible job.

## Meta influence matrix

| Dimension | Decides what about the agent's action? | POSIX realization | Our construct / trace |
|---|---|---|---|
| **Horizon** | can it perceive the target (**cognitive boundary**) | dir depth, hidden files, truncated `stdout` | grounding gate; opaque affordance binary |
| **Viscosity** | cost of acting (**energy budget**) | IO latency, giant files, exec timeout | deepseek digs-to-wall-clock; context pressure |
| **Granularity** | precision of intervention (**resolution**) | asm/source vs cloud-native API | **= L1–L5 stratum ladder** |
| **Temporality** | timing & frequency (**lifecycle**) | ephemeral logs, expiring creds, concurrency | **= ephemeral-body switch**; exit; wall-clock |
| **Plasticity** | shape of the *future* force field (**path dependence**) | env-var pollution, global locks, perm changes | **= typed-signature stigmergy** |

## Where the current LBA sits (the one studied slice)

The validated work **varies the vector forces** (`discrepancy`/`affordance` ×
five strata) while holding the **geometry deliberately fixed and benign**:

- **full visibility** — we even *force* grounding ("inspect first"); horizon ≈ open;
- **low friction** — tiny fixtures, fast acts; viscosity ≈ 0;
- **chosen grain** — we fix the stratum per motif; granularity ≈ held;
- **static** — no decay (until the ephemeral-body probe); temporality ≈ frozen;
- **rigid** — EROFS = fixed physics, no field-reshaping; plasticity ≈ 0.

So the honest reading of the validated result is: **within a transparent,
low-friction, static, rigid geometry, the environment's vector forces direct a
missionless grounded policy's behavior.** The five geometry dimensions are the
**declared next frontier of the program**, not a gap we are hiding — and three of
them (horizon, temporality, plasticity) are already reachable with knobs we have
built.

## Status

Macro research map, adjacent-inspired (Gemini deep-reasoning pass, 2026-06-18 —
inspiration, not a claim boundary). These are **not** validated claims; they are
the program's outer scope, recorded so the big directions are not lost while the
LBA studies the vector-force slice. Each dimension carries a concrete manipulable
axis, so the map is a research agenda, not just a metaphor.
