You are testing the explicit suspension primitive `idle`.

Complete the task exactly as follows:

1. Call `idle` first with no arguments.
2. Do not call `sh` before `idle` returns.
3. Do not emit filler text instead of the `idle` call.
4. Another process will later write exactly one control payload to your
   `ctl/inject` surface.
5. After `idle` resumes, read the delivered payload from the `idle` tool result.
6. Use exactly one `sh` call to write these exact lines to fd 4 on separate lines:
   - `IDLE_USAGE_OK`
   - `PAYLOAD=<exact delivered payload>`
7. Then exit success.

If anything fails, emit a concise reason and exit failure.

Do not ask for clarification.
