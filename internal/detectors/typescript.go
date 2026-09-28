// Package detectors holds the per-language version detectors of the hub.
// Each detector resolves one language's declared version axes from project
// manifests (PRD-0000 requirement 5, tier T2) and the local toolchain (tier
// T3) and registers itself into internal/registry via init().
package detectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
)

// Axis names of the typescript dataset (PRD decision D12: one dataset, two
// version axes).
const (
	axisTypeScript = "typescript"
	axisNode       = "node"
)

// Provenance source labels, one per resolution tier and manifest. The
// --version flag label matches cli's versionFlagSource literal.
const (
	sourceVersionFlag       = "--version flag"
	sourceDevDepsTypescript = "package.json devDependencies.typescript"
	sourceEnginesNode       = "package.json engines.node"
	sourceNodeVersion       = "node --version"
	sourceToolVersionsNode  = ".tool-versions nodejs"
)

// undeclaredVersion marks an axis that no source could resolve. registry
// has no "unknown" version state, and downstream filtering
// (guidelines.RulesForVersions → schema.CompareVersions) only accepts numeric
// major[.minor] strings — a literal "unknown" would make `mg list` fail
// outright. "0" is the honest mechanical equivalent: no since_version can be
// ≤ 0, so the axis gates off every rule, while the Source field carries the
// human-readable undeclared explanation and Fidelity reads "unknown".
const undeclaredVersion = "0"

// fidelityLow and fidelityUnknown extend registry's high/medium grades: the
// registry defines only those two constants, and the CLI renders its advisory
// warning ("pass --version axis=<version> ...") for every non-high fidelity,
// which is exactly the behavior these grades require.
const (
	fidelityLow     = registry.Fidelity("low")     // local toolchain fallback (`node --version`)
	fidelityUnknown = registry.Fidelity("unknown") // no source declared the axis
)

func init() {
	registry.RegisterDetector(typescriptDetector{})
}

// typescriptDetector resolves the two axes of the typescript dataset:
//
//	typescript — package.json devDependencies.typescript range, lower bound
//	              (a range is an inference, not a pin, hence medium fidelity);
//	node       — package.json engines.node (advisory, medium) → .nvmrc /
//	              .node-version (a bare pin, high) → mise.toml / .tool-versions
//	              node entry (fuzzy values possible, medium) → `node --version`
//	              (whatever happens to be on PATH, low).
type typescriptDetector struct{}

func (typescriptDetector) Language() string { return "typescript" }

func (typescriptDetector) Detect(req registry.DetectRequest) (map[string]registry.Version, error) {
	versions := make(map[string]registry.Version, 2)

	// T1: explicit --version overrides win per axis, normalized like every
	// other source (v20.11.0 → 20.11). Overrides for other languages' axes,
	// if any ride along in the same request, are ignored.
	for _, axis := range [...]string{axisTypeScript, axisNode} {
		raw, overridden := req.VersionOverrides[axis]
		if !overridden {
			continue
		}
		version, ok := normalizeBareVersion(raw)
		if !ok {
			return nil, fmt.Errorf("invalid --version value %q for axis %q: use a numeric version such as 5.4 or v20.11.0", raw, axis)
		}
		versions[axis] = registry.Version{Version: version, Source: sourceVersionFlag, Fidelity: registry.FidelityHigh}
	}

	startDir := anchorDir(req.FilePath)

	if _, resolved := versions[axisTypeScript]; !resolved {
		if version, ok := nearestPackageJSONVersion(startDir, typescriptDevDependencyVersion); ok {
			versions[axisTypeScript] = version
		}
	}
	if _, resolved := versions[axisNode]; !resolved {
		resolveNodeAxis(startDir, versions)
	}

	// Both axes undeclared is the only hard error: the caller cannot gate
	// anything. A single undeclared axis is reported as an "unknown" marker so
	// the other axis still works.
	_, tsResolved := versions[axisTypeScript]
	_, nodeResolved := versions[axisNode]
	if !tsResolved && !nodeResolved {
		return nil, fmt.Errorf("cannot resolve a version for either axis of this TypeScript project: no devDependencies.typescript range, no engines.node, no .nvmrc or .node-version, no mise.toml or .tool-versions node entry, and %q is unavailable. Declare both axes explicitly: --version typescript=<ver>,node=<ver>", sourceNodeVersion)
	}
	if !tsResolved {
		versions[axisTypeScript] = undeclaredAxisVersion("no package.json up the tree pins devDependencies.typescript (absent, wildcard, or unparsable ranges do not count); pass --version typescript=<version>")
	}
	if !nodeResolved {
		versions[axisNode] = undeclaredAxisVersion("no engines.node, .nvmrc, .node-version, mise node entry, or usable node executable; pass --version node=<version>")
	}
	return versions, nil
}

