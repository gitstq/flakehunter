package ingest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFixtures(t *testing.T) {
	cases := []struct {
		path string
		want Format
	}{
		{"../../testdata/junit/pytest-mixed.xml", FormatJUnit},
		{"../../testdata/gotest/events.jsonl", FormatGoTest},
		{"../../testdata/tap/basic.tap", FormatTAP},
		{"../../testdata/jest/report.json", FormatJest},
	}
	for _, c := range cases {
		data, err := os.ReadFile(c.path)
		if err != nil {
			t.Fatal(err)
		}
		got := Detect(filepath.Base(c.path), data)
		if got != c.want {
			t.Errorf("%s: want format %s, got %s", c.path, c.want, got)
		}
	}
}

func TestDetectUnknown(t *testing.T) {
	if got := Detect("notes.txt", []byte("just some notes, nothing to see")); got != "" {
		t.Errorf("want empty format, got %s", got)
	}
}

func TestParseFileRoundTrip(t *testing.T) {
	results, format, err := ParseFile("../../testdata/junit/pytest-mixed.xml")
	if err != nil {
		t.Fatal(err)
	}
	if format != FormatJUnit {
		t.Errorf("want junit, got %s", format)
	}
	if len(results) != 5 {
		t.Errorf("want 5 results, got %d", len(results))
	}
}

func TestWalkPaths(t *testing.T) {
	files, err := WalkPaths([]string{"../../testdata"})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Errorf("want 4 fixture files, got %d: %v", len(files), files)
	}
}
