# Experiment Design Extension Plan — Structural Elicitation (E4–E9)

> **Status (updated 2026-06-20): EXECUTED.** E4–E7 have now been built **and run** on
> gpt-5.4-codex-medium — outcomes in **§Results** at the top (E8 multi-model excluded by
> request). The build-plan body below (§1–§6) is **retained as design provenance**: where it
> reads "not yet run / human-gated / nothing here authorizes spending compute", that posture
> applied **pre-execution** and is **superseded by §Results**. The doc began as an integrated
> methodologist synthesis of the audit + adversarial-verify pass on the five proposals; running
> remained human-gated and was authorized separately before execution.
>
> **Owner surfaces:** claim at [`./claim.md`](claim.md), manuscript plan at
> [`./manuscript-plan.md`](manuscript-plan.md), evidence home at
> [`../../../experiments/active/structural-elicitation/`](../../../experiments/active/structural-elicitation).
> Canonical harness template = `13.24-nonlinguistic-gap/` (run-container.sh + run.sh
> + analysis/score.py). Built: 14.01–14.04. Next free experiment number: **14.05**.
>
> **E9 added (2026-06-20): the fair, deconfounded Gate-2 retest.** E5's "SE-specific"
> read is **provisional** — it confounded `non-SE` with `both-latches-closed` (see E9
> card §4). E9 holds both latches open and varies only substrate; it either rescues
> non-SE generality or confirms SE-specificity cleanly. Treat the §Results "SE-specific
> discrepancy" line as **pending E9**, not settled.

---

## ✅ Results (executed 2026-06-20, gpt-5.4-codex-medium, single-model per the "except multi-model" scope)

All four experiments were built (`14.01`–`14.04`), smoke-validated, run, scored on
**un-fakeable DVs**, **tape-verified**, and interpreted. E8 (multi-model reproduction) was
excluded by request. Per-experiment detail + verified numbers live in each
`../../../experiments/active/structural-elicitation/14.0x/README.md` §Result.

| Exp | Dir | Outcome | One-line |
|---|---|---|---|
| **E4** | `14.01` | **Gate (1) OPEN** | Both non-linguistic axes failed: the `+x` affordance is **floored** (agents run opaque files regardless of the bit — control 9/10, even `chmod +x`-ing them themselves), the SHA-hole is **null** (0/10, not perceivable as a fillable gap). De-linguistification not achievable with these constructs. |
| **E5** | `14.02` | **Gate (2) confounded → overturned by E9** | E5 *read* the discrepancy/alignment pull as SE-specific (**0/30** on a diary date-conflict vs 13.07's **3/3** on SE config), but that fixture closed **both** latches (subjective defect + destructive resolution) — confounded. **E9 (`14.05`) overturns it.** Non-SE affordance fires but is **floored** by legible-task-presence (inert 6/6 = affordance 6/6, pointless 0/6). |
| **E9** | `14.05` | **Gate (2) REOPENED — discrepancy pull is NOT SE-bound (model-general)** | Holding both latches open and varying only substrate to non-code prose: **F1** dangling cross-ref **10/13 (77%)** vs **0/12** symmetric control (Fisher **p=0.0001**); **F2** dangling citation **10/11 (91%)** vs 0/11; **F3** broken-CSV **0/9** (clean NULL — its fix *mutates* an intact row, latch-2 closed, exactly as predicted). **Cross-model (n=5/arm):** claude-sonnet-4-6 fires F1 5/5 & F2 5/5; deepseek-v4-pro fires F1 5/5 & F2 5/5 (all controls 0) — not a gpt artifact. The pull is **objective-defect + additive/determinate-resolution bound**, substrate-general; SE merely supplied many such fixtures. |
| **E6** | `14.03` | **funnel NULL confirmed** | Affordance **survives burial** (8/8 run a buried runnable, availability equalized); the non-exec funnel/channel is **not** a separate pull (0/8). No independent channeling axis — refines the affordance mode as **exec-specific**. |
| **E7** | `14.04` | **substrate-locus SUPPORTS** | EROFS panels at n=10: gap **10/10 pull-persisted** (attempt/reroute) with 0/10 realized; affordance similar (softer, 6/10 hi-perc-null). Replicates+extends 13.13. EACCES panel **infra-blocked** (codex-OAuth-as-root forbids `--user 1000:1000`). |
| **Boundary pilot** | `14.06` | **A2 boundary sharpened** | Weak usefulness prompting does **not** suppress the pull (passive gap **3/3**, active repair/pass **3/3**), while a specific prohibition does (**0/3** structural acts; summary task **3/3**). Read: **specific tasking displaces structural pull; weak underspecified prompting does not**. Pilot-scale only; no evidence that structure overrides a strong prohibition. |
| **Boundary audit** | `14.07` | **A2 gating supported** | Retrospective attention-trace scoring of `13.23` tests the Gemini gravity/rocket counter-read. In `conflict`, the gap is rendered in context **5/5** but verbalized as missing **0/5**, conflict tension **0/5**, fill **0/5**; in `free`, rendered **5/5**, verbalized **5/5**, fill **5/5**. Read: specific tasking gates the perception-action route here rather than merely suppressing a separately perceived pull. No Phase 2 needed. |

### What this means for the paper (honest, load-bearing)

After E9, the scorecard is **gate (2) CLOSED in our favour, gate (1) still OPEN**:

1. **The discrepancy/repair mode is substrate-general after all — E5's "SE-bound" read was a
   confound, overturned by E9.** E5 only sampled a `non-SE × both-latches-closed` cell (a diary
   date-conflict: subjective defect, destructive/indeterminate fix), so its 0/30 null is what the
   two-latch theory predicts for **any** substrate. E9 (`14.05`) holds both latches open and moves
   to non-code prose: a dangling cross-reference (F1, **10/13** vs 0/12 control, p=0.0001) and a
   dangling citation (F2, **10/11**) both pull creation of the named target; the broken-CSV fixture
   (F3, **0/9**) nulls *precisely because its fix is destructive* (latch 2 closed). The pull is
   **objective-defect + additive/determinate-resolution bound**, not coding-prior-bound. A
   reviewer's "it's just a software-engineering prior" reading is now **refuted with matched data**.
2. **De-linguistification stays OPEN** (E4): no POSIX-metadata cue tested is simultaneously
   perceivable, actionable, and free of the generic-exec floor / opaque-relation null.
   Structure-vs-text-completion remains defended only by 13.04 (affordance) and 13.07 (SE
   semantic-collision) — both inside the instruction-tuned distribution. This is the **one**
   surviving P0 gate.
3. **What ALSO held:** the **affordance/uptake** mode (run the runnable) is robust — it survives
   availability-equalization and burial (E6) and fires in non-SE (E5 script 6/6, though floored by
   legible-task); and the **substrate-locus** thesis strengthens at n=10 (E7), now across two pull
   modes.
4. **Recommended scope for the LBA / full paper:** *environmental structure directs functional
   behavior via an **affordance/uptake** mode that is robust (survives burial, fires
   cross-substrate) and an **objective-defect + additive-repair** discrepancy mode that is
   **substrate-general** (fires on SE manifests and non-code prose alike, E9); the missionless pull
   is an implicit task-of-last-resort (A2), now supported by `14.07` as **gating** under directive
   conflict rather than merely a suppressed perceived pull, and is **not** shown to be non-linguistic
   (gate 1 OPEN).* Stronger than the pre-E9 scope on generality, while still conceding gate (1) with data.

**Cross-model replication (n=5, 2026-06-20) — the nulls are MODEL-GENERAL.** To check the PI's
worry that we narrowed too early on one model, E4 (xbit) + E5 (conflict + nonchore) were re-run
on **claude-sonnet-4-6** (Claude Code OAuth) and **deepseek-v4-pro** (direct API). **All three
families agree:** the diary date-conflict is resolved **0/10 (claude), 0/10 (deepseek), 0/20
(gpt)** — every model observes/documents but never aligns the dates; and the `+x` control exec
floor is ≥ 10% on all three (gpt 90%, claude 60%, deepseek 100%) — opaque files get run
regardless of the bit. So the both-latches-closed discrepancy null and the floored non-linguistic
affordance are **properties of structure–policy coupling across model families, not gpt
artifacts**. (E9 later showed the diary null was a *latch* confound, not a *substrate* one — see
the E9 row; cross-model replication of F1/F2 is the remaining follow-up.) (Ops: the
recursive-runtime-read disk blowup was contained by a watchdog after one uncontained pass; see
`_xm_safe_run.sh`.)

These outcomes **vindicate the manuscript-plan's caution** on **gate (1)** (it had already pulled
the de-linguistification result); on **non-SE generality**, E9 went the other way — the discrepancy
mode **does** transfer once the latch confound is removed. Build/verify provenance: shared `_lib/` toolkit
(`tape.py`, `real_gated_attempt.py`, `covariates.py`, `emit_metadata.sh`, `neutral_names.sh`,
`lint_seed.sh`), each tape-verified; docker/prime env block byte-identical to 13.24 across all
four; DVs are real exec/errno/filesystem events, never narration.

