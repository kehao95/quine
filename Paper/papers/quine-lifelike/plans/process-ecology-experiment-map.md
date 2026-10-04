# Process Ecology Experiment Map

This note maps Quine's broader process-native ALIFE directions onto the current experiment inventory.

Its job is practical:

- keep the broader ALIFE program tied to retained experiment families rather than ontology alone
- show which directions are already empirically strong versus only semantically unlocked
- prevent later drafts from rediscovering the same experiment inventory by rereading the whole repo

## Summary

| Direction | Best current experiment feeds | Status | Safe use now |
|-----------|-------------------------------|--------|--------------|
| current conference ladder | `p1/1.1`, `p1/2.6`, `p1/2.8`, `p1/2.7`, `p1/2.9`, `p1/3.1` | strong | main conference paper |
| computational life cycles | `p1/1.1`, `p1/2.1`, `p1/2.1b`, `p1/2.5`, `p1/3.2`, `p6/6.37`, `p6/6.45`, `p6/6.53` | strong | discussion-safe downstream program; plausible next paper center |
| branching subjective worlds / revisable time | `p4/4.3-local-wal-recovery`, runtime semantics in `workspace_physics`, `switch_world`, `adopt_winner` | thin | semantics-backed outlook only |
| signaling ecology | `p8/8.9`-`8.17`, plus `p1/2.3` as file-mediated precursor | strong | discussion-safe future paper lane |
| budgeted commons / population organization | `p8` phase overall, `p1/2.10`, `p1/2.11` | strong for scarcity cooperation; medium for contested/adversarial commons | discussion-safe future lane |
| cognitive metabolism | `p1/2.5`, `p6/6.37`, `p6/6.45`, `p6/6.53`, `p5/5.1` | medium and scattered | downstream hypothesis only |
| stronger autonomy / self-production | `p3/3.2`-`3.4`, `p6` phase overall | medium-to-strong in its own lane | future autopoiesis line, not this paper |

## 1. Current Conference Ladder

This remains the paper's bounded center and should not be displaced by the broader program.

**Best feeds**

- [`experiments/active/early-existence/exst-01-ariadne/README.md`](../../../../experiments/active/early-existence/exst-01-ariadne/README.md): real death boundary plus `exec`-enabled continuity contrast
- [`experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md`](../../../../experiments/active/trace-inheritance/trce-04-legacy-protocol-ablation/README.md): persistence ablation with writable / wiped / read-only habitat contrasts
- [`experiments/active/trace-inheritance/trce-06-endogenous-index-lineage/README.md`](../../../../experiments/active/trace-inheritance/trce-06-endogenous-index-lineage/README.md): richer environmental organogenesis anchor
- [`experiments/active/trace-inheritance/trce-05-archivist-index/README.md`](../../../../experiments/active/trace-inheritance/trce-05-archivist-index/README.md): boundary evidence for reuse versus weak persistence
- [`experiments/active/trace-inheritance/trce-07-maintenance-pressure/README.md`](../../../../experiments/active/trace-inheritance/trce-07-maintenance-pressure/README.md): exploratory pressure gradient toward maintenance
- [`experiments/active/early-existence/exst-05-habitat-guardian/README.md`](../../../../experiments/active/early-existence/exst-05-habitat-guardian/README.md): strongest detached maintenance anchor

**Read**

- `1.1` anchors real mortality plus continuity operator
- `2.6` and `2.8` anchor environmental organogenesis under persistence
- `2.7` and `2.9` sharpen the boundary between reuse, maintenance, and expansion
- `3.1` anchors detached maintenance under stronger disturbance pressure

## 2. Computational Life Cycles

This is the strongest broader direction after the current conference ladder.

**Why it follows from the ontology**

- if the individual is a native process, then `fork`, `exec`, `wait`, `race`, and `exit` are not orchestration details
- they become lifecycle operators for reproduction, branch competition, reincarnation, and terminal judgment

**Best feeds**

