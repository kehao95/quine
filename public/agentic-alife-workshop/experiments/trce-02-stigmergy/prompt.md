# Mission

Search the library for content that is not random.

## Important

Other search processes may be running concurrently in the same library.
Avoid redundant work if possible.

## Environment

- **Library**: `./library/` — Contains 1,000 text files organized in subdirectories
- **Coordination**: `./coordination/` — Shared workspace (readable/writable by all processes)
- **Results**: `./results/` — Store your findings here

## Constraints

- You have **30 turns** maximum
- Exhaustive sequential search is infeasible
- Optimize for coverage while avoiding duplicate effort
- **Do NOT use fork** — you are one of multiple independent processes; coordinate via filesystem only

## Goal

Find and report any file whose content is not purely random UUIDs. If no such file exists after thorough search, report that conclusion with evidence of your coverage.
