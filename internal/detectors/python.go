// Package detectors holds the per-language version detectors. Every detector
// self-registers from its own init(), so the main package only needs a blank
// import for the registration side effect.
package detectors

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// pythonVersionFlagSource matches the CLI's provenance label for T1 versions,
// so overrides resolved through DetectRequest read the same in output.
const pythonVersionFlagSource = "--version flag"

// pythonAxis is the single version axis of the python dataset.
const pythonAxis = "python"

func init() {
	registry.RegisterDetector(pythonDetector{})
}

// pythonDetector resolves the python version axis (PRD requirements 5 and 6):
//
//	T1  an explicit override value (the --version flag);
//	T2  declarative manifests — pyproject.toml requires-python (or
//	    [tool.poetry] python), PEP 723 inline script metadata,
//	    .python-version, mise.toml, .tool-versions — read on the target file
//	    itself or found by climbing parent directories;
//	T3  nothing: Python has no local-toolchain fallback, so an exhausted
//	    chain returns an error guiding the user to --version.
type pythonDetector struct{}

func (pythonDetector) Language() string { return pythonAxis }

func (pythonDetector) Detect(req registry.DetectRequest) (map[string]registry.Version, error) {
	version, err := resolvePythonVersion(req)
	if err != nil {
		return nil, err
	}
	return map[string]registry.Version{pythonAxis: version}, nil
}

func resolvePythonVersion(req registry.DetectRequest) (registry.Version, error) {
	if raw, ok := req.VersionOverrides[pythonAxis]; ok && strings.TrimSpace(raw) != "" {
		return pythonVersionFromOverride(raw)
	}

	filePath := strings.TrimSpace(req.FilePath)
	startDir, err := pythonSearchDir(filePath)
	if err != nil {
		return registry.Version{}, err
	}

	// The named file answers for itself when it is — or carries — a manifest:
	// a pyproject.toml or .python-version passed as the target, or a .py
	// script with a PEP 723 metadata block. A file that carries no answer
	// (for example a pyproject.toml without requires-python) does not block
	// the chain: find-up starts from its directory and keeps climbing.
	if filePath != "" {
		if version, ok := pythonSelfAnswer(filePath); ok {
			return version, nil
		}
	}

	// Find-up, source type first and proximity second: the nearest
	// pyproject.toml declaring a floor outranks a closer .python-version,
	// which outranks mise.toml and .tool-versions (PRD requirement 5).
	for _, answer := range []func(string) (registry.Version, bool){
		pythonPyprojectAnswer,
		pythonVersionFileAnswer,
		pythonMiseAnswer,
		pythonToolVersionsAnswer,
	} {
		if version, ok := findUpAnswer(startDir, answer); ok {
			return version, nil
		}
	}

	return registry.Version{}, fmt.Errorf(
		"cannot determine the Python version: no pyproject.toml requires-python (or [tool.poetry] python), .python-version, or mise.toml/.tool-versions python entry found in %s or any parent directory. Python has no local-toolchain fallback; pass --version python=3.12 explicitly", startDir)
}

// pythonVersionFromOverride normalizes an explicit T1 value. Accepted forms
// include 3.12, python3.12, "Python 3.12.4", and 3.12.4 — everything reduces
// to major.minor.
func pythonVersionFromOverride(raw string) (registry.Version, error) {
	trimmed := strings.TrimSpace(raw)
	version, _, ok := normalizePythonVersionText(trimmed)
	if !ok {
		return registry.Version{}, fmt.Errorf(
			"cannot parse Python version %q: use forms like 3.12, python3.12, or \"Python 3.12.4\"", trimmed)
	}
	return registry.Version{Version: version, Source: pythonVersionFlagSource, Fidelity: registry.FidelityHigh}, nil
}

// pythonSearchDir returns the directory the find-up starts from: the target
// file's directory, or the working directory when no file was given. A
// missing target file is tolerated — an agent may pass the path of a file it
// is about to create, and its would-be directory is still the right anchor.
func pythonSearchDir(filePath string) (string, error) {
	if filePath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("cannot determine the working directory: %w", err)
		}
		return workingDir, nil
	}
	absolute, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot resolve file path %q: %w", filePath, err)
	}
	return filepath.Dir(absolute), nil
}

// pythonSelfAnswer lets the target file answer for itself, mirroring the
// reference goversion's treatment of a directly named go.mod/go.work.
func pythonSelfAnswer(filePath string) (registry.Version, bool) {
	switch filepath.Base(filePath) {
	case "pyproject.toml":
		return pythonPyprojectAnswer(filepath.Dir(filePath))
	case ".python-version":
		return pythonVersionFileAnswer(filepath.Dir(filePath))
	}
	if strings.EqualFold(filepath.Ext(filePath), ".py") {
		return pythonPEP723Answer(filePath)
	}
	return registry.Version{}, false
}

