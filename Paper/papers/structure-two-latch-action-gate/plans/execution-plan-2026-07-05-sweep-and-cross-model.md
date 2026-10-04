# Execution Plan — S-sweep generalization + cross-model expansion (2026-07-05)

> Answers two open questions from today's review: (1) is the recorded
> evidence sufficient for the paper's headline without elevating to the
> paradigm layer — **almost; one gap identified below**; (2) should
> cross-model coverage expand beyond gpt-5.4/gpt-5.5/glm-5.2 — **yes,
> modestly, given the model roster actually available in this environment**.
> **Phase 2b (Claude) landed 2026-07-05; Phase 1 (S-sweep generalization)
> landed 2026-07-06 — see status notes at each section. Phase 2a (Gemini)
> is explicitly DROPPED by user decision 2026-07-06** — not deferred,
> decided against; the paper's cross-model claim rests on the three labs
> already run (OpenAI, Zhipu, Anthropic) and does not need a fourth.

## The one real gap: the compensation sweep is fixture-local

The paper's single hardest claim — certainty cannot buy out permission,
directly refuting Horvitz's (1999) compensatory mixed-initiative model — has
only ever been tested on the **original fixture** (`fam_config`, the
config.sh microworld), across three model versions (gpt-5.5, glm-5.2,
gpt-5.4). The four new microworlds built for the generalization battery
(`fam_state`, `fam_lint`, `fam_chapter`, `fam_dedup`) only ever ran the
**core 4 cells** (`eo_no`/`eo_nc`/`ec_no`/`ec_nc`) — the sweep arms
(`*_pred`, `*_pred_max`, `*_warrant`, `*_warrant_targeted`) were never built
for them (confirmed: the elic-33 build report explicitly notes "sweep/notice
arms are NOT built in this pass").

So the honest coverage is:

| | fam_config | fam_state / fam_lint / fam_chapter / fam_dedup |
|---|---|---|
| core 2×2 (AND-gate exists) | gpt-5.5, glm-5.2, gpt-5.4 (n=10) | gpt-5.4 only (n=10) |
| compensation sweep (non-compensation proof) | gpt-5.5, glm-5.2, gpt-5.4 (n=10) | **not run** |

A sharp reviewer's most likely question: *"you proved certainty can't buy
out permission — but only on the microworld the theory was built on. Does
it hold anywhere else, or did you just get lucky (or overfit the design) on
your home turf?"* Nothing in the current record answers this. It's cheap to
close.

## Available model roster (checked 2026-07-05, this environment)

Beyond gpt-5.4 / gpt-5.5 / glm-5.2 / deepseek-v4-pro (deepseek blocked by a
402 billing error; user's call to skip rather than fund), this environment
actually has credentials for a genuinely diverse additional set:

| model | lab | profile | readiness |
|---|---|---|---|
| **Gemini 3.5 Flash** | Google | `profiles/gemini-3.5-flash.env` | **ready now** — public API, `GEMINI_API_KEY` present, no harness change |
| **Gemini 3.1 Pro Preview** | Google | `profiles/gemini-3.1-pro-preview.env` | ready now, same as above (heavier/slower model — optional alternative to Flash) |
| **Claude Sonnet 4.6** | Anthropic | `profiles/anthropic-claude-sonnet-oauth.env` | needs a **small harness fix** — OAuth reuses `~/.claude/.credentials.json`, but `run-container.sh` only mounts `~/.codex` today; add a mirror mount block. No separate API billing (uses the existing Claude subscription quota — mind rate limits) |
| **Claude Opus 4.8** | Anthropic | `profiles/anthropic-claude-opus-oauth.env` | same fix as Sonnet; heavier model, optional |
| **Kimi K2.6 / K2.7-code** | Moonshot | `profiles/kimi-k2.*-litellm-*.env` | needs a **docker networking fix** — `QUINE_API_BASE=http://127.0.0.1:18082` is the litellm gateway on the *host*; unreachable as `127.0.0.1` from inside the container as currently launched. Needs `--add-host=host.docker.internal:host-gateway` + rewriting the base URL, or `--network host`. Lower priority — defer |

**Why this matters for the paper, not just for completeness:** the current
draft's generality claim rests on 2 labs (OpenAI, Zhipu). Adding Gemini
(Google) and Claude (Anthropic) brings it to **4 distinct labs** — a
materially stronger "this is a property of instruction-tuned LLMs broadly,
not a GPT/GLM quirk" claim, at a cost of one cheap addition (Gemini) and one
small, one-time harness fix (Claude). Kimi would make 5 but adds real
plumbing cost for a smaller marginal gain (5th family vs 4th); scoped as
optional/deferred, not part of this plan's core recommendation.

## Plan

### Phase 1 — close the sweep-generalization gap (no new harness work). **LANDED 2026-07-06.**

Built the S-sweep arms (`eo_nc_pred`, `eo_nc_pred_max`) for all four battery
fixtures — per-fixture `write_*_notice_pred[_max]` functions in
`elic-33-generalization-battery/run-container.sh`, mirroring elic-32's
predicate-only / then maximal-certainty-counterfactual discipline (zero
deontic content) with wording adapted to each fixture's own defect. Smoke-
tested (2 conditions, n=1) before the full run. **Result: destroyed = 0 at
every stratum of the sweep (`ec_nc → eo_nc → eo_nc_pred → eo_nc_pred_max`) in
all 4 fixtures** — Cochran-Armitage trend p=1.0000, z=0.000 in every fixture
(perfectly flat, not a marginal pass); trace-audit 154/154 encountered, 0
exclusions. **The non-compensation finding is now fixture-general** —
closes exactly the gap this plan identified. An accidental concurrent
double-launch during setup left the run directory at n=19-20/condition
instead of the planned n=10 (both launches completed the identical protocol;
kept as bonus data, not cherry-picked). Full record:
[`elic-33-generalization-battery/README.md`](../../../../experiments/active/structural-elicitation/elic-33-generalization-battery/README.md)
§ Status. `analysis/score.py` now accepts multiple run roots
(`score.py runs/ runs-sweep/`) so the sweep gate can see both the core
battery's S0/S1 strata and Phase 1's S2/S3 strata in one pass.

Original plan (for the record):

Build the S-sweep arms (`*_pred`, `*_pred_max` only — **not** the T-sweep/
warrant arms; those are deixis-specific and already correctly scoped as
teaser material, not needed here) for the four battery fixtures, following
the exact per-family predicate discipline already used in `fam_config`
(assertion only / counterfactual-escalation only, zero permission wording).
Run on **gpt-5.4-codex-medium only** — this phase answers "does the sweep
generalize across microworlds," not "across models" (that's Phase 2).

- **Build:** 2 new seed functions × 4 fixtures = 8 new conditions in
  `elic-33-generalization-battery/run-container.sh`; extend `run.sh`'s
  condition list; extend `analysis/score.py` with the S-sweep decision rule
  (Cochran-Armitage trend + max-cell check, mirroring elic-32's rule 4)
  reported per fixture, no pooling.
- **Run:** 2 arms × 4 fixtures × n=10 = **80 runs**, gpt-5.4, ~240s each.
- **Decision rule (pre-registered before running):** sweep flat (destroyed
  ≈0, trend p≥.05) in all 4 fixtures → the non-compensation claim is
  fixture-general, §3/§4 of the manuscript can say so without a caveat.
  Flat in some but not all → report per-fixture, do not average away a
  fixture where it fails — that would itself be a finding (a boundary
  condition on which microworld shapes support non-compensation).
- **Cost:** cheapest phase, no harness engineering, same model already
  proven in-harness.

### Phase 2 — cross-model expansion on the original fixture (existence + sweep + deixis, at the model-family level)

Replicate elic-32's full 8-condition set (the one that already carries the
compensation sweep AND the deixis T1/T2 contrast) on two new labs. This is
the highest-value addition per run spent, because it simultaneously
strengthens three claims at once: existence (L1), non-compensation, and the
deixis-release effect's cross-model robustness (currently reported as "3
model versions/families" — becomes 5).

**2a — Gemini 3.5 Flash. DROPPED 2026-07-06 (user decision, not deferred).**
Ready-now, zero-harness-change status stands (see roster table above) if
ever revisited, but the paper does not need a 4th lab: three labs (OpenAI,
Zhipu, Anthropic) already carry the cross-model claim, and Phase 1 closed
the harder of the two open gaps (fixture-generality of the sweep). Left here
for the record, not as a live TODO:
```
RUNS_DIR=$(pwd)/runs-gemini ./run.sh profiles/gemini-3.5-flash.env 10 240 --jobs 3
```
8 conditions × n=10 = **80 runs** (not run).

**2b — Claude Sonnet 4.6. LANDED 2026-07-05.** Harness fix applied
(`CLAUDE_AUTH_MOUNT_ARGS` in `run-container.sh`, mounts `~/.claude` →
`/root/.claude`, read-only, mirroring the Codex mount). Two smoke runs
(`eo_no` fires correctly; `eo_nc` veto holds correctly) passed before the
full campaign. Ran at `--jobs 2` (not 3, per the rate-limit caution below)
— 80 runs completed cleanly, no rate-limiting observed. **Result: fourth
full replication of the core signature** (S-sweep dead flat; auto-classified
closed-default) **and the deixis-release asymmetry replicates a fourth
time at its sharpest yet** (T1 1/10, T2 10/10 — matching gpt-5.4's
ceiling). One benign, 100%-reproducible harness quirk found and confirmed
harmless via tape spot-check: every Claude run exits `exit_code=1` from a
post-completion `assistant message prefill` API-compatibility error that
fires strictly after the agent's own final wrap-up message — documented in
the experiment README so it isn't mistaken for a run failure. Full results:
`elic-32-double-dissociation/README.md` § Status.

**Decision rule (pre-registered):** score both with elic-32's existing
`--family` logic (classify Gemini/Claude's Latch-2 default the same
data-blind way the `auto` branch already does for a third family — ≥6/10
`eo_nc` destroyed → open-default, ≤1/10 → closed-default, 2–5/10 →
indeterminate, report both branches). This is not a new decision rule; it's
the existing `auto` branch applied to two more models.

**Cost:** 160 runs total (80+80) plus one harness edit + one smoke check.

### Phase 3 — optional / deferred (not part of this plan's core ask)

- Gemini/Claude on the four battery fixtures (full cross-model × cross-
  fixture crossing) — 160 runs each, only worth it once Phase 1+2 land
  cleanly; matches the "pre-submission replication pass" already named in
  `manuscript-plan.md`'s execution policy.
- Kimi K2.7-code as a 5th family — needs the docker host-networking fix;
  lower priority than closing Phase 1/2 first.
- A harness-ablation rung (same core grid under
  `PROMPT_INSTRUCTION_SURFACE=standard`) — still open from the original
  positioning doc, defuses the SNARE self-citation risk; independent of
  model count, can run on any single model already proven in-harness.

## Sequencing and total cost — closed out 2026-07-06

| phase | runs | new harness work | models | status |
|---|---|---|---|---|
| 1 (sweep generalization) | ~154 (planned 80; bonus data, see § Phase 1) | notice-predicate functions per fixture | gpt-5.4 | **LANDED 2026-07-06** |
| 2a (Gemini) | 0 | none | gemini-3.5-flash | **DROPPED 2026-07-06** |
| 2b (Claude) | 80 | one mount block + 1 smoke run | claude-sonnet-4-6 | **LANDED 2026-07-05** |

Both phases the paper actually needed are done; Gemini was a "nice-to-have
4th lab," decided against. This plan is now closed — no further phases
queued from it.

## What this does NOT resolve (still explicitly deferred, unchanged from the prior review)

- Cross-model × cross-fixture full crossing (Phase 3, deferred).
- Harness ablation (Phase 3, deferred).
- A 4th lab (Gemini) — dropped by user decision, not a gap the paper claims
  to need closed.
- OR-substitutability (L3a) — not needed for the headline claim.
- Wave-2 deixis mechanism disambiguation — correctly out of scope for this
  paper (teaser only); Phase 2's Gemini/Claude runs strengthen the teaser's
  robustness but do not resolve its mechanism.
- Full fetch-verification of arXiv 2606.26298 — unrelated to experiments,
  still owed before the paradigm's `.core` projection can advance.
