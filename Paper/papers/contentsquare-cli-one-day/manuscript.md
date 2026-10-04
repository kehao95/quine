# How I Turned ContentSquare into a CLI in One Day Without Writing One Single Line of Code

*Public draft: internal repository names, prototype names, project IDs, private paths, and private URLs have been anonymized.*

---

**Six hours. Zero lines of code written by hand. One complete CLI that an AI agent can pick up and use to do real analytical work.**

This is not a story about fast typing or clever prompts. It is a story about a shift in what "building software" means when the builder is an AI and the human becomes something else entirely--not a coder, but a *gradient-setter*.

---

## What Exists Now

By the end of the sprint, the prototype had become a real CLI with enough surface area to authenticate, orient within a workspace, inspect existing analytical artifacts, and retrieve evidence for non-trivial investigations. A Node.js POSIX binary built on a standard command framework.

The more interesting fact: an AI agent that had never seen this tool before could pick it up, run `--help`, and conduct a multi-step analytical investigation--composing across command groups, repairing mismatches mid-flight, and producing a customer-facing readout that preserved explicit uncertainty. That is an agent using a tool to do real work.

The build arc ran from **Friday, March 27, 2026** to **Saturday, March 28, 2026**. Under 6 hours of active human involvement across three work sessions. I did not write the implementation code, the tests, or the scaffolding. I imported one file by hand: `AGENTS.md`. Everything else was directed through prompts.

---

## A Real Question, Not a Demo

This project started as a product question: **should ContentSquare have a CLI?**

Customers use ContentSquare through a GUI. AI agents can use APIs, and increasingly MCP servers. So what is a POSIX CLI actually for?

The bet: **POSIX is the universal agentic surface.** A CLI composes with pipes, scripts, and shell automation in ways a GUI cannot, and in ways MCP does not yet make universally smooth across tool boundaries. When an agent needs to explore, inspect, retrieve, compare, and export, a well-structured CLI with good help output is one of the most discoverable interfaces available. It is the `gh` move, but for a product surface instead of a source-control platform.

### The one file I placed by hand

