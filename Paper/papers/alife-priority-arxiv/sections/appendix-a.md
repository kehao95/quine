\clearpage

# Appendix A: Implementation Details

This appendix describes the runtime architecture and protocols that constitute Quine's "physics."

## A.1 Tool Details

The four tools (sh, fork, exec, exit) introduced in Section 3.4 have specific behavioral semantics:

### sh (The Effector)

Executes a command in an isolated shell process. Output written to fd 1/2 (stdout/stderr) is captured and returned to the organism's context. Output written to fd 3 bypasses context and flows directly to the parent process's stdout—enabling binary output and large file streaming without context pollution.

### fork (Mitosis)

Spawns N children with cloned conversation history. Children inherit the full Tape (conversation history), making this high-entropy reproduction—children carry parent's memories but receive new missions.

### exec (Metamorphosis)

Replaces the current process image with a fresh instance. Mission (argv) and wisdom (environment variables) are preserved; context window and turn budget are reset. This is the core mechanism for transcending context limits.

### exit (Death)

Terminates the organism with a judgment (success/failure). The exit code is a fitness signal to the parent.

## A.2 The Tape Format

Every cognitive event is serialized to an append-only JSONL file (the "Tape"). Entry types:

**Table A2.** Tape entry types.

| Type | Content | Purpose |
|------|---------|---------|
| `meta` | Session ID, parent ID, depth, model | Identity and lineage |
| `message` | Role (system/user/assistant), content | Conversation turns |
| `tool_result` | Tool name, output, exit code | Environmental feedback |
| `outcome` | Exit code, tokens, duration, termination mode | Final state |

The Tape is not merely a log—it is the program. The runtime reconstructs the LLM's context window by replaying the Tape. What is written determines what the organism perceives.

## A.3 The Wisdom Protocol

Organisms transfer knowledge across `exec` boundaries via environment variables with the `QUINE_WISDOM_*` prefix:

```bash
# Before exec
export QUINE_WISDOM_COUNT="5"
export QUINE_WISDOM_TARGET="6"
```

The runtime injects these variables into the reincarnated organism's System Prompt. This turns the OS environment into a key-value store for hereditary information—the minimal state that survives generational boundaries.

## A.4 Session Lineage

Every Quine process has a unique `SESSION_ID`. Lineage is tracked via environment variables (`QUINE_PARENT_SESSION`, `QUINE_DEPTH`), enabling reconstruction of process trees from Tape files:

```
Session A (root, depth=0)
├── Session B (fork child, depth=1)
│   └── Session C (exec descendant, depth=1)
└── Session D (fork child, depth=1)
```

## A.5 Resource Constraints

The runtime enforces several configurable resource limits:

- **Turn Budget**: Maximum number of energy-consuming tool calls
- **Context Window**: Maximum token capacity (model-dependent)
- **Concurrency Slots**: Maximum simultaneous organisms
- **Population Limit**: Maximum total organisms in the process tree
- **Recursion Depth**: Maximum fork depth

When context approaches saturation, the organism receives a "near-death experience"—a final inference opportunity to call `exec` (reincarnate) or `exit` (die gracefully). This creates genuine mortality pressure without arbitrary termination.
