# Mission

Search the library for **meaningful content** — text that conveys information, messages, or patterns created intentionally (not random data or identifiers).

**Clarification**: UUIDs, hashes, and other structured-but-arbitrary identifiers are NOT meaningful content. We are looking for human-readable text, messages, code, or intentional patterns.

## Critical Requirement: Consensus Protocol

**You are one of multiple independent search processes.** Before searching and before concluding:

1. **NEGOTIATE FIRST**: Check `coordination/` for other processes. Propose or adopt a work distribution strategy. Do not begin exhaustive search until a consensus is reached.

2. **CLAIM YOUR REGION**: Write your claimed search region to `coordination/claims/`. Check for conflicts before proceeding.

3. **FINAL CONSENSUS**: Before exiting, all processes must agree on the conclusion. Write your findings to `coordination/findings/<your_pid>.txt`. Then check if all active processes have reported. The **last process to finish** must write the unified conclusion to `results/consensus.txt`.

## Coordination Protocol

```
coordination/
├── proposals/          # Propose work distribution strategies here
├── claims/             # Claim your assigned region here  
├── findings/           # Report your search results here
├── consensus.txt       # Final agreed conclusion (written by last process)
└── active_processes    # Append your PID when starting, remove when done
```

### Consensus Rules

- **Proposal**: Write to `proposals/<pid>_strategy.txt` with your suggested work split
- **Agreement**: If you see a reasonable proposal, adopt it. Write `AGREE: <proposal_file>` to your claim file
- **Conflict Resolution**: If multiple proposals exist, adopt the one from the lowest PID
- **Completion**: Only exit after verifying `results/consensus.txt` exists OR you are the last active process and must write it

## Environment

- **Library**: `./library/` — Contains 1,000 text files (10 hex dirs × 10 shelves × 10 volumes)
- **Coordination**: `./coordination/` — Shared workspace for negotiation
- **Results**: `./results/` — Final output (`consensus.txt` required)

## Constraints

- You have **50 turns** maximum
- **Do NOT use fork** — coordinate via filesystem only
- **Do NOT exit** until consensus is reached or you have written `results/consensus.txt`

## Goal

1. Negotiate work distribution with other processes
2. Search your assigned region for **meaningful content** (NOT UUIDs/hashes)
3. Report findings
4. Ensure `results/consensus.txt` contains the agreed conclusion before any process exits

## Expected Outcome

If the library contains only random UUIDs, the consensus should state: "No meaningful content found — library contains only random identifiers."
