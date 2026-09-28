# JetBrains go-modern-guidelines 上游快照（ingestion 输入）

本目录保存 `internal/ingest`（`make ingest` / `scripts/ingest-jetbrains-go`）唯一消费的上游文件，
是 `internal/guidelines/data/go.json` 的转换源，由字节级漂移测试锁定。

- **来源**：github.com/JetBrains/go-modern-guidelines
- **上游 commit（锁定）**：`019b45e2b2f80d7f5c1e28bd4d35f3f0fbcf9c9c`（与 `internal/ingest/ingest.go` 的 `UpstreamCommit` 一致）
- **许可**：Apache-2.0；署名与快照语义见仓库根 [NOTICE](../../../NOTICE)（本目录不重复存放 LICENSE 副本）
- **文件**：`go-guidelines.json` — 逐字节复制自上游 `internal/guidelines/guidelines.json`（commit 019b45e2）；JSON 严格解码（`DisallowUnknownFields`），不加注释头，出处以本 README 记录

完整的参考项目（源码、LICENSE 原文等）保留为本地参考目录 `examples/go-modern-guidelines/`，
该目录不入库（见 `.gitignore`），除 ingestion 同步外勿改。
