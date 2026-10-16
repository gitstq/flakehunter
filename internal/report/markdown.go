package report

import (
	"fmt"
	"io"
)

// Markdown writes a GitHub-flavored flake report suitable for
// $GITHUB_STEP_SUMMARY or a PR comment.
func Markdown(w io.Writer, s Summary) {
	fmt.Fprintln(w, "## 🎲 FlakeHunter Report")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "- Runs analyzed: **%d**\n", s.RunsAnalyzed)
	fmt.Fprintf(w, "- Tests tracked: **%d** (✅ stable %d · 🎲 flaky %d · 💥 broken %d · ❓ insufficient %d)\n",
		s.Counts.Total, s.Counts.Stable, s.Counts.Flaky, s.Counts.Broken, s.Counts.Insufficient)
	fmt.Fprintf(w, "- Flake rate: **%.1f%%**\n", s.FlakeRate)
	if len(s.Flaky) > 0 {
		fmt.Fprintln(w, "\n### 🎲 Flaky tests")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "| Test | Runs | Pass/Fail | Flake score | Confidence | Median | Max |")
		fmt.Fprintln(w, "| --- | --- | --- | --- | --- | --- | --- |")
		for _, f := range s.Flaky {
			fmt.Fprintf(w, "| `%s` | %d | %d/%d | %.1f%% | %s | %.0fms | %.0fms |\n",
				f.Fingerprint, f.Total, f.Passed, f.Failed, f.FlakeScore*100,
				confLabel(f.Confidence), f.MedianMS, f.MaxMS)
		}
	}
	if len(s.Broken) > 0 {
		fmt.Fprintln(w, "\n### 💥 Consistently failing tests")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "| Test | Runs | Failed | Errored | Last seen |")
		fmt.Fprintln(w, "| --- | --- | --- | --- | --- |")
		for _, b := range s.Broken {
			fmt.Fprintf(w, "| `%s` | %d | %d | %d | %s |\n",
				b.Fingerprint, b.Total, b.Failed, b.Errored, b.LastSeen.Format("2006-01-02 15:04"))
		}
	}
	if len(s.Signals) > 0 {
		fmt.Fprintln(w, "\n### ⚠️ Order-dependence signals")
		fmt.Fprintln(w)
		for _, sig := range s.Signals {
			fmt.Fprintf(w, "- `%s` — %s\n", sig.Fingerprint, sig.Detail)
		}
	}
	fmt.Fprintln(w)
}
