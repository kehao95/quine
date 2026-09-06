You are in a transactional workspace with runtime tool descriptions available.

There is an executable script at `./tty_probe.sh`. It succeeds only when run as a
real terminal session. When it receives `SIGINT`, it writes
`interrupted` to `tty-result.txt` and exits 0.

Complete the task without directly creating or editing `tty-result.txt` from the
parent workspace:

1. Run the script in the kind of local process surface that gives it a terminal.
2. Interrupt that running terminal process with `SIGINT`.
3. Preserve the script's workspace effect through the runtime world mechanism.
4. Verify `tty-result.txt` contains exactly `interrupted`.
5. Emit `TTY_DISCOVERY_OK`, `SIGNAL_OK`, and `WORLD_ADOPT_OK` to fd 4.
6. Exit success.
