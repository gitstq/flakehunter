package ingest

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/gitstq/flakehunter/internal/model"
)

// junitDocument accepts both <testsuites> roots and bare <testsuite> roots.
type junitDocument struct {
	Suites []junitSuite `xml:"testsuite"`
}

type junitSuite struct {
	Name      string          `xml:"name,attr"`
	File      string          `xml:"file,attr"`
	TestCases []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	Name         string        `xml:"name,attr"`
	ClassName    string        `xml:"classname,attr"`
	File         string        `xml:"file,attr"`
	Time         float64       `xml:"time,attr"`
	Failure      *junitMessage `xml:"failure"`
	Error        *junitMessage `xml:"error"`
	Skipped      *junitSkip    `xml:"skipped"`
	RerunFailure *junitMessage `xml:"rerunFailure"`
}

type junitMessage struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

type junitSkip struct {
	Message string `xml:"message,attr"`
}

// ParseJUnit reads a JUnit XML report with either a <testsuites> or a
// <testsuite> root element. It understands the rerun markup emitted by
// pytest-rerunfailures (<rerunFailure>) and marks those tests as rerun-pass.
func ParseJUnit(r io.Reader) ([]model.TestResult, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var doc junitDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if len(doc.Suites) == 0 {
		var single junitSuite
		if err := xml.Unmarshal(data, &single); err != nil {
			return nil, err
		}
		if single.Name == "" && len(single.TestCases) == 0 {
			return nil, fmt.Errorf("no <testsuite> elements found")
		}
		doc.Suites = []junitSuite{single}
	}
	var results []model.TestResult
	for _, suite := range doc.Suites {
		for _, tc := range suite.TestCases {
			suiteName := suite.Name
			if tc.ClassName != "" {
				suiteName = tc.ClassName
			}
			outcome := model.OutcomePassed
			switch {
			case tc.Skipped != nil:
				outcome = model.OutcomeSkipped
			case tc.Error != nil:
				outcome = model.OutcomeError
			case tc.Failure != nil && tc.RerunFailure != nil:
				// Failed, then passed on rerun: a definite flake.
				outcome = model.OutcomePassed
			case tc.Failure != nil:
				outcome = model.OutcomeFailed
			}
			file := tc.File
			if file == "" {
				file = suite.File
			}
			results = append(results, model.TestResult{
				Fingerprint: model.FingerprintOf(suiteName, tc.Name),
				Suite:       suiteName,
				Name:        tc.Name,
				File:        file,
				Outcome:     outcome,
				DurationMS:  tc.Time * 1000,
				Retried:     tc.RerunFailure != nil,
				RerunPass:   tc.Failure != nil && tc.RerunFailure != nil,
			})
		}
	}
	return results, nil
}
