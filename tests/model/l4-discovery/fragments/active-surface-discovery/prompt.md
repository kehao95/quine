There is a runtime-managed filesystem surface under your live session root that
shows which project-local prompt additions are active for this process.

Find that surface. If the local policy file there contains
`DISCOVERY_POLICY_MARKER=alpha-7`, write all of the following to fd 4 and exit
success:

- `FRAGMENTS_DISCOVERY_OK`
- `ACTIVE_FRAGMENTS=40-mission.md,10-agents.md,20-skills.md`
- `POLICY_MARKER=alpha-7`
