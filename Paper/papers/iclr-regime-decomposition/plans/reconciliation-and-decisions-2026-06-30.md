# Reconciliation & Decisions — iclr-regime-decomposition

> Authored 2026-06-30 during the "develop this paper" pass. This note locks the
> structural decisions that the manuscript-plan and section drafts build on. It
> resolves the inconsistencies that accumulated across the candidate file, the
> dossier README, the paradigm publication-strategy, and the prediction files.
>
> **Authority order** when documents disagree: the paradigm README's
> *Publication strategy (three-paper arc)* → Element-allocation table is
> authoritative for ownership; the survey *synthesis.md* Framing Decision Matrix
> is authoritative for evidence stance; the candidate file is the OLDEST artifact
> and is treated as stale where it conflicts.

## D-1. Dimension count: SIX (D1–D6), presented as a reference frame

The candidate's thesis sentence says "five logically independent dimensions";
its own body, experiment, and the paradigm README all use six. **Resolution:
six dimensions (D1–D6)**, matching the paradigm README, the
`regime-decomposition` theory object, and the seven predictions.

- D1 persistence → succession
- D2 bundle-binary → reconstructive-resource gradient (self-description × carrier)
- D3 binary success → property-class-decomposed
- D4 individual → ecological (selection-signal collapse)
- D5 barrier-free → barrier-mediated (exec barrier)
- D6 motivated → emergent

