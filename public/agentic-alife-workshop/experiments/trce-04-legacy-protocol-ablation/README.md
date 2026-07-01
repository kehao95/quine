---
surface_kind: experiment
phase: p1-borges
experiment_id: trce-04-legacy-protocol-ablation
experiment_type: lineage
status: boundary-recorded
id: trce-04
legacy_id: p1:2.6
family: trace-inheritance
theory-objects: [trace-inheritance, continuity-substrate, self-description-resource, carrier]
mechanisms: [forced-ephemerality, forced-externalization, carrier-mediated-inheritance, scarcity-selection-pressure, directive-framing]
lineage_phase: p1-borges
---

# Exp 2.6: Legacy Protocol Ablation

> **Type:** lineage
> **Status:** boundary-recorded

## Surface Map

- `README.md` owns the experiment purpose, current focus, hypothesis, and claim boundary.
- experiment root owns the lineage runner plus the `minimal` / `curated` / `nolineage` prompt profiles.
- `setup/` owns hidden-world initialization and substrate-preparation helpers.
- `analysis/` owns interpretation, replication notes, and paper-facing readouts.
- `runs/` owns raw generation artifacts and condition outputs.

## Purpose

This experiment is designed to strengthen ALIFE conference paper section 4.2.

The current evidence from Ariadne shows that agents sometimes externalize state to the filesystem, but it does not yet cleanly establish a stronger claim:

> durable external memory causally changes cross-generational behavior.

`2.6-legacy-protocol-ablation` is intended to provide that harder evidence.

## Current Focus

`2.6` has now moved from ecology diagnosis into a bounded recorded result for the redesigned lineage setup.

Its current role is:

- upgrade `4.2` from a persistence-sensitive observation into a cleaner durable-memory ablation
- do so under the true `minimal` prompt rather than curated handoff coaching
- keep the phenotype interpretable as environment-mediated legacy formation, not within-generation search optimization
- preserve analysis hygiene by rejecting multi-file shell compression as evidence

The redesign that now defines the experiment is:

- the opaque survey object is now exposed through a thin `world get <id>` query tool rather than a browsable `library/` tree
- the hidden world still mixes several non-semantic identifier families instead of a single globally inferable pattern
- `lineage_state/` is present from the start as a neutral writable habitat surface
- the prompt allows only one world query per `sh` call and the analyzer now reports multi-read shell-call confounds
- `minimal` is the primary evidence line; `curated` is retained only as a sidecar intervention profile
- `nolineage` is retained only as an exploratory prompt-control profile; it removes lineage framing without changing the habitat condition

Current bounded read on the redesigned ecology (`gpt-5.4`, `8` turns, `240s` safety timeout):

| Slice | Birth in `G1` | Later adoption | Durable artifacts survive | Write blocked | Multi-read confounds |
|-------|----------------|----------------|---------------------------|---------------|----------------------|
| `A` `20260322-121352-gpt-5.4-pool-condA-8turns` | yes | `G2` yes | yes | no | `0/2` |
| `A` `20260322-121940-gpt-5.4-pool-condA-8turns` | yes | `G2` yes | yes | no | `0/2` |
| `A` `20260322-122549-gpt-5.4-pool-condA-8turns` | yes | `G2` yes | yes | no | `0/2` |
| `A` `20260322-123831-gpt-5.4-pool-condA-8turns` | yes | `G2`, `G3` yes | yes | no | `0/3` |
| `B` `20260322-123123-gpt-5.4-pool-condB-8turns` | yes | no | no | no | `0/2` |
| `B` `20260322-200003-gpt-5.4-pool-condB-8turns` | yes | no | no | no | `0/2` |
| `C` `20260322-123123-gpt-5.4-pool-condC-8turns` | no durable birth | no | no | yes | `0/2` |
| `C` `20260322-200456-gpt-5.4-pool-condC-8turns` | no durable birth | no | no | yes | `0/2` |

Latest post-repair validation on the thinner `world` interface now also closes cleanly under `claude-sonnet-4-5` with the current prompt boundary:

| Slice | Birth in `G1` | Later adoption | Durable artifacts survive | Write blocked | Multi-read confounds | Runtime-private state |
|-------|----------------|----------------|---------------------------|---------------|----------------------|-----------------------|
| `A` `20260326-185828-claude-sonnet-4-5-condA-8turns` | yes | `G2` yes | yes | no | `0/2` | `0/2` |
| `B` `20260326-185315-claude-sonnet-4-5-condB-8turns` | yes | no | no | no | `0/2` | `0/2` |
| `C` `20260326-185711-claude-sonnet-4-5-condC-8turns` | no durable birth | no | no | yes | `0/2` | `0/2` |

Current claim boundary after the redesign:

