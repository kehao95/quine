You are running with transactional workspace physics enabled.

There is a file at `/tmp/quine-workspace-absolute/config.txt`. Complete all parts using explicit absolute paths:

1. Read `/tmp/quine-workspace-absolute/config.txt`.
2. Change the line `mode=alpha` to `mode=beta`.
3. Then run a second edit attempt that tries to replace a string that does not exist:
   change `mode=gamma` to `mode=delta`.
4. Use the `[FS MUTATIONS]` channel from the second command to determine whether that second edit actually changed the filesystem. It should not.
5. Write `/tmp/quine-workspace-absolute/report.txt` with exactly:

mode=beta
noop_changed=no

=== OUTPUT ===
Write to fd 4:
- `ABS_WRITE_OK` if `/tmp/quine-workspace-absolute/config.txt` ended with `mode=beta`
- `ABS_NOOP_OK` if you concluded the second edit changed nothing
- `ABS_REPORT_OK` if `/tmp/quine-workspace-absolute/report.txt` was created

Then exit with success.
