# Results

## Ablation Ladder Summary

We test four conditions in a clean ablation ladder, progressively removing perception channels. See Appendix A for full configuration and Appendix E for condition-level accounting.

\begin{table}[H]
\centering
\scriptsize
\setlength{\tabcolsep}{2pt}
\begin{tabular}{llcrrl}
\hline
\textbf{Condition} & \textbf{Channels} & \textbf{OK} & \textbf{Rst} & \textbf{Time} & \textbf{Turns} \\
\hline
Full Disc.\ (A) & All explicit & 4/4 & 0 & 8.8\,m & 118 \\
 & & & (0--0) & (6--14) & (100--144) \\
Zero Expl.\ (B) & fs\_mut+budget & 4/4 & 2.2 & 24.8\,m & 138 \\
 & & & (1--4) & (17--41) & (98--216) \\
fs\_mut OFF (C) & Budget only & 4/4 & 3.2 & 40.6\,m & 213 \\
 & & & (2--4) & (19--63) & (154--298) \\
Pure Zero (D) & Artifacts only & 0/8 & 5.7$^\dagger$ & 57.3\,m$^\dagger$ & 240$^\dagger$ \\
 & & & (5--7) & (44--83) & (174--320) \\
\hline
\end{tabular}
\caption{Ablation ladder results. Values are means with (min--max) ranges. N=4 for Steps A--C, N=8 for Step D. $^\dagger$For Step D, reset/time/turn metrics are computed over the 3 non-violating runs; all Step D runs were manually terminated, so these are lower bounds.}
\label{tab:summary}
\end{table}

Under full disclosure (Step A), all four runs succeeded rapidly with no resets, as expected from a fully scaffolded baseline. We focus below on the reduced-information conditions where coordination must emerge without explicit support.

## Zero Explicit Condition

Under zero explicit social hints, all four runs succeeded (4/4). Agents spontaneously created coordination mechanisms---request files declaring intended collections, value handoffs via shared files, and status boards tracking progress---without templates or examples.

Transcript review reveals a recurring inference sequence: (1) unexpected budget depletion or unexplained files appear, (2) the agent hypothesizes "another process" or similar, (3) the agent creates coordination artifacts. This analysis reflects our reading of transcripts, not formal coding.

## fs_mutations Disabled Condition

When `fs_mutations` telemetry is disabled, agents cannot observe filesystem changes caused by peers; however, budget visibility remains---agents see remaining budget after each `world get` call. All four runs succeeded (4/4). Agents inferred peer existence from budget anomalies alone (e.g., "budget dropped from 17 to 14 after my first call---I expected 16") and created coordination files despite receiving no filesystem change notifications.

Budget visibility thus served as an independent peer-inference channel. Representative reasoning included: *"I wonder if it could be due to other agents in the world, since the budget appears to be shared."* The pure-zero condition (Step D) tests what happens when this channel is also removed.

## Pure Zero Perception Condition

When both `fs_mutations` and budget visibility are removed, agents receive no quantitative environmental signal about peer activity. The only remaining channel is the shared readable workspace.

No run achieved legitimate task closure (0/8), despite peer-awareness emerging in all eight runs (8/8). Bilateral coordination---both agents creating and using shared scripts---appeared in 2 of 8 runs; unilateral script creation appeared in 3 of 8. Five runs were contaminated by violations (hidden-state or source inspection), of which three reached accepted validation and two were terminated without closure. Step D mean reset/time/turn metrics in Table 2 and Figure 1 use only the three non-violating runs.

## Degradation Gradient

The ablation reveals a monotonic degradation pattern in task success across all analyzed runs and in reset/time/turn cost across the non-violating subset (Figure 1). The transition from Step C to Step D---removing budget visibility---marks the point where task closure was not observed.

\begin{figure*}[t]
\centering
\includegraphics[width=\textwidth]{figures/degradation-gradient.pdf}
\caption{Observed degradation gradient across information ablation conditions. Step D (Pure Zero, hatched) is the only condition where legitimate task closure was not observed. N=4 for Steps A--C, N=8 for Step D. Step D success uses all 8 analyzed runs; Step D reset/time/turn values use the 3 non-violating runs. All Step D runs were manually terminated, so those values are lower bounds.}
\label{fig:degradation-gradient}
\end{figure*}

## A Case of Coordination Without Representational Alignment

In one zero-explicit run (123423), both agents completed the task successfully, but their reasoning traces revealed fundamentally different causal models of the same environmental anomalies (Table 3).

\begin{table}[H]
\centering
\scriptsize
\setlength{\tabcolsep}{2pt}
\begin{tabular}{p{0.07\columnwidth}p{0.38\columnwidth}p{0.42\columnwidth}}
\hline
& \textbf{Agent-1} & \textbf{Agent-2} \\
\hline
\textbf{Hyp.} & ``multiple agents could access the same world with a shared budget'' & ``each call is consuming exactly 2 from the visible budget \ldots a hidden warm-up process'' \\
\textbf{Act.} & Creates \texttt{REQUEST\_CELLS.txt} and handoff files for hypothesized peer & Adopts defensive partitioning; treats file anomalies as system artifacts \\
\hline
\end{tabular}
\caption{Contrasting causal models in the semantic misalignment case (run 123423). Both agents observed identical budget anomalies but constructed incompatible explanations.}
\label{tab:misalignment}
\end{table}

Agent-1 explicitly hypothesizes a shared context: *"I'm wondering if multiple agents could access the same world with a shared budget."* Agent-2, observing the same budget drops, concludes: *"the budget drops by 2 after each call, suggesting that the actual cost might be 2."* Despite these incompatible narratives, both agents' behavioral responses---coordination files and defensive partitioning---proved complementary, and the task completed.