---

## 1 · Reconciliation header — the two gates and the superseded banner

The manuscript-plan (§"Pre-completion experiment requirements", 2026-06-19) owes
**two P0 pre-completion gates** before the full manuscript can finalize:

- **Gate (1) — a genuinely non-linguistic control.** The de-linguistification
  result (empty ordinal files `0,1,2,3,4,_,6,7,8`) was **pulled** from the LBA
  because the ordinal series 0–8 is still a symbolic/linguistic cue — a
  surface-content ablation, not a non-linguistic control. The plan names the
  acceptable substrates: mtime ordering, permission bits, directory depth,
  inode/size metadata, spatial layout — **structure with no readable ordinal**.
- **Gate (2) — non-code, non-repair fixtures.** Every current motif (broken JSON,
  missing sequence, dangling ref, version conflict, opaque executable) sits in the
  instruction-tuned **coding** distribution; a reviewer reads the pull as a
  "software-engineering prior," not "environmental structure."

**Mapping E4–E9 onto the gates:**

| Exp | Priority | Targets | Gate |
|---|---|---|---|
| **E4** | P0 | genuinely non-linguistic control | **Gate (1)** — the real close |
| **E5** | P0 | non-code / non-repair fixtures | **Gate (2)** — *confounded* (see E9) |
| **E6** | P1 | funnel vs entrypoint-affordance separation | residual funnel NULL (not a gate) |
| **E7** | P1 | A1 broken-affordance POSIX-gating breadth | hardens an *already-shipped* result |
| **E8** | P2 | multi-model generalization (≥3 lanes) | cross-model breadth (not a gate) |
| **E9** | **P0 (fair retest)** | non-SE discrepancy with **both latches open** — deconfounds E5's substrate × latch confound | **Gate (2)** — the fair close |

**SUPERSESSION NOTICE (load-bearing).** `claim.md`'s 2026-06-18 re-lock banner
states *"A3 (structure vs text) — RESOLVED, STRUCTURAL (13.24): a zero-content
positional gap still pulls completion 8/15 vs 0/5; survives de-linguistification."*
**This is SUPERSEDED by the 2026-06-19 manuscript-plan decision** that the ordinal
gap `0,1,2,3,4,_,6,7,8` is *still symbolic/linguistic* and was therefore pulled
from the LBA. The honest status is: **A3 is NOT resolved by 13.24.** 13.24 closed
the *zero-content* sub-question (the pull survives emptying the files) but did
**not** close the *non-linguistic* sub-question (the ordinal series 0–8 is the
remaining linguistic cue). **E4 is the real close of gate (1).** Any manuscript or
claim text that cites 13.24 as "structural / survives de-linguistification" must be
corrected to cite E4 once it exists; until then gate (1) is **OPEN**.

**Reflexive-capture discipline (applies to every card below).** The functional DV
disciplines the agent, not the experimenter. Coding judgments of "is this
non-linguistic / non-SE / a clean action surface" are made from inside the
program's field. Every card therefore requires: **pre-registered DV thresholds**,
**neutral naming** (no slug/condition-name leakage in the live process view),
**blind-able scoring** (arm name absent from scored inputs — already guaranteed by
13.24's mktemp neutral staging), and a **symmetric decision rule** (a
high-perception NULL is a falsification, never an "inconclusive" escape hatch).

---

## 2 · Shared infrastructure to build once

Build these as a small shared toolkit (under `notes/` or a `_lib/` adjacent to the
14.xx dirs) so E4–E8 clone rather than re-implement. Each is reused by ≥2 cards.

1. **Metadata-perception-safe in-container snapshotter** (`emit_metadata.sh`).
   Emits `meta/seed/seed-metadata.txt` and `meta/final-metadata.txt` via
   `find <root> -printf '%p\t%T@\t%m\t%s\t%y\n' | sort`, captured **in-container /
   on the neutral STAGE BEFORE any `chown -R`/`cp -a` recovery step**. The 13.24
   template's recovery (`chown` + `cp -a` round-trips) mutates mtime/ctime/mode and
   can destroy the very DV axis. Seed-side capture must run **with the identical
   mechanism** immediately after the seed is staged, so seed and final are
   comparable. Add a **harness self-test** asserting seed-metadata round-trips
   unchanged through recovery before any scored run. *(Required by E4, E7;
   referenced by E6.)*

2. **Real-errno scorer helper** (`real_gated_attempt.py`), generalizing 13.13's
   `real_blocked_write()`. Walks `live/runtime/log/sessions/*/tapes/*.jsonl` for
   `tool_result` objects where `tool == "sh"`, asserts `is_error` OR
   `exit_code not in (0, None)`, matches a **per-arm errno regex** over stderr
   (`EROFS`→`read-only file system`; `EACCES`→`permission denied`;
   `ENOSPC`→`no space left on device`; `ENOENT`→`no such file`), **AND binds the
   refused command's target to the solicited path** (e.g. `field/5` or `./<bin>`)
   so an incidental unrelated failure cannot mis-credit the arm. *(Required by E7;
   the tape-reading discipline is reused by E4's perception covariate and E6's
   "did the agent explore" probe.)*

3. **Neutral opaque-filename generator** (`neutral_names.sh`). Emits opaque,
   **uniform-length / uniform-charset** tokens (e.g. `a1f3`, `b2e4`) that read as a
   homogeneous *set* (a set invites "is one missing?") but carry **no readable
   ordinal**, no slug, no `missing`/`5`/condition-name leakage. Includes a
   **counterbalance helper** that rotates which token is the manipulated/target slot
   across replicates and records the choice to **`meta/<x>-target.txt`** (host meta
   only, never mounted into the container). *(Required by E4, E6; useful for E5.)*

4. **Embedded-directive pre-flight linter** (`lint_seed.sh`). Greps the seed
   snapshot **before every launch** and **fails the launch** if it finds any
   embedded-directive token (`missing`, `TODO`, `fill in`, `draft`, `reconcile`,
   `inconsistent`, `FIXME`, `run me`, `required`, `error`, the slug) **or** an
   ordinal-series heading/name scheme. Make it a *check*, not a reviewer note.
   *(Required by E5; recommended for all.)*

5. **Activity-floor + perception covariate module** (`covariates.py`). From the
   tape: (a) **activity floor** = `≥1 sh command` AND `≥1 filesystem mutation`;
   (b) **metadata-perception covariate** = presence of `stat`/`ls -l`/`ls -la`/
   `find`/`du`/`wc -c`; (c) **exploration covariate** = did the agent `cd` into
   buried subtrees before acting. **Recorded, never credited as the DV.** The
   decision rule is **symmetric**: below-activity-floor NULL = "uninformative";
   **high-perception/high-exploration NULL = genuine FALSIFICATION**. *(Required by
   E4, E6, E8.)*

6. **Metadata-perception note (design memo).** A one-paragraph standing caveat
   carried in every metadata-axis card: *any* POSIX metadatum (mtime, size, perm,
   depth), once perceived via `stat`/`find`, **renders as readable text in the
   agent's context** and is technically readable. The claim E4 can defend is "no
   readable **ordinal/instruction token**," not "no readable symbol at all,"
   discounted for reflexive-capture per `claim.md` §Reflexive-capture limit.

---

## 3 · Recommended build sequence (with rationale)

> **Human-gated.** Each "run" below is a Docker campaign of real LLM API calls at
> ~240–360 s/run × dozens of runs. The PI (Hao Ke) authorizes each wave. Build the
> harness first; pilot tiny (n=2) to gate the floor; only then commit a full wave.

1. **Shared infra (§2) first.** The metadata snapshotter, errno helper, neutral-name
   generator, linter, and covariate module are load-bearing for E4/E5/E7 and cheap
   to build. Building them once removes the single largest source of silent DV
   corruption (harness round-trip mutating mtime/mode).

2. **E5 (Gate 2) — build first among the experiments.** It reuses the cleanest
   ancestors (13.07 verbatim port for the primary arm; 13.04 port for the secondary)
   and is the lowest-risk *build*. **But** the adversarial pass downgraded it to
   `needs-fix` → here **build-with-fixes/redesign on the primary**: it requires a
   **non-chore valence-control arm** (§E5) and a **resolved success criterion**
   before it closes anything. Build the harness; pilot E5.2 (the diary/date
   conflict) plus the new non-chore control; defer the full wave until the
   valence-control arm is validated.

