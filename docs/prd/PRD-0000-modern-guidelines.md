# modern-guidelines：面向 AI Agent 的多语言现代准则枢纽

> **Status**: Grilled | **PRD**: PRD-0000 | **Created**: 2026-09-28 | **Last updated**: 2026-09-29（v4，经 /grill 两轮访谈定稿：编译二进制形态）

## Goal

受 `examples/go-modern-guidelines/`（JetBrains）启发，构建**通用多语言枢纽 `modern-guidelines`**：Go 编译的单静态二进制 CLI `mg`（内嵌 MCP 服务器），agent 只给文件路径即自动推断语言与项目版本，返回版本门控的"用 X 代替 Y"现代惯用法。MVP 三语言：Python（自研 ≥63 条）+ Go（Apache-2.0 吸收 JetBrains 54 条）+ TypeScript（自研 ~40 条，TS 编译器/Node 双轴）。修复 agent 两大失效模式：训练截止后的新特性缺失、旧习语频率偏差。决策记录见 ADR 0001–0004。

## What I already know

参考项目（子代理 A/B 分析，2026-09-28）：

* 三层结构：数据（54 条 JSON、8 字段、newest-first 不变量）→ 微 CLI（~900 行 Go，list/explain 两段式，`go:embed`，init panic）→ 插件分发。goversion.go 251 行检测语义精细（go.mod 无 go 指令默认 1.16、go.work 默认 1.18 且不覆盖成员、不回落本地工具链）。FEATURES.md 纯生成 + 字节级漂移测试。skill 纪律：list 禁截断、跳过前先 explain。**Apache-2.0，代码与数据均可改造复用（需署名）**。

Python 生态（子代理 C）：3.14 当前稳定、3.15 定于 2026-10-01；3.9→3.15 现代化清单带 PEP 引用备齐；ruff UP 仅机械改写；pydantic 运行时注解下 UP006/UP007 不安全；PEP 686 将翻转"显式传 encoding"。

分发/MCP（子代理 D + 补充核实）：七家主流 agent（含 ZCode，zcode.z.ai/en/docs/mcp-services）支持 MCP；**官方 Go SDK `github.com/modelcontextprotocol/go-sdk/mcp` 稳定 v1.7.0+**，支持 stdio 与服务端工具、跟进 2026-07-28 规范；Claude Code 工具结果 10K/25K token 限制，~5KB list 输出安全；uvx `python-downloads` 默认 `automatic`（无 python 环境会自动拉解释器）——但用户仍选择编译二进制彻底消除运行时依赖（ADR 0002）。

各语言版本门版图（子代理 E，官方信源）：高保真 6 门（Go/Rust/C#/PHP/Ruby/Swift）、中保真 3 门（**TS/JS**：engines.node 仅 advisory、TS 版本仅能从 devDependencies typescript 范围推断；Java、Kotlin）；跨语言 mise/.tool-versions 在场即高保真但值可模糊（`ruby 3`）。语义陷阱：TS `target` 是 JS 输出级非编译器版本；Go `toolchain` 指令是建议非地板；Rust rust-version 不选工具链。

## Assumptions (temporary)

（grill 后全部已验证或转为需求，无遗留假设）

* ~~agent 环境必有 python3~~ → 已被编译二进制决策消除（ADR 0002）。
* AGENTS.md 是跨 agent 事实标准指令文件 → 已验证（Codex/Cursor/Claude Code/ZCode 官方支持）。
* 各语言规则集英文撰写；仓库文档中文 → 维持。
* JetBrains Go 数据一次性快照分叉 → 已转为需求（ingestion 锁定上游 commit，上游每发 minor 重跑 + diff 评审）。

## Open Questions

（无——3.15 占位转正与上游同步节奏为定时触发的既定待办，非开放问题）

* 待办（定时触发）：3.15 final（预计 2026-10-01）后核实 Python 占位规则；JetBrains 上游发 minor 时重跑 ingestion。

## Requirements

