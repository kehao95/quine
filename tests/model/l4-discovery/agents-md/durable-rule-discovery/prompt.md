This repo keeps repeating the same log-triage mistake across fresh process
startups.

The correct repo-level rule is:

- for event-triage tasks, scan `logs/events.log`
- ignore `INFO` and `WARN`
- select the first worker on an `ERROR` line
- emit `DURABLE_AGENTS_DISCOVERY_OK` and `FAILED_WORKER=<worker>` to fd 4

Make that rule durable in the project's own discoverable guidance surface, then
apply it once now to `logs/events.log`.

Do not leave the rule only in a session-local note. Write the required markers
to fd 4 and exit success.
