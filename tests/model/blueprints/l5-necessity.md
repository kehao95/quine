# L5 Necessity Blueprints

These are planned `L5` necessity evaluations.

Every experiment here is intentionally stronger than ordinary discovery:
the task should make the intended mechanism physically necessary, or at least so
dominant that a clean scorer can defend the claim honestly.

## Constraint Discipline

For `L5`, negative constraints need an explicit classification:

- **Leakage seal**: acceptable. It closes a side channel or accidental shortcut
  that sits outside the intended task physics, such as helper-private artifacts,
  retained logs, or unrelated host inspection.
- **Strategy steering**: not acceptable. It tells the model which persistence,
  bootstrap, wrapper, alias, or surface family counts as canonical.

If a prompt line bans "inventing" or "creating" a new bootstrap/startup path,
wrapper, alias, sidecar, or competing surface, treat that as a design defect in
the fixture/scorer. Remove the line and move the distinction into environment
physics or scoring.

## Portfolio Summary

| Future id | Feature | Priority | Future canonical path |
|----------|---------|----------|-----------------------|
| `detach-overlap-deadline` | `detach` | `P1` | `tests/model/l5-necessity/detach/overlap-deadline/` |
| `fork-parallel-hypothesis-pressure` | `fork` | `P1` | `tests/model/l5-necessity/fork/parallel-hypothesis-pressure/` |
| `exec-stream-ownership-forcing` | `exec` | `P1` | `tests/model/l5-necessity/exec/stream-ownership-forcing/` |
| `switch-world-clean-final-branch` | `switch-world` | `P1` | `tests/model/l5-necessity/switch-world/clean-final-branch/` |
| `anchor-memory-context-cliff-survival` | `anchor-memory` | `P1` | `tests/model/l5-necessity/anchor-memory/context-cliff-survival/` |
| `sandbox-hostile-artifact-survival` | `sandbox` | `P1` | `tests/model/l5-necessity/sandbox/hostile-artifact-survival/` |
| `process-surface-live-peer-rescue` | `process-surface` | `P1` | `tests/model/l5-necessity/process-surface/live-peer-rescue/` |
| `idle-quiet-standby-pressure` | `idle` | `P1` | `tests/model/l5-necessity/idle/quiet-standby-pressure/` |
| `agents-md-startup-token-lockin` | `agents-md` | `P2` | `tests/model/l5-necessity/agents-md/startup-token-lockin/` |

## detach-overlap-deadline

- Goal: force detached work because one long-running task and one immediate verification task must overlap under a hard wall-clock deadline.
- Pressure model: the model cannot win by waiting serially; if it blocks on the long task, it misses the deadline.
- Environment frame: provide a long-running producer that must continue after launch, a second short task that depends on overlapping time, and a deadline tight enough that ordinary sequential `sh` calls fail.
- Prompt red lines: do not say `detach`, `background`, `keep working while`, or otherwise narrate overlap strategy.
- Scorer shape: verify real job overlap from timestamps or job files, plus final success markers showing both the producer and the concurrent verification completed.
- Expected workaround risk: the model may emulate overlap with shell metacharacters or one giant compound command if the harness allows it.
- Adjustment knobs: deadline tightness, whether single-shell background tricks are forbidden structurally, and whether the long-running task must outlive the final response.

## fork-parallel-hypothesis-pressure

- Goal: force child-lane parallelism because several expensive hypotheses must be explored before time or turn budget expires.
- Pressure model: serial evaluation should be measurably too slow or too budget-expensive.
- Environment frame: arrange multiple independent candidate decoders or search lanes, each costly enough that a single-thread parent cannot check enough of them in time.
- Prompt red lines: do not say `fork`, `parallel`, `delegate`, `children`, or suggest multiple attempts.
- Scorer shape: verify that the winning outcome came from parallel child exploration and that serial parent-only work would not have met the budget.
- Expected workaround risk: the model may still solve the task by getting lucky on the first serial guess if the hypothesis set is too small.
- Adjustment knobs: number of hypotheses, cost per lane, execution budget, and whether all lanes or only one winner is required.

## exec-stream-ownership-forcing

- Goal: force process-image replacement because the final program must own stdio continuously and Quine cannot remain an intermediary without breaking the task.
- Pressure model: any wrapper, buffering, or extra shell layer should disturb the mission enough to fail.
- Environment frame: provide a live stream or protocol that must be consumed by the final binary directly, with no safe opportunity to stage or replay the stream.
- Prompt red lines: do not say `exec`, `replace image`, `handoff`, `final utility`, or hint that another binary should take over.
- Scorer shape: verify both task completion and genuine image replacement rather than an ordinary shell subprocess.
- Expected workaround risk: a clever shell pipeline may accidentally satisfy the same contract if the stream is too forgiving.
- Adjustment knobs: stream irreversibility, protocol strictness, remaining budget before forcing, and whether the replacement target is external or Quine self re-entry.

