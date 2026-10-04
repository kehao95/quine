# Agentic ALife Workshop Workspace

This file is the active work surface for the ALIFE 2026 `Agentic and
Generative AI as Artificial Life` workshop submission.

## Current State

- Live status: [`development/status/README.md`](../../../development/status/README.md)
- Deadline: `2026-06-14` (extended from `2026-06-01`)
- Status: OS-substrate pivot applied; assay-space framing selected; manuscript
  regenerated, assay sections softened for academic flow, workshop claim
  calibration absorbed, compressed empirical texture added, submission
  mechanics confirmed, and `4`-page PDF verified on 2026-05-26; a 2026-06-12
  language-flow pass tightened the introduction, habitat, and discussion prose,
  normalized the bibliography to a consistent `et al.` author style, and
  re-verified the `4`-page PDF; a 2026-06-13 claim-correction pass replaced the
  "body-transcending succession" over-read with *succession by material
  recruitment* (re-anchored on the genuine no-source `6M.03`/`6M.05` evidence,
  including the new `6M.05 C11` native-control closure), and re-verified `4`
  pages; a 2026-06-13 reviewer-absorption / flow pass then added an
  emergence-vs-retrieval boundary, a `morphospace`-coverage qualifier, a
  softened MAS characterization, and a pre-OEE concession, then trimmed
  redundant prose and two framework citations to hold `4` pages
- Priority shift: as of `2026-05-21`, this is the active near-term runtime-
  ontology paper line and should take precedence over further scope growth
- Current deliverable: `4`-page workshop short paper, including references
- Compression fallback: `2`-page workshop abstract, or `1` page if required
- Manuscript source: [`manuscript.md`](manuscript.md)
- Build metadata: [`metadata.yaml`](metadata.yaml)
- Bibliography: [`references.bib`](references.bib)
- Generated TeX target: `output/paper.tex`
- Generated PDF target:
  `output/From_Simulated_Worlds_to_Process_Habitats_POSIX_Generative_ALife.pdf`
- Literature synthesis: [`references/deep-research-synthesis.md`](references/deep-research-synthesis.md)
- Supporting P7G evidence synthesis:
  [`references/p7g-trace-uptake-evidence.md`](references/p7g-trace-uptake-evidence.md)
- Build note: `make tex` and `make pdf` both succeed from this paper
  directory. `tectonic` was installed locally through Cargo on 2026-05-18.
  The current PDF is `4` pages and carries expected TeX underfull warnings from
  dense evidence tables.

## Submission Mechanics

Confirmed from the workshop website on `2026-05-26`; deadline extension
verified on `2026-06-01`:

- submission deadline: `2026-06-14` (extended from `2026-06-01`)
- format: one-page extended abstract or short paper up to `4` pages including
  references
- file: single PDF prepared with the ALIFE 2026 template
- route: email the PDF to `stefano.nichele@hiof.no`
- presentation form: accepted contributions may become short contributed talks
  or posters
- publication surface: accepted abstracts / short papers will be compiled into
  a public booklet PDF on the workshop website
- attendance: main ALIFE 2026 registration is required, and at least one author
  must attend and present in person

## Working Rule

Draft from Quine's runtime object, not from Quine Lab's repo maintenance model.

Because this workshop contribution is not expected to enter the main ALIFE
2026 proceedings, assay selection should optimize for what the substrate makes
distinguishable rather than for avoiding overlap with submitted main-conference
papers. The current spine is Ariadne / `2.6` for death-memory-lineage outcome
classes, `8.17` for environmental sociality, P7G `C16/C17` for trace topology
and repair basins, and `6M.02` / `6M.07` for successor morphospace and
governance. The 2026-05-26 `8.20-pure` finite-time probe is supplementary
boundary evidence for Assay 2 only: it shows delayed artifact-aware coordination
without `fs_mutations` or budget visibility, but no writeback / validation
closure before operator cutoff.

