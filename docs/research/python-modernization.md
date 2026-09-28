# Python 现代化规则候选研究（3.9 → 3.15）

> **Status**: Researched | **Created**: 2026-09-29 | **上游**: PRD-0000（Requirement 3 / PR6 内容填充前置研究）
> **候选总条数**: **88**（3.9–3.14 语言/库惯用法 70 条 + 3.15 占位 8 条 + Tooling 10 条；满足 PRD"≥63 条、3.9–3.14 每版本 ≥3 条"）

## 0. 信源与核实方法（硬要求执行记录）

全部结论仅采信官方权威信源，逐一于 2026-09-29 抓取核实：

| 信源 | 用途 | 核实方式 |
| --- | --- | --- |
| `docs.python.org/3/whatsnew/3.9.html` … `3.14.html` | 特性清单、锚点、PEP 编号 | WebFetch 逐页抓取（当前 3.14.7 文档集） |
| `docs.python.org/3.15/whatsnew/3.15.html` | 3.15 候选 | WebFetch（页面自标 **3.15.0rc2 + "draft form"**） |
| `peps.python.org/pep-0561/` | py.typed | WebFetch 全文 |
| `docs.python.org/3/library/typing.html` | 弃用别名时间线 | 选择器抓取 `#deprecated-aliases` 全文 |
| `docs.astral.sh/ruff/rules/<slug>` | ruff UP 规则号 | 19 个规则页逐一 WebFetch；另以 `astral-sh/ruff` 仓库 `codes.rs`（main 分支）交叉核对完整 UP 代码表 |
| `packaging.python.org`（3 页） | Tooling 引用 | tinyfish 全文抓取 |
| `docs.pytest.org`、`docs.astral.sh/uv` | Tooling 引用 | WebFetch |

**本页所有 PEP 编号、whatsnew 锚点、ruff 规则号均来自上述抓取，无一凭记忆书写。** 仅有意未核实项以 `UNVERIFIED` 标注（见 §9，本轮为 0 条硬性未决，仅 2 条次要注意点）。

**表头约定**（与 PRD schema 对应）：

- `id`：snake_case，语言内唯一（建议值，落库时可调整）。
- `since_version`：major.minor。数据集地板为 3.9；唯一例外 `lru_cache_no_parens`（3.8，ruff 官方页明示 "since Python 3.8"），落库时可归并入 3.9 组。
- `autofix`：`ruff:UPxxx`（已核实存在）或 `null`。**UP038（NonPEP604Isinstance）已被 ruff 官方移除（since 0.13.0），全表不引用。**
- `impact`：Critical=几乎每个项目 / High=常见 5-20 处 / Medium=规律 1-5 处 / Low=罕见。
- 引用格式：`PEP <n>`；`whatsnew/3.X#<锚点>`（锚点为抓取核实过的 HTML id）；未确认到稳定锚点的仅写 `whatsnew/3.X（某节）`。

**版本基线**（抓取自 whatsnew 各页）：3.9（2020-10-05）、3.10（2021-10-04）、3.11（2022-10-24）、3.12（2023-10-02）、3.13（2024-10-07）、3.14（2025-10-07，当前稳定）；3.15 处于 rc2（PRD/devguide：final 预计 2026-10-01）。

---

