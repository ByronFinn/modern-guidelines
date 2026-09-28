package cli

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

var errWriteFailed = errors.New("write failed")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errWriteFailed
}

// datasetFor returns the embedded dataset of a language. Tests derive ids,
// versions, and guideline lines from the dataset at runtime instead of
// hardcoding sample content: the content PRs overwrite the sample data
// entirely.
func datasetFor(t *testing.T, language string) registry.Dataset {
	t.Helper()
	dataset, ok := registry.DatasetForLanguage(language)
	if !ok {
		t.Fatalf("dataset for %q is not registered", language)
	}
	return dataset
}

// listLines returns the guideline lines of a list output, without the
// provenance and axis header lines.
func listLines(t *testing.T, output string) []string {
	t.Helper()
	var lines []string
	for line := range strings.SplitSeq(output, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func TestListWithVersionFlagIncludesAllRulesAtNewestVersion(t *testing.T) {
	dataset := datasetFor(t, "python")
	newest := dataset.Rules()[0].SinceVersion // newest-first is a schema invariant

	var stdout bytes.Buffer
	if err := Run([]string{"list", "--version", "python=" + newest}, &stdout); err != nil {
		t.Fatal(err)
	}

	wantProvenance := "# python " + newest + " (source: --version flag, fidelity: high)"
	output := stdout.String()
	if !strings.HasPrefix(output, wantProvenance+"\n") {
		t.Fatalf("list output missing provenance first line %q:\n%s", wantProvenance, output)
	}

	// Every rule applies at the dataset's newest version: expect one line per
	// rule, in dataset order.
	var wantLines []string
	for _, rule := range dataset.Rules() {
		wantLines = append(wantLines, rule.ID+": "+rule.Guideline)
	}
	if got := listLines(t, output); !slices.Equal(got, wantLines) {
		t.Fatalf("list lines = %v\nwant %v\noutput:\n%s", got, wantLines, output)
	}
}

func TestListWithVersionFlagFiltersOlderRules(t *testing.T) {
	dataset := datasetFor(t, "python")
	rules := dataset.Rules()
	newest, oldest := rules[0].SinceVersion, rules[len(rules)-1].SinceVersion
	newer, err := schema.CompareVersions(newest, oldest)
	if err != nil {
		t.Fatal(err)
	}
	if newer == 0 {
		t.Skip("sample dataset has a single version tier")
	}

	var stdout bytes.Buffer
	if err := Run([]string{"list", "--version", "python=" + oldest}, &stdout); err != nil {
		t.Fatal(err)
	}

	var wantLines []string
	for _, rule := range rules {
		comparison, err := schema.CompareVersions(oldest, rule.SinceVersion)
		if err != nil {
			t.Fatal(err)
		}
		if comparison >= 0 {
			wantLines = append(wantLines, rule.ID+": "+rule.Guideline)
		}
	}
	got := listLines(t, stdout.String())
	if !slices.Equal(got, wantLines) {
		t.Fatalf("list lines = %v\nwant %v", got, wantLines)
	}
	if slices.Contains(got, rules[0].ID+": "+rules[0].Guideline) {
		t.Fatalf("newest rule still listed at version %s:\n%s", oldest, stdout.String())
	}
}

func TestListMultiAxisDataset(t *testing.T) {
	dataset := datasetFor(t, "typescript")
	newest := make(map[string]string)
	for _, rule := range dataset.Rules() {
		axis := rule.Axis
		if axis == "" {
			axis = dataset.Language()
		}
		if _, ok := newest[axis]; !ok {
			newest[axis] = rule.SinceVersion // newest-first per axis
		}
	}
	versionFlag := "typescript=" + newest["typescript"] + ",node=" + newest["node"]

	var stdout bytes.Buffer
	if err := Run([]string{"list", "--version", versionFlag}, &stdout); err != nil {
		t.Fatal(err)
	}

	output := stdout.String()
	for _, want := range []string{
		"# typescript " + newest["typescript"] + " (source: --version flag, fidelity: high)",
		"# node " + newest["node"] + " (source: --version flag, fidelity: high)",
		"# axis: typescript",
		"# axis: node",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("list output missing %q:\n%s", want, output)
		}
	}
	// Grouped per axis: the node group must start after the typescript group.
	if strings.Index(output, "# axis: typescript") > strings.Index(output, "# axis: node") {
		t.Fatalf("axis groups out of order:\n%s", output)
	}
}

func TestListVersionSpecErrors(t *testing.T) {
	tests := []struct {
		name string
		spec string
		want string
	}{
		{"bare version", "3.12", "use axis=version"},
		{"missing value", "python=", "use axis=version"},
		{"missing axis", "=3.12", "use axis=version"},
		{"malformed version", "python=3.x", `invalid --version entry "python=3.x"`},
		{"three-part version", "python=3.12.1", "invalid --version entry"},
		{"duplicate axis", "python=3.12,python=3.13", `duplicate axis "python"`},
		{"entries only", ",", "at least one axis=version"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run([]string{"list", "--version", test.spec}, &stdout)
			if err == nil {
				t.Fatalf("list --version %q: expected error", test.spec)
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("list --version %q error = %v, want substring %q", test.spec, err, test.want)
			}
		})
	}
}

