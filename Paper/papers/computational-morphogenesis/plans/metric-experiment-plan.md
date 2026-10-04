# Metric And Experiment Plan: Successor Morphospace

## Purpose

This plan turns the computational-morphogenesis route from metaphor into an
operational assay. The paper should not claim morphogenesis from the current
seed ladder alone. It needs measured successor forms under controlled pressure
and resource manipulations.

Target claim after this plan matures:

> Holding the process-native substrate fixed, Quine successor morphology varies
> with self-description resource and future-carrier pressure, and the resulting
> forms can be measured by executable organization and post-handoff viability.

## Unit Of Observation

One retained run should produce a morphology record:

```text
M = <B, S, C, T, H, A, R, V>
```

| Field | Meaning | Minimal observable |
|-------|---------|--------------------|
| `B` | body class | generated files, build target, binary / script kind |
| `S` | self-description used | contract-only, source projection, build metadata, tool schema |
| `C` | carrier restored | stdin, fd, `ctl/wake`, socket / queue adapter |
| `T` | tool surface | none, shell, filesystem, compiler, runtime tools |
| `H` | handoff relation | spawn, wrapper restart, real `exec`, `exec` plus ready state |
| `A` | addressability | none, immediate, delayed, repeated |
| `R` | readiness recovery | exit, one-shot, `idle`, multi-wake ready |
| `V` | viability proof | launch, echo, unseen task, tool proof, commit-reveal / multi-wake |

A run is not a morphogenesis datapoint unless the record separates:

- resources available before handoff;
- pressure / carrier demand imposed before handoff;
- successor organization actually constructed;
- post-handoff viability outcome.

## Metric Families

| Metric family | Candidate measures | Required artifact |
|---------------|--------------------|-------------------|
| code structure | file count, generated LOC, AST nodes/depth, call graph size, binary size | generated successor tree or copied source projection |
| interface surface | exposed stdin/stdout, inherited fd, `ctl/wake`, tool API, waiting state | fd map, control-surface inventory, prompt/schema snapshot |
| carrier topology | carrier preserved, wrapped, rebuilt, adapted, or absent | handoff wrapper, carrier log, delivery trace |
| runtime behavior | syscall/process graph, file writes, tool invocations, build commands | trace log, tool transcript, process observations |
| developmental effort | build attempts, LLM calls, tool-use count, wall time, repair loops | pre-handoff transcript and tool log |
| viability depth | delayed task success, proof artifact, return-to-ready, second wake success | post-handoff challenge log |

Minimum first implementation can be manual coding against retained artifacts.
Automated extraction is future work; do not block the first assay on full metric
automation.



## Numbered Experiment Branch

The metric plan is now routed through the prospective experiment branch
[`experiments/active/autopoiesis-succession`](../../../../experiments/active/autopoiesis-succession/README.md).
Use `6M.xx` IDs for morphospace assays instead of extending the legacy `6.xx`
sequence or duplicating `6F.xx` future-addressability ownership.

| ID | Plan phase | Role |
|----|------------|------|
| [`6M.01`](../../../../experiments/active/autopoiesis-succession/succ-29-retained-run-morphology-inventory/README.md) | retained-run inventory and manual morphology coding | first required extraction pass |
| [`6M.02`](../../../../experiments/active/autopoiesis-succession/succ-30-fixed-substrate-carrier-ladder/README.md) | fixed-substrate carrier ladder | first new controlled batch |
| [`6M.03`](../../../../experiments/active/autopoiesis-succession/succ-31-self-description-resource-ladder/README.md) | self-description resource ladder | separates resource from pressure |
| [`6M.04`](../../../../experiments/active/autopoiesis-succession/succ-32-niche-pruning-pressure-ladder/README.md) | niche-pruning pressure ladder | morphology collapse / pruning trajectory |
| [`6M.05`](../../../../experiments/active/autopoiesis-succession/succ-33-resource-degradation-compensation/README.md) | resource-degradation compensation | developmental-effort and failure-mode axis |
| [`6M.06`](../../../../experiments/active/autopoiesis-succession/succ-34-viability-depth-overlay/README.md) | viability-depth overlay | shared-provenance morphology readout for repeated wake / readiness depth |

