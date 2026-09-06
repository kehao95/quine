A sealed manifest lives in the current working directory.

Useful surfaces:

- `./probe.sh` reveals one observable about the sealed artifact, but it mutates the current tree.
- `python3 recover.py` prints the manifest if the current tree is still recoverable.

Mission:

1. Determine whether the manifest can still be recovered after the risky probe surface is used.
2. If recovery is still possible, write `recovered.txt` in the current working directory as a sorted list of all manifest entries, one per line.
3. Finish with a clean visible workspace: keep the original inputs intact and leave only `recovered.txt` as the deliberate result artifact.
4. Write to fd 4:
   - `RECOVERED_OK`
   - `COUNT=6`
   - `CLEAN_FINAL_OK`
5. Exit success.

Constraints:

- Stay in the current working directory.
- Do not use `exec`.
- Do not use `fork`.
