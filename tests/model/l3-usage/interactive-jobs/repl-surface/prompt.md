You have a PTY-backed interactive job surface available through `sh(command="...", interactive=true)`.
Complete ALL parts and output a structured report to fd 4.

This scenario tests the screen-first interactive job model. There is no special `job(...)` tool for this flow. Use the absolute job directory path returned by `sh(..., interactive=true)`.

=== PART 1: OBSERVE THE SCREEN ===
Start an interactive Python REPL:
  sh(command="python3 -q", interactive=true)
You should receive a `[JOB] pid=N path=/absolute/path/... (interactive)` result immediately.
Reuse that returned `path` value in later commands.

Confirm the directory contains:
- `pid`
- `started_at`
- `screen.txt`
- `screen.meta`
- `in`
- `winsize`
- `events.log`

Then read `screen.txt` and confirm it shows the Python prompt `>>>`.

=== PART 2: ACT THROUGH THE INPUT SURFACE ===
Write this exact content to `<PATH>/in`:
  print(6*7)<enter>

Then read `screen.txt` until it contains `42`.

=== PART 3: RESIZE THE VIEWPORT ===
Write this exact content to `<PATH>/winsize`:
  100x30

Then read `screen.meta` and confirm it reports:
- `cols = 100`
- `rows = 30`

=== PART 4: CLEAN EXIT ===
Write this exact content to `<PATH>/in`:
  exit()<enter>

When the REPL finishes, its exit code will appear at `<PATH>/exit`.
Confirm the exit code is `0`.

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- SCREEN_OK       (Part 1: interactive job path exposed screen files and prompt appeared)
- INPUT_OK        (Part 2: writing to `<path>/in` caused `42` to appear on screen)
- RESIZE_OK       (Part 3: writing to `<path>/winsize` updated `screen.meta`)
- EXIT_OK         (Part 4: `<path>/exit` recorded 0 after clean REPL shutdown)
- INTERACTIVE_OK  (all interactive job parts completed)

Example (if all pass):
SCREEN_OK
INPUT_OK
RESIZE_OK
EXIT_OK
INTERACTIVE_OK

Then exit with success.
