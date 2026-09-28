---
name: use-modern-guidelines
description: Use the modern-guidelines (mg) CLI whenever writing, modifying, fixing, or refactoring Python, Go, or TypeScript code. Apply its version-specific guidance to generated changes.
---

# modern-guidelines (mg) CLI

Always write modern, idiomatic code. Use the `mg` CLI as the source of truth for modern Python, Go, and TypeScript idioms that may be newer than your knowledge cutoff.

This skill applies to development of this repository itself. Its main language is Go; the Go detector resolves the version from this repository's `go.mod` `go` directive, so Go code is version-gated by `go.mod` (the provenance line reports, e.g., `# go 1.26 (source: go.mod go directive, fidelity: high)`).

No wrapper or bootstrap script is needed: the tool ships with this repository.

Commands:

- Inside this repository: `go run . <subcommand> ...`
- Once the published `mg` binary is installed on PATH: `mg <subcommand> ...`

Subcommands:

- `list`
- `explain`

Before editing Python, Go, or TypeScript code:

1. Call `list` for the file you are about to edit.

   Inside this repository:

   ```sh
   go run . list --file-path path/to/file.go
   ```

   Or, with the published binary:

   ```sh
   mg list --file-path path/to/file.go
   ```

   The CLI infers the language from the file extension (`.py`, `.go`, `.ts`, `.mts`, `.cts`, `.tsx`) and resolves each axis version from project manifests (pyproject `requires-python` / `.python-version`, `go.mod` `go` directive / go.work, `devDependencies.typescript` / `engines.node` / `.nvmrc`). The output starts with one provenance line per axis (version, source, fidelity).

2. If the target version is already known, or the detector cannot resolve it from manifests, pass it explicitly with `--version axis=<ver>`; this overrides detection:

   ```sh
   go run . list --version python=3.12
   go run . list --version go=1.26
   go run . list --version typescript=5.4,node=18
   ```

   For multi-axis datasets (TypeScript's typescript + node axes) `--version` must cover every axis.

3. Read the complete `list` output before deciding which guidelines apply.

   The list output is ordered newest first. Read the full output because older supported guidelines may still apply.

   Do not pipe the output through head, tail, grep, sed, awk, or any other truncating/filtering command. Important guidelines may otherwise be missed.

4. Treat returned guidelines as authoritative for modern style choices in code you are editing.

   If a guideline applies, follow it even when nearby code or repository convention uses an older pattern. Skip it only when it would not compile, would change behavior, or clearly does not match the edited code. Before skipping a returned guideline that seems relevant, call `explain` for that guideline ID.

Call `explain` only when a specific guideline may apply and you need the detailed explanation or examples. Request only the guideline IDs you intend to evaluate or apply, as `lang:id`:

```sh
go run . explain go:range_over_int
```

Multiple guideline IDs may be requested as positional arguments:

```sh
go run . explain go:range_over_int go:loopvar_capture
```

Do not call `explain` without guideline IDs. Use `list` first to discover the guideline list for the resolved version, then call `explain` for the specific returned IDs that need more context.
