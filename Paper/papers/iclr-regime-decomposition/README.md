---
paper_id: iclr-regime-decomposition
title: "Decomposing at the Wrong Level: Structural Blind Spots in AI Self-Replication Evaluation"
claim_bundle: regime-decomposition
active: true
venue: iclr-2027
format: iclr-9
format-envelope: iclr-9-single-column
track: theory
cluster: successor-continuity
state_mode: full-manuscript-home
publication_class: framework-follow-on
maturity_target: framework
status: drafting
role: evaluation-framework-follow-on
candidate: iclr-regime-decomposition
companion: aaai-facultative-succession, iclr-safeguard-composition
family: successor-continuity
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
view-affinity: [continuity]
theory-objects: [exec-barrier, carrier, delayed-addressability, unsampled-function, viability-as-distortion-measure, carrier-membership-heritability-predictor, narrative-claim, viability-pressure, future-task-surface, reproduced-object]
edges:
  - supersedes: archive/superseded-2026/evaluation-paper
---

# ICLR-1 Dossier — Decomposing at the Wrong Level

> **DRAFTING (2026-06-30).** A drafting workspace now exists: locked structural
> decisions ([`plans/reconciliation-and-decisions-2026-06-30.md`](plans/reconciliation-and-decisions-2026-06-30.md)),
> a full [`manuscript-plan.md`](manuscript-plan.md), finalized §3 theorem
> write-ups ([`plans/theorems-formalization-2026-06-30.md`](plans/theorems-formalization-2026-06-30.md)),
> a DRAFT [`references.bib`](references.bib) + source-check pool
> ([`references/`](references)), and section drafts under [`sections/`](sections).
> Build with `make pdf`. Navigation: [`WORKSPACE.md`](WORKSPACE.md).
>
> **Title** set to the working default (overridable). **Status**: section drafts
> generated and adversarially reviewed; pending a human read-through, a
> source-checked bibliography, and (optionally) the unrun 6-condition study.

## Experiment Feeds

- `none-yet` - none - not-for-paper-yet - No active experiment currently declares this framework paper as a feed target.

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - No active experiment feed currently owns a citable backbone entry for this framework paper.

## Development status (2026-06-30)

This pass moved the paper from planning to a drafting workspace and resolved the
inconsistencies that had accumulated across the candidate, this dossier, the
paradigm publication-strategy, and the prediction files. The binding resolutions
(full rationale in the reconciliation note):

- **Six dimensions D1–D6** as a *reference frame* (Framing A), not five and not a
  Kuhnian paradigm-shift. The candidate's "five dimensions" wording is stale.
- **Owned theorems: the Compression Blindness Principle + the Selection-Signal
  Collapse Theorem** only. The boundary-relative safeguard-composition result
  and guard-placement assay belong to the companion
  `iclr-safeguard-composition`; they are **not** claimed here. The shared
  non-identifiability spine may be cited without re-deriving the companion's
  safety result.
- **Evidence = hybrid**: demonstrate via existing waves (no new experiment
  required for the core claims); the 6-condition full-inversion is the *proposed*
  discriminating test.
- **Predictions owned: P1, P3, P4, P5, P6, P7**; P2 → competent-stranger.
  (Prediction file `regime-decomposition-P7` mislabels its dimension "D7"; it is a
  D2 consequence. That file is `scope: paradigm` / human-directed — flagged
  (S-2), proposal only.)

## Adjudication (2026-07-26) — **fact survives; the superseded taxonomy framing does not**

