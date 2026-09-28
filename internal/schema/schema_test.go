package schema

import (
	"encoding/json"
	"strings"
	"testing"
)

func validRule() Rule {
	return Rule{
		ID:           "dict_merge",
		SinceVersion: "3.9",
		Autofix:      Autofix{Present: true, Tool: "ruff", Rule: "RUF005"},
		Category:     "Idioms",
		Impact:       "High",
		Guideline:    "Use `dict | other` to merge mappings.",
		Details:      "The union operator merges in one expression and keeps types precise; dict(**a, **b) silently drops duplicate keys.",
		Examples: []Example{{
			Before: []string{"merged = {**first, **second}"},
			After:  []string{"merged = first | second"},
		}},
	}
}

func validDataset() Dataset {
	return Dataset{
		Language: "python",
		Rules:    []Rule{validRule()},
	}
}

func multiAxisDataset() Dataset {
	return Dataset{
		Language: "typescript",
		Axes:     []string{"typescript", "node"},
		Rules: []Rule{
			{
				ID:           "satisfies_operator",
				SinceVersion: "5.4",
				Axis:         "typescript",
				Autofix:      Autofix{Present: true},
				Category:     "Type System",
				Impact:       "Medium",
				Guideline:    "Use `satisfies` to validate a value against a type without widening it.",
				Details:      "satisfies checks assignability while keeping the literal inferred type; a plain annotation widens it.",
				Examples: []Example{{
					Before: []string{"const config: Record<string, string | number> = { retries: 3 }"},
					After:  []string{"const config = { retries: 3 } satisfies Record<string, string | number>"},
				}},
			},
			{
				ID:           "global_fetch",
				SinceVersion: "18",
				Axis:         "node",
				Autofix:      Autofix{Present: true},
				Category:     "Standard Library",
				Impact:       "High",
				Guideline:    "Use the global `fetch` instead of dependency HTTP clients.",
				Details:      "fetch is stable since Node 18; dropping the client removes an install and an API surface.",
				Examples: []Example{{
					Before: []string{"import axios from \"axios\"", "const res = await axios.get(url)"},
					After:  []string{"const res = await fetch(url)"},
				}},
			},
			{
				ID:           "const_type_parameters",
				SinceVersion: "4.9",
				Axis:         "typescript",
				Autofix:      Autofix{Present: true},
				Category:     "Type System",
				Impact:       "Low",
				Guideline:    "Prefer `const` type parameters over `as const` call sites.",
				Details:      "The modifier keeps literal inference at the declaration instead of at every call site.",
				Examples: []Example{{
					Before: []string{"const routes = asConst([\"/a\", \"/b\"] as const)"},
					After:  []string{"const routes = asConstConst([\"/a\", \"/b\"])"},
				}},
			},
			{
				ID:           "structured_clone",
				SinceVersion: "17",
				Axis:         "node",
				Autofix:      Autofix{Present: true},
				Category:     "Standard Library",
				Impact:       "Medium",
				Guideline:    "Use `structuredClone` instead of JSON round-trips for deep copies.",
				Details:      "structuredClone preserves types such as Map, Set, and Dates that JSON.parse(JSON.stringify(x)) silently corrupts.",
				Examples: []Example{{
					Before: []string{"const copy = JSON.parse(JSON.stringify(value))"},
					After:  []string{"const copy = structuredClone(value)"},
				}},
			},
		},
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestParseAcceptsSingleAxisDataset(t *testing.T) {
	dataset := validDataset()
	// Newest-first: the newer rule must appear first.
	dataset.Rules = append([]Rule{{
		ID:           "pep695_type_params",
		SinceVersion: "3.12",
		Autofix:      Autofix{Present: true},
		Category:     "Typing",
		Impact:       "Low",
		Guideline:    "Use PEP 695 type parameter syntax.",
		Details:      "The new syntax removes the TypeVar boilerplate and scopes parameters to the definition.",
		References:   []string{"PEP 695", "typing docs"},
		Examples: []Example{
			{
				Before: []string{"T = TypeVar(\"T\")", "def first(items: list[T]) -> T: ..."},
				After:  []string{"def first[T](items: list[T]) -> T: ..."},
			},
			{
				Before: []string{"class Box(Generic[T]): ..."},
				After:  []string{"class Box[T]: ..."},
			},
		},
	}}, dataset.Rules...)

	parsed, err := Parse(mustJSON(t, dataset))
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Language != "python" || len(parsed.Rules) != 2 {
		t.Fatalf("Parse() = %#v", parsed)
	}
	if axes := parsed.EffectiveAxes(); len(axes) != 1 || axes[0] != "python" {
		t.Fatalf("EffectiveAxes() = %#v, want [python]", axes)
	}
	if parsed.Provenance != nil {
		t.Fatalf("native dataset has provenance: %#v", parsed.Provenance)
	}
	newer := parsed.Rules[0]
	if !newer.Autofix.Present || newer.Autofix.String() != "" {
		t.Fatalf("explicit null autofix = %#v", newer.Autofix)
	}
	if len(newer.References) != 2 || len(newer.Examples) != 2 {
		t.Fatalf("rule = %#v", newer)
	}
	older := parsed.Rules[1]
	if !older.Autofix.Present || older.Autofix.String() != "ruff:RUF005" {
		t.Fatalf("autofix = %#v", older.Autofix)
	}
}

func TestParseAcceptsMultiAxisDataset(t *testing.T) {
	data, err := json.Marshal(multiAxisDataset())
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	axes := parsed.EffectiveAxes()
	if len(axes) != 2 || axes[0] != "typescript" || axes[1] != "node" {
		t.Fatalf("EffectiveAxes() = %#v", axes)
	}
	if len(parsed.Rules) != 4 {
		t.Fatalf("Parse() = %#v", parsed)
	}
	// Node-axis rules carry bare major versions and may interleave with
	// typescript-axis rules; per-axis grouping makes both orders valid.
	if parsed.Rules[1].SinceVersion != "18" || parsed.Rules[1].Axis != "node" {
		t.Fatalf("node rule = %#v", parsed.Rules[1])
	}
}

func TestParseRejectsInvalidRule(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Rule)
		want   string
	}{
		{"id charset", func(r *Rule) { r.ID = "Unsafe-Id" }, "invalid rule id"},
		{"id empty", func(r *Rule) { r.ID = "" }, "invalid rule id"},
		{"version three parts", func(r *Rule) { r.SinceVersion = "3.10.1" }, "invalid since_version"},
		{"version prefix", func(r *Rule) { r.SinceVersion = "v3.10" }, "invalid since_version"},
		{"version word", func(r *Rule) { r.SinceVersion = "latest" }, "invalid since_version"},
		{"version empty", func(r *Rule) { r.SinceVersion = "" }, "invalid since_version"},
		{"version negative", func(r *Rule) { r.SinceVersion = "-3.10" }, "invalid since_version"},
		{"version missing minor", func(r *Rule) { r.SinceVersion = "3." }, "invalid since_version"},
		{"autofix no colon", func(r *Rule) { r.Autofix = Autofix{Present: true, Tool: "ruff-RUF005"} }, "invalid autofix"},
		{"autofix leading colon", func(r *Rule) { r.Autofix = Autofix{Present: true, Tool: "", Rule: "RUF005"} }, "invalid autofix"},
		{"autofix two colons", func(r *Rule) { r.Autofix = Autofix{Present: true, Tool: "ruff", Rule: "U:P:007"} }, "invalid autofix"},
		{"category empty", func(r *Rule) { r.Category = "" }, "has no category"},
		{"impact invalid", func(r *Rule) { r.Impact = "Huge" }, "invalid impact"},
		{"impact empty", func(r *Rule) { r.Impact = "" }, "invalid impact"},
		{"guideline blank", func(r *Rule) { r.Guideline = "   " }, "has no guideline text"},
		{"details blank", func(r *Rule) { r.Details = "  " }, "has no details"},
		{"references over limit", func(r *Rule) { r.References = []string{"a", "b", "c"} }, "references"},
		{"reference empty", func(r *Rule) { r.References = []string{"  "} }, "reference 1 is empty"},
		{"examples missing", func(r *Rule) { r.Examples = nil }, "has no examples"},
		{"before blank", func(r *Rule) { r.Examples[0].Before = []string{"   "} }, "empty before snippet"},
		{"after missing", func(r *Rule) { r.Examples[0].After = nil }, "empty after snippet"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dataset := validDataset()
			test.mutate(&dataset.Rules[0])

			_, err := Parse(mustJSON(t, dataset))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Parse() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestParseRejectsMissingAutofixKey(t *testing.T) {
	var raw map[string]any
	if err := json.Unmarshal(mustJSON(t, validDataset()), &raw); err != nil {
		t.Fatal(err)
	}
	rules, ok := raw["rules"].([]any)
	if !ok || len(rules) != 1 {
		t.Fatalf("unexpected marshaled dataset: %#v", raw)
	}
	delete(rules[0].(map[string]any), "autofix")

	_, err := Parse(mustJSON(t, raw))
	if err == nil || !strings.Contains(err.Error(), "has no autofix key") {
		t.Fatalf("Parse() error = %v, want missing autofix key error", err)
	}
}

func TestParseRejectsNonStringAutofix(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"number", 42},
		{"object", map[string]any{"tool": "ruff"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var raw map[string]any
			if err := json.Unmarshal(mustJSON(t, validDataset()), &raw); err != nil {
				t.Fatal(err)
			}
			raw["rules"].([]any)[0].(map[string]any)["autofix"] = test.value

			_, err := Parse(mustJSON(t, raw))
			if err == nil || !strings.Contains(err.Error(), "autofix must be a JSON string") {
				t.Fatalf("Parse() error = %v, want autofix type error", err)
			}
		})
	}
}

func TestParseRejectsDuplicateIDs(t *testing.T) {
	dataset := validDataset()
	dataset.Rules = append(dataset.Rules, validRule())

	_, err := Parse(mustJSON(t, dataset))
	if err == nil || !strings.Contains(err.Error(), "is duplicated") {
		t.Fatalf("Parse() error = %v, want duplicate id error", err)
	}
}

func TestParseRejectsInvalidDatasetHeader(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Dataset)
		want   string
	}{
		{"language missing", func(d *Dataset) { d.Language = "" }, "has no language"},
		{"language charset", func(d *Dataset) { d.Language = "Python" }, "invalid dataset language"},
		{"rules missing", func(d *Dataset) { d.Rules = nil }, "has no rules"},
		{"duplicate axis", func(d *Dataset) { d.Language = "typescript"; d.Axes = []string{"typescript", "typescript"} }, "duplicate axis"},
		{"axis charset", func(d *Dataset) { d.Axes = []string{"Node"} }, "invalid axis"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dataset := validDataset()
			test.mutate(&dataset)

			_, err := Parse(mustJSON(t, dataset))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Parse() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestParseRejectsIncompleteProvenance(t *testing.T) {
	tests := []struct {
		name       string
		provenance Provenance
		want       string
	}{
		{"no source", Provenance{License: "Apache-2.0", Commit: "abc", Date: "2026-09-29"}, `"source"`},
		{"no license", Provenance{Source: "JetBrains/go-modern-guidelines", Commit: "abc", Date: "2026-09-29"}, `"license"`},
		{"no commit", Provenance{Source: "JetBrains/go-modern-guidelines", License: "Apache-2.0", Date: "2026-09-29"}, `"commit"`},
		{"no date", Provenance{Source: "JetBrains/go-modern-guidelines", License: "Apache-2.0", Commit: "abc"}, `"date"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dataset := validDataset()
			provenance := test.provenance
			dataset.Provenance = &provenance

			_, err := Parse(mustJSON(t, dataset))
			if err == nil || !strings.Contains(err.Error(), "provenance is missing") || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Parse() error = %v, want provenance is missing %s", err, test.want)
			}
		})
	}
}

func TestParseExternalRequiresProvenance(t *testing.T) {
	_, err := ParseExternal(mustJSON(t, validDataset()))
	if err == nil || !strings.Contains(err.Error(), "externally absorbed but has no provenance") {
		t.Fatalf("ParseExternal() error = %v, want missing provenance error", err)
	}

	dataset := validDataset()
	dataset.Provenance = &Provenance{
		Source:  "JetBrains/go-modern-guidelines",
		License: "Apache-2.0",
		Commit:  "019b45e2b2f80d7f5c1e28bd4d35f3f0fbcf9c9c",
		Date:    "2026-09-29",
	}
	parsed, err := ParseExternal(mustJSON(t, dataset))
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Provenance.Source != "JetBrains/go-modern-guidelines" {
		t.Fatalf("provenance = %#v", parsed.Provenance)
	}
}

func TestParseRejectsMultiAxisViolations(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Dataset)
		want   string
	}{
		{"axis missing", func(d *Dataset) { d.Rules[0].Axis = "" }, "has no axis"},
		{"axis undeclared", func(d *Dataset) { d.Rules[0].Axis = "browser" }, "not declared in dataset axes"},
		{"newest-first violated within axis", func(d *Dataset) {
			d.Rules = append(d.Rules, Rule{
				ID: "satisfies_late", SinceVersion: "5.5", Axis: "typescript", Autofix: Autofix{Present: true},
				Category: "Type System", Impact: "Low", Guideline: "g", Details: "d",
				Examples: []Example{{Before: []string{"a"}, After: []string{"b"}}},
			})
		}, "ordered newest-first"},
		{"axis set on single-axis dataset", func(d *Dataset) {
			single := validDataset()
			single.Rules[0].Axis = "python"
			*d = single
		}, "must not set axis"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dataset := multiAxisDataset()
			test.mutate(&dataset)

			_, err := Parse(mustJSON(t, dataset))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Parse() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestParseRejectsNewestFirstOnSingleAxis(t *testing.T) {
	dataset := validDataset()
	dataset.Rules = append(dataset.Rules, Rule{
		ID: "newer_rule", SinceVersion: "3.10", Autofix: Autofix{Present: true},
		Category: "Idioms", Impact: "Low", Guideline: "g", Details: "d",
		Examples: []Example{{Before: []string{"a"}, After: []string{"b"}}},
	})

	_, err := Parse(mustJSON(t, dataset))
	if err == nil || !strings.Contains(err.Error(), "ordered newest-first") {
		t.Fatalf("Parse() error = %v, want ordering error", err)
	}
}

func TestParseRejectsInvalidJSON(t *testing.T) {
	_, err := Parse([]byte("{"))
	if err == nil || !strings.Contains(err.Error(), "parse dataset JSON") {
		t.Fatalf("Parse() error = %v, want JSON error", err)
	}
}

func TestParseRejectsUnknownField(t *testing.T) {
	var raw map[string]any
	if err := json.Unmarshal(mustJSON(t, validDataset()), &raw); err != nil {
		t.Fatal(err)
	}
	raw["unknown_field"] = true

	_, err := Parse(mustJSON(t, raw))
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("Parse() error = %v, want unknown field error", err)
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		left  string
		right string
		want  int
	}{
		{"18", "18", 0},
		{"18", "18.0", 0},
		{"5.4", "5.10", -1},
		{"3.9", "3.10", -1},
		{"1.16", "1.9", 1},
		{"18", "17.9", 1},
		{"5.4", "18", -1},
	}
	for _, test := range tests {
		t.Run(test.left+"_vs_"+test.right, func(t *testing.T) {
			got, err := CompareVersions(test.left, test.right)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("CompareVersions(%q, %q) = %d, want %d", test.left, test.right, got, test.want)
			}
		})
	}

	invalid := []struct{ left, right string }{
		{"latest", "1.0"},
		{"1.0", "v5"},
		{"1.2.3", "1.0"},
		{"", "1.0"},
	}
	for _, test := range invalid {
		t.Run("invalid_"+test.left, func(t *testing.T) {
			if _, err := CompareVersions(test.left, test.right); err == nil || !strings.Contains(err.Error(), "invalid version") {
				t.Fatalf("CompareVersions(%q, %q) error = %v, want invalid version error", test.left, test.right, err)
			}
		})
	}
}
