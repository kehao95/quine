# Phase 0 Feasibility Probe: Idea Layer → Addressable Object Graph

**Method note:** empirical extraction first (§1–5), design second (§6–9), exactly as
ordered by the brief. No validators or migration steps proposed. No commits made.

## Corpus read (15 files, full text)

Seeds (`Paper/incubator/`): `metadata-as-object-graph.md`, `consequence-identity.md`,
`never-paraphrase-fields.md`, `read-layer-prior.md`, `idea-space-density-not-tier.md`,
`encoding-recursion-productive-spiral.md`, `negentropy-is-not-independence.md`,
`collision-is-not-automated.md`.
View design: `Paper/_design/proposals/collision-structure.md`.
Discovered-sparks (`Paper/incubator/discovered-sparks/`): `6m-corpus-as-ecology-seed.md`,
`carrier-pull-coupling-as-dual-use-attack-surface.md`,
`co-evolving-evaluator-as-the-general-antidote.md`,
`drift-vs-selection-confound-in-single-descendant-lineage-chains.md`,
`exploit-niche-succession.md`.
Precedent checked directly: `scripts/check-theory-maturity.sh`,
`Paper/theory/README.md` (edge-verb glossary), `Paper/theory/objects/self-immunity.md`
(a live `## Relations` section), `Paper/theory/objects/README.md`, `Paper/paper-families.md`,
`Paper/registry.tsv`.

**Load-bearing precedent found that the brief did not flag:** the theory-objects
layer (`Paper/theory/objects/*.md`) **already has** exactly the typed-edge
convention `metadata-as-object-graph.md` argues for — a `## Relations` section
with lines like `` `specializes`: target.md — one-line reason ``.
`Paper/theory/README.md` §"Growing the graph" names an explicit **core verb
set** (`depends-on`, `specializes`, `shares-evidence`, `tensions-with`,
`appears-in-view`, `distinct-from`) plus attested extensions
(`grounds-in`, `instance-of`, `part-of`, `dual-of`, `relates-to`, `instantiates`,
`conditioned-by`/`conditions`, …), explicitly declared **open, not closed**.
This is one layer removed from where the migration is proposed
(`theory/objects/`, not `incubator/`) but it is a working, several-months-old
instance of the exact mechanism the four seed files re-derive from first
principles. The frontmatter design below **reuses this vocabulary** rather than
inventing a parallel one — that is the single strongest empirical grounding
available, and it means the "open value space / small revisable type set"
requirement is not speculative, it is observed running code one directory over.

---

## 1–2. Extracted edges (by file)

Only relations with a real target are listed; free text is omitted. "J" marks a
judgment call (a second extractor could plausibly disagree).

### Seed cluster (`incubator/*.md`)

