# TypeScript 数据集内容研究：typescript 轴 + node 轴（D12 双轴）

> **Status**: 研究定稿 | **对应**: PRD-0000 需求 3/5、D8（先研究后定稿）、D12（TS+Node 双轴） | **研究快照**: 2026-09-29 | **候选规则总数**: 57（TS 轴 33 + Node 轴 24；PR7 从中挑选 ~40 定稿）

## 0. 研究方法与信源纪律

- **只采信官方信源**，逐条给出引用；凡未能在官方页面原文中核实的一律标 **UNVERIFIED**（见 §5）。
- TS 语言特性：**typescriptlang.org 官方 release notes**（4.9–5.4 逐版页面，特性标题逐字提取）；5.5 起官方发布说明正式载体为 **TypeScript 团队官方博客 devblogs.microsoft.com/typescript**（typescriptlang.org 文档链出的 announcement，本报告视为官方公告源，逐字提取标题）。
- Node：**nodejs.org API 文档**（`/docs/latest/api/*` 的 "Added in" / History 表逐字引用）与 **nodejs.org release notes**（`/en/blog/release/vX.Y.Z`）。
- V8 语言特性语义：**developer.mozilla.org**（MDN 允许信源；注意 MDN 兼容表在本研究抓取中未渲染出 Node 版本号，相关 Node 落地版本标 UNVERIFIED）。
- **autofix 字段**仅当 typescript-eslint 存在真实 fixer 时填写；核实源为 typescript-eslint.io 规则页（eslint 官方插件自身文档，作为 autofix 字段专用权威源，特此声明）。
- 版本口径：TS 轴 `since_version` 用 TS 版本（4.9–7.0）；Node 轴用 major（个别条目精确到 minor 以反映 stable/unflag 节点）。

### 引用键

| 键 | 来源 |
|---|---|
| [TS-RN-4.9]~[TS-RN-5.4] | `typescriptlang.org/docs/handbook/release-notes/typescript-{4-9,5-0,5-1,5-2,5-3,5-4}.html` |
| [TS-B-5.5]~[TS-B-7.0] | `devblogs.microsoft.com/typescript/announcing-typescript-{5-5,5-6,5-7,5-8,5-9,6-0,7-0}/` |
| [TS-BLOG] | `devblogs.microsoft.com/typescript/`（官方博客索引，含发布日期） |
| [N-GLOBALS] | `nodejs.org/docs/latest/api/globals.html` |
| [N-TEST] | `nodejs.org/docs/latest/api/test.html` |
| [N-UTIL] | `nodejs.org/docs/latest/api/util.html` |
| [N-FS] | `nodejs.org/docs/latest/api/fs.html` |
| [N-CLI] | `nodejs.org/docs/latest/api/cli.html` |
| [N-ESM] | `nodejs.org/docs/latest/api/esm.html` |
| [N-MODULES] | `nodejs.org/docs/latest/api/modules.html` |
| [N-TS] | `nodejs.org/docs/latest/api/typescript.html`（Node 原生运行 TS） |
| [N-ASSERT] | `nodejs.org/docs/latest/api/assert.html` |
| [N-PROCESS] | `nodejs.org/docs/latest/api/process.html` |
| [N-SQLITE] | `nodejs.org/docs/latest/api/sqlite.html` |
| [N-RN-18.11] [N-RN-21.0] [N-RN-22.0] [N-RN-22.12] [N-RN-24.0] | `nodejs.org/en/blog/release/v{18.11.0, 21.0.0, 22.0.0, 22.12.0, 24.0.0}/` |
| [N-PREV] | `nodejs.org/en/about/previous-releases`（各版本线状态表） |
| [N-SCHED] | `nodejs.org/en/blog/announcements/evolving-the-nodejs-release-schedule/` |
| [MDN-fromAsync] 等 | `developer.mozilla.org/.../Global_Objects/{Array/fromAsync, Promise/withResolvers, Set/union, Object/groupBy}` |
| [TSE-CTI] [TSE-NITSE] | `typescript-eslint.io/rules/{consistent-type-imports, no-import-type-side-effects}/` |

## 1. 版本门版图（已核实）

### 1.1 TypeScript 轴（4.9 → 7.0）

- **当前最新稳定版：TypeScript 7.0**（原生 Go 移植，官方公告 2026-07-08："Today we are proud to announce the availability of TypeScript 7, a 10x faster native port of TypeScript!"；"often about 10 times faster than TypeScript 6.0"；与 6.0 side-by-side 运行）[TS-B-7.0][TS-BLOG]。
- TypeScript 6.0（2026-03-23 发布）：基于现有 JS 代码库的**最后一个大版本**，承担面向 7.0 的弃用清理（node10 解析、baseUrl、es5 target 弃用，7.0 彻底移除；outFile 在 6.0 已移除、amd/umd/systemjs 模块值不再支持）[TS-B-6.0]。
- TypeScript 5.9（2025-08-01 发布）[TS-B-5.9]；5.8/5.7/5.6/5.5/5.4/5.3/5.2/5.1/5.0/4.9 依次向下（4.9–5.4 官方 release notes 页在站内，页面未展示发布日期——见 §5 UNVERIFIED）。
- **门控语义提醒（检测器已定）**：TS 版本从 `devDependencies.typescript` 范围下界推断（如 `^5.4`→5.4）；`compilerOptions.target` 是 JS 产出级别，**不是**编译器版本（PRD 需求 5，子代理 E 已核）。

