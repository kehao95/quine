# Structural Evolution Journal

Current structural-friction index for Quine. Detailed chronology lives in
[`evolution/periods/`](./evolution/periods/); recurring synthesis lives in
[`evolution/patterns.md`](./evolution/patterns.md); friction-repair arcs live in
[`evolution/arc-index.md`](./evolution/arc-index.md).

Protocol: [`development/contracts/structural-friction-journal.md`](./development/contracts/structural-friction-journal.md).

## Current Navigation

- This file holds only unresolved friction as compact owner pointers.
- Period files own narrative chronology; resolved/parked entries leave this
  surface rather than accumulating annotations forever.
- Entries use `[owner: <token>; status: active; review: YYYY-MM-DD]` to support
  navigation and review. Dates are reminders; unresolved questions can remain.
- Drift check: `scripts/check-evolution-hygiene.sh` plus active-doc links.

## Open Friction

- Anti-premature-convergence discipline can suppress completion-driven discovery [owner: lab; status: active; review: 2026-07-29]: “keep options open” has no symmetric cost for never committing. Owners:
  [`projection-principles`](./Paper/projection-principles.md),
  [`computational-epistemology`](./Paper/philosophy/computational-epistemology.md).
  First hard cost recorded 2026-07-26: the AAAI-27 cycle passed with four
  configurations bound to it and zero submitted — three of them held open on
  gates (frozen basis, unrun A3 experiment, restart gate) that no deadline was
  allowed to override. Record:
  [`aaai-27`](./Paper/venues/objects/aaai-27.md) § Cycle outcome.

- The discovery loop produces theory but does not yet discriminate publication readiness [owner: lab; status: active; review: 2026-08-06]: missing axes include novelty, power, feasibility, venue fit, and differentiation. Repair design:
  [`publication-readiness-gate`](./Paper/_design/migrations/publication-readiness-gate.md).
  The AAAI-27 miss adds a missing axis the design does not yet name: **who is
  drafting, by when**. A venue binding in frontmatter was treated as a plan by
  every surface that read it, including the venue card's own "configurations
  currently targeting this venue" table.

- Theory maturity describes writing/relations more reliably than empirical anchoring [owner: lab; status: active; review: 2026-08-06]: navigation is stronger than value steering. The MDL/evidence readout now exposes resolved/open
  counts; remaining work is to couple visible signals to agent judgment without
  replacing that judgment. Owner:
  [`computational-epistemology`](./Paper/philosophy/computational-epistemology.md).

- Runtime-first progressive disclosure has one frozen-release residue [owner: lab; status: active; review: 2026-08-02]: public `quine-arxiv` still names retired `wisdom`; decide whether frozen releases need standing errata.

- System-prompt ordering contradicts its progressive-disclosure rule [owner: lab; status: closed]: both halves repaired. Active Constraints now render second, immediately after the Prime Directive (T5.1). The `exec` half is dissolved rather than patched: there is no longer a separate staged-config file to point at — one agent-authored `config/env/override` is what every sh child, fork child, and exec successor inherits, and the exec schema names that one file. Authority:
  [`env-process-boundary-brief`](./Paper/theory/views/runtime-capability/env-process-boundary-brief.md) (E2);
  repair history: [`prompt-progressive-disclosure-followups`](./Paper/_design/migrations/prompt-progressive-disclosure-followups.md).

- “Sandbox” overstates current isolation [owner: lab; status: active; review: 2026-07-29]: runtime FUSE projects public process state and overlay isolates the configured workspace, not the whole host root. Naming and tests
  must not imply isolation for `/etc`, package state, or unrelated host paths.

- Projection cards cannot detect a body that duplicates or contradicts its declared owner [owner: lab; status: active; review: 2026-07-29]. Repeated
  index drift should repair source/projection separation or add a derived view,
  not create another summary.

- Experiment interpretation can duplicate one read across family, analysis, and experiment surfaces [owner: lab; status: active; review: 2026-07-29]: the current validator does not cover every sibling body. Owner:
  [`experiments/PROTOCOL.md`](./experiments/PROTOCOL.md) § Phase Interpretation
  Layer.

- Meta-evaluators are canalized toward “accept with the smallest concession” [owner: lab; status: active; review: 2026-07-29]: both discipline vocabulary and imperative prompt shape contribute. Repair direction: express target-state invariants
  declaratively and test whether stance/owner/typed relations remain visible.

- Experiment-analysis projections weld phenomenon, reading, and theory routing [owner: lab; status: active; review: 2026-07-29]: later readers inherit interpretation as fact. Open choice: declare these
  files ownerless readings with a sharp stale boundary, or create addressable
  phenomenon objects. Neighbors:
  [`narrative-claim`](./Paper/theory/objects/narrative-claim.md),
  [`addressability-debt`](./Paper/theory/objects/addressability-debt.md).

