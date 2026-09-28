// Package cli implements the mg command line interface: list, explain,
// languages, the built-in MCP server, --version, and help (PRD-0000 requirements
// 4, 5, and 7). All output goes to an injected io.Writer.
//
// Every list output starts with one provenance line per axis — version,
// source, and fidelity — and advisory or ambiguous resolutions (anything not
// high fidelity) get a warning recommending an explicit --version. The lines
// render from whatever the registry resolves: the --version path fills them
// in already, and manifest/toolchain detectors plug into the same code path
// once their packages register.
package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/ByronFinn/modern-guidelines/internal/guidelines"
	"github.com/ByronFinn/modern-guidelines/internal/mcpserver"
	"github.com/ByronFinn/modern-guidelines/internal/registry"
	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// listVersionSourceConflictError follows the reference project's error
// style: one sentence naming the accepted sources.
const listVersionSourceConflictError = "list accepts only one version source: --version, --file-path, or one positional file path. Run list -h for usage"

// versionFlagSource is the provenance source recorded for versions that came
// from the --version flag (resolution tier T1).
const versionFlagSource = "--version flag"

var version = detectedVersion()

// Run executes the CLI with args, writes command output to stdout, and
// returns user-facing errors.
func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return printUsage(stdout)
	}

	switch args[0] {
	case "-h", "--help", "help":
		return printUsage(stdout)
	case "--version", "version":
		return printVersion(stdout)
	case "list":
		return runList(args[1:], stdout)
	case "explain":
		return runExplain(args[1:], stdout)
	case "languages":
		return runLanguages(args[1:], stdout)
	case "mcp":
		return runMCP(args[1:], stdout)
	default:
		return fmt.Errorf("unknown command %q\n\nRun with --help for usage.", args[0])
	}
}

func detectedVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "dev"
	}
	return info.Main.Version
}

func runList(args []string, stdout io.Writer) error {
	var filePath string
	var versionSpec string
	var langOverride string
	fs := newFlagSet("list", "list [--version axis=ver[,...] | --file-path <path> | <path>]")
	fs.StringVar(&filePath, "file-path", "", "Optional path to a project file; infers the language from its extension and resolves the axis versions.")
	fs.StringVar(&versionSpec, "version", "", "Optional explicit axis versions, comma-separated, covering every axis of one language: --version python=3.12 or --version typescript=5.4,node=18.")
	fs.StringVar(&langOverride, "lang", "", "Optional language override when --file-path or a positional path has an unknown extension.")
	parsed, err := parseFlagSet(fs, args, stdout)
	if err != nil {
		return err
	}
	if !parsed {
		return nil
	}

	positional := fs.Args()
	sources := 0
	if versionSpec != "" {
		sources++
	}
	if filePath != "" {
		sources++
	}
	if len(positional) > 0 {
		sources++
	}
	if sources > 1 {
		return errors.New(listVersionSourceConflictError)
	}
	if len(positional) > 1 {
		return errors.New("list accepts at most one positional file path. Run list -h for usage")
	}
	if len(positional) == 1 {
		filePath = positional[0]
	}
	if versionSpec != "" && langOverride != "" {
		return errors.New("list --lang cannot be combined with --version; the axis names already select the language. Run list -h for usage")
	}
	if versionSpec == "" && filePath == "" {
		return errors.New("list requires a version source: --version axis=ver[,...], --file-path <path>, or one positional file path. Run list -h for usage")
	}

	dataset, versions, err := resolveListVersions(versionSpec, filePath, langOverride)
	if err != nil {
		return err
	}
	return writeListOutput(stdout, dataset, versions)
}

// resolveListVersions resolves the target dataset and its per-axis versions
// from exactly one version source: the --version flag, or a file path with a
// registered detector.
func resolveListVersions(versionSpec, filePath, langOverride string) (registry.Dataset, map[string]registry.Version, error) {
	if versionSpec != "" {
		overrides, err := parseVersionSpec(versionSpec)
		if err != nil {
			return nil, nil, err
		}
		dataset, err := datasetForVersionAxes(overrides)
		if err != nil {
			return nil, nil, err
		}
		versions := make(map[string]registry.Version, len(overrides))
		for axis, value := range overrides {
			versions[axis] = registry.Version{Version: value, Source: versionFlagSource, Fidelity: registry.FidelityHigh}
		}
		return dataset, versions, nil
	}

	language, known := registry.LanguageForPath(filePath)
	if langOverride != "" {
		language, known = langOverride, true
	}
	if !known {
		return nil, nil, unknownExtensionError(filePath)
	}
	dataset, ok := registry.DatasetForLanguage(language)
	if !ok {
		return nil, nil, fmt.Errorf("no guideline dataset is installed for language %q. Installed languages: %s. Run mg languages", language, strings.Join(installedLanguages(), ", "))
	}
	detector, ok := registry.DetectorForLanguage(language)
	if !ok {
		return nil, nil, fmt.Errorf("no version detector is available for language %q yet; pass --version instead, for example %s", language, versionExampleFor(dataset))
	}
	versions, err := detector.Detect(registry.DetectRequest{FilePath: filePath})
	if err != nil {
		return nil, nil, err
	}
	return dataset, versions, nil
}