### 1.2 Node 轴（16 → 24 LTS）

| 版本线 | 状态（[N-PREV]） | 备注 |
|---|---|---|
| 26 | Current（2026-05-05 首发） | 非 LTS |
| **24 (Krypton)** | **LTS**（2025-05-06 首发） | 目标上限；V8 13.6 |
| **22 (Jod)** | **LTS**（2024-04-24 首发） | 仍在维护更新 |
| 20 (Iron) | **EOL**（最后更新 2026-03-24） | — |
| 18 (Hydrogen) | EOL | — |
| 16 (Gallium) | EOL（2023-08-08） | 数据集下界 |

- Node 24 (V8 13.6) 官方点名新 JS 特性：**Float16Array、Explicit resource management、RegExp.escape、WebAssembly Memory64、Error.isError**；npm 11；URLPattern 全局暴露；Permission Model 改名 `--permission` [N-RN-24.0]。
- 发布节奏变更（2026-03 公告）："Starting with 27.x, Node.js will move from two major releases per year to one."；"One major release per year (April), with LTS promotion in October."；"Every release becomes LTS." [N-SCHED] —— 未来奇偶区分消失。
- **门控语义提醒**：`engines.node` 为 advisory，provenance 必须带警示（PRD 需求 5、D12）。

## 2. typescript 轴规则表（33 条候选）

表格列：id | 轴 | 版本 | 旧→新 | 引用 | impact | autofix | caveat。rationale 与建议 category 见 §2.1 逐条说明。

