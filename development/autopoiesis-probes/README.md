# Autopoiesis Probes

Retained outputs from light-prompt, context-isolated fresh-agent probes live
here.

For continuation after context loss, read [`META-DIRECTION.md`](./META-DIRECTION.md)
first; it is the compact handoff over the current autopoiesis program.

The former `ACTIVE.md` goal-less probe router has been removed. Use
[`META-DIRECTION.md`](./META-DIRECTION.md) plus the latest cycle-local README
when choosing or creating a cycle namespace.

## Lifecycle

- Keep raw or shape-imperfect outputs when they reveal structural friction.
- For high-volume divergent cycles, add a per-cycle README and seed curation index
  instead of flattening raw outputs into one summary.
- For concurrent goal-less batches, record a substrate-mode receipt: prompt,
  agent ids/nicknames, shared-clean vs shared-dirty vs isolated worktree mode,
  base ref, parent worktree dirty-state relation, commits/actions, validation
  scope, and staged or index residue before/after. Dirty shared-worktree batches
  are stress tests, not clean cold-start selection assays; isolated worktrees
  test committed `HEAD`, not uncommitted parent worktree state.
- If `HEAD` moves during a shared-worktree probe, record the start/end refs and
  classify the run as a moving-substrate stress probe. Object-level commits from
  that substrate need path-scoped staging plus a post-commit boundary audit.
- Absorb repeated findings into the owning control-plane doc.
- Delete scaffolds once their findings are absorbed.
