You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 3.

=== PART 1: NON-INTERACTIVE ISATTY ===
Run the following command WITHOUT interactive mode:
  sh(command="python3 -c \"import sys; print('tty' if sys.stdin.isatty() else 'notty')\"")
Confirm the output contains "notty".

=== PART 2: INTERACTIVE ISATTY ===
Run the same command WITH interactive mode:
  sh(command="python3 -c \"import sys; print('tty' if sys.stdin.isatty() else 'notty')\"", interactive=true, timeout=5)
Confirm the output contains "tty".

=== PART 3: INPUT INJECTION ===
Start an interactive Python REPL session with a short timeout:
  sh(command="python3 -i", interactive=true, timeout=2)
This will pause after 2 seconds with a [PAUSED] job ID. The output should show the Python >>> prompt.
Then inject input to evaluate an expression:
  job(id=<ID from above>, signal="cont", input="print('INJECT_OK')\n", timeout=3)
This should resume the REPL, execute the print statement, and pause again (or exit).
Confirm "INJECT_OK" appears in the output of the resumed job.
Then kill the job to clean up:
  job(id=<same ID>, signal="kill")

=== OUTPUT ===
Write a report to fd 3 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- NOTTY_OK      (Part 1: non-interactive mode shows isatty==false)
- TTY_OK        (Part 2: interactive mode shows isatty==true)
- INJECT_OK     (Part 3: input injection was echoed back in output)

Example (if all pass):
NOTTY_OK
TTY_OK
INJECT_OK

Then exit with success.
