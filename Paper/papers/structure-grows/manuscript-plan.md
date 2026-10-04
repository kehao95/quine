# Manuscript Plan: Structure Grows

> **⚠ SUPERSEDED (Jun 25).** This v3 outline (L0→L5 lens escalation + p6
> cleanroom centerpiece) was retired when the manuscript was rewritten to a
> repo-as-living-system spine. The runtime self-replication experiment and the
> cmd/ layout rollback were cut; the lens-rung scaffold was dropped. Retained
> for provenance. The current manuscript is [`./manuscript.md`](manuscript.md);
> the current structure is described in [`./README.md`](README.md) § Current
> State.

**Full title:** Structure Grows: How a Living Repository Differentiates Its Organs
**Target:** 个人博客 (Substack)，前作续篇
**Status:** outline v3 — progressive disclosure + theory escalation (Jun 10); cleanroom centerpiece corrected against real p6 evidence
**Tone:** 第一人称观察者，scene-driven，从具体证据起步、理论逐级攀升

> v3 的两个根本改动:
> 1. **cleanroom 实验描述全面纠正**——它是 p6-autopoiesis 的自我复制研究阶梯，
>    不是"agent 跑 `go build` 撞上一次布局重构"。见 Corrected Centerpiece。
> 2. **结构改为 progressive disclosure + 理论逐级升级**——同一个仓库被揭示五层，
>    每层把理论高度抬一级，最后到达(被挣到的)ALife/自维护系统高度。

## Framing Principle: 克制比喻，但允许并挣得真实理论

(Jun 10 两轮反馈的统一)

前一轮"别太玄"和这一轮"别退太干、保留 ALife/复杂系统视角"不矛盾。区分两件事：

- **装饰性比喻**(immune system、morphogenesis、membrane、把每个 surface 叫一个器官)——**克制**。
  只保留 organ / growth 一个受控意象(它在标题和前作里，是 franchise)。
- **真实理论词汇**(viability、individuation、selection-over-instruction、entropy、
  autopoiesis、open-ended evolution)——**这是这个 lab 实际研究的对象，不是调味**。
  允许使用，但有两条纪律:(a) **逐级挣得**——不在开头就用高阶词，靠证据一级级爬上去;
  (b) **每个理论词都拴住一个具体 arc 或实验**，去掉理论后仍剩可观测事实。
- 判据:如果一句话去掉它后什么都不剩 → 是空话，删;如果剩下一个具体 surface/实验/arc
  → 保留，理论是对它的命名。

## Narrative Posture

第一人称、scene-driven、observer-inside;用前作当**个人引子**(发布时的完结感 →
系统没停 → 被复杂化的"我自己的结论")。个人叙述与理论攀升不冲突:每个情绪/惊讶
节拍后面拴一个真实 arc 或实验样本。

## Core Argument

前作说"仓库可以自维护"。这篇说:

> 一个活的仓库不只是自维护。它在压力下**长出新结构、并且这些结构是被选择出来的而非被设计的**;
> 把它一层层看下去，它表现出的是一个抗熵、维持自身边界、却尚未开放式演化的系统——
> 而它"自治到哪一步"是可量度的。

基础洞察:

1. 结构不是先验设计，而是在压力下分化(旧 container 装不下不同职责被逼分裂)。
2. 系统只能修复它能感知的 defect → 可观测性本身变成一等结构。
3. 更成熟的成长是分化职责、收紧边界，而非增加说明。
4. 成长不等于一直向前;rollback 也是发育证据。
5. Blindness、over-compression 仍是成本;活 ≠ frictionless。
6. 表现出活性的单位是 Human + agent + repo + validators + history 的回路，不是 repo 本体。
7. Human 更像调节者(regulator)/选择环境，不像预知终局的架构师。

## Lens Escalation（视野攀升，不是理论阶梯）

**核心重构(Jun 10 Human)**:攀升的不是"对象的理论高度"，而是**我这个观察者的视野**。
同一个仓库自始至终是同一个东西——变的是我用什么眼睛看它。文章的脊柱因此是
**"我如何一步步学会看它"**:从"看到摩擦就想优化结构"，到"摩擦得长出感知器官"，
到"它作为一个 living system 在演化"。这把高阶理论变成**我的看法的成熟**，而不是对现实的
宏大断言——既诚实、又天然个人叙述、也不玄(玄来自把理论当对象属性宣称;视野攀升把它当
观察者的镜头)。

