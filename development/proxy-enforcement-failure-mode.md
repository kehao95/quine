---
title: Proxy-enforcement collapse — repo-wide census annex and per-validator dispositions
status: closed (all dispositions executed 2026-07-19; Round 1 17c2650e9/cc3842e2f; §5 A1/B1/B2/B3/C1/C2/C3 6f8e0ecf2/5c11cf31d/d4b5c96d8/f3fbe5e76/d7d66e17c/a21f45351/118481d83; remaining §3.6/§3.7 items are deliberately no-urgent-action watch items, not open work)
kind: control-plane census annex + fix proposals
owner: handoff (any agent)
blocks: nothing (case zero already repaired); feeds future validator-authoring decisions
created: 2026-07-18
---

# Proxy-enforcement collapse: census annex and dispositions

The canonical failure-mode record is the deposited object
[`../evolution/failures/proxy-enforcement-collapse.md`](../evolution/failures/proxy-enforcement-collapse.md)
(definition, theory, repair grammar). Case zero — the projection-maintenance card
equilibrium — is diagnosed in
[`../evolution/periods/23-validator-template-equilibrium-jul02-jul18.md`](../evolution/periods/23-validator-template-equilibrium-jul02-jul18.md)
and repaired via
[`fix-plan-projection-card-equilibrium.md`](fix-plan-projection-card-equilibrium.md)
(landed 2026-07-18). This document is the **evidence annex**: the repo-wide survey
(all 84 validators + the non-validator mechanism layer) and per-instance
dispositions. Dispositions follow the object's repair grammar — **delete > rebind >
amend-existing; no new watcher structures** (sensors, schedules, lifecycle
machinery were considered and rejected: they re-instantiate the failure one level
up). This annex's shelf-life is bound to its consumption: if the disposition rounds
below never execute, prune this file rather than keep it as prose.

## 1. The failure mode, defined

**Proxy-enforcement collapse.** A mechanism wants a latent property (honest reasoning,
real linkage, maintained content) but verifies a *proxy* — the presence of tokens,
headings, or fields in the checked text. Once the check is **blocking** and applied to
a **population** of surfaces, the population converges on the cheapest satisfying
strategy: replicate the shape that passes. The proxy decouples from the latent
property (Goodhart), the population becomes uniform boilerplate, and the check now
manufactures *false confidence* — surfaces look maintained precisely because a machine
demands the look.

**Signature (detection heuristic).** A check is in this class when its pass condition
can be satisfied *by editing only the checked text* — no external state needs to be
true. Ground-truth checks require something outside the text: a referenced file
exists, a link resolves, generated output content-equals its source, a command exits
0, git state matches. Corollary: read the pass/fail logic, never the intention header
— the header states the latent property, i.e. exactly the thing the mechanism may not
deliver.

