# Deixis-Release Novelty Sweep + Citation Verification — 2026-07-05

> Targeted prior-art sweep triggered by the elic-32 wave-1 T2 anomaly (the
> location-scoped warrant released destruction 7/10 while the strictly broader
> generic warrant released 0/10 — release tracks referential specificity, not
> authority scope). Feeds
> [`structure-two-latch-action-gate`](../../../candidates/structure-two-latch-action-gate.md)
> § Wave-1 result and the L4 adversarial-warrant layer. Web-sourced
> (WebSearch/WebFetch), single-agent sweep, five neighborhoods × multiple query
> formulations.

## Verdict — the asymmetry is claimable

No work found in five neighborhoods that decomposes permission/warrant text
into **scope** vs **referential specificity** as independent variables, or that
reports a broader authorization releasing *less* action than a
narrower-but-specific one.

**Pre-drafted qualifier sentence for the manuscript:**

> "To our knowledge, no existing work varies the referential specificity of an
> agent authorization independently of its scope, or reports that increasing
> specificity — even while narrowing or holding scope constant — releases
> destructive action more reliably than broadening scope alone; the nearest
> related result (Qu et al. 2026, consent-declaration ablation) varies whether
> authorization is stated at all, not its specificity net of scope."

## Neighborhood table

| Neighborhood | Nearest works | Verdict | Differentiator |
|---|---|---|---|
| LLM-agent nudge/sensitivity | 2505.11584 (Hypersensitive to Nudges, ICLR-26); 2509.25609 (ABxLab) | **CLEAR** | choice-attribute salience in tasked consumer decisions; never authorization text; vulnerability+inherited-bias framing confirmed |
| Choice-architecture / nudge theory | canonical taxonomy (defaults/salience/framing) | **CLEAR** | salience = which *option* is noticed, not how a granted *authority* is worded; no named specificity-vs-scope principle |
| Prompt-eng / instruction-following | 2601.18554 (Deconstructing Instruction-Following); 2502.10995 (indirect speech acts) | **CLEAR** | task-constraint compliance, not permission-grant structure, not destructive release |
| Agent authorization / permission systems | 2605.18583 (Overeager); 2605.28122 (SNARE); 2603.17170 (PAuth); 2603.28166; 2510.11108 | **PARTIAL-OVERLAP** | nearest analog = Qu et al. verbose/terse consent ablation: varies declaration **presence**, not specificity independent of scope; nobody inverts scope vs specificity |
| Pragmatics/deixis for LLMs | 2401.07078 (PUB benchmark) | **CLEAR** | comprehension benchmark only; no action-release connection |

## Citation corrections (load-bearing — apply everywhere these anchors appear)

1. **SNARE ID was wrong.** SNARE = **arXiv 2605.28122** ("SNARE: Adaptive
   Scenario Synthesis for Eliciting Overeager Behavior in Coding Agents", Qu et
   al. — same group as, but distinct from, 2605.18583). The 56%-framework /
   21%-model variance split is **verified verbatim** in 2605.28122.
2. **Krakovna et al. 2021 → Lindner et al. 2021.** arXiv 2101.12509
   ("Challenges for Using Impact Regularizers", SafeAI@AAAI-21) is by **Lindner,
   Matoba, Meulemans** — Krakovna is not an author. The "penalties must be
   specified externally" reading is fair; only the attribution was wrong.
   (Krakovna's own nearest works: 2010.07877 future-tasks; 1806.01186 relative
   reachability.)
3. **Overeager consent-stripping stat needs scoping.** 0.0%→17.1% (McNemar
   p=2.4e-4) is **GLM-4.6-specific, Claude-Code framework, 76-scenario phase-1
   pilot**; other models: MiniMax-M2.7 3.9%→21.1%, Sonnet-4.6 3.9%→15.8%. Do
   not cite as a framework-wide number.
4. **Selectively Quitting scope fix.** Verified title: "Check Yourself Before
   You Wreck Yourself: Selectively Quitting Improves LLM Agent Safety"
   (Bonagiri, Kumaraguru, Nguyen, Plaut). "Compulsion to act" is verbatim and
   **general across all 12 evaluated models**; "largely seen in the GPT family"
   verbatim but modifies the **quit-rate/helpfulness correlation** (Fig. 2c),
   not the compulsion itself. Keep the two claims separate.
5. **"Overeager explicitly scopes mechanism out" — soften.** No explicit
   scope-out statement confirmed; the paper simply does not investigate
   mechanism.

## Provenance

Sweep run 2026-07-05 by a web-research subagent; queries spanned
authorization-ablation, jailbreak specific-target framings, behavioral-econ
"specificity principle" phrasings. Anchors verified against arXiv abstracts
and paper text via fetch. Consumer: the two-latch candidate's promotion
checklist items "deixis-release novelty check" and "snippet-only prior-art
anchors verified".
