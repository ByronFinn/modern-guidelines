<!-- Generated from internal/guidelines/data/python.json; DO NOT EDIT. -->

# Modern Python Guidelines Explained

This file provides a more detailed description of the features supported in the Modern Python Guidelines.

**Autofix legend:**

- [x] — a mechanical autofix exists; the guideline below names its `tool:rule`
- [ ] — no autofix available

**Impact legend:**

- Critical — found in almost every project, dozens of occurrences
- High — found often, 5–20 occurrences per project
- Medium — found regularly, 1–5 occurrences per project
- Low — found rarely or in specific code

## Guidelines

| Category | Guideline | Autofix | Python | Impact |
|----------|-----------|------------|----|--------|
| IO | [`utf8_default_encoding`](#utf8_default_encoding) | [ ] | 3.15 | High |
| Builtins | [`sentinel_builtin`](#sentinel_builtin) | [ ] | 3.15 | Medium |
| Builtins | [`frozendict_builtin`](#frozendict_builtin) | [ ] | 3.15 | Medium |
| Imports | [`lazy_import_keyword`](#lazy_import_keyword) | [ ] | 3.15 | Medium |
| Syntax | [`comprehension_unpacking`](#comprehension_unpacking) | [ ] | 3.15 | Medium |
| Async | [`taskgroup_cancel_api`](#taskgroup_cancel_api) | [ ] | 3.15 | Low |
| Typing | [`typing_type_form`](#typing_type_form) | [ ] | 3.15 | Low |
| Typing | [`closed_typeddict`](#closed_typeddict) | [ ] | 3.15 | Low |
| Typing | [`deferred_annotations`](#deferred_annotations) | [x] | 3.14 | High |
| Collections | [`map_strict`](#map_strict) | [ ] | 3.14 | Medium |
| Files | [`pathlib_copy_move`](#pathlib_copy_move) | [ ] | 3.14 | Medium |
| Strings | [`t_string_templates`](#t_string_templates) | [ ] | 3.14 | Medium |
| Concurrency | [`forkserver_start_method`](#forkserver_start_method) | [ ] | 3.14 | Medium |
| Typing | [`typing_union_repr_unified`](#typing_union_repr_unified) | [ ] | 3.14 | Low |
| Errors | [`except_multiple_no_brackets`](#except_multiple_no_brackets) | [ ] | 3.14 | Low |
| Concurrency | [`subinterpreters_stdlib`](#subinterpreters_stdlib) | [ ] | 3.14 | Low |
| Data | [`pickle_protocol5_default`](#pickle_protocol5_default) | [ ] | 3.14 | Low |
| Files | [`zstandard_stdlib`](#zstandard_stdlib) | [ ] | 3.14 | Low |
| Runtime | [`free_threaded_supported`](#free_threaded_supported) | [ ] | 3.14 | Low |
| Typing | [`typing_type_is`](#typing_type_is) | [ ] | 3.13 | Medium |
| Typing | [`typeddict_readonly`](#typeddict_readonly) | [ ] | 3.13 | Medium |
| Typing | [`typevar_defaults`](#typevar_defaults) | [ ] | 3.13 | Medium |
| Typing | [`generator_default_params`](#generator_default_params) | [x] | 3.13 | Medium |
| Data | [`copy_replace`](#copy_replace) | [ ] | 3.13 | Medium |
| Deprecation | [`warnings_deprecated_decorator`](#warnings_deprecated_decorator) | [ ] | 3.13 | Medium |
| Deprecation | [`dead_batteries_replacement`](#dead_batteries_replacement) | [ ] | 3.13 | Medium |
| Typing | [`anystr_type_params`](#anystr_type_params) | [ ] | 3.13 | Low |
| Collections | [`batched_strict`](#batched_strict) | [ ] | 3.13 | Low |
| Runtime | [`process_cpu_count`](#process_cpu_count) | [ ] | 3.13 | Low |
| Runtime | [`locals_snapshot_semantics`](#locals_snapshot_semantics) | [ ] | 3.13 | Low |
| Typing | [`type_alias_statement`](#type_alias_statement) | [x] | 3.12 | High |
| Typing | [`generic_class_params`](#generic_class_params) | [x] | 3.12 | Medium |
| Typing | [`generic_function_params`](#generic_function_params) | [x] | 3.12 | Medium |
| Typing | [`typeddict_kwargs_unpack`](#typeddict_kwargs_unpack) | [ ] | 3.12 | Medium |
| Typing | [`override_decorator`](#override_decorator) | [ ] | 3.12 | Medium |
| Collections | [`itertools_batched`](#itertools_batched) | [ ] | 3.12 | Medium |
| Files | [`pathlib_walk`](#pathlib_walk) | [ ] | 3.12 | Medium |
| Deprecation | [`distutils_removed`](#distutils_removed) | [ ] | 3.12 | Medium |
| Typing | [`typing_abc_direct_import`](#typing_abc_direct_import) | [x] | 3.12 | Low |
| Strings | [`fstring_quote_reuse`](#fstring_quote_reuse) | [ ] | 3.12 | Low |
| Files | [`tomllib_stdlib`](#tomllib_stdlib) | [ ] | 3.11 | High |
| Async | [`asyncio_task_group`](#asyncio_task_group) | [ ] | 3.11 | High |
| Typing | [`typing_self`](#typing_self) | [ ] | 3.11 | Medium |
| Typing | [`typeddict_required_not_required`](#typeddict_required_not_required) | [ ] | 3.11 | Medium |
| Typing | [`typing_assert_never`](#typing_assert_never) | [ ] | 3.11 | Medium |
| Errors | [`exception_group_except_star`](#exception_group_except_star) | [ ] | 3.11 | Medium |
| Errors | [`exception_add_note`](#exception_add_note) | [ ] | 3.11 | Medium |
| Async | [`asyncio_timeout_cm`](#asyncio_timeout_cm) | [ ] | 3.11 | Medium |
| Errors | [`timeout_error_alias`](#timeout_error_alias) | [x] | 3.11 | Medium |
| Enums | [`enum_strenum`](#enum_strenum) | [x] | 3.11 | Medium |
| Datetime | [`datetime_fromisoformat_iso8601`](#datetime_fromisoformat_iso8601) | [ ] | 3.11 | Medium |
| Typing | [`typevar_tuple_variadic`](#typevar_tuple_variadic) | [ ] | 3.11 | Low |
| Async | [`asyncio_runner`](#asyncio_runner) | [ ] | 3.11 | Low |
| Datetime | [`datetime_utc_alias`](#datetime_utc_alias) | [x] | 3.11 | Low |
| Typing | [`union_type_annotation`](#union_type_annotation) | [x] | 3.10 | Critical |
| Data | [`dataclass_slots`](#dataclass_slots) | [ ] | 3.10 | High |
| Patterns | [`structural_pattern_matching`](#structural_pattern_matching) | [ ] | 3.10 | High |
| Collections | [`zip_strict`](#zip_strict) | [ ] | 3.10 | High |
| Files | [`explicit_encoding_open`](#explicit_encoding_open) | [ ] | 3.10 | High |
| Typing | [`isinstance_union`](#isinstance_union) | [ ] | 3.10 | Medium |
| Typing | [`type_alias_explicit`](#type_alias_explicit) | [ ] | 3.10 | Medium |
| Data | [`dataclass_kw_only`](#dataclass_kw_only) | [ ] | 3.10 | Medium |
| Syntax | [`parenthesized_context_managers`](#parenthesized_context_managers) | [ ] | 3.10 | Medium |
| Collections | [`itertools_pairwise`](#itertools_pairwise) | [ ] | 3.10 | Medium |
| Typing | [`typeguard_narrowing`](#typeguard_narrowing) | [ ] | 3.10 | Low |
| Typing | [`paramspec_callable`](#paramspec_callable) | [ ] | 3.10 | Low |
| Integers | [`int_bit_count`](#int_bit_count) | [ ] | 3.10 | Low |
| Typing | [`builtin_generic_annotation`](#builtin_generic_annotation) | [x] | 3.9 | Critical |
| Collections | [`dict_merge_operator`](#dict_merge_operator) | [ ] | 3.9 | Critical |
| Tooling | [`pyproject_single_config`](#pyproject_single_config) | [ ] | 3.9 | Critical |
| Tooling | [`no_setup_py_cli`](#no_setup_py_cli) | [ ] | 3.9 | Critical |
| Strings | [`str_removeprefix_removesuffix`](#str_removeprefix_removesuffix) | [ ] | 3.9 | High |
| Tooling | [`requires_python_metadata`](#requires_python_metadata) | [ ] | 3.9 | High |
| Tooling | [`src_layout`](#src_layout) | [ ] | 3.9 | High |
| Tooling | [`py_typed_marker`](#py_typed_marker) | [ ] | 3.9 | High |
| Tooling | [`uv_project_manager`](#uv_project_manager) | [ ] | 3.9 | High |
| Tooling | [`ruff_lint_format`](#ruff_lint_format) | [ ] | 3.9 | High |
| Typing | [`annotated_metadata`](#annotated_metadata) | [ ] | 3.9 | Medium |
| Typing | [`typing_text_str`](#typing_text_str) | [x] | 3.9 | Medium |
| Datetime | [`zoneinfo_stdlib`](#zoneinfo_stdlib) | [ ] | 3.9 | Medium |
| Async | [`asyncio_to_thread`](#asyncio_to_thread) | [ ] | 3.9 | Medium |
| Tooling | [`spdx_license_expression`](#spdx_license_expression) | [ ] | 3.9 | Medium |
| Tooling | [`universal_lockfile`](#universal_lockfile) | [ ] | 3.9 | Medium |
| Tooling | [`pytest_config_in_pyproject`](#pytest_config_in_pyproject) | [ ] | 3.9 | Medium |
| Data | [`functools_cache`](#functools_cache) | [x] | 3.9 | Medium |
| Data | [`lru_cache_no_parens`](#lru_cache_no_parens) | [x] | 3.9 | Medium |
| Algorithms | [`graphlib_topological_sort`](#graphlib_topological_sort) | [ ] | 3.9 | Low |
| Integers | [`math_lcm_multiarg`](#math_lcm_multiarg) | [ ] | 3.9 | Low |

---

<a id="utf8_default_encoding"></a>

## `utf8_default_encoding`

**Category: IO · Python 3.15+ · Impact: High · Autofix: none**

On 3.15+, rely on UTF-8 as the default text encoding and pass `encoding` only for non-UTF-8 data or legacy locale semantics.

PEP 686 flips the default: I/O opened without an explicit encoding now uses UTF-8, so blanket `encoding="utf-8"` becomes redundant and the 3.10 `explicit_encoding_open` advice reverses for projects floored at 3.15 or newer; keep passing it explicitly below 3.15. Non-UTF-8 data still needs an explicit encoding; use `encoding="locale"` or `PYTHONUTF8=0` / `-X utf8=0` for legacy behavior. Windows is most affected. Written from the 3.15 rc2 draft whatsnew; re-verify after final.

### Example

**Before:**

```python
def read_config(path):
    # defensive utf-8 was the rule before 3.15
    with open(path, encoding="utf-8") as fh:
        return fh.read()
```

**After:**

```python
def read_config(path):
    # UTF-8 is the default since 3.15 (PEP 686)
    with open(path) as fh:
        return fh.read()


def read_legacy(path):
    # pre-3.15 locale data still needs an explicit encoding
    with open(path, encoding="locale") as fh:
        return fh.read()
```

---

<a id="sentinel_builtin"></a>

## `sentinel_builtin`

**Category: Builtins · Python 3.15+ · Impact: Medium · Autofix: none**

Use the built-in `sentinel` type instead of module-level `object()` sentinels.

PEP 661 ships a pickleable built-in sentinel with a stable repr and stable identity, replacing hand-rolled `_MISSING = object()` markers that survive no pickle round-trip and print as bare object addresses. It composes with `|` so signatures can read `sentinel | None`. Written from the 3.15 rc2 draft whatsnew; re-verify the exact API after the final release.

### Example

**Before:**

```python
_MISSING = object()


def get(key, default=_MISSING):
    if default is _MISSING:
        raise KeyError(key)  # unpicklable default, opaque repr
```

**After:**

```python
def get(key, default=sentinel):
    if default is sentinel:
        raise KeyError(key)
```

---

<a id="frozendict_builtin"></a>

## `frozendict_builtin`

**Category: Builtins · Python 3.15+ · Impact: Medium · Autofix: none**

Store immutable mapping data in the built-in `frozendict` instead of `types.MappingProxyType` or third-party frozen maps.

PEP 814 adds an immutable, hashable, insertion-order-preserving built-in mapping. Unlike `MappingProxyType`, which is only a view over an underlying dict that can still mutate, `frozendict` owns its storage, so it is safe to share, cache, and use as a dict key or set element. Written from the 3.15 rc2 draft whatsnew; re-verify after the final release.

### Example

**Before:**

```python
from types import MappingProxyType

DEFAULTS = MappingProxyType(base_defaults)  # still mutable via base_defaults
```

**After:**

```python
DEFAULTS = frozendict(base_defaults)  # immutable, hashable, order-preserving
```

---

<a id="lazy_import_keyword"></a>

## `lazy_import_keyword`

**Category: Imports · Python 3.15+ · Impact: Medium · Autofix: none**

Defer expensive optional imports with the `lazy import` soft keyword instead of function-local imports.

PEP 810's `lazy import` loads the module on first attribute access, so heavyweight optional dependencies stop taxing startup while the import stays at module top level where linters and type checkers can see it. Load errors surface at first use rather than at import time. Written from the 3.15 rc2 draft whatsnew; re-verify against the final release.

### Example

**Before:**

```python
def render(chart):
    import matplotlib.pyplot as plt  # hidden from linters, hit every call
    ...
```

**After:**

```python
lazy import matplotlib.pyplot as plt  # loaded on first use


def render(chart):
    ...
```

---

<a id="comprehension_unpacking"></a>

## `comprehension_unpacking`

**Category: Syntax · Python 3.15+ · Impact: Medium · Autofix: none**

Flatten and merge inside comprehensions with unpacking (`[*xs for xs in groups]`, `{**d for d in dicts}`) instead of hand-written accumulation loops.

PEP 798 permits unpacking in comprehensions and generator expressions: `[*xs for xs in groups]` flattens nested iterables and `{**d for d in dicts}` merges mappings, replacing manual accumulation loops or `itertools.chain.from_iterable` in the eager cases. Generator expressions take the same syntax. Written from the 3.15 rc2 draft whatsnew; re-verify after the final release.

### Example

**Before:**

```python
from itertools import chain

flat = list(chain.from_iterable(chunks))
merged = {}
for part in configs:
    merged.update(part)
```

**After:**

```python
flat = [*chunk for chunk in chunks]
merged = {**part for part in configs}
```

---

<a id="taskgroup_cancel_api"></a>

## `taskgroup_cancel_api`

**Category: Async · Python 3.15+ · Impact: Low · Autofix: none**

Stop a running `asyncio.TaskGroup` with its `cancel()` method instead of cancelling inner tasks or signalling through side channels.

TaskGroups previously had no first-party way for surrounding code to stop the group early; 3.15 adds `TaskGroup.cancel()` for early task-group termination, propagating cancellation to every child task from one call. Prefer it over cancelling inner tasks one by one or routing through a shared event flag. Written from the 3.15 rc2 draft whatsnew; re-verify after the final release.

### Example

**Before:**

```python
stop = asyncio.Event()

async with asyncio.TaskGroup() as tg:
    tg.create_task(poll(stop))
    tg.create_task(refresh(stop))

# stopping the group meant flipping the event and
# cancelling inner tasks one by one
```

**After:**

```python
async with asyncio.TaskGroup() as tg:
    tg.create_task(poll())
    tg.create_task(refresh())
    if should_stop:
        tg.cancel()  # early termination of every child task
```

---

<a id="typing_type_form"></a>

## `typing_type_form`

**Category: Typing · Python 3.15+ · Impact: Low · Autofix: none**

Annotate class-or-factory parameters with `typing.TypeForm` instead of `type[T]`.

`type[T]` accepts only what `type(...)` produces, silently excluding metaclass-made classes and factory callables; PEP 747's `TypeForm[T]` describes the full class-like argument space and matches what checkers infer from `x = SomeClass`. It is a static-checking construct with no runtime effect, aligned with the earlier typing_extensions experiment. Written from the 3.15 rc2 draft whatsnew; re-verify after final.

### Example

**Before:**

```python
from typing import TypeVar

T = TypeVar("T", bound=Plugin)


def register(cls: type[T]) -> None: ...  # rejects factories, metaclass classes
```

**After:**

```python
from typing import TypeForm


def register(cls: TypeForm[Plugin]) -> None: ...
```

---

<a id="closed_typeddict"></a>

## `closed_typeddict`

**Category: Typing · Python 3.15+ · Impact: Low · Autofix: none**

Reject unknown keys statically with `class Config(TypedDict, closed=True)` instead of hand-rolled `__extra_items__` conventions.

PEP 728 adds closed TypedDicts: `closed=True` makes type checkers flag extra keys, and `extra_items` types whatever unknown keys a caller may still pass, replacing ad-hoc validation docstrings and custom `__extra_items__` markers. It binds static checkers only; runtime dicts keep accepting any keys. Written from the 3.15 rc2 draft whatsnew; re-verify after the final release.

### Example

**Before:**

```python
from typing import TypedDict


class Config(TypedDict):
    retries: int
    # unknown keys pass static checks unchecked
```

**After:**

```python
from typing import TypedDict


class Config(TypedDict, closed=True):
    retries: int
```

---

<a id="deferred_annotations"></a>

## `deferred_annotations`

**Category: Typing · Python 3.14+ · Impact: High · Autofix: ruff:UP037**

Write annotations as bare names instead of quoted forward references; PEP 649 lazy evaluation makes `from __future__ import annotations` unnecessary.

Annotations are evaluated lazily by default since 3.14, so quoted forward references like `-> "Node"` and a blanket `from __future__ import annotations` become unnecessary. `__annotations__` access semantics changed — code reading annotations directly should move to `annotationlib`; the future import stays legal but no longer needed. The autofix only strips quoted annotations.

### Example

**Before:**

```python
from __future__ import annotations

class Node:
    def children(self) -> list["Node"]: ...
```

**After:**

```python
class Node:
    def children(self) -> list[Node]: ...
```

---

<a id="map_strict"></a>

## `map_strict`

**Category: Collections · Python 3.14+ · Impact: Medium · Autofix: none**

Pass `strict=True` to `map()` when mapping over multiple iterables that must have equal lengths.

`map()` still stops at the shortest input by default; `map(f, a, b, strict=True)` (3.14) raises `ValueError` when one iterable is exhausted before the others, matching `zip(strict=True)`. This changes the failure mode of code that relied on silent truncation, so audit multi-iterable `map()` calls before flipping the flag; the ruff check is B912, not a UP autofix.

### Example

**Before:**

```python
totals = list(map(operator.add, revenues, costs))
# a longer revenues list is silently truncated
```

**After:**

```python
totals = list(map(operator.add, revenues, costs, strict=True))
# raises ValueError on length mismatch
```

---

<a id="pathlib_copy_move"></a>

## `pathlib_copy_move`

**Category: Files · Python 3.14+ · Impact: Medium · Autofix: none**

Copy and move files or directory trees with `Path.copy()`, `Path.copy_into()`, `Path.move()`, and `Path.move_into()` instead of mixing `shutil` calls into `pathlib` code.

3.14 adds tree-aware copy/move methods to `Path` (3.13 has none of them); each returns the destination `Path`. `copy(target, *, follow_symlinks=True, preserve_metadata=False)` replaces an existing target file and skips metadata by default; `move(target)` overwrites a file target but raises `OSError` for a non-empty directory target and falls back to copy-plus-delete across filesystems. These methods do not exist before 3.14.

### Example

**Before:**

```python
import shutil

shutil.copytree("dist/site", "/srv/www/site")
shutil.move("dist/site.html", "/srv/www/site.html")
```

**After:**

```python
Path("dist/site").copy("/srv/www/site")
Path("dist/site.html").move("/srv/www/site.html")
```

---

<a id="t_string_templates"></a>

## `t_string_templates`

**Category: Strings · Python 3.14+ · Impact: Medium · Autofix: none**

Build structured, injection-safe interpolations with `t"..."` template strings instead of f-string concatenation for SQL, HTML, or shell fragments.

PEP 750 t-strings evaluate to `string.templatelib.Template` — an object holding static text plus `Interpolation` parts, not a `str`. A renderer supplied by the database driver, HTML library, or shell builder decides how each interpolation is quoted, moving safety from discipline into the library. F-strings remain correct when the finished string is wanted immediately; passing a Template where a str is expected fails loudly, which is the point.

### Example

**Before:**

```python
query = f"SELECT * FROM users WHERE name = '{name}'"
rows = db.execute(query)  # quoting is manual and fragile
```

**After:**

```python
query = t"SELECT * FROM users WHERE name = {name}"
rows = db.execute(query)  # the driver renders interpolations safely
```

---

<a id="forkserver_start_method"></a>

## `forkserver_start_method`

**Category: Concurrency · Python 3.14+ · Impact: Medium · Autofix: none**

Expect `forkserver` as the Unix default start method in 3.14, and request `fork` explicitly — only — where its semantics are required.

On Unix platforms other than macOS, 3.14 replaces the `fork` default with `forkserver` in both multiprocessing and concurrent.futures, closing the fork-with-threads hazard: workers are forked from a clean server process. Costs: no copy-on-write sharing of parent memory, arguments must pickle, and startup is slower. Code relying on fork's memory sharing or unpicklable globals must opt back in via `get_context("fork")` or `ProcessPoolExecutor(mp_context=...)`.

### Example

**Before:**

```python
# 3.13 and earlier on Unix: default start method "fork"
with multiprocessing.Pool(4) as pool:  # children share memory CoW
    results = pool.map(heavy, chunks)
```

**After:**

```python
# 3.14: Unix default is "forkserver". Keep fork only if required:
ctx = multiprocessing.get_context("fork")
with ctx.Pool(4) as pool:
    results = pool.map(heavy, chunks)
```

---

<a id="typing_union_repr_unified"></a>

## `typing_union_repr_unified`

**Category: Typing · Python 3.14+ · Impact: Low · Autofix: none**

Do not depend on `repr()` or caching differences between `typing.Union` and `|` unions — 3.14 unifies them.

3.14 completes PEP 604 convergence: `repr(typing.Union[int, str])` now renders exactly `int | str`, and the runtime no longer caches legacy `Union` instances. This is an observable behavior change — snapshot tests or code that serializes type reprs will see different output, per the official porting notes.

### Example

**Before:**

```python
assert repr(Union[int, str]) == "typing.Union[int, str]"  # breaks on 3.14
```

**After:**

```python
assert repr(Union[int, str]) == "int | str"  # unified spelling
```

---

<a id="except_multiple_no_brackets"></a>

## `except_multiple_no_brackets`

**Category: Errors · Python 3.14+ · Impact: Low · Autofix: none**

List multiple exception types in `except A, B:` without parentheses when the clause binds no name.

PEP 758 drops the mandatory tuple for both `except` and `except*` in 3.14; `except* A, B:` was likewise a SyntaxError before. Parentheses remain required whenever an `as` clause follows, so the grammar cannot misread the bound name as an exception type. The bracketed tuple form stays fully valid; treat the short form as style, not a forced migration.

### Example

**Before:**

```python
try:
    process(row)
except (ValueError, KeyError):
    log.warning("skipped malformed row")
```

**After:**

```python
try:
    process(row)
except ValueError, KeyError:  # no `as` clause: brackets optional
    log.warning("skipped malformed row")
```

---

<a id="subinterpreters_stdlib"></a>

## `subinterpreters_stdlib`

**Category: Concurrency · Python 3.14+ · Impact: Low · Autofix: none**

Fan out isolated Python work with `InterpreterPoolExecutor` / `concurrent.interpreters` instead of multiprocessing when startup and IPC dominate (PEP 734).

PEP 734 ships per-interpreter GILs: `concurrent.futures.InterpreterPoolExecutor` runs calls on a pool of subinterpreters inside one process, so there is no process-spawn cost per worker and each worker holds its own GIL. Interpreters share nothing — arguments and results still cross by pickling or shared buffers — and C extensions must be per-interpreter compatible before you switch.

### Example

**Before:**

```python
from multiprocessing import Pool

with Pool(4) as pool:  # process spawn + pickle per task
    results = pool.map(heavy, chunks)
```

**After:**

```python
from concurrent.futures import InterpreterPoolExecutor

with InterpreterPoolExecutor() as pool:  # one GIL per worker
    results = list(pool.map(heavy, chunks))
```

---

<a id="pickle_protocol5_default"></a>

## `pickle_protocol5_default`

**Category: Data · Python 3.14+ · Impact: Low · Autofix: none**

Drop the explicit `protocol=5` from `pickle.dump`/`dumps` on 3.14, where protocol 5 became the default.

3.14 raises the default pickle protocol from 4 to 5, so the explicit argument is redundant and out-of-band buffers (PEP 574) no longer need a protocol opt-in. Only the default moved: data written with protocol 5 still cannot be read by interpreters older than 3.8, so payloads consumed by 3.7-era runtimes must pin `protocol=4` explicitly rather than trust the writer's default.

### Example

**Before:**

```python
blob = pickle.dumps(obj, protocol=5)  # pinned since 3.8 days
```

**After:**

```python
blob = pickle.dumps(obj)  # default protocol is 5 since 3.14
# pin protocol=4 explicitly if 3.7 runtimes must read it
```

---

<a id="zstandard_stdlib"></a>

## `zstandard_stdlib`

**Category: Files · Python 3.14+ · Impact: Low · Autofix: none**

Compress with the `compression.zstd` module (PEP 784) instead of the third-party `zstandard` package; `tarfile`, `zipfile`, and `shutil` accept zstd too.

PEP 784 adds `compression.zstd` with a `zlib`-like surface — `zstd.compress(data)` and `zstd.decompress(blob)` — and threads the codec through `tarfile`, `zipfile`, and `shutil`, letting archive code drop a dependency. It is not a drop-in for the `zstandard` package: names and streaming idioms differ, so port call sites deliberately, and floors below 3.14 still need the third-party package as a conditional dependency.

### Example

**Before:**

```python
import zstandard

blob = zstandard.ZstdCompressor().compress(data)
```

**After:**

```python
from compression import zstd

blob = zstd.compress(data)
```

---

<a id="free_threaded_supported"></a>

## `free_threaded_supported`

**Category: Runtime · Python 3.14+ · Impact: Low · Autofix: none**

Treat the free-threaded build (`python3.14t`) as a supported deployment target since 3.14 (PEP 779), after auditing extensions and shared state.

PEP 779 moves free-threaded CPython from experimental to officially supported: the `t` build ships alongside the GIL build and is a legitimate production target. The official whatsnew puts single-threaded overhead at roughly five to ten percent. The audit is the real work — C extensions need free-threaded builds, and code that leaned on the GIL for atomicity needs explicit locking before switching.

### Example

**Before:**

```python
# CI matrix
python-version: ["3.12", "3.13", "3.14"]  # GIL builds only
```

**After:**

```python
# CI matrix
python-version: ["3.14", "3.14t"]  # 3.14t: free-threaded (PEP 779)
```

---

<a id="typing_type_is"></a>

## `typing_type_is`

**Category: Typing · Python 3.13+ · Impact: Medium · Autofix: none**

Return `TypeIs[T]` from narrowing predicates instead of `TypeGuard[T]`.

`TypeIs` narrows like `isinstance`: a True result narrows to the checked type and the else branch narrows to the complement, fixing `TypeGuard`'s one-directional narrowing. It is a static-checking contract only — nothing is validated at runtime, so predicates must still be implemented correctly.

### Example

**Before:**

```python
def is_str_list(values: list[object]) -> TypeGuard[list[str]]: ...
```

**After:**

```python
def is_str_list(values: list[object]) -> TypeIs[list[str]]: ...
```

---

<a id="typeddict_readonly"></a>

## `typeddict_readonly`

**Category: Typing · Python 3.13+ · Impact: Medium · Autofix: none**

Mark write-once TypedDict keys with `ReadOnly[...]` instead of naming conventions or docstrings.

PEP 705 adds `typing.ReadOnly` so individual TypedDict keys are checked as settable only at construction; checkers then reject later assignment. Purely a static contract — no runtime immutability is enforced; covariance support is the PEP's main motivation.

### Example

**Before:**

```python
class Config(TypedDict):
    # "id must not be overwritten" — convention only
    id: str
    retries: int
```

**After:**

```python
class Config(TypedDict):
    id: ReadOnly[str]
    retries: int
```

---

<a id="typevar_defaults"></a>

## `typevar_defaults`

**Category: Typing · Python 3.13+ · Impact: Medium · Autofix: none**

Give TypeVars defaults (`TypeVar("T", default=...)`) instead of hand-written placeholder aliases.

PEP 696 lets a TypeVar fall back to a default type, removing placeholder TypeVar aliases and the unions they force onto every signature. `ParamSpec` and `TypeVarTuple` also accept defaults, but support is checker-dependent — verify your toolchain resolves the fallback before relying on it.

### Example

**Before:**

```python
T = TypeVar("T")
DefaultT = TypeVar("DefaultT")

def first(items: list[T], fallback: DefaultT) -> T | DefaultT: ...
```

**After:**

```python
T = TypeVar("T", default=None)

def first(items: list[T]) -> T | None: ...
```

---

<a id="generator_default_params"></a>

## `generator_default_params`

**Category: Typing · Python 3.13+ · Impact: Medium · Autofix: ruff:UP043**

Write `Generator[int]` instead of `Generator[int, None, None]` when the extra types are `None`.

Since 3.13 the send and return type parameters of `Generator` and `AsyncGenerator` default to `None`, so trailing `None` arguments are redundant noise; ruff UP043 removes them and the fix is safe. Keep the explicit arguments when those types genuinely differ from `None`.

### Example 1

**Before:**

```python
def poll() -> Generator[int, None, None]: ...
```

**After:**

```python
def poll() -> Generator[int]: ...
```

### Example 2

**Before:**

```python
async def watch() -> AsyncGenerator[Event, None]: ...
```

**After:**

```python
async def watch() -> AsyncGenerator[Event]: ...
```

---

<a id="copy_replace"></a>

## `copy_replace`

**Category: Data · Python 3.13+ · Impact: Medium · Autofix: none**

Use `copy.replace(obj, **changes)` for immutable updates instead of per-type `_replace` helpers and hand-written `replace()` methods.

3.13 introduces a `__replace__` protocol and `copy.replace()` as one uniform update operation; namedtuple, dataclass, datetime, ZoneInfo and other stdlib types already implement it. The target type must implement `__replace__` — custom immutable classes need to add it themselves before call sites can switch.

### Example

**Before:**

```python
new_dt = dt.replace(year=2031)   # datetime-specific method
new_rec = record._replace(x=1)   # namedtuple private API
```

**After:**

```python
import copy

new_dt = copy.replace(dt, year=2031)
new_rec = copy.replace(record, x=1)
```

---

<a id="warnings_deprecated_decorator"></a>

## `warnings_deprecated_decorator`

**Category: Deprecation · Python 3.13+ · Impact: Medium · Autofix: none**

Mark deprecated APIs with `@warnings.deprecated("...")` instead of hand-rolled DeprecationWarning decorators.

PEP 702 makes the signal uniform: one decorator on a function, class, method, or property emits `DeprecationWarning` at runtime on use and lets type checkers flag call sites statically, replacing bespoke decorator stacks that only covered functions. Runtime warnings fire on call for functions and on attribute access for classes; the static half needs checker support, so keep the changelog entry too.

### Example

**Before:**

```python
import warnings

def deprecated(msg):
    def wrap(fn):
        def inner(*a, **kw):
            warnings.warn(msg, DeprecationWarning, stacklevel=2)
            return fn(*a, **kw)
        return inner
    return wrap

@deprecated("use fetch_v2")
def fetch(url): ...
```

**After:**

```python
import warnings

@warnings.deprecated("use fetch_v2")
def fetch(url): ...
```

---

<a id="dead_batteries_replacement"></a>

## `dead_batteries_replacement`

**Category: Deprecation · Python 3.13+ · Impact: Medium · Autofix: none**

Replace the 19 modules removed by PEP 594 (`cgi`, `crypt`, `telnetlib`, `pipes`, `imghdr`, ...) with their documented successors before moving to 3.13.

3.13 deletes the PEP 594 dead batteries, and `import cgi` now raises `ModuleNotFoundError`. Each module has a mapped successor: password hashing moves from `crypt` to `hashlib` or passlib, `telnetlib` to asyncio streams or telnetlib3, `pipes` to `subprocess`, `imghdr` to filetype or puremagic. Several successors are third-party, so budget the new dependencies instead of assuming stdlib-for-stdlib swaps.

### Example

**Before:**

```python
import cgi
import crypt
import telnetlib  # all removed in 3.13 (PEP 594)
```

**After:**

```python
from hashlib import pbkdf2_hmac  # replaces crypt password hashing
from urllib.parse import parse_qs  # replaces cgi form helpers
# telnetlib -> asyncio streams; PEP 594 maps all 19 modules
```

---

<a id="anystr_type_params"></a>

## `anystr_type_params`

**Category: Typing · Python 3.13+ · Impact: Low · Autofix: none**

Replace `typing.AnyStr` with explicit type parameters bound to `str` or `bytes`.

`typing.AnyStr` is deprecated since 3.13 and will be removed in 3.18, joining `typing.ByteString` (removal 3.17) among top-level typing aliases with announced removal versions. Declare the constraint explicitly instead. The PEP 695 constraint syntax shown needs 3.12+; on older interpreter floors use a TypeVar constrained to `(str, bytes)`.

### Example

**Before:**

```python
def concat(a: AnyStr, b: AnyStr) -> AnyStr: ...
```

**After:**

```python
def concat[S: (str, bytes)](a: S, b: S) -> S: ...
```

---

<a id="batched_strict"></a>

## `batched_strict`

**Category: Collections · Python 3.13+ · Impact: Low · Autofix: none**

Pass `strict=True` to `itertools.batched()` when the input length must divide evenly by the batch size.

From 3.13, `batched(iterable, n, strict=True)` raises `ValueError` if the final batch would be shorter than `n`, extending the 3.10 `zip(strict=True)` guarantee to chunking loops. Use it only where uneven division is genuinely an error — streams of unknown length (sockets, generators) legitimately end mid-batch, and strict mode would turn a normal end-of-stream into a crash.

### Example

**Before:**

```python
for page in batched(rows, page_size):
    # a short final page is passed through silently
    render(page)
```

**After:**

```python
for page in batched(rows, page_size, strict=True):
    render(page)
```

---

<a id="process_cpu_count"></a>

## `process_cpu_count`

**Category: Runtime · Python 3.13+ · Impact: Low · Autofix: none**

Size parallelism with `os.process_cpu_count()` instead of `os.cpu_count()` when affinity or container limits matter.

`os.process_cpu_count` returns the CPUs actually usable by this process — after `sched_setaffinity`, and it is the value `PYTHON_CPU_COUNT` overrides — while `os.cpu_count` reports the machine total and can massively overstate a cgroup-limited container. Concurrency primitives account for it from 3.13, so explicit worker sizing should read the same number.

### Example

**Before:**

```python
workers = min(8, os.cpu_count() or 1)
# a 128-core host limited to 2 CPUs yields 8 workers
```

**After:**

```python
workers = min(8, os.process_cpu_count() or 1)
# respects affinity and container CPU limits
```

---

<a id="locals_snapshot_semantics"></a>

## `locals_snapshot_semantics`

**Category: Runtime · Python 3.13+ · Impact: Low · Autofix: none**

Assign variables directly instead of writing through `locals()`; PEP 667 made the snapshot semantics explicit in 3.13.

PEP 667 replaces old undefined behavior: `locals()` now formally returns a snapshot, while `frame.f_locals` is a write-through proxy shared across accesses. Debuggers, tracers, and REPLs gain introspection that actually sticks, but code doing `locals()["x"] = value` must switch to direct assignment — the write-through path belongs to tooling, and optimized module and class scopes never propagated writes anyway.

### Example

**Before:**

```python
def run():
    locals()["mode"] = "fast"  # undefined write-back
    dispatch()
```

**After:**

```python
def run():
    mode = "fast"
    dispatch()
# tooling that must mutate: use frame.f_locals (PEP 667 proxy)
```

---

<a id="type_alias_statement"></a>

## `type_alias_statement`

**Category: Typing · Python 3.12+ · Impact: High · Autofix: ruff:UP040**

Declare type aliases with the `type` statement instead of `TypeAlias` assignments.

The `type X = ...` statement creates a `TypeAliasType` with lazy evaluation and unambiguous alias semantics. Plain assignments like `X = int` can be used with `isinstance`, while the `type` statement product raises `TypeError` there; variance and bound information may be lost, and the fix is marked unsafe.

### Example

**Before:**

```python
UserId: TypeAlias = int
Vector: TypeAlias = list[float]
```

**After:**

```python
type UserId = int
type Vector = list[float]
```

---

<a id="generic_class_params"></a>

## `generic_class_params`

**Category: Typing · Python 3.12+ · Impact: Medium · Autofix: ruff:UP046**

Declare generic classes inline (`class Box[T]:`) instead of module-level TypeVars plus `Generic[T]`.

PEP 695 binds type parameters to the class declaration, dropping module-level TypeVars and the `Generic` base class. Variance is inferred and is not guaranteed to match an explicit `covariant`/`contravariant` declaration, and ruff's fix is unsafe because the old TypeVar assignment is not removed automatically.

### Example

**Before:**

```python
T = TypeVar("T")

class Box(Generic[T]):
    def __init__(self, item: T) -> None: ...
    def get(self) -> T: ...
```

**After:**

```python
class Box[T]:
    def __init__(self, item: T) -> None: ...
    def get(self) -> T: ...
```

---

<a id="generic_function_params"></a>

## `generic_function_params`

**Category: Typing · Python 3.12+ · Impact: Medium · Autofix: ruff:UP047**

Declare generic functions inline (`def first[T](...)`) instead of module-level TypeVars.

The same PEP 695 mechanism binds type parameters to the function: `def first[T](items: list[T]) -> T`. Inferred variance may differ from the old explicit TypeVar declaration, and ruff notes some type checkers' PEP 695 support is still maturing — confirm your checker handles it before bulk rewrites.

### Example

**Before:**

```python
T = TypeVar("T")

def first(items: list[T]) -> T: ...
```

**After:**

```python
def first[T](items: list[T]) -> T: ...
```

---

<a id="typeddict_kwargs_unpack"></a>

## `typeddict_kwargs_unpack`

**Category: Typing · Python 3.12+ · Impact: Medium · Autofix: none**

Type `**kwargs` with `**kwargs: Unpack[SomeTypedDict]` instead of `Any` or a plain dict.

PEP 692 types each keyword argument through a TypedDict, so `**kwargs: Unpack[ConnOpts]` gives callers per-key checking with `Required`/`NotRequired` control over optionality. Purely a type-checker contract — no runtime validation of kwargs happens; keep validating values as before.

### Example

**Before:**

```python
def connect(**kwargs: Any) -> Connection: ...
```

**After:**

```python
class ConnOpts(TypedDict):
    host: str
    timeout: NotRequired[float]

def connect(**kwargs: Unpack[ConnOpts]) -> Connection: ...
```

---

<a id="override_decorator"></a>

## `override_decorator`

**Category: Typing · Python 3.12+ · Impact: Medium · Autofix: none**

Mark overriding methods with `@typing.override` so misspellings fail static checks.

`@typing.override` asserts the decorated method exists on a base class, so a misspelled override like `get_colour` fails static analysis instead of silently becoming a new method. Enforced by type checkers only — the decorator does nothing at runtime.

### Example

**Before:**

```python
class TruncatingFilter(Filter):
    def aply(self, text: str) -> str: ...  # typo: never called
```

**After:**

```python
class TruncatingFilter(Filter):
    @override
    def apply(self, text: str) -> str: ...
```

---

<a id="itertools_batched"></a>

## `itertools_batched`

**Category: Collections · Python 3.12+ · Impact: Medium · Autofix: none**

Chunk iterables with `itertools.batched(it, n)` instead of hand-written `range(0, len(xs), n)` slicing loops.

`batched()` yields tuples of up to `n` items from any iterable — generators included — with no sequence requirement, no index arithmetic, and no intermediate list slices. The final batch may be shorter than `n`; if uneven division is an error, 3.13 adds `strict=True` to raise `ValueError` instead (the ruff check is B911, not a UP autofix).

### Example

**Before:**

```python
for i in range(0, len(records), 100):
    flush(records[i : i + 100])
```

**After:**

```python
from itertools import batched

for group in batched(records, 100):
    flush(group)
```

---

<a id="pathlib_walk"></a>

## `pathlib_walk`

**Category: Files · Python 3.12+ · Impact: Medium · Autofix: none**

Traverse directory trees with `Path.walk()` instead of re-joining `os.walk()` strings into `Path` objects.

`Path.walk(top_down=True, on_error=None, follow_symlinks=False)` yields `(dirpath, dirnames, filenames)` where `dirpath` is a `Path` and the name lists are strings — join with `dirpath / name`. It is not a rename of `os.walk()`: with `follow_symlinks=False` (the default), symlinks to directories land in `filenames`, whereas `os.walk()` categorizes them as directories; pruning mutates `dirnames` in place and only works when `top_down` is true.

### Example

**Before:**

```python
import os

for root, dirs, files in os.walk("src"):
    for name in files:
        path = Path(root) / name
        if path.suffix == ".py":
            yield path
```

**After:**

```python
for root, dirs, files in Path("src").walk():
    for name in files:
        path = root / name
        if path.suffix == ".py":
            yield path
```

---

<a id="distutils_removed"></a>

## `distutils_removed`

**Category: Deprecation · Python 3.12+ · Impact: Medium · Autofix: none**

Migrate off `distutils` before 3.12, where it was removed (PEP 632); build metadata belongs in pyproject.toml.

PEP 632 removed distutils in 3.12, so `import distutils` fails unless setuptools happens to be installed — and venvs created by 3.12+ no longer preinstall setuptools, so even the compatibility shim is not guaranteed. Move build configuration to `[build-system]` and `[project]`, file operations to `shutil`, and packaging helpers to dedicated libraries; the PEP maps every removed API to a replacement.

### Example

**Before:**

```python
from distutils.core import setup

setup(name="mypkg", version="1.0")
```

**After:**

```python
# pyproject.toml
[build-system]
requires = ["setuptools>=77"]
build-backend = "setuptools.build_meta"

[project]
name = "mypkg"
version = "1.0"
```

---

<a id="typing_abc_direct_import"></a>

## `typing_abc_direct_import`

**Category: Typing · Python 3.12+ · Impact: Low · Autofix: ruff:UP035**

Import ABCs like `Hashable` and `Sized` from `collections.abc`, not `typing`.

The `typing.Hashable`, `typing.Sized` and sibling re-exports are deprecated since 3.12; import the ABCs from `collections.abc` (ruff UP035 rewrites the imports). `typing.ByteString` is stricter — deprecated since 3.9 and removed in 3.17 — and should migrate to `collections.abc.Buffer` (PEP 688), not merely be re-imported.

### Example

**Before:**

```python
from typing import Hashable, Iterable, Sized
```

**After:**

```python
from collections.abc import Hashable, Iterable, Sized
```

---

<a id="fstring_quote_reuse"></a>

## `fstring_quote_reuse`

**Category: Strings · Python 3.12+ · Impact: Low · Autofix: none**

Reuse the same quote character and nest expressions freely inside f-strings (PEP 701).

PEP 701 formalized f-string parsing: expressions may contain strings with the enclosing quote, arbitrary nesting, multi-line spans, comments, backslashes, and `\N{...}` escapes — the old quote-alternation and temp-variable workarounds are obsolete. These are a `SyntaxError` before 3.12, and the new token stream confused tools built on the old tokenizer, so keep formatters and coverage tools 3.12-aware.

### Example

**Before:**

```python
playlist = "This is the playlist: {}".format(", ".join(songs))
```

**After:**

```python
playlist = f"This is the playlist: {", ".join(songs)}"
```

---

<a id="tomllib_stdlib"></a>

## `tomllib_stdlib`

**Category: Files · Python 3.11+ · Impact: High · Autofix: none**

Read TOML with the standard library `tomllib` instead of the third-party `tomli`.

PEP 680 adds a TOML 1.0 reader: `tomllib.load(fh)` expects a file opened in binary mode, `tomllib.loads(text)` takes a `str`; both return plain dicts and raise `TOMLDecodeError` on malformed input. The module is read-only — writing stays with `tomli-w`. Passing a text-mode file raises `TypeError`, and on floors below 3.11 keep the `try: import tomllib / except ModuleNotFoundError: import tomli as tomllib` shim.

### Example

**Before:**

```python
import tomli

with open("pyproject.toml", "rb") as fh:
    config = tomli.load(fh)
```

**After:**

```python
import tomllib

with open("pyproject.toml", "rb") as fh:
    config = tomllib.load(fh)
```

---

<a id="asyncio_task_group"></a>

## `asyncio_task_group`

**Category: Async · Python 3.11+ · Impact: High · Autofix: none**

Run concurrent coroutines under `async with asyncio.TaskGroup()` instead of `create_task` plus `gather`.

The group awaits every child on exit and, as soon as one fails, cancels the siblings and re-raises all failures as an `ExceptionGroup` (a `BaseExceptionGroup` when a leaf derives only from BaseException); handle it with `except*` per `exception_group_except_star`. Tasks can no longer be forgotten the way bare `create_task` allowed. Fetch results from the task objects you keep; there is no `gather`-style combined return.

### Example

**Before:**

```python
tasks = [asyncio.create_task(fetch(url)) for url in urls]
results = await asyncio.gather(*tasks)
# one crash leaves siblings running or errors unobserved
```

**After:**

```python
async with asyncio.TaskGroup() as tg:
    tasks = [tg.create_task(fetch(url)) for url in urls]
results = [t.result() for t in tasks]
```

---

<a id="typing_self"></a>

## `typing_self`

**Category: Typing · Python 3.11+ · Impact: Medium · Autofix: none**

Annotate instance-returning methods with `typing.Self` instead of class-name strings or custom TypeVars.

`typing.Self` refers to the enclosing class, replacing quoted class-name forward references and self-referential TypeVars in `__enter__`, `classmethod` constructors, and fluent builders; subclasses keep precise return types. It is mainly a static construct, though pydantic v2 does support `Self` at runtime.

### Example

**Before:**

```python
class Builder:
    def reset(self) -> "Builder": ...
```

**After:**

```python
class Builder:
    def reset(self) -> Self: ...
```

---

<a id="typeddict_required_not_required"></a>

## `typeddict_required_not_required`

**Category: Typing · Python 3.11+ · Impact: Medium · Autofix: none**

Control TypedDict key optionality per key with `Required[...]`/`NotRequired[...]` instead of class-level `total=`.

`total=` makes every key required or every key optional in one shot; PEP 655 items override it per key, so mixed required/optional shapes need no artificial subclass split. A static-checker contract only — keys are never validated at runtime.

### Example

**Before:**

```python
class Job(TypedDict, total=False):
    name: str  # actually required, but unchecked
    retries: int
```

**After:**

```python
class Job(TypedDict):
    name: Required[str]
    retries: NotRequired[int]
```

---

<a id="typing_assert_never"></a>

## `typing_assert_never`

**Category: Typing · Python 3.11+ · Impact: Medium · Autofix: none**

Close exhaustive `match`/`if` chains with `typing.assert_never(value)` instead of ad-hoc raises.

Passing the exhaustively narrowed value to `typing.assert_never` makes checkers prove every enum member or Literal case is handled; adding a new member becomes a visible error. At runtime it raises `AssertionError` — a developer tripwire, not a user-facing error type.

### Example

**Before:**

```python
def render(mode: Literal["dark", "light"]) -> str:
    if mode == "dark":
        return "#000"
    elif mode == "light":
        return "#fff"
    raise AssertionError("unreachable")
```

**After:**

```python
def render(mode: Literal["dark", "light"]) -> str:
    if mode == "dark":
        return "#000"
    elif mode == "light":
        return "#fff"
    assert_never(mode)
```

---

<a id="exception_group_except_star"></a>

## `exception_group_except_star`

**Category: Errors · Python 3.11+ · Impact: Medium · Autofix: none**

Propagate multiple failures as an `ExceptionGroup` and split them by type with `except*` instead of collapsing them into one error or swallowing them via `gather(return_exceptions=True)`.

`except*` runs once per matching subgroup and re-raises the non-matching remainder as a new group; the handler binds that subgroup, not the original leaves. It cannot share a `try` with plain `except` (a SyntaxError). Hierarchy matters: `ExceptionGroup` subclasses `Exception`, but `BaseExceptionGroup` does not, so `except Exception` misses base-only groups. This is the handler half of `asyncio_task_group`.

### Example

**Before:**

```python
results = await asyncio.gather(fetch(a), fetch(b), return_exceptions=True)
failed = [r for r in results if isinstance(r, Exception)]
if failed:
    raise RuntimeError(f"{len(failed)} requests failed")  # types lost
```

**After:**

```python
try:
    async with asyncio.TaskGroup() as tg:
        tg.create_task(fetch(a))
        tg.create_task(fetch(b))
except* TimeoutError as group:
    retry([exc for exc in group.exceptions])
```

---

<a id="exception_add_note"></a>

## `exception_add_note`

**Category: Errors · Python 3.11+ · Impact: Medium · Autofix: none**

Attach operational context to a propagating exception with `exc.add_note(...)` instead of re-raising a wrapped replacement.

PEP 678 appends notes to the live exception object; each note renders in the default traceback, and repeated calls accumulate rather than overwrite. The original type and traceback stay intact, unlike `raise AppError(...) from exc`, which forces every caller to unwrap a new type; keep explicit `raise ... from` chaining for the cases where a genuinely different exception must be raised.

### Example

**Before:**

```python
try:
    body = json.loads(raw)
except json.JSONDecodeError as exc:
    raise ValueError(f"bad payload from {source}") from exc
    # callers now catch ValueError and lose the decode detail
```

**After:**

```python
try:
    body = json.loads(raw)
except json.JSONDecodeError as exc:
    exc.add_note(f"payload source: {source}")
    raise  # same type, same traceback, richer report
```

---

<a id="asyncio_timeout_cm"></a>

## `asyncio_timeout_cm`

**Category: Async · Python 3.11+ · Impact: Medium · Autofix: none**

Bound a whole block of awaits with `async with asyncio.timeout(...)` instead of wrapping single calls in `asyncio.wait_for`.

The context manager guards every await inside the block, loops, retries, and helper calls alike, and raises the built-in `TimeoutError` on expiry (the same class as `asyncio.TimeoutError` since 3.11; see `timeout_error_alias`). `wait_for` is not deprecated and still fits one-shot wrapping. On expiry the current await is cancelled first, so cleanup code runs before the error surfaces.

### Example

**Before:**

```python
async def sync_document(doc):
    remote = await asyncio.wait_for(download(doc), timeout=5)
    for part in remote.parts:
        await asyncio.wait_for(merge(doc, part), timeout=5)
```

**After:**

```python
async def sync_document(doc):
    async with asyncio.timeout(5):
        remote = await download(doc)
        for part in remote.parts:
            await merge(doc, part)  # every await is covered
```

---

<a id="timeout_error_alias"></a>

## `timeout_error_alias`

**Category: Errors · Python 3.11+ · Impact: Medium · Autofix: ruff:UP041**

Catch and raise the built-in `TimeoutError` instead of the `asyncio.TimeoutError` and `socket.timeout` aliases.

`socket.timeout` (since 3.10) and `asyncio.TimeoutError` (since 3.11) are aliases of the built-in `TimeoutError`, so the built-in is the only name worth writing; ruff:UP041 rewrites both. Unifying shrinks `except` tuples and removes imports that existed only for the alias. `ssl.SSLError` is not in this family: it is an OSError subclass but not an alias, so UP024 (which covers OSError aliases) does not flag it.

### Example

**Before:**

```python
import asyncio

try:
    data = await fetch(url)
except asyncio.TimeoutError:
    data = b""
```

**After:**

```python
try:
    data = await fetch(url)
except TimeoutError:  # same class as socket.timeout and asyncio.TimeoutError
    data = b""
```

---

<a id="enum_strenum"></a>

## `enum_strenum`

**Category: Enums · Python 3.11+ · Impact: Medium · Autofix: ruff:UP042**

Declare string enums with `enum.StrEnum` instead of `class X(str, enum.Enum)`.

StrEnum is the purpose-built str-enum base: members are genuine `str` values and `str(member)` / f-strings render the value (`"red"`), while since 3.11 the str-mixin form renders `Color.RED` instead. Display, logging, and serialization that depend on either rendering must be re-tested; ruff:UP042 rewrites the base class but cannot fix downstream output changes.

### Example

**Before:**

```python
class Color(str, Enum):
    RED = "red"

f"{Color.RED}"  # 3.11: 'Color.RED'
```

**After:**

```python
class Color(StrEnum):
    RED = "red"

f"{Color.RED}"  # 'red'
```

---

<a id="datetime_fromisoformat_iso8601"></a>

## `datetime_fromisoformat_iso8601`

**Category: Datetime · Python 3.11+ · Impact: Medium · Autofix: none**

Parse ISO 8601 timestamps with `datetime.fromisoformat()` instead of `strptime` format strings.

3.11 rewrote `fromisoformat` to accept most ISO 8601: the `Z` suffix, basic dates like `20260928`, week dates, and offsets that previously raised `ValueError` — so the format-string zoo and the `replace("Z", "+00:00")` workaround go away. The docs deliberately say most: non-standard legacy formats still need `strptime`, and floors on 3.10 or older parse only strict `isoformat()` output.

### Example

**Before:**

```python
dt = datetime.strptime("2026-09-28T12:00:00+00:00",
              "%Y-%m-%dT%H:%M:%S%z")
# the "Z" spelling needs a replace() workaround first
```

**After:**

```python
dt = datetime.fromisoformat("2026-09-28T12:00:00Z")
```

---

<a id="typevar_tuple_variadic"></a>

## `typevar_tuple_variadic`

**Category: Typing · Python 3.11+ · Impact: Low · Autofix: none**

Type arity-generic signatures with `TypeVarTuple` (`*Ts`) instead of TypeVar-plus-overload towers.

PEP 646's `TypeVarTuple` and unpacking express signatures that are generic over arity, replacing one-TypeVar-per-shape overload towers; the main payoff is array/numpy-style APIs. On 3.12+ prefer the inline PEP 695 `*Ts` syntax over the assignment form shown once your floor allows it.

### Example

**Before:**

```python
# one TypeVar per arity plus overloads
def scale2(v: tuple[T1, T2], k: float) -> tuple[T1, T2]: ...
def scale3(v: tuple[T1, T2, T3], k: float) -> tuple[T1, T2, T3]: ...
```

**After:**

```python
Ts = TypeVarTuple("Ts")

def scale(*values: *Ts, k: float) -> tuple[*Ts]: ...
```

---

<a id="asyncio_runner"></a>

## `asyncio_runner`

**Category: Async · Python 3.11+ · Impact: Low · Autofix: none**

Run many coroutines on one event loop with `asyncio.Runner` instead of hand-rolled `new_event_loop`/`run_until_complete` plumbing.

`asyncio.Runner` bundles what `asyncio.run` sets up and tears down per call, loop creation, signal handling, and async-generator shutdown, around a loop you reuse across several `run()` invocations. That suits CLIs and test harnesses that execute many top-level coroutines. One-shot entry points should stay on `asyncio.run`.

### Example

**Before:**

```python
loop = asyncio.new_event_loop()
try:
    loop.run_until_complete(warmup())
    loop.run_until_complete(process(jobs))
finally:
    loop.run_until_complete(loop.shutdown_asyncgens())
    loop.close()
```

**After:**

```python
with asyncio.Runner() as runner:
    runner.run(warmup())
    runner.run(process(jobs))  # one loop, managed lifecycle
```

---

<a id="datetime_utc_alias"></a>

## `datetime_utc_alias`

**Category: Datetime · Python 3.11+ · Impact: Low · Autofix: ruff:UP017**

Write `datetime.UTC` instead of the equivalent long alias `datetime.timezone.utc`.

3.11 adds `datetime.UTC` as a plain alias of `datetime.timezone.utc`; ruff UP017 rewrites the long spelling. Nothing is deprecated — the old form stays valid and both names denote the same singleton object, so this is consistency polish rather than a migration, and mixed use across a codebase is harmless.

### Example

**Before:**

```python
def now_utc() -> datetime:
    return datetime.now(datetime.timezone.utc)
```

**After:**

```python
def now_utc() -> datetime:
    return datetime.now(datetime.UTC)
```

---

<a id="union_type_annotation"></a>

## `union_type_annotation`

**Category: Typing · Python 3.10+ · Impact: Critical · Autofix: ruff:UP007**

Write unions with `|` (`X | None`, `A | B`) instead of `Optional[X]` and `Union[A, B]`.

PEP 604's `|` unions are the standard spelling from 3.10 on; ruff UP007 rewrites `Union` and UP045 rewrites `Optional`. Unsafe around runtime-evaluated annotations — `types.UnionType` exists only from 3.10, and libraries such as pydantic must resolve it; ruff marks the fix unsafe on pre-3.10 targets (`lint.pyupgrade.keep-runtime-typing = true` disables it).

### Example 1

**Before:**

```python
def parse(value: Union[int, str]) -> None: ...
```

**After:**

```python
def parse(value: int | str) -> None: ...
```

### Example 2

**Before:**

```python
def find(key: str) -> Optional[Item]: ...
```

**After:**

```python
def find(key: str) -> Item | None: ...
```

---

<a id="dataclass_slots"></a>

## `dataclass_slots`

**Category: Data · Python 3.10+ · Impact: High · Autofix: none**

Generate `__slots__` with `@dataclass(slots=True)` instead of hand-maintained slot tuples.

`@dataclass(slots=True)` derives `__slots__` from the field list, removing hand-maintained slot tuples and the per-instance `__dict__`. The decorator rebuilds the class and returns a new class object — interactions with inheritance chains, `super()` calls, and field defaults need regression tests before switching.

### Example

**Before:**

```python
@dataclass
class Point:
    __slots__ = ("x", "y")
    x: int
    y: int
```

**After:**

```python
@dataclass(slots=True)
class Point:
    x: int
    y: int
```

---

<a id="structural_pattern_matching"></a>

## `structural_pattern_matching`

**Category: Patterns · Python 3.10+ · Impact: High · Autofix: none**

Replace if/elif chains that destructure and dispatch on a value's shape with `match`/`case` patterns.

PEP 634 matches structure: sequence, mapping, literal, and class-attribute patterns test and bind in one `case`. A bare `case name:` is a capture pattern that matches anything and never fails — literal alternatives must be written as literals or dotted paths (`case 404:`, `case Color.RED:`), and class patterns need `__match_args__` or explicit keyword subpatterns.

### Example 1

**Before:**

```python
def handle(command):
    if isinstance(command, tuple) and command[0] == "go":
        direction, distance = command[1], command[2]
        return go(direction, distance)
    elif command == "stop":
        return stop()
    raise ValueError(command)
```

**After:**

```python
def handle(command):
    match command:
        case ("go", direction, distance):
            return go(direction, distance)
        case "stop":
            return stop()
    raise ValueError(command)
```

### Example 2

**Before:**

```python
match command:
    case go:  # meant to match the word "go"; actually captures anything
        move()
```

**After:**

```python
match command:
    case "go":
        move()
```

---

<a id="zip_strict"></a>

## `zip_strict`

**Category: Collections · Python 3.10+ · Impact: High · Autofix: none**

Pass `strict=True` to `zip()` whenever the zipped iterables must have equal lengths.

PEP 618: by default `zip()` stops silently at the shortest input, hiding length bugs such as mismatched columns or dropped records. `zip(a, b, strict=True)` raises `ValueError` the moment one iterable finishes before the other. Code that intentionally truncates to the shorter input will now raise, so migrate those call sites consciously; the ruff check is B905, not a UP autofix.

### Example

**Before:**

```python
names = ["alice", "bob"]
scores = [90, 75, 88]
pairs = list(zip(names, scores))
# the third score is silently dropped
```

**After:**

```python
names = ["alice", "bob"]
scores = [90, 75, 88]
pairs = list(zip(names, scores, strict=True))
# raises ValueError on the length mismatch
```

---

<a id="explicit_encoding_open"></a>

## `explicit_encoding_open`

**Category: Files · Python 3.10+ · Impact: High · Autofix: none**

Open text files with an explicit `encoding="utf-8"` instead of relying on the platform-dependent locale default.

Text I/O opened without an `encoding` argument picks the platform locale, so the same `open(path)` decodes bytes as cp1252 on Windows and UTF-8 elsewhere; PEP 597 added the `-X warn_default_encoding` switch and `EncodingWarning` to flush these sites out. Passing `encoding="utf-8"` documents the on-disk contract. Floor check matters: from 3.15 PEP 686 makes UTF-8 the default and this rule reverses — see `utf8_default_encoding`.

### Example

**Before:**

```python
with open("config.json") as fh:
    config = json.load(fh)  # cp1252 here, utf-8 there
```

**After:**

```python
with open("config.json", encoding="utf-8") as fh:
    config = json.load(fh)
```

---

<a id="isinstance_union"></a>

## `isinstance_union`

**Category: Typing · Python 3.10+ · Impact: Medium · Autofix: none**

Check multiple types with `isinstance(x, int | str)` to match PEP 604 annotation style.

`isinstance` and `issubclass` accept `X | Y` unions from 3.10, keeping runtime checks visually aligned with PEP 604 annotations. The tuple form is not wrong and there is deliberately no autofix here — ruff removed UP038 outright (since 0.13.0) rather than keeping it as a rewrite; adopt the union form only for consistency.

### Example

**Before:**

```python
if isinstance(value, (int, str)): ...
```

**After:**

```python
if isinstance(value, int | str): ...
```

---

<a id="type_alias_explicit"></a>

## `type_alias_explicit`

**Category: Typing · Python 3.10+ · Impact: Medium · Autofix: none**

Mark alias assignments explicitly with `X: TypeAlias = ...` (pre-3.12 floors) instead of bare assignments.

PEP 613's explicit `X: TypeAlias = ...` tells checkers an assignment is a type alias, not an ordinary value. A transitional form — `typing.TypeAlias` is deprecated since 3.12 and ruff UP040 converts it to the `type` statement; projects on a 3.12+ floor should adopt `type X = ...` directly and skip this form.

### Example

**Before:**

```python
Vector = list[float]  # alias or accidental value?
```

**After:**

```python
Vector: TypeAlias = list[float]
```

---

<a id="dataclass_kw_only"></a>

## `dataclass_kw_only`

**Category: Data · Python 3.10+ · Impact: Medium · Autofix: none**

Force keyword-only construction with `@dataclass(kw_only=True)` or the `KW_ONLY` sentinel instead of hand-written `__init__` validation.

`@dataclass(kw_only=True)` makes every field keyword-only, and the `KW_ONLY` sentinel switches only the fields declared after it — no more hand-written `__init__` just to force kwargs. The sentinel itself does not become a field, and inherited fields still order by the dataclass MRO rules, so review subclasses.

### Example

**Before:**

```python
@dataclass
class Conn:
    def __init__(self, *, host: str, port: int) -> None:
        self.host = host
        self.port = port
```

**After:**

```python
@dataclass(kw_only=True)
class Conn:
    host: str
    port: int
```

---

<a id="parenthesized_context_managers"></a>

## `parenthesized_context_managers`

**Category: Syntax · Python 3.10+ · Impact: Medium · Autofix: none**

Combine several context managers in one parenthesized `with (a as x, b as y):` statement.

Since 3.10 the `with` statement accepts an arbitrarily nested parenthesized group, so `contextlib.ExitStack` scaffolding and rightward-drifting nested `with` blocks collapse into one statement. This is pure sugar — managers still enter top-to-bottom and exit in reverse — and on interpreters older than 3.10 the parenthesized form is a `SyntaxError`, unlike the `ExitStack` fallback.

### Example

**Before:**

```python
with open("in.txt") as src:
    with open("out.txt", "w") as dst:
        dst.write(src.read())
```

**After:**

```python
with (
    open("in.txt") as src,
    open("out.txt", "w") as dst,
):
    dst.write(src.read())
```

---

<a id="itertools_pairwise"></a>

## `itertools_pairwise`

**Category: Collections · Python 3.10+ · Impact: Medium · Autofix: none**

Walk adjacent items with `itertools.pairwise(it)` instead of `zip(xs, xs[1:])`.

`pairwise()` accepts any iterable — generators, iterators, file objects — and pulls items lazily, while the `zip(xs, xs[1:])` idiom demands a re-sliceable sequence and copies all but the first element. It yields `(previous, current)` tuples and yields nothing for inputs with fewer than two items. It is 3.10-only, so code shared with older floors keeps the slice idiom.

### Example

**Before:**

```python
for prev, cur in zip(values, values[1:]):
    if cur < prev:
        report_drop(prev, cur)
```

**After:**

```python
from itertools import pairwise

for prev, cur in pairwise(values):
    if cur < prev:
        report_drop(prev, cur)
```

---

<a id="typeguard_narrowing"></a>

## `typeguard_narrowing`

**Category: Typing · Python 3.10+ · Impact: Low · Autofix: none**

Annotate user-defined narrowing predicates with `TypeGuard[T]` instead of plain `bool` returns.

A `TypeGuard[T]` return annotation tells checkers that a True result narrows the argument inside the taken branch. Narrowing is one-directional — the else branch does not get the inverted type; from 3.13 prefer `TypeIs[T]`, which narrows both branches like `isinstance`.

### Example

**Before:**

```python
def is_str_list(values: list[object]) -> bool: ...

if is_str_list(items):
    reveal_type(items)  # list[object] — no narrowing
```

**After:**

```python
def is_str_list(values: list[object]) -> TypeGuard[list[str]]: ...

if is_str_list(items):
    reveal_type(items)  # list[str]
```

---

<a id="paramspec_callable"></a>

## `paramspec_callable`

**Category: Typing · Python 3.10+ · Impact: Low · Autofix: none**

Preserve decorated signatures with `ParamSpec`/`Concatenate` instead of `Callable[..., Any]`.

`ParamSpec` captures a wrapped callable's full parameter list so decorated functions keep signature checking instead of collapsing to `Callable[..., Any]`. Requires checker support (pyright/mypy); on 3.12+ the mechanism composes with PEP 695's inline `**P` syntax.

### Example

**Before:**

```python
def logged(func: Callable[..., Any]) -> Callable[..., Any]: ...

@logged
def add(a: int, b: int) -> int: ...  # signature no longer checked
```

**After:**

```python
P = ParamSpec("P")
R = TypeVar("R")

def logged(func: Callable[P, R]) -> Callable[P, R]: ...

@logged
def add(a: int, b: int) -> int: ...  # signature preserved
```

---

<a id="int_bit_count"></a>

## `int_bit_count`

**Category: Integers · Python 3.10+ · Impact: Low · Autofix: none**

Count set bits with `x.bit_count()` instead of `bin(x).count("1")`.

`int.bit_count()` returns the population count of the absolute value — the same result as `bin(x).count("1")`, but without materializing the number's full binary text, which matters for very large integers in hot loops. It is a method on `int`, so `float`, `Decimal`, or numpy scalars raise `AttributeError` until converted with `int()`.

### Example

**Before:**

```python
ones = bin(mask).count("1")
```

**After:**

```python
ones = mask.bit_count()
```

---

<a id="builtin_generic_annotation"></a>

## `builtin_generic_annotation`

**Category: Typing · Python 3.9+ · Impact: Critical · Autofix: ruff:UP006**

Annotate with builtin generics (`list[int]`, `dict[str, int]`) instead of `typing.List` and `typing.Dict`.

PEP 585 made builtin collections subscriptable at runtime in 3.9, so `list[int]` replaces `typing.List[int]`; the typing aliases are deprecated (ruff UP006 rewrites them, UP035 removes the dead imports). Ruff marks the fix unsafe only when the runtime target is below 3.9, where `list[int]` breaks libraries that evaluate annotations at runtime, such as pydantic; on 3.9+ targets the rewrite is runtime-legal — still verify framework support before bulk fixes.

### Example

**Before:**

```python
from typing import Dict, List

def scores(rows: List[Row]) -> Dict[str, int]: ...
```

**After:**

```python
def scores(rows: list[Row]) -> dict[str, int]: ...
```

---

<a id="dict_merge_operator"></a>

## `dict_merge_operator`

**Category: Collections · Python 3.9+ · Impact: Critical · Autofix: none**

Merge and update dicts with the `|` and `|=` operators instead of `{**a, **b}`, `copy()` plus `update()`, or `collections.ChainMap`.

PEP 584 operators: `a | b` returns a new dict with the right operand's keys winning; `a |= b` updates `a` in place. The right operand must be mapping-like (it needs a `keys()` method), otherwise `TypeError`. `|` never mutates, so `merged = a | b` does not reproduce code that relied on `a.update(b)` aliasing through a second reference.

### Example

**Before:**

```python
merged = {**defaults, **overrides}

settings = defaults.copy()
settings.update(overrides)
```

**After:**

```python
merged = defaults | overrides

settings = defaults | overrides
```

---

<a id="pyproject_single_config"></a>

## `pyproject_single_config`

**Category: Tooling · Python 3.9+ · Impact: Critical · Autofix: none**

Keep all project configuration in `pyproject.toml`, build metadata in `[build-system]` and `[project]` plus per-tool `[tool.*]` sections, instead of scattering setup.py, setup.cfg, pytest.ini, and .flake8 files.

One file becomes the single source of truth: `[build-system]` should always be present, new projects declare metadata in the `[project]` table, and each tool reads its own `[tool.<name>]` section. Migration is incremental because every tool section is independent. Main exception: Poetry before 2.0 (pre-2025) does not consume the `[project]` table.

### Example

**Before:**

```python
# setup.py          -> setup(...)
# setup.cfg         -> [metadata], [flake8]
# pytest.ini        -> [pytest]
# .isort.cfg, .pre-commit-config.yaml, requirements.txt ...
```

**After:**

```python
# pyproject.toml: one file, every tool
[build-system]
requires = ["setuptools>=77"]
build-backend = "setuptools.build_meta"

[project]
name = "mypkg"

[tool.pytest.ini_options]
testpaths = ["tests"]
```

---

<a id="no_setup_py_cli"></a>

## `no_setup_py_cli`

**Category: Tooling · Python 3.9+ · Impact: Critical · Autofix: none**

Build and install through `pip` and `python -m build`; never invoke `python setup.py <command>` from the command line.

The four historical invocations, `setup.py install`, `develop`, `sdist`, and `bdist_wheel`, must not be run: they execute arbitrary project code with stale isolation semantics and were deprecated in setuptools 58.3.0. Setup.py itself is not deprecated; it remains valid as a build-backend customization file. Only its life as a command-line interface is over.

### Example

**Before:**

```python
$ python setup.py install
$ python setup.py develop
$ python setup.py sdist bdist_wheel
```

**After:**

```python
$ python -m pip install .
$ python -m pip install --editable .
$ python -m build
```

---

<a id="str_removeprefix_removesuffix"></a>

## `str_removeprefix_removesuffix`

**Category: Strings · Python 3.9+ · Impact: High · Autofix: none**

Trim affixes with `str.removeprefix()` and `str.removesuffix()` instead of conditional slicing.

`s.removeprefix(p)` returns `s` unchanged when the prefix is absent, replacing the `s[len(p):] if s.startswith(p) else s` idiom and its off-by-one risks; `bytes`, `bytearray`, and `UserString` gained the same methods. These are literal-affix operations, not `lstrip`/`rstrip` — `"www.python.org".removeprefix("w")` strips the single prefix `"w"` once, never a character class, and never repeats.

### Example

**Before:**

```python
name = "www.python.org"
if name.startswith("www."):
    name = name[len("www."):]
```

**After:**

```python
name = "www.python.org".removeprefix("www.")
```

---

<a id="requires_python_metadata"></a>

## `requires_python_metadata`

**Category: Tooling · Python 3.9+ · Impact: High · Autofix: none**

Declare supported interpreters with `requires-python = ">=3.9"` in `[project]`, not only with trove classifiers.

Classifiers such as `Programming Language :: Python :: 3.12` are search and browsing facets; installers do not enforce them. `requires-python` is the machine-checked constraint that pip and uv evaluate before resolving. Keep the two in sync and let one floor drive both; this field is also what downstream tooling reads to detect the minimum supported version.

### Example

**Before:**

```python
[project]
classifiers = [
  "Programming Language :: Python :: 3.9",
  "Programming Language :: Python :: 3.10",
]
```

**After:**

```python
[project]
requires-python = ">=3.9"
classifiers = [
  "Programming Language :: Python :: 3.9",
  "Programming Language :: Python :: 3.10",
]
```

---

<a id="src_layout"></a>

## `src_layout`

**Category: Tooling · Python 3.9+ · Impact: High · Autofix: none**

Lay packages out under `src/` (`src/<pkg>/`) instead of the flat repository-root layout.

The src layout stops the working directory from shadowing the installed package during tests, so the suite exercises what users install rather than accidental local imports, and it forces the importable surface to be declared explicitly. Development then requires an editable install (`pip install -e .` or `uv sync`); running tests off the raw checkout no longer works.

### Example

**Before:**

```python
mypkg/__init__.py
mypkg/core.py
tests/test_core.py
```

**After:**

```python
src/mypkg/__init__.py
src/mypkg/core.py
tests/test_core.py
```

---

<a id="py_typed_marker"></a>

## `py_typed_marker`

**Category: Tooling · Python 3.9+ · Impact: High · Autofix: none**

Ship a `py.typed` marker file inside every typed package so type checkers consume its inline types.

PEP 561: packages distributing inline types MUST add a `py.typed` marker; without it, checkers ignore or stub your annotations no matter how complete they are. Single-file modules are not supported, so promote the code to a package first. The file must actually ship: declare it as package data, for example `[tool.setuptools.package-data]` with `mypkg = ["py.typed"]`.

### Example

**Before:**

```python
src/mypkg/__init__.py
src/mypkg/models.py  # annotated, but checkers cannot use it
```

**After:**

```python
src/mypkg/__init__.py
src/mypkg/models.py
src/mypkg/py.typed

[tool.setuptools.package-data]
mypkg = ["py.typed"]
```

---

<a id="uv_project_manager"></a>

## `uv_project_manager`

**Category: Tooling · Python 3.9+ · Impact: High · Autofix: none**

Manage environments, dependencies, and tools with `uv` instead of gluing pip, venv, pip-tools, and pipx together.

uv is one fast Rust binary covering project environments, locked installs, interpreter downloads, and tool running, replacing the pip/venv/pip-tools/pipx chain (and optionally poetry, pyenv, twine, virtualenv). Its `.python-version` pin and `requires-python` awareness match how version-aware tooling already reasons about the project, so adoption does not fight existing metadata.

### Example

**Before:**

```python
$ python -m venv .venv
$ source .venv/bin/activate
$ pip install -r requirements-dev.txt
$ pipx install ruff
```

**After:**

```python
$ uv sync                  # environment + locked deps in one step
$ uv run pytest
$ uvx ruff check           # ephemeral tool runner
```

---

<a id="ruff_lint_format"></a>

## `ruff_lint_format`

**Category: Tooling · Python 3.9+ · Impact: High · Autofix: none**

Replace the flake8/isort/black/pyupgrade chain with ruff for both linting (`ruff check --fix`) and formatting (`ruff format`), configured under `[tool.ruff]` in pyproject.toml.

One Rust tool covers lint rules (including the pyupgrade UP set this dataset cites in `autofix`), import sorting, and a black-compatible formatter, collapsing four tools and their config shards into one section. UP006/UP007/UP045 autofixes are unsafe around runtime-annotation libraries such as pydantic, and UP038 has been removed outright (since ruff 0.13.0).

### Example

**Before:**

```python
# .flake8, .isort.cfg, pyproject [tool.black],
# plus pyupgrade --py39-plus in CI: four tools, four configs
```

**After:**

```python
# pyproject.toml
[tool.ruff]
target-version = "py39"

[tool.ruff.lint]
select = ["E", "F", "I", "UP"]

# $ ruff check --fix && ruff format
```

---

<a id="annotated_metadata"></a>

## `annotated_metadata`

**Category: Typing · Python 3.9+ · Impact: Medium · Autofix: none**

Attach metadata via `Annotated[T, meta]` instead of inventing custom wrapper types.

`typing.Annotated[T, meta]` attaches metadata directly to a type instead of forcing bespoke wrapper classes per constraint. Static type checkers ignore the extra objects — the consumers are runtime frameworks (pydantic constraints, FastAPI dependencies) that opt in to reading them; the plain type must stay the first argument.

### Example

**Before:**

```python
# bespoke wrapper type per constraint
def set_port(port: PortInt) -> None: ...
```

**After:**

```python
from typing import Annotated

def set_port(port: Annotated[int, Field(ge=0, le=65535)]) -> None: ...
```

---

<a id="typing_text_str"></a>

## `typing_text_str`

**Category: Typing · Python 3.9+ · Impact: Medium · Autofix: ruff:UP019**

Use `str` instead of the Python 2 compatibility alias `typing.Text`.

`typing.Text` exists only for Python 2/3 compatibility and is listed with the deprecated typing aliases; `str` is the only spelling 3.9+ code should use, and ruff UP019 replaces it everywhere. The typing docs state no removal is currently planned, but keeping it only preserves dead compatibility surface.

### Example

**Before:**

```python
from typing import Text

def read(path: str) -> Text: ...
```

**After:**

```python
def read(path: str) -> str: ...
```

---

<a id="zoneinfo_stdlib"></a>

## `zoneinfo_stdlib`

**Category: Datetime · Python 3.9+ · Impact: Medium · Autofix: none**

Attach IANA time zones with `zoneinfo.ZoneInfo` instead of the third-party `pytz`.

PEP 615 ships the tz database reader in the standard library: `ZoneInfo("Europe/Berlin")` instances are cached and interned per key, and slot straight into `datetime(..., tzinfo=...)` and `astimezone()`. Pytz idioms do not port — `pytz.localize()` and `normalize()` have no equivalent because `ZoneInfo` handles DST transitions itself — and platforms without system tzdata (notably Windows) need the `tzdata` package as a fallback.

### Example

**Before:**

```python
import pytz

naive = datetime(2026, 9, 28, 12, 0)
berlin = pytz.timezone("Europe/Berlin").localize(naive)
```

**After:**

```python
from zoneinfo import ZoneInfo

berlin = datetime(2026, 9, 28, 12, 0, tzinfo=ZoneInfo("Europe/Berlin"))
```

---

<a id="asyncio_to_thread"></a>

## `asyncio_to_thread`

**Category: Async · Python 3.9+ · Impact: Medium · Autofix: none**

Offload blocking calls with `asyncio.to_thread(fn, *args, **kwargs)` instead of `loop.run_in_executor(None, fn, *args)`.

`to_thread` targets the default executor, forwards keyword arguments (which `run_in_executor` cannot), and copies the current contextvars into the worker thread, so logging correlation and request-scoped state survive the hop. Reserve explicit `run_in_executor` with a `ProcessPoolExecutor` for CPU-bound work; `to_thread` is for blocking I/O and other GIL-releasing calls.

### Example

**Before:**

```python
loop = asyncio.get_running_loop()
data = await loop.run_in_executor(None, read_blocking, path)
```

**After:**

```python
data = await asyncio.to_thread(read_blocking, path, mode="rb")
```

---

<a id="spdx_license_expression"></a>

## `spdx_license_expression`

**Category: Tooling · Python 3.9+ · Impact: Medium · Autofix: none**

Declare licensing with SPDX expressions, `license = "MIT"` plus `license-files = ["LICEN[CS]E*"]`, instead of the deprecated table form or license classifiers.

PEP 639 deprecates `license = {text = "..."}` and license classifiers in favor of valid SPDX license expressions, with `license-files` globs replacing per-file wiring. The new syntax needs a recent backend, setuptools >= 77.0.3 or hatchling >= 1.27, so raise the `build-system.requires` floor in the same change or older builds will fail to parse the table.

### Example

**Before:**

```python
[project]
license = { text = "MIT" }
classifiers = ["License :: OSI Approved :: MIT License"]
```

**After:**

```python
[project]
license = "MIT"
license-files = ["LICEN[CS]E*"]
```

---

<a id="universal_lockfile"></a>

## `universal_lockfile`

**Category: Tooling · Python 3.9+ · Impact: Medium · Autofix: none**

Commit a universal `uv.lock` and reproduce environments with `uv sync` / `uv run` instead of hand-pinned requirements.txt matrices.

uv's lockfile is cross-platform: one resolution records the wheel and sdist choices for every platform, so Linux CI and macOS laptops install identical dependency graphs from the same file. Commit it to version control and cache it in CI. Do not keep a parallel pip-tools hash-pinning flow on the same project; one lock authority only.

### Example

**Before:**

```python
# requirements-linux.txt
requests==2.32.3

# requirements-macos.txt (drifting)
requests==2.32.0
```

**After:**

```python
# uv.lock: committed, resolved once for every platform
$ uv sync   # installs the exact locked graph everywhere
```

---

<a id="pytest_config_in_pyproject"></a>

## `pytest_config_in_pyproject`

**Category: Tooling · Python 3.9+ · Impact: Medium · Autofix: none**

Configure pytest in `pyproject.toml` under `[tool.pytest.ini_options]` instead of maintaining pytest.ini or setup.cfg sections.

pytest >= 6.0 reads `[tool.pytest.ini_options]`, and 9.0 adds a native typed `[tool.pytest]` table, so the test runner needs no dedicated file. Mind precedence: pytest picks the first of pytest.toml, pytest.ini, pyproject.toml, tox.ini, setup.cfg, and first match wins with no merging; delete the old file when migrating or the new section is silently ignored.

### Example

**Before:**

```python
# pytest.ini
[pytest]
testpaths = tests
addopts = -q
```

**After:**

```python
# pyproject.toml
[tool.pytest.ini_options]
testpaths = ["tests"]
addopts = "-q"
```

---

<a id="functools_cache"></a>

## `functools_cache`

**Category: Data · Python 3.9+ · Impact: Medium · Autofix: ruff:UP033**

Cache pure function results with `@functools.cache` instead of `@functools.lru_cache(maxsize=None)`.

`functools.cache` is the 3.9 name for the unbounded cache previously spelled `lru_cache(maxsize=None)`: a drop-in replacement with identical semantics and one less argument to mistype; ruff UP033 rewrites it. The cache remains unbounded — on instance methods it keeps every `self` alive for the process lifetime, so cap growth with `lru_cache(maxsize=...)` on methods of short-lived classes.

### Example

**Before:**

```python
@functools.lru_cache(maxsize=None)
def fib(n: int) -> int: ...
```

**After:**

```python
@functools.cache
def fib(n: int) -> int: ...
```

---

<a id="lru_cache_no_parens"></a>

## `lru_cache_no_parens`

**Category: Data · Python 3.9+ · Impact: Medium · Autofix: ruff:UP011**

Write the bare decorator `@functools.lru_cache` instead of the empty call `@functools.lru_cache()`.

Since 3.8 `lru_cache` works without parentheses when it gets no arguments, so the empty call is dead weight; ruff UP011 strips it. The rewrite is exact only for the truly empty call — a call passing `maxsize` or `typed` keeps its parentheses, and the bare form means the defaults (`maxsize=128`). Filed at the 3.9 dataset floor even though the syntax is one minor older.

### Example

**Before:**

```python
@functools.lru_cache()
def expensive(x: int) -> int: ...
```

**After:**

```python
@functools.lru_cache
def expensive(x: int) -> int: ...
```

---

<a id="graphlib_topological_sort"></a>

## `graphlib_topological_sort`

**Category: Algorithms · Python 3.9+ · Impact: Low · Autofix: none**

Order dependency graphs with `graphlib.TopologicalSorter` instead of a hand-rolled Kahn or DFS toposort.

Pass a `{node: predecessors}` mapping: `static_order()` returns one flat ordering, while `prepare()` plus repeated `get_ready()`/`done()` emits each node the moment its dependencies finish — the incremental mode maps directly onto worker pools. `CycleError` is raised by `prepare()`/`static_order()`, not at construction, so an invalid graph stays silent until you actually order it.

### Example

**Before:**

```python
def topo(graph):
    order, done = [], set()
    while len(done) < len(graph):
        for node, preds in graph.items():
            if node not in done and preds <= done:
                order.append(node)
                done.add(node)
    return order
```

**After:**

```python
from graphlib import TopologicalSorter

order = list(TopologicalSorter(graph).static_order())
```

---

<a id="math_lcm_multiarg"></a>

## `math_lcm_multiarg`

**Category: Integers · Python 3.9+ · Impact: Low · Autofix: none**

Compute shared divisors and periods with `math.gcd(*values)` and `math.lcm(*values)` instead of reduce chains or hand-written lcm helpers.

Since 3.9 both functions accept any number of arguments, replacing `functools.reduce(math.gcd, values)` folds and per-pair lcm helpers. `lcm()` with no arguments returns 1, the identity element. Both still accept integers only — `float` or `Fraction` operands raise `TypeError`, so numeric code outside the integers keeps its own helper.

### Example

**Before:**

```python
def lcm2(a, b):
    return a * b // gcd(a, b)

step = lcm2(lcm2(4, 6), 9)  # hand-chained
```

**After:**

```python
from math import lcm

step = lcm(4, 6, 9)
```
