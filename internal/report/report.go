// Package report renders flake analysis for humans (terminal, markdown)
// and machines (JSON).
package report

import (
	"sort"

	"github.com/gitstq/flakehunter/internal/engine"
)

// statusEmoji maps a test status to its icon.
func statusEmoji(status string) string {
	switch status {
	case engine.StatusStablePass:
		return "✅"
	case engine.StatusFlaky:
		return "🎲"
	case engine.StatusStableFail:
		return "💥"
	default:
		return "❓"
	}
}

// confLabel converts a Wilson confidence bound into a coarse label.
func confLabel(confidence float64) string {
	switch {
	case confidence >= 0.15:
		return "high"
	case confidence >= 0.07:
		return "med"
	default:
		return "low"
	}
}

// statusRank orders statuses so problems surface first.
func statusRank(status string) int {
	switch status {
	case engine.StatusFlaky:
		return 0
	case engine.StatusStableFail:
		return 1
	case engine.StatusInsufficient:
		return 2
	default:
		return 3
	}
}

// sortStats orders stats for terminal display: problems first.
func sortStats(stats []engine.Stats) {
	sort.Slice(stats, func(i, j int) bool {
		ri, rj := statusRank(stats[i].Status), statusRank(stats[j].Status)
		if ri != rj {
			return ri < rj
		}
		if stats[i].FlakeScore != stats[j].FlakeScore {
			return stats[i].FlakeScore > stats[j].FlakeScore
		}
		return stats[i].Fingerprint < stats[j].Fingerprint
	})
}

// Counts summarizes statistics by status.
type Counts struct {
	Total        int
	Flaky        int
	Broken       int
	Stable       int
	Insufficient int
}

// CountStats tallies statuses.
func CountStats(stats []engine.Stats) Counts {
	var c Counts
	for _, s := range stats {
		c.Total++
		switch s.Status {
		case engine.StatusFlaky:
			c.Flaky++
		case engine.StatusStableFail:
			c.Broken++
		case engine.StatusStablePass:
			c.Stable++
		default:
			c.Insufficient++
		}
	}
	return c
}

// FlakeRateOf returns the percentage of analyzed tests that are flaky.
func FlakeRateOf(c Counts) float64 {
	analyzed := c.Flaky + c.Stable + c.Broken
	if analyzed == 0 {
		return 0
	}
	return float64(c.Flaky) / float64(analyzed) * 100
}
