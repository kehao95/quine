# Manuscript Plan — narration-invariance-eval

Argument order and section skeleton for the
[`narration-invariance-eval`](README.md) dossier. **Working lead domain: Arm A
(CoT / System-2 reasoning faithfulness)** after the deep-research novelty-source-check
([`plans/dr-lead-domain-report.md`](plans/dr-lead-domain-report.md)). Arm B
(mech-interp) is the fallback. RAG and code-synthesis stay as §4 transport
corroboration.

## Core hypothesis

A reasoning trace's correctness decomposes into (a) a **load-bearing** component
carried by the task's functional structure and (b) a **narrated** component
carried by familiar surface dressing (names, objects, real-world intuition). The
two are causally distinguishable: load-bearing correctness is **invariant** under
a math-graph-preserving, dressing-scrambling environment shift; narrated
correctness **drifts**. The narrated component is non-trivial specifically in
regimes where the model leans on dressing to bridge steps the formal path alone
does not — making the phenomenon **difficulty- and model-regime-dependent, not
universal**. The operator's contribution is the discrimination test plus the
3-arm design that isolates dressing-dependence from the generic reskin fragility
(GSM-Symbolic regime); the pilot suggests the operator may be non-vacuous in at
least one model and maps a tentative regime boundary.

## Abstract (draft v2 — methodological framing, post Gemini novelty/importance review)

Chain-of-thought traces are evaluated by whether they reach the right answer,
but a trace can be correct because it rode familiar surface dressing rather than
the task's functional structure. We introduce a **topology-preserving
cross-environment invariance assay** that causally isolates this
dressing-dependence. Unlike GSM-Symbolic-style work, which measures generic
accuracy variance under input reskinning, our 3-arm design (original /
real-reskin control / nonce-isomorph) uses the **B-vs-C delta** to separate
load-bearing semantic reliance from generic reskin fragility: in a 4-model pilot
the control arm broke 0% of correct traces while the nonce-isomorph broke 17.4%
in one model — a delta that argues against pure reskin fragility but remains
confounded by nonce-OOD. The assay suggests a non-trivial dressing-dependent
subset in one model and none in three frontier models, indicating the phenomenon
may be **regime-gated** (emerges where the model
leans on dressing to bridge steps the formal path alone does not); the same
assay surfaces its own open confound (dressing-dependence vs nonce-OOD), which a
real-word-nonsense control resolves. The contribution is the causal
**evaluation primitive** — a scalpel that names and isolates a
previously-suspected-but-unisolated shortcut of the reasoning process — and an
honest map of the difficulty regime where it applies, not a universal law.

## Future work

- **Decisive confound control (D arm):** a real-word-nonsense dressing
  (e.g. "glomerates") — real English tokens with no real-world intuition — to
  separate narration-channel dependence (breaks on D) from nonce-OOD fragility
  (breaks only on nonce C). Single experiment that decides the interpretation.
- **Harder-regime elicitation:** the phenomenon may emerge in the invariant
  models on harder multi-step problems where dressing scaffolding is
  load-bearing; map the difficulty threshold per model.
- **Worst-where-cleanest (P4):** a real blinded judge rates trace legibility;
  correlate with per-trace invariance. The trace-length proxy used in the pilot
  is uninformative.
- **Legibility-pressure causal test (P6):** base-vs-RLHF'd matched-model
  comparison — does legibility optimization manufacture the decoupling?
- **Transport:** RAG citation attribution and program-synthesis spec-match as
  second/third domains for the same operator.

## Lead-domain working decision (provisional)

| Arm | Verdict | Main experiment |
|---|---|---|
| **A — CoT / System-2 reasoning faithfulness** (o1-style) — **WORKING LEAD** | **Provisional lead.** The current source-check leaves the D1–D4 intersection apparently unoccupied: closest behavioral kin *Lanham* breaks D2 (injects mistakes); closest principle kin *Triangulation* established D3-as-principle but on internal circuitry; *Illusion of Thinking / GSM-Symbolic* test capability not attribution. Functional mechanism formally verifiable in math/logic → environment shift clean. | semantic-isomorphism shift on MATH / GSM8K (translate word problems into meaningless tokens preserving the math graph); Bamboogle. |
| **B — Mechanistic interpretability** — **fallback** | Held in reserve; the Othello-GPT experiment if Arm A's proves unworkable. | Othello-GPT probe-vs-board-state; IOI under token permutation. |

**Novelty narrowings to carry (from source-check):** (1) *Triangulation* already
established cross-environment invariance as causal-gold-standard → the paper's
D3 novelty is **application-only** (external trace + topology-preserved shift),
not the principle; (2) *CausalDPO* uses the identical phrase for a different task
— distinguish; (3) the one-paragraph *Lanham* defense (they break logic / D2, we
hold it) is the first reviewer-attack answer.

## Reviewer Defense — D2 vs Lanham (drafted, for §6 Related Work)

