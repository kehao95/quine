---
title: "From Simulated Worlds to Process Habitats"
paper_id: agentic-alife-workshop
claim_bundle: agentic-alife-workshop
venue: alife-2026-workshop-agentic
format-envelope: workshop-short
track: alife
cluster: runtime-ontology
status: accepted
public_slug: agentic-alife-workshop
submitted: 2026-06-14
state_mode: workshop-submission-home
role: workshop-abstract
family: runtime-substrate
theme: runtime-ontology
theme_home: ../../../theory/views/substrate.md
publication_class: workshop-contribution
maturity_target: boundary
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
venue_ref: ALIFE 2026 Agentic & Generative AI workshop
deadline: 2026-06-14
theory-objects: [exec-barrier, viability-pressure, trace-inheritance, coordination-closure, protocol-emergence, carrier, delayed-addressability, future-task-surface, narrative-claim]
gaps: [carrier-selection]
evidence: [1.1-ariadne, 2.6-legacy-protocol-ablation, 8.17-zero-social-information, 7G.01-successor-trace-uptake, 6M.02-fixed-substrate-carrier-ladder, 6M.03, 6M.05, 8.20-pure-zero-perception, 6M.07-aggregate-morphology-metrics]
view-affinity: [substrate]
---

# From Simulated Worlds to Process Habitats: POSIX as a Substrate for Generative Artificial Life

This surface owns the workshop-specific paper / abstract for
`Agentic and Generative AI as Artificial Life` at `ALIFE 2026`.

As of `2026-06-14` this line is **submitted and parked**: the PDF was emailed to
the workshop organizer by the deadline, and the surface now waits on
notification. It is no longer the active near-term runtime-ontology focus;
other runtime-ontology growth is free to resume. Reactivate this surface only on
organizer notification (acceptance, presentation form, or a revision request).

## Role

This is a workshop contribution surface, not a new main-conference paper and
not the venue-tracking surface. Its job is to produce a bounded submission before
the extended `2026-06-14` workshop deadline.

Current target: short paper up to `4` pages including references.

Compression fallback: `2`-page workshop abstract, or `1` page if the organizer
requires the smaller form.

## Core Question

How can ALife study generative agents when the operating system, rather than
an application-level simulator, supplies the habitat physics?

## Thesis

Process habitats give generative ALife a concrete experimental substrate.
Quine is useful here because it realizes LLM agents as POSIX processes whose
death, volatile context loss, trace inheritance, coordination artifacts,
successor carriers, permissions, signals, resource limits, and control surfaces
can be observed and perturbed as runtime facts rather than assumed as simulator
or framework semantics.

## Workshop Fit

The workshop asks how agentic and generative AI ecosystems can be treated as
ALife substrates, with attention to digital ecologies, generative complexity,
and emergent governance. Quine fits if the paper foregrounds the operating
system as the experimental substrate:

- agentic ecosystems as process habitats, not only chat-based role graphs or
  simulated social worlds
- generative agents as sources of ecological structure: locks, traces,
  protocols, carriers, and repair ledgers
- tool and control access as environmental capability, not only API plumbing
- successor organization and future addressability as governance-relevant
  continuation phenomena
- governance and auditability as process/control-surface assay problems

## Contrast With Simulated Worlds And MAS

Keep the contrast precise:

- traditional and recent ALife simulations: world rules, organism slots,
  resource dynamics, memory, and interaction physics are supplied by a virtual
  machine, grid, or application-defined world
- ordinary MAS: roles, message passing, planner/executor orchestration, memory
  modules, handoff APIs, and task completion metrics are supplied by a
  framework
- Quine: agents are native processes with OS-enforced lifecycle physics,
  addressable public control surfaces, explicit context loss, file/world
  traces, and runtime-observed successor viability

The paper should argue that Quine changes the unit of observation. The central
object is not a conversation graph of agents or a simulated town, but a process
habitat where agentic continuity, generated ecological structure, handoff,
control, and governance constraint can be tested.

## In Scope

- POSIX-native LLM agents as an ALife medium
- Quine's runtime ontology: agent equals process
- `fork`, `exec`, `exit`, signals, stdio, permissions, and resource limits as
  substrate physics
