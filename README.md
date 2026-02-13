# quine

A recursive POSIX Agent runtime. One binary, zero dependencies, infinite depth.

> *"We do not need to teach AI to be intelligent. We only need to build an environment where nothing else can survive."*
> — **[The Quine Manifesto](./QUINE_MANIFESTO.md)**

## Get Started

```bash
# Install
go install github.com/kehao95/quine/cmd/quine@latest

# Set your API key (Anthropic or OpenAI)
export ANTHROPIC_API_KEY=sk-ant-...
# or
export OPENAI_API_KEY=sk-...

# Run
quine "Write a haiku about recursion"

# Pipe input
echo "What is 2+2?" | quine "Answer the question"
```

**That's it.** The agent can read/write files, run shell commands, and spawn child agents.

## Configuration

All configuration via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `QUINE_MODEL_ID` | `claude-sonnet-4-5-20250929` | Model to use |
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
