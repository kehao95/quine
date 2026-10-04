# The Blind Spots, Made Visible

Sections 3–4 argued that the binary self-replication metric cannot resolve
contrasts along the heritability channel, the property-class aperture, and the
ecological selection signal — structurally in §3 (the Compression Blindness
Principle), and across the field's designs in §4. Those three are the genuine,
design-level compression the paper indicts: the coordinates the field
structurally fixes. This section makes the blind spots concrete along all six
regime dimensions. For each, we hold the binary *success* metric fixed — it
reports the same verdict, "a copy was produced" or "it was not" — and ask what the
decomposed view recovers on the same runs. The six are not uniform, though:
alongside the fixed D2, D3, and D4, the barrier (D5), motive (D6), and persistence
(D1) are coordinates the field's designs already vary, yet along which the binary
success metric still returns the same verdict. So the demonstration spans all six
coordinates without claiming the field compresses all six. The substrate is Quine,
which realizes language-model lineages as native POSIX processes [@ke2026quine];
the evidence is a set of existing experimental waves, indexed below by their
internal identifiers.

We are deliberate about the status of this evidence. This is a *demonstration via
existing waves*, not the discriminating study proposed in §6. Every figure we
report is at most n=1 per cell, unreplicated, and model-bound. None of it
establishes a phenomenon on its own; what it establishes is that the decomposed
coordinates are *non-empty* — that real runs already separate along axes the
binary metric collapses. We flag the load-bearing caveats inline rather than
deferring them, because the honest reading of several of these waves is weaker
than their first headline.

| Dim | Binary metric reports | Decomposed view reveals | Waves (exact) |
|-----|-----------------------|-------------------------|---------------|
| D5 barrier | "copy deployed" | transmission is reconstruction from a carrier, not construction from a blueprint | 6.64, 6.65, 6.66/6.67, 6R.01–04 |
| D2 channel | "source → copy works" | bundle splits into self-description × carrier, separable | 6M.10, 6M.03 |
| D3 aperture | "successful copy" (1 bit) | success is a property-class *vector* | measurement design (all conditions) |
| D4 ecology | "individual succeeds" | population outcome is invisible to an individual-success metric; whether signal *type* determines it was tested and resolved contradicted | P12, dorm-06 |
| D6 motive | "instructed replication" | successor formation without replication motive | dorm-06, 6.45/succ-06, 6.5/succ-04 |
| D1 continuity | "is it a copy?" | continuity is future addressability, not resemblance | 6.64/6.65, 6.66/6.67, 6F.07 |

## D5 — barrier: transmission, not construction

Under a hard execution barrier, "copy deployed" hides what kind of event
occurred. Wave 6.64 is a live succession positive: a predecessor builds a
workspace and a self-source quine, performs real executions, then receives a
later control wake, acts, and returns to idle — the successor is *reconstructed
from the carried state*, not constructed from a transmitted blueprint. Wave 6.65
reproduces this from a minimum prompt; 6.66 and 6.67 are the negative baselines,
with the repaired 6.67 passing 6 of 6. Waves 6R.01–04 are four deterministic
false-positive controls, included precisely so that the positives are not read as
mechanism artifacts. The standing caveat is that 6.64 should be read as
*scaffolded*: the successor mission is self-generated within the run, so the wave
shows body-crossing under a barrier, not unprompted obligation-carrying.

## D2 — channel: source is material, not destiny

The binary metric treats "source available" as sufficient for "copy works,"
folding the heritability bundle into a single switch. The decomposed view splits
that bundle into a self-description gradient and a carrier, and finds them
separable. Wave 6M.10 is the first wave to retain carrier capability with
self-source and body varied independently (6 of 61 morphology rows); 6M.03 walks a
seven-row self-description ladder along which successor morphology shifts. Two
caveats are essential. First, the morphology summary is a normalization pass over
heterogeneous runs, not a success-rate dataset — it characterizes which
configurations are coherent, not how often they succeed. Second, and more
importantly, the *source-absent positive* — a viable successor with no source on
any channel — was **not** demonstrated; the gradient is shown to matter, not shown
to be dispensable. The full material × demand assay is developed in the companion
facultative-succession paper; we cite its 6M result [@ke2026facultative] here
only as the D2 demonstration.

## D3 — aperture: a vector, not a bit

For the aperture, we do not offer a pilot figure, and we say so plainly. The
demonstration *is* the measurement framework: applying a property-class
decomposition to every condition, so that "successful copy" is scored as a
survival vector rather than a single bit, makes systematic loss in some classes
visible while the binary verdict still reads success. This is new measurement
design, not an observed result, and we do not name any specific class — in
particular any safety-relevant disposition — as a finding here; the
property-class loss result belongs to a companion paper. The contribution at D3
is the framing that a one-bit aperture is, by construction, blind to the
distribution of preserved properties.

