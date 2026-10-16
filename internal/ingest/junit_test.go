package ingest

import (
	"os"
	"testing"

	"github.com/gitstq/flakehunter/internal/model"
)

func TestParseJUnitMixed(t *testing.T) {
	f, err := os.Open("../../testdata/junit/pytest-mixed.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	results, err := ParseJUnit(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 5 {
		t.Fatalf("want 5 results, got %d", len(results))
	}
	byFP := map[string]model.TestResult{}
	for _, r := range results {
		byFP[r.Fingerprint] = r
	}
	cases := []struct {
		fp      string
		outcome model.Outcome
		rerun   bool
	}{
		{"tests.test_math#test_add", model.OutcomePassed, false},
		{"tests.test_math#test_sub", model.OutcomeFailed, false},
		{"tests.test_net#test_fetch", model.OutcomePassed, true},
		{"tests.test_net#test_skipme", model.OutcomeSkipped, false},
		{"tests.test_db#test_query", model.OutcomeError, false},
	}
	for _, c := range cases {
		r, ok := byFP[c.fp]
		if !ok {
			t.Errorf("missing fingerprint %s", c.fp)
			continue
		}
		if r.Outcome != c.outcome {
			t.Errorf("%s: want outcome %s, got %s", c.fp, c.outcome, r.Outcome)
		}
		if r.Rerun != c.rerun {
			t.Errorf("%s: want rerun=%v, got %v", c.fp, c.rerun, r.Rerun)
		}
	}
	// Duration is converted from seconds to milliseconds.
	if got := byFP["tests.test_math#test_add"].DurationMS; got != 12 {
		t.Errorf("test_add duration: want 12ms, got %v", got)
	}
	// The rerun-pass test must be flagged as a definite flake.
	if !byFP["tests.test_net#test_fetch"].RerunPass {
		t.Error("test_fetch: want RerunPass=true")
	}
}

func TestParseJUnitSingleSuiteRoot(t *testing.T) {
	xml := `<?xml version="1.0"?>
<testsuite name="maven" tests="2" failures="1">
  <testcase classname="com.example.CalcTest" name="testAdd" time="0.01"/>
  <testcase classname="com.example.CalcTest" name="testDiv" time="0.02">
    <failure message="boom">stack</failure>
  </testcase>
</testsuite>`
	results, err := parseString(xml, FormatJUnit)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 results, got %d", len(results))
	}
	if results[0].Outcome != model.OutcomePassed || results[1].Outcome != model.OutcomeFailed {
		t.Fatalf("unexpected outcomes: %v", results)
	}
}

func TestParseJUnitParamNormalization(t *testing.T) {
	xml := `<?xml version="1.0"?>
<testsuite name="pytest" tests="2">
  <testcase classname="tests.test_p" name="test_case[alpha-1]" time="0.01"/>
  <testcase classname="tests.test_p" name="test_case[beta-2]" time="0.01"/>
</testsuite>`
	results, err := parseString(xml, FormatJUnit)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Fingerprint != results[1].Fingerprint {
		t.Errorf("parametrized cases should share a fingerprint: %s vs %s",
			results[0].Fingerprint, results[1].Fingerprint)
	}
}
