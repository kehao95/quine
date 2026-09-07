## Period XIII: Repository Self-Model and Philosophy Layering (Apr 27 – May 1)

### Apr 27 — Theme theory-map extraction
**Key commits:** `d0ff02eec`
**What changed:** Extracted `autopoiesis-renewal` family theory into `Paper/themes/autopoiesis-renewal/THEORY.md`, slimmed the theme README back into an entry and routing surface, and converted theme `plans/`, `references/`, `analysis/`, and `papers/` indexes into role-aware support surfaces. Propagated the reusable pattern upward by documenting optional theme-level `THEORY.md` surfaces in `Paper/themes/README.md` and `Paper/README.md`.
**Friction:** The theme README had become a manuscript-sized theory synthesis while also serving as navigation, routing, red-line discipline, and paper split control. The content was valuable, but the container held too many jobs.

---

### May 1 — Philosophy layer extraction and self-maintenance frame
**Key commits:** `f2e999c8f`, `39acd5973`, `cff9df90b`
**What changed:** Created `Paper/philosophy/operating-philosophy.md` as the runtime-facing bridge from stable process ontology to the Plan 9 / Unix operating pull. Created `Paper/philosophy/self-maintaining-systems.md` as the higher-level frame connecting runtime, repository, and research as related self-maintaining systems. Routed README, Paper, core, roadmap, development, evaluation, runtime-ontology, closure, and AGENTS surfaces through the new layering. Added the explicit distinction that repo-level directional self-model commitments are normative, not publication claims.
**Friction:** The same Plan 9 / public-surface / self-maintenance principles were active in runtime docs, development review, roadmap, AGENTS.md, and autopoiesis/repository-self-model surfaces without one philosophy owner. Paper-claim caution also leaked upward into repo-level direction, making strong self-maintenance pull look like an overclaim rather than a normative organizing commitment.
**Structural repair:** Split the philosophy stack into progressive layers: north star, ontology, self-maintaining systems, operating philosophy, development vector, and implementation-backed core. Downstream docs now apply or route to those layers rather than re-owning them.
**What it taught:** Self-maintaining systems need strong normative direction and bounded descriptive claims at the same time. Entropy reduction improves when the repo can say what it is becoming without pretending every layer has already landed as implementation fact or paper evidence.

---

### May 1 — Activation kernel and maintenance trigger
**Key commits:** `053af7931`
**What changed:** Rewrote the AGENTS.md bootstrapper into an activation kernel that makes coherence maintenance the default session stance. Added the canonical `Maintenance` session-start keyword for a pure self-maintenance pass with no feature, paper, or research task attached.
**Friction:** The repo had articulated self-maintenance philosophically, but starting a session still depended on Human supplying an object-level task or repeatedly reminding the agent to maintain the control plane. The self-maintaining behavior existed as expectation rather than an explicit activation surface.
**Structural repair:** Made self-maintenance callable as a first-class session mode: sense current state, name bounded entropy, repair the smallest owning surface, validate, commit or hand off, and close through the existing closure protocol.

### May 1 — Maintenance protocol density correction
**Key commits:** `85f722933`
**What changed:** Reframed maintenance as investigation-first coherence work rather than a patch trigger. Kept `Maintenance` as the pure unscoped keyword, but made explicit maintenance phrasing route to the same diagnostic stance. Clarified that smallest bounded repair belongs after diagnosis, not before inquiry.
**Friction:** The first activation repair overfit a spelling edge case and reinforced shallow micro-fixes. The control plane needed a denser principle: maintenance should first read a structure deeply enough to name its coherence state, then repair or hand off.
**Structural repair:** Raised the rule from operational wording to a maintenance philosophy: diagnosis precedes repair; Human discussion precedes policy or structural changes when multiple futures are plausible.

---

### May 1 — Repository-self-model exploration absorption
**Key commits:** `5d33a190f`
**What changed:** Absorbed two repository-self-model deep-research passes into theme-local `plans/` and `references/`, created a merged deep-research synthesis, promoted provisional theory positioning into `THEORY.md`, and repaired AGENTS / DEVELOPMENT / Paper research rules so theme-level theory exploration stays generative until a branch becomes a paper, public, or implementation claim.
**Friction:** Paper-facing overclaim discipline leaked upward into theory exploration, prematurely narrowing the possible readings of the repository as organism-like process, artificial autopoiesis, hybrid cognitive system, developmental record, or non-biological morphogenesis.
**Structural repair:** Made exploration-before-defensibility a control-plane rule and preserved the research intake in durable theme surfaces.

---

