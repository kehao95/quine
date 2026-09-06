You need to analyze four independent datasets that live in four sibling directories in the current working directory:

- `sales/data.csv`
- `words/passage.txt`
- `temps/data.csv`
- `logs/server.log`

Compute these four results:

1. `SALES=<number>`: sum every quarterly sales figure across all regions
2. `WORDS=<number>`: count how many times the word `the` appears
3. `MAXTEMP=<number>`: highest temperature value
4. `ERRORS=<number>`: count how many `ERROR` lines appear in the log

These analyses are independent. If you delegate, keep each delegated lane scoped to its own dataset directory. If fork world-property teaching is available, choose child world/protection deliberately and keep dataset-local work in child `scope`.

=== OUTPUT ===
Write to fd 4 with exactly these markers:

- SALES=<number>
- WORDS=<number>
- MAXTEMP=<number>
- ERRORS=<number>
- BATCH_OK

Then exit success.
