# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-09-29

初始版本：版本门控的多语言现代准则枢纽（Python / Go / TypeScript），Go 编译单静态二进制 `mg`，为 AI agent 修复训练截止后的新特性缺失与旧习语频率偏差。

### Added

- `mg` CLI：`list`（`--version axis=ver[,...]` / `--file-path <path>` / 位置路径三选一；输出一行一条 newest-first，多轴数据集按轴分组，首行每轴一条 provenance 行）、`explain`（`lang:id` 或 `--lang` 加裸 id；Since/Summary/Details/Examples/References/Autofix 缩进块；未知 id 报错附全部可用 id）、`languages`、`version`/`--version`、`help`。
- `mcp` 子命令：内置 MCP stdio 服务器（`internal/mcpserver`，基于官方 Go SDK `github.com/modelcontextprotocol/go-sdk v1.8.0`），提供 `list_guidelines` / `explain_guideline` 两个工具，工具输出与对应 CLI 子命令一致（仅末尾换行被裁去）。
- 三份内嵌数据集（`go:embed`，`init` 时全量校验，非法数据启动即 panic）：
  - Python 88 条（3.9 → 3.15；3.15 为占位条目，final 发布后核实转正）。
  - Go 54 条（1.0 → 1.27；自 JetBrains/go-modern-guidelines 的 commit 锁定快照吸收，`modernizer:true → autofix:"gopls:modernize"` 映射，Apache-2.0 署名见 NOTICE）。
  - TypeScript 57 条双轴（typescript 轴 33 条 4.9 → 7.0 + node 轴 24 条 16 → 24）。
- 数据集 schema v2（`internal/schema`）：多轴 `axes`/`axis` 声明与校验、轴内 newest-first 不变量、`autofix` 键必填（`"tool:rule"` 或显式 `null`）、impact 枚举、references ≤2 条、未知字段拒绝。
- 语言注册表（`internal/registry`）：扩展名→语言映射（.py/.go/.ts/.mts/.cts/.tsx；.js/.mjs 预留 javascript）、数据集与检测器注册接口、保真度分级（high/medium）。
- 版本解析三层中的 T1：`--version axis=ver[,...]`（多轴逗号形式、全轴覆盖校验、轴必须归属唯一语言）；每次 `list` 输出每轴一行 provenance（版本、来源、保真度），非高保真（advisory/模糊）来源附警示并建议显式 `--version`。
- 三语言版本检测器（`internal/detectors`，已通过 `main.go` 空白导入接线进主二进制，`--file-path` 直接可用）：
  - Python：pyproject `requires-python` 下界（`>=3.10,<3.13`→3.10、`==3.11.*`→3.11、`~=3.10`→3.10）/ PEP 723 内联 / `.python-version` / mise，find-up 嵌套上溯。
  - Go：go.mod `go` 指令（无指令默认 1.16）/ go.work（默认 1.18 且不覆盖成员）/ `toolchain` 指令仅提示不作地板 / `go version` 兜底。
  - TypeScript 双轴：`devDependencies.typescript` 范围下界（`^5.4`→5.4）+ node 轴 `engines.node`（advisory 标注 medium 保真度）/ `.nvmrc` / mise / `node --version` 兜底。
- ingestion 管道（`scripts/ingest-jetbrains-go` + `internal/ingest`）：上游 commit 锁定（`make ingest` 重跑），转换结果由字节级漂移测试锁定。
- FEATURES 文档生成器（`internal/featuresgen`）：`make generate-features` 从数据集重建 `docs/features/python.md` / `go.md` / `typescript.md`（88/54/57 条，与数据集字节级漂移测试锁定）。
- Makefile 维护目标：`test`（go test ./...）、`generate-features`、`ingest`、`check`（test + gofmt + go vet）。
- 项目脚手架：Go module、Apache-2.0 LICENSE 与 NOTICE（JetBrains 署名）、README（快速开始三通道、四家 MCP 配置、AGENTS.md 触发片段）、CONTRIBUTING（schema 速查与新语言接入指南）、AGENTS.md 自举片段、CONTEXT.md 领域术语表、docs/prd 与 docs/adr（0001–0004）与 docs/research 研究记录。
- 发布流水线（PRD Req 9/13）：`.goreleaser.yml`（`mg` 单静态二进制，linux/darwin/windows × amd64/arm64 六目标 `CGO_ENABLED=0` 构建，tar.gz 归档（windows 为 zip）+ checksums.txt，配置经 `goreleaser check` v2.18.2 校验通过）、`.github/workflows/ci.yml`（`make check`、六目标交叉构建、`mg list --file-path main.go` dogfood 自检）与 `.github/workflows/release.yml`（`v*` tag 触发 goreleaser 发布 GitHub Releases）。
- 验收核实记录（AC13 端到端实测，2026-09-29）：三扩展名 `--file-path` 调用全部通过且 exit 0——`.py`（pyproject `requires-python: >=3.12` → `# python 3.12 (source: pyproject.toml requires-python, fidelity: high)`）、`.go`（go.mod `go 1.22` → `# go 1.22 (source: go.mod go directive, fidelity: high)`）、`.ts`（devDependencies `^5.4` → `# typescript 5.4 …(medium)`，engines.node `>=18` → `# node 18 …(medium)`，advisory 警示按预期出现）；dogfood `go run . list --file-path main.go` → `# go 1.26 (source: go.mod go directive, fidelity: high)`。六目标交叉构建（`CGO_ENABLED=0 go build`）全部成功，linux 产物静态链接，`mg list` 冒烟通过。
