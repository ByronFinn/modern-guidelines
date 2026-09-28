// Package guidelines embeds every guideline dataset under data/, validates
// each one through internal/schema at init time (invalid data panics the
// process, mirroring the reference project), registers them in
// internal/registry, and renders the list and explain output
// (PRD-0000 requirements 1, 7, and 8).
//
// Rendering semantics follow the reference project
// examples/go-modern-guidelines/internal/guidelines: list output is one
// "id: guideline" line per applicable rule, newest-first; explain output is
// an indented detail block per requested id, with the schema v2 additions
// References (rendered only when present) and Autofix (always rendered).
// Multi-axis datasets additionally group list output by axis.
//
//go:generate go run ../featuresgen data ../../docs/features
package guidelines

import (
	"embed"
	"fmt"
	"strings"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

//go:embed data/*.json
var dataFS embed.FS

func init() {
	// ReadDir returns entries sorted by file name, so registration (and any
	// startup panic for invalid data) happens in a deterministic order.
	entries, err := dataFS.ReadDir("data")
	if err != nil {
		panic(fmt.Sprintf("guidelines: read embedded dataset directory: %v", err))
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := dataFS.ReadFile("data/" + entry.Name())
		if err != nil {
			panic(fmt.Sprintf("guidelines: read embedded dataset %s: %v", entry.Name(), err))
		}
		dataset, err := schema.Parse(data)
		if err != nil {
			panic(fmt.Sprintf("guidelines: load embedded dataset %s: %v", entry.Name(), err))
		}
		registry.RegisterDataset(newLoadedDataset(dataset))
	}
}

// loadedDataset adapts a validated schema.Dataset to registry.Dataset and
// keeps an O(1) id index for explain lookups.
type loadedDataset struct {
	dataset   schema.Dataset
	rulesByID map[string]schema.Rule
}

func newLoadedDataset(dataset schema.Dataset) *loadedDataset {
	rulesByID := make(map[string]schema.Rule, len(dataset.Rules))
	for _, rule := range dataset.Rules {
		rulesByID[rule.ID] = rule
	}
	return &loadedDataset{dataset: dataset, rulesByID: rulesByID}
}

func (d *loadedDataset) Language() string               { return d.dataset.Language }
func (d *loadedDataset) Axes() []string                 { return d.dataset.EffectiveAxes() }
func (d *loadedDataset) Provenance() *schema.Provenance { return d.dataset.Provenance }
func (d *loadedDataset) Rules() []schema.Rule           { return d.dataset.Rules }

func (d *loadedDataset) Rule(id string) (schema.Rule, bool) {
	rule, ok := d.rulesByID[id]
	return rule, ok
}

// RulesForVersions returns the rules of dataset that apply when every axis is
// at its resolved version (resolved[axis] >= since_version), in list-output
// order: multi-axis datasets group rules by axis following the declared axis
// order, and within each group rules stay newest-first (dataset order, a
// schema validation invariant). It errors when an axis has no resolved
// version: rules must never render ungated.
func RulesForVersions(dataset registry.Dataset, resolved map[string]string) ([]schema.Rule, error) {
	axes := dataset.Axes()
	for _, axis := range axes {
		if _, ok := resolved[axis]; !ok {
			return nil, fmt.Errorf("no resolved version for axis %q of language %q", axis, dataset.Language())
		}
	}
	applicable := make([]schema.Rule, 0, len(dataset.Rules()))
	for _, axis := range axes {
		for _, rule := range dataset.Rules() {
			if ruleAxis(dataset, rule) != axis {
				continue
			}
			comparison, err := schema.CompareVersions(resolved[axis], rule.SinceVersion)
			if err != nil {
				return nil, fmt.Errorf("rule %q: %w", rule.ID, err)
			}
			if comparison >= 0 {
				applicable = append(applicable, rule)
			}
		}
	}
	return applicable, nil
}

// ListText renders the short list output for dataset at the resolved
// per-axis versions: one "id: guideline" line per applicable rule. Multi-axis
// datasets get one "# axis: <name>" header per axis group.
func ListText(dataset registry.Dataset, resolved map[string]string) (string, error) {
	rules, err := RulesForVersions(dataset, resolved)
	if err != nil {
		return "", err
	}
	multiAxis := len(dataset.Axes()) > 1
	currentAxis := ""
	var b strings.Builder
	for _, rule := range rules {
		if multiAxis {
			if axis := ruleAxis(dataset, rule); axis != currentAxis {
				currentAxis = axis
				fmt.Fprintf(&b, "# axis: %s\n", axis)
			}
		}
		b.WriteString(rule.ID)
		b.WriteString(": ")
		b.WriteString(rule.Guideline)
		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}

// ExplainText renders the detailed guidance blocks for the requested ids of
// one dataset. Unknown ids produce a single error listing every unknown id
// together with the available ids of the language.
func ExplainText(dataset registry.Dataset, ids []string) (string, error) {
	requestedIDs := NormalizeIDs(ids)
	if len(requestedIDs) == 0 {
		return "", fmt.Errorf("explain command requires at least one guideline id. Run list to list available ids")
	}

	var unknownIDs []string
	for _, id := range requestedIDs {
		if _, ok := dataset.Rule(id); !ok {
			unknownIDs = append(unknownIDs, id)
		}
	}
	if len(unknownIDs) > 0 {
		return "", fmt.Errorf(
			"unknown %s guideline ids: %s. Available ids: %s",
			dataset.Language(), strings.Join(unknownIDs, ", "), strings.Join(availableIDs(dataset), ", "),
		)
	}

	blocks := make([]string, 0, len(requestedIDs))
	for _, id := range requestedIDs {
		rule, _ := dataset.Rule(id)
		blocks = append(blocks, ruleBlock(dataset, rule))
	}
	return strings.Join(blocks, "\n\n"), nil
}

// NormalizeIDs trims, comma-splits, drops empty parts, and dedupes guideline
// id references while preserving first-appearance order.
func NormalizeIDs(values []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, value := range values {
		for part := range strings.SplitSeq(value, ",") {
			part = strings.TrimSpace(part)
			if part == "" || seen[part] {
				continue
			}
			seen[part] = true
			result = append(result, part)
		}
	}
	return result
}

func ruleBlock(dataset registry.Dataset, rule schema.Rule) string {
	var b strings.Builder
	b.WriteString(rule.ID)
	b.WriteString(":\n")
	fmt.Fprintf(&b, "  Since: %s %s\n\n", ruleAxis(dataset, rule), rule.SinceVersion)
	b.WriteString("  Summary:\n")
	b.WriteString(indentLines(rule.Guideline, "    "))
	b.WriteString("\n\n  Details:\n")
	b.WriteString(indentLines(rule.Details, "    "))
	writeExamples(&b, rule.Examples)
	writeReferences(&b, rule.References)
	writeAutofix(&b, rule.Autofix)
	return b.String()
}

func writeExamples(b *strings.Builder, examples []schema.Example) {
	if len(examples) == 0 {
		return
	}
	b.WriteString("\n\n  Examples:")
	for index, example := range examples {
		if len(examples) > 1 {
			fmt.Fprintf(b, "\n\n  Example %d:", index+1)
		}
		b.WriteString("\n\n  Before:\n")
		b.WriteString(indentLines(strings.Join(example.Before, "\n"), "    "))
		b.WriteString("\n\n  After:\n")
		b.WriteString(indentLines(strings.Join(example.After, "\n"), "    "))
	}
}

// writeReferences renders the optional References section; it is omitted
// entirely when the rule has none (absorbed Go data carries no references).
func writeReferences(b *strings.Builder, references []string) {
	if len(references) == 0 {
		return
	}
	b.WriteString("\n\n  References:")
	for _, reference := range references {
		b.WriteString("\n")
		b.WriteString(indentLines(reference, "    "))
	}
}

// writeAutofix renders the mandatory Autofix section: the "tool:rule" value,
// or "none" for an explicit null.
func writeAutofix(b *strings.Builder, autofix schema.Autofix) {
	b.WriteString("\n\n  Autofix:\n")
	if autofix.IsNull() {
		b.WriteString("    none")
		return
	}
	b.WriteString(indentLines(autofix.String(), "    "))
}

func indentLines(s, indent string) string {
	var b strings.Builder
	first := true
	for line := range strings.SplitSeq(s, "\n") {
		if !first {
			b.WriteByte('\n')
		}
		first = false
		b.WriteString(indent)
		b.WriteString(line)
	}
	return b.String()
}

// ruleAxis is the effective axis of a rule: its explicit axis on multi-axis
// datasets, or the language itself on single-axis datasets (where the axis
// field must be empty per schema validation).
func ruleAxis(dataset registry.Dataset, rule schema.Rule) string {
	if rule.Axis != "" {
		return rule.Axis
	}
	return dataset.Language()
}

func availableIDs(dataset registry.Dataset) []string {
	rules := dataset.Rules()
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	return ids
}
