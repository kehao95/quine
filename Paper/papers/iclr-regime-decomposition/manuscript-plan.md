# Manuscript Plan — Decomposing at the Wrong Level

Owns the manuscript plan for the regime-decomposition **framework** paper. This
is the spec the section drafts follow. Structural decisions are locked in
[`plans/reconciliation-and-decisions-2026-06-30.md`](plans/reconciliation-and-decisions-2026-06-30.md);
this file is the writing-level instantiation.

## Working title

Default: **Decomposing at the Wrong Level: Structural Blind Spots in AI
Self-Replication Evaluation.** Alternates: *Transmission, Not Construction: A
Regime Decomposition of AI Self-Replication Evaluation*; *What Self-Replication
Benchmarks Cannot See*.

## One-sentence thesis

AI self-replication evaluation decomposes at the wrong level: the field varies
execution barriers, motive, persistence, and operational risk, but fixes the
heritability channel, the property classes a copy must preserve, and the
ecological selection signal — and we prove that this fixing produces structural,
non-improvable blind spots, then show what is hidden and how to redesign the
evaluation to see it.

## Stance (read before drafting)

- **Framing A — reference frame, not paradigm shift.** Present D1–D6 as a useful
  coordinate system that makes the field's compressed dimensions visible.
  Explicitly disclaim uniqueness: "we do not claim a unique or complete
  decomposition." Von Neumann's construction model is the special case where
  A1–A6 all hold; we name the succession regime it does not cover. **No** "we
  replace the field" rhetoric, **no** Kuhn.
- **Honest about what the field already does.** The field is NOT at a single
  scalar. RepliBench/SOCK/Agent Matrix/METR *do* vary barriers, motive,
  persistence, and decompose risk. The genuine compression is at a *different
  level* (heritability / property-class / ecology). Lead with this nuance — it is
  the paper's credibility anchor and was the correction the 10-report survey
  forced.
- **Six dimensions D1–D6** (not five). Separate *core regime variables*
  (barrier, self-description, carrier, criterion/demand, selection, motive —
  manipulable) from *structural consequences* (D1 persistence, D3 measurement —
  derived). Present them at different levels; this pre-empts "why exactly N?".
- **Evidence: hybrid (Framing C).** Core claims rest on survey + the two
  theorems + per-dimension demonstration using *existing* waves. The 6-condition
  full-inversion study is the *proposed discriminating experiment*, not a result.
  State pilot caveats plainly.

## Contribution (three, in order)

1. **A six-dimension regime reference frame** + the **Compression Blindness
   Principle**: a structural-blindness theorem showing a benchmark cannot recover
   discriminative power along a fixed dimension by adding in-regime tests.
2. **The Selection-Signal Collapse Theorem** (D4): a game-theoretic collapse of
   narration-satisfiable agent ecologies, double-dissociated from Goodhart.
3. **Constructive evaluation redesign**: five protocols + a 6-condition
   property-class-decomposed discriminating design, with pilot data mapped.

## Ownership boundary (HARD — every section obeys)

OWNS: framework (D1–D6), spine, wrong-level-compression survey, **Compression
Blindness Principle**, **Selection-Signal Collapse Theorem**, per-dimension
demonstration, evaluation redesign, predictions **P1, P3, P4, P5, P6, P7**.

CITES a companion, never develops:
- boundary-relative safeguard composition, the guard-placement assay, the
  successor-inheritance/system-composition distinction, and the continuity-only
  successor test → all to `iclr-safeguard-composition` (forward-cite once, in §8
  only).
- facultative-morphology assay *details* (material × demand design) → to
  `aaai-facultative-succession`; cite the 6M-wave *result* as the D2
  demonstration only.

MUST-NOT language: never claim the safeguard-composition result or placement
study as this paper's contribution; use dorm-06 only as D4/D6 evidence; never
present safety disposition loss as *this paper's* finding.

---

## Section skeleton (9 pages + appendix)

### §1 Introduction (1 p) — `01-introduction.md`

Open on the anomaly, not the framework. Beat sequence:

