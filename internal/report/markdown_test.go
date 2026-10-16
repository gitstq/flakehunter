package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/model"
)

func TestBuildSummaryBuckets(t *testing.T) {
	stats := []engine.Stats{
		{Fingerprint: "s#a", Status: engine.StatusStablePass, Total: 10},
		{Fingerprint: "s#b", Status: engine.StatusFlaky, Total: 10, FlakeScore: 0.4, Confidence: 0.2},
		{Fingerprint: "s#c", Status: engine.StatusStableFail, Total: 10},
		{Fingerprint: "s#d", Status: engine.StatusInsufficient, Total: 1},
	}
	summary := BuildSummary(10, stats, nil)
	if summary.Counts.Total != 4 || summary.Counts.Flaky != 1 || summary.Counts.Stable != 1 ||
		summary.Counts.Broken != 1 || summary.Counts.Insufficient != 1 {
		t.Fatalf("unexpected counts: %+v", summary.Counts)
	}
	if summary.FlakeRate < 33 || summary.FlakeRate > 34 {
		t.Errorf("want flake rate ~33.3%%, got %v", summary.FlakeRate)
	}
	if summary.RunsAnalyzed != 10 {
		t.Errorf("want RunsAnalyzed=10, got %d", summary.RunsAnalyzed)
	}
}

func TestMarkdownContainsTables(t *testing.T) {
	stats := []engine.Stats{
		{
			Fingerprint: "tests.test_net#test_fetch",
			Status:      engine.StatusFlaky,
			Total:       10,
			Passed:      6,
			Failed:      4,
			FlakeScore:  0.4,
			Confidence:  0.2,
			MedianMS:    120,
			MaxMS:       900,
		},
	}
	summary := BuildSummary(10, stats, []engine.Signal{
		{Type: engine.SignalOrderDep, Fingerprint: "tests.test_net#test_fetch", Detail: "fails after s#x"},
	})
	var buf bytes.Buffer
	Markdown(&buf, summary)
	out := buf.String()
	for _, want := range []string{
		"## 🎲 FlakeHunter Report",
		"### 🎲 Flaky tests",
		"| `tests.test_net#test_fetch` |",
		"### ⚠️ Order-dependence signals",
		"fails after s#x",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("markdown output missing %q, got:\n%s", want, out)
		}
	}
}

func TestScanSummaryCounters(t *testing.T) {
	run := model.Run{
		RunID:     "r1",
		Timestamp: time.Now().UTC(),
		Source:    "testdata/junit/pytest-mixed.xml (junit)",
		Results: []model.TestResult{
			{Fingerprint: "s#a", Outcome: model.OutcomePassed},
			{Fingerprint: "s#b", Outcome: model.OutcomeFailed},
			{Fingerprint: "s#c", Outcome: model.OutcomeError},
			{Fingerprint: "s#d", Outcome: model.OutcomeSkipped},
			{Fingerprint: "s#e", Outcome: model.OutcomePassed, RerunPass: true},
		},
	}
	summary := BuildScanSummary(run, nil)
	if summary.Total != 5 || summary.Passed != 2 || summary.Failed != 1 ||
		summary.Errored != 1 || summary.Skipped != 1 || summary.RerunPasses != 1 {
		t.Fatalf("unexpected scan summary counters: %+v", summary)
	}
	if len(summary.Failures) != 2 {
		t.Errorf("want 2 failing fingerprints, got %v", summary.Failures)
	}
	var buf bytes.Buffer
	ScanMarkdown(&buf, summary)
	out := buf.String()
	for _, want := range []string{
		"## 🔍 FlakeHunter Scan",
		"Rerun-pass (definite flakes): **1** 🎲",
		"### ❌ Failing tests",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("scan markdown missing %q, got:\n%s", want, out)
		}
	}
}
