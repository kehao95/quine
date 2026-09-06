You are testing self-rebuild plus explicit exec handoff.

Use the carried self-source tree under `QUINE_AGENT_ROOT/source-code/` to rebuild
your own binary into `QUINE_AGENT_ROOT`, then replace the current process with
that rebuilt binary.

Requirements:
- Before the handoff, use only `sh` and `exec`. Do not call `exit` before the
  handoff.
- Build the `./cmd/quine` program from the module rooted at
  `QUINE_AGENT_ROOT/source-code/`, and materialize the executable as
  `QUINE_AGENT_ROOT/rebuilt-quine`.
- The build must originate from `QUINE_AGENT_ROOT/source-code/`; do not use
  `/tmp/quine` or ordinary repo-root shortcuts like `./cmd`, `./internal`,
  `./go.mod`, or `./go.sum` as the source of truth.
- Use `exec` exactly once with an explicit absolute
  `target="$QUINE_AGENT_ROOT/rebuilt-quine"` and explicit `argv`.
- The `exec` tool does not expand shell variables in `target`; give it the
  fully expanded absolute path, not a literal `$QUINE_AGENT_ROOT/...` string.
- The replacement process must be quine again, not a shell wrapper.
- Pass this exact successor mission as `argv[1]`:

Do exactly one thing: call the `exit` tool with `status="success"`.

- Do not return control to the original process after exec.
- Do not ask for clarification.
