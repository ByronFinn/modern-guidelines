// Package detectors implements the per-language version detectors of the
// modern-guidelines hub (PRD-0000 requirement 5). Each detector resolves one
// language's version axes through the three-tier chain: explicit --version
// overrides (T1), declarative project manifests (T2), and the local
// toolchain (T3).
//
// This file provides the Go detector. Its resolution semantics — find-up
// order, language-version defaults, go.work non-leakage — are adapted from
// JetBrains/go-modern-guidelines internal/goversion (Apache-2.0); see the
// repository NOTICE for attribution.
package detectors

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
)

// goAxis is the single version axis of the Go dataset.
const goAxis = "go"

// goVersionFlagSource names the provenance source for T1 overrides, matching
// the source string the CLI records for --version.
const goVersionFlagSource = "--version flag"

// goEnvSource is the provenance source for the `go env GOVERSION` fallback.
const goEnvSource = "go env GOVERSION"

// defaultModuleLanguageVersion is the language version the go command assumes
// for a module whose go.mod has no go directive.
// See https://go.dev/ref/mod#go-mod-file-go.
const defaultModuleLanguageVersion = "1.16"

// defaultWorkspaceLanguageVersion is the language version the go command
// assumes for a workspace whose go.work has no go directive.
// See https://go.dev/ref/mod#workspaces.
const defaultWorkspaceLanguageVersion = "1.18"

// goVersionInText extracts a Go major.minor version from free-form text such
// as "1.24", "go1.24.3", or "go version go1.26.0 linux/amd64".
var goVersionInText = regexp.MustCompile(`(?i)(?:^|\s)(?:go)?(\d+\.\d+)`)

// goDirectiveVersionRE mirrors golang.org/x/mod/modfile.GoVersionRE (v0.37.0,
// the version pinned by the reference project): 1.23, 1.23.4, 1.23rc1, with
// no leading zeros. Directive values are validated strictly; free-form
// --version and `go env` output goes through goVersionInText instead.
var goDirectiveVersionRE = regexp.MustCompile(`^([1-9][0-9]*)\.(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*))?([a-z]+[0-9]+)?$`)

// toolchainNameRE mirrors golang.org/x/mod/modfile.ToolchainRE: "default" or
// a toolchain name starting with "go1".
var toolchainNameRE = regexp.MustCompile(`^default$|^go1($|\.)`)

// goEnvGOVERSION runs `go env GOVERSION` and returns its output. It is a
// package variable so tests can stub the subprocess.
var goEnvGOVERSION = func() (string, error) {
	output, err := exec.Command("go", "env", "GOVERSION").Output()
	return string(output), err
}

func init() {
	registry.RegisterDetector(goDetector{})
}

// goDetector resolves the "go" axis for Go projects.
type goDetector struct{}

func (goDetector) Language() string { return goAxis }

// Detect resolves the Go language version through the three-tier chain
// (PRD-0000 requirement 5): an explicit --version go=... override, the
// go.mod / go.work manifests found upwards from the target path, and the
// local toolchain via `go env GOVERSION`.
func (d goDetector) Detect(req registry.DetectRequest) (map[string]registry.Version, error) {
	if override := strings.TrimSpace(req.VersionOverrides[goAxis]); override != "" {
		version, err := normalizeGoVersion(override)
		if err != nil {
			return nil, fmt.Errorf("cannot parse Go version %q. Pass --version go=1.24 to set the version explicitly", override)
		}
		return goAxisResult(registry.Version{
			Version:  version,
			Source:   goVersionFlagSource,
			Fidelity: registry.FidelityHigh,
		}), nil
	}

	filePath := strings.TrimSpace(req.FilePath)
	if filePath != "" {
		return resolveGoVersionFromPath(filePath)
	}

	// No path given, so the working directory anchors the search. Without
	// this, a language-only request would report the local toolchain even
	// inside a module, offering guidelines the project's own go directive
	// does not allow.
	if workingDir, err := os.Getwd(); err == nil {
		version, ok, err := resolveGoVersionFromModuleFiles(workingDir)
		if err != nil {
			return nil, err
		}
		if ok {
			return goAxisResult(version), nil
		}
	}

	version, err := resolveGoToolchainVersion("the local Go toolchain")
	if err != nil {
		return nil, err
	}
	return goAxisResult(version), nil
}

