You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 4.

The sh tool supports a "stdin" parameter: sh(command="...", stdin="...").
Use this to pipe data into a command without shell quoting — the content is
passed verbatim as the process's standard input.

=== PART 1: WRITE A FILE WITH SHELL-HOSTILE CONTENT ===
Create the directory /tmp/quine-stdin-test/, then write a file
/tmp/quine-stdin-test/config.ini with EXACTLY this content
(including the quotes, backslash, and dollar sign):

  api_key = "abc123"
  base_path = C:\data\logs
  greeting = $HOME is your home

Use sh(command="cat > /tmp/quine-stdin-test/config.ini", stdin="<the content above>")
Do NOT use echo or printf with shell escaping — use the stdin parameter.

Then verify the file was written correctly:
  sh(command="cat /tmp/quine-stdin-test/config.ini")
Confirm all three lines appear verbatim (quotes, backslash, and dollar sign intact).

=== PART 2: PIPE DATA TO A PROGRAM ===
Count the number of words in a sentence using wc:
  sh(command="wc -w", stdin="the quick brown fox jumps over the lazy dog\n")
The answer should be 9.
Confirm you see "9" in the output.

=== PART 3: MULTI-LINE SCRIPT EXECUTION ===
Pipe a Python script to the interpreter via stdin:
  sh(command="python3 -", stdin="import sys\nfor i in range(1, 4):\n    print(f'LINE_{i}')\n")
Confirm the output contains LINE_1, LINE_2, and LINE_3.

=== CLEANUP ===
Remove the test directory:
  sh(command="rm -rf /tmp/quine-stdin-test")

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- WRITE_OK      (Part 1: config.ini written with all special chars intact)
- WORD_OK       (Part 2: wc -w via stdin returned 9)
- SCRIPT_OK     (Part 3: python3 - via stdin printed LINE_1 through LINE_3)

Example (if all pass):
WRITE_OK
WORD_OK
SCRIPT_OK

Then exit with success.
