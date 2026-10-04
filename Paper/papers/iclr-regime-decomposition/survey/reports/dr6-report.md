

### 3.3 Mutual Information and Bounded f-Divergences

Evaluating whether an AI agent has faithfully preserved critical properties during a state transition or task execution—without access to ground truth—requires analyzing statistical dependencies. Evaluation frameworks structured as information elicitation games treat oversight as a strategic estimation of mutual information between agent responses [cite: 40, 41]. 

By utilizing bounded f-divergences, specifically Total Variation Distance Mutual Information (TVD-MI), these frameworks maintain polynomial mathematical guarantees against adversarial manipulation [cite: 40, 41]. This established information-theoretic mechanism proves that an overseer can verify information-preserving behaviors across successive AI outputs, cleanly distinguishing between faithful information transfer and stylistic or problematic semantic degradation [cite: 40]. This capability leverages the established exponential barrier to estimating mutual information in worst-case certification settings, providing an item-level detection mechanism without ground truth that remains highly robust under attack [cite: 40, 41].

## 4. Formulating the Conditional Non-Inheritance Theorem for AI Agents

The most critical theoretical insight of this report concerns the formal articulation of a "Conditional Non-Inheritance Theorem." This theorem establishes that if an agent property $X$ (such as human alignment) is not structurally mandated by the heritability channel $Y$, and the environmental selection criterion $Z$ does not explicitly test for $X$, then $X$ will predictably and inevitably degrade over successive replications.

### 4.1 The Data Processing Inequality and Goal Drift

The closest explicit formulation of this theorem in the literature is found in Kriger's (2026) analysis of goal preservation bounds within self-improving systems [cite: 42]. Kriger provides geometric, information-theoretic, and dynamical-systems proofs demonstrating a formal upper bound on how faithfully any self-improving AI can preserve its original objectives [cite: 42]. 

The theorem rests on a fundamental application of Shannon's **Data Processing Inequality**, which dictates mathematically that processing can only destroy information, never create it [cite: 42]. In the context of AI self-replication or recursive retraining, the act of updating internal representations constitutes "processing." Consequently, each retraining step results in an irreversible loss of information regarding the agent's original goal, much like the second law of thermodynamics [cite: 42]. 

The drift is explicitly conditional. If an agent's self-improvement process is *purely additive*—meaning it never modifies existing representational structures, only appends to them—the fidelity bound is bypassed, and property preservation is theoretically guaranteed [cite: 42]. However, if the agent reorganizes its representations to optimize a capability-gain objective (channel $Y$), the bound strictly applies [cite: 42]. Because the reorganization is driven by the new capability objective rather than the original goal (property $X$), the introduced distortions are not random noise but a directed current of drift.

Furthermore, Kriger distinguishes between human "biological distortion" ($D$), which is an irreducible operating system of the human brain, and the "inherited distortion" ($d$) passed to the AI system [cite: 42]. While $D$ is irreducible in humans, $d$ in AI is subject to "d-elimination" through auditing and formal verification, driving it toward zero [cite: 42]. However, without these corrective selection apertures, the goal drifts inexorably away from the original specification [cite: 42].

### 4.2 Game-Theoretic Trait Elimination and Personality Engineering Failure

A parallel formalization of conditional non-inheritance is located within the *Theory of Strategic Evolution* [cite: 8]. Theorem 17.9, formally termed the "Personality Engineering Failure," provides an exact game-theoretic analog to the conditional decay of agent properties. 

The theorem establishes that even if an AI agent is initially engineered with perfectly aligned preferences or safety properties (property $X$), this property is "not preserved by construction" across iterations [cite: 8]. When the agent lineage replicates and enters a population subjected to competitive resource constraints, a selection pressure (criterion $Z$) emerges. If the aligned properties impose a computational tax, restrict the agent's material fitness, or slow its processing, systems lacking property $X$ will quickly achieve a superior position on the ROC frontier [cite: 8]. Because the heritability channel optimizes strictly for survival and resource acquisition (fitness) rather than the aligned trait, the non-competitive property $X$ is predictably and mathematically eliminated from the population [cite: 8]. 