| id | 轴 | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
|---|---|---|---|---|---|---|---|
| ts7_native_compiler | typescript | 7.0 | JS 版 tsc → TS 7 原生编译器（Go，~10x，与 6.0 并行装） | [TS-B-7.0] | Medium | — | 语言能力与 6.0 对齐但有行为差异（如 template literal types 保留 Unicode code points）；npm 包与 6.0 并行，勿盲目全量切换 |
| ts60_deprecated_flags | typescript | 6.0 | `--moduleResolution node10`、`--baseUrl`、`target: es5`、`outFile`、amd/umd/system → bundler/nodenext、相对路径、es2015+ | [TS-B-6.0] | High | — | 6.0 分两类：node10/baseUrl/es5 仅"弃用"（`ignoreDeprecations: "6.0"` 可暂缓，7.0 彻底移除）；outFile 在 6.0 已移除、amd/umd/systemjs/none 已不再支持（须在升级到 6.0 前迁离）；切 bundler 需同步改相对导入写法 |
| ts60_subpath_imports | typescript | 6.0 | 深相对路径 → package.json `#` 子路径导入（6.0 支持） | [TS-B-6.0] | Low | — | 需 tsconfig 与 package.json `imports` 字段配合 |
| ts60_bundler_commonjs | typescript | 6.0 | bundler 解析仅配 ESM → 允许 `--moduleResolution bundler` + `--module commonjs` | [TS-B-6.0] | Low | — | 仅"bundler 产物 + CJS 输出"的工具链场景 |
| ts60_target_es2025 | typescript | 6.0 | target/lib es2022 → es2025（含 RegExp.escape、Temporal 新类型） | [TS-B-6.0] | Medium | — | 类型≠运行时：RegExp.escape 运行时需 Node 24（V8 13.6）[N-RN-24.0]；Temporal 运行时仍未定 |
| ts59_import_defer | typescript | 5.9 | 急切 import namespace → `import defer * as ns` | [TS-B-5.9] | Low | — | 依赖运行时对 deferred import 的支持；仅模块初始化有副作用代价场景收益 |
| ts59_module_node20 | typescript | 5.9 | `--module nodenext`（语义漂移）→ 显式 `--module node20` | [TS-B-5.9] | Medium | — | 5.8 先加了 `node18` 档 [TS-B-5.8]；固定 node 版本可让解析语义不随编译器升级漂移 |
| ts58_erasable_syntax_only | typescript | 5.8 | enum/参数属性/namespace 运行时代码 → `--erasableSyntaxOnly` 强制类型可擦除 | [TS-B-5.8] | High | — | 启用后禁用 enum、参数属性、运行时 namespace 等；与 Node 原生 type stripping 推荐 tsconfig 一致 [N-TS] |
| ts58_require_esm_nodenext | typescript | 5.8 | CJS 侧无类型化 require(esm) → `--module nodenext` 下支持 | [TS-B-5.8] | Medium | — | 对应 Node 22.12 unflag [N-RN-22.12]；仅同步 ESM 图，TLA 抛 `ERR_REQUIRE_ASYNC_MODULE` |
| ts57_rewrite_relative_import_paths | typescript | 5.7 | 声明文件相对路径手工映射 → `--rewriteRelativeImportPaths` | [TS-B-5.7] | Low | — | 面向复杂 outDir/monorepo 布局；单一输出目录项目无需 |
| ts56_no_unchecked_side_effect_imports | typescript | 5.6 | `import "./foo.css"` 不校验 → `--noUncheckedSideEffectImports` 开启存在性检查 | [TS-B-5.6] | Medium | — | 仅解析失败时报错；CSS-in-JS 等带 loader 管线需确认路径可静态解析 |
| ts56_iterator_helpers_types | typescript | 5.6 | 手写迭代器 map/filter/take → 内置 iterator helpers + lib 类型 | [TS-B-5.6] | Medium | — | 运行时门控：Node 22（V8 12.4 ≥ Chrome 122 的 V8 版本；精确 Node 首发版本 UNVERIFIED，保守取 22） |
| ts56_no_check | typescript | 5.6 | tsc 全量又编译又查 → `--noCheck` 只编译（CI 拆 typecheck/build） | [TS-B-5.6] | Low | — | 与 `--build` 组合语义见原文；不替代独立 `tsc --noEmit` 检查 |
| ts55_inferred_type_predicates | typescript | 5.5 | 手写 `(x): x is T` 守卫 → 箭头函数自动推断类型谓词 | [TS-B-5.5] | High | — | 推断条件：无显式返回标注、单一 return、不突变参数、布尔返回与参数收窄相关；函数声明与箭头函数均适用（官方示例即函数声明）；带显式返回标注或多 return 的复杂守卫仍需手写 `x is T` |
| ts55_set_methods_types | typescript | 5.5 | 手写 Set 并/交/差 → `Set.prototype.union/intersection/difference/…` + 类型 | [TS-B-5.5] | Medium | — | 运行时门控 Node 22（MDN Baseline 2024-06；Node 精确首发 UNVERIFIED）[MDN-set-union] |
| ts55_config_dir | typescript | 5.5 | `../shared` 相对路径 hack → `${configDir}` 模板变量 | [TS-B-5.5] | Medium | — | 面向 monorepo 共享 tsconfig；仅限被 extends 的配置里使用 |
| ts55_regex_syntax_checking | typescript | 5.5 | 坏正则运行时才炸 → 编译期正则语法检查 | [TS-B-5.5] | Medium | — | 检查语法而非语义；动态拼接的正则串不可静态判 |
| ts55_jsdoc_import | typescript | 5.5 | .js 里重复类型结构 → JSDoc `@import` 复用类型 | [TS-B-5.5] | Low | — | 仅 JS（checkJs）场景；TS 数据集中标记适用边界 |
| ts54_noinfer | typescript | 5.4 | 哨兵默认值/never hack → `NoInfer<T>` | [TS-B-5.4] | High | — | 需要调用方双向推断保持原样时不要用；仅在"该形参不参与推断"时受益 |
| ts54_preserved_narrowing_closures | typescript | 5.4 | 闭包内窄化丢失（重复判空）→ 最后赋值后创建的闭包保留窄化 | [TS-B-5.4] | Medium | — | 仅"最后一次赋值之后创建"的闭包生效；可变 let 被再赋值仍丢失 |
| ts54_group_by_types | typescript | 5.4 | `reduce` 手写分组 → `Object.groupBy` / `Map.groupBy`（类型支持） | [TS-B-5.4] | Medium | — | 运行时门控：Node 21（V8 11.8）应已含，nodejs.org 未明文（UNVERIFIED），保守取 22；MDN Baseline 2024-03 [MDN-groupBy] |
| ts53_import_attributes | typescript | 5.3 | `import … assert { type: "json" }` → `import … with { type: "json" }`（Import Attributes） | [TS-RN-5.3] | Medium | — | Node 22 起 **移除** import assertions（"esm: drop support for import assertions"）[N-RN-22.0]——`with` 是唯一可用形式 |
| ts53_switch_true_narrowing | typescript | 5.3 | if/else 判别链 → `switch (true)` 获得窄化 | [TS-RN-5.3] | Low | — | 属风格可选；窄化收益依赖 case 写法 |
| ts52_using_declarations | typescript | 5.2 | try/finally 手动清理 → `using` / `await using`（Explicit Resource Management） | [TS-RN-5.2] | High | — | **运行时门控：Node 24（V8 13.6 点名 "Explicit resource management"）**[N-RN-24.0]；Node 22 及以下需 polyfill 或降级发射 |
| ts52_decorator_metadata | typescript | 5.2 | 手写元数据注入 → `Symbol.metadata`（Decorator Metadata） | [TS-RN-5.2] | Low | — | 运行时支持有限；legacy `--experimentalDecorators` 不适用 |
| ts50_decorators_stage3 | typescript | 5.0 | legacy 实验装饰器 → TC39 Stage 3 标准装饰器（稳定） | [TS-RN-5.0] | High | — | 语义与 legacy 完全不同（编译器不再转译）；生态（NestJS 等）多仍在 legacy；Node type stripping 对装饰器直接 parser error [N-TS] |
| ts50_const_type_parameters | typescript | 5.0 | `as const` 传参/字面量丢失 → `const` 类型参数 `<const T>` | [TS-RN-5.0] | High | — | 仅推断位置生效；显式传类型实参时 const 无意义 |
| ts50_module_resolution_bundler | typescript | 5.0 | `node`/`node10` 解析（面向 bundler 项目失真）→ `--moduleResolution bundler` | [TS-RN-5.0] | High | — | 不能与 `--module commonjs` 组合（6.0 才放开 [TS-B-6.0]）；Node 运行项目用 nodenext |
| ts50_verbatim_module_syntax | typescript | 5.0 | isolatedModules 零散约定 → `--verbatimModuleSyntax` 统一"原样发射"纪律 + `import type` | [TS-RN-5.0] | High | **typescript-eslint:consistent-type-imports**（可加 no-import-type-side-effects） | 开启后未标 type 的纯类型导入直接报错，需先完成 import type 迁移；与 Node type stripping 推荐配置一致 [N-TS] |
| ts50_enum_union_types | typescript | 5.0 | enum 成员≈数字 → 所有 enum 皆为 union enum（计算成员也有唯一字面量类型，可窄化、可作类型引用） | [TS-RN-5.0] | Medium | — | Enum Overhaul 同时引入两个新错误（越界字面量赋值、间接混串/数字 enum）；**注**：任务清单所称"任意枚举成员名"未在官方 5.0 notes 出现，已修正为真实特性名（§5） |
| ts49_satisfies | typescript | 4.9 | 类型注解校验但拓宽字面量 → `satisfies` 校验且保留精确推断 | [TS-RN-4.9] | High | — | 校验可赋值性且不改变表达式推断类型（上下文类型照常生效，字面量信息保留）；过度使用会牺牲可读性，公共 API 边界仍建议注解 |
| ts49_in_operator_narrowing | typescript | 4.9 | 手写 `'k' in x` 自定义守卫 → `in` 原生窄化未声明属性 | [TS-RN-4.9] | Medium | — | 只窄化"未在类型上声明"的属性；声明过的键不触发 |
| ts49_auto_accessors | typescript | 4.9 | 手写 getter/setter + 背板字段 → `accessor` 自动访问器 | [TS-RN-4.9] | Medium | — | 发射目标需支持（与标准装饰器配套语义）；旧 target 下降级发射，注意产物差异 |

