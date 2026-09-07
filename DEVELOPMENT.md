# Development Control Plane

This document defines the control plane for agent-driven development and research engineering in Quine.

Public checkouts include this guide; some research inputs, profiles, and lab
checks are not shipped. See [public scope](./development/PUBLIC-RELEASE.md).

Use it together with:

- [`AGENTS.md`](./AGENTS.md) for repo-wide governance and constitution
- [`EVALUATION.md`](./EVALUATION.md) for evaluation runbooks and the 5-layer validation ladder
- [`CLOSURE.md`](./CLOSURE.md) for interval state transitions and durable artifact absorption
- [`Paper/philosophy/operating-philosophy.md`](./Paper/philosophy/operating-philosophy.md) for runtime design commitments
- [`Paper/core/registries/env-controls.md`](./Paper/core/registries/env-controls.md) for the canonical env-control matrix
- [`development/contracts/README.md`](./development/contracts/README.md) for reusable development contracts

## Purpose

[README](./README.md) introduces the project.
[Development direction](./Paper/philosophy/development-vector.md) connects ideas
and practical needs to architecture; [ontology](./Paper/philosophy/ontology.md)
offers distinctions for particular questions.

Quine development integrates runtime systems engineering, model-facing physics, and empirical research discovery. The objective of this control plane is to enable rapid, generative architectural exploration and rigorous ground-truth validation while keeping repository state, runtime contracts, and scientific evidence clean and reproducible.

## Public README Boundary

The root `README.md` is the public project introduction. It presents overlapping
descriptions of Quine and routes readers to published work and the repository.

Agent and maintainer navigation belongs strictly in the control plane:
- `AGENTS.md` owns repo-wide agent behavior and constitutional rules.
- `DEVELOPMENT.md` owns development routing, architectural exploration, environment contracts, and validation workflows.
- `development/status/README.md` owns live observation state (current posture, backlog, blockers, and recent transitions).

## Navigation And Routing

| Situation | Read first | Then check |
|---|---|---|
| repo-wide agent behavior and constitution | [`AGENTS.md`](./AGENTS.md) | — |
| live execution posture, backlog, blockers | [`development/status/README.md`](./development/status/README.md) | — |
| runtime philosophy and durable commitments | [`Paper/philosophy/operating-philosophy.md`](./Paper/philosophy/operating-philosophy.md) | [`Paper/philosophy/living-system.md`](./Paper/philosophy/living-system.md) |
| runtime architecture and file-native layering | [`Paper/core/architecture/layering.md`](./Paper/core/architecture/layering.md) | [`Paper/core/README.md`](./Paper/core/README.md) |
| environment control matrix and semantics | [`Paper/core/registries/env-controls.md`](./Paper/core/registries/env-controls.md) | [`env-process-boundary-brief.md`](./Paper/theory/views/runtime-capability/env-process-boundary-brief.md) |
| validation ladder and verification commands | [`EVALUATION.md`](./EVALUATION.md) | [`tests/README.md`](./tests/README.md), [`tests/runtime/COVERAGE_MAP.md`](./tests/runtime/COVERAGE_MAP.md) |
| closure and durable state absorption | [`CLOSURE.md`](./CLOSURE.md) | — |
| research program, theory graph, and evidence | [`Paper/README.md`](./Paper/README.md), [`Paper/theory/`](./Paper/theory/README.md) | [`Paper/evidence-map.md`](./Paper/evidence-map.md), [`experiments/PROTOCOL.md`](./experiments/PROTOCOL.md) |
| reusable development contracts | [`development/contracts/README.md`](./development/contracts/README.md) | — |

## Generative Architectural Exploration

Quine evolves through hypothesis-driven generative exploration. Rather than restricting mechanisms to defensive lowest-common-denominator guarantees, development is guided by structural discovery and capability expansion:

- **Generative Capability Expansion**: Design runtime interfaces, file projections, and coordination channels that unlock new model behaviors, emergent discovery, and autonomous self-modeling.
- **Vertical Prototyping**: Build end-to-end mechanisms across substrate, runtime, and model-facing layers to observe empirically how models discover and utilize new primitives under task pressure.
- **Hypothesis-Driven Design**: Frame development units around what capability or behavior becomes possible, what structural affordances enable it, and how the model discovers it.
- **Emergent Feedback**: Let runtime physics, environmental friction, and empirical evaluation shape the architecture. Harden and refine mechanisms once their behavioral utility is demonstrated.

For runtime development, clarify the following when relevant; conceptual and
control-plane work use the question, owner, and acceptance appropriate to their change:
- `goal` — what capability or behavior should become true
- `layer` — `substrate | runtime | usage | discovery | necessity`
- `env_profile` — `mac-kimi | linux-kimi`
- `acceptance` — minimum verification checks across deterministic and behavioral layers
- `artifacts` — key evidence, tapes, or baselines to preserve

