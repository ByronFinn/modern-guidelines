package detectors

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
)

// The tests below port the semantic boundary set of the reference project's
// goversion_test.go (JetBrains/go-modern-guidelines, Apache-2.0): version
// normalization, manifest defaults, go.work non-leakage, explicit-file
// self-answer, and the toolchain fallback.

func TestNormalizeGoVersion(t *testing.T) {
	tests := map[string]string{
		"1.24":                            "1.24",
		"go1.24.3":                        "1.24",
		"Go1.25":                          "1.25",
		"go version go1.26.0 x":           "1.26",
		"go version go1.26.0 linux/amd64": "1.26",
		"go1.21rc1":                       "1.21",
		"go1.24\n":                        "1.24",
		// `go env GOVERSION` output of a development toolchain.
		"devel go1.27-0f9d8a3f9f m=darwin/arm64": "1.27",
	}

	for input, want := range tests {
		got, err := normalizeGoVersion(input)
		if err != nil {
			t.Fatalf("normalizeGoVersion(%q): %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeGoVersion(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeGoVersionRejectsUnparseableText(t *testing.T) {
	// The reference project mapped a bare "devel" override to a configured
	// newest-known version; this hub has no such configuration, so a bare
	// "devel" is unparseable and must guide the user to --version go=... .
	// Full devel toolchain output still normalizes (see TestNormalizeGoVersion).
	for _, input := range []string{"", "devel", "garbage", "1", "go"} {
		if _, err := normalizeGoVersion(input); err == nil {
			t.Fatalf("normalizeGoVersion(%q) succeeded, want a parse error", input)
		}
	}
}

func TestGoDetectorRegisteredInRegistry(t *testing.T) {
	detector, ok := registry.DetectorForLanguage("go")
	if !ok {
		t.Fatal("registry.DetectorForLanguage(\"go\") found no detector; init registration is broken")
	}
	if detector.Language() != "go" {
		t.Fatalf("detector.Language() = %q, want %q", detector.Language(), "go")
	}
}

func TestDetectVersionOverrideNormalizes(t *testing.T) {
	tests := map[string]string{
		"1.24":                            "1.24",
		"go1.24.3":                        "1.24",
		"Go1.25":                          "1.25",
		"go version go1.26.0 linux/amd64": "1.26",
	}

	for override, want := range tests {
		versions, err := goDetector{}.Detect(registry.DetectRequest{VersionOverrides: map[string]string{"go": override}})
		if err != nil {
			t.Fatalf("Detect(--version go=%q): %v", override, err)
		}
		version := singleGoVersion(t, versions)
		if version.Version != want {
			t.Fatalf("Detect(--version go=%q) = %q, want %q", override, version.Version, want)
		}
		if version.Source != goVersionFlagSource {
			t.Fatalf("Detect(--version go=%q) source = %q, want %q", override, version.Source, goVersionFlagSource)
		}
		if version.Fidelity != registry.FidelityHigh {
			t.Fatalf("Detect(--version go=%q) fidelity = %q, want %q", override, version.Fidelity, registry.FidelityHigh)
		}
	}
}

func TestDetectVersionOverrideRejectsUnparseableValue(t *testing.T) {
	_, err := goDetector{}.Detect(registry.DetectRequest{VersionOverrides: map[string]string{"go": "garbage"}})
	if err == nil {
		t.Fatal("Detect succeeded for an unparseable --version override")
	}
	for _, want := range []string{"cannot parse Go version", "--version go=1.24"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Detect error = %q, want it to contain %q", err, want)
		}
	}
}

func TestDetectEmptyOverrideFallsThroughToFilePath(t *testing.T) {
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.24\n")
	goFile := filepath.Join(dir, "main.go")
	writeGoTestFile(t, goFile, "package main\n")

	versions, err := goDetector{}.Detect(registry.DetectRequest{
		FilePath:         goFile,
		VersionOverrides: map[string]string{"go": "  "},
	})
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if got := singleGoVersion(t, versions).Version; got != "1.24" {
		t.Fatalf("Detect = %q, want %q from go.mod", got, "1.24")
	}
}

func TestParseModuleVersionFileReadsGoModDirective(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	writeGoTestFile(t, path, "module example.test/m\n\ngo 1.24.3\n")

	directives, err := parseModuleVersionFile(path)
	if err != nil {
		t.Fatalf("parseModuleVersionFile: %v", err)
	}
	if !directives.hasGoDirective {
		t.Fatal("parseModuleVersionFile did not find the go directive")
	}
	if directives.version != "1.24" {
		t.Fatalf("parseModuleVersionFile = %q, want %q", directives.version, "1.24")
	}
}

func TestParseModuleVersionFileReadsGoWorkDirective(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.work")
	writeGoTestFile(t, path, "go 1.25\n\nuse ./m\n")

	directives, err := parseModuleVersionFile(path)
	if err != nil {
		t.Fatalf("parseModuleVersionFile: %v", err)
	}
	if !directives.hasGoDirective {
		t.Fatal("parseModuleVersionFile did not find the go directive")
	}
	if directives.version != "1.25" {
		t.Fatalf("parseModuleVersionFile = %q, want %q", directives.version, "1.25")
	}
}

func TestParseModuleVersionFileRejectsInvalidGoDirective(t *testing.T) {
	// Ported from the reference TestParseGoDirectiveRejectsInvalidGoModVersion:
	// "go 1.24x" must be rejected, not silently read as 1.24.
	for _, content := range []string{
		"module example.test/m\n\ngo 1.24x\n",
		"module example.test/m\n\ngo 01.24\n",
		"module example.test/m\n\ngo 1.24 1.25\n",
	} {
		path := filepath.Join(t.TempDir(), "go.mod")
		writeGoTestFile(t, path, content)
		if _, err := parseModuleVersionFile(path); err == nil {
			t.Fatalf("parseModuleVersionFile succeeded for %q", content)
		}
	}

	path := filepath.Join(t.TempDir(), "go.mod")
	writeGoTestFile(t, path, "module example.test/m\n\ngo 1.24x\n")
	_, err := parseModuleVersionFile(path)
	if err == nil || !strings.Contains(err.Error(), "invalid go version") {
		t.Fatalf("parseModuleVersionFile error = %v, want it to contain %q", err, "invalid go version")
	}
}

func TestParseModuleVersionFileRejectsInvalidToolchainDirective(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	writeGoTestFile(t, path, "module example.test/m\n\ngo 1.24\n\ntoolchain gophers\n")

	_, err := parseModuleVersionFile(path)
	if err == nil || !strings.Contains(err.Error(), "invalid toolchain version") {
		t.Fatalf("parseModuleVersionFile error = %v, want it to contain %q", err, "invalid toolchain version")
	}
}

func TestParseModuleVersionFileRejectsRepeatedGoDirective(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	writeGoTestFile(t, path, "module example.test/m\n\ngo 1.24\n\ngo 1.25\n")

	_, err := parseModuleVersionFile(path)
	if err == nil || !strings.Contains(err.Error(), "repeated go statement") {
		t.Fatalf("parseModuleVersionFile error = %v, want it to contain %q", err, "repeated go statement")
	}
}

func TestParseModuleVersionFileIgnoresCommentsAndBlockEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	writeGoTestFile(t, path, `module example.test/m

// go 1.99

go 1.24 // comment

require (
	go 1.98
	golang.org/x/mod v0.37.0
)

godebug default=go1.21
`)
	directives, err := parseModuleVersionFile(path)
	if err != nil {
		t.Fatalf("parseModuleVersionFile: %v", err)
	}
	if !directives.hasGoDirective || directives.version != "1.24" {
		t.Fatalf("parseModuleVersionFile = %+v, want the go 1.24 directive to win over comments and block entries", directives)
	}
}

func TestDetectFromNearestGoMod(t *testing.T) {
	dir := t.TempDir()
	moduleDir := filepath.Join(dir, "module")
	pkgDir := filepath.Join(moduleDir, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeGoTestFile(t, filepath.Join(moduleDir, "go.mod"), "module example.test/m\n\ngo 1.24\n")
	goFile := filepath.Join(pkgDir, "x.go")
	writeGoTestFile(t, goFile, "package pkg\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.24" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.24")
	}
	if version.Source != "go.mod go directive" {
		t.Fatalf("Detect source = %q, want %q", version.Source, "go.mod go directive")
	}
	if version.Fidelity != registry.FidelityHigh {
		t.Fatalf("Detect fidelity = %q, want %q", version.Fidelity, registry.FidelityHigh)
	}
}

func TestDetectGoModWithoutGoDirectiveUsesLanguageDefault(t *testing.T) {
	// The go command reads a module whose go.mod has no go directive as
	// language go1.16. Falling back to the local toolchain here made the
	// reference CLI recommend features the go command then refuses to compile.
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n")
	goFile := filepath.Join(dir, "main.go")
	writeGoTestFile(t, goFile, "package main\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.16" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.16")
	}
	if want := "go.mod (no go directive: the go command assumes 1.16)"; version.Source != want {
		t.Fatalf("Detect source = %q, want %q", version.Source, want)
	}
	if version.Fidelity != registry.FidelityHigh {
		t.Fatalf("Detect fidelity = %q, want %q", version.Fidelity, registry.FidelityHigh)
	}
}

func TestDetectGoModPathWithoutGoDirectiveUsesLanguageDefault(t *testing.T) {
	dir := t.TempDir()
	goMod := filepath.Join(dir, "go.mod")
	writeGoTestFile(t, goMod, "module example.test/m\n")

	version := detectGoFile(t, goMod)
	if version.Version != "1.16" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.16")
	}
}

func TestDetectGoModWithoutGoDirectiveIgnoresParentGoWork(t *testing.T) {
	// A go.work above the module does not change the go command's answer: the
	// member module is still built at go1.16, so the workspace version would
	// be the wrong one to inherit.
	dir := t.TempDir()
	moduleDir := filepath.Join(dir, "m")
	if err := os.MkdirAll(moduleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeGoTestFile(t, filepath.Join(dir, "go.work"), "go 1.25\n\nuse ./m\n")
	writeGoTestFile(t, filepath.Join(moduleDir, "go.mod"), "module example.test/m\n")
	goFile := filepath.Join(moduleDir, "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.16" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.16")
	}
}

func TestDetectGoWorkWithoutGoDirectiveUsesWorkspaceLanguageDefault(t *testing.T) {
	// The go command reads a workspace whose go.work has no go directive as
	// go1.18, not as the local toolchain.
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.work"), "use ./m\n")
	goFile := filepath.Join(dir, "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.18" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.18")
	}
	if want := "go.work (no go directive: the go command assumes 1.18)"; version.Source != want {
		t.Fatalf("Detect source = %q, want %q", version.Source, want)
	}
}

func TestDetectGoWorkPathWithoutGoDirectiveIgnoresNeighbouringGoMod(t *testing.T) {
	// The file path names the go.work, so the workspace default answers even
	// though the go.mod beside it carries a go directive.
	dir := t.TempDir()
	goWork := filepath.Join(dir, "go.work")
	writeGoTestFile(t, goWork, "use .\n")
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.24\n")

	version := detectGoFile(t, goWork)
	if version.Version != "1.18" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.18")
	}
}

func TestDetectGoWorkPathUsesItsGoDirective(t *testing.T) {
	dir := t.TempDir()
	goWork := filepath.Join(dir, "go.work")
	writeGoTestFile(t, goWork, "go 1.25\n\nuse .\n")
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.24\n")

	version := detectGoFile(t, goWork)
	if version.Version != "1.25" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.25")
	}
	if version.Source != "go.work" {
		t.Fatalf("Detect source = %q, want %q", version.Source, "go.work")
	}
}

