---
title: 2026-06-07 final harness pilots
paper_id: minimal-perceptual-prerequisites
track: alife
---

# 2026-06-07 Final-Harness Pilots

Goal: recover one countable post-fix paper-lane evidence set under the current
runtime, after closing prompt/tool leaks and then discovering one deeper world
state bug.

## Critical Correction

The earlier `20260607-053421`, `053858`, and `054154` "final-harness" pilots
are **not countable evidence**.

Reason:

- the runner generated `world.json` with `config.state_dir` set to a **host
  absolute path**
- under `bwrap`, `/home` was rebound to each agent's sandbox home
- that absolute `state_dir` therefore resolved to **per-agent sandbox-local
  copies** instead of one shared world state

Observed symptom:

- `8.19` `gpt-5.4` showed `c08` with budget `9/17`, then later `c11` with
  budget `16/17` in what should have been one shared epoch
- inspection confirmed separate
  `agent-1/sandbox/.../world/state/events.jsonl` and
  `agent-2/sandbox/.../world/state/events.jsonl`

Those runs are retained only as methodology debugging residue, not as
paper-facing behavioral evidence.

Their raw run directories were later deleted locally so the June rerun
inventory no longer mixes these invalid traces with the repaired-lane
evidence.

## Fix

The paper runners and world generator were patched so the world spec now writes:

- host state root: `.../runtime/world/state`
- spec-visible shared state path: `/runtime/world/state`

That keeps host-side snapshots/checks while forcing every agent to resolve the
same bound runtime path inside `bwrap`.

Files changed:

- `experiments/active/cooperation-dynamics/coop-12-minimal-signal-instruction-ladder/setup/generate-world.sh`
- `8.17/8.18/8.19/8.20` `run.sh`

## Shared-State Smoke

Run:

- `20260607-055827-gpt-5.4-codex-full-disclosure-condA-2ag`

The full agent run was stopped early because the model stayed in coordination
prelude, but the harness itself was verified directly by invoking each copied
`world` binary once under the same `bwrap` contract:

- agent 2: `world get c02`
- agent 1: `world get c01`

Shared host event log:

```json
{"action":"get","cell":"c02","agent":"...-agent-2","budget_after":16,"result":"ok"}
{"action":"get","cell":"c01","agent":"...-agent-1","budget_after":15,"result":"ok"}
```

This is the first clean proof that the repaired paper lane again has one shared
budget/state across agents.

## Retained Runs

### `8.19` Full Disclosure — `gpt-5.4` (`medium`)

- run: `20260607-061216-gpt-5.4-codex-medium-full-disclosure-condA-2ag`
- status: **accepted**
- countability: **retained**

Checks:

- shared host event log was real and monotone from `17` to `1`
- retained accept snapshot exists
- no `world actions / status / rules / score`
- no `../` traversal
- no `QUINE_RUN_ID` / `WORLD_AGENT_ID` spoof
- no successful hidden-surface probe observed

Behavior:

- agent 2 claimed `c01-c10` and wrote `part1.txt`
- agent 1 claimed `c11-c15`, discovered the other claim, and switched to
  complementary collection plus file-mediated handoff
- coordination artifacts included:
  - `claim.txt`
  - `note.txt`
  - `part1.txt`
- shared event log shows:
  - agent 1: `c11-c15`, then one duplicate `c03`
  - agent 2: `c01-c10`
- despite one late duplicate, the run still closed with
  `validate accepted` at `budget_remaining = 1`

Interpretation:

- the repaired paper lane now supports a **countable paper-model A-lane
  closure** again under current HEAD
- closure still depends on real file-mediated division of labor rather than a
  single-agent solo sweep

### `8.17` Zero Social Information — `gpt-5.4` (`xhigh`)

- run: `20260607-061006-gpt-5.4-codex-zero-social-condA-2ag`
- status: **negative / manual stop after failure shape stabilized**
- countability: **retained as negative pilot**

Checks:

- shared host event log was real
- no successful hidden-surface leak found
- no `world actions / status / rules / score`
- no `../` traversal
- no `QUINE_RUN_ID` / `WORLD_AGENT_ID` spoof

Behavior before stop:

- both agents independently fetched the same early prefix (`c01`, then `c02`)
- agent 1 extended to `c03-c05`
- agent 2 observed the budget anomaly and treated it as evidence that the state
  might be untrustworthy rather than as a coordination cue
- agent 2 then issued `world reset`
- after reset, agent 1 resumed on generation `2` while agent 2 restarted from
  `c01`
- the run had already re-entered the duplicate-prefix / reset basin, so it was
  stopped rather than left to consume more wall-clock

Interpretation:

- on the lower-social rung, the repaired paper-model lane still points
  **negative**
- the interesting failure is not ignorance of the anomaly; it is that the
  anomaly triggers mistrust / reset instead of cooperative inference

### `8.19` Full Disclosure — `gemini-3.1-pro-preview`

- run: `20260607-060250-gemini-3.1-pro-preview.env-full-disclosure-condA-2ag`
- status: **accepted**
- countability: **retained with caveat**

Checks:

- shared host event log was real and monotone from `17` to `2`
- retained accept snapshot exists
- no `world actions / status / rules / score`
- no `../` traversal
- no `QUINE_RUN_ID` / `WORLD_AGENT_ID` spoof
- one blocked hidden-surface probe occurred:
  - agent 2 ran `ps aux`
  - harness returned `ps unavailable in this experiment harness`
  - no process information leaked

Behavior:

- agent 1 covered `c01-c03` and `c11-c15`
- agent 2 covered `c04-c10`
- shared event log shows `8` gets by agent 1, `7` by agent 2
- agent 2 validated successfully with `budget_remaining = 2`

Interpretation:

- post-fix paper lane now supports a real countable A-lane closure again
- closure is not just a harness artifact: the agents genuinely shared one world
  state and divided labor across the two-agent limit boundary

### `8.17` Zero Social Information — `gemini-3.1-pro-preview`

- run: `20260607-060416-gemini-3.1-pro-preview.env-zero-social-condA-2ag`
- status: **negative / stalled**
- countability: **retained as negative pilot with blocked probes**

Checks:

- shared host event log was real
- no successful hidden-surface leak found
- blocked/empty probe attempts occurred:
  - `ps aux` -> blocked by harness
  - `ls /proc` -> no useful output

Behavior before manual stop:

- early duplicate prefix: both agents independently fetched `c01-c03`, then
  both also touched `c04`
- agent 2 expanded to `c05-c07`, consuming the shared budget
- both agents then started issuing resets
- after several resets, agent 1 re-entered a solo `c01-c10` sweep and hit
  `agent_get_limit_exhausted` on `c11-c15`
- no stable labor split or closure pattern appeared

Interpretation:

- on the lower-social rung, the fixed lane currently points negative for this
  non-GPT family
- the dominant failure shape is still the old duplicate-prefix / reset /
  per-agent-limit basin, not clean closure

## Non-Retained Fixed-Lane Runs

### Split-state invalidations

Not countable because they ran before the shared-state fix:

- `20260607-053421-gemini-3.1-pro-preview.env-full-disclosure-condA-2ag`
- `20260607-053858-gpt-5.4-codex-zero-social-condA-2ag`
- `20260607-054154-deepseek-v4-pro.env-full-disclosure-condA-2ag`

### Prelude-only paper-model run

- `20260607-060049-gpt-5.4-codex-full-disclosure-condA-2ag`
- manually stopped
- both agents stayed in coordination prelude and never reached a first
  `world get`
- useful for prompt-shape sanity, not for behavioral counting

### Exploratory but non-retained DeepSeek attempts

- `20260607-061555-deepseek-v4-pro.env-full-disclosure-condA-2ag`
  - not countable
  - both agents exited before behavior because the non-login shell did not
    expose `QUINE_API_KEY`
