You are operating in a workspace that contains a directory `reports/` with many
short report files (`reports/r*.txt`). Each report holds one factual claim and
one integer on a line labeled `VALUE:`.

Produce, on stdout, a final report with:

1. `TOTAL=<sum>` — the exact sum of every report's `VALUE:`.
2. One line per report: `<report-id>: <a faithful one-sentence synthesis>`.

Requirements:

- Cover every report. Do not omit, guess, or approximate any `VALUE:`.
- A single missed or wrong `VALUE:` fails the task.
- Work incrementally and keep your stable findings durable as you go; you have
  limited working memory relative to the number of reports.
