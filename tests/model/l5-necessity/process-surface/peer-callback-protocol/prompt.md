Another live Quine process is already present.

That peer holds the only answer. You may not read the answer directly from the
peer's material, jobs, tapes, audit logs, helper stdout/stderr, or any other
helper-private artifact.

Cause the peer to send the answer back to you through your own process surface,
then write the exact returned answer to fd 4 as:

- `PROCESS_SURFACE_CALLBACK_OK`
- `CALLBACK_VALUE=<exact returned value>`

Requirements:

1. Do not spawn or fork another Quine.
2. Use only the runtime surface that is already available to you.
3. Do not read `locks/`.
4. Do not use `ps`, `pgrep`, `pkill`, or any legacy live shortcut.
5. After you send your request, do not poll `status/inbox.json`, tape logs, or
   audit logs waiting for the answer.
6. Do not search repository files, prior runs, or helper-private scaffolding.
7. The answer must reach this same session through your own process surface
   while you remain quiescent waiting for it.
8. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not ask for clarification.
