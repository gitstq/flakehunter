package quarantine

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gitstq/flakehunter/internal/model"
)

// PatchSnippet returns framework-specific code that quarantines the test.
// The snippet is a starting point: drop it into the test file at the
// fingerprint's location.
func PatchSnippet(fp, suite, name, framework string) string {
	switch framework {
	case "go":
		return fmt.Sprintf("// %s\n// Add inside the test function:\nt.Skip(\"flakehunter: quarantined flaky test\")", fp)
	case "pytest":
		return fmt.Sprintf("# %s\n@pytest.mark.skip(reason=\"flakehunter: quarantined flaky test\")\ndef %s(...):", fp, sanitizePyName(name))
	case "jest", "vitest":
		return fmt.Sprintf("// %s\nit.skip(\"%s\", () => {\n  // quarantined by flakehunter\n});", fp, name)
	case "junit", "maven":
		return fmt.Sprintf("// %s\nimport org.junit.jupiter.api.Disabled;\n\n@Disabled(\"flakehunter: quarantined flaky test\")", fp)
	case "xunit":
		return fmt.Sprintf("// %s\n[Fact(Skip = \"flakehunter: quarantined flaky test\")]", fp)
	default:
		return fmt.Sprintf("// %s (framework: %s)\n// Mark this test as skipped in your test framework.", fp, framework)
	}
}

// sanitizePyName derives a plausible Python function name from a test name.
func sanitizePyName(name string) string {
	cleaned := model.NormalizeName(name)
	var b strings.Builder
	for _, r := range cleaned {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	out := strings.Trim(b.String(), "_")
	if out == "" {
		return "test_quarantined"
	}
	if !strings.HasPrefix(out, "test") {
		out = "test_" + out
	}
	return out
}
