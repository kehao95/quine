# Deep Research Query Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


## Fact-Finding vs Generative Theory Sparks

- **`dr` (Gemini Deep Research CLI / Web Search):** External literature searches, prior-art mapping, empirical citations, and claim verification.
- **`agy` (Gemini CLI):** Lateral theory sparks, divergent analogies, and cross-paradigm seeds (owned by [`theory-spark-query.md`](./theory-spark-query.md)).

## Query Modes

Every deep research query operates in one of two modes:

1. **Gap Survey (Field-First):**
   - *Question:* "What does the field leave unsolved, unasked, or assumed?"
   - *Guideline:* Search on the field's own terms, independent of current repo evidence. Target unstated assumptions and limitations rather than filtering only for what current data can close.
2. **Novelty Check (Claim-First):**
   - *Question:* "Has anyone already articulated claim X?"
   - *Dual Extraction:* When finding nearest-neighbor or pre-empting papers, always extract their stated limitations and future-work items—convert defense into offense in a single pass.

## Proportional Brief Formats

Tailor the structure of the research query to the scope of the investigation:

### A. Quick Exploratory Queries
For bounded questions, quick prior-art checks, or sanity lookups, avoid boilerplate overhead. Provide:
1. **Core Question:** The precise question to answer.
2. **Context & Motivation:** What decision or text this informs.
3. **Expected Deliverable:** Desired verdict (`open` / `pre-empted` / `differentiating nuances`) and key citations.

### B. Formal / Manuscript Research Briefs
For extensive field sweeps or major paper positioning:
- **Context & Core Research Question**
- **Detailed Sub-questions & Suggested Search Anchors**
- **Relevance & Scope Filters** (in-scope vs out-of-scope)
- **Synthesis Deliverables & Red Lines** (key distinctions to extract, non-claims to preserve)

## Operational Guidelines

- **Preserve Session IDs:** When running long-running async research jobs (`dr`), always capture and record the interaction ID for retrieval.
- **Verification Budget:** Avoid bundling multiple disjoint topics into a single verification run; focus searches on coherent question clusters.