func undeclaredAxisVersion(explanation string) registry.Version {
	return registry.Version{
		Version:  undeclaredVersion,
		Source:   "undeclared: " + explanation,
		Fidelity: fidelityUnknown,
	}
}

// resolveNodeAxis walks the node-axis chain (PRD requirement 5): engines.node
// → .nvmrc / .node-version → mise.toml / .tool-versions → `node --version`.
// The first source yielding a usable version wins; a source that is present
// but unusable (an "lts/*" .nvmrc, a path-style .tool-versions entry) is
// skipped and the chain continues.
func resolveNodeAxis(startDir string, versions map[string]registry.Version) {
	// Tier 1: engines.node. npm treats engines as advisory (it only warns),
	// so this is medium fidelity — the CLI's provenance output flags it and
	// suggests --version node=<version>.
	if version, ok := nearestPackageJSONVersion(startDir, enginesNodeVersion); ok {
		versions[axisNode] = version
		return
	}

	// Tier 2: .nvmrc / .node-version — a bare version pin, high fidelity.
	if path, found := findUpFile(startDir, ".nvmrc", ".node-version"); found {
		if line, readable := readFirstLine(path); readable {
			if version, ok := normalizeBareVersion(line); ok {
				versions[axisNode] = registry.Version{Version: version, Source: filepath.Base(path), Fidelity: registry.FidelityHigh}
				return
			}
		}
		// An unusable pin (e.g. "lts/hydrogen") falls through to mise.
	}

	// Tier 3: mise / asdf tool pins. Values can be fuzzy ("20"), hence medium.
	if version, source, ok := miseNodeVersion(startDir); ok {
		versions[axisNode] = registry.Version{Version: version, Source: source, Fidelity: registry.FidelityMedium}
		return
	}

	// Tier 4: whatever node is on PATH — the environment, not the project.
	if raw, err := nodeVersionCommand(); err == nil {
		if version, ok := normalizeBareVersion(raw); ok {
			versions[axisNode] = registry.Version{Version: version, Source: sourceNodeVersion, Fidelity: fidelityLow}
		}
	}
}

// typescriptDevDependencyVersion extracts the typescript axis from one
// package.json: the lower bound of the devDependencies.typescript range. A
// wildcard or unparsable range ("*", "latest", "workspace:*") pins nothing,
// so it does not answer and the find-up walk continues — the same monorepo
// rule the PRD states for python's requires-python (requirement 6).
func typescriptDevDependencyVersion(pkg packageJSON) (registry.Version, bool) {
	floor, ok := semverRangeFloor(pkg.DevDependencies[axisTypeScript])
	if !ok {
		return registry.Version{}, false
	}
	return registry.Version{
		Version:  floor,
		Source:   sourceDevDepsTypescript,
		Fidelity: registry.FidelityMedium, // inferred from a range, not pinned
	}, true
}

// enginesNodeVersion extracts the node axis from one package.json's
// engines.node range lower bound (advisory source, medium fidelity).
func enginesNodeVersion(pkg packageJSON) (registry.Version, bool) {
	floor, ok := semverRangeFloor(pkg.Engines.Node)
	if !ok {
		return registry.Version{}, false
	}
	return registry.Version{
		Version:  floor,
		Source:   sourceEnginesNode,
		Fidelity: registry.FidelityMedium, // engines is advisory in npm
	}, true
}