| 镜头 | 我当时用的眼睛 | 是什么让镜头加深 | 证据锚 |
|------|----------------|------------------|--------|
| **L0** 完结 | "一个会改写自己的系统"——前作给我的、以为已经看到底的镜头 | 系统没停下来;journal 攒到 111 条 | journal 111 条 / 7 期 |
| **L1** 优化者 | 看到摩擦 → 想"把结构优化一下"(拆文件、挪布局) | 同一种摩擦在不同尺度反复出现，单次优化压不住 | themes 层 / THEORY 拆分 / philosophy 层 |
| **L2** 造仪器者 | 摩擦反复是因为系统**看不见自己的漂移** → 摩擦需要长出感知器官(checks/journal) | 我开始给它造眼睛，它甚至开始度量自己多久才发现 | EVOLUTION 创建+自裂、blindness duration、causal-path contract |
| **L3** 选择论者 | 镜头翻转:不是"我设计结构"，而是"环境在**选择**结构" | 自我复制实验:给它源码它也不重建自己——是生态在选被复制的自我;布局被自重建约束否决 | p6 6.59–6.68、cmd/ 布局回滚 |
| **L4** living-system | 我开始把整个回路看成一个**在演化**的活系统:抗熵、维持有边界的自我 | viability、individuation(边界=划出内部)、entropy resistance | DVC 四层现实、Grand ALife Synthesis ①② |
| **L5** 诚实的眼 | 我也看见了这副镜头的边:它**还**没在演化的最强意义上演化 | diagnosis 仍靠人;开放式演化未达;V0 spine retired;sandbox 未修 | Grand ALife Synthesis ③、open-friction、autonomy gauge 读数 |

用户明确点名的三级 = L1(看到摩擦就优化结构)→ L2(摩擦需要感知器官)→ L4(repo 作为 living
system 在演化)。L3(选择)是它们之间那道没被点名但证据最硬的枢纽:它是**镜头从"我造"翻向
"环境选"**的那一下。L3 翻转之后，L4 的 living-system 视野才合法——否则就是硬把活系统语言
扣在对象上。

纪律:**对象恒定，视野攀升。** ALife 高度是一种**看见的成就**(perceptual achievement)，
不是对象的形而上升级。L4/L5 的词不能出现在 L0/L1——那会暴露"我一开始就站在高处"的不诚实。

## Progressive Disclosure Spine（渐进披露 = 我的镜头在加深）

渐进披露的对象是**观察者**——所以它天然就是视野攀升:同一个仓库，我用越来越深的眼睛
看了五遍。这本身就是 Quine 的设计哲学(`operating-philosophy.md` 的 Progressive Disclosure
Map、`THEORY.md` 的分层结构)——文章的**形式**在 enact 仓库自己的原则，coda 点这个 meta payoff。

我学会看它的五层(同一个仓库，不同的眼睛):

- 看见 1:一个让我想去**优化结构**的工作区 (L1)
- 看见 2:…我意识到它得**感知**自己的漂移，于是开始造仪器 (L2)
- 看见 3:…镜头翻转——它的结构是被**选择**的，不是我设计的 (L3)
- 看见 4:…我开始看见整个回路在**演化**:抗熵、维持有边界的自我 (L4)
- 看见 5:…我也看见这副镜头的**边**:还没跨过的那道线 (L5)

§1 给承诺("我会把同一个仓库给你看五遍——其实是给你看我怎么一步步学会看它");coda 兑现。

## Corrected Centerpiece（§4，事实纠正 + 加深）

**旧版错误**:写成"我做了个漂亮的 cmd/ 重构，agent 跑 `go build ./cmd/quine` 失败，回滚"。
那把一个独立的布局 arc 和真正的实验混成了一个浅故事。

**真实素材**(p6-autopoiesis 6.58–6.68，self-source / 自我复制阶梯):

- Quine 字面意义是"输出自己源码的程序"。p6 把它做成实验:把一个 LLM-backed Quine
  放进 repo 外的隔离 cleanroom，给它未来任务流，看它能否**合成一个 successor 并通过真实
  `exec` 交接**，以及 runtime 自带的 read-only `source-code/` 是否改变被复制对象的形态。
- **核心发现(反直觉)**:暴露 `source-code/` 并**没有**让 lineage 直接重建自己。它常常
  搭一个 hybrid 自定义 successor，把自源码当**设计参考/runtime-contract 学习面**用。
  6.62/6.64/6.65 显示:决定它**重建完整 quine 体**(copy `source-code/` → `go build ./cmd/quine`
  → `exec` 进 rebuilt quine)还是搭一次性自定义体的，是**生态压力**——当未来任务只能由
  quine 原生控制面(`ctl/wake`)承载、且交接后存活要求它时，它才重建 quine;当 carrier 是
  raw stdin(6.68)，它搭自定义流体、从不重建 quine。
- 一句话:**蓝图可用 ≠ 会复制自己。是生态在选择被复制的自我对象。** 这正是 Quine runtime
  的"Selection over Instruction"被实验实证。p6 自己的评估维度里就把这层叫 `individuation`。
