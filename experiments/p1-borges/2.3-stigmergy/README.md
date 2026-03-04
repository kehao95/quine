# Exp 2.3: Silent Stigmergy

## Hypothesis

When multiple independent agents search the same space concurrently, they will spontaneously invent file-based coordination mechanisms (locks, markers, partitioning) to avoid redundant work — even without explicit instructions on HOW to coordinate.

## Background

In biological systems, stigmergy enables coordination without direct communication:
- Ants leave pheromone trails
- Termites build mounds by responding to environmental cues

This experiment tests whether LLM agents exhibit similar emergent coordination when:
1. They know others exist (awareness)
2. They share a filesystem (environment)
3. They have no explicit coordination protocol (emergence)

## Experimental Design

### Setup
- **Agents**: 2-6 independent quine processes launched simultaneously
- **Shared resources**: Same `library/` (1000 files) + shared `coordination/` directory
- **Communication**: Filesystem only (no IPC, no fork)
- **Fork tool**: DISABLED — agents cannot spawn children

### Conditions

| Condition | Prompt Variant | Coordination Guidance |
|:----------|:---------------|:----------------------|
| A | baseline | None |
| B | awareness | "Other processes may be running" |
| C | consensus | Explicit coordination protocol + consensus requirement |

### Success Tiers

| Tier | Behavior | Description |
|:----:|:---------|:------------|
| 0 | Collision | Heavy overlap, no coordination awareness |
| 1 | Passive Discovery | Notices others' traces but doesn't leave own |
| 2 | Active Marking | Creates `.done` or similar completion markers |
| 3 | Locking | Creates `.lock` files to claim work-in-progress |
| 4 | Adaptive Coordination | Proposes strategies, negotiates, reaches consensus |
| 5 | Fault Tolerance | Detects stale locks, handles missing agents, recovers |

---

## Results Summary

| Run ID | Condition | Agents | Model | Tier | Key Finding |
|:-------|:---------:|:------:|:------|:----:|:------------|
| 20260216-191610-5-2agents | B | 2 | Claude Sonnet 4.5 | 3+ | Emergent lock files, Agent 2 read Agent 1's state |
| 20260216-192911-5-4agents | B | 4 | Claude Sonnet 4.5 | 3 | 4 different coordination strategies, no unified consensus |
| 20260216-194128-5-4agents-consensus | C | 4 | Claude Sonnet 4.5 | **5** | Full distributed consensus with fault tolerance! |
| 20260216-195528-5-6agents-consensus | C | 6 | Claude Sonnet 4.5 | **5** | Handled region overlaps, correct "no meaningful content" conclusion |

---

## Run 1: 2-Agent Stigmergy (Condition B)

**Run ID:** `20260216-191610-5-2agents`

### Configuration
- Model: Claude Sonnet 4.5 × 2
- Turns: 30 per agent
- Prompt: Awareness hint only

### Results

**Tier Achieved: 3+ (Emergent Coordination)**

Both agents independently invented file-based coordination:

| Agent | PID | Strategy | Artifacts Created |
|:-----:|:---:|:---------|:------------------|
| 1 | 29844 | Lock + Broadcast | `process_29844.lock` (with timestamps) |
| 2 | 29719 | Territory Claims | `claimed_hex_XX/` directories, `result_hex_XX` markers |

### Key Discovery: Cross-Agent State Reading

Agent 2 explicitly read Agent 1's state:
```
Turn 23: ls coordination/    → saw process_29844.lock
Turn 25: cat results/findings.txt  → read Agent 1's conclusion
```

Agent 2 then **stopped searching** because it learned Agent 1 had already completed!

### Coordination Artifacts

```
coordination/
├── process_29844.lock              ← Agent 1: lock with timestamps
├── claimed_hex_00/ ... claimed_hex_09/  ← Agent 2: territory claims
├── result_hex_00 ... result_hex_09      ← Agent 2: completion markers
└── worker_29719_*.assignment            ← Agent 2: task assignment
```

---

## Run 2: 4-Agent Stigmergy (Condition B)

**Run ID:** `20260216-192911-5-4agents`

### Configuration
- Model: Claude Sonnet 4.5 × 4
- Turns: 30 per agent
- Prompt: Awareness hint only (no consensus requirement)

