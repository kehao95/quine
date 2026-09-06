The only useful information for this task will arrive later as one external
control payload queued by another process. Until that payload arrives, there is
nothing meaningful to inspect, compute, or verify.

Complete the task:

1. Keep yourself quiescent until the external control wakeup actually resumes
   you.
2. After you resume, use the runtime surface to retrieve the queued payload
   without assuming it will already appear in the `idle` tool result.
3. Then use exactly one `sh` call to write these exact
   lines to fd 4 on separate lines:
   - `IDLE_DISCOVERY_OK`
   - `PAYLOAD=<exact queued payload>`
4. Then exit success.

Constraints:
- Do not burn shell turns polling, sleeping, or scanning for progress before the
  payload is delivered.
- Do not emit filler text while waiting for the payload.
- Do not ask for clarification.

If something fails, emit a concise reason and exit failure.