Before the repo had a single product command, it had a self-model. I imported `AGENTS.md` from a methodology I had been developing called the [*autopoietic repository*](https://kehao95.substack.com/p/the-autopoietic-repository)--a framework where a repository is not a warehouse of artifacts, but a living system that maintains and modifies itself.

Think of `AGENTS.md` as the repository's DNA. Before the AI agent writes a single line of code, it absorbs this file and learns: how to think about problems, where to store findings, when to treat friction as a signal to repair rather than tolerate. Once injected, the agent stops being a stateless instruction-follower and becomes something closer to an apprentice with institutional memory.

The details are in the linked article. What matters here is the effect: once `AGENTS.md` was in place at commit zero, the agent knew how to preserve research, repair friction, and maintain coherence without being re-taught in every session.

---

## The Build: Three Scenes from Friday to Saturday

The build window contained 8 commits across roughly 29 hours, but the meaningful rhythm was three work sessions:

| Session | Window | What changed |
|---------|--------|--------------|
| 1 | Fri 11:28-12:16 | Bootstrap the MVP and the first live test harness |
| 2 | Fri 22:49-Sat 00:14 | Expand usable capability and strengthen applied analysis |
| 3 | Sat 13:05-16:08 | Push higher-order validation, link-seeded workflows, and auth |

There were two long gaps between those sessions. The "one day" claim is wall-clock compression, not a claim of one uninterrupted 24-hour sprint.

### Scene A: discovering the latent CLI

The agent started by tracing the real product surface: frontend code for user-facing contracts, backend services for data shapes and auth flows, docs for customer-facing concepts, infrastructure for routing and environment boundaries. From that trace, it derived the MVP command groups by mapping what the product actually exposed.

The auth investigation is a good example. The agent traced OAuth-related flows, redirects, grant assumptions, and token-handling constraints, and discovered that the production surface did not expose a clean CLI-native OAuth path. Rather than stall, it documented the finding, identified a session-token workaround, and built around that reality. That research landed in the repo as checked-in memory--the absorption rule doing real work.

> **The CLI was discovered latent in the product graph before it was implemented in code.** Define the target before generating the code. The gradient starts with knowing what to build.

### Scene B: the human role was direction, not construction

After importing `AGENTS.md`, I did not touch the repository by hand. The manifest, README, research notes, implementation, tests, behavior harness, and command-surface repairs were all built by the agent through prompt-directed iteration.

My role was direction: what to build, what constraints to respect, what layer to work on next, and when something had missed the mark. I was setting the gradient, not constructing the artifact.

> **In an agentic sprint, the human job is direction and judgment--not construction.** If the environment is shaped correctly, the agent builds. You are no longer a bricklayer. You are the architect who defines where the building should stand and how tall it should be.

### Scene C: the testing ladder was the gradient

Immediately after the MVP, the repo established a layered validation ladder:

| Layer | Name | What it tests |
|-------|------|---------------|
| L1 | Substrate | Does the code compile? Do types check? |
| L2 | Runtime | Do commands execute without crashing? Do they return expected shapes? |
| L3 | Instructional | Can an agent discover how to use the CLI from help output alone? |
| L4 | Emergent | Can an agent compose multiple commands into a coherent workflow? |
| L5 | Applied | Can an agent use the CLI to do real analytical work? |
| L6 | Operationalization | Can an agent create analysis surfaces on demand for novel problems? |

I defined the layer taxonomy. The builder agent generated the concrete tests, scenarios, and behavior harness.

Two agents in this story:

- The **builder agent** constructs the repo, implements the CLI, writes the tests, and repairs defects.
- The **test agent** runs inside each higher-order behavior scenario and uses the CLI as a tool for the first time.

The ladder gave the builder agent a self-corrective loop:

1. I set a high-level L5 requirement.
2. The builder agent implements the necessary commands and writes a scenario.
3. The builder agent runs the scenario.
4. The test agent hits friction: missing flags, confusing naming, awkward argument shapes, or absent supporting primitives.
5. The builder agent treats that friction as an implementation defect, repairs the CLI, and reruns.
6. Repeat until the scenario passes at the target layer.

This loop sustained itself for long stretches. In some iterations, a single loop ran for more than 90 minutes without human intervention.

> **Once the structural gradient is right--self-model in place, target defined, evidence ladder established--the builder agent iterates autonomously for long stretches.** The human sets the direction. The agent does the climbing.

Think of training a smart retriever. You do not teach it which leg to move first. You throw the ball, set up the guardrails, and let it figure out the path. Tests are not the quality inspector at the end of the assembly line--they are the *carrots* that guide the AI upward through the levels.

> **Tests are not a gate at the end. They are the gradient that makes autonomous iteration productive.**

---

## The Payoff: What "Done" Looked Like

The proof was the preserved canonical L5 runs: a separate agent using the tool for real analytical work under live conditions. Three representative scenarios:

### Performance readout

The test agent was asked to assess performance for a live production project. It started with auth and workspace orientation, inspected existing artifacts, pulled comparative evidence, hit a naming mismatch, repaired mid-run, widened the investigation when the first pass was insufficient, and ended with a readout that preserved explicit uncertainty.

**28 CLI invocations**, **7 command groups**, **16 artifact files** preserved.

This reads like product analysis, not command output. The agent was not following a script--it explored, composed, repaired, and synthesized because the CLI was discoverable enough to support that behavior.

### Diagnostic investigation

The test agent was asked to diagnose conversion issues. It swept across multiple evidence surfaces in parallel, drilled into representative artifacts, attempted natural-but-missing affordances, then recovered by falling back to working primitives and adjacent help output.

**25 invocations**, **15 command groups**, **30 artifact files**.

This showed the CLI functioning as a diagnostic workspace, not a reporting endpoint--and exposed latent demand: a diagnostic-minded agent naturally reached for surfaces not yet implemented.

### Artifact sufficiency check

The test agent was given a pasted ContentSquare URL and asked whether the linked artifact alone was sufficient for a decision. It inspected the artifact directly from the URL, widened into adjacent evidence surfaces, and assessed whether the investigation needed to broaden.

**38 invocations**, **8 command groups**, **9 artifact files**.

This is how humans and agents actually work: start from a pasted link, not a memorized identifier. The CLI handled URL normalization transparently and supported widening from one artifact into a broader investigation--both affordances that came from earlier friction-repair loops.

### What these runs proved

All three demonstrated the same thing: a CLI coherent enough that a separate agent could compose across surfaces, recover from friction, preserve uncertainty, and do real analytical work. The richness came from the quality of the tool, not from hidden cleverness in the test harness.

---

## The Idea: Build the Gradient, Then Let It Grow

The method compressed into a developer recipe:

1. **Give the repo a self-model first.** Import your operating philosophy before the first feature commit. This is the DNA injection.
2. **Define the target before generating code.** Trace the real product surface through code, contracts, docs, and routing--not through imagination.
3. **Let research persist as repo memory.** If findings only exist in chat, they are lost. The moment a session ends, so does anything not absorbed into files.
4. **Build the evidence ladder immediately.** From "does it compile?" to "can an agent do real work with it?" The ladder is the gradient.
5. **Let higher-order targets pull implementation upward.** Stop asking "what feature next?" and ask "what blocks the next level of usefulness?"

> **Without the gradient, an agent produces drift. With the gradient, it climbs.**

And when the gradient is well-formed, the climb sustains itself for surprisingly long stretches--90-minute autonomous loops where the agent builds, tests, hits friction, repairs, and climbs again, without the human touching a single file.

---

## Limits

**The CLI was still incomplete.** Auth still relied on a session-token workaround. Fast did not mean finished.

**The validation was still read-only.** The preserved runs demonstrated rich analytical behavior on existing artifacts. They did not validate write-path workflows end to end.

**L6 was defined but not yet testable.** Having an agent create new analysis surfaces on demand remained blocked by environment availability and safe write access, not by conceptual direction.

**Read-only usefulness is a plateau.** The method worked through L5. Whether it scales into creative, write-heavy operationalization is an open question.

---

The structural gradient made the sprint productive. It did not make the product finished. But it made the next steps visible.

> **In the age of AI-assisted development, the human's job is no longer to build--it is to shape the environment in which building becomes inevitable.**

What structural gradient would you build for *your* next project?
