# Idea-Layer Edge Migration — Independent QA + Graph-State Report

Scope of review: full `Paper/incubator/**/*.md` corpus, checked against
`Paper/_design/proposals/idea-layer-migration-runbook.md` and
`idea-layer-object-graph.md`. Today's migration = the 87 files shown as
`M` under `Paper/incubator/` in `git status` (the pre-existing 14-file pilot,
already committed and git-clean, is excluded from "today's batch" numbers but
included in the full-corpus graph analysis for Part B, since Part B needs the
whole 101-file population). Nothing in `Paper/` was modified by this review;
nothing was committed. Scripts used: `qa_parse.py`, `edge_extract.py`,
`check_missed_citations.py` / `_2.py`, `check_reverse.py`, `full_m1.py`,
`tier1_frontier.py` (all in this scratchpad dir).

## Headline verdict

**Every hard invariant holds. Commit-ready on invariants.** One real,
mechanically-verified **completeness gap** (not an invariant violation):
**20 stable-arXiv-id citations across 8 files** were left in prose instead of
promoted to `external-occupancy` edges, in violation of runbook §3.5. This is
cheap to fix (append missing edge lines) and does not by itself block a
commit, but should be either patched or explicitly logged as a known gap
before the human commits.

---

## PART A — Migration QA

### 1. Coverage / skips

- 103 `Paper/incubator/**/*.md` files total; 2 correctly excluded as
  index/README (`Paper/incubator/README.md`,
  `Paper/incubator/divergent-directions/README.md`) — never carried `edges:`,
  never should.
- **101/101 in-scope files carry `edges:`.** Of these, 87 were touched today
  (git `M`), 14 are the pre-existing, already-committed pilot batch.
- **0 empty `edges: []` blocks** anywhere (no bulk stamping).
- 3 files have exactly **1 edge** (the corpus minimum):
  `aaai-c1-exec-unfakeable-writeup.md` (`superseded-by` a candidate — file is
  an explicit "provenance pointer, do not extend," so 1 edge is exactly
  correct), and two `divergent-directions/*/research-brief.md` files
  (`05-signal-design`, `08-succession-complexity`) each with a single
  `external-occupancy` edge. Read both research briefs in full: each really
  does contain only one arXiv/DOI-bearing citation among many
  author-year-only citations (correctly left in prose per §3.5) — **not**
  under-stamped filler, genuinely minimal.
- Edge-count distribution across the corpus: min 1, median 13, mean 12.7,
  max 28 — no suspicious clustering at a "filler" floor.
- **No file that should have been skipped was wrongly stamped** — every file
  checked has load-bearing content behind its edges.

### 2. Resolution

- Full corpus: 1280 edges parsed, **0 non-resolving in-repo edges**. Kind
  breakdown: 786 repo-path, 298 external-arxiv, 173 spawns-slug, 21
  external-doi, 1 external-owasp, 1 `file#anchor` (and that one anchor
  **does** resolve — validated against the target's actual headings). The
  migration is conspicuously anchor-averse (1 anchor in the whole corpus),
  consistent with the runbook's caution to prefer file-level edges.
- Today's 87-file batch alone: 1163 edges, 1003 in-repo (all resolve), 160
  spawns/external (all well-formed by grammar check).
- arXiv-edge sanity: 298 external-arxiv edges resolve to 249 *distinct* ids,
  max reuse count 4 — a natural citation distribution, no sign of
  copy-paste/hallucinated repetition.
- **Reverse check (hallucination guard):** for every external-occupancy edge
  in every file, confirmed the id string also appears somewhere else in that
  same file's own prose (i.e., no edge invents a citation not actually
  discussed). **0 misses** — every external edge traces to real prose.
