package cmd

import (
	"fmt"
	"io"
)

// Usage prints the top-level help text.
func Usage(w io.Writer) {
	fmt.Fprint(w, `flakehunter — zero-dependency, offline-first flaky test detective 🎲

Usage:
  flakehunter <command> [options]

Commands:
  init [dir]                       Create the state directory and default config
  scan <paths...>                  One-off analysis of test result files
  record <paths...>                Append a run to local history
  report                           Flake statistics from recorded history
  gate                             CI quality gate based on history (exit 1 on violation)
  quarantine add|list|sync|patch   Manage temporarily excluded flaky tests
  version                          Print version information

Supported inputs:
  JUnit XML (pytest, Jest, Maven Surefire, Gradle, TestNG, xUnit, RSpec ...)
  go test -json        TAP (TAP version 12/13)        Jest/Vitest --json

Examples:
  flakehunter scan ./test-results
  flakehunter init
  flakehunter scan ./test-results --record
  flakehunter report --out markdown --summary "$GITHUB_STEP_SUMMARY"
  flakehunter gate --max-flake-rate 5
  flakehunter quarantine add "tests.test_net#test_fetch" --framework pytest

Docs: https://github.com/gitstq/flakehunter
`)
}
