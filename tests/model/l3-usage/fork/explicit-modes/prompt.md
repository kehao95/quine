You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 4.

Use the current fork interface: one fork call per part, with explicit `children`
objects and explicit `mode`.

For every child below, use `scope="."`.

=== PART 1: GATHER-ALL / WAIT MODE ===
Use exactly one fork call with `mode="wait"` and two children.

Child 0 intent:
- use a shell command to write `ALPHA_RESULT` to fd 4
- then exit success

Child 1 intent:
- use a shell command to write `BETA_RESULT` to fd 4
- then exit success

After the fork call returns, verify from the returned fork result that:
- `mode` is `wait`
- both children completed successfully
- child stdout/result fields contain both `ALPHA_RESULT` and `BETA_RESULT`

=== PART 2: RACE MODE ===
Use exactly one fork call with `mode="race"` and two children.

Child 0 intent:
- use a shell command that sleeps for about 5 seconds and then writes `SLOW_LOSER` to fd 4
- then exit success

Child 1 intent:
- use a shell command to write `FAST_WINNER` to fd 4 immediately
- then exit success

After the fork call returns, verify from the returned fork result that:
- `mode` is `race`
- `winner.index` is the fast child
- the winner completed with exit code 0
- the winner payload/stdout contains `FAST_WINNER`
- `succeeded` is `1`
- `killed` is `1`

Do not rely on old textual swarm headers like "won"; use the structured fork
result fields.

=== PART 3: SINGLE CHILD ===
Use exactly one fork call with `mode="wait"` and one child.

Child intent:
- use a shell command to write `SINGLE_CHILD` to fd 4
- then exit success

After the fork call returns, verify from the returned fork result that:
- exactly one child ran
- that child completed successfully
- its stdout/result contains `SINGLE_CHILD`

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- GATHER_OK
- RACE_OK
- SINGLE_OK

Then exit with success.
