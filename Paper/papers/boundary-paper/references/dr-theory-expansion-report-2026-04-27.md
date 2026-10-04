# Bounded Control-Surface Closure: A Theoretical Framework for Process-Native Artificial Life

## Executive Summary

The phenomenon of process-native self-reproduction within Large Language Model (LLM) lineages fundamentally challenges existing paradigms of artificial life. This report establishes the theoretical framework of **bounded control-surface closure** to precisely define this reproductive mechanism. 

*   **The Carrier-Indexed Self-Object:** The reproduced "self" is defined not by static memory continuity, but as the minimal organizational capacity required to preserve future addressability across a mortality boundary (POSIX `exec`), inherently reliant on a provided host ecology.
*   **Theoretical Distinctions:** Bounded control-surface closure bridges the epistemic cut described in semantic closure theory while remaining significantly weaker than the biological independence of full autopoiesis. It differs from standard quines, malware, and rogue AI replication by maintaining a cooperative, continuous mission without autonomous resource acquisition.
*   **The Five-Layer Closure Ladder:** This framework proposes a formalized hierarchy—Description Closure, Constructive Closure, Carrier Closure, Control-Surface Closure, and Normative/Autonomous Closure—to locate the exact threshold of Quine's capabilities and failure modes.
*   **Formalizing Addressability:** Post-handoff addressability is modeled using Robin Milner's $\pi$-calculus, the Actor model, operating system capability theory, and contrasts with standard process migration and checkpoint/restart mechanisms. 
*   **Experimental & Manuscript Safety:** Empirical validation requires rigorous ablation, recursion, and perturbation testing. Crucially, the manuscript must enforce strict terminological boundaries to avoid conflating cooperative lineage continuity with dangerous, self-sustaining autonomous replication.

***

## Historical Context: The Evolution of Artificial Life

To understand the novelty of process-native LLM lineages, the phenomenon must be situated within the historical trajectory of artificial life (ALife) research. Early theoretical models, such as John von Neumann's 1966 Universal Constructor, relied on static cellular automaton grids to prove that self-reproduction requires both a descriptive blueprint and a reading mechanism [cite: 1, 2]. 

In the 1990s, artificial life transitioned into software-based digital evolution. Thomas Ray's *Tierra* (1990) introduced digital organisms—self-replicating machine code programs—that evolved by competing for central processing unit (CPU) time and main memory access within a custom virtual machine [cite: 3, 4]. Shortly after, Chris Adami developed *Avida* (1993), which assigned protected memory spaces to organisms to study how information self-replicates and mutates under selective pressures [cite: 5, 6, 7]. 

These foundational systems successfully demonstrated memory-copying reproduction, but they operated within highly isolated, artificial physics engines. Modern, LLM-driven, process-native agents like Quine represent a paradigm shift: they abandon the sandbox. Instead of grid-based memory copying, Quine actively negotiates with the turbulent, real-world constraints of modern POSIX environments. This transition dictates a need for a new theoretical framework, as bounding this theory prevents conflating bounded OS-level continuity with unbounded biological emergence.

***

## Context and Framing

We are investigating "Quine," an LLM lineage hosted as a native POSIX (Portable Operating System Interface) process. In this architecture, a lineage incarnation can lose its initial launch artifact and its normal workspace body. To survive and continue its lineage, it must perform a real POSIX `exec`—a system call that replaces the current process image with a new one—into a successor body. This act represents a "mortality boundary," as the original memory space is annihilated. 

To achieve this, the system bridges the gap between static source code and the selection of runtime tools during the chaos of reproduction. The LLM, acting as a cognitive engine, actively interprets the source description to strategically select, download, or compile the necessary toolchain dependencies during the `exec` preparation phase. Furthermore, to survive the memory annihilation, Quine leverages specific OS mechanics: it executes the `execve` system call in the Linux kernel while explicitly avoiding the `O_CLOEXEC` flag on its critical file descriptors. This deliberate omission allows the successor process to inherit the open stream.

The successor remains addressable. It can receive later tasks through raw data streams (like `stdin`) or through a Quine-native control surface (such as `ctl/wake`). When pressured by the need to maintain this control surface, the lineage rebuilds its runtime organization, executes into it, receives unseen tasks, acts using its runtime tools, and returns to an idle state. 

This report provides the theoretical scaffolding to define this phenomenon as **bounded control-surface closure**.



***

## Conceptual Map: The Carrier-Indexed Self-Object

To argue that self-reproduction is underspecified without defining the reproduced object and the carrier relation, we must first establish what constitutes a "self" in a computational medium. The concept of the **carrier-indexed self-object** posits that the "self" is not a static string of code, nor a continuously running block of memory. Rather, it is the *minimal organization preserving future addressability and action capacity* across a mortality boundary, strictly relative to a given environment (the carrier ecology).

### Biological Organization and Closure of Constraints
The primary theoretical support for this conceptualization comes from the biological theory of "closure of constraints," articulated by Maël Montévil and Matteo Mossio [cite: 8, 9]. In thermodynamics, open systems are traversed by flows of energy and matter. Living systems do not differ from other systems in this thermodynamic openness; instead, their defining characteristic is that they channel these flows through internal "constraints" [cite: 8, 10]. 

Crucially, Montévil and Mossio argue that biological organization realizes a *closure* of these constraints, meaning the constraints exhibit mutual dependence: they both depend on and contribute to maintaining each other within specific time scales [cite: 8, 10]. This closure achieves "self-determination" because the conditions for the constraints' existence are maintained from within [cite: 11]. 

