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

| Variable | Required | Description |
|----------|----------|-------------|
| `QUINE_MODEL_ID` | ✓ | Model name sent to the API |
| `QUINE_API_TYPE` | ✓ | Wire protocol: `openai` or `anthropic` |
| `QUINE_API_BASE` | ✓ | API base URL |
| `QUINE_API_KEY` | ✓ | API key |
| `QUINE_CONTEXT_WINDOW` | | Context window size in tokens (default 128000) |
| `QUINE_MAX_DEPTH` | | Max recursion depth (default 5) |
| `QUINE_MAX_TURNS` | | Max conversation turns, 0 = unlimited (default 20) |
| `QUINE_DATA_DIR` | | Session log directory (default `.quine/`) |

> **Tip:** Every line in your `.env` must start with `export` so that `source .env` propagates variables to child processes.

## Usage

```bash
# Run a single task
quine "Write a haiku about recursion"

# Pipe input
echo "What is 2+2?" | quine "Answer the question"
```

**That's it.** The agent can read/write files, run shell commands, and spawn child agents.

## Design Principles

- **Zero external dependencies** — stdlib only
- **Everything is an environment variable** — no flags, no files, no magic
- **The agent owns its lifecycle** — it calls `exit`, not the runtime
- **Unix is the API** — pipes, processes, and files are the coordination primitives
- **Fractal architecture** — a tree of identical processes, scale-invariant

## License

GPLv2
