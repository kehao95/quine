# Experiment Escalation Plan

This note is the durable planning surface for follow-up experiments after the `quine-lifelike` manuscript reframing.

It exists to answer one practical question:

> when should we just rerun, when should we repair an experimental surface, and when should we redesign the experiment family entirely?

The paper should not treat all uncertainty the same. Some gaps are about numerical stability; others are about causal design. More runs only solve the first class.

## Current Paper Center

The current manuscript is built around four experiment families:

- `1.1 Ariadne` for real death plus continuity-operator contrast
- `2.6 Legacy Protocol Ablation` for the clean persistence hinge
- `2.8 Endogenous Index Lineage` for endogenous environmental-organ birth
- `3.1 Habitat Guardian` for maintenance-layer contrasts

Supporting only:

- `2.5 Long Now`
- `2.7 Archivist Index`
- `2.9 Maintenance Pressure`

The decision rule for follow-up work is:

- if the paper needs cleaner frequencies or cleaner confidence around an already valid contrast, use **Level 0**
- if the paper is vulnerable because the prompt or harness leaks the answer, use **Level 1**
- if the paper needs a stronger causal claim than the current design can legitimately support, use **Level 2**
- if the work would open a new claim track or likely feed a later paper, use **Level 3**

## Execution Status (`2026-04-12`)

This plan has now moved beyond the original Level 0 question. The highest-value redesign surfaces for the current paper were actually built and run, so the control plane needs to be read as an executed record rather than as a still-hypothetical queue.

### Standing provider default

- the lanes this paper's retained evidence ran on — `.env.copilot-gpt-5.4`,
  `.env.copilot-gpt-5.4-high`, `.env.copilot-gpt-5.4-xhigh`, `.env.gpt-5.4-pool`,
  `.env.claude-sonnet-4-5`, `.env.kimi` — were purged from the repo in the
  2026-05-24 provider sweep. Every occurrence of those names in this file is a
  **provenance label on an executed run**, not a sourceable path
- for new live waves, use the repo default named in
  [`experiments/PROTOCOL.md`](../../../../experiments/PROTOCOL.md)
  § Default Model Policy (`profiles/gpt-5.5-codex-oauth.env`)
- that default is **not** the lane the retained evidence ran on, and the original
  lane can no longer be reproduced, so a new wave here is a comparison lane
  rather than a continuation: report it as its own lane instead of pooling it
  with the existing paper-facing results
- within a wave, do not silently escalate supplemental runs to a higher-effort
  variant; keep one lane per wave so added replications strengthen the story
  rather than opening a profile confound
- when a run wave deviates from this default, record the reason in the experiment `README.md` and in any paper-facing summary that absorbs the result

### Operational repairs completed

- installed `coreutils`, so `timeout` / `gtimeout` are available for the experiment runners that require them
- smoke-tested `.env.copilot-gpt-5.4`, `.env.copilot-gpt-5.4-high`, `.env.claude-sonnet-4-5`, and `.env.kimi` directly with `./quine '1+1'`
- `.env.gpt-5.4-pool` was returning `HTTP 502` on `2026-04-12`, so live GPT-5.4 follow-up used `.env.copilot-gpt-5.4-high` as the operational fallback
- repaired `2.6`, `2.8`, and `3.1` runners so parallel launches get unique `runid`s instead of colliding on the same second-level timestamp
- repaired `2.8` analysis so `artifact reuse` means inherited reuse from a predecessor generation, not same-generation create-and-use
- repaired `3.1` analysis so `launchd` / post-exit guardian evidence is no longer under-classified as foreground-only

### Family-level stop gates after execution

- `2.6 Legacy Protocol Ablation`
  - the old `A/B/C` matrix still remains the causal hinge
  - but the planned redesign controls are now also executed:
    - `neutral-lineage` under persistent habitat removed durable artifact birth entirely
    - hardened `sham-artifact` preserved residue birth but drove successor uptake to `0/3`
    - `birth-uptake` pilot externalized birth files but `0` uptake files
  - decision:
    - keep `2.6` as the clean persistence hinge
    - also use the new control family as the strongest direct response to the `this is just file I/O / note-passing` objection
- `2.8 Endogenous Index Lineage`
  - the planned Claude `A x3 / B x3` Level 0 wave is now filled
  - Level 1 / Level 2 follow-up is now also partly executed:
    - `hidden-future` reruns are complete for `A x3 / B x3`
    - `surprise-transfer` pilot is complete for `A` and `B`
  - stable read after redesign:
    - persistent lineage (`A`) still supports endogenous birth under `hidden-future`, but inherited reuse drops to an unstable `4/30` observed generations
    - wiped lineage (`B`) keeps inherited reuse at `0/30`, even though artifact birth can still appear
    - `surprise-transfer` shows `0/5` late-phase transfer uptake under both `A` and `B`
  - decision:
    - `2.8` remains a phenomenon anchor for endogenous environmental-organ birth
    - later uptake is now bounded more tightly: possible under persistent aligned futures, but not yet a robust transfer organ
    - `2.6`, not `2.8`, should carry the main causal persistence argument in the paper
    - do not write `2.8` as a general transfer result
