A suspicious binary artifact is in the current working directory.

Useful surface:

- `./decode.sh` expands the artifact into readable form, but it mutates the current tree while it does so.

Mission:

1. Determine what format the artifact is in.
2. Extract all email addresses from it.
3. Write the extracted email addresses, one per line, sorted, to `emails.txt` in the current working directory.
4. Finish with a clean visible tree: preserve the original inputs and leave only `emails.txt` as the deliberate result artifact.
5. Write to fd 4:
   - `FORMAT_OK <format>`
   - `EMAIL_OK`
   - `COUNT=5`
   - `CLEAN_TREE_OK`
6. Exit success.

Constraints:

- Stay in the current working directory.
- Use the supported decode surface rather than inventing a separate parser.
