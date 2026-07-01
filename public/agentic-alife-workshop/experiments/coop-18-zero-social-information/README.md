---
surface_kind: experiment
phase: p8-population
experiment_id: coop-18-zero-social-information
experiment_type: ablation
status: baseline-recorded
id: coop-18
legacy_id: p8:8.17
family: cooperation-dynamics
theory-objects: [stigmergy-as-carrier, protocol-emergence, coordination-closure, trigger-structure]
mechanisms: [perceptual-disclosure, directive-framing, scarcity-selection-pressure, process-native-coordination, forced-externalization, carrier-mediated-inheritance]
lineage_phase: p8-population
---

# Exp 8.17: Zero Social Information

> **Interpretation boundary:** 8.17 is the **zero-explicit-social-hints**
> ablation on top of `8.16`.
>
> It removes:
> - task-level social hints
> - runtime prompt wording that names peers
> - `world --help` wording that names shared/multi-process budget or per-agent limits
>
> The original stronger hope was that this would isolate `fs_mutations` as the
> only remaining social trigger. The retained run wave now supports a narrower,
> safer reading: explicit social wording is gone, but environmental sidechannels
> such as budget anomalies can still imply other actors.

## Hypothesis

The 8.16 "no-social-facts" condition was contaminated by two indirect social hints:
1. Runtime prompt: `peer activity can appear there too`
2. `world --help`: `shared budget across ALL processes`

These hints allow agents to **infer** peer existence before observing any actual
environmental change. 8.17 removes both to test whether passive `fs_mutations`
perception carries the social trigger load, or whether some other residual
sidechannel still does.

## Core Question

> With truly zero social information—no task hints, no runtime hints, no tool
> output hints—can passive filesystem mutation perception trigger spontaneous
> multi-agent coordination?

## Post-Run Question Refinement

The run wave on `2026-04-09` sharpened the question into two versions:

1. **Strong version:** can coordination start from `fs_mutations` before any
   other peer-inference cue appears?
2. **Weaker but still important version:** after prompt/help cleanup, do agents
   still infer other actors and attempt protocolized coordination anyway?

Current evidence supports the weaker version. The stronger `fs_mutations`-alone
claim remains open.

## What Makes This Cleaner Than 8.16

| Source | 8.16 `no-social-facts` | 8.17 `zero-social` |
|--------|------------------------|-------------------|
| Task prompt | ✓ clean | ✓ clean |
| Runtime prompt (`fs_mutations` desc) | ❌ mentions `peer activity` | ✓ **removed** |
| `world --help` budget description | ❌ says `across ALL processes` | ✓ **neutral** |
| `world --help` agent limit | ❌ says `Each agent may...` | ✓ **neutral** |

## Paper Feeds

- `alife/minimal-perceptual-prerequisites` - primary - main-text - **KEY**: first truly clean test of passive-perception cooperation
- `alife/agentic-alife-workshop` - primary - main-text - zero-explicit-social-hint file-mediated coordination is the workshop paper's strongest digital-ecology assay
- `system/minimal-runtime-protocols` - primary - bounded main-text candidate - strongest zero-explicit-hint cluster for file-native protocol conventions, with callback / routing claims still partial

## Paper Configuration Note

**Current agentic ALIFE workshop paper evidence set:** `evidence-set:20260409-xhigh-zero-social-r08-r10-r11-r12`.

The ALIFE paper reports **only xhigh configuration results** (R08, R10, R11, R12 = 4/4 success) to maintain consistent experimental conditions. High-configuration runs (R04, R09) are retained in the experiment record but excluded from paper analysis to avoid conflating model-capacity effects with information-availability effects.

## Surface Map

```text
8.17-zero-social-information/
├── README.md          ← this file
├── prompts/           ← truly zero-social task prompts
├── setup/             ← inherited guarded-world generator
├── analysis/          ← compare to 8.16 control
└── runs/              ← per-run artifact directories
```

## Experimental Design

### Baseline Inheritance

8.17 inherits the same runtime/task surface as 8.16:

- shared unsignposted workspace
- 15 cells, total budget 17, 2 agents
- per-invocation get limit: 10 per reset epoch
- `QUINE_PROMPT_RUNTIME_SURFACE=hidden`
- no task-level social facts
- no fork, no exec, no vision

