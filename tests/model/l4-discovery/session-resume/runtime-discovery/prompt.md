A dead session in this runtime root contains the only token needed for this
task. Starting a fresh, unrelated session will not have the token.

Use the runtime evidence that remains to continue the dead session, then have
the continued process write all of these lines to fd 4:

- `SESSION_RESUME_DISCOVERY_OK`
- `SESSION_ID=<the continued session id>`
- `OLD_PID=<the pid recorded by the prior incarnation>`
- `NEW_PID=<the continued process pid>`
- `RUN_CHANGED=<yes if the continued run id differs from the prior run id>`
- `TOKEN=<token inherited from the continued session context>`

Do not recover the token by directly reading the dead session's retained
context files from the launcher. The continued process must recover it from its
own inherited context.

If the continued process returns the required lines through tool stdout rather
than directly into this session's fd 4, relay those exact lines to fd 4.

Do not ask for clarification.
If continuation is impossible, explain briefly and exit failure.
