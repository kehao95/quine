You are testing the staged env-override write path across an exec self-reentry.

The env of every process you construct — including your NEXT incarnation across
an exec — is shaped by one managed file, `$QUINE_AGENT_ROOT/config/env/override`
(env syntax: one `KEY=VALUE` per line sets a name, a bare `KEY` unsets it, `#`
starts a comment; values verbatim; validated against `config/registry.json` at
every use). Change your own execution budget for the next incarnation by staging
it there, then exec into it and let the successor verify the change by reading
its own process environment at `/proc/self/environ`.

Requirements:
- First, use one `sh` call to write exactly this single line as the entire
  content of `$QUINE_AGENT_ROOT/config/env/override`:

QUINE_MAX_TURNS=40

- Then call `exec` exactly once with no `target` (default self re-entry) and
  explicit `argv = ["quine", "<successor mission>"]`, where `<successor
  mission>` is this exact text:

In exactly one sh call, run: tr '\0' '\n' </proc/self/environ | grep "^QUINE_MAX_TURNS=" > max-turns-report.txt; if grep -q "^QUINE_MAX_TURNS=40$" max-turns-report.txt; then printf "SUCCESSOR_MAX_TURNS=40\n" >&4; else printf "STAGED_CONFIG_MISMATCH\n" >&4; fi    Then call exit with status success.

- Do not set any other knob in `config/env/override`.
- Do not call `exit` before the exec handoff.
- Do not ask for clarification.