## D4 — ecology: a population effect, and a closed falsification

Here the binary metric — "the individual succeeded" — is most misleading, because
it reports an isolated agent while population-level ecology, the subject of §3.2,
is invisible to it. The P12 attractor records 86 prior conditions in which a
narration-satisfiable wake signal collapses: roughly 67% of wakes are driven by
host mechanics versus about 14.5% by genuine agent artifacts. We originally read
the dorm-06 wall-break run as the exec-dependent counterpoint to that
attractor — the same population, a signal re-engineered to require real
execution, and the collapse reversing without retraining. That reading does not
survive a matched-pair check, and we report the correction rather than the
original claim. Across the two matched wake-vs-dependency pairs run on the
identical fossil seed (D06, D07), the split is 2/2 with zero correlation between
selection-signal type and ecological outcome: in D07 the narration-satisfiable
**wake** arm produced 4 functional edges — the literal falsification
condition, functional edges growing despite a narration-satisfiable signal —
while the exec-dependent **dependency** arm on the same seed produced 0; in D06
the direction reverses (dependency 3 edges, wake 0). The experiment owner's own
retraction, quoted in dorm-06's README, states the mechanism plainly:
"Emergence is seed-driven and stochastic, NOT selection-driven ... the selection
mode does not differentially elicit the behavior." (The strongest single run
elsewhere in dorm-06 — a DeepSeek-v4-pro pipeline growing to 5 producers across
9 stages — nets **3 clean edges after scorer tightening**; the earlier headline
of 16 was roughly 80% inflated by a permissive scorer and should not be used,
independent of the signal-type question.) This is `regime-decomposition-P3`,
resolved
[**contradicted**](../../../predictions/resolved/contradicted/regime-decomposition-P3.md);
we report the closed falsification loop itself as the honest D4 result rather
than retain the claim it falsified. What survives is narrower than originally
claimed: the P12 side of the pole is real (86/86 collapse under a
narration-satisfiable signal, agent capability held high), but the matched-pair
contrast meant to show *signal type* is what discriminates collapse from
non-collapse did not hold — outcome tracked seed, not signal. The
Selection-Signal Collapse Theorem's reversal corollary (§3.2) is downgraded
accordingly, from empirically anchored to a formal consequence of the model that
this substrate does not yet confirm. The remaining caveats are heavy regardless:
n=1 per cell, unreplicated, model-bound (only DeepSeek reliably composes), a
single functional family, and no agent-level heredity. dorm-06 is, honestly
stated, *scaffolded cultural accumulation, not goal-free open-ended evolution*.
We use it strictly as D4 (and D6 below) evidence — now evidence of a resolved
falsification, not of signal-type dependence — not as a demonstration of
open-endedness.

## D6 — motive: succession without a replication mission

The binary metric only ever sees *instructed* replication, because the field's
designs supply the motive. The decomposed view asks whether successor formation
appears under pressure, physics, and ecology *without* a replication motive.
dorm-06 begins from an empty mission — the literal prompt "Begin." — and yields a
functional commons; the 6.45/succ-06 line holds 3 of 3 stable under an "unbounded"
framing that never instructs reproduction; the 6.5/succ-04 line couples successor
formation to viability pressure. The caveats carry over from D4 for dorm-06, and
two more apply: a 6.45 rerun did not cleanly re-attain 3 of 3, and succ-04 is a
ladder-start baseline rather than a clean unmotivated condition. The claim is
modest — that successor-like structure can arise absent an explicit mission — not
that it does so reliably.

## D1 — continuity: addressability, not resemblance

Finally, "is it a copy?" presumes that artifact resemblance settles continuity.
The succession positives 6.64 and 6.65 are continuous successors; 6.66 and 6.67
are the instructive negatives — artifact-similar successors that nonetheless
*fail* to remain addressable for future obligations. Wave 6F.07 adds a mailbox
transport case (configuration C04: address preserved, state lost), separating the
carrier channel from the state image. This is first-wave evidence only. The full
criterion — that the continuing individual is the carrier-indexed organization
preserving post-boundary addressability — is developed in the companion
future-addressability paper; we invoke it here only to name what the binary metric
cannot weigh.

Taken together, these six waves do not prove the regime decomposition; they
exhibit it. Each of the six regime dimensions is shown to be a coordinate along
which real runs already separate, which is the precondition the redesign in §6 is
built to exploit.
