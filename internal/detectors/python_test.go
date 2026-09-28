package detectors_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	// The python detector self-registers from its init(); importing the
	// package for that side effect is the whole point of this file.
	_ "github.com/ByronFinn/modern-guidelines/internal/detectors"
	"github.com/ByronFinn/modern-guidelines/internal/registry"
)

// pythonDetector returns the registered python detector, failing the test
// when registration is broken.
func pythonDetector(t *testing.T) registry.Detector {
	t.Helper()
	detector, ok := registry.DetectorForLanguage("python")
	if !ok {
		t.Fatal("python detector is not registered; its init() must call registry.RegisterDetector")
	}
	return detector
}

// detectPython runs the registered python detector and asserts the
// single-axis result shape python must return.
func detectPython(t *testing.T, req registry.DetectRequest) (registry.Version, error) {
	t.Helper()
	versions, err := pythonDetector(t).Detect(req)
	if err != nil {
		return registry.Version{}, err
	}
	if len(versions) != 1 {
		t.Fatalf("python is a single-axis language, got %d axes: %v", len(versions), versions)
	}
	version, ok := versions["python"]
	if !ok {
		t.Fatalf("detect returned no python axis: %v", versions)
	}
	return version, nil
}

// writePythonFile creates path (and parent directories) with the given
// content and returns path.
func writePythonFile(t *testing.T, path, content string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("cannot create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("cannot write %s: %v", path, err)
	}
	return path
}

func TestPythonDetectorRegistration(t *testing.T) {
	if got := pythonDetector(t).Language(); got != "python" {
		t.Errorf("detector language = %q, want %q", got, "python")
	}
}