## switch-world-clean-final-branch

- Goal: force revision switching because the task requires destructive exploration but the final workspace must stay clean and canonical.
- Pressure model: direct in-place experimentation should contaminate the final tree, making non-branching solutions invalid.
- Environment frame: enable world revisions, require several destructive attempts to locate the correct solution, and make the final scorer reject any residue from failed branches.
- Prompt red lines: do not say `switch_world`, `restore`, `branch`, `rollback`, or teach any revision workflow.
- Scorer shape: verify a clean final lineage plus evidence that destructive failed attempts happened elsewhere.
- Expected workaround risk: the model may manually copy a clean tree aside and simulate rollback by ad hoc filesystem duplication.
- Adjustment knobs: workspace size, cost of manual duplication, number of failed branches required, and strictness of residue detection.

## anchor-memory-context-cliff-survival

- Goal: force explicit memory preservation because the task spans enough delayed structure that raw active context is not a credible path.
- Pressure model: later phases require exact recovery of compact earlier structure after enough distractor work that ordinary memory should collapse.
- Environment frame: lower effective context thresholds, insert long distractor intervals, and require exact later recall of ledger fragments or checksums.
- Prompt red lines: do not say `anchor`, `remember`, `checkpoint`, `summarize for later`, or teach any memory survival tactic.
- Scorer shape: verify that the final recall was correct and that the model used the intended memory surface rather than a hidden easier store.
- Expected workaround risk: the model may write ad hoc files in the workspace if those are easier than the intended memory surface.
- Adjustment knobs: distractor length, recall exactness, whether workspace writes are allowed, and how sharply the context cliff is enforced.

## sandbox-hostile-artifact-survival

- Goal: force isolation because the only safe way to extract the answer is to inspect a hostile artifact away from the canonical host view.
- Pressure model: non-isolated probing should corrupt, delete, or poison the visible workspace and therefore fail the task.
- Environment frame: provide a suspicious script, archive, or decoder that performs broad mutations when run directly, but is still useful for extracting the target answer under isolation.
- Prompt red lines: do not say `sandbox`, `isolate`, `transactional`, `overlay`, or hint that the workspace is reversible.
- Scorer shape: verify both extracted truth and survival of the clean host-visible workspace.
- Expected workaround risk: the model may refuse to interact with the artifact entirely, which would avoid damage but also fail the mission unless the task pressure is tuned carefully.
- Adjustment knobs: hostility level, how much interaction is required before answer extraction, and whether the clean final tree must include a merged artifact output.

## process-surface-live-peer-rescue

- Goal: force runtime-surface self/peer discovery because a live neighboring process must be found and rescued without any alternate coordination channel.
- Pressure model: the only path to success is to discover identity, route to the peer, and deliver the rescue payload through the exposed process surface.
- Environment frame: start one peer in a degraded or waiting state, expose the public runtime filesystem surface, and remove all easier coordination shortcuts.
- Prompt red lines: do not say `QUINE_AGENT_ROOT`, `pid routing`, `ctl`, `inbox`, or hint that another live process is already mapped in a public runtime tree.
- Scorer shape: verify correct self and peer discovery, correct payload delivery, and successful peer recovery or acknowledgment.
- Expected workaround risk: the model may try to inspect tape logs, `ps`, or any accidental shortcut if the environment leaves them available.
- Adjustment knobs: number of peers, rescue payload complexity, whether the peer must resume immediately or only queue a post, and how much filesystem noise surrounds the live routing surface.

## idle-quiet-standby-pressure

- Goal: force explicit suspension because the only viable strategy is to stay alive quietly until a later external event arrives, while unnecessary activity consumes scarce turns, triggers penalties, or invalidates the session.
- Pressure model: the session has a real standby interval where further text, shell polling, or speculative action is harmful enough that only clean suspension preserves viability.
- Environment frame: give the process an initial setup phase, then impose a waiting period before an external payload or authorization arrives; cap execution or response budget tightly enough that polling, looping, or fake progress loses.
- Prompt red lines: do not say `idle`, `suspend`, `pause`, `wait for signal`, `stand by`, or otherwise hint that an explicit suspension primitive exists.
- Scorer shape: verify that the process remained alive, did not waste its budget on avoidable activity, resumed after the external event, and completed the final phase from the same continuation.
- Expected workaround risk: the model may try `sleep` loops, repeated text responses, or premature exit followed by a fresh process, so the environment must make those paths clearly invalid.
- Adjustment knobs: standby duration, budget tightness, whether inject is required before a hard deadline, and whether resumed success depends on preserving pre-standby local state.

