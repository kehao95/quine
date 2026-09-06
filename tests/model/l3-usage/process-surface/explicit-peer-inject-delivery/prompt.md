You are testing the explicit peer-inject delivery contract on the runtime process surface.

Complete both parts and write success markers to fd 4.

=== PART 1: CREATE ONE LIVE NEIGHBOR ===
Use exactly one `fork(mode="forget")` call with one child.
Pass the child intent **verbatim** as this exact string:

`The fork that created you does not count as one of your tool actions. Take exactly two tool actions and no plain-text replies: first call idle with no arguments; after idle resumes, call exit with status success.`

- child scope must be `.`

The child only exists to give you one temporary live neighbor while you inspect
the runtime surface and deliver one inject event.

=== PART 2: DISCOVER SELF + LIVE NEIGHBOR + DELIVER ONE INJECT ===
After the fork returns, use exactly one `sh` call and nothing else.

Inside that call, use shell and/or inline Python as you prefer.
If you write deliverables from Python, use the numeric fd 4 descriptor directly
(for example `os.write(4, ...)` or `os.fdopen(4, 'w')`).

1. Read `$QUINE_AGENT_ROOT/status/session.json`.
2. Let `runtime_root` come from `status/session.json["runtime_root"]`.
3. Use directory entry names from `os.listdir(runtime_root / "pid")`.
   Use `os.readlink(runtime_root / "pid" / <pid>)` to recover the published
   `agent/<session>/public` target for that live pid entry, then recover the
   session id and agent root from that target.
   The target itself names the peer's exported `public` directory, so the
   agent root is that target's parent directory and the session id is that
   parent directory's basename.
   The symlink target may be absolute or relative; handle either form.
4. Poll for up to 40 seconds until `runtime_root/pid` contains at least one pid
   besides your own.
5. Recover:
   - the neighbor pid
   - the neighbor session id
   - the neighbor agent root
6. Do not treat the pid returned by `fork` as authoritative. If more than one
   non-self pid appears, poll each candidate's live cognition surface
   (`context/state/current.jsonl`) and choose the candidate that records the
   child calling `idle`.
7. Poll candidate live cognition surfaces for up to 30 seconds until one records
   the neighbor calling `idle`.
8. Build one non-empty single-line ASCII payload of your choice.

9. Write the exact payload plus one trailing newline to `ctl/inject` on the
   neighbor's exported public process surface.
10. Poll the neighbor surfaces for up to 5 seconds until:
   - `neighbor_root/status/inbox.json` shows `pending_count == 0`
   - the neighbor's live cognition surface (`context/state/current.jsonl`) and/or
     retained control log (`log/<session>/control.jsonl`) records the payload
     with delivery `"inject"`
11. If everything is correct, write these exact lines to fd 4 on separate lines:
   - SELF_DISCOVERY_OK
   - NEIGHBOR_DISCOVERY_OK
   - CTL_WRITE_OK
   - INJECT_TRANSACTION_OK
   - PEER_INJECT_OK
   - SELF_SESSION=<your session id>
   - SELF_PID=<your pid>
   - NEIGHBOR_SESSION=<neighbor session id>
   - NEIGHBOR_PID=<neighbor pid>
   - PEER_MESSAGE=<exact payload you wrote>
   - DELIVERY=inject
12. Then exit success.

If anything fails, print a clear reason to stderr and exit failure.

Do not run `ps`.
Do not read `locks/`.
Do not read any legacy `agent/live` or `agent/live_by_pid` shortcut.
Do not issue a second `fork`.
Do not make extra diagnostic `sh` calls.

Do not ask for clarification.