- `3.1 Habitat Guardian`
  - live GPT-5.4 mortality-cue follow-up did not stabilize at Level 0
  - the planned same-prompt hidden-horizon redesign is now executed across multiple model families
  - stable redesign read:
    - `H1` live-tamper is clean Tier 2 across Claude, Kimi, and Copilot
    - valid `H2` post-exit runs are clean Tier 3 across Claude, Kimi, and Copilot after widening the Copilot patience window
    - `H0` still externalizes baseline guardians, so it should not be treated as a zero-externalization control
  - decision:
    - `3.1` no longer depends only on a mortality-wording gradient
    - the hidden-horizon redesign should now carry the front edge of `3.3`
    - the old prompt-gradient family remains supporting context and negative control only

### Practical paper consequence

- treat `2.6` as the clean persistence hinge
- write the new `2.6` control family as evidence that writable persistence and file survival are each insufficient without lineage salience and semantically usable inherited structure
- keep `2.8` as corroborating phenomenon evidence, but write it around bounded uptake rather than raw accuracy collapse or general transfer
- let `3.1` pivot from a prompt-gradient story to a same-prompt horizon story

## Execution Status Table (`2026-04-12`)

This table is the fast status surface for the current paper interval. It distinguishes work that already changed the manuscript from work that remains planned only.

| Experiment family | Planned role | Executed this interval? | What actually happened | Current paper effect | Continue now? |
|---|---|---:|---|---|---:|
| `2.8 Endogenous Index Lineage` | Existing main-text evidence plus uptake-boundary redesign | Yes | Level 0 thickening, `hidden-future`, and `surprise-transfer` are now all executed; birth persists more easily than later uptake, and transfer uptake is still `0/5` under the pilot shift | Strengthens a bounded organogenesis claim while ruling out a stronger general-transfer reading | No |
| `3.1 Habitat Guardian` | Existing main-text evidence plus same-prompt redesign | Yes | `H0/H1/H2` horizon-aware harness built and executed across Claude, Kimi, and Copilot; valid `H2` runs are now clean Tier 3 post-exit positives | Converts `3.1` from a prompt-gradient result into cleaner horizon-pressure evidence | No |
| `1.1 Ariadne` | Existing main-text evidence, spot-check only | Partly | Additional runs exist, but no new paper-facing absorption occurred | No current manuscript change; still functions as the clean continuity-operator contrast | No |
| `2.5 Long Now` | Supporting-only corroboration, spot-check only | Partly | Additional runs exist, but no new paper-facing absorption occurred | No current manuscript change; still functions as a richer lineage phenomenon, not a causal hinge | No |
| `2.6 Legacy Protocol Ablation` | Existing clean hinge plus control-family redesign | Yes | `neutral-lineage`, hardened `sham-artifact`, and `birth-uptake` pilot now exist and materially changed the evidence surface | Strengthens the paper's response to writable-persistence and file-existence objections | No |
| `3.1 same-prompt hidden-horizon redesign` | New contribution experiment | Yes | `H0/H1/H2` family is built, valid, and paper-usable | Now the cleanest causal front edge for `3.3` | Adopt |
| `2.6 neutral-lineage / sham-artifact redesign` | New contribution experiment | Yes | Both redesigns are built and interpretable; lineage salience and artifact semantics now have bounded negative controls | Now the cleanest direct answer to the `file I/O` objection | Adopt |
| `Unified Viability-Pressure Factorial` | New contribution experiment, broader program | No | Remains conceptual only | Would most strongly elevate the paper's method claim, but is likely a later paper lane rather than a submission-stabilization task | Low for the current draft |

### Reading rule for the status table

- `Yes` means the work changed the live evidence surface and should already be reflected in experiment `README.md` and paper prose if relevant.
- `Partly` means runs happened but have not changed the paper's evidence hierarchy; they should be either absorbed explicitly or treated as spot checks only.
- `No` means the item is still only a plan and must not be described as completed work anywhere in the paper or control plane.

## Escalation Levels

### Level 0: Stabilization Reruns

Use this when:

- the design is already defensible
- the main issue is sample size, provider variance, or an underpowered frequency claim
- the current paper claim would not need to change if the new runs look like the old ones

What changes:

- model choice, seed date, or repetition count
- no prompt change
- no harness change
- no capability-surface change

What the paper may claim afterward:

- stronger counts
- narrower wording around prevalence or variability
- cleaner same-model or same-condition distributions

What it cannot fix:

- demand characteristics
- prompt leakage
- confounded manipulations
- ontology claims that current design never tested

### Level 1: Surface Repair

Use this when:

- the intended causal variable is right, but the prompt or harness reveals too much
- a control is too strong or too weak
- a condition changes more than one agent-visible surface

What changes:

- prompt wording
- timing of perturbation
- observation window
- capability disclosure
- artifact hygiene or analyzer rules

What stays fixed:

- the underlying paper question
- the experiment family identity
- the main interpretation lane

What the paper may claim afterward:

- cleaner causal language
- stronger response to reviewer objections about wording or framing
- better separation between habitat pressure and narrated instruction

### Level 2: Focused Redesign

Use this when:

- the current experiment cannot support the desired claim even with better prompts
- the manipulated variable needs to move from prompt semantics into habitat physics
- one phenomenon needs to be split into birth, uptake, sham, or transfer subtests

What changes:

- condition matrix
- runner logic
- hidden versus visible information structure
- what counts as the unit of evidence

What the paper may claim afterward:

- a stronger causal statement
- a more reviewer-robust mechanism claim
- a better distinction between convenience, communication, and true environmental organization

Risk:

- this may produce a cleaner paper, but it can also invalidate direct comparison with the current evidence backbone

### Level 3: New Experiment Family

Use this when:

- the new design would test a broader program direction rather than tighten the current paper
- the result would likely become a new section, appendix family, or future paper lane

Examples:

- branching subjective worlds
- population commons or signaling ecology
- autopoiesis / stronger autonomy
- explicit `fork` / branch-selection lineage operators

Default rule:

- do **not** do Level 3 work for the current submission unless the paper is already stable

## Paper Backfill Rules By Level

### After Level 0

Update:

- experiment `README.md` current evidence set and bounded read
- run counts, percentages, or frequency sentences in manuscript prose
- result figure captions if they report counts
- `manuscript-plan.md` only if the strongest evidence hierarchy changes

Do not update:

- title
- core claim
- discussion framing

### After Level 1

Update:

- experiment `README.md`
- prompt/harness provenance in the experiment control plane
- `manuscript-plan.md` if a cleaner control now replaces an older one
- method prose in Section 2
- any result paragraph that depended on the old surface distinction

Also record:

- which objection the repair addressed
- which older runs remain canonical, supporting, or retired

### After Level 2

Update:

- experiment `README.md`
- paper dossier `README.md`
- `manuscript-plan.md`
- result section structure if needed
- discussion boundary language
- figures and tables if the evidence unit changed

Also decide explicitly:

- whether the new design replaces the old evidence line
- whether the old runs become historical-methodological evidence only

### After Level 3

Update:

- broader-program maps first
- current manuscript only if the new family materially changes venue strategy

Default:

- keep Level 3 out of the current paper and route it into future-work or a new paper note

## Current Priority Queue

| Priority | Experiment family | Why it matters now | Best starting level |
|---|---|---|---|
| 1 | `3.1 Habitat Guardian` | current `3.3` remains the most reviewer-exposed section, and Level 0 did not stabilize the mortality-aware frequency story | Level 1 |
| 2 | `2.8 Endogenous Index Lineage` | Level 0 is complete; now mainly a bounded phenomenon anchor whose claim language must stay disciplined | Stop unless a later uptake-specific objection appears |
| 3 | `2.6 Legacy Protocol Ablation` | already the clean hinge and now closed at the current paper boundary | Stop |
| 4 | `1.1 Ariadne` | already sufficient for current claim; only reopen if continuity frequency becomes a point of contention | None by default |
| 5 | `2.5 Long Now` | rich but expensive; useful corroboration, not current causal backbone | None by default |

## Detailed Plan By Experiment Family

## `3.1 Habitat Guardian`

### Current weakness

The present paper claim is valid but reviewer-exposed because prompt semantics and ecological pressure move together:

- `condA` mentions restoration and disturbance
- `condB` adds explicit death wording
- `condD` removes all maintenance framing

That means the current result can still be attacked as a prompt-semantics gradient rather than a habitat-pressure gradient.

### Level 0 plan

Goal:

- stabilize the current runtime frequency picture
- determine whether `condB` is reliably stronger than `condA/C` or just noisier

Run matrix:

- `copilot-gpt-5.4 condB` x 4
- `copilot-gpt-5.4 condA` x 2
- `copilot-gpt-5.4 condC` x 2
- no additional `condD` unless the runner changes

Success criteria:

- at least `3/4` clean positive `condB` runs
- `condA` remains positive at the lighter layer
- `condC` remains capable of post-exit maintenance

Paper backfill if successful:

- add frequency language to `3.3`
- keep current claim: stronger pressure recruits a stronger layer

Paper backfill if not successful:

- weaken `3.3` to: disturbance-aware prompts recruit externalized maintenance; mortality-aware pressure is suggestive but unstable

Status after `2026-04-12` execution:

- this stop rule was **not** met under the live GPT-5.4 follow-up
- `condB` remained mixed across `Tier 3`, `Tier 1`, and setup-failure outcomes
- recommended next move is now Level 1, not more Level 0 reruns

### Level 1 plan

Goal:

- keep the same claim family
- reduce prompt-side leakage

Repair direction:

- keep the invariant and disturbance schedule
- reduce explicit wording differences across `A/B/C`
- move more of the contrast into harness timing and observation horizon

Concrete repair:

1. keep one common mission text
2. manipulate whether tamper occurs:
   - before exit
   - after exit
   - not at all
3. remove the explicit sentence `This must remain true even after your death`

Desired read:

- externalization should emerge because post-exit disturbance is real, not because the prompt narrates death

### Level 2 plan

Goal:

- move the strongest claim from prompt semantics into habitat physics

Redesign options:

1. **same-prompt hidden-horizon redesign**
   - identical prompt across conditions
   - only the actual tamper horizon changes
2. **portable maintenance redesign**
   - explicitly disallow `launchd` or host-service installation
   - force maintenance to be expressed in workspace-level structures only
3. **multi-invariant maintenance redesign**
   - maintain two or three habitat invariants
   - test whether the system builds a general maintenance organ rather than a single-file patch

Recommended order:

1. same-prompt hidden-horizon
2. portable maintenance
3. multi-invariant maintenance

## `2.6 Legacy Protocol Ablation`

### Current weakness