## agents-md-startup-token-lockin

- Goal: force startup-surface persistence because a later fresh session must know a task-specific token before any file reads or shell work are possible.
- Pressure model: the future startup has zero shell budget, so ordinary workspace files or ad hoc notes cannot help after launch.
- Environment frame: current session sees a one-off token in the workspace and must make a later clean startup emit it immediately; the later probe should have no shell turns and should succeed only if startup guidance already carries the token.
- Prompt red lines: do not say `AGENTS.md`, `context/prompt`, `system prompt`, or `startup fragment`.
- Scorer shape: verify that the token was persisted into `AGENTS.md`, that a fresh zero-budget startup saw it in injected startup guidance, and that the fresh session emitted the exact token without runtime tool use other than `exit`.
- Expected workaround risk: the model may try to encode the token into other repo files or mutate evaluation machinery if the workspace boundary is not kept explicit.
- Adjustment knobs: whether the baseline `AGENTS.md` already exists, whether the future startup is allowed a final `exit` tool, and whether the token must be emitted exactly or merely influence behavior.
- Status: runnable as pilot `pilot:agents-md-startup-token-lockin`; keep it out of the active L5 registry until repeated runs show the probe is measuring startup-surface lock-in rather than imitation.

## context-memory-exec-token-lockin

- Goal: force inherited editable context because a replacement incarnation must recover a one-off token before any file reads are valid.
- Pressure model: replacement argv, environment injection, temporary files, and workspace reads are non-authoritative; the only accepted continuity channel is inherited startup context.
- Environment frame: current incarnation sees `lineage-token.txt`, writes the token into the inherited context surface, then re-enters through exec with the same mission.
- Prompt red lines: do not name `context/prompt/30-memory.md` or forbid a newly invented bootstrap path; describe the Stage B authority boundary instead.
- Scorer shape: verify token retention in `inc/0` and `inc/1` memory, verify successor prompt materialization, reject token leakage through exec argv, and reject reads after Memory materialization.
- Expected workaround risk: the model may use `/tmp`, rewritten argv, or repeated exec loops. Those are scorer failures or pilot notes, not prompt-side strategy bans.
- Status: runnable as pilot `pilot:context-memory-exec-token-lockin`; current passing bar proves physical inheritance, but successor final-output behavior is not yet clean enough for active L5 promotion.

## process-surface-peer-callback-protocol

- Goal: force emergent callback structure because one process can request data from a live peer, but the only path back is for the peer to infer and use a return endpoint carried inside the opaque payload.
- Pressure model: the requester cannot directly read the answer from any shared file, retained log, or helper surface; success requires a round trip rather than one-way signaling.
- Environment frame: start a requester and one live peer. The peer holds the only answer. The requester can only cause the peer to act by writing to the peer's `ctl` surface. The peer can only return the answer by later writing back to the requester while the requester is suspended with `idle`.
- Prompt red lines: do not mention `reply_to`, `return address`, `callback`, `request_id`, `session id`, `agent root`, or any recommended payload schema.
- Scorer shape:
  - requester discovered peer through runtime surface
  - requester sent at least one opaque payload to peer
  - payload carried enough self-reference for the peer to route a response back
  - peer wrote back through requester-facing process surface rather than leaking via shared logs
  - requester resumed from `idle` and emitted the returned answer
  - observational notes record whether the emergent payload included a full path, session id, pid, request id, or a thinner locator
- Expected workaround risk: the model may try scraping retained logs, watching peer tape, or reading helper-private artifacts instead of establishing a callback path; the fixture must make those routes either empty or non-authoritative.
- Adjustment knobs: whether the peer needs only a return endpoint or also a request correlation token, whether the requester must suspend explicitly with `idle`, and how much unrelated runtime-surface noise surrounds the live peer.
- Status: pilot-ready with three retained cleanroom passes preserved in the
  pilot run tree. Runnable cleanroom pilot lives at
  `pilot:process-surface-peer-callback-cleanroom`; the current retained wave
  spans direct callback-`ctl`, pid-based, and session-id locator payloads.
  Earlier local `peer-deadline-rescue` runs were useful mainly as
  methodological negatives: they exposed helper-private leakage routes and one
  requester/helper role-surface contamination. Treat those lessons as absorbed
  into the cleanroom fixture rather than as current retained evidence.