func TestListUnknownVersionAxis(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--version", "ruby=3.3"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	for _, want := range []string{"unknown version axis", "known axes:", "node", "python", "typescript"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func TestListAxesSpanningMultipleLanguages(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--version", "python=3.12,node=18"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(err.Error(), "exactly one language") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListMissingAxisOfMultiAxisDataset(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--version", "typescript=5.4"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	for _, want := range []string{"missing the node axis", `"typescript"`, "--version typescript=5.4,node=<version>"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func TestListRejectsConflictingVersionSources(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"version flag with file-path flag", []string{"list", "--version", "python=3.12", "--file-path", "one.py"}},
		{"version flag with positional path", []string{"list", "--version", "python=3.12", "one.py"}},
		{"version flag with two positional paths", []string{"list", "--version", "python=3.12", "one.py", "two.py"}},
		{"file-path flag with positional path", []string{"list", "--file-path", "one.py", "two.py"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(test.args, &stdout)
			if err == nil {
				t.Fatal("expected error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(err.Error(), "accepts only one version source") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestListRejectsMultiplePositionalPaths(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "one.py", "two.py"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(err.Error(), "at most one positional file path") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRequiresVersionSource(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"bare list", []string{"list"}},
		{"lang without source", []string{"list", "--lang", "python"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(test.args, &stdout)
			if err == nil {
				t.Fatal("expected error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(err.Error(), "requires a version source") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestListRejectsLangWithVersionFlag(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--lang", "python", "--version", "python=3.12"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(err.Error(), "--lang cannot be combined with --version") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRejectsUnknownExtension(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"flag form", []string{"list", "--file-path", "script.rb"}},
		{"positional form", []string{"list", "script.rb"}},
		{"no extension", []string{"list", "Makefile"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(test.args, &stdout)
			if err == nil {
				t.Fatal("expected error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(err.Error(), "cannot identify the language") {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(err.Error(), "Supported languages: go, python, typescript") {
				t.Fatalf("error does not list supported languages: %v", err)
			}
		})
	}
}

func TestListRejectsLanguageWithoutDataset(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--file-path", "app.js"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	for _, want := range []string{`no guideline dataset is installed for language "javascript"`, "Installed languages:"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func TestListFilePathWithoutDetectorGuidesToVersionFlag(t *testing.T) {
	if _, ok := registry.DetectorForLanguage("python"); ok {
		t.Skip("a python detector is registered; manifest resolution applies")
	}

	var stdout bytes.Buffer

	err := Run([]string{"list", "--file-path", "src/app.py"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	for _, want := range []string{`no version detector is available for language "python"`, "--version python=<version>"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func TestLanguagesListsDatasetsAndAxes(t *testing.T) {
	var stdout bytes.Buffer

	if err := Run([]string{"languages"}, &stdout); err != nil {
		t.Fatal(err)
	}
	want := "go (axes: go)\npython (axes: python)\ntypescript (axes: typescript, node)\n"
	if stdout.String() != want {
		t.Fatalf("languages output = %q, want %q", stdout.String(), want)
	}
}

func TestExplainFullChain(t *testing.T) {
	dataset := datasetFor(t, "python")
	rule := dataset.Rules()[0]

	var stdout bytes.Buffer
	if err := Run([]string{"explain", "python:" + rule.ID}, &stdout); err != nil {
		t.Fatal(err)
	}

	output := stdout.String()
	if !strings.HasPrefix(output, rule.ID+":\n") {
		t.Fatalf("explain output does not start with the id header:\n%s", output)
	}
	for _, want := range []string{
		"  Since: python " + rule.SinceVersion,
		"  Summary:",
		"  Details:",
		"  Examples:",
		"  Autofix:",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("explain output missing %q:\n%s", want, output)
		}
	}
	if len(rule.References) > 0 && !strings.Contains(output, "  References:") {
		t.Fatalf("explain output missing References section for a rule with references:\n%s", output)
	}
	if len(rule.References) == 0 && strings.Contains(output, "References") {
		t.Fatalf("explain output rendered References for a rule without references:\n%s", output)
	}
}

func TestExplainWithLangFlag(t *testing.T) {
	dataset := datasetFor(t, "go")
	rule := dataset.Rules()[0]

	var stdout bytes.Buffer
	if err := Run([]string{"explain", "--lang", "go", rule.ID}, &stdout); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(stdout.String(), rule.ID+":\n") {
		t.Fatalf("explain --lang output:\n%s", stdout.String())
	}
}

func TestExplainMixedLanguages(t *testing.T) {
	pythonDataset := datasetFor(t, "python")
	goDataset := datasetFor(t, "go")
	pythonID := pythonDataset.Rules()[0].ID
	goID := goDataset.Rules()[0].ID

	var stdout bytes.Buffer
	if err := Run([]string{"explain", "python:" + pythonID, "go:" + goID}, &stdout); err != nil {
		t.Fatal(err)
	}

	output := stdout.String()
	pythonIndex := strings.Index(output, "  Since: python ")
	goIndex := strings.Index(output, "  Since: go ")
	if pythonIndex == -1 || goIndex == -1 {
		t.Fatalf("explain output missing a language block:\n%s", output)
	}
	if pythonIndex > goIndex {
		t.Fatalf("explain blocks not in requested order:\n%s", output)
	}
}

func TestExplainUnknownIDListsAvailableIDs(t *testing.T) {
	dataset := datasetFor(t, "python")
	firstID := dataset.Rules()[0].ID

	var stdout bytes.Buffer

	err := Run([]string{"explain", "python:zzz_no_such_rule_ever"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	for _, want := range []string{
		"unknown python guideline ids: zzz_no_such_rule_ever",
		"Available ids:",
		firstID,
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func TestExplainRequiresIDs(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"explain"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "requires at least one guideline id") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExplainRequiresLanguageForBareIDs(t *testing.T) {
	dataset := datasetFor(t, "python")
	id := dataset.Rules()[0].ID

	var stdout bytes.Buffer

	err := Run([]string{"explain", id}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "requires a language") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExplainRejectsMalformedReferences(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"empty language", []string{"explain", ":some_rule"}, "invalid guideline reference"},
		{"empty id", []string{"explain", "python:"}, "invalid guideline reference"},
		{"unknown language prefix", []string{"explain", "ruby:some_rule"}, `unknown language "ruby"`},
		{"unknown lang flag", []string{"explain", "--lang", "ruby", "some_rule"}, `unknown language "ruby"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(test.args, &stdout)
			if err == nil {
				t.Fatal("expected error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestProvenanceWarnsOnAdvisorySources(t *testing.T) {
	dataset := datasetFor(t, "typescript")

	tests := []struct {
		name    string
		axis    string
		version registry.Version
		want    string
	}{
		{
			name: "high fidelity has no warning",
			axis: "typescript",
			version: registry.Version{
				Version: "5.4", Source: "package.json devDependencies.typescript", Fidelity: registry.FidelityHigh,
			},
			want: "# typescript 5.4 (source: package.json devDependencies.typescript, fidelity: high)",
		},
		{
			name: "advisory source carries a warning",
			axis: "node",
			version: registry.Version{
				Version: "18", Source: "package.json engines.node", Fidelity: registry.FidelityMedium,
			},
			want: "# node 18 (source: package.json engines.node, fidelity: medium) [advisory: pass --version node=<version> to pin it explicitly]",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := writeProvenance(&stdout, dataset, map[string]registry.Version{
				"typescript": {Version: "5.4", Source: "test", Fidelity: registry.FidelityHigh},
				"node":       {Version: "18", Source: "test", Fidelity: registry.FidelityHigh},
				test.axis:    test.version,
			})
			if err != nil {
				t.Fatal(err)
			}
			output := stdout.String()
			if !strings.Contains(output, test.want+"\n") {
				t.Fatalf("provenance output missing %q:\n%s", test.want, output)
			}
			if strings.Count(output, "[advisory:") > 1 {
				t.Fatalf("unexpected extra warnings:\n%s", output)
			}
		})
	}
}

func TestProvenanceRequiresEveryAxis(t *testing.T) {
	dataset := datasetFor(t, "typescript")

	var stdout bytes.Buffer
	err := writeProvenance(&stdout, dataset, map[string]registry.Version{
		"typescript": {Version: "5.4", Source: "test", Fidelity: registry.FidelityHigh},
	})
	if err == nil || !strings.Contains(err.Error(), `no resolved version for axis "node"`) {
		t.Fatalf("writeProvenance() error = %v, want missing-axis error", err)
	}
}

func TestVersion(t *testing.T) {
	oldVersion := version
	version = "v0.1.0"
	t.Cleanup(func() {
		version = oldVersion
	})

	var stdout bytes.Buffer

	if err := Run([]string{"--version"}, &stdout); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "v0.1.0\n" {
		t.Fatalf("version output = %q", stdout.String())
	}
}

func TestUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"frobnicate"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `unknown command "frobnicate"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHelpListsCommands(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no arguments", nil},
		{"help flag", []string{"--help"}},
		{"help word", []string{"help"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			if err := Run(test.args, &stdout); err != nil {
				t.Fatal(err)
			}
			output := stdout.String()
			for _, want := range []string{"list ", "explain ", "languages", "mcp", "--version"} {
				if !strings.Contains(output, want) {
					t.Fatalf("help output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestSubcommandHelpWritesUsageToStdout(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "list", args: []string{"list", "-h"}, want: "Usage of list:"},
		{name: "explain", args: []string{"explain", "-h"}, want: "Usage of explain:"},
		{name: "languages", args: []string{"languages", "-h"}, want: "Usage of languages:"},
		{name: "mcp", args: []string{"mcp", "-h"}, want: "Usage of mcp:"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(test.args, &stdout)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(stdout.String(), test.want) {
				t.Fatalf("help output missing %q:\n%s", test.want, stdout.String())
			}
		})
	}
}

func TestInvalidSubcommandFlagReturnsSingleErrorWithUsage(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--df"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	output := err.Error()
	if count := strings.Count(output, "flag provided but not defined: -df"); count != 1 {
		t.Fatalf("parse error appeared %d times, want once:\n%s", count, output)
	}
	if !strings.Contains(output, "Usage of list:") {
		t.Fatalf("parse error should include usage:\n%s", output)
	}
}

func TestRunReturnsOutputWriteErrors(t *testing.T) {
	dataset := datasetFor(t, "python")
	newest := dataset.Rules()[0].SinceVersion
	firstID := dataset.Rules()[0].ID

	tests := []struct {
		name string
		args []string
	}{
		{name: "top-level help", args: []string{"--help"}},
		{name: "version output", args: []string{"--version"}},
		{name: "subcommand help", args: []string{"list", "-h"}},
		{name: "list output", args: []string{"list", "--version", "python=" + newest}},
		{name: "explain output", args: []string{"explain", "python:" + firstID}},
		{name: "languages output", args: []string{"languages"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Run(test.args, failingWriter{})
			if !errors.Is(err, errWriteFailed) {
				t.Fatalf("Run() error = %v, want %v", err, errWriteFailed)
			}
		})
	}
}
