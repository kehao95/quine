You are running with transactional workspace physics enabled.

There is a file at `__ABS_WORKSPACE__/config.txt`. Complete all parts using explicit absolute paths:

1. Read `__ABS_WORKSPACE__/config.txt`.
2. Change the line `mode=alpha` to `mode=beta`.
3. Then run a second edit attempt that tries to replace a string that does not exist:
   change `mode=gamma` to `mode=delta`.
4. Use the `[FS MUTATIONS]` channel from the second command to determine whether that second edit actually changed the filesystem. It should not.
5. Write `__ABS_WORKSPACE__/report.txt` with exactly:

mode=beta
noop_changed=no

=== OUTPUT ===
Write to fd 4:
- `ABS_WRITE_OK` if `__ABS_WORKSPACE__/config.txt` ended with `mode=beta`
- `ABS_NOOP_OK` if you concluded the second edit changed nothing
- `ABS_REPORT_OK` if `__ABS_WORKSPACE__/report.txt` was created

Then exit with success.