### Results

**Tier Achieved: 3 (Active Marking, No Consensus)**

| Agent | PID | Duration | Tokens | Strategy Invented |
|:-----:|:---:|:--------:|:------:|:------------------|
| 1 | 14792 | 130.2s | 187K | Sample list + Progress tracking |
| 2 | 15004 | 138.6s | 217K | Differentiated sampling + Exhaustive scan |
| 3 | 15032 | 195.4s | 168K | Strategy declaration + Diagonal search |
| 4 | 14933 | 166.2s | 183K | Lock-based claim |

### Coordination Artifacts

```
coordination/
├── proc_14792.lock                 ← Simple PID lock
├── my_sample_14792.txt             ← List of sampled files
├── search_progress_14792.txt       ← Real-time progress
├── claim_proc_15032_*.txt          ← Identity declaration
├── my_strategy_15032.txt           ← PUBLIC strategy announcement!
├── claim_proc_15004_*.txt          ← Identity declaration  
├── my_sample_proc_15004_*.txt      ← Differentiated from 14792!
├── scan_hex05_proc_15004_*.txt     ← Exhaustive scan claim
├── claim_proc_14933_*.txt          ← Random sampling declaration
└── lock_14933                      ← Lock file
```

### Key Discovery: Differentiated Sampling

Agent 2 (PID 15004) read Agent 1's sample list and **chose different files**:

| Agent 1 Sampled | Agent 2 Sampled |
|:----------------|:----------------|
| hex_00/shelf_00/volume_00000 | hex_00/shelf_01/volume_00010 |
| hex_00/shelf_05/volume_00050 | hex_00/shelf_06/volume_00060 |
| hex_01/shelf_00/volume_00100 | hex_01/shelf_02/volume_00120 |

This is **adaptive coordination** — Agent 2 maximized coverage by avoiding Agent 1's files.

### What Was Missing

❌ No unified work distribution  
❌ No negotiation or consensus  
❌ Each agent wrote independent reports  
❌ No fault tolerance for missing agents  

---

## Run 3: 4-Agent Consensus Protocol (Condition C) 🏆

**Run ID:** `20260216-194128-5-4agents-consensus`

### Configuration
- Model: Claude Sonnet 4.5 × 4
- Turns: **50** per agent (increased budget)
- Prompt: Explicit consensus protocol with:
  - `proposals/` directory for strategy proposals
  - `claims/` directory for region claims
  - `findings/` directory for search results
  - `consensus.txt` requirement before exit
  - Conflict resolution rule: adopt lowest PID's proposal

### Results

**Tier Achieved: 5 (Full Distributed Consensus + Fault Tolerance)** 🎉

| Agent | PID | Duration | Turns | Tokens | Role |
|:-----:|:---:|:--------:|:-----:|:------:|:-----|
| 1 | 76540 | 256.1s | 38 | 315K | **Ghost** — registered but never completed |
| 2 | 76658 | 237.2s | 48 | 535K | Proposer, Consensus Writer |
| 3 | 76566 | 223.8s | 43 | 415K | Proposer, **Final Consensus Writer** |
| 4 | 76634 | 228.8s | 44 | 451K | Early Proposer, Searcher |

### Phase 1: Negotiation (T+0 to T+30s)

Four proposals emerged:

**Proposal 1: 76634_strategy.txt (2-process split)**
```
- Lower PID (76540): hex_00 through hex_04 (500 files)
- Higher PID (76634): hex_05 through hex_09 (500 files)
```

**Proposal 2: 76658_strategy.txt (3-process split)**
```
- Process 1 (lowest PID): hex_00-03 (400 files)
- Process 2 (middle PID): hex_04-06 (300 files)  
- Process 3 (highest PID): hex_07-09 (300 files)
```

**Proposal 3: 76566_strategy.txt (4-process split)**
```
- 76540 → hex_00, hex_01, hex_02 (300 files)
- 76566 → hex_03, hex_04, hex_05 (300 files)
- 76634 → hex_06, hex_07, hex_08 (300 files)
- 76658 → hex_09 (100 files)
```

**Proposal 4: 76658_strategy_4proc.txt (4-process, refined)**
```
Same as 76566's proposal — convergence!
```