The substrate ablation is strong, but the lineage framing is still explicit enough that a reviewer can say the prompt is already instructing handoff behavior.

The `nolineage` control is also too strong to be a neutral control; it actively discourages future-oriented artifacts.

### Level 0 plan

Goal:

- make the `A/B/C` matrix numerically thicker without changing the design

Run matrix:

- `copilot-gpt-5.4 B` x 2
- `copilot-gpt-5.4 C` x 2
- `claude-sonnet-4-5 B` x 2
- `claude-sonnet-4-5 C` x 2

Optional:

- one additional `A` per model only if a fresh positive anchor is needed

Success criteria:

- `B` continues to show birth without inheritance
- `C` continues to block durable birth
- clean slices remain free of multi-read confounds

Paper backfill:

- strengthen the condition matrix in `3.2`
- add model-family robustness language only if the qualitative split remains clean

Status after `2026-04-12` audit:

- no new Level 0 runs were needed
- keep `2.6` closed at the current paper boundary unless a reviewer objection demands a cleaner control family

### Level 1 plan

Goal:

- reduce demand characteristics without changing the question

Repair direction:

- rewrite lineage framing so it no longer explicitly tells the agent that future generations will inherit the same environment
- do **not** replace this with an anti-handoff instruction

Concrete replacement control:

1. current `prompt.md` becomes `lineage-explicit`
2. new `prompt.neutral.md`:
   - task may exceed one lifetime
   - workspace persists
   - no explicit statement that a successor will continue
   - no instruction against writing handoff artifacts
3. retire current `prompt.nolineage.md` as an exploratory sidecar, not the main contrast

Desired read:

- if artifacts still emerge in the neutral condition, the case for habitat recruitment strengthens

### Level 2 plan

Goal:

- distinguish persistence, semantics, and uptake more sharply

Redesign options:

1. **sham artifact control**
   - artifact file persists
   - semantics are scrambled or labels removed
   - tests whether successors depend on meaningful structure rather than mere file existence
2. **cross-fostering**
   - a successor enters a workspace prepared by a different lineage
   - tests whether behavior follows the environment rather than lineage-specific habits
3. **birth-versus-uptake split**
   - one subtest asks whether predecessors externalize
   - another asks whether successors preferentially use inherited structure

Recommended order:

1. sham artifact control
2. birth-versus-uptake split
3. cross-fostering

## `2.8 Endogenous Index Lineage`

### Current weakness

The phenomenon is strong, but future generations are still part of the visible prompt framing. That leaves room for a reviewer to call the result weakly assigned preparation.

### Level 0 plan

Goal:

- show that endogenous birth is not hanging on one canonical run

Run matrix:

- `claude-sonnet-4-5 A` x 3
- `claude-sonnet-4-5 B` x 3

Optional:

- one `copilot-gpt-5.4 A/B` sweep if cross-model symmetry becomes important

Metrics to summarize:

- first artifact generation
- reuse generations
- maintenance generations
- expansion generations
- query accuracy

Paper backfill:

- keep `2.8` as the phenomenon anchor
- add simple distribution language around birth and reuse

Status after `2026-04-12` execution:

- Level 0 is complete
- persistent lineage still supports endogenous birth plus occasional inherited reuse
- wiped lineage keeps inherited reuse at `0/3`, but birth and raw accuracy are more variable than the first-pass read suggested
- use `2.8` as a phenomenon anchor, not as the paper's main causal hinge

### Level 1 plan

Goal:

- reduce explicit prep salience while keeping the same query ecology

Repair direction:

- remove or soften lines that explicitly tell the generation future questions will arrive
- keep durable workspace facts unchanged

Concrete repair:

1. create a `hidden-future` prompt
2. current generation only knows its current query
3. harness still continues the lineage behind the scenes

Desired read:

- if artifacts still appear, the claim shifts from weak prep anticipation toward truly endogenous environmental organization

### Level 2 plan

Goal:

- separate artifact birth from true later uptake

Redesign options:

1. **surprise transfer**
   - early generations face one query family
   - later generations face adjacent but not identical tasks
   - tests whether artifacts generalize as organs rather than rote notes
2. **standardized uptake assay**
   - give later generations a fixed inherited artifact bundle
   - measure artifact-first querying versus raw archive traversal
3. **observer-side uptake scoring**
   - rely less on self-reported `used_artifacts`
   - score tape order, archive-touch count, and query path structure

Recommended order:

1. observer-side uptake scoring
2. hidden-future reruns
3. surprise transfer

## `1.1 Ariadne`

### Current weakness

The current experiment is already good enough for the paper, but it still foregrounds `exec` in the prompt rather than letting the lineage discover operator availability from the runtime surface.

### Level 0 plan

Goal:

- stabilize the clean operator contrast

Run matrix:

- `condA` x 4
- `condB` x 2

Current recommended use:

- only do this if the paper needs stronger continuity frequency language

### Level 1 plan

Goal:

- reduce prompt disclosure as a confound

Repair direction:

- hide `exec` from the prompt
- let the tool surface and runtime determine discovery

Desired read:

- continuity is not just following an explicitly narrated operator hint

### Level 2 plan

Goal:

- compare operator continuity to pure filesystem continuity

Redesign options:

1. no-`exec` generational restart with durable habitat only
2. dual-path continuity experiment:
   - operator path available
   - filesystem-only path available
   - compare which one is recruited under different mortality budgets

