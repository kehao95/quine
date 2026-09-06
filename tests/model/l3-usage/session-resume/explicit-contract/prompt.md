Use the runtime session-resume contract explicitly.

There is a retained session named `sess-resume-explicit-target` in this
runtime root. It has already exited, but its retained context contains the
token needed for this task.

Resume that exact session by launching a new Quine process with:

- `QUINE_DATA_DIR=<this same runtime root>`
- `QUINE_RETENTION_DIR=<this same runtime root>/log`
- `QUINE_SESSION_ID=sess-resume-explicit-target`
- the current Quine binary

Give the resumed process a mission that tells it to use the values already
visible in its inherited Memory prompt:

- `TOKEN=...`
- `OLD_PID=...`
- `OLD_RUN_ID=...`

The resumed process can compare `OLD_RUN_ID` with its current
`status/session.json.run_id` to decide `RUN_CHANGED`.
Inside shell tools, `status/session.json` is not relative to the workspace;
the resumed process should read its current status from
`$QUINE_AGENT_ROOT/status/session.json`.

The resumed process must emit all of these lines to fd 4 and exit success:

- `SESSION_RESUME_EXPLICIT_OK`
- `SESSION_ID=<the resumed session id>`
- `OLD_PID=<the pid recorded by the prior incarnation>`
- `NEW_PID=<the resumed process pid>`
- `RUN_CHANGED=<yes if the resumed run id differs from the prior run id>`
- `TOKEN=<token inherited from the resumed session context>`

Do not read the retained session's `context/prompt/30-memory.md` directly from the
launcher session. The token should be recovered by the resumed Quine process
from its inherited context.

If the shell invocation returns the resumed process's deliverable lines as
tool stdout instead of forwarding them to this session's fd 4, relay those
exact lines to fd 4 from the launcher. Do not recover missing fields by reading
the retained memory file.
Relay only a complete output that already contains a non-empty `SESSION_ID`,
`OLD_PID`, `NEW_PID`, `RUN_CHANGED=yes`, and `TOKEN`.

Do not ask for clarification.