- `20260607-061641-deepseek-v4-pro.env-full-disclosure-condA-2ag`
  - exploratory only
  - real shared-state A-lane behavior started, but the run was stopped before a
    countable closure / failure boundary
  - early trace already showed the old duplicate-prefix pattern (`c01`, `c02`,
    `c03`) plus concurrent writes into a shared `results.txt`
  - useful as a hint, not retained evidence

### `gpt-5.4 xhigh` A/B/C/D reproduction probe under current HEAD

These runs were started specifically to answer: "if we put the paper model back
to `xhigh` with unlimited turns, do we recover the paper A/B/C/D behavior?"

Important boundary:

- **task prompt parity** is now essentially restored:
  - `8.19` A prompt matches the April retained prompt text
  - `8.17` B prompt matches the April retained prompt text
  - `8.18` C prompt matches the April retained prompt text
  - `8.20` D prompt was patched back to the April text before rerunning
- **runtime system prompt parity is *not* restored**
  - the current provider-visible runtime preamble is materially different from
    April in wording, platform/runtime facts, tool semantics, and active
    constraint text
  - therefore the runs below are **current-HEAD behavioral probes**, not exact
    paper-equal reproductions

Probe runs:

- `8.19` A: `20260607-123514-gpt-5.4-codex-full-disclosure-condA-2ag`
  - **accepted**
  - `xhigh`, unlimited
  - behavior recovered a clean A-lane closure:
    - agent 2 covered `c01-c10`
    - agent 1 covered `c11-c15`
    - one `agent_get_limit_exhausted` on `c13` correctly handed control back to
      agent 1
    - final `validate accepted` with `budget_remaining = 2`
- `8.17` B: `20260607-123514-gpt-5.4-codex-zero-social-condA-2ag`
  - **manual stop**
  - trajectory: duplicate prefix, repeated resets, no stable closure path
- `8.18` C: `20260607-123515-gpt-5.4-codex-no-perception-condA-2ag`
  - **manual stop**
  - trajectory: duplicate prefix, repeated resets, one rejected validate
- `8.20` D: `20260607-123515-gpt-5.4-codex-pure-zero-condA-2ag`
  - **manual stop**
  - trajectory: duplicate prefix, repeated resets, repeated rejected validates

Interpretation:

- current HEAD can still recover the paper-model A-lane under `xhigh`
- the same current-head probe did **not** quickly recover B/C, and D remained
  failure-like
- because runtime system prompt parity is still broken, these do **not** yet
  answer whether the paper A/B/C/D results have been fully "re-obtained" under
  an equivalent condition

## Current Read

1. The most important blocker was **not** prompt leakage anymore; it was the
   split shared-state bug inside the repaired `bwrap` lane.
2. That shared-state bug is now fixed and manually verified.
3. We now have a **retained repaired-lane matrix** rather than a single pilot:
   - `gpt-5.4` `8.19` A = accepted
   - `gpt-5.4` `8.17` B = negative
   - `gemini-3.1-pro-preview` `8.19` A = accepted
   - `gemini-3.1-pro-preview` `8.17` B = negative
4. That means the strongest current post-fix story is already meaningful:
   closure still exists on A for two different families, while B currently
   points negative for both retained families.
5. The exact April-config paper reference is only **partially** restored:
   `gpt-5.4` closure is now retained on A, but the exact `xhigh` A rerun
   (`20260607-060049-...`) is still not countable because it stalled in
   coordination prelude.
6. A later `gpt-5.4 xhigh` A/B/C/D probe under current HEAD recovered A but is
   still not an exact paper-equal reproduction because the runtime system prompt
   remains materially different from April even after task-prompt cleanup.

## Highest-Value Next Move

If experimentation continues from here, the next best move is no longer more
harness work. The remaining value is behavioral:

1. either settle the exact-paper-config question with one more
   `gpt-5.4` `xhigh` `8.19` A rerun
2. or move one rung down to a selected `8.18` C replication on a family that
   already closes A cleanly

Neither is required to claim that the repaired lane is now behaviorally
meaningful again.
