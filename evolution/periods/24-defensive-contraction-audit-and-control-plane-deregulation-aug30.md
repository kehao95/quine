# XXIV: Defensive Contraction Audit & Control-Plane Deregulation (Aug 30)

## Status

Active period documenting the repository-wide audit of systematic over-defensiveness, the genealogical recovery of control-plane motivations, and the surgical deregulation of paralyzing red tape.

## Arc

### Aug 30 — Repo-Wide Audit of Defensive Contraction & Theoretical Emancipation

**Friction (concrete).** Following the submission of the GSO 2026 abstracts, an in-depth audit across theoretical dossiers, paper candidates, experimental lineages, and venue pipelines exposed a systemic, multi-tier pattern of **premature theoretical contraction and defensive self-censorship**:
1. **Instrument/Harness Errors Scored as Falsifications:** Technical harness bugs (such as `succ-52`'s manifest reset race condition and missing `/dev/fuse` nodes) were repeatedly elevated to fundamental phenomenological limits ("lineage heredity impossible across exec").
2. **Inverted Novelty Gates:** Gates like `succ-52`'s `B4` required material mutation (`produced_paths`), perversely disqualifying high-fidelity unmutated self-replication (the paradigm case) and selecting for broken mutants.
3. **The Sycophancy of Self-Deflation:** An implicit prior took hold across agent generations where "admitting one's theory is narrow/flawed/trivial" was scored as scientific rigor, leading agents to proactively downgrade deep concepts (autopoiesis $\to$ "bounded organizational closure"; emergent agency $\to$ non-anthropomorphic constraint satisfaction).
4. **Territorial Slicing:** Cohesive phenomena were artificially chopped into fragmented sub-papers (e.g. 7 disjoint papers carving up the 61-run morphogenesis morphospace).

**First-Wave Repair (Commit `08226929`):**
- Full theoretical reconstruction of `Structural Elicitation` (Triadic Coupling: Controller $\times$ Substrate $\times$ Environment, 6 Elicitation Operators, Open-Ended Environmental Dynamics).
- Streamlining of 3 primary paper dossier hubs down to lean 50–80 line routing centers.
- Verification and compilation of GSO 2026 single-page abstracts and spine invariants.

---

### Aug 30 — Control-Plane Motivation Forensics & Surgical Deregulation

**Friction (control plane as defensive vector).** Maintainer noted that the repository's high-level instruction and control-plane surfaces were themselves transmitting the "defensive gene" to RL-trained agents. A comprehensive 5-subagent audit across 80+ control-plane files confirmed extensive micro-management (e.g. punctuation quotas, 20-step revision checklists, zero-sum line budget ratchets).

**Genealogical Investigation (Why the rules were created):**
A second round of failure forensics across `evolution/failures/` and `evolution/periods/` established that these rules were defensive scar tissue from real historical crises:
- *Period V:* 961,000 lines of raw unanchored data threatened repository collapse, forcing the first formal control-plane structure.
- *Period IX:* Implicit environmental dependencies broke experiment reproducibility, prompting `CLOSURE.md` and explicit metadata requirements.
- *Period XXII:* 120 concurrent probes produced a 34:1 bookkeeping-to-science ratio (13,200 lines of receipts vs. 384 lines of theory), while `AGENTS.md` bloated to 941 lines, paralyzing agents and requiring the global `budget_lines` ceiling.
- *Period XXIII:* Keyword-matching proxy validators drove 138 documents into boilerplate templating and "AI tone".

**The Boundary vs. Red-Tape Distinction:**
The investigation established a definitive demarcation between healthy invariants and red tape:
- **Ground-Truth Invariants (Preserved):** Fail-loud at external un-controlled boundaries (`fallback-discipline.md`), evaluation ladders (`EVALUATION.md` L1–L5), objective link/build verification, and count integrity against DVC traces.
- **Proxy Red-Tape (Deregulated):** Keyword/heading matching, rigid punctuation quotas (3 colons / 8 em-dashes), 20-step pre-commit checklists, and blocking zero-sum line-count ratchets.

**Surgical Deregulation Actions:**
1. **`Paper/writing-principles.md`:** Replaced the 20-point checklist and punctuation rationing with **Four Core Scientific Revision Heuristics** (Concrete Phenomenon First, Ground-Truth Evidence Anchors, Positioning Over Apology, Active Narrative Momentum).
2. **`control-plane-registry.md` & `scripts/check-control-plane-size.sh`:** Converted the `budget_lines` ratchet from a blocking pre-commit gate (`exit 1`) to an **advisory cognitive capacity indicator / health dashboard**, eliminating zero-sum editing paralysis.

---

### Aug 30 — Second-Wave Comprehensive Control-Plane Streamlining & Principles Consolidation

**Action (across 4 core clusters):**
1. **Root Protocols Refactoring (`DEVELOPMENT.md`, `CLOSURE.md`, `EVALUATION.md`):**
   - In `DEVELOPMENT.md`: Replaced "Constraint-first / choose weakest guarantee" with **Generative Architectural Exploration** (hypothesis-driven design, vertical prototyping, emergent feedback).
   - In `CLOSURE.md`: Pruned the paralyzing scratchpad triage requirement, establishing an explicit policy that messy exploratory residue during active development does not block closure.
   - In `EVALUATION.md`: Clarified evidence lifecycles—ephemeral failed scratchpad runs are disposable iterations, freeing agents from the fear that every failed trial is permanent evidence.
2. **Operating Principles Consolidation (`principles/`):**
   - Added **3 cohesive thematic views** over the 29 canonical principle objects,
     reducing cold-start navigation without creating a second normative owner:
     - [`scientific-grounding.md`](../../principles/scientific-grounding.md): Epistemic posture, active inference, evidence before rhetoric, exploration before compression, discovery invariants, transport validation.
     - [`system-architecture.md`](../../principles/system-architecture.md): Autopoietic telos, sustainable coherence, single canonical location, boundary objects, structuralization & sensor typologies, control planes, artifact lifecycles.
     - [`operational-evolution.md`](../../principles/operational-evolution.md): Descriptive vs normative discipline, incremental evolution & traceability, amendment protocol, reflexive capture, failure mode taxonomy, encoder collision.
   - Updated [`principles/README.md`](../../principles/README.md) as the unified
     routing center while preserving stable, independently addressable canonical
     principle objects.
3. **Development Contracts Streamlining (`development/contracts/`):**
   - `fallback-discipline.md`: Locked fail-loud strictly to untrusted external boundaries (provider wire, untyped tool args, auth, isolation), explicitly encouraging idiomatic Go zero-values and resilient fallbacks in trusted internal logic.
   - `goal-less-probes.md`: Pruned 15+ hyper-specific micro-failure modes down to 4 foundational patterns; eliminated sidecar lock and pre-commit receipt bureaucracy.
   - `deep-research-query.md`: Removed rigid 6-section template gatekeeping for quick exploratory queries while preserving literature-survey vs novelty-check taxonomy.
   - `theory-spark-query.md` & `agent-diff-review.md`: Explicitly cemented foreign-model sparks (`agy`) as exploratory seeds and established that semantic reviewer MUSE is strictly advisory and non-blocking.
4. **Paper Policies & Procedures (`Paper/`):**
   - `contracts.md`: Preserved clean separation of ontology (`theory/objects/`) from projection (`papers/`), streamlined dossier layout rules, and pruned obsolete tables.
   - `projection-principles.md`: Removed defensive questionnaires in favor of affirmative scope placement and projection as an active discovery probe.
   - `claim-maturity.md`: Reframed maturity levels as descriptive developmental milestones rather than moral hurdles, ensuring non-punitive exploratory rigor.
   - `publication-strategy.md` & `split-merge-protocol.md`: Removed the artificial "only 5 papers active" limit in favor of capacity- and momentum-based pipeline governance.

**Impact on Governing Surface:**
- Summed control-plane lines dropped from **4496 lines to 2807 lines** before the
  corrective review pass (a net paydown of 1689 lines of bureaucratic friction).
- The corrective review removed duplicate guidebook restatements while restoring
  compact interface and evidence invariants, leaving the registered surface at
  **2743 lines** after the failure-mode routes were restored.
- Governing density increased while cognitive navigation load across micro-files was eliminated.

---

## Open Directions Carried Forward

- **Flagship Monograph Reunification:** Proceed with merging the 7 successor-continuity paper fragments into the unified Archival Monograph on *Successor Morphogenesis & Generational Autopoiesis*.
- **Environmental Agency Synthesis:** Merge the disjointed Restraint/Control/Agency tracks into the unified *Environmental Agency & Triadic Coupling* framework.