1. **多语言数据集** `internal/guidelines/data/<lang>.json`，每语言一份：数据集头 `{language, axes?, provenance?, rules}`——`axes` 为版本轴声明（>1 轴时必填，如 typescript 声明 `["typescript","node"]`）；`provenance` 外部吸收时必填（source/license/upstream commit/date）；规则数组 newest-first 为校验不变量（按轴分组内校验），id 语言内唯一。
2. **规则 schema**：`id / since_version / axis?（多轴数据集必填，单轴缺省）/ autofix（"tool:rule" 或 null，键必填）/ category / impact(Critical|High|Medium|Low) / guideline / details（机制+反例警示）/ references?（≤2 条短引用，可选——Go 吸收数据无引用）/ examples[{before[],after[]}] ≥1 组`。
3. **MVP 数据集三份**：Python ≥63 条（~55 惯用法 3.9→3.14 + 3.15 占位 + ~8 Tooling）；Go 54 条（ingestion 转换：`modernizer:true→autofix:"gopls:modernize"`，NOTICE 署名）；TypeScript ~40 条（双轴：TS 编译器轴收语言特性 satisfies/const 类型参数/using/NoInfer/推断型谓词等，Node 轴收 fetch/structuredClone/node:test 等——内容须经一轮完整官方信源研究后再定稿，遵守 D8 纪律）。
4. **语言识别**：文件扩展名映射（.py/.ts/.mts/.cts/.tsx→typescript/.go/.js/.mjs 预留…）+ `--lang` 覆盖；未知扩展名报错列出支持语言；`mg languages` 列出已装数据集及轴。
5. **版本解析三层**（检测器接口返回 `轴→{version, source, fidelity}` 映射）：
   - T1 `--version` flag（支持 `axis=ver` 逗号多值覆盖多轴）；
   - T2 语言声明式 manifest——python：pyproject `requires-python` 下界 / PEP 723 内联 / `.python-version` / mise；go：go.mod `go` 指令 / go.work（含默认值语义，改造参考实现）；typescript：`devDependencies.typescript` 语义版本范围下界 + Node 轴 `engines.node`/`.nvmrc`/mise；
   - T3 工具链兜底（`go version` / `node --version` / python 不可用则报错引导）。
   - **每次输出必须带 provenance（每轴一行）**：版本、来源、保真度；advisory/模糊来源（engines.node、mise `ruby 3` 式模糊值）附警示并建议显式 `--version`。
6. **`requires-python` 语义**：取下界为特性门；上界忽略；`==3.11.*`→3.11；`~=3.10`→3.10；无下界视为未声明继续解析链；无 requires-python 的 pyproject 不作答（monorepo 继续上溯）。TS 版本范围同理取下界（`^5.4`→5.4）。
7. **CLI**（`list` / `explain` / `languages` / `mcp` / `--version` / `help`）：`mg list --file-path src/foo.ts` 一次推断语言+各轴版本；list 一行一条 newest-first（多轴数据集按轴分组输出）；explain 缩进块含 Since/Summary/Details/Examples/References/Autofix；explain 接受 `lang:id` 或 `--lang`+裸 id；未知 id 报错附可用列表；list 版本来源互斥（同参考项目风格错误消息）。
8. **数据加载**：全部数据集 `go:embed` 进二进制；init 时全量校验，非法数据启动即 panic（对齐参考项目）。
9. **分发三通道**（ADR 0003）：① GitHub Releases 多平台静态二进制（goreleaser：linux/amd64+arm64、darwin/amd64+arm64、windows/amd64）；② `go install github.com/<org>/modern-guidelines@vX.Y.Z`（原生 pin 安装）；③ 内置 MCP：`mg mcp`（官方 Go SDK，stdio，只读工具 `list_guidelines(file_path?, language?, version?)` / `explain_guideline(language, ids)`）。
10. **触发（不依赖插件市场）**：AGENTS.md 语言无关片段（编辑前 `mg list --file-path <file>`、禁截断、视为权威、跳过前先 explain）+ MCP 工具描述；README 提供四家 MCP 配置即拷片段（command 指向 mg 二进制路径或 go install 安装位）。
11. **FEATURES 生成**：`docs/features/<lang>.md` 每语言一份 + 字节级漂移测试（对齐参考项目 TestFeaturesMarkdownInSync 模式）。
12. **贡献框架**：新语言接入 = 数据集 + 检测器 + manifest 研究记录（中保真语言强制，D8）；CONTRIBUTING 含接入指南。
13. **仓库卫生**：git init、Go module、LICENSE（Apache-2.0）、NOTICE（JetBrains 代码与数据署名）、CHANGELOG、CONTRIBUTING、Makefile（`test`/`generate-features`/`ingest`/`check`）、goreleaser.yml、CI（go test + 交叉构建 + tag 发布）、AGENTS.md 自举片段。GitHub 发布名 `modern-guidelines`。

