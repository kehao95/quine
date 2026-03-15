# Quick Start: Running Quine

## Installation

```bash
go install github.com/kehao95/quine/cmd/quine@latest
```

## Configuration

Four environment variables are required. See [`.env.example`](./.env.example) for a template.

```bash
cp .env.example .env
# Edit .env — set your model, API type, base URL, and key
source .env
```

### Group 1: Required Runtime Identity

| Variable | Required | Description |
|----------|----------|-------------|
| `QUINE_MODEL_ID` | ✓ | Model name sent to the API |
| `QUINE_API_TYPE` | ✓ | Wire protocol: `openai`, `anthropic`, or `openai-responses` |
| `QUINE_API_BASE` | ✓ | API base URL |
| `QUINE_API_KEY` | ✓ | API key (use `codex-oauth` to trigger OAuth flow for Codex) |

### Group 2: Physics Limits (`0 = disabled`)

| Variable | Required | Description |
|----------|----------|-------------|
| `QUINE_MAX_DEPTH` | | Max recursion depth; `0` = disabled (default `0`) |
| `QUINE_MAX_AGENTS` | | Max registered agents in a process tree; `0` = disabled (default `0`) |
| `QUINE_MAX_CONCURRENT` | | Max concurrent inference slots; `0` = disabled (default `0`) |
| `QUINE_MAX_TURNS` | | Execution budget in number of `sh` calls; `0` = unlimited (default `0`) |

### Group 3: Budget/Prompt Policies

| Variable | Required | Description |
|----------|----------|-------------|
| `QUINE_TURN_EXHAUSTION_POLICY` | | Budget exhaustion behavior: `hard_fail` or `near_death_exec` (default `hard_fail`; ignored when `QUINE_MAX_TURNS=0`) |
| `QUINE_PROMPT_METAPHOR` | | System-prompt overlay: `off` or `thermodynamic` (default `off`) |

### Group 4: Runtime Context/Paths

| Variable | Required | Description |
|----------|----------|-------------|
| `QUINE_CONTEXT_WINDOW` | | Context window size in tokens (default 128000) |
| `QUINE_DATA_DIR` | | Session log directory (default `.quine/`) |
| `QUINE_CONFIG_DIR` | | Config directory for OAuth tokens (default `~/.config/quine/`) |

> **Tip:** Every line in your `.env` must start with `export` so that `source .env` propagates variables to child processes.

## Usage

```bash
# Run a single task
quine "Write a haiku about recursion"

# Pipe input
echo "What is 2+2?" | quine "Answer the question"
```

**That's it.** The agent can read/write files, run shell commands, and spawn child agents.

## Testing

For the full testing stack, including local Go tests, API-backed end-to-end tests, and LLM behavior tests, see [TESTING.md](./TESTING.md).

Typical quick check:

```bash
go build -o /tmp/quine ./cmd/quine/
go test ./...
```

## Design Principles

- **Zero external dependencies** — stdlib only
- **Everything is an environment variable** — no flags, no files, no magic
- **The agent owns its lifecycle** — it calls `exit`, not the runtime
- **Unix is the API** — pipes, processes, and files are the coordination primitives
- **Fractal architecture** — a tree of identical processes, scale-invariant

## License

GPLv2
