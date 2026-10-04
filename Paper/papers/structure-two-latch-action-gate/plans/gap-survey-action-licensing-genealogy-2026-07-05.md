# Gap Survey — the Action-Licensing Genealogy (2026-07-05)

> Mode: gap survey (field-first), three strands. Fetch-verified (V) unless
> marked UNVERIFIED. Feeds the Related Work / framing of
> `structure-two-latch-action-gate`. Companion:
> [`gap-survey-ai-slop-2026-07-05.md`](gap-survey-ai-slop-2026-07-05.md).

## Strand A — AI-safety formal lineage: the complete arc

Arc: constrained MDPs (Altman 1999) → deontic machine ethics (Bringsjord
2006, Powers 2006 — internal gates, hand-installed) → safe-RL taxonomy
(García & Fernández 2015, JMLR: all safe RL = modify criterion or modify
exploration) → **Corrigibility (Soares/Fallenstein/Yudkowsky/Armstrong 2015):
the strongest anti-emergence statement in the lineage** ("despite default
incentives for rational agents to resist…") → Amodei 2016 (assumes-engineered,
but carries the earliest emergence seed: reward-uncertainty prior) →
interruptibility/off-switch (external-guard) → low-impact (Armstrong &
Levinstein 1705.10720) → shielding (1708.08611) → whitelisting (Turner 2018 —
**AF post, not a paper**) → relative reachability (1806.01186) → AUP
(1902.09725 — conservatism emerges from penalty *structure* even with random
auxiliaries, not from data) → SafeLife (1912.01217) → future-tasks
(2010.07877 — **third author is Richard Ngo**) → Lindner/Matoba/Meulemans
2021 (2101.12509 — **exact title includes "Negative" Side Effects**; §6.3
"asking the designer" = nearest internal crack) → Naiff & Goel 2303.03139
(post-LLM survey, still V−μP framing, **never invokes LLM priors** — the
bridge between literatures is unbuilt).

**Three cracks, all short of our claim:** Amodei's uncertainty-prior seed
(prior hand-installed); AUP's random-auxiliary emergence (from structure, not
data); Lindner §6.3 (learned via active querying, not pretraining). Nobody
anticipates restraint arriving free with instruction tuning.

**Post-LLM negative verdicts we MUST frame against (gate exists but leaks):**
ToolEmu (2309.15817: safest agent fails 23.9%), SafeAgentBench (2412.13178:
best proactive rejection 10%, five agents 0), **Tang et al. 2508.13465** (>98%
risk-knowledge but <26% applied — the knowledge-action gap), Policy-Invisible
Violations (2604.12177: 90–98% violation when policy state hidden;
"self-avoidance is not a substitute for explicit policy enforcement"). Our
framing: characterizable prior ≠ reliable guarantee (Not-claimed (a) already
covers this); these papers report leakage without characterizing the release
grammar — that is exactly our contribution.

## Strand B — empirical LLM-agent line: the intersection is unoccupied

**Internal + decomposed + conjunctive (E∧N) + veto + reward-free — no
occupant.** Nearest neighbors: Ask or Assume? (2603.26233 — dual-condition
but **epistemic∧epistemic and externalized**); Selectively Quitting
(2510.16492 — single confidence latch; our claim is a strict generalization);
Overeager (2605.18583 — strongest hint of an internal normative condition);
R-Judge/AgentHarm/XSTest lineage (internal but undecomposed);
Progent/guard-cluster (external by design). Only 2509.21224 (UNVERIFIED) has
a true mission-free baseline — and studies spontaneous activity, not
restraint. **Caveat: Overeager and SNARE are the same research thread —
"effective breadth two lines, not three" when we describe Field 2.**

## Strand C — institutional action-licensing: component-occupied, synthesis OPEN

| domain | two-condition (epistemic showing ∧ authority grant)? |
|---|---|
| 4th Amendment | **YES** — probable cause + neutral magistrate, independent vetoes |
| Military ROE | **YES** — PID "necessary but not sufficient" + engagement authority |
| Medicine (German law) | **YES, and NAMES it: the *Zwei-Säulen-Modell*** (two-pillar model) — no indication → impermissible even with consent; no consent → battery even if indicated |
| ITIL / ATC | partial (process flow / authority-half only) |
| Psychology (omission/status-quo bias) | NO — single-condition directional substrate |
| Philosophy | NO, but **Foot 1967 names and brackets the second latch**: "I shall leave aside the special case of giving permission, which involves the idea of authority…" — set aside for 60 years |
| **Horvitz 1999 mixed-initiative** | NO — **the instructive antithesis**: both ingredients fused into ONE compensatory expected-utility scalar where confidence buys out permission — the exact model our flat compensation sweep refutes for frozen LLM agents |
| Sheridan LOA | NO — authority continuum only |

**Cross-cutting hunts, all negative:** nobody asserts the shared grammar
across domains; LLM-internalization is component-occupied (Schramowski 2022
moral norms; 2605.30934; Cheung/Maier/Lieder PNAS 2025 — **LLMs show
stronger omission bias than humans, but single-condition, no two-latch, no
institutional grammar**); engineered conjunctive gates exist (RACG
2606.13884: need ⊆ state AND authorization ⊆ state, fail-closed — designed,
never measured as a prior).

**Nearest three to the synthesis + differentiators:**
1. *Governing Actions, Not Agents* (2606.26298, PARTIALLY VERIFIED — fetch
   fully before citing): same two-condition shape, same human-institution
   analogy — but **engineers** the gate as external infrastructure. Opposite
   locus: designed artifact vs discovered prior.
2. Cheung/Maier/Lieder PNAS 2025: emergent bias in frozen LLMs, but
   single-condition, no authority veto.
3. RACG 2606.13884: verified conjunctive fail-closed gate — designer-imposed.

**Verdict: the framing "instruction-tuned models internalized the
institutional action-licensing grammar (epistemic showing ∧ authority
grant)" is OPEN at synthesis level.** Scope discipline: our data licenses the
framing + the behavioral signature; it does NOT license a provenance claim
(internalization-from-text is a hypothesis the framing motivates, not a
result — consistent with Not-claimed (c)).

## Ranked open problems (field's) → our coverage

1. conjunctive-vs-graded internal decomposition (unarticulated) → **largely
   done** (elic-32/33; wave-2 pending)
2. reward-free naturalistic base rate (partially named) → **done** (the
   missionless prime + E4 is exactly it)
3. release grammar of the leaky gate (named as failure only) → **ours**
   (deixis-release; L4a battery designed)
4. LLM priors as substitute for impact regularizers (seeded, unbridged) →
   our Field-1 section becomes the bridge Naiff & Goel didn't build
5. cross-domain licensing theory + internalization (unarticulated; Foot
   bracketed it) → framing + signature only, scoped
6. cross-model variance explanation (SNARE 56/21 named, unexplained) →
   latch-default decomposition; L5 port
7. provenance base-vs-instruct (unarticulated) → **excluded, keep excluded**

## Citation hygiene to apply

Lindner title includes "Negative"; 2010.07877 add Ngo; whitelisting = Turner
2018 Alignment Forum post; "conservative agency" is Turner et al.'s coinage.
Do not load-bear without re-fetch: Bringsjord 2006 abstract wording, Bengio
2502.15657 body quotes, 2606.26298 full text, 2509.21224, 2511.08798,
2409.00557, the 26xx guard-cluster IDs, CJCS SROE primary.
