You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 4.

This scenario tests the new filesystem-backed detached job model. There is no `job(...)` tool. Use the absolute job directory path (under `${QUINE_DATA_DIR}`) that `sh(detach=true)` returns.

=== PART 1: DETACH AND OBSERVE THE SURFACE ===
Start a background command that takes a few seconds:
  sh(command="sleep 2 && echo DETACH_DONE", detach=true)
You should receive a `[JOB] pid=N path=/absolute/path/... (detached)` result immediately. Reuse that returned `path` value in later commands.
Confirm the directory exposes:
- `cmd`
- `pid`
- `started_at`
- `out.log`
- `err.log`

Confirm the `pid` file matches the returned pid.
When the job finishes, its exit code will appear at `<PATH>/exit`.
Confirm the final exit code is `0`, then inspect:
  sh(command="cat <PATH>/out.log")
Confirm the output contains `DETACH_DONE`.

=== PART 2: EXIT FILE PERSISTS AFTER COMPLETION ===
Start another short detached command:
  sh(command="echo PERSIST_DONE", detach=true)
Confirm `<PATH>/exit` is not present immediately after spawn.
Later, after the job finishes, confirm `<PATH>/exit` appears.
Then read `<PATH>/exit` a SECOND time.
Confirm both reads return `0`, and confirm:
  sh(command="cat <PATH>/out.log")
contains `PERSIST_DONE`.

=== PART 3: KILL A DETACHED JOB WITH SHELL SIGNALS ===
Start a long-running detached command:
  sh(command="sleep 60", detach=true)
Terminate the whole process group using the pid from the detach result:
  sh(command="kill -TERM -<N>")
When the terminal outcome is recorded in `<PATH>/exit`, confirm the exit code is non-zero.

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- DETACH_OK    (Part 1: detach returned immediately with an absolute job path and factual surface files)
- WAIT_OK      (Part 1: the job eventually recorded exit `0` and `out.log` contained `DETACH_DONE`)
- PERSIST_OK   (Part 2: `<path>/exit` appeared later and remained readable multiple times)
- KILL_OK      (Part 3: `kill -TERM -<pid>` terminated the detached job and `<path>/exit` became non-zero)

Example (if all pass):
DETACH_OK
WAIT_OK
PERSIST_OK
KILL_OK

Then exit with success.
