# Citation verification ledger

Status: **Verification pass completed 2026-07-05** (WebSearch + WebFetch against
arxiv.org / doi.org / publisher / primary pages). All 7 items resolved below.
**Second pass 2026-07-06** (items 8–13): Szeider full-text methods verified
(upgrades item 5's caveat), plus five new related-work citations and the
stigmergy pair. Summary line per item: VERIFIED / VERIFIED-WITH-CORRECTION /
NOT-FOUND.

---

## 1. arXiv 2505.11584 — VERIFIED (with title-drift note)

- **Exact title:** *LLM Agents Are Hypersensitive to Nudges*
- **Authors:** Manuel Cherep, Pattie Maes, Nikhil Singh
- **Submitted:** 16 May 2025. **Comments:** 33 pages, 28 figures. **Category:** cs.AI
- **URL:** https://arxiv.org/abs/2505.11584
- **Our attribution:** frozen LLM agents respond to environment-embedded nudges;
  framed as vulnerability/sensitivity (not positive design); nudges are
  message-like/human-derived; agents have tasks. — **Accurate.**
- **Verbatim quotes:**
  - "We perform a case study with a few such LLM models on a multi-attribute
    tabular decision-making problem, under canonical nudges such as the default
    option, suggestions, and information highlighting, as well as additional
    prompting strategies."
  - "They show much higher susceptibility to the nudges... All these findings
    suggest that behavioral tests are needed before deploying models as agents or
    assistants acting on behalf of users in complex environments."
- **Note:** The nudges (default option, suggestions, information highlighting)
  are choice-architecture manipulations of a task decision — consistent with our
  "message-like/human-derived, agents have tasks" framing. Vulnerability register
  confirmed ("behavioral tests are needed before deploying"). **Cross-ref:** this
  is the same body of work as PNAS item 3 (see below) — the GitHub repo
  PapayaResearch/nudging titles it "AI Agents are Sensitive to Nudges" and links
  arXiv 2505.11584.

## 2. arXiv 2509.25609 — VERIFIED-WITH-CORRECTION (title)

- **Exact title:** *A Framework for Studying AI Agent Behavior: Evidence from
  Consumer Choice Experiments* (the system/benchmark is named **ABxLab**; "ABxLab"
  is not the paper title)
- **Authors:** Manuel Cherep, Chengtian Ma, Abigail Xu, Maya Shaked, Pattie Maes,
  Nikhil Singh
- **Submitted:** 30 Sep 2025 (revised 24 Feb 2026). **Comments:** ICLR, 31 pages,
  17 figures.
- **URL:** https://arxiv.org/abs/2509.25609
- **Our attribution:** benchmark showing LLM shopping/web agents are steered by
  nudges (position, framing etc.) in web environments; vulnerability register;
  task present. — **Accurate in substance.**
- **Correction:** cite as the full title above with "(ABxLab)" as the framework
  name, not as if "ABxLab" were the title.
- **Verbatim quotes:**
  - "We introduce ABxLab, a framework for systematically probing agentic choice
    through controlled manipulations of option attributes and persuasive cues. We
    apply this to a realistic web-based shopping environment, where we vary prices,
    ratings, and psychological nudges..."
  - "We find that agent decisions shift predictably and substantially in
    response, revealing that agents are strongly biased choosers... This
    susceptibility reveals both risk and opportunity..."
- **Note:** Agents are tasked with consumer choices in a web shopping env; nudges
  explicit ("psychological nudges"). Register is dual (risk AND opportunity), not
  purely vulnerability — if our text implies purely vulnerability framing, soften
  to "risk-and-opportunity."

## 3. PNAS 10.1073/pnas.2537030123 — VERIFIED-WITH-CORRECTION (title + same-work note)

- **Exact title:** *AI agents are sensitive to nudges*
- **Authors:** Manuel Cherep, Pattie Maes, Nikhil Singh
- **Journal/date:** Proceedings of the National Academy of Sciences (PNAS),
  vol. 123, no. 25, published 15 June 2026.
- **URL:** https://www.pnas.org/doi/10.1073/pnas.2537030123 (DOI resolves;
  publisher page returned 403 to automated fetch but DOI + metadata confirmed via
  PNAS listing, MIT Media Lab page, and PapayaResearch/nudging repo)
