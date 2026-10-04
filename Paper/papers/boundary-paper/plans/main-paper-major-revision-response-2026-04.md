# Main Paper Major-Revision Response: April 2026 Peer Review

This note records the April 2026 peer-review critique of the main boundary
paper and converts it into a manuscript repair plan. It is a planning surface,
not manuscript prose.

## Review Read

The review judges the paper as theoretically strong but not yet
benchmark-grade as an empirical paper:

- theoretical value is high because the paper turns `self` from a fixed
  referent into an experimental variable
- the strongest concepts are `identity collapse`, `identity climb`, and
  carrier-governed continuity
- the main weakness is empirical auditability: methods, prompts, logs, run
  accounting, exclusion rules, physical traces, and failure cases are not yet
  visible enough in the manuscript
- the right current posture is a strong conceptual framework with promising
  demonstrations, not a fully reproducible benchmark-grade empirical result

Best one-sentence thesis to preserve:

> Self-replication is underspecified until we say what object is reproduced
> and what carrier relation makes the successor continuous with the lineage.

## Required Repairs

### P0: Methods And Evidence Standard

The main paper needs a formal methods / evidence-standard layer before strong
empirical conclusions are allowed.

It must define and report:

- environment: OS, cleanroom setup, compiler, permissions, binary unlinking,
  `/proc/self/exe` status, network policy, provider access, and toolchain
  assumptions
- agent configuration: model/provider, temperature or determinism settings,
  context budget, tool access, system/runtime prompts, and runtime contract
- experiment matrix: condition, deprivation, prompt difference, `n`, success
  criterion, and failure modes for every promoted experiment
- trace instrumentation: file-access provenance, fd inheritance map, process
  tree, pre/post-`exec` evidence, stdout/fd capture, and hash comparisons
- scoring rubric: body copy, fresh synthesis, contract preservation,
  reproduced scaffolded interpretive organization, handoff continuity,
  carrier boundary, delayed addressability, return-to-idle viability, and
  control-surface closure
- exclusion and retention rules: what `retained run` means, total runs,
  excluded runs, whether exclusion was pre-registered, and what was learned
  from excluded or failed samples
- raw artifact routing: prompts, logs, hashes, minimal traces, successor
  source snippets, fd maps, and scorer outputs

### P0: Empirical Claim Downgrades

Use the conservative forms until the methods layer is complete:

- `6.59` demonstrates autonomous post-handoff closure on one repaired
  cleanroom lane, not regime-general closure
- `6.26` supports reproduced scaffolded interpretive organization, not
  intrinsic cognitive reproduction or cognitive autopoiesis
- `6.48 -> 6.50` supports renewal under viability pressure plus an externally
  supplied responsibility norm, not pure endogenous teleology
- `6.63 -> 6.65` supports operational selection by future-task carrier, not
  population-level evolutionary selection
- `6.64/6.65` support bounded control-surface closure under successor-mission
  scaffolding, not full operational or material closure

### P1: Prompt-Compliance Controls

`semantic target-shaping` is currently a strong interpretation, but the review
flags ordinary prompt compliance as the obvious alternative explanation. The
main paper should not present the mechanism as settled until controls are
either run or explicitly deferred.

Candidate controls:

- preserve target meaning while changing body salience
- preserve body salience while changing target-shaping
- paraphrase target-shaping across equivalent wordings
- adversarially prohibit body copying without supplying an organizational
  target
- compare against no-lineage ordinary coding-agent baselines
- test model-family variation
- use scorer-blind trace classification for reproduced-object level

### P1: Related Work Expansion And Verification

The related-work layer should be strengthened around:

- classical quines
- von Neumann self-reproducing automata
- Tierra / Avida / computational autopoiesis
- computer viruses and worms
- process migration, checkpoint/restart, hot code swap, actor identity, and
  distributed-system continuity
- RepliBench, SOCK, DGM, and related self-improving-agent systems

Citation discipline:

- verify primary sources before adding new claims to manuscript prose
- distinguish RepliBench/SOCK as risk and persistence benchmarks that already
  decompose practical replication capacities
- distinguish DGM-style work as inspectable/editable-codebase self-improvement,
  not mortal process-lineage continuity
- check whether `@tian2024selfimprovement` is actually about codebase editing
  before using it in an inspectable-source comparison

### P1: Safety And Release Posture

Because the paper concerns self-replication, persistence, and continuation,
the main paper needs a compact safety/release subsection:

- containment boundaries
- what artifacts are redacted or not released
- why experiments are process-local / cleanroom-limited
- what capability claims are intentionally not made

## Manuscript Address Strategy

The LBA should stay compact and only absorb the language downgrades. The main
paper must carry the major-revision burden.

Recommended main-paper structure:

1. Introduction: preserve the two questions and carrier-governed thesis.
2. Framework: define the two axes and avoid mystifying Quine.
3. Operational definitions: pass/fail criteria for each self-production level.
4. Methods / evidence standard: environment, run accounting, traces, scoring,
   retention/exclusion, and artifact routing.
5. What is reproduced? Axis 1 results.
6. Does it continue? Axis 2 results.
7. Control-surface closure: carrier contrast and delayed addressability.
8. Discussion: bounded claim, non-claims, safety posture, and next controls.

The current sectionized draft can keep its basic argument, but it needs a new
methods/evidence section before the result sections and should treat the
current PDF as a conceptual draft rather than submission-ready empirical prose.

## Light Formalization Candidate

Use only if it clarifies rather than bloats the paper:

- `B_t`: body artifact at time `t`
- `C_t`: runtime contract / normative self-description
- `I_t`: interpretive coupling to provider and environment
- `K_t`: task or mission state
- `H_t`: handoff carrier set, such as surviving fds, env, cwd, argv, streams,
  or control-surface handles
- `S_{t+1}`: successor process after `exec`

Relations:

- body copy: `B_{t+1}` is derived from `B_t`
- fresh synthesis: `B_{t+1}` is generated without reading or copying `B_t`
- contract preservation: behavior of `S_{t+1}` on held-out tasks satisfies
  `C_t`
- reproduced scaffolded interpretive organization: `S_{t+1}` reconstructs
  the coupling needed for novel mission interpretation
- lineage continuity: a surviving causal carrier path from `S_t` to
  `S_{t+1}` enables live task continuation

## Closure Criteria For This Review

This review is addressed only when:

- manuscript plan and drafting index require a methods/evidence-standard
  section
- section prose no longer relies on undefined `retained runs`
- promoted experiments have run-level accounting or are explicitly marked as
  frontier demonstrations
- related-work claims cite verified primary sources
- safety/release posture is visible
- the paper consistently uses reproduced scaffolded interpretive organization
  rather than bare `cognitive reproduction`