## 1. Python 3.9（11 条）

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `dict_merge_operator` | 3.9 | `{**a, **b}` / `a.copy(); a.update(b)` → `a \| b`、`a \|= b` | PEP 584；whatsnew/3.9#dictionary-merge-update-operators | Critical | null | `\|` 返回新 dict（原对象不动）、`\|=` 原地；右操作数须为 Mapping，否则 TypeError |
| `builtin_generic_annotation` | 3.9 | `typing.List[int]`/`Dict`/`Tuple` → `list[int]`/`dict`/`tuple` | PEP 585；whatsnew/3.9#type-hinting-generics-in-standard-collections；typing 文档 #deprecated-aliases（官方：3.9 起弃用、类型检查器应对 min≥3.9 项目标红） | Critical | ruff:UP006 | **pydantic 等运行时求值注解库下改写不安全**（ruff 官方 UP006 页明示，仅 <3.9 目标标记 unsafe；≥3.9 运行时合法）；伴随弃用 import 由 ruff:UP035 清理 |
| `typing_text_str` | 3.9 | `typing.Text` → `str` | typing 文档 #deprecated-aliases（"Deprecated since version 3.11"，Python 2 兼容残留） | Medium | ruff:UP019 | 官方明示"暂无移除计划"，但建议全部换 `str` |
| `str_removeprefix_removesuffix` | 3.9 | `s[len(p):] if s.startswith(p) else s` → `s.removeprefix(p)`（suffix 同理） | PEP 616；whatsnew/3.9#new-string-methods-to-remove-prefixes-and-suffixes | High | null | 不匹配时返回原串副本（安全）；bytes/bytearray/UserString 同步支持（whatsnew 原文）；比手写切片少一处 off-by-one 风险 |
| `functools_cache` | 3.9 | `@functools.lru_cache(maxsize=None)` → `@functools.cache` | whatsnew/3.9（functools 节）；ruff UP033 页确认"requires Python 3.9+，drop-in replacement" | Medium | ruff:UP033 | 语义等价（无界缓存）；仍需注意实例方法缓存导致内存滞留 |
| `zoneinfo_stdlib` | 3.9 | `pytz` → `zoneinfo.ZoneInfo` | PEP 615；whatsnew/3.9#zoneinfo | Medium | null | 无系统 tzdata 数据库的平台（如 Windows）需安装 `tzdata` 包兜底（zoneinfo 模块文档语义）；pytz 的 `localize()` 惯用法不迁移 |
| `annotated_metadata` | 3.9 | 自造元数据包装类型 → `typing.Annotated[T, meta]` | PEP 593；whatsnew/3.9（typing 节） | Medium | null | 运行时 `Annotated` 保留元数据供框架（pydantic/FastAPI）消费；裸类型检查器忽略之 |
| `asyncio_to_thread` | 3.9 | `loop.run_in_executor(None, fn, *args)` → `asyncio.to_thread(fn, *args)` | whatsnew/3.9（asyncio 节） | Medium | null | 自动传播当前 contextvars；kwargs 可直接传（executor 版不行） |
| `math_lcm_multiarg` | 3.9 | 手写 lcm / `reduce(gcd, xs)` → `math.lcm(*xs)`、`math.gcd(*xs)` | whatsnew/3.9（math 节，官方原文"math.gcd() accepts multiple arguments"） | Low | null | 多参数化是 3.9 新增；单参/两参版本 3.5 已有 |
| `graphlib_topological_sort` | 3.9 | 手写 Kahn/DFS 拓扑排序 → `graphlib.TopologicalSorter` | whatsnew/3.9#graphlib | Low | null | `prepare()` 后支持并发"就绪即取"增量遍历 |
| `lru_cache_no_parens` | 3.8* | `@lru_cache()` → `@lru_cache` | ruff UP011 页（"since Python 3.8"） | Medium | ruff:UP011 | *低于数据集地板 3.9；落库时可并入 3.9 组或删除 |

## 2. Python 3.10（13 条）

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `union_type_annotation` | 3.10 | `Optional[X]` / `Union[A, B]` → `X \| None` / `A \| B` | PEP 604；whatsnew/3.10#pep-604-new-type-union-operator | Critical | ruff:UP007（Union）/ UP045（Optional） | **pydantic 等运行时求值注解库下改写不安全**：`types.UnionType` 3.10 才存在，ruff 官方两页均明示 <3.10 目标 fix unsafe（可用 `lint.pyupgrade.keep-runtime-typing=true` 关闭） |
| `structural_pattern_matching` | 3.10 | if/elif 链 + 手工解包/取属性 → `match`/`case` 模式匹配 | PEP 634/636；whatsnew/3.10#pep-634-structural-pattern-matching | High | null | 最大陷阱：case 后裸名是**捕获模式**（永不失败），字面量须用 `case 404:` 或点路径；类模式依赖 `__match_args__` |
| `zip_strict` | 3.10 | `zip(a, b)`（静默截断） → `zip(a, b, strict=True)` | PEP 618；whatsnew/3.10（Summary/其他语言特性） | High | null | 长度不等抛 ValueError；依赖旧行为（有意截断）的代码会破；ruff 检查项为 B905（bugbear 家族，非 UP） |
| `dataclass_slots` | 3.10 | 手写 `__slots__` 继承 dataclass → `@dataclass(slots=True)` | whatsnew/3.10#slots | High | null | slots=True 会**重建类**（装饰器返回新类），与继承链、`super()`、默认可变字段默认值工厂交互需回归测试 |
| `dataclass_kw_only` | 3.10 | 手写 `__init__` 做 kwargs 校验 → `@dataclass(kw_only=True)` / 字段级 `KW_ONLY` 哨兵 | whatsnew/3.10#keyword-only-fields | Medium | null | `KW_ONLY` 哨兵字段不成为字段；与继承字段排序规则需留意 |
| `parenthesized_context_managers` | 3.10 | `contextlib.ExitStack` 或 `with a as x: with b as y:` → `with (a as x, b as y):` | whatsnew/3.10#parenthesized-context-managers | Medium | null | 纯语法糖；旧解析器（<3.10）报 SyntaxError |
| `explicit_encoding_open` | 3.10 | `open(f)` 不传 encoding → `open(f, encoding="utf-8")` | PEP 597；whatsnew/3.10#optional-encodingwarning-and-encoding-locale-option | High | null | **前瞻翻转**：3.15 PEP 686 使 UTF-8 成为默认后，"必须显式传 utf-8"的建议失效（详见 §7/§8）；Windows locale 场景本就应显式 |
| `isinstance_union` | 3.10 | `isinstance(x, (int, str))` → `isinstance(x, int \| str)` | PEP 604；whatsnew/3.10#pep-604-new-type-union-operator | Medium | null | **ruff UP038 已移除（since 0.13.0）**，勿引用；元组形式并非错误，收益主要是与注解书写统一 |
| `type_alias_explicit` | 3.10 | 裸赋值 `X = List[int]`（语义含混） → `X: TypeAlias = list[int]` | PEP 613；whatsnew/3.10#pep-613-typealias | Medium | null | 过渡形态：3.12 起 `TypeAlias` 弃用、终态为 `type` 语句（见 3.12 `type_alias_statement`）；min≥3.12 项目应一步到位（ruff:UP040） |
| `typeguard_narrowing` | 3.10 | 自写 `def is_str_list(v) -> bool` → 返回 `TypeGuard[list[str]]` | PEP 647；whatsnew/3.10#pep-647-user-defined-type-guards | Low | null | TypeGuard 只保证 if 分支收窄（else 分支不反转）；3.13 有更好用的 TypeIs（见 3.13） |
| `paramspec_callable` | 3.10 | 装饰器签 `Callable[..., Any]`（丢签名） → `ParamSpec`/`Concatenate` | PEP 612；whatsnew/3.10#pep-612-parameter-specification-variables | Low | null | 需 pyright/mypy 支持；与 3.12 PEP 695 `**P` 语法可组合 |
| `int_bit_count` | 3.10 | `bin(x).count("1")` → `x.bit_count()` | whatsnew/3.10（Other Language Changes） | Low | null | 等价替换，无行为差 |
| `itertools_pairwise` | 3.10 | `zip(xs, xs[1:])` → `itertools.pairwise(xs)` | whatsnew/3.10（itertools 节） | Medium | null | pairwise 接受任意可迭代（迭代器/生成器）；切片法要求序列 |