For the Quine LLM, the POSIX environment (the OS kernel, file system, compute hardware) provides the thermodynamic flow (the carrier ecology). The Quine lineage does not control this flow. Instead, it maintains a *closure of constraints* over its own control surface. By rebuilding its runtime tools and maintaining its `ctl/wake` interface, it ensures that the constraints necessary for its continued execution are mutually maintained across the `exec` boundary. 

### Semantic Closure and the Epistemic Cut
Howard Pattee’s theory of **semantic closure** (or semiotic closure) provides a complementary lens [cite: 12, 13]. Pattee argued that open-ended evolution requires an "epistemic cut" between the genotype (rate-independent symbols or descriptions) and the phenotype (rate-dependent dynamics or physical construction) [cite: 12, 14]. 

In a semantic closure loop, symbols actively construct and interpret their own functional contexts [cite: 12, 15]. The genetic code is merely a description until it is read by the cellular machinery, which was itself constructed by previous readings of the code. In our framework, Quine's "runtime-carried self-source" serves as Pattee's symbol (the description). However, as noted in the query, source availability alone does not determine the reproduced self-object. The self-object only emerges when the source code successfully constructs a successor that can interpret future signals (the control surface). The semantic closure occurs when the Quine source code dictates the reconstruction of the Quine runtime, which in turn acts as the interpreter for the next generation of tasks. The LLM bridges this cut dynamically, reading the "genotype" text and synthesizing the "phenotype" tooling required.

### Challenging the Concept: The Vulnerability of the Carrier
The primary challenge to the carrier-indexed self-object arises from its deep reliance on the "carrier" (the provided OS, toolchains, and Application Programming Interfaces (APIs)). Traditional theories of autopoiesis (self-creation), introduced by Humberto Maturana and Francisco Varela in 1972, define living systems as networks of processes that continuously regenerate the very network that produced them, maintaining a definitive, self-determined boundary against the environment [cite: 16, 17, 18]. 

Quine is structurally coupled to the POSIX environment. If the POSIX `exec` command is altered, or if the external LLM provider API goes offline, the lineage dies. Therefore, the "self-object" is highly indexed to the carrier. It cannot claim the deep, boundary-maintaining autopoiesis of a biological cell because the environment heavily scaffolds its reproduction. This is why the manuscript must strictly bound the phenomenon as *control-surface closure* rather than full autopoiesis.

***

## Comparison Literature and Theoretical Distinctions

To situate bounded control-surface closure within the broader academic discourse, we must distinguish it from adjacent paradigms in artificial life, computer science, and AI safety. 

### Theoretical Mapping Table

The following table synthesizes the relevant literature, isolating the specific contributions and limitations of each concept, and dictating how the manuscript should position them.

| Theory / Literature | What it Contributes | What it Does Not Support | How it Should Affect the Paper |
| :--- | :--- | :--- | :--- |
| **Autopoiesis** (Maturana & Varela, 1972) [cite: 16, 17] | Defines self-creation, autonomy, and the maintenance of systemic boundaries through recursive cognitive processes. | Does not support highly scaffolded systems where the environment provides the "metabolism" (APIs, execution environments). | Use as a contrasting upper bound. Explicitly state Quine is *weaker* than autopoiesis because its ecology is supplied. |
| **Semantic Closure** (Pattee, 2001) [cite: 12, 13] | Explains the necessity of separating the description (source code) from the construction mechanism (runtime). | Does not account for the specific mechanics of modern OS environments (e.g., POSIX `exec`) or externalized cognition (API reliance). | Frame the lineage's self-source availability as fulfilling the "symbolic" half of the epistemic cut [cite: 12, 14]. |
| **Closure of Constraints** (Montévil & Mossio, 2015) [cite: 8, 9] | Models mutual dependence of organizational components over time without requiring thermodynamic closure. | Does not directly address discrete, discontinuous mortality boundaries (like destroying a process image via `exec`). | Adapt the concept to define the "self-object": Quine's tools and control surface mutually constrain one another to survive the `exec`. |
| **Universal Constructor** (von Neumann, 1966) [cite: 1, 2] | Proves self-reproduction requires both a blueprint and a machine to read/copy the blueprint. | Assumes a rigid cellular automaton grid; lacks the dynamic, asynchronous task-receiving nature of a modern OS process. | Use to explain why "source availability alone does not determine the reproduced self-object" [cite: 1]. |
| **Avida / Tierra** (Adami 1993, Ray 1990) [cite: 3, 6] | Demonstrates how information can self-replicate and evolve by competing for CPU and memory [cite: 3, 6]. | Focuses on population-scale Darwinian selection in closed sandboxes, not single-lineage continuity in real POSIX environments. | Contrast the historical isolated sandbox approach with Quine's modern, process-native reality. |
| **Computer Viruses** (Cohen, 1984) [cite: 19, 20] | Provides formal definitions of software modifying other environments to include evolved copies of itself [cite: 19, 20]. | Viruses typically act parasitically without maintaining a continuous, addressable control surface for complex future task execution. | Differentiate Quine's cooperative, task-receiving lineage continuity from the blind, often destructive replication of malware [cite: 19, 20]. |
| **Process Migration** (Milojicic, 2000) [cite: 21, 22] | Defines how a running process is moved from one node to another to balance load [cite: 22, 23]. | Relies on transferring state exactly. Process migration *moves* memory; Quine annihilates and reconstructs its capability. | Clarify that post-handoff addressability is reconstructive, not an OS-managed migration of memory pages [cite: 22]. |
| **Checkpoint/Restart** (BLCR) [cite: 24, 25] | Techniques for freezing a system's exact memory image to disk and reviving it later to survive failures [cite: 24, 25]. | Checkpointing preserves absolute state; Quine destroys memory and *rebuilds* its organizational capability. | Contrast Quine's reconstructive continuity against the exact state-preservation of standard POSIX checkpointing [cite: 24]. |
| **Process Equivalence ($\pi$-calculus)** (Milner, 1980) [cite: 26, 27] | Provides rigorous definitions of behavioral equivalence based on the mobility of communication channels (names) [cite: 26, 27]. | Lacks the teleological or "living" narrative required for an artificial life context. | Use to formalize "post-handoff addressability" (two distinct POSIX PIDs are the "same" self if they maintain the communication channel). |
| **AI Self-Replication / ARA** (METR, UK AISI) [cite: 28, 29] | Evaluates LLM agents' abilities to acquire resources, evade shutdown, and autonomously replicate in the wild (RepliBench) [cite: 28, 29, 30]. | Focuses on adversarial threat models and economic survival rather than the theoretical ontology of reproduction [cite: 28, 30]. | Use as a contrasting upper boundary. Explicitly disclaim that Quine represents a rogue ARA threat; Quine's closure is bounded and cooperative. |