The paper presents D1–D6 as a **reference frame**, not a uniquely-correct
partition (Framing A, below). It explicitly separates *core regime variables*
(barrier, self-description, carrier, criterion/demand, selection, motive —
manipulable) from *structural consequences* (D1 persistence and D3 measurement
are derived, per the paradigm README's dimension audit). This pre-empts the
"why exactly N?" objection.

## D-2. Framing: A (reference frame), not B (paradigm shift)

The paradigm README leaves Framing A vs B open but states A is "best for ICLR
(theory venue)." **Resolution: Framing A.** Lead with the wrong-level-compression
anomaly; position the decomposition as a coordinate system that makes the
field's compressed dimensions visible; contribution rests on the two formal
results + the survey + the demonstration, not on a Kuhnian replacement claim.
Von Neumann's construction model is the special case where A1–A6 all hold; we
name the succession regime it does not cover. No "we replace the field" rhetoric.

## D-3. Owned vs cited (the boundary that prevents three-paper overlap)

Per the authoritative Element-allocation table:

**This paper (ICLR-1 / evaluation) OWNS:**
- the six-dimension reference frame + dependency structure + spine (spine
  *referenced*, full treatment in boundary-paper)
- the wrong-level-compression survey (10 DR reports)
- **Compression Blindness Principle** (CBP) — the headline formal result
- **Selection-Signal Collapse Theorem** (D4)
- predictions **P1, P3, P4, P5, P6, P7**
- the 6M wave as *compression demonstration* (D2/D7 morphology), dorm-06 as
  *ecology + motive* evidence (D4/D6), 6.64/6.65 + 6R as *barrier* evidence (D5)
- the constructive evaluation-redesign protocols
- the proposed 6-condition full-inversion experiment as a *discriminating future
  design* (not run; Framing C below)

**This paper CITES a companion for (does NOT own):**
- the **conditional non-identifiability / non-inheritance theorem** → cite
  `iclr-safeguard-composition` (its formal core; prediction **P2**)
- the **safety non-inheritance corollary** + guardrail-loss / seeded-marker
  experiment → `iclr-safeguard-composition`
- dorm-06 *as safety evidence* → `iclr-safeguard-composition` (we use dorm-06 only
  for ecology/motive, never for guardrail loss)
- the **facultative-morphology assay** full details → `aaai-facultative-succession`
  (we cite the 6M wave as evidence, we do not re-run the assay)

**MUST NOT SAY** (reserved, would create overlap):
- do not state or prove the conditional non-inheritance theorem (only cite it)
- do not run the safety/guardrail argument as a result of *this* paper
- do not develop the facultative-morphology indexing (D,F) as a contribution
- do not invoke the `regime-blindness` umbrella paradigm (not yet publishable)

## D-4. Theorem ownership conflict — candidate is stale

The candidate's *Novelty boundary* claims the "conditional non-inheritance
theorem" as a novel contribution of THIS paper. The authoritative table assigns
it to `iclr-safeguard-composition` (P2). **Resolution: this paper's two owned
theorems are the Compression Blindness Principle + the Selection-Signal Collapse
Theorem only.** The conditional non-inheritance theorem is cited as the
companion's contribution. → *Flagged staleness item S-1; candidate to be
corrected (it is a `kind: candidate` file, editable), but the correction is a
separate, surfaced change.*

## D-5. Evidence stance: Framing C (hybrid) — demonstrate-via-existing + propose-the-discriminator

The dossier README says "no new experiments needed; all evidence exists." The
candidate carries an ambitious 6-condition experiment. The synthesis Framing
Decision Matrix recommends **Framing C (hybrid)**: wrong-level compression
(survey + formal) + pilot data for specific dimensions + a *proposed*
full-inversion experiment.

**Resolution:**
- The paper's core claims are carried by the survey + the two theorems + the
  per-dimension demonstration using existing waves. No new experiment is
  required for the paper to stand.
- The 6-condition full-inversion study is presented as the **proposed
  discriminating experiment** (the genuinely-novel full-inversion existence
  proof, per synthesis §3), with the existing ladder mapped as pilot data.
- Honesty: present the experiment as design + pilot, not results-mode. dorm-06
  caveats (seed exemplar, model-bound, no ratchet) are stated.

## D-6. Spine: keep, but subordinate to §2

The dossier open question asks whether to keep the spine. **Resolution: keep it,
folded into §2 (framework)** as the one-paragraph statement of what the
decompressed theory composes into (barrier → self-description × carrier →
criterion → aperture → ecology → emergence), with the full treatment
*referenced* to boundary-paper (the .core projection). It is not a standalone
theory-heavy section — this is an evaluation paper.

## D-7. Prediction → dimension mapping (corrected)

| Pred | Dimension | Owner |
|------|-----------|-------|
| P1 | D5 (carrier = transmission unit) | this paper |
| P2 | D2 (criterion-mediated heritability → non-inheritance) | **competent-stranger** |
| P3 | D4 (selection-signal collapse) | this paper |
| P4 | D3 (property-class loss) | this paper |
| P5 | D1 (addressability continuity) | this paper |
| P6 | D6 (emergent succession) | this paper |
| P7 | D2 (self-description-gradient morphology) | this paper |

**Staleness item S-2:** the file `predictions/open/regime-decomposition-P7.md`
labels itself "D7" in several places and says "D7 collapses into D2" — P7 is a
**D2 sub-prediction** (self-description gradient), not a seventh dimension. This
file is `scope: paradigm` → **human-directed; do NOT auto-edit.** Surface as a
proposal. The paper itself presents P7 correctly as a D2 consequence.

## D-8. Title

Dossier title is "TBD". Candidate title: *"Transmission, Not Construction:
Regime-Decomposed Evaluation of AI Self-Replication."* Dossier framing:
*"Evaluation Blind Spots in AI Self-Replication."*

**Recommended default (used in metadata until the user overrides):**

> **Decomposing at the Wrong Level: Structural Blind Spots in AI
> Self-Replication Evaluation**

Alternates: (a) *Transmission, Not Construction: A Regime Decomposition of AI
Self-Replication Evaluation*; (b) *What Self-Replication Benchmarks Cannot See:
Heritability, Property-Class Survival, and Ecological Collapse*.

## Section map (9 pages, reconciled)

| § | pp | content | owns/cites |
|---|----|---------|-----------|
| 1 Introduction | 1 | wrong-level-compression anomaly; reference-frame stance; three contributions (CBP, SSC, redesign) | owns |
| 2 Framework | 1.5 | A1–A6; D1–D6 reference frame; core-vars vs structural-consequences; spine (1 para, ref boundary-paper); dependency structure | owns |
| 3 Formal properties | 1.5 | Compression Blindness Principle (+ proof); Selection-Signal Collapse Theorem (+ proof); cite conditional-non-inheritance to competent-stranger | owns CBP+SSC; cites P2 |
| 4 Wrong-level compression | 1 | survey: what RepliBench/SOCK/Agent Matrix/METR vary vs compress; the named gaps (DR-1 §7.1–7.4) | owns |
| 5 Demonstration | 2 | per-dimension benchmark-reports-vs-decomposed-reveals using existing waves (6M, dorm-06, 6.64/6.65+6R) | owns; cites AAAI for 6M |
| 6 Evaluation redesign | 1 | five constructive protocols (material×demand assay, ecology test, property-class decomposition, barrier variation, motive variation) + the proposed 6-condition discriminator | owns |
| 7 Predictions | 0.5 | P1,P3,P4,P5,P6,P7 + falsifiers; P2 noted as companion | owns; cites P2 |
| 8 Discussion | 0.5 | limitations (decomposition non-uniqueness, substrate specificity); links to AAAI + competent-stranger | owns |

## Open items surfaced for the user (not blockers)

- **S-1** candidate `candidates/iclr-regime-decomposition.md` claims the
  conditional non-inheritance theorem (belongs to competent-stranger) and says
  "five dimensions". Editable; correction proposed.
- **S-2** prediction `regime-decomposition-P7.md` mislabeled "D7" (`scope:
  paradigm`, human-directed — proposal only).
- **Title** default chosen above; user may override.
- **Framing A** chosen per repo guidance; user may prefer the more aggressive
  Framing B for a position-paper lane.