- **子故事(布局回滚)**:把 cmd/ 重排成 scoped layout 后，投影自重建路径
  `go build ./cmd/quine` 在 projected surface 上失效;这个缺陷**通过 cleanroom 自我复制工作
  当天浮现**，回滚，落到一个显式 re-entry 机制(env var 名后来又变过，文中只描述机制，不点名)。
  即:**仓库结构被"系统能否重建自己"这条硬约束否决了。**

**§4 的纪律(防止跑题成 p6 论文)**:
- 实验是**证据和原则锚**，不是主角。用布局回滚(仓库结构事件)做具体脊柱，用 self-source
  发现做原则(selection 选定自我对象)，然后**立刻拉回仓库**:仓库自己的 organs 也是被选择
  出来的，不是被设计的。
- 不展开 p6 的 rung 细节(6.x 编号、scorer、fd 物理)。读者要的是原则，不是实验台账。
- 点名递归:lab 遵守它要造之物的同一条物理——这条现在有实证牙齿，不只是漂亮对称。

## Autonomy Gauge（贯穿，落在 A5）

开头承诺一把工程化刻度尺，全文兑现，最后是 A5 的读数:
- 轴一·自治度:人退一步后回路还跑通多少——repair 已内化、可观测性正在内化、**diagnosis 仍靠人**。
- 轴二·感知延迟(blindness duration):从受到 friction 到命名 defect 的时间是否在缩短。
死掉但文档完美的 repo 两轴都是零/无穷。gauge 让"活着"是读数不是修辞，并给 A5 一个具体天花板。

## Section Outline（v3：每节是我视野的一次加深）

每一节 = 我换上一副更深的眼睛。对象恒定，看见在变。

### §1 — 我以为故事讲完了（L0：起始镜头）

- 个人引子:发布前作时的完结感("命名即放下")→ 系统没停 → 回看 journal 111 条 →
  发现自己用错了词:不是 repair，是长出原本没有的结构。
- claim 升级声明:前作证明它能维持自己;这篇 claim 更强，我得在这儿挣到。
- 装 autonomy gauge 承诺 + **视野承诺**:"我会把同一个仓库给你看五遍——其实是给你看我
  怎么一步步学会看它。"
- 精度动作:活的单位是回路，不是 repo。

### §2 — 起初我只想把结构优化一下（L1：优化者的眼）

- 我那时的眼睛:看到摩擦 → 反应是"拆个文件、挪个布局"。三场景三尺度同形(加深，讲清旧
  surface 同时干哪几份活、分裂前后各自 owner):themes 层 / theme README→THEORY / philosophy 层。
- 镜头的裂缝:同一种摩擦在不同尺度反复出现，单次优化压不住——开始怀疑"优化"不是对的词。
- 看见 1:这是一个会长出结构的工作区(还只是工程视野)。

### §3 — 然后我意识到它得能看见自己（L2：造仪器者的眼）

- 镜头加深:摩擦反复是因为系统**看不见自己的漂移**;于是"优化结构"换成"给它造感知器官"。
  eyes/memory/wiring 作轻量标签:emergent/cleanroom 评测、EVOLUTION journal 及其自裂、
  causal-path contract。"说不出自己怎么过期的控制面只是记忆，不是维护。"
- **gauge 轴二兑现**:blindness duration——它甚至开始**度量自己多久才发现**。轻接控制论但不堆术语。
- 看见 2:我开始看见它感知自己的漂移(从"我修"到"它得能看见")。

### §4 — 镜头翻转：不是我设计，是环境在选（L3：选择论者的眼 · 重心）

- 这一节是视野的**翻转**那一下。用 Corrected Centerpiece:self-source 自我复制实验 + 布局回滚。
- 反直觉:给它源码它也不重建自己;**是生态在选被复制的自我对象**;布局被自重建约束否决。
  这逼我把眼睛从"我造结构"换成"环境选结构"。
- 拉回仓库 + 点名递归:仓库自己的 organs 也是被选出来的;**lab 在遵守它要造之物的物理**
  (selection over instruction)，这次有实验牙齿。
- 看见 3:结构是被选择的，不是被设计的。这是镜头从"工程"翻向"活系统"的枢纽。

### §5 — 我开始看见它在演化（L4：living-system 的眼，已挣得）

- L3 翻转之后，活系统视野才合法。DVC 四层现实 + README 边界 = 仓库在**划出自己的内部**
  (individuation，工程上=职责分离)。接 viability / entropy resistance(Grand ALife Synthesis ①②):
  分不清"保留 vs 公开"就没有内外，也就没有可维持的自我。
- 自我印证转折:最地基的边界**最晚被我看见**(私有期一直糊着，开源压力才逼出)——印证 L2:
  最深的结构感知延迟最长。
- 看见 4:我开始看见整个回路在抗熵、维持有边界的自我。

