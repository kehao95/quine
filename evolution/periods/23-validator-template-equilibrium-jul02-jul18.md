# XXIII: Validator-Template Equilibrium & Control-Plane Steering Drift (Jul 2–18)

## Status

This period file is opened early to record a structural-friction discovery. The period spans
Jul 2–18 (today); ~200 commits (jul02–jul15) are not yet narrated here. Full commit-arc
coverage is open debt (below). The load-bearing content of this file today is the friction
entry, not the commit chronology.

## Arc

The period's first recorded event is a root-cause reattribution: the widely-felt "docs read
as AI-generated / unnatural" symptom was traced from prose style to a steering-layer
validator equilibrium. The fix landed the same day (below).

### Jul 18 — "AI tone" traced to a validator-produced templating equilibrium

**Friction (concrete).** Maintainer reported that root control-plane documents
([`DEVELOPMENT.md`](../../DEVELOPMENT.md), [`CLOSURE.md`](../../CLOSURE.md),
[`EVALUATION.md`](../../EVALUATION.md), [`AGENTS.md`](../../AGENTS.md)) read as increasingly
unnatural / AI-generated. Surface symptoms: run-on compound sentences, a mechanically
repeated five-field "Projection Maintenance" card at the head of every control-plane file,
twisted post-modifying noun-stacks, and jargon stacking. An attempt to rewrite
`DEVELOPMENT.md` for naturalness failed the size budget (812 → 1114 lines) *without*
addressing the cause — which exposed that the symptom was not prose style.