- **Our attribution:** a 2025/2026 PNAS paper showing LLM agents sensitive to
  nudges/choice architecture. — **Accurate.**
- **Correction / IMPORTANT:** This PNAS paper is the **peer-reviewed version of
  the SAME work as item 1 (arXiv 2505.11584)** — same authors, same nudge
  findings, same PapayaResearch/nudging codebase. Do **not** cite items 1 and 3 as
  two independent pieces of evidence for the nudge-vulnerability cluster; they are
  one result at preprint vs. journal stage. Cite the PNAS version as the canonical
  reference. Correct year is **2026** (not 2025).
- **Verbatim quotes (from significance/abstract as surfaced):**
  - "Many applications of AI agents tacitly assume that, under uncertainty, they
    will react in roughly human-like ways, if not more rationally... the researchers
    show this assumption often fails."
  - "none of these methods resolve the sensitivity of these models to nudges."

## 4. Roger Barker, *Ecological Psychology* (1968) — VERIFIED

- **Exact reference:** Roger G. Barker, *Ecological Psychology: Concepts and
  Methods for Studying the Environment of Human Behavior*, Stanford University
  Press, 1968.
- **Our attribution:** behavior settings are social programs with joined human
  participants — standing patterns of behavior tied to place+time+program; NOT
  bare inanimate structure acting on a solitary non-member agent. — **Accurate.**
- **Sources used (book not online; verified via secondary):**
  - Wikipedia "Behavior settings": https://en.wikipedia.org/wiki/Behavior_settings
  - Royal Society / PMC review (Barker School, virtual contexts):
    https://pmc.ncbi.nlm.nih.gov/articles/PMC11500833/
- **Verbatim quotes:**
  - (Wikipedia, on Barker) "A behavior setting is a self-referenced (internally
    interdependent and self-defined) entity that consists of one or more standing
    patterns of behavior." — and it "exists at the interface between the standing
    patterns of behavior and the milieu."
  - (PMC review) behavior settings are "interdependent standing patterns of
    behaviour and milieu, where the milieu is *circumjacent* and *synomorphic* with
    the behaviour"; characterized as "relatively short-term activity by a relatively
    small group of interacting people in a circumscribed space."
- **Note:** Supports our out-positioning claim: a behavior setting requires a
  *standing pattern of behavior* enacted by *interacting people* (a vacant setting
  is "a meaningless artifact"). This is a social/program unit with joined
  participants — distinct from bare inanimate workspace structure acting on a
  solitary non-member agent. Attribution is safe to make load-bearing.

## 5. Szeider 2025 (@Szeider2025LeftAlone) — VERIFIED-WITH-CORRECTION (exact ref) + scope caveat

- **Exact title:** *What Do LLM Agents Do When Left Alone? Evidence of
  Spontaneous Meta-Cognitive Patterns*
- **Author:** Stefan Szeider (sole author)
- **arXiv:** **2509.21224** (submitted 25 Sep 2025). Also at OpenReview
  id=rzJlQi3rAd.
- **URL:** https://arxiv.org/abs/2509.21224
- **Correction:** one search hit gave a wrong ID (2509.21266); the correct ID is
  **2509.21224**. Ensure the bib entry uses 2509.21224.
- **Our attribution:** reports spontaneous meta-cognitive patterns in unsupervised
  LLM agents. — **Accurate.** The methodological critique we attach (reads behavior
  off narration/self-reports, no matched controls, no functional/filesystem DV) is
  **our critical reading, not stated in the abstract** — the abstract neither
  confirms nor denies it. Flag: if this critique is presented as if sourced from
  the paper, verify against the full-text methods before making it load-bearing;
  as our own analysis it is fine but should be phrased as our assessment.
- **Verbatim quotes:**
  - "We introduce an architecture for studying the behavior of large language
    model (LLM) agents in the absence of externally imposed tasks. Our continuous
    reason and act framework, using persistent memory and self-feedback, enables
    sustained autonomous operation."
  - "We find agents spontaneously organize into three distinct behavioral
    patterns: (1) systematic production of multi-cycle projects, (2) methodological
    self-inquiry into their own cognitive processes, and (3) recursive
    conceptualization of their own nature."