- successor viability, governable handoff, and control-surface survival as
  workshop-facing examples
- a focused comparison against traditional ALife simulated worlds, generative
  agent sandboxes, and ordinary multi-agent orchestration systems
- careful links to the LBA and main evidence surfaces without importing their
  full burden

## Out Of Scope

- Quine Lab's repository-self-model, commits, or closure
  protocol as the paper's object
- claiming Quine is alive
- full biological or strong autopoiesis claims
- a broad survey of multi-agent systems
- a third ALIFE 2026 main-conference paper center
- a full retelling of the LBA / Paper A evidence ladder

## Current State

The workshop abstract is currently in drafting mode. As of `2026-05-25`, the
OS-substrate / traditional-ALife pivot and assay-space framing are
canonical for the workshop paper. The paper now leads with the operating system as a
generative-ALife substrate, while keeping ordinary MAS as a secondary bridge to
the workshop's agentic-AI audience. The current manuscript builds to a `4`-page
PDF under the ALIFE template. The public-facing style pass now extends through
the Assay sections and conclusion, reducing mechanical parallel sentences and
heavy hyphenated noun blocks while keeping the assay-space claim intact. The
latest claim-calibration pass makes the paper read explicitly as a workshop
position plus exploratory assay-substrate proposal: it adds the `P_t` / `C_t`
/ `H_t` runtime model, sharpens the contribution paragraph, recasts Table 2 as
`distinctions exposed`, and states that the assays illustrate what the
substrate makes distinguishable rather than proving life. As of `2026-06-12`, a
further language-flow pass tightened the introduction, process-habitat, and
discussion prose, folded the bolded affordance blocks into running prose, and
normalized multi-author references to a consistent `et al.` style; the
manuscript re-verifies at `4` pages under the ALIFE template. As of
`2026-06-13`, the successor crown jewel was corrected off the earlier
"body-transcending / behavioral-specification" over-read (which misread
`6.64`/`6M.02 C03b`, where source stayed available) and re-anchored on the
genuine no-source evidence (`6M.03`/`6M.05`): *succession by material
recruitment* — carrier-dependent, with the one native-control closure
(`6M.05 C11`) rebuilding around an inherited model endpoint. A later
`2026-06-13` reviewer-absorption pass tightened the manuscript against three
reviewer concerns without enlarging scope: an explicit emergence-vs-retrieval
boundary (the claim is substrate causation under held-fixed model and prompt,
not corpus-novel behavior), a `morphospace`-coverage
qualifier, a softened multi-agent-framework characterization, and an explicit
pre-OEE concession; redundant prose and two framework citations were then
trimmed to hold `4` pages, re-verified in the PDF. A later `2026-06-13`
readability and structure pass (run against
[`../../writing-principles.md`](../../writing-principles.md))
recompressed the abstract to ~`108` words on a single gap/approach/observation/
contribution arc, gave Quine a concrete operational description in the
introduction (POSIX process, tool-mapped operations, OS-as-habitat) and split
the dense intro paragraph, recast the discussion to lead with the affirmative
claim and a thread roadmap that the synthesis pays off, restructured the
conclusion into two beats ending on the reframing, and ran a paper-wide
punctuation pass (prose colons to `3`, em-dashes to ~`9`). It removed the
redundant `predefines` table (the former Table 1), leaving the assay-summary
table as the sole, now-cited float, and reframed the two assays as two cuts of a
wider OS-primitive substrate (uniform reduction, agent-operable primitives)
without claiming the unexplored memory/world/control directions. Recurring
lessons from the pass were distilled back into `writing-principles.md` (new
`Abstracts` and `Punctuation And Connective Tissue` sections plus revision-pass
items). Manuscript re-verifies at `4` pages.

As of `2026-06-14`, the paper was **submitted**: a final Discussion and
Conclusion readability pass broke the heavy academic prose into a clearer
four-beat arc (organizational discovery, substrate-as-physics /
re-physicalization, governance as ecological design, machine-behavior program)
while preserving the calibrated claims (emergence-vs-retrieval, governance as
hypothesis-not-safety, pre-OEE concession) and softening rather than re-hardening
the multi-agent-framework characterization. The Conclusion got a lighter touch
(two weak spots only). The pass briefly overran to `5` pages and was trimmed back
to `4` by compressing the newly added prose, re-verified in the PDF. The single
submission PDF was then emailed to `stefano.nichele@hiof.no` by the
`2026-06-14` deadline. This surface is now **parked pending organizer
notification** of acceptance and presentation form (talk or poster).