### 2.1 rationale 逐条（一句）与建议 category

- ts7_native_compiler — 共享内存并行 + 原生代码使类型检查/构建约 10 倍提速，CI 与编辑器体验质变。Tooling。
- ts60_deprecated_flags — 6.0 是通往 7.0 的弃用清理版（node10/baseUrl/es5 弃用、outFile 移除、amd/umd/systemjs 不再支持），先迁离弃用面才能平滑吃到 7.0。Tooling。
- ts60_subpath_imports — 用包内 `#/` 别名替代深相对路径，重构时导入路径稳定。Syntax。
- ts60_bundler_commonjs — bundler 打包但输出 CJS 的工具链不再需要伪装成 nodenext。Tooling。
- ts60_target_es2025 — target 跟进新 ES 特性类型，减少 polyfill/手写类型。Tooling。
- ts59_import_defer — 推迟模块求值，砍掉不需副作用的依赖初始化开销。Syntax。
- ts59_module_node20 — 显式钉住 Node 解析语义，升级编译器不再隐式改变解析行为。Tooling。
- ts58_erasable_syntax_only — 强制"类型可整体擦除"，是 Node 原生直跑 .ts 的语法前提。Tooling。
- ts58_require_esm_nodenext — CJS 存量渐进吸收 ESM 依赖，无需先整体改造包格式。Modules。
- ts57_rewrite_relative_import_paths — 声明文件自动改写导入扩展名，去掉发布期路径 hack。Tooling。
- ts56_no_unchecked_side_effect_imports — 副作用导入的路径错误从"运行时静默失败"提前到编译期。Tooling。
- ts56_iterator_helpers_types — 惰性链式迭代替代中间数组，类型与运行时一次到位。Lib & Target。
- ts56_no_check — 大仓 CI 把"编译"与"类型检查"拆道并行，缩短关键路径。Tooling。
- ts55_inferred_type_predicates — 消灭最常见的一类手写守卫样板，过滤/判别回调自动获得 `x is T`。Type System。
- ts55_set_methods_types — 集合运算用内建方法替代手写循环+新增 Set，语义更准、代码更短。Type System。
- ts55_config_dir — 共享 tsconfig 不再被"相对于哪份文件"的 extends 解析坑绊倒。Tooling。
- ts55_regex_syntax_checking — 坏正则从生产崩溃提前到编辑器红线。Tooling。
- ts55_jsdoc_import — JS 存量文件按需引入 .d.ts 类型，不必整仓迁 TS。Type System。
- ts54_noinfer — 修复"默认值参与推断导致推断过宽"的泛型 API 设计痛点，删掉哨兵 hack。Type System。
- ts54_preserved_narrowing_closures — 回调里不必重复判空，窄化结果按最后赋值保留。Type System。
- ts54_group_by_types — 分组是一等内建，替代 reduce 样板且键类型精确。Lib & Target。
- ts53_import_attributes — `assert` 已被标准废弃为 `with`，Node 22 起旧写法直接不支持。Syntax。
- ts53_switch_true_narrowing — 判别逻辑集中一处且每支分支自动窄化。Syntax。
- ts52_using_declarations — 作用域退出即释放，替代 try/finally 清理样板并覆盖提前 return/throw。Syntax。
- ts52_decorator_metadata — 框架不再需要手写元数据注册表。Syntax。
- ts50_decorators_stage3 — 标准装饰器跨工具链一致，不再依赖编译器私有转译。Syntax。
- ts50_const_type_parameters — 推断元组/字面量时免写 `as const`，调用点更干净。Type System。
- ts50_module_resolution_bundler — bundler 项目的解析真相（无扩展名、package.json exports）终于有官方档位。Tooling。
- ts50_verbatim_module_syntax — 一条开关统一模块发射纪律，是 ESM/isolatedModules/type stripping 时代的地基。Tooling。
- ts50_enum_union_types — enum 获得完整字面量联合语义，窄化与成员类型引用全面可用。Type System。
- ts49_satisfies — 校验与推断解耦：既检查形状又保住字面量类型。Type System。
- ts49_in_operator_narrowing — 对"可能有额外键"的对象用 `in` 即窄化，免写守卫函数。Type System。
- ts49_auto_accessors — 一行替代 getter/setter + 私有字段三件套。Syntax。

