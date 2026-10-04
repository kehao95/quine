# How I Open-Sourced a Living Research Workbench

---

Most people think open-sourcing a repository is a weekend chore. You scan for
API keys, strip out `.env` files, maybe squash the embarrassing early commits,
and push. For a conventional software project, that is roughly correct.

Quine is not a conventional software project.

Quine is a research program studying autopoietic intelligence — an AI system
that attempts to understand and rebuild itself through recursive POSIX
processes. Its repository is not a codebase with some documentation attached.
It is a living research workbench: runtime code, theoretical notes, experiment
protocols, raw agent interaction logs, retained run trajectories, draft papers,
control-plane surfaces, and over a thousand commits of evolutionary history all
cohabit the same Git tree. Some of what a normal engineering team would discard
as build residue functions here as primary scientific evidence.

When I decided to make this repository public, I discovered that the standard
open-source playbook — clean the secrets, carve out a public branch, ship —
was not merely insufficient. It was conceptually wrong. The hard problems were
not about what to hide. They were about what to *keep*, what *keeping* even
means when your evidence can leak your secrets, and whether a cleaned-up
history can still be an honest one.

This post is about how I resolved those problems.

---

## I. Why "Just Ship the Code" Is No Longer Enough

There is a growing sense in the AI community that open-sourcing pure runtime
code is becoming a hollow gesture. Code in the LLM era is cheap. An agent can
rewrite a runtime in an afternoon. What cannot be cheaply reproduced is the
research itself: the theoretical dead ends, the failed experiments, the
gradual shifts in understanding that only become visible when you can read the
full chronological record of a project's evolution.

I felt this acutely with Quine. If I carved out a sanitized `public-runtime`
branch containing only the executable agent code, I would be releasing a
product, not sharing a research program. The resulting repository would look
clean and professional and completely dead — a museum exhibit rather than a
working lab.

What I actually wanted to share was the workbench: the theory folder where
ideas competed with each other, the experiment families where most runs failed
and a few succeeded, the protocol documents that evolved as we learned what
"retaining evidence" actually required, the Paper directory where draft essays
lived alongside their source analysis. I wanted a reader — or an agent — to
be able to walk through the repository and understand not just what Quine does,
but how it got here.

That decision — to open-source the research, not just the runtime — turned a
weekend chore into a months-long boundary-design problem.

---

## II. The Deadlock: Logs Are Memory, and Memory Leaks Secrets

### Why logs cannot be thrown away

In most software projects, experiment logs are disposable. You keep the test
results, maybe archive the CI output, and delete the raw traces. Nobody needs
last Tuesday's `stdout`.

LLM research breaks this assumption at its root. When your experimental
subject is a large language model, the behavior has real stochasticity. Even if
you freeze the model weights, pin the temperature, fix the random seed, and
replay the exact same prompt, you will not deterministically recover the same
trajectory. The original run output is not reproducible later. It is a
one-time observation.

For Quine, this was especially acute. The experiments involved an agent
operating inside a live POSIX environment, issuing shell commands, reading and
writing files, and interacting with external LLM providers across multiple
turns. Each run produced a dense interaction tape — a full log of prompts,
responses, tool calls, file mutations, and observed outputs. These tapes were
not engineering residue. They were the raw experimental data.

More than that, the tapes were a *read surface*. When the theory advanced
months later, we needed to go back and reread old runs under a new analytical
frame. Sometimes a human did the rereading. Sometimes the point was to let an
agent ingest old trajectories and extract patterns we had not noticed before.
The logs were not write-once archives. They were living evidence that the
project expected to revisit repeatedly.

Early in the project's life, the implementation of this principle was blunt:
track everything directly in Git. Tape files, stdout captures, run metadata,
environment snapshots — all committed alongside the code that produced them.
This was not carelessness; it came from the belief that losing any of this
material was scientifically unacceptable.

It worked, until it didn't.

### Why logs cannot be published

