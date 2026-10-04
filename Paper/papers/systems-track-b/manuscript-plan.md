# Manuscript Plan: Track B Full System Paper (OSDI/SOSP Target)

**Target**: OSDI / SOSP (requires Layer 1-5 implementation)  
**Length**: ~12-14 pages + appendices  
**Status**: PLANNING — waiting for Layer 2+ implementation

---

## Strategic Context: Realistic Paper Planning

### 论文价值评估

| 内容 | 实际分量 | 合理目标 |
|------|----------|----------|
| 五层理论框架 | Position paper | arXiv → **HotOS** |
| Layer 1 alone | Feature / demo | **Workshop** / **ATC short** |
| Layer 1 + 2 | Substantial system | **EuroSys** / **ATC** |
| Layer 1-5 complete | Full system | **OSDI** / **SOSP** |

**核心认知**：一个完整的 Agent OS 才值一篇 OSDI，不是每层都能单独成文。

### 发表路径

```
Track A (现在)
│   "Quine: Architecting the Agent OS from POSIX Processes to Plan 9 Primitives"
│   arXiv preprint → HotOS submission (3-4 pages)
│   定位: Position paper — 提出问题和解决思路
│
├── [可选] 中间产出 (Layer 1+2 完成后)
│   Workshop paper 或 EuroSys/ATC
│   验证部分架构，积累评估数据
│
└── Track B: 完整系统论文 (Layer 1-5 实现后)
    "Quine: A Plan 9-Inspired Agent Operating System"
    目标: OSDI / SOSP
    需要: 完整实现 + real workload evaluation + 框架对比
```

---

## Track B: 完整系统论文 Outline

### Title

**"Quine: A Plan 9-Inspired Agent Operating System"**

或：

**"The Agent OS: From POSIX Processes to Plan 9 Primitives"**

---

### 前置条件

Track B 投稿前需要完成：