**Blindness duration.** Template enforcement landed ~May 25 (arc-index "May 25 — Index
projection card" arc: the five fields became required card fields). The "AI tone" was felt
for ~8 weeks and repeatedly misattributed to authoring style before being traced to the
validator mechanism on Jul 18.

**Root cause / mechanism.**
`check-control-plane-projection-maintenance.sh` (deleted by the repair; frozen copy:
[`public/computational-neoplasm/artifacts/check-control-plane/`](../../public/computational-neoplasm/artifacts/check-control-plane/check-control-plane-projection-maintenance.sh))
enforced staleness disclosure by **substring-matching five keywords** (`source owner`,
`projection role`, `freshness trigger`, `drift check`, `absorption owner`) inside a
`## Projection Maintenance` section. Substring matching has a decisive structural property:
it **cannot reward honesty, only token-presence**. Under that rule the cheapest way to
satisfy the check is verbatim copy of the canonical card, which carries no actual staleness
reasoning. The population's best-response is therefore template replication. As of Jul 18 the
card appears in **116 files** (60 machine-enforced via the validator's candidate list, ~56
adopted by voluntary contagion beyond enforcement), and the causal-path 4-tuple
(`reader / effector` / `behavior effect` / `perception / validation` / `absorption owner`)
in 26 more.

This is a genuine equilibrium, not loose metaphor: (1) it is a **coordination game** — once
verbatim cards are the expected shape, natural-prose authors deviate and get "repaired" back
to the card by consistency pressure; (2) it is a **mechanism-design** failure — the
validator's equilibrium ("contain the tokens") is not the designer's intent ("honestly
disclose how this surface goes stale"). The strict "Nash equilibrium" label is loose for the
single-author-vs-validator slice (that is a decision problem, not a game), but the
system-level fixed point of "rule + best-responding population" is a real bad equilibrium.

It couples **perversely with `check-control-plane-size.sh`**: the 200-line per-file
threshold incentivizes splitting, and each split creates a new file that needs its own card,
so the two checks jointly reward template proliferation. Splitting does not reduce the
governing surface — it breeds more cards.

**Structural repair.** Direction decided Jul 18, superseding the three options below: a
requirement-level meta-analysis found the check was never a valid hook — it enforced the
two sub-needs that are not machine-checkable ("the author thought about staleness";
"surfaces without causal paths should not exist"), while the machine-checkable sub-needs
(drift visibility) were already served by the ground-truth validators
(`check-active-doc-links.sh`, `gen-env-docs -check`, `check-runtime-doc-sync.sh`), to
which the audit showed the card validator added zero verification. The chosen repair is
**option 0: remove the validator and the card population; demote the causal-path
contract to authoring guidance** — executed and landed the same day (see Key commits
below). Full analysis and phased plan:
[`development/fix-plan-projection-card-equilibrium.md`](../../development/fix-plan-projection-card-equilibrium.md).
The historical options, recorded when option 1 was still the recommendation:

1. **Keep the intent, change the mechanism** — move the five staleness fields into
   frontmatter, validator checks frontmatter keys / ground-truth-bound signals instead of
   prose tokens. Rejected on meta-analysis: it would still machine-enforce a declaration
   nobody consumes — a smaller fig leaf, but a fig leaf.
2. **Narrow enforcement scope** to the top 5–6 root control-plane files; let the rest go
   natural. Cuts the contagion source but leaves the mechanism flawed.
3. **Accept the template, rewrite prose within it.** Logically doomed: individual effort
   cannot escape a bad equilibrium (coordination failure).

**What it taught (generalizable lesson).**

- A validator that checks its target by **proxy tokens** rather than **ground truth** will,
  once enforced, collapse onto its proxy and decouple from the latent variable (Goodhart's
  law). The five keywords were a proxy for "the author thought about staleness"; once
  enforced they became the target and stopped carrying the thought.
- The fix direction, from mechanism design: make the **honest declaration the cheapest
  satisfying strategy** (incentive compatibility), and bind it to **falsifiability-face**
  signals — measures checkable against reality (a drift command that actually runs, a
  frontmatter key whose referenced file exists, a link that resolves) rather than
  self-referential tokens. The repo already has this instinct in
  [`check-active-doc-links.sh`](../../scripts/check-active-doc-links.sh).
- **Recursive drift:** the repo's own theme is self-reference and drift, yet the steering
  layer's validators are themselves an unchecked drift surface. "Projection maintenance" was
  designed to keep docs honest about change; its implementation made docs uniform. The
  meta-mechanism (how validators change) is not itself incentive-compatible. This maps onto
  the repo's four layers: the validator is layer ② (steering); its bad equilibrium *is*
  drift; the repair strengthens layer ③ (falsifiability face), not the prose.

**Key commits.** The enforcing origin is the May 25 "Index projection card" arc
([`arc-index.md`](../arc-index.md)) and the introduction of
`check-control-plane-projection-maintenance.sh`. Repair landed Jul 18 across four commits:
`32117633d` (decommission the validator: script, prek.toml hooks, Makefile targets, README
inventory mentions), `699f7bdf2` (demote `control-plane-causal-path.md` to authoring
guidance), `ae5702af5` (sweep the card from 136 genuine docs, salvaging the genuinely
runnable drift-check lines into doc intro prose), and `e0596564c` (the two holdouts:
`experiments/protocol/README.md`'s card — its path trips a pre-existing, unrelated
dossier-link gap between `Paper/papers/emergent-self-reproduction/README.md` and several
`succ-*` experiments, so it landed with the pre-existing failures bypassed and documented —
plus the `budget_lines` ratchet in `control-plane-registry.md` from 5039 to 4888, the
post-sweep total). Full plan:
[`development/fix-plan-projection-card-equilibrium.md`](../../development/fix-plan-projection-card-equilibrium.md).
This repair was executed by two independent concurrent agent sessions on the same shared
worktree (one did Phase 1+2, the other took over for Phase 3+4); every subagent salvage claim
was independently re-verified against actual file content, not trusted from self-reports —
three of the ~19 claimed salvage lines had in fact been silently dropped by the shared-index
race and were restored before committing.

### Jul 19 — residual AI-tone measurement after the sweep

With the card population removed, a follow-up measurement asked what "AI tone"
remains and where. Em-dash density per 1k words: old core control-plane docs
2.1–3.8; `Paper/writing-principles.md` 0; the July diagnosis/repair-layer docs
themselves (this file, the fix plans, the census annex) 12.3–24.1. Triple-hyphen
compounds run ~20–24/1k in theory/analysis docs. Reading the post-sweep root
docs confirms the acute mechanical source is gone.

The residue is therefore not mechanism-driven: it is (1) authorial idiolect of
freshly agent-written analysis prose (em-dash asides, triadic lists) — densest
in the repair layer that diagnosed the problem; (2) vocabulary debt from the
repo's naming method (each coined pattern name becomes a term later docs must
carry); (3) corpus self-imitation — nearly all prose is agent-written, so new
sessions learn tone from the existing corpus with no heterogeneous baseline.
A size-budget hypothesis (lines-based metering rewards packed sentences) could
not be tested past the 80-column hard-wrap convention and stays unproven.

Repair: no style validator (prose naturalness is the same unenforceable-need
class as the retired card; a checker would recreate it). The norm owner was
extended instead: [`Paper/writing-principles.md`](../../Paper/writing-principles.md)
now scopes its Style And Fluency and Punctuation sections over internal prose.
Human hand-edits remain the only heterogeneous baseline.

### Jul 19 — proxy-enforcement census and family dispositions

Case zero's mechanism (substring/heading presence, not ground truth) is one instance
of a general failure mode, deposited as
[`evolution/failures/proxy-enforcement-collapse.md`](../failures/proxy-enforcement-collapse.md).
A repo-wide census of all 84 `check-*.sh` classified 38 ground-truth / 12
structural-shape proxy / 7 token proxy / 22 mixed, with 21 proxy-carrying validators
wired as blocking hooks — full table in
[`development/proxy-enforcement-failure-mode.md`](../../development/proxy-enforcement-failure-mode.md).

Per-validator dispositions, applying the requirement test (machine-checkable → bind;
not machine-checkable → guidance only; already served → don't duplicate), landed as
seven commits: `17c2650e9` / `cc3842e2f` (Round 1 — rebind `check-experiment-task-readmes.sh`'s
inverted literal-string match to real link resolution; delete the consumerless section
half of `check-experiment-family.sh`); `6f8e0ecf2` (delete the 10-section presence-only
half of `check-paper-dossier-shape.sh` — the two sections with real consumers stay
independently enforced by other validators); `5c11cf31d` / `d4b5c96d8` (delete
`check-experiment-current-state.sh` / `check-experiment-surface-map.sh` — both
consumer-hunted to zero: the surface-map case is proxy collapse in vivo, 29 of 32
triggered READMEs satisfied the heading while leaving 53 structural subdirs
unlinked); `f3fbe5e76` (rebind `check-experiment-note-index.sh` from heading-presence
to actual link-reachability of every sibling `NOTE-*.md`); `d7d66e17c` (delete
`check-dormant-curation.sh` — born-vacuous, 0 of 63 theory objects ever dormant);
`a21f45351` (self-label `check-spine-invariant-sync.sh`'s two checks as COARSE
proxies rather than force an unviable generate-and-compare rebind). The
authoring-guidance amendment landed as `118481d83`: a compact requirement-test
paragraph in
[`control-plane-causal-path.md`](../../development/contracts/control-plane-causal-path.md),
paid for by tightening ~20 lines of existing prose in the same file to hold the
4888-line control-plane budget at net zero.

Net: scripts 84 → 80, blocking proxy-carriers 21 → 16, per-dossier required
sections 10 → 0, per-experiment-README heading obligations −2. Remaining census
items (the doc-skeleton family, the budget de-registration escape, Paper-Feeds
null-token density, free-text experiment `status:`, friction review-date renewal)
are recorded as deliberate watch items, not open debt — the requirement test found
no urgent action justified for any of them today.
[`proxy-enforcement-collapse.md`](../failures/proxy-enforcement-collapse.md) is now
`status: repaired`. Executed by the same shared-worktree concurrent-session pattern
as the Jul 18 repair above — four items landed via one dispatched agent group, three
via another, both independently gate-verified against the same pre-existing baseline
failure set before and after.

### Jul 19 — prose-density campaign abandoned; the residual AI tone is left to fix-on-touch

A planned campaign would have rewritten the AI-dense strata to bring em-dash density into the
2-4 per 1k band. It was designed, redesigned once, and then abandoned without running a batch.
Three reasons, in increasing weight.

First, the redesign did not survive its own test. Version 2 replaced plain-language expansion
with terminology mapping, on the argument that expanding a coined term inlines its definition at
every use site and pays the description length repeatedly. The mapping research (78 terms, five
domain clusters, recorded in [`../../development/term-dispositions.tsv`](../../development/term-dispositions.tsv))
came back with almost no clean replacements: the repo mostly coins for genuinely new compounds
rather than renaming standard things, so the dominant disposition was to add a citation at the
defining surface and keep the term. That action adds text and leaves sentence structure
untouched, which does not serve a density goal at all.

Second, the campaign's own target was a proxy. Em-dash density is a measurable stand-in for
"reads as machine-written"; the latent variable is whether the maintainer experiences a document
as slop. The documents actually named in the original friction (the root control-plane files)
returned to the healthy band from the card sweep alone. The measured residual sits in theory
objects and analysis notes that were never the complaint. Driving the metric down across strata
nobody flagged is the same move as the retired card validator, one level up, committed by the
agent that had just finished diagnosing it. That is the reason of record.

Third, a single rewriter pass over one corpus risks replacing one uniform voice with another.

What survives: the norm in [`Paper/writing-principles.md`](../../Paper/writing-principles.md),
which now scopes internal prose, applied when a document is edited for other reasons. The
mapping table survives as a reference index with four content-correctness flags carried in it,
not as a work queue.

## Open debt carried forward

- **Commit coverage jul02–jul15 (~200 commits) is not narrated here.** This period file was
  opened early to record the friction above; the commit arc needs a separate tracing pass
  before this period can be closed.
- The pre-existing dossier-link gap between `emergent-self-reproduction` and the `succ-*`
  experiments (see Key commits above) predates this repair and still needs its own fix.
- The validator-equilibrium friction itself is now repaired, case zero and the
  repo-wide proxy-enforcement family alike; see the Jul 18/19 arcs in
  [`arc-index.md`](../arc-index.md).
