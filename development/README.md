# Development Surfaces

Development control-plane support surfaces.

Public readers: [release scope and evidence availability](./PUBLIC-RELEASE.md)
explains historical records, redacted copies, and references to unreleased work.

[`../DEVELOPMENT.md`](../DEVELOPMENT.md) remains the canonical entrypoint for
development workflow. This directory holds extracted contract organs that were
large enough to deserve direct addressability without overloading the root.

## Surfaces

| Surface | Role |
|---------|------|
| [`contracts/`](./contracts/) | reusable development contracts extracted from the root entrypoint |
| [`status/`](./status/) | live observation state: execution posture, backlog, blockers, transitions |
| [`open-source/`](./open-source/) | open-source release scope register and public/private boundary decisions |

## Causal Path

| Field | Value |
|-------|-------|
| reader / effector | coding, research, and maintenance agents |
| behavior effect | lets agents jump directly from a development friction to the owning contract |
| perception / validation | active-doc links, root entry-doc checks, hook failures, repeated workflow correction |
| absorption owner | `DEVELOPMENT.md` for routing; local contract files for detailed rules |