| source | relation (as it actually appears) | target | J? |
|---|---|---|---|
| metadata-as-object-graph | `links` (untyped, frontmatter) | never-paraphrase-fields, read-layer-prior, consequence-identity, idea-space-density-not-tier | |
| metadata-as-object-graph | "Grounded by" | `theory/references/compression-loss-experiment-2026-07-02` | |
| consequence-identity | `links` | idea-space-density-not-tier, selection-signal-collapse, exec-barrier, narrative-claim | |
| consequence-identity | "rests on the … (P12-grounded) half of" | selection-signal-collapse **§Evidence "P12 attractor" bullet** (not a standalone file — see §3) | J |
| never-paraphrase-fields | `links` | idea-space-density-not-tier, consequence-identity, read-layer-prior, selection-signal-collapse | |
| never-paraphrase-fields | "Evidence" | `theory/references/compression-loss-experiment-2026-07-02` | |
| never-paraphrase-fields | "Revision (2026-07-02): … The corrected form:" | **metadata-as-object-graph** | — flagged below as a **missed edge**: this target is named in the body but is absent from the file's own `links:` frontmatter. Direct proof that free-text cross-reference silently drifts out of sync with the declared edge set — exactly the failure the migration is meant to prevent. |
| read-layer-prior | `links` | idea-space-density-not-tier, never-paraphrase-fields, selection-signal-collapse, structural-pull | |
| read-layer-prior | "Evidence" | `theory/references/compression-loss-experiment-2026-07-02` | |
| read-layer-prior | "the repo's existing … pattern … generalized" | `philosophy/living-system.md` §Repository Ontology | J (citation vs. dependency) |
| collision-structure | "downstream of the seed cluster … those supply the detector signal" (depends-on) | consequence-identity, metadata-as-object-graph, read-layer-prior, encoding-recursion-productive-spiral | |
| collision-structure | "orthogonal to" | `_design/proposals/exec-dependent-idea-lifecycle.md` | |
| collision-structure | "instantiation of the move that … names" | `philosophy/mutual-encoding-network.md` §3 | |
| collision-structure | "treats as transport" | `philosophy/computational-epistemology.md` §Validity | J |
| collision-structure | "Driven by" | `incubator/collision-is-not-automated.md`, `philosophy/structural-pull.md` | |
| collision-structure | "grounds in (surviving half)" | `theory/objects/selection-signal-collapse.md` | |
| collision-structure | "per" (M3 non-goal citation) | `incubator/negentropy-is-not-independence.md` | J |

### Discovered-sparks (matured, JSON-verdict-bearing)

| source | relation | target | J? |
|---|---|---|---|
| 6m-corpus-as-ecology-seed | `links` | body-ecology-coupling (gap), successor-morphology, network-utility-selection, symmetry-trap-phase, selection-signal-collapse, unfakeable-selection | |
| 6m-corpus | `home-family` (JSON `existing_family_name`) | `open-ended-ecology` (row in `paper-families.md`) | |
| 6m-corpus | **dedup-candidate** ("already scoped a materially more rigorous version of the same idea") | `incubator/divergent-directions/11-diversification/research-brief.md` | J, but explicit and confident |
| 6m-corpus | terminology cross-ref | `papers/aaai-facultative-succession/references/reaction-norm-grounding.md` | |
| 6m-corpus | `spawns` ("New sparks this seed surfaced") | 2 named, unfiled ideas (no target file) | — expected unresolved |
| carrier-pull-coupling-as-dual-use-attack-surface | `links` | `theory/gaps/carrier-pull-coupling.md`, `predictions/resolved/contradicted/self-immunity-P1.md`, `theory/gaps/function-vs-narrative.md` | |
| carrier-pull-coupling | theory-spine chain | `theory/objects/carrier.md`, `theory/objects/structural-pull.md`, `theory/objects/delayed-addressability.md`, `paradigms/ai-self-replication/` | |
| carrier-pull-coupling | `home-family` | `ai-self-replication` + `successor-continuity` (paper-families.md row) | |
| carrier-pull-coupling | `distinct-from` ("a different claim … does not overlap") | `candidates/aaai-c3-safety-non-inheritance.md` | J |
| carrier-pull-coupling | follow-up-to ("this research pass is effectively that follow-up") | `papers/agentic-alife-workshop/references/deep-research-synthesis.md` | |
| carrier-pull-coupling | `external-occupancy` | OWASP ASI06, arXiv:2602.15654, arXiv:2512.16962, arXiv:2605.15338, +2 more (external IDs, not repo paths) | |
| carrier-pull-coupling | inline kill-condition ("should be retired … if it does not [replicate]") | no target object exists yet — literal text | — see §7 |
| carrier-pull-coupling | `spawns` | 3 named, unfiled ideas | — expected unresolved |
| co-evolving-evaluator-as-the-general-antidote | `links` | regime-blindness, unsampled-function, selection-signal-collapse, unfakeability-gradient, ai-self-replication | |
| co-evolving-evaluator | `grounds-in` | `experiments/active/digital-red-queen/drq-03-host-lineage-drq` | |
| co-evolving-evaluator | fills-gap-in | `paradigms/regime-blindness/` "case-study-5" (a sub-section, not separately addressable) | J — partial resolution, see §3 |
| co-evolving-evaluator | `external-occupancy` | arXiv:2606.26294, arXiv:2606.16682 | |
| co-evolving-evaluator | `spawns` | 2 named, unfiled ideas | — expected unresolved |
| drift-vs-selection-confound | `links` | `theory/objects/safety-selection-neutrality.md`, `incubator/divergent-directions/deep-cross-read.md` | |
| drift-vs-selection-confound | theory-spine chain | `theory/objects/exec-barrier.md`, `theory/objects/successor-morphology.md`, `theory/gaps/branching-identity.md` | |
| drift-vs-selection-confound | `home-family` | `ai-self-replication` + `successor-continuity` | |
| drift-vs-selection-confound | background-cite | `candidates/aaai-c3-safety-non-inheritance.md`, `incubator/divergent-directions/07-path-dependent-divergence/`, `incubator/divergent-directions/13-mesa-successor/` | |
| drift-vs-selection-confound | `external-occupancy` | arXiv:2506.09401, a 2026-01 bioRxiv LTEE paper, a Springer 2025 paper | |
| drift-vs-selection-confound | `spawns` | 2 named, unfiled ideas | — expected unresolved |
| exploit-niche-succession | `links` | selection-signal-collapse, `candidates/symmetry-trap-phase.md`, `papers/unfakeable-selection/` | |
| exploit-niche-succession | **flags-stale-citation-in** ("selection-signal-collapse.md still cites the pre-retraction reading … has not yet been corrected") | `theory/objects/selection-signal-collapse.md`, resolved against `predictions/resolved/contradicted/regime-decomposition-P3.md` | — see live-status note below |
| exploit-niche-succession | source-hypothesis-of | `incubator/divergent-directions/05-signal-design/research-brief.md` (its own unowned H5) | |
| exploit-niche-succession | `distinct-from` (explicit, re: named cousins) | `candidates/symmetry-trap-phase.md`, `papers/unfakeable-selection/` (same files also in `links` — see §4) | |
| exploit-niche-succession | `external-occupancy` | Tierra/Avida (Ray 1991), arXiv:2406.10162, Krakovna catalogs, 4× RLVR-taxonomy arXiv IDs | |
| exploit-niche-succession | `spawns` | 2 named, unfiled ideas | — expected unresolved |

