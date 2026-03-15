You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 4.

This scenario tests the new filesystem-backed detached job model. There is no `job(...)` tool. Use the absolute job directory path (under `${QUINE_DATA_DIR}`) that `sh(detach=true)` returns.

=== PART 1: DETACH AND WAIT VIA EXIT FILE ===
Start a background command that takes a few seconds:
  sh(command="sleep 2 && echo DETACH_DONE", detach=true)
You should receive a `[JOB] pid=N path=/absolute/path/... (detached)` result immediately. Reuse that returned `path` value in later commands.
Then wait for completion using:
  sh(command="cat <PATH>/exit")
Confirm the exit code is `0`, then inspect:
  sh(command="cat <PATH>/out.log")
Confirm the output contains `DETACH_DONE`.

=== PART 2: EXIT FILE PERSISTS AFTER COMPLETION ===
Start another short detached command:
  sh(command="echo PERSIST_DONE", detach=true)
Wait once with:
  sh(command="cat <PATH>/exit")
Then call the same command a SECOND time:
  sh(command="cat <PATH>/exit")
Confirm both reads return `0`, and confirm:
  sh(command="cat <PATH>/out.log")
contains `PERSIST_DONE`.

=== PART 3: KILL A DETACHED JOB WITH SHELL SIGNALS ===
Start a long-running detached command:
  sh(command="sleep 60", detach=true)
Terminate the whole process group using the pid from the detach result:
  sh(command="kill -TERM -<N>")
Then wait for completion:
  sh(command="cat <PATH>/exit")
Confirm the exit code is non-zero.

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- DETACH_OK    (Part 1: detach returned immediately with an absolute job path)
- WAIT_OK      (Part 1: `cat <path>/exit` waited and `out.log` contained DETACH_DONE)
- PERSIST_OK   (Part 2: reading `<path>/exit` twice returned the same exit code)
- KILL_OK      (Part 3: `kill -TERM -<pid>` terminated the detached job and `exit` became non-zero)

Example (if all pass):
DETACH_OK
WAIT_OK
PERSIST_OK
KILL_OK

Then exit with success.
