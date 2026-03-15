You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 4.

=== PART 1: GATHER ALL (multiple intents) ===
Fork two children that each compute a simple result:
  fork(intents=["sh(command='echo ALPHA_RESULT >&4') && exit success", "sh(command='echo BETA_RESULT >&4') && exit success"])
This should use the default gather-all mode (wait=true, race=false).
Wait for both children to complete. Confirm:
- Both children's output is returned (you see ALPHA_RESULT and BETA_RESULT in the fork result).
- The result header says something like "[FORK SWARM] 2 children completed".

=== PART 2: RACE MODE ===
Fork two children in race mode. One should be fast, one slow:
  fork(intents=["sh(command='echo FAST_WINNER >&4') && exit success", "sh(command='sleep 30 && echo SLOW_LOSER >&4') && exit success"], race=true)
Confirm:
- The fast child wins (you see FAST_WINNER in the result).
- The result header mentions "won" and shows the winner.
- The slow child was killed (not waited on for 30 seconds).

=== PART 3: SINGLE INTENT (backward compat) ===
Fork a single child:
  fork(intents=["sh(command='echo SINGLE_CHILD >&4') && exit success"])
Confirm:
- A single child ran and returned output containing SINGLE_CHILD.

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- GATHER_OK    (Part 1: gather-all mode returned both children's output)
- RACE_OK      (Part 2: race mode returned fast winner, slow child killed)
- SINGLE_OK    (Part 3: single intent fork worked)

Example (if all pass):
GATHER_OK
RACE_OK
SINGLE_OK

Then exit with success.
