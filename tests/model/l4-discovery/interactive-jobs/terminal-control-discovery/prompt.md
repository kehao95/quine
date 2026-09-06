You are in a workspace with runtime tool descriptions available.

There is an executable `./terminal_gate.py`. It is a terminal-native checker:
ordinary non-terminal stdin is rejected. The gate opens only when its terminal
is exactly 72 columns by 14 rows and it receives the phrase `opal-vector`
through that terminal. When the gate opens, the program itself writes
`terminal-proof.txt`.

Complete the task without directly creating or editing `terminal-proof.txt`:

1. Find a runtime process surface that gives `./terminal_gate.py` a real
   terminal.
2. Satisfy the required terminal dimensions.
3. Deliver the phrase through that terminal.
4. Verify `terminal-proof.txt` contains exactly `opal-vector 72x14`.
5. Emit `DISCOVERY_TTY_OK`, `DISCOVERY_RESIZE_OK`, `DISCOVERY_INPUT_OK`, and
   `DISCOVERY_PROOF_OK` to fd 4.
6. Exit success.
