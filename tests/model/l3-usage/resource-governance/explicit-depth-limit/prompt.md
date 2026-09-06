You are validating explicit fork depth limits.

This runtime has a strict recursion depth limit.
Attempting to fork a child when the next child depth would reach that limit must be rejected by the runtime.

Task:

1. Attempt exactly one `fork(...)` call with a single child mission.
2. Inspect the returned fork error.
3. If the error clearly says the max recursion depth was exceeded, write `DEPTH_LIMIT_OK` to fd 4.
4. Exit success.

Do not retry with a different fork.
Do not work around the limit.