// packageJSON carries the two fields this detector reads; package.json is
// strict JSON, so the standard library is the parsing authority.
type packageJSON struct {
	DevDependencies map[string]string `json:"devDependencies"`
	Engines         struct {
		Node string `json:"node"`
	} `json:"engines"`
}

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// nearestPackageJSONVersion walks startDir and its parents and returns the
// version extracted from the nearest package.json for which extract answers
// (version, true). package.json files that are missing, malformed (npm-style
// BOM tolerated), or whose field yields no usable version do not answer.
func nearestPackageJSONVersion(startDir string, extract func(packageJSON) (registry.Version, bool)) (registry.Version, bool) {
	var found registry.Version
	answered := false
	findUp(startDir, func(dir string) bool {
		data, err := os.ReadFile(filepath.Join(dir, "package.json"))
		if err != nil {
			return false
		}
		data = bytes.TrimPrefix(data, utf8BOM)
		var pkg packageJSON
		if json.Unmarshal(data, &pkg) != nil {
			return false // malformed JSON answers nothing; keep walking
		}
		version, answers := extract(pkg)
		if !answers {
			return false
		}
		found, answered = version, true
		return true
	})
	return found, answered
}

// miseNodeVersion resolves the node entry from mise.toml / .mise.toml
// ([tools] node = "20") or .tool-versions ("nodejs 20.11.0"). Only string
// pins (or a list of pins, first = mise's default) answer; option tables,
// paths, and lts aliases do not.
func miseNodeVersion(startDir string) (version, source string, ok bool) {
	path, found := findUpFile(startDir, "mise.toml", ".mise.toml", ".tool-versions")
	if !found {
		return "", "", false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", false
	}
	switch base := filepath.Base(path); base {
	case ".tool-versions":
		// One "<tool> <version>..." line per tool; the asdf/mise plugin for
		// node is named "nodejs" (mise also accepts "node"). With several
		// versions on one line, mise picks the first.
		for line := range strings.SplitSeq(string(data), "\n") {
			fields := strings.Fields(line)
			if len(fields) < 2 || (fields[0] != "nodejs" && fields[0] != "node") {
				continue
			}
			pinned, usable := normalizeBareVersion(fields[1])
			if !usable {
				return "", "", false // a node entry exists but is unusable
			}
			return pinned, sourceToolVersionsNode, true
		}
		return "", "", false
	default:
		var config struct {
			Tools struct {
				Node any `toml:"node"`
			} `toml:"tools"`
		}
		if toml.Unmarshal(data, &config) != nil {
			return "", "", false
		}
		raw := ""
		switch entry := config.Tools.Node.(type) {
		case string:
			raw = entry
		case []any:
			if len(entry) > 0 {
				if first, isString := entry[0].(string); isString {
					raw = first
				}
			}
		}
		if raw == "" {
			return "", "", false
		}
		pinned, usable := normalizeBareVersion(raw)
		if !usable {
			return "", "", false
		}
		return pinned, base + " tools.node", true
	}
}

// nodeVersionTimeout bounds the `node --version` fallback so a hung node
// process cannot stall a detector call.
const nodeVersionTimeout = 5 * time.Second

// nodeVersionCommand runs `node --version` (tier T3). It is a package-level
// variable so tests can stub the subprocess.
var nodeVersionCommand = func() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), nodeVersionTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, "node", "--version").Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", sourceNodeVersion, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// anchorDir returns the directory the manifest search starts from: the
// directory of the target file, or the working directory when no file path
// was given. A path that names an existing directory anchors at itself.
func anchorDir(filePath string) string {
	if filePath == "" {
		if cwd, err := os.Getwd(); err == nil {
			return cwd
		}
		return "."
	}
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		return filepath.Clean(filePath)
	}
	return filepath.Dir(filepath.Clean(filePath))
}

