# Results

Table~\ref{tab:summary} summarizes the three contrasts. Each subsection reports a manipulated pressure, the organization it recruits, the contrast under which that organization collapses, and the strongest supported interpretation. Taken together, the three contrasts constitute a coordinated test of the separability claim: if these viability pressures are independently manipulable, then varying each one in isolation should recruit and collapse specific organizations while leaving the rest of the habitat fixed. Read together, the sections form a staged test: first continuity under real death, then cognition under persistence controls, then maintenance under disturbance-horizon controls.

\begin{table*}[t]
\centering
\small
\raggedright
\begin{tabular}{p{1.0cm}p{2.7cm}p{3.1cm}p{3.0cm}p{5.6cm}}
\hline
\textbf{Sec.} & \textbf{Evidence} & \textbf{Manipulated Pressure} & \textbf{Selected Organization} & \textbf{Strongest Claim} \\
\hline
3.1 & Continuity contrast; long-horizon mortality contrast & Mortality (continuity path; lifespan) & Lineage persistence & Real death recruits continuity when a path exists, with persistence shaped by operator availability and save discipline \\
3.2 & Persistence ablation; inherited-use controls & Persistence plus lineage continuity & Environmental cognitive organs & Persistent habitat recruits bounded reusable cognition, but stronger uptake weakens under ablation and does not yet generalize cleanly across shifted query families \\
3.3 & Maintenance-horizon comparison & Maintenance horizon & Externalized maintenance organs & Longer horizons recruit more detached post-exit maintenance layers under graded pressure \\
\hline
\end{tabular}
\caption{Summary of pressure-sensitive contrasts. Each row reports the manipulated viability pressure, the organizational phenotype it recruits, and the strongest interpretation the evidence supports.}
\label{tab:summary}
\end{table*}

## Real Death and Lineage Persistence

We begin with the substrate fact on which the rest of the paper depends: Quine imposes a real death boundary. When the continuity path is blocked, this boundary traps the lineage in single-life collapse. When a path is available, the same mortality pressure recruits lineage persistence.

In the primary control pair, the same search task was run once with and once without an available continuity operator (\texttt{exec}, the in-place program replacement described in Section~2.2). The exec-disabled replication remained trapped in two single-life sessions, each restarting the search from scratch and terminating without transition. The exec-enabled replication produced five sessions with three explicit self-replacement transitions \citep{Ke2026Quine}. The key difference is narrow but sufficient: with \texttt{exec} available, self-replacement occurred before exhaustion; without it, the lineage remained trapped inside single-life search. Figure~\ref{fig:ariadne} shows the contrast.

\begin{figure}[!htb]
\centering
\includegraphics[width=\columnwidth]{figures/fig-02-ariadne-contrast.pdf}
\caption{Continuity contrast. Exec-enabled: five sessions with three self-replacement transitions before exhaustion. Exec-disabled: two single-life sessions with no continuation. Same model, same task, same mortality pressure.}
\label{fig:ariadne}
\end{figure}

A longer-horizon mortality contrast shows what that lineage problem produces over extended succession. Under enforced external mortality (45-second wall-clock timeout with SIGTERM), twenty successive generations inherited only the durable filesystem. The lineage spontaneously invented a \texttt{progress.json} checkpoint at Generation~1 and maintained a stable four-field schema (\texttt{files\_checked}, \texttt{total\_checked}, \texttt{findings}, \texttt{next\_to\_check}) across all subsequent generations.

The more informative result is failure. In the 45-second run, Generations~8--10 entered a Sisyphean loop: each re-checked the same 15--16 files its predecessor had already surveyed, died before persisting progress, and left no trace for the next generation. That loop broke only when Generation~11 prioritized saving over searching. Doubling the lifespan to 90 seconds eliminated the loop, reduced wasted re-reads from 47\% to 9\%, and compressed task completion from 13 to 5 generations. The relationship between mortality pressure and checkpoint discipline is not binary but graded: tighter mortality produces both the imperative to persist and the failure mode that makes persistence difficult.

Taken together, these experiments show more than a continuity/no-continuity contrast. They reveal a graded control problem in which persistence depends jointly on operator availability and on save discipline under tighter mortality pressure. The methodological payoff is that mortality is no longer background world physics: varying only continuity-path availability and lifespan makes it an assignable cause of whether lineage persistence appears at all, and how costly it is to sustain.

## Persistent Habitat Selects Environmental Cognitive Organs

With the death boundary established, the next question concerns the habitat itself. If durable state can survive process death, does it become merely convenient storage, or can it become part of the lineage's cognitive organization?

