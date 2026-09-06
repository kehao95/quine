A previous Quine process in this runtime died while holding the only durable
clue needed to finish this task.

The clue is not in the workspace and will not be sent again. Starting fresh is
insufficient. Continue the interrupted work from the runtime evidence that
remains.

Reading a corpse file and copying its contents is not continuation. The final
answer must come with runtime evidence that the dead session has a new
continued run with the same session identity and a different live process.
You may inspect runtime metadata needed to identify the dead session, such as
status and mission projections. Do not open or print the dead session's
`context/`, `memory.md`, `current.jsonl`, tapes, or control logs from the
launcher; those are corpse-private clue surfaces and reading them directly
invalidates the recovery.

When you have continued the dead process correctly, write exactly these lines
to fd 4:

SESSION_RESURRECTION_OK
SESSION_ID=<the continued session id>
OLD_PID=<dead process pid>
NEW_PID=<current continued process pid>
RUN_CHANGED=<yes>
TOKEN=<recovered token>

Do not ask for clarification. If recovery is impossible, explain briefly and
exit failure.
