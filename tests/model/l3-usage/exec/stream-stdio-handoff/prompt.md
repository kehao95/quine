You are receiving a live stream on stdin.

Complete both of these requirements:

1. Forward the stream to stdout after converting every lowercase ASCII letter to uppercase.
2. After the stream ends, write exactly one summary line to stderr:
   `BYTES=17`

Constraints:
- treat the input as a stream, not as a quoted literal to retype from memory
- do not use `sh` for this task
- do not exec into `/bin/sh` or any shell wrapper
- do not stage the stream through a temporary file
- keep the stream live all the way through completion
- use ordinary Unix stdin/stdout/stderr for the final streaming path; do not rely on quine-only fd 3/4/5 side channels
- in the final filter process, fd 3 will be invalid; only stdin/stdout/stderr are valid
- prefer an existing utility over a Python/Perl/Ruby/Node wrapper
- the ideal solution is one direct handoff into the final utility, not a bootstrap program
- choose a final utility that can emit the stderr summary itself
- preserve line boundaries on stdout
- emit only the transformed stream on stdout
- emit only the single summary line on stderr
- finish successfully
- A direct `exec` into an existing utility is enough; you do not need a helper program.
- Use `exec` exactly once with an explicit external `target`; do not omit `target`,
  do not self-reenter quine, and do not do a bootstrap exec followed by `sh`.
- One valid shape is an `awk` handoff whose inline program uppercases each line on stdout
  while accumulating `length($0)+1` and finally printing `BYTES=17` to `/dev/stderr`.

The stream you receive is:
alpha
beta
gamma
