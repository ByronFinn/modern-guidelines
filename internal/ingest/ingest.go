// Package ingest converts the Apache-2.0 JetBrains go-modern-guidelines
// snapshot committed under internal/testdata/upstream/ into this
// repository's schema v2 Go dataset (PRD-0000 requirement 3, ADR 0004).
//
// The conversion is a pure, deterministic function of the upstream JSON so
// the drift test in this package can re-execute it and byte-compare against
// the committed internal/guidelines/data/go.json, mirroring the reference
// project's TestFeaturesMarkdownInSync pattern. scripts/ingest-jetbrains-go
// is a thin wrapper around this package; both share the exact same code, so
// a regeneration and a test re-run can never disagree.
//
// Field mapping (upstream guidelines.json -> schema v2 data/go.json):
//
//	id            -> id                       verbatim
//	since_version -> since_version            verbatim; newest-first upstream
//	                                         order is preserved (schema checks it)
//	modernizer    -> autofix                  true -> "gopls:modernize",
//	                                         false -> null; the key is always
//	                                         present (mandatory in schema v2)
//	category      -> category                 verbatim
//	impact        -> impact                   verbatim
//	guideline     -> guideline                verbatim
//	details       -> details                  verbatim
//	examples      -> examples                 verbatim, before/after line
//	                                         arrays preserved
//	references    -> (omitted)                optional in schema v2; the
//	                                         upstream data carries none
//
// The dataset header carries the locked provenance below.
package ingest

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// Locked conversion inputs and outputs, as paths relative to the repository
// root (the ingest script runs from there; tests resolve them from this
// package's directory).
const (
	// UpstreamPath is the committed copy of the JetBrains snapshot (the
	// local copy of the upstream commit locked by UpstreamCommit); its
	// provenance is documented in internal/testdata/upstream/README.md.
	UpstreamPath = "internal/testdata/upstream/go-guidelines.json"

	// OutputPath is the schema v2 dataset that `make ingest` overwrites.
	OutputPath = "internal/guidelines/data/go.json"
)

// Provenance of the absorbed dataset (ADR 0004: one-time snapshot, commit
// locked, resynced per upstream minor release together with the drift test's
// locked rule count).
const (
	UpstreamSource  = "github.com/JetBrains/go-modern-guidelines"
	UpstreamLicense = "Apache-2.0"
	UpstreamCommit  = "019b45e2b2f80d7f5c1e28bd4d35f3f0fbcf9c9c"
	IngestDate      = "2026-09-29"
)

// datasetLanguage is the single language this package ingests.
const datasetLanguage = "go"

// autofixModernize marks rules the gopls modernize analyzer can rewrite;
// every other rule keeps an explicit null autofix.
const autofixTool, autofixRule = "gopls", "modernize"

// upstreamRule is one entry of the JetBrains guidelines.json document
// (the upstream repository's internal/schema). Decoding is strict: a new
// upstream field fails the conversion instead of being silently dropped, so
// schema-mapping changes surface during the ADR 0004 diff review.
type upstreamRule struct {
	ID           string           `json:"id"`
	SinceVersion string           `json:"since_version"`
	Modernizer   bool             `json:"modernizer"`
	Category     string           `json:"category"`
	Impact       string           `json:"impact"`
	Guideline    string           `json:"guideline"`
	Details      string           `json:"details"`
	Examples     []schema.Example `json:"examples"`
}

// ConvertDataset transforms the upstream guidelines JSON into a schema v2
// dataset. It performs only the mechanical field mapping; full validation is
// done by Validate on the encoded result.
func ConvertDataset(upstream []byte) (schema.Dataset, error) {
	var rules []upstreamRule
	decoder := json.NewDecoder(bytes.NewReader(upstream))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&rules); err != nil {
		return schema.Dataset{}, fmt.Errorf("parse upstream guidelines JSON: %w", err)
	}
	if len(rules) == 0 {
		return schema.Dataset{}, fmt.Errorf("upstream guidelines are empty")
	}

	dataset := schema.Dataset{
		Language: datasetLanguage,
		Provenance: &schema.Provenance{
			Source:  UpstreamSource,
			License: UpstreamLicense,
			Commit:  UpstreamCommit,
			Date:    IngestDate,
		},
		Rules: make([]schema.Rule, 0, len(rules)),
	}
	for _, rule := range rules {
		autofix := schema.Autofix{Present: true} // explicit null unless modernizable
		if rule.Modernizer {
			autofix.Tool, autofix.Rule = autofixTool, autofixRule
		}
		dataset.Rules = append(dataset.Rules, schema.Rule{
			ID:           rule.ID,
			SinceVersion: rule.SinceVersion,
			Autofix:      autofix,
			Category:     rule.Category,
			Impact:       rule.Impact,
			Guideline:    rule.Guideline,
			Details:      rule.Details,
			Examples:     rule.Examples,
		})
	}
	return dataset, nil
}

// Marshal renders a dataset as the deterministic go.json document: two-space
// indentation, no HTML escaping (upstream snippets contain <, >, and &, which
// must stay literal for reviewable diffs), trailing newline. The drift test
// relies on this byte-stable form.
func Marshal(dataset schema.Dataset) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dataset); err != nil {
		return nil, fmt.Errorf("encode dataset JSON: %w", err)
	}
	return buffer.Bytes(), nil
}

// Validate decodes and validates an encoded dataset with the external-dataset
// rules (complete provenance required). Every conversion must pass it before
// being committed.
func Validate(encoded []byte) (schema.Dataset, error) {
	return schema.ParseExternal(encoded)
}

// Convert is the full pipeline: map the upstream JSON, encode it
// deterministically, and validate the encoded result. Its output is exactly
// what the ingest script writes and what the committed dataset must contain.
func Convert(upstream []byte) ([]byte, error) {
	dataset, err := ConvertDataset(upstream)
	if err != nil {
		return nil, err
	}
	encoded, err := Marshal(dataset)
	if err != nil {
		return nil, err
	}
	if _, err := Validate(encoded); err != nil {
		return nil, fmt.Errorf("converted dataset does not pass schema v2 validation: %w", err)
	}
	return encoded, nil
}
