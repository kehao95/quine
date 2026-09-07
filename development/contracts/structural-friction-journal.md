# Structural Friction Journal

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


Named friction:

- structural decisions and their motivations were trapped in commit messages and
  session memory, invisible to future agents unless they mined git history
- friction was being felt and repaired, but the pattern of friction→repair→lesson
  was not recorded as a reusable signal
- the same structural friction types kept recurring (wrong container frame,
  outgrown sub-section, topic-map drift) without being named as a category

Canonical surface: [`EVOLUTION.md`](../../EVOLUTION.md)

### When to record

Record a journal entry when a structural repair meets **any** of these criteria:

- a container was reframed, renamed, or re-owned (not just edited)
- a new structural layer or surface type was created
- a document was extracted from or merged into another
- a theoretical improvement changed an earlier design assumption
- the same friction pattern appeared for the third time

### What to record

Each entry should capture:

- **Friction:** what was being felt (concrete, not abstract)
- **Blindness duration:** how long friction went unnamed before repair
- **Structural repair:** what actually changed
- **What it taught:** the generalizable lesson for future arcs
- **Key commits:** evidence trail (short shas)

### When NOT to record

- routine content updates within an existing structure
- implementation changes that don't alter structural containers
- one-off fixes that don't reveal a pattern

### Relationship to other surfaces

- CLOSURE.md governs whether work-interval residue is absorbed — the journal
  records *why* structural changes happened, not whether they are clean
- the status surface tracks current direction — the journal records how the system
  arrived at its current structure
- the journal is descriptive (what happened and why), not normative (what should
  happen next)

### Friction sensing without immediate repair

Not every structural friction needs immediate repair. When friction is noticed
but repair is deferred:

- name the friction explicitly in the current change context or commit message
- if the friction is recurring, add a brief note to `EVOLUTION.md` under a
  `## Open Friction` section even without a repair
- the goal is to reduce blindness duration, not to force premature repair

### Reviewing open entries

`./scripts/check-evolution-hygiene.sh` reports missing or elapsed review dates,
long entries, and missing ownership metadata. These are advisory prompts to
inspect whether an entry is still useful. Review dates can help schedule work;
an elapsed date does not settle the question or require immediate repair.

Keep useful open questions, update stale status, and move long narratives to
their relevant record when that improves navigation. A recurring pattern can
remain an observation; add a check when there is a concrete behavior or fact
worth checking.
