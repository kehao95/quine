# Manuscript Plan: Agent Mechanisms as OS Primitives

The dossier owns thesis, scope, evidence posture, and state decisions. This
file owns the public argument shape.

## Argument Plan

The post should make one move:

```text
The harness paper is right that prompt-only agents need executable, verifiable,
stateful structure. Quine accepts that diagnosis but lowers the answer into a
different substrate: agent mechanisms should map to OS primitives such as
processes, files, IO, permissions, namespaces, signals, and IPC.
```

The argument should proceed from agreement to ontological divergence to
mechanism mapping, while keeping current Quine implementation slices as
examples rather than as the north star.

## Section Skeleton

### 1. The Nearby Paper

- introduce `Code as Agent Harness` as a valuable adjacent framing
- summarize the shared diagnosis: prompt chains are not enough
- name the paper's central move: code as executable, verifiable, stateful
  harness around agent behavior
- transition: Quine agrees with the pressure but not with the first-order
  object

### 2. The Ontological Fork

- harness view: the agent is managed by an external code system
- Quine view: the design target is an operating-system ontology for agents
- explain why process, file, IO, permission, namespace, signal, and IPC are the
  right conceptual level for agent mechanisms
- anchor the sentence: `Quine does not just build a harness around the agent;
  it lowers agent mechanisms into OS primitives.`

### 3. Mechanisms Mapped To OS Primitives

- identity maps to process identity, PID, process tree, executable body, and
  runtime roots
- state and memory map to files, directories, append-only logs, tapes, and
  named state surfaces; anchor memory is one conforming probe, not the final
  memory doctrine
- tool use maps to IO and process interaction rather than only function calls
- verification maps to exit codes, fd separation, tests, traces, world
  mutations, and observable file state
- safety maps to permissions, namespaces, resource limits, filesystem overlays,
  and capability boundaries rather than advisory prompt policy
- multi-agent coordination maps to IPC: pipes, files, sockets, signals,
  process supervision, and process-external `ctl` surfaces
- recovery maps to restart, `exec`, checkpointed files, retained logs, and
  namespace rollback/adoption

### 4. Why This Matters

- a harness can still leave the agent as an app-level object
- Quine makes mechanisms inspectable, composable, permissionable, persistent,
  and discoverable as environmental facts
- the evaluation question shifts from "did the workflow run?" to "can the
  model discover and use OS-level affordances under pressure?"
- connect this to Quine's L1-L5 behavior acceptance ladder

### 5. Principle Versus Current Implementation

- this is not a claim that code harnesses are useless
- this is not a claim that Quine has fully implemented every downstream
  mechanism
- current Quine mechanisms such as anchors, overlay worlds, fork/spawn
  relations, and `ctl` surfaces are implementation probes that conform to the
  OS-primitive direction
- the stronger claim is architectural: substrate choice changes the shape of
  agent mechanisms

### 6. Closing

- restate the contrast:
  - code-as-harness moves agents beyond text chains
  - OS-primitive lowering moves agent mechanisms into computational physics
- close with the phrase `from harness to habitat` only if the concrete
  mechanics have already been established

## Evidence Placement

- section 1: cite adjacent arXiv paper
- section 2: cite Quine north star and README
- section 3: cite operating philosophy, core architecture, and runtime surface
  registry
- section 4: cite `EVALUATION.md`
- section 5: cite current non-claims from operating philosophy where useful

## Writing Notes

- keep the post shorter than an academic position paper
- avoid a long comparison table unless it materially improves readability
- use one compact table only after the ontological fork is explained in prose
- make the piece self-contained for readers who have not read Quine docs
