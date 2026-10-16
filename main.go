// Command flakehunter is a zero-dependency, offline-first flaky test
// detective. It parses local test result files (JUnit XML, go test -json,
// TAP, Jest/Vitest JSON), scores flakiness with Wilson confidence bounds,
// and gates CI pipelines on flake thresholds.
package main

import (
	"fmt"
	"os"

	"github.com/gitstq/flakehunter/cmd"
)

func main() {
	if err := dispatch(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "flakehunter: %v\n", err)
		os.Exit(1)
	}
}

// dispatch routes the first argument to its subcommand.
func dispatch(args []string) error {
	if len(args) == 0 {
		cmd.Usage(os.Stdout)
		return nil
	}
	switch args[0] {
	case "init":
		return cmd.InitCmd(args[1:])
	case "scan":
		return cmd.ScanCmd(args[1:])
	case "record":
		return cmd.RecordCmd(args[1:])
	case "report":
		return cmd.ReportCmd(args[1:])
	case "gate":
		return cmd.GateCmd(args[1:])
	case "quarantine":
		return cmd.QuarantineCmd(args[1:])
	case "version", "--version", "-v":
		return cmd.VersionCmd(args[1:])
	case "help", "--help", "-h":
		cmd.Usage(os.Stdout)
		return nil
	default:
		cmd.Usage(os.Stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}
