package engine

import (
	"testing"
	"time"

	"github.com/gitstq/flakehunter/internal/model"
)

// runWith builds a single-result run for fingerprint fp with the given outcome.
func runWith(fp string, outcome model.Outcome, minute int) model.Run {
	return model.Run{
		RunID:     fp,
		Timestamp: time.Date(2026, 10, 16, 10, minute, 0, 0, time.UTC),
		Results:   []model.TestResult{{Fingerprint: fp, Outcome: outcome}},
	}
}

// find returns the stats entry for fp.
func find(stats []Stats, fp string) (Stats, bool) {
	for _, s := range stats {
		if s.Fingerprint == fp {
			return s, true
		}
	}
	return Stats{}, false
}

func TestAnalyzeClassifiesFlaky(t *testing.T) {
	fp := "suite#flaky_test"
	runs := []model.Run{
		runWith(fp, model.OutcomePassed, 0),
		runWith(fp, model.OutcomeFailed, 1),
		runWith(fp, model.OutcomePassed, 2),
		runWith(fp, model.OutcomeFailed, 3),
	}
	stats := Analyze(runs, Options{MinRuns: 3})
	s, ok := find(stats, fp)
	if !ok {
		t.Fatal("flaky test missing from stats")
	}
	if s.Status != StatusFlaky {
		t.Fatalf("want status flaky, got %s", s.Status)
	}
	if s.FlakeScore < 0.49 || s.FlakeScore > 0.51 {
		t.Errorf("want flake score ~0.5, got %v", s.FlakeScore)
	}
	if s.Confidence < 0.15 {
		t.Errorf("want high confidence (>=0.15) for 4 runs, got %v", s.Confidence)
	}
	if s.Passed != 2 || s.Failed != 2 || s.Total != 4 {
		t.Errorf("unexpected counters: %+v", s)
	}
}

func TestAnalyzeStableCases(t *testing.T) {
	runs := []model.Run{
		runWith("suite#green", model.OutcomePassed, 0),
		runWith("suite#green", model.OutcomePassed, 1),
		runWith("suite#green", model.OutcomePassed, 2),
		runWith("suite#red", model.OutcomeFailed, 0),
		runWith("suite#red", model.OutcomeError, 1),
		runWith("suite#red", model.OutcomeFailed, 2),
		runWith("suite#new", model.OutcomePassed, 0),
	}
	stats := Analyze(runs, Options{MinRuns: 3})
	green, _ := find(stats, "suite#green")
	if green.Status != StatusStablePass {
		t.Errorf("green: want stable_pass, got %s", green.Status)
	}
	red, _ := find(stats, "suite#red")
	if red.Status != StatusStableFail {
		t.Errorf("red: want stable_fail, got %s", red.Status)
	}
	fresh, _ := find(stats, "suite#new")
	if fresh.Status != StatusInsufficient {
		t.Errorf("new: want insufficient, got %s", fresh.Status)
	}
}

func TestWilsonLowerBound(t *testing.T) {
	// 1 of 2 flips is weak evidence; 5 of 10 is strong evidence.
	if got := wilsonLower(1, 2); got > 0.1 {
		t.Errorf("wilsonLower(1,2) too high: %v", got)
	}
	if got := wilsonLower(5, 10); got < 0.15 {
		t.Errorf("wilsonLower(5,10) too low: %v", got)
	}
	if got := wilsonLower(0, 0); got != 0 {
		t.Errorf("wilsonLower(0,0) want 0, got %v", got)
	}
}

func TestAnalyzeMedianDuration(t *testing.T) {
	base := time.Date(2026, 10, 16, 10, 0, 0, 0, time.UTC)
	mk := func(d float64) model.TestResult {
		return model.TestResult{Fingerprint: "suite#timed", Outcome: model.OutcomePassed, DurationMS: d}
	}
	runs := []model.Run{
		{RunID: "r1", Timestamp: base, Results: []model.TestResult{mk(10), mk(20), mk(30)}},
		{RunID: "r2", Timestamp: base.Add(time.Minute), Results: []model.TestResult{mk(40)}},
	}
	stats := Analyze(runs, Options{MinRuns: 2})
	s, _ := find(stats, "suite#timed")
	if s.MedianMS != 25 {
		t.Errorf("median: want 25, got %v", s.MedianMS)
	}
	if s.MaxMS != 40 {
		t.Errorf("max: want 40, got %v", s.MaxMS)
	}
}
