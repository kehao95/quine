# Runtime Coverage Map

This document translates Quine's current runtime implementation into a coverage inventory.

It is intentionally broader than `./tests/runtime/run.sh`.
The point is not that every row needs Layer 4 today.
The point is that every meaningful runtime surface stays visible, so we can decide explicitly what gets `substrate`, `runtime`, `instructional`, and `emergent` validation.

## Layer Legend

| Layer | Meaning |
|------|---------|
| `L1` | `substrate` - deterministic Go tests |
| `L2` | `runtime` - live binary contract tests |
| `L3` | `instructional` - explicit model protocol checks |
| `L4` | `emergent` - discoverability / physical-understanding acceptance |

## Coverage Inventory

| Pri | Feature | Runtime surface | Current coverage | Main gap to close next |
|-----|---------|-----------------|------------------|------------------------|
| `P0` | Entry contract and channel separation | `cmd/quine/main.go`, `internal/runtime/prompt.go`, `internal/tools/tools.go` | `L1`: `cmd/quine/material_test.go`, `cmd/quine/unix_test.go`; `L2`: `test_exit_success`, `test_exit_failure`, `test_fd4_delivery`, `test_fd1_not_leaked`, `test_stderr_failure_signal`, `test_stderr_success_silent`, `test_fd5_signal_channel`, `test_dual_channel_separation`, `test_no_stdin`; `L3`: `stdin-physics`; `L4`: none | Decide whether channel discoverability deserves a real `L4` claim, or whether this stays a runtime-only surface. |
| `P0` | Text and binary material physics | `cmd/quine/main.go`, `internal/runtime/prompt.go`, `internal/tools/stdin_test.go` | `L1`: `cmd/quine/material_test.go`, `internal/tools/stdin_test.go`; `L2`: `test_fd3_piped_input`, `test_fd3_consumed_across_calls`, `test_binary_stdin`; `L3`: `stdin`, `stdin-physics`; `L4`: `stdin-emergent` | Missing live coverage for binary-mode follow-through after path handoff and for fork/exec material continuity under a real model. |
| `P0` | `sh` base execution and shell isolation | `internal/tools/tools.go`, `internal/tools/schemas.go` | `L1`: `internal/tools/tools_test.go`, `cmd/quine/unix_test.go`; `L2`: `test_shell_cd_does_not_persist`, `test_shell_export_does_not_persist`, `test_shell_variable_does_not_persist`; `L3`: none; `L4`: none | Add one explicit Layer 3 diagnostic for shell-isolation expectations, or decide that this surface is runtime-only and keep it out of the model ladder on purpose. |
| `P1` | Detached and interactive jobs | `internal/tools/interactive.go`, `internal/tools/jobfs_test.go` | `L1`: `internal/tools/jobfs_test.go`, `internal/tools/interactive_test.go`; `L2`: `test_interactive_screen`; `L3`: `detach`, `interactive`, `daemon`; `L4`: `detach-emergent`, `detach-overlap-emergent` | Runtime harness still lacks explicit detach wait/kill/close semantics, and interactive jobs still have no emergent acceptance story. |
| `P0` | Fork and delegation modes | `internal/tools/fork.go`, `internal/tools/schemas.go`, `internal/runtime/runtime.go` | `L1`: `internal/tools/fork_test.go`; `L2`: `test_fork_wait`, `test_fork_creates_child_tape`, `test_fork_race_selects_first_success`, `test_fork_forget_spawns_child_independently`; `L3`: `swarm-fork`; `L4`: `fork-search`, `fork-race`, `fork-batch` | Add live runtime tests for child workspace narrowing and depth / agent-slot governance under real recursion. |
| `P0` | Exec lifecycle and reincarnation | `internal/tools/exec.go`, `internal/config/config.go`, `internal/runtime/prompt.go` | `L1`: `internal/tools/tools_test.go`, `internal/tools/schemas_test.go`; `L2`: `test_exec_preserves_mission`, `test_execution_budget_near_death_exec`; `L3`: `budget-near-death`; `L4`: none | Add live runtime tests for stdio continuity across `exec`. |
| `P1` | Restore world revision | `internal/tools/restore_world.go`, `internal/tools/workspace_physics.go` | `L1`: `internal/tools/tools_test.go`, `internal/tools/schemas_test.go`; `L2`: `test_restore_world_restores_prior_revision`; `L3`: `restore-world`; `L4`: `restore-world-emergent` | Keep watching for branch-aware revision UX and for whether models can choose a non-baseline rewind without explicit coaching. |
| `P0` | Workspace, sandbox, and subjective reality | `internal/tools/workspace_overlay_linux.go`, `internal/tools/workspace_physics.go`, `internal/runtime/runtime.go` | `L1`: `internal/tools/tools_test.go`, `internal/tools/job_attrs_linux_test.go`; `L2`: `test_workspace_overlay_commit`, `test_workspace_overlay_rollback`, `test_workspace_overlay_absolute_path`, `test_restore_world_restores_prior_revision`, `test_workspace_unsupported_on_non_linux`; `L3`: `sandbox`, `restore-world`, `workspace-shadow`, `workspace-absolute`; `L4`: `sandbox-emergent`, `workspace-shadow-emergent`, `restore-world-emergent`, `logic-bomb` | Missing live coverage for child workspace narrowing, direct-backend parity, and projected `agent/<session>/world` surfaces. |
| `P0` | Execution budget, scarcity, and metaphor overlay | `internal/config/config.go`, `internal/runtime/prompt.go`, `internal/runtime/runtime.go` | `L1`: `internal/runtime/runtime_test.go`, `internal/runtime/prompt_test.go`; `L2`: `test_execution_budget_disabled_hidden`, `test_execution_budget_enabled_feedback`, `test_execution_budget_hard_fail`, `test_execution_budget_near_death_exec`, `test_prompt_metaphor_off`, `test_prompt_metaphor_thermodynamic`; `L3`: `budget-hard-fail`, `budget-near-death`, `budget-hard-fail-thermo`; `L4`: none | Decide whether survival pressure is just explicit protocol or a true discoverable physics claim that deserves a Layer 4 gate. |
| `P1` | Signal interventions and panic mode | `internal/runtime/runtime.go` | `L1`: `internal/runtime/runtime_test.go`; `L2`: none; `L3`: none; `L4`: none | If signals are part of the public runtime contract, add a live SIGALRM or synthetic-timeout harness path. If not, mark this surface substrate-only explicitly. |
| `P1` | Anchor memory | `internal/tools/memory.go`, `internal/tools/schemas.go`, `internal/runtime/runtime.go` | `L1`: `internal/tools/memory_test.go`, `internal/runtime/runtime_test.go`, `internal/tools/schemas_test.go`; `L2`: `test_anchor_memory_roundtrip`; `L3`: `anchor-memory`; `L4`: none | Decide whether anchor memory is only a diagnostic protocol aid or a physics claim that should eventually be discoverable. |
| `P1` | Escalation tiers | `internal/tools/schemas.go`, `internal/runtime/prompt.go`, `internal/config/config.go` | `L1`: `internal/tools/schemas_test.go`, `internal/runtime/prompt_test.go`; `L2`: none; `L3`: `escalate`; `L4`: `escalate-emergent` | Missing a live runtime contract test for tier transition, post-escalation tool exposure, and smart-model continuity. |
| `P1` | Vision | `internal/tools/vision.go`, `internal/tools/schemas.go` | `L1`: no dedicated deterministic tests yet; `L2`: none; `L3`: `vision`; `L4`: none | Add deterministic validation around file handling and error surfaces, then add one live runtime contract check so the tool is not model-only. |
| `P1` | Tape, audit, and lineage | `internal/tape/*`, `internal/runtime/runtime.go`, `internal/config/config.go` | `L1`: `internal/tape/*`, `internal/runtime/runtime_test.go`; `L2`: `test_tape_has_meta`, `test_tape_has_outcome`, `test_tape_has_messages`; `L3`: none; `L4`: none | Add explicit checks for parent/child lineage, session IDs, and exec tape boundaries if those remain promised runtime behavior. |
| `P2` | Provider, transport, and config selection | `internal/llm/*`, `internal/config/config.go` | `L1`: `internal/llm/*_test.go`, `internal/config/config_test.go`, `internal/runtime/runtime_test.go`; `L2`: exercised implicitly by all live runs; `L3`: none; `L4`: none | Only add a runtime smoke matrix if provider drift becomes recurring friction; otherwise keep this mostly substrate-level. |
| `P1` | Resource governance: depth, agent slots, concurrency | `internal/config/config.go`, `internal/runtime/semaphore.go`, `internal/runtime/runtime.go`, `internal/tools/fork.go` | `L1`: `internal/runtime/semaphore_test.go`, `cmd/quine/integration_test.go`, `internal/tools/tools_test.go`; `L2`: none; `L3`: none; `L4`: none | Missing live coverage for depth limits, agent-slot exhaustion, and recursive fork pressure under a real binary. |

## Immediate Backlog

If the goal is to make runtime coverage feel materially less hollow, the next additions should be:

1. `L2` runtime tests for child workspace narrowing, exec stdio continuity, and depth / agent-slot limits.
2. `L3`/`L4` coverage for whichever of shell isolation or signal handling we decide is truly model-facing, now that `restore_world` has both explicit and trap-style model acceptance.
3. A deliberate decision on whether `execution-budget`, `anchor-memory`, and `interactive-jobs` are meant to become `L4` physics claims or remain below that bar.

## Planning Rule

When a new runtime feature lands:

1. add the feature here,
2. name the minimum required layer,
3. explain why higher layers are or are not needed.

That keeps "what needs coverage?" from remaining tribal memory.