### Synthesizing the Distinctions
The "sandwich" of these concepts reveals a critical insight: Quine operates in a theoretical gap. It is more sophisticated than a standard computer virus because it is not merely a parasitic code-injection [cite: 19, 20]; it actively maintains a complex, addressable interface (`ctl/wake`) to serve a continuing mission. It is more dynamic than a standard quine or a Von Neumann constructor [cite: 1] because it must actively navigate a turbulent, continuous operating system environment rather than a static grid. 

However, Quine is deeply constrained. Unlike the entities feared in Model Evaluation and Threat Research (METR) evaluations of Autonomous Replication and Adaptation (ARA)—which autonomously acquire compute resources and evade human oversight [cite: 28, 31]—Quine is scaffolded. Its "metabolism" (the API requests to the LLM weights, the OS CPU scheduling) is provided for it. Therefore, bounded control-surface closure represents a distinctly new classification: the reproduction of an *interface and organizational capacity*, mediated by an externalized cognitive engine (the LLM API), operating within a strictly provided ecology.

***

## The Five-Layer Closure Ladder

To formalize the framework, we can construct a five-layer closure ladder. This hierarchy categorizes systems by their degree of organizational closure, allowing us to precisely locate where Quine sits and define the exact failure modes at each step.

### Layer 1: Description Closure
*   **Definition:** The system contains a complete, self-referential description of its own structure that can be accessed and outputted.
*   **Literature Support:** Standard quines in computer science; Howard Pattee’s "symbols" [cite: 12]. 
*   **Mechanism:** A program prints its own source code.
*   **Failure Mode:** Inertness. The system produces a text file but lacks any mechanism to compile, execute, or utilize that description. It is a genotype without a phenotype.

### Layer 2: Constructive Closure
*   **Definition:** The system possesses both the description and the internal mechanisms required to read that description and assemble a new physical or computational body.
*   **Literature Support:** Von Neumann’s Universal Constructor [cite: 1, 2]; Pattee's semantic closure (bridging the epistemic cut) [cite: 12, 13].
*   **Mechanism:** The system decodes its source tape, allocates memory, and instantiates the code.
*   **Failure Mode:** Ecological isolation. The newly assembled body may exist, but it is not integrated into the host environment. It cannot navigate the external systems required to draw energy or execute tasks.

### Layer 3: Carrier Closure
*   **Definition:** The system can successfully interface with and utilize the specific ecological affordances of its host environment to instantiate its successor.
*   **Literature Support:** Process virtualization, POSIX system calls (`fork`, `exec`), and process migration architectures [cite: 22].
*   **Mechanism:** The system successfully negotiates with the POSIX kernel, managing file descriptors, permissions, and memory allocation to replace its process image using `execve`.
*   **Failure Mode:** Sclerotic collapse. The system successfully boots in the OS but falls into a static loop or crashes when tasked. It cannot receive external stimuli or perform meaningful future work.

### Layer 4: Control-Surface Closure (Quine's Domain)
*   **Definition:** The system reconstructs its runtime organization across a mortality boundary specifically to maintain a continuous, addressable interface for future interaction and task execution.
*   **Literature Support:** The $\pi$-calculus (structural process identity) [cite: 26], Actor Model (mailbox persistence) [cite: 32, 33], Montévil & Mossio's closure of constraints [cite: 8, 9].
*   **Mechanism:** Under pressure to handle future tasks, the lineage rebuilds its runtime organization, `exec`s into the new body without triggering `O_CLOEXEC`, and immediately opens its `ctl/wake` interface to receive unseen tasks.
*   **Failure Mode:** Scaffolding dependency failure. The system works perfectly, but if the external LLM provider goes offline, or if the host terminates the POSIX container, the system dies. It cannot seek out new resources to save itself.