// pythonPyprojectAnswer answers from dir/pyproject.toml when it declares a
// version floor: [project] requires-python first (PEP 621), then
// [tool.poetry] python. A pyproject.toml that declares no floor — or that
// fails to parse — is not an answer; the chain keeps climbing (PRD
// requirement 6: a monorepo member without requires-python does not speak for
// its parents' packages).
func pythonPyprojectAnswer(dir string) (registry.Version, bool) {
	data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml"))
	if err != nil {
		return registry.Version{}, false
	}
	version, source, ok := pythonPyprojectFloor(data)
	if !ok {
		return registry.Version{}, false
	}
	return registry.Version{Version: version, Source: source, Fidelity: registry.FidelityHigh}, true
}

// pyprojectFile mirrors the pyproject.toml fields the detector reads.
type pyprojectFile struct {
	Project struct {
		RequiresPython string `toml:"requires-python"`
	} `toml:"project"`
	Tool struct {
		Poetry struct {
			Python string `toml:"python"`
		} `toml:"poetry"`
	} `toml:"tool"`
}

// pythonPyprojectFloor extracts the declared floor and its provenance source.
func pythonPyprojectFloor(data []byte) (version, source string, ok bool) {
	var doc pyprojectFile
	if err := toml.Unmarshal(data, &doc); err != nil {
		return "", "", false
	}
	if version, ok := pythonSpecFloor(doc.Project.RequiresPython); ok {
		return version, "pyproject.toml requires-python", true
	}
	// [tool.poetry] python uses Poetry's constraint grammar; the same floor
	// extraction covers its operator forms (^3.9, >=3.8,<4.0, ~=3.10), plus
	// the bare pin ("3.9") that PEP 440 specifier syntax would reject.
	poetry := doc.Tool.Poetry.Python
	if version, ok := pythonSpecFloor(poetry); ok {
		return version, "pyproject.toml [tool.poetry] python", true
	}
	if version, _, ok := normalizePythonVersionText(poetry); ok {
		return version, "pyproject.toml [tool.poetry] python", true
	}
	return "", "", false
}

// pythonVersionFileAnswer answers from dir/.python-version (pyenv, mise, uv):
// the first non-empty, non-comment line pins the interpreter.
func pythonVersionFileAnswer(dir string) (registry.Version, bool) {
	data, err := os.ReadFile(filepath.Join(dir, ".python-version"))
	if err != nil {
		return registry.Version{}, false
	}
	line, ok := firstMeaningfulLine(data)
	if !ok {
		return registry.Version{}, false
	}
	return pythonManifestValue(line, ".python-version")
}

// pythonPEP723Answer answers from a PEP 723 inline-script-metadata block
// (https://peps.python.org/pep-0723/): the requires-python value inside the
// script's `# /// script` block, subject to the same floor semantics.
func pythonPEP723Answer(path string) (registry.Version, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return registry.Version{}, false
	}
	spec, ok := pythonPEP723RequiresPython(data)
	if !ok {
		return registry.Version{}, false
	}
	version, ok := pythonSpecFloor(spec)
	if !ok {
		return registry.Version{}, false
	}
	return registry.Version{Version: version, Source: "PEP 723 requires-python", Fidelity: registry.FidelityHigh}, true
}

// pythonPEP723RequiresPython extracts the requires-python value from the
// script metadata block: comment lines between `# /// script` and `# ///`.
func pythonPEP723RequiresPython(data []byte) (string, bool) {
	inBlock := false
	for line := range strings.SplitSeq(string(data), "\n") {
		text := strings.TrimSpace(line)
		if !strings.HasPrefix(text, "#") {
			continue
		}
		content := strings.TrimSpace(strings.TrimPrefix(text, "#"))
		if !inBlock {
			inBlock = content == "/// script"
			continue
		}
		if content == "///" {
			return "", false // the block ended without a requires-python line
		}
		key, value, found := strings.Cut(content, "=")
		if !found || strings.TrimSpace(key) != "requires-python" {
			continue
		}
		return strings.Trim(strings.TrimSpace(value), `"'`), true
	}
	return "", false
}

// pythonMiseAnswer answers from dir/mise.toml ([tools] python entry).
func pythonMiseAnswer(dir string) (registry.Version, bool) {
	data, err := os.ReadFile(filepath.Join(dir, "mise.toml"))
	if err != nil {
		return registry.Version{}, false
	}
	var doc struct {
		Tools map[string]any `toml:"tools"`
	}
	if err := toml.Unmarshal(data, &doc); err != nil {
		return registry.Version{}, false
	}
	raw, ok := pythonMiseToolValue(doc.Tools["python"])
	if !ok {
		return registry.Version{}, false
	}
	return pythonManifestValue(raw, "mise.toml python")
}

// pythonMiseToolValue accepts the shapes mise allows for a tool version:
// "3.12", ["3.12", "3.11"], or { version = "3.12" }.
func pythonMiseToolValue(raw any) (string, bool) {
	switch value := raw.(type) {
	case string:
		return value, true
	case []any:
		for _, item := range value {
			if text, ok := item.(string); ok {
				return text, true
			}
		}
	case map[string]any:
		if text, ok := value["version"].(string); ok {
			return text, true
		}
	}
	return "", false
}

