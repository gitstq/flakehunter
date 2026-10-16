package engine

import (
	"strings"
	"testing"
	"time"

	"github.com/gitstq/flakehunter/internal/model"
)

func TestRunSignalsRerunFlip(t *testing.T) {
	run := model.Run{
		RunID:     "r1",
		Timestamp: time.Now().UTC(),
		Results: []model.TestResult{
			{Fingerprint: "s#a", Outcome: model.OutcomeFailed, Order: 0},
			{Fingerprint: "s#a", Outcome: model.OutcomePassed, Order: 1},
			{Fingerprint: "s#b", Outcome: model.OutcomePassed, Order: 2},
		},
	}
	signals := RunSignals(run)
	if len(signals) != 1 {
		t.Fatalf("want 1 signal, got %d: %+v", len(signals), signals)
	}
	if signals[0].Type != SignalRerunFlip || signals[0].Fingerprint != "s#a" {
		t.Errorf("unexpected signal: %+v", signals[0])
	}
}

func TestRunSignalsTimingSpike(t *testing.T) {
	run := model.Run{
		RunID:     "r1",
		Timestamp: time.Now().UTC(),
		Results: []model.TestResult{
			{Fingerprint: "s#slow", Outcome: model.OutcomePassed, DurationMS: 10, Order: 0},
			{Fingerprint: "s#slow", Outcome: model.OutcomePassed, DurationMS: 12, Order: 1},
			{Fingerprint: "s#slow", Outcome: model.OutcomePassed, DurationMS: 500, Order: 2},
		},
	}
	signals := RunSignals(run)
	found := false
	for _, sig := range signals {
		if sig.Type == SignalTimingSpike && sig.Fingerprint == "s#slow" {
			found = true
		}
	}
	if !found {
		t.Errorf("want a timing-spike signal for s#slow, got %+v", signals)
	}
}

func TestRunSignalsCleanRun(t *testing.T) {
	run := model.Run{
		RunID:     "r1",
		Timestamp: time.Now().UTC(),
		Results: []model.TestResult{
			{Fingerprint: "s#a", Outcome: model.OutcomePassed, DurationMS: 10, Order: 0},
			{Fingerprint: "s#b", Outcome: model.OutcomePassed, DurationMS: 11, Order: 1},
		},
	}
	if signals := RunSignals(run); len(signals) != 0 {
		t.Errorf("want no signals, got %+v", signals)
	}
}

func TestOrderSignalsCorrelation(t *testing.T) {
	base := time.Date(2026, 10, 16, 10, 0, 0, 0, time.UTC)
	mkRun := func(id string, minute int, pred, bFP string, bFails bool) model.Run {
		outcome := model.OutcomePassed
		if bFails {
			outcome = model.OutcomeFailed
		}
		return model.Run{
			RunID:     id,
			Timestamp: base.Add(time.Duration(minute) * time.Minute),
			Results: []model.TestResult{
				{Fingerprint: pred, Outcome: model.OutcomePassed, Order: 0},
				{Fingerprint: bFP, Outcome: outcome, Order: 1},
			},
		}
	}
	runs := []model.Run{
		mkRun("r1", 0, "s#a", "s#b", true),
		mkRun("r2", 1, "s#a", "s#b", true),
		mkRun("r3", 2, "s#c", "s#b", false),
		mkRun("r4", 3, "s#c", "s#b", false),
	}
	signals := OrderSignals(runs)
	if len(signals) != 1 {
		t.Fatalf("want 1 order signal, got %d: %+v", len(signals), signals)
	}
	if signals[0].Fingerprint != "s#b" {
		t.Errorf("want signal for s#b, got %s", signals[0].Fingerprint)
	}
	if !strings.Contains(signals[0].Detail, "s#a") {
		t.Errorf("detail should name the predecessor: %q", signals[0].Detail)
	}
}

func TestOrderSignalsNoCorrelation(t *testing.T) {
	base := time.Date(2026, 10, 16, 10, 0, 0, 0, time.UTC)
	mkRun := func(id string, minute int, pred string, bFails bool) model.Run {
		outcome := model.OutcomePassed
		if bFails {
			outcome = model.OutcomeFailed
		}
		return model.Run{
			RunID:     id,
			Timestamp: base.Add(time.Duration(minute) * time.Minute),
			Results: []model.TestResult{
				{Fingerprint: pred, Outcome: model.OutcomePassed, Order: 0},
				{Fingerprint: "s#b", Outcome: outcome, Order: 1},
			},
		}
	}
	// Failures alternate between two predecessors: no strong correlation.
	runs := []model.Run{
		mkRun("r1", 0, "s#a", true),
		mkRun("r2", 1, "s#a", false),
		mkRun("r3", 2, "s#c", true),
		mkRun("r4", 3, "s#c", false),
	}
	if signals := OrderSignals(runs); len(signals) != 0 {
		t.Errorf("want no order signals, got %+v", signals)
	}
}

func TestPhiCoefficient(t *testing.T) {
	if got := phiCoefficient(2, 0, 0, 2); got != 1 {
		t.Errorf("perfect correlation want 1, got %v", got)
	}
	if got := phiCoefficient(1, 1, 1, 1); got != 0 {
		t.Errorf("no correlation want 0, got %v", got)
	}
	if got := phiCoefficient(0, 0, 0, 0); got != 0 {
		t.Errorf("degenerate table want 0, got %v", got)
	}
}