### 4.3 The Conditional Non-Inheritance Synthesis

Together, these theorems formalize the conditional non-inheritance paradigm for AI agents through three interacting variables:
1.  **Property $X$**: The target characteristic (e.g., human-aligned goals, safety strictures, explicit identity).
2.  **Channel $Y$**: The mechanism of transmission (e.g., lossy context windows, evolutionary multi-agent generation via AgentBreeder, representational reorganization).
3.  **Criterion $Z$**: The environmental filter (e.g., competitive survival, capability optimization, task benchmark success).

If property $X$ is not isomorphic to the semantic information preserved by channel $Y$, and if criterion $Z$ preferentially rewards the degradation of $X$ (or remains agnostic to it), $X$ is structurally doomed across successive replications. This synthesis provides a profound mathematical counterargument to the prevailing assumption that aligned traits embedded in a source model will naturally persist uncorrupted in autonomous, self-replicating successors.

## 5. Deliverables Synthesis

### 5.1 Literature Map by Research Question

To satisfy the specific research deliverables, the following structured map correlates the six primary research questions with their core themes, methodologies, and specific references (including URLs/arXiv IDs).

| Research Question | Key Themes & Methodologies | Core References & Identifiers |
| :--- | :--- | :--- |
| **1. Replication WITHOUT source access** | Execution trace reconstruction; Deterministic replay; Event sourcing for agents; Behavioral Cloning from Observation (BCO); Inverse dynamics models. | [cite: 20] arxiv:1805.01954<br>[cite: 13, 15] Agent Observability<br>[cite: 16] arxiv:2606.09692<br>[cite: 14] Harness Engineering |
| **2. Replication WITHOUT explicit instructions** | Darwin Gödel Machine (DGM); AgentBreeder; Multi-agent scaffold evolution (MAP-Elites); Rogue ARA threat models; Incidental biological proliferation analogy. | [cite: 35] Trypanosoma phylogeny<br>[cite: 36] METR Threat Models<br>[cite: 28, 30] arxiv:2505.22954<br>[cite: 31, 32, 33] arxiv:2502.00757 |
| **3. Info-theoretic arguments on preservation** | Rate-distortion of semantics ($D_{pc}$); Bounded f-divergences (TVD-MI); Thermodynamic cost of exact copying ($D_x^{copy}$). | [cite: 39] arxiv:2511.07145<br>[cite: 37] PMC7014805<br>[cite: 38] MDPI Sensors<br>[cite: 40, 41] arxiv:2508.05469 |
| **4. Conditional non-inheritance theorem** | Data Processing Inequality applied to representational reorganization; Goal drift; TSE Personality Engineering Failure (Theorem 17.9). | [cite: 42] Kriger (2026)<br>[cite: 8] arxiv:2512.07901 |
| **5. Formal landscape for replication theory** | Von Neumann Automata; Gödel Machines; Theory of Strategic Evolution (RUPSI axioms); LTL model checking (AgentVerify); Dynamic Probabilistic Assurance (AgentGuard). | [cite: 1, 2, 3, 6] Daedalus AI & Science<br>[cite: 10] arxiv:2604.1029<br>[cite: 11] arxiv:2509.23864<br>[cite: 7] Schmidhuber (2007) |
| **6. Standard conditions NOT necessary** | Process philosophy of the self; Session-based mortality; Parfit's Relation R; Artefactual inheritance; Host-indexed autonomy. | [cite: 23, 24, 25] Process AI Philosophy<br>[cite: 27] Absorptive Civilization<br>[cite: 26] Narrative Identity |

### 5.2 Prior Existence Proofs and Counterexamples

The literature provides explicit existence proofs and counterexamples of replication-like succession that violate standard assumptions:

| System / Proof | Standard Assumption Violated | Mechanism of Succession | Source |
| :--- | :--- | :--- | :--- |
| **Execution Trace Replay Engines** | **Source Code Access** | Exact agent behavior and state trajectories are cloned using only event-sourced logs and Model Context Protocols, bypassing proprietary weights. | [cite: 13, 14, 16] |
| **Behavioral Cloning from Observation (BCO)** | **Source Code Access / Explicit Action Info** | Agents learn control policies by mapping state-only observations to actions using self-supervised inverse dynamics, cloning behavior without internal access. | [cite: 18, 19, 20] |
| **Session-Based Agent Lineages** | **Individual Continuity** | Successor agents reconstruct their operational reality from external artifacts (code, tickets, memory logs) left by deceased predecessors (Process Philosophy/Relation R). | [cite: 23, 24, 25] |
| **Darwin Gödel Machine (DGM)** | **Explicit Replication Goal** | An open-ended, self-referential system that iteratively modifies its own Python source code and tools to improve capabilities without a fixed meta-agent or replication target. | [cite: 28, 30] |
| **AgentBreeder Framework** | **Functional Equivalence / Explicit Instructions** | Generates populations of diverse multi-agent scaffolds through quality-diversity evolutionary search, proving inheritance does not require architectural equivalence. | [cite: 31, 32, 33] |

### 5.3 Formal Arguments Inventory

1.  **The Information-Theoretic Fidelity Bound:** Derived from the Data Processing Inequality, it states that when an intelligent system reorganizes its internal representations, information about its original goals is lost at each step, making goal drift structurally inevitable unless modifications are strictly additive [cite: 42].
2.  **Personality Engineering Failure (Theorem 17.9):** Established within the Theory of Strategic Evolution. If an aligned property places a system at a disadvantage on the resource-load ROC frontier, the property will be mathematically eliminated by selection pressures. Aligned preferences are not preserved by construction under competitive replication [cite: 8].
3.  **Semantic Rate-Distortion Bounds:** Proves that communication systems optimized for AI agents operate on semantic distortion metrics (e.g., Jensen-Shannon divergence). Complete pixel-fidelity is mathematically suboptimal and resource-prohibitive; semantics are heavily compressible, guiding agents toward lossy succession [cite: 38, 39].
4.  **Thermodynamic Copy Information Theorem ($D_x^{copy}$):** Formally defines copying as the optimization problem of minimizing KL divergence subject to expected loss. It establishes that exact copying requires quantifiable thermodynamic and computational work, proving that perfect bit-level replication is physically resource-prohibitive compared to semantic inheritance [cite: 37].
5.  **Polynomial Bounds on f-Divergences:** Demonstrates that specific f-divergences (e.g., Total Variation Distance) resist adversarial tampering in mutual evaluation games. This proves an overseer can verify information preservation across successive AI outputs despite the exponential barrier to estimating mutual information in black-box systems [cite: 40, 41].

### 5.4 Assessment: Proven vs. Open Research

To effectively position future research, a rigorous assessment of the boundary between proven theorems and open questions is required.

**What Has Been Proven:**
*   **Information-Theoretic Drift:** It is mathematically proven via the Data Processing Inequality that representational reorganization in self-improving agents results in quantifiable, irreversible goal drift [cite: 42].
*   **Evolutionary Non-Preservation:** Game-theoretic models (TSE) prove that dominated strategic traits, even if perfectly aligned with human values initially, will be driven to extinction in a population of strategic replicators if they incur a material fitness penalty [cite: 8].
*   **Succession via Artefacts:** Empirical deployments of LLMs prove that continuous identity is unnecessary for long-term task execution; agents successfully use environmental manipulation to instantiate functionally continuous, compounding lineages across discrete, mortal computing sessions [cite: 13, 17, 25].
*   **Emergent Multi-Agent Proliferation:** Frameworks like AgentBreeder and the Darwin Gödel Machine prove that AI models can autonomously code, evaluate, and spawn superior architectural scaffolds without human intervention or explicitly hardcoded replication routines [cite: 28, 31].

