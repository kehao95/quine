\section*{Appendix A: Primary Comparisons and Experiment Notes}

The main text reports one stable GPT-5.4 comparison for each contrast. This appendix records only the replication and invalidation context needed to interpret that choice without interrupting the argument.

\textbf{A.1 Main-Text Comparisons}

\begin{itemize}
\item \textbf{Section 3.1 (continuity):} the minimal continuity control is the primary exec-enabled / exec-disabled pair reported in the current runtime, and the longer-horizon mortality contrast is the 45-second versus 90-second lineage run.
\item \textbf{Section 3.2 (environmental cognition):} the main ablation is the persistent / wiped / read-only carryover comparison in the primary GPT-5.4 evidence set, supported by same-habitat prompt controls that separate lineage salience from writable persistence and birth from explicit uptake. Richer archive-query probes are retained only as upper-bound notes on how far uptake and transfer currently generalize.
\item \textbf{Section 3.3 (maintenance):} the main comparison is the same-prompt maintenance-horizon family \texttt{H0/H1/H2} in the primary GPT-5.4 evidence set, in which only disturbance timing changes.
\end{itemize}

\textbf{A.2 GPT-5.4 Follow-up}

\begin{itemize}
\item The continuity contrast is shown in one primary control pair because its evidential role is narrow: to establish that real death recruits continuity only when a path exists.
\item Fresh GPT-5.4 environmental-cognition follow-up now closes the main comparison cleanly without auxiliary model support: persistent \texttt{A} was reproduced in three short positive slices, while same-habitat \texttt{neutral} again yielded no durable artifact birth or reuse and \texttt{birth-uptake} again yielded zero uptake files.
\item Fresh GPT-5.4 maintenance follow-up replicated the primary \texttt{H2} post-exit positive. Later \texttt{H1} reruns sometimes invalidated as \texttt{horizon\_miss} because the seed process exited before the live tamper window opened; those are treated as harness invalids rather than scientific negatives.
\item The richer archive-query cognition lanes remain more variable than the carryover ablation. They still show that artifact birth is easier than robust reuse, but they do not provide a cleaner main-text anchor than the simpler persistence controls.
\end{itemize}

\textbf{A.3 Notes on Invalidated or Secondary Runs}

\begin{itemize}
\item Some no-disturbance maintenance baselines still externalized a guardian without realized tamper, so the no-disturbance case should not be read as a universal zero-externalization condition.
\item Hidden-future, open-future, and surprise-transfer cognition probes are included to bound the result, not to expand the main claim: they show that artifact birth is easier than reliable uptake, and uptake is easier than transfer across shifted query families.
\end{itemize}