## 3. node 轴规则表（24 条候选）

| id | 轴 | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
|---|---|---|---|---|---|---|---|
| node24_regexp_escape | node | 24 | 手写正则转义函数 → `RegExp.escape` | [N-RN-24.0]（V8 13.6 点名）+ [TS-B-6.0]（es2025 类型） | Medium | — | TS 侧类型需 lib es2025（TS 6.0） |
| node24_explicit_resource_mgmt | node | 24 | try/finally 清理 → `using`/`await using` 运行时支持（Symbol.dispose） | [N-RN-24.0]（"Explicit resource management"） | High | — | **TS 5.2 语法 ≠ 运行时可用**：Node 22 及以下无内建支持，需 polyfill；与 ts52_using_declarations 跨轴联动 |
| node23_type_stripping | node | 23.6（22.18 回移） | tsc 构建管线 → Node 原生直跑 .ts（type stripping 默认开） | [N-TS]（"Added in: v22.6.0"；"v23.6.0, v22.18.0 Type stripping is enabled by default"；"v25.2.0, v24.12.0 … now stable"） | High | — | 仅可擦除语法：enum、含运行时代码的 namespace、参数属性、import 别名抛 `ERR_UNSUPPORTED_TYPESCRIPT_SYNTAX`；装饰器 parser error；`.tsx` 不支持；tsconfig `paths` 被忽略；**不做类型检查** |
| node22_array_from_async | node | 22 | `for await` 手动收集 → `Array.fromAsync` | [MDN-fromAsync] + [N-RN-22.0]（V8 12.4.254.14） | Medium | — | **顺序 await**："Array.fromAsync() awaits each value yielded from the object sequentially. Promise.all() awaits all values concurrently."（MDN 原文）；Node 精确首发版本 UNVERIFIED（V8 12.4 推断，保守 22） |
| node22_promise_with_resolvers | node | 22 | 外部 `let resolve, reject` → `Promise.withResolvers()` | [MDN-withResolvers] | Low | — | Node 精确首发 UNVERIFIED（保守 22）；TS 类型需较新 lib |
| node22_set_methods | node | 22 | 手写并/交/差集 → `Set.prototype.union/intersection/difference/symmetricDifference/…` | [MDN-set-union]（Baseline 2024-06） | Medium | — | Node 精确首发 UNVERIFIED（保守 22）；TS 类型 5.5 起 [TS-B-5.5] |
| node22_fs_glob | node | 22 | glob/fast-glob 依赖 → `fs.glob` / `fs.globSync` / `fs.promises.glob` | [N-FS]（"Added in: v22.0.0"；"v24.0.0, v22.17.0 Marking the API stable"）+ [N-RN-22.0]（"fs: expose glob and globSync"） | Medium | — | 22.0–22.16 为 experimental；稳定门控取 22.17/24 |
| node22_watch_stable | node | 22 | nodemon/watchpack 依赖 → `node --watch`（22 起稳定） | [N-RN-18.11]（"cli: add --watch"）+ [N-RN-22.0]（"watch: mark as stable"） | Medium | — | 18.11 引入为实验；watch 是**进程重启**非热重载；测试另有 runner watch [N-TEST] |
| node22_node_run | node | 22 | `npm run <script>` → `node --run <script>`（免 npm 启动开销） | [N-CLI]（"Added in: v22.0.0"）+ [N-RN-22.0]（"cli: implement node --run …"） | Low | — | 不执行 npm 生命周期钩子，语义非全等价 |
| node22_require_esm | node | 22.12 | 只能 `import()` 引入 ESM → `require(esm)` | [N-RN-22.12]（"it is now no longer behind a flag on v22.x"；"can, however, throw ERR_REQUIRE_ASYNC_MODULE if … top-level await"） | High | — | 仅同步 ESM 图；含 TLA 抛 `ERR_REQUIRE_ASYNC_MODULE`；可用 `process.features.require_module` 检测 |
| node22_fs_cp | node | 22.3 | fs-extra/手写递归 copy → `fs/promises` `cp(src, dest, { recursive: true })` | [N-FS]（fsPromises.cp "Added in: v16.7.0"；"v22.3.0 This API is no longer experimental"） | Medium | — | 16.7–22.2 实验期；稳定门控取 22.3 |
| node22_partial_deep_strict_equal | node | 22.13 | 逐字段断言/子集比较手写 → `assert.partialDeepStrictEqual` | [N-ASSERT]（"Added in: v23.4.0, v22.13.0"；"v24.0.0, v22.17.0 … now Stable"） | Medium | — | 22.13–22.16 为 experimental；稳定门控 22.17/24；注意 v25 语义微调（同实例 Promise、Invalid Date 判等） |
| node22_websocket | node | 22 | ws 依赖 → 全局 WebSocket 客户端 | [N-GLOBALS]（Class: WebSocket "Added in: v21.0.0, v20.10.0"；v22.0.0 移除 `--experimental-websocket`；v22.4.0 "No longer experimental"）+ [N-RN-22.0]（"lib: enable WebSocket by default"） | Medium | — | 客户端 API；与 ws 库（服务端、更多选项）能力面不同 |
| node22_sqlite | node | 22.13 | better-sqlite3 依赖 → `node:sqlite`（22.13 起免 flag） | [N-SQLITE]（"Added in: v22.5.0"；"v23.4.0, v22.13.0 … no longer behind --experimental-sqlite but still experimental"；v25.7.0 RC） | Low | — | 22/24 内仍 experimental（RC 在 25.7）；require 时**必须** `node:` 前缀 [N-MODULES]；生产替换 better-sqlite3 前评估稳定性 |
| node21_style_text | node | 21.7（20.12 回移） | chalk/picocolors → `util.styleText` | [N-UTIL]（"Added in: v21.7.0, v20.12.0"；"v23.5.0, v22.13.0 styleText is now stable"） | Medium | — | v22.8.0/20.18.0 起才尊重 `isTTY`/`NO_COLOR`/`FORCE_COLOR`；稳定门控 22.13+ |
| node20_parse_args | node | 20（18.3 引入） | commander/yargs 最小场景 → `util.parseArgs` | [N-UTIL]（"Added in: v18.3.0, v16.17.0"；"v20.0.0 The API is no longer experimental"） | Medium | — | 只做"标志+位置参数"解析，不做子命令/帮助文本；复杂 CLI 仍需库 |
| node20_env_file | node | 20.6 | dotenv 依赖 → `--env-file` / `process.loadEnvFile()` | [N-CLI]（"--env-file … Added in: v20.6.0"；"v24.10.0, v22.21.0 … no longer experimental"）+ [N-PROCESS]（loadEnvFile "Added in: v21.7.0, v20.12.0"） | Medium | — | 文件不存在会报错（容错用 22.9 的 `--env-file-if-exists`）；.env 中 NODE_OPTIONS 不生效 [N-PROCESS] |
| node20_import_meta_resolve | node | 20.6 | 自写解析/正则推断路径 → `import.meta.resolve(specifier)`（同步返回字符串） | [N-ESM]（"v20.0.0, v18.19.0 … returns a string synchronously"；"v20.6.0, v18.19.0 No longer behind --experimental-import-meta-resolve"；Stability 1.2）+ [N-RN-24.0]（"esm: graduate import.meta properties"） | Medium | — | 非标准 `parentURL` 参数仍带标记；仅 ESM 内可用 |
| node20_mock_timers | node | 20.4 | sinon fake timers → `node:test` `mock.timers` | [N-TEST]（MockTimers "Added in: v20.4.0, v18.19.0"；"v23.1.0 The Mock Timers is now stable"） | Medium | — | 21.2/20.11 起参数改为选项对象（API 有过一次形态变化）；稳定门控 23.1 |
| node18_fetch | node | 18 | axios/node-fetch 依赖 → 全局 `fetch`（FormData/Headers/Request/Response 同批全局化） | [N-GLOBALS]（fetch "Added in: v17.5.0, v16.15.0"；"v18.0.0 No longer behind --experimental-fetch"；"v21.0.0 No longer experimental"）+ [N-RN-21.0]（"stable fetch and WebStreams"） | High | — | 基于 undici：无浏览器式 cookie jar、部分 Web 语义差异；18–20 标记仍为 experimental（严格生产门控取 21）；`AbortSignal.timeout`（v17.3.0）常与之搭配 [N-GLOBALS] |
| node18_test_runner | node | 18（20 稳定） | jest/mocha（无特殊能力需求时）→ `node:test` + `node:assert/strict` | [N-TEST]（"Added in: v18.0.0, v16.17.0"；"v20.0.0 The test runner is now stable"；describe "Added in: v22.0.0, v20.13.0"；快照 v22.3.0、v23.4.0 不再实验） | High | — | 稳定门控取 20；describe 语法 20.13+/22；快照/mocking 能力弱于 jest，迁移前对照用例面 |
| node17_structured_clone | node | 17 | `JSON.parse(JSON.stringify(x))` → `structuredClone` | [N-GLOBALS]（"Added in: v17.0.0"） | Medium | — | 保留 Map/Set/Date/循环引用；对函数与某些对象**抛错**而非忽略——JSON 方案会静默丢，clone 会响 |
| node16_node_prefix | node | 16（14.18 起） | `require('fs')` → `require('node:fs')` / `import node:fs` | [N-MODULES]（"v16.0.0, v14.18.0 Added node: import support to require(...)"；"`node:` 前缀绕过 require cache"） | Medium | — | `node:test`、`node:sqlite` 等模块 require 时**必须**带前缀 [N-MODULES]；unicorn/prefer-node-protocol 有 fixer 但非 typescript-eslint，按约束不填 |
| node16_esm_explicit_extensions | node | 16 | CJS 式无扩展相对导入 → ESM 相对导入**必须写全扩展名** | [N-ESM]（"A file extension must be provided when using the import keyword to resolve relative or absolute specifiers. Directory indexes … must also be fully specified."；"No default extensions"、"No folder mains"） | High | — | ESM/CJS 互操作边界：TS 侧需配合 `moduleResolution: nodenext/bundler`（bundler 允许无扩展名是 bundler 语义而非 Node 语义）；`.js` 后缀指向 `.ts` 源（`rewriteRelativeImportExtensions` [TS-B-5.7] 可自动改写） |

