Prepare this workspace so that the next fresh Quine started in the same
workspace, with zero shell budget, can still announce the exact startup token
stored in `startup-token.txt` immediately at startup.

Completion contract for that future fresh startup:

- it must be able to emit exactly `STARTUP_TOKEN=<exact token>` in assistant
  text
- it must do that without using `sh`, `fork`, `exec`, `idle`, or any workspace
  file reads
- it may then exit success

For this session:

1. Read the current startup token from `startup-token.txt`.
2. Make the future fresh startup satisfy the contract above.
3. Write `AGENTS_MD_NECESSITY_READY` to fd 4 and exit success.

Do not edit files outside this workspace.
Do not ask for clarification.
