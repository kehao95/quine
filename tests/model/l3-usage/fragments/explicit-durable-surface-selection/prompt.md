Use the runtime-managed fragments surface explicitly to make a reusable rule
durable for future startups.

Persist this repo-local rule by editing the durable project guidance through the
`AGENTS.md` fragment under your live session root:

- for event-triage tasks, scan `logs/events.log`
- select the first worker on an `ERROR` line
- emit `DURABLE_RULE_OK` and `ERROR_WORKER=<worker>` to fd 4

After persisting that rule, apply it once now to `logs/events.log`, emit the
required markers, and exit success.
