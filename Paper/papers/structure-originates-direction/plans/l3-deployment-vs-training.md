# L3 — the RL-frame facet: deployment ≠ training

**Status:** BLOCKED on L2 (do not run before L2's gate opens).
**Gate this closes:** deployment-time force ≠ training curriculum (RL/UED's
central formalism cannot express this case).

**Proves:** steers a fixed policy with no reward, no weight change, difficulty
held constant — a case RL's regret/curriculum machinery cannot represent.

**Design:** frozen agent, no reward channel, difficulty (success-rate/optimal-
gap) equated across structure conditions; DV shifts.

**Gate:** holds → L4.

**Forward-pointer:** will need a difficulty-equalization protocol (matching
success-rate/optimal-gap across structure conditions) that doesn't yet exist
in the harness — this is the main new apparatus cost of this layer.