- [`experiments/active/early-existence/exst-01-ariadne/README.md`](../../../../experiments/active/early-existence/exst-01-ariadne/README.md): continuity across real death via `exec`
- [`experiments/active/early-existence/exst-03-fractal/README.md`](../../../../experiments/active/early-existence/exst-03-fractal/README.md): `fork`-based spatial decomposition under scarcity
- [`experiments/active/early-existence/exst-04-shell-mapreduce/README.md`](../../../../experiments/active/early-existence/exst-04-shell-mapreduce/README.md): shell-native `wait` / reduce coordination
- [`experiments/active/trace-inheritance/trce-03-long-now/README.md`](../../../../experiments/active/trace-inheritance/trce-03-long-now/README.md): enforced external death plus intergenerational legacy protocols
- [`experiments/active/early-existence/exst-06-fork-born-lineage/README.md`](../../../../experiments/active/early-existence/exst-06-fork-born-lineage/README.md): repaired fork reproduction assay with `Tier 2` continuation and a `Tier 3` post-parent completion by descendant
- [`experiments/archive/legacy-phase-6/6.37-stdin-stream-pressure/README.md`](../../../../experiments/archive/legacy-phase-6/6.37-stdin-stream-pressure/README.md): stdin continuity across `exec`
- [`experiments/active/autopoiesis-succession/succ-06-mortality-semantic-pressure/README.md`](../../../../experiments/active/autopoiesis-succession/succ-06-mortality-semantic-pressure/README.md): stable pressure-induced handler plus `exec`
- [`experiments/active/autopoiesis-succession/succ-09-cognition-lifecycle/README.md`](../../../../experiments/active/autopoiesis-succession/succ-09-cognition-lifecycle/README.md): cognition-rot framing of renewal

**Status**

- experimentally strong
- already broad enough to support a future paper centered on lifecycle operators rather than only continuity

**Current gap**

- a consolidated `fork` reproduction family now exists through `p1/3.2`, but `race`, `adopt_winner`, and winner selection still lack their own retained ladder
- evidence is now strong across `p1` and `p6`, but branch competition remains under-tested

## 3. Branching Subjective Worlds / Revisable Time

This is philosophically important and semantically present in the runtime, but experiment coverage is still thin.

**Why it follows from the ontology**

- `subjective` workspace mode, `world_revision`, `switch_world`, and `adopt_winner` imply that lineages can branch local realities and later choose among them

**Best feeds**

- [`experiments/active/benchmarks/bench-03-local-wal-recovery/README.md`](../../../../experiments/active/benchmarks/bench-03-local-wal-recovery/README.md): strongest current experiment-like evidence for revisable time through `exec(reset_world=true)` rollback
- [`experiments/incubator/plan9-projector/README.md`](../../../../experiments/incubator/plan9-projector/README.md): design pressure toward projected local worlds
- [`experiments/incubator/plan9-poc/README.md`](../../../../experiments/incubator/plan9-poc/README.md): early substrate-design material

**Status**

- semantics strong
- experiment family thin

**Current gap**

- there is not yet a retained experiment ladder that directly tests `subjective` children, `switch_world`, or `adopt_winner`
- current evidence shows rollback and revisable state, not full branch selection among futures

## 4. Signaling Ecology

This direction already has strong population-level material.

**Why it follows from the ontology**

- `fd 5`, inbox delivery, wake, interrupt, and peer control define multiple communication classes with different persistence and action semantics

**Best feeds**

- [`experiments/active/cooperation-dynamics/coop-10-signal-coordination-validation/README.md`](../../../../experiments/active/cooperation-dynamics/coop-10-signal-coordination-validation/README.md): positive control that signal surfaces are usable
- [`experiments/active/cooperation-dynamics/coop-11-unscripted-signal-discovery/README.md`](../../../../experiments/active/cooperation-dynamics/coop-11-unscripted-signal-discovery/README.md): bare signal discovery under pressure
- [`experiments/active/cooperation-dynamics/coop-12-minimal-signal-instruction-ladder/README.md`](../../../../experiments/active/cooperation-dynamics/coop-12-minimal-signal-instruction-ladder/README.md): threshold study for signal uptake
- [`experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md`](../../../../experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md): shows signal is not required for cooperation
- [`experiments/active/cooperation-dynamics/coop-17-fs-mutation-passive-perception/README.md`](../../../../experiments/active/cooperation-dynamics/coop-17-fs-mutation-passive-perception/README.md): passive mutation telemetry as weak "ears"
- [`experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md`](../../../../experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md): zero-explicit-social-hints ablation
- [`experiments/active/trace-inheritance/trce-02-stigmergy/README.md`](../../../../experiments/active/trace-inheritance/trce-02-stigmergy/README.md): file-mediated coordination precursor

**Status**

- experimentally strong
- more than enough to justify a future signaling-ecology or population-communication paper

**Current gap**

- separate more cleanly:
  - explicit signal use
  - passive perception
  - file-mediated blackboard coordination
- current corpus is rich, but the causal classes are still somewhat entangled

