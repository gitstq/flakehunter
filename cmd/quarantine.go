package cmd

import (
	"flag"
	"fmt"

	"github.com/gitstq/flakehunter/internal/engine"
	"github.com/gitstq/flakehunter/internal/quarantine"
)

// QuarantineCmd implements `flakehunter quarantine add|list|sync|patch`.
func QuarantineCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("quarantine: expected one of add|list|sync|patch")
	}
	switch args[0] {
	case "add":
		return quarantineAdd(args[1:])
	case "list":
		return quarantineList(args[1:])
	case "sync":
		return quarantineSync(args[1:])
	case "patch":
		return quarantinePatch(args[1:])
	default:
		return fmt.Errorf("quarantine: unknown subcommand %q (expected add|list|sync|patch)", args[0])
	}
}

// quarantineAdd implements `flakehunter quarantine add <fingerprint>...`.
func quarantineAdd(args []string) error {
	fs := flag.NewFlagSet("quarantine add", flag.ContinueOnError)
	var (
		dir       = fs.String("dir", defaultDir, "state directory")
		reason    = fs.String("reason", "flaky", "why the test is quarantined")
		framework = fs.String("framework", "pytest", "framework for patch hints: pytest|jest|vitest|go|junit|xunit")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) == 0 {
		return fmt.Errorf("quarantine add: missing test fingerprint")
	}
	man, err := quarantine.LoadManifest(*dir)
	if err != nil {
		return err
	}
	for _, fp := range rest {
		suite, name := splitFingerprint(fp)
		man.Add(fp, suite, name, *reason, *framework)
		fmt.Printf("🎲 quarantined %s\n", fp)
	}
	return man.Save(*dir)
}

// quarantineList implements `flakehunter quarantine list`.
func quarantineList(args []string) error {
	dir := defaultDir
	if len(args) > 0 {
		dir = args[0]
	}
	man, err := quarantine.LoadManifest(dir)
	if err != nil {
		return err
	}
	if len(man.Entries) == 0 {
		fmt.Println("no quarantined tests")
		return nil
	}
	fmt.Printf("🎲 %d quarantined test(s):\n", len(man.Entries))
	for _, e := range man.Entries {
		fmt.Printf("   %s — %s (framework: %s, added %s, stable runs: %d)\n",
			e.Fingerprint, e.Reason, e.Framework, e.AddedAt.Format("2006-01-02"), e.StableRuns)
	}
	return nil
}

// quarantineSync implements `flakehunter quarantine sync`: drop entries
// whose tests have been stable for enough consecutive runs.
func quarantineSync(args []string) error {
	fs := flag.NewFlagSet("quarantine sync", flag.ContinueOnError)
	dir := fs.String("dir", defaultDir, "state directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	store, cfg, err := openStore(*dir)
	if err != nil {
		return err
	}
	runs, err := store.Load()
	if err != nil {
		return err
	}
	if len(runs) == 0 {
		return fmt.Errorf("quarantine sync: no history in %s", *dir)
	}
	stats := engine.Analyze(runs, engine.Options{MinRuns: cfg.MinRuns})
	status := map[string]string{}
	for _, s := range stats {
		status[s.Fingerprint] = s.Status
	}
	man, err := quarantine.LoadManifest(*dir)
	if err != nil {
		return err
	}
	removed := man.Sync(status, cfg.QuarantineStableRuns)
	if err := man.Save(*dir); err != nil {
		return err
	}
	if len(removed) == 0 {
		fmt.Println("no quarantined test became stable — list unchanged")
		return nil
	}
	for _, e := range removed {
		fmt.Printf("✅ un-quarantined %s (stable for %d runs)\n", e.Fingerprint, cfg.QuarantineStableRuns)
	}
	return nil
}

// quarantinePatch implements `flakehunter quarantine patch <fingerprint>`.
func quarantinePatch(args []string) error {
	fs := flag.NewFlagSet("quarantine patch", flag.ContinueOnError)
	framework := fs.String("framework", "pytest", "target framework: pytest|jest|vitest|go|junit|xunit")
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) == 0 {
		return fmt.Errorf("quarantine patch: missing test fingerprint")
	}
	fp := rest[0]
	suite, name := splitFingerprint(fp)
	fmt.Println(quarantine.PatchSnippet(fp, suite, name, *framework))
	fmt.Println("// (drop the snippet above into the test file for " + fp + ")")
	return nil
}