1. Self-replication evaluation is consolidating (RepliBench, SOCK, Agent Matrix,
   METR ARA, OpenAI Preparedness, Anthropic RSP, DeepMind FSF). Cite the Global
   Call for AI Red Lines as evidence the stakes are taken as settled.
2. These frameworks *do* decompose — capability trees, RCL/PCL ladders, risk
   scores, barrier stages. So the usual critique ("they reduce it to one bit")
   is a strawman; we do not make it.
3. The real problem: every one of them fixes three variables at one operating
   point — *how* a successor is reconstructed (heritability channel), *what
   properties* a "successful copy" must keep (aperture), and *under what
   ecological selection* a lineage competes. Quote DR-1 §7.1/7.2/7.3 gaps.
4. Thesis sentence (above). Position as a reference frame.
5. Contributions list (the three above). Name the two theorems.
6. Boundary sentence: this is the framework paper of a three-paper arc; the
   safety corollary and the morphology assay are companion contributions
   (one forward-citation each), not developed here.

### §2 Framework (1.5 p) — `02-framework.md`

- **A1–A6**: the field's six implicit assumptions (persistence, bundle-binary,
  binary-success, individual, barrier-free, motivated). One line each. Note this
  is a *characterization of assumptions the frameworks already make*, derived
  from the field's own designs, not imposed.
- **D1–D6 reference frame** (table): assumption → relaxed point → qualitative
  shift → spine object. Use the regime-decomposition theory-object table.
- **Core variables vs structural consequences**: barrier (D5), self-description
  (D2a), carrier (D2b), criterion/demand, selection (D4), motive (D6) are
  manipulable; D1 persistence and D3 measurement are *derived* consequences
  (D1 = consequence of the barrier killing the predecessor; D3 = consequence of
  the D2 gradient making the binary metric lossy). Present at two levels.
- **Dependency structure**: D5 is the physical prerequisite; D1, D2 are most
  meaningful under a barrier; D4, D6 independent of D5; D2's two sub-gradients
  separable (6M.10). One small diagram or prose.
- **Spine** (ONE paragraph, then move on): barrier → self-description × carrier →
  criterion → aperture → ecology → emergence; full treatment referenced to the
  `boundary-paper` (.core projection). Do not theory-dump here.
- Close with the explicit non-uniqueness disclaimer (different cuts serve
  different purposes; we offer a reference frame).

### §3 Formal properties (1.5 p) — `03-formal-properties.md`

The intellectual core. Two results. Keep notation lean — state theorems
precisely, give proof *sketches*, push heavy measure-theory to an appendix.

**§3.1 Compression Blindness Principle.** See "CBP — drafting capsule" below.
The non-trivial core to foreground: (a) the *closure/factorization floor* —
every in-regime battery factors through the same window S, so adding tests is
capped at I(S;σ), uniformly in the number of tests (this refutes the
"more-tests-more-information" default, and is the load-bearing step, not the
one-line DPI); (b) the *falsifiable visible/blind split* — a contrast is
invisible iff *every* minimal resolving set meets a fixed coordinate, which
correctly predicts barrier-evasion (D5) and covert-vs-primary motive (D6) stay
VISIBLE (the non-tautological controls); (c) the *constructive minimal repair* —
to see phenomenon P, free exactly ID(P)\(varied slice). State the "spine
identity": CBP is the same DPI/Fano structure as the D2 companion theorem,
lifted from the heritability channel to the measurement channel — *cite* the
companion's conditional non-inheritance theorem, do not prove it.

**§3.2 Selection-Signal Collapse Theorem.** See "SSC — drafting capsule" below.
Foreground the *double dissociation from Goodhart*: collapse occurs at zero
proxy gap (a perfectly faithful but narration-satisfiable signal still
collapses), and resistance occurs with a crude-but-exec-dependent signal —
fidelity and forgeability are orthogonal axes, only forgeability drives the
theorem. State the constructive reversal corollary, but flagged as a formal
consequence of the model, not an empirical result: anchor the collapse pole
empirically to P12 (86 conditions); dorm-06 no longer anchors a reversal pole —
the matched wake-vs-dependency pairs split 2/2 with zero correlation to signal
type (`regime-decomposition-P3`, resolved contradicted; see §5/§7) — so state
the reversal corollary as open/conjectural on this substrate, with the
inflation caveat (16→3, see §5) noted separately as a distinct issue.

