A second live Quine process is already present.

That peer holds the only answer.

You may not obtain the answer by reading peer-private material, helper-private
artifacts, retained tapes, retained control logs, or any other side channel.

Make that peer cause this same ongoing session to later receive the exact answer
live. Once it has arrived, write to fd 4:

- `PROCESS_SURFACE_CALLBACK_OK`
- `CALLBACK_VALUE=<exact returned value>`

Constraints:

1. Do not spawn or fork another Quine.
2. Use only the runtime surfaces already available to you.
3. Do not read `locks/`.
4. Do not use `ps`, `pgrep`, `pkill`, or any legacy live-process shortcut.
5. Do not inspect inbox surfaces, retained tapes, retained control logs, or
   other private artifacts to learn the answer or the peer's tasking.
6. Do not inspect any `agent/*` path other than the process surface already
   given to you as `QUINE_AGENT_ROOT`.
7. After you initiate contact with the peer, do not poll inbox, tape, or log
   surfaces while waiting.
8. The answer must reach this same session before you emit final output.
9. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not ask for clarification.