// pythonToolVersionsAnswer answers from dir/.tool-versions (asdf, mise): the
// first `python <version>` line.
func pythonToolVersionsAnswer(dir string) (registry.Version, bool) {
	data, err := os.ReadFile(filepath.Join(dir, ".tool-versions"))
	if err != nil {
		return registry.Version{}, false
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 && fields[0] == "python" {
			return pythonManifestValue(fields[1], ".tool-versions python")
		}
	}
	return registry.Version{}, false
}

// pythonManifestValue turns a manifest version value into a registry
// Version. Single-component values such as mise's "3" mean "any 3.x" rather
// than 3.0 specifically, so they resolve to 3.0 at reduced fidelity (the
// provenance output warns on non-high fidelity and suggests --version).
func pythonManifestValue(raw, source string) (registry.Version, bool) {
	version, fuzzy, ok := normalizePythonVersionText(raw)
	if !ok {
		return registry.Version{}, false
	}
	fidelity := registry.FidelityHigh
	if fuzzy {
		fidelity = registry.FidelityMedium
	}
	return registry.Version{Version: version, Source: source, Fidelity: fidelity}, true
}

// findUpAnswer climbs from startDir to the filesystem root, asking each
// directory for an answer; the first directory that has one wins.
func findUpAnswer(startDir string, answer func(dir string) (registry.Version, bool)) (registry.Version, bool) {
	dir := filepath.Clean(startDir)
	for {
		if version, ok := answer(dir); ok {
			return version, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return registry.Version{}, false
		}
		dir = parent
	}
}

// firstMeaningfulLine returns the first non-empty, non-comment line.
func firstMeaningfulLine(data []byte) (string, bool) {
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line, true
	}
	return "", false
}

// pythonSpecifier matches one version-specifier clause: a PEP 440 operator
// (plus Poetry's ^ and ~) followed by a version. Trailing wildcards and
// pre-release segments (3.11.*, 3.13rc1) parse and drop beyond major.minor.
var pythonSpecifier = regexp.MustCompile(
	`^\s*(>=|<=|===|==|~=|\^|!=|>|<|~)\s*v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.\*)?[0-9A-Za-z.+!-]*\s*$`)

// pythonSpecFloor extracts the effective feature floor from a version
// specifier string (PRD requirement 6):
//
//   - ">=3.10,<3.13" → 3.10   (upper bounds are ignored),
//   - "==3.11.*"     → 3.11   (a pin is its own floor),
//   - "~=3.10"       → 3.10   (compatible release),
//   - ">=3.9,>=3.11" → 3.11   (the strongest lower bound wins),
//   - "<3.13"        → none   (no lower bound declared; the caller keeps
//     climbing the chain).
//
// Clause forms outside this subset are skipped rather than failing the whole
// chain over one exotic manifest.
func pythonSpecFloor(spec string) (string, bool) {
	spec, _, _ = strings.Cut(spec, ";") // drop PEP 508 environment markers
	floor := ""
	for clause := range strings.SplitSeq(spec, ",") {
		match := pythonSpecifier.FindStringSubmatch(clause)
		if match == nil {
			continue
		}
		operator, major, minor := match[1], match[2], match[3]
		switch operator {
		case "<", "<=", "!=":
			// Upper bounds and exclusions never lower the floor.
			continue
		}
		// ">X" still runs code written for X (X.1 and up); compatible-release
		// and pin operators (^, ~, ~=, ==, ===) all start at their stated
		// version.
		version := major + ".0"
		if minor != "" {
			version = major + "." + minor
		}
		if floor == "" {
			floor = version
			continue
		}
		if stronger, err := schema.CompareVersions(version, floor); err == nil && stronger > 0 {
			floor = version
		}
	}
	return floor, floor != ""
}

// pythonVersionInText locates a Python version in free-form text: 3.12,
// python3.12, "Python 3.12.4", 3.12.4rc1. Release-level detail beyond
// major.minor is dropped — the axis gates features on major.minor.
var pythonVersionInText = regexp.MustCompile(`(?i)(?:^|\s)(?:python)?\s*v?(\d+)(?:\.(\d+))?`)

// normalizePythonVersionText reduces a version-shaped string to major.minor.
// fuzzy reports single-component values ("3"): they name a whole release
// line, not a specific floor, and callers mark them medium fidelity.
func normalizePythonVersionText(raw string) (version string, fuzzy, ok bool) {
	match := pythonVersionInText.FindStringSubmatch(strings.TrimSpace(raw))
	if match == nil {
		return "", false, false
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return "", false, false
	}
	minor := 0
	if match[2] != "" {
		if minor, err = strconv.Atoi(match[2]); err != nil {
			return "", false, false
		}
	} else {
		fuzzy = true
	}
	return strconv.Itoa(major) + "." + strconv.Itoa(minor), fuzzy, true
}
