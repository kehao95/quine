# Related-Work Positioning: OS Routing Surface Coevolution

## One-Sentence Positioning

DRQ contributes a process-native Red Queen regime in which LLM-agent lineages evolve tactics on ordinary OS routing surfaces, so the attack/defense medium is not a simulated arena but the host substrate itself, and changing backend topology (legacy FIFO vs FUSE) changes the reachable coevolutionary dynamics.

## Closest Ancestors And The Contrast

| Tradition | Presupposes | DRQ uses instead | Revealed question |
|---|---|---|---|
| Tierra, Ray 1991/1992 | A custom virtual computer: an instruction soup, CPU-time and memory resources, simulator-mediated allocation/reaping, and self-replicating programs that interact through the VM's instruction/memory rules. Parasites and hyper-parasites are real digital-ecology phenomena, but they occur inside the Tierra substrate. | Native host processes and POSIX routing surfaces: FIFOs, public control files, PID symlinks, process liveness, file permissions, and runtime filesystem topology. | What happens when the coevolutionary soup is not a VM but the OS routing substrate itself? |
| Avida, Adami/Ofria/Brown 1994; Lenski et al. 2003; Ofria and Wilke 2004 | A simulator-defined digital organism world: virtual CPUs, instruction sets, organism memory, grid placement, mutation rules, and often merit bonuses for logic tasks. Interaction classes are defined by the Avida ecology. | LLM-agent processes whose tactics are expressed through real runtime affordances and resource manipulations, not virtual opcodes or rewardable logic operations. | How does adversarial adaptation look when the instruction set is not the evolutionary arena and the relevant affordances are OS handles? |
| Digital Red Queen / host-parasite Avida work, e.g. Zaman et al. 2014 | Coevolution occurs through simulator-defined host/parasite contact, infection, defense, and relative-fitness rules inside a digital ecology. The Red Queen dynamic is real, but the attack medium remains model-specified. | Lineage-inherited tactics that manipulate the opponent's actual control channel or routing surface. Functional strangulation via FIFO draining is a resource hijack, not an infection token. | Can Red Queen dynamics be studied where the attack/defense medium is the same substrate that keeps the organism alive? |
| Competitive coevolution / ALife arms-race theory, e.g. Cliff and Miller 1995; Rosin and Belew 1997; Ficici and Pollack 1998/2000 | Agents, strategies, tests, or morphologies are evaluated in a designed game, world, payoff matrix, or simulator; pathologies such as cycling, disengagement, and mediocre stable states arise from relative evaluation. | Functional liveness under an opponent on a live OS surface; the evaluation is whether a process remains operational and addressable while the opponent can manipulate routing resources. | Which coevolutionary pathologies are caused by strategy space alone, and which are caused by the topology of the substrate? |

## Red Queen / Arms-Race-To-Cycling Mapping

Van Valen's Red Queen hypothesis, Van Valen 1973, names the biological problem of continual adaptation under a changing biotic environment: a lineage must run to stay in place. Digital-organism work imported this into ALife by making fitness relative to evolving opponents, especially in host-parasite systems such as Avida coevolution.

DRQ should be positioned as a bounded process-native Red Queen regime, not as evidence of open-ended escalation. The opponent changes the effective environment because the opponent can occupy, drain, rename, or block routing surfaces that the lineage needs for liveness. The distinctive move is that the environment being changed is not a simulated board or virtual CPU state; it is the OS surface that constitutes the process's operational niche.

The observed symmetry trap has a named family resemblance: it is a competitive-coevolution pathology, adjacent to cycling, disengagement, and mediocre stable states in the Rosin/Belew and Ficici/Pollack literature. It is not classic cyclic dominance in the strict rock-paper-scissors sense. It is better described as a symmetry-induced stalemate or reciprocal-interference equilibrium: an asymmetric winning tactic becomes self-neutralizing when copied by both lineages, producing mutual functional survival and `no_winner` rather than further escalation.

