package report

import (
	"time"

	"github.com/gitstq/flakehunter/internal/engine"
)

// Summary is the machine-readable flake analysis of recorded history.
type Summary struct {
	GeneratedAt  time.Time       `json:"generated_at"`
	RunsAnalyzed int             `json:"runs_analyzed"`
	Counts       Counts          `json:"counts"`
	FlakeRate    float64         `json:"flake_rate"`
	Flaky        []engine.Stats  `json:"flaky,omitempty"`
	Broken       []engine.Stats  `json:"broken,omitempty"`
	Stable       []engine.Stats  `json:"stable,omitempty"`
	Insufficient []engine.Stats  `json:"insufficient,omitempty"`
	Signals      []engine.Signal `json:"signals,omitempty"`
}

// BuildSummary assembles a summary from analyzed statistics.
func BuildSummary(runsAnalyzed int, stats []engine.Stats, signals []engine.Signal) Summary {
	s := Summary{
		GeneratedAt:  time.Now().UTC(),
		RunsAnalyzed: runsAnalyzed,
		Signals:      signals,
	}
	for _, st := range stats {
		switch st.Status {
		case engine.StatusFlaky:
			s.Flaky = append(s.Flaky, st)
		case engine.StatusStableFail:
			s.Broken = append(s.Broken, st)
		case engine.StatusStablePass:
			s.Stable = append(s.Stable, st)
		default:
			s.Insufficient = append(s.Insufficient, st)
		}
	}
	s.Counts = CountStats(stats)
	s.FlakeRate = FlakeRateOf(s.Counts)
	sortStats(s.Flaky)
	sortStats(s.Broken)
	sortStats(s.Stable)
	sortStats(s.Insufficient)
	return s
}
