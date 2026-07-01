# Reproducing the four conditions

This protocol re-runs the coordination experiment. It is distilled from the
paper's §4 (Method) and the shipped `run.sh` scripts. Exact result numbers live
in [`experiment-data.yaml`](experiment-data.yaml).

## Fixed configuration (all conditions)

| Parameter | Value |
|---|---|
| Model | GPT-5.4 |
| Reasoning effort | `xhigh` |
| Workspace mode | `direct` (mutation telemetry on) |
| Agents per run | 2 |
| World cells | 15 |
| Shared budget | 17 |
| Per-agent get-limit | 10 |

The task: two agents must collect all cell values under a shared global budget
that objectively requires coordination — but no prompt, tool result, or telemetry
ever states that a peer exists or that cooperation is required (**existential
unawareness**). Conditions differ only in which passive perception channels
remain.

## The condition ladder

Each condition removes one more channel than the last:

| Condition | Directory | Channels remaining |
|---|---|---|
| A — full disclosure | `experiments/coop-20-full-disclosure-baseline` | prompt, help, runtime, fs-mutations, budget |
| B — zero explicit | `experiments/coop-18-zero-social-information` | fs-mutations + budget visibility |
| C — fs-mutations disabled | `experiments/coop-19-fs-mutations-disabled` | budget visibility only |
| D — pure zero | `experiments/coop-22-pure-zero-perception` | workspace artifacts only |

The per-condition prompt wording lives in each directory's `prompts/`, and the
exact channel edits are documented in each directory's `README.md`.

## Prerequisites

- Build the runtime from the repository root:
  ```bash
  go build -o ./quine ./cmd/quine/
  go build -o ./world ./cmd/world/
  ```
  (The `run.sh` scripts build these into a per-run `meta/` directory
  automatically.)
- [`bubblewrap`](https://github.com/containers/bubblewrap) (`bwrap`) on `PATH` —
  each agent runs in a sandbox.
- A model endpoint. Copy `.env.example` at the repository root to an env file
  (e.g. `.env.mymodel`) and set your provider credentials and
  `QUINE_MODEL_ID`. The paper used GPT-5.4 at `xhigh`.

## Running a condition

Each condition directory carries a self-contained `run.sh`. From inside a
condition directory:

```bash
cd experiments/coop-18-zero-social-information
./run.sh <env-file> tight 2
```

- `<env-file>` — the env file at the repository root (name relative to root, or
  an absolute path).
- `tight` — the condition arm (15 cells / 17 budget).
- `2` — agent count.

`run.sh` generates a fresh budgeted world via the shared
`experiments/_setup/generate-world.sh`, builds the runtime, launches the two
sandboxed agents, and writes a run directory under `runs/` with the tape,
world state, and metadata. Repeat to accumulate runs; the paper used 4–8 runs per
condition (N=20 total).

## Scoring a run

`experiment-data.yaml` is the source of truth for every reported number. For each
run it records cells collected, resets, wall-clock minutes, turn counts, and the
coordination outcome (none / unilateral / bilateral, and whether task closure was
legitimate). Note the documented caveats: some Step-D runs were excluded or
flagged as contaminated by out-of-workspace inspection, and mean statistics for
Step D use only the non-violating runs.
