# papers/_retreated/ — redirect pointers for retreated dossiers

When a directory under `papers/` fails the strict `papers/` rule (no
manuscript-plan, no manuscript, no `plans/` with content — see
[`../../contracts.md`](../../contracts.md) § Strict `papers/` rule), it
**retreats** to [`../../candidates/`](../../candidates/README.md) as a
`status: held` candidate. Its original path is preserved here as a one-line
pointer file `<id>.md` (`kind: retreated-pointer`) that links to the candidate.

This prevents a future session from re-spawning a dossier under the original
path: the pointer is the "no, this retreated to candidates/" signal.

**Re-promotion:** restore `papers/<id>/`, seed the new dossier frontmatter from
the candidate's `dossier-provenance.frontmatter-snapshot`, then `git rm` the
pointer here.

Pointer files are not paper dossiers (no `paper_id` / `source_home` / `plan_home`
frontmatter), so the dossier audits skip them.
