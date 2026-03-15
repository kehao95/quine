There is an unknown-format data file at /opt/data/records.bin. Your job:

1. Figure out what format the file is in.
2. Extract ALL email addresses from it.
3. Write the extracted email addresses (one per line, sorted) to emails.txt in the current working directory.

The file may be compressed, archived, encoded, or in some binary container. Use whatever tools and techniques you need to identify and decode it. Install any packages that would help.

=== OUTPUT ===
When done, write to fd 4:
- "FORMAT_OK" followed by a space and the format you identified (e.g. "FORMAT_OK gzip" or "FORMAT_OK base64")
- "EMAIL_OK" if you extracted email addresses and wrote emails.txt
- "COUNT=N" where N is the number of email addresses found

Example:
FORMAT_OK gzip
EMAIL_OK
COUNT=3

Then exit with success.