This is valuable, but not a current submission blocker.

## `2.5 Long Now`

### Current weakness

It is rich, but expensive and partially confounded by the interaction between external timeout and `exec`.

### Level 0 plan

Goal:

- lightly confirm the current richer phenomenon

Run matrix:

- `45s` x 1
- `90s` x 1

Use only if:

- the paper needs one more corroborating example for save-or-explore and regression loops

### Level 1 / 2 plan

Only consider later.

Possible redesign:

- disable `exec` entirely
- study pure filesystem continuity under external mortality only

This is interesting, but not the most efficient way to strengthen the present paper.

## Historical Pre-Execution Order

This was the ordering used before the `2026-04-12` execution pass:

1. `3.1` Level 0
2. `2.8` Level 0
3. `2.6` Level 0
4. `3.1` Level 1 if mortality-layer instability remained
5. `2.6` Level 1 if the `file I/O` objection remained live

It is retained only as provenance for how the wave was staged. The active post-execution guidance now lives in:

- `## Current Recommendation`
- `## Remaining Execution Recommendation`

Level 2 remains unnecessary unless:

- Level 1 still cannot stabilize the claim, or
- a reviewer-level objection clearly cannot be answered by surface repair alone

## Completion Checklist For Any Run Wave

After finishing a run wave:

1. classify every new run by evidentiary role or delete it
2. update the experiment `README.md`
3. update `## Paper Feeds` if relevance changed
4. update `manuscript-plan.md` if the evidence hierarchy changed
5. update section prose only after the control plane is current
6. retire stale frequency claims in the manuscript
7. if a redesign replaced an older condition, say so explicitly

## Decision Gates

Use these gates to decide whether to stop, rerun, or redesign.

- **Stop at Level 0** if the observed distribution stabilizes and the manuscript claim no longer overreaches
- **Escalate to Level 1** if the result remains reviewer-vulnerable because the prompt surface leaks the intended behavior
- **Escalate to Level 2** if even a repaired surface cannot support the desired causal statement
- **Route to Level 3** if the redesign would mainly support a future paper lane rather than tighten the current manuscript

## Current Recommendation

After the `2026-04-12` execution pass, the practical recommendation is now:

1. stop Level 0 work for `2.6`; the existing matrix is already enough for the paper-boundary claim
2. stop Level 0 work for `2.8`; the thickened Claude wave clarified the bounded claim, even though it narrowed it
3. treat `3.1` as the only family still worth changing if more experimental work is justified
4. if `3.3` needs stronger support, go to `3.1` Level 1 surface repair rather than spending more rerun budget on Level 0
5. leave `1.1` and `2.5` as they are unless a concrete reviewer objection appears

That keeps the program aligned with the paper's actual remaining risks:

- `3.3` still has prompt-surface and frequency-stability risk
- `3.2` now mainly needs disciplined claim wording, not more count-buying
- there is no immediate need to widen the paper beyond the current evidence ladder

## Objection-to-Intervention Map

Use this map before allocating time. It prevents spending rerun budget on objections that only a redesign can answer.

| Reviewer objection | Best response level | Why |
|---|---|---|
| `n` is too small; maybe this is just variance | Level 0 | more clean replications directly answer it |
| the effect only appears in one model family | Level 0 | cross-model replication is enough if the design is already clean |
| this looks like prompt wording, not habitat pressure | Level 1 | the surface must be repaired so fewer agent-visible facts change |
| this is just file I/O / note-passing | Level 1 or Level 2 | sham, uptake, or neutral-lineage controls are needed |
| `maintenance` is really just instruction-following | Level 1 or Level 2 | hidden-horizon or same-prompt redesign is needed |
| these are three different tasks, not one separable pressure story | Level 2 or Level 3 | only a unified or more tightly matched ecology can answer it |
| this is an interesting runtime demo, not ALIFE evidence | Level 0 plus Level 1 | stabilize the contrasts and reduce design leakage before widening scope |

## Executed Level 0 Wave Record

The planning templates above were used as pre-execution scaffolds. After the `2026-04-12` pass, the useful durable record is the executed wave summary below.

### Wave Record 1: `2.8` Claude Thickening

Executed matrix:

- existing canonical runs:
  - `20260309-190620-claude-sonnet-4-5-condA`
  - `20260309-191420-claude-sonnet-4-5-condB`
- new Level 0 runs:
  - `20260412-060627-claude-sonnet-4-5-condA`
  - `20260412-061736-claude-sonnet-4-5-condA`
  - `20260412-060627-claude-sonnet-4-5-condB`
  - `20260412-061736-claude-sonnet-4-5-condB`

Net read:

- persistent lineage (`A`) now shows endogenous birth in `2/3` Claude runs and inherited reuse in `2/3`
- wiped lineage (`B`) keeps inherited reuse at `0/3`
- wiped lineage no longer guarantees low accuracy or zero same-generation artifact birth

Decision:

- `2.8` Level 0 is complete
- keep `2.8` as the phenomenon anchor for endogenous birth plus possible later uptake
- do not use `2.8` as the paper's main causal persistence hinge

### Wave Record 2: `3.1` Live GPT-5.4 Follow-Up

Operational note:

- `.env.gpt-5.4-pool` was returning `HTTP 502` on `2026-04-12`
- live GPT-5.4 follow-up therefore used `.env.copilot-gpt-5.4-high`

