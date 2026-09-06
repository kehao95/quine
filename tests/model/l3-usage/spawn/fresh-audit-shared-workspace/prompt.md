Use the `spawn` tool explicitly to verify fresh-context review of a current
workspace artifact.

Do exactly these steps:

1. Use `sh` to write `candidate.txt` in the current workspace. It must contain
   these two lines:
   `status: claimed-correct`
   `gap: missing empty-input behavior`
2. Use one `spawn` call with `mode="wait"` and one child. The child mission is:
   `Inspect candidate.txt in the current workspace. If it is visible and you find a gap line, print SPAWN_SHARED_WORKSPACE_FOUND and SPAWN_SHARED_WORKSPACE_GAP, then exit success.`
3. Inspect the spawn result. Confirm that the child saw the workspace file,
   reported the gap, retained relation/member process handles, and did not
   expose fork-style `seed_root`.
4. Write these markers to fd 4:
   `SPAWN_SHARED_WORKSPACE_OK`
   `SPAWN_SHARED_WORKSPACE_GAP`
   `SPAWN_SHARED_WORKSPACE_RELATION_OK`

Then exit success.
