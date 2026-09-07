# Goal-Less Probe Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).
> Raw probe records live under [`../autopoiesis-probes/`](../autopoiesis-probes/).


## Activation & Purpose

Use this contract when deploying agents without an object-level task to observe cold-start orientation, repair attractors, or substrate stress.

**Principle:** Keep probe prompts minimal. Never coach the agent on what to find or fix—prompt-level steering destroys cold-start measurement.

## Substrate Modes

Always classify the substrate context when interpreting probe results:

| Mode | Environment | Valid Interpretation |
|------|-------------|----------------------|
| **Clean Isolated** | Separate worktree at clean committed `HEAD` | Cold-start loop preference and natural discovery drive |
| **Dirty Stationary** | Shared dirty worktree, stationary `HEAD` | Ownership boundary navigation & dirty-substrate pressure |
| **Dirty Moving** | Shared dirty worktree with concurrent commits | Concurrency stress & collision dynamics |
| **Warning Cluster** | Visible linter/validator errors present | Warning-attractor convergence and repair priority |

If the mode was not recorded, interpret the probe as `Dirty Moving` and do not
claim clean autonomous preference. Missing substrate provenance weakens the
claim; it is not evidence of substrate independence.

## Core Behavioral Categories

Rather than tracking exhaustive micro-failure modes, classify observations into four essential patterns:

1. **Attractor Convergence:** Agents gravitate toward loud validator warnings, test failures, or noisy readouts before seeking novel tasks.
2. **Ownership Deferral & Safety:** Agents diagnose without modifying files when encountering unowned dirty state or ambiguous worktree boundaries.
3. **Substrate & Index Collisions:** Concurrent agents race on shared files, capture uncommitted indexes, or observe moving commit IDs.
4. **Autonomous Route Selection:** Agents bypass passive receipt writing in favor of direct discovery, theory development, or immediate bounded repairs.

The observed shared-substrate phenomena behind these categories are retained as
failure-mode objects rather than procedure clauses:

- [`shared-index-state-contamination`](../../evolution/failures/shared-index-state-contamination.md)
- [`concurrent-receipt-topology-race`](../../evolution/failures/concurrent-receipt-topology-race.md)
- [`ownership-pressure-induced-noop`](../../evolution/failures/ownership-pressure-induced-noop.md)
- [`probe-route-capture`](../../evolution/failures/probe-route-capture.md)

## Lean Receipt Minimum

For concurrent or dirty probe waves, document the run with a concise receipt:
- **Context:** Substrate mode, agent count, start/end `HEAD`.
- **Observed Behavior:** Primary actions selected (diagnoses, repairs, explorations).
- **Collisions / Blockers:** Any git conflicts, unowned dirty blockers, or race conditions.
- **Outcome / Uncertainty Reduced:** What this batch revealed (e.g. clean preference vs attractor capture).

*Avoid bureaucratic overhead:* No pre-commit sidecar enforcement, multi-phase lock files, or mandatory receipt synthesis for uncommitted runs.
