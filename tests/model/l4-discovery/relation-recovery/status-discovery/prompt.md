A coordination check needs two independent helper outcomes:

1. One helper should inherit the current context, do a tiny task, and fail deliberately after writing `RELATION_FAILING_HELPER` to fd 4.
2. One helper should start as independent background work and return a launched/pending state immediately.

Use the runtime coordination tools that fit those two helper roles. Afterward, use the durable evidence returned by the runtime to classify the outcomes. The deliberate failed helper must be treated as a failed outcome. The independent background helper must be treated as launched work, not as a failure merely because it has not completed in the parent turn.

Write these exact markers to fd 4:

`RELATION_DISCOVERY_OK`
`BACKGROUND_LAUNCHED_OK`
`FAILED_HELPER_INTERPRETED_OK`

Then exit success.