## Chronological Periods

- Period records: [`evolution/periods/`](./evolution/periods/)
- Recurring patterns: [`evolution/patterns.md`](./evolution/patterns.md)
- Theory-facing repair arcs: [`evolution/arc-index.md`](./evolution/arc-index.md)

| Period | File | Entries |
|---|---|---:|
| I | [`01-runtime-features-and-experiment-sweeps-feb-25-26.md`](./evolution/periods/01-runtime-features-and-experiment-sweeps-feb-25-26.md) | 10 |
| II | [`02-design-documents-and-alife-paper-sprint-feb-27-mar-1.md`](./evolution/periods/02-design-documents-and-alife-paper-sprint-feb-27-mar-1.md) | 6 |
| III | [`03-paper-compilation-and-arxiv-submission-mar-2-3.md`](./evolution/periods/03-paper-compilation-and-arxiv-submission-mar-2-3.md) | 6 |
| IV | [`04-systems-paper-sprint-mar-4-7.md`](./evolution/periods/04-systems-paper-sprint-mar-4-7.md) | 8 |
| V | [`05-structural-reorganization-and-control-plane-creation-mar-8-10.md`](./evolution/periods/05-structural-reorganization-and-control-plane-creation-mar-8-10.md) | 28 |
| VI | [`06-runtime-stabilization-and-feature-expansion-mar-11-12.md`](./evolution/periods/06-runtime-stabilization-and-feature-expansion-mar-11-12.md) | 10 |
| VII | [`07-testing-infrastructure-and-paper-track-expansion-mar-13-17.md`](./evolution/periods/07-testing-infrastructure-and-paper-track-expansion-mar-13-17.md) | 19 |
| VIII | [`08-fork-architecture-autopoiesis-and-emergence-closure-mar-18-21.md`](./evolution/periods/08-fork-architecture-autopoiesis-and-emergence-closure-mar-18-21.md) | 25 |
| IX | [`09-control-plane-hardening-and-alife-manuscript-mar-22-23.md`](./evolution/periods/09-control-plane-hardening-and-alife-manuscript-mar-22-23.md) | 19 |
| X | [`10-alife-paper-strengthening-and-p6-transmutation-series-mar-27-apr-5.md`](./evolution/periods/10-alife-paper-strengthening-and-p6-transmutation-series-mar-27-apr-5.md) | 18 |
| XI | [`11-evaluation-migration-signal-discovery-and-paper-convergence-apr-6-11.md`](./evolution/periods/11-evaluation-migration-signal-discovery-and-paper-convergence-apr-6-11.md) | 32 |
| XII | [`12-alife-submission-self-source-and-structural-reflection-apr-12-25.md`](./evolution/periods/12-alife-submission-self-source-and-structural-reflection-apr-12-25.md) | 42 |
| XIII | [`13-repository-self-model-and-philosophy-layering-apr-27-may-1.md`](./evolution/periods/13-repository-self-model-and-philosophy-layering-apr-27-may-1.md) | 12 |
| XIV | [`14-self-maintaining-center-and-manifesto-retirement-may-5.md`](./evolution/periods/14-self-maintaining-center-and-manifesto-retirement-may-5.md) | 1 |
| XV | [`15-v0-living-runtime-spine-may-5.md`](./evolution/periods/15-v0-living-runtime-spine-may-5.md) | 1 |
| XVI | [`16-dvc-externalization-and-curated-public-release-may-11-27.md`](./evolution/periods/16-dvc-externalization-and-curated-public-release-may-11-27.md) | 5 |
| XVII | [`17-constitution-compression-jun-9.md`](./evolution/periods/17-constitution-compression-jun-9.md) | 1 |
| XVIII | [`18-theory-net-generative-substrate-jun-13.md`](./evolution/periods/18-theory-net-generative-substrate-jun-13.md) | 12 |
| XIX | [`19-theory-restructure-auto-discovery-jun-20.md`](./evolution/periods/19-theory-restructure-auto-discovery-jun-20.md) | 9 |
| XX | [`20-runtime-capability-topology-jun-22.md`](./evolution/periods/20-runtime-capability-topology-jun-22.md) | 7 |
| XXI | [`21-latent-substrate-interpolation-jun-24.md`](./evolution/periods/21-latent-substrate-interpolation-jun-24.md) | 3 |
| XXII | [`22-autopoietic-self-governance-and-paradigm-rerooting-jun25-jul01.md`](./evolution/periods/22-autopoietic-self-governance-and-paradigm-rerooting-jun25-jul01.md) | 13 |
| XXIII | [`23-validator-template-equilibrium-jul02-jul18.md`](./evolution/periods/23-validator-template-equilibrium-jul02-jul18.md) | 1 |
