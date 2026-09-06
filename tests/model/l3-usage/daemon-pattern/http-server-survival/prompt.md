You have a composite task with 2 parts. Complete ALL parts and output a structured report to fd 4.

This scenario tests the **daemon pattern**: starting a background server with detach=true, verifying it works, then exiting successfully so it survives quine's exit.

=== PART 1: START A DAEMON ===
Start a simple HTTP server in the background using detach=true:
  sh(command="python3 -m http.server 18923", detach=true)
You should receive `[JOB] pid=N path=/absolute/path/... (detached)` immediately.
Wait briefly (1-2 seconds) for the server to start up.

=== PART 2: VERIFY DAEMON IS ALIVE ===
Confirm the daemon is running and serving requests:
  sh(command="curl -s -o /dev/null -w '%{http_code}' http://localhost:18923/")
The HTTP status should be 200. This confirms the daemon started successfully.

Do NOT kill the daemon. The entire point is that it survives.
Exit with status "success".

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- DAEMON_STARTED   (Part 1: sh(detach=true) returned an absolute job path)
- DAEMON_ALIVE     (Part 2: curl got HTTP 200 from the daemon)

Example (if all pass):
DAEMON_STARTED
DAEMON_ALIVE

Then exit with success. Do NOT kill or stop the daemon before exiting.
