// Package registry wires the pluggable pieces of the modern-guidelines hub:
// file-extension based language identification, the per-language guideline
// datasets, and the per-language version detectors.
//
// Datasets and detectors self-register from their own packages via init(),
// for example:
//
//	func init() { registry.RegisterDetector(pythonDetector{}) }
//
// The registry never imports concrete implementations; the main package
// imports them for their registration side effects.
package registry

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// extensionToLanguage maps file extensions (with leading dot) to languages.
//
// javascript is a placeholder: its extensions are recognized, but no dataset
// is registered for it — Node idioms live on the typescript dataset's node
// axis (PRD out-of-scope note).
var extensionToLanguage = map[string]string{
	".py":  "python",
	".go":  "go",
	".ts":  "typescript",
	".mts": "typescript",
	".cts": "typescript",
	".tsx": "typescript",
	".js":  "javascript",
	".mjs": "javascript",
}

// LanguageForExtension returns the language for a file extension such as
// ".py" or "py" (case-insensitive, optional leading dot).
func LanguageForExtension(extension string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(extension))
	normalized = strings.TrimPrefix(normalized, ".")
	if normalized == "" {
		return "", false
	}
	language, ok := extensionToLanguage["."+normalized]
	return language, ok
}

// LanguageForPath returns the language for a file path based on its extension.
func LanguageForPath(path string) (string, bool) {
	return LanguageForExtension(filepath.Ext(path))
}

// SupportedExtensions returns every recognized file extension, sorted, with
// its leading dot. Callers use it to list supported languages in error
// messages for unknown extensions.
func SupportedExtensions() []string {
	return slices.Sorted(maps.Keys(extensionToLanguage))
}

// Dataset is the read-only view of one language's validated guideline
// dataset. The guidelines package provides the concrete implementation on
// top of the embedded data.
type Dataset interface {
	// Language is the dataset language, matching the "language" header field.
	Language() string
	// Axes are the effective version axes: the declared axes for multi-axis
	// datasets, or the language itself for single-axis datasets
	// (schema.Dataset.EffectiveAxes).
	Axes() []string
	// Provenance is the external-absorption provenance, or nil for native
	// datasets.
	Provenance() *schema.Provenance
	// Rules are all rules in dataset order (newest-first per axis, as
	// guaranteed by schema validation).
	Rules() []schema.Rule
	// Rule returns the rule with the given id; ids are unique per language.
	Rule(id string) (schema.Rule, bool)
}

// Fidelity grades how trustworthy a resolved version is (PRD decision D8).
type Fidelity string

const (
	// FidelityHigh marks versions from authoritative project manifests
	// (go.mod go directive, pyproject requires-python, ...).
	FidelityHigh Fidelity = "high"
	// FidelityMedium marks versions inferred from advisory or ambiguous
	// sources (engines.node, mise); the provenance output must warn on these.
	FidelityMedium Fidelity = "medium"
)

// Version is one resolved version gate for one axis, with its provenance.
type Version struct {
	// Version is the resolved version in dot-separated numeric form
	// (major[.minor]), e.g. "3.10" or "18".
	Version string
	// Source names where the version came from, e.g. "go.mod go directive".
	Source string
	// Fidelity grades the trustworthiness of the source.
	Fidelity Fidelity
}

// DetectRequest carries the inputs a detector may need to resolve versions.
// New inputs are added as fields, never as signature changes.
type DetectRequest struct {
	// FilePath anchors manifest discovery (find-up); may be empty.
	FilePath string
	// VersionOverrides maps axis names to explicitly requested versions
	// (T1: --version axis=ver, comma-separated for multiple axes).
	VersionOverrides map[string]string
}

// Detector resolves the declared version axes of one language (PRD
// requirement 5). It returns one entry per resolved axis; when no version
// can be resolved it returns an error guiding the user to declare one.
type Detector interface {
	Language() string
	Detect(req DetectRequest) (map[string]Version, error)
}

var (
	registryMu sync.RWMutex
	datasets   = make(map[string]Dataset)
	detectors  = make(map[string]Detector)
)

// RegisterDataset registers a guideline dataset for its language. It panics
// on nil datasets, empty languages, and duplicate languages — invalid wiring
// fails at process start, mirroring the reference project's init discipline.
func RegisterDataset(dataset Dataset) {
	if dataset == nil {
		panic("registry: RegisterDataset called with a nil dataset")
	}
	language := dataset.Language()
	if language == "" {
		panic("registry: RegisterDataset called with an empty dataset language")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := datasets[language]; exists {
		panic(fmt.Sprintf("registry: dataset for language %q is already registered", language))
	}
	datasets[language] = dataset
}

// DatasetForLanguage returns the registered dataset for a language.
func DatasetForLanguage(language string) (Dataset, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	dataset, ok := datasets[language]
	return dataset, ok
}

// Datasets returns every registered dataset, sorted by language.
func Datasets() []Dataset {
	registryMu.RLock()
	defer registryMu.RUnlock()
	ordered := make([]Dataset, 0, len(datasets))
	for _, language := range slices.Sorted(maps.Keys(datasets)) {
		ordered = append(ordered, datasets[language])
	}
	return ordered
}

// RegisterDetector registers a version detector for its language, with the
// same panic discipline as RegisterDataset.
func RegisterDetector(detector Detector) {
	if detector == nil {
		panic("registry: RegisterDetector called with a nil detector")
	}
	language := detector.Language()
	if language == "" {
		panic("registry: RegisterDetector called with an empty detector language")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := detectors[language]; exists {
		panic(fmt.Sprintf("registry: detector for language %q is already registered", language))
	}
	detectors[language] = detector
}

// DetectorForLanguage returns the registered detector for a language.
func DetectorForLanguage(language string) (Detector, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	detector, ok := detectors[language]
	return detector, ok
}
