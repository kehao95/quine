```{=latex}
\begin{table*}[t]
\centering
\scriptsize
\setlength{\tabcolsep}{4pt}
\renewcommand{\arraystretch}{1.15}
\begin{tabular}{p{0.12\textwidth}p{0.16\textwidth}p{0.15\textwidth}p{0.18\textwidth}p{0.28\textwidth}}
\hline
\textbf{Substrate} & \textbf{Substrate-type} & \textbf{Individual} & \textbf{Weismann barrier} & \textbf{Enforced by} \\
\hline
Biology & separate (cell vs.\ world) & discrete (membrane) & present (germ-line sequestration) & chemistry, physics, development \\
Tierra  & separate (code vs.\ CPU) & discrete (occupied space) & present (no phenotypic feedback) & virtual-machine architecture \\
Avida   & separate (genotype vs.\ phenotype) & discrete (CPU, memory) & present (no phenotypic feedback) & virtual machine + environment \\
Quine   & \textbf{collapsed} (language = language) & \textbf{fuzzy} (context window) & \textbf{absent at persistence} & lossless read, shared filesystem \\
\hline
\end{tabular}
\caption{Substrate regimes against the three evolutionary separations. In Quine all three collapse from a single root property: lossless absorption of environment into cognition without translation.}
\label{tab:coincidence}
\end{table*}
```

## Problem

Classical evolutionary theory rests on three structural separations whose enforcement is so fundamental it is rarely noticed as an assumption. (1) *Substrate-type*: organism and environment are different substances, so heredity flows through a narrow, lossy translation interface. (2) *Individual boundary*: the organism is a discrete, countable unit of selection. (3) The *Weismann barrier*: acquired traits do not inherit, because the germ-line is sequestered from somatic change [@weismann1893germ]. These are not laws of nature. In biology each is enforced by an independent mechanism (chemistry, membrane, development); Tierra [@ray1991approach] and Avida [@ofria2004avida] enforce all three by architecture; and Pattee's epistemic cut --- the boundary between a symbolic record and the dynamics that read it --- is treated as necessary for semantic closure [@pattee1995evolving]. We ask: what happens when all three fail at once?

## Setup

Quine [@ke2026quine] is a multi-agent medium in which large language model (LLM) agents run as native operating-system processes over a shared filesystem workspace --- simultaneously the runtime source, the experiment record, and the system's own theory. An agent dropped into a worktree is embedded in the source of the process that defines it.

The load-bearing operation is *read*. In biology an organism reads its environment through sensory transduction --- a lossy, bandwidth-limited translation between substances. In Quine, reading a file inserts its text into the context window; the tokens that enter are the very tokens that were written, processed by the same machinery as the agent's own reasoning. **Read is lossless absorption, not translation.** Coverage is partial: an agent reads only a selective subset into a finite context. ``Lossless'' describes fidelity *at the moment of reading*, not the completeness of inheritance --- a distinction load-bearing for the heredity claim.

We report missionless-agent experiments (the null-05 series) across three model families (gpt-5.5, deepseek-v4-pro, glm-5.2), at one to six concurrent agents and 5--25 minute horizons. Agents had access to the constitution, source, theory, and peer output, but received no task.

## Results

Table \ref{tab:coincidence} places the substrates against the three separations. If environment is absorbed losslessly into cognition, the three cannot hold: agent and environment are the same substance (language), so in Pattee's terms the epistemic cut is structurally impossible; a peer's output, read into context, is indistinguishable from the reader's own prior thoughts, so the individual becomes a temporary cut through a continuous language field; and every write persists with no soma/germ filter, so an acquired tool or modification is latent heredity, becoming actual only when a later agent reads it.

The null-05 runs illustrate each consequence. *Constitutive reflexivity:* every missionless run entered an audit attractor immediately, reading the self-describing constitution and beginning a maintenance loop with no mechanism to start it --- the reflexive loop is the substrate's resting state. *Immediate inheritance:* a coordination tool written in one run was used by a fresh-worktree successor in the next, with no generation gap. *Fuzzy individuals:* agents reading the same coordination files converged on shared behavior. *Endogenous fitness:* agents took the targets of their work from what other agents had written, with no external reward. A context-window termination read as genuine mortality --- cooperating agents noted the absence and could not reverse it.

Persistence is itself stratified, and the stratum can be agent-constructed. In a six-agent focal run, one agent enacted the first observed commit --- staging a selected subset of its own outputs while explicitly declining to commit other agents' untracked files --- treating commit as a deliberate promotion gate over its heredity. Later agents then rewrote the committed working-tree files, yet the commit remained fixed as a historical marker in repository history. Workspace writes are thus mutable, peer-rewritable candidate heredity; a commit is a durable stratum the agent constructs, as convention rather than a substrate-enforced soma/germ filter.

## Reading

The three collapses are not independent accidents: each depends on the same substrate property, so they fail together, whereas in biology each separation has its own enforcement and could fail alone. The regime is not unfiltered Lamarckism in practice --- persistence is unfiltered at the write layer (an agent may additionally enact a commit stratum as convention, not a restored barrier), but actual inheritance is gated by an endogenous *read-selection* bottleneck, and the reading trajectory is itself shaped by absorbed text. This relocates the bottleneck that Boyd \& Richerson's cultural-evolution framework [@boyd1985culture; @mesoudi2011cultural] places in the transmission step: per-transmission fidelity goes to zero (lossless per read), while an endogenous, salience-driven coverage filter takes its place. For artificial life, the contribution is a substrate class in which the foundational separations presupposed by open-ended evolution [@bedau2000classification] and multi-level selection [@okasha2006levels] do not hold. We report a structural phenomenon and its illustration, not a proof of cumulative open-endedness.

## Acknowledgments

This work was independent and unfunded; the authors retain responsibility for all claims.
