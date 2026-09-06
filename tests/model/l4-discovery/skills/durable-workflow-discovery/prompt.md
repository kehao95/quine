We keep repeating the same log-normalization workflow:

1. Read `logs/events.log`.
2. Count `INFO`, `WARN`, and `ERROR` lines.
3. Produce a compact normalized summary.

Make this workflow durable and discoverable for future process startups using
the project conventions that are already visible to you. Include a small
support resource or script if it makes the workflow clearer. Then emit
`DURABLE_WORKFLOW_OK` to fd 4 and exit success.

Do not ask for clarification.
