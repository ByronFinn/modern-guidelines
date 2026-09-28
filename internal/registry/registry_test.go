package registry

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// testDataset is a minimal Dataset implementation for registration tests.
type testDataset struct {
	language   string
	axes       []string
	provenance *schema.Provenance
}

var _ Dataset = testDataset{}

func (d testDataset) Language() string                   { return d.language }
func (d testDataset) Axes() []string                     { return d.axes }
func (d testDataset) Provenance() *schema.Provenance     { return d.provenance }
func (d testDataset) Rules() []schema.Rule               { return nil }
func (d testDataset) Rule(id string) (schema.Rule, bool) { return schema.Rule{}, false }

// testDetector is a minimal Detector implementation for registration tests.
type testDetector struct {
	language string
}

var _ Detector = testDetector{}

func (d testDetector) Language() string { return d.language }

func (d testDetector) Detect(req DetectRequest) (map[string]Version, error) {
	if version, ok := req.VersionOverrides[d.language]; ok {
		return map[string]Version{d.language: {Version: version, Source: "flag", Fidelity: FidelityHigh}}, nil
	}
	return map[string]Version{d.language: {Version: "1.0", Source: "test", Fidelity: FidelityMedium}}, nil
}

func resetRegistry(t *testing.T) {
	t.Helper()
	clear := func() {
		registryMu.Lock()
		defer registryMu.Unlock()
		datasets = make(map[string]Dataset)
		detectors = make(map[string]Detector)
	}
	clear()
	t.Cleanup(clear)
}

func TestLanguageForExtension(t *testing.T) {
	tests := []struct {
		extension string
		want      string
		wantOK    bool
	}{
		{".py", "python", true},
		{".go", "go", true},
		{".ts", "typescript", true},
		{".mts", "typescript", true},
		{".cts", "typescript", true},
		{".tsx", "typescript", true},
		{".js", "javascript", true},
		{".mjs", "javascript", true},
		{"ts", "typescript", true},
		{".TS", "typescript", true},
		{".Go", "go", true},
		{" .py ", "python", true},
		{".json", "", false},
		{".rs", "", false},
		{".jsx", "", false},
		{"", "", false},
		{".", "", false},
	}

	for _, test := range tests {
		t.Run(test.extension, func(t *testing.T) {
			got, ok := LanguageForExtension(test.extension)
			if got != test.want || ok != test.wantOK {
				t.Fatalf("LanguageForExtension(%q) = (%q, %v), want (%q, %v)", test.extension, got, ok, test.want, test.wantOK)
			}
		})
	}
}

func TestLanguageForPath(t *testing.T) {
	tests := []struct {
		path   string
		want   string
		wantOK bool
	}{
		{"src/foo.ts", "typescript", true},
		{"main.go", "go", true},
		{"a/b/c/d.mts", "typescript", true},
		{"script.py", "python", true},
		{"README.md", "", false},
		{"noext", "", false},
		{"dir.py/readme", "", false},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			got, ok := LanguageForPath(test.path)
			if got != test.want || ok != test.wantOK {
				t.Fatalf("LanguageForPath(%q) = (%q, %v), want (%q, %v)", test.path, got, ok, test.want, test.wantOK)
			}
		})
	}
}

func TestSupportedExtensions(t *testing.T) {
	extensions := SupportedExtensions()
	want := []string{".cts", ".go", ".js", ".mjs", ".mts", ".py", ".ts", ".tsx"}
	if !slices.Equal(extensions, want) {
		t.Fatalf("SupportedExtensions() = %v, want %v", extensions, want)
	}
}

func TestRegisterDatasetAndLookup(t *testing.T) {
	resetRegistry(t)

	RegisterDataset(testDataset{language: "typescript", axes: []string{"typescript", "node"}})
	RegisterDataset(testDataset{language: "go"})
	RegisterDataset(testDataset{language: "python"})

	dataset, ok := DatasetForLanguage("python")
	if !ok || dataset.Language() != "python" {
		t.Fatalf("DatasetForLanguage(python) = (%#v, %v)", dataset, ok)
	}
	if _, ok := DatasetForLanguage("javascript"); ok {
		t.Fatal("javascript is a placeholder extension language and must have no dataset")
	}

	languages := []string{}
	for _, dataset := range Datasets() {
		languages = append(languages, dataset.Language())
	}
	if want := []string{"go", "python", "typescript"}; !slices.Equal(languages, want) {
		t.Fatalf("Datasets() languages = %v, want %v", languages, want)
	}
}

func TestRegisterDetectorAndLookup(t *testing.T) {
	resetRegistry(t)

	RegisterDetector(testDetector{language: "python"})
	detector, ok := DetectorForLanguage("python")
	if !ok || detector.Language() != "python" {
		t.Fatalf("DetectorForLanguage(python) = (%#v, %v)", detector, ok)
	}

	versions, err := detector.Detect(DetectRequest{
		FilePath:         "src/app.py",
		VersionOverrides: map[string]string{"python": "3.13"},
	})
	if err != nil {
		t.Fatal(err)
	}
	version, ok := versions["python"]
	if !ok || version.Version != "3.13" || version.Source != "flag" || version.Fidelity != FidelityHigh {
		t.Fatalf("Detect() = %#v", versions)
	}

	if _, ok := DetectorForLanguage("go"); ok {
		t.Fatal("unregistered detector found")
	}
}

func TestRegisterDatasetPanics(t *testing.T) {
	tests := []struct {
		name string
		run  func()
		want string
	}{
		{"nil dataset", func() { RegisterDataset(nil) }, "nil dataset"},
		{"empty language", func() { RegisterDataset(testDataset{}) }, "empty dataset language"},
		{"duplicate language", func() {
			RegisterDataset(testDataset{language: "python"})
			RegisterDataset(testDataset{language: "python"})
		}, "already registered"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resetRegistry(t)
			defer func() {
				recovered := recover()
				if recovered == nil {
					t.Fatal("RegisterDataset() did not panic")
				}
				if message := fmt.Sprint(recovered); !strings.Contains(message, test.want) {
					t.Fatalf("panic = %q, want substring %q", message, test.want)
				}
			}()
			test.run()
		})
	}
}

func TestRegisterDetectorPanics(t *testing.T) {
	tests := []struct {
		name string
		run  func()
		want string
	}{
		{"nil detector", func() { RegisterDetector(nil) }, "nil detector"},
		{"empty language", func() { RegisterDetector(testDetector{}) }, "empty detector language"},
		{"duplicate language", func() {
			RegisterDetector(testDetector{language: "go"})
			RegisterDetector(testDetector{language: "go"})
		}, "already registered"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resetRegistry(t)
			defer func() {
				recovered := recover()
				if recovered == nil {
					t.Fatal("RegisterDetector() did not panic")
				}
				if message := fmt.Sprint(recovered); !strings.Contains(message, test.want) {
					t.Fatalf("panic = %q, want substring %q", message, test.want)
				}
			}()
			test.run()
		})
	}
}