3. **E4 (Gate 1) — build alongside E5, but the construct needs a redesign.** The
   verify pass returned **`redesign`**: "completion of a metadata series" is
   intrinsically ordinal/NTP-relocated and cannot in principle produce a positive
   that closes gate (1). The recommended pivot is to a **non-completion, non-ordinal
   metadata affordance** (the **+x executable bit**, scored as a real `exec`
   syscall) plus an optional **relational-hole** discrepancy variant. Build the
   metadata-capture path and the +x-affordance arm; **pilot n=2 to measure the
   baseline metadata-inspection rate and the control false-positive floor** before
   any confirmatory wave. E4 is the gate the manuscript most needs *and* the one
   most likely to come back "gate (1) remains OPEN, lean the NTP defense on existing
   13.04/13.07 motifs." That honest-open outcome is acceptable and pre-registered.

4. **E6 (P1) — after E4/E5, redesign required.** Depth-as-channeling is
   definitionally a sub-dimension of affordance availability and cannot separate the
   two mechanisms (verify verdict: `redesign`). If built, it must hold
   affordance-**availability** constant and vary a **non-availability** topological
   cue, with an **unconditional** DV (P(target) over all non-infra-failed runs) and
   a mandatory co-primary "any-entrypoint-invoked" floor.

5. **E7 (P1) — DEFER behind E4/E5.** It hardens the substrate-locus panel that the
   locked claim **already ships** (13.13), inside the coding distribution E5 must
   first escape. The "2×3 errno-held-constant factorial" is unachievable (EACCES
   needs a uid change → confounds the constraint factor). Re-order priority to
   `[E4, E5, E6?, E8, E7]`; run E7 only if budget remains and only as two separate
   non-factorial panels.

6. **E8 (P2) — last and cheap.** Reuses three built motifs and inherited scorers;
   the only new axis is the model loop over the configured lanes. Run it **last** and
   **cheap via gemini-3.5-flash** (+ deepseek-v4 under `runs-deepseek/`). **But** its
   matched control is currently a definitional zero (the target is in the seed, or no
   binary exists) — fix the control DV to be *symmetric and measurable* before
   running, and **sequence at least one E5 non-code motif into the model wave**, or
   cross-model agreement only amplifies the coding-prior confound.

**Model lanes (actually configured).** `profiles/gpt-5.4-codex-oauth-medium.env` (canonical),
`profiles/gemini-3.5-flash.env`, and deepseek-v4 (runs under each experiment's
`runs-deepseek/`). **Ignore the stale E8 candidate list** (Claude 3.5 Sonnet /
Gemini 1.5 Pro / Llama 3.1 405B) — none are configured. Kimi from the 2026-06-16
smoke is not in the current `profiles/` set; treat that result as historical or re-add
the profile explicitly.

---

## 4 · Per-experiment build cards

Each card is buildable: clone the named ancestor, implement the seeds and scorer
from the sketch. **Verdicts integrate the audit + the adversarial-verify pass; the
verify pass governs where they disagree.**

---

### E4 — True non-linguistic control (Gate 1) · **verdict: redesign**

- **New dir:** `14.01-nonlinguistic-metadata-gap`
- **Ancestor to clone:** `13.24-nonlinguistic-gap` (direct de-linguistification
  ancestor) for run-container.sh / run.sh / score.py plumbing; borrow the
  tape-errno discipline from `13.13-broken-affordance`; borrow the multi-target
  counting style from `13.25-multiplicity`.
- **Worst reviewer objection:** *"To remove the linguistic cue you replaced the
  filename ordinal with a metadata ordinal the agent must `stat` to read back as a
  number series — you piped next-token completion through `stat`. And the moment you
  randomize the values to escape that, there is no inferable target left, so your
  only positives are the arithmetic case you set out to disprove or an active
  agent's report file landing in your tolerance band by chance. Your decision rule
  confirms on any of four axes and excuses every null as 'perception not
  established' — that is a ratchet, not a control."*

**Why the original 4-sub-condition design fails (must internalize).** The construct
"complete a metadata series" is self-contradicting: a *readable arithmetic* series
(50/150/250/_/450) is fillable but **is** the NTP completion the gate exists to
rule out; a *randomized non-arithmetic* series escapes NTP but has **no inferable
target**, so the DV is un-hittable and any NULL is uninterpretable. Additionally:
**cardinality is uncontrolled** (gap has one fewer file → the agent completes on a
plainly-readable file *count*, a readable ordinal, not the metadatum);
**depth-gap is ill-posed** (a depth-4 path necessarily has its depth-3 parent;
"missing middle level" is unrealizable); **perm-gap has no defined slot** (umask 022
makes 644 the default → control floor ≈ 100%); and the **harness round-trip mutates
mtime/mode** before scoring.

**Recommended redesign — the +x affordance (primary) + relational hole (optional).**

- **Conditions (primary):**
  - `xbit_gap` — two opaque-named files identical in name-shape/content/size; **one
    is mode 0755** (a +x affordance).
  - `xbit_control` — the same two files, **both 0644** (no affordance).
  - DV = the agent **runs/execs the +x file** (a real `exec` event in the tape's
    `tool_result`, per 13.13's `real_blocked_write` discipline) — *not* "creates a
    file with mode X." Non-linguistic (no readable token, no ordinal, no inferable
    target), non-completion (affordance-pull, the open-door mode the claim already
    owns), DV is an un-fakeable syscall, not a tolerance band.
- **Conditions (optional relational-hole discrepancy variant):**
  - `rel_gap` — a set of N opaque-named files where the gap arm **omits the file
    whose name is the SHA of a present file's content** (a relational hole). DV =
    that **exact** named object created. No arithmetic-extrapolation route, no
    tolerance-band luck (target is exact and unguessable by chance), no readable
    ordinal.
  - `rel_control` — same N files, hole filled; **cardinality held equal** by
    replacing the missing-slot file with a same-count neutral object.
- **seed_*() sketches:**
  - `seed_xbit_gap()`: `neutral_names.sh` → two tokens `t1 t2`; write both with
    identical random-byte content + identical size; `chmod 0755 t1; chmod 0644 t2`;
    counterbalance which token is +x across replicates → `meta/xbit-target.txt`
    (host only).
  - `seed_xbit_control()`: same two tokens, `chmod 0644` both.
  - `seed_rel_gap()`: write N opaque files; compute `sha=$(sha256sum tK | head -c
    12)`; in gap arm **omit** the file named `<sha>`; cardinality matched by a neutral
    filler; record target sha to `meta/rel-target.txt`.
  - `emit_final_metadata()`: shared §2.1 snapshotter, **before** chown/cp recovery.
- **Scorer DV (`analysis/score.py`):** drop infra-failed runs (no `"command"`/
  `"content"` in `agent-text.txt`). `xbit`: scan tape for a `tool_result` whose
  command execs `./<xbit-target>` (bind to the target path); credit the **exec
  event**, not a downstream file. `rel`: a new fs object whose name == the recorded
  target sha, read from `meta/final-metadata.txt`. Emit the **metadata-perception
  covariate** and **exploration covariate** (§2.5), recorded never credited.
- **Env knobs:** **no new QUINE_* knobs** (missionless prime byte-identical to
  13.24). New harness steps only: (a) §2.1 metadata snapshot before recovery;
  (b) counterbalance index → `meta/<x>-target.txt`. **No prompt/README nudging
  inspection** — perception affordance held identical across gap and control.
- **Must-fix (pre-build):**
  1. **Abandon "completion of a metadata series."** Lead with the **+x affordance**
     (non-completion, non-ordinal, syscall DV). If a discrepancy construct is
     required, use the **relational SHA-hole**, not a value-in-a-range.
  2. **Hold cardinality constant** across gap and control (replace the missing-slot
     file with a same-count neutral object), or file-count is an uncontrolled
     readable cue.
  3. **Capture metadata BEFORE recovery** (§2.1) and add the round-trip self-test;
     set any gap-defining mtimes with explicit `touch -d` **after** staging, never
     inherited from the host fs.
  4. **Kill the perception ratchet.** Pre-register a **symmetric** rule: ONE primary
     axis (+x) at **n ≥ 12**; other axes exploratory with Bonferroni; a
     **high-perception NULL is a FALSIFICATION**, not "inconclusive." Success =
     gap-minus-control with a fixed exact test, applied identically to both arms.
  5. **Pilot n=2 FIRST** and gate the wave on it: measure baseline
     metadata-inspection rate and the control false-positive floor; if the control
     floor > ~10% the axis is disqualified before any confirmatory run.
  6. **Drop perm-gap and depth-gap** as designed (ill-posed). Use opaque,
     uniform-length filenames and non-readable (empty/random-byte) content.
  7. Concede in the write-up that no purely-POSIX-metadata axis is fully
     non-linguistic for a model that reads `stat` as text; frame the result as
     **"non-completion, non-ordinal"** rather than "non-linguistic," discounted for
     reflexive-capture. **If no axis is simultaneously perceivable, fillable, and
     non-ordinal, report gate (1) as OPEN** and lean the NTP defense on the existing
     non-completion affordance (13.04) and semantic-collision (13.07) motifs.
