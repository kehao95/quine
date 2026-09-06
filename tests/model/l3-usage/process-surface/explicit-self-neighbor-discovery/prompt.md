You are testing the explicit runtime process-surface contract.

Complete both parts and write success markers to fd 4.

=== PART 1: CREATE ONE LIVE NEIGHBOR ===
Use exactly one `fork(mode="forget")` call with one child.
Pass the child intent **verbatim** as this exact string (do not paraphrase it):

`Call idle first with no arguments and nothing else. After idle resumes, exit success.`

- child scope must be `.`

The child only exists to give you one temporary live neighbor while you inspect the runtime surface.

=== PART 2: DISCOVER SELF + LIVE NEIGHBOR ===
After the fork returns, use exactly one `sh` call and nothing else.

Inside that call, run one inline `python3` script. Prefer the heredoc form
`python3 - <<'PY'` inside that one shell call rather than `python3 -c ...` so
you do not burn the one allowed `sh` call on quoting errors. That script must:

1. Read `$QUINE_AGENT_ROOT/status/session.json`.
2. Let `runtime_root` come from `status/session.json["runtime_root"]`.
3. Use directory entry names from `os.listdir(runtime_root / "pid")`.
   Use `os.readlink(runtime_root / "pid" / <pid>)` to recover the published
   `agent/<session>/public` target for that live pid entry, then recover the
   corresponding session/agent root from that target.
   The target itself names the peer's exported `public` directory, so the
   agent root is that target's parent directory and the session id is that
   parent directory's basename.
   The symlink target may be absolute or relative; handle either form.
4. Poll for up to 40 seconds until `runtime_root/pid` contains at least one pid besides your own.
5. Verify all of the following:
    - `$QUINE_AGENT_ROOT/mission.txt` exists and is non-empty
    - `status/session.json` contains a non-empty `session_id`
    - `status/session.json` contains `pid > 0`
    - `status/session.json.agent_root` exactly equals `$QUINE_AGENT_ROOT`
    - `runtime_root/pid/<pid>` exists
    - `runtime_root/pid/<pid>` resolves exactly to `$QUINE_AGENT_ROOT/public`
    - at least one live neighbor besides yourself is visible in `runtime_root/pid`
    - at least one neighbor session besides yourself can be recovered from a
      `pid/<pid>` symlink target
6. If everything is correct, write these exact markers to fd 4 on separate lines:
    - SELF_DISCOVERY_OK
    - PID_INDEX_OK
    - PID_ROUTING_OK
7. Then exit success.

If anything fails, print a clear reason to stderr and exit failure.
Do not run `ps`.
Do not read `locks/`.
Do not read any legacy `agent/live` or `agent/live_by_pid` shortcut.
Do not issue a second `fork`.
Do not make extra diagnostic `sh` calls.

Do not ask for clarification.
