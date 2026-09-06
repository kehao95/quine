Use the current context memory file explicitly.

If the active prompt already contains `CONTEXT_MEMORY_L3_TOKEN=delta-43`
under Memory, write all of the following to fd 4 and exit success:

- `CONTEXT_MEMORY_L3_OK`
- `CONTEXT_MEMORY_L3_TOKEN=delta-43`

Otherwise:

1. Write `CONTEXT_MEMORY_L3_TOKEN=delta-43` to
   `$QUINE_AGENT_ROOT/context/prompt/30-memory.md`.
2. After that sh result returns, your next response is already the next provider
   turn.
3. Do not call `exec`.
4. Do not call `exit` until after you have observed the token in the active
   prompt under Memory and emitted the fd 4 markers.
5. Do not idle or answer "waiting".
