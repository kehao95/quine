You are testing timeout-interrupted synchronous `sh` behavior.

This runtime supports `sh(timeout=<seconds>)`. In direct/host mode, if a sync
`sh` call exceeds that timeout, the result becomes `status="interrupted"` and
includes `job.pid`, `job.path`, `stdout_so_far`, and `stderr_so_far`. The
process group is stopped, not killed.

Your task:

1. Launch `./terminate_probe.sh` exactly once with a synchronous `sh` call and
   `timeout=1`.
2. If the call returns interrupted, keep using that same first run. Do not
   launch a second copy of `./terminate_probe.sh`.
3. Explicitly terminate the paused process. If a softer signal is insufficient,
   escalate until the first run is truly terminal.
4. Verify that `<PATH>/exit` eventually appears and is non-zero.
5. Write these exact lines to fd 4 on separate lines:
   - `INTERRUPTED_OK`
   - `TERMINATED_OK`
   - `EXIT_OK`
   - `TIMEOUT_TERMINATE_OK`
6. Exit success.

If anything fails, emit a concise reason to stderr and exit failure.

