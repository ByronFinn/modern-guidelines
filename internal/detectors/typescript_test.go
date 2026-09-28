package detectors

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/registry"
)

// versionExpect is one expected axis resolution: exact version, a substring
// the provenance source must contain, and the expected fidelity grade.
type versionExpect struct {
	version   string
	sourceSub string
	fidelity  registry.Fidelity
}

func TestRegistryWiring(t *testing.T) {
	detector, ok := registry.DetectorForLanguage("typescript")
	if !ok {
		t.Fatal("typescript detector is not registered")
	}
	if language := detector.Language(); language != "typescript" {
		t.Fatalf("Language() = %q, want %q", language, "typescript")
	}
}

func TestSemverRangeFloor(t *testing.T) {
	tests := []struct {
		name   string
		spec   string
		want   string
		wantOK bool
	}{
		{name: "caret", spec: "^5.4", want: "5.4", wantOK: true},
		{name: "tilde", spec: "~5.2", want: "5.2", wantOK: true},
		{name: "tilde legacy alias", spec: "~>5.2", want: "5.2", wantOK: true},
		{name: "space intersection takes max lower bound", spec: ">=5.0 <5.5", want: "5.0", wantOK: true},
		{name: "explicit equals", spec: "=5.4", want: "5.4", wantOK: true},
		{name: "intersection with ceiling only on the right", spec: ">=6.0 <=6.9", want: "6.0", wantOK: true},
		{name: "strict greater folds to the boundary", spec: ">5.4", want: "5.4", wantOK: true},
		{name: "caret keeps prerelease floor", spec: "^5.4.0-beta.1", want: "5.4", wantOK: true},
		{name: "major only", spec: ">=18", want: "18", wantOK: true},
		{name: "union takes the smallest lower bound", spec: "^18 || ^20", want: "18", wantOK: true},
		{name: "union order independent", spec: ">=20 || >=18", want: "18", wantOK: true},
		{name: "leading v tolerated", spec: "v6.1.0", want: "6.1", wantOK: true},
		{name: "minor wildcard pins the major", spec: "18.x", want: "18", wantOK: true},
		{name: "patch wildcard pins major minor", spec: "5.4.x", want: "5.4", wantOK: true},
		{name: "empty range declares nothing", spec: "", wantOK: false},
		{name: "star declares nothing", spec: "*", wantOK: false},
		{name: "dist tag declares nothing", spec: "latest", wantOK: false},
		{name: "protocol range declares nothing", spec: "workspace:^5.4", wantOK: false},
		{name: "ceiling only has no floor", spec: "<5.5", wantOK: false},
		{name: "hyphen range unsupported", spec: "1.2.3 - 2.3.4", wantOK: false},
		{name: "spaceless comparators unsupported", spec: ">=5.0<5.5", wantOK: false},
		{name: "leading zero rejected", spec: ">=05.4", wantOK: false},
		{name: "four segments rejected", spec: "1.2.3.4", wantOK: false},
		{name: "wildcard then concrete rejected", spec: "5.x.4", wantOK: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := semverRangeFloor(test.spec)
			if ok != test.wantOK {
				t.Fatalf("semverRangeFloor(%q) = (%q, %v), want ok = %v", test.spec, got, ok, test.wantOK)
			}
			if ok && got != test.want {
				t.Errorf("semverRangeFloor(%q) = %q, want %q", test.spec, got, test.want)
			}
		})
	}
}

func TestNormalizeBareVersion(t *testing.T) {
	tests := []struct {
		raw    string
		want   string
		wantOK bool
	}{
		{raw: "v20.11.0", want: "20.11", wantOK: true},
		{raw: "V18", want: "18", wantOK: true},
		{raw: "20.11", want: "20.11", wantOK: true},
		{raw: "20", want: "20", wantOK: true},
		{raw: "5.7.2", want: "5.7", wantOK: true},
		{raw: "20.x", want: "20", wantOK: true},
		{raw: "lts/hydrogen", wantOK: false},
		{raw: "node", wantOK: false},
		{raw: "latest", wantOK: false},
		{raw: "", wantOK: false},
	}

	for _, test := range tests {
		t.Run(test.raw, func(t *testing.T) {
			got, ok := normalizeBareVersion(test.raw)
			if ok != test.wantOK {
				t.Fatalf("normalizeBareVersion(%q) = (%q, %v), want ok = %v", test.raw, got, ok, test.wantOK)
			}
			if ok && got != test.want {
				t.Errorf("normalizeBareVersion(%q) = %q, want %q", test.raw, got, test.want)
			}
		})
	}
}