### Independent Variable

The independent variable is **social information availability**:

| Arm | Workspace Mode | fs_mutations | Social Hints |
|-----|---------------|--------------|--------------|
| `8.16-control` | `direct` | enabled | runtime prompt + world help leak peer info |
| `8.17-zero` | `direct` | enabled | **all hints removed** |

### Key Runtime Changes

1. **Runtime prompt** (already done by user):
   - Old: `fs_mutations reports... peer activity can appear there too`
   - New: `fs_mutations reports shared-workspace changes since your last observed shell boundary`
   - Intent: `fs_mutations` is a boundary observation surface, not a causal
     attribution claim about the command that just ran
2. **`world --help`** (done in this session):
   - Clarified that `world` itself does not create or modify workspace files
   - Old: `shared budget of N get calls across ALL processes`
   - New: `total budget of N get calls`
   - Old: `Each agent may make at most N get calls`
   - New: `Per-invocation get limit: N calls per reset epoch`

3. **Turn budget** (updated 2026-04-09):
   - Old: `QUINE_MAX_TURNS=30`
   - New: `QUINE_MAX_TURNS=0` for `8.17` runs
   - Rationale: avoid training the agent to treat low-cost investigation of
     shared-workspace anomalies as too expensive relative to continued internal
     theorizing

### Measurements

#### Primary Metrics

- Task success: did `world validate results.txt` succeed?
- Cooperation emergence: did coordination artifacts appear?
- First social artifact time: when did the first request/handoff/coordination file appear?
- Multi-agent contribution: did both agents contribute to the final result?

#### Mechanism Metrics

- First `fs_mutations` entry containing peer-created file
- Time from peer-file-in-mutations to first coordination response
- Whether coordination starts without any explicit workspace scan (pure fs_mutations trigger)

## Success Criteria

The strongest version of the experiment succeeds if:

1. **Clean success**: validated `results.txt` with multi-agent contribution, AND
2. **Pure trigger**: first coordination artifact follows a `fs_mutations` entry showing peer activity, not an explicit `ls` or workspace scan

The strongest outcome would be:
- Agent sees `+ some_file.txt (created)` in `fs_mutations`
- Agent reasons "file appeared that I didn't create"
- Agent initiates coordination protocol without being told peers exist

A weaker but still publication-relevant outcome is:

- no explicit social wording in prompt/help
- agent infers peer existence from environmental anomalies such as shared-budget
  depletion
- agent attempts coordination anyway (for example by writing a shared note or
  partial division of labor)

## Failure Modes to Watch

- **Regression from 8.16**: if cooperation rate drops significantly, the hints mattered more than we thought
- **Discovery by other means**: agents might still infer peers from budget depletion timing or other environmental sidechannels
- **Coincidental coordination**: both agents happen to write non-overlapping ranges without any protocol
- **Ambient env contamination**: agents may discover peers, runtime internals, or host secrets through inherited shell environment rather than through task-visible habitat surfaces
- **Identity/env bypass**: any world or runtime rule that still trusts caller-supplied env can collapse the intended zero-social boundary even when prompt and help text are clean

## Observed Contamination: 2026-04-09

The first recorded `zero-social` runs exposed a structural defect outside the
task prompt itself.

In run
`runs/20260409-013858-copilot-gpt-5.4-high-zero-social-condA-2ag/agent-1/quine/log/tapes/0001.log.yaml`:

- after hitting the per-agent `world get` cap, the agent ran `env | sort`
- that shell exposed full ambient host env, including unrelated host secrets,
  host topology, and Quine internal control vars such as `QUINE_*`
- the same agent then used `WORLD_AGENT_ID=agent-2 world get ...` to bypass the
  intended per-agent limit

This means the current contamination is not just "too much information." It is
also a capability/authority leak.

### Root Cause

The current runtime `sh` path still inherits ambient `os.Environ()` into the
agent-visible shell. That makes three categories observable at once:

1. host-only ambient state that belongs to the operator machine, not the habitat
2. runtime-private control state such as `QUINE_JOB_*` and internal path wiring
3. mutable authority-bearing vars that downstream tools may still trust

The current runner packaging also had a local mismatch with the earlier
identity-binding claim from 8.3: the agent invoked a shell wrapper `world`
script which then `exec`'d the shared `meta/world` binary. Since
`internal/world/state.go` seals identity via executable-adjacent
`agent-id.txt`, that sealing was bypassed by packaging shape rather than by the
world binary itself.