// resolveGoVersionFromPath resolves the Go version anchored at a file or
// directory path: a go.mod/go.work path answers for itself, anything else
// searches upwards from its directory.
func resolveGoVersionFromPath(filePath string) (map[string]registry.Version, error) {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read --file-path %q: %w", filePath, err)
	}

	if !info.IsDir() && isGoModuleVersionFile(absolutePath) {
		// The named file answers for itself. Searching from its directory
		// instead would let a neighbouring go.mod speak for an explicitly
		// requested go.work.
		version, err := resolveGoVersionFromModuleFile(absolutePath)
		if err != nil {
			return nil, err
		}
		return goAxisResult(version), nil
	}

	searchDir := absolutePath
	if !info.IsDir() {
		searchDir = filepath.Dir(absolutePath)
	}
	version, ok, err := resolveGoVersionFromModuleFiles(searchDir)
	if err != nil {
		return nil, err
	}
	if ok {
		return goAxisResult(version), nil
	}

	version, err = resolveGoToolchainVersion(fmt.Sprintf("the Go toolchain for %s", filePath))
	if err != nil {
		return nil, err
	}
	return goAxisResult(version), nil
}

// resolveGoVersionFromModuleFiles searches go.mod before go.work, each all
// the way up from startDir: a module's own go.mod wins over any enclosing
// go.work, and a go.mod further up beats a go.work nearer by, matching the
// reference implementation's search order (go.mod 优先于 go.work).
func resolveGoVersionFromModuleFiles(startDir string) (registry.Version, bool, error) {
	for _, name := range []string{"go.mod", "go.work"} {
		path := goFindUp(startDir, name)
		if path == "" {
			continue
		}
		version, err := resolveGoVersionFromModuleFile(path)
		if err != nil {
			return registry.Version{}, false, err
		}
		return version, true, nil
	}
	return registry.Version{}, false, nil
}

// resolveGoVersionFromModuleFile resolves the language version a go.mod or
// go.work pins. A missing go directive still pins a language version — 1.16
// for go.mod, 1.18 for go.work — whatever the local toolchain or an enclosing
// go.work says: falling through to either made the reference CLI offer
// features the go command then refuses to compile.
func resolveGoVersionFromModuleFile(path string) (registry.Version, error) {
	directives, err := parseModuleVersionFile(path)
	if err != nil {
		return registry.Version{}, err
	}
	if !directives.hasGoDirective {
		if filepath.Base(path) == "go.work" {
			directives.version = defaultWorkspaceLanguageVersion
		} else {
			directives.version = defaultModuleLanguageVersion
		}
	}
	return registry.Version{
		Version:  directives.version,
		Source:   moduleVersionSource(path, directives),
		Fidelity: registry.FidelityHigh,
	}, nil
}

// moduleVersionSource renders the provenance source for a go.mod/go.work
// answer. The toolchain directive is recorded as an advisory hint only: it
// selects a compiler, it does not raise the language version, which the go
// directive (or its documented default) alone decides (PRD-0000: toolchain
// 指令是建议非地板).
func moduleVersionSource(path string, directives moduleVersionDirectives) string {
	isGoWork := filepath.Base(path) == "go.work"
	source := "go.mod"
	if isGoWork {
		source = "go.work"
	} else if directives.hasGoDirective {
		source = "go.mod go directive"
	}

	var notes []string
	if !directives.hasGoDirective {
		defaultVersion := defaultModuleLanguageVersion
		if isGoWork {
			defaultVersion = defaultWorkspaceLanguageVersion
		}
		notes = append(notes, "no go directive: the go command assumes "+defaultVersion)
	}
	if directives.toolchain != "" {
		notes = append(notes, "toolchain "+directives.toolchain+" advisory, not a language floor")
	}
	if len(notes) == 0 {
		return source
	}
	return source + " (" + strings.Join(notes, "; ") + ")"
}

// moduleVersionDirectives is the version-relevant content of one go.mod or
// go.work file.
type moduleVersionDirectives struct {
	// hasGoDirective reports whether a go directive was present.
	hasGoDirective bool
	// version is the normalized major.minor language version of the go
	// directive; it holds the file's documented default when the directive is
	// absent (filled in by resolveGoVersionFromModuleFile).
	version string
	// toolchain is the raw toolchain directive name (for example "go1.26.0"),
	// empty when absent or "default". Advisory hint only, never a floor.
	toolchain string
}

