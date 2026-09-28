package guidelines

import (
	"slices"
	"strings"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// testRule builds a minimal valid rule; axis stays empty for single-axis
// datasets.
func testRule(id, sinceVersion, axis string) schema.Rule {
	return schema.Rule{
		ID:           id,
		SinceVersion: sinceVersion,
		Axis:         axis,
		Autofix:      schema.Autofix{Present: true},
		Category:     "Idioms",
		Impact:       "Medium",
		Guideline:    "Use `" + id + "`.",
		Details:      "Details for " + id + ".",
		Examples:     []schema.Example{{Before: []string{"before " + id}, After: []string{"after " + id}}},
	}
}

func testDataset(language string, axes []string, rules ...schema.Rule) registry.Dataset {
	return newLoadedDataset(schema.Dataset{Language: language, Axes: axes, Rules: rules})
}

func TestEmbeddedDatasetsRegistered(t *testing.T) {
	// Languages and axes are PRD-level facts (D3, D12), stable across the
	// content PRs that will overwrite the sample rule sets.
	wantAxes := map[string][]string{
		"python":     {"python"},
		"go":         {"go"},
		"typescript": {"typescript", "node"},
	}
	for language, axes := range wantAxes {
		t.Run(language, func(t *testing.T) {
			dataset, ok := registry.DatasetForLanguage(language)
			if !ok {
				t.Fatalf("dataset for %q is not registered", language)
			}
			if dataset.Language() != language {
				t.Fatalf("Language() = %q, want %q", dataset.Language(), language)
			}
			if got := dataset.Axes(); !slices.Equal(got, axes) {
				t.Fatalf("Axes() = %v, want %v", got, axes)
			}
			if len(dataset.Rules()) == 0 {
				t.Fatal("registered dataset has no rules")
			}
			for _, rule := range dataset.Rules() {
				indexed, ok := dataset.Rule(rule.ID)
				if !ok || indexed.ID != rule.ID {
					t.Fatalf("O(1) index missing rule %q", rule.ID)
				}
			}
		})
	}
}

func TestLoadedDatasetRuleIndexIntegrity(t *testing.T) {
	dataset := newLoadedDataset(schema.Dataset{
		Language: "python",
		Rules:    []schema.Rule{testRule("a_rule", "3.12", ""), testRule("b_rule", "3.10", "")},
	})
	if len(dataset.rulesByID) != len(dataset.Rules()) {
		t.Fatalf("index has %d entries, want %d", len(dataset.rulesByID), len(dataset.Rules()))
	}
	if _, ok := dataset.Rule("missing_rule"); ok {
		t.Fatal("Rule(missing_rule) reported found")
	}
}

func TestRulesForVersionsFiltersPerAxis(t *testing.T) {
	dataset := testDataset("python", nil,
		testRule("new_rule", "3.12", ""),
		testRule("old_rule", "3.10", ""),
	)

	tests := []struct {
		name    string
		version string
		wantIDs []string
	}{
		{"newer than all", "3.14", []string{"new_rule", "old_rule"}},
		{"boundary includes equal", "3.12", []string{"new_rule", "old_rule"}},
		{"between tiers", "3.11", []string{"old_rule"}},
		{"older than all", "3.9", nil},
		{"bare major normalizes", "4", []string{"new_rule", "old_rule"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules, err := RulesForVersions(dataset, map[string]string{"python": test.version})
			if err != nil {
				t.Fatal(err)
			}
			gotIDs := make([]string, 0, len(rules))
			for _, rule := range rules {
				gotIDs = append(gotIDs, rule.ID)
			}
			if !slices.Equal(gotIDs, test.wantIDs) {
				t.Fatalf("RulesForVersions(python %s) = %v, want %v", test.version, gotIDs, test.wantIDs)
			}
		})
	}
}

func TestRulesForVersionsMultiAxis(t *testing.T) {
	// Deliberately interleaved rule order: valid per schema (newest-first
	// holds within each axis) and the grouping must regroup by axis.
	dataset := testDataset("typescript", []string{"typescript", "node"},
		testRule("node_fetch", "18", "node"),
		testRule("ts_satisfies", "4.9", "typescript"),
		testRule("node_clone", "17", "node"),
	)

	tests := []struct {
		name     string
		resolved map[string]string
		wantIDs  []string
	}{
		{
			name:     "grouped by axis order",
			resolved: map[string]string{"typescript": "5.0", "node": "20"},
			wantIDs:  []string{"ts_satisfies", "node_fetch", "node_clone"},
		},
		{
			name:     "per-axis gating",
			resolved: map[string]string{"typescript": "5.0", "node": "17"},
			wantIDs:  []string{"ts_satisfies", "node_clone"},
		},
		{
			name:     "typescript axis below its rule",
			resolved: map[string]string{"typescript": "4.8", "node": "18"},
			wantIDs:  []string{"node_fetch", "node_clone"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules, err := RulesForVersions(dataset, test.resolved)
			if err != nil {
				t.Fatal(err)
			}
			gotIDs := make([]string, 0, len(rules))
			for _, rule := range rules {
				gotIDs = append(gotIDs, rule.ID)
			}
			if !slices.Equal(gotIDs, test.wantIDs) {
				t.Fatalf("RulesForVersions(%v) = %v, want %v", test.resolved, gotIDs, test.wantIDs)
			}
		})
	}
}

func TestRulesForVersionsRequiresEveryAxis(t *testing.T) {
	dataset := testDataset("typescript", []string{"typescript", "node"},
		testRule("ts_satisfies", "4.9", "typescript"),
		testRule("node_fetch", "18", "node"),
	)

	_, err := RulesForVersions(dataset, map[string]string{"typescript": "5.0"})
	if err == nil || !strings.Contains(err.Error(), `no resolved version for axis "node"`) {
		t.Fatalf("RulesForVersions() error = %v, want missing-axis error", err)
	}
}

func TestListTextSingleAxis(t *testing.T) {
	dataset := testDataset("python", nil,
		testRule("new_rule", "3.12", ""),
		testRule("old_rule", "3.10", ""),
	)

	got, err := ListText(dataset, map[string]string{"python": "3.12"})
	if err != nil {
		t.Fatal(err)
	}
	want := "new_rule: Use `new_rule`.\nold_rule: Use `old_rule`."
	if got != want {
		t.Fatalf("ListText() =\n%q\nwant\n%q", got, want)
	}
}

func TestListTextMultiAxisGroupsByAxis(t *testing.T) {
	dataset := testDataset("typescript", []string{"typescript", "node"},
		testRule("node_fetch", "18", "node"),
		testRule("ts_satisfies", "4.9", "typescript"),
		testRule("node_clone", "17", "node"),
	)

	got, err := ListText(dataset, map[string]string{"typescript": "5.0", "node": "20"})
	if err != nil {
		t.Fatal(err)
	}
	want := "# axis: typescript\nts_satisfies: Use `ts_satisfies`.\n" +
		"# axis: node\nnode_fetch: Use `node_fetch`.\nnode_clone: Use `node_clone`."
	if got != want {
		t.Fatalf("ListText() =\n%q\nwant\n%q", got, want)
	}
}

func TestListTextPropagatesFilterErrors(t *testing.T) {
	dataset := testDataset("python", nil, testRule("a_rule", "3.12", ""))

	if _, err := ListText(dataset, map[string]string{"go": "1.22"}); err == nil || !strings.Contains(err.Error(), `axis "python"`) {
		t.Fatalf("ListText() error = %v, want missing-axis error", err)
	}
}

func TestExplainTextFormat(t *testing.T) {
	rule := schema.Rule{
		ID:           "union_type_syntax",
		SinceVersion: "3.10",
		Autofix:      schema.Autofix{Present: true, Tool: "ruff", Rule: "UP007"},
		Category:     "Typing",
		Impact:       "High",
		Guideline:    "Use `X | Y` instead of `typing.Union[X, Y]`.",
		Details:      "PEP 604 union operator.\nSecond line of details.",
		References:   []string{"PEP 604", "typing docs"},
		Examples: []schema.Example{
			{Before: []string{"x: Union[int, str]"}, After: []string{"x: int | str"}},
			{Before: []string{"y: Union[a, b]"}, After: []string{"y: a | b"}},
		},
	}
	dataset := testDataset("python", nil, rule)

	got, err := ExplainText(dataset, []string{"union_type_syntax"})
	if err != nil {
		t.Fatal(err)
	}
	want := `union_type_syntax:
  Since: python 3.10

  Summary:
    Use ` + "`X | Y` instead of `typing.Union[X, Y]`." + `

  Details:
    PEP 604 union operator.
    Second line of details.

  Examples:

  Example 1:

  Before:
    x: Union[int, str]

  After:
    x: int | str

  Example 2:

  Before:
    y: Union[a, b]

  After:
    y: a | b

  References:
    PEP 604
    typing docs

  Autofix:
    ruff:UP007`
	if got != want {
		t.Fatalf("ExplainText() =\n%s\nwant\n%s", got, want)
	}
}

func TestExplainTextOmitsReferencesAndRendersNullAutofix(t *testing.T) {
	dataset := testDataset("python", nil, schema.Rule{
		ID:           "old_style",
		SinceVersion: "3.9",
		Autofix:      schema.Autofix{Present: true},
		Category:     "Idioms",
		Impact:       "Low",
		Guideline:    "Summary text.",
		Details:      "Details text.",
		Examples:     []schema.Example{{Before: []string{"before"}, After: []string{"after"}}},
	})

	got, err := ExplainText(dataset, []string{"old_style"})
	if err != nil {
		t.Fatal(err)
	}
	want := `old_style:
  Since: python 3.9

  Summary:
    Summary text.

  Details:
    Details text.

  Examples:

  Before:
    before

  After:
    after

  Autofix:
    none`
	if got != want {
		t.Fatalf("ExplainText() =\n%s\nwant\n%s", got, want)
	}
	if strings.Contains(got, "References") {
		t.Fatalf("References section rendered for a rule without references:\n%s", got)
	}
}

func TestExplainTextUsesRuleAxisInSince(t *testing.T) {
	dataset := testDataset("typescript", []string{"typescript", "node"},
		testRule("node_fetch", "18", "node"),
		testRule("ts_satisfies", "4.9", "typescript"),
	)

	tests := []struct {
		id      string
		wantSin string
	}{
		{"node_fetch", "  Since: node 18"},
		{"ts_satisfies", "  Since: typescript 4.9"},
	}
	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			got, err := ExplainText(dataset, []string{test.id})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(got, test.wantSin) {
				t.Fatalf("ExplainText(%q) missing %q:\n%s", test.id, test.wantSin, got)
			}
		})
	}
}

