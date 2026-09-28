// Command ingest-jetbrains-go converts the vendored, commit-locked
// JetBrains go-modern-guidelines snapshot (Apache-2.0) into the schema v2
// dataset at internal/guidelines/data/go.json, overwriting it in place
// (PRD-0000 requirement 3, ADR 0004).
//
// Run it from the repository root:
//
//	go run ./scripts/ingest-jetbrains-go
//
// or via `make ingest`. The conversion logic lives in internal/ingest and is
// byte-locked by that package's drift test, so regenerating the dataset can
// never silently diverge from the committed one.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ByronFinn/modern-guidelines/internal/ingest"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("ingest-jetbrains-go: ")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	upstream, err := os.ReadFile(ingest.UpstreamPath)
	if err != nil {
		return fmt.Errorf("read upstream snapshot: %w", err)
	}

	dataset, err := ingest.ConvertDataset(upstream)
	if err != nil {
		return err
	}
	encoded, err := ingest.Marshal(dataset)
	if err != nil {
		return err
	}
	if _, err := ingest.Validate(encoded); err != nil {
		return fmt.Errorf("refusing to write an invalid dataset: %w", err)
	}

	if err := os.WriteFile(ingest.OutputPath, encoded, 0o644); err != nil {
		return fmt.Errorf("write dataset: %w", err)
	}
	fmt.Printf("wrote %d rules to %s (source %s@%s, %s)\n",
		len(dataset.Rules), ingest.OutputPath, ingest.UpstreamSource,
		ingest.UpstreamCommit, ingest.UpstreamLicense)
	return nil
}
