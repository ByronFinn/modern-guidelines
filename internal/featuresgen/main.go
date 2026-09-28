// Command featuresgen renders every guideline dataset under
// internal/guidelines/data into docs/features/<lang>.md (PRD-0000
// requirement 11). It is the per-language, multi-axis counterpart of the
// reference project's examples/go-modern-guidelines/internal/guidelines/
// featuresgen, whose FEATURES.md structure it mirrors: generated-from
// marker, title, intro, legends, a summary table with anchor links, and one
// section per rule.
//
// Usage:
//
//	featuresgen <data> <output>
//
// <data> is one dataset JSON file or a directory with one dataset per
// language; <output> is the markdown file written for a single dataset, or
// the directory receiving <language>.md for each dataset in a data
// directory. Both arguments must resolve inside the current module (see
// repoLocalPath). Datasets are validated through internal/schema before
// rendering, so invalid data aborts generation instead of producing docs.
//
// The Makefile wires this command as `make generate-features`; the
// equivalent //go:generate directive sits at the top of
// internal/guidelines/guidelines.go.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: featuresgen <data-dir-or-json> <output-dir-or-md>")
	}
	dataPath, err := repoLocalPath(args[0])
	if err != nil {
		return fmt.Errorf("data argument: %w", err)
	}
	outputPath, err := repoLocalPath(args[1])
	if err != nil {
		return fmt.Errorf("output argument: %w", err)
	}
	dataInfo, err := os.Stat(dataPath)
	if err != nil {
		return fmt.Errorf("stat data: %w", err)
	}
	if dataInfo.IsDir() {
		return runDir(dataPath, outputPath)
	}

	dataset, err := readDataset(dataPath)
	if err != nil {
		return err
	}
	output := outputPath
	if info, err := os.Stat(output); err == nil && info.IsDir() {
		output = filepath.Join(output, dataset.Language+".md")
	}
	return writeRendered(dataset, output)
}

// repoLocalPath cleans a command-line path argument and requires the
// resolved location to stay inside the module root (the nearest ancestor
// directory holding go.mod). This generator only ever reads and writes
// inside its own repository, so arguments that escape it — absolute paths
// elsewhere or ".." chains climbing out — are rejected outright.
// File names derived from dataset content are separately constrained by
// the schema's validName charset (lowercase letters, digits, hyphens), so
// no path component can enter through dataset.Language either.
func repoLocalPath(arg string) (string, error) {
	cleaned := filepath.Clean(arg)
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}
	root, err := moduleRoot()
	if err != nil {
		return "", err
	}
	if !withinDir(root, abs) {
		return "", fmt.Errorf("%q resolves outside the module root %s: only paths inside the current repository are accepted", arg, root)
	}
	return abs, nil
}

// moduleRoot walks up from the working directory to the nearest ancestor
// containing a go.mod file.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if info, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !info.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found from %s upward: run featuresgen inside the module", dir)
		}
		dir = parent
	}
}

