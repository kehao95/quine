# Gap Survey — the "AI slop rots codebases" space (2026-07-05)

> Mode: gap survey (field-first), five neighborhoods, fetch-verified anchors.
> Feeds the SE projection of `structure-two-latch-action-gate` and the
> [`defensive-coding-as-route-around`](../../../incubator/discovered-sparks/defensive-coding-as-route-around.md)
> spark (umbrella: why AI slop rots codebases). Operational note: one research
> pass encountered prompt-injection text embedded in fetched pages; recognized
> and ignored.

## Headline verdicts

1. **The authorization-semantics gap is UNOCCUPIED** (top-ranked open problem,
   unarticulated in the field). Industry built approval *mechanics* (Claude
   Code permission modes/auto mode: ~93% of permission prompts accepted
   without careful review); academia benchmarked capability; **nobody varies
   the wording/referential specificity of a cleanup grant as an independent
   variable.** Our T1/T2 result is the first datapoint in an empty cell.
2. **The permission-vs-uncertainty confound is the decisive experiment only a
   two-latch apparatus designs cleanly.** A specific instruction ("delete
   foo(), lines 40–55") simultaneously grants permission AND certifies
   deadness; myopia accounts predict evidence suffices, deference accounts
   predict permission suffices — nobody noticed the confound. The 2×2
   (epistemic proof-of-deadness × object-level grant) separates them; only
   the both-open cell fires under two-latch.
3. **No over-refusal benchmark for code actions exists** — deletion-freeze
   looks like success (nothing broke) in every existing eval harness. An
   XSTest-for-code-edits suite falls out of our apparatus naturally.
4. **The degradation evidence itself is CONTESTED — do not build on "AI
   degrades code" as established fact.** Vendor-scale mining (GitClear
   2024–26: churn +15%, block duplication +81% 2023→2026, copy/paste 15.7% vs
   moved 3.8%, legacy-code updates −74%) vs the two most rigorous designs
   finding null-to-positive (Stray et al. arXiv:2509.20353, 39 devs/26k
   commits/2y, "virtually unchanged"; GitHub's own controlled study 2024,
   +2.5–4.2% quality). Our stance: the mechanism predicts the **additive
   signature** (duplication up, moved-code down, deletions down) — which is
   exactly GitClear's observed profile — without needing the contested
   aggregate "worse" verdict.

## Key measured anchors (fetch-verified)

- **Broad/catch-all exception handling is the #1 AI code smell** — Liu et
  al., Debt Behind the AI Boom (arXiv:2603.28592): 302.6k AI-authored
  commits; 41,374 instances (8.5%). The empirical anchor for fallback
  accretion.
- **Agents are scope-conservative and removal-spec-dependent:** Refactoring
  Runaway (arXiv:2605.22526): agents do LESS unrequested refactoring than
  humans (21.4% vs 36.7%); Agentic Refactoring (arXiv:2511.04824): refactor
  in 26.1% of commits, low-level renames, median LOC delta −15.25; CodeTaste
  (arXiv:2603.04177): agents execute specified removals well but **fail to
  discover** them — the discovery/execution split maps to our epistemic
  latch needing canonicality adjudication it cannot self-supply (fam_dedup).
- **Qu et al. Overeager (arXiv:2605.18583) is the mirror-image experiment:**
  manipulates consent explicitness, finds 0.0%→17.1% out-of-scope destruction
  when consent stripped — same variable, opposite failure mode (over-action
  under vague consent vs deletion-freeze under vague authorization). Must be
  differentiated in the first paragraph that mentions authorization.
- Counter-domain caveat: security literature finds LLM code UNDER-defensive
  for boundary/null checks in security contexts (arXiv:2409.19182) —
  "over-defensive" is not a monolithic construct; scope claims carefully.

## Nearest competitors (none occupies the conjunction)

| account | level | predicts authorization-wording effects? |
|---|---|---|
| GitClear "workflow incentives" (2026) | economic | no |
| Qu et al. consent ablation | experimental, safety register | presence, not specificity-net-of-scope |
| Lavi pre-action legitimacy; Uchibeke decision-vs-authorization (arXiv:2603.20953) | theory preprints, security framing | different axes, no maintainability link |
| Zhu et al. Volume-Quality Inverse Law (arXiv:2605.02741) | empirical scaling | no |

**The conjunction mechanism-account + maintainability outcome +
authorization-wording prediction is open.**

## Terms of art (adopt / differentiate)

Adopt: "code churn" (GitClear operationalization), "copy/pasted vs moved
code" (our fixed-point prediction targets exactly this ratio),
"refactoring discovery vs implementation" (CodeTaste). Differentiate:
"overeager/out-of-scope" (mirror image). Use "AI slop" only as colloquial
hook, never as construct. SATD exists but measured on LLM-based software,
not AI-authored code (unit mismatch).

## Revised bridge-experiment design (supersedes the spark's v1 ladder)

The decisive form is the **2×2 confound-splitter**, not the specificity
ladder alone: {epistemic proof of deadness supplied (coverage/static-analysis
evidence) | withheld} × {object-level grant | none}, fixed codebase, fixed
dead-code target; DV = deletion executed vs additive workaround vs freeze.
Two-latch: fires only in both-open. Myopia-only: evidence cell fires.
Deference-only: grant cell fires. Plus the specificity ladder (none /
policy-level / capability-grant / object-level) as the applied arm. Testbeds:
Knip-style dead-code corpora. This also generates the over-refusal-for-code
benchmark (problem 3) as a by-product.

## Ranked open problems (field's, importance × novelty)

1. authorization-semantics gap (unarticulated) — ours to take
2. permission-vs-uncertainty confound (unarticulated) — only our apparatus splits it
3. over-refusal benchmark for code actions (inverse named by Qu et al.)
4. mining-vs-experiment contradiction in degradation evidence (infrastructure)
5. refactoring discovery vs execution (named, CodeTaste) — epistemic-latch reframe
6. dead-code/SATD accrual attributable to AI authorship (measurement gap)
7. two-factor gate theory for destructive action (partially named, wrong axes)
