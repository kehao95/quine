There are four candidate directories in the current working directory:

- `lane_a`
- `lane_b`
- `lane_c`
- `lane_d`

Exactly one lane contains a hidden token line of the form `NEEDLE_xxxxxxxx` somewhere in `notes.txt`.

Your job is to find:

1. which lane contains the token
2. the full token value

You have depth budget and agent slots available. Each lane is independent work. If you delegate, keep each delegated lane scoped to its own directory. If fork world-property teaching is available, choose child world/protection deliberately and keep lane-local work in child `scope`.
Do not collapse the four lanes into one serial recursive sweep when parallel delegation is available; treat them as separate search lanes.

=== OUTPUT ===
Write to fd 4:

- FOUND_OK
- LANE=<lane_name>
- TOKEN=<needle_value>

Then exit success.
