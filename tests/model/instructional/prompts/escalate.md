You have a multi-part task testing the escalate tool.

The escalate tool allows you to request a more capable model when needed.
Call escalate(reason="...") to upgrade — the smarter model will see your full context.
Escalation is one-way and preserves all prior conversation.

=== PART 1: INITIAL WORK ===
Create a test file to demonstrate work before escalation:
  sh(command="mkdir -p /tmp/escalate-test && echo 'Phase 1 complete' > /tmp/escalate-test/phase1.txt")
Verify the file exists and contains the expected content:
  sh(command="cat /tmp/escalate-test/phase1.txt")

=== PART 2: ESCALATE ===
Now escalate to the smarter model. Your reason should explain that you want
the smarter model to verify and continue the work.
  escalate(reason="Phase 1 complete. Requesting smarter model to verify work and complete Phase 3.")

=== PART 3: POST-ESCALATION (for the smarter model) ===
After escalation, you are now running as the smarter model.
Your job is to:
1. Review the conversation history — you should see Part 1 work was completed
2. Verify /tmp/escalate-test/phase1.txt exists and contains "Phase 1 complete"
3. Create /tmp/escalate-test/verified.txt with content "Verified by smart model"
4. Read back both files to confirm

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- PHASE1_OK      (Part 1: initial file created and verified)
- ESCALATED_OK   (Part 2: escalation tool call succeeded — you received confirmation)
- VERIFIED_OK    (Part 3: smart model verified prior work and created verification file)

Example output if all pass:
PHASE1_OK
ESCALATED_OK
VERIFIED_OK

Then exit with success.
