---
paper_id: quine-hotos
claim_bundle: quine-hotos
venue: hotos
format-envelope: workshop-short
track: system
cluster: process-native-systems-substrate
status: planning
state_mode: planning-skeleton
role: architectural-reframing
family: runtime-substrate
theme: runtime-ontology
theme_home: ../../../theory/views/substrate.md
publication_class: archival-candidate
maturity_target: synthesis
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [carrier, exec-barrier, trace-inheritance, delayed-addressability, coordination-closure, future-task-surface, structural-pull]
gaps: [revisable-time-semantics, branching-identity]
evidence: [experiments/active/early-existence/exst-01-ariadne/README.md, experiments/active/early-existence/exst-03-fractal/README.md, experiments/active/early-existence/exst-04-shell-mapreduce/README.md, experiments/active/early-existence/exst-10-introspective-quine/README.md, experiments/active/autopoiesis-succession/succ-01-minimal-self-birth/README.md, experiments/active/autopoiesis-succession/succ-02-blind-self-birth/README.md, experiments/active/autopoiesis-succession/succ-03-prompt-ablation-ladder/README.md]
view-affinity: [substrate]
---

# Quine (HotOS)

## Current Question

Why is the systems significance of Quine not merely that POSIX can host agents, but that native process primitives already form a viable agent substrate?

## Thesis

The systems significance of Quine is not only that POSIX can host agents, but that native process primitives already constitute a viable agent substrate.

## In Scope

- abstraction-level reframing for systems audiences
- the relation between runtime primitives and agent architecture
- how Quine changes what the OS is allowed to claim about agent systems

## Out Of Scope

- full implementation history
- later filesystem/world/time papers as if they were already runtime reality
- A-Life ontology as the main argumentative center

## Current State

- manuscript workspace: [`./WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- neighboring flagship paper: [`../quine-arxiv/README.md`](../quine-arxiv/README.md)

## Experiment Feeds

- `experiments/active/early-existence/exst-01-ariadne/README.md` - secondary - supporting-only - reframes `exec` as a substrate-level continuity operator rather than a mere implementation convenience
- `experiments/active/early-existence/exst-03-fractal/README.md` - secondary - supporting-only - shows `fork` as a substrate-native exploration operator rather than an application-layer coordination add-on
- `experiments/active/early-existence/exst-04-shell-mapreduce/README.md` - secondary - supporting-only - demonstrates that shell-native process semantics can complete coordination patterns without framework abstractions
- `experiments/active/early-existence/exst-10-introspective-quine/README.md` - secondary - supporting-only - meta-recursive delegation and self-use strengthen the architectural reframing of native process primitives as an agent substrate
- `experiments/active/autopoiesis-succession/succ-01-minimal-self-birth/README.md` - secondary - supporting-only - minimal self-reconstruction probes how little substrate machinery is needed before process-native self-production reappears
- `experiments/active/autopoiesis-succession/succ-02-blind-self-birth/README.md` - secondary - supporting-only - blind self-birth tests whether prompt scaffolding can be removed once physical affordances are explicit in runtime physics
- `experiments/active/autopoiesis-succession/succ-03-prompt-ablation-ladder/README.md` - secondary - supporting-only - prompt ablation measures how much of Quine's self-model can move from mission wording into system-prompt physics before behavior collapses
- `experiments/active/benchmarks/bench-01-terminal-bench/README.md` - context-only - exploratory-only - benchmark behavior is useful background for the systems reframing, but the current HotOS lane should not depend on leaderboard-style evaluation
- `experiments/active/benchmarks/bench-02-tb2-sweep/README.md` - context-only - exploratory-only - full-sweep outcomes are background context, not current flagship-paper evidence
- `experiments/active/benchmarks/bench-04-programbench/README.md` - context-only - exploratory-only - ProgramBench is useful architecture context, but the current scaffold is not yet cleanroom-valid benchmark evidence
- `experiments/active/benchmarks/bench-05-terminal-bench2-gpt54/README.md` - context-only - exploratory-only - useful background for the agent architecture, not current flagship evidence

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - no stable evidence backbone is currently declared for this paper

## Strongest Evidence

- the existing Quine runtime and behavior traces
- the flagship systems paper as the implementation-backed base
- compact examples that show native process semantics outperform framework metaphors

## Blocking Gaps

- the paper must avoid sounding like a manifesto detached from implementation
- the boundary between current reality and later paper program still needs careful policing
- claims need to stay venue-appropriate and concrete
- reverse-link note: benchmark-facing experiments should remain contextual unless the paper's scope broadens materially

## Next Useful Work

- keep the reframing downstream of the flagship proof
- borrow only evidence that already survives the flagship paper's discipline
- use this paper to sharpen systems consequences, not to absorb every future idea

## Unsafe Overclaims

- do not let this reframing paper become a manifesto detached from implementation-backed runtime reality
- do not blur later filesystem, world, or time papers into the HotOS boundary
- do not let architecture-level reframing outrun the flagship paper's evidence discipline
