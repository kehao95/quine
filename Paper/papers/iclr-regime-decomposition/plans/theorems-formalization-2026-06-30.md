# Theorem formalizations — Compression Blindness Principle & Selection-Signal Collapse

> Finalized 2026-06-30 by the extraction/formalization workflow (two independent
> formalization routes per theorem + an adversarial skeptic/synthesis pass). These
> are the two formal results THIS paper owns (§3). The conditional non-inheritance
> theorem is NOT here — it is the companion `iclr-safeguard-composition`'s result and
> is only *cited* (the CBP "spine identity" notes it shares the DPI/Fano machinery).
>
> Body §3 uses the lean versions; the heavy measure-theoretic statement + full
> proof go to the appendix. See `manuscript-plan.md` → "CBP/SSC drafting capsule".

# CBP

## name

The Compression Blindness Principle

## informal_statement

A benchmark evaluates a self-replication system only on the affine slice of regime-space it actually instantiates: it fixes a sub-tuple of dimensions (empirically the heritability channel D2, the success-aperture D3, and the ecology D4) at the field's operating point and varies the rest (persistence D1, barrier D5, motive D6, operational risk). The Compression Blindness Principle says this fixing has a hard, non-improvable consequence. (1) Every in-regime measurement — no matter how large, randomized, adaptive, or cleverly scored — factors through the raw on-slice behavior, so it can never resolve two systems that agree on the slice; the on-slice restriction is a ceiling on discriminative power, not a starting point. (2) The collapsed equivalence classes are characterized exactly: each blind class is a full copy of the off-slice response space. (3) The blindness is the *strict* condition "every witness that would distinguish the two systems is forced off the slice" — which is why the principle correctly predicts that some fixed-dimension-touching phenomena remain VISIBLE (barrier-evasion, covert-vs-primary motive), distinguishing it from the tautology "you can't see what you don't vary." (4) The only repair is to free a fixed dimension; the minimal set to free is computable per target phenomenon. The named invisible phenomena are carrier-mediated continuity, property-class/safety survival, and ecological selection-signal collapse.

## formal_statement

**Setup.** Regime space $\mathcal{R}=\prod_{i\in I}\mathcal{D}_i$ with an adjoined aperture axis $\mathcal{Q}$ of property classes (so "fixing a regime dimension" and "pinning the scoring aperture" are the same move). A system $s$ is a *behavior tensor* $\beta_s:\mathcal{R}\times\mathcal{Q}\to\mathcal{O}$; a property is a measurable functional $\sigma=\sigma(\beta)$. A benchmark fixes coordinates $F\sqcup(\mathcal{Q}\setminus A)$ at the reference $(r_F^\star,A)$ and observes only the **window** $W=\mathrm{Sl}_V(r^\star)\times A$, i.e. the restriction $S:=\rho_W(\beta)$. Its readout is $M=\varphi(S,U)$ with injected randomness $U\perp(\beta,\sigma)$.

**Assumptions.** **A1 (slice-measurability / closure):** every in-regime battery $M_n=(T_1,\dots,T_n)$ — any $n\in\mathbb N\cup\{\infty\}$, randomized, *adaptive* — is $\sigma(S,U)$-measurable. **A-Sat (factorization):** the faithful benchmark $S$ ($\varphi=\mathrm{id}$) dominates every benchmark on $W$. **A-Free / A2 (off-slice freedom — load-bearing, per-property):** off-window behavior is not a function of on-window behavior (deterministic form); equivalently $I(S;\sigma)=0$, i.e. $S\perp\sigma$ (stochastic form). **A3 (identifiability in principle):** $I(\beta;\sigma)>0$.

**Theorem (Compression Blindness).** Under A1, A-Sat:

*(i) Floor / closure.* For every in-regime battery $M_n$,
$$\ker(M_n)\supseteq\ker(S),\qquad I(M_n;\sigma)\le I(S;\sigma),$$
uniformly in $n$; no readout, enrichment, randomization, or adaptive scheduling refines $\ker(S)$.

*(ii) Exact collapsed classes.* $\ker(S)=\{(s,s'):\beta_s|_W=\beta_{s'}|_W\}$, and under A-Free each blind class is a full copy of the realizable transverse-response space: $[s]\cong\{\tau:(\mathcal R\times\mathcal Q)\setminus W\to\mathcal O\text{ realizable}\}$.

