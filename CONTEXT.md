# CONTEXT — modern-guidelines 领域词汇表

> 本文件是术语表，只记领域概念定义，不含实现细节。与 `docs/prd/PRD-0000-modern-guidelines.md`、`docs/adr/` 配套使用。

| 术语 | 定义 |
|---|---|
| **Language Hub（语言枢纽）** | 本项目的产品概念：单一入口（`mg` CLI / 单个 MCP 服务器）服务多语言的现代准则查询，语言是插拔单元而非独立工具。 |
| **Dataset（数据集）** | 一个语言的全部规则集合，`data/<lang>.json` 一份文件，含数据集头（language/axes/provenance）与规则数组。 |
| **Axis（版本轴）** | 数据集内的版本维度。单轴数据集（python、go）隐含语言版本轴；多轴数据集（typescript）显式声明 `axes: ["typescript","node"]`，每条规则标注所属轴。 |
| **Detector（检测器）** | 为一种语言解析项目版本的语言特异策略，输入文件路径，输出"轴→{版本, 来源, 保真度}"映射。 |
| **since_version** | 规则的最低适用版本（下界门控）：解析出的项目版本 ≥ since_version 时规则适用。语义是"代码必须兼容的最低版本"，故 `requires-python` 取下界而非上界。 |
| **autofix** | 规则可被哪个工具的哪条规则自动改写，格式 `tool:rule`（如 `ruff:UP007`、`gopls:modernize`）或 `null`。键必填——强制作者显式声明"可自动改写/不可"。 |
| **Provenance line（来源行）** | `list` 输出中标注版本解析结果的行：每轴一行，含版本、来源文件、保真度；advisory 或模糊来源附警示。多语言下保真参差的透明度补偿机制。 |
| **T1 / T2 / T3** | 版本解析三层：T1 = 显式 `--version` flag（精确）；T2 = 语言声明式 manifest（pyproject/go.mod/package.json 等，高保真）；T3 = 本地工具链兜底（`go version` 等，低保真，必附警示）。 |
| **Two-stage consumption（两段式消费）** | agent 的推荐使用模式：先 `list` 全量浏览（一行一条），再对相关规则 `explain` 取详情——省 token 且保证不漏。配套纪律：list 输出禁止截断；跳过疑似相关规则前必须 explain。 |
| **Drift test（漂移测试）** | 字节级快照测试：生成物（FEATURES 文档、ingestion 产物）必须与数据源同步，改数据不重新生成即测试失败。 |
| **Ingestion（吸收）** | 一次性导入外部数据集并转换为本仓 schema 的过程（当前唯一实例：JetBrains Go 54 条，Apache-2.0）。锁定上游 commit，上游发 minor 时重跑 + diff 评审——是快照分叉，不是持续 vendor。 |
| **Impact** | 规则出现频率分级：Critical（几乎每个项目都有）/ High（常见，5-20 处）/ Medium（规律出现，1-5 处）/ Low（罕见或特定场景）。跨数据集共用同一标尺。 |
