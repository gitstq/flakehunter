package report

import (
	"fmt"
	"io"
	"os"

	"github.com/gitstq/flakehunter/internal/engine"
)

// ANSI colors, disabled when NO_COLOR is set.
const (
	colReset  = "\033[0m"
	colRed    = "\033[31m"
	colGreen  = "\033[32m"
	colYellow = "\033[33m"
	colDim    = "\033[2m"
)

// useColor reports whether ANSI colors should be emitted.
func useColor() bool {
	return os.Getenv("NO_COLOR") == ""
}

// paint wraps s in an ANSI color when colors are enabled.
func paint(code, s string) string {
	if !useColor() {
		return s
	}
	return code + s + colReset
}

// statusColor maps a status to its ANSI color.
func statusColor(status string) string {
	switch status {
	case engine.StatusFlaky:
		return colYellow
	case engine.StatusStableFail:
		return colRed
	case engine.StatusStablePass:
		return colGreen
	default:
		return colDim
	}
}

// Terminal writes a human-friendly flake report with a colored table.
func Terminal(w io.Writer, stats []engine.Stats, signals []engine.Signal, runsAnalyzed int) {
	rows := make([]engine.Stats, len(stats))
	copy(rows, stats)
	sortStats(rows)
	counts := CountStats(rows)
	fmt.Fprintf(w, "🎲 FlakeHunter — %d runs, %d tests tracked (✅ %d · 🎲 %d · 💥 %d · ❓ %d)\n\n",
		runsAnalyzed, counts.Total, counts.Stable, counts.Flaky, counts.Broken, counts.Insufficient)
	fmt.Fprintf(w, "%-3s %-6s %-8s %-5s %-7s %-9s %s\n",
		"", "RUNS", "P/F/E", "CONF", "FLAKE", "MEDIAN", "TEST")
	for _, s := range rows {
		flake := "-"
		conf := "-"
		if s.Status == engine.StatusFlaky {
			flake = fmt.Sprintf("%.1f%%", s.FlakeScore*100)
			conf = confLabel(s.Confidence)
		}
		name := s.Fingerprint
		if s.Quarantined {
			name += " (quarantined)"
		}
		fmt.Fprintf(w, "%s %-6d %d/%d/%d %-5s %-7s %-9s %s\n",
			paint(statusColor(s.Status), statusEmoji(s.Status)),
			s.Total, s.Passed, s.Failed, s.Errored,
			conf, flake, fmt.Sprintf("%.0fms", s.MedianMS), name)
	}
	if len(signals) > 0 {
		fmt.Fprintf(w, "\n⚠️  Order-dependence signals:\n")
		for _, sig := range signals {
			fmt.Fprintf(w, "   %s — %s\n", sig.Fingerprint, sig.Detail)
		}
	}
	fmt.Fprintln(w)
	if counts.Flaky > 0 {
		fmt.Fprintf(w, "🎲 %d flaky tests, %.1f%% flake rate. Quarantine with: flakehunter quarantine add <fingerprint>\n",
			counts.Flaky, FlakeRateOf(counts))
	}
	if counts.Insufficient > 0 {
		fmt.Fprintf(w, "❓ %d tests need more runs before classification (see min_runs in config).\n", counts.Insufficient)
	}
	if counts.Flaky == 0 && counts.Broken == 0 && counts.Insufficient == 0 {
		fmt.Fprintln(w, "✅ All tracked tests look stable.")
	}
}