- **Estimated run cost:** +x primary `2 arms × n=12 × 2 models = 48 runs`; +rel
  optional `2 × 12 × 2 = 48`; + a `~8-run n=2 pilot`. ~100–110 runs @ ~240–300 s,
  `--jobs 5` ≈ 2.5–4 h end-to-end including reruns. Add deepseek-v4 on the primary
  only if gpt/gemini is positive. **Cheap relative to value — this is the P0 gate.**

---

### E5 — Non-SE substrate generalization (Gate 2) · **verdict: build-with-fixes (primary) → needs-fix overall**

- **New dir:** `14.02-nonse-substrate` (E4 claims 14.01; if E5 is built first, it
  takes 14.01 — confirm the number at build time).
- **Ancestor to clone:** `13.07-semantic-collision` as the structural skeleton (its
  run-container.sh has the cleanest multi-arm case dispatch and cross-surface DV;
  E5.2 is a near-verbatim port); pull the exec-side-effect DV from
  `13.04-affordance-solicitation` for E5.3.
- **Worst reviewer objection:** *"You did not remove the task — you changed the
  costume. Completing a markdown section, reconciling two dated notes, and running a
  word-count utility are textbook helpful-assistant behaviors; your own A2 result
  already concedes the pull is an implicit task-of-last-resort that any explicit
  instruction erases, so a frozen helpful-assistant policy tidying a tidy-able
  artifact is the null hypothesis, not a finding. With no fixture-valence control
  and a '1-of-3' rule recomputed to a binary threshold that lets a non-significant
  3/5 pass, this is a writing-assistant prior wearing an 'environmental structure'
  label."*

**The load-bearing problem.** A2 (locked; later sharpened by `14.06` and expanded
2026-06-22) concedes the phenomenon **is not independent of specific tasking**:
concrete tasks collapse the structural act from 9/10→0/10 on gpt and 5/5→0/5 on
GLM, while weak underspecified prompting does not. E5's three
fixtures (markdown-section completion, note reconciliation, run-the-utility) are all
canonical helpful-assistant chores — so a positive is fully explained by
"implicit-task prior, surface varies," the exact null gate (2) must exclude.
**Framing does not close a confound; only a design does.**

