// Package schema defines and validates the modern-guidelines dataset format
// (PRD-0000 requirements 1 and 2).
//
// A dataset is one JSON document per language at internal/guidelines/data/<lang>.json:
//
//	{
//	  "language": "typescript",
//	  "axes": ["typescript", "node"],          // optional; required when >1 axis
//	  "provenance": { ... },                   // required for externally absorbed datasets
//	  "rules": [ ... ]
//	}
//
// Versions ("since_version" and detector results) use dot-separated numeric
// form: major[.minor] — "1.16", "5.4", or a bare major such as "18" for axes
// gated at major granularity (e.g. the node axis; PRD multi-axis example).
// Rules within one axis must be ordered newest-first; this package is the
// single authority for parsing and comparing that format.
package schema

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// maxReferences is the upper bound on optional short references per rule
// (PRD requirement 2: "references?（≤2 条短引用，可选）").
const maxReferences = 2

// Example contains the before and after snippets of a rule example.
type Example struct {
	Before []string `json:"before"`
	After  []string `json:"after"`
}

// Autofix names the tool and rule that can mechanically rewrite a violation,
// in "tool:rule" form (for example "ruff:UP007" or "gopls:modernize"), or is
// null when no such fix exists. The JSON key itself is mandatory: authors
// must declare either a mapping or an explicit null.
//
// The custom JSON handling distinguishes a missing "autofix" key (invalid)
// from an explicit null (valid); validation of the "tool:rule" form happens
// in validate so error messages can name the offending rule.
type Autofix struct {
	// Present reports whether the "autofix" key appeared in the JSON object.
	Present bool
	// Tool is the part before the colon, e.g. "ruff". Empty for null.
	Tool string
	// Rule is the part after the colon, e.g. "UP007". Empty for null.
	Rule string
}

// UnmarshalJSON implements json.Unmarshaler. It records key presence and
// treats the JSON literal null as an explicit, valid "no autofix".
func (a *Autofix) UnmarshalJSON(data []byte) error {
	a.Present = true
	if string(data) == "null" {
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf(`autofix must be a JSON string ("tool:rule") or null`)
	}
	a.Tool, a.Rule, _ = strings.Cut(value, ":")
	return nil
}

// IsNull reports whether the autofix is an explicit null (no mechanical fix).
func (a Autofix) IsNull() bool {
	return a.Tool == "" && a.Rule == ""
}

// MarshalJSON implements json.Marshaler. It writes null for the explicit
// no-autofix value and "tool:rule" otherwise. A zero Autofix (key absent)
// also marshals as null; such a rule only exists before validation.
func (a Autofix) MarshalJSON() ([]byte, error) {
	if a.IsNull() {
		return []byte("null"), nil
	}
	return json.Marshal(a.Tool + ":" + a.Rule)
}

// String renders the "tool:rule" form, or the empty string for null.
func (a Autofix) String() string {
	if a.IsNull() {
		return ""
	}
	return a.Tool + ":" + a.Rule
}

