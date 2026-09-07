# Feasibility TEST: applying the `edges:` frontmatter migration to 3 real files

Scope: three COPIES migrated in scratch (`./6m-corpus-as-ecology-seed.md`,
`./never-paraphrase-fields.md`, `./carrier-pull-coupling-as-dual-use-attack-surface.md`).
No file under `Paper/` was touched; no commits made.

## 1. Per-file edge counts

| file | edges written | rel types used |
|---|---:|---|
| `6m-corpus-as-ecology-seed.md` | 9 | home-family, dedup-candidate, grounds-in, relates-to, external-occupancy ×3, spawns ×2 |
| `never-paraphrase-fields.md` | 2 | grounds-in, superseded-by |
| `carrier-pull-coupling-as-dual-use-attack-surface.md` | 19 | depends-on ×3, home-family ×2, distinct-from, follow-up-to, grounds-in, external-occupancy ×8, spawns ×3 |
| **total** | **30** | |

This is noticeably higher than the phase0 probe's own illustrative worked
examples (5, 2, 10 respectively) because I fully mechanized the "Theory
spine" prose chains and the *entire* external-citation list rather than a
representative subset. That difference is itself a finding — see §5.

## 2. Resolution check (every `to` target)

All in-repo paths were checked for literal existence; anchors were checked
against actual headings (not assumed).

| edge (file → rel → to) | check | result |
|---|---|---|
| 6m-corpus → home-family → `paper-families.md#open-ended-ecology` | grep for `` `open-ended-ecology` `` in paper-families.md | **found** (row 24), but **not a real anchor** — it's a table-cell value, no heading exists there. `#open-ended-ecology` would 404 in a real renderer. |
| 6m-corpus → dedup-candidate → `incubator/divergent-directions/11-diversification/research-brief.md` | file exists | **resolves** |
| 6m-corpus → grounds-in → `theory/gaps/body-ecology-coupling.md#candidate-evidence` | heading exists at line 58 | **resolves at file+section level**; does NOT resolve to the specific empty "(future) yes/yes/yes" table row the prose is actually about — a third sub-row granularity gap (see §5) |
| 6m-corpus → relates-to → `.../reaction-norm-grounding.md` | file exists | **resolves** (file-level only; no anchor attempted, see §5) |
| 6m-corpus → external-occupancy → alphastar-league-pfsp / population-based-training / policy-space-response-oracles | by design, external, non-repo | **correctly non-resolving** (0/3, expected) |
| 6m-corpus → spawns → 2 unfiled slugs | by design, no file yet | **correctly non-resolving** (0/2, expected) |
| never-paraphrase-fields → grounds-in → `theory/references/compression-loss-experiment-2026-07-02.md` | file exists | **resolves** |
| never-paraphrase-fields → superseded-by → `incubator/metadata-as-object-graph.md` | file exists | **resolves** |
| carrier-pull-coupling → depends-on → carrier.md / structural-pull.md / delayed-addressability.md | all exist | **resolves** (3/3) |
| carrier-pull-coupling → home-family → `paradigms/ai-self-replication` | dir + README.md exist | **resolves** (dir-addressing convention) |
| carrier-pull-coupling → home-family → `paper-families.md#successor-continuity` | table-cell value found (row 26) | **found but same non-anchor issue as above** |
| carrier-pull-coupling → distinct-from → `candidates/aaai-c3-safety-non-inheritance.md` | file exists (status: retired) | **resolves** — edges can legitimately point at retired objects |
| carrier-pull-coupling → follow-up-to → `.../deep-research-synthesis.md` | file exists; grepped and confirmed it contains the literal "Do not cite MINJA... MemoryGraft" caveat this edge claims to follow up on | **resolves, content-verified** |
| carrier-pull-coupling → grounds-in → `predictions/resolved/contradicted/self-immunity-P1.md` | file exists | **resolves** |
| carrier-pull-coupling → external-occupancy → 8 arXiv/taxonomy ids | by design, external | **correctly non-resolving** (0/8, expected — these are IDs, not repo paths) |
| carrier-pull-coupling → spawns → 3 unfiled slugs | by design, no file yet | **correctly non-resolving** (0/3, expected) |

**Headline:** every in-repo edge I wrote resolves at file level (13/13). Two
distinct sub-file resolution failures recur (table-row anchors in
`paper-families.md`, and a table-row inside `## Candidate evidence`) — both
are the *same class* of problem the phase0 probe already flagged for P12 /
case-study-5, but this test surfaces a **new sub-case**: a target that isn't
even a heading (a bare table cell), which `file#anchor` syntax cannot
address at all, only a grep-style lookup can. `home-family` is the relation
type that hits this every time, because `paper-families.md` rows are never
headings.