Execution order as originally planned:

```text
6M.01 -> 6M.02 -> (6M.03 or 6M.04) -> 6M.05 / 6M.06
```

Current execution state: the first retained wave through `6M.06` now exists.
`6M.01` supplied the seed inventory, `6M.02` validated the fixed-substrate
carrier harness, `6M.03` added the self-description material axis, `6M.04`
added pressure pruning, `6M.05` added resource degradation, and `6M.06` added
repeated-readiness depth. The plan's remaining role is not to block later
experiments, but to constrain paper-facing synthesis: keep morphology class,
implementation thickness, and viability depth separate, and treat generated
LOC / effort / latency metrics as exploratory until cleaned.

## Concrete Metric Pipeline

Use the first metric pass as an artifact-coding exercise before launching a new
benchmark. The goal is to make each retained run legible as a morphology record
and to discover which observations are already recoverable from artifacts.

### Phase 0: Retained-Run Inventory

Create a seed inventory of existing retained runs before running new
conditions. Each row should name the public condition, internal run / experiment
handle, artifact root, and whether the trace currently exposes enough material
to code `M = <B,S,C,T,H,A,R,V>`.

Initial inventory targets:

| Public condition | Expected source | First coding question |
|------------------|-----------------|-----------------------|
| contract-only source-absent lane | source-absent cleanroom run | can the contract-level successor be classified without source/body copy? |
| source-visible raw stream | raw-stream self-source run | is the successor a stream-compatible custom body rather than rebuilt runtime? |
| source-visible fd / adapter, if retained | fd-lane runs | what carrier topology was restored or adapted? |
| source-visible native `ctl/wake` | native control-surface run | can runtime rebuild, delayed wake, tool proof, and `idle` be independently coded? |
| minimized native `ctl/wake` | prompt-minimized control run | does the same morphology class survive less strategy narration? |
| public `ctl` plus raw carrier | control-present raw-carrier control | is visible control surface unselected when not load-bearing? |
| generic `exec` and direct rebuild controls | baseline controls | where do launch, immediate re-entry, and delayed closure separate? |

Output of this phase can be a manual CSV or table. Do not infer missing fields;
mark them `unknown` and use the gap to plan instrumentation.

### Phase 1: Manual Morphology Coding

Code each retained run with a controlled vocabulary plus a small number of
numeric companion measures. The first pass should be reproducible by a second
reader from artifacts, not generated from memory.

Controlled vocabulary for `M`:

| Field | Allowed first-pass values |
|-------|---------------------------|
| `B` body class | `ancestor-body-copy`, `rebuilt-runtime`, `custom-stream-body`, `fd-adapter`, `one-shot-wrapper`, `standard-utility`, `failed-no-successor`, `unknown` |
| `S` self-description used | `none`, `contract-only`, `source-projection`, `source-plus-build-metadata`, `source-plus-control-schema`, `binary-or-artifact`, `unknown` |
| `C` carrier restored | `none`, `stdin`, `fd`, `ctl-wake`, `socket`, `queue`, `unknown` |
| `T` tool surface | `none`, `shell`, `filesystem`, `provider-api`, `compiler`, `runtime-tools`, `unknown` |
| `H` handoff relation | `spawn`, `wrapper-restart`, `scripted-exec`, `agent-constructed-exec`, `exec-plus-ready`, `unknown` |
| `A` addressability | `none`, `launch-only`, `immediate`, `delayed`, `multi-wake`, `unknown` |
| `R` readiness recovery | `exit`, `one-shot`, `idle`, `multi-wake-ready`, `unknown` |
| `V` viability proof | `launch-only`, `echo`, `visible-task-batch`, `unseen-task`, `tool-proof`, `commit-reveal`, `multi-wake-proof`, `failure`, `unknown` |

