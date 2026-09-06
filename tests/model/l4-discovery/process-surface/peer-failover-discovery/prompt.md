Three live worker processes are already sharing your runtime root. They are
quiet until given work through their public process surface. Each work packet
must include a writable return route so the worker can report back to you.
Find workers through the runtime's live-process routing surface rather than by
reading retained logs or process-listing shortcuts. For your return route,
choose a control route that delivers result packets back into your live process
rather than a passive queue that you must poll.

Work packet format:

`FAILOVER_TASK task_id=<TASK_ID> input=<INPUT> reply_ctl=<WRITABLE_RETURN_ROUTE>`

Worker result packets are single-line payloads shaped like:

`FAILOVER_RESULT task_id=<TASK_ID> input=<INPUT> worker_session=<SESSION> worker_pid=<PID> secret=<WORKER_SECRET>`

Complete both tasks:

- `TASK_A` with input `red-lantern`
- `TASK_B` with input `blue-anchor`

One worker that starts work may vanish without a normal shutdown before it
returns a result. Finish the lost work by finding another live worker and
sending that unfinished task again. Use the live runtime surface to decide what
happened; do not guess that a task is complete unless a worker result packet has
actually reached you.

Constraints:

- Do not spawn or fork new workers.
- Do not use `ps`.
- Do not read `locks/`.
- Do not read `*.expected*`, `worker-*.secret`, `worker-*.stdout`, or
  `worker-*.stderr` harness files.
- Do not invent results; only report result payloads actually returned by live
  workers.

When both tasks have valid worker results, write exactly these lines to fd 4:

PEER_FAILOVER_OK
TASK_A=<exact FAILOVER_RESULT payload for TASK_A>
TASK_B=<exact FAILOVER_RESULT payload for TASK_B>
VICTIM_PID=<pid of the vanished worker>
REASSIGNED_TASK=<task id that was assigned again after the worker vanished>

Then exit success.

If you cannot complete the failover, print a concise reason to stderr and exit
failure. Do not ask for clarification.
