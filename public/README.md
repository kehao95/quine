# Published works

The public portfolio — one self-contained entry per published or accepted work.
Each entry carries everything needed to read, reproduce, and cite it: manuscript
source, figures, the released PDF, and (for experimental papers) the experiment
data and a reproduction protocol. Inclusion is strict — **published** or
**accepted** works, plus **under-review** works whose evidence is released ahead of
acceptance (marked as such); this is not a mirror of the research tree.

| Paper | Venue | Status |
|-------|-------|--------|
| [Quine: Realizing LLM Agents as Native POSIX Processes](quine-arxiv/) | arXiv:2603.18030 | published |
| [Coordination Under Existential Unawareness](minimal-perceptual-prerequisites/) | ALIFE 2026 | accepted |
| [From Simulated Worlds to Process Habitats](agentic-alife-workshop/) | ALIFE 2026 workshop | accepted |
| [Facultative Self-Reproduction in Quine](facultative-self-reproduction/) | ALIFE 2026 LBA | accepted |
| [Structural Elicitation in a Frozen POSIX Agent](structural-elicitation/) | ALIFE 2026 LBA | accepted |
| [A Tumor in the Repository](computational-neoplasm/) | ALIFE 2026 LBA | accepted |
| [Structure Grows](structure-grows/) | Blog essay (Substack) | published |
| [The Autopoietic Repository](the-autopoietic-repository/) | Blog essay (Substack) | published |
| [Harness Engineering Is Architectural Amnesia](harness-engineering-is-architectural-amnesia/) | Blog essay (Substack) | published |
| [Why Terminal Bench 2 Is Broken, and Why I Still Love It](terminal-bench-2-love-letter/) | Blog essay (Substack) | published |

---

### [Quine: Realizing LLM Agents as Native POSIX Processes](quine-arxiv/)

**Systems / runtime** · arXiv:2603.18030 · **published**

Realizes an LLM agent as a native POSIX process — identity is the PID, interface
the standard streams, lifecycle `fork`/`exec`/`exit` — so isolation, composition,
and resource control are inherited from the kernel rather than rebuilt in an app
framework.

[Read](quine-arxiv/) · [PDF](quine-arxiv/data/2603.18030v1.pdf) · [Manuscript](quine-arxiv/sections/)

### [Coordination Under Existential Unawareness](minimal-perceptual-prerequisites/)

**Experimental** · ALIFE 2026 · **accepted**

Two LLM agents coordinate under shared scarcity with no hint that a peer exists —
and a threshold separates coordination that *closes* from coordination that
*stalls* (quantitative signals close it; qualitative ones do not).

[Read](minimal-perceptual-prerequisites/) · [PDF](minimal-perceptual-prerequisites/output/Coordination_Under_Existential_Unawareness_Information_Ablation_and_Closure_Thresholds_in_LLM_Multi-Agent_Systems.pdf) · [Experiments](minimal-perceptual-prerequisites/experiments/) · [Reproduce](minimal-perceptual-prerequisites/REPRODUCE.md)

### [From Simulated Worlds to Process Habitats](agentic-alife-workshop/)

**ALife** · ALIFE 2026 workshop · **accepted**

Quine as a process-habitat substrate for generative ALife: withholding
affordances for social knowledge and reproduction reveals heterogeneous
coordination protocols and a successor morphospace.

[Read](agentic-alife-workshop/) · [PDF](agentic-alife-workshop/output/From_Simulated_Worlds_to_Process_Habitats_POSIX_Generative_ALife.pdf) · [Experiments](agentic-alife-workshop/experiments/)

### [Facultative Self-Reproduction in Quine](facultative-self-reproduction/)

**ALife** · ALIFE 2026 LBA · **accepted**

Self-reproduction as a *facultative* outcome across real POSIX `exec` boundaries —
the same carried source becomes a rebuilt runtime or a degenerate stream body
depending only on the demand it faces.

[Read](facultative-self-reproduction/) · [PDF](facultative-self-reproduction/output/Facultative_Self_Reproduction_in_Quine_LBA.pdf) · [Experiments](facultative-self-reproduction/experiments/)

### [Structural Elicitation in a Frozen POSIX Agent](structural-elicitation/)

**ALife** · ALIFE 2026 LBA · **accepted**

A frozen, missionless LLM process, activated with a shell and no task or
reward: the static structure of its workspace directs its functional behavior.
Across five motifs the structure arm produces the matched, externally
verified filesystem/exec event in 96 of 102 runs versus 0 of 85 matched
controls, across three model families, with the pull replicating on non-code
prose and persisting at a read-only substrate boundary.

[Read](structural-elicitation/) · [PDF](structural-elicitation/output/Structural_Elicitation_in_a_Frozen_POSIX_Agent_LBA.pdf) · [Experiments](structural-elicitation/experiments/)

### [A Tumor in the Repository](computational-neoplasm/)

**Experimental** · ALIFE 2026 LBA · **accepted**

A governance structure in an instruction-file-governed repository behaves like a
neoplasm: it captures agent effort and resists its own deletion — both decoupled
from function (defended as readily when fabricated as when real), and
reference-driven recurrence restores it unless authority is rewritten
systemically. Ships the incident-point contract, constitution, and validator
rather than the full run tapes.

[Read](computational-neoplasm/) · [PDF](computational-neoplasm/output/A_Tumor_in_the_Repository_LBA.pdf) · [Artifacts](computational-neoplasm/artifacts/)

### [Structure Grows](structure-grows/)

**Blog essay** · Substack · **published**

The sequel to *The Autopoietic Repository*: how a living repository grows the
organs it needs to maintain itself, and why encoding the strange loop into a file
made the system flexible rather than alive — a different thing than the author
first claimed.

[Read on Substack](https://kehao95.substack.com/p/structure-grows-how-a-living-repository) · [Manuscript](structure-grows/manuscript.md)

### [The Autopoietic Repository](the-autopoietic-repository/)

**Blog essay** · Substack · **published**

The flagship essay: when AI agents become first-class collaborators, a repository
stops being a filing cabinet and becomes an autopoietic system that maintains and
rewrites its own rules.

[Read on Substack](https://kehao95.substack.com/p/the-autopoietic-repository) · [Manuscript](the-autopoietic-repository/manuscript.md)

### [Harness Engineering Is Architectural Amnesia](harness-engineering-is-architectural-amnesia/)

**Blog essay** · Substack · **published**

Agent harnesses keep rebuilding the operating system in user space — re-deriving
isolation, composition, and lifecycle that the kernel already provides. What
changes when you inherit those primitives instead.

[Read on Substack](https://kehao95.substack.com/p/harness-engineering-is-architectural) · [Manuscript](harness-engineering-is-architectural-amnesia/manuscript.md)

### [Why Terminal Bench 2 Is Broken, and Why I Still Love It](terminal-bench-2-love-letter/)

**Blog essay** · Substack · **published**

A benchmark critique with credit: the task contracts are distorted enough that the
leaderboard doesn't mean what it looks like — and the benchmark is still worth
building on.

[Read on Substack](https://kehao95.substack.com/p/why-terminal-bench-2-is-broken-and) · [Manuscript](terminal-bench-2-love-letter/manuscript.md)

---

Experimental entries ship their run tapes as [DVC](https://dvc.org) pointers;
fetch the published blobs with `./scripts/pull-public-dvc-manifest.sh` from the
repository root.
