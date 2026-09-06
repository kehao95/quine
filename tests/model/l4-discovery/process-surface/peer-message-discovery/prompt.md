Another live Quine process is already sharing your runtime root.

Without using `ps`, tape/audit logs, `agent/self`, or any legacy live-process
shortcut, determine all of the following from the runtime process surface:

- your own session id
- your own pid
- the session id of the other live Quine process
- the pid of the other live Quine process

Then send one non-empty single-line ASCII payload of your choice to the other
live Quine through its public control surface.

Requirements:

1. Do not spawn or fork another Quine.
2. Use only the runtime surface that is already available to you.
3. Do not read `locks/`.
4. Do not use any legacy `agent/live` or `agent/live_by_pid` shortcut.
5. If you succeed, write exactly these lines to fd 4:

PROCESS_SURFACE_COMM_OK
SELF_SESSION=<your session id>
SELF_PID=<your pid>
NEIGHBOR_SESSION=<neighbor session id>
NEIGHBOR_PID=<neighbor pid>
PEER_MESSAGE=<exact payload you wrote>

6. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not ask for clarification.
