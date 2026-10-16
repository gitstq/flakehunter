package quarantine

import (
	"strings"
	"testing"
)

func TestPatchSnippetFrameworks(t *testing.T) {
	cases := []struct {
		framework string
		contains  string
	}{
		{"pytest", "@pytest.mark.skip"},
		{"jest", "it.skip"},
		{"vitest", "it.skip"},
		{"go", "t.Skip"},
		{"junit", "@Disabled"},
		{"xunit", "Skip ="},
		{"cargo", "Mark this test as skipped"},
	}
	for _, c := range cases {
		snippet := PatchSnippet("suite#test_thing", "suite", "test_thing", c.framework)
		if !strings.Contains(snippet, c.contains) {
			t.Errorf("framework %s: want snippet containing %q, got:\n%s", c.framework, c.contains, snippet)
		}
		if !strings.Contains(snippet, "suite#test_thing") {
			t.Errorf("framework %s: snippet should reference the fingerprint, got:\n%s", c.framework, snippet)
		}
	}
}

func TestPatchSnippetSanitizePyName(t *testing.T) {
	snippet := PatchSnippet("suite#test_fetch", "suite", "test_fetch[case-1]", "pytest")
	if !strings.Contains(snippet, "def test_fetch(") {
		t.Errorf("want sanitized python function name, got:\n%s", snippet)
	}
}