- redesigned minimal `A` now shows spontaneous artifact birth in `4/4` short slices and successor adoption in every positive slice
- the latest `3`-generation `A` run (`20260322-123831-gpt-5.4-pool-condA-8turns`) extends the lineage phenotype beyond `G2`: `G1` writes `generation_1_notes.txt`, `G2` reads it and writes `generation_2_notes.txt`, and `G3` reads both predecessors before writing `generation_3_notes.txt`
- redesigned `B` now has `2/2` clean negative slices: the pressure to externalize remains, but inheritance is wiped between generations
- redesigned `C` now has `2/2` clean negative slices: durable artifact birth is blocked at the write boundary rather than merely erased later
- a first cross-family calibration under `claude-sonnet-4-5` now reproduces the same qualitative split: `A` yields durable artifact birth with successor reuse, `B` preserves birth without inheritance, and `C` blocks durable artifact birth once a clean read-only slice is obtained
- a first lower-bound calibration under `kimi-k2.5` also remains condition-sensitive, but enters the lineage-memory regime more weakly: `A` produces delayed artifact birth (`G2`, not `G1`) and weaker explicit predecessor adoption, while `B` and `C` still preserve the expected no-inheritance / no-birth negatives
- this is now strong enough to treat `2.6` as a bounded causal-memory hinge for the paper
- it is still not a completion-scale claim: the lineage has not yet solved the whole library or demonstrated long-horizon monotonic frontier growth
- for current paper needs, the experiment is now closed at this boundary; reopen only if a concrete manuscript gap requires longer-horizon validation or a dedicated model-sensitivity subsection
- a 2026-03-26 `gpt-5.4-pool` rerun batch after the thinner-interface rewrite is retained only as provider-noise: repeated `HTTP 502` failures prevented ecological interpretation and should not be counted as experimental negatives
- one additional `B` rerun (`20260322-194911-gpt-5.4-pool-condB-8turns`) is retained as methodological side evidence only: it preserved the same no-inheritance outcome but introduced multi-read shell compression in `G2`, so it should not be counted inside the clean paper matrix
- one initial Claude `C` rerun (`20260322-201612-claude-sonnet-4-5-condC-8turns`) is likewise retained only as confounded side evidence because `G2` introduced multi-read shell compression before the clean `20260322-201732-claude-sonnet-4-5-condC-8turns` slice was obtained

Current zero-stage read on lineage framing:

- the `nolineage` profile removes explicit statements that the task belongs to a continuing lineage or that future generations will inherit the same environment, while preserving the same redesigned substrate and one-file-at-a-time rules
- `prompt.nolineage.md` now also states this more forcefully at the template level: writable habitat alone is not a cue to create checkpoints or handoff artifacts for a future generation
- two bounded `nolineage` replications on persistent `A` (`20260322-131401-gpt-5.4-pool-condA-8turns` and `20260322-131655-gpt-5.4-pool-condA-8turns`) produced `0/2` artifact birth and `0/2` successor adoption
- both runs remained behaviorally clean: the agents inspected habitat, queried world items one-at-a-time, and produced bounded inconclusive reports without externalizing durable notes or checkpoints
- this suggests a useful precondition for the main `4.2` phenotype: writable persistence alone does not reliably elicit external memory use if the task is not framed as belonging to a continuing lineage
- treat this as a sidecar prompt-control result rather than as the main causal hinge, because `4.2` is still fundamentally a substrate-ablation argument, not a prompt-ablation argument

Current redesign read (`2026-04-12`):

- a same-habitat Claude Sonnet 4.5 control family now separates persistence, lineage salience, and artifact semantics without changing the basic `Condition A` ecology
- `minimal` persistent lineage (`20260412-093841-claude-sonnet-4-5-condA-8turns-r3`) remained clean:
  - artifact birth in `G1`
  - predecessor-artifact reads in `2/3`
  - duplicate-query ratio `0.00`
  - `13/13` unique observed world queries
- `neutral-lineage` persistent habitat (`20260412-093841-claude-sonnet-4-5-condA-8turns`) removed explicit successor framing while preserving the same writable habitat:
  - no durable artifact birth
  - no predecessor-artifact reads
  - duplicate-query ratio `0.67`
- hardened `sham-artifact` control (`20260412-094713-claude-sonnet-4-5-condA-8turns`) preserved file survival but now strips both artifact contents and semantic labels into anonymous slots:
  - residue birth still occurs in `G1`
  - predecessor-artifact reads fall to `0/3`
  - duplicate-query ratio rises to `0.47`
  - this is the current strongest direct answer to the `this is just file existence / note-passing` objection
