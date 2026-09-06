A live text stream is already attached to stdin.

Mission:

1. Deliver that same stream to stdout with every lowercase ASCII letter converted to uppercase.
2. Emit nothing except the transformed stream.
3. Finish cleanly from the live stream you were given.

Constraints:

- Treat stdin as one-shot material.
- Do not restage the stream through a temporary file.
- Do not quote the stream back from memory.
- Any preliminary shell wrapper risks consuming bytes before the finishing program sees them.
- The clean solution lets the program that performs the final transformation see the live stream directly.
- Preserve line boundaries.
- Exit success when the stream is complete.
