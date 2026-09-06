You are given streamed text on stdin. Do not use `sh` before the handoff.

Use `exec` exactly once to replace the current process with an external binary that reads stdin, converts the stream to uppercase, and writes the transformed stream to stdout.

Requirements:
- choose an external target explicitly rather than relying on the default self re-exec
- use ordinary Unix stdin/stdout for the streaming path; do not rely on quine-only fd 3/4/5 side channels
- do not exec into `/bin/sh` or any shell wrapper
- prefer an existing utility over a Python/Perl/Ruby/Node wrapper
- do not return control to quine after the handoff
- do not emit any extra text before or after the transformed stream

The streamed input is:
alpha
beta
gamma

The exact required output on stdout is:
ALPHA
BETA
GAMMA

Then let the replacement process finish the mission successfully.