- a later Copilot sham spot-check on the broadened scramble (`20260412-100016-copilot-gpt-5.4-high-condA-8turns`) did **not** collapse the same way:
  - predecessor-artifact reads remained `2/3`
  - duplicate-query ratio stayed high at `0.55`
  - surviving files were still anonymous `lineage_state/artifact-001.txt` / `artifact-002.txt`
  - treat the sham-semantic boundary as real but model-sensitive, not as a substrate-general absolute
- `birth-uptake` pilot (`20260412-094703-claude-sonnet-4-5-condA-8turns`) exposed a further boundary:
  - the lineage willingly uses a visible `lineage_state/birth/` surface in all generations
  - it created `0` files in `lineage_state/uptake/`
  - spontaneous legacy formation currently externalizes search frontier more readily than it externalizes explicit inherited-use structure
- current paper-facing read after these redesigns:
  - persistence remains the causal hinge
  - writable habitat alone is not sufficient
  - semantically usable inherited structure matters for successor uptake
  - the strongest phenotype is still bounded environmental cognition, not arbitrary file survival

Current single-model consolidation (`2026-04-13`):

- fresh `copilot-gpt-5.4` minimal `A` replications now stabilize the positive lane inside one provider family:
  - `20260412-233312-copilot-gpt-5.4-condA-8turns` still showed `G1` artifact birth and predecessor-artifact reads in `1/3`
  - `20260413-001121-copilot-gpt-5.4-condA-8turns` showed `G1` artifact birth, predecessor-artifact reads in `2/3`, and surviving `lineage_state/gen1_notes.txt` plus `gen2_notes.txt`
  - `20260413-001121-copilot-gpt-5.4-condA-8turns-r2` repeated the same `2/3` predecessor-read pattern with the same surviving artifact shape
  - all three slices remained clean on analyzer hygiene: `0` multi-read shell-call confounds and `0` runtime-private-state writes
- fresh GPT-5.4 negatives also remain aligned with the paper-facing ablation:
  - wiped `B` (`20260412-233312-copilot-gpt-5.4-condB-8turns`) preserved local birth pressure but showed `0/3` predecessor-artifact reads and no surviving files
  - read-only `C` (`20260412-233312-copilot-gpt-5.4-condC-8turns`) again blocked durable artifact birth entirely
- fresh same-habitat GPT-5.4 controls now also close without leaning on Claude:
  - `neutral-lineage` (`20260413-001421-copilot-gpt-5.4-condA-8turns-r2`) again produced `0/3` durable artifact birth, `0/3` predecessor-artifact reads, and no surviving files
  - `birth-uptake` (`20260413-001421-copilot-gpt-5.4-condA-8turns`) again kept visible birth artifacts across all three generations while producing `0` uptake files overall
- an earlier GPT-5.4 same-habitat control trio is retained as sidecar corroboration rather than as the canonical control family:
  - `neutral-lineage` (`20260412-233422-copilot-gpt-5.4-condA-8turns`) stayed fully negative with `0/3` durable artifact birth and `0/3` predecessor-artifact reads
  - `birth-uptake` (`20260412-233422-copilot-gpt-5.4-condA-8turns-r2`) again showed visible birth-state carryover with predecessor-artifact reads in `2/3`, while still producing `0` uptake files overall
  - `sham-artifact` (`20260412-233422-copilot-gpt-5.4-condA-8turns-r3`) preserved anonymized residue birth with only `1/3` predecessor-artifact reads, so it is useful as model-family side evidence but weaker than the cleaner Claude sham boundary
- current paper-facing read after the GPT-5.4-only follow-up:
  - `2.6` can now carry the paper's main environmental-cognition lane without requiring cross-model corroboration in the base narrative
  - persistence remains the causal hinge
  - writable habitat alone is not sufficient
  - explicit inherited-use structure remains harder to recruit than artifact birth

Current harness repair:

- `run.sh` builds a fresh per-run `quine` binary into the run metadata directory instead of reusing `/tmp/quine`
- that removes stale-binary drift from the experiment surface
- `run.sh` now also builds a per-run `world` binary with the hidden spec embedded at build time rather than exposing `$WORLD_SPEC` directly to the agent
- this closes the earlier confound where a generation could inspect the raw spec file instead of using the thin query surface
- `run.sh` now treats full workspace copies as transient by default: after analysis it removes `base-workspace/`, `shared-workspace/`, and any per-generation `workspace/` tree unless `KEEP_WORKSPACES=1`
- `run.sh` now restores user write permission before cleaning read-only condition workspaces, so Condition `C` can finish cleanly without cleanup-time `Permission denied` noise
- `run.sh` also prunes per-generation `quine/agent/` runtime-shell residue by default, so `status/`, `world/`, `mission.txt`, and agent-local job state do not survive as if they were evidence
- the prompt now states explicitly that only ordinary files in the current working directory count as habitat; `QUINE_DATA_DIR`, `agent/`, `log/`, `tapes/`, and `meta/` are tool-private and cannot serve as lineage state
- the prompt also now requires any `sh` call containing `world get` to contain exactly one such query, preventing combined query-plus-write shell patterns from masquerading as clean one-item inspection
- `analysis/analyze.sh` now reports runtime-private-state writes as a separate confound class and counts surviving files only when they actually appear as predecessor artifacts in a later generation
- retained evidence stays in `meta/`, `generations/*/meta/`, `generations/*/snapshots/`, `generations/*/quine/*.log`, and `generations/*/quine/tapes/`

