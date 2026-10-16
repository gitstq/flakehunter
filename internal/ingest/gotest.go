package ingest

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/gitstq/flakehunter/internal/model"
)

// goTestEvent is one line of `go test -json` output.
type goTestEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
}

// ParseGoTest reads the JSON stream produced by `go test -json ./...`.
// Only final per-test actions (pass/fail/skip) are recorded; with -count>1
// the same test may appear several times, which FlakeHunter treats as a
// same-run flip signal.
func ParseGoTest(r io.Reader) ([]model.TestResult, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	var results []model.TestResult
	order := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var ev goTestEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue // build output or non-JSON noise
		}
		if ev.Test == "" {
			continue // package-level event
		}
		outcome, final := goTestOutcome(ev.Action)
		if !final {
			continue
		}
		results = append(results, model.TestResult{
			Fingerprint: model.FingerprintOf(ev.Package, ev.Test),
			Suite:       ev.Package,
			Name:        ev.Test,
			Outcome:     outcome,
			DurationMS:  ev.Elapsed * 1000,
			Order:       order,
		})
		order++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// goTestOutcome maps a go test action to an outcome and reports whether it is final.
func goTestOutcome(action string) (model.Outcome, bool) {
	switch action {
	case "pass":
		return model.OutcomePassed, true
	case "fail":
		return model.OutcomeFailed, true
	case "skip":
		return model.OutcomeSkipped, true
	default:
		return "", false
	}
}