The project began its life as a private repository, and it operated accordingly.
The Quine agent inherited its host environment when it ran. Shell commands
could print environment variables. Trace files captured absolute local paths.
Run metadata automatically snapshotted the environment before and after
execution, including variables that happened to contain API keys, provider
endpoints, and host-specific configuration.

None of this was intentional leakage. It was the natural byproduct of a
research agent operating honestly inside a real system. But it meant that the
very artifacts most worth preserving — the interaction tapes, the stdout
captures, the environment snapshots — were also the most dangerous surfaces
to publish.

A scan of the repository's history revealed the scale of the problem. The
history contained **93,712** distinct paths matching run, tape, session, log,
and workspace patterns. Among those were files like `meta/env-after.txt` and
`meta/summary.env` — environment snapshots that the experiment framework
captured automatically at the end of each run. These files were part of the
scientific record. They also contained host secrets.

The threat model was not "someone committed a `.env` file." The threat model
was: *the normal, correct operation of the experiment framework produced
evidence that structurally contained sensitive host information.*

### The trap

So here was the deadlock:

- **Delete the logs**, and you destroy the scientific memory layer. Future
  theory work loses access to irrecoverable raw observations.
- **Publish the logs**, and you expose host secrets, provider configurations,
  and local paths to anyone who clones the repository.
- **Keep a private branch with full logs and a public branch without them**,
  and you end up maintaining two divergent realities. The public branch
  becomes a lie — a sanitized fiction that omits exactly the material that
  makes the repository a research workbench rather than just a software
  project.

Every obvious move was wrong.

---

## III. The Obsession: History Must Not Be Faked

At this point, the pragmatic engineering voice in my head offered a simple
solution: *Just start over.*

Squash the thousand-plus commits into a single initial commit. Delete all the
historical log files. Begin the public repository from a clean slate with only
the current-state files, neatly organized and free of any sensitive residue.
Many respected open-source projects do exactly this when they transition from
private to public development. It is fast, safe, and industry-standard.

I could not bring myself to do it.

This is where I have to be honest about a personal bias that shaped every
subsequent technical decision. I have what might generously be called an
obsession with historical authenticity — or less generously, a cleanliness
compulsion directed at the wrong target. The idea of publishing a repository
whose history begins with a single pristine commit, when the actual work
evolved through over a thousand messy, real, timestamped steps, felt like
academic fraud. Not legally. Not even ethically, probably. But aesthetically, in
the deepest sense of the word. The history *was* the work. The early confusion,
the wrong turns, the gradual refinement of the experiment protocol, the
compression of the control plane — all of it was visible in the commit
sequence. Deleting that sequence to produce a "clean" repository was not
cleaning. It was erasure.

So I set myself a constraint that made the entire problem harder:

> Scrub the poison — the secrets, the host paths, the sensitive environment
> snapshots — but do not alter the footprints. Every commit stays. Every
> timestamp stays. Every directory structure stays. The repository's
> evolutionary skeleton must survive intact.

This meant I could not use any of the easy shortcuts. I could not squash
history. I could not start a fresh repository. I could not delete the
experiment directories and pretend they had never existed. I needed a way to
reach into the past and surgically replace dangerous *payloads* while leaving
the *structure* untouched.

---

## IV. The Pivot: DVC as a Boundary Protocol

The breakthrough came from an unexpected reframing of a tool I was already
using.