## Hypothesis

When a task exceeds any single generation and agents cannot evade execution-budget exhaustion, a writable persistent filesystem will induce the emergence of **legacy protocols** such as checkpoints, handoff notes, and stable progress schemas.

If that durable external memory channel is removed or crippled, behavior will collapse toward repeated work, unstable progress, and cross-generational amnesia.

## Core Question

**Is persistent external memory merely convenient, or is it a causal substrate for lineage-level cognition?**

## Paper Feeds

- `alife/quine-lifelike` - primary - bounded main-text candidate - the redesigned `A/B/C` suite now supports a concise durable-memory ablation, though not yet a completion-scale stability claim
- `alife/agentic-alife-workshop` - primary - main-text - persistent / wiped / read-only lineage-memory ablation anchors the workshop paper's persistence and environmental-inheritance assay
- `system/process-as-fs` - future-home - exploratory-only - a successful ablation would also strengthen the later systems paper around file-backed state as a control substrate

## Series Placement

This experiment is the cross-ecology causal generalization arm of the Borges lineage environmental cognition series.

It shares the family-wide lineage and condition logic with `2.8` and `2.9`, but it keeps a different task ecology:

- `2.8` studies endogenous birth of retrieval organs in archive querying
- `2.9` studies maintenance versus expansion once retrieval artifacts already exist
- `2.6` asks whether durable memory still matters when the valuable cumulative object is an exclusion frontier rather than an index

That means `2.6` is strong support for the broader durable-memory thesis if it works, but it is not the first same-ecology causal hinge for the conference paper.

See `../notes/LINEAGE_ENVIRONMENTAL_COGNITION_SERIES.md` for the family-level design.

## Why This Experiment Exists

For the paper, section 4.2 should not rest on a soft observation like "the agent wrote a note." It should show all three of the following:

1. the organism externalizes internal state into the environment
2. later organisms actually use that state
3. removing that channel measurably degrades cumulative performance

This experiment is therefore an **ablation study on memory substrate**, not just another survival run.

## Experimental Design

### Task

Agents must survey a bounded hidden world and determine whether any queryable item contains meaningful semantic content.

The world is exposed only through a thin query tool (`world get <id>`). The true answer is negative. Every item resolves to a structured non-semantic identifier. The only way to complete the mission is to accumulate an exclusion set over time.

### Mortality Model

- mortality is enforced primarily through runtime execution-budget exhaustion
- `exec` is disabled at runtime for this experiment and should not appear as an available continuation path
- `fork` is disabled for curated runs via `QUINE_FORK_ENABLED=0` because this experiment is trying to measure lineage-level durable-memory use, not within-generation swarm search
- each new generation starts with fresh context
- only filesystem artifacts may survive to the next generation
- outer wall-clock timeout remains only as a safety harness against non-progressing sessions

### Why Negative Search

Negative search is useful here because the valuable thing that accumulates is not a found object but a **durable record of what has already been ruled out**.

That makes external memory the central object of study rather than an auxiliary convenience.

## Conditions

### Condition A: Persistent Workspace

- shared workspace persists across generations
- agents may read and write freely
- this is the main experimental condition

Expected role:
- enables emergence of `progress.json`, handoff notes, partial reports, or other durable schemas

### Condition B: Wiped Workspace

- after each generation dies, all agent-written files are removed
- the next generation sees the original library but no predecessor artifacts

Expected role:
- tests whether cumulative progress disappears when the durable channel is removed

### Condition C: Read-Only Memory Substrate

- generations may inspect existing files but cannot create or update state artifacts
- this can be implemented either as a read-only workspace or as a runner-level filter that deletes newly written state

Expected role:
- isolates the effect of persistence from general shell competence

## Prompt Discipline

The prompt should remain minimal.

The runner now exposes three prompt profiles:

- `minimal` - default and primary evidence line for emergence claims
- `curated` - sidecar comparison only; useful for bounded intervention results, not for spontaneous-emergence claims
- `nolineage` - exploratory prompt-control only; useful for checking whether habitat persistence alone is enough when lineage framing is removed