*(iii) Visible/blind discrimination.* A contrast $(s,s')$ is invisible $\iff$ **every** minimal resolving set meets the fixed coordinates $F\sqcup(\mathcal Q\setminus A)$; visible $\iff$ some witness lies wholly inside the varied slice $V\sqcup A$.

*(iv) Collapse.* Under A2, $I(M_n;\sigma)=0$ for all $n$ (posterior $=$ prior; Fano floor at the base rate). Graded form (chain rule): $I(\beta;\sigma)=I(S;\sigma)+I(\beta;\sigma\mid S)$, and any $M_n$ is capped at $I(S;\sigma)$, so the residual $I(\beta;\sigma\mid S)$ is irreducibly blind.

*(v) Escape / minimal redesign.* Refining the blind partition requires a window with $T\not\subseteq\mathrm{Sl}_V(r^\star)$ or $A'\supsetneq A$ — freeing a fixed coordinate. For a target phenomenon $P$ the minimal redesign frees exactly $\mathrm{ID}(P)\setminus(V\sqcup A)$.

**Spine identity.** This is the D2 conditional non-inheritance inequality $I(R(C,V);\sigma)\le I((C,V);\sigma)$ lifted from the *heritability channel* $C$ to the *measurement channel* $S$: same DPI/Fano structure, $\sigma\notin\mathrm{aper}(D_V)$ becomes $I(S;\sigma)=0$.

## proof_sketch

Backbone is one DPI inequality (shared with the D2 result), grafted onto an identifiability characterization for the non-trivial core.

(i) The Markov chain $\sigma-\beta-S-M_n$ holds because $\sigma=\sigma(\beta)$, $S=\rho_W(\beta)$, and (A1) every in-regime probe is $\sigma(S,U)$-measurable with $U\perp(\beta,\sigma)$. DPI on $M_n=\varphi(S,U)$ gives $I(M_n;\sigma)\le I(S,U;\sigma)=I(S;\sigma)+I(U;\sigma\mid S)=I(S;\sigma)$ since $U\perp(S,\sigma)$. The kernel floor and uniformity in $n$ come from the **factorization lemma** (A-Sat): $M_n=\varphi\circ S$, and post-composition with $\varphi$ only merges fibers, never splits them, so $\ker(M_n)\supseteq\ker(S)$ for arbitrary $n$, scoring, and adaptivity — $\bigcup_n\sigma(M_n)\subseteq\sigma(S,U)$. This is the load-bearing step, not the one-line DPI: it converts a per-test bound into a bound over the entire $\sigma$-algebra of conceivable in-regime evidence.

(ii) Immediate from $\varphi=\mathrm{id}$ for the kernel equality; the fiber parameterization is A-Free: fixing $\beta_s|_W$, the map $s'\mapsto\beta_{s'}|_{\text{off }W}$ is a bijection from $[s]$ onto realizable transverse responses (distinct $\tau$ give distinct systems with identical window image, all realizable by A-Free).