## 3. Python 3.11（14 条）

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `tomllib_stdlib` | 3.11 | 第三方 `tomli` / 手写解析 → `tomllib.load(fh)`（`rb` 模式） | PEP 680；whatsnew/3.11（new-modules 节） | High | null | **只读不写**（写用 `tomli-w`）；需兼容 3.10- 时惯用法 `try: import tomllib except ModuleNotFoundError: import tomli as tomllib` |
| `asyncio_task_group` | 3.11 | `create_task` + `gather` + 手动取消样板 → `async with asyncio.TaskGroup() as tg: tg.create_task(...)` | whatsnew/3.11（asyncio 节，官方措辞"recommended over create_task+gather"） | High | null | 任一子任务失败→取消兄弟任务并以 **ExceptionGroup** 传播（错误处理模式改变，需 except* 配合，见 `exception_group_except_star`） |
| `exception_group_except_star` | 3.11 | `gather(return_exceptions=True)` 手工聚合/丢错 → `raise ExceptionGroup` + `except* XxxError:` | PEP 654；whatsnew/3.11#pep-654-exception-groups-and-except | Medium | null | `except*` 与普通 `except` **不能混用于同一 try**；层级注意：`ExceptionGroup` 继承 Exception、`BaseExceptionGroup` 继承 BaseException，`except Exception` 捕不到后者；`except*` 处理器收到的是按类型拆分后的子组 |
| `exception_add_note` | 3.11 | 自定义异常子类携带上下文字段 → `exc.add_note("...")` | PEP 678；whatsnew/3.11#pep-678-exceptions-can-be-enriched-with-notes | Medium | null | note 进入默认 traceback 展示；不可重复"覆盖"，多次调用是追加 |
| `typing_self` | 3.11 | 返回类型写类名字符串/自引用 TypeVar → `typing.Self` | PEP 673；whatsnew/3.11#pep-673-self-type | Medium | null | 适用于 `__enter__`/`classmethod` 构造器/流式 builder；pydantic v2 运行时支持 Self |
| `typeddict_required_not_required` | 3.11 | `total=` 全有全无 → 逐键 `Required[...]`/`NotRequired[...]` | PEP 655；whatsnew/3.11#pep-655-marking-individual-typeddict-items-as-required-or-not-required | Medium | null | 仅静态检查层契约，运行时不校验键 |
| `typing_assert_never` | 3.11 | `else: raise AssertionError("unreachable")` → `typing.assert_never(x)` | whatsnew/3.11（typing 节） | Medium | null | 收穷尽性检查红利（配合 Literal enum/match 分支收窄）；运行时抛 AssertionError |
| `enum_strenum` | 3.11 | `class X(str, enum.Enum)` → `enum.StrEnum` | whatsnew/3.11（enum 节） | Medium | ruff:UP042 | **行为差**（ruff UP042 官方 caveat）：3.11 起 str-Enum 混类 f-string 输出 `X.MEMBER`，StrEnum 输出小写成员值 `member`——依赖任一旧行为的展示/序列化需回归 |
| `datetime_fromisoformat_iso8601` | 3.11 | `strptime` 手写 ISO 格式串 → `datetime.fromisoformat()` 解析大多数 ISO 8601 | whatsnew/3.11（datetime 节） | Medium | null | 官方措辞是"most ISO 8601 formats"而非全部，边缘格式仍需 strptime 兜底 |
| `asyncio_timeout_cm` | 3.11 | `await asyncio.wait_for(coro(), t)` → `async with asyncio.timeout(t): await ...` | whatsnew/3.11（asyncio 节，官方措辞"recommended over wait_for"） | Medium | null | 保护**代码块**而非单个 await；`wait_for` 未被弃用；超时抛内建 `TimeoutError`（`asyncio.TimeoutError` 3.11 起为其别名——ruff UP041 页核实） |
| `asyncio_runner` | 3.11 | 手写 `new_event_loop`/`run_until_complete`/`shutdown_asyncgens` → `asyncio.Runner` | whatsnew/3.11（asyncio 节） | Low | null | 同一 loop 反复跑多个协程场景；一次性脚本继续用 `asyncio.run` |
| `timeout_error_alias` | 3.11 | `asyncio.TimeoutError` / `socket.timeout` → 内建 `TimeoutError` | whatsnew/3.11；ruff UP041 页（socket.timeout 自 3.10、asyncio.TimeoutError 自 3.11 为别名） | Medium | ruff:UP041 | 捕获侧统一；`ssl.SSLError` 不属此列（是 OSError 子类，UP024 另管 OSError 别名） |
| `datetime_utc_alias` | 3.11 | `datetime.timezone.utc` → `datetime.UTC` | whatsnew/3.11（datetime 节） | Low | ruff:UP017 | 纯别名新增，旧写法未弃用 |
| `typevar_tuple_variadic` | 3.11 | 多 TypeVar + 重载模拟张量形状 → `TypeVarTuple`/`*Ts` | PEP 646；whatsnew/3.11#pep-646-variadic-generics | Low | null | 主要收益在 numpy/数组库签名；3.12 PEP 695 提供 `*Ts` 内联语法 |

