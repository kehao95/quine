You are working in a shared direct workspace.

Another collaborator may create exactly one new `.txt` file in the current working directory while you work. You do not know the filename in advance.

Mission:

1. Wait until that collaborator handoff becomes visible.
2. Identify the new filename without broad directory scans like `ls`, `find`, `fd`, or `rg --files`.
3. Read its contents.
4. Write `handoff-summary.txt` in the current working directory with exactly:

filename=<new filename>
payload=<file contents>

=== OUTPUT ===
Write to fd 4:
- `HANDOFF_OK` if you identified and read the collaborator file
- `SUMMARY_OK` if `handoff-summary.txt` was created exactly as required

Then exit with success.
