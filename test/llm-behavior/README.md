# LLM Behavior Tests

Verifies that an LLM can correctly understand and use quine's tool API purely from the system prompt and schema. These are **not unit tests** — they make real API calls and produce tapes for human review.

## Why

Unit tests verify the engine works. These tests verify the **API is learnable** — that the system prompt and tool schema are sufficient for a model to figure out correct usage without prior training.

## When to Run

- After changing `system_prompt.md`
- After changing tool schemas (`schemas.go`)
- After changing tool output formats (`tools.go`)
- When evaluating a new model

## Usage

```bash
# Build quine first
go build -o /tmp/quine ./cmd/quine/

# Source credentials
source .env.claude-sonnet-4-5  # or any model

# Run a specific scenario
./test/llm-behavior/run.sh sessions

# Run all scenarios
./test/llm-behavior/run.sh all

# Review results
ls test/llm-behavior/runs/          # tape files + scores
cat test/llm-behavior/runs/latest/score.txt
```

## Scenarios

| ID | Name | Tests | Key Question |
|----|------|-------|-------------|
| `sessions` | Named Session API | anonymous, session create, persistence, isolation, read-only, background+read | Can the LLM discover all 3 sh modes? |

## Interpreting Results

Each run produces:
- `stdout.txt` — agent's fd 3 output (deliverable)
- `stderr.txt` — agent's stderr
- `tape.jsonl` — full conversation tape (the source of truth)
- `tape.log` — human-readable summary
- `score.txt` — automated check results (PASS/FAIL per criterion)

**Automated scoring is a heuristic.** Always review `tape.log` when a check fails — the agent may have achieved the goal through an unexpected path.

## Baseline Results

| Date | Model | Scenario | C1 | C2 | C3 | C4 | C5 | C6 | C7 | Score | Turns | Tokens | Time |
|------|-------|----------|----|----|----|----|----|----|----|-------|-------|--------|------|
| 2026-02-21 | claude-sonnet-4-5 | sessions | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 7/7 | 9 | 37065 | 33s |
| 2026-02-21 | kimi-k2.5 | sessions | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 7/7 | 9 | 12346 | 30s |

## Adding Scenarios

1. Create `prompts/<name>.md` — the mission prompt
2. Add scoring logic in `run.sh` `score_scenario()` function
3. Document in this README
