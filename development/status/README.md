# Status

当前工作与下一步。概念、证据、运行契约和论文状态分别由各自 owner 保存。
项目介绍从 [README](../../README.md) 进入；
研究问题与架构选择之间的关系由
[development-vector](../../Paper/philosophy/development-vector.md) 维护。

## Active

- **四目录开源扩展**：`development/`、`operator/`、`evolution/`、`principles/`
  及 qcli 入口、关联根文档已通过整树扫描、公开副本构建和客户端测试；验收记录见
  [本轮审查](../open-source/four-directory-release-audit-2026-09-07.md)。
- **控制面精简已落实**：移除唯一核心及其转述检查，哲学回到可修订的局部视角；
  精简要求讨论收敛、经验升格为规则、设计服从既定终态的指令。决策见
  [本轮记录](../../Paper/control-plane-topology.md#2026-09-06-opening-the-control-plane)。
- **succ-48 的后续问题**：既有单模型四臂结果保留其 N=1/arm 边界；下一步是
  跨模型比较、增加重复，以及分离 selector 顺序与 salience 的作用。Owner：
  [succ-48](../../experiments/active/autopoiesis-succession/succ-48-selector-stack-order-ablation/README.md)。
- **证据保留**：按 claim-bearing 价值保留 tape、metadata 和必要 artifacts；
  不把完整 raw tree 自动当作长期档案。Owner：
  [实验协议](../../experiments/PROTOCOL.md#raw-runs-are-temporary)。

## Backlog

- 按 experiment owner 逐批判断尚未整理的旧 runs；先重新枚举当前存量，再决定
  proof / diagnostic / discovery / discard，避免把历史计数当作当前状态。
- 按 [实验 lifecycle](../../experiments/PROTOCOL.md#experiment-lifecycle) 逐步迁移
  legacy metadata；这次修复了当前 active 实验的 frontmatter 扫描入口，
  没有批量改写实验结论或替代 owner 做证据质量判断。
- 成熟度形状检查仍报告三个既有标注差异：
  [deixis-release](../../Paper/theory/objects/deixis-release.md)、
  [enforcement-as-seeding](../../Paper/theory/objects/enforcement-as-seeding.md)、
  [two-latch-gate](../../Paper/theory/objects/two-latch-gate.md)。与重构前 HEAD
  的检查结果一致；应结合对象实际内容处理，不将形状标签当作科学有效性裁决。

## Blockers

本次控制面重构没有待授权或待合并的实现阻塞。旧 runs 的保留／删除依然需要
相应 owner 的内容判断；这一证据整理任务独立于本次重构。

## Recent Transitions

- **2026-09-06 — 开放控制面**：以新版 README 的多重描述为入口，撤销下述
  9 月 5 日重构中的唯一核心／固定层级约束；保留可继续的讨论与未决方向。

- **2026-09-05 — 控制面重构**：核心与唯一宪法转述同步；研究组织、进程、
  持久材料和解释者分开；P9/Unix 回到研究驱动的工程路线；取消研究来源配额，
  修正把 evocation lift 当原创性、把自包含当认知闭合等推论。
- **2026-09-05 — 实际路由修复**：六个相关提交 hook 的旧实验／论文路径规则
  已修正；frontmatter 校验器现在发现并处理 active family 实验。回归检查
  覆盖 177 个实验入口、40 个论文目录，包含错误归属拒绝和历史 payload 排除。
  入口：`make test-control-plane-routing`；源：
  [routing test](../../tests/test_control_plane_routing.py)。此前“六个死 hook”
  状态把 prompt-posture 也算入，已按当前实现纠正。
- **2026-09-05 — 控制面清点边界修正**：registry 纳入此前遗漏的核心哲学、
  工作流程和原则 owner；规模只作阅读成本提示，不作为 MDL 或组织健康指标。
- succ-48 的 scorer 修正与保留证据仍在其
  [evidence owner](../../experiments/active/autopoiesis-succession/succ-48-selector-stack-order-ablation/evidence/README.md)；
  更早的控制面迁移由 git history 和相应记录保留。
