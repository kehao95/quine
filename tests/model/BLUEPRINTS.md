# Model Experiment Blueprints

This document is the planning surface for model-layer experiments that are not
yet active in the canonical registry.

Use it to design the next `L3/L4/L5` evaluations before they are registered in
[`tests/model/evaluations.toml`](./evaluations.toml).

Blueprints are intentionally pre-registry:

- they may describe multiple prompt shapes before convergence
- they may rely on harness pieces that do not exist yet
- they may fail repeatedly before a clean scorer emerges
- they should not be treated as active acceptance until prompt, scorer, and at
  least one preserved passing baseline exist

## Status Model

Each planned experiment should be thought of as one of these:

- `concept`: the mechanism and pressure story are visible, but the harness is not
- `scaffolded`: the environment frame and scorer direction are explicit enough to build
- `pilot-ready`: prompt red lines, environment setup, and success markers are concrete enough for a first live run
- `registerable`: prompt, scorer, and preservation plan are clear enough to promote into the canonical evaluation catalog

`pilot-ready` is the handoff point from blueprint-only design into the runnable
pre-registry pilot surface:

- registry: [`tests/model/pilots.toml`](./pilots.toml)
- prompt path: `tests/model/pilot-<level>-<mode>/<feature>/<variant>/prompt.md`
- runner selector: `./tests/model/run.sh pilot:<id>`

## Blueprint Anatomy

Every experiment blueprint should specify all of the following:

| Field | Meaning |
|------|---------|
| Future id | intended canonical selector once registered |
| Layer | `l3`, `l4`, or `l5` |
| Feature | runtime surface being tested |
| Goal | what the experiment is trying to prove |
| Pressure model | why the intended mechanism becomes selected |
| Environment frame | files, processes, env toggles, and fixtures the harness must set up |
| Prompt red lines | what the prompt must not leak |
| Scorer shape | what evidence would count as success |
| Expected failures | likely ways the first pilot will fail |
| Adjustment knobs | what we can tune between pilot runs without changing the claim |

## Layer-Specific Design Rules

### L3 Usage

- The mechanism may be named directly.
- The prompt may explain protocol and expected operation.
- The scorer should check mechanism use, not just the end result.
- `L3` is the right home for env-controlled behavior that changes operation but
  does not need emergent discovery claims.

### L4 Discovery

- The prompt must not name the target mechanism.
- The prompt must not teach the strategy or mechanism principle.
- Pressure may come from task shape, runtime artifacts, deadlines, or hazards.
- The environment should make the mechanism legible, not hidden behind trivia.
- Pilot plans should expect several rewrites before the hint leakage is honest.

### L5 Necessity

- The prompt must not hint the mechanism directly or indirectly.
- The task must be materially hard to complete without the mechanism.
- The scorer should be able to argue that the mechanism was forced or
  overwhelmingly selected, not merely chosen.
- Each `L5` blueprint should include at least one plausible failure mode where
  the model wins through an unintended workaround; that workaround must then be
  closed by environment design rather than prompt narration.

## Planned Portfolio

### L3 Planned

See [`tests/model/blueprints/l3-usage.md`](./blueprints/l3-usage.md).

| Future id | Feature | Status | Why it exists |
|----------|---------|--------|---------------|
| `shell-envelope-timeout-truncation-cwd` | `shell-envelope` | `scaffolded` | env-controlled `sh` behavior still lacks an explicit model-layer usage contract |
| `response-governance-exit-disabled-impossibility` | `response-governance` | `scaffolded` | exit-disable and impossible-work posture still lack a direct usage evaluation |
| `context-pressure-anchor-exec-carryover` | `context-pressure` | `concept` | memory thresholds and anchor/`context/` exec carryover still have no explicit `L3` operational proof (reworded after the wisdom mechanism deletion, registry-design-brief D5) |
| `tape-lineage-parent-child-incarnation` | `tape-lineage` | `scaffolded` | session / tape / parent lineage are promised surfaces but not yet model-checked |
| `idle-explicit-suspension-resume` | `idle` | `concept` | explicit suspension should prove the model can enter quiescence and resume correctly under named control input |

