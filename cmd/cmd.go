// Package cmd wires FlakeHunter subcommands to the engine.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/ingest"
	"github.com/gitstq/flakehunter/internal/model"
)

// defaultDir is the default state directory.
const defaultDir = ".flakehunter"

// nowRunID returns a sortable run identifier.
func nowRunID() string {
	return time.Now().UTC().Format("20060102T150405.000000000")
}

// commitFromEnv returns the current commit when running inside GitHub Actions.
func commitFromEnv() string {
	return os.Getenv("GITHUB_SHA")
}

// collect parses every detectable test-result file found in paths and
// re-numbers the results so their order is stable within the run.
func collect(paths []string) ([]model.TestResult, []string, error) {
	files, err := ingest.WalkPaths(paths)
	if err != nil {
		return nil, nil, err
	}
	var results []model.TestResult
	var used []string
	for _, file := range files {
		rs, format, err := ingest.ParseFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "flakehunter: skipping %s: %v\n", file, err)
			continue
		}
		if len(rs) == 0 {
			continue
		}
		used = append(used, fmt.Sprintf("%s (%s)", file, format))
		results = append(results, rs...)
	}
	for i := range results {
		results[i].Order = i
	}
	return results, used, nil
}

// openStore loads the config and opens the history store.
func openStore(dir string) (*engine.Store, engine.Config, error) {
	cfg, err := engine.LoadConfig(dir)
	if err != nil {
		return nil, cfg, err
	}
	store, err := engine.OpenStore(dir, cfg.MaxRuns)
	if err != nil {
		return nil, cfg, err
	}
	return store, cfg, nil
}

// writeFile creates path and lets write fill it.
func writeFile(path string, write func(w io.Writer)) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	write(f)
	return f.Close()
}

// splitFingerprint splits "suite#name" into its parts.
func splitFingerprint(fp string) (suite, name string) {
	if idx := strings.Index(fp, "#"); idx >= 0 {
		return fp[:idx], fp[idx+1:]
	}
	return "", fp
}
