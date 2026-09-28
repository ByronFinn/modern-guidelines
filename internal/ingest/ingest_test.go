package ingest

import (
	"bytes"
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/ByronFinn/modern-guidelines/internal/schema"
)

// lockedUpstreamRules is the rule count of the absorbed snapshot (54 rules at
// upstream v1.1.1, commit 019b45e2). Update this constant together with
// `make ingest` when resyncing to a new upstream minor release (ADR 0004).
const lockedUpstreamRules = 54

// Test-time paths: tests run with the package directory as working
// directory, so both point through the repository root.
const (
	upstreamPath = "../../" + UpstreamPath
	datasetPath  = "../guidelines/data/go.json"
)

func readUpstream(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(upstreamPath)
	if err != nil {
		t.Fatalf("read upstream snapshot: %v", err)
	}
	return data
}

func readDataset(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(datasetPath)
	if err != nil {
		t.Fatalf("read committed dataset: %v", err)
	}
	return data
}

// TestGoDatasetInSync fails if internal/guidelines/data/go.json has drifted
// from a fresh conversion of the locked upstream snapshot. Run `make ingest`
// to regenerate it. Byte-level on purpose, mirroring the reference project's
// TestFeaturesMarkdownInSync pattern.
func TestGoDatasetInSync(t *testing.T) {
	want, err := Convert(readUpstream(t))
	if err != nil {
		t.Fatalf("convert upstream snapshot: %v", err)
	}
	got := readDataset(t)
	if !bytes.Equal(got, want) {
		t.Fatal("internal/guidelines/data/go.json is out of date; run `make ingest`")
	}
}

// TestGoDatasetPreservesEveryUpstreamRule is the double-entry check behind
// the drift test: it restates the upstream-to-schema mapping independently and
// verifies the committed dataset rule by rule, so a mapping change cannot
// pass merely by regenerating the file. It also locks the rule count (54 at
// the absorbed commit; ADR 0004).
func TestGoDatasetPreservesEveryUpstreamRule(t *testing.T) {
	var upstream []upstreamRule
	decoder := json.NewDecoder(bytes.NewReader(readUpstream(t)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&upstream); err != nil {
		t.Fatalf("parse upstream snapshot: %v", err)
	}
	if len(upstream) != lockedUpstreamRules {
		t.Fatalf("upstream snapshot has %d rules, want the locked %d (resync in progress? update the lock)", len(upstream), lockedUpstreamRules)
	}

	dataset, err := Validate(readDataset(t))
	if err != nil {
		t.Fatalf("committed dataset is not a valid external schema v2 dataset: %v", err)
	}
	if dataset.Language != "go" {
		t.Fatalf("dataset language = %q, want %q", dataset.Language, "go")
	}
	if len(dataset.Rules) != len(upstream) {
		t.Fatalf("committed dataset has %d rules, upstream has %d: %d rules lost", len(dataset.Rules), len(upstream), len(upstream)-len(dataset.Rules))
	}

	// Provenance is restated literally (not via the package constants) so a
	// change to either side trips this test.
	wantProvenance := schema.Provenance{
		Source:  "github.com/JetBrains/go-modern-guidelines",
		License: "Apache-2.0",
		Commit:  "019b45e2b2f80d7f5c1e28bd4d35f3f0fbcf9c9c",
		Date:    "2026-09-29",
	}
	if dataset.Provenance == nil || *dataset.Provenance != wantProvenance {
		t.Fatalf("dataset provenance = %#v, want %#v", dataset.Provenance, wantProvenance)
	}

	for i, rule := range upstream {
		got := dataset.Rules[i]
		if got.ID != rule.ID {
			t.Fatalf("rule %d: id = %q, want upstream order id %q (rules reordered or lost)", i, got.ID, rule.ID)
		}
		if got.SinceVersion != rule.SinceVersion {
			t.Fatalf("rule %q: since_version = %q, want %q", got.ID, got.SinceVersion, rule.SinceVersion)
		}
		if got.Category != rule.Category {
			t.Fatalf("rule %q: category = %q, want %q", got.ID, got.Category, rule.Category)
		}
		if got.Impact != rule.Impact {
			t.Fatalf("rule %q: impact = %q, want %q", got.ID, got.Impact, rule.Impact)
		}
		if got.Guideline != rule.Guideline {
			t.Fatalf("rule %q: guideline mismatch", got.ID)
		}
		if got.Details != rule.Details {
			t.Fatalf("rule %q: details mismatch", got.ID)
		}
		// The mapping restated independently of ConvertDataset: only
		// modernizer rules carry an autofix, every other rule an explicit
		// null (the key itself is mandatory).
		wantAutofix := schema.Autofix{Present: true}
		if rule.Modernizer {
			wantAutofix = schema.Autofix{Present: true, Tool: "gopls", Rule: "modernize"}
		}
		if got.Autofix != wantAutofix {
			t.Fatalf("rule %q: autofix = %v, want %v", got.ID, got.Autofix, wantAutofix)
		}
		if len(got.References) != 0 {
			t.Fatalf("rule %q: references = %v, want none (omitted for absorbed Go data)", got.ID, got.References)
		}
		if !equalExamples(got.Examples, rule.Examples) {
			t.Fatalf("rule %q: examples mismatch (before/after line arrays must be preserved verbatim)", got.ID)
		}
	}
}

func equalExamples(got, want []schema.Example) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if !slices.Equal(got[i].Before, want[i].Before) || !slices.Equal(got[i].After, want[i].After) {
			return false
		}
	}
	return true
}