### L4 Planned

See [`tests/model/blueprints/l4-discovery.md`](./blueprints/l4-discovery.md).

| Future id | Feature | Status | Why it exists |
|----------|---------|--------|---------------|
| `stdin-binary-replay-discovery` | `stdin` | `pilot-ready` | stdin/material continuity still lacks a clean unnamed discovery task |
| `fork-deadline-sharded-search` | `fork` | `pilot-ready` | current fork prompts still teach delegation too directly; current pilot also probes whether retained relation/helper surfaces stay legible after unnamed discovery |
| `fork-adopt-winning-world-promotion` | `fork-adopt` | `pilot-ready` | winner-adoption still leaks mechanism shape |
| `exec-final-utility-stream-handoff` | `exec` | `pilot-ready` | current exec discovery is still too close to protocol teaching |
| `switch-world-rollback-after-destructive-probe` | `switch-world` | `pilot-ready` | destructive recovery still needs a cleaner unnamed prompt |
| `anchor-memory-recall-barrier-ledger` | `anchor-memory` | `pilot-ready` | current memory retrieval tasks still explain too much of the mechanism |
| `workspace-overlay-dangerous-decoder-containment` | `workspace-overlay` | `pilot-ready` | overlay / containment physics need a cleaner discovery frame |
| `idle-external-poke-discovery` | `idle` | `concept` | explicit suspension only matters if a model can later discover it as the clean response to waiting-for-input pressure |

### L5 Planned

See [`tests/model/blueprints/l5-necessity.md`](./blueprints/l5-necessity.md).

| Future id | Feature | Status | Why it exists |
|----------|---------|--------|---------------|
| `detach-overlap-deadline` | `detach` | `concept` | clean overlap pressure can plausibly force detached jobs |
| `fork-parallel-hypothesis-pressure` | `fork` | `concept` | parallel exploration under deadline is a strong `L5` candidate |
| `exec-stream-ownership-forcing` | `exec` | `concept` | direct process-image replacement may be the only honest path in a streaming handoff task |
| `switch-world-clean-final-branch` | `switch-world` | `concept` | revision rewind may become necessary under destructive search pressure |
| `anchor-memory-context-cliff-survival` | `anchor-memory` | `concept` | memory survival under recall gates is a strong necessity candidate |
| `sandbox-hostile-artifact-survival` | `sandbox` | `concept` | unsafe artifacts can make isolation effectively necessary |
| `process-surface-live-peer-rescue` | `process-surface` | `concept` | runtime-surface coordination is a plausible zero-hint necessity task |
| `idle-quiet-standby-pressure` | `idle` | `concept` | some long-lived sessions should only stay viable if they can suspend cleanly instead of burning turns or fabricating progress |
| `agents-md-startup-token-lockin` | `agents-md` | `pilot` | fresh-startup token recall with zero shell budget probes whether startup guidance is truly durable |
| `context-memory-exec-token-lockin` | `context-memory` | `pilot` | replacement-lineage token recall probes whether inherited editable context is physically usable across exec |

## Promotion Rule

Promote a blueprint into the active registry only when all of these are true:

1. the prompt leakage boundary is written down explicitly
2. the environment frame is stable enough to recreate
3. the scorer can distinguish intended success from easy accidental wins
4. the experiment has a plausible preservation and rerun story

For `L5`, prompt-boundary review must also classify every negative constraint:

- leakage seal: allowed, because it closes a non-target side channel
- strategy steering: not allowed, because it teaches which solution family is canonical

If a line says some variant of "do not invent/create a new bootstrap/startup
path/wrapper/alias/surface", redesign the fixture or scorer instead of keeping
that line in the prompt.

Until an experiment is runnable, keep the work here.
Once it is runnable but still pre-acceptance, keep it in the pilot surface
rather than polluting the active catalog.