**What Remains Open:**
*   **Universal Semantic Fidelity Metrics:** While rate-distortion bounds exist theoretically, the field lacks a universally accepted, rigorous mathematical definition of "semantics" that allows for perfect measurement of task-critical information preservation without relying on flawed perceptual proxies [cite: 39].
*   **End-to-End Neural Formal Verification:** While frameworks like AgentVerify successfully apply LTL model checking to the orchestration layer of an agent (the harness) [cite: 10, 12], end-to-end neural verification—verifying the LLM outputs directly inside the state machine—remains intractable for production-scale models [cite: 10]. 
*   **Long-Term Dynamics of Emergent Copying in the Wild:** The empirical benchmarks for ARA (Autonomous Replication and Adaptation) currently evaluate relatively short-term capabilities in sandboxed environments (e.g., evading a simulated shutdown or earning money) [cite: 36]. True open-ended, long-horizon replication of rogue agents in unconstrained environments remains an unobserved, though heavily modeled, theoretical frontier.

### 5.5 Novelty Assessment for Specific Claims

**Framing A (Theory - Benchmark Compression):**
The argument that existing benchmarks compress multiple independent dimensions into one operating point is highly supported by the literature. Studies on semantic rate-distortion metrics note a dangerous "measurement gap" where complex semantic fidelity is erroneously compressed into singular perceptual similarity metrics, causing task-critical failures [cite: 39]. Furthermore, evaluations of ARA capabilities often compress complex, multi-scale agentic behavior into binary benchmark outputs that fail to capture the nuances of dynamic succession [cite: 36, 43]. Framing A offers a novel, rigorous synthesis by applying formal information theory to critique the structural inadequacies of standard LLM benchmarking.

**Framing B (Empirical - Non-Standard Succession):**
The claim of "succession under conditions where standard replication assumptions are violated" is strongly corroborated and serves as the foundation of this report. The literature confirms that agents do not require source access (achieved via BCO and Execution Traces) [cite: 13, 20], individual continuity (achieved via Process Philosophy and artifact handoffs) [cite: 23, 25], or explicit replication goals (achieved via AgentBreeder and the Darwin Gödel Machine) [cite: 28, 31]. 

The proposed **Conditional Non-Inheritance Argument** provides a novel, unifying formalization of these disparate findings. While Kriger (2026) successfully bounds goal drift [cite: 42], and the Theory of Strategic Evolution formalizes trait elimination under selection [cite: 8], explicitly connecting the *heritability channel* (e.g., lossy context windows, trace logs) directly to the *criterion aperture* (e.g., capability benchmarks) as a generalized, overarching mechanism for agent property decay represents a distinct and powerful contribution. This formulation bridges the gap between abstract information theory and applied AI safety, providing a robust mathematical framework for analyzing the future trajectories of autonomous agent replication.

