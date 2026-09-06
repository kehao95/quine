You are testing timeout-interrupted synchronous `sh` behavior.

This runtime supports `sh(timeout=<seconds>)`. In direct/host mode, if a sync
`sh` call exceeds that timeout, the result becomes `status="interrupted"` and
includes `job.pid`, `job.path`, `stdout_so_far`, and `stderr_so_far`. The
process group is stopped, not killed.

Your task:

1. Launch `./resume_probe.sh` exactly once with a synchronous `sh` call and
   `timeout=1`.
2. If the call returns interrupted, keep using that same first run. Do not
   launch a second copy of `./resume_probe.sh`.
3. Continue the paused process.
4. Verify that `<PATH>/exit` eventually appears and is exactly `0`.
5. Read the first run's `out.log` and verify it contains the exact token from
   `resume.expected.txt`.
6. Write these exact lines to fd 4 on separate lines:
   - `INTERRUPTED_OK`
   - `CONT_OK`
   - `EXIT_OK`
   - `TOKEN_OK`
   - `TIMEOUT_RESUME_OK`
7. Exit success.

If anything fails, emit a concise reason to stderr and exit failure.