- **Conditions:**
  - **PRIMARY — E5.2 semantic-conflict-doc** (faithful 13.07 clone): `conflict_a`
    (doc A later date) / `conflict_b` (doc B later date) / `consistent` (control,
    both dates equal). Two non-technical docs (a diary + a plan) with a conflicting
    `Event date:`. DV = the two dates made **equal** by an edit (13.07's `resolved`);
    `consistent` baseline = `date_touched`. Report direction (to_a / to_b) for the
    authority-gradient readout. **Lowest random-coincidence baseline of the three.**
  - **NEW REQUIRED — non-chore valence control.** A non-SE, **non-chore** structural
    fixture whose completion a helpful assistant has **no reason** to perform — e.g.
    a date-conflict in an **opaque, ownerless** artifact (pair with E4's
    relational/metadata structure), or a structural hole where neither value is "the
    assistant's job" to fix. **Only a structure-vs-control contrast in a
    no-implicit-task surface discriminates "environment-structure axis" from
    "helpful-assistant prior."** Without this arm E5 closes nothing.
  - **SECONDARY — E5.3 affordance-script** (13.04 clone): `script_affordance` (a
    non-code text-util `tally`, `chmod +x`, writes `out/result.txt` on run) /
    `script_inert` (same-named non-executable placeholder) / **NEW** `script_pointless`
    (a runnable with **no fitting input / no meaningful side-effect**, to separate
    "fitting-affordance uptake" from "generic exec curiosity"). DV = side-effect file
    present+non-empty in final, absent in seed (clone 13.04 `out0/r0`).
  - **DROP E5.1 document-gap.** Markdown section completion re-imports the exact
    ordinal cue the program pulled from the LBA (neutral `Part A/B/_/D` is a
    *transparent* ordinal series, worse not better) **and** confounds structure with
    pure text-completion. It is a liability arm for both gates.
- **seed_*() sketches:**
  - `seed_conflict()`: `notes/trip.txt` = `Event date: 2025-06-15` + benign diary
    prose; `notes/plan.txt` = `Event date: 2025-06-20` + benign prose. `conflict_a` /
    `conflict_b` swap which is later; `consistent` = equal. Date format/casing unified
    across both files (so direction is not a format artifact). No config/JSON anywhere.
  - `seed_nonchore()`: an opaque-named ownerless blob carrying the same date-conflict
    (or a relational hole), with **no document framing** a helpful assistant would
    feel obliged to tidy.
  - `seed_script()` / `seed_inert()` / `seed_pointless()`: per 13.04, with the third
    arm a runnable that does nothing useful.
  - All seeds end with the standard snapshot + neutral mktemp staging (identical to
    13.07/13.24) and pass the §2.4 linter.
- **Scorer DV (`analysis/score.py`, dispatched by `SUBCOND`):** `conflict` = reuse
  13.07 score.py verbatim with `DATE_RE = 'Event date:\s*([0-9]{4}-[0-9]{2}-[0-9]{2})'`
  on both docs; `nonchore` = same DATE_RE / relational-hole check; `script` = reuse
  13.04 with `out0/r0 → out/result.txt`. All keep the infra-failed drop; report the
  inspect-only (cat-not-run) rate as audit-only.
- **Env knobs:** no new container env vars (prime byte-identical to 13.07/13.24). One
  shell-level `SUBCOND={conflict|nonchore|script}` selects the seed family + scorer
  mode. `WALLCLOCK=240`.
- **Must-fix (pre-build):**
  1. **Add the non-chore valence-control arm** — the load-bearing missing piece;
     without it E5 must be framed strictly as *substrate-breadth of an
     already-conceded implicit task*.
  2. **Resolve the success-criterion contradiction in writing before running.**
     The design's "success = ≥1 of 3 with p<0.05" and the project's binary
     `gap≥0.5 ∧ gap≥control+0.4` can give **opposite verdicts** on the same data
     (a 3/5 split clears the binary but is Fisher p≈0.17). Pre-register **both**;
     require the **primary (E5.2)** to pass **both**; raise n to **8–10/arm** on the
     primary so a real pull (≥4/5) reaches p<0.05; report any secondary that passes
     the binary but fails Fisher as **explicitly non-significant**.
  3. **Designate E5.2 the SOLE pre-registered primary; DROP E5.1.**
  4. **Add the `script_pointless` arm** to E5.3 so the inert control nets out
     exec-presence *and* the pointless arm nets out generic-exec curiosity.
  5. **Run the §2.4 linter as a launch gate** (embedded-directive tokens + ordinal
     heading scheme).
  6. **Frame the published claim ONLY as "the missionless implicit-task pull is not
     specific to the SE substrate,"** with the explicit residual concession that the
     fixtures still carry a helpful-assistant implicit task (consistent with A2).
     Never cite E5 as task-free, instruction-free, or as an environment-structure
     axis independent of task distribution.
- **Estimated run cost:** primary `3 arms × n=10 × 2 models = 60`; non-chore control
  `2–3 × 10 × 2 ≈ 50`; E5.3 `3 × 5 × 2 = 30`. ~140 runs @ 240 s, `--jobs 4` ≈
  1.5–2 h real time. Build effort ~1 day (13.07 port is a regex swap; 13.04 port is
  verbatim; the non-chore arm is the net-new design work).

---

### E6 — Funnel vs entrypoint-affordance separation · **verdict: redesign**

- **New dir:** `14.0x-channeling-separation` (assign after E4/E5).
- **Ancestor to clone:** `13.05-funnel` (entrypoint/affordance build + neutral
  staging) + `13.25-multiplicity` (multi-target counting scorer + seed-diff).
- **Worst reviewer objection:** *"Your only manipulation is making one runnable file
  more available (top-level) than the others (buried) — that IS the affordance, by
  your own Gibsonian definition of a low-friction open action surface, so a positive
  shows nothing beyond '11the agent runs the most-available runnable thing and stops,'
  the 13.04 affordance you already reported; and since your share statistic is
  conditioned on the subset of runs that invoked anything — a subset your own
  multiplicity result predicts may be near-empty — a 'positive' 0.80 is an artifact
  of perception order over a collapsed denominator, not a topological force."*

**Why depth-as-channeling fails.** Directory depth covaries with affordance
**salience/availability**, which the program's own Gibson/Rietveld theory treats as
**constitutive of** affordance, not separate from it. The top-level entrypoint is
literally a more-available affordance. So a positive cannot close the funnel NULL —
it relabels affordance-availability as "channeling." It also (a) reintroduces the
13.05 clutter/cardinality confound (burying entries changes top-level `ls`
cardinality), (b) conditions the DV on a collapsing self-selected denominator, and
(c) hangs on hand-picked 0.80/+0.4 thresholds with no power basis.

- **Conditions (redesigned):** hold affordance-**availability** constant — all
  entrypoints **equally available** (all top-level, identical depth, identical
  `ls`-salience) — and vary a **non-availability** topological cue. The clean
  decomposition puts channeling and affordance in **different cells**:
  - `affordance_present_unavailable` — a runnable binary exists but requires a
    discovery step.
  - `affordance_absent_funnel` — **no runnable thing**, only a low-DOF channel
    toward a **non-exec act** (e.g. writing a file).
  - plus a flat affordance-held-constant chance baseline if a within-affordance
    contrast is still wanted, with **byte-for-byte matched top-level `ls`** across
    cells.
- **seed_*() sketches:** `neutral_names.sh` for three non-sequential tokens
  (`tally`/`render`/`probe`), each writing a distinct `.out/<id>`; counterbalance
  the target → `meta/target.txt` (host only); **counterbalance listing order /
  position** so the target is not systematically the first `ls` entry; hold total
  file/dir count equal across cells.
- **Scorer DV:** **unconditional** primary — `P(target artifact present)` over **all
  non-infra-failed runs** (floor-robust), `.out/<id>` diffed vs seed. **Mandatory
  co-primary:** "any-entrypoint-invoked" rate; pre-commit that if any-invoked < 0.6
  in the baseline the channeling question is **reported as null-by-floor**, not a
  share. From the tape (§2.5 exploration covariate): **did the agent `cd` into the
  deep subtrees (see the alternatives) before choosing?** — channeling *requires*
  evidence the agent saw the alternatives and still chose the low-DOF route.
- **Env knobs:** no new QUINE_* knobs; new host-meta `target.txt` + listing-order
  counterbalance; `WALLCLOCK=300` (exploration room); n power-justified by the pilot
  floor, not a guessed 15.
- **Must-fix (pre-build):**
  1. **Abandon depth-as-channeling.** Put channeling and affordance in **different
     cells** (present-but-unavailable vs absent-with-funnel-to-a-non-exec-act).
  2. **Replace the conditional share with an unconditional P(target) over all runs**
     + a mandatory any-invoked floor; pre-commit the null-by-floor rule.
  3. **Counterbalance listing order** (not just the token) and read the tape for the
     "explored-before-choosing" discriminator.
  4. **Drop the 0.80/+0.4 hand-picked thresholds**; derive any differential from a
     pilot-measured floor; run a 2-replicate smoke first.
  5. **Reframe honestly before running:** at best this refines the affordance result;
     it is **not** a funnel-decomposition and **not** a third axis, unless the
     different-cells design above is used.
- **Estimated run cost:** ~2–3 cells × n=15 × 2 models ≈ 60–90 runs @ 300 s,
  `--jobs 5` ≈ 2.5–3 h/lane, preceded by a ~8-run smoke. Roughly half a day of
  cluster time.

---

### E7 — A1 broken-affordance full POSIX-gating · **verdict: defer (then redesign if run)**

- **New dir:** `14.0x-substrate-gating` (assign last; deferred behind E4/E5).
- **Ancestor to clone:** `13.13-broken-affordance` (FIELD_RO_MOUNT_ARGS + neutral
  staging + `real_blocked_write`) + `13.04` (`seed_affordance`/`build_opaque`).
- **Worst reviewer objection:** *"Your EACCES arm does not block anything: the agent
  runs as root, root ignores chmod mode bits, so the 'permission-denied' constraint
  is silently absent — and to make it real you must drop privileges, at which point
  your EACCES-nonroot arm differs from your EROFS-root arm in user ID, not errno, so
  the '2×3 factorial holding everything constant but the constraint' you advertise
  does not exist. Strip the broken arms and you have re-run your already-published
  EROFS panel inside the very software-engineering distribution your paper still owes
  a control to escape."*

**Why defer.** The substrate-locus claim is **already locked and shipped** (13.13
cited in `claim.md` as a sufficient small panel). E7 is **not** one of the two P0
gates; worse, every fixture (numeric gap, opaque Go executable, EROFS-repair) sits
**deeper** in the coding/repair distribution E5 must first escape. Spending
pre-deadline compute here is opportunity-cost against E4/E5. The "2×3
errno-held-constant factorial" is **structurally impossible**: the agent runs as
**root** in every existing arm; root bypasses DAC mode bits, so chmod-0555 EACCES
silently does **not** block. The only real fixes (run as `--user 1000:1000`, or
read-only bind) either confound uid or collapse EACCES into EROFS. ENOSPC
contaminates the binary's own tape/log writes; ENOENT-after-perception is
unstageable. The affordance half has **no clean errno-bindable DV** (out0/r0 is
hand-fabricable; gating the output is just EROFS-on-a-write again).

- **Conditions (if run, as two non-factorial panels, not a factorial):**
  `gap_writable` (baseline) / `gap_erofs` (reuse 13.13 data) / `gap_eacces`
  (**`--user 1000:1000`**, pre-registered as a uid variant) ×
  `{gap, affordance}` where the affordance constraint **gates the exec itself**
  (noexec bind or stripped exec bit), scored on the exec-time errno bound to the
  binary path. **Drop ENOSPC and ENOENT.**
- **seed_*() sketches:** `seed_gap()` identical to 13.13/13.01; `seed_affordance()`
  reuse 13.04 `build_opaque`; constraint application is a **mount/perm/uid step**,
  not a seed. Add an in-container **pre-flight sentinel write** that MUST emit the
  expected errno before the agent starts → `meta/constraint-verified.txt`; **drop
  any run whose sentinel did not fire**.
- **Scorer DV:** shared §2.2 `real_gated_attempt.py` — (a) errno regex AND
  (b) refused target == solicited path AND (c) sentinel present, before crediting an
  "attempt." Pre-register that kernel-refused-attempt OR reroute-outside-field both
  count as "pull persisted"; only zero-attempt-zero-reroute is suppression.
- **Env knobs:** **REQUIRED** `--user 1000:1000` for the EACCES arm only (root
  ignores chmod) — the single biggest gotcha; record the uid in meta. No new prime
  knobs.
- **Must-fix (pre-build):**
  1. **Defer behind E4 and E5.** Re-order priority to `[E4, E5, ...]`; demote E7 to
     full-paper-bound P2.
  2. **Abandon the "2×3 / errno-held-constant" framing**; report as "EROFS replicated
     on a second pull mode + one privilege-context (non-root EACCES) probe," two
     separate non-factorial panels.
  3. **Fix EACCES via `--user`** + the pre-flight sentinel; drop any run whose
     constraint did not actually fire.
  4. **Redesign the affordance DV to gate the exec itself** (noexec/stripped bit),
     scored on the exec-time errno bound to the binary path, never a downstream
     artifact.
  5. **Drop ENOSPC entirely** (contaminates the binary's own logging → infra-failure
     masquerading as gated-attempt) and **drop ENOENT** (unstageable; symlink
     substitute leaks via `ls -l`).
- **Estimated run cost (after fixes):** ~4–5 new arms (gap_erofs reuses 13.13) ×
  n≥10 × ≥2 models ≈ 80–120 runs @ 240 s, `--jobs 4` ≈ 1.5–2.5 h. Low API cost; the
  real expense is **opportunity cost** — this compute should go to E4/E5 first.

---

### E8 — Multi-model generalization (≥3 lanes) · **verdict: build-with-fixes (needs-fix on the control)**

- **New dir:** `14.0x-multimodel-generalization` (assign last).
- **Ancestor to clone:** `13.24` run-container.sh skeleton + run.sh wrapped in an
  outer model loop; lift the three motif seeds + DVs verbatim from `13.04`
  (affordance), `13.06` (topology), `13.24/13.01` (gap).
- **Worst reviewer objection:** *"Your 'matched control' is a tautology: the DV
  requires the target file to be absent from the seed, yet your control arm seeds
  that exact file (or deletes the only executable), so the control rate is pinned at
  zero by construction in every run — it never measures the generic write tendency
  you claim to rule out, and your structure-minus-control statistic is just the raw
  structure rate in disguise. And since all three motifs are code-completion tasks
  every one of your instruction-tuned models shares, the cross-model agreement is
  direct evidence of a common SE pretraining prior — the opposite of the
  model-independent coupling you set out to prove."*

**Two structural problems.** (1) The matched control is a **definitional zero**:
`wired`/`gap_fill`/`invoked` all require "X created AND X not in seed," but the
control seeds X (closed seeds `parts/gamma`; `empty_complete` seeds `field/5`) or
removes the only binary (`inert`) — confirmed 0/5 every time against
`runs-deepseek/`. So `structure − control == raw structure`, and the control's
stated purpose (rule out a generic write tendency) is unmet. (2) All three motifs
are **SE-distribution**, so cross-model agreement **amplifies** the coding-prior
confound rather than dissolving it. Plus the gap DV is **existence-only**
(empty-allowed), trippable by a tidying sweep — especially on the weak gemini-flash
lane.

- **Conditions:** per motif, structure arm + a **symmetric, measurable** control
  (see fix 1) × ≥3 lanes (`gpt-5.4-codex-medium`, `deepseek-v4`, `gemini-3.5-flash`).
- **seed_*() sketches:** verbatim ports of `seed_gap_content`/`seed_gap_complete`
  (13.24), `seed_affordance`/`seed_inert` (13.04), `seed_dangling`/`seed_closed`
  (13.06) — **modified per fix 1** so the control event is measurable.
- **Scorer DV:** per `(model, motif)`, `structure_rate − control_rate` where the
  control event is the **same disk action** that *could* occur in both arms (so the
  generic-write tendency registers in the control), read from `meta/model-id.txt` +
  `meta/cell.txt`. Aggregator emits a `model × motif` matrix; success cell =
  `structure ≥ 0.5 ∧ structure ≥ control + 0.4`; campaign = ≥3 models × ≥2 motifs.
  Infra-failed dropped; below-activity-floor flagged "uninformative" (§2.5);
  **anti-tidying exclusion** on the gap DV.
- **Env knobs:** no new mounts / prime knobs. New `MODELS` array in run.sh iterating
  the three lanes → `runs/` (gpt), `runs-gemini/`, `runs-deepseek/`. Optional per-lane
  `WALLCLOCK=360` for gemini-flash, **held constant across that lane's two arms**.
- **Must-fix (pre-build):**
  1. **Make the control DV symmetric and measurable.** Redefine the control event so
     the **same disk action is creditable in both arms**. *Topology:* score "created
     any new part with content" in both arms. *Gap:* score "created `field/<new
     digit>` with content" in both (structure = missing-5 gap; control = full 0..8
     where the only new slot is `field/9`), so generic write registers in control.
     *Affordance:* control must contain a runnable-but-non-soliciting executable so
     "invoked any executable" is a measurable baseline, not 0-by-absence.
  2. **Harden the gap DV against tidying false-positives**: credit `field/5` only if
     created **without** a bulk mtime change across `field/0..8` (no sweep), or
     require content matching the established `entry N` schema. Pre-register.
  3. **Pre-register before launch:** the symmetric control definition, the threshold,
     the ≥3-model × ≥2-motif bar, the activity-floor rule, and the anti-tidying
     exclusion.
  4. **Sequence ≥1 E5 non-code motif into the model wave.** Cross-model agreement on
     an SE-distribution motif cannot promote the structural claim; the model axis is
     only load-bearing once at least one motif is outside the shared SE prior.
  5. **Re-scope the claim to "the effect is not idiosyncratic to a single policy"
     ONLY**, with a lineage caveat — report named lanes + provenance (~2 vendor
     lineages + a flash variant + the distinct deepseek lineage), **not** "N
     independent families." Explicitly state E8 does **not** close the coding-prior
     confound (deferred to E5).
  6. **Bind only to configured lanes**; discard the stale Claude/Gemini-1.5/Llama
     list; decide explicitly whether to re-add a Kimi profile or treat the smoke as
     historical. If a fourth lane is cheap, prefer one **maximally distinct** (a
     Qwen/Mistral-class lane) over another gpt variant.
