# 0002 - 工具用 Go 编译为单静态二进制（推翻早期 Python 实现）

Date: 2026-09-29

## Status

Accepted

## Context

项目最初决策（D1 v1）用 Python 编写工具（自举、stdlib-only、3.8 兼容），分发走 PyPI/uvx + zipapp。通用多语言化（ADR 0001）后，"消费者环境必有 python3"的假设失效：写 Go/Rust 代码的 agent 沙箱可能既无 python 也无 uv。已核实 uvx 在无 python 时会自动下载受管解释器（`python-downloads` 默认 automatic），但前提仍是 uv 在场。grill 访谈中用户裁决：彻底消除运行时依赖，回到编译二进制路线；实现语言在 Go 与 Rust 之间选择。

## Decision

工具用 **Go** 实现，goreleaser 交叉编译为单静态二进制（linux/amd64+arm64、darwin/amd64+arm64、windows/amd64），全部数据集 `go:embed` 内嵌；MCP 服务器（官方 Go SDK `github.com/modelcontextprotocol/go-sdk/mcp`，v1.7+ 稳定）编译进同一二进制（`mg mcp`）。选择 Go 而非 Rust 的决定性理由：参考项目（Apache-2.0）的 goversion/cli/schema/featuresgen 代码与测试语义可直接改造复用——尤其 Go 检测器的 go.mod/go.work 边界语义（无 go 指令默认 1.16 等）已被上游用测试固化；交叉编译 trivial；`go install` 提供原生 pin 安装通道。

## Consequences

### Positive

* 零运行时依赖：任何 agent 沙箱（有 curl 或 go 任一即可）可用，运行时缺口问题彻底消除。
* 参考项目代码与测试语义直接复用（NOTICE 署名），Go 检测器风险最低。
* 三通道共享一个构建物：Releases 二进制、`go install@版本`、内置 MCP。
* MCP SDK 依赖编译期引入，早期"SDK 重依赖须隔离在 extra"的问题不复存在。

### Negative

* 自举叙事弱化：仓库主语言为 Go，"用现代 Python 写的 Python 准则工具"不再成立——Python 知识的载体只剩数据集（可接受：数据才是产品核心，工具是通用渲染器）。
* python/typescript 检测器需用 Go 重写（TOML/PEP 440/PEP 723/semver 解析），引入编译期第三方依赖（如 TOML 库）。
* 多平台构建矩阵与 goreleaser 维护。
* Python 社区贡献者门槛升高（数据贡献不受影响，检测器贡献受影响）。

## Alternatives Considered

* **Python 实现 + uvx/pyz/MCP（v3 方案）**：自举与贡献面好，但"无 uv 且无 python"沙箱的残余缺口被用户裁决为不可接受。
* **Rust 实现**：内存安全叙事好，但参考代码零复用、交叉编译配置重、`cargo install` 恰恰需要 Rust 工具链（违背零依赖初衷）、MCP SDK 生态较年轻。
* **hosted 远程 MCP endpoint**（补 Python 实现的缺口）：零本地依赖，但引入部署/域名/认证/运维，偏离静态数据集的简单性——保留为 post-MVP 通道。

## References

* docs/prd/PRD-0000-modern-guidelines.md（Decision D1 v4、D7、D10）
* docs.astral.sh/uv/concepts/python-versions/（uvx 自动下载 python 事实）
* github.com/modelcontextprotocol/go-sdk（官方 Go SDK，v1.7.0+，Apache-2.0）
* examples/go-modern-guidelines/（可复用代码来源，Apache-2.0，NOTICE 署名）
