package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// TestFeaturesMarkdownInSync fails if any docs/features/<lang>.md has
// drifted from the generator output. Run `make generate-features` to
// regenerate them (PRD-0000 requirement 11, mirroring the reference
// project's TestFeaturesMarkdownInSync).
func TestFeaturesMarkdownInSync(t *testing.T) {
	entries, err := os.ReadDir("../guidelines/data")
	if err != nil {
		t.Fatal(err)
	}
	generated := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		generated++
		dataset, err := readDataset(filepath.Join("../guidelines/data", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		want := render(dataset)
		got, err := os.ReadFile("../../docs/features/" + dataset.Language + ".md")
		if err != nil {
			t.Fatalf("docs/features/%s.md does not exist; run `make generate-features`: %v", dataset.Language, err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("docs/features/%s.md is out of date; run `make generate-features`", dataset.Language)
		}
	}
	if generated == 0 {
		t.Fatal("no datasets found under ../guidelines/data")
	}
}

const singleAxisDataset = `{
  "language": "go",
  "rules": [
    {
      "id": "safe_id",
      "since_version": "1.25",
      "autofix": "gopls:modernize",
      "category": "Testing",
      "impact": "High",
      "guideline": "Use ` + "`t.Context()`" + ` instead of context.Background in tests.",
      "details": "Keep ` + "`t`" + ` as code.",
      "examples": [{"before": ["old()"], "after": ["new()"]}]
    }
  ]
}`

const multiAxisDataset = `{
  "language": "typescript",
  "axes": ["typescript", "node"],
  "rules": [
    {
      "id": "ts_rule",
      "since_version": "5.5",
      "axis": "typescript",
      "autofix": null,
      "category": "Tooling",
      "impact": "Medium",
      "guideline": "Adopt the new compiler.",
      "details": "Faster checks.",
      "examples": [{"before": ["tsc --noEmit"], "after": ["tsgo --noEmit"]}]
    },
    {
      "id": "node_rule",
      "since_version": "18",
      "axis": "node",
      "autofix": null,
      "category": "Modules",
      "impact": "High",
      "guideline": "Prefer node: prefixed core imports.",
      "details": "Explicit is better.",
      "examples": [{"before": ["require(\"fs\")"], "after": ["require(\"node:fs\")"]}]
    }
  ]
}`

func renderJSON(t *testing.T, data string) []byte {
	t.Helper()
	dataset, err := schema.Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	return render(dataset)
}

func TestRenderSingleAxisDataset(t *testing.T) {
	output := renderJSON(t, singleAxisDataset)

	for _, want := range []string{
		"<!-- Generated from internal/guidelines/data/go.json; DO NOT EDIT. -->\n",
		"# Modern Go Guidelines Explained\n",
		"This file provides a more detailed description of the features supported in the Modern Go Guidelines.\n",
		"| Category | Guideline | Autofix | Go | Impact |\n",
		"| Testing | [`safe_id`](#safe_id) | [x] | 1.25 | High |\n",
		"**Category: Testing · Go 1.25+ · Impact: High · Autofix: gopls:modernize**\n",
		"Use `t.Context()` instead of context.Background in tests.",
		"Keep `t` as code.",
		"### Example\n\n**Before:**\n\n```go\nold()\n```\n",
		"\n**After:**\n\n```go\nnew()\n```\n",
	} {
		if !bytes.Contains(output, []byte(want)) {
			t.Errorf("render() output does not contain %q", want)
		}
	}
	if bytes.Contains(output, []byte("| Axis |")) {
		t.Error("render() added an Axis column to a single-axis dataset")
	}
	if bytes.Contains(output, []byte("## Axis:")) {
		t.Error("render() added an axis heading to a single-axis dataset")
	}
}

func TestRenderMultiAxisDataset(t *testing.T) {
	output := renderJSON(t, multiAxisDataset)

	for _, want := range []string{
		"<!-- Generated from internal/guidelines/data/typescript.json; DO NOT EDIT. -->\n",
		"# Modern TypeScript Guidelines Explained\n",
		"This dataset gates its guidelines on 2 independent version axes: typescript, node.",
		"| Category | Guideline | Autofix | Axis | Version | Impact |\n",
		"| Tooling | [`ts_rule`](#ts_rule) | [ ] | TypeScript | 5.5 | Medium |\n",
		"| Modules | [`node_rule`](#node_rule) | [ ] | Node | 18 | High |\n",
		"## Axis: TypeScript\n",
		"## Axis: Node\n",
		"**Category: Tooling · TypeScript 5.5+ · Impact: Medium · Autofix: none**\n",
		"**Category: Modules · Node 18+ · Impact: High · Autofix: none**\n",
		"### Example\n\n**Before:**\n\n```typescript\ntsc --noEmit\n```\n",
	} {
		if !bytes.Contains(output, []byte(want)) {
			t.Errorf("render() output does not contain %q", want)
		}
	}
	if tsAxis, nodeAxis := bytes.Index(output, []byte("## Axis: TypeScript")), bytes.Index(output, []byte("## Axis: Node")); tsAxis == -1 || nodeAxis == -1 || tsAxis > nodeAxis {
		t.Error("render() did not group the typescript axis before the node axis")
	}
	if tsRule, nodeRule := bytes.Index(output, []byte("## `ts_rule`")), bytes.Index(output, []byte("## `node_rule`")); tsRule == -1 || nodeRule == -1 || tsRule > nodeRule {
		t.Error("render() did not order rule sections by axis group")
	}
}

func TestRenderNumbersMultipleExamples(t *testing.T) {
	output := renderJSON(t, `{
  "language": "python",
  "rules": [
    {
      "id": "two_examples",
      "since_version": "3.12",
      "autofix": null,
      "category": "Syntax",
      "impact": "Low",
      "guideline": "Use the new syntax.",
      "details": "Two fiddles.",
      "examples": [
        {"before": ["a()"], "after": ["b()"]},
        {"before": ["c()"], "after": ["d()"]}
      ]
    }
  ]
}`)

	for _, want := range []string{"### Example 1\n", "### Example 2\n"} {
		if !bytes.Contains(output, []byte(want)) {
			t.Errorf("render() output does not contain %q", want)
		}
	}
	if bytes.Contains(output, []byte("### Example\n")) {
		t.Error("render() left an unnumbered Example heading for a multi-example rule")
	}
}

func TestWriteCodeBlock(t *testing.T) {
	var b strings.Builder
	writeCodeBlock(&b, []string{`fmt.Println("hello")`}, "go")

	want := "```go\nfmt.Println(\"hello\")\n```\n"
	if got := b.String(); got != want {
		t.Fatalf("writeCodeBlock() = %q, want %q", got, want)
	}
}

// chdirTempModule moves the test into a fresh temporary directory that
// carries a go.mod, so run()'s repo-local path checks treat that directory
// as the module root. It returns the module root.
func chdirTempModule(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module featuresgen.test\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	return root
}

func TestRunWritesPerLanguageFiles(t *testing.T) {
	base := chdirTempModule(t)
	dataDir := filepath.Join(base, "data")
	outputDir := filepath.Join(base, "out")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, dataset := range []string{singleAxisDataset, multiAxisDataset} {
		parsed, err := schema.Parse([]byte(dataset))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dataDir, parsed.Language+".json"), []byte(dataset), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := run([]string{dataDir, outputDir}); err != nil {
		t.Fatal(err)
	}
	for _, language := range []string{"go", "typescript"} {
		got, err := os.ReadFile(filepath.Join(outputDir, language+".md"))
		if err != nil {
			t.Fatalf("run() did not write %s.md: %v", language, err)
		}
		dataset, err := readDataset(filepath.Join(dataDir, language+".json"))
		if err != nil {
			t.Fatal(err)
		}
		if want := render(dataset); !bytes.Equal(got, want) {
			t.Errorf("run() wrote %s.md that differs from render(): %d bytes vs %d", language, len(got), len(want))
		}
	}
}

func TestRunSingleDatasetToFile(t *testing.T) {
	base := chdirTempModule(t)
	dataDir := filepath.Join(base, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dataPath := filepath.Join(dataDir, "go.json")
	if err := os.WriteFile(dataPath, []byte(singleAxisDataset), 0o644); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(base, "Go.md")

	if err := run([]string{dataPath, outputPath}); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(outputPath); err != nil || len(got) == 0 {
		t.Fatalf("run() did not write the output file: %v", err)
	}
}

func TestRunRejectsPathsOutsideModuleRoot(t *testing.T) {
	base := chdirTempModule(t)
	dataDir := filepath.Join(base, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "go.json"), []byte(singleAxisDataset), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, args := range [][]string{
		{dataDir, "../../outside"},
		{filepath.Join(base, "..", "elsewhere"), filepath.Join(base, "out")},
		{"/tmp/featuresgen-should-not-write-here", filepath.Join(base, "out")},
	} {
		if err := run(args); err == nil || !strings.Contains(err.Error(), "outside the module root") {
			t.Errorf("run(%q) error = %v, want module-root containment error", args, err)
		}
	}
}

func TestRunRejectsWrongArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"only-one"}, {"a", "b", "c"}} {
		if err := run(args); err == nil || !strings.Contains(err.Error(), "usage: featuresgen") {
			t.Errorf("run(%q) error = %v, want usage error", args, err)
		}
	}
}

func TestDisplayName(t *testing.T) {
	for name, want := range map[string]string{
		"go":         "Go",
		"python":     "Python",
		"typescript": "TypeScript",
		"node":       "Node",
		"rust":       "Rust",
		"":           "",
	} {
		if got := displayName(name); got != want {
			t.Errorf("displayName(%q) = %q, want %q", name, got, want)
		}
	}
}