- **Forward check (completeness gap — the one real defect class):** swept
  every migrated file for arXiv-id-shaped tokens (`NNNN.NNNNN`) in prose and
  diffed against that file's own `external-occupancy` edges. After manually
  eliminating false positives (DOI substrings, e.g. `...5646.2007...` inside
  `doi:10.1111/j.1558-5646.2007.00203.x` are not arXiv ids), **8 files leave
  a genuine, unambiguous, parenthetically-cited arXiv id in prose with no
  corresponding edge** (20 ids total):
  - `discovered-sparks/barrier-order-as-an-attack-surface.md` — missing 2508.17155, 2603.00476
  - `discovered-sparks/bottleneck-degree-as-a-dial-for-fork-experiments.md` — missing 2503.17378, 2604.09588
  - `discovered-sparks/coordination-as-oversight-evasion-safety-relevant-specialization.md` — missing 2410.00031, 2510.04303, 2511.02997, 2601.11369, 2602.15198, 2602.20628, 2603.15655, 2604.01151 (8 ids — this file's Novelty note cites a whole AI-control/collusion literature pool and only captured a subset)
  - `discovered-sparks/persistence-hazard-transfer.md` — missing 2503.17378, 2605.06760
  - `discovered-sparks/recursive-self-similar-goodhart-ladder.md` — missing 2510.02840 ("Take Goodhart Seriously" — sits in the *same sentence* as arXiv:2410.09638, which *did* get an edge; inconsistent treatment within one paragraph)
  - `discovered-sparks/salience-governance-as-a-fourth-resource-control-axis.md` — missing 2603.19469, 2606.10749
  - `discovered-sparks/singular-learning-theory-successor-selection.md` — missing 2512.00686
  - `discovered-sparks/web-scale-stigmergic-eval-contamination-bridge.md` — missing 2508.13180, 2606.05241
  All 20 verified by hand (line-context read) to be real, well-formed arXiv
  ids cited as external literature, not incidental numbers. This is a
  **§3.5 compliance gap** (stable id ⇒ should be an edge), concentrated in
  the "Novelty note" sections of long literature-survey sparks where the
  agent evidently captured the *first* or *most central* citation per
  paragraph but not every one in a dense citation-dump sentence.

### 3. `home-family`

**0 occurrences anywhere in the corpus.** Clean.

### 4. Verb-set conformance

- Distinct `rel` values used corpus-wide: `dedup-candidate, depends-on,
  distinct-from, external-occupancy, follow-up-to, grounds-in, originates-in,
  refines, relates-to, specializes, spawns, superseded-by, tensions-with`.
- `refines` (1 instance, in `metadata-as-object-graph.md`) is **not** in the
  approved verb list given for this QA — but that file is one of the
  **pre-existing 14-pilot files** (git-clean, untouched today), and `refines`
  is explicitly named as an *attested* verb in `idea-layer-object-graph.md`
  itself (§ Relation-type set: "refines / superseded-by (refines attested...)").
  **0 invented verbs in today's 87-file batch.**
- Today's batch uses exactly 8 verbs: `depends-on` (385), `external-occupancy`
  (301), `relates-to` (193), `spawns` (160), `grounds-in` (83),
  `originates-in` (34), `dedup-candidate` (4), `superseded-by` (3) — total
  1163. **Notably, today's batch never once uses `distinct-from`,
  `tensions-with`, `follow-up-to`, or `specializes`** (all 4 uses of those
  verbs corpus-wide are from the pre-existing pilot). Per the runbook this is
  only correct if none of the 87 files' prose explicitly asserted a
  distinct-from/tensions-with/follow-up claim — plausible for a first pass
  over sparks that mostly *build on* rather than *contend with* existing
  material, but worth the human's attention: it means the M1 false-positive
  guard (declared `distinct-from` suppresses a candidate pair, see Part B §8)
  got **zero help** from today's data.

### 5. Additive-only

- `git diff` across all 87 changed files: **0 deleted lines in every file**
  (verified programmatically, not sampled). Pure insertion in every case.
- `edges:` is the **last top-level frontmatter key** in all 87 files (i.e.
  added after existing keys, per runbook §2.2) — verified programmatically,
  0 exceptions.
- Since there are 0 deletions anywhere, `links:` cannot have been touched in
  any file (a modification would show as a paired deletion+addition).
- Scope check: `git status --short` outside `Paper/incubator/` shows only
  pre-existing, unrelated in-flight work (public-branch-redesign branch
  files, `operator/qcli/*`, etc.) — **none** attributable to this migration.
  Exactly 87 files under `Paper/incubator/` are modified, 0 added/deleted/
  renamed.

### 6. YAML validity

- All **87** files touched today parse as valid YAML frontmatter (strict
  PyYAML `safe_load`). **0 failures.**
- Caveat (not a defect of this migration): 4 of the pre-existing 14-pilot
  files (`collision-is-not-automated.md`, `consequence-identity.md`,
  `encoding-recursion-productive-spiral.md`,
  `negentropy-is-not-independence.md`) fail *strict* YAML parsing because
  their (pre-existing, untouched) `thesis`/`description` prose contains an
  unquoted `: ` mid-sentence (e.g. "...collision is exactly that: valued but
  with no surface..."), which is ambiguous to a strict YAML parser. This
  predates today's migration (git-clean, not modified), is unrelated to the
  `edges:` work, and their `edges:` blocks are recoverable by a regex
  extractor regardless. Flag for the human as a latent hygiene item (quote
  those scalars), not something this migration caused or should fix.

### 7. Re-discovery quality spot-check (read in full, not sampled)

Read 4 migrated files in depth, chosen because they cite
selection-signal-collapse / regime-decomposition / dorm-06 — exactly the
contested/nuanced-status nodes named as hard cases:

- **`discovered-sparks/carrier-memetics-mechanism-correction.md`** — 23
  edges. All 7 arXiv ids mentioned in its dense "Novelty note" paragraph
  were captured 1:1, zero hallucination, zero omission. All 12 repo-path
  targets (theory objects, a gap, a paper report, a prediction, a divergent-
  direction brief) verified to exist on disk. Its `depends-on ->
  theory/objects/selection-signal-collapse` edge is judgment-correct: that
  node's own frontmatter records it as *partially* contradicted (the
  exec-dependent-resists half was contradicted by `regime-decomposition-P3`
  on 2026-07-01; the narration-satisfiable-collapse half survives) — the
  file's prose treats it as live, surviving theory, so `depends-on` (not a
  stale affirming edge) is the right call. No silent drift.