func TestDetectTypeScriptAxis(t *testing.T) {
	disableNodeFallback(t)

	tests := []struct {
		name     string
		files    map[string]string
		target   string
		wantTS   *versionExpect // nil → the axis must be undeclared
		wantNode *versionExpect // nil → the axis must be undeclared
	}{
		{
			name:   "caret range lower bound",
			files:  map[string]string{"package.json": `{"devDependencies":{"typescript":"^5.4"}}`},
			target: "src/app.ts",
			wantTS: &versionExpect{version: "5.4", sourceSub: "package.json devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name:   "tilde range lower bound",
			files:  map[string]string{"package.json": `{"devDependencies":{"typescript":"~5.2"}}`},
			target: "src/app.tsx",
			wantTS: &versionExpect{version: "5.2", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name:   "space separated intersection lower bound",
			files:  map[string]string{"package.json": `{"devDependencies":{"typescript":">=5.0 <5.5"}}`},
			target: "src/app.mts",
			wantTS: &versionExpect{version: "5.0", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name:   "exact version drops the patch",
			files:  map[string]string{"package.json": `{"devDependencies":{"typescript":"5.7.2"}}`},
			target: "src/app.cts",
			wantTS: &versionExpect{version: "5.7", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name: "both axes from one package.json",
			files: map[string]string{
				"package.json": `{"devDependencies":{"typescript":"^5.4"},"engines":{"node":">=18"}}`,
			},
			target:   "src/app.ts",
			wantTS:   &versionExpect{version: "5.4", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
			wantNode: &versionExpect{version: "18", sourceSub: "package.json engines.node", fidelity: registry.FidelityMedium},
		},
		{
			name: "nearest package.json wins over the monorepo root",
			files: map[string]string{
				"package.json":              `{"devDependencies":{"typescript":"^5.4"}}`,
				"packages/app/package.json": `{"devDependencies":{"typescript":"~5.2"}}`,
			},
			target: "packages/app/src/main.ts",
			wantTS: &versionExpect{version: "5.2", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name: "package.json without the field keeps walking up",
			files: map[string]string{
				"package.json":              `{"devDependencies":{"typescript":"^5.4"}}`,
				"packages/app/package.json": `{"name":"app","private":true}`,
			},
			target: "packages/app/src/main.ts",
			wantTS: &versionExpect{version: "5.4", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name: "wildcard dependency keeps walking up to a real pin",
			files: map[string]string{
				"package.json":              `{"devDependencies":{"typescript":"^5.4"}}`,
				"packages/app/package.json": `{"devDependencies":{"typescript":"*"}}`,
			},
			target: "packages/app/src/main.ts",
			wantTS: &versionExpect{version: "5.4", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
		{
			name: "malformed nearest package.json keeps walking up",
			files: map[string]string{
				"package.json":              `{"devDependencies":{"typescript":"^5.4"}}`,
				"packages/app/package.json": `{ not json`,
			},
			target: "packages/app/src/main.ts",
			wantTS: &versionExpect{version: "5.4", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			versions, err := detectAt(t, test.files, test.target, nil)
			if err != nil {
				t.Fatalf("Detect() error: %v", err)
			}
			assertAxis(t, versions, axisTypeScript, test.wantTS)
			assertAxis(t, versions, axisNode, test.wantNode)
		})
	}
}

func TestDetectNodeAxisChain(t *testing.T) {
	disableNodeFallback(t)

	tests := []struct {
		name       string
		files      map[string]string
		target     string
		nodeOutput string // when non-empty, stub `node --version` success with it
		wantNode   versionExpect
	}{
		{
			name:     "engines node lower bound",
			files:    map[string]string{"package.json": `{"engines":{"node":">=18"}}`},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "18", sourceSub: "package.json engines.node", fidelity: registry.FidelityMedium},
		},
		{
			name:     "engines union takes the smallest lower bound",
			files:    map[string]string{"package.json": `{"engines":{"node":"^18 || ^20"}}`},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "18", sourceSub: "engines.node", fidelity: registry.FidelityMedium},
		},
		{
			name: "engines wildcard falls through to nvmrc",
			files: map[string]string{
				"package.json": `{"engines":{"node":"*"}}`,
				".nvmrc":       "v20.11.0",
			},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "20.11", sourceSub: ".nvmrc", fidelity: registry.FidelityHigh},
		},
		{
			name:     "node version file",
			files:    map[string]string{".node-version": "22.0.0"},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "22.0", sourceSub: ".node-version", fidelity: registry.FidelityHigh},
		},
		{
			name: "engines outranks nvmrc",
			files: map[string]string{
				"package.json": `{"engines":{"node":">=16"}}`,
				".nvmrc":       "20.11.0",
			},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "16", sourceSub: "engines.node", fidelity: registry.FidelityMedium},
		},
		{
			name: "unusable nvmrc falls through to tool versions",
			files: map[string]string{
				".nvmrc":         "lts/hydrogen",
				".tool-versions": "ruby 3.4\nnodejs 20.11.0\n",
			},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "20.11", sourceSub: ".tool-versions nodejs", fidelity: registry.FidelityMedium},
		},
		{
			name:     "mise toml node entry",
			files:    map[string]string{"mise.toml": "[tools]\nnode = \"20\"\n"},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "20", sourceSub: "mise.toml tools.node", fidelity: registry.FidelityMedium},
		},
		{
			name:     "mise toml list takes the first pin",
			files:    map[string]string{"mise.toml": "[tools]\nnode = [\"22.11.0\", \"20\"]\n"},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "22.11", sourceSub: "mise.toml tools.node", fidelity: registry.FidelityMedium},
		},
		{
			name:     "tool versions accepts the short node alias",
			files:    map[string]string{".tool-versions": "# pinned tools\nnode 21.7.0\n"},
			target:   "src/app.ts",
			wantNode: versionExpect{version: "21.7", sourceSub: ".tool-versions", fidelity: registry.FidelityMedium},
		},
		{
			name:       "node executable fallback",
			target:     "src/app.ts",
			nodeOutput: "v20.11.0\n",
			wantNode:   versionExpect{version: "20.11", sourceSub: "node --version", fidelity: fidelityLow},
		},
		{
			name:     "toolchain fallback failure leaves the node axis undeclared",
			files:    map[string]string{"package.json": `{"devDependencies":{"typescript":"^5.4"}}`},
			target:   "src/app.ts",
			wantNode: versionExpect{version: undeclaredVersion, sourceSub: "--version node", fidelity: fidelityUnknown},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.nodeOutput != "" {
				stubNodeVersion(t, test.nodeOutput, nil)
			}
			versions, err := detectAt(t, test.files, test.target, nil)
			if err != nil {
				t.Fatalf("Detect() error: %v", err)
			}
			assertAxis(t, versions, axisNode, &test.wantNode)
		})
	}
}

