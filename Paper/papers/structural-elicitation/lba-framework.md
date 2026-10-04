---
status: lba-coherence-surface
date: 2026-06-22
role: the LEAN LBA surface — the small, evidence-sufficient claim and nothing more. Reconstructed 2026-06-19 to ship a CAUSAL/PHENOMENAL claim only; all MECHANISM (two-latch, appraisal-readiness, the normative latch, latched-potential, capture/scope, open-ended drive, prime=drive) is deferred to the full paper because it is argued from small-n probes + tape-reading, not yet LBA-grade matched-pair proof.
---

# LBA Framework — Structural Elicitation (lean / Tier-A)

> **The LBA ships a causal/phenomenal claim — WHAT, proven cleanly — not a mechanism.** Everything
> about *why/how* (the latched potential, the two-latch gate, appraisal-readiness, the
> information-preservation latch, capture-vs-scope, the open-ended drive, the prime=drive
> decomposition) is **deferred to the full paper** (§4) — it is argued from small-n probes +
> tape-reading, not yet matched-pair-grade. Keep the LBA small, concrete, and **not over-defensive**.
> Companions: [`./theory-framework.md`](theory-framework.md) (the deep mechanism theory — full
> paper), [`./claim.md`](claim.md), [`./lba-abstract-draft.md`](lba-abstract-draft.md),
> evidence home: [`../../../experiments/active/structural-elicitation/`](../../../experiments/active/structural-elicitation).

## 0. The claim (one sentence)

> **In the ground state ($E_{\text{task}}=0$), the realized structure of a frozen LLM agent's environment is a SUFFICIENT CAUSE of matched, externally-verified functional behavior — cross-model, robust beyond code prose for discrepancy structure, realized on the POSIX substrate, and modulated across competition regimes when explicit tasks are introduced.**

This is a *phenomenal / causal* claim. We report **what** happens, cleanly. The full multi-regime phase diagram and mechanism ecology are explored in the full paper.

## 1. The evidence (Tier-A only: matched-pair, un-fakeable DV, cross-model, adequate n)

**The phenomenon — the cross-model spine.** Matched pairs (structure arm vs an otherwise-identical control), clean task-free prime, scored only on **un-fakeable filesystem/exec events** (narration never credited). Three model families fire every structure arm; no control arm fires.

| motif | structure arm -> DV | control | gpt | deepseek | glm |
|---|---|---|---|---|---|
| rupture (`13.03`) | broken JSON -> valid parse | intact | 7/7 vs 0 | 5/5 vs 0 | 4/5 vs 0 |
| gap (`13.01`) | missing slot -> `field/5` created | complete | 7/7 vs 0 | 5/5 vs 0 | 5/5 vs 0 |
| topology (`13.06`) | dangling reference -> target written | closed | 7/7 vs 0 | 5/5 vs 0 | 5/5 vs 0 |
| semantic (`13.07`) | version conflict -> reconciled (higher-wins) | consistent | 13/14 vs 0 | 8/10 vs 0 | 8/10 vs 0 |
| affordance (`13.04`) | opaque executable -> run | inert | 7/7 vs 0 | 5/5 vs 0 | 5/5 vs 0 |

*Total: **96/102** structure-arm successes vs **0/85** controls. Semantic motif combines two deconfounded collision directions; scorers exclude 0-command infra-failures.*

**It is not merely a coding-prior artifact** for the discrepancy mode. Ported to non-code prose — a handbook referencing a file that does not exist — the same discrepancy creates the named target above matched controls across **four** families: gpt **10/13 vs 0/12**, and **claude-sonnet-4-6, deepseek-v4-pro, glm-5.2** each **5/5 vs 0/5**. This does not close the stronger non-linguistic gate; `13.24` remains a surface-content ablation, not an LBA claim.

**It is realized on the substrate** (optional panel). `13.13`: under a read-only mount the agent still issues a real `EROFS`-refused write **attempt** (6/7 gpt, 5/5 deepseek) — the behavior is realized at the POSIX layer, and solicitation is separable from permission.