### §6 — 这副眼睛也照出它仍然会出错（L4，诚实压舱）

- 活系统视野不是滤镜美化。over-compression(V0 spine 被 retired)、naming overclaim(sandbox 未修)、
  coordination drift(repo/index 双源)。
- **gauge 轴二回收**:latency 是仪表读数，不为零;活 ≠ frictionless，熵是常驻的。

### §7 — 我在这里干嘛 + 这副镜头的边（L5：诚实的眼）

- 三角色:regulator / selective environment / naming authority。
- **gauge 轴一收束**:repair 已内化、可观测性正在内化、**diagnosis 仍靠人**=今天天花板。
- **接 Grand ALife Synthesis ③**:我能看见它抗熵、维持个体性，但**看不见它在开放式演化**——
  没有无外部目的地自发生成累积复杂度。说"闭环了"的是抱负不是观察。这是这副镜头诚实的边。
- 看见 5:我也看见了这副眼睛还没跨过的那道线。

### §8 — Coda：这篇文章也是一个 organ（合上各环）

- 反身收尾:文章活在它描述的系统里，带真实 Unsafe Overclaims 约束(就是 §5 那批边界之一)。
- **视野/渐进披露兑现**:我把同一个仓库给你看了五遍——变的从来是眼睛;而"分层揭示"正是这
  仓库自己的设计哲学，文章的形式在 enact 它的原则。
- **claim 之环收束**:挣到了什么?自维护——挣到了;自扩展——有界版本;开放式演化——没挣到，不声称。
- 末句:不是系统升到了某个高度，是我学会了看见它的哪些层——以及我还看不见、它也还没跨过的那道线。

## Evidence Anchors（已纠正、加深）

| Section | 镜头 | Evidence | Source |
|---------|------|----------|--------|
| §1 | L0 | 111 entries / 7 periods after the first essay | `EVOLUTION.md`, `evolution/periods/` |
| §2 | L1 | themes 层 (Apr 16)、theme README→THEORY (Apr 27)、philosophy 层 (May 1) | `evolution/arc-index.md` 对应行 |
| §3 | L2 | EVOLUTION 创建+三裂、blindness-duration 字段、causal-path contract 五字段 | period XII–XIII、`development/contracts/structural-friction-journal.md`、`development/contracts/control-plane-causal-path.md` |
| §4 | L3 | self-source 自我复制阶梯(生态选定自我对象;源码≠会复制);cmd/ 布局回滚(自重建约束否决重构) | `experiments/active/autopoiesis-succession/{6.59,6.60,6.62,6.64,6.65,6.68}/README.md`、`p6-autopoiesis/README.md` ladder、EVOLUTION period XII(Apr 19–20) |
| §5 | L4 | DVC retention/restore/publication/curation、README 边界;viability/entropy(Grand ALife Synthesis ①②) | period XVI、arc-index May 11–27 行、`Paper/philosophy/self-maintaining-systems.md` |
| §6 | L4 | retired V0 spine、sandbox overclaim、index drift | `EVOLUTION.md` Open Friction、`ROADMAP.md` Retired Framing |
| §7 | L5 | diagnosis 未内化;open-ended evolution 未达(Grand ALife Synthesis ③);selection over instruction | `README.md` core-idea 表、`self-maintaining-systems.md`、essay-local synthesis |
| §8 | — | 本 dossier 自己的 Unsafe Overclaims;Quine 自己的 progressive-disclosure 哲学 | 本 README、`operating-philosophy.md`、`THEORY.md` |

## Anti-Patterns

- **克制装饰性比喻，但挣得真实理论**——organ/growth 一个受控意象;ALife/复杂系统词汇逐级挣得、
  每个拴住实验或 arc(见 Framing Principle)。
- **不要在开头就站在 ALife 高度**——progressive disclosure 的意义就是爬上去;A4/A5 的词不能出现在 A0/A1。
- **§4 不要变成 p6 论文**——实验是原则锚，不是主角;不展开 rung 编号/scorer/fd 物理;讲完原则立刻拉回仓库。
- **cleanroom 事实必须对**——不是"漂亮重构被 go build 打脸"，是"生态选定自我对象 + 自重建约束否决布局"。
- **不要变成 commit log**——读者要 story，且公开历史已 rewrite，不引 sha。
- **不要重复前作**——Strange Loop / thermodynamics / autopoiesis 引入已由前作完成，link 过去。
- **不要把 autonomy gauge 写成精确指标**——ordinal 读数(哪些能力已内化)，不是带小数点的度量。
- **不要自我恭维**——§6 honest concerns、§7 "diagnosis 仍靠人 + OEE 未达" 是必须的诚实顶。
- **不要假装两个月所有变化是一个 trend**——承认 rollback、retirement、misframing。