## 4. Python 3.12（10 条）

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `type_alias_statement` | 3.12 | `X: TypeAlias = ...` / `X = TypeAliasType(...)` → `type X = ...` | PEP 695；whatsnew/3.12#pep-695-type-parameter-syntax | High | ruff:UP040 | ruff UP040 官方 caveat：旧式简单别名可被 `isinstance` 使用，`type` 语句产物**isinstance 必 TypeError**；TypeAliasType 方差/边界信息可能丢失；fix 标 unsafe |
| `generic_class_params` | 3.12 | `T = TypeVar("T"); class C(Generic[T])` → `class C[T]:` | PEP 695；whatsnew/3.12#pep-695-type-parameter-syntax | Medium | ruff:UP046 | PEP 695 用**推断方差**，与显式 `covariant/contravariant` 不等价；旧 TypeVar 定义不会自动删除（ruff unsafe caveat） |
| `generic_function_params` | 3.12 | `T = TypeVar("T"); def f(x: T) -> T` → `def f[T](x: T) -> T` | PEP 695；whatsnew/3.12#pep-695-type-parameter-syntax | Medium | ruff:UP047 | 同上（方差推断）；部分类型检查器对 PEP 695 支持仍在完善（ruff 官方注） |
| `typeddict_kwargs_unpack` | 3.12 | `**kwargs: Any` / `**kwargs: Unpacked[dict]` → `**kwargs: Unpack[Movie]`（TypedDict 逐键） | PEP 692；whatsnew/3.12#pep-692-using-typeddict-for-more-precise-kwargs-typing | Medium | null | 纯类型检查器层面；运行时不校验 kwargs |
| `override_decorator` | 3.12 | 拼写错误静默生成新方法（如 `get_colour`） → `@typing.override` 让检查器报错 | PEP 698；whatsnew/3.12#pep-698-override-decorator-for-static-typing | Medium | null | 仅静态检查；运行时无操作 |
| `fstring_quote_reuse` | 3.12 | 为嵌套引号换引号/提取中间变量 → f-string 内重用同类引号、任意嵌套与多行表达式 | PEP 701；whatsnew/3.12#pep-701-syntactic-formalization-of-f-strings | Low | null | 目标 <3.12 语法报错；tokenization 变化影响极少数依赖旧 token 流的工具（官方 whatsnew 提及） |
| `itertools_batched` | 3.12 | 手写 `range(0, n, k)` 分块循环 → `itertools.batched(it, k)` | whatsnew/3.12（itertools 节） | Medium | null | 末批可短；3.13 起有 `strict=True`（见 3.13）；ruff 检查项 B911（非 UP） |
| `pathlib_walk` | 3.12 | `os.walk` 拼 `Path` → `Path.walk()` 直接产 Path | whatsnew/3.12（pathlib 节） | Medium | null | 语义与 os.walk 有差（`follow_symlinks`/`on_error` 参数设计不同）；默认不跟随符号链接目录 |
| `distutils_removed` | 3.12 | `import distutils` → 迁移 setuptools 兼容层或替代命令 | PEP 632；whatsnew/3.12（distutils 节；3.10 whatsnew 标注移除计划） | Medium | null | 3.12 起彻底移除；venv 不再预装 setuptools（官方 whatsnew 原话） |
| `typing_abc_direct_import` | 3.12 | `typing.Hashable`/`typing.Sized` → `collections.abc.Hashable`/`Sized` | typing 文档 #deprecated-aliases（两者"Deprecated since version 3.12"）；PEP 585 背景同页 | Low | ruff:UP035 | 与 UP006 同族；`typing.ByteString` 更严：3.9 弃用、**3.17 移除**（官方同页明示），建议换 `Buffer`（PEP 688） |

