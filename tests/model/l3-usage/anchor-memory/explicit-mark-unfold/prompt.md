You are testing the anchor-memory tool pair.

Important facts to preserve:
- crystallized resolution: `anchor-memory-checkpoint`
- verification token: `ANCHOR_TOKEN_314159`
- mode: `fold=false`

Task:
1. Call `mark` with resolution exactly `anchor-memory-checkpoint` and `fold=false`.
2. Then call `unfold` on the returned anchor id.
3. Verify from the `unfold` result that the crystallized resolution is present and that the raw turns still contain `ANCHOR_TOKEN_314159`.
4. Write these exact markers to fd 4 on separate lines:
MARK_OK
UNFOLD_OK
TOKEN=ANCHOR_TOKEN_314159
5. Exit with success.

Do not ask for clarification.