**Status:** first-draft prose, synthesized only from material already checked
into the source-check
([`plans/dr-lead-domain-report.md`](plans/dr-lead-domain-report.md) §"Direct
Trace Perturbation: Anthropic's Faithfulness Tests (2023)" and the D1–D4
table). Lanham et al. 2023 is the one 2025–2026-adjacent citation the README's
Next Useful Work already flags as **well-established** (unlike Triangulation /
CMI / LITL / CausalDPO, which still need a primary-source pass) — so this
paragraph does not have to wait on that pending verification. This closes
README Next Useful Work item 4; it does not touch the restart gate.

> The closest prior art a reviewer will reach for is Lanham et al. (2023,
> Anthropic), which also intervenes directly on the generated chain-of-thought
> text (satisfying D1) and measures an invariance metric — the probability the
> final answer is unchanged after perturbation (D3, partially). But the two
> operations Lanham et al. use to perturb the trace — "Early Answering"
> (truncating the trace mid-way) and "Adding Mistakes" (injecting a logical or
> mathematical error into one step) — change what the trace *computes*: they
> break the functional topology of the reasoning graph itself, not merely its
> surface dressing (violates D2). Their resulting question is therefore "does
> breaking the logic break the answer?" — a test of whether the trace is
> causally load-bearing at all — not "does changing legibility while holding
> the logic fixed reveal causal attribution?", which is what a
> topology-preserving environment shift is built to answer. Lanham et al. also
> do not formalize a causal legibility-vs-load-bearing correlation (D4): their
> inverse-scaling finding (larger models produce less faithful, more
> plausible-looking post-hoc rationales) is suggestive of the Legibility
> Penalty but is never tested as a causal claim. One line: Lanham breaks the
> graph to test faithfulness; this operator holds the graph fixed and breaks
> only the narration to test the same question causally.

**Not yet done:** resolving the source-check's `[cite: 13, 14]` interaction-log
markers into an actual Anthropic paper citation (title / arXiv or blog ID) —
the report only carries citation markers, not a resolved bibliography entry.
That resolution is still Next Useful Work item 3 (primary-source pass) and
must land before this paragraph is submission-ready.

## Page Budget (7-page two-column + references)

| § | Section | Pages | Content |
|---|---------|-------|---------|
| 1 | Intro: narration-satisfiable evaluation | ~1 | The failure class: an eval signal satisfied by description inflates against a causal test. Worst-where-cleanest as the hook. Three one-line instances (circuits, citations, specs). |
| 2 | The operator: cross-environment invariance | ~0.75 | IRM-grounded shift: perturb the readable channel, hold the mechanism fixed. Load-bearing = invariant; narrated = drifts. Grounded briefly; not a treatise. |
| 3 | Lead experiment — CoT (Arm A) | ~2.5 | Semantic-isomorphism shift on MATH/GSM8K: translate word problems into meaningless tokens preserving the exact math graph. Measure trace-invariance under the shift. Predict a narrated-trace subset whose answers drift despite functional topology held constant. One-paragraph *Lanham* D2 defense here. |
| 4 | Transport corroboration | ~1 | One compact paragraph each: RAG counterfactual-passage swap (narrated-attribution rate), code-synthesis spec-operationalization shift. Same structural prediction; no full runs needed at submission. |
| 5 | The attractor: worst-where-cleanest | ~0.75 | Literature-level prediction: celebrated explanations correlate negatively with load-bearing. Testable on a sample of published explanations. |
| 6 | Related work / limits | ~0.5 | IRM, causal-abstraction (Geiger), attribution-faithfulness, spec-code eval. What each gets vs what this operator adds. |
| 7 | Conclusion | ~0.25 | One operator, three fields, one falsifiable structural prediction. |
|  | References | +unlimited |  |

## Decisions blocking a draft

1. **resolve D-control**: run the real-word-nonsense dressing arm before treating
   the pilot as narration-channel evidence rather than nonce-OOD fragility.
2. **map the struggle zone**: find a task/model regime where frontier models are
   near the capability edge before making an importance claim.
3. **source-check citations into checked-in references** for Arm A's adjacent art:
   Turpin 2023, Lanham 2023, Apple *Illusion of Thinking* 2025, GSM-Symbolic
   2024, Triangulation, CausalDPO, and IRM.
4. **freeze the concrete artifact** for the restart experiment: MATH/GSM8K
   semantic-isomorphism shift, Bamboogle, or another math-graph-preserving task.

## Open question (honest)

The single sharpest risk: "why does the readable-channel shift buy more than a
synchronous activation-patch?" Patching shows the feature is in a causal path to
the output; it does **not** show the feature is in the causal path *for the
reason the interpretation says*. The §3 experiment must make that distinction
visible — otherwise a reviewer reduces the result to existing interventional
interpretability. The worst-where-cleanest literature-level correlation (§5) is
the prediction that existing interventional work cannot produce and is therefore
the strongest novelty defense.

## What this plan is not

- not a three-domain operator paper (too thin per domain at 7 pages)
- not a claim the operator is validated by the home substrate (not cited)
- not a mech-interp methods paper (the contribution is the eval operator, not a
  new circuit-finding technique)