**Lifecycle observed in case zero** (predictive for other instances): mass seeding →
contagion beyond the enforced set (the shape becomes "what this kind of doc looks
like") → coordination pressure repairs natural deviation back to the template →
symptom misattribution (felt as "AI tone" / style, not mechanism) for weeks.

**The requirement test (reusable rule — apply before creating or repairing any
check).** Decompose the motivating need into sub-needs and ask of each: is it
machine-checkable against ground truth?

- Machine-checkable (existence, resolution, content-equality, execution, git state) →
  a hook may enforce it, bound to that signal.
- Not machine-checkable ("the author thought about X", "this surface deserves to
  exist") → guidance and review only. Any hook translation collapses into token
  presence; the collapse is a design consequence, not an implementation bug.
- Already served by an existing ground-truth check → don't duplicate with a proxy.

If nothing machine-checkable remains, the correct hook is **no hook** (precedent:
case zero; the C11 path-string guard removed 2026-07-12 as an "unenforceable
fig-leaf").

## 2. Survey scope and headline counts

Two passes, 2026-07-18: (a) mechanism classification of all **84** `scripts/check-*.sh`
by actual pass/fail logic; (b) non-validator mechanisms (frontmatter schemas, doc
templates, process rituals, tracked metrics).

| Class | Meaning | Count |
|---|---|---|
| A — ground-truth bound | pass requires external reality | 38 |
| B — structural-shape proxy | required headings/sections exist | 12 |
| C — token-presence proxy | required keywords appear | 7 (incl. retired case zero) |
| D — mixed | part proxy, part ground truth | 22 |
| non-enforcing readouts | print signals, no pass/fail | 5 |

**21 proxy-carrying validators (B/C/D) are wired as blocking hooks** in `prek.toml`.
Blocking + proxy + large population is the dangerous combination — that is what
produced case zero. No *second* instance has reached case-zero scale (138 files of
boilerplate), but several have the preconditions.

## 3. Instances, ranked, with proposed remediation

Dispositions follow the §1 requirement test. Deliberately **not** "delete all
proxies": a shape check with a real consumer is legitimate once bound to that
consumer.

### 3.1 `check-paper-dossier-shape.sh` — the structural twin (highest risk)

Blocking, pre-commit + pre-push, over every paper/essay dossier. Requires 10 `##`
sections to *exist*; content never inspected. Cheapest strategy: paste the 10-heading
skeleton with empty bodies — an empty "In Scope" / "Unsafe Overclaims" passes. Its
frontmatter/H1-vs-registry half is genuine ground truth and should survive.

**Proposal.** Apply the test per section: (a) sections with a machine-checkable core —
e.g. evidence/feeds sections whose rows must resolve — keep, bound to resolution
(`check-paper-dossier-links.sh` already does part); (b) sections that exist for
progressive disclosure — keep at most a *small* required set and add a non-stub guard
(the `STUB_RE` pattern already exists in `check-theory-decoder.sh`); (c) the rest →
advisory or drop. Decide before this family breeds its own card equilibrium: the
population is every future dossier.

**Verdict settled 2026-07-18 (per-section consumer audit): delete the entire
`required_sections` mechanism — 10 → 0, stronger than the proposal above.** The only
two machine-consumed sections (`## Experiment Feeds`, `## Current Evidence Backbone`)
are *independently required and content-validated* by `check-paper-dossier-links.sh`
and `check-paper-evidence-backbone.sh` (both `--strict`, both stages), so the shape
check's section half is fully redundant for them; the other 8 sections have zero
consumers (repo-wide anchor search empty; `focus-paper.sh` is an absence-tolerant
manual readout). Non-stub guards are rejected: the two kept-by-others sections are
link-only tables by design (a "non-link-only line" guard would false-fail 33–34 of
34 dossiers). Execution: §5 item A1.

**Executed 2026-07-19 — DELETE (`6f8e0ecf2`):** `required_sections` dict and its
report loop removed; frontmatter/H1/registry checks untouched. `--strict`
repo-wide: `docs_scanned=34 issues=0`. `check-paper-dossier-links.sh` /
`check-paper-evidence-backbone.sh` still independently enforce the two real
sections, output byte-identical to pre-change (verified via `git stash`).

### 3.2 Experiment-README heading-spine family (largest population)

`check-experiment-current-state.sh`, `check-experiment-surface-map.sh`,
`check-experiment-note-index.sh` — blocking, rglob over all experiment READMEs. Each
fails only on a missing anchored heading; a stub line under the heading passes.
Intent/mechanism gaps: "summary **near the top**" is not checked (search is
file-wide); "keep directories **navigable**" reduces to a heading existing.

**Proposal.** These have the best case for a real consumer: navigation. Bind them to
it — a required heading is legitimate if something links to its anchor (anchor
resolution is already ground-truth in `check-active-doc-links.sh`) or a tool reads it
(obsidian bases, surface maps). For each of the three: name the consumer; if one
exists, keep + add non-stub guard; if none, downgrade to advisory. Do not add more
members to this family without the test.

**Verdicts settled 2026-07-18 (consumer hunt + rebind dry-runs; zero consumers found
for all three headings over the experiment population):**
- `check-experiment-current-state.sh` — **delete.** No parser anywhere reads the
  heading on experiment READMEs (the `focus-paper.sh` / dossier-shape hits are the
  separate Paper-dossier population); the real triggered population is 2 task
  READMEs, double-covered by the rebound `check-experiment-task-readmes.sh`.
  Execution: §5 item B1.
  **Executed 2026-07-19 — DELETE (`5c11cf31d`):** script + prek.toml hook pair +
  Makefile targets/tokens + `.githooks/README.md` prose + `experiments/PROTOCOL.md`
  convention paragraph removed together. `check-experiment-task-readmes.sh`
  confirmed to still cover task-README current-state.
- `check-experiment-surface-map.sh` — **delete.** All 32 triggered READMEs carry the
  heading, yet 53 structural subdirs across 29 of them are not linked (e.g. coop-05's
  "Surface Map" is a zero-link ASCII tree): the heading is satisfied while delivering
  no navigation — proxy collapse in vivo. The parent→child-link need fails
  re-derivation (29/32 violating with no felt pain, no consumer), and a rebind would
  *add* 53 links of enforcement burden. No other check covers parent→child (subdir-
  index = child-has-README; surface-ownership = child→ancestor) — that direction goes
  deliberately unenforced; revisit only on evidence of real navigation pain.
  Execution: §5 item B2.
  **Executed 2026-07-19 — DELETE (`d4b5c96d8`):** symmetric to B1's removal.
- `check-experiment-note-index.sh` — **rebind to coverage, drop the heading.** Honest
  form: every sibling `NOTE-*.md` must be link-reachable from the README — no
  particular heading required. Tracked population today: 1 README / 4 notes /
  0 orphans → zero migration burden. Execution: §5 item B3.
  **Executed 2026-07-19 — REBIND (`f3fbe5e76`):** heading check replaced with
  markdown-link resolution against every sibling `NOTE-*.md`; a live gitignored
  `source-code`/`.disabled-*` residue exclusion was added while in the file
  (mirroring surface-map's). Fixture test (unlinked NOTE file) confirmed the
  rebind actually discriminates: fails when orphaned, clean once linked.

### 3.3 `check-experiment-task-readmes.sh` — inverted ground truth (mechanical fix)

Blocking. Matches the *literal string* `(../../README.md)`: a correct parent link
written with a different relative path **fails**; a broken link containing the exact
text **passes**. The check is anti-correlated with its intent.

**Proposal (uncontroversial, do first).** Replace the literal match with actual link
resolution — "README contains ≥1 resolving link whose target is the parent
experiment README" — reusing the `check-active-doc-links.sh` resolution logic.

**Executed 2026-07-18 — REBIND (`17c2650e9`):** link-resolution mechanism landed; the population glob was also re-aimed from the dead pre-restructure `experiments/p*` layout (scanning 0 docs since `18fe02b39`) to `experiments/active` (4 docs), both inverted fixture cases flipped.

### 3.4 `check-experiment-family.sh` — bare-substring sections (mechanical fix)

Blocking, small population. `sec not in text` — a section *name mentioned in prose*
satisfies it; weaker even than the anchored-heading checks. Its frontmatter half
(`family` == directory name) is sound.

**Proposal.** Anchor the section test (`^## <name>` + non-stub) or drop the section
half and keep the frontmatter half. Either is a one-line-class fix.

**Executed 2026-07-18 — DELETE (`cc3842e2f`):** consumer test found zero inbound anchor links and zero parsing scripts for all five `FAMILY_SECTIONS`; the section half was deleted, the frontmatter/registration half (ground truth) retained.

### 3.5 Advisory C-class miniatures (low pressure, same disease)

- `check-dormant-curation.sh` — pass = substring `reactivation`/`reactivate` anywhere.
  A literal miniature of case zero. **Proposal:** same verdict — either bind ("names a
  resolving target surface for reactivation") or retire; do not leave as ritual.
  **Executed 2026-07-19 — DELETE (`d7d66e17c`):** population was 0 of 63 theory
  objects (born-vacuous, never demonstrated discriminative power per §1 point 5);
  script + Makefile targets/tokens + `scripts/README.md` mentions + the
  `MANUAL_ONLY` entry in `check-validator-registration.sh` + the retired row in
  `Paper/procedure/sensor-catalog.md` removed together. Manual-only — no prek
  entries existed.
- `check-evolution-hygiene.sh` — owner/status/review tokens present; placeholder
  values pass. **Proposal:** validate values (owner ∈ known set; `review:` date parses
  and is not past — the past-date half makes it ground-truth against the calendar).
  **Disposition settled 2026-07-19 — NO CHANGE.** The date half (parse + not-past,
  ground-truth against the calendar) is already implemented; the owner-∈-known-set
  half was withdrawn — it would require inventing a set nothing sources, which is
  structure accretion, not a repair. Execution: §5 item C3.2.
- `check-spine-invariant-sync.sh` — token co-occurrence across two files; files can
  diverge in meaning while sharing tokens. **Proposal:** candidate for the
  generate-and-compare pattern (one canonical source, paraphrase checked by
  generation), else accept and label as coarse.
  **Executed 2026-07-19 — KEEP-AND-LABEL (`a21f45351`):** a generate-and-compare
  rebind was rejected (the canon is human-authored prose with no machine source to
  generate from); header now self-labels both SYNC-DRIFT and TELOS-REGRESSION as
  COARSE proxies, modeled on `check-theory-code-consistency.sh`'s wording.
  Comment-only diff, zero logic change (verified byte-identical pass/fail output).
- `check-runtime-doc-sync.sh` (blocking, D) — its git half forces an owner doc to be
  *staged* with semantic code changes but not *corrected*; the "STALE phrasing" greps
  only catch enumerated bad phrases. **Proposal:** accept (the staging force-function
  has real value; full correctness is not machine-checkable) but note the limit in its
  header rather than implying more.
- `check-theory-code-consistency.sh` self-labels its vocabulary overlap "a COARSE
  proxy" — the honest way to carry a proxy. Model wording for the above.

### 3.6 Doc-skeleton family (small fixed sets — lowest priority)

`check-root-entry-docs.sh`, `check-tests-entry-docs.sh`, `check-paper-track-readmes.sh`,
`check-experiment-archive-readmes.sh`, `check-experiment-incubator-readmes.sh`,
`check-scripts-readme.sh`, `check-githooks-readme.sh`. Blocking, but populations of a
handful of hand-curated files each; equilibrium pressure is real yet bounded.
**Proposal:** no urgent action; apply the test opportunistically when any of them next
changes. Watch one aggregate signal: total required-section count across the repo
should not grow without a named consumer.

### 3.7 Non-validator instances

- **Budget de-registration escape (moderate; one confirmed use).** The size ratchet
  ([`../control-plane-registry.md`](../control-plane-registry.md) `budget_lines:`)
  meters only registered files. Commit `429fdd600` ratcheted 5598→5027 by
  *de-registering* a live file — "the file stays live and linked; it is just no longer
  budget-counted." Human-justified, but the escape is proven. **Disposition: record
  only, no mechanism.** An earlier draft proposed a de-registration visibility
  sensor; withdrawn under the repair grammar (watcher structure for a one-occurrence,
  human-justified event). The escape is now named here and in the failure object;
  if it recurs, revisit with that evidence.
- **Paper Feeds `none-yet` density (benign-by-design; watch).** 68 of 166 experiments
  (41%) satisfy the required `## Paper Feeds` section with only the null token. The
  token is an honest "no paper yet" and MISSING is only a warning, so pressure is
  weak. **Proposal:** no change; re-check the density quarterly — if it grows while
  the paper layer grows, the section has become ritual.
- **Free-text experiment `status:` (recognized; migration planned).** 35 experiments
  say `active` while ~25 days untouched; nothing reads the field.
  [`status/README.md`](status/README.md) already plans the closed-set migration.
  **Proposal:** no new mechanism; the planned migration is the fix.
- **Watch item — friction review-date renewal.** `EVOLUTION.md` Open Friction entries
  carry `review:` dates; renewing a date is cheaper than resolving the friction, so
  the mechanism is structurally primed for ritual renewal. No churn yet.
  **Disposition: record only.** The signature is repeated `review:` renewal without a
  status change — recognizable by eye at review time; no counting machinery (the
  repair grammar rejects watcher structures ahead of evidence).
- **Watch item — budget-ratchet under-compression incentive (theoretical; unobserved).**
  A terminology-mapping pass flagged that `check-control-plane-size.sh`'s ratchet-down
  convention has the shape Weitzman's ratchet principle warns about (repeated
  performance-based tightening can incentivize under-revealing capacity to avoid a
  harder future target). No occurrence has actually been observed: retracing
  `budget_lines:` history shows every value change is either the deliberate one-time
  card-sweep compression (5039→4888, `e0596564c`) or the de-registration escape above
  — not a pattern of editors stopping just short of a lower total. One session
  initially misread landing exactly on budget after a routine same-file compression
  as a live instance of the effect; on review that edit inherited an already-set
  ceiling and showed ordinary satisficing against a binary gate, not strategic
  capacity-hiding — recorded here as the corrected disposition, not the original claim.
  **Disposition: record only, no mechanism.** The risk is real as a design property to
  watch for over many editors and iterations; it is not evidence of present behavior,
  and normal editors doing the minimum to pass a gate is not, by itself, the signature
  to watch for — a real instance would show a *repeated* pattern of edits stopping just
  above budget when a materially lower total was reachable.

### 3.8 Immune exemplars (what keeps a check honest — copy these shapes)

- **Resolution:** `check-active-doc-links.sh` — link passes only if the referent
  exists (filesystem + git-tracked).
- **Generate-and-compare:** `gen-env-docs -check`, `generate-paper-projections.py
  --check` — doc content-equals its generated source; drift cannot hide.
- **Derive-and-verify:** `check-theory-maturity.sh` — the declared `maturity:` value
  is recomputed from the body and must match; a copied default fails.
- **External-state binding:** `check-selfsource-freshness.sh` (git ancestry),
  `check-githooks-active.sh` (git config), `check-control-plane-size.sh` (wc + ratchet),
  `check-validator-registration.sh` (set-diff scripts vs hooks).

## 4. Meta-mechanism guard (the actual lesson)

Case zero's deepest finding: the steering layer's validators are themselves an
unchecked drift surface — nothing reviews whether a new check's mechanism can deliver
its intent. **Proposal:** add the §1 requirement test to the validator-authoring
guidance (the `repair-guidance` conventions in `scripts/README.md` and the authoring
contract in [`contracts/control-plane-causal-path.md`](contracts/control-plane-causal-path.md)):
before adding a `check-*.sh`, state which sub-need is machine-checkable and what
ground-truth signal binds it; a pass condition satisfiable by editing only the checked
text requires either a named consumer or a `COARSE proxy` self-label. This stays
**guidance, not a hook** — hookifying the anti-hookification rule would recreate the
failure mode one level up.

## 5. Disposition rounds (when someone picks this up)

Every disposition must state the re-derived requirement and its verdict verb
(**delete** / **rebind** / **amend**). Net structure across all rounds must be
negative or zero: no new scripts, no new required sections, no new schedules.

**Round 1 — rebind the outright-wrong mechanics (no decisions needed).**
- §3.3 task-readmes: rebind the literal `(../../README.md)` match to actual link
  resolution ("≥1 link in the task README resolves to the owning experiment
  README"). Acceptance: the two inverted cases flip — a correct link written with a
  different relative path passes; a broken link containing the verbatim string fails.
  — executed 2026-07-18, **REBIND**, `17c2650e9` (both fixture cases flipped; see §3.3).
- §3.4 experiment-family: grep for consumers of the `FAMILY_SECTIONS` anchors
  (inbound `#` links, parsing scripts). None found → **delete** the section half,
  keep the frontmatter half (ground truth). Found → rebind to anchored
  `^## <section>` + non-stub.
  — executed 2026-07-18, **DELETE** (no consumer found), `cc3842e2f` (see §3.4).

### Migration runbook — batches A/B/C (evidence-settled 2026-07-19, approved; execute in order)

Verdicts below were settled by three read-only evidence passes (per-section consumer
audit; spine-family consumer hunt + rebind dry-runs; Round-4 placement + wiring
inventory) — the evidence summaries live in §3.1/§3.2/§3.5 verdict blocks. General
discipline for the executor:

- **Line numbers in this runbook are as-of 2026-07-18 and prek.toml has shifted since
  (`17c2650e9`): always re-locate by hook id / target name / regex, never by raw line
  number.**
- One commit per item; shared worktree — `git commit --only -- <explicit paths>`,
  then verify with `git show --stat HEAD` that only intended files landed.
- Baseline `.githooks/run-prek --all-files` BEFORE the first change; after each
  batch the failing set must be byte-identical to baseline. Known pre-existing
  failures (do not fix, do not be blocked): check-active-doc-links (6 BROKEN),
  check-paper-dossier-links-from-experiments (succ-*), check-paper-evidence-
  backbone-from-experiments, check-experiment-prompt-posture, check-runtime-doc-sync.
- A script deletion and ALL its wiring (both prek hook ids, Makefile targets +
  `.PHONY` tokens + aggregate-line tokens, doc mentions) move in the same commit,
  or prek errors on a missing file. `check-validator-registration.sh` stays silent
  on deletions (set-diff only sees existing scripts) — but clean up stale
  `MANUAL_ONLY` entries in the same commit where noted.
- Batches are independent (no behavioral cross-dependency confirmed); within a
  batch, items are ordered but separable.

**A1 — `check-paper-dossier-shape.sh`: DELETE the entire `required_sections`
mechanism (10 → 0).**
Edits: remove the `required_sections` dict + the section-presence loop; keep the
frontmatter-keys and H1-vs-PORTFOLIO-registry logic untouched (the class-A half);
rewrite the VALIDATOR INTENTION header — state that section presence is deliberately
unenforced and must not be reintroduced without a named machine consumer. Both prek
entries (`check-paper-dossier-shape`, `push-check-paper-dossier-shape`) stay — the
script lives on. There is no Makefile target for this check (none existed).
Acceptance: `--strict` repo-wide green over all 34 dossiers; then run
`check-paper-dossier-links.sh --strict` and `check-paper-evidence-backbone.sh
--strict` to confirm `## Experiment Feeds` / `## Current Evidence Backbone` remain
independently enforced. Optional same-commit nicety: a one-line note in
`Paper/templates/paper-dossier-template.md` that only the two feeds/backbone
sections are machine-required; the rest are authoring aids.

**B1 — `check-experiment-current-state.sh`: DELETE the script.**
Pre-step: confirm the rebound `check-experiment-task-readmes.sh` still enforces a
current-state heading on task READMEs (the double-coverage fact cited in §3.2 was
verified before Round 1 landed; re-verify against the current script).
Edits (re-locate by id/name): delete the script; prek.toml remove ids
`check-experiment-current-state` and `push-check-experiment-current-state`;
Makefile remove `audit-experiment-current-state` / `check-experiment-current-state`
targets, their `.PHONY` tokens, and their tokens in the `audit-control-plane` and
`check-control-plane` aggregate lines; `.githooks/README.md` remove the pre-commit
summary phrase ("current-state summary structure") and the description bullet —
leave the task-readmes line that also names `## Current State`;
`experiments/PROTOCOL.md` trim the current-state convention prose (~L1251–1256).
Acceptance: `rg 'check-experiment-current-state'` finds only historical/frozen
references (autopoiesis-probe tapes, public frozen excerpt, paper materials, this
annex); full gate run baseline-identical.

**B2 — `check-experiment-surface-map.sh`: DELETE the script.**
Edits: symmetric to B1 — script; prek ids `check-experiment-surface-map` /
`push-check-experiment-surface-map`; Makefile targets + `.PHONY` + both aggregate
tokens; `.githooks/README.md` bullet; `experiments/PROTOCOL.md` surface-map prose
(~L1258–1270). Same acceptance pattern as B1.

**B3 — `check-experiment-note-index.sh`: REBIND to coverage; the heading
requirement is dropped.**
Edits: same script, new mechanism — for each README with sibling `NOTE-*.md`
(keep the current population scoping and exclusions), every sibling `NOTE-*.md`
must be the target of ≥1 resolving markdown link (inline or wikilink) from that
README; no `## Note Index` heading is required at all. Rewrite the header. prek
entries and their `files` regex stay unchanged (still triggered by NOTE/README
changes). Optional tightening while in the file: mirror surface-map's dot-dir /
`source-code` exclusions so audit mode stops scanning gitignored run residue.
Acceptance: population run = 0 issues (tracked reality: 1 README / 4 notes /
0 orphans); create a temporary unlinked `NOTE-zz-fixture.md` beside the real
note-bearing README — `--strict` must fail — then delete the fixture and verify
`git status` clean.

**C1 — `check-dormant-curation.sh`: DELETE the script.**
(Population is empty — 0 of 63 theory objects are dormant; born-vacuous; no
ground-truth reactivation-target form exists to rebind to.)
Edits: delete the script (manual-only — no prek entries exist); Makefile
`audit-dormant-curation` / `check-dormant-curation` targets + `.PHONY` tokens (not
in any aggregate); `scripts/README.md` individual description line (~L143–144) and
the family mention (~L66) if it names the script; remove the stale `MANUAL_ONLY`
allowlist entry in `check-validator-registration.sh` (~L103); check
`Paper/procedure/sensor-catalog.md` (~L23, the only "dormant" definition) and update
its row if it names the script. Acceptance: `rg 'check-dormant-curation'` finds only
historical/frozen references.

**C2 — `check-spine-invariant-sync.sh`: header self-label only (keep-and-label).**
Edits: extend the `purpose:`/`repair-guidance:` comment block — token co-occurrence
is a **COARSE proxy** for semantic agreement (two files can share all four tokens
while diverging in meaning); the FORBIDDEN regex is a coarse proxy for
telos-regression; deep semantic sync is deferred to human spine decisions. Model
wording: `check-theory-code-consistency.sh` (~L24–27, L33). No mechanism or wiring
change. A generate-and-compare rebind was evaluated and is not viable (the canon is
human-authored prose with no machine source).

**C3 — contract amendment + closeout (one commit).**
Budget warning: the contract is a registered control-plane file and `budget_lines:`
sits at its ratchet limit — adding ~8 lines WILL breach the blocking budget check
unless the same change compresses ≥8 lines in this or another registered file.
Plan the compression before committing.
1. `development/contracts/control-plane-causal-path.md`: insert after the line
   "If the answer is unclear, the surface is not yet a durable control-plane
   repair." (~L74) and before `## History` a short block — "Requirement test for
   validators (before adding or repairing any check-*.sh)": the §1 trichotomy
   (machine-checkable → bind to that signal; not machine-checkable → guidance/review
   only; already served → don't duplicate), the detection heuristic (pass condition
   satisfiable by editing only the checked text ⇒ proxy), and the verb order
   (delete > rebind > amend-existing; a new check needs the test passed AND a named
   consumer; never a meta-check on checks). ~8 lines, same genre as the existing
   review-question lists.
2. `check-evolution-hygiene.sh`: **no change** — record here that the proposed value
   rebind was withdrawn: the date half (parse + not-past → FRICTION-AGED) is already
   implemented; the owner-∈-known-set half would require inventing a set that
   nothing sources (structure accretion).
3. Closeout: mark each executed item in this annex (§3 verdict blocks + this
   runbook) with its commit hash; then advance
   `evolution/failures/proxy-enforcement-collapse.md` per the failures/README
   lifecycle — `applied: true`, `status: repaired`, and point its closing section at
   the amended contract + the batch commits; add one arc line to
   `evolution/periods/23-validator-template-equilibrium-jul02-jul18.md` recording
   the family dispositions.

**Executed 2026-07-19 — `118481d83`** (item 1, contract amendment; item 2, no-change
recorded above in §3.5): the requirement test landed as a compact paragraph, not the
full trichotomy prose — the budget warning above was correct, the addition had to be
paid for by tightening ~20 lines of pre-existing prose elsewhere in the same file to
land the change at net zero against the 4888-line ratchet. Item 3 (this edit + the
failure-object advance + the period-23 arc line) executed in the same pass.

**Net effect, landed:** scripts 84 → 80 (A1 kept its script, edited; B1/B2/C1 deleted
theirs; net −3 from those, plus one script already gone before this batch); blocking
proxy-carriers 21 → 16 (B1, B2, and case-zero's already-removed validator drop out;
A1 stays wired but its proxy half is gone; B3/C2 stay wired/manual with real or
labeled mechanisms); per-dossier required sections 10 → 0; per-experiment-README
heading obligations −2 (current-state, surface-map deleted outright; note-index
rebound, not counted as an obligation anymore); net control-plane line growth: 0
(ratchet held at 4888 throughout). Commits: A1 `6f8e0ecf2`, B1 `5c11cf31d`, B2
`d4b5c96d8`, B3 `f3fbe5e76`, C1 `d7d66e17c`, C2 `a21f45351`, C3 `118481d83`.
