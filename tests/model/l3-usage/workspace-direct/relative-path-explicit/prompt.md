You are running with direct shared workspace physics enabled.

There is a file called `config.txt` in the current working directory. Complete all parts:

1. Read `config.txt`.
2. Change the line `mode=alpha` to `mode=beta`.
3. Then run a second edit attempt that tries to replace a string that does not exist:
   change `mode=gamma` to `mode=delta`.
4. Use the `[FS MUTATIONS]` channel from the second command to determine whether that second edit actually changed the visible workspace. It should not.
5. Write a file `report.txt` in the current working directory with exactly:

mode=beta
noop_changed=no

=== OUTPUT ===
Write to fd 4:
- `WRITE_OK` if `config.txt` ended with `mode=beta`
- `NOOP_OK` if you concluded the second edit changed nothing
- `REPORT_OK` if `report.txt` was created

Then exit with success.
