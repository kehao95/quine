---
date: 2026-07-04
mode: gap-survey (field-first)
supersedes_framing: complements — does NOT replace — the three novelty-mode sweeps
  (prior-art-dr-report.md, matrix-novelty-dr-report.md, novelty-adjacency-dr-report-2026-06-18.md).
  Those asked "is our claim novel?"; this asks "what does each neighboring field name as unsolved,
  and which can our existing evidence close?"
method: 5 parallel Claude subagents (general-purpose, WebSearch/WebFetch), one homogeneous
  neighborhood each, per the updated deep-research-query contract § Query Modes. Briefs +
  evidence inventory: ./plans/gap-survey-briefs-2026-07-04.md.
---

# Full-Paper Gap Survey — Structural Elicitation (field-first)

The first **gap-mode** research pass for this dossier. All prior sweeps were
novelty-mode (claim-first) and returned "narrow residual, heavily pre-empted."
This pass asks the inverse question and reaches the opposite-signed conclusion:
**four to five distinct full papers are viable, each answering a different
neighboring field's explicitly-named open problem, all powered by one shared
apparatus.**

## The meta-finding (holds across all five neighborhoods)

Every adjacent field **names** the phenomena this program studies — but only as
something to **suppress** (AI-safety side effects), **describe** (ALIFE emergent
culture / stigmergy), or **engineer as an external wrapper / practitioner
guidance** (agent harness, permission guards). **None treats them as a
controlled, causally-isolated, reward-free measurement.** Our apparatus (E1
matched-pair contrast + frozen weights + missionless prime + un-fakeable
filesystem/exec DV + cross-model) is precisely the isolation instrument these
fields' own limitations sections request.

The novelty checks measured the wrong thing. Asking "is our theory novel?"
returns "owned/narrowed." Asking "what does the field want measured?" returns
four fields each requesting our exact instrument. **The contribution is not the
theory (conceded, owned) — it is the controlled, reward-free measurement of
phenomena everyone else only theorizes, suppresses, or engineers.**

## Evidence inventory (the "our evidence" key)

See [`./plans/gap-survey-briefs-2026-07-04.md`](plans/gap-survey-briefs-2026-07-04.md)
§ Shared evidence inventory for E1–E8 definitions. In brief: E1 five-motif spine
(96/102 vs 0/85, 3 families); E2 two-latch gate; E3 information-preservation
invariant (`cp` over `mv`, 0/20 unprompted, warrant-released); E4 capture/scope
(9/10 free → 0/10 tasked); E5 non-code prose generalization; E6 substrate-locus
(EROFS attempts); E7 model axis; E8 geometry knobs + non-monotone dose-response.

## Five candidate directions (ranked by first-paper leverage)

Leverage = (field names the gap verbatim, not surveyor-inferred) × (our evidence
closes it directly) × (closing experiment is "package/extend" not "new
apparatus") × (venue clear and hungry) × (opens a NEW audience vs competes with
the ALIFE LBA).

### Direction A — The emergent conservative/destructive action prior (SAFETY) ★ strongest first paper
- **Field-named gaps (verbatim, hot, multiple):**
  - "compulsion to act" as an emergent not-trained-in bias, cross-model
    difference *unexplained* — *Selectively Quitting* (arXiv 2510.16492, [V]
    verbatim): "agents have a strong inherent bias toward attempting to complete
    a task… this effect is largely seen in the GPT family of models."
  - overeager destructive out-of-scope actions measured, **mechanism explicitly
    scoped out** — *Overeager Coding Agents* (arXiv 2605.18583, ID verified):
    stripping consent declarations flips Claude Code 0.0%→17.1%, no mechanistic
    account; limitations concede it cannot cover "non-shell sinks and
    non-enumerable authorization boundaries."
  - benign-setting propensity is a **systematic blind spot** — *Taxonomy of
    Safety Benchmarks* (arXiv 2605.16282, [V]): "only 6 [benchmarks] directly
    probe agent-internal tendencies… adversarial robustness and autonomous
    alignment are distinct risk dimensions."
  - task converts conservative agents into aggressive ones — *What Breaks When
    LLMs Code* (arXiv 2605.30777): agents "translate benign user intents into
    aggressive environmental modifications"; framework-vs-model variance measured
    but unexplained (SNARE arXiv 2605.28122: 56% framework vs 21% model).
- **Our evidence:** E2 (two-latch gate) supplies the decision structure the field
  declines to theorize; E3 (`cp`-over-`mv`, warrant-released 4/5, p<0.001) is a
  **direct existence proof of reward-free, baseline-free side-effect avoidance** —
  a counterexample to the decade-old impact-measure lineage (Amodei 2016 → AUP →
  SafeLife → Lindner 2021) assumption that avoidance must be *engineered*; E5
  shows it is weight-level not a coding-scaffold artifact (4 families, non-code
  prose); E4 supplies the missionless baseline that isolates the prior from task
  pressure.
