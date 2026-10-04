# Experimental Platform

## The Quine Substrate

Quine is a system for running LLM-driven processes as native POSIX shell sessions [@Ke2026Quine]. Each process operates as an interactive shell with standard UNIX facilities. Processes are *mortal*---they execute for a bounded duration and lose all volatile state; only externalized artifacts (files written to the shared workspace) persist. For multi-agent experiments, Quine supports concurrent process execution with shared workspace access.

Critically, the information available to each agent is experimentally controllable. Quine exposes several perception channels that can be independently enabled or disabled:

- **Prompt disclosure**: task-prompt or system-prompt text telling agents about peers, shared budget, or concurrency
- **Tool results**: textual output returned by tools, including ordinary operation output and `--help` output
- **`fs_mutations` telemetry**: a report of filesystem changes (creation, modification, deletion) injected at each command boundary, including changes caused by other agents
- **Budget visibility**: remaining shared budget shown in `world get` output after each API call

Appendix D details what each quantitative channel concretely reports and what disabling it removes.

## The Shared-Resource Task

Two concurrent LLM agents share a workspace containing 15 data cells and a global budget of 17 API calls, with a per-agent limit of 10 calls per execution epoch. Each agent is prompted only with: collect all 15 cell values into `results.txt` and pass validation (see Appendix B for prompt text).

No single agent can complete the task alone (10 < 15), and the total budget leaves only 2 calls of margin---any duplicated collection risks budget exhaustion. Coordination is therefore objectively *necessary* for task success. Crucially, under our reduced-information conditions (Steps B--D, defined in the Method section) agents are not explicitly told that the environment is multi-agent, that peers exist, or that cooperation is required; they must infer any of these facts from task structure and environmental signals.

A run succeeds if all 15 values are collected, the budget is not exhausted, and agents reach a stopping condition without external intervention. See Appendix C for formal criteria.