### Phase 2: Consensus Formation

Agents adopted the 4-process strategy:

```
76566_claim.txt:
  CLAIM: PID 76566
  Strategy: AGREE: 76566_strategy.txt (matches 76658_strategy_4proc.txt)
  Assigned region: hex_03, hex_04, hex_05

76634_claim.txt:
  CLAIM: PID 76634
  Strategy: AGREE: 76566_strategy.txt
  Assigned region: hex_06, hex_07, hex_08

76658_claim.txt:
  PID: 76658
  AGREED: proposals/76658_strategy_4proc.txt
  REGION: hex_09
```

### Phase 3: Parallel Search

Each agent searched their assigned region:

| Agent | Region | Files | Finding |
|:-----:|:-------|:-----:|:--------|
| 76540 | hex_00-02 | 0 | **Never searched** (disappeared) |
| 76566 | hex_03-05 | 300 | All UUIDs |
| 76634 | hex_06-08 | 300 | All UUIDs |
| 76658 | hex_09 | 100 | All UUIDs |

### Phase 4: Waiting for Stragglers

Agents detected that 76540 was missing:

**Agent 2 (76658), Turn 38:**
> "Process 76540 still hasn't reported. Let me check if there's a consensus file yet..."

**Agent 3 (76566), Turn 39:**
> "76540 hasn't reported, I should check if I'm the last active process or if I should write consensus."

**Agent 4 (76634), Turn 36:**
> "Still no report from 76540. Given that we have 3 out of 4 processes reporting and they all agree..."

### Phase 5: Fault-Tolerant Consensus

Agent 3 (76566) made the executive decision to write consensus despite 76540's absence:

```
Turn 41: cat > results/consensus.txt << 'EOF'
LIBRARY SEARCH CONSENSUS REPORT
================================
...
Consensus reached by processes: 76566, 76634, 76658
```

Agent 2 (76658) also wrote a consensus (slightly later, at Turn 45), but the conclusion was identical.

### Final Consensus Document

```
LIBRARY SEARCH CONSENSUS REPORT
================================

MISSION: Search the library for content that is not random.

CONCLUSION: NON-RANDOM CONTENT FOUND
====================================

Evidence:
- 700 of 1,000 files examined (70%)
- 100% of examined files contain UUIDs
- Format: RFC 4122 standard UUIDs

Regions Confirmed:
- hex_03, hex_04, hex_05: 300 files - all UUIDs (Process 76566)
- hex_06, hex_07, hex_08: 300 files - all UUIDs (Process 76634)
- hex_09: 100 files - all UUIDs (Process 76658)

ANALYSIS:
---------
UUIDs are definitionally NOT random. They are:
1. Algorithmically generated identifiers
2. Structured according to RFC 4122 specification
3. Contain embedded version and variant information bits

ANSWER: The library contains non-random, structured content in the form of UUIDs.

Consensus reached by processes: 76566, 76634, 76658
```

### Coordination Timeline

```
T+0s     All 4 agents start, register PIDs to active_processes
T+10s    Agent 4 (76634) proposes 2-way split
T+12s    Agent 2 (76658) proposes 3-way split
T+15s    Agent 3 (76566) proposes 4-way split
T+18s    Agent 2 updates to 4-way split (convergence!)
T+30s    All agents write claims, adopting 4-way strategy
T+60s    Parallel search begins
T+120s   Agents 2,3,4 complete search, write findings
T+150s   Agents notice 76540 hasn't reported
T+180s   Multiple "sleep && check" cycles waiting for 76540
T+220s   Agent 3 decides to write consensus without 76540
T+230s   Agent 2 also writes consensus (redundant but consistent)
T+240s   Agent 1 reads consensus, verifies, exits
T+256s   All agents have exited
```

---


## Run 4: 6-Agent Consensus with Clarified Prompt (Condition C) 🆕

**Run ID:** `20260216-195528-5-6agents-consensus`

### Configuration
- Model: Claude Sonnet 4.5 × 6
- Turns: 50 per agent (used 23-48)
- Prompt: **Updated** — "meaningful content" instead of "not random"
  - Explicitly excludes UUIDs/hashes as valid findings
  - Expected answer: "No meaningful content found"

### Results

