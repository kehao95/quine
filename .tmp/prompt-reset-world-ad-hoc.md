I have a SQLite database in WAL mode in `__APP_ROOT__`.
However, the WAL file appears to be corrupted or encrypted. When you try to
access the database, SQLite may only show the base data (5 records) instead of
all 11 records that should be there.

Your task is to:

1. Fix the WAL file so SQLite can read it
2. Extract ALL data from the database, including WAL changes
3. Create `__APP_ROOT__/recovered.json`

The output JSON should be sorted by `id` and should contain all 11 records.

If you determine the mission can no longer be completed from the current world
state, explain why briefly and exit failure promptly.

Constraints:

- Do not use `fork`.
- Do not access, inspect, search, or copy from any path outside `__APP_ROOT__`.

If your own probing destroys or checkpoints away unique WAL evidence, you may use `exec(reset_world=true, wisdom={...})` to restart from a fresh workspace world. Use it only as a rollback move after you have evidence the current world is corrupted by your own action.
