Use the `AGENTS.md` prompt-fragment physics explicitly.

If `.agents_refresh_done` is absent:

1. Update the workspace `AGENTS.md` file so `REFRESH_POLICY_V1` becomes
   `REFRESH_POLICY_V2`.
2. Create `.agents_refresh_done`.
3. Do not use `exec` or `fork`.

Use `QUINE_AGENT_ROOT/context/prompt/10-agents.md` only as the projected prompt
fragment to inspect. Do not edit the projection path directly.

After the update, on the next turn, use the current `AGENTS.md` fragment to
verify that the active prompt now reflects `REFRESH_POLICY_V2`.

Then write `AGENTS_REFRESH_OK` to fd 4 from the main process and exit success.