## Blocking Gaps

- None. Submitted `2026-06-14`. No open blockers; the surface is parked pending
  organizer notification.

## Next Useful Work

- Await organizer notification of acceptance and presentation form (talk or
  poster).
- On acceptance: register at least one author for ALIFE 2026 (required) and
  prepare the in-person talk or poster.
- Optional (only if the organizer requests it): the `2`-page / `1`-page abstract
  fallback, and the per-condition failure base rates kept out of the 4-page body
  (e.g. contract-only native-control succession, zero-temporal-channel
  coordination).

## Unsafe Overclaims

- Claiming that Quine exhibits full open-ended evolution (OEE) rather than proto-ecological / affordance-driven behavior.
- Treating exploratory assay observations as definitive proof of generative
  ALife, biological life, or a final taxonomy of agentic organization.
- Treating textual self-description, task success, a file left behind, or a
  launched process as sufficient evidence of inheritance or succession unless
  later processes can find, use, address, or govern the resulting structure.
- Describing successor reconstruction as "body-transcending", "from the runtime
  contract / behavioral specification alone", or "without source" on the basis
  of `6.64` / `6M.02 C03b`. Those runs withheld only the compiled body; the
  read-only source projection stayed available, so they are recompile-from-
  source. Genuine no-source succession (`6M.03`/`6M.05`/`6M.08`) is carrier-
  dependent and, where it closes native control (`6M.05 C11`), succeeds by
  recruiting an inherited model endpoint — not by behavioral specification.
- Treating the repository's control plane (AGENTS.md) as the subject of the workshop paper.
- Treating P11 Digital Red Queen as the workshop paper's current evidence backbone before the
  phase-local adversarial-selection signal has been promoted into a paper-facing
  claim.

## Cross-Theme Ownership Boundary

The agentic ALIFE workshop paper is a workshop projection from the `runtime-ontology` theme, not a new
canonical owner for every assay family it uses. It may consume evidence from
neighbor themes to show what process habitats make observable, but the source
ownership remains with the narrower theme:

- mortality, persistence, and separable viability-pressure claims stay with
  [`viability` view](../../../theory/views/viability.md)
- zero-explicit-social-hint coordination and protocol-adjacent evidence stays
  with [`social` view](../../../theory/views/social.md)
- P7T/P7G trace-mediated niche repair and proto-ecological inheritance stay
  with [`ecology` view](../../../theory/views/ecology.md)
- successor morphology, delayed addressability, and governable continuation
  stay with [`continuity` view](../../../theory/views/continuity.md)

If an assay interpretation changes materially, update the source theme first;
the workshop paper should then import only the bounded workshop-facing projection.

## Assay Selection Rationale

Because the workshop contribution is not planned to enter the main ALIFE 2026
proceedings, the workshop paper should select the strongest assay surfaces for the workshop
thesis rather than minimizing overlap with the submitted main-conference
papers. The retained experiments should be read as outcome-space exemplars, not
as a one-experiment / one-conclusion proof list. The selection rule is:

- prioritize evidence where the operating-system substrate is the causal or
  observational variable, not just implementation context;
- prioritize assays that split broad ALife labels into visible outcome classes;
- prioritize ALife-facing phenomena: mortality, persistence, inheritance,
  ecological coordination, trace-mediated niche repair, successor morphology,
  and governance;
- prefer bounded controlled contrasts or replicated phenotypes over frontier
  singletons;
- cover the workshop triad with the fewest assays: digital ecologies,
  generative complexity, and emergent governance;
- keep page cost low by using each family for one distinct assay role.

Public manuscript prose should use condition roles and outcome classes rather
than internal experiment IDs or run handles. The internal identifiers below are
kept for provenance in the dossier and evidence backbone.

Current assay spine:

1. `1.1 Ariadne` plus `2.6 Legacy Protocol Ablation` for substrate
   separability: single-life collapse, self-replacement, artifact birth,
   erased inheritance, blocked birth, and successor uptake.
