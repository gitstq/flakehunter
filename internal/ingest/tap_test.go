package ingest

import (
	"os"
	"testing"

	"github.com/gitstq/flakehunter/internal/model"
)

func TestParseTAPFile(t *testing.T) {
	f, err := os.Open("../../testdata/tap/basic.tap")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	results, err := ParseTAP(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 4 {
		t.Fatalf("want 4 results, got %d", len(results))
	}
	want := []model.Outcome{
		model.OutcomePassed,
		model.OutcomeFailed,
		model.OutcomeSkipped, // # SKIP directive
		model.OutcomeSkipped, // # TODO directive (expected failure)
	}
	for i, w := range want {
		if results[i].Outcome != w {
			t.Errorf("result %d (%s): want %s, got %s", i, results[i].Name, w, results[i].Outcome)
		}
	}
	if results[0].Name != "adds numbers correctly" {
		t.Errorf("unexpected first test name: %q", results[0].Name)
	}
}

func TestParseTAPWithoutDash(t *testing.T) {
	stream := "TAP version 13\n1..2\nok 1 simple name\nnot ok 2 other name\n"
	results, err := parseString(stream, FormatTAP)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 results, got %d", len(results))
	}
	if results[0].Name != "simple name" || results[1].Name != "other name" {
		t.Errorf("names not stripped of test numbers: %q, %q", results[0].Name, results[1].Name)
	}
}