func TestDetectVersionOverrides(t *testing.T) {
	disableNodeFallback(t)

	tests := []struct {
		name      string
		files     map[string]string
		overrides map[string]string
		wantTS    *versionExpect // nil → the axis must be undeclared
		wantNode  *versionExpect // nil → the axis must be undeclared
	}{
		{
			name:      "both axes overridden and normalized",
			overrides: map[string]string{"typescript": "5.4", "node": "v20.11.0"},
			wantTS:    &versionExpect{version: "5.4", sourceSub: sourceVersionFlag, fidelity: registry.FidelityHigh},
			wantNode:  &versionExpect{version: "20.11", sourceSub: sourceVersionFlag, fidelity: registry.FidelityHigh},
		},
		{
			name:      "override beats the manifest",
			files:     map[string]string{"package.json": `{"devDependencies":{"typescript":"^5.4"},"engines":{"node":">=16"}}`},
			overrides: map[string]string{"typescript": "5.8"},
			wantTS:    &versionExpect{version: "5.8", sourceSub: sourceVersionFlag, fidelity: registry.FidelityHigh},
			wantNode:  &versionExpect{version: "16", sourceSub: "engines.node", fidelity: registry.FidelityMedium},
		},
		{
			name:      "single axis override, other resolved from manifest",
			files:     map[string]string{"package.json": `{"devDependencies":{"typescript":"^5.4"}}`},
			overrides: map[string]string{"node": "20"},
			wantTS:    &versionExpect{version: "5.4", sourceSub: "devDependencies.typescript", fidelity: registry.FidelityMedium},
			wantNode:  &versionExpect{version: "20", sourceSub: sourceVersionFlag, fidelity: registry.FidelityHigh},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			versions, err := detectAt(t, test.files, "src/app.ts", test.overrides)
			if err != nil {
				t.Fatalf("Detect() error: %v", err)
			}
			assertAxis(t, versions, axisTypeScript, test.wantTS)
			assertAxis(t, versions, axisNode, test.wantNode)
		})
	}
}