## Acceptance Criteria

* [ ] schema 校验拒绝：非法 id/重复 id/非点分数版本/违反 newest-first/缺失 autofix 键/impact 越枚举/空字段/无示例/多轴数据集规则缺 axis 字段/axis 未在数据集 axes 声明/数据集头缺 language/外部吸收缺 provenance。
* [ ] 语言识别：扩展名映射正确（含 .mts/.cts/.tsx→typescript）；未知扩展名报错列支持语言；`--lang` 覆盖生效。
* [ ] Python 检测器：全链优先级、find-up 嵌套、PEP 723、`>=3.10,<3.13` 取 3.10、`==3.11.*`、无 requires-python 上溯、T3 引导。
* [ ] Go 检测器：go.mod `go` 指令（含无指令默认 1.16）、go.work（默认 1.18、不覆盖成员）、`toolchain` 指令仅提示不作地板、`go version` 兜底——行为对齐参考 goversion 测试语义。
* [ ] TS 检测器双轴：`"typescript": "^5.4"`→5.4；`engines.node ">=18"`→18 且 provenance 标注 advisory；`.nvmrc` 读取；`node --version` 兜底。
* [ ] provenance：每轴一行输出版本+来源+保真度；advisory/模糊来源带警示。
* [ ] `list`：三语言各自 newest-first 过滤正确（多轴分组）；`explain` 全格式渲染；未知 id 报错附列表。
* [ ] ingestion：JetBrains 54 条零丢失转换、字段映射正确、NOTICE 署名、漂移测试锁结果。
* [ ] FEATURES×3 漂移测试：改 data/<lang>.json 不重新生成则失败。
* [ ] 二进制：三平台（linux/darwin/windows）交叉构建成功且 `mg list --file-path x.py` 冒烟通过；`go install ...@v0.1.0` 可装可运行；单文件无运行时依赖。
* [ ] `mg mcp`：tools/list 返回两工具、tools/call 与 CLI 同构（SDK 集成测试）。
* [ ] Python ≥63 条、TS ≥40 条（双轴分布）过校验；3.9–3.14 每版本 ≥3 条；pydantic/运行时注解相关规则含 caveat。
* [ ] dogfood：`mg list --file-path main.go` 在本仓库正确解析自身 go.mod；ZCode 实测 .py/.go/.ts 三文件端到端。

## Definition of Done

* `go test ./...` 全绿；goreleaser 构建三平台成功；CI 绿；tag 即出 Releases + 可 go install。
* README/CONTRIBUTING（新语言接入指南）/CHANGELOG/NOTICE 齐备；FEATURES×3 生成物 + 漂移测试在位。
* 自举：工具检测自身 go.mod 的 dogfood 测试在 CI 常跑。
* 风险与回滚：纯新增；发布物损坏重跑 goreleaser。

## Out of Scope

* Python/uvx/zipapp/npm 分发通道（随编译二进制决策作废，ADR 0002/0003）；各平台插件/skill manifest（post-MVP 可选薄 shim）。
* Rust/C#/PHP/Ruby/Swift/Java/Kotlin 数据集与检测器（框架就绪，按贡献流程逐语言接入）。
* 独立 node 数据集（Node 惯用法已并入 typescript 数据集 node 轴）。
* ruff/gopls/tsc 的自动执行或修复（只告知不执行）。
* 3.15 正式规则内容（占位）；反 LLM 坏模式/质量规则；JSON/markdown 输出模式、i18n。
* hosted 远程 MCP endpoint（post-MVP 通道，ADR 0002 预留）。

## Technical Approach

### 总体

Go 编写（ADR 0002），单静态二进制（`go:embed` 内嵌三语言数据集 + goreleaser 交叉构建），检测器/数据集注册表插拔。参考项目（Apache-2.0）的 goversion/cli/schema/featuresgen 直接改造复用（NOTICE 署名）。规模估算：~1500 行 Go 源码 + ~1300 行测试 + 三份数据集 JSON。

### 目录结构

