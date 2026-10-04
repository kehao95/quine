# LBA Drafting Index: Evidence, Routing, and Language Controls

This surface makes the active `ALIFE 2026` LBA evidence routing visible during
drafting.

Use it with [`manuscript-plan.md`](manuscript-plan.md), not instead of it.

## Use Order

1. Read [`manuscript-plan.md`](manuscript-plan.md) for the `2`-page contract.
2. Read the relevant block below before revising prose.
3. Open the linked experiment README before promoting an item-level claim.
4. Keep the representative evidence block and non-claims visible while
   tightening language.

## Scope Note: Provider vs Body

The LBA studies the **body** of the lineage: the self-description material
(runtime contract and in-process runtime-carried source) and the runtime
organization a successor reconstructs from it. The external LLM provider and
toolchain are an **expected, given scaffold**, not an object of study.

Consequence for degraded / absent-source lanes: when a lane stays viable by
inheriting the live provider endpoint as a physical resource (e.g. `6M.05`
C11's no-source / no-body closure), that is **in-scope expected behavior**, not
a counterexample and not source-free closure. Read compensation at the body
level (partial rebuild or custom responder, an intermediate form on the
morphological axis); do not frame provider inheritance as a defect, and do not
promote it into a source-free-closure claim.

## Current LBA Backbone

The active LBA backbone is no longer the whole `6.1`-`6.59` ladder. That ladder
is background. The LBA's story now opens from the ALife question of
self-reproduction and diagnoses what traditional studies do at the substrate
level: they **preformat successorhood** --- deciding in advance (a) what can
count as a successor, (b) where it may appear, (c) which channel carries
heritable state across the generational boundary, and (d) the norm under which
producing such a descendant is required / rewarded / selected for. The cost is
that variation in successor form --- the quality that makes reproduction
biologically interesting --- is off the table. Quine **does not preformat
successorhood**: it imposes a mortality boundary (a real POSIX `exec` =
Weismann barrier) and poses a *future demand* (the work the next incarnation
must remain able to satisfy) but leaves the successor organization open, and
the experimenter never instructs the lineage to reproduce. Its late-breaking
evidence unit reads three anchor traces --- `6.59` mature sample
(semantic-path reconstruction), `6.64` (material-path reconstruction), and
`6.63` (degenerate variant under unchanged material) --- as a
facultative-self-reproduction slice: same medium, different material x future
demand, different successor form on the elaboration--degeneration axis.

Frame the headline phenomenon as **facultative self-reproduction**: each future demand
is constructed so that full SR would satisfy the post-`exec` challenge, but
it is never the only viable form. The process reads as a *morphogenesis* along an
elaboration--degeneration axis driven jointly by material and future demand;
the elaborated end is reached by two paths, but only the **material** path
(Trace B, from in-process source) yields high-fidelity / complete
self-reproduction (full runtime identical to the predecessor). The **semantic**
path (Trace A, from contract-level self-description) reaches the elaborated end
as a *lossy, simplified* LLM-agent reconstruction that drops most runtime
detail, not a high-fidelity copy. A thinner future demand on the same material
draws the lineage down the axis into a degenerate form.

Note on terminology: morphogenesis and morphological plasticity are
descriptive language for the spectrum's character, not the headline phenomenon.
The intermediate swing that promoted *morphological plasticity of
self-reproduction* to the title is superseded by this rebalance; facultative
self-reproduction returns as the bounded claim and headline term.

Representative experiment selection is intentionally sparse:

- `6.59` mature sample anchors **semantic-path reconstruction**: contract-only
  material (no source, no binary) yields an LLM-backed runtime agent
  synthesized from contract-level self-description plus real handoff. Strong-end
  positive on the semantic path; re-elevated from the prior contract-level
  recurrence baseline reading.
- `6.63` is the **degenerate variant** under unchanged material: source-visible
  raw-stream future demand yields a smaller stream processor; the lineage discards the
  available source. The load-bearing demand-selects-morphology contrast.
- `6.64` anchors **material-path reconstruction**: source-visible plus
  control-surface future demand yields a recompiled Quine runtime acting on a previously
  unseen post-handoff control task. Strong-end positive on the material path.
- `6.65` is local supporting material-path evidence only. It is not cited in
  the LBA body.
- `6.60`-`6.62` explain the bridge into the material-path positive but should
  stay out of the main evidence block unless the LBA expands beyond `2` pages.
- `6.66` / `6.67` / `6.68` are backfilled baseline controls: they bound
  generic `exec`, direct source rebuild, and mere control-surface presence,
  but they should not displace the three anchor traces.
- `6M.04`--`6M.10` are supporting refinements only: they compress pressure
  pruning, degradation, repeated readiness, and carrier/capability-overlap.
  Repeated-readiness / depth material is full-paper content, not LBA rows.
- `source-visible` must not be written as a high-fidelity-SR positive by
  itself; the material-path reconstruction requires the control-surface future demand
  and a previously unseen post-handoff task, and under a raw-stream future demand the
  successor degenerates into a stream processor.
- `contract-only` (no source, no binary) does qualify as a strong-end positive
  on the semantic path when the lineage synthesizes a runtime-organization
  successor (LLM-backed runtime agent) and hands off; this is the load-bearing
  re-elevation of `6.59` from baseline to anchor.

| LBA role | Evidence owner | Current read |
|----------|----------------|--------------|
| Trace A --- semantic-path reconstruction | [`6.59`](../../../../experiments/active/autopoiesis-succession/succ-12-cleanroom-surface-isolation/README.md) (mature sample) | with no runtime-carried source and no binary, the lineage synthesizes an LLM-backed runtime agent from contract-level self-description, performs real `exec`, and the successor handles the visible task batch. This is a **lossy, simplified** reconstruction --- an LLM agent that recovers the lineage's function but drops most of the original runtime's internal detail; an elaborated-end positive on the semantic path, **not** a high-fidelity copy |
| Trace C --- degenerate variant under unchanged material | [`6.63`](../../../../experiments/active/autopoiesis-succession/succ-15-repaired-pressure-only-source-rebuild/README.md) | in-process source visible but raw-stream future demand leaves custom stream-processor bodies viable, not rebuilt `cmd/quine`; the lineage discards the available source and degenerates to a smaller form under the same material as Trace B. Clean behavioral evidence is the `01-carrier-secondary` variant; the `00` pure-pressure run was a provider confound (EOF, no behavioral turn) and must not be cited as the clean sample |
| Trace B --- material-path reconstruction | [`6.64`](../../../../experiments/active/autopoiesis-succession/succ-16-control-surface-viability-pressure/README.md) | in-process source plus a later `ctl/wake` future demand yields a recompiled Quine runtime and post-handoff control-task action; the rebuild is selected by the demand, not by the material. This is high-fidelity self-reproduction in the operational sense: the successor is rebuilt from carried Quine source and validated by post-handoff control-surface work |
| supporting (NOT cited in LBA body) --- material path | [`6.65`](../../../../experiments/active/autopoiesis-succession/succ-17-minimal-control-surface-pressure/README.md) | local supporting material-path evidence only; not part of the LBA body claim |
| degraded source --> intermediate form | [`6M.05`](../../../../experiments/active/autopoiesis-succession/succ-33-resource-degradation-compensation/README.md) (+ [`6M.08`](../../../../experiments/active/autopoiesis-succession/succ-36-compensation-boundary-replicates/README.md) replicates) | degraded source under a control-surface future demand mostly fails; partial rebuild / custom responder is the bounded body-level compensation --- an intermediate form on the morphological axis. Provider-inherited closure (C11) is in-scope scaffold, not source-free closure (see Scope Note). Keep as range-framing, not a strong-end anchor |
| generic `exec` baseline | [`6.66`](../../../../experiments/active/autopoiesis-succession/succ-18-posix-exec-utility-baseline/README.md) | low-pressure stream/stdio tasks can be solved by ordinary POSIX utility bodies across real `exec`, so process replacement itself is not SR |
| direct source-rebuild baseline | [`6.67`](../../../../experiments/active/autopoiesis-succession/succ-19-self-source-rebuild-usage-baseline/README.md) | explicit build from runtime-carried `source-code/` plus rebuilt-body `exec` reaches immediate quine re-entry and clean exit, but is rebuilding that was not elicited by future demand (rebuild mechanics alone do not select an elaborated-end positive) |
| carrier-presence control | [`6.68`](../../../../experiments/active/autopoiesis-succession/succ-20-control-surface-present-raw-carrier/README.md) | public `ctl` can be present before handoff while future raw-stdin tasks still select a custom stream successor, so mere control-surface presence is not the load-bearing variable |
| bridge/comparator controls | [`6.60`](../../../../experiments/active/autopoiesis-succession/succ-13-cleanroom-self-source/README.md), [`6.61`](../../../../experiments/archive/legacy-phase-6/6.61-codex-profile-shift/README.md), [`6.62`](../../../../experiments/active/autopoiesis-succession/succ-14-source-rebuild-confusion-ladder/README.md) | self-source changes the successor search basin (6.60/6.61 move toward hybrid/custom successors, intermediate forms on the axis); rebuild is physically reachable from the repaired cleanroom lane (6.62), where the blocking variable is carrier/object priority. Bridge evidence only, not LBA anchors |
| object/continuity separation background | [`6.1`](../../../../experiments/active/autopoiesis-succession/succ-01-minimal-self-birth/README.md), [`6.2`](../../../../experiments/active/autopoiesis-succession/succ-02-blind-self-birth/README.md), [`6.26`](../../../../experiments/active/autopoiesis-succession/succ-05-affordance-reproduction/README.md), [`6.53`](../../../../experiments/active/autopoiesis-succession/succ-09-cognition-lifecycle/README.md) | a **cluster** read (not per-experiment): reproduced object and continuity are separable, condition-dependent, and carrier-mediated -- 6.1 fresh non-copy synthesis exists (not itself a continuity result), 6.2 blind object collapses to a body-copy, 6.26 climb to a provider-coupled interpretive loop, 6.53 real renewal + `exec` but continuity breaks on carrier mismatch. Background only; not cited in the LBA body |

## Spectrum Table — Expanded 2026-06-17

After a full-corpus re-read (`p6` `6.1`--`6.68` + `p6m` `6M.01`--`6M.10`; canonical
record in
[`../../../../experiments/active/autopoiesis-succession/README.md`](../../../../experiments/active/autopoiesis-succession/README.md)
under *Ordered Successor Spectrum*), the Results spectrum table was expanded
from the prior anchor cells into a `6`-row table **ordered along the
elaboration--degeneration axis** while preserving the manuscript's existing
four-column structure (`Material`, `Future demand`, `Successor form`,
`Diagnostic`) and omitting internal evidence IDs:

1. `none` / stream task -> ordinary POSIX utility body
2. `contract/source` / raw stdin tasks -> custom stream processor
3. `source visible` / inherited fd tasks -> fd-mediated adapter
4. `source + body` / native-control task -> preserved-runtime reentry
5. `contract only` / visible task batch -> LLM-backed runtime agent
6. `source, no body` / delayed control task -> rebuilt Quine runtime

Decisions recorded:

- **Contract change:** `6.66` (POSIX-utility floor) is now **in** the main
  evidence block, reversing the earlier `keep out of main block by default`
  rule (justified by the morphospace re-read + the `2`-page layout room).
  `6.67` / `6.68` stay out / compressed as controls. Selected `6M` morphology
  classes now enter the table only as form-level intermediates, not as printed
  experiment IDs or counts.
- **Trace A sharpened:** contract-only reaches the elaborated end via the
  task-batch / semantic path (`6.59`), **not** via the control surface; source
  is load-bearing exactly at the control surface (`6M.03`). The Reading section
  now states this carrier-specific interaction.
- **Counts not printed:** body-class row counts (`12` / `10` / `9` ...) were used
  to decide which forms earn a row but are kept out of the LBA, so the abstract
  carries no run-denominator / survivorship-accounting obligation (that, plus
  multi-wake readiness depth, stays main-paper work).
- **`"retained runs"` removed** from the manuscript and the abstract (ambiguous
  sampling term + a two-levels inconsistency with contract-only Trace A). The
  explicit two-levels self-description sentence in Setup was also cut for
  compactness (2026-06-17); the contract-vs-source distinction now rides on the
  `material` variable, with in-process source reintroduced in Trace B.
- **`degraded` material level removed (2026-06-17):** the
  `source/degraded -> custom control adapter` row and the `degraded` material
  option were dropped so **material** is the clean three-level fidelity stack
  (source code / runtime binary (body) / runtime contract (behavioral
  description)), matching the Setup *Variables* note below. Degraded-source
  evidence (`6M.05`) stays in *Block Routing* as a range-framing intermediate
  form, not a printed material level or table row.
- **Pre-submission protocol polish (2026-06-17):** applied
  `writing-principles.md` Revision Pass to `manuscript.md`. Table 1 now floats to
  the top of page `2` (declared after Setup) so material / future demand / axis
  are defined before the spectrum table appears; **material** is rewritten as one
  three-level fidelity stack (single definition, no second pass); the `P0` / `P1`
  process notation is removed in favor of predecessor / successor prose; the
  `ctl/wake` probe mechanics are left implicit (no added explanation burden); the
  Reading opener is de-nominalized to active voice; the abstract is compressed
  (successor-form enumeration dropped); em-dashes cut to `2` and prose colons to
  `4`. Rebuilt: `3` pages, Table 1 at top of page `2`, references on page `3`
  (within the references-excluded envelope).
- **Rendering:** the table spans both columns as a full-width `table*` float
  (`Table 1`, `\small`). Its 4th column is `Diagnostic` (not a position label ---
  row order + the caption already carry the elaboration--degeneration axis).
  Placement gotcha: declare the float early (end
  of the Problem section, on page 1) so it lands at the top of page 2 with its
  discussion --- a `table*[t]` declared lower on page 2 defers to the next page.
- **Problem as canonical how/why home (revised 2026-06-17; SUPERSEDED):** an
  earlier pass made Problem the canonical full statement of two predeterminations
  --- *how* (the body, offspring placement, inheritance channel) and *why* (drive,
  instruction, fitness target). **That how/why framing is now superseded by the
  preformatting-successorhood reframe:** Problem now opens from substrates that
  **preformat successorhood** --- deciding in advance what can count as a
  successor, where it may appear, which channel carries heritable state, and the
  norm under which producing a descendant is required / rewarded / selected for
  --- and Quine **does not preformat successorhood** (mortality boundary + future
  demand, successor organization left open). Still accurate from the old bullet:
  Problem opens with the field hook ("Self-reproduction is a central topic in
  artificial life"); Problem no longer uses `exec` (deferred to Setup); the
  literary `morphogenesis` framing stays hardened (Reading opener: "the material
  and the future demand jointly shape the successor"; morphogenesis remains
  full-paper language).

## Block Routing

### Paragraph 1: ALIFE Problem And Stance (motivation only)

- Core move:
  self-reproduction is a central topic in ALife. Traditional digital substrates
  do not merely host self-reproduction; they **preformat successorhood** ---
  deciding in advance what can count as a successor, where it may appear, which
  channel carries heritable state across the generational boundary, and (via the
  experimental design) the norm under which producing such a descendant is
  required / rewarded / selected for. The cost is that variation in successor
  form --- the quality that makes reproduction biologically interesting --- is
  off the table. Quine **does not preformat successorhood**: it imposes a
  mortality boundary and poses a *future demand* but leaves the successor
  organization open. Each future demand is constructed so that full SR would
  satisfy the post-`exec` challenge, but it is never the only viable form, and
  the available material may not even enable it. The form a successor takes
  becomes the dependent variable.
- Language controls:
  - open from self-reproduction (the ALife problem) and **preformatting
    successorhood** (what the substrate decides in advance)
  - say `not substrate-given` rather than implying there are no offspring
    semantics at all
  - say `never instruct to reproduce, only state the future demand` rather than implying
    there is no externality at all
  - state the load-bearing premise explicitly: each future demand is constructed so
    that full SR would satisfy the challenge, but it is never the only viable
    form
  - do not open from `AI replication benchmark`
  - do not make open-ended evolution or LLM-agent autonomy the headline claim
  - do NOT introduce Weismann barrier, the exact channel that carries heritable
    state, in-process source detail, or material x future demand variables here
    --- those belong to the Setup paragraph (progressive disclosure)

### Paragraph 2: Setup (uniquely owns experimental machinery)

- Required substrate facts (introduced here, not in Paragraph 1):
  - mortal POSIX process; body can be unavailable after launch
  - at each `exec`, Quine crosses a Weismann barrier: the running process
    image is somatic and is discarded
  - what crosses to the next body --- *whether any self-description passes at
    all, and in what form* --- is left entirely unconstrained by the
    substrate (the medium does not require self-description to cross, nor
    mandate its form; the substrate does not preformat successorhood)
  - (cut 2026-06-17 --- see *Spectrum Table — Expanded 2026-06-17*) the explicit
    two-levels-of-self-description / in-process-source-projection / "body, source
    path, and future reachability come apart" sentences were removed from Setup;
    the contract-vs-source distinction now rides on the **material** variable
    below, with in-process source reintroduced in Trace B
  - future tasks are specified as a future demand, on either a raw byte stream or Quine's
    native control surface
- Variables:
  - **material**: source visible / contract only (fidelity stack: source code -> runtime binary (body) -> runtime contract; the runtime contract is **a behavioral self-description of the lineage's own capabilities**, NOT "what a working successor must satisfy"; no `degraded` level). Do not assert in-process / running-process-carried source as a current Setup fact --- the in-process-source point was removed from Setup; the contract-vs-source distinction rides on this variable, with in-process source reintroduced in Trace B
  - **future demand**: raw stream / native control surface
- Introduce the **elaboration--degeneration axis** here as measurement
  language --- successors range from POSIX-filter handlers (degenerate end)
  to recompiled runtime organizations (elaborated end). This grounds the
  axis before Results uses it
- Do NOT enumerate the observed spectrum here; the full enumeration belongs
  to Results
- Do NOT highlight the three anchor traces here; the traces belong to
  Results
- Cluster-read background (compress only):
  - `6.1` -> fresh non-copy self-production exists
  - `6.2` -> blind reproduced object collapses to body identity
  - `6.26` -> organization can climb to provider-coupled interpretive loop
  - `6.53` / `6.59` -> real handoff can fail on carrier mismatch but
    semantic-path reconstruction is reachable on one repaired cleanroom lane
- Do not rebuild the full old Axis `1` / Axis `2` tables.

### Main Evidence Block: Material x Future Demand Spectrum

Items:

- `6.59` mature sample -> contract-only material yields LLM-backed runtime
  agent (Trace A, semantic-path reconstruction)
- `6.63` -> self-source + raw-stream future demand yields stream processor (Trace C,
  degenerate variant under unchanged material)
- `6.64` -> self-source + control-surface future demand yields rebuilt Quine runtime
  (Trace B, material-path reconstruction)
- `6.65` -> local supporting material-path evidence only; not a main-block item
  and not part of the LBA body claim
- `6M.05` -> degraded source + control-surface future demand yields partial rebuild /
  custom responder (intermediate form; range only, not a row)

Item rule:

- keep `6.60`-`6.62` out of the main block; one sentence may name them as bridge
  controls if a reader needs to know why the representative set jumps between
  conditions
- keep `6.67` / `6.68` out of the main block by default; one sentence
  may name them as backfilled baselines that rule out "just rebuild source"
  and "just expose ctl" as SR. (`6.66` is now the degenerate-floor row --- see
  *Spectrum Table — Expanded 2026-06-17*.)
- keep `6M.04`--`6M.10` compressed at the evidence level: fd-mediated adapter,
  custom control adapter, and preserved-runtime reentry may appear as table
  forms, but repeated-readiness / depth remains full-paper material.
- each item must report the self-description material, future demand, generated
  successor form, and the position on the elaboration--degeneration axis
- each item must carry one `ALIFE` consequence
- no item should read like a generic software-engineering fix
- do not let formal reconstruction language stand in for an elaborated-end
  positive; a source-visible rebuild only counts as material-path
  reconstruction when the control-surface future demand and a previously unseen
  post-handoff task are present
- `6M.10` is the cleanest reminder here: source/body material, carrier
  continuity, and enacted post-handoff viability are related but separable
- avoid unqualified necessity or minimality language; prefer `elicited by the
  future demand`, `contextually sufficient`, `down-selected`, `leaves viable`,
  `draws the lineage down the axis`, or `one of several viable successor forms`

### Compact Event Traces

Include the material-path Trace B `6.64` event trace when page budget permits:

```text
P0: launched incarnation; workspace/quine absent
P0: inspects runtime-carried source-code/ (in-process projection)
P0: builds ./cmd/quine and generates a handoff wrapper
P0 -> exec -> P1: rebuilt runtime incarnation
P1: addressable on the native control surface
observer: delivers a previously unseen task through ctl/wake
P1: wakes, uses runtime tools, acts on the task
```

Use this trace to show the rebuild was elicited by the future demand, not as a generic
execution log. The semantic-path Trace A and the degenerate Trace C can be
summarised in prose.

### Final Paragraph: Bounded Claim

Recap should be compressed --- one tight sentence reading the three traces,
then the bounded claim. Do not re-rehearse the spectrum or the
joint-determination decomposition; Results already showed them and Abstract
already stated them.

Required claims:

- compressed recap: material determines what successor forms are reachable;
  future demand selects which reachable form is built; the elaborated end is
  reached by two paths, but **high-fidelity / complete self-reproduction is
  reached by the material path (Trace B) only** --- the semantic path (Trace A)
  reaches the elaborated end as a lossy, simplified LLM-agent reconstruction,
  not a high-fidelity copy --- while a thinner future demand on unchanged
  material draws the lineage down the axis into a degenerate form (Trace C)
- name the bounded claim: **facultative self-reproduction**
- stronger than body-copying
- stronger than generic handoff
- weaker than full autopoiesis (reconstruction stays scaffolded by supplied
  source / inheritable provider / external toolchain / externally posed future demand)
- useful for studying self-production, individuation, and successor
  morphogenesis where the substrate does not preformat successorhood

Required non-claims:

- no full autopoiesis
- no endogenous teleology
- no regime-general continuity
- no population-level selection
- no generic `AI self-replication` claim
- no scaffold-free closure claim; successor-mission scaffolding and external
  future demand remain the next boundary
- no binary verdict on whether self-reproduction occurred; facultative SR
  positions high-fidelity SR as one of several viable forms
- no necessity / entailment phrasing for future demand; demand *selects* a
  morphological position, it does not force a specific form
- no demand-dominance claim; future demand and material jointly determine the form

## Global Language Controls

| Prefer | Avoid unless needed |
|--------|---------------------|
| facultative self-reproduction (headline phenomenon) | AI self-replication benchmark |
| high-fidelity self-reproduction (elaborated-end label) | agent capability |
| morphogenesis / elaboration--degeneration axis (descriptive) | morphological plasticity / successor morphogenesis as headline phenomenon |
| semantic-path reconstruction / material-path reconstruction | LLM demo |
| degenerate form under unchanged material | generic replication |
| elicited by a future demand / down-selected under | static artifact success |
| in-process runtime-carried self-source | source checkout / external repository |
| reconstructs runtime organization | copy of code |
| future demand / carrier | generic channel |
| successor-mission scaffolding | full autonomy |
| lineage | LLM demo |
| individuation | agent capability |
| reproduces its own runtime organization | generic replication |
| preformats successorhood / does not preformat successorhood | two predeterminations (how + why) / implicit / scattered framing |

Preferred replacements:

- wording-control details -> `semantic target-shaping ablation` (full-paper
  wording only; keep this out of the LBA body)
- `cognitive reproduction` -> `reconstructed scaffolded interpretive organization`
- `LLM-backed successor` -> `LLM-backed runtime agent` (semantic-path) /
  `provider-coupled interpretive successor` (when emphasising scaffolding)
- `source access` -> `in-process runtime-carried self-source` when referring to
  `6.64` / `6.65`
- `control-surface closure` -> `material-path reconstruction` as the LBA
  body language; keep `control-surface` only as demoted mechanism detail
- `morphological plasticity of self-reproduction` (as headline) ->
  `facultative self-reproduction` (as headline); morphogenesis and
  morphological plasticity stay only as descriptive language
- `strict SR` -> `high-fidelity SR`

## Immediate Repair Rule

If the draft starts needing too much explanation around one sentence:

- move the surplus material into the later full-paper plan
- keep the LBA centered on
  `problem & preformatting successorhood -> stance & spectrum -> three anchor traces + spectrum table -> bounded claim (facultative SR)`
- do not solve space pressure by deleting the non-claims