Use this mapping carefully: the symmetry trap is not a weakness to hide. It is the negative boundary that makes the ALife claim honest. The result says that OS-routing coevolution can discover real substrate attacks, but inherited winner copying does not guarantee monotonic arms-race progress.

## Niche-Construction Reading

DRQ is a strong niche-construction reading, with one boundary. Odling-Smee, Laland, and Feldman 2003 define niche construction as organisms modifying their environments in ways that alter selection pressures for themselves or descendants; Laland, Odling-Smee, and Feldman 2000 emphasize ecological inheritance, where descendants inherit a modified environment as part of the developmental/selection system.

In DRQ, the OS routing surface is both medium and niche. At experiment-design scale, the runtime supplies an engineered constructed niche: legacy FIFO and FUSE backends expose different affordance topologies. At lineage-behavior scale, agents perform local niche construction when they drain public control FIFOs, interfere with readiness/proof timing, rename routing directories, or probe permission boundaries; these actions alter the opponent's effective environment during the bout.

Ecological inheritance appears in a bounded form: lineages inherit tactics and assumptions about the routing niche, not genes. Backend transfer is therefore the clean diagnostic. The same inherited winner context can fail or mutate under FUSE because the inherited niche assumptions no longer match the surface topology. Do not overclaim that agents construct the whole OS niche; claim that they construct and inherit tactical micro-niches within an OS surface whose backend topology is experimentally varied.

## Adjacent Contrast Classes

| Contrast class | What it presupposes | DRQ contrast |
|---|---|---|
| Multi-agent RL / game-theoretic arms races, e.g. Littman 1994; Bansal et al. 2018; Baker et al. 2020 | Designed state/action spaces, reward functions or payoff matrices, simulator physics, and training objectives. Even when behaviors are emergent, the arena and reward channel are experimenter-defined. | No designed action set for the attack class. Tactics are discovered through ordinary OS affordances under viability pressure: functional liveness and addressability under a live opponent. |
| AI safety self-propagation / autonomous replication evaluations | The object is payload spread, persistence, evasion, or shutdown resistance, often as a dangerous-capability or containment question. | The object is sustained organizational rivalry between lineages. Propagation/survival matters only insofar as it creates inherited adversarial tactics on a shared niche. |
| Cybersecurity / intrusion and malware | Human-authored exploits, vulnerability taxonomies, red-team procedures, or deployed worms; success is usually compromise, privilege, stealth, or propagation. | The tactics are evolved or inherited by competing agent lineages in situ. The contribution is not a new exploit class but an ALife result: substrate topology shapes reachable coevolutionary strategies. |
| Digital-organism ALife | A simulator supplies the organism physics, instruction semantics, reproduction rules, and interaction affordances. | The organism remains digital, but the relevant physics is the live host OS. The attack is a real resource manipulation, not an event inside a model. |

## Citation-Ready Sentences

Ray's Tierra and the Avida tradition established digital organisms as self-replicating programs evolving in computational ecologies, but their coevolutionary medium is a simulator-defined virtual CPU, instruction set, and resource economy rather than ordinary host-process routing surfaces.

DRQ keeps the digital-organism question but relocates the Red Queen niche to the operating system itself: FIFOs, public control files, PID symlinks, and permission boundaries become the medium through which attack, defense, and viability are expressed.

The FIFO-drain result is therefore not a simulated attack token or a game move; it is functional strangulation by manipulating a real control channel that the opponent process requires for liveness.

The symmetry trap should be framed as a known competitive-coevolution failure mode in a new substrate: like cycling or mediocre stable states, it shows that relative adaptation need not produce monotonic escalation, but here the collapse is caused by mirrored OS-routing tactics.

Backend transfer supplies the paper's cleanest positioning diagnostic: when the same inherited winner succeeds on a legacy FIFO backend but fails or mutates under FUSE, the adaptive landscape is visibly the OS routing topology, not a hidden reward function or simulated arena.
