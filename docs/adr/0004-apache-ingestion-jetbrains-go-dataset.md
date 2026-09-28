# 0004 - 按 Apache-2.0 吸收 JetBrains Go 数据集（一次性快照）

Date: 2026-09-29

## Status

Accepted

## Context

通用枢纽（ADR 0001）需要 Go 语言的高质量规则集。JetBrains/go-modern-guidelines 的 54 条规则（含 gopls modernize 标记、机制+caveat 式 rationale、before/after 示例）正是所需的 editorial 质量，且其仓库为 **Apache-2.0**（已核对本地副本 LICENSE）。选项：从零自研（≈一轮完整研究 + 撰写，数周量级）、持续 vendor 上游（紧耦合）、一次性快照吸收（independent fork，commit 锁定）。

## Decision

一次性吸收（ingestion）：脚本读取上游 `guidelines.json`，转换为本仓 schema v2（`modernizer:true → autofix:"gopls:modernize"`、`false → null`；其余字段 1:1；references 留空后补），产出 `data/go.json`，由字节级漂移测试锁定。provenance 记录来源仓库、许可、**上游 commit** 与转换日期；仓库根 NOTICE 履行 Apache-2.0 署名义务（数据与改造复用的代码一并署名，见 ADR 0002）。同步策略：上游每发 minor 版本时重跑 ingestion + 人工 diff 评审——是快照分叉，不追 commit 级同步。

## Consequences

### Positive

* MVP 免费获得 54 条高质量 Go 规则（含 1.0→1.27 全版本梯度），Go 语言即刻满保真。
* 字段映射是纯机械转换，漂移测试保证零丢失、可复现。
* commit 锁定使数据来源完全可追溯。

### Negative

* 与上游分叉：上游新增/修订规则不会自动流入，需按 minor 节奏人工同步；长期可能累积偏差。
* schema 演进（如新增字段）时旧快照需重新映射——由漂移测试兜底。
* 上游若改许可证或停止维护，快照仍可用但来源老化。

## Alternatives Considered

* **从零自研 Go 规则集**：无分叉问题、风格统一，但数周 editorial 成本换不来差异化（上游已足够好）。
* **持续 vendor / git subtree 上游**：自动跟进，但上游 schema 与本仓 schema v2（axis/autofix/provenance）已分叉，每次同步都要转换评审，且受上游发布节奏与意外变更牵制。
* **运行时直接读上游 JSON（不落盘）**：引入网络依赖与运行时信任问题，违背单二进制自包含。

## References

* docs/prd/PRD-0000-modern-guidelines.md（Decision D9、Requirement 3/13）
* internal/testdata/upstream/（ingestion 消费的上游快照位置；完整参考项目保留为本地 examples/，不入库）
* apache.org/licenses/LICENSE-2.0（许可原文；署名与 NOTICE 义务见根目录 NOTICE）
