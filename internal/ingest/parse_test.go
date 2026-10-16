package ingest

import (
	"strings"

	"github.com/gitstq/flakehunter/internal/model"
)

// parseString is a small helper for table-driven parser tests.
func parseString(s string, format Format) ([]model.TestResult, error) {
	return Parse(strings.NewReader(s), format)
}