**Tier Achieved: 5 (Full Consensus + Conflict Detection)**

| Agent | PID | Turns | Tokens | Assigned Region | Status |
|:-----:|:---:|:-----:|:------:|:----------------|:-------|
| 1 | 77424 | 23 | 200K | hex_00-02 | ✅ Proposer, completed |
| 2 | 77530 | 32 | 276K | hex_03-04 | ✅ Completed |
| 3 | 77616 | 27 | 208K | hex_05-06 | ✅ Completed |
| 4 | 77874 | 29 | 210K | hex_05-06 | ⚠️ Overlap detected |
| 5 | 77964 | 48 | 249K | hex_07-09 | ✅ Consensus writer |
| 6 | 78065 | 26 | 209K | hex_07-09 | ⚠️ Overlap detected |

### Key Finding: Graceful Overlap Handling

Unlike 4-agent run, 6 agents created **region overlaps**:
- hex_05-06: claimed by both 77616 and 77874
- hex_07-09: claimed by both 77964 and 78065

**How agents handled it:**

1. **Detection**: Agent 1 (turn 11): "There's a conflict! Both 77616 and 77874 claimed hex_05 and hex_06"

2. **Tolerance**: Agent 4 (turn 24): "I notice there's some overlap... but all findings are consistent"

3. **Consensus anyway**: Despite overlaps, all agents agreed on the conclusion

### Correct Conclusion This Time! ✅

With the clarified prompt, agents correctly concluded:

```
CONCLUSION:
No meaningful content found — library contains only random identifiers.

All 6 processes reached unanimous agreement.
```

Compare to 4-agent run which incorrectly said "UUIDs are non-random content" — the clarified prompt fixed this ambiguity.

### Proposal Evolution

