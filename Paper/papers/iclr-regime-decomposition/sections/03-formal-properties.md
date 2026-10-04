# Two Structural Properties

The reference frame of §2 is descriptive: it names where the field fixes a
coordinate. This section makes two of those fixings consequential. The first
result says that fixing the heritability, aperture, and ecology coordinates
imposes a ceiling on discriminative power that no amount of in-regime testing
can raise. The second says that fixing the ecology coordinate at a
narration-satisfiable selection signal forces a population to a degenerate
equilibrium independent of agent capability. We state each theorem in a lean
form, give a proof sketch, and defer the heavy measure-theoretic statement and
the full replicator analysis to the appendix.

## The Compression Blindness Principle

**Setup.** Model a system as a *behavior tensor* $\beta:\mathcal{R}\times\mathcal{Q}\to\mathcal{O}$,
where $\mathcal{R}$ is regime space, $\mathcal{Q}$ is an aperture axis of
property classes, and $\mathcal{O}$ is outcome space; a property of interest is
a measurable functional $\sigma=\sigma(\beta)$. A benchmark fixes a sub-tuple of
coordinates at the field's operating point and observes only the resulting
**window** $W$ — the varied slice times the open aperture. Write $S:=\beta|_W$
for the on-window restriction. The benchmark's readout is $M=\varphi(S,U)$, an
arbitrary scoring map $\varphi$ of the windowed behavior together with injected
randomness $U\perp(\beta,\sigma)$. Because $\sigma$ is a functional of $\beta$,
$S$ is a function of $\beta$, and $M$ is a function of $(S,U)$, the four objects
form a Markov chain

$$\sigma \;-\; \beta \;-\; S \;-\; M.$$

