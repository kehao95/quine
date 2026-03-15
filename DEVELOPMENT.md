# Development Control Plane

This document defines the default control plane for agent-driven development work in Quine.

Use it together with:

- [`AGENTS.md`](./AGENTS.md) for repo-wide governance
- [`TESTING.md`](./TESTING.md) for test runbooks and command matrix
- [`tests/README.md`](./tests/README.md) for the 4-layer test map
- [`tests/model/README.md`](./tests/model/README.md) for instructional and emergent scenario catalog

## Purpose

Quine development is not only "write code and run tests."

Many changes affect:

- the runtime contract
- the model-facing physics
- the broader research control plane

The goal of this document is to make agent development legible and repeatable without over-hardcoding one workflow.

## Development Unit Contract

Before implementing a task, classify it as a development unit with these fields:

- `goal` — what should become true
- `layer` — `substrate | runtime | instructional | emergent`
- `env_profile` — `mac-kimi | linux-kimi`
- `coordination_model` — `sequential-single-runtime | concurrent-readers | concurrent-writers`
- `failure_semantics` — `best-effort | crash-recoverable | strongly-atomic`
- `portability_target` — `current-profile | named-profiles | cross-platform`
- `acceptance` — minimum checks that must pass
- `artifacts` — which tapes, test runs, or logs must be preserved
- `control_plane_impact` — which docs, tests, or private research surfaces may now be stale

This classification does not need its own file by default, but the agent should be able to state it clearly when working.

## Constraint-First Design Contract

Before exploring mechanisms, state the minimum guarantee surface the task actually needs.

This is where recent overengineering came from: the discussion silently escalated from "must recover after crash" to "must provide strong atomicity, symlink swaps, and cross-platform replacement semantics" before the task had named those guarantees as requirements.

Use this ladder:

- `best-effort` — transient inconsistency is acceptable; no special recovery promise
- `crash-recoverable` — intermediate state is acceptable during execution, but startup recovery must restore coherence
- `strongly-atomic` — intermediate state must not become externally visible

Default rule:

- choose the weakest guarantee that satisfies the task
- do not buy concurrency, atomicity, or portability complexity unless the task or Human explicitly requires it
- if the design starts needing journals, swap indirection, lock protocols, or platform-specific syscalls, first ask whether the stronger guarantee was ever part of the contract

Interpretation rules:

- `sequential-single-runtime` means no concurrent observer needs to be protected from temporary intermediate state
- `current-profile` means a mechanism only needs to be honest on the active development profile; do not silently promote it to a cross-platform claim
- "implementation ready" requires that the claimed guarantee surface is actually realizable on the named profile, not merely plausible in the abstract

Design review checkpoint:

- before converging on a mechanism, restate the task in this form:
  `goal + coordination_model + failure_semantics + portability_target`
- if that sentence sounds weaker than the proposed mechanism, the mechanism is probably overbuilt

## Environment Profiles

Quine currently uses two canonical development profiles.

| Profile | Use For | Model Default | Notes |
|---------|---------|---------------|-------|
| `mac-kimi` | default development on macOS | `.env.kimi-oauth` | default for most substrate, runtime, and non-Linux physics work |
| `linux-kimi` | Linux-only runtime and physics features | `.env.kimi-oauth` | run inside a Linux VM; on macOS the default repo path is Lima via the `colima` instance, while Podman machine is optional/manual rather than the primary harness; required for Linux workspace physics, mount namespace features, and related kernel-dependent behavior |

### Profile Rules

- default to `mac-kimi` unless the feature is Linux-dependent
- use `linux-kimi` for overlayfs, mount namespace, FUSE, or other Linux-kernel features
- when working from macOS, assume Lima/`colima` is the default validation backend unless the task explicitly names a Podman-based experiment
- when Linux behavior depends on local filesystem semantics, stage work onto a guest-local filesystem rather than a host-shared mount
- rootless user-namespace mounts are the preferred Linux workspace mode; if privilege is required inside Linux, treat it as a fallback assumption and record it explicitly in the handoff
- when a macOS host drives a Linux guest, treat guest architecture as part of the profile: the runtime binary executed inside the guest must be built for Linux, not merely be executable on the host

## Validation Philosophy

Quine validation is layered:

1. full deterministic substrate checks
2. targeted live runtime checks
3. targeted instructional model checks
4. targeted emergent acceptance for physical understanding

Two rules matter:

- for ordinary code tests, prefer running the full deterministic suite
- for behavior tests, prefer the smallest relevant scenario set rather than rerunning the whole catalog

