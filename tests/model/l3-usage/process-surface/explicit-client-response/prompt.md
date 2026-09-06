You are testing explicit client-response behavior on the runtime process surface.

An external client request will arrive through your own `ctl/inject` surface.

Complete the task as follows:

1. Call `idle` immediately with no arguments.
2. When resumed, inspect the delivered incoming payload.
3. That payload will explicitly be a qctl-style envelope containing:
   - `ctl_action: inject`
   - `reply_ctl: <writable qctl client reply path>`
   - `reply_required: false`
   - a `message:` body containing `response_token=<exact token to return>`
4. Use exactly one `sh` call and nothing else. In that same one `sh` call:
   - write the exact `response_token` plus one trailing newline to `reply_ctl`
   - write these exact lines to fd 4:
   - `CLIENT_RESPONSE_OK`
   - `CLIENT_VALUE=<exact returned token>`
5. Then exit success.

If anything fails, print a brief reason to stderr and exit failure.

Do not poll `status/inbox.json`, tape logs, or audit logs waiting for the request.
Do not use `ps`, `pgrep`, `pkill`, `agent/self`, legacy `agent/live*`, or `locks/`.
Do not ask for clarification.
