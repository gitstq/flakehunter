// Package gate turns flake analysis into a CI pass/fail decision.
package gate

import (
	"fmt"
	"sort"

	"github.com/gitstq/flakehunter/internal/engine"
)

// Options configures the quality gate.
type Options struct {
	// MaxFlakeRate is the maximum percentage of analyzed tests allowed to be flaky.
	MaxFlakeRate float64
	// MaxFlaky is the maximum absolute number of flaky tests allowed.
	MaxFlaky int
	// FailOnBroken fails the gate when consistently failing tests exist.
	FailOnBroken bool
	// Quarantined lists fingerprints excluded from gating.
	Quarantined map[string]bool
}

// Result is the gate decision plus diagnostics.
type Result struct {
	Pass        bool     `json:"pass"`
	FlakeRate   float64  `json:"flake_rate"`
	FlakyCount  int      `json:"flaky_count"`
	BrokenCount int      `json:"broken_count"`
	Analyzed    int      `json:"analyzed"`
	Violations  []string `json:"violations,omitempty"`
	Annotations []string `json:"annotations,omitempty"`
}

// Evaluate applies gate options to analyzed statistics.
func Evaluate(stats []engine.Stats, opts Options) Result {
	res := Result{Pass: true}
	for _, s := range stats {
		if opts.Quarantined[s.Fingerprint] {
			continue
		}
		switch s.Status {
		case engine.StatusInsufficient:
			continue
		case engine.StatusFlaky:
			res.FlakyCount++
			res.Analyzed++
			res.Annotations = append(res.Annotations, fmt.Sprintf(
				"::warning title=FlakeHunter::flaky test %s (flake %.1f%%, %d runs)",
				s.Fingerprint, s.FlakeScore*100, s.Total))
		case engine.StatusStableFail:
			res.BrokenCount++
			res.Analyzed++
			res.Annotations = append(res.Annotations, fmt.Sprintf(
				"::error title=FlakeHunter::consistently failing test %s (%d runs)",
				s.Fingerprint, s.Total))
		default:
			res.Analyzed++
		}
	}
	if res.Analyzed > 0 {
		res.FlakeRate = float64(res.FlakyCount) / float64(res.Analyzed) * 100
	}
	if res.FlakeRate > opts.MaxFlakeRate {
		res.Pass = false
		res.Violations = append(res.Violations, fmt.Sprintf(
			"flake rate %.1f%% exceeds limit %.1f%%", res.FlakeRate, opts.MaxFlakeRate))
	}
	if res.FlakyCount > opts.MaxFlaky {
		res.Pass = false
		res.Violations = append(res.Violations, fmt.Sprintf(
			"%d flaky tests exceed limit %d", res.FlakyCount, opts.MaxFlaky))
	}
	if opts.FailOnBroken && res.BrokenCount > 0 {
		res.Pass = false
		res.Violations = append(res.Violations, fmt.Sprintf(
			"%d consistently failing tests", res.BrokenCount))
	}
	sort.Strings(res.Annotations)
	return res
}
