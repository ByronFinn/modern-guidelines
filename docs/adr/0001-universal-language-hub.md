# 0001 - 通用多语言枢纽（而非每语言独立工具）

Date: 2026-09-29

## Status

Accepted

## Context

参考项目 JetBrains/go-modern-guidelines 只服务 Go 一门语言。本项目最初也定位为 Python 单语言工具，但用户明确要求"所有语言一个入口，agent 按语言调用即可获取对应版本准则"。两条路线：每语言一个独立工具（共享 schema 契约），或单一工具内做多语言（数据集+检测器插拔）。版本门版图研究（对照 9 门主流语言官方文档）证实：所有主流语言都存在可静态读取的项目版本声明来源（Go/Rust/C#/PHP/Ruby/Swift 高保真，TS/Java/Kotlin 中保真），通用检测在技术上无结构性障碍。

## Decision

构建通用多语言枢纽 `modern-guidelines`：单一 CLI（`mg`）与单一 MCP 服务器；一门语言 = 一个数据集（`data/<lang>.json`，含可选多版本轴）+ 一个检测器（`detectors/<lang>.go`），核心（schema/注册表/加载/过滤/CLI/MCP/文档生成）语言无关。MVP 三语言：Python（自研 ≥63 条）、Go（Apache-2.0 吸收 JetBrains 54 条，见 ADR 0004）、TypeScript（自研 ~40 条，TS 编译器/Node 双轴）。其余语言经贡献框架接入，中保真语言接入前强制先完成 manifest 研究记录。

## Consequences

### Positive

* agent 单一入口：`mg list --file-path x` 自动推断语言+版本，无需知道装了哪个工具；一次 MCP 配置覆盖全部语言。
* 工程一次成型：加语言 = 数据集 + 检测器，核心零改动。
* 通用检测能力（三层解析 + provenance 输出）对所有语言复用。
* 多语言数据集互相校准（同一 impact 标尺、同一两段式消费纪律）。

### Negative

* 核心 schema 须容纳语言差异（版本轴、autofix 通用化），比单语言略复杂。
* 检测保真参差（TS/Java/Kotlin 中保真）——需 provenance 透明度补偿。
* 单一发布节奏绑定所有语言的数据更新。
* 内容 editorial 成本随语言数线性增长，是真正的扩展瓶颈（每语言 ≈ 一轮完整官方信源研究）。

## Alternatives Considered

* **每语言独立工具族**：生态原生分发（go install/cargo/npm），但 agent 需知道调哪个工具，N 份 CI/发布/维护，无法单一 MCP 入口——被用户核心诉求直接否决。
* **数据仓库 + 各官方工具互认（federate）**：契约靠约定无强制校验，仍是 N 个入口，被本方案严格支配。
* **维持 Python 单语言工具**：与用户三轮演进方向相悖。

## References

* docs/prd/PRD-0000-modern-guidelines.md（Decision D2）
* 参考项目 examples/go-modern-guidelines/（单语言形态对照）
* 版本门版图研究记录（PRD Technical Notes，子代理 E）