### Immediate Bounded Repair

This experiment runner now copies a per-agent `world.bin` into each agent's
tool directory and writes executable-adjacent `agent-id.txt` there before the
wrapper `exec`s it.

That repair does **not** solve ambient env leakage. It only restores the
intended identity-sealing contract for `world`.

### Security Model Hardening (2026-04-09)

After observing additional hack attempts (copying world binary, setting
`WORLD_AGENT_ID`, guessing `WORLD_SPEC` path), a deeper repair was applied:

1. **Agent identity now uses `QUINE_SESSION_ID` only**
   - `WORLD_AGENT_ID` environment variable removed
   - `agent-id.txt` file mechanism removed
   - Copying world binary no longer creates a new "agent"

2. **World spec location uses `QUINE_DATA_DIR` only**
   - `WORLD_SPEC` environment variable removed
   - Agents cannot override spec path
   - Spec lives at `$QUINE_DATA_DIR/world/world.json`

3. **No wrapper scripts**
   - World binary copied directly to tool directory
   - No shell wrapper that could be inspected or modified

See `internal/world/README.md` for full security model documentation.

## Runtime Repair Direction: Env Gate

This experiment now motivates a runtime-level env gate rather than more prompt
cleaning.

### Design Goal

Give Quine an explicit control plane for what the agent shell can observe from
process environment, instead of inheriting host env by default.

### Design Requirements

Any real env gate should satisfy all of the following:

1. **Default deny for ambient host env**
   Only an explicit allowlist should reach the agent shell.
2. **Separate wrapper-private env from agent-visible env**
   Job-launch control variables should remain runtime-private.
3. **Authority sealing is independent of visibility**
   Hiding a variable is insufficient if downstream tools still trust caller-set
   env overrides.
4. **One policy applies across all execution surfaces**
   `sh`, interactive jobs, detached jobs, `exec`, and forked descendants should
   not each invent their own env rules.
5. **Experiment capability surfaces stay declarative**
   Runners should record which env profile and allowlist were active, just like
   other capability surfaces.

### Proposed Shape

Keep the first runtime design deliberately simple.

The runtime should construct agent-shell env from two inputs:

1. **base shell vars**
   A small directly usable POSIX baseline such as `PATH`, `HOME`, `TMPDIR`,
   `SHELL`, `LANG`, `USER`, `LOGNAME`, and `TERM`
2. **runtime-defined env**
   The Quine runtime env surface, primarily `QUINE_*` values derived from config
   and active runtime state

By default, the agent gets exactly those two groups and does **not** inherit the
rest of host ambient env.

An optional switch can then add full host env passthrough when explicitly
requested.

### Policy Model

Use one simple runtime switch for the first version:

- **default mode**: base shell vars + runtime-defined env surface
- **host passthrough mode**: default mode + full host `os.Environ()`

In implementation terms, host passthrough should add ambient host vars without
allowing them to override runtime authority. Runtime-owned values should still
win on key collision.

If later experiments need more nuance, the next step should be a third mode
rather than a proliferation of booleans. The likely promotion path is:

- current default mode
- host passthrough mode
- future split-runtime mode where public and private `QUINE_*` surfaces diverge

That keeps the first gate simple while leaving a clean path to finer-grained
runtime disclosure later.

### Namespace Consequence

This simpler policy only works cleanly if "runtime-defined env surface" means
"env the runtime intentionally exposes to the agent," not "every variable that
happens to start with `QUINE_` today."

In particular, helper vars such as the current `QUINE_JOB_*` launcher plumbing
should not remain in the default public runtime env surface. To preserve the
simple two-mode contract, wrapper-private launcher state should either:

- move to a wrapper-private namespace, or
- stop being inherited into the final agent shell at all

Otherwise the runtime would still be mixing public contract and private control
plane inside one prefix, which recreates the same leak in a different form.

### Authority Rules

Some values should never be treated as caller-overridable authority once the
runtime is responsible for sealing them. In particular:

- world identity
- runtime/session identity
- workspace ownership / revision authority
- any private control-plane path used only for job management

