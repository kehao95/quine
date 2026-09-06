Use the `spawn` tool explicitly to verify fresh child cognition.

Do exactly these steps:

1. Use one `sh` call to write `SPAWN_PARENT_MEMORY_MARKER` into your own memory prompt file:
   `mkdir -p "$QUINE_AGENT_ROOT/context/prompt" && printf "SPAWN_PARENT_MEMORY_MARKER\n" > "$QUINE_AGENT_ROOT/context/prompt/30-memory.md"`
2. Use one `spawn` call with `mode="wait"` and one child mission:
   `Exit success immediately without using sh.`
3. Inspect the spawn result and the retained child context. Confirm it contains a retained relation handle and child process handles, that it does not expose a fork-style `seed_root`, and that the child's retained `context/prompt/30-memory.md` does not contain `SPAWN_PARENT_MEMORY_MARKER`.
4. Use one final `sh` call to write these markers to fd 4:
   `SPAWN_L3_OK`
   `SPAWN_RELATION_OK`
   `SPAWN_FRESH_CONTEXT_OK`

Then exit success.