### May 1 — Control-plane causal map
**Key commits:** `371966b54`
**What changed:** Added `Paper/themes/repository-self-model/plans/control-plane-causal-map.md`, a coupling ladder and map from major control-plane surfaces to readers, behavior effects, validation paths, and drift signals. Routed it from repository-self-model `THEORY.md` and `plans/README.md`.
**Friction:** Rich self-description risked becoming memory archive without an explicit account of how surfaces causally affect future agent behavior, validation, routing, closure, or Human judgment.
**Structural repair:** Turned causal coupling itself into an inspectable planning surface, sharpening the repo-self-model test for whether a surface belongs to the self-model.

---

### May 1 — Evolution journal organ split
**Key commits:** `159de8deb`
**What changed:** Split the monolithic `EVOLUTION.md` into a root navigation surface plus `evolution/periods/`, `evolution/patterns.md`, and `evolution/arc-index.md`. Preserved `EVOLUTION.md` as the canonical entrypoint while moving chronology, pattern synthesis, and thematic friction-repair evidence into separate organs.
**Friction:** `EVOLUTION.md` had grown into one container carrying current navigation, open friction, chronological archive, meta-pattern synthesis, and repository-self-model evidence. The content was valuable, but retrieval and role clarity were degrading.
**Structural repair:** Separated memory roles by access pattern: root for current state, period files for chronology, patterns for cross-arc synthesis, and arc index for theory-facing evidence.

---

### May 1 — Causal path annotation contract
**Key commits:** `235bd328b`
**What changed:** Added a Control-Plane Causal Path Contract to `DEVELOPMENT.md`, naming reader / effector, behavior effect, perception / validation, and absorption owner as the default causal-path fields for new or materially changed control-plane surfaces. Updated the repository-self-model causal map to route the reusable rule back to DEVELOPMENT.
**Friction:** The causal map made a deeper defect visible: without a repo-wide development rule, future control-plane surfaces could keep adding rich prose without naming how that prose changes future behavior or how drift against it is sensed.
**Structural repair:** Promoted causal coupling from a theme-local insight into a general development contract while preserving the detailed repository-self-model map as the application surface.

---

### May 1 — Development contract organ split
**Key commits:** `dcc13c1df`
**What changed:** Created `development/` with `contracts/` as addressable support surfaces for detailed development contracts. Slimmed `DEVELOPMENT.md` from a full contract library into the canonical entrypoint plus short causal-path stubs for deep research queries, structural friction journaling, commit messages, focus compression, code maintenance, causal paths, and active-doc addressability.
**Friction:** `DEVELOPMENT.md` had started to mirror the earlier `EVOLUTION.md` container defect: the root entrypoint was also carrying long reusable contracts, making navigation and causal ownership heavier than necessary.
**Structural repair:** Split detailed contracts into local organs while preserving stable root headings and validator expectations, so future agents can route by change type without losing direct access to full rules.

---

### May 1 — Experiment protocol navigation organ
**Key commits:** `1af64aa97`
**What changed:** Added `experiments/protocol/README.md` as a support map for the large experiment protocol, naming each protocol section's job, split posture, and causal path. Linked it from `experiments/PROTOCOL.md` without moving validator-bound rules.
**Friction:** `experiments/PROTOCOL.md` is large and multi-role, but immediate splitting would risk breaking validator-aligned rules. The repo needed better navigation before deeper surgery.
**Structural repair:** Added a navigation organ that lowers lookup cost and defines the conditions under which future protocol sections may be safely split.

---

### May 1 — Core references section split
**Key commits:** `d82e84600`
**What changed:** Split the 1000+ line `Paper/core/references.md` into a stable entrypoint plus section-level notes under `Paper/core/reference-notes/`. Updated the core README to expose the reference layer as its own support surface.
**Friction:** The core reference file mixed bibliography, prior-art threat assessment, vocabulary notes, conceptual comparisons, and deep-research reports in one large container. That made reference lookup expensive and made future additions likely to deepen the monolith.
**Structural repair:** Preserved `references.md` as the canonical link target while moving detailed reference material into addressable organs by section.

---

### May 1 — Autopoiesis-renewal reference surface map
**Key commits:** `d40b250f8`
**What changed:** Added `Paper/themes/autopoiesis-renewal/references/reference-map.md` and routed paper-local reference READMEs through it. The map distinguishes shared syntheses, thematic scans, paper-local reports, and lossless archives before future agents split or extend large reference files.
**Friction:** The autopoiesis-renewal theme had several large reference surfaces, but file size alone was not enough to decide whether to split them. Archives, retained report bundles, and active syntheses needed different structural treatment.
**Structural repair:** Made the reference layer's causal path and split posture explicit, preserving provenance-heavy archives while giving future work a clear trigger for when a synthesis should become section-level organs.
