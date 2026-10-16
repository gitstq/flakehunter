// Package model defines the core data types shared across FlakeHunter.
package model

import (
	"strings"
	"time"
)

// Outcome is the result of a single test execution.
type Outcome string

// Supported test outcomes.
const (
	OutcomePassed  Outcome = "passed"
	OutcomeFailed  Outcome = "failed"
	OutcomeSkipped Outcome = "skipped"
	OutcomeError   Outcome = "error"
)

// TestResult is one test execution observed in a single run.
type TestResult struct {
	Fingerprint string  `json:"fingerprint"`
	Suite       string  `json:"suite,omitempty"`
	Name        string  `json:"name"`
	File        string  `json:"file,omitempty"`
	Outcome     Outcome `json:"outcome"`
	DurationMS  float64 `json:"duration_ms,omitempty"`
	Order       int     `json:"order"`
	Retried     bool    `json:"retried,omitempty"`
	RerunPass   bool    `json:"rerun_pass,omitempty"`
}

// Run is the collection of test results produced by one CI run or local invocation.
type Run struct {
	RunID     string       `json:"run_id"`
	Timestamp time.Time    `json:"timestamp"`
	Source    string       `json:"source"`
	Commit    string       `json:"commit,omitempty"`
	Results   []TestResult `json:"results"`
}

// NormalizeName strips parametrized suffixes such as "[case-1]" so that
// parametrized variants of the same test share a single fingerprint.
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	for {
		open := strings.Index(name, "[")
		if open < 0 {
			return strings.TrimSpace(name)
		}
		closeIdx := strings.Index(name[open:], "]")
		if closeIdx < 0 {
			return strings.TrimSpace(name)
		}
		name = name[:open] + name[open+closeIdx+1:]
	}
}

// NormalizeSuite collapses whitespace inside a suite name.
func NormalizeSuite(suite string) string {
	return strings.Join(strings.Fields(suite), " ")
}

// FingerprintOf builds the stable identity of a test:
// normalized suite + "#" + normalized name.
func FingerprintOf(suite, name string) string {
	return NormalizeSuite(suite) + "#" + NormalizeName(name)
}
