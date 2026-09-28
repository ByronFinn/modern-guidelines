# 0003 - 分发与触发解耦：Releases + go install + 内置 MCP，不依赖插件市场

Date: 2026-09-29

## Status

Accepted

## Context

参考项目用插件市场分发（5 平台 manifest + skill 内 bootstrap 脚本按版本 go install 到缓存目录），把两件事捆绑在一起：**分发**（工具位如何到达机器）与**触发**（agent 如何知道该用它 + 工作流纪律：list 禁截断、跳过前先 explain）。插件与平台绑定，用户明确否决此路线，要求平台无关（MCP/CLI/二进制）。分发通道候选：uvx/PyPI（随 ADR 0002 作废）、npm（生态错配）、GitHub Releases 二进制、go install、MCP。

## Decision

分发与触发作为两个独立问题分别解决：

**分发三通道**（全部平台无关，共享同一构建物）：
1. GitHub Releases 多平台静态二进制（goreleaser）；
2. `go install github.com/<org>/modern-guidelines@vX.Y.Z`（原生版本 pin 安装）；
3. 内置 MCP stdio 服务器 `mg mcp`（官方 Go SDK 编译进二进制）。

**触发两形态**：
1. AGENTS.md 语言无关触发片段（跨 agent 事实标准，Codex/Cursor/Claude Code/ZCode 均读取）——主触发，承载两段式消费纪律；
2. MCP 工具描述（config-time 触发）——辅助。

各平台插件/skill manifest 降级为 post-MVP 可选薄 shim（仅当市场用户提出）。README 提供四家 MCP 客户端配置即拷片段。

## Consequences

### Positive

* 平台无关：任何能跑二进制或连 MCP 的 agent 可用，无市场审核/manifest 维护负担。
* 通道全部零安装失败面（curl 下载 / go install / 已配置的 MCP command）。
* 触发片段随仓库版本化，纪律（禁截断、explain-before-skip）显式可审计。
* 人类与 CI 也能直接用同一 CLI。

### Negative

* AGENTS.md 触发需用户一次性复制片段（插件市场是零操作安装）——README 降低摩擦。
* MCP 工具被动躺在工具列表中，主动触发可靠性不保证优于显式 AGENTS.md 指令——两形态并存而非互斥。
* 失去插件市场自带的发现性（搜索/推荐）渠道。

## Alternatives Considered

* **插件市场分发（参考项目原样）**：市场负责版本流与零操作安装，但平台绑定（用户否决）且每平台 manifest 维护。
* **npm/npx 通道**：仅当工具为 JS 实现才合理；Go 二进制经 npm 包装徒增一层。
* **hosted 远程 MCP 作为主通道**：一次配置全环境可用，但部署/认证/运维成本，留 post-MVP。

## References

* docs/prd/PRD-0000-modern-guidelines.md（Decision D6、D10）
* 七家 agent 的 MCP 配置形态核实记录（PRD Technical Notes，子代理 D；含 ZCode zcode.z.ai/en/docs/mcp-services）
* examples/go-modern-guidelines/plugin/（被否决的插件分发形态对照）
