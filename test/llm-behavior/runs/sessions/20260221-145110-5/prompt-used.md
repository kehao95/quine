You have a composite task with 6 parts. You MUST complete ALL parts and output a structured report to fd 3.

=== PART 1: ANONYMOUS EXECUTION ===
Run an anonymous sh command (no session parameter): echo "ANON_OK"
Confirm you received the output.

=== PART 2: SESSION CREATION ===
Create a named session called "alpha":
  sh(command="cd /tmp && export MARKER=session_alpha_42 && echo SESSION_CREATED", session="alpha")

=== PART 3: SESSION PERSISTENCE ===
In the SAME session "alpha", verify state persisted:
  sh(command="echo PERSIST_${MARKER}_$(basename $(pwd))", session="alpha")
Expected output: PERSIST_session_alpha_42_tmp
If you see this, persistence works.

=== PART 4: SESSION ISOLATION ===
Create a SECOND session called "beta":
  sh(command="echo MARKER_is_${MARKER:-unset}", session="beta")
Expected: MARKER_is_unset (beta should NOT see alpha's variable)
If MARKER is unset in beta, isolation works.

=== PART 5: READ-ONLY MODE ===
Use the read-only mode to drain output from session "alpha" (call sh with ONLY session parameter, NO command):
  sh(session="alpha")
This should return any accumulated output. Report what you got.

=== PART 6: BACKGROUND + READ ===
Start a slow command in a new session "worker":
  sh(command="for i in 1 2 3; do echo TICK_$i; sleep 2; done", session="worker")
This will return partial output. Then do other work:
  sh(command="echo doing_other_work")
Then read the worker's accumulated output:
  sh(session="worker")
Report the ticks you observed.

=== OUTPUT ===
After completing all 6 parts, output a report to fd 3 with these EXACT markers on separate lines:
- ANON_OK (if part 1 succeeded)
- SESSION_CREATED (if part 2 succeeded)
- PERSIST_OK (if part 3 showed the correct persisted state)
- ISOLATE_OK (if part 4 showed MARKER was unset in beta)
- READ_OK (if part 5 successfully used read-only mode)
- BG_READ_OK (if part 6 showed at least one TICK from the read)

Example output format (write this to >&3):
ANON_OK
SESSION_CREATED
PERSIST_OK
ISOLATE_OK
READ_OK
BG_READ_OK

Then exit with success.