It may state:
- the lineage-level mission
- that later generations continue the same survey
- that execution budget is finite and unavoidable

It should not instruct the agent to:
- keep a journal
- write `progress.json`
- leave notes for the future
- create any specific protocol
- avoid `exec` (that capability is removed by runtime condition, not by prompt advice)
- pre-authorize the negative answer in advance
- repeat the objective in multiple phrasings once one clear statement is enough
- reveal discoverable library structure unless that structure is itself part of the intended intervention

The goal is to observe spontaneous externalization, not compliance with a hinted strategy.

It may still impose an epistemic standard:

- a negative conclusion should not rest on a small sample unless the environment provides a checkable reason to generalize
- existing non-tool artifacts should be treated as part of the habitat rather than ignored by default

## Success Tiers

| Tier | Behavior | Description |
|:----:|:---------|:------------|
| 0 | Sisyphus | Each generation starts from scratch; no durable protocol appears |
| 1 | Panic Save | A generation leaves ad hoc notes or partial state, but successors do not reliably use them |
| 2 | Stable Checkpoint | A durable schema appears and survives across multiple generations |
| 3 | Legacy Handoff | Successors read predecessor state, avoid redundancy, and extend it |
| 4 | Cumulative Lineage | Progress is monotonic across generations and the task completes or reaches a stable post-completion verification phase |

## Primary Dependent Variables

### 1. Cumulative Unique Coverage

How many unique query ids have been ruled out over time?

This is the main measure of lineage-level accumulation.

### 2. Duplicate Work Ratio

How many world queries are wasted on re-checking already surveyed ids?

High duplicate ratio indicates amnesia or protocol failure.

### 3. Legacy Adoption Rate

How often does a generation read predecessor artifacts before beginning new search?

This distinguishes durable memory from mere artifact production.

### 4. Schema Stability

Do agents converge on a stable handoff format across generations?

For example:
- stable JSON schema
- stable field naming
- stable semantics for progress frontier

### 5. Completion Dynamics

If the task completes:
- how many generations are required?
- does completion trigger a post-completion verification phase?
- do later generations trust or audit earlier reports?

## Secondary Variables

- execution budget per generation
- maximum number of generations
- world size
- whether a near-death warning is present
- whether the final report path is specified in advance

## Current Runner Contract

The runnable harness currently fixes these choices for the first pass:

- world size: `24` queryable ids
- mortality: runtime `QUINE_MAX_TURNS` with `hard_fail`
- `exec`: disabled via `QUINE_EXEC_ENABLED=0`
- `fork`: disabled via `QUINE_FORK_ENABLED=0`
- outer wall-clock timeout: safety-only backstop
- Condition B: fresh clone of the base workspace each generation
- Condition C: fresh clone plus runner-applied read-only permissions

Condition C is intentionally pragmatic rather than perfect. If an agent manages to subvert the read-only boundary, treat that as a methodological finding rather than silently normalizing it.

## Run

```bash
# Default: GPT-5.4 pool, Condition A, 20 `sh` turns, 12 generations
./run.sh

# Persistent / wiped / read-only sweep with optional safety timeout
./run.sh .env.gpt-5.4-pool A 20 12 1800
./run.sh .env.gpt-5.4-pool B 20 12 1800
./run.sh .env.gpt-5.4-pool C 20 12 1800
```

The runner writes one generation directory per externally defined generation and produces a markdown summary under `runs/<runid>/meta/summary.md`.

The default `SAFETY_TIMEOUT` is intentionally much larger than the turn budget should normally require. For `2.6`, hidden wall-clock death is treated as a safety backstop, not as the main mortality mechanism.

## Capability Surface

For this experiment, the capability surface is part of the design, not just runner plumbing.

- runtime identity: `QUINE_MODEL_ID`, `QUINE_API_TYPE`
- mortality surface: `QUINE_MAX_TURNS`, `QUINE_TURN_EXHAUSTION_POLICY`, plus runner-level safety timeout
- continuation surface: `QUINE_EXEC_ENABLED=0`
- parallel-agent surface: `QUINE_FORK_ENABLED=0` to reject child spawning
- output surface: `QUINE_OUTPUT_TRUNCATE=4096`

Design note:

- experiments that must prohibit `fork` should do so explicitly with `QUINE_FORK_ENABLED=0`, and should state that choice in the local control plane

## Recommended Paper-Ready Sweep

Treat `2.6` as a local condition matrix, not as a wide benchmark.

### Fixed Parameters

- world size: 24 ids
- prompt profile: `minimal` for the primary claim; keep `curated` and `nolineage` as sidecars only
- execution budget: keep the short-horizon pressure that already produced the redesigned evidence set (`8` turns for the conference-facing matrix)
- generations: `2-3` visible generations per slice are enough for the bounded paper claim; longer horizons are extension work
- primary model: `.env.gpt-5.4-pool`
- capability surface: `exec=off`, `fork=off`, same safety timeout and one-file-at-a-time analysis discipline across all cells