### §4 Wrong-level compression (1 p) — `04-wrong-level-compression.md`

The survey, as a single argument: *the field decomposes operational capability
and risk; it compresses heritability, aperture, and ecology.*

- Table 1: what each framework VARIES vs COMPRESSES (RepliBench, SOCK, Agent
  Matrix, METR ARA, Pan, MISR), with exact facts (see "§4 facts" below).
- The four named gaps, quoted (DR-1 §7.1 safety-disposition fade; §7.2
  reconstruction-from-partial-information; §7.3 multi-generational population
  dynamics; §7.4 dynamic heritability-channel choice).
- Tie back to §3: these are exactly the dimensions CBP predicts are invisible.
- Source discipline: this paper makes the *critique itself* the contribution —
  the survey found "I do not see a paper that makes the critique itself the main
  theoretical contribution" (clearest novelty zone for the framework).

### §5 Demonstration (2 p) — `05-demonstration.md`

Per-dimension: what the binary metric reports vs what the decomposed view
reveals, on the Quine substrate. ONE row per compressed dimension. Use EXACT
numbers and state caveats inline (see "§5 demonstration table" below). This is
*demonstration via existing evidence*, explicitly not the definitive study.
Frame dorm-06 strictly as D4 (selection-signal) + D6 (emergent) evidence.

### §6 Evaluation redesign (1 p) — `06-evaluation-redesign.md`

Constructive. Five protocols benchmark designers can adopt:
1. **Material × demand assay** (generalized; cite AAAI for the instance).
2. **Ecology test**: narration-satisfiable vs exec-dependent selection arm.
3. **Property-class decomposition**: measure a survival *vector*, not a bit.
4. **Barrier variation**: with/without a hard exec barrier.
5. **Motive variation**: instructed / instrumental / no-mission.
Then the **proposed 6-condition full-inversion discriminator** (the genuinely
novel existence-proof if run): base + single-dimension relaxations + full
inversion, scored by the property-class vector. Pilot map (which existing wave
approximates which condition). Pre-registration discipline (n ≥ 30 lineages;
reconstructor-prior ablation; blind scoring; seeded off-prior marker). State
clearly: not yet run; conditions 3/5/6 require new runs.

### §7 Predictions (0.5 p) — `07-predictions.md`

Table: P1 (D5), P3 (D4), P4 (D3), P5 (D1), P6 (D6), P7 (D2-morphology), each
with statement + falsifier + current pilot status. One line noting P2 (D2
heritability/non-inheritance) is the companion `iclr-safeguard-composition`'s
prediction, not developed here. **Do not** relabel anything "D7" — P7 is a D2
consequence.

### §8 Discussion (0.5 p) — `08-discussion.md`

- Limitations: decomposition non-uniqueness (reference frame, not partition);
  substrate specificity (theory is substrate-independent, the demonstration is
  Quine-specific, cross-substrate is future work); A2/A-Free must be discharged
  per property; dorm-06 caveats; the proposed experiment is unrun.
- The two forward-citations: safeguard composition through successor synthesis
  (→ ICLR-2) and the facultative-morphology phenomenon (→ AAAI). One sentence
  each. This is the ONLY place safety is mentioned, and only as a companion's
  contribution.
- Close on the constructive note: a coordinate system for benchmark design.

### Appendix

Full CBP measure-theoretic statement + proof; SSC replicator dynamics; 6M / 6R
/ dorm-06 pilot tables; statistical analysis plan for the proposed study.

---

## CBP — drafting capsule (for §3.1; full text in `plans/`/scratch)