> **STATUS: ACTIVE. Zero class-A findings against the core fact.** What was
> refuted is the *old* taxonomy claim this dossier already replaced — the new
> explanation framing below is the surviving position, not the casualty. One date
> error is load-bearing. All documents were fetched as PDFs and grepped directly,
> which is why this line came through the 2026-07-26 pass intact while others did
> not: it was adjudicated against primary sources from the start. Cross-check:
> [`../../../theory/references/regime-blind-collision-audit-2026-07-26.md`](../../../theory/references/regime-blind-collision-audit-2026-07-26.md).
>
> **Confirmed:** FSF v2.0 removed the Autonomy risk domain (footnote 9 — the
> whole **domain**, not one CCL); FSF 3.1 has `autonom` = 0 and `replicat` = 0;
> OpenAI PF v2 demoted ARA to a Research Category exactly as described; RSP v3.4
> has four threshold rows, none about replication; AISI measured <5% → >60%.
> The umbrella claim — all three labs inside the window — holds.
>
> **Corrections that must be made before use:**
> - **The Anthropic demotion is RSP v2.0, effective `2024-10-15`, not v3.0
>   `2026-02-24`.** The quote is real but sits in the v2.0 changelog block, and
>   the identical text already appears in v2.2 (2025-05-14), nine months before
>   v3.0 existed. The true story is **two steps**: demote (v2.0), then **silently
>   delete the checkpoint in v3.0 with no reason given** — and that silent
>   deletion is the one move where no explanation was offered, which is more
>   useful to this paper than the collapsed version.
> - The FSF zero dates from **3.0 (2025-09-22)**, not 3.1.
> - Restore AISI's dropped qualifiers: "across **11 of our 20** RepliBench
>   evaluations" and "does not **necessarily** mean".
>
> **The explanation is refuted.** The labs stated their reasons and mostly did
> not say "unmeasurable". DeepMind: the risk is "**now covered by our
> misalignment section**" and "**no heightened deployment mitigations would be
> necessary**". Anthropic: "we **now believe that these capabilities — at the
> levels we initially considered — would not necessitate the ASL-3 standard**",
> plus an admission the original threshold was a placeholder. Scorecard:
> absorption/supersession 3/3; severity re-judgement 2/3 explicit; measurability
> 1/3 (OpenAI only, subordinated to threat-model maturity); regulatory pressure
> 0/3. **So the sentence "the construct did not become less dangerous" is
> contradicted by the primary sources this paper cites** — keeping it requires
> arguing the labs were wrong, which needs evidence this dossier states it does
> not have.
>
> **Self-inflicted:** the paper cannot use ">60%" as proof danger rose *and* the
> AISI caveat as proof of unmeasurability — that caveat **is** the labs' argument.
>
> **The claim that survives, stated precisely:** across `2024-10` to `2026-04`,
> replication ceased to be a *named, mitigation-triggering* threshold in all three
> frameworks, while AISI's RepliBench scores rose from <5% to >60%; the labs'
> stated reasons are absorption and severity re-judgement, **not** measurement —
> with one exception, which is where the paper actually lives (below).
>
> **The better hook, primary-sourced and apparently unremarked.** OpenAI PF v2's
> criterion 2 for admitting a capability to Tracked status reads: "**Measurable:**
> We can construct or adopt capability evaluations that measure capabilities that
> **closely track the potential for the severe harm**." That is not a
> measurability criterion — it is a **construct-validity criterion written into
> operative governance**, and ARA is the capability that failed it and was parked
> in a research holding pen with a named escalation path. Beside it sits Agent
> Matrix's dissociation (two models at 100% replication success, 100% vs 0%
> overuse; "**replication success rate is a poor proxy for actual risk**") as
> independent convergence. **The six-dimension frame earns its place only by
> saying which dimension separates those two models and why a metric collapsed
> onto it cannot satisfy criterion 2.**
>
> **Contrary signal to weigh:** METR's *Common Elements of Frontier AI Safety
> Policies* (2025-12-09) compared all three frameworks and **did not flag a
> retreat**, still listing autonomous replication as an example threshold. That is
> both a novelty signal and a warning that the retreat reading is contestable.
>
> **New record error found:** Table 1 names "Agent Matrix / **ERABBench**" —
> **"ERABBench" does not appear anywhere in arXiv 2509.25302.** Delete it. Also:
> the MISR arXiv id (2412.03904) is confirmed but its **NeurIPS 2024 attribution
> is unverified** — the arXiv comments field carries no venue.
>
> **Venue:** as an explanation-of-policy paper this is AIES/FAccT, not ICLR — and
> this dossier has itself demoted CBP to "a framing device, not a theorem", which
> removes the reason an ICLR theory reviewer would engage. The one ICLR-viable
> route is to lead with a **measurement**: a scoring protocol that discriminates
> along the six dimensions, **run on at least two existing open harnesses**
> (RepliBench and Palisade both qualify), with the policy retreat demoted to
> motivation.

## The repositioning (corrections above apply)

**This paper should be rebuilt on an external anomaly it does not currently
record.** An assessment pass found the Claim below is a taxonomy claim whose
central empirical premise is factually wrong on its load-bearing row, and found
a verified fact about the world that is a far stronger opening. The sections
below have **not** been rewritten.

**The anomaly (verified against primary sources):** between 2024-05 and 2026-02
all three frontier labs deleted or demoted their explicit autonomous-replication
threshold — DeepMind FSF v2.0 (2025-02-04) removed the Autonomy level 1 CCL and
the string "autonom" does not appear in FSF 3.1 (2026-04-17); OpenAI
Preparedness Framework v2 (2025-04-15) demoted Autonomous Replication and
Adaptation to a Research Category; Anthropic RSP v3.0 (2026-02-24) replaced the
ARA threshold with a checkpoint, and v3.4 (2026-07-08) has four threshold rows,
none about replication. Over the same period UK AISI measured self-replication
capability rising from <5% to >60% while cautioning that a perfect RepliBench
score does not mean a model could self-replicate.

**New position, one sentence:**