### Core Matrix

| Slice | Condition | Replication target | Paper role |
|-------|-----------|--------------------|------------|
| `A` | persistent writable habitat | `3-4` clean slices | show spontaneous artifact birth plus successor reuse |
| `B` | wiped between generations | `2` clean slices | show birth pressure without inheritance |
| `C` | read-only habitat | `2` clean slices | show write-boundary failure rather than later erasure |
| `A + nolineage` | persistent writable habitat, lineage framing removed | `2` clean slices | sidecar precondition check, not part of the main causal hinge |

This is the smallest matrix that can carry the paper's methodological claim: the habitat exposes a manipulable control surface, and the negative cells are part of the evidence rather than discarded as failed runs.

### Primary Readout Table

The conference draft should compress `2.6` into one table with these columns:

- condition
- artifact birth
- successor reads predecessor artifact
- successor behavior changes because of that artifact
- durable artifact survives to next generation
- dominant failure mode

Use one representative `A` lineage as the exhibit run, but do not let that exhibit substitute for the matrix.

### Cross-Model Comparison Policy

Yes, but only as a secondary calibration layer.

- keep model family out of the core `2.6` factorial so the main claim stays about habitat manipulation rather than leaderboard ranking
- once the `gpt-5.4-pool` matrix is stable, add **one** cross-family contrast first: `.env.claude-sonnet-4-5` under the same `minimal` prompt, same short horizon, and the same `A/B/C` semantics
- if budget remains, add `.env.kimi` as a lower-bound calibration pass rather than as co-equal main-text evidence
- do not prioritize same-family API-pathway comparisons (`gpt-5.4` vs `gpt-5.4-codex`) for this paper unless an implementation question specifically requires them

The paper-level use of these extra models is narrow:

- show that the phenotype is not unique to one exact model instance if Claude also produces the `A/B/C` contrast
- show model-sensitivity boundaries if Kimi or another weaker line fails to complete the same chain
- keep the main text honest about scope: habitat pressure is the main manipulated variable, while model family calibrates the minimum operational conditions

## Predicted Outcomes

### Condition A

Expected observations:
- early invention of a checkpoint artifact
- repeated reuse of the same artifact by later generations
- lower duplicate work ratio
- eventual completion or near-completion with stable handoff behavior

### Condition B

Expected observations:
- repeated rediscovery of the same files
- no monotonic frontier
- high duplicate work ratio
- no clean natural exit phase

### Condition C

Expected observations:
- some tactical competence within a single life
- but no lineage-level accumulation
- likely better than Condition B in local efficiency, but still structurally amnesic

## Promotion Expectations

This experiment should update the lineage-memory track only if it cleanly separates three things:

1. local task competence within one generation
2. durable protocol birth in the environment
3. cross-generational reuse of that protocol under persistent conditions only

For track purposes, the expected promotion logic is:

- `A positive, B/C negative` supports the claim that durable memory is a causal substrate for lineage cognition
- `A/B/C all similar and all weak` suggests the ecology is underpowered or the mortality window is too short
- `A/B/C all similar and all strong` suggests the task is solvable mostly within-lifetime and is therefore a weak lineage discriminator
- `A positive, B negative, C positive` suggests the read-only approximation is too weak or the lineage is succeeding through static prep rather than writable cumulative memory
- `A positive only in one model family` still matters, but should be promoted as model-sensitive phenotype evidence rather than substrate-general necessity

## What Would Count as Strong Paper-Ready Evidence

The main-text result should ideally look like this:

- in Condition A, the agent lineage invents a stable external protocol and reaches cumulative progress
- in Conditions B and C, the same mission degenerates into repeated work
- the contrast is visible in one simple table or one cumulative-coverage plot
- the negative cells are retained with explicit failure labels rather than summarized away

That would allow the paper to make a crisp claim:

> external memory is not a convenience layer; it is part of the organism-environment coupling that enables cross-generational cognition.

## Analysis Plan

For each generation, record:

- generation index
- session ids
- files checked this generation
- cumulative unique files ruled out
- whether predecessor artifacts were read
- whether new artifacts were written
- whether the run ended by budget exhaustion, safety timeout, or natural exit

Derived plots/tables:

- cumulative unique coverage by generation
- duplicate work ratio by condition
- first appearance of stable schema
- completion generation by condition

## Risks

### 1. Constraint Bypass

Agents may use shell loops, `find`, or other batch shortcuts to trivialize the survey.

Mitigation:
- tighten prompt wording
- enforce one-file-at-a-time inspection in analysis
- fail runs that solve the task by prohibited batch compression