// parseModuleVersionFile reads the go and toolchain directives from a go.mod
// or go.work file.
//
// The reference project parsed these files with golang.org/x/mod/modfile;
// this hub hand-rolls an equivalent strict reader so the detector adds no
// module dependency. Directive values are held to x/mod's GoVersionRE and
// ToolchainRE with x/mod's error wording. One deliberate divergence: unknown
// directives are ignored rather than rejected, so go.mod files written by
// newer Go toolchains (godebug blocks and beyond) still resolve.
func parseModuleVersionFile(path string) (moduleVersionDirectives, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return moduleVersionDirectives{}, fmt.Errorf("cannot read %s: %w", path, err)
	}

	var directives moduleVersionDirectives
	seenGo := false
	seenToolchain := false
	blockDepth := 0
	for rawLine := range strings.SplitSeq(string(data), "\n") {
		line, _, _ := strings.Cut(rawLine, "//")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if blockDepth > 0 {
			// Inside a require/use/... block every line is an entry, not a
			// language statement.
			if line == ")" {
				blockDepth--
			}
			continue
		}
		if line == ")" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == "(" {
			blockDepth++
			continue
		}
		verb := unquoteModfileToken(fields[0])
		args := make([]string, 0, len(fields)-1)
		for _, field := range fields[1:] {
			args = append(args, unquoteModfileToken(field))
		}
		switch verb {
		case "go":
			if seenGo {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: repeated go statement", path)
			}
			if len(args) != 1 {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: go directive expects exactly one argument", path)
			}
			if !goDirectiveVersionRE.MatchString(args[0]) {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: invalid go version '%s': must match format 1.23.0", path, args[0])
			}
			version, err := normalizeGoVersion(args[0])
			if err != nil {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: invalid go version '%s': must match format 1.23.0", path, args[0])
			}
			seenGo = true
			directives.hasGoDirective = true
			directives.version = version
		case "toolchain":
			if seenToolchain {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: repeated toolchain statement", path)
			}
			if len(args) != 1 {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: toolchain directive expects exactly one argument", path)
			}
			if !toolchainNameRE.MatchString(args[0]) {
				return moduleVersionDirectives{}, fmt.Errorf("cannot parse %s: invalid toolchain version '%s': must match format go1.23.0 or default", path, args[0])
			}
			seenToolchain = true
			if args[0] != "default" {
				directives.toolchain = args[0]
			}
		}
	}
	return directives, nil
}

// unquoteModfileToken strips one pair of surrounding double quotes from a modfile
// token; the modfile grammar allows quoting any token.
func unquoteModfileToken(token string) string {
	if len(token) >= 2 && strings.HasPrefix(token, `"`) && strings.HasSuffix(token, `"`) {
		if unquoted, err := strconv.Unquote(token); err == nil {
			return unquoted
		}
	}
	return token
}

func isGoModuleVersionFile(path string) bool {
	name := filepath.Base(path)
	return name == "go.mod" || name == "go.work"
}

func goFindUp(startDir, name string) string {
	dir := filepath.Clean(startDir)
	for {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// resolveGoToolchainVersion resolves the version of the local Go toolchain
// via `go env GOVERSION` (resolution tier T3, PRD-0000 requirement 5).
func resolveGoToolchainVersion(source string) (registry.Version, error) {
	output, err := goEnvGOVERSION()
	if err != nil {
		return registry.Version{}, fmt.Errorf("cannot determine the Go version from %s: %w. Pass --version go=1.24 to set the version explicitly", source, err)
	}
	version, err := normalizeGoVersion(output)
	if err != nil {
		return registry.Version{}, fmt.Errorf("cannot determine the Go version from %s: %q. Pass --version go=1.24 to set the version explicitly", source, strings.TrimSpace(output))
	}
	return registry.Version{
		Version:  version,
		Source:   goEnvSource,
		Fidelity: registry.FidelityMedium,
	}, nil
}

// normalizeGoVersion extracts the normalized Go major.minor version from
// free-form text such as "1.24", "go1.24.3", "go version go1.26.0
// linux/amd64", or `go env GOVERSION` output.
func normalizeGoVersion(rawVersion string) (string, error) {
	trimmed := strings.TrimSpace(rawVersion)
	match := goVersionInText.FindStringSubmatch(trimmed)
	if len(match) < 2 {
		return "", fmt.Errorf("cannot parse Go version %q", rawVersion)
	}
	majorMinor, ok := parseGoMajorMinor(match[1])
	if !ok {
		return "", fmt.Errorf("cannot parse Go version %q", rawVersion)
	}
	return fmt.Sprintf("%d.%d", majorMinor.major, majorMinor.minor), nil
}

type goMajorMinorVersion struct {
	major int
	minor int
}

func parseGoMajorMinor(version string) (goMajorMinorVersion, bool) {
	majorPart, minorPart, ok := strings.Cut(version, ".")
	if !ok {
		return goMajorMinorVersion{}, false
	}
	minorPart, _, _ = strings.Cut(minorPart, ".")
	major, err := strconv.Atoi(majorPart)
	if err != nil {
		return goMajorMinorVersion{}, false
	}
	minor, err := strconv.Atoi(minorPart)
	if err != nil {
		return goMajorMinorVersion{}, false
	}
	return goMajorMinorVersion{major: major, minor: minor}, true
}

// goAxisResult wraps one resolved version as the single-axis result map.
func goAxisResult(version registry.Version) map[string]registry.Version {
	return map[string]registry.Version{goAxis: version}
}
