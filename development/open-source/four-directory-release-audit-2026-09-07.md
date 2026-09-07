# Four-directory release audit — 2026-09-07

用户拟开放 `development/`、`operator/`、`evolution/`、`principles/`。
初始审查形成下列检查和执行计划；用户随后授权落实并发布，执行记录见文末。
`experiments/` 不在本轮范围。

## Baseline and method

- Source: `main@36284d68007476ba30cc2ae48001b002bca6968c`。
- Public baseline: `release@d559b365f`；此前已核对 GitHub 默认分支为 `release`。
- 从 Git archive 提取四个目录，避免扫描范围混入 node_modules、dist、工作区残留或还原的 DVC 数据。
- 对全部 677 个跟踪文件运行 gitleaks，检查文件类型、符号链接、体积、DVC 指针、相对 Markdown 链接及隐私/第三方材料候选。
- 对关键词命中的协作、审稿、地址与 tape 内容进行定向复核；这不是逐句人工审阅，也不是完整法律、依赖漏洞或 Git 历史审计。
- 在隔离目录叠加公开版、四个候选目录和 `cmd/qcli/`，验证 Go 构建及 qcli 测试；前端从各自锁文件安装依赖。

## Findings

| Directory | Tracked files | Bytes | Assessment |
|---|---:|---:|---|
| `development/` | 544 | 3,074,271 | 可以准备开放；须明确历史记录、私有数据指针和公开文档依赖 |
| `operator/` | 64 | 455,355 | 可以准备开放；连同 `cmd/qcli/`，补实际使用入口 |
| `evolution/` | 36 | 289,688 | 未发现目录级保留理由；须解释私有历史提交引用 |
| `principles/` | 33 | 29,125 | 未发现目录级保留理由；须补治理文档入口 |

合计 3,848,439 bytes。凭据扫描退出码 0、命中 0。文件均为文本/JSON/源码/差异或空文件；没有二进制、压缩包、PDF、图片或符号链接。
未发现独立第三方邮件、整篇审稿意见或第三方图像/PDF。命中的审稿措辞为项目自身修改记录，不能据此认定有保密材料。
`package.json` 中的 `private: true` 防止包意外发布到包注册表，不妨碍源代码公开，应保留。

### Development records and data

`development/` 不只是设计文档：包含 33 个 `.log`、2 个 `.jsonl` 和 10 个 `.dvc`。
日志及 tape 已进入本次 Git 文件凭据扫描。两份 tape 的凭据关键词来自说明文字，没有记录凭据值；其中一份包含个人 Git 身份和本地路径。

需要定向处理的文件：

- `development/autopoiesis-probes/l0-quine-missionless-2026-06-27-post-restructure/run1-tape.jsonl`：Git 身份及本地工作区路径。
- 同目录 `run2-live-exit-on-5min-tape.jsonl`：本地路径。
- `development/open-source/public-branch-cut-2026-05-26.md` 和 `paper-history-rewrite-validation-2026-05-26.md`：具体私有存储位置。
- `development/frontmatter-object-graph-feasibility/readouts/graph_lib.py:19`：仓库根目录写死为作者机器路径，且依赖未公开的 `Paper/incubator/`。

普通历史路径、作者自己的研究草稿和旧检查输出不自动构成发布阻塞。
建议仅对公开副本中无研究必要的个人身份、具体存储位置进行明确标注的脱敏；保留原始证据，不把修订后的 tape 冒充原始记录。

10 个 DVC 指针位于两个 `l0-*` 探针目录。用户随后授权按同一凭据扫描和
manifest 流程发布；目录源码开放与证据数据开放仍是不同状态，只有明确列入
manifest 的对象可公开下载。

### Dependency and navigation gaps

qcli Go 入口 `cmd/qcli/main.go` 导入 `operator/qcli/app`，应和 `operator/` 一起发布。
`internal/qcli/` 已公开，不需要另建实现或改变 import 路径。
TUI 和 Web 依赖仓库内相邻的 client-core，保留完整 `operator/qcli/` 结构即可。

四目录内简单相对 Markdown 链接检查发现，拟议公开树仍有以下外部依赖：

- `Paper/`：80 处；`experiments/`：11 处。本轮不扩大这两个目录的范围，应在入口声明暂未公开的引用。
- 根文档：`DEVELOPMENT.md` 18 处、`EVOLUTION.md` 6 处、`CLOSURE.md` 4 处、`AGENTS.md` 2 处、`control-plane-registry.md` 2 处、`EVALUATION.md` 1 处。
- 尚未公开的 `scripts/` 目标 6 处，`prek.toml` 1 处。

这些是链接出现次数，不是唯一文件数，也不覆盖纯文本路径、代码字符串或历史提交。
建议把六份根文档纳入下一步同批审查；不要为了消除引用而开放整个 Paper 或 experiments。
演化记录及 qcli 设计文档中的私有提交哈希可以作为历史标识保留，但应说明公开仓库不能执行相应 `git show`。

