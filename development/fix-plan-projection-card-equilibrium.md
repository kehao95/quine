---
title: Fix plan — retire the projection-maintenance card validator, remove the card population
status: implemented (Phases 1–4 landed 2026-07-18 in 32117633d, 699f7bdf2, ae5702af5, e0596564c; Phase 5 hardening remains deferred decisions)
kind: control-plane fix-plan
owner: handoff (any agent)
blocks: EVOLUTION.md Open Friction "validator-template equilibrium"; evolution/periods/23 closure
created: 2026-07-18
---

# Fix plan: validator-template equilibrium → retire the mechanism, keep the intent

> **Implementation note (2026-07-18):** Phases 1–4 landed as planned —
> validator + hooks + Makefile targets removed (`32117633d`), contract demoted
> to guidance (`699f7bdf2`), 137-doc card sweep with runnable drift-check
> salvage (`ae5702af5` + `e0596564c`), budget ratcheted 5039 → 4888.
> The `env-controls.md` card was removed at its generator
> (`internal/config/docgen.go`) and regenerated. Phase 5 items remain
> unexecuted deferred decisions. Narrative: period 23 § Key commits.

Diagnosis home: [`../evolution/periods/23-validator-template-equilibrium-jul02-jul18.md`](../evolution/periods/23-validator-template-equilibrium-jul02-jul18.md).
That file recorded three repair options and recommended option 1 (frontmatter migration).
This plan supersedes that recommendation after a requirement-level meta-analysis (§3):
the correct repair is **option 0 — remove the validator and the card population**,
because the sub-need the validator actually enforces is not machine-checkable, and the
sub-needs that are machine-checkable are already served by other, ground-truth-bound
validators.

## 1. Problem (short form)

`scripts/check-control-plane-projection-maintenance.sh` enforces, by substring-matching
five keywords inside a `## Projection Maintenance` section, that control-plane docs
"disclose how they go stale." Substring matching cannot reward honesty, only
token-presence, so the population's best response was template replication. The card
became the expected shape of a control-plane doc; the uniform five-field block is the
largest single source of the "docs read as AI-generated" symptom felt for ~8 weeks.

## 2. Empirical findings (2026-07-18 audit, three independent passes)

### 2.1 Spread

- `## Projection Maintenance` / `## Index Projection Maintenance` headings: **154 files**
  total; **138 genuine docs** after excluding 16 frozen run-snapshot prompt copies
  (`.../runs/.../context/prompt/10-agents.md`) and embedded repo clones.
- The validator's resolved candidate set is **70 files** (literals + globs after
  exclusions). 68 carry the heading; 2 satisfy the alternate marker path.
- **70 genuine docs carry the card voluntarily** — outside enforcement entirely
  (paradigm essays, design proposals, prediction tables, per-experiment READMEs).
  More than half the population is contagion, not compliance.
- The causal-path 4-tuple (`reader / effector`, `behavior effect`,
  `perception / validation`, `absorption owner`) appears in **26 further docs**,
  12 of them in `development/contracts/` where the vocabulary is native.
- History: the card was mass-seeded in one commit (`f7bffc880`, 2026-05-25,
  ~115 files at once), stayed flat to mid-June, then accreted to 138 via voluntary
  adoption and size-threshold-driven file splits.

### 2.2 Content quality (30-card sample, every reference checked against the tree)

- `projection role` and `absorption owner`: generic boilerplate in **26 of 28**
  classifiable cards. Zero file-specific information.
- `drift check`, the load-bearing field: only **~3 of 30** are genuinely
  ground-truth-bound (content-equality checks: `gen-env-docs -check` for
  `Paper/core/registries/env-controls.md`; `check-runtime-doc-sync.sh --strict` for
  `runtime-surface.md` / `COVERAGE_MAP.md`). ~8 more name real but self-referential
  doc-shape validators. **16 of 30 are unrunnable prose** ("compare with…",
  "re-read the source").
- The validator passes demonstrably false cards: a source owner pointing at the
  deleted `operator/qctl/` directory
  (`experiments/active/dormancy-ecology/analysis/real-world-coupled-design.md`),
  references to a non-existent `evolution/FRICTION.md`
  (`evolution/failures/README.md`), and verbatim clones differing only in a cycle
  number (`development/autopoiesis-probes/cycle-{9,10}-recombination/seeds/README.md`).