func TestExplainTextMultipleIDs(t *testing.T) {
	dataset := testDataset("python", nil,
		testRule("new_rule", "3.12", ""),
		testRule("old_rule", "3.10", ""),
	)

	got, err := ExplainText(dataset, []string{"old_rule", "new_rule", "old_rule", " new_rule "})
	if err != nil {
		t.Fatal(err)
	}
	newIndex := strings.Index(got, "new_rule:\n")
	oldIndex := strings.Index(got, "old_rule:\n")
	if newIndex == -1 || oldIndex == -1 {
		t.Fatalf("ExplainText() missing a block:\n%s", got)
	}
	// Requested order wins; duplicates were dropped.
	if oldIndex > newIndex {
		t.Fatalf("ExplainText() blocks not in requested order:\n%s", got)
	}
	if strings.Count(got, "Since:") != 2 {
		t.Fatalf("ExplainText() rendered %d blocks, want 2:\n%s", strings.Count(got, "Since:"), got)
	}
	if !strings.Contains(got, "\n\n") {
		t.Fatalf("ExplainText() blocks not separated by a blank line:\n%s", got)
	}
}

func TestExplainTextUnknownIDs(t *testing.T) {
	dataset := testDataset("python", nil,
		testRule("new_rule", "3.12", ""),
		testRule("old_rule", "3.10", ""),
	)

	_, err := ExplainText(dataset, []string{"nope", "also_nope", "old_rule"})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{
		"unknown python guideline ids: nope, also_nope",
		"Available ids: new_rule, old_rule",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func TestExplainTextRequiresIDs(t *testing.T) {
	dataset := testDataset("python", nil, testRule("a_rule", "3.12", ""))

	tests := []struct {
		name string
		ids  []string
	}{
		{"none", nil},
		{"blank only", []string{"  ", ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ExplainText(dataset, test.ids)
			if err == nil || !strings.Contains(err.Error(), "requires at least one guideline id") {
				t.Fatalf("ExplainText(%v) error = %v, want requires-ids error", test.ids, err)
			}
		})
	}
}

func TestNormalizeIDs(t *testing.T) {
	got := NormalizeIDs([]string{" a , b ", "a", "", "c"})
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Fatalf("NormalizeIDs() = %v, want %v", got, want)
	}
}
