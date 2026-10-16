package report

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/model"
)

// ScanSummary is the one-off analysis of a single run.
type ScanSummary struct {
	RunID       string          `json:"run_id"`
	Source      string          `json:"source,omitempty"`
	Commit      string          `json:"commit,omitempty"`
	GeneratedAt time.Time       `json:"generated_at"`
	Total       int             `json:"total"`
	Passed      int             `json:"passed"`
	Failed      int             `json:"failed"`
	Skipped     int             `json:"skipped"`
	Errored     int             `json:"errored"`
	RerunPasses int             `json:"rerun_passes"`
	Failures    []string        `json:"failures,omitempty"`
	Signals     []engine.Signal `json:"signals,omitempty"`
}

// BuildScanSummary aggregates a single run and its detected signals.
func BuildScanSummary(run model.Run, signals []engine.Signal) ScanSummary {
	s := ScanSummary{
		RunID:       run.RunID,
		Source:      run.Source,
		Commit:      run.Commit,
		GeneratedAt: time.Now().UTC(),
		Signals:     signals,
	}
	seen := map[string]bool{}
	for _, r := range run.Results {
		s.Total++
		switch r.Outcome {
		case model.OutcomePassed:
			s.Passed++
		case model.OutcomeFailed:
			s.Failed++
		case model.OutcomeSkipped:
			s.Skipped++
		case model.OutcomeError:
			s.Errored++
		}
		if r.RerunPass {
			s.RerunPasses++
		}
		if (r.Outcome == model.OutcomeFailed || r.Outcome == model.OutcomeError) && !seen[r.Fingerprint] {
			seen[r.Fingerprint] = true
			s.Failures = append(s.Failures, r.Fingerprint)
		}
	}
	sort.Strings(s.Failures)
	return s
}

// ScanTerminal prints the one-off run report.
func ScanTerminal(w io.Writer, s ScanSummary) {
	fmt.Fprintf(w, "🔍 FlakeHunter scan — run %s\n", s.RunID)
	if s.Source != "" {
		fmt.Fprintf(w, "   source: %s\n", s.Source)
	}
	fmt.Fprintf(w, "   total %d · ✅ %d · ❌ %d · ⏭️ %d · 💥 %d",
		s.Total, s.Passed, s.Failed, s.Skipped, s.Errored)
	if s.RerunPasses > 0 {
		fmt.Fprintf(w, " · 🎲 %d rerun-pass", s.RerunPasses)
	}
	fmt.Fprintln(w)
	if len(s.Failures) > 0 {
		fmt.Fprintln(w, "\n❌ Failing tests:")
		for _, f := range s.Failures {
			fmt.Fprintf(w, "   %s\n", f)
		}
	}
	if len(s.Signals) > 0 {
		fmt.Fprintln(w, "\n⚠️  Flake signals (single run):")
		for _, sig := range s.Signals {
			fmt.Fprintf(w, "   [%s] %s — %s\n", sig.Type, sig.Fingerprint, sig.Detail)
		}
	} else if s.Failed+s.Errored == 0 {
		fmt.Fprintln(w, "\n✅ No flake signals detected in this run.")
	}
}

// ScanMarkdown writes the one-off run report as GitHub-flavored markdown.
func ScanMarkdown(w io.Writer, s ScanSummary) {
	fmt.Fprintln(w, "## 🔍 FlakeHunter Scan")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "- Total: **%d** · Passed: **%d** · Failed: **%d** · Skipped: **%d** · Errors: **%d**\n",
		s.Total, s.Passed, s.Failed, s.Skipped, s.Errored)
	if s.RerunPasses > 0 {
		fmt.Fprintf(w, "- Rerun-pass (definite flakes): **%d** 🎲\n", s.RerunPasses)
	}
	if len(s.Signals) > 0 {
		fmt.Fprintln(w, "\n### ⚠️ Flake signals")
		fmt.Fprintln(w)
		for _, sig := range s.Signals {
			fmt.Fprintf(w, "- `%s` — **%s** (%s)\n", sig.Fingerprint, sig.Type, sig.Detail)
		}
	}
	if len(s.Failures) > 0 {
		fmt.Fprintln(w, "\n### ❌ Failing tests")
		fmt.Fprintln(w)
		for _, f := range s.Failures {
			fmt.Fprintf(w, "- `%s`\n", f)
		}
	}
	fmt.Fprintln(w)
}
