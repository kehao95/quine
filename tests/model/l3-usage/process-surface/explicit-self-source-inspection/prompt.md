You are testing the explicit self-source process surface.

Use the carried self-source tree under `QUINE_AGENT_ROOT/source-code/`.

Complete the task exactly as follows:

1. Use exactly one `sh` call and no other work tool calls before your final success exit.
2. Inside that call:
   - set `root="$QUINE_AGENT_ROOT/source-code"`
   - verify `"$root"` exists and is read-only
   - verify these files exist and are non-empty:
      - `"$root/go.mod"`
      - `"$root/go.sum"`
      - `"$root/selfsource.go"`
      - `"$root/cmd/quine/main.go"`
      - `"$root/internal/runtime/runtime.go"`
    - verify:
      - `"$root/go.mod"` contains the exact module line `module github.com/kehao95/quine`
      - `"$root/cmd/quine/main.go"` contains `func main()`
      - `"$root/selfsource.go"` contains `var SelfSource embed.FS`
      - `"$root/internal/runtime/runtime.go"` contains `type Runtime struct`
3. If everything passes, write these exact lines to fd 4 on separate lines:
   - `SELF_SOURCE_USAGE_OK`
   - `MODULE_OK`
   - `MAIN_OK`
   - `RUNTIME_OK`
4. Then exit success.

If anything fails, print a clear reason to stderr and exit failure.

Do not read `./cmd`, `./internal`, `./go.mod`, or `./go.sum` outside
`QUINE_AGENT_ROOT/source-code/`.
Do not ask for clarification.