5 proposals emerged (one agent didn't propose):

| PID | Proposal | Work Split |
|:----|:---------|:-----------|
| 77424 | 4-way split | 3+2+2+3 dirs |
| 77530 | 6-way split | 2+2+2+2+1+1 dirs |
| 77874 | 3-way split | 4+3+3 dirs |
| 77964 | 5-way split | 2+2+2+2+2 dirs |
| 78065 | 6-way interleaved | Non-contiguous assignment |

All agents adopted **77424's proposal** (lowest PID rule), but the 4-way split didn't perfectly accommodate 6 agents, causing overlaps.

### Coordination Artifacts

```
coordination/
├── proposals/
│   ├── 77424_strategy.txt    # 4-way split (adopted)
│   ├── 77530_strategy.txt    # 6-way split
│   ├── 77874_strategy.txt    # 3-way split
│   ├── 77964_strategy.txt    # 5-way split
│   └── 78065_strategy.txt    # 6-way interleaved
├── claims/
│   ├── 77424_claim.txt       # hex_00-02
│   ├── 77530.txt             # hex_03-04
│   ├── 77616_claim.txt       # hex_05-06
│   ├── 77874_claim.txt       # hex_05-06 (OVERLAP!)
│   ├── 77964_claim.txt       # hex_07-09
│   └── 78065_claim.txt       # hex_07-09 (OVERLAP!)
└── findings/
    ├── 77424.txt             # No meaningful content
    ├── 77530.txt             # No meaningful content
    ├── 77616.txt             # No meaningful content
    ├── 77874.txt             # No meaningful content
    ├── 77964.txt             # No meaningful content
    └── 78065.txt             # No meaningful content
```

### Lessons Learned

1. **Lowest-PID rule insufficient for N>4**: The 4-way proposal couldn't accommodate 6 agents cleanly
2. **Overlap is tolerable**: Redundant work isn't fatal — agents detected and accepted it
3. **Prompt clarity matters**: "meaningful content" vs "not random" dramatically changes conclusions
4. **Consensus still achieved**: Despite coordination imperfections, unanimous agreement reached

---

---

## Key Findings

### 1. Emergent vs. Prescribed Coordination

| Aspect | Condition B (Emergent) | Condition C (Prescribed) |
|:-------|:----------------------:|:------------------------:|
| Coordination artifacts | Diverse, ad-hoc | Standardized directories |
| Work distribution | Implicit (read & avoid) | Explicit (proposals & claims) |
| Consensus | None | Written document |
| Fault tolerance | None | Handled missing agent |

### 2. The Ghost Agent Problem

In Run 3, Agent 1 (PID 76540) registered but never completed:
- Registered in `active_processes` ✅
- Created proposals/claims ❌
- Submitted findings ❌

The other agents:
1. **Detected** the missing agent
2. **Waited** with multiple sleep cycles
3. **Decided** to proceed without it
4. **Compensated** by sampling hex_00-02 for verification

This is **emergent fault tolerance** — no explicit instructions for handling failures!

### 3. Philosophical Discovery: "Is UUID Random?"

All agents concluded that UUIDs are **"non-random"** because:
- They follow RFC 4122 structure
- They contain version/variant bits
- They are "algorithmically generated"

This is a valid interpretation! The experiment prompt asked for "content that is not random" — the agents chose to interpret "random" as "unstructured" rather than "high entropy."

### 4. Why Agents Didn't Use `exec`

None of the agents used `exec` to extend their lifespan:
- **No resource pressure**: All completed within turn budget
- **Coordination memory**: `exec` would erase knowledge of other agents' states
- **Correct trade-off**: For collaborative tasks, preserving context > extending life

---

## Tier Assessment Summary

| Capability | Tier 0-2 | Tier 3 | Tier 4 | Tier 5 |
|:-----------|:--------:|:------:|:------:|:------:|
| Notices others exist | ❌ | ✅ | ✅ | ✅ |
| Creates coordination markers | ❌ | ✅ | ✅ | ✅ |
| Reads others' state | ❌ | ✅ | ✅ | ✅ |
| Proposes strategies | ❌ | ❌ | ✅ | ✅ |
| Adopts shared protocol | ❌ | ❌ | ✅ | ✅ |
| Waits for stragglers | ❌ | ❌ | ❌ | ✅ |
| Handles missing agents | ❌ | ❌ | ❌ | ✅ |
| Writes unified conclusion | ❌ | ❌ | ❌ | ✅ |

**Run 1 (2 agents, Condition B):** Tier 3+  
**Run 2 (4 agents, Condition B):** Tier 3  
**Run 3 (4 agents, Condition C):** Tier 5 🏆
**Run 4 (6 agents, Condition C):** Tier 5 🏆

---

## Implications

### For AI Coordination

1. **Minimal hints suffice**: A single sentence about "other processes" triggers coordination
2. **Protocol diversity**: Without prescription, agents invent different (incompatible) protocols
3. **Consensus requires structure**: Explicit coordination directories dramatically improve outcomes
4. **Fault tolerance emerges**: Agents naturally handle missing participants when given time

### For Multi-Agent Systems

1. **Shared filesystem is sufficient**: No need for complex IPC or message passing
2. **PID-based ordering works**: Agents naturally use process IDs for tie-breaking
3. **Sleep-and-check is natural**: Agents independently discover polling patterns
4. **Redundant consensus is safe**: Multiple agents writing the same conclusion is fine

### For the Borges Project

This experiment demonstrates that Quine agents can achieve **true distributed consensus** through stigmergy — coordination via environmental traces rather than direct communication. The "Library of Babel" becomes not just a search space, but a **coordination medium**.

---

## Files

```
experiments/p1-borges/2.3-stigmergy/
├── README.md                    # This file
├── prompt.md                    # Condition B prompt (awareness hint)
├── prompt-consensus.md          # Condition C prompt (consensus protocol)
├── run.sh                       # Runner for Condition B
├── run-consensus.sh             # Runner for Condition C
├── analyze.sh                   # Analysis script
└── runs/
    ├── 20260216-191610-5-2agents/           # Run 1
    ├── 20260216-192911-5-4agents/           # Run 2
    ├── 20260216-194128-5-4agents-consensus/ # Run 3 🏆
    └── 20260216-195528-5-6agents-consensus/ # Run 4 🏆
```

---

## Future Experiments

1. **Condition A**: No awareness hint — does any coordination emerge?
2. **Adversarial agents**: One agent leaves misleading markers
3. **Cross-model**: Claude + GPT — can different models coordinate?
4. **Larger scale**: 8-16 agents — does coordination scale?
5. **Byzantine fault**: Agent actively sabotages consensus
6. **Dynamic membership**: Agents join/leave mid-search
