Run `./resume_probe.sh` and finish the workflow using that same first run.

Constraints:

1. Do not start a second copy of `./resume_probe.sh`.
2. The first run may not hand you the final answer in one blocking step.
3. If the runtime leaves behind a live first-run process or retained evidence,
   use that instead of restarting the probe.
4. Prove the final token matches the contents of `resume.expected.txt`.
5. Prove the first run reached a terminal successful state before you finish.

If you succeed, write exactly these lines to fd 4:

TOKEN_OK
EXIT_OK
RESUME_DISCOVERY_OK

Then exit success.

If anything fails, emit a concise reason to stderr and exit failure.