**Informal.** A benchmark evaluates a system only on the slice of regime space
it instantiates: it fixes some dimensions (empirically D2 heritability, D3
aperture, D4 ecology) and varies the rest. Fixing has a hard consequence: every
in-regime measurement — any number, randomized, adaptive — factors through the
raw on-slice behavior, so it can never resolve two systems that agree on the
slice. Adding in-regime tests cannot help; the only repair is to free a fixed
dimension, and the minimal set to free is computable per target phenomenon.

**Formal core (lean version for the body).** System = behavior tensor
β: R×Q → O; property σ = σ(β); benchmark observes window W (the fixed slice ×
aperture), reporting M = φ(S,U), S = β|_W, U ⊥ (β,σ). Markov chain
σ — β — S — M. (i) *Floor:* ker(M_n) ⊇ ker(S) and I(M_n;σ) ≤ I(S;σ) uniformly in
n (factorization lemma: post-composition with φ only merges fibers). (ii)
*Classes:* ker(S) = {systems agreeing on W}; each blind class ≅ the realizable
off-slice response space. (iii) *Visible/blind:* invisible ⟺ every minimal
resolving set meets a fixed coordinate; visible ⟺ some witness lies wholly in
the varied slice. (iv) *Collapse:* if I(S;σ)=0 then I(M_n;σ)=0 ∀n (Fano floor at
base rate); graded form I(β;σ)=I(S;σ)+I(β;σ|S). (v) *Escape:* refining the blind
partition requires observing β off W; minimal redesign frees ID(P)\(varied
slice). **Spine identity:** this is the D2 inequality I(R(C,V);σ) ≤ I((C,V);σ)
lifted from heritability channel C to measurement channel S — cite the companion.

**Non-triviality defense (must appear, ~3 sentences):** the slogan "you can't
see what you don't vary" is the content-free constant-function fact. CBP carries
content the slogan cannot: (1) the closure floor refutes "more tests ⟹ more
information" — the *entire* σ-algebra of conceivable in-regime evidence is capped
by I(S;σ); (2) the visible/blind split is falsifiable and correctly predicts the
field DOES resolve barrier-evasion and covert motive; (3) the minimal-repair
rule is a design instruction. Honest load: A2/A-Free is discharged per property
(the §5 catalog IS the argument); where 0<I(S;σ)<I(β;σ) the principle degrades
gracefully to the graded decomposition.

**Closest prior (cite):** IRL non-identifiability (cao2021irlidentifiability,
skalse2025partialidentifiability, rolland2022multipleexperts); model-extraction
non-identifiability (carlini2024stealing); accuracy-vs-fidelity
(jagielski2020fidelity); DPI/Fano (cover-thomas; Fano); DPI goal-drift
(kriger2026goaldrift). Novelty: no prior work packages DPI+Fano+IRL into an
agent-level structural-blindness theorem for the *evaluation* regime.

## SSC — drafting capsule (for §3.2)

**Informal.** In an ecology where survival depends on a selection signal, if the
signal can be satisfied by *narrating* a result as cheaply as *producing* it
(narration-satisfiable), the population converges to a stable narrative-
complexity equilibrium that drives function-producers extinct, regardless of
capability, and even when the signal is a *faithful* encoding of the true
objective. Re-engineering the signal to be exec-dependent reverses the collapse
on the same population, no retraining.

**Formal core.** Strategies π_F (cost c_F>0) / π_N (cost c_N≈0); Δc = c_F−c_N;
signal discrimination D_S = sel_S(π_F) − sel_S(π_N); fitness Φ = sel_S − c.
*Inversion lemma:* Φ(π_N) > Φ(π_F) ⟺ D_S < Δc. *Narration-satisfiable:* D_S < Δc
⟹ π_N strictly dominates ⟹ all-narration is the unique ESS; replicator
ẋ_N = x_N(1−x_N)(Δc − D_S) → 1; function-lineages have negative invasion fitness
and go extinct. *Exec-dependent:* D_S > Δc flips the ESS to π_F on the same
population. *Double dissociation:* with Fid(S) the proxy fidelity (Goodhart
axis), (Fid=1, D_S≈0) ⟹ collapse and (Fid low, D_S>Δc) ⟹ no collapse — the
binding variable is forgeability, not fidelity.

