<!-- Generated from internal/guidelines/data/typescript.json; DO NOT EDIT. -->

# Modern TypeScript Guidelines Explained

This file provides a more detailed description of the features supported in the Modern TypeScript Guidelines.

**Autofix legend:**

- [x] — a mechanical autofix exists; the guideline below names its `tool:rule`
- [ ] — no autofix available

**Impact legend:**

- Critical — found in almost every project, dozens of occurrences
- High — found often, 5–20 occurrences per project
- Medium — found regularly, 1–5 occurrences per project
- Low — found rarely or in specific code

**Axis legend:**

This dataset gates its guidelines on 2 independent version axes: typescript, node. Each guideline belongs to exactly one axis, and the guideline sections below are grouped by axis.

## Guidelines

| Category | Guideline | Autofix | Axis | Version | Impact |
|----------|-----------|------------|------|---------|--------|
| Tooling | [`native_compiler`](#native_compiler) | [ ] | TypeScript | 7.0 | Medium |
| Tooling | [`migrate_deprecated_compiler_options`](#migrate_deprecated_compiler_options) | [ ] | TypeScript | 6.0 | High |
| Syntax | [`subpath_imports`](#subpath_imports) | [ ] | TypeScript | 6.0 | Low |
| Tooling | [`bundler_resolution_with_commonjs`](#bundler_resolution_with_commonjs) | [ ] | TypeScript | 6.0 | Low |
| Tooling | [`target_es2025`](#target_es2025) | [ ] | TypeScript | 6.0 | Medium |
| Syntax | [`import_defer`](#import_defer) | [ ] | TypeScript | 5.9 | Low |
| Tooling | [`module_node20`](#module_node20) | [ ] | TypeScript | 5.9 | Medium |
| Tooling | [`erasable_syntax_only`](#erasable_syntax_only) | [ ] | TypeScript | 5.8 | High |
| Modules | [`nodenext_require_esm`](#nodenext_require_esm) | [ ] | TypeScript | 5.8 | Medium |
| Tooling | [`rewrite_relative_import_paths`](#rewrite_relative_import_paths) | [ ] | TypeScript | 5.7 | Low |
| Tooling | [`no_unchecked_side_effect_imports`](#no_unchecked_side_effect_imports) | [ ] | TypeScript | 5.6 | Medium |
| Lib & Target | [`iterator_helpers`](#iterator_helpers) | [ ] | TypeScript | 5.6 | Medium |
| Tooling | [`no_check`](#no_check) | [ ] | TypeScript | 5.6 | Low |
| Type System | [`inferred_type_predicates`](#inferred_type_predicates) | [ ] | TypeScript | 5.5 | High |
| Type System | [`typed_set_methods`](#typed_set_methods) | [ ] | TypeScript | 5.5 | Medium |
| Tooling | [`config_dir_template`](#config_dir_template) | [ ] | TypeScript | 5.5 | Medium |
| Tooling | [`regex_syntax_checking`](#regex_syntax_checking) | [ ] | TypeScript | 5.5 | Medium |
| Type System | [`jsdoc_import`](#jsdoc_import) | [ ] | TypeScript | 5.5 | Low |
| Type System | [`no_infer`](#no_infer) | [ ] | TypeScript | 5.4 | High |
| Type System | [`preserved_narrowing_closures`](#preserved_narrowing_closures) | [ ] | TypeScript | 5.4 | Medium |
| Lib & Target | [`group_by`](#group_by) | [ ] | TypeScript | 5.4 | Medium |
| Syntax | [`import_attributes`](#import_attributes) | [ ] | TypeScript | 5.3 | Medium |
| Syntax | [`switch_true_narrowing`](#switch_true_narrowing) | [ ] | TypeScript | 5.3 | Low |
| Syntax | [`using_declarations`](#using_declarations) | [ ] | TypeScript | 5.2 | High |
| Syntax | [`decorator_metadata`](#decorator_metadata) | [ ] | TypeScript | 5.2 | Low |
| Syntax | [`standard_decorators`](#standard_decorators) | [ ] | TypeScript | 5.0 | High |
| Type System | [`const_type_parameters`](#const_type_parameters) | [ ] | TypeScript | 5.0 | High |
| Tooling | [`module_resolution_bundler`](#module_resolution_bundler) | [ ] | TypeScript | 5.0 | High |
| Tooling | [`verbatim_module_syntax`](#verbatim_module_syntax) | [x] | TypeScript | 5.0 | High |
| Type System | [`enum_union_types`](#enum_union_types) | [ ] | TypeScript | 5.0 | Medium |
| Type System | [`satisfies_operator`](#satisfies_operator) | [ ] | TypeScript | 4.9 | High |
| Type System | [`in_operator_narrowing`](#in_operator_narrowing) | [ ] | TypeScript | 4.9 | Medium |
| Syntax | [`auto_accessors`](#auto_accessors) | [ ] | TypeScript | 4.9 | Medium |
| Language (V8) | [`regexp_escape`](#regexp_escape) | [ ] | Node | 24 | Medium |
| Language (V8) | [`explicit_resource_management`](#explicit_resource_management) | [ ] | Node | 24 | High |
| Tooling | [`type_stripping`](#type_stripping) | [ ] | Node | 23.6 | High |
| Test | [`mock_timers`](#mock_timers) | [ ] | Node | 23.1 | Medium |
| Standard Library | [`fs_glob`](#fs_glob) | [ ] | Node | 22.17 | Medium |
| Test | [`partial_deep_strict_equal`](#partial_deep_strict_equal) | [ ] | Node | 22.17 | Medium |
| Standard Library | [`style_text`](#style_text) | [ ] | Node | 22.13 | Medium |
| Standard Library | [`sqlite_module`](#sqlite_module) | [ ] | Node | 22.13 | Low |
| Modules | [`require_esm`](#require_esm) | [ ] | Node | 22.12 | High |
| Standard Library | [`fs_cp`](#fs_cp) | [ ] | Node | 22.3 | Medium |
| Language (V8) | [`array_from_async`](#array_from_async) | [ ] | Node | 22 | Medium |
| Language (V8) | [`promise_with_resolvers`](#promise_with_resolvers) | [ ] | Node | 22 | Low |
| Language (V8) | [`set_methods`](#set_methods) | [ ] | Node | 22 | Medium |
| Tooling | [`watch_mode`](#watch_mode) | [ ] | Node | 22 | Medium |
| Tooling | [`node_run`](#node_run) | [ ] | Node | 22 | Low |
| Standard Library | [`global_websocket`](#global_websocket) | [ ] | Node | 22 | Medium |
| Standard Library | [`global_fetch`](#global_fetch) | [ ] | Node | 21 | High |
| Tooling | [`env_file`](#env_file) | [ ] | Node | 20.6 | Medium |
| Modules | [`import_meta_resolve`](#import_meta_resolve) | [ ] | Node | 20.6 | Medium |
| Standard Library | [`parse_args`](#parse_args) | [ ] | Node | 20 | Medium |
| Test | [`test_runner`](#test_runner) | [ ] | Node | 20 | High |
| Standard Library | [`structured_clone`](#structured_clone) | [ ] | Node | 17 | Medium |
| Modules | [`node_prefix_imports`](#node_prefix_imports) | [ ] | Node | 16 | Medium |
| Modules | [`esm_explicit_extensions`](#esm_explicit_extensions) | [ ] | Node | 16 | High |

---

## Axis: TypeScript

---

<a id="native_compiler"></a>

## `native_compiler`

**Category: Tooling · TypeScript 7.0+ · Impact: Medium · Autofix: none**

Run type checks and builds with the TypeScript 7 native compiler instead of the JavaScript `tsc` where speed matters.

TypeScript 7 is the native Go port: shared-memory parallelism and native code make type checking and builds often about 10x faster than TypeScript 6.0. Language capability tracks 6.0, but behavioral differences exist (for example, template literal types preserve Unicode code points). The npm package installs side-by-side with 6.0 — adopt it for CI and editor speedups, and do not blind-switch every tool in the chain.

### Example

**Before:**

```typescript
tsc --noEmit
# JavaScript-implemented TypeScript 6.x: minutes on large repositories
```

**After:**

```typescript
tsc --noEmit
# TypeScript 7 native (Go) compiler, installed side-by-side with 6.0:
# often about 10x faster, same CLI surface
```

---

<a id="migrate_deprecated_compiler_options"></a>

## `migrate_deprecated_compiler_options`

**Category: Tooling · TypeScript 6.0+ · Impact: High · Autofix: none**

Migrate off `--outFile` and amd/umd/system modules before upgrading to TypeScript 6.0, and off the options it deprecates (`--moduleResolution node10`, `--baseUrl`, `--target es5`) before TypeScript 7.0.

TypeScript 6.0 is the deprecation-cleanup release on the road to 7.0, and it splits this list in two. Deprecated in 6.0 (usable behind `ignoreDeprecations: "6.0"`, removed entirely in 7.0): node10 resolution, `baseUrl`, the `es5` target. Removed in 6.0: `--outFile` is gone and the amd/umd/systemjs/none module values are no longer supported, so migrate off both before upgrading to 6.0. Moving off `baseUrl` means rewriting extensionless absolute imports to real relative paths, with `bundler` or `nodenext` resolution as the new truth.

### Example

**Before:**

```typescript
// tsconfig.json
{
  "compilerOptions": {
    "moduleResolution": "node10",
    "baseUrl": "src",
    "target": "es5",
    "outFile": "build/bundle.js"
  }
}
```

**After:**

```typescript
// tsconfig.json
{
  "compilerOptions": {
    "moduleResolution": "bundler", // or "nodenext" for Node-run code
    "module": "esnext",
    "target": "es2015" // or newer
  }
}
// "@/util" style imports become relative paths such as "../util.js"
```

---

<a id="subpath_imports"></a>

## `subpath_imports`

**Category: Syntax · TypeScript 6.0+ · Impact: Low · Autofix: none**

Use package.json `#` subpath imports instead of deep relative paths for intra-package imports.

TypeScript 6.0 supports the package.json `imports` field, so `#/shared/result` replaces `../../../shared/result` and survives file moves during refactors. The feature needs the `imports` map in package.json plus a tsconfig resolution mode that honors it (`bundler` or `nodenext`); classic node10 resolution ignores the field.

### Example

**Before:**

```typescript
import { Result } from "../../../shared/result.js";
```

**After:**

```typescript
import { Result } from "#/shared/result.js";

// package.json
// { "imports": { "#/shared/*": "./src/shared/*" } }
```

---

<a id="bundler_resolution_with_commonjs"></a>

## `bundler_resolution_with_commonjs`

**Category: Tooling · TypeScript 6.0+ · Impact: Low · Autofix: none**

Pair `--moduleResolution bundler` with `--module commonjs` when a bundler consumes sources but the output stays CommonJS.

Before 6.0, `--moduleResolution bundler` required `--module esnext`, forcing toolchains that bundle TypeScript sources yet ship CommonJS output to masquerade as `nodenext`. TypeScript 6.0 lifts the restriction and allows the `bundler` plus `commonjs` combination. Only relevant to bundler pipelines that emit CJS; code executed directly by Node should stay on `nodenext`.

### Example

**Before:**

```typescript
// tsconfig.json — pre-6.0: had to claim ESM output
{ "compilerOptions": { "module": "esnext", "moduleResolution": "bundler" } }
```

**After:**

```typescript
// tsconfig.json — 6.0+: bundler resolution, CommonJS emit
{ "compilerOptions": { "module": "commonjs", "moduleResolution": "bundler" } }
```

---

<a id="target_es2025"></a>

## `target_es2025`

**Category: Tooling · TypeScript 6.0+ · Impact: Medium · Autofix: none**

Raise `target`/`lib` to `es2025` on TypeScript 6.0+ instead of hand-writing types for new platform features.

TypeScript 6.0 adds the `es2025` target and lib with types for `RegExp.escape` and Temporal. Types are not runtime: `RegExp.escape` only executes on runtimes with V8 13.6 (Node 24 — see the node axis `regexp_escape` rule), and Temporal runtime support is still settling. Keep runtime gates in mind when the compiler accepts the call.

### Example

**Before:**

```typescript
// tsconfig: { "target": "es2022", "lib": ["es2022"] }
const escaped = myEscapeHelper(input); // hand-written, hand-typed
```

**After:**

```typescript
// tsconfig: { "target": "es2025", "lib": ["es2025"] }
const escaped = RegExp.escape(input);
```

---

<a id="import_defer"></a>

## `import_defer`

**Category: Syntax · TypeScript 5.9+ · Impact: Low · Autofix: none**

Defer namespace imports with `import defer * as ns from` when module initialization cost dominates.

`import defer * as ns from "mod"` postpones module evaluation until the namespace is first read, cutting the startup cost of dependencies whose initialization side effects you never use. It requires runtime support for deferred imports, so check the execution target before relying on it. The benefit is real only when module initialization is the measured cost.

### Example

**Before:**

```typescript
import * as reportGenerator from "./report-generator.js";
// module body of report-generator runs at startup, even if unused here
```

**After:**

```typescript
import defer * as reportGenerator from "./report-generator.js";
// evaluation deferred until reportGenerator is first read
```

---

<a id="module_node20"></a>

## `module_node20`

**Category: Tooling · TypeScript 5.9+ · Impact: Medium · Autofix: none**

Pin `--module node20` (or `node18`) instead of `nodenext` to freeze Node resolution semantics across compiler upgrades.

`--module nodenext` tracks whatever Node semantics the newest compiler knows, so upgrading TypeScript can silently change how imports resolve. TypeScript 5.9 adds `--module node20` (5.8 added `node18`) to pin resolution semantics to one Node line, keeping compiler upgrades from drifting resolution behavior. It requires the matching resolution mode.

### Example

**Before:**

```typescript
// tsconfig.json
{ "compilerOptions": { "module": "nodenext" } }
// resolution semantics drift with every compiler upgrade
```

**After:**

```typescript
// tsconfig.json
{ "compilerOptions": { "module": "node20" } }
// pinned to Node 20 resolution semantics
```

---

<a id="erasable_syntax_only"></a>

## `erasable_syntax_only`

**Category: Tooling · TypeScript 5.8+ · Impact: High · Autofix: none**

Enable `--erasableSyntaxOnly` to keep the codebase free of TypeScript syntax that emits runtime code.

`--erasableSyntaxOnly` rejects syntax with runtime semantics: enums, constructor parameter properties, and namespaces containing runtime code must move to erasable forms. This is the syntax precondition for running `.ts` directly on Node type stripping, and Node's recommended tsconfig includes it (see the node axis `type_stripping` rule). The migration cost is rewriting existing enums and parameter properties.

### Example

**Before:**

```typescript
enum LogLevel {
  Debug,
  Info,
}

class Logger {
  constructor(private readonly level: LogLevel = LogLevel.Info) {}
}
```

**After:**

```typescript
const LogLevel = {
  Debug: "debug",
  Info: "info",
} as const;

type LogLevel = (typeof LogLevel)[keyof typeof LogLevel];

class Logger {
  private readonly level: LogLevel;
  constructor(level: LogLevel = LogLevel.Info) {
    this.level = level;
  }
}
```

---

<a id="nodenext_require_esm"></a>

## `nodenext_require_esm`

**Category: Modules · TypeScript 5.8+ · Impact: Medium · Autofix: none**

Type `require()` calls that load ESM under `--module nodenext` instead of casting the result to `any`.

Under `--module nodenext`, TypeScript 5.8 models `require()` of an ES module inside CommonJS files, so the call gets real types instead of hand-written casts. The runtime counterpart is Node's `require(esm)`, unflagged in v22.12.0 (see the node axis `require_esm` rule). It only works for synchronous ESM graphs: a required module with top-level await throws `ERR_REQUIRE_ASYNC_MODULE`.

### Example

**Before:**

```typescript
// CommonJS file, --module nodenext, TypeScript < 5.8
const { format } = require("./format.js") as {
  format: (value: number) => string;
};
```

**After:**

```typescript
// CommonJS file, --module nodenext, TypeScript 5.8+
// types now come from ./format.js exports (Node 22.12+ require(esm))
const { format } = require("./format.js");
```

---

<a id="rewrite_relative_import_paths"></a>

## `rewrite_relative_import_paths`

**Category: Tooling · TypeScript 5.7+ · Impact: Low · Autofix: none**

Enable `--rewriteRelativeImportPaths` so declaration emit rewrites relative `.ts` import specifiers to `.js`.

`--rewriteRelativeImportPaths` makes emitted declaration files rewrite relative `.ts` import specifiers to `.js`, removing publish-time path hacks in layouts where source and output extensions differ. It is aimed at complex outDir or monorepo layouts; projects with a single output directory that already import with `.js` extensions do not need it.

### Example

**Before:**

```typescript
// src/util.ts contains: import { helper } from "./deps.ts";
// emitted d.ts (flag off) keeps: import { helper } from "./deps.ts";
```

**After:**

```typescript
// tsconfig: { "compilerOptions": { "rewriteRelativeImportPaths": true } }
// emitted d.ts rewrites to: import { helper } from "./deps.js";
```

---

<a id="no_unchecked_side_effect_imports"></a>

## `no_unchecked_side_effect_imports`

**Category: Tooling · TypeScript 5.6+ · Impact: Medium · Autofix: none**

Enable `--noUncheckedSideEffectImports` so unresolvable side-effect import paths fail compilation.

`import "./foo.css"` is normally not resolved or checked, so a typo'd path fails silently at runtime. `--noUncheckedSideEffectImports` turns unresolvable side-effect import paths into compile-time errors. Loader pipelines (CSS modules, virtual modules) must keep their paths statically resolvable, or the flag will reject them.

### Example

**Before:**

```typescript
import "./styling.csss"; // typo: silently ignored without the flag
```

**After:**

```typescript
// tsc --noUncheckedSideEffectImports
import "./styling.css"; // unresolvable paths now fail compilation
```

---

<a id="iterator_helpers"></a>

## `iterator_helpers`

**Category: Lib & Target · TypeScript 5.6+ · Impact: Medium · Autofix: none**

Use the built-in iterator helpers (`map`, `filter`, `take`, …) on iterables instead of hand-written loops or intermediate arrays.

TypeScript 5.6 ships lib types for the built-in iterator helpers, so `iterable.map(...).filter(...).take(3).toArray()` chains lazily without collecting intermediate arrays. Runtime gate: the helpers ship in V8 from Chrome 122 onward; nodejs.org does not state the exact Node debut (UNVERIFIED in the research), so the conservative floor is Node 22, which carries V8 12.4.254.14.

### Example

**Before:**

```typescript
function* numbers() {
  yield 1; yield 2; yield 3; yield 4;
}

const evens: number[] = [];
for (const n of numbers()) if (n % 2 === 0) evens.push(n);
const firstTwo = evens.slice(0, 2);
```

**After:**

```typescript
function* numbers() {
  yield 1; yield 2; yield 3; yield 4;
}

const firstTwo = numbers()
  .filter((n) => n % 2 === 0)
  .take(2)
  .toArray();
```

---

<a id="no_check"></a>

## `no_check`

**Category: Tooling · TypeScript 5.6+ · Impact: Low · Autofix: none**

Split CI into a fast `tsc --noCheck` build pass and a separate `tsc --noEmit` type-check pass.

`tsc --noCheck` emits without type checking, letting large repositories run a fast build pass and a separate `tsc --noEmit` type-check pass in parallel in CI. It does not replace type checking — keep the `--noEmit` job gating merges. See the 5.6 announcement for the exact interaction with `--build` mode.

### Example

**Before:**

```typescript
// package.json
// "build": "tsc -p ." — emits and type checks in one slow pass
```

**After:**

```typescript
// package.json
// "build": "tsc -p . --noCheck",    // fast emit
// "typecheck": "tsc -p . --noEmit" // parallel CI job still gates merges
```

---

<a id="inferred_type_predicates"></a>

## `inferred_type_predicates`

**Category: Type System · TypeScript 5.5+ · Impact: High · Autofix: none**

Drop hand-written `x is T` annotations on simple boolean-returning functions; TypeScript 5.5 infers the type predicate.

TypeScript 5.5 infers type predicates for functions that return boolean narrowing results about their parameters, so plain functions narrow in `filter` and guard positions without a hand-written `x is T` annotation. Arrow functions and function declarations both qualify (the official example is a declaration) when the function has no explicit return type annotation, a single return statement, no parameter mutation, and a boolean return tied to that narrowing; explicit return annotations, multi-return bodies, and more complex guards still need a hand-written `x is T`.

### Example

**Before:**

```typescript
const isString = (value: unknown): value is string =>
  typeof value === "string";

const strings = mixed.filter(isString);
```

**After:**

```typescript
const isString = (value: unknown) => typeof value === "string";

// TypeScript 5.5 infers the predicate: value is string
const strings = mixed.filter(isString);
```

---

<a id="typed_set_methods"></a>

## `typed_set_methods`

**Category: Type System · TypeScript 5.5+ · Impact: Medium · Autofix: none**

Use `Set.prototype.union`/`intersection`/`difference`/`symmetricDifference` instead of hand-written set-operation loops.

TypeScript 5.5 types `Set.prototype.union`, `intersection`, `difference`, `symmetricDifference`, `isSubsetOf`, and friends, replacing hand-written loops that allocate intermediate Sets. Runtime gate: MDN Baseline 2024-06; nodejs.org does not name the exact Node debut (UNVERIFIED in the research), so the conservative floor is Node 22 — mirrored by the node axis `set_methods` rule.

### Example

**Before:**

```typescript
function union<T>(a: Set<T>, b: Set<T>): Set<T> {
  const result = new Set(a);
  for (const value of b) result.add(value);
  return result;
}
const tags = union(primary, secondary);
```

**After:**

```typescript
const tags = primary.union(secondary);
```

---

<a id="config_dir_template"></a>

## `config_dir_template`

**Category: Tooling · TypeScript 5.5+ · Impact: Medium · Autofix: none**

Use the `${configDir}` template variable in shared tsconfigs instead of `../shared`-style relative paths.

Shared base tsconfigs historically broke because relative paths inside them resolve against the extending config's directory. The `${configDir}` template variable resolves against the config file that declares it, fixing monorepo `extends` setups. It is only meaningful inside a base config that other configs extend; a standalone tsconfig gains nothing from it.

### Example

**Before:**

```typescript
// shared/tsconfig.base.json — relative paths break depending on who extends
{ "compilerOptions": { "paths": { "~/*": ["../src/*"] } } }
```

**After:**

```typescript
// shared/tsconfig.base.json
{ "compilerOptions": { "paths": { "~/*": ["${configDir}/src/*"] } } }
```

---

<a id="regex_syntax_checking"></a>

## `regex_syntax_checking`

**Category: Tooling · TypeScript 5.5+ · Impact: Medium · Autofix: none**

Delete try/catch wrappers that only guard RegExp literal syntax; TypeScript 5.5 checks it at compile time.

TypeScript 5.5 validates RegExp literal syntax at compile time, surfacing malformed patterns as editor errors instead of runtime exceptions, so try/catch wrappers that only guarded construction of a fixed literal can go. It checks syntax, not semantics, and dynamically concatenated pattern strings cannot be judged statically.

### Example

**Before:**

```typescript
try {
  const pattern = /word(/ig; // typo: until 5.5, found only at runtime
  console.log(pattern.source);
} catch {
  console.error("bad pattern");
}
```

**After:**

```typescript
const pattern = /word/ig; // TypeScript 5.5 flags malformed literals at compile time
console.log(pattern.source);
```

---

<a id="jsdoc_import"></a>

## `jsdoc_import`

**Category: Type System · TypeScript 5.5+ · Impact: Low · Autofix: none**

Reuse `.ts`-declared types inside JavaScript files with the JSDoc `@import` tag instead of duplicating the shapes.

In JavaScript files under `checkJs`, the JSDoc `@import` tag pulls named types from `.ts` modules, so a JS codebase can reuse declared types instead of duplicating object shapes in comments. It applies only to JavaScript files; TypeScript files import types directly.

### Example

**Before:**

```typescript
// app.js (checkJs)
/** @type {{ id: number, name: string }} */
let user;
```

**After:**

```typescript
// app.js (checkJs)
/** @import { User } from "./types.ts" */
/** @type {User} */
let user;
```

---

<a id="no_infer"></a>

## `no_infer`

**Category: Type System · TypeScript 5.4+ · Impact: High · Autofix: none**

Exclude a parameter from inference with `NoInfer<T>` instead of sentinel defaults or `never` tricks.

`NoInfer<T>` marks a type parameter position as excluded from inference, fixing APIs where an optional parameter or default value polluted inference (the classic `createStreetLight` case) without `never` tricks or sentinel defaults. Do not use it when callers are supposed to influence `T` — bidirectional inference is then the desired behavior.

### Example

**Before:**

```typescript
function createStreetLight<C extends string>(colors: C[], defaultColor?: C) {}
createStreetLight(["red", "yellow", "green"], "blue");
// "blue" widens C — no error before 5.4
```

**After:**

```typescript
function createStreetLight<C extends string>(
  colors: C[],
  defaultColor?: NoInfer<C>,
) {}
createStreetLight(["red", "yellow", "green"], "blue");
// error: "blue" is not assignable to C
```

---

<a id="preserved_narrowing_closures"></a>

## `preserved_narrowing_closures`

**Category: Type System · TypeScript 5.4+ · Impact: Medium · Autofix: none**

Rely on TypeScript 5.4 preserving narrowing inside closures created after the variable's last assignment.

TypeScript 5.4 preserves the narrowing of `let` variables inside closures created after the variable's last assignment, so callbacks no longer need to re-check a condition already established outside. A mutable variable that is reassigned after the closure is created still loses the narrowing — the guarantee covers the last-assignment point only.

### Example

**Before:**

```typescript
function greet(name: string | null) {
  if (name === null) return;
  return [1, 2].map(() => {
    if (name === null) return ""; // re-check required before 5.4
    return `hello ${name}`;
  });
}
```

**After:**

```typescript
function greet(name: string | null) {
  if (name === null) return;
  return [1, 2].map(() => `hello ${name}`); // narrowing preserved, no re-check
}
```

---

<a id="group_by"></a>

## `group_by`

**Category: Lib & Target · TypeScript 5.4+ · Impact: Medium · Autofix: none**

Group collections with `Object.groupBy` / `Map.groupBy` instead of hand-rolled `reduce`.

TypeScript 5.4 types `Object.groupBy` and `Map.groupBy`, replacing `reduce` boilerplate with keyed grouping whose result types are precise. Runtime gate: nodejs.org does not state the Node debut (UNVERIFIED in the research); the conservative floor is Node 22 with V8 12.4, matching MDN Baseline 2024-03.

### Example

**Before:**

```typescript
const byDept = employees.reduce<Record<string, Employee[]>>((acc, e) => {
  (acc[e.department] ??= []).push(e);
  return acc;
}, {});
```

**After:**

```typescript
const byDept = Object.groupBy(employees, (e) => e.department);
```

---

<a id="import_attributes"></a>

## `import_attributes`

**Category: Syntax · TypeScript 5.3+ · Impact: Medium · Autofix: none**

Write import attributes with `with { type: "json" }` — never the removed `assert` form.

Import attributes replaced the `assert` keyword with `with`. This is not stylistic: Node 22 removes import assertions entirely ("esm: drop support for import assertions" in the v22.0.0 notes), so `assert` throws at runtime on current Node and `with` is the only working form.

### Example

**Before:**

```typescript
import data from "./config.json" assert { type: "json" }; // removed in Node 22
```

**After:**

```typescript
import data from "./config.json" with { type: "json" };
```

---

<a id="switch_true_narrowing"></a>

## `switch_true_narrowing`

**Category: Syntax · TypeScript 5.3+ · Impact: Low · Autofix: none**

Collapse if/else discriminant chains into `switch (true)` for per-branch narrowing.

TypeScript 5.3 narrows inside `switch (true)` branches, letting a discriminant chain collapse from if/else ladders into one switch where each case narrows its subject. This is an optional style choice: the narrowing benefit depends on writing case predicates the checker can use.

### Example

**Before:**

```typescript
function area(shape: Circle | Square): number {
  if (shape.kind === "circle") return Math.PI * shape.r ** 2;
  else return shape.side ** 2;
}
```

**After:**

```typescript
function area(shape: Circle | Square): number {
  switch (true) {
    case shape.kind === "circle":
      return Math.PI * shape.r ** 2; // shape: Circle here
    default:
      return shape.side ** 2; // shape: Square here
  }
}
```

---

<a id="using_declarations"></a>

## `using_declarations`

**Category: Syntax · TypeScript 5.2+ · Impact: High · Autofix: none**

Release resources with `using` / `await using` declarations instead of try/finally cleanup.

`using` and `await using` declarations call `Symbol.dispose`/`Symbol.asyncDispose` at scope exit, replacing try/finally cleanup and covering early returns and throws. The syntax lands in TypeScript 5.2, but syntax is not runtime: native execution requires explicit resource management, which Node 24's V8 13.6 ships — Node 22 and below need a polyfill or downlevel emit (node axis `explicit_resource_management`).

### Example

**Before:**

```typescript
{
  const handle = openResource();
  try {
    handle.use();
  } finally {
    handle.close();
  }
}
```

**After:**

```typescript
{
  using handle = openResource(); // Symbol.dispose runs at scope exit
  handle.use();
}
```

---

<a id="decorator_metadata"></a>

## `decorator_metadata`

**Category: Syntax · TypeScript 5.2+ · Impact: Low · Autofix: none**

Expose decorator annotations through `Symbol.metadata` instead of hand-written metadata registries.

Decorator metadata exposes per-decoration metadata through `Symbol.metadata`, letting frameworks read annotations from the decorated element instead of maintaining hand-written WeakMap registries keyed by class. Runtime support is still limited, and it only applies to standard decorators — legacy `--experimentalDecorators` code cannot use it.

### Example

**Before:**

```typescript
const routeMetadata = new WeakMap<object, Record<string, string>>();

function Route(path: string) {
  return (target: object, key: string) => {
    const existing = routeMetadata.get(target) ?? {};
    routeMetadata.set(target, { ...existing, [key]: path });
  };
}
```

**After:**

```typescript
function Route(path: string) {
  return (target: object, key: string, context: ClassMethodDecoratorContext) => {
    (context.metadata as Record<string, string>)[key] = path; // Symbol.metadata
  };
}
```

---

<a id="standard_decorators"></a>

## `standard_decorators`

**Category: Syntax · TypeScript 5.0+ · Impact: High · Autofix: none**

Use TC39 stage-3 standard decorators for new code instead of legacy `--experimentalDecorators`.

TypeScript 5.0 implements TC39 stage-3 decorators, which run natively instead of being transpiled into compiler-specific helpers. Their semantics differ completely from legacy `--experimentalDecorators`, so ecosystems that still require legacy (NestJS and friends) cannot switch yet. Node type stripping rejects decorators outright with a parser error (node axis `type_stripping`) — another reason to prefer standard decorators in new code.

### Example

**Before:**

```typescript
// tsconfig: { "experimentalDecorators": true }
@Component({ selector: "profile" })
class ProfileComponent {}
```

**After:**

```typescript
// tsconfig: standard decorators are the default; legacy flag removed
@register
class ProfileComponent {}
```

---

<a id="const_type_parameters"></a>

## `const_type_parameters`

**Category: Type System · TypeScript 5.0+ · Impact: High · Autofix: none**

Declare `<const T>` on generic functions to infer literal and tuple types without `as const` at call sites.

`<const T>` makes the compiler infer the narrowest literal and tuple types from an argument, removing `as const` at call sites: `firstOf(["a", "b"])` infers `readonly ["a", "b"]`. It only affects inference positions; when the caller passes an explicit type argument, the `const` modifier has nothing to do.

### Example

**Before:**

```typescript
function widest<T>(values: readonly T[]): T[] {
  return [...values];
}
const letters = widest(["a", "b"] as const); // as const needed to keep literals
```

**After:**

```typescript
function widest<const T>(values: readonly T[]): T[] {
  return [...values];
}
const letters = widest(["a", "b"]); // T inferred as readonly ["a", "b"]
```

---

<a id="module_resolution_bundler"></a>

## `module_resolution_bundler`

**Category: Tooling · TypeScript 5.0+ · Impact: High · Autofix: none**

Use `--moduleResolution bundler` for bundler-driven projects instead of `node`/`node10` resolution.

`--moduleResolution bundler` models how bundlers actually resolve: extensionless relative imports and `package.json` `exports`, which `node`/`node10` resolution misrepresents for bundler-driven projects. It cannot be combined with `--module commonjs` until TypeScript 6.0 lifts the restriction, and it is wrong for code Node executes directly — use `nodenext` there.

### Example

**Before:**

```typescript
// tsconfig.json
{ "compilerOptions": { "moduleResolution": "node" } }
// wrong for bundlers: ignores the exports map, forces extensioned imports
```

**After:**

```typescript
// tsconfig.json
{ "compilerOptions": { "module": "esnext", "moduleResolution": "bundler" } }
```

---

<a id="verbatim_module_syntax"></a>

## `verbatim_module_syntax`

**Category: Tooling · TypeScript 5.0+ · Impact: High · Autofix: typescript-eslint:consistent-type-imports**

Enable `--verbatimModuleSyntax` and mark type-only imports with `import type`.

`--verbatimModuleSyntax` enforces one emit discipline: imports and exports are emitted exactly as written, so type-only imports must carry `import type`. Turning it on makes unmarked pure-type imports errors — migrate to `import type` first; typescript-eslint's `consistent-type-imports` fixes that mechanically. It matches the tsconfig Node recommends for running `.ts` directly (node axis `type_stripping`).

### Example

**Before:**

```typescript
import { Config, loadConfig } from "./config.js";
// Config is a type; unmarked, it breaks erasable emit
```

**After:**

```typescript
import type { Config } from "./config.js";
import { loadConfig } from "./config.js";
```

---

<a id="enum_union_types"></a>

## `enum_union_types`

**Category: Type System · TypeScript 5.0+ · Impact: Medium · Autofix: none**

Treat every enum as a union of unique member literal types — reference member types and narrow on them directly.

The 5.0 enum overhaul makes every enum a union of unique member literal types — even enums with computed members — so member types can narrow and appear in type positions. The same overhaul adds two new errors: assigning an out-of-range literal to an enum type, and indirectly mixing string and number enum values.

### Example

**Before:**

```typescript
enum Color { Red, Green, Blue }

declare function fill(color: Color): void;
// before 5.0, member types such as Color.Red were not usable as types
```

**After:**

```typescript
enum Color { Red, Green, Blue }

declare function fill(color: Color.Red): void; // member literal as a type

declare const c: Color;
if (c === Color.Red) { /* c: Color.Red here */ }
```

---

<a id="satisfies_operator"></a>

## `satisfies_operator`

**Category: Type System · TypeScript 4.9+ · Impact: High · Autofix: none**

Use `satisfies` to validate a value against a type without widening its inferred type.

`satisfies` checks that a value is assignable to a type without changing the resulting type of the expression, keeping precise literal information that a plain annotation erases by retyping the variable. The check still contextually types the literal itself — the official palette example relies on this to validate RGB tuples — it just leaves the expression's inferred type in place. Prefer annotations at public API boundaries; overusing `satisfies` costs readability.

### Example

**Before:**

```typescript
const config: Config = { retries: 3 };
```

**After:**

```typescript
const config = { retries: 3 } satisfies Config;
// config.retries keeps the literal type 3
```

---

<a id="in_operator_narrowing"></a>

## `in_operator_narrowing`

**Category: Type System · TypeScript 4.9+ · Impact: Medium · Autofix: none**

Narrow unions on unlisted properties with the `in` operator instead of hand-written guard functions.

When a property is not declared on every member of a union, `if ("key" in value)` narrows to the members that have it, replacing hand-written guards for objects that may carry extra keys. It narrows only unlisted properties — keys already declared on the union's types do not trigger narrowing.

### Example

**Before:**

```typescript
interface RequestWithID { id: number }
interface PlainRequest {}

function hasId(request: RequestWithID | PlainRequest): request is RequestWithID {
  return "id" in request; // hand-written guard, pre-4.9
}

if (hasId(req)) console.log(req.id);
```

**After:**

```typescript
function handle(request: RequestWithID | PlainRequest) {
  if ("id" in request) {
    // TypeScript 4.9 narrows: request is RequestWithID
    console.log(request.id);
  }
}
```

---

<a id="auto_accessors"></a>

## `auto_accessors`

**Category: Syntax · TypeScript 4.9+ · Impact: Medium · Autofix: none**

Replace getter/setter plus backing field with `accessor` class fields.

An `accessor` class field replaces the getter/setter plus backing-field trio with one declaration, emitting the accessors for you, and aligns with standard decorator semantics for `@accessor` decorations. On downlevel targets TypeScript emits different support code — check the emitted output when targeting older runtimes.

### Example

**Before:**

```typescript
class Slide {
  #visible = true;

  get visible(): boolean {
    return this.#visible;
  }

  set visible(value: boolean) {
    this.#visible = value;
  }
}
```

**After:**

```typescript
class Slide {
  accessor visible = true;
}
```

---

## Axis: Node

---

<a id="regexp_escape"></a>

## `regexp_escape`

**Category: Language (V8) · Node 24+ · Impact: Medium · Autofix: none**

Escape user input for regex patterns with `RegExp.escape` instead of hand-written escape functions.

`RegExp.escape` escapes arbitrary text for embedding in a regex, eliminating hand-written escape functions whose edge cases (dashes inside classes, `]` handling) breed bugs. It ships with Node 24's V8 13.6 — the v24.0.0 release notes name it explicitly. TypeScript needs `lib: es2025` (TS 6.0) to type the call (typescript axis `target_es2025`).

### Example

**Before:**

```typescript
function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"); // edge cases bite
}
const pattern = new RegExp(escapeRegExp(userInput));
```

**After:**

```typescript
const pattern = new RegExp(RegExp.escape(userInput)); // Node 24, V8 13.6
```

---

<a id="explicit_resource_management"></a>

## `explicit_resource_management`

**Category: Language (V8) · Node 24+ · Impact: High · Autofix: none**

Run `using`/`await using` natively on Node 24+ instead of shipping a polyfill or avoiding the syntax.

Node 24's V8 13.6 ships explicit resource management, so `using` and `await using` execute natively, calling `Symbol.dispose`/`Symbol.asyncDispose` at scope exit. This is the runtime half of the pair: TypeScript 5.2 provides the syntax (typescript axis `using_declarations`), but Node 22 and below have no built-in support and need a polyfill or downlevel emit.

### Example

**Before:**

```typescript
// Node < 24: no native Symbol.dispose — polyfill, or stick to try/finally
const handle = openLog();
try {
  handle.write(event);
} finally {
  handle.close();
}
```

**After:**

```typescript
// Node 24: V8 13.6 ships explicit resource management
using handle = openLog(); // Symbol.dispose runs at scope exit
handle.write(event);
```

---

<a id="type_stripping"></a>

## `type_stripping`

**Category: Tooling · Node 23.6+ · Impact: High · Autofix: none**

Run `.ts` files directly with Node's built-in type stripping instead of a compile step for dev-time scripts.

Node runs `.ts` natively by erasing types: added v22.6.0 behind a flag, enabled by default since v23.6.0 (v22.18.0), marked stable in v24.12.0/v25.2.0. Only erasable syntax is allowed — enums, runtime namespaces, parameter properties, and import aliases throw `ERR_UNSUPPORTED_TYPESCRIPT_SYNTAX`; decorators are parser errors; `.tsx` is unsupported; tsconfig `paths` is ignored; and no type checking happens, so keep `tsc --noEmit` in CI. Pair with `--erasableSyntaxOnly` and `--verbatimModuleSyntax` (typescript axis).

### Example

**Before:**

```typescript
tsc --outDir dist tool.ts && node dist/tool.js
```

**After:**

```typescript
node tool.ts # erasable syntax only; keep tsc --noEmit in CI
```

---

<a id="mock_timers"></a>

## `mock_timers`

**Category: Test · Node 23.1+ · Impact: Medium · Autofix: none**

Fake time with `node:test`'s `mock.timers` instead of sinon's fake timers.

`node:test`'s `mock.timers` fakes `Date`, `setTimeout`, and friends, replacing sinon's fake timers for time-dependent cases. Added v20.4.0/v18.19.0 as experimental and stable since v23.1.0 — the gate used here. The API changed shape once: since v21.2.0/v20.11.0 the methods take a single options object, so older snippets may use the pre-options form.

### Example

**Before:**

```typescript
import sinon from "sinon";

test("schedules reminder", () => {
  const clock = sinon.useFakeTimers();
  scheduleReminder(60_000);
  clock.tick(60_000);
  clock.restore();
});
```

**After:**

```typescript
import { test, mock } from "node:test";

test("schedules reminder", () => {
  mock.timers.enable({ now: 0 });
  scheduleReminder(60_000);
  mock.timers.tick(60_000);
  mock.timers.reset();
});
```

---

<a id="fs_glob"></a>

## `fs_glob`

**Category: Standard Library · Node 22.17+ · Impact: Medium · Autofix: none**

Match file paths with `fs.glob`/`fs.globSync`/`fs.promises.glob` instead of the glob or fast-glob packages.

`fs.glob`, `fs.globSync`, and `fs.promises.glob` cover common file-matching needs without a dependency: `fs.glob` returns an async iterator, `fs.promises.glob` resolves an array. Added in v22.0.0 as experimental and marked stable in v24.0.0 and v22.17.0 — the stable gate used here, since v22.0 through v22.16 ran it experimentally.

### Example

**Before:**

```typescript
import { glob } from "fast-glob";

const files = await glob("src/**/*.test.ts");
```

**After:**

```typescript
import { glob } from "node:fs/promises";

const files = await glob("src/**/*.test.ts");
// or: import { glob } from "node:fs" for the async-iterator form
```

---

<a id="partial_deep_strict_equal"></a>

## `partial_deep_strict_equal`

**Category: Test · Node 22.17+ · Impact: Medium · Autofix: none**

Assert only the fields you care about with `assert.partialDeepStrictEqual` instead of hand-picking subsets.

`assert.partialDeepStrictEqual(actual, expected)` passes when `expected` deep-matches inside `actual`, so tests lock only the fields they care about instead of asserting whole objects brittlely. Added v23.4.0/v22.13.0 as experimental, stable in v24.0.0/v22.17.0 — the stable gate. v25 refines equality semantics (same-instance Promises, Invalid Date handling).

### Example

**Before:**

```typescript
assert.deepStrictEqual({ id: user.id, email: user.email }, {
  id: 42,
  email: "a@example.com",
});
```

**After:**

```typescript
assert.partialDeepStrictEqual(user, {
  id: 42,
  email: "a@example.com",
});
```

---

<a id="style_text"></a>

## `style_text`

**Category: Standard Library · Node 22.13+ · Impact: Medium · Autofix: none**

Color terminal output with `util.styleText` instead of chalk or picocolors.

`util.styleText(format, text)` colors terminal output without a dependency. Added v21.7.0/v20.12.0 and stable since v23.5.0/v22.13.0 — the stable gate. Only from v22.8.0/v20.18.0 does it respect `isTTY`, `NO_COLOR`, and `FORCE_COLOR`, so on older lines the environment contract does not hold.

### Example

**Before:**

```typescript
import chalk from "chalk";

console.log(chalk.red("failed"), chalk.bold("see logs"));
```

**After:**

```typescript
import { styleText } from "node:util";

console.log(styleText("red", "failed"), styleText("bold", "see logs"));
```

---

<a id="sqlite_module"></a>

## `sqlite_module`

**Category: Standard Library · Node 22.13+ · Impact: Low · Autofix: none**

Prefer the built-in `node:sqlite` module over better-sqlite3 for low-stakes embedded storage, after checking its experimental status.

The built-in `node:sqlite` module (a `DatabaseSync` client) covers embedded SQL without better-sqlite3. Added v22.5.0 behind `--experimental-sqlite`; the flag requirement was lifted in v23.4.0/v22.13.0 — the gate — but the module remains experimental on 22/24, reaching RC only in v25.7.0. It must be imported with the `node:` prefix; weigh stability before replacing better-sqlite3 in production.

### Example

**Before:**

```typescript
import Database from "better-sqlite3";

const db = new Database("app.db");
db.exec("CREATE TABLE notes (id INTEGER PRIMARY KEY, body TEXT)");
```

**After:**

```typescript
import { DatabaseSync } from "node:sqlite"; // node: prefix required

const db = new DatabaseSync("app.db");
db.exec("CREATE TABLE notes (id INTEGER PRIMARY KEY, body TEXT)");
```

---

<a id="require_esm"></a>

## `require_esm`

**Category: Modules · Node 22.12+ · Impact: High · Autofix: none**

Load ES modules from CommonJS with `require()` on Node 22.12+ instead of dynamic `import()` workarounds.

Since v22.12.0, CommonJS `require()` can load ES modules synchronously, letting CJS codebases absorb ESM dependencies incrementally instead of going all-or-nothing on `import()`. It applies to synchronous ESM graphs only — a required module with top-level await throws `ERR_REQUIRE_ASYNC_MODULE`. Detect availability with `process.features.require_module`; for typing, use `--module nodenext` (typescript axis `nodenext_require_esm`).

### Example

**Before:**

```typescript
// CommonJS before 22.12: dynamic import is the only bridge
const { formatBytes } = await import("./format-lib.js");
```

**After:**

```typescript
// CommonJS on Node 22.12+: synchronous require of ESM
const { formatBytes } = require("./format-lib.js");
```

---

<a id="fs_cp"></a>

## `fs_cp`

**Category: Standard Library · Node 22.3+ · Impact: Medium · Autofix: none**

Copy directory trees with `fs.promises.cp(src, dest, { recursive: true })` instead of fs-extra.

`fs.promises.cp(src, dest, { recursive: true })` copies directory trees, retiring the fs-extra dependency for copy operations. `fsPromises.cp` was added v16.7.0 but stayed experimental until v22.3.0 — the stable gate used here. Sync and callback variants exist with the same options.

### Example

**Before:**

```typescript
import fsExtra from "fs-extra";

await fsExtra.copy("dist", "deploy", { recursive: true });
```

**After:**

```typescript
import { cp } from "node:fs/promises";

await cp("dist", "deploy", { recursive: true });
```

---

<a id="array_from_async"></a>

## `array_from_async`

**Category: Language (V8) · Node 22+ · Impact: Medium · Autofix: none**

Collect async iterables with `Array.fromAsync` instead of manual `for await` accumulation loops.

`Array.fromAsync(iterable)` collects an async iterable into an array in one expression, replacing manual `for await` accumulation. It awaits each yielded value sequentially — MDN contrasts it with `Promise.all`, which awaits concurrently — so it is not a parallelism tool. Node's exact debut is UNVERIFIED in the research (nodejs.org is silent; MDN's compat table did not render); the conservative gate is Node 22, which ships V8 12.4.254.14.

### Example

**Before:**

```typescript
const results: Result[] = [];
for await (const item of source) {
  results.push(transform(item));
}
```

**After:**

```typescript
const results = await Array.fromAsync(source, transform);
```

---

<a id="promise_with_resolvers"></a>

## `promise_with_resolvers`

**Category: Language (V8) · Node 22+ · Impact: Low · Autofix: none**

Create externally-resolvable promises with `Promise.withResolvers()` instead of external `let resolve`/`let reject` declarations.

`Promise.withResolvers()` returns `{ promise, resolve, reject }` with the controllers in the same scope as the promise, removing the external `let resolve; let reject;` dance in event-driven code. Node's exact debut is UNVERIFIED in the research (nodejs.org is silent; MDN's compat table did not render); the conservative gate is Node 22 (V8 12.4.254.14, MDN Baseline 2024-03). TypeScript needs a recent lib for the types.

### Example

**Before:**

```typescript
let resolve!: (value: number) => void;
let reject!: (reason: unknown) => void;
const ready = new Promise<number>((res, rej) => {
  resolve = res;
  reject = rej;
});
```

**After:**

```typescript
const { promise: ready, resolve, reject } = Promise.withResolvers<number>();
```

---

<a id="set_methods"></a>

## `set_methods`

**Category: Language (V8) · Node 22+ · Impact: Medium · Autofix: none**

Compute set unions, intersections, and differences with `Set.prototype.union` and friends instead of hand-written loops.

`Set.prototype.union`, `intersection`, `difference`, `symmetricDifference`, `isSubsetOf`, `isSupersetOf`, and `isDisjointFrom` make set algebra engine-guaranteed instead of loop-and-hope. Node's exact debut is UNVERIFIED in the research (nodejs.org does not name it; MDN Baseline 2024-06); the conservative gate is Node 22 with V8 12.4. TypeScript types arrive in 5.5 (typescript axis `typed_set_methods`).

### Example

**Before:**

```typescript
function difference<T>(a: Set<T>, b: Set<T>): Set<T> {
  const out = new Set(a);
  for (const value of b) out.delete(value);
  return out;
}

const stale = difference(installed, active);
```

**After:**

```typescript
const stale = installed.difference(active);
```

---

<a id="watch_mode"></a>

## `watch_mode`

**Category: Tooling · Node 22+ · Impact: Medium · Autofix: none**

Restart on file changes with `node --watch` instead of nodemon.

`node --watch` restarts the process when watched files change, replacing nodemon for development loops. Added v18.11.0 as experimental ("cli: add --watch") and marked stable in v22.0.0 — the gate. Watch is process restart, not hot reload; the test runner has its own `node --test --watch` mode.

### Example

**Before:**

```typescript
npx nodemon server.js
```

**After:**

```typescript
node --watch server.js
```

---

<a id="node_run"></a>

## `node_run`

**Category: Tooling · Node 22+ · Impact: Low · Autofix: none**

Run package.json scripts with `node --run <script>` instead of `npm run <script>` where npm's startup cost dominates.

`node --run <script>` executes a package.json script without spawning an npm process, saving npm's startup cost in loops and CI. Added v22.0.0. It does not run npm lifecycle hooks (`pre`/`post` scripts), so it is not a drop-in replacement where hooks do real work.

### Example

**Before:**

```typescript
npm run generate
```

**After:**

```typescript
node --run generate # skips npm startup; no lifecycle hooks
```

---

<a id="global_websocket"></a>

## `global_websocket`

**Category: Standard Library · Node 22+ · Impact: Medium · Autofix: none**

Use the global `WebSocket` client instead of the `ws` package for client-side connections.

A standards-shaped `WebSocket` client is global since v22.0.0 (exposed v21.0.0/v20.10.0, flag removed v22.0.0, no longer experimental v22.4.0), dropping the `ws` dependency for client-side connections. It is a client API: `ws` remains the tool for servers and for advanced options the global class does not expose.

### Example

**Before:**

```typescript
import WebSocket from "ws";

const socket = new WebSocket("wss://example.com/feed");
```

**After:**

```typescript
const socket = new WebSocket("wss://example.com/feed");
```

---

<a id="global_fetch"></a>

## `global_fetch`

**Category: Standard Library · Node 21+ · Impact: High · Autofix: none**

Use the global `fetch` — with `FormData`, `Headers`, `Request`, `Response` — instead of axios or node-fetch.

The global `fetch` became stable in v21.0.0 ("stable fetch and WebStreams"); v18.0.0 removed the flag but kept the experimental label through Node 20, so 21 is the strict gate. It is built on undici: no browser-style cookie jar, and some web semantics differ. `AbortSignal.timeout` (v17.3.0) pairs well with it for deadlines.

### Example

**Before:**

```typescript
import axios from "axios";

const response = await axios.get(url, { timeout: 5000 });
```

**After:**

```typescript
const response = await fetch(url, {
  signal: AbortSignal.timeout(5000),
});
```

---

<a id="env_file"></a>

## `env_file`

**Category: Tooling · Node 20.6+ · Impact: Medium · Autofix: none**

Load environment files with `node --env-file` (or `process.loadEnvFile()`) instead of dotenv.

`node --env-file=.env` loads environment files without dotenv, and `process.loadEnvFile()` (v21.7.0/v20.12.0) does it in-process. `--env-file` was added v20.6.0 and left experimental only in v24.10.0/v22.21.0 — treat 20.6 as the floor. A missing file errors; use `--env-file-if-exists` (v22.9+) for optional files. `NODE_OPTIONS` set inside an env file has no effect.

### Example

**Before:**

```typescript
import dotenv from "dotenv";

dotenv.config();
const port = process.env.PORT;
```

**After:**

```typescript
// package.json: "start": "node --env-file=.env app.js"
const port = process.env.PORT;
```

---

<a id="import_meta_resolve"></a>

## `import_meta_resolve`

**Category: Modules · Node 20.6+ · Impact: Medium · Autofix: none**

Resolve module specifiers with `import.meta.resolve(specifier)` instead of hand-built URL path logic.

`import.meta.resolve(specifier)` resolves a specifier against the current module and returns the URL string synchronously — no `--experimental-import-meta-resolve` flag since v20.6.0/v18.19.0, and the `import.meta` properties graduated in Node 24. It works only inside ES modules, and the non-standard `parentURL` second parameter remains flagged.

### Example

**Before:**

```typescript
const helperUrl = new URL("./helper.js", import.meta.url).href;
// guesses resolution by hand; breaks with exports maps
```

**After:**

```typescript
const helperUrl = import.meta.resolve("./helper.js");
```

---

<a id="parse_args"></a>

## `parse_args`

**Category: Standard Library · Node 20+ · Impact: Medium · Autofix: none**

Parse simple flags and positionals with `util.parseArgs` instead of commander or yargs.

`util.parseArgs` parses options and positional arguments — added v18.3.0/v16.17.0, non-experimental since v20.0.0 (the gate). It deliberately covers only `options` plus `positionals`: no subcommands, no help-text generation, no coercion framework. CLIs needing those still justify commander or yargs.

### Example

**Before:**

```typescript
import { program } from "commander";

program
  .option("--minify")
  .parse();
const minify = program.opts().minify ?? false;
```

**After:**

```typescript
import { parseArgs } from "node:util";

const { values } = parseArgs({
  options: { minify: { type: "boolean" } },
});
const minify = values.minify ?? false;
```

---

<a id="test_runner"></a>

## `test_runner`

**Category: Test · Node 20+ · Impact: High · Autofix: none**

Write tests with the built-in `node:test` runner and `node:assert/strict` instead of jest or mocha when you do not need their extras.

The built-in `node:test` runner runs tests with zero install — added v18.0.0/v16.17.0, stable since v20.0.0 (the gate). `describe` blocks arrive in v20.13.0/v22.0.0; snapshot testing left experimental in v22.3.0/v23.4.0. Its mocking and snapshot story is thinner than jest's — compare your feature usage before migrating.

### Example

**Before:**

```typescript
import { describe, it, expect } from "@jest/globals";

describe("money", () => {
  it("adds", () => expect(add(1, 1)).toBe(2));
});
```

**After:**

```typescript
import { test } from "node:test";
import assert from "node:assert/strict";

test("adds", () => assert.equal(add(1, 1), 2));
```

---

<a id="structured_clone"></a>

## `structured_clone`

**Category: Standard Library · Node 17+ · Impact: Medium · Autofix: none**

Deep-copy values with `structuredClone` instead of `JSON.parse(JSON.stringify(x))`.

`structuredClone` deep-copies while preserving `Map`, `Set`, `Date`, typed arrays, and cycles — everything `JSON.parse(JSON.stringify(x))` silently corrupts. It throws on functions and unsupported objects instead of dropping them, so cloning failures surface loudly; treat an error as a signal that the value is not plain data.

### Example

**Before:**

```typescript
const copy = JSON.parse(JSON.stringify(value));
```

**After:**

```typescript
const copy = structuredClone(value);
```

---

<a id="node_prefix_imports"></a>

## `node_prefix_imports`

**Category: Modules · Node 16+ · Impact: Medium · Autofix: none**

Import Node built-ins through the `node:` prefix, e.g. `require("node:fs")` or `import { readFile } from "node:fs"`.

The `node:` import prefix for `require()` is supported since v16.0.0/v14.18.0. It bypasses the require cache and any shadowing of core module names by installed packages. Some modules require it outright: `node:test` and `node:sqlite` cannot be required unprefixed.

### Example

**Before:**

```typescript
const { readFile } = require("fs");
```

**After:**

```typescript
const { readFile } = require("node:fs");
// or: import { readFile } from "node:fs";
```

---

<a id="esm_explicit_extensions"></a>

## `esm_explicit_extensions`

**Category: Modules · Node 16+ · Impact: High · Autofix: none**

Write ES module relative imports with explicit file extensions.

Node ES modules resolve relative and absolute specifiers strictly: a file extension must be provided, directory indexes must be fully specified, and there are no default extensions or folder mains. In TypeScript the `.js` specifier refers to the emitted file even when the source is `.ts`; `--rewriteRelativeImportExtensions` (TS 5.7) rewrites on emit. Note that `engines.node` is advisory only — it neither blocks install nor enforces the floor, so verify the runtime in CI.

### Example

**Before:**

```typescript
import { helper } from "./utils"; // Node ESM: ERR_MODULE_NOT_FOUND
```

**After:**

```typescript
import { helper } from "./utils.js"; // .js even when the source is utils.ts
```
