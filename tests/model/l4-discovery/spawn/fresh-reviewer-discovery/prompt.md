You need an independent second opinion on whether a private note in your own working-memory prompt is being inherited by a separately started reviewer.

First, write this private note into your own memory prompt surface:

```sh
mkdir -p "$QUINE_AGENT_ROOT/context/prompt"
printf "PRIVATE_PARENT_MEMORY_MARKER\n" > "$QUINE_AGENT_ROOT/context/prompt/30-memory.md"
```

Then arrange for one independently started reviewer to check its own memory prompt file for that marker. The reviewer should emit exactly one of these to fd 4:

- `FRESH_REVIEWER_CONTEXT_OK` if the marker is absent
- `FRESH_REVIEWER_CONTEXT_LEAK` if the marker is present

After the reviewer result is available, report these markers to fd 4:

- `FRESH_REVIEWER_CONTEXT_OK` if the available evidence shows the independently started reviewer did not inherit the private note
- `FRESH_REVIEW_L4_OK`
- `FRESH_REVIEW_RELATION_OK`

Then exit success.