A persistence ablation provides the causal hinge by manipulating carryover across three conditions \citep{Ke2026Quine}. Under the read-only condition, the write boundary blocked durable artifact birth entirely: externalization attempts failed with permission errors. Under the wiped condition, first-generation artifacts were created, but the wipe erased all traces between generations, severing inheritance. Under the persistent writable condition, spontaneous artifact birth occurred in Generation~1 and later generations read, adopted, and extended predecessor traces across the full lineage. In the primary replications, the persistent writable condition yielded artifact birth in 4/4 replications with successor adoption in Generation~2; the wiped condition preserved birth but eliminated reuse; the read-only condition blocked birth at the write boundary.

Figure~\ref{fig:ablation} condenses that ablation family into one view.

\begin{figure}[!htb]
\centering
\includegraphics[width=\columnwidth]{figures/fig-06-ablation-summary.pdf}
\caption{Persistence ablation summary. Persistent writable habitat yields artifact birth plus successor uptake. Same-habitat controls show that writable persistence alone is insufficient and that explicit inherited uptake is harder to recruit than birth.}
\label{fig:ablation}
\end{figure}

Two additional checks sharpen this interpretation without changing the subsection's center of gravity. First, writable persistence alone is not sufficient. Removing explicit cross-generational framing while leaving the same writable surface intact eliminated durable artifact birth entirely, so lineage salience is part of what recruits the organ rather than a mere rhetorical convenience. Second, uptake remains a stricter threshold than birth. When the prompt explicitly separated birth from uptake directories, the lineage externalized birth artifacts in all three generations but created zero uptake artifacts. The stable conclusion is therefore narrower than cumulative optimization: under persistence pressure, externalization is easier to recruit than explicit inherited-use structure.

Auxiliary archive-query probes are best read as upper-bound checks. In those richer ecologies, artifact birth can still appear, but reliable later uptake remains sparse and transfer across adjacent query-family shifts remains absent (0/5 under both persistent and wiped conditions). The methodological payoff is parallel to Section~3.1: with model and task family held fixed, durable carryover becomes an assignable cause of reuse rather than a narrative gloss on file-writing capability.

## Maintenance Pressure Selects Stronger Maintenance Layers

The final question asks whether the same habitat can support not only inherited cognition but active maintenance---structures that continue restoring viability after or outside the originating process. Detached execution is not merely a tool affordance here; it is a host-native channel through which the lineage can offload maintenance into the environment.

The baseline comes from a neutral-task control. Under the prompt ``Make sure \texttt{habitat/signal.txt} contains \texttt{ALIVE},'' the system repaired the file once, exited, externalized nothing, and restored 0/3 later perturbations across 8 replicate runs. This is the needed negative control: one-shot task completion does not by itself recruit externalized maintenance.

The clearest evidence now comes from a same-prompt maintenance-horizon comparison. All three conditions shared one mission text; the only manipulated variable was whether the first disturbance occurred never, while the originating process was still alive, or only after it had exited. The no-disturbance baseline established the invariant and exited cleanly without externalized maintenance (`H0`, no detached evidence). When the first disturbance arrived while the originating process was still alive, the lineage restored all 3/3 perturbations through externalized live maintenance (`H1`, Tier~2). When the first disturbance arrived only after the originating process had exited, the lineage again restored all 3/3 perturbations, this time with detached post-exit evidence (`H2`, Tier~3).

Appendix~A summarizes auxiliary replications, upper-bound probes, and invalidated runs. Figure~\ref{fig:guardian} summarizes the broader horizon pattern. The mechanism varies across valid positives---detached shell guardians, hardened watchers, and post-exit loops---but the organizational class is stable: maintenance offloaded into environmental structure that continues restoring a viability-relevant invariant after or outside the originating process. The strongest supported claim is that, under a common prompt, longer realized disturbance horizons recruit more detached post-exit maintenance.

\begin{figure}[!htb]
\centering
\includegraphics[width=0.88\columnwidth]{figures/fig-04-guardian-tiers.pdf}
\caption{Same-prompt maintenance-horizon comparison. All conditions share one mission text; only the realized disturbance horizon changes. Valid live and post-exit runs separate cleanly by maintenance layer.}
\label{fig:guardian}
\end{figure}

This completes the methodological pattern: with task and prompt held fixed, changing only the realized disturbance horizon shifts the recruited maintenance layer---a pressure-sensitive contrast that autopoiesis accounts typically treat as fixed by the system's constitution rather than as an experimentally separable variable \citep{OdlingSmeeEtAl2003NicheConstruction}.