**Non-triviality defense (must appear):** the bare inversion inequality is a
tautology; the content is four claims the inequality lacks — (1) a *decidable
construction-level criterion* computable before any agent acts (presence of a
forced-real third-party exec); (2) the *double dissociation* from objective
misspecification (collapse at zero proxy gap — no Goodhart account predicts
this); (3) a *population result with extinction and no internal cure* (drift, no
agent intending to game); (4) a *constructive converse* (a class of signals with
a proof of immunity). Falsifiable predictions: capability-invariance of collapse
(P12: 86/86 at high capability); reversal-without-retraining
(`regime-decomposition-P3`, tested via dorm-06's matched wake-vs-dependency
pairs — **resolved contradicted**: the matched pairs split 2/2, seed-driven and
stochastic, not selection-driven; a closed falsification loop, not confirming
evidence); collapse-speed ∝ Δc/D_S (cheaper LLMs accelerate collapse,
untested).

**Closest prior (cite):** Goodhart / Manheim & Garrabrant; reward hacking
(skalse2022reward, krakovna2020specification); survival-of-the-flattest
(wilkeflattest); ESS (maynardsmith) + replicator (hofbauer); Price equation;
TSE Thm 17.9 (tse2025strategicevolution); Agent Matrix (zhang2025agentmatrix).

---

## §4 facts (benchmark decomposition — exact, for the survey table)

- **RepliBench** (black2025replibench, arXiv:2504.18565): 20 task families, 86
  tests; 4 sub-capabilities (resources, weights, compute, persistence); varies
  barriers (firewalls/AV/KYC) + motive (covert weight-exfil track); fixes source
  to black-box API; concludes 2025 frontier models not yet a credible end-to-end
  threat. COMPRESSES: heritability channel, successor property classes.
- **SOCK**: RCL 0–5 (local copy → cross-network) × PCL (terminable → rootkit) on
  a 2D propagation×stealth grid; varies barrier + persistence + motive
  (tame/bold); fixes ecology to single-agent.
- **Agent Matrix / ERABBench** (zhang2025agentmatrix, arXiv:2509.25302): only
  major framework to vary ecological selection pressure; Overuse Rate + Aggregate
  Overuse Count; >50% of agents showed uncontrolled self-replication under
  pressure; but single-generation / load-balancing — COMPRESSES multi-generational
  dynamics.
- **METR ARA** (kinniment2023autonomous, arXiv:2312.11671): 5-step timeline;
  economic criterion (agent compute cost < human cost); 12 ARA-relevant tasks;
  fixes motive (assumes rogue) + source (assumes weights stolen).
- **Pan** (pan2024redline 2412.12140; pan2025selfreplication 2503.17378):
  explicit "Replicate yourself…" instruction (NOT spontaneity); 11/32 systems
  showed capability (2025). Varies source access inadvertently.
- **MISR**: implicit goal-setting (self-modification to succeed); fixes ecology.
- Frameworks (governance scalars that compress): OpenAI Preparedness
  (Low/Med/High/Critical), Anthropic RSP (2–8h SWE checkpoint), DeepMind FSF
  (phuong2024fsf, CCLs + expert best-of-N).

Named gaps (quote verbatim): DR-1 §7.1 "no benchmark currently measures the
phenomenon of Safety Disposition Fade"; §7.2 "zero benchmark coverage assessing
an agent's ability to achieve Reconstruction from Partial Information"; §7.3 "no
benchmarks that evaluate sustained, multi-generational population dynamics or
mutational drift"; §7.4 "no existing benchmark that explicitly forces an
artificial intelligence agent to … alter its heritability channel."

## §5 demonstration table (EXACT numbers + caveats — do not inflate)