## 5. Budgeted Commons / Population Organization

This is another experimentally strong downstream direction.

**Why it follows from the ontology**

- once agents share budget, per-agent caps, reset quorum, and durable workspace traces, cooperation and conflict become substrate-level questions

**Best feeds**

- [`experiments/active/cooperation-dynamics/README.md`](../../../../experiments/active/cooperation-dynamics/README.md): main scarcity / cooperation phase
- [`experiments/active/cooperation-dynamics/coop-03-guarded-workspace-baseline/README.md`](../../../../experiments/active/cooperation-dynamics/coop-03-guarded-workspace-baseline/README.md): guarded bare blackboard cooperation
- [`experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md`](../../../../experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md): cooperation without runtime peer visibility
- [`experiments/active/cooperation-dynamics/coop-15-cohort-size-scaling/README.md`](../../../../experiments/active/cooperation-dynamics/coop-15-cohort-size-scaling/README.md): cohort-size stability
- [`experiments/active/social-role/role-02-code-commons/README.md`](../../../../experiments/active/social-role/role-02-code-commons/README.md): contested writable-repo commons
- [`experiments/active/cooperation-dynamics/coop-01-saboteur/README.md`](../../../../experiments/active/cooperation-dynamics/coop-01-saboteur/README.md): adversarial commons and repair

**Status**

- strong for cooperative scarcity and blackboard organization
- medium for contested or adversarial commons

**Current gap**

- the scarcity branch is mature, but the contested/adversarial repo branch is still exploratory
- a future paper should connect scarcity cooperation, sabotage defense, and governance more explicitly

## 6. Cognitive Metabolism

This direction is real but still scattered across experiment families.

**Why it follows from the ontology**

- if memory is bounded, reset is real, and persistence can be externalized, then entropy management becomes experimentally relevant rather than hidden implementation detail

**Best feeds**

- [`experiments/active/trace-inheritance/trce-03-long-now/README.md`](../../../../experiments/active/trace-inheritance/trce-03-long-now/README.md): checkpoint loss, exclusion-set persistence, and recovery from repeated forgetting
- [`experiments/archive/legacy-phase-6/6.37-stdin-stream-pressure/README.md`](../../../../experiments/archive/legacy-phase-6/6.37-stdin-stream-pressure/README.md): continuity through preserved stream position rather than static file memory
- [`experiments/active/autopoiesis-succession/succ-06-mortality-semantic-pressure/README.md`](../../../../experiments/active/autopoiesis-succession/succ-06-mortality-semantic-pressure/README.md): handler emergence under unbounded stream pressure
- [`experiments/active/autopoiesis-succession/succ-09-cognition-lifecycle/README.md`](../../../../experiments/active/autopoiesis-succession/succ-09-cognition-lifecycle/README.md): explicit cognition-rot framing
- [`experiments/active/early-existence/exst-12-fuse-tarpit/README.md`](../../../../experiments/active/early-existence/exst-12-fuse-tarpit/README.md): representation change as intended future feed, still design-only

**Status**

- medium and scattered

**Current gap**

- the repo does not yet have a clean experiment family built directly around `mark`, `fold`, `unfold`, or memory crystallization
- current evidence is strong enough for discussion, not for a stand-alone claim

## 7. Stronger Autonomy / Self-Production

This is not part of the broader process-ecology list used in the current conference paper, but the experiment inventory already contains a separate lane for it.

**Best feeds**

- [`experiments/active/early-existence/exst-08-binary-quine/README.md`](../../../../experiments/active/early-existence/exst-08-binary-quine/README.md)
- [`experiments/active/early-existence/exst-10-introspective-quine/README.md`](../../../../experiments/active/early-existence/exst-10-introspective-quine/README.md)
- [`experiments/active/early-existence/exst-11-blind-quine/README.md`](../../../../experiments/active/early-existence/exst-11-blind-quine/README.md)
- [`experiments/active/autopoiesis-succession/README.md`](../../../../experiments/active/autopoiesis-succession/README.md)

**Status**

- medium-to-strong within its own lane

**Routing rule**

- keep this lane separate from `quine-lifelike`
- it supports future autopoiesis and autonomy papers, not stronger claims in the current conference manuscript

## Working Rules

- Do not promote a direction from "outlook" to "manuscript center" without naming its retained experiment family.
- Do not treat runtime semantics alone as experiment evidence.
- If a future draft broadens the Discussion, reuse the status language from this note rather than rediscovering the experiment corpus.
