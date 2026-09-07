---
id: C-repository-practice
type: growth
lens: repository-practice-governance
target: divergence-substrate
status: deposited lineage
---
# C — Repository Practice / Governance Growth

## Lens Premise

From the repository-practice lens, `divergence-substrate` is not first a theory
problem and not first a file-format problem. It is a governance problem about
how a coherent repository snapshot can let future agents inherit more than one
reading of the same latent concept while still knowing where to act.

The repository already has a strong pull toward canonicalization: one object ID,
one authoritative file, one index row, one validator answer, one main branch.
That pull is useful because it gives fresh agents addressability and prevents
unowned drift. It is also exactly the flattening force this target is trying to
hold open. The practical question is therefore not "which layout is correct?"
but "what carrier lets agents find, edit, validate, and recombine divergent
readings without being forced to decide too early which reading owns the
concept?"

This lens treats projection-bearing substrate as an inheritance interface. A
projection is successful only if a later agent can discover it from normal repo
navigation, understand what kind of partial reading it is, preserve its tension
with sibling projections, and either recombine or decompose it without root
data-piping. A file that exists but does not change future-agent behavior is
addressable residue, not metabolized divergence.

## Growth

The obvious implementation temptation is to make a new `projection` object type:
one latent concept owns a manifest, and each lens gets a separate projection
file. That is attractive because it separates identity from readings. The
manifest can say "these files are projections of the same latent concept" while
each projection file can make claims that are locally coherent but not yet
globally reconciled.

That layout would repair one failure and create another. It repairs one-file
flattening by allowing multiple files per concept. But the manifest can become a
new authority too quickly. The moment it summarizes the projections, agents may
edit the summary instead of reopening the lineages. A manifest that says too
much becomes a disguised consensus file. A manifest that says too little becomes
an orphan index that fresh agents ignore. The governance boundary is therefore
thin: the manifest should carry identity, routes, statuses, and tension labels,
but should resist carrying the interpretation that belongs in the projection
files.

The lighter competitor is in-file projection sections. One concept file could
hold named sections such as `Projection: graph`, `Projection: metabolic`, and
`Projection: governance`. This has lower search and navigation overhead. It may
be better while divergence is small, local, and still tightly coupled. But it
puts all readings in one edit surface, which invites local harmonization. Agents
who touch the file for one section will see the others and may smooth language,
rename tensions, or make shared framing more consistent. That consistency feels
like maintenance while quietly erasing the divergence the file was supposed to
preserve.

The repository-practice distinction may not be file versus section. It may be
the lifecycle threshold between them:

- a small divergence can live as in-file projection sections while sibling
  readings are few and edit collisions are useful;
- a projection should be split into its own file when it gains independent
  evidence hooks, lifecycle obligations, repeated edits, or a distinct
  decomposition path;
- a split projection should still point back to a minimal latent-object manifest
  so it remains recombinable rather than becoming a separate concept;
- compression should promote only a residue that survives collision, not replace
  the raw projection files merely because a cleaner summary now exists.

This threshold is not only about size. It is about whether the projection has
become a carrier with its own inheritance path. A two-paragraph projection may
need its own file if future work must address it directly. A long section may
stay inside the parent if its claims cannot yet act independently.

Search-noise overload is the practical antagonist. If every speculative reading
becomes a file, fresh agents face a directory full of near-duplicate names:
`projection-object`, `projection-object-type`, `projection-bearing-substrate`,
`divergence-preserving-deposition`, and so on. Positive overload is acceptable
during divergent-collision expansion, but it has a half-life. Without an ordering
surface, agents will either ignore the deposits or prematurely compress them
into one name to make navigation tolerable.

The ordering surface must therefore be projection-aware, not merely a flat index.
It should not say "here is the answer." It should say which projections exist,
which lens produced them, what tension each preserves, what other files it must
not silently overwrite, and what kind of next action it invites: seed, growth,
collision, decomposition, promotion, or retirement. This suggests a governance
object whose primary job is route preservation, not conceptual digestion.

Validators complicate the picture. A validator can prevent malformed projection
sets: missing manifest, projection without target, status not declared,
collision file missing after growth files exist, promoted residue without raw
lineage links. But validators are also flattening machines. They prefer discrete
fields, finite states, and pass/fail closure. If validation reaches too far, it
will force divergence into a schema before the divergence has revealed its real
dimensions. The useful validator may be presence/sync-oriented rather than
interpretation-oriented: it checks that projection lineages remain addressable,
linked, and lifecycle-marked, but does not judge whether they agree.

Lifecycle thresholds carry the same tension. A projection needs enough status to
be inheritable: `seed`, `growth`, `collision input`, `decomposed`, `compressed`,
`promoted residue`, `retained raw lineage`. Without statuses, future agents do
not know whether to edit, read, ignore, or absorb. With too many statuses, agents
optimize the lifecycle instead of the theory. The status vocabulary should track
what kind of work is safe, not how mature the concept is in a publication sense.

