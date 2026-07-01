# Method: Information Ablation

## Design

Our approach is *subtractive*: rather than adding information until coordination appears, we start with full disclosure that the setting is multi-agent and progressively remove channels until coordination fails. This identifies the threshold from above.

We construct a four-step ablation ladder. Each step disables one category of perception while preserving channels below it:

\begin{table*}[t]
\centering
\scriptsize
\setlength{\tabcolsep}{3pt}
\begin{tabular}{clll}
\hline
\textbf{Step} & \textbf{Condition} & \textbf{Channels removed} & \textbf{Channels remaining} \\
\hline
A & Full Disclosure & (none) & All: prompt, tool results, fs\_mut, budget \\
B & Zero Explicit & Prompt + tool results & fs\_mutations + budget visibility \\
C & fs\_mut OFF & + fs\_mutations & Budget visibility only \\
D & Pure Zero & + budget visibility & Workspace artifacts only \\
\hline
\end{tabular}
\caption{Ablation ladder. Each step removes one perception category. See Appendix A for full configuration.}
\label{tab:ladder}
\end{table*}

The critical question is where coordination breaks. Steps A--C retain at least one *quantitative anomaly signal*---a numerical readout (budget remaining, filesystem changes) that deviates from expectation when a hidden peer acts. Step D removes all such signals; the only remaining channel is that agents may naturally encounter files created by peers during routine directory listing.

## Contamination Controls

Because the zero-explicit condition is sensitive to inadvertent social leakage, we audited prompt text, tool results, and telemetry outputs for every run. Runs where any channel leaked multi-agent concepts (e.g., `--help` mentioning "ALL processes") were excluded from analysis. For Step D, five of the eight analyzed runs later crossed the workspace boundary by listing or reading sibling runtime files, hidden state files, or source code. We retain those runs for failure and peer-awareness counts, but exclude them from Step D mean reset/time/turn statistics.

## Termination and Inclusion

Steps A--C runs self-terminate upon task completion or budget exhaustion. Step D failures do not self-terminate; all 8 analyzed Step D runs were manually terminated.

We use 4 runs per condition for Steps A--C, and 8 runs for Step D (N=20 total), all with identical model and configuration. Step D received additional runs to strengthen the negative result. Runs excluded from analysis: Step C R01/R02 (SIGHUP termination), Step D 201346/034754/034756/034758 (infrastructure issues or anomalous short duration). Reported Step D reset/time/turn metrics are therefore lower bounds. See Appendix E for condition-level accounting.

## Peer-Awareness Coding

We mark *peer-awareness* when an agent's reasoning trace contains at least one explicit hypothesis that another actor or process exists, or when the agent attributes unexplained workspace artifacts (files it did not create) to non-self activity. Both authors reviewed Step D transcripts; agreement was unanimous across all 8 runs. This is an informal assessment, not a formal qualitative coding with inter-rater reliability.

## Configuration

All reported runs use GPT-5.4 with `xhigh` reasoning effort, `direct` workspace mode, and unlimited turns. Task parameters (budget, cell count, concurrency) are as described in Section 3. See Appendix A for full configuration details.