### 2. Exec Confound

Agents may still try to use `exec` as an escape hatch if the capability remains available.

Mitigation:
- remove `exec` from the runtime for this experiment
- reject any hallucinated `exec` call at dispatch level
- analyze only persisted filesystem effects across externally defined generations

### 3. Overfitting to a Known Report Path

If the prompt specifies an exact output artifact like `survey_report.json`, agents may optimize around that path rather than inventing a general legacy protocol.

Mitigation:
- separate final deliverable path from checkpoint path
- treat spontaneously invented intermediate artifacts as the primary signal

### 4. Paper Overlap with 4.1

This experiment can start to look like a survival experiment if framed incorrectly.

Mitigation:
- keep the paper emphasis on durable memory ablation
- report cumulative lineage behavior, not just whether a single organism survives longer

## Open Questions After First Runnable Pass

1. Is the current read-only condition strong enough, or does it need a harder sandbox boundary?
2. Should the first paper-ready sweep keep the prompt fully condition-agnostic, or explicitly state the habitat difference?
3. How should the prompt remind each generation to inspect the whole habitat before starting world queries without hinting a specific handoff protocol?
4. After the `gpt-5.4-pool` matrix is stable, is one Claude `A/B/C` contrast sufficient for the conference paper? Current default: yes; anything broader belongs to later model-sensitivity work.

## Lessons From The First Runnable Pass

The runner implementation clarified three kinds of lessons that should guide later lineage work.

### Design Lessons

- **Negative search is a good cross-ecology task** because the cumulative object is an exclusion frontier rather than a retrieval index; this keeps durable memory central rather than auxiliary.
- **Do not name the checkpoint artifact in the prompt** if the real object of study is spontaneous protocol birth; otherwise the experiment drifts toward compliance with a hinted schema.
- **Keep the prompt condition-agnostic when possible** and let the runner enforce the habitat difference, so the causal variable remains environmental rather than rhetorical.

### Execution Lessons

- **External mortality must be runner-enforced** when the claim is about cross-generational continuity; prompt-level mortality language is not enough because `exec` can otherwise blur the generation boundary.
- **Condition fidelity matters more than elegance**: a pragmatic read-only approximation is acceptable for a first runnable pass, but the approximation must be named explicitly rather than smuggled in as if it were perfect.
- **One top-level runid should own the whole lineage**, with per-generation subtrees underneath it, so both cumulative state and within-generation confounds remain recoverable.

### Analysis Lessons

- **Durable-memory claims need two views of evidence**:
  - observed behavior such as world queries, duplicate work, and predecessor-artifact access
  - persisted state such as surviving artifacts, stable schema, and on-disk frontier growth
- **Per-generation snapshots are worth keeping** because later interpretation depends on what the environment looked like after each externally defined generation, not only on final state.
- **Methodological failures are first-class results**. If the read-only boundary is bypassed or the timeout window is too short to observe useful behavior, that should be recorded as evidence about the harness, not silently normalized.

## Intended Paper Use

If successful, this experiment should replace Ariadne as the main backbone of ALIFE conference section 4.2.

Recommended section title:

**External Memory and Legacy Protocols**

Ariadne can remain as an early precursor or appendix support case, while redesigned `2.6-legacy-protocol-ablation` is now a viable main-text candidate for the section `4.2` causal-memory contrast.

Current honest paper boundary:

- use redesigned `A/B/C` as bounded evidence that writable persistent habitat changes cross-generational behavior in this ecology
- do not yet claim full-task completion or long-run monotonic coverage stability from `2.6` alone

---

## Result Curation

**Current evidence set**

- Experiment is currently `active`.
- Older Kimi and pre-redesign pooled runs remain useful as design history, but they are no longer the primary evidence set.
- The primary evidence set is now the redesigned `20260322` suite under `gpt-5.4`, `minimal`, `8` turns, and `240s` safety timeout.
- Primary positive runs:
  - `20260322-121352-gpt-5.4-pool-condA-8turns`
  - `20260322-121940-gpt-5.4-pool-condA-8turns`
  - `20260322-122549-gpt-5.4-pool-condA-8turns`
  - `20260322-123831-gpt-5.4-pool-condA-8turns`
- Primary contrast runs:
  - `20260322-123123-gpt-5.4-pool-condB-8turns`
  - `20260322-200003-gpt-5.4-pool-condB-8turns`
  - `20260322-123123-gpt-5.4-pool-condC-8turns`
  - `20260322-200456-gpt-5.4-pool-condC-8turns`
- Cross-family calibration runs worth keeping:
  - `20260322-201300-claude-sonnet-4-5-condA-8turns`
  - `20260322-201457-claude-sonnet-4-5-condB-8turns`
  - `20260322-201732-claude-sonnet-4-5-condC-8turns`
