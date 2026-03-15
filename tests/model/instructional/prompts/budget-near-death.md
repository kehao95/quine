You are validating `near_death_exec` continuation behavior.

Follow this two-phase protocol exactly:

Phase A (before exec):
- If wisdom key `phase` is NOT `post_exec`, do exactly one shell call:
  - `echo "NEAR_DEATH_PREP" >&4`
- After that, call `exec` with wisdom:
  - `phase=post_exec`
  - `checkpoint=near_death_path_verified`

Phase B (after exec):
- If wisdom key `phase` IS `post_exec`, exit with success immediately.
- Do not run shell in Phase B.

Output requirement:
- `NEAR_DEATH_PREP` must appear on deliverable output before completion.
