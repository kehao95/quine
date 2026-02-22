You have a composite task with 5 parts. Complete ALL parts and output a structured report to fd 3.

=== PART 1: NORMAL COMPLETION ===
Run a command that finishes within its budget:
  sh(command="echo NORMAL_OK", output_limit=1000)
Confirm you received "[EXIT CODE] 0" with "NORMAL_OK" in the output.
(A command that finishes before the limit should NOT produce [PAUSED].)

=== PART 2: OUTPUT LIMIT PAUSE ===
Run a command that will exceed the output limit:
  sh(command="seq 1 10000", output_limit=200)
This should return [PAUSED] with a job ID. Extract the job ID from the response.
Do NOT skip this — the [PAUSED] response is expected and correct behavior.

=== PART 3: RESUME PAUSED JOB ===
Resume the paused job from Part 2 with a larger budget:
  job(id=<ID from Part 2>, signal="cont", output_limit=50000)
This should eventually complete with "[EXIT CODE] 0".
Confirm you see higher numbers (e.g. numbers > 100) in the continued output.

=== PART 4: KILL A JOB ===
Start a long-running command with a small output limit:
  sh(command="seq 1 100000", output_limit=100)
It will pause. Then kill it:
  job(id=<new ID>, signal="kill")
Confirm you received a "killed" confirmation.

=== PART 5: READ WITHOUT RESUMING ===
Start another paused job:
  sh(command="yes | head -1000", output_limit=50)
It will pause. Read its accumulated output WITHOUT resuming:
  job(id=<new ID>)   ← no signal parameter
Confirm you see output (the "y" lines). Then kill it to clean up:
  job(id=<same ID>, signal="kill")

=== OUTPUT ===
Write a report to fd 3 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- NORMAL_OK      (Part 1: normal completion worked, no [PAUSED])
- PAUSE_OK       (Part 2: [PAUSED] was returned with a job ID)
- RESUME_OK      (Part 3: resumed job completed with exit code 0)
- KILL_OK        (Part 4: kill returned "killed" confirmation)
- READ_OK        (Part 5: read-only job() returned output)

Example (if all pass):
NORMAL_OK
PAUSE_OK
RESUME_OK
KILL_OK
READ_OK

Then exit with success.