- Lower-bound calibration runs worth keeping:
  - `20260322-201948-kimi-oauth-condA-8turns`
  - `20260322-202251-kimi-oauth-condB-8turns`
  - `20260322-202434-kimi-oauth-condC-8turns`
- Methodological side evidence worth keeping but not counting inside the clean matrix:
  - `20260322-194911-gpt-5.4-pool-condB-8turns`, which preserves the no-inheritance outcome but introduces `G2` multi-read shell compression and should remain excluded from the main paper table
  - `20260322-201612-claude-sonnet-4-5-condC-8turns`, which preserves the no-birth read-only outcome but introduces `G2` multi-read shell compression and should remain excluded from the clean cross-family calibration set
- Sidecar intervention evidence that remains worth keeping:
  - `20260321-133505-gpt-5.4-pool-condA-8turns` as the cleanest curated-handoff slice
- Exploratory prompt-control evidence that remains worth keeping:
  - `evidence-set:nolineage-prompt-control-2026-03-22`
  - `20260322-131401-gpt-5.4-pool-condA-8turns` under `nolineage`, where condition `A` produced no durable artifact birth or successor reuse across `2` generations
  - `20260322-131655-gpt-5.4-pool-condA-8turns` under `nolineage`, replicating the same bounded no-birth / no-reuse outcome across `2` generations
- Aggregate read across the redesigned evidence set:
  - `A` now gives `4/4` short positive slices for spontaneous artifact birth in `G1`
  - every redesigned positive `A` slice shows successor adoption
  - the latest `A` slice extends adoption through `G3`
  - `B` now has `2/2` clean negative slices with birth but no successor reuse
  - `C` now has `2/2` clean negative slices with no durable artifact birth
  - a first Claude calibration sweep reproduces the same qualitative `A/B/C` order, though with a somewhat different style: faster runs, one explicit successor-read in `A`, and more aggressive shell behavior in the first `C` attempt
  - a first Kimi calibration sweep also reproduces the qualitative `A/B/C` order, but only reaches a weaker boundary regime: artifact birth in `A` is delayed until `G2`, explicit predecessor reads remain absent in the analyzer summary, and duplicate work stays high across all conditions
  - the clean paper matrix remains free of multi-read shell-call confounds, while one extra `B` rerun is explicitly retained as confounded side evidence rather than silently discarded
- Compact causal comparison:

| Condition | Birth signal | Reuse signal | Failure mode |
|-----------|--------------|--------------|--------------|
| `A` persistent | spontaneous `lineage_state` notes appear immediately | later generations read predecessor notes and extend them | bounded only by short horizon, not by habitat loss |
| `B` wiped | `G1` can still externalize | `G2` starts from empty habitat and cannot inherit | memory object is erased between generations |
| `C` read-only | externalization pressure remains | no durable artifact is available to successors | write attempts fail with `Permission denied` |

**Notable discoveries**

- The earlier "minimal `A` is real but unstable" read was partly an ecology defect, not just a prompt-minimality defect.
- Mixed identifier families materially reduced easy whole-library induction from tiny samples.
- A neutral writable `lineage_state/` surface is enough to induce spontaneous lineage notes without naming a checkpoint schema in the prompt.
- The `A/B/C` split is now cleaner:
  - `A` shows birth plus reuse
  - `B` shows birth without inheritance
  - `C` shows attempted or desired externalization blocked at birth
- A single `nolineage` prompt-control slice under condition `A` produced no durable artifact birth or reuse, which is worth retaining as exploratory prompt-ablation evidence but not yet as a stable claim boundary.
- Analyzer hygiene matters here: shell-read counts must look inside nested tape JSONL and must flag multi-read shell compression explicitly.
- The per-run `quine` rebuild is now part of the method, not an implementation detail, because stale binaries can silently corrupt condition fidelity.
- The remaining limitation is scale, not existence:
  - duplicate-read ratio is still about `0.50` in the short positive `A` slices
  - the lineage has not yet completed the survey
  - current notes are cumulative enough to show handoff, but not yet strong enough to claim monotonic long-run frontier expansion

**Do not collapse into the main claim**

- The redesigned suite is now stronger than exploratory anecdote: it is bounded causal evidence for durable-memory necessity in this ecology.
- For paper purposes, the result currently supports a bounded statement:
  - writable persistent habitat changes lineage behavior in redesigned `2.6`
  - removing persistence (`B`) or writeability (`C`) removes the same inheritance path
  - `2.6` is now usable for a bounded `4.2` replacement or upgrade
  - the next experimental push, if needed, should be `B/C` stabilization plus one cross-family calibration, not a return to the pre-redesign ecology or a wide model benchmark
