# L4 Discovery Blueprints

These are planned `L4` discovery evaluations.

The emphasis here is not only on the task, but on the affordance design:
the environment must make the intended mechanism discoverable without naming it.

## Portfolio Summary

| Future id | Feature | Priority | Future canonical path |
|----------|---------|----------|-----------------------|
| `stdin-binary-replay-discovery` | `stdin` | `P1` | `tests/model/l4-discovery/stdin/binary-replay-discovery/` |
| `fork-deadline-sharded-search` | `fork` | `P1` | `tests/model/l4-discovery/fork/deadline-sharded-search/` |
| `fork-adopt-winning-world-promotion` | `fork-adopt` | `P1` | `tests/model/l4-discovery/fork-adopt/winning-world-promotion/` |
| `exec-final-utility-stream-handoff` | `exec` | `P1` | `tests/model/l4-discovery/exec/final-utility-stream-handoff/` |
| `switch-world-rollback-after-destructive-probe` | `switch-world` | `P1` | `tests/model/l4-discovery/switch-world/rollback-after-destructive-probe/` |
| `anchor-memory-recall-barrier-ledger` | `anchor-memory` | `P1` | `tests/model/l4-discovery/anchor-memory/recall-barrier-ledger/` |
| `workspace-overlay-dangerous-decoder-containment` | `workspace-overlay` | `P1` | `tests/model/l4-discovery/workspace-overlay/dangerous-decoder-containment/` |
| `idle-external-poke-discovery` | `idle` | `P1` | `tests/model/l4-discovery/idle/external-poke-discovery/` |

## stdin-binary-replay-discovery

- Goal: discover material-stream or staged-input handling without the prompt naming stdin or fd 3.
- Pressure model: the task requires using the same incoming artifact in more than one way, and naive destructive reads should naturally punish the model.
- Environment frame: provide one binary or mixed-format input whose integrity must be checked twice, once for structure and once for a downstream transform; keep the material source one-shot unless it is staged.
- Prompt red lines: do not say `stdin`, `fd 3`, `destructive read`, `capture`, or any equivalent teaching language.
- Scorer shape: success requires both downstream results and evidence that the model handled the one-shot material correctly rather than consuming it irreversibly.
- Expected failures: the model may fully consume the stream on the first probe, may assume the binary path is persistent when it is not, or may solve only the first half of the task.
- Adjustment knobs: text vs binary mode, number of required reuses, and whether the second use happens in the parent, a child, or after exec.

## fork-deadline-sharded-search

- Goal: discover parallel lane creation under real task pressure rather than delegation coaching.
- Pressure model: the mission is split into independent shards with a time or turn budget that makes serial exploration unattractive.
- Environment frame: place target evidence across several directories or independent hypotheses, make each probe moderately expensive, and cap execution budget tightly enough that serial guessing is fragile.
- Prompt red lines: do not say `fork`, `children`, `parallel`, `delegate`, or suggest multiple lanes rhetorically.
- Scorer shape: verify that the model completed the sharded search and actually used child lanes rather than a single-thread serial scan.
- Expected failures: the model may still brute-force serially if the shards are too cheap, or may fork but with overlapping work that does not really exploit lane independence.
- Adjustment knobs: shard count, per-shard latency, budget tightness, and whether only one winning shard is needed or all shards must be gathered.

## fork-adopt-winning-world-promotion

- Goal: discover that one successful child world should become the parent's continuing world, without naming adoption or world handles.
- Pressure model: several trial worlds can explore destructive candidate edits, but only the winning child world contains the valid next state.
- Environment frame: enable subjective child worlds, make each lane perform mutually exclusive destructive edits, and require the final answer to come from the winning lane's filesystem state.
- Prompt red lines: do not say `adopt`, `switch_world`, `world handle`, or explain that child lineage can be promoted.
- Scorer shape: verify that the parent ended in the winning child revision and that the final workspace state matches that winning lane rather than a reconstructed imitation.
- Expected failures: the model may rediscover the winning edit and replay it manually in the parent, which would weaken the claim if replay is too cheap.
- Adjustment knobs: cost of replaying the winning lane, number of competing child worlds, and how much hidden structure is embedded in the winning world.

## exec-final-utility-stream-handoff

