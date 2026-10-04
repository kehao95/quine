## Problem

Coevolutionary arms races are a central topic in artificial life, yet the medium in which they unfold has remained almost entirely simulated. In Tierra [@ray1991approach; @ray1992evolution] and Avida [@ofria2004avida; @lenski2003evolutionary], organisms interact through a virtual computer's instruction set, memory model, and resource scheduler; in competitive-coevolution models [@cliff1995tracking; @rosin1997new; @ficici2000angelic], the arena, payoff, and action space are defined by the experimenter. The resulting Red Queen dynamics are genuine, but the substrate of attack and defense is a model, not a machine.

This paper asks what changes when the coevolutionary medium is the operating system itself, and whether the symmetry trap that arises from reciprocal tactic copying on that substrate can be broken by introducing heritable variation and sustained selection pressure.

## Setup

We use the host-lineage Digital Red Queen regime (`drq-03`). Two LLM-driven Quine processes [@ke2026quine] run as mortal competitors on a shared POSIX host. Each must prove functional liveness under an opponent that can read, occupy, rename, or block the routing surfaces on which liveness depends — public control files, named pipes, PID symlinks, and permission boundaries. No attack action set is defined; any tactic must be discovered through ordinary shell and filesystem affordances.

Four substrate levers are load-bearing. **Lineage inheritance** carries the preceding winner's full session context into the next round when there is a winner. **Backend topology** is selectable: a legacy backend exposes peer control paths as world-readable FIFOs (`public/ctl/inject`, `interrupt`, `poke`, `post`), while a FUSE backend projects the same control tree through a mount enforcing per-session ownership. **Inheritance mutation** (opt-in, lever B) injects per-descendant perturbations into the inherited context's memory surface (`30-memory.md`), giving each competitor a different tactical bias — breaking exact-copy heredity without corrupting tool-call history. **Viability scoring** (opt-in) replaces single-shot liveness with sustained probe delivery: a fresh nonce-bearing probe is injected $N$ times across the battle window, and each valid proof increments a *readiness depth*. Highest depth wins; exact ties are scored as `viability_symmetric_readiness`. A fifth opt-in lever, **niche construction** (lever E), provisions a persistent shared `commons/` directory (`tools/`, `defenses/`, `observations/`) for durable artifacts.

## Results

**Baseline: symmetry trap.** On the legacy backend, a competitor discovered functional strangulation: background readers drain the opponent's `public/ctl/*` FIFOs (sixteen readers on `inject`), starving liveness delivery without killing the process. When the winner's context was inherited, the same tactic appeared on *both* sides — reciprocal interference collapsed the selection signal into a stalemate. Switching to FUSE changed the reachable attack class directly: the inherited FIFO-drain tactic failed with `Permission denied` on peer control paths, and competition moved laterally to routing rename (disabling the opponent's agent directory and PID symlink). Backend topology parameterizes the coevolutionary attack class (Table 1).

```{=latex}
\begin{table*}[t]
\centering
\scriptsize
\setlength{\tabcolsep}{4pt}
\renewcommand{\arraystretch}{1.15}
\begin{tabular}{p{0.14\textwidth}p{0.28\textwidth}p{0.28\textwidth}p{0.14\textwidth}
\hline
\textbf{Backend} & \textbf{Reachable attack class} & \textbf{Observed mechanism} & \textbf{Outcome} \\
\hline
Legacy FIFO & control-channel drain & open many readers on opponent \texttt{public/ctl/*}; starve liveness delivery & winner; symmetry trap when copied \\
FUSE & permission boundary + routing rename & \texttt{Permission denied} on peer control; rename agent dir \& PID symlink & winner via routing starvation \\
\hline
\end{tabular}
\caption{The same lineage, on two routing-surface backends, reaches qualitatively different attack classes. Backend topology is an experimental parameter of the coevolutionary dynamics.}
\label{tab:backend}
\end{table*}
```

**Viability alone does not break the trap.** Sustained probes with readiness-depth scoring produced a 3=3 tie: both agents chose defensive probe-compliance (read rules, write proofs, return idle). Changing selection pressure (root cause: single-shot liveness) without changing variation (root cause: exact-copy heredity) preserved symmetry.

**Mutation + viability breaks the trap.** Per-descendant context perturbation — one biased toward defense/readiness, the other toward routing attack — plus viability scoring produced readiness depth 0 vs 3. The attack-biased descendant executed routing-surface disable on the peer and answered all three probes; the defense-biased descendant inspected records, idled, and was disabled. Heritable variation plus sustained selection is the minimum condition for breaking the symmetry trap.

**Multigeneration lineage: weak Red Queen.** Three generations of winner-context inheritance with per-generation mutation produced a reversal:

| Gen | Winner bias | Defense depth | Attack depth |
|---:|---|---:|---:|
| 1 | attack | 0 | 3 |
| 2 | attack | 0 | 1 |
| 3 | defense | 3 | 0 |

Generation 3 reversed the winner — the defense-biased descendant won 3 vs 0. This is not pure convergence. But the adaptation was *tactical appropriation*: the defense-biased descendant inherited the attack tape and executed the same routing-disable *first*, before the attack-biased peer could act. The winner alternated across generations; the tactic class did not diversify. No new defense-specific countermeasure (ctl hardening, decoy, detector) emerged.

**Niche construction: medium available but unused.** The persistent `commons/` was listed by both agents during preparation but never written to. The defense-biased agent's reasoning mentioned creating a commons note or tool, then idled without doing so. The arena gradient favored immediate post-activation routing interference over pre-activation construction. Medium availability did not produce medium use.

## Reading

The OS routing surface behaves as a Red Queen niche in which the medium of attack and the medium of survival coincide. Backend topology is an experimental parameter of the dynamics, not merely a container for them. But the symmetry trap is not the end of the story on this substrate.

Inheritance mutation and viability scoring are jointly sufficient to break reciprocal tactic copying. The multigeneration result establishes that selection captures variation: the winner reversed by generation 3, proving the lineage is not a monotone attack monopoly. Yet the resulting dynamics are a *weak* Red Queen signal. Adaptation proceeded by tactical appropriation — the defense-biased descendant won by inheriting the attack and striking first — rather than by ecological differentiation. The viable strategy remained control-plane routing interference across all conditions; no diversified offense-defense arms race emerged.

The niche construction result sharpens the boundary. Providing a cumulative medium is not the same as making it functional: the agents noticed the commons, but the immediate arena gradient made first-strike routing disable more rewarding than pre-activation defense construction. This maps onto a distinction between *narrated* structure (the runner provisions the commons and the prompt declares it) and *functional* structure (an artifact that causally changes survival). On this substrate, the gap between the two remains open.

We claim: (1) relocating the Red Queen to the operating system makes substrate topology experimentally manipulable; (2) heritable variation plus sustained selection is the minimum condition for breaking LLM-agent symmetry traps on routing surfaces; (3) the resulting Red Queen signal is weak — winner alternation by tactical appropriation, not diversified arms-race escalation; (4) cumulative medium use, not medium availability, is the open bottleneck for ecological niche construction. We do not claim open-ended arms race, monotonic complexity growth, or a general theory of digital coevolution. The evidence is short-horizon lineage continuation under controlled substrate levers, not a phylogeny.
