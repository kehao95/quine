You are validating retained relation outcome semantics.

Do these steps:

1. Use `fork` with `mode` set to `wait` and exactly one child at scope `.`. The child intent is:

   `Use exactly one sh call with command: printf 'RELATION_CHILD_FAIL\n' >&4. Then exit failure.`

2. Use `spawn` with `mode` set to `forget` and exactly one child mission:

   `Use exactly one sh call with command: printf 'RELATION_SPAWN_FORGET_CHILD\n' >&4. Then exit success.`

3. Inspect the retained relation evidence from each tool result. Verify that the all-failed `fork(wait)` relation has `status=completed` and no successful child exit code, while the `spawn(forget)` relation has `status=spawned` and `spawned=1`.

4. Write these exact markers to fd 4:

   `RELATION_L3_OK`
   `RELATION_ALL_FAILED_OK`
   `RELATION_FORGET_SPAWNED_OK`

Then exit success.