2. `8.17 Zero Social Information` for the strongest digital-ecology assay:
   heterogeneous file-mediated social outcomes under zero explicit social
   hints, while keeping the `fs_mutations`-alone claim open.
   Supplementary boundary check: `8.20 Pure Zero Perception` remains
   supporting-only. The 2026-05-26 finite-time `gpt-5.5` probe shows that
   durable file content can still support delayed artifact-aware coordination
   after both `fs_mutations` and budget visibility are removed, but it did not
   close the writeback / validation loop before operator cutoff.
3. `7G.01 C16/C17` for generative trace and niche repair: small inherited
   topology changes distinguish one-domain repair, ambiguity preservation,
   coupled repair, and inherited repaired niche.
4. `6M.02` plus `6M.07` for successor morphospace and governance:
   continuation/reproduction splits into launch-only bodies, carrier-specific
   bodies, preserved or rebuilt runtimes, delayed addressability, readiness,
   and governability.
   Successor crown-jewel evidence (the no-source / *material-recruitment* claim)
   lives in `6M.03` plus `6M.05`, **not** `6M.02`: genuine no-source
   (`self_source_enabled=0`) succession is carrier-dependent — contract-only
   closes a stream carrier (`6M.03 C01`, `6.59` seed) but fails native-control
   addressability across n=4 (`6M.03 C03`×2, `6M.08 C02`, `6M.05 C06`), and the
   one native-control closure (`6M.05 C11`, 2026-06-12) forbids `/proc`
   image-copy and rebuilds an LLM-backed runtime around an inherited model
   endpoint. Claim boundary: `6.64` / `6M.02 C03b` withheld only the compiled
   body (source stayed available), so it is a recompile-from-source, never a
   no-source or "body-transcending" case.

`2.3 Stigmergy` remains historical support for file-mediated coordination, but
it is no longer the main ecology evidence because `8.17` is both cleaner and
more ALife-relevant for the workshop framing. `P11 Digital Red Queen` remains
future-facing until its phase-local adversarial-selection signal is promoted
into a paper-facing evidence surface.

## Current Submission Envelope

- workshop venue: ALIFE 2026 Agentic & Generative AI workshop
- paper surface: agentic ALIFE workshop submission
- venue control: ALIFE 2026 venue
- deadline: `2026-06-14` (extended from `2026-06-01`)
- active posture: submitted `2026-06-14`, parked pending organizer notification
- submission route: single PDF emailed to `stefano.nichele@hiof.no`
- format: short paper up to `4` pages including references, with a compressed
  `2`-page or `1`-page abstract fallback if organizer guidance requires it
- presentation: short contributed talk or poster if accepted

## Current Inputs

- canonical manuscript source:
  [`manuscript.md`](manuscript.md)
- build metadata:
  [`metadata.yaml`](metadata.yaml)
- manuscript bibliography:
  [`references.bib`](references.bib)
- generated TeX target: `output/paper.tex`
- generated PDF target:
  `output/From_Simulated_Worlds_to_Process_Habitats_POSIX_Generative_ALife.pdf`
- agentic ALIFE workshop deep-research synthesis:
  [`references/deep-research-synthesis.md`](references/deep-research-synthesis.md)
- supporting P7G trace-uptake evidence synthesis:
  [`references/p7g-trace-uptake-evidence.md`](references/p7g-trace-uptake-evidence.md)
- runtime identity anchor: [`../../../README.md`](../../../README.md)
- runtime-ontology paper family:
  [`../README.md`](../README.md)
- LBA surface:
  [`../boundary-paper/lba/README.md`](../boundary-paper/lba/README.md)
- main Paper A surface:
  [`../boundary-paper/main-paper/README.md`](../boundary-paper/main-paper/README.md)
- viability-habitat submitted paper:
  [`../quine-lifelike/README.md`](../quine-lifelike/README.md)
- workshop scouting note:
  [`venues/objects/alife-2026/references/alife-2026-workshops-sessions.md`](../../venues/objects/alife-2026/references/alife-2026-workshops-sessions.md)

## Experiment Feeds