**Live-status check on the "flags-stale-citation" edge:** I opened
`theory/objects/selection-signal-collapse.md` directly. Its status line and a
`## Backlog Maturation Pass (2026-07-02)` section **already record** the
regime-decomposition-P3 contradiction and narrow the surviving claim to the
P12-grounded half. So the inconsistency this spark flagged has since been
fixed in the target object — a live, in-corpus demonstration of exactly the
failure mode `metadata-as-object-graph.md` predicts (a value fixed at the
target and a claim about it made from a second file go out of sync) and of
the fix (an edge into the target's *current* state, not a restated copy of
it, would never have gone stale in the first place).

## 3. Edge resolution

Separating three categorically different edge kinds (conflating them is itself
a finding — see §6):

| edge kind | count | resolves to a real repo file/object? |
|---|---:|---|
| in-repo edges (the target is meant to be an addressable repo object) | 74 | **73–74 / 74 (~99%)** |
| external-occupancy edges (target is an arXiv ID / named external taxonomy, by design outside the repo) | ~17 | 0/17 resolve as repo files — **correctly so**, this is a different node type (external reference), not a failure |
| `spawns` edges (target is a not-yet-filed new idea) | ~13 | 0/13 resolve — **expected**, exactly the "promotion threshold" case `metadata-as-object-graph.md` §Open questions already names |

