// Command mg is the modern-guidelines language hub: a single static binary
// that infers a project's language and version axes from file paths and
// serves version-gated modern-idiom guidelines over a CLI and a built-in
// MCP server (PRD-0000).
//
// The internal/cli package runs the subcommands; importing internal/cli in
// turn loads internal/guidelines, whose init embeds and validates the
// datasets and registers them in internal/registry. The blank import of
// internal/detectors registers the per-language version detectors.
package main

import (
	"fmt"
	"os"

	"github.com/ByronFinn/modern-guidelines/internal/cli"
	_ "github.com/ByronFinn/modern-guidelines/internal/detectors"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
