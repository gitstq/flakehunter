package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/model"
	"github.com/gitstq/flakehunter/internal/report"
)

// ScanCmd implements `flakehunter scan`: a one-off analysis of the given
// test result files. Scan is diagnostic and never fails the build by itself.
func ScanCmd(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	var (
		out         = fs.String("out", "terminal", "output format: terminal|json|markdown")
		summaryFile = fs.String("summary", "", "also write a markdown summary to this file (e.g. $GITHUB_STEP_SUMMARY)")
		record      = fs.Bool("record", false, "append this run to local history")
		dir         = fs.String("dir", defaultDir, "state directory")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	paths := fs.Args()
	if len(paths) == 0 {
		return fmt.Errorf("scan: no test result files or directories given")
	}
	results, used, err := collect(paths)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return fmt.Errorf("scan: no test results found in %s", strings.Join(paths, ", "))
	}
	run := model.Run{
		RunID:     nowRunID(),
		Timestamp: time.Now().UTC(),
		Source:    strings.Join(used, ", "),
		Commit:    commitFromEnv(),
		Results:   results,
	}
	signals := engine.RunSignals(run)
	summary := report.BuildScanSummary(run, signals)
	if *record {
		store, _, err := openStore(*dir)
		if err != nil {
			return err
		}
		if err := store.Append(run); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "flakehunter: recorded run %s (%d results)\n", run.RunID, len(results))
	}
	switch *out {
	case "json":
		if err := report.WriteJSON(os.Stdout, summary); err != nil {
			return err
		}
	case "markdown":
		report.ScanMarkdown(os.Stdout, summary)
	default:
		report.ScanTerminal(os.Stdout, summary)
	}
	if *summaryFile != "" {
		if err := writeFile(*summaryFile, func(w io.Writer) { report.ScanMarkdown(w, summary) }); err != nil {
			return err
		}
	}
	return nil
}
