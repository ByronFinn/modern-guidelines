# modern-guidelines (`mg`)

**版本门控的多语言现代准则枢纽，为 AI agent 而建。** `mg` 是一个 Go 编译的单静态二进制（内嵌 MCP 服务器）：agent 给一个文件路径，`mg` 自动推断语言与项目版本，返回该版本下适用的"用 X 代替 Y"现代惯用法。它修复 agent 的两大失效模式——**训练截止后的新特性缺失**（如 Python 3.12 的 `type` 语句、Go 1.22 的 `range over int`）与**旧习语的频率偏差**（训练语料里旧写法出现得更多，agent 于是更倾向输出旧写法）。当前内置 Python、Go、TypeScript 三套数据集；数据集插拔，语言可按[贡献指南](CONTRIBUTING.md)逐语言接入。

> **语言说明**：本仓库的文档（README/CONTRIBUTING/ADR 等）用中文撰写；各语言规则数据集本身用英文撰写（`guideline`/`details` 等字段是给 agent 消费的准则文本，保持英文）。

## 当前状态（诚实声明）

以下行为以当前代码为准，均已实测核实：

- **`--file-path` 已接线进主二进制**：`main.go` 空白导入 `internal/detectors`，三语言检测器随启动注册，语言与各轴版本直接从项目 manifest 解析，无需 `--version` 兜底。实测 `go run . list --file-path main.go` 输出首行为 `# go 1.26 (source: go.mod go directive, fidelity: high)`（三层优先级见[版本检测](#版本检测三层--provenance)）。
- **`mg mcp` 是真实的 MCP stdio 服务器**：实现在 `internal/mcpserver`，基于官方 Go SDK `github.com/modelcontextprotocol/go-sdk v1.8.0`，提供 `list_guidelines` 与 `explain_guideline` 两个工具，工具输出与对应 CLI 子命令一致（仅末尾换行被裁去）。
- **首个 Release 尚未发布**：GitHub Releases 下载与 `go install` 通道需等 `v0.1.0` tag 发布后可用。

## 快速开始（三通道）

### 通道 1：GitHub Releases 下载二进制

从 Releases 页下载对应平台的静态二进制（linux/darwin/windows，amd64/arm64），无任何运行时依赖：

```console
# 占位链接：v0.1.0 发布后生效
https://github.com/ByronFinn/modern-guidelines/releases
```

下载后把 `mg` 放进 `PATH` 即可。

### 通道 2：go install

需要本地 Go 工具链：

```console
go install github.com/ByronFinn/modern-guidelines@v0.1.0
```

> 注意：`@v0.1.0` 需等首个 tag 发布（见"当前状态"）。模块路径 org 已定（`ByronFinn`），不会再调整。

### 通道 3：MCP

`mg mcp` 启动内置的 stdio MCP 服务器（同一二进制，无额外安装）。工具与 CLI 同构：

- `list_guidelines(file_path?, language?, version?)` —— 单个 `file_path` 同时推断语言与各轴版本
- `explain_guideline(language, ids)` —— 取规则详情

`mg mcp` 即装即用（见上方"当前状态"）。[四家客户端配置见下文](#mcp-配置)。

## CLI 速览

```console
$ mg languages
go (axes: go)
python (axes: python)
typescript (axes: typescript, node)
```

`list` 给一个语言在其解析版本下适用的全部规则，一行一条、newest-first；输出首行是 provenance 行（版本来源与保真度）：

```console
$ mg list --version python=3.12
# python 3.12 (source: --version flag, fidelity: high)
type_alias_statement: Declare type aliases with the `type` statement instead of `TypeAlias` assignments.
generic_class_params: Declare generic classes inline (`class Box[T]:`) instead of module-level TypeVars plus `Generic[T]`.
...
```

多轴数据集（如 typescript 的 typescript+node 双轴）按轴分组，每轴一行 provenance：

```console
$ mg list --version typescript=5.4,node=18
# typescript 5.4 (source: --version flag, fidelity: high)
# node 18 (source: --version flag, fidelity: high)
# axis: typescript
no_infer: Exclude a parameter from inference with `NoInfer<T>` instead of sentinel defaults or `never` tricks.
...
# axis: node
...
```

`explain` 取一条规则的完整详情（机制、反例警示、before/after 示例、引用、autofix 映射）：

```console
$ mg explain python:dict_merge_operator
dict_merge_operator:
  Since: python 3.9

  Summary:
    Merge and update dicts with the `|` and `|=` operators instead of ...

  Details:
    PEP 584 operators: `a | b` returns a new dict with the right operand's
    keys winning; ...

  Examples:
  Before: ...
  After: ...

  References:
    PEP 584

  Autofix:
    none
```

其他已核实的行为细节：

- `list` 的版本来源互斥：`--version`、`--file-path`、位置参数路径三选一，组合即报错。
- 多轴数据集用 `--version` 时必须覆盖全部轴（如 `--version typescript=5.4,node=18`），缺轴报错并给出补全示例。
- 未知扩展名报错并列出支持语言；`--lang` 可覆盖按扩展名的语言识别（不能与 `--version` 组合）。
- `explain` 接受 `lang:id` 或 `--lang <language>` 加裸 id；未知 id 报错并列出该语言全部可用 id。
- `mg version` / `mg --version` 打印版本（未发布构建为 `dev`）。

## 版本检测三层 + provenance

版本解析按三层优先级（术语见 [CONTEXT.md](CONTEXT.md)）：

| 层 | 来源 | 保真度 |
|---|---|---|
| **T1** | 显式 `--version axis=ver`（多轴逗号分隔） | 精确 |
| **T2** | 语言声明式 manifest：python 的 pyproject `requires-python` 下界 / PEP 723 内联 / `.python-version` / mise；go 的 go.mod `go` 指令 / go.work；typescript 的 `devDependencies.typescript` 范围下界 + node 轴的 `engines.node` / `.nvmrc` / mise | 高（声明式）|
| **T3** | 本地工具链兜底：`go version` / `node --version`；python 不可用则报错引导 | 低（本机 ≠ 项目）|

每次 `list` 输出**每轴一行 provenance**：版本、来源、保真度。advisory 或模糊来源（如 `engines.node`、mise 的 `ruby 3` 式模糊值）会被标注 `medium` 保真度并附警示，建议显式 `--version` 钉住。警示行的实际渲染格式（由 `internal/cli` 渲染逻辑输出、测试锁定）：

```console
# node 18 (source: package.json engines.node, fidelity: medium) [advisory: pass --version node=<version> to pin it explicitly]
```

多语言下保真度参差（Go 高保真、TS/Node 中保真）用这种透明度补偿——见 PRD 决策 D8。

## AGENTS.md 触发片段（即拷即用）

把下面片段加入你仓库的 `AGENTS.md`（Codex / Cursor / Claude Code / ZCode 等均读取该文件），即可让 agent 在编辑代码前自动查询准则：

```markdown
## modern-guidelines (mg)

Before editing any source file in this repository, run:

    mg list --file-path <file>

- Read the COMPLETE output; never truncate or sample it.
- Treat the listed guidelines as authoritative: they take precedence over
  older conventions found elsewhere in this repository and over your own
  training priors.
- Before skipping any guideline that seems relevant, run:

    mg explain <lang>:<id>

- If mg cannot resolve a version from the project's manifests, pin it
  explicitly, e.g.:

    mg list --version python=3.12
    mg list --version typescript=5.4,node=18
    mg list --version go=1.26
```

## 支持语言

| 语言 | 规则数 | 版本跨度 | 说明 |
|---|---|---|---|
| python | 88 | 3.9 → 3.15 | 自研；3.15 条目为占位（3.15 final 预计 2026-10-01，落地后核实转正） |
| go | 54 | 1.0 → 1.27 | 吸收自 [JetBrains/go-modern-guidelines](https://github.com/JetBrains/go-modern-guidelines)（Apache-2.0，commit 锁定快照，署名见 [NOTICE](NOTICE)） |
| typescript | 57 | typescript 轴 33 条（4.9 → 7.0）+ node 轴 24 条（16 → 24） | 自研双轴：TS 编译器轴收语言特性，node 轴收运行时 API（决策 D12） |

`.py`、`.go`、`.ts`、`.mts`、`.cts`、`.tsx` 按扩展名识别（`.js`/`.mjs` 已预留映射到 javascript，暂无数据集）。

## MCP 配置

四家 MCP 客户端的即拷配置，`command` 均指向 `mg` 二进制（按安装方式替换为绝对路径），`args` 均为 `["mcp"]`：

**Claude Code**（项目根 `.mcp.json`）：

```json
{
  "mcpServers": {
    "modern-guidelines": {
      "command": "mg",
      "args": ["mcp"]
    }
  }
}
```

**Cursor**（`.cursor/mcp.json`）：

```json
{
  "mcpServers": {
    "modern-guidelines": {
      "command": "mg",
      "args": ["mcp"]
    }
  }
}
```

**Codex**（`~/.codex/config.toml`）：

```toml
[mcp_servers.modern-guidelines]
command = "mg"
args = ["mcp"]
```

**ZCode**（`~/.zcode/cli/config.json` 的 `mcp.servers`）：

```json
{
  "mcp": {
    "servers": {
      "modern-guidelines": {
        "command": "mg",
        "args": ["mcp"]
      }
    }
  }
}
```

> 以上配置即拷即用：`mg mcp` 是真实服务器（见[当前状态](#当前状态诚实声明)）。

## 目录导航

| 路径 | 内容 |
|---|---|
| [docs/prd/PRD-0000-modern-guidelines.md](docs/prd/PRD-0000-modern-guidelines.md) | 需求与决策总纲（D1–D12） |
| [docs/adr/](docs/adr/) | 决策记录：0001 通用多语言枢纽、0002 Go 编译单静态二进制、0003 分发与触发解耦、0004 Apache-2.0 吸收 JetBrains Go 数据集 |
| [docs/research/](docs/research/) | 内容研究记录：`python-modernization.md`（88 条候选，逐条官方引用）、`typescript-modernization.md`（双轴 57 条）——新语言接入的研究范例 |
| [docs/features/](docs/features/) | 各语言 FEATURES 文档（生成物，勿手改：`make generate-features` 由 `internal/featuresgen` 从数据集生成 python/go/typescript.md，88/54/57 条，漂移测试字节级锁定） |
| [CONTEXT.md](CONTEXT.md) | 领域术语表（dataset / axis / detector / provenance line / T1–T3 / 漂移测试等） |
| [CONTRIBUTING.md](CONTRIBUTING.md) | 贡献指南：schema 速查、规则工作流、新语言接入、ingestion 同步 |
| [NOTICE](NOTICE) | JetBrains 代码与数据署名 |
| `internal/testdata/upstream/` | JetBrains 上游快照的 ingestion 输入（`go-guidelines.json`，commit 锁定；完整参考项目保留为本地 `examples/`，不入库） |
| `internal/` | 源码：`cli`（子命令与 provenance 渲染）、`mcpserver`（内置 MCP stdio 服务器）、`schema`（数据集校验、版本比较权威）、`registry`（扩展名/数据集/检测器注册表）、`guidelines`（embed 加载过滤渲染，`data/*.json` 为数据集）、`detectors`（三语言版本检测器）、`featuresgen`（FEATURES 文档生成器）、`ingest`（JetBrains 吸收转换） |

## 许可证

Apache-2.0（见 [LICENSE](LICENSE)）。Go 数据集衍生自 JetBrains/go-modern-guidelines，署名与快照语义见 [NOTICE](NOTICE)。
