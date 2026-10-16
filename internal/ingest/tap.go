package ingest

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/gitstq/flakehunter/internal/model"
)

// tapLinePattern matches "ok 1 - description" / "not ok 2 - description"
// including optional TAP directives ("# SKIP", "# TODO").
var tapLinePattern = regexp.MustCompile(`(?i)^(not ok|ok)\b\s*(\d+)?\s*-?\s*(.*)$`)

// ParseTAP reads a Test Anything Protocol stream (TAP version 12/13).
func ParseTAP(r io.Reader) ([]model.TestResult, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	var results []model.TestResult
	order := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if tapPlanPattern.MatchString(line) {
			continue // plan line "1..N"
		}
		match := tapLinePattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		ok := strings.EqualFold(match[1], "ok")
		number := match[2]
		name := strings.TrimSpace(match[3])
		directive := ""
		if idx := strings.Index(name, "#"); idx >= 0 {
			directive = strings.ToUpper(strings.TrimSpace(name[idx+1:]))
			name = strings.TrimSpace(name[:idx])
		}
		if name == "" {
			name = "test-" + number
		}
		outcome := model.OutcomePassed
		if !ok {
			outcome = model.OutcomeFailed
		}
		switch {
		case strings.HasPrefix(directive, "SKIP"):
			outcome = model.OutcomeSkipped
		case strings.HasPrefix(directive, "TODO"):
			// Expected failure: not a real signal.
			outcome = model.OutcomeSkipped
		}
		results = append(results, model.TestResult{
			Fingerprint: model.FingerprintOf("tap", name),
			Suite:       "tap",
			Name:        name,
			Outcome:     outcome,
			Order:       order,
		})
		order++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