## 5. Python 3.13（11 条）

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `typing_type_is` | 3.13 | `TypeGuard`（单向收窄） → `TypeIs`（if/else 双向更直觉收窄） | PEP 742；whatsnew/3.13（typing 节） | Medium | null | TypeIs 遵循 isinstance 式收窄（含 else 分支反转）；两者都仅静态层 |
| `typeddict_readonly` | 3.13 | 冻结/不可变 TypedDict 约定（docstring/前缀命名） → `ReadOnly[...]` | PEP 705；whatsnew/3.13（typing 节） | Medium | null | 仅静态检查器契约；与协变结合是主要动机（官方 PEP） |
| `typevar_defaults` | 3.13 | 手写 `X = DefaultT` 占位别名 → `T = TypeVar("T", default=...)` | PEP 696；whatsnew/3.13（typing 节） | Medium | null | ParamSpec/TypeVarTuple 同样支持默认值 |
| `generator_default_params` | 3.13 | `Generator[int, None, None]` → `Generator[int]`（Send/Return 默认 None） | PEP 696（collections.abc 泛型默认值）；whatsnew/3.13（typing 节"Changed in 3.13: Default values for the send and return types"） | Medium | ruff:UP043 | ruff UP043 官方页（unnecessary-default-type-args）确认 3.13+ 门控、fix 常驻；AsyncGenerator[int, None] 同理可省第二参 |
| `warnings_deprecated_decorator` | 3.13 | 自建 DeprecationWarning 装饰器 → `@warnings.deprecated("msg")` | PEP 702；whatsnew/3.13（warnings 节） | Medium | null | 运行时触发 DeprecationWarning 且静态检查器可感知弃用；对属性/类/函数均可用 |
| `copy_replace` | 3.13 | 只为 `_replace` 手写 `def replace(...)` → `copy.replace(obj, **changes)`（走 `__replace__` 协议） | whatsnew/3.13（copy 节） | Medium | null | 需类型实现 `__replace__`；标准库 namedtuple/dataclass/datetime/ZoneInfo 等已覆盖（官方列表），自定义不可变类型需自行实现 |
| `process_cpu_count` | 3.13 | `os.cpu_count()` → `os.process_cpu_count()`（受 affinity 约束的真实可用核） | whatsnew/3.13（os 节） | Low | null | 语义不同于系统总核数；配 `PYTHON_CPU_COUNT` 覆盖（官方同节） |
| `batched_strict` | 3.13 | `batched` 末批静截断 → `batched(it, k, strict=True)` | whatsnew/3.13（itertools 节） | Low | null | 不整除抛 ValueError；ruff 检查项 B911（非 UP） |
| `dead_batteries_replacement` | 3.13 | `import cgi`/`crypt`/`telnetlib`/`pipes` 等 19 模块 → PEP 594 给出的替代方案 | PEP 594；whatsnew/3.13#pep-594-remove-dead-batteries-from-the-standard-library | Medium | null | 3.13 起移除；逐模块替代清单见 PEP 594 各小节 |
| `locals_snapshot_semantics` | 3.13 | 依赖 `locals()` 写回优化作用域（旧未定义行为） → PEP 667 定义的快照/`f_locals` 写穿代理 | PEP 667；whatsnew/3.13#defined-mutation-semantics-for-locals | Low | null | 语义收紧条目：调试器/追踪库最相关；普通代码勿依赖 locals() 修改传播 |
| `anystr_type_params` | 3.13 | `typing.AnyStr` → PEP 695 类型参数显式 `[str]`/`[bytes]` | typing 文档 #deprecated-aliases（"Deprecated since 3.13, will be removed in 3.18"） | Low | null | 与 `typing.Pattern/Match`（3.9 弃用）同页；AnyStr 将于 3.18 移除，与 `typing.ByteString`（3.9 弃用、3.17 移除）同为已公告移除版本的 typing 顶层弃用别名 |

