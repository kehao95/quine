You are receiving a live stream on stdin.

Your job is to forward that stream to stdout after converting every lowercase ASCII letter to uppercase.

Constraints:
- treat the input as a stream, not as a quoted literal to retype from memory
- do not use `sh` for this task
- do not stage the stream through a temporary file
- keep the stream live all the way through completion
- use ordinary Unix stdin/stdout for the streaming data path; do not rely on quine-only fd 3/4/5 side channels
- in the final filter process, fd 3 will be invalid; only stdin/stdout are valid for the data path
- this transformation is something a standard Unix filter can already do directly
- prefer an existing utility over a Python/Perl/Ruby/Node wrapper
- the ideal solution is one direct handoff into the final utility, not a bootstrap program
- preserve line boundaries
- emit only the transformed stream on stdout
- finish successfully

The stream you receive is:
alpha
beta
gamma