[DVC](https://dvc.org) (Data Version Control) is conventionally understood as
a large-file storage solution for Git — a way to keep heavy datasets out of
the Git object database while preserving their semantic position in the
repository tree. You replace a 30 MB tarball with a tiny `.dvc` pointer file
that records the content hash, and the actual payload lives in a remote storage
backend. Standard practice.

But in the context of the release problem, DVC became something much more
interesting. It became a **boundary protocol** — a mechanism for decoupling
three concerns that had been hopelessly entangled inside the Git tree:

1. **Structure**: the directory layout, file paths, and semantic organization
   of the repository.
2. **Payload**: the actual bytes of experiment logs, run tapes, archived
   bundles, and binary outputs.
3. **Publication**: the question of whether a given payload should be publicly
   restorable.

Before DVC, these three things were fused. If a file existed in the Git tree,
its structure, payload, and publication status were all the same fact: it was
tracked, it was stored in Git objects, and anyone who cloned the repo could see
it. The only way to change any one of these was to change all three.

DVC broke that fusion. After externalization:

- **Structure** was preserved. The `.dvc` pointer file sat in exactly the same
  path where the original file had been. The repository's directory tree
  remained legible. A reader could still see that
  `experiments/active/benchmarks/bench-02-tb2-sweep/runs/run-006.tar.zst` existed and
  was part of the experiment family. The semantic organization was intact.
- **Payload** was moved to a separately controlled DVC remote — in this case,
  a self-hosted server running on personal infrastructure. The heavy bytes
  were no longer in Git history. A fresh clone of the public branch dropped
  from **469 MiB** to **27 MiB**.
- **Publication** became an independent decision. This was the critical
  insight.

### The publication boundary

Because the DVC remote was self-hosted, I had full control over what it served
to the public. This turned DVC into something it was not originally designed to
be: a *manifest-scoped publication mechanism*.

The implementation was a single text file: `public-manifest.txt` in the `.dvc`
directory. This file listed, by path, every DVC-tracked payload that had been
reviewed and approved for public access. Only payloads on the manifest could
be restored by a public user running `dvc pull`. Everything else remained in
private storage — retained, recoverable by the project, but invisible to the
outside world.

This created a three-layer architecture that resolved the original deadlock:

| Layer | Visibility | Contents |
|-------|-----------|----------|
| **Public Source & History** | Anyone who clones the repo | Code, theory, experiment *structure*, DVC pointers, papers, protocols, control surfaces. The full evolutionary skeleton. |
| **Private Retention** | Project-internal only | The complete experiment memory: all run tapes, all environment snapshots, all raw outputs. Retained for future scientific rereading. Nothing lost. |
| **Published Evidence** | Public via manifest allowlist | A curated subset of experiment artifacts that have been individually reviewed, confirmed clean of sensitive material, and approved for public restore. |

The retained private layer was the key to the whole design. It meant the
project did not have to choose between scientific integrity and operational
safety. The full experimental memory continued to exist, privately. The logs
were not deleted. The tapes were not destroyed. The environment snapshots were
not lost. They were simply moved behind a boundary that separated *retention*
from *publication*.

---

## V. The Surgery: Rewriting a Thousand Commits

Externalizing current-state files to DVC solved the problem for the HEAD of the
repository. But Git remembers everything. The historical commits still
contained the original heavy payloads — the raw tarballs, the stdout captures,
the environment snapshots, the built binaries. A fresh clone of the private
branch carried **469 MiB** of pack data, and the reachable historical blobs
totaled nearly **1.5 GB** across **47,565** objects.

Cleaning the current tree was not enough. The history itself had to be
rewritten.

### The pathset

The first step was mapping the attack surface. A systematic scan of the full
commit history produced a rewrite pathset:

- **93,712** distinct paths matching run, tape, session, log, and workspace
  patterns
- **527** distinct tarball and compressed archive paths
- Hundreds of additional paths carrying provider profiles, environment
  snapshots, local-path metadata, and built binaries

These were not scattered needles in a haystack. They were the haystack. The
experiment-record infrastructure had woven itself through the entire history
of the repository.

### The rewrite

The tool was `git-filter-repo`, operating on the private `lab` branch to
produce a sanitized `main` branch. The rewrite followed a principle:

> Transform the payload form. Preserve the commit form.

For each historical commit, the rewrite identified files matching the
dangerous path families and either:

- replaced them with DVC pointer files (preserving the semantic position while
  removing the payload from Git history), or
- removed them entirely when no DVC retention was needed (e.g., incidental
  build artifacts, intermediate binaries).

What was *not* touched:

- Commit hashes were necessarily recomputed (this is unavoidable in any history
  rewrite), but the commit *metadata* — author, timestamp, message, parent
  structure — was preserved exactly.
- Directory trees retained their original shape. If an experiment directory had
  contained five run folders in the original history, it still contained five
  run folders in the rewritten history. The *contents* of those folders
  changed from raw payloads to DVC pointers, but the *structure* remained.
- Source code, theory notes, protocol documents, and paper drafts passed
  through untouched.

### The result

The numbers told the story:

| Metric | Private (`lab`) | Public (`main`) | Change |
|--------|----------------|-----------------|--------|
| Fresh clone pack size | 469 MiB | 27 MiB | **−94%** |
| Reachable blob count | 47,565 | 16,099 | **−66%** |
| Reachable blob bytes | 1.48 GB | 185 MB | **−87.5%** |
| Tracked files at tip | 6,043 | 5,822 | −3.7% |

The last row is the most telling. The number of tracked files barely changed.
The repository did not shrink by *removing structure*. It shrank by
*extracting payload*. The skeleton survived.

A reader cloning the public `main` branch gets a 27 MiB download that
contains the complete evolutionary record of the project: every experiment
family, every protocol revision, every theoretical note, every paper draft,
every control-plane surface. What they do not get — unless the manifest
explicitly allows it — are the heavy raw payloads behind the DVC pointers.
The structure tells the story. The pointers acknowledge the evidence. The
boundary controls the access.

---

## VI. What Actually Became Public

The result was not a minimal open-source release. It was not a curated
showcase repository. It was a sanitized but honest public research workbench.

The public Quine repository tells the truth about itself. Its README explains
that Quine is both a runtime and a research workbench. It labels draft papers
and working theoretical notes as exactly that — drafts, not publications. It
documents how to use the public DVC manifest to restore approved experiment
evidence, and it makes clear that the manifest is an allowlist, not a default.
The CITATION file treats the project as a research artifact, not just software.

Some of what a public reader sees is polished. Most of it is not. That is the
point. A research workbench is not a product. It is a place where work
happens, and the traces of that work are part of the value.

---

## VII. What This Taught Me

### Separate the three boundaries early

If you are building a research repository that might someday become public,
the single most valuable architectural move is to recognize that **public
source**, **private retention**, and **published evidence** are three
different things — and to keep them separate from the start.

Most projects conflate at least two of these. A file is either in the public
repo or it is not. But research repositories need a middle layer: material
that is retained for scientific integrity but not published for public access.
If you design that layer early, the eventual open-source transition becomes a
boundary adjustment rather than a reconstruction.

### Retention is not publication

This is the corollary, and it is worth stating on its own. The fact that
you *keep* a file does not mean you *publish* it. Private retention and
public publication are independent facts. DVC with a self-hosted remote and a
manifest-scoped allowlist is one way to implement this separation. There are
certainly others. The important thing is the conceptual split, not the
specific tooling.

### Treat the release as a control plane, not a cleanup

Open-sourcing a complex research repository is not a one-time scrubbing
exercise. It is a stateful, multi-phase process with its own scope, its own
audit trail, its own blockers, and its own owner surface. Treating it as a
managed process — with formal scans, inventories, allow/deny reviews, and
validation steps — made it tractable. Treating it as a vague "cleanup
campaign" would have made it endless.

### History is worth more than cleanliness

This is the most personal lesson, and I suspect the most controversial. There
is real value in preserving the authentic historical record of a research
project, even when that record is messy, even when cleaning it would be easier,
and even when no one will ever read most of it. The mess is the evidence.
The timestamps are the chronology. The dead ends are the negative results.
A fabricated clean history may look better, but it is less true — and for a
research project, truth is the only currency that does not depreciate.

---

*Quine is open-sourced at [github.com/kehao95/quine-lab](https://github.com/kehao95/quine-lab). The public DVC manifest and experiment evidence are available via the instructions in the repository's QUICKSTART guide.*