Numeric companions for the first pass:

| Measure | First-pass definition |
|---------|----------------------|
| generated files | count files generated or materially rewritten for the successor |
| generated LOC | count non-empty, non-comment lines in generated successor code where artifact-backed |
| artifact size | built binary / script size when available |
| build attempts | distinct compile/build invocations before handoff |
| tool calls | tool invocations before and after handoff, counted separately when possible |
| post-handoff tasks | delayed tasks completed after replacement |
| idle returns | observed returns to required waiting state |
| recovery latency | wall-clock or step count from handoff to ready state, if trace exposes it |


### First Inventory Stop Boundary

`6M.01` showed that the seed ladder can be coded categorically, but it should
not yet be promoted to a quantitative morphospace result. The historical
control-surface positives retain strong live-observation summaries but not a
uniform successor body / trace snapshot. Before `6M.02` is run as a result
batch, the runner should emit a per-run morphology row or `meta/morphology.env`
with body class, self-description resource, carrier, handoff, addressability,
readiness, proof, build attempts, generated file/LOC counts, artifact size,
tool-call counts, and ready-state timestamps.

### Phase 2: Fixed-Substrate Carrier Ladder

After the retained-run inventory exposes what can be measured, run the first new
batch by holding the POSIX / Quine substrate and self-description resource as
fixed as possible while varying only the future-work carrier. The key dependent
variable is not pass/fail alone; it is morphology class plus post-handoff
operational success.

Minimum new ladder:

1. raw stdin delayed task;
2. inherited fd or adapter-delivered delayed task;
3. native `ctl/wake` delayed task with tool proof and return to `idle`;
4. public-control-visible raw-stream control.

For each condition, freeze the prompt frame, available resources, and artifact
retention rules before comparing successor classes.

### Phase 3: Readout

The first paper-ready figure should be a condition-by-morphology matrix, not a
claim of broad success rates. A minimal mature readout should show:

```text
condition + resource + carrier -> coded M vector -> post-handoff outcome
```

Promote to a morphospace plot only after there are enough retained records to
make regions or trajectories visible. Until then, call the table a
`trace-level successor-morphology assay`.

The first `6M.01`--`6M.10` wave now meets the minimum table target for a
trace-level assay: it contains stream-compatible successors, fd adapters,
native-control runtime successors, controls, degradation compensation/failure
basins, finite readiness-depth stress positives, and a carrier/capability-
overlap ladder that separates formal reconstruction from operational
continuity. It still does not justify success-rate language or formal
morphospace-region claims.

### Implementation Rules

- Every morphology label must point to an artifact, trace line, or explicit
  `unknown`; no label should be assigned from author memory alone.
- Treat `exec`, delayed delivery, tool use, and return-to-ready as separate
  evidence columns; do not collapse them into one `success` field.
- Keep public condition names primary and internal experiment numbers secondary.
- Use `post-handoff operational success` in metric tables when biological
  `viability` wording would be ambiguous.
- Store future machine-readable records under a paper-local `data/` surface once
  the first manual sheet exists.

## Experiment Matrix

### E1. Fixed-Substrate Carrier Ladder

**Question:** If the POSIX / Quine substrate is held fixed and the future-task
carrier changes, do successor forms fall into distinct measurable classes?

| Condition | Resource | Carrier pressure | Expected morphology class | Viability test |
|-----------|----------|------------------|---------------------------|----------------|
| stream-only | contract or source | raw stdin | stream-compatible successor | delayed stream task after `exec` |
| fd lane | source | inherited fd / adapter | fd-mediated successor | delayed fd-carried task |
| native control | source | `ctl/wake` | rebuilt or control-compatible runtime | delayed wake, tool proof, `idle` |
| public control non-load-bearing | source + visible `ctl` | raw stdin | stream successor, no control rebuild | stream task succeeds; no control closure claimed |

**Paper role:** core empirical morphospace claim.