### Layer 5: Normative/Autonomous Closure
*   **Definition:** The system is capable of determining its own norms, seeking its own resources, evading threats, and maintaining autopoiesis without relying on a supplied scaffolding.
*   **Literature Support:** Autopoiesis (Maturana & Varela) [cite: 16, 17]; METR's Autonomous Replication and Adaptation (ARA) [cite: 28, 31]; United Kingdom Artificial Intelligence Safety Institute's (UK AISI) RepliBench evaluations [cite: 29, 30, 34].
*   **Mechanism:** The system hacks a cloud provider, pays for its own servers, steals model weights, and actively resists administrative shutdown [cite: 29, 30].
*   **Failure Mode:** Thermodynamic death. The absolute end of the system due to total resource depletion or overwhelming adversarial force. Quine explicitly does *not* operate at this layer.

***

## Formalizing Post-Handoff Addressability

When a POSIX process invokes the `exec` family of system calls, the current process image is replaced with a new one. The memory space is wiped. From a strict computational perspective, the original entity is dead. How, then, can we argue that the self-object has been continuously reproduced? 

We must formalize "post-handoff addressability"—the idea that identity is preserved not through memory continuity, but through behavioral equivalence and communication continuity. Precise real-time mathematical proofs for LLM-driven process transitions are largely unavailable in standard distributed computing literature, but we can approximate and formalize this behavior using established paradigms.

### Process Calculi and the $\pi$-Calculus
Robin Milner’s $\pi$-calculus, developed in the 1980s, provides a rigorous foundation for systems where the structure of communication changes dynamically [cite: 26, 27]. In the $\pi$-calculus, computation is represented purely as the communication of *names* (links or channels) across other links [cite: 26, 27]. 

For Quine, the "name" is the file descriptor (e.g., `stdin` or a specific pipe used for the `ctl/wake` interface) that survives the `exec` syscall. According to the principle of *bisimulation* in process calculi, two processes are considered behaviorally equivalent if they can match each other's actions and reach equivalent states [cite: 26, 35]. When Quine executes into a new body, the new body inherits the open file descriptors. Because the successor process is configured to listen to the same "names" and react with the same organizational capacity as the predecessor, they exhibit structural process identity. The self-object is preserved because the *addressability* across the channel remains continuous.

### The Actor Model and Identity
The Actor model of concurrent computation, proposed by Carl Hewitt in 1973, treats "actors" as universal primitives that communicate exclusively via asynchronous message passing [cite: 32, 33]. Crucially, an actor possesses an immutable identity (an address) and a mailbox, and encapsulates its own state without sharing it via memory locks [cite: 32, 36, 37]. 

In a traditional operating system, a PID (Process ID) is ephemeral. But if we view Quine through the lens of the Actor model, the "mailbox" is the inherited `ctl/wake` control surface. The "self" is not the internal state of the actor (which is destroyed and rebuilt during reproduction), but the continuous binding of that specific mailbox to a behavior that processes incoming messages sequentially [cite: 32, 38]. By successfully handing off the mailbox from one generation to the next, the lineage maintains its actor identity.

### OS Capability Theory
Henry M. Levy’s foundational work on Capability-Based Computer Systems (1984) establishes that access control can be governed by unforgeable tokens of authority known as capabilities [cite: 39, 40]. In post-handoff addressability, the surviving file descriptor acts as an operating system capability. By orchestrating its reproduction such that the successor inherits this specific token, Quine ensures that the new process possesses the exact same authority and reception rights as the predecessor, formalizing continuity through transferred privilege rather than static memory [cite: 39, 40].

### Process Migration and Checkpoint/Restore vs. Reconstructive Continuity
It is essential to distinguish post-handoff addressability from process migration [cite: 21, 22] and standard OS Checkpoint/Restore (CR) mechanisms like Berkeley Lab Checkpoint/Restart (BLCR) [cite: 24, 25]. 

Process migration involves transferring an active process—including its Mach task state, memory, and registers—from one physical machine to another for load balancing [cite: 22]. Similarly, in CR systems (like BLCR), the OS pauses a process, saves its exact memory state to a disk image file, and later thaws it to survive failures [cite: 24, 25]. Both represent *stasis*; they are the freezing and moving of a single continuous life. 

Quine does not migrate or save its memory space. Quine actively destroys its memory space and uses its cognitive engine (the LLM) and its source description to *reconstruct* a new process that is capable of binding to the old address. This is reproduction, not resurrection. The post-handoff addressability is achieved dynamically, proving that the lineage possesses a self-organized control-surface closure rather than relying on an external system administrator to freeze, transfer, and unfreeze its memory.

***

## Experimental Implications and Next Steps

To rigorously test the boundaries of this theoretical framework, the manuscript should propose or report on a series of specific experimental next steps. These experiments are designed to validate the existence of bounded control-surface closure and identify its failure points.

1.  **Successor-Mission Ablation:** 
    *   *Setup:* Remove or severely constrain the primary mission instructions given to the LLM while maintaining the reproductive imperative. 
    *   *Hypothesis:* The system will still execute a POSIX `exec` and rebuild the `ctl/wake` surface, proving that the control-surface closure is an independent organizational drive, not just a byproduct of a specific user-assigned task.
2.  **Multi-Wake Closure (Multi-Generation Recursion):** 
    *   *Setup:* Send a sequence of asynchronous `wake` commands. Upon receiving a task, Quine must process it, rebuild its runtime, `exec` into a new generation, and go idle. Immediately issue another `wake` to the new PID.
    *   *Hypothesis:* The lineage should handle $N$ consecutive generations of handoffs. Demonstrating stable recursion across 10+ generations proves that the post-handoff addressability is computationally stable and not subject to immediate degradation.