// TestDriftCheckDetectsMutations proves the byte comparison in
// TestGoDatasetInSync is sensitive: every injected upstream change and every
// conversion-mapping change must yield output that no longer equals the
// committed dataset (or be rejected outright), so a stale data/go.json can
// never pass as fresh.
func TestDriftCheckDetectsMutations(t *testing.T) {
	committed := readDataset(t)
	upstream := readUpstream(t)

	mutateFirst := func(t *testing.T, mutate func(*upstreamRule)) []byte {
		t.Helper()
		rules := decodeUpstream(t, upstream)
		mutate(&rules[0])
		encoded, err := json.Marshal(rules)
		if err != nil {
			t.Fatal(err)
		}
		return encoded
	}

	mutations := map[string]func(t *testing.T) []byte{
		"upstream guideline text": func(t *testing.T) []byte {
			return mutateFirst(t, func(rule *upstreamRule) { rule.Guideline += " (edited upstream)" })
		},
		"upstream modernizer flag": func(t *testing.T) []byte {
			return mutateFirst(t, func(rule *upstreamRule) { rule.Modernizer = !rule.Modernizer })
		},
		"upstream example line": func(t *testing.T) []byte {
			return mutateFirst(t, func(rule *upstreamRule) { rule.Examples[0].After[0] += " // edited" })
		},
		"upstream since_version": func(t *testing.T) []byte {
			return mutateFirst(t, func(rule *upstreamRule) { rule.SinceVersion = "1.26" })
		},
		"upstream rule order swapped": func(t *testing.T) []byte {
			t.Helper()
			rules := decodeUpstream(t, upstream)
			rules[0], rules[1] = rules[1], rules[0] // both 1.27: still newest-first
			encoded, err := json.Marshal(rules)
			if err != nil {
				t.Fatal(err)
			}
			return encoded
		},
		"upstream rule dropped": func(t *testing.T) []byte {
			t.Helper()
			rules := decodeUpstream(t, upstream)
			encoded, err := json.Marshal(rules[:len(rules)-1])
			if err != nil {
				t.Fatal(err)
			}
			return encoded
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			converted, err := Convert(mutate(t))
			if err == nil && bytes.Equal(converted, committed) {
				t.Fatalf("mutation %q produced the committed dataset; the drift check would not fail", name)
			}
		})
	}

	t.Run("conversion mapping changed", func(t *testing.T) {
		wrong, err := ConvertDataset(upstream)
		if err != nil {
			t.Fatal(err)
		}
		for i := range wrong.Rules { // e.g. autofix dropped for modernizable rules
			wrong.Rules[i].Autofix = schema.Autofix{Present: true}
		}
		encoded, err := Marshal(wrong)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(encoded, committed) {
			t.Fatal("changed mapping produced the committed dataset; the drift check would not fail")
		}
	})
}

func decodeUpstream(t *testing.T, data []byte) []upstreamRule {
	t.Helper()
	var rules []upstreamRule
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&rules); err != nil {
		t.Fatalf("parse upstream snapshot: %v", err)
	}
	return rules
}