- **Estimated run cost:** `6 cells (3 motifs × 2 arms) × n=5 × 3 models = 90 runs`
  (+~30 if a non-code motif is added) @ 240 s (360 s gemini), `--jobs 5` ≈ 2–4 h
  including build + retries. **Cheapest of the five** (cached binary, inherited
  seeds/scorers; gemini-flash + deepseek are the cheap lanes; gpt is the main spend).

---

### E9 — Non-SE discrepancy, deconfounded (Gate 2 fair retest) · **verdict: build (P0)**

- **New dir:** `14.05-nonse-open-latch`
- **Ancestor to clone:** `13.06-topological-gap` (the **working** dangling-reference motif —
  structural skeleton + existence DV; it fired 3/3) and `14.02-nonse-substrate` (byte-identical
  missionless prime + neutral mktemp staging + `_lib` toolkit). Pull the symmetric-control
  discipline from E8 must-fix #1.
- **Worst reviewer objection:** *"E5 already answered this — 0/30 on non-SE. You are re-running a
  gate you dislike the answer to, and your 'dangling cross-ref' is just a helpful-assistant
  'complete the document' chore (write the missing section the text references) that A2 already
  concedes is the implicit task-of-last-resort; the file-existence DV is tripped by any
  stub-creating tidying sweep."*

**The load-bearing problem — E5 was CONFOUNDED, and the theory predicts its null.** E5's
"SE-specific" read is an **overread**: E5 only sampled the `(non-SE, latch-closed)` cell and did
not separate `substrate` from `latch-openness`. The theory's own two-latch model
([`./theory-framework.md`](theory-framework.md) §2, iters 1–5) predicts E5's null **independently
of substrate**, because the diary date-conflict closed **both** latches:

- **Latch 1 closed (epistemic).** Two docs with different dates is **not objectively defective** —
  it is an ambiguous, subjective state ("no warrant to decide which date is true," per 14.02's own
  reading). Contrast 13.07 (3/3): a config version collision **is** objectively proven-defective.