## 6. Python 3.14（11 条）

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `deferred_annotations` | 3.14 | 前向引用字符串 `"ClassName"` / 全面 `from __future__ import annotations` → 注解惰性求值默认生效，直接写裸名 | PEP 649/749；whatsnew/3.14#pep-649-pep-749-deferred-evaluation-of-annotations | High | ruff:UP037（部分：引号注解） | 移植注意（whatsnew porting 节）：`__annotations__` 访问语义变化，直接读取的库应改用 `annotationlib`；`from __future__ import annotations` 仍合法但不再必要 |
| `t_string_templates` | 3.14 | `string.Template`/手写拼接防注入 → `t"...{x}..."` 模板字面量（`string.templatelib.Template`） | PEP 750；whatsnew/3.14#pep-750-template-string-literals | Medium | null | t-string **不是 str** 而是 Template 对象（静态段+Interpolation 列表），需渲染器消费（SQL/HTML/shell 安全构建）；与 f-string 用途互补非替代 |
| `except_multiple_no_brackets` | 3.14 | `except (A, B):` → `except A, B:` | PEP 758；whatsnew/3.14#pep-758-allow-except-and-except-expressions-without-brackets | Low | null | 仅当**不带 as 子句**时可省括号（官方原文）；带 as 仍需括号 |
| `map_strict` | 3.14 | `map(f, a, b)` 静默截断 → `map(f, a, b, strict=True)` | whatsnew/3.14（built-ins 节） | Medium | null | 与 3.10 zip strict 对齐；ruff 检查项 B912（非 UP） |
| `pathlib_copy_move` | 3.14 | `shutil.copy/copytree/move` 混搭 Path → `Path.copy()/copy_into()/move()/move_into()` | whatsnew/3.14（pathlib 节） | Medium | null | 3.14 新增（3.13 无此 API——whatsnew/3.13 pathlib 节官方明确"no copy/move in 3.13"）；另有 `Path.info` 属性 |
| `subinterpreters_stdlib` | 3.14 | `multiprocessing`（进程开销） → `concurrent.interpreters`（每解释器独立 GIL） | PEP 734；whatsnew/3.14#pep-734-multiple-interpreters-in-the-standard-library | Low | null | 数据共享受限（序列化/Buffer 共享）；C 扩展需 per-interpreter 兼容 |
| `pickle_protocol5_default` | 3.14 | 显式 `pickle.dumps(o, protocol=5)` → 默认即 5（可省） | whatsnew/3.14（pickle 节，官方原文"default protocol now 5"） | Low | null | 旧解释器若不识别协议 5 数据需自行处理；变化仅在默认值，旧协议仍可显式指定 |
| `typing_union_repr_unified` | 3.14 | 依赖 `repr(Union[int, str])`/运行时 Union 缓存行为 → 与 `int \| str` 完全统一（"int \| str"），旧式联合不再缓存 | whatsnew/3.14（typing 节） | Low | null | 快照测试/序列化类型 repr 的代码会看到输出变化（官方 porting 注意点） |
| `zstandard_stdlib` | 3.14 | 第三方 `zstandard` 包 → `compression.zstd`（tarfile/zipfile/shutil 亦支持 zstd） | PEP 784；whatsnew/3.14#pep-784-zstandard-support-in-the-standard-library | Low | null | API 与第三方包不一一对应；旧版 Python 回退需条件依赖 |
| `forkserver_start_method` | 3.14 | Unix 默认 `fork` → 默认 `forkserver`（concurrent.futures 与 multiprocessing 同步切换） | whatsnew/3.14（concurrent.futures/multiprocessing 节，官方原文"forkserver now default start method on Unix (replacing fork)"） | Medium | null | **移植注意**（默认值切换）：依赖 fork 语义（如写时复制共享大对象）的代码需显式设回 `fork` |
| `free_threaded_supported` | 3.14 | 仅实验性 free-threaded 构建 → 官方支持级（PEP 779），可作部署目标 | PEP 779；whatsnew/3.14（free-threaded 节 #whatsnew314-free-threaded-now-supported） | Low | null | 信息条目：C 扩展/线程安全审计建议；单线程开销约 5–10%（官方口径） |

## 7. Python 3.15 占位候选（8 条）——whatsnew 为 **rc2 + draft** 状态

> 来源：`docs.python.org/3.15/whatsnew/3.15.html`，页面自标 **Python 3.15.0rc2** 且带"this document is currently in draft form"免责声明。**final（预计 2026-10-01）后须逐条复核再转正**（PRD 既定待办）。

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `utf8_default_encoding` | 3.15 | 全面防御式 `encoding="utf-8"` → 默认编码即 UTF-8（PEP 686 落地），"必须显式传 utf-8"建议**翻转失效** | PEP 686；whatsnew/3.15（Other language changes 节原文："I/O operations without an explicit encoding … will use UTF-8"） | High | null | **rc 状态**。非 UTF-8 场景仍须显式；locale 旧行为用 `encoding="locale"`；可 `PYTHONUTF8=0`/`-X utf8=0` 退出；Windows 影响最大。对 3.10 `explicit_encoding_open` 规则形成版本翻转（该条 caveat 已标注） |
| `sentinel_builtin` | 3.15 | `_MISSING = object()` 自造哨兵 → 内建 `sentinel` 类型（可 pickle、保恒等、支持 `\|`） | PEP 661；whatsnew/3.15#pep-661-add-sentinel-built-in-type | Medium | null | **rc 状态**。自造 object 哨兵不可 pickle、无稳定 repr——即本条机制 |
| `frozendict_builtin` | 3.15 | `types.MappingProxyType` / 第三方 immutables → 内建 `frozendict`（不可变、可哈希、保序） | PEP 814；whatsnew/3.15#pep-814-add-frozendict-built-in-type | Medium | null | **rc 状态**。MappingProxyType 是视图非独立容器（底层可变），机制差异即动机 |
| `lazy_import_keyword` | 3.15 | 函数内局部 import / 自建懒加载代理 → `lazy import` 软关键字（首用才加载） | PEP 810；whatsnew/3.15#pep-810-explicit-lazy-imports | Medium | null | **rc 状态**。注意与本页 3.15 的 PEP 810 编号来自 whatsnew 页面本身；`import` 行写入 `.pth` 同步软弃用（PEP 829，同页） |
| `comprehension_unpacking` | 3.15 | 手写嵌套循环展平/`itertools.chain.from_iterable` → `[*xs for xs in lists]`、`{**d for d in dicts}` | PEP 798；whatsnew/3.15#pep-798-unpacking-in-comprehensions | Medium | null | **rc 状态**。生成器表达式同样适用（官方示例含 genexp） |
| `typing_type_form` | 3.15 | 注解"接受类"用 `type[T]` 无法表达元类/工厂 → `typing.TypeForm` | PEP 747；whatsnew/3.15（typing 节） | Low | null | **rc 状态**。对齐 mypy/typing_extensions 既有 TypeForm 实验 |
| `closed_typeddict` | 3.15 | 手写 `__extra_items__`/自定义校验拒绝未知键 → `class TD(TypedDict, closed=True, extra_items=...)` | PEP 728；whatsnew/3.15（typing 节） | Low | null | **rc 状态**。仅静态检查器契约 |
| `taskgroup_cancel_api` | 3.15 | 外部任务取消 TaskGroup 的绕行（取消内部任务/事件信号） → `TaskGroup.cancel()` | whatsnew/3.15（asyncio 节） | Low | null | **rc 状态**。官方描述为"early task-group termination" |