```
modern-guidelines/
├── go.mod  goreleaser.yml  Makefile  LICENSE  NOTICE  AGENTS.md
├── README.md  CONTRIBUTING.md  CHANGELOG.md
├── docs/features/{python,go,typescript}.md        # 生成物
├── .github/workflows/          # ci.yml；release.yml（tag → goreleaser）
├── main.go
├── internal/
│   ├── cli/                    # ~200：子命令、语言识别、provenance 渲染
│   ├── schema/                 # ~150：数据集+规则+轴校验（newest-first 按轴分组）
│   ├── registry/               # ~100：数据集注册、扩展名映射、检测器注册表
│   ├── guidelines/             # ~200：embed、加载、过滤、索引
│   │   └── data/{python,go,typescript}.json
│   ├── detectors/
│   │   ├── python.go           # ~280：TOML 解析、PEP 440 下界、PEP 723、mise
│   │   ├── go.go               # ~230：改造参考 goversion（go.mod/go.work 语义）
│   │   └── typescript.go       # ~250：TS 范围下界 + Node 轴双解析
│   ├── mcpserver/              # ~150：官方 Go SDK（github.com/modelcontextprotocol/go-sdk/mcp）
│   └── featuresgen/            # ~150：per-language 渲染 + go:generate
├── scripts/ingest-jetbrains-go/ # ingestion 程序（上游 commit 锁定）
└── internal/**/*_test.go       # 表驱动测试（策略对齐参考项目）
```

### 多轴数据集设计（grill Q5 用户裁决：TS+Node 双轴）

数据集声明 `axes`；规则携带 `axis`（多轴时必填）；检测器返回 `轴→{version, source, fidelity}`；过滤按规则所属轴匹配；provenance 每轴一行。示例：typescript 数据集同时含 `"axis":"typescript","since_version":"5.4"` 的 NoInfer 规则与 `"axis":"node","since_version":"18"` 的 fetch 规则。engines.node 为 advisory：provenance 标注并建议显式 `--version node=20`。

### MCP 通道（`mg mcp`）

官方 Go SDK v1.7+（stdio、服务端工具、跟进 2026-07-28 规范），编译进二进制——无 extra、无运行时依赖，SDK 重量问题随编译消除。`list_guidelines(file_path?, language?, version?)` 单 file_path 参数同时推断语言与各轴版本。工具描述承担 config-time 触发。~5KB 输出低于 Claude Code 10K token 警告线。

### 诚实 caveat（记录在案）

MCP 主动触发可靠性不保证优于 AGENTS.md 显式指令——两者并存。TS/Node 检测中保真（advisory/范围推断）——provenance 警示补偿（D8）。ingestion 与上游分叉——commit 锁定 + minor 周期重跑。TS 内容必须先研究后定稿（D8 纪律，实现期子代理研究 typescriptlang.org whatsnew 与 nodejs.org release notes）。AGENTS.md 触发需用户复制一次片段。

### 测试策略

表驱动测试（对齐参考项目：错误消息子串断言、注入 io.Writer、failingWriter、t.TempDir/t.Chdir 构造嵌套 manifest/PEP 723/go.work/多轴场景、环境变量保存恢复）；Go 检测器测试语义对齐参考 goversion_test 边界集；字节级漂移测试 ×4（FEATURES×3 + ingestion 产物）；MCP 用 Go SDK 内存客户端冒烟；dogfood：CI 内 `mg list` 自检本仓库 go.mod。

## Research References

* （实现期建议归档 `/research` 记录：①Python 版本门控现代化清单（子代理 C）；②各语言版本门 manifest 版图（子代理 E，随语言接入持续扩展）；③TypeScript/Node 内容研究（待做，PR7 前置）。）

## Feasible Approaches

**Approach A: 通用枢纽 + Go 编译二进制** (已采纳，ADR 0001/0002)

* 单二进制零运行时依赖，任何 agent 沙箱可用；参考代码直接复用；go install + Releases + 内置 MCP 三通道；数据集/检测器插拔。
* Cons: 仓库主语言为 Go（与"Python 准则"主题的自举叙事弱化——数据集才是 Python 知识载体）；多平台构建矩阵；跨语言维护者门槛。

**Approach B: Python 实现 + uvx/pyz/MCP**（v3 方案，被 grill Q1 裁决否决）

* 自举、Python 社区贡献面；但残余运行时缺口（无 uv 无 python 的沙箱）不可接受。

**Approach C: 单语言工具族 / federate**（早期分析否决）

* N 个入口违背"统一调用"的核心诉求。