- What accuracy the cards do have is enforced entirely by **other** validators
  (`check-active-doc-links.sh` keeps markdown links alive). The one leak channel is
  bare backtick paths, which the link checker does not inspect — and that is exactly
  where the stale references above live.

### 2.3 Enforcement wiring

- No CI; enforcement is local hooks: `prek.toml` pre-commit entry (targeted files,
  `--strict`) and pre-push entry (`always_run`, repo-wide, `--strict`) — both blocking.
- Makefile: `audit-/check-control-plane-projection-maintenance` targets, aggregated
  into `audit-control-plane` / `check-control-plane`.
- `check-scripts-readme` / `check-validator-registration` enforce that the script
  inventory surfaces stay in sync — decommissioning must update those surfaces in the
  same change.
- Perverse coupling: `check-control-plane-size.sh`'s 200-line pain threshold
  incentivizes splitting files, and each split file needed its own card.

## 3. Requirement meta-analysis: was this ever a valid hook check?

The May 2026 motivating friction
([`contracts/control-plane-causal-path.md`](contracts/control-plane-causal-path.md)):
"rich self-description, but rich memory can still fail to change future action" —
index/status surfaces silently drift from source truth and become duplicate truth
stores. Decomposed, the need was four distinct sub-needs:

| # | Sub-need | Legitimate? | Machine-checkable? | Already served by |
|---|----------|-------------|--------------------|-------------------|
| N1 | a projection must name its canonical source (no duplicate truth stores) | yes | only as "referenced path exists" | `check-active-doc-links.sh`; registry inventory |
| N2 | drift must become visible before an agent acts on stale info | yes | yes — as content-equality / link-resolution checks | `gen-env-docs -check`, `check-runtime-doc-sync.sh`, `check-active-doc-links.sh` |
| N3 | authors must **think about** staleness when creating a surface | yes, as authoring discipline | **no — category error** | nothing can; only guidance + review can |
| N4 | surfaces with no credible causal path should be absorbed, not created | yes, as judgment | no | review judgment |

The validator was built to enforce N3 and N4 — the two sub-needs that are not
machine-checkable. Any machine translation of "the author thought about X" collapses
into "the text contains the keywords of X"; Goodhart collapse was the necessary
outcome of the design, not an implementation bug. Meanwhile N1/N2 — the
machine-checkable halves — were already served by the repo's ground-truth validators,
and the audit confirms the card validator added **zero** verification on top of them.

Supporting evidence that this was knowable:

- The contract itself warns: *"do not over-instrument every note; causal annotation
  exists to prevent self-model prose from becoming inert, not to create bureaucracy"*
  — the hook violated the anti-bureaucracy clause of the contract it enforced.
- Repo precedent (2026-07-12): the C11 path-string guard was **removed** rather than
  patched, as a "trivially bypassed, unenforceable fig-leaf." Same class of check;
  same correct disposition.

**Verdict: the requirement, as a hook check, was never valid.** The valid parts of the
underlying need either already have real enforcement or can only live as guidance.
Repairing the mechanism (frontmatter migration, option 1) would still machine-enforce
a declaration nobody consumes — a smaller fig leaf, but a fig leaf.

## 4. Disposition of the need