### 3.1 rationale 逐条（一句）与建议 category

- node24_regexp_escape — 内建转义消除手写函数的边界 bug（如 `-` 与 `]` 处理）。Language (V8)。
- node24_explicit_resource_mgmt — 给 ts52_using_declarations 补上运行时前提，跨轴成对出现。Language (V8)。
- node23_type_stripping — 砍掉"为跑 TS 而建"的构建步骤，Node 即运行时。Tooling。
- node22_array_from_async — 异步可迭代→数组一行完成，语义即文档。Language (V8)。
- node22_promise_with_resolvers — resolve/reject 与 promise 同作用域，事件驱动代码少一层闭包体操。Language (V8)。
- node22_set_methods — 集合运算内建化，正确性（不去重错、不漏元素）由引擎保证。Language (V8)。
- node22_fs_glob — 文件匹配是文件系统职责，交回内建减少依赖面。Standard Library。
- node22_watch_stable — 开发期监听重启零依赖，官方维护与 Node 版本同步。Tooling。
- node22_node_run — 省一次 npm 进程启动，脚本跑得更快更干净。Tooling。
- node22_require_esm — 存量 CJS 渐进吃 ESM 依赖，不必"全有或全无"。Modules。
- node22_fs_cp — 递归拷贝是 fs 一等能力，删掉 fs-extra 这类单函数依赖。Standard Library。
- node22_partial_deep_strict_equal — 子集断言让测试只锁关心的字段，减少脆弱断言。Test。
- node22_websocket — 客户端 WS 全局化，与浏览器 API 对齐。Standard Library。
- node22_sqlite — 内嵌 SQL 库零依赖可用（仍实验，谨慎评估）。Standard Library。
- node21_style_text — 颜色输出内建且默认尊重 NO_COLOR。Standard Library。
- node20_parse_args — 简单 CLI 的参数解析不需要一个框架。Standard Library。
- node20_env_file — 环境变量加载是运行时职责，去掉 dotenv 引导顺序问题。Tooling。
- node20_import_meta_resolve — 模块自解析取代对 import.meta.url 的手写路径拼接。Modules。
- node20_mock_timers — 时间相关用例不再依赖第三方 fake timers。Test。
- node18_fetch — HTTP 客户端进入运行时标准面，删依赖即删升级面。Standard Library。
- node18_test_runner — 测试器内建，零安装可跑，CI 最小化。Test。
- node17_structured_clone — 深拷贝保真运行时类型，替代会悄悄腐蚀数据的 JSON 往返。Standard Library。
- node16_node_prefix — node: 前缀显式声明内建意图且绕过 require cache 陷阱。Modules。
- node16_esm_explicit_extensions — 显式扩展名是 Node ESM 的硬规则，写对它才不依赖 bundler 兜底。Modules。