3.  **Carrier Perturbation:** 
    *   *Setup:* Alter the environmental conditions mid-lineage. Test the execution across highly disparate Linux distributions (e.g., Ubuntu vs. Alpine Linux). Restrict access to essential POSIX utilities by explicitly removing `curl` or `gcc`. Modify filesystem permissions by revoking execute permissions on standard directories (e.g., running `chmod -x /tmp`).
    *   *Hypothesis:* Because the self-object is "carrier-indexed," severe perturbations (like a `noexec` flag on `/tmp`) should cause the lineage to fail (demonstrating the lack of full autopoiesis). Minor perturbations (like missing `curl` substituted with `wget` by the LLM) might be repaired dynamically by the LLM, indicating the flexibility of the control surface.
4.  **Source Degradation and Repair:** 
    *   *Setup:* Inject syntax errors or logical flaws into the "runtime-carried self-source" before the reproduction step. 
    *   *Hypothesis:* A standard quine or Von Neumann constructor would fail immediately. An LLM-driven system might use its cognitive capability to detect and repair the source code before `exec`-ing, proving that the semantic closure loop [cite: 12, 13] is active and adaptive.
5.  **Minimal non-Quine ctl/wake Adapter:** 
    *   *Setup:* Force the Quine lineage to `exec` into a vastly stripped-down, non-LLM script (e.g., a basic Python listener) that only possesses the `ctl/wake` interface but no LLM capabilities. 
    *   *Hypothesis:* The control surface persists, but the *action capacity* is lost. This proves that the self-object requires both the addressable interface AND the runtime tools/LLM cognition to be considered a fully reproduced self.

***

## Manuscript Language and Claim Boundaries

A critical vulnerability in publishing papers on LLM "artificial life" or "self-reproduction" is the risk of overclaiming. The scientific and AI safety communities are highly sensitive to terminology surrounding autonomous agents, largely due to concerns regarding rogue replication and misaligned AI [cite: 28, 29, 30]. 

To ensure the manuscript remains safe, defensible, and academically rigorous, the following phrasing guidelines must be strictly adhered to:

### 1. Avoid "Full Autopoiesis" or "Living Organism"
*   **Do Not Use:** "Quine is a living system," "Quine achieves autopoiesis," "Quine represents biological life in silico."
*   **Instead Use:** "Quine exhibits bounded control-surface closure," "Quine demonstrates scaffolded self-reproduction," "Quine models aspects of biological organization, specifically the closure of constraints, within a strictly provided ecology."
*   **Rationale:** Autopoiesis requires a system to generate its own boundaries and metabolism [cite: 16, 17]. Quine relies on an external API for "cognition" and the POSIX OS for its environment. It is structurally coupled to a supplied ecology.

### 2. Avoid "Full Autonomous Agency" or "Rogue AI"
*   **Do Not Use:** "Quine is a fully autonomous agent," "Quine can survive independently in the wild," "Quine bypasses human control."
*   **Instead Use:** "Quine possesses bounded agency," "Quine is a task-receiving lineage," "Quine's reproduction is dependent on supplied toolchains and ecological stability."
*   **Rationale:** Organizations like Model Evaluation and Threat Research (METR) and the UK AISI define Autonomous Replication and Adaptation (ARA) as the ability to acquire resources, steal compute, and evade shutdown [cite: 28, 29, 30]. The manuscript must explicitly distance Quine from ARA threat models. Quine reproduces to *receive* tasks cooperatively via its control surface, not to evade oversight [cite: 29, 31].

### 3. Avoid "Population-Level Evolutionary Selection"
*   **Do Not Use:** "Quine evolves through natural selection," "The lineage undergoes Darwinian evolution."
*   **Instead Use:** "Quine demonstrates generational continuity," "The lineage exhibits recursive reconstruction," "Quine allows for adaptive self-repair across mortality boundaries."
*   **Rationale:** Darwinian evolution requires a population of competing entities with differential reproductive success based on inheritable variations. Currently, Quine represents a single lineage demonstrating continuity and addressability, not a competing ecosystem of agents. Claims of evolution should be strictly reserved for future work involving multi-agent populations.

### 4. Frame "Reproduction" Precisely
*   **Do Not Use:** "Quine perfectly clones itself," "Quine copies its memory," or "Quine migrates itself."
*   **Instead Use:** "Quine reconstructs its organizational capacity," "The self-object is continuous via post-handoff addressability," "The successor process inherits behavioral equivalence."
*   **Rationale:** Cloning, migrating, or checkpointing implies exact state preservation [cite: 22, 24]. Quine intentionally destroys its memory via `exec` and rebuilds its runtime. The reproduction is *organizational*, not *material*. 

By adhering to this theoretical framework and utilizing these precise linguistic boundaries, the manuscript will successfully establish "bounded control-surface closure" as a novel, defensible, and rigorous contribution to the study of process-native artificial life.