- **`discovered-sparks/recursive-self-similar-goodhart-ladder.md`** — 17
  edges, otherwise excellent (correct `depends-on` to
  `selection-signal-collapse`, `exec-barrier`, `narrative-claim`,
  `network-utility-selection`, `encoding-pressure`; correct `relates-to` to
  `papers/unfakeable-selection` and `candidates/iclr-regime-decomposition`).
  This is the file with the one concrete miss noted in §2 above
  (arXiv:2510.02840 uncaptured next to arXiv:2410.09638 which was captured,
  same sentence).
- **`discovered-sparks/closed-sandbox-oee-impossibility.md`** — its
  `dedup-candidate -> candidates/aaai-c1-exec-unfakeable-edge` edge points at
  a candidate file whose own frontmatter still says `status: open` while
  `registry.tsv` says `retired` (2026-07-01) — this is a **pre-existing
  inconsistency in the candidate file itself** (out of migration scope:
  `candidates/*` is explicitly untouchable per the runbook), not migration
  drift: the edge target is still the correct, current, addressable node
  (the file didn't move), and a dedup-candidate claim about a *historical*
  duplication remains true regardless of the target's later retirement.
  Correctly resolved, no fix needed on the migration's part.
- **`discovered-sparks/cultural-attraction-rival-to-the-discrete-replicator-weismann-reading.md`**
  — clean 22-edge file, used as a comparison point for the Part B cluster
  check below.

**Overall quality impression: high.** Resolution discipline is careful
(re-verified against the live repo, not copied from prose), current-node
judgment is correct on the two contested cases actually tested, and no
silently-stale reference was found in the sample. The one real defect
(§2/§3.5 misses) is an *omission* pattern (dense multi-citation sentences
losing a tail id), not a *fabrication* or *drift* pattern.

---

## PART B — Current graph state

### 8. Over-density (specificity-weighted M1, full 101-file corpus)

Reused `m1_tightening.py`'s method verbatim: `score(a,b) = Σ 1/in-degree(shared
target)` over all edges (excluding `home-family`, of which there are none),
threshold ≥ 1.2, run over all 101 files / 1280 edges / 664 distinct targets.

- 2107 pairs share ≥1 target (out of C(101,2)=5050 possible pairs).
- **4 pairs score ≥ 1.2**, forming **3 connected-component clusters**
  touching **7/101 files (~7%)**:

  | Cluster | Size | Members | Score | Shared basis |
  |---|---|---|---|---|
  | A | 3 | `06-phase-diagram-capability`, `14-bureaucratization`, `15-reaction-norm-benchmark` (all `divergent-directions/*/research-brief.md`) | 2.0 / 2.0 | 4 shared arXiv ids each pairwise (Reflexion/Generative-Agents/Voyager/MemGPT for 06↔14; BIG-bench/HELM/ChatbotArena/AgentBench for 06↔15) |
  | B | 2 | `carrier-memetics-mechanism-correction.md`, `cultural-attraction-rival-...-weismann-reading.md` | 1.971 | 3 shared arXiv ids (cultural-evolution-of-LLM-agents lit) + `papers/cultural-evolution-past-brain-world-barrier` + `theory/objects/{carrier,exec-barrier}` + `theory/views/substrate` |
  | C | 2 | `sematectonic-vs-marker-trace-taxonomy.md`, `web-scale-stigmergic-eval-contamination-bridge.md` | 1.525 | `candidates/stigmergic-trace-causality`, a benchmark workspace, `theory/gaps/{carrier-selection,function-vs-narrative}`, `theory/objects/{carrier,stigmergy-as-carrier}` |

  **Sanity-check by reading (as required, not trusting the score):**
  - **Cluster A is a demonstrable false positive.** Read all three briefs:
    they are independent literature reviews over the same ML-agent subfield
    on genuinely different research questions (capability phase diagrams vs.
    organizational bureaucratization vs. reaction-norm benchmarking). The
    shared arXiv ids are generic, well-known background citations (Reflexion,
    Generative Agents, Voyager, MemGPT; BIG-bench, HELM, Chatbot Arena,
    AgentBench) that any two research-briefs surveying this subfield would
    both cite — they got a high specificity weight only because *most other*
    incubator files (biology/evolution-flavored sparks) never cite ML
    benchmark papers, not because these three files share a real claim. This
    is a **new false-positive mode** the 14-file pilot never exercised
    (the pilot didn't include the `divergent-directions` literature-heavy
    research-brief population).
  - **Clusters B and C read as plausible true candidates.** B: both files
    are from the same 2026-07 backlog-maturation pass and both propose
    mechanisms for carrier-mediated cultural/memetic inheritance — genuinely
    adjacent, worth a human collision check (not obviously a duplicate, but
    same neighborhood, sharing real theory-node substance, not just
    literature). C: both explicitly extend the *same* stigmergy/marker-trace
    theoretical toolkit (`stigmergy-as-carrier`, `trace-inheritance`) toward
    adjacent applications — same-referent family, correctly surfaced.

  **Verdict on the detector at full scale:** it still works (2/3 clusters are
  real, non-trivial candidates), but it has a **known, now-demonstrated
  blind spot**: pairs of literature-heavy `research-brief` files can score
  above threshold purely from shared generic background citations. A cheap
  fix worth flagging to the human: discount `external-occupancy` targets in
  the specificity score (or require ≥1 shared *non-external* target) when
  both files are `kind: research-brief`/long literature surveys.

### 9. In-degree distribution

Top hubs (idea-layer in-degree), full 101-file corpus:

| in-degree | target |
|---|---|
| 43 | `theory/objects/carrier` |
| 42 | `theory/objects/exec-barrier` |
| 40 | `paradigms/ai-self-replication` |
| 26 | `theory/objects/selection-signal-collapse` |
| 23 | `theory/gaps/function-vs-narrative` |
| 23 | `theory/objects/self-description-resource` |
| 21 | `theory/objects/delayed-addressability` |
| 15 | `theory/objects/emergent-succession` |
| 15 | `theory/objects/unsampled-function` |
| 14 | `theory/objects/structural-pull` |
| 13 | `theory/objects/successor-morphology` |
| 13 | `theory/objects/self-immunity` |
| 11 | `theory/views/substrate` |
| 11 | `core/mechanisms/unfakeability-gradient` |
| 11 | `papers/structural-elicitation/claim` |

These are exactly the paradigm's own core spine objects — a sane result (the
idea layer's attention concentrates on the theory it's built to extend, not
on noise).

**Tier-1 under-attention frontier** — cross-referenced against all 58
`Paper/theory/objects/*.md` (57 real objects + 1 README, excluded) and all 14
`Paper/theory/gaps/*.md` (71 addressable nodes total):

- **17 nodes with ZERO incoming idea-edges**: `theory/gaps/dossier-direction-genesis`, `theory/gaps/function-vs-safety`, `theory/objects/affordance-null-ecology`, `theory/objects/autopoietic-membrane`, `theory/objects/continuity-substrate`, `theory/objects/divergence-substrate`, `theory/objects/dossier-shaping-dimensions`, `theory/objects/event-driven-cognition`, `theory/objects/evolutionary-recombination`, `theory/objects/graded-axis-shaping`, `theory/objects/metabolic-transition-topology`, `theory/objects/negative-constraint-cut`, `theory/objects/reproductive-code-topology`, `theory/objects/research-state-machine`, `theory/objects/self-containing-runtime`, `theory/objects/sequence-carrier-topology`, `theory/objects/trigger-structure`
- **12 nodes with EXACTLY 1 incoming idea-edge**: `theory/gaps/revisable-time-semantics`, `theory/objects/addressability-debt`, `theory/objects/closure-signal-tradeoff`, `theory/objects/deposition-compression`, `theory/objects/explanation-blind-spot-duality`, `theory/objects/identity-criterion`, `theory/objects/latent-substrate`, `theory/objects/reflexive-substrate-loop`, `theory/objects/relational-existence`, `theory/objects/reproduced-object`, `theory/objects/selector-stack-topology`, `theory/objects/viability-constitutive-safety`
- **Total Tier-1 frontier: 29/71 nodes (41%)** have 0–1 incoming idea-edges —
  a real, actionable "cold node" list now that coverage is full.

### 10. State verdict

**Yes — the migrated graph now carries both functions computably over the
whole corpus, with one caveat on each:**

- **Over-density detection**: real and computable. Produces a small (3-
  cluster, 7-file) but mixed signal: 2 clusters read as genuine same-
  neighborhood candidates on manual check, 1 (the largest, densest-scoring)
  is a demonstrable false positive from shared generic literature between
  independent research-briefs — a full-scale-only failure mode not visible
  in the 14-file pilot. The mechanism is sound but needs the literature-
  discount refinement noted in §8 before its output should be trusted
  un-reviewed.
- **Tier-1 under-attention**: real and clean. 29/71 (41%) existing theory
  objects/gaps are near-invisible to the idea layer — a genuine, directly
  actionable frontier list, exactly the readout the migration was built to
  unlock.
- **Tier-2 caveat still holds, unchanged**: this is still structurally blind
  to genuinely empty regions with no node at all — a relation graph cannot
  surface what has no addressable object, and completing this migration does
  not and cannot change that (per the runbook's own explicit non-claim).

---

## Summary for commit decision

**Hard invariants: all clean.** 0 `home-family`, 0 invented verbs in today's
batch, 0 non-additive files, 0 non-resolving in-repo edges, 87/87 valid YAML,
scope clean (only `Paper/incubator/**` touched).

**One real defect class, small and non-blocking:** 20 stable-arXiv-id
citations across 8 files left in prose without an `external-occupancy` edge
(§3.5 gap) — a completeness miss, not a resolution or invariant failure.
Recommend either a 10-minute mechanical patch (append the 20 missing edge
lines) or logging it as a known follow-up before the human commits.