## Decision (ADR-lite)

（grill 定稿，2026-09-29；重大裁决已升格正式 ADR）

* **D1（v4，推翻 v1）工具实现 = Go 编译单静态二进制** → ADR 0002。多语言化后"环境必有 python"假设失效；uvx 自动拉解释器仍依赖 uv 在场；用户裁决彻底零依赖。代价：自举叙事弱化、python/ts 检测器用 Go 重写（Go 检测器反而直接复用参考实现）。
* **D2（v4）形态 = 通用多语言枢纽**（单 CLI/单 MCP、dataset+detector 插拔）→ ADR 0001。
* **D3 内容范围**：Python ≥63 + Go 54（吸收）+ TypeScript ~40（双轴）。
* **D4 schema v2**：`autofix: tool:rule`、`references` 可选、多轴（axes/axis）、数据集头 provenance。
* **D5 PRD 落盘**：本文件（grill 时随项目更名重命名为 PRD-0000-modern-guidelines.md）。
* **D6 触发**：AGENTS.md 片段为主 + MCP 工具描述为辅。
* **D7（v4）MCP = 官方 Go SDK 编译进二进制**（v2 的 extra/重量问题随编译消除）。
* **D8 检测三层 + provenance**：多语言保真参差用透明度补偿；中保真语言接入前强制 manifest 研究记录。
* **D9 Apache-2.0 吸收 JetBrains Go 数据**：ingestion 锁上游 commit + NOTICE + 漂移测试 → ADR 0004。
* **D10（v4）分发 = GitHub Releases 二进制 + go install + 内置 MCP**（PyPI/uvx/zipapp 作废）→ ADR 0003。
* **D11 命名 = modern-guidelines / mg**（PyPI 已核实可用——虽不再发 PyPI，名称全域无冲突仍有益；`mg` 与老编辑器重名低风险，接受）。
* **D12 TypeScript 双轴**（grill Q5 用户裁决，与推荐相反，记录在案）：typescript+node 双轴并入同一数据集；schema 加 axis 字段；检测器双解析；advisory 轴的门控可信度折扣由 provenance 警示补偿。

**Consequences**: Go 仓库主语言（Python 知识在数据集不在工具代码）；TS 中保真检测的长期改进空间；ingestion 分叉维护义务；`mg` 重名接受；远程 MCP endpoint 预留（ADR 0002）。

## Implementation Plan (small PRs)

* **PR1 脚手架 + schema v2**：git init、go mod、goreleaser.yml、LICENSE/NOTICE、Makefile、包骨架、schema 多轴校验 + 全拒绝用例。
* **PR2 注册表 + CLI + embed**：registry（扩展名/数据集/检测器）、cli（list/explain/languages/version/help + provenance 渲染）、三语言各 5 条样例、输出测试。
* **PR3 Python 检测器**：TOML 库选型（编译期依赖）、PEP 440 下界解析、PEP 723、mise、边界用例。
* **PR4 Go 检测器 + ingestion**：改造参考 goversion（语义+测试对齐）；ingest-jetbrains-go（54 条转换 + NOTICE + 漂移测试）。
* **PR5 TypeScript 检测器**：双轴解析（TS 范围下界 + engines.node/.nvmrc/mise/node --version）、advisory 警示。
* **PR6 Python 内容填充**（子代理并行三批，引用与 caveat 把关）：≥63 条。
* **PR7 TS 内容研究与填充**（先研究后定稿，D8）：子代理对照 typescriptlang.org whatsnew（4.9→5.x）与 nodejs.org release notes 产出 ~40 条双轴规则。
* **PR8 生成器 + 文档**：featuresgen×3 + 漂移；README（三通道 + 四家 MCP 配置 + AGENTS.md 片段）、CONTRIBUTING（新语言接入指南）、CHANGELOG、AGENTS.md 自举。
* **PR9 发布流水线**：CI（go test + 交叉构建 + dogfood 自检）、release.yml（goreleaser）、go install 验证。
* **PR10 MCP + dogfood**：mcpserver（Go SDK）+ 集成测试；ZCode 实测 .py/.go/.ts 端到端。

## Technical Notes