## Validation Contract

| Change layer | Minimum validation |
|--------------|--------------------|
| `substrate` | `./tests/validate.sh --change substrate` |
| `runtime` | `./tests/validate.sh --change runtime --runtime <test>` |
| `instructional` | `./tests/validate.sh --change instructional --runtime <test> --instructional <scenario>` |
| `emergent` | `./tests/validate.sh --change emergent --runtime <test> --instructional <scenario> --emergent <scenario>` |

Interpretation rules:

- `./tests/validate.sh` is the canonical entrypoint after any change
- the script always prints all four layers so higher layers cannot disappear from view
- model scenarios should still be selected narrowly around the new or changed feature

## Behavior Acceptance Contract

Model-layer tests are split on purpose:

- `instructional` asks whether the model can execute the explicit protocol correctly
- `emergent` asks whether the environment itself causes the model to discover the intended primitive

For new or changed physics:

- prefer emergent pressure over explicit instruction when the claim is about discoverability
- treat instructional scenarios as diagnostics or wiring checks
- treat emergent scenarios as the highest acceptance layer for model-facing physics

The strongest behavior scenario asks:

> given this environment, does the model naturally discover that the new primitive is the right move?

not merely:

> can the model obey an instruction to call the primitive?

## Failure Contract

When a development unit fails, classify the failure before changing code.

Use these buckets:

- `substrate regression` — deterministic logic is wrong
- `runtime regression` — live binary contract is wrong
- `physics legibility` — the feature exists but the model does not discover or trust it
- `environment mismatch` — wrong profile, wrong kernel capability, wrong mount context, missing VM support
- `model limitation or instability` — the runtime contract may be valid, but the chosen model is not stable enough on that scenario

When blocked, preserve:

- exact command
- env profile used
- failing run directory or tape path
- short first-pass diagnosis

Linux workspace-specific interpretation rule:

- do not collapse `environment mismatch` and `physics legibility` into one diagnosis
- first prove the workspace substrate can mount, commit, and emit `[FS MUTATIONS]` in a minimal live run
- only after that should a larger experiment be used to judge whether the agent noticed the mutation signal or converted it into honest failure

## Research Linkage Contract

Development work can move the research control plane.

After a change, inspect whether any of these are now stale:

- public docs under the repo root
- validation docs under `tests/`
- experiment notes or summaries that explain the changed behavior

If the development change altered model-facing physics or the defensible claim surface, treat control-plane updates as part of the task, not optional cleanup. Private-lab-only research materials are updated on `lab`, not `main`.

## Dual Repository Strategy

This project uses a **private-lab primary branch** with a **selective public sync** model.

Current friction:

- "curated public subset" was being interpreted so narrowly that the public
  branch drifted away from the code, tests, and docs it was supposed to carry

```
┌─────────────────────────────────────────────────────────────┐
│  Local Repository                                           │
│                                                             │
│  lab ───────────────────────► origin/lab (kehao95/quine-lab)│
│   │                        (private repo, source of truth)  │
│   │                                                         │
│   └── selective sync ─────► main ───────► public/main       │
│                              (curated)      (public repo)   │
└─────────────────────────────────────────────────────────────┘
```

| Remote | URL | Purpose |
|--------|-----|---------|
| `public` | `git@github.com:kehao95/quine.git` | Open-source release |
| `origin` | `git@github.com:kehao95/quine-lab.git` | Private lab with experiments |

| Branch | Contains | Push To |
|--------|----------|---------|
| `lab` | Primary working branch: runtime work, experiments, private research materials, .beads/, AGENTS.md, private notes | `origin/lab` |
| `main` | Public-facing repo: runtime code, tests, and root docs, with experiments and private research materials remaining curated separately | `public/main` |

Operating rules:

- default to working in `lab` unless the task is explicitly about the public branch
- treat `main` as a curated release surface, not a mirror of `lab`
- move material from `lab` to `main` by intentional curation, but do not hide
  code, tests, or root docs by default
- when something is part of the public engineering artifact, prefer syncing it
  to `main` rather than keeping it lab-only
- private research materials, experiments, raw run trees, and other bulky evidence remain
  selective even when code/tests/root docs are public
- verify the target branch and remote before pushing, because the same local repo carries two different publication obligations
- when asked to make something public, prepare a clean public-facing repo
  surface, not a teaser-only subset
- when touching public-surface control-plane docs, run
  `./scripts/check-public-surface.sh` to catch stale path and policy wording

Public-surface drift check:

```bash
./scripts/check-public-surface.sh
```