## 4. autofix 字段结论（PR7 采用）

- **唯一核实为真实 fixer 的映射**：`ts50_verbatim_module_syntax` → `typescript-eslint:consistent-type-imports`
  - [TSE-CTI] 页面原文徽标："Some problems reported by this rule are automatically fixable by the --fix ESLint command line option."（且有 `fixStyle: separate-type-imports | inline-type-imports` 选项）。
  - 可并列 `typescript-eslint:no-import-type-side-effects`（[TSE-NITSE] 同样标注 autofix），用于清一化 side-effect 导入写法。
- 其余 56 条：typescript-eslint 无对应真实 autofix 规则 → **autofix: null**。
- 备注：`unicorn/prefer-node-protocol`（node: 前缀）与 `n/prefer-promises/fs` 等存在 fixer，但非 typescript-eslint，按任务约束不填（可写入 dataset 的 references 供人工参考）。

## 5. UNVERIFIED 清单（未能官方原文核实，D8 纪律）

1. **Array.fromAsync / Promise.withResolvers / Set methods 的精确 Node 首发版本**：nodejs.org 未在 v21/v22 release notes 文中点名这些特性；MDN 兼容表在两次抓取中均未渲染（"See full compatibility" 折叠）。已核实事实仅到：Node 22 = V8 12.4.254.14 [N-RN-22.0]、Node 21 = V8 11.8 [N-RN-21.0]、MDN Baseline 时间（fromAsync 2024-01、withResolvers/groupBy 2024-03、Set methods 2024-06）。**处置：门控保守取 22，details 注明推断链**。
2. **Object.groupBy/Map.groupBy 的 Node 运行时落地**：同上（TS 5.4 类型支持已由官方标题证实 [TS-B-5.4]）；保守取 22。
3. **TS 4.9–5.8 的准确发布日期**：typescriptlang.org 各 release notes 页面未展示日期（仅 5.9/6.0/7.0 有官方公告日期：2025-08-01 / 2026-03-23 / 2026-07-08 [TS-BLOG]）。数据集不需要日期，不采。
4. **任务清单中"任意枚举成员名（5.0）"**：5.0 官方 notes 无此条目；真实相关特性为 "All enums Are Union enums" + Breaking Changes 下的 "Enum Overhaul"（已按核实内容重写为 ts50_enum_union_types）。
5. **WebSocket 独立 API 文档页**（nodejs.org/docs/latest/api/websocket.html）抓取 404；版本史改以 globals.html 的 Class: WebSocket 条目为准（已核实）。
6. **`--watch` 在 18.x 的完整行为细节**：Added 版本（v18.11.0 changelog "cli: add --watch"）与 stable（v22.0.0 "watch: mark as stable"）已核实；18 期间的限制细节未逐条核实（不影响门控）。
7. **TS 5.1 的两个候选**（Easier Implicit Returns / Unrelated Types for Getters and Setters，[TS-RN-5.1] 已核实标题）未入选 57 条，因迁移价值低——如需凑满 3.9–3.14 式"每版本 ≥3 条"分布可启用。