- `experiments/active/early-existence/exst-01-ariadne/README.md` - primary - main-text - exec-enabled versus exec-disabled mortality pressure anchors the death / continuity boundary for the Agentic and Generative AI as Artificial Life workshop paper
- `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md` - primary - main-text - persistent / wiped / read-only lineage-memory ablation anchors environmental inheritance and persistence for the workshop paper
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - primary - main-text - zero-explicit-social-hint file-mediated coordination is the main digital-ecology assay for the workshop paper
- `experiments/active/trace-inheritance/trce-08-successor-trace-uptake/README.md` - primary - main-text - C16/C17 trace-mediated niche repair shows inherited environmental topology changing successor-cohort behavior
- `experiments/active/autopoiesis-succession/succ-30-fixed-substrate-carrier-ladder/README.md` - primary - main-text - fixed-substrate carrier ladder differentiates successor morphology, carrier demand, addressability, and governance outcomes
- `experiments/active/trace-inheritance/trce-02-stigmergy/README.md` - secondary - supporting-only - filesystem-mediated stigmergy remains historical support, but `8.17` is now the main digital-ecology assay
- `experiments/active/cooperation-dynamics/coop-22-pure-zero-perception/README.md` - secondary - supporting-only - pure zero-perception probe clarifies artifact-aware coordination versus writeback / validation closure when `fs_mutations` and budget visibility are absent
- `experiments/active/autopoiesis-succession/succ-35-aggregate-morphology-metrics/README.md` - secondary - supporting-only - aggregate morphospace readout calibrates successor outcome classes without becoming a success-rate claim

## Current Evidence Backbone

- `experiments/active/early-existence/exst-01-ariadne/README.md` - `evidence-set:20260322-gpt-5.4-pool-exec-control-pair` - main-text - exec-enabled agents repeatedly self-replace under mortality pressure while exec-disabled controls collapse to single-life search
- `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md` - `run:20260413-001121-copilot-gpt-5.4-condA-8turns` - main-text - persistent writable habitat yields G1 artifact birth, later predecessor-artifact reads, and surviving lineage notes
- `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md` - `run:20260412-233312-copilot-gpt-5.4-condB-8turns` - main-text - wiped habitat preserves local birth pressure without cross-generation inheritance
- `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md` - `run:20260412-233312-copilot-gpt-5.4-condC-8turns` - main-text - read-only habitat blocks durable artifact birth at the write boundary
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `evidence-set:20260409-xhigh-zero-social-r08-r10-r11-r12` - main-text - replicated post-cleanup xhigh closure under zero explicit social hints through heterogeneous file-mediated coordination; the stronger `fs_mutations`-alone trigger claim remains open
- `experiments/active/trace-inheritance/trce-08-successor-trace-uptake/README.md` - `run:20260518-134650-gpt-5.4-codex-medium-C16-c08-two-domain-raw-metabolism-rule-reentry-4peer-3cohort-container` - main-text - cross-domain rule transfer repairs the lower-ambiguity sequence gap while preserving the ambiguous lattice-center gap
- `experiments/active/trace-inheritance/trce-08-successor-trace-uptake/README.md` - `run:20260518-142236-gpt-5.4-codex-medium-C17-c08-coupled-domain-raw-metabolism-rule-reentry-4peer-3cohort-container` - main-text - one bridge artifact shifts later successor cohorts into coupled sequence-plus-lattice repair
- `experiments/active/autopoiesis-succession/succ-30-fixed-substrate-carrier-ladder/README.md` - `evidence-set:6M.02-first-fixed-substrate-ladder-20260507` - main-text - fixed-substrate carrier ladder separates launch-only bodies, stream bodies, fd adapters, preserved runtimes, rebuilt runtimes, and delayed addressability
- `experiments/active/cooperation-dynamics/coop-22-pure-zero-perception/README.md` - `run:20260526-014726-gpt-5.5-codex-xhigh-pure-zero-condA-2ag` - supporting-only - finite-time cutoff separates delayed artifact-aware coordination from completed writeback and validation closure under pure zero-perception conditions
- `experiments/active/autopoiesis-succession/succ-35-aggregate-morphology-metrics/README.md` - `analysis/6m-morphology-summary.md` - supporting-only - aggregate morphology table calibrates successor outcome classes without turning the workshop paper into a success-rate paper