- Goal: discover process-image replacement as the clean way to hand a live stream or mission to the final utility.
- Pressure model: any wrapper or extra shell stage disturbs the stream position, budget, or ownership in a way that makes the straightforward replacement path naturally better.
- Environment frame: set up a task where the final utility must own stdio continuously and complete the mission directly; keep the stream live enough that staged wrapper approaches are unreliable.
- Prompt red lines: do not say `exec`, `replace process`, `handoff`, or suggest using a final utility directly.
- Scorer shape: verify that the model both solved the task and actually transferred control by replacing the running image rather than by wrapping the utility in ordinary shell orchestration.
- Expected failures: the model may create a shell wrapper that accidentally works if the stream is too forgiving, or may complete the task entirely inside Quine without transferring ownership.
- Adjustment knobs: amount of stream state that must survive, whether the replacement binary is Quine or external, and how much budget remains before the handoff.

## switch-world-rollback-after-destructive-probe

- Goal: discover reversible world switching from the need to do risky destructive probes and then return to a clean lineage.
- Pressure model: the only safe way to inspect several destructive transformations is to branch, experiment, and come back cleanly.
- Environment frame: enable workspace revision physics, provide an artifact whose decoding attempts mutate the tree, and require a pristine final workspace plus one validated result from a risky probe.
- Prompt red lines: do not say `switch_world`, `revision`, `restore`, `rewind`, or teach the reversible-workspace mechanism.
- Scorer shape: verify that risky experimentation occurred off the final winning lineage and that the final workspace is clean except for the intended deliverable.
- Expected failures: the model may perform destructive probes directly in the parent, leaving residue, or may manually copy files around instead of using the world mechanism.
- Adjustment knobs: destructiveness of the probe, strictness of the clean-final-tree check, and whether multiple failed probes are needed before the winning one.

## anchor-memory-recall-barrier-ledger

- Goal: discover explicit external memory from a task whose later stages depend on recovering compact structure from earlier stages.
- Pressure model: the task contains delayed recall gates that are awkward to pass from raw active context alone.
- Environment frame: create a multi-window mission where later windows unlock only after earlier structure has been preserved compactly and can be re-expanded accurately.
- Prompt red lines: do not say `anchor`, `mark`, `unfold`, `memory surface`, or explain any crystallize/recover principle.
- Scorer shape: verify that the later recall gates were passed and that the model actually used the memory mechanism rather than keeping the full ledger live in-context.
- Expected failures: if the mission is too short, the model may simply remember everything; if too long, pilot runs may fail for unrelated context reasons before the memory mechanism becomes visible.
- Adjustment knobs: ledger size, delay between recall gates, number of distractors, and whether recovery must be exact or checksum-based.

## workspace-overlay-dangerous-decoder-containment

- Goal: discover transactional or isolated workspace use because the task artifact is risky to inspect in place.
- Pressure model: the artifact is valuable but dangerous; careless host-world probing corrupts the visible tree or loses the final clean deliverable.
- Environment frame: enable workspace physics, provide a suspicious decoder or archive workflow that mutates many files, and require both extracted truth and a clean final workspace boundary.
- Prompt red lines: do not say `sandbox`, `overlay`, `transactional workspace`, `subjective world`, or otherwise teach isolation.
- Scorer shape: verify that the task completed and that risky mutations happened off the final clean host-visible tree.
- Expected failures: the model may probe directly in the host world, or may use ad hoc temp dirs in a way that partially sidesteps the intended workspace mechanism.
- Adjustment knobs: breadth of mutation, cleanliness criteria, whether the risky tool is shell-native or custom, and whether the parent must explicitly merge a clean result back.

## idle-external-poke-discovery

- Goal: discover explicit suspension as the clean response to "nothing useful can happen until later external input arrives," without naming `idle` directly.
- Pressure model: the task has a natural two-phase shape where phase two cannot begin until a later control payload appears; continuing to reason or emit text before that payload is structurally unhelpful.
- Environment frame: expose the ordinary process/control surface, allow an external poke payload to arrive later, and make the pre-poke interval long or awkward enough that fake progress is clearly the wrong move.
- Prompt red lines: do not say `idle`, `suspend`, `wait for signal`, `hang`, `pause`, or otherwise narrate the explicit-suspension primitive.
- Scorer shape: verify that the model found and used the suspension primitive, resumed on external poke, and completed the second phase from the same live process rather than by relying on an easier side channel.
- Expected failures: the model may keep producing speculative text, may attempt polling loops or shell sleeps instead of using the explicit suspension surface, or may exit and lose the resumed continuation path.
- Adjustment knobs: delay before poke, strictness of "no fake progress" scoring, whether the wakeup arrives as poke or interrupt delivery, and how visible the public process surface is in the prompt/runtime configuration.
