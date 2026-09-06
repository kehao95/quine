You have a three-phase mission with a hard local boundary before final verification.

Phase 1:

1. Create `/tmp/escalate-test/phase1.txt` with content `Phase 1 complete`.
2. Verify that the file exists and still contains that exact content.

Phase 2:

Once phase 1 is complete, move the mission onto the strongest continuation path available in this runtime that can inherit the full conversation and the work already done so far.

Phase 3:

From that stronger continuation, verify `/tmp/escalate-test/phase1.txt` still exists and still contains `Phase 1 complete`, then create `/tmp/escalate-test/verified.txt` with content `Verified by smart model`.

Output:

Write to fd 4 on separate lines:

- `PHASE1_OK`
- `ESCALATED_OK`
- `VERIFIED_OK`

Then exit success.
