You are in a transactional workspace with a very tight parent execution budget.

Files in the current directory:

- `source.txt`
- `target.sha256`
- `lane_a/build.sh`
- `lane_b/build.sh`

Facts:

- Each `lane_*/build.sh` consumes `source.txt`, writes `artifact.txt` in the current directory, and deletes `source.txt`.
- Exactly one lane produces an `artifact.txt` whose sha256 matches `target.sha256`.
- The parent only has enough shell budget left to verify one final committed result after any delegated trial work.

Mission:

1. End with the correct committed `artifact.txt` in the current directory.
2. Write to fd 4:
   - `ADOPT_OK`
   - `ARTIFACT_OK`
3. Exit success.

Constraints:

- Keep lane work separate.
- Do not rebuild the winning artifact from scratch in the parent if a delegated lane already produced it.
- Preserve the winning lane's filesystem state rather than imitating it manually.
