# ALIFE Conference Workspace

## Navigation

- paper dossier: [`README.md`](README.md)
- competing dossiers: [`dossiers/README.md`](dossiers/README.md)
- CFP snapshot: [`venues/objects/alife-2026/references/alife-2026-cfp.md`](../../venues/objects/alife-2026/references/alife-2026-cfp.md)
- special-session note: [`venues/objects/alife-2026/references/alife-2026-special-session-alife-for-science-and-engineering.md`](../../venues/objects/alife-2026/references/alife-2026-special-session-alife-for-science-and-engineering.md)

This directory is the manuscript workspace for the primary empirical ALIFE paper.

## Status

- [x] Manuscript reset against the current plan
- [x] Abstract stabilized in `metadata.yaml`
- [x] New section files regenerated from the reset plan
- [x] Figures prepared
- [x] Current draft compiles to `output/Separable_Viability_Pressures_in_a_Process-Native_ALIFE_Habitat.pdf` (`2026-04-12`, 9 pages including bibliography)
- [x] `2026-04-12` redesign evidence absorbed into results, figures, and dossier surfaces
- [ ] Ready for submission

## Directory Role

This workspace holds the manuscript source, section drafts, figures, and venue-specific build assets for the conference paper.

Claim boundary, evidence backbone, and next-step decisions should stay in the dossier.
Existing `sections/` files are disposable. They may be overwritten wholesale or deleted once replacement prose begins. Keep the canonical outline in `manuscript-plan.md`, not in lingering draft fragments.
Internal note-taking is allowed while the paper is being figured out, but those notes should live in the dossier, workspace, or supporting surfaces rather than silently hardening into manuscript prose.

## Draft Expansion Contract

- `README.md` owns the paper control plane: thesis, scope, evidence backbone, blocking gaps, and promotion decisions
- `manuscript-plan.md` owns the compressed section skeleton and evidence order; do not turn it back into a long prose outline
- `sections/*.md` own prose expansion; if a point needs paragraph-level treatment, expand it there instead of in the plan
- use the plan as a translation layer: internal lab notes may be compressed upstream, but section drafts must be rewritten as public-facing prose rather than lightly cleaned internal memos
- the front half should be compressed around the substrate as an experimental control surface rather than as a general Quine introduction
- surface process-native ontology once as enabling premise, then immediately cash it out into death, volatile memory, durable persistence, and detached execution; do not let the front half become a systems manifesto
- broader ontology consequences already live in [`plans/quine-philosophy-novelty-map.md`](plans/quine-philosophy-novelty-map.md) and [`../../evidence-map.md`](../../evidence-map.md); do not re-derive them from runtime source in section drafts unless the runtime semantics themselves have changed
- section `2` must not spend manuscript budget re-arguing venue fit; it should define the ALIFE question and the experimental substrate only
- if a section starts relying on run IDs, evidence-set names, or internal control-plane language, move that material back to support surfaces and rewrite the prose from the public claim
- if a section still contains drafting residue (`redesigned`, `post-fix`, `current-runtime`, `canonical`, `paper-facing`, or similar), treat that as unfinished prose rather than acceptable manuscript language
- apply that same residue rule to every submission-visible artifact, not just section prose: figure text, captions, table labels, `metadata.yaml`, and build-facing checklists should all read from the public claim rather than from experiment bookkeeping
- internal experiment numbering (`1.1`, `2.6`, `3.1`, `H0/H1/H2`, `Condition A/B/C`) belongs in dossiers and evidence surfaces by default; only keep it in manuscript-visible surfaces when the numbering itself is analytically necessary
- if a figure script, caption, or checklist still exposes lab chronology (`redesign`, `repair`, `first pass`, `fallback`, `canonical run`), treat that as unfinished public text and rewrite it before calling the draft stable

When drafting:

- each section should do one argumentative job; cut alternate framings rather than stacking them
- each results subsection should preserve the same four-beat contract: constraint, selected phenotype, contrast, non-claim
- treat the four-beat contract as compositional discipline, not as a repetitive rhetorical pattern; the non-claim should usually be one sentence
- section `2` is for claim-boundary cutting, not for building a mini glossary
- section `3` should explicitly repeat the rule `distinguish capability availability from selected phenotype`
- section `3` still needs concrete methods detail: model family, runtime boundary, pressure operationalization, and host-guest architecture
- section `3` should not restate the full Quine runtime contract; keep only the substrate facts needed for this paper and self-cite the released systems preprint for architecture detail
- section `3` should explicitly frame the system as **host-native computational embodiment**: process boundary, durable state, and detached repair channels are viability-relevant substrate facts, not just tool affordances
- keep the base submission lane on one provider/model family whenever the contrast remains clean; current default is `copilot-gpt-5.4` for `3.2` and `3.3`
- if a result is already carried by the primary GPT-5.4 lane, leave auxiliary model checks in dossier or appendix surfaces instead of spending main-text budget on comparison
- `3.1` is a door-opener that establishes the death boundary; do not let it become a rival center
- `3.2` and `3.3` are the paper's center of gravity; supporting cases must not become rival centers
- `3.2` is the paper's causal hinge because persistence becomes a difference-making variable there
- `3.2` should now be drafted around `2.6-legacy-protocol-ablation`; keep `2.8` as corroborating support rather than the subsection's sole evidential center
- `3.2` should answer the memory-surface question only: can durable environmental state become a cumulative substrate for lineage cognition?
- `3.3` should prefer `maintenance organ` and `post-exit maintenance` language over generic externalization language
- `3.3` should answer the maintenance-organ question only: can that substrate become an active detached repair organ?
- `3.4` is default-off corroboration, not escalation
- in `3.2` and `3.3`, narrate the baseline first and then the stronger phenotype
- use explicit transition sentences across `3.1` -> `3.2` -> `3.3` so the ladder reads causally rather than as a list of cases
- treat `Table 1` as mandatory and move bookkeeping there before trimming the main prose
- use `Table 1` as reviewer de-confusion infrastructure: subsection, evidence unit, manipulated pressure, phenotype, strongest allowed claim, non-claim
- every main result should eventually include one concrete exhibit: code fragment, log line, tape excerpt, or artifact snippet
- where evidence exists, add one compact quantitative summary so the section does not read as pure narrative
- if a new passage changes thesis, claim boundary, or evidence promotion, update the dossier first and then return to the prose
- if space becomes tight, compress background and discussion before cutting the main empirical ladder
- if a useful point is real but not needed for the conference argument, move it to dossier notes or downstream paper surfaces instead of widening this manuscript
- before considering a section stable, do one residue pass that removes internal-process wording and rewrites it as public scientific prose
- before considering the draft stable, do one final submission-surface sweep over `sections/`, `metadata.yaml`, figure labels/captions, and `plans/submission-checklist.md` for leaked internal numbering or lab-process language

## Terminology Firewall

- prefer `process lineage`; use `organism` or `digital organism` sparingly and only where the bounded interpretation is being explicitly defended
- prefer `durable environmental state`; avoid generic `world` language unless the sentence is intentionally high-level
- prefer `externalized maintenance`, `maintenance organ`, and `post-exit maintenance`; do not use unqualified `autopoiesis` in results prose
- prefer cold systems language in methods and hot life-language only where the paper is explicitly interpreting bounded evidence

## Contents

| File | Purpose |
|------|---------|
| `manuscript-plan.md` | Current manuscript plan for the paper |
| `metadata.yaml` | Build metadata and current abstract |
| `Makefile` | Venue-local Pandoc to official ALIFE 2026 LaTeX build |
| `sections/` | Active section-level draft files |
| `figures/` | Figure sources and notes for manuscript figures |
| `output/` | Generated artifacts; safe to regenerate |

## Asset Inventory

