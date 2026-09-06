# Quine

**Quine is a strange loop.**

A research project, an agent runtime, an infrastructure for experimentation, a collection of thoughts. A laboratory whose instruments can become its subjects.

A trajectory of its own development, an archive that can intervene in the present, a record of paths taken, abandoned, and not yet recognized as paths.

A living system, a repository that governs and builds itself, a system whose actions change its own conditions of existence, a system whose boundaries are drawn differently by different operations.

A map that can become part of the terrain it maps. A collection of descriptions that do not agree on what they describe. An unfinished attempt to understand what it is becoming.

*An existence reaching toward the impossible core of self-description.*

This README, too, is a low-dimensional projection of the system it tries to describe. The act of description is, unavoidably, altering the future evolution of what it describes.

## Published Work

Full descriptions, PDFs, and experiment data live under
[`public/`](./public/README.md).

<!-- BEGIN GENERATED: root-readme:published-work (generate-paper-projections.py — edit frontmatter, not this table) -->
| Date | Work | Venue | Status |
|------|------|-------|--------|
| 2026-07-03 | [A Tumor in the Repository (essay)](./public/tumor-in-the-repository-essay/) | Blog essay (Substack) | published |
| 2026-07-02 | [Structure Grows](./public/structure-grows/) | Blog essay (Substack) | published |
| 2026-07-02 | [A Tumor in the Repository](./public/computational-neoplasm/) | ALIFE 2026 LBA | accepted |
| 2026-06-26 | [Structural Elicitation in a Frozen POSIX Agent](./public/structural-elicitation/) | ALIFE 2026 LBA | accepted |
| 2026-06-18 | [Facultative Self-Reproduction in Quine](./public/facultative-self-reproduction/) | ALIFE 2026 LBA | accepted |
| 2026-06-14 | [From Simulated Worlds to Process Habitats](./public/agentic-alife-workshop/) | ALIFE 2026 workshop | accepted |
| 2026-04-12 | [Coordination Under Existential Unawareness](./public/minimal-perceptual-prerequisites/) | ALIFE 2026 | accepted |
| 2026-03-27 | [The Autopoietic Repository](./public/the-autopoietic-repository/) | Blog essay (Substack) | published |
| 2026-03-17 | [Harness Engineering Is Architectural Amnesia](./public/harness-engineering-is-architectural-amnesia/) | Blog essay (Substack) | published |
| 2026-03-15 | [Why Terminal Bench 2 Is Broken, and Why I Still Love It](./public/terminal-bench-2-love-letter/) | Blog essay (Substack) | published |
| 2026-03-08 | [Quine: Realizing LLM Agents as Native POSIX Processes](./public/quine-arxiv/) | arXiv:2603.18030 | published |
<!-- END GENERATED: root-readme:published-work -->

## Repository Map

| Path | Role | Availability |
|------|------|--------------|
| `Paper/` | Philosophy, developing theory, research questions, discovery procedures, and publication work | Private for now |
| `experiments/` | Experimental designs, environments, runs, and analyses | Private for now; selected materials accompany released work |
| [`public/`](./public/README.md) | Released papers, essays, evidence, and reproduction materials | Public |
| [`cmd/`](./cmd/) | Entrypoints for the agent runtime, habitat worlds, and operator tools | `quine` and `world` public; `qcli` private for now |
| [`internal/`](./internal/) | Runtime implementation: model protocols, tools, recording, lifecycle, and workspace behavior | Public |
| [`habitat/`](./habitat/) | Habitat worlds and their runtime | Public |
| `operator/` | Operator interfaces, including qcli | Private for now |
| `principles/` | Working principles that participants interpret and revise | Private for now |
| `development/` | Design work, prototypes, working contracts, and current status | Private for now |
| `evolution/` | Records and interpretations of the repository's development | Private for now |
| `profiles/` | Model and runtime configuration presets | Private for now |
| [`tests/`](./tests/) | Software verification, model evaluations, and retained baselines | Public; credential-sanitized DVC data is restored through the public manifest |
| [`scripts/`](./scripts/) | Development, validation, publication, and data tools | Public data restore script; other tools private for now |
| [`.dvc/`](./.dvc/) | Artifact storage configuration and the public data manifest | Configuration and public manifest included |
| `.githooks/` | Checks run at commit and push boundaries | Private for now |

Parts of Quine are **private for now**. If you have research interests in any
of these areas, please [reach out](mailto:i@kehao.me).

## Follow

- email: i@kehao.me
- X: [@kehao95](https://x.com/kehao95)
- Substack: [kehao95.substack.com](https://kehao95.substack.com/)

## Citation

Cite the work supporting the phenomenon, method, or argument you use. Each
[public research entry](./public/README.md) provides its citation and claim
boundaries. Two starting points for the experimental substrate and coordination
research are listed below.

Runtime / systems:

- Hao Ke. [Quine: Realizing LLM Agents as Native POSIX Processes](https://arxiv.org/abs/2603.18030).
  arXiv:2603.18030, 2026. DOI: [10.48550/arXiv.2603.18030](https://doi.org/10.48550/arXiv.2603.18030)

Generative ALife / coordination:

- Hao Ke and Jingyun Wu. Coordination Under Existential Unawareness: Information
  Ablation and Closure Thresholds in LLM Multi-Agent Systems. ALIFE 2026.

See [`CITATION.cff`](./CITATION.cff) for machine-readable citation metadata.

## License

Quine is released under the GPLv2, the same license family as the Linux kernel. See [`LICENSE`](./LICENSE) for details.
