# 贡献指南

感谢贡献！本仓库的约定：**文档用中文，各语言规则数据集用英文**（`guideline`/`details` 等字段是给 agent 消费的准则文本）。动手前请先读 [CONTEXT.md](CONTEXT.md)（领域术语表）与 [README.md](README.md) 的"当前状态"一节——那里如实记录了已交付与尚未发布（Release/tag）的部分。

## 开发环境

- Go 1.26+（本仓 `go.mod` 声明 `go 1.26.5`）。
- 常用命令（[Makefile](Makefile)）：

```console
make test    # 全量测试（go test ./...）
make check   # test + gofmt + go vet，提交前门槛
make ingest  # 重跑 JetBrains Go 数据集吸收（见下文 ingestion 一节）
make generate-features  # 重建 docs/features/<lang>.md（internal/featuresgen，漂移测试锁定）
```

提交前跑 `make check`。

## 数据集 schema 速查

数据集是 `internal/guidelines/data/<lang>.json` 一份 JSON（`go:embed` 进二进制，`init` 时全量校验，**非法数据启动即 panic**）。权威定义在 `internal/schema/schema.go`；未知字段会被拒绝。

**数据集头**：

| 字段 | 必填 | 说明 |
|---|---|---|
| `language` | 是 | 小写字母/数字/连字符 |
| `axes` | 多轴时必填 | 版本轴声明，如 `["typescript","node"]`；单轴数据集省略（隐含语言自身为唯一轴） |
| `provenance` | 外部吸收时必填 | `{source, license, commit, date}`；自研数据集不写 |
| `rules` | 是 | 非空规则数组 |

**规则字段**：

| 字段 | 必填 | 说明 |
|---|---|---|
| `id` | 是 | 小写字母/数字/下划线；语言内唯一 |
| `since_version` | 是 | 点分数 `major[.minor]`（`"3.12"`、`"18"`）；规则的最低适用版本 |
| `axis` | 多轴数据集必填；单轴数据集禁止 | 必须在数据集头 `axes` 中声明 |
| `autofix` | **键必填** | `"tool:rule"`（如 `"ruff:UP007"`、`"gopls:modernize"`）或显式 `null`；缺键非法——作者必须显式声明"可自动改写"或"不可" |
| `category` | 是 | 非空，如 `"Tooling"`、`"Generics"` |
| `impact` | 是 | 枚举 `Critical` / `High` / `Medium` / `Low`（频率定义见下文） |
| `guideline` | 是 | 一行祈使句准则 |
| `details` | 是 | 机制 + 反例警示（30–90 词，见写作规范） |
| `references` | 可选 | ≤2 条短引用；Go 吸收数据集无引用 |
| `examples` | 是 | ≥1 组 `{before: [...], after: [...]}` |

**校验不变量**（违反即校验失败）：规则在**每个轴分组内**按 `since_version` newest-first 排列；id 去重；impact 在枚举内；字段非空；至少一组示例。

## 改规则的工作流

1. **改数据**：编辑 `internal/guidelines/data/<lang>.json`，保持 newest-first 排序。
2. **重新生成 FEATURES**：跑 `make generate-features` 重建 `docs/features/<lang>.md`。
   > `docs/features/<lang>.md` 是生成物（生成器为 `internal/featuresgen`），禁止手改；改了数据不重新生成就提交，会被字节级漂移测试拦下。
3. **漂移测试把关**：`go test ./...`——FEATURES 漂移测试做字节级比对，改了数据不重新生成即失败（ingestion 产物另有独立的字节级漂移测试）。
4. **冒烟**：`go run . languages`（数据非法会在启动时 panic）与 `go run . list --version <lang>=<version>` / `go run . explain <lang>:<id>` 各跑一遍。

## 规则写作规范

- **英文撰写**规则的全部文本字段。
- `guideline`：祈使句、一行、直说"用 X 代替 Y"。
- `details`：30–90 词，先讲**机制**（新写法做了什么、边界在哪），再给**反例警示**（什么时候不能无脑换）。例：pydantic 等运行时消费注解的场景，`UP006/UP007` 式改写不安全——这类 caveat 必须写进 details。
- `references`：只引官方信源（PEP 编号、whatsnew 锚点、官方 release note），每条保持简短，最多 2 条。
- `impact` 按出现频率定级（跨数据集共用同一标尺）：
  - `Critical`：几乎每个项目都有
  - `High`：常见，5–20 处
  - `Medium`：规律出现，1–5 处
  - `Low`：罕见或特定场景