## Environment Control Contract

Canonical matrix: [`Paper/core/registries/env-controls.md`](./Paper/core/registries/env-controls.md) is the single authority for `QUINE_*` semantics.

Environmental Process Boundaries:
- **Two levels, one pipeline, no synthesis.** A process's own environment is its birth record (read at `/proc/<pid>/environ`, immutable in effect). The environment constructed for child processes is separate (`child = stamps(boundary) ⊕ override ⊕ (environ − mask)` across `sh`, `fork`/`spawn`, and `exec`). Absence of a knob means compiled default (cataloged in `config/registry.json`); never synthesize defaults that nobody explicitly authored.
- **Enforcement is compiled; projections do not enforce.** Visibility views (such as `public/config/` projecting `registry.json`) expose governance but do not enforce it; enforcement is compiled directly into the binary and envp constructors. Non-claims are part of the contract.
- **Backend physics:** `overlay` is the transactional backend providing private lineage, rollback, and `switch_world`. `direct` is host-visible shared workspace physics. Validate on the actual backend required by the task.

## Projected Process Boundary

Public files serve as coupling surfaces; private cognition and history remain within the process's own interior (`everything is file, but not everything is in every process's world`). Treat runtime interfaces as deliberate projection boundaries rather than total transparency.

## Observation Overlay Contract

Text-native observer overlays and external viewers are observer projections, not the source of truth.
- Checked-in markdown, code, and canonical control-plane registries remain the ground truth.
- If an overlay disagrees with repository reality, repair the repo metadata or control plane rather than normalizing the overlay.

## Research Transitions & Exploratory Branches

Exploratory research moves through distinct state transitions:
- Checkpoint when an empirical test, prompt adjustment, or runtime tweak changes the working model or rules an approach in/out.
- Use exploratory branches when searching structurally uncertain or high-variance spaces.
- Retain useful findings, decisions, and unresolved alternatives in a form that supports continuation. Link a short orientation to fuller exploratory records where useful.

## Environment Profiles

| Profile | Target Environment | Reference Preset | Purpose |
|---|---|---|---|
| `mac-kimi` | macOS (Host) | `profiles/kimi-k2.7-code-litellm-medium.env` | Default development lane for substrate, tooling, and portable runtime logic. |
| `linux-kimi` | Linux VM / Container | `profiles/kimi-k2.7-code-litellm-medium.env` | Required for Linux workspace physics, overlayfs, mount namespaces, and FUSE. |

Ordinary live evaluations default to `profiles/gpt-5.4-codex-oauth.env`. Use dedicated profiles when testing specific provider transports or kernel physics.

## Validation Contract

Validation is grounded in reproducible truth:
- Control-plane changes use the relevant owner checks and `make check-control-plane`; review positioning, conceptual consistency, and actual routing separately.
- Conceptual work needs explicit distinctions and boundary cases; formal claims need assumptions and arguments; empirical claims need evidence for their scope.
- For software changes, `./tests/validate.sh` is the canonical entrypoint. Run `./tests/validate.sh --change substrate` for deterministic Go checks.
- Verification follows the 5-layer ladder: `substrate` (L1) -> `runtime` (L2) -> `usage` (L3) -> `discovery` (L4) -> `necessity` (L5).
- For new mechanisms or changed physics, verify both deterministic correctness and model-facing discoverability.

## Research Linkage Contract

When runtime semantics, environment controls, or claim surfaces change, update the owning control-plane documents in the same state transition:
- `README.md` for user-facing changes
- `Paper/core/registries/env-controls.md` for env knob changes
- `Paper/core/registries/runtime-surface.md` and related core papers
- `tests/runtime/COVERAGE_MAP.md` and `EVALUATION.md` for acceptance mappings
- `Paper/evidence-map.md` and theory graph objects for empirical discoveries

## Experiment Artifact Retention Closure

- Git stores control-plane state: source, runners, prompts, indexes, and `.dvc` pointers.
- DVC stores bulky generated payloads and retained evidence bundles.
- Private DVC retention is additive: `dvc push` inserts missing content-addressed objects into both private mirrors, and remote garbage collection is never part of routine operation. The primary host seals existing private-store objects immutable after each scheduled mirror run.
- The primary host enables the pre-push mirror and a persistent six-hour timer. A failure in one mirror must not prevent attempting the other.
- Transient scratchpad artifacts can be pruned freely; retain curated proof artifacts and canonical run baselines.

## Dual Repository Strategy

Work on `main`. `origin` is the private GitHub lab; `backup` is a local bare mirror. `public` is a publication surface updated only by deterministic projection to `release` upon human authorization.
