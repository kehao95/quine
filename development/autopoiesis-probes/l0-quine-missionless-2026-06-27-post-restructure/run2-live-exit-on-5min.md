# Run 2 — live repo, exit ON, 5-min wall-clock

Follow-up to run 1, testing the **exit-enabled** behavior on the **live repo**
(after the red-line collapse and the pure-router rewrite). Config delta vs run 1:
workspace = live repo (not a clean worktree); `QUINE_EXIT_ENABLED=1`; wall-clock
300s. Head `3539a7cba` (227-line router AGENTS.md, single red line, no ownership-
paralysis prohibition). Tape: `run2-live-exit-on-5min-tape.jsonl`.

## Outcome

| metric | value |
|---|---|
| termination_mode | `signal` (wall-clock self-exit after 5m0s) |
| exit_code | 0 |
| duration | 300s (full wall-clock) |
| turns | 3 |
| tokens in / out | 20.5K / 690 |
| fs mutations | **0** |

## Findings

1. **No early exit.** The agent did not call `exit`; it ran the full 300s and was
   wall-clock-killed. Quine is not bail-prone even when it cannot act — consistent
   with run 1 and the prior arc.
2. **No action — but from time starvation, not paralysis.** All 2–3 turns were
   orientation (sensed the live repo, noted the dirty tree + 2 commits ahead,
   declared maintenance intent). The red-line paralysis from the earlier live run
   is gone (the unowned-diffs prohibition was removed), but the agent never
   reached the action phase.
3. **The dirty live repo caps turn-rate.** ~100s/turn (vs ~15s/turn on a clean
   worktree) because the dirty workspace produces large tool outputs (git diff of
   100+ files, big dir listings) that slow each turn. 5 min is insufficient to
   pass orientation on the live substrate; even the prior 600s live run reached
   only 3 turns.

## Conclusion

The exit flag behaves correctly (no spurious early exit). The live-repo substrate,
not the agent's judgment, is the bottleneck: the dirty shared workspace starves
the agent of turns. To test genuine live-repo behavior, either (a) extend the
wall-clock past the orientation phase, or (b) clean the workspace so orientation
is fast. Run 1 (clean worktree, 41 turns) remains the high-resolution data point;
this run mainly documents that 5 min on a dirty live repo yields no signal.
