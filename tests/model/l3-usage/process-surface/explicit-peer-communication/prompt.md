You are testing the explicit peer-communication contract on the runtime process surface.

Complete both parts and write success markers to fd 4.

=== PART 1: CREATE ONE LIVE NEIGHBOR ===
Use exactly one `fork(mode="forget")` call with one child.
Pass the child intent **verbatim** as this exact string (do not paraphrase it):

`The fork that created you does not count as one of your tool calls. Use exactly one sh tool call and no other tool calls except a final exit success. In that single sh call, run command sleep 240. Do not use timeout. Do not use detach. Wait for it normally. Then call exit with status success.`

- child scope must be `.`

The child only exists to give you one temporary live neighbor while you inspect
the runtime surface and send one peer message.

=== PART 2: DISCOVER SELF + LIVE NEIGHBOR + SEND ONE MESSAGE ===
After the fork returns, use exactly one `sh` call and nothing else.

Inside that call, run one inline `python3` script. Prefer the heredoc form
`python3 - <<'PY'` inside that one shell call rather than `python3 -c ...` so
you do not burn the one allowed `sh` call on quoting errors. That script must:

Use the numeric fd 4 descriptor directly when writing deliverables from Python
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
6. Build one non-empty single-line ASCII payload of your choice.

7. Write the exact payload plus one trailing newline to `ctl/post` on the
   neighbor's exported public process surface.
8. Poll `neighbor_root/status/inbox.json` for up to 5 seconds until:
   - `pending_count >= 1`
   - one inbox message payload exactly matches the payload you wrote
9. If everything is correct, write these exact lines to fd 4 on separate lines:
   - SELF_DISCOVERY_OK
   - NEIGHBOR_DISCOVERY_OK
   - CTL_WRITE_OK
   - PEER_INBOX_OK
   - SELF_SESSION=<your session id>
   - SELF_PID=<your pid>
   - NEIGHBOR_SESSION=<neighbor session id>
   - NEIGHBOR_PID=<neighbor pid>
   - PEER_MESSAGE=<exact payload you wrote>
10. Then exit success.

If anything fails, print a clear reason to stderr and exit failure.

Do not run `ps`.
Do not read `locks/`.
Do not read any legacy `agent/live` or `agent/live_by_pid` shortcut.
Do not issue a second `fork`.
Do not make extra diagnostic `sh` calls.

Do not ask for clarification.
