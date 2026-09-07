## Period XV: V0 Living-Runtime Spine (May 5)

### May 5 — Process/control contract to runtime-owned V0 spine

**Key commits:** pending

**What changed:** Landed a narrow V0 runtime spine that turns the
self-maintaining-system center into inspectable runtime-owned process surfaces:
`status/contract.json` now exposes `process-control/v0`;
`agent/<session>/world/{status.json,resources.json,events.jsonl}` exposes
`world/v0`; and `habitat/world/v0` gives the budgeted experiment world
state/rules/actions/score surfaces. The follow-up correction removed client
UI ownership from contract, world, maintenance, paper, and evaluation surfaces.

**Friction:** The repo had enough runtime organs to discuss self-maintenance,
but they were still split between prose, public control files, workspace-world
projection, experiment-local habitats, and client behavior. The first repair
overcorrected by letting the client appear to own diagnostics, maintenance
launch, and acceptance/paper boundaries. That made the living system legible,
but at the wrong layer.

**Structural repair:** Promoted the spine into machine-readable runtime surfaces
and updated the roadmap/status/evaluation registries so future agents encounter
`process-control/v0 -> world/v0 -> habitat/world/v0` as a bounded V0 sequence
rather than as unrelated feature ideas. Client UI remains outside the spine: it
can consume and lower into runtime-owned control surfaces, but it is not a
second control plane or an acceptance owner.

**What it taught:** A self-maintaining repository needs inspectable organs before
it can maintain itself through runtime behavior. The repair did not widen the
claim to total private worlds, callbacks, or revisable time; it made the current
body describable, diagnosable, and ready for live acceptance pressure.