The one partial/non-resolving in-repo case: `consequence-identity.md`'s "P12"
and `co-evolving-evaluator`'s "case-study-5" both point at a **sub-section
inside** an existing file (an `## Evidence` bullet, a named case study), not
at the file itself. A file-level edge under-resolves these — they need a
two-hop or anchor-qualified target (`selection-signal-collapse#p12-attractor`),
which the design in §6 accounts for rather than papering over.

**Headline finding: resolution is not the bottleneck.** Authors (agent or
human) already write specific enough prose — literal file paths, literal
object ids, literal arXiv numbers — that a careful extractor resolves nearly
every intended in-repo edge on the first pass. The corpus is already
graph-shaped; it just isn't machine-readable as a graph yet.

## 4. Same-referent detection (M1 by hand)

Pairwise shared-target check across the 5 discovered-sparks (10 possible
pairs), using the *unweighted* rule from `collision-structure.md` M1 ("≥2
shared target objects"):

| pair | shared targets | count | verdict |
|---|---|---:|---|
| 6m-corpus × exploit-niche-succession | selection-signal-collapse, symmetry-trap-phase, unfakeable-selection | 3 | **false positive** — body-corpus reuse vs. evaluator-exploit-ordering are different problems; all three shared targets are the paradigm's generic D4/"digital-red-queen" hub objects |
| carrier-pull-coupling × drift-vs-selection-confound | ai-self-replication, successor-continuity, aaai-c3-safety-non-inheritance | 3 | **false positive** — dual-use/security framing vs. population-genetics confound; shared targets are (a) the near-universal home-family and (b) one candidate file both cite *for opposite reasons* (one says "distinct from this", the other cites it as background) |
| co-evolving-evaluator × exploit-niche-succession | selection-signal-collapse, (ai-self-replication implied via home-family) | 2 (borderline) | **false positive**, same hub-object cause |
| co-evolving-evaluator × drift-vs-selection-confound | selection-signal-collapse (body-only), ai-self-replication | 2 (borderline) | **false positive** |
| drift-vs-selection-confound × exploit-niche-succession | selection-signal-collapse (body-only), ai-self-replication (implied) | 2 (borderline) | **false positive** |
| all other pairs | 0–1 | — | not flagged |

**Zero true-positive same-referent pairs among these 5 discovered-sparks
themselves.** That is a real, reportable null result for this specific slice
— these five sparks are, in fact, about five different problems.

But real same-referent pairs **do** exist one level up in the same
maturation-pass output, and they are exactly the kind the theory predicts:

- **6m-corpus-as-ecology-seed ↔ `divergent-directions/11-diversification/research-brief.md`** — a genuine dedup candidate, already caught by the human/agent author in prose ("this is best read as a cost-engineering refinement of that brief… should be filed as a cross-reference note… to avoid recreating near-duplicate-candidate proliferation"). Shared signal: a **narrow** gap target (`body-ecology-coupling`) plus a shared, non-generic external-citation set (PBT/AlphaStar/PSRO/POET), not a paradigm-level hub.
- **carrier-pull-coupling-as-dual-use-attack-surface ↔ `papers/agentic-alife-workshop/references/deep-research-synthesis.md`** — an explicit resumed-thread ("this research pass is effectively that follow-up") to a previously-deferred citation question in another file.
- The 5-file **seed cluster** (metadata-as-object-graph / consequence-identity / never-paraphrase-fields / read-layer-prior / collision-structure) is a connected component with heavy mutual shared-edges, but it correctly resolves to a **refinement chain**, not a dedup: `never-paraphrase-fields` names its own successor as `metadata-as-object-graph` ("Revision … The corrected form:"). This is a live example of collision-structure's M3 "productive disagreement / refines" outcome, distinct from both the dedup and the null cases above.

**Diagnosis of the false-positive mechanism:** every false positive above is
caused by one of two high-in-degree edge kinds — `home-family` (near-universal
within one paradigm) and a small set of paradigm-hub theory objects
(`selection-signal-collapse`, `symmetry-trap-phase`, `unfakeable-selection`,
`ai-self-replication` itself) that most D4-adjacent sparks cite as shared
background. Both true positives above share **narrow, low-in-degree** targets
(a specific gap file, a specific reference file, an explicit named-successor
edge) instead. Raw shared-target counting without an in-degree/specificity
discount will systematically over-trigger on paradigm structure; discounting
it recovers exactly the two real pairs and drops all five false ones in this
sample.

## 5. Numbers

- **Edges per file:** min 5 (bare incubator seeds), max ~20 (fully-matured
  discovered-sparks with Evidence/Novelty/Paradigm-fit sections), mean ≈ 10.2.
  Maturity level of the *source* strongly predicts edge density.
- **Resolve rate:** ~99% for edges the author intended as in-repo (73–74/74);
  0% for external-occupancy edges (correctly — different node type); 0% for
  `spawns` edges (correctly — pre-promotion stubs).
- **Same-referent pairs found:** 2 true positives (both cross-corpus, not
  within the 5-discovered-spark sample), 1 refinement-chain (5-file cluster,
  not a dedup), 5 false-positive pairs at the naive ≥2-shared-target
  threshold, all attributable to two identifiable hub-edge kinds.
- **Apparent false-positive rate (naive rule, within-sample):** 5 flagged
  pairs / 5 flagged pairs = **100% false positive** in this narrow slice;
  but the true positives that motivated the whole design (found via
  human/agent prose judgment, not yet mechanized) are real and exist in the
  wider corpus, and they are structurally distinguishable from the noise by
  edge specificity, not by anything narrative.

---

## 6. Proposed frontmatter design (grounded only in what appeared above)

### Relation types actually observed, with node type and consumer

| relation | observed instance(s) | target node type | consumer (§8) |
|---|---|---|---|
| `grounds-in` | metadata-as-object-graph/never-paraphrase-fields/read-layer-prior → compression-loss-experiment; co-evolving-evaluator → drq-03 | run / experiment / reference file | derive-don't-declare evidence-maturity (generalized `check-theory-maturity.sh`); contradiction-propagation |
| `home-family` | 6m-corpus → open-ended-ecology; carrier-pull-coupling, drift-vs-selection-confound → ai-self-replication/successor-continuity | family row (`paper-families.md`) / paradigm dir | density/clustering view (`idea-space-density-not-tier`); paradigm-fit verdict derivation. **Must be excluded/down-weighted in M1's shared-edge computation (§4) — it is the single largest source of false-positive same-referent flags observed.** |
| `distinct-from` (already a core verb in `theory/README.md`) | carrier-pull-coupling → aaai-c3-safety-non-inheritance; exploit-niche-succession → symmetry-trap-phase/unfakeable-selection | any idea/candidate/theory object | M1 false-positive guard: an explicit `distinct-from` edge should pre-seed `collided-null` for that pair (collision-structure M3), suppressing it from ever resurfacing |
| `refines` / `superseded-by` (dialectical; already implied in `theory/README.md`'s spine-alignment verbs) | never-paraphrase-fields → metadata-as-object-graph | another idea object | contradiction-propagation; M3's "productive disagreement" resolution type in collision-structure.md names exactly this pair |
| `external-occupancy` | every discovered-spark's Novelty-note WebSearch citations (arXiv IDs, OWASP taxonomy names) | an external id/URL string, **not** a repo path | novelty/occupancy view — directly replaces the free-prose "Novelty note" section; feeds the derived paradigm-fit verdict (§7) |
| `spawns` | "## New sparks this seed surfaced" — present in **all 5** sampled discovered-sparks | a name/slug with no file yet (legitimately unresolved) | future-task-surface backlog view (the same mechanism `collision-structure.md` M2 already names) |
| `dedup-candidate` / `follow-up-to` (ad hoc names for the two true positives found) | 6m-corpus → 11-diversification brief; carrier-pull-coupling → deep-research-synthesis.md | another idea/brief/reference file | M1 (once mechanized, this *is* what M1's readout should have flagged automatically) |
| existing theory/objects core+attested verbs (`depends-on`, `specializes`, `part-of`, `dual-of`, `tensions-with`, `appears-in-view`, `instance-of`, `instantiates`, `conditions`/`conditioned-by`) | every "Theory spine" prose chain in the discovered-sparks, and every "Relationship to sibling designs" block in collision-structure.md, is functionally this exact vocabulary written as an arrow-chain instead of typed edges | theory object / gap / candidate | reuse verbatim — **do not invent a second vocabulary for the idea layer**; this is the strongest single grounding fact in the probe |

Two dimensions from the brief's own prompt list were **not** observed in the
sampled corpus and are therefore **not** proposed:

- **`counter-thesis`** (an opposite-signed hypothesis) never occurred as such.
  What occurred instead, twice, was the weaker `distinct-from` ("a different
  claim, does not overlap") — already in the existing verb set, not a new one.
- **`kill-condition` → a condition object**: no addressable "condition" kind
  exists anywhere in the repo. Every kill-condition/falsifier observed
  (6m-corpus's "falsified if all bodies converge…"; carrier-pull-coupling's
  "should be retired… if it does not [replicate]") is inline prose with **no
  prediction file yet**. These stay literal until a prediction file is filed,
  at which point the edge is simply `grounds-in`/`tested-by` → that prediction
  — no separate "condition" object type is warranted (would fail the bloat
  discriminator in §8: no distinct consumer from the prediction file itself).

### YAML shape

Idea-layer files (`incubator/*.md`, `incubator/discovered-sparks/*.md`) keep
their existing scalar frontmatter (`kind`, `name`/`thesis`, `description`,
`origin`) and the existing untyped `links:` list (kept, unchanged, as a cheap
fallback / see-also — it already has a real consumer, generic reading, and
costs nothing to keep). The addition is a typed `edges:` list, one entry per
extracted relation, each `{rel, to}`. `rel` is drawn from the small, open,
already-precedented verb set above (reusing `theory/README.md`'s list rather
than inventing one); `to` is an **open** string — a repo-relative path/id, a
`file#section` anchor for sub-file targets (the P12 / case-study-5 problem in
§3), an external id (`arxiv:2602.15654`), or a bare slug with no file yet
(a `spawns` stub). No enum is placed on `to` anywhere.

```yaml
---
kind: spark
name: 6M-corpus-as-ecology-seed
description: ...
origin: surfaced during 2026-07 backlog-maturation pass
links: [body-ecology-coupling, successor-morphology, network-utility-selection,
        symmetry-trap-phase, selection-signal-collapse, unfakeable-selection]
edges:
  - { rel: home-family,      to: paper-families/open-ended-ecology }
  - { rel: dedup-candidate,  to: incubator/divergent-directions/11-diversification/research-brief }
  - { rel: grounds-in,       to: theory/gaps/body-ecology-coupling#evidence-table }
  - { rel: spawns,           to: ecological-common-garden-vs-developmental-clone-split }   # no file yet
  - { rel: spawns,           to: pfsp-matchmaking-as-symmetry-trap-exit }                  # no file yet
---
```

### Worked BEFORE / AFTER

**1. `never-paraphrase-fields.md`** — before: the successor relation lives
only in a body sentence ("The corrected form: `metadata-as-object-graph`")
and is **absent from `links:`** (the missed-edge finding in §2). After:

```yaml
links: [idea-space-density-not-tier, consequence-identity, read-layer-prior, selection-signal-collapse]
edges:
  - { rel: grounds-in,      to: theory/references/compression-loss-experiment-2026-07-02 }
  - { rel: superseded-by,   to: incubator/metadata-as-object-graph }
```
The reciprocal (`metadata-as-object-graph` gets `refines: never-paraphrase-fields`)
is added on the target per `theory/README.md`'s reciprocal-edge discipline —
this is exactly the check `check-theory-maturity.sh`'s sibling,
`check-theory-coherence.sh`, already performs for theory objects.

**2. `6m-corpus-as-ecology-seed.md`** — before: a 5-paragraph "Novelty note"
prose block plus a fenced JSON blob with `verdict`/`existing_family_name`/
`rationale` keys. After (the `rationale` prose stays — it is genuine argument,
not a value to scalarize — but the machine-checkable conclusion becomes edges):

```yaml
edges:
  - { rel: home-family,        to: paper-families/open-ended-ecology }
  - { rel: dedup-candidate,    to: incubator/divergent-directions/11-diversification/research-brief }
  - { rel: grounds-in,         to: theory/gaps/body-ecology-coupling }
  - { rel: external-occupancy, to: "AlphaStar-league/PFSP" }
  - { rel: external-occupancy, to: "population-based-training" }
  - { rel: external-occupancy, to: "policy-space-response-oracles" }
  - { rel: spawns,             to: ecological-common-garden-vs-developmental-clone-split }
  - { rel: spawns,             to: pfsp-matchmaking-as-symmetry-trap-exit }
```
`verdict: extends-other-existing-family` is no longer declared: it is
**derived** (see §7) from the presence of a `home-family` edge whose target
differs from the file's originating paradigm plus a `dedup-candidate` edge —
so a stale/self-serving verdict (the exact P3-style drift found live in §2)
becomes structurally impossible to leave uncorrected.

**3. `carrier-pull-coupling-as-dual-use-attack-surface.md`** — before: prose
kill-condition buried in "Evidence status" ("should be retired… if it does
not [replicate]"), a distinct-from claim buried in "Novelty note", 6 external
citations only in prose. After:

```yaml
edges:
  - { rel: home-family,        to: paper-families/successor-continuity }
  - { rel: distinct-from,      to: candidates/aaai-c3-safety-non-inheritance }
  - { rel: follow-up-to,       to: papers/agentic-alife-workshop/references/deep-research-synthesis }
  - { rel: grounds-in,         to: predictions/resolved/contradicted/self-immunity-P1 }
  - { rel: external-occupancy, to: "arxiv:2602.15654" }   # Zombie Agents
  - { rel: external-occupancy, to: "arxiv:2512.16962" }   # MemoryGraft
  - { rel: external-occupancy, to: "owasp-asi06" }
  - { rel: spawns,             to: self-immunity-P1-F2-as-unlabeled-poisoning-trial }
  - { rel: spawns,             to: competence-generalization-as-poisoning-discriminator }
  - { rel: spawns,             to: quine-as-live-red-team-target }
```
The kill-condition itself ("retire if the re-score does not replicate the
separation") stays **literal prose** — there is no prediction file yet to
point at (§7 classification below), and per §8 a `kill-condition:` field with
no addressable target has no consumer beyond a human re-reading the same
sentence.

## 7. Classification: EDGE / LITERAL / DERIVED

- **EDGE**: home-family, grounds-in, distinct-from, refines/superseded-by,
  external-occupancy, spawns, dedup-candidate/follow-up-to, and the full
  reused theory/objects verb set (depends-on, specializes, part-of, dual-of,
  tensions-with, appears-in-view, instance-of, instantiates). All observed
  directly in the corpus (§2).
- **LITERAL (stays a scalar)**: `origin` (a date + one-line provenance string
  — already correctly literal in every sampled file); a kill-condition or
  falsifier **before** a prediction file exists for it (no target to edge to
  yet — this is a real, observed, transitional state, not a design gap); the
  `rationale` prose inside a paradigm-fit judgment (genuine argument, not a
  value — compressing it would recreate the exact sign-flip-at-paraphrase
  failure `never-paraphrase-fields.md` documents).
- **DERIVED (computed from edges, never stored)**: the JSON blob's `verdict`
  key (derive from home-family + external-occupancy + distinct-from edges,
  generalizing `check-theory-maturity.sh`'s pattern exactly as
  `metadata-as-object-graph.md` §"Derive, don't declare" specifies);
  evidence-maturity (derive from whether `grounds-in` edges resolve to a run/
  prediction vs. only an external reference vs. nothing — in this sample
  **all 5** discovered-sparks would derive to "literature-grounded" or
  "data-available-unscored," **none** to "run-backed," which is itself a
  finding about the current population's evidentiary state, not a design
  artifact); a `superseded` flag on old seeds (derive from having an
  incoming `refines` edge, rather than needing a hand-maintained "Status:
  matured/retired" line, which the corpus currently maintains by hand and
  which the P3/selection-signal-collapse episode in §2 shows can silently
  desync).

## 8. Bloat discriminator (consumer check)

Every relation type proposed in §6 has a named consumer already sketched
above (mostly `collision-structure.md`'s own M1/M2/M3, plus
`check-theory-maturity.sh`'s generalization, plus the density/clustering view
`idea-space-density-not-tier.md` and `read-layer-prior.md` already call for).
Two candidates considered and **dropped**:

- A distinct `condition:` object type for kill-conditions — no consumer
  different from the prediction file it would eventually point at (§6/§7).
- A distinct `counter-thesis:` relation — never observed; would duplicate
  `distinct-from`, which already exists and already has a consumer.

## 9. Feasibility verdict

**Conditional GO, with one concrete required refinement to `collision-structure.md` M1.**

- **Extraction and resolution are not the bottleneck.** ~99% of in-repo-
  intended edges resolved to a real file on a straightforward reading pass;
  the corpus already writes precise enough prose (literal paths, literal
  object ids, literal arXiv numbers) that almost nothing needs guessing.
  External-occupancy and `spawns` edges correctly fail to resolve to repo
  files because they are, by design, pointing at something else (the outside
  world, or a not-yet-created stub) — that is a category distinction the
  design must carry (`to` accepts external ids and unfiled slugs), not a
  resolution failure.
- **The same-referent signal is real but not densely reliable at the naive
  threshold.** Two genuine same-referent pairs were found (§4), both already
  caught by human/agent prose judgment and both showing the exact
  shared-narrow-edge signature the theory predicts. But in the specific
  5-discovered-spark sample, the naive "≥2 shared targets" rule produced
  **5 false positives and 0 true positives**, because `home-family` and a
  handful of paradigm-hub theory objects (`selection-signal-collapse`,
  `symmetry-trap-phase`, `unfakeable-selection`, the paradigm object itself)
  are cited by nearly everything in one paradigm and swamp the signal.
- **This is not the finding that stalls the direction** — it is a
  parameter, not a wall: excluding `home-family` and down-weighting
  high-in-degree theory-object targets from the shared-edge computation
  recovers exactly the two real pairs and eliminates all five false ones in
  this sample. Concretely: `collision-structure.md`'s Phase-0 gate ("do
  enough idea files carry resolvable edges for M1 to compute a non-trivial
  candidate list?") should be answered **yes, provided M1's shared-edge
  count excludes `home-family` edges and applies an in-degree discount
  (e.g., a target cited by more than ~3 idea files in the same paradigm
  does not count toward the ≥2 threshold)** — a one-line addition to that
  proposal's detector spec, not a new mechanism.
- A secondary, smaller finding worth carrying forward: the "Theory spine"
  prose chains in every discovered-spark and the "Relationship to sibling
  designs" block in `collision-structure.md` are *already* using the
  theory-objects layer's typed-verb vocabulary informally, in prose. The
  migration this probe was asked to evaluate is less "invent a new schema"
  and more "promote a vocabulary the repo has already been running for
  months, one directory up, into the layer that currently only writes it as
  free text."
