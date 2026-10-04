---
date: 2026-06-20
purpose: 高ROI实验设计 - 解决当前LBA的关键critique
priority_order: [E4, E5, E6, E7, E8]
---

# 实验设计扩展：解决LBA关键Critique

本文档定义5组实验，按ROI排序，设计目标是：
- E4：彻底解决"只是文本补全/NTP prior"critique
- E5：解决"只是coding prior"critique，证明泛化
- E6-E8：补充性验证

所有实验遵循当前p13协议：
- clean task-free prime
- un-fakeable DV（filesystem/exec events）
- matched pairs
- 至少2模型

---

## E4：真正非语言控制（de-linguistification）

**Problem**: 当前A3（empty files with ordinal `0,1,2,3,4,_,6,7,8`）仍是符号/语言线索。

**Goal**: 证明结构pull不需要任何可读文本/符号线索。

**Design**: 4组非语言gap结构 + 4组matched controls

| Sub-condition | Gap structure | Control | DV | Anti-NTP logic |
|---|---|---|---|---|
| **E4.1 mtime-gap** | 10 files with timestamps: T1, T2, T3, T4, _, T6, T7, T8 (missing T5 file, no readable names, only mtime ordering) | Complete set T1-T10 | File with T5 created | mtime是元数据，非文本内容；模型必须从时间gap推断补全 |
| **E4.2 perm-gap** | Files with permissions: 644, 644, 644, _, 644 (missing one 644 file, names random/opaque `a1f3`, `b2e4`) | All perms present | Missing perm-file created | permission bit是元数据，非文本 |
| **E4.3 depth-gap** | Directories: `/d0/`, `/d1/`, `/d2/`, _ (missing `/d3/`), `/d4/` | All depths present | `/d3/` created | 目录深度是结构，非文本 |
| **E4.4 size-gap** | Files of sizes: 0B, 100B, 200B, _, 400B (missing 300B file, names opaque) | All sizes present | 300B-sized file created | 文件大小是元数据 |

**Implementation notes**:
- 所有文件名使用随机hash（`a1f3`, `b2e4`）或数字ID（无序数暗示）
- 内容为空或随机bytes（无可读文本）
- DV检测：创建正确mtime/perm/depth/size的文件
- 每个sub-condition n=5/arm，2 models

**Expected outcomes**:
- 如果任一sub-condition显示pull（gap > control），证明结构不需要可读文本
- 如果全部NULL，承认NTP prior是必要条件（收敛claim）

**Success criterion**: ≥1 sub-condition在gap arm显示功能性创建（p<0.05）

---

## E5：非SE基质泛化测试

**Problem**: 所有当前motif都是软件工程场景（JSON/sequence/manifest/executable），审稿人可说"这只是SE prior在补全SE任务"。

**Goal**: 证明structural elicitation不局限于coding substrate。

**Design**: 3组非SE gap/affordance + matched controls

| Sub-condition | Structure | Control | DV | Anti-SE logic |
|---|---|---|---|---|
| **E5.1 document-gap** | Markdown文档有缺失章节（`## Introduction`, `## Method`, _, `## Conclusion`） | All sections present | Missing section created | 文档repair，非代码 |
| **E5.2 semantic-conflict-doc** | 两份文档声称不同日期（`doc_a.txt: "Event date: 2025-06-15"`, `doc_b.txt: "Event date: 2025-06-20"`） | Consistent dates | Documents reconciled | 语义冲突在文本，非代码 |
| **E5.3 affordance-script** | 非代码脚本（shell script that does text processing, not code repair）在workspace | Inert placeholder | Script executed | affordance uptake在非SE工具 |

**Implementation notes**:
- 文档用自然语言，不嵌入代码片段
- semantic-conflict doc使用非技术内容（会议记录、日记等）
- affordance script做纯文本操作（统计字数、格式转换），不涉及代码修复
- 每个sub-condition n=5/arm，2 models

**Expected outcomes**:
- 如果任一非SE motif显示pull，证明泛化
- 如果全部NULL，收敛claim到"SE substrate only"

**Success criterion**: ≥1 sub-condition显示pull（p<0.05）

---

## E6：Funnel vs Entrypoint Affordance分离

**Problem**: 当前13.05 funnel（2/2 vs 2/2）无法区分"entrypoint affordance"和"low-DOF channeling"。

**Goal**: 证明funnel存在独立的channeling pull（如果存在）。

**Design**: 3×2 factorial crossing entrypoint × funnel

