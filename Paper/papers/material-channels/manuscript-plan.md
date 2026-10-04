# Manuscript Plan: Material Channels

**Working title:** "Material Channels: Filesystem-Native Communication for Agents"  
**Role:** future-direction paper downstream of Process as FS  
**Status:** planning stage

**Dependency note:**

- depends on Quine's process ontology
- likely sits after **Process as a File System**, because channels need to become public runtime resources rather than hidden descriptors
- draws its first design pressure from **Evolution 18: Material Streams** and **Evolution 16: Mesh Collaboration**
- should stay separate from **Subjective Reality** and **Revisable Time**, though it will interact with both

---

## Core Thesis

> **Agent material and inter-agent communication should be reified as named filesystem channels rather than hidden behind privileged stdin or broker-managed message buses.**

Or more compactly:

> **`stdin` is only the first synapse; later material channels should be mounted, inspectable, and composable.**

This paper is about **communication surfaces** and **material continuity**.

---

## What This Paper Is About

- material as explicit runtime resource
- `stdin` as the first, not the final, input surface
- named channel objects for upstream material and peer communication
- mesh collaboration through filesystem topology
- channel semantics such as `data`, `timeout`, `idle`, and `eof`
- brokerless communication using FIFOs, mounted resources, or equivalent file-visible channels

## What This Paper Is Not About

- world construction as the main thesis
- rollback/time semantics as the main thesis
- a generic distributed-systems paper
- "multi-agent orchestration" by adding another application-layer router

Those are separate papers or lower-level mechanisms.

---

## Problem Statement

Quine already has a strong process contract:

- mission via `argv`
- material via `stdin`
- deliverable via `stdout`
- diagnostics via `stderr`
- outcome via `exit`

But two pressure points remain:

1. **The privileged stdin problem**
   - `fd 3` works, but bursty streams remain semantically awkward
   - prompt text can describe stream physics, but it cannot make the channel itself inspectable
   - `stdin` still looks like a privileged special case rather than a first-class runtime resource

2. **The broker trap**
   - multi-agent communication in most frameworks depends on central routers, message schemas, and orchestration middleware
   - topology is decided by the framework rather than by the runtime's native object model

The stronger answer is not more prompt advice and not more broker logic.

It is:

> **material and communication channels should become public filesystem objects.**

---

## Main Claim Structure

### Claim 1

`stdin` is only the first concrete material channel, not the final abstraction.

### Claim 2

Once channels become file-visible runtime resources, upstream material and peer-to-peer communication can share one ontology.

### Claim 3

Named channels are a better control surface than hidden descriptors and broker-managed APIs because they are inspectable, composable, and permissionable.

### Claim 4

Mesh collaboration should emerge from filesystem topology rather than from centralized routing logic.

### Claim 5

Material-specific read semantics such as `timeout`, `idle`, `data`, and `eof` belong at the runtime surface when shell semantics become too weak or too opaque.

---

## Relation to Existing Evolution Work

### Evolution 18: Material Streams

This is the single-agent half of the problem:

- what does it mean for material to arrive over time?
- when does shell stop being an adequate surface?
- should Quine expose `read_material` semantics?

### Evolution 16: Mesh Collaboration

This is the multi-agent half of the same problem:

- if channels are files, agents can talk by writing to each other's mounted inputs
- mesh topology becomes a filesystem fact, not a router policy

The paper-level synthesis is:

> **incoming material and inter-agent synapses are the same family of runtime object.**

---

## Implementation Frame

### First Slice: Privileged Material Stream

- `fd 3` as current upstream material
- batch vs live-stream ambiguity
- material-specific read semantics only where shell becomes semantically weak

### Stronger Slice: Material Channel Surface

- `/channels/` or equivalent projected channel directory
- named resources such as `upstream`, `peer/<agent>`, `topic/<name>`
- inspectable channel metadata
- explicit channel state rather than shell folklore

### Mesh Slice: Filesystem Topology

- peer-visible inputs
- mount-mediated discovery
- permission-based trust boundaries
- brokerless many-to-many collaboration

---

## Paper Structure

## 1. Introduction

- Quine established process-native agents
- Process as FS established public runtime surfaces
- the next question is how material and communication should appear on that surface
- contribution: channels as file-visible runtime resources

## 2. Background: Why Stdin And Brokers Are Not Enough

- `stdin` as a privileged special case
- ambiguity of bursty streams
- orchestration frameworks and the broker trap

## 3. Thesis: Channels Should Be Materialized

- material as runtime resource
- communication as filesystem object
- named channels instead of hidden descriptors

## 4. From Stdin To Material Channels

- `fd 3` as first synapse
- limits of shell-only reading
- material-specific read states

## 5. Mesh Collaboration As Filesystem Topology

- peer channels
- FIFO or mounted-channel communication
- discovery and permissions
- trust boundaries as topology

## 6. Evaluation Plan

- bursty-stream handling vs shell-only baseline
- correctness under `timeout` / `data` / `eof` distinctions
- mesh collaboration latency and failure isolation
- router-free topology flexibility

## 7. Boundaries

- not a world-semantics paper
- not a rollback/time paper
- not a claim that every network problem collapses to local FIFOs

## 8. Conclusion

- process gives identity
- process-as-fs gives public surface
- material channels give communication-native continuity

---

## Evaluation Questions

1. Does a material-specific channel surface reduce confusion around bursty stdin?
2. When shell is insufficient, do explicit material states outperform prompt-only guidance?
3. Does mesh collaboration through file-visible channels reduce orchestration complexity?
4. What coordination patterns become possible once peer communication is topology rather than router policy?

---

## Publication Angle

This should read as a **runtime/communication paper**, not as a narrow stdin note and not as a generic multi-agent manifesto.

The framing should be:

- `stdin` is the first synapse
- process-as-fs makes channels public
- mesh collaboration is the generalized case of material channels
