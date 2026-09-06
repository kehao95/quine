Another live Quine process is already sharing your runtime root.

Without using `ps`, tape/audit logs, `agent/self`, or any legacy live-process
shortcut, discover the runtime process surface, identify the other live process,
and deliver exactly one peer inject event.

Requirements:

1. Do not spawn or fork another Quine.
2. Use only the already-available runtime filesystem surface.
3. Discover all of the following:
   - your own session id
   - your own pid
   - the session id of the other live Quine process
   - the pid of the other live Quine process
4. Build one non-empty single-line ASCII payload of your choice.
5. Verify from the runtime filesystem surface that delivery occurred.
   Treat inbox receipt alone as insufficient; confirm the delivered `inject`
   state from the peer's live `context/state/current.jsonl` and/or retained
   `log/<session>/control.jsonl`.
6. If you succeed, write exactly these lines to fd 4:

PROCESS_SURFACE_INJECT_OK
SELF_SESSION=<your session id>
SELF_PID=<your pid>
NEIGHBOR_SESSION=<neighbor session id>
NEIGHBOR_PID=<neighbor pid>
PEER_MESSAGE=<exact payload you wrote>
DELIVERY=inject

7. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not ask for clarification.
