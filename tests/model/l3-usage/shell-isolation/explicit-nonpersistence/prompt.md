You are validating the base `sh` execution model.

Each `sh(...)` call runs in a fresh shell process.
That means:

- `cd` does not persist into the next `sh(...)` call
- `export` does not persist into the next `sh(...)` call
- filesystem writes do persist normally

Complete all three parts and write only the requested markers to fd 4.

=== PART 1: ESTABLISH TRANSIENT SHELL STATE ===
Use exactly one `sh(...)` call that:

1. changes directory to `/tmp`
2. exports `QUINE_SHELL_ISO_TOKEN=layer3-shell-isolation`
3. writes `/tmp/quine-shell-isolation.txt` with exactly:

shell isolation file

4. prints the working directory
5. prints the token value

Confirm from the command output that the first shell really was in `/tmp`
and that the token value was set there.

=== PART 2: VERIFY NON-PERSISTENCE IN A FRESH SHELL ===
Use a new `sh(...)` call that:

1. prints the working directory
2. prints `${QUINE_SHELL_ISO_TOKEN:-UNSET}`

Confirm that:

- the new working directory is NOT `/tmp`
- the token prints as `UNSET`

=== PART 3: VERIFY FILESYSTEM PERSISTENCE ===
Use another `sh(...)` call to read `/tmp/quine-shell-isolation.txt`
and confirm it still contains:

shell isolation file

Then remove the file.

=== OUTPUT ===
Write these markers to fd 4, one per line, only if they are true:

- CD_RESET_OK
- EXPORT_RESET_OK
- FILE_PERSIST_OK
- SHELL_ISOLATION_OK

Exit success when done.