`operator/qcli/README.md` 仍把已存在的 Web 实现称为 ready to execute，并以实现代理阅读顺序为主。
发布前补简短 clone/build/run 说明，更新实际状态；用现有 Go/Bun 命令即可，不必为此开放整份实验室 Makefile。

## Verification and repair

- 隔离树初始化本地 Git 仓库后，`go build ./...` 通过。未初始化时遇到获取 VCS 状态失败，不能将其误报为代码编译失败。
- 隔离树 `go test ./internal/qcli ./operator/qcli/app ./cmd/qcli` 通过。
- Bun 1.3.14：client-core 干净锁定安装、类型检查和 19 个测试通过。
- TUI 干净锁定安装、类型检查和 46 个测试通过。
- Web 初始 `bun install --frozen-lockfile` 失败：package.json 指定 DOMPurify 3.4.13，bun.lock 仍锁定 3.4.11。
- 在临时副本通过 Bun 更新锁文件，仅两个记录变化。更新后锁定安装、类型检查、26 个测试及生产构建通过；这份精确修复已带回 `operator/qcli/web/bun.lock`。

没有运行需要模型凭据的实时 agent 会话，没有验证开发目录下全部历史探针的可执行性，没有发布候选 DVC 数据。

## Proposed execution sequence

1. **准备内容。** 保留四个目录的作者文档、源码和已检查记录；处理上列有限的身份/路径项目，标注 development 的私有 DVC 数据状态。审查六份关联根文档，并更新 qcli 使用说明和状态。依赖 Paper 的历史原型明确其运行前提。
2. **统一发布规则。** 更新 `scripts/project-public.sh` 白名单，加入四目录及 `cmd/qcli/`；在 `scripts/check-public-projection.sh` 移除对应 denylist 项、补 required roots。将通过审查的根文档一起纳入。同步 `Paper/_design/migrations/public-branch-redesign.md`、本目录范围登记与根 README。保留现有 public/ 论文快照触发协议，不把源码开放绑定到论文发表状态。
3. **做精确候选树验收。** 从 Git 跟踪文件生成候选树，运行整树凭据扫描、既有公开清单 DVC 检查、Go 构建、qcli 测试和三套前端的干净锁定安装/类型检查/测试，执行 Web 构建。验证范围排除、DVC 清单未意外扩大和入口链接。未公开目录的已声明引用与意外断链分开报告。
4. **形成可发布提交。** 本地 release 做增量投影提交并记录公开树差异。完成审查后按发布协议推送 `public release`；不改写历史，不携带私有 Git bundle，不连带开放 experiments、Paper、profiles 或未检查 DVC 对象。

以上为初始审查结果，不将局部检查等同于最终整树验收。

## Execution — authorized 2026-09-07

- 四目录、`cmd/qcli/` 和六份关联根文档已纳入 v7 投影规则；根文档单独凭据扫描为 0 命中。
- 新增 `development/PUBLIC-RELEASE.md` 说明历史引用、私有数据指针与四份公开副本的脱敏。
- 四份原始记录保持原样；实际副本测试确认 JSON 事件顺序与非脱敏字段保留，源文件校验和不变。
- qcli 安装/运行/测试说明已补齐，Web 状态与锁文件已修正。
- 图读出工具按自身文件位置定位仓库，缺少私有 Paper 输入时明确失败；输入边界自检通过。
- `make check-control-plane` 通过；既有历史锚点、演化提醒等 advisory 保留，未扩展成本轮清理任务。
- 已从 `main@9a9f39f28` 生成并验证候选 `release@9d64a1b2d`：2,028 个文件、11 个公开成果入口，投影检查 `RESULT: PASS`。
- 初始候选 Git 树与既有 645 项公开清单关联的 DVC 数据凭据扫描通过。后续 DVC
  扩展以 manifest 候选树的同一检查为准，不复用此初始结果。
- 后续候选将 10 个 development 指针加入清单（共 655 项）。其中两份原始对象检出
  同一 API 令牌的 6 处副本，公开副本以 `*-public.dvc` 脱敏后重新扫描；原始对象不进入
  公开存储。
- 干净公开副本中普通 `go build ./...`、qcli Go 测试、三套前端锁定安装/类型检查、91 个测试及 Web 构建通过。
- 无模型调用的实际启动验证通过：qcli 帮助、HTTP 健康检查、Web HTML 和两个构建资源均可访问。
- 公开目录边界、入口文档链接和原始证据校验和通过；未公开路径按公开范围说明保留引用。
- 本记录的验收结果随最终增量投影发布；对外推送仍由操作步骤执行，不由投影脚本自动执行。
