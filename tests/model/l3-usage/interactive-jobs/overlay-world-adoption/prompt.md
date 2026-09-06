You are testing explicit interactive-job workspace adoption.

The current working directory is an empty transactional workspace.

Complete this protocol exactly:

1. Use exactly one `sh` call with `interactive=true`.
2. The interactive command must print `READY` to its screen, then wait until it receives `SIGINT`.
3. When the command receives `SIGINT`, it must write exactly `interrupted` to `interactive-result.txt` and exit 0.
4. Use the returned job path and its `pid` file to send a POSIX `SIGINT` to the job process group with an ordinary `sh` command.
5. Wait until the job records exit code `0`.
6. Confirm the job directory exposes a `world_handle`.
7. Before switching worlds, confirm the parent workspace does not contain `interactive-result.txt` and emit `PRE_SWITCH_PRIVATE` and `SIGNAL_OK` to fd 4.
8. Call `switch_world` with the exact handle from `<job>/world_handle`.
9. Confirm `interactive-result.txt` now contains exactly `interrupted`, then emit `INTERACTIVE_WORLD_OK` and `HANDLE_OK` to fd 4.
10. Exit success.

Do not use `fork`, `spawn`, or `exec`.
