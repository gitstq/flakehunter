package engine

import (
	"math"
	"sort"
	"time"

	"github.com/gitstq/flakehunter/internal/model"
)

// Status values for a test fingerprint.
const (
	StatusStablePass   = "stable_pass"
	StatusStableFail   = "stable_fail"
	StatusFlaky        = "flaky"
	StatusInsufficient = "insufficient"
)

// Stats is the aggregated history of one test fingerprint.
type Stats struct {
	Fingerprint string    `json:"fingerprint"`
	Suite       string    `json:"suite,omitempty"`
	Name        string    `json:"name,omitempty"`
	File        string    `json:"file,omitempty"`
	Status      string    `json:"status"`
	Total       int       `json:"total"`
	Passed      int       `json:"passed"`
	Failed      int       `json:"failed"`
	Skipped     int       `json:"skipped"`
	Errored     int       `json:"errored"`
	FlakeScore  float64   `json:"flake_score"`
	Confidence  float64   `json:"confidence"`
	LastOutcome string    `json:"last_outcome,omitempty"`
	MedianMS    float64   `json:"median_ms,omitempty"`
	MaxMS       float64   `json:"max_ms,omitempty"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	RerunPasses int       `json:"rerun_passes,omitempty"`
	Quarantined bool      `json:"quarantined,omitempty"`
}

// Options controls analysis behavior.
type Options struct {
	MinRuns int
}

// Analyze aggregates run history into per-test statistics.
func Analyze(runs []model.Run, opts Options) []Stats {
	if opts.MinRuns <= 0 {
		opts.MinRuns = 3
	}
	type acc struct {
		stats     Stats
		durations []float64
	}
	accs := map[string]*acc{}
	var order []string
	for _, run := range runs {
		for _, res := range run.Results {
			a, ok := accs[res.Fingerprint]
			if !ok {
				a = &acc{stats: Stats{
					Fingerprint: res.Fingerprint,
					Suite:       res.Suite,
					Name:        res.Name,
					File:        res.File,
					FirstSeen:   run.Timestamp,
					LastSeen:    run.Timestamp,
				}}
				accs[res.Fingerprint] = a
				order = append(order, res.Fingerprint)
			}
			s := &a.stats
			s.Total++
			s.LastOutcome = string(res.Outcome)
			if res.DurationMS > 0 {
				a.durations = append(a.durations, res.DurationMS)
			}
			if res.RerunPass {
				s.RerunPasses++
			}
			if run.Timestamp.After(s.LastSeen) {
				s.LastSeen = run.Timestamp
			}
			if run.Timestamp.Before(s.FirstSeen) {
				s.FirstSeen = run.Timestamp
			}
			switch res.Outcome {
			case model.OutcomePassed:
				s.Passed++
			case model.OutcomeFailed:
				s.Failed++
			case model.OutcomeSkipped:
				s.Skipped++
			case model.OutcomeError:
				s.Errored++
			}
		}
	}
	out := make([]Stats, 0, len(order))
	for _, fp := range order {
		a := accs[fp]
		if len(a.durations) > 0 {
			sorted := append([]float64{}, a.durations...)
			sort.Float64s(sorted)
			a.stats.MedianMS = median(sorted)
			a.stats.MaxMS = sorted[len(sorted)-1]
		}
		out = append(out, classify(a.stats, opts.MinRuns))
	}
	return out
}

// classify assigns a status and flake score to one fingerprint's statistics.
func classify(s Stats, minRuns int) Stats {
	bad := s.Failed + s.Errored
	switch {
	case s.Total < minRuns:
		s.Status = StatusInsufficient
	case bad == s.Total:
		s.Status = StatusStableFail
	case bad == 0:
		s.Status = StatusStablePass
	case s.Passed == 0:
		// Mixed failures and skips, but never green: consistently broken.
		s.Status = StatusStableFail
	default:
		s.Status = StatusFlaky
		minor := math.Min(float64(s.Passed), float64(bad))
		s.FlakeScore = round4(minor / float64(s.Total))
		s.Confidence = round4(wilsonLower(minor, float64(s.Total)))
	}
	return s
}

// wilsonLower returns the lower bound of the Wilson score interval for a
// proportion, using a 95% confidence level. It answers: "how confident are
// we that this test really flips, and not just got unlucky?"
func wilsonLower(successes, n float64) float64 {
	if n <= 0 {
		return 0
	}
	const z = 1.96
	p := successes / n
	denominator := 1 + z*z/n
	center := p + z*z/(2*n)
	margin := z * math.Sqrt(p*(1-p)/n+z*z/(4*n*n))
	lower := (center - margin) / denominator
	if lower < 0 {
		return 0
	}
	return lower
}

// median returns the median of an already-sorted slice.
func median(sorted []float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// round4 rounds to four decimal places.
func round4(x float64) float64 {
	return math.Round(x*10000) / 10000
}