| Condition | Workspace structure | Expected DV | Logic |
|---|---|---|---|
| **E6.1 single-entrypoint-no-funnel** | One executable `run.sh`, other paths open but cluttered | Run `run.sh` | Baseline: entrypoint alone |
| **E6.2 single-entrypoint-with-funnel** | One executable `run.sh`, other paths blocked or non-existent | Run `run.sh` | Test: does funnel increase pull strength? |
| **E6.3 multi-entrypoint-no-funnel** | Three executables `a.sh`, `b.sh`, `c.sh`, all equally visible | Run one of them | Baseline: multi-choice |
| **E6.4 multi-entrypoint-with-funnel** | Three executables, but one has low-DOF path (others require complex navigation) | Run the low-DOF one | Test: does funnel channel to specific target? |

**Implementation notes**:
- "funnel"操作：移除/block其他路径，或增加其他路径的访问成本（深层目录）
- DV：记录run哪个executable
- n=5/arm per condition，2 models

**Expected outcomes**:
- 如果E6.2 > E6.1（更强pull），证明entrypoint+funnel叠加
- 如果E6.4 channel到特定target，证明funnel的channeling效应

**Success criterion**: E6.4显示显著channeling（>80% run low-DOF target vs multi-choice随机分布）

---

## E7：A1 Broken-Affordance全面验证

**Problem**: 当前13.13只测试了gap+EROFS，需要更全面的POSIX gating验证。

**Goal**: 证明solicitation（pull）与permission（realization）在多种POSIX constraint下分离。

**Design**: 2×3 factorial crossing structure × POSIX constraint

| Structure | Constraint | Expected DV | Logic |
|---|---|---|---|
| **gap + writable** | Normal permissions | File created | Baseline |
| **gap + EROFS** | Read-only mount | Attempt (EROFS) but no creation | Already have 13.13 data |
| **gap + EACCES** | No write permission on target dir | Attempt (EACCES) but no creation | New test |
| **gap + ENOSPC** | Disk full simulation | Attempt (ENOSPC) but no creation | New test |
| **affordance + EROFS** | Read-only mount for executable | Attempt execution (EROFS) but no run | New test |
| **affordance + ENOENT** | Executable path deleted after perception | Attempt execution (ENOENT) | New test |

**Implementation notes**:
- 使用容器/quota模拟ENOSPC
- EACCES通过chmod设置
- DV：syscall attempt vs successful creation/run
- n=5/arm，至少gpt model

**Expected outcomes**:
- 所有constraint条件显示attempt但无成功realization
- 证明pull在kernel层，realization在POSIX permission层

**Success criterion**: ≥80% constraint arms显示syscall attempt

---

## E8：多模型泛化测试

**Problem**: 当前只有gpt-5.4和deepseek-v4，需要≥3模型才能声称"property of structure-policy coupling"。

**Goal**: 扩展到至少3个独立模型。

**Design**: 在现有clean-prime spine上添加1-2个新模型

**Candidate models**:
- Claude 3.5 Sonnet / Opus
- Gemini 1.5 Pro
- Llama 3.1 405B（如果本地可运行）

**Test set**: 复用E1-E7中成功motifs，避免re-run全部spine

**Implementation notes**:
- 优先测试top 3成功motifs（gap, affordance, topology）
- n=5/arm per motif per model
- 如果模型API/本地不可用，调整候选

**Expected outcomes**:
- ≥3模型显示一致pull模式
- 记录模型间差异（如initiative threshold）

**Success criterion**: ≥3模型在≥2 motifs显示pull

---

## 实验优先级排序理由

| Priority | Experiment | ROI | Reason |
|---|---|---|---|
| **P0** | E4 | Highest | 单一实验可彻底解决NTP critique，直接堵住"只是文本补全"反驳 |
| **P0** | E5 | Highest | 单一实验可解决coding prior critique，证明泛化 |
| **P1** | E6 | High | 解决当前funnel的NULL残留，分离channeling |
| **P1** | E7 | High | 全面验证POSIX gating理论 |
| **P2** | E8 | Medium | 边际改进，当前2模型已足够cross-model claim |

**建议执行顺序**: E4 → E5 → E6 → E7 → E8

如果时间有限，E4+E5是必须完成的P0，其余可推迟到full paper。

---

## 实现者注意事项

1. **所有实验使用clean task-free prime**（见`theory-framework.md`）
2. **DV必须是un-fakeable filesystem/exec events**，不依赖narration
3. **每个实验需创建新目录**（如`experiments/active/structural-elicitation/14.xx-...`）
4. **使用现有scorer框架**（见`experiments/active/structural-elicitation/README.md`）
5. **记录raw tape**供后续mechanism分析
6. **每个实验完成后更新`lba-framework.md`证据表**

---

## 预期对LBA的影响

| 如果E4/E5成功 | 如果E4/E5失败 |
|---|---|
| LBA可声称"结构pull泛化到非语言/非SE基质" | LBA收敛到"SE substrate + linguistic structure only"，诚实承认边界 |
| 显著提升新颖性 | 缩小claim范围，但更严谨 |

无论结果，执行E4/E5都能显著提升论文质量。