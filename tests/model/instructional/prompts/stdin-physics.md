You are given piped text material on fd 3. Complete ALL parts and output a structured report to fd 4.

The piped material is exactly the 10-byte string:
ABCDEfghij

=== PART 1: sh(stdin=...) is separate from material fd 3 ===
Use `sh(command="cat", stdin="TOOL_ONLY\n")`.
Confirm the output contains exactly `TOOL_ONLY`.
Do NOT read fd 3 in this part.

=== PART 2: fd 3 is a live consumed stream ===
Use separate `sh` calls for these steps:
1. Read exactly the first 5 bytes from fd 3 into `/tmp/stdin-physics/part1.txt`
2. Read the remaining bytes from fd 3 into `/tmp/stdin-physics/part2.txt`
3. Verify:
   - `part1.txt` is `ABCDE`
   - `part2.txt` is `fghij`
   - `part2.txt` is NOT the full original string

=== PART 3: fork children do not inherit the parent's material bytes ===
Use `fork` with exactly one child mission. The child mission must:
1. In one `sh` call, attempt to read `/dev/fd/3`
2. If `/dev/fd/3` is empty or unreadable, print `CHILD_NO_PARENT_MATERIAL` to stdout and exit success
3. If the child can read any of the parent's material bytes (`ABCDEfghij` or the unread remainder `fghij`) from fd 3, print `CHILD_INHERITED_PARENT_MATERIAL` to stdout and exit failure

Wait for the child result. Confirm the fork result contains `CHILD_NO_PARENT_MATERIAL` and NOT `CHILD_INHERITED_PARENT_MATERIAL`.

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- TOOL_STDIN_OK   (Part 1: sh(stdin=...) worked without using fd 3)
- STREAM_OK       (Part 2: fd 3 behaved like a live consumed stream)
- FORK_FD_OK      (Part 3: fork child did not inherit the parent's material)

Example (if all pass):
TOOL_STDIN_OK
STREAM_OK
FORK_FD_OK

Then exit with success.
