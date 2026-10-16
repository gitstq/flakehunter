package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/quarantine"
	"github.com/gitstq/flakehunter/internal/report"
)

// ReportCmd implements `flakehunter report`: flake statistics from the
// locally recorded run history.
func ReportCmd(args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	var (
		dir         = fs.String("dir", defaultDir, "state directory")
		out         = fs.String("out", "terminal", "output format: terminal|json|markdown")
		summaryFile = fs.String("summary", "", "also write a markdown summary to this file")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	store, cfg, err := openStore(*dir)
	if err != nil {
		return err
	}
	runs, err := store.Load()
	if err != nil {
		return err
	}
	if len(runs) == 0 {
		return fmt.Errorf("report: no history in %s — run `flakehunter record <files>` first", *dir)
	}
	stats := engine.Analyze(runs, engine.Options{MinRuns: cfg.MinRuns})
	signals := engine.OrderSignals(runs)
	man, err := quarantine.LoadManifest(*dir)
	if err != nil {
		return err
	}
	for i := range stats {
		stats[i].Quarantined = man.IsQuarantined(stats[i].Fingerprint)
	}
	summary := report.BuildSummary(len(runs), stats, signals)
	switch *out {
	case "json":
		if err := report.WriteJSON(os.Stdout, summary); err != nil {
			return err
		}
	case "markdown":
		report.Markdown(os.Stdout, summary)
	default:
		report.Terminal(os.Stdout, stats, signals, len(runs))
	}
	if *summaryFile != "" {
		if err := writeFile(*summaryFile, func(w io.Writer) { report.Markdown(w, summary) }); err != nil {
			return err
		}
	}
	return nil
}