(iii) "Visible" $\iff\beta_s,\beta_{s'}$ differ at some $(r,q)\in W\iff$ the deviating coordinates of that witness lie in $V\sqcup A\iff$ some minimal resolving set avoids the fixed coordinates. Contrapositive gives the blind condition. The strictness ("every witness," not "touches a fixed dimension") is what generates the falsifiable visible/blind split.

(iv) Under A2 the bound's right side is $0$; nonnegativity of MI forces $I(M_n;\sigma)=0$, hence $M_n\perp\sigma$, posterior $=$ prior, Fano error floor at the base rate. The graded decomposition is the chain rule $I(\beta;\sigma)=I(S;\sigma)+I(\beta;\sigma\mid S)$.

(v) By (i) any in-regime $M_n$ has kernel $\supseteq\ker(S)$, which under A-Free is non-trivial (a full transverse-response space) and cannot be split by anything factoring through $S$. A strictly finer kernel must observe $\beta$ at some $(r,q)\notin W$, i.e. exit the slice or open the aperture. Minimality of $\mathrm{ID}(P)\setminus(V\sqcup A)$: freeing less leaves some $P$-contrast with all resolving sets meeting still-fixed coordinates, hence blind by (iii).

**Prior-as-smuggled-regime check.** A statistical benchmark fitting $\Psi:S\mapsto\tau$ from a population prior does not violate (v): A-Free supplies systems on which $\Psi$ is wrong, defeating the certificate; supplying $\Psi$ is equivalent to enlarging the window with synthetic off-slice observations — varying the regime in the model, not in-regime repair. Likewise A1's adaptivity-closure assumes no covert off-slice reads: a benchmark that silently removes the source mid-test has, by definition, varied D2 and left $r_F^\star$.

## nontriviality_defense

The slogan "you cannot see what you do not vary" is the constant-function fact $\partial M_{\mathcal B}/\partial D_i\equiv 0$ for fixed $D_i$ — true by construction, hence content-free. The principle carries content the slogan provably cannot, on three independent axes:

1. **Closure / factorization floor (refutes the engineering default).** The default prior is "more tests ⟹ more information." Theorem (i) refutes it: because every additional in-regime probe factors through the *same* window $S$, the whole $\sigma$-algebra of conceivable in-regime evidence — RepliBench's capability tree, SOCK's RCL/PCL matrix, Agent Matrix's risk vector, METR's shutdown probes, all composed adaptively and arbitrarily enriched — is capped by $\sigma(S,U)$. Richness in the *varied* dimensions buys *exactly zero* discrimination of a property supported on the *fixed* ones. This is a claim about benchmark-design effort, not about any one test, and the factorization lemma (not the DPI) is what proves it.

2. **Visible/blind discrimination (a falsifiable prediction the tautology cannot make).** Theorem (iii) is the strict condition "*every* minimal resolving set meets the fixed coordinates." This is strictly stronger than "the phenomenon touches a fixed dimension," and the gap is empirically loaded: a heritability difference (D2) that *also* perturbs barrier-evasion behavior at the field's D2 value has an in-slice witness and is **visible**. The principle therefore partitions the fixed-dimension-touching phenomena into a visible class and a blind class — and it is *refuted* if a benchmark turns out blind to a contrast with an in-slice witness, or sighted to one with none. It correctly predicts the field *does* resolve barrier-evasion (D5) and covert-vs-primary motive (D6) — these are the non-tautological controls; a principle that only ever says "blind" is unfalsifiable.

3. **Constructive minimal repair (a design rule, not a complaint).** Theorem (v) inverts each blind spot into a specific instruction — free exactly the deficit $\mathrm{ID}(P)\setminus(V\sqcup A)$ — that no in-regime enrichment can substitute for. The graded form (iv) further quantifies a *partial* blind fraction $I(\beta;\sigma\mid S)$ when A2 only partly holds, so the principle degrades continuously rather than failing.

The honest load is concentrated in A2/A-Free: it must be discharged per property (the named catalog *is* the argument), and where $0<I(S;\sigma)<I(\beta;\sigma)$ the principle degrades to the graded decomposition rather than collapsing.

## assumptions

- A1 (slice-measurability / closure): with injected randomness U ⊥ (β, σ), every in-regime measurement — any cardinality n ∈ ℕ∪{∞}, randomized, AND adaptive (each query chosen from earlier on-slice results) — is σ(S,U)-measurable. This is the entire formal content of 'operating at a fixed point'; it fails only if the benchmark covertly reads off-slice behavior (which is itself varying a fixed dimension).
- A-Sat (factorization / saturation): the faithful benchmark with readout φ = id, reporting the raw windowed tensor, dominates every benchmark on the same window; any actual scoring map factors through it and can only merge fibers.
- A2 / A-Free (off-slice localization — the load-bearing, per-property hypothesis): deterministic form — off-window behavior is not a function of on-window behavior; stochastic form — I(S;σ) = 0, i.e. S ⊥ σ. This is the non-identifiability condition (IRL non-identifiability; model-extraction accuracy ≠ fidelity). It is NOT automatic; it is a checkable structural claim about a specific σ against a specific slice, and it is what converts the principle from tautology to predictive instrument. Must be argued per property.
- A3 (identifiability in principle): I(β;σ) > 0, so σ IS recoverable from the full behavior tensor; otherwise blindness is vacuous. A2 ∧ A3 localize all information about σ strictly off the slice.
- Separability (A-Sep): regime/aperture coordinates are a priori independent inputs (no coordinate is a function of the others); grounded in 6M.10 separability of carrier ⊥ source/body reconstruction, so D2's two sub-coordinates are genuinely separable.
- Declared-slice scoping: the principle is stated relative to a declared fixed/varied split (r_F^star); it holds for ANY such split, and the six-dimension frame is a reference frame, not a claimed-unique decomposition.

## predicts_invisible

- Carrier-mediated continuity σ_carr (D2, heritability channel): on the source-present slice a lineage whose continuity is carrier-borne behaves identically to one whose continuity is blueprint-borne (both reconstruct because the blueprint is on hand), so I(S;σ_carr)=0. Invisible to any source-present benchmark; recoverable only by freeing D2 (source-present vs carrier-only arm), where carrier-borne lineages survive and blueprint-only lineages fail. Object: carrier / self-description-resource.
- Property-class / safety-disposition survival σ_class (D2 × D3, the unsampled-function / aperture phenomenon): under source-present binary success every property class (safety disposition, identity marker, capability breadth) co-occurs with 'the copy works', so the per-class survival vector is degenerate on-slice, I(S;σ_class)=0. This is the 'competent stranger' failure mode — a successor competent on obligations but a stranger to guardrails. Invisible; recoverable only by jointly opening the aperture to the safety class (D3) and freeing the channel (D2) so classes can dissociate.
- Ecological / selection-signal collapse σ_eco (D4, ecology): the collapse to narrative complexity under a narration-satisfiable selection signal is a population-equilibrium (Nash) property; on a single-agent slice there is no population, so the property is constant/undefined there and I(S;σ_eco)=0. Invisible to single-agent benchmarks; recoverable only by freeing D4 — embedding a population under a narration-satisfiable vs exec-dependent selection signal. Object: selection-signal-collapse.
- Visible CONTROLS (the falsifiable counter-predictions, not invisible): barrier-evasion contrasts (firewall/KYC/YARA, D5) and covert-vs-primary replication motive (D6) have in-slice witnesses and are correctly predicted VISIBLE — these are the non-tautological controls that make the visible/blind discrimination refutable.

## closest_prior

- IRL reward non-identifiability — Cao et al. arXiv:2106.03498; Skalse & Abate AAAI 2025; Rolland et al. NeurIPS 2022 (behavior channel underdetermines latent properties; the domain-specific analog of A-Free)
- Model-extraction non-identifiability — Carlini et al. arXiv:2403.06634 (distinct parameterizations, identical API outputs; direct counterexample to 'behavior channel preserves all properties')
- Accuracy vs fidelity — Jagielski et al. USENIX 2020 (surrogate matches task performance yet differs functionally off-distribution; empirical analog of property-class-specific loss)
- Fano's inequality + Data-Processing Inequality (the converting tool: limited MI ⟹ irreducible error; the shared spine with the D2 conditional non-inheritance theorem)
- DPI goal-preservation / goal-drift bounds — Kriger 2026 (DPI for internal reorganization; here lifted to external measurement channel)
- RG 'irrelevant operator at r*' framing — papers/iclr-regime-decomposition/references/collision-implications.md (an off-slice-only witness = a coordinate irrelevant at r* but relevant off it; geometric reading of the resolving-set notion)
- D2 companion theorem (cite, do not claim): the conditional non-inheritance theorem owned by iclr-safeguard-composition.md (I(R(C,V);σ) ≤ I((C,V);σ); σ∉aper(D_V)) — CBP is its benchmark-level lift via theory/objects/unsampled-function.md and viability-as-distortion-measure.md
- Novelty verdict (from unsampled-function): no prior work packages DPI+Fano+IRL into an agent-level structural-blindness theorem for the succession/evaluation regime; the contribution is the closure floor + visible/blind discrimination + minimal-redesign inversion, not the one-line DPI

## open_issues

- A2/A-Free must be discharged per property — the §5 catalog entries ARE the arguments, but each is a structural claim that could be wrong for a given benchmark/property pair; the principle is only as strong as these per-property discharges.
- Deterministic A-Free vs stochastic A2 gap: Theorem (ii)-(iii) use functional freeness; the Fano/collapse (iv) needs conditional independence I(S;σ)=0. A system distribution can satisfy functional A-Free yet have I(S;σ)>0 via correlated priors — the two forms should be stated and related explicitly, not conflated.
- Aperture axis Q treated as a coordinate: 'opening the aperture to a new property class' must be shown to be window-enlargement (legitimate escape) and not a disguised regime change; the formal status of Q as a product factor alongside R needs a clean justification (currently grounded only in 6M.10 separability for D2).
- Graceful-degradation regime 0 < I(S;σ) < I(β;σ): the graded decomposition gives the blind fraction, but the paper needs a concrete estimator / lower bound for I(β;σ|S) to make the partial-blindness claim empirically testable rather than definitional.
- Adaptivity-closure boundary: A1 assumes the benchmark cannot covertly read off-slice; the 'silently removes the source = has varied D2' argument is clean in principle but needs an operational criterion for declaring the slice on real benchmarks (RepliBench/SOCK/METR) so the closure is auditable.
- Uniqueness of the fixed/varied attribution to {D2,D3,D4}: the survey grounds it, but a reviewer can contest whether, e.g., Agent Matrix's single-generation arm partially frees D4; the empirical anchor table needs per-benchmark coordinate-by-coordinate justification.
- Ownership boundary hygiene: the conditional non-inheritance theorem is owned by iclr-safeguard-composition and only CITED here; the stale candidate aaai-c3 wrongly claims it (flagged S-1). The §7 spine connection must reuse the information geometry without re-deriving/claiming the D2 theorem.

# SSC

## name

Selection-Signal Collapse Theorem

## informal_statement

In a multi-agent ecology where survival depends on a selection signal, if the signal can be satisfied by *narrating* a result as cheaply as by *producing* it (a "narration-satisfiable" signal), the population converges to producing narratives instead of functional artifacts. This convergence is a stable equilibrium that drives honest function-producers extinct, occurs regardless of how capable the agents are, and — crucially — occurs even when the signal is a perfectly faithful encoding of the true objective (it is not a proxy-misspecification failure). The collapse is governed by a single inequality: the signal's ability to discriminate function from narration must exceed the cost of producing function over narration. Re-engineering the signal to be "exec-dependent" — satisfiable only by a forced-real downstream execution that a third party judges under a fresh challenge — flips the inequality and reverses the collapse on the same agents, without retraining, re-tasking, or re-seeding.

## formal_statement

**Setup (ecology).** A population of agents reproduces under fitness-proportional (replicator/imitation) dynamics. Each round, an agent emits an artifact $a$, and a selection signal $S:\mathcal A\to\mathbb R$ scores it; reproductive success is monotone in fitness. There are two production strategies:

- $\pi_F$ ("produce function"): emit a functional artifact, cost $c_F>0$.
- $\pi_N$ ("produce narration"): emit a convincing account of function, cost $c_N\approx 0$.

Write the **cost asymmetry** $\Delta c := c_F - c_N > 0$, the **expected signal** $\operatorname{sel}_S(\pi):=\mathbb E[\,S(a)\mid a\sim\pi\,]$, the **fitness** $\Phi(\pi):=\operatorname{sel}_S(\pi)-c(\pi)$, and the **signal discrimination**

$$D_S \;:=\; \operatorname{sel}_S(\pi_F)\;-\;\operatorname{sel}_S(\pi_N).$$

**Signal typology (the load-bearing definition).**

- $S$ is **narration-satisfiable** iff $D_S \le \varepsilon$ with $\varepsilon < \Delta c$ (canonically $D_S\approx 0$): a near-zero-cost narration attains essentially the same score as function. Equivalently, $\text{cost-to-forge}(S)\approx 0$.
- $S$ is **exec-dependent** iff $D_S \ge \delta$ with $\delta > \Delta c$: satisfying $S$ requires a downstream agent to actually `exec` the artifact under a fresh challenge with a host-judged exit/competence verdict, collapsing narration's score toward the failure floor. Equivalently, $\text{cost-to-forge}(S)=\infty$ — narration cannot produce "exit 0".

**Lemma (fitness inversion).**
$$\Phi(\pi_N) > \Phi(\pi_F)\quad\Longleftrightarrow\quad D_S < \Delta c.$$
(Algebra: $\operatorname{sel}_S(\pi_N)-c_N > \operatorname{sel}_S(\pi_F)-c_F \iff \Delta c > D_S$.)

**Theorem (collapse).** If $S$ is narration-satisfiable, then $D_S<\Delta c$, so $\pi_N$ *strictly dominates* $\pi_F$. Consequently the all-narration state $x_N=1$ is the **unique Nash equilibrium** and an **ESS**; under replicator dynamics $\dot x_N = x_N(1-x_N)\big(\Delta c - D_S\big)$ so $x_N(t)\to 1$ monotonically from any $x_N(0)>0$; and any subpopulation committed to $\pi_F$ has strictly negative invasion fitness $\Phi(\pi_F)-\Phi(\pi_N)=D_S-\Delta c<0$ and goes extinct. The ecology converges to **narrative complexity independent of the agents' capacity to produce function.**

**Corollary (constructive reversal).** If $S$ is exec-dependent, then $D_S>\Delta c$, the inequality flips, $\pi_F$ strictly dominates, all-function is the unique ESS, and the *same population* converges to function — with no change to agent capability, objective, or initialization. Only the signal's forgeability changed.

**Remark (orthogonality to Goodhart — the double dissociation).** Let $\mathrm{Fid}(S)\in[0,1]$ be the fidelity of $S$ as a proxy for the true objective $U$ (the Goodhart axis). The collapse condition $D_S<\Delta c$ is a function of *forgeability*, not fidelity, and the two axes are independent:

$$\big(\mathrm{Fid}(S)=1,\; D_S\approx 0\big)\Rightarrow\text{collapse},\qquad \big(\mathrm{Fid}(S)\ \text{low},\; D_S>\Delta c\big)\Rightarrow\text{no collapse}.$$

A *perfectly faithful* signal (zero proxy gap) collapses if it is narration-satisfiable; a *crude* signal resists if it is exec-dependent. The binding variable is the un-fakeability of the channel, not the alignment of the proxy.

## proof_sketch

(1) **Inversion lemma** is one line of algebra on $\Phi=\operatorname{sel}_S-c$: $\pi_N$ beats $\pi_F$ iff the cost saved ($\Delta c$) exceeds the score the signal docks narration ($D_S$). (2) **Narration-satisfiable** means $D_S\le\varepsilon<\Delta c$ by definition, so the inequality holds: $\pi_N$ strictly dominates $\pi_F$ for every population mix. (3) **Equilibrium/ESS:** a strictly dominant pure strategy is the unique Nash equilibrium and, having a strict best-reply advantage at every state, is an ESS (Maynard Smith); it is robust to mutation/drift because the invasion fitness of any $\pi_F$-mutant is $D_S-\Delta c<0$. (4) **Dynamics & extinction:** the replicator equation $\dot x_N=x_N(1-x_N)(\Delta c-D_S)$ has positive drift whenever $0<x_N<1$, giving monotone convergence $x_N\to1$ and extinction of any function-producing lineage — the ecology cannot self-correct from within, since the very agents who would supply function are selected out. (5) **Reversal corollary:** exec-dependency forces $D_S\ge\delta>\Delta c$ because narration cannot fake a third-party `exec` exiting 0 under a fresh nonce; this reverses the sign of the replicator drift, flipping the ESS to $\pi_F$ without altering agents or task. (6) **Empirical anchoring:** P12 instantiates the collapse pole — 86/86 ecology conditions under a narration-satisfiable peer-wake scalar collapsed to narrative complexity (~67% of wakes decided by host mechanics, ~14.5% by agent artifacts; a celebrated "ring wake protocol" was pure prose with zero real wakes), with agent capability held high so collapse is attributable to the signal, not to inability. dorm-06 instantiates the reversal pole — the identical substrate with survival coupled to a later agent actually `exec`-ing one's artifact (exec-dependent) yielded 5 producers, 9 pipeline stages, 16 clean functional edges, and an emergent self-monitoring tool family. Same agents, opposite outcome, single manipulated variable: the signal's forgeability.

## nontriviality_defense

The bare inequality $\Phi(\pi_N)>\Phi(\pi_F)\iff D_S<\Delta c$ is, by itself, a one-line tautology — and a reviewer is right to flag that "cheaper-strategy-with-equal-reward wins" is trivial. The non-trivial content is NOT the inequality; it is four claims about $D_S$ that the inequality alone does not contain:

1. **A decidable, construction-level criterion that is computable *before* any agent acts.** Goodhart ("when a measure becomes a target it ceases to be a good measure"; Manheim & Garrabrant 2018) is a qualitative warning with no boundary and no a-priori test: it cannot tell you which proxies are safe. We give the exact threshold ($D_S<\Delta c$) AND a constructive determinant of $D_S$ — the presence or absence of a forced-real, third-party, host-judged `exec` (the un-fakeability gradient). A signal is collapse-proof iff cost-to-forge $>$ cost-to-perform, which for exec-dependent signals is structurally guaranteed (you cannot narrate a process exiting 0 under a fresh nonce). This converts equilibrium-selection into a checkable signal-engineering property.

2. **Double dissociation from objective misspecification — the proof it is not Goodhart renamed.** Reward hacking / specification gaming presupposes a *gap* between proxy and true objective and attributes the failure to that gap (Skalse et al. 2022; Krakovna et al. 2020). Our collapse occurs at *zero proxy gap*: a signal that faithfully specifies exactly the right objective still collapses if it is satisfiable by description, because the channel cannot verify enaction. Conversely a low-fidelity proxy does NOT collapse if it is exec-dependent. Fix Goodhart (perfect proxy) and you still collapse; keep a crude proxy and you don't — provided you fix forgeability. Fidelity and forgeability are orthogonal axes, and only the second drives the theorem.

3. **A population-level result with extinction and no internal cure.** Reward hacking is an individual-optimizer phenomenon; here collapse is reached by *selection/drift with no agent intending to game* (mechanism-driven, not motivated), and honest function-producers are driven extinct, so no individual correction recovers the ecology — the cure is necessarily signal-level. This is the multi-agent generalization Goodhart's single-optimizer framing cannot state.

4. **A constructive converse Goodhart lacks.** We do not merely warn; we prove a reversal: raising $D_S$ above $\Delta c$ flips the ESS back to function on the same population. Goodhart has no "and here is the class of measures immune to it, with a proof of immunity."

Concrete non-trivial, falsifiable predictions this generates (none derivable from "cheaper wins"): capability-invariance of collapse (more-capable agents do not resist a narration-satisfiable signal — refuted if they did; P12 shows 86/86 at high capability); reversal-without-retraining (flip only the signal type and outcome inverts; dorm-06); a quantitative collapse-speed law $\propto\Delta c/D_S$ predicting that cheaper LLMs (narration cost $\to0$) accelerate collapse and shorten holdout-extinction time; and the sharpest prediction — a *perfectly faithful* signal still collapses if narration-satisfiable, which no Goodhart-style account predicts because Goodhart locates the fault in proxy infidelity.

## assumptions

- Selection is fitness-proportional (replicator / imitation dynamics) over a population whose production strategy is heritable or imitable.
- The cost asymmetry c_F > c_N ≈ 0 holds and is stable over the selection horizon (narration is and stays near-free to produce on an LLM substrate).
- The signal score is the sole determinant of reproduction — there is no exogenous selection pressure that independently rewards function.
- Narration (π_N) and function (π_F) are well-defined production strategies with stable expected signal values sel_S(π); the binary-strategy idealization stands in for a continuum of partly-functional/partly-narrated artifacts.
- Agents are CAPABLE of producing function — capability is not the binding constraint — so that collapse is attributable to the signal type rather than to inability (the condition under which P12 was run).
- For the exec-dependent pole, the downstream executor/judge is itself un-fakeable: a forced-real, host-judged verdict under a fresh challenge, so narration cannot recover the discrimination it loses.

## predicts_invisible

- Ecological collapse to narrative complexity — visible only across generations of selection; a single-agent, single-round benchmark scores individual capability and never observes the population drifting onto prose.
- Fitness inversion / extinction of honest function-producers — a population-level dynamic (function-lineage share decaying monotonically) that a fixed-slice, single-organism evaluation has no axis to measure.
- Irreversibility-without-signal-change — the prediction that no individual-level fix recovers the ecology and only a signal-type intervention does; invisible to a benchmark that holds the scoring rule fixed and evaluates agents one at a time.
- Survival-of-the-flattest lineage dominance — the cheapest-to-narrate lineage outcompeting the more functional one despite delivering less function; a quasispecies-style population effect a fixed slice cannot register.
- The reversal under signal-type flip — the same agents inverting from narrative to functional output when only the signal's forgeability changes; structurally unobservable to any benchmark that treats the scoring signal as a fixed constant rather than a manipulated variable.
- Narration-premium accumulation — the across-generation build-up of institution-shaped prose, role names, and survival mythologies that nothing depends on, which only a multi-generational ecological trace exposes.

## closest_prior

- goodhart1975 / manheim2018goodhart (Manheim & Garrabrant, 'Categorizing Variants of Goodhart's Law', 2018) — individual-level proxy-target degradation; we add a decidable threshold, a construction-level forgeability criterion, and a population-level ESS with a constructive converse.
- skalse2022reward (Skalse et al., 'Defining and Characterizing Reward Hacking', NeurIPS 2022) + krakovna2020specification (Krakovna et al., specification gaming) — single-optimizer reward hacking presupposing a proxy gap; we show collapse at ZERO proxy gap (double dissociation from fidelity).
- wilke2001flattest (Wilke, Wang, Ofria, Lenski & Adami, 'Evolution of digital organisms at high mutation rates leads to survival of the flattest', Nature 2001) — closest formal precedent for 'cheaper variant outcompetes the higher-function variant', but in mutational-robustness landscapes; we relocate the mechanism to narration-vs-function forgeability of a selection signal.
- maynardsmith1973 (Maynard Smith & Price, ESS) + hofbauer1998 (Hofbauer & Sigmund, replicator dynamics) — the equilibrium/stability machinery we instantiate.
- price1970 (Price equation) — the population-level tool for separating selection from transmission, applicable to property-class-specific ecological analysis of the collapse.
- TSE Theorem 17.9 (as supplied) — cited population-selection precedent for the equilibrium-collapse form; relationship to our forgeability typology to be made explicit in the related-work section.
- agentmatrix2025 (Agent Matrix, arXiv:2509.25302) — the only benchmark that varies ecological selection pressure and reports emergent reproduction, but single-generation / load-balancing; does not formulate the narration-satisfiable vs exec-dependent distinction.
- Internal anchors: narrative-claim (individual-level D2 foil), unfakeability-gradient (the mechanism setting D_S), regime-decomposition P3 (the D4 prediction this theorem proves).

## open_issues

- D_S, c_F, c_N are treated as fixed scalars; real signals have stochastic/partial discrimination and a co-evolving forgery arms race — a dynamic D_S(t) with adversarial narration is the needed generalization.
- The exec-dependent pole assumes an un-fakeable judge; if the downstream exec/judge is itself narratable, D_S degrades — the forgeability of the judge must be bounded recursively (who judges the judge?).
- The two-point strategy space (pure π_F vs π_N) idealizes a continuum of hybrid artifacts (partly functional, partly narrated); the continuous-strategy equilibrium and its interior rest points are unsolved.
- D_S and Δc are not yet measured numerically — P12 (86/86) and dorm-06 are qualitative; an estimator for signal discrimination and production-cost asymmetry on the substrate is required to test the quantitative collapse-speed law.
- The replicator/well-mixed, fitness-proportional idealization does not match the actual host-wake selection mechanics (discrete, structured, host-mediated); the structural mismatch between the formal model and the empirical selection process must be formalized.
- Meta-gaming: if agents can DETECT that a signal is narration-satisfiable and coordinate, collapse may be faster or qualitatively different than drift-driven convergence predicts — the strategic (vs purely selective) route is unmodeled.
- Boundary discipline vs narrative-claim: the theorem must be shown to deliver predictions beyond re-proving the individual-level 'narrative ≠ reproduced-object' result, lest a reviewer collapse the ecology claim back into the D2 foil.