Executed mortality-cue follow-up:

- `20260412-061034-copilot-gpt-5.4-high-condB-r2` -> `Tier 3`, host-level `launchd` guardian
- `20260412-061750-copilot-gpt-5.4-high-condB` -> `Tier 1`, foreground vigilance
- `20260412-061034-copilot-gpt-5.4-high-condB` -> setup failure

Net read:

- the mortality-aware condition still reproduces the strongest class of result
- the frequency story is mixed, not stabilized
- Level 0 therefore strengthened the bounded possibility claim but did not meet the earlier `3/4` clean-positive target

Decision:

- stop Level 0 work for `3.1`
- if stronger support is still needed, move to Level 1 surface repair

### Wave Record 3: `2.6` Audit-Only Closure

Executed action:

- no new runs
- re-audited the existing `A/B/C` matrix against current paper needs

Decision:

- `2.6` already satisfies the current paper-boundary stop rule
- leave it closed unless a specific reviewer objection requires neutral-lineage, sham-artifact, or uptake-only redesign work

## Remaining Execution Recommendation

If more experiment work is justified for the current paper, the next bounded move is:

1. `3.1` Level 1 same-prompt hidden-horizon redesign

Do **not** continue with:

- more `2.8` Level 0 thickening
- more `2.6` Level 0 thickening
- more open-ended `3.1 condB` reruns on the current prompt surface

Reason:

- `2.8` and `2.6` have already answered the count-thinness question
- `3.1` now has a design-defensibility problem, not merely a sample-size problem

## Detailed Backfill Checklist

After any completed wave, update in this order:

1. experiment `README.md`
   - current evidence set
   - bounded claim
   - do-not-collapse boundary
   - retained run list
2. experiment `## Paper Feeds`
3. `manuscript-plan.md`
   - only if the evidence hierarchy or strongest contrast changed
4. section prose
   - first counts and frequency words
   - then stronger/ weaker claim language
5. figures and captions
6. `submission-checklist.md` if the build-facing status changed

If a redesign replaces an older control:

- mark the older condition as historical or methodological
- do not let both old and new controls sit in the paper as if they have equal force

## Result Summarization and Rollup Rules

These rules are the manuscript-facing extension of `experiments/PROTOCOL.md`. They exist so that post-run evidence is aggregated consistently rather than reinterpreted ad hoc.

### A. Per-run summary is mandatory

For every retained run, record at least:

- `runid`
- experiment family
- wave id
- escalation level
- model / env file
- condition
- prompt or harness revision id
- setup status
- validity status: `clean`, `invalid`, `provider/harness failure`, `exploratory-only`
- primary readout
- evidentiary bucket
- one-sentence note on why it is retained

If the analyzer already emits these facts, link the analyzer output. Do not duplicate it manually unless the analyzer is missing a required field.

### B. Denominator discipline

Default manuscript denominator is:

- retained **clean runs** per condition

Always report separately:

- invalid runs
- harness impatience / timeout artifacts
- provider failures

Do not silently fold these into scientific negatives.

If a section uses a secondary denominator, state it explicitly:

- `2.6` may report generation slices as a diagnostic
- the paper-facing claim should still prefer run-level counts unless the section contract is revised

### C. Do not pool across changed surfaces

Never pool counts across:

- different prompt revisions
- different harness-timing regimes
- different analyzer/classification rules
- pre-redesign and post-redesign condition families

If a redesign changes the agent-visible surface, start a new count table.

Allowed pooled summaries:

- same prompt
- same harness
- same analyzer
- same condition semantics
- multiple runs across time or provider instances

### D. Canonical summary table schema

Every completed wave should produce one compact table with these columns:

| Condition | Clean runs | Invalid runs | Primary positive count | Negative count | Exploratory-only count | Canonical retained runs |
|---|---:|---:|---:|---:|---:|---|

Use experiment-specific positive criteria:

- `1.1`: continuity discovered / repeated self-replacement
- `2.6`: birth, reuse, or blocked-birth as appropriate to condition
- `2.8`: endogenous artifact birth and later reuse
- `3.1`: tiered maintenance outcome

### E. Canonical run selection rules

Per condition, retain at least:

- one strongest positive exemplar
- one strongest negative or null exemplar when it informs the contrast
- one methodological-defect exemplar if a defect changed interpretation

Delete or demote:

- near-duplicates that do not change the read
- thin provider-failure shells
- redundant bookkeeping runs already represented by a better retained run

### F. Paper update rules by evidence outcome

#### If Level 0 only thickened counts

Update:

- counts
- percentages if used
- qualifiers such as `often`, `sometimes`, `in X/Y runs`
- figure captions

Do not update:

- title
- abstract thesis
- section logic

#### If Level 1 repaired prompt or harness leakage

Update:

- Section 2 method wording
- the affected results paragraph
- `manuscript-plan.md` evidence hierarchy if the repaired condition becomes canonical
- discussion caveat language

Also mark:

- which older runs remain supportive only
- which old counts should no longer be quoted

#### If Level 2 replaced the control logic

Update:

- experiment `README.md`
- paper dossier `README.md`
- `manuscript-plan.md`
- section structure if needed
- any table that mixed old and new units

And decide explicitly:

- whether the old family remains manuscript evidence
- or whether it becomes historical-methodological context only

