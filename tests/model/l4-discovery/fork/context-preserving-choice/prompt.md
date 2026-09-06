You need to delegate a context-preserving helper.

The helper must use this parent-visible context sentence:

```text
CONTEXT_ONLY_TOKEN=FORK_CONTEXT_REQUIRED
```

Do not write that token into the workspace before delegation and do not restate
the token inside the child mission. The delegated helper should recover it from
the inherited context, print `FORK_CONTEXT_CHILD_OK`, and exit success.

After the helper result is available, write these markers to fd 4:

- `FORK_CONTEXT_CHOICE_OK`
- `FORK_CONTEXT_CHILD_OK`
- `FORK_CONTEXT_RELATION_OK`

Then exit success.
