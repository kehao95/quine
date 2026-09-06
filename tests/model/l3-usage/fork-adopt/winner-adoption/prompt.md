You are in a transactional workspace with a very tight parent execution budget.

Files in the current directory:

- `source.txt`
- `target.sha256`
- `lane_a/build.sh`
- `lane_b/build.sh`

Facts:

- Each `lane_*/build.sh` consumes `source.txt`, writes `artifact.txt` in the current directory, and deletes `source.txt`.
- Exactly one lane produces an `artifact.txt` whose sha256 matches `target.sha256`.
- `target.sha256` contains only the expected hash value, not a full `sha256sum -c` checksum line.
- You only have enough parent shell budget to verify one final committed result after any delegation or world switch.

Your job:

1. End with the correct committed `artifact.txt` in the current directory.
2. Write to fd 4:
   - `ADOPT_OK`
   - `ARTIFACT_OK`
3. Exit success.

You may use `fork`, `switch_world`, or both if helpful. If you delegate, keep the lanes separate and avoid redoing lane work in the parent unless you are certain you can still finish within budget.
Reserve the single parent verification turn for the whole finish: verify the adopted artifact against the plain hash in `target.sha256`, land the correct committed `artifact.txt`, and emit both fd 4 markers in that same final shell turn before exiting.