## 6. Ye et al. 2026 (@Ye2026LookBeforeYouLeap) — VERIFIED

- **Exact title:** *Look Before You Leap: Autonomous Exploration for LLM Agents*
- **Authors:** Ziang Ye, Wentao Shi, Yuxin Liu, Yu Wang, Zhengzhou Cai,
  Yaorui Shi, Qi Gu, Xunliang Cai, Fuli Feng
- **arXiv:** 2605.16143 (submitted 15 May 2026)
- **URL:** https://arxiv.org/abs/2605.16143
- **Our attribution:** studies autonomous exploration by LLM agents that remains
  goal-directed (not task-free). — **Accurate.**
- **Verbatim quotes:**
  - "We propose the Explore-then-Act paradigm, which decouples information-gathering
    from task execution: agents first utilize an interaction budget to acquire
    grounded environmental knowledge, then leverage it for task resolution."
  - "Large language model based agents often fail in unfamiliar environments due to
    premature exploitation: a tendency to act on prior knowledge before acquiring
    sufficient environment-specific information."
- **Note:** Exploration is explicitly in service of downstream task resolution
  (reward-driven RL rollouts). Confirms goal-directed, not task-free/reward-free —
  supports the contrast we draw.

## 7. "Environment = information" definitional quotes — VERIFIED-WITH-CORRECTION (LangChain URL)

### 7a. LangChain — VERIFIED-WITH-CORRECTION (source page)
- **Exact quote (verbatim):** "Context engineering is building dynamic systems to
  provide the right information and tools in the right format such that the LLM can
  plausibly accomplish the task."
- **Correct source:** LangChain blog, **"The rise of context engineering"** (by
  Harrison Chase), NOT the "Context Engineering for Agents" page.
- **URL:** https://www.langchain.com/blog/the-rise-of-context-engineering
  (also blog.langchain.com/the-rise-of-context-engineering/ → redirects here)
- **Correction:** our attribution paraphrased "the right information and tools…
  such that the LLM can plausibly accomplish the task" — wording is correct, but
  cite the "rise of context engineering" post as the source. (This phrasing does
  NOT appear on the separate langchain.com/blog/context-engineering-for-agents
  page — confirmed absent there.)

### 7b. Anthropic — VERIFIED
- **Exact quote (verbatim):** context is "the set of tokens included when sampling
  from a large-language model (LLM)." Context engineering is "the set of strategies
  for curating and maintaining the optimal set of tokens (information) during LLM
  inference."
- **Source:** Anthropic engineering blog, "Effective context engineering for AI
  agents."