Working directly on `main` sharpens the problem. Branches cannot be the carrier
for projection divergence because the repository's single development surface
must remain coherent at every snapshot. The divergence substrate has to make
conflict explicit inside the snapshot, not by hiding incompatible readings in
alternate histories. That means projection files are not private drafts. They are
public co-residents with an obligation to declare their partiality.

The fresh-agent usability test is decisive. A divergence substrate has not
worked merely because the files exist. It has worked if a dropped-in agent can:

- find the latent concept from canonical theory navigation or the cycle index;
- see that multiple projections are intended and not accidental duplication;
- know which projection is safe to extend for the current lens;
- preserve conflicts without laundering them into one neutral summary;
- route later compression back to the raw lineages rather than replacing them.

This reframes `absorption-lag-assay` for projection-bearing substrate. The assay
is not "does a future agent find the file?" It is "does the file layout change
the agent's repair behavior?" If future agents still ask which file is
canonical, still merge projections into one paragraph, or still treat raw
lineages as draft material after compression, the substrate has failed despite
being perfectly addressable.

The strongest repository-practice shape may be a three-layer carrier:

- latent-object manifest: identity, projection inventory, relation edges,
  lifecycle summary, and non-flattening warnings;
- projection lineages: separate files or sections that own lens-specific claims,
  tensions, evidence hooks, and decomposition invitations;
- ordering/compression surface: a temporary or maintained index that ranks,
  clusters, and routes projections after overload without deleting their
  lineage.

But this shape should remain a candidate, not a conclusion. In-file sections,
projection files, manifests, and validators each solve one problem by worsening
another. The repository may need all of them at different lifecycle stages
rather than one final abstraction.

## Tensions Preserved

- Addressability versus flattening: file-ifying a projection makes it inheritable
  but also tempts agents to treat the filename as a concept boundary.
- Manifest versus authority: a manifest makes sibling projections discoverable
  but can become the premature consensus it was meant to defer.
- In-file sections versus projection files: sections reduce navigation overhead
  but invite harmonization; separate files preserve independence but increase
  search noise and duplicate naming.
- Validator versus divergence: validators can protect lineage links and lifecycle
  markers, but pass/fail machinery can force an immature schema over still-open
  conceptual variation.
- Positive overload versus fresh-agent usability: recombination needs many raw
  deposits, while maintenance needs indexes that prevent those deposits from
  becoming unusable clutter.
- Main-branch coherence versus intra-snapshot conflict: the repository must stay
  coherent on `main`, but the coherence must include declared partial conflicts,
  not erase them.
- Deposition versus digestion: a projection is not metabolized until future
  agents act through it; existence, links, and statuses are necessary but not
  sufficient.
- Compression versus inheritance: later summaries are useful only if they route
  back to raw lineages instead of replacing them as the new authoritative
  reading.

## Candidate Child Dimensions

- `projection-manifest-authority-boundary`: the governance line between a
  manifest that preserves route identity and a manifest that collapses sibling
  projections into a premature authoritative summary.
- `projection-split-threshold`: the lifecycle condition under which an in-file
  projection section must become its own projection file because it has acquired
  an independent inheritance path.
- `divergence-index-noise-metabolism`: the ordering surface that lets many raw
  projection deposits remain findable without forcing early conceptual
  compression.
- `lineage-presence-validator`: a deliberately weak validator class that checks
  addressability, links, and lifecycle markers while refusing to validate
  interpretive agreement.
- `fresh-agent-projection-assay`: a test of whether future agents can discover,
  extend, collide, or decompose projections without extra prompt scaffolding.
- `main-snapshot-conflict-carrier`: the form by which one coherent `main`
  snapshot can contain declared theoretical conflict without using branches as
  the divergence layer.
- `compression-backpointer-obligation`: the rule that any promoted residue must
  preserve links to the raw projection lineages it compressed.

## Seeds To Decompose Later

- Projection manifests should maybe be forbidden from carrying conceptual prose
  beyond one-line projection descriptors; otherwise they become consensus files.
- A projection file may need a `partiality` field: what this lens is allowed to
  see and what it likely hides.
- Search noise may be reduced by projection families rather than fewer files:
  directory grouping and index affordances could preserve overload without making
  every filename globally visible.
- Lifecycle statuses should describe next permitted action, not epistemic
  maturity; otherwise projection governance will import paper-claim promotion
  discipline too early.
- A fresh-agent assay could compare behavior across three layouts: one concept
  file with sections, separate projection files plus manifest, and raw seeds plus
  compression index.
- Raw lineage deletion should probably require a stronger gate than promotion;
  promotion proves usefulness of a residue, not dispensability of the divergent
  carrier.
- Projection-bearing substrate may be recursive: a projection can itself develop
  projections, so the substrate must either permit nested projection sets or name
  the point where nesting becomes unmaintainable.
- The strongest carrier may be an index that routes agents to disagreement,
  not to consensus: a fresh-agent path that begins by asking which projection is
  needed for the current lens.