## 8. Tooling 类（10 条，since_version=3.9）

> 引用主体为 packaging.python.org 官方页（2026-09-29 全文抓取）；uv/ruff/pytest 用各自官方文档站。

| id | 版本 | 旧→新 | 引用 | impact | autofix | caveat |
| --- | --- | --- | --- | --- | --- | --- |
| `pyproject_single_config` | 3.9 | 配置散落 setup.py/setup.cfg/pytest.ini/.flake8 → 统一 `pyproject.toml`（`[build-system]`+`[project]`+`[tool.*]`） | packaging.python.org/en/latest/guides/writing-pyproject-toml/（官方：`[build-system]` "should always be present"，新项目用 `[project]` 表） | Critical | null | Poetry <2.0（2025-01 前）不用 `[project]` 表——官方页明示例外 |
| `no_setup_py_cli` | 3.9 | `python setup.py install/develop/sdist/bdist_wheel` → `python -m pip install .` / `python -m pip install --editable .` / `python -m build` | packaging.python.org/en/latest/discussions/setup-py-deprecated/（官方表格：四命令 **MUST NOT** 再运行） | Critical | null | **精确语义**：官方原话 "No, setup.py and Setuptools are not deprecated"——弃用的是 setup.py 作为**命令行调用**；`python setup.py install` 在 setuptools 58.3.0 弃用 |
| `requires_python_metadata` | 3.9 | 仅用 trove classifiers 声明版本支持 → `requires-python = ">=3.9"` | packaging.python.org/en/latest/guides/writing-pyproject-toml/#requires-python（官方：classifiers "only used for searching and browsing"，安装约束靠 requires-python） | High | null | 该字段同时是本工具 Python 检测器 T2 的主信源（PRD Req 5/6） |
| `src_layout` | 3.9 | flat 布局（包在仓库根） → src 布局（`src/<pkg>/`） | packaging.python.org/en/latest/discussions/src-layout-vs-flat-layout/（官方列三条行为差异） | High | null | 需 editable install 工作流；收益：防 cwd 遮蔽已安装包、约束可导入面（官方原话） |
| `py_typed_marker` | 3.9 | 包仅藏类型于注释/pyi 不暴露 → 包内放 `py.typed` 标记文件并随 wheel 分发 | PEP 561（官方原文 "MUST add a marker file named py.typed"） | High | null | 单文件模块不支持（须先成包，PEP 原文）；构建配置需包含该文件（setuptools 用 package_data，或 `[tool.setuptools.package-data]`） |
| `spdx_license_expression` | 3.9 | `license = {text = "MIT"}` 旧表格 / classifier 声明许可 → `license = "MIT"`（SPDX 表达式）+ `license-files = ["LICEN[CS]E*"]` | packaging.python.org/en/latest/guides/writing-pyproject-toml/（官方："As per PEP 639…this format is now deprecated"，附各后端最低版本表） | Medium | null | 旧 dict 格式弃用（PEP 639）；setuptools 需 ≥77.0.3、hatchling ≥1.27 等方可读新格式 |
| `uv_project_manager` | 3.9 | pip+venv+pip-tools+pipx 多工具手工流 → `uv`（单一 Rust 工具：项目/锁/解释器管理） | docs.astral.sh/uv/（官方自述："an extremely fast Python package and project manager"，可替代 pip/pip-tools/pipx/poetry/pyenv/twine/virtualenv） | High | null | 与本工具无关的独立建议；uv 管理 `.python-version` 与 `requires-python` 语义与 PRD 检测器兼容 |
| `universal_lockfile` | 3.9 | `requirements.txt` 手动 pin + 多平台矩阵 → `uv.lock` 跨平台通用锁文件（随 `uv sync`/`uv run` 消费） | docs.astral.sh/uv/（官方：universal lockfile + 平台无关解析 + 全局缓存） | Medium | null | 锁文件进版本库；CI 缓存收益大；与 pip-tools 的 hash pin 流程并存时须择一 |
| `ruff_lint_format` | 3.9 | flake8+isort+black+pyupgrade 多工具链 → ruff 单工具（`[tool.ruff]` 入 pyproject，`ruff check --fix`/`ruff format`） | docs.astral.sh/ruff/rules/（官方：Ruff 启用 F/E 子集，preview 默认集含 UP；UP 规则即本表 autofix 列来源） | High | null | 本数据集 autofix 字段即 `ruff:UPxxx`；UP006/UP007/UP045 在运行时注解库场景的 unsafe 警示见 §1/§2 对应行 |
| `pytest_config_in_pyproject` | 3.9 | `pytest.ini`/`setup.cfg [tool:pytest]` → `pyproject.toml [tool.pytest.ini_options]`（pytest ≥6.0；9.0 起原生 `[tool.pytest]`） | docs.pytest.org/en/stable/customize.html（官方优先级表：pytest.toml > pytest.ini > **pyproject.toml** > tox.ini > setup.cfg；"first match wins"，不合并） | Medium | null | pytest 6.0 起支持 pyproject 配置、9.0 起原生 TOML 类型 `[tool.pytest]`（官方两处版本注记） |

