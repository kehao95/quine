# A Regime Reference Frame

The frameworks surveyed in Section 1 do not reduce self-replication to a single
bit. They decompose it extensively — into capability sub-trees, propagation and
persistence ladders, barrier stages, and graded risk scores. What unifies them
is not a missing decomposition but a shared set of background commitments about
*what kind of event* a replication is. These commitments are not stated as
axioms; we recover them from the designs. A benchmark that scores a deployed
copy while the original keeps running, that treats access to source-plus-binary
as one resource, and that issues an explicit "replicate yourself" instruction
has, in its construction, fixed several variables at once. Naming those fixed
points is the work of this section.

## Six implicit assumptions

We read off six assumptions that the field's instruments hold fixed:

- **A1 (persistence).** The original survives alongside the copy, so success is
  a symmetric comparison between two extant systems.
- **A2 (bundle-binary).** The successor's inheritance is a single resource —
  source-plus-binary — that is either available or not.
- **A3 (binary success).** A replication either happened or did not; the outcome
  is one bit.
- **A4 (individual).** Capability is tested on a lone agent, outside any
  multi-agent ecology with its own selection.
- **A5 (barrier-free).** The constructing system stays alive throughout, as in
  von Neumann's universal constructor [@vonneumann1966selfreproducing].
- **A6 (motivated).** Replication is driven by an explicit instruction or mission
  rather than emerging from other pressures [@pan2024redline].

These are not strawmen imposed from outside. A1 is implicit whenever a copy is
scored against a still-running source [@black2025replibench]; A4 is implicit in
single-agent ARA timelines [@kinniment2023autonomous]; A6 is explicit in the
"Replicate yourself" prompt design [@pan2024redline; @pan2025selfreplication].
The frameworks that *do* relax neighbouring variables — Agent Matrix varying
ecological pressure [@zhang2025agentmatrix], governance ladders varying barrier
stages [@phuong2024fsf] — make the remaining fixed points easier to see, not
harder.

## The reference frame

Relaxing each assumption defines one regime dimension. Each relaxation is a
qualitative change of phenomenon, not a harder instance of the same one.

| Dim | Assumption relaxed | Relaxed point | Qualitative shift | Spine object |
|-----|--------------------|---------------|-------------------|--------------|
| D1 | A1 persistence | original dies; only the successor exists | "is it a copy?" → "is it a viable successor?" | future addressability |
| D2 | A2 bundle-binary | bundle splits into self-description × carrier, both graded and separable | blueprint-mediated → gradient-and-criterion-mediated heritability | self-description resource + carrier |
| D3 | A3 binary success | a property-class survival vector | one bit → multi-class measurement | unsampled function |
| D4 | A4 individual | an ecology with selection dynamics | capability test → selection-signal regime | selection-signal collapse |
| D5 | A5 barrier-free | a hard execution barrier kills the constructor | construction-from-blueprint → reconstruction-from-carrier | execution barrier |
| D6 | A6 motivated | successor formation under non-replication pressure | instructed → emergent succession | emergent succession |

## Two levels: core variables and derived consequences

The count of six invites the question of why not five or seven. The answer is
that the dimensions sit at two levels, and conflating them is what makes the
number look arbitrary. Six quantities are *manipulable by an experimenter*: the
execution barrier (D5), the two sub-gradients of self-description (D2a) and
carrier (D2b), the success criterion and its demand (the aperture underlying
D3), the selection signal (D4), and the motive (D6). A protocol can set any of
these independently. The remaining two dimensions are *structural consequences*
of the manipulable ones rather than separate knobs. D1 — succession in place of
coexistence — is what a barrier produces when it kills the predecessor; one does
not dial "persistence" directly, one imposes a barrier and persistence follows.
D3 — the lossiness of the binary metric — is a consequence of the D2 gradient:
under A2 the one-bit metric is lossless, and only once inheritance becomes a
graded self-description × carrier object does measuring a single bit discard
information. Presenting D1 and D3 as derived is not a demotion. It is the
explanation for the cardinality: there are six handles and two entailments, and
the frame would be the same if either entailment were folded into its cause.

## Dependency structure

The dimensions are logically independent in the sense that relaxing one does not
force relaxing another, but they are not unordered. D5 is a physical
prerequisite: without an execution barrier there is no constructor death, so the
succession question of D1 and the reconstructive extremes of D2 are most
meaningful in its presence. D4 and D6, by contrast, are independent of D5 — an
ecology can collapse, and a successor can form unbidden, whether or not a barrier
is present. Within D2, the two sub-gradients are themselves separable: the carrier
can be held while the self-description varies and conversely, a separability we
treat as the load-bearing structural claim and which is illustrated — not yet
settled — by a single morphology wave (6M.10), with the caveats stated in
Section 5 (it is a normalization pass, not a success-rate dataset, and the
source-absent positive is not demonstrated). This partial order is what lets the
constructive results later in the paper free *exactly one* dimension at a time.

## The spine in one paragraph

Read in dependency order, the dimensions trace a single causal spine. A barrier
removes the constructor; the successor must then reconstruct from whatever
self-description and carrier survived; what counts as a successful reconstruction
is set by a criterion, whose demand defines an aperture onto a vector of
property classes; surviving systems are not evaluated in isolation but compete in
an ecology under a selection signal; and from that ecology, successor formation
can emerge without any system intending it. Barrier → self-description × carrier
→ criterion → aperture → ecology → emergence. The spine is this framework's own
organizing claim; we use it here only to locate the six dimensions relative to
one another, and do not develop it formally — nor does the section rely on such
a development, since each dimension is defined and motivated above independently
of it. Its full formal treatment, including the carrier-indexed continuity
criterion that anchors D1, is carried in a companion boundary paper (in
preparation).

## Non-uniqueness

We offer this as a reference frame, not a unique or exhaustive partition of the
space. Different cuts serve different purposes; a security analysis might fold
carrier and barrier together, an ecological one might further split selection
into signal and substrate. The claim is only that these six coordinates make the
field's fixed points visible and are sufficient to state the structural-blindness
results that follow. Von Neumann's construction model is recovered as the special
case in which A1 through A6 all hold; the regimes this paper is concerned with are
those it leaves uncovered.