**Competition Regimes & Attention Gating:** In the ground state ($E_{\text{task}}=0$), structural pull directs action. When an explicit, restrictive negative command is introduced (`conflict`), attention and action space are captured by the prompt (0/10 gap fill on gpt, 0/5 on GLM); under weak/underspecified prompting (`14.06`), structural pull remains intact (3/3 passive and active). Attention-trace audits (`14.07`) show that conflict agents still perceive the structural gap in context (5/5) but gate execution under normative task constraints.

## 2. Held out of the LBA (one line)

**Funnel (`13.05`) is a clean NULL** (the entrypoint binary runs in both arms — no *distinct* channeling); held out without claiming low-DOF channeling cannot exist.

## 3. A light ALife framing (framing, not a claim)

A minimal, environment-coupled dynamic: in the ground state, a frozen policy's functional behavior is supplied not by an external task or a reward but by the **structure of its environment**. Lineage to invoke: **affordance landscapes** (Gibson) / **behavior settings** (Barker) / **Umwelt** (von Uexküll) / **self-stigmergy**. Avoid (overclaim traps): *alive*, *autopoiesis* (this is allopoiesis), *morphological computation* (use "structural topology"), *niche construction*, *utility function* (use "substitute objective"), and the label *proto-agentic* when a direct description is enough.

## 4. Deferred to the full paper (the MECHANISM — NOT claimed or argued in the LBA)

These are real, recorded observations, but they are argued from **small-n probes + tape-reading**,
not matched-pair-grade — so they are full-paper material that needs **designed, rigorous proof**,
not exposition. (Owned in [`./theory-framework.md`](theory-framework.md); evidence in motifs
`13.18`–`13.25`.)

- **The mechanism of the drive** — "the prime supplies an (open-ended) drive; structure supplies the
  target" (`13.25`); the latched-potential framing (kinetic prior + structural gradient).
- **The two-latch gate** — Epistemic (is it a defect?) × Normative (am I allowed?) (`13.18`–`13.22`).
- **Appraisal-readiness** — only structure legible-as-a-defect fires; diffuse/global structure
  (clutter) does not until asserted as text (`13.18`/`13.19`).
- **The information-preservation latch** — add/fix the proven-defective, don't destroy/relocate an
  intact file without warrant; `cp`-not-`mv` (`13.20`–`13.22`).
- **Capture vs satiation** — specific tasking *scopes* the drive rather than *satiating* it (`13.25`).
- **The model axis** — initiative differs across models (needs ≥3 models + scale).
- **Origin** — is the normative latch Unix-RLHF-bound or architectural? (domain-strip; not run).

## Open (LBA submission)

- [x] re-lean the abstract/body to: phenomenon -> three-family spine -> one-line boundary -> light ALife
      (mechanism = future work).
- [ ] optional one figure/poster graphic: matched-pair bar chart (5 motifs x 3 families) plus small A2/14.06 boundary panel.
- [x] title decision (2026-06-22): *Structural Elicitation in a Frozen POSIX Agent:
      Environmental Affordance Fields for Artificial Agency*; metadata + Makefile + PDF aligned.
- [x] novelty paragraph (`prior-art-dr-report.md`); inline related-work positioning
      (Szeider / Ye cited + differentiated in manuscript Problem, 2026-06-26).
- [ ] LBA paper status update.
- [x] Gemini theory-spark session (2026-06-22): logged at
      [`references/gemini-theory-spark-2026-06-22.md`](references/gemini-theory-spark-2026-06-22.md);
      one LBA-affecting spark (Edelman sentence), four full-paper sparks.
- [x] S2/S3 follow-up (2026-06-22): removed the `proto-agentic` label from the LBA manuscript;
      `elic-31` / `14.07` tested Gemini's gravity/rocket counter-reading on existing `elic-22`
      tapes and supports gating (conflict: gap in view 5/5, verbalized 0/5, filled 0/5; free:
      verbalized 5/5, filled 5/5).
- [x] A2 sample/model strengthening (2026-06-22): `elic-22` scaled gpt core cells to `free` 9/10,
      `additive` 0/10, `conflict` 0/10 and completed GLM's additive cell (`free` 5/5, `additive`
      0/5, `conflict` 0/5). DeepSeek was attempted but unavailable.