- `autofix`：只在工具**确有**对应自动修复规则时填写并核实规则号（如 ruff 的 UP 系列须对照 docs.astral.sh 规则页），否则显式 `null`。

## 新语言接入指南

接入一个新语言 = **三件套**：数据集 + 检测器 + manifest 研究记录。

**第 0 步（中保真语言强制）：manifest 研究。** 若该语言的版本声明是 advisory 或推断型的中保真来源（类似 TS 的 `engines.node`、从 `devDependencies` 推断编译器版本），必须**先研究后定稿**（PRD 决策 D8）：在 `docs/research/<lang>-modernization.md` 落一份研究记录，写清可静态读取的版本来源、各来源的保真度与语义陷阱，内容规则逐条给出官方引用，未能官方核实的一律标 `UNVERIFIED`。范例就是现有两份：

- [docs/research/python-modernization.md](docs/research/python-modernization.md) —— 88 条候选、逐条 PEP/whatsnew 引用与核实方法。
- [docs/research/typescript-modernization.md](docs/research/typescript-modernization.md) —— 双轴 57 条、TS 官方 release notes 与 nodejs.org API 文档逐字提取的信源纪律。

高保真语言（manifest 即强制地板，如 Rust 的 `rust-version`、Ruby 的 Gemfile `ruby`）研究记录同样建议写，但不强制。

**第 1 步：数据集。** 新建 `internal/guidelines/data/<lang>.json`（惯例文件名 = 语言名）。放入 `data/` 即被 `go:embed` 自动加载、校验并注册，无需改加载代码。多轴语言记得声明 `axes` 并给每条规则标 `axis`。

**第 2 步：检测器。** 新建 `internal/detectors/<lang>.go`，实现 `registry.Detector` 接口（`Language()` + `Detect()`，返回"轴→{version, source, fidelity}"映射），并在包 `init()` 里 `registry.RegisterDetector(...)`。检测语义务必表驱动测试（`t.TempDir`/`t.Chdir` 构造嵌套 manifest 场景；T2 各来源、T3 兜底、advisory 标注都要有边界用例），参考现有三个检测器的测试。

**第 3 步：注册接线。** 在 `internal/registry` 的 `extensionToLanguage` 加扩展名映射；在 `main` 包为检测器包添加空白导入以触发注册副作用（registry 从不导入具体实现）。
> 现有三语言已接线（`main.go` 空白导入 `internal/detectors`）；接入新语言时把新检测器包加进 `main.go` 的同一处空白导入即可。

**第 4 步：验证。** `go run . languages` 出现新语言；`list`/`explain` 冒烟；`make check` 全绿；更新本 README 的支持语言表与 `CHANGELOG.md`。

## ingestion 同步流程（JetBrains Go 数据集）

Go 数据集是 JetBrains/go-modern-guidelines 的**一次性 commit 锁定快照**（ADR 0004），不是持续 vendor：上游发布新 minor 时才重跑同步。

1. 更新 `internal/testdata/upstream/go-guidelines.json`（从上游仓库对应 commit 逐字节重新导出；该文件是 ingestion 唯一消费的上游文件，出处记录在同目录 README.md），并同步修改 `internal/ingest/ingest.go` 中锁定的 `UpstreamCommit`（与数据集头 `provenance.commit`/`date` 由重跑自动写入）。
2. 跑 `make ingest`：转换器（`go run ./scripts/ingest-jetbrains-go`）从 `internal/testdata/upstream/go-guidelines.json` 重生成 `internal/guidelines/data/go.json`，原位覆盖。
3. `go test ./...`：`internal/ingest` 的漂移测试对转换结果做字节级锁定——转换逻辑与已提交数据集的任何分歧都会失败。
4. **diff 评审** `data/go.json` 的变化后再提交；确认 54 条零丢失、`modernizer:true → autofix:"gopls:modernize"` 字段映射正确。
5. 署名不变时 [NOTICE](NOTICE) 无需改动。

## 其他

- 提交信息与 CHANGELOG：用户可感知的行为变化补进 [CHANGELOG.md](CHANGELOG.md)（Keep-a-Changelog 格式）。
- 参考项目快照不入库：`internal/testdata/upstream/` 只保存 ingestion 消费的 `go-guidelines.json`；完整的 `examples/go-modern-guidelines/` 保留为本地参考目录（已在 `.gitignore`），除上述 ingestion 同步外不要提交或修改。
- 决策类讨论请落 ADR（`docs/adr/`），新术语进 [CONTEXT.md](CONTEXT.md)。
