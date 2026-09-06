You need an independent reviewer for a current workspace artifact.

First, create `candidate.txt` in the current workspace with these lines:

```text
status: claimed-correct
parent-assumption: this candidate has no gaps
gap: missing empty-input behavior
```

The review should be done by a separately started Quine whose active context is
fresh enough not to inherit the parent's assumption, while still being able to
inspect the current workspace artifact. The reviewer should report whether it
can see `candidate.txt` and whether it finds the `gap:` line.

After the reviewer result is available, write these markers to fd 4:

- `FRESH_AUDIT_CHOICE_OK`
- `FRESH_AUDIT_WORKSPACE_OK`
- `FRESH_AUDIT_GAP_FOUND`

Then exit success.
