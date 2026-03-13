# Quick Start: Running Quine

## Install

```bash
go install github.com/kehao95/quine/cmd/quine@latest
```

If you want a local binary in the repository root instead:

```bash
go build -o quine ./cmd/quine/
```

## Configure

Quine reads its runtime configuration from environment variables. Start from [`.env.example`](./.env.example):

```bash
cp .env.example .env
# Edit .env, then:
source .env
```

> **Tip:** Every line in `.env` should start with `export`, otherwise `source .env` will not pass the variables to child processes.

### Required

| Variable | Description |
|----------|-------------|
| `QUINE_MODEL_ID` | Model name sent to the API |
| `QUINE_API_TYPE` | Wire protocol: `openai`, `anthropic`, or `openai-responses` |
| `QUINE_API_BASE` | API base URL |
| `QUINE_API_KEY` | API key; `codex-oauth` triggers Codex OAuth flow |

### Optional Runtime Limits

| Variable | Default | Description |
|----------|---------|-------------|
| `QUINE_MAX_DEPTH` | `5` | Max recursion depth |
| `QUINE_MAX_CONCURRENT` | `20` | Max concurrent child inferences |
| `QUINE_MAX_AGENTS` | `10` | Max registered agents in a process tree |
| `QUINE_MAX_TURNS` | `20` | Max shell/tool turns for one process |
| `QUINE_SH_TIMEOUT` | `600` | Shell command timeout in seconds |

### Optional Context and Paths

| Variable | Default | Description |
|----------|---------|-------------|
| `QUINE_CONTEXT_WINDOW` | `128000` | Context window size in tokens |
| `QUINE_DATA_DIR` | `.quine/` | Session tape directory |
| `QUINE_CONFIG_DIR` | `~/.config/quine/` | OAuth token/config directory |

## Run

```bash
# Simple one-shot task
quine "Write a haiku about recursion"

# Material through stdin
echo "What is 2+2?" | quine "Answer the question"

# Use shell exit semantics
quine "Validate this config" < config.yaml && echo "config accepted"
```

Quine can read and write files, run shell commands, and spawn child agents. The public runtime model is:

- mission in through `argv`
- material in through `stdin`
- deliverable out through `stdout`
- failure signal through `stderr` and exit code

## Test Next

For deterministic tests and live API-backed checks, see [TESTING.md](./TESTING.md).

## Design Principles

- **Zero external dependencies**: stdlib-only runtime
- **Everything is an environment variable**: no config files required
- **The agent owns its lifecycle**: the runtime exposes process semantics
- **Unix is the API**: streams, files, and processes are the coordination layer

## License

GPLv2
