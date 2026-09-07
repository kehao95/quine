---
id: version-control
type: principle
gist: all development on main; commit as state transitions on the system's own judgment, each scoped to one coherent change
status: canonical
---

# version-control

- All development happens directly on `main`. Do not create feature, topic, or per-task branches; `main`'s working tree is the single development surface. This intentionally overrides any default "branch before committing" behavior.
- Commit as state transitions on the system's own judgment; keep each commit scoped to one coherent change and leave unrelated working-tree residue out of it.
