You are testing explicit child-world adoption.

The current working directory is an empty managed workspace.

Follow this protocol exactly:

1. Use exactly one `fork(mode="wait")` call with one child at `scope="."`.
2. The child intent must:
   - use exactly one `sh` call with command `printf "child-line\n" > adopted.txt && echo CHILD_DONE >&4`
   - then exit success
3. From the parent, use exactly one `sh` call that verifies `adopted.txt` is absent and emits `PRE_SWITCH_PRIVATE` to fd 4.
4. Read the child world handle from the fork result and call `switch_world` with that exact handle.
5. Use exactly one final `sh` call that verifies `adopted.txt` now contains `child-line`, then writes both:
   - `ADOPT_SWITCH_OK`
   - `HANDLE_OK`
   to fd 4.
6. Exit success.

Important:
- Do not use `exec`.
- Do not create or rewrite `adopted.txt` in the parent before switching.