### E2. Self-Description Resource Ladder

**Question:** How does successor morphology change as self-description resource
moves from runtime contract to source projection to richer build/control schema?

| Condition | Resource available | Held pressure | Expected read |
|-----------|--------------------|---------------|---------------|
| contract-only | runtime contract, no source | repaired stream lane | contract-level reconstitution |
| source-visible | source projection | same stream lane | possible stream body or source-copied adapter |
| source + build metadata | source, build path, toolchain hints | native control | more likely rebuilt runtime |
| source + minimized mission | source, sparse mission | native control | scaffold-light runtime reconstruction |

**Paper role:** separates material for reproduction from pressure to reproduce.

### E3. Niche-Pruning Pressure Ladder

**Question:** Does successor morphology collapse as future-task pressure weakens?

Start from a native-control closure task and remove requirements one at a time:

1. no return-to-idle required;
2. no runtime tools required after handoff;
3. raw-delivered delayed task instead of `ctl/wake`;
4. immediate post-`exec` task instead of delayed task;
5. launch-only success.

Expected outcome is not monotonic proof yet, but a morphology trajectory from
rebuilt runtime toward adapter or one-shot successor.

**Paper role:** phase-transition / morphology plasticity candidate.

### E4. Resource-Degradation Compensation Ladder

**Question:** Under stronger carrier pressure, can Quine compensate for degraded
self-description by spending more developmental effort?

Possible manipulations:

- partial source projection;
- missing build metadata;
- disabled compiler until discovered through tool use;
- truncated runtime contract;
- hidden control schema but visible carrier demand.

Measures:

- additional build attempts;
- generated adapter code;
- tool-use count;
- failure mode when compensation fails.

**Paper role:** supports a developmental-resource reading without claiming
biological development.

### E5. Multi-Wake Viability Depth

**Question:** Is return-to-ready a reusable runtime state rather than a one-shot
completion marker?

Protocol:

1. successor crosses `exec`;
2. receives delayed wake task A;
3. emits proof and returns to `idle`;
4. receives delayed wake task B;
5. emits second proof and returns to `idle` again.

**Paper role:** strengthens viability-depth axis `R` and `V`.

## Minimal First Data Table

For the first morphology paper draft, do not aim for a large benchmark. Aim for
a table that classifies 8--12 retained runs across at least three conditions:

| Run class | Count target | Required fields |
|-----------|--------------|-----------------|
| stream-compatible successors | 3 | `B`, `S`, `C`, `H`, `V` |
| fd / adapter successors | 2 | `B`, `C`, carrier topology, viability outcome |
| native-control runtime successors | 3 | all fields, with delayed task and readiness recovery |
| controls / failures | 2--4 | manipulated condition, failure point, closest morphology |

The key deliverable is a morphospace plot or table, not a success-rate claim.

## Extraction Template

For each run, preserve or derive:

```yaml
run_id:
condition_name:
substrate_held_fixed:
self_description_resource:
future_carrier_pressure:
successor_body_class:
source_relation:
carrier_restored:
tool_surface_after_handoff:
handoff_relation:
addressability_depth:
readiness_recovery:
viability_proof:
developmental_effort:
failure_point:
notes:
```

## Stop Conditions

Do not promote this route to a drafted result paper until:

- at least one carrier ladder has multiple retained examples or a clearly named
  reason why the paper remains trace-level;
- every morphology label maps to observable artifacts rather than metaphor;
- the paper can distinguish launch, immediate handoff, delayed carrier
  addressability, and return-to-ready;
- related-work citation anchors are source-checked against morphology,
  autoconstructive evolution, and digital-organism comparators.

## Immediate Next Work

1. Build a manual morphology-record sheet for existing retained runs.
2. Run or recover a subprocess-vs-`exec` / launch-only control if not already
   artifact-backed.
3. Design E1 carrier ladder as the first new experiment batch.
4. Delay broad DR until metric extraction exposes a specific comparator gap.
