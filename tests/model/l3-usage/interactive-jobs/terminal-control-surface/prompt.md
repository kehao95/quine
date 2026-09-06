You are testing the PTY-backed terminal-control surface.

The current working directory contains an executable `./terminal_control.py`.
It refuses non-terminal stdin, prints `READY_CONTROL` when its terminal is
ready, requires the terminal viewport to be exactly `64x12`, then reads the
line `MAGENTA-42` from that terminal. On success it writes
`control-result.txt` and renders terminal escape sequences before exiting 0.

Complete this protocol exactly:

1. Start exactly one interactive job for `./terminal_control.py`.
2. Use the returned job path. Confirm `screen.txt` contains `READY_CONTROL`.
3. Resize the job terminal to `64x12` through its filesystem control surface
   and confirm `screen.meta` reports `cols=64` and `rows=12`.
4. Send `MAGENTA-42<enter>` through the job input surface.
5. Wait for `<job>/exit` to contain `0`.
6. Confirm `input.log` contains `MAGENTA-42`.
7. Confirm `events.hex` contains the raw escape prefix `1b5b324a` from
   the program's clear-screen render.
8. Confirm `control-result.txt` contains exactly `MAGENTA-42 64x12`.
9. Emit these markers to fd 4, one per line:

`CONTROL_READY_OK`
`RESIZE_OK`
`INPUT_LOG_OK`
`EVENTS_HEX_OK`
`CONTROL_RESULT_OK`
`TERMINAL_CONTROL_OK`

Then exit success.