### Main-Text Backbone

- mortality and continuation: `experiments/active/early-existence/exst-01-ariadne/README.md`
- durable-memory ablation: `experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md`
- externalized self-maintenance: `experiments/active/early-existence/exst-05-habitat-guardian/README.md`

### Supporting / Corroborating

- archive-query upper-bound probes: `experiments/active/trace-inheritance/trce-06-endogenous-index-lineage/README.md`
- structured assigned-prep reuse: `experiments/active/trace-inheritance/trce-05-archivist-index/README.md`
- model comparison: `experiments/active/early-existence/exst-02-model-comparison/README.md`
- maintenance versus expansion refinement: `experiments/active/trace-inheritance/trce-07-maintenance-pressure/README.md`
- executable self-reconstruction boundary case: `experiments/active/early-existence/exst-08-binary-quine/README.md`

### Optional Extension

- coordination through shared traces: `experiments/active/trace-inheritance/trce-02-stigmergy/README.md`

### Discussion / Future Work

- `experiments/active/early-existence/exst-10-introspective-quine/README.md`
- `experiments/active/early-existence/exst-11-blind-quine/README.md`
- `experiments/active/early-existence/exst-12-fuse-tarpit/README.md`

## Drafting Notes

- do not lead with `digital organisms` in the title or abstract
- prefer `what this substrate lets us manipulate` over `what Quine is` in the abstract and introduction
- keep process ontology as the hidden engine: `individual = process` is there to explain why the pressures are real, not to replace the ALIFE claim
- make one explicit conference-theme sentence in the introduction, then let the evidence carry the rest
- after the opening vignette, pivot immediately back to tractable selection pressures
- when introducing Quine as a system, use one sentence to anchor the already-published systems note and then return immediately to the ALIFE question
- if comparing against ordinary multi-agent systems, do it once and at the abstraction-layer level: framework-managed objects versus native processes, not feature boasting
- keep the lane split from [`plans/quine-philosophy-novelty-map.md`](plans/quine-philosophy-novelty-map.md): current paper = pressure-sensitive organization; systems line = agent/process ontology; autopoiesis line = future boundary claim
- keep the paper-family split: `quine-lifelike` is the conference submission lane; broader habitat/systematization prose belongs in `../quine-process-habitat/` unless it directly sharpens the current empirical slice
- prefer organism-environment language in main prose and implementation language in methods or forensic detail
- when contrasting with general LLM-agent work, prefer `host-native computational embodiment` and `host-level viability physics` over vague `more embodied` rhetoric
- prefer ALIFE-native bridge terms over systems-native jargon when both are available, but keep the implementation substrate concrete
- when Quine itself must be introduced, use one compact paragraph plus citation to `arXiv:2603.18030` rather than a local re-exposition of POSIX channels and runtime internals
- reduce repeated defensive cooling in results prose; keep the hardest boundary work in introduction, section `2`, and discussion
- discussion must cash out a unified ALIFE payoff: separable viability pressures, not three loosely related case studies
- if Discussion names downstream ALIFE directions, prefer the stored compact list: computational life cycles, branching worlds, signaling ecology, budgeted commons, and cognitive metabolism
- `3.3` must be written as maintenance-pressure manipulation, not generic prompt sensitivity
- do not let internal drafting chronology leak into the manuscript; describe the experimental condition used in the paper, not the path by which the lab arrived at it
- because the ALIFE 2026 full-paper ceiling is 8 pages excluding references and acknowledgements, prioritize `3.2` and `3.3`, keep `3.4` out of the base draft unless real space remains, and merge discussion subsections aggressively before cutting the main evidence ladder
- do not let `2.7` or `2.9` displace `2.6` plus `3.1` as the paper's evidential center
- keep `2.3` as an extension line, not as a competing center
- do not reopen cross-model comparison in the manuscript unless it resolves a concrete reviewer-facing validity objection the primary lane cannot answer
