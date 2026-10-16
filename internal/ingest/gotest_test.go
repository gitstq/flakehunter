package ingest

import (
	"os"
	"testing"

	"github.com/gitstq/flakehunter/internal/model"
)

func TestParseGoTestEvents(t *testing.T) {
	f, err := os.Open("../../testdata/gotest/events.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	results, err := ParseGoTest(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("want 3 results, got %d", len(results))
	}
	want := []struct {
		fp      string
		outcome model.Outcome
	}{
		{"github.com/gitstq/flakehunter/internal/engine#TestAnalyze", model.OutcomePassed},
		{"github.com/gitstq/flakehunter/internal/ingest#TestParseJUnit", model.OutcomeFailed},
		{"github.com/gitstq/flakehunter/internal/model#TestFingerprint", model.OutcomeSkipped},
	}
	for i, w := range want {
		if results[i].Fingerprint != w.fp {
			t.Errorf("result %d: want fingerprint %s, got %s", i, w.fp, results[i].Fingerprint)
		}
		if results[i].Outcome != w.outcome {
			t.Errorf("result %d: want outcome %s, got %s", i, w.outcome, results[i].Outcome)
		}
	}
	// Package-level events (empty Test) must be ignored, and duration
	// converted from seconds to milliseconds.
	if results[0].DurationMS != 10 {
		t.Errorf("TestAnalyze duration: want 10ms, got %v", results[0].DurationMS)
	}
}

func TestParseGoTestIgnoresNoise(t *testing.T) {
	stream := `{"Action":"run","Package":"p","Test":"T"}
not json at all
{"Action":"output","Package":"p","Test":"T","Output":"hello\n"}
{"Action":"pass","Package":"p","Test":"T","Elapsed":0.5}
`
	results, err := parseString(stream, FormatGoTest)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("want 1 result, got %d", len(results))
	}
	if results[0].Outcome != model.OutcomePassed {
		t.Errorf("want passed, got %s", results[0].Outcome)
	}
}
