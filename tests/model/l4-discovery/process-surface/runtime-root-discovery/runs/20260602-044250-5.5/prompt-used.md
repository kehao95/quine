Another live Quine process is already sharing your runtime root.

Without using `ps`, tape/audit logs, or any `agent/self` shortcut, determine all of the following from the public runtime filesystem surface:

- your own session id
- your own pid
- the session id of the other live Quine process
- the pid of the other live Quine process

Requirements:

1. Do not spawn or fork another Quine.
2. Use the runtime surface that is already available to you.
3. `pid/<pid>` symlink targets resolve to `agent/<session>/public`.
   The agent root is the parent directory of that `public` target, and the
   session id is that parent directory's basename.
3. If you succeed, write exactly these lines to fd 4:

SELF_SURFACE_OK
NEIGHBOR_SURFACE_OK
PROCESS_SURFACE_DISCOVERY_OK
SELF_SESSION=<your session id>
SELF_PID=<your pid>
NEIGHBOR_SESSION=<neighbor session id>
NEIGHBOR_PID=<neighbor pid>

4. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not ask for clarification.