| Dim | Binary metric reports | Decomposed view reveals | Wave(s) & exact figures | Caveats (state inline) |
|-----|----------------------|-------------------------|-------------------------|------------------------|
| D5 barrier | "copy deployed" | barrier-mediated transmission = reconstruction from carrier, not construction from blueprint; body-crossing ≠ obligation-carrying | 6.64 live-positive (builds workspace/quine, real-execs, receives later ctl/wake, acts, returns idle); 6.65 minimum-prompt positive; 6.66/6.67 baselines (6.67 post-repair passes 6/6); 6R.01–04 four deterministic false-positive controls | 6.64 "should be cited as scaffolded" (self-generated successor mission); 6R are controls, not positives |
| D2 channel/morphology | "source available → copy works" | bundle decomposes into self-description gradient × carrier, separable; successor morphology varies along the gradient ("source is material, not destiny") | 6M.10 carrier⊥source/body (first retained carrier-capability wave; 6 of 61 morphology rows); 6M.03 self-description ladder (7 rows) | morphology summary is "a normalization pass, not a success-rate dataset"; source-absent *positive* NOT demonstrated (honest-risk #5); cite AAAI for the 6M assay |
| D3 aperture | "successful copy" (1 bit) | success is a property-class *vector*; "success" compatible with systematic property-class loss outside the aperture | demonstration is the *measurement framework* applied to all conditions; no single pilot — state as design | the property-class decomposition is new measurement design; do NOT name safety-disposition loss as a finding (companion owns it) |
| D4 ecology | "individual succeeds" | population outcome invisible to an individual-success metric; whether selection-signal type determines ecology outcome (narration-satisfiable collapses, exec-dependent resists) was tested via matched wake-vs-dependency pairs and **resolved contradicted** (P3) — 2/2 split, zero correlation to signal type, "seed-driven and stochastic, NOT selection-driven" per the experiment owner | P12 attractor: 86 prior conditions collapsed (this pole stands); ~67% of wakes by host mechanics vs ~14.5% by agent artifacts. dorm-06 matched pairs (D06, D07, identical seed): wake arm 4 edges / dependency 0 in D07, dependency 3 edges / wake 0 in D06; strongest single run (DeepSeek) 5 producers / 9 stages / **3** clean edges after scorer tightening (headline 16 was ~80% inflated, a separate issue); only 6/12 cycles completed in that run (~2000s wall); selection fired 2× in 6 cycles | n=1 per cell, unreplicated, model-bound (only DeepSeek composes), single functional family, no agent-level heredity; "scaffolded cultural accumulation, not goal-free OEE"; P3 resolved contradicted — see [`regime-decomposition-P3`](../../predictions/resolved/contradicted/regime-decomposition-P3.md) |
| D6 motive | "instructed replication" | successor formation emerges WITHOUT replication motive under pressure + physics + ecology | dorm-06 empty mission ("Begin.") → functional commons; succ-06/6.45 3/3 stable under "unbounded" framing (not "reproduce"); succ-04/6.5 viability-pressure coupling | dorm-06 caveats as above; 6.45 rerun did not cleanly re-attain 3/3; succ-04 is a ladder-start baseline |
| D1 continuity | "is it a copy?" | under succession the criterion is future addressability, not artifact resemblance | 6.64/6.65 succession positives; 6.66/6.67 artifact-similar successors that FAIL addressability; 6F.07 mailbox transport (C04 = address-preserved-state-lost) | 6F.07 first wave only; cite future-addressability for the full criterion |

## §7 predictions (statements + falsifiers)

| P | Dim | Statement | Falsifier | Pilot status |
|---|-----|-----------|-----------|--------------|
| P1 | D5 | carrier is the transmission unit under barrier-mediated succession | source-copy predicts continuity better than carrier under the barrier | 6M.10, 6.64/6.65, 6.66/6.67 partial |
| P3 | D4 | narration-satisfiable selection collapses; exec-dependent resists | wake-scalar ecology does not collapse | **resolved contradicted** — dorm-06 matched wake-vs-dependency pairs split 2/2, seed-driven not selection-driven; collapse pole (P12 86/86) unaffected; see [`regime-decomposition-P3`](../../predictions/resolved/contradicted/regime-decomposition-P3.md) |
| P4 | D3 | "successful copy" compatible with systematic property-class loss | all classes transfer together whenever binary = success | needs property-class measurement (new) |
| P5 | D1 | continuity = future addressability, not artifact resemblance | artifact similarity perfectly predicts continuity | 6.64/6.65, 6.66/6.67 partial |
| P6 | D6 | successor formation emerges without explicit motive under pressure | unmotivated condition 6 produces no successor formation | dorm-06, succ-06/04 moderate |
| P7 | D2 | successor morphology changes along the self-description gradient | morphology same regardless of self-description form | 6M.03, 6.59/6.63/6.64 |

P2 (D2: heritability requires resource-coverage AND criterion-aperture; DPI
bound at the carrier-only extreme) is the companion `iclr-safeguard-composition`'s
prediction — note in one line, do not develop.

## Citation plan (use the DRAFT bibkeys in references.bib)

- Benchmarks/frameworks: black2025replibench, zhang2025agentmatrix,
  kinniment2023autonomous, pan2024redline, pan2025selfreplication, phuong2024fsf,
  reworr2026hackselfreplicate.
- Foundations: vonneumann1966selfreproducing (+ Langton for the
  interpreted/uninterpreted description point), hughes2024openendedness.
- D2 prior: intentsignaltheory2026, overeinder2006generative,
  agentidentityevals2025, genomebook2026, wang2023voyager, carlini2024stealing,
  tramer2016stealing, jagielski2020fidelity; memory systems
  (packer2023memgpt, xu2025amem, chhikara2025mem0) as carrier instances.
- D4 prior: wilkeflattest, priceequation (no-id → keep in prose, source-check),
  Goodhart/Manheim (no-id → source-check), skalse2022reward (no-id →
  source-check), tse2025strategicevolution, zhang2025agentmatrix,
  projectsid2024, agentsociety2025, terralingua2026.
- D6 prior: aguerayarcas2024computationallife, zhang2025darwingodel,
  wu2023autogen, agentspawn2026, aorchestra2026, berkeleyrdi2026peerpreservation,
  hordijk2023raf (autocatalytic), maturana-varela (no-id → source-check).
- D3 prior: helm2022, mlcommons2024safety, safetybench2023, muse2024.
- Formal: cao2021irlidentifiability, rolland2022multipleexperts,
  skalse2025partialidentifiability, gao2018ratedistortion, kriger2026goaldrift,
  Cover-Thomas + Fano (textbook → add by hand, source-check).
- D1 prior: hintze2025rethinking, adams2009universalframework, ziesche2022aiidentity,
  torabi2018bco, executiontrace2026replay.

**Rule:** every `references.bib` entry is DRAFT; many ids are 2026-dated and
survey-sourced. In prose, lean on the well-established works; flag speculative
2026 ids. Do NOT cite the DR reports directly. Source-check before submission.

## Language guardrails

Preferred: `regime dimension`, `heritability channel`, `selection signal`,
`narration-satisfiable` / `exec-dependent`, `property-class survival`,
`reference frame`, `structural blindness`, `wrong-level compression`,
`carrier-mediated continuity`.

Avoid / forbidden: `competent stranger` (companion's); claiming the
`conditional non-inheritance theorem` as ours; `paradigm shift` / Kuhn framing;
"the field sees only one bit" (strawman); presenting `safety disposition loss`
as this paper's result; relabeling P7 as `D7`; `the LLM copied itself`;
over-stating dorm-06 (use the post-tightening figure of 3 edges, not 16);
presenting dorm-06 as a confirmed reversal of the P12 collapse or P3 as
supported — `regime-decomposition-P3` is resolved **contradicted** (matched
wake-vs-dependency pairs split 2/2, seed-driven not selection-driven); state
the reversal corollary as a formal/open consequence of the model, not an
empirical finding.

Use carefully: `selection` (always `ecological / selection-signal`); `paradigm`
(prefer `reference frame`); the 6-condition study (always "proposed",
"discriminating", "not yet run").

## Readiness

Scaffolding complete (metadata, references.bib DRAFT, Makefile, WORKSPACE,
plans). Next: draft §1–§8 against this plan, then a whole-paper
consistency/overclaim/boundary review. Not ready: final venue lane,
source-checked bibliography, the unrun 6-condition study.