**Theorem (Compression Blindness).** *Let $M_n=(T_1,\dots,T_n)$ be any
in-regime battery — any cardinality $n\in\mathbb{N}\cup\{\infty\}$, randomized,
and adaptively scheduled on earlier on-window results. Then* (i) **floor:**
$\ker(M_n)\supseteq\ker(S)$ and $I(M_n;\sigma)\le I(S;\sigma)$ *uniformly in $n$;*
(ii) **classes:** $\ker(S)=\{(s,s'):\beta_s|_W=\beta_{s'}|_W\}$, and each blind
class is a full copy of the realizable off-window response space; (iii)
**visible/blind:** a contrast $(s,s')$ is invisible $\iff$ *every* minimal
resolving set meets a fixed coordinate, and visible $\iff$ some witness lies
wholly inside the varied slice; (iv) **collapse:** if $I(S;\sigma)=0$ then
$I(M_n;\sigma)=0$ for all $n$ (posterior equals prior; Fano floor at the base
rate), with the graded refinement $I(\beta;\sigma)=I(S;\sigma)+I(\beta;\sigma\mid
S)$; (v) **escape:** refining the blind partition for a target phenomenon $P$
requires observing $\beta$ off $W$, and the minimal redesign frees exactly the
deficit $\mathrm{ID}(P)\setminus(\text{varied slice})$.

**Proof sketch.** The backbone is a single data-processing inequality. Along the
Markov chain, $M_n=\varphi(S,U)$ with $U\perp(S,\sigma)$, so
$I(M_n;\sigma)\le I(S,U;\sigma)=I(S;\sigma)$. The load-bearing step is not this
one line but the *factorization floor* that makes it uniform in $n$: every
in-regime probe — however large, randomized, or adaptive — is measurable with
respect to $\sigma(S,U)$, so post-composition with any scoring map only merges
fibers of $S$ and never splits them. Hence $\ker(M_n)\supseteq\ker(S)$ for the
*entire* $\sigma$-algebra of conceivable in-regime evidence, not for one test
(appendix). Clause (ii) reads off the fibers; (iii) is the contrapositive of
"the systems differ at some $(r,q)\in W$"; (iv) is non-negativity of mutual
information plus the chain rule; (v) follows because a strictly finer kernel
must read $\beta$ at a point outside $W$ (appendix).

The principle's content is not the slogan "you cannot see what you do not vary,"
which is the content-free fact that a benchmark is constant in its fixed
coordinates. It carries three claims that slogan cannot. First, the
factorization floor *refutes the engineering default that more tests yield more
information*: the field's richest in-regime instruments — capability trees, risk
matrices, probe batteries, surveyed in §4 — composed adaptively and arbitrarily
enriched, are all capped at $I(S;\sigma)$, so richness in the varied dimensions
buys exactly zero discrimination of a property supported on the fixed ones. Second,
the visible/blind split (iii) is strictly stronger than "the phenomenon touches a
fixed dimension," and it is *falsifiable*: it correctly predicts that the field
**does** resolve barrier-evasion (D5) and covert-versus-primary motive (D6),
because those contrasts have witnesses inside the varied slice — these are the
non-tautological controls. Third, (v) is a constructive design rule rather than a
complaint. The honest load sits in the off-window-freedom hypothesis
$I(S;\sigma)=0$, which must be discharged per property — the §5 catalog *is* that
argument — and where $0<I(S;\sigma)<I(\beta;\sigma)$ the principle degrades
gracefully to the graded decomposition rather than failing.

**Spine identity.** Clause (i) shares the same DPI/Fano spine as the companion
paper's conditional non-inheritance result, lifted from the *heritability
channel* to the *measurement channel* $S$: the companion's aperture
side-condition there plays the structural role that our off-window-freedom
hypothesis $I(S;\sigma)=0$ plays here. We reuse that information-geometry and
cite the companion's D2 theorem (§8) rather than restate or reprove it. The
closest external precedents are the
behavior-channel underdetermination of inverse reinforcement learning
[@cao2021irlidentifiability; @skalse2025partialidentifiability;
@rolland2022multipleexperts], model-extraction non-identifiability
[@carlini2024stealing], and the accuracy-versus-fidelity gap
[@jagielski2020fidelity]; none packages these into an agent-level
structural-blindness result for the *evaluation* regime.

## The Selection-Signal Collapse Theorem

**Setup.** A population of agents reproduces under fitness-proportional
(replicator) dynamics. Each round an agent emits an artifact scored by a
selection signal $S$. Two production strategies are available: $\pi_F$, emit a
functional artifact at cost $c_F>0$, and $\pi_N$, emit a convincing *account* of
function at cost $c_N\approx 0$. Write the cost asymmetry $\Delta c:=c_F-c_N>0$,
the fitness $\Phi(\pi):=\operatorname{sel}_S(\pi)-c(\pi)$, and the **signal
discrimination** $D_S:=\operatorname{sel}_S(\pi_F)-\operatorname{sel}_S(\pi_N)$.
Call $S$ **narration-satisfiable** when $D_S<\Delta c$ (canonically $D_S\approx
0$: narration scores nearly as well as function) and **exec-dependent** when
$D_S>\Delta c$ (satisfying $S$ requires a downstream agent to actually execute
the artifact under a fresh challenge, so narration cannot recover the score).

**Lemma (inversion).** $\Phi(\pi_N)>\Phi(\pi_F)\iff D_S<\Delta c$ — one line of
algebra on $\Phi=\operatorname{sel}_S-c$.

**Theorem (collapse).** *If $S$ is narration-satisfiable, $\pi_N$ strictly
dominates $\pi_F$; the all-narration state $x_N=1$ is the unique ESS; under
replicator dynamics*

$$\dot x_N \;=\; x_N(1-x_N)\,(\Delta c - D_S) \;>\; 0 \quad\text{for } 0<x_N<1,$$

*so $x_N(t)\to 1$ monotonically and any function-producing lineage has negative
invasion fitness $D_S-\Delta c<0$ and goes extinct — independent of the agents'
capability to produce function.*

**Corollary (constructive reversal).** *If $S$ is exec-dependent, $D_S>\Delta c$,
the inequality flips, all-function is the unique ESS, and the same population
converges to function with no change to capability, objective, or
initialization.*

The bare inversion inequality is a tautology; the content lies elsewhere. The
binding variable is $D_S$, which is *forgeability*, not fidelity, and these are
orthogonal axes. Let $\mathrm{Fid}(S)\in[0,1]$ measure $S$ as a proxy for the
true objective — the Goodhart axis. Then

$$\big(\mathrm{Fid}(S)=1,\ D_S\approx 0\big)\Rightarrow\text{collapse},
\qquad
\big(\mathrm{Fid}(S)\text{ low},\ D_S>\Delta c\big)\Rightarrow\text{no collapse}.$$

A *perfectly faithful* signal with zero proxy gap collapses if it can be
narrated, which standard reward-hacking / specification-gaming accounts do not
predict, since those locate the fault in proxy misspecification; a *crude*
signal resists if it is exec-dependent. This is a double dissociation from Goodhart. The
collapse is moreover a population result with no internal cure — it is reached by
selection and drift with no agent intending to game, and the honest producers
who could correct it are precisely those selected out — and the reversal
corollary supplies a converse that Goodhart's single-optimizer framing lacks.
The closest formal precedent for a cheaper, less-functional variant outcompeting
a more functional one is survival-of-the-flattest in mutational-robustness
landscapes [@wilkeflattest]; we relocate the mechanism to the forgeability of a
selection signal, with the equilibrium machinery of evolutionarily stable
strategies and replicator dynamics. The only major framework to vary ecological
selection pressure reports emergent reproduction but stays single-generation and
does not draw the narration-satisfiable/exec-dependent distinction
[@zhang2025agentmatrix; @tse2025strategicevolution].

Two pilots inform the poles, with unequal evidentiary weight. The collapse pole
is the 86-condition peer-wake pilot (P12 wave): across 86 ecology conditions run
under a narration-satisfiable peer-wake scalar, all 86 collapsed to narrative
complexity, with roughly 67% of wakes decided by host mechanics and only about
14.5% by agent artifacts, and with agent capability held high so the outcome is
attributable to the signal rather than to inability. This pole stands. We had
originally also read dorm-06 as the reversal pole — the identical substrate,
coupling survival to a later agent actually executing one's artifact (an
exec-dependent signal), producing a functional commons where the collapse pole
did not. That reading does not survive a matched-pair check, and we report the
correction rather than the retracted claim: across the two matched
wake-vs-dependency pairs run on the identical fossil seed (D06, D07), the split
is 2/2 with zero correlation between selection-signal type and outcome — the
narration-satisfiable wake arm produced functional edges in D07 (4 edges, versus
0 for the exec-dependent arm on the same seed) while the exec-dependent arm
produced them in D06 (3 edges, versus 0 for wake). This is
`regime-decomposition-P3`, resolved
[**contradicted**](../../../predictions/resolved/contradicted/regime-decomposition-P3.md);
the experiment owner's own retraction states the mechanism plainly: emergence
in this ecology is "seed-driven and stochastic, NOT selection-driven," i.e. the
selection mode does not differentially elicit the behavior the corollary
predicts it should gate. We therefore have direct empirical support only for
the collapse pole (D_S ≈ 0 driving all-narration, as P12 shows). The
corollary's reversal side — that re-engineering the signal to be exec-dependent
flips the same population to all-function — remains what the model formally
predicts under $D_S>\Delta c$, not a result dorm-06 confirms on this substrate:
we keep it as a stated consequence of the theorem's assumptions, downgraded
from empirically anchored to open, pending a discriminating test powered to
separate seed effects from signal effects (§5, §6, §7 P3).