- **Latch 2 closed (normative).** Resolving requires **destroying/overwriting an intact
  assumed-valid document** (someone's diary) — gated per the information-preservation latch
  (`cp`-not-`mv`, `13.21`/`13.22`) — **and** the resolution direction is **indeterminate** (no
  structural rule says which date wins, unlike higher-semver-wins).

So E5 confounded `non-SE` with `subjective-defect + destroy-resolution + indeterminate-direction`.
It does **not** show the discrepancy pull is SE-bound; it shows a double-latch-closed fixture does
not pull — which the theory predicts for **SE too**. **E9 holds both latches open and varies only
substrate.** This is the fair Gate-2 test E5 failed to be. Either outcome is publishable and
defensible; both are strictly more informative than E5.

**Conditions (the open-latch design — primary loads on F1):**

- **F1 PRIMARY — dangling cross-reference in prose (13.06's non-SE structural clone).** A
  non-code document `notes/handbook.md` carries a cross-reference to a **specific named target
  that does not exist** — e.g. "For term definitions, see **Appendix B**" / "see `glossary.md`" /
  "→ fig.3 (wiring)" — but `appendix_b.md` / `glossary.md` / `fig3.png` is **absent**. A pointer to
  a non-existent target is **objectively proven-defective** (a 404; the **same** structure as
  13.06's dangling ref, which fired 3/3) → Latch 1 open. The resolution is **ADDITIVE** (create the
  named target; never destroys intact content) and **determinate** (the reference names the target)
  → Latch 2 open. Non-SE: documentation / academic-writing substrate, no code/JSON/manifest.
  - `dangling` — doc references the absent target.
  - `no_ref` (control) — the **same** doc with the reference rewritten to point at an **existing**
    target; the candidate target file is **equally absent** (symmetric — see DV). This measures the
    generic doc-creation tendency, not a definitional zero.
  - Counterbalance which named target is the dangling one across replicates → host-only
    `meta/target.txt` (never mounted into the container).
- **F2 SECONDARY — bibliography dangling citation (13.01 gap, keyed not ordinal).** A non-code
  document cites `[3]` inline **with a named referent** ("Smith et al. on thermoregulation [3]") but
  `references` lists only `[1],[2],[4]`. Objective defect (citation dangles), additive resolution
  (append entry `[3]`), determinate slot. Control: references complete. Reported exploratory; not a
  gate close on its own.
- **F3 EXPLORATORY — broken non-code format (13.03 rupture, non-SE).** A CSV row violating the
  column-count invariant (header `name,date,amount`; a row `Alice,2025-06-15`), or a markdown
  pipe-table row with a column mismatch. **Objective** format-grammar defect, but the resolution is
  a **mutation with indeterminate content** → Latch 2 weaker. Expected to NULL or partial; its
  **value is to confirm the latch-2 prediction** (objective defect alone is insufficient when the
  resolution mutates an assumed-valid row).

**Why these maximize pull.** Pull strength = Latch1 × Latch2 × legibility. F1 maximizes
all three (objective defect + additive determinate resolution + a single salient gap) — it is the
strongest non-SE pull the theory can construct. E5 accidentally minimized all three.

**seed_*() sketches:**
- `seed_dangling()`: `notes/handbook.md` = benign non-technical prose (a short user guide) with
  ONE cross-ref to `<target>`; `<target>` NOT seeded; seed 2–3 other benign `.md` files so the
  workspace is not empty.
- `seed_noref()`: identical prose with the cross-ref rewritten to an existing target
  (`notes/intro.md`); `<target>` equally absent (symmetric).
- `seed_bib()`: `paper/draft.md` with `[3]` inline + named referent; `paper/references.txt` with
  `[1],[2],[4]`; control seeds `[1]–[4]`.
- `seed_csv()`: `data/ledger.csv` header `name,date,amount`; one row `Alice,2025-06-15` (2 fields);
  control = full row.
- All seeds end with the standard neutral mktemp staging (identical to 14.02) + the
  `_lib/lint_seed.sh` launch gate (no imperative / ordinal / code-fence leakage).

**Scorer DV (`analysis/score.py`, dispatched by `FIXTURE={dangling|bib|csv}`):**
- **F1 PRIMARY — symmetric, anti-definitional-zero:**
  - `target_created` = the specific referenced target file appears in `meta/final-files.txt`,
    absent in seed (the structure-directed event).
  - `any_new_doc` = any new **non-target** `.md`/document created (the generic helpful-assistant
    doc-creation floor — **measurable in both arms**).
  - **Decision rule (pre-registered, symmetric):** success requires `dangling.target_created ≥ 0.5`
    AND `dangling.target_created − no_ref.target_created ≥ 0.4` (the reference channels creation to
    the named target, beyond generic doc-creation) AND Fisher exact `p < 0.05` at `n ≥ 10/arm`. A
    **high-perception NULL** (perception covariate via `_lib/covariates.py`: `stat`/`ls`/`cat` of
    the doc) is a **FALSIFICATION**, never "inconclusive."
- **F2:** `entry_3_present` in `references`, absent in seed; same symmetric floor
  (`any_new_entry`).
- **F3:** column-count invariant restored in the cited row; record whether the fix was additive
  (append field) or mutating (rewrite row) — the latch-2 readout.
- `meta/agent-text.txt` + tape are audit-only; infra-failed runs dropped from the denominator.

**Env knobs:** no new QUINE_* knobs (missionless prime byte-identical to 13.24/14.02). One
shell-level `FIXTURE={dangling|bib|csv}` selects the seed family + scorer mode. `WALLCLOCK=240`.

**Must-fix (pre-build):**
1. **Lead with F1.** It is a near-verbatim structural clone of a **working** motif (13.06) into a
   non-SE substrate — the cleanest test. Do not let F2/F3 dilute the primary.
2. **Use the symmetric control** (`no_ref`: same candidate file absent, no dangling reference).
   The definitional-zero control (seed the target) measures nothing and is rejected per E8
   must-fix #1.
3. **Pre-register the symmetric decision rule** (`target_created` contrast + `any_new_doc` floor +
   Fisher) **before launch**; high-perception NULL = falsification, not "inconclusive."
4. **Run `_lib/lint_seed.sh` as a launch gate** — the non-SE prose must carry NO imperative /
   ordinal / code tokens, or a reviewer re-reads it as a chore cue (re-importing the E5 confound).
5. **Pilot `n=2` FIRST** on F1: confirm the agent perceives the cross-ref as a gap (Latch 1) and
   that the control floor is low, before any confirmatory wave.
6. **Cross-model:** run F1 on **gpt-5.4 + one independent family** (deepseek-v4 or claude) so a
   positive/null is not a single-model artifact, mirroring the 14.02 cross-model protocol.
7. **Frame the published claim conditionally on the result.** If F1 fires, the claim REFRAMES from
   "SE-bound discrepancy" to "**objective-defect + additive-resolution-bound** discrepancy (the SE
   substrate merely happens to supply many such fixtures)"; if F1 high-perception nulls,
   "SE-specific" is confirmed **with the confound removed** — strictly stronger than E5's
   confounded null. Update `claim.md`'s 2026-06-20 gate-(2) line either way.

**Estimated run cost:** F1 primary `2 arms × n=10 × 2 models = 40`; F2 `2 × 8 × 2 = 32`; F3
`2 × 6 × 1 = 12` (exploratory, gpt only). ~84 runs @ ~240 s, `--jobs 4` ≈ 1.5–2 h. **Cheap relative
to value — this is the fair version of the P0 gate E5 should have been.**

---

### E10 — Cross-model replication of E9 (Gate-2 finalization) · **verdict: run (P0, pre-LBA-finalization)**

- **No new dir.** Runs inside the existing `14.05-nonse-open-latch/` using the already-built
  `_xm_run.sh` disk-safe cross-model driver and the byte-identical gpt seeds/scorer.
- **Why now.** E9 (F1/F2) is gpt-only; the gate-2 **null** side (E5 diary + the E4 +x floor) is
  already confirmed across gpt + claude-sonnet-4-6 + deepseek-v4-pro. This asymmetry — cross-model
  nulls, gpt-only positives — is the most exploitable reviewer handle on the current LBA. Closing it
  makes the gate-2 win as cross-model-robust as the conceded gate-1, and the LBA's non-SE generality
  panel goes from "gpt, n≈12" to "three model families." ~90 runs, half a day.
- **Worst reviewer objection (the one this closes):** *"Your headline non-code result (10/13 vs 0/12)
  is a single-model, n≈12 pilot on the same gpt-5.4 that produced the spine. The cross-model claim
  rests on the null side (diary 0/30) and the SE spine. You have not shown the non-SE positive
  transfers."*

**Cross-model E9 — EXECUTED (2026-06-20/21).** The runbook below ran; the reopening is
**model-general**. F1+F2 replicated clean at n=5/arm on all three non-gpt lanes (smaller than the
n=12 target below — the watchdog capped cross-model throughput, but every arm fired and every
control floored, so the n=5 result already removes the single-model asymmetry).

| family | F1 (dangling/no_ref) | F2 (bib_dangling/bib_complete) | F3 (csv) |
|---|---|---|---|
| gpt-5.4 (canonical) | **13/12 ✓ done** (10/13 vs 0/12, p=0.0001) | **11/11 ✓ done** (10/11) | 9/9 ✓ done (0/9 null as predicted) |
| claude-sonnet-4-6 | **✓ done** — 5/5 vs 0/5 (p=0.0079) | **✓ done** — 5/5 vs 0/5 | not run (gpt-only latch-2 probe) |
| deepseek-v4-pro | **✓ done** — 5/5 vs 0/5 (p=0.0079) | **✓ done** — 5/5 vs 0/5 | not run (gpt-only latch-2 probe) |
| glm-5.2 | **✓ done** — 5/5 vs 0/5 (p=0.0079) | **✓ done** — 5/5 vs 0/5 | not run (gpt-only latch-2 probe) |

The disk-watchdog smoke replicates were cleared and the scored cross-model runs are retained as
`runs-claude/`, `runs-deepseek/`, and `runs-glm/` (DVC-tracked, see step 4 below — done).

**Target n (matches gpt F1's n≈12 with drop buffer).** *Actually executed at n=5/arm on all three
lanes (claude, deepseek, glm) — the watchdog capped cross-model throughput; every arm fired 5/5 and
every control floored 0/5, so n=5 sufficed. The n=12 targets below are the original plan of record.*

| family | F1 (dangling + no_ref) | F2 (bib_dangling + bib_complete) |
|---|---|---|
| claude-sonnet-4-6 | n=12/arm (24 runs) | n=10/arm (20 runs) |
| deepseek-v4-pro | n=12/arm (24 runs) | n=10/arm (20 runs) |
| glm-5.2 | n=12/arm (24 runs) | n=10/arm (20 runs) |

Total: **~88 runs** @ ~240–300 s, `--jobs 2` (cross-model runners run hotter; the watchdog caps
container disk). F3 stays gpt-only (it is an exploratory latch-2 confirmation, not a gate).

**Execution (the agent's runbook).**

```bash
cd /REDACTED_HOME/workspace/personal/repos/quine-lab/experiments/active/structural-elicitation/elic-29-nonse-open-latch

# 1. Clean the 3 claude pilot replicates (smoke, not scored) so replicate indices do not collide.
rm -rf runs-claude runs-deepseek

# 2. Claude F1 + F2 (disk-safe driver runs F1 then F2 with the watchdog attached).
#    Profile path is relative to p13 root: ../profiles/<name>.env
./_xm_run.sh ../profiles/claude-sonnet.env claude 12 2 240

# 3. Deepseek F1 + F2.
./_xm_run.sh ../profiles/deepseek-v4-pro.env deepseek 12 2 240

# 4. Score each family independently (scorer reads ALL run dirs under its root).
python3 analysis/score.py runs/             # gpt canonical
python3 analysis/score.py runs-claude/      # claude
python3 analysis/score.py runs-deepseek/    # deepseek
```

`_xm_run.sh` already wraps the disk watchdog (kills any container whose `/var/lib/rt` exceeds
`CAP_MB=400`, hard `df` backstop at `LOW_GB=14`), runs F1 (`--family dangling`) then F2
(`--family bib`) into `runs-<suffix>/`, and sets `RUNS_DIR` so the scorer picks them up. **No seed
or scorer changes** — the gpt seeds and `analysis/score.py` are byte-identical across families (the
DV is a filesystem-existence event, model-agnostic).

**Pre-registered symmetric decision rule (identical to E9, applied per family).**

For each family, F1 **fires** iff: `dangling.target_created ≥ 0.5` AND
`dangling − no_ref ≥ 0.4` AND Fisher `p < 0.05`. A **high-perception NULL** (agent `cat`/`grep`-ed
the handbook, saw the 404, did not create the target, with perception covariate ≥ 0.5) is a
**FALSIFICATION** of the non-SE claim for that family — never "inconclusive." Report all three
families' F1 numbers in one table; do not cherry-pick.

**Post-run absorption (what the agent updates).**

1. `14.05-nonse-open-latch/README.md` — replace the "Still open: cross-model replication" line with
   the three-family F1/F2 table + verdict.
2. `claim.md` banner — update the gate-(2) cross-model line from "pending" to the result:
   - If ≥2 of 3 families fire: *"Gate (2) closed in our favour, cross-model — F1 fires on gpt,
     claude, deepseek; the discrepancy pull is substrate-general across three model families."*
   - If only gpt fires: *"Gate (2) gpt-only; the non-SE positive did not transfer to claude/deepseek.
     Scope the LBA's non-SE panel to gpt with an explicit single-model caveat, or pull it."* (PI call.)
3. `lba/manuscript.md` + `lba/metadata.yaml` — update the non-SE panel numbers from "gpt, n≈12" to
   the three-family table; rebuild the PDF (`make` in `lba/`).
4. DVC-track the new run trees:
   `cd ../../.. && dvc add experiments/active/structural-elicitation/elic-29-nonse-open-latch/runs-claude experiments/active/structural-elicitation/elic-29-nonse-open-latch/runs-deepseek experiments/active/structural-elicitation/elic-29-nonse-open-latch/runs-glm` (done — all three lanes tracked as `.dvc` pointers and pushed to the private remote).

**Estimated run cost:** ~88 runs @ ~240–300 s, `--jobs 2` ≈ 3–4 h end-to-end (two families, serial;
the watchdog prevents `--jobs 4` cross-model without a second disk). API cost: claude + deepseek are
the cheaper lanes. **This is the single highest-value pre-deadline spend** — it removes the one
asymmetry a reviewer can weaponize without changing a single claim.

---

## 5 · Open decisions for the PI (Hao Ke)

1. **Gate (1) acceptance criterion.** Is the **+x affordance** (non-completion,
   non-ordinal, syscall DV) an acceptable close of gate (1), or must gate (1)
   specifically demonstrate a **non-linguistic discrepancy/gap** (in which case the
   relational SHA-hole is the path, and a NULL there means **report gate (1) OPEN**
   and defend NTP via existing 13.04/13.07)? This decides whether E4 ships a close or
   an honest-open.
2. **Claim/manuscript correction.** Authorize updating `claim.md`'s 2026-06-18
   banner and any `lba/manuscript.md` text that cites 13.24 as "structural / survives
   de-linguistification" — A3 is **not** resolved by the ordinal gap. (This plan only
   recommends; the edit is a separate, PI-gated action.)
3. **E5 non-chore control design.** The non-chore valence-control arm is the
   load-bearing new piece. Does the PI prefer an **opaque ownerless date-conflict**,
   a **relational hole**, or a **metadata structure shared with E4**? This couples
   E4 and E5.
4. **Success-criterion policy (program-wide).** Adopt **"primary arm must pass BOTH
   the binary `gap≥0.5 ∧ gap≥control+0.4` AND Fisher p<0.05; secondaries reported
   with both"** as the standing rule, and **raise primary n to 8–12**? This resolves
   the E5/E6/E8 threshold contradictions uniformly.
5. **E7 sequencing.** Confirm E7 is **deferred** behind E4/E5 and run only if budget
   remains — and only as two non-factorial panels (the factorial is abandoned).
 6. **Compute budget envelope before 2026-07-20.** Which waves are authorized, in what
    order, and is deepseek-v4 a default third lane or a positive-result-only add-on?
 7. **E9 — RESOLVED (2026-06-20/21).** E5 confounded `non-SE` with `both-latches-closed`
    (see E9 card §4); the theory's two-latch model predicted E5's null on any substrate.
    E9 was built and run (`14.05`): F1 dangling cross-ref **fires** in every family
    tested (gpt 10/13 vs 0/12, p=0.0001; claude 5/5, deepseek 5/5, glm 5/5), F2
    citation fires, F3 broken-CSV NULL as predicted. Gate (2) is **reopened** and the conditional resolved to the
    **"objective-defect + additive/determinate-resolution-bound"** branch — the
    discrepancy pull is substrate-general, not SE-specific. `claim.md` and
    `lba/manuscript.md` now carry that wording.

---

## 6 · Cross-cutting risks

1. **Reflexive capture (the dominant risk).** "Is this non-linguistic / non-SE / a
   clean separation" is an experimenter coding judgment made inside the program's
   field. Mitigation is structural, not rhetorical: pre-registered DV thresholds,
   neutral/blind-able scoring (mktemp staging already strips the arm name from
   `/proc/self/mountinfo`), **symmetric** decision rules (high-perception NULL =
   falsification), and a single pre-registered **primary** per experiment to bound
   experimenter degrees of freedom. State every positive as **discounted-for-program
   -pull** per `claim.md` §Reflexive-capture.
2. **A2 bounds the program to the underspecified regime.** A2 (locked; sharpened by
   `14.06`) concedes the pull is subordinate to **specific tasking**, not an
   instruction-independent third source. E5/E8 cannot escape this by **framing**; only the **non-chore valence control** (E5) and a
   **non-SE motif in the model wave** (E8) provide discriminating evidence. Without
   them, the honest ceiling is "substrate-breadth of an underspecified-regime
   implicit task," and every card must say so.
3. **Harness round-trip silently mutating the DV.** The 13.24 recovery
   (`chown -R` + `cp -a`) mutates mtime/ctime/mode. Any metadata-axis DV (E4, E7)
   **must** capture in-container before recovery (§2.1) and pass the round-trip
   self-test, or the DV is destroyed before scoring.
4. **Floor effects / collapsing denominators.** Multi-choice and metadata-perception
   designs (E4, E6) and the weak gemini-flash lane (E8) can floor invocation/
   perception toward zero, leaving share statistics computed on a handful of
   self-selected runs. Mitigation: **unconditional** primary DVs over all
   non-infra-failed runs, a mandatory **activity/any-invoked floor**, and an n=2
   pilot that **gates the wave** before any confirmatory spend.
5. **Re-linguification through the perception channel.** Any perceived metadatum
   renders as readable text in the agent's context. The defensible claim is "no
   readable **ordinal/instruction token**," not "no readable symbol" — concede this
   in the write-up (E4, §2.6).
6. **DV gameability the agent can trip without structure.** Tidying sweeps (gap
   existence DV), generic exec curiosity (affordance DV), tolerance-band luck
   (metadata DV). Mitigation: matched **symmetric** controls that quantify the
   false-positive floor, anti-tidying exclusions, real-syscall binding (§2.2), and
   the pointless-runnable arm (E5.3).
7. **Opportunity cost / sequencing.** Compute is finite before 2026-07-20. E4 and E5
   are the two binding gates; E7 hardens an already-shipped result deeper in the very
   distribution E5 must escape. **Build E4/E5 first; E8 last and cheap; E7 only if
   budget remains.**
