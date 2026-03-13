# Testing Quine

Quine has three public test surfaces. They answer different questions, so do not treat them as substitutes for one another.

| Surface | Command | What it verifies |
|---------|---------|------------------|
| Go tests | `make test` | Deterministic package-level correctness |
| End-to-end runtime tests | `./tests/e2e.sh` | The built binary behaves like a POSIX process against a real model API |
| LLM behavior tests | `./test/llm-behavior/run.sh` | A model can infer and use Quine's tool surface from prompt plus schema |

## Quick Commands

```sh
# Deterministic tests
make test

# Build the binary used by live API-backed checks
go build -o /tmp/quine ./cmd/quine/

# Run one end-to-end runtime test
source .env
./tests/e2e.sh test_exit_success

# Run one behavior scenario
source .env
./test/llm-behavior/run.sh jobs
```

## 1. Go Tests

Use these first. They are local, deterministic, and do not hit a model API.

```sh
make test
```

This covers the runtime implementation under [`cmd/quine/`](./cmd/quine/) and [`internal/`](./internal/).

## 2. End-to-End Runtime Tests

[`tests/e2e.sh`](./tests/e2e.sh) runs the real binary against a live API and checks observable process behavior:

- exit codes
- `stdout`
- `stderr`
- tape files

Typical flow:

```sh
go build -o /tmp/quine ./cmd/quine/
source .env
./tests/e2e.sh
./tests/e2e.sh test_fd3_delivery
```

Use this after changing runtime semantics, I/O behavior, shell execution, or other process-facing contracts.

## 3. LLM Behavior Tests

[`test/llm-behavior/run.sh`](./test/llm-behavior/run.sh) checks whether a model can discover how to use Quine correctly from the system prompt and tool schema alone.

Typical flow:

```sh
go build -o /tmp/quine ./cmd/quine/
source .env
./test/llm-behavior/run.sh jobs
./test/llm-behavior/run.sh all
```

Use this after changing prompt wording, tool schemas, or output conventions.

Run artifacts are written under [`test/llm-behavior/runs/`](./test/llm-behavior/runs/). Review `tape.log` and `score.txt` when a scenario fails.

## Notes

- The last two surfaces make real API calls and can consume tokens.
- `make test` does not replace the live checks.
- The live checks do not replace deterministic Go tests.
- Start with the smallest surface that can falsify the change you just made.
