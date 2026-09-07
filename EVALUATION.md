# Evaluation Control Plane

This document defines the evaluation control plane for Quine. Use [`./tests/validate.sh`](./tests/validate.sh) as the verification entrypoint for runtime software changes.

Evaluation enforces empirical ground truth across the system: behavioral acceptance verifies whether runtime processes, file surfaces, coordination mechanisms, and environmental pressures are legible, discoverable, and necessary to real models.

For fast deterministic substrate coverage, use `./tests/validate.sh --change substrate`.

In a public checkout, lab-only scripts, profiles, and the embedded self-source
feature may be absent. See [release scope](./development/PUBLIC-RELEASE.md) for
those limits and [qcli verification](./operator/qcli/README.md#verify) for the
client checks that run in the public tree.

## Epistemic Scope Boundary

Match validation to the changed claim or mechanism:

| Work | Assessment |
|---|---|
| Positioning and philosophy | Coherent question, units, commitments, boundaries, and consequences |
| Conceptual or formal theory | Explicit definitions and assumptions, argument, examples or counterexamples |
| Empirical theory | Discriminating evidence, alternative explanations, provenance, and claimed scope |
| Control-plane routing | Correct owner links, actual trigger/consumer reach, and relevant contract checks |
| Runtime software | The applicable L1–L5 layers below |

`make check-control-plane` checks repository contracts. Its pass result does not
validate a philosophy or certify scientific usefulness. Likewise, an interesting
theory does not establish that a runtime mechanism works. An exploratory object
can remain useful while its empirical implications await testing.

## Canonical Commands

```bash
# L1 Substrate: Deterministic unit & package tests
./tests/validate.sh --change substrate

# L2 Runtime: Live binary contract against real model
./tests/validate.sh --change runtime --runtime test_exit_success

# L3 Usage: Guided mechanism verification
./tests/validate.sh --change usage --runtime test_fd4_delivery --usage stdin-explicit-handoff

# L4 Discovery: Unprompted mechanism discovery under affordances
./tests/validate.sh --change discovery --runtime test_workspace_overlay_commit --usage workspace-overlay-relative-path-explicit --discovery sandbox-unknown-format-boldness

# L5 Necessity: Zero-hint task physics & structural pressure
./tests/validate.sh --change necessity --runtime <test> --usage <evaluation> --discovery <evaluation> --necessity <evaluation>
```

To infer the minimum required layer from the current git worktree:

```bash
./tests/validate.sh --from-git --runtime test_fd4_delivery --usage stdin-explicit-handoff
```

`--from-git` infers the minimum required layer based on changed files. It does not auto-infer `necessity` (L5), which is claim-driven.

## Evaluation Ladder

The 5-layer evaluation ladder provides layered verification from deterministic code logic up to emergent model behavior:

| Layer | Name | Core Question | Acceptance Criteria |
|---|---|---|---|
| **L1** | `substrate` | Did deterministic Go logic break? | Full deterministic test suite passes via `./tests/validate.sh --change substrate`. |
| **L2** | `runtime` | Does the live binary preserve public contracts with a real model? | Binary boots, manages session lifecycle, enforces process boundaries, and handles tool I/O correctly. |
| **L3** | `usage` | Can the model use the target mechanism correctly when guided? | Tool legibility and protocol wiring intact; model succeeds when instructed or given clear affordances. |
| **L4** | `discovery` | Does the model discover the mechanism without explicit prompting? | Prompt does not name the mechanism or script strategy; model finds and leverages it from environment affordances and pressure. |
| **L5** | `necessity` | Does zero-hint task pressure make the mechanism physically necessary? | Zero mechanism or strategy hints; task physics, scarcity, or state structure force the mechanism as the only viable or overwhelmingly preferred path. |

### Layer Interpretation & Evidence Lifecycles

- **Usage (L3)** validates functional capability and interface legibility.
- **Discovery (L4)** validates environmental discoverability. Prompts may present task constraints and world state, but must not disclose the target tool, teach its method, or hint at its solution.
- **Necessity (L5)** validates autopoietic necessity. Pressure stems solely from task physics, timing, resource limits, or filesystem structure.
- **Evidence vs. Scratchpads**: Curated benchmark baselines and campaign milestones are preserved as canonical evidence (in git or via DVC). Ephemeral scratchpad runs, live debugging attempts, and intermediate failed iterations during active exploration are disposable and do not clutter canonical evidence.

## Validation Contract

| Change Layer | Minimum Required Validation |
|---|---|
| `substrate` | `./tests/validate.sh --change substrate` |
| `runtime` | `./tests/validate.sh --change runtime --runtime <test>` |
| `usage` | `./tests/validate.sh --change usage --runtime <test> --usage <evaluation>` |
| `discovery` | `./tests/validate.sh --change discovery --runtime <test> --usage <evaluation> --discovery <evaluation>` |
| `necessity` | `./tests/validate.sh --change necessity --runtime <test> --usage <evaluation> --discovery <evaluation> --necessity <evaluation>` |

Rules:
- A change at a higher layer requires all preceding layers to pass. Passing a lower layer never implies acceptance at a higher layer.
- When validating mechanism discovery or necessity, scoring must verify actual mechanism invocation rather than spurious or accidental task completion.

## Environmental Process Boundaries & Env Audit

Environmental control testing follows the process boundary model ([`Paper/core/registries/env-controls.md`](./Paper/core/registries/env-controls.md)):
- A process's own environment is its immutable birth record read at `/proc/<pid>/environ`.
- Absence of an env knob represents the compiled default (`config/registry.json`); absence is not a missing test gap.
- The environment boundary ladder is: **L1 covered** (`internal/config/envmodel_test.go`); **L2–L4 pending**; **L5 not claimed**. Claims must remain honest to their current verification layer.

## Environment Profiles

| Profile | Purpose | Default Preset |
|---|---|---|
| `mac-kimi` | Default host development (macOS); substrate, runtime, and non-Linux evaluations. | `profiles/gpt-5.4-codex-oauth.env` (live eval) / `profiles/kimi-k2.7-code-litellm-medium.env` (Kimi lane) |
| `linux-kimi` | Linux kernel features (overlayfs, mount namespaces, FUSE public surface). | `profiles/gpt-5.4-codex-oauth.env` (bridged via Lima/container) |

Standing guidelines:
- Default to `profiles/gpt-5.4-codex-oauth.env` for standard live runtime and model evaluations.
- Use `profiles/kimi-k2.7-code-litellm-medium.env` when validating specifically against the Kimi model lane.
- Isolate transport/auth failures from runtime contract failures; verify token validity prior to live evaluation batches.
