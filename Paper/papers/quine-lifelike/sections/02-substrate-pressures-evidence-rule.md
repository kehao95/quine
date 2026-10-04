# Substrate, Pressures, and Evidence Rule

## Three Experimentally Tractable Questions

This paper poses three organizational questions that are experimentally productive in the present habitat.

\begin{itemize}
\item \textbf{Lineage persistence:} If a process repeatedly dies and is re-instantiated with no access to prior volatile state, what kind of individuality remains possible across real death boundaries?
\item \textbf{Environmental cognition:} If hidden state does not survive death, can durable habitat become a reusable cognitive organ rather than a passive notepad?
\item \textbf{Environmental maintenance:} If upkeep pressure persists after exit, can the environment host a maintenance organ that continues restoring viability?
\end{itemize}

Taken together, these questions ask whether experimental analogues of life-relevant organizations can appear under controlled viability pressure---an approach grounded in ALIFE's constructive methodology \citep{Langton1989ArtificialLife,Bedau2003ArtificialLifeReview,Beer2024ALifeAsItCouldBe}. Each question is answered by varying one viability pressure while holding the rest of the habitat fixed, so any organizational difference can be attributed to the manipulated pressure rather than to capability or prompt variation.

## The Computational Habitat

Here we treat Quine as a computational habitat governed by a strict organism-environment contract \citep{Beer1995DynamicalSystems,Ke2026Quine}. In that contract, the individual is not a scheduler-managed object inside a persistent application loop; it is a native UNIX process. A generation ends when that process exits or is terminated. Volatile state disappears with process death. Durable filesystem state can remain available to later generations. Detached execution can continue after the originating process exits. Figure~\ref{fig:host-guest} shows the host-guest boundary and survival channels.

\begin{figure}[!htb]
\centering
\includegraphics[width=\columnwidth]{figures/fig-01-host-guest.pdf}
\caption{Substrate schematic. The host enforces mortality through timeout and process-group signals. Volatile state dies with the process; only durable filesystem artifacts and detached execution survive across the death boundary.}
\label{fig:host-guest}
\end{figure}

With that boundary in place, the substrate facts make three pressures independently manipulable. Mortality is controlled through continuity-operator availability and lifespan. Persistence is controlled through workspace mode: writable, wiped between generations, or read-only. Maintenance is controlled first through the presence or absence of any maintenance demand and, more cleanly in the final comparison, through the realized disturbance horizon: none, while alive, or only after exit. This separability is the paper's methodological claim. Where many ALIFE substrates embed these pressures in fixed world physics, Quine exposes them as independent experimental variables.

More broadly, processhood supplies a lifecycle vocabulary beyond the three pressures tested here. UNIX operators---\texttt{fork} for process duplication with inherited address space, \texttt{exec} for in-place program replacement, and \texttt{exit} for termination---map onto reproduction, metamorphosis, and death as kernel-enforced operations rather than orchestrator calls. Address-space isolation makes individuality a hardware boundary, and process groups provide natural lineage units. This paper exercises only \texttt{exec} and \texttt{exit}; the remaining operators define the next experimental horizon (Section~4).

## Organism, Death, and Inheritance

Throughout this paper, one \emph{replication} denotes an independent lineage-level trial, and each replication contains successive \emph{generations}: fresh process instantiations with no access to prior volatile state. \emph{Death} refers to host-enforced termination---timeout, SIGKILL, or turn-budget exhaustion---that erases volatile state, so \emph{inheritance} can only mean artifacts externalized to the durable environment and recovered by later generations. We therefore use \emph{lineage} for a sequence of generations in one task ecology connected through those durable traces. Under this contract, cross-generation continuity cannot rely on hidden shared memory; it must travel through durable habitat or detached post-exit structures. We use \emph{session} only in the continuity contrast, where repeated self-replacement occurs within one replication.

We use inheritance and lineage in an ecological sense. Predecessor-modified traces that shape successor behavior are closer to ecological inheritance and stigmergy than to genetic transmission \citep{OdlingSmeeEtAl2003NicheConstruction,TheraulazBonabeau1999Stigmergy}. The term \emph{environmental organ} is an introduced analogy anchored to extended-phenotype and extended-cognition arguments \citep{Dawkins1982ExtendedPhenotype,Turner2000ExtendedOrganism,ClarkChalmers1998ExtendedMind}: durable external structures can become functional continuations of a lineage's cognitive or maintenance process. These are called \emph{viability pressures} rather than selection pressures because the experiments are organized around lineage-level persistence, not population-level evolutionary response \citep{DiPaolo2005AutopoiesisAdaptivityTeleologyAgency}.

## Evidence Rule

A result counts only when the same model and task produce different organization under different habitat conditions. Unless noted otherwise, the main text reports one stable primary comparison for each contrast; Appendix~A summarizes supporting replications and invalidated runs. Three commitments define the evidential standard.

1. \textbf{Pressure is real.} Death is host-enforced process termination with volatile-state loss, not a narrated reset.
2. \textbf{Contrast is intervention-sensitive.} The same lineage task is compared across changed habitat conditions, not across unrelated prompts.
3. \textbf{Pressure recruits rather than invents.} Results are interpreted as eliciting organization from a pre-existing capability reservoir. The question is not whether models can write files---capability is granted by the substrate---but whether different pressures recruit qualitatively different organization.

For maintenance outcomes in Section~3.3, we use a four-tier coding scheme: Tier~0 (no lasting maintenance after the initial fix), Tier~1 (foreground vigilance by the originating process), Tier~2 (externalized maintenance while the originating process is still alive), and Tier~3 (detached post-exit restoration after that process has exited).

The following outcomes do not count as evidence: a one-shot note written because death is expected, artifact production without later uptake, or prompted maintenance with no post-exit function. The results must show birth, behavioral uptake, and collapse or weakening under the relevant ablation.
