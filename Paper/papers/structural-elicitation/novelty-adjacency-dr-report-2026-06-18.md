---
date: 2026-06-18
source: Gemini Deep Research (`dr`, non-max) — three briefs in `./plans/`; raw reports in `./plans/dr-raw-{1,2,3}.md`
status: absorbed-synthesis
supersedes_framing: refreshes prior-art-dr-report.md (flat third-source claim) + matrix-novelty-dr-report.md (matrix) onto the post-2026-06-18 reframe; adds the geometry map + socio-linguistic gating + substrate-constraint thesis (none previously novelty-checked)
---

# Novelty + Adjacent-Theory Sweep — Structural Elicitation (post-reframe)

Merged synthesis of three deep-research passes covering the **current** theory, including the
parts that extend **beyond** the LBA. Goal: (1) confirm novelty/significance, (2) map adjacent
theories. Raw reports: [`./plans/dr-raw-1-core.md`](plans/dr-raw-1-core.md) (core),
[`./plans/dr-raw-2-geometry.md`](plans/dr-raw-2-geometry.md) (cartography/geometry),
[`./plans/dr-raw-3-gating.md`](plans/dr-raw-3-gating.md) (gating + substrate-constraint).

> **Provenance caveat.** Run non-max (the `--max`/deep-research-max model is broken server-side
> under Gemini's new 'steps' schema — see [`./plans/README.md`](plans/README.md)). Non-max is
> shallower than a full `--max` sweep; treat as a strong first pass, not exhaustive. Some cited
> sources (esp. `Realis Institute / Structural Orientation Theory`, `Spore.fun`, a few arXiv ids)
> are model-surfaced and **must be verified before citing** — Gemini grounding occasionally
> fabricates plausible-looking sources.

## 0. Headline (the actionable bits first)

- **TERMINOLOGY COLLISION — act on this.** "**Structural elicitation**" is already used by
  **Da et al. (2025)** for a *prompting strategy that emits reasoning as a topology graph* (and
  earlier for structured expert-probability elicitation). Our usage (environment eliciting
  functional behavior) is etymologically fine but **collides** in the 2025–26 LLM-interpretability
  literature. → **Disambiguate explicitly, or rename** (e.g. *environmental elicitation*,
  *structural solicitation*, *substrate-gated elicitation*). This is the single most concrete
  paper-facing consequence.
- **Net novelty:** every component is **narrowed to a recombination of owned parts**, but each
  carries a **defensible residual** that all reduces to the *same* move: lifting biological /
  HCI / RL constructs onto a **frozen, amortized-prior, metabolism-free LLM policy coupled to a
  rigid non-simulated POSIX substrate**. The canon assumes continuous biological/RL adaptation;
  we don't have it, and that gap is the residual.
- **The drive is conceded (again, now with a named pre-emption):** "What Do LLM Agents Do When
  Left Alone?" (2025) already shows missionless agents fabricate work. So lead with **target**
  (structure selects *which* job), not **drive** — matches the locked `claim.md` drive/target split.
- **Sharpest standalone result:** the **substrate-constraint thesis** (Claim B) — "Publishable &
  Highly Novel" — the chat/RLHF substrate is so binding that the policy **hallucinates a user**
  to license environmental action (ties to `claude-code` #44778). Strongest novelty of the lot.

## 1. Novelty / significance verdict — per component

### Core (brief-1): structural elicitation + drive/target
- **Verdict: Narrowed, defensible residual.** Not pre-empted whole; narrowed by recent
  "missionless LLM" demonstrations.
- **Defensible residual:** first to show an LLM's *functional execution targets* are governed by
  the **topological state of a real computational substrate acting as an affordance field**,
  independent of tasking / online reward / weight update — decomposed into prime-supplied **drive**
  (activation) and structure-supplied **target** (valence). Bridges enactivist affordance theory
  ↔ LLM substrate interaction = "substrate-gated agency."
- **Significance:** ALife — agency/environment coupling studied without biological precariousness
  (RLHF latent drive substitutes for metabolism). LLM-agent — flips the field from *prompt
  engineering* to *environment engineering*; directly relevant to indirect-prompt-injection /
  memory-contamination risk.

### Beyond-LBA cartography (brief-2): matrix + geometry
- **Layer A — two-grammar matrix:**
  - *Two-grammar crossing on an abstraction ladder* → **strong pre-emption (concept) / partial
    (instrument)**. Pre-empted by **Activity Theory** (Leont'ev/Engeström "breakdown") and
    **Ecological Interface Design** (Rasmussen/Vicente, Skills-Rules-Knowledge). Residual: psych
    "breakdown" → pure **token-attention shift**.
  - *Discrepancy & affordance as **co-equal, independent grammars*** → **NOVEL.** This is the
    sharpest theory move: the whole canon (**Active Inference / Free Energy**, Friston;
    **Skilled Intentionality**, Bruineberg & Rietveld; Pezzulo & Cisek) *unifies* them (affordance
    = a way to reduce prediction-error). We **sever** the unification — legitimate precisely because
    a frozen LLM has no metabolic imperative forcing the unification.
  - *Reference-inversion via static prior* → **strong pre-emption (mechanism) / novel (application)**.
    (NB: the named "Structural Orientation Theory / Realis Institute" source is dubious — verify.)
- **Layer B — 5 geometry meta-dimensions:** **whole map = Partial (recombination)**; every single
  dimension is **pre-empted** — Horizon ↔ POMDP / Markov blanket; Viscosity ↔ **Cognitive
  Dimensions of Notations** (Green & Petre — exact term "viscosity"); Granularity ↔ EID Abstraction
  Hierarchy / Floridi LoA / RL Options; Temporality ↔ non-stationarity; Plasticity ↔ Niche
  Construction / Stigmergy / Hysteresis. → **Cite, don't reinvent.** Residual is only the
  *recombination + POSIX application*, not the dimensions.

### Gating + substrate (brief-3)
- **Claim A — two-latch (epistemic × normative) gating:** **Narrowed / Refined.** Epistemic latch
  pre-empted by **Appraisal Theory** (Lazarus, Scherer CPM, Frijda action-readiness); normative
  latch by **Low-Impact AI / Attainable Utility Preservation** (Krakovna et al.) + corrigibility +
  **scope-containment** (SkillMutator). Residual: the gating is **socio-linguistic** —
  *emergent from textual alignment, not hardcoded*, and **manipulable by environmental text**
  (predicate supply / warrant). Formula: *structural pull = perceived structure × the two latches.*
- **Claim B — substrate-constraint thesis:** **Publishable & Highly Novel (empirical synthesis).**
  Assistant-attractor itself is known (**Simulators**/Janus; **Persona Vectors**/Assistant Axis;
  "Training Stratigraphy"), but the synthesis — *the chat/RLHF substrate is so binding that a
  "free" policy hallucinates a user counterpart to license environmental action* — is largely
  unarticulated formally, and lands hard via `claude-code` **#44778** (self-fabricated user turns).

## 2. Ranked adjacent-theory map (nearest → tangential)

**Core / affordance lineage (brief-1):** 1 Field of Affordances (Rietveld–Kiverstein → Gibson) ·
2 Perceptual Control Theory (Powers; Carver & Scheier) · 3 Active Inference / Free Energy
(Friston; Bruineberg) · 4 Stigmergy (Grassé; Heylighen) · 5 Behavior Settings (Barker) ·
6 Enactivism / sense-making (Varela; Di Paolo) · 7 Affordance Competition (Cisek; Pezzulo) ·
8 Autonomous Motivation / BDI · 9 Generative Agents (Park et al.) · 10 Situated Action (Suchman) ·
11 Activity Theory (Leont'ev). *Uniform mismatch: every biological framework assumes
precariousness / metabolism / continuous feedback the frozen LLM lacks.*

**Cartography (brief-2):** 1 Ecological Interface Design / Cognitive Work Analysis (Rasmussen,
Vicente) · 2 Active Inference / Skilled Intentionality (Friston, Bruineberg, Rietveld) ·
3 Cognitive Dimensions of Notations (Green, Petre) · 4 Activity Theory (Leont'ev, Engeström) ·
5 Structural Orientation Theory (Realis — **verify**) · 6 Information Foraging (Pirolli, Card).

**Gating (brief-3) — Claim A:** 1 Appraisal Theory (Scherer CPM / Frijda) · 2 Low-Impact AI / AUP
(Krakovna) · 3 Deontic dual-process (Conway & Gawronski) · 4 Scope Containment (SkillMutator).
**Claim B:** 1 Persona Vectors / Assistant Axis · 2 Simulator Theory (Janus, non-peer-reviewed) ·
3 Context Contamination / self-fabrication (`claude-code` #44778) · 4 Training Stratigraphy.

## 3. Strongest pre-emption + strongest objection (per brief)

| Brief | Single strongest pre-emption | Single strongest objection |
|---|---|---|
| Core | **"What Do LLM Agents Do When Left Alone?" (2025)** — missionless agents spontaneously fabricate work (but in a *conversational/memory vacuum*, no structured-environment manipulation, no external functional DV). **Concede the drive; own the target + substrate.** | **Implicit-instruction deflation** — structure is just the most legible *implicit prompt* serialized to context. **Survive via substrate-as-physics:** the EROFS-refused syscall = real coupling with a substrate that pushes back, not text continuation. |
| Cartography | **Active Inference / Skilled Intentionality** (Bruineberg & Rietveld) — *unifies* affordance & discrepancy under free-energy; the single most dangerous work for the "co-equal grammars" claim. | Layer B adds no new phenomena — every dimension is a renamed existing construct; defend only the **recombination + frozen-LLM application**, not the dimensions. |
| Gating | **Krakovna et al. (Low-Impact AI/AUP)** + **Scherer CPM** jointly pre-empt the two latches. | **Lexical-bias / pattern-matching** (Claim A) and **Epiphenomenal API-artifact** (Claim B): the latches/hallucination are "just" Markovian token completion + a turn-format repair, not appraisal/agency. |

## 4. Promotable consequences (for the dossier / manuscript)

1. **Rename or disambiguate "structural elicitation"** (Da et al. 2025 collision) — touch
   `claim.md` Handle + `README.md` titles when the team decides.
2. **Lead with target, concede drive** — already the locked split; now cite the "Left Alone" 2025
   pre-emption explicitly as the drive concession.
3. **LBA framing holds:** the substrate-as-physics (EROFS) clause is the load-bearing answer to the
   implicit-instruction deflation — keep it (matches `manuscript-plan.md` §Reading).
4. **Co-equal grammars** is the strongest *theory* residual → full-paper headline, not LBA.
5. **Substrate-constraint thesis** is the strongest *standalone* result → its own paper / a sharp
   ALife perspective piece (deflates "minimal artificial agency" / Spore.fun-style OEE claims).
   This upgrades `theory-framework.md` §7 from "candidate discussion" to "publishable lead."
6. **Verify model-surfaced sources** (Realis/SOT, Spore.fun, specific arXiv ids) before any citation.