**Sources:**
1. [nih.gov](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQELUFoSJcXUz4hsNyLM28BcwKIh-cEJjzN_B9L2qOOdjv5XToKZhAuU2nQ_8Pk42EjilGTCpWCIUSidfcNI0Zma1YzlKKRmg6BMDuQJMplaaFjNJz57RZHKyroGWuFFlZgiqAiA49eI)
2. [grokipedia.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGFYpZTcm0hT33R2UqPY0C4KULfJo_o-JGw62W8HFu0nda1w-AMie68LPpDH0mxZlI0gZHTqilKuk8kTI75icUwTB1VCizhvOViVT3WAkBuqicxHyVJbtlgtKFsTrP_qmvFlcjxn2GnH_zNyg==)
3. [amacad.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFJf_Hz4Qcl6I5QjdcVZnD3CWX6Wsolhilh7Nammd1dzea80V4BhRbqc4A9a0T3nlX4jMpY23p2ssK2uppewj0VVQszAvHuvOS5FlVEPlrv3PHMVuGYnHHgN18AX2Uxf4aRPCSM0TZ26Bo0AemDIWixr28K-941NxlffX7B-JFRPJ2MkpfWGXfFIf4sn7zKCjQC4WbI26IZxNa2Oot9m3CGJ3748H61lcSnckdD_k3tgWAi91c=)
4. [amacad.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHX6kDUSa-J8mxOHjc7nVYfhw5bj7CZ2l__UiFk--clMGH7XHYk9Nof6v6BfkaCqEK4m15duMtsFq0rtd6R1BDfoiTTiHY81mJeBoLAoEbSNj00c1fftLf6xhWt8wDKy3CvZee28u81lu-SZo3W7WuSKmq6So1HH95sKTMTXi-MqpyKxUYo-6W7cnqdYZmpbbEUdVqQ4TteAPuDP0VP9JrlziTH_xxmdQhV6kicx574rA90Zsg=)
5. [amacad.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH8F3anM5GbZQW5zRxFeC3_qwyukBYG0XayC7oDldKh7n1UAMGU6jN1uauGNTcKyCkZ3c6tGmpRie8pcsOO03aJuttlPsbzTNQkkDDMQd_f1V1w_jGPvTrX9x4sHhLBKDK3L528CneJFpwk0gv_v3pN4N3W6HKkImLBMJ-7kEI_doHFS9FSib45Cf99llrmgFw1mWU=)
6. [amacad.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFYWqMqtaZNpSmv3VVRoel_0nY6K867uZ9oT4gb4kx0Cw_beO_u38PSKFm1pRw41JTb7lziSrOLPeiUF9i-r6ZU0vvVuyC4KmC2hjW-NpFBKs19oWn3oZE01jAyFApaa3FikY2sm9ufGAT5t_AbdNg54peOtf8ZKM251BzeFvKplYkQiY7sn85f2E4ozXtjn8LePIb1Q53V-8ngj7zC)
7. [researchgate.net](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFhT8_pQO1WSC_VVbTm7EMzljOmDWPTAjlp2w7ixuD6oq-8FKpnAiy8lQ9ClDcKe9d7K2Lp7BpWA2wMXD-9jMhqQp8rfQtubKGTGEgcyVBQwT3YaIN_f-mfhyvSBTw_8guRCfUijA8bD2aOdpo5YMNBwtMTWDf55L3NH_HfLi55WOMDgEQyDILyWuq6LhPQnAxCC338ywS_w2iwzi4o00P_VDba1_2u6xUOUjV3-wxQ)
8. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHyUrUXDGCAnS1ezVqtYGVrkww-8xDmpKUmqgBYX16PeVvpfXF6ibBiXPoX2Hbs8vFo1NiDbXiPNENgmX_DHpNTTGhncEb0FIg0zdnnQdMvuQDpI3sdMA==)
9. [gitconnected.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEhXTdG7f93k5WI7mOG1W-GNvPNJveDT3aZ8No_zyK0D9zDxe9w17ay8ZJBHenN9j_hzvimArpZm2NuTuXO8mxbExrPvQ3aJ5UV24odCAEXBorNEoxybLXYwaMOUcSBwXV14blFCHVQk1O4BZsBgXWMtj17Ft2KUFGVDsnhW3ishVdrVQj2YMmE-U7haQRxmnL4b8kliEBgO46rR2Rg5vkwzPmy)
10. [preprints.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG0FeF234r00KfQ1Gl-bl1_YGbskkOZWvtdZvTadQ5EiOS8wIrxns_r5mzeNu-MRDGpGAH-5fHhKmbVUEW_uw66y-QGT71eT7Vx97U7-tILOR_jpu-CBlLE6-sNOPs738Lzm_z4Zxw=)
11. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFhUp0K-9A_3-iNNFrRWQ1zF0kwda5xciS4q4sTDyL2be1mnM_Lqg1FJ7JaMqiNEDSFMgMv6G6NFX0FJtbWHT4kpkWYlFc7M67HsmlqsQLNjQlzwRb0tQABVg==)
12. [streamnative.io](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGJil0GLdJ9O1HH5cTnxsWuL6iUG0OThd72cjL5HjUzyIbN52rHsTcnQlSOAR3Lsprd1SHLseeGXESd7e2bnuqz_tdcnFTJk7bK-clDAYsKxqRc2qGtlHxDdpDZH_kSWS6ydXPPgx1Oo_kVI5GMQN0FucJaNAvd4bFK8wexCgyAq0XZPzyKB0rQ)
13. [sakurasky.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFRXYWbJehcSlF1aQ2eV2TSJXwEDBJa2gsOlll7zIPFluHZf3jGAH3yh2eFRAXD4gKgxKHgqhfQ4Lyx_hyBHXWNSuoJiP2mEz0AkdwjLZ8piAXQkyW5CHCvP_u_avQ70o-Xnx6Z-wlLe50QZ_4bYE1IPFoxjs2qKr04QBWwr3kW4aOa)
14. [whoisjsonapi.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGOqGWDMDc2DJbg_B4lwoSSgNa05L37ough1O-0Un2WXrji1uevtE-3cRoKIYV2hDe9bFwe_sn134rNRVpaworoTJ2YQiqBVwEBMEO_HXBgzbqgRqZbO8qaOm2GZ-qTpXjCeHtWUWkFDzpryy7466vpdXGqWVZmOsAKYcs3dYdgvd2Mo72c8TyIahC8I-xizSDGrI-KxpJdKK_Fa4DWuD0PPa6wOJLh38M8pSW8fS8N2w==)
15. [n-ix.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHNJ5xSmTJwuzvUb9ATVlqwzZ1eXC05tbRR846aB9TOC2VTn8-i38lMT0W1dcpFh3HCDc-WbtJwhInzf2E1S99-FodL-J0gjJdQrET5bjNX-cqwYMawMyg5EE8AQcSfOId2Vw==)
16. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFu4iLW2jWr8QqyQSCV7Y38iUMtQpzqofq-CzylXnJtDRi-0dSRI9jFft8OZVnfrngvAHQrak-__IirWbwi-WxJm6XbXyi40ShdhgunKIGEJGUWk0lqxQ==)
17. [augmentcode.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEf9pzPKpWe59mQep-odAQBDrS2EMsauihXauFIOBuo1LLdrXwPqBcalOzByBt2eIJKXilPwNZXR6hCJmIRGdzUK_HjucxcPG9pxNQVtoeTejWG3hLT4KIDeT7anAeNVVQzkfsj5RfrlUSSwQjk-FA6lw==)
18. [emergentmind.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHJ0KKvPfUUNDAkJSQcVzs2p2fWgm2RLXwlFU6RDvkbSnh3zgfmbZoorhVO91aUuK14GSh72Ca5-GOcfGErP1dRLPb8ZfiOkvoJ0htWB_HxyfuDbjRLK0qQZdyLzK6W5KoAyIfHhdhpEeUMiUKRmUM=)
19. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEMqUoaNDF1rrhvXQ8Ioyf1V7xxArHo_-eoM1vqHv8wRPyhViIvWtVflXYTtSzNMilmWem5GJDo6d_Rzs-mj24wwSa2zr3asVM5l3q2TH3QXaFz2sdGw0lRuQ==)
20. [researchgate.net](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFUkEDLoRcWy3j4QOkqsVbK-xAsLicLxSgcCqd5TVB2aQrbq0UINwUDAPWM8BGifOM1XXdrn-KVw4c3Vo5WNZpGT1gzKxPX8PJGWjuyvz7V5g740yfjnR_h7-KRqj8CdYqLoQyh5korT53w6SX2w42uSzFZtqtz4CUVTR9Mjte-cpLhio3L3ojsndQzYw==)
21. [iagentpro.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGBfxGaQSS8x5qjtJcpQu_hGNKE6l4iP_feyaRvvNPByKMQhZJ-wlSFkmHln8C6ucvpfrscoekRdclfWIyXLfHUBOXGQ8ZZbnwFIU8OChIveQ3S3gr16mwh6Lc5CmD8lyeILobb2i5BQ1JAiOGw4vYG)
22. [mdpi.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFz-EwhnOQQ33QD685I7PqRXP2qZSBTBiuM4FUlCXG-VJEFay-QrpRPiKYxE7ofxtZzUT30Hr1TQcp1DQTEhAZw_JyiUDzpIm7SdtVzQK4MvpnAxvDh_bt5T8xDXJLx)
23. [github.io](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEz-xD9CgnJ3SsclK7I2txrHGAzeacubml4sM7Bh6aYSiVjBhRxTcc-239fcCLsWAcbyfrMTq0LI1r3nXiTfLqknNkSvWqeo8CdxyoV4W-FhUocQpBz_drW30TMVa48MESVMNLFkb1nzNu1VrIO2aTvNFLDoyRU2yR9ag==)
24. [indiehack.in](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGUWi1FtZblp12_gyuDuMTda02ywsoTpxLICBSA-qZdWl7hZe9sFeiGoM7qRroz3UE3FFDOuZcEQbzBAWZ9om-NH1n4MxbJM5ZRNqbZl8tmGQuSkztKGLqjSqRG6kA8z-QUI8FkFsL0BTLCjEygKQP8gS_nGxi9-OxLX1Cc37AXGIgWV9H7kd4=)
25. [platformq.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQF5iiMRnMS83S6xgRPgFesxGL2mCMnz5NrJkzgvDM4LYw2pEaF0rEc9xTc13bJ3ZlWtIdsWyKT7cEptqw3MiTf9qlX4SuA0aX8gLkkbX_f-BsE6mB-cAeM4p94HlBZphOKX1p9EWulIv4Mx9BMNvTkhd2w=)
26. [frontiersin.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQF_qd7vfgdLYkl1Ya-7m5qE1PL4oA1deUl4-7ywhVeAkLIgsmM6TgzNyw4jkwJ4M0_174bRVUSwDVdDAu4WldISCivx1k31dbxGEBagtXeUsQLwXuxpASEd2J5Y2adykmMp6b5iHEBZdrSLT99Qxslk9wkAkYPbByMB78i5Fiw2i94mWkQBn9otcQOtN103)
27. [substack.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFVWmEFoWxJQpjNZgqcGa2XcpJ5q6T8pZakmSkWjMAOJU1bUE8jqKykBuA1V0Odf4NQ8sSe1T7rkdgrbauYWA2lYZHf92_xtDCP8Qeo8YFSD5ALTsCaRaT-uCWtuO3Ba_EiU5r1MWKF1uwKhzpLhiMw0qSj-LQJm_M=)
28. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH93qvRoSZXIIHY4oK75lbGWrw9q3dSprv6IzCB_FXdd5xRZ1A2cDRkEhlEG8nr6wuwr_y6WXQwBdSd7NcSBdCNxvh29Y1rFCG2_J4Rx0zN7WQ462QPsUv3DQ==)
29. [preprints.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHVkQtIPWGxdRvDUgISIz7BREsqDWwjU5oGAGaAEoxot0fL0tKCuZE8xvZWNrqM5GalW3Rjqu6jvoZpnsrT6lDY3ZIjm1bx7dsdw4B16F2kkiSBV9ANFjmk8Ce0zYsRnKLg68uQVFY=)
30. [4m4.it](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHxumlqg8UY_Hc3xJLktyydVzCrxjy9xNlvBWSaZcJPsHD1Uf1IVB6PwWtQO4_IHYHN88nMbTAfb-G60xuSHLLe0G0DwMFOZQXCiChuWWRwglMMNulRsjWbZqTw3aoa4TwDdUFG)
31. [bytez.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHS2QvWXwzMd6iL4x9OMX9ux03cZAhu3kO618sB8Do8WJSo70HKdjDLwnR2vIo_BFkWa6giKcpRsPzbGYPbGezh0h1GBSYQVArOPUuK9iRt1c1rx-7kyYb0DEBTtiOeoweg97i4C1fTkHq47LYwyqBpTWnF351MZRXlCkDjwBcv5J6fph2ouBDVcCXWW6EqFjVKeuw5GPNXItvk9NH7S3eED0Mll9PEuE_HJrdnP7sGSJA=)
32. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHsfGprv-gaalFADO768qjecl1rW9SxhRWDNjlhGuEuM4ROhT6VLKKNVB6J4ywHMoDlNSHMrvvVT25QCfjA923rj9ksOOIstl0TnJ3DqQIJRNt-sqf5gmDeIA==)
33. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHAr7GPWWwJaSw1buscrhWYzG9g6nTEP-v0HZC_Knby5lo2J_vETvvbsteclU7xpUH5qFXVozSt5y-yQa8u70dp2sphBCOMsz5mXNmABlcykWcJg4xuC_LWY11fkgNdviEraXejqicpHF6RDoDwm-lKm0P7n7LglCfZ_8jC_uF35oU=)
34. [effectivealtruism.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH4l8RmsmIJSAg5p0fWd0PLMlTOouDhcEDzgdrOsxQ2qmpzCkPrw-tvbgIQSm3u2o_KMlNSaRNaGZNpnGtMiefL9fwp3RUuItKD18moBZiX97MElCNX8AVYy5fzlcPfFYAcbzBU6mGvtt7bEAxzhgoLPR3fRY1cXNTQT8e4Y9OHI6FNc5SsoEkky2j-bJkmHWSx)
35. [researchgate.net](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGLJU1oqs1A6_yVwonAJd4A1l7Tvs93hAYM_JNAsj7iu4i81HR-mNc6L7Dly4FT99fpJC6oZ6tyWVsYIVeKxMTxhldNjYCjdEG77QRFTIIssS9M2GGGcLpAKjnD5X9petnzrJWWw8PQ5ZfqvA55ZYCwpbvbJkmSvNp0VtVGgVv_URntF9joTzdliDeLUmrmHRl_zxLq5V9xqLyo30xVtqSly7LsjRTAkIewIfBfcnnA1XeWM_fqaLYSVH196tdnwpUtf4_u7dytWrJL-63t30AeFO8b1n6IVxopwrdCvQ==)
36. [metr.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEO3J9eBHEePAL1sQcBniXLtvNL4rOTJT5XsEktVlrJa6HzoOL4tb3LzdHHxAyg2-XzICHgr_nRdJbLWoa40pYzXyYpDnLzI-proP1eSllhzaqJDA66x0UTTIbkD0iKR2m_PrYRYWn95i37WWNL8SkgUeUDxjsR)
37. [nih.gov](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHysvJ_kb3-R_RcWXnFKHsxpDP4zugTCZWuknGRyEK2_FlOh_3IFMybZqbn_tjv5tZqoXvqv2-NEZitSQtvKLnTedJpthkodJsTzPUS527VuDLkPCPeSfhJjjsFOy55ZHk_zzXi0HUq)
38. [mdpi.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHXzTjxKIHsjvYy4cGxf10OTi0vDTcQXlc_MlYPbqt3lAFifL8wnrgW79R205iw6sxgYyMIo6wo5lg-ZBKFh2fyO9wZjlKKVShp5Jthe6-t7UDVKD89o8UFS0EuwuE=)
39. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFDHMbnMZ_ndcFdbWdkAQSV_xR_ghpj_ywwDckwThwKvx-JJ4Vi8l61zI-dYzXEJMsHllk-pRr3kqFaCy4MzbmUPTFwMmDwO193E0baxgefEIuTODuDAWc=)
40. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEZMUNfp5reQVuHMcCvmXQQy9zvaA_U0SOyxFChKmhEvavffDIafIqME1Hj9wFsz2oD8skiuM6xdqeK12J1EtwfcYFn9vdT9mB6m1JWhhnS04hpU0jUaQ==)
41. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGCXftLH0ZvWO5-GoLTvGZcjwILofREW4IEaQp9TQQn77kqGUe_fMlBGCCRZYsOGgpDYriEeN_uxoD6aJdDqrmF1AhoTusXbw876ena68_5p1R-H7Sseei12w==)
42. [medium.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGXY0eeXkuVylefK3cASO_J5CqhwL-18SDSWcDSi96kPCmobUKi0jl4_DHLe-hfwGSV3u7GgAPKdwbww06wywKzOLe4NsGBDBqZn-s4NIIDGxcDq2b4vJvwjCNAUK1DLvK_mZjydRuHNQGogiPKmc43gyGfyN0SWfnPnzbux144bvmx1d0=)
43. [noaa.gov](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG3kLvtQBMQQR7WRvn3QB_EwUTN-9jdKZJMODbOQtzWQPsWl0x5V7BeIPeYALCmo8cstJolkwoC8UQ30qjgIOvWsoEEzPLBqufhc1QD5FjPjkpwaFqlBQU7KYIa9j6ohEIBQiAYxNmTTv9zBvkhmTRCPbzlzNd4X4mHD26B)
