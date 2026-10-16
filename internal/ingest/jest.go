package ingest

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/gitstq/flakehunter/internal/model"
)

// jestReport is the shape produced by `jest --json` and Vitest's JSON reporter.
type jestReport struct {
	NumTotalTests int            `json:"numTotalTests"`
	TestResults   []jestTestFile `json:"testResults"`
}

type jestTestFile struct {
	Name             string          `json:"name"`
	AssertionResults []jestAssertion `json:"assertionResults"`
}

type jestAssertion struct {
	Title          string   `json:"title"`
	FullName       string   `json:"fullName"`
	Status         string   `json:"status"`
	Duration       float64  `json:"duration"`
	AncestorTitles []string `json:"ancestorTitles"`
}

// ParseJest reads JSON reports produced by Jest (--json) or Vitest (--reporter=json).
func ParseJest(r io.Reader) ([]model.TestResult, error) {
	var report jestReport
	if err := json.NewDecoder(r).Decode(&report); err != nil {
		return nil, err
	}
	var results []model.TestResult
	order := 0
	for _, file := range report.TestResults {
		for _, a := range file.AssertionResults {
			name := a.FullName
			if name == "" {
				parts := append(append([]string{}, a.AncestorTitles...), a.Title)
				name = strings.Join(parts, " > ")
			}
			results = append(results, model.TestResult{
				Fingerprint: model.FingerprintOf(file.Name, name),
				Suite:       file.Name,
				Name:        name,
				Outcome:     jestOutcome(a.Status),
				DurationMS:  a.Duration,
				Order:       order,
			})
			order++
		}
	}
	return results, nil
}

// jestOutcome maps a Jest status string to a normalized outcome.
func jestOutcome(status string) model.Outcome {
	switch status {
	case "passed":
		return model.OutcomePassed
	case "failed":
		return model.OutcomeFailed
	default:
		// pending, todo, skipped, disabled
		return model.OutcomeSkipped
	}
}