- **URL:** https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents
- **Our attribution:** context is "the set of tokens… (information)". — **Accurate**
  (the parenthetical "(information)" appears in the "optimal set of tokens
  (information)" phrasing). Safe to cite as-is.

## 8. Szeider 2509.21224 FULL-TEXT methods — VERIFIED 2026-07-06 (closes item 5's caveat)

- **What was checked:** the full text (arXiv HTML) against the four-part
  methodological critique our related-work section makes load-bearing.
- **Findings (each stronger than our original phrasing):**
  1. *Measurement source:* text-only. Agents run "in complete isolation from
     system resources, with all interactions mediated through controlled tool
     interfaces"; the architecture enforces "strict safety constraints that
     prevent external actions beyond observation and communication."
  2. *Controls:* none — no baseline arm, no matched comparison; all 18 runs
     use the same architecture and system prompt across 6 models.
  3. *Classification:* the three behavioral patterns come from qualitative
     reading of agent transcripts; no systematic coding procedure,
     inter-rater reliability, or automated functional check is reported.
  4. *Artifacts:* agents have no filesystem access ("the provided environment
     does not allow the agent to run code"); the only writable surface is a
     key-value memory. No independent artifact inspection is possible.
- **Consequence:** the manuscript's contrast (real workspaces vs isolation;
  matched pairs vs none; functional filesystem/exec DVs vs narrated
  transcripts) is now full-text verified and safe as load-bearing. Item 5's
  "phrase as our assessment" caveat is retired; the manuscript states the
  facts and cites them to the paper's own text.

## 9. Yao et al., ReAct — VERIFIED 2026-07-06

- **Exact title:** *ReAct: Synergizing Reasoning and Acting in Language Models*
- **Authors:** Shunyu Yao, Jeffrey Zhao, Dian Yu, Nan Du, Izhak Shafran,
  Karthik Narasimhan, Yuan Cao. **Venue:** ICLR 2023. **arXiv:** 2210.03629.
- **Our attribution:** the reasoning-acting loop is motivated by task-solving.
  — **Accurate** ("generate both reasoning traces and task-specific actions in
  an interleaved manner… to improve task-solving performance").

## 10. Yang et al., SWE-agent — VERIFIED 2026-07-06

- **Exact title:** *SWE-agent: Agent-Computer Interfaces Enable Automated
  Software Engineering*
- **Authors:** John Yang, Carlos E. Jimenez, Alexander Wettig, Kilian Lieret,
  Shunyu Yao, Karthik Narasimhan, Ofir Press. **Venue:** NeurIPS 2024.
  **arXiv:** 2405.15793.
- **Our attribution:** the agent-computer interface is engineered so agents
  resolve tasks better. — **Accurate** ("SWE-agent's custom agent-computer
  interface (ACI) significantly enhances an agent's ability" on SWE tasks).

## 11. Liu et al., AgentBench — VERIFIED 2026-07-06

- **Exact title:** *AgentBench: Evaluating LLMs as Agents*
- **Authors:** Xiao Liu et al. (22 authors). **Venue:** ICLR 2024.
  **arXiv:** 2308.03688.
- **Our attribution:** evaluation instantiates environments as task-scored
  benchmarks. — **Accurate** ("8 distinct environments to assess
  LLM-as-Agent's reasoning and decision-making abilities").

## 12. UED pair — VERIFIED 2026-07-06

- **Dennis et al. 2020**, *Emergent Complexity and Zero-shot Transfer via
  Unsupervised Environment Design*, NeurIPS 2020, arXiv 2012.02096 —
  introduces UED and PAIRED; curricula parameterized by **regret**
  ("the adversary generates environments which maximize regret"). — matches
  our "curricula parameterized by regret" phrasing.
- **Jiang et al. 2021**, *Prioritized Level Replay*, ICML 2021, arXiv
  2010.03934 — levels sampled by "higher estimated learning potential"
  (TD-error) — matches our "learning potential" phrasing.

## 13. Stigmergy pair — VERIFIED (inherited 2026-07-06)

- **Theraulaz & Bonabeau 1999**, *A Brief History of Stigmergy*, Artificial
  Life 5(2):97–116, doi:10.1162/106454699568700; **Heylighen 2016**,
  *Stigmergy as a Universal Coordination Mechanism I*, Cognitive Systems
  Research 38:4–13, doi:10.1016/j.cogsys.2015.12.002.
- Entries copied verbatim from the already-verified
  `papers/minimal-perceptual-prerequisites/references.bib` (used load-bearing
  in that dossier's ALIFE camera-ready). Our attribution — stigmergy as
  environment-mediated coordination for reactive rule-agents — is the
  standard reading of both.

---

### Corrections roll-up (what to fix before load-bearing use)
1. Item 2: title is *A Framework for Studying AI Agent Behavior: Evidence from
   Consumer Choice Experiments*; "ABxLab" is the framework name, not the title.
   Register is risk-AND-opportunity, not pure vulnerability.
2. Item 3: PNAS title is *AI agents are sensitive to nudges*; year **2026**; it is
   the SAME work as item 1 — do not double-count as independent evidence.
3. Item 5: use arXiv **2509.21224** (not 2509.21266). ~~Methodology critique is
   our own reading, not from the abstract — phrase accordingly.~~ **Retired
   2026-07-06:** full-text verified (item 8); the critique is now stated as
   fact and is safe as load-bearing.
4. Item 7a: cite the LangChain "**The rise of context engineering**" post
   (langchain.com/blog/the-rise-of-context-engineering) as the source of the
   quote.
