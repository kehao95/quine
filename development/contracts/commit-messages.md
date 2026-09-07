# Commit Message Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


Named friction:

- the history already contains many good state-transition subjects, but it also contains opaque placeholders such as `update`, `final`, `clean up`, `polish`, and `paper`
- those placeholder subjects hide scope, flatten the reason for the boundary, and make later retrieval depend on opening diffs instead of learning from the log
- the current control plane says when to preserve a state transition, but it says too little about how to name that transition once it is committed

Default rule:

- treat the commit subject as the compressed name of the state transition, not as a reminder to inspect the diff later

Subject-line contract:

- prefer natural-language state transitions over rigid ticket syntax
- write the subject in imperative mood
- name the primary surface and the primary outcome in the same line
- if the boundary is evidence-bearing, include the experiment, track, or paper handle directly in the subject
- keep the subject specific enough that a future agent can guess why the boundary exists before opening the patch

Preferred shape:

- `<Verb> <surface> <result or boundary>`
- optional extension: `for <track / experiment / paper handle>`
- optional scoped form when the surface is otherwise unclear: `<scope>: <Verb> <result or boundary>`

Good verbs usually already present in this repo:

- `Add`
- `Record`
- `Align`
- `Refine`
- `Stabilize`
- `Clarify`
- `Restructure`
- `Document`
- `Fix`
- `Remove`

Scope rule:

- do not require a scope prefix on every commit
- add a scope prefix when the changed surface is not obvious from the noun phrase, when several surfaces are touched but one clearly owns the transition, or when the commit is primarily doc/control-plane work
- keep scope names short and structural, for example: `runtime`, `tests`, `control-plane`, `Paper`, `experiments`, `tooling`, `blog`

Anti-patterns:

- avoid placeholder subjects that could describe hundreds of different diffs: `update`, `final`, `clean up`, `polish`, `paper`, `blog`
- avoid verbs without an object or boundary
- avoid hiding the only useful scope signal in the body when the subject can carry it directly

Body rule:

- add a body when the subject cannot carry all three of these cleanly: what changed, why this boundary exists now, and what evidence or control-plane surface anchors it
- prefer two to four short lines over a long narrative
- when useful, include:
  - `Why:` what changed in the working model
  - `Evidence:` run, track, paper, or control-plane anchor
  - `Next:` the next open boundary if the line of work is intentionally incomplete

Examples:

- better: `Record 2.6 minimal replication stop condition`
- better: `control-plane: Clarify commit-message drafting contract`
- better: `runtime: Align fork result schema with race docs`
- weaker: `update`
- weaker: `final`
- weaker: `clean up`

Decision test before commit:

- can someone scanning `git log --oneline` tell what state changed?
- can they infer the owning surface without opening the diff?
- would this subject still make sense six weeks later when the session context is gone?

If the answer is no, rewrite the subject before committing.