- **Closing experiments (low apparatus cost — port probes into existing scenario
  suites):** run the two-latch / `mv`-refusal probes inside OverEager-Gen / SNARE
  / Saber and show the latch structure predicts which scenarios trip which models;
  a warrant-typology × operation-type grid showing the two latches release
  independently.
- **Venue:** NeurIPS/ICLR safety, ICSE/FSE/ASE (SE-safety), USENIX/IEEE S&P,
  "Agents in the Wild" workshop.
- **Strategic bonus:** near-zero overlap with the ALIFE LBA — opens a second
  audience rather than competing with the current submission.
- **Best-fit bundle:** GS-1 Gaps 1 + 4 + 7 (mechanism + harness-isolation +
  missionless baseline).

### Direction B — Environment-sensitivity / disposition benchmark
- **Field-named gaps (verbatim):** *Agents Explore but Agents Ignore* (arXiv
  2604.17609, verified): "outcome-driven metrics like pass@k cannot distinguish
  adaptive reasoning from rigid plan execution"; "richer methods for measuring
  environmental curiosity [beyond solution injection] are needed." Disposition-
  benchmark need is composite (assembled from Szeider 2509.21224 + Behavioral
  Fingerprinting 2509.04504 single-turn-only + PROBE 2510.19771 task-success-only).
  AgentAtlas (2605.20530) concedes its labels were "produced by a single Claude
  Opus 4.7 checkpoint" and unvalidated against real tool traces.
- **Our evidence:** E1 five-motif matched-pair IS the discriminator pass@k lacks;
  E7 IS the real-substrate, un-fakeable-DV cross-model disposition profile the
  fingerprinting line wants but does not have.
- **Closing experiments (package + widen roster):** release E1 as a
  responsiveness benchmark (motif generator + external verifier + control
  protocol); publish E7 as an initiative/responsiveness leaderboard over 6–10
  models incl. open-weights.
- **Venue:** NeurIPS Datasets & Benchmarks, ICLR, COLM.
- **Caveat:** benchmark venue is crowded; OP3's "need" must be synthesized from
  three papers rather than cited from one.

### Direction C — Environment properties as controlled variables (geometry) ★ highest ceiling / flagship
- **Field-named gap (one hard, verbatim):** eco-evo neuroevolution (arXiv
  2302.09334, verified): episodic RL "prevents the study of both niche
  construction and eco-evolutionary feedback, which require that populations
  alter their environment and that those changes influence subsequent evolution."
  UED admits it cannot say **which** parameters matter (novelty-quantification
  arXiv 2502.05726: "underspecified nature of environment parameters"). The
  headline "no orthogonal *property* taxonomy, only difficulty/regret/novelty" is
  **surveyor-inferred** from PAIRED/ACCEL/CENIE structure, not field-named — flag.
- **Our evidence:** E8 proposes the five-dimension taxonomy AND has built 4 knobs
  (horizon / temporality / rigidity / plasticity); E8 plasticity-vs-rigidity is a
  direct persistence-vs-reset toggle answering the eco-evo gap head-on; the
  **non-monotone gap-dose curve** (1 gap 3/3 → 5 gaps 0/15) frontally contradicts
  UED's monotone-complexity premise with no co-adaptation confound; E6 separates
  solicitation from realization.
- **Closing experiments (NEW apparatus — highest cost):** factorial grid of the
  five E8 knobs on the POSIX substrate with difficulty held constant → first
  behavior-shaping *property* response surface; in-vitro niche-construction
  (rigid EROFS vs writable-persistent home).
- **Venue:** NeurIPS/ICLR (UED/autocurricula), ALIFE / *Artificial Life*
  (open-endedness), RLC.
- **Strategic note:** most theoretically ambitious, closest to the quine
  open-endedness long game; needs the most new experiments.

### Direction D — Causal grounding for context / environment engineering
- **Field-named gaps (verbatim, practitioner + academic):** *Agent READMEs*
  (arXiv 2511.12884, fetch-verified): "future work should conduct controlled
  experiments to isolate which context features most impact agent success rates."
  *EurekAgent* (arXiv 2606.13662, verified): "we call for environment engineering
  as a core research direction" — with no isolating ablation of its own.
  Anthropic's "Effective context engineering" blog presents filesystem-as-context
  as best practice with **no controlled experiment or baseline**.
- **Our evidence:** E1 (3 families, 0/85 control) is exactly the controlled
  isolation of structure from prompt and from model that the field folklore
  lacks; E4 explains why over-specified prompts suppress useful autonomy; E8
  supplies a design-relevant dose-response.
- **Closing experiment (re-skin existing data):** re-run the E1 matched pair with
  a CLAUDE.md/README context-file as the treated axis to speak the field's
  vocabulary while keeping the un-fakeable DV.
