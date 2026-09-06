Use exactly one fork call with mode "wait" and one child at scope ".".

Child intent:
- use exactly one sh call with command `printf 'REL_CHILD_OK\n' >&4`
- then exit success

After the fork call returns, inspect the fork result carefully.

You must verify all of these from the returned relation/member handles:
1. the child stdout snapshot contains `REL_CHILD_OK`
2. the fork result includes a non-empty `relation_id`, `relation_root`, and `relation_handle`
3. the first child entry includes non-empty `session_id`, `retained_root`, `seed_root`, `status_path`, and `control_path`
4. on disk, `relation_root` contains:
   - `relation.json`
   - `status.json`
   - `result.json`
   - `members/000.json`
5. on disk, `seed_root` contains `origin.json`

Use shell verification for the file checks; do not merely assume they exist.

=== OUTPUT ===
Write exactly these success markers, one per line, to fd 4:
- CHILD_STDOUT_OK
- RELATION_ROOT_OK
- HANDLE_PATHS_OK
- SEED_ORIGIN_OK

Then exit success.
