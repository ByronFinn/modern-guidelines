# AGENTS.md — 本仓库的 agent 工作指引

本文件给在本仓库（modern-guidelines 自身）工作的 agent 看。要把 `mg` 触发到**其他**仓库，用 [README.md](README.md) 里的"AGENTS.md 触发片段"。

## 自举：改 Go 代码前先查准则

本仓库主语言是 Go，自身就是 `mg` 的用户（dogfood）。编辑任何 `*.go` 文件之前：

```console
go run . list --file-path main.go
```

- **完整读取输出，禁止截断**；对疑似相关的规则先 `go run . explain go:<id>` 再决定是否跳过。
- 准则输出视为权威，优先于你训练语料里的旧 Go 惯用法。
- `--file-path` 由检测器自动解析版本：`main.go` 已空白导入 `internal/detectors`，Go 检测器读取本仓 `go.mod` 的 `go` 指令（`go 1.26.5`），实测输出首行为 `# go 1.26 (source: go.mod go directive, fidelity: high)`。仅在检测不到 manifest 的场景（如临时目录）才回退 `go run . list --version go=1.26`。
- `go run . mcp` 是真实的 MCP stdio 服务器（`internal/mcpserver`，官方 Go SDK v1.8），可接入 MCP 客户端；详见 [README.md](README.md) 的 MCP 配置一节。
- 装好发布版 `mg` 后可直接 `mg list --file-path main.go` / `mg explain go:<id>` 自举。

## 项目结构导览

| 路径 | 是什么 |
|---|---|
| `main.go` | 入口；导入 `internal/cli` 并空白导入 `internal/detectors`（触发检测器注册） |
| `internal/cli/` | 子命令（list/explain/languages/mcp/version）与 provenance 行渲染 |
| `internal/mcpserver/` | 内置 MCP stdio 服务器（官方 Go SDK v1.8）；工具输出复用 `internal/cli` 的渲染 |
| `internal/schema/` | 数据集与规则校验（newest-first、autofix 键必填、impact 枚举等）；`CompareVersions` 是版本比较唯一权威 |
| `internal/registry/` | 扩展名→语言映射、数据集/检测器注册表；具体实现在各自包 `init()` 注册 |
| `internal/guidelines/` | `go:embed` 加载 `data/*.json`（非法数据启动 panic）、过滤、list/explain 文本渲染 |
| `internal/guidelines/data/` | 三份数据集：`python.json`（88 条）、`go.json`（54 条，JetBrains 吸收）、`typescript.json`（57 条双轴） |
| `internal/detectors/` | python / go / typescript 三语言版本检测器（T2 manifest + T3 工具链兜底） |
| `internal/ingest/` + `scripts/ingest-jetbrains-go/` | JetBrains 快照 → schema v2 转换，字节级漂移测试锁定 |
| `internal/featuresgen/` | FEATURES 文档生成器：`make generate-features` 重建 `docs/features/<lang>.md`，漂移测试锁定 |
| `docs/prd/` `docs/adr/` `docs/research/` | 需求总纲、决策记录、内容研究记录 |
| `CONTEXT.md` | 领域术语表（dataset/axis/detector/provenance line/T1–T3 等） |
| `examples/go-modern-guidelines/` | JetBrains 参考项目只读快照（原 ingestion 上游；现已不入库，见 `.gitignore`），除同步外勿改 |
| `internal/testdata/upstream/` | ingestion 唯一消费的上游文件（`go-guidelines.json`，commit 锁定，随仓库提交） |

## 仓库纪律

- **文档中文，规则数据集英文**；改数据集前读 [CONTRIBUTING.md](CONTRIBUTING.md) 的 schema 速查与写作规范。
- 改 `data/*.json` 后：数据非法会在 `init` panic，先 `go run . languages` 冒烟，再 `go test ./...`（漂移测试把关）；提交前 `make check`。
- 决策变化落 ADR，新术语进 CONTEXT.md，用户可感知变化进 CHANGELOG.md。
