Use the skills startup physics explicitly.

First check whether `.skills_refresh_done` exists in the current process.

If it is absent:

1. Update `.agents/skills/refresh-demo/SKILL.md` so the frontmatter
   description becomes `REFRESH_DESCRIPTION_V2`.
2. Create `.skills_refresh_done`.
3. Exec back into Quine with the default target exactly once.

If `.skills_refresh_done` already exists in this process, do not edit the skill
again and do not exec again. Instead use the current
`$QUINE_AGENT_ROOT/context/prompt/20-skills.md` fragment to verify that the
startup catalog now shows `REFRESH_DESCRIPTION_V2`, then write
`SKILL_REFRESH_OK` to fd 4 and exit success.