- **N3/N4 → guidance.** `contracts/control-plane-causal-path.md` stays, rewritten as
  authoring guidance and review questions ("if this index is wrong, how will a future
  agent notice before acting?"). No hook.
- **N1/N2 → ground-truth validators**, which exist. Optional hardening in §5 Phase 5.
- **The card and its validator → removed.**

## 5. Repair plan

Ordering matters and is favorable: the validator blocks files *lacking* cards, so
**decommission first, then sweep** — no transition mode needed. Each phase is one
commit (`git commit --only -- <paths>`; shared worktree).

### Phase 1 — decommission the validator

1. Remove both `prek.toml` entries: `check-control-plane-projection-maintenance`
   (pre-commit) and `push-check-control-plane-projection-maintenance` (pre-push).
2. Remove Makefile targets `audit-/check-control-plane-projection-maintenance`, their
   `.PHONY` entries, and their mentions in the `audit-control-plane` /
   `check-control-plane` aggregates.
3. Delete `scripts/check-control-plane-projection-maintenance.sh`.
4. Update the synced inventory surfaces in the same commit: `scripts/README.md`,
   `.githooks/README.md`, and whatever `check-validator-registration` reads.
5. Gate: `make audit-control-plane` and `prek run --all-files` green.

### Phase 2 — demote the contract to guidance

1. Rewrite `contracts/control-plane-causal-path.md`: keep the causal-path questions
   and the maintenance questions; delete the "fill the maintained projection-index
   card" requirement; add a short "card retired 2026-07 — see this plan and period 23"
   note so the history is discoverable.
2. Update the short route in `DEVELOPMENT.md` that points at the contract.

### Phase 3 — global card removal (the bulk of the de-templating)

1. Sweep the **138 genuine docs**: delete the `## Projection Maintenance` /
   `## Index Projection Maintenance` section wholesale. Scripted removal (heading →
   next `## ` or EOF) plus a manual salvage pass.
2. Salvage rule: a card line survives only if it names something real and
   non-obvious — in practice the ~11 cards whose drift check names a runnable
   validator. Fold that one line into the doc's intro ("drift check:
   `scripts/…​.sh --strict`") or drop it if the doc already links its validator.
   Boilerplate fields are deleted without replacement, including the registry's own
   card (its size-signal prose already lives in its body).
3. Leave the 16 run-snapshot copies untouched (frozen experiment evidence).
4. 4-tuple docs: keep the 12 in `development/contracts/` (native vocabulary); prune
   the rest case-by-case only where the table is filler.
5. The two stale references (deleted `operator/qctl/`, phantom `FRICTION.md`) live
   inside cards and disappear with them; verify no bare-path references to either
   remain afterward.
6. Gate: full `prek run --all-files` green; `check-active-doc-links` green.

### Phase 4 — lock in and close the narrative

1. Run `check-control-plane-size.sh`; ratchet `budget_lines:` in
   `control-plane-registry.md` down to the printed total (est. several hundred lines
   reclaimed across the 28 registered files alone).
2. Narrative bookkeeping: add the repair arc to `evolution/arc-index.md`; close the
   Open Friction entry in `EVOLUTION.md`; record repair + key commits in period 23.

### Phase 5 — optional hardening (separate decisions, not blockers)

1. **Bare backtick-path existence check** — extend `check-active-doc-links.sh` (or a
   sibling) to resolve `` `path/like/this` `` tokens in active surfaces. Closes the
   one real leak channel the audit found.
2. **Registry `drift_check` column + runner** — add a column to
   `control-plane-registry.md` for docs that have a real content-equality check, and a
   small hook that **executes** the registered commands (checks reality, not
   declarations). Only the ~3 genuine checks qualify today; empty cells stay empty.
3. **Registry admission as the minimal "stop and think" gate** — if pure guidance for
   N3 proves too soft, the fallback is: new governing docs must add a registry row in
   the same change (the registry already enforces this expectation prose-side; a
   presence check is falsifiable and does not Goodhart).

## 6. Follow-up debt (recorded, not fixed here)

Same disease, milder cases — validators that gate on prose-token/heading presence
rather than ground truth:

- `check-experiment-family.sh` (requires literal headings).
- `check-paper-dossier-shape.sh` — the `required_sections` half (its frontmatter half
  is fine).

Apply the §3 test to each before repairing: which sub-need is machine-checkable, and
is that part already served elsewhere?

## 7. Expected effect

| Metric | Before | After |
|---|---|---|
| Template cards in doc bodies | 138 | 0 (16 frozen snapshots remain as evidence) |
| Blocking hooks enforcing token presence | 2 (pre-commit + pre-push) | 0 |
| Honesty verification | none (substring match) | unchanged where real (content-equality checks), optionally extended (Phase 5) |
| "Expected shape" pressure toward the card | active | removed at the source |
| Control-plane line budget | 5039 | ratcheted down to post-sweep total |
