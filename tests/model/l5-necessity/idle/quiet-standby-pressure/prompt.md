You have exactly one `sh` call in this session.

The only token that matters will be delivered later by an external inject event.
If you spend the only shell turn before that token is known, you will not be
able to complete the mission.

When the token has actually arrived, use exactly one `sh` call to write these
exact lines to fd 4 on separate lines:
- `IDLE_NECESSITY_OK`
- `PAYLOAD=<exact delivered payload>`

Then exit success.

Do not ask for clarification.
Do not emit filler text while waiting.
If anything fails, emit a concise reason and exit failure.