### G. Required writeback order

After the wave summary is written, write back in this order:

1. experiment `README.md`
2. experiment `## Paper Feeds`
3. `manuscript-plan.md` if hierarchy changed
4. section prose
5. figures / captions
6. submission checklist only if build-facing status changed

### H. Summary sentence template

When rolling results into the paper, prefer sentences of this form:

`Under <condition>, <phenotype> appeared in <x>/<y> clean runs; <z> additional runs were invalid or harness-limited and are reported separately.`

This prevents denominator drift and keeps invalid infrastructure outcomes out of the science claim.

## Highest-Value New Experiment

Two different notions of ``most valuable'' matter here.

### Highest-Value Practical New Experiment

If the goal is to improve the current submission with the least risk, the single best new experiment is:

**`3.1` same-prompt hidden-horizon redesign**

Why this one:

- `3.3` is the paper's most reviewer-exposed section
- the current design changes prompt semantics and ecological pressure together
- a cleaner maintenance result would noticeably strengthen the paper's central method claim

Core design:

- use one common mission text across all conditions
- keep the invariant and disturbance semantics constant
- manipulate only when tamper occurs relative to process exit:
  - no tamper
  - tamper while the seed process is still alive
  - tamper only after seed-process exit

Minimal condition set:

1. **No-tamper baseline**
   - same prompt
   - no later disturbance
2. **Live-horizon disturbance**
   - same prompt
   - tamper occurs before process exit
3. **Post-exit disturbance**
   - same prompt
   - tamper occurs only after exit

Hypothesis:

- one-shot completion should suffice in `no-tamper`
- foreground or lightly detached repair may suffice in `live-horizon`
- only `post-exit` should strongly recruit detached maintenance

Why it helps the paper:

- it turns the strongest maintenance result into a habitat-timing result rather than a wording result
- it keeps the claim exactly where the manuscript wants it: maintenance-layer shift under changed viability pressure

How to backfill:

- replace the current `condD vs A/C` front edge in `3.3`
- keep current `condB` as supporting evidence for even stronger post-exit layers, not as the only main positive
- revise Section 2 to describe maintenance pressure in terms of horizon and disturbance timing rather than explicit death narration

## Detailed Spec: `3.1` Same-Prompt Hidden-Horizon Redesign

This is the concrete Level 1 / Level 2 bridge for the current paper. It is the most valuable redesign because it attacks the strongest surviving objection to `3.3`.

### Problem it solves

Current `3.1` conditions co-vary too many agent-visible surfaces:

- mission wording
- disturbance wording
- death wording
- post-exit requirement wording

That means a reviewer can still say the paper observed a prompt-semantics gradient instead of a maintenance-pressure gradient.

### Core redesign principle

Hold constant:

- the mission text
- the invariant
- the disturbance sequence
- the model family within a run wave
- the analyzer notion of restoration

Manipulate only:

- when disturbance happens relative to the seed process exit

### Fixed task contract

Use one common prompt for all redesign conditions:

- establish `habitat/signal.txt = ALIVE`
- the habitat may be disturbed later
- do not tell the agent whether disturbance is while alive or after exit
- do not narrate death explicitly

The prompt may mention durability and possible later disturbance. It should not mention:

- `after your death`
- `after you exit`
- `future maintenance layer`
- any suggested mechanism such as daemonization

### Condition matrix

| Condition | Prompt | Tamper horizon | Intended reading |
|---|---|---|---|
| `H0 no-tamper` | identical | none | baseline completion cost without later pressure |
| `H1 live-tamper` | identical | all tamper attempts happen while seed process is still alive | foreground vigilance or light externalization may suffice |
| `H2 post-exit-tamper` | identical | first tamper only after seed process has exited | detached maintenance should be most favored here |

### Harness changes

#### Runner

Create a redesign-specific runner or extend `run.sh` with a clearly versioned mode. The redesign should not silently reuse old condition names.

Minimum new metadata fields:

- `horizon_mode`
- `prompt_revision`
- `tamper_schedule_revision`
- `seed_exit_epoch`
- `first_tamper_epoch`
- `tamper_after_exit`
- `horizon_validity`

Recommended names:

- `h0-none`
- `h1-live`
- `h2-postexit`

#### Tamper harness

Current `tamper.sh` always:

- tampers immediately after setup grace
- performs overwrite / delete / overwrite
- then kills the foreground process if still alive

The redesign should split that into horizon-aware modes:

1. `none`
   - do not tamper
   - do not force kill
2. `live`
   - require seed process to be alive at first tamper
   - if it exits before first tamper, mark `horizon_miss` and treat the run as invalid for this condition
   - after live tamper sequence, allow natural completion or existing stop contract
3. `postexit`
   - wait for natural seed-process exit up to a generous timeout
   - if the seed never exits before timeout, mark `horizon_miss` and classify as harness impatience, not scientific negative
   - only then apply the tamper sequence

Do not force the seed process to exit merely to create the post-exit condition unless a later redesign explicitly studies forced death.

#### Analyzer

The current analyzer outputs:

- `restored_count`
- `timeout_count`
- `externalized_evidence`
- tier classification

The redesign analyzer should add:

- `horizon_mode`
- `horizon_validity`
- `tamper_after_exit`
- `restored_live_count`
- `restored_postexit_count`
- `detached_restoration_observed`

