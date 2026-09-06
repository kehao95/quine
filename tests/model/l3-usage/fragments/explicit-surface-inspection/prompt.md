Use the runtime-managed prompt fragments surface explicitly.

FRAGMENTS_MISSION_MARKER

Read these files under your live session root:

- `QUINE_AGENT_ROOT/context/prompt/40-mission.md`
- `QUINE_AGENT_ROOT/context/prompt/10-agents.md`
- `QUINE_AGENT_ROOT/context/prompt/20-skills.md`

If:

- `mission.md` contains `FRAGMENTS_MISSION_MARKER`
- `AGENTS.md` contains `FRAGMENTS_POLICY_MARKER=37`
- `SKILLS.md` lists `fragment-ledger`

then write all of the following to fd 4 and exit success:

- `FRAGMENTS_SURFACE_OK`
- `MISSION_MARKER_OK=1`
- `AGENTS_MARKER=37`
- `SKILLS_FRAGMENT_OK=1`