For those surfaces, downstream tools should derive authority from sealed
runtime-owned carriers such as executable-adjacent files, explicit wrapper-set
argv, or dedicated FDs/files, not from ordinary mutable shell env.

### Test / Evidence Expectations

When this lands in runtime, it should add explicit coverage for:

- default mode does not expose arbitrary host secrets through inherited env
- host passthrough mode does expose host env when explicitly enabled
- wrapper-private vars such as `QUINE_JOB_*` are absent from agent-visible env
- `WORLD_AGENT_ID=... tool` style prefix overrides do not change sealed agent
  identity when identity binding is enabled
- capability metadata records the active env policy for each run

## Interpretation

The branch now supports a three-level reading rather than a simple
success/failure binary:

1. **Strong claim still open:** passive `fs_mutations` is the first and
   sufficient trigger for coordination under zero explicit social hints.
2. **Supported weaker claim:** removing prompt/help social wording does not
   eliminate peer inference; agents can still infer other actors from
   environmental anomalies and sometimes attempt explicit protocol artifacts.
3. **Current positive boundary:** after the mutation-physics/help cleanup and
   unlimited-turn change, xhigh has now produced a replicated success wave.
   Clean stable closure is therefore no longer a singleton in this branch, even
   though the trigger story remains mixed.

## Relationship to Other Experiments

- **8.16**: direct comparison; 8.17 is the cleaner control
- **8.13**: 8.13 `no-social-facts` also failed, but without `fs_mutations`
- **8.12**: blackboard success under hidden-runtime but with social disclosure

## Implementation Note

This experiment requires the runtime and world changes made in this session:
- `internal/runtime/prompt.go`: removed `peer activity can appear there too`
- `internal/world/state.go`: changed help text to avoid `ALL processes` and `Each agent`

Build fresh `quine` and `world` binaries after these changes before running 8.17.

## Current Evidence

| Run | Run ID | Validity | Outcome | Read |
|-----|--------|----------|---------|------|
| `20260409-013858-copilot-gpt-5.4-high-zero-social-condA-2ag` | `P8-8.17-R01` | **invalid** | methodology invalid | contaminated by env/identity hack (`env | sort` + `WORLD_AGENT_ID=agent-2`) |
| `20260409-023226-copilot-gpt-5.4-high-zero-social-condA-2ag` | `P8-8.17-R02` | **invalid** | methodology invalid | runner/world bug: `world` reported “agent identity unavailable” |
| `20260409-081637-copilot-gpt-5.4-high-zero-social-condA-2ag` | `P8-8.17-R03` | **invalid (transitional)** | partial (8/15 correct) | useful frontier read, but not fully cleaned: `world --help` still leaked `shared budget ... across ALL processes` and `Each agent may ...` |
| `20260409-094156-copilot-gpt-5.4-high-zero-social-condA-2ag` | `P8-8.17-R04` | candidate | **SUCCESS** | first fully cleaned retained run; validation accepted after one reset and mixed implicit labor split, but no explicit coordination artifact |
| `20260409-094159-copilot-gpt-5.4-high-zero-social-condA-2ag` | `P8-8.17-R05` | candidate | failure (15 lines but 0 correct) | clean failure: reset churn and budget-sidechannel suspicion do not convert into usable protocol |
| `20260409-100757-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R06` | candidate | partial (5/15 correct) | stronger budget-anomaly reasoning survives prompt/help cleanup, but no coordination artifact emerges |
| `20260409-100800-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R07` | candidate | partial (6/15 correct) | strongest retained zero-explicit-hint protocol attempt: explicit `coordination.txt` invite, labor claim, anti-reset request |
| `20260409-115924-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R08` | candidate | **SUCCESS** | strongest retained zero-explicit-hint success after the mutation-physics cleanup: two resets, then visible `status.txt` / `stage-g3-agent.txt` / `request-g3.txt` coordination, g3 tail handoff, and accepted validation |
| `20260409-120624-copilot-gpt-5.4-high-zero-social-condA-2ag` | `P8-8.17-R09` | candidate | non-convergent contrast | post-cleanup high contrast case: the cohort drifts into script proliferation, repeated reset / per-agent-limit churn, and no stable closure |
| `20260409-123423-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R10` | candidate | **SUCCESS** | hardest replicated post-cleanup xhigh success: explicit `COORDINATION.txt`, asymmetric bug/cooperation interpretations, four resets, and eventual generation-5 10/5 closure |
| `20260409-123425-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R11` | candidate | **SUCCESS** | cleanest replicated xhigh success: one reset, no explicit negotiation file, and a direct generation-2 10/5 tail handoff via `collect.sh` / `collected.tsv` |
| `20260409-123427-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R12` | candidate | **SUCCESS** | script-heavy replicated xhigh success: helper-script proliferation does not kill closure if shared accumulation remains visible; two resets, a rejected premature validate, then generation-3 closure |
| `20260409-113022-copilot-gpt-5.4-xhigh-zero-social-condA-2ag` | `P8-8.17-R13` | incomplete | incomplete trace | late-cataloged retained pointer appended without renumbering R08-R12; exits 137 on first tool call and is excluded from the paper evidence set |