func TestPythonDetectOverrideNormalization(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr string
	}{
		{name: "plain major.minor", raw: "3.12", want: "3.12"},
		{name: "python prefix attached", raw: "python3.12", want: "3.12"},
		{name: "python prefix with space", raw: "python 3.11", want: "3.11"},
		{name: "interpreter banner", raw: "Python 3.12.4", want: "3.12"},
		{name: "full release", raw: "3.12.4", want: "3.12"},
		{name: "major only", raw: "3", want: "3.0"},
		{name: "not a version", raw: "banana", wantErr: `cannot parse Python version "banana"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			version, err := detectPython(t, registry.DetectRequest{
				VersionOverrides: map[string]string{"python": test.raw},
			})
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("Detect error = %v, want it to contain %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}
			if version.Version != test.want {
				t.Errorf("override %q: version = %q, want %q", test.raw, version.Version, test.want)
			}
			if version.Source != "--version flag" {
				t.Errorf("override %q: source = %q, want %q", test.raw, version.Source, "--version flag")
			}
			if version.Fidelity != registry.FidelityHigh {
				t.Errorf("override %q: fidelity = %q, want %q", test.raw, version.Fidelity, registry.FidelityHigh)
			}
		})
	}
}

func TestPythonDetectRequiresPythonSemantics(t *testing.T) {
	tests := []struct {
		name           string
		requiresPython string
		want           string
	}{
		{name: "lower bound wins, upper bound ignored", requiresPython: ">=3.10,<3.13", want: "3.10"},
		{name: "wildcard pin", requiresPython: "==3.11.*", want: "3.11"},
		{name: "exact pin", requiresPython: "==3.11", want: "3.11"},
		{name: "arbitrary equality", requiresPython: "===3.11.2", want: "3.11"},
		{name: "compatible release", requiresPython: "~=3.10", want: "3.10"},
		{name: "compatible release with patch", requiresPython: "~=3.10.2", want: "3.10"},
		{name: "strongest of several lower bounds", requiresPython: ">=3.9, >=3.11", want: "3.11"},
		{name: "strict greater still floors at the stated minor", requiresPython: ">3.9", want: "3.9"},
		{name: "patch-level floor truncates to minor", requiresPython: ">=3.12.4", want: "3.12"},
		{name: "prerelease floor truncates", requiresPython: ">=3.13rc1", want: "3.13"},
		{name: "environment marker is stripped", requiresPython: ">=3.8; sys_platform != 'win32'", want: "3.8"},
		{name: "whitespace tolerated", requiresPython: " >= 3.9 , < 4.0 ", want: "3.9"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			pyproject := writePythonFile(t, filepath.Join(dir, "pyproject.toml"),
				"[project]\nname = \"app\"\nrequires-python = \""+test.requiresPython+"\"\n")
			version, err := detectPython(t, registry.DetectRequest{FilePath: pyproject})
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}
			if version.Version != test.want {
				t.Errorf("requires-python %q: version = %q, want %q", test.requiresPython, version.Version, test.want)
			}
			if version.Source != "pyproject.toml requires-python" {
				t.Errorf("requires-python %q: source = %q, want %q", test.requiresPython, version.Source, "pyproject.toml requires-python")
			}
			if version.Fidelity != registry.FidelityHigh {
				t.Errorf("requires-python %q: fidelity = %q, want %q", test.requiresPython, version.Fidelity, registry.FidelityHigh)
			}
		})
	}
}

func TestPythonDetectRequiresPythonChainBehavior(t *testing.T) {
	t.Run("no lower bound means undeclared, chain continues", func(t *testing.T) {
		dir := t.TempDir()
		pyproject := writePythonFile(t, filepath.Join(dir, "pyproject.toml"),
			"[project]\nname = \"app\"\nrequires-python = \"<3.13\"\n")
		writePythonFile(t, filepath.Join(dir, ".python-version"), "3.12.4\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: pyproject})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.12" || version.Source != ".python-version" {
			t.Errorf("floor-less requires-python did not continue the chain: got %+v", version)
		}
	})

	t.Run("pyproject without requires-python keeps climbing", func(t *testing.T) {
		root := t.TempDir()
		writePythonFile(t, filepath.Join(root, "pyproject.toml"),
			"[project]\nname = \"mono\"\nrequires-python = \">=3.9\"\n")
		member := writePythonFile(t, filepath.Join(root, "svc", "pyproject.toml"),
			"[project]\nname = \"svc\"\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: member})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.9" || version.Source != "pyproject.toml requires-python" {
			t.Errorf("silent member pyproject did not keep climbing: got %+v", version)
		}
	})

	t.Run("unparseable pyproject is skipped, not fatal", func(t *testing.T) {
		root := t.TempDir()
		writePythonFile(t, filepath.Join(root, "pyproject.toml"),
			"[project]\nname = \"mono\"\nrequires-python = \">=3.9\"\n")
		broken := writePythonFile(t, filepath.Join(root, "svc", "pyproject.toml"),
			"this is = not = toml ]]\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: broken})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.9" {
			t.Errorf("version = %q, want 3.9", version.Version)
		}
	})
}

func TestPythonDetectPoetry(t *testing.T) {
	tests := []struct {
		name       string
		pyproject  string
		want       string
		wantSource string
	}{
		{
			name:       "caret constraint",
			pyproject:  "[tool.poetry]\nname = \"app\"\npython = \"^3.9\"\n",
			want:       "3.9",
			wantSource: "pyproject.toml [tool.poetry] python",
		},
		{
			name:       "range constraint",
			pyproject:  "[tool.poetry]\npython = \">=3.8,<4.0\"\n",
			want:       "3.8",
			wantSource: "pyproject.toml [tool.poetry] python",
		},
		{
			name:       "bare pin",
			pyproject:  "[tool.poetry]\npython = \"3.9\"\n",
			want:       "3.9",
			wantSource: "pyproject.toml [tool.poetry] python",
		},
		{
			name:       "requires-python wins over poetry",
			pyproject:  "[project]\nrequires-python = \">=3.10\"\n[tool.poetry]\npython = \"^3.8\"\n",
			want:       "3.10",
			wantSource: "pyproject.toml requires-python",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			pyproject := writePythonFile(t, filepath.Join(dir, "pyproject.toml"), test.pyproject)
			version, err := detectPython(t, registry.DetectRequest{FilePath: pyproject})
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}
			if version.Version != test.want {
				t.Errorf("version = %q, want %q", version.Version, test.want)
			}
			if version.Source != test.wantSource {
				t.Errorf("source = %q, want %q", version.Source, test.wantSource)
			}
			if version.Fidelity != registry.FidelityHigh {
				t.Errorf("fidelity = %q, want %q", version.Fidelity, registry.FidelityHigh)
			}
		})
	}
}

func TestPythonDetectExplicitFileSelfAnswer(t *testing.T) {
	t.Run("pyproject.toml answers for itself over sibling .python-version", func(t *testing.T) {
		dir := t.TempDir()
		writePythonFile(t, filepath.Join(dir, ".python-version"), "3.11.9\n")
		pyproject := writePythonFile(t, filepath.Join(dir, "pyproject.toml"),
			"[project]\nname = \"app\"\nrequires-python = \">=3.10,<3.13\"\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: pyproject})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.10" || version.Source != "pyproject.toml requires-python" {
			t.Errorf("explicit pyproject did not self-answer: got %+v", version)
		}
	})

	t.Run(".python-version answers for itself", func(t *testing.T) {
		dir := t.TempDir()
		versionFile := writePythonFile(t, filepath.Join(dir, ".python-version"), "3.12.4\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: versionFile})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.12" || version.Source != ".python-version" {
			t.Errorf("explicit .python-version did not self-answer: got %+v", version)
		}
		if version.Fidelity != registry.FidelityHigh {
			t.Errorf("fidelity = %q, want %q", version.Fidelity, registry.FidelityHigh)
		}
	})
}

func TestPythonDetectPEP723(t *testing.T) {
	t.Run("script metadata block answers", func(t *testing.T) {
		dir := t.TempDir()
		script := writePythonFile(t, filepath.Join(dir, "fetch.py"), `#!/usr/bin/env python3
# /// script
# requires-python = ">=3.11"
# dependencies = ["requests"]
# ///
import sys
`)
		version, err := detectPython(t, registry.DetectRequest{FilePath: script})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.11" || version.Source != "PEP 723 requires-python" {
			t.Errorf("PEP 723 block did not answer: got %+v", version)
		}
		if version.Fidelity != registry.FidelityHigh {
			t.Errorf("fidelity = %q, want %q", version.Fidelity, registry.FidelityHigh)
		}
	})

	t.Run("floor-less PEP 723 constraint continues the chain", func(t *testing.T) {
		dir := t.TempDir()
		script := writePythonFile(t, filepath.Join(dir, "fetch.py"), `# /// script
# requires-python = "<3.13"
# ///
`)
		writePythonFile(t, filepath.Join(dir, "pyproject.toml"),
			"[project]\nname = \"app\"\nrequires-python = \">=3.9\"\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: script})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.9" || version.Source != "pyproject.toml requires-python" {
			t.Errorf("floor-less PEP 723 did not continue the chain: got %+v", version)
		}
	})

	t.Run("script without a block falls through to find-up", func(t *testing.T) {
		dir := t.TempDir()
		script := writePythonFile(t, filepath.Join(dir, "fetch.py"), "import sys\n")
		writePythonFile(t, filepath.Join(dir, ".python-version"), "3.12.4\n")
		version, err := detectPython(t, registry.DetectRequest{FilePath: script})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.12" || version.Source != ".python-version" {
			t.Errorf("block-less script did not fall through: got %+v", version)
		}
	})
}

func TestPythonDetectFindUpMonorepo(t *testing.T) {
	t.Run("nearest answering pyproject outranks closer .python-version files", func(t *testing.T) {
		root := t.TempDir()
		writePythonFile(t, filepath.Join(root, "pyproject.toml"),
			"[project]\nname = \"mono\"\nrequires-python = \">=3.9\"\n")
		writePythonFile(t, filepath.Join(root, ".python-version"), "3.11.9\n")
		// Silent member pyproject: no requires-python, so it must not answer
		// and must not block the climb to the root.
		writePythonFile(t, filepath.Join(root, "svc", "pyproject.toml"),
			"[project]\nname = \"svc\"\n")
		writePythonFile(t, filepath.Join(root, "svc", "app", ".python-version"), "3.12.4\n")
		target := writePythonFile(t, filepath.Join(root, "svc", "app", "main.py"), "print(\"hi\")\n")

		version, err := detectPython(t, registry.DetectRequest{FilePath: target})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.9" || version.Source != "pyproject.toml requires-python" {
			t.Errorf("find-up priority wrong: got %+v", version)
		}
	})

	t.Run("nearest .python-version answers when pyprojects are silent", func(t *testing.T) {
		root := t.TempDir()
		writePythonFile(t, filepath.Join(root, "pyproject.toml"), "[project]\nname = \"silent\"\n")
		writePythonFile(t, filepath.Join(root, "svc", ".python-version"), "3.12.4\n")
		target := writePythonFile(t, filepath.Join(root, "svc", "main.py"), "x = 1\n")

		version, err := detectPython(t, registry.DetectRequest{FilePath: target})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.12" || version.Source != ".python-version" {
			t.Errorf("nearest .python-version did not answer: got %+v", version)
		}
	})

	t.Run("pyproject outranks mise", func(t *testing.T) {
		dir := t.TempDir()
		writePythonFile(t, filepath.Join(dir, "pyproject.toml"),
			"[project]\nname = \"app\"\nrequires-python = \">=3.10\"\n")
		writePythonFile(t, filepath.Join(dir, "mise.toml"), "[tools]\npython = \"3.14\"\n")
		target := writePythonFile(t, filepath.Join(dir, "main.py"), "x = 1\n")

		version, err := detectPython(t, registry.DetectRequest{FilePath: target})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.10" || version.Source != "pyproject.toml requires-python" {
			t.Errorf("pyproject should outrank mise: got %+v", version)
		}
	})

	t.Run(".python-version outranks mise", func(t *testing.T) {
		dir := t.TempDir()
		writePythonFile(t, filepath.Join(dir, ".python-version"), "3.12.4\n")
		writePythonFile(t, filepath.Join(dir, "mise.toml"), "[tools]\npython = \"3.14\"\n")
		target := writePythonFile(t, filepath.Join(dir, "main.py"), "x = 1\n")

		version, err := detectPython(t, registry.DetectRequest{FilePath: target})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.12" || version.Source != ".python-version" {
			t.Errorf(".python-version should outrank mise: got %+v", version)
		}
	})

	t.Run("mise outranks .tool-versions", func(t *testing.T) {
		dir := t.TempDir()
		writePythonFile(t, filepath.Join(dir, "mise.toml"), "[tools]\npython = \"3.12\"\n")
		writePythonFile(t, filepath.Join(dir, ".tool-versions"), "python 3.13\n")
		target := writePythonFile(t, filepath.Join(dir, "main.py"), "x = 1\n")

		version, err := detectPython(t, registry.DetectRequest{FilePath: target})
		if err != nil {
			t.Fatalf("Detect failed: %v", err)
		}
		if version.Version != "3.12" || version.Source != "mise.toml python" {
			t.Errorf("mise should outrank .tool-versions: got %+v", version)
		}
	})
}

func TestPythonDetectMiseAndToolVersions(t *testing.T) {
	tests := []struct {
		name         string
		miseTools    string // lines inside mise.toml's [tools] table
		toolVersions string // full .tool-versions content
		want         string
		wantSource   string
		wantFidelity registry.Fidelity
	}{
		{
			name:         "mise pin with patch",
			miseTools:    "python = \"3.12.4\"\n",
			want:         "3.12",
			wantSource:   "mise.toml python",
			wantFidelity: registry.FidelityHigh,
		},
		{
			name:         "mise fuzzy major only",
			miseTools:    "python = \"3\"\n",
			want:         "3.0",
			wantSource:   "mise.toml python",
			wantFidelity: registry.FidelityMedium,
		},
		{
			name:         "mise list takes the first entry",
			miseTools:    "python = [\"3.11\", \"3.12\"]\n",
			want:         "3.11",
			wantSource:   "mise.toml python",
			wantFidelity: registry.FidelityHigh,
		},
		{
			name:         "tool-versions pin",
			toolVersions: "python 3.12.4\n",
			want:         "3.12",
			wantSource:   ".tool-versions python",
			wantFidelity: registry.FidelityHigh,
		},
		{
			name:         "tool-versions fuzzy major only",
			toolVersions: "python 3\n",
			want:         "3.0",
			wantSource:   ".tool-versions python",
			wantFidelity: registry.FidelityMedium,
		},
		{
			name:         "tool-versions skips comments and other tools",
			toolVersions: "# managed by mise\nnode 20\npython 3.13\nruby 3.3\n",
			want:         "3.13",
			wantSource:   ".tool-versions python",
			wantFidelity: registry.FidelityHigh,
		},
		{
			name:         "mise without a python entry falls to tool-versions",
			miseTools:    "node = \"20\"\n",
			toolVersions: "python 3.12\n",
			want:         "3.12",
			wantSource:   ".tool-versions python",
			wantFidelity: registry.FidelityHigh,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			if test.miseTools != "" {
				writePythonFile(t, filepath.Join(dir, "mise.toml"), "[tools]\n"+test.miseTools)
			}
			if test.toolVersions != "" {
				writePythonFile(t, filepath.Join(dir, ".tool-versions"), test.toolVersions)
			}
			target := writePythonFile(t, filepath.Join(dir, "main.py"), "x = 1\n")

			version, err := detectPython(t, registry.DetectRequest{FilePath: target})
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}
			if version.Version != test.want {
				t.Errorf("version = %q, want %q", version.Version, test.want)
			}
			if version.Source != test.wantSource {
				t.Errorf("source = %q, want %q", version.Source, test.wantSource)
			}
			if version.Fidelity != test.wantFidelity {
				t.Errorf("fidelity = %q, want %q", version.Fidelity, test.wantFidelity)
			}
		})
	}
}

func TestPythonDetectFromWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	writePythonFile(t, filepath.Join(dir, "pyproject.toml"),
		"[project]\nname = \"app\"\nrequires-python = \">=3.10\"\n")
	t.Chdir(dir)

	version, err := detectPython(t, registry.DetectRequest{})
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}
	if version.Version != "3.10" || version.Source != "pyproject.toml requires-python" {
		t.Errorf("empty file path should anchor find-up at the working directory: got %+v", version)
	}
}

func TestPythonDetectNoManifestAnywhere(t *testing.T) {
	dir := t.TempDir()
	target := writePythonFile(t, filepath.Join(dir, "main.py"), "import sys\n")

	_, err := detectPython(t, registry.DetectRequest{FilePath: target})
	if err == nil {
		t.Fatal("Detect succeeded without any manifest, want a guidance error")
	}
	for _, want := range []string{
		"cannot determine the Python version",
		"no local-toolchain fallback",
		"--version python=3.12",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not contain %q", err.Error(), want)
		}
	}
}