| 层级 | 核心实现 | 状态 |
|------|----------|------|
| Layer 1 (Identity) | POSIX process mapping | ✅ Done (Quine) |
| Layer 2 (Space) | User Namespace + OverlayFS + Fifth Channel | ❌ TODO |
| Layer 3 (Structure) | Synthetic filesystem (/mnt/agent/*) | ❌ TODO |
| Layer 4 (Time) | Append-only history + Schrödinger's Process | ❌ TODO |
| Layer 5 (Mesh) | 9P protocol + cross-mounted FIFOs | ❌ TODO |

**最低要求**：Layer 1-3 完整实现，Layer 4-5 至少有 prototype。

---

### Abstract (~200 words)

- 问题：Agent 框架在应用层重新发明 OS 抽象
- 方案：五层 Agent OS 架构，从 POSIX 到 Plan 9
- 实现：Quine 运行时（完整五层）
- 评估：vs LangChain/AutoGen/OpenHands on real workloads
- 贡献：首个将 Plan 9 哲学应用于 LLM agent 的完整系统

---

### 1. Introduction (~1.5 pages)

#### 1.1 The Problem

- 框架增殖，各自重新发明 OS 抽象
- 应用层实现的根本局限

#### 1.2 The Solution: A Five-Layer Agent OS

| Layer | Name | Core Mechanism |
|-------|------|----------------|
| 1 | Identity | Agent = POSIX process |
| 2 | Space | Per-process namespace (subjective reality) |
| 3 | Structure | Synthetic filesystem (publicity of mind) |
| 4 | Time | Append-only history (Schrödinger's process) |
| 5 | Mesh | 9P protocol (Unix swarm) |

#### 1.3 Contributions

1. 五层 Agent OS 架构
2. Quine：完整实现
3. 从架构推导的系统属性
4. 与主流框架的 real workload 对比

---

### 2. Background (~1.5 pages)

#### 2.1 Agent Framework Landscape

LangChain, AutoGen, CrewAI, OpenHands 等的结构性问题

#### 2.2 Plan 9 Philosophy

- Everything is a file
- Per-process namespaces
- 9P protocol
- Why it's relevant to agents

#### 2.3 The POSIX Gap

POSIX 解决了什么，没解决什么（引出 Layer 2-5 的必要性）

---

### 3. Layer 1: Identity (~1.5 pages)

*基于 Track A §2-4，扩展实现细节*

- POSIX 映射表
- Host-Guest 架构
- Tool interface (sh, fork, exec, exit)
- 结构性后果（fault isolation, instruction-data separation, unbounded continuation）

---

### 4. Layer 2: Space (~1.5 pages)

- **问题**：全局文件系统污染 + 盲区（silent failure）
- **方案**：User Namespace + OverlayFS
- **实现**：
  - Copy-on-write 隔离
  - Fifth Channel（filesystem diff 注入 tool result）
  - 并发现实（parallel agents with private uppers）
- **评估**：污染防御、盲区消除、COW overhead

---

### 5. Layer 3: Structure (~1.5 pages)

- **问题**：进程内存是黑盒，无法观测/干预 agent 认知
- **方案**：Synthetic filesystem 投射 agent 状态
- **实现**：
  - `/mnt/agent/memory/`, `/mnt/agent/status`
  - Jobs as directories
  - FUSE-based implementation
- **评估**：实时干预效果、调试效率 vs 日志分析

---

### 6. Layer 4: Time (~1.5 pages)

- **问题**：Mind-world desynchronization（context 可回退，filesystem 不可）
- **方案**：Append-only history + observation-driven execution
- **实现**：
  - Three invariants (history is append-only, derivatives never replace, context is constructed)
  - Schrödinger's Process（read() 驱动执行）
  - Unified rollback（mind + world 同步恢复）
- **评估**：rollback 效率、分支探索开销

---

### 7. Layer 5: Mesh (~1.5 pages)

- **问题**：POSIX pipes 是 local/linear，分布式需要 HTTP/gRPC（语义断裂）
- **方案**：9P protocol + cross-mounted FIFOs
- **实现**：
  - Agents as directories
  - Network-transparent file I/O
  - Broker-less collaboration
- **评估**：网络透明性、swarm 扩展性（10 → 1000 agents）

---

### 8. Evaluation (~2.5 pages)

#### 8.1 Microbenchmarks

| Operation | Quine | LangChain | AutoGen | OpenHands |
|-----------|-------|-----------|---------|-----------|
| Agent spawn | ? | ? | ? | ? |
| Context reset | ? | ? | ? | ? |
| Fault recovery | ? | ? | ? | ? |
| Concurrent agents | ? | ? | ? | ? |

#### 8.2 Real Workload Benchmarks

- SWE-Bench subset
- WebArena subset
- Custom multi-step coding tasks

#### 8.3 Case Study

用 Quine 构建一个实际应用，展示五层协作

#### 8.4 Ablation Study

逐层移除，量化每层的贡献

---

### 9. Related Work (~1 page)

- Agent frameworks (LangChain, AutoGen, CrewAI, OpenHands)
- OS-level AI (AIOS, agentOS, MemGPT)
- Historical (Plan 9, Erlang/OTP, Softbots)

---

### 10. Discussion (~0.5 pages)

- Limitations
- When to use / not use
- Broader implications

---

### 11. Conclusion (~0.5 pages)

- 完整的 Agent OS 架构
- Plan 9 哲学在 AI 时代的复兴
- Open source contribution

---

## Estimated Length

| Section | Pages |
|---------|-------|
| Abstract | 0.25 |
| Introduction | 1.5 |
| Background | 1.5 |
| Layer 1-5 (5 sections) | 7.5 |
| Evaluation | 2.5 |
| Related Work | 1.0 |
| Discussion + Conclusion | 1.0 |
| **Total** | **~14 pages** |

---

## 中间产出选项

如果完整系统需要较长时间，可以考虑中间产出：

### Option A: Layer 1+2 → EuroSys/ATC

**Title**: "Subjective Reality: Namespace Isolation for LLM Agents"

- Layer 1 (POSIX mapping) + Layer 2 (Space isolation)
- 聚焦 Fifth Channel 和污染防御
- ~10-12 pages

### Option B: Workshop Paper → HotOS/APSys

**Title**: "Towards an Agent Operating System"

- 五层架构 overview + Layer 1 validation
- ~4-6 pages
- 为完整论文预热

### Option C: Demo/Tool Paper → USENIX ATC Artifact

- Quine 作为开源工具发布
- 工具论文，不需要完整评估
- ~4 pages

---

## Layer 2-5 实现笔记

（待实现时填充详细设计）

### Layer 2 (Space) 实现要点

- [ ] User Namespace 创建和管理
- [ ] OverlayFS 挂载和清理
- [ ] Filesystem diff 计算
- [ ] Fifth Channel 注入到 tool result
- [ ] 并发 fork 的 upper layer 管理

### Layer 3 (Structure) 实现要点

- [ ] FUSE filesystem 框架选择
- [ ] /mnt/agent/* 目录结构设计
- [ ] Memory files 的读写语义
- [ ] Jobs directory 实现
- [ ] 权限控制

### Layer 4 (Time) 实现要点

- [ ] Append-only log 存储格式
- [ ] Snapshot 机制
- [ ] Schrödinger's Process 的 blocking read 实现
- [ ] Unified rollback 原子性保证

### Layer 5 (Mesh) 实现要点

- [ ] 9P server/client 实现（或使用现有库）
- [ ] FIFO 跨挂载
- [ ] 网络透明的 mount
- [ ] Swarm discovery 机制