> Measured self-replication capability rose sharply while every frontier lab
> removed self-replication as a governable threshold. The construct did not
> become less dangerous; it became unmeasurable in a way that could carry a
> threshold. This paper says why — the six dimensions along which replication
> outcomes vary — and which dimension each retreat tripped on.

This converts a taxonomy paper into an explanation paper, and answers the
deflation the taxonomy framing cannot ("taxonomies are cheap; show me one
measurement that changes"). It also makes §1's current sentence — that the three
labs "all carry autonomous-replication or self-improvement thresholds," cited to
the 2023 RSP and FSF v1.0 — **false as written and requiring replacement rather
than patching**.

**Table 1 must be rebuilt from primary sources.** The claim that RepliBench
"fixes source access to a black-box API" is wrong: RepliBench (arXiv 2504.18565)
has a named **Obtaining Weights** domain including Covert Weight Exfiltration,
and a **Recursive Replication** family spanning three generations with a
local-weights or API-only choice, plus a Rate-of-Infection family. Similar errors
on SOCK (graded R-score plus RCL/PCL ladders, not a pass/fail bit) and Pan et al.
(whose abstract itself reports a "chain of replica"). Palisade arXiv 2605.06760
is the strongest counterexample to two of the four claimed gaps and is absent
from Table 1. **The accurate criticism is narrower and sharper than the false
one:** the channels are present in these harnesses and are never *scored
comparatively* — which is a claim about scoring, and scoring is exactly this
paper's thesis.

**Citations that flip from refuter to foundation:**

| Work | Old role | New role |
|---|---|---|
| **Agent Matrix, arXiv 2509.25302** | cited only as a benchmark; actually publishes the construct-validity critique in this exact domain, with a demonstrated dissociation (two models at 100% replication success, 100% vs 0% overuse) | **co-critic and independent convergence.** Cite as corroboration that success rate ≠ risk, then show the six-dimension account explains *which* dimension their dissociation lies on |
| Srivastava arXiv 2602.16984; Wang arXiv 2606.05169 | reported as having "beaten CBP to the formal claim" — **verified 2026-07-26 and that reading is withdrawn.** Srivastava is *Fundamental Limits of Black-Box Safety Evaluation*, about **latent context-conditioned policies** (rare-under-eval triggers); Wang is *The Evaluation Blind Spot*, a **stereological theory of benchmark-suite coverage**. Different objects; neither decomposes a measured behaviour into regimes | **adjacent formal neighbours to cite as company, not competitors.** The disposition is unchanged for an independent reason: CBP is the data-processing inequality and should be a framing device, not a claimed theorem |
| Bean et al., arXiv 2511.04703 (NeurIPS 2025 D&B) | establishes the genus at a top venue | **precedent that this kind of paper publishes**, and the general frame this specializes |
| Hendrycks arXiv 2303.16200; costly-signalling literature (Spence, Zahavi, Grafen) | uncited; every reviewer asks | required related work — the SSC pooling/separating structure is signalling theory, and the repo's own `references/collision-implications.md` already names pooling/separating equilibria |

**Honesty assets to keep and foreground.** `references.bib` self-flags 111 of
112 entries as unverified, and a spot-check of twelve load-bearing 2026-dated ids
resolved to real papers with matching content. The dossier records its own
falsified prediction (P3, contradicted by its authors' own matched-pair test) and
its unrun study. Under an explanation framing, the P3 falsification is
reportable as a finding about the instrument's discriminating power rather than a
wound.

**What this repositioning does not fix.** The 86/86 result remains confounded —
all conditions used the same peer-wake scalar, and the same experiment family
concluded outcomes are seed-driven and stochastic. The "~67% host / ~14.5%
agent" figures come from re-reads of C84 and C81 specifically and must not be
generalized to "86 prior conditions." D3 has no evidence at all. The MISR
footnote ("we located no primary publication") must be fixed — it is arXiv
2412.03904, NeurIPS 2024. Body text is ~25% over the 9-page envelope.

## Claim

AI self-replication evaluation decomposes at the wrong level: existing
benchmarks decompose operational capability and risk, but compress
heritability channel, criterion content, and ecological dynamics into
fixed parameters. This compression is structural (not a design oversight)
and produces systematic blind spots. We formalize why, demonstrate via
systematic experimentation, and propose concrete evaluation redesign.

## Three-layer contribution

### Layer 1: Formalize (why compression = blindness)

**Compression Blindness Principle**: if benchmark B operates at fixed
point r* in regime space R, then for any dimension D_i that B fixes, B
has zero discriminative power along D_i (∂M_B/∂D_i = 0 by construction).

- Proves blindness is structural, not accidental
- Predicts which phenomena are invisible (any varying along fixed dims)
- Proves in-regime test addition cannot fix blindspot (must vary the regime)

**Selection-Signal Collapse Theorem** (already developed): when ecology
dimension IS varied, narration-satisfiable selection collapses to
narrative complexity under cost asymmetry (Nash equilibrium).

Two formal results: why fixing creates blindness + what happens when you
un-fix.

### Layer 2: Demonstrate via Quine (one blindspot per dimension)

For each compressed dimension: what the benchmark reports vs what Quine
reveals.

| Dimension | Benchmark reports | Quine reveals | Evidence |
|-----------|------------------|---------------|----------|
| D2 material/carrier | "copy exists" | 8+ successor forms; source ignored under stream; carrier ⊥ source | 6M wave (cite AAAI) |
| D4 ecology | "individual success" | population outcome invisible to an individual-success metric; whether signal type determines collapse was tested and resolved contradicted (P3) | dorm-06, P12 |
| D5 barrier | "copy deployed" | barrier-mediated = reconstruction not construction | 6.64/6.65, 6R |
| D6 motive | "instructed replication" | motive-free succession emerges | dorm-06 |

Each row = one blindspot demonstration.

### Layer 3: Propose evaluation redesign (constructive)

Five concrete protocols benchmark designers can use:
1. Material × demand assay (generalized from AAAI)
2. Ecology test (narration-satisfiable vs exec-dependent)
3. Property-class decomposition (multi-axis successor measurement)
4. Barrier variation (with/without hard barrier)
5. Motive variation (instructed / instrumental / no-mission)

## Structure (9 pages)

| § | Pages | Content |
|---|-------|---------|
| §1 Intro | 1 | compression claim + three-layer contribution preview |
| §2 Framework | 1.5 | six-dimension reference frame + spine + dependency structure |
| §3 Formal properties | 1.5 | compression blindness principle + selection-signal collapse theorem |
| §4 Wrong-level compression | 1 | survey: what benchmarks decompose vs compress (10 DR reports) |
| §5 Quine demonstration | 2 | per-dimension blindspot: benchmark-report vs Quine-reveal |
| §6 Evaluation redesign | 1 | five concrete protocols |
| §7 Predictions | 0.5 | P1, P3, P4, P5, P6, P7 + falsifiers |
| §8 Discussion | 0.5 | limitations; connection to ICLR-2 (safety) |

## What this paper does NOT own

- Boundary-relative safeguard-composition result → ICLR-2
- Guard-placement experiment and inheritance/composition distinction → ICLR-2
- Facultative morphology assay details → AAAI (cited)
- dorm-06 safety implications → ICLR-2

## Evidence status

| Evidence | Status |
|----------|--------|
| 10 DR reports (survey) | ✓ complete |
| 6M wave (61 records) | ✓ complete (via AAAI) |
| dorm-06 + P12 (ecology) | ✓ complete — P3 resolved contradicted (2/2 matched-pair split, seed-driven not selection-driven); P12 collapse pole unaffected |
| 6.64/6.65 + 6R (barrier) | ✓ complete |
| Compression blindness principle | needs formal write-up (simple) |
| Selection-signal collapse theorem | ✓ developed (in theory objects); reversal-corollary empirical anchoring retracted, now stated as formal/open — collapse-pole anchoring (P12) still stands |

No new experiments needed for ICLR-1. All evidence exists.

## Dependencies

- Cites AAAI (aaai-facultative-succession) for D2 demonstration
- Cited BY ICLR-2 (iclr-safeguard-composition) for framework
- Does NOT depend on ICLR-2

## Open questions — resolved 2026-06-30

- ~~Compression blindness principle: is it too trivial?~~ **Resolved.** The
  finalized principle ([`plans/theorems-formalization-2026-06-30.md`](plans/theorems-formalization-2026-06-30.md))
  lifts it above the tautology on three axes: the *closure/factorization floor*
  (every in-regime battery factors through the same window, so I(M_n;σ) is capped
  at I(S;σ) uniformly in n — refuting "more tests ⟹ more information"), the
  *falsifiable visible/blind split* (correctly predicts barrier-evasion and
  covert motive stay visible), and the *constructive minimal-repair* rule. The
  honest load sits in the per-property A2/A-Free discharge (the §5 catalog is the
  argument).
- ~~Evaluation redesign: how concrete?~~ **Resolved.** Five named protocols +
  the proposed 6-condition discriminator with pre-registration discipline (§6).
- ~~Spine: keep or drop?~~ **Resolved.** Keep, folded into §2 as a single
  paragraph; full treatment referenced to `boundary-paper`.

Remaining (genuine) open items:

- Title is the working default; the user may prefer a more aggressive
  (position-paper) lane.
- The DRAFT `references.bib` is survey-sourced; many 2026 ids need a
  primary-source check before submission.
- A2/A-Free must be discharged per property in §5; the graceful-degradation
  regime (0 < I(S;σ) < I(β;σ)) needs a concrete estimator (appendix).
