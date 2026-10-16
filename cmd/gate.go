package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/gate"
	"github.com/gitstq/flakehunter/internal/quarantine"
	"github.com/gitstq/flakehunter/internal/report"
)

// errGateFailed is returned when the quality gate rejects the run.
var errGateFailed = errors.New("quality gate failed")

// GateCmd implements `flakehunter gate`: a CI quality gate that exits
// non-zero when flake thresholds are exceeded.
func GateCmd(args []string) error {
	fs := flag.NewFlagSet("gate", flag.ContinueOnError)
	var (
		dir          = fs.String("dir", defaultDir, "state directory")
		maxFlakeRate = fs.Float64("max-flake-rate", -1, "maximum flake rate in percent (default: from config)")
		maxFlaky     = fs.Int("max-flaky", -1, "maximum number of flaky tests (default: from config)")
		failOnBroken = fs.Bool("fail-on-broken", false, "fail when consistently failing tests exist")
		summaryFile  = fs.String("summary", "", "also write a markdown summary to this file")
		quiet        = fs.Bool("quiet", false, "only print violations and annotations")
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
		fmt.Println("flakehunter: no recorded history yet — nothing to gate")
		return nil
	}
	stats := engine.Analyze(runs, engine.Options{MinRuns: cfg.MinRuns})
	signals := engine.OrderSignals(runs)
	man, err := quarantine.LoadManifest(*dir)
	if err != nil {
		return err
	}
	quarantined := map[string]bool{}
	for _, e := range man.Entries {
		quarantined[e.Fingerprint] = true
	}
	opts := gate.Options{
		MaxFlakeRate: cfg.MaxFlakeRate,
		MaxFlaky:     cfg.MaxFlaky,
		FailOnBroken: *failOnBroken,
		Quarantined:  quarantined,
	}
	if *maxFlakeRate >= 0 {
		opts.MaxFlakeRate = *maxFlakeRate
	}
	if *maxFlaky >= 0 {
		opts.MaxFlaky = *maxFlaky
	}
	res := gate.Evaluate(stats, opts)
	if !*quiet {
		report.Terminal(os.Stdout, stats, signals, len(runs))
	}
	for _, annotation := range res.Annotations {
		fmt.Println(annotation)
	}
	if *summaryFile != "" {
		summary := report.BuildSummary(len(runs), stats, signals)
		if err := writeFile(*summaryFile, func(w io.Writer) { report.Markdown(w, summary) }); err != nil {
			return err
		}
	}
	if res.Pass {
		fmt.Printf("✅ gate passed — %d/%d analyzed tests flaky (%.1f%%), %d broken\n",
			res.FlakyCount, res.Analyzed, res.FlakeRate, res.BrokenCount)
		return nil
	}
	for _, v := range res.Violations {
		fmt.Fprintf(os.Stderr, "❌ gate violation: %s\n", v)
	}
	return errGateFailed
}