func writeListOutput(w io.Writer, dataset registry.Dataset, versions map[string]registry.Version) error {
	resolved := make(map[string]string, len(versions))
	for axis, version := range versions {
		resolved[axis] = version.Version
	}
	text, err := guidelines.ListText(dataset, resolved)
	if err != nil {
		return err
	}
	if err := writeProvenance(w, dataset, versions); err != nil {
		return err
	}
	if text == "" {
		return nil
	}
	_, err = fmt.Fprintln(w, text)
	return err
}

// writeProvenance renders one provenance line per axis (PRD requirement 5:
// version, source, fidelity per axis). Resolutions that are not high
// fidelity — advisory or ambiguous sources such as engines.node or fuzzy
// mise values — carry a warning suggesting an explicit --version.
func writeProvenance(w io.Writer, dataset registry.Dataset, versions map[string]registry.Version) error {
	for _, axis := range dataset.Axes() {
		resolved, ok := versions[axis]
		if !ok {
			return fmt.Errorf("no resolved version for axis %q of language %q", axis, dataset.Language())
		}
		line := fmt.Sprintf("# %s %s (source: %s, fidelity: %s)", axis, resolved.Version, resolved.Source, resolved.Fidelity)
		if resolved.Fidelity != registry.FidelityHigh {
			line += fmt.Sprintf(" [advisory: pass --version %s=<version> to pin it explicitly]", axis)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// parseVersionSpec parses a comma-separated "axis=version" list (T1). Every
// version is validated through schema.CompareVersions, the single version
// authority.
func parseVersionSpec(spec string) (map[string]string, error) {
	overrides := make(map[string]string)
	for part := range strings.SplitSeq(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		axis, value, found := strings.Cut(part, "=")
		axis = strings.TrimSpace(axis)
		value = strings.TrimSpace(value)
		if !found || axis == "" || value == "" {
			return nil, fmt.Errorf("invalid --version entry %q: use axis=version, for example --version python=3.12 or --version typescript=5.4,node=18", part)
		}
		if _, err := schema.CompareVersions(value, value); err != nil {
			return nil, fmt.Errorf("invalid --version entry %q: %v", part, err)
		}
		if _, exists := overrides[axis]; exists {
			return nil, fmt.Errorf("duplicate axis %q in --version", axis)
		}
		overrides[axis] = value
	}
	if len(overrides) == 0 {
		return nil, errors.New("--version requires at least one axis=version entry, for example --version python=3.12. Run list -h for usage")
	}
	return overrides, nil
}

// datasetForVersionAxes finds the single dataset whose axes cover the given
// overrides, and requires every axis of that dataset to be set.
func datasetForVersionAxes(overrides map[string]string) (registry.Dataset, error) {
	known := knownAxes()
	var unknown []string
	for axis := range overrides {
		if !slices.Contains(known, axis) {
			unknown = append(unknown, axis)
		}
	}
	if len(unknown) > 0 {
		slices.Sort(unknown)
		return nil, fmt.Errorf("unknown version axis %s in --version; known axes: %s. Run list -h for usage", strings.Join(quoteAll(unknown), ", "), strings.Join(known, ", "))
	}

	var candidates []registry.Dataset
	for _, dataset := range registry.Datasets() {
		covers := true
		for axis := range overrides {
			if !slices.Contains(dataset.Axes(), axis) {
				covers = false
				break
			}
		}
		if covers {
			candidates = append(candidates, dataset)
		}
	}
	if len(candidates) != 1 {
		return nil, errors.New("--version axes must belong to exactly one language. Run list -h for usage")
	}

	dataset := candidates[0]
	var missing []string
	example := make([]string, 0, len(dataset.Axes()))
	for _, axis := range dataset.Axes() {
		if value, ok := overrides[axis]; ok {
			example = append(example, axis+"="+value)
			continue
		}
		missing = append(missing, axis)
		example = append(example, axis+"=<version>")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("--version is missing the %s axis of language %q; set every axis, for example --version %s", strings.Join(missing, ", "), dataset.Language(), strings.Join(example, ","))
	}
	return dataset, nil
}

func runExplain(args []string, stdout io.Writer) error {
	var language string
	fs := newFlagSet("explain", "explain [--lang <language>] <lang:id | id>...")
	fs.StringVar(&language, "lang", "", "Optional language for bare guideline ids, for example --lang python. Explicit lang:id references override it.")
	parsed, err := parseFlagSet(fs, args, stdout)
	if err != nil {
		return err
	}
	if !parsed {
		return nil
	}

	refs, err := ruleRefs(fs.Args(), language)
	if err != nil {
		return err
	}

	var blocks []string
	for _, group := range groupRefsByLanguage(refs) {
		block, err := guidelines.ExplainText(group.dataset, group.ids)
		if err != nil {
			return err
		}
		blocks = append(blocks, block)
	}
	_, err = fmt.Fprintln(stdout, strings.Join(blocks, "\n\n"))
	return err
}

// ruleRef is one resolved (dataset, id) explain reference.
type ruleRef struct {
	dataset registry.Dataset
	id      string
}

// ruleRefs resolves guideline references: each is either "lang:id" or a bare
// id qualified by the --lang flag. An explicit lang:id prefix wins over
// --lang for that reference.
func ruleRefs(args []string, languageFlag string) ([]ruleRef, error) {
	ids := guidelines.NormalizeIDs(args)
	if len(ids) == 0 {
		return nil, errors.New("explain command requires at least one guideline id. Run list to list available ids")
	}
	if languageFlag != "" {
		if _, ok := registry.DatasetForLanguage(languageFlag); !ok {
			return nil, fmt.Errorf("unknown language %q. Available languages: %s. Run mg languages", languageFlag, strings.Join(installedLanguages(), ", "))
		}
	}

	refs := make([]ruleRef, 0, len(ids))
	seen := make(map[ruleRef]bool)
	for _, id := range ids {
		language, ruleID := languageFlag, id
		if prefix, rest, found := strings.Cut(id, ":"); found {
			if prefix == "" || rest == "" {
				return nil, fmt.Errorf("invalid guideline reference %q: use lang:id or bare ids with --lang", id)
			}
			language, ruleID = prefix, rest
		}
		if language == "" {
			return nil, fmt.Errorf("explain requires a language for id %q: use lang:id or pass --lang <language>", ruleID)
		}
		dataset, ok := registry.DatasetForLanguage(language)
		if !ok {
			return nil, fmt.Errorf("unknown language %q in %q. Available languages: %s. Run mg languages", language, id, strings.Join(installedLanguages(), ", "))
		}
		ref := ruleRef{dataset: dataset, id: ruleID}
		if seen[ref] {
			continue
		}
		seen[ref] = true
		refs = append(refs, ref)
	}
	return refs, nil
}

// languageGroup batches explain references of one language so unknown-id
// errors list every unknown id of that language together.
type languageGroup struct {
	dataset registry.Dataset
	ids     []string
}

func groupRefsByLanguage(refs []ruleRef) []languageGroup {
	var groups []languageGroup
	indexByLanguage := make(map[string]int, len(refs))
	for _, ref := range refs {
		language := ref.dataset.Language()
		if index, ok := indexByLanguage[language]; ok {
			groups[index].ids = append(groups[index].ids, ref.id)
			continue
		}
		indexByLanguage[language] = len(groups)
		groups = append(groups, languageGroup{dataset: ref.dataset, ids: []string{ref.id}})
	}
	return groups
}

func runLanguages(args []string, stdout io.Writer) error {
	fs := newFlagSet("languages", "languages")
	parsed, err := parseFlagSet(fs, args, stdout)
	if err != nil {
		return err
	}
	if !parsed {
		return nil
	}
	for _, dataset := range registry.Datasets() {
		if _, err := fmt.Fprintf(stdout, "%s (axes: %s)\n", dataset.Language(), strings.Join(dataset.Axes(), ", ")); err != nil {
			return err
		}
	}
	return nil
}

// runMCP serves the built-in read-only MCP server over stdio (PRD-0000
// requirement 9). The two tools reuse the exact list and explain code paths
// by dispatching back into Run with synthesized flags, so CLI and MCP output
// stay byte-identical; mcpserver owns only the protocol layer and cannot
// import this package (this subcommand already imports mcpserver, so a
// back-dependency would be an import cycle).
func runMCP(args []string, stdout io.Writer) error {
	fs := newFlagSet("mcp", "mcp")
	parsed, err := parseFlagSet(fs, args, stdout)
	if err != nil {
		return err
	}
	if !parsed {
		return nil
	}
	return mcpserver.Run(context.Background(), os.Stdin, os.Stdout, mcpserver.Deps{
		List: func(versionSpec, filePath, language string) (string, error) {
			return runCommandOutput(mcpListArgs(versionSpec, filePath, language))
		},
		Explain: func(language string, ids []string) (string, error) {
			return runCommandOutput(append([]string{"explain", "--lang", language}, ids...))
		},
	})
}

// mcpListArgs rebuilds the list flag arguments from the list_guidelines tool
// parameters: versionSpec → --version, filePath → --file-path, language →
// --lang. Flag-level validation (one version source, axis coverage) then
// happens in the shared runList path.
func mcpListArgs(versionSpec, filePath, language string) []string {
	args := []string{"list"}
	if versionSpec != "" {
		args = append(args, "--version", versionSpec)
	}
	if filePath != "" {
		args = append(args, "--file-path", filePath)
	}
	if language != "" {
		args = append(args, "--lang", language)
	}
	return args
}

// runCommandOutput runs one CLI command against a buffer and returns its raw
// output (trailing newline included; mcpserver trims it into the tool text).
func runCommandOutput(args []string) (string, error) {
	var output bytes.Buffer
	if err := Run(args, &output); err != nil {
		return "", err
	}
	return output.String(), nil
}

func newFlagSet(name, usage string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n  %s\n\nFlags:\n", name, usage)
		fs.PrintDefaults()
	}
	return fs
}

func parseFlagSet(fs *flag.FlagSet, args []string, stdout io.Writer) (bool, error) {
	var output bytes.Buffer
	fs.SetOutput(&output)

	err := fs.Parse(args)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, flag.ErrHelp) {
		_, err := fmt.Fprint(stdout, output.String())
		return false, err
	}
	if output.Len() > 0 {
		return false, errors.New(strings.TrimRight(output.String(), "\n"))
	}
	return false, err
}

func printUsage(w io.Writer) error {
	_, err := fmt.Fprint(w, `mg provides version-gated modern coding guidelines for AI agents, so they can write up-to-date code despite their knowledge cutoff.

Commands:
  list [--version axis=ver[,...] | --file-path <path> | <path>]
      Return the guidelines of one language for its resolved axis versions,
      one per line, newest first; multi-axis datasets are grouped by axis and
      the output starts with one provenance line per axis.

  explain [--lang <language>] <lang:id | id>...
      Return detailed guidance, before/after examples, references, and the
      autofix mapping for specific guideline ids.

  languages
      List the installed guideline datasets and their version axes.

  mcp
      Serve the built-in MCP server over stdio (see README for client config).

  version, --version
      Print the mg version.
`)
	return err
}

func printVersion(w io.Writer) error {
	_, err := fmt.Fprintln(w, version)
	return err
}

// unknownExtensionError reports an unidentifiable path and lists the
// languages that actually have datasets installed (PRD requirement 4:
// 未知扩展名报错列出支持语言).
func unknownExtensionError(filePath string) error {
	detail := fmt.Sprintf("unknown extension %q", filepath.Ext(filePath))
	if filepath.Ext(filePath) == "" {
		detail = "no file extension"
	}
	return fmt.Errorf("cannot identify the language of %q: %s. Supported languages: %s. Run mg languages", filePath, detail, strings.Join(installedLanguages(), ", "))
}

func installedLanguages() []string {
	datasets := registry.Datasets()
	languages := make([]string, 0, len(datasets))
	for _, dataset := range datasets {
		languages = append(languages, dataset.Language())
	}
	return languages
}

func knownAxes() []string {
	seen := make(map[string]bool)
	var axes []string
	for _, dataset := range registry.Datasets() {
		for _, axis := range dataset.Axes() {
			if !seen[axis] {
				seen[axis] = true
				axes = append(axes, axis)
			}
		}
	}
	slices.Sort(axes)
	return axes
}

func versionExampleFor(dataset registry.Dataset) string {
	parts := make([]string, 0, len(dataset.Axes()))
	for _, axis := range dataset.Axes() {
		parts = append(parts, axis+"=<version>")
	}
	return "--version " + strings.Join(parts, ",")
}

func quoteAll(values []string) []string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("%q", value))
	}
	return quoted
}