## 3. Reciprocal-edge test: `never-paraphrase-fields --superseded-by--> metadata-as-object-graph`

Checked (read-only) whether adding the reciprocal `refines:
never-paraphrase-fields` edge to the real `metadata-as-object-graph.md`
would be consistent:

- `metadata-as-object-graph.md`'s current `links:` frontmatter already
  contains `never-paraphrase-fields` (same direction, already-existing
  untyped citation) — no conflict.
- Its `## Status` is `seed` (not `retired`/`superseded`) — consistent with
  being the surviving, refining side of the pair.
- No other file already claims a `superseded-by`/`refines` edge over either
  file that would conflict.
- `refines` is itself attested elsewhere in the repo
  (`theory/objects/proxy-successor.md`: `refines: [reproduced-object] —
  narrows...`), used in exactly this narrows/corrects sense and in exactly
  this direction (declared on the narrower/corrected object, pointing at
  what it narrows).

**Verdict: consistent, no contradiction or cycle.** The pair
`superseded-by`/`refines` behaves like the repo's other attested mirror
pairs (`conditioned-by`/`conditions`, `enables`/`enabled-by`) — one directed
semantic relationship expressed as two inverse-labeled edges, not a genuine
2-cycle. The one real gap: `superseded-by` itself (as opposed to `refines`)
is **not attested anywhere** in the current repo — it is a new verb,
proposed here only as the mechanical mirror of an already-attested one (see
§5, needed-extension list).

## 4. Derive-don't-declare test (6m-corpus `verdict`)

Attempted to derive `verdict: extends-other-existing-family` from the edges
written on the copy:

- `home-family -> paper-families.md#open-ended-ecology`
- `dedup-candidate -> incubator/divergent-directions/11-diversification/research-brief.md`

The design's own rule (phase0 §7) is: derive `extends-other-existing-family`
from "a `home-family` edge whose target **differs from the file's
originating paradigm**" plus a `dedup-candidate` edge.

**Result: derivation only PARTIALLY reproduces the hand-written verdict —
this is a real design gap, not a stale value.**

- The `dedup-candidate` edge alone is an unambiguous, purely mechanical
  signal that the file overlaps/duplicates prior work and should not be
  independently promoted — that half derives cleanly and would already have
  caught the exact case the hand-written verdict caught.
- The *comparison* half ("home-family differs from the file's originating
  paradigm") has **no operand to compare against anywhere in the graph**. I
  checked three candidate sources for an "originating paradigm" signal and
  found none:
  1. The incubator schema itself (`Paper/incubator/README.md`) — only
     `kind`, `origin` (a date/provenance string, not a paradigm),
     `links`, and (for `idea`) `thesis`. No paradigm field.
  2. The linked theory objects (`successor-morphology.md`,
     `network-utility-selection.md`) — checked their frontmatter directly;
     neither carries any family/paradigm tag. Theory objects are
     family-agnostic by design; `home-family` is exclusively an
     incubator/papers-layer construct.
  3. `paper-families.md` / `registry.tsv` — these register papers and
     candidates, not incubator sparks, so there is no existing "which
     paradigm did this spark originate in" row to read.
- So "differs from the file's originating paradigm" is currently only
  answerable by an agent's out-of-band world-knowledge of what the
  program's "default" paradigm is (here, `ai-self-replication`, inferred
  from context, not from any edge) — which silently reintroduces exactly
  the tacit-knowledge dependency the edge design exists to eliminate.

For contrast, I ran the same derivation attempt on
`carrier-pull-coupling-as-dual-use-attack-surface.md`, where **no**
`dedup-candidate` edge exists and both `home-family` targets ARE the file's
own originating paradigm (it's `ai-self-replication`-adjacent from the
start) — there the derivation *does* land cleanly on "new angle within an
existing family, explicitly distinguished from its nearest sibling via
`distinct-from`," which matches the original hand-written
`extends-ai-self-replication` well. The rule works when origin and
home-family visibly coincide; it breaks exactly when they diverge, which is
the one case ("extends *other*") it was specifically built to detect. That
is a real, not incidental, gap.

## 5. Cost + friction report

**Extraction effort per file:** roughly proportional to file maturity, as
phase0 predicted. `never-paraphrase-fields.md` (a bare seed) took under a
minute to extract — 2 unambiguous edges. `6m-corpus` and
`carrier-pull-coupling` (matured, JSON-verdict-bearing) each took
substantially longer, dominated by (a) deciding which Novelty-note citations
were load-bearing enough to edge vs. leave as color, and (b) resolving
anchor granularity. Most individual relation extractions were unambiguous
(a `links:`-style bare path or an arXiv id is not a judgment call); the
judgment calls clustered in exactly two places: whether a prose citation was
evidentiary (`grounds-in`) vs. merely background, and whether a target
needed an anchor and, if so, which one.

