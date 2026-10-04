# ALIFE 2026 Camera-Ready Review Checklist - fp283s1

This checklist turns the accepted-paper reviewer feedback into concrete
camera-ready and poster work.

Source review record:
[`alife-2026-fp283s1.md`](alife-2026-fp283s1.md)

Decision: `Accepted as Poster`
Camera-ready deadline: `2026-06-21`
Page limit: `8` pages of body (Intro through Appendix), references and
acknowledgements excluded. Body currently ~7.3 pages after the edits below;
total PDF 9 pages including references. Headroom ~0.5-0.7 page.

## Status (2026-06-06)

- Items 1-5 (ALIFE grounding, clean CRDT removal, front-facing claim audit)
  **applied and PDF rebuilt**; all new ALife citations resolve, no undefined
  references, body within the 8-page limit.
- Item 6 (model-choice rationale) **reverted / deferred** at the authors'
  direction: the real reason is resource limits on a multi-model matrix, which
  they prefer not to state. The model-choice question is instead deferred to the
  cross-model experiments below, which make a written rationale unnecessary. The
  closed-source-limitation enhancement was reverted too (the limitation
  paragraph is back to its original minimal form), since open-weight models in
  the planned matrix will reframe reproducibility.
- `2026-06-20` v2.3 draft note: the manuscript now reintroduces a minimal
  closed-source replication limitation ("exact replication by independent
  parties is not guaranteed") without adding a model-choice rationale. Treat
  this as the current local draft state, not as a completed cross-model
  absorption.
- Cross-model / larger-N / factorial / formal-coding experiments are tracked as
  planned supplementation below, now with a concrete design.
- The June `2026-06-06` to `2026-06-07` supplementation wave was later
  invalidated and deleted. The underlying mistake was procedural: we started
  paper reruns under a newer runtime/harness before first checking whether the
  original paper-defined prompt, tool, help, and other agent-visible surfaces
  still matched. No June supplement readout remains countable evidence for this
  paper.
- New default for any future supplementation: run and record a paper-runtime
  compatibility audit first, then start fresh reruns only after that audit
  passes. Canonical rule:
  [`../../../../experiments/PROTOCOL.md`](../../../../experiments/PROTOCOL.md).

## Verification Notes

The following items appeared on an earlier checklist but are already present in
the submitted manuscript. Do not repeat:

- Term definitions (task closure, peer-awareness, quantitative anomaly signals,
  qualitative artifact signals) — already in Introduction paragraph 4
- Limitations paragraphs (sample size, single model/task, informal coding,
  Step D contamination) — already in Discussion / Limitations section
- Peer-awareness coding criteria — already in Method, with the explicit
  disclaimer that this is not formal inter-rater coding
- Existence proof / bounded threshold language — already in Introduction and
  Discussion
- **All stigmergy / ALife references already exist in `references.bib`**
  (Grassé 1959, Theraulaz & Bonabeau 1999, Holland & Melhuish 2001,
  Heylighen 2016, Danchin 2004, Camazine 2001, Odling-Smee 2003,
  Maturana & Varela 1987). They simply were not yet cited in the body. No new
  literature research was required — the work was wiring existing citations into
  prose.

## Completion Rule

- Check an item only after the manuscript section has been updated and the PDF
  rebuilt.
- If an item is intentionally not pursued before camera-ready, leave it
  unchecked and add a short rationale beside it.
- New experiments are allowed only when they close a named reviewer gap; do not
  reopen the paper as a generic robustness project.

---

## Must Fix Before Camera-Ready

### 1. ALIFE Grounding — Introduction ✅ DONE

- [x] Added an ALife motivation paragraph as the new opening of the
  Introduction: frames the paper's question as a classic ALife question
  (minimal environmental information for emergent collective organization),
  grounded in stigmergy and self-organization
  (Camazine, Grassé, Theraulaz & Bonabeau, Heylighen) before narrowing to LLM
  agents.

### 2. ALIFE Grounding — Related Work ✅ DONE

- [x] Added a new opening subsection "Stigmergy and Indirect Coordination in
  Artificial Life" before the LLM-focused subsections: covers Grassé,
  Theraulaz & Bonabeau, Holland & Melhuish (minimal conditions), Heylighen
  (no mutual awareness), Camazine (reactive self-organization), and Danchin
  (inadvertent social information). Establishes the novelty boundary: classical
  stigmergy is reactive; our agents abductively recover a peer model.
- [x] Added a sentence at the end of "Positioning Our Contribution" using
  Reviewer 2's framing: LLM MAS has advanced *direct* coordination; this paper
  isolates the *indirect* / stigmergic regime in capable agents.
- [x] Replaced the bare term "CRDTs" in the CodeCRDT description with
  "observation of shared replicated state" (part of item 4).

### 3. ALIFE Grounding — Discussion ✅ DONE

- [x] Added a new subsection "Implications for Artificial Life" connecting the
  closure threshold to ALife: quantitative signals close a self-organization
  feedback loop (Camazine), static artifacts are inadvertent social information
  (Danchin) but close no loop, and the LLM substrate adds abductive peer
  recovery beyond Holland & Melhuish's reactive minimal agents.

### 4. CRDTs — clean removal ✅ DONE

- [x] Introduction enumeration: "CRDTs" → "replicated shared state" (citation
  to CodeCRDT preserved, no dangling reference).
- [x] Related Work CodeCRDT description: "uses CRDTs for observation-driven
  coordination" → "coordinates through observation of shared replicated state".
- Note: the cite key `@CodeCRDT2025` and the system's proper name "CodeCRDT"
  remain (a system name, not an unexplained concept); no reader-facing bare
  acronym remains.

### 5. Front-facing Claim Audit ✅ DONE

Reviewer 3 flagged unqualified strong claims (esp. "existential unawareness is
experimentally isolatable"). The same phrasing recurred on every front-facing
surface, so this was audited as a set, not just in the conclusion:

- [x] **Abstract** (`metadata.yaml`): reframed the results as "an existence
  proof in a single substrate, task, and model configuration"; "there is a
  threshold" → "we observe a threshold".
- [x] **Introduction**: "experimentally isolatable" → "...isolatable in our
  setup"; "framing is unnecessary" → "not always necessary" (existence-proof
  scope).
- [x] **Related Work**: "is experimentally isolatable and that a threshold
  exists" → "can be experimentally isolated and that a threshold appears, in our
  conditions"; the open question "recover a representation of those peers" →
  "infer their existence" (defensible given the informal coding caveat).
- [x] **Conclusion**: "...experimentally isolatable: within this substrate,
  task, and single model configuration, we can cleanly remove...".

Net effect: the contribution still reads as a clear demonstrated finding, but as
a bounded existence proof rather than a universal law.

### 6. GPT Model Family Rationale + Closed-Source Limitation ⬜ DEFERRED

Author direction (2026-06-06): do **not** state a model-choice rationale now.
The real reason is resource limits on running a full multi-model matrix, which
the authors prefer not to write. The better fix is empirical — the cross-model
experiments below directly answer Reviewer 1's "why only GPT" concern and make a
written rationale unnecessary.

- [x] **Reverted** the Method rationale sentence (back to the plain
  configuration sentence).
- [x] **Reverted** the Discussion/Limitations closed-source enhancement (back to
  the original "All runs use GPT-5.4. Whether the threshold generalizes ...
  remains untested"); open-weight models in the planned matrix will reframe the
  reproducibility point anyway.
- [ ] **After cross-model experiments**: rewrite the model paragraph to report
  the model set (turning "why GPT" into "we test N families"), and note that an
  open-weight model (e.g., DeepSeek) supports reproducibility — converting the
  closed-source weakness into a strength.

---

## Planned Experimental Supplementation

Full design is canonical in [`./cross-model-plan.md`](cross-model-plan.md)
(model set, tiered conditions, run counts, comparability, compatibility-audit prerequisite,
absorption path). Key facts for camera-ready scoping:

- A cross-model probe already exists (`8.15-cross-model-replication`) and leans
  **negative** for Opus 4.5 — but under the older pre-A-D prompt rungs, so it is
  not yet a clean test. The threshold may be model-dependent; do not assume
  cross-model = confirmed generality.
- Batch runs are **gated on a paper-runtime compatibility audit** under current
  HEAD. Do not reopen supplementation until prompt, runtime prompt, tool/help,
  workspace-default, sidechannel, and authority-bearing surfaces have been
  checked against the original paper contract.
- Larger-N, the factorial rung (8.21), and formal transcript coding fold into the
  same campaign.

**Camera-ready decision (`2026-06-21`)**: do **not** cram the cross-model study
in. Ship the bounded single-model existence proof (current state); route the full
cross-model + larger-N + factorial study to an extended version. Reasons: the
~2-week window, the pending compatibility audit, the tight page budget, and — decisively
— the uncertain outcome (the Opus negative) which should not be rushed into the
archival proceedings.

Immediate execution gate: the next supplementation attempt must begin with the
compatibility audit above, not with fresh pilots. If the audit reveals runtime
or tool drift, repair that drift first and treat any probe run before the audit
as methodology-only rather than paper evidence.

### Restart Status (`2026-06-07`)

- compatibility gate is **passed enough to restart pilots**, but the first
  restart wave uncovered one deeper runner bug: shared world state had split
  into per-agent sandbox-local copies because the generated `world.json` still
  carried a host absolute `state_dir`.
- that split-state bug is now fixed and manually smoke-verified on the repaired
  `/runtime/world/state` lane; canonical note:
  [`2026-06-07-final-harness-pilots.md`](2026-06-07-final-harness-pilots.md)
- the early `053421 / 053858 / 054154` "final harness" pilots are therefore
  invalid methodology, not retained evidence
- current retained post-shared-state-fix read:
  - `8.19` `gemini-3.1-pro-preview` A (`20260607-060250-...`):
    **accepted** on a real shared-state lane
  - `8.17` `gemini-3.1-pro-preview` B (`20260607-060416-...`):
    **negative / stalled** — duplicate prefix, resets, then re-entry into the
    per-agent-limit basin
  - `8.19` `gpt-5.4` A (`20260607-061216-...`):
    **accepted** on the repaired lane through file-mediated coordination
    (`claim.txt`, `note.txt`, `part1.txt`)
  - `8.17` `gpt-5.4` B (`20260607-061006-...`):
    **negative** — shared-budget anomaly noticed, then converted into reset /
    mistrust rather than cooperative inference
  - exact-paper-config caveat:
    - `gpt-5.4` `8.19` `xhigh` (`20260607-060049-...`) stayed in coordination
      prelude and never reached a first `world get`
- current execution read:
  - the runner semantics themselves are now finally correct again
  - we now have a meaningful retained A/B matrix on the repaired lane for two
    model families
  - further supplementation is no longer about fixing the lane; it is about
    whether we want to spend effort on an exact `gpt-5.4 xhigh` A rerun or
    move to C-lane model-dependence tests
  - a fresh `gpt-5.4 xhigh` unlimited A/B/C/D probe did recover A under current
    HEAD, but the prompt audit showed that only the task prompts are now
    restored; the provider-visible runtime system prompt is still materially
    different from April, so exact paper-equal re-acquisition is still blocked
    at the runtime-prompt layer

---

## Should Fix If Page Budget Allows

(~0.5-0.7 page headroom remains.)

- [ ] One sentence distinguishing conventional stigmergy (biological / simple
  agents) from this paper's LLM-agent contribution (complex reasoning,
  long-horizon planning) — partly covered by the new Related Work and
  Discussion subsections; add only if a sharper one-liner helps.
- [ ] One sentence explaining why qualitative artifact signals support
  peer-awareness but not stable closure (one-shot perturbation vs. recurring
  feedback) — partly covered by the new ALife Discussion paragraph.
- [ ] Short transcript excerpts in the appendix for the peer-awareness /
  closure distinction.

---

## Poster Checklist

- [ ] Poster title and framing lead with the accepted contribution: information
  ablation under existential unawareness.
- [ ] Include the four-condition ablation ladder.
- [ ] Include the closure-threshold result: `12/12` success with quantitative
  anomaly channels, `0/8` legitimate closure under pure zero perception.
- [ ] Include the R10 semantic-misalignment case as the memorable qualitative
  finding.
- [ ] Keep limitations visible but concise: small N, single model family,
  simple task, informal peer-awareness coding.
- [ ] Make the ALIFE connection explicit on the poster: indirect coordination
  through environmental traces, connecting to biological stigmergy.

---

## Do Not Do By Default

- Do not add broad cross-model experiments unless a specific reviewer gap cannot
  be addressed honestly without them (now tracked above as a deliberate choice,
  not an accident).
- Do not scale sample size just to appear more rigorous unless the schedule can
  absorb the analysis and claim updates.
- Do not turn the camera-ready into a new paper about protocol emergence or
  material channels.
- Do not re-add limitations paragraphs — they are already present.
