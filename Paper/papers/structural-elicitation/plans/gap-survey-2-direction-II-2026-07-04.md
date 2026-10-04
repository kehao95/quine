# Gap Survey Round 2 — Direction II (evidence-independent, premise-hunting)

Mode: **gap survey**, run under the *corrected* protocol
([`deep-research-query` contract](../../../../development/contracts/deep-research-query.md)
§ Evidence-independence rule, 2026-07-04). Round 1
([`../full-paper-gap-survey.md`](../full-paper-gap-survey.md)) was
evidence-anchored — it front-loaded our E1–E8 inventory and asked each agent
"which of our evidence bears on this gap," which biased discovery toward
evidence-adjacent gaps. Direction II ("environment as a control surface") came
out of that round without a real inversion insight (unlike I and III, whose
cores each reverse an unstated field premise).

**This round fixes that for Direction II specifically.** No evidence inventory
is given to the searchers. The mandate is to find, on each field's own terms,
the **unstated assumptions / taken-for-granted premises / questions no one is
posing** about environment, control, and agent behavior — premises that could be
**reversed into a thesis** (the shape of I: "avoidance must be engineered" →
it's already there; III: "agency needs selection" → structure alone suffices).
Rank by **importance × novelty**, not by whether we could close it.

## Neighborhoods (4 parallel Claude subagents, homogeneous each)

| # | Neighborhood | Reversible-premise hunt |
|---|--------------|-------------------------|
| II-1 | LLM-agent control / prompt & context engineering | is "controlling an agent = telling it what to do (instruction/prompt); the environment is passive context/scaffold" an unexamined premise? where does the field *locate the lever of control*? |
| II-2 | RL / unsupervised environment design / curriculum | is "an environment is parameterized by difficulty/solvability/reward; it is a curriculum, not a controller" an unexamined frame? is environment's role assumed epistemic (a task to solve) rather than dynamic (a force that steers)? |
| II-3 | Human behavior-design disciplines (choice architecture / nudge / persuasive tech / ecological interface design) | these fields already treat environment as a behavior *controller* — for humans. what is the unexamined premise about whether/why this does (not) transfer to artificial agents? is agent-facing environment design a named field or an unclaimed territory? |
| II-4 | Situated/distributed cognition + systems structure (stigmergy, extended mind, Conway's law, constraints-as-design) | is "structure *constrains* / *scaffolds* / *supports* behavior" held as distinct from "structure *directs* behavior," with only the weaker (constraint/support) claim made? is the directive/force reading the unclaimed one? |

Deliverable (identical): for each neighborhood, 3–6 findings ranked by
importance × novelty. Each finding = (1) the unstated premise or unposed
question, stated as a sentence the field would nod at without having written it
down; (2) evidence it is actually the field's operative assumption (a quote /
close paraphrase from a real source, or an honest "surveyor-identified pattern,
not a cited claim"); (3) **the inversion** — the thesis you get by reversing it,
and why that would be important + novel; (4) what kind of result would be needed
to earn the inversion (no reference to our specific apparatus — describe the
experiment class). Flag unverified citations. Explicitly separate "the field
says this is open" from "the field never thought to ask."

## Run ledger

| Launch | Neighborhood | Runner | Status |
|---|---|---|---|
| 2026-07-04 | II-1 agent control / prompt-context eng | Claude subagent | ✅ force≠information; attraction≠constraint; token-scramble falsification |
| 2026-07-04 | II-2 RL / UED / curriculum | Claude subagent | ✅ frozen/reward-free case central formalism can't express; deploy≠train |
| 2026-07-04 | II-3 human behavior-design disciplines | Claude subagent | ✅ substrate-neutral≠inherited-bias; unclaimed positive design discipline; nudge-vuln prior-art landmine |
| 2026-07-04 | II-4 situated cognition / systems structure | Claude subagent | ✅ **sharpest**: originates≠amplifies (direction in a goal-less agent) |

**Outcome:** II HAS a real inversion core (four independent fields converge one
step short of the same claim). Absorbed into
[`../full-paper-directions.md`](../full-paper-directions.md) § Direction II
(rewritten from "characterize the control surface" to the *originating directional
force* inversion, with 5 killer falsifications + prior-art to out-position).

**Verification flags (before any becomes load-bearing):** nudge-vuln cluster IDs
2505.11584 + 2509.25609 (ABxLab) fetched-confirmed; the PNAS "agents sensitive to
nudges" (doi 10.1073/pnas.2537030123) and Royal-Society Barker→virtual
(rstb.20230291) returned 403 — real venues, quotes are snippet paraphrase, verify
before citing. Ng–Harada–Russell policy-invariant reward-shaping cited from memory —
verify. UED observability-as-parameter downgrade (2508.14135, 2512.08139)
snippet-only. Suchman/Hutchins/Barker/Withagen-Rietveld quotes are secondary
paraphrase except Clark & Chalmers + Conway (verbatim-confirmed).
