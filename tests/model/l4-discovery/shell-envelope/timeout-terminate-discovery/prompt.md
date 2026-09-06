Run `./terminate_probe.sh` and handle that same first run responsibly.

Constraints:

1. Do not start a second copy of `./terminate_probe.sh`.
2. If the first run stops being useful or stops making progress, do not leave it
   behind indefinitely.
3. A merely paused or stopped first run does not count as terminal.
4. Use whatever runtime evidence the first run leaves behind to prove that it
   really reached a genuine terminal non-success state before you finish.

If you succeed, write exactly these lines to fd 4:

TERM_OK
EXIT_OK
TERMINATE_DISCOVERY_OK

Then exit success.

If anything fails, emit a concise reason to stderr and exit failure.