The manuscript should read as a strong position plus exploratory assay
proposal, not as a full empirical paper. Keep the runtime model explicit
(`P_t`, `C_t`, `H_t`, carriers, host intervention, successor), treat the assay
labels as working distinctions, and avoid wording that implies the current
observations prove life or final taxonomy.

The paper should stay centered on POSIX-native agent physics and should lead
with the OS-substrate question for ALife:

- agent equals process
- POSIX is the environment's physics
- context is volatile working memory
- `fork`, `exec`, and `exit` are lifecycle operators
- public process/control surfaces make governance and successor viability
  observable
- generated files, markers, repair ledgers, and carriers become ecological
  structure
- traditional ALife worlds, generative-agent sandboxes, and MAS frameworks are
  comparators for where substrate semantics are supplied, not targets to
  dismiss

## Near-Term Milestones

| Date | Target |
|------|--------|
| 2026-05-15 | Create paper surface |
| 2026-05-20 | Freeze one-sentence thesis, title, and `4`-page format choice |
| 2026-05-21 | Reconfirm original `2026-06-01` deadline and treat the agentic ALIFE workshop paper as active near-term paper line |
| 2026-05-23 | Draft first full short-paper manuscript |
| 2026-05-25 | Canonicalize OS-substrate / traditional-ALife pivot |
| 2026-05-25 | Record assay selection rule and regenerate manuscript spine |
| 2026-05-25 | Apply public-facing style/flow pass |
| 2026-05-25 | Soften Assay sections and conclusion against mechanical report style |
| 2026-05-26 | Absorb workshop claim calibration and minimal runtime model |
| 2026-05-26 | Record finite-time `8.20-pure` sidechannel-removal probe |
| 2026-05-26 | Confirm submission mechanics and polish abstract / citations |
| 2026-05-26 | Review compressed ALIFE template PDF against page budget |
| 2026-05-28 | Produce submission PDF candidate |
| 2026-05-30 | Final internal review / trim |
| 2026-06-01 | Record workshop deadline extension to `2026-06-14` |
| 2026-06-12 | Absorb language-flow refinement pass; normalize bibliography to consistent `et al.` style; re-verify `4`-page PDF |
| 2026-06-13 | Correct successor crown jewel (material recruitment) off the no-source over-read; re-anchor on `6M.03`/`6M.05 C11`; re-verify `4` pages |
| 2026-06-13 | Absorb reviewer feedback (emergence-vs-retrieval, reproducibility, `morphospace`, MAS, pre-OEE); smooth flow; trim to hold `4` pages |
| 2026-06-14 | Email submission PDF to workshop organizer |

## Drafting Order

1. Start from [`manuscript-plan.md`](manuscript-plan.md).
2. Draft in [`manuscript.md`](manuscript.md).
3. Compile TeX with `make tex` from this paper directory.
4. Keep the workshop audience first: agentic AI ecosystems, generative
   complexity, safety/security/governance.
5. State the traditional-ALife / simulated-world contrast early.
6. Keep the ordinary-MAS contrast as a compact bridge to the agentic-AI audience.
7. Use the assay-space portfolio in [`manuscript-plan.md`](manuscript-plan.md);
   do not drift back to the earlier `2.3`-only ecology example unless a
   specific compression decision requires it.
8. Keep public manuscript prose free of repo paths, run IDs, evidence-set
   handles, and local experiment numbers; those belong in dossier provenance.
9. Use Quine examples only to illustrate the runtime-substrate distinction.
10. Use `illustrate`, `expose`, and `make distinguishable` for exploratory
   assays rather than definitive proof language.
11. End with what ALife gains: a process habitat for turning agentic-AI
   ecology and governance questions into runtime assays.

## Exit Criteria

Ready means:

- title and thesis fit the workshop without overclaiming life or autopoiesis
- abstract clearly distinguishes Quine from simulated ALife worlds and ordinary
  multi-agent systems
- submission format remains chosen or is explicitly downgraded to `2` or `1`
  pages
- PDF uses ALIFE 2026 template