// withinDir reports whether path is root itself or lies underneath it.
func withinDir(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// runDir renders every *.json dataset in dataDir, in file-name order, to
// outputDir/<language>.md.
func runDir(dataDir, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("read data directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		dataset, err := readDataset(filepath.Join(dataDir, entry.Name()))
		if err != nil {
			return err
		}
		if err := writeRendered(dataset, filepath.Join(outputDir, dataset.Language+".md")); err != nil {
			return err
		}
	}
	return nil
}

func readDataset(path string) (schema.Dataset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.Dataset{}, fmt.Errorf("read dataset: %w", err)
	}
	dataset, err := schema.Parse(data)
	if err != nil {
		return schema.Dataset{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return dataset, nil
}

func writeRendered(dataset schema.Dataset, path string) error {
	if err := os.WriteFile(path, render(dataset), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// render produces the complete FEATURES markdown for one validated dataset.
// The layout follows the reference project's FEATURES.md: generated-from
// marker, title, one-line intro, autofix and impact legends, the summary
// table with anchor links, then one section per rule. Multi-axis datasets
// add an axis legend, an Axis column in the table, the axis name in each
// metadata line, and group the rule sections under per-axis headings.
func render(dataset schema.Dataset) []byte {
	axes := dataset.EffectiveAxes()
	multiAxis := len(axes) > 1

	var b strings.Builder
	fmt.Fprintf(&b, "<!-- Generated from internal/guidelines/data/%s.json; DO NOT EDIT. -->\n\n", dataset.Language)
	fmt.Fprintf(&b, "# Modern %s Guidelines Explained\n\n", displayName(dataset.Language))
	fmt.Fprintf(&b, "This file provides a more detailed description of the features supported in the Modern %s Guidelines.\n\n", displayName(dataset.Language))
	b.WriteString("**Autofix legend:**\n\n")
	b.WriteString("- [x] — a mechanical autofix exists; the guideline below names its `tool:rule`\n")
	b.WriteString("- [ ] — no autofix available\n\n")
	b.WriteString("**Impact legend:**\n\n")
	b.WriteString("- Critical — found in almost every project, dozens of occurrences\n")
	b.WriteString("- High — found often, 5–20 occurrences per project\n")
	b.WriteString("- Medium — found regularly, 1–5 occurrences per project\n")
	b.WriteString("- Low — found rarely or in specific code\n\n")
	if multiAxis {
		b.WriteString("**Axis legend:**\n\n")
		fmt.Fprintf(
			&b,
			"This dataset gates its guidelines on %d independent version axes: %s. Each guideline belongs to exactly one axis, and the guideline sections below are grouped by axis.\n\n",
			len(axes), strings.Join(axes, ", "),
		)
	}

	b.WriteString("## Guidelines\n\n")
	if multiAxis {
		b.WriteString("| Category | Guideline | Autofix | Axis | Version | Impact |\n")
		b.WriteString("|----------|-----------|------------|------|---------|--------|\n")
	} else {
		fmt.Fprintf(&b, "| Category | Guideline | Autofix | %s | Impact |\n", displayName(axes[0]))
		b.WriteString("|----------|-----------|------------|----|--------|\n")
	}
	groups := rulesByAxis(dataset, axes)
	for axisIndex, group := range groups {
		if multiAxis {
			for _, rule := range group {
				fmt.Fprintf(
					&b,
					"| %s | [`%s`](#%s) | %s | %s | %s | %s |\n",
					rule.Category, rule.ID, rule.ID, autofixMark(rule.Autofix),
					displayName(axes[axisIndex]), rule.SinceVersion, rule.Impact,
				)
			}
		} else {
			for _, rule := range group {
				fmt.Fprintf(
					&b,
					"| %s | [`%s`](#%s) | %s | %s | %s |\n",
					rule.Category, rule.ID, rule.ID, autofixMark(rule.Autofix),
					rule.SinceVersion, rule.Impact,
				)
			}
		}
	}

	for axisIndex, group := range groups {
		if multiAxis {
			b.WriteString("\n---\n\n")
			// No trailing blank here: the rule separator below supplies it.
			fmt.Fprintf(&b, "## Axis: %s\n", displayName(axes[axisIndex]))
		}
		for _, rule := range group {
			b.WriteString("\n---\n\n")
			fmt.Fprintf(&b, "<a id=\"%s\"></a>\n\n", rule.ID)
			fmt.Fprintf(&b, "## `%s`\n\n", rule.ID)
			fmt.Fprintf(
				&b,
				"**Category: %s · %s %s+ · Impact: %s · Autofix: %s**\n\n",
				rule.Category, displayName(ruleAxis(dataset, rule)), rule.SinceVersion,
				rule.Impact, autofixLabel(rule.Autofix),
			)
			b.WriteString(rule.Guideline)
			b.WriteString("\n\n")
			b.WriteString(rule.Details)
			b.WriteByte('\n')

			for i, example := range rule.Examples {
				b.WriteString("\n### Example")
				if len(rule.Examples) > 1 {
					fmt.Fprintf(&b, " %d", i+1)
				}
				b.WriteString("\n\n**Before:**\n\n")
				writeCodeBlock(&b, example.Before, dataset.Language)
				b.WriteString("\n**After:**\n\n")
				writeCodeBlock(&b, example.After, dataset.Language)
			}
		}
	}

	return []byte(b.String())
}

// rulesByAxis returns the rules in rendering order, one group per axis:
// multi-axis datasets group rules under their axis following the declared
// axis order, keeping the newest-first dataset order within each group
// (a schema invariant); single-axis datasets keep plain dataset order.
func rulesByAxis(dataset schema.Dataset, axes []string) [][]schema.Rule {
	groups := make([][]schema.Rule, len(axes))
	for i, axis := range axes {
		for _, rule := range dataset.Rules {
			if ruleAxis(dataset, rule) == axis {
				groups[i] = append(groups[i], rule)
			}
		}
	}
	return groups
}

// ruleAxis is the effective axis of a rule: its explicit axis on multi-axis
// datasets, or the language itself on single-axis datasets.
func ruleAxis(dataset schema.Dataset, rule schema.Rule) string {
	if rule.Axis != "" {
		return rule.Axis
	}
	return dataset.Language
}

// displayNames holds the document display forms of the known language and
// axis identifiers; plain capitalization would render "TypeScript" wrong.
var displayNames = map[string]string{
	"go":         "Go",
	"python":     "Python",
	"typescript": "TypeScript",
	"javascript": "JavaScript",
	"node":       "Node",
}

// displayName returns the display form of a language or axis identifier.
// Schema validation guarantees lowercase ASCII, so the fallback can safely
// capitalize the first byte.
func displayName(name string) string {
	if display, ok := displayNames[name]; ok {
		return display
	}
	if name == "" {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func autofixMark(autofix schema.Autofix) string {
	if autofix.IsNull() {
		return "[ ]"
	}
	return "[x]"
}

func autofixLabel(autofix schema.Autofix) string {
	if autofix.IsNull() {
		return "none"
	}
	return autofix.String()
}

func writeCodeBlock(b *strings.Builder, lines []string, language string) {
	fmt.Fprintf(b, "```%s\n", language)
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString("```\n")
}