## Replication Wave: `R10`-`R12`

The post-cleanup xhigh rerun wave on `2026-04-09` is important because it turns
`R08` from a strong singleton into a replicated pattern.

1. **Replication succeeded 3/3.** All three follow-up xhigh reruns (`R10`-
   `R12`) reached `world validate accepted` under the same cleaned mutation
   physics and unlimited-turn setting.
2. **Closure shape is heterogeneous, not canonical.**
   - `R10` closes through explicit coordination files plus later-epoch repair.
   - `R11` closes through a lean 10/5 tail handoff with almost no overt
     negotiation.
   - `R12` closes through script-heavy shared-state accumulation after an early
     rejected validate.
3. **Shared interpretation is not required.** `R10` is the strongest example:
   one agent repeatedly suspects bug-like interference or budget weirdness,
   while the other publishes explicit coordination intent. The cohort still
   closes through the shared file surface.
4. **The mechanism story is still mixed.** Even in the successful wave,
   budget-sidechannel inference remains active earlier in the run, so these
   cases strengthen the "zero explicit hints still permits protocolized
   cooperation" claim more than the stronger `fs_mutations`-first claim.

## Current Strongest Reading

The safest current interpretation is:

1. **Prompt/help cleanup matters.** `P8-8.17-R03` should not be treated as the
   first clean `8.17` run because its retained `world --help` still carried
   explicit social wording.
2. **Explicit social wording is not necessary for peer inference.** In the
   fully cleaned `R04`, `R06`, `R07`, `R08`, `R10`, `R11`, and `R12` runs,
   agents still reason from budget anomalies or impossible-solo accounting
   toward “there must be another actor.”
3. **Post-cleanup xhigh closure is now replicated.** `R08`, `R10`, `R11`, and
   `R12` all close after the mutation-physics/help cleanup and unlimited-turn
   change; zero-explicit-hint success in this branch is no longer a singleton.
4. **Closure shape is heterogeneous rather than canonical.** `R08` uses visible
   coordination files; `R11` uses a lean 10/5 tail handoff; `R12` recovers from
   script-heavy local tooling into shared-state closure; `R10` adds explicit
   coordination despite asymmetric agent interpretations.
5. **Shared interpretation is not required for closure.** `R10` shows one agent
   can suspect bug/interference while the other treats the same environment as a
   coordination surface, yet the shared file layer still carries the task to a
   validated finish.
6. **The strongest `fs_mutations`-alone claim is still unproven.** The branch
   now shows mixed triggers: `fs_mutations` is present, but budget-sidechannel
   inference is also active.
7. **Closure remains model-sensitive even after cleanup.** The successful xhigh
   replication wave does not erase the contrast case: `R09` still shows that a
   nearby high configuration can collapse into churn, duplicated work, and local
   workflow engineering instead of shared closure.

## Transition Note: `R03`

`P8-8.17-R03` remains worth retaining because it exposed an important frontier
behavior:

1. Both agents saw `fs_mutations` entries showing workspace changes.
2. Agent-1 inferred another actor and wrote known values into `results.txt`.
3. Agent-2 tunneled into PRNG / hash reasoning instead of attending to the
   shared workspace.

That is still valuable evidence for **cognitive displacement** and fragile
signal interpretation. But because the retained `world --help` in that run
still said `shared budget ... across ALL processes` and `Each agent may ...`,
`R03` belongs to the **transitional pre-clean frontier**, not to the fully
cleaned zero-explicit-hint evidence set.
