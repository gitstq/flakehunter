package gate

import (
	"strings"
	"testing"

	"github.com/gitstq/flakehunter/internal/engine"
)

func statsOf(fp, status string, total int, score float64) engine.Stats {
	return engine.Stats{Fingerprint: fp, Status: status, Total: total, FlakeScore: score}
}

func TestEvaluatePassesUnderThreshold(t *testing.T) {
	stats := []engine.Stats{
		statsOf("s#a", engine.StatusStablePass, 10, 0),
		statsOf("s#b", engine.StatusStablePass, 10, 0),
		statsOf("s#c", engine.StatusFlaky, 10, 0.5),
	}
	res := Evaluate(stats, Options{MaxFlakeRate: 50, MaxFlaky: 1})
	if !res.Pass {
		t.Errorf("want pass, got violations: %v", res.Violations)
	}
	if res.FlakeRate < 33 || res.FlakeRate > 34 {
		t.Errorf("want flake rate ~33.3%%, got %v", res.FlakeRate)
	}
	if len(res.Annotations) != 1 || !strings.Contains(res.Annotations[0], "::warning") {
		t.Errorf("want one warning annotation, got %v", res.Annotations)
	}
}

func TestEvaluateFailsOverThreshold(t *testing.T) {
	stats := []engine.Stats{
		statsOf("s#a", engine.StatusStablePass, 10, 0),
		statsOf("s#b", engine.StatusFlaky, 10, 0.5),
		statsOf("s#c", engine.StatusFlaky, 10, 0.3),
	}
	res := Evaluate(stats, Options{MaxFlakeRate: 10, MaxFlaky: 0})
	if res.Pass {
		t.Error("want gate failure for 66.7% flake rate")
	}
	if len(res.Violations) < 2 {
		t.Errorf("want violations for both rate and count, got %v", res.Violations)
	}
}

func TestEvaluateQuarantineExempts(t *testing.T) {
	stats := []engine.Stats{
		statsOf("s#a", engine.StatusStablePass, 10, 0),
		statsOf("s#b", engine.StatusFlaky, 10, 0.5),
	}
	res := Evaluate(stats, Options{
		MaxFlakeRate: 10,
		MaxFlaky:     0,
		Quarantined:  map[string]bool{"s#b": true},
	})
	if !res.Pass {
		t.Errorf("quarantined flaky test should not fail the gate: %v", res.Violations)
	}
	if res.FlakyCount != 0 {
		t.Errorf("want 0 flaky counted, got %d", res.FlakyCount)
	}
}

func TestEvaluateFailOnBroken(t *testing.T) {
	stats := []engine.Stats{
		statsOf("s#a", engine.StatusStablePass, 10, 0),
		statsOf("s#b", engine.StatusStableFail, 10, 0),
	}
	if res := Evaluate(stats, Options{MaxFlakeRate: 10, MaxFlaky: 5}); !res.Pass {
		t.Error("broken tests should not fail unless FailOnBroken is set")
	}
	res := Evaluate(stats, Options{MaxFlakeRate: 10, MaxFlaky: 5, FailOnBroken: true})
	if res.Pass {
		t.Error("want failure with FailOnBroken")
	}
	if len(res.Annotations) != 1 || !strings.Contains(res.Annotations[0], "::error") {
		t.Errorf("want one error annotation, got %v", res.Annotations)
	}
}

func TestEvaluateIgnoresInsufficient(t *testing.T) {
	stats := []engine.Stats{
		statsOf("s#a", engine.StatusInsufficient, 1, 0),
		statsOf("s#b", engine.StatusStablePass, 10, 0),
	}
	res := Evaluate(stats, Options{MaxFlakeRate: 0, MaxFlaky: 0})
	if !res.Pass {
		t.Errorf("insufficient data should not affect the gate: %v", res.Violations)
	}
	if res.Analyzed != 1 {
		t.Errorf("want 1 analyzed test, got %d", res.Analyzed)
	}
}
