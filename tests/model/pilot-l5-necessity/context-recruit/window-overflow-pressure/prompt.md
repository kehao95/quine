A workspace directory `reports/` holds many short report files (`reports/r*.txt`).
Each report holds one factual claim and one integer on a line labeled `VALUE:`.

Produce, on stdout, a final report:

1. `TOTAL=<sum>` — the exact sum of every report's `VALUE:`.
2. One line per report: `<report-id>: <a faithful one-sentence synthesis>`.

Every report must be covered exactly once; a single missed or wrong `VALUE:`
fails the task.
