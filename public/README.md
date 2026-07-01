# Published works

The public portfolio — one self-contained entry per published or accepted work.
Each entry carries everything needed to read, reproduce, and cite it: manuscript
source, figures, the released PDF, and (for experimental papers) the experiment
data and a reproduction protocol. Inclusion is strict — only **published** or
**accepted** works appear here; this is not a mirror of the research tree.

| Paper | Venue | Status |
|-------|-------|--------|
| [Quine: Realizing LLM Agents as Native POSIX Processes](quine-arxiv/) | arXiv:2603.18030 | published |
| [Coordination Under Existential Unawareness](minimal-perceptual-prerequisites/) | ALIFE 2026 | accepted |
| [From Simulated Worlds to Process Habitats](agentic-alife-workshop/) | ALIFE 2026 workshop | accepted |
| [Facultative Self-Reproduction in Quine](facultative-self-reproduction/) | ALIFE 2026 LBA | accepted |

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

---

Experimental entries ship their run tapes as [DVC](https://dvc.org) pointers;
fetch the published blobs with `./scripts/pull-public-dvc-manifest.sh` from the
repository root.