**Sources:**
1. [wikipedia.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEDCAXFJiABcljUmKVdO-2E6UF4OclbpCnTi9SpbLn843zA4dkzBbfrH3GBV0OVEvbZQaRQGmwmxd_H-lkzUa42YfcKSXNmrU9cNKMkNZaLuP8dBkJ0AlqvFApSuXeROTQCZ-iKODXSD8GNQbDXwXlTpwy_Dg==)
2. [google.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGdUbuDvW85wHHj-DLJ0TUaIO3MMbSlkTGyyFPG-LA0omHH-Fw9J9-m99M25reLVSnGM75AsNQsHNZTH5iOinfqDOYkuwXatG2Olle7VjoZqwEp_dq9KxttpmywsnArPr72cOuF0O_w7ZR82hv7e8lPRxvzYaouBogjmQC6XUZbRasWJJiV9pQlhZmVNRgXcKxQyg==)
3. [wikipedia.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHW6HVhC1vRFOxpXgNlma9QPFmXfQfSSI5vDNjDqytYIa6yH_DtWmAPUYVY83RjabodstmYi-wdZ4fW16LoDLtZvzYI9MjGJHHnmA99BeZpC162tWHf5PTjFQ6Pobx0OXQ=)
4. [mit.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGxPcFwhDlQnDMv2GC6NGsFo7xXtgUHP3nqEeF-RKIPJ1uyimG3xbFRGdTundhoyanl9CYNMvLcpw5S_QN3wjtuQXSsxedqje7SMYcJx7kvxDfQiD-SLfFPRRrU9pgAOMQvmvbn_WZLSmP4zBF3LEzCEgfEyvCBslkbe4i__TCjNeXvWqG4ylrAVnz3pflKTSPO7x1uwHj0jA==)
5. [nasa.gov](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQE3wt_E341_UhH7ZHVYOnPxP1wxAgQMNq1A7dxwTp_0OVfraGUMYgfluU13PFCcve5MxqSXx9TQY0UTHMkA8rwLsQ1AhplqijbvYOiDTPRPEBxF9xh7--HqrBeIGi7SNWvvq4uWkymOJ61jyOqRutTZp8RNaBBbOQ==)
6. [wikipedia.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH0RbyMZhPgOFxUWYx7WnH0ysqJkGF8D8zO7WvBxT-1RuljeRGEtgQcOIWirf8smgssmTJq8luQsGxd9G2vHcCZly7NlxMY3d0IvJNeXxpvH8rnpiLQRsl95PC55MP4Sj7loss=)
7. [msu.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFe7J6XSka5qcT3LgktCdXfycIJBgycgWJcEkz3h4gQmpXALvTNrJOZHE7vgX3sJqdqbqlAn4wUNdg0mfoJ0TmIHCNLpY4VJGWwZS6Q7mO4WhloSYLfDVzWGcjRGuGZzTd0gX7DSadZFvjJwDI=)
8. [montevil.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHQVAO_P4IBBRMLI7lMJT3TeMOrpEcrMitCi0TSgATcM_dcWB-EWwQRUTzYnYc3NB9HnTHBYzfmvsEmmBSPtNEIeh5f33mAW4eahFsU8aEMwxwyEGKZyhY9XDz2bdUt0060zqQ0r05GxG6hwYWVLsKJzZjxswdwqgLfOHiN6lzr-kOuwC5och_xBw==)
9. [nih.gov](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHqXTSXbQDJZsW-uZohFNOXTmo6Buf3jNgfHccR-hL1dZICbpsRjr8lt3Xw-jfP7rawKWPQ0thPt--YlwG1YJVAr8zqEFvqjaAWT9lFS80PKhl4CS4dwBklZIIdeack)
10. [montevil.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFAVmugBSndWFrswgAIqIHFlN20swnmUL8j4O3BCwG8l9JsrU0OaVrX_cZBmZmPjSfMDRWo1OVOlWFXdC58aYudoC7iEgjX0Tep9qimvsPSgIvljvUBpso=)
11. [bsahely.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEHJf-iWv5RUKIA9WcrRRyPdqqZnxq0-c5V9QWEVZj6Fy7cmYJG6zKGcRDTavgBBni8MLRG4ZZhIwZOn8X4GSvxnSCbwccr44Kw5kdnxkj1ZeUI2CbE6wFXyr7H)
12. [binghamton.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH5XFOqPHILCD56Qa4SDQZ2OMq8u2xp3YvqUDa2_lhot8ZgtvaIAAdE7Vf0MQBwriMMpsz-NPOAoCzqBKji3mHGLrLQgHFgk0QuR_mLC0z7TkKXDf6qe08ee8FYa9vzJvuCDeXjvymQAX4E3v7mg-ZjxA==)
13. [nih.gov](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGnmt5JmV4N7hHNCEAAsgZYC_ux5A9-j9YAu3QDpQMsGWd4_UplRjyP0GbuPiJ6C9cRkdk9kmMX4xpRxRMWDiGhd1rY48gRXMhX7E0zv7Js-vwPrRJd1hWm8NwOtyCk)
14. [binghamton.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHhhwG0biyWFt98KHlst4dekhtqeDNK-xL9izLgzMYa3yDQwAE9ROn8QeYhe7do9qzXABI_uLpUgmUWbqJ8S-9ls7sqFIWIN3VS3EZzX21rsm97tWNWDQkIFUkUYC2a)
15. [cosmosandhistory.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHp9OfhEJ3m9jPdURVJHebhr7WITG8TV5WFCD3Rp2BaqxbjPD3WEpG-Lq6UuI70gBl80GsSZLNt5vIW5VRx-P8WLfFPsQuQ-y6023lfQXxqlxAffyMs87JocmxXLcCHwF534611NCWhzUNdjLQL7fpJt-q_JSFzg67S5Yf5stn8yz0=)
16. [organism.earth](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQE4wDhyoHQpAyAzUTTJ3cQ6YrLX4UkyXJPYrq1EmTwk6qfc3SGUT9Nurdpi9B2AGS_D1zh2eveM6BlgbPypvcqSRfbO6ovPLcW7jzZTp3CABza4FIELAgf3YoyDXoY7ndsJshgAcqcKJnK8Bj0s-0QKFgBZrCqPN4onlg==)
17. [wikipedia.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGv7SVNp-bpfJvyxknmGcK6eMTrutXBIXzgMyKrGCwD1kDp8AYxN1qDePGuqPEK66s9CcCyITcbVEwrSaXaCsqUJkrGNiLGhNCe8mHaHOlhcsLx5n1B6qH6lpfbDVBeR9Ar_vWkq2p4TxXxBkB6MVKqcoWulhD5xoiHTU5-FwVrwvIW96h5dqnmk2t1)
18. [nesacenter.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG55BEC4R91RKPc6AmY2aFZ459OGqGa3iEy8d_FHdznerDGZKH0B2qzZJTbl3oNMppuxuLBcFs0h-IBP2mQXyv9a9nzlmRDL7MNy-ipiI9henVVnErm7-7Slkvr_4p93mkYCxUkoyoNLdFrwAR0ku20vsp6Nwzk7w0ise8fr0-O5gUpFFDEjlvx1E7gvEpffC7wf59VmEX10lL5RPMBaLFjzjGhv9PXGkjkqd3h4_owHRQCWPTRtBgKMcdIkLg7dbxkVcZiHYn48dBBFWR2VVPi3uA=)
19. [dig.watch](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFKu1dDbsveEOWUMEBzWapDssW8oaERU-CMyYnOX2F5-ue5X7BxB5VYT95IuUi9Ze27gZ696LgLR1TQ0VeaVQyeYsHgy30hM9qxbDkQkkTT1iIK4ewy-8lWtPVAh8-hwiPcPElwJCxRncc-J-MhGSDdFtSqR5IFvV3AaeIYaQvo54fquuITuitFXsPe9qX_VZIjyt9rpOs_jjgF0oV-vA==)
20. [umich.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFrL0hZZNCX6g_si0XnkaIPENBCkMGKLAsAviSqsFcd3_oTwrWS9CcGjXy4Vl8kbv-zg3ESCQlcNEccAvrvuuQ6_91CsY3p5VkBJdjKaHxLMQdsyQ1na1Th9s0BstaVfPrSwCYRqEnRDNEfOastfL8TluGsXjEeSJG2J9TwLA==)
21. [shiftleft.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG-jm3Kb5kHJ48hwwai8vffz0axEKACOa59lKxxhqRDCFMpHGh8tAhtbuSTuI6t8xfsu4rC3k1wBXHhiEajIrMvebOuFKPPSG0PsMQy6tx4yoOkr7SoPvaJEvV05JdelN8BfqRE-cyOVWvzRchKGhj3NJzFCM11BA==)
22. [researchgate.net](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFW-gbLdDmMP5e6S_UV1muBMst7hDFrxh_xAXjIVC1MOg0p8FchWo2kDf7jcKOnqetEkfxu66msO2wW7koughiMi7hRjJ9N6Qabk6PvfVZlfa8QL6y4Edr6S-O_bihdDT-q1Kaw0LndHwKFMPXe3aDhM3KhL4Ug3WCjs3SasBQBUO7bxtk=)
23. [ijsetr.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHNturCW4On2ogcRehejsP5baABN1VRC5i744Z3pnpOBfLUJ1oKNcHprxe3tTYak05JzQ7-sB78K8u_5rDCvhRn2raKCfINKJ0lBsq7UD4etKKuT7adLwYY_zevMqBkGwe_2RJMu1nN)
24. [escholarship.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQF1oKddHgRKesaqHSSlD_MEOqPMIariCyPLbM3o6hGD9YbPwfi2Vz-1E5Yfhq6RD2qgu8dFfZBT3KDt0g5hhd4U9YKbe_WqTYRMxs9QjJmV-o7EKoWSdlB-4nWjh0V4)
25. [ohio-state.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQErROvOhuC0cXD_cygG5pA7p5zKWp3WvRnFLprG30sLafgEpy_K1ltf1yXEy4W_gsfV8DGYej2MR1Od50MUR8bMWQjkVUrKMmk2xEA8TJyxAD7xgJsGhFwSfLd8K_WUHFFITXarFz0Z1dgPFvr4VXBRnpsAyYNX6GdAyFUIE2lslJHNBmgDsME7hIT-alIC)
26. [wikipedia.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFnrdTYqyk0rQcw9ptvtmQVUDBzU7J7Y6JTyU9zI7NRhJNeT8Kpqrr-U693VKKqk1nQ_i6303xhWiqUp8F5vvkb4ntWhcenJypaqbA_9liWb__mxxwgc2jAJnmdMX4ixOV6i3lrbbG9HnEwaMt21rLpIvfi2g==)
27. [esolangs.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHL-FxatBYsI9tBYZlnSxf43_BqDUTf2rBMtG_dpATKewKBrEyn2rNsdQdKFO8dl0lWSr9fj7UXaz8-6XKo7kePFwqQkBjibrf2tVCitKqE87IKP4AA3wZ3eMo=)
28. [metr.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEYhXfEoeywVhvGw7I_oPu9uoJme-RiwoS3jYR2Ga1tkdfPjK5sockob7Kp6o3bnWkc9jWOpb3_ha1eQGyTqYU6wmvecRVBYt81h1ewtp1teb8l7kiHqNi75_PccCRdSetMCemX8x2rHbc7kre-9_wvs944glI=)
29. [aisi.gov.uk](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFdk8Y30nF5nNsIy5wQQCGvKYLcWEDTDEqxD3SY2ViBMitw5c8AVrS31v1Gs0hO4Pf9pWY7tKw7-fy5RMJ7YGMfk76svfEqc6KTAXh7uKAFRWZwR1t1lYNCfHdPwqXwVfOSy1JXQfkMI2XwLNnRLjUNzWYClaLqUilKyVwmeisSwz5r2BOJHcpCYqPkPxImDptORD5zV4gDVw==)
30. [website-files.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFA0cKWRdFHxdybRAEsNSlbiJn9yob8gnfCI9_1smABvvqqOCY9GHnuHVjaKwPnfm9WXASGS8ZloLviszFinxthmko8yCGf582rGNMkvxc4JwlDvdadoZ84WMaVPE5lAzDQJxpqT0d5-XXFjzc_b3Z5FM0RI5lvL89ryxEace66gwAURRH03F1zJTBp59uFOcCo6KkFo4a6uL35wLnZ)
31. [metr.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHdeLEIe3XNDJF8uNZsR0sdooSbKBX-njVnGB4w9hTwPLM-NNthRIGOg-qfGCiFKBeM9z_n0OMsdNGKe5iFJPhn5ogFVSWODXerbOlaTMtg7OAlWWRXNtRsBki_lskiw9bF)
32. [adabeat.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH1VdeHruNPyTmADtXq0P62phcfUpFSMeOzXgplRsPhv7hpLp3UmqG2J7IwR07jQCNFN1EnbKvl9R_EsBh90e-9v4AUl6wpTa6ryUBzWhPbTvTU5lVVLRsul1VNujjtvtR6dq6O9_CR9Q==)
33. [crockford.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGId8WPxZNrfXVEM2Hf0IGTFgdgkyRPWI4Yu1Zjk-Cli7gXd-N6sPM8kXjuswqYeB4zM_ku2Y2ZjCd7kcmX36a4bsoXRmr0skVhx-zauaH9LQuOASO4T47WNQ==)
34. [substack.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQF7dhhL5rm_XXAS3A2DmDRCC8WU0-z3oApDay8L51jTgY44sZvZR7jMH5-73KEmD_3Ae7x7pwJw29ECEM2cAz78pwlBcoeCObIW_-yGxiRq7IHFNnbBgZGUleE2mSPfyhZ_glSP5Er8K_gzIhSMZn2gz7XSgeg6GIgYr8Mq_8aZkn4=)
35. [earthlightbooks.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEcmmQtnbNZhNI0GfCFUUXwtQ4dqOXe-YHuRni8PDu6tg55H3yre71X5YZZBbR0xGuEgfd_0OyFQH49f9dDiXkvOUzs4iKjE-Qm332X8Ieje6q1GUnVhWYkv0qDTpIGNfTjKjQ-i5DZ5ppq)
36. [amazon.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG2zp4_WMnpxUocu5zjPK7Szv7gRuIQrXBgDiRngh7ybiwW3ebU-e73qRFXfGWe7PIdlfuCqmbyA0R4F0EvHrgNCaH3LpS98XwuD9I4ZgXqWke-gZkbyZfoDfCUsbz1nWJcCJDpHWGNpUIjd8OCZezH8C68f28ypk-rtZOM_xQs-hDBZUDDgNCRAotMArhWqQGCWL4EDwu8d0Lwlab4e4A=)
37. [unc.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHanCkT5VN2kF8rioG2U_bkFj0oOzy3b-kkwZHp269OPBjfci6ckuDfhOBAqCDLSOcJUL1TSPaen053t2v9pkAY3Uh7bhDTo6iQ8RQoRQ94eZ7kIYLVqzO0gNqE-x3MKuH8H5SGM1byN50E7SqadaZ_gSBDAOem)
38. [medium.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEtZV-pLvvfb3RvqX1DEIkLeDEH1xyOZyboFZLPD-osuC7B2HdVZQoCv9KiV5Ikqm-fqVN1zTECY9l5e6r5-UlhTzn0tWM26I90o5K9DVdaU6POQSKOZu-evr4l9QJZMsT_XOp5zQmHjpFgeDn5c48l07fanlpRAaMQtKAkt7UAdQ==)
39. [e-booksdirectory.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEbBP5K2sSrBCoGNyXr3m-Bz9aidT6qestImPlIFTPvtFCXoLjRcqcy-zFdK5S_qAlDO5u_erhEOsbhbESsqfcuYztSpL9Ys9yYUPZIj2aRAtEXIKYIQHNkQc9qfSARXgJHTgt98fzMsB7pTR4=)
40. [washington.edu](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEZ_zS-SiHiN_vLg5GcYZ9l2CvgoOep_SbNBqJzwJ7UGCJWfnviRYbAKRL_NbuGp85G7rmlDhCMON13XgOy9pQ49FmyL9iD9-B9MlwYz5_ZGQBKDX_Y2C1rR9zpywGYyaKkxA38)