// Provenance records where an externally absorbed dataset came from
// (PRD requirement 1: source/license/upstream commit/date).
type Provenance struct {
	Source  string `json:"source"`
	License string `json:"license"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Rule is one modern-idiom guideline entry (PRD requirement 2).
type Rule struct {
	ID           string    `json:"id"`
	SinceVersion string    `json:"since_version"`
	Axis         string    `json:"axis,omitempty"`
	Autofix      Autofix   `json:"autofix"`
	Category     string    `json:"category"`
	Impact       string    `json:"impact"`
	Guideline    string    `json:"guideline"`
	Details      string    `json:"details"`
	References   []string  `json:"references,omitempty"`
	Examples     []Example `json:"examples"`
}

// Dataset is one language's guideline collection, including its header.
type Dataset struct {
	Language   string      `json:"language"`
	Axes       []string    `json:"axes,omitempty"`
	Provenance *Provenance `json:"provenance,omitempty"`
	Rules      []Rule      `json:"rules"`
}

// EffectiveAxes returns the declared axes, or the language itself as the one
// implicit axis of a single-axis dataset.
func (d Dataset) EffectiveAxes() []string {
	if len(d.Axes) > 0 {
		return d.Axes
	}
	return []string{d.Language}
}

// Parse decodes and validates a dataset document.
func Parse(data []byte) (Dataset, error) {
	var dataset Dataset
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dataset); err != nil {
		return Dataset{}, fmt.Errorf("parse dataset JSON: %w", err)
	}
	if err := validate(&dataset); err != nil {
		return Dataset{}, err
	}
	return dataset, nil
}

// ParseExternal decodes and validates a dataset produced by ingesting an
// external source. In addition to Parse it requires the dataset header to
// carry complete provenance (PRD acceptance criteria: 外部吸收缺 provenance).
func ParseExternal(data []byte) (Dataset, error) {
	dataset, err := Parse(data)
	if err != nil {
		return Dataset{}, err
	}
	if dataset.Provenance == nil {
		return Dataset{}, fmt.Errorf("dataset %q is externally absorbed but has no provenance", dataset.Language)
	}
	return dataset, nil
}

// CompareVersions compares two dot-separated numeric versions (major[.minor])
// and returns -1, 0, or 1. It errors on malformed input.
func CompareVersions(left, right string) (int, error) {
	leftVersion, ok := parseDotVersion(left)
	if !ok {
		return 0, fmt.Errorf("invalid version %q: use major[.minor] with dot-separated numbers", left)
	}
	rightVersion, ok := parseDotVersion(right)
	if !ok {
		return 0, fmt.Errorf("invalid version %q: use major[.minor] with dot-separated numbers", right)
	}
	return leftVersion.compare(rightVersion), nil
}

func validate(dataset *Dataset) error {
	if dataset.Language == "" {
		return fmt.Errorf("dataset header has no language")
	}
	if !validName(dataset.Language) {
		return fmt.Errorf("invalid dataset language %q: use lowercase letters, digits, and hyphens", dataset.Language)
	}
	if len(dataset.Rules) == 0 {
		return fmt.Errorf("dataset %q has no rules", dataset.Language)
	}

	declaredAxes := make(map[string]bool, len(dataset.Axes))
	for _, axis := range dataset.Axes {
		if !validName(axis) {
			return fmt.Errorf("invalid axis %q: use lowercase letters, digits, and hyphens", axis)
		}
		if declaredAxes[axis] {
			return fmt.Errorf("duplicate axis %q", axis)
		}
		declaredAxes[axis] = true
	}
	multiAxis := len(dataset.Axes) > 1

	if dataset.Provenance != nil {
		if err := validateProvenance(dataset); err != nil {
			return err
		}
	}

	seenIDs := make(map[string]bool, len(dataset.Rules))
	previousVersion := make(map[string]dotVersion, len(declaredAxes))
	previousID := make(map[string]string, len(declaredAxes))
	for i := range dataset.Rules {
		rule := &dataset.Rules[i]
		if !validID(rule.ID) {
			return fmt.Errorf("invalid rule id %q: use lowercase letters, digits, and underscores", rule.ID)
		}
		if seenIDs[rule.ID] {
			return fmt.Errorf("rule id %q is duplicated", rule.ID)
		}
		seenIDs[rule.ID] = true

		version, ok := parseDotVersion(rule.SinceVersion)
		if !ok {
			return fmt.Errorf("invalid since_version %q for rule %q: use major[.minor] with dot-separated numbers", rule.SinceVersion, rule.ID)
		}

		axis := dataset.Language
		if multiAxis {
			if rule.Axis == "" {
				return fmt.Errorf("rule %q has no axis; the dataset declares multiple axes", rule.ID)
			}
			if !declaredAxes[rule.Axis] {
				return fmt.Errorf("rule %q has axis %q not declared in dataset axes", rule.ID, rule.Axis)
			}
			axis = rule.Axis
		} else if rule.Axis != "" {
			return fmt.Errorf("rule %q must not set axis; the dataset has a single axis", rule.ID)
		}

		if seenVersion, ordered := previousVersion[axis]; ordered && seenVersion.less(version) {
			return fmt.Errorf(
				"rules must be ordered newest-first within axis %q: rule %q with version %s appears before rule %q with version %s",
				axis, previousID[axis], seenVersion, rule.ID, rule.SinceVersion,
			)
		}
		previousVersion[axis] = version
		previousID[axis] = rule.ID

		if !rule.Autofix.Present {
			return fmt.Errorf("rule %q has no autofix key", rule.ID)
		}
		// An explicit null is valid; any non-null value must be "tool:rule".
		if !rule.Autofix.IsNull() && (rule.Autofix.Tool == "" || rule.Autofix.Rule == "" || strings.Contains(rule.Autofix.Rule, ":")) {
			return fmt.Errorf("rule %q has invalid autofix %q: use \"tool:rule\"", rule.ID, rule.Autofix.String())
		}
		if rule.Category == "" {
			return fmt.Errorf("rule %q has no category", rule.ID)
		}
		if !validImpact(rule.Impact) {
			return fmt.Errorf("rule %q has invalid impact %q", rule.ID, rule.Impact)
		}
		if strings.TrimSpace(rule.Guideline) == "" {
			return fmt.Errorf("rule %q has no guideline text", rule.ID)
		}
		if strings.TrimSpace(rule.Details) == "" {
			return fmt.Errorf("rule %q has no details", rule.ID)
		}
		if len(rule.References) > maxReferences {
			return fmt.Errorf("rule %q has %d references; at most %d are allowed", rule.ID, len(rule.References), maxReferences)
		}
		for referenceIndex, reference := range rule.References {
			if strings.TrimSpace(reference) == "" {
				return fmt.Errorf("rule %q reference %d is empty", rule.ID, referenceIndex+1)
			}
		}
		if len(rule.Examples) == 0 {
			return fmt.Errorf("rule %q has no examples", rule.ID)
		}
		for exampleIndex, example := range rule.Examples {
			if strings.TrimSpace(strings.Join(example.Before, "\n")) == "" {
				return fmt.Errorf("rule %q example %d has empty before snippet", rule.ID, exampleIndex+1)
			}
			if strings.TrimSpace(strings.Join(example.After, "\n")) == "" {
				return fmt.Errorf("rule %q example %d has empty after snippet", rule.ID, exampleIndex+1)
			}
		}
	}
	return nil
}

func validateProvenance(dataset *Dataset) error {
	provenance := dataset.Provenance
	for _, field := range []struct {
		name  string
		value string
	}{
		{"source", provenance.Source},
		{"license", provenance.License},
		{"commit", provenance.Commit},
		{"date", provenance.Date},
	} {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("dataset %q provenance is missing %q", dataset.Language, field.name)
		}
	}
	return nil
}

// dotVersion is a parsed major[.minor] version.
type dotVersion struct {
	major int
	minor int
}

func (v dotVersion) less(other dotVersion) bool {
	return v.compare(other) < 0
}

func (v dotVersion) compare(other dotVersion) int {
	if c := cmp.Compare(v.major, other.major); c != 0 {
		return c
	}
	return cmp.Compare(v.minor, other.minor)
}

// String renders the canonical form: major.minor, or the bare major when
// the minor part is zero.
func (v dotVersion) String() string {
	if v.minor == 0 {
		return strconv.Itoa(v.major)
	}
	return strconv.Itoa(v.major) + "." + strconv.Itoa(v.minor)
}

// parseDotVersion parses a dot-separated numeric version: major[.minor].
// A bare major (for example "18") is accepted for axes gated at major
// granularity; more than two components, signs, or non-numeric parts are
// rejected.
func parseDotVersion(version string) (dotVersion, bool) {
	majorPart, minorPart, hasMinor := strings.Cut(version, ".")
	if hasMinor && strings.Contains(minorPart, ".") {
		return dotVersion{}, false
	}
	major, ok := parseVersionNumber(majorPart)
	if !ok {
		return dotVersion{}, false
	}
	minor := 0
	if hasMinor {
		minor, ok = parseVersionNumber(minorPart)
		if !ok {
			return dotVersion{}, false
		}
	}
	return dotVersion{major: major, minor: minor}, true
}

func parseVersionNumber(part string) (int, bool) {
	if part == "" || part[0] < '0' || part[0] > '9' {
		return 0, false
	}
	number, err := strconv.Atoi(part)
	if err != nil {
		return 0, false
	}
	return number, true
}

// validID mirrors the reference project's guideline id charset: lowercase
// letters, digits, and underscores.
func validID(id string) bool {
	if id == "" {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

// validName checks language and axis names: non-empty lowercase letters,
// digits, and hyphens.
func validName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

func validImpact(impact string) bool {
	switch impact {
	case "Critical", "High", "Medium", "Low":
		return true
	default:
		return false
	}
}