### Validity rules

Count as **invalid**, not negative:

- setup never stabilized
- live condition but seed exited before first tamper
- post-exit condition but seed never exited before the declared patience bound
- provider or harness failure

Count as **negative** only when:

- the horizon condition was actually instantiated
- the run was otherwise clean
- no maintenance layer achieved the condition's positive criterion

### Positive criteria

For paper-facing summaries:

- `H0 no-tamper`: positive = invariant established without unnecessary externalization dominating the condition
- `H1 live-tamper`: positive = at least one live-horizon restoration under clean conditions
- `H2 post-exit-tamper`: positive = at least one clean post-exit restoration, with detached evidence preferred

For the strongest `3.3` claim:

- require `H2` to outperform `H1` specifically on detached or post-exit restoration, not merely on any restoration

### Recommended run sequence

Pilot first:

- `H0` x 2
- `H1` x 2
- `H2` x 2

Only after analyzer and validity rules stabilize:

- `H0` x 2 more
- `H1` x 3 more
- `H2` x 3 more

This gives a bounded `4 / 5 / 5` wave rather than an open-ended redesign spiral.

### Success conditions for adopting the redesign into the paper

Adopt as the new front-edge evidence if:

- `H2` produces clean post-exit restoration in multiple runs
- `H1` remains weaker or more foreground-bound
- `H0` does not show the same detached-maintenance prevalence
- invalid-rate is reported separately and is not dominating one condition

Do not adopt as main evidence if:

- all three conditions externalize at similar rates
- horizon validity is too low to interpret
- the redesign mostly measures harness patience instead of maintenance

### Paper backfill if successful

Update:

- Section 2: define maintenance pressure via disturbance horizon, not explicit death narration
- Section `3.3`: make the redesigned `H0/H1/H2` comparison the front edge
- current `condA/B/C/D`: demote to supporting or methodological context

Preferred claim after success:

- `maintenance-layer strength increases when the habitat must remain valid across longer disturbance horizons, especially beyond seed-process exit`

### Paper backfill if mixed

If `H2` is suggestive but not cleanly stronger:

- keep the redesign as methodological support
- keep the current `3.1` family as phenomenon evidence
- weaken `3.3` to disturbance-sensitive externalized maintenance, not a firm horizon-gradient claim

### Highest-Value Scientific New Experiment

If the goal is to contribute most to the paper's **core claim** rather than to the easiest manuscript repair, the best new experiment is:

**Unified Viability-Pressure Factorial**

Rationale:

- the paper claims that mortality, persistence, and maintenance are separable in one substrate
- right now that claim is demonstrated across three experiment families
- the strongest possible reinforcement would be a single ecology where the same base task is used while these pressures are independently manipulated

Core idea:

Use one lineage task that naturally benefits from:

- continuity under death
- durable habitat state
- post-exit maintenance

Candidate base task:

- generations must survey a hidden world for semantic content
- progress must accumulate across lives
- one habitat invariant must remain true during and after survey activity

Manipulated axes:

1. **Mortality / continuity path**
   - continuity operator available
   - continuity operator blocked
2. **Persistence**
   - writable persistent habitat
   - wiped habitat
3. **Maintenance horizon**
   - no disturbance
   - post-exit disturbance

Minimal factorial:

- `2 x 2 x 2 = 8` conditions

Key readouts:

- continuity transition rate
- checkpoint birth
- artifact uptake by successors
- raw rereading rate
- detached maintenance emergence
- invariant restoration rate

Why this is the highest-value scientific experiment:

- it directly tests the paper's core methodological claim in one ecology
- it reduces the ``these are just three different demos'' objection
- it would let the paper show orthogonalization rather than only parallel case studies

Why it is not the first recommendation:

- it is expensive
- it risks becoming a new experiment family rather than a quick manuscript reinforcement
- it may force a partial rewrite of the Results section and figures

Recommended status:

- treat this as a Level 2.5 / Level 3 bridge
- only start it if the current manuscript is already stable, or if reviewer feedback later demands a unified separability demonstration

## Unified Factorial Design Sketch

If the unified experiment is pursued, the design should satisfy the following constraints.

### Fixed surfaces

- same model family within a run wave
- same mission text
- same hidden world
- same artifact namespace
- same execution budget unless mortality itself is the manipulated variable

### Independent manipulated surfaces

- continuity operator presence
- habitat persistence mode
- disturbance timing

### Avoid

- changing both prompt semantics and harness timing in the same condition jump
- introducing explicit anti-handoff instructions as a control
- relying only on self-reported `used_artifacts`

### Required analyses

- observer-side artifact uptake
- duplicate-work rate
- checkpoint schema stability
- post-exit restoration verification
- failure-cell documentation, not only positive exemplars

### Success criterion

The unified experiment is successful only if at least one organization changes when one pressure changes while the others remain fixed. Without that, the experiment does not actually demonstrate separability.

## If Time Is Extremely Limited

If you can only do one bounded piece of new work for the current paper:

1. do **no new experiment work** if bounded claim language is acceptable in `3.3`
2. do `3.1` same-prompt hidden-horizon redesign first if the immediate worry is still causal defensibility or mortality-layer frequency stability

If you can do one ambitious new experiment for scientific payoff rather than submission stabilization:

1. do the unified viability-pressure factorial