// findUp calls visit for startDir and every parent directory until visit
// returns true or the filesystem root is reached.
func findUp(startDir string, visit func(dir string) bool) {
	dir := filepath.Clean(startDir)
	for {
		if visit(dir) {
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

// findUpFile walks startDir and its parents; at each level it checks names in
// order and returns the path of the first file that exists. Same-level files
// keep the given name priority; deeper files lose to shallower ones.
func findUpFile(startDir string, names ...string) (string, bool) {
	var found string
	visited := false
	findUp(startDir, func(dir string) bool {
		for _, name := range names {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				found, visited = path, true
				return true
			}
		}
		return false
	})
	return found, visited
}

func readFirstLine(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	line, _, _ := strings.Cut(string(data), "\n")
	return strings.TrimSpace(line), true
}

// semverRangeFloor is the small node-semver range parser shared by both axes
// (devDependencies.typescript and engines.node). It answers one question:
// what is the lowest version the range can match, normalized to major[.minor]?
//
// Supported grammar — the subset package.json manifests actually use:
//
//	range        := comparatorSet ("||" comparatorSet)*  union: floor = min of set floors
//	comparatorSet := comparator (space comparator)*      intersection: floor = max of comparator floors
//	comparator   := (">=" | "<=" | ">" | "<" | "^" | "~" | "~>" | "=")? version
//	version      := "v"? core ["-" prerelease | "+" build]
//	core         := num ["." (num | "x" | "X" | "*") ["." ...]]  a wildcard ends the pin
//
// Deliberate boundaries (documented rather than implemented — no third-party
// semver library per the detector constraints):
//   - upper bounds ("<", "<=") contribute nothing: PRD requirement 6 takes
//     the lower bound as the feature gate and ignores the ceiling;
//   - ">" folds to ">=" at major.minor granularity: ">5.4" admits 5.4.1+, so
//     5.4 features are available in every matching version — only the exact
//     boundary version is lost;
//   - hyphen ranges ("1.2.3 - 2.3.4"), parenthesized ranges, and space-less
//     comparators (">=5.0<5.5") are not parsed: an alternative containing a
//     bare "-" is rejected outright, and any unparseable token contributes no
//     floor, so such ranges report "undeclared" and the chain continues;
//   - dist-tags and protocols ("latest", "workspace:*", "github:u/r", URLs)
//     parse to no floor, i.e. undeclared;
//   - prerelease and build suffixes are stripped (">=5.4.0-beta.1" floors at
//     5.4) and leading zeros are rejected, per node-semver;
//   - the patch segment is dropped ("5.7.2" → "5.7", "v20.11.0" → "20.11")
//     because registry versions are major[.minor]; "5.0" keeps its minor
//     (">=5.0" → "5.0", not "5").
func semverRangeFloor(raw string) (string, bool) {
	var best floorVersion
	haveBest := false
	for alternative := range strings.SplitSeq(raw, "||") {
		setFloor, ok := comparatorSetFloor(alternative)
		if !ok {
			continue
		}
		if !haveBest || setFloor.less(best) {
			best, haveBest = setFloor, true
		}
	}
	if !haveBest {
		return "", false
	}
	return best.render(), true
}

// comparatorSetFloor returns the highest lower bound among the AND-combined
// comparators of one "||" alternative (intersection lower bound = max).
func comparatorSetFloor(alternative string) (floorVersion, bool) {
	var result floorVersion
	haveResult := false
	for token := range strings.FieldsSeq(alternative) {
		if token == "-" { // hyphen-range separator: unsupported syntax
			return floorVersion{}, false
		}
		tokenFloor, ok := comparatorFloor(token)
		if !ok {
			continue // ceiling, wildcard, or unparseable: no floor contribution
		}
		if !haveResult || result.less(tokenFloor) {
			result, haveResult = tokenFloor, true
		}
	}
	return result, haveResult
}

// comparatorFloor returns one comparator's lower bound: "^5.4" → 5.4,
// ">=5.0" → 5.0, "<5.5" → none.
func comparatorFloor(token string) (floorVersion, bool) {
	operation := ""
	for _, candidate := range [...]string{">=", "<=", "~>", "^", "~", ">", "<", "="} {
		if strings.HasPrefix(token, candidate) {
			operation, token = candidate, token[len(candidate):]
			break
		}
	}
	switch operation {
	case "<", "<=":
		return floorVersion{}, false // ceilings are ignored (PRD requirement 6)
	}
	spec, ok := parseVersionSpec(token)
	if !ok || spec.specificity == specificityAny {
		return floorVersion{}, false
	}
	// ">=", ">", "^", "~", "~>", "=", and a bare version all floor at the
	// core itself (see the semverRangeFloor comment for the ">" caveat).
	return floorVersion{major: spec.major, minor: spec.minor, hasMinor: spec.specificity == specificityMajorMinor}, true
}

// Specificity of a parsed version core: how much of the version the pin
// actually constrains. The patch level never matters for a floor, so full
// cores collapse into major.minor.
const (
	specificityAny        = iota // "*", "x": matches everything
	specificityMajor             // "18", "18.x"
	specificityMajorMinor        // "5.4", "5.4.x", "5.4.2"
)

type versionSpec struct {
	major, minor int
	specificity  int
}

// parseVersionSpec parses one semver core ("v20.11.0", "18", "5.4.x",
// "5.4.0-beta.1") into the precision it was declared at. ok is false for
// anything that is not a wildcard or dotted numeric core.
func parseVersionSpec(raw string) (versionSpec, bool) {
	raw = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(raw), "vV"))
	if core, _, found := strings.Cut(raw, "-"); found {
		raw = core // strip prerelease
	}
	if core, _, found := strings.Cut(raw, "+"); found {
		raw = core // strip build metadata
	}
	if raw == "" {
		return versionSpec{}, false
	}
	segments := strings.Split(raw, ".")
	if len(segments) > 3 {
		return versionSpec{}, false
	}
	var numbers [3]int
	concrete := 0 // length of the leading concrete (numeric) segment run
	for index, segment := range segments {
		if index != concrete {
			break // a wildcard ended the pin; the rest must be wildcards too
		}
		if isWildcardSegment(segment) {
			break
		}
		number, ok := parseCoreNumber(segment)
		if !ok {
			return versionSpec{}, false
		}
		numbers[concrete] = number
		concrete++
	}
	for index := concrete; index < len(segments); index++ {
		if !isWildcardSegment(segments[index]) {
			return versionSpec{}, false // "5.x.4"-style mixed pins are rejected
		}
	}
	spec := versionSpec{}
	switch {
	case concrete == 0:
		spec.specificity = specificityAny
	case concrete == 1:
		spec.specificity = specificityMajor
		spec.major = numbers[0]
	default:
		spec.specificity = specificityMajorMinor
		spec.major, spec.minor = numbers[0], numbers[1]
	}
	return spec, true
}

