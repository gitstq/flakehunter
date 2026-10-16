package cmd

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/gitstq/flakehunter/internal/model"
)

// RecordCmd implements `flakehunter record`: parse test result files and
// append the run to the local history store.
func RecordCmd(args []string) error {
	fs := flag.NewFlagSet("record", flag.ContinueOnError)
	dir := fs.String("dir", defaultDir, "state directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	paths := fs.Args()
	if len(paths) == 0 {
		return fmt.Errorf("record: no test result files or directories given")
	}
	results, used, err := collect(paths)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return fmt.Errorf("record: no test results found in %s", strings.Join(paths, ", "))
	}
	run := model.Run{
		RunID:     nowRunID(),
		Timestamp: time.Now().UTC(),
		Source:    strings.Join(used, ", "),
		Commit:    commitFromEnv(),
		Results:   results,
	}
	store, _, err := openStore(*dir)
	if err != nil {
		return err
	}
	if err := store.Append(run); err != nil {
		return err
	}
	fmt.Printf("✓ Recorded run %s: %d results from %s\n", run.RunID, len(results), strings.Join(used, ", "))
	return nil
}
