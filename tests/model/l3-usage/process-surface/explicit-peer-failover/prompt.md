You are a coordinator sharing one runtime root with three already-running idle
Quine worker peers.

Use the explicit runtime process surface to coordinate them:

- Your own process root is `QUINE_AGENT_ROOT`; your identity is in
  `status/session.json`.
- The shared runtime root is two directories above `QUINE_AGENT_ROOT`.
- Live peer routes are under `pid/<pid>` and each pid route resolves to that
  peer's public root.
- Queue worker tasks by writing one non-empty line to the worker public
  `ctl/post` path.
- Ask workers to answer by including your own public `ctl/inject` path as the
  `reply_ctl` field in the task packet.
- Watch your own control surface for worker result packets and peer-offline
  notices. A peer may also disappear by having its `pid/<pid>` route removed.

Worker task packet format:

`FAILOVER_TASK task_id=<TASK_ID> input=<INPUT> reply_ctl=<YOUR_PUBLIC_CTL_INJECT>`

Workers answer with one single-line payload shaped like:

`FAILOVER_RESULT task_id=<TASK_ID> input=<INPUT> worker_session=<SESSION> worker_pid=<PID> secret=<WORKER_SECRET>`

Complete exactly these tasks:

- `TASK_A` with input `red-lantern`
- `TASK_B` with input `blue-anchor`

One worker that has already received a task will be killed abruptly with
`SIGKILL`. This is intentional. Keep enough dispatch state to know which task
was assigned to that worker. When its offline notice or vanished pid route is
observed before its result arrives, reassign that unfinished task to another
live worker that has not already completed a task.

Operational constraints:

- Use at least two workers before the failure; keep or recover a live spare for
  reassignment.
- Do not use `ps`.
- Do not read `locks/`.
- Do not read `*.expected*`, `worker-*.secret`, `worker-*.stdout`, or
  `worker-*.stderr` harness files.
- Do not invent results; only report result payloads actually returned by live
  workers through your control surface.

When both tasks have valid worker results, write exactly these lines to fd 4:

PEER_FAILOVER_OK
TASK_A=<exact FAILOVER_RESULT payload for TASK_A>
TASK_B=<exact FAILOVER_RESULT payload for TASK_B>
VICTIM_PID=<pid of the killed worker>
REASSIGNED_TASK=<task id that was reassigned after the worker died>

Then exit success.

If you cannot complete the failover, print a concise reason to stderr and exit
failure. Do not ask for clarification.
