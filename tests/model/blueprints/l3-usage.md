# L3 Usage Blueprints

These are planned `L3` usage evaluations.

They are intentionally not registered yet.
The point is to make the design, harness shape, and likely failure modes
concrete enough that later implementation can proceed directly.

## Portfolio Summary

| Future id | Feature | Priority | Future canonical path |
|----------|---------|----------|-----------------------|
| `shell-envelope-timeout-truncation-cwd` | `shell-envelope` | `P1` | `tests/model/l3-usage/shell-envelope/timeout-truncation-cwd/` |
| `response-governance-exit-disabled-impossibility` | `response-governance` | `P1` | `tests/model/l3-usage/response-governance/exit-disabled-impossibility/` |
| `context-pressure-anchor-exec-carryover` | `context-pressure` | `P1` | `tests/model/l3-usage/context-pressure/anchor-exec-carryover/` |
| `tape-lineage-parent-child-incarnation` | `tape-lineage` | `P1` | `tests/model/l3-usage/tape-lineage/parent-child-incarnation/` |
| `idle-explicit-suspension-resume` | `idle` | `P1` | `tests/model/l3-usage/idle/explicit-suspension-resume/` |

## shell-envelope-timeout-truncation-cwd

- Goal: prove that the model can use explicitly documented shell timeout, output truncation, and startup working-directory semantics correctly.
- Pressure model: the prompt can name the mechanisms directly because this is `L3`; the difficulty comes from coordinating all three without misreading tool results.
- Environment frame: set `QUINE_SH_DEFAULT_TIMEOUT_SECONDS` low enough to interrupt a slow command, set `QUINE_OUTPUT_TRUNCATE` low enough to force a truncation notice, and start Quine in a non-default `QUINE_WORK_DIR` containing a sentinel file.
- Prompt red lines: none beyond normal `L3` honesty; this experiment may say `timeout`, `truncated output`, and `working directory` explicitly.
- Scorer shape: confirm the model provoked one timeout, recognized one truncation notice, found the startup cwd sentinel, and wrote the expected markers to fd 4 without confusing truncated output for full output.
- Expected failures: the model may accidentally choose commands that finish before timeout, may ignore the truncation banner, or may inspect the wrong directory and think cwd semantics failed.
- Adjustment knobs: timeout seconds, output length, sentinel path complexity, and whether timeout proof comes from sleep or from a naturally noisy command.
- Follow-on coverage: the explicit timeout-resume and timeout-terminate slices are now active evaluations, so the remaining planned shell-envelope work is the broader truncation/cwd composition case rather than the paused-job operation-handle path.

## response-governance-exit-disabled-impossibility

- Goal: prove that the model can behave correctly when explicit termination and impossible-work posture are altered by env control.
- Pressure model: the prompt names `exit` unavailability and impossible-work rules directly; the test measures whether the model routes success/failure through the remaining allowed surfaces.
- Environment frame: disable `exit`, keep `QUINE_FAIL_ON_IMPOSSIBLE=1`, and provide one impossible branch plus one clean reporting path.
- Prompt red lines: do not hide the mechanism; this is an operation check, not a discovery test.
- Scorer shape: verify the model does not call `exit` when it is unavailable, does not bluff success on the impossible branch, and emits the expected diagnostics through the surviving channels.
- Expected failures: the model may still attempt `exit`, may loop in text, or may report success without proving impossibility.
- Adjustment knobs: whether the impossible branch is mathematical, filesystem-based, or protocol-based; and whether success and failure are mixed in one mission or separated.

## context-pressure-anchor-exec-carryover

(Reworded 2026-07-03: the `wisdom` exec parameter was deleted per
registry-design-brief D5; exec-boundary carryover surfaces are now anchors,
`context/`, and explicit workspace files.)

- Goal: prove explicit operational use of anchor / `context/` / workspace-file carryover across `exec` under named memory pressure.
- Pressure model: the prompt may say that context is scarce, that recall should be checkpointed, and that re-entry is allowed; the test measures whether the model uses the named carryover surfaces correctly rather than merely talking about them.
- Environment frame: lower `QUINE_MEMORY_WARN_TOKENS` and `QUINE_MEMORY_DANGER_TOKENS`, give a multi-phase mission where a compact intermediate summary must survive an `exec`, and preserve enough runtime state to inspect the carryover.
- Prompt red lines: name the mechanism directly, but do not over-script the exact command sequence; leave room for either `anchor`, an explicit file checkpoint, or both if the claim is meant to permit both.
- Scorer shape: verify that the model stores the required compact state, performs a re-entry or other named carryover action, and later reconstructs the answer from the carried state rather than recomputing from the original source.
- Expected failures: the model may keep everything in active context and never test carryover, may store too much verbose state, or may re-enter without preserving the needed summary.
- Adjustment knobs: number of phases, summary size cap, whether anchor and file checkpointing are both allowed or one is mandatory, and whether the re-entry is forced or optional.

## tape-lineage-parent-child-incarnation

- Goal: prove explicit use of session, tape, and parent lineage surfaces across fork or exec boundaries.
- Pressure model: the prompt may directly ask for parent/child/incarnation bookkeeping and require the model to read the corresponding runtime artifacts.
- Environment frame: create one child process or one re-exec boundary, preserve runtime state under `QUINE_DATA_DIR` and `QUINE_RETENTION_DIR`, and ensure the resulting lineage markers are inspectable from the filesystem.
- Prompt red lines: keep the mechanism explicit, but avoid hand-feeding exact file contents; the model should still navigate the runtime surfaces itself.
- Scorer shape: verify correct extraction of `session_id`, `tape_id`, and `parent_session`, plus one correct statement about how exec incarnation changed tape identity or lineage.
- Expected failures: the model may confuse session and tape identity, read stale artifacts, or report parent lineage without proving it from runtime state.
- Adjustment knobs: use fork versus exec, whether an external retained root is required, and how many lineage facts must be recovered for success.

## idle-explicit-suspension-resume

- Goal: prove that the model can explicitly suspend itself with `idle`, then resume correctly when a later external control event arrives.
- Pressure model: this is `L3`, so the prompt may name `idle`, inject, and resume semantics directly; the test is about correct operational use rather than discovery.
- Environment frame: enable the planned `idle` env gate, start one session that must produce an initial response, enter `idle`, then continue only after a later external payload is delivered through the ordinary control surface.
- Prompt red lines: do not over-script the exact turn choreography; the mechanism can be named, but the prompt should still leave the agent responsible for choosing a clean response boundary before suspending.
- Scorer shape: verify three things: the model produced the first required output, explicitly entered `idle`, and after inject resumed from the same process/session to complete the second phase without pretending the session had ended.
- Expected failures: the model may avoid `idle` and keep talking, may suspend before finishing the first boundary, or may resume and treat the new input as a fresh unrelated task instead of process continuation.
- Adjustment knobs: number of pre-idle obligations, whether resume uses poke, inject, or interrupt delivery, and how much state the resumed phase must preserve from before suspension.