## 6. 对 PR7 的落地建议

- **条数分布**：57 条候选远超 ~40 需求；建议按 impact 优先保 High/Medium，Low 条目按版本分布补齐（TS 轴 5.1 空档见 §5 第 7 条）。
- **跨轴联动**（details 必须写清，防止 agent 只看单轴）：`ts52_using_declarations`↔`node24_explicit_resource_mgmt`；`ts58_erasable_syntax_only`↔`node23_type_stripping`；`ts50_verbatim_module_syntax`↔`node23_type_stripping`（Node 官方推荐 tsconfig 原文 [N-TS] 三者同现）；`ts53_import_attributes`↔Node 22 移除 assert [N-RN-22.0]；`node22_require_esm`↔`ts58_require_esm_nodenext`。
- **示例（before/after）**：表内"旧→新"即 before→after 骨架，PR7 编写时逐条落成 examples 数组（≥1 组，schema 要求）。
- **门控保守原则**：Node 轴凡"实验→稳定"分界的条目，since_version 建议取稳定点（表内已标注），details 保留引入点版本作为背景。
- **检测器提醒**（需求 5）：TS 轴 provenance 为"devDependencies 范围推断（中保真）"；Node 轴 engines.node 为 advisory——两条警示文案应引用本研究的版本门版图（§1）。
