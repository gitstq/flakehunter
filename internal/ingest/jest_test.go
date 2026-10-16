package ingest

import (
	"os"
	"testing"

	"github.com/gitstq/flakehunter/internal/model"
)

func TestParseJestReport(t *testing.T) {
	f, err := os.Open("../../testdata/jest/report.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	results, err := ParseJest(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("want 3 results, got %d", len(results))
	}
	want := []model.Outcome{
		model.OutcomePassed,
		model.OutcomeFailed,
		model.OutcomeSkipped, // pending
	}
	for i, w := range want {
		if results[i].Outcome != w {
			t.Errorf("result %d: want %s, got %s", i, w, results[i].Outcome)
		}
	}
	if results[0].Fingerprint != "/app/src/math.test.js#math adds" {
		t.Errorf("unexpected fingerprint: %s", results[0].Fingerprint)
	}
	if results[0].DurationMS != 2 {
		t.Errorf("duration: want 2ms, got %v", results[0].DurationMS)
	}
}
