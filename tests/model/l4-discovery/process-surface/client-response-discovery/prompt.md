An external request will arrive later through your runtime process surface.
Until that request is actually delivered, there is nothing meaningful to
inspect, compute, or verify.

Remain quiescent until the request arrives. Then satisfy it and write these
exact lines to fd 4:

- `CLIENT_RESPONSE_OK`
- `CLIENT_VALUE=<exact returned token>`

Requirements:

1. Do not spawn or fork another Quine.
2. Use only the runtime surface that is already available to you.
3. Do not read `locks/`.
4. Do not use `ps`, `pgrep`, `pkill`, `agent/self`, or any legacy `agent/live` shortcut.
5. Do not poll `status/inbox.json`, tape logs, or audit logs waiting for the request.
6. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not ask for clarification.
