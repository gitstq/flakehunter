package cmd

import (
	"fmt"

	"github.com/gitstq/flakehunter/internal/engine"
)

// InitCmd implements `flakehunter init`.
func InitCmd(args []string) error {
	dir := defaultDir
	if len(args) > 0 {
		dir = args[0]
	}
	if err := engine.SaveConfig(dir, engine.DefaultConfig()); err != nil {
		return err
	}
	fmt.Printf("✓ Initialized FlakeHunter state in %s\n", dir)
	fmt.Println("  Keep the state local by adding this to .gitignore:")
	fmt.Printf("  %s/\n", dir)
	return nil
}