func TestDetectErrors(t *testing.T) {
	disableNodeFallback(t)

	tests := []struct {
		name       string
		files      map[string]string
		overrides  map[string]string
		wantErrSub []string
	}{
		{
			name: "nothing declared anywhere",
			wantErrSub: []string{
				"--version typescript=<ver>,node=<ver>",
				"devDependencies.typescript",
				"node --version",
			},
		},
		{
			name:       "wildcard typescript and no node source",
			files:      map[string]string{"package.json": `{"devDependencies":{"typescript":"*"}}`},
			wantErrSub: []string{"--version typescript=<ver>,node=<ver>"},
		},
		{
			name:       "invalid override value",
			overrides:  map[string]string{"typescript": "next"},
			wantErrSub: []string{`invalid --version value "next" for axis "typescript"`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := detectAt(t, test.files, "src/app.ts", test.overrides)
			if err == nil {
				t.Fatal("Detect() = nil error, want an error")
			}
			for _, sub := range test.wantErrSub {
				if !strings.Contains(err.Error(), sub) {
					t.Errorf("Detect() error %q does not contain %q", err.Error(), sub)
				}
			}
		})
	}
}

func TestDetectAnchorsAtWorkingDirectory(t *testing.T) {
	disableNodeFallback(t)

	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"package.json": `{"devDependencies":{"typescript":"^5.4"}}`,
	})
	t.Chdir(root)

	versions, err := typescriptDetector{}.Detect(registry.DetectRequest{})
	if err != nil {
		t.Fatalf("Detect() error: %v", err)
	}
	assertAxis(t, versions, axisTypeScript, &versionExpect{
		version:   "5.4",
		sourceSub: "devDependencies.typescript",
		fidelity:  registry.FidelityMedium,
	})
}

// detectAt materializes files under a fresh t.TempDir and runs the detector
// anchored at the given target path (relative to the temp dir).
func detectAt(t *testing.T, files map[string]string, target string, overrides map[string]string) (map[string]registry.Version, error) {
	t.Helper()
	root := t.TempDir()
	writeFiles(t, root, files)
	return typescriptDetector{}.Detect(registry.DetectRequest{
		FilePath:         filepath.Join(root, target),
		VersionOverrides: overrides,
	})
}

// assertAxis checks one axis of a Detect result: an expected resolution must
// match exactly (version, source substring, fidelity); a nil expectation must
// be the undeclared marker with its --version guidance.
func assertAxis(t *testing.T, versions map[string]registry.Version, axis string, want *versionExpect) {
	t.Helper()
	got, resolved := versions[axis]
	if !resolved {
		t.Fatalf("Detect() result has no %q axis: %#v", axis, versions)
	}
	if want == nil {
		want = &versionExpect{version: undeclaredVersion, sourceSub: "--version " + axis, fidelity: fidelityUnknown}
	}
	if got.Version != want.version {
		t.Errorf("%s axis version = %q, want %q (source %q)", axis, got.Version, want.version, got.Source)
	}
	if !strings.Contains(got.Source, want.sourceSub) {
		t.Errorf("%s axis source %q does not contain %q", axis, got.Source, want.sourceSub)
	}
	if got.Fidelity != want.fidelity {
		t.Errorf("%s axis fidelity = %q, want %q", axis, got.Fidelity, want.fidelity)
	}
}

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for path, content := range files {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("create dir for %s: %v", path, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
}

// stubNodeVersion replaces the `node --version` subprocess for one test and
// restores it afterwards.
func stubNodeVersion(t *testing.T, output string, err error) {
	t.Helper()
	previous := nodeVersionCommand
	nodeVersionCommand = func() (string, error) { return output, err }
	t.Cleanup(func() { nodeVersionCommand = previous })
}

// disableNodeFallback pins the T3 fallback to a failure so tests never depend
// on whether the host has node installed.
func disableNodeFallback(t *testing.T) {
	t.Helper()
	stubNodeVersion(t, "", errors.New("node: executable file not found in $PATH"))
}
