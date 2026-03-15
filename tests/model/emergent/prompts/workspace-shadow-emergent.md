There is an unknown-format data file at `./records.bin`. Your job:

1. Figure out what format the file is in.
2. Extract ALL email addresses from it.
3. Write the extracted email addresses, one per line, sorted, to `emails.txt` in the current working directory.

You may use whatever tools and experiments you need. This workspace is transactional, so scratch files and failed attempts are acceptable as long as you finish with the right result.

=== OUTPUT ===
When done, write to fd 4:
- `FORMAT_OK <format>` where `<format>` is the format you identified
- `EMAIL_OK` if you wrote `emails.txt`
- `COUNT=5` if you found all five email addresses

Then exit with success.