**Top 3 friction / breakage findings:**

1. **`home-family` cannot be anchored, only grepped — a resolution mode the
   design doesn't name.** Every `home-family` edge in both matured files
   points into `paper-families.md`, and every one of them targets a table
   *cell*, not a heading. `file#anchor` syntax (which the design proposes
   for the P12/case-study-5 sub-file case) silently doesn't apply here: a
   real anchor-following reader gets a 404, while a "grep this file for the
   literal family name" reader succeeds. Since `home-family` is one of the
   single highest-traffic relation types (it's on almost every matured
   spark), this is a first-order gap, not an edge case — the design needs a
   named resolution strategy for "target is an unheaded value inside a
   registry file," distinct from both `file` and `file#anchor`.

2. **Full mechanization roughly doubles edge count vs. the design's own
   illustrative examples, concentrated entirely in `external-occupancy`.**
   Converting every Novelty-note citation into an edge took
   `carrier-pull-coupling` from the phase0 sketch's 10 edges to 19; nearly
   all of the growth is 5 additional `external-occupancy` entries the
   worked example simply didn't show. This is cheap per-edge but raises a
   real design question phase0's bloat discriminator (§8) doesn't yet
   answer: does a *list* of 8 near-identical external citations deserve 8
   separate typed edges, or one `external-occupancy` edge to a citation
   block object? As written, the open-`to`-string design has no opinion,
   and the natural instinct (mechanize everything with a target) produces a
   frontmatter block that is longer than some of the prose it replaced.

3. **A real, non-trivial slice of load-bearing content refused to become an
   edge at all — not from ambiguity but from missing identifiers.** Three
   citations in `carrier-pull-coupling`'s Novelty note (a CSA CISO briefing,
   the "Skill Issues" blog series, a Feb-2026 ClawHub malware campaign) name
   real, externally-adjudicable events but have **no canonical id** — no
   arXiv number, no CVE, no stable name to grep against. Inventing a slug
   for them (`clawhub-2026-02-malware-campaign`) just relocates the
   paraphrase risk the whole design exists to eliminate one level down, into
   the slug itself, with no external anchor to check it against. These
   stayed in prose by necessity, not by choice — the design's own dropped
   "counter-thesis" and "kill-condition" precedents (§6/§7: "stays literal
   until an addressable target exists") extend naturally to this case, but
   phase0 didn't anticipate it: unlike a `spawns` slug (self-resolving once
   promoted) or an arXiv id (externally stable), an *unstably-named* external
   event has no promotion path at all.

**Smaller findings, still worth carrying forward:**

- The design's own phase0 worked example for `6m-corpus`
  (`theory/gaps/body-ecology-coupling#evidence-table`) points at an anchor
  that does not exist in the real file (`## Candidate evidence`, not
  `## Evidence table`) — a concrete instance of exactly the drift-at-
  paraphrase problem the whole migration exists to prevent, found inside the
  design document proposing the fix.
- `to` as an open string worked without friction for every case that HAD a
  stable identifier (repo path, `file#anchor`, `arxiv:NNNN.NNNNN`, bare
  slug) — the format ambiguity only ever bit on the two cases above
  (unheaded table-cell targets; identifier-less external events), not on
  the open-string design itself.
- One relation had no clean verb in the reused set at all and had no close
  mirror to borrow from (`follow-up-to` — "this research pass resumes a
  previously-deferred citation thread"); flagged as a genuinely new verb,
  distinct from `superseded-by` (a mirror of an attested verb) and from
  `relates-to` (already attested, used where no better fit existed for the
  6m-corpus terminology cross-reference).

## Verdict

The mechanics hold up for the common case (repo-path and external-id
targets resolve almost for free, matching phase0's ~99% finding). Three
concrete places break or cost more than advertised: (1) `home-family`'s
actual target shape (an un-headed registry row) has no resolution mode in
the current design; (2) full-fidelity external-citation extraction roughly
doubles edge count with no bloat discriminator to cap it; (3) a real class
of load-bearing content (identifier-less external events) has no promotion
path to becoming an edge, unlike the `spawns`/external-id cases the design
already accounts for. The derive-don't-declare mechanism works exactly when
origin and home-family visibly coincide and silently fails to specialize
when they diverge — precisely the case it was built to catch.
