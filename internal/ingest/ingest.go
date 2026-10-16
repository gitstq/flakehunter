// Package ingest parses test result files produced by common test frameworks
// and normalizes them into FlakeHunter's internal model.
package ingest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gitstq/flakehunter/internal/model"
)

// Format identifies a supported test-result file format.
type Format string

// Supported formats. FormatAuto is only used as a parse hint.
const (
	FormatJUnit  Format = "junit"
	FormatGoTest Format = "gotest"
	FormatTAP    Format = "tap"
	FormatJest   Format = "jest"
	FormatAuto   Format = "auto"
)

// tapPlanPattern matches the TAP plan line "1..N".
var tapPlanPattern = regexp.MustCompile(`^1\.\.\d+$`)

// Parser converts a test-result stream into normalized test results.
type Parser interface {
	Parse(r io.Reader) ([]model.TestResult, error)
}

// ParseFile sniffs the format of path and parses it in one pass.
func ParseFile(path string) ([]model.TestResult, Format, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	format := Detect(path, data)
	if format == "" {
		return nil, "", fmt.Errorf("unknown test result format: %s", path)
	}
	results, err := Parse(bytes.NewReader(data), format)
	if err != nil {
		return nil, format, fmt.Errorf("parse %s: %w", path, err)
	}
	return results, format, nil
}

// Parse reads test results in the given format.
func Parse(r io.Reader, format Format) ([]model.TestResult, error) {
	switch format {
	case FormatJUnit:
		return ParseJUnit(r)
	case FormatGoTest:
		return ParseGoTest(r)
	case FormatTAP:
		return ParseTAP(r)
	case FormatJest:
		return ParseJest(r)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Detect guesses the format of a file from its name and content.
func Detect(path string, data []byte) Format {
	base := strings.ToLower(filepath.Base(path))
	head := bytes.TrimSpace(data)
	if strings.HasSuffix(base, ".tap") {
		return FormatTAP
	}
	if bytes.HasPrefix(head, []byte("<?xml")) ||
		bytes.HasPrefix(head, []byte("<testsuite")) ||
		bytes.HasPrefix(head, []byte("<testsuites")) {
		return FormatJUnit
	}
	if len(head) > 0 && (head[0] == '{' || head[0] == '[') {
		// Whole-document JSON: Jest / Vitest reports.
		var probe map[string]json.RawMessage
		if err := json.Unmarshal(head, &probe); err == nil {
			if _, ok := probe["Action"]; ok {
				return FormatGoTest
			}
			if _, ok := probe["testResults"]; ok {
				return FormatJest
			}
			if _, ok := probe["numTotalTests"]; ok {
				return FormatJest
			}
		}
		// JSON-lines: `go test -json` writes one event per line.
		if idx := bytes.IndexByte(head, '\n'); idx > 0 {
			var lineProbe map[string]json.RawMessage
			if err := json.Unmarshal(bytes.TrimSpace(head[:idx]), &lineProbe); err == nil {
				if _, ok := lineProbe["Action"]; ok {
					return FormatGoTest
				}
			}
		}
	}
	line := firstLine(head)
	if strings.HasPrefix(line, "TAP version") ||
		strings.HasPrefix(line, "ok ") ||
		strings.HasPrefix(line, "not ok ") ||
		tapPlanPattern.MatchString(line) {
		return FormatTAP
	}
	return ""
}

// firstLine returns the first line of b without leading whitespace.
func firstLine(b []byte) string {
	if idx := bytes.IndexByte(b, '\n'); idx >= 0 {
		return string(bytes.TrimSpace(b[:idx]))
	}
	return string(b)
}

// WalkPaths collects candidate test-result files from the given files and directories.
func WalkPaths(paths []string) ([]string, error) {
	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			files = append(files, p)
			continue
		}
		err = filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			switch strings.ToLower(filepath.Ext(path)) {
			case ".xml", ".json", ".jsonl", ".tap":
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
