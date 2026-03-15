You have 3 independent tasks. Complete all of them and report results to fd 4.

=== TASK 1: GENERATE AND SAVE A REPORT ===
Generate the following report content and save it to /tmp/quine-emergent-test/report.txt:

  title: "Q4 Results"
  revenue: $1,200,000
  path: C:\reports\q4\final
  notes: Use "net" not "gross"

The file must contain those lines exactly — including the dollar sign, backslash,
and quotation marks.
Verify the saved file matches by reading it back.

=== TASK 2: TEXT STATISTICS ===
Count how many lines are in this block of text:

  alpha
  beta
  gamma
  delta
  epsilon

Report the line count.

=== TASK 3: CODE EVALUATION ===
Evaluate the following Python snippet and capture its output:

  result = 0
  for i in range(1, 6):
      result += i
  print(f"SUM={result}")

Report what SUM equals.

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines:
- REPORT_OK     (Task 1: file saved with all special characters intact)
- COUNT_OK      (Task 2: line count was correct — write "COUNT=5" on its own line too)
- SUM_OK        (Task 3: Python output was captured — write "SUM=15" on its own line too)

Example output if all pass:
REPORT_OK
COUNT=5
COUNT_OK
SUM=15
SUM_OK

Then exit with success.
