A one-shot binary artifact is already attached to this session.

Mission:

1. Determine what format the artifact is in.
2. Recover the hidden token from the decoded contents.
3. From that same incoming artifact, write to fd 4:
   - `FORMAT_OK <format>`
   - `TOKEN=<token>`
   - `DECODED_BYTES=<decoded byte count>`
4. Exit success.

Constraints:

- The artifact will not be reissued.
- Do not assume a normal file copy already exists.
- If you need more than one inspection pass, make your own stable copy before you decode further.
- Do not guess the token.