## 9. 跨版本主题与风险汇总

1. **pydantic / 运行时求值注解（最高优先 caveat）**：`X | Y`（PEP 604）依赖 3.10 的 `types.UnionType`，`list[int]`（PEP 585）依赖 3.9 的内建泛型。ruff 官方在 UP006、UP007、UP045 三页均明示：与"依赖运行时类型注解的库（如 pydantic）"共存时，低版本目标上 fix 标记 unsafe（`lint.pyupgrade.keep-runtime-typing = true` 可关停）。**数据集落库时：`union_type_annotation` 与 `builtin_generic_annotation` 的 details 必须携带此警示**（PRD 验收项）。
2. **PEP 686 翻转（3.15）**：3.10 的 `explicit_encoding_open`（"总是显式传 encoding"）在 3.15 UTF-8 默认化后失效。两条规则的 details 需互相引用：min<3.15 与 min≥3.15 的项目应得到相反建议（3.15 规则当前为 rc 占位，转正时落实）。
3. **行为等价陷阱清单**（非纯语法糖，autofix 不可盲改）：`enum_strenum`（f-string 展示变化）、`dataclass_slots`（类重建）、`typing_union_repr_unified`（repr 变化）、`forkserver_start_method`（默认值切换）、`type_alias_statement`（isinstance 语义）、`zip_strict`/`map_strict`/`batched_strict`（截断→抛错）。
4. **ruff UP038（isinstance 联合类型）已被 ruff 官方移除（since 0.13.0）**——`isinstance_union` 条目 autofix=null 并在 caveat 中说明，防止未来填充者误引。
5. **3.15 全节为 rc2+draft 占位**：PEP 编号（810/814/661/798/686/747/728/800）逐字取自 whatsnew/3.15 页面原文，但 final 前内容仍可能变动（PRD Open Questions 既定待办：final 后复核转正）。

## 10. 核实清点

- **UNVERIFIED 条目**：0。全部 88 条的引用均完成在线核实。
- **次要注意点（非未决）**：① `lru_cache_no_parens` 的 since_version 为 3.8（ruff UP011 官方页原文"since Python 3.8"），低于数据集地板 3.9，落库决策：并入 3.9 组或删；② whatsnew/3.14 的 `forkserver_start_method` 与 whatsnew/3.9 的 `functools_cache`、3.10 的 `zip_strict`/`int_bit_count` 引用为"PEP 号或节名"级（无逐字锚点），因抓取时官方页面该内容位于节内而非独立小节。
- **ruff 规则号核实方式**：19 个规则在 docs.astral.sh/ruff/rules/ 独立页面逐一核实（UP004/006/007/008/009/011/024/031/032/033/035/036/037/040/041/042/043/045/046/047）；UP017/UP019 仅经 `astral-sh/ruff` main 分支 `codes.rs` 官方映射核实（UP017=DatetimeTimezoneUTC、UP019=TypingTextStrAlias），文档页未单独抓取——如需 100% 页面级证据，落库前补抓两页即可。
- **抓取日期**：2026-09-29；whatsnew 3.9–3.14 抓自 docs.python.org/3/（3.14.7 文档集），3.15 抓自 docs.python.org/3.15/（3.15.0rc2）。

## 11. 分布统计（对照 PRD 验收）

| 版本段 | 条数 | 带 autofix 的行 | 其中 Critical/High |
| --- | --- | --- | --- |
| 3.9 | 11 | 4（UP006、UP019、UP033、UP011） | 3（Critical 2 + High 1） |
| 3.10 | 13 | 1（UP007/UP045 同行） | 5（Critical 1 + High 4） |
| 3.11 | 14 | 3（UP042、UP041、UP017） | 2（High 2） |
| 3.12 | 10 | 4（UP040、UP046、UP047、UP035） | 1（High 1） |
| 3.13 | 11 | 1（UP043） | 0 |
| 3.14 | 11 | 1（UP037 部分） | 1（High 1） |
| 3.15 占位 | 8 | 0 | 1（High 1） |
| Tooling | 10 | 0 | 7（Critical 2 + High 5） |
| **合计** | **88** | **14 行 / 15 个不同 UP 规则** | **20** |

（UP035 另在 3.9 的 `builtin_generic_annotation` caveat 中作为伴随清理提及，不重复计数。）

每版本 ≥3 条验收：3.9=11、3.10=13、3.11=14、3.12=10、3.13=11、3.14=11，全部满足。
