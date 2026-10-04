# Subjective Reality Novelty Read (2026-04-15)

This note records the current literature-based novelty judgment for the `subjective-reality` dossier.

It is a working read, not a final paper-position commitment.

## Current Read

Broad novelty for `subjective-reality` is constrained but not dead.

The literature already strongly covers:

- sandboxed execution with mounted workspaces
- explicit world state in evaluation harnesses
- capability shaping and tool/resource isolation
- branching / lineage / revert / merge over workspaces, repos, or context
- inspectable artifact traces and versioned workspace histories
- per-process projected worlds at the OS level via namespaces

So the following broad claims are not safely defensible as headline novelty:

- bounded workspace projection is new
- child worlds / branching are new
- inspectable runtime artifacts are new
- capability shaping is new
- per-process projected worlds are new

## Design Lineage

`Plan 9` private namespaces and related Unix namespace ideas are part of the deliberate inspiration lineage for this dossier, not an accidental after-the-fact similarity.

That does **not** make the overlap go away.

It means the dossier should treat them as:

- explicit conceptual ancestors
- a real prior-art boundary on broad systems claims
- support for narrowing the contribution toward agent-runtime semantics rather than generic namespace novelty

## Strongest Prior-Art Threats

- `AWCP` - environment/workspace projection as explicit delegation substrate
- `GCC` - inspectable filesystem-native branch / merge / trace substrate
- `SWE-Adept` - checkpoint reversion and branch switching over workspaces
- `AgentStepper` - fine-grained inspectable mutation history
- `Policy-Invisible Violations` - speculative copy-on-write overlay over world state
- `Resilient Write` - transactional write visibility and durable write-path semantics
- `IsolateGPT`, `Progent`, `Conseca`, `AC4A` - capability shaping as an explicit systems/security object
- `Plan 9` namespaces and Linux namespaces - projected per-process worlds at the OS level

## Claim-Level Read

### Weak / already-covered claim space

- projected workspace or mounted local root
- branching / child-world lineage in the broad sense
- inspectable world-like artifacts
- capability/resource restriction in the broad sense

### Narrower claim space that may still survive

- projected world as a first-class runtime ontology rather than only containment
- observation and mutation visibility as part of world semantics, not just logging
- capability/resource shaping as part of the operative world object rather than only security policy
- a unified child-world semantics across:
  - workspace state
  - capability surface
  - observation surface
- a stable inspectable `world` object that the runtime manages as one semantic surface rather than scattered logs, branches, and permissions
- evaluation explicitly targeting blindness reduction, mutation awareness, or discoverability of world semantics

## Repo-Level Implication

This dossier is in a stronger state than `revisable-time`:

- real runtime slice exists
- real deterministic tests exist
- real coverage rows already exist
- the main issue is not lack of substrate but lack of a frozen paper-facing evidence backbone

So the current best judgment is:

- keep the dossier active
- keep the claim narrow
- avoid inflating the current implemented slice into a total theory of capability/resource/world projection
- treat `Plan 9` / namespace lineage as declared ancestry and design guidance, while making the paper claim live at the agent-runtime semantics layer

## Most Defensible Residual Gap

If this lane becomes publishable, the strongest surviving thesis is likely to be a narrow conjunction:

> a projected operative world treated as a first-class runtime object whose semantics unify scope, capability, observation, mutation visibility, revision addressing, and child-world adoption.

That is narrower than “sandboxing” and narrower than “branching workspaces.”

## Scope Posture For Now

Current dossier posture should remain:

- active
- narrow
- implementation-backed where possible

The paper should not currently widen into:

- a total ontology of all capability/resource projection
- a generic sandboxing paper
- a merged world-plus-time paper
