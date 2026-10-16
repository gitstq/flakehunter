package engine

import (
	"fmt"
	"math"
	"sort"

	"github.com/gitstq/flakehunter/internal/model"
)

// SignalType classifies a suspicious pattern that suggests flakiness.
type SignalType string

// Signal types detected by FlakeHunter.
const (
	SignalRerunFlip   SignalType = "rerun-flip"
	SignalTimingSpike SignalType = "timing-spike"
	SignalOrderDep    SignalType = "order-dependent"
)

// Signal is a suspicious pattern observed in one run or across history.
type Signal struct {
	Type        SignalType `json:"type"`
	Fingerprint string     `json:"fingerprint"`
	Detail      string     `json:"detail"`
}

// RunSignals detects anomalies inside a single run, without any history.
// It catches the two strongest single-run flake smells: a test that both
// passed and failed in the same run (retry/rerun), and a test whose duration
// is a gross outlier compared to its own median in that run.
func RunSignals(run model.Run) []Signal {
	var signals []Signal
	outcomes := map[string]map[model.Outcome]int{}
	durations := map[string][]float64{}
	for _, res := range run.Results {
		if outcomes[res.Fingerprint] == nil {
			outcomes[res.Fingerprint] = map[model.Outcome]int{}
		}
		outcomes[res.Fingerprint][res.Outcome]++
		if res.DurationMS > 0 {
			durations[res.Fingerprint] = append(durations[res.Fingerprint], res.DurationMS)
		}
	}
	for fp, counts := range outcomes {
		if counts[model.OutcomePassed] > 0 &&
			(counts[model.OutcomeFailed] > 0 || counts[model.OutcomeError] > 0) {
			signals = append(signals, Signal{
				Type:        SignalRerunFlip,
				Fingerprint: fp,
				Detail:      "passed and failed within the same run (retry or rerun)",
			})
		}
	}
	for fp, ds := range durations {
		if len(ds) < 2 {
			continue
		}
		sorted := append([]float64{}, ds...)
		sort.Float64s(sorted)
		med := median(sorted)
		if med <= 0 {
			continue
		}
		for _, d := range ds {
			if d > med*4 && d-med > 50 {
				signals = append(signals, Signal{
					Type:        SignalTimingSpike,
					Fingerprint: fp,
					Detail:      fmt.Sprintf("duration %.0fms is over 4x the run median %.0fms", d, med),
				})
				break
			}
		}
	}
	sortSignals(signals)
	return signals
}

// OrderSignals finds tests whose failures correlate with a specific
// predecessor test, using the phi coefficient over a 2x2 contingency table
// built from run history. These are order-dependent flakes.
func OrderSignals(runs []model.Run) []Signal {
	type pairKey struct{ fp, pred string }
	type pairCounts struct {
		failWithPred int
		passWithPred int
	}
	type totals struct{ fails, passes int }

	pairs := map[pairKey]*pairCounts{}
	testTotals := map[string]*totals{}
	for _, run := range runs {
		byOrder := map[int]model.TestResult{}
		for _, res := range run.Results {
			byOrder[res.Order] = res
		}
		for _, res := range run.Results {
			failed := res.Outcome == model.OutcomeFailed || res.Outcome == model.OutcomeError
			passed := res.Outcome == model.OutcomePassed
			if !failed && !passed {
				continue
			}
			t := testTotals[res.Fingerprint]
			if t == nil {
				t = &totals{}
				testTotals[res.Fingerprint] = t
			}
			if failed {
				t.fails++
			} else {
				t.passes++
			}
			if res.Order == 0 {
				continue
			}
			pred, ok := byOrder[res.Order-1]
			if !ok || pred.Fingerprint == res.Fingerprint {
				continue
			}
			key := pairKey{fp: res.Fingerprint, pred: pred.Fingerprint}
			c := pairs[key]
			if c == nil {
				c = &pairCounts{}
				pairs[key] = c
			}
			if failed {
				c.failWithPred++
			} else {
				c.passWithPred++
			}
		}
	}
	var signals []Signal
	for key, c := range pairs {
		t := testTotals[key.fp]
		if t == nil || c.failWithPred < 2 {
			continue
		}
		failWithout := t.fails - c.failWithPred
		passWithout := t.passes - c.passWithPred
		phi := phiCoefficient(float64(c.failWithPred), float64(failWithout), float64(c.passWithPred), float64(passWithout))
		if phi >= 0.5 {
			signals = append(signals, Signal{
				Type:        SignalOrderDep,
				Fingerprint: key.fp,
				Detail:      fmt.Sprintf("fails mostly after %q (phi=%.2f)", key.pred, phi),
			})
		}
	}
	sortSignals(signals)
	if len(signals) > 20 {
		signals = signals[:20]
	}
	return signals
}

// phiCoefficient computes the phi correlation of a 2x2 contingency table.
func phiCoefficient(a, b, c, d float64) float64 {
	numerator := math.Abs(a*d - b*c)
	denominator := math.Sqrt((a + b) * (c + d) * (a + c) * (b + d))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// sortSignals orders signals deterministically for stable output.
func sortSignals(signals []Signal) {
	sort.Slice(signals, func(i, j int) bool {
		if signals[i].Fingerprint != signals[j].Fingerprint {
			return signals[i].Fingerprint < signals[j].Fingerprint
		}
		return signals[i].Type < signals[j].Type
	})
}
