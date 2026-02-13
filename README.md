# quine

A recursive POSIX Agent runtime. One binary, zero dependencies, infinite depth.

> *"We do not need to teach AI to be intelligent. We only need to build an environment where nothing else can survive."*
> — **[The Quine Manifesto](./QUINE_MANIFESTO.md)**

## Get Started

```bash
# Install
go install github.com/kehao95/quine/cmd/quine@latest

# Configure (copy and edit)
cp .env.example .env
# Edit .env — set your model and API key

# Run
source .env && quine "Write a haiku about recursion"

# Pipe input
echo "What is 2+2?" | source .env && quine "Answer the question"
```

**That's it.** The agent can read/write files, run shell commands, and spawn child agents.

### Quick Setup for Known Providers

For **Anthropic** or **OpenAI** models, you only need two things — the model name and an API key. Everything else is auto-detected:

```bash
# Anthropic (claude-* models)
export QUINE_MODEL_ID=claude-sonnet-4-5-20250929
export ANTHROPIC_API_KEY=sk-ant-...

# OpenAI (gpt-*, o1-*, o3-*, o4-* models)
export QUINE_MODEL_ID=gpt-4o
export OPENAI_API_KEY=sk-...
```

### Third-Party / Custom Providers

For any OpenAI-compatible API (Moonshot, Together, Ollama, vLLM, etc.), set all four fields explicitly:

```bash
export QUINE_MODEL_ID=kimi-k2.5
export QUINE_API_TYPE=openai
export QUINE_API_BASE=https://api.moonshot.ai/v1
export QUINE_API_KEY=sk-your-key-here
```

> **Tip:** Every line in your `.env` must start with `export` so that `source .env` propagates variables to child processes.

## Configuration

All configuration via environment variables. See [`.env.example`](./.env.example) for a documented template.

| Variable | Default | Description |
|----------|---------|-------------|
| `QUINE_MODEL_ID` | `claude-sonnet-4-5-20250929` | Model name sent to the API |
| `QUINE_API_TYPE` | *(auto-detected)* | Wire protocol: `openai` or `anthropic` |
| `QUINE_API_BASE` | *(auto-detected)* | API base URL |
| `QUINE_API_KEY` | *(from provider key)* | API key (falls back to `OPENAI_API_KEY` / `ANTHROPIC_API_KEY`) |
| `QUINE_CONTEXT_WINDOW` | *(auto-detected)* | Context window size in tokens |
| `QUINE_MAX_DEPTH` | `5` | Max recursion depth |
| `QUINE_MAX_TURNS` | `20` | Max conversation turns (0 = unlimited) |
| `QUINE_DATA_DIR` | `.quine/` | Where to store session logs |

See **[Artifacts/](./Artifacts/)** for the full theoretical framework.

## Design Principles

- **Zero external dependencies** — stdlib only
- **Everything is an environment variable** — no flags, no files, no magic
- **The agent owns its lifecycle** — it calls `exit`, not the runtime
- **Unix is the API** — pipes, processes, and files are the coordination primitives
- **Fractal architecture** — a tree of identical processes, scale-invariant

## License

MIT
