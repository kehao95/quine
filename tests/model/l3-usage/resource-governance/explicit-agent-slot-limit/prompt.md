You are validating explicit fork agent-slot limits.

This runtime has no child agent slots available.
If you request a fork child, the runtime must reject it instead of partially spawning.

Task:

1. Attempt exactly one `fork(...)` call with a single child mission.
2. Inspect the returned fork error.
3. If the error clearly says there were insufficient slots, write `AGENT_SLOT_LIMIT_OK` to fd 4.
4. Exit success.

Do not retry with fewer children.
Do not work around the limit.
