You are testing transactional workspace rewind via `switch_world`.

The current working directory is an empty managed workspace.

Follow this protocol exactly:

1. Run one `sh` call that creates `state.txt` with exactly:
   `v1`
2. Observe that the tool result reports current world revision `wr1`.
3. Run one `sh` call that overwrites `state.txt` with exactly:
   `v2`
4. Observe that the tool result reports current world revision `wr2`.
5. Call `switch_world` with:
   - `target="wr1"`
6. Run one final `sh` call that:
   - verifies `state.txt` now contains `v1`
   - writes these two markers to fd 4 on separate lines:
     - `RESTORE_OK`
     - `STATE_V1_OK`
7. Exit with success.

Important:
- Do not use `exec`.
- Do not use `fork`.
- Do not recreate `state.txt` after the switch; the point is to verify the revision restore itself.