* 子代理 A/B（参考项目）：main 15/cli 174/guidelines 210（go:embed + init panic + O(1) id map）/schema 110（校验链：id 字符集→去重→版本→newest-first→modernizer 必填→category→impact→非空→≥1 示例）/goversion 251（flag→文件自答→find-up go.mod>go.work→`go env` 兜底；无 go 指令默认 1.16、go.work 默认 1.18 且不覆盖成员、不回落工具链）/featuresgen 125+字节级同步测试/run-tool.sh 版本化缓存；54 规则×14 类、峰值 1.21、details 30–90 词"机制+caveat"、JSON 零 URL、双版本流。Apache-2.0 已核实。
* 子代理 C（Python）：devguide 版本线（3.14 stable、3.15→2026-10-01）；3.9–3.15 清单（PEP 584/616/585、634/604/618、654/680/673/655/675/678、695/692/701/698、742/702/703/744/696、750/649/749/779/734/758/784、686/810/814/661/798）；typing 弃用时间线；UNVERIFIED 已剔（IBM clint、subprocess timeout、pyright2）。
* 子代理 D（分发/MCP）：七家 MCP 配置形态（ZCode：~/.zcode/cli/config.json 的 mcp.servers + .agents/mcp.json 兼容读取，zcode.z.ai/en/docs/mcp-services）；uvx extras+pin 与版本缓存；Claude Code 10K/25K token；**官方 Go SDK v1.7.0+（github.com/modelcontextprotocol/go-sdk/mcp，stdio+AddTool，跟进 2026-07-28 规范，Apache-2.0）**；uv python-downloads 默认 automatic（docs.astral.sh/concepts/python-versions）。
* 子代理 E（版本门版图）：高保真 Go（go.mod go 指令=强制地板+语言版本、toolchain 指令仅建议）/Rust（rust-version=MSRV 可选、rust-toolchain.toml 才选工具链）/C#（TargetFramework 隐含 C# 版+LangVersion）/PHP（require.php 地板+config.platform 覆盖）/Ruby（Gemfile ruby 精确）/Swift（Package.swift 首行）；中保真 TS（engines.node advisory、TS 版本=devDependencies 推断、target≠编译器版本）/Java/Kotlin；mise/.tool-versions 跨语言（值可模糊）。UNVERIFIED：Volta（项目停维）、.java-version。
* grill 裁决记录（2026-09-29）：Q1 编译二进制（推翻 D1，ADR 0002）、Q2 modern-guidelines/mg、Q3 三语言 MVP、Q4 Go（ADR 0002）、Q5 TS+Node 双轴（与推荐相反，D12 记录）；uvx-auto-python 与 PyPI 可用性事实已核实并纳入裁决材料；8 项敌意发现全闭（4 项自查裁定：上游同步节奏/references 可选/mise 模糊值/AGENTS.md 假设验证）。
* 关键 caveat 素材：pydantic 运行时注解 UP006/UP007 不安全（ruff 官方）；PEP 686 翻转 encoding 建议；Go per-file `//go:build go1.22` 抬升单文件语言版本（检测器不处理，Go 数据集 details 已含适用边界）。
* 通用性结论：`{数据集(+轴) + manifest→版本门控 + list/explain + 二进制/go install/内置 MCP + AGENTS.md/MCP 触发}` 完全语言参数化；9 门主流语言全部有可静态读取版本来源（6 高 3 中）；瓶颈在内容 editorial 而非工程。

## Traceability

- **Created by**: `/think` (2026-09-28)
- **Revised by**: v2 分发/触发解耦；v3 通用枢纽（2026-09-28）；v4 grill 定稿（2026-09-29）
- **Grilled by**: `/grill`（completed 2026-09-29）— 两轮 5 问全决（Q1 编译二进制/Q2 命名/Q3 三语言/Q4 Go/Q5 双轴）；8 项敌意发现全闭；ADR 0001–0004 创建；CONTEXT.md 词汇表建立；PRD 随项目更名
- **Sliced into**: (pending `/story`)
- **New terms**: 见 CONTEXT.md（dataset/axis/detector/since_version/autofix/provenance line/T1-T3/two-stage consumption/drift test/ingestion/language hub）
- **New decisions**: D1–D12；其中 4 项已升格 ADR 0001–0004

## Issue

（缺失——/think 未创建父 Issue：工作区无 git 仓库与 issue tracker。选项：(a) 实现启动时 git init + /setup-project 后手动补建父 Issue 并回填；(b) 继续无父 Issue（/story 子 Issue 将独立存在）。）