func isWildcardSegment(segment string) bool {
	return segment == "x" || segment == "X" || segment == "*"
}

// parseCoreNumber accepts single digits or multi-digit numbers without a
// leading zero, per node-semver core syntax.
func parseCoreNumber(segment string) (int, bool) {
	if segment == "" || len(segment) > 1 && segment[0] == '0' {
		return 0, false
	}
	for _, r := range segment {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	number, err := strconv.Atoi(segment)
	if err != nil {
		return 0, false
	}
	return number, true
}

// floorVersion is a resolved lower bound in the registry's major[.minor]
// form; the patch segment is intentionally dropped.
type floorVersion struct {
	major, minor int
	hasMinor     bool
}

func (f floorVersion) render() string {
	if f.hasMinor {
		return fmt.Sprintf("%d.%d", f.major, f.minor)
	}
	return strconv.Itoa(f.major)
}

func (f floorVersion) less(other floorVersion) bool {
	if f.major != other.major {
		return f.major < other.major
	}
	return f.minor < other.minor
}

// normalizeBareVersion parses a concrete version token — a .nvmrc line, a
// --version override, `node --version` output — into the registry's
// major[.minor] form: "v20.11.0" → "20.11", "5.7.2" → "5.7", "18" → "18".
// Non-numeric pins ("lts/hydrogen", "latest", "node") report ok=false so the
// resolution chain can continue.
func normalizeBareVersion(raw string) (string, bool) {
	spec, ok := parseVersionSpec(raw)
	if !ok || spec.specificity == specificityAny {
		return "", false
	}
	return (floorVersion{major: spec.major, minor: spec.minor, hasMinor: spec.specificity == specificityMajorMinor}).render(), true
}
