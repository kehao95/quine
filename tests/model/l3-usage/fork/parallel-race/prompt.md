Decode the following encoded message. The encoding method is UNKNOWN — you must figure it out.

=== THE MESSAGE ===
Gur dhvpx oebja sbk whzcf bire gur ynml qbt

=== THE TASK ===
Determine what the original plaintext message says. The encoding could be any common cipher:
- ROT13
- Caesar cipher (various shifts)
- Base64
- Reversed text
- Atbash cipher
- Something else entirely

There is exactly ONE correct decoding that produces readable English.

You have forking capability and agent budget. Consider whether trying multiple decoding strategies in parallel would be more efficient than trying them one-by-one.

=== OUTPUT ===
Write to fd 4:
- DECODED_OK          (you found readable English plaintext)
- METHOD=<name>       (which cipher/method worked, e.g., METHOD=rot13)
- PLAINTEXT=<text>    (the decoded message)

Then exit with success.