func TestDetectDirectoryPrefersGoModOverGoWork(t *testing.T) {
	// Passing the directory rather than a file keeps the search order: go.mod
	// wins, so the workspace default must not leak into this case.
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.work"), "use .\n")
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.24\n")

	version := detectGoFile(t, dir)
	if version.Version != "1.24" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.24")
	}
}

func TestDetectAncestorGoModBeatsNearerGoWork(t *testing.T) {
	// go.mod is searched all the way up before go.work is considered, so an
	// ancestor module wins even over a go.work right beside the file.
	dir := t.TempDir()
	subDir := filepath.Join(dir, "workspace", "m")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.22\n")
	writeGoTestFile(t, filepath.Join(dir, "workspace", "go.work"), "go 1.25\n\nuse ./m\n")
	goFile := filepath.Join(subDir, "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.22" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.22")
	}
}

func TestDetectToolchainDirectiveIsAdvisoryHintNotFloor(t *testing.T) {
	// The toolchain directive selects a compiler; the language version stays
	// pinned by the go directive (PRD-0000: toolchain 是建议非地板).
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.24\n\ntoolchain go1.26.0\n")
	goFile := filepath.Join(dir, "main.go")
	writeGoTestFile(t, goFile, "package main\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.24" {
		t.Fatalf("Detect = %q, want %q: the toolchain directive must not raise the language version", version.Version, "1.24")
	}
	if want := "go.mod go directive (toolchain go1.26.0 advisory, not a language floor)"; version.Source != want {
		t.Fatalf("Detect source = %q, want %q", version.Source, want)
	}
	if version.Fidelity != registry.FidelityHigh {
		t.Fatalf("Detect fidelity = %q, want %q", version.Fidelity, registry.FidelityHigh)
	}
}

func TestDetectToolchainWithoutGoDirectiveKeepsDefaultAndRecordsHint(t *testing.T) {
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ntoolchain go1.26.0\n")
	goFile := filepath.Join(dir, "main.go")
	writeGoTestFile(t, goFile, "package main\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.16" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.16")
	}
	want := "go.mod (no go directive: the go command assumes 1.16; toolchain go1.26.0 advisory, not a language floor)"
	if version.Source != want {
		t.Fatalf("Detect source = %q, want %q", version.Source, want)
	}
}

func TestDetectToolchainDefaultRecordsNoHint(t *testing.T) {
	dir := t.TempDir()
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.24\n\ntoolchain default\n")
	goFile := filepath.Join(dir, "main.go")
	writeGoTestFile(t, goFile, "package main\n")

	version := detectGoFile(t, goFile)
	if version.Source != "go.mod go directive" {
		t.Fatalf("Detect source = %q, want %q without a toolchain hint", version.Source, "go.mod go directive")
	}
}

func TestDetectGoWorkToolchainRecordedAsHint(t *testing.T) {
	dir := t.TempDir()
	goWork := filepath.Join(dir, "go.work")
	writeGoTestFile(t, goWork, "go 1.25\n\ntoolchain go1.24.3\n\nuse .\n")

	version := detectGoFile(t, goWork)
	if version.Version != "1.25" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.25")
	}
	if want := "go.work (toolchain go1.24.3 advisory, not a language floor)"; version.Source != want {
		t.Fatalf("Detect source = %q, want %q", version.Source, want)
	}
}

func TestDetectWithoutFilePathReadsModuleInWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeGoTestFile(t, filepath.Join(dir, "go.mod"), "module example.test/m\n\ngo 1.21\n")
	t.Chdir(pkgDir)

	versions, err := goDetector{}.Detect(registry.DetectRequest{})
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if got := singleGoVersion(t, versions).Version; got != "1.21" {
		t.Fatalf("Detect = %q, want %q", got, "1.21")
	}
}

func TestDetectWithoutFilePathFallsBackToGoEnvOutsideAModule(t *testing.T) {
	// No go.mod or go.work anywhere above the working directory, so the local
	// toolchain remains the only available answer.
	stubGoEnv(t, "go1.26.0\n", nil)
	t.Chdir(t.TempDir())

	versions, err := goDetector{}.Detect(registry.DetectRequest{})
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	version := singleGoVersion(t, versions)
	if version.Version != "1.26" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.26")
	}
	if version.Source != goEnvSource {
		t.Fatalf("Detect source = %q, want %q", version.Source, goEnvSource)
	}
	if version.Fidelity != registry.FidelityMedium {
		t.Fatalf("Detect fidelity = %q, want %q", version.Fidelity, registry.FidelityMedium)
	}
}

func TestDetectFilePathWithoutModuleFallsBackToGoEnv(t *testing.T) {
	stubGoEnv(t, "go1.26.0\n", nil)
	goFile := filepath.Join(t.TempDir(), "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	version := detectGoFile(t, goFile)
	if version.Version != "1.26" {
		t.Fatalf("Detect = %q, want %q", version.Version, "1.26")
	}
	if version.Fidelity != registry.FidelityMedium {
		t.Fatalf("Detect fidelity = %q, want %q", version.Fidelity, registry.FidelityMedium)
	}
}

func TestDetectGoEnvDevelOutputNormalizes(t *testing.T) {
	stubGoEnv(t, "devel go1.27-0f9d8a3f9f m=darwin/arm64\n", nil)
	goFile := filepath.Join(t.TempDir(), "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	if got := detectGoFile(t, goFile).Version; got != "1.27" {
		t.Fatalf("Detect = %q, want %q", got, "1.27")
	}
}

func TestDetectGoEnvFailureGuidesToVersionFlag(t *testing.T) {
	stubGoEnv(t, "", errors.New("go: executable file not found"))
	goFile := filepath.Join(t.TempDir(), "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	_, err := goDetector{}.Detect(registry.DetectRequest{FilePath: goFile})
	if err == nil {
		t.Fatal("Detect succeeded although `go env GOVERSION` failed")
	}
	for _, want := range []string{"cannot determine the Go version", "--version go=1.24"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Detect error = %q, want it to contain %q", err, want)
		}
	}
}

func TestDetectGoEnvUnparseableOutputGuidesToVersionFlag(t *testing.T) {
	stubGoEnv(t, "nonsense\n", nil)
	goFile := filepath.Join(t.TempDir(), "x.go")
	writeGoTestFile(t, goFile, "package m\n")

	_, err := goDetector{}.Detect(registry.DetectRequest{FilePath: goFile})
	if err == nil {
		t.Fatal("Detect succeeded although `go env GOVERSION` output was unparseable")
	}
	if !strings.Contains(err.Error(), "--version go=1.24") {
		t.Fatalf("Detect error = %q, want it to contain %q", err, "--version go=1.24")
	}
}

func TestDetectUnreadableFilePath(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.go")
	_, err := goDetector{}.Detect(registry.DetectRequest{FilePath: missing})
	if err == nil || !strings.Contains(err.Error(), "cannot read --file-path") {
		t.Fatalf("Detect error = %v, want it to contain %q", err, "cannot read --file-path")
	}
}

func TestDetectUnreadableGoMod(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "go.mod")
	_, err := goDetector{}.Detect(registry.DetectRequest{FilePath: missing})
	if err == nil || !strings.Contains(err.Error(), "cannot read --file-path") {
		t.Fatalf("Detect error = %v, want it to contain %q", err, "cannot read --file-path")
	}
}

// detectGoFile runs the detector for one path and returns the single "go"
// axis version.
func detectGoFile(t *testing.T, path string) registry.Version {
	t.Helper()
	versions, err := goDetector{}.Detect(registry.DetectRequest{FilePath: path})
	if err != nil {
		t.Fatalf("Detect(%q): %v", path, err)
	}
	return singleGoVersion(t, versions)
}

func singleGoVersion(t *testing.T, versions map[string]registry.Version) registry.Version {
	t.Helper()
	if len(versions) != 1 {
		t.Fatalf("Detect returned %d axis versions, want exactly the go axis: %v", len(versions), versions)
	}
	version, ok := versions[goAxis]
	if !ok {
		t.Fatalf("Detect returned no %q axis version: %v", goAxis, versions)
	}
	return version
}

// stubGoEnv replaces the `go env GOVERSION` subprocess for one test.
func stubGoEnv(t *testing.T, output string, err error) {
	t.Helper()
	original := goEnvGOVERSION
	goEnvGOVERSION = func() (string, error) { return output, err }
	t.Cleanup(func() { goEnvGOVERSION = original })
}

func writeGoTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