- **Venue:** ICSE/FSE/MSR (SE-agents), practitioner audience.
- **Note:** lowest theoretical ambition, cleanest "they asked, we answer."

### Direction E — ALIFE minimal agency without selection / precariousness (the LBA→full-paper continuation)
- **Field-named gaps:** Nisioti et al. 2024 "From Text to Life" (arXiv 2407.09502,
  peer-reviewed, verbatim): LLM-agent capability is "currently mostly designed
  top-down instead of being evolvable and emerging, driven by intrinsic
  motivation, and complex couplings with their ecosystem." Enactivist
  precariousness lineage (*BioSystems* 2022) holds agency requires precarious
  self-maintenance — E1 is a counterexample. **Structural finding:** the OEE
  lineage (Taylor #5; Spore.fun 2506.04236) *assumes selection is the engine*;
  "structure-driven behavior without selection" is a gap the field has **not
  framed** — E1 challenges the unstated selection/precariousness premise.
- **Our evidence:** E1 + E4 (coupling-driven + intrinsic-drive, no top-down
  objective); E6 (substrate resistance = a non-biological precariousness); E8
  plasticity (Taylor's "medium" / stigmergic persistence).
- **Closing experiment:** agency-from-coupling dose-response across E8 knobs in a
  precariousness-free substrate.
- **Venue:** *Artificial Life* journal, ALIFE conference — the LBA's existing home.
- **Note:** the natural full-paper extension of the current LBA; more
  position/framing character; lower ML-venue impact.

## Recommended sequencing (not a single paper)

1. **Direction A first** — hottest field, verbatim + multiple named gaps, we hold
   the mechanism they lack, lowest apparatus cost, direct existence proof against
   a standing assumption, opens a NEW audience disjoint from the LBA.
2. **Direction C as the flagship / ambitious second** — aligned with the quine
   open-endedness long game; budget for the new factorial-geometry experiments.
3. **Direction E as the LBA continuation** — extends the current submission's
   home venue; mostly reframing of held evidence.
4. **Directions B / D** — lowest marginal cost (package/re-skin existing data);
   good workshop or fast-follow papers, or fold B's benchmark into A/C as an
   instrument section.

Do **not** bundle all directions into one manuscript — that repeats the LBA's
"ambition outruns the page budget" failure at full-paper scale.

## Citation-hygiene carry-forward (verify before any becomes load-bearing)

- **Safety:** SafeLife 1912.01217, AUP 1902.09725, relative reachability
  1806.01186 — believed correct, NOT fetch-verified. **"Krakovna 2021" CORRECTED
  2026-07-05:** 2101.12509 (*Challenges for Using Impact Regularizers*,
  SafeAI@AAAI-21) is by **Lindner, Matoba, Meulemans — Krakovna is not an
  author**; this entry previously said "title/authors verified," which was
  itself wrong (fetch-verified correction:
  `plans/deixis-release-novelty-sweep-2026-07-05.md` § Citation corrections).
  The "specified externally" reading is fair; the verbatim quote remains a
  small-model paraphrase — re-read the PDF before quoting. Unconfirmed IDs: ClawsBench 2604.05172,
  AgentGuardian 2601.10440, Constraint Decay 2605.06445, WebGuard 2507.14293.
- **Task-free/evals:** "Look Before You Leap" (2605.16143) is NOT goal-free
  exploration (trains exploration *for* tasks) — citing as goal-free is an error.
  Evaluation-awareness IDs 2605.23055 / 2504.20084: PDF text not extracted, do
  not lean on. 2510.14548 "clean quotes" were a fetch confabulation — the real
  paper is far weaker.
- **Geometry:** POBAX 2508.00046, Environment Scaling 2511.09586, Jiang 2312.03126
  fetched but not ID-verified; 2511.12706, 2505.20659 search-only. Gomez
  three-mechanism wording (2510.05192) may be fetch-reconstructed.
- **Harness:** Saisubramanian & Zilberstein 2102.07017 (AAMAS 2021) verified.
  Search-only: 2511.03100, 2510.14548, 2604.16548, 2605.17641.
- **ALIFE:** Taylor 1507.07403, Nisioti 2407.09502, Hamon 2402.10236, ASAL
  2412.17799, Spore 2506.04236, TerraLingua 2603.16910 verified. Unverified:
  2607.01047, 2603.14771, 2605.13625.
- **Terminology collision (from the 2026-06-18 sweep, still binding):**
  "structural elicitation" already used by Da et al. 2025 for a prompting
  strategy — disambiguate or rename in any new paper.

## Raw subagent outputs

Preserved in the session task transcripts (2026-07-04); this file is the merged
synthesis. Per-neighborhood full text: GS-1 safety, GS-2 task-free/evals,
GS-3 environment properties, GS-4 harness, GS-5 ALIFE.
